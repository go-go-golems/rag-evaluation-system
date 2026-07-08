---
Title: 'Calendar and Scheduling Widgets: Analysis, Design, and Implementation Guide'
Ticket: RAGEVAL-SCHEDULE-WIDGETS
Status: active
Topics:
    - ui-dsl
    - widget-ir
    - design-system
    - react
    - frontend-architecture
    - intern-guide
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/StepList
      Note: Canonical six-file widget folder layout to copy
    - Path: repo://packages/rag-evaluation-site/src/context/types.ts
      Note: ContextStyleSet palette contract + DTO precedent the scheduling domain mirrors
    - Path: repo://packages/rag-evaluation-site/src/widgets/cellRenderers.tsx
      Note: renderCell/CellSpec — the defunctionalized render-lambda template the engines copy
    - Path: repo://packages/rag-evaluation-site/src/widgets/defaultRegistry.ts
      Note: per-module registries; where scheduleWidgetRegistry/timeWidgetRegistry get merged
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/index.ts
      Note: Widget IR barrel after ir.ts split; re-exports core/actions/cells/engines/props
    - Path: repo://packages/rag-evaluation-site/src/widgets/registry.ts
      Note: defineWidget, RenderContext, createWidgetRegistry (duplicate-type gotcha), mergeWidgetRegistries
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-06T18:48:57.627411111-04:00
WhatFor: ""
WhenToUse: ""
---




# Calendar and Scheduling Widgets: Analysis, Design, and Implementation Guide

> **Audience:** a new intern who has *not* worked in this codebase before.
> By the end of this document you should understand (a) how the existing widget
> system works, (b) why we build scheduling UI as *generic engines + domain
> presets* instead of one-off components, and (c) exactly what to build first.
> Read it top to bottom once; then use the "API Reference" and "File Reference"
> sections as lookup tables while you code.

---

## 0. How to read this document

This guide is long on purpose. It is structured as a funnel:

1. **Part A — The world as it is.** The design system, the Widget IR, the
   registry, and the one pattern that already lives in the code:
   "lambda-as-data". You must understand this before writing a line.
2. **Part B — The idea.** The three-layer architecture (base engine → IR adapter
   → DSL preset) and the *cell contract* that makes cells swappable.
3. **Part C — The widgets.** Each generic engine and each domain preset, with
   ASCII mockups, prop tables, and pseudocode.
4. **Part D — The backend.** Server actions, REST endpoints, and what the server
   computes so the widgets stay dumb.
5. **Part E — Build order.** What to implement first, the file checklist per
   widget, and the Storybook states expected.

Terms in `monospace` are real identifiers you can grep for. Paths like
`packages/rag-evaluation-site/src/...` are real files.

---

# PART A — The world as it is

## A.1 What package are we in?

Everything lives in **`packages/rag-evaluation-site`**. This is a *strict design
system* — a reusable UI library, not an app. The golden rules
(`GUIDELINES.md`) that constrain us:

- **Layered.** Code is sorted into `foundation → atoms → layout → molecules →
  organisms → widgets`. Lower layers never import higher layers.
- **Presentational only.** Package components **must not** import RTK Query hooks,
  app stores, router state, or backend services. They take data as props and emit
  callbacks. The *app* (in `packages/web`) decides where data comes from.
- **Tokenized typography.** Use `Text`, `Caption`, `CodeText`, `StatusText` and
  the `--rag-font-role-*` CSS variables. Do **not** hand-write `font-size: 13px`.
- **`data-rag-*` attributes** on every public component (used for visual-diff
  extraction and IR inspection).
- **React first, Widget IR later.** Stabilize the React component and its
  Storybook states *before* exposing it to the DSL.

Layer definitions (memorize these — they decide where your file goes):

```
foundation  text/semantic roles: Text, Caption, CodeText, StatusText, Divider
atoms       small controls & markers: Button, TextInput, Tag, MeterBar, IconButton
layout      generic structure: AppShell, Panel, Stack, Inline, SplitPane, TileGrid
molecules   reusable data/content patterns: DataTable, MetadataGrid, StepList
organisms   feature panels w/ DTO props: CourseLessonPanel, TranscriptReaderPanel
widgets     the IR + WidgetRenderer + per-DSL registries
```

> **Litmus test for layer.** If a component answers *"where do regions go?"* it is
> **layout**. If it answers *"what domain data is shown?"* it is a **molecule /
> organism**. If it answers *"how are cells arranged in space or time?"* — which
> is what our grids and calendars do — it is a **generic molecule** that must
> stay domain-blind.

## A.2 The file layout convention (every public component)

Look at `src/components/molecules/StepList/`. Every DSL-exposed component is a
folder with exactly these files:

```
StepList/
  StepList.tsx            # the React component (props in, JSX out)
  StepList.module.css     # local anatomy only (tokens, no raw colors)
  StepList.stories.tsx    # Storybook — the canonical review surface
  StepList.widget.tsx     # the IR adapter (maps IR props -> React props)
  StepList.widget.yaml    # the widget manifest (type, module, helper, slots)
  index.ts                # re-exports
```

You will create this same six-file set for each new widget. The Storybook file is
**not optional** — it is how we (and you) review the widget.

## A.3 The Widget IR — UI as JSON

The core idea of the "widget" layer: a UI can be described as a **JSON tree** and
rendered by a single component. This lets a JavaScript DSL (running in Goja, a Go
JS engine) *emit UI* without shipping React.

The node model lives in **`src/widgets/ir.ts`**:

```ts
// src/widgets/ir.ts
type WidgetNode = TextNode | ElementNode | ComponentNode;

interface TextNode      { kind: "text";      text: string; }
interface ElementNode   { kind: "element";   tag: string; attrs?; children?: WidgetNode[]; }
interface ComponentNode { kind: "component"; type: RagWidgetType | string;
                          props?: WidgetProps; children?: WidgetNode[]; }
```

A `ComponentNode` is the interesting one: `type` is a string key like
`"StepList"`, `props` is a JSON-compatible object, `children` are more nodes.
There is a giant union `RagWidgetType` listing every registered widget type, and
a matching `*WidgetProps` interface for each.

Helper constructors at the bottom of the file:

```ts
text("hi")                              // -> TextNode
element("div", { id: "x" }, [ ... ])    // -> ElementNode
component("StepList", { items: [...] }) // -> ComponentNode
```

## A.4 The renderer — how JSON becomes React

**`src/widgets/WidgetRenderer.tsx`** walks the node tree. The important function:

```tsx
// src/widgets/WidgetRenderer.tsx  (paraphrased)
function renderComponentNode(node, ctx, registry) {
  const adapter = registry.get(node.type);            // <-- flat lookup by string
  if (!adapter) return <UnknownWidget node={node} />;
  return adapter.render(node.props ?? {}, renderChildren(...), ctx, node);
}
```

Two things to notice, because the whole cross-DSL composition story depends on
them:

- **Lookup is a flat `registry.get(node.type)`.** The renderer does not care which
  DSL/module an adapter came from. Any node type in the registry can nest any
  other. This is *why* you can freely mix `schedule.dsl` and `calendar.dsl` nodes
  in one tree.
- **The renderer hands each adapter a `RenderContext` (`ctx`).** This is the
  bridge that turns serializable specs back into live behavior.

The `RenderContext` (from `src/widgets/registry.ts`):

```ts
interface RenderContext {
  renderNode(node): ReactNode;                 // render a child node
  renderChildren(children?): ReactNode[];      // render a node array
  renderValue(value): ReactNode;               // render a RenderableValue (node OR string)
  bindAction(action, context): (() => void) | undefined;
  dispatchAction(action, context): void;       // run an ActionSpec
}
```

## A.5 The adapter + registry

An **adapter** connects an IR `type` string to a real React component. Defined
with `defineWidget` (`src/widgets/registry.ts`). Real example
(`src/components/molecules/StepList/StepList.widget.tsx`):

```tsx
export const stepListWidget = defineWidget<StepListWidgetProps>({
  type: "StepList",
  module: "ui.dsl",
  render: (props, _children, ctx) => (
    <StepList
      items={props.items.map(item => ({
        ...item,
        index: ctx.renderValue(item.index),      // spec value -> ReactNode
        title: ctx.renderValue(item.title),
      }))}
      onItemSelect={
        props.onItemSelectAction
          ? (itemId) => ctx.dispatchAction(props.onItemSelectAction!, { itemId, value: itemId })
          : undefined                            // ActionSpec -> real callback
      }
    />
  ),
});
```

Adapters are grouped into **per-module registries** and merged
(`src/widgets/defaultRegistry.ts`):

```ts
export const uiWidgetRegistry      = createWidgetRegistry([ stepListWidget, panelWidget, ... ]);
export const dataWidgetRegistry    = createWidgetRegistry([ dataTableWidget ]);
export const courseWidgetRegistry  = createWidgetRegistry([ courseLessonPanelWidget, ... ]);
// ...
export const defaultWidgetRegistry = mergeWidgetRegistries(
  uiWidgetRegistry, dataWidgetRegistry, contextWindowWidgetRegistry,
  courseWidgetRegistry, cmsWidgetRegistry,
);
```

> **CRITICAL GOTCHA.** `createWidgetRegistry` **throws on duplicate `type`**:
> `if (byType.has(adapter.type)) throw new Error("Duplicate widget adapter ...")`.
> The `type` string is a **global namespace** even though `module` groups the
> builders. Two DSLs must never both register `"Grid"`. We solve this with
> namespaced type strings (see B.4).

## A.6 The pattern that already exists: "lambda-as-data"

This is the single most important idea to internalize. **You cannot serialize a
JavaScript closure into JSON** and send it across the Goja boundary. So the IR
represents callbacks and renderers as *data that an interpreter executes*. This is
formally called **defunctionalization**. The codebase already has three of these:

| Concept | "Real" closure | Defunctionalized spec (JSON) | Interpreter |
|---|---|---|---|
| Event handler | `onClick={() => navigate("/x")}` | `{ kind: "navigate", to: "/x" }` | `dispatchWidgetAction` |
| Cell renderer | `cell={(row) => <StatusText/>}` | `{ kind: "status", field: "status" }` | `renderCell` |
| Value accessor | `(row) => row.id` | `{ field: "id" }` or `"${a}-${b}"` | `getPath` / `renderTemplate` |

Look at **`src/widgets/actions.ts`** — `ActionSpec` is a tagged union:

```ts
type ActionSpec =
  | { kind: "navigate"; to: string; params?; }
  | { kind: "download"; to: string; }
  | { kind: "server";   name: string; payload?; }   // POST /api/widget/actions/<name>
  | { kind: "event";    event: string; detail?; }    // window CustomEvent (print, fullscreen...)
  | { kind: "copy";     value?; field?; };           // clipboard
```

And **`src/widgets/cellRenderers.tsx`** — `renderCell(spec, row, ...)` is a big
`switch (spec.kind)` interpreting `CellSpec` variants (`field`, `number`,
`status`, `caption`, `template`, `link`, `linkButton`, `actionButton`,
`constant`). **This is our template.** Every generic engine we build will accept a
`CellSpec`-shaped renderer and an `ActionSpec`-shaped callback, exactly like
`DataTable` already does.

The fourth mechanism is the **slot**: any prop typed `RenderableValue` /
`WidgetNode` is a hole you can drop a whole subtree into. `Panel.actions`,
`SlideShell.primary`, `EmptyState.actionSlot` all do this. Slots are how a caller
"passes in their own list of grid cells."

## A.7 The palette contract (reuse, do not reinvent)

Colors for anything "categorical" go through the existing `styleKey +
ContextStyleSet` contract (`src/context/types.ts`). Do **not** invent a new color
system for events/availability.

```ts
interface ContextVisualStyle { fill: string; line?; labelColor?; pattern?; ... }
interface ContextStyleSet {
  styles: Record<string /*styleKey*/, ContextVisualStyle>;
  legend: ContextLegendItemSpec[];
  fallbackStyle?: ContextVisualStyle;
}
```

A datum carries a `styleKey` (a *lookup key*, never a raw color); the `styleSet`
maps it to a visual style. A calendar event's color and a poll slot's "best"
highlight both resolve through this. `ContextBudgetBar` and the context diagrams
already consume it — read one for a worked example.

---

# PART B — The idea: three layers + one contract

## B.1 The problem with one-off scheduling widgets

The naive approach is to write a `PollGrid` that knows about `response.cells`, a
`CalendarMonth` that knows about `events`, a `BookingDay` that knows about slots —
each a bespoke component. That is three grids, three calendars' worth of
sticky-header / scroll / keyboard-nav / selection logic, duplicated and subtly
divergent. It also produces zero reuse: nobody outside "Doodle" can use a
`PollGrid`.

## B.2 The three-layer architecture

Instead, split every widget into three roles:

```
        ┌─────────────────────────────────────────────────────────────┐
        │ LAYER 3 — DSL preset  (schedule.dsl / calendar.dsl)           │
        │   availabilityMatrix(poll)   monthCalendar(events)            │
        │   one arg, domain vocabulary, opinionated defaults baked in   │
        └───────────────▲─────────────────────────────────────────────┘
                        │ configures (accessor/render/style/action specs)
        ┌───────────────┴─────────────────────────────────────────────┐
        │ LAYER 2 — IR adapter  (MatrixGrid.widget.tsx)                 │
        │   maps serializable specs  ->  base component's real lambdas  │
        │   via ctx.renderNode / ctx.renderValue / ctx.dispatchAction   │
        └───────────────▲─────────────────────────────────────────────┘
                        │ renders
        ┌───────────────┴─────────────────────────────────────────────┐
        │ LAYER 1 — Base component  (MatrixGrid.tsx)                    │
        │   plain React, real render-prop lambdas                       │
        │   owns ONLY spatial mechanics (headers, scroll, selection)    │
        │   domain-blind AND cell-blind                                 │
        └───────────────────────────────────────────────────────────────┘
```

- **Layer 1 (base)** is the reusable engine. It knows geometry, not meaning. It
  exposes real React render-prop lambdas so ordinary app code can use it directly.
- **Layer 2 (adapter)** is the only place specs become closures. It exists so the
  DSL / IR can drive the same base component.
- **Layer 3 (preset)** is a thin function that calls the engine with domain
  defaults. `availabilityMatrix(poll)` is ~20 lines that build a `MatrixGrid` node
  with the availability palette, the `yes→ifneedbe→no` cell, and the
  `poll.toggleCell` server action pre-wired.

The payoff: **one engine, many one-line skins.** A RACI chart, a feature-compare
table, and a Doodle poll are all `MatrixGrid` with different specs.

## B.3 The cell contract — the seam that makes cells swappable

For a base grid to host *any* cell (availability toggle, calendar day, rating,
avatar), it must talk to cells through a **stable payload contract**. The grid
owns geometry + the `onAction` plumbing; the cell owns everything visual and
semantic.

```ts
// The contract every cell honors. The grid calls the cell with this:
interface CellRenderPayload<Row = JsonObject, Col = MatrixColumnSpec> {
  row: Row;                 // opaque row datum
  col: Col;                 // opaque column
  value: unknown;           // value at (row,col), already resolved by the accessor
  selected: boolean;
  editable: boolean;
  onAction: (payload: JsonObject) => void;   // grid-owned; cell just calls it
}
```

Any component matching this shape is a valid cell. `AvailabilityCell`,
`DayCell`, `RatingCell`, `RsvpCell` become **independent atoms** — none knows
about the grid, and the grid knows about none of them. **The contract is the
seam.**

There are two ways a caller supplies cells (support both, mirroring existing
patterns):

```yaml
# Mode A — "bring your own renderer" (homogeneous, data-driven)
#   same shape as DataTable columns[].cell today
MatrixGrid:
  cell: CellSpec            # ONE renderer applied per (row,col); DSL swaps it

# Mode B — "bring your own cells" (heterogeneous / prebuilt)
#   same shape as Panel.actions / SlideShell.primary slots (RenderableValue)
MatrixGrid:
  cells: WidgetNode[][]     # an explicit matrix of nodes you built yourself
```

Mode B is the literal "pass in your own list of grid cells" — the grid becomes a
pure spatial container and renders whatever it is handed.

> **The one discipline that keeps this from rotting.** The moment `MatrixGrid`
> grows an `if (kind === "availability")`, the layering is broken — that knowledge
> belongs in the cell atom or the DSL preset, never the engine. Same rule as
> "layout must not know domain nouns", extended to "engines must not know cell
> nouns".

## B.4 Cross-DSL composition (and the naming rule)

Because the renderer is a flat `registry.get(node.type)` over a merged registry,
a single IR tree crosses DSLs freely:

```
CalendarWeekPanel            (calendar.dsl)
└─ slot: eventDetail
   └─ Panel                  (ui.dsl)
      └─ MeetingPollPanel    (schedule.dsl)
         └─ SegmentedBar     (ui.dsl)
```

You assemble the app surface with `mergeWidgetRegistries(uiWidgetRegistry,
timeWidgetRegistry, scheduleWidgetRegistry, calendarWidgetRegistry)` — exactly
how `defaultWidgetRegistry` merges today. To respect the duplicate-`type` gotcha
(A.5), **namespace type strings once more than one DSL specializes the same base
engine**:

```
type: "schedule/AvailabilityMatrix"   // schedule.dsl
type: "calendar/MonthMatrix"          // calendar.dsl
```

Both specialize the same `MatrixGrid` base, under collision-proof names.

---

# PART C — The widgets

Naming: **generic engines/atoms** go in `ui.dsl` / a new `time.dsl` (they are
domain-blind). **Domain presets** go in `schedule.dsl` / `calendar.dsl`.

## C.1 Inventory & generic-vs-domain split

| Domain need | Generic engine (build once) | Module | Reuse beyond scheduling |
|---|---|---|---|
| Doodle poll grid | **`MatrixGrid`** | `data.dsl` | RACI, compare table, rubric, confusion matrix |
| Poll result bar | **`SegmentedBar`** (generalizes `ContextBudgetBar`) | `ui.dsl` | quota, sentiment, budget |
| Month picker/heatmap | **`MonthGrid`** | `time.dsl` | date picker, activity heatmap |
| Week/day calendar | **`TimeGrid`** | `time.dsl` | planner, resource lanes, gantt |
| Availability toggle | **`CycleCell`** (n-state) | `ui.dsl` | RSVP, kanban state, rating |
| "5/8" tally | **`RatioBadge`** | `ui.dsl` | any count/total |
| Overlapping avatars | **`AvatarStack`** | `ui.dsl` | any people list |
| Date marker, slot pill | **`DateTile`, `TimeRangeChip`** | `time.dsl` | already generic |
| Poll / results / booking panels | **stay domain presets** | `schedule.dsl` | — |
| Month/week calendar panels | **thin domain wrappers** | `calendar.dsl` | — |

## C.2 `MatrixGrid` (the flagship engine)

**Purpose.** A rows × columns grid where each cell is a *value at (row, col)*,
rendered by a swappable cell, with optional sticky headers, an aggregate footer
row, selection, and keyboard navigation. Domain-blind.

**ASCII (as a Doodle poll, via the `availabilityMatrix` preset):**

```
              Thu Jul9   Fri Jul10  Fri Jul10  Sat Jul11   <- column headers (slot chips)
              14:00      10:00      16:00      09:00
           ┌──────────┬──────────┬──────────┬──────────┐
 Alice   💬│    ✓     │    ~     │    ✕     │    ✓     │   <- ParticipantRow = row header + cells
 Bob       │    ✓     │    ✓     │    ✕     │    ·     │
 Chen      │    ✓     │    ✓     │    ✓     │    ✓     │
 You    ✎  │   [✓]    │   [~]    │   [✕]    │   [ ]    │   <- editable row (CycleCell)
           ├──────────┼──────────┼──────────┼──────────┤
 tally     │  4 ★     │    3     │    1     │    2     │   <- aggregate footer
           └──────────┴──────────┴──────────┴──────────┘
```

**Base props (Layer 1, real React):**

```ts
interface MatrixGridProps {
  rows: JsonObject[];
  columns: MatrixColumnSpec[];               // { id; header: ReactNode; meta? }
  // Mode A: a renderer; Mode B: an explicit matrix. Provide one.
  renderCell?: (p: CellRenderPayload) => ReactNode;
  cells?: ReactNode[][];
  valueAt?: (row: JsonObject, col: MatrixColumnSpec) => unknown;   // accessor
  renderRowHeader?: (row: JsonObject) => ReactNode;
  footer?: { render: (col: MatrixColumnSpec) => ReactNode };
  selectedCell?: { rowKey: string; colId: string } | null;
  editableRowKey?: string;
  onCell?: (p: { rowKey: string; colId: string; value: unknown }) => void;
  stickyHeader?: boolean;
}
```

**IR props (Layer 2, serializable — note everything is a spec):**

```ts
interface MatrixGridWidgetProps extends BaseWidgetProps {
  rows: JsonObject[];
  columns: Array<{ id: string; header: RenderableValue; meta?: JsonObject }>;
  valueAt: { rowField: string; colField?: string } | { template: string };  // accessor spec
  cell: CellSpec | CycleCellSpec;            // render spec (Mode A)
  cells?: WidgetNode[][];                    // slot matrix (Mode B)
  rowHeader?: CellSpec;
  colorBy?: StyleBySpec;                     // NEW: value -> styleKey -> style
  footer?: { aggregate: "count" | "sum" | "custom"; render: CellSpec };
  editableRowKey?: string;
  selectedCell?: { rowKey: string; colId: string };
  onCellAction?: ActionSpec;                 // fired with { rowKey, colId, value }
  onColumnHeaderAction?: ActionSpec;
}
```

**Adapter pseudocode (Layer 2):**

```tsx
defineWidget<MatrixGridWidgetProps>({
  type: "MatrixGrid", module: "data.dsl",
  render: (props, _children, ctx) => (
    <MatrixGrid
      rows={props.rows}
      columns={props.columns.map(c => ({ ...c, header: ctx.renderValue(c.header) }))}
      valueAt={(row, col) => resolveAccessor(props.valueAt, row, col)}
      renderCell={
        props.cells
          ? undefined                                       // Mode B handled below
          : (p) => renderCellSpec(props.cell, p, ctx)       // Mode A: spec -> node
      }
      cells={props.cells?.map(r => r.map(n => ctx.renderNode(n)))}   // Mode B
      renderRowHeader={props.rowHeader ? (row) => renderCellSpec(props.rowHeader, { row }, ctx) : undefined}
      editableRowKey={props.editableRowKey}
      onCell={
        props.onCellAction
          ? (p) => ctx.dispatchAction(props.onCellAction!, p)  // spec -> callback
          : undefined
      }
    />
  ),
});
```

**`NEW` `StyleBySpec` (the one new spec type):**

```ts
// The defunctionalized "colorFn". Add to ir.ts.
interface StyleBySpec {
  field?: string;                    // which value to key on (defaults to cell value)
  styleSet: ContextStyleSet;         // key -> ContextVisualStyle
  map?: Record<string, string>;      // optional value -> styleKey remap
  fallbackStyleKey?: string;
}
```

Adding `StyleBySpec` once makes `MatrixGrid`, `MonthGrid`, `SegmentedBar`, and the
existing context diagrams all recolorable from data.

## C.3 `CycleCell` (an atom that honors the cell contract)

**Purpose.** A cell that cycles through N states on click (`yes → ifneedbe → no →
unknown`), each state mapped to a `styleKey`. This is the availability toggle,
but also RSVP, kanban state, star rating.

```
┌───────┐ ┌───────┐ ┌───────┐ ┌───────┐
│   ✓   │ │   ~   │ │   ✕   │ │   ·   │
└───────┘ └───────┘ └───────┘ └───────┘
  yes      ifneedbe    no       unknown
```

```ts
interface CycleCellSpec {                    // IR form
  kind: "cycle";
  states: string[];                          // ordered ring of state ids
  glyphs?: Record<string, RenderableValue>;  // state -> glyph
  styleKey?: string;                         // legend key into styleSet
}
interface CycleCellProps {                   // React form (honors CellRenderPayload)
  state: string; states: string[]; styleSet?: ContextStyleSet;
  readOnly?: boolean; onCycle?: (next: string) => void;
}
```

## C.4 `SegmentedBar` (generalizes `ContextBudgetBar`)

**Purpose.** A proportional bar of colored segments with optional counts and
markers. Poll result heat bar, quota, sentiment split.

```
Thu Jul 9 14:00  ▓▓▓▓▓▓▓▒▒░░  6 yes · 1 maybe · 1 no   ★ best
```

```ts
interface SegmentedBarWidgetProps extends BaseWidgetProps {
  segments: Array<{ value: number; styleKey: string; label?: RenderableValue }>;
  styleSet: ContextStyleSet;
  total?: number;
  showCounts?: boolean;
  markers?: Array<{ at: number; styleKey: string; label?: RenderableValue }>;
  onSegmentAction?: ActionSpec;
}
```

## C.5 `MonthGrid` (engine behind month pickers & heatmaps)

**Purpose.** A month of day cells; each day carries a value/markers; the day is
rendered by a swappable spec. Date picker, poll slot picker, GitHub-style
activity heatmap.

```
      July 2026
 Mo Tu We Th Fr Sa Su
     1  2  3  4  5  6
  7  8 [9]10 11 12 13     [9] selected; dots/heat = markers[dateISO]
 14 15 16 17 18 19 20
 21 22 23 24 25 26 27
 28 29 30 31
```

```ts
interface MonthGridWidgetProps extends BaseWidgetProps {
  monthISO: string;
  markers?: Record<string /*dateISO*/, { count: number; styleKey?: string }>;
  dayRender?: CellSpec | { kind: "heat"; styleSet: ContextStyleSet } | { kind: "dots" };
  selectedDateISO?: string;
  minDateISO?: string; maxDateISO?: string;
  onDaySelectAction?: ActionSpec;
  onMonthChangeAction?: ActionSpec;
  dayContent?: WidgetNode;                    // slot
}
```

## C.6 `TimeGrid` (engine behind week/day calendars)

**Purpose.** Hour-ruled day columns with positioned event blocks. Week calendar,
day planner, resource lanes.

```
      Mon      Tue      Wed
 09  ▕█ standup                   block position derived from start/end
 10       ▕███ review
 11  ▕█            ▕████ 1:1
 12            ▕█
```

```ts
interface TimeGridWidgetProps extends BaseWidgetProps {
  daysISO: string[];                          // one column per date
  blocks: Array<{ id: string; dayISO: string; startISO: string; endISO: string;
                  styleKey: string; label: RenderableValue }>;
  styleSet: ContextStyleSet;
  hourStart?: number; hourEnd?: number;
  onBlockSelectAction?: ActionSpec;
  onSlotCreateAction?: ActionSpec;            // drag-to-create -> server action
}
```

## C.7 Supporting atoms

- **`RatioBadge`** — `{ count, total, kind? }` → "5/8" with tone. Thin `Caption`
  wrapper.
- **`AvatarStack`** — `{ people: {name,avatarUrl?}[], max? }` → overlapping avatars
  + overflow count.
- **`DateTile`** — `{ dateISO, emphasis? }` → torn-calendar-page tile (month over
  big day number; `--rag-font-role-metric` + `--rag-font-role-label`).
- **`TimeRangeChip`** — `{ startISO, endISO, tz, tone? }` → "Thu Jul 9 · 14:00–15:00".

## C.8 Domain presets (Layer 3)

These are thin builders in `schedule.dsl` / `calendar.dsl`. They do **not** add
new base components — they configure the engines above.

```
availabilityMatrix(poll)   -> MatrixGrid  (cell=CycleCell(yes/ifneedbe/no),
                                            palette=availability, onCell=poll.toggleCell)
pollResults(poll, tallies) -> Stack of SegmentedBar (one per slot) + finalize buttons
bookingPage(bookingType, days, slots) -> SplitPane[ MonthGrid | TileGrid<TimeRangeChip> ]
monthCalendar(events)      -> MonthGrid   (dayRender=event pills, colorBy=event.colorKey)
weekCalendar(events)       -> TimeGrid    (blocks from events)
```

And the panels that carry Doodle vocabulary (`MeetingPollPanel`,
`PollResultsPanel`, `PollBuilderPanel`, `BookingPagePanel`, `EventDetailDialog`,
`SchedulingStudioShell`) are organisms composed from the presets + existing
layout (`SidebarShell`, `Panel`, `SplitPane`, `FormPanel`, `DialogShell`). See the
prior design discussion for their full ASCII; they are Layer-3 compositions and
should be built *after* the engines are stable.

## C.9 The scheduling domain model (DTOs)

A new `src/scheduling/types.ts` (mirrors `src/context/types.ts`) owns the DTOs.
Abbreviated:

```ts
type AvailabilityState = "yes" | "ifneedbe" | "no" | "unknown";
interface TimeSlot   { id; startISO; endISO; tz; allDay?; label?; }
interface PollOption { id; slot: TimeSlot; note?; }
interface ParticipantResponse { id; name; avatarUrl?; comment?;
                                cells: Record<string /*optionId*/, AvailabilityState>; }
interface MeetingPoll { id; title; description?; location?; organizer;
                        options: PollOption[]; responses: ParticipantResponse[];
                        settings: PollSettings; status: "open"|"finalized"|"closed"; }
interface SlotTally  { optionId; yes; ifneedbe; no; score; isBest?; atCapacity?; }  // server-computed
interface CalendarEvent { id; title; startISO; endISO; allDay?; colorKey; location?; attendees?; }
// + BookingType / BookableDay / BookableSlot for the 1:1 flow
```

Plus `src/scheduling/fixtures.ts` with sample data for Storybook (again mirroring
`src/context/fixtures.ts`).

---

# PART D — The backend (Doodle-like server)

The package stays presentational; the **`packages/web`** app owns data. Two halves.

## D.1 Widget server actions

Interactive widgets fire `{ kind: "server", name, payload }` actions.
`dispatchWidgetAction` (`src/widgets/actions.ts`) POSTs to
`/api/widget/actions/<name>` and expects `{ ok, refresh?, toast?, patch?, data? }`.
If `refresh` is true it dispatches a `popstate` to re-fetch.

```
Widget (onCellAction)  --POST /api/widget/actions/poll.toggleCell-->  server
                       <--   { ok, patch: { tallies } }            ---
```

| Action `name` | Payload | Returns |
|---|---|---|
| `poll.toggleCell` | `{ pollId, responseId, optionId, state }` | `{ ok, patch: { tallies } }` |
| `poll.submitResponse` | `{ pollId, name, comment, cells }` | `{ ok, refresh, toast }` |
| `poll.create` / `poll.update` | draft `MeetingPoll` | `{ ok, data: { id }, refresh }` |
| `poll.addSlot` / `poll.removeSlot` | `{ pollId, slot | optionId }` | `{ ok, patch: { options } }` |
| `poll.finalize` | `{ pollId, optionId }` | `{ ok, refresh, toast }` |
| `poll.remind` | `{ pollId }` | `{ ok, toast }` |
| `booking.daySlots` | `{ bookingTypeId, dateISO, tz }` | `{ ok, patch: { slots } }` |
| `booking.book` | `{ slotId, name, email }` | `{ ok, data: { bookingId }, refresh }` |
| `calendar.range` | `{ startISO, endISO, view }` | `{ ok, patch: { events } }` |

## D.2 REST resources (app-owned, RTK Query)

```
POST   /api/polls                       create
GET    /api/polls/:id                   poll + responses + computed tallies
PUT    /api/polls/:id                   edit
POST   /api/polls/:id/responses         participant submit
PATCH  /api/polls/:id/responses/:rid    cell edits (optimistic)
POST   /api/polls/:id/finalize          { optionId }
GET    /api/booking-types/:id/availability?date=&tz=
POST   /api/bookings                    { slotId, name, email }
GET    /api/calendar?start=&end=&view=
```

## D.3 What the server computes (so widgets stay dumb)

- **Tally + ranking** (`SlotTally.yes/ifneedbe/no/score/isBest`) — a `GROUP BY
  option_id, state` aggregate. Widgets only *render* it.
- **Timezone normalization** — store UTC (`startISO`); the server resolves per-day
  slot lists so DST is never a frontend concern.
- **Hidden polls / capacity / deadline** — the server strips other votes when
  `hideVotesUntilResponded`, and rejects over-capacity / late submits.
- **Notifications** — invites, reminders, finalize, booking confirmations are
  triggered by the corresponding server actions.

Data model sketch: `polls`, `poll_options`, `responses`, `response_cells`
(participant × option × state), `booking_types`, `availability_rules`,
`bookings`, and a `calendar_events` view.

---

# PART E — Build order & checklists

## E.1 Recommended order

Build **engines bottom-up**, prove each in Storybook before the next:

1. **Atoms first** (fast wins, no dependencies): `DateTile`, `RatioBadge`,
   `TimeRangeChip`, `AvatarStack`.
2. **`CycleCell`** — the first cell-contract citizen; proves the payload shape.
3. **`SegmentedBar`** — self-contained; also validates the palette reuse.
4. **`MatrixGrid`** base — the flagship; wire `CycleCell` as its cell in a story.
5. *(later)* `MonthGrid`, `TimeGrid`.
6. *(later)* `StyleBySpec` in `ir.ts` + adapters + registries.
7. *(later)* domain DTOs, presets, and organism panels.

> **This ticket's first coding milestone stops after step 4** (a small set of base
> widgets) so the reviewer can check looks in Storybook before we go wide.

## E.2 Per-widget file checklist (do all six)

```
[ ] X.tsx            React component, props typed, data-rag-* attribute present
[ ] X.module.css     local anatomy only; tokens (--rag-*), no raw colors/fonts
[ ] X.stories.tsx    states: default, empty, dense/overflow, selected, disabled,
                     error/warning, alternate direction (where applicable)
[ ] X.widget.tsx     defineWidget adapter (add AFTER React API is stable)
[ ] X.widget.yaml    manifest: type, module, helper, props, slots, actions, status
[ ] index.ts         re-exports
```

Storybook title prefixes (from GUIDELINES):

```
Design System/Atoms/<Atom>            e.g. CycleCell, DateTile, RatioBadge
Component Library/Molecules/<Comp>    e.g. SegmentedBar, MatrixGrid
```

## E.3 Definition of done for the first milestone

- Atoms + `CycleCell` + `SegmentedBar` + `MatrixGrid` (base) render in Storybook.
- `MatrixGrid` demonstrates **both** injection modes: a story using `renderCell`
  (Mode A) and a story passing an explicit `cells` matrix (Mode B), with
  `CycleCell` as the cell in at least one.
- Typecheck passes:
  `pnpm --dir packages/rag-evaluation-site typecheck`.
- No IR adapters required yet (React-first). Adapters/manifests land in a later
  milestone once looks are approved.

---

## API Reference (quick lookup)

| Symbol | File | Role |
|---|---|---|
| `WidgetNode`, `ComponentNode` | `src/widgets/ir.ts` | IR node model |
| `component()`, `element()`, `text()` | `src/widgets/ir.ts` | IR constructors |
| `ActionSpec`, `dispatchWidgetAction` | `src/widgets/actions.ts` | defunctionalized event handlers |
| `CellSpec`, `renderCell` | `src/widgets/cellRenderers.tsx` | defunctionalized render lambda |
| `RenderContext` | `src/widgets/registry.ts` | `renderNode/renderValue/dispatchAction` bridge |
| `defineWidget`, `createWidgetRegistry`, `mergeWidgetRegistries` | `src/widgets/registry.ts` | adapters & registries |
| `WidgetRenderer` | `src/widgets/WidgetRenderer.tsx` | JSON tree → React |
| `defaultWidgetRegistry` + per-module registries | `src/widgets/defaultRegistry.ts` | where adapters are registered |
| `ContextStyleSet`, `ContextVisualStyle` | `src/context/types.ts` | palette contract to reuse |
| `ContextBudgetBar` | `src/components/molecules/ContextBudgetBar/` | worked example `SegmentedBar` generalizes |
| `DataTable` + `.widget.tsx` | `src/components/molecules/DataTable/` | worked example of `CellSpec` in a grid |
| `StepList` folder | `src/components/molecules/StepList/` | canonical six-file layout |

## File Reference (what to create)

```
src/scheduling/types.ts                          NEW  domain DTOs
src/scheduling/fixtures.ts                        NEW  Storybook sample data
src/components/atoms/DateTile/*                    NEW  atom (6 files)
src/components/atoms/RatioBadge/*                  NEW  atom
src/components/atoms/TimeRangeChip/*               NEW  atom
src/components/atoms/AvatarStack/*                 NEW  atom
src/components/atoms/CycleCell/*                    NEW  cell-contract atom
src/components/molecules/SegmentedBar/*            NEW  engine
src/components/molecules/MatrixGrid/*              NEW  flagship engine
src/components/molecules/MonthGrid/*              (later) engine
src/components/molecules/TimeGrid/*               (later) engine
src/widgets/ir.ts                                 EDIT add StyleBySpec, *WidgetProps, RagWidgetType entries
src/widgets/cellRenderers.tsx                     EDIT add CycleCellSpec handling
src/widgets/defaultRegistry.ts                    EDIT new scheduleWidgetRegistry/timeWidgetRegistry
```

## Open Questions

- Priority ordering of the three product flows (group poll / 1:1 booking / full
  calendar)? Milestone 1 is flow-agnostic (engines), so this can be answered
  later, but it decides which presets we build first.
- Adopt namespaced `type` strings (`schedule/AvailabilityMatrix`) from day one, or
  only when a second DSL specializes the same engine?
- Does `MatrixGrid` need virtualized scrolling for large polls (>50 participants),
  or is CSS sticky + overflow sufficient for v1? (Recommend: sufficient for v1.)

## References

- `packages/rag-evaluation-site/GUIDELINES.md` — design-system rules.
- `src/widgets/ir.ts`, `actions.ts`, `cellRenderers.tsx`, `registry.ts`,
  `WidgetRenderer.tsx`, `defaultRegistry.ts` — the IR mechanism.
- `src/context/types.ts` — the palette + DTO precedent to mirror.
- Sibling ticket docs under `ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY...` and
  `...RAGEVAL-CONTEXT-WINDOWS-DESIGN...` (design-system unification + context
  viewer integration) for prior art on this layering.
