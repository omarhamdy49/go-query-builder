# Filtering (WHERE)

## Basic Operators

```go
rows, _ := querybuilder.
  QB().
  Table("users").
  Where("status", "active").
  Where("age", ">=", 18).
  Get(ctx)
```

## `WhereIn`, `WhereNotNull`, `OrWhere`

```go
rows, _ := querybuilder.
  QB().
  Table("users").
  WhereIn("role", []any{"user", "admin", "moderator"}).
  WhereNotNull("email").
  OrWhere("status", "premium").
  Get(ctx)
```

## Ranges: `between`

```go
rows, _ := querybuilder.
  QB().
  Table("users").
  Where("age", "between", []any{18, 65}).
  Get(ctx)
```

## JSON Queries

```go
rows, _ := querybuilder.
  QB().
  Table("users").
  Where("metadata->theme", "dark").
  WhereJsonContains("preferences", `{"notifications": true}`).
  Get(ctx)
```

## Full-text Search

```go
posts, _ := querybuilder.
  QB().
  Table("posts").
  WhereFullText([]string{"title", "content"}, "golang tutorial").
  Where("status", "published").
  Get(ctx)
```
