#!/usr/bin/env bash
set -euo pipefail

repo=$(git rev-parse --show-toplevel)
work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
mkdir -p "$work/docs" "$work/state" "$work/indexes"
printf '# Workflow V3 intake smoke\nA bounded document for canonical intake acceptance.\n' >"$work/docs/document.md"

cd "$repo"
GOWORK=off go build -o "$work/rag-eval" ./cmd/rag-eval
bin="$work/rag-eval"
"$bin" source create --db "$work/domain.db" --id smoke --name Smoke --type filesystem --output json >/dev/null
"$bin" source scan --db "$work/domain.db" --source-id smoke --dir "$work/docs" --output json >/dev/null
"$bin" intake submit \
  --db "$work/domain.db" \
  --workflow-db "$work/state/workflow.db" \
  --artifact-root "$work/state/artifacts" \
  --index-root "$work/indexes" \
  --run-id smoke-run \
  --source-ids smoke \
  --skip-preprocessing \
  --skip-chunk-enrichment \
  --skip-embeddings \
  --index-id smoke-index >"$work/submission.json"

for _ in $(seq 1 10); do
  "$bin" intake run-once \
    --db "$work/domain.db" \
    --workflow-db "$work/state/workflow.db" \
    --artifact-root "$work/state/artifacts" \
    --index-root "$work/indexes" >>"$work/dispatch.ndjson"
  status=$("$bin" intake status --workflow-db "$work/state/workflow.db" --artifact-root "$work/state/artifacts" smoke-run | jq -r '.snapshot.status')
  [[ "$status" == succeeded ]] && break
done

"$bin" intake status --workflow-db "$work/state/workflow.db" --artifact-root "$work/state/artifacts" smoke-run >"$work/status.json"
"$bin" intake ops --workflow-db "$work/state/workflow.db" --artifact-root "$work/state/artifacts" --workflow-id smoke-run >"$work/operations.json"

jq -e '.submission.runId == "smoke-run" and .submission.status == "running"' "$work/submission.json" >/dev/null
jq -e '.snapshot.status == "succeeded" and .operations.nodeStatuses.succeeded == 3' "$work/status.json" >/dev/null
jq -e '.run.operations.externalOperations.admitted == 0 and .run.operations.externalOperations.completed == 0' "$work/operations.json" >/dev/null
[[ $(wc -l <"$work/dispatch.ndjson") -eq 3 ]]
find "$work/indexes" -type f -print -quit | grep -q .

printf 'PASS: built-binary intake run completed through Workflow V3 with 3 succeeded nodes and no provider operations\n'
