#!/usr/bin/env bash
set -Eeuo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(git -C "$script_dir" rev-parse --show-toplevel)"
example_dir="$repo_root/examples/xgoja/widget-site"
ticket_dir="$repo_root/ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design"
run_dir="$ticket_dir/various/runtime"

printf '[stop] stopping example devctl service\n' >&2
(
  cd "$example_dir"
  devctl down >/dev/null 2>&1 || true
)

if [[ -f "$run_dir/storybook-static.pid" ]]; then
  pid="$(cat "$run_dir/storybook-static.pid")"
  if [[ -n "$pid" ]] && kill -0 "$pid" >/dev/null 2>&1; then
    printf '[stop] stopping storybook static pid=%s\n' "$pid" >&2
    kill "$pid" >/dev/null 2>&1 || true
  fi
  rm -f "$run_dir/storybook-static.pid"
fi
