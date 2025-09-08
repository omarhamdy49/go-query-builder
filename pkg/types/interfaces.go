// Package types contains core types, interfaces, and constants used throughout the query builder.
package types

import (
	"context"
	"database/sql/driver"
)

// QueryExecutor defines the interface for executing database queries.
type QueryExecutor interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
	Begin() (Tx, error)
	BeginTx(ctx context.Context, opts *TxOptions) (Tx, error)
}

// Rows represents the result of a query.
type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Columns() ([]string, error)
	Err() error
}

// Row represents a single row result from a query.
type Row interface {
	Scan(dest ...interface{}) error
}

// Result represents the result of an exec query.
type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

// Tx represents a database transaction.
type Tx interface {
	QueryExecutor
	Commit() error
	Rollback() error
}

// TxOptions holds transaction configuration options.
type TxOptions struct {
	Isolation int
	ReadOnly  bool
}

// DB represents a database connection.
type DB interface {
	QueryExecutor
	Driver() Driver
	Close() error
	Ping() error
	Stats() DBStats
}

// DBStats holds database connection statistics.
type DBStats struct {
	OpenConnections int
	InUse           int
	Idle            int
}

// QueryBuilder defines the interface for building SQL queries fluently.
type QueryBuilder interface {
	From(table string) QueryBuilder
	Select(columns ...string) QueryBuilder
	SelectRaw(raw string, bindings ...interface{}) QueryBuilder
	SelectAs(column, alias string) QueryBuilder
	Distinct() QueryBuilder
	Where(column string, args ...interface{}) QueryBuilder
	OrWhere(column string, args ...interface{}) QueryBuilder
	WhereNot(column string, args ...interface{}) QueryBuilder
	OrWhereNot(column string, args ...interface{}) QueryBuilder
	WhereRaw(raw string, bindings ...interface{}) QueryBuilder
	OrWhereRaw(raw string, bindings ...interface{}) QueryBuilder
	WhereBetween(column string, values []interface{}) QueryBuilder
	WhereNotBetween(column string, values []interface{}) QueryBuilder
	WhereIn(column string, values []interface{}) QueryBuilder
	WhereNotIn(column string, values []interface{}) QueryBuilder
	WhereNull(column string) QueryBuilder
	WhereNotNull(column string) QueryBuilder
	WhereExists(query QueryBuilder) QueryBuilder
	WhereNotExists(query QueryBuilder) QueryBuilder
	WhereDate(column string, args ...interface{}) QueryBuilder
	OrWhereDate(column string, args ...interface{}) QueryBuilder
	WhereTime(column string, args ...interface{}) QueryBuilder
	OrWhereTime(column string, args ...interface{}) QueryBuilder
	WhereDay(column string, args ...interface{}) QueryBuilder
	OrWhereDay(column string, args ...interface{}) QueryBuilder
	WhereMonth(column string, args ...interface{}) QueryBuilder
	OrWhereMonth(column string, args ...interface{}) QueryBuilder
	WhereYear(column string, args ...interface{}) QueryBuilder
	OrWhereYear(column string, args ...interface{}) QueryBuilder
	WherePast(column string) QueryBuilder
	WhereFuture(column string) QueryBuilder
	WhereNowOrPast(column string) QueryBuilder
	WhereNowOrFuture(column string) QueryBuilder
	WhereToday(column string) QueryBuilder
	WhereBeforeToday(column string) QueryBuilder
	WhereAfterToday(column string) QueryBuilder
	WhereTodayOrBefore(column string) QueryBuilder
	WhereTodayOrAfter(column string) QueryBuilder
	WhereJsonContains(column string, value interface{}) QueryBuilder
	OrWhereJsonContains(column string, value interface{}) QueryBuilder
	WhereJsonLength(column string, args ...interface{}) QueryBuilder
	OrWhereJsonLength(column string, args ...interface{}) QueryBuilder
	WhereJsonPath(column, path string, args ...interface{}) QueryBuilder
	OrWhereJsonPath(column, path string, args ...interface{}) QueryBuilder
	WhereFullText(columns []string, value string) QueryBuilder
	OrWhereFullText(columns []string, value string) QueryBuilder
	WhereAny(columns []string, args ...interface{}) QueryBuilder
	OrWhereAny(columns []string, args ...interface{}) QueryBuilder
	WhereAll(columns []string, args ...interface{}) QueryBuilder
	OrWhereAll(columns []string, args ...interface{}) QueryBuilder
	WhereNone(columns []string, args ...interface{}) QueryBuilder
	OrWhereNone(columns []string, args ...interface{}) QueryBuilder
	WhereColumn(first, second string, args ...interface{}) QueryBuilder
	OrWhereColumn(first, second string, args ...interface{}) QueryBuilder
	Join(table, first string, args ...interface{}) QueryBuilder
	LeftJoin(table, first string, args ...interface{}) QueryBuilder
	RightJoin(table, first string, args ...interface{}) QueryBuilder
	CrossJoin(table string) QueryBuilder
	OrderBy(column string, direction ...OrderDirection) QueryBuilder
	OrderByDesc(column string) QueryBuilder
	OrderByRaw(raw string) QueryBuilder
	GroupBy(columns ...string) QueryBuilder
	GroupByRaw(raw string) QueryBuilder
	Having(column string, args ...interface{}) QueryBuilder
	OrHaving(column string, args ...interface{}) QueryBuilder
	HavingRaw(raw string) QueryBuilder
	OrHavingRaw(raw string) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	Take(limit int) QueryBuilder
	Skip(offset int) QueryBuilder
	Union(query QueryBuilder) QueryBuilder
	UnionAll(query QueryBuilder) QueryBuilder
	ForUpdate() QueryBuilder
	ForShare() QueryBuilder
	When(condition bool, callback ConditionalFunc) QueryBuilder
	Unless(condition bool, callback ConditionalFunc) QueryBuilder
	Tap(callback ConditionalFunc) QueryBuilder
	Scope(scopes ...ScopeFunc) QueryBuilder
	Debug() QueryBuilder
	ToSQL() (string, []interface{}, error)
	Get(ctx context.Context) (Collection, error)
	First(ctx context.Context) (map[string]interface{}, error)
	Find(ctx context.Context, id interface{}) (map[string]interface{}, error)
	Pluck(ctx context.Context, column string) ([]interface{}, error)
	Count(ctx context.Context) (int64, error)
	Sum(ctx context.Context, column string) (interface{}, error)
	Avg(ctx context.Context, column string) (interface{}, error)
	Min(ctx context.Context, column string) (interface{}, error)
	Max(ctx context.Context, column string) (interface{}, error)
	Insert(ctx context.Context, values map[string]interface{}) error
	InsertBatch(ctx context.Context, values []map[string]interface{}) error
	Update(ctx context.Context, values map[string]interface{}) (int64, error)
	Delete(ctx context.Context) (int64, error)
	Paginate(ctx context.Context, page int, perPage int) (PaginationResult, error)
	SimplePaginate(ctx context.Context, page int, perPage int) (PaginationResult, error)
	// Async methods
	GetAsync(ctx context.Context) <-chan AsyncResult
	CountAsync(ctx context.Context) <-chan AsyncCountResult
	PaginateAsync(ctx context.Context, page int, perPage int) <-chan AsyncPaginationResult
	Clone() QueryBuilder
}

// ConditionalFunc represents a function that can conditionally modify a query builder.
type ConditionalFunc func(QueryBuilder) QueryBuilder

// ScopeFunc represents a function that applies a scope to a query builder.
type ScopeFunc func(QueryBuilder) QueryBuilder

// ChunkFunc represents a function that processes data chunks.
type ChunkFunc func(Collection) error

// LazyFunc represents a function for lazy data processing.
type LazyFunc func(map[string]interface{}) error

// AsyncResult holds the result of an asynchronous operation.
type AsyncResult struct {
	Data  Collection
	Error error
}

// AsyncCountResult holds the result of an asynchronous count operation.
type AsyncCountResult struct {
	Count int64
	Error error
}

// AsyncPaginationResult holds the result of an asynchronous pagination operation.
type AsyncPaginationResult struct {
	Result PaginationResult
	Error  error
}

// JSONValue wraps a value for JSON handling.
type JSONValue struct {
	Val interface{}
}

// Value implements the driver.Valuer interface.
func (j JSONValue) Value() (driver.Value, error) {
	return j.Val, nil
}

// Collection defines the interface for working with collections of data.
type Collection interface {
	ToSlice() []map[string]interface{}
	Pluck(column string) []interface{}
	First() map[string]interface{}
	Count() int
	IsEmpty() bool
	Each(fn func(item map[string]interface{}) bool)
	Filter(predicate func(item map[string]interface{}) bool) Collection
	Map(mapper func(item map[string]interface{}) map[string]interface{}) Collection
}
