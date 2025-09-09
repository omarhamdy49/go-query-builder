package clauses

import (
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// WhereClause represents a WHERE clause in an SQL query.
type WhereClause struct {
	Type     string
	Column   string
	Operator types.Operator
	Value    interface{}
	Values   []interface{}
	Query    types.QueryBuilder
	Boolean  types.BooleanOperator
	Raw      string
}

// NewWhereClause creates a new WHERE clause with the specified column, operator, and value.
func NewWhereClause(column string, operator types.Operator, value interface{}) *WhereClause {
	return &WhereClause{
		Type:     "basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  types.And,
	}
}

// NewWhereRawClause creates a new WHERE clause using raw SQL.
func NewWhereRawClause(raw string) *WhereClause {
	return &WhereClause{
		Type:    "raw",
		Raw:     raw,
		Boolean: types.And,
	}
}

// NewWhereBetweenClause creates a new BETWEEN WHERE clause.
func NewWhereBetweenClause(column string, values []interface{}, not bool) *WhereClause {
	operator := types.OpBetween
	if not {
		operator = types.OpNotBetween
	}
	
	return &WhereClause{
		Type:     "between",
		Column:   column,
		Operator: operator,
		Values:   values,
		Boolean:  types.And,
	}
}

// NewWhereInClause creates a new IN or NOT IN WHERE clause.
func NewWhereInClause(column string, values []interface{}, not bool) *WhereClause {
	operator := types.OpIn
	if not {
		operator = types.OpNotIn
	}
	
	return &WhereClause{
		Type:     "in",
		Column:   column,
		Operator: operator,
		Values:   values,
		Boolean:  types.And,
	}
}

// NewWhereNullClause creates a new IS NULL or IS NOT NULL WHERE clause.
func NewWhereNullClause(column string, not bool) *WhereClause {
	operator := types.OpIsNull
	if not {
		operator = types.OpIsNotNull
	}
	
	return &WhereClause{
		Type:     "null",
		Column:   column,
		Operator: operator,
		Boolean:  types.And,
	}
}

// NewWhereExistsClause creates a new EXISTS or NOT EXISTS WHERE clause with a subquery.
func NewWhereExistsClause(query types.QueryBuilder, not bool) *WhereClause {
	operator := types.OpExists
	if not {
		operator = types.OpNotExists
	}
	
	return &WhereClause{
		Type:     "exists",
		Operator: operator,
		Query:    query,
		Boolean:  types.And,
	}
}

// NewWhereJSONContainsClause creates a new WHERE clause for JSON containment checks.
func NewWhereJSONContainsClause(column string, value interface{}) *WhereClause {
	return &WhereClause{
		Type:     "json",
		Column:   column,
		Operator: types.OpJSONContains,
		Value:    value,
		Boolean:  types.And,
	}
}

// NewWhereJSONLengthClause creates a new WHERE clause for JSON length comparisons.
func NewWhereJSONLengthClause(column string, operator types.Operator, value interface{}) *WhereClause {
	return &WhereClause{
		Type:     "json_length",
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  types.And,
	}
}

// NewWhereFullTextClause creates a new WHERE clause for full-text search across multiple columns.
func NewWhereFullTextClause(columns []string, value string) *WhereClause {
	return &WhereClause{
		Type:     "fulltext",
		Values:   make([]interface{}, len(columns)),
		Operator: types.OpFullText,
		Value:    value,
		Boolean:  types.And,
	}
}

// SetBoolean sets the boolean operator (AND/OR) for the WHERE clause and returns the clause.
func (w *WhereClause) SetBoolean(boolean types.BooleanOperator) *WhereClause {
	w.Boolean = boolean
	return w
}

// IsComplex returns true if the WHERE clause contains complex operations like EXISTS or subqueries.
func (w *WhereClause) IsComplex() bool {
	return w.Type == "exists" || w.Type == "subquery"
}

// HasSubQuery returns true if the WHERE clause contains a subquery.
func (w *WhereClause) HasSubQuery() bool {
	return w.Query != nil
}

// IsRaw returns true if the WHERE clause contains raw SQL.
func (w *WhereClause) IsRaw() bool {
	return w.Type == "raw"
}