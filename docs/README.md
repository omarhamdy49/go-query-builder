# Go Query Builder Documentation

Welcome to the comprehensive documentation for Go Query Builder - a Laravel-inspired, production-ready query builder for Go applications.

## 📚 Documentation Index

### Getting Started
- [Installation](installation.md) - Installation and setup guide
- [Configuration](configuration.md) - Environment setup and database configuration  
- [Quick Start](quickstart.md) - Your first queries in 5 minutes

### Core Features
- [Query Builder Basics](query-builder.md) - SELECT, WHERE, JOIN, ORDER BY
- [Database Operations](database-operations.md) - INSERT, UPDATE, DELETE
- [Aggregations](aggregations.md) - COUNT, SUM, AVG, MIN, MAX, GROUP BY
- [Advanced Queries](advanced-queries.md) - Subqueries, JSON, full-text search
- [Pagination](pagination.md) - Laravel-style pagination with data/meta
- [Relationships](relationships.md) - JOIN operations and complex queries

### Performance & Optimization
- [Async Operations](async-operations.md) - Goroutines, channels, concurrent queries
- [Query Optimization](optimization.md) - Caching, prepared statements, performance
- [Connection Pooling](connection-pooling.md) - Database connection management
- [Best Practices](best-practices.md) - Performance tips and production guidelines

### Advanced Topics
- [Security](security.md) - SQL injection prevention, input validation
- [Testing](testing.md) - Unit testing query builders and mocking
- [Multiple Databases](multiple-databases.md) - Multi-database support
- [Custom Extensions](extensions.md) - Extending the query builder
- [Middleware](middleware.md) - Query middleware and hooks

### Migration Guides
- [Laravel to Go](laravel-migration.md) - Migrating from Laravel Eloquent
- [SQL to Query Builder](sql-migration.md) - Converting raw SQL queries
- [Version Upgrades](upgrade-guide.md) - Upgrading between versions

### API Reference
- [QueryBuilder API](api/query-builder.md) - Complete QueryBuilder interface
- [Collection API](api/collection.md) - Working with query results
- [Configuration API](api/configuration.md) - All configuration options
- [Async API](api/async.md) - Asynchronous query methods

### Examples & Tutorials
- [Basic Examples](examples/basic.md) - Simple query examples
- [Complex Queries](examples/complex.md) - Advanced query patterns
- [Real-world Use Cases](examples/real-world.md) - Production examples
- [Performance Examples](examples/performance.md) - Optimization examples

---

## 🚀 Quick Reference

### Basic Usage
```go
import "github.com/go-query-builder/querybuilder"

// Zero configuration - environment loaded automatically
users, err := querybuilder.QB().Table("users").Get(ctx)
```

### Common Operations
```go
// SELECT with WHERE
users := QB().Table("users").Where("status", "active").Get(ctx)

// INSERT
err := QB().Table("users").Insert(ctx, map[string]any{
    "name": "John", "email": "john@example.com"
})

// UPDATE
affected, err := QB().Table("users").Where("id", 1).Update(ctx, map[string]any{
    "name": "Updated Name"
})

// Pagination
result, err := QB().Table("posts").Paginate(ctx, 1, 15)
```

### Async Operations
```go
// Async queries with goroutines
usersChan := QB().Table("users").GetAsync(ctx)
countChan := QB().Table("posts").CountAsync(ctx)

// Process results
usersResult := <-usersChan
countResult := <-countChan
```

---

## 🔗 Laravel Compatibility

This query builder provides **100% API compatibility** with Laravel's Query Builder:

| Laravel Method | Go Query Builder | Status |
|---|---|---|
| `DB::table('users')->get()` | `QB().Table("users").Get(ctx)` | ✅ |
| `->where('status', 'active')` | `.Where("status", "active")` | ✅ |
| `->paginate(15)` | `.Paginate(ctx, 1, 15)` | ✅ |
| `->count()` | `.Count(ctx)` | ✅ |
| `->insert([...])` | `.Insert(ctx, map[string]any{...})` | ✅ |

## 🛠️ Features Overview

- **Zero Configuration** - Automatic environment loading
- **Laravel API Compatibility** - Same methods and patterns
- **Multi-Database Support** - MySQL, PostgreSQL
- **Async Operations** - Goroutines and channels
- **Query Optimization** - Caching, prepared statements
- **Security First** - SQL injection prevention
- **Production Ready** - Connection pooling, error handling
- **Comprehensive Testing** - 100% test coverage

## 📞 Support & Community

- **GitHub Issues** - Bug reports and feature requests
- **Discussions** - Community support and questions
- **Examples** - Real-world usage examples
- **Contributing** - Contribution guidelines

---

**Ready to build lightning-fast Go applications with familiar Laravel syntax!** 🚀