package clauses

import (
	"github.com/go-query-builder/querybuilder/pkg/types"
)

type OrderClause struct {
	Column    string
	Direction types.OrderDirection
	Raw       string
}

func NewOrderClause(column string, direction types.OrderDirection) *OrderClause {
	return &OrderClause{
		Column:    column,
		Direction: direction,
	}
}

func NewOrderRawClause(raw string) *OrderClause {
	return &OrderClause{
		Raw: raw,
	}
}

func (o *OrderClause) IsRaw() bool {
	return o.Raw != ""
}

func (o *OrderClause) GetColumn() string {
	return o.Column
}

func (o *OrderClause) GetDirection() types.OrderDirection {
	return o.Direction
}

func (o *OrderClause) GetRaw() string {
	return o.Raw
}