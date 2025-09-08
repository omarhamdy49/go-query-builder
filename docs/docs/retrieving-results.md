# Retrieving Results

## `Get`, `First`, `Find`

```go
users, _ := querybuilder.QB().Table("users").Get(ctx)
first,  _ := querybuilder.QB().Table("users").First(ctx)
one,    _ := querybuilder.QB().Table("users").Find(ctx, 42)
```

## Selecting Columns

```go
admins, _ := querybuilder.
  QB().
  Table("users").
  Select("id", "name", "email").
  Where("role", "admin").
  Get(ctx)
```

## Ordering, Limiting, Offsetting

```go
rows, _ := querybuilder.
  QB().
  Table("posts").
  OrderBy("created_at", "desc").
  Limit(10).
  Offset(20).
  Get(ctx)
```

## Result Collections

```go
rows.Each(func(row map[string]any) bool {
  fmt.Println(row["name"])
  return true
})

slice := rows.ToSlice()       // []map[string]any
names := rows.Pluck("name")   // []any
first := rows.First()         // map[string]any or nil
n     := rows.Count()         // int
empty := rows.IsEmpty()       // bool

active := rows.Filter(func(r map[string]any) bool {
  return r["status"] == "active"
})
display := rows.Map(func(r map[string]any) map[string]any {
  r["display"] = fmt.Sprintf("%s <%s>", r["name"], r["email"])
  return r
})
```
