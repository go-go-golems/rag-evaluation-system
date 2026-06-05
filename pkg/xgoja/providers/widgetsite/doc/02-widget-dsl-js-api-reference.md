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
  widget.button({ label: "Run" }),
  widget.button({ label: "Cancel", variant: "secondary" })
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
| `selectInput` | `SelectInput` | Select control. |
| `stack` | `Stack` | Vertical stack layout. |
| `statusText` | `StatusText` | Status label with optional icon. |
| `tabList` | `TabList` | Tab selector. |
| `textInput` | `TextInput` | Text input control. |

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
  label: "Refresh",
  variant: "primary",
  action: { kind: "server", name: "refresh", payload: { source: "demo" } }
})
```

The React app posts server actions to `/api/widget/actions/{name}`. The exact action contract is owned by the WidgetRenderer server or jsverb site. The DSL only preserves the action object.

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
