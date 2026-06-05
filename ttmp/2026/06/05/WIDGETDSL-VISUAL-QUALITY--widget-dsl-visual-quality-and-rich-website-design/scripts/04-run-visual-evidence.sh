#!/usr/bin/env bash
set -Eeuo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
repo_root="$(git -C "$script_dir" rev-parse --show-toplevel)"
ticket_dir="$repo_root/ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design"
out_dir="${1:-$ticket_dir/sources/visual-evidence/$(date +%Y%m%d-%H%M%S)}"

"$ticket_dir/scripts/01-start-widget-and-storybook.sh"
trap '"$ticket_dir/scripts/02-stop-widget-and-storybook.sh"' EXIT

mkdir -p "$out_dir"
css-visual-diff verbs \
  --repository "$ticket_dir/scripts" \
  widget-dsl visual capture-evidence \
  "$out_dir" \
  --widgetUrl "${WIDGET_SITE_URL:-http://127.0.0.1:18791}" \
  --storybookUrl "${STORYBOOK_URL:-http://127.0.0.1:6007}" \
  --output json

# css-visual-diff artifact bundles can emit helper Markdown files named index.md
# or computed-css.md. docmgr treats every Markdown file in the ticket as a
# managed document, so convert generated helper Markdown to text unless it is
# the ticket-level summary with explicit frontmatter.
find "$out_dir" -type f -name '*.md' ! -name '01-visual-evidence-summary.md' -print0 | while IFS= read -r -d '' md_file; do
  mv "$md_file" "${md_file%.md}.txt"
done
printf '[evidence] wrote %s\n' "$out_dir" >&2
