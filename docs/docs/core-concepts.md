# Core Concepts

## Global Entrypoint (`QB`)

`QB()` returns the singleton query entrypoint — think of it like Laravel’s `DB` facade.

```go
q := querybuilder.QB()
```

## Selecting a Table (`Table`)

```go
query := querybuilder.QB().Table("users")
```

## Using Models (`Table(User{})`)

```go
type User struct {
  ID    int    `db:"id"`
  Name  string `db:"name"`
  Email string `db:"email"`
}

func (User) TableName() string { return "users" }

adults, err := querybuilder.
  QB().
  Table(User{}).
  Where("age", ">=", 18).
  Get(ctx)

one, err := querybuilder.
  QB().
  Table(&User{}).
  Find(ctx, 1)
```

## Convenience Builder (`TableBuilder`)

```go
rows, err := querybuilder.
  TableBuilder("users").
  Where("age", ">", 21).
  Get(ctx)
```

## Multiple Connections

```go
pg := querybuilder.Config{
  Driver:   querybuilder.PostgreSQL,
  Host:     "localhost",
  Port:     5432,
  Database: "analytics_db",
  Username: "postgres",
  Password: "password",
}
querybuilder.QB().AddConnection("analytics", pg)

mysqlUsers, _ := querybuilder.QB().Table("users").Get(ctx)
pgRows, _    := querybuilder.Connection("analytics").Table("events").Get(ctx)
```
