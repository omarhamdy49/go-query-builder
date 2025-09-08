package clauses

import (
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

type UnionClause struct {
	Query types.QueryBuilder
	All   bool
}

func NewUnionClause(query types.QueryBuilder) *UnionClause {
	return &UnionClause{
		Query: query,
		All:   false,
	}
}

func NewUnionAllClause(query types.QueryBuilder) *UnionClause {
	return &UnionClause{
		Query: query,
		All:   true,
	}
}

func (u *UnionClause) IsUnionAll() bool {
	return u.All
}

func (u *UnionClause) GetQuery() types.QueryBuilder {
	return u.Query
}