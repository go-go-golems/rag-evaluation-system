#!/usr/bin/env bash
set -euo pipefail
repo=$(git rev-parse --show-toplevel)
cd "$repo"

fail_if_found() {
  local label=$1 pattern=$2
  shift 2
  if rg -n "$pattern" "$@"; then
    printf 'FAIL: %s\n' "$label" >&2
    exit 1
  fi
}

[[ ! -d internal/workflow ]] || { echo 'FAIL: internal/workflow still exists' >&2; exit 1; }
fail_if_found 'legacy Scraper engine imports remain' 'github\.com/go-go-golems/scraper/pkg/(engine|database|events|metrics|js|steps|webengine)(/|\")' --glob '*.go' --glob '!ttmp/**' .
fail_if_found 'legacy workflow API routes remain' '/api/v1/workflows|/api/v1/workflow-operations|workflow_ops' --glob '!ttmp/**' --glob '!docs/**' --glob '!cmd/rag-eval/doc/**' .
fail_if_found 'legacy engine database flags remain' 'engine-db|rag-eval-workflows' --glob '!ttmp/**' .
fail_if_found 'manual intake operation retry remains' 'retryWorkflowOp|RetryWorkflow|/retry' web/src internal/api cmd/rag-eval/cmds/intake
fail_if_found 'deleted legacy UI panels remain exported' 'Workflow(Op|Summary|List)|QueueHealthPanel' web/src

rg -n '/api/v1/intake/runs' internal/api web/src >/dev/null
rg -n 'rag-intake-v1@1\.0\.0|rag-intake-request/v1' pkg/ragintakeworkflow >/dev/null
printf 'PASS: no legacy RAG intake lifecycle, route, flag, retry, or UI path remains\n'
