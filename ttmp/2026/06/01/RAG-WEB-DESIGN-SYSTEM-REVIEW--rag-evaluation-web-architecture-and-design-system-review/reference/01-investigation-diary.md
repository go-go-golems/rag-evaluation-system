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
      Note: App shell and navigation event evidence
    - Path: web/src/components/foundation/index.ts
      Note: Phase 1 foundation primitive barrel
    - Path: web/src/components/layout/index.ts
      Note: Phase 1 layout primitive barrel
    - Path: web/src/components/molecules/CoveragePanel/CoveragePanel.tsx
      Note: Phase 2 extraction diary evidence
    - Path: web/src/components/molecules/QueryPresetList/QueryPresetList.tsx
      Note: Phase 2 extraction diary evidence
    - Path: web/src/components/organisms/ResultInspectorPanel/ResultInspectorPanel.tsx
      Note: Phase 2 inspector extraction diary evidence
    - Path: web/src/components/organisms/RetrievalResultsPanel/RetrievalResultsPanel.tsx
      Note: Phase 2 extraction diary evidence
    - Path: web/src/components/organisms/SearchControlsPanel/SearchControlsPanel.tsx
      Note: Phase 2 extraction diary evidence
    - Path: web/src/components/search/SearchView.module.css
      Note: Phase 2 layout primitive migration diary evidence
    - Path: web/src/components/search/SearchView.tsx
      Note: |-
        First recommended vertical-slice evidence
        Phase 2 Search Workbench composition update
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
