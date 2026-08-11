#!/usr/bin/env bash
set -euo pipefail

# Preconditions:
# - Mac llama-server is bound to 127.0.0.1:8012 with --embedding --pooling rank --rerank.
# - tmux session rag-reranker-mimimi forwards local 127.0.0.1:18012 to that port.
# This script is intentionally a bounded contract probe, not a TTC experiment.

endpoint="${RERANK_ENDPOINT:-http://127.0.0.1:18012/v1/rerank}"

curl -fsS --max-time 60 -w '\n__CURL_TOTAL_SECONDS__=%{time_total}\n' "$endpoint" \
  -H 'Content-Type: application/json' \
  -d @- <<'JSON'
{
  "model": "qllama/bge-reranker-v2-m3:q4_k_m",
  "query": "How does TTC calculate a payroll adjustment?",
  "documents": [
    "A payroll adjustment corrects wages, deductions, or time records after a payroll calculation.",
    "TTC offers drought-tolerant cypress trees for privacy planting.",
    "The weather forecast predicts rain during the weekend."
  ],
  "top_n": 3
}
JSON
