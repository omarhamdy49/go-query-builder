package types

import (
	"context"
	"database/sql/driver"
	"time"
)

type QueryExecutor interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (Result, error)
	Begin() (Tx, error)
	BeginTx(ctx context.Context, opts *TxOptions) (Tx, error)
}

type Rows interface {
	Next() bool
	Scan(dest ...interface{}) error
	Close() error
	Columns() ([]string, error)
	Err() error
}

type Row interface {
	Scan(dest ...interface{}) error
}

type Result interface {
	LastInsertId() (int64, error)
	RowsAffected() (int64, error)
}

type Tx interface {
	QueryExecutor
	Commit() error
	Rollback() error
}

type TxOptions struct {
	Isolation int
	ReadOnly  bool
}

type DB interface {
	QueryExecutor
	Driver() Driver
	Close() error
	Ping() error
	Stats() DBStats
}

type DBStats struct {
	OpenConnections int
	InUse          int
	Idle           int
}

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
	Clone() QueryBuilder
}

type ConditionalFunc func(QueryBuilder) QueryBuilder
type ScopeFunc func(QueryBuilder) QueryBuilder
type ChunkFunc func(Collection) error
type LazyFunc func(map[string]interface{}) error

type JSONValue struct {
	Value interface{}
}

func (j JSONValue) Value() (driver.Value, error) {
	return j.Value, nil
}

type Collection interface {
	ToSlice() []map[string]interface{}
	Pluck(column string) []interface{}
	First() map[string]interface{}
	Count() int
	IsEmpty() bool
	Each(func(map[string]interface{}) bool)
	Filter(func(map[string]interface{}) bool) Collection
	Map(func(map[string]interface{}) map[string]interface{}) Collection
}