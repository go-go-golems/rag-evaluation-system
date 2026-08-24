# RAG v2 study playbook

A RAG study is pure domain authoring. It describes canonical variants, RAG factors, immutable input bindings, requested measures, and desired replicate counts. It does not schedule processes, contact providers, retry work, or persist scientific runs.

## Ownership

- RAG-eval validates authoring, resolves immutable domain inputs, expands semantic cells, and lowers each execution into Workflow V3.
- Researchctl owns cases, replicates, ordering, concurrency, resume, and laboratory custody.
- Scraper Workflow V3 owns node scheduling, attempts, leases, retries, cancellation, effects, provider operations, and execution artifacts.
- Geppetto remains behind RAG-owned provider adapters.

## Validate and explain

```bash
rag-eval study validate study.js --inputs inputs.json
rag-eval study explain study.js --inputs inputs.json
```

## Compile

```bash
artifact_root="$PWD/laboratory/artifacts"

rag-eval study compile study.js \
  --inputs inputs.json \
  --artifact-root "$artifact_root" \
  --output-dir "$artifact_root/inputs/my-study" \
  --experiment-id EXP-RAG
```

Compilation creates `manifest.json`, `researchctl-plan.js`, and one immutable case directory per expanded RAG cell. The bundle is byte-stable for the same canonical study, inputs, provider authority, and compiler version.

For provider-backed operators, pass reviewed host authority at compile time:

```bash
rag-eval study compile study.js \
  --inputs inputs.json \
  --artifact-root "$artifact_root" \
  --output-dir "$artifact_root/inputs/my-study" \
  --experiment-id EXP-RAG \
  --provider-config provider-config.yaml
```

`--provider-fixture` exists only for deterministic tests. Provider credentials and request content are not embedded in the plan.

## Validate and execute the plan

```bash
researchctl experiment validate-plan \
  "$artifact_root/inputs/my-study/researchctl-plan.js"

researchctl experiment run-plan \
  "$artifact_root/inputs/my-study/researchctl-plan.js" \
  --project project.js \
  --database laboratory.db \
  --runner-command rag-workflow-runner \
  --runner-name scraper-workflow-runner \
  --runner-version v1 \
  --runner-arg=--state-root \
  --runner-arg="$PWD/workflow-state" \
  --runner-arg=--artifact-root \
  --runner-arg="$PWD/workflow-artifacts" \
  --max-attempts 2 \
  --output json
```

Repeat the same command to resume. Do not write a RAG-owned loop around cases or replicates.

## Removed paths

The following are intentionally unavailable:

- `rag-eval study run`;
- `rag-eval preview`;
- `rag-worker`;
- the old RAG preparation Workflow package;
- post-hoc direct-worker progress and run adapters.

For a one-query investigation, compile a one-query evaluation dataset as one Researchctl case. Do not restore a second preview lifecycle.
