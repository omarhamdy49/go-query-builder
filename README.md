# Go Query Builder

A powerful, Laravel-inspired query builder for Go, providing a fluent interface for database operations. Built with clean architecture, security, and performance in mind.

## Features

### 🚀 **Core Features**
- **Laravel-inspired API**: Familiar syntax for Laravel developers
- **Multiple Database Support**: MySQL and PostgreSQL with optimized drivers
- **Type Safety**: Full Go type safety with interfaces
- **Clean Architecture**: Well-organized packages following SOLID principles
- **Security First**: Built-in SQL injection protection and input validation

### 📊 **Query Operations**
- **SELECT**: Get, First, Find, Pluck with full collection support
- **WHERE**: Complex conditions, JSON queries, date helpers, full-text search
- **JOINS**: Inner, Left, Right, Cross joins with multiple conditions
- **AGGREGATES**: Count, Sum, Avg, Min, Max
- **UNIONS**: Union and Union All with subqueries
- **ORDERING/GROUPING**: Multiple columns, raw expressions
- **LIMITING**: Limit, Offset, Take, Skip

### 🔧 **Advanced Features**
- **JSON Queries**: JSON path queries, contains, length operations
- **Date/Time Helpers**: WhereToday, WhereMonth, WherePast, etc.
- **Full-Text Search**: MySQL MATCH AGAINST and PostgreSQL tsvector
- **UPSERT Operations**: MySQL ON DUPLICATE KEY, PostgreSQL ON CONFLICT
- **Chunking & Streaming**: Memory-efficient large dataset processing
- **Pagination**: Length-aware, simple, and cursor-based pagination
- **Transactions**: Full transaction support with rollback/commit
- **Conditional Building**: When, Unless, Tap for dynamic queries

### 🛡️ **Security & Performance**
- **SQL Injection Protection**: Parameterized queries and input validation
- **Connection Pooling**: Optimized connection management
- **Lazy Loading**: Efficient memory usage for large datasets
- **Debug Support**: Query logging and performance monitoring
- **Context Support**: Proper context handling for timeouts/cancellation

## Installation

```bash
go get github.com/go-query-builder/querybuilder
```

## Quick Start

### Database Connection

```go
package main

import (
    "context"
    "log"
    
    "github.com/go-query-builder/querybuilder"
)

func main() {
    config := querybuilder.Config{
        Driver:   querybuilder.MySQL, // or querybuilder.PostgreSQL
        Host:     "localhost",
        Port:     3306,
        Database: "mydb",
        Username: "user",
        Password: "password",
    }

    db, err := querybuilder.NewConnection(config)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    ctx := context.Background()
    
    // Your queries here...
}
```

### Basic Queries

```go
// Get all users
users, err := querybuilder.Table(db, db.Driver(), "users").Get(ctx)

// Get specific user
user, err := querybuilder.Table(db, db.Driver(), "users").
    Where("id", 1).
    First(ctx)

// Get user names
names, err := querybuilder.Table(db, db.Driver(), "users").
    Pluck(ctx, "name")

// Count users
count, err := querybuilder.Table(db, db.Driver(), "users").
    Where("active", true).
    Count(ctx)
```

### WHERE Clauses

```go
qb := querybuilder.Table(db, db.Driver(), "users")

// Basic WHERE
users, err := qb.Where("age", ">", 18).
    Where("status", "active").
    OrWhere("role", "admin").
    Get(ctx)

// WHERE IN
users, err := qb.WhereIn("role", []interface{}{"admin", "editor", "user"}).Get(ctx)

// WHERE BETWEEN
users, err := qb.WhereBetween("age", []interface{}{18, 65}).Get(ctx)

// WHERE NULL
users, err := qb.WhereNull("deleted_at").Get(ctx)

// Raw WHERE
users, err := qb.WhereRaw("YEAR(created_at) = ?", 2023).Get(ctx)
```

### Date/Time Queries

```go
qb := querybuilder.Table(db, db.Driver(), "orders")

// Today's orders
orders, err := qb.WhereToday("created_at").Get(ctx)

// This month's orders
orders, err := qb.WhereMonth("created_at", time.Now().Month()).Get(ctx)

// Past orders
orders, err := qb.WherePast("shipped_at").Get(ctx)

// Orders from specific year
orders, err := qb.WhereYear("created_at", 2023).Get(ctx)
```

### JSON Queries

```go
qb := querybuilder.Table(db, db.Driver(), "products")

// JSON contains
products, err := qb.WhereJsonContains("metadata", map[string]interface{}{
    "featured": true,
}).Get(ctx)

// JSON path query
products, err := qb.WhereJsonPath("settings", "$.notifications.email", true).Get(ctx)

// JSON array length
products, err := qb.WhereJsonLength("tags", ">", 3).Get(ctx)
```

### JOINS

```go
// Inner Join
users, err := querybuilder.Table(db, db.Driver(), "users").
    Select("users.name", "profiles.bio").
    Join("profiles", "users.id", "profiles.user_id").
    Get(ctx)

// Left Join with conditions
users, err := querybuilder.Table(db, db.Driver(), "users").
    LeftJoin("orders", "users.id", "orders.user_id").
    Where("orders.status", "completed").
    Get(ctx)

// Multiple joins
posts, err := querybuilder.Table(db, db.Driver(), "posts").
    Join("users", "posts.author_id", "users.id").
    LeftJoin("categories", "posts.category_id", "categories.id").
    Select("posts.title", "users.name", "categories.name").
    Get(ctx)
```

### Aggregates

```go
qb := querybuilder.Table(db, db.Driver(), "orders")

// Count
count, err := qb.Where("status", "completed").Count(ctx)

// Sum
total, err := qb.Sum(ctx, "amount")

// Average
avg, err := qb.Avg(ctx, "amount")

// Min/Max
min, err := qb.Min(ctx, "amount")
max, err := qb.Max(ctx, "amount")
```

### INSERT Operations

```go
qb := querybuilder.Table(db, db.Driver(), "users")

// Single insert
err := qb.Insert(ctx, map[string]interface{}{
    "name":     "John Doe",
    "email":    "john@example.com",
    "age":      30,
})

// Batch insert
users := []map[string]interface{}{
    {"name": "Jane", "email": "jane@example.com"},
    {"name": "Bob", "email": "bob@example.com"},
}
err := qb.InsertBatch(ctx, users)
```

### UPDATE Operations

```go
qb := querybuilder.Table(db, db.Driver(), "users")

// Update specific records
rowsAffected, err := qb.Where("id", 1).Update(ctx, map[string]interface{}{
    "name":       "John Updated",
    "updated_at": time.Now(),
})

// Increment/Decrement (would be available through execution package)
// rowsAffected, err := executor.Increment(ctx, qb, "views", 1)
```

### DELETE Operations

```go
qb := querybuilder.Table(db, db.Driver(), "users")

// Delete with conditions
rowsAffected, err := qb.Where("status", "inactive").Delete(ctx)

// Delete all (be careful!)
rowsAffected, err := qb.Delete(ctx)
```

### Transactions

```go
tx, err := db.Begin()
if err != nil {
    return err
}
defer func() {
    if r := recover(); r != nil {
        tx.Rollback()
    }
}()

// Perform operations within transaction
qb := querybuilder.Table(tx, db.Driver(), "accounts")

_, err = qb.Where("id", 1).Update(ctx, map[string]interface{}{
    "balance": 1000,
})
if err != nil {
    tx.Rollback()
    return err
}

err = tx.Commit()
```

### Pagination

```go
import "github.com/go-query-builder/querybuilder/pkg/pagination"

paginator := pagination.NewPaginator(db, db.Driver())
qb := querybuilder.Table(db, db.Driver(), "posts")

// Length-aware pagination
result, err := paginator.Paginate(ctx, qb, 1, 10) // page 1, 10 per page
if err != nil {
    return err
}

fmt.Printf("Page %d of %d (%d total items)\n", 
    result.CurrentPage, result.LastPage, result.Total)

// Simple pagination (no total count)
simple, err := paginator.SimplePaginate(ctx, qb, 1, 10)
if err != nil {
    return err
}

fmt.Printf("Has more pages: %t\n", simple.HasMore)
```

### Chunking Large Datasets

```go
import "github.com/go-query-builder/querybuilder/pkg/execution"

executor := execution.NewQueryExecutor(db, db.Driver())
qb := querybuilder.Table(db, db.Driver(), "large_table")

// Process in chunks of 1000
err := executor.Chunk(ctx, qb, 1000, func(collection querybuilder.Collection) error {
    fmt.Printf("Processing %d records\n", collection.Count())
    
    // Process each record
    for _, record := range collection.ToSlice() {
        // Process record...
    }
    
    return nil
})
```

### Conditional Query Building

```go
qb := querybuilder.Table(db, db.Driver(), "users")

includeInactive := true
filterByAge := false

users, err := qb.
    When(includeInactive, func(q querybuilder.QueryBuilder) querybuilder.QueryBuilder {
        return q.OrWhere("status", "inactive")
    }).
    Unless(filterByAge, func(q querybuilder.QueryBuilder) querybuilder.QueryBuilder {
        return q.Where("age", ">=", 18)
    }).
    Tap(func(q querybuilder.QueryBuilder) querybuilder.QueryBuilder {
        fmt.Println("Query building completed")
        return q
    }).
    Get(ctx)
```

### Raw Queries

```go
qb := querybuilder.Table(db, db.Driver(), "analytics")

results, err := qb.
    SelectRaw("DATE(created_at) as date, COUNT(*) as count").
    WhereRaw("created_at >= ?", time.Now().AddDate(0, -1, 0)).
    GroupByRaw("DATE(created_at)").
    HavingRaw("COUNT(*) > ?", 10).
    OrderByRaw("date DESC").
    Get(ctx)
```

### Debug and Monitoring

```go
qb := querybuilder.Table(db, db.Driver(), "users").Debug()

users, err := qb.Where("active", true).Get(ctx)

// Get debug information
if debugInfo := qb.GetDebugInfo(); debugInfo != nil {
    fmt.Printf("SQL: %s\n", debugInfo.SQL)
    fmt.Printf("Bindings: %+v\n", debugInfo.Bindings)
    fmt.Printf("Duration: %v\n", debugInfo.Duration)
}
```

## Database Support

### MySQL
```go
config := querybuilder.Config{
    Driver:   querybuilder.MySQL,
    Host:     "localhost",
    Port:     3306,
    Database: "mydb",
    Username: "user",
    Password: "password",
    Charset:  "utf8mb4",
}
```

### PostgreSQL
```go
config := querybuilder.Config{
    Driver:   querybuilder.PostgreSQL,
    Host:     "localhost",
    Port:     5432,
    Database: "mydb",
    Username: "user",
    Password: "password",
    SSLMode:  "disable",
}
```

## Security Features

The query builder includes built-in security features:

- **Parameterized Queries**: All user inputs are properly escaped
- **Input Validation**: Automatic validation of table names, column names, and values
- **SQL Injection Prevention**: Protection against common attack vectors
- **Query Length Limits**: Configurable limits to prevent abuse

```go
import "github.com/go-query-builder/querybuilder/pkg/security"

validator := security.NewSecurityValidator()

// Validate table name
err := validator.ValidateTableName("users")

// Validate column name  
err := validator.ValidateColumnName("email")

// Validate raw SQL
err := validator.ValidateRawSQL("SELECT * FROM users WHERE id = ?")
```

## Architecture

The package follows clean architecture principles:

```
pkg/
├── types/          # Interfaces and type definitions
├── database/       # Database connection management
├── query/          # Core query builder implementation
├── clauses/        # Query clause structures
├── execution/      # Query execution logic
├── pagination/     # Pagination utilities
└── security/       # Security validation
```

## Contributing

Contributions are welcome! Please ensure:

1. **Tests**: Write comprehensive tests for new features
2. **Documentation**: Update documentation for API changes
3. **Security**: Follow security best practices
4. **Performance**: Consider performance implications
5. **Compatibility**: Maintain backward compatibility when possible

## License

MIT License - see LICENSE file for details.

## Acknowledgments

- Inspired by Laravel's Eloquent Query Builder
- Built with security and performance best practices
- Follows Go idioms and conventions