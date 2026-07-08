---
Title: 'The Widget DSL by Example: TypeScript Usage, Go Counterparts, and the Case for Opaque Types'
Ticket: RAGEVAL-WIDGET-DECOMPOSITION
Status: active
Topics:
    - design-system
    - widget-ir
    - ui-dsl
    - react
    - frontend-architecture
    - intern-guide
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/cells.ts
      Note: The real precise CellSpec union the DSL .d.ts should adopt (vs the open bag today)
    - Path: repo://pkg/widgetdsl/typescript.go
      Note: Emits Props=Record<string,any> today; §8.1 proposes emitting the real per-widget props types + branded specs
    - Path: repo://pkg/widgetdsl/v2_builders.go
      Note: attachV2Ref/mustV2Ref — the existing opaque-handle pattern §8.2 generalizes to all specs
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-07T14:54:20.91334424-04:00
WhatFor: ""
WhenToUse: ""
---


# The Widget DSL by Example: TypeScript Usage, Go Counterparts, and the Case for Opaque Types

> **What this is.** A gallery of worked examples across every DSL module, ordered
> from a single line of text to a full CRM screen. Each example shows three things:
> the **TypeScript** an author writes, the **types** in play, and the **Go
> counterpart** that produces the result. Reading it, you should be able to look at
> any DSL call and know exactly what type each argument is and which Go function
> builds it.
>
> **The thread running through it.** Almost every argument in the DSL today is an
> untyped bag — `Record<string, any>` in TypeScript, `map[string]any` in Go. That
> looseness is where bad API usage hides: a misspelled prop, a `CellSpec` handed to a
> slot that wanted an `ActionSpec`, a raw object passed where a node was required —
> none of these are caught until runtime, if at all. As the examples escalate, a
> "**map-watch**" note flags each place a map still lurks, and Part 8 turns those
> observations into a concrete proposal: **opaque (branded) types and Go-side
> builders that make the wrong call unrepresentable, or at least immediately
> rejected.** This document is both a tutorial and that argument, told through code.
>
> Read after `design-doc/02` (the DSL redesign): this document is its evidence base.

## Executive Summary

The DSL is easy to call and hard to call *correctly*, because nearly every argument
is an untyped map — `Record<string, any>` on the TypeScript side, `map[string]any`
on the Go side — even though the underlying IR is precisely typed. This document
walks the whole DSL by example, from `ui.text("hi")` to a fully assembled CRM deal
board, showing for each call the TypeScript, the types, and the Go function behind
it, and marking every place an untyped map survives. It then proposes closing those
gaps with **branded/opaque types** (so a `CellSpec`, an `ActionSpec`, a `WidgetNode`
can only come from a builder, never a hand-rolled object) and **Go-side builder
functions returning sealed structs** (so the runtime rejects fabricated specs with a
clear error). The result is a DSL where the wrong call does not compile, or fails
loudly at the boundary, instead of producing subtly wrong UI.

---

## Part 0 — How to read the examples

Every example uses the same three-part layout, so you can scan them quickly.

- **TS** — what an author writes. This is the author-facing surface (the generated
  `.d.ts`), not the internal React.
- **Types** — the TypeScript types of the arguments and the return, named so you can
  find them in `src/widgets/ir/`.
- **Go** — the function in `pkg/widgetdsl/` that runs when the call is made, and the
  shape it emits.

Two vocabulary notes before we start. A **spec** is a small tagged object that stands
in for a function the IR will interpret later — `{ kind: "status", field: "stage" }`
is a `CellSpec`; `{ kind: "server", name: "deal.move" }` is an `ActionSpec`. A
**node** is a rendered UI element in IR form — `{ kind: "component", type: "Panel",
props, children }`. Specs go *inside* props; nodes go inside `children`. Confusing the
two — putting a spec where a node belongs, or vice versa — is one of the bad usages we
want the types to catch.

The current author-facing types are deliberately loose. The generated declarations
type a node as `WidgetNode { kind: string; [key: string]: any }` and every props
argument as `Props = Record<string, any>` (`typescript.go:26`). So *today* the
compiler will accept almost anything. Throughout, where an example would benefit from
a tighter type, it is shown as **(proposed)** beside the current form, so you can see
the improvement the Part 8 design would deliver.

---

## Part 1 — The type vocabulary (what everything returns)

Before the examples, meet the handful of types every call trades in. Each is a real
TypeScript type in `src/widgets/ir/`, with a Go counterpart that is currently a plain
`map[string]any`.

| Concept | TypeScript type (`src/widgets/ir/…`) | Produced by (TS) | Go builder | Go emits |
|---|---|---|---|---|
| a UI node | `WidgetNode` = `TextNode｜ElementNode｜ComponentNode` (`core.ts`) | `ui.panel(...)`, `ui.text(...)` | `componentFactory` (`module.go:565`) | `map[string]any{kind,type,props,children}` |
| an event | `ActionSpec` (`actions.ts`) | `action.server(...)` | `actionObject` (`module.go:318`) | `map[string]any{kind:"server",…}` |
| a table cell | `CellSpec` (`cells.ts`) | `cell.status(...)` | `cellObject` (`module.go:272`) | `map[string]any{kind:"status",…}` |
| a color rule | `StyleBySpec` (`engines.ts`) | `style.by(...)` *(proposed)* | *(none yet)* | — |
| a value path | `AccessorSpec` *(proposed, Part 4 of design-doc/01)* | `at(...)` *(proposed)* | *(none yet)* | — |
| a typed field | `FieldSpec` *(proposed, CRM ticket)* | `field.currency(...)` *(proposed)* | *(none yet)* | — |
| a schema handle | v2 `SchemaHandle` (opaque) | `v2.schema(...)` | `v2_builders.go` (`__widgetdsl_v2_ref`) | opaque JS handle |

Notice the last row: the v2 builder already hides a Go pointer behind an opaque JS
handle (`attachV2Ref`/`mustV2Ref`, `v2_builders.go`). That is the one place the DSL
already uses an opaque type, and it is the model Part 8 generalizes.

Here are the two node/spec shapes written out, because every example returns one of
them:

```ts
// src/widgets/ir/core.ts  — a node (what children hold)
type WidgetNode = TextNode | ElementNode | ComponentNode;
interface ComponentNode { kind: "component"; type: string; props?: object; children?: WidgetNode[]; }

// src/widgets/ir/actions.ts — a spec (what props hold)
type ActionSpec =
  | { kind: "navigate"; to: string }
  | { kind: "server";   name: string; payload?: object }
  | { kind: "copy"; value?: string } /* … */;
```

---

## Part 2 — `ui.dsl`: composing nodes

### Example 1 — the simplest possible call

```ts
// TS
ui.text("Hello");
// Types:  ui.text(value: string | number | boolean): WidgetNode
// Go:     runtime.text (module.go) → map[string]any{ "kind": "text", "text": "Hello" }
```

A leaf. `ui.text` is one of the primitives installed on every module
(`runtime.install`, `module.go:236`). Nothing loose here — the argument is a scalar.

### Example 2 — a component with props and children

```ts
// TS
ui.panel({ title: "Overview", density: "condensed" },
  ui.text("Body text"));
// Types (today):    ui.panel(props?: Props, ...children: WidgetChild[]): WidgetNode
//                   Props = Record<string, any>            ← the map
// Types (proposed): ui.panel(props?: PanelWidgetProps, ...children: WidgetNode[]): WidgetNode
// Go:               componentFactory("Panel") (module.go:565)
//                   → { kind:"component", type:"Panel", props:{title,density}, children:[…] }
```

- **map-watch.** `title` and `density` are real, documented props
  (`PanelWidgetProps` in `props.ts`), but the author-facing type is `Record<string,
  any>`, so `ui.panel({ titel: "typo" })` compiles and silently renders a titleless
  panel. The proposed signature swaps in the real `PanelWidgetProps` and the typo
  becomes a compile error. This is the single most common bad usage the opaque-type
  work removes.

### Example 3 — layout composition

```ts
// TS
ui.stack({ gap: "md" },
  ui.inline({ gap: "sm", justify: "between" },
    ui.text("Left"),
    ui.button({ variant: "primary" }, ui.text("Right"))),
  ui.divider(),
  ui.caption("footnote"));
// Types:  each is (props?, ...children) → WidgetNode; children nest freely
// Go:     four componentFactory calls; the tree is built depth-first and serialized
```

Composition is just nesting calls; the return of one call is a child argument of the
next. This is the substrate the whole DSL rests on and it is genuinely clean.

### Example 4 — a page (a top-level document)

```ts
// TS
ui.page({
  id: "overview", title: "Overview",
  sections: [ ui.panel({ title: "A" }, ui.text("…")),
              ui.panel({ title: "B" }, ui.text("…")) ],
});
// Types (today):    ui.page(options: Props): WidgetPage
// Go:               runtime.page (module.go:576) — wraps sections in a Stack unless a root node is given
```

- **map-watch.** `ui.page` takes one big options map. `id`/`title`/`sections`/`root`
  are the only meaningful keys, but the type does not say so. A `PageOptions` interface
  would name them.

---

## Part 3 — `data.dsl`: tables, cells, and actions

### Example 5 — a data table with typed cells

```ts
// TS
data.dataTable({
  rows: deals,
  getRowKey: "id",
  columns: [
    { id: "title",  header: "Deal",   cell: data.cell.field("title") },
    { id: "amount", header: "Amount", cell: data.cell.number("amount", { format: "fixed", digits: 0 }) },
    { id: "stage",  header: "Stage",  cell: data.cell.status("stage") },
  ],
});
// Types:  data.dataTable(props: DataTableWidgetProps): WidgetNode
//         DataTableColumnSpec = { id: string; header: RenderableValue; cell: CellSpec; … }   (cells.ts)
//         CellSpec = FieldCellSpec | NumberCellSpec | StatusCellSpec | …                     (cells.ts)
// Go:     dataTable helper (dataHelpers) builds the DataTable node;
//         data.cell.field/number/status → cellObject (module.go:272), each a map{kind,…}
```

- **map-watch.** Two different looseness levels sit side by side here. The `columns`
  entries are *structurally* typed on the TS IR side (`DataTableColumnSpec`), which is
  good — but the `cell` value is produced by `data.cell.status(...)`, and its return
  type in the generated `.d.ts` today is `CellSpec { kind: string; [k]: any }`, an open
  bag. So `cell: { kind: "sttatus", field: "stage" }` (a hand-rolled, misspelled spec)
  is accepted. If `CellSpec` were the real discriminated union, only the builders could
  produce a valid one.

### Example 6 — cells that carry actions

```ts
// TS
{ id: "open", header: "", cell: data.cell.actionButton("Open", data.action.navigate("/deal/${id}")) }
// Types:  data.cell.actionButton(label: RenderableValue, action: ActionSpec): CellSpec
//         data.action.navigate(to: string): ActionSpec
// Go:     cellObject.actionButton + actionObject.navigate (module.go:272 / :318)
//         → { kind:"actionButton", label, action:{ kind:"navigate", to } }
```

- **map-watch.** `actionButton` wants an `ActionSpec` as its second argument, and here
  it gets one from `data.action.navigate`. But because both `CellSpec` and `ActionSpec`
  are open bags today, nothing stops `data.cell.actionButton("Open", data.cell.field("x"))`
  — a *cell* spec where an *action* spec belongs. Distinct opaque types for `CellSpec`
  and `ActionSpec` make that a type error.

### Example 7 — the `${…}` template convention

```ts
// TS
data.action.navigate("/deal/${id}");     // ${id} interpolates against the row at dispatch time
data.cell.template("${first} ${last}");
// Types:  strings, interpreted later by the IR (actions.ts interpolate / cellRenderers.tsx renderTemplate)
// Go:     stored verbatim; the frontend interpolates against the row/context
```

- **map-watch (subtle).** `${id}` is a string microformat with no type at all — a
  typo like `${idd}` interpolates to empty. `design-doc/01` Part 4 proposes one
  `AccessorSpec` to replace these; an `at("id")` builder would at least make the field
  reference a value rather than a string fragment.

### Example 8 — a recipe with a render function (the good pattern)

```ts
// TS
data.recipes.masterDetailTable({
  rows: deals, selectedKey, columns,
  onRowSelect: data.action.navigate("/deal/${id}"),
  detail: (row) => ui.panel({ title: row.title },
    ui.keyValueStrip({ items: [ { key: "Amount", value: fmtMoney(row.amount) } ] })),
});
// Types (today):    detail?: (row: any) => WidgetNode
// Types (proposed): detail?: (row: Deal) => WidgetNode
// Go:               masterDetailTableRecipe (module.go:853); detail invoked via detailNode (module.go:882)
//                   → DashboardGrid(two-up)[ Panel[DataTable], <detail(selectedRow)> ]
```

This is the one place the current DSL lets you pass a real function and composes with
it. It is the seed of `design-doc/02`'s whole proposal.

- **map-watch.** The options object is a single untyped bag mixing data (`rows`),
  specs (`onRowSelect`), and a function (`detail`). A `MasterDetailOptions<Row>`
  interface — generic in the row type — would make `detail`'s `row` parameter `Deal`,
  not `any`, giving you completion inside the render function.

### Example 9 — the intent grammar (describe records, not components)

```ts
// TS
data.collection(deals, {
  schema: data.schema({
    title:  data.f.primary(),
    stage:  data.f.status(),
    amount: data.f.measure(),
    owner:  data.f.short(),
  }),
  verb: "show",
});
// Types:  data.f.<role>() : FieldRole ;  data.schema(fields): SchemaHandle
//         data.collection(rows, opts): WidgetNode
// Go:     grammar.go — f.<role> (fieldRoleObject :64), schema (schemaCtor :79, tags __ragSchema),
//         collectionVerb (:265) → SectionBlock > Stack > DataTable (+ optional detail)
```

The grammar is a step up in abstraction: you name each field's *role* and the grammar
picks the components. Roles are `primary`, `status`, `measure`, `date`, `tags`, and so
on (`fieldRoles`, `grammar.go:25`).

- **map-watch.** `data.schema({...})` returns a value tagged internally with
  `__ragSchema`, but the author-facing type is loose. And the roles are stringly
  chosen behind the `f.*` methods. This is the closest the classic DSL gets to opaque
  types — the `__ragSchema` tag is a brand in spirit — and Part 8 makes it real.

---

## Part 4 — Engine verbs: composition + typing together (proposed surface)

These use the `design-doc/02` verbs. They are the payoff: the place where a fully
typed unit-renderer function makes both the composition and the types shine. Shown as
proposed TS with the intended types.

### Example 10 — a matrix (the availability poll) in one call

```ts
// TS (proposed)
ui.matrix(responses, {
  columns: options.map(o => ({ id: o.id, header: slotLabel(o) })),
  valueAt: at.map("cells"),                              // row.cells[colId]
  cell: cell.cycle(["yes", "ifneedbe", "no"], { styleSet: availabilityPalette }),
  editableRow: "you",
  onCell: act.server("poll.toggleCell"),
});
// Types:  ui.matrix<Row>(rows: Row[], opts: MatrixOptions<Row>): WidgetNode
//         MatrixOptions<Row> = { columns: ColumnSpec[]; valueAt: AccessorSpec;
//                                cell: CellSpec | CycleCellSpec; editableRow?: string;
//                                onCell?: ActionSpec; colorBy?: StyleBySpec }
//         cell.cycle(states: string[], o?: {glyphs?; styleSet?}): CycleCellSpec        (engines.ts)
// Go:     ui.matrix helper → MatrixGrid node; cell.cycle → cellObject.cycle (proposed, module.go:272)
```

- **map-watch: mostly closed.** Every field of `MatrixOptions` is a named type; the
  only remaining bag would be `columns`, and even that is `ColumnSpec[]`. `cell.cycle`
  returns a `CycleCellSpec`, not `any`, so passing `cell: act.server("x")` (an action
  where a cell belongs) is a type error. This is the target state.

### Example 11 — a heatmap via `colorBy` (StyleBySpec)

```ts
// TS (proposed)
ui.matrix(people, {
  columns: skills.map(s => ({ id: s.id, header: s.name })),
  valueAt: at.map("ratings"),
  cell: cell.value(),                                    // render the number
  colorBy: style.by(heatPalette),                        // value → styleKey → color
});
// Types:  cell.value(): ValueCellSpec ;  style.by(set: ContextStyleSet, o?): StyleBySpec   (engines.ts)
// Go:     ui.matrix + cell.value + style.by (proposed builders); StyleBySpec resolved by styleBy.ts
```

- **map-watch.** `style.by` returns a `StyleBySpec` — a distinct type from `CellSpec`
  and `ActionSpec` — so it can only be assigned to `colorBy`, never to `cell` or
  `onCell`. The slot and the spec type are locked together.

### Example 12 — a kanban board with a typed card renderer

```ts
// TS (proposed)
ui.board(deals, {
  columns: pipeline.stages.map(s => ({ id: s.id, header: `${s.name} · ${fmtMoney(sum(s))}` })),
  columnOf: at("stageId"),
  card: (deal: Deal) => ui.panel({ tone: "raised" },
    ui.text(deal.title),
    ui.inline({ gap: "xs" }, ui.avatar(deal.owner), ui.caption(fmtMoney(deal.amount)))),
  onMove: act.server("deal.move"),
  selected: select.single("id"),
});
// Types:  ui.board<Card>(cards: Card[], opts: BoardOptions<Card>): WidgetNode
//         BoardOptions<Card> = { columns: ColumnSpec[]; columnOf: AccessorSpec;
//                                card: (c: Card) => WidgetNode; onMove?: ActionSpec;
//                                selected?: SelectionSpec }
// Go:     ui.board helper; card invoked per datum via the generalized unit-renderer (module.go:882)
```

- **map-watch: closed, and generic.** Because `BoardOptions<Card>` is generic in the
  card type, `deal` inside the `card` function is `Deal`, so `deal.titel` is a compile
  error *inside the renderer*. The looseness that plagued Example 2 is gone precisely
  where authors spend the most time — writing the card.

### Example 13 — a timeline whose row type switches on a discriminant

```ts
// TS (proposed)
ui.timeline(activities, {
  item: (a: Activity) => {
    switch (a.kind) {
      case "email": return ui.inline({ gap: "xs" }, ui.icon("mail"), ui.text(a.title));
      case "call":  return ui.inline({ gap: "xs" }, ui.icon("phone"), ui.text(`${a.title} · ${a.durationMin}m`));
      default:      return ui.text(a.title);
    }
  },
  onOpen: act.navigate("/activity/${id}"),
});
// Types:  ui.timeline<Item>(items: Item[], opts: { item: (i: Item) => WidgetNode; onOpen?: ActionSpec }): WidgetNode
//         Activity is a discriminated union on `kind` (crm/types.ts) → exhaustiveness checking
// Go:     ui.timeline helper; item invoked per activity
```

- **map-watch: closed, with a bonus.** Because `Activity` is a discriminated union, the
  `switch` gets exhaustiveness checking — add a new `kind` and TypeScript flags the
  missing case. No map anywhere in the author's code.

### Example 14 — a typed record (CRM), read mode

```ts
// TS (proposed)
ui.record(contact, {
  sections: [
    { label: "Details", fields: [
      field.email("email"),
      field.phone("phone"),
      field.user("owner"),
      field.select("stage", { options: stageOptions, styleSet: stagePalette }),
    ]},
  ],
  mode: "read",
  onFieldChange: act.server("field.update"),
});
// Types:  field.<type>(key: string, o?): FieldSpec ;  ui.record<R>(rec: R, opts: RecordOptions): WidgetNode
//         FieldSpec discriminated on `type` (crm ticket Part 4)
// Go:     ui.record + field.* builders (proposed crm.dsl); each field.* → a FieldSpec map today
```

- **map-watch.** `field.select` needs `options` and `styleSet`; those belong only to
  `select`/`multiselect`, and a discriminated `FieldSpec` would make `field.email("x",
  { options })` a type error — you cannot give an email field select options.

---

## Part 5 — the domain modules (`context_window`, `course`, `cms`)

A few representative calls so the whole surface is covered.

### Example 15 — a context-window diagram (`context_window.dsl`)

```ts
// TS
cw.contextStripDiagram({ snapshot, styleSet });
// Types:  cw.contextStripDiagram(props: ContextStripDiagramWidgetProps): WidgetNode  (props.ts)
//         snapshot: ContextWindowSnapshot ; styleSet: ContextStyleSet                 (context/types.ts)
// Go:     contextWindowHelpers["contextStripDiagram"] → componentFactory("ContextStripDiagram")
```

### Example 16 — a palette-driven style set (`context_window.dsl`)

```ts
// TS
const palette = cw.paletteStyleSet("signalOrangeCyan", { entries });
cw.contextDiagramPanel({ snapshot, styleSet: palette, initialView: "strip" });
// Types:  cw.paletteStyleSet(name: PaletteName, o): ContextStyleSet
// Go:     buildPaletteStyleSet (module.go:393) — the one place a ContextStyleSet is constructed Go-side
```

- **map-watch.** `paletteStyleSet` returns a `ContextStyleSet`, a real type — this is a
  good precedent. Its reuse (`module.go:393`) is exactly what the proposed
  `cell.cycle`/`style.by` builders need to construct their palettes.

### Example 17 — a course studio shell (`course.dsl`)

```ts
// TS
course.courseStudioShell(
  { sections: navSections, activeItemId: "overview", title: "Course" },
  course.courseLessonPanel({ course: courseData, onAgendaItemSelectAction: act.navigate("/lesson/${id}") }));
// Types:  each helper (props: XWidgetProps): WidgetNode ;  courseData: ContextCourse (context/types.ts)
// Go:     courseHelpers → componentFactory for each; the shell holds the panel as a child
```

### Example 18 — a media library (`cms.dsl`) with a recipe

```ts
// TS
cms.recipes.mediaLibrary({
  assets, page: 1, pageCount: 4,
  onAssetOpen: act.navigate("/asset/${id}"),
  onFilesSelected: act.event("cms.upload"),
});
// Types (today):  a single options bag → WidgetNode
// Go:             mediaLibraryRecipe (module.go) — ~20 lines of copyIfPresent + an action-map loop
```

- **map-watch.** This recipe is the poster child for `design-doc/01` Part 6's
  "declarative recipe" idea: the options bag maps option names to prop names by hand.
  A typed `MediaLibraryOptions` plus a declarative recipe descriptor removes both the
  bag and the hand-mapping.

---

## Part 6 — `data.v2.dsl`: the fluent, already-opaque surface

### Example 19 — the fluent builder

```ts
// TS
data.v2
  .collection("deals", deals)
  .schema(s => s.field("title").primary().field("amount").currency().field("stage").status())
  .select("id").edit()
  .table(t => t.columns("title", "amount", "stage"))
  .toIR();
// Types:  collection(name, rows): CollectionBuilder ;  .schema/.select/.edit/.table return the builder;
//         .toIR(): WidgetNode ;  the builder carries an opaque ref (__widgetdsl_v2_ref)
// Go:     v2_builders.go — attachV2Ref (:348) hides a *v2spec.CollectionSpec pointer behind a JS handle;
//         .toIR() validates (v2/spec/validate.go) then lowers (v2/spec/lower.go)
```

This is the one surface that already does what Part 8 recommends everywhere. The
builder is an **opaque handle**: you cannot fabricate one, you can only obtain it from
`data.v2.collection(...)` and transform it with its methods. `.toIR()` validates the
accumulated spec (`validate.go`) before producing a node, so an incoherent collection
(two key fields, a missing action target) is rejected with a diagnostic — not rendered
wrong. The lesson of this document is: **make the rest of the DSL feel like this.**

- **map-watch: none.** There is no `Record<string, any>` in sight; every step is a
  typed method on an opaque builder.

---

## Part 7 — Two full screens, end to end

### Example 20 — a dashboard (composition of existing engines)

```ts
// TS (proposed)
ui.dashboard({ recipe: "metrics" },
  ui.statTile({ label: "Open pipeline", value: fmtMoney(total), delta: +0.12 }),
  ui.statTile({ label: "Win rate", value: "38%", delta: -0.03 }),
  ui.panel({ title: "Pipeline" }, ui.segmentedBar({ segments: stageSegments, styleSet: stagePalette })),
  ui.panel({ title: "My tasks" }, ui.list(tasks, { item: (t: Task) => ui.checkboxRow({ checked: t.done }, ui.text(t.title)) })));
// Types:  ui.dashboard(props, ...children) ; ui.statTile(props: StatTileWidgetProps) ;
//         ui.list<T>(items: T[], opts: { item: (t:T)=>WidgetNode })
// Go:     dashboard → DashboardGrid; statTile/segmentedBar/list helpers; list.item per-datum
```

### Example 21 — the capstone: a CRM deal board screen

```ts
// TS (proposed) — the whole screen, fully typed, no untyped map in author code
ui.studioShell({ sections: crmNav, activeItemId: "deals", title: "Sales" },
  ui.panel({ title: "Pipeline", actions: ui.button({ variant: "primary" }, ui.text("+ Deal")) },
    ui.board<Deal>(deals, {
      columns: pipeline.stages.map(s => ({ id: s.id, header: `${s.name} · ${fmtMoney(sumStage(deals, s.id))}` })),
      columnOf: at("stageId"),
      selected: select.single("id"),
      onMove: act.server("deal.move"),
      card: (deal: Deal) => ui.panel({ tone: "raised" },
        ui.text(deal.title),
        ui.inline({ gap: "xs", justify: "between" },
          ui.avatar(deal.owner),
          ui.caption(fmtMoney(deal.amount)),
          ui.when(deal.status === "won", ui.tag({ label: "won" })))),
    })));
// Types in play:  WidgetNode, BoardOptions<Deal>, AccessorSpec, SelectionSpec, ActionSpec,
//                 StudioShellWidgetProps, PanelWidgetProps, ButtonWidgetProps
// Go:             studioShell + panel + board (+ per-card unit renderer) + button/avatar/caption/tag/when
```

Read the capstone against Example 2. The only place the author touches free-form data
is `deal.title`/`deal.amount` — and those are typed `Deal` fields, checked by the
compiler, because `ui.board<Deal>` threads the type into the `card` function. Every
spec (`at`, `select.single`, `act.server`) is a distinct opaque type bound to exactly
one option slot. There is no `Record<string, any>` anywhere in what the author wrote.
That is the destination.

---

## Part 8 — The case for opaque types and Go-side builders

The examples make the problem concrete: the DSL is loosely typed at exactly the points
where mistakes are easy and expensive. This part is the proposal, in two halves —
TypeScript and Go — plus the table of what changes.

### 8.1 TypeScript: brand every spec and node

A **branded** (or nominal, or opaque) type is a structural type tagged with a phantom
marker so that two types with the same shape are nonetheless incompatible, and so that
a value of the type can only be produced by code that has the marker. The mechanism is
a `unique symbol`:

```ts
declare const brand: unique symbol;
type Brand<T, B extends string> = T & { readonly [brand]: B };

export type WidgetNode = Brand<{ kind: string }, "WidgetNode">;
export type ActionSpec = Brand<{ kind: string }, "ActionSpec">;
export type CellSpec   = Brand<{ kind: string }, "CellSpec">;
export type StyleBySpec = Brand<{ kind: "styleBy" }, "StyleBySpec">;
export type FieldSpec  = Brand<{ kind: string }, "FieldSpec">;
```

With these, three things become impossible to express, which is to say three bad
usages become compile errors:

- **A hand-rolled spec.** `cell: { kind: "status", field: "stage" }` no longer
  type-checks against `cell: CellSpec`, because a plain object lacks the brand. The
  *only* way to get a `CellSpec` is to call `cell.status(...)`. Typos in the spec shape
  are now unrepresentable, not merely discouraged.
- **A spec in the wrong slot.** `onCell: cell.value()` fails, because `onCell` wants an
  `ActionSpec` and `cell.value()` returns a `CellSpec` — different brands, no
  assignment. Example 6's cell-where-an-action-belongs bug cannot be written.
- **A raw object as a child.** `ui.panel({}, { kind: "component", type: "X" })` fails,
  because `children` wants `WidgetNode` and a bare object is not branded. Children must
  come from DSL calls.

Alongside branding, replace `Props = Record<string, any>` with the real per-widget
props interfaces the IR already defines (`PanelWidgetProps`, `DataTableWidgetProps`,
…). The generator that emits the `.d.ts` (`typescript.go`) already knows each widget's
props type name from its manifest (`design-doc/01` Part 5); it simply emits
`Record<string, any>` today. Emitting the real type name is a small change with a large
effect — Example 2's `titel` typo becomes an error.

### 8.2 Go: sealed builder returns and a boundary brand

The TypeScript brand protects authors who compile. But the DSL also runs untyped
scripts, and the Go runtime must not trust its input. Today a Go builder returns
`map[string]any` and the runtime consumes `map[string]any`, so a script that passes a
fabricated `{ kind: "status" }` where a cell belongs is indistinguishable from one the
builder produced. Two changes fix this.

First, **Go builders return sealed types, not bare maps.** Define a marker interface
with an unexported method so no code outside the package can implement it, and have the
builders return concrete types satisfying it:

```go
// pkg/widgetdsl — a sealed spec
type CellSpec interface{ isCellSpec() }
type statusCell struct{ Field string }
func (statusCell) isCellSpec() {}
func (r *runtime) cellStatus(field string) CellSpec { return statusCell{Field: field} }  // was: map[string]any
```

Internally this replaces the `map[string]any` soup (`cellObject` at `module.go:272`,
`actionObject` at `module.go:318`) with typed values that only serialize to a map at
the very edge, when the node tree is handed to goja. The compiler now enforces, inside
the Go code, that a cell builder returns a cell.

Second, **brand the values that cross the goja boundary**, exactly as v2 already does.
The v2 builder attaches a hidden `__widgetdsl_v2_ref` data-property to its JS handles
(`attachV2Ref`/`mustV2Ref`, `v2_builders.go`) and refuses to proceed if a handle lacks
it. Generalize that: every spec object a builder returns to a script carries a hidden
brand naming its kind, and every verb that *consumes* a spec checks the brand and
panics with a clear message if it is missing or wrong. A script that hand-rolls a
`{ kind: "status" }` and passes it as an action now fails at the call with
"expected an ActionSpec built by action.*, got an unbranded object" instead of emitting
subtly wrong IR. The opaque handle is the runtime counterpart of the TypeScript brand.

### 8.3 What changes, concretely

| Today | Proposed | Bad usage it catches |
|---|---|---|
| `Props = Record<string, any>` | real `XWidgetProps` per helper | misspelled / unknown props |
| `CellSpec = { kind: string; [k]: any }` | branded `CellSpec`, builder-only | hand-rolled or misspelled cell specs |
| `ActionSpec`, `StyleBySpec`, `FieldSpec` open | branded, builder-only, distinct | a spec in the wrong option slot |
| `WidgetNode = { kind: string; [k]: any }` | branded `WidgetNode` | raw objects passed as children |
| `detail: (row: any) => …` | `detail: (row: Row) => …` (generic) | field typos inside render functions |
| Go builders return `map[string]any` | sealed Go types, serialized at the edge | a Go builder returning the wrong shape |
| goja specs are anonymous maps | hidden per-kind brand + verb-side check | fabricated specs from untyped scripts |

The through-line is one idea stated twice, once per language: **a spec should be
obtainable only from its builder, and recognizable as such by whatever consumes it.**
TypeScript brands give authors compile-time enforcement; Go sealed types plus boundary
brands give the runtime enforcement for scripts that never see the compiler. Together
they turn the whole family of "wrong shape in the wrong place" bugs — the ones the
examples kept flagging — from silent wrong-UI into loud, early errors.

### 8.4 Cost and staging

This is additive and can land behind the examples that motivate it. Brand the specs and
switch the `.d.ts` to real props types first (pure TypeScript, no runtime change).
Introduce the Go sealed types next, serializing at the edge so the emitted IR is
byte-identical. Add the boundary brand-check last, as a hardening pass, once the
builders are the only source of specs. At each stage the emitted JSON is unchanged, so
the frontend never notices — only the ways to *get it wrong* disappear.

---

## Open questions

1. **How aggressively do we brand?** Branding `WidgetNode` is the highest-value and
   the most disruptive (every helper's children type changes). Brand specs first and
   nodes later, or all at once?
2. **Ergonomics of brands for authors.** Branded types can produce noisier error
   messages. Is the safety worth the occasional cryptic mismatch, or do we brand only
   the specs most often confused (`CellSpec` vs `ActionSpec`)?
3. **Go sealing vs. serialization cost.** Sealed Go types add allocation and a
   serialize-at-edge step. Is that measurable at the tree sizes we emit?
4. **Boundary-brand strictness.** Should an unbranded spec from a script be a hard
   panic, or a warning that still renders? Hard is safer; soft eases migration of
   existing scripts.
5. **Where does `at`/`AccessorSpec` sit** — its own branded type, or a plain string
   with a builder? (Example 7.)

## References

- `design-doc/02-redesigning-the-widget-dsl-...` — the composition-first redesign these
  examples exercise; this document is its evidence base.
- `design-doc/01-...` — Part 4 (the specs to unify), Part 5 (manifest-driven `.d.ts`),
  Part 6 (the DSL analysis).
- `RAGEVAL-CRM-WIDGETS` — the `field.*`/`FieldSpec` examples (Part 4, 14).
- `RAGEVAL-SCHEDULE-WIDGETS` — the `matrix`/`board`/`calendar` engines the verbs expose.
- `pkg/widgetdsl` — `componentFactory` (`module.go:565`), `cellObject` (`:272`),
  `actionObject` (`:318`), `recipesObject`/`masterDetailTableRecipe` (`:495`/`:853`),
  `detailNode` (`:882`), `buildPaletteStyleSet` (`:393`), `grammar.go`,
  `v2_builders.go` (`attachV2Ref`/`mustV2Ref` — the existing opaque handle),
  `typescript.go` (`TypeScriptModule` `:14`), `v2/spec/{validate,lower}.go`.
- `src/widgets/ir/{core,actions,cells,engines,props}.ts` — the real, precise IR types
  the DSL `.d.ts` should adopt.
