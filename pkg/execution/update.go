package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

func (e *QueryExecutor) UpdateJson(ctx context.Context, qb QueryBuilderInterface, column string, path string, value interface{}) (int64, error) {
	table := qb.GetTable()
	if table == "" {
		return 0, fmt.Errorf("no table specified for update")
	}

	var sql string
	var bindings []interface{}

	switch e.driver {
	case types.MySQL:
		sql = fmt.Sprintf("UPDATE %s SET %s = JSON_SET(%s, %s, %s)", 
			table, column, column, e.getPlaceholder(1), e.getPlaceholder(2))
		bindings = append(bindings, path, value)
	case types.PostgreSQL:
		sql = fmt.Sprintf("UPDATE %s SET %s = jsonb_set(%s, %s, %s)", 
			table, column, column, e.getPlaceholder(1), e.getPlaceholder(2))
		pathArray := strings.Split(strings.Trim(path, "$"), ".")
		jsonValue, _ := json.Marshal(value)
		bindings = append(bindings, fmt.Sprintf("{%s}", strings.Join(pathArray, ",")), string(jsonValue))
	default:
		return 0, fmt.Errorf("JSON updates not supported for driver: %s", e.driver)
	}

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
		return 0, fmt.Errorf("failed to execute JSON update: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (e *QueryExecutor) UpdateJsonRemove(ctx context.Context, qb QueryBuilderInterface, column string, path string) (int64, error) {
	table := qb.GetTable()
	if table == "" {
		return 0, fmt.Errorf("no table specified for update")
	}

	var sql string
	var bindings []interface{}

	switch e.driver {
	case types.MySQL:
		sql = fmt.Sprintf("UPDATE %s SET %s = JSON_REMOVE(%s, %s)", 
			table, column, column, e.getPlaceholder(1))
		bindings = append(bindings, path)
	case types.PostgreSQL:
		sql = fmt.Sprintf("UPDATE %s SET %s = %s - %s", 
			table, column, column, e.getPlaceholder(1))
		pathArray := strings.Split(strings.Trim(path, "$"), ".")
		bindings = append(bindings, fmt.Sprintf("{%s}", strings.Join(pathArray, ",")))
	default:
		return 0, fmt.Errorf("JSON updates not supported for driver: %s", e.driver)
	}

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
		return 0, fmt.Errorf("failed to execute JSON remove: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (e *QueryExecutor) Increment(ctx context.Context, qb QueryBuilderInterface, column string, value ...interface{}) (int64, error) {
	amount := 1
	if len(value) > 0 {
		if v, ok := value[0].(int); ok {
			amount = v
		}
	}

	table := qb.GetTable()
	if table == "" {
		return 0, fmt.Errorf("no table specified for increment")
	}

	sql := fmt.Sprintf("UPDATE %s SET %s = %s + %s", 
		table, column, column, e.getPlaceholder(1))
	bindings := []interface{}{amount}

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
		return 0, fmt.Errorf("failed to execute increment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (e *QueryExecutor) Decrement(ctx context.Context, qb QueryBuilderInterface, column string, value ...interface{}) (int64, error) {
	amount := 1
	if len(value) > 0 {
		if v, ok := value[0].(int); ok {
			amount = v
		}
	}

	table := qb.GetTable()
	if table == "" {
		return 0, fmt.Errorf("no table specified for decrement")
	}

	sql := fmt.Sprintf("UPDATE %s SET %s = %s - %s", 
		table, column, column, e.getPlaceholder(1))
	bindings := []interface{}{amount}

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
		return 0, fmt.Errorf("failed to execute decrement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("failed to get rows affected: %w", err)
	}

	return rowsAffected, nil
}

func (e *QueryExecutor) UpdateOrCreate(ctx context.Context, qb QueryBuilderInterface, attributes map[string]interface{}, values map[string]interface{}) (map[string]interface{}, error) {
	clone := qb.Clone()
	for column, value := range attributes {
		clone = clone.Where(column, value)
	}

	record, err := e.First(ctx, clone.(QueryBuilderInterface))
	if err == nil {
		mergedValues := make(map[string]interface{})
		for k, v := range values {
			mergedValues[k] = v
		}
		
		_, updateErr := e.Update(ctx, clone.(QueryBuilderInterface), mergedValues)
		if updateErr != nil {
			return nil, fmt.Errorf("failed to update record: %w", updateErr)
		}
		
		for k, v := range mergedValues {
			record[k] = v
		}
		return record, nil
	}

	insertValues := make(map[string]interface{})
	for k, v := range attributes {
		insertValues[k] = v
	}
	for k, v := range values {
		insertValues[k] = v
	}

	insertErr := e.Insert(ctx, qb, insertValues)
	if insertErr != nil {
		return nil, fmt.Errorf("failed to create record: %w", insertErr)
	}

	return insertValues, nil
}

func (e *QueryExecutor) UpdateOrInsert(ctx context.Context, qb QueryBuilderInterface, attributes map[string]interface{}, values map[string]interface{}) error {
	clone := qb.Clone()
	for column, value := range attributes {
		clone = clone.Where(column, value)
	}

	_, err := e.First(ctx, clone.(QueryBuilderInterface))
	if err == nil {
		mergedValues := make(map[string]interface{})
		for k, v := range values {
			mergedValues[k] = v
		}
		
		_, updateErr := e.Update(ctx, clone.(QueryBuilderInterface), mergedValues)
		return updateErr
	}

	insertValues := make(map[string]interface{})
	for k, v := range attributes {
		insertValues[k] = v
	}
	for k, v := range values {
		insertValues[k] = v
	}

	return e.Insert(ctx, qb, insertValues)
}