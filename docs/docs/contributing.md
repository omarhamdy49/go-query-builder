# Contributing

1. Fork the repo and create a feature branch.
2. Write tests (`go test ./...`).
3. Ensure docs examples compile where applicable.
4. Open a PR with a clear description and reference to issues.

## Docs

- Local preview: `pip install -r requirements.txt && mkdocs serve`
- Versioned deploy: `./scripts/release.sh v0.1 latest`
