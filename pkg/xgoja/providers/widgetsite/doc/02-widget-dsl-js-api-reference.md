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
- `formRow`, `selectInput`, `textInput`, `tabList`
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

`context_window.dsl` exports:

- `contextKindSwatch`, `annotationBadge`, `transcriptRoleBadge`
- `contextLegend`, `contextBudgetBar`, `contextStripDiagram`, `contextStackDiagram`, `contextTreemap`, `contextDiagramPanel`
- `transcriptSessionHeader`, `transcriptMessageCard`, `annotationNoteCard`, `annotationRailPanel`, `transcriptReaderPanel`, `transcriptWorkspacePanel`
- `anchoredCommentCard`, `anchoredCommentRail`
- `contextUploadDropArea`

Example:

```js
const contextWindow = require("context_window.dsl")

contextWindow.contextDiagramPanel({
  snapshot,
  initialView: "budget",
  selectedPartId: "retrieval"
})
```

Action contexts:

- annotation selection: `{ annotationId, value, componentType }`
- anchored comment selection: `{ commentId, value, componentType }`
- upload selection: `{ files, fileNames, fileCount, componentType }`

### Context-window recipes

- `contextWindow.recipes.contextDiagram({ snapshot, view?, initialView?, selectedPartId? })`
- `contextWindow.recipes.annotatedTranscript({ transcript?, title?, subtitle?, messages?, annotations?, selectedAnnotationId?, showNotes?, onAnnotationSelect? })`

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
