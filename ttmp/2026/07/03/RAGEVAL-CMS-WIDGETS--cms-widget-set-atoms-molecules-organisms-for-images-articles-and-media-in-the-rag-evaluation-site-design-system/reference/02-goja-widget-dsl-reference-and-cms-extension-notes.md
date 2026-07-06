---
Title: Goja Widget DSL reference and CMS extension notes
Ticket: RAGEVAL-CMS-WIDGETS
Status: active
Topics:
    - design-system
    - frontend
    - cms
    - goja
    - dsl
    - widget-ir
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: cmd/rag-eval/xgoja.yaml
      Note: Evidence DSL not selected in-repo
    - Path: internal/api/dsl_handlers.go
      Note: Only in-repo IR producer (hardcoded demo page)
    - Path: packages/rag-evaluation-site/src/components/organisms/MediaLibraryPanel/MediaLibraryPanel.widget.tsx
      Note: cms.dsl flagship adapter (upload serialization + navigate-based state)
    - Path: packages/rag-evaluation-site/src/widgets/WidgetRenderer.cms.stories.tsx
      Note: 7 IR stories rendering the CMS from JSON
    - Path: packages/rag-evaluation-site/src/widgets/uploadSerialization.ts
      Note: Shared File→JSON serialization for upload actions
    - Path: pkg/widgetdsl/module.go
      Note: 'The entire goja DSL: specs'
    - Path: pkg/widgetdsl/module_test.go
      Note: Module-boundary tests (helper presence/absence)
    - Path: pkg/widgetdsl/registrar.go
      Note: engine.RuntimeModuleRegistrar wiring
    - Path: pkg/widgetdsl/typescript.go
      Note: d.ts generation; drift (missing linkButton/actionButton/download)
    - Path: pkg/xgoja/providers/widgetsite/doc/02-widget-dsl-js-api-reference.md
      Note: Canonical JS API reference (drift noted)
    - Path: pkg/xgoja/providers/widgetsite/provider.go
      Note: rag-widget-site xgoja provider + embedded help
ExternalSources: []
Summary: 'Line-anchored reference for the four goja UI DSL modules (ui.dsl, data.dsl, context_window.dsl, course.dsl) derived from the React WidgetRenderer: registration paths, helper surfaces, recipes, TS typing generation, observed drift, and how a cms.dsl module would be added.'
LastUpdated: 2026-07-03T18:30:00-04:00
WhatFor: Understand how JS scripts author Widget IR pages via goja, and what Phase 6 (cms.dsl) of the CMS design concretely requires.
WhenToUse: Read alongside design-doc/01 §4.6/§6.9/Phase 6, or before touching pkg/widgetdsl or the widgetsite xgoja provider.
---



# Goja Widget DSL reference and CMS extension notes

This document covers the JavaScript authoring layer that sits on top of the React `WidgetRenderer`: goja native modules that let scripts (jsverbs, xgoja binaries, embedded engines) build Widget IR JSON without touching React. It complements design-doc/01 §4.6 (Widget IR).

> **Update (2026-07-03, Phase 5 landed):** `cms.dsl` now exists as the **fifth module** (14 helpers: mediaThumb, tag, contentStatusBadge, meterBar, tileGrid, assetTile, breadcrumbs, pagination, searchField, emptyState, markdownEditor, mediaLibraryPanel, articleListPanel, cmsShell; recipes `mediaLibrary`/`articleList`; provider entry in widgetsite). The §5 extension sketch below is implemented. The §4 drift items are **fixed**: typescript.go now declares `cell.linkButton`/`cell.actionButton`/`action.download`; `contextGroupedStripDiagram` has a helper; doc/02 documents all of the above plus the new `confirm` field on every action (`ActionSpec.confirm`, handled centrally in `dispatchWidgetAction` with context interpolation). Boundary tests cover the new module (`module_test.go` TestCmsModuleExportsHelpersRecipesAndBoundaries).

## 1. What it is, in one paragraph

`pkg/widgetdsl` implements four `require()`-able goja modules — `ui.dsl`, `data.dsl`, `context_window.dsl`, `course.dsl` (module.go:14–19) — whose exports are thin factories over the Widget IR node shape `{ kind: "component", type, props?, children? }`. The helper names are camelCase mirrors of the React component types (`ui.panel` → `Panel`, `course.markdownArticle` → `MarkdownArticle`), grouped by the exact same module boundaries as the TS-side registries in `packages/rag-evaluation-site/src/widgets/defaultRegistry.ts` (`uiWidgetRegistry`, `dataWidgetRegistry`, `contextWindowWidgetRegistry`, `courseWidgetRegistry`). JS produces JSON; React owns CSS modules, a11y, and behavior. There is deliberately no "bucket" compatibility module — scripts import the domain module they use.

## 2. Where it lives and how it gets into a runtime

| File | Role |
|---|---|
| `pkg/widgetdsl/module.go` (1104 lines) | The whole DSL: module specs, helper maps, base constructors, `action`/`cell` objects, style helpers, 8 recipes, child/props normalization |
| `pkg/widgetdsl/registrar.go` | `Registrar` (`ID() = "widget-dsl"`) for go-go-goja `engine.NewRuntimeFactoryBuilder()` composition |
| `pkg/widgetdsl/typescript.go` | Generates `.d.ts` module declarations via go-go-goja `tsgen/spec` |
| `pkg/widgetdsl/module_test.go` | Runtime tests: helper presence/absence per module, recipe output shape |
| `pkg/xgoja/providers/widgetsite/provider.go` | xgoja provider, package id `rag-widget-site`: exposes the four modules + embedded help (`doc/*.md`) to generated binaries |
| `pkg/xgoja/providers/widgetsite/doc/01..03-*.md` | Glazed help pages: getting started, JS API reference, SPA bundling |

Three registration paths:

1. **Auto-registration into go-go-goja**: `init()` calls `modules.Register(&module{spec})` for each spec (module.go:181–185), so any go-go-goja engine that scans registered native modules can `require("ui.dsl")`.
2. **Explicit registrar**: `widgetdsl.NewRegistrar()` implements `engine.RuntimeModuleRegistrar` and calls `Register(reg)` on the require registry (registrar.go:14–24).
3. **xgoja buildspec**: external generated binaries list package `rag-widget-site` (import `…/pkg/xgoja/providers/widgetsite`) and select modules by name in `xgoja.yaml` (doc/01:37–57). Each module ships its generated TypeScript typings (`TypeScript: widgetdsl.TypeScriptModule(name)`, provider.go:24–52).

**Observed non-usage (worth knowing):** this repo's own `cmd/rag-eval/xgoja.yaml` does *not* select `rag-widget-site` (its packages are go-go-goja core/host/http, goja-text, geppetto), and none of the `cmd/rag-eval/jsverbs/*.js` scripts require a `*.dsl` module (they use `db`, `fs`, `yaml`, `markdown`, `sanitize`, `express`). So today the DSL's consumers are the tests, the help docs, and *external* xgoja apps — inside this repo the only live IR producer is the hardcoded Go demo page (`internal/api/dsl_handlers.go`). The serving pattern the docs prescribe (`/api/widget/pages/{id}` from a jsverb + SPA fallback, doc/01:154–182) is a design intention, not yet wired here.

## 3. Module surfaces

### 3.1 Shared constructors (all four modules)

Every module exports (module.go:191–216):

- `text(value)` → `{ kind: "text", text: String(value) }`
- `element(tag, attrs?, ...children)` → raw host element node
- `component(type, props?, ...children)` → component node by explicit type string (the escape hatch for types without a helper)
- `fragment(...children)` → normalized child array
- one camelCase helper per component in the module's map (module.go:200–202)
- `action` object (all modules)
- `ui.dsl` only: `page(options)`; `data.dsl` only: `cell` object; `context_window.dsl` only: style helpers; per-module `recipes`

**Props-vs-child disambiguation** (module.go:834–839, 863–870): the first argument after the type is treated as props only if it exports as a plain object *and* does not look like a widget node (`kind` ∈ {text, element, component}). So `ui.panel(ui.caption(...))` works — the caption becomes a child, not props.

**Child normalization** (module.go:796–832): `null`/`undefined` dropped; arrays flattened recursively; widget nodes passed through; everything else stringified into text nodes. This is why `ui.statusText({...}, "Rows: 2")` renders text.

### 3.2 Helper inventory (helper → React component type)

- **`ui.dsl`** (module.go:33–64, 30 helpers): `appShell appNav button caption checkList codeText dashboardGrid divider figureBlock formPanel formRow inline keyPointList keyValueStrip metadataGrid panel personSummary scrollRegion sectionBlock selectInput sidebarNav sidebarShell splitPane stack statusText stepList tabList textBlock(→Text) textInput textareaInput`.
- **`data.dsl`** (module.go:66–68): `dataTable`.
- **`context_window.dsl`** (module.go:70–90, 19 helpers): swatch/badge atoms, all context diagrams (`contextBudgetBar contextStripDiagram contextStackDiagram contextTreemap contextDiagramPanel contextTurnPagerPanel contextLegend`), transcript/annotation surfaces, `contextUploadDropArea`. **Gap:** `ContextGroupedStripDiagram` has a React widget adapter but *no* DSL helper — reachable only via `component("ContextGroupedStripDiagram", …)`.
- **`course.dsl`** (module.go:92–104, 11 helpers): `contextStudioNavIcon courseLessonPanel courseSlidePanel courseStepNav courseStudioShell documentListPanel documentPreviewToolbar handoutDocumentShell markdownArticle richArticle slideShell`.

Module boundaries are enforced by tests: `ui.dataTable`, `data.page`, and `contextWindow.contextStudioNavIcon` must be `undefined` (module_test.go:15–60).

### 3.3 `ui.page(options)` (module.go:518–551)

Returns `{ schemaVersion: "0.1.0"?, id, title, meta?, root }`. If `root` is a widget node, it is used as-is; otherwise `sections: WidgetNode[]` are wrapped in a `Stack` with `gap` (default `"lg"`). This output is exactly the `WidgetPageResponse` shape `useWidgetPage`/`getDslPage` consume.

### 3.4 `action.*` (module.go:264–290)

`action.server(name, opts?)`, `.navigate(to, opts?)`, `.download(to, opts?)`, `.event(event, opts?)`, `.copy(value)` — producing the same `ActionSpec` union the React dispatcher handles (`packages/rag-evaluation-site/src/widgets/actions.ts:24–96`). Recipes additionally accept a bare string where an action is expected and normalize it to `{ kind: "server", name }` (`normalizeActionSpec`, module.go:1001–1022) — e.g. `onSelect: "handout-select"`.

### 3.5 `data.cell.*` (module.go:218–262)

`field number status caption template link linkButton actionButton constant` — mirroring the TS `CellSpec` union (ir.ts:561–634) rendered by `cellRenderers.tsx`. Callback `cell` functions cannot cross the JSON boundary; specs are mandatory.

### 3.6 Context style helpers (`context_window.dsl` only, module.go:292–342)

`visualStyle(opts)` (defaults `pattern:"none"`, `fill:"var(--mac-surface)"`), `legendItem(id,label,opts?)` (supports `hidden:true`), `styleSet(opts)` (normalizes empty `legend`/`styles`), `contextPart(id,label,styleKey,tokens,opts?)`, `contextSnapshot(opts)` (defaults id/title/limit/parts), and `paletteStyleSet({palette, entries})`.

`paletteStyleSet` (module.go:344–423) is the interesting one: four named palettes hardcoded in Go — `"Dusty Magenta / Blue"` (default), `"Signal Orange / Cyan"`, `"Slate / Coral"`, `"Cobalt / Sand"` — each defining paper/ink/grid/shadow + 3 accents. Entries pick `accent: a|b|c|grid|shadow|ink`, a `.pattern_*` name, and fill/line percentages; fills are emitted as CSS `color-mix(in srgb, <accent> N%, <paper>)` strings, solid entries get white labels. This enforces the GUIDELINES.md contract ("DSL examples should use `contextWindow.styleSet(...)` or `paletteStyleSet(...)`"; `styleKey + styleSet`, never `kind`).

### 3.7 Recipes (module.go:441–782)

Recipes are Go-side macro expansions returning full IR subtrees (higher-level "semantic recipes" per GUIDELINES.md:268):

| Recipe | Module | Expands to |
|---|---|---|
| `recipes.metrics({items, recipe?})` | ui | `DashboardGrid` of condensed `Panel`s each holding a `StatusText` metric (:553–583); grid recipe auto-picks two-up/three-up/four-up by count (:991–999) |
| `recipes.actionToolbar({title?, actions, caption?})` | ui | `Panel > Inline` of `Button`s with normalized actions (:585–625) |
| `recipes.masterDetailTable({rows, columns, selectedKey?, detail?, …})` | data | `DashboardGrid` of a `Panel>DataTable` plus a detail node; `detail` may be a **JS callback invoked at build time** with the selected row (:734–782) — the one place a function is allowed, because it runs before serialization |
| `recipes.contextDiagram({snapshot, styleSet | palette+entries, view?})` | context_window | `ContextDiagramPanel` (:627–648); throws if neither styleSet nor palette+entries given |
| `recipes.annotatedTranscript({transcript?|messages…, onAnnotationSelect?})` | context_window | `TranscriptWorkspacePanel` with `onAnnotationSelectAction` (:650–667) |
| `recipes.courseStudio({sections, main?, onNavigate?})` | course | `CourseStudioShell` wrapping `main` (:669–686) |
| `recipes.courseSlide({slide, snapshot, …, onNext?})` | course | `CourseSlidePanel` with `on*Action` props (:688–711) |
| `recipes.handout({bundle?|documents…, onSelect?, onDownload?…})` | course | `HandoutDocumentShell` with `onDocumentSelectAction` etc. (:713–732) — **the read-side CMS entry point that already exists** |

Note the callback-to-action naming convention: React props `onDocumentSelect` (function) become IR props `onDocumentSelectAction` (ActionSpec); the widget adapters translate.

### 3.8 TypeScript typings (`typescript.go`)

`TypeScriptModule(name)` emits a `spec.Module` with raw `.d.ts` lines: loose base types (`WidgetNode { kind: string; [key: string]: any }`), `page` when applicable, one `export function <helper>(props?, ...children): WidgetNode` per helper (sorted), plus `cell`/`action`/style-helper/recipe declarations. Props are intentionally untyped (`Props = Record<string, any>`) — "component props remain open-ended by design" (typescript.go:10–13).

## 4. Observed drift (documentation/typing bugs to fix opportunistically)

1. **`typescript.go` omits runtime exports**: the `cell` declaration lacks `linkButton` and `actionButton` (both exist at module.go:248–257), and the `action` declaration lacks `download` (exists at module.go:276–280). Generated typings will mark valid calls as errors.
2. **`doc/02-widget-dsl-js-api-reference.md` omits `action.download`** (its Actions section lists server/navigate/event/copy) and omits `cell.linkButton`/`cell.actionButton` from the data.dsl list.
3. **`ContextGroupedStripDiagram` has no helper** in `contextWindowHelpers` despite having a React widget adapter and Storybook coverage ("Context Grouped Strip By Turn") — either add the helper or document `component()` usage.
4. **The DSL is dormant in-repo**: not selected in `cmd/rag-eval/xgoja.yaml`, unused by jsverbs, and the only served page is Go-hardcoded (`internal/api/dsl_handlers.go` recognizes id `"demo"` only). Anyone "testing the DSL end-to-end" must either run `pkg/widgetdsl` tests or build an external xgoja app.
5. `action.event` TS signature says `event(name, options?)` while the docs' runtime uses `{ detail }` in options — consistent, but note `event` special-cases `"print"`/`"fullscreen"` in the React dispatcher (actions.ts:43–61).

## 5. What this means for the CMS (`cms.dsl`, design-doc Phase 6)

Adding the CMS module is mechanical once the React organisms are stable. Concretely:

```go
// pkg/widgetdsl/module.go
const CmsModuleName = "cms.dsl"

var cmsHelpers = map[string]string{
    "mediaThumb":         "MediaThumb",
    "tag":                "Tag",
    "contentStatusBadge": "ContentStatusBadge",
    "meterBar":           "MeterBar",
    "assetTile":          "AssetTile",
    "uploadQueueList":    "UploadQueueList",
    "mediaLibraryPanel":  "MediaLibraryPanel",
    "articleListPanel":   "ArticleListPanel",
    "assetDetailPanel":   "AssetDetailPanel",
    "cmsShell":           "CmsShell",
    // deliberately NOT: MarkdownEditor / ArticleEditorPanel / AssetPickerDialog —
    // editing surfaces are interactive-state-heavy; expose them to IR only if a
    // server-action round-trip story exists (see D-6 in design-doc/01).
}

moduleSpecs = append(moduleSpecs, moduleSpec{
    name: CmsModuleName, helpers: cmsHelpers, action: true,
    recipes: []string{"mediaLibrary", "articleList"},
    doc: "cms.dsl provides media, asset, and article-management helpers.",
})
```

Plus, following the existing pattern end to end:

1. **Recipes**: `recipes.mediaLibrary({assets, query?, page?, onAssetSelect?, onFilesSelected?})` → `MediaLibraryPanel` with `on*Action` props; `recipes.articleList({articles, statusFilter?, onRowAction?})` → `ArticleListPanel`. Callback→`on*Action` normalization via `normalizeActionSpec` exactly like `handoutRecipe` (module.go:713–732).
2. **TS side**: new `cmsWidgetRegistry` in `defaultRegistry.ts`, `"cms.dsl"` added to the `WidgetModule` union (registry.ts:5), `X.widget.tsx`/`.widget.yaml` per component translating `onAssetSelectAction: ActionSpec` → `onAssetSelect: (id) => dispatch`.
3. **Provider**: one more `providerapi.Module` entry in `pkg/xgoja/providers/widgetsite/provider.go` with `TypeScript: widgetdsl.TypeScriptModule(CmsModuleName)`.
4. **Docs**: extend `doc/02-widget-dsl-js-api-reference.md` (and fix the §4 drift while there).
5. **Tests**: extend `module_test.go` boundary assertions (`ui.mediaLibraryPanel` must be `undefined`, etc.).
6. **Actions caveat**: DSL-driven CMS pages can *select/navigate/download* fine, but *mutations* (upload, save, publish) need the unimplemented `/api/widget/actions/{name}` route in Go — which is why design-doc D-6 routes CMS mutations through RTK Query containers first. Upload is doubly special: `File` objects cannot cross the JSON boundary at all, so `contextUploadDropArea`'s pattern (files surface in the *action context* `{files, fileNames, fileCount}`, doc/02:219–223) only works with a real handler behind it.

Target authoring experience once wired:

```js
const ui = require("ui.dsl")
const cms = require("cms.dsl")

return ui.page({
  id: "media",
  title: "Media library",
  sections: [
    cms.recipes.mediaLibrary({
      assets: db.query("SELECT … FROM media_assets ORDER BY updated_at DESC"),
      onAssetSelect: ui.action.navigate("/pages/asset-$value"),
    }),
  ],
})
```

## 6. File references

- `pkg/widgetdsl/module.go` — module names :14–19; helper maps :33–104; specs :106–137; install :191–216; cell :218–262; action :264–290; style helpers :292–342; paletteStyleSet + palettes :344–402; recipes :441–782; page :518–551; child/props normalization :796–870; normalizeActionSpec :1001–1022.
- `pkg/widgetdsl/typescript.go` — d.ts generation (drift: missing linkButton/actionButton/download).
- `pkg/widgetdsl/registrar.go` — engine registrar (`widget-dsl`).
- `pkg/widgetdsl/module_test.go` — module-boundary and recipe tests.
- `pkg/xgoja/providers/widgetsite/provider.go` — `rag-widget-site` provider + help source.
- `pkg/xgoja/providers/widgetsite/doc/01-widget-dsl-getting-started.md`, `02-widget-dsl-js-api-reference.md`, `03-widget-dsl-spa-bundling.md` — embedded help.
- `cmd/rag-eval/xgoja.yaml`, `cmd/rag-eval/jsverbs/*.js` — evidence the DSL is not yet consumed in-repo.
- `packages/rag-evaluation-site/src/widgets/{ir.ts,actions.ts,registry.ts,defaultRegistry.ts}` — the TS mirror the DSL targets.
- `internal/api/dsl_handlers.go` — the only in-repo IR producer today (hardcoded "demo").
