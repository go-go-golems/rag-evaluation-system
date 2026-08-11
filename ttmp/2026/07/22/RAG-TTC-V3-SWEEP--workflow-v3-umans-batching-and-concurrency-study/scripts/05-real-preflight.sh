#!/usr/bin/env bash
# Non-submitting real TTC preflight. This script intentionally requires
# operator-owned paths and never copies or prints provider configuration.
set -euo pipefail

: "${REAL_PROVIDER_CONFIG:?set to host-only provider YAML}"
: "${REAL_SPECIFICATION:?set to current canonical researchctl specification}"
: "${REAL_ARTIFACT_ROOT:?set to the specification artifact root}"

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../../../.." && pwd)
cd "$repo_root"
GOWORK=off go run ./cmd/rag-ttc-v3-sweep \
  --profile real \
  --provider-config "$REAL_PROVIDER_CONFIG" \
  --specification "$REAL_SPECIFICATION" \
  --artifact-root "$REAL_ARTIFACT_ROOT" \
  --chunks 16 --concurrency 1,2 --maximum-requests 60 \
  --prior-generation-requests 61 --maximum-generation-retries 8

# This command returns before any provider submission. Only after it succeeds
# and an operator explicitly approves the printed cumulative limits may a
# separate real invocation include --execute-real and all exact authority flags.
