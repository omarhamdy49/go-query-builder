// Package clauses provides SQL clause structures and builders for the go-query-builder.
// It includes implementations for GROUP BY, HAVING, JOIN, ORDER BY, SELECT, UNION, and WHERE clauses.
package clauses

// GroupClause represents a GROUP BY clause in an SQL query.
type GroupClause struct {
	Column string
	Raw    string
}

// NewGroupClause creates a new GROUP BY clause for the specified column.
func NewGroupClause(column string) *GroupClause {
	return &GroupClause{
		Column: column,
	}
}

// NewGroupRawClause creates a new GROUP BY clause using raw SQL.
func NewGroupRawClause(raw string) *GroupClause {
	return &GroupClause{
		Raw: raw,
	}
}

// IsRaw returns true if this is a raw GROUP BY clause.
func (g *GroupClause) IsRaw() bool {
	return g.Raw != ""
}

// GetColumn returns the column name for this GROUP BY clause.
func (g *GroupClause) GetColumn() string {
	return g.Column
}

// GetRaw returns the raw SQL for this GROUP BY clause.
func (g *GroupClause) GetRaw() string {
	return g.Raw
}