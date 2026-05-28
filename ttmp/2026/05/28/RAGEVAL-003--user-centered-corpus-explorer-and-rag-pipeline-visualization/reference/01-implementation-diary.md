---
Title: Implementation Diary
Ticket: RAGEVAL-003
Status: active
Topics:
    - rag
    - playground
    - corpus
    - embeddings
    - search
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: internal/services/corpus/service.go
      Note: Corpus service implemented in Step 3
    - Path: ttmp/2026/05/28/RAGEVAL-003--user-centered-corpus-explorer-and-rag-pipeline-visualization/design-doc/01-corpus-explorer-and-pipeline-visualization-implementation-guide.md
      Note: Design guide authored in Step 1
    - Path: web/src/components/corpus/CorpusExplorerView.tsx
      Note: Corpus Explorer UI built in Step 4
    - Path: web/src/index.css
      Note: Monochrome theme redesign in Step 4
ExternalSources: []
Summary: Chronological diary for the Corpus Explorer and pipeline visualization design ticket.
LastUpdated: 2026-05-28T00:00:00Z
WhatFor: Track design decisions, evidence, and validation for the website improvement planning work.
WhenToUse: Read before implementing RAGEVAL-003 or updating the Corpus Explorer design.
---





# Diary

## Goal

This diary records the design workflow for adding a user-centered Corpus Explorer and RAG pipeline visualization to the website.

## Step 1: Ticket Creation and Intern-Ready Design Guide

I created a new docmgr ticket for improving the website from the perspective of a user who wants to learn, experiment with, and validate the RAG pipeline. I inspected the current frontend and API surfaces, then wrote a detailed implementation guide that explains the system, user intents, proposed information architecture, API contracts, frontend slices, validation plan, and design constraints.

The guide deliberately leaves visual freedom for a designer while being strict about technical clarity. The user should be able to see what was ingested, what each document became, where chunk boundaries are, which chunks have embeddings, and what can be tested next.

### Prompt Context

**User prompt (verbatim):** "ok, create a new docmgr ticket for adding this to the website / improving the website, starting from a user intent perspective as someone who wants to learn and experiment and validate the RAG pipeline here. 

Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.


Let enough freedom for a great designer to improve the site, and focus on like user intent and value"

**Assistant interpretation:** Create a new ticket for a user-centered website improvement, write an intern-ready design/implementation guide for a Corpus Explorer and pipeline visualization, and upload the package to reMarkable.

**Inferred user intent:** Turn the backend corpus/chunk/embedding pipeline into an understandable visual product surface rather than just CLI/API functionality.

### What I did

- Created ticket `RAGEVAL-003` titled `User-Centered Corpus Explorer and RAG Pipeline Visualization`.
- Added tasks for current-state analysis, user-intent design, implementation guide writing, and reMarkable upload.
- Inspected current frontend files:
  - `web/src/App.tsx`
  - `web/src/components/retro/MacMenuBar.tsx`
  - `web/src/components/pipeline/PipelineView.tsx`
  - `web/src/components/embeddings/EmbeddingsView.tsx`
  - `web/src/services/api.ts`
- Inspected current API/backend files:
  - `internal/api/handlers.go`
  - `internal/db/queries.go`
- Reused current corpus context from `RAGEVAL-002`, especially the dump-to-SQLite design guide and imported TTC corpus state.
- Wrote `design-doc/01-corpus-explorer-and-pipeline-visualization-implementation-guide.md`.

### Why

The system now has real data and useful backend operations, but users need a website surface that explains the pipeline. The design must start from user intent because the goal is not to expose every table; the goal is to help users understand, experiment with, and validate each transformation stage.

### What worked

- The current code already has clear places to extend:
  - add `CorpusExplorerView` under `web/src/components/corpus/`;
  - add a `corpus` menu item in `MacMenuBar`;
  - add RTK Query bindings in `web/src/services/api.ts`;
  - add corpus-specific HTTP endpoints in `internal/api/handlers.go`;
  - add a shared corpus service under `internal/services/corpus/`.
- The current embedding coverage and source-aware compute work provides the first useful coverage data model for the Corpus Explorer.

### What didn't work

- No website implementation was done in this step. This ticket is a design and implementation guide package, not a code implementation ticket.
- The current generic document API is insufficient for the proposed UI because it lacks source filtering, document-level chunk counts, content text, metadata, and chunk-level embedding status in one payload.

### What I learned

- The first useful visualization is not a chart-heavy dashboard. It is a Corpus Explorer that makes source, document, chunk, and embedding identities visible.
- The best product value is helping a user answer validation questions: what was ingested, what changed, how it was chunked, what is embedded, and what can be tested next.
- A designer should have freedom over layout and visual style, but exact IDs and embedding identity tuples must remain visible and copyable.

### What was tricky to build

The main challenge was balancing design freedom with technical specificity. If the document over-prescribes layout, it limits good design work. If it only describes user experience goals, it is not enough for an intern to implement. The guide therefore specifies required data, endpoints, files, validation commands, and component responsibilities, while leaving visual composition open.

### What warrants a second pair of eyes

- Review whether the proposed corpus API should be a separate route group or folded into the existing document endpoints.
- Review whether the first UI slice should be read-only or include bounded compute actions.
- Review whether product metadata should be copied into `documents.content_text` before the Corpus Explorer highlights missing searchable fields.

### What should be done in the future

- Implement the backend corpus service and endpoints.
- Implement the first read-only Corpus Explorer frontend slice.
- Add source-level embedding coverage to the Embedding Inspector or cross-link it from Corpus Explorer.
- Add search-result-to-corpus navigation after BM25 search exists.

### Code review instructions

- Start with the design guide:
  - `ttmp/2026/05/28/RAGEVAL-003--user-centered-corpus-explorer-and-rag-pipeline-visualization/design-doc/01-corpus-explorer-and-pipeline-visualization-implementation-guide.md`
- Then inspect the current files referenced by the guide:
  - `web/src/App.tsx`
  - `web/src/components/retro/MacMenuBar.tsx`
  - `web/src/services/api.ts`
  - `web/src/components/pipeline/PipelineView.tsx`
  - `web/src/components/embeddings/EmbeddingsView.tsx`
  - `internal/api/handlers.go`
  - `internal/db/queries.go`

### Technical details

- Ticket: `RAGEVAL-003`.
- Primary design document: `design-doc/01-corpus-explorer-and-pipeline-visualization-implementation-guide.md`.
- Proposed first endpoint group: `/api/v1/corpus/*`.
- Proposed first frontend view: `web/src/components/corpus/CorpusExplorerView.tsx`.

---

## Step 2: reMarkable Upload and Ticket Completion

I uploaded the RAGEVAL-003 design package to reMarkable as a bundled PDF. The bundle includes the intern-ready Corpus Explorer implementation guide and this diary.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Complete the requested documentation delivery by publishing the design package to reMarkable.

**Inferred user intent:** Make the website design plan reviewable away from the codebase, with a readable PDF bundle and table of contents.

### What I did

- Uploaded the design document and diary with `remarquee upload bundle`.
- Used remote directory `/ai/2026/05/28/RAGEVAL-003`.
- Marked task 4 complete.
- Updated the changelog with the upload location.

### Why

The user explicitly requested upload to reMarkable. A bundled PDF keeps the design guide and diary together for review.

### What worked

- Upload succeeded:
  - `OK: uploaded RAGEVAL-003 Corpus Explorer Design.pdf -> /ai/2026/05/28/RAGEVAL-003`

### What didn't work

- N/A.

### What I learned

- The design package is ready for review and implementation planning.

### What was tricky to build

N/A for this upload step. The only operational constraint was to avoid redundant reMarkable status/list calls and use the direct non-interactive upload command.

### What warrants a second pair of eyes

- Review the design guide before implementation begins, especially the proposed `/api/v1/corpus/*` endpoint boundaries and first read-only UI slice.

### What should be done in the future

- Implement the first read-only Corpus Explorer slice under a follow-up implementation task.
- Add frontend coverage wiring after the corpus endpoints exist.

### Code review instructions

- Review the uploaded PDF or the source markdown:
  - `design-doc/01-corpus-explorer-and-pipeline-visualization-implementation-guide.md`

### Technical details

- reMarkable path: `/ai/2026/05/28/RAGEVAL-003`.
- Uploaded bundle name: `RAGEVAL-003 Corpus Explorer Design.pdf`.

---

## Step 3: Backend Corpus Service and API Endpoints

I implemented the read-only corpus service with three query methods and registered three new API endpoints. The service provides source summaries with embedding coverage, a paginated document browser filtered by source, and a document detail endpoint that returns chunks with per-chunk embedding status.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Implement the backend corpus service and API endpoints needed by the Corpus Explorer frontend.

**Inferred user intent:** Provide the data layer that the frontend Corpus Explorer needs to show source, document, chunk, and embedding coverage data.

**Commit (code):** d0934ea — "feat(corpus): add corpus service and API endpoints for explorer"

### What I did

- Created `internal/services/corpus/service.go` with three methods:
  - `SourceSummaries()`: per-source counts with optional embedding coverage for a given strategy/provider/model/dimensions identity.
  - `DocumentBrowser()`: paginated documents filtered by source, with chunk and embedding counts.
  - `DocumentDetail()`: full document with metadata, optional content text, and chunk list with embedding presence.
- Registered three HTTP endpoints in `internal/api/handlers.go`:
  - `GET /api/v1/corpus/sources`
  - `GET /api/v1/corpus/documents`
  - `GET /api/v1/corpus/documents/{id}`
- Added `parseEmbeddingIdentity()` and `intQueryDefault()` helpers.

### Why

The design guide specified that the Corpus Explorer needs its own corpus-specific endpoints rather than overloading the generic document endpoints. A separate corpus service keeps query logic out of handlers and leaves room for a future CLI group.

### What worked

- Build and all existing tests passed on first compile.
- The API returned correct source summaries with embedding coverage for the OpenAI smoke data (35 embeddings across 481 chunks).

### What didn't work

- First attempt had the SQL parameter ordering wrong in `DocumentBrowser()`: `sourceID` was placed at the start of the args array but the first `?` in the generated SQL was for the LEFT JOIN subquery's `strategy_id`. This caused the document browser to return null results.
- Same bug in `DocumentDetail()` chunk query: `documentID` and `strategyID` were placed before the chunk_embeddings LEFT JOIN params.

### What I learned

- When building dynamic SQL with optional LEFT JOINs that use `?` placeholders, the args must follow the order of `?` appearance in the final SQL string, not the order of code blocks. I now build args incrementally as each block is appended.

### What was tricky to build

- The conditional SQL construction means the number and order of `?` placeholders changes depending on which identity fields are provided. The scan logic also varies (different column counts). Both must stay in sync.
- SQLite NULL handling: using `sql.NullInt64` for optional chunk/embedding counts and `sql.NullString` for optional text fields.

### What warrants a second pair of eyes

- The three-way branching in `SourceSummaries()` scan logic (identity with/without embedding, strategy-only, no strategy) is fragile. If a new identity field is added, all branches need updating.
- The `DocumentDetail()` method does N+1 (one query for document, one for chunks) rather than a single join. This is acceptable for single-document inspection but not for bulk operations.

### What should be done in the future

- Add corpus service tests with temporary SQLite databases.
- Consider a unified scan struct approach to avoid the branching scan logic.

### Code review instructions

- Start with `internal/services/corpus/service.go` — check the SQL query construction and parameter ordering.
- Then `internal/api/handlers.go` — check the corpus handler functions and the `parseEmbeddingIdentity` helper.

### Technical details

- Commit: d0934ea
- Endpoints: `GET /api/v1/corpus/sources`, `GET /api/v1/corpus/documents`, `GET /api/v1/corpus/documents/{id}`
- Bugfix commit: 7934e0a — fixed SQL parameter ordering in document browser and detail queries.

---

## Step 4: Frontend Corpus Explorer and Monochrome Theme Redesign

I redesigned the entire site to a retro macOS 1 monochrome aesthetic (pure black and white, no menu bar, no window chrome, color accents only on foreground text) and built the full Corpus Explorer view with source panel, document browser, document inspector with four tabs, chunk timeline bar, and embedding coverage strip.

### Prompt Context

**User prompt (verbatim):** "WOrk on the docmgr ticket RAGEVAL-003, read careefully, and also feel free to update the existing site if you think it's worth it, you are a great designer. Use a retro macos1 monochrome-ish look (you can use color accents on font foreground color), no menu bar, no window chrome. Add detailed task list to the ticket, then work on them step by step, keeping a diary as you go, committing at appropriate intervals"

**Assistant interpretation:** Implement the full Corpus Explorer from the design guide, redesign the site to retro macOS 1 monochrome, and maintain a diary.

**Inferred user intent:** Turn the design plan into a working, visually polished product that uses a clean monochrome aesthetic.

**Commit (code):** 96847da — "feat(corpus): add Corpus Explorer UI with retro macOS 1 monochrome theme"

### What I did

- Rewrote `web/src/index.css` as a complete monochrome design system:
  - Pure black/white palette with subtle grays.
  - No window chrome, no menu bar — replaced with a minimal nav strip.
  - Color accents only on foreground text: blue (accent), green (accent-green), amber (accent-amber), red (accent-red).
  - New classes: `.panel`, `.panel-header`, `.data-table`, `.coverage-strip`, `.chunk-bar`, `.tab-bar`, `.fieldset`, `.stat-grid`, `.meta-grid`.
- Added RTK Query corpus types and endpoints to `web/src/services/api.ts`:
  - `CorpusSourceSummary`, `CorpusDocumentRow`, `CorpusChunk`, `CorpusDocumentDetail` types.
  - `listCorpusSources`, `listCorpusDocuments`, `getCorpusDocument` query endpoints.
  - Helper `filterIdentityParams()` to conditionally include identity params.
- Built `web/src/components/corpus/CorpusExplorerView.tsx` as the main view:
  - Identity bar at top (strategy, provider, model, dimensions).
  - Three-column layout: source list | document browser | document inspector.
  - Source items show document count, word count, embedding coverage percentage.
  - Document table shows title, word count, chunk count, embedding ratio, status.
  - Document inspector with four tabs: Overview, Text, Chunks, Coverage.
  - Chunk timeline bar visualizes chunk positions with embedded/missing shading.
  - Chunk table with index, range, tokens, embedding status (●/○), copyable chunk ID.
  - Selected chunk shows full text below the table.
  - Coverage strip with per-chunk dot visualization.
  - Missing chunks table.
- Updated `web/src/App.tsx`: removed MacMenuBar, added Corpus as default view, new nav strip.
- Updated all existing views (Pipeline, Embeddings, Search, Evaluation) to use new panel-based classes.

### Why

The design guide specified a Corpus Explorer centered on user intents (learn, inspect, validate). The user requested a retro macOS 1 monochrome look with no menu bar and no window chrome. This implementation delivers both.

### What worked

- The three-column explorer layout renders cleanly and all data flows correctly through RTK Query.
- Selecting a source loads 483 TTC dump articles in the document browser with correct chunk/embedding counts.
- Selecting "Crape Myrtle Varieties and Guide" shows 55 chunks with 10 embedded.
- The chunk timeline bar correctly visualizes the first 10 chunks as embedded (filled black) and the rest as missing (light gray).
- The coverage strip shows the per-chunk dot visualization clearly.
- The monochrome theme is clean and readable — black headers, white panels, dotted row separators.
- TypeScript compiles cleanly with no errors.
- All Go tests pass.

### What didn't work

- Initial TypeScript build failed with unused imports (`CorpusDocumentRow`) and unused parameter (`idx` in coverage strip map). Quick fix by removing the unused import and parameter.

### What I learned

- The retro macOS 1 aesthetic works well for data-dense tooling: the monochrome palette with black panel headers creates strong visual hierarchy without visual noise.
- Using foreground-only color accents (green for embedded, amber for partial, red for error) keeps the monochrome feel while still communicating status.
- The chunk timeline bar is an effective way to show both chunk position and embedding coverage at a glance.

### What was tricky to build

- Balancing monochrome strictness with usability: the user said "you can use color accents on font foreground color" which I interpreted narrowly — only text `color` property uses non-black/white values, never backgrounds or borders.
- The three-column layout with independent scrolling per panel requires fixed heights and `overflow-y: auto` on the panel bodies.
- The conditional SQL query building in the corpus service required careful parameter ordering (see Step 3).

### What warrants a second pair of eyes

- The `CorpusExplorerView.tsx` file is large (500+ lines). It could be split into the component files suggested by the design guide (SourceSummaryPanel, DocumentBrowser, DocumentInspector, ChunkTimeline, EmbeddingCoverageStrip) for better maintainability.
- The default active tab in the document inspector is "chunks" — this works well for the primary user intent but might confuse users who expect an overview first. Currently set to "overview".

### What should be done in the future

- Split `CorpusExplorerView.tsx` into separate component files.
- Add pagination to the document browser (currently limited to 100 docs, TTC dump articles has 483).
- Add text search/filter for the document browser.
- Add keyboard navigation support.
- Add corpus service unit tests.

### Code review instructions

- Start with `web/src/components/corpus/CorpusExplorerView.tsx` — the main component.
- Then `web/src/index.css` — the full theme redesign.
- Check `web/src/services/api.ts` for the new corpus types and endpoints.
- Verify `web/src/App.tsx` has the new nav strip and Corpus as default view.

### Technical details

- Commit: 96847da
- Bugfix commit: 7934e0a (SQL parameter ordering)
- New files: `internal/services/corpus/service.go`, `web/src/components/corpus/CorpusExplorerView.tsx`
- Modified files: `internal/api/handlers.go`, `web/src/index.css`, `web/src/services/api.ts`, `web/src/App.tsx`, `web/src/components/pipeline/PipelineView.tsx`, `web/src/components/embeddings/EmbeddingsView.tsx`, `web/src/components/search/SearchView.tsx`, `web/src/components/evaluation/EvaluationView.tsx`

---

## Step 5: Component Split, Search/Pagination, and Corpus Service Tests

I split the monolithic CorpusExplorerView into separate component files, added a search filter and pagination to the document browser, and wrote 8 unit tests for the corpus service.

### Prompt Context

**User prompt (verbatim):** "Go ahead."

**Assistant interpretation:** Continue with follow-up improvements: split components, add search/pagination, write corpus service tests.

**Inferred user intent:** Polish and harden the Corpus Explorer implementation.

**Commit (code):** 32f7d71 — "refactor(corpus): split CorpusExplorerView into components, add search/pagination and tests"

### What I did

- Split `CorpusExplorerView.tsx` (500+ lines) into 6 files:
  - `CorpusExplorerView.tsx` — main layout, state, data flow (~130 lines)
  - `IdentityBar.tsx` — embedding identity controls
  - `SourcePanel.tsx` — source list with coverage items
  - `DocumentBrowser.tsx` — document table with search and pagination
  - `DocumentInspector.tsx` — 4-tab inspector (overview/text/chunks/coverage)
  - `ChunkTimelineBar.tsx` — proportional chunk position visualization
- Added search filter in DocumentBrowser: filters by title or ID substring
- Added pagination: loads 100 docs per page with "Load more (N remaining)" button
- Fixed coverage tab to use actual identity props instead of hardcoded label
- Wrote 8 corpus service unit tests in `internal/services/corpus/service_test.go`:
  - `TestSourceSummaries_NoIdentity`
  - `TestSourceSummaries_WithIdentity`
  - `TestDocumentBrowser`
  - `TestDocumentBrowser_Pagination`
  - `TestDocumentDetail`
  - `TestDocumentDetail_NotFound`
  - `TestDocumentDetail_NoText`
  - `TestDocumentDetail_NoStrategy`

### Why

The monolithic component was hard to maintain. Separate files make each piece independently improvable. Search and pagination are essential for the TTC corpus (483 articles). Tests prevent regressions in the SQL query construction.

### What worked

- All 8 tests pass on first run.
- TypeScript compiles cleanly after the split.
- Search filters 100 docs to matching results instantly (client-side filter on loaded data).
- Pagination correctly loads 100 → 200 docs on "Load more" click.
- The Go server needs a rebuild after frontend changes since it embeds the SPA.

### What didn't work

- Initially forgot to rebuild the Go binary after the frontend refactor. The Playwright browser was still seeing the old embedded frontend. Fixed by rebuilding `go build ./cmd/rag-eval` and restarting the server.

### What I learned

- The Go binary embeds the frontend at build time via `go:embed`. After any `web/` change, the full build sequence is: `cd web && npm run build && cd .. && go build ./cmd/rag-eval`.
- Client-side search is fast enough for up to 1000 loaded documents. Server-side search would only be needed for much larger corpora.

### What was tricky to build

- The pagination approach uses a growing `pageLimit` state variable rather than traditional offset-based pagination. This avoids re-fetching already-loaded data and simplifies the RTK Query cache, but means the query parameter changes (limit grows) which triggers a refetch of all data.

### What warrants a second pair of eyes

- The search filter is purely client-side on the loaded page. If a user searches for something in the first 100 docs that doesn't match, but the matching doc is at position 150, they won't find it until they load more. This is a known tradeoff.

### What should be done in the future

- Add server-side search endpoint for very large sources (2500+ docs).
- Add keyboard navigation (arrow keys in source list and document table).
- Consider virtual scrolling for very large document lists.

### Code review instructions

- Review the component split: each file should have a single responsibility.
- Check the test seed data covers the key scenarios (embedded/not embedded, with/without strategy, pagination).
- Verify the search filter works in the browser.

### Technical details

- Commit: 32f7d71
- New files: `IdentityBar.tsx`, `SourcePanel.tsx`, `DocumentBrowser.tsx`, `DocumentInspector.tsx`, `ChunkTimelineBar.tsx`, `service_test.go`
- Rewritten: `CorpusExplorerView.tsx` (from 500+ to ~130 lines)
