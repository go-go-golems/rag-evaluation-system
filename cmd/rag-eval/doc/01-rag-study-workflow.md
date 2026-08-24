---
Title: "Compile RAG v2 studies for Researchctl and Workflow V3"
Slug: "rag-study-workflow"
Short: "Validate, explain, and compile pure RAG studies into immutable Workflow V3 experiment plans."
Topics:
- rag
- studies
- evaluation
- researchctl
- workflow-v3
Commands:
- rag-eval study validate
- rag-eval study explain
- rag-eval study compile
Flags:
- inputs
- artifact-root
- ttc-database
- output-dir
- experiment-id
- provider-config
- provider-fixture
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

RAG-eval owns domain authoring and semantic compilation. Researchctl owns cases, factors, replicates, ordering, resume, and laboratory custody. Scraper Workflow V3 owns production execution, node attempts, retries, leases, cancellation, provider operations, and artifacts.

There is no direct `rag-eval study run`, `rag-eval preview`, or `rag-worker` path. Compilation emits one immutable bundle and a pure Researchctl experiment plan targeting `scraper-workflow-execution/v2`.

## Author and inspect a study

A study script must export pure `rag-study/v2` data through `require("rag")`. It may describe domain variants, factors, requested measures, and desired replicate counts. It must not contact providers or start execution.

```bash
rag-eval study validate study.js --inputs inputs.json
rag-eval study explain study.js --inputs inputs.json
```

## Compile the Workflow bundle

Both the resolved input envelopes and generated Workflow inputs live under the Researchctl artifact root. The output directory must be contained by that root; path escapes are rejected.

```bash
artifact_root="$PWD/laboratory/artifacts"

rag-eval study compile study.js \
  --inputs inputs.json \
  --artifact-root "$artifact_root" \
  --output-dir "$artifact_root/inputs/my-study" \
  --experiment-id EXP-RAG
```

The output directory contains:

```text
manifest.json
researchctl-plan.js
cell-.../execution.json
cell-.../corpus.json
cell-.../queries.json
cell-.../domain-config.json
```

`execution.json` is canonical `rag-pipeline-execution/v2`. `domain-config.json` is canonical `scraper-workflow-execution/v2`. Query inputs use the bounded set-input archive contract. Every file carries a digest and byte count in `manifest.json`.

Provider-required operators must be bound during compilation:

```bash
rag-eval study compile study.js \
  --inputs inputs.json \
  --artifact-root "$artifact_root" \
  --output-dir "$artifact_root/inputs/my-study" \
  --experiment-id EXP-RAG \
  --provider-config provider-config.yaml
```

`--provider-fixture` is deterministic test-only authority. It is mutually exclusive with `--provider-config`. Provider secrets and request bodies are never written into the Workflow input bundle.

## Execute with Researchctl

Initialize the laboratory with a project that declares the same experiment ID, then execute the generated plan through the RAG-owned Workflow runner:

```bash
researchctl lab init --project project.js --database laboratory.db

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

For provider-backed bundles, pass the matching host configuration to the runner with `--runner-arg=--provider-config` and a separate path argument. Researchctl process retries remain scientific-attempt custody; Workflow node and provider-operation retries remain inside each subordinate Workflow run.

Re-running the same `run-plan` command resumes immutable plan items rather than creating a RAG-owned replicate loop.

## Hard-cut diagnostics

| Error | Meaning |
|---|---|
| `RAG_WORKFLOW_PROVIDER_REQUIRED` | The study contains a provider operator but compilation did not bind provider authority. |
| `RAG_WORKFLOW_STUDY_OUTPUT_BOUNDARY` | The output directory escapes the Researchctl artifact root. |
| `RAG_INPUT_DIGEST` | A resolved or staged immutable input no longer matches its declared digest. |
| `RAG_WORKFLOW_DATASET_DIGEST` | The query dataset does not match the canonical execution identity. |

Use `researchctl experiment validate-plan` before execution and `rag-workflow-inspect` for canonical Workflow observation reprojection.
