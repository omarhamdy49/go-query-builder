# Getting Started

## Installation

```bash
go get github.com/omarhamdy49/go-query-builder
```

## Configuration

The builder reads standard environment variables (or a `.env` if you load it) to configure the default connection.

```env
DB_DRIVER=mysql        # mysql or postgresql
DB_HOST=localhost
DB_PORT=3306           # 3306 for MySQL, 5432 for PostgreSQL
DB_USER=your_username
DB_PASSWORD=your_password
DB_NAME=your_database

# Optional
DB_SSL_MODE=disable
DB_CHARSET=utf8mb4
DB_TIMEZONE=UTC

# Pooling
DB_MAX_OPEN_CONNS=25
DB_MAX_IDLE_CONNS=5
DB_MAX_LIFETIME=5m
DB_MAX_IDLE_TIME=2m
```

## First Query

```go
ctx := context.Background()

users, err := querybuilder.
  QB().
  Table("users").
  Get(ctx)
if err != nil { /* handle */ }

users.Each(func(u map[string]any) bool {
  fmt.Println(u["id"], u["name"])
  return true
})
```
