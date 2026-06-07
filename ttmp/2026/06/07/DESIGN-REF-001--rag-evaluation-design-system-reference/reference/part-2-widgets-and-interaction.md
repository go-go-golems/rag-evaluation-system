# Part II: Widget System, Hierarchy & Interaction

## 1. Widget IR: The Data Contract

All UI in the RAG evaluation system is described by a JSON-serialisable
intermediate representation called **Widget IR**. Every page, every panel, every
interactive control is a tree of Widget IR nodes.

### 1.1 Node types

```ts
type WidgetNode = TextNode | ElementNode | ComponentNode;

interface TextNode  { kind: 'text';   text: string }
interface ElementNode { kind: 'element'; tag: string; attrs?: JsonObject; children?: WidgetNode[] }
interface ComponentNode { kind: 'component'; type: string; props?: WidgetProps; children?: WidgetNode[] }
```

**Key rule**: the DSL (JavaScript) produces data (Widget IR), React consumes data. No
functions or closures cross the boundary. This keeps the JavaScript authoring
side pure-data and the rendering side pure-React.

### 1.2. Component Catalog

Each component below maps one-to-one to a React component.  The `type` string in the
Widget IR determines which React component the `WidgetRenderer` renders.

| Type | React Mapping | Key Props | Description |
|------|--------------|-----------|-------------|
| `Panel` | `<Panel>` | title, density, fill | Bordered card with optional header and actions |
| `Stack` | `<Stack>` | gap, align | Vertical layout; like a flex-col |
| `Inline` | `<Inline>` | gap, justify, wrap | Horizontal layout |
| `DataTable` | `<DataTable>` | columns, rows, onRowSelect, getRowKey | Tabular data with server-sortable columns |
| `StatusText` | `<StatusText>` | status (ok/error/warning/...) | Status badge with optional icon |
| `Button` | `<Button>` | variant, size, action | Clickable CTA |
| `Caption` | `<Caption>` | tone, transform | Small styled text label (e.g. "3 results") |
| `MetadataGrid` | `<MetadataGrid>` | items (key/value pairs) | Two-column key-value display |
| `AppShell` | `<AppShell>` | header, sidebar | Full-page chrome with nav |
| `DashboardGrid` | `<DashboardGrid>` | cols, recipe | Responsive CSS grid |
| `FormRow` | `<FormRow>` | label, control | Label + control pair |
| `SelectInput` | `<SelectInput>` | value, options | Dropdown select |
| `TextInput` | `<TextInput>` | value, placeholder | Single-line text input |
| `ScrollRegion` | `<ScrollRegion>` | axis | Scrollable container |
| `TabList` | `<TabList>` | items, onChange | Tab bar for sub-navigation |

### 1.3. Interaction Hierarchy

```
AppShell
 ├── AppNav (top-level navigation)
 ├── DashboardGrid
 │   ├── Panel (KPI Metrics)
 │   │   └── Stack
 │   │       ├── StatusText (ok/warning/error)
 │   │       └── MetricValue
 └── Panel (data detail)
     ├── DataTable
     └── MetadataGrid
```

The hierarchy flows from the `AppShell` at the top, through navigation (`AppNav`),
through content panels, down to atomic display components. Actions on lower
components propagate up via callbacks (onAction → handler).  The DSL never sends
functions across the wire — actions are serialisable specs (see §4).

### 1.4. Actions

Actions are serialisable specifications of user intent. The DSL produces action objects:

| kind | Usage | Example |
|------|-------|---------|
| `navigate` | Client-side page change | Navigate to another page within the SPA |
| `server` | RPC call to Go backend | Submit form, toggle a flag |
| `event` | Dispatch a DOM custom event | Show/hide panels |
| `copy` | Copy text to clipboard | Copy-to-clipboard button |

Action flow:

```
  [User clicks button]
     → action: { kind: "server", name: "reindex", payload: { id: "..." }
     → POST /api/widget/actions/reindex
     → Go server runs Goja, returns updated Widget IR
     → React re-renders from the new Widget IR
```

The **data contract** is: Goja scripts produce Widget IR; React consumes it. No
cross-boundary function calls or closures are permitted. This is the most
important invariant of the system.

### 1.5. The Page Lifecycle

1. Browser requests `/pages/overview`.
2. Server-side (Go) builds a Goja VM, executes the configured script.
3. Script calls `rag.panel(...)`, `rag.dataTable(...)`, etc., producing Widget IR.
4. The IR is returned as JSON.
5. The React `WidgetRenderer` walks the IR tree, mapping `component` nodes to React components.

All scripts are **trusted** (they run server-side in Goja). The system is
designed for internal tools where script authors are the same team that deploys
the server. The security boundary is at the API: the JS scripts cannot reach the
host; they can only produce Widget IR.

### 1.6. Semantic Recipes

For common patterns, the DSL provides semantic recipe functions that produce
pre-composed Widget IR:

```javascript
// Metrics row with 3 tiles
rag.recipes.metrics({ items: [
  { label: "Queries", value: "1284", trend: "up" },
  { label: "Latency", value: "243ms", trend: "down" },
  { label: "Uptime", value: "99.8%", trend: "flat" },
]})

// Toolbar with actions
rag.recipes.actionToolbar({
  title: "Queue",
  actions: [
    { label: "Run All", action: { kind: "server", name: "runAll" } }
  ]
})

// Master-detail table
rag.recipes.masterDetailTable({
  columns,
  rows,
  selectedKey: state.selectedId,
  onRowSelect: "select-user",
  detail: (row) => rag.panel({}, rag.metadataGrid({ items: [...] }))
})
```

### 1.7. File Locations

| Artifact | Path |
|---|---|
| Widget IR types | `web/src/widgets/ir.ts` |
| Widget renderer | `web/src/widgets/WidgetRenderer.tsx` |
| DSL Go module | `pkg/widgetdsl/` |
| Runner (Goja host) | `pkg/widgetrunner/` |
| HTTP server | `pkg/widgetserver/` |
| Embedded SPA | `pkg/defaultspa/` |
| NPM package | `packages/rag-evaluation-site/` |
| xgoja provider | `pkg/xgoja/providers/widgetsite/` |

---

## Quick Start for New Contributors

1. **Run the dev server**: `cd web && npm run dev`
2. **Add a widget**: edit `widgets/ir.ts` to add a new type, then implement
   a `case` in `WidgetRenderer.tsx`.
3. **Add a recipe**: write a function in the Goja DSL that produces Widget IR.
4. **Test it**: the WidgetRenderer renders your new node without custom React code.

