#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$repo_root"

sites=(
  examples/xgoja/widget-site
  examples/xgoja/doodle-site
  examples/xgoja/workshop-crm-site
  examples/xgoja-widgetdsl-v3
)

for site in "${sites[@]}"; do
  printf '\n==> Widget DSL v3 site smoke: %s\n' "$site"
  make -C "$site" smoke
done

printf '\nAll Widget DSL v3 generated-host smoke suites passed.\n'
