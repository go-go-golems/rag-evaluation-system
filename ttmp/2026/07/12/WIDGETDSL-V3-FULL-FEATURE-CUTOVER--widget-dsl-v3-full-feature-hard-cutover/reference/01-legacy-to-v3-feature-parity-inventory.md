---
Title: Legacy to v3 Feature Parity Inventory
Ticket: WIDGETDSL-V3-FULL-FEATURE-CUTOVER
Status: active
Topics:
    - widget-dsl
    - ui-dsl
    - widget-ir
    - goja
    - xgoja
    - react
    - design-system
    - frontend-architecture
    - typescript
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://cmd/widgetdsl-migration-checker
      Note: Parser-backed legacy import and raw escape-hatch checker
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/engines.ts
      Note: Generic data time CRM and ActivityFeed engine contracts
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/props.ts
      Note: Widget IR public prop contracts
    - Path: repo://pkg/widgetdsl/module.go
      Note: Legacy helper maps recipes and split-module installation evidence
    - Path: repo://pkg/widgetdsl/testdata/v3/examples
      Note: Golden v3 example corpus and raw-use migration evidence
    - Path: repo://pkg/widgetdsl/typescript.go
      Note: Current handwritten TypeScript declarations and builder surface
    - Path: repo://pkg/widgetdsl/v3_descriptors.go
      Note: Partial v3 descriptor inventory and parity-generation foundation
ExternalSources: []
Summary: Evidence-backed inventory of legacy split-module helpers, current Widget DSL v3 namespaces, React Widget adapters, raw escape-hatch use, and the hard-cutover disposition of each capability family.
LastUpdated: 2026-07-12T19:35:00-04:00
WhatFor: Use this inventory to decide whether a legacy capability should become a typed v3 primitive, a semantic domain view, an engine-level helper, an internal lowering detail, or be removed.
WhenToUse: During Widget DSL v3 parity implementation, migration reviews, descriptor work, example rewrites, and the final removal of legacy provider modules.
---


# Legacy to v3 Feature Parity Inventory

## Goal

This document answers a narrower question than the primary design guide: what exists today, where does it exist, what is missing from v3, and what should happen to it during the hard cutover?

“Full feature set” does **not** mean copying every legacy component factory into the top level of v3. The legacy modules expose 87 direct component helpers because they grew as component catalogs. V3 should preserve the useful behavior while reorganizing it into a coherent grammar:

1. typed generic primitives under `widget.ui`;
2. schemas, collections, engines, and shaping under `widget.data`;
3. task-level views and intents under domain namespaces;
4. `widget.raw` only for experimental or genuinely unmodeled components.

## Evidence Baseline

The generated inventory in `sources/01-generated-runtime-inventory.md` reports:

- 41 direct `ui.dsl` component helpers;
- 1 direct `data.dsl` component helper plus cell helpers and recipes;
- 20 direct `context_window.dsl` component helpers;
- 11 direct `course.dsl` component helpers;
- 14 direct `cms.dsl` component helpers;
- 17 current `widget.ui` exports;
- 6 current `widget.data` exports;
- semantic v3 namespaces for CMS, course, context, schedule, time, and CRM;
- 87 registered React Widget adapters across registry groups.

The migration checker reports 11 raw-component findings in the 41 checked v3 golden examples. Several findings target components for which typed v3 APIs now exist, demonstrating example/runtime drift. Others identify genuine v3 gaps such as Markdown article and upload-drop-area helpers.

## Status Vocabulary

| Status | Meaning |
|---|---|
| Native v3 | A typed v3 helper or builder exists and should remain. |
| Semantic v3 | V3 intentionally represents the capability through a higher-level view. |
| Partial v3 | Some behavior exists, but important options or typed methods are missing. |
| Raw only | React/Widget IR supports it, but v3 requires `widget.raw.component`. |
| New design | No adequate current contract; implement the proposed typed design. |
| Internal only | Keep as a lowering target but do not expose as everyday public vocabulary. |
| Remove | Do not carry the legacy API forward. |

## Core Module and Transport Surface

| Capability | Legacy | Current v3 | Hard-cutover disposition |
|---|---|---|---|
| Page wrapper | `ui.page` | `widget.page` builder | Native v3; delete legacy form after migration. |
| Text node | shared `text` | `widget.raw.text` | Keep raw kernel; add `widget.ui.text` for styled `Text`. |
| Element node | shared `element` | `widget.raw.element` | Keep only as explicit low-level escape hatch. |
| Component node | shared `component` | `widget.raw.component` | Keep escape hatch; CI should budget/document uses. |
| Fragment | shared `fragment` | `widget.raw.fragment` | Keep kernel helper. |
| Actions | `action.server/navigate/download/event/copy` | `widget.act.*` | Native v3; extend action union rather than adding ad hoc events. |
| Bindings | ad hoc paths/templates | `widget.bind.field/path/map/template/context/const` | Native v3; make one canonical binding vocabulary. |
| Runtime module packaging | seven provider modules | one `widget.dsl` plus legacy modules | Final provider should expose only `widget.dsl`. |
| TypeScript declarations | per-module handwritten descriptors | one v3 descriptor plus handwritten declaration lines | Replace partial descriptor inventory with complete parity enforcement. |

## Generic UI Parity

### Structure and layout

| Legacy helper | Current v3 expression | Status | Target |
|---|---|---|---|
| `appShell` | Page builder/default app shell | Semantic v3 | Keep AppShell internal; expose page shell policy, not raw AppShell props. |
| `appNav` | Page breadcrumbs/navigation or domain shell | Partial v3 | Add typed page navigation builder if general app navigation remains needed. |
| `dashboardGrid` | Section metrics or raw | Partial v3 | Add `ui.dashboard`/metric composition only if not fully covered by page section metrics. |
| `fieldGrid` | Collection/editor lowerer | Internal only | Keep as form-layout lowering; optional typed layout helper for advanced forms. |
| `inline` | `widget.ui.inline` | Native v3 | Keep. |
| `panel` | `widget.ui.card`/`callout` | Semantic v3 | Do not restore generic `panel`; document card versus callout intent. |
| `scrollRegion` | raw | Raw only | Add `widget.ui.scroll` because wide tables and bounded panes need it. |
| `sectionBlock` | `widget.page(...).section(...)` | Semantic v3 | Keep SectionBlock internal. |
| `sidebarShell` | course shell or raw | Partial v3 | Add a general typed app/sidebar shell builder if non-course apps need it. |
| `splitPane` | `widget.ui.splitPane` | Native v3 | Keep. |
| `stack` | `widget.ui.stack` | Native v3 | Keep. |
| `tabList` | raw | Raw only | Add typed `widget.ui.tabs` with action context. |
| `tileGrid` | domain media views or raw | Partial v3 | Keep internal for CMS; expose `data.cards`/tile arrangement only when collection grammar supports it. |

### Text, status, and compact data display

| Legacy helper | Current v3 expression | Status | Target |
|---|---|---|---|
| `caption` | `widget.ui.caption` | Native v3 | Keep. |
| `codeText` | raw | Raw only | Add `widget.ui.code`. |
| `divider` | raw | Raw only | Add `widget.ui.divider` or section policy; avoid repeated raw component calls. |
| `statusText` | `widget.ui.status` | Native v3 | Keep semantic name. |
| `textBlock` | raw text or raw `Text` | Partial v3 | Add `widget.ui.text`; raw text nodes cannot express typography roles. |
| `tag` | `widget.ui.badge` | Semantic v3 | Decide whether name should be `tag` or `badge`; one term only after cutover. |
| `meterBar` | raw/domain views | Raw only | Add typed `widget.ui.meter` or `data.measure` when used independently. |
| `keyValueStrip` | section metrics or raw | Partial v3 | Add `widget.ui.summary` if metric builders do not cover compact key/value strips. |
| `metadataGrid` | `widget.ui.metadata` | Semantic v3 | Keep. |
| `personSummary` | raw | Raw only | Add only if used outside CRM/domain views. |

### Content/document composition

| Legacy helper | Current v3 expression | Status | Target |
|---|---|---|---|
| `markdownArticle` | raw `MarkdownArticle` | Raw only | Add `widget.ui.markdownArticle(source, configure?)`. |
| `richArticle` | course/domain view or raw | Partial v3 | Add `widget.ui.richArticle(blocks, configure?)`; course may wrap it. |
| `figureBlock` | raw | Raw only | Add typed content helper if article/handout slots need it. |
| `checkList` | raw | Raw only | Add `widget.ui.checkList`. |
| `stepList` | raw | Raw only | Add `widget.ui.stepList` or keep domain-specific if usage is narrow. |
| `keyPointList` | raw | Raw only | Consider one generalized semantic list builder instead of three overlapping list APIs. |

### Forms, search, and navigation

| Legacy helper | Current v3 expression | Status | Target |
|---|---|---|---|
| `formPanel` | `widget.ui.form` | Native v3 | Keep; improve typed submit/action result integration. |
| `formRow` | `widget.ui.formRow` | Native v3 | Keep. |
| `textInput` | `widget.ui.textInput` | Native v3 | Keep. |
| `textareaInput` | `widget.ui.textareaInput` | Native v3 | Keep. |
| `selectInput` | `widget.ui.selectInput` | Native v3 | Keep. |
| `searchField` | CMS views or raw | Raw only | Add collection-level search builder; keep SearchField as lowering component. |
| `pagination` | CMS views or raw | Raw only | Add collection-level pagination builder; extend component with page-size control. |
| `breadcrumbs` | page builder `.breadcrumb` | Semantic v3 | Keep page-owned breadcrumbs. |
| `sidebarNav` | domain shell or raw | Partial v3 | General app shell/navigation builder should own it. |
| `uploadDropArea` | CMS media view or raw | Raw only | Add `widget.ui.upload` or a typed file-input intent shared by CMS/context. |
| `emptyState` | `widget.ui.emptyState` | Native v3 | Keep. |

## Data Grammar Parity

### Collections and records

| Legacy capability | Current v3 | Status | Target |
|---|---|---|---|
| Direct `dataTable(props)` | `data.collection(...).table()` | Semantic v3 | Do not restore raw table as normal API. |
| Ordered schema and field roles | `data.fields` | Native v3 | Keep and extend field types/options. |
| Selectable table | collection selection + table row action | Native/partial | Add keyboard policy and stable focus behavior. |
| Master-detail | collection `.masterDetail()` | Native v3 | Keep; add named detail slots rather than raw callback leakage. |
| Create/edit/remove/reorder | collection edit/actions | Partial v3 | Preserve typed actions; improve method consistency. |
| URL selection | `data.selection.urlParam` | Native v3 | Keep as explicit server-visible state. |
| Search/filter shaping | none on collection | New design | Add `collection.search(...)` and typed ordered filter specs. |
| Pagination/page size | none on collection | New design | Add `collection.paginate(...)` and server page-result contract. |
| Table keyboard navigation | none | New design | Add scoped keyboard builder and row command specs. |
| Conditional semantic styling | none | New design | Add predicates and semantic tone/decoration rules, not arbitrary CSS. |
| Column preferences | URL workaround/raw | New design | Add client presentation preference contract after collection core. |
| Card/tile arrangement | isolated domain widgets | Partial v3 | Extend arrangement grammar instead of adding one-off collection APIs. |
| Board arrangement | `crm.pipelineBoard` | Domain v3 | Keep CRM view; generic BoardEngine remains data engine. |
| Matrix | `data.matrix` | Native v3 | Keep engine helper. |
| Record timeline | `crm.activityFeed` | Native but misplaced | Promote to `data.activityFeed` or `data.timeline`; CRM supplies presets. |

### Cell renderer parity

| Legacy cell | Current v3 | Target |
|---|---|---|
| `field` | Native | Keep. |
| `status` | Native | Keep. |
| `template` | Native | Keep, but prefer typed bindings where possible. |
| `number` | Schema-derived only | Add engine-level typed cell for explicit tables/matrices. |
| `caption` | Schema-derived only | Add semantic caption cell. |
| `link` | Missing | Add typed link cell. |
| `linkButton` | Missing | Add typed link-button cell if action links remain distinct. |
| `actionButton` | Collection action columns | Keep internal plus typed action-column builder. |
| `constant` | `value` is not equivalent in all contexts | Add explicit constant/renderable cell. |
| `cycle` | Native v3 | Keep for matrices. |
| `value` | Native v3 | Keep for explicit matrix values. |

## Domain Namespace Parity

### Context

V3 correctly replaces many component factories with two primary views:

- `widget.context.diagram(snapshot, builder)`;
- `widget.context.workspace(session, builder)`.

The following legacy capabilities remain available only through raw components or as internal details: anchored comment card/rail, annotation badge/note/rail, individual budget/strip/grouped-strip/stack/treemap renderers, context turn pager, upload drop area, transcript message card, transcript role badge, transcript session header, and reader panel.

Hard-cutover policy:

- keep individual diagram molecules internal to `context.diagram`;
- expose annotation/comment/transcript subviews only when they are legitimate named slots or standalone tasks;
- share upload through generic UI/file input instead of a context-specific duplicate;
- ensure every context view carries `styleSet` and uses `styleKey`.

### Course

V3 semantic views are stronger than the legacy factories:

- `course.shell` replaces direct course studio shell construction;
- `course.landing` replaces hand-composed lesson landing panels;
- `course.slideDeck` replaces direct slide-panel composition;
- `course.handouts` replaces document list/preview shell assembly;
- metadata, agenda, and material helpers cover admin tasks.

Direct legacy helpers such as `CourseStepNav`, `DocumentPreviewToolbar`, `SlideShell`, and `ContextStudioNavIcon` should usually remain lowering details or named slots. Markdown/rich article rendering belongs under generic UI content, not course.

### CMS

V3 semantic views cover the major workflows:

- `cms.mediaLibrary`;
- `cms.articleQueue`;
- `cms.markdownEditor`;
- CMS intents.

Legacy atoms such as MediaThumb, AssetTile, ContentStatusBadge, MeterBar, SearchField, Pagination, Tag, and TileGrid should not all become CMS top-level functions. Generic pieces move to `ui`/`data`; media-specific units become slots or internal lowerings of `mediaLibrary`.

### CRM and activity

CRM is a v3-only namespace. It adds typed fields, pipeline definitions, BoardEngine views, record fields, tasks, stats, funnel, intents, and ActivityFeed.

`ActivityFeed` itself is domain-blind. Its required item contract is:

```ts
interface ActivityFeedItemSpec {
  id: string;
  kind: string;
  title: RenderableValue;
  body?: RenderableValue;
  atISO: string;
  actor?: { id?: string; name: string; avatarUrl?: string };
  meta?: JsonObject;
}
```

Hard-cutover recommendation:

```js
widget.data.activityFeed(activities, feed => feed
  .groupByDay(true)
  .glyph("note", "📝")
  .onOpen(openActivity)
  .onLoadMore(loadEarlier));
```

Then CRM may provide `crm.activityFeed` only if it adds CRM vocabulary or defaults. Since v3 is unreleased, the cleaner option is to move the generic helper now and update examples. The React adapter already lives in `dataWidgetRegistry`, reinforcing generic ownership.

### Schedule and time

These are new v3 strengths with no legacy equivalent:

- availability poll, summary, and booking picker;
- month/week views and time formatting/ranges;
- MatrixGrid, MonthGrid, and TimeGrid lowerings;
- schedule/time intents.

They establish the desired pattern: product-level view, scoped builder, domain intents, generic engine underneath.

## Cross-cutting Capabilities Missing from Both Public Surfaces

These are not simple parity items; they require coordinated renderer and DSL design:

- keyboard row navigation and row-scoped command bindings;
- generic controlled FormDialog/overlay lifecycle;
- structured navigate actions with typed query parameters;
- browser-local presentation preferences;
- progressive search panels with draft versus submitted state;
- page-size selection and server collection pagination;
- conditional semantic row/cell styling;
- first-class toast/live-region and field-error result handling;
- focus restoration across page refresh and mutation;
- command discovery/help;
- generic action undo policy.

The primary design guide incorporates these requirements into the v3 hard-cutover rather than treating them as application hacks.

## Current Raw Escape-Hatch Findings

The migration checker found 11 raw uses in v3 examples:

| Example | Raw target | Interpretation |
|---|---|---|
| `08-markdown-article.js` | MarkdownArticle | Genuine missing v3 UI helper. |
| `21-crm-record-fields.js` | RecordFieldList | Stale example; typed CRM helper exists. |
| `22-crm-field-renderers.js` | FieldRenderer | Engine-level helper missing or example should use record fields. |
| `23-crm-board.js` | BoardEngine | Stale example; typed pipeline board exists. |
| `24-activity-feed.js` | ActivityFeed | Stale example; CRM helper exists, generic data helper recommended. |
| `25-ui-form.js` | FormRow/TextInput/TextareaInput | Stale example; typed UI helpers exist. |
| `26-upload-drop-area.js` | ContextUploadDropArea | Genuine missing generic upload helper. |

The example hosts themselves produced no legacy-import or raw-component findings. This means the host path is already close to the intended end state, while the golden corpus needs cleanup and parity enforcement.

## Hard-cutover Exit Criteria

The legacy provider modules can be removed only when:

1. every stable registered Widget component has an explicit disposition in this inventory;
2. every promoted public capability has a typed v3 helper, declaration, descriptor, docs, and tests;
3. semantic views replace legacy recipes in all first-party hosts;
4. golden examples have zero unexplained raw-component findings;
5. migration checker reports no legacy imports in first-party authoring source;
6. provider exposes only `widget.dsl` plus help;
7. legacy TypeScript declarations and docs are deleted or archived;
8. browser interaction tests cover navigation, forms, upload, collections, search, pagination, dialogs, and keyboard commands.

## Key Evidence Files

- `pkg/widgetdsl/module.go`
- `pkg/widgetdsl/v3.go`
- `pkg/widgetdsl/v3_crm.go`
- `pkg/widgetdsl/v3_descriptors.go`
- `pkg/widgetdsl/typescript.go`
- `pkg/widgetdsl/v2/spec/types.go`
- `pkg/widgetdsl/v2/spec/validate.go`
- `pkg/widgetdsl/v2/spec/lower.go`
- `packages/rag-evaluation-site/src/widgets/defaultRegistry.ts`
- `packages/rag-evaluation-site/src/widgets/ir/props.ts`
- `packages/rag-evaluation-site/src/widgets/ir/engines.ts`
- `pkg/xgoja/providers/widgetsite/provider.go`
- `pkg/widgetdsl/testdata/v3/examples/`
- `sources/01-generated-runtime-inventory.md`
- `sources/02-v3-example-migration-check.txt`
- `sources/03-v3-host-migration-check.txt`
