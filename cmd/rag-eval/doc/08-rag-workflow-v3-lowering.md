---
Title: "RAG v2 Workflow V3 lowering"
Slug: "rag-workflow-v3-lowering"
Short: "Compile provider-free RAG v2 executions into durable Scraper Workflow V3 runs."
Topics:
- rag
- workflow
- experiments
Commands:
- rag-workflow-fixture
- rag-workflow-runner
- rag-workflow-inspect
Flags: []
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

Canonical `rag-pipeline-execution/v2` remains the domain authority. `pkg/ragworkflow` validates that execution and deterministically lowers its closed provider-free operator set into a `scraper-workflow-v3-plan/v1`. Scraper owns node attempts, leases, retries, effects, cancellation, artifacts, and canonical workflow observations. Researchctl may execute the resulting `scraper-workflow-execution/v2` configuration as one immutable laboratory attempt.

This backend does not implement provider calls. Structured representations and real embedding operators fail closed with `RAG_V2_OPERATOR_PROVIDER_REQUIRED`. Add those operations through the separate Geppetto workflow-operations integration rather than weakening the provider-free package.

## Closed package and schemas

The linked package identity is `rag-v2-provider-free@1.0.0`, with immutable task implementation digests and trusted module alias `rag:workflow`.

Important artifacts are:

- `rag-pipeline-execution/v2`: canonical execution identity;
- `rag-corpus/v1`: immutable corpus input;
- `researchctl-set-input-archive/v1`: bounded query archive supplied by Researchctl;
- `scraper-workflow-item-manifest/v1`: staged Workflow map input;
- `rag-workflow-query/v1`: query plus exact dataset-manifest digest;
- `rag-workflow-prepared/v1`: reusable preparation state and fingerprints;
- `rag-workflow-query-result/v1`: one query result;
- `rag-workflow-reduction/v1`: bounded reduction state;
- `rag-workflow-result/v1`: terminal domain result.

Unknown fields, schemas, operators, versions, task identities, package identities, digest mismatches, and stale preparation/index fingerprints are errors. Query maps admit at most 10,000 items and materialize at most eight items ahead. Reductions use fan-in 16. The runner also enforces archive byte/item limits before staging.

## Preparation and query boundaries

Preparation includes corpus loading, unit preparation, chunking, raw representation, optional deterministic fixture embedding, lexical/vector index construction, and index publication. Its fingerprint covers the canonical pipeline, corpus binding, and preparation-affecting operator configuration. Query text and evaluation labels do not enter that fingerprint.

Each staged query item repeats the canonical dataset manifest digest. Query execution rejects a digest that differs from the execution binding. This prevents a valid query archive from being silently used with the wrong experiment cell.

## Observations

Scraper publishes `scraper-workflow-observations/v1`. The RAG projector adds one requested metric per `(metric name, query ID)` scope and one bounded `rag.query` trace per query. Query traces include identities, retrieval/ranking evidence, evaluation values, and usage but omit query text and corpus text. Generic Workflow snapshots do not contain task payloads, artifact locators, provider bodies, arbitrary failure messages, or capabilities.

`workflow-observations.json`, domain `result.json`, and external-operation evidence are verified Researchctl artifacts. Reprojection is deterministic from the terminal Workflow database.

## Commands

Generate deterministic provider-free fixtures:

```text
rag-workflow-fixture --out /tmp/rag-fixture
```

Run as a Researchctl process runner. The program reads one `researchctl-runner-stdio/v1` request from standard input and writes NDJSON frames to standard output:

```text
rag-workflow-runner \
  --state-root /var/lib/rag-workflows \
  --artifact-root /var/lib/rag-artifacts \
  --poll-interval 10ms
```

Reproject observations in a fresh process:

```text
rag-workflow-inspect \
  --workflow-db /var/lib/rag-workflows/<run>.db \
  --artifact-root /var/lib/rag-artifacts/<run> \
  --run-id <run>
```

The runnable two-case/two-replicate plan is `examples/rag-workflow/researchctl-plan.js`. Input URIs are relative to the Researchctl artifact root.

## Validation

Run focused semantic and lifecycle tests:

```text
GOWORK=off go test ./pkg/ragworkflow -count=1
bash ttmp/2026/07/22/RAG-V2-WORKFLOW-LOWERING--*/scripts/02-smoke-rag-workflow-lowering.sh
```

The smoke regenerates fixtures and the JavaScript plan, builds all binaries, compares Workflow results with `ragengine`, injects a runner crash, resumes immutable runs, forces timeout cancellation, and compares fresh observation projection with the exported artifact.

## See also

- `rag-v2-api-reference`
- `rag-study-workflow`
- `scraper-workflow-v3-observations` in Scraper
- `scraper-researchctl-runner` in Scraper
