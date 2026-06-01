# Tasks

## DONE

- [x] Create RAG-local docmgr ticket with `docmgr --root ttmp`.
- [x] Inspect the React web app source tree.
- [x] Inspect the RTK Query API layer and Redux store.
- [x] Inspect the global CSS/design-system state.
- [x] Inspect Go server, embedded SPA, API route registry, workflow/artifact handlers, and database concepts needed to understand the web app.
- [x] Review prior RAGEVAL web/retrieval/workflow docs for context.
- [x] Run static scans for file size, class usage, inline style pressure, custom events, and missing Storybook/CSS-module/DMETA artifacts.
- [x] Attempt frontend build and record dependency failure.
- [x] Write the intern-facing design/architecture/review report.
- [x] Write the investigation diary.
- [x] Relate key files to the ticket docs.
- [x] Run `docmgr --root ttmp doctor`.
- [x] Upload the report bundle to reMarkable.
- [x] Phase 0: install frontend dependencies and restore a passing local web build baseline.
- [x] Phase 0: add Storybook package scripts/config so component review can start before large refactors.
- [x] Phase 1: extract current retro visual tokens into `web/src/styles/tokens.css`.
- [x] Phase 1: add foundation primitives for text, code identifiers, statuses, dividers, and visually hidden labels.
- [x] Phase 1: add `Caption` foundation primitive for compact metadata labels and statuses.
- [x] Phase 1: add layout primitives for panels, stacks, inline toolbars, and dashboard grid recipes.
- [x] Phase 1: add `AppShell`, `ScrollRegion`, `TabList`, and `FormRow` layout primitives.
- [x] Phase 2: start Storybook immediately with co-located foundation primitive stories.
- [x] Phase 2: start Storybook immediately with co-located layout primitive stories.
- [x] Phase 2: validate the new Storybook setup with `pnpm build-storybook`.

## TODO

- [x] Phase 3: extract `SearchControlsPanel` with interactive and loading Storybook stories.
- [x] Phase 3: extract `RetrievalResultsPanel` with empty, loading, hybrid results, and error Storybook stories.
- [x] Phase 3: extract `CoveragePanel` with complete and sparse coverage Storybook stories.
- [x] Phase 3: extract `QueryPresetList` with TTC query Storybook story.
- [x] Phase 3: refactor `SearchView.tsx` to compose extracted controls/results/coverage/preset components while preserving behavior.
- [x] Phase 3: extract `ResultInspectorPanel` and add stories for detail/BM25/document-loading inspector states.
- [x] Phase 3: migrate Search Workbench page shell from inline flex layout to `DashboardGrid`/`Stack` plus `SearchView.module.css`.
- [x] Phase 4: add Storybook review surfaces for the Search Workbench organisms/molecules extracted so far.
- [x] Phase 4: add Storybook review surfaces for `Caption`, `AppShell`, `ScrollRegion`, `TabList`, `FormRow`, `MetadataGrid`, and `DataTable`.
- [ ] Phase 4 optional follow-up: expand `ResultInspectorPanel` stories with explicit chunk-tab, document-tab, and no-document states once tab state can be controlled externally or via story interactions.
- [x] Phase 5: add initial documentation-only `dmeta-ir` core/interaction/Web MDS YAML for the Search Workbench vertical slice.
- [x] Phase 3: add reusable `MetadataGrid` and `DataTable` molecules as extraction targets for inspector/results/corpus/workflow views.
- [ ] Phase 3 follow-up: adopt `MetadataGrid` and `DataTable` inside Search Workbench, Corpus Explorer, and Workflow views.
- [ ] Deferred beyond requested scope (Phase 6): generated React scaffold output and promotion metadata.
- [ ] Deferred beyond requested scope (Phase 7): executable validators for the documentation-only IR catalog.
