#!/usr/bin/env bash
set -euo pipefail
export LC_ALL=C
root=$(git rev-parse --show-toplevel)
cd "$root"
[[ ! -d cmd/rag-ttc-v3-sweep ]]
[[ ! -d internal/workflowv3ttc ]]
[[ ! -d internal/workflow ]]
files=$(find experiments/ttc-scripted -maxdepth 1 -type f -printf '%f\n' | sort)
expected=$'README.md\nanalysis.js\ninputs.json\npipeline.js\nproject.js\nstudy.js'
[[ "$files" == "$expected" ]] || { printf 'unexpected TTC source inventory:\n%s\n' "$files" >&2; exit 1; }
if rg -n 'database/sql|sqlite3|os\.Exec|exec\.Command|python|pandas|matplotlib|researchctl/internal|workflowv3sqlite' experiments/ttc-scripted; then
  echo 'TTC workload source contains lifecycle or custom-analysis infrastructure' >&2
  exit 1
fi
rg -n 'researchctl-analysis-spec/v1' experiments/ttc-scripted/analysis.js >/dev/null
rg -n '\.replicates\(3\)' experiments/ttc-scripted/study.js >/dev/null
rg -n 'factors\.enum\("collapse", \["chunk", "unit"\]\)' experiments/ttc-scripted/study.js >/dev/null
printf 'PASS: TTC contains only JS/data/docs workload source and no bespoke lifecycle or analysis infrastructure\n'
