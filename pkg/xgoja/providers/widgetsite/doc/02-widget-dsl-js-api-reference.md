---
Title: "Widget DSL JavaScript API Reference"
Slug: widget-dsl-js-api-reference
Short: "Reference for Widget DSL modules, helpers, recipes, actions, and table cell specifications."
Topics:
- xgoja
- widget-dsl
- widget-ir
- react
- javascript
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: Reference
---

This reference documents the JavaScript API exposed by the `rag-widget-site` provider. The API creates JSON-compatible Widget IR consumed by the React WidgetRenderer.

For new code, prefer the parallel `widget.dsl` v3 module. The split module reference below remains useful for existing scripts during migration. The v3 namespace inventory is generated in the ticket API reference at `reference/05-widget-dsl-v3-api-reference.md`.

## Module names

| Module | Purpose |
| --- | --- |
| `widget.dsl` | Preferred v3 module with `raw`, `act`, `bind`, `ui`, `data`, `cms`, `course`, `context`, `schedule`, and `time` namespaces. |
| `ui.dsl` | Generic page, layout, primitive, foundation, and UI recipe helpers. |
| `data.dsl` | Legacy/current data-display widgets, `cell.*` helpers, and v1 data recipes. |
| `data.v2.dsl` | Experimental hard-cutover typed/fluent builders for schemas, tables, selectable tables, master-detail editors, and row actions. |
| `context_window.dsl` | Context-window diagrams, transcript, annotation, anchored-comment, and upload helpers. |
| `course.dsl` | Course, lesson, slide, handout, and course-studio helpers. |
| `cms.dsl` | Media, asset, and article-management helpers. |

There is no compatibility bucket module. New v3 pages use one import:

```js
const widget = require("widget.dsl")
```

Legacy split-module pages import the domains they use explicitly:

```js
const ui = require("ui.dsl")
const data = require("data.dsl")
const dataV2 = require("data.v2.dsl")
const contextWindow = require("context_window.dsl")
const course = require("course.dsl")
const cms = require("cms.dsl")
```

## Shared constructors

Each module exports these low-level constructors:

| Helper | Description |
| --- | --- |
| `text(value)` | Creates `{ kind: "text", text: String(value) }`. |
| `element(tag, attrs?, ...children)` | Creates a plain host element node. |
| `component(type, props?, ...children)` | Creates a component node by explicit type string. |
| `fragment(...children)` | Normalizes child values into an array of Widget IR child nodes. |

`ui.dsl` additionally exports `page(options)`. `page(...)` owns the page wrapper and creates a `Stack` root from `sections` when no explicit `root` is supplied.

## Actions

All four domain modules export `action` because actions can be attached to widgets across domains:

```js
ui.action.server("refresh", { payload: { force: true } })
ui.action.navigate("/pages/$value")
ui.action.download("/api/handouts/$value/download.md")
ui.action.event("widget:selected", { detail: { source: "story" } })
ui.action.copy("copy this")
```

Every action accepts an optional `confirm` string (via the options object). When present the renderer shows a confirmation prompt — with `${path}` / `$name` interpolation against the action context — before dispatching. Use it for destructive actions:

```js
ui.action.server("admin-delete-course-material", { confirm: "Delete ${file}?" })
```

## `ui.dsl` helpers

`ui.dsl` exports generic visual helpers:

- `appShell`, `appNav`
- `button`, `caption`, `codeText`, `divider`, `statusText`, `textBlock`, `tag`, `meterBar`
- `inline`, `stack`, `dashboardGrid`, `panel`, `scrollRegion`, `sectionBlock`, `fieldGrid`, `tileGrid`, `sidebarShell`, `splitPane`
- `formRow`, `selectInput`, `textInput`, `textareaInput`, `tabList`, `searchField`, `uploadDropArea`
- `metadataGrid`, `keyValueStrip`, `checkList`, `stepList`, `personSummary`, `figureBlock`, `keyPointList`, `sidebarNav`
- `breadcrumbs`, `pagination`, `emptyState`, `markdownArticle`, `richArticle`

Generic primitives that used to live only in `cms.dsl` (`tag`, `meterBar`, `tileGrid`, `searchField`, `breadcrumbs`, `pagination`, `emptyState`), `context_window.dsl` (`uploadDropArea`, formerly `contextUploadDropArea`), and `course.dsl` (`markdownArticle`, `richArticle`) are now first-class `ui.dsl` exports. The old module-local names still work but are deprecated aliases — import from `ui.dsl` in new code.

### `ui.section(title, options?, ...children)` — flat sectioning

Document structure without boxes: an uppercase label, an optional 1px rule, and content. Use `section` for page structure; reserve `panel` for interactive tools and selected-item cards.

- `options.level` — `1 | 2 | 3` heading scale (default 1)
- `options.anchor` — DOM id for in-page links
- `options.caption` — muted description line under the label
- `options.actions` — widget node(s) shown at the right of the label row
- `options.rule` — default `true`; `options.density` — default `"flush"`

```js
ui.section("Media library", { level: 1, anchor: "media", caption: "Files under course/media." },
  cms.recipes.mediaLibrary({ assets, onFilesSelected: "admin-upload-course-material" }))
```

Example:

```js
const ui = require("ui.dsl")

ui.page({
  id: "overview",
  title: "Overview",
  sections: [
    ui.panel({ title: "Status" },
      ui.inline({ gap: "sm" },
        ui.statusText({ status: "ready", icon: true }),
        ui.caption({ tone: "muted" }, "Rendered by React")
      )
    )
  ]
})
```

### UI recipes

- `ui.recipes.metrics({ items, recipe? })`
- `ui.recipes.actionToolbar({ title?, actions, caption?, gap?, wrap? })`

## `data.dsl` helpers

`data.dsl` exports:

- `dataTable(props)`
- `cell.field(field, options?)`
- `cell.number(field, options?)`
- `cell.status(field, options?)`
- `cell.caption(field, options?)`
- `cell.template(template)`
- `cell.link(hrefField, labelField, options?)`
- `cell.linkButton(hrefField, labelField, options?)`
- `cell.actionButton(label, action, options?)`
- `cell.constant(value)`

Example:

```js
const ui = require("ui.dsl")
const data = require("data.dsl")

ui.panel({ title: "Rows" },
  data.dataTable({
    rows: [{ id: "a", name: "Alpha", status: "running" }],
    getRowKey: "id",
    columns: [
      { id: "name", header: "Name", cell: data.cell.field("name") },
      { id: "status", header: "Status", cell: data.cell.status("status", { icon: true }) }
    ]
  })
)
```

### The v2 data grammar: `data.v2.dsl`

Use this form for new examples and new page code. It is the hard-cutover typed/fluent direction; it replaces v1 option bags with builder handles and validation terminals.

```js
const data = require("data.v2.dsl")

const sessionSchema = data.schema("Session")
  .field("sessionId", data.f.key().label("ID").width("14ch"))
  .field("title", data.f.primary().label("Title").required().maxLength(120))
  .field("turnCount", data.f.count().label("Turns"))
  .field("status", data.f.status().label("Status"))
  .field("body", data.f.prose().label("Body").rows(3))
  .build()

data.collection("sessions", sessions)
  .schema(sessionSchema)
  .table()
  .toIR()
```

Selectable URL-backed table:

```js
data.collection("sessions", sessions)
  .schema(sessionSchema)
  .select(s => s.urlParam("selected", query.selected))
  .table(t => t.rowSelect(data.action.navigate("/pages/sessions?selected=${row.sessionId}")))
  .toIR()
```

Master-detail editor with native form submit:

```js
data.collection("agenda", agenda)
  .schema(agendaSchema)
  .edit(e => e
    .selectUrl("agenda", query.agenda)
    .submitPost("/settings/agenda-item")
    .create({ label: "New agenda item" }))
  .masterDetail()
  .toIR()
```

Row/server actions:

```js
data.collection("agenda", agenda)
  .schema(agendaSchema)
  .edit(e => e.actions(a => a
    .reorder(data.action.server("admin-reorder-course-agenda"))
    .remove(data.action.server("admin-delete-agenda-item")
      .confirm("Delete ${row.title}?"))))
  .masterDetail()
  .toIR()
```

Live demo pages in `go-go-course`:

- `/pages/dsl-examples-table`
- `/pages/dsl-examples-selectable-table`
- `/pages/dsl-examples-master-detail`
- `/pages/dsl-examples-actions`

### Legacy data grammar: `data.dsl` `schema`, `f`, `record`, `collection`

The following v1 option-bag grammar remains documented as legacy/current runtime behavior for existing pages such as the admin Course CMS. Do **not** use these examples for new v2 demos or new hard-cutover page code; prefer `data.v2.dsl` above.

Intent-level authoring: declare what the records look like and how they should be shown or edited; the grammar compiles to `DataTable`/`FormPanel`/`FieldGrid`/`SectionBlock` IR. See ticket RAGEVAL-UI-GRAMMAR design-doc 02.

**Field roles** (`data.f.*`): `key`, `primary`, `short`, `prose`, `count`, `size`, `measure`, `date`, `status`, `tags`, `media`, `href`. A role decides the summary cell (prose/media are elided from tables; status renders as StatusText; count/size/measure as numbers), the editor control (prose → stacked textarea; key → read-only text input), and grid batching. Options per field: `label`, `width`, `placeholder`, `required`, `maxLength`, `rows`, `hint`, `readOnly`, `editable`.

**`data.schema(fields)`** — ordered field specs:

```js
const agendaSchema = data.schema({
  id:          data.f.key({ hint: "Stable anchor. Leave blank for a generated ID." }),
  number:      data.f.short({ label: "Time", width: "6ch", placeholder: "14h30" }),
  duration:    data.f.short({ width: "8ch" }),
  title:       data.f.primary({ required: true, maxLength: 160 }),
  description: data.f.prose({ rows: 4, maxLength: 800 }),
})
```

**`data.record(values, options)`** — one record. `verb: "edit"` (default) compiles to a `FormPanel` whose rows derive from the schema — consecutive short fields batch into a `FieldGrid`, prose fields become stacked textareas; `verb: "show"` compiles to a `MetadataGrid`. `submit: data.formPost("/settings/…")` wires the native form post; `title`, `status`, `statusMessage`, `submitLabel`, `footer` pass through.

**`data.collection(rows, options)`** — records through an arrangement:

- `verb`: `"show" | "edit" | "pick" | "manage"`
- `arrange`: `"table" | "master-detail"`
- `select: data.urlParam("agenda", query.agenda)` — selection state lives in the URL; row clicks navigate `?agenda=<key>`; the value `"__new"` opens an empty editor
- `submit: data.formPost(...)` — per-record save for the detail editor
- `open` — action for an Open column / row activation when there is no `select`
- `reorder` — action dispatched with `payload.direction: "up" | "down"` and the row in context
- `remove` — action (with `confirm: "Delete ${row.title}?"`) for a Delete column
- `create: true` — a "New item" button navigating `?<param>=__new`
- `title`/`caption` — wraps everything in a flat `SectionBlock` (level 2, ruled); `empty`, `getRowKey` as in `dataTable`

```js
ui.section("Agenda", { level: 1 },
  data.collection(agenda, {
    schema: agendaSchema,
    title: "Workshop agenda",
    verb: "edit",
    arrange: "master-detail",
    select: data.urlParam("agenda", query.agenda),
    submit: data.formPost("/settings/agenda-item"),
    reorder: "admin-reorder-course-agenda",
    remove: { kind: "server", name: "admin-delete-agenda-item", confirm: "Delete ${row.title}?" },
    create: true,
    empty: "No agenda items yet.",
  }))
```

The summary table elides prose, keys render muted, and the detail editor for the selected row derives from the same schema — roughly 250 px of table plus one editor, replacing an unrolled stack of per-item panels.

### Data recipes

- `data.recipes.masterDetailTable({ rows, columns, getRowKey?, selectedKey?, onRowSelect?, detail?, detailTitle? })`

Use `ui.dsl` helpers inside `detail` callbacks when you need UI nodes:

```js
data.recipes.masterDetailTable({
  rows,
  columns,
  detail: row => ui.panel({ title: "Selected" }, row ? row.name : "No selection")
})
```

## `context_window.dsl` helpers

`context_window.dsl` exports context-window diagrams, transcript surfaces, annotations, anchored comments, upload widgets, and the style-set helpers that make those widgets palette-aware.

Component helpers:

- `contextStyleSwatch`, `annotationBadge`, `transcriptRoleBadge`
- `contextLegend`, `contextBudgetBar`, `contextStripDiagram`, `contextGroupedStripDiagram`, `contextStackDiagram`, `contextTreemap`, `contextDiagramPanel`
- `transcriptSessionHeader`, `transcriptMessageCard`, `annotationNoteCard`, `annotationRailPanel`, `transcriptReaderPanel`, `transcriptWorkspacePanel`
- `anchoredCommentCard`, `anchoredCommentRail`
- `contextUploadDropArea`

Style-set helpers:

- `visualStyle(options)` returns a JSON-compatible `ContextVisualStyle` object. Useful fields are `fill`, `line`, `stroke`, `labelColor`, `pattern`, `strokeWidth`, and `vars`.
- `legendItem(id, label, options?)` returns a `ContextLegendItemSpec`; pass `hidden: true` for style keys that should render but not appear in the visible legend.
- `styleSet(options)` returns a `ContextStyleSet` and normalizes missing `legend` / `styles` to empty containers.
- `paletteStyleSet(options)` builds a palette-derived `ContextStyleSet` from `palette` plus `entries`.
- `contextPart(id, label, styleKey, tokens, options?)` creates a context-window part with required `styleKey`.
- `contextSnapshot(options)` creates a normalized context-window snapshot.

The hard-cutover contract is `styleKey + styleSet`. Context-window parts do not use `kind`, and context diagram widgets must receive a `styleSet` explicitly or through `recipes.contextDiagram({ palette, entries, ... })`.

Example with a caller-defined legend:

```js
const contextWindow = require("context_window.dsl")

const styleSet = contextWindow.styleSet({
  legend: [
    contextWindow.legendItem("prompt", "Prompt"),
    contextWindow.legendItem("evidence", "Evidence"),
    contextWindow.legendItem("answer", "Answer"),
    contextWindow.legendItem("free", "Headroom", { hidden: true }),
  ],
  styles: {
    prompt: contextWindow.visualStyle({ pattern: "checker", fill: "#f2eee8", line: "#5f7f89", stroke: "#111314", labelColor: "#111314" }),
    evidence: contextWindow.visualStyle({ pattern: "stipple", fill: "#f6e6df", line: "#c46a55", stroke: "#111314", labelColor: "#111314" }),
    answer: contextWindow.visualStyle({ pattern: "solid", fill: "#5f7f89", stroke: "#111314", labelColor: "#ffffff" }),
    free: contextWindow.visualStyle({ pattern: "none", fill: "#f2f3ef", stroke: "#b8bdbb", labelColor: "#111314" }),
  },
})

const snapshot = contextWindow.contextSnapshot({
  id: "rag-window",
  title: "RAG answer window",
  limit: 32000,
  parts: [
    contextWindow.contextPart("p", "Prompt", "prompt", 1400),
    contextWindow.contextPart("e", "Evidence", "evidence", 9200),
    contextWindow.contextPart("a", "Draft", "answer", 1800),
    contextWindow.contextPart("h", "Headroom", "free", 19600),
  ],
})

contextWindow.contextDiagramPanel({
  snapshot,
  styleSet,
  initialView: "budget",
  selectedPartId: "e",
})
```

Example with a preferred palette:

```js
const styleSet = contextWindow.paletteStyleSet({
  palette: "Signal Orange / Cyan",
  entries: [
    { id: "guardrails", label: "Guardrails", accent: "b", pattern: "checker", fillPct: 16, linePct: 70 },
    { id: "chat", label: "Chat", accent: "grid", pattern: "none" },
    { id: "commands", label: "Commands", accent: "a", pattern: "stipple", fillPct: 16, linePct: 60 },
    { id: "free", label: "Free", accent: "grid", pattern: "none", hidden: true },
  ],
})
```

Transcript widgets also accept `styleSet` so role title bars, note chips, and side-note headers share the same palette. Transcript bodies remain neutral by design; palette colors should live in chrome and small affordances with explicit `labelColor` foregrounds.

Action contexts:

- annotation selection: `{ annotationId, value, componentType }`
- anchored comment selection: `{ commentId, value, componentType }`
- upload selection: `{ files, fileNames, fileCount, componentType }`

### Context-window recipes

- `contextWindow.recipes.contextDiagram({ snapshot, styleSet, view?, initialView?, selectedPartId? })`
- `contextWindow.recipes.contextDiagram({ snapshot, palette, entries, view?, initialView?, selectedPartId? })`
- `contextWindow.recipes.annotatedTranscript({ transcript?, title?, subtitle?, messages?, annotations?, selectedAnnotationId?, showNotes?, styleSet?, onAnnotationSelect? })`

## `course.dsl` helpers

`course.dsl` exports:

- `contextStudioNavIcon`
- `courseStepNav`, `markdownArticle`, `documentListPanel`, `documentPreviewToolbar`
- `courseLessonPanel`, `courseSlidePanel`, `courseStudioShell`, `handoutDocumentShell`
- `slideShell`

Example:

```js
const course = require("course.dsl")

course.courseStudioShell({
  sections: [{ id: "course", label: "Course", items: [
    { id: "slides", label: "Slides", icon: course.contextStudioNavIcon({ id: "slides" }) }
  ] }],
  activeItemId: "slides",
  title: "Studio"
},
  course.courseSlidePanel({ slide, snapshot, index: 0, total: slides.length })
)
```

### Course recipes

- `course.recipes.courseStudio({ sections, activeItemId?, title?, subtitle?, main?, onNavigate? })`
- `course.recipes.courseSlide({ slide, snapshot, index?, total?, visualSide?, onPrevious?, onNext? })`
- `course.recipes.handout({ bundle?, intro?, documents?, selectedDocumentId?, title?, onSelect?, onDownload?, onDownloadAll? })`

## `cms.dsl` helpers

`cms.dsl` exports content-management widgets built for admin/CMS pages:

- `mediaThumb`, `tag`, `contentStatusBadge`, `meterBar`
- `tileGrid`, `assetTile`, `breadcrumbs`, `pagination`, `searchField`, `emptyState`
- `markdownEditor` — browser-local editing state with a live Markdown preview; give it `name` + `defaultValue` so the value participates in `formPanel({ method: "post", formAction })` form posts
- `mediaLibraryPanel`, `articleListPanel`, `cmsShell`

Assets are `CmsAsset`-shaped objects (`{ id, kind: "image"|"file", title, filename, mime, size, src, tags, status, createdAt, updatedAt }`); articles are `CmsArticleSummary`-shaped (`{ id, slug, title, status, tags, updatedAt, author?, excerpt? }`). Content status is `draft | published | scheduled | archived`.

`mediaLibraryPanel`'s `onFilesSelectedAction` receives serialized files in the action context (`{ files, fileNames, fileCount }` with utf8 text or base64 payloads) — the same contract as `contextUploadDropArea`. Selection, paging, and filtering pair naturally with query-parameter navigation:

```js
const cms = require("cms.dsl")

cms.recipes.mediaLibrary({
  assets,
  selectedAssetIds: query.asset ? [query.asset] : [],
  onAssetSelect: ui.action.navigate("?asset=$assetId"),
  onPageChange: ui.action.navigate("?page=$page"),
  onFilesSelected: "admin-upload-course-material",
})
```

### CMS recipes

- `cms.recipes.mediaLibrary({ assets, selectedAssetIds?, selectionMode?, query?, kindFilter?, page?, pageCount?, uploads?, showStatusBadges?, emptyMessage?, title?, minTileWidth?, onAssetSelect?, onAssetOpen?, onQuerySubmit?, onKindFilterChange?, onPageChange?, onFilesSelected? })`
- `cms.recipes.articleList({ articles, selectedArticleId?, statusFilter?, query?, page?, pageCount?, emptyMessage?, title?, maxVisibleTags?, onArticleSelect?, onCreate?, onRowAction?, onStatusFilterChange?, onQuerySubmit?, onPageChange? })`

`onRowAction` dispatches with `{ articleId, rowAction }` where `rowAction` is `edit | publish | archive | delete`; archive/delete confirm inside the panel before dispatching.

## JSON compatibility rules

- Widgets are plain objects: `{ kind: "component", type, props?, children? }`.
- Children can be strings, numbers, widget nodes, arrays, or fragments.
- Renderable fields accept strings, numbers, or widget nodes.
- Table cells must use `data.cell.*` specs; JavaScript render functions cannot cross the JSON boundary.
- Use action specs instead of callback functions for renderer events.
