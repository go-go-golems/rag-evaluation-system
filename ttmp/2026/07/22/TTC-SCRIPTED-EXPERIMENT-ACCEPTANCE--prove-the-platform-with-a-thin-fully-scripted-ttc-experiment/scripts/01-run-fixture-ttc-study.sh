#!/usr/bin/env bash
set -euo pipefail
rag_root=$(git rev-parse --show-toplevel)
research_root=${RESEARCHCTL_REPO:-/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl}
ttc_db=${RAG_TTC_DATABASE:-$rag_root/data/ttc-wordpress-rag.sqlite}
[[ -f "$ttc_db" ]] || { echo "TTC database not found: $ttc_db" >&2; exit 2; }
work=$(mktemp -d)
trap '[[ ${KEEP_TTC_SCRIPTED_WORK:-0} == 1 ]] && echo "preserved: $work" >&2 || rm -rf "$work"' EXIT
mkdir -p "$work/bin"
cat >"$work/inputgen.go" <<'GO'
package main
import("context";"database/sql";"encoding/json";"fmt";"os";"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators";_ "github.com/mattn/go-sqlite3")
func main(){ctx:=context.Background();root,source:=os.Args[1],os.Args[2];db,err:=sql.Open("sqlite3",source);if err!=nil{panic(err)};defer db.Close();rows,err:=db.QueryContext(ctx,`SELECT doc_id,title,substr(content_text,1,2400) FROM documents WHERE status='publish' AND length(content_text)>300 ORDER BY doc_id LIMIT 3`);if err!=nil{panic(err)};defer rows.Close();records:=[]ragoperators.SourceRecord{};queries:=[]ragoperators.Query{};for rows.Next(){var id,title,text string;if err:=rows.Scan(&id,&title,&text);err!=nil{panic(err)};records=append(records,ragoperators.SourceRecord{ID:id,Text:title+"\n"+text});queries=append(queries,ragoperators.Query{ID:fmt.Sprintf("query-%02d",len(queries)+1),Text:title,RelevantIDs:[]string{id},Grades:map[string]float64{id:1}})};if len(records)!=3{panic("expected three TTC records")};corpus:=ragoperators.NewCorpusArtifact(ragoperators.Corpus{SchemaVersion:"rag-corpus-data/v1",Records:records},"ttc-scripted");dataset:=ragoperators.NewEvaluationArtifact(ragoperators.EvaluationDataset{SchemaVersion:"rag-evaluation-data/v1",Queries:queries},"ttc-scripted","acceptance","candidate","unit",corpus.Manifest.Digest);write:=func(name string,value any){body,err:=json.Marshal(value);if err!=nil{panic(err)};if err=os.WriteFile(root+"/"+name,body,0644);err!=nil{panic(err)}};write("corpus.json",corpus);write("dataset.json",dataset);write("inputs.json",map[string]any{"inputs":map[string]any{"corpus":map[string]any{"uri":root+"/corpus.json"},"evaluation-dataset":map[string]any{"uri":root+"/dataset.json"}}})}
GO
(cd "$rag_root" && GOWORK=off go run "$work/inputgen.go" "$work" "$ttc_db")
(cd "$rag_root" && GOWORK=off go build -o "$work/bin/rag-eval" ./cmd/rag-eval)
(cd "$rag_root" && GOWORK=off go build -o "$work/bin/rag-workflow-runner" ./cmd/rag-workflow-runner)
(cd "$research_root" && GOWORK=off go build -o "$work/bin/researchctl" ./cmd/researchctl)
"$work/bin/rag-eval" study compile "$rag_root/experiments/ttc-scripted/study.js" --inputs "$work/inputs.json" --artifact-root "$work/artifacts" --output-dir "$work/artifacts/inputs/compiled-study" --experiment-id EXP-TTC-SCRIPTED >"$work/bundle.json"
cp "$rag_root/experiments/ttc-scripted/project.js" "$work/project.js"
"$work/bin/researchctl" lab init --project "$work/project.js" --database "$work/lab.db" >/dev/null
plan="$work/artifacts/inputs/compiled-study/researchctl-plan.js"
"$work/bin/researchctl" experiment validate-plan "$plan" --output json >"$work/validated.json"
run=(experiment run-plan "$plan" --project "$work/project.js" --database "$work/lab.db" --runner-command "$work/bin/rag-workflow-runner" --runner-name scraper-workflow-runner --runner-version v1 --runner-arg=--state-root --runner-arg="$work/workflow-state" --runner-arg=--artifact-root --runner-arg="$work/workflow-artifacts" --runner-arg=--poll-interval --runner-arg=2ms --max-attempts 1 --timeout 60s --output json)
"$work/bin/researchctl" "${run[@]}" >"$work/execution.json"
"$work/bin/researchctl" "${run[@]}" >"$work/resume.json"
analysis=(analysis run "$rag_root/experiments/ttc-scripted/analysis.js" --project "$work/project.js" --database "$work/lab.db" --experiment EXP-TTC-SCRIPTED --output-root "$work/analysis" --output json)
"$work/bin/researchctl" "${analysis[@]}" >"$work/analysis-first.json"
"$work/bin/researchctl" "${analysis[@]}" >"$work/analysis-second.json"
cmp <(jq -S . "$work/analysis-first.json") <(jq -S . "$work/analysis-second.json")
python3 - "$work" <<'PY'
import json,sys
r=sys.argv[1];b=json.load(open(r+'/bundle.json'));v=json.load(open(r+'/validated.json'));e=json.load(open(r+'/execution.json'));resume=json.load(open(r+'/resume.json'));a=json.load(open(r+'/analysis-first.json'))
assert b['schemaVersion']=='rag-workflow-study-bundle/v1' and len(b['cases'])==2
assert v['valid'] and v['caseCount']==2
assert (e['scheduled'],e['executed'],e['failed'])==(6,6,0)
assert (resume['executed'],resume['resumed'])==(0,6)
assert len(a['rows'])==2 and all(x['runs']==3 and x['failedRuns']==0 and x['mrr']['n']==3 for x in a['rows'])
assert all(x['failedOperations']['value']==0 for x in a['rows'])
print('PASS: 2 TTC cases x 3 replicates executed, resumed, and regenerated deterministic quality analysis')
PY
