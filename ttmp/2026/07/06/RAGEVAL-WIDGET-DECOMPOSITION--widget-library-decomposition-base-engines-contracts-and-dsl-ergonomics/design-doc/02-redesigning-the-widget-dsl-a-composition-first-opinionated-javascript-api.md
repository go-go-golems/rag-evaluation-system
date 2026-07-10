---
Title: 'Redesigning the Widget DSL: A Composition-First, Opinionated JavaScript API'
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
    - Path: repo://pkg/widgetdsl/module.go
      Note: detailNode (:882) is the unit-renderer prototype; componentFactory/cellObject/actionObject/recipesObject are the surfaces the redesign extends
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-07T14:45:17.296138119-04:00
WhatFor: ""
WhenToUse: ""
---


# Redesigning the Widget DSL: A Composition-First, Opinionated JavaScript API

> **What this document is.** The companion analysis (`design-doc/01`) surveyed the
> whole widget library and gave the Go/Goja DSL one part (Part 6) among many. This
> document zooms in on that DSL — the JavaScript API the goja modules expose to
> script authors — and proposes a redesign. The goal is a DSL whose center of
> gravity is *composition*: assembling screens out of engines and units, the way the
> React and IR layers have evolved to work, rather than filling option bags on
> flat helpers. The redesign stays deliberately **opinionated** — it offers a small
> set of curated, high-level verbs and makes the low-level escape hatch explicit and
> rare — and it is designed to be applied *backwards* onto the existing modules
> without breaking the scripts already written against them.
>
> **Who should read it.** The engineer who owns `pkg/widgetdsl`, and anyone who
> writes DSL scripts. It assumes you have read `design-doc/01` Part 6 and, ideally,
> the architecture textbook (`reference/02-...`). Where it cites `module.go:565`,
> that is a real symbol you can open.

## Executive Summary

The Go/Goja Widget DSL grew as five overlapping authoring styles — flat component
helpers, spec builders, hand-written recipes, an intent grammar, and a fluent typed
experiment — and as a result composition is awkward, the surface lags the IR types,
and the raw `component()` escape hatch has become the default path for anything
interesting. This document proposes a composition-first, opinionated redesign: a small
set of **engine verbs** (`matrix`, `board`, `timeline`, `record`, `list`, `calendar`),
each taking data plus a **unit renderer** — a real JavaScript function the runtime calls
to build each card/cell/row and then serializes — which exposes the library's
engine/contract/preset pattern directly as the authoring surface. The everyday surface
is curated and small; the low-level substrate is retained but renamed `raw()` to make
escaping visible. The redesign is additive and lands incrementally onto the existing
modules, converging the two collection surfaces and generating the whole vocabulary from
the manifest set so it never drifts again.

---

## Part 0 — How to read this, and the one question it answers

The DSL's job is narrow and worth stating precisely: it is a set of JavaScript
functions, implemented in Go and run inside an embedded engine (goja), whose return
values are **Widget IR nodes** — plain JSON objects of the form
`{ kind: "component", type, props, children }`. A script calls the functions; the
returned JSON is handed to the frontend, which renders it. The DSL never renders
anything itself. It is an *authoring surface* for the IR.

That framing produces the single question this document answers: **what is the most
composable, most opinionated authoring surface we can put in front of the IR?**
"Composable" means screens are built by nesting small pieces, and the language makes
nesting natural. "Opinionated" means the language offers a few good ways to do a
thing rather than every possible way — it curates. Those two goals are in tension
only if you design badly; the evolved React style (engines that arrange, units that
render, presets that configure) resolves the tension, and the core proposal here is
to expose that same three-part shape *as the DSL*.

The document proceeds by teaching the DSL as it is (Part 1), naming exactly where
authoring hurts today (Part 2), stating the principles that resolve those hurts
(Part 3), proposing the concrete new surface (Part 4), showing how it stays
opinionated (Part 5), and explaining how to retrofit it onto the running system
without a rewrite (Part 6). Parts 7–8 are for the implementer.

---

## Part 1 — The DSL as it is today

The best way to understand the current DSL is to watch it build something. A script
first acquires a module and then calls its functions:

```js
const ui = require("ui.dsl");
ui.panel({ title: "Overview" },
  ui.text("Hello"),
  ui.button({ variant: "primary" }, ui.text("Save")));
```

`ui.panel(props, ...children)` returns `{ kind:"component", type:"Panel", props:{…},
children:[…] }`. The mechanism is uniform: a module's `moduleSpec` (`module.go:23`)
holds a `helpers` map from a JavaScript name to an IR component `type`, and
`runtime.install` (`module.go:236`) turns each entry into a function built by
`componentFactory` (`module.go:565`). The first argument is treated as props if it
is a plain object that does not itself look like a node; everything after is a child.
So `panel`, `button`, `stepList`, and the rest are all the same factory over
different type strings. This part is clean and should survive the redesign untouched;
it is the substrate everything else sits on.

Around that substrate are **four additional authoring styles**, and the fact that
there are four — each with its own philosophy — is the seed of the problems in
Part 2. It is worth seeing one example of each so the redesign has something concrete
to improve.

The first extra style is **spec builders**: small functions that construct the
defunctionalized specs (the data-descriptions-of-functions from the IR). Cells and
actions each get a namespace:

```js
const data = require("data.dsl");
data.dataTable({
  rows,
  columns: [
    { id: "name",   header: "Name",   cell: data.cell.field("name") },
    { id: "stage",  header: "Stage",  cell: data.cell.status("stage") },
    { id: "open",   header: "",       cell: data.cell.actionButton("Open", data.action.navigate("/x/${id}")) },
  ],
});
```

`data.cell.*` (`cellObject`, `module.go:272`) and `data.action.*` (`actionObject`,
`module.go:318`) return the `CellSpec`/`ActionSpec` objects the IR interprets.

The second extra style is **recipes**: composite builders that emit a whole subtree
at once, namespaced under `module.recipes`. This one already contains the seed of the
good idea — it accepts a *function* to render the detail pane:

```js
data.recipes.masterDetailTable({
  rows, columns, selectedKey,
  onRowSelect: data.action.navigate("/deal/${id}"),
  detail: (row) => ui.panel({ title: row.title }, ui.text(row.summary)),   // ← a real JS function
});
```

Look closely at `detail`. `masterDetailTableRecipe` (`module.go:853`) receives that
JavaScript function, invokes it with a row, and serializes whatever node it returns
(`detailNode`, `module.go:882`, via `goja.AssertFunction`). This is the crucial
capability the redesign leans on and Part 3 returns to: **the DSL runs in a real JS
engine, so authors can pass real functions at authoring time, even though the output
must be plain JSON.** The recipe calls the function, gets a node, and embeds it. Only
one builder in the entire DSL currently exploits this.

The third extra style is the **data grammar** (`grammar.go`): higher-level intent
verbs. Instead of describing a table, you describe a *collection of records with
typed fields* and let the grammar choose the components:

```js
data.collection(rows, {
  schema: data.schema({
    name:   data.f.primary(),
    stage:  data.f.status(),
    amount: data.f.measure(),
  }),
  verb: "edit",
  arrange: "master-detail",
});
```

`f.<role>()` (`grammar.go`) names a field by *role* (`primary`, `status`, `measure`,
`date`, `tags`, …), and `collectionVerb` (`grammar.go:265`) compiles the whole thing
into a `SectionBlock` > `DataTable` (+ a detail form) subtree.

The fourth extra style is the **typed v2 experiment** (`data.v2.dsl`,
`v2_builders.go`): the same collection idea expressed as a fluent chain that keeps a
typed intermediate and validates it before lowering to IR:

```js
const v2 = require("data.v2.dsl");
v2.collection("deals", rows)
  .schema(s).select("id").edit()
  .table(t => t.columns("name", "stage", "amount"))
  .toIR();
```

So: a substrate of flat component helpers, plus spec builders, plus recipes, plus a
grammar, plus a fluent typed builder. Five surfaces. Each is individually reasonable.
Together they are the problem.

---

## Part 2 — Where authoring hurts

Five concrete frictions, each grounded in the styles above. These are the things the
redesign must fix.

**1. Composition is a second-class citizen.** The substrate composes by nesting
`component(props, ...children)` calls, which is fine for static trees but awkward the
moment structure depends on data. To build a board you would iterate stages, and
within each iterate deals, hand-nesting `component` calls and threading keys and
props yourself — and there is no verb for "columns of draggable cards," so you would
reach for the raw substrate every time. The one place the DSL *does* offer real
compositional power — passing a `detail(row)` function to a recipe — is available in
exactly one recipe and nowhere else. Composition is possible but never encouraged.

**2. Two philosophies of the same thing.** The classic grammar
(`data.collection(rows, { …option bag… })`) and the v2 fluent builder
(`v2.collection(…).schema().edit().toIR()`) are two authoring models for one
operation, and they disagree on everything: option bag versus method chain, untyped
versus typed, immediate versus `.toIR()`-terminated. An author must learn both and
choose, and the two share no vocabulary. Internally they are *also* two compilers
(`grammar.go` and `v2/spec/lower.go`) — `design-doc/01` Part 6 #5 covers that
duplication. The authoring split is the user-facing half of the same wound.

**3. The surface lags the IR.** The React and IR layers have grown a vocabulary the
DSL cannot speak. There is no `cell.cycle` or `cell.value`, no `style.by`, no engine
verb for a matrix, board, calendar, or timeline, and (per the CRM kit) no `field`
builder. An author who wants an availability grid, a heatmap, a pipeline board, or a
typed record has to drop to `component("MatrixGrid", { …hand-written props… })` — the
raw substrate — reproducing by hand what a preset does in TypeScript. The friendly
layer stops exactly where the interesting widgets begin.

**4. The escape hatch is unlabeled and overused.** `component(type, props, …)` is the
"I'll assemble the raw node myself" tool. It is necessary — every good DSL needs an
escape hatch — but here it is not marked as one, and because of frictions 1 and 3 it
is the *default* path for anything the curated helpers do not cover. An escape hatch
that authors live in is a sign the curated surface is too small.

**5. Inconsistency and drift.** Recipes are hand-written per case
(`mediaLibraryRecipe`, `articleListRecipe` are twenty-plus lines of `copyIfPresent`
plus an action-map loop each), module capabilities are wired by `if spec.name == …`
chains in two places that must be kept in sync (`runtime.install` and
`typescript.go`), and helpers can name component types the manifests never verify
exist (`design-doc/01` Part 5). The surface is not generated from one description, so
it drifts and it is laborious to extend.

---

## Part 3 — Design principles

Six principles resolve Part 2. They are stated as commitments, each with its
rationale, because the concrete API in Part 4 is just their mechanical consequence.

**Composition is the primary verb, and functions are how you compose.** The DSL runs
in a real JavaScript engine, and the `detail(row)` pattern proves the runtime can
accept an author's function, call it, and serialize the node it returns. Make that
the norm, not the exception. Every engine verb takes a **unit renderer** — a function
from one datum to a node — and the verb calls it to build each child. This is the
authoring-time expression of the engine/contract/preset pattern: the verb is the
engine, the function you pass is the unit, and the arguments you fill are the preset.
Authors compose by writing small functions, and the output is still plain JSON.

**Opinionated by default, raw by exception.** The everyday surface is a small set of
high-level verbs (`matrix`, `board`, `list`, `timeline`, `calendar`, `record`,
`collection`). They bake in the good layout and interaction decisions and expose only
the choices that matter. The low-level `component()` substrate remains, but it is
renamed to advertise its role — `raw()` — so that reaching for it is a visible,
deliberate act rather than the path of least resistance. The measure of success is
that a well-written script contains almost no `raw()`.

**One authoring model.** Pick option-bag-with-a-render-function and use it
everywhere. A verb is `verb(data, options)` where `options` carries the unit renderer
and the specs. No parallel fluent dialect for some operations and option bags for
others. (The typed intermediate that v2 introduced is valuable, but it belongs
*under* the surface as the shared lowering, not as a second surface authors must
choose — `design-doc/01` Part 6 #5.)

**Specs are first-class and complete.** The defunctionalized specs get one coherent
set of builder namespaces that covers everything the IR can express today:
`cell.*` (including `cycle` and `value`), `act.*` (actions), `style.by(...)`,
`at(...)` (a value accessor), `field.*` (typed CRM fields), and `select.*`
(selection). Parity with the IR is a standing requirement, ideally enforced by a
test, so the surface never again lags the types.

**Presets are authored declaratively, not hand-coded.** A recipe today is a bespoke
Go function. Most recipes are the same shape: take some data, map a few option names
onto prop names, wrap in a component. Express that shape as data — a small descriptor
— and interpret it once, so a new preset is a table entry, not a new function.

**Everything is generated from, or checked against, one description.** The helper
maps, the TypeScript declarations, and the manifest set should agree by construction,
not by discipline. This is the DSL half of `design-doc/01` Part 5.

---

## Part 4 — The proposed surface

This part is the concrete redesign. It is additive to the substrate: `component`,
`text`, `element`, `fragment`, and `page` stay exactly as they are; the new surface
sits on top and is what authors reach for.

### 4.1 Composition primitives

Keep `component`/`text`/`element`/`fragment`/`page`. Rename `component` to `raw` as an
exported alias to signal "escape hatch" (keep `component` as a deprecated alias so no
script breaks). Add three tiny composition helpers that remove the most common
hand-nesting:

```js
ui.when(cond, node)              // node or nothing — conditional children
ui.map(items, fn)               // items.map(fn) that flattens to children (a typed fragment)
ui.slot(node | fn)              // marks a value as a renderable slot (node, or a function to call)
```

`ui.map` is what you reach for instead of hand-writing a `for` loop that pushes
`component` calls; `ui.when` replaces the `cond ? node : undefined` ternary that
litters scripts; `ui.slot` names the "this prop is a subtree or a function that
returns one" convention that recipes already rely on informally.

### 4.2 Spec builders, unified and complete

One set of namespaces, available on every module that needs them:

```js
// value accessors — the ONE way to say "read this out of a datum" (design-doc/01 Part 4 #1)
at("amount")                    // { field: "amount" }
at("owner.name")                // dotted path
at.tmpl("${first} ${last}")     // template

// cells — the existing set PLUS the evolved ones
cell.field("name"), cell.status("stage"), cell.currency("amount"), cell.template("${a}"),
cell.cycle(["yes","ifneedbe","no"], { glyphs }),        // NEW — availability/RSVP
cell.value({ colorBy: style.by(palette) }),             // NEW — heatmap value cell

// colors — value → styleKey → visual style (design-doc/01 Part 4)
style.by(styleSet, { field: "stage", map, fallback })

// actions — unchanged kinds, tidier names
act.server("deal.move", { payload }), act.navigate("/deal/${id}"), act.copy("x"), act.event("print")

// fields — the CRM typed-field builder (from RAGEVAL-CRM-WIDGETS Part 4)
field.currency("amount"), field.select("stage", { options }), field.relation("company", { object:"company" })

// selection — one shape instead of ~14 bespoke fields (design-doc/01 Part 4 #2)
select.single("id"), select.multi("id")
```

These are thin builders returning the plain spec objects the IR already defines; the
value is a single, discoverable, typed vocabulary instead of five scattered partial
ones.

### 4.3 Engine verbs — the heart of the redesign

This is where composition becomes the primary verb. Each engine gets a verb of the
form `verb(data, options)`, where `options` always contains a **unit renderer** (a
function the runtime calls to build each unit) plus the specs that configure the
engine. The verb is the engine; the function is the unit; the options are the preset.

```js
// a board: columns of draggable cards. `card` is called per datum → a node.
ui.board(deals, {
  columns: pipeline.stages.map(s => ({ id: s.id, header: s.name })),
  columnOf: at("stageId"),
  card: (deal) => ui.panel({ tone: "raised" },              // ← unit renderer, a real function
    ui.text(deal.title), ui.caption(fmtMoney(deal.amount))),
  onMove: act.server("deal.move"),
  selected: select.single("id"),
})

// a matrix (the scheduling grid), the availability poll in four lines:
ui.matrix(responses, {
  columns: options.map(o => ({ id: o.id, header: slotLabel(o) })),
  valueAt: at.map("cells"),                                 // row.cells[colId]
  cell: cell.cycle(["yes","ifneedbe","no"], { styleSet: availability }),
  editableRow: "you",
  onCell: act.server("poll.toggleCell"),
})

// a timeline: one stream, a renderer per kind
ui.timeline(activities, {
  item: (a) => a.kind === "email" ? emailRow(a) : defaultRow(a),   // ← unit renderer switches on kind
  onOpen: act.navigate("/activity/${id}"),
})

// a typed record: fields described with the field.* builders, rendered read or edit
ui.record(contact, {
  sections: [
    { label: "Details", fields: [field.email("email"), field.phone("phone"), field.user("owner")] },
  ],
  mode: "read",
  onFieldChange: act.server("field.update"),
})
```

Compare the board verb to what an author writes today for the same screen: a nested
pair of loops over stages and deals, hand-built `component("...")` calls, manual key
and prop threading, and no drag wiring at all because there is no board type to target.
The verb collapses that to "here is the data, here is how one card looks, here is what
a move means." The *engine* — columns, drag, drop targets, selection — is supplied by
the verb and written once on the Go side.

### 4.4 The unit-renderer convention

The unit renderer is the redesign's load-bearing idea, so pin down its contract. A
unit renderer is a JavaScript function the runtime invokes with a **payload** and
whose return value must be a widget node (validated, exactly as `detailNode` validates
today at `module.go:882`). The payload mirrors the React contract for that engine:

```
  board.card(payload)     payload = { card, columnId, selected }
  matrix.cell(payload)    payload = { row, col, value, selected, editable }
  timeline.item(payload)  payload = { activity, isLast }
  record.field(payload)   payload = { spec, value, mode }   // usually you pass field.* specs instead
```

Two conveniences keep it ergonomic. First, the renderer may be given the raw datum
directly when the engine's payload is trivial (the board `card: (deal) => …` form is
sugar for `card: ({ card }) => …`). Second, where an author does not want a custom
unit at all, they pass a **spec instead of a function** — `cell: cell.status("stage")`
— and the runtime supplies a default renderer for that spec. Function for bespoke,
spec for standard: one axis, author's choice.

This convention is what makes the surface both composable and opinionated. Composable,
because the unit is an arbitrary function the author writes with the full nesting
vocabulary. Opinionated, because the *arrangement* around the unit — the thing that is
hard to get right — is not the author's to write; the verb owns it.

### 4.5 Before and after

The redesign's value is clearest side by side. The pipeline board today (sketch):

```js
// BEFORE — no board verb, so assemble raw
ui.panel({ title: "Sales" },
  ui.inline({},
    ...pipeline.stages.map(stage =>
      component("Panel", { title: stage.name, density: "condensed" },
        ...deals.filter(d => d.stageId === stage.id).map(d =>
          component("Panel", { tone: "raised" },
            component("Text", {}, text(d.title)),
            component("Caption", {}, text(fmtMoney(d.amount)))))))));
// …and dragging a card between columns is simply not expressible.
```

```js
// AFTER — one verb; drag is built in
ui.board(deals, {
  columns: pipeline.stages.map(s => ({ id: s.id, header: `${s.name} · ${fmtMoney(sum(s))}` })),
  columnOf: at("stageId"),
  card: (d) => ui.panel({ tone: "raised" }, ui.text(d.title), ui.caption(fmtMoney(d.amount))),
  onMove: act.server("deal.move"),
});
```

The "after" is shorter, but that is not the point — the point is that the "before"
cannot express drag at all, and the "after" gets it because the *engine* provides it.
Composition-first means the author writes only the card; opinionated means the author
cannot get the board layout wrong, because they do not write it.

---

## Part 5 — How this stays opinionated

Composability and opinionation pull in opposite directions in most systems; here they
are made to cooperate by *where* each lives. The author composes freely inside the
unit — the card, the cell, the timeline row — and the language composes freely at the
top — nesting verbs and panels. But the *arrangement between units* — the grid, the
board, the timeline spine, the field layout — is not composed by the author at all; it
is chosen by the verb. The surface is opinionated precisely at the layer that is hard
to get right and permissive precisely at the layer that is cheap to get wrong.

Three mechanisms keep it that way:

- **A small, curated verb set.** There are on the order of a dozen engine verbs, each
  corresponding to a real engine in the library. New verbs are added deliberately when
  a new engine lands, not invented per screen. The verb set *is* the set of arrangement
  patterns the product supports, and keeping it small keeps the product coherent.
- **The escape hatch is named and rare.** `raw()` (today's `component`) is always
  available, so nothing is impossible, but its name announces that you have left the
  curated path. A lint or review norm — "a script under review should contain no
  `raw()` except in a clearly justified spot" — turns the exception into a visible
  signal rather than a silent default.
- **Presets are declared, not coded.** The composite presets (`contactRecord`,
  `pipelineBoard`, `masterDetailTable`) are expressed as descriptors interpreted by one
  generic runner, so the library of opinionated starting points grows by data-entry.
  A descriptor names an engine verb, a data source, a unit renderer, and an action map:

```
recipe pipelineBoard:
  verb: board
  data: pipeline.deals
  columns: from pipeline.stages
  card: DealCard
  actions: { onMove: deal.move }
```

This is the mechanism `design-doc/01` Part 6 recommended for collapsing the hand-written
recipes (`mediaLibraryRecipe`, `articleListRecipe`) into data.

---

## Part 6 — Applying it backwards, without breaking scripts

The redesign is additive and can land incrementally; nothing here requires a flag day.
The migration has a clear order because each step stands on the previous.

First, **add the spec builders** (`at`, `style.by`, `cell.cycle`/`value`, `field.*`,
`select.*`). These are pure additions to the existing `cell`/`action` objects and break
nothing; they immediately close the "surface lags the IR" gap (Part 2 #3) for authors
who assemble with the substrate.

Second, **add the engine verbs** (`matrix`, `board`, `timeline`, `record`, `list`,
`calendar`) as new helpers that internally call the same `buildComponent` the substrate
uses, plus the unit-renderer machinery generalized from `detailNode` (`module.go:882`).
Existing scripts keep working; new scripts get the composable surface. This is the bulk
of the value and it touches no existing behavior.

Third, **rename `component` → `raw`** with `component` kept as a deprecated alias, and
introduce the `ui.when`/`ui.map`/`ui.slot` helpers. Cosmetic and non-breaking.

Fourth, **converge the two collection surfaces**. Make `data.collection` (the classic
grammar) build the same typed `CollectionSpec` the v2 path builds and lower through one
compiler (`design-doc/01` Part 6 #5), then express both as one `collection` verb in the
new style. The v2 fluent surface can remain as a deprecated alias during migration and
be removed once callers move. This is the only step that touches existing internals, and
it is the one that repays the most duplication.

Fifth, **make the surface generated**. Drive the helper maps, the `.d.ts`, and the recipe
descriptors from the manifest set (`design-doc/01` Part 5) so the verb set, the exported
functions, and the author-facing types agree by construction. At this point "apply it
backwards" is complete: every module — including the old `course.dsl`, `cms.dsl`,
`context_window.dsl` — exposes the same composition-first vocabulary, because the
vocabulary is generated from one description rather than hand-written per module.

The end state is that a script for a course studio, a CMS media library, a context
diagram, a scheduling poll, and a CRM board all read the same way: acquire a module,
call an engine verb, pass a unit renderer, done.

---

## Part 7 — Implementation notes (Go side)

For the implementer, the shape of the work in `pkg/widgetdsl`:

- **Unit-renderer machinery.** Generalize `detailNode` (`module.go:882`) into a shared
  `callUnitRenderer(fn, payload) → validated node` used by every engine verb. It already
  does the `goja.AssertFunction` + `isWidgetNodeExport` validation; lift it out of the
  master-detail recipe and make it the common path.
- **Engine verbs** are ordinary helpers registered like any other, except their factory
  wraps the unit-renderer call and builds the engine's props from the option bag. They
  live next to `componentFactory` (`module.go:565`).
- **Spec builders** extend `cellObject` (`module.go:272`) and `actionObject`
  (`module.go:318`); `style.by`, `at`, `field.*`, `select.*` are new sibling objects
  built the same way. `cell.cycle`/`cell.value` reuse `buildPaletteStyleSet`
  (`module.go:393`) for the `ContextStyleSet` a `StyleBySpec` needs.
- **Capability descriptors.** Replace the `if spec.name == …` chains in `runtime.install`
  (`module.go:236-269`) and the mirror in `typescript.go:73-116` with declarative fields
  on `moduleSpec` (`design-doc/01` Part 6 #4), so a module's verb set, spec objects, and
  recipes are data both the runtime and the `.d.ts` generator iterate.
- **Declarative recipes.** Interpret a `recipeSpec{ verb, data, unit, actionMap }`
  generically instead of hand-writing each recipe (Part 5).
- **Typed declarations.** Emit the real IR prop and spec types in the generated `.d.ts`
  (`design-doc/01` Part 6 #6) so authors get completion on the new verbs and specs; the
  runtime output stays plain JSON.

None of this changes what the DSL *outputs* — every verb still returns the same IR node
shape the frontend already renders. The redesign is entirely in the authoring surface
and the internal wiring.

---

## Part 8 — Open questions

1. **How much does the unit renderer see?** The proposal hands each unit renderer a
   payload mirroring the React contract. Should it also receive helpers (an `h`/`ui`
   handle) so the function is self-contained, or rely on module closure? Closure is
   simpler; an injected handle is more portable across module boundaries.
2. **Do we keep the grammar's *roles*, or move to the CRM *field types*?** The classic
   grammar names fields by behavioral role (`primary`, `measure`, `status`); the CRM kit
   names them by type (`currency`, `select`, `relation`). Converging on one vocabulary is
   desirable (this is the CRM ticket's open question too); the DSL redesign should adopt
   whichever wins.
3. **Fluent chains for the few cases that want them?** A small number of builders (a
   query/segment builder, perhaps) read naturally as chains. Is a chain ever worth
   breaking the "one authoring model" rule, or should everything be option bags?
4. **Verb granularity.** Is `list` distinct from `collection`, or is `collection` just
   `list` with a schema? Fewer verbs is more opinionated; the line wants deciding.
5. **Deprecation window** for `component`→`raw` and the v2 fluent surface — one release,
   or a longer coexistence?

## References

- `design-doc/01-widget-library-decomposition-analysis-and-design.md` — Part 4 (the specs
  to unify), Part 5 (manifest-as-source-of-truth), and Part 6 (the DSL analysis this
  document deepens).
- `reference/02-the-widget-system-a-new-intern-s-guide-...` — the IR and rendering machinery.
- `RAGEVAL-CRM-WIDGETS` — the CRM kit, whose `field.*` builders this redesign incorporates.
- `RAGEVAL-SCHEDULE-WIDGETS` — the scheduling engines (`matrix`, `board`-adjacent, `calendar`)
  the engine verbs expose.
- `pkg/widgetdsl/module.go` — `moduleSpec` (:23), `runtime.install` (:236), `componentFactory`
  (:565), `cellObject` (:272), `actionObject` (:318), `recipesObject` (:495),
  `masterDetailTableRecipe` (:853) and `detailNode` (:882, the unit-renderer prototype),
  `buildPaletteStyleSet` (:393); `grammar.go` (`collectionVerb` :265); `v2_builders.go`;
  `typescript.go` (`TypeScriptModule` :14).

