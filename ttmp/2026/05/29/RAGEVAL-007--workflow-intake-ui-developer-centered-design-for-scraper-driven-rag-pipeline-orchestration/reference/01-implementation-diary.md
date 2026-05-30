---
title: "Implementation diary"
doc-type: reference
ticket: RAGEVAL-007
status: active
intent: long-term
topics: [rag, ui, workflow, scraper, intake, design]
created: 2026-05-29
---

# Diary

## Goal

Capture the investigation and design work for RAGEVAL-007: creating a workflow intake UI for the scraper-driven RAG pipeline.

## Step 1: Codebase investigation and ticket creation

Explored the full rag-evaluation-system codebase to understand the architecture: the two-database model (corpus DB for domain data, engine DB for workflow orchestration), the scraper engine library, the IntakeRunner, the submit/DAG logic, and the existing React SPA. Then created the RAGEVAL-007 ticket with a design doc and diary.

### Prompt Context

**User prompt (verbatim):** "Create a new docmgr ticket to create a great UI for the workflow scraper approach to rag embedding in 2026-05-27--rag-evaluation-system/ using the scraper/ framework. We already have a great UI, but now we want to start with a proper design workflow. Let's analyze the needs of a developer working on the intake workflows and parts of it, and what affordances they need. Then derive views and screens to enable them to do this work and to access this data to help them iterate. Use ascii screenshots and sketch out input props needed and which / where the data comes from. Create a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable. Kepe a diary as you work."

**Assistant interpretation:** Create a docmgr ticket for designing a workflow intake UI. Analyze developer needs, derive screens with ASCII wireframes and input props, write a comprehensive intern-ready design/implementation guide, store in ticket, upload to reMarkable, and keep a diary.

**Inferred user intent:** Produce a complete design specification that a new intern can pick up and implement the workflow UI end-to-end, including backend API, frontend components, and data flow diagrams.

### What I did

- Explored the full directory structure of the rag-evaluation-system workspace
- Read key Go source files: workflow engine.go, intake_runner.go, ops.go, submit.go, constants.go, api/handlers.go
- Read scraper engine model types, store interface, scheduler, engineview services
- Read React SPA: App.tsx, api.ts, store/index.ts, CorpusExplorerView.tsx, PipelineView.tsx
- Read corpus DB queries and database schema
- Examined existing ttmp tickets (RAGEVAL-001 through RAGEVAL-006) for context
- Created RAGEVAL-007 ticket via docmgr
- Created design doc: `01-intake-workflow-ui-design-screens-affordances-and-data-flow.md`
- Created this diary document

### Why

The existing UI covers corpus exploration and search but has zero visibility into the workflow engine. A developer iterating on intake pipelines needs to see workflow progress, diagnose failures, and submit new runs—all from the browser.

### What worked

- The `engineview` service package already provides all the read/mutation operations needed (ListWorkflows, WorkflowOps, RetryOp, CancelWorkflow, ListQueues). No new backend logic needed—just HTTP handlers that call existing service methods.
- The existing React SPA has a clear pattern (RTK Query + Bootstrap-styled panels) that the new Workflows view can follow directly.
- The intake DAG structure is predictable (fan-out chunk → fan-in embed/bm25), so the DAG visualization can use a simple hierarchical layout without needing a graph library.

### What didn't work

- Nothing failed during investigation. The codebase is well-structured and the gap is clear.

### What I learned

- The two-database architecture is the key insight: the engine DB is invisible in the current UI, which is why developers must use CLI commands to monitor workflows.
- The `IntakeOpInput` struct has many optional fields that vary by operation type—the UI needs to display different result schemas for chunk vs embed vs bm25 ops.
- The scheduler's `Observer` interface could power real-time updates in the future, but polling is sufficient for phase 1.

### What was tricky to build

- The design doc needed to be comprehensive enough for an intern while staying grounded in actual code. I anchored every screen and API endpoint to specific Go functions and SQL tables, with line references to the source files.

### What warrants a second pair of eyes

- The DAG visualization approach: I chose CSS-based absolute positioning with SVG edges. For 500+ ops this might need virtualization or a different approach. The collapsed view design needs iteration.
- The Submit Intake modal: the form has many fields. The grouping and progressive disclosure need UX review.
- Engine DB concurrency: SQLite WAL mode should handle concurrent reads + writes, but this needs verification.

### What should be done in the future

- SSE/WebSocket for real-time op status updates (instead of polling)
- Workflow templates (save common intake configurations)
- Workflow comparison (side-by-side run diffs)
- Bulk retry for all failed ops in a workflow
- Op artifact viewer (for preprocess_document artifacts)

### Code review instructions

- Start with `ttmp/2026/05/29/RAGEVAL-007--.../design/01-intake-workflow-ui-design-screens-affordances-and-data-flow.md`
- Cross-reference the scraper engineview services against the proposed API endpoints
- Verify the RTK Query type definitions match the Go struct field names and JSON tags

### Technical details

- Engine DB path: `state/rag-eval-workflows.db` (configurable via `--engine-db`)
- Corpus DB path: `data/rag-eval.db` (configurable via `--db-path`)
- Workflow ID format: `intake-YYYYMMDDTHHMMSS`
- Strategy ID format: `{strategy}-{chunkSize}-{overlap}` (e.g. `fixed-1200-150`)
- Op ID format: `{workflowId}:{operation}:{documentId}` or `{workflowId}:{operation}`

## Step 3: Backend API implementation (Phase 1)

Added 5 new endpoints to `workflow_artifact_handlers.go` (extending the 3 existing ones):
- `GET /api/v1/workflows/{id}/results/{opId}` — inspect op result data
- `POST /api/v1/workflows/{id}/retry/{opId}` — retry a failed op
- `POST /api/v1/workflows/{id}/cancel` — cancel a running workflow
- `POST /api/v1/workflows/intake` — submit new intake workflow from browser
- `GET /api/v1/queues` — queue health status

All wired through the existing `handler.engineDB` field. No new handler struct needed — just added methods to the existing `handler` type.

Committed as `faf6ab`.

## Step 4: Frontend implementation (Phases 2-6)

Created `WorkflowsView.tsx` with 5 sub-components:
1. **WorkflowsList** — table with status filter, polling every 2s
2. **QueueHealthWidget** — queue status table, polling every 5s
3. **SubmitIntakeModal** — form with source IDs, chunking params, embedding config, BM25 toggle
4. **WorkflowDetail** — header with stats, compact op graph, ops-by-group table, op inspector
5. **OpGroupRow** — clickable row per (operation, status) group

Added RTK Query endpoints and TypeScript types to `api.ts`.
Added "Workflows" tab to `App.tsx` nav.

**Critical fix:** Changed ops endpoint to return grouped summary (by operation+status) instead of full op list. A workflow with 3117 docs creates 6236+ ops — returning all would explode context. The grouped response returns ~4 rows instead of 6236, with a `sample` op per group for drill-down.

Committed as `387713d`.

### What worked
- The grouped ops approach keeps the response tiny even for huge workflows
- The op graph shows chunk fan-out + embed/bm25 fan-in clearly
- The inspector shows sample op input fields for drill-down
- Submit intake modal works end-to-end (POST → workflow created → navigate to detail)

### What didn't work
- First attempt used raw op list — 6236 ops in one response crashed the browser context
- TypeScript field names had to match actual Go JSON tags (PascalCase for nested op fields, camelCase for wrapper fields)

### What was tricky to build
- The Go JSON serialization uses different casing conventions for nested structs vs top-level: `WorkflowRun` fields serialize as PascalCase (because Go's `encoding/json` uses exported field names by default), while `engineview.WorkflowOp` wrapper fields use `json:"status"` tags with camelCase. Had to inspect actual API responses to align TypeScript types.
- The `QueueStatus` struct in the scraper uses `json:"inFlight"` etc (camelCase), not PascalCase.

### What warrants a second pair of eyes
- The `handleWorkflowOps` grouping logic — it creates `opGroup` structs inline and does a map+order pattern. Should be tested with edge cases (empty ops, single op, mixed statuses).
- The `handleSubmitIntake` handler applies defaults that overlap with `SubmitIntakeWorkflow()`'s own defaults. Double-defaulting could cause issues if they diverge.

### What should be done in the future
- Add a dedicated endpoint to fetch ops for a specific group (e.g., all failed chunk ops) with pagination
- Add op result detail view (GET /api/v1/workflows/{id}/results/{opId})
- Real-time updates via SSE instead of polling
- Workflow templates for common configurations

## Step 2: Doc validation and reMarkable upload

Ran `docmgr doctor`, fixed vocabulary warnings by adding `design`, `intake`, `scraper`, `ui` topics. Removed 3 broken related file entries for scraper engineview files (paths outside project root don't resolve correctly in docmgr). Doctor passes cleanly after fixes.

Uploaded bundle to reMarkable at `/ai/2026/05/29/RAGEVAL-007` — confirmed with `remarquee cloud ls`.

No code was implemented. The deliverable is the design document only.

## Step 5: Phase 7 Polish

Polished the Workflows UI with proper retro Mac CSS design system classes and fixed a critical React crash.

### Prompt Context
**User prompt (verbatim):** (see Step 1)
**Assistant interpretation:** Continue implementation — Phase 7 polish.
**Inferred user intent:** Make the UI look production-quality with consistent styling.
**Commit (code):** 1851167 — "feat(ui): Phase 7 polish — progress bar, proper CSS classes, styled modal, op inspector fix"

### What I did
- Added CSS classes: `status-running`, `status-error`, `status-canceled` (were referenced but never defined)
- Added `progress-bar` / `progress-fill` CSS with done/running/failed color bands and centered label
- Added `op-graph`, `op-node`, `op-arrow`, `op-fan-out` classes for the DAG visualization
- Added `modal-overlay`, `modal-panel`, `form-row`, `form-section` classes for the submit modal
- Rewrote `WorkflowsView.tsx` to use all `.btn`, `.input`, `.select`, `.form-section`, `.form-row` classes
- Added `ProgressBar` component in workflow detail header
- Added `OpGraph` component as a standalone sub-component with proper fan-out/fan-in layout
- Added selectable row highlighting (`.selectable` → `.selected` on click)
- Added error display in submit modal
- Added available source names in modal
- Fixed React #310 crash caused by `inspectedSample` in `useCallback` dependency array

### Why
The raw HTML elements (bare `<button>`, `<input>`, `<fieldset>`) didn't match the retro Mac design system. Missing CSS classes caused unstyled or broken status colors. The progress bar gives an at-a-glance view of how far a workflow has progressed.

### What worked
- Progress bar with done/running/failed bands is immediately readable
- The `.form-row` pattern (label right-aligned + input) is much cleaner than stacked labels
- Using `.selectable` rows with hover highlight makes the table interactive
- Source names list in modal helps discoverability

### What didn't work
- **React #310 crash**: `inspectedSample` (derived from polled `groups` array) was in `useCallback` deps, causing a new function on every poll (every 2s). Each new function triggered a re-render, which changed `inspectedSample` reference, which created another new function → infinite loop. Fixed by using a `useRef` to hold the current sample and removing it from deps.
- Variable ordering error: had `selectedGroup`/`inspectedSample` before `groups` was declared, causing TypeScript block-scoped variable errors.

### What was tricky to build
- The React #310 debug cycle: minified React errors give no component name or line reference. Had to reason about the dependency chain: `useCallback([inspectedSample])` → `inspectedSample` from `groups.find(...)` → `groups` from `useGetWorkflowOpsQuery` polling → new `groups` array every 2s → new `inspectedSample` → new `handleRetry` → re-render. The fix was to break the dependency by using `useRef`.

### What warrants a second pair of eyes
- The `inspectedSampleRef` pattern: assigning to `.current` during render is technically a side effect. It works because it's idempotent (same value), but a `useEffect` would be more correct. Low priority since it's just a cache.

### What should be done in the future
- Replace polling with SSE for real-time updates
- Add per-group drill-down (click a group → paginated list of individual ops)
- Add op result detail endpoint integration
- Workflow templates

### Code review instructions
- Start at `web/src/components/workflows/WorkflowsView.tsx` — review the `inspectedSampleRef` pattern and `ProgressBar` component
- Check `web/src/index.css` for the new CSS additions under `/* ===== Progress Bar ===== */` and `/* ===== Op Graph ===== */`
- Run: `cd web && pnpm build && cd .. && go build ./cmd/rag-eval`
- Verify: open http://127.0.0.1:8772 → Workflows tab → click workflow → check progress bar + op graph + inspector

## Step 6: devctl plugin with dynamic port allocation

Created a devctl NDJSON v2 plugin that replaces the inline `.devctl.yaml` services with a proper orchestrated pipeline.

### Prompt Context
**User prompt (verbatim):** "create a devctl setup to run the vite + backend setup"
**Assistant interpretation:** Create a devctl plugin that orchestrates the Go backend + Vite dev server.
**Inferred user intent:** Replace manual tmux-based startup with devctl-managed services.
**Commit (code):** 3216dc2 — "feat(devctl): add plugin with dynamic port allocation"

Follow-up prompt: "find a free port for each process as part of the planning."

### What I did
1. Created `plugins/rag-eval.py` — a Python NDJSON v2 devctl plugin with:
   - `config.mutate`: finds free ports using `socket.bind()` (prefers 8772/5173, falls back to OS-allocated)
   - `validate.run`: checks for go, pnpm, node_modules, state dir
   - `build.run`: runs `pnpm install` (if needed) → `pnpm build` → `go build` in order
   - `launch.plan`: launches backend with health check + Vite dev server with `--port` flag
   - `command.run`: exposes `build` and `web-build` dynamic commands
2. Rewrote `.devctl.yaml` from inline services to plugin reference
3. Updated `web/vite.config.ts` to read `RAG_EVAL_BACKEND_PORT` env var for proxy target

### Why
The existing `.devctl.yaml` used hardcoded ports and inline service definitions. When ports were occupied (e.g. stale processes), Vite would silently pick a different port, breaking the proxy. The plugin approach lets `config.mutate` find free ports and propagate them through the entire pipeline.

### What worked
- Dynamic port allocation: when 8772/5173 are occupied, falls back to OS-assigned ports and propagates them correctly
- The merged config flows from `config.mutate` → `launch.plan` via `input.config` with nested structure
- Vite reads `RAG_EVAL_BACKEND_PORT` env to set proxy target dynamically
- `pnpm dev --port N` sets the Vite port explicitly (no more silent fallback)

### What didn't work
- First attempt used `config.get("services.backend.port")` which treats the dotted key as a single dict key — but the config is nested (`{"services": {"backend": {"port": N}}}`). Fixed with a `cfg_get()` helper that traverses dot paths.
- `artifacts` field in `build.run` was initially an array but devctl expects `map[string]string`. Fixed to `{"rag-eval-bin": "rag-eval", ...}`.

### What was tricky to build
- The `find_free_port()` function has a TOCTOU window: the port is free when checked but might be taken by the time the service starts. This is inherent to port allocation. The function tries the preferred port first, then asks the OS for any free port via `bind(0)`.
- The devctl config merge passes nested objects, not flat dotted-key maps. The `input.config` in `launch.plan` has `{"services": {"backend": {"port": N}}}` not `{"services.backend.port": N}`. Required a recursive key traversal helper.

### What warrants a second pair of eyes
- The TOCTOU window in port allocation — in rare cases both the health checker and the actual backend could race. Not a real issue since devctl starts services sequentially.
- The `find_free_port` function uses `SO_REUSEADDR` which on Linux allows binding to a port in TIME_WAIT. On macOS this can cause issues.

### What should be done in the future
- Add `prepare.run` op for running DB migrations
- Add profile support (e.g. `devctl up -p prod` uses embedded SPA only, no Vite)
- Persist allocated ports in `.devctl/state.json` for reattach scenarios

### Code review instructions
- Start at `plugins/rag-eval.py` — review `find_free_port()`, `cfg_get()`, and the `launch.plan` service definitions
- Check `web/vite.config.ts` — review `RAG_EVAL_BACKEND_PORT` env var reading
- Run: `devctl plan && devctl up && devctl status && devctl down`
- Verify: `curl -s http://127.0.0.1:5173/api/v1/health` returns OK
