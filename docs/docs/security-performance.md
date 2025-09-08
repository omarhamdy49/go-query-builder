# Security & Performance

## Security

- Prepared statements & bound parameters
- Input sanitization & validation
- Optional static analysis in CI (e.g., `gosec`)

## Performance Tips

- Prefer indexed `WHERE` columns
- Use `Limit()` and explicit `Select(...)`
- Batch writes with `InsertBatch`
- Tune connection pooling via env vars
