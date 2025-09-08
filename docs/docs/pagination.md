# Pagination

The paginator returns an object with `Data` (rows) and `Meta` (page info), plus helpers.

```go
result, err := querybuilder.
  QB().
  Table("users").
  Where("status", "active").
  OrderBy("created_at", "desc").
  Paginate(ctx, 1, 15)

fmt.Println("Page:", result.Meta.CurrentPage, "of", result.Meta.LastPage)
if result.HasMorePages() {
  if next := result.GetNextPageNumber(); next != nil {
    fmt.Println("Next page:", *next)
  }
}
```

Typical JSON shape:

```json
{
  "data": [ { "id": 1, "name": "John" }, { "id": 2, "name": "Jane" } ],
  "meta": {
    "current_page": 1, "next_page": 2, "per_page": 15,
    "total": 150, "last_page": 10, "from": 1, "to": 15
  }
}
```
