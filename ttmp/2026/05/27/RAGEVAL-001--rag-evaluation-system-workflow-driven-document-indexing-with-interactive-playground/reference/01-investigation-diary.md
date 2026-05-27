---
title: "Investigation diary"
doc_type: reference
status: active
intent: long-term
topics: [rag, embeddings, evaluation, workflow, playground, search]
ticket: RAGEVAL-001
created_at: 2026-05-27
updated_at: 2026-05-27
---

# Diary

## Goal

Investigation diary for the RAG Evaluation System (RAGEVAL-001). Captures the research phase: studying existing codebases, understanding the architecture, and producing the design document.

## Step 1: Codebase Research and Architecture Mapping

I studied the four existing repositories in the workspace to understand their capabilities, integration points, and patterns that the RAG Evaluation System should reuse. This research forms the foundation for the design document.

### Prompt Context

**User prompt (verbatim):** "Create a new docmgr ticket to create a RAG evaluation system. I have never really done rag stuff, and as such I would like to see how embedding similarities behave, how to do reranking, etc... I want to build the following: - a workflow system to index / transform / compute embeddings / extract structured data / index documents, using the scraper system in ./scraper, as for example used in 2026-05-20--book-ocr - a RAG playground / web UI where I can exercise different parts of this workflow, and be able to query / interact / view "under the hood" data I want to use ./geppetto and potentially ./pinocchio (at least for the pinocchio profiles), and you can look at ./2026-05-21--readwise-viewer for some examples of using embeddings. Use react + rtk-query + tailwind + storybook for the UI, with a retro macos1 style. Create a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet points and pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and then upload to remarkable. Keep a diary as you work. Work within 2026-05-27--rag-evaluation-system/ and ask for permission to work outside if necessary. Study things and write up a great detailed rich document."

**Assistant interpretation:** Create a docmgr ticket for a RAG evaluation system with a workflow engine (using scraper) and a web playground (React + RTK Query + Tailwind + Storybook, retro macOS style). Study the existing codebases (scraper, geppetto, pinocchio, book-ocr, readwise-viewer) and produce a comprehensive design document suitable for a new intern.

**Inferred user intent:** The user wants a complete design and implementation guide for a RAG evaluation system that teaches RAG concepts while specifying concrete architecture. They want to understand how embeddings, similarity search, reranking, and evaluation work — and they want an interactive playground to explore these concepts hands-on.

### What I did

- Created docmgr ticket RAGEVAL-001
- Created design doc and diary doc in the ticket workspace
- Studied all four existing repositories in detail:
  - **scraper**: Read engine model types, store interface, scheduler, JS runner, site manifests, API handlers, frontend package.json
  - **geppetto**: Read embeddings provider interface, OpenAI/Ollama providers, settings factory, disk cache, config/settings, engine profiles
  - **pinocchio**: Read README, profile config format, JS runtime examples
  - **book-ocr**: Read README for workflow-backed OCR patterns
  - **readwise-viewer**: Read README, design doc for hybrid search implementation
- Wrote the comprehensive design document (82KB) covering:
  - RAG concepts for newcomers (Section 1)
  - Current-state architecture of all four repos (Section 2)
  - Proposed system architecture (Section 3)
  - Data model design with full SQL schema (Section 4)
  - Workflow design with site manifests and verb/script examples (Section 5)
  - Search service design with Bleve mappings and hybrid search (Section 6)
  - HTTP API design with full endpoint listing and JSON contracts (Section 7)
  - Frontend design with four main views and component lists (Section 8)
  - Go backend architecture with package layout and server bootstrap (Section 9)
  - Evaluation metrics (Section 10)
  - Reranking deep dive (Section 11)
  - Embedding similarity deep dive (Section 12)
  - 30-day implementation plan in 7 phases (Section 13)
  - Key files reference (existing and new) (Section 14)
  - Risks, open questions, alternatives (Sections 15-16)
  - Testing strategy and glossary (Sections 17-18)

### Why

The user explicitly asked for a "detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is." This required thorough study of existing codebases to ground the design in real integration points and patterns, rather than abstract architecture.

### What worked

- The Readwise Viewer design doc was an excellent reference for the Bleve hybrid search architecture. Its concrete code patterns (mapping definitions, query construction, cache schema) could be directly adapted.
- The Book OCR README clearly demonstrated the scraper workflow pattern (discover → per-page ops → assemble → validate) which maps directly to the RAG pipeline (discover → ingest → extract → chunk → embed → index).
- Geppetto's embedding provider interface is clean and well-designed for our use case. The `SettingsFactory` from inference settings is exactly what we need for profile-driven embedding computation.

### What didn't work

- No blockers during the research phase.

### What I learned

- The scraper engine is more than a scraping tool — it's a general-purpose durable workflow engine. The RAG pipeline is conceptually similar to scraping (fetch → transform → store → index), so scraper is a natural fit.
- Geppetto's profile system and the Pinocchio profile registry format are the same system. Using Pinocchio's `~/.config/pinocchio/profiles.yaml` as the configuration source means the RAG system automatically picks up the same embedding models the user already has configured.
- Bleve vector search requires `-tags vectors` and FAISS native library setup. This is a known pain point. The design keeps vector code behind build tags and makes BM25-only search fully functional without FAISS.

### What was tricky to build

- The data model design needed to support multi-configuration comparison (same document, different chunking/embedding/indexing strategies) without making the schema overly complex. The `strategy_id` foreign key pattern solves this but requires careful documentation.
- Balancing the "for a new intern" requirement with technical depth. The document needs to explain RAG concepts from scratch (what is an embedding, what is cosine similarity) while also providing concrete code and schema for implementation. I structured it with progressive depth: concepts first, then architecture, then concrete code.

### What warrants a second pair of eyes

- The SQLite schema (Section 4) is the foundation everything builds on. Verify that the composite primary keys on `chunk_embeddings` and `chunk_enrichments` correctly support multi-strategy storage.
- The workflow design (Section 5) assumes that scraper's Goja runtime can call Geppetto embedding APIs via host functions. This needs a concrete host API implementation.
- The evaluation metrics computation (Section 10) should have concrete test cases with known expected values.
- The `SearchDocument` type's `SearchText` field structure (Title/Author/Source/Tags/Summary/Content) should be validated against actual embedding model behavior.

### What should be done in the future

- Implement Phase 0 (project scaffolding) and Phase 1 (SQLite schema + document ingestion) to validate the architecture.
- Add t-SNE/UMAP visualization of embedding clusters in the playground.
- Investigate Ollama's new `/api/embed` batch endpoint for more efficient local embedding computation.
- Add support for image embeddings (CLIP) if multi-modal documents are needed.
- Consider adding a "prompt playground" where users can experiment with the enrichment prompt (summary, entities, hypothetical questions) and see how changes affect retrieval quality.

### Code review instructions

- Start with the design doc: `ttmp/2026/05/27/RAGEVAL-001--.../design-doc/01-rag-evaluation-system-architecture-and-implementation-guide.md`
- Key sections to review: Section 4 (data model), Section 6 (search service), Section 11 (reranking)
- Validate that the SQL schema composite keys support the intended multi-configuration comparison workflow
- Check that the API contracts (Section 7) are consistent with the frontend components (Section 8)

### Technical details

- Document size: ~83KB, 18 major sections
- SQLite tables designed: 9 (sources, documents, chunks, chunking_strategies, chunk_embeddings, chunk_enrichments, search_indexes, eval_queries, eval_runs, eval_results)
- API endpoints designed: ~25
- Frontend views designed: 4 (Pipeline Explorer, Embedding Inspector, Search Sandbox, Evaluation Dashboard)
- Implementation phases: 7, estimated 30 days

---

## Step 2: Phase 0 Implementation — Project Scaffolding

Implemented the Go server skeleton, React frontend with retro macOS theme, and development tooling. The system now starts, serves the SPA, and responds to API calls.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Proceed with implementing the design, committing at appropriate intervals, keeping a diary.

**Inferred user intent:** Build the system end-to-end, not just design it.

**Commit (code):** 308752e — "feat: scaffold Go server with SQLite, API handlers, and SPA embedding"
**Commit (code):** 14a4620 — "feat: add React frontend with retro macOS theme and RTK Query"
**Commit (code):** e375665 — "feat: add Makefile and devctl config for local development"

### What I did

- Created Go module `github.com/go-go-golems/rag-evaluation-system`
- Implemented `cmd/rag-eval/main.go`: HTTP server with graceful shutdown, go:embed SPA
- Implemented `internal/config/config.go`: env-based configuration
- Implemented `internal/db/db.go`: SQLite with WAL mode, 10-table schema, auto-migration
- Implemented `internal/api/handlers.go`: health, sources CRUD, documents list/detail, chunks list
- Implemented `internal/web/spa.go`: go:embed SPA handler with fallback routing
- Set up React 19 + Vite + Tailwind CSS 4 + RTK Query + Recharts frontend
- Created retro macOS design system in `src/index.css`: windows, title bars, buttons, inputs, tables, scrollbars, menu bar, tab bar
- Created `MacWindow`, `MacMenuBar`, `MacButton` retro components
- Created `PipelineView` with live sources/documents tables (RTK Query powered)
- Created placeholder views for Embeddings, Search, Evaluation
- Added Makefile and devctl.yaml

### Why

Phase 0 establishes the working skeleton. Before adding RAG-specific logic, we need a server that starts, a database that stores data, API endpoints that return JSON, and a frontend that renders.

### What worked

- The retro macOS theme looks authentic. The classic Mac System 7 design language (3D beveled buttons, horizontal stripe title bars, platinum gray background) translates well to CSS.
- RTK Query + SQLite via Go API is a clean stack. The API proxy in Vite dev config makes local development seamless.
- SQLite WAL mode with `MaxOpenConns(1)` is the right pattern for single-writer Go+SQLite.

### What didn't work

- `go:embed` with relative paths was tricky. The embed path is relative to the Go source file, not the module root. Solved by putting `spa.go` in `internal/web/` alongside the `dist/` directory.
- Vite `react-ts` template created vanilla TypeScript (not React) the first time. Had to manually add React dependencies and create `.tsx` files.

### What I learned

- `go:embed` patterns need the embed directory to contain actual files matching the glob pattern. An empty directory won't compile.
- Vite 8 + Tailwind CSS 4 uses `@import "tailwindcss"` instead of the old `@tailwind base/components/utilities` directives.
- The `pnpm create vite` template selection can be flaky; better to verify the template produced React files before installing dependencies.

### What was tricky to build

- The go:embed path resolution. The directive resolves relative to the package directory, so `cmd/rag-eval/main.go` can't embed `../../web/dist/`. Moving the embed to a dedicated `internal/web/` package and importing it from main solved this cleanly.
- TypeScript import path resolution in the React app. Components in nested directories need `../retro/MacWindow` not `./retro/MacWindow`.

### What warrants a second pair of eyes

- The `internal/web/spa.go` embeds `dist/*` which only works when the frontend has been built. During Go development without frontend changes, the placeholder `index.html` is sufficient, but this should be documented.
- SQLite `MaxOpenConns(1)` means only one database connection. This is correct for write safety but may bottleneck under concurrent reads. Consider `MaxOpenConns(2)` with `MaxIdleConns(1)` if read contention appears.

### What should be done in the future

- Add Storybook stories for retro components.
- Add document ingestion API (POST /api/v1/documents) and file-system scanner.
- Implement chunking service (Phase 2).

### Code review instructions

- Build: `make build` then `make run`
- Test API: `curl http://127.0.0.1:8772/api/v1/health`
- Test frontend: `make web-dev` then open http://127.0.0.1:5173
- Key files: cmd/rag-eval/main.go, internal/db/db.go, internal/api/handlers.go, web/src/index.css

### Technical details

- Go binary: ~10MB with SQLite + embedded SPA
- SQLite schema: 10 tables with foreign keys and ON DELETE CASCADE
- API endpoints implemented: 5 (health, list/create sources, list documents, get document, list chunks)
- Frontend bundle: ~270KB JS + ~10KB CSS (gzip: ~86KB + ~3KB)

---

## Step 3: Glazed CLI Pivot, Chunking Incident, and Audit Recovery

The implementation moved beyond Phase 0 into Glazed CLI commands, ingestion, and chunking. During validation, a fixed-size chunking bug caused unbounded chunk generation near the end of a document, which likely caused the observed memory spike and killed process.

I corrected the chunker termination and validation behavior, added regression tests, committed the fix, and wrote a recovery report so the next engineer can continue from a factual state instead of continuing the sprint momentum blindly.

### Prompt Context

**User prompt (verbatim):** "Ok, we're taking over for your little brother who has lost the saucce. Assess, correct, go over the whole plan actually, and make a full report describing what is good, what is bad, what they didn't know, what they should have done, what to do next time, and what to do now.

Save in the ticket, and upload to remarkable."

**Assistant interpretation:** Stop feature expansion, diagnose the implementation quality and memory incident, fix the immediate dangerous bug, write a candid recovery report in the ticket, and publish it to reMarkable.

**Inferred user intent:** Regain control of a runaway implementation session, preserve useful work, identify incorrect assumptions, and establish a safe next plan.

**Commit (code):** cbae145 — "feat: add chunking commands and guard overlap loops"

### What I did

- Inspected the uncommitted chunking implementation and current ticket state.
- Confirmed that fixed-size chunking could loop at the tail when overlap was enabled.
- Added overlap validation and explicit termination when `end >= totalRunes`.
- Added regression tests in `internal/chunking/chunker_test.go`.
- Fixed SQLite DB parent directory creation formatting and behavior in `internal/db/db.go`.
- Validated with constrained commands:
  - `GOMAXPROCS=2 GOMEMLIMIT=768MiB go test ./internal/chunking ./internal/db ./internal/ingest -count=1 -timeout 20s`
  - `GOMAXPROCS=2 GOMEMLIMIT=768MiB go build ./cmd/rag-eval`
- Smoke-tested CLI source creation, scan, and chunk apply against `/tmp/rag-eval-audit.db`.
- Created `analysis/01-implementation-audit-and-recovery-plan.md`.

### Why

The memory incident made it unsafe to continue adding embeddings or search. The right response was to isolate the immediate correctness bug, add tests, document the deeper architectural gaps, and give the next implementation pass a stabilization plan.

### What worked

- The stack trace pointed directly at `generateChunkID` inside `FixedSizeChunker.Chunk`, with an absurd chunk index for a small document.
- A small unit test reproduced the infinite-loop shape without needing a large corpus.
- After the fix, the chunker test suite passes quickly and CLI chunking completes successfully on a scanned source file.

### What didn't work

- An initial attempt to constrain the process with `ulimit -v 1048576` caused `runtime/cgo: pthread_create failed: Resource temporarily unavailable`; that was too aggressive for this Go+cgo+SQLite binary and produced a runtime abort unrelated to the original application bug.
- The current schema cannot support multiple chunking strategies per document because `chunks` lacks a real `strategy_id` column and has `UNIQUE(document_id, chunk_index)`.
- The CLI and HTTP chunking code still duplicate behavior instead of sharing a service layer.

### What I learned

- Fixed-size chunking with overlap must treat `end == total_length` as a terminal condition; otherwise the tail can be emitted forever.
- Workflow-oriented RAG operations must be idempotent from the start because retries are part of the intended runtime model.
- Glazed CLI lockstep is valuable, but it must share domain services with HTTP handlers to avoid behavior drift.

### What was tricky to build

The tricky part was distinguishing a resource-limit failure from the real memory incident. The `ulimit` experiment failed because the Go runtime could not create enough threads under the artificial virtual-memory cap. The real bug was found by reading the stack trace and seeing repeated chunk ID generation at high indexes for a small document, then reasoning through the overlap tail case.

The fix required both a terminal condition (`end >= totalRunes`) and invalid parameter rejection (`overlap >= chunk_size`). Without both, one class of runaway could remain.

### What warrants a second pair of eyes

- `internal/db/db.go` chunk schema: `chunks` should include `strategy_id` and unique `(document_id, strategy_id, chunk_index)` before embeddings are implemented.
- `cmd/rag-eval/cmds/chunk/apply.go` should move persistence behavior into an `internal/services/chunking` service shared by CLI and HTTP.
- `internal/api/handlers.go` should not duplicate chunk application logic.
- Write semantics for source scan and chunk apply need idempotency decisions before scraper workflow integration.

### What should be done in the future

- Freeze Phase 2 expansion until chunk identity and service-layer sharing are fixed.
- Reconcile tasks so completed checkboxes do not overclaim Storybook/UI work.
- Add retry/idempotency notes to every future workflow operation.
- Add bounded output modes for chunking, embeddings, and search commands.

### Code review instructions

- Start with `internal/chunking/chunker.go`, especially fixed-size termination and overlap validation.
- Review `internal/chunking/chunker_test.go` to confirm the memory incident now has a regression test.
- Review `analysis/01-implementation-audit-and-recovery-plan.md` for the full recovery plan.
- Validate with:
  - `GOMAXPROCS=2 GOMEMLIMIT=768MiB go test ./internal/chunking ./internal/db ./internal/ingest -count=1 -timeout 20s`
  - `GOMAXPROCS=2 GOMEMLIMIT=768MiB go build ./cmd/rag-eval`

### Technical details

- Corrected commit: `cbae145`.
- New report: `ttmp/2026/05/27/RAGEVAL-001--rag-evaluation-system-workflow-driven-document-indexing-with-interactive-playground/analysis/01-implementation-audit-and-recovery-plan.md`.
- Core bug pattern: after a final capped chunk reaches `totalRunes`, subtracting overlap can produce another start whose next end is still `totalRunes`, repeating forever.

---

## Step 4: P0 Stabilization — Strategy-Aware Chunks and Shared Services

After the audit, I continued with the highest-priority recovery items instead of moving to embeddings. The goal was to make ingestion/chunking safe enough for future workflow retries and to stop CLI and HTTP from carrying separate copies of domain behavior.

This step changed chunk identity, made chunk application rerun-safe, made source/document ingestion idempotent, bounded chunk output by default, and introduced source/document/chunking services that both Glazed CLI commands and HTTP handlers use.

### Prompt Context

**User prompt (verbatim):** "continue. commit at appropriate intervals"

**Assistant interpretation:** Continue implementing the recovery/stabilization plan, making focused commits after each validated unit.

**Inferred user intent:** Keep momentum, but avoid another uncontrolled implementation burst by committing coherent, validated changes.

**Commit (code):** 3846818 — "fix: make chunks strategy-aware and rerun-safe"
**Commit (code):** 5f4e8a2 — "refactor: share chunk application service across CLI and HTTP"
**Commit (code):** 9b366ef — "fix: make source and document ingestion idempotent"
**Commit (code):** 320bb2e — "feat: bound chunk output for CLI and HTTP"
**Commit (code):** a82dc57 — "refactor: share source service across CLI and HTTP"
**Commit (code):** 2278179 — "refactor: share document service across CLI and HTTP"

### What I did

- Added first-class `strategy_id` to the `chunks` schema.
- Added `internal/db/migrations.go` to upgrade older development DBs whose chunks table lacked `strategy_id`.
- Changed chunk IDs to include `(document_id, strategy_id, chunk_index)`.
- Added `DeleteChunksForDocumentStrategy` so chunk application rebuilds derived state for one document/strategy pair.
- Changed `InsertChunkingStrategy`, `InsertSource`, and `InsertDocument` to upsert.
- Changed scanner document IDs to use stable relative paths and store that path as `external_id`.
- Added `internal/services/chunking.Service` and routed CLI/HTTP chunk application through it.
- Added `internal/services/source.Service` and routed CLI/HTTP source create/scan through it.
- Added `internal/services/document.Service` and routed CLI/HTTP document list/get/chunks through it.
- Added bounded chunk output controls:
  - `rag-eval chunk apply --emit preview|full|none --limit N --preview-runes N`
  - `rag-eval document chunks --emit preview|full|none --limit N --preview-runes N`
  - HTTP chunk endpoint now defaults to summary output and only returns chunks when requested.

### Why

Embeddings depend on stable chunk identity. If the same document can be chunked with several strategies, `strategy_id` must be part of the relational identity, not buried inside JSON. Future scraper workflow retries also require idempotent writes; plain inserts would fail or duplicate work on retry.

### What worked

- CLI chunk application can now be rerun for the same `(document_id, strategy_id)` without duplicate/primary-key failures.
- Source scans can be rerun without duplicating documents.
- CLI and HTTP chunking behavior now shares `internal/services/chunking.Service`.
- The HTTP chunk endpoint no longer returns full chunk text by default.

### What didn't work

- The schema still uses hand-rolled migration helpers rather than a real migration framework/version table. This is acceptable for the current dev ticket but should be improved before long-lived deployments.
- Existing old DBs with chunk embeddings were not deeply tested through the table-rebuild migration; the current migration is intended for development DB recovery.

### What I learned

- The earlier design invariant (`strategy_id` as first-class identity) had to be enforced before embeddings; otherwise embeddings would be attached to ambiguous chunks.
- Output bounding belongs in the first implementation of any command that can emit content-bearing rows.
- Service extraction is easiest immediately after behavior exists but before frontend/workflow layers depend on it.

### What was tricky to build

The tricky part was changing `chunks` identity after the first schema already existed. New DBs can use the corrected schema directly, but old dev DBs need repair. I added an `ensureChunksStrategyID` helper that detects missing `strategy_id`, creates a legacy strategy, rebuilds the table, and preserves existing rows.

Another tricky detail was making chunk IDs strategy-aware. Without including `strategy_id` in the hash, two different chunking strategies for the same document would produce colliding chunk IDs for chunk index `0`.

### What warrants a second pair of eyes

- `internal/db/migrations.go`: table rebuild strategy and foreign-key handling should be reviewed before treating it as production migration infrastructure.
- `internal/services/chunking/service.go`: chunk application deletes and rebuilds derived chunks for the strategy; this is good for retries but must be considered when embeddings/enrichments exist.
- HTTP error mapping is still coarse. Service errors should eventually map to typed API errors (404 vs 400 vs 500).

### What should be done in the future

- Add a migration version table or migrate package before schemas become more complex.
- Add tests for service-layer behavior with temporary SQLite databases.
- Decide whether rechunking should cascade-delete embeddings/enrichments or refuse when downstream derived state exists.
- Extend bounded output patterns to embeddings/search when those commands are implemented.

### Code review instructions

- Start with `internal/db/db.go` and `internal/db/migrations.go` for chunk schema identity.
- Review `internal/services/chunking/service.go` for the canonical chunk application path.
- Review `cmd/rag-eval/cmds/chunk/apply.go` and `internal/api/handlers.go` to confirm they are now adapters.
- Validate with:
  - `GOMAXPROCS=2 GOMEMLIMIT=768MiB go test ./internal/chunking ./internal/db ./internal/ingest ./internal/services/chunking ./internal/services/source ./internal/services/document -count=1 -timeout 20s`
  - `GOMAXPROCS=2 GOMEMLIMIT=768MiB go build ./cmd/rag-eval`

### Technical details

- Completed stabilization task IDs: 44, 45, 46.
- Partial bounded-output task: chunking outputs are bounded; embeddings/search outputs remain future work because those features are not implemented yet.
- Smoke tests used temporary databases under `/tmp` and constrained Go memory/CPU using `GOMAXPROCS=2 GOMEMLIMIT=768MiB`.

---

## Step 5: Service-Layer Test Coverage and Detailed Continuation Tasks

I added a more detailed continuation task list before coding so the ticket now has explicit near-term work for service tests, migration tests, and Phase 3 embedding slices. Then I implemented the test coverage that should have existed before moving from chunking into embeddings.

The tests use temporary SQLite databases, not mocks, because the highest-risk behavior here is relational identity, idempotent upserts, and migration behavior. This catches the real failure modes we care about before wiring more workflow or embedding layers on top.

### Prompt Context

**User prompt (verbatim):** "add tasks and phases to the ticket, if necessary, detailed task list, then work on them step by step, committing at appropriate intervals, keeping a detailed diary as you work"

**Assistant interpretation:** Reconcile and extend the ticket task list, then continue implementation one validated task group at a time with focused commits and diary updates.

**Inferred user intent:** Keep implementation disciplined and auditable after the recovery pass, with enough task granularity for handoff and review.

**Commit (docs):** 632f601 — "docs: add detailed continuation tasks"
**Commit (code):** 978f680 — "test: cover service idempotency and chunk migrations"

### What I did

- Added detailed tasks for:
  - source service idempotency tests,
  - chunk service strategy/idempotency tests,
  - document service tests,
  - legacy chunk migration tests,
  - Phase 3 embedding research and implementation slices.
- Added `internal/services/source/service_test.go`.
- Added `internal/services/chunking/service_test.go`.
- Added `internal/services/document/service_test.go`.
- Added `internal/db/migrations_test.go`.
- Marked tasks 48, 49, 50, and 51 complete.

### Why

The previous sprint had enough behavior to be useful, but not enough service-layer tests to safely start embeddings. Embeddings will depend on stable document/chunk identity and rerun-safe derived state, so these tests are the safety rail for Phase 3.

### What worked

- Temporary SQLite tests are fast enough and exercise the real schema.
- Source create and scan idempotency are covered with real filesystem files.
- Chunk application is covered for rerunning the same strategy and keeping different strategies separate.
- Document service tests confirm list/get/chunks behavior over real stored data.
- The migration test confirms a legacy `chunks` table without `strategy_id` is upgraded and preserves strategy metadata from `boundaries_json`.

### What didn't work

- No major blockers. The only tradeoff is some duplicated test setup helpers across service packages. That is acceptable for now because it keeps tests local and avoids introducing a production test utility package.

### What I learned

- The service boundary is now testable enough to support the embedding phase.
- The migration helper is adequate for development DB recovery, but it still should not be mistaken for a full production migration framework.
- Testing against SQLite directly is worth the small overhead because schema mistakes were the primary risk.

### What was tricky to build

The migration test required constructing a partial legacy database by hand, then running the current `Migrate()` path against it. This matters because `CREATE TABLE IF NOT EXISTS` alone would not alter an old table; the explicit `ensureChunksStrategyID` path is the behavior under test.

### What warrants a second pair of eyes

- `internal/db/migrations_test.go` should be reviewed alongside `internal/db/migrations.go` to confirm the legacy upgrade path is good enough for development DBs.
- The duplicated `openTestQueries` helpers may be worth consolidating later if more services are added.

### What should be done in the future

- Start Phase 3.0 by rereading current Geppetto embedding/profile APIs and Pinocchio profile wiring before implementing the embedding provider factory.
- Keep embedding service tests temporary-SQLite-backed as well.
- Add CLI/HTTP smoke tests for embeddings only after service-level tests exist.

### Code review instructions

- Review tests first:
  - `internal/services/source/service_test.go`
  - `internal/services/chunking/service_test.go`
  - `internal/services/document/service_test.go`
  - `internal/db/migrations_test.go`
- Validate with:
  - `GOMAXPROCS=2 GOMEMLIMIT=768MiB go test ./internal/db ./internal/ingest ./internal/chunking ./internal/services/source ./internal/services/chunking ./internal/services/document -count=1 -timeout 30s`
  - `GOMAXPROCS=2 GOMEMLIMIT=768MiB go build ./cmd/rag-eval`

### Technical details

- Completed task IDs: 48, 49, 50, 51.
- Added Phase 3 continuation task IDs: 52-56.

---

## Step 6: Phase 3 Start — Embedding Plan, Provider Resolver, and Compute Service

I moved into Phase 3 only after adding service-layer and migration tests. The first Phase 3 step was intentionally split into research/planning, provider construction, and service-level embedding persistence so no CLI or HTTP adapter could hide core behavior problems.

The embedding implementation now has a Geppetto-backed provider resolver, Pinocchio/profile-registry support following the Geppetto smoke-test pattern, vector BLOB encoding helpers, and a compute service that batches chunks, skips fresh embeddings using `text_hash`, validates dimensions, and upserts into SQLite.

### Prompt Context

**User prompt (verbatim):** (see Step 5)

**Assistant interpretation:** Continue executing detailed tasks in order, committing after coherent validated units and updating the diary.

**Inferred user intent:** Proceed into embeddings only after documenting and validating the backend foundation.

**Commit (docs):** da55bd9 — "docs: add Phase 3 embedding integration plan"
**Commit (code):** d204a5e — "feat: add Geppetto embedding provider resolver"
**Commit (code):** 785462b — "feat: add embedding compute service with staleness checks"

### What I did

- Read Geppetto embedding/provider APIs:
  - `geppetto/pkg/embeddings/embeddings.go`
  - `geppetto/pkg/embeddings/settings_factory.go`
  - `geppetto/pkg/embeddings/settings_validation.go`
  - `geppetto/cmd/examples/embedding-profile-smoke/main.go`
  - `geppetto/pkg/steps/ai/settings/settings-inference.go`
- Wrote `analysis/02-phase-3-embedding-integration-plan.md`.
- Implemented `internal/services/embedding/provider.go`:
  - direct config resolution,
  - profile-backed resolution via Geppetto engineprofiles,
  - default Pinocchio profile registry path,
  - validation through `embeddings.ValidateInferenceSettingsForEmbeddings`,
  - provider construction via `embeddings.NewSettingsFactoryFromInferenceSettings`.
- Implemented `internal/services/embedding/vector.go` for deterministic little-endian float32 vector encoding/decoding.
- Added DB helpers for chunk embedding staleness checks and upsert persistence.
- Implemented `internal/services/embedding/service.go`:
  - list chunks by strategy,
  - compute text hashes,
  - skip fresh embeddings unless forced,
  - batch calls to the provider,
  - verify vector dimensions,
  - upsert embedding blobs.
- Added fake-provider SQLite tests for the embedding service.

### Why

The original requirement was to use Geppetto for embeddings and Pinocchio at least for profile configuration. The safest way to satisfy that is to use Geppetto's provider interface and settings factory while keeping network calls out of unit tests.

### What worked

- Geppetto already exposes the right seam: `embeddings.Provider` with single and batch generation plus model metadata.
- The `embedding-profile-smoke` example provided a clear pattern for profile registry resolution and base-profile overlays.
- Fake-provider tests let us verify batching, staleness, vector dimensions, and SQLite persistence without OpenAI/Ollama.

### What didn't work

- This step does not yet expose embeddings through CLI or HTTP; it deliberately stops at the provider/service layer.
- The provider resolver constructs real OpenAI/Ollama providers, but live requests are not tested here and should remain opt-in smoke tests.

### What I learned

- OpenAI profile validation should fail early if no `openai-api-key` is configured.
- Ollama provider construction does not require a network call; the network call happens when embeddings are generated.
- Keeping `ProviderType` explicit in the compute request avoids ambiguity in the `chunk_embeddings` identity key.

### What was tricky to build

Profile-backed resolution has three modes: explicit profile, base profile plus embedding overlay, and direct config. The implementation keeps these separate so errors remain understandable. It also returns a `Close` function for profile registry chains, mirroring the Geppetto smoke-test pattern.

The other tricky part was staleness: the service must compare `text_hash` for the exact `(chunk_id, strategy_id, provider, model, dimensions)` identity, not just by chunk ID.

### What warrants a second pair of eyes

- `internal/services/embedding/provider.go`: profile/base-profile merge behavior should be reviewed against real Pinocchio profile files.
- `internal/services/embedding/service.go`: confirm the primary key and staleness semantics before adding search indexes.
- `internal/db/queries.go`: embedding upsert shape should be reviewed before adding vector search or export tooling.

### What should be done in the future

- Add Glazed `embedding compute` and `embedding similarity` commands using the service.
- Add HTTP embedding endpoints as adapters over the same service.
- Add an opt-in live smoke command for Ollama/OpenAI, but do not run it in unit tests.
- Add frontend Embedding Inspector only after CLI/HTTP contracts are stable.

### Code review instructions

- Start with `ttmp/.../analysis/02-phase-3-embedding-integration-plan.md`.
- Then review:
  - `internal/services/embedding/provider.go`
  - `internal/services/embedding/service.go`
  - `internal/services/embedding/service_test.go`
  - `internal/db/queries.go`
- Validate with:
  - `GOMAXPROCS=2 GOMEMLIMIT=1024MiB go test ./internal/db ./internal/ingest ./internal/chunking ./internal/services/source ./internal/services/chunking ./internal/services/document ./internal/services/embedding -count=1 -timeout 60s`
  - `GOMAXPROCS=2 GOMEMLIMIT=1024MiB go build ./cmd/rag-eval`

### Technical details

- Completed task IDs: 52, 53, 54.
- New Geppetto dependency added to `go.mod`/`go.sum`.
- Embeddings are encoded as little-endian float32 BLOBs.
- Unit tests use fake providers and temporary SQLite databases; no network calls are made.

---

## Step 7: Embedding Adapters — Glazed Compute Command and HTTP Compute Endpoint

After the embedding service was tested, I added the first CLI and HTTP adapters. Both adapters resolve a Geppetto embedding provider, call the shared embedding compute service, and emit bounded summary output rather than vectors.

This keeps the lockstep principle intact: the CLI and HTTP paths are wrappers around the same service, not independent implementations.

### Prompt Context

**User prompt (verbatim):** (see Step 5)

**Assistant interpretation:** Continue stepwise from the detailed task list, committing each adapter slice after validation.

**Inferred user intent:** Make embeddings usable from operator tooling and HTTP while preserving the shared-service architecture.

**Commit (code):** c0d4bd2 — "feat: add Glazed embedding compute command"
**Commit (code):** 2a1f752 — "feat: add embedding compute HTTP endpoint"

### What I did

- Added `cmd/rag-eval/cmds/embedding/root.go`.
- Added `cmd/rag-eval/cmds/embedding/compute.go`.
- Registered `rag-eval embedding compute` in the root command.
- Added `POST /api/v1/embeddings/compute`.
- Added task 57 for the CLI compute command and task 58 for the HTTP endpoint.
- Marked tasks 57 and 58 complete.

### Why

The embedding service is only useful if operators and the web UI can call it. The CLI command provides a scriptable operator surface; the HTTP endpoint is the backend contract for the future Embedding Inspector.

### What worked

- The CLI help renders correctly with Glazed output flags.
- The HTTP endpoint can construct an Ollama provider and return a zero-work summary for a missing/empty strategy without making a provider network call.
- Both adapters route through `internal/services/embedding.Service`.

### What didn't work

- No live embedding provider smoke test was run. That is intentional: live OpenAI/Ollama calls should be explicit and opt-in.
- `embedding similarity` is not implemented yet, so the broader task 55 remains open.
- The frontend Embedding Inspector is not implemented yet, so task 56 remains open.

### What I learned

- For a strategy with no chunks, the service returns a valid summary without provider batch calls. This is useful for HTTP smoke testing without requiring Ollama/OpenAI availability.
- The CLI adapter should remain summary-oriented; vectors are persisted, not printed.

### What was tricky to build

The adapter has many provider/profile flags. To avoid duplicating semantics, all provider resolution is delegated to `internal/services/embedding.ResolveProvider`. The command only translates Glazed settings into `ProviderConfig` and emits the result summary.

### What warrants a second pair of eyes

- HTTP error mapping is still coarse: provider-resolution errors are 400, compute errors are 500. Typed service errors would improve this.
- The CLI flag names should be reviewed before they become public API.

### What should be done in the future

- Add `embedding similarity` once vector retrieval helpers exist.
- Add frontend Embedding Inspector first slice.
- Add optional live smoke documentation for Ollama/OpenAI.

### Code review instructions

- Review `cmd/rag-eval/cmds/embedding/compute.go` and `internal/api/handlers.go`.
- Validate with:
  - `GOMAXPROCS=2 GOMEMLIMIT=1024MiB go build ./cmd/rag-eval`
  - `./rag-eval embedding compute --help`
  - HTTP smoke against an empty strategy using direct Ollama config.

### Technical details

- Completed task IDs: 57, 58.
- Broader tasks 55 and 56 remain open because similarity and frontend work are not complete.
