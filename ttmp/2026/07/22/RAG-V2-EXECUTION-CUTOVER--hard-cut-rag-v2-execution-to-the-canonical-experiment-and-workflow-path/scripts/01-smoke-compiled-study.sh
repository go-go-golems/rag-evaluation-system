#!/usr/bin/env bash
set -euo pipefail

rag_root=$(git rev-parse --show-toplevel)
research_root=${RESEARCHCTL_REPO:-/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl}
[[ -d "$research_root" ]] || { echo "RESEARCHCTL_REPO is not a directory" >&2; exit 1; }
work=$(mktemp -d)
cleanup() { if [[ "${KEEP_SMOKE_WORK:-0}" == "1" ]]; then echo "preserved smoke work: $work" >&2; else rm -rf "$work"; fi; }
trap cleanup EXIT
mkdir -p "$work/bin"

cat > "$work/inputgen.go" <<'GO'
package main
import (
  "encoding/json"
  "os"
  "github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)
func main() {
  root := os.Args[1]
  corpus := ragoperators.NewCorpusArtifact(ragoperators.Corpus{SchemaVersion:"rag-corpus-data/v1", Records:[]ragoperators.SourceRecord{{ID:"source-1",Text:"Reciprocal rank fusion combines retrieval channels."}}}, "cutover-fixture")
  dataset := ragoperators.NewEvaluationArtifact(ragoperators.EvaluationDataset{SchemaVersion:"rag-evaluation-data/v1",Queries:[]ragoperators.Query{{ID:"q1",Text:"What combines retrieval channels?",RelevantIDs:[]string{"source-1"},Grades:map[string]float64{"source-1":1}}}}, "cutover-fixture", "smoke", "candidate", "unit", corpus.Manifest.Digest)
  write := func(name string, value any) { body,err:=json.Marshal(value);if err!=nil{panic(err)};if err=os.WriteFile(root+"/"+name,body,0644);err!=nil{panic(err)} }
  write("corpus.json",corpus);write("dataset.json",dataset)
  write("inputs.json",map[string]any{"inputs":map[string]any{"corpus":map[string]any{"uri":root+"/corpus.json"},"evaluation-dataset":map[string]any{"uri":root+"/dataset.json"}}})
}
GO

(cd "$rag_root" && GOWORK=off go run "$work/inputgen.go" "$work")
sed 's/\.replicates(1)/.replicates(2)/' "$rag_root/examples/rag-v2/06-raw-study.js" > "$work/study.js"
(cd "$rag_root" && GOWORK=off go build -o "$work/bin/rag-eval" ./cmd/rag-eval)
(cd "$rag_root" && GOWORK=off go build -o "$work/bin/rag-workflow-runner" ./cmd/rag-workflow-runner)
(cd "$research_root" && GOWORK=off go build -o "$work/bin/researchctl" ./cmd/researchctl)

"$work/bin/rag-eval" study compile "$work/study.js" \
  --inputs "$work/inputs.json" \
  --artifact-root "$work/artifacts" \
  --output-dir "$work/artifacts/inputs/compiled-study" \
  --experiment-id EXP-RAG-WORKFLOW > "$work/bundle.json"

cp "$rag_root/examples/rag-workflow/researchctl-project.js" "$work/project.js"
"$work/bin/researchctl" lab init --project "$work/project.js" --database "$work/lab.db" >/dev/null
"$work/bin/researchctl" experiment validate-plan "$work/artifacts/inputs/compiled-study/researchctl-plan.js" --output json > "$work/validated.json"

run_args=(experiment run-plan "$work/artifacts/inputs/compiled-study/researchctl-plan.js" \
  --project "$work/project.js" --database "$work/lab.db" \
  --runner-command "$work/bin/rag-workflow-runner" --runner-name scraper-workflow-runner --runner-version v1 \
  --runner-arg=--state-root --runner-arg="$work/workflow-state" \
  --runner-arg=--artifact-root --runner-arg="$work/workflow-artifacts" \
  --runner-arg=--poll-interval --runner-arg=2ms \
  --max-attempts 1 --timeout 30s --output json)
"$work/bin/researchctl" "${run_args[@]}" > "$work/result.json"
"$work/bin/researchctl" "${run_args[@]}" > "$work/resume.json"

python3 - "$work" <<'PY'
import hashlib,json,os,sys
root=sys.argv[1]
bundle=json.load(open(root+'/bundle.json'));validated=json.load(open(root+'/validated.json'))
result=json.load(open(root+'/result.json'));resume=json.load(open(root+'/resume.json'))
assert bundle['schemaVersion']=='rag-workflow-study-bundle/v1'
assert bundle['experimentId']=='EXP-RAG-WORKFLOW' and len(bundle['cases'])==1
assert validated['valid'] and validated['caseCount']==1
assert (result['scheduled'],result['executed'],result['resumed'],result['failed'])==(2,2,0,0)
assert (resume['scheduled'],resume['executed'],resume['resumed'],resume['failed'])==(2,0,2,0)
for item in result['items']:
  execution=item['execution'];assert execution['status']=='succeeded' and len(execution['attempts'])==1
  selected=execution['export']['attempts'][-1]
  assert any(m['name']=='rag.mrr' and m['scope']=='rag.query.q1' for m in selected['metrics'])
  assert all(a['verification']['status']=='verified' for a in selected['artifacts'])
  operations=next(a for a in selected['artifacts'] if a['kind']=='scraper-workflow-external-operations')
  assert operations['metadata']['recordCount']==0
for case in bundle['cases']:
  for key in ('execution','corpus','queries','domainConfig'):
    value=case[key];body=open(root+'/artifacts/'+value['path'],'rb').read()
    assert len(body)==value['sizeBytes']
    assert 'sha256:'+hashlib.sha256(body).hexdigest()==value['digest']
print(json.dumps({'bundleSchema':bundle['schemaVersion'],'validatedCases':1,'executed':2,'resumed':2,'status':'succeeded','metric':'rag.mrr','externalOperations':0},sort_keys=True))
PY
