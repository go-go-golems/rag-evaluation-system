---
Title: 'The Widget System: A New Intern''s Guide to the Architecture'
Ticket: RAGEVAL-WIDGET-DECOMPOSITION
Status: active
Topics:
    - design-system
    - widget-ir
    - ui-dsl
    - react
    - frontend-architecture
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-06T21:43:41.382538663-04:00
WhatFor: ""
WhenToUse: ""
---

# The Widget System: A New Intern's Guide to the Architecture

This guide teaches you how the RAG-Evaluation widget system works. Not the API —
the architecture: what the pieces are, why they exist, and how a description of a
user interface travels from a line of JavaScript on a server, through a tree of
JSON, into React components on a screen. By the end you will be able to open any
widget in the codebase and understand what it is doing and why, and you will
recognize the one design pattern the whole library is organized around.

Read the chapters in order. Each one builds the vocabulary the next one assumes.
The code excerpts are real; the file paths are real; you can open every one of
them as you read.

---

## Chapter 1 — Three representations of one user interface

Most UI code has a single representation: React components. You write JSX, React
renders it, and that is the whole story. This codebase is different, and the
difference is the first thing you need to understand, because everything else is
a consequence of it.

Here a user interface exists in **three** forms at different moments in its life:

1. **A React component tree** — ordinary `.tsx` files under
   `packages/rag-evaluation-site/src/components/`. This is what the browser paints.
2. **A Widget IR tree** — a tree of plain JSON objects that *describes* a UI
   without being one. "IR" means intermediate representation. It lives in
   `packages/rag-evaluation-site/src/widgets/`.
3. **A DSL call** — a line of JavaScript, run server-side inside a Go program, that
   *produces* the JSON tree. This lives in `pkg/widgetdsl/` and is written in Go.

The reason there are three forms is that the people who need to *describe* a UI
are not always running a browser. A Go service, or a script an analyst writes,
needs to emit a dashboard or a report and hand it to a frontend to render. You
cannot send a React component over the network — a React component is a function,
and functions do not serialize. You *can* send JSON. So the system introduces a
JSON representation of a UI (form 2), a way to produce it from a script (form 3),
and a single component that turns it back into React (form 1).

```
   pkg/widgetdsl (Go/Goja)          Widget IR (JSON)              React (browser)
   ─────────────────────            ────────────────              ───────────────
   ui.panel({title:"x"},    ──►     { kind:"component",    ──►    <Panel title="x">
     ui.text("hi"))                   type:"Panel",                <>hi</>
                                      props:{title:"x"},         </Panel>
                                      children:[
                                        {kind:"text",text:"hi"}
                                      ] }
        form 3                            form 2                       form 1
      (produces)                     (serializable)                 (rendered)
```

The middle form is the pivot. It is the contract between "something that wants to
describe a UI" and "something that can draw one." Chapters 2 and 3 explain how
form 2 becomes form 1. Chapter 6 explains how form 3 produces form 2. Everything
in between is detail hanging off this spine.

One consequence is worth stating now, because it will keep coming up. Anything
that must cross the JSON boundary has to *be* JSON. A color function, a click
handler, a "how do I render this cell" rule — none of these can be a JavaScript
closure, because closures do not survive `JSON.stringify`. The system's answer to
this constraint is the single most important idea in the codebase, and Chapter 4
is devoted to it. Hold the question in mind until then: *how do you send a
function through JSON?*

**Key points:**

- A UI here exists as React (form 1), as a JSON IR tree (form 2), and as a
  server-side DSL call that produces the IR (form 3).
- The IR exists because a UI often needs to be *described* by something that is
  not a browser, and JSON is what can be sent over a network; React functions
  cannot.
- The IR tree is the contract between producer and renderer. Learn it first.
- Nothing that crosses the JSON boundary can be a function. Remember the question:
  how do you send a function through JSON?

---

## Chapter 2 — The IR tree and the loop that renders it

The Widget IR is a tree of three kinds of node. The definitions are in
`src/widgets/ir/core.ts`, and they are small enough to hold in your head:

```ts
type WidgetNode = TextNode | ElementNode | ComponentNode;

interface TextNode      { kind: "text";      text: string; }
interface ElementNode   { kind: "element";   tag: string; attrs?; children?: WidgetNode[]; }
interface ComponentNode { kind: "component"; type: string; props?; children?: WidgetNode[]; }
```

A `TextNode` is a string. An `ElementNode` is a raw HTML tag — a `div`, a `span`
— with attributes. A `ComponentNode` is the interesting one: it names a widget by
a string `type` (`"Panel"`, `"DataTable"`, `"MatrixGrid"`) and carries a bag of
`props`. The `type` is a string, not a class or a function reference, precisely
because this object has to be JSON.

The component that turns this tree into React is `WidgetRenderer`
(`src/widgets/WidgetRenderer.tsx`). It is a recursive tree walk, and its core is a
switch on `node.kind`:

```tsx
function renderWidgetNode(node, ctx, registry) {
  switch (node.kind) {
    case "text":      return node.text;
    case "element":   return createElement(node.tag, node.attrs, renderChildren(...));
    case "component": return renderComponentNode(node, ctx, registry);
  }
}
```

Text renders to a string. An element renders to a React element with the same tag.
A component node needs one more step, and that step is the heart of the system:

```tsx
function renderComponentNode(node, ctx, registry) {
  const adapter = registry.get(node.type);          // look up "Panel" -> its adapter
  if (!adapter) return <UnknownWidget node={node} />;
  return adapter.render(node.props ?? {}, renderChildren(node.children), ctx, node);
}
```

Read that first line carefully, because it explains a property you will rely on
constantly. To render a `ComponentNode`, the renderer takes its `type` string and
looks it up in a **flat registry** — a `Map` from `"Panel"` to the code that
renders a Panel. The renderer does not know what a Panel is. It knows how to look
one up. This is why any widget can contain any other widget: nesting is just a
node whose `children` contain more nodes, and each node is resolved independently
by the same `registry.get`. There is no coupling between a parent widget type and
its children's types. A calendar node can contain a poll node can contain a
button node, and the renderer treats them identically.

Consider a concrete tree and trace it by hand:

```json
{ "kind": "component", "type": "Panel", "props": { "title": "Results" },
  "children": [
    { "kind": "component", "type": "Button", "children": [ { "kind": "text", "text": "Finalize" } ] }
  ] }
```

The renderer sees `kind:"component"`, `type:"Panel"`. It calls
`registry.get("Panel")`, gets the Panel adapter, and calls its `render` — passing
`{title:"Results"}` as props and the rendered children. To render those children,
it recurses: the child is `kind:"component"`, `type:"Button"`, so
`registry.get("Button")`, and *its* child is a `TextNode` that renders to the
string `"Finalize"`. The result is `<Panel title="Results"><Button>Finalize</Button></Panel>`.
The tree walk bottoms out at text and builds React back up on the way out.

**Key points:**

- The IR has three node kinds: text, element (raw HTML tag), and component (a
  widget named by a `type` string).
- `type` is a string because the node must be JSON. The renderer resolves it
  through a flat `registry.get(type)` lookup.
- Because resolution is per-node and flat, any widget can nest any other. There is
  no parent/child type coupling. This is the property that makes composition free.
- Rendering is a recursive tree walk that bottoms out at text nodes and rebuilds
  React on the way back up.

---

## Chapter 3 — Adapters, the registry, and the render context

Chapter 2 left a gap: `registry.get("Panel")` returns "the code that renders a
Panel," but what is that code, exactly? It is an **adapter**. An adapter is the
bridge between a JSON `ComponentNode` and a real React component. Every widget has
one, in a file named `X.widget.tsx`.

An adapter is created with `defineWidget` (`src/widgets/registry.ts`), and its
shape is three fields: a `type` string, a `module` label, and a `render`
function. Here is a real one, for the `StepList` molecule
(`src/components/molecules/StepList/StepList.widget.tsx`), lightly trimmed:

```tsx
export const stepListWidget = defineWidget<StepListWidgetProps>({
  type: "StepList",
  module: "ui.dsl",
  render: (props, _children, ctx) => (
    <StepList
      items={props.items.map((item) => ({
        ...item,
        title: ctx.renderValue(item.title),        // JSON value -> ReactNode
        description: ctx.renderValue(item.description),
      }))}
      onItemSelect={
        props.onItemSelectAction
          ? (itemId) => ctx.dispatchAction(props.onItemSelectAction!, { itemId })
          : undefined                               // ActionSpec -> real callback
      }
    />
  ),
});
```

The adapter's job is translation. On the way in, `props` is plain JSON: `item.title`
might be a string, or it might itself be a `WidgetNode` (a nested subtree). The
adapter cannot assume; it calls `ctx.renderValue(item.title)`, which renders a
node if it is one and passes a string through if it is not. On the way out, the
React `StepList` wants a real `onItemSelect` callback, but the JSON only carried an
`onItemSelectAction` — a *description* of what to do. The adapter converts that
description into a live function with `ctx.dispatchAction`. We will unpack
`ActionSpec` and `dispatchAction` in Chapter 4; for now notice the pattern: **the
adapter is the one place where JSON specs become live React behavior.**

The `ctx` parameter is a `RenderContext`, and it is the only capability an adapter
is given. Its full surface is five methods (`src/widgets/registry.ts`):

| Method | What it does |
|---|---|
| `renderNode(node)` | render one child `WidgetNode` to React |
| `renderChildren(children)` | render an array of child nodes |
| `renderValue(value)` | render a `RenderableValue` — a node if it is one, else the string |
| `bindAction(action, ctx)` | return a `() => void` that runs an `ActionSpec` |
| `dispatchAction(action, ctx)` | run an `ActionSpec` now |

Notice what is *not* in that table: the adapter never receives the registry, and
never receives the app's global state or router. It cannot reach outside its own
render. This is a deliberate boundary. An adapter can render its children and turn
its specs into behavior, and nothing more. That restriction is what keeps the
widget package presentational and reusable — the same adapter works in the real
app, in Storybook, and in a test, because it depends only on the five methods it
is handed.

Adapters are collected into registries and merged. In
`src/widgets/defaultRegistry.ts` you will find per-module registries —
`uiWidgetRegistry`, `dataWidgetRegistry`, `contextWindowWidgetRegistry`, and so on
— each built with `createWidgetRegistry([...])` and then combined with
`mergeWidgetRegistries(...)` into the `defaultWidgetRegistry` that the renderer
uses. The merge is a flat union keyed by `type`. One rule follows from this and it
bites people: `createWidgetRegistry` throws if two adapters claim the same `type`
string. The `type` namespace is global across every module. Two DSLs cannot both
register `"Grid"`.

**Key points:**

- An adapter (`X.widget.tsx`, made with `defineWidget`) translates a JSON
  `ComponentNode` into a real React component.
- The adapter is the single place where JSON becomes behavior: `ctx.renderValue`
  turns value-or-node into React, `ctx.dispatchAction` turns an `ActionSpec` into a
  callback.
- An adapter's only capability is the five-method `RenderContext`. It cannot see
  the registry, the store, or the router. That boundary is what keeps widgets
  reusable.
- Registries merge flatly by `type`; the `type` namespace is global, and duplicate
  types are a hard error.

---

## Chapter 4 — How to send a function through JSON

Chapter 1 posed the question and Chapter 3 showed it in passing: the JSON carried
`onItemSelectAction`, not an `onClick` function. This chapter answers the question
directly, because the technique it uses — **defunctionalization** — is the idea the
entire IR is built on, and once you see it you will see it everywhere.

You cannot serialize a closure. So instead of sending the function, you send a
*data description of what the function should do*, and you keep the actual doing in
one interpreter on the other side. A click handler is not a function in the IR; it
is an `ActionSpec`, a tagged union (`src/widgets/ir/actions.ts`):

```ts
type ActionSpec =
  | { kind: "navigate"; to: string }
  | { kind: "download"; to: string }
  | { kind: "server";   name: string; payload? }   // POST /api/widget/actions/<name>
  | { kind: "event";    event: string; detail? }    // browser CustomEvent
  | { kind: "copy";     value?; field? };            // clipboard
```

The interpreter is `dispatchWidgetAction` (`src/widgets/actions.ts`), a switch on
`action.kind`. When the adapter wired `onItemSelect` to `ctx.dispatchAction(action, ...)`,
it was deferring to this switch:

```ts
function dispatchWidgetAction(action, context, onAction) {
  if (action.kind === "copy")     { navigator.clipboard?.writeText(...); return; }
  if (action.kind === "navigate") { window.history.pushState({}, "", action.to); ... }
  if (action.kind === "server")   { fetch(`/api/widget/actions/${action.name}`, ...); }
  // event, download ...
}
```

The value of this indirection is that the set of things a UI can *do* is now a
closed, inspectable vocabulary. There are five kinds of action, defined in one
place, interpreted in one place. A script author on the Go side cannot invent a
new kind of side effect; they can only choose from the five and fill in the data.
That is a strong guarantee, and it is only possible because the action is data
rather than code.

The same move is made three times in the codebase, for the three "function slots"
a UI has:

| The function you want | Reified as data | Defined in | Interpreted by |
|---|---|---|---|
| a click handler | `ActionSpec` | `ir/actions.ts` | `actions.ts:dispatchWidgetAction` |
| "how do I render this table cell?" | `CellSpec` | `ir/cells.ts` | `cellRenderers.tsx:renderCell` |
| "what color is this value?" | `StyleBySpec` | `ir/engines.ts` | `styleBy.ts:resolveStyleByVars` |

`CellSpec` is worth seeing because it makes the pattern concrete for rendering,
not just for events. A table column needs to know how to draw each cell. In
ordinary React you would pass a function `(row) => <StatusText.../>`. In the IR you
pass a `CellSpec`:

```ts
type CellSpec =
  | { kind: "field";  field: string }               // print row[field]
  | { kind: "status"; field: string }               // render <StatusText> of row[field]
  | { kind: "template"; template: string }           // interpolate "${a}-${b}"
  | { kind: "actionButton"; label; action: ActionSpec }
  // ... link, number, caption, constant
```

and `renderCell(spec, row, ...)` (`src/widgets/cellRenderers.tsx`) is the switch
that turns it into a React node. Same shape as actions: a closed union, one
interpreter.

There is one more mechanism that is not a tagged union but serves the same goal:
the **slot**. Any prop typed `RenderableValue` (`WidgetNode | string | number |
boolean | null`) is a hole into which the producer can drop an entire subtree. When
an adapter calls `ctx.renderValue(prop)`, it is honoring a slot: if the prop is a
node, render it; otherwise treat it as text. Slots are how a widget accepts
arbitrary caller-supplied content without knowing what it is — the JSON equivalent
of `props.children`.

**Key points:**

- You send a function through JSON by *not* sending the function: you send a data
  description (a tagged union) and keep one interpreter for it. This is
  defunctionalization.
- The three function slots of a UI are all reified this way: `ActionSpec` (events),
  `CellSpec` (cell rendering), `StyleBySpec` (color).
- Because each is a closed union with a single interpreter, the set of things a UI
  can do or render is inspectable and bounded — a property you get *only* by making
  the function into data.
- A `RenderableValue` slot is the fourth mechanism: a prop that can hold a whole
  subtree, rendered via `ctx.renderValue`.

---

## Chapter 5 — The pattern the whole library is organized around

You now have the machinery: nodes, the renderer, adapters, and specs. This chapter
assembles them into the one pattern you should carry with you into every file. It
has three parts — an **engine**, a **contract**, and a **preset** — and the
clearest place to see it is `MatrixGrid`, the generic grid that backs the
scheduling widgets.

Start with the problem the pattern solves. A "Doodle" availability poll is a grid:
people down the side, time slots across the top, and in each cell a yes/no/maybe
toggle. A month calendar is also a grid: days in a 7-column layout, an event chip
in each cell. A feature-comparison table is a grid too. Written naively, these are
three separate components, each re-implementing sticky headers, scrolling,
selection, and keyboard navigation. The pattern's claim is that they are one
engine with three different cells.

The **engine** is `MatrixGrid` (`src/components/molecules/MatrixGrid/MatrixGrid.tsx`).
It knows exactly one thing: how to lay out rows and columns and put *something* in
each cell. It does not know what a cell means. It does not know about votes, or
events, or features. Its knowledge is purely spatial.

The **contract** is how the engine talks to a cell without knowing what the cell
is. `MatrixGrid` hands every cell a fixed payload, `MatrixCellPayload`
(`MatrixGrid.tsx:18`):

```ts
interface MatrixCellPayload {
  row; col; value;                 // which datum, resolved for you
  rowKey; rowIndex; colIndex;
  selected; editable;
  onAction: (extra?) => void;      // tell the grid something happened
}
```

Any component that accepts this payload is a valid cell. The engine computes
*where* the cell goes and *what value* belongs there; the cell decides how that
value looks and what happens on interaction, reporting back through `onAction`.
This is the same seam as `MatrixCellPayload`'s cousins elsewhere —
`MonthGridDayPayload` for calendar days, `TimeGridBlockPayload` for calendar
events — and it is the exact analogue of the `SegmentPayload` that the context
diagrams *should* use but do not yet (Chapter 8).

The **preset** is a thin function that configures the engine with domain meaning.
`availabilityMatrix(poll)` (`src/widgets/presets/scheduling.ts`) is one: it takes a
`MeetingPoll`, and returns a `MatrixGrid` IR node whose rows are the respondents,
whose columns are the slots, whose cell spec is a `cycle` cell over
`yes/ifneedbe/no`, and whose `onCellAction` is a `server` action named
`poll.toggleCell`. The preset is where the words "vote" and "slot" live. The engine
never hears them.

Put the three together and trace an availability poll from data to pixels:

```
availabilityMatrix(poll)                         [preset: knows "poll", "slot", "vote"]
   returns  { type:"MatrixGrid", props:{ rows, columns, cell:{kind:"cycle",...},
                                          onCellAction:{kind:"server", name:"poll.toggleCell"} } }
        │
        ▼  WidgetRenderer -> registry.get("MatrixGrid") -> MatrixGrid.widget.tsx
   adapter builds renderCell = (payload) => <CycleCell value={payload.value}
                                              onCycle={next => payload.onAction({value:next})}/>
        │
        ▼  <MatrixGrid> lays out rows×cols            [engine: knows only geometry]
   for each (row,col): renderCell(payload)  ───►  <CycleCell/>   [cell: knows "yes/no/maybe"]
        │
        ▼  user clicks a cell -> onCycle -> payload.onAction({value:"no"})
   -> MatrixGrid onCell -> ctx.dispatchAction(poll.toggleCell, {rowKey, colId, value:"no"})
        │
        ▼  POST /api/widget/actions/poll.toggleCell   [back to the ActionSpec interpreter, Ch.4]
```

Every idea from Chapters 2–4 appears in that trace: the preset emits an IR node
(Ch.1–2), the renderer resolves it by `type` (Ch.2), the adapter turns specs into
behavior (Ch.3), the cell spec and the server action are defunctionalized specs
(Ch.4). The pattern is not new machinery; it is the disciplined *use* of the
machinery you already understand.

Why go to this trouble instead of writing three grid components? Because the
engine is where the hard, reusable work lives — layout, selection, keyboard
handling — and you want to write and debug it once. The cell is small and specific.
The preset is a dozen lines. When a fourth grid-shaped UI appears, you write a new
preset and a new cell, and you inherit a correct engine for free. The cost of the
pattern is one indirection — the payload contract — and that indirection is the
whole point: it is the line between "reusable" and "domain-specific," drawn
explicitly.

**Key points:**

- The pattern is **engine + contract + preset**. The engine owns geometry and
  nothing else; the preset owns domain meaning; the contract (a cell payload) is
  the seam that lets the engine place a cell it knows nothing about.
- `MatrixGrid` + `MatrixCellPayload` + `availabilityMatrix` is the reference
  instance. `MonthGrid`, `TimeGrid`, and `SegmentedBar` follow the same shape.
- The pattern reuses the Chapters 2–4 machinery; it does not add new machinery. It
  is a discipline for *where the knowledge lives*.
- The payoff is that the expensive engine is written once and every new
  domain UI is a small preset plus a small cell.

---

## Chapter 6 — The other side of the JSON: the Go/Goja DSL

Everything so far has been about turning IR JSON into React (form 2 → form 1).
This chapter is about producing the JSON in the first place (form 3 → form 2),
which happens in Go, in `pkg/widgetdsl/`. You need to understand it because the IR
has two authors — TypeScript preset functions like `availabilityMatrix`, and
server-side scripts — and the Go side is how the second author writes.

The Go program embeds a JavaScript engine (Goja). A script running in that engine
calls `require("ui.dsl")` and gets back an object of helper functions —
`panel`, `text`, `button`, and so on. Calling `ui.panel({title:"x"})` returns a
plain `map[string]any` that is exactly an IR component node. The DSL is, in one
sentence, a set of JavaScript functions that build IR JSON.

The mechanism is smaller than you would expect. A module is described by a
`moduleSpec` (`pkg/widgetdsl/module.go`), and the interesting field is a map:

```go
var uiHelpers = map[string]string{     // JS helper name -> IR component type
    "panel":  "Panel",
    "button": "Button",
    "stepList": "StepList",
    // ...
}
```

When the module loads, `runtime.install` (`module.go:236`) walks that map and, for
each entry, installs a JavaScript function built by `componentFactory`
(`module.go:565`):

```go
func (r *runtime) componentFactory(componentType string) func(goja.FunctionCall) goja.Value {
    return func(call goja.FunctionCall) goja.Value {
        props, childStart := propsAndChildStart(call.Arguments, 0)
        return r.vm.ToValue(r.buildComponent(componentType, props, call.Arguments[childStart:]))
    }
}
```

`buildComponent` emits `{kind:"component", type:componentType, props, children}`.
That is the entire classic path. Adding a new widget to the DSL is, for the common
case, adding **one line** to a helper map — the string `"stepList": "StepList"` is
enough to make `ui.stepList({...})` produce a correct `StepList` node, because the
factory is generic over the type.

Two richer surfaces sit alongside the plain helpers. The first is builders for the
specs from Chapter 4: `cell.field(...)`, `cell.status(...)`, `action.server(...)`
(`cellObject` at `module.go:272`, `actionObject` at `module.go:318`). These exist
because a `CellSpec` or `ActionSpec` is a small tagged object, and a builder
function that returns `{kind:"status", field:"x"}` is friendlier than writing the
literal. The second is **recipes** (`recipesObject`, `module.go:495`): composite
builders that emit a whole subtree of several components at once —
`masterDetailTableRecipe` (`module.go:853`) returns a Panel containing a DataTable
plus a detail pane. A recipe is the Go-side equivalent of a TypeScript preset:
where a plain helper maps one call to one node, a recipe maps one call to a
configured *tree*.

The relationship to the TypeScript side is a parallel, and it is where drift
lives. `availabilityMatrix` is a TypeScript preset; the equivalent Go recipe does
not exist yet, so a script cannot yet build an availability poll except through the
generic `component("MatrixGrid", props)` escape hatch. The DSL also generates
TypeScript type declarations for its authors (`typescript.go`) and checks them
against the real compiler in a test (`typescript_fixture_test.go`) — so the DSL's
promises to script authors are verified, even though the promises are currently
looser than the TypeScript IR they target.

**Key points:**

- The DSL (`pkg/widgetdsl/`) is a set of JavaScript functions, implemented in Go,
  that build IR JSON. A script calls them; the frontend renders what they return.
- The classic path is a `map[name]→type` plus a generic `componentFactory`. Adding
  a widget to the DSL is usually one map entry.
- `cell.*` / `action.*` builders construct the Chapter 4 specs; **recipes** are
  composite builders that emit a configured subtree — the Go analogue of a
  TypeScript preset.
- The Go and TypeScript authors target the same IR but are maintained separately,
  which is where they drift (Chapter 7).

---

## Chapter 7 — The manifests, and why the same fact lives in five places

There is one more file per widget you have not met: `X.widget.yaml`, the
**manifest**. Open `MatrixGrid.widget.yaml` and you will see the widget described
as data:

```yaml
type: MatrixGrid
module: data.dsl
helper: matrixGrid
props: MatrixGridWidgetProps
reactComponent: MatrixGrid
slots: [columns.header, cornerCell, cell.cycle.glyphs, cells]
actions: [onCellAction]
```

Every field here is a fact the rest of the system also needs. The `type` appears
in the `RagWidgetType` union in `ir/core.ts`. The `props` type appears in the
`WidgetProps` union in `ir/props.ts`. The `module`/`helper` pair is what the Go DSL
must expose. The `slots` are exactly the props an adapter must pass through
`ctx.renderValue`. The `actions` are exactly the props an adapter must wire through
`ctx.dispatchAction`. The manifest is, in effect, a machine-readable specification
of the adapter and its registrations.

Here is the problem, and it is worth understanding because it is the motivation for
much of the companion analysis document. **Nothing reads the manifest.** A search
of the codebase finds no TypeScript that consumes `.widget.yaml` to generate
anything. The `RagWidgetType` union is hand-written. The `WidgetProps` union is
hand-written. The registry import list is hand-written. The Go helper maps are
hand-written. The manifest sits alongside all of them, describing the same facts,
read by nothing but a linter.

When the same fact is written by hand in several places, the copies drift. And they
have. There are 81 adapters but only 79 manifests — `ContextStyleSwatch` and
`RichArticle` have adapters with no manifest. The scheduling manifests name a
`time.dsl` module and helpers like `matrixGrid` that **do not exist** in the Go
code, because the manifest was written to describe an intention that the Go side
never implemented, and no check caught the gap. The manifest says `helper:
matrixGrid`; the Go `moduleSpec` has never heard of it; nothing compared the two.

This is not a bug in any one file. It is the predictable consequence of a design
where a single fact has five hand-maintained homes and no source of truth. The
companion document (`design-doc/01-...`) proposes the fix — make the manifest the
generator input and derive the unions, the registry, and the Go maps from it — but
you do not need that proposal to learn the lesson. The lesson is a way of reading
the codebase: when you see the same widget named in `core.ts`, in `props.ts`, in
`defaultRegistry.ts`, in `X.widget.yaml`, and in `module.go`, you are looking at
five copies of one fact, and your first question should be whether they still
agree.

**Key points:**

- `X.widget.yaml` describes a widget as data: its `type`, `props` type,
  `module`/`helper`, `slots`, and `actions` — a machine-readable spec of the
  adapter and its registrations.
- The same facts are *also* hand-written in `RagWidgetType` (`core.ts`), the
  `WidgetProps` union (`props.ts`), `defaultRegistry.ts`, and the Go helper maps.
- Nothing reads the manifest, so these copies drift: 81 adapters vs 79 manifests,
  and scheduling manifests that name Go helpers which do not exist.
- Read defensively: one widget named in five files is one fact copied five times.
  Ask whether the copies agree.

---

## Chapter 8 — Reading the library: healthy decomposition and its absence

You can now read any widget in the codebase. This final chapter gives you the lens
to read them *critically* — to tell, at a glance, whether a widget follows the
engine + contract + preset discipline of Chapter 5 or predates it. This is the
substance of the companion analysis document, compressed into a way of seeing.

The newest code follows the pattern. The oldest code does not, and the clearest
example is the family of context-window diagrams:
`ContextStripDiagram`, `ContextStackDiagram`, `ContextTreemap`, `ContextBudgetBar`,
and `ContextGroupedStripDiagram`. All five draw the same thing — a set of colored,
labeled, selectable segments sized by token count — and all five draw it with their
own copy of the same code. A helper called `patternClass` is duplicated
byte-for-byte across seven files. Helpers called `styleName` and `formatTokens` are
duplicated across five. The block of code that makes a segment keyboard-navigable —
the `role`, `tabIndex`, `aria-pressed`, `onClick`, `onKeyDown` cluster — appears
five times, differing only in whether the layout is horizontal or vertical.

Measured against Chapter 5, the diagnosis is precise. There is an engine here — call
it a segment engine — but it was never extracted; instead it was copied. There is a
contract here too, and this is the telling detail: `context/types.ts:80` defines
`ContextDiagramSegment`, a struct with exactly the fields a segment payload would
need, and it has **zero references** anywhere in the codebase. Someone scaffolded
the contract, intending to build the engine, and stopped. The five diagrams are
what "before Chapter 5" looks like: the engine smeared across five files because the
seam that would have unified them was drawn on paper and never cut.

The same reading applies up and down the layers, and once you have the lens you will
spot each one:

- **Atoms.** `ContentStatusBadge`, `TranscriptRoleBadge`, and `StatusText` are the
  same engine three times — an enum mapped to a glyph and a class, rendered as a
  pill. The engine (a generic badge) was never named; each domain wrote its own.
- **Molecules.** `DataTable` is a strict subset of `MatrixGrid` — a grid engine that
  predates the general one and was never folded in. `StepList`, `KeyPointList`, and
  `CheckList` are one list engine written three times.
- **Organisms.** `CmsShell` and `CourseStudioShell` are the same sidebar-and-content
  shell, ~95% identical, because the shell engine was never extracted into a preset
  pair.
- **The IR itself.** "Read a value out of an object by a path" is implemented four
  times (`cellRenderers.tsx`, `actions.ts`, `MatrixGrid.widget.tsx`, `styleBy.ts`),
  three of them re-deriving the same `${field}` regex. The accessor was never made
  a spec, so each interpreter grew its own.

Every one of these is the same story as the context diagrams: a reusable engine
that exists implicitly, spread across the widgets that need it, because the seam was
never drawn. Learning to see the seam — to look at two widgets and ask "is there one
engine here wearing two costumes?" — is the single most useful skill for working in
this codebase. The scheduling widgets show what it looks like when the seam is cut
cleanly. The context diagrams show what it looks like when it is not. Most of the
library is somewhere between, and improving it is largely the work of finding the
implicit engines and making them explicit.

**Key points:**

- The lens: for any two similar widgets, ask whether one engine is hiding inside
  both. If the same layout/selection/rendering logic is copied, the engine exists
  implicitly and was never extracted.
- The context diagrams are the canonical example: `patternClass` copied across seven
  files, and a contract (`ContextDiagramSegment`, `context/types.ts:80`) scaffolded
  but never used — a refactor drawn and abandoned.
- The same pattern recurs in atoms (the badge family), molecules (`DataTable` vs
  `MatrixGrid`, the list family), organisms (`CmsShell` vs `CourseStudioShell`), and
  the IR (four copies of a path accessor).
- The scheduling widgets are the "seam cut cleanly" reference. Improving the library
  is largely the work of finding implicit engines and drawing their seams.

---

## Where to go next

You have the whole architecture now: the three forms of a UI (Ch.1), the IR tree and
its renderer (Ch.2), adapters and the render context (Ch.3), defunctionalized specs
(Ch.4), the engine + contract + preset pattern (Ch.5), the Go DSL that produces the
IR (Ch.6), the manifests and their drift (Ch.7), and the lens for reading the
library critically (Ch.8).

To make it concrete, do this in order:

1. Open `src/widgets/WidgetRenderer.tsx` and read `renderComponentNode`. Confirm the
   `registry.get(type)` lookup from Chapter 2.
2. Open `src/components/molecules/StepList/StepList.widget.tsx` and match every line
   to Chapter 3.
3. Open `src/widgets/presets/scheduling.ts` and read `availabilityMatrix`, then open
   `MatrixGrid.widget.tsx`, and trace one cell click through both, using the diagram
   in Chapter 5.
4. Open any two context diagrams side by side and find the duplicated
   `patternClass`. You are now reading the library the way Chapter 8 describes.

The companion document, `design-doc/01-widget-library-decomposition-analysis-and-design.md`,
turns Chapter 8's lens into a ranked, file-by-file plan for cutting the seams. Read
it once you are comfortable with the architecture here.

## Related

- `design-doc/01-widget-library-decomposition-analysis-and-design.md` — the ranked decomposition plan this guide's Chapter 8 summarizes.
- `packages/rag-evaluation-site/src/widgets/` — the IR, renderer, registry, and adapters (Chapters 2–4).
- `packages/rag-evaluation-site/src/components/molecules/MatrixGrid/` — the reference engine + contract (Chapter 5).
- `pkg/widgetdsl/module.go` — the Go DSL runtime (Chapter 6).
- Sibling ticket `RAGEVAL-SCHEDULE-WIDGETS` — the worked example of the pattern.
