# Changelog

## 2026-05-29

- Initial workspace created


## 2026-05-29

Created workflow/scraper integration design guide and diary for modular retryable intake prototype.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/29/RAGEVAL-006--workflow-scraper-integration-for-modular-retryable-intake/design-doc/01-workflow-scraper-intake-integration-design-and-implementation-guide.md — Primary design deliverable
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/29/RAGEVAL-006--workflow-scraper-integration-for-modular-retryable-intake/reference/01-implementation-diary.md — Investigation and delivery diary


## 2026-05-29

Validated RAGEVAL-006 docs with docmgr doctor and uploaded design bundle to reMarkable at /ai/2026/05/29/RAGEVAL-006.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/29/RAGEVAL-006--workflow-scraper-integration-for-modular-retryable-intake/design-doc/01-workflow-scraper-intake-integration-design-and-implementation-guide.md — Uploaded in bundle
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/29/RAGEVAL-006--workflow-scraper-integration-for-modular-retryable-intake/reference/01-implementation-diary.md — Uploaded in bundle


## 2026-05-29

Expanded RAGEVAL-006 into detailed implementation phases and task breakdowns before starting Phase 0.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/29/RAGEVAL-006--workflow-scraper-integration-for-modular-retryable-intake/tasks.md — Detailed phase/task breakdown


## 2026-05-29

Phase 0: added scraper dependency spike and custom echo runner with scheduler integration tests.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/go.mod — Added scraper/go-go-goja local compatibility wiring for Phase 0
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/workflow/echo_runner.go — Custom scraper runner compatibility spike
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/workflow/echo_runner_test.go — Temporary SQLite scheduler integration tests


## 2026-05-29

Phase 1 first slice: added typed intake runner dispatch and durable chunk_document operation over the existing chunking service.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/workflow/intake_runner.go — Intake runner dispatch and chunk_document service wrapper
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/workflow/intake_runner_test.go — Workflow integration tests for chunk_document
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/workflow/ops.go — Typed op input/output structs for first workflow operation


## 2026-05-29

Completed Phase 1 by adding compute_embeddings and build_bm25 workflow operations with dependency-chain tests.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/workflow/intake_runner.go — Added compute_embeddings and build_bm25 service-backed operations
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/workflow/intake_runner_test.go — Added workflow dependency-chain tests for embedding freshness and BM25 metadata
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/workflow/ops.go — Expanded workflow op contracts for embeddings and BM25


## 2026-05-29

Completed Phase 2 by adding workflow submission, worker, status, and ops CLI commands plus fake-provider smoke coverage.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/cmd/rag-eval/cmds/workflow/root.go — Adds operator-facing workflow command group
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/workflow/engine.go — Centralizes scheduler construction and runner registration
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/workflow/submit.go — Builds durable intake workflow topology for chunking embeddings and BM25


## 2026-05-29

Tightened Phase 2 downstream scoping by propagating document ID filters through embedding and BM25 service paths.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/db/queries.go — Adds chunk listing filter by document IDs for workflow-scoped embeddings
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/db/search_queries.go — Adds document ID filter for workflow-scoped BM25 indexing
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/embedding/service.go — Accepts document IDs in embedding compute requests
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/search/bm25.go — Passes document IDs into BM25 chunk selection

