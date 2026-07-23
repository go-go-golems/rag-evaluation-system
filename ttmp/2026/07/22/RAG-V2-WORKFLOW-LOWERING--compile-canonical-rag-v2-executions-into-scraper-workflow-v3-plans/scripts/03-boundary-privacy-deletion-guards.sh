#!/usr/bin/env bash
set -euo pipefail
rag_root=$(git rev-parse --show-toplevel)
scraper_root=${SCRAPER_REPO:-/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/scraper}
research_root=${RESEARCHCTL_REPO:-/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl}
fail_match() {
  local message=$1; shift
  if "$@"; then echo "$message" >&2; exit 1; fi
}
[[ ! -e "$rag_root/cmd/rag-ttc-v3-sweep" && ! -e "$rag_root/internal/workflowv3ttc" ]] || { echo "deleted TTC custody path returned" >&2; exit 1; }
fail_match "ragworkflow imports Researchctl implementation" rg -n 'github.com/go-go-golems/researchctl' "$rag_root/pkg/ragworkflow" "$rag_root/cmd/rag-workflow-"*
fail_match "ragworkflow imports provider implementation" rg -n 'geppetto|pinocchio|pkg/ragproviders' "$rag_root/pkg/ragworkflow" "$rag_root/cmd/rag-workflow-"*
fail_match "ragworkflow imports Scraper legacy workflow" rg -n 'github.com/go-go-golems/scraper/pkg/workflow"' "$rag_root/pkg/ragworkflow" "$rag_root/cmd/rag-workflow-"*
fail_match "Scraper generic integration contains RAG/TTC semantics" rg -n -i '\b(rag|ttc|retrieval|embedding)\b' "$scraper_root/pkg/researchrunner" "$scraper_root/pkg/workflowv3product"
fail_match "Researchctl generic integration contains RAG/TTC semantics" rg -n -i '\b(rag|ttc|retrieval|embedding)\b' "$research_root/pkg/lab/processrunner" "$research_root/internal/labsqlite"
fail_match "production lowering contains unfinished marker" rg -n 'TODO|FIXME|HACK' "$rag_root/pkg/ragworkflow" "$rag_root/cmd/rag-workflow-"*
fail_match "module contains local replacement" rg -n '^replace github.com/go-go-golems/scraper => (\.|/)' "$rag_root/go.mod"
fail_match "fixture contains privacy canary" rg -n 'SECRET-|Authorization:|Bearer ' "$rag_root/examples/rag-workflow"

tmp=$(mktemp -d); trap 'rm -rf "$tmp"' EXIT
(cd "$rag_root" && GOWORK=off go run ./cmd/rag-workflow-fixture --out "$tmp/fixture" >/dev/null)
diff -rq "$rag_root/examples/rag-workflow/provider-free" "$tmp/fixture" >/dev/null
python3 "$rag_root/ttmp/2026/07/22/RAG-V2-WORKFLOW-LOWERING--compile-canonical-rag-v2-executions-into-scraper-workflow-v3-plans/scripts/01-generate-researchctl-plan.py" --fixture-root "$tmp/fixture" --out "$tmp/plan.js"
diff -u "$rag_root/examples/rag-workflow/researchctl-plan.js" "$tmp/plan.js" >/dev/null
printf '%s\n' 'boundary/privacy/deletion guards passed'
