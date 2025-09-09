package execution

import (
	"context"
	"fmt"

	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// Upsert inserts new records or updates existing ones based on conflict resolution.
func (e *QueryExecutor) Upsert(ctx context.Context, qb QueryBuilderInterface, values []map[string]interface{}, options types.UpsertOptions) error {
	if len(values) == 0 {
		return fmt.Errorf("no values provided for upsert")
	}

	table := qb.GetTable()
	if table == "" {
		return fmt.Errorf("no table specified for upsert")
	}

	switch e.driver {
	case types.MySQL:
		return e.upsertMySQL(ctx, table, values, options)
	case types.PostgreSQL:
		return e.upsertPostgreSQL(ctx, table, values, options)
	default:
		return fmt.Errorf("upsert not supported for driver: %s", e.driver)
	}
}

func (e *QueryExecutor) upsertMySQL(ctx context.Context, table string, values []map[string]interface{}, options types.UpsertOptions) error {
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
			rowPlaceholders = append(rowPlaceholders, "?")
		}
		
		allBindings = append(allBindings, rowBindings...)
		valueSets = append(valueSets, "("+joinStrings(rowPlaceholders, ", ")+")")
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		table,
		joinColumns(columns),
		joinStrings(valueSets, ", "))

	updateColumns := options.UpdateColumns
	if len(updateColumns) == 0 {
		updateColumns = columns
	}

	var updateParts []string
	for _, column := range updateColumns {
		updateParts = append(updateParts, fmt.Sprintf("%s = VALUES(%s)", column, column))
	}

	sql += " ON DUPLICATE KEY UPDATE " + joinStrings(updateParts, ", ")

	_, err := e.executor.ExecContext(ctx, sql, allBindings...)
	if err != nil {
		return fmt.Errorf("failed to execute MySQL upsert: %w", err)
	}

	return nil
}

func (e *QueryExecutor) upsertPostgreSQL(ctx context.Context, table string, values []map[string]interface{}, options types.UpsertOptions) error {
	firstRow := values[0]
	columns := make([]string, 0, len(firstRow))
	for column := range firstRow {
		columns = append(columns, column)
	}

	var allBindings []interface{}
	var valueSets []string
	bindingPos := 1

	for _, row := range values {
		var rowBindings []interface{}
		var rowPlaceholders []string
		
		for _, column := range columns {
			value, exists := row[column]
			if !exists {
				value = nil
			}
			rowBindings = append(rowBindings, value)
			rowPlaceholders = append(rowPlaceholders, fmt.Sprintf("$%d", bindingPos))
			bindingPos++
		}
		
		allBindings = append(allBindings, rowBindings...)
		valueSets = append(valueSets, "("+joinStrings(rowPlaceholders, ", ")+")")
	}

	sql := fmt.Sprintf("INSERT INTO %s (%s) VALUES %s",
		table,
		joinColumns(columns),
		joinStrings(valueSets, ", "))

	conflictTarget := options.ConflictTarget
	if len(conflictTarget) == 0 {
		return fmt.Errorf("conflict target must be specified for PostgreSQL upsert")
	}

	sql += fmt.Sprintf(" ON CONFLICT (%s)", joinColumns(conflictTarget))

	switch options.ConflictAction {
	case types.DoNothing:
		sql += " DO NOTHING"
	case types.DoUpdate:
		updateColumns := options.UpdateColumns
		if len(updateColumns) == 0 {
			updateColumns = make([]string, 0)
			for _, column := range columns {
				found := false
				for _, conflict := range conflictTarget {
					if column == conflict {
						found = true
						break
					}
				}
				if !found {
					updateColumns = append(updateColumns, column)
				}
			}
		}

		if len(updateColumns) > 0 {
			var updateParts []string
			for _, column := range updateColumns {
				updateParts = append(updateParts, fmt.Sprintf("%s = EXCLUDED.%s", column, column))
			}
			sql += " DO UPDATE SET " + joinStrings(updateParts, ", ")
		} else {
			sql += " DO NOTHING"
		}
	default:
		sql += " DO NOTHING"
	}

	_, err := e.executor.ExecContext(ctx, sql, allBindings...)
	if err != nil {
		return fmt.Errorf("failed to execute PostgreSQL upsert: %w", err)
	}

	return nil
}

// InsertOrIgnore inserts records but ignores duplicates without raising an error.
func (e *QueryExecutor) InsertOrIgnore(ctx context.Context, qb QueryBuilderInterface, values interface{}) error {
	table := qb.GetTable()
	if table == "" {
		return fmt.Errorf("no table specified for insert or ignore")
	}

	switch v := values.(type) {
	case map[string]interface{}:
		return e.insertOrIgnoreSingle(ctx, table, v)
	case []map[string]interface{}:
		return e.insertOrIgnoreBatch(ctx, table, v)
	default:
		return fmt.Errorf("invalid values type for insert or ignore")
	}
}

func (e *QueryExecutor) insertOrIgnoreSingle(ctx context.Context, table string, values map[string]interface{}) error {
	if len(values) == 0 {
		return fmt.Errorf("no values provided for insert or ignore")
	}

	columns := make([]string, 0, len(values))
	bindings := make([]interface{}, 0, len(values))
	placeholders := make([]string, 0, len(values))

	for column, value := range values {
		columns = append(columns, column)
		bindings = append(bindings, value)
		placeholders = append(placeholders, e.getPlaceholder(len(bindings)))
	}

	var sql string
	switch e.driver {
	case types.MySQL:
		sql = fmt.Sprintf("INSERT IGNORE INTO %s (%s) VALUES (%s)",
			table,
			joinColumns(columns),
			joinStrings(placeholders, ", "))
	case types.PostgreSQL:
		sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) ON CONFLICT DO NOTHING",
			table,
			joinColumns(columns),
			joinStrings(placeholders, ", "))
	default:
		return fmt.Errorf("insert or ignore not supported for driver: %s", e.driver)
	}

	_, err := e.executor.ExecContext(ctx, sql, bindings...)
	if err != nil {
		return fmt.Errorf("failed to execute insert or ignore: %w", err)
	}

	return nil
}

func (e *QueryExecutor) insertOrIgnoreBatch(ctx context.Context, table string, values []map[string]interface{}) error {
	if len(values) == 0 {
		return fmt.Errorf("no values provided for batch insert or ignore")
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

	var sql string
	switch e.driver {
	case types.MySQL:
		sql = fmt.Sprintf("INSERT IGNORE INTO %s (%s) VALUES %s",
			table,
			joinColumns(columns),
			joinStrings(valueSets, ", "))
	case types.PostgreSQL:
		sql = fmt.Sprintf("INSERT INTO %s (%s) VALUES %s ON CONFLICT DO NOTHING",
			table,
			joinColumns(columns),
			joinStrings(valueSets, ", "))
	default:
		return fmt.Errorf("insert or ignore not supported for driver: %s", e.driver)
	}

	_, err := e.executor.ExecContext(ctx, sql, allBindings...)
	if err != nil {
		return fmt.Errorf("failed to execute batch insert or ignore: %w", err)
	}

	return nil
}

// Replace replaces records using MySQL's REPLACE INTO statement.
func (e *QueryExecutor) Replace(ctx context.Context, qb QueryBuilderInterface, values map[string]interface{}) error {
	if e.driver != types.MySQL {
		return fmt.Errorf("replace is only supported for MySQL")
	}

	table := qb.GetTable()
	if table == "" {
		return fmt.Errorf("no table specified for replace")
	}

	if len(values) == 0 {
		return fmt.Errorf("no values provided for replace")
	}

	columns := make([]string, 0, len(values))
	bindings := make([]interface{}, 0, len(values))
	placeholders := make([]string, 0, len(values))

	for column, value := range values {
		columns = append(columns, column)
		bindings = append(bindings, value)
		placeholders = append(placeholders, "?")
	}

	sql := fmt.Sprintf("REPLACE INTO %s (%s) VALUES (%s)",
		table,
		joinColumns(columns),
		joinStrings(placeholders, ", "))

	_, err := e.executor.ExecContext(ctx, sql, bindings...)
	if err != nil {
		return fmt.Errorf("failed to execute replace: %w", err)
	}

	return nil
}