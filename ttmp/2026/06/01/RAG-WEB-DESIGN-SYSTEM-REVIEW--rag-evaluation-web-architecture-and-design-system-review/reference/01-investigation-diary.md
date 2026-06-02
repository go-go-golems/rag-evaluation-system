---
Title: Investigation Diary
Ticket: RAG-WEB-DESIGN-SYSTEM-REVIEW
Status: active
Topics:
    - rag
    - react
    - design-system
    - css
    - dmeta
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: dmeta-ir/README.md
      Note: Phase 3 diary evidence for catalog boundary
    - Path: dmeta-ir/core-model/archetypes/10-search-workbench.yaml
      Note: Phase 3 core-model diary evidence
    - Path: dmeta-ir/interactions/representations/10-search-workbench.yaml
      Note: Phase 3 Interaction IR diary evidence
    - Path: dmeta-ir/meta-design-systems/web/widgets/search-workbench.yaml
      Note: Phase 3 Web MDS diary evidence
    - Path: ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/design-doc/01-rag-evaluation-web-architecture-and-design-system-review.md
      Note: Primary report produced by the investigation
    - Path: web/package.json
      Note: Phase 0/2 tooling changes
    - Path: web/src/App.tsx
      Note: |-
        App shell and navigation event evidence
        Diary evidence for AppShell adoption
    - Path: web/src/components/corpus/DocumentBrowser.stories.tsx
      Note: Diary evidence for Corpus Storybook coverage
    - Path: web/src/components/corpus/DocumentBrowser.tsx
      Note: Diary evidence for Corpus document table adoption
    - Path: web/src/components/corpus/DocumentInspector.stories.tsx
      Note: |-
        Diary evidence for Corpus Storybook coverage
        Diary evidence for Redux Provider story fix
    - Path: web/src/components/corpus/DocumentInspector.tsx
      Note: Diary evidence for Corpus inspector primitive adoption
    - Path: web/src/components/corpus/DocumentInspector/DocumentInspector.stories.tsx
      Note: Moved Corpus story remains co-located with widget
    - Path: web/src/components/corpus/SourcePanel.stories.tsx
      Note: Diary evidence for Corpus Storybook coverage
    - Path: web/src/components/corpus/SourcePanel.tsx
      Note: |-
        Diary evidence for Corpus source design-system adoption
        Diary evidence for SourcePanel move-back
    - Path: web/src/components/embeddings/EmbeddingsView.module.css
      Note: Diary evidence for embeddings view CSS ownership
    - Path: web/src/components/embeddings/EmbeddingsView.tsx
      Note: Diary evidence for embeddings view primitive adoption
    - Path: web/src/components/foundation/Caption/Caption.tsx
      Note: Diary evidence for expanded primitive extraction
    - Path: web/src/components/foundation/index.ts
      Note: Phase 1 foundation primitive barrel
    - Path: web/src/components/layout/TabList/TabList.tsx
      Note: Diary evidence for expanded layout primitives
    - Path: web/src/components/layout/index.ts
      Note: Phase 1 layout primitive barrel
    - Path: web/src/components/molecules/CoveragePanel/CoveragePanel.tsx
      Note: |-
        Phase 2 extraction diary evidence
        Diary evidence for support molecule primitive adoption
    - Path: web/src/components/molecules/DataTable/DataTable.tsx
      Note: Diary evidence for reusable molecule extraction
    - Path: web/src/components/molecules/MetadataGrid/MetadataGrid.tsx
      Note: Diary evidence for reusable molecule extraction
    - Path: web/src/components/molecules/QueryPresetList/QueryPresetList.tsx
      Note: |-
        Phase 2 extraction diary evidence
        Diary evidence for support molecule primitive adoption
    - Path: web/src/components/organisms/QueueHealthPanel/QueueHealthPanel.tsx
      Note: Diary evidence for workflow organism extraction
    - Path: web/src/components/organisms/ResultInspectorPanel/ResultInspectorPanel.module.css
      Note: Diary evidence for remaining inspector anatomy CSS
    - Path: web/src/components/organisms/ResultInspectorPanel/ResultInspectorPanel.tsx
      Note: |-
        Phase 2 inspector extraction diary evidence
        Diary evidence for MetadataGrid/TabList adoption
    - Path: web/src/components/organisms/RetrievalResultsPanel/RetrievalResultsPanel.tsx
      Note: |-
        Phase 2 extraction diary evidence
        Diary evidence for DataTable adoption
    - Path: web/src/components/organisms/SearchControlsPanel/SearchControlsPanel.tsx
      Note: |-
        Phase 2 extraction diary evidence
        Diary evidence for controls primitive adoption
    - Path: web/src/components/organisms/WorkflowListPanel/WorkflowListPanel.tsx
      Note: Diary evidence for workflow organism extraction
    - Path: web/src/components/organisms/WorkflowOpGraphPanel/WorkflowOpGraphPanel.tsx
      Note: Diary evidence for workflow detail extraction
    - Path: web/src/components/organisms/WorkflowOpGroupsPanel/WorkflowOpGroupsPanel.tsx
      Note: Diary evidence for workflow detail extraction
    - Path: web/src/components/organisms/WorkflowOpInspectorPanel/WorkflowOpInspectorPanel.tsx
      Note: Diary evidence for sampled-op inspector extraction
    - Path: web/src/components/organisms/WorkflowOpResultPanel/WorkflowOpResultPanel.tsx
      Note: Diary evidence for workflow op result extraction
    - Path: web/src/components/organisms/WorkflowSummaryPanel/WorkflowSummaryPanel.tsx
      Note: Diary evidence for workflow detail extraction
    - Path: web/src/components/pages/PipelinePage/PipelinePage.stories.tsx
      Note: Diary evidence for required page story coverage
    - Path: web/src/components/pages/SearchWorkbenchPage/SearchWorkbenchPage.tsx
      Note: Diary evidence for page-boundary extraction
    - Path: web/src/components/pipeline/PipelineOverview.stories.tsx
      Note: Diary evidence for required component story coverage
    - Path: web/src/components/search/SearchView.module.css
      Note: Phase 2 layout primitive migration diary evidence
    - Path: web/src/components/search/SearchView.tsx
      Note: |-
        First recommended vertical-slice evidence
        Phase 2 Search Workbench composition update
    - Path: web/src/components/workflows/WorkflowsView.tsx
      Note: |-
        Diary evidence for sampled-op inspector composition
        Diary evidence for op result container/presentation split
        Diary evidence for intake modal primitive adoption
    - Path: web/src/index.css
      Note: CSS/design-system evidence
    - Path: web/src/services/api.ts
      Note: API contract evidence
    - Path: web/src/styles/tokens.css
      Note: Phase 1 token extraction
ExternalSources: []
Summary: Chronological diary for the RAG Evaluation web architecture and design-system/DMETA review.
LastUpdated: 2026-06-01T00:00:00-04:00
WhatFor: Use this diary to understand how the RAG web review was investigated, what commands were run, what failed, and how to continue.
WhenToUse: Read before continuing the RAG web design-system or DMETA IR implementation work.
---






















# Diary

## Goal

Capture the investigation and delivery steps for a detailed intern-facing review of the `2026-05-27--rag-evaluation-system` web application, with recommendations based on the design-system and DMETA IR playbook.

## Step 1: Create the Ticket and Map the Web System

I created a new docmgr ticket inside the RAG Evaluation System repository using the repository-local `ttmp` root, then inspected the React app, CSS, RTK Query API layer, Go API handlers, embedded-SPA server, and prior RAGEVAL design documents. The review treats the current React/CSS implementation as evidence and evaluates it against the layered design-system + DMETA approach from the previous TTC documents.

The investigation showed that the web app is already a useful developer dashboard: it has search, corpus exploration, workflows, embeddings, pipeline overview, and a placeholder evaluation view. The main issue is not missing product value; it is missing architectural ownership. Global CSS, inline styles, large page files, and implicit domain concepts need to be turned into tokens, primitives, component layers, DMETA IR, and validators.

### Prompt Context

**User prompt (verbatim):** "ok, let's go through 2026-05-27--rag-evaluation-system and create a new ticket in that 2026-05-27--rag-evaluation-system/ttmp (using docmgr --root ...) and study the whole web stuff in there and make a full report, according to these review and organization document we just wrote.

Create  a detailed analysis / design / review that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Create a new RAG-local docmgr ticket, audit the whole web/UI/API surface, write a detailed technical report for an intern using the design-system/DMETA guidance just produced, and upload it to reMarkable.

**Inferred user intent:** The user wants a durable onboarding and architecture report that explains the current RAG web system and gives a concrete path to formalize it into a design system plus DMETA IR.

**Commit (code):** N/A — documentation/review step.

### What I did

- Ran `docmgr --root ttmp` from inside `2026-05-27--rag-evaluation-system` to create ticket `RAG-WEB-DESIGN-SYSTEM-REVIEW`.
- Created:
  - `design-doc/01-rag-evaluation-web-architecture-and-design-system-review.md`
  - `reference/01-investigation-diary.md`
- Inspected frontend files under `web/src`, including:
  - `App.tsx`
  - `services/api.ts`
  - `store/index.ts`
  - `index.css`
  - `SearchView.tsx`
  - `CorpusExplorerView.tsx`
  - `DocumentInspector.tsx`
  - `WorkflowsView.tsx`
  - `EmbeddingsView.tsx`
  - `PipelineView.tsx`
  - `EvaluationView.tsx`
  - `components/retro/*`
- Inspected backend/server files needed to understand the web system:
  - `cmd/rag-eval/cmds/serve/server.go`
  - `cmd/rag-eval/cmds/serve/root.go`
  - `internal/web/spa.go`
  - `internal/api/handlers.go`
  - `internal/api/workflow_artifact_handlers.go`
  - `internal/db/db.go`
- Inspected prior ticket documents for context:
  - RAGEVAL-003 corpus explorer design
  - RAGEVAL-004 retrieval UX/status design
  - RAGEVAL-007 workflow UI design
- Ran static scans for file sizes, global class usage, inline style counts, custom events, and absence of `.module.css`, `.stories.tsx`, and `dmeta-ir` files.
- Attempted a frontend build.

### Why

- The report needs to be evidence-backed and useful for a new intern.
- The RAG app is a complex dashboard candidate, so the review must explain product concepts, UI views, API data flow, CSS/design-system state, and future DMETA IR boundaries.
- The previous TTC playbook says to begin with the current UI as evidence, not as architecture; this investigation follows that method.

### What worked

- The repository already has a local `ttmp` and existing RAGEVAL ticket history.
- The React source tree is compact enough to audit completely.
- The backend API route registry and RTK Query layer line up well, making data flow understandable.
- Prior RAGEVAL docs provide useful product context for corpus, retrieval, and workflow UX.

### What didn't work

- `pnpm build` failed because frontend dependencies were not installed in `web/node_modules`:

```text
> web@0.0.0 build /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web
> tsc && vite build

sh: 1: tsc: not found
 ELIFECYCLE  Command failed.
 WARN  Local package.json exists, but node_modules missing, did you mean to install?
```

- I did not run `pnpm install` because the task is a documentation/review pass and the failure itself is useful validation evidence for the report.

### What I learned

- The web app is not merely a placeholder. It already exposes real RAG workflows: BM25/vector/hybrid search, corpus inspection, document/chunk/embedding coverage, workflow orchestration, queue health, and embedding similarity.
- The highest leverage refactor is not visual redesign; it is ownership extraction: tokens, foundation primitives, layout primitives, component layers, Storybook, DMETA IR, and validators.
- `SearchView.tsx` is the best first vertical slice because it touches query composition, retrieval results, coverage, result inspection, and corpus navigation.

### What was tricky to build

- The tricky part was avoiding a superficial CSS cleanup report. The current CSS is coherent and useful; the problem is not that it exists, but that one global CSS file owns too many layers. The report therefore recommends preserving the retro monochrome language while extracting ownership boundaries.
- Another tricky part was deciding how much backend context to include. The report focuses on backend files only when they explain the web app: server wiring, API routes, database concepts, corpus/search/workflow handlers, and endpoint contracts.

### What warrants a second pair of eyes

- Whether the first DMETA vertical slice should be Search-first, as recommended, or Corpus-first because Corpus Explorer has the richest explanatory model.
- Whether the retro `components/retro/*` files should be deleted or revived as part of the future primitive layer.
- Whether cross-view navigation should move to React Router, route params, or a small typed dashboard navigation store.

### What should be done in the future

- Install frontend dependencies and run `pnpm build` to separate dependency issues from TypeScript issues.
- Implement the first Search Workbench design-system/DMETA vertical slice.
- Add a `dmeta-ir` directory with human-readable YAML fields for core model, Interaction IR, Web MDS templates, and lowering rules.

### Code review instructions

- Start with the design doc:
  - `ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/design-doc/01-rag-evaluation-web-architecture-and-design-system-review.md`
- Compare claims against:
  - `web/src/App.tsx`
  - `web/src/services/api.ts`
  - `web/src/index.css`
  - `web/src/components/search/SearchView.tsx`
  - `web/src/components/corpus/CorpusExplorerView.tsx`
  - `web/src/components/workflows/WorkflowsView.tsx`
  - `internal/api/handlers.go`
  - `internal/api/workflow_artifact_handlers.go`
- Validate ticket hygiene with:
  - `docmgr --root ttmp doctor --ticket RAG-WEB-DESIGN-SYSTEM-REVIEW --stale-after 30`

### Technical details

Useful commands from the investigation:

```bash
cd 2026-05-27--rag-evaluation-system

docmgr --root ttmp status --summary-only
find web/src -maxdepth 5 -type f | sort
rg -n "rag:navigate|className=|style=\{\{" web/src -S
find . -maxdepth 3 \( -name '*dmeta*' -o -name 'dmeta-ir' -o -name '*.stories.tsx' -o -name '*.module.css' \) -print
cd web && pnpm build
```

## Step 2: Validate and Upload the Report

I completed the delivery loop by resolving docmgr vocabulary warnings, rerunning `docmgr doctor`, and uploading the report bundle to reMarkable. The final bundle includes both the primary review and this investigation diary.

The upload completed successfully and returned the expected `OK: uploaded` message.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Finish ticket hygiene and deliver the intern-facing RAG web review as a reMarkable PDF bundle.

**Inferred user intent:** The user wants the report stored in the RAG ticket workspace and available for offline reading/review.

**Commit (code):** N/A — documentation/review step.

### What I did

- Related key web, API, server, database, and ticket files to the design doc and diary.
- Added missing vocabulary entries for ticket topics:
  - `react`
  - `css`
  - `design-system`
  - `dmeta`
- Ran:
  - `docmgr --root ttmp doctor --ticket RAG-WEB-DESIGN-SYSTEM-REVIEW --stale-after 30`
- Ran a reMarkable dry run for the report bundle.
- Uploaded the report bundle to:
  - `/ai/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW`

### Why

- The ticket should pass docmgr validation before handoff.
- The reMarkable bundle is the requested delivery artifact for reading and review.

### What worked

- After adding vocabulary entries, `docmgr doctor` reported `All checks passed`.
- The dry run showed the expected report and diary inputs.
- The real upload returned: `OK: uploaded RAG Web Architecture Design System Review.pdf -> /ai/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW`.

### What didn't work

- The first `docmgr doctor` run warned that `css`, `design-system`, `dmeta`, and `react` were unknown topics in the RAG repository vocabulary. I resolved this by adding the topic slugs to `ttmp/vocabulary.yaml` with `docmgr --root ttmp vocab add`.

### What I learned

- The RAG repository's docmgr vocabulary was narrower than the TTC design-system vocabulary, so design-system review tickets need explicit topic additions before doctor passes.

### What was tricky to build

- The subtle part was using the RAG repository's local docmgr root and vocabulary rather than the outer workspace's TTC docmgr configuration. Running docmgr from inside `2026-05-27--rag-evaluation-system` with `--root ttmp` used the correct local `.ttmp.yaml` and vocabulary.

### What warrants a second pair of eyes

- The added vocabulary entries should be acceptable for future RAG web/design-system tickets; if the project prefers `ui` or `design` only, the topic taxonomy can be revisited.

### What should be done in the future

- Commit the ticket documents and vocabulary update.
- Start a follow-up implementation ticket for the Search Workbench vertical slice.

### Code review instructions

- Confirm `docmgr --root ttmp doctor --ticket RAG-WEB-DESIGN-SYSTEM-REVIEW --stale-after 30` passes.
- Confirm the uploaded PDF path in the final handoff.

### Technical details

Validation and upload commands:

```bash
cd 2026-05-27--rag-evaluation-system

docmgr --root ttmp doctor --ticket RAG-WEB-DESIGN-SYSTEM-REVIEW --stale-after 30

remarquee upload bundle \
  2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/design-doc/01-rag-evaluation-web-architecture-and-design-system-review.md \
  2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md \
  --name "RAG Web Architecture Design System Review" \
  --remote-dir "/ai/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW" \
  --toc-depth 2 \
  --non-interactive
```

## Step 3: Implement Phases 0, 1, and the Start of Phase 2

I started the implementation work requested after the review by making the web package buildable locally, adding Storybook scaffolding, extracting the first design-system token file, and adding the first foundation and layout primitives. This is intentionally a bottom-up step: it does not split the Search page yet, but it creates the primitives and Storybook surface that the Search Workbench refactor should use next.

I also updated the ticket tasks so Storybook work starts in Phase 2 immediately, not after all React refactoring is done. The new stories are co-located with the primitives so later extracted organisms can follow the same pattern.

### Prompt Context

**User prompt (verbatim):** "ok, hase 0 1 2. don't forget to add tasks to write the storybook stories, start right away in phase 2 already."

**Assistant interpretation:** Begin implementation for phases 0, 1, and 2 of the RAG web design-system plan, and make sure Storybook stories are represented as explicit tasks and started immediately during Phase 2.

**Inferred user intent:** The user wants implementation momentum, not another planning-only pass, and wants Storybook to become part of the design-system loop before large component extraction.

**Commit (code):** N/A at time of writing — implementation changes pending commit.

### What I did

- Installed web dependencies with `pnpm install`.
- Added Storybook dependencies with `pnpm add -D storybook @storybook/react-vite`.
- Added package scripts:
  - `typecheck`
  - `storybook`
  - `build-storybook`
- Added `.storybook/main.ts` and `.storybook/preview.ts`.
- Added `storybook-static` to `web/.gitignore`.
- Extracted the current retro `:root` tokens from `web/src/index.css` into `web/src/styles/tokens.css`.
- Added RAG role token aliases such as `--rag-font-role-body`, `--rag-font-role-code`, and `--rag-font-role-label`.
- Added CSS module typing support to `web/src/vite-env.d.ts`.
- Added foundation primitives with co-located stories:
  - `Text`
  - `CodeText`
  - `StatusText`
  - `Divider`
  - `VisuallyHidden`
- Added layout primitives with co-located stories:
  - `Panel`
  - `Stack`
  - `Inline`
  - `DashboardGrid`
- Added component barrels:
  - `web/src/components/foundation/index.ts`
  - `web/src/components/layout/index.ts`
  - `web/src/components/index.ts`
- Updated ticket tasks to include Phase 0/1/2 completion and Phase 2 follow-up Storybook tasks for Search Workbench organisms.

### Why

- Phase 0 needed a working local package baseline before any serious refactor.
- Phase 1 needed explicit token/foundation/layout ownership so the next component extractions do not keep adding global classes and inline styles.
- Phase 2 needed Storybook immediately so primitives are reviewable before the Search Workbench is split into organisms and molecules.

### What worked

- `pnpm install` completed from the existing lockfile.
- `pnpm add -D storybook @storybook/react-vite` completed and updated `package.json` / `pnpm-lock.yaml`.
- `pnpm typecheck` passed after fixing CSS module typing and the `Panel` prop name conflict.
- `pnpm build` passed.
- `pnpm build-storybook` passed.

### What didn't work

- The first `pnpm typecheck` failed because `noUncheckedIndexedAccess` made CSS module class lookups type as `string | undefined` under the initial declaration:

```text
src/components/foundation/Text/Text.tsx(21,3): error TS2322: Type 'string | undefined' is not assignable to type 'string'.
```

- The first `pnpm typecheck` also failed because `PanelProps` extended `HTMLAttributes<HTMLDivElement>` while redefining the native `title` attribute as `ReactNode`:

```text
src/components/layout/Panel/Panel.tsx(4,18): error TS2430: Interface 'PanelProps' incorrectly extends interface 'HTMLAttributes<HTMLDivElement>'.
```

- `StatusText.stories.tsx` initially needed default `args` because `status` is a required component prop.
- `go test ./...` failed because this checkout is missing the sibling `../scraper` module required by the local replace directive:

```text
github.com/go-go-golems/scraper@v0.0.0 (replaced by ./scraper): reading ../scraper/go.mod: open /home/manuel/workspaces/2026-05-27/ttc-design-system/scraper/go.mod: no such file or directory
```

### What I learned

- The web package was buildable once dependencies were installed; the earlier build failure was environmental rather than a TypeScript failure.
- Storybook 10.4.1 builds successfully against the current React/Vite package stack.
- The existing retro global CSS can be preserved while introducing token aliases and primitive components; no visual reset is required for this first step.

### What was tricky to build

- The tricky part was introducing CSS Modules under strict TypeScript with `noUncheckedIndexedAccess`. The immediate pragmatic fix was to type CSS module imports as `any` in `vite-env.d.ts`. A future stricter typed-css-module setup can improve this, but it should not block Phase 1/2 scaffolding.
- Another subtle issue was avoiding a generated build-artifact commit. `pnpm build` rewrote `internal/web/dist/index.html` to reference ignored hashed assets, so I reverted that tracked build artifact and kept the source/package changes only.

### What warrants a second pair of eyes

- Whether `CodeText`, `StatusText`, and `Text` should remain separate primitives or whether `CodeText`/`StatusText` should become roles on a richer `Text` primitive later.
- Whether `DashboardGrid` recipes should stay minimal until the Search Workbench extraction proves the exact recipes.
- Whether the CSS module typing should be tightened after the first refactor.

### What should be done in the future

- Split `SearchView.tsx` into Search Workbench page/organism/molecule components.
- Add Storybook stories for `SearchControlsPanel`, `RetrievalResultsPanel`, `ResultInspectorPanel`, and `CoveragePanel` during the Phase 2 follow-up.
- Add initial `dmeta-ir` YAML in Phase 3 after the first component boundaries stabilize.

### Code review instructions

- Start with `web/src/styles/tokens.css` to see how current retro tokens were preserved.
- Review foundation primitives under `web/src/components/foundation/`.
- Review layout primitives under `web/src/components/layout/`.
- Check Storybook setup under `web/.storybook/` and co-located `*.stories.tsx` files.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

Repository-level Go validation attempted but blocked by the missing sibling scraper checkout:

```bash
cd 2026-05-27--rag-evaluation-system
go test ./...
```

## Step 4: Extract the First Search Workbench Components and Stories

I continued Phase 2 by extracting the first Search Workbench organisms and molecules from `SearchView.tsx`. This keeps the current runtime behavior intact while creating named component boundaries and Storybook review surfaces. The result is intentionally incremental: controls, results, coverage, and query presets are now reusable components, while the larger `ResultInspector` remains inside `SearchView.tsx` for the next pass.

This step follows the playbook's vertical-slice approach. The Search view is no longer one monolithic file for every visible obligation. It now starts to express a page composition: query controls, optional coverage warning, query presets, retrieval results, and selected-result inspector.

### Prompt Context

**User prompt (verbatim):** "continue."

**Assistant interpretation:** Continue the Phase 2 implementation from the previous turn by extracting Search Workbench components and expanding Storybook coverage.

**Inferred user intent:** The user wants concrete progress through the Search Workbench design-system vertical slice, with Storybook stories created as components are extracted.

**Commit (code):** N/A at time of writing — implementation changes pending commit.

### What I did

- Extracted `SearchControlsPanel` to:
  - `web/src/components/organisms/SearchControlsPanel/SearchControlsPanel.tsx`
  - `web/src/components/organisms/SearchControlsPanel/SearchControlsPanel.stories.tsx`
- Extracted `RetrievalResultsPanel` to:
  - `web/src/components/organisms/RetrievalResultsPanel/RetrievalResultsPanel.tsx`
  - `web/src/components/organisms/RetrievalResultsPanel/RetrievalResultsPanel.stories.tsx`
- Extracted `CoveragePanel` to:
  - `web/src/components/molecules/CoveragePanel/CoveragePanel.tsx`
  - `web/src/components/molecules/CoveragePanel/CoveragePanel.stories.tsx`
- Extracted `QueryPresetList` to:
  - `web/src/components/molecules/QueryPresetList/QueryPresetList.tsx`
  - `web/src/components/molecules/QueryPresetList/QueryPresetList.stories.tsx`
- Added molecule and organism barrels:
  - `web/src/components/molecules/index.ts`
  - `web/src/components/organisms/index.ts`
- Updated `web/src/components/index.ts` to export molecules and organisms.
- Refactored `SearchView.tsx` to import and compose the extracted components.
- Updated ticket tasks to mark the extracted components and stories as done, while leaving `ResultInspectorPanel` extraction as the next Phase 2 task.

### Why

- Search is the recommended first vertical slice from the review.
- Storybook should start during Phase 2, not after all refactors are complete.
- Extracting controls/results/coverage first reduces the size and responsibility of `SearchView.tsx` without touching API behavior or navigation behavior.

### What worked

- `SearchView.tsx` now composes `SearchControlsPanel`, `CoveragePanel`, `QueryPresetList`, and `RetrievalResultsPanel`.
- Storybook now includes organism-level stories for the search controls and retrieval results.
- Storybook now includes molecule-level stories for coverage states and query presets.
- Validation passed:
  - `pnpm typecheck`
  - `pnpm build`
  - `pnpm build-storybook`

### What didn't work

- The first typecheck after adding stories failed because Storybook's typed `StoryObj<typeof meta>` required `args` for stories whose components have required props. I fixed this by using looser `StoryObj` typing for stateful render-only stories in `SearchControlsPanel.stories.tsx` and `QueryPresetList.stories.tsx`.

### What I learned

- Search Workbench extraction can proceed safely in smaller units; it does not require rewriting the whole page at once.
- `ResultInspector` is the next high-value extraction because it still owns a lot of tabbed detail/chunk/document behavior.
- Stories for data-heavy RAG components need realistic mock DTOs, not placeholder text, to be useful for review.

### What was tricky to build

- The main tricky part was preserving behavior while moving JSX. `SearchControlsPanel` needs many state setter props because behavior still lives in `SearchView.tsx`. This is acceptable for the first extraction, but a future `useSearchWorkbenchState` hook could make the page cleaner.
- Another subtle part was Storybook typing for required callback-heavy components. The stateful story pattern is useful, but stricter Storybook generic inference expects required props in `args`. I used `StoryObj` for these stories to keep the examples concise.

### What warrants a second pair of eyes

- Whether `SearchControlsPanel` should stay an organism or become a molecule once the full Search Workbench page exists.
- Whether `CoveragePanel` belongs under molecules or organisms; it is currently a small dependent panel, so molecule is a reasonable starting point.
- Whether `RetrievalResultsPanel` should depend on a reusable `DataTable` molecule before more result-table variants appear.

### What should be done in the future

- Extract `ResultInspectorPanel` and add stories for detail, chunk, document, loading, and missing-document states.
- Replace the remaining Search view inline layout with `DashboardGrid` and `Stack` after the inspector is extracted.
- Start Phase 3 DMETA IR after the Search Workbench component boundaries are stable enough to model.

### Code review instructions

- Start with `web/src/components/search/SearchView.tsx` to see the new composition.
- Review extracted components under:
  - `web/src/components/organisms/SearchControlsPanel/`
  - `web/src/components/organisms/RetrievalResultsPanel/`
  - `web/src/components/molecules/CoveragePanel/`
  - `web/src/components/molecules/QueryPresetList/`
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 5: Extract the Result Inspector and Move Search Layout onto Primitives

I finished the next Phase 2 increment by pulling the selected-result inspector out of `SearchView.tsx` and into a dedicated `ResultInspectorPanel` organism. This keeps document-fetching orchestration in the Search page while moving the tabbed result/detail/chunk/document presentation into a reusable component with Storybook examples.

I also moved the Search Workbench page shell away from a page-level inline flex layout and onto the new `DashboardGrid` and `Stack` layout primitives. The remaining page-local CSS now only owns Search-view sizing and scroll behavior in `SearchView.module.css`.

### Prompt Context

**User prompt (verbatim):** "continue, keep a diary, commit at appropriate intervals."

**Assistant interpretation:** Continue the RAG Phase 2 implementation, maintain the investigation diary as work proceeds, and create focused commits after validated increments.

**Inferred user intent:** The user wants momentum through the Search Workbench design-system slice without losing implementation traceability or leaving uncommitted validated work behind.

**Commit (code):** N/A at time of writing — implementation changes pending commit.

### What I did

- Extracted `ResultInspectorPanel` to:
  - `web/src/components/organisms/ResultInspectorPanel/ResultInspectorPanel.tsx`
  - `web/src/components/organisms/ResultInspectorPanel/ResultInspectorPanel.stories.tsx`
  - `web/src/components/organisms/ResultInspectorPanel/index.ts`
- Updated `web/src/components/organisms/index.ts` to export `ResultInspectorPanel`.
- Refactored `SearchView.tsx` so it:
  - fetches selected document detail with `useGetCorpusDocumentQuery`,
  - dispatches the existing `rag:navigate-to-chunk` event from the page layer,
  - passes `result`, `docDetail`, `onClose`, and `onOpenInCorpus` into `ResultInspectorPanel`.
- Added `web/src/components/search/SearchView.module.css` for Search page shell sizing/scroll behavior.
- Replaced Search page shell inline flex layout with:
  - `DashboardGrid recipe="searchWorkbench"`,
  - `Stack` for the left controls rail,
  - `Stack` for the results column.
- Added Storybook examples for:
  - a populated hybrid inspector,
  - document-loading state,
  - BM25-only result state.

### Why

- `ResultInspectorPanel` was the largest remaining chunk of presentational UI inside `SearchView.tsx`.
- Keeping the RTK Query hook in `SearchView.tsx` preserves the current data-flow boundary while making the inspector storyable with mock data.
- Moving the outer shell to `DashboardGrid` proves the layout primitives can host real application pages, not just isolated Storybook demos.

### What worked

- The extraction preserved TypeScript correctness.
- The Storybook build now includes `ResultInspectorPanel`.
- Validation passed:
  - `pnpm typecheck`
  - `pnpm build`
  - `pnpm build-storybook`

### What didn't work

- `pnpm build` again rewrote `internal/web/dist/index.html` because Vite emits the SPA into the Go embed directory. I reverted that file before committing, as in the previous step, because the emitted hashed assets are ignored and the tracked HTML would otherwise point at ignored build artifacts.
- I did not add separate stories that open the chunk/document tabs by default because `ResultInspectorPanel` currently owns its active tab internally. Those states are reachable interactively in Storybook, but explicit initial-tab stories should wait for a small controlled-tab API or play-function interaction.

### What I learned

- The useful boundary is container/page versus presentational organism: Search page owns API fetching and cross-view navigation events; `ResultInspectorPanel` owns result-detail rendering.
- The new `DashboardGrid` primitive already has a `searchWorkbench` recipe that matches this page well enough to remove the top-level inline flex shell.
- The inspector's internal tab state is a small obstacle to exhaustive static story coverage.

### What was tricky to build

- The main tricky part was not moving the RTK Query hook into the presentational component. Doing so would make Storybook harder to isolate and would entangle the organism with live API assumptions. I kept the query in `SearchView.tsx` and passed `docDetail` as data.
- The `Open in Corpus` action also needed to remain behavior-preserving. I moved only the event dispatch callback to the page and passed it as `onOpenInCorpus`, keeping the custom event contract unchanged.

### What warrants a second pair of eyes

- Whether `ResultInspectorPanel` should expose a controlled `activeTab`/`defaultTab` API for stronger Storybook state coverage.
- Whether the `rag:navigate-to-chunk` custom event should be replaced by typed route state or a dashboard navigation store in a later phase.
- Whether page shell sizing belongs in `SearchView.module.css` or should become a reusable app-frame primitive once more pages migrate.

### What should be done in the future

- Add explicit chunk-tab and document-tab stories via controlled tab props or Storybook play interactions.
- Start Phase 3 `dmeta-ir` YAML now that the Search Workbench component boundary is stable enough to model.
- Keep replacing inline internals gradually, but avoid abstracting table rows or metadata grids until repetition across pages proves the need.

### Code review instructions

- Start with `web/src/components/search/SearchView.tsx` to verify that API fetching and navigation behavior stayed in the page layer.
- Review `web/src/components/organisms/ResultInspectorPanel/ResultInspectorPanel.tsx` for the presentational split.
- Review `web/src/components/search/SearchView.module.css` for the new page-local layout ownership.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 6: Seed the Search Workbench DMETA IR Catalog

I started Phase 3 by adding the first RAG-local `dmeta-ir` catalog. This is intentionally a seed catalog rather than a generator-ready implementation: it names the Search Workbench core concepts, visible Interaction IR obligations/actions, Web MDS templates, lowering rules, token/recipe references, React target metadata, and the first instantiation.

The catalog follows the boundary rules from the TTC design-system work. Core model files name target-neutral retrieval/corpus concepts; Interaction IR files name visible obligations and actions; Web MDS files choose Web dashboard templates, component levels, coarse layout recipes, CSS ownership, and Storybook contracts.

### Prompt Context

**User prompt (verbatim):** (same as Step 5)

**Assistant interpretation:** Continue the implementation after the Search Workbench component extraction, keeping documentation and commits aligned.

**Inferred user intent:** The user wants the RAG design-system work to move from React extraction into the DMETA modeling layer while preserving an implementation diary.

**Commit (code):** N/A at time of writing — implementation changes pending commit.

### What I did

- Added `dmeta-ir/README.md` explaining the layer boundaries.
- Added explicit manifests:
  - `dmeta-ir/core-model/00-index.yaml`
  - `dmeta-ir/interactions/00-index.yaml`
  - `dmeta-ir/meta-design-systems/web/00-index.yaml`
- Added core model seed files:
  - `dmeta-ir/core-model/archetypes/10-search-workbench.yaml`
  - `dmeta-ir/core-model/capabilities/10-search-workbench.yaml`
- Added Interaction IR seed files:
  - `dmeta-ir/interactions/representations/10-search-workbench.yaml`
  - `dmeta-ir/interactions/actions/10-search-workbench.yaml`
  - `dmeta-ir/interactions/elaboration-rules/10-search-workbench.yaml`
- Added Web MDS seed files:
  - `dmeta-ir/meta-design-systems/web/component-system.yaml`
  - `dmeta-ir/meta-design-systems/web/lowering-rules.yaml`
  - `dmeta-ir/meta-design-systems/web/widgets/search-workbench.yaml`
  - `dmeta-ir/meta-design-systems/web/style/tokens.yaml`
  - `dmeta-ir/meta-design-systems/web/style/recipes.yaml`
  - `dmeta-ir/meta-design-systems/web/targets/react.yaml`
- Added the first instantiation:
  - `dmeta-ir/instantiations/rag-evaluation-dashboard.yaml`
- Validated that every YAML file parses with Python/PyYAML:
  - `python3 - <<'PY' ... yaml.safe_load(...) ... PY`

### Why

- The review called out that the RAG repository had no `dmeta-ir`, which meant the important concepts only existed in React and Go code.
- The Search Workbench now has extracted React boundaries, making it a good moment to formalize the IR names behind those components.
- Explicit manifests and human-readable fields reduce future ambiguity and avoid glob-based catalog loading.

### What worked

- All YAML files parsed successfully.
- The catalog maps directly onto the components already extracted in Phase 2:
  - `SearchControlsPanel`
  - `CoveragePanel`
  - `QueryPresetList`
  - `RetrievalResultsPanel`
  - `ResultInspectorPanel`
- The Web MDS templates include Storybook contracts that mirror existing stories.

### What didn't work

- There is not yet a repository-local DMETA validator for the RAG catalog. I used a parse-only YAML check for this step and recorded a follow-up to add duplicate-ID, manifest-reference, and human-readable-field validation.
- There is not yet a generator path for RAG React scaffolds, so the instantiation records selected representations/actions/templates but does not produce generated output.

### What I learned

- The Search Workbench slice is now cohesive enough to model cleanly: query composition, coverage, presets, retrieval results, and result inspection each have clear IR representations.
- The current custom event for `open_in_corpus` belongs in action `human_notes`; it is a known implementation detail, not a target-neutral contract.
- Web MDS can stay coarse: component level, role, layout primitive/recipe, composition, CSS ownership, and stories are enough for this seed.

### What was tricky to build

- The main tricky part was avoiding React leakage into core or Interaction IR. For example, the core `RetrievalResult` archetype describes ranked evidence, while `RetrievalResultsPanel` only appears in Web MDS.
- Another tricky point was naming coverage. I modeled `embedding_coverage_summary` as an Interaction IR representation because it is a visible obligation that changes how users interpret vector/hybrid search quality, not just a visual widget.

### What warrants a second pair of eyes

- Whether `query_preset_list` should remain a representation or become supporting Web-only content if presets stay purely demo-oriented.
- Whether `copy_artifact_identity` belongs in the first slice or should be generalized later across corpus/workflow inspectors.
- Whether Web MDS `css_ownership` values should align exactly with TTC's `none`, `anatomy_only`, `widget_internal`, and `page_chrome` vocabulary before validators are added.

### What should be done in the future

- Add a validation script for duplicate IDs, manifest references, required human-readable fields, and lowering-rule references.
- Add generated scaffold output only after the RAG catalog has a tooling consumer.
- Extend the catalog to corpus explorer and workflow monitor after the Search Workbench slice is validated.

### Code review instructions

- Start with `dmeta-ir/README.md` for the layer contract.
- Then review:
  - `dmeta-ir/core-model/archetypes/10-search-workbench.yaml`
  - `dmeta-ir/interactions/representations/10-search-workbench.yaml`
  - `dmeta-ir/meta-design-systems/web/widgets/search-workbench.yaml`
  - `dmeta-ir/meta-design-systems/web/lowering-rules.yaml`
- Validate YAML parsing with:
  - `python3 - <<'PY' ... yaml.safe_load(...) ... PY`

### Technical details

Parse validation command:

```bash
cd 2026-05-27--rag-evaluation-system
python3 - <<'PY'
from pathlib import Path
import yaml
for p in Path('dmeta-ir').rglob('*.yaml'):
    with p.open() as f:
        yaml.safe_load(f)
print('yaml ok')
PY
```

## Step 7: Add a Local DMETA IR Validator

I added the first repository-local validator for the new RAG `dmeta-ir` catalog. The validator is intentionally lightweight and scriptable: it parses YAML, checks explicit manifest references, enforces human-readable fields on major objects, detects duplicate IDs, and verifies the most important representation/action/template references.

This gives the new catalog a safety net before any generator exists. It also turns the phase-3 follow-up from a prose recommendation into an executable repository command.

### Prompt Context

**User prompt (verbatim):** (same as Step 5)

**Assistant interpretation:** Continue with the next safe increment after seeding the DMETA IR catalog, maintaining diary and committing validated changes.

**Inferred user intent:** The user wants the newly created IR not to become ungoverned YAML; it should have at least basic validation immediately.

**Commit (code):** N/A at time of writing — implementation changes pending commit.

### What I did

- Added `scripts/validate-dmeta-ir.py`.
- Added `make dmeta-validate` to the root `Makefile`.
- The validator checks:
  - all `dmeta-ir/**/*.yaml` files parse,
  - `00-index.yaml` manifest targets exist,
  - major catalog objects have non-empty `summary`, `description`, and `why`,
  - IDs are globally unique,
  - key capability/representation/action/template references resolve.
- Ran:
  - `make dmeta-validate`

### Why

- The catalog now has enough structure to enforce basic invariants.
- Human-readable YAML fields were a key design decision; enforcing them early prevents the IR from becoming opaque implementation metadata.
- Explicit manifests and duplicate-ID checks match the TTC lessons learned.

### What worked

- `make dmeta-validate` passed:
  - `DMETA IR validation passed (15 YAML files, 36 IDs)`

### What didn't work

- The first implementation of the validator incorrectly treated every list named `representations`, `actions`, or `capabilities` as a list of major objects. In this catalog those names also appear in `consumes` and `outputs` blocks as string references. The validator initially failed with many false positives such as:
  - `dmeta-ir/meta-design-systems/web/lowering-rules.yaml: representations[0] should be an object`
- I fixed this by treating only list items with an `id` field as major catalog objects and by skipping semantic reference validation inside manifest files.

### What I learned

- YAML validators need to understand context; field names like `actions` can mean object collections or reference lists.
- The current catalog already has enough reference structure to catch useful errors without a full DMETA loader.

### What was tricky to build

- The tricky part was avoiding over-validation. A strict schema would be premature because the RAG IR is still a seed. The validator now enforces durable invariants while leaving room for the catalog shape to evolve.
- Another subtle issue was global duplicate ID scope. I chose global IDs across all major objects because earlier TTC work established duplicate IDs across merged catalog files should fail during load.

### What warrants a second pair of eyes

- Whether global ID uniqueness should include token-family and recipe IDs forever, or whether future tooling should namespace some categories.
- Whether `PyYAML` is acceptable as an implicit operator dependency or should be replaced by a small Go validator later.

### What should be done in the future

- Add validation for selected template/action/representation closure against instantiations.
- Add validation for Web MDS composition dependency closure.
- Move this logic into the eventual Go DMETA tooling if/when RAG generation starts.

### Code review instructions

- Review `scripts/validate-dmeta-ir.py` for validation scope and false-positive avoidance.
- Review `Makefile` for the new `dmeta-validate` target.
- Validate with:
  - `make dmeta-validate`

### Technical details

Validation output:

```text
python3 scripts/validate-dmeta-ir.py
DMETA IR validation passed (15 YAML files, 36 IDs)
```

## Step 8: Restrict the Current Scope to Documentation-Only IR Through Phase 5

I adjusted the work after the user clarified that the IR is currently documentation-only and should not be validated or used for generated scaffolds. The implementation now stops at the Phase 5 outcome: a hand-authored Search Workbench DMETA-style YAML catalog that documents the intended core, interaction, and Web MDS layers.

This step deliberately removes the executable IR validator added in the prior step and leaves Phase 6 scaffolding plus Phase 7 validation as deferred, out-of-scope work. The diary keeps the earlier validator attempt as historical context because it explains what was tried and why it was backed out.

### Prompt Context

**User prompt (verbatim):** "no need to validate or scaffold from the IR< currently we are just creating the yaml as documentation purposes. 

do up to phase 5 (included), but none of phase 6 and 7"

**Assistant interpretation:** Keep the RAG DMETA IR as documentation only, complete work through the report's Phase 5, and avoid Phase 6 generated scaffolds and Phase 7 validators.

**Inferred user intent:** The user wants useful YAML documentation and design-system/component extraction, but does not want premature tooling around the IR yet.

**Commit (code):** N/A at time of writing — scope correction pending commit.

### What I did

- Removed `scripts/validate-dmeta-ir.py`.
- Removed the `make dmeta-validate` target from `Makefile`.
- Updated `dmeta-ir/README.md` to state that the catalog is documentation-only for now.
- Updated `dmeta-ir/instantiations/rag-evaluation-dashboard.yaml` to state that no validation or generated scaffolds are required in the current scope.
- Updated `tasks.md` so:
  - Search component extraction maps to Phase 3,
  - Storybook review surfaces map to Phase 4,
  - documentation-only `dmeta-ir` maps to Phase 5,
  - generated scaffolds are deferred as Phase 6,
  - executable validators are deferred as Phase 7.
- Removed stale docmgr related-file entries for the deleted validator script and Makefile validator target.

### Why

- The user's scope clarification supersedes the previous validator follow-up.
- The YAML files are meant to document design intent right now, not drive code generation or enforce repository contracts.
- Removing premature tooling keeps the work aligned with the requested phase boundary.

### What worked

- The repository now contains the Phase 5 documentation artifact without Phase 6/7 tooling.
- The ticket tasks now match the phase numbering in the report.
- `docmgr --root ttmp doctor --ticket RAG-WEB-DESIGN-SYSTEM-REVIEW --stale-after 30` passed after removing stale file relations.

### What didn't work

- The previous step went one phase too far by adding a validator. That is now backed out rather than expanded.

### What I learned

- For this RAG effort, the near-term value is a readable design/IR reference, not executable IR governance.
- Phase boundaries matter: validators and scaffolds should wait until the user explicitly asks for Phase 6 or Phase 7 work.

### What was tricky to build

- The tricky part was preserving the diary trail while removing the actual validator. I left the previous diary entry intact as historical record, then added this corrective step so future readers understand why the script disappeared.
- Docmgr related files also needed cleanup because removed files should not remain in active document frontmatter.

### What warrants a second pair of eyes

- Whether the Phase 5 YAML should keep `generation_policy` fields as future intent or soften that language further while the catalog remains documentation-only.
- Whether the task list should move the optional `ResultInspectorPanel` tab-state stories into a separate future ticket.

### What should be done in the future

- Do not add IR validators unless Phase 7 is explicitly requested.
- Do not add generated React scaffolds or promotion metadata unless Phase 6 is explicitly requested.
- Keep improving the YAML as human-readable documentation if additional Search Workbench concepts need clarification.

### Code review instructions

- Confirm `scripts/validate-dmeta-ir.py` is removed.
- Confirm `Makefile` no longer exposes `dmeta-validate`.
- Review `dmeta-ir/README.md` and `dmeta-ir/instantiations/rag-evaluation-dashboard.yaml` for documentation-only wording.
- Review `tasks.md` for Phase 3/4/5 completion and Phase 6/7 deferral.

### Technical details

The scope boundary is now:

```text
Phase 3: split SearchView into reusable pieces — done
Phase 4: add Storybook/component review surface — done
Phase 5: author documentation-only DMETA IR vertical slice — done
Phase 6: generated scaffolds/metadata — deferred
Phase 7: validators — deferred
```

## Step 9: Expand the Reusable Design-System Primitive Set

I resumed the design-system extraction with a broader pass over the primitives and reusable molecules called out in the design doc. This step adds the missing low-level pieces that are needed before refactoring the remaining Search, Corpus, and Workflow surfaces more aggressively.

The new components are intentionally small and Storybook-backed. They do not redesign the dashboard; they provide named ownership for recurring caption text, app shell structure, scroll regions, tabs, form rows, metadata grids, and dense data tables.

### Prompt Context

**User prompt (verbatim):** "ok, it's great, now go full in in terms of extracting desig system and components and pages."

**Assistant interpretation:** Continue beyond the minimal Search Workbench slice and aggressively extract the reusable design-system/component/page building blocks identified in the design doc, while still avoiding Phase 6 scaffolding and Phase 7 validators unless requested.

**Inferred user intent:** The user wants the React design-system architecture to become real across more reusable primitives and components, not just documented as a plan.

**Commit (code):** N/A at time of writing — implementation changes pending commit.

### What I did

- Added foundation primitive:
  - `web/src/components/foundation/Caption/`
- Added layout primitives:
  - `web/src/components/layout/AppShell/`
  - `web/src/components/layout/ScrollRegion/`
  - `web/src/components/layout/TabList/`
  - `web/src/components/layout/FormRow/`
- Added reusable molecules:
  - `web/src/components/molecules/MetadataGrid/`
  - `web/src/components/molecules/DataTable/`
- Added co-located Storybook stories for each new primitive/molecule.
- Updated foundation/layout/molecule barrels.
- Updated ticket tasks to mark the new reusable pieces and Storybook surfaces as complete.

### Why

- The design doc explicitly identified these as reusable pieces needed before wider page extraction.
- `MetadataGrid` and `DataTable` are useful extraction targets for Search result inspection, Corpus Explorer document details, workflow lists, and future evaluation result views.
- `TabList`, `ScrollRegion`, and `FormRow` give the app named primitives for common dashboard mechanics that were previously embedded in page-local/global CSS patterns.

### What worked

- TypeScript passed after adding the new components.
- Storybook built successfully and now includes stories for the expanded primitive/molecule set.

### What didn't work

- I have not yet replaced all existing call sites with these primitives. This step establishes the reusable surface first; the next step should adopt `MetadataGrid`, `DataTable`, `TabList`, and `ScrollRegion` inside real pages/components.

### What I learned

- The next high-leverage refactor is not another primitive; it is adoption. `ResultInspectorPanel`, `RetrievalResultsPanel`, `DocumentInspector`, and `WorkflowsView` all contain patterns that can now move onto `MetadataGrid` and `DataTable`.
- `AppShell` should be adopted carefully because `App.tsx` owns navigation and route-level behavior; it may need one small integration pass rather than a blind replacement.

### What was tricky to build

- The main tricky part was keeping the primitives generic without creating a giant prop-driven framework. `DataTable` owns table mechanics but leaves cell rendering to callers. `MetadataGrid` owns key/value structure but leaves values as React nodes.
- `TabList` is controlled rather than stateful so real components and Storybook can show explicit tab states later.

### What warrants a second pair of eyes

- Whether `DataTable` should support column-level truncation and sticky headers in a more formal API before widespread adoption.
- Whether `AppShell` should model the existing Mac-style menu/navigation more explicitly or stay as a bare layout primitive.
- Whether `Caption` overlaps too much with `Text size="metadata"` or whether the mono/uppercase metadata role justifies a separate primitive.

### What should be done in the future

- Refactor `RetrievalResultsPanel` to use `DataTable`.
- Refactor `ResultInspectorPanel` to use `MetadataGrid`, `TabList`, and `ScrollRegion`.
- Refactor Corpus Explorer and Workflow views to use `DataTable`, `MetadataGrid`, `Panel`, `Stack`, and `ScrollRegion`.
- Extract page-level folders such as `SearchWorkbenchPage` once the current component adoption is stable.

### Code review instructions

- Review new primitives under `web/src/components/foundation/Caption/` and `web/src/components/layout/`.
- Review new molecules under `web/src/components/molecules/MetadataGrid/` and `web/src/components/molecules/DataTable/`.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build-storybook
```

## Step 10: Adopt Reusable Molecules in the Search Workbench Panels

I followed the primitive extraction with an adoption pass inside the Search Workbench. `RetrievalResultsPanel` now uses the reusable `DataTable`, and `ResultInspectorPanel` now uses `MetadataGrid`, `TabList`, and `ScrollRegion` instead of carrying all table, metadata, tab, and scroll structure locally.

This step is the first proof that the newly extracted components are not just Storybook examples. They can replace real UI structure while preserving the dense retro dashboard behavior.

### Prompt Context

**User prompt (verbatim):** (same as Step 9)

**Assistant interpretation:** Continue the full design-system extraction by adopting the new reusable components in real Search Workbench panels.

**Inferred user intent:** The user wants the reusable design-system pieces to actively reduce monolithic page/component code, not remain unused abstractions.

**Commit (code):** N/A at time of writing — implementation changes pending commit.

### What I did

- Refactored `RetrievalResultsPanel` to use:
  - `Panel`
  - `ScrollRegion`
  - `DataTable`
  - `Caption`
- Refactored `ResultInspectorPanel` to use:
  - `Panel`
  - `ScrollRegion`
  - `TabList`
  - `MetadataGrid`
  - `Caption`
  - `StatusText`
- Added `ResultInspectorPanel.module.css` for inspector-local anatomy that is not a reusable primitive yet:
  - section spacing,
  - chunk/document text blocks,
  - neighbor blocks,
  - coverage strip dots.
- Added `defaultTab` support to `ResultInspectorPanel` so Storybook can show explicit detail/chunk/document states.
- Expanded `ResultInspectorPanel` stories with `ChunkTab` and `DocumentTab`.
- Updated tasks to mark Search Workbench adoption complete for these reusable molecules/primitives.

### Why

- `DataTable` and `MetadataGrid` are only valuable if real views use them.
- Search Workbench is still the safest vertical slice for adoption because its behavior is already covered by component stories and build checks.
- Controlled/default tab support makes Storybook more useful without coupling stories to browser click interactions.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed.
- Storybook now includes explicit inspector tab states and the panels compile against reusable molecules.

### What didn't work

- `pnpm build` again rewrote `internal/web/dist/index.html`; I reverted that generated embed artifact before preparing the commit.
- I did not yet move Corpus Explorer or Workflow views onto the new molecules. Those are the next broader adoption targets.

### What I learned

- `ResultInspectorPanel` still needs component-specific anatomy CSS for the coverage strip and text blocks. That is fine: not every repeated-looking pattern deserves a primitive immediately.
- `DataTable` works for retrieval results, but broader Corpus/Workflow adoption will reveal whether it needs sorting, column width, wrapping, or empty-state APIs.

### What was tricky to build

- The main tricky part was preserving selected-row behavior while moving retrieval results to `DataTable`. The row selection callback still toggles the selected result exactly as before.
- Another tricky point was making `ResultInspectorPanel` storyable without making all tab state externally controlled. `defaultTab` is a small compromise: stories can start in a specific tab, while the component remains internally interactive.

### What warrants a second pair of eyes

- Whether `defaultTab` is sufficient or whether `activeTab/onTabChange` should be added before more controlled stories are needed.
- Whether `DataTable` should support semantic row labels or keyboard navigation before it is used more broadly.
- Whether the `Panel` body/fill behavior should be tightened for full-height panel layouts.

### What should be done in the future

- Adopt `DataTable` and `MetadataGrid` in Corpus Explorer document/source/chunk surfaces.
- Adopt `DataTable` and `StatusText` in Workflow list/op graph surfaces.
- Consider extracting a reusable `CoverageStrip` only after corpus/search coverage strips repeat enough.

### Code review instructions

- Start with `web/src/components/organisms/RetrievalResultsPanel/RetrievalResultsPanel.tsx` to review the `DataTable` adoption.
- Then review `web/src/components/organisms/ResultInspectorPanel/ResultInspectorPanel.tsx` and `.module.css` for the inspector split.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 11: Promote Search to a Page Boundary and Adopt More Layout Primitives

I continued the broader extraction by creating an explicit Search Workbench page boundary and adopting more of the layout/foundation primitives in application code. The old `SearchView` path remains as a compatibility re-export, while the actual page implementation now lives under `components/pages/SearchWorkbenchPage/`.

This step also moves the root app shell onto the new `AppShell` primitive and rewrites `SearchControlsPanel` to use `Panel`, `Stack`, `FormRow`, `ScrollRegion`, and `Caption`. This reduces page-local inline layout and makes the controls panel follow the same design-system vocabulary as the result and inspector panels.

### Prompt Context

**User prompt (verbatim):** (same as Step 9)

**Assistant interpretation:** Continue extracting pages and adopting the new design-system primitives across real application boundaries.

**Inferred user intent:** The user wants the component/page architecture to match the design doc's intended structure, not just expose low-level primitives.

**Commit (code):** N/A at time of writing — implementation changes pending commit.

### What I did

- Added page boundary:
  - `web/src/components/pages/SearchWorkbenchPage/SearchWorkbenchPage.tsx`
  - `web/src/components/pages/SearchWorkbenchPage/SearchWorkbenchPage.module.css`
  - `web/src/components/pages/SearchWorkbenchPage/index.ts`
  - `web/src/components/pages/index.ts`
- Replaced `web/src/components/search/SearchView.tsx` with a compatibility re-export of `SearchWorkbenchPage`.
- Removed the old `web/src/components/search/SearchView.module.css` after moving the page CSS.
- Updated `web/src/components/index.ts` to export pages.
- Updated `web/src/App.tsx` to use `AppShell` for the root application shell.
- Refactored `SearchControlsPanel` to use:
  - `Panel`
  - `Stack`
  - `FormRow`
  - `ScrollRegion`
  - `Caption`
- Added `SearchControlsPanel.module.css` for remaining controls-specific anatomy.

### Why

- The design doc's target structure includes `components/pages/SearchWorkbenchPage/`, not a generic `components/search/SearchView.tsx` as the main implementation location.
- `AppShell` was extracted in the previous step; adopting it in `App.tsx` proves it can own real application shell structure.
- `SearchControlsPanel` had many inline styles and global layout patterns that now have named primitives.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed.
- Existing imports of `SearchView` remain valid because `SearchView.tsx` re-exports the page component.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted that generated embed artifact before preparing the commit.
- I have not yet moved Corpus Explorer or Workflow views into page folders. Search is the first page boundary.

### What I learned

- The compatibility re-export pattern is useful for incremental page moves: application imports do not need to change immediately, but the source tree can still move toward the target architecture.
- `SearchControlsPanel` still has small control-specific CSS for button grouping and narrow inputs, which is appropriate component anatomy rather than a reason to over-expand primitives.

### What was tricky to build

- Moving the page file required fixing relative imports from `../molecules`/`../organisms` to `../../...` and service imports from `../../services/api` to `../../../services/api`.
- Preserving compatibility through `SearchView.tsx` avoids a large rename blast radius while still creating the new page folder.

### What warrants a second pair of eyes

- Whether `SearchView` should be removed entirely in a later cleanup once imports have migrated to `SearchWorkbenchPage`.
- Whether `AppShell` should own the existing nav strip more formally rather than receiving it as arbitrary header content.
- Whether `SearchControlsPanel` should expose smaller subcomponents for retriever selection or source filters if reused elsewhere.

### What should be done in the future

- Move Corpus Explorer and Workflow views into page folders.
- Adopt `DataTable`, `MetadataGrid`, `Panel`, and `ScrollRegion` in Corpus Explorer and Workflow pages.
- Consider adding `SearchWorkbenchPage.stories.tsx` once page-level API mocking is available.

### Code review instructions

- Start at `web/src/components/pages/SearchWorkbenchPage/SearchWorkbenchPage.tsx` for the moved page boundary.
- Review `web/src/components/search/SearchView.tsx` to confirm it is only a compatibility re-export.
- Review `web/src/App.tsx` for `AppShell` adoption.
- Review `web/src/components/organisms/SearchControlsPanel/SearchControlsPanel.tsx` for primitive adoption.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 12: Adopt Design-System Molecules in Corpus Surfaces

I extended the adoption pass from Search into the Corpus Explorer supporting components. Source selection, document browsing, and document inspection now use the reusable design-system primitives and molecules instead of bespoke panel/table/tab/metadata structures.

This is a broader proof point for the design system: the same `Panel`, `ScrollRegion`, `Caption`, `StatusText`, `DataTable`, `MetadataGrid`, and `TabList` components now support both Search Workbench and Corpus exploration workflows.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue the broad design-system extraction and component/page refactor work from the previous step.

**Inferred user intent:** The user wants sustained progress applying the new reusable primitives beyond Search and into the rest of the dashboard.

**Commit (code):** N/A at time of writing — implementation changes pending commit.

### What I did

- Refactored `SourcePanel` to use:
  - `Panel`
  - `ScrollRegion`
  - `Caption`
  - `SourcePanel.module.css`
- Refactored `DocumentBrowser` to use:
  - `Panel`
  - `ScrollRegion`
  - `DataTable`
  - `Caption`
  - `StatusText`
  - `DocumentBrowser.module.css`
- Refactored `DocumentInspector` to use:
  - `TabList`
  - `ScrollRegion`
  - `MetadataGrid`
  - `DataTable`
  - `Panel`
  - `Caption`
  - `StatusText`
  - `DocumentInspector.module.css`
- Preserved the existing corpus behaviors:
  - source selection,
  - document filtering,
  - document selection,
  - chunk row selection,
  - chunk timeline selection,
  - missing coverage list,
  - artifact loading and workflow navigation.

### Why

- Corpus Explorer had the same repeated global classes and inline styles as Search: panel headers, data tables, metadata grids, tab bars, scroll regions, and status text.
- Refactoring Corpus after Search checks whether the primitives generalize across views.
- The design doc explicitly called out corpus inspector tabs as Phase 4/5 review targets.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed.
- The corpus components now share the same table/metadata/tab primitives as the Search Workbench.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before preparing the commit.
- A surprising unrelated `go.mod` diff appeared removing local `replace` directives. I reverted `go.mod` because it was unrelated to the React design-system work and would also disturb the known local `../scraper` setup.
- Workflow views are still largely unreworked and remain the largest remaining design-system adoption target.

### What I learned

- `DataTable` is reusable across retrieval results, document rows, chunk rows, missing chunks, and artifact rows without needing a more complex API yet.
- `MetadataGrid` maps cleanly to document overview and artifact detail metadata.
- `TabList` can replace both Search inspector tabs and Document inspector tabs with the same controlled-tab pattern.

### What was tricky to build

- The main tricky part was preserving chunk selection behavior while moving chunk rows to `DataTable`. I kept selection keyed by chunk ID but converted it back to the existing selected index state so `ChunkTimelineBar` remains unchanged.
- Artifact detail had a local selected-index state that was not previously connected to row selection. I preserved its existing behavior without introducing new artifact-row selection semantics in this pass.

### What warrants a second pair of eyes

- Whether `DocumentInspector` should be split into smaller organisms/molecules now that it has clearer sections.
- Whether `DataTable` should gain optional wrapping cells for long IDs/text before Workflow adoption.
- Whether the corpus coverage strip should eventually share an extracted `CoverageStrip` with `ResultInspectorPanel`.

### What should be done in the future

- Refactor `WorkflowsView` next; it still contains many local panels, tables, status classes, and metadata grids.
- Consider extracting `DocumentInspector` section components if it remains large after the design-system adoption.
- Add Storybook stories for Corpus source/document/inspector states.

### Code review instructions

- Review `web/src/components/corpus/SourcePanel.tsx` for source list primitive adoption.
- Review `web/src/components/corpus/DocumentBrowser.tsx` for `DataTable` adoption.
- Review `web/src/components/corpus/DocumentInspector.tsx` for `TabList`, `MetadataGrid`, and `DataTable` adoption.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 13: Add Corpus Storybook Coverage and Re-check Foundation Guidance

I added Storybook stories for the Corpus components that were just moved onto the reusable design-system primitives. I also read the TTC foundation implementation guide and used it to recalibrate the RAG work: primitives should own typography/accessibility/separators/layout mechanics, but component-specific anatomy should remain local and CSS Modules should not be eliminated wholesale.

This step focuses on review surfaces. The Corpus refactor is now visible in Storybook with realistic source, document, and inspector states instead of relying only on the live app.

### Prompt Context

**User prompt (verbatim):** "also read /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--ttc-design-system/ttmp/2026/06/01/TTC-FOUNDATION-SYSTEM--ttc-react-foundation-primitives-and-token-documentation/design-doc/01-react-foundation-system-implementation-guide.md for guidelines on how to extract foundation stuff and clean up the css"

**Assistant interpretation:** Read the TTC foundation guide and apply its design-system extraction rules while continuing the RAG component/story work.

**Inferred user intent:** The user wants the RAG extraction to follow the proven TTC foundation approach: small targeted primitives, Storybook documentation, CSS reduction without over-abstracting, and local anatomy CSS where appropriate.

**Commit (code):** N/A at time of writing — Storybook coverage pending commit.

### What I did

- Read the TTC foundation guide at:
  - `/home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--ttc-design-system/ttmp/2026/06/01/TTC-FOUNDATION-SYSTEM--ttc-react-foundation-primitives-and-token-documentation/design-doc/01-react-foundation-system-implementation-guide.md`
- Added Corpus stories:
  - `web/src/components/corpus/SourcePanel.stories.tsx`
  - `web/src/components/corpus/DocumentBrowser.stories.tsx`
  - `web/src/components/corpus/DocumentInspector.stories.tsx`
- Covered realistic states:
  - source list interactive selection/loading,
  - document browser interactive selection/no-source/loading,
  - document inspector overview/highlighted chunk/no content.
- Ran validation:
  - `pnpm typecheck`
  - `pnpm build-storybook`
  - `pnpm build`

### Why

- The user explicitly asked to build many Storybook stories for the extracted components.
- The TTC guide says Storybook should become the review surface before broad CSS cleanup.
- Corpus components are now using reusable primitives, so they need stories to make those states inspectable.

### What worked

- `pnpm typecheck` passed.
- `pnpm build-storybook` passed and now includes Corpus component stories.
- `pnpm build` passed.

### What didn't work

- Storybook produces the existing large iframe chunk warning. This is not new and does not block the current extraction work.
- `pnpm build` again rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.

### What I learned

- The TTC guide reinforces the current RAG direction: avoid a generic Box primitive, keep primitive APIs small, and only refactor where the primitive clarifies ownership.
- Local CSS modules such as `DocumentInspector.module.css` and `SearchControlsPanel.module.css` are appropriate because they own component-specific anatomy, not reusable foundation concerns.
- Story coverage should accompany each extraction increment, not trail behind it.

### What was tricky to build

- `DocumentInspector` depends on RTK Query for artifacts when the artifacts tab is active, so the initial stories focus on overview/chunk/no-content states that do not require mocking the API store.
- This suggests that richer document artifact stories may need an API mocking pattern later, but that should be done carefully rather than blocking current component coverage.

### What warrants a second pair of eyes

- Whether Corpus stories should move under `Applications/Corpus` or stay under `Component Library/Corpus`.
- Whether we should add an MSW/RTK Query mocking setup before writing stories for API-active tabs.
- Whether `DocumentInspector` should expose a `defaultTab` prop like `ResultInspectorPanel` to make chunk/coverage/artifact stories more explicit.

### What should be done in the future

- Add Storybook coverage for Workflow list/detail/op states as workflow extraction proceeds.
- Add `defaultTab` support to `DocumentInspector` if we need explicit coverage/text/artifacts stories without play functions.
- Continue CSS cleanup by moving repeated foundation/layout concerns into primitives while keeping component anatomy local.

### Code review instructions

- Review the new Corpus stories under `web/src/components/corpus/`.
- Confirm the stories use realistic DTO-shaped data and no live API assumptions.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build-storybook`
  - `cd web && pnpm build`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build-storybook
pnpm build
```

## Step 14: Reclassify Corpus Source Selection and Add Workflow Organism Stories

I corrected the component taxonomy after the user pointed out that the source selector is an organism, not a corpus-local helper. The source selector owns a meaningful dashboard region with data summaries, selection state, preprocessing coverage, and scroll behavior, so it now lives under `components/organisms/CorpusSourcePanel/` with a compatibility re-export at the old corpus path.

I also started the Workflow extraction in the same style: dense workflow surfaces are being moved into organism-level panels with Storybook stories. `WorkflowsView` now delegates queue health and workflow list rendering to reusable organisms while still owning API polling and high-level page state.

### Prompt Context

**User prompt (verbatim):** "sourcepanel is an organism, no?"

**Assistant interpretation:** Correct the component hierarchy so SourcePanel is modeled and organized as an organism, and continue rigorous Storybook/design-system organization.

**Inferred user intent:** The user wants taxonomy discipline: components should live in folders that match their actual design-system level, not where they happened to originate.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Moved source selection into an organism:
  - `web/src/components/organisms/CorpusSourcePanel/CorpusSourcePanel.tsx`
  - `web/src/components/organisms/CorpusSourcePanel/CorpusSourcePanel.module.css`
  - `web/src/components/organisms/CorpusSourcePanel/CorpusSourcePanel.stories.tsx`
  - `web/src/components/organisms/CorpusSourcePanel/index.ts`
- Kept `web/src/components/corpus/SourcePanel.tsx` as a compatibility re-export.
- Updated Storybook title from a corpus-local grouping to:
  - `Component Library/Organisms/CorpusSourcePanel`
- Added workflow organism components and stories:
  - `web/src/components/organisms/QueueHealthPanel/`
  - `web/src/components/organisms/WorkflowListPanel/`
- Added `web/src/components/workflows/workflowFormat.ts` for workflow status/time formatting helpers shared by workflow organisms and the legacy workflow view.
- Refactored `WorkflowsView` so API polling remains in the page/view layer while `QueueHealthPanel` and `WorkflowListPanel` own reusable presentation.
- Updated `DocumentInspector.stories.tsx` to wrap stories in the Redux `Provider`, because the component uses RTK Query hooks even when most stories do not actively fetch artifact data.

### Why

- Per the TTC foundation guidance, primitives own foundation/layout concerns, but feature-level dashboard regions should be organisms.
- Corpus source selection is a region-level component, not a molecule.
- Workflow queue health and workflow list are also region-level organisms and need Storybook review surfaces.
- Storybook stories that render RTK Query components must provide the Redux store to avoid runtime provider errors.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed.
- Storybook now includes Workflow organism stories and the corrected `CorpusSourcePanel` organism story.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.
- Workflow detail/op inspector surfaces remain partially legacy and still need a follow-up extraction pass.

### What I learned

- Component taxonomy is not just folder hygiene; it controls Storybook navigation and future DMETA/Web MDS naming.
- API-owning views can remain containers while organism components receive DTO-shaped data and callbacks, making them easy to story.
- RTK Query hooks in components require explicit Provider setup in Storybook unless the hook is moved out to a container.

### What was tricky to build

- The main tricky part was preserving existing imports of `SourcePanel` while moving the implementation. A compatibility re-export avoids a broad rename blast radius.
- Another subtlety was avoiding a global Storybook provider for every story just to support `DocumentInspector`; adding a local decorator keeps the dependency explicit for the component that needs it.

### What warrants a second pair of eyes

- Whether `DocumentInspector` should eventually become an organism folder too, with the current corpus path as a re-export.
- Whether workflow detail should be split into `WorkflowDetailPanel`, `WorkflowOpGroupsPanel`, and `WorkflowOpInspectorPanel` organisms.
- Whether Storybook titles should use `Component Library/Organisms/...` consistently for all organisms, including the earlier Search Workbench panels.

### What should be done in the future

- Extract workflow detail/op inspector surfaces next.
- Add stories for workflow detail states, op graph states, and op result states.
- Consider a shared Storybook decorator for RTK Query/API-backed stories if more components need providers.

### Code review instructions

- Review `web/src/components/organisms/CorpusSourcePanel/` for the source selector reclassification.
- Review `web/src/components/corpus/SourcePanel.tsx` to confirm it is only a compatibility re-export.
- Review `web/src/components/organisms/QueueHealthPanel/` and `WorkflowListPanel/` for workflow organism extraction.
- Review `web/src/components/corpus/DocumentInspector.stories.tsx` for the local Redux provider decorator.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 15: Move SourcePanel Back and Extract Workflow Detail Panels

I moved `SourcePanel` back into the Corpus feature area after the user clarified the desired organization. The important design-system ownership remains intact: `SourcePanel` still uses shared foundation/layout primitives, but its source-specific feature identity now lives next to the other Corpus components.

I then continued the Workflow extraction by pulling the workflow detail header/progress, op graph, and op groups table into Storybook-backed organism components. `WorkflowsView` still owns polling and workflow-detail orchestration, while the new organisms own reusable presentation.

### Prompt Context

**User prompt (verbatim):** "ok move corpus source panel back i guess. continue"

**Assistant interpretation:** Move the source panel implementation back to the Corpus folder and continue the rigorous design-system extraction work.

**Inferred user intent:** The user wants practical organization over rigid taxonomy: feature-local Corpus components may stay grouped together as long as foundation/layout ownership is clean and Storybook remains organized.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Moved `CorpusSourcePanel` back to:
  - `web/src/components/corpus/SourcePanel.tsx`
  - `web/src/components/corpus/SourcePanel.module.css`
  - `web/src/components/corpus/SourcePanel.stories.tsx`
- Removed the organism-level `CorpusSourcePanel` folder and export.
- Restored the story title to `Component Library/Corpus/SourcePanel`.
- Added workflow detail organisms:
  - `WorkflowSummaryPanel`
  - `WorkflowOpGraphPanel`
  - `WorkflowOpGroupsPanel`
- Added Storybook stories for the new workflow detail organisms.
- Refactored `WorkflowsView` to use the new workflow detail organisms.
- Kept typography/status ownership in foundation primitives (`Caption`, `StatusText`) and layout/table ownership in layout/molecule primitives (`Panel`, `DataTable`, `MetadataGrid`).

### Why

- The user explicitly asked to move the source panel back.
- The TTC foundation guidance says not to over-abstract or reorganize for its own sake; keep component-specific anatomy local and use primitives for shared foundation/layout concerns.
- Workflow detail had repeated progress, op graph, and table structures that were better as Storybook-reviewable organisms.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed.
- Storybook now includes workflow detail organism stories in addition to queue/list workflow stories.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.
- Workflow op inspector/result detail is still partly legacy and should be the next extraction target.

### What I learned

- The source panel can be feature-local without violating design-system ownership. Folder location and design-system level are related, but not identical.
- Workflow detail is now much easier to reason about because progress/header, graph, and group table are separate components with stories.

### What was tricky to build

- The tricky part was preserving the existing workflow selection semantics. `WorkflowOpGroupsPanel` emits selected groups, and `WorkflowsView` still maps that to the stable `operation|status` key used to find the sampled op.
- Another tricky part was avoiding foundation leakage: workflow-specific progress bar and graph anatomy live in workflow organism CSS modules, while typography/status text uses foundation primitives.

### What warrants a second pair of eyes

- Whether `WorkflowSummaryPanel` should expose the progress bar as a separate molecule later.
- Whether workflow graph styling should remain organism-specific or become a reusable graph primitive after more graph-like views appear.
- Whether `SourcePanel` should remain in `components/corpus` permanently or move later if a broader organism taxonomy is reintroduced.

### What should be done in the future

- Extract `WorkflowOpInspectorPanel` and `WorkflowOpResultSection` with stories.
- Add API-mocked stories for artifact/result loading states if needed.
- Keep Storybook titles consistent by feature area unless a component is genuinely generic.

### Code review instructions

- Confirm `web/src/components/corpus/SourcePanel.tsx` is the real source-panel implementation again.
- Review `web/src/components/organisms/WorkflowSummaryPanel/`, `WorkflowOpGraphPanel/`, and `WorkflowOpGroupsPanel/`.
- Review `web/src/components/workflows/WorkflowsView.tsx` for the new workflow-detail composition.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 16: Extract the Workflow Op Inspector Panel

I extracted the sampled workflow operation inspector into a dedicated organism. The new `WorkflowOpInspectorPanel` owns the selected operation's identity, status, queue, dedup key, input metadata, and retry affordance, while `WorkflowsView` still owns the polled data, selected group state, retry mutation, and op result fetching.

This keeps the design-system ownership consistent: reusable foundation primitives own caption/status text, `MetadataGrid` owns key/value structure, `Panel` owns the shell, and the workflow-specific inspector CSS owns only local anatomy such as the accent rail and section spacing.

### Prompt Context

**User prompt (verbatim):** (same as Step 15)

**Assistant interpretation:** Continue extracting workflow detail surfaces into organized, Storybook-backed components while keeping foundation/layout ownership rigorous.

**Inferred user intent:** The user wants the workflow dashboard to receive the same component extraction and Storybook discipline as Search and Corpus.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Added `web/src/components/organisms/WorkflowOpInspectorPanel/`.
- Added `WorkflowOpInspectorPanel.tsx`, `.module.css`, `index.ts`, and Storybook stories.
- Refactored `WorkflowsView` to render `WorkflowOpInspectorPanel` for the selected sampled op.
- Kept `OpResultSection` in `WorkflowsView` for now and passed it as children to the inspector panel.
- Removed now-unused local status/op formatting helpers from `WorkflowsView` after moving that presentation into organisms.

### Why

- The sampled op inspector was one of the last large workflow-detail regions still using ad hoc panel/metadata/status CSS directly in the page file.
- Extracting it makes the next step smaller: only the op result section remains as the main workflow-detail presentation target.
- Storybook can now review succeeded and failed/retryable sampled op states without live API calls.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed.
- The new story uses DTO-shaped `WorkflowOp` data and does not need a Redux provider because it is presentation-only.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.
- `OpResultSection` still uses RTK Query directly in `WorkflowsView`, so it remains a follow-up extraction target.

### What I learned

- Passing the op result section as `children` is a useful incremental boundary: the inspector can become storyable immediately without moving the data-fetching hook.
- Workflow detail is now mostly composed of named organisms, which makes the remaining legacy pieces easier to identify.

### What was tricky to build

- The tricky part was preserving retry behavior. `WorkflowsView` still looks up the selected group by the stable key and calls the retry mutation; `WorkflowOpInspectorPanel` only renders the retry button and invokes the callback.
- Another subtlety was avoiding over-generalization: the accent rail and inspector spacing stay in an organism CSS module rather than becoming new primitives.

### What warrants a second pair of eyes

- Whether `OpResultSection` should become `WorkflowOpResultPanel` or a smaller molecule used inside `WorkflowOpInspectorPanel`.
- Whether retry controls should remain inside the inspector organism or become an action slot supplied by the container.

### What should be done in the future

- Extract `WorkflowOpResultSection` and add stories for no-result, records, artifacts, emitted ops, and error states.
- Consider replacing remaining workflow form rows in `SubmitIntakeModal` with `FormRow` after op result extraction.

### Code review instructions

- Review `web/src/components/organisms/WorkflowOpInspectorPanel/WorkflowOpInspectorPanel.tsx`.
- Review the selected-op block in `web/src/components/workflows/WorkflowsView.tsx` to confirm behavior is preserved.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 17: Extract the Workflow Op Result Panel

I extracted the workflow op result renderer into a Storybook-backed `WorkflowOpResultPanel`. The container hook remains in `WorkflowsView`, but the result presentation now owns its reusable metadata, records, artifacts, emitted op, loading, no-result, and error states.

This finishes the main workflow-detail extraction target identified in the previous step. Workflow details now compose named organisms for summary/progress, op graph, op groups, sampled op inspection, and op result rendering.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue the workflow extraction, focusing on the remaining legacy op result surface.

**Inferred user intent:** The user wants the workflow view to keep moving toward organized, Storybook-reviewable design-system components.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Added `web/src/components/organisms/WorkflowOpResultPanel/`.
- Implemented result states for:
  - loading,
  - 404/no result,
  - completed metadata,
  - JSON data,
  - records written via `DataTable`,
  - artifacts via `DataTable`,
  - emitted op IDs,
  - result errors.
- Added Storybook stories for:
  - complete result,
  - loading,
  - no result,
  - error result.
- Refactored `OpResultSection` in `WorkflowsView` to fetch data and delegate rendering to `WorkflowOpResultPanel`.
- Updated tasks to mark op result extraction and Storybook coverage complete.

### Why

- `OpResultSection` was the last large workflow-detail presentation block still rendering bespoke metadata grids and tables directly inside `WorkflowsView`.
- Keeping the RTK Query hook in `WorkflowsView` preserves the container/presentation split and keeps the new result panel easy to story with plain DTO-shaped props.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed.
- Storybook now includes `WorkflowOpResultPanel` review states.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.

### What I learned

- Workflow details are now mostly container/presentation split: `WorkflowsView` owns API calls and mutations, while extracted organisms own rendering.
- `DataTable` continues to work for dense operational data without adding custom table variants.

### What was tricky to build

- The tricky part was preserving the special 404 behavior. A missing result is not a hard failure; it means no result data was recorded yet. `WorkflowOpResultPanel` now models that as a normal no-result state.
- Another subtle point was not moving the RTK Query hook into the result panel, because doing that would make Storybook harder and mix data fetching into the reusable organism.

### What warrants a second pair of eyes

- Whether `WorkflowOpResultPanel` should be a molecule rather than organism. It is currently an organism because it owns a multi-section operational panel, but it is nested inside another organism.
- Whether the JSON data block needs a shared `CodeBlock` foundation/utility component after more JSON views are found.

### What should be done in the future

- Refactor `SubmitIntakeModal` with `FormRow`, `Stack`, `Panel`, and `Caption`.
- Add a workflow page-level story once RTK Query mocking is standardized.
- Consider splitting `WorkflowsView` into a `WorkflowPage` container folder after modal cleanup.

### Code review instructions

- Review `web/src/components/organisms/WorkflowOpResultPanel/WorkflowOpResultPanel.tsx`.
- Review `OpResultSection` in `web/src/components/workflows/WorkflowsView.tsx` to confirm it only fetches and delegates rendering.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 18: Refactor the Intake Modal Form Through Layout/Foundation Primitives

I continued the workflow cleanup by refactoring the intake workflow modal's form structure through the shared design-system primitives. The modal still lives inside `WorkflowsView` because it owns mutation state and submission behavior, but its field layout now uses `Panel`, `Stack`, `FormRow`, and `Caption` instead of repeated `form-row` divs and ad hoc text classes.

This follows the TTC foundation guidance: keep behavior and feature-specific form grouping local, but move repeated foundation and layout concerns into primitives. The result is not a new abstraction for every form field; it is a cleaner use of the already-extracted layout vocabulary.

### Prompt Context

**User prompt (verbatim):** (same as Step 17)

**Assistant interpretation:** Continue cleaning remaining workflow UI surfaces after extracting the op result panel.

**Inferred user intent:** The user wants the workflow dashboard to keep reducing repeated CSS/layout patterns while preserving behavior.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Refactored `SubmitIntakeModal` inside `web/src/components/workflows/WorkflowsView.tsx` to use:
  - `Panel` for the modal shell,
  - `Stack` for vertical field grouping,
  - `FormRow` for label/control rows,
  - `Caption` for strategy and available-source metadata.
- Preserved existing behavior:
  - source input parsing,
  - strategy id derivation,
  - chunking settings,
  - embedding toggles/settings,
  - BM25 toggle,
  - error display,
  - submit/cancel actions.

### Why

- The modal had many repeated form-row/layout classes that overlapped with the new layout primitives.
- This was the main remaining workflow form structure after the list/detail/op panels were extracted.
- Keeping the mutation logic in the view avoids turning a presentational form into an API-aware component prematurely.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed.

### What didn't work

- The first edit accidentally left a duplicate `);` after the modal return block, causing TypeScript parse errors:
  - `TS1128: Declaration or statement expected.`
- I removed the extra closing expression and reran validation successfully.
- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.

### What I learned

- `FormRow` is useful for workflow forms without needing a large form framework.
- It is still appropriate for `WorkflowsView` to own modal submission state until/unless a storyable presentational `SubmitIntakeWorkflowForm` becomes necessary.

### What was tricky to build

- The tricky part was replacing a large JSX block without changing form state wiring. Each control still calls the same `set` helper or `setSourceInput` callback.
- Another subtlety was not over-extracting: fieldsets remain local because their grouping is feature-specific, while row layout and caption text move to shared primitives.

### What warrants a second pair of eyes

- Whether `SubmitIntakeModal` should become its own component once API mocking/story infrastructure is added.
- Whether `FormRow` needs a compact/numeric-width convention before more forms adopt it.

### What should be done in the future

- Add a story for a presentational intake form only if the form is split from RTK Query mutation state.
- Continue scanning Workflow/Cross-page CSS for remaining repeated foundation concerns.

### Code review instructions

- Review `SubmitIntakeModal` in `web/src/components/workflows/WorkflowsView.tsx`.
- Confirm behavior is unchanged and only form layout/foundation ownership changed.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 19: Clean Up Search Support Molecules Against Foundation Ownership

I audited remaining global panel/table/text class usage after the workflow extraction and cleaned up two Search support molecules that still had old global structure embedded in them. `CoveragePanel` and `QueryPresetList` now use the shared design-system primitives instead of rendering their own global `panel`, `panel-header`, `data-table`, `meta-grid`, and `text-mono` structure.

This is a small but important CSS-ownership pass. It follows the TTC foundation guidance by moving typography/status/table/metadata ownership into foundation/layout/molecule primitives while keeping only component-specific anatomy in local CSS modules.

### Prompt Context

**User prompt (verbatim):** "go ahead"

**Assistant interpretation:** Proceed with the next cleanup step: audit remaining global CSS usage and reduce foundation/layout debt where it is clearly safe.

**Inferred user intent:** The user wants the design-system extraction to continue beyond large page splits into rigorous CSS ownership cleanup.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Refactored `CoveragePanel` to use:
  - `Panel`
  - `StatusText`
  - `Caption`
  - `MetadataGrid`
  - `DataTable`
- Added `CoveragePanel.module.css` for the remaining sparse-coverage warning anatomy.
- Refactored `QueryPresetList` to use:
  - `Panel`
  - `Stack`
  - `Caption`
- Added `QueryPresetList.module.css` for query button hover/click anatomy.
- Ran a scan of remaining global panel/table/text class usage to identify next cleanup candidates.

### Why

- These molecules were already extracted but still carried old global classes internally.
- They are low-risk because they are small and covered by existing Storybook stories.
- The design-system contract should be used consistently inside new/recently extracted components.

### What worked

- `pnpm typecheck` passed.
- `pnpm build-storybook` passed.
- `pnpm build` passed.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.

### What I learned

- The remaining global CSS debt is now concentrated in older page/view areas (`EmbeddingsView`, `PipelineView`, some corpus identity/artifact bars, and Storybook demo snippets), rather than the main Search/Corpus/Workflow extracted panels.
- Small molecules are good places to enforce the design-system rules because they should not need bespoke panel/table/metadata implementations.

### What was tricky to build

- The only tricky point was avoiding circular barrel imports from inside `CoveragePanel`; I imported `DataTable` and `MetadataGrid` directly from their sibling folders instead of from the molecule barrel.

### What warrants a second pair of eyes

- Whether `CoveragePanel` should eventually use a shared `CoverageStrip` if more visual coverage strips are normalized.
- Whether `QueryPresetList` should be an organism if presets become stateful or source-aware; it is currently still a small molecule.

### What should be done in the future

- Continue global CSS cleanup in `EmbeddingsView`, `PipelineView`, `IdentityBar`, and `ArtifactIdentityBar`.
- Replace story demo snippets that use raw global panel classes with actual `Panel` primitives.
- Consider adding design-system overview stories for token/foundation usage once component cleanup stabilizes.

### Code review instructions

- Review `web/src/components/molecules/CoveragePanel/CoveragePanel.tsx` and `.module.css`.
- Review `web/src/components/molecules/QueryPresetList/QueryPresetList.tsx` and `.module.css`.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build-storybook`
  - `cd web && pnpm build`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build-storybook
pnpm build
```

## Step 20: Refactor the Embeddings View Through Shared Primitives

I continued the global CSS ownership cleanup by targeting `EmbeddingsView`, which was one of the largest remaining view-level users of raw `panel`, `panel-header`, `panel-body`, `data-table`, `text-dim`, and inline flex/grid structures. The view still owns its RTK Query hooks, local form state, and mutation handlers, but its visible layout now composes `Panel`, `DashboardGrid`, `Stack`, `FormRow`, `Caption`, `CodeText`, `MetadataGrid`, and `DataTable`.

This keeps the refactor inside the current scope: it improves design-system adoption without creating generated DMETA scaffolds, validators, or premature presentational extraction for API-heavy controls.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue the next design-system cleanup pass after the Search support molecule cleanup.

**Inferred user intent:** The user wants systematic reduction of remaining global CSS debt and stronger reuse of the shared primitives across older dashboard views.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Rewrote `web/src/components/embeddings/EmbeddingsView.tsx` to use shared primitives:
  - `Panel` for all major cards,
  - `DashboardGrid` for compute/similarity two-column layout,
  - `Stack` for vertical spacing,
  - `FormRow` for labeled controls,
  - `Caption` for explanatory and muted text,
  - `CodeText` for chunk/model identifiers,
  - `MetadataGrid` for snapshot and compute result facts,
  - `DataTable` for similarity matches.
- Added `web/src/components/embeddings/EmbeddingsView.module.css` for view-specific anatomy that is not a general primitive:
  - controls grid,
  - fieldset styling,
  - short numeric inputs,
  - checkbox row,
  - action/table spacing.
- Preserved all existing data/query/mutation behavior.

### Why

- `EmbeddingsView` was a high-value cleanup target because it contained several repeated global class patterns in one file.
- It is still API-heavy, so extracting a Storybook-friendly organism should wait until there is a clearer mocking/provider pattern for page-level views.
- Moving to primitives now improves consistency while keeping behavioral ownership stable.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.
- Storybook still emits the known large iframe chunk warning; build succeeds.

### What I learned

- The existing primitives cover most of the embeddings inspector without adding a new generic `Box` or page-specific layout component.
- `MetadataGrid` is a better fit than an ad hoc two-column result table for compute summaries.

### What was tricky to build

- The main sharp edge was preserving the dependent state behavior for chunks while replacing large JSX blocks. The effects that default `strategyId`, `documentId`, `chunkIDA`, and `chunkIDB` were left unchanged.
- Another subtle point was keeping feature-specific fieldsets local: they describe this view's information grouping, not a reusable design-system primitive.

### What warrants a second pair of eyes

- Confirm the visual spacing of the controls grid after replacing inline styles with the module class.
- Confirm whether `EmbeddingsView` should later split into presentational panels if/when API mocking patterns are standardized for Storybook.

### What should be done in the future

- Continue cleanup in `PipelineView`, `IdentityBar`, and `ArtifactIdentityBar`.
- Consider page-level stories only after standardized Redux/RTK Query mocking is available.

### Code review instructions

- Start with `web/src/components/embeddings/EmbeddingsView.tsx` and compare the preserved hooks/effects/handlers against the previous behavior.
- Review `web/src/components/embeddings/EmbeddingsView.module.css` for view-only CSS ownership.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 21: Add Storyable Pipeline Boundaries and Organize Corpus Widgets

I continued the cleanup by converting the lightweight pipeline dashboard from raw global classes into a storyable presentation component and an explicit page boundary. I also applied the user's organizational rule to Corpus widgets: each public widget now lives in its own subdirectory with its component, CSS module, and Storybook story co-located where applicable.

This step also records a strengthened process rule: every reusable extraction or page boundary added from here must get a Storybook story in the same change. For API-heavy views, the pattern is to split a DTO-shaped presentational component/page from the RTK Query container and story the presentational boundary.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue reducing old global CSS usage and improve design-system/component organization.

**Inferred user intent:** The user wants the refactor to keep progressing systematically through remaining older views.

**User prompt (verbatim):** "add storybook stories as you go, never forget that"

**Assistant interpretation:** Treat Storybook stories as mandatory for every reusable component/page extraction or reorganization step.

**Inferred user intent:** The user wants Storybook to remain the review surface and contract for design-system work, not an afterthought.

**User prompt (verbatim):** "for pages too"

**Assistant interpretation:** Page-level boundaries also need Storybook stories, not only atoms/molecules/organisms.

**Inferred user intent:** The user wants page composition states to be reviewable in Storybook even when the runtime container uses RTK Query.

**User prompt (verbatim):** "organize corpus widgets in subdirectories per widget as well."

**Assistant interpretation:** Move Corpus widgets from a flat directory into per-widget folders, matching the reusable component organization used elsewhere.

**Inferred user intent:** The user wants consistent file organization across component areas so stories, styles, and component code stay co-located.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Added `web/src/components/pipeline/PipelineOverview.tsx` as a DTO-shaped presentational component.
- Added `web/src/components/pipeline/PipelineOverview.stories.tsx` with populated/loading/empty states.
- Added `web/src/components/pages/PipelinePage/PipelinePage.tsx` as an explicit page boundary.
- Added `web/src/components/pages/PipelinePage/PipelinePage.stories.tsx` with populated/loading/empty states.
- Updated `PipelineView` to remain the RTK Query container and delegate rendering to `PipelinePage`.
- Moved Corpus widgets into per-widget directories:
  - `ArtifactIdentityBar/ArtifactIdentityBar.tsx`
  - `ChunkTimelineBar/ChunkTimelineBar.tsx`
  - `DocumentBrowser/DocumentBrowser.tsx`, `.module.css`, `.stories.tsx`
  - `DocumentInspector/DocumentInspector.tsx`, `.module.css`, `.stories.tsx`
  - `IdentityBar/IdentityBar.tsx`
  - `SourcePanel/SourcePanel.tsx`, `.module.css`, `.stories.tsx`
- Updated imports after the directory move.

### Why

- `PipelineView` was still using raw `panel`, `data-table`, status, and text classes directly.
- The page-story requirement means API-heavy containers should delegate to storyable presentational/page boundaries.
- Corpus had already grown enough widgets that a flat folder made ownership harder to scan.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed, including the new Pipeline overview/page stories and the moved Corpus stories.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.
- Storybook still emits the known large iframe chunk warning; build succeeds.

### What I learned

- The container/page/presentational split works well for page stories without introducing API mocking yet.
- Corpus story files survived the subdirectory move cleanly once service/store import paths were adjusted.

### What was tricky to build

- Moving Corpus widgets required careful relative-import updates because each component moved one directory deeper.
- The page story needed to avoid live RTK Query dependencies, so `PipelinePage` receives DTO-shaped props and `PipelineView` owns fetching.

### What warrants a second pair of eyes

- Whether `PipelineOverview` should live under `components/pipeline` long-term or under `components/organisms/PipelineOverview` if it becomes reusable outside the page.
- Whether the Corpus identity bars should also receive Storybook stories and primitive cleanup in a follow-up step.

### What should be done in the future

- Add stories for Corpus identity/artifact bars when they are cleaned up.
- Continue applying the story-required rule to all future page and component boundaries.
- Consider adding a SearchWorkbench page story once RTK Query mocking or a presentational page split exists.

### Code review instructions

- Review `web/src/components/pipeline/PipelineOverview.tsx` and its story first.
- Review `web/src/components/pages/PipelinePage/PipelinePage.tsx` and its story next.
- Review `web/src/components/corpus/CorpusExplorerView.tsx` import changes and the moved Corpus widget directories.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 22: Write RAG React Design-System Guidelines

I paused implementation and converted the TTC foundation-system guidance into a RAG-specific design-system guideline document. The new guide keeps the RAG dashboard's retro monochrome identity, but raises the standard for ownership, file layout, Storybook coverage, page boundaries, and CSS cleanup so this project can move closer to the TTC design-system discipline.

This is intentionally a project guideline, not more code. The immediate value is that future refactors now have a local RAG rulebook instead of relying on memory or on a TTC document written for a different product surface.

### Prompt Context

**User prompt (verbatim):** "now read this and go back threough /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--ttc-design-system/ttmp/2026/06/01/TTC-FOUNDATION-SYSTEM--ttc-react-foundation-primitives-and-token-documentation/design-doc/01-react-foundation-system-implementation-guide.md and create a similar set of guidelines for our project, we still have too much ad-hoc css and for example no full fledged design system like ttc-design-system, we should be much closer"

**Assistant interpretation:** Read the TTC foundation guide and author an equivalent RAG-local React/design-system guideline document that addresses remaining ad-hoc CSS and missing design-system rigor.

**Inferred user intent:** The user wants the RAG project to stop accumulating piecemeal cleanup and instead follow explicit, reusable design-system rules comparable to TTC.

**User prompt (verbatim):** "create your own set of guidelines for react and design system and guidelines and file layouts and store in the ticket"

**Assistant interpretation:** Store a new RAG-specific design-system guideline document inside the current docmgr ticket.

**Inferred user intent:** The user wants this guidance persisted as ticket documentation, not just summarized in chat.

**Commit (code):** N/A at time of writing — documentation-only changes pending commit.

### What I did

- Read the TTC foundation-system implementation guide.
- Created a new ticket design doc:
  - `design-doc/02-rag-react-design-system-guidelines.md`
- Wrote RAG-specific rules for:
  - target React/design-system stack,
  - CSS ownership,
  - file/folder layout,
  - Storybook requirements for components and pages,
  - container vs presentational boundaries,
  - foundation/layout/molecule/organism/page responsibilities,
  - DMETA boundaries,
  - global CSS reduction,
  - validation and commit checklist,
  - recommended next work.
- Updated ticket tasks to mark the guideline work complete.

### Why

- The RAG project has improved primitives, but still lacks a full TTC-like foundation documentation surface and clear local rules.
- The user explicitly called out too much ad-hoc CSS and asked for project-local guidelines.
- A written guideline prevents future cleanup from becoming subjective or inconsistent.

### What worked

- `docmgr doc add` created the design doc under the existing RAG ticket.
- The TTC guide translated cleanly into a RAG-specific stack with `Text`, `CodeText`, `Caption`, `StatusText`, `Divider`, `VisuallyHidden`, `Panel`, `Stack`, `DashboardGrid`, `DataTable`, and `MetadataGrid` as the current foundation.

### What didn't work

- N/A; this was documentation authoring only.

### What I learned

- RAG should not blindly mirror TTC's editorial primitives like `Heading` and `Eyebrow`; it should add them only if repeated dashboard heading semantics require them.
- The biggest missing RAG piece is now a foundation overview Storybook documentation surface plus continued migration of legacy global CSS into primitives/modules.

### What was tricky to build

- The main tricky part was preserving TTC's discipline without importing TTC's product assumptions. RAG is a dense developer dashboard, not a commerce/editorial landing page, so the guideline emphasizes panels, tables, metadata, status, code text, API containers, and page stories.
- Another subtlety was keeping the DMETA boundary strict: the guidelines are React/design-system rules and do not authorize RAG IR validators or scaffold generation.

### What warrants a second pair of eyes

- Whether the recommended future primitives (`MetricGrid`, `ProgressBar`, `IdentitySelectorBar`, `CoverageStrip`) are framed narrowly enough to avoid premature abstraction.
- Whether Storybook title taxonomy should be adjusted before adding more page-level stories.

### What should be done in the future

- Add RAG foundation overview Storybook docs.
- Add stories and primitive cleanup for `IdentityBar`, `ArtifactIdentityBar`, and `ChunkTimelineBar`.
- Continue shrinking `web/src/index.css` only after `rg` confirms consumers are gone.

### Code review instructions

- Review `design-doc/02-rag-react-design-system-guidelines.md` first.
- Compare its rules against the TTC foundation guide and the current RAG component tree.
- Confirm it preserves these constraints:
  - no generic `Box`,
  - Storybook mandatory for components and pages,
  - CSS Modules remain for anatomy,
  - RAG DMETA remains documentation-only for now.

### Technical details

Key commands:

```bash
cd 2026-05-27--rag-evaluation-system
docmgr --root ttmp doc add --ticket RAG-WEB-DESIGN-SYSTEM-REVIEW --doc-type design-doc --title "RAG React Design System Guidelines"
```

## Step 23: Audit Components Against Guidelines and Add Foundation Overview Stories

I turned the new RAG React design-system guidelines into concrete ticket work. First I added detailed continuation tasks for auditing, foundation documentation, Corpus identity cleanup, global CSS cleanup, and page-story extraction. Then I completed the first two requested steps: a component-by-component guideline audit and the missing TTC-style foundation Storybook overview.

This gives the project a clearer next path. We now have both a rulebook and an audit that says which components already comply, which ones are partial, and what should be extracted after the foundation layer is visible in Storybook.

### Prompt Context

**User prompt (verbatim):** "alright, add detailed tasks to the ticket to continue the audit, add all the stuff above, and then let's do the full 1. and 2."

**Assistant interpretation:** Add detailed follow-up tasks to the ticket, then complete the guideline audit and foundation Storybook documentation steps.

**Inferred user intent:** The user wants the guidelines translated into actionable project management and wants immediate progress on auditing/foundation docs before more component cleanup.

**User prompt (verbatim):** "after doing the foundations (typography, etc...) , go through each components and figure out what should be extracted further"

**Assistant interpretation:** The audit should not stop at generic compliance; it should identify concrete component-by-component extraction opportunities after the foundation docs are in place.

**Inferred user intent:** The user wants a systematic walkthrough of all prior work to plan deeper cleanup, not isolated refactors.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Added detailed ticket tasks for:
  - guideline inventory/audit,
  - audit analysis doc,
  - foundation overview Storybook docs,
  - Corpus identity/timeline stories and cleanup,
  - global CSS cleanup,
  - page story boundaries.
- Created `analysis/01-rag-design-system-guideline-audit.md`.
- Audited each component/page for:
  - story coverage,
  - legacy global class usage,
  - inline style pressure,
  - layer compliance,
  - extraction opportunities.
- Added `web/src/components/foundation/Foundation.stories.tsx`.
- Added `web/src/components/foundation/Foundation.stories.module.css`.
- Foundation overview stories cover:
  - colors,
  - typography,
  - status tones,
  - spacing,
  - borders/radii,
  - accessibility.
- Checked tasks 56, 57, 58, and 59.

### Why

- The user explicitly asked for detailed ticket tasks and a full audit before continuing component cleanup.
- RAG had individual primitive stories but no TTC-style overview documentation surface for the design-system foundation.
- The audit is needed so future cleanup happens in a deliberate order rather than opportunistically.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed and included the new foundation overview stories.
- The audit clearly identified Corpus identity widgets, Search page story split, Workflows page split, Embeddings page split, and Evaluation placeholder cleanup as the next sequence.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.
- Storybook still emits the known large iframe chunk warning and plugin timing warning; build succeeds.

### What I learned

- The strongest layers are now foundation/layout/molecules plus most Search and Workflow organisms.
- The biggest design-system debt is concentrated in `IdentityBar`, `ArtifactIdentityBar`, `ChunkTimelineBar`, API-heavy page boundaries, and remaining global CSS blocks.
- Foundation overview stories make the missing design-system surface visible without adding new app-facing primitives.

### What was tricky to build

- The audit needed to distinguish pure RTK Query containers from presentational components. Containers do not always need stories if they delegate to a storyable page/presentation boundary.
- Foundation docs needed some story-only visual styling, but the guideline forbids creating a generic `Box` or app-facing demo primitive. I kept that styling in `Foundation.stories.module.css`.

### What warrants a second pair of eyes

- Whether the audit's proposed extraction order matches product priorities.
- Whether `SearchWorkbenchPage` should be split before Corpus identity cleanup, or vice versa.
- Whether the Storybook foundation taxonomy should be `Design System/Foundation/Overview` or split into separate title groups later.

### What should be done in the future

- Start the next cleanup with Corpus identity/timeline widgets and required stories.
- Then split API-heavy pages into storyable boundaries.
- Delete global CSS blocks only after `rg` confirms no remaining consumers.

### Code review instructions

- Review `analysis/01-rag-design-system-guideline-audit.md` for the component-by-component extraction plan.
- Review `web/src/components/foundation/Foundation.stories.tsx` and `.module.css` for the foundation overview documentation.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 24: Refactor Corpus IdentityBar Through Shared Primitives

I started the high-priority Corpus identity cleanup identified by the guideline audit. `IdentityBar` now uses the design-system primitives for its panel shell, wrapping selector layout, labeled controls, and totals caption instead of global `panel`/`text-*` classes and inline flex styles.

This is the first concrete post-audit cleanup pass. It keeps the component feature-local under Corpus, but gives it the same co-located folder/story/module shape as the rest of the reusable dashboard widgets.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue with the next recommended cleanup from the audit, starting with Corpus identity widgets.

**Inferred user intent:** The user wants the guideline audit to turn into implementation work, with Storybook coverage as we go.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Refactored `web/src/components/corpus/IdentityBar/IdentityBar.tsx` to use:
  - `Panel`,
  - `Inline`,
  - `FormRow`,
  - `Caption`.
- Added `IdentityBar.module.css` for component-specific sizing and totals placement.
- Added `IdentityBar.stories.tsx` with interactive, sparse coverage, and no-strategies states.
- Added `IdentityBar/index.ts`.

### Why

- `IdentityBar` was one of the audit's highest-priority non-compliant widgets.
- It used old global panel/text classes and many inline styles.
- It is important RAG UX: it controls the embedding identity under which Corpus data is inspected.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed and included the new IdentityBar stories.

### What didn't work

- The first Storybook story attempt omitted a required `onChange` arg in default meta args, causing TypeScript errors. I added a noop default and reran validation successfully.
- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.

### What I learned

- `FormRow` works for compact inline identity controls, but it may need future tuning if identity bars become denser.
- `IdentityBar` should stay feature-local for now; a generic `IdentitySelectorBar` should wait until Artifact and Embeddings selectors converge on the same shape.

### What was tricky to build

- The main subtlety was keeping the control semantics unchanged while replacing label/flex markup. `FormRow` is not a native `<label>`, so the existing control behavior is preserved but not made more semantically coupled yet.
- The no-strategies state is story-covered but still renders an empty select; this matches current behavior but may warrant UX improvement later.

### What warrants a second pair of eyes

- Whether the wrapped inline `FormRow` layout is compact enough for narrow dashboard widths.
- Whether we should add IDs/html labels to `FormRow` controls in a future accessibility pass.

### What should be done in the future

- Refactor `ArtifactIdentityBar` next and decide whether the two identity selector styles should share a small private subcomponent.
- Continue with `ChunkTimelineBar` after artifact identity selectors.

### Code review instructions

- Review `web/src/components/corpus/IdentityBar/IdentityBar.tsx` for unchanged identity update behavior.
- Review `web/src/components/corpus/IdentityBar/IdentityBar.stories.tsx` for the required Corpus identity states.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 25: Refactor Corpus Artifact Identity Selectors

I continued the Corpus identity cleanup by refactoring the document-processing and chunk-enrichment identity selector bars. These controls now use shared foundation/layout primitives plus a local CSS Module for selector button anatomy instead of raw `text-mono` classes, inline flex layouts, and global button styling.

This keeps the artifact identity controls feature-local, but makes their states visible in Storybook. The components are still in one file because they share a small style and conceptual role; a future split into separate folders is only warranted if they diverge or become reusable outside Corpus.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue the next high-priority post-audit cleanup after `IdentityBar`, namely `ArtifactIdentityBar`.

**Inferred user intent:** The user wants systematic application of the guidelines to Corpus identity controls, with stories as changes are made.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Refactored `web/src/components/corpus/ArtifactIdentityBar/ArtifactIdentityBar.tsx` to use:
  - `Inline`,
  - `Caption`,
  - `CodeText`.
- Removed direct `text-mono`/`text-dim` usage and global `btn`/`btn-primary` usage from these controls.
- Added `ArtifactIdentityBar.module.css` for local selector anatomy and selected state.
- Added `ArtifactIdentityBar.stories.tsx` covering:
  - document-processing selector,
  - chunk-enrichment selector,
  - combined stacked selectors,
  - empty preprocessing identities,
  - empty enrichment identities.
- Added `ArtifactIdentityBar/index.ts`.

### Why

- The guideline audit marked this file as a high-priority non-compliant area.
- These selectors define which document-processing and chunk-enrichment artifact identity the Corpus dashboard is currently inspecting.
- They need Storybook coverage because they encode important RAG state-selection behavior.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed and included the new artifact identity stories.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.
- Storybook still emits the known large iframe chunk warning; build succeeds.

### What I learned

- The document-processing and enrichment selectors are similar enough to share styling, but their selected identity shapes differ enough that a generic identity-selector abstraction would be premature.
- Replacing global `.btn` here also helps the later global CSS shrink pass because these buttons were one of the remaining active consumers of generic button styling.

### What was tricky to build

- The tricky part was preserving the exact selection keys and callback payloads. Document-processing selection includes `artifact_type`, `prompt_version`, `provider`, and `model`; chunk enrichment selection only includes `strategy_id` and `prompt_version`.
- Empty enrichment currently renders nothing by design. The story documents that behavior explicitly.

### What warrants a second pair of eyes

- Whether selected buttons with nested `CodeText` are visually clear enough in the retro monochrome style.
- Whether the empty enrichment behavior should eventually render a caption like preprocessing does.

### What should be done in the future

- Clean `ChunkTimelineBar` next: move global chunk-bar CSS into a module and add stories.
- Revisit whether artifact identity selectors should split into separate files after any future behavior divergence.

### Code review instructions

- Review `ArtifactIdentityBar.tsx` for unchanged selected identity comparisons and `onChange` payloads.
- Review the new stories for coverage of available and empty identity states.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```

## Step 26: Add Button Atom and Move ChunkTimelineBar CSS Local

The selected artifact identity buttons exposed a design-system gap: selected buttons were implemented ad hoc, and nesting `CodeText` inside a black selected button produced unreadable black-on-black text. I fixed the immediate contrast issue by introducing a reusable `Button` atom with explicit default, primary, selected, disabled, normal, and compact states, then migrated the artifact identity selectors to use it.

I also completed the next Corpus cleanup target by moving `ChunkTimelineBar` out of global `chunk-bar-*` CSS and into a local CSS Module. The only remaining inline styles in that component are dynamic geometry (`left` and `width`), which the guidelines explicitly allow.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue the next high-priority cleanup after artifact identity bars, namely `ChunkTimelineBar`.

**Inferred user intent:** The user wants the audit cleanup sequence to continue with stories and local CSS modules.

**User prompt (verbatim):** "selected style in say, Artifact IdentityBar is all black : /tmp/pi-clipboard-873bd5fb-4e01-4720-b568-fe7e0ab2f07d.png (this issue might appear in other pleaces, it feels like Button should be an Atom maybe?"

**Assistant interpretation:** The selected artifact identity button has a contrast/readability bug, and the broader pattern suggests buttons should be represented as a reusable atom rather than ad hoc/global button classes.

**Inferred user intent:** The user wants the design system to address button states systematically, not just patch one selected-state CSS bug.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Added a reusable Button atom:
  - `web/src/components/atoms/Button/Button.tsx`
  - `Button.module.css`
  - `Button.stories.tsx`
  - `Button/index.ts`
  - `atoms/index.ts`
- Exported atoms from `web/src/components/index.ts`.
- Migrated `ArtifactIdentityBar` selected identity options to `Button size="compact" selected`.
- Removed the nested `CodeText` inside selected identity buttons to avoid black-on-black text.
- Simplified `ArtifactIdentityBar.module.css` to only local bar/label/empty anatomy.
- Added `ChunkTimelineBar.module.css` and moved the timeline CSS out of `web/src/index.css`.
- Converted timeline chunks from generic divs to buttons with `aria-label` and `aria-pressed`.
- Added `ChunkTimelineBar.stories.tsx` with mixed coverage, all embedded, selected chunk, and empty states.
- Added `ChunkTimelineBar/index.ts`.
- Verified `rg -n "chunk-bar" src` returns no consumers after deleting the global block.

### Why

- The screenshot showed the selected artifact identity button text was unreadable.
- Button states are a repeated atom-level concern, not a local Corpus concern.
- `ChunkTimelineBar` was one of the last high-priority Corpus widgets using global CSS.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed.
- The Storybook build now includes `Design System/Atoms/Button` and `Component Library/Corpus/ChunkTimelineBar`.

### What didn't work

- The first `ChunkTimelineBar` story attempt omitted required `onSelect` in default args, causing Storybook/TypeScript errors. I added a noop default and reran validation successfully.
- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.

### What I learned

- The selected-state contrast issue was not just an ArtifactIdentityBar bug; it showed that button visual states should be owned by an atom.
- `ChunkTimelineBar` is a good example of the inline-style rule: dynamic geometry belongs inline, static geometry/anatomy belongs in a CSS Module.

### What was tricky to build

- The selected identity label originally used `CodeText`, whose CSS forces `color: var(--mac-text)`. That was incompatible with a selected black button. Rather than adding cross-component override hacks, the selected button now owns the text style directly.
- The timeline needed to remain absolutely positioned, so only `left` and `width` stayed inline.

### What warrants a second pair of eyes

- Whether the Button atom should immediately replace all remaining global `.btn` consumers or whether that should happen in smaller focused passes.
- Whether compact selected buttons need a stronger focus-visible style for keyboard users.

### What should be done in the future

- Migrate remaining `.btn`/`.btn-primary` consumers to the Button atom.
- Continue global CSS cleanup by removing button globals only after `rg` confirms no active consumers.
- Move on to Search/Corpus/Workflow page-boundary story splits after Corpus identity/timeline cleanup.

### Code review instructions

- Start with `web/src/components/atoms/Button/Button.tsx` and `Button.module.css`.
- Review `ArtifactIdentityBar.tsx` to confirm selected identity payload behavior is unchanged and selected text remains readable.
- Review `ChunkTimelineBar.tsx` and `.module.css` to confirm only dynamic geometry remains inline.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
rg -n "chunk-bar" src || true
```

## Step 27: Migrate Remaining Global Button Consumers

After adding the Button atom, I completed the first global CSS retirement pass by replacing the remaining active `.btn` and `.btn-primary` consumers in components and stories. The only remaining match is the legacy `retro/MacButton` class name, which uses a different `mac-btn` namespace and is part of the separate retro-component decision.

This turns the Button atom from a one-off fix for ArtifactIdentityBar into the actual system-owned button contract. It also allowed the old global `.btn` block to be removed from `web/src/index.css` safely.

### Prompt Context

**User prompt (verbatim):** "ok, continue"

**Assistant interpretation:** Continue applying the Button atom and global CSS cleanup after fixing the selected ArtifactIdentityBar contrast issue.

**Inferred user intent:** The user wants the design-system cleanup to proceed systematically, not leave old global button styles around after introducing the atom.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Migrated active `.btn` / `.btn-primary` consumers to `Button` in:
  - `SearchControlsPanel`,
  - `ResultInspectorPanel`,
  - `WorkflowSummaryPanel`,
  - `WorkflowOpInspectorPanel`,
  - `EmbeddingsView`,
  - `WorkflowsView`,
  - `DocumentBrowser`,
  - `DocumentInspector`,
  - `Inline.stories`,
  - `VisuallyHidden.stories`.
- Moved small static spacing previously inline on migrated buttons into CSS Modules where needed.
- Deleted the global `.btn` / `.btn-primary` CSS block from `web/src/index.css`.
- Verified no active `.btn` or `.btn-primary` consumers remain outside legacy `retro/MacButton`.

### Why

- Once `Button` exists, global button classes should no longer be the app-facing contract.
- Removing `.btn` from active components reduces `index.css` and makes button states Storybook-owned.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed.
- `rg` confirmed no active `.btn`/`.btn-primary` consumers remain, except unrelated legacy `mac-btn` naming.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.
- Storybook still emits known iframe chunk and plugin timing warnings; build succeeds.

### What I learned

- The Button atom is immediately useful across Search, Corpus, Workflow, Embeddings, and stories.
- The global CSS can shrink safely when migration is backed by `rg` checks.

### What was tricky to build

- The tricky part was distinguishing real `.btn` consumers from `retro/MacButton`'s `mac-btn` namespace. I left `retro/MacButton` untouched because the audit already treats the retro folder as a separate keep-or-delete decision.
- Several buttons had small inline font-size/margin tweaks. I avoided expanding the Button API for those and either used `size="compact"` or local CSS Module spacing.

### What warrants a second pair of eyes

- Whether `Button size="compact"` should become the default in dense panels.
- Whether `copy-btn` should become a `Button` variant or a separate `IconButton` atom in a later pass.

### What should be done in the future

- Consider an `IconButton` atom for close/copy buttons currently using `copy-btn`.
- Continue retiring global style families in `index.css` with the same migrate-then-delete pattern.

### Code review instructions

- Review the Button atom consumers to confirm behavior is unchanged.
- Confirm `web/src/index.css` no longer contains the global button block.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`
  - `cd web && rg -n 'className="btn|btn-primary|\\.btn|`btn' src/components src/index.css src/App.tsx || true`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
rg -n 'className="btn|btn-primary|\.btn|`btn' src/components src/index.css src/App.tsx || true
```

## Step 28: Add IconButton Atom and Remove copy-btn Globals

I added the second atom implied by the button cleanup: `IconButton` for compact close, copy, and back actions. The previous `copy-btn` global class was doing atom work from `index.css`; it is now replaced by a co-located atom with Storybook coverage.

This continues the migrate-then-delete pattern from the Button atom: first move all consumers to the atom, verify no active global consumers remain with `rg`, and only then delete the global CSS block.

### Prompt Context

**User prompt (verbatim):** "add the atom, continue"

**Assistant interpretation:** Add the compact icon/copy button atom discussed in the previous turn and keep migrating global CSS consumers.

**Inferred user intent:** The user wants small action buttons to be design-system owned, not implemented through global `copy-btn` classes.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Added `IconButton` atom:
  - `web/src/components/atoms/IconButton/IconButton.tsx`
  - `IconButton.module.css`
  - `IconButton.stories.tsx`
  - `IconButton/index.ts`
- Exported `IconButton` from `web/src/components/atoms/index.ts`.
- Migrated `copy-btn` consumers in:
  - `WorkflowsView`,
  - `WorkflowOpInspectorPanel`,
  - `WorkflowSummaryPanel`,
  - `CorpusExplorerView`,
  - `DocumentInspector`,
  - `Panel.stories`.
- Replaced the `Panel` story's raw `meta-grid` demo with `MetadataGrid` while touching that story.
- Deleted the global `.copy-btn` block from `web/src/index.css`.
- Verified `rg -n "copy-btn|\\.copy-btn" src src/index.css` returns no consumers.

### Why

- Close/copy/back controls are repeated atom-level UI states.
- Keeping them in global CSS prevents `index.css` from shrinking and hides the contract from Storybook.
- The Button atom should not be overloaded for terse icon-like actions with no border.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed and included the new `IconButton` stories.
- `copy-btn` has no remaining active consumers.

### What didn't work

- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.
- Storybook still emits the known large iframe chunk and plugin timing warnings; build succeeds.

### What I learned

- `IconButton` and `Button` give us a clean atom split: bordered action buttons vs compact chromeless utility actions.
- Several stories still contained raw global examples; touching atoms is a good moment to clean those stories too.

### What was tricky to build

- `IconButton` needs an accessible `label` prop because many usages are glyph-only (`✕`, `⧉`). Existing text usages like `← Back` still benefit from the same atom because they use the same chromeless action style.
- The old `copy-btn` was also used for chunk ID text buttons, so the atom supports arbitrary children rather than requiring only an icon.

### What warrants a second pair of eyes

- Whether `IconButton` should require `label` when children are non-text glyphs. TypeScript cannot easily enforce that without making common text cases clumsy.
- Whether `IconButton` should have a selected or danger variant later; not needed yet.

### What should be done in the future

- Continue deleting global CSS families only after migrating all consumers.
- Consider the next cleanup family: raw `panel`/`panel-header` fallback in page containers, or `text-*` globals.

### Code review instructions

- Review `web/src/components/atoms/IconButton/IconButton.tsx` and story coverage.
- Review migrated close/copy/back usages for unchanged behavior.
- Confirm `web/src/index.css` no longer contains `.copy-btn`.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`
  - `cd web && rg -n "copy-btn|\\.copy-btn" src src/index.css || true`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
rg -n "copy-btn|\.copy-btn" src src/index.css || true
```

## Step 29: Add EvaluationPage Boundary and Clean Layout Stories

I continued the guideline cleanup by replacing the raw Evaluation placeholder panel with a small storyable page boundary. `EvaluationView` is now only the runtime wrapper, while `EvaluationPage` owns the presentational placeholder and has Storybook coverage.

I also cleaned two layout primitive stories that still demonstrated layout recipes with raw global `panel` markup. Those stories now compose the actual `Panel` primitive, so the Storybook examples reinforce the guidelines instead of preserving legacy global CSS patterns.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue the post-atom cleanup by reducing remaining global panel/text usage and adding page story boundaries where low-risk.

**Inferred user intent:** The user wants the design-system guideline audit applied broadly, including placeholders and stories, not just major feature widgets.

**Commit (code):** N/A at time of writing — changes pending commit.

### What I did

- Added `web/src/components/pages/EvaluationPage/` with:
  - `EvaluationPage.tsx`,
  - `EvaluationPage.module.css`,
  - `EvaluationPage.stories.tsx`,
  - `index.ts`.
- Updated `web/src/components/evaluation/EvaluationView.tsx` to delegate to `EvaluationPage`.
- Exported `EvaluationPage` from `web/src/components/pages/index.ts`.
- Replaced raw `.panel` snippets in:
  - `Stack.stories.tsx`,
  - `DashboardGrid.stories.tsx`.

### Why

- `EvaluationView` was still a raw global `panel` placeholder with no story.
- Layout stories should model the real design-system primitives, not legacy global class recipes.
- This is a low-risk way to continue shrinking active global panel usage while improving page-level Storybook coverage.

### What worked

- `pnpm typecheck` passed.
- `pnpm build` passed.
- `pnpm build-storybook` passed and included `Pages/EvaluationPage`.

### What didn't work

- I initially used `Text tone="secondary"`, but RAG `Text` supports `muted` rather than `secondary`; TypeScript caught this and I corrected it.
- `pnpm build` rewrote `internal/web/dist/index.html`; I reverted the generated embed artifact before committing.

### What I learned

- Page boundaries can be added incrementally even for placeholders, as long as runtime containers remain simple wrappers.
- Story files are still a source of legacy global CSS debt and should be included in cleanup scans.

### What was tricky to build

- The only notable mismatch was terminology drift between TTC-style tone names and RAG foundation tone names. RAG uses `muted`, not `secondary`, so future docs/examples should use the RAG primitive API exactly.

### What warrants a second pair of eyes

- Whether `EvaluationPage` should remain in `components/pages` while it is only a placeholder; it probably should because it establishes the future page contract.
- Whether the evaluation capability text should be linked to future DMETA documentation when that work starts.

### What should be done in the future

- Add storyable page boundaries for Corpus, Workflows, and Embeddings.
- Continue replacing raw panel/text classes in active page containers.

### Code review instructions

- Review `web/src/components/pages/EvaluationPage/EvaluationPage.tsx` and story first.
- Confirm `EvaluationView` remains a no-logic runtime wrapper.
- Review layout story changes to ensure they still demonstrate `Stack` and `DashboardGrid` recipes clearly.
- Validate with:
  - `cd web && pnpm typecheck`
  - `cd web && pnpm build`
  - `cd web && pnpm build-storybook`

### Technical details

Validation commands that passed:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
```
