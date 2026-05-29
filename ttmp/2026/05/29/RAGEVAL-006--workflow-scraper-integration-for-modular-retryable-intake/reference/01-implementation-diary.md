---
Title: Implementation Diary
Ticket: RAGEVAL-006
Status: active
Topics:
    - rag
    - workflow
    - ingestion
    - embeddings
    - corpus
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../scraper/pkg/doc/topics/scraper-queue-policies-and-rate-limiting.md
      Note: Queue policy reference used for LLM and embedding op pacing recommendations
    - Path: ../../../../../../../scraper/pkg/doc/topics/scraper-runtime-model.md
      Note: Scraper runtime reference used to distinguish submit verbs and execution ops
    - Path: cmd/rag-eval/cmds/workflow/root.go
      Note: Phase 2 workflow CLI root
    - Path: cmd/rag-eval/main.go
      Note: Registers workflow CLI command group
    - Path: go.mod
      Note: Phase 0 scraper and go-go-goja dependency compatibility wiring
    - Path: internal/db/queries.go
      Note: Supports document-scoped chunk selection for workflow embeddings
    - Path: internal/db/search_queries.go
      Note: Supports document-scoped chunk selection for workflow BM25
    - Path: internal/services/embedding/service.go
      Note: Propagates document IDs in compute requests
    - Path: internal/services/search/bm25.go
      Note: Uses document IDs when building workflow-scoped BM25 indexes
    - Path: internal/services/search/service.go
      Note: Includes document IDs in BM25 request/result contracts
    - Path: internal/workflow/constants.go
      Note: Phase 0 runner kind and queue constants
    - Path: internal/workflow/echo_runner.go
      Note: Phase 0 custom scraper runner spike
    - Path: internal/workflow/echo_runner_test.go
      Note: Phase 0 scheduler/store integration tests
    - Path: internal/workflow/engine.go
      Note: Phase 2 worker scheduler construction
    - Path: internal/workflow/intake_runner.go
      Note: |-
        Phase 1 Go-native intake runner and chunk_document dispatch
        Completed Phase 1 runner operations for chunking embeddings and BM25
    - Path: internal/workflow/intake_runner_test.go
      Note: |-
        Phase 1 durable chunk_document workflow test
        Completed Phase 1 dependency-chain workflow tests
    - Path: internal/workflow/ops.go
      Note: |-
        Phase 1 typed intake op contracts
        Completed Phase 1 op contracts for chunking embeddings and BM25
    - Path: internal/workflow/submit.go
      Note: Phase 2 workflow submission service
    - Path: internal/workflow/submit_test.go
      Note: Phase 2 fake-provider smoke workflow test
    - Path: ttmp/2026/05/29/RAGEVAL-006--workflow-scraper-integration-for-modular-retryable-intake/design-doc/01-workflow-scraper-intake-integration-design-and-implementation-guide.md
      Note: Primary design guide produced in Step 1
ExternalSources: []
Summary: Chronological diary for the workflow/scraper integration design ticket.
LastUpdated: 2026-05-29T18:30:00Z
WhatFor: Track the design investigation, evidence, decisions, and delivery steps for scraper-backed modular retryable intake.
WhenToUse: Read before implementing workflow orchestration, LLM preprocessing, chunk enrichment, or scraper runner integration in rag-eval.
---







# Diary

## Goal

This diary records the RAGEVAL-006 design workflow: creating the workflow/scraper integration ticket, studying current rag-eval intake services and scraper engine primitives, and writing an intern-ready implementation guide for modular, retryable intake with LLM pre/post-processing.

## Step 1: Ticket Creation, Architecture Mapping, and Design Guide

I created RAGEVAL-006 to capture workflow/scraper integration as a bounded prototype rather than an immediate full migration. The design keeps current rag-eval services as the domain implementation and proposes using scraper as a durable orchestration layer above them.

The key recommendation is to add workflow support now, but add it incrementally. The first implementation should be a Go-native scraper runner that calls existing services, then a small LLM preprocessing/enrichment prototype over one or two TTC documents. The direct CLI/service path should remain the debugging fallback.

### Prompt Context

**User prompt (verbatim):** "create the ticket, and then Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Create a new docmgr ticket for integrating scraper workflows into rag-eval intake, then write and deliver a detailed design/implementation guide suitable for a new intern, including architecture, decisions, pseudocode, APIs, file references, and reMarkable upload.

**Inferred user intent:** The user wants a concrete plan for making rag-eval intake more modular, retryable, and structured, especially to support LLM document preprocessing and chunk postprocessing without losing debuggability.

### What I did

- Created ticket `RAGEVAL-006 -- Workflow/Scraper Integration for Modular Retryable Intake`.
- Added a design document:
  - `ttmp/2026/05/29/RAGEVAL-006--workflow-scraper-integration-for-modular-retryable-intake/design-doc/01-workflow-scraper-intake-integration-design-and-implementation-guide.md`
- Added this implementation diary:
  - `ttmp/2026/05/29/RAGEVAL-006--workflow-scraper-integration-for-modular-retryable-intake/reference/01-implementation-diary.md`
- Added tasks for architecture mapping, guide writing, diary/file relations, validation, and reMarkable upload.
- Inspected current rag-eval intake files:
  - `internal/services/source/service.go`
  - `internal/services/chunking/service.go`
  - `internal/services/embedding/service.go`
  - `internal/services/embedding/provider.go`
  - `internal/services/search/bm25.go`
  - `internal/db/db.go`
  - `internal/api/handlers.go`
  - `cmd/rag-eval/cmds/embedding/compute.go`
- Inspected scraper architecture docs and code:
  - `../scraper/pkg/doc/topics/scraper-architecture-overview.md`
  - `../scraper/pkg/doc/topics/scraper-runtime-model.md`
  - `../scraper/pkg/doc/topics/scraper-queue-policies-and-rate-limiting.md`
  - `../scraper/pkg/doc/topics/scraper-js-api-reference.md`
  - `../scraper/pkg/engine/model/types.go`
  - `../scraper/pkg/engine/store/store.go`
  - `../scraper/pkg/engine/scheduler/scheduler.go`
  - `../scraper/pkg/engine/runner/runner.go`
  - `../scraper/sites/slashdot/verbs/seed.js`
  - `../scraper/sites/slashdot/scripts/seed.js`
- Wrote a detailed design guide covering:
  - current rag-eval intake state;
  - scraper engine primitives;
  - recommended Go-native runner approach;
  - LLM document preprocessing artifact table;
  - chunk enrichment using existing `chunk_enrichments` schema;
  - prototype workflow graph;
  - CLI/API sketches;
  - runner pseudocode;
  - retry/idempotency rules;
  - queue policy recommendations;
  - testing strategy;
  - risks, alternatives, and open questions.

### Why

The current intake path is service-backed and idempotent enough to be orchestrated. LLM preprocessing and chunk enrichment introduce slow, provider-backed, failure-prone work where durable retry and per-operation observability become useful. A workflow system is appropriate now, but only if it wraps existing services rather than replacing them.

### What worked

- The current rag-eval service comments and behavior already anticipate future workflow operations.
- Chunking is rerun-safe for exact document/strategy pairs, which is essential for workflow retries.
- Embedding compute already has source filters, limits, batching, and text-hash freshness checks.
- Scraper exposes a clean runner interface that can support a Go-native rag-eval runner.
- Scraper's queue policy and retry model align well with provider-backed LLM work.

### What didn't work

- No implementation was performed in this step. This is a design and delivery ticket.
- There is not yet a document-level preprocessing artifact table in rag-eval.
- The current `chunk_enrichments` table exists but does not yet have query helpers or service code.
- It is not yet proven that importing scraper as a Go dependency will be frictionless; Phase 0 of the guide explicitly calls for dependency/version validation.

### What I learned

- The strongest integration seam is scraper's Go `Runner` interface, not scraper's JavaScript site scripting layer.
- JavaScript site workflows remain useful as a later option, but starting there would move too much domain behavior out of existing tested Go services.
- The workflow integration should be measured by debuggability: if a workflow op fails, the equivalent direct CLI/service command should still reproduce the failure.

### What was tricky to build

The tricky design point was choosing how much of scraper to reuse. Scraper's normal site model uses manifests, JS submit verbs, and JS execution scripts. Rag-eval already has Go services for source scanning, chunking, embeddings, and BM25 indexing. The design resolves this by using scraper's engine, scheduler, store, retry, and queue machinery while adding a Go-native `rag-eval/intake` runner.

Another tricky point is LLM preprocessing storage. Overwriting `documents.content_text` would make experiments destructive and hard to compare. The guide recommends a new `document_processing_artifacts` table so preprocessing outputs are explicit derived artifacts with prompt/provider/model/hash identity.

### What warrants a second pair of eyes

- Review whether the first implementation should embed scraper as a Go library or call the scraper binary as a subprocess.
- Review whether workflow state should live in a separate engine DB or inside `data/rag-eval.db`.
- Review the proposed `document_processing_artifacts` schema before implementation.
- Review the recommendation to start with a Go-native runner instead of scraper JS site scripts.
- Review dependency version compatibility between rag-eval and scraper before coding.

### What should be done in the future

- Run Phase 0 dependency and minimal custom-runner spike.
- Add document preprocessing artifact persistence.
- Add chunk enrichment service over `chunk_enrichments`.
- Add bounded workflow submit and worker commands.
- Run a two-document TTC guide smoke workflow with fake providers first, then one explicit live provider smoke.

### Code review instructions

- Start with the design guide:
  - `ttmp/2026/05/29/RAGEVAL-006--workflow-scraper-integration-for-modular-retryable-intake/design-doc/01-workflow-scraper-intake-integration-design-and-implementation-guide.md`
- Then review the main rag-eval service boundaries:
  - `internal/services/source/service.go`
  - `internal/services/chunking/service.go`
  - `internal/services/embedding/service.go`
  - `internal/services/search/bm25.go`
- Then review scraper integration seams:
  - `../scraper/pkg/engine/model/types.go`
  - `../scraper/pkg/engine/store/store.go`
  - `../scraper/pkg/engine/scheduler/scheduler.go`
  - `../scraper/pkg/engine/runner/runner.go`

### Technical details

- Ticket: `RAGEVAL-006`.
- Recommended first runner kind: `rag-eval/intake`.
- Recommended first workflow scope: two `ttc-dump-guides` documents, `fixed-1200-150`, first one or two chunks enriched, embedding limit around 10, BM25 smoke index only.
- Recommended new table: `document_processing_artifacts`.
- Recommended existing table to wire: `chunk_enrichments`.

---

## Step 2: Phase 0 Dependency Boundary and Echo Runner Spike

I expanded the ticket tasks into explicit implementation phases before coding, then started Phase 0. The goal was deliberately small: prove that rag-eval can import scraper engine packages, register a custom Go-native runner, create a workflow in a scraper engine SQLite database, and execute one op through the scraper scheduler.

This step does not implement production intake operations yet. It creates the compatibility seam and a tested `internal/workflow` package skeleton so Phase 1 can replace the echo operation with real service-backed operations.

### Prompt Context

**User prompt (verbatim):** "ok, add detailed phases and for each phases tasks, and then work on them step by step, keeping a diary as you work, committing at appropriate intervals."

**Assistant interpretation:** Expand RAGEVAL-006 into detailed phased implementation tasks, then begin executing those phases incrementally with diary updates and focused commits.

**Inferred user intent:** Move from design to implementation while keeping the workflow integration disciplined, inspectable, and easy to review.

### What I did

- Rewrote `tasks.md` with detailed phases 0 through 7 and per-phase task checklists.
- Added a local development dependency on scraper in `go.mod`:
  - `require github.com/go-go-golems/scraper v0.0.0`
  - `replace github.com/go-go-golems/scraper => ../scraper`
- Added a compatibility replacement for `go-go-goja`:
  - `replace github.com/go-go-golems/go-go-goja => github.com/go-go-golems/go-go-goja v0.4.16`
- Created `internal/workflow/constants.go` with:
  - runner kind `rag-eval/intake`
  - queue names `rag-eval:cpu`, `rag-eval:llm`, `rag-eval:embedding`, `rag-eval:index`
- Added `internal/workflow/echo_runner.go`:
  - `EchoRunner` implements scraper's `runner.Runner` interface;
  - operation `echo` returns compact JSON data;
  - unknown operations return non-retryable `model.OpError` with code `unknown_operation`.
- Added `internal/workflow/echo_runner_test.go`:
  - creates a temporary scraper engine SQLite store;
  - registers `EchoRunner` in a scraper runner registry;
  - creates a workflow with one op;
  - runs `scheduler.RunOnce`;
  - asserts the op result exists and the workflow reaches `succeeded`;
  - tests unknown-operation non-retryable error behavior.

### Why

Before implementing real workflow intake ops, we need to know whether scraper can be embedded in the rag-eval Go module at all. A minimal custom runner is the smallest useful proof. It validates the dependency boundary, the scheduler/store/runner path, and the shape of future Go-native operations without touching production intake behavior.

### What worked

- `EchoRunner` successfully implements scraper's `runner.Runner` interface.
- A scraper workflow can be created and run from a rag-eval test.
- The scheduler leases and executes the op, persists the result, and updates workflow status to `succeeded`.
- Focused validation passed:
  - `GOMAXPROCS=2 GOMEMLIMIT=1024MiB go test ./internal/workflow -count=1 -timeout 60s`
- Broader relevant validation passed:
  - `GOMAXPROCS=2 GOMEMLIMIT=1024MiB go test ./internal/db ./internal/ingest ./internal/chunking ./internal/services/source ./internal/services/chunking ./internal/services/document ./internal/services/embedding ./internal/services/search ./internal/workflow -count=1 -timeout 120s`
- Build passed:
  - `GOMAXPROCS=2 GOMEMLIMIT=1024MiB go build ./cmd/rag-eval`

### What didn't work

- The first focused `go test ./internal/workflow` failed after adding the scraper dependency:

```text
# github.com/go-go-golems/scraper/pkg/sites/registry
../scraper/pkg/sites/registry/registry.go:23:38: undefined: gggengine.ModuleSpec
../scraper/pkg/sites/registry/registry.go:31:38: undefined: gggengine.RuntimeModuleRegistrar
# github.com/go-go-golems/scraper/pkg/js/runtime
../scraper/pkg/js/runtime/executor.go:32:38: undefined: gggengine.ModuleSpec
../scraper/pkg/js/runtime/executor.go:33:38: undefined: gggengine.RuntimeModuleRegistrar
../scraper/pkg/js/runtime/executor.go:64:3: gggengine.NewBuilder().WithRequireOptions(require.WithLoader(loader)).WithModules(e.config.Modules...).WithRuntimeModuleRegistrars undefined (type *engine.FactoryBuilder has no field or method WithRuntimeModuleRegistrars)
FAIL	github.com/go-go-golems/rag-evaluation-system/internal/workflow [build failed]
```

- Cause: rag-eval's existing dependency graph selected `github.com/go-go-golems/go-go-goja v0.6.0`, while the local scraper code expects the older module-composition API present in the 0.4.x line.
- Fix for the spike: add a replace to `github.com/go-go-golems/go-go-goja v0.4.16`.
- This is acceptable for Phase 0 but must be reviewed before treating the scraper integration as stable.

### What I learned

- Importing scraper as a Go library is possible, but dependency compatibility is a real implementation concern.
- Scraper's runner package compiles all runner files, including JS runner code, so even a Go-native custom runner test is affected by scraper's go-go-goja API expectations.
- The custom runner seam is still the right integration point: the scheduler/store path works once the module version conflict is resolved.

### What was tricky to build

The tricky part was that the compile failure did not come from the custom runner code. It came from scraper's JS/site packages being compiled as part of the scraper runner package dependency graph. The immediate solution was a module replacement, but the longer-term solution should be a deliberate dependency decision: either keep the replacement, update scraper to the newer go-go-goja API, or isolate the engine/runner packages so Go-native runner usage does not require compiling JS site integration code.

### What warrants a second pair of eyes

- Review the `go-go-goja` replacement before merging widely. It may affect Geppetto or Pinocchio code paths outside the tested package set.
- Review whether scraper should be updated to work with `go-go-goja v0.6.0` instead of pinning rag-eval lower.
- Review whether the scraper engine should expose a smaller package boundary that does not compile JS site integration for Go-native runners.
- Review the `EchoRunner` result/error shape before Phase 1 turns it into a real dispatch runner.

### What should be done in the future

- Complete Phase 0 commit.
- In Phase 1, replace the echo-only runner with typed operation dispatch for real service calls.
- Add a dependency compatibility note to the design guide if the `go-go-goja` replacement remains necessary.
- Consider opening a scraper-side cleanup ticket to decouple engine runner interfaces from JS runner dependencies.

### Code review instructions

- Start with `go.mod` and inspect the two replace directives.
- Review `internal/workflow/constants.go` for runner kind and queue naming.
- Review `internal/workflow/echo_runner.go` for the custom runner shape.
- Review `internal/workflow/echo_runner_test.go` for the temporary engine DB + scheduler integration path.
- Validate with:
  - `GOMAXPROCS=2 GOMEMLIMIT=1024MiB go test ./internal/workflow -count=1 -timeout 60s`
  - `GOMAXPROCS=2 GOMEMLIMIT=1024MiB go test ./internal/db ./internal/ingest ./internal/chunking ./internal/services/source ./internal/services/chunking ./internal/services/document ./internal/services/embedding ./internal/services/search ./internal/workflow -count=1 -timeout 120s`
  - `GOMAXPROCS=2 GOMEMLIMIT=1024MiB go build ./cmd/rag-eval`

### Technical details

- New package: `internal/workflow`.
- Runner kind: `rag-eval/intake`.
- Phase 0 runner operation: `echo`.
- Temporary engine DB: created under `t.TempDir()` in tests.
- Scraper store: `github.com/go-go-golems/scraper/pkg/engine/store/sqlite`.
- Scraper scheduler: `github.com/go-go-golems/scraper/pkg/engine/scheduler`.

---

## Step 3: Phase 1 First Slice — Typed Intake Runner and `chunk_document` Operation

After the Phase 0 echo runner proved that scraper can execute a rag-eval custom runner, I started Phase 1 with the safest real operation: `chunk_document`. This operation uses the existing chunking service, does not call external providers, and is already rerun-safe for a document/strategy pair.

The result is the first real durable intake op. A scraper workflow can now run `rag-eval/intake` with operation `chunk_document`, open the rag-eval app database, call `internal/services/chunking.Service.Apply`, persist chunks, and store a compact op result containing document ID, strategy ID, and chunk count.

### Prompt Context

**User prompt (verbatim):** (same as Step 2)

**Assistant interpretation:** Continue implementing RAGEVAL-006 phase by phase after the detailed task expansion.

**Inferred user intent:** Make steady implementation progress while preserving the ability to review and debug each slice independently.

### What I did

- Added `internal/workflow/ops.go` with first typed intake op input/output structures.
- Added operation constants:
  - `echo`
  - `chunk_document`
- Added `internal/workflow/intake_runner.go`:
  - `IntakeRunner` implements scraper's `runner.Runner` interface;
  - dispatches by `IntakeOpInput.Operation`;
  - validates missing operation and unknown operations;
  - opens/migrates the rag-eval app DB from `db_path`;
  - calls `chunking.Service.Apply` for `chunk_document`;
  - returns compact `ChunkDocumentOutput` JSON.
- Added `internal/workflow/intake_runner_test.go`:
  - seeds a temporary rag-eval app DB with source/document content;
  - creates a temporary scraper engine DB;
  - registers `IntakeRunner`;
  - creates a workflow containing one `chunk_document` op;
  - runs `scheduler.RunOnce`;
  - asserts chunks were written to the app DB and workflow status is `succeeded`;
  - checks missing `db_path` returns a non-retryable op error.
- Updated Phase 1 tasks to mark the first chunk operation slice complete.

### Why

`chunk_document` is the right first real operation because it validates the main architecture without introducing provider complexity. It proves that a workflow op can mutate rag-eval's canonical SQLite state through an existing service and that workflow results remain compact.

### What worked

- Focused workflow tests pass.
- The broader intake/search/workflow test set passes.
- The `rag-eval` binary still builds.
- The existing chunking service needed no production behavior changes.

Validation commands:

```text
GOMAXPROCS=2 GOMEMLIMIT=1024MiB go test ./internal/workflow -count=1 -timeout 60s
GOMAXPROCS=2 GOMEMLIMIT=1024MiB go test ./internal/db ./internal/ingest ./internal/chunking ./internal/services/source ./internal/services/chunking ./internal/services/document ./internal/services/embedding ./internal/services/search ./internal/workflow -count=1 -timeout 120s
GOMAXPROCS=2 GOMEMLIMIT=1024MiB go build ./cmd/rag-eval
```

### What didn't work

- The first `chunk_document` workflow test expected `result.Error == nil`, but scraper's result read path returned a non-nil zero-value `OpError` for a successful result. I adjusted the assertion to treat a nil error or an error with empty code as success for now.
- This deserves a second look in scraper's result storage/read behavior before relying on `result.Error != nil` as the only success test.

### What I learned

- The service-wrapper approach works: workflow ops do not need to duplicate chunking logic.
- Result shape should be intentionally compact. Returning full chunks through workflow results would duplicate canonical data and make engine rows heavy.
- Tests that cross scraper engine DB and rag-eval app DB are straightforward enough to use as the main validation pattern for future ops.

### What was tricky to build

The subtle issue was scraper result error decoding. Successful op completion stores no meaningful error, but the retrieved result may still contain a non-nil zero-value error pointer. For this slice, the test avoids failing on an empty error code. A future scraper-side fix or helper should normalize that behavior so result inspection is less surprising.

### What warrants a second pair of eyes

- Review `internal/workflow/intake_runner.go` error classification. `chunk_document_failed` is currently non-retryable. Some SQLite busy/temporary errors may deserve retryable classification later.
- Review the use of `db_path` in every op input. It is explicit and easy to debug, but repeated DB open/migrate per op may be inefficient.
- Review scraper result error decoding behavior.

### What should be done in the future

- Add remaining Phase 1 operations: `compute_embeddings`, `build_bm25`, `select_documents`, and `summarize_workflow`.
- Add dependency-chain tests, especially chunk → embedding.
- Decide whether to keep `EchoRunner` after `IntakeRunner` fully supports an echo/no-op operation.

### Code review instructions

- Review:
  - `internal/workflow/ops.go`
  - `internal/workflow/intake_runner.go`
  - `internal/workflow/intake_runner_test.go`
- Validate with the same commands listed above.

### Technical details

- New real operation: `chunk_document`.
- Output shape: `{document_id, strategy_id, chunk_count}`.
- DB behavior: per-op open + migrate using `internal/db.OpenDB` and `internal/db.Migrate`.

---

## Step 4: Complete Phase 1 — Embedding and BM25 Workflow Operations

I completed the rest of Phase 1 by extending the Go-native intake runner beyond chunking. The workflow runner can now compute embeddings through the existing embedding service and build BM25 indexes through the existing search service. Both operations keep workflow results compact and write canonical artifacts to the rag-eval database or index directory.

This finishes the first real workflow runner slice. The runner still does not expose operator CLI commands; that is Phase 2. At this point, tests can create scraper workflows programmatically with dependencies such as chunking before embedding or chunking before BM25 indexing.

### Prompt Context

**User prompt (verbatim):** "continue. do the full next phase"

**Assistant interpretation:** Finish the remaining Phase 1 tasks after the initial `chunk_document` slice: add embedding and BM25 workflow operations, dependency-chain tests, diary updates, and a focused commit.

**Inferred user intent:** Get the complete service-backed runner phase done before moving to workflow CLI/user-facing submission commands.

### What I did

- Extended `IntakeOpInput` in `internal/workflow/ops.go` with fields for:
  - strategy IDs and source filters;
  - embedding provider/profile settings;
  - batch size, limit, force;
  - BM25 index root and index ID.
- Added operation constants:
  - `compute_embeddings`
  - `build_bm25`
- Added compact output structs:
  - `ComputeEmbeddingsOutput`
  - `BuildBM25Output`
- Extended `internal/workflow/intake_runner.go`:
  - added `ProviderResolver` seam for tests and real provider resolution;
  - implemented `compute_embeddings` over `embedding.Service.Compute`;
  - implemented `build_bm25` over `search.Service.BuildBM25`;
  - added default provider resolution through `embedding.ResolveProvider`;
  - added `openOpQueries` to normalize per-op DB open/migration errors.
- Extended `internal/workflow/intake_runner_test.go`:
  - added fake embedding provider and resolver;
  - added `chunk_document -> compute_embeddings -> compute_embeddings` dependency-chain test;
  - verifies the second embedding op skips fresh vectors instead of recomputing;
  - added `chunk_document -> build_bm25` dependency-chain test;
  - verifies `search_indexes` metadata is written.
- Updated Phase 1 tasks to complete.

### Why

The workflow integration needs to wrap the existing services that already define the intake pipeline. Completing `compute_embeddings` and `build_bm25` proves that the runner can handle both provider-backed operations and derived filesystem index operations. The dependency-chain tests prove scraper dependencies can order these steps correctly.

### What worked

- `compute_embeddings` works with an injected fake provider, so tests do not call live providers.
- The repeated embedding workflow op correctly reports `computed=0` and skips fresh rows after the first embedding op writes vectors.
- `build_bm25` works against a temporary index root and writes `search_indexes` metadata.
- Validation passed:

```text
GOMAXPROCS=2 GOMEMLIMIT=1024MiB go test ./internal/workflow -count=1 -timeout 60s
GOMAXPROCS=2 GOMEMLIMIT=1024MiB go test ./internal/db ./internal/ingest ./internal/chunking ./internal/services/source ./internal/services/chunking ./internal/services/document ./internal/services/embedding ./internal/services/search ./internal/workflow -count=1 -timeout 120s
GOMAXPROCS=2 GOMEMLIMIT=1024MiB go build ./cmd/rag-eval
```

### What didn't work

- No new blockers in this step.
- The earlier scraper result error decoding caveat still exists in the tests: success checks ignore a non-nil `OpError` when its code is empty.
- Live provider-backed workflow execution is not tested yet. That belongs after the CLI submission path and explicit smoke controls exist.

### What I learned

- The existing service boundaries are strong enough for workflow orchestration.
- Embedding provider injection is essential for workflow tests; otherwise every workflow test would depend on OpenAI/Ollama credentials.
- Search indexing as a workflow op fits the same derived-state model as chunking and embeddings.

### What was tricky to build

The main subtlety was preserving compact workflow results. The embedding service returns enough data to report considered/computed/skipped counts without storing vectors in the scraper engine database. The BM25 service returns enough data to locate the derived index without copying index contents into workflow artifacts. This keeps scraper workflow state as orchestration metadata and leaves canonical artifacts in rag-eval DB/index directories.

### What warrants a second pair of eyes

- `compute_embeddings_failed` is currently marked retryable. This is right for transient provider failures but too broad for configuration errors that reach the service layer. Phase 2 or 3 should add finer error classification.
- `build_bm25_failed` is currently non-retryable. Filesystem/transient errors may deserve retry in future.
- Provider resolution errors are non-retryable. This is correct for missing API keys and bad profiles, but not necessarily for registry IO errors.

### What should be done in the future

- Start Phase 2 by adding operator-facing workflow CLI commands.
- Add workflow status/ops inspection so users do not need test helpers to inspect scraper engine state.
- Consider adding a dry-run or planning mode before live provider workflow runs.

### Code review instructions

- Review:
  - `internal/workflow/ops.go`
  - `internal/workflow/intake_runner.go`
  - `internal/workflow/intake_runner_test.go`
- Pay special attention to provider resolver injection and compact result output.
- Validate with the commands listed above.

### Technical details

- Completed Phase 1 operations:
  - `chunk_document`
  - `compute_embeddings`
  - `build_bm25`
- Dependency-chain tests:
  - chunk → embed → fresh-check embed;
  - chunk → BM25 index.

---

## Step 5: Complete Phase 2 — Operator-Facing Workflow CLI

I added the first operator-facing CLI layer for the scraper-backed intake workflow. The CLI can now submit durable intake workflows, run a local worker cycle or continuous worker, and inspect workflow status and operation state from the scraper engine database.

This phase turns the Phase 1 programmatic workflow runner into something an operator can actually use. The commands intentionally keep direct service commands available: if a workflow op fails, the same chunking, embedding, or search services can still be exercised through existing non-workflow command paths.

### Prompt Context

**User prompt (verbatim):** "phase 2 entirely"

**Assistant interpretation:** Implement all Phase 2 tasks, including workflow CLI registration, submission, worker execution, status/ops inspection, smoke coverage, diary/task updates, and a focused commit.

**Inferred user intent:** Move from internal workflow runner tests to an operator-usable prototype that can create, run, and inspect durable intake workflows.

### What I did

- Added `internal/workflow/submit.go`:
  - `SubmitIntakeRequest`
  - `SubmitIntakeResult`
  - `SubmitIntakeWorkflow`
  - document selection by explicit document IDs or source IDs plus limit;
  - workflow creation with chunk ops per document and optional embedding/BM25 downstream ops.
- Added `internal/workflow/engine.go`:
  - `WorkerConfig`
  - `NewIntakeScheduler`
  - shared scheduler construction for CLI and tests.
- Added `cmd/rag-eval/cmds/workflow` command package:
  - `workflow submit-intake`
  - `workflow run-once`
  - `workflow run-worker`
  - `workflow status`
  - `workflow ops`
- Registered the workflow command group from `cmd/rag-eval/main.go`.
- Added `internal/workflow/submit_test.go` as a smoke/integration test:
  - submits a workflow against temporary rag-eval and scraper engine databases;
  - runs scheduler cycles with a fake embedding provider seam;
  - verifies the workflow reaches `succeeded`.
- Added document-ID filtering through embedding and BM25 service paths so explicitly selected workflow documents do not accidentally trigger downstream work over every chunk with the same strategy.
- Marked Phase 2 tasks complete.

### Why

Phase 1 proved that scraper can orchestrate rag-eval service operations. Phase 2 needed to expose that capability through stable commands so a local operator can submit intake work, run workers, and inspect status without writing test code or manually manipulating the scraper engine database.

### What worked

- `rag-eval workflow --help` shows the new command group and subcommands.
- The workflow submit helper can create a full chunk → embedding/BM25 workflow.
- The worker helper registers the `rag-eval/intake` runner and runs scraper scheduler cycles.
- Status and ops commands reuse scraper's `engineview` read service rather than duplicating read-side SQL.
- Validation passed:

```text
GOMAXPROCS=2 GOMEMLIMIT=1024MiB go test ./internal/workflow -count=1 -timeout 60s
GOMAXPROCS=2 GOMEMLIMIT=1024MiB go test ./internal/db ./internal/ingest ./internal/chunking ./internal/services/source ./internal/services/chunking ./internal/services/document ./internal/services/embedding ./internal/services/search ./internal/workflow ./cmd/rag-eval/cmds/workflow -count=1 -timeout 120s
GOMAXPROCS=2 GOMEMLIMIT=1024MiB go build ./cmd/rag-eval
go run ./cmd/rag-eval workflow --help
```

### What didn't work

- I did not run a live provider workflow from the CLI. The smoke coverage uses a fake provider via the internal scheduler seam so tests stay deterministic and credential-free.
- `workflow status` and `workflow ops` currently emit JSON only. This is useful for debugging and automation but not yet as nice as the Glazed table output used by other commands.

### What I learned

- Scraper's `engineview` package already provides enough read-side functionality for first-pass status and op inspection.
- Keeping scheduler construction in `internal/workflow` avoids duplicating runner registration logic between CLI commands and tests.
- The submission layer is the right place to encode workflow topology: one chunk op per document, followed by aggregate embedding and BM25 ops that depend on all chunk ops.

### What was tricky to build

The tricky part was drawing the boundary between reusable workflow logic and CLI logic. A second sharp edge was downstream scoping: embedding and BM25 services originally filtered only by strategy/source, which was too broad for explicit `--document-ids` workflows, so I added document-ID filters to the shared service/query layer rather than special-casing workflow code.

The broader reusable boundary remains important: The durable workflow topology belongs in `internal/workflow/submit.go`, because tests and future APIs should use the same shape as the CLI. Cobra files should only parse flags, call the service layer, and print JSON. Worker construction follows the same pattern: the CLI should not know the runner registration details beyond the engine DB and worker tuning flags.

### What warrants a second pair of eyes

- The default workflow ID uses timestamp seconds. That is human-readable, but rapid repeated submissions in the same second could collide unless the user provides `--workflow-id`.
- Embedding and BM25 ops currently depend directly on all chunk ops. For very large document sets, this is simple but may create very large dependency lists.
- `submit-intake` uses selected documents ordered by word count when `--document-ids` is omitted. This is useful for bounded smoke runs but should be reviewed for production default behavior.
- CLI output is raw JSON rather than Glazed output.

### What should be done in the future

- Add Glazed/table output for `workflow status` and `workflow ops`.
- Add a `--dry-run` mode to print the workflow topology without creating it.
- Add a safer unique default workflow ID with sub-second precision or a short random suffix.
- Consider source-level workflow fan-out for large corpora.

### Code review instructions

- Start with `internal/workflow/submit.go` to review the workflow topology and document/source scoping.
- Then review `internal/workflow/engine.go` for scheduler/runner registration.
- Review `internal/db/queries.go`, `internal/db/search_queries.go`, `internal/services/embedding/service.go`, and `internal/services/search/bm25.go` for document-ID filter propagation.
- Then review `cmd/rag-eval/cmds/workflow/*.go` for flag mapping and command behavior.
- Validate with the commands listed in `What worked`.

### Technical details

New commands:

```text
rag-eval workflow submit-intake
rag-eval workflow run-once
rag-eval workflow run-worker
rag-eval workflow status
rag-eval workflow ops
```

The submit command maps to existing direct-debug paths:

- `chunk_document` ↔ `rag-eval chunk apply`
- `compute_embeddings` ↔ `rag-eval embedding compute`
- `build_bm25` ↔ `rag-eval search index bm25`
