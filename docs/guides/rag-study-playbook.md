# RAG study playbook

This playbook takes a pure JavaScript study from authoring to scientifically honest researchctl evidence.

## Author

Start from `experiments/rag-sol2/study.js` or `examples/rag-v2/06-raw-study.js`. JavaScript composes typed values only; it does not open catalogs, call providers, allocate runs, or persist results.

Keep fixed variables explicit in the pipeline. Use variants for coherent pipeline alternatives and factors for controlled substitutions. Check expected cell count before running.

## Bind immutable inputs

`inputs.json` may contain RAG-owned catalog references for author convenience. `rag-eval` resolves these to immutable envelope bytes and manifests before generic submission. Record:

- corpus manifest and record count;
- evaluation manifest, split, status and relevance target;
- exact model/prompt/embedding/reranker manifests;
- operator versions and configs;
- replicate/seed policy;
- requested measure versions/configs.

A candidate, smoke or preview dataset is not an adjudicated benchmark.

## Validate and explain

```bash
rag-eval study validate study.js \
  --inputs inputs.json --ttc-database data/rag-eval.db --output json

rag-eval study explain study.js \
  --inputs inputs.json --ttc-database data/rag-eval.db --output json
```

Validation checks typed graphs, configs, bindings, lineage policies, factors, measures and stable cell identities. Explain output should be reviewed for variants, factor combinations, representation channels, collapse scope, provider requirements and expected cell count.

## Compile without executing

```bash
rag-eval study compile study.js \
  --inputs inputs.json --ttc-database data/rag-eval.db \
  --output-dir compiled-specs --output json
```

Retain canonical specifications when peer review or delayed execution matters. Recompilation from identical authoring and immutable inputs must reproduce identity.

## Preview

```bash
rag-eval preview study.js \
  --inputs inputs.json --ttc-database data/rag-eval.db \
  --query 'What is reciprocal rank fusion?' --variant raw \
  --researchctl-command researchctl --worker-command rag-worker
```

Preview creates a one-query candidate dataset but uses the normal compiler, adapter, worker and researchctl laboratory. It is diagnostic evidence, never a benchmark result.

## Execute

```bash
rag-eval study run study.js \
  --inputs inputs.json --ttc-database data/rag-eval.db \
  --project project.yaml --experiment-id EXP-RAG \
  --researchctl-command researchctl --worker-command rag-worker \
  --spec-output-dir compiled-specs --output json
```

The adapter checks worker capability before allocation. Researchctl verifies generic input/artifact custody and owns run/attempt/retry/timestamp/persistence/export lifecycle. The worker revalidates canonical RAG config, envelope bytes and lineage before executing.

## Inspect and export

Use researchctl's generic commands:

```bash
researchctl lab runs list --project project.yaml --output json
researchctl lab runs show RUN_ID --project project.yaml --output json
researchctl lab export RUN_ID --project project.yaml --output run-export.json
```

Check terminal status, attempt count, requested measures, failure counts, trace kind, artifact digests/sizes and candidate/frozen labels. A clean terminal state alone is insufficient—inspect aggregate metrics and invariants.

## Scientific review

Before any quality claim:

1. verify relevance target maps to evaluated identity;
2. confirm generated representations are not cited evidence;
3. confirm one vote per collapse key/channel before fusion;
4. inspect hydration and source citations;
5. inspect provider/model/prompt/tokenization/truncation/request identities;
6. report failures, abstentions, latency, cost and storage;
7. separate exploratory candidate queries from holdouts;
8. compare only runs with compatible immutable inputs and measurement definitions.

## Reproducibility and custody

Generic file digest and RAG manifest digest are separate and both matter. Input files are read-only; acceptance tests compare staged bytes before/after worker execution. Export reconstruction must reproduce canonical specification bytes.

Do not copy researchctl state into RAG tables, bypass the worker for previews, or introduce a second study runner.
