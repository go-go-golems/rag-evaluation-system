---
title: "Widget Hierarchy and Interaction Reference"
ticket: RAGEVAL-DS-RAGEVAL-QWEN36-HIGH
topics:
  - frontend-architecture
  - design-system
  - widget-ir
  - goja
  - composition
  - interaction
what_for: >
  Comprehensive reference for understanding the widget composition model, information
  hierarchy, interaction patterns, the Widget IR data format, the Goja authoring DSL,
  and the WidgetRenderer — the system that turns IR into React UI.
---

# RAG Evaluation System — Widget Hierarchy and Interaction Reference

> **Audience:** Frontend engineers, Goja script authors, and designers who need to
> understand the widget composition model, information hierarchy, interaction patterns,
> and the Widget IR data format.
>
> **Status:** Active · **Target version:** `web/` (RAG Evaluation System, Go 2025)

## 1. Executive Summary

The RAG Evaluation System uses a **layered component architecture** where UI is built
by composing named components from increasingly specific layers. Each layer has a
single responsibility, and components at each layer consume primitives from layers
below them without redefining their visual behavior.

The system also supports **Widget IR** — a JSON-compatible data format that allows
trusted JavaScript (Goja) scripts to author pages using the same component vocabulary
without importing React. The rule is simple: **Goja authors data; React renders UI**.

```
Design layers (composition):
  tokens → foundation → atoms → layout → molecules → organisms → pages → containers

Data format (Goja authoring):
  Widget IR → WidgetRenderer → React components → browser
```

## 2. The Composition Hierarchy

### 2.1 Layer Map

```
tokens.css
  │
  ├─ foundation/         ← Text, CodeText, StatusText, Caption, Divider, VisuallyHidden
  │     (text roles, status tones, accessibility primitives)
  │
  ├─ atoms/              ← Button, IconButton, TextInput, SelectInput, CheckboxRow, ErrorCallout
  │     (basic HTML control appearances)
  │
  ├─ layout/             ← Panel, Stack, Inline, DashboardGrid, AppShell, ScrollRegion, TabList, FormRow
  │     (structure: panels, grids, stacks, scrolling, forms, tabs)
  │
  ├─ molecules/          ← AppNav, DataTable, MetadataGrid, CoveragePanel, QueryPresetList
  │     (reusable data/display patterns that span features)
  │
  ├─ organisms/          ← SearchControlsPanel, RetrievalResultsPanel, ResultInspectorPanel,
  │                       QueueHealthPanel, WorkflowListPanel, WorkflowSummaryPanel,
  │                       WorkflowOpGraphPanel, WorkflowOpGroupsPanel, WorkflowOpInspectorPanel,
  │                       WorkflowOpResultPanel
  │     (feature panels with DTO-shaped props and callbacks)
  │
  ├─ pages/              ← SearchWorkbenchPage, PipelinePage, EvaluationPage, DslPreviewPage
  │     (page composition boundaries)
  │
  └─ corpus/             ← SourcePanel, DocumentBrowser, DocumentInspector, IdentityBar,
  │                       ArtifactIdentityBar, ChunkTimelineBar
          (corpus-domain widgets that adopt shared primitives)
```

### 2.2 The Ownership Rule (Recap)

Every visual decision has one owner:

| Layer | Owns | Does not own |
|---|---|---|
| `tokens.css` | color, font, border, spacing, semantic aliases | component selectors, feature anatomy |
| **Foundation** | text roles, code text, captions, status text, dividers, hidden text | grids, panels, API behavior |
| **Atoms** | buttons, icon buttons, inputs, selects, checkbox rows, error callouts | page state, data fetching, table rendering |
| **Layout** | panels, stacks, inline groups, dashboard grids, shells, tabs, form rows, scroll regions | domain formatting, API behavior |
| **Molecules** | reusable tables, metadata grids, coverage summaries, preset lists, app navigation | RTK Query hooks, page orchestration |
| **Organisms** | feature panels with DTO-shaped props and callbacks | global layout policy, app-wide state |
| **Pages** | composition of organisms and page-local state | low-level typography, raw CSS globals |
| **Containers** | RTK Query, mutations, navigation events, side effects | presentational styling decisions |

### 2.3 How Layers Compose

Each layer consumes from layers below. The composition direction is always **bottom-up**:

```
Page
  → Panel (layout)
       → Caption (foundation) for title
       → Stack (layout) for children spacing
            → DataTable (molecule) for data display
                 → Button (atom) for row actions
                     → text content (foundation/atoms)
            → Panel (layout) for detail section
                 → MetadataGrid (molecule)
                 → StatusText (foundation)
                 → ErrorCallout (atom) for errors
```

A Panel does not define button colors. A DataTable does not define panel headers.
Each component delegates to the appropriate layer.

## 3. Information Hierarchy

### 3.1 Density Scale

The system uses a **compact density model** optimized for dense operational dashboards
where many data points must be visible simultaneously.

| Density | Use case | Examples |
|---|---|---|
| **Condensed** | Maximum data in minimum vertical space | Panel density, table cells |
| **Normal** | Standard readable density | Most panels, forms, lists |
| **Spacious** | Less common; reserved for focused single-task views | — |

**Rule:** Default to normal density. Only use condensed when you are certain that
vertical space is constrained and the content is well-known to the user.

### 3.2 Panel as the Primary Content Container

The `Panel` is the workhorse container for all content:

```
<section class="root">
  <div class="header">                    ← Black bar, uppercase title
    <span>Search</span>
    <div class="actions">{actions}</div>  ← Optional: buttons, icons
  </div>
  <div class="body|condensed">            ← White area, content goes here
    {children}
  </div>
</section>
```

| Property | Normal | Condensed |
|---|---|---|
| Body padding | `8px` | `4px 6px` |
| Header | Always shown | Always shown |
| Fill mode | Default: content-sized | `fill=true`: 100% height |

**Rule:** Every panel has a header when it has a title or actions. Panels without titles
or actions render as a plain `section` with no header strip.

### 3.3 Metadata Display Pattern

The `MetadataGrid` is the standard pattern for key-value information display:

```
<dl class="root|compact">
  <div class="row">
    <dt class="key">Source</dt>
    <dd class="value">openai/text-embedding-3-small</dd>
  </div>
  <div class="row">
    <dt class="key">Model</dt>
    <dd class="value">
      text-embedding-3-small
      <button class="copyButton" onclick="copy(...)">⧉</button>
    </dd>
  </div>
</dl>
```

- Keys use `--rag-font-role-label` (bold mono)
- Values use `--rag-font-role-metadata` (mono)
- Copy button appears when `copyValue` is set (defaults to copying the value text)
- Compact mode reduces row spacing

### 3.4 Status Display Pattern

Status is always displayed with `StatusText`, which handles both the icon and color:

```tsx
<StatusText status="succeeded" icon>Succeeded</StatusText>
<StatusText status="running" icon>Processing</StatusText>
<StatusText status="failed">Failed</StatusText>
```

- `icon=true` prepends the Unicode icon: `✔ succeeded`, `● running`, `✘ failed`
- Icon is omitted by default — use it when you want to make status instantly scannable
- Color is determined by the status value (see Section 2.3 of the Styling reference)

### 3.5 Caption Pattern

Captions are used for secondary labels, hints, helper text, and categorical markers:

```tsx
<Caption>Caption — small, muted descriptive text.</Caption>
<Caption tone="accent" transform="uppercase">Retriever</Caption>
<Caption tone="warning">Partial completion</Caption>
<Caption transform="uppercase">Limits</Caption>
```

**Where captions appear:**
- Above form control groups as section labels (`<Caption transform="uppercase">Limits</Caption>`)
- Below controls as hints or helper text
- As inline metadata values in tables (via `CaptionCellSpec`)
- As error/warning labels in callouts

## 4. Interaction Hierarchy

### 4.1 Action Specification (Widget IR)

Actions in Widget IR are serialized data structures, not JavaScript callbacks. This allows
them to cross the HTTP boundary between Goja and React without serialization issues.

```ts
type ActionSpec =
  | NavigateActionSpec     // { kind: "navigate", to: string, params?: JsonObject }
  | ServerActionSpec       // { kind: "server", name: string, payload?: JsonObject }
  | EventActionSpec        // { kind: "event", event: string, detail?: JsonObject }
  | CopyActionSpec;        // { kind: "copy", value?: string, field?: string }
```

### 4.2 Action Types

| Kind | React behavior | Backend path |
|---|---|---|
| `navigate` | Router navigation to `to` path | None (client-side only) |
| `server` | POST to `/api/widget/actions/{name}` | Goja action function |
| `event` | Dispatch browser event `widget:{event}` | Optional handler |
| `copy` | Copies `value` to clipboard | None (client-side only) |

### 4.3 Action Binding in WidgetRenderer

The WidgetRenderer binds actions to component handlers:

```tsx
// Button binding
<Button
  onClick={bindAction(props.action, { componentType: 'Button' }, onAction)}
>
  {children}
</Button>

// DataTable row select binding
<DataTable
  onRowSelect={(row) => bindAndRun(rowSelectAction, { row, rowKey, componentType: 'DataTable' }, onAction)}
/>
```

The binding context includes `componentType` and `value` so that action functions can
inspect what triggered the action:

```js
// In Goja script
exports.actions = {
  "select-query": function(ctx, payload) {
    // ctx = { action: "select-query", value: "row-key-42", componentType: "DataTable" }
    // payload = { row: {...}, rowKey: "row-key-42" }
    return { refresh: true, data: { selectedId: ctx.value } };
  }
};
```

### 4.4 Server Actions

Server actions follow this flow:

```
React click → POST /api/widget/actions/{name} → Goja action function →
  { ok, refresh, toast, patch, data } → React refresh or toast event
```

The action result determines what happens next:

| Result field | Effect |
|---|---|
| `ok: true` | Action succeeded |
| `ok: false` | Action failed |
| `refresh: true` | Reload the current page from the server |
| `toast: "message"` | Dispatch `widget:toast` browser event (not yet rendered) |
| `patch: {...}` | Apply client-side state patch |
| `data: {...}` | Additional data returned to the client |

**Rule:** Missing `ok` defaults to `true`. Simple actions can return `{ refresh: true }`
without also writing `{ ok: true }`.

### 4.5 Cell Actions (DataTable)

Table cells can carry actions through the cell specification. A cell spec is serialized
data that the WidgetRenderer converts into a render function:

```ts
type CellSpec =
  | { kind: "field", field: string, fallback?: string }
  | { kind: "number", field: string, format?: "integer"|"fixed", digits?: number, fallback?: string }
  | { kind: "status", field: string, icon?: boolean, fallback?: RagStatus|string }
  | { kind: "caption", field: string, tone?: CaptionTone, fallback?: string }
  | { kind: "template", template: string }
  | { kind: "link", hrefField: string, labelField: string, target?: string, fallbackLabel?: string }
  | { kind: "constant", value: RenderableValue };
```

**Important:** Table cells cannot be JavaScript functions. They are data specifications
that the WidgetRenderer converts into render functions internally.

### 4.6 Form Interaction Pattern

Forms use `FormRow` for each field:

```tsx
<FormRow
  label="Index ID"
  control={<TextInput value={indexId} onChange={e => setIndexId(e.target.value)} />}
  hint="The BM25 index to search against"
/>
```

Form rows are typically wrapped in `Stack` with a gap:

```tsx
<Stack gap="sm">
  <FormRow label="..." control={<TextInput ... />} />
  <FormRow label="..." control={<SelectInput ... />} />
</Stack>
```

**Rule:** Form controls own their appearance (atoms). The form row owns the label/control
pairing and hint. The panel owns the overall layout. Never set form control styling
inside a panel's CSS Module.

### 4.7 Tab Navigation Pattern

Tabs use `TabList` with an `onChange` action:

```tsx
<TabList
  items={[
    { id: "summary", label: "Summary" },
    { id: "artifacts", label: "Artifacts" },
    { id: "chunks", label: "Chunks" },
  ]}
  activeId="artifacts"
  onChange={(id) => setSelectedTab(id)}
  ariaLabel="Document detail tabs"
/>
```

The tab bar is typically placed inside a `Panel` header via the `actions` prop, or
below the title as the first element in the panel body.

## 5. Widget IR Data Format

### 5.1 The Three Node Kinds

Widget IR is a tree of nodes. Every node is one of three kinds:

```ts
type WidgetNode = TextNode | ElementNode | ComponentNode;
```

| Kind | Purpose | Example |
|---|---|---|
| `text` | Plain text content | `{ kind: "text", text: "Hello" }` |
| `element` | Generic HTML element | `{ kind: "element", tag: "div", attrs: { class: "foo" }, children: [...] }` |
| `component` | React component instance | `{ kind: "component", type: "Panel", props: { title: "Test" }, children: [...] }` |

### 5.2 Component Node Props

Component nodes map to real React components. The `type` field is the component name:

```json
{
  "kind": "component",
  "type": "Panel",
  "props": {
    "title": "Retrieval Results",
    "density": "condensed"
  },
  "children": [
    { "kind": "text", "text": "Table content" }
  ]
}
```

### 5.3 Recognized Widget Types

The `RagWidgetType` union includes all widget types that the WidgetRenderer knows how
to render:

```
AppShell, AppNav, Button, Caption, DashboardGrid, DataTable, FormRow, Inline,
MetadataGrid, Panel, ScrollRegion, SelectInput, Stack, StatusText, TabList, TextInput
```

Any other `type` value renders as an `ErrorCallout` with "Unknown widget: {type}".

### 5.4 Renderable Values

Many props accept `RenderableValue`:

```ts
type RenderableValue = WidgetNode | string | number | boolean | null;
```

This means you can nest Widget IR nodes inside any prop that accepts `RenderableValue`:

```json
{
  "type": "AppNav",
  "props": {
    "brand": { "kind": "text", "text": "◉ RAG Eval" },
    "items": [
      { "id": "search", "label": "Search" },
      { "id": "corpus", "label": { "kind": "text", "text": "Corpus" } }
    ]
  }
}
```

### 5.5 Widget IR Factory Functions

The IR package provides factory functions for building IR programmatically:

```ts
// Text node
text("Hello World")                    // → { kind: "text", text: "Hello World" }

// Element node
element("div", { className: "foo" })   // → { kind: "element", tag: "div", attrs: { className: "foo" } }

// Component node
component("Panel", { title: "Test" }, [text("content")])
// → { kind: "component", type: "Panel", props: { title: "Test" }, children: [{ kind: "text", text: "content" }] }

// Check if value is a widget node
isWidgetNode(value)                    // → boolean
```

## 6. Goja Widget DSL

### 6.1 The DSL Philosophy

The Widget DSL lets trusted JavaScript scripts author pages using the same component
vocabulary as the React app — without importing React. The rule is:

> **Goja authors data; React renders UI**

The DSL returns JSON-compatible objects (Widget IR), not HTML or React elements.

### 6.2 Loading the DSL

```js
const rag = require("widget.dsl");
// or
const rag = require("rag.dsl");  // RAG-oriented alias
```

### 6.3 Low-Level Helpers

| Helper | Purpose | Example |
|---|---|---|
| `text(value)` | Create text node | `rag.text("Hello")` |
| `element(tag, attrs?, children?)` | Create HTML element node | `rag.element("div", { className: "foo" })` |
| `component(type, props?, children?)` | Create component node | `rag.component("Panel", { title: "X" })` |
| `fragment(children?)` | Create fragment (array children) | `rag.fragment([rag.text("a"), rag.text("b")])` |

### 6.4 Component Helpers

These are convenience wrappers around `component()`:

| Helper | Maps to |
|---|---|
| `rag.panel(props?, ...children)` | `{ type: "Panel" }` |
| `rag.stack(props?, ...children)` | `{ type: "Stack" }` |
| `rag.inline(props?, ...children)` | `{ type: "Inline" }` |
| `rag.dashboardGrid(props?, ...children)` | `{ type: "DashboardGrid" }` |
| `rag.appShell(props?, ...children)` | `{ type: "AppShell" }` |
| `rag.appNav(props?, ...children)` | `{ type: "AppNav" }` |
| `rag.button(props?, ...children)` | `{ type: "Button" }` |
| `rag.caption(props?, ...children)` | `{ type: "Caption" }` |
| `rag.statusText(props?, ...children)` | `{ type: "StatusText" }` |
| `rag.dataTable(props?, ...children)` | `{ type: "DataTable" }` |
| `rag.metadataGrid(props?, ...children)` | `{ type: "MetadataGrid" }` |
| `rag.tabList(props?, ...children)` | `{ type: "TabList" }` |
| `rag.formRow(props?, ...children)` | `{ type: "FormRow" }` |
| `rag.scrollRegion(props?, ...children)` | `{ type: "ScrollRegion" }` |
| `rag.textInput(props?, ...children)` | `{ type: "TextInput" }` |
| `rag.selectInput(props?, ...children)` | `{ type: "SelectInput" }` |

### 6.5 Cell Helpers

Table cell specs:

| Helper | Maps to |
|---|---|
| `rag.cell.field(field, fallback?)` | `{ kind: "field", field, fallback }` |
| `rag.cell.number(field, format?, digits?, fallback?)` | `{ kind: "number", field, format, digits, fallback }` |
| `rag.cell.status(field, icon?, fallback?)` | `{ kind: "status", field, icon, fallback }` |
| `rag.cell.caption(field, tone?, fallback?)` | `{ kind: "caption", field, tone, fallback }` |
| `rag.cell.template(template)` | `{ kind: "template", template }` |
| `rag.cell.link(hrefField, labelField, target?, fallbackLabel?)` | `{ kind: "link", ... }` |
| `rag.cell.constant(value)` | `{ kind: "constant", value }` |

### 6.6 Action Helpers

| Helper | Purpose |
|---|---|
| `rag.action.server(name, options?)` | Create a server action spec |
| `rag.action.navigate(to, options?)` | Create a navigate action spec |
| `rag.action.event(event, options?)` | Create an event action spec |
| `rag.action.copy(value)` | Create a copy action spec |

### 6.7 Example: Simple Page

```js
const rag = require("widget.dsl");

exports.pages = {
  demo(ctx) {
    return {
      schemaVersion: "0.1.0",
      id: ctx.pageId,
      title: "Demo",
      root: rag.panel({ title: "Hello World" },
        rag.text("This is a simple page.")
      )
    };
  }
};
```

### 6.8 Example: DataTable with Cell Specs

```js
const rag = require("widget.dsl");

exports.pages = {
  queries(ctx) {
    const rows = ctx.data.rows; // from database query
    return {
      schemaVersion: "0.1.0",
      id: "queries",
      title: "Query Queue",
      root: rag.panel({ title: "Queries" },
        rag.dataTable({
          rows: rows,
          getRowKey: "id",
          columns: [
            { id: "id",     header: "ID",    cell: rag.cell.field("id") },
            { id: "status", header: "Status", cell: rag.cell.status("status", true) },
            { id: "score",  header: "Score", cell: rag.cell.number("score", "fixed", 4) },
          ]
        })
      )
    };
  }
};
```

### 6.9 Example: Page with Actions

```js
const rag = require("widget.dsl");

exports.pages = {
  demo(ctx) {
    return {
      schemaVersion: "0.1.0",
      id: "demo",
      title: "Actions Demo",
      root: rag.panel({ title: "Actions" },
        rag.button({
          label: "Click me",
          variant: "primary",
          action: rag.action.server("increment", { payload: { source: "demo" } })
        })
      )
    };
  }
};

exports.actions = {
  increment(ctx, payload) {
    clicks = (clicks || 0) + 1;
    return { ok: true, refresh: true, toast: "Clicked " + clicks + " times" };
  }
};
```

### 6.10 Semantic Recipes

Recipes are macros that expand into Widget IR for common dashboard patterns:

```js
rag.recipes.metrics({ items: [{ label, value, status? }] })
// → DashboardGrid of condensed Panels, each with a StatusText

rag.recipes.actionToolbar({ title, actions, caption })
// → Inline toolbar with ActionSpec buttons

rag.recipes.masterDetailTable({ rows, columns, selectedKey, onRowSelect, detail })
// → Two-column DashboardGrid: DataTable panel + detail panel
```

**Important:** Recipe callbacks (like `detail`) are evaluated only while constructing the
server response. They return Widget IR nodes, which are exported to JSON-compatible data
before crossing the HTTP boundary.

## 7. WidgetRenderer Architecture

### 7.1 The Renderer

The WidgetRenderer receives a Widget IR node and renders it as React elements:

```tsx
function renderWidgetNode(node: WidgetNode, onAction?): ReactNode {
  if (node.kind === 'text') return node.text;
  if (node.kind === 'element') return <node.tag>{renderChildren(node.children)}</node.tag>;
  // switch on node.type → real React component
  return renderComponentNode(node, onAction);
}
```

### 7.2 Component Dispatch Table

| IR type | WidgetRenderer function | React component |
|---|---|---|
| `AppShell` | `renderAppShell` | `AppShell` |
| `AppNav` | `renderAppNav` | `AppNav` |
| `Button` | `renderButton` | `Button` |
| `Caption` | `renderCaption` | `Caption` |
| `DashboardGrid` | `renderDashboardGrid` | `DashboardGrid` |
| `DataTable` | `renderDataTable` | `DataTable` |
| `FormRow` | `renderFormRow` | `FormRow` |
| `Inline` | `renderInline` | `Inline` |
| `MetadataGrid` | `renderMetadataGrid` | `MetadataGrid` |
| `Panel` | `renderPanel` | `Panel` |
| `ScrollRegion` | `renderScrollRegion` | `ScrollRegion` |
| `SelectInput` | `renderSelectInput` | `SelectInput` |
| `Stack` | `renderStack` | `Stack` |
| `StatusText` | `renderStatusText` | `StatusText` |
| `TabList` | `renderTabList` | `TabList` |
| `TextInput` | `renderTextInput` | `TextInput` |
| (unknown) | default case | `ErrorCallout` |

### 7.3 Cell Renderer Pipeline

Table cells use a two-stage pipeline:

```
CellSpec (data) → renderCell() → ReactNode (UI)
```

| CellSpec kind | Output |
|---|---|
| `field` | Plain text (stringified value) |
| `number` | Formatted number (integer or fixed decimal) |
| `status` | `<StatusText>` with icon and color |
| `caption` | `<Caption>` with tone |
| `template` | String with `${field}` and `$field` interpolation |
| `link` | `<a href>` element |
| `constant` | Rendered renderable value |

### 7.4 Row Key Resolution

Row keys can be specified in three ways:

```ts
type RowKeySpec = string | { field: string } | { template: string };

// "id" → row["id"]
// { field: "id" } → row["id"]
// { template: "row-${prefix}-${id}" } → interpolated path
```

The template syntax supports both `${path}` and `$path` (dot-separated paths work):

```
"${user.name}"  → row.user.name
"$user.email"   → row.user.email
```

## 8. Typical Page Layout Patterns

### 8.1 Search Workbench (The Canonical Pattern)

```
AppShell (default shell)
  → DashboardGrid (recipe: "searchWorkbench")
       → SearchControlsPanel (panel: "Search", condensed)
            → Stack (gap: "sm")
                 → TextInput + Button (query row)
                 → Caption (retriever label)
                 → Button group (bm25/vector/hybrid)
                 → FormRows (index ID, strategy, profile)
                 → ScrollRegion with CheckboxRows (source filter)
                 → FormRows (limits)
       → RetrievalResultsPanel (panel: "Results")
            → DataTable with cell specs
       → ResultInspectorPanel (panel: "Inspector")
            → MetadataGrid (metadata)
            → DataTable (chunks)
            → Divider
            → MetadataGrid (source document)
```

### 8.2 Corpus Explorer Pattern

```
AppShell
  → DashboardGrid (recipe: "corpusExplorer")
       → SourcePanel (left: source list)
       → DocumentBrowser (center: document list)
            → DataTable + TextInput (search)
       → DocumentInspector (right: detail)
            → TabList
            → Panel (metadata)
            → Panel (artifacts)
            → Panel (chunks)
```

### 8.3 Workflow Detail Pattern

```
AppShell
  → DashboardGrid (recipe: "twoColumn")
       → WorkflowOpGraphPanel (left: graph)
       → WorkflowOpInspectorPanel (right: detail)
            → Panel (metadata)
            → ErrorCallout (if failed)
            → MetadataGrid (operation details)
            → WorkflowOpResultPanel (result)
```

### 8.4 Generic Dashboard Page Pattern

```tsx
// JSX example
<DashboardGrid recipe="twoColumn">
  <Panel title="Queue Health" density="condensed">
    <StatusText status={health} icon>{health}</StatusText>
    <DataTable rows={items} columns={columns} getRowKey={r => r.id} />
  </Panel>
  <Panel title="Details">
    <MetadataGrid items={metadata} />
    <Stack gap="sm">
      <FormRow label="Field" control={<TextInput value={value} />} />
    </Stack>
  </Panel>
</DashboardGrid>
```

## 9. File Reference

| File | Purpose |
|---|---|
| `web/src/widgets/ir.ts` | Widget IR TypeScript types and factory functions |
| `web/src/widgets/WidgetRenderer.tsx` | The WidgetRenderer — maps IR nodes to React components |
| `web/src/widgets/cellRenderers.tsx` | CellSpec → ReactNode conversion pipeline |
| `web/src/widgets/actions.ts` | Action binding and dispatch |
| `web/src/widgets/index.ts` | Widget IR exports |
| `pkg/widgetdsl/` | Goja native module: `require("widget.dsl")` and `require("rag.dsl")` |
| `pkg/widgetrunner/` | Script loading, page lookup, action invocation |
| `pkg/widgetserver/` | HTTP server: `/api/widget/pages`, `/api/widget/actions`, `/api/widget/schema` |
| `packages/rag-evaluation-site/` | Reusable React npm package with WidgetRenderer |
| `pkg/xgoja/providers/widgetsite/` | xgoja provider wrapping the DSL |

## 10. Quick Reference: Component Props

### 10.1 Panel

| Prop | Type | Default | Description |
|---|---|---|---|
| `title` | ReactNode | — | Panel header text (black bar) |
| `actions` | ReactNode | — | Right-aligned actions in header |
| `density` | `"normal"` \| `"condensed"` | `"normal"` | Body padding |
| `fill` | boolean | `false` | 100% height |

### 10.2 Stack

| Prop | Type | Default | Description |
|---|---|---|---|
| `gap` | `"xs"` \| `"sm"` \| `"md"` \| `"lg"` \| `"xl"` | — | Vertical/horizontal spacing |
| `align` | `"start"` \| `"center"` \| `"end"` \| `"stretch"` | — | Cross-axis alignment |

### 10.3 Inline

| Prop | Type | Default | Description |
|---|---|---|---|
| `gap` | `"xs"` \| `"sm"` \| `"md"` \| `"lg"` \| `"xl"` | — | Horizontal spacing |
| `justify` | `"start"` \| `"center"` \| `"end"` \| `"between"` | — | Main-axis alignment |
| `wrap` | boolean | `false` | Allow wrapping |

### 10.4 DataTable

| Prop | Type | Description |
|---|---|---|
| `columns` | `ColumnSpec[]` | Column definitions with header, cell spec, alignment |
| `rows` | `JsonObject[]` | Data rows |
| `getRowKey` | string \| `{ field }` \| `{ template }` | Row key resolution |
| `selectedKey` | string \| null | Currently selected row |
| `onRowSelect` | `(row) => void` \| ActionSpec | Row selection callback or action |
| `emptyMessage` | ReactNode | Message when rows is empty |

### 10.5 MetadataGrid

| Prop | Type | Default | Description |
|---|---|---|---|
| `items` | `[{key, value, copyValue?}]` | — | Key-value pairs |
| `density` | `"normal"` \| `"compact"` | `"normal"` | Row spacing |

### 10.6 DashboardGrid

| Prop | Type | Default | Description |
|---|---|---|---|
| `recipe` | `"searchWorkbench"` \| `"corpusExplorer"` \| `"twoColumn"` | `"twoColumn"` | Grid layout pattern |

### 10.7 StatusText

| Prop | Type | Default | Description |
|---|---|---|---|
| `status` | `RagStatus \| string` | — | One of 10 status values |
| `icon` | boolean | `false` | Show Unicode icon |

### 10.8 Caption

| Prop | Type | Default | Description |
|---|---|---|---|
| `tone` | `"muted"` \| `"accent"` \| `"success"` \| `"warning"` \| `"danger"` \| `"inherit"` | `"muted"` | Color tone |
| `transform` | `"none"` \| `"uppercase"` | `"none"` | Text transformation |
| `truncate` | boolean | `false` | Single-line with ellipsis |

### 10.9 Button

| Prop | Type | Default | Description |
|---|---|---|---|
| `variant` | `"default"` \| `"primary"` | `"default"` | Visual variant |
| `size` | `"normal"` \| `"compact"` | `"normal"` | Padding and font size |
| `selected` | boolean | `false` | Selected state (same as primary) |
| `disabled` | boolean | `false` | Disabled state |

## 11. Working Rules

1. **Always decide the layer first.** Before writing a component, determine which layer owns its visual responsibility using the ownership table.

2. **Compose from bottom up.** A Panel contains a Stack, which contains a DataTable, which contains Buttons. Each component consumes from lower layers.

3. **Never redefine primitives in molecules or organisms.** If a DataTable needs a panel header, use the `Panel` component with `title` and `actions` props.

4. **Use status vocabulary, not custom colors.** For operational status, always use `StatusText` with one of the 10 recognized statuses.

5. **Data attributes for tooling.** Always set `data-rag-*` attributes on public components for Playwright smokes and visual diff tools.

6. **Actions are data, not callbacks.** In Widget IR, actions are serializable objects. In React, actions are bound through the `onAction` handler.

7. **Cell specs replace callbacks.** Table cells use `CellSpec` objects, not JavaScript functions. This is required because functions cannot cross the HTTP boundary.

8. **Storybook is the review surface.** Every public component should have at least one story demonstrating its default state and key variants.

9. **Recipes are macros, not a second renderer.** `rag.recipes.*` expand into Widget IR. They are authoring conveniences, not rendering layers.

10. **The standalone token bridge is required.** When extracting components to a standalone package, ensure the token bridge maps `--mac-*` to `--rag-*` variables.
