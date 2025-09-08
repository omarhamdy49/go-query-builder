# Joins

## `Join`

```go
rows, _ := querybuilder.
  QB().
  Table("users").
  Select("users.name", "posts.title", "posts.created_at").
  Join("posts", "users.id", "posts.author_id").
  Where("posts.status", "published").
  OrderBy("posts.created_at", "desc").
  Get(ctx)
```

## `LeftJoin`

```go
rows, _ := querybuilder.
  QB().
  Table("users").
  Select("users.name", "posts.title", "categories.name AS category").
  LeftJoin("posts", "users.id", "posts.author_id").
  LeftJoin("categories", "posts.category_id", "categories.id").
  Where("users.status", "active").
  Get(ctx)
```
