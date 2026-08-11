#!/usr/bin/env bash
set -euo pipefail

# Print the deterministic Markdown candidate set for the first go-go-parc
# snapshot. The script does not mutate the vault or create an experiment run.
VAULT_ROOT="${VAULT_ROOT:-/home/manuel/code/wesen/go-go-golems/go-go-parc}"

if [[ ! -d "$VAULT_ROOT" ]]; then
  printf 'vault not found: %s\n' "$VAULT_ROOT" >&2
  exit 1
fi

{
  rg --files "$VAULT_ROOT/Projects/2026/07" -g '*.md' 2>/dev/null || true
  rg --files "$VAULT_ROOT/Research/playbooks" -g '*.md' 2>/dev/null || true
  rg --files "$VAULT_ROOT/Research/Institute/Guidelines" -g '*.md' 2>/dev/null || true
} \
  | sed "s#^$VAULT_ROOT/##" \
  | sort -u
