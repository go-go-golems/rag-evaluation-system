#!/usr/bin/env bash
set -euo pipefail

root=$(git rev-parse --show-toplevel)
cd "$root"
manifest=ttmp/2026/07/22/RAG-V2-EXECUTION-CUTOVER--hard-cut-rag-v2-execution-to-the-canonical-experiment-and-workflow-path/analysis/01-execution-path-deletion-manifest.json

for path in cmd/rag-worker cmd/rag-preparation-smoke cmd/rag-eval/cmds/preview cmd/rag-eval/cmds/workflow internal/preparationworkflow pkg/researchctladapter/run.go pkg/researchctladapter/progress.go; do
  [[ ! -e "$path" ]] || { echo "superseded path remains: $path" >&2; exit 1; }
done

if rg -n 'pkg/ragengine' cmd --glob '*.go' --glob '!**/*_test.go'; then
  echo "production command imports semantic harness directly" >&2
  exit 1
fi
while IFS= read -r importer; do
  case "$importer" in
    pkg/ragproduct/runtime.go|pkg/ragproviders/provider_set.go|pkg/ragworkflow/*.go) ;;
    *) echo "unexpected production ragengine importer: $importer" >&2; exit 1 ;;
  esac
done < <(rg -l 'pkg/ragengine' pkg --glob '*.go' --glob '!**/*_test.go' | sort)
if rg -n 'ExecuteSpecification|researchctl.*execute-spec|rag-worker/v2' cmd pkg internal --glob '*.go' --glob '!**/*_test.go'; then
  echo "direct RAG execution lifecycle remains" >&2
  exit 1
fi
if rg -n '^\s*rag-eval (study run|preview)(\s|$)|--worker-command|cmd/rag-worker' README.md docs cmd/rag-eval/doc experiments web --glob '*.md' --glob '*.tsx'; then
  echo "active documentation teaches a removed command" >&2
  exit 1
fi

help=$(GOWORK=off go run ./cmd/rag-eval study --help)
[[ "$help" == *"compile"* ]] || { echo "study compile missing" >&2; exit 1; }
[[ "$help" != *"**run**"* ]] || { echo "study run remains in help" >&2; exit 1; }
root_help=$(GOWORK=off go run ./cmd/rag-eval --help)
[[ "$root_help" != *"**preview**"* ]] || { echo "preview remains in root help" >&2; exit 1; }
[[ "$root_help" != *"**workflow**"* ]] || { echo "ambiguous legacy workflow command remains" >&2; exit 1; }
[[ "$root_help" == *"**intake**"* ]] || { echo "document intake command missing" >&2; exit 1; }

python3 - "$manifest" <<'PY'
import json,os,sys
value=json.load(open(sys.argv[1]));assert value['schemaVersion']=='rag-v2-execution-cutover-manifest/v1'
entries=value['entries'];paths=[x['path'] for x in entries]
assert len(paths)==len(set(paths)) and entries
assert all(x['classification'] and x['action'] and x['replacement'] and x['proof'] for x in entries)
assert not any(x['classification']=='unresolved' for x in entries)
for item in entries:
    if item['classification']=='superseded-delete': assert not os.path.exists(item['path']),item['path']
print('deletion manifest classifications verified:',len(entries))
PY

if GOWORK=off go list ./... | grep -E 'cmd/rag-worker|cmd/rag-preparation-smoke|internal/preparationworkflow|cmd/rag-eval/cmds/preview'; then
  echo "superseded package remains in go list" >&2
  exit 1
fi

echo "RAG v2 execution cutover guards passed"
