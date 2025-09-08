# Quick Start Guide

Get up and running with Go Query Builder in just 5 minutes!

## 📦 Installation

```bash
go get github.com/go-query-builder/querybuilder
```

## ⚙️ Environment Setup

Create a `.env` file or set environment variables:

```bash
# Database Configuration  
DB_DRIVER=mysql              # mysql or postgresql
DB_HOST=localhost
DB_PORT=3306                 # 3306 for MySQL, 5432 for PostgreSQL
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database

# Optional: Connection Pool Settings
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_MAX_LIFETIME=5m
```

## 🚀 Your First Query

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
    
    // Zero configuration needed! Environment loaded automatically ✨
    users, err := querybuilder.QB().Table("users").Get(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d users\n", users.Count())
}
```

## 📋 Core Operations

### SELECT Queries

```go
// Basic SELECT
users, err := querybuilder.QB().Table("users").Get(ctx)

// SELECT with WHERE
activeUsers, err := querybuilder.QB().Table("users").
    Where("status", "active").
    Where("age", ">=", 18).
    Get(ctx)

// SELECT with JOIN
userPosts, err := querybuilder.QB().Table("users").
    Select("users.name", "posts.title").
    Join("posts", "users.id", "posts.author_id").
    Where("posts.status", "published").
    Get(ctx)

// SELECT with ORDER and LIMIT
recentUsers, err := querybuilder.QB().Table("users").
    OrderBy("created_at", "desc").
    Limit(10).
    Get(ctx)
```

### INSERT Operations

```go
// Single INSERT
err := querybuilder.QB().Table("users").Insert(ctx, map[string]any{
    "name":       "John Doe",
    "email":      "john@example.com", 
    "age":        30,
    "status":     "active",
    "created_at": time.Now(),
})

// Batch INSERT
users := []map[string]any{
    {"name": "Alice", "email": "alice@test.com", "age": 25},
    {"name": "Bob", "email": "bob@test.com", "age": 30},
}
err := querybuilder.QB().Table("users").InsertBatch(ctx, users)
```

### UPDATE Operations

```go
// UPDATE with WHERE
rowsAffected, err := querybuilder.QB().Table("users").
    Where("id", 1).
    Update(ctx, map[string]any{
        "name":       "Updated Name",
        "updated_at": time.Now(),
    })

// Bulk UPDATE
rowsAffected, err := querybuilder.QB().Table("users").
    Where("status", "inactive").
    Update(ctx, map[string]any{
        "status": "archived",
    })
```

### DELETE Operations

```go
// DELETE with WHERE
rowsAffected, err := querybuilder.QB().Table("users").
    Where("status", "banned").
    Delete(ctx)

// Soft DELETE (update deleted_at)
rowsAffected, err := querybuilder.QB().Table("users").
    Where("id", 1).
    Update(ctx, map[string]any{
        "deleted_at": time.Now(),
    })
```

## 📊 Aggregations

```go
// COUNT
total, err := querybuilder.QB().Table("users").Count(ctx)

// COUNT with WHERE
activeCount, err := querybuilder.QB().Table("users").
    Where("status", "active").
    Count(ctx)

// Other aggregations
avgAge, err := querybuilder.QB().Table("users").Avg(ctx, "age")
maxAge, err := querybuilder.QB().Table("users").Max(ctx, "age") 
totalSalary, err := querybuilder.QB().Table("employees").Sum(ctx, "salary")

// GROUP BY with aggregations
roleStats, err := querybuilder.QB().Table("users").
    Select("role", "COUNT(*) as count", "AVG(age) as avg_age").
    GroupBy("role").
    Having("COUNT(*)", ">", 5).
    Get(ctx)
```

## 📄 Pagination

```go
// Laravel-style pagination
result, err := querybuilder.QB().Table("users").
    Where("status", "active").
    OrderBy("created_at", "desc").
    Paginate(ctx, 1, 15) // page 1, 15 per page

// Access pagination data
fmt.Printf("Page %d of %d\n", result.Meta.CurrentPage, result.Meta.LastPage)
fmt.Printf("Showing %d-%d of %d users\n", 
    result.Meta.From, result.Meta.To, result.Meta.Total)

// Helper methods
if result.HasMorePages() {
    nextPage := *result.GetNextPageNumber()
    fmt.Printf("Next page: %d\n", nextPage)
}

// Iterate through results
result.Data.Each(func(user map[string]any) bool {
    fmt.Printf("User: %s\n", user["name"])
    return true // continue iteration
})
```

## ⚡ Async Operations

```go
// Async queries with goroutines
usersChan := querybuilder.QB().Table("users").GetAsync(ctx)
countChan := querybuilder.QB().Table("posts").CountAsync(ctx)

// Do other work while queries run...
fmt.Println("Queries running in background...")

// Collect results
usersResult := <-usersChan
countResult := <-countChan

if usersResult.Error == nil {
    fmt.Printf("Users loaded: %d\n", usersResult.Data.Count())
}
if countResult.Error == nil {
    fmt.Printf("Post count: %d\n", countResult.Count)
}
```

## 🎪 Using Models

```go
// Define your models
type User struct {
    ID     int    `json:"id" db:"id"`
    Name   string `json:"name" db:"name"`
    Email  string `json:"email" db:"email"`
    Status string `json:"status" db:"status"`
}

func (u User) TableName() string {
    return "users"
}

// Use models with query builder
users, err := querybuilder.QB().Table(User{}).
    Where("status", "active").
    Get(ctx)

// Or use table names directly
users, err := querybuilder.QB().Table("users").
    Where("status", "active").
    Get(ctx)
```

## 🗄️ Multiple Databases

```go
// Add additional database connections
pgConfig := querybuilder.Config{
    Driver:   querybuilder.PostgreSQL,
    Host:     "localhost",
    Port:     5432,
    Database: "analytics_db",
    Username: "postgres",
    Password: "password",
}

querybuilder.QB().AddConnection("analytics", pgConfig)

// Use different connections
mysqlUsers := querybuilder.QB().Table("users").Get(ctx)                    // default
pgAnalytics := querybuilder.Connection("analytics").Table("events").Get(ctx) // postgres
```

## 🔧 Working with Results

```go
users, err := querybuilder.QB().Table("users").Get(ctx)

// Collection methods
userCount := users.Count()
isEmpty := users.IsEmpty()
firstUser := users.First()
userSlice := users.ToSlice()

// Functional operations
activeUsers := users.Filter(func(user map[string]any) bool {
    return user["status"] == "active"
})

names := users.Pluck("name")

// Iteration
users.Each(func(user map[string]any) bool {
    fmt.Printf("User: %s (%s)\n", user["name"], user["email"])
    return true // continue
})
```

## 🎯 Laravel Comparison

### Laravel Eloquent
```php
// Laravel
$users = DB::table('users')
    ->where('status', 'active')
    ->where('age', '>=', 18)
    ->orderBy('created_at', 'desc')
    ->paginate(15);

$count = DB::table('posts')->count();

DB::table('users')->insert([
    'name' => 'John',
    'email' => 'john@example.com'
]);
```

### Go Query Builder
```go
// Go - Same API!
users, err := querybuilder.QB().Table("users").
    Where("status", "active").
    Where("age", ">=", 18).
    OrderBy("created_at", "desc").
    Paginate(ctx, 1, 15)

count, err := querybuilder.QB().Table("posts").Count(ctx)

err := querybuilder.QB().Table("users").Insert(ctx, map[string]any{
    "name":  "John",
    "email": "john@example.com",
})
```

## 🚀 Next Steps

Now that you've got the basics down:

1. **[Query Builder Guide](query-builder.md)** - Dive deeper into query building
2. **[Async Operations](async-operations.md)** - Master concurrent queries  
3. **[Performance Optimization](optimization.md)** - Speed up your queries
4. **[Security Guide](security.md)** - Keep your app secure
5. **[Real-world Examples](examples/real-world.md)** - Production use cases

## 💡 Pro Tips

- **Use context everywhere** - Always pass `context.Context` for timeouts and cancellation
- **Leverage async queries** - Use `GetAsync()` for concurrent operations  
- **Enable query caching** - Set up optimization for better performance
- **Use prepared statements** - Automatically enabled for security
- **Monitor query performance** - Use built-in optimization tools

---

**You're now ready to build amazing Go applications with familiar Laravel syntax!** 🎉

[← Back to Documentation](README.md) | [Configuration Guide →](configuration.md)