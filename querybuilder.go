package querybuilder

import (
	"log"
	"reflect"
	"strings"
	"sync"

	"github.com/omarhamdy49/go-query-builder/pkg/config"
	"github.com/omarhamdy49/go-query-builder/pkg/database"
	"github.com/omarhamdy49/go-query-builder/pkg/query"
	"github.com/omarhamdy49/go-query-builder/pkg/types"
)

// Global singleton instance
var (
	builderInstance *Builder
	builderOnce     sync.Once
)

// Builder represents the singleton query builder instance
type Builder struct {
	connections map[string]types.DB
	defaultConn string
	mu          sync.RWMutex
}

// Tabler interface for models that can provide table names
type Tabler interface {
	TableName() string
}

// init automatically loads environment configuration when package is imported
func init() {
	// Initialize the singleton builder with environment config
	GetBuilder()
}

// GetBuilder returns the singleton Builder instance with automatic environment loading
func GetBuilder() *Builder {
	builderOnce.Do(func() {
		builderInstance = &Builder{
			connections: make(map[string]types.DB),
			defaultConn: "default",
		}
		
		// Automatically load environment configuration
		if err := builderInstance.loadEnvironmentConfig(); err != nil {
			log.Printf("Warning: Failed to load environment config: %v", err)
			log.Printf("You can manually add connections using AddConnection()")
		}
	})
	return builderInstance
}

// loadEnvironmentConfig loads configuration from environment variables
func (b *Builder) loadEnvironmentConfig() error {
	cfg, err := config.LoadFromEnv()
	if err != nil {
		return err
	}
	
	// Create default connection
	db, err := database.NewConnection(cfg)
	if err != nil {
		return err
	}
	
	b.mu.Lock()
	b.connections["default"] = db
	b.defaultConn = "default"
	b.mu.Unlock()
	
	return nil
}

// AddConnection adds a named database connection
func (b *Builder) AddConnection(name string, config types.Config) error {
	db, err := database.NewConnection(config)
	if err != nil {
		return err
	}
	
	b.mu.Lock()
	defer b.mu.Unlock()
	
	b.connections[name] = db
	if len(b.connections) == 1 {
		b.defaultConn = name
	}
	
	return nil
}

// Connection switches to a different database connection
func (b *Builder) Connection(name string) *Builder {
	// Create a copy to avoid affecting the singleton state
	newBuilder := &Builder{
		connections: b.connections,
		defaultConn: name,
		mu:          sync.RWMutex{},
	}
	
	return newBuilder
}

// Table creates a query builder for a table using model, pointer, or string
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

// extractTableName extracts table name from various input types
func (b *Builder) extractTableName(table interface{}) string {
	switch t := table.(type) {
	case string:
		return t
	case Tabler:
		return t.TableName()
	default:
		// Try to extract table name from struct type
		return b.deriveTableNameFromType(t)
	}
}

// deriveTableNameFromType converts struct type to table name
func (b *Builder) deriveTableNameFromType(model interface{}) string {
	t := reflect.TypeOf(model)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	
	// Convert struct name to snake_case and pluralize
	name := t.Name()
	snakeName := b.toSnakeCase(name)
	return b.pluralize(snakeName)
}

// toSnakeCase converts CamelCase to snake_case
func (b *Builder) toSnakeCase(str string) string {
	var result strings.Builder
	for i, r := range str {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteByte('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}

// pluralize adds 's' to make table name plural (simple implementation)
func (b *Builder) pluralize(str string) string {
	if strings.HasSuffix(str, "y") {
		return strings.TrimSuffix(str, "y") + "ies"
	}
	if strings.HasSuffix(str, "s") || strings.HasSuffix(str, "x") || strings.HasSuffix(str, "z") {
		return str + "es"
	}
	return str + "s"
}

// getConnectionNames returns list of available connection names
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
	
	for _, db := range b.connections {
		if err := db.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Legacy functions for backward compatibility
func NewConnection(config types.Config) (types.DB, error) {
	return database.NewConnection(config)
}

func NewQueryBuilder(executor types.QueryExecutor, driver types.Driver) types.QueryBuilder {
	return query.NewBuilder(executor, driver)
}

func Table(executor types.QueryExecutor, driver types.Driver, table string) types.QueryBuilder {
	return query.Table(executor, driver, table)
}

// Convenience functions that use the singleton builder
func QB() *Builder {
	return GetBuilder()
}

func TableBuilder(table interface{}) types.QueryBuilder {
	return GetBuilder().Table(table)
}

func Connection(name string) *Builder {
	return GetBuilder().Connection(name)
}

type Config = types.Config
type Driver = types.Driver
type Collection = types.Collection
type QueryBuilder = types.QueryBuilder
type PaginationResult = types.PaginationResult
type AggregateResult = types.AggregateResult
type DebugInfo = types.DebugInfo
type OrderDirection = types.OrderDirection

const (
	MySQL      = types.MySQL
	PostgreSQL = types.PostgreSQL
	Asc        = types.Asc
	Desc       = types.Desc
)

var (
	NewCollection = types.NewCollection
)