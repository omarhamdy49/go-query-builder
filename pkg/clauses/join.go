package clauses

import (
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// JoinClause represents a JOIN clause in an SQL query.
type JoinClause struct {
	Type      types.JoinType
	Table     string
	First     string
	Operator  types.Operator
	Second    string
	Clauses   []*WhereClause
}

// NewJoinClause creates a new JOIN clause with the specified parameters.
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

// NewCrossJoinClause creates a new CROSS JOIN clause for the specified table.
func NewCrossJoinClause(table string) *JoinClause {
	return &JoinClause{
		Type:    types.CrossJoin,
		Table:   table,
		Clauses: make([]*WhereClause, 0),
	}
}

// AddWhereClause adds a WHERE clause to this JOIN.
func (j *JoinClause) AddWhereClause(clause *WhereClause) *JoinClause {
	j.Clauses = append(j.Clauses, clause)
	return j
}

// Where adds a WHERE condition to this JOIN clause.
func (j *JoinClause) Where(column string, operator types.Operator, value interface{}) *JoinClause {
	clause := NewWhereClause(column, operator, value)
	return j.AddWhereClause(clause)
}

// OrWhere adds an OR WHERE condition to this JOIN clause.
func (j *JoinClause) OrWhere(column string, operator types.Operator, value interface{}) *JoinClause {
	clause := NewWhereClause(column, operator, value)
	clause.SetBoolean(types.Or)
	return j.AddWhereClause(clause)
}

// IsCrossJoin returns true if this is a CROSS JOIN.
func (j *JoinClause) IsCrossJoin() bool {
	return j.Type == types.CrossJoin
}

// HasAdditionalClauses returns true if this JOIN has additional WHERE clauses.
func (j *JoinClause) HasAdditionalClauses() bool {
	return len(j.Clauses) > 0
}

// GetTable returns the table name for this JOIN.
func (j *JoinClause) GetTable() string {
	return j.Table
}

// GetType returns the JOIN type.
func (j *JoinClause) GetType() types.JoinType {
	return j.Type
}