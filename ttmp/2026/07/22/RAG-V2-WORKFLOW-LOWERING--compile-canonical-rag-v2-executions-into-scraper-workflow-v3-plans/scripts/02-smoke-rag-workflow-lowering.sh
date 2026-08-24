#!/usr/bin/env bash
set -euo pipefail

rag_root=$(git rev-parse --show-toplevel)
research_root=${RESEARCHCTL_REPO:-/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl}
[[ -d "$research_root" ]] || { echo "RESEARCHCTL_REPO is not a directory" >&2; exit 1; }
work=$(mktemp -d)
cleanup() { if [[ "${KEEP_SMOKE_WORK:-0}" == "1" ]]; then echo "preserved smoke work: $work" >&2; else rm -rf "$work"; fi; }
trap cleanup EXIT
mkdir -p "$work/bin"

(cd "$rag_root" && GOWORK=off go build -o "$work/bin/rag-workflow-runner" ./cmd/rag-workflow-runner)
(cd "$rag_root" && GOWORK=off go build -o "$work/bin/rag-workflow-fixture" ./cmd/rag-workflow-fixture)
(cd "$rag_root" && GOWORK=off go build -o "$work/bin/rag-workflow-inspect" ./cmd/rag-workflow-inspect)
(cd "$research_root" && GOWORK=off go build -o "$work/bin/researchctl" ./cmd/researchctl)

"$work/bin/rag-workflow-fixture" --out "$work/generated-fixture" >/dev/null
diff -rq "$rag_root/examples/rag-workflow/provider-free" "$work/generated-fixture"
python3 "$rag_root/ttmp/2026/07/22/RAG-V2-WORKFLOW-LOWERING--compile-canonical-rag-v2-executions-into-scraper-workflow-v3-plans/scripts/01-generate-researchctl-plan.py" \
  --fixture-root "$work/generated-fixture" --out "$work/generated-plan.js"
diff -u "$rag_root/examples/rag-workflow/researchctl-plan.js" "$work/generated-plan.js"

cp "$rag_root/examples/rag-workflow/researchctl-project.js" "$work/project.js"
cp "$rag_root/examples/rag-workflow/researchctl-plan.js" "$work/plan.js"
mkdir -p "$work/artifacts/inputs/rag-workflow"
cp -R "$rag_root/examples/rag-workflow/provider-free/case-a" "$work/artifacts/inputs/rag-workflow/"
cp -R "$rag_root/examples/rag-workflow/provider-free/case-b" "$work/artifacts/inputs/rag-workflow/"
"$work/bin/researchctl" lab init --project "$work/project.js" --database "$work/lab.db" >/dev/null
"$work/bin/researchctl" experiment explain-plan "$work/plan.js" --output json > "$work/explain.json"

cat > "$work/bin/crash-once-runner" <<EOF
#!/usr/bin/env python3
import json, os, subprocess, sys
try: os.mkdir("$work/crash-claimed")
except FileExistsError: os.execv("$work/bin/rag-workflow-runner", ["$work/bin/rag-workflow-runner", *sys.argv[1:]])
request=sys.stdin.buffer.read()
child=subprocess.Popen(["$work/bin/rag-workflow-runner",*sys.argv[1:]],stdin=subprocess.PIPE,stdout=subprocess.PIPE,stderr=subprocess.DEVNULL)
child.stdin.write(request);child.stdin.close()
for line in child.stdout:
    sys.stdout.buffer.write(line);sys.stdout.buffer.flush()
    frame=json.loads(line)
    if frame.get("type")=="event" and frame.get("event",{}).get("type")=="workflow.submitted":
        child.kill();child.wait();sys.exit(17)
sys.exit(child.wait())
EOF
chmod +x "$work/bin/crash-once-runner"

run_args=(experiment run-plan "$work/plan.js" --project "$work/project.js" --database "$work/lab.db" \
  --runner-command "$work/bin/crash-once-runner" --runner-name scraper-workflow-runner --runner-version v1 \
  --runner-arg=--state-root --runner-arg="$work/workflow-state" \
  --runner-arg=--artifact-root --runner-arg="$work/workflow-artifacts" \
  --runner-arg=--poll-interval --runner-arg=2ms --max-attempts 2 --timeout 30s --output json)
"$work/bin/researchctl" "${run_args[@]}" > "$work/result.json"
"$work/bin/researchctl" "${run_args[@]}" > "$work/resume.json"

python3 - "$work/explain.json" "$work/cancel-spec.json" <<'PY'
import json,sys
value=json.load(open(sys.argv[1]))
json.dump(value['schedule'][0]['specification'],open(sys.argv[2],'w'),sort_keys=True,separators=(',',':'))
PY
"$work/bin/researchctl" lab init --project "$work/project.js" --database "$work/cancel-lab.db" >/dev/null
mkdir -p "$work/cancel-artifacts/inputs/rag-workflow"
cp -R "$rag_root/examples/rag-workflow/provider-free/case-a" "$work/cancel-artifacts/inputs/rag-workflow/"
set +e
"$work/bin/researchctl" experiment execute-spec "$work/cancel-spec.json" --project "$work/project.js" \
  --database "$work/cancel-lab.db" --experiment-id EXP-RAG-WORKFLOW \
  --runner-command "$work/bin/rag-workflow-runner" --runner-name scraper-workflow-runner --runner-version v1 \
  --runner-arg=--state-root --runner-arg="$work/cancel-state" \
  --runner-arg=--artifact-root --runner-arg="$work/cancel-workflow-artifacts" \
  --runner-arg=--capacity --runner-arg=blocked.resource=1 --runner-arg=--poll-interval --runner-arg=2ms \
  --runner-cancel-grace 3s --max-attempts 1 --timeout 300ms --output json > "$work/cancel-result.json" 2> "$work/cancel-stderr.txt"
cancel_exit=$?
set -e
[[ $cancel_exit -ne 0 ]] || { echo "timeout execution unexpectedly succeeded" >&2; exit 1; }

mkdir -p "$work/reprojected"
for database in "$work"/workflow-state/*.db "$work"/cancel-state/*.db; do
  read -r run_id status < <(python3 - "$database" <<'PY'
import sqlite3,sys
print(*sqlite3.connect(sys.argv[1]).execute('select run_id,status from v3_runs').fetchone())
PY
)
  if [[ "$status" == "succeeded" || "$status" == "failed" || "$status" == "canceled" ]]; then
    base=$(dirname "$database")
    if [[ "$base" == "$work/workflow-state" ]]; then artifact_base="$work/workflow-artifacts"; else artifact_base="$work/cancel-workflow-artifacts"; fi
    "$work/bin/rag-workflow-inspect" --workflow-db "$database" --artifact-root "$artifact_base/$run_id" --run-id "$run_id" > "$work/reprojected/$run_id.json"
  fi
done

python3 - "$work" "$rag_root" <<'PY'
import glob,json,os,sqlite3,sys
root,rag_root=sys.argv[1:]
result=json.load(open(root+'/result.json'));resume=json.load(open(root+'/resume.json'))
assert (result['scheduled'],result['executed'],result['resumed'],result['failed'])==(4,4,0,0)
assert (resume['scheduled'],resume['executed'],resume['resumed'],resume['failed'])==(4,0,4,0)
attempt_counts=[]
for item in result['items']:
    execution=item['execution'];assert execution['status']=='succeeded';attempt_counts.append(len(execution['attempts']))
    selected=execution['export']['attempts'][-1]
    metrics=[m for m in selected['metrics'] if m['name']=='rag.mrr']
    assert len(metrics)==2 and {m['scope'] for m in metrics}=={'rag.query.q1','rag.query.q2'}
    assert len(selected['metrics'])==24 and len(selected['traces'])==5 and len(selected['artifacts'])==4
    assert all(a['verification']['status']=='verified' for a in selected['artifacts'])
    operations=next(a for a in selected['artifacts'] if a['kind']=='scraper-workflow-external-operations')
    operation_manifest=next(a for a in selected['artifacts'] if a['kind']=='scraper-workflow-external-operation-manifest')
    assert operations['sizeBytes']==0 and operations['metadata']['recordCount']==0
    assert operation_manifest['metadata']['recordCount']==0
    output=next(a for a in selected['artifacts'] if a['kind']=='scraper-workflow-output')
    final=json.load(open(os.path.join(root,'artifacts',output['uri'])))
    case=item['caseId'];parity=json.load(open(os.path.join(rag_root,'examples/rag-workflow/provider-free',case,'ragengine-parity.json')))
    assert final['executionDigest']==parity['executionDigest']
    assert [x['queryId'] for x in final['results']]==[x['queryId'] for x in parity['queries']]
    for got,want in zip(final['results'],parity['queries']):
        assert got['trace']['results']==want['results']
        assert got['metrics']==want['metrics']
    observation_artifact=next(a for a in selected['artifacts'] if a['kind']=='scraper-workflow-observations')
    observation=json.load(open(os.path.join(root,'artifacts',observation_artifact['uri'])))
    reprojected=json.load(open(os.path.join(root,'reprojected',observation['runId']+'.json')))
    assert observation==reprojected
assert sorted(attempt_counts)==[1,1,1,2]
workflow_databases=glob.glob(root+'/workflow-state/*.db');assert len(workflow_databases)==5
cancel=json.load(open(root+'/cancel-result.json'));assert cancel['status']=='failed'
terminal=cancel['export']['attempts'][0]['terminalSummary']['payload'];assert terminal['kind']=='timeout'
cancel_db=glob.glob(root+'/cancel-state/*.db');assert len(cancel_db)==1
conn=sqlite3.connect(cancel_db[0]);run_id,status=conn.execute('select run_id,status from v3_runs').fetchone();assert status=='canceled'
canceled=json.load(open(os.path.join(root,'reprojected',run_id+'.json')));assert canceled['runStatus']=='canceled'
summary={
 'contractFixturesMatched':True,
 'matrix':{'scheduled':4,'executed':4,'researchAttemptCounts':sorted(attempt_counts),'workflowRunsIncludingCrashedAttempt':5},
 'ragParity':{'cases':2,'queriesPerCase':2,'metric':'rag.mrr','matched':True},
 'workflowEvidence':{'metricsPerRun':24,'tracesPerRun':5,'artifactsPerRun':4,'externalOperationRecords':0,'restartReprojectionMatched':True},
 'resume':{'executed':0,'resumed':4},
 'timeout':{'researchStatus':'failed','researchFailureKind':'timeout','workflowStatus':'canceled','durableCanceledObservation':True},
}
print(json.dumps(summary,indent=2,sort_keys=True))
PY
