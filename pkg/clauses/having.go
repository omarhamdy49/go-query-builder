package clauses

import (
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// HavingClause represents a HAVING clause in an SQL query.
type HavingClause struct {
	Column   string
	Operator types.Operator
	Value    interface{}
	Boolean  types.BooleanOperator
	Raw      string
}

// NewHavingClause creates a new HAVING clause with the specified column, operator, and value.
func NewHavingClause(column string, operator types.Operator, value interface{}) *HavingClause {
	return &HavingClause{
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  types.And,
	}
}

// NewHavingRawClause creates a new HAVING clause using raw SQL.
func NewHavingRawClause(raw string) *HavingClause {
	return &HavingClause{
		Raw:     raw,
		Boolean: types.And,
	}
}

// SetBoolean sets the boolean operator for this HAVING clause.
func (h *HavingClause) SetBoolean(boolean types.BooleanOperator) *HavingClause {
	h.Boolean = boolean
	return h
}

// IsRaw returns true if this is a raw HAVING clause.
func (h *HavingClause) IsRaw() bool {
	return h.Raw != ""
}

// GetColumn returns the column name for this HAVING clause.
func (h *HavingClause) GetColumn() string {
	return h.Column
}

// GetOperator returns the operator for this HAVING clause.
func (h *HavingClause) GetOperator() types.Operator {
	return h.Operator
}

// GetValue returns the value for this HAVING clause.
func (h *HavingClause) GetValue() interface{} {
	return h.Value
}

// GetBoolean returns the boolean operator for this HAVING clause.
func (h *HavingClause) GetBoolean() types.BooleanOperator {
	return h.Boolean
}

// GetRaw returns the raw SQL for this HAVING clause.
func (h *HavingClause) GetRaw() string {
	return h.Raw
}