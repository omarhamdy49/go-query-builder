# Project Structure & Testing

## Structure

```
pkg/
  types/        # Interfaces & types
  database/     # Connection management
  query/        # Core builder
  clauses/      # WHERE/JOIN/etc. clause shapes
  execution/    # Statement execution
  pagination/   # Paginator
  security/     # Validations & guards
  config/       # Env/Config loader
examples/       # Usage samples
querybuilder.go # Singleton API surface
```

## Testing

```bash
go test ./...
```

Consider adding linters and security scanners in CI:

```bash
golangci-lint run
gosec ./...
```
