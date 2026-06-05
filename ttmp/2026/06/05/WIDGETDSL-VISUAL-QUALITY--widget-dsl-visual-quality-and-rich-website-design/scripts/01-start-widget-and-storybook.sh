#!/usr/bin/env bash
set -Eeuo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(git -C "$script_dir" rev-parse --show-toplevel)"
example_dir="$repo_root/examples/xgoja/widget-site"
storybook_dir="$repo_root/web/storybook-static"
ticket_dir="$repo_root/ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design"
run_dir="$ticket_dir/various/runtime"
mkdir -p "$run_dir"

widget_port="${WIDGET_SITE_PORT:-18791}"
storybook_port="${STORYBOOK_PORT:-6007}"

printf '[start] repo=%s\n' "$repo_root" >&2
printf '[start] widget site via devctl on preferred port %s\n' "$widget_port" >&2
(
  cd "$example_dir"
  devctl up --force --timeout 10m
)

printf '[start] storybook static on port %s\n' "$storybook_port" >&2
if curl -fsS "http://127.0.0.1:${storybook_port}/index.html" >/dev/null 2>&1; then
  printf '[start] existing storybook server is already responding\n' >&2
else
  if [[ ! -f "$storybook_dir/index.html" ]]; then
    printf '[start] missing %s; run (cd web && pnpm build-storybook)\n' "$storybook_dir/index.html" >&2
    exit 1
  fi
  (
    cd "$storybook_dir"
    python3 -m http.server "$storybook_port" --bind 127.0.0.1 >"$run_dir/storybook-static.log" 2>&1 &
    echo $! >"$run_dir/storybook-static.pid"
  )
fi

for i in $(seq 1 100); do
  if curl -fsS "http://127.0.0.1:${storybook_port}/index.html" >/dev/null 2>&1; then
    break
  fi
  sleep 0.1
done

curl -fsS "http://127.0.0.1:${widget_port}/healthz" >/dev/null
curl -fsS "http://127.0.0.1:${storybook_port}/index.html" >/dev/null

cat >"$run_dir/servers.env" <<EOF
WIDGET_SITE_URL=http://127.0.0.1:${widget_port}
STORYBOOK_URL=http://127.0.0.1:${storybook_port}
EOF

printf '[start] widget=%s\n' "http://127.0.0.1:${widget_port}" >&2
printf '[start] storybook=%s\n' "http://127.0.0.1:${storybook_port}" >&2
