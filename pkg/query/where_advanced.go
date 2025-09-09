package query

import (
	"fmt"
	"strings"
	"time"

	"github.com/omarhamdy49/go-query-builder/pkg/clauses"
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

func (qb *Builder) WhereDate(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "date", types.And, args...)
}

func (qb *Builder) OrWhereDate(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "date", types.Or, args...)
}

func (qb *Builder) WhereTime(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "time", types.And, args...)
}

func (qb *Builder) OrWhereTime(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "time", types.Or, args...)
}

func (qb *Builder) WhereDay(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "day", types.And, args...)
}

func (qb *Builder) OrWhereDay(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "day", types.Or, args...)
}

func (qb *Builder) WhereMonth(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "month", types.And, args...)
}

func (qb *Builder) OrWhereMonth(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "month", types.Or, args...)
}

func (qb *Builder) WhereYear(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "year", types.And, args...)
}

func (qb *Builder) OrWhereYear(column string, args ...interface{}) types.QueryBuilder {
	return qb.addDateWhere(column, "year", types.Or, args...)
}

func (qb *Builder) WherePast(column string) types.QueryBuilder {
	return qb.Where(column, "<", time.Now())
}

func (qb *Builder) WhereFuture(column string) types.QueryBuilder {
	return qb.Where(column, ">", time.Now())
}

func (qb *Builder) WhereNowOrPast(column string) types.QueryBuilder {
	return qb.Where(column, "<=", time.Now())
}

func (qb *Builder) WhereNowOrFuture(column string) types.QueryBuilder {
	return qb.Where(column, ">=", time.Now())
}

func (qb *Builder) WhereToday(column string) types.QueryBuilder {
	return qb.WhereDate(column, time.Now().Format("2006-01-02"))
}

func (qb *Builder) WhereBeforeToday(column string) types.QueryBuilder {
	return qb.WhereDate(column, "<", time.Now().Format("2006-01-02"))
}

func (qb *Builder) WhereAfterToday(column string) types.QueryBuilder {
	return qb.WhereDate(column, ">", time.Now().Format("2006-01-02"))
}

func (qb *Builder) WhereTodayOrBefore(column string) types.QueryBuilder {
	return qb.WhereDate(column, "<=", time.Now().Format("2006-01-02"))
}

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

func (qb *Builder) WhereJSONContains(column string, value interface{}) types.QueryBuilder {
	clause := clauses.NewWhereJSONContainsClause(column, value)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) OrWhereJSONContains(column string, value interface{}) types.QueryBuilder {
	clause := clauses.NewWhereJSONContainsClause(column, value)
	clause.SetBoolean(types.Or)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) WhereJSONLength(column string, args ...interface{}) types.QueryBuilder {
	return qb.addJSONLengthWhere(column, types.And, args...)
}

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

func (qb *Builder) WhereJSONPath(column, path string, args ...interface{}) types.QueryBuilder {
	return qb.addJSONPathWhere(column, path, types.And, args...)
}

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

func (qb *Builder) WhereFullText(columns []string, value string) types.QueryBuilder {
	clause := clauses.NewWhereFullTextClause(columns, value)
	clause.SetBoolean(types.And)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) OrWhereFullText(columns []string, value string) types.QueryBuilder {
	clause := clauses.NewWhereFullTextClause(columns, value)
	clause.SetBoolean(types.Or)
	qb.wheres = append(qb.wheres, clause)
	return qb
}

func (qb *Builder) WhereAny(columns []string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereAny(columns, types.And, false, args...)
}

func (qb *Builder) OrWhereAny(columns []string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereAny(columns, types.Or, false, args...)
}

func (qb *Builder) WhereAll(columns []string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereAll(columns, types.And, false, args...)
}

func (qb *Builder) OrWhereAll(columns []string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereAll(columns, types.Or, false, args...)
}

func (qb *Builder) WhereNone(columns []string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereAny(columns, types.And, true, args...)
}

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

func (qb *Builder) WhereColumn(first, second string, args ...interface{}) types.QueryBuilder {
	return qb.addWhereColumn(first, second, types.And, args...)
}

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