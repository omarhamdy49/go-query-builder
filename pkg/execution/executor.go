package execution

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// QueryExecutor handles the execution of database queries built by query builders.
type QueryExecutor struct {
	executor types.QueryExecutor
	driver   types.Driver
}

// NewQueryExecutor creates a new QueryExecutor with the specified executor and driver.
func NewQueryExecutor(executor types.QueryExecutor, driver types.Driver) *QueryExecutor {
	return &QueryExecutor{
		executor: executor,
		driver:   driver,
	}
}

// QueryBuilderInterface defines the methods required by query builders for execution.
type QueryBuilderInterface interface {
	ToSQL() (string, []interface{}, error)
	Clone() types.QueryBuilder
	GetTable() string
}

// Get executes the query and returns all matching rows as a collection.
func (e *QueryExecutor) Get(ctx context.Context, qb QueryBuilderInterface) (types.Collection, error) {
	sql, bindings, err := qb.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build SQL: %w", err)
	}

	rows, err := e.executor.QueryContext(ctx, sql, bindings...)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer func() { _ = rows.Close() }()

	return e.scanRows(rows)
}

// First executes the query and returns the first matching row.
func (e *QueryExecutor) First(ctx context.Context, qb QueryBuilderInterface) (map[string]interface{}, error) {
	clone := qb.Clone()
	limitedQB := clone.Limit(1)
	
	collection, err := e.Get(ctx, limitedQB.(QueryBuilderInterface))
	if err != nil {
		return nil, err
	}

	if collection.IsEmpty() {
		return nil, sql.ErrNoRows
	}

	return collection.First(), nil
}

// Find finds a record by its primary key ID.
func (e *QueryExecutor) Find(ctx context.Context, qb QueryBuilderInterface, id interface{}) (map[string]interface{}, error) {
	clone := qb.Clone()
	findQB := clone.Where("id", id).Limit(1)
	
	collection, err := e.Get(ctx, findQB.(QueryBuilderInterface))
	if err != nil {
		return nil, err
	}

	if collection.IsEmpty() {
		return nil, sql.ErrNoRows
	}

	return collection.First(), nil
}

// Pluck returns all values from a single column as a slice.
func (e *QueryExecutor) Pluck(ctx context.Context, qb QueryBuilderInterface, column string) ([]interface{}, error) {
	clone := qb.Clone()
	pluckQB := clone.Select(column)
	
	collection, err := e.Get(ctx, pluckQB.(QueryBuilderInterface))
	if err != nil {
		return nil, err
	}

	return collection.Pluck(column), nil
}

// Count returns the number of rows that match the query conditions.
func (e *QueryExecutor) Count(ctx context.Context, qb QueryBuilderInterface) (int64, error) {
	result, err := e.aggregate(ctx, qb, types.Count, "*")
	if err != nil {
		return 0, err
	}
	
	if count, ok := result.(int64); ok {
		return count, nil
	}
	
	// Handle different numeric types that might be returned
	switch v := result.(type) {
	case int:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case []uint8:
		// MySQL driver sometimes returns numbers as byte slices
		if str := string(v); str != "" {
			if count, err := strconv.ParseInt(str, 10, 64); err == nil {
				return count, nil
			}
		}
		return 0, fmt.Errorf("failed to parse count from byte slice: %s", string(v))
	case string:
		// Handle string representations
		if count, err := strconv.ParseInt(v, 10, 64); err == nil {
			return count, nil
		}
		return 0, fmt.Errorf("failed to parse count from string: %s", v)
	default:
		return 0, fmt.Errorf("unexpected count result type: %T", result)
	}
}

// Sum returns the sum of values in the specified column.
func (e *QueryExecutor) Sum(ctx context.Context, qb QueryBuilderInterface, column string) (interface{}, error) {
	return e.aggregate(ctx, qb, types.Sum, column)
}

// Avg returns the average value of the specified column.
func (e *QueryExecutor) Avg(ctx context.Context, qb QueryBuilderInterface, column string) (interface{}, error) {
	return e.aggregate(ctx, qb, types.Avg, column)
}

// Min returns the minimum value of the specified column.
func (e *QueryExecutor) Min(ctx context.Context, qb QueryBuilderInterface, column string) (interface{}, error) {
	return e.aggregate(ctx, qb, types.Min, column)
}

// Max returns the maximum value of the specified column.
func (e *QueryExecutor) Max(ctx context.Context, qb QueryBuilderInterface, column string) (interface{}, error) {
	return e.aggregate(ctx, qb, types.Max, column)
}

func (e *QueryExecutor) aggregate(ctx context.Context, qb QueryBuilderInterface, fn types.AggregateFunction, column string) (interface{}, error) {
	clone := qb.Clone()
	aggregateQB := clone.SelectRaw(fmt.Sprintf("%s(%s) as aggregate", fn, column))
	
	sql, bindings, err := aggregateQB.ToSQL()
	if err != nil {
		return nil, fmt.Errorf("failed to build aggregate SQL: %w", err)
	}

	row := e.executor.QueryRowContext(ctx, sql, bindings...)
	
	var result interface{}
	if err := row.Scan(&result); err != nil {
		return nil, fmt.Errorf("failed to scan aggregate result: %w", err)
	}

	return result, nil
}

// Insert executes an INSERT statement with the provided values.
func (e *QueryExecutor) Insert(ctx context.Context, qb QueryBuilderInterface, values map[string]interface{}) error {
	if len(values) == 0 {
		return fmt.Errorf("no values provided for insert")
	}

	table := qb.GetTable()
	if table == "" {
		return fmt.Errorf("no table specified for insert")
	}

	columns := make([]string, 0, len(values))
	bindings := make([]interface{}, 0, len(values))
	placeholders := make([]string, 0, len(values))

	for column, value := range values {
		columns = append(columns, column)
		bindings = append(bindings, value)
		placeholders = append(placeholders, e.getPlaceholder(len(bindings)))
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		table,
		joinColumns(columns),
		joinStrings(placeholders, ", "))

	_, err := e.executor.ExecContext(ctx, sql, bindings...)
	if err != nil {
		return fmt.Errorf("failed to execute insert: %w", err)
	}

	return nil
}

// InsertBatch executes a batch INSERT statement with multiple rows of values.
func (e *QueryExecutor) InsertBatch(ctx context.Context, qb QueryBuilderInterface, values []map[string]interface{}) error {
	if len(values) == 0 {
		return fmt.Errorf("no values provided for batch insert")
	}

	table := qb.GetTable()
	if table == "" {
		return fmt.Errorf("no table specified for insert")
	}

	firstRow := values[0]
	columns := make([]string, 0, len(firstRow))
	for column := range firstRow {
		columns = append(columns, column)
	}

	var allBindings []interface{}
	var valueSets []string

	for _, row := range values {
		var rowBindings []interface{}
		var rowPlaceholders []string
		
		for _, column := range columns {
			value, exists := row[column]
			if !exists {
				value = nil
			}
			rowBindings = append(rowBindings, value)
			rowPlaceholders = append(rowPlaceholders, e.getPlaceholder(len(allBindings)+len(rowBindings)))
		}
		
		allBindings = append(allBindings, rowBindings...)
		valueSets = append(valueSets, "("+joinStrings(rowPlaceholders, ", ")+")")
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		table,
		joinColumns(columns),
		joinStrings(valueSets, ", "))

	_, err := e.executor.ExecContext(ctx, sql, allBindings...)
	if err != nil {
		return fmt.Errorf("failed to execute batch insert: %w", err)
	}

	return nil
}

// Update executes an UPDATE statement and returns the number of affected rows.
func (e *QueryExecutor) Update(ctx context.Context, qb QueryBuilderInterface, values map[string]interface{}) (int64, error) {
	if len(values) == 0 {
		return 0, fmt.Errorf("no values provided for update")
	}

	table := qb.GetTable()
	if table == "" {
		return 0, fmt.Errorf("no table specified for update")
	}

	setParts := make([]string, 0, len(values))
	bindings := make([]interface{}, 0, len(values))

	for column, value := range values {
		setParts = append(setParts, fmt.Sprintf("%s = %s", column, e.getPlaceholder(len(bindings)+1)))
		bindings = append(bindings, value)
	}

	sql := fmt.Sprintf("UPDATE %s SET %s", table, joinStrings(setParts, ", "))

	whereSQL, whereBindings, err := e.buildWhereClause(qb)
	if err != nil {
		return 0, fmt.Errorf("failed to build where clause: %w", err)
	}

	if whereSQL != "" {
		sql += " WHERE " + whereSQL
		bindings = append(bindings, whereBindings...)
	}

	result, err := e.executor.ExecContext(ctx, sql, bindings...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

// Delete executes a DELETE statement and returns the number of affected rows.
func (e *QueryExecutor) Delete(ctx context.Context, qb QueryBuilderInterface) (int64, error) {
	table := qb.GetTable()
	if table == "" {
		return 0, fmt.Errorf("no table specified for delete")
	}

	sql := fmt.Sprintf("DELETE FROM %s", table)
	var bindings []interface{}

	whereSQL, whereBindings, err := e.buildWhereClause(qb)
	if err != nil {
		return 0, fmt.Errorf("failed to build where clause: %w", err)
	}

	if whereSQL != "" {
		sql += " WHERE " + whereSQL
		bindings = append(bindings, whereBindings...)
	}

	result, err := e.executor.ExecContext(ctx, sql, bindings...)
	if err != nil {
		return 0, fmt.Errorf("failed to execute delete: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (e *QueryExecutor) scanRows(rows types.Rows) (types.Collection, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to get columns: %w", err)
	}

	var results []map[string]interface{}
	
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		
		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			if b, ok := val.([]byte); ok {
				row[col] = string(b)
			} else {
				row[col] = val
			}
		}
		results = append(results, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return types.NewCollection(results), nil
}

func (e *QueryExecutor) buildWhereClause(qb QueryBuilderInterface) (string, []interface{}, error) {
	sql, bindings, err := qb.ToSQL()
	if err != nil {
		return "", nil, err
	}
	
	whereStart := findWhereClause(sql)
	if whereStart == -1 {
		return "", nil, nil
	}
	
	whereClause := sql[whereStart+6:]
	if orderBy := findOrderByClause(whereClause); orderBy != -1 {
		whereClause = whereClause[:orderBy]
	}
	if groupBy := findGroupByClause(whereClause); groupBy != -1 {
		whereClause = whereClause[:groupBy]
	}
	
	return whereClause, bindings, nil
}

func (e *QueryExecutor) getPlaceholder(position int) string {
	switch e.driver {
	case types.PostgreSQL:
		return fmt.Sprintf("$%d", position)
	default:
		return "?"
	}
}

func findWhereClause(sql string) int {
	return findKeyword(sql, "WHERE")
}

func findOrderByClause(sql string) int {
	return findKeyword(sql, "ORDER BY")
}

func findGroupByClause(sql string) int {
	return findKeyword(sql, "GROUP BY")
}

func findKeyword(sql, keyword string) int {
	upperSQL := toUpper(sql)
	return indexOf(upperSQL, keyword)
}

func joinColumns(columns []string) string {
	return joinStrings(columns, ", ")
}

func joinStrings(strs []string, _sep string) string {
	separator := ", "
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += separator + strs[i]
	}
	return result
}

func toUpper(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'a' && c <= 'z' {
			c = c - 'a' + 'A'
		}
		result[i] = c
	}
	return string(result)
}

func indexOf(s, substr string) int {
	if len(substr) == 0 {
		return 0
	}
	if len(substr) > len(s) {
		return -1
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}