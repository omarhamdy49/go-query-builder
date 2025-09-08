## Run locally
cd docs/
python3 -m venv .venv && source .venv/bin/activate
pip install -r requirements.txt
mkdocs serve

## Versioned deploy (optional)
./scripts/release.sh v0.1 latest