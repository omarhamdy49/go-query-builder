## Run locally
cd docs/
python3 -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
mkdocs serve
`Serving on http://127.0.0.1:8000/go-query-builder/`


## Versioned deploy (optional)
./scripts/release.sh v0.1 latest