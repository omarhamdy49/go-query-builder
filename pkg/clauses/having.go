package clauses

import (
	"github.com/go-query-builder/querybuilder/pkg/types"
)

type HavingClause struct {
	Column   string
	Operator types.Operator
	Value    interface{}
	Boolean  types.BooleanOperator
	Raw      string
}

func NewHavingClause(column string, operator types.Operator, value interface{}) *HavingClause {
	return &HavingClause{
		Column:   column,
		Operator: operator,
		Value:    value,
		Boolean:  types.And,
	}
}

func NewHavingRawClause(raw string) *HavingClause {
	return &HavingClause{
		Raw:     raw,
		Boolean: types.And,
	}
}

func (h *HavingClause) SetBoolean(boolean types.BooleanOperator) *HavingClause {
	h.Boolean = boolean
	return h
}

func (h *HavingClause) IsRaw() bool {
	return h.Raw != ""
}

func (h *HavingClause) GetColumn() string {
	return h.Column
}

func (h *HavingClause) GetOperator() types.Operator {
	return h.Operator
}

func (h *HavingClause) GetValue() interface{} {
	return h.Value
}

func (h *HavingClause) GetBoolean() types.BooleanOperator {
	return h.Boolean
}

func (h *HavingClause) GetRaw() string {
	return h.Raw
}