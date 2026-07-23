#!/usr/bin/env bash
set -euo pipefail

[[ "${RAG_WORKFLOW_REAL_PROVIDER_ACCEPT:-}" == "1" ]] || { echo "set RAG_WORKFLOW_REAL_PROVIDER_ACCEPT=1 to authorize provider contacts" >&2; exit 2; }
: "${RAG_EMBEDDING_BASE_URL:?set RAG_EMBEDDING_BASE_URL}"
: "${RAG_RERANKER_BASE_URL:?set RAG_RERANKER_BASE_URL}"

root=$(git rev-parse --show-toplevel)
source_config=${RAG_PROVIDER_CONFIG:-$root/experiments/real-provider-v2/provider-config.yaml}
ttc_database=${RAG_TTC_DATABASE:-$root/data/ttc-wordpress-rag.sqlite}
[[ -f "$source_config" ]] || { echo "provider config not found" >&2; exit 2; }
[[ -f "$ttc_database" ]] || { echo "TTC database not found" >&2; exit 2; }
curl --fail --silent --show-error --max-time 5 "$RAG_EMBEDDING_BASE_URL/api/tags" >/dev/null
curl --fail --silent --show-error --max-time 5 "$RAG_RERANKER_BASE_URL/health" >/dev/null

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
cp -R "$(dirname "$source_config")/manifests" "$(dirname "$source_config")/schemas" "$work/"
cp "$source_config" "$work/provider-config.yaml"
cd "$root"
RAG_PROVIDER_CONFIG="$work/provider-config.yaml" \
RAG_TTC_DATABASE="$ttc_database" \
GOWORK=off GONOSUMDB=github.com/go-go-golems/scraper \
go test ./pkg/ragworkflow -run TestAuthorizedRealProviderWorkflowAcceptance -count=1 -v
