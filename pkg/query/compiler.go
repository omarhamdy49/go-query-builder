package query

import (
	"fmt"
	"strings"
	"time"

	"github.com/omarhamdy49/go-query-builder/pkg/clauses"
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// SQLCompiler compiles query builder objects into SQL statements for different database drivers.
type SQLCompiler struct {
	driver    types.Driver
	debug     bool
	debugInfo *types.DebugInfo
}

// NewSQLCompiler creates a new SQL compiler for the specified database driver.
func NewSQLCompiler(driver types.Driver) *SQLCompiler {
	return &SQLCompiler{
		driver: driver,
		debug:  false,
	}
}

// Debug enables debug mode to collect compilation information.
func (c *SQLCompiler) Debug() *SQLCompiler {
	c.debug = true
	return c
}

// GetDebugInfo returns debug information from the last compilation.
func (c *SQLCompiler) GetDebugInfo() *types.DebugInfo {
	return c.debugInfo
}

// CompileSelect compiles a query builder into a SELECT SQL statement.
func (c *SQLCompiler) CompileSelect(qb *Builder) (string, []interface{}, error) {
	start := time.Now()
	var parts []string
	var bindings []interface{}

	selects := c.compileSelects(qb.GetSelects(), qb.IsDistinct())
	parts = append(parts, "SELECT "+selects)

	if table := qb.GetTable(); table != "" {
		parts = append(parts, "FROM "+table)
	}

	if joins := qb.GetJoins(); len(joins) > 0 {
		joinSQL, joinBindings := c.compileJoins(joins)
		parts = append(parts, joinSQL)
		bindings = append(bindings, joinBindings...)
	}

	if wheres := qb.GetWheres(); len(wheres) > 0 {
		whereSQL, whereBindings := c.compileWheres(wheres)
		parts = append(parts, "WHERE "+whereSQL)
		bindings = append(bindings, whereBindings...)
	}

	if groups := qb.GetGroups(); len(groups) > 0 {
		parts = append(parts, "GROUP BY "+c.compileGroups(groups))
	}

	if havings := qb.GetHavings(); len(havings) > 0 {
		havingSQL, havingBindings := c.compileHavings(havings)
		parts = append(parts, "HAVING "+havingSQL)
		bindings = append(bindings, havingBindings...)
	}

	if unions := qb.GetUnions(); len(unions) > 0 {
		unionSQL, unionBindings := c.compileUnions(unions)
		parts = append(parts, unionSQL)
		bindings = append(bindings, unionBindings...)
	}

	if orders := qb.GetOrders(); len(orders) > 0 {
		parts = append(parts, "ORDER BY "+c.compileOrders(orders))
	}

	if limit := qb.GetLimit(); limit != nil {
		parts = append(parts, c.compileLimit(*limit))
	}

	if offset := qb.GetOffset(); offset != nil {
		parts = append(parts, c.compileOffset(*offset))
	}

	if lock := qb.GetLock(); lock != nil {
		parts = append(parts, string(*lock))
	}

	sql := strings.Join(parts, " ")
	bindings = append(bindings, qb.GetBindings()...)

	if c.debug {
		c.debugInfo = &types.DebugInfo{
			SQL:      sql,
			Bindings: bindings,
			Driver:   c.driver,
			Duration: time.Since(start),
		}
	}

	return sql, bindings, nil
}

func (c *SQLCompiler) compileSelects(selects []*clauses.SelectClause, distinct bool) string {
	if len(selects) == 0 {
		if distinct {
			return "DISTINCT *"
		}
		return "*"
	}

	var selectParts []string
	for _, sel := range selects {
		if sel.IsRaw() {
			selectParts = append(selectParts, sel.GetRaw())
		} else if sel.HasAlias() {
			selectParts = append(selectParts, fmt.Sprintf("%s AS %s", sel.GetColumn(), sel.GetAlias()))
		} else {
			selectParts = append(selectParts, sel.GetColumn())
		}
	}

	result := strings.Join(selectParts, ", ")
	if distinct {
		return "DISTINCT " + result
	}
	return result
}

func (c *SQLCompiler) compileJoins(joins []*clauses.JoinClause) (string, []interface{}) {
	var parts []string
	var bindings []interface{}

	for _, join := range joins {
		if join.IsCrossJoin() {
			parts = append(parts, fmt.Sprintf("%s %s", join.GetType(), join.GetTable()))
		} else {
			parts = append(parts, fmt.Sprintf("%s %s ON %s %s %s",
				join.GetType(), join.GetTable(), join.First, join.Operator, join.Second))
		}

		for _, clause := range join.Clauses {
			clauseSQL, clauseBindings := c.compileWhereClause(clause)
			parts = append(parts, "AND "+clauseSQL)
			bindings = append(bindings, clauseBindings...)
		}
	}

	return strings.Join(parts, " "), bindings
}

func (c *SQLCompiler) compileWheres(wheres []*clauses.WhereClause) (string, []interface{}) {
	var parts []string
	var bindings []interface{}

	for i, where := range wheres {
		clauseSQL, clauseBindings := c.compileWhereClause(where)
		
		if i == 0 {
			parts = append(parts, clauseSQL)
		} else {
			parts = append(parts, strings.ToUpper(string(where.Boolean))+" "+clauseSQL)
		}
		
		bindings = append(bindings, clauseBindings...)
	}

	return strings.Join(parts, " "), bindings
}

func (c *SQLCompiler) compileWhereClause(where *clauses.WhereClause) (string, []interface{}) {
	var bindings []interface{}

	switch where.Type {
	case "basic":
		bindings = append(bindings, where.Value)
		return fmt.Sprintf("%s %s %s", where.Column, where.Operator, c.getParameterPlaceholder()), bindings
	case "raw":
		return where.Raw, bindings
	case "between":
		bindings = append(bindings, where.Values...)
		return fmt.Sprintf("%s %s %s AND %s", where.Column, where.Operator, 
			c.getParameterPlaceholder(), c.getParameterPlaceholder()), bindings
	case "in":
		placeholders := c.getInPlaceholders(len(where.Values))
		bindings = append(bindings, where.Values...)
		return fmt.Sprintf("%s %s (%s)", where.Column, where.Operator, placeholders), bindings
	case "null":
		return fmt.Sprintf("%s %s", where.Column, where.Operator), bindings
	case "exists":
		subSQL, subBindings, _ := where.Query.ToSQL()
		bindings = append(bindings, subBindings...)
		return fmt.Sprintf("%s (%s)", where.Operator, subSQL), bindings
	case "json":
		return c.compileJSONWhereClause(where, &bindings)
	case "json_length":
		return c.compileJSONLengthWhereClause(where, &bindings)
	case "fulltext":
		return c.compileFullTextWhereClause(where, &bindings)
	default:
		return "", bindings
	}
}

func (c *SQLCompiler) compileJSONWhereClause(where *clauses.WhereClause, bindings *[]interface{}) (string, []interface{}) {
	switch c.driver {
	case types.MySQL:
		*bindings = append(*bindings, where.Column, where.Value)
		return fmt.Sprintf("JSON_CONTAINS(%s, %s)", c.getParameterPlaceholder(), c.getParameterPlaceholder()), *bindings
	case types.PostgreSQL:
		*bindings = append(*bindings, where.Value)
		return fmt.Sprintf("%s @> %s", where.Column, c.getParameterPlaceholder()), *bindings
	default:
		return "", *bindings
	}
}

func (c *SQLCompiler) compileJSONLengthWhereClause(where *clauses.WhereClause, bindings *[]interface{}) (string, []interface{}) {
	switch c.driver {
	case types.MySQL:
		*bindings = append(*bindings, where.Value)
		return fmt.Sprintf("JSON_LENGTH(%s) %s %s", where.Column, where.Operator, c.getParameterPlaceholder()), *bindings
	case types.PostgreSQL:
		*bindings = append(*bindings, where.Value)
		return fmt.Sprintf("jsonb_array_length(%s) %s %s", where.Column, where.Operator, c.getParameterPlaceholder()), *bindings
	default:
		return "", *bindings
	}
}

func (c *SQLCompiler) compileFullTextWhereClause(where *clauses.WhereClause, bindings *[]interface{}) (string, []interface{}) {
	columns := make([]string, len(where.Values))
	for i, col := range where.Values {
		columns[i] = col.(string)
	}
	
	switch c.driver {
	case types.MySQL:
		*bindings = append(*bindings, where.Value)
		return fmt.Sprintf("MATCH(%s) AGAINST(%s)", strings.Join(columns, ", "), c.getParameterPlaceholder()), *bindings
	case types.PostgreSQL:
		*bindings = append(*bindings, where.Value)
		return fmt.Sprintf("to_tsvector(%s) @@ plainto_tsquery(%s)", strings.Join(columns, " || ' ' || "), c.getParameterPlaceholder()), *bindings
	default:
		return "", *bindings
	}
}

func (c *SQLCompiler) compileGroups(groups []*clauses.GroupClause) string {
	var parts []string
	for _, group := range groups {
		if group.IsRaw() {
			parts = append(parts, group.GetRaw())
		} else {
			parts = append(parts, group.GetColumn())
		}
	}
	return strings.Join(parts, ", ")
}

func (c *SQLCompiler) compileHavings(havings []*clauses.HavingClause) (string, []interface{}) {
	var parts []string
	var bindings []interface{}

	for i, having := range havings {
		var havingSQL string
		if having.IsRaw() {
			havingSQL = having.GetRaw()
		} else {
			havingSQL = fmt.Sprintf("%s %s %s", having.GetColumn(), having.GetOperator(), c.getParameterPlaceholder())
			bindings = append(bindings, having.GetValue())
		}

		if i == 0 {
			parts = append(parts, havingSQL)
		} else {
			parts = append(parts, strings.ToUpper(string(having.GetBoolean()))+" "+havingSQL)
		}
	}

	return strings.Join(parts, " "), bindings
}

func (c *SQLCompiler) compileOrders(orders []*clauses.OrderClause) string {
	var parts []string
	for _, order := range orders {
		if order.IsRaw() {
			parts = append(parts, order.GetRaw())
		} else {
			parts = append(parts, fmt.Sprintf("%s %s", order.GetColumn(), order.GetDirection()))
		}
	}
	return strings.Join(parts, ", ")
}

func (c *SQLCompiler) compileUnions(unions []*clauses.UnionClause) (string, []interface{}) {
	var parts []string
	var bindings []interface{}

	for _, union := range unions {
		unionSQL, unionBindings, _ := union.GetQuery().ToSQL()
		
		if union.IsUnionAll() {
			parts = append(parts, "UNION ALL ("+unionSQL+")")
		} else {
			parts = append(parts, "UNION ("+unionSQL+")")
		}
		
		bindings = append(bindings, unionBindings...)
	}

	return strings.Join(parts, " "), bindings
}

func (c *SQLCompiler) compileLimit(limit int) string {
	return fmt.Sprintf("LIMIT %d", limit)
}

func (c *SQLCompiler) compileOffset(offset int) string {
	return fmt.Sprintf("OFFSET %d", offset)
}

func (c *SQLCompiler) getParameterPlaceholder() string {
	switch c.driver {
	case types.PostgreSQL:
		return "$%d"
	default:
		return "?"
	}
}

func (c *SQLCompiler) getInPlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = c.getParameterPlaceholder()
	}
	return strings.Join(placeholders, ", ")
}