# FAQ

**Does it require manual DB initialization?**  
No—config is auto-loaded (env / `.env`), and you can add named connections.

**Can I mix DBs?**  
Yes—use `AddConnection("name", cfg)` + `Connection("name").Table("...")`.

**Does it support JSON and full-text queries?**  
Yes—`WhereJsonContains`, JSON path filters (e.g., `metadata->theme`), and `WhereFullText`.

**How do I paginate like Laravel?**  
Use `Paginate(ctx, page, perPage)` and read `result.Meta` and helpers.
