# API Cheatsheet

## Entrypoints

```go
QB()
QB().Table("users")
QB().Table(User{})
TableBuilder("users")
Connection("analytics").Table("events")
```

## Select / Retrieve

```go
Select("col", "col2")
Get(ctx)
First(ctx)
Find(ctx, id)
OrderBy("col", "desc")
Limit(10)
Offset(20)
```

## Filter

```go
Where("col", "value")
Where("col", ">", 10)
OrWhere("col", "value")
WhereIn("col", []any{...})
WhereNotNull("col")
Where("age", "between", []any{18, 65})
Where("metadata->theme", "dark")
WhereJsonContains("prefs", `{"notifications": true}`)
WhereFullText([]string{"title","body"}, "query")
```

## Join

```go
Join("posts", "users.id", "posts.author_id")
LeftJoin("categories", "posts.category_id", "categories.id")
```

## Aggregation / Grouping

```go
Count(ctx)
Avg(ctx, "col")
Max(ctx, "col")
Min(ctx, "col")
Sum(ctx, "col")
GroupBy("role").Having("COUNT(*)", ">", 5)
```

## Writing

```go
Insert(ctx, map[string]any{...})
InsertBatch(ctx, []map[string]any{...})
Update(ctx, map[string]any{...})
Delete(ctx)
```

## Pagination

```go
p, _ := QB().Table("users").OrderBy("id").Paginate(ctx, 2, 20)
p.Data.Each(func(row map[string]any) bool { return true })
p.Meta
p.OnFirstPage()
p.HasMorePages()
p.GetNextPageNumber()
```
