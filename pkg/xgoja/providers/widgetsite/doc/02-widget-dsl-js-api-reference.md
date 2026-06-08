---
Title: "Widget DSL JavaScript API Reference"
Slug: widget-dsl-js-api-reference
Short: "Reference for widget.dsl and rag.dsl constructors, component helpers, children, and table cell specifications."
Topics:
- xgoja
- widget-dsl
- rag-widget-site
- widget-ir
- api-reference
Commands:
- xgoja build
- xgoja list-modules
Flags: []
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

This reference documents the JavaScript API exposed by `require("widget.dsl")` and `require("rag.dsl")`. Both module names return the same API. The helpers create Widget IR, a JSON-compatible representation consumed by the React WidgetRenderer.

The API is intentionally small. It gives scripts enough structure to describe pages while keeping rendering decisions in React. Every function returns a plain value: an object, an array, a string-normalized text node, or a cell specification. If a value cannot be serialized to JSON, it should not be part of Widget IR.

## Module names

| Module | Purpose |
| --- | --- |
| `widget.dsl` | Canonical Widget IR authoring module. |
| `rag.dsl` | RAG-oriented alias with the same exports. |

Use the canonical name unless a script is written specifically for RAG examples:

```js
const widget = require("widget.dsl")
```

## Core node constructors

The core constructors build the three basic Widget IR node forms: text nodes, element nodes, and component nodes.

### `text(value)`

`text(value)` converts a JavaScript value into a Widget IR text node.

```js
widget.text("hello")
```

Result:

```json
{ "kind": "text", "text": "hello" }
```

Most code does not need to call `text` directly. When a string or number is passed as a child to a component helper, the module converts it to a text node automatically.

### `element(tag, attrs?, ...children)`

`element` creates a generic element node. Use it sparingly. Prefer component helpers when the desired UI maps to a real RAG component.

```js
widget.element("span", { className: "note" }, "small note")
```

Result:

```json
{
  "kind": "element",
  "tag": "span",
  "attrs": { "className": "note" },
  "children": [{ "kind": "text", "text": "small note" }]
}
```

The second argument is treated as attributes only when it is a plain object and not already a Widget IR node. Children may be nodes, arrays of nodes, strings, numbers, booleans, or nullish values. Nullish children are ignored.

### `component(type, props?, ...children)`

`component` creates a component node by explicit component type. It is useful for advanced cases or new renderer components that do not yet have a named helper.

```js
widget.component("Panel", { title: "Manual" }, "Body")
```

Result:

```json
{
  "kind": "component",
  "type": "Panel",
  "props": { "title": "Manual" },
  "children": [{ "kind": "text", "text": "Body" }]
}
```

### `fragment(...children)`

`fragment` flattens children into an array. It is useful when a helper wants to return multiple sibling nodes.

```js
const controls = widget.fragment(
  widget.button({}, "Run"),
  widget.button({ variant: "primary" }, "Cancel")
)
```

`fragment` returns an array, not a node. Passing a fragment as a child flattens it into the parent children array.

## Component helpers

Component helpers are convenience functions for component nodes. Each helper has the same call shape:

```js
widget.helperName(props?, ...children)
```

The first argument is props when it is a plain object. Otherwise, it is treated as the first child.

| Helper | Component type | Typical use |
| --- | --- | --- |
| `appShell` | `AppShell` | Top-level application layout. |
| `appNav` | `AppNav` | Navigation structure. |
| `button` | `Button` | Action button. |
| `caption` | `Caption` | Secondary text, notes, and muted labels. |
| `dashboardGrid` | `DashboardGrid` | Dashboard card grid layout. |
| `dataTable` | `DataTable` | Tabular data with serializable columns. |
| `formRow` | `FormRow` | Label/control form row. |
| `inline` | `Inline` | Inline layout for compact controls. |
| `metadataGrid` | `MetadataGrid` | Key-value metadata display. |
| `panel` | `Panel` | Section container with optional title. |
| `scrollRegion` | `ScrollRegion` | Scrollable content region. |
| `sectionBlock` | `SectionBlock` | Editorial section with optional label/caption/divider. |
| `selectInput` | `SelectInput` | Select control. |
| `sidebarShell` | `SidebarShell` | Two-column shell with sidebar, optional header/footer, and optional content padding. |
| `slideShell` | `SlideShell` | Generic teaching/presentation slide layout with slots. |
| `splitPane` | `SplitPane` | Two-pane layout with ratios, divider, and optional gutters. |
| `stack` | `Stack` | Vertical stack layout. |
| `statusText` | `StatusText` | Status label with optional icon. |
| `tabList` | `TabList` | Tab selector. |
| `textBlock` | `Text` | Foundation text component; use when typography props are needed. |
| `textInput` | `TextInput` | Text input control. |
| `contextKindSwatch` | `ContextKindSwatch` | Visual marker for a context part kind. |
| `annotationBadge` | `AnnotationBadge` | Compact semantic badge. |
| `transcriptRoleBadge` | `TranscriptRoleBadge` | Badge for transcript message roles. |
| `contextLegend` | `ContextLegend` | Legend for context diagram colors/kinds. |
| `contextBudgetBar` | `ContextBudgetBar` | Budget bar for a context-window snapshot. |
| `contextStripDiagram` | `ContextStripDiagram` | Strip view for context-window parts. |
| `contextStackDiagram` | `ContextStackDiagram` | Stacked view for context-window parts. |
| `contextTreemap` | `ContextTreemap` | Treemap view for context-window parts. |
| `contextDiagramPanel` | `ContextDiagramPanel` | Switchable context diagram panel. |
| `transcriptSessionHeader` | `TranscriptSessionHeader` | Title/status header for transcript sessions. |
| `transcriptMessageCard` | `TranscriptMessageCard` | Single transcript message card. |
| `annotationNoteCard` | `AnnotationNoteCard` | Single transcript annotation note. |
| `annotationRailPanel` | `AnnotationRailPanel` | Annotation list/rail for transcripts. |
| `transcriptReaderPanel` | `TranscriptReaderPanel` | Transcript reader with optional annotation chips. |
| `transcriptWorkspacePanel` | `TranscriptWorkspacePanel` | Transcript reader plus notes workspace. |
| `anchoredCommentCard` | `AnchoredCommentCard` | Single spatially anchored comment. |
| `anchoredCommentRail` | `AnchoredCommentRail` | List/rail of anchored comments. |
| `keyValueStrip` | `KeyValueStrip` | Horizontal key/value summary strip. |
| `checkList` | `CheckList` | Checklist content block. |
| `stepList` | `StepList` | Ordered/agenda-style selectable step list. |
| `personSummary` | `PersonSummary` | Instructor/person summary block. |
| `figureBlock` | `FigureBlock` | Figure wrapper with caption/legend slots. |
| `keyPointList` | `KeyPointList` | Teaching key-points list. |
| `sidebarNav` | `SidebarNav` | Grouped sidebar navigation. |
| `courseStepNav` | `CourseStepNav` | Course agenda navigation. |
| `markdownArticle` | `MarkdownArticle` | Dependency-free markdown subset renderer. |
| `documentListPanel` | `DocumentListPanel` | Handout/document list panel. |
| `documentPreviewToolbar` | `DocumentPreviewToolbar` | File/format/size toolbar for document preview. |
| `courseLessonPanel` | `CourseLessonPanel` | Course landing/lesson panel. |
| `courseSlidePanel` | `CourseSlidePanel` | Course slide with context visual. |
| `courseStudioShell` | `CourseStudioShell` | Sidebar course studio shell. |
| `handoutDocumentShell` | `HandoutDocumentShell` | Full handout document browser/reader. |

Example:

```js
widget.panel({ title: "Evaluation" },
  widget.stack({ gap: "md" },
    widget.statusText({ status: "running", icon: true }, "Running"),
    widget.caption({ tone: "muted" }, "The page is rendered by React.")
  )
)
```

The component helper does not validate every prop. It preserves props as JSON-compatible data and lets the renderer decide how to interpret them. This keeps the authoring layer lightweight but means authors should use props supported by the React component library.

## Semantic component props and DTO shapes

The newer context-viewer helpers use JSON DTOs that match the React package types. The most common shapes are listed below. All renderable fields accept strings, numbers, or Widget IR nodes unless noted.

### Context-window DTOs

```js
const snapshot = {
  id: "ctx",
  title: "Context Window",
  subtitle: "Optional subtitle",
  limit: 16000,
  selectedPartId: "retrieval",
  parts: [
    { id: "system", label: "System prompt", kind: "system", tokens: 620, note: "Pinned" },
    { id: "retrieval", label: "Retrieved docs", kind: "retrieval", tokens: 7100 }
  ]
}
```

Valid context kinds include `system`, `instruction`, `context`, `conversation`, `summary`, `retrieval`, `tool`, `result`, `generated`, `annotation`, `course`, `active`, `evicted`, `empty`, and `other`.

```js
widget.contextDiagramPanel({ snapshot, initialView: "budget", selectedPartId: "retrieval" })
widget.contextBudgetBar({ snapshot })
widget.contextStripDiagram({ snapshot, selectedPartId: "system" })
widget.contextStackDiagram({ snapshot })
widget.contextTreemap({ snapshot })
widget.contextLegend({ compact: true })
```

Diagram views are `strip`, `stack`, `budget`, and `treemap`.

### Transcript and annotation DTOs

```js
const transcript = {
  title: "Context review session",
  subtitle: "Optional subtitle",
  messages: [
    { id: "m1", role: "user", text: "Question", tokens: 12, timestamp: "09:14", annotationIds: ["a1"] },
    { id: "m2", role: "assistant", text: "Answer", tokens: 34, timestamp: "09:15" }
  ],
  annotations: [
    { id: "a1", targetMessageId: "m1", kind: "context", label: "Budget pressure", text: "Retrieved docs dominate.", confidence: 0.92 }
  ],
  selectedAnnotationId: "a1"
}
```

Roles are `system`, `developer`, `user`, `assistant`, `tool`, and `other`.

```js
widget.transcriptWorkspacePanel({
  title: transcript.title,
  messages: transcript.messages,
  annotations: transcript.annotations,
  selectedAnnotationId: "a1",
  showNotes: true,
  onAnnotationSelectAction: widget.action.server("note-selected")
})

widget.splitPane({
  ratio: "course",
  divider: true,
  gutter: "lg",
  left: widget.transcriptReaderPanel({ messages: transcript.messages, annotations: transcript.annotations, showAnnotationChips: true }),
  right: widget.annotationRailPanel({ annotations: transcript.annotations, selectedAnnotationId: "a1" })
})
```

`onAnnotationSelectAction` receives context with `annotationId`, `value`, and `componentType`.

### Anchored comments

```js
const comments = [
  { id: "c1", anchorX: 28, anchorY: 18, author: "Reviewer", text: "Clarify this region.", time: "09:22", status: "open" }
]

widget.anchoredCommentRail({
  title: "Review comments",
  comments,
  selectedCommentId: "c1",
  onCommentSelectAction: widget.action.server("comment-selected")
})
```

Comment statuses are `open` and `resolved`. `onCommentSelectAction` receives `commentId`, `value`, and `componentType`.

### Course and slide DTOs

```js
const course = {
  id: "context-workshop",
  kicker: "Workshop",
  title: "Context Window Engineering",
  tagline: "Learn to see and budget context windows.",
  instructor: { name: "RAG Systems Team", role: "Instructor", bio: "Builds evaluation tooling." },
  outcomes: ["Read diagrams", "Annotate transcripts"],
  agenda: [
    { id: "map", number: "01", title: "Map the window", description: "Identify context slices.", duration: "20 min" }
  ]
}

const slide = {
  id: "slide-budget",
  number: "01",
  kicker: "Budget",
  title: "Context budget is a product decision",
  view: "budget",
  snapshotId: snapshot.id,
  notes: ["Pinned content protects invariants."]
}
```

```js
widget.courseLessonPanel({ course, activeAgendaItemId: "map", onAgendaItemSelectAction: widget.action.server("agenda-selected") })
widget.courseSlidePanel({ slide, snapshot, index: 0, total: 3, visualSide: "right" })
widget.courseStudioShell({ sections, activeItemId: "slides" },
  widget.courseSlidePanel({ slide, snapshot, index: 0, total: 3 })
)
```

For custom slides, use `slideShell` with slot props:

```js
widget.slideShell({
  eyebrow: "custom slide",
  counter: "02 / 03",
  title: "Direct Widget IR composition",
  primary: widget.figureBlock({ frame: "bordered", caption: "Budget" }, widget.contextBudgetBar({ snapshot })),
  secondary: widget.keyPointList({ items: [
    { id: "json", title: "JSON", text: "xgoja returns Widget IR data." }
  ]}),
  footer: widget.inline({ justify: "between" }, widget.button({}, "Previous"), widget.button({ variant: "primary" }, "Next"))
})
```

### Handout/document DTOs

```js
const bundle = {
  intro: "Reference material for the workshop.",
  docs: [
    { id: "guide", title: "Guide", file: "guide.md", format: "markdown", size: "8 KB", description: "Checklist", body: "# Guide\n\nContent." }
  ]
}
```

```js
widget.handoutDocumentShell({
  title: "Workshop handout",
  intro: bundle.intro,
  documents: bundle.docs,
  selectedDocumentId: "guide",
  onDocumentSelectAction: widget.action.server("document-selected"),
  onDownloadAction: widget.action.server("document-download"),
  onDownloadAllAction: widget.action.server("download-all")
})

widget.splitPane({
  ratio: "sidebar",
  divider: true,
  gutter: "lg",
  left: widget.documentListPanel({
    title: "Documents",
    items: bundle.docs.map(doc => ({ id: doc.id, title: doc.title, format: doc.format, size: doc.size, description: doc.description })),
    selectedItemId: "guide",
    showItemDescriptions: true
  }),
  right: widget.stack({ gap: "sm" },
    widget.documentPreviewToolbar({ file: "guide.md", format: "markdown", size: "8 KB" }),
    widget.markdownArticle({ source: bundle.docs[0].body }))
})
```

`MarkdownArticle` supports the dependency-free subset used by the examples: headings, paragraphs, lists, code-ish text, and simple tables.

### Layout details for new examples

```js
widget.splitPane({ ratio: "sidebar", divider: true, gutter: "lg", left, right })
widget.sidebarShell({ sidebarWidth: 204, contentPadding: "lg", sidebar, header, footer }, mainContent)
widget.sectionBlock({ label: "Overview", density: "spacious", divider: "bottom" }, body)
```

`splitPane.ratio` may be `balanced`, `leftNarrow`, `rightNarrow`, `course`, or `sidebar`. `gutter` may be `none`, `md`, or `lg`.

`sidebarShell.contentPadding` may be `none`, `md`, or `lg`.

## Status values

`statusText` and status table cells display a label and, optionally, an icon. The renderer has icons for known status values:

| Status | Meaning in examples |
| --- | --- |
| `pending` | Work has not started. |
| `ready` | Work can start. |
| `running` | Work is in progress. |
| `succeeded` | Work completed successfully. |
| `done` | Work is complete. |
| `partial` | Work completed with partial results. |
| `warning` | Work completed with a warning. |
| `failed` | Work failed. |
| `error` | Work produced an error. |
| `canceled` | Work was canceled. |

Use `succeeded`, not `success`, when you want the success checkmark icon. Unknown statuses still render their label, but the icon fallback is `?`.

## Data tables and cell specifications

`dataTable` expects serializable `rows` and serializable `columns`. A column's `cell` field is a cell specification created by `widget.cell.*` helpers.

```js
widget.dataTable({
  rows: [
    { id: 1, name: "Fast Growing Trees", status: "succeeded" },
    { id: 2, name: "Arborvitae Spacing", status: "running" }
  ],
  getRowKey: "id",
  columns: [
    { id: "id", header: "ID", cell: widget.cell.field("id") },
    { id: "name", header: "Name", cell: widget.cell.field("name") },
    { id: "status", header: "Status", cell: widget.cell.status("status") }
  ]
})
```

The JSON boundary is the reason cell helpers exist. A React table column can use functions, closures, and component renderers. Widget IR cannot. A cell specification states what the renderer should do using plain data.

### `cell.field(field, options?)`

Displays the value from `row[field]`.

```js
widget.cell.field("name")
```

### `cell.number(field, options?)`

Displays a numeric value. Options may include formatting hints understood by the renderer.

```js
widget.cell.number("score", { maximumFractionDigits: 2 })
```

### `cell.status(field, options?)`

Displays `row[field]` through `StatusText`.

```js
widget.cell.status("status", { icon: true })
```

### `cell.caption(field, options?)`

Displays `row[field]` as caption-style secondary text.

```js
widget.cell.caption("description", { tone: "muted" })
```

### `cell.template(template)`

Displays a template string evaluated by the renderer's table-cell implementation. Use this only when field cells are not expressive enough.

```js
widget.cell.template("${name} (${status})")
```

### `cell.link(hrefField, labelField, options?)`

Displays a link where the href and label come from row fields.

```js
widget.cell.link("url", "title", { target: "_blank" })
```

### `cell.constant(value)`

Displays a constant renderable value for every row.

```js
widget.cell.constant(widget.caption({ tone: "muted" }, "fixed"))
```

The value is normalized like a child. A string becomes a text node; a Widget IR node is preserved; an array is flattened.

## Child normalization

Every component helper accepts children. The module normalizes children with these rules:

| Input | Result |
| --- | --- |
| `null` or `undefined` | Ignored. |
| Widget IR node | Preserved. |
| Array-like value | Flattened recursively. |
| String, number, boolean, other value | Converted to `{ kind: "text", text: String(value) }`. |

This means conditional rendering can be written naturally:

```js
widget.panel({ title: "Run" },
  hasWarning ? widget.statusText({ status: "warning", icon: true }, warning) : null,
  rows.length === 0 ? "No rows yet" : widget.dataTable({ rows, columns })
)
```

The null branch is ignored, so no empty placeholder node appears in the final IR.

## Actions

Buttons and other interactive widgets may carry action specifications in props. A server action is serializable:

```js
widget.button({
  variant: "primary",
  action: { kind: "server", name: "refresh", payload: { source: "demo" } }
}, "Refresh")
```

The React app posts server actions to `/api/widget/actions/{name}`. The exact action contract is owned by the WidgetRenderer server or jsverb site. The DSL only preserves the action object.

The `action` namespace provides small constructors for common action specifications:

```js
widget.action.server("add-query", { payload: { owner: "research" } })
widget.action.navigate("/pages/actions")
widget.action.event("widget:selected", { detail: { id: 1 } })
widget.action.copy("copied value")
```

The constructors return plain JSON-compatible action objects. Low-level literal action objects remain valid.

## Page helper and semantic recipes

The `page` helper builds the standard page wrapper and can turn a list of sections into a stacked root:

```js
return widget.page({
  id: "actions",
  title: "Actions",
  meta: { shell: "app", maxWidth: "wide" },
  sections: [
    widget.panel({ title: "Intro" }, "Rendered by React")
  ]
})
```

`page({ root })` preserves an explicit root. `page({ sections })` creates a `Stack` root with `gap: "lg"` by default.

The `recipes` namespace expands common RAG dashboard patterns into ordinary Widget IR:

### `recipes.metrics({ items, recipe? })`

Creates a dashboard grid of condensed metric panels:

```js
widget.recipes.metrics({ items: [
  { label: "Total", value: 12, status: "ready" },
  { label: "Running", value: 3, status: "running" }
]})
```

### `recipes.actionToolbar({ title, actions, caption? })`

Creates a panel containing inline buttons. `action` may be a full action spec or a string shorthand for a server action name:

```js
widget.recipes.actionToolbar({
  title: "Queue controls",
  actions: [
    { label: "Add query", variant: "primary", action: "add-query", payload: { owner: "research" } },
    { label: "Reset", action: widget.action.server("reset-demo") }
  ],
  caption: "Actions refresh the React app after the server mutates state."
})
```

### `recipes.masterDetailTable(options)`

Creates a two-up dashboard grid containing a `DataTable` panel and a detail panel. `onRowSelect` may be a full action spec or a string shorthand for a server action name. `detail` may be a callback that receives the selected row and returns Widget IR:

```js
widget.recipes.masterDetailTable({
  title: "Query queue",
  rows,
  columns,
  selectedKey,
  onRowSelect: "select-query",
  detail: row => widget.panel({ title: "Selected" }, row ? row.name : "No selection")
})
```

### `recipes.contextDiagram({ snapshot, view?, initialView?, selectedPartId? })`

Creates a `ContextDiagramPanel` from a serializable context-window snapshot:

```js
widget.recipes.contextDiagram({ snapshot, view: "budget", selectedPartId: "retrieval" })
```

### `recipes.annotatedTranscript({ transcript, messages?, annotations?, selectedAnnotationId?, showNotes?, onAnnotationSelect? })`

Creates a `TranscriptWorkspacePanel`. `transcript` may contain `title`, `subtitle`, `messages`, `annotations`, and `selectedAnnotationId`; top-level fields can override the bundle:

```js
widget.recipes.annotatedTranscript({
  transcript,
  selectedAnnotationId: "a1",
  onAnnotationSelect: "note-selected"
})
```

### `recipes.courseStudio({ sections, activeItemId?, title?, subtitle?, main?, onNavigate? })`

Creates a `CourseStudioShell`. The optional `main` field should be a Widget IR node, often a course slide or lesson panel:

```js
widget.recipes.courseStudio({
  sections,
  activeItemId: "slides",
  main: widget.recipes.courseSlide({ slide, snapshot, index: 0, total: 3 })
})
```

### `recipes.courseSlide({ slide, snapshot, index?, total?, visualSide?, onPrevious?, onNext? })`

Creates a `CourseSlidePanel` from a slide DTO and matching context-window snapshot:

```js
widget.recipes.courseSlide({ slide, snapshot, index: 0, total: slides.length, visualSide: "right" })
```

### `recipes.handout({ bundle, documents?, selectedDocumentId?, title?, onSelect?, onDownload?, onDownloadAll? })`

Creates a `HandoutDocumentShell`. `bundle.docs` is used as the document list when present:

```js
widget.recipes.handout({
  bundle,
  selectedDocumentId: "guide",
  onSelect: "document-selected",
  onDownload: "document-download"
})
```

Recipes are convenience helpers only. They return the same JSON-compatible Widget IR that could be written manually with low-level component helpers.

## Page object contract

A page endpoint should return a wrapper object around the root node:

```js
{
  schemaVersion: "0.1.0",
  id: "demo",
  title: "Demo",
  root: widget.panel({ title: "Demo" }, "Hello")
}
```

The required fields are:

| Field | Meaning |
| --- | --- |
| `schemaVersion` | Widget IR schema version. Current examples use `0.1.0`. |
| `id` | Page identifier used by the React app and API route. |
| `title` | Human-readable page title. |
| `root` | Widget IR node rendered by React. |

Additional metadata may be added when the consuming server and frontend agree on it.

## Error behavior

The module fails early for invalid core constructor calls:

| Call | Error |
| --- | --- |
| `element()` | Missing tag. |
| `element("")` | Empty tag. |
| `component()` | Missing component type. |
| `component("")` | Empty component type. |

Component helper calls are permissive. Unknown props are passed through as data. Unknown component types are possible through `component(type, ...)`, but they will only render if the React `WidgetRenderer` knows how to render that type.

## See Also

- `widget-dsl-getting-started`
- `tutorial-http-serve-jsverbs`
- `tutorial-static-assets-http-server`
- `buildspec-reference`
