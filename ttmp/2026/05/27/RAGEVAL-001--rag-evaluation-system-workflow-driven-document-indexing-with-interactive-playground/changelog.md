# Changelog

## 2026-05-27

- Initial workspace created


## 2026-05-27

Created design document and investigation diary for RAG Evaluation System. Document covers: RAG concepts, current-state architecture of 4 repos, proposed system architecture, SQLite data model, workflow design, search service, HTTP API, frontend design (retro macOS), Go backend architecture, evaluation metrics, reranking deep dive, embedding similarity exploration, 30-day implementation plan.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/27/RAGEVAL-001--rag-evaluation-system-workflow-driven-document-indexing-with-interactive-playground/design-doc/01-rag-evaluation-system-architecture-and-implementation-guide.md — Primary design document (82KB)


## 2026-05-27

Uploaded design doc + diary bundle to reMarkable at /ai/2026/05/27/RAGEVAL-001


## 2026-05-27

Phase 0.1: Go server scaffolding with SQLite schema, API handlers, SPA embedding (commit 308752e)

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/cmd/rag-eval/main.go — Main server entrypoint with graceful shutdown
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/api/handlers.go — Health
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/db/db.go — SQLite schema with 10 tables and migration runner


## 2026-05-27

Phase 0 complete: Go server + React frontend + retro macOS theme + devctl + Makefile (commits 308752e, 14a4620, e375665)


## 2026-05-27

Audit recovery: fixed fixed-size chunker runaway overlap loop, added regression tests, committed cbae145, and wrote implementation audit/recovery report.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/chunking/chunker.go — Corrected overlap validation and tail termination
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/chunking/chunker_test.go — Regression coverage for the memory incident
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/27/RAGEVAL-001--rag-evaluation-system-workflow-driven-document-indexing-with-interactive-playground/analysis/01-implementation-audit-and-recovery-plan.md — Full assessment and recovery plan


## 2026-05-27

Uploaded audit recovery bundle to reMarkable at /ai/2026/05/27/RAGEVAL-001 as RAGEVAL-001 Audit Recovery Plan.pdf

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/27/RAGEVAL-001--rag-evaluation-system-workflow-driven-document-indexing-with-interactive-playground/analysis/01-implementation-audit-and-recovery-plan.md — Uploaded report bundle source


## 2026-05-27

P0 stabilization: made chunks first-class strategy-aware derived state, added rerun-safe chunk rebuilds, made source/document ingestion idempotent, and extracted source/document/chunking services shared by CLI and HTTP (commits 3846818, 5f4e8a2, 9b366ef, 320bb2e, a82dc57, 2278179).

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/db/db.go — Chunk strategy_id schema and migration entry point
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/db/migrations.go — Development DB upgrade for old chunks schema
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/chunking/service.go — Shared chunk application service
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/document/service.go — Shared document read service
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/source/service.go — Shared source create/scan service


## 2026-05-27

Added detailed next tasks for service-layer tests, legacy chunk migration coverage, and Phase 3 embedding work breakdown.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/27/RAGEVAL-001--rag-evaluation-system-workflow-driven-document-indexing-with-interactive-playground/tasks.md — Detailed continuation tasks added


## 2026-05-27

Added temporary-SQLite tests covering source idempotency, chunk strategy identity/idempotency, document service behavior, and legacy chunk migration (commit 978f680).

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/db/migrations_test.go — Legacy chunks strategy_id migration test
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/chunking/service_test.go — Chunk apply rerun-safety and strategy isolation tests
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/document/service_test.go — Document service list/get/chunks tests
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/source/service_test.go — Source create/scan idempotency tests


## 2026-05-27

Completed Phase 3.0 research plan for Geppetto/Pinocchio embedding integration before coding provider/service layers.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/27/RAGEVAL-001--rag-evaluation-system-workflow-driven-document-indexing-with-interactive-playground/analysis/02-phase-3-embedding-integration-plan.md — Phase 3 embedding integration plan


## 2026-05-27

Implemented Phase 3 provider resolver and embedding compute service with Geppetto provider construction, Pinocchio/profile resolution support, vector encoding, batching, text_hash freshness checks, and fake-provider SQLite tests (commits d204a5e, 785462b).

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/db/queries.go — Chunk embedding query/upsert helpers
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/embedding/provider.go — Geppetto/Pinocchio embedding provider resolver
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/embedding/service.go — Embedding compute service with batching and staleness checks
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/embedding/service_test.go — Fake-provider SQLite tests for embedding persistence
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/embedding/vector.go — Float32 vector BLOB encoding helpers


## 2026-05-27

Added Glazed embedding compute command backed by the embedding provider resolver and compute service (commit c0d4bd2).

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/cmd/rag-eval/cmds/embedding/compute.go — Glazed embedding compute command
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/cmd/rag-eval/main.go — Registered embedding command group


## 2026-05-27

Added HTTP POST /api/v1/embeddings/compute backed by the embedding provider resolver and compute service (commit 2a1f752).

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/api/handlers.go — Embedding compute HTTP endpoint

