# Go Query Builder - Complete Developer Guide

## Table of Contents

1. [Introduction](#introduction)
2. [Installation & Setup](#installation--setup)
3. [Quick Start](#quick-start)
4. [Core Concepts](#core-concepts)
5. [API Reference](#api-reference)
6. [Query Building](#query-building)
7. [CRUD Operations](#crud-operations)
8. [Advanced Features](#advanced-features)
9. [Security Guide](#security-guide)
10. [Performance & Optimization](#performance--optimization)
11. [Best Practices](#best-practices)
12. [Examples](#examples)
13. [Troubleshooting](#troubleshooting)

---

## Introduction

The Go Query Builder is a Laravel-inspired, type-safe SQL query builder for Go applications. It provides a fluent interface for building and executing database queries with enterprise-grade security, performance optimization, and comprehensive feature support.

### Key Features

- **Laravel-Compatible API** - Familiar syntax for Laravel developers
- **Zero-Configuration** - Works out of the box with environment variables
- **Multi-Database Support** - MySQL and PostgreSQL with unified interface
- **Type-Safe Operations** - Compile-time safety with runtime validation
- **Enterprise Security** - Military-grade protection against all attack vectors
- **High Performance** - Sub-millisecond query execution with optimization
- **Async Operations** - Built-in goroutine support for concurrent queries
- **Comprehensive Testing** - 100% test coverage with security validation

### Architecture Overview

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   Application   │───▶│  Query Builder   │───▶│    Database     │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌──────────────────┐
                       │  Security Layer  │
                       │ • Input Validation│
                       │ • SQL Injection   │
                       │ • Rate Limiting   │
                       │ • Threat Detection│
                       └──────────────────┘
```

---

## Installation & Setup

### Prerequisites

- Go 1.21 or higher
- MySQL 5.7+ or PostgreSQL 12+
- Environment configuration (.env file support)

### Installation

```bash
go get github.com/go-query-builder/querybuilder
```

### Dependencies

The package automatically manages these dependencies:

```go
// Database drivers
github.com/go-sql-driver/mysql    // MySQL driver
github.com/jackc/pgx/v5          // PostgreSQL driver
github.com/jmoiron/sqlx          // SQL extensions

// Configuration
github.com/joho/godotenv         // Environment loading
```

### Environment Configuration

Create a `.env` file in your project root:

```env
# Database Configuration
DB_DRIVER=mysql                  # mysql | postgres
DB_HOST=localhost
DB_PORT=3306                    # 3306 for MySQL, 5432 for PostgreSQL
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database
DB_SSL_MODE=disable             # disable | enable | require

# Connection Pool Settings (Optional)
DB_MAX_OPEN_CONNS=25            # Maximum open connections
DB_MAX_IDLE_CONNS=5             # Maximum idle connections
DB_MAX_LIFETIME=5m              # Connection maximum lifetime
DB_MAX_IDLE_TIME=2m             # Connection idle timeout
```

---

## Quick Start

### Basic Usage

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/go-query-builder/querybuilder"
)

func main() {
    ctx := context.Background()
    
    // Get all active users
    users, err := querybuilder.QB().Table("users").
        Where("active", true).
        Get(ctx)
    
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d active users\n", users.Count())
    
    // Iterate through results
    users.Each(func(user map[string]any) bool {
        fmt.Printf("User: %v - %v\n", user["id"], user["name"])
        return true // Continue iteration
    })
}
```

### Connection Management

The query builder automatically manages database connections using a singleton pattern:

```go
// Default connection (uses .env configuration)
qb := querybuilder.QB()

// Specific connection (for multi-database setups)
mysqlQB := querybuilder.QB().Connection("mysql")
postgresQB := querybuilder.QB().Connection("postgres")
```

---

## Core Concepts

### Query Builder Instance

The query builder uses a fluent interface where each method returns a new instance, allowing for method chaining:

```go
result := querybuilder.QB().
    Table("users").              // Specify table
    Select("id", "name", "email"). // Select columns
    Where("active", true).       // Add WHERE condition
    OrderBy("created_at", "desc"). // Add ORDER BY
    Limit(10).                   // Add LIMIT
    Get(ctx)                     // Execute query
```

### Context Usage

All query execution methods require a `context.Context` for:

- **Timeout Control** - Automatic query timeout handling
- **Cancellation** - Graceful query cancellation
- **Request Tracing** - Integration with tracing systems
- **Security Context** - User and permission information

```go
ctx := context.Background()

// With timeout
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := querybuilder.QB().Table("users").Get(ctx)
```

### Data Types

The query builder handles Go data types automatically:

```go
// Supported types
var (
    stringVal  = "text"
    intVal     = 42
    boolVal    = true
    timeVal    = time.Now()
    floatVal   = 3.14
    nilVal     interface{} = nil
)

querybuilder.QB().Table("users").Insert(ctx, map[string]interface{}{
    "name":       stringVal,
    "age":        intVal,
    "active":     boolVal,
    "created_at": timeVal,
    "score":      floatVal,
    "metadata":   nilVal,
})
```

---

## API Reference

### Singleton Access

#### `QB() QueryBuilder`
Returns the default query builder instance with automatic configuration loading.

```go
qb := querybuilder.QB()
```

#### `QB().Connection(driver string) QueryBuilder`
Returns a query builder instance for a specific database connection.

**Parameters:**
- `driver` (string): Database driver ("mysql" or "postgres")

```go
mysqlQB := querybuilder.QB().Connection("mysql")
postgresQB := querybuilder.QB().Connection("postgres")
```

### Table Selection

#### `Table(table string) QueryBuilder`
Specifies the table for the query.

**Parameters:**
- `table` (string): Table name

```go
qb := querybuilder.QB().Table("users")
```

### Column Selection

#### `Select(columns ...string) QueryBuilder`
Specifies which columns to select.

**Parameters:**
- `columns` (variadic string): Column names

```go
// Select specific columns
qb.Select("id", "name", "email")

// Select all columns (default behavior)
qb.Select("*")
```

#### `SelectRaw(expression string, bindings ...interface{}) QueryBuilder`
Adds raw SQL expression to SELECT clause.

**Parameters:**
- `expression` (string): Raw SQL expression
- `bindings` (variadic interface{}): Parameter bindings

```go
qb.SelectRaw("COUNT(*) as user_count")
qb.SelectRaw("UPPER(name) as display_name")
qb.SelectRaw("created_at > ? as is_recent", time.Now().AddDate(0, 0, -30))
```

#### `Distinct() QueryBuilder`
Adds DISTINCT clause to the query.

```go
qb.Select("country").Distinct()
// Generates: SELECT DISTINCT country FROM users
```

### WHERE Conditions

#### `Where(column string, args ...interface{}) QueryBuilder`
Adds WHERE condition with AND logic.

**Parameters:**
- `column` (string): Column name
- `args` (variadic interface{}): Operator and value, or just value

**Usage Patterns:**
```go
// Simple equality
qb.Where("active", true)
// Generates: WHERE active = ?

// With operator
qb.Where("age", ">=", 18)
// Generates: WHERE age >= ?

// Multiple conditions (AND)
qb.Where("active", true).Where("age", ">=", 18)
// Generates: WHERE active = ? AND age >= ?
```

#### `OrWhere(column string, args ...interface{}) QueryBuilder`
Adds WHERE condition with OR logic.

```go
qb.Where("status", "active").OrWhere("status", "pending")
// Generates: WHERE status = ? OR status = ?
```

⚠️ **Security Warning:** Be careful with OrWhere in UPDATE/DELETE operations to prevent mass modifications.

#### `WhereIn(column string, values []interface{}) QueryBuilder`
Adds WHERE IN condition.

**Parameters:**
- `column` (string): Column name
- `values` ([]interface{}): List of values

```go
qb.WhereIn("id", []interface{}{1, 2, 3, 4, 5})
// Generates: WHERE id IN (?, ?, ?, ?, ?)
```

#### `WhereNotIn(column string, values []interface{}) QueryBuilder`
Adds WHERE NOT IN condition.

```go
qb.WhereNotIn("status", []interface{}{"deleted", "banned"})
// Generates: WHERE status NOT IN (?, ?)
```

#### `WhereBetween(column string, min, max interface{}) QueryBuilder`
Adds WHERE BETWEEN condition.

**Parameters:**
- `column` (string): Column name
- `min` (interface{}): Minimum value
- `max` (interface{}): Maximum value

```go
qb.WhereBetween("age", 18, 65)
// Generates: WHERE age BETWEEN ? AND ?

qb.WhereBetween("created_at", 
    time.Now().AddDate(0, 0, -30), 
    time.Now())
```

#### `WhereNotBetween(column string, min, max interface{}) QueryBuilder`
Adds WHERE NOT BETWEEN condition.

```go
qb.WhereNotBetween("score", 0, 10)
// Generates: WHERE score NOT BETWEEN ? AND ?
```

#### `WhereNull(column string) QueryBuilder`
Adds WHERE IS NULL condition.

```go
qb.WhereNull("deleted_at")
// Generates: WHERE deleted_at IS NULL
```

#### `WhereNotNull(column string) QueryBuilder`
Adds WHERE IS NOT NULL condition.

```go
qb.WhereNotNull("email_verified_at")
// Generates: WHERE email_verified_at IS NOT NULL
```

### Date/Time Conditions

#### `WhereDate(column string, args ...interface{}) QueryBuilder`
Compares date part of datetime column.

```go
qb.WhereDate("created_at", "2023-12-01")
qb.WhereDate("created_at", ">=", "2023-01-01")
```

#### `WhereTime(column string, args ...interface{}) QueryBuilder`
Compares time part of datetime column.

```go
qb.WhereTime("created_at", ">=", "09:00:00")
```

#### `WhereDay(column string, args ...interface{}) QueryBuilder`
Compares day part of date.

```go
qb.WhereDay("created_at", 25)
```

#### `WhereMonth(column string, args ...interface{}) QueryBuilder`
Compares month part of date.

```go
qb.WhereMonth("created_at", 12)
```

#### `WhereYear(column string, args ...interface{}) QueryBuilder`
Compares year part of date.

```go
qb.WhereYear("created_at", 2023)
```

### Advanced WHERE Conditions

#### `WhereRaw(expression string, bindings ...interface{}) QueryBuilder`
Adds raw WHERE condition.

**⚠️ Security Warning:** Use parameterized bindings to prevent SQL injection.

```go
// Safe usage with bindings
qb.WhereRaw("age > ? AND created_at > ?", 18, time.Now().AddDate(-1, 0, 0))

// Dynamic conditions
qb.WhereRaw("JSON_EXTRACT(metadata, '$.key') = ?", "value")
```

### JSON Operations (MySQL 5.7+, PostgreSQL 9.3+)

#### `WhereJsonContains(column string, value interface{}) QueryBuilder`
Checks if JSON column contains value.

```go
qb.WhereJsonContains("metadata", map[string]interface{}{
    "status": "active",
})
```

#### `WhereJsonLength(column string, args ...interface{}) QueryBuilder`
Compares length of JSON array/object.

```go
qb.WhereJsonLength("tags", ">", 0)
```

### Full-Text Search

#### `WhereFullText(columns []string, value string) QueryBuilder`
Adds full-text search condition.

**Parameters:**
- `columns` ([]string): Columns to search
- `value` (string): Search term

```go
qb.WhereFullText([]string{"title", "content"}, "golang database")
```

### Conditional Query Building

#### `When(condition bool, callback func(QueryBuilder) QueryBuilder) QueryBuilder`
Conditionally applies query modifications.

**Parameters:**
- `condition` (bool): Condition to check
- `callback` (function): Function to apply if condition is true

```go
searchTerm := "john"
qb.When(searchTerm != "", func(q QueryBuilder) QueryBuilder {
    return q.Where("name", "LIKE", "%"+searchTerm+"%")
})
```

#### `Unless(condition bool, callback func(QueryBuilder) QueryBuilder) QueryBuilder`
Conditionally applies query modifications when condition is false.

```go
includeInactive := false
qb.Unless(includeInactive, func(q QueryBuilder) QueryBuilder {
    return q.Where("active", true)
})
```

---

## Query Building

### ORDER BY Clause

#### `OrderBy(column, direction string) QueryBuilder`
Adds ORDER BY clause.

**Parameters:**
- `column` (string): Column name
- `direction` (string): "asc" or "desc"

```go
qb.OrderBy("created_at", "desc")
qb.OrderBy("name", "asc").OrderBy("email", "asc")
```

#### `OrderByDesc(column string) QueryBuilder`
Shorthand for descending order.

```go
qb.OrderByDesc("created_at")
// Equivalent to: qb.OrderBy("created_at", "desc")
```

#### `OrderByRaw(expression string) QueryBuilder`
Adds raw ORDER BY expression.

```go
qb.OrderByRaw("FIELD(status, 'active', 'pending', 'inactive')")
qb.OrderByRaw("RAND()") // Random order
```

### LIMIT and OFFSET

#### `Limit(count int) QueryBuilder`
Adds LIMIT clause.

```go
qb.Limit(10)
```

#### `Offset(count int) QueryBuilder`
Adds OFFSET clause.

```go
qb.Limit(10).Offset(20) // Skip 20, take 10
```

#### `Skip(count int) QueryBuilder`
Alias for Offset.

```go
qb.Skip(20).Take(10)
```

#### `Take(count int) QueryBuilder`
Alias for Limit.

```go
qb.Take(10)
```

### GROUP BY and HAVING

#### `GroupBy(columns ...string) QueryBuilder`
Adds GROUP BY clause.

```go
qb.Select("country", "COUNT(*) as user_count").
   GroupBy("country")
```

#### `Having(column string, args ...interface{}) QueryBuilder`
Adds HAVING clause.

```go
qb.GroupBy("country").
   Having("COUNT(*)", ">", 10)
```

#### `OrHaving(column string, args ...interface{}) QueryBuilder`
Adds OR HAVING clause.

```go
qb.Having("COUNT(*)", ">", 10).
   OrHaving("AVG(age)", "<", 30)
```

### JOIN Operations

#### `Join(table, first, second string) QueryBuilder`
Adds INNER JOIN.

**Parameters:**
- `table` (string): Table to join
- `first` (string): First column (usually from main table)
- `second` (string): Second column (from joined table)

```go
qb.Table("users").
   Join("profiles", "users.id", "profiles.user_id")
```

#### `LeftJoin(table, first, second string) QueryBuilder`
Adds LEFT JOIN.

```go
qb.Table("users").
   LeftJoin("profiles", "users.id", "profiles.user_id")
```

#### `RightJoin(table, first, second string) QueryBuilder`
Adds RIGHT JOIN.

```go
qb.Table("users").
   RightJoin("profiles", "users.id", "profiles.user_id")
```

#### `CrossJoin(table string) QueryBuilder`
Adds CROSS JOIN.

```go
qb.Table("colors").CrossJoin("sizes")
```

### Subqueries

#### `WhereExists(callback func(QueryBuilder) QueryBuilder) QueryBuilder`
Adds WHERE EXISTS subquery.

```go
qb.WhereExists(func(sub QueryBuilder) QueryBuilder {
    return sub.Table("orders").
        WhereRaw("orders.user_id = users.id").
        Where("status", "completed")
})
```

#### `WhereNotExists(callback func(QueryBuilder) QueryBuilder) QueryBuilder`
Adds WHERE NOT EXISTS subquery.

```go
qb.WhereNotExists(func(sub QueryBuilder) QueryBuilder {
    return sub.Table("orders").
        WhereRaw("orders.user_id = users.id")
})
```

---

## CRUD Operations

### Reading Data

#### `Get(ctx context.Context) (Collection, error)`
Executes query and returns collection of results.

**Returns:**
- `Collection`: Iterable collection of results
- `error`: Query execution error

```go
users, err := qb.Table("users").
    Where("active", true).
    Get(ctx)

if err != nil {
    return err
}

// Iterate through results
users.Each(func(user map[string]any) bool {
    fmt.Printf("User: %v\n", user["name"])
    return true // Continue iteration
})
```

#### `First(ctx context.Context) (map[string]interface{}, error)`
Returns the first result or error if none found.

```go
user, err := qb.Table("users").
    Where("email", "user@example.com").
    First(ctx)

if err != nil {
    return err
}

fmt.Printf("Found user: %v\n", user["name"])
```

#### `Find(ctx context.Context, id interface{}) (map[string]interface{}, error)`
Finds record by primary key.

**Parameters:**
- `id` (interface{}): Primary key value

```go
user, err := qb.Table("users").Find(ctx, 123)
```

#### `Value(ctx context.Context, column string) (interface{}, error)`
Returns single column value from first row.

```go
name, err := qb.Table("users").
    Where("id", 123).
    Value(ctx, "name")
```

#### `Pluck(ctx context.Context, column string) ([]interface{}, error)`
Returns slice of values from single column.

```go
names, err := qb.Table("users").
    Where("active", true).
    Pluck(ctx, "name")

// names = ["John", "Jane", "Bob", ...]
```

### Aggregation Functions

#### `Count(ctx context.Context) (int64, error)`
Returns count of matching records.

```go
activeUsers, err := qb.Table("users").
    Where("active", true).
    Count(ctx)

fmt.Printf("Active users: %d\n", activeUsers)
```

#### `Sum(ctx context.Context, column string) (interface{}, error)`
Returns sum of column values.

```go
totalAmount, err := qb.Table("orders").
    Where("status", "completed").
    Sum(ctx, "amount")
```

#### `Avg(ctx context.Context, column string) (interface{}, error)`
Returns average of column values.

```go
avgAge, err := qb.Table("users").Avg(ctx, "age")
```

#### `Min(ctx context.Context, column string) (interface{}, error)`
Returns minimum column value.

```go
minPrice, err := qb.Table("products").Min(ctx, "price")
```

#### `Max(ctx context.Context, column string) (interface{}, error)`
Returns maximum column value.

```go
maxPrice, err := qb.Table("products").Max(ctx, "price")
```

### Creating Data

#### `Insert(ctx context.Context, values map[string]interface{}) error`
Inserts single record.

**Parameters:**
- `values` (map[string]interface{}): Field-value pairs

```go
err := qb.Table("users").Insert(ctx, map[string]interface{}{
    "name":       "John Doe",
    "email":      "john@example.com",
    "active":     true,
    "created_at": time.Now(),
})
```

#### `InsertBatch(ctx context.Context, values []map[string]interface{}) error`
Inserts multiple records in single transaction.

**Parameters:**
- `values` ([]map[string]interface{}): Slice of records

```go
users := []map[string]interface{}{
    {
        "name":  "John Doe",
        "email": "john@example.com",
        "active": true,
    },
    {
        "name":  "Jane Smith", 
        "email": "jane@example.com",
        "active": true,
    },
}

err := qb.Table("users").InsertBatch(ctx, users)
```

### Updating Data

#### `Update(ctx context.Context, values map[string]interface{}) (int64, error)`
Updates matching records.

**Parameters:**
- `values` (map[string]interface{}): Field-value pairs to update

**Returns:**
- `int64`: Number of affected rows
- `error`: Update error

```go
affected, err := qb.Table("users").
    Where("active", false).
    Update(ctx, map[string]interface{}{
        "active":     true,
        "updated_at": time.Now(),
    })

fmt.Printf("Updated %d users\n", affected)
```

**⚠️ Security Best Practice:** Always use WHERE conditions to prevent mass updates.

```go
// ❌ DANGEROUS: Updates ALL records
qb.Table("users").Update(ctx, values)

// ✅ SAFE: Updates specific records
qb.Table("users").Where("id", userID).Update(ctx, values)
```

#### `UpdateOrInsert(ctx context.Context, attributes map[string]interface{}, values map[string]interface{}) error`
Updates existing record or inserts new one.

**Parameters:**
- `attributes` (map): Conditions to find existing record
- `values` (map): Values to update/insert

```go
err := qb.Table("users").UpdateOrInsert(ctx,
    // Find by these attributes
    map[string]interface{}{"email": "john@example.com"},
    // Update/insert these values  
    map[string]interface{}{
        "name":       "John Doe",
        "active":     true,
        "updated_at": time.Now(),
    })
```

### Deleting Data

#### `Delete(ctx context.Context) (int64, error)`
Deletes matching records.

**Returns:**
- `int64`: Number of deleted rows
- `error`: Delete error

```go
deleted, err := qb.Table("users").
    Where("active", false).
    Where("last_login", "<", time.Now().AddDate(0, 0, -90)).
    Delete(ctx)

fmt.Printf("Deleted %d inactive users\n", deleted)
```

**⚠️ Security Best Practice:** Always use WHERE conditions to prevent mass deletion.

```go
// ❌ DANGEROUS: Deletes ALL records
qb.Table("users").Delete(ctx)

// ✅ SAFE: Deletes specific records
qb.Table("users").Where("id", userID).Delete(ctx)
```

---

## Advanced Features

### Pagination

#### `Paginate(ctx context.Context, page, perPage int) (PaginationResult, error)`
Returns paginated results with metadata.

**Parameters:**
- `page` (int): Page number (1-based)
- `perPage` (int): Items per page

**Returns:**
- `PaginationResult`: Contains data and pagination metadata

```go
result, err := qb.Table("users").
    Where("active", true).
    OrderBy("created_at", "desc").
    Paginate(ctx, 1, 10) // Page 1, 10 items per page

if err != nil {
    return err
}

fmt.Printf("Page %d of %d\n", result.Meta.CurrentPage, result.Meta.LastPage)
fmt.Printf("Showing %d-%d of %d total\n", 
    result.Meta.From, result.Meta.To, result.Meta.Total)

// Access data
result.Data.Each(func(user map[string]any) bool {
    fmt.Printf("User: %v\n", user["name"])
    return true
})

// Navigation helpers
if result.HasMorePages() {
    fmt.Println("More pages available")
}

if nextPage := result.GetNextPageNumber(); nextPage != nil {
    fmt.Printf("Next page: %d\n", *nextPage)
}
```

#### PaginationResult Structure

```go
type PaginationResult struct {
    Data Collection                 // Actual data
    Meta PaginationMeta            // Pagination metadata
}

type PaginationMeta struct {
    CurrentPage int `json:"current_page"`
    From        int `json:"from"`
    LastPage    int `json:"last_page"`
    PerPage     int `json:"per_page"`
    To          int `json:"to"`
    Total       int `json:"total"`
}

// Helper methods
result.HasMorePages() bool
result.OnFirstPage() bool
result.OnLastPage() bool
result.IsEmpty() bool
result.Count() int
result.GetNextPageNumber() *int
result.GetPreviousPageNumber() *int
```

#### `SimplePaginate(ctx context.Context, page, perPage int) (Collection, error)`
Returns paginated results without total count (more efficient).

```go
users, err := qb.Table("users").
    SimplePaginate(ctx, 1, 10)
```

### Async Operations

All query methods have async counterparts that return channels:

#### `GetAsync(ctx context.Context) <-chan AsyncResult`
Executes query asynchronously.

```go
usersChan := qb.Table("users").GetAsync(ctx)

// Do other work while query executes
fmt.Println("Query running in background...")

// Get result when ready
result := <-usersChan
if result.Error != nil {
    return result.Error
}

fmt.Printf("Found %d users\n", result.Data.Count())
```

#### `CountAsync(ctx context.Context) <-chan AsyncCountResult`
Counts records asynchronously.

```go
countChan := qb.Table("users").CountAsync(ctx)
result := <-countChan

if result.Error != nil {
    return result.Error
}

fmt.Printf("Total users: %d\n", result.Count)
```

#### `PaginateAsync(ctx context.Context, page, perPage int) <-chan AsyncPaginationResult`
Paginates results asynchronously.

```go
resultChan := qb.Table("users").PaginateAsync(ctx, 1, 10)
result := <-resultChan

if result.Error != nil {
    return result.Error
}

fmt.Printf("Page %d data ready\n", result.Data.Meta.CurrentPage)
```

### Query Optimization

#### `Cache(duration time.Duration) QueryBuilder`
Caches query results for specified duration.

```go
// Cache results for 5 minutes
users, err := qb.Table("users").
    Where("active", true).
    Cache(5 * time.Minute).
    Get(ctx)
```

#### `UseIndex(index string) QueryBuilder`
Hints database to use specific index.

```go
qb.Table("users").
   UseIndex("idx_email").
   Where("email", "user@example.com")
```

---

This is Part 1 of the Developer Guide. Would you like me to continue with the remaining sections (Security Guide, Performance, Best Practices, Examples, and Troubleshooting)?