package clauses

import (
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// OrderClause represents an ORDER BY clause in an SQL query.
type OrderClause struct {
	Column    string
	Direction types.OrderDirection
	Raw       string
}

// NewOrderClause creates a new ORDER BY clause for the specified column and direction.
func NewOrderClause(column string, direction types.OrderDirection) *OrderClause {
	return &OrderClause{
		Column:    column,
		Direction: direction,
	}
}

// NewOrderRawClause creates a new ORDER BY clause using raw SQL.
func NewOrderRawClause(raw string) *OrderClause {
	return &OrderClause{
		Raw: raw,
	}
}

// IsRaw returns true if this is a raw ORDER BY clause.
func (o *OrderClause) IsRaw() bool {
	return o.Raw != ""
}

// GetColumn returns the column name for this ORDER BY clause.
func (o *OrderClause) GetColumn() string {
	return o.Column
}

// GetDirection returns the order direction for this ORDER BY clause.
func (o *OrderClause) GetDirection() types.OrderDirection {
	return o.Direction
}

// GetRaw returns the raw SQL for this ORDER BY clause.
func (o *OrderClause) GetRaw() string {
	return o.Raw
}