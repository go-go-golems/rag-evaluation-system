---
Title: RAG Design System Guideline Audit
Ticket: RAG-WEB-DESIGN-SYSTEM-REVIEW
Status: active
Topics:
    - rag
    - react
    - design-system
    - css
    - dmeta
DocType: analysis
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ttmp/2026/06/01/RAG-WEB-DESIGN-SYSTEM-REVIEW--rag-evaluation-web-architecture-and-design-system-review/design-doc/02-rag-react-design-system-guidelines.md
      Note: Guidelines applied by this audit
    - Path: web/src/components/atoms/Button/Button.stories.tsx
      Note: Button atom Storybook states
    - Path: web/src/components/atoms/Button/Button.tsx
      Note: Button atom added after selected-state contrast issue
    - Path: web/src/components/corpus/ArtifactIdentityBar/ArtifactIdentityBar.module.css
      Note: Artifact identity local selector anatomy
    - Path: web/src/components/corpus/ArtifactIdentityBar/ArtifactIdentityBar.stories.tsx
      Note: Artifact identity required Storybook coverage
    - Path: web/src/components/corpus/ArtifactIdentityBar/ArtifactIdentityBar.tsx
      Note: |-
        High-priority artifact identity cleanup target
        Artifact identity cleanup completed after audit
    - Path: web/src/components/corpus/ChunkTimelineBar/ChunkTimelineBar.module.css
      Note: Timeline local CSS replacing global chunk-bar styles
    - Path: web/src/components/corpus/ChunkTimelineBar/ChunkTimelineBar.stories.tsx
      Note: Timeline Storybook coverage
    - Path: web/src/components/corpus/ChunkTimelineBar/ChunkTimelineBar.tsx
      Note: |-
        High-priority timeline CSS-module cleanup target
        ChunkTimelineBar cleanup completed after audit
    - Path: web/src/components/corpus/IdentityBar/IdentityBar.module.css
      Note: IdentityBar local anatomy CSS
    - Path: web/src/components/corpus/IdentityBar/IdentityBar.stories.tsx
      Note: IdentityBar required Storybook coverage
    - Path: web/src/components/corpus/IdentityBar/IdentityBar.tsx
      Note: |-
        High-priority identity selector cleanup target
        IdentityBar cleanup completed after audit
    - Path: web/src/components/embeddings/EmbeddingsView.tsx
      Note: Button atom migration evidence
    - Path: web/src/components/foundation/Foundation.stories.module.css
      Note: Story-only foundation overview styling
    - Path: web/src/components/foundation/Foundation.stories.tsx
      Note: Foundation overview Storybook docs added after audit
    - Path: web/src/components/organisms/SearchControlsPanel/SearchControlsPanel.tsx
      Note: Button atom migration evidence
    - Path: web/src/components/workflows/WorkflowsView.tsx
      Note: Button atom migration evidence
    - Path: web/src/index.css
      Note: Legacy button global removed after migration
ExternalSources: []
Summary: Component-by-component audit of the current RAG React dashboard against the RAG React Design System Guidelines, including remaining CSS debt, missing Storybook stories, page-boundary gaps, and extraction recommendations.
LastUpdated: 2026-06-01T19:30:00-04:00
WhatFor: Use this audit to plan the next cleanup passes after foundation Storybook docs and to decide what should be extracted from each component.
WhenToUse: Read before touching RAG React components, deleting global CSS, adding stories, or extracting new primitives/molecules/organisms/pages.
---






# RAG Design System Guideline Audit

## Executive Summary

The RAG frontend is now much closer to a real design system than when the ticket started, but it is not done. The strongest areas are foundation primitives, layout primitives, reusable table/metadata molecules, Search organisms, Workflow result/detail organisms, Pipeline presentation/page boundaries, and the newly organized Corpus widget folders. The weakest areas are still the identity selector widgets in Corpus, page-level Storybook boundaries, remaining global CSS consumers, and the absence of a TTC-like foundation overview story surface.

This audit applies `design-doc/02-rag-react-design-system-guidelines.md` to the current component tree. It answers three questions:

1. What already follows the guidelines?
2. What still violates them?
3. After foundation docs are in place, what should be extracted or cleaned up component-by-component?

The highest-priority cleanup sequence is:

```text
Foundation overview stories
  -> Corpus IdentityBar / ArtifactIdentityBar / ChunkTimelineBar stories and primitive cleanup
  -> SearchWorkbenchPage presentational/story split
  -> Workflows page/panel split for remaining inline/modal state
  -> Embeddings page story boundary
  -> Evaluation placeholder cleanup
  -> global CSS deletion after rg confirms no consumers
```

## Audit Method

I scanned the React tree for:

- missing co-located stories;
- direct legacy global class usage;
- inline style pressure;
- public components that are not folder-based;
- API-heavy containers that need storyable presentational boundaries;
- candidates for extraction into foundation/layout/molecule/organism/page layers.

Representative scan command:

```bash
cd 2026-05-27--rag-evaluation-system/web
python - <<'PY'
from pathlib import Path
for p in sorted(Path('src/components').rglob('*.tsx')):
    if p.name.endswith('.stories.tsx'):
        continue
    txt = p.read_text()
    story = p.with_name(p.stem + '.stories.tsx').exists()
    legacy = []
    for token in ['panel','panel-header','panel-body','data-table','meta-grid','tab-bar','tab-item','text-mono','text-dim','text-small','status-','form-row','stat-grid','chunk-bar','coverage-','op-graph','fieldset']:
        if token in txt:
            legacy.append(token)
    print(f'{p}|story={story}|inline={txt.count("style={{")}|legacy={",".join(legacy) if legacy else "-"}')
PY
```

## Current Layer Health

### Strong / mostly compliant

These areas mostly follow the guidelines:

- `web/src/components/foundation/*`
- `web/src/components/layout/*`
- `web/src/components/molecules/DataTable/*`
- `web/src/components/molecules/MetadataGrid/*`
- `web/src/components/molecules/CoveragePanel/*`
- `web/src/components/molecules/QueryPresetList/*`
- `web/src/components/organisms/SearchControlsPanel/*`
- `web/src/components/organisms/ResultInspectorPanel/*`
- `web/src/components/organisms/WorkflowOpGraphPanel/*`
- `web/src/components/organisms/WorkflowOpInspectorPanel/*`
- `web/src/components/organisms/WorkflowOpResultPanel/*`
- `web/src/components/pages/PipelinePage/*`
- `web/src/components/pipeline/PipelineOverview.*`

Why these are strong:

- They are folder-based where public/reusable.
- They generally have stories.
- They use primitives instead of raw global panels/tables/status text.
- Their CSS Modules are mostly local anatomy.

### Partially compliant

These are usable but still need cleanup:

- `SearchWorkbenchPage`: page exists, but no page story and still contains a small raw inspector fallback panel.
- `CorpusExplorerView`: still a runtime container with raw inspector fallback panel and inline layout.
- `DocumentInspector`: has a story, but still calls RTK Query internally for artifacts; eventually split artifact data fetching from presentation.
- `EmbeddingsView`: uses primitives well, but is still a container/view without a storyable page boundary.
- `WorkflowsView`: many organisms extracted, but the container still has inline styles, modal internals, and no page story boundary.
- `RetrievalResultsPanel`: has stories but still uses `text-dim` for empty/loading states.
- `WorkflowOpGroupsPanel`: has stories but still uses `text-dim text-mono` for empty state.
- `WorkflowSummaryPanel`: has stories but progress-bar styling remains tied to old global CSS unless already moved elsewhere; verify consumers before deleting globals.

### Non-compliant / next cleanup targets

These need near-term attention:

- `IdentityBar`: no story, raw global panel/text classes, many inline styles.
- `ArtifactIdentityBar`: no story, raw global text classes, inline styles, two separate identity selectors in one file.
- `ChunkTimelineBar`: no story, global `chunk-bar-*` classes.
- `EvaluationView`: placeholder with raw global panel/text classes and no story/page boundary.
- `retro/*`: stale or disconnected Mac components with no stories; either document as legacy or remove/replace.
- layout stories: some demos still use raw `.panel` snippets rather than `Panel`.

## Component-by-Component Audit and Extraction Plan

### Foundation package

Files:

- `foundation/Text`
- `foundation/CodeText`
- `foundation/Caption`
- `foundation/StatusText`
- `foundation/Divider`
- `foundation/VisuallyHidden`

Status: **mostly compliant**.

Recommended extraction/cleanup:

- Add `Foundation.stories.tsx` overview docs for colors, typography, status tones, spacing, borders/radii, and accessibility.
- Do not add `Heading`/`Eyebrow` yet unless dashboard headings start repeating beyond `Panel` titles and `Text`/`Caption`.
- Do not add generic `Box`.

### Layout package

Files:

- `Panel`
- `Stack`
- `Inline`
- `DashboardGrid`
- `AppShell`
- `ScrollRegion`
- `TabList`
- `FormRow`

Status: **mostly compliant**.

Recommended extraction/cleanup:

- Replace raw `.panel` snippets in layout stories with actual `Panel` components.
- Consider adding more `DashboardGrid` recipes only after repeated page layouts appear.
- Do not add a generic grid/box primitive for one-off layouts.

### Molecules

#### `DataTable`

Status: **compliant**.

Next extraction:

- Add optional `density` only if multiple consumers need compact vs normal row height.
- Add built-in loading/empty rendering only if it reduces repeated code in at least two components.

#### `MetadataGrid`

Status: **compliant**.

Next extraction:

- Use it inside `IdentityBar` for corpus totals.
- Keep copy behavior local to this molecule.

#### `CoveragePanel`

Status: **mostly compliant**.

Next extraction:

- If chunk/source coverage strips repeat, extract `CoverageStrip`; do not do this until `ChunkTimelineBar`/coverage usages prove overlap.

#### `QueryPresetList`

Status: **compliant**.

Next extraction:

- Keep as molecule unless presets become source-aware or API-backed.

### Search

#### `SearchWorkbenchPage`

Status: **partially compliant**.

Issues:

- No Storybook story for the page boundary.
- Still has raw `panel-header` fallback around inspector area.
- Runtime page owns RTK Query and mutations, which makes storying harder.

Recommended extraction:

1. Split a DTO-shaped `SearchWorkbenchPresentation` or `SearchWorkbenchPageFrame` from the API-heavy page.
2. Story page states:
   - empty query/no results;
   - loading search;
   - BM25 results;
   - vector sparse coverage;
   - selected result with document detail;
   - search error.
3. Replace any remaining raw panel fallback with `Panel` and `Caption`.

#### `SearchView`

Status: **compatibility wrapper**.

Recommendation:

- Keep as re-export only until callers are updated.
- Do not add logic here.

#### `SearchControlsPanel`, `RetrievalResultsPanel`, `ResultInspectorPanel`

Status: **mostly compliant**.

Recommended cleanup:

- `RetrievalResultsPanel`: replace any remaining `className="text-dim"` empty/loading copy with `Caption` or `Text`.
- Review inline styles; move static padding/centering to module CSS if repeated.

### Corpus

#### `CorpusExplorerView`

Status: **container, partially compliant**.

Issues:

- RTK Query-heavy container.
- Still has a raw fallback document inspector panel.
- No page story boundary.

Recommended extraction:

1. Keep `CorpusExplorerView` as runtime container.
2. Extract `CorpusExplorerPage` with DTO-shaped props:
   - identity values;
   - source summaries;
   - selected source;
   - document rows;
   - selected document detail;
   - processing identity/coverage summaries;
   - callbacks.
3. Add `Pages/CorpusExplorerPage` stories for loading, empty, selected document, and missing document states.
4. Replace raw fallback panels with `Panel`/`Caption`.

#### `SourcePanel`

Status: **mostly compliant**.

Recommended cleanup:

- Add `index.ts` inside the new folder if missing.
- Keep feature-local unless reused outside Corpus.

#### `DocumentBrowser`

Status: **mostly compliant**.

Recommended cleanup:

- Add `index.ts` inside the new folder if missing.
- Review inline styles; move static height/flex story wrappers to story CSS only if they grow.

#### `DocumentInspector`

Status: **partially compliant**.

Issues:

- Has Storybook with Redux provider because it calls RTK Query for processing artifacts.

Recommended extraction:

1. Split artifact fetching into `DocumentInspectorContainer` or parent page.
2. Make `DocumentInspector` fully DTO-shaped.
3. Keep current story until split is done.
4. Add specific artifact-loading/error stories after split.

#### `IdentityBar`

Status: **non-compliant and high priority**.

What it is:

- The embedding identity selector for Corpus exploration.
- It controls strategy/provider/model/dimensions and displays corpus totals.

Issues:

- No story.
- Raw `.panel`, `.panel-header`, `.panel-body-condensed`.
- Raw `text-mono text-dim` labels.
- Many inline flex styles.

Recommended extraction/cleanup:

1. Add `IdentityBar.stories.tsx` with:
   - default identity;
   - no strategies;
   - sparse embedding totals;
   - changed dimensions/model.
2. Refactor with:
   - `Panel` for shell;
   - `Inline` for wrapping selector row;
   - `FormRow` or local selector row component for labeled controls;
   - `Caption`/`CodeText` for labels and values;
   - `MetadataGrid` for totals if the layout remains compact.
3. Move remaining selector-row anatomy into `IdentityBar.module.css`.
4. Consider future molecule `IdentitySelectorBar` only if Embeddings/Corpus/Artifacts converge on the same shape.

#### `ArtifactIdentityBar`

Status: **non-compliant and high priority**.

What it is:

- A pair of identity selectors for document preprocessing and chunk enrichment coverage.
- It determines which artifact profile/version the Corpus source panel and document inspector are showing.

Issues:

- No stories.
- Two exported components in one file with ad-hoc inline layout.
- Raw `text-mono text-dim` labels.
- Not folder/story/index complete.

Recommended extraction:

1. Split into:
   - `DocProcessingIdentityBar` or `DocumentProcessingIdentityBar`;
   - `ChunkEnrichmentIdentityBar`.
2. Add stories for both:
   - identities available;
   - no identities;
   - profile/version changes;
   - selected stale/missing profile if representable.
3. Use `Panel` only if each bar is a standalone panel; otherwise use `Inline`, `FormRow`, `Caption`, and `CodeText` inside the parent page.
4. Move static layout to CSS Modules.

#### `ChunkTimelineBar`

Status: **non-compliant and high priority**.

Issues:

- No story.
- Uses global `chunk-bar-container`, `chunk-bar`, `embedded`, `not-embedded`, `selected` classes.
- Inline dynamic geometry is valid, but static styles should be local.

Recommended extraction/cleanup:

1. Add `ChunkTimelineBar.module.css`.
2. Move global chunk timeline classes from `index.css` into the module.
3. Keep inline `left`/`width` because those are dynamic geometry.
4. Add stories:
   - all embedded;
   - mixed coverage;
   - selected chunk;
   - no chunks returns null can be covered with a small wrapper note if useful.

### Workflows

#### `WorkflowsView`

Status: **container, partially compliant**.

Issues:

- Many panels extracted, but runtime view still owns modal, detail composition, inline styles, and no page story.

Recommended extraction:

1. Extract `WorkflowsPage` DTO-shaped composition from `WorkflowsView`.
2. Extract `SubmitIntakeWorkflowForm` from `SubmitIntakeModal` if a storyable form is useful.
3. Keep API mutations/polling in `WorkflowsView`.
4. Add page stories:
   - no workflows;
   - running workflow selected;
   - failed op selected;
   - queue backlog;
   - submit modal open.

#### Workflow organisms

Status: **mostly compliant**.

Cleanup:

- `WorkflowOpGroupsPanel`: replace `text-dim text-mono` empty state with `Caption`/`CodeText`.
- `WorkflowSummaryPanel`: verify progress CSS is local; if still global, move to module or future `ProgressBar`.
- `WorkflowListPanel`: move static inline widths/spacing if they are not dynamic.

### Embeddings

#### `EmbeddingsView`

Status: **good primitive adoption, missing page story boundary**.

Recommended extraction:

1. Keep `EmbeddingsView` as container.
2. Extract `EmbeddingsPage` or `EmbeddingsInspectorPage` with DTO-shaped props for strategies, documents, chunks, compute result, similarity result, loading/error flags, and callbacks.
3. Story states:
   - initial/default;
   - compute result;
   - similarity matches;
   - sparse/no chunks;
   - API errors.
4. Keep `fieldset` styling local unless repeated elsewhere.

### Pipeline

#### `PipelineOverview` and `PipelinePage`

Status: **compliant**.

Recommended cleanup:

- Good model for future page splits.
- Add `index.ts` under `components/pipeline` only if it becomes a public package surface.

#### `PipelineView`

Status: **container wrapper, acceptable without story**.

Recommendation:

- No story needed for pure RTK Query container if `PipelinePage` and `PipelineOverview` are covered.

### Evaluation

#### `EvaluationView`

Status: **non-compliant placeholder**.

Issues:

- Raw global panel classes.
- No story/page boundary.

Recommended cleanup:

1. Replace with `EvaluationPage` placeholder using `Panel`, `Text`, and `Caption`.
2. Add `Pages/EvaluationPage` story.
3. Later define actual evaluation organisms when backend concepts are ready.

### Retro components

Files:

- `retro/MacButton.tsx`
- `retro/MacMenuBar.tsx`
- `retro/MacWindow.tsx`

Status: **legacy/unclear**.

Recommendation:

- Decide whether these are still part of the system.
- If retained, move into the design-system hierarchy with stories and module CSS.
- If unused, delete them in a focused cleanup commit.

## Global CSS Deletion Plan

Do not delete global CSS first. Delete only after active consumers are gone.

Current blocks to retire:

```text
.panel / .panel-header / .panel-body / .panel-body-condensed
.data-table
.status-*
.text-mono / .text-dim / .text-small / .text-bold / .text-content
.meta-grid
.tab-bar / .tab-item
.stat-grid / .stat-label / .stat-value
.coverage-strip / .coverage-dot
.chunk-bar-container / .chunk-bar
.op-graph / .op-node / .op-arrow / .op-fan-out
.form-row / .form-section
.fieldset if no longer used globally
```

Suggested order:

1. Remove component/story consumers.
2. Run `rg` for each class family.
3. Delete the global block.
4. Run `pnpm typecheck`, `pnpm build`, and `pnpm build-storybook`.

## Post-Foundation Extraction Roadmap

After foundation overview stories are added, walk each area in this order:

1. Corpus identity widgets:
   - `IdentityBar`
   - `ArtifactIdentityBar`
   - `ChunkTimelineBar`
2. Search page split:
   - storyable `SearchWorkbenchPresentation` or equivalent;
   - remove fallback raw panel.
3. Workflows page split:
   - storyable `WorkflowsPage`;
   - optional storyable submit intake form.
4. Embeddings page split:
   - storyable `EmbeddingsPage`.
5. Evaluation placeholder:
   - storyable `EvaluationPage`.
6. Retro component decision:
   - story and keep, or delete.
7. Global CSS deletion:
   - delete only unused blocks.

## Validation Checklist

For each cleanup batch:

```bash
cd 2026-05-27--rag-evaluation-system/web
pnpm typecheck
pnpm build
pnpm build-storybook
cd ..
git checkout -- internal/web/dist/index.html
```

Ticket hygiene:

```bash
cd 2026-05-27--rag-evaluation-system
docmgr --root ttmp doctor --ticket RAG-WEB-DESIGN-SYSTEM-REVIEW --stale-after 30
```
