---
Title: RAG intake Workflow V3 cutover design and implementation guide
Ticket: RAG-INTAKE-WORKFLOW-V3-CUTOVER
Status: active
Topics:
  - rag
  - rag-eval
  - intake
  - scraper
  - workflow
DocType: design-doc
Intent: long-term
Summary: Replace RAG-eval document intake's old Scraper engine with a versioned Workflow V3 task package and product/API surfaces, then delete every old intake caller.
WhatFor: Implement the final RAG-owned migration gate before Scraper legacy-engine deletion.
WhenToUse: Read before changing RAG intake commands, APIs, task execution, or old Scraper engine imports.
---

# RAG intake Workflow V3 cutover

## Goal

Move the active document-intake product from Scraper's old `pkg/engine` and `pkg/workflow` generation to canonical Workflow V3 without preserving an adapter or dual backend. The cutover must retain supported preprocessing, chunking, chunk enrichment, embedding, and BM25 behavior; durable execution and read models become Workflow V3-native.

This is the final RAG-eval deletion gate for `SCRAPER-LEGACY-CLEANUP`.

## Current path

```text
rag-eval intake submit-intake / POST /api/v1/workflows/intake
  -> internal/workflow.SubmitIntakeWorkflow
  -> old engine WorkflowRun + OpSpec rows
  -> old scheduler leases IntakeRunner
  -> RAG services mutate immutable/domain tables
  -> old engineview API + workflow frontend DTOs
```

Operations are:

- `preprocess_document` per selected document;
- `chunk_document` per selected document;
- `enrich_chunk` for selected chunks;
- one `compute_embeddings` operation;
- one `build_bm25` operation.

The old path currently owns eight direct old-Scraper imports across intake commands, API, and `internal/workflow`.

## Target path

```text
rag-eval intake submit
  -> strict rag-intake-request/v1 artifact
  -> deterministic WorkflowPlan V3
  -> workflowv3product.Application.Submit
  -> rag-intake-v1@1.0.0 task package
  -> trusted rag:intake host module
  -> existing RAG domain services
  -> Workflow V3 attempts/leases/retries/cancellation/artifacts/observations

rag-eval intake list/show/cancel/worker
API /api/v1/intake/runs/...
frontend intake views
  -> workflowv3product read models and canonical observations
```

Scraper remains domain-neutral. RAG-eval owns the intake request schema, plan construction, task package, provider resolution, and output schemas.

## Contracts

### Request

`rag-intake-request/v1` is strict and bounded. It contains immutable selection and semantic settings, but no API keys or bearer secrets. Provider configuration is host-owned.

```go
type Request struct {
  SchemaVersion string
  DatabasePath string
  DocumentIDs []string
  SourceIDs []string
  Strategy ChunkStrategy
  Preprocessing PreprocessingPolicy
  Enrichment EnrichmentPolicy
  Embedding EmbeddingPolicy
  BM25 BM25Policy
}
```

The request digest and task-package catalog participate in Workflow identity. Unknown fields, unsupported provider/model/profile identities, empty selections, path escapes, oversized lists, and duplicate IDs fail before submission.

### Task package

Package identity: `rag-intake-v1@1.0.0`.

Task identities:

```text
rag.intake.preprocess-document/v1
rag.intake.chunk-document/v1
rag.intake.enrich-chunk/v1
rag.intake.compute-embeddings/v1
rag.intake.build-bm25/v1
rag.intake.publish/v1
```

Each task accepts the canonical request plus a bounded operation descriptor and emits a typed result artifact. JavaScript task entrypoints are thin calls into a fresh lease-scoped trusted `rag:intake` module. Domain services remain ordinary Go code and do not learn Workflow persistence concepts.

### Plan

Plan construction is deterministic:

```pseudo
validate request
resolve and sort document IDs
nodes = []
for document:
  optional preprocess(document)
  chunk(document), depending on preprocess when enabled
for eligible chunk:
  optional enrich(chunk), depending on its document chunk node
optional embed(all selected documents), depending on all chunk nodes
optional bm25(all selected documents), depending on all chunk nodes
publish(all terminal results)
compile exact Workflow V3 plan
```

Do not use a lazy map for database-mutating document operations unless map item identity and database idempotency are proved. Explicit bounded nodes are preferable for the initial cutover.

### Provider custody

Embedding and generation provider contacts must use the durable external-operation APIs already established by `RAG-GEPPETTO-WORKFLOW-OPERATIONS`. Secrets remain in host configuration. Fake deterministic providers are used first; authorized real providers may be contacted for acceptance.

### API/read models

Replace old `/api/v1/workflows` projections with intake-specific Workflow V3 endpoints:

```text
GET  /api/v1/intake/runs
GET  /api/v1/intake/runs/{runId}
GET  /api/v1/intake/runs/{runId}/observations
POST /api/v1/intake/runs
POST /api/v1/intake/runs/{runId}/cancel
```

There is no per-operation manual retry endpoint. Workflow retry policy is immutable plan identity; an operator creates an explicit new run for rerun semantics.

## Implementation sequence

1. Characterize current outputs, retryability, CLI/API behavior, database writes, and frontend fields.
2. Add `pkg/ragintakeworkflow` with strict codecs, package, bundle, module, and deterministic plan builder.
3. Compose `workflowv3product.Application` in RAG-eval with exact capacities and host provider services.
4. Switch CLI submit/worker/list/show/cancel to V3.
5. Switch API and frontend to intake-specific V3 DTOs.
6. Prove fake-provider parity, database idempotency, restart, lease renewal, retry, cancellation, artifact verification, privacy, and observations.
7. Run authorized provider acceptance where configured.
8. Delete `internal/workflow`, old command/API DTOs, old engine-view handlers, stale docs/tests, and all RAG old-Scraper imports.
9. Add guards rejecting old imports/routes/commands.
10. Close this ticket, then execute `RAG-EVAL-LEGACY-CLEANUP`, `RESEARCHCTL-LEGACY-CLEANUP`, convergence acceptance, and `SCRAPER-LEGACY-CLEANUP` gates in dependency order.

## Deletion gate

The old RAG intake code may be removed only when all are true:

- fake-provider result/database parity passes;
- provider-backed operations use authoritative custody;
- restart and cancellation produce terminal canonical observations;
- CLI and API can submit, inspect, list, cancel, and run workers;
- frontend uses only V3 intake DTOs;
- full RAG-eval tests, race, lint, builds, web checks, and smokes pass;
- repository search finds no production import of Scraper `pkg/engine` or old `pkg/workflow`;
- downstream Scraper and Researchctl acceptance remains green.

## Scraper legacy cleanup dependency

Completing this ticket removes RAG-eval as an external consumer. Scraper can delete its old engine only after its own root commands, worker/API/services, site stack, JavaScript runtime, metrics/events, and frontend consumers also move or receive an explicit delete disposition. Do not claim this RAG cutover alone makes Scraper deletion safe.

## Primary files

- `cmd/rag-eval/cmds/intake/`
- `internal/workflow/`
- `internal/api/workflow_artifact_handlers.go`
- `internal/api/handlers.go`
- `web/src/components/workflows/`
- `pkg/ragintakeworkflow/` (new)
- Scraper `pkg/workflowv3product/`, `pkg/workflowv3runtime/`, `pkg/workflowv3sqlite/`

## Validation

```bash
GOWORK=off go test ./... -count=1
GOWORK=off go test -race ./pkg/ragintakeworkflow ./internal/api ./cmd/rag-eval/cmds/intake
make lint
GOWORK=off go build ./...
pnpm --dir web typecheck
pnpm --dir web build
rg -n 'scraper/pkg/(engine|workflow)(/|")' --glob '*.go' --glob '!ttmp/**'
docmgr doctor --ticket RAG-INTAKE-WORKFLOW-V3-CUTOVER --stale-after 30
```
