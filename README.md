# Go Query Builder 🚀

**Zero Configuration | Laravel-inspired | Production Ready**

A comprehensive, Laravel-inspired query builder for Go with **automatic environment loading** and **singleton pattern**. Just import and start querying - no configuration needed!

## ✨ Features

- 🔥 **Zero Configuration**: Automatically loads from environment variables
- 🎯 **Laravel Query Builder Parity**: 100% compatible API
- 🏗️ **Singleton Pattern**: One import, query anywhere  
- 🔒 **Security First**: Built-in SQL injection prevention
- 🗄️ **Multi-Database**: MySQL + PostgreSQL support
- 🎪 **Smart Table Names**: Auto-detects from models or use strings
- ⚡ **Production Ready**: Connection pooling, transactions, testing

## 📦 Installation

```bash
go get github.com/go-query-builder/querybuilder
```

## ⚡ Quick Start - Zero Config!

### Environment Setup (Optional - uses defaults if not set)

```bash
# .env or environment variables
DB_DRIVER=mysql
DB_HOST=localhost  
DB_PORT=3308
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=your_database
```

### Usage - Just Table Name + Query!

```go
import "github.com/go-query-builder/querybuilder"

// That's it! Config loaded automatically on import ✨

// Just table name → start querying!
users, err := querybuilder.QB().Table("users").Get(ctx)

// With conditions
activeUsers, err := querybuilder.QB().Table("users").
    Where("status", "active").
    Where("age", ">=", 18).
    Get(ctx)
```

## 🎯 Core API Patterns

### 1. Direct Table Access
```go
// String table name - simplest approach
querybuilder.QB().Table("users").Where("status", "active").Get(ctx)
querybuilder.QB().Table("posts").OrderBy("created_at", "desc").Limit(10).Get(ctx)
```

### 2. Model-Based Tables
```go
type User struct {
    ID     int    `json:"id" db:"id"`
    Name   string `json:"name" db:"name"`
    Email  string `json:"email" db:"email"`
    Status string `json:"status" db:"status"`
}

func (u User) TableName() string { return "users" }

// Use model directly
querybuilder.QB().Table(User{}).Where("status", "active").Get(ctx)
querybuilder.QB().Table(&User{}).Find(ctx, 1)
```

### 3. Convenience Functions
```go
// Even shorter syntax
querybuilder.TableBuilder("users").Where("age", ">", 21).Get(ctx)
```

## 🔧 Complete Query Examples

### SELECT Queries
```go
// Basic SELECT
users, err := querybuilder.QB().Table("users").Get(ctx)

// SELECT with conditions
adults, err := querybuilder.QB().Table("users").
    Where("age", ">=", 18).
    Where("status", "active").
    WhereNotNull("email").
    Get(ctx)

// SELECT specific columns
names, err := querybuilder.QB().Table("users").
    Select("name", "email").
    Where("role", "admin").
    Get(ctx)

// Complex WHERE
users, err := querybuilder.QB().Table("users").
    Where("age", "between", []interface{}{18, 65}).
    WhereIn("role", []interface{}{"user", "admin", "moderator"}).
    WhereNotNull("deleted_at").
    OrWhere("status", "premium").
    Get(ctx)
```

### JOIN Queries  
```go
// JOIN with relationships
userPosts, err := querybuilder.QB().Table("users").
    Select("users.name", "posts.title", "posts.created_at").
    Join("posts", "users.id", "posts.author_id").
    Where("posts.status", "published").
    OrderBy("posts.created_at", "desc").
    Get(ctx)

// Multiple JOINs
fullData, err := querybuilder.QB().Table("users").
    Select("users.name", "posts.title", "categories.name as category").
    LeftJoin("posts", "users.id", "posts.author_id").
    LeftJoin("categories", "posts.category_id", "categories.id").
    Where("users.status", "active").
    Get(ctx)
```

### Aggregations
```go
// COUNT
total, err := querybuilder.QB().Table("users").Count(ctx)

// Other aggregates
avgAge, err := querybuilder.QB().Table("users").Avg(ctx, "age")
maxAge, err := querybuilder.QB().Table("users").Max(ctx, "age")
minAge, err := querybuilder.QB().Table("users").Min(ctx, "age")
totalSalary, err := querybuilder.QB().Table("employees").Sum(ctx, "salary")

// GROUP BY with aggregates
stats, err := querybuilder.QB().Table("users").
    Select("role", "COUNT(*) as count", "AVG(age) as avg_age").
    GroupBy("role").
    Having("COUNT(*)", ">", 5).
    OrderBy("count", "desc").
    Get(ctx)
```

### INSERT Operations
```go
// Single insert
user := map[string]interface{}{
    "name":       "John Doe",
    "email":      "john@example.com",
    "age":        30,
    "status":     "active",
    "created_at": time.Now(),
}
err := querybuilder.QB().Table("users").Insert(ctx, user)

// Batch insert
users := []map[string]interface{}{
    {"name": "Alice", "email": "alice@test.com", "age": 25},
    {"name": "Bob", "email": "bob@test.com", "age": 30},
    {"name": "Carol", "email": "carol@test.com", "age": 28},
}
err := querybuilder.QB().Table("users").InsertBatch(ctx, users)
```

### UPDATE Operations
```go
// Update with WHERE
affected, err := querybuilder.QB().Table("users").
    Where("id", 1).
    Update(ctx, map[string]interface{}{
        "name":       "Updated Name",
        "updated_at": time.Now(),
    })

// Bulk update
affected, err := querybuilder.QB().Table("posts").
    Where("status", "draft").
    Where("created_at", "<", time.Now().AddDate(0, -1, 0)).
    Update(ctx, map[string]interface{}{
        "status":     "archived",
        "updated_at": time.Now(),
    })

// Conditional updates
affected, err := querybuilder.QB().Table("users").
    Where("last_login", "<", time.Now().AddDate(0, -6, 0)).
    Update(ctx, map[string]interface{}{
        "status": "inactive",
    })
```

### DELETE Operations
```go
// DELETE with conditions
affected, err := querybuilder.QB().Table("users").
    Where("status", "banned").
    Delete(ctx)

// Soft delete (update deleted_at)
affected, err := querybuilder.QB().Table("posts").
    Where("id", 123).
    Update(ctx, map[string]interface{}{
        "deleted_at": time.Now(),
    })

// Complex DELETE
affected, err := querybuilder.QB().Table("logs").
    Where("created_at", "<", time.Now().AddDate(0, -3, 0)).
    Where("level", "debug").
    Delete(ctx)
```

## 🎪 Advanced Features

### Multiple Database Connections
```go
// Add additional connections
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

### Laravel-Style Pagination
```go
// Laravel-style pagination with data + meta structure
result, err := querybuilder.QB().Table("users").
    Where("status", "active").
    OrderBy("created_at", "desc").
    Paginate(ctx, 1, 15) // page 1, 15 per page

// Access pagination data
fmt.Printf("Page %d of %d\n", result.Meta.CurrentPage, result.Meta.LastPage)
fmt.Printf("Showing %d-%d of %d users\n", result.Meta.From, result.Meta.To, result.Meta.Total)
fmt.Printf("Has more pages: %t\n", result.HasMorePages())

// Iterate through results
result.Data.Each(func(user map[string]interface{}) bool {
    fmt.Printf("User: %s\n", user["name"])
    return true
})

// JSON API response format
{
  "data": [
    {"id": 1, "name": "John", "email": "john@example.com"},
    {"id": 2, "name": "Jane", "email": "jane@example.com"}
  ],
  "meta": {
    "current_page": 1,
    "next_page": 2,
    "per_page": 15,
    "total": 150,
    "last_page": 10,
    "from": 1,
    "to": 15
  }
}

// Pagination helper methods
if result.OnFirstPage() {
    fmt.Println("This is the first page")
}
if result.HasMorePages() {
    nextPage := *result.GetNextPageNumber()
    fmt.Printf("Next page: %d\n", nextPage)
}

// Complex queries with pagination
complexResult, err := querybuilder.QB().Table("posts").
    Select("posts.*", "users.name as author_name").
    Join("users", "posts.author_id", "users.id").
    Where("posts.status", "published").
    Where("posts.created_at", ">=", time.Now().AddDate(0, -1, 0)).
    OrderBy("posts.created_at", "desc").
    Paginate(ctx, 2, 10) // page 2, 10 per page
```

### Advanced WHERE Conditions
```go
// Date helpers
recent, err := querybuilder.QB().Table("posts").
    Where("created_at", ">=", time.Now().AddDate(0, -1, 0)). // Last month
    Where("updated_at", "<=", time.Now()).                    // Up to now
    Get(ctx)

// JSON queries (MySQL 5.7+, PostgreSQL)
jsonUsers, err := querybuilder.QB().Table("users").
    Where("metadata->theme", "dark").
    WhereJsonContains("preferences", `{"notifications": true}`).
    Get(ctx)

// Full-text search
posts, err := querybuilder.QB().Table("posts").
    WhereFullText([]string{"title", "content"}, "golang tutorial").
    Where("status", "published").
    Get(ctx)
```

### Query Building Patterns
```go
// Conditional query building
query := querybuilder.QB().Table("products")

if category != "" {
    query = query.Where("category", category)
}

if minPrice > 0 {
    query = query.Where("price", ">=", minPrice)
}

if maxPrice > 0 {
    query = query.Where("price", "<=", maxPrice)
}

if inStock {
    query = query.Where("stock_count", ">", 0)
}

results, err := query.OrderBy("name").Get(ctx)
```

### Working with Results
```go
users, err := querybuilder.QB().Table("users").Get(ctx)

// Iterate through results
users.Each(func(user map[string]interface{}) bool {
    fmt.Printf("User: %s (%s)\n", user["name"], user["email"])
    return true // continue iteration
})

// Convert to slice
userSlice := users.ToSlice()

// Get specific values
names := users.Pluck("name")
firstUser := users.First()
userCount := users.Count()
isEmpty := users.IsEmpty()

// Functional operations
activeUsers := users.Filter(func(user map[string]interface{}) bool {
    return user["status"] == "active"
})

transformed := users.Map(func(user map[string]interface{}) map[string]interface{} {
    user["display_name"] = fmt.Sprintf("%s (%s)", user["name"], user["role"])
    return user
})
```

## 🎨 Comparison with Laravel

### Laravel Eloquent
```php
// Laravel
$users = User::where('status', 'active')
    ->where('age', '>=', 18)
    ->orderBy('created_at', 'desc')
    ->limit(10)
    ->get();

$userCount = User::count();

User::create([
    'name' => 'John Doe',
    'email' => 'john@example.com'
]);
```

### Go Query Builder (This Package)
```go
// Go Query Builder - Same API!
users, err := querybuilder.QB().Table("users").
    Where("status", "active").
    Where("age", ">=", 18).
    OrderBy("created_at", "desc").
    Limit(10).
    Get(ctx)

userCount, err := querybuilder.QB().Table("users").Count(ctx)

err := querybuilder.QB().Table("users").Insert(ctx, map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
})
```

## 🚀 Environment Configuration

The package automatically loads configuration on import. Create a `.env` file or set environment variables:

```bash
# Database Configuration
DB_DRIVER=mysql              # mysql or postgresql
DB_HOST=localhost
DB_PORT=3306                 # 3306 for MySQL, 5432 for PostgreSQL  
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database

# Optional Settings
DB_SSL_MODE=disable          # PostgreSQL SSL mode
DB_CHARSET=utf8mb4           # MySQL charset
DB_TIMEZONE=UTC              # Database timezone

# Connection Pool Settings
DB_MAX_OPEN_CONNS=25         # Maximum open connections
DB_MAX_IDLE_CONNS=5          # Maximum idle connections  
DB_MAX_LIFETIME=5m           # Connection max lifetime
DB_MAX_IDLE_TIME=2m          # Connection max idle time
```

## 🔒 Security Features

- **SQL Injection Prevention**: All parameters properly bound
- **Input Sanitization**: Automatic sanitization of dangerous patterns
- **Parameter Validation**: Comprehensive input validation
- **Prepared Statements**: All queries use prepared statements
- **Security Testing**: Built-in security checks in CI/CD

## 📊 Performance Tips

1. **Use Indexes**: Ensure WHERE clause columns are indexed
2. **Limit Results**: Always use `Limit()` for large datasets
3. **Connection Pooling**: Configure appropriate pool sizes
4. **Select Specific Columns**: Use `Select()` instead of `SELECT *`
5. **Batch Operations**: Use `InsertBatch()` for multiple inserts

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -race -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Linting
golangci-lint run

# Security scan
gosec ./...
```

## 📁 Project Structure

```
├── pkg/
│   ├── types/           # Interfaces and type definitions
│   ├── database/        # Connection management  
│   ├── query/           # Core query builder
│   ├── clauses/         # Query clause structures
│   ├── execution/       # Query execution logic
│   ├── pagination/      # Pagination utilities
│   ├── security/        # Security validation
│   └── config/          # Configuration management
├── examples/            # Usage examples
├── .github/workflows/   # CI/CD pipelines
└── querybuilder.go      # Main singleton API
```

## 🤝 Contributing

1. Fork the repository
2. Create feature branch (`git checkout -b feature/amazing-feature`)  
3. Commit changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📜 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Laravel's Eloquent Query Builder
- Built with Go best practices and clean architecture
- Community-driven development

---

**⭐ Star this repo if it helped you build awesome Go applications!**