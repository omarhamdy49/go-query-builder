# Aggregates & Grouping

## Aggregates

```go
total, _  := querybuilder.QB().Table("users").Count(ctx)
avgAge, _ := querybuilder.QB().Table("users").Avg(ctx, "age")
maxAge, _ := querybuilder.QB().Table("users").Max(ctx, "age")
minAge, _ := querybuilder.QB().Table("users").Min(ctx, "age")
payroll,_ := querybuilder.QB().Table("employees").Sum(ctx, "salary")
```

## Grouping & Having

```go
stats, _ := querybuilder.
  QB().
  Table("users").
  Select("role", "COUNT(*) AS count", "AVG(age) AS avg_age").
  GroupBy("role").
  Having("COUNT(*)", ">", 5).
  OrderBy("count", "desc").
  Get(ctx)
```
