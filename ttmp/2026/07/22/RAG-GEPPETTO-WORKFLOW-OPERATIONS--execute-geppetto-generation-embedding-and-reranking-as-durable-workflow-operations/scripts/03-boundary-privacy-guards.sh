#!/usr/bin/env bash
set -euo pipefail
rag_root=$(git rev-parse --show-toplevel)
scraper_root=${SCRAPER_REPO:-/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/scraper}
research_root=${RESEARCHCTL_REPO:-/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl}
fail_match(){ local message=$1; shift; if "$@"; then echo "$message" >&2; exit 1; fi; }

fail_match "Scraper imports RAG or Geppetto implementation" rg -n 'rag-evaluation-system|/geppetto/' "$scraper_root/pkg"
fail_match "Researchctl imports RAG, Scraper, or Geppetto implementation" rg -n 'rag-evaluation-system|go-go-golems/scraper|/geppetto/' "$research_root/pkg" "$research_root/internal"
fail_match "RAG Workflow imports Scraper legacy workflow" rg -n 'github.com/go-go-golems/scraper/pkg/workflow"' "$rag_root/pkg/ragworkflow" "$rag_root/pkg/ragworkflowops"
fail_match "RAG Workflow imports Researchctl implementation" rg -n 'github.com/go-go-golems/researchctl' "$rag_root/pkg/ragworkflow" "$rag_root/pkg/ragworkflowops"
fail_match "module contains local Scraper replacement" rg -n '^replace github.com/go-go-golems/scraper => (\.|/)' "$rag_root/go.mod"
fail_match "production provider workflow contains unfinished marker" rg -n 'TODO|FIXME|HACK' "$rag_root/pkg/ragworkflow" "$rag_root/pkg/ragworkflowops" "$rag_root/cmd/rag-workflow-runner" "$rag_root/cmd/rag-workflow-provider-fixture"
fail_match "canonical provider fixture contains secret-like material" rg -n 'SECRET-|Authorization:|Bearer |api[_-]?key' "$rag_root/examples/rag-workflow-provider"
fail_match "generic operation contract persists raw provider payload fields" rg -n 'Raw(Request|Response|Prompt|Evidence)|Authorization|Credential|APIKey|Bearer' "$scraper_root/pkg/workflowv3/external_operation.go"

if (cd "$scraper_root" && GOWORK=off go list -deps ./pkg/workflowv3runtime ./pkg/workflowv3sqlite ./pkg/researchrunner) | grep -E 'rag-evaluation-system|go-go-golems/geppetto' >/dev/null; then
  echo "generic Scraper dependency graph includes RAG or Geppetto" >&2; exit 1
fi

tmp=$(mktemp -d); trap 'rm -rf "$tmp"' EXIT
(cd "$rag_root" && GOWORK=off go run ./cmd/rag-workflow-provider-fixture --output "$tmp/fixture" >/dev/null)
diff -rq "$rag_root/examples/rag-workflow-provider" "$tmp/fixture" >/dev/null
printf '%s\n' 'provider boundary/privacy guards passed'
