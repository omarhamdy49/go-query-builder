package clauses

import (
	"github.com/go-query-builder/querybuilder/pkg/types"
)

type JoinClause struct {
	Type      types.JoinType
	Table     string
	First     string
	Operator  types.Operator
	Second    string
	Clauses   []*WhereClause
}

func NewJoinClause(joinType types.JoinType, table, first string, operator types.Operator, second string) *JoinClause {
	return &JoinClause{
		Type:     joinType,
		Table:    table,
		First:    first,
		Operator: operator,
		Second:   second,
		Clauses:  make([]*WhereClause, 0),
	}
}

func NewCrossJoinClause(table string) *JoinClause {
	return &JoinClause{
		Type:    types.CrossJoin,
		Table:   table,
		Clauses: make([]*WhereClause, 0),
	}
}

func (j *JoinClause) AddWhereClause(clause *WhereClause) *JoinClause {
	j.Clauses = append(j.Clauses, clause)
	return j
}

func (j *JoinClause) Where(column string, operator types.Operator, value interface{}) *JoinClause {
	clause := NewWhereClause(column, operator, value)
	return j.AddWhereClause(clause)
}

func (j *JoinClause) OrWhere(column string, operator types.Operator, value interface{}) *JoinClause {
	clause := NewWhereClause(column, operator, value)
	clause.SetBoolean(types.Or)
	return j.AddWhereClause(clause)
}

func (j *JoinClause) IsCrossJoin() bool {
	return j.Type == types.CrossJoin
}

func (j *JoinClause) HasAdditionalClauses() bool {
	return len(j.Clauses) > 0
}

func (j *JoinClause) GetTable() string {
	return j.Table
}

func (j *JoinClause) GetType() types.JoinType {
	return j.Type
}