package clauses

import (
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// UnionClause represents a UNION clause in an SQL query.
type UnionClause struct {
	Query types.QueryBuilder
	All   bool
}

// NewUnionClause creates a new UNION clause for the specified query.
func NewUnionClause(query types.QueryBuilder) *UnionClause {
	return &UnionClause{
		Query: query,
		All:   false,
	}
}

// NewUnionAllClause creates a new UNION ALL clause for the specified query.
func NewUnionAllClause(query types.QueryBuilder) *UnionClause {
	return &UnionClause{
		Query: query,
		All:   true,
	}
}

// IsUnionAll returns true if this is a UNION ALL clause.
func (u *UnionClause) IsUnionAll() bool {
	return u.All
}

// GetQuery returns the query for this UNION clause.
func (u *UnionClause) GetQuery() types.QueryBuilder {
	return u.Query
}