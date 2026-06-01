# Changelog

## 2026-06-01

- Initial workspace created


## 2026-06-01

Created intern-facing RAG web architecture/design-system/DMETA review and investigation diary

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/design-doc/01-rag-evaluation-web-architecture-and-design-system-review.md — Primary review deliverable
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Investigation diary


## 2026-06-01

Validated RAG web review ticket and uploaded report bundle to reMarkable at /ai/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/design-doc/01-rag-evaluation-web-architecture-and-design-system-review.md — Uploaded primary report
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Updated delivery diary


## 2026-06-01

Implemented Phase 0/1/early-2 web design-system scaffolding: dependencies, Storybook config, tokens, foundation primitives, layout primitives, and co-located primitive stories

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded Phase 0/1/2 implementation diary
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/package.json — Added typecheck and Storybook scripts/dependencies
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/foundation/index.ts — Foundation primitive exports
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/layout/index.ts — Layout primitive exports
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/styles/tokens.css — Extracted retro tokens and added RAG role aliases


## 2026-06-01

Validated web Phase 0/1/2 with pnpm typecheck, pnpm build, and pnpm build-storybook; noted go test is blocked in this checkout by missing ../scraper replace target

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded validation results and Go test blocker
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/package.json — Web validation scripts


## 2026-06-01

Continued Phase 2 by extracting SearchControlsPanel, RetrievalResultsPanel, CoveragePanel, and QueryPresetList with co-located Storybook stories

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded Phase 2 Search Workbench extraction
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/molecules/CoveragePanel/CoveragePanel.tsx — Extracted embedding coverage molecule
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/molecules/QueryPresetList/QueryPresetList.tsx — Extracted query preset molecule
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/RetrievalResultsPanel/RetrievalResultsPanel.tsx — Extracted retrieval results organism
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/SearchControlsPanel/SearchControlsPanel.tsx — Extracted search controls organism
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/search/SearchView.tsx — Refactored to compose extracted Search Workbench components


## 2026-06-01

Completed the next Phase 2 increment by extracting ResultInspectorPanel and moving the Search Workbench shell onto DashboardGrid/Stack

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded ResultInspector extraction and layout primitive adoption
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/ResultInspectorPanel/ResultInspectorPanel.stories.tsx — Storybook examples for inspector states
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/ResultInspectorPanel/ResultInspectorPanel.tsx — Extracted selected-result inspector organism
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/search/SearchView.module.css — Page-local sizing and scroll ownership for Search Workbench shell
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/search/SearchView.tsx — Search page now owns data orchestration and composes ResultInspectorPanel with layout primitives


## 2026-06-01

Started Phase 3 by adding a RAG-local dmeta-ir catalog for the Search Workbench vertical slice

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/dmeta-ir/README.md — Documents RAG DMETA IR layer boundaries
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/dmeta-ir/core-model/archetypes/10-search-workbench.yaml — Core Search Workbench semantics
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/dmeta-ir/instantiations/rag-evaluation-dashboard.yaml — First Search Workbench instantiation
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/dmeta-ir/interactions/actions/10-search-workbench.yaml — Search Workbench actions
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/dmeta-ir/interactions/representations/10-search-workbench.yaml — Search Workbench visible obligations
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/dmeta-ir/meta-design-systems/web/lowering-rules.yaml — Interaction-to-Web lowering rules
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/dmeta-ir/meta-design-systems/web/widgets/search-workbench.yaml — Web MDS Search Workbench templates


## 2026-06-01

Added a local dmeta-ir validator and Makefile target for YAML parsing, manifest targets, required human fields, duplicate IDs, and key references

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/Makefile — Adds make dmeta-validate target
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/scripts/validate-dmeta-ir.py — Local validator for RAG DMETA IR seed catalog
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded validator implementation and initial false positives


## 2026-06-01

Adjusted scope to documentation-only DMETA IR through Phase 5; removed premature Phase 7 validator tooling and deferred scaffolds/validators

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/dmeta-ir/README.md — Documents that the RAG DMETA IR catalog is documentation-only for current scope
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/dmeta-ir/instantiations/rag-evaluation-dashboard.yaml — Clarifies no validation or generated React scaffolds are required now
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Records scope correction and validator removal
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/tasks.md — Aligns task list with phases 3-5 and defers phases 6-7


## 2026-06-01

Expanded the reusable RAG design system with Caption, AppShell, ScrollRegion, TabList, FormRow, MetadataGrid, and DataTable plus stories

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded expanded design-system primitive extraction
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/foundation/Caption/Caption.tsx — Foundation caption primitive for compact metadata labels
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/layout/AppShell/AppShell.tsx — Application shell layout primitive
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/layout/FormRow/FormRow.tsx — Reusable form row primitive
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/layout/ScrollRegion/ScrollRegion.tsx — Reusable scroll ownership primitive
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/layout/TabList/TabList.tsx — Controlled tab-list layout primitive
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/molecules/DataTable/DataTable.tsx — Reusable dense data table molecule
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/molecules/MetadataGrid/MetadataGrid.tsx — Reusable key/value metadata molecule


## 2026-06-01

Adopted DataTable, MetadataGrid, TabList, and ScrollRegion inside Search Workbench result and inspector panels

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded Search Workbench reusable component adoption
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/ResultInspectorPanel/ResultInspectorPanel.module.css — Inspector-specific anatomy after reusable primitive adoption
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/ResultInspectorPanel/ResultInspectorPanel.stories.tsx — Expanded explicit inspector tab stories
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/ResultInspectorPanel/ResultInspectorPanel.tsx — Inspector now uses MetadataGrid/TabList/ScrollRegion and supports defaultTab stories
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/RetrievalResultsPanel/RetrievalResultsPanel.tsx — Retrieval results now use reusable DataTable/Panel/ScrollRegion


## 2026-06-01

Moved Search into a SearchWorkbenchPage boundary and adopted AppShell/FormRow/Panel/Stack/ScrollRegion/Caption in app shell and controls

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded page-boundary and controls primitive adoption
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/App.tsx — Root application now uses AppShell
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/SearchControlsPanel/SearchControlsPanel.module.css — Controls-specific anatomy after primitive adoption
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/SearchControlsPanel/SearchControlsPanel.tsx — Controls panel now uses reusable design-system primitives
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/pages/SearchWorkbenchPage/SearchWorkbenchPage.tsx — Explicit Search Workbench page boundary
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/search/SearchView.tsx — Compatibility re-export for moved SearchWorkbenchPage


## 2026-06-01

Adopted reusable design-system primitives and molecules in Corpus source, document, and inspector surfaces

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded Corpus design-system adoption
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/DocumentBrowser.tsx — Document browser now uses DataTable/Panel/ScrollRegion
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/DocumentInspector.module.css — Document-inspector-specific anatomy after primitive adoption
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/DocumentInspector.tsx — Document inspector now uses TabList/MetadataGrid/DataTable
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/SourcePanel.tsx — Source list now uses Panel/ScrollRegion/Caption


## 2026-06-01

Added Corpus Storybook coverage and aligned RAG foundation extraction with the TTC foundation guide

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded Corpus Storybook coverage and TTC foundation guide takeaways
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/DocumentBrowser.stories.tsx — Storybook coverage for Corpus document browser states
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/DocumentInspector.stories.tsx — Storybook coverage for Corpus document inspector states
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/SourcePanel.stories.tsx — Storybook coverage for Corpus source panel states
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--ttc-design-system/ttmp/2026/06/01/TTC-FOUNDATION-SYSTEM--ttc-react-foundation-primitives-and-token-documentation/design-doc/01-react-foundation-system-implementation-guide.md — Foundation extraction and CSS cleanup guidelines applied to RAG work


## 2026-06-01

Reclassified Corpus source selection as an organism, added workflow list/queue organisms with Storybook stories, and added Redux provider to DocumentInspector stories

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/DocumentInspector.stories.tsx — Redux Provider wrapper for RTK Query-backed document inspector stories
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/SourcePanel.tsx — Compatibility re-export for moved CorpusSourcePanel
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/CorpusSourcePanel/CorpusSourcePanel.tsx — Organism-level corpus source selector
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/QueueHealthPanel/QueueHealthPanel.tsx — Workflow queue health organism
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/WorkflowListPanel/WorkflowListPanel.tsx — Workflow list organism
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/workflows/WorkflowsView.tsx — Delegates workflow list and queue presentation to organisms


## 2026-06-01

Moved SourcePanel back into Corpus and extracted workflow summary, op graph, and op group organisms with Storybook stories

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded SourcePanel move-back and workflow detail extraction
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/SourcePanel.tsx — Feature-local Corpus source panel using shared primitives
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/WorkflowOpGraphPanel/WorkflowOpGraphPanel.tsx — Workflow op graph organism
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/WorkflowOpGroupsPanel/WorkflowOpGroupsPanel.tsx — Workflow op groups table organism
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/WorkflowSummaryPanel/WorkflowSummaryPanel.tsx — Workflow detail summary/progress organism
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/workflows/WorkflowsView.tsx — Workflow detail now composes extracted organisms


## 2026-06-01

Extracted WorkflowOpInspectorPanel with Storybook stories and refactored WorkflowsView sampled-op inspection through it

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded workflow op inspector extraction
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/WorkflowOpInspectorPanel/WorkflowOpInspectorPanel.stories.tsx — Storybook coverage for succeeded and failed sampled op states
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/WorkflowOpInspectorPanel/WorkflowOpInspectorPanel.tsx — Sampled workflow op inspector organism
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/workflows/WorkflowsView.tsx — Delegates sampled op inspection to WorkflowOpInspectorPanel


## 2026-06-01

Extracted WorkflowOpResultPanel with Storybook states and delegated op result rendering from WorkflowsView

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded workflow op result extraction
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/WorkflowOpResultPanel/WorkflowOpResultPanel.stories.tsx — Storybook coverage for workflow op result states
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/WorkflowOpResultPanel/WorkflowOpResultPanel.tsx — Workflow op result renderer using MetadataGrid and DataTable
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/workflows/WorkflowsView.tsx — OpResultSection now fetches and delegates rendering


## 2026-06-01

Refactored SubmitIntakeModal through Panel, Stack, FormRow, and Caption while preserving workflow submission behavior

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded intake modal primitive adoption
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/workflows/WorkflowsView.tsx — Workflow intake modal now uses shared form/layout primitives


## 2026-06-01

Cleaned up CoveragePanel and QueryPresetList to use shared foundation/layout/molecule primitives instead of global panel/table/text classes

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded Search support molecule CSS ownership cleanup
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/molecules/CoveragePanel/CoveragePanel.tsx — Coverage support molecule now uses Panel/MetadataGrid/DataTable/StatusText
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/molecules/QueryPresetList/QueryPresetList.tsx — Query presets now use Panel/Stack/Caption with local button anatomy


## 2026-06-01

Refactored EmbeddingsView through shared Panel/FormRow/MetadataGrid/DataTable primitives while preserving API-owned behavior

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded embeddings view primitive adoption
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/embeddings/EmbeddingsView.module.css — View-specific embeddings anatomy after primitive adoption
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/embeddings/EmbeddingsView.tsx — Embeddings view now composes shared design-system primitives


## 2026-06-01

Added storyable PipelineOverview/PipelinePage boundaries and organized Corpus widgets into per-widget directories

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded Pipeline stories and Corpus organization
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/CorpusExplorerView.tsx — Updated Corpus imports after per-widget folder organization
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/pages/PipelinePage/PipelinePage.tsx — Storyable pipeline page boundary
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/pipeline/PipelineOverview.tsx — Storyable pipeline presentation boundary


## 2026-06-01

Added RAG React design-system guidelines modeled after the TTC foundation guide, covering CSS ownership, file layout, Storybook, pages, and cleanup phases

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/../2026-05-27--ttc-design-system/ttmp/2026/06/01/TTC-FOUNDATION-SYSTEM--ttc-react-foundation-primitives-and-token-documentation/design-doc/01-react-foundation-system-implementation-guide.md — Reference guide used as source model
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/design-doc/02-rag-react-design-system-guidelines.md — RAG-local design-system guideline document


## 2026-06-01

Added detailed follow-up tasks, completed the guideline audit, and added RAG foundation overview Storybook docs

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/analysis/01-rag-design-system-guideline-audit.md — Component-by-component guideline audit and extraction plan
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/foundation/Foundation.stories.tsx — Foundation overview stories


## 2026-06-01

Refactored Corpus IdentityBar through shared primitives and added Storybook states

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded IdentityBar cleanup
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/IdentityBar/IdentityBar.stories.tsx — IdentityBar Storybook coverage
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/IdentityBar/IdentityBar.tsx — IdentityBar primitive adoption


## 2026-06-01

Refactored Corpus artifact identity selectors through shared primitives and added Storybook states

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded artifact identity cleanup
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/ArtifactIdentityBar/ArtifactIdentityBar.stories.tsx — Artifact identity Storybook coverage
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/ArtifactIdentityBar/ArtifactIdentityBar.tsx — Artifact identity primitive adoption


## 2026-06-01

Added Button atom, fixed selected artifact identity contrast, and moved ChunkTimelineBar styling into a CSS Module

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded Button atom and ChunkTimelineBar cleanup
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/atoms/Button/Button.tsx — Reusable Button atom for selected/primary/default states
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/ArtifactIdentityBar/ArtifactIdentityBar.tsx — Artifact identity selectors now use Button atom
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/corpus/ChunkTimelineBar/ChunkTimelineBar.tsx — Chunk timeline now uses local CSS module and accessible buttons
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/index.css — Removed unused global chunk-bar styles


## 2026-06-01

Migrated remaining active global button consumers to the Button atom and removed the legacy global .btn CSS block

### Related Files

- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/reference/01-investigation-diary.md — Recorded global button migration
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/atoms/Button/Button.tsx — Button atom now owns active button styling
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/organisms/SearchControlsPanel/SearchControlsPanel.tsx — Search controls use Button atom
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/components/workflows/WorkflowsView.tsx — Workflow actions use Button atom
- /home/manuel/workspaces/2026-05-27/ttc-design-system/2026-05-27--rag-evaluation-system/web/src/index.css — Removed legacy global button styles

