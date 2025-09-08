// Package querybuilder provides a fluent SQL query builder for Go applications.
// It supports multiple database drivers and provides a clean, intuitive API for building complex queries.
package querybuilder

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/omarhamdy49/go-query-builder/pkg/config"
	"github.com/omarhamdy49/go-query-builder/pkg/database"
	"github.com/omarhamdy49/go-query-builder/pkg/query"
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// Global singleton instance.
var (
	builderInstance *Builder
	builderOnce     sync.Once
)

// Builder represents the singleton query builder instance.
type Builder struct {
	connections map[string]types.DB
	defaultConn string
	mu          sync.RWMutex
}

// Tabler interface for models that can provide table names.
type Tabler interface {
	TableName() string
}

// GetBuilder returns the singleton Builder instance with automatic environment loading.
func GetBuilder() *Builder {
	builderOnce.Do(func() {
		builderInstance = &Builder{
			connections: make(map[string]types.DB),
			defaultConn: "default",
			mu:          sync.RWMutex{},
		}

		// Automatically load environment configuration
		if err := builderInstance.loadEnvironmentConfig(); err != nil {
			log.Printf("Warning: Failed to load environment config: %v", err)
			log.Printf("You can manually add connections using AddConnection()")
		}
	})

	return builderInstance
}

// loadEnvironmentConfig loads configuration from environment variables.
func (b *Builder) loadEnvironmentConfig() error {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return fmt.Errorf("failed to load config from env: %w", err)
	}

	// Create default connection
	conn, err := database.NewConnection(cfg)
	if err != nil {
		return fmt.Errorf("failed to create database connection: %w", err)
	}

	b.mu.Lock()
	b.connections["default"] = conn
	b.defaultConn = "default"
	b.mu.Unlock()

	return nil
}

// AddConnection adds a named database connection.
func (b *Builder) AddConnection(name string, config *types.Config) error {
	conn, err := database.NewConnection(*config)
	if err != nil {
		return fmt.Errorf("failed to create connection %s: %w", name, err)
	}

	b.mu.Lock()
	defer b.mu.Unlock()

	b.connections[name] = conn
	if len(b.connections) == 1 {
		b.defaultConn = name
	}

	return nil
}

// Connection switches to a different database connection.
func (b *Builder) Connection(name string) *Builder {
	// Create a copy to avoid affecting the singleton state
	newBuilder := &Builder{
		connections: b.connections,
		defaultConn: name,
		mu:          sync.RWMutex{},
	}

	return newBuilder
}

// Table creates a query builder for a table using model, pointer, or string.
func (b *Builder) Table(table interface{}) types.QueryBuilder {
	b.mu.RLock()
	conn, exists := b.connections[b.defaultConn]
	b.mu.RUnlock()

	if !exists {
		log.Fatalf("Connection '%s' not found. Available connections: %v",
			b.defaultConn, b.getConnectionNames())
	}

	tableName := b.extractTableName(table)

	return query.Table(conn, conn.Driver(), tableName)
}

// extractTableName extracts table name from various input types.
func (b *Builder) extractTableName(table interface{}) string {
	switch tableType := table.(type) {
	case string:
		return tableType
	case Tabler:
		return tableType.TableName()
	default:
		// Try to extract table name from struct type
		return b.deriveTableNameFromType(tableType)
	}
}

// deriveTableNameFromType converts struct type to table name.
func (b *Builder) deriveTableNameFromType(model interface{}) string {
	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	// Convert struct name to snake_case and pluralize
	name := modelType.Name()
	snakeName := b.toSnakeCase(name)

	return b.pluralize(snakeName)
}

// toSnakeCase converts CamelCase to snake_case.
func (b *Builder) toSnakeCase(str string) string {
	var result strings.Builder

	for idx, char := range str {
		if idx > 0 && 'A' <= char && char <= 'Z' {
			result.WriteByte('_')
		}

		result.WriteRune(char)
	}

	return strings.ToLower(result.String())
}

// pluralize adds 's' to make table name plural (simple implementation).
func (b *Builder) pluralize(str string) string {
	if strings.HasSuffix(str, "y") {
		return strings.TrimSuffix(str, "y") + "ies"
	}

	if strings.HasSuffix(str, "s") || strings.HasSuffix(str, "x") || strings.HasSuffix(str, "z") {
		return str + "es"
	}

	return str + "s"
}

// getConnectionNames returns list of available connection names.
func (b *Builder) getConnectionNames() []string {
	b.mu.RLock()
	defer b.mu.RUnlock()

	names := make([]string, 0, len(b.connections))
	for name := range b.connections {
		names = append(names, name)
	}

	return names
}

// Close closes all database connections
func (b *Builder) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, dbConn := range b.connections {
		if err := dbConn.Close(); err != nil {
			return fmt.Errorf("failed to close connection: %w", err)
		}
	}

	return nil
}

// NewConnection creates a new database connection for backward compatibility.
func NewConnection(config *types.Config) (types.DB, error) {
	conn, err := database.NewConnection(*config)
	if err != nil {
		return nil, fmt.Errorf("failed to create new connection: %w", err)
	}

	return conn, nil
}

// NewQueryBuilder creates a new query builder instance.
func NewQueryBuilder(executor types.QueryExecutor, driver types.Driver) types.QueryBuilder {
	return query.NewBuilder(executor, driver)
}

// Table creates a query builder for a specific table.
func Table(executor types.QueryExecutor, driver types.Driver, table string) types.QueryBuilder {
	return query.Table(executor, driver, table)
}

// QB returns the singleton query builder instance.
func QB() *Builder {
	return GetBuilder()
}

// TableBuilder creates a query builder for the given table using the singleton instance.
func TableBuilder(table interface{}) types.QueryBuilder {
	return GetBuilder().Table(table)
}

// Connection switches to a named connection using the singleton instance.
func Connection(connectionName string) *Builder {
	return GetBuilder().Connection(connectionName)
}

// Config is an alias for types.Config.
type Config = types.Config

// Driver is an alias for types.Driver.
type Driver = types.Driver

// Collection is an alias for types.Collection.
type Collection = types.Collection

// QueryBuilder is an alias for types.QueryBuilder.
type QueryBuilder = types.QueryBuilder

// PaginationResult is an alias for types.PaginationResult.
type PaginationResult = types.PaginationResult

// AggregateResult is an alias for types.AggregateResult.
type AggregateResult = types.AggregateResult

// DebugInfo is an alias for types.DebugInfo.
type DebugInfo = types.DebugInfo

// OrderDirection is an alias for types.OrderDirection.
type OrderDirection = types.OrderDirection

// Database driver constants.
const (
	// MySQL database driver.
	MySQL = types.MySQL
	// PostgreSQL database driver.
	PostgreSQL = types.PostgreSQL
	// Asc represents ascending order.
	Asc = types.Asc
	// Desc represents descending order.
	Desc = types.Desc
)

// Collection factory functions.
var (
	// NewCollection creates a new collection instance.
	NewCollection = types.NewCollection
)
