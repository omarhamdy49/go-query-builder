package query

import (
	"fmt"
	"strings"
	"time"

	"github.com/omarhamdy49/go-query-builder/pkg/clauses"
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// WhereDate adds a WHERE clause that matches a date portion of a datetime column.
func (qb *Builder) WhereDate(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "date", types.And, args...)
}

// OrWhereDate adds an OR WHERE clause that matches a date portion of a datetime column.
func (qb *Builder) OrWhereDate(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "date", types.Or, args...)
}

// WhereTime adds a WHERE clause that matches a time portion of a datetime column.
func (qb *Builder) WhereTime(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "time", types.And, args...)
}

// OrWhereTime adds an OR WHERE clause that matches a time portion of a datetime column.
func (qb *Builder) OrWhereTime(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "time", types.Or, args...)
}

// WhereDay adds a WHERE clause that matches the day of the month from a datetime column.
func (qb *Builder) WhereDay(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "day", types.And, args...)
}

// OrWhereDay adds an OR WHERE clause that matches the day of the month from a datetime column.
func (qb *Builder) OrWhereDay(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "day", types.Or, args...)
}

// WhereMonth adds a WHERE clause that matches the month from a datetime column.
func (qb *Builder) WhereMonth(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "month", types.And, args...)
}

// OrWhereMonth adds an OR WHERE clause that matches the month from a datetime column.
func (qb *Builder) OrWhereMonth(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "month", types.Or, args...)
}

// WhereYear adds a WHERE clause that matches the year from a datetime column.
func (qb *Builder) WhereYear(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "year", types.And, args...)
}

// OrWhereYear adds an OR WHERE clause that matches the year from a datetime column.
func (qb *Builder) OrWhereYear(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "year", types.Or, args...)
}

// WherePast adds a WHERE clause that matches dates in the past.
func (qb *Builder) WherePast(column string) types.QueryBuilder {
	return qb.Where(column, "<", time.Now())
}

// WhereFuture adds a WHERE clause that matches dates in the future.
func (qb *Builder) WhereFuture(column string) types.QueryBuilder {
	return qb.Where(column, ">", time.Now())
}

// WhereNowOrPast adds a WHERE clause that matches dates in the past or now.
func (qb *Builder) WhereNowOrPast(column string) types.QueryBuilder {
	return qb.Where(column, "<=", time.Now())
}

// WhereNowOrFuture adds a WHERE clause that matches dates in the future or now.
func (qb *Builder) WhereNowOrFuture(column string) types.QueryBuilder {
	return qb.Where(column, ">=", time.Now())
}

// WhereToday adds a WHERE clause that matches today's date.
func (qb *Builder) WhereToday(column string) types.QueryBuilder {
	return qb.WhereDate(column, time.Now().Format("2006-01-02"))
}

// WhereBeforeToday adds a WHERE clause that matches dates before today.
func (qb *Builder) WhereBeforeToday(column string) types.QueryBuilder {
	return qb.WhereDate(column, "<", time.Now().Format("2006-01-02"))
}

// WhereAfterToday adds a WHERE clause that matches dates after today.
func (qb *Builder) WhereAfterToday(column string) types.QueryBuilder {
	return qb.WhereDate(column, ">", time.Now().Format("2006-01-02"))
}

// WhereTodayOrBefore adds a WHERE clause that matches today's date or earlier.
func (qb *Builder) WhereTodayOrBefore(column string) types.QueryBuilder {
	return qb.WhereDate(column, "<=", time.Now().Format("2006-01-02"))
}

// WhereTodayOrAfter adds a WHERE clause that matches today's date or later.
func (qb *Builder) WhereTodayOrAfter(column string) types.QueryBuilder {
	return qb.WhereDate(column, ">=", time.Now().Format("2006-01-02"))
}

func (qb *Builder) addDateWhere(column, dateType string, boolean types.BooleanOperator, args ...interface{}) types.QueryBuilder {
	if len(args) == 0 {
		return qb
	}

	var operator types.Operator
	var value interface{}

	switch len(args) {
	case 1:
		operator = types.OpEqual
		value = args[0]
	case 2:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	default:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	}

	var raw string
	switch qb.driver {
	case types.MySQL:
		switch dateType {
		case "date":
			raw = fmt.Sprintf("DATE(%s) %s ?", column, operator)
		case "time":
			raw = fmt.Sprintf("TIME(%s) %s ?", column, operator)
		case "day":
			raw = fmt.Sprintf("DAY(%s) %s ?", column, operator)
		case "month":
			raw = fmt.Sprintf("MONTH(%s) %s ?", column, operator)
		case "year":
			raw = fmt.Sprintf("YEAR(%s) %s ?", column, operator)
		}
	case types.PostgreSQL:
		switch dateType {
		case "date":
			raw = fmt.Sprintf("DATE(%s) %s ?", column, operator)
		case "time":
			raw = fmt.Sprintf("EXTRACT(HOUR FROM %s) %s ?", column, operator)
		case "day":
			raw = fmt.Sprintf("EXTRACT(DAY FROM %s) %s ?", column, operator)
		case "month":
			raw = fmt.Sprintf("EXTRACT(MONTH FROM %s) %s ?", column, operator)
		case "year":
			raw = fmt.Sprintf("EXTRACT(YEAR FROM %s) %s ?", column, operator)
		}
	}

	clause := clauses.NewWhereRawClause(raw)
	clause.SetBoolean(boolean)
	qb.wheres = append(qb.wheres, clause)
	qb.bindings = append(qb.bindings, value)
	return qb
}

// WhereJSONContains adds a WHERE clause for JSON containment checks.
func (qb *Builder) WhereJSONContains(column string, value interface{}) types.QueryBuilder {
	clause := clauses.NewWhereJSONContainsClause(column, value)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

// OrWhereJSONContains adds an OR WHERE clause for JSON containment checks.
func (qb *Builder) OrWhereJSONContains(column string, value interface{}) types.QueryBuilder {
	clause := clauses.NewWhereJSONContainsClause(column, value)
	clause.SetBoolean(types.Or)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

// WhereJSONLength adds a WHERE clause for JSON array/object length comparisons.
func (qb *Builder) WhereJSONLength(column string, args ...interface{}) types.QueryBuilder {
	return qb.addJSONLengthWhere(column, types.And, args...)
}

// OrWhereJSONLength adds an OR WHERE clause for JSON array/object length comparisons.
func (qb *Builder) OrWhereJSONLength(column string, args ...interface{}) types.QueryBuilder {
	return qb.addJSONLengthWhere(column, types.Or, args...)
}

func (qb *Builder) addJSONLengthWhere(column string, boolean types.BooleanOperator, args ...interface{}) types.QueryBuilder {
	if len(args) == 0 {
		return qb
	}

	var operator types.Operator
	var value interface{}

	switch len(args) {
	case 1:
		operator = types.OpEqual
		value = args[0]
	case 2:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	default:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	}

	clause := clauses.NewWhereJSONLengthClause(column, operator, value)
	clause.SetBoolean(boolean)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

// WhereJSONPath adds a WHERE clause for JSON path-based comparisons.
func (qb *Builder) WhereJSONPath(column, path string, args ...interface{}) types.QueryBuilder {
	return qb.addJSONPathWhere(column, path, types.And, args...)
}

// OrWhereJSONPath adds an OR WHERE clause for JSON path-based comparisons.
func (qb *Builder) OrWhereJSONPath(column, path string, args ...interface{}) types.QueryBuilder {
	return qb.addJSONPathWhere(column, path, types.Or, args...)
}

func (qb *Builder) addJSONPathWhere(column, path string, boolean types.BooleanOperator, args ...interface{}) types.QueryBuilder {
	if len(args) == 0 {
		return qb
	}

	var operator types.Operator
	var value interface{}

	switch len(args) {
	case 1:
		operator = types.OpEqual
		value = args[0]
	case 2:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	default:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	}

	var raw string
	switch qb.driver {
	case types.MySQL:
		raw = fmt.Sprintf("JSON_EXTRACT(%s, '%s') %s ?", column, path, operator)
	case types.PostgreSQL:
		if strings.Contains(path, "[") {
			raw = fmt.Sprintf("%s #>> '%s' %s ?", column, path, operator)
		} else {
			raw = fmt.Sprintf("%s ->> '%s' %s ?", column, path, operator)
		}
	}

	clause := clauses.NewWhereRawClause(raw)
	clause.SetBoolean(boolean)
	qb.wheres = append(qb.wheres, clause)
	qb.bindings = append(qb.bindings, value)
	return qb
}

// WhereFullText adds a WHERE clause for full-text search across multiple columns.
func (qb *Builder) WhereFullText(columns []string, value string) types.QueryBuilder {
	clause := clauses.NewWhereFullTextClause(columns, value)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

// OrWhereFullText adds an OR WHERE clause for full-text search across multiple columns.
func (qb *Builder) OrWhereFullText(columns []string, value string) types.QueryBuilder {
	clause := clauses.NewWhereFullTextClause(columns, value)
	clause.SetBoolean(types.Or)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

// WhereAny adds a WHERE clause that matches if any of the specified columns meet the criteria.
func (qb *Builder) WhereAny(columns []string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereAny(columns, types.And, false, args...)
}

// OrWhereAny adds an OR WHERE clause that matches if any of the specified columns meet the criteria.
func (qb *Builder) OrWhereAny(columns []string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereAny(columns, types.Or, false, args...)
}

// WhereAll adds a WHERE clause that matches if all of the specified columns meet the criteria.
func (qb *Builder) WhereAll(columns []string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereAll(columns, types.And, false, args...)
}

// OrWhereAll adds an OR WHERE clause that matches if all of the specified columns meet the criteria.
func (qb *Builder) OrWhereAll(columns []string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereAll(columns, types.Or, false, args...)
}

// WhereNone adds a WHERE clause that matches if none of the specified columns meet the criteria.
func (qb *Builder) WhereNone(columns []string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereAny(columns, types.And, true, args...)
}

// OrWhereNone adds an OR WHERE clause that matches if none of the specified columns meet the criteria.
func (qb *Builder) OrWhereNone(columns []string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereAny(columns, types.Or, true, args...)
}

func (qb *Builder) addWhereAny(columns []string, boolean types.BooleanOperator, not bool, args ...interface{}) types.QueryBuilder {
	if len(columns) == 0 || len(args) == 0 {
		return qb
	}

	var operator types.Operator
	var value interface{}

	switch len(args) {
	case 1:
		operator = types.OpEqual
		value = args[0]
	case 2:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	default:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	}

	conditions := make([]string, len(columns))
	for i, column := range columns {
		conditions[i] = fmt.Sprintf("%s %s ?", column, operator)
	}

	var raw string
	if not {
		raw = fmt.Sprintf("NOT (%s)", strings.Join(conditions, " OR "))
	} else {
		raw = fmt.Sprintf("(%s)", strings.Join(conditions, " OR "))
	}

	clause := clauses.NewWhereRawClause(raw)
	clause.SetBoolean(boolean)
	qb.wheres = append(qb.wheres, clause)

	for i := 0; i < len(columns); i++ {
		qb.bindings = append(qb.bindings, value)
	}

	return qb
}

func (qb *Builder) addWhereAll(columns []string, boolean types.BooleanOperator, not bool, args ...interface{}) types.QueryBuilder {
	if len(columns) == 0 || len(args) == 0 {
		return qb
	}

	var operator types.Operator
	var value interface{}

	switch len(args) {
	case 1:
		operator = types.OpEqual
		value = args[0]
	case 2:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	default:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
		value = args[1]
	}

	conditions := make([]string, len(columns))
	for i, column := range columns {
		conditions[i] = fmt.Sprintf("%s %s ?", column, operator)
	}

	var raw string
	if not {
		raw = fmt.Sprintf("NOT (%s)", strings.Join(conditions, " AND "))
	} else {
		raw = fmt.Sprintf("(%s)", strings.Join(conditions, " AND "))
	}

	clause := clauses.NewWhereRawClause(raw)
	clause.SetBoolean(boolean)
	qb.wheres = append(qb.wheres, clause)

	for i := 0; i < len(columns); i++ {
		qb.bindings = append(qb.bindings, value)
	}

	return qb
}

// WhereColumn adds a WHERE clause that compares two columns.
func (qb *Builder) WhereColumn(first, second string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereColumn(first, second, types.And, args...)
}

// OrWhereColumn adds an OR WHERE clause that compares two columns.
func (qb *Builder) OrWhereColumn(first, second string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereColumn(first, second, types.Or, args...)
}

func (qb *Builder) addWhereColumn(first, second string, boolean types.BooleanOperator, args ...interface{}) types.QueryBuilder {
	var operator types.Operator

	switch len(args) {
	case 0:
		operator = types.OpEqual
	case 1:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
	default:
		operator = types.Operator(fmt.Sprintf("%v", args[0]))
	}

	raw := fmt.Sprintf("%s %s %s", first, operator, second)
	clause := clauses.NewWhereRawClause(raw)
	clause.SetBoolean(boolean)
	qb.wheres = append(qb.wheres, clause)
	return qb
}