# Writing Data

## Insert (single)

```go
err := querybuilder.
  QB().
  Table("users").
  Insert(ctx, map[string]any{
    "name":       "John Doe",
    "email":      "john@example.com",
    "age":        30,
    "status":     "active",
    "created_at": time.Now(),
  })
```

## Bulk Insert (`InsertBatch`)

```go
batch := []map[string]any{
  {"name": "Alice", "email": "alice@test.com", "age": 25},
  {"name": "Bob",   "email": "bob@test.com",   "age": 30},
  {"name": "Carol", "email": "carol@test.com", "age": 28},
}
err := querybuilder.QB().Table("users").InsertBatch(ctx, batch)
```

## Update

```go
affected, err := querybuilder.
  QB().
  Table("users").
  Where("id", 1).
  Update(ctx, map[string]any{
    "name":       "Updated Name",
    "updated_at": time.Now(),
  })
```

## Delete

```go
affected, err := querybuilder.
  QB().
  Table("users").
  Where("status", "banned").
  Delete(ctx)
```
