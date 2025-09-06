package querybuilder

import (
	"github.com/go-query-builder/querybuilder/pkg/database"
	"github.com/go-query-builder/querybuilder/pkg/query"
	"github.com/go-query-builder/querybuilder/pkg/types"
)

func NewConnection(config types.Config) (types.DB, error) {
	return database.NewConnection(config)
}

func NewQueryBuilder(executor types.QueryExecutor, driver types.Driver) types.QueryBuilder {
	return query.NewBuilder(executor, driver)
}

func Table(executor types.QueryExecutor, driver types.Driver, table string) types.QueryBuilder {
	return query.Table(executor, driver, table)
}

type Config = types.Config
type Driver = types.Driver
type Collection = types.Collection
type QueryBuilder = types.QueryBuilder
type PaginationResult = types.PaginationResult
type AggregateResult = types.AggregateResult
type DebugInfo = types.DebugInfo

const (
	MySQL      = types.MySQL
	PostgreSQL = types.PostgreSQL
)

var (
	NewCollection = types.NewCollection
)