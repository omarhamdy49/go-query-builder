// Package types contains core types, interfaces, and constants used throughout the query builder.
package types

// Driver represents different database drivers supported by the query builder.
type Driver string

// Database drivers.
const (
	// MySQL driver.
	MySQL Driver = "mysql"
	// PostgreSQL driver.
	PostgreSQL Driver = "postgres"
)

// Operator represents SQL comparison and logical operators.
type Operator string

// SQL operators.
const (
	// OpEqual represents the = operator.
	OpEqual Operator = "="
	// OpNotEqual represents the != operator.
	OpNotEqual Operator = "!="
	// OpGreaterThan represents the > operator.
	OpGreaterThan Operator = ">"
	// OpGreaterThanOrEqual represents the >= operator.
	OpGreaterThanOrEqual Operator = ">="
	// OpLessThan represents the < operator.
	OpLessThan Operator = "<"
	// OpLessThanOrEqual represents the <= operator.
	OpLessThanOrEqual Operator = "<="
	// OpLike represents the LIKE operator.
	OpLike Operator = "LIKE"
	// OpNotLike represents the NOT LIKE operator.
	OpNotLike Operator = "NOT LIKE"
	// OpILike represents the ILIKE operator.
	OpILike Operator = "ILIKE"
	// OpNotILike represents the NOT ILIKE operator.
	OpNotILike Operator = "NOT ILIKE"
	// OpIn represents the IN operator.
	OpIn Operator = "IN"
	// OpNotIn represents the NOT IN operator.
	OpNotIn Operator = "NOT IN"
	// OpBetween represents the BETWEEN operator.
	OpBetween Operator = "BETWEEN"
	// OpNotBetween represents the NOT BETWEEN operator.
	OpNotBetween Operator = "NOT BETWEEN"
	// OpIsNull represents the IS NULL operator.
	OpIsNull Operator = "IS NULL"
	// OpIsNotNull represents the IS NOT NULL operator.
	OpIsNotNull Operator = "IS NOT NULL"
	// OpExists represents the EXISTS operator.
	OpExists Operator = "EXISTS"
	// OpNotExists represents the NOT EXISTS operator.
	OpNotExists Operator = "NOT EXISTS"
	// OpJSONContains represents the JSON_CONTAINS operator.
	OpJSONContains Operator = "JSON_CONTAINS"
	// OpJSONExtract represents the JSON_EXTRACT operator.
	OpJSONExtract Operator = "JSON_EXTRACT"
	// OpFullText represents the MATCH operator.
	OpFullText Operator = "MATCH"
)

// JoinType represents different types of SQL joins.
type JoinType string

// SQL join types.
const (
	// InnerJoin represents INNER JOIN.
	InnerJoin JoinType = "INNER JOIN"
	// LeftJoin represents LEFT JOIN.
	LeftJoin JoinType = "LEFT JOIN"
	// RightJoin represents RIGHT JOIN.
	RightJoin JoinType = "RIGHT JOIN"
	// CrossJoin represents CROSS JOIN.
	CrossJoin JoinType = "CROSS JOIN"
	// FullJoin represents FULL JOIN.
	FullJoin JoinType = "FULL JOIN"
)

// OrderDirection represents SQL order directions.
type OrderDirection string

// Order directions.
const (
	// Asc represents ascending order.
	Asc OrderDirection = "ASC"
	// Desc represents descending order.
	Desc OrderDirection = "DESC"
)

// LockType represents different types of row locking in SQL.
type LockType string

// Row locking types.
const (
	// ForUpdate represents FOR UPDATE lock.
	ForUpdate LockType = "FOR UPDATE"
	// ForShare represents FOR SHARE lock.
	ForShare LockType = "FOR SHARE"
	// ForUpdateNW represents FOR UPDATE NOWAIT lock.
	ForUpdateNW LockType = "FOR UPDATE NOWAIT"
	// ForShareNW represents FOR SHARE NOWAIT lock.
	ForShareNW LockType = "FOR SHARE NOWAIT"
	// ForUpdateSL represents FOR UPDATE SKIP LOCKED lock.
	ForUpdateSL LockType = "FOR UPDATE SKIP LOCKED"
	// ForShareSL represents FOR SHARE SKIP LOCKED lock.
	ForShareSL LockType = "FOR SHARE SKIP LOCKED"
)

// BooleanOperator represents logical operators for combining conditions.
type BooleanOperator string

// Boolean operators.
const (
	// And represents the AND operator.
	And BooleanOperator = "AND"
	// Or represents the OR operator.
	Or BooleanOperator = "OR"
)

// AggregateFunction represents SQL aggregate functions.
type AggregateFunction string

// Aggregate functions.
const (
	// Count represents the COUNT function.
	Count AggregateFunction = "COUNT"
	// Sum represents the SUM function.
	Sum AggregateFunction = "SUM"
	// Avg represents the AVG function.
	Avg AggregateFunction = "AVG"
	// Min represents the MIN function.
	Min AggregateFunction = "MIN"
	// Max represents the MAX function.
	Max AggregateFunction = "MAX"
)

// ConflictAction represents actions to take when handling conflicts in upsert operations.
type ConflictAction string

// Conflict resolution actions.
const (
	// DoNothing represents DO NOTHING action.
	DoNothing ConflictAction = "DO NOTHING"
	// DoUpdate represents DO UPDATE action.
	DoUpdate ConflictAction = "DO UPDATE"
)
