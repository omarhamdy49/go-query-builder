package types

type Driver string

const (
	MySQL      Driver = "mysql"
	PostgreSQL Driver = "postgres"
)

type Operator string

const (
	OpEqual              Operator = "="
	OpNotEqual           Operator = "!="
	OpGreaterThan        Operator = ">"
	OpGreaterThanOrEqual Operator = ">="
	OpLessThan           Operator = "<"
	OpLessThanOrEqual    Operator = "<="
	OpLike               Operator = "LIKE"
	OpNotLike            Operator = "NOT LIKE"
	OpILike              Operator = "ILIKE"
	OpNotILike           Operator = "NOT ILIKE"
	OpIn                 Operator = "IN"
	OpNotIn              Operator = "NOT IN"
	OpBetween            Operator = "BETWEEN"
	OpNotBetween         Operator = "NOT BETWEEN"
	OpIsNull             Operator = "IS NULL"
	OpIsNotNull          Operator = "IS NOT NULL"
	OpExists             Operator = "EXISTS"
	OpNotExists          Operator = "NOT EXISTS"
	OpJsonContains       Operator = "JSON_CONTAINS"
	OpJsonExtract        Operator = "JSON_EXTRACT"
	OpFullText           Operator = "MATCH"
)

type JoinType string

const (
	InnerJoin JoinType = "INNER JOIN"
	LeftJoin  JoinType = "LEFT JOIN"
	RightJoin JoinType = "RIGHT JOIN"
	CrossJoin JoinType = "CROSS JOIN"
	FullJoin  JoinType = "FULL JOIN"
)

type OrderDirection string

const (
	Asc  OrderDirection = "ASC"
	Desc OrderDirection = "DESC"
)

type LockType string

const (
	ForUpdate   LockType = "FOR UPDATE"
	ForShare    LockType = "FOR SHARE"
	ForUpdateNW LockType = "FOR UPDATE NOWAIT"
	ForShareNW  LockType = "FOR SHARE NOWAIT"
	ForUpdateSL LockType = "FOR UPDATE SKIP LOCKED"
	ForShareSL  LockType = "FOR SHARE SKIP LOCKED"
)

type BooleanOperator string

const (
	And BooleanOperator = "AND"
	Or  BooleanOperator = "OR"
)

type AggregateFunction string

const (
	Count AggregateFunction = "COUNT"
	Sum   AggregateFunction = "SUM"
	Avg   AggregateFunction = "AVG"
	Min   AggregateFunction = "MIN"
	Max   AggregateFunction = "MAX"
)

type ConflictAction string

const (
	DoNothing ConflictAction = "DO NOTHING"
	DoUpdate  ConflictAction = "DO UPDATE"
)