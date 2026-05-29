# Tasks

## Completed design/delivery setup

- [x] Create RAGEVAL-006 ticket workspace.
- [x] Map current rag-eval intake services and scraper workflow architecture.
- [x] Write intern-ready workflow/scraper integration design guide.
- [x] Record investigation diary.
- [x] Relate key rag-eval and scraper files to the design document.
- [x] Validate ticket docs with `docmgr doctor`.
- [x] Upload design package to reMarkable.

## Phase 0 — Dependency boundary and custom runner spike

Goal: prove rag-eval can import scraper engine packages and run a tiny custom runner without touching production intake behavior.

- [x] Add a local development dependency on `github.com/go-go-golems/scraper` using a temporary `replace ../scraper` entry if needed.
- [x] Create `internal/workflow` package skeleton with constants for runner kind and queue names.
- [x] Implement a no-op or echo custom runner that satisfies scraper's `runner.Runner` interface.
- [x] Add a temporary-SQLite integration test that creates a scraper engine store, registers the custom runner, creates a workflow with one op, runs `scheduler.RunOnce`, and asserts workflow success.
- [x] Add a second test for retry classification or non-retryable unknown operations if the runner dispatch layer exists.
- [x] Validate with focused `go test ./internal/workflow -count=1` and then broader `go test ./internal/...` if dependency resolution is stable.
- [x] Commit Phase 0 as a focused dependency/spike commit.

## Phase 1 — Go-native `rag-eval/intake` runner over existing services

Goal: turn existing idempotent services into durable workflow operations while preserving direct CLI/service debuggability.

- [ ] Define typed op input/output structs for `scan_source`, `select_documents`, `chunk_document`, `compute_embeddings`, `build_bm25`, and `summarize_workflow`.
- [ ] Implement operation dispatch with strict validation and compact `OpResult.Data` JSON.
- [ ] Implement `chunk_document` by calling `internal/services/chunking.Service.Apply`.
- [ ] Implement `compute_embeddings` using a fake/test provider seam first, then the existing provider resolver for real profiles.
- [ ] Implement `build_bm25` by calling `internal/services/search.Service.BuildBM25`.
- [ ] Add temporary rag-eval DB + scraper engine DB integration tests for chunk → embedding dependency flow.
- [ ] Add tests proving retries do not duplicate chunks or fresh embeddings.
- [ ] Commit Phase 1 as the first real workflow runner slice.

## Phase 2 — Workflow submission and worker CLI

Goal: give operators a rag-eval-native way to submit and run bounded intake workflows.

- [ ] Add `cmd/rag-eval/cmds/workflow/root.go` and register it from `cmd/rag-eval/main.go`.
- [ ] Add `rag-eval workflow submit-intake` with explicit source/document/chunk/embedding/BM25 limits.
- [ ] Add `rag-eval workflow run-worker` or `run-once` for local worker cycles over the scraper engine DB.
- [ ] Add `rag-eval workflow status` and `rag-eval workflow ops` for basic inspection.
- [ ] Ensure every workflow op can be reproduced by an equivalent existing direct command or documented service call.
- [ ] Smoke-test a workflow over one or two existing TTC guide documents with fake providers.
- [ ] Commit Phase 2 as the operator-facing workflow prototype.

## Phase 3 — Document preprocessing artifacts

Goal: add non-destructive document-level LLM preprocessing as derived state.

- [ ] Add `document_processing_artifacts` schema with document, artifact type, prompt version, provider, model, input hash, output text/JSON, status, and error metadata.
- [ ] Add DB query helpers for upsert, lookup, freshness, and coverage counts.
- [ ] Add `internal/services/documentprocessing` with fake-provider tests.
- [ ] Add workflow op `preprocess_document` with prompt version and provider identity.
- [ ] Add a direct CLI or test harness path so preprocessing can be debugged outside workflow.
- [ ] Validate that preprocessing never overwrites `documents.content_text` in the first implementation.
- [ ] Commit Phase 3 as the document preprocessing artifact slice.

## Phase 4 — Chunk enrichment over `chunk_enrichments`

Goal: add retryable LLM postprocessing for chunk summaries, topics, entities, and hypothetical questions.

- [ ] Add DB helpers for `chunk_enrichments` lookup/upsert/freshness/coverage.
- [ ] Add `internal/services/chunkenrichment` with strict output validation and fake-provider tests.
- [ ] Add workflow op `enrich_chunk` that skips fresh enrichments by `text_hash`.
- [ ] Add bounded fan-out from selected documents to first N chunks per document.
- [ ] Add direct debugging command or test harness for one chunk enrichment.
- [ ] Commit Phase 4 as the chunk enrichment slice.

## Phase 5 — Bounded live-provider workflow smoke

Goal: prove the workflow can run real provider-backed preprocessing/enrichment safely.

- [ ] Add explicit live smoke documentation and cost controls.
- [ ] Run one document preprocessing or one chunk enrichment with a live profile only after fake-provider tests pass.
- [ ] Run a two-document workflow with strict document/chunk/embedding limits.
- [ ] Record exact commands, provider/model, prompt versions, counts, and failures in the diary.
- [ ] Commit any fixes from live smoke separately from generated data.

## Phase 6 — Workflow visibility in API/Corpus Explorer

Goal: make workflow-produced artifacts visible enough for users to understand intake state.

- [ ] Add read-only API endpoints for workflow summaries or adapt scraper engineview services.
- [ ] Add document preprocessing artifact coverage by source and prompt version.
- [ ] Add chunk enrichment coverage by strategy and prompt version.
- [ ] Link Corpus Explorer document/chunk views to preprocessing/enrichment status.
- [ ] Add UI smoke tests or manual browser validation notes.
- [ ] Commit Phase 6 as the first workflow visibility slice.

## Phase 7 — Retrospective and expansion decision

Goal: decide whether to expand scraper workflows beyond the prototype.

- [ ] Compare direct CLI debugging against workflow debugging for the prototype.
- [ ] Document which operations became easier and which became harder.
- [ ] Decide whether to migrate more intake commands, keep workflow optional, or redesign the runner boundary.
- [ ] Update the design guide with observed outcomes.
- [ ] Upload an updated RAGEVAL-006 bundle to reMarkable if the implementation materially changes the design.
