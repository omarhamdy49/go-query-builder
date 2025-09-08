#!/usr/bin/env bash
set -euo pipefail

VERSION="${1:-}"
if [[ -z "$VERSION" ]]; then
  echo "Usage: scripts/release.sh <version> [alias]" >&2
  exit 1
fi

ALIAS="${2:-latest}"

pip install -r requirements.txt

git fetch origin gh-pages:gh-pages || true
git config user.name "github-actions[bot]"
git config user.email "41898282+github-actions[bot]@users.noreply.github.com"

mike deploy --push --update-aliases "$VERSION" "$ALIAS"
mike set-default --push "$ALIAS"

echo "Deployed version $VERSION and set default alias to $ALIAS"
