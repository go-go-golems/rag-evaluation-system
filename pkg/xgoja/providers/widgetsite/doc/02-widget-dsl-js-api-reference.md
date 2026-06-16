---
Title: "Widget DSL JavaScript API Reference"
Slug: widget-dsl-js-api-reference
Short: "Reference for split Widget DSL modules, helpers, recipes, actions, and table cell specifications."
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

## Module names

| Module | Purpose |
| --- | --- |
| `ui.dsl` | Generic page, layout, primitive, foundation, and UI recipe helpers. |
| `data.dsl` | Data-display widgets, `cell.*` helpers, and data recipes. |
| `context_window.dsl` | Context-window diagrams, transcript, annotation, anchored-comment, and upload helpers. |
| `course.dsl` | Course, lesson, slide, handout, and course-studio helpers. |

There is no compatibility bucket module. Import the domain modules explicitly:

```js
const ui = require("ui.dsl")
const data = require("data.dsl")
const contextWindow = require("context_window.dsl")
const course = require("course.dsl")
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
ui.action.event("widget:selected", { detail: { source: "story" } })
ui.action.copy("copy this")
```

## `ui.dsl` helpers

`ui.dsl` exports generic visual helpers:

- `appShell`, `appNav`
- `button`, `caption`, `codeText`, `divider`, `statusText`, `textBlock`
- `inline`, `stack`, `dashboardGrid`, `panel`, `scrollRegion`, `sectionBlock`, `sidebarShell`, `splitPane`
- `formRow`, `selectInput`, `textInput`, `textareaInput`, `tabList`
- `metadataGrid`, `keyValueStrip`, `checkList`, `stepList`, `personSummary`, `figureBlock`, `keyPointList`, `sidebarNav`

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
- `contextLegend`, `contextBudgetBar`, `contextStripDiagram`, `contextStackDiagram`, `contextTreemap`, `contextDiagramPanel`
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

## JSON compatibility rules

- Widgets are plain objects: `{ kind: "component", type, props?, children? }`.
- Children can be strings, numbers, widget nodes, arrays, or fragments.
- Renderable fields accept strings, numbers, or widget nodes.
- Table cells must use `data.cell.*` specs; JavaScript render functions cannot cross the JSON boundary.
- Use action specs instead of callback functions for renderer events.
