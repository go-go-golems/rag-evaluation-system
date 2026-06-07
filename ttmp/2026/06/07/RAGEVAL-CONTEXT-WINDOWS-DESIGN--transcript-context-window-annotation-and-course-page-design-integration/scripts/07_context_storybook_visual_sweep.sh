#!/usr/bin/env bash
set -euo pipefail

TICKET_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
REPO_ROOT="$(cd "$TICKET_DIR/../../../../.." && pwd)"
OUT_DIR="$TICKET_DIR/sources/visual-parity/context-story-sweep"
STORYBOOK_DIR="$OUT_DIR/storybook-static"
REPORTS_DIR="$OUT_DIR/reports"
PORT="${CSSVD_CONTEXT_STORYBOOK_PORT:-18627}"
BASE_URL="http://127.0.0.1:${PORT}"
INDEX_JSON="$OUT_DIR/index.json"
STORY_LIST="$OUT_DIR/story-list.tsv"
SUMMARY="$OUT_DIR/summary.tsv"

rm -rf "$OUT_DIR"
mkdir -p "$REPORTS_DIR"

cd "$REPO_ROOT"
pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir "$STORYBOOK_DIR"

python3 -m http.server "$PORT" --directory "$STORYBOOK_DIR" > "$OUT_DIR/http-server.log" 2>&1 &
server_pid=$!
trap 'kill "$server_pid" 2>/dev/null || true' EXIT

for _ in $(seq 1 80); do
  if curl -fsS "$BASE_URL/index.json" > "$INDEX_JSON" 2>/dev/null; then
    break
  fi
  sleep 0.1
done

if [[ ! -s "$INDEX_JSON" ]]; then
  echo "failed to fetch Storybook index.json from $BASE_URL" >&2
  exit 1
fi

jq -r '.entries | to_entries[] | select(.value.type == "story") | select(.value.title | test("Context")) | [.key, .value.title, .value.name] | @tsv' "$INDEX_JSON" > "$STORY_LIST"
printf 'story_id\ttitle\tname\tchanged_pixels\tchanged_percent\treport_dir\n' > "$SUMMARY"

while IFS=$'\t' read -r story_id title name; do
  story_url="$BASE_URL/iframe.html?id=$story_id"
  safe_id="${story_id//[^A-Za-z0-9_.-]/_}"
  compare_dir="$REPORTS_DIR/$safe_id"

  css-visual-diff compare \
    --url1 "$story_url" \
    --selector1 '#storybook-root' \
    --url2 "$story_url" \
    --selector2 '#storybook-root' \
    --viewport-w 1440 \
    --viewport-h 1000 \
    --wait-ms1 1500 \
    --wait-ms2 1500 \
    --threshold 30 \
    --out "$compare_dir" >/dev/null

  report_md="$compare_dir/compare.md"
  prefixed_report_md="$compare_dir/01-compare.md"
  if [[ -f "$report_md" ]]; then
    {
      cat <<FM
---
Title: Context Storybook Visual Sweep - ${title} - ${name}
Ticket: RAGEVAL-CONTEXT-WINDOWS-DESIGN
Status: active
Topics: [design-system, frontend-architecture, react, rag, storybook]
DocType: reference
Intent: short-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: css-visual-diff self-compare baseline for context Storybook story ${story_id}.
LastUpdated: 2026-06-07T00:00:00Z
WhatFor: Visual baseline artifact for context-viewer component Storybook coverage.
WhenToUse: During context-viewer implementation and visual regression review.
---

FM
      cat "$report_md"
    } > "$prefixed_report_md"
    rm "$report_md"
  fi

  changed_pixels="$(jq -r '.pixel_diff.changed_pixels // 0' "$compare_dir/compare.json")"
  changed_percent="$(jq -r '.pixel_diff.changed_percent // 0' "$compare_dir/compare.json")"
  printf '%s\t%s\t%s\t%s\t%s\t%s\n' "$story_id" "$title" "$name" "$changed_pixels" "$changed_percent" "$compare_dir" >> "$SUMMARY"
done < "$STORY_LIST"

printf 'Context Storybook visual sweep complete.\nStories: %s\nSummary: %s\nReports: %s\n' "$(($(wc -l < "$STORY_LIST")))" "$SUMMARY" "$REPORTS_DIR"
