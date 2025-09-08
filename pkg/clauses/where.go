package clauses

import (
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

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

func NewWhereClause(column string, operator types.Operator, value interface{}) *WhereClause {
	return &WhereClause{
		Type:     "basic",
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  types.And,
	}
}

func NewWhereRawClause(raw string) *WhereClause {
	return &WhereClause{
		Type:    "raw",
		Raw:     raw,
		Boolean: types.And,
	}
}

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

func NewWhereJsonContainsClause(column string, value interface{}) *WhereClause {
	return &WhereClause{
		Type:     "json",
		Column:   column,
		Operator: types.OpJsonContains,
		Value:    value,
		Boolean:  types.And,
	}
}

func NewWhereJsonLengthClause(column string, operator types.Operator, value interface{}) *WhereClause {
	return &WhereClause{
		Type:     "json_length",
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  types.And,
	}
}

func NewWhereFullTextClause(columns []string, value string) *WhereClause {
	return &WhereClause{
		Type:     "fulltext",
		Values:   make([]interface{}, len(columns)),
		Operator: types.OpFullText,
		Value:    value,
		Boolean:  types.And,
	}
}

func (w *WhereClause) SetBoolean(boolean types.BooleanOperator) *WhereClause {
	w.Boolean = boolean
	return w
}

func (w *WhereClause) IsComplex() bool {
	return w.Type == "exists" || w.Type == "subquery"
}

func (w *WhereClause) HasSubQuery() bool {
	return w.Query != nil
}

func (w *WhereClause) IsRaw() bool {
	return w.Type == "raw"
}