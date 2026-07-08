---
Title: 'Scheduling and Calendar Widgets: Implementation and IR Handoff for DSL Wiring'
Ticket: RAGEVAL-SCHEDULE-WIDGETS
Status: active
Topics:
    - ui-dsl
    - widget-ir
    - design-system
    - react
    - frontend-architecture
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/engines.ts
      Note: New IR specs (StyleBySpec, CycleCellSpec, ValueCellSpec, engine props)
    - Path: repo://packages/rag-evaluation-site/src/widgets/presets/scheduling.ts
      Note: Preset node shapes to port to Go recipes
    - Path: repo://pkg/widgetdsl/module.go
      Note: DSL runtime — helper maps, componentFactory, cellObject, recipes (edit sites for wiring)
    - Path: repo://pkg/widgetdsl/typescript.go
      Note: TS .d.ts codegen from helper maps (keep parity)
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-06T20:50:01.934195415-04:00
WhatFor: ""
WhenToUse: ""
---


# Scheduling and Calendar Widgets: Implementation and IR Handoff for DSL Wiring

> **Audience:** the Goja/`go-go-goja` engineer who will expose these widgets in
> the Widget IR DSL (`pkg/widgetdsl/`). It documents every widget I implemented
> on the React/TypeScript side, its exact IR contract, the actions it emits, and
> a precise, file-and-symbol-level plan for wiring it into the Go DSL runtime.
>
> **TL;DR of the split:** the TS package (`packages/rag-evaluation-site`) owns the
> React components, the IR adapters, the `.widget.yaml` manifests, and TS preset
> builders. The Go package (`pkg/widgetdsl`) owns the DSL modules that emit IR
> JSON from user scripts. The two meet at the Widget IR node shape
> (`{ kind:"component", type, props, children }`). **All the TS side is done and
> green; none of it is wired into `pkg/widgetdsl` yet** — that is your task.

## Goal

Give you everything needed to add `MatrixGrid`, `SegmentedBar`, `MonthGrid`,
`TimeGrid` (and, optionally, the scheduling/calendar *preset recipes*) to the
Widget IR DSL, without having to reverse-engineer the TypeScript. Every IR prop,
every action payload, and every Go edit site is enumerated below.

## Context — how the two sides connect

```
User script (Goja/JS)                         pkg/widgetdsl/module.go
  require("data.dsl").matrixGrid({ ... })  ─┐  (componentFactory → buildComponent)
                                            │
                                            ▼  emits IR JSON
   { kind:"component", type:"MatrixGrid", props:{...}, children:[] }
                                            │  serialized to the browser
                                            ▼
  WidgetRenderer.registry.get("MatrixGrid")     packages/rag-evaluation-site
   → matrixGridWidget.render(props, _, ctx)      src/components/molecules/MatrixGrid/
   → <MatrixGrid …/> (React)                      MatrixGrid.widget.tsx → MatrixGrid.tsx
```

Two facts that drive everything:

1. **Component `type` is open-ended.** `pkg/widgetdsl` and the TS renderer both
   treat `props` as free-form JSON (`typescript.go` says so explicitly; `v2/spec`
   only validates node *kinds*, action *kinds*, and cell *field kinds* — **not**
   component types). So exposing a new engine is mostly: pick a helper name, map
   it to the `type` string. No schema registration of the type is required.
2. **The helper→type map auto-generates the DSL function.**
   `runtime.install` (module.go ~L234) iterates `spec.helpers` and calls
   `componentFactory(componentType)` (module.go ~L565), which returns a goja fn
   `props → buildComponent(type, props, children)`. So **adding one map entry
   creates the whole `matrixGrid(props, …children)` helper.**

---

## Part 1 — Where everything lives (file index)

### TypeScript side (`packages/rag-evaluation-site/src/`)

| Concern | Path | Key symbols |
|---|---|---|
| IR node model | `widgets/ir/core.ts` | `WidgetNode`, `ComponentNode`, `RagWidgetType`, `BaseWidgetProps`, `text()`, `element()`, `component()` |
| IR action specs | `widgets/ir/actions.ts` | `ActionSpec`, `ServerActionSpec`, `PayloadTemplateSpec`, `TemplatePartSpec` |
| IR cell specs | `widgets/ir/cells.ts` | `CellSpec`, `RowKeySpec`, `DataTableColumnSpec` |
| **IR engine specs** | `widgets/ir/engines.ts` | `StyleBySpec`, `CycleCellSpec`, `ValueCellSpec`, `MatrixCellSpec`, `MatrixColumnWidgetSpec`, `MatrixValueSpec`, `MatrixGridWidgetProps`, `SegmentedBarWidgetProps`, `MonthGridWidgetProps`, `TimeGridWidgetProps` |
| IR component props (rest) + `WidgetProps` union | `widgets/ir/props.ts` | all other `*WidgetProps` |
| IR barrel | `widgets/ir/index.ts` | re-exports all of the above |
| Renderer | `widgets/WidgetRenderer.tsx` | `WidgetRenderer`, `renderComponentNode` |
| Adapter/registry infra | `widgets/registry.ts` | `defineWidget`, `RenderContext`, `createWidgetRegistry`, `mergeWidgetRegistries`, `WidgetModule` |
| Registry wiring | `widgets/defaultRegistry.ts` | `uiWidgetRegistry`, `dataWidgetRegistry`, `timeWidgetRegistry`, `defaultWidgetRegistry` |
| Cell renderer interp | `widgets/cellRenderers.tsx` | `renderCell`, `rowKey` |
| **StyleBySpec resolver** | `widgets/styleBy.ts` | `resolveStyleByVars` |
| **Presets (TS)** | `widgets/presets/scheduling.ts` | `availabilityMatrix`, `pollResults`, `monthCalendar`, `weekCalendar` |
| Domain DTOs | `scheduling/types.ts` | `MeetingPoll`, `ParticipantResponse`, `SlotTally`, `CalendarEvent`, `BookingType`, `BookableDay`, `BookableSlot`, `AvailabilityState` |
| Domain palettes | `scheduling/palettes.ts` | `availabilityStyleSet`, `AVAILABILITY_STATES`, `AVAILABILITY_GLYPHS`, `eventStyleSet` |
| Domain fixtures | `scheduling/fixtures.ts` | `sampleTeamSyncPoll`, `sampleTeamSyncTallies`, `sampleWeekEvents`, `sampleBooking*` |
| IR renderer stories | `widgets/WidgetRenderer.scheduling.stories.tsx`, `widgets/WidgetRenderer.calendar.stories.tsx` | reference IR node examples |

### Go DSL side (`pkg/widgetdsl/`)

| Concern | Path | Key symbols |
|---|---|---|
| Module + helper maps | `module.go` | `uiHelpers`, `dataHelpers`, `contextWindowHelpers`, `cmsHelpers`, `courseHelpers`, `moduleSpecs`, `runtime.install`, `componentFactory`, `cellObject`, `actionObject`, `recipesObject` |
| Engine registrar | `registrar.go` | `Registrar.RegisterRuntimeModule`, `Register` |
| TS decl codegen | `typescript.go` | `TypeScriptModule` |
| TS parity test | `typescript_fixture_test.go`, `typescript_test.go` | compiles generated `.d.ts` against the TS package's `tsc` |
| v2 typed model | `v2/spec/{types,validate,lower}.go` | `NodeSpec`, cell/action-kind validation |
| Grammar | `grammar.go` | DSL grammar |

---

## Part 2 — New IR types (all in `widgets/ir/engines.ts`)

These are the reusable spec vocabulary the engines consume. They are the pieces
that need matching DSL builder ergonomics.

### `StyleBySpec` — defunctionalized color function
```ts
interface StyleBySpec {
  field?: string;                  // datum field to key on; default = the cell value
  styleSet: ContextStyleSet;       // styleKey -> ContextVisualStyle (existing palette contract)
  map?: Record<string, string>;    // optional value -> styleKey remap before lookup
  fallbackStyleKey?: string;
}
```
Resolved by `resolveStyleByVars(spec, value, row?)` in `widgets/styleBy.ts` →
CSS vars (`--ctx-fill` etc.) via `contextVisualStyleToCssVars`. **DSL note:** this
is a plain JSON prop; no special builder is required, but a convenience
`styleBy(styleSet, opts)` helper would be nice-to-have. It is consumed today only
by `MatrixGrid.colorBy` (see `ValueCellSpec`).

### Cell specs for `MatrixGrid`
```ts
interface CycleCellSpec {                    // n-state toggle cell (availability, RSVP…)
  kind: "cycle";
  states: string[];                          // ordered ring; click advances
  glyphs?: Record<string, RenderableValue>;  // state -> glyph
  styleSet?: ContextStyleSet;                // palette override; else grid.styleSet
}
interface ValueCellSpec { kind: "value"; }   // render resolved (row,col) value as text (+ colorBy)
type MatrixCellSpec = CellSpec | CycleCellSpec | ValueCellSpec;   // CellSpec = the DataTable cell union
```

### `MatrixValueSpec` — accessor for value at (row, col)
```ts
type MatrixValueSpec = { mapField: string } | { template: string };
// mapField: value = row[mapField][col.id]      (e.g. { mapField: "cells" })
// template: interpolate ${field}/${colId} against the row
// (omitted): value = row[col.id]
```

### `MatrixColumnWidgetSpec`
```ts
interface MatrixColumnWidgetSpec { id: string; header: RenderableValue; meta?: JsonObject; }
```

---

## Part 3 — Per-widget reference

For each: React component, IR adapter, manifest binding, full IR props, emitted
action context, and a copy-paste IR node example.

### 3.1 `MatrixGrid` — generic rows×columns grid engine (flagship)

- **React:** `components/molecules/MatrixGrid/MatrixGrid.tsx` → `MatrixGrid`,
  types `MatrixColumnSpec`, `MatrixCellPayload`, `MatrixRow`.
- **Adapter:** `components/molecules/MatrixGrid/MatrixGrid.widget.tsx` → export
  `matrixGridWidget` (`defineWidget`, `type:"MatrixGrid"`, `module:"data.dsl"`).
- **Manifest:** `MatrixGrid.widget.yaml` → `module: data.dsl`, `helper: matrixGrid`.
- **Registered in:** `dataWidgetRegistry` (`defaultRegistry.ts`).

**IR props (`MatrixGridWidgetProps`, `engines.ts`):**

| Prop | Type | Notes |
|---|---|---|
| `rows` | `JsonObject[]` | opaque row data |
| `columns` | `MatrixColumnWidgetSpec[]` | `{ id, header, meta? }` |
| `valueAt?` | `MatrixValueSpec` | default `row[col.id]` |
| `cell?` | `MatrixCellSpec` | Mode A renderer (`cycle` / `value` / `CellSpec`) |
| `cells?` | `WidgetNode[][]` | Mode B: explicit `[rowIndex][colIndex]` node matrix |
| `rowHeader?` | `CellSpec` | left header cell per row |
| `styleSet?` | `ContextStyleSet` | palette for `cycle` cells |
| `colorBy?` | `StyleBySpec` | tints `value` cells |
| `footer?` | `{ header?: RenderableValue; cell: CellSpec }` | footer `cell` is evaluated against each **column's `meta`** as the row |
| `getRowKey?` | `RowKeySpec` | default `row.id` |
| `editableRowKey?` | `string` | only this row's cells are interactive |
| `selectedCell?` | `{ rowKey, colId }` | selection outline |
| `cornerCell?` | `RenderableValue` | top-left cell |
| `stickyHeader?` | `boolean` | default true |
| `ariaLabel?` | `string` | |
| `onCellAction?` | `ActionSpec` | fired on cell change |

**Emitted action context** (`onCellAction`): `{ rowKey, colId, value, componentType:"MatrixGrid" }`.
Compose the server payload with template parts, e.g.
`{ pollId, responseId:{kind:"path",path:"rowKey"}, optionId:{kind:"path",path:"colId"}, state:{kind:"path",path:"value"} }`.

**Example IR node (a Doodle poll):**
```json
{ "kind":"component", "type":"MatrixGrid", "props":{
  "rows":[{"id":"alice","name":"Alice","cells":{"s1":"yes","s2":"no"}}],
  "columns":[{"id":"s1","header":{"kind":"text","text":"Thu 14:00"},"meta":{"yes":4,"total":6}}],
  "valueAt":{"mapField":"cells"},
  "cell":{"kind":"cycle","states":["yes","ifneedbe","no","unknown"],"glyphs":{"yes":"✓","no":"✕"}},
  "styleSet":{ "...ContextStyleSet..." },
  "rowHeader":{"kind":"field","field":"name"},
  "editableRowKey":"you",
  "footer":{"header":{"kind":"text","text":"yes"},"cell":{"kind":"template","template":"${yes}/${total}"}},
  "onCellAction":{"kind":"server","name":"poll.toggleCell","payload":{ "...template parts..." }}
}}
```

### 3.2 `SegmentedBar` — proportional segmented bar

- **React:** `components/molecules/SegmentedBar/SegmentedBar.tsx` → `SegmentedBar`,
  types `SegmentedBarSegment`, `SegmentedBarMarker`.
- **Adapter:** `SegmentedBar.widget.tsx` → `segmentedBarWidget` (`type:"SegmentedBar"`, `module:"ui.dsl"`).
- **Manifest:** `module: ui.dsl`, `helper: segmentedBar`.
- **Registered in:** `uiWidgetRegistry`.

**IR props (`SegmentedBarWidgetProps`):**

| Prop | Type |
|---|---|
| `segments` | `{ value:number; styleKey:string; label?:RenderableValue }[]` |
| `styleSet` | `ContextStyleSet` |
| `total?` | `number` (default = Σ values) |
| `showCounts?` | `boolean` |
| `markers?` | `{ at:number; styleKey?:string; label?:RenderableValue }[]` |
| `size?` | `"sm"｜"md"｜"lg"` |
| `onSegmentAction?` | `ActionSpec` |

**Emitted action context:** `{ styleKey, index, value:styleKey, componentType:"SegmentedBar" }`.

### 3.3 `MonthGrid` — month-of-days engine

- **React:** `components/molecules/MonthGrid/MonthGrid.tsx` → `MonthGrid`,
  types `MonthGridDayPayload`, `MonthGridDayMarker`.
- **Adapter:** `MonthGrid.widget.tsx` → `monthGridWidget` (`type:"MonthGrid"`, `module:"time.dsl"`).
- **Manifest:** `module: time.dsl`, `helper: monthGrid`. **⚠ `time.dsl` does not
  exist in `pkg/widgetdsl` yet — see Part 5.**
- **Registered in:** `timeWidgetRegistry`.

**IR props (`MonthGridWidgetProps`):**

| Prop | Type |
|---|---|
| `monthISO` | `string` (`2026-07` or full ISO) |
| `markers?` | `Record<dateISO, { count?:number; styleKey?:string; label?:RenderableValue }>` |
| `styleSet?` | `ContextStyleSet` |
| `selectedDateISO?` | `string` |
| `todayISO?` | `string` (omit → no today highlight; the engine never calls `Date.now`) |
| `minDateISO?` / `maxDateISO?` | `string` |
| `weekStartsOn?` | `0｜1` (default 1) |
| `showHeader?` | `boolean` |
| `onDaySelectAction?` | `ActionSpec` → ctx `{ dateISO, value:dateISO }` |
| `onMonthChangeAction?` | `ActionSpec` → ctx `{ monthISO, value:monthISO }` |

### 3.4 `TimeGrid` — week/day time-grid engine

- **React:** `components/molecules/TimeGrid/TimeGrid.tsx` → `TimeGrid`,
  types `TimeGridBlock`, `TimeGridColumnSpec`, `TimeGridBlockPayload`,
  helper `packColumn` (lane packing for overlaps).
- **Adapter:** `TimeGrid.widget.tsx` → `timeGridWidget` (`type:"TimeGrid"`, `module:"time.dsl"`).
- **Manifest:** `module: time.dsl`, `helper: timeGrid`. **⚠ same `time.dsl` note.**
- **Registered in:** `timeWidgetRegistry`.

**IR props (`TimeGridWidgetProps`):**

| Prop | Type |
|---|---|
| `days` | `Array<string ｜ { dayISO:string; header?:RenderableValue }>` |
| `blocks` | `{ id; dayISO; startISO; endISO; styleKey; label:RenderableValue; allDay?; meta? }[]` |
| `styleSet` | `ContextStyleSet` |
| `hourStart?` / `hourEnd?` | `number` (default 8 / 20) |
| `hourHeight?` | `number` px/hour |
| `nowISO?` | `string` (omit → no now-line) |
| `selectedBlockId?` | `string` |
| `onBlockSelectAction?` | `ActionSpec` → ctx `{ blockId, value:blockId }` |
| `onSlotCreateAction?` | `ActionSpec` → ctx `{ dayISO, hour, value:dayISO }` |

> **Positioning note:** the engine reads only the wall-clock `HH:MM` substring of
> `startISO`/`endISO`; it does **no** timezone conversion. Any UTC→local
> conversion belongs in the preset/producer, not the widget.

### 3.5 Atoms (not standalone IR widget types)

`DateTile`, `RatioBadge`, `CycleCell` (`components/atoms/…`) are React atoms used
*inside* engines/organisms. Only `CycleCell` has an IR surface, and it is exposed
**via `MatrixGrid`'s `cell:{kind:"cycle"}` spec**, not as its own `type`. Do not
add DSL helpers for these unless you want them as free-standing nodes.

### 3.6 Organisms (built, not yet IR-exposed)

`MeetingPollPanel`, `PollResultsPanel`, `CalendarMonthPanel`, `CalendarWeekPanel`,
`BookingPagePanel` (`components/organisms/…`) are React-first and have **no
`.widget.tsx` adapter yet** (per "React first, Widget IR later"). They are out of
scope for DSL wiring now; the equivalent DSL surface is the **preset recipes**
below.

---

## Part 4 — Presets (TS) and what IR they emit

`widgets/presets/scheduling.ts`. These are the model for the DSL **recipes** you
may implement (Part 5, step 5). Each returns a `WidgetNode`.

| Preset | Signature | Emits |
|---|---|---|
| `availabilityMatrix` | `(poll, { tallies?, editableResponseId? })` | one `MatrixGrid` node (cycle cells, availability palette, tally footer, `poll.toggleCell` action) |
| `pollResults` | `(poll, tallies)` | `Stack` of `{ Caption + SegmentedBar }` per option |
| `monthCalendar` | `(events, monthISO)` | one `MonthGrid` node with density markers colored by `event.colorKey` |
| `weekCalendar` | `(events, daysISO)` | one `TimeGrid` node with blocks from events |

Read these four functions as the exact prop shapes to reproduce in Go.

---

## Part 5 — DSL wiring plan (`pkg/widgetdsl`) — the actual handoff

All edits are in `pkg/widgetdsl/`. Steps 1–2 make the **engines** callable; steps
3–4 add the **cell/style builders**; step 5 adds the **preset recipes**; step 6
keeps codegen/tests green.

### Step 1 — Register the engine helpers (`module.go`)

- `dataHelpers` (module.go ~L82): add `"matrixGrid": "MatrixGrid"`.
- `uiHelpers` (module.go ~L35): add `"segmentedBar": "SegmentedBar"`.

That alone makes `require("data.dsl").matrixGrid({...})` and
`require("ui.dsl").segmentedBar({...})` emit correct IR, because
`runtime.install` (~L234) wires each map entry through `componentFactory` (~L565).

### Step 2 — Resolve `time.dsl` and register `monthGrid`/`timeGrid`

The TS manifests bind `MonthGrid`/`TimeGrid` to a `time.dsl` module that **does
not exist** in `module.go` (only `ui/data/data.v2/context_window/course/cms`).
Pick one:

- **(Recommended) Add a `time.dsl` module.** In `module.go`:
  ```go
  const TimeModuleName = "time.dsl"
  var timeHelpers = map[string]string{ "monthGrid": "MonthGrid", "timeGrid": "TimeGrid" }
  // append to moduleSpecs:
  { name: TimeModuleName, helpers: timeHelpers, action: true,
    doc: "time.dsl provides calendar month/week time-grid helpers." }
  ```
  `Register` (~L207) and `init` (~L225) iterate `moduleSpecs`, so registration is
  automatic. Also add it in `registrar.go` if that enumerates modules explicitly
  (it calls `Register(reg)`, so likely no change).
- **(Alternative) Drop `time.dsl`:** put `monthGrid`/`timeGrid` in `dataHelpers`
  or `uiHelpers` and change the two TS `.widget.yaml` `module:` fields to match.

Whichever you choose, **the TS `.widget.yaml` `module:` and the Go module must
agree** — that is the one inconsistency this handoff introduces.

### Step 3 — Add the `cycle` and `value` cell builders (`module.go` `cellObject`)

`cellObject()` (~L272) builds the `cell.*` object for `data.dsl`. Add:
```go
setExport(cell, "cycle", func(states []string, options ...goja.Value) map[string]any {
    out := map[string]any{"kind": "cycle", "states": states}
    mergeOptions(out, exportOptions(options))   // glyphs, styleSet
    return out
})
setExport(cell, "value", func() map[string]any { return map[string]any{"kind": "value"} })
```
These mirror `CycleCellSpec` / `ValueCellSpec` exactly. (Existing `cell.field`,
`cell.template`, etc. already cover the `CellSpec` branch of `MatrixCellSpec`.)

### Step 4 — (Optional) `styleBy` convenience + `colorBy` passthrough

`StyleBySpec` is a plain prop, so `matrixGrid({ colorBy: { styleSet, map } })`
already works. Optionally add a `styleBy(styleSet, options)` helper (ui.dsl) that
returns `{ styleSet, field?, map?, fallbackStyleKey? }` for symmetry with `cell`.

### Step 5 — Preset recipes (`schedule.dsl` / `calendar.dsl`)

The presets are **composite** (they build a configured node), so they map to the
existing **recipe** mechanism, not the helper map. See `masterDetailTableRecipe`
(module.go ~L853) and `recipesObject` (~L495) for the pattern. Recommended:

- Add `ScheduleModuleName = "schedule.dsl"` with `recipes: ["availabilityMatrix","pollResults"]`.
- Add `CalendarModuleName = "calendar.dsl"` with `recipes: ["monthCalendar","weekCalendar"]`.
- Implement `r.availabilityMatrixRecipe`, `r.pollResultsRecipe`, etc. as
  `func(goja.FunctionCall) goja.Value` that reproduce `presets/scheduling.ts`
  (build the same `MatrixGrid`/`Stack`/`MonthGrid`/`TimeGrid` node maps). Wire
  them into `recipesObject` (~L495 switch).

These recipes need the availability/event palettes; port `scheduling/palettes.ts`
(`availabilityStyleSet`, `AVAILABILITY_STATES`, `AVAILABILITY_GLYPHS`,
`eventStyleSet`) to Go constants, or accept the styleSet as an argument.

### Step 6 — Keep codegen + tests green

- `typescript.go` `TypeScriptModule` (~L14) generates `.d.ts` **from the helper
  maps**, so new helpers auto-appear. Recipes and the new `cell.cycle/value`
  builders may need manual decl lines (check how existing `cell`/`recipes` decls
  are emitted).
- `typescript_fixture_test.go` compiles the generated `.d.ts` against the TS
  package's `tsc` (`packages/rag-evaluation-site/node_modules/.bin/tsc`) — run
  `go test ./pkg/widgetdsl/...` after edits.
- `v2/spec/validate.go` validates **cell field kinds** (~L112) and action kinds
  (~L185) but not component types. If `cycle`/`value` cells will be authored via
  the **v2** path, add them to the field-kind allowlist there; the classic
  `data.dsl` path is open-map and needs no change.

### Wiring summary table

| TS widget `type` | TS module (`.widget.yaml`) | Go helper to add | Go map / site |
|---|---|---|---|
| `MatrixGrid` | `data.dsl` | `matrixGrid` | `dataHelpers` (module.go ~L82) |
| `SegmentedBar` | `ui.dsl` | `segmentedBar` | `uiHelpers` (module.go ~L35) |
| `MonthGrid` | `time.dsl` ⚠ | `monthGrid` | new `timeHelpers` (+ `time.dsl` module) |
| `TimeGrid` | `time.dsl` ⚠ | `timeGrid` | new `timeHelpers` |
| cycle cell | (cell spec) | `cell.cycle` | `cellObject` (module.go ~L272) |
| value cell | (cell spec) | `cell.value` | `cellObject` |
| `availabilityMatrix` | preset | `recipes.availabilityMatrix` | new `schedule.dsl` recipe |
| `pollResults` | preset | `recipes.pollResults` | new `schedule.dsl` recipe |
| `monthCalendar` | preset | `recipes.monthCalendar` | new `calendar.dsl` recipe |
| `weekCalendar` | preset | `recipes.weekCalendar` | new `calendar.dsl` recipe |

---

## Part 6 — Actions / events contract (for the server side)

Every interactive widget emits an `ActionSpec` with a context object. The server
action dispatcher (envelope `POST /api/widget/actions/<name>` →
`{ ok, refresh?, toast?, patch?, data? }`, defined client-side in
`widgets/actions.ts` `dispatchWidgetAction`) is **not implemented in this repo**
(a grep of the Go tree found no `/api/polls` or `widget/actions` handler). The
proposed `name`s and payloads are in the design doc
(`design-doc/01-…-implementation-guide.md`, Part D): `poll.toggleCell`,
`poll.submitResponse`, `poll.finalize`, `booking.daySlots`, `booking.book`,
`calendar.range`, etc. Context keys emitted per widget are listed in Part 3.

---

## Part 7 — Known inconsistencies / open items

1. **`time.dsl` mismatch** (Part 5, step 2) — must be reconciled.
2. **Presets not yet in the DSL** — only reachable as TS today; recipes are the
   Go equivalent (Part 5, step 5).
3. **`StyleBySpec` has one consumer** (`MatrixGrid.colorBy` + `ValueCellSpec`).
   If you want month/segment coloring driven by `colorBy` too, that is a small
   follow-up in the respective adapters (not yet done).
4. **Organism-level IR** (`MeetingPollPanel`, …) is deliberately deferred.
5. **Namespaced types** (`schedule/AvailabilityMatrix`) are **not** used — all
   types are bare and globally unique. Keep new DSL types unique to avoid the
   `createWidgetRegistry` duplicate-`type` panic.

## Verification performed on the TS side

- `pnpm --dir packages/rag-evaluation-site typecheck` — passes.
- `pnpm --dir packages/rag-evaluation-site build-storybook` — passes.
- `pnpm --dir packages/rag-evaluation-site build` — passes (library packaging).

## Related

- `design-doc/01-calendar-and-scheduling-widgets-analysis-design-and-implementation-guide.md` — the design + backend spec.
- `reference/01-implementation-diary.md` — chronological build log (Steps 1–8).
- `pkg/widgetdsl/module.go` — the DSL runtime you will edit.
- `packages/rag-evaluation-site/src/widgets/ir/engines.ts` — the new IR specs.
- `packages/rag-evaluation-site/src/widgets/presets/scheduling.ts` — preset shapes to port.
