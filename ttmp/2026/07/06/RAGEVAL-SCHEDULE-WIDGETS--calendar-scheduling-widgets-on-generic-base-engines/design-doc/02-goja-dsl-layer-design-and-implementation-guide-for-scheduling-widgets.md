---
Title: Goja DSL Layer Design and Implementation Guide for Scheduling Widgets
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
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/engines.ts
      Note: Browser-side IR contracts the Goja DSL must emit
    - Path: repo://packages/rag-evaluation-site/src/widgets/presets/scheduling.ts
      Note: TypeScript preset shapes to mirror in Go recipes
    - Path: repo://pkg/widgetdsl/grammar.go
      Note: Existing intent-level data.dsl grammar used as comparison for scheduling recipes
    - Path: repo://pkg/widgetdsl/module.go
      Note: Main runtime edit site for module constants, helper maps, cell/style helpers, and recipes
    - Path: repo://pkg/widgetdsl/module_test.go
      Note: Runtime export and JSON-shape test patterns for new modules/helpers/recipes
    - Path: repo://pkg/widgetdsl/typescript.go
      Note: Generated TypeScript declarations that must match new runtime exports
    - Path: repo://pkg/widgetdsl/typescript_fixture_test.go
      Note: TypeScript fixture pattern for declaration parity
ExternalSources: []
Summary: Intern-facing design and implementation guide for wiring the new scheduling/calendar Widget IR engines into pkg/widgetdsl Goja modules.
LastUpdated: 2026-07-06T21:16:01.067156973-04:00
WhatFor: Use this to implement MatrixGrid, SegmentedBar, MonthGrid, TimeGrid, scheduling presets, calendar presets, cells, palettes, and TypeScript declarations in the Goja Widget DSL runtime.
WhenToUse: Read before editing pkg/widgetdsl for RAGEVAL-SCHEDULE-WIDGETS, especially when adding time.dsl, schedule.dsl, calendar.dsl, or new recipe helpers.
---


# Goja DSL Layer Design and Implementation Guide for Scheduling Widgets

> **Audience:** a new intern or Goja engineer who understands JavaScript but has
> not worked in this repository's Widget DSL runtime before.
>
> **Goal:** make the React/TypeScript scheduling widgets usable from Goja scripts
> through `require("*.dsl")`, while keeping the DSL output as plain Widget IR JSON
> that the existing `WidgetRenderer` already knows how to render.

---

## 0. Executive summary

The TypeScript side already has the scheduling/calendar widgets:

- `MatrixGrid` in `data.dsl` terms.
- `SegmentedBar` in `ui.dsl` terms.
- `MonthGrid` and `TimeGrid` in a new `time.dsl` module.
- TS presets for scheduling/calendar use cases: `availabilityMatrix`,
  `pollResults`, `monthCalendar`, `weekCalendar`.
- IR adapters and manifests that define the browser-side contracts.

The Goja side is not wired yet. The update should be intentionally small at the
engine layer and slightly richer at the preset layer:

1. **Expose engine helpers by map entries.** Add `matrixGrid`, `segmentedBar`,
   `monthGrid`, and `timeGrid` to the right helper maps. These helpers already
   compile to `{ kind:"component", type, props, children }` via the existing
   `componentFactory` mechanism.
2. **Add `time.dsl`.** The TS manifests for `MonthGrid` and `TimeGrid` already
   say `module: time.dsl`; the Go runtime should match that.
3. **Extend `data.dsl.cell`.** Add `cell.cycle(states, options?)` and
   `cell.value()` so `MatrixGrid` can be authored ergonomically.
4. **Add a generic `styleBy` convenience helper.** It should produce the
   `StyleBySpec` JSON consumed by `MatrixGrid.colorBy`.
5. **Add domain recipe modules.** Add `schedule.dsl` and `calendar.dsl` with
   recipe helpers that emit the same IR shapes as the TS presets.
6. **Update TypeScript declarations and tests.** `TypeScriptModule` mostly
   derives helper declarations automatically, but the new cell/style/preset
   surface needs explicit declaration coverage and runtime/fixture tests.

The most important principle: **the Goja DSL should not render React.** It only
constructs serializable Widget IR. The React package remains the interpreter.

---

# Part A — The system you are extending

## A.1 The full data flow

The Widget DSL is a way for JavaScript running inside Goja to author UI as JSON.
That JSON is sent to the browser, where the React `WidgetRenderer` looks up
component adapters and renders real React components.

```
Goja script                         pkg/widgetdsl                  Browser TS/React
────────────────────────────────────────────────────────────────────────────────────
const data = require("data.dsl")

const node = data.matrixGrid({        module.go helpers             WidgetRenderer
  rows, columns, cell, ...       ──▶  componentFactory        ──▶   registry.get("MatrixGrid")
})                                     buildComponent                matrixGridWidget.render(...)
                                      emits JSON                    <MatrixGrid />

Resulting Widget IR:
{
  "kind": "component",
  "type": "MatrixGrid",
  "props": { ... },
  "children": []
}
```

The browser does **not** know or care whether the node came from `data.dsl`,
`time.dsl`, `schedule.dsl`, or hand-written JSON. It only cares that the
component `type` has a registered adapter.

## A.2 Where the Goja DSL runtime lives

The Go runtime lives under `pkg/widgetdsl/`:

| File | What it owns | Why it matters for this task |
|---|---|---|
| `pkg/widgetdsl/module.go` | Module names, helper maps, module registration, generic node builders, cell/action helpers, recipes | Main file to edit for new modules/helpers/recipes |
| `pkg/widgetdsl/typescript.go` | Generated `.d.ts` declarations for Goja DSL modules | Must keep JS authoring types in sync with runtime exports |
| `pkg/widgetdsl/module_test.go` | Runtime export and JSON-serializability tests | Add tests that new modules and helpers exist |
| `pkg/widgetdsl/grammar.go` | Existing `data.dsl` higher-level grammar (`schema`, `record`, `collection`) | Useful model for intent-level DSL helpers, but not the primary edit site |
| `pkg/widgetdsl/grammar_test.go` | Tests for the old data grammar | Pattern for verifying emitted IR shape |
| `pkg/widgetdsl/v2_builders.go` | Typed fluent builder experiment for `data.v2.dsl` | Useful design inspiration; do not mix scheduling into it in phase 1 |
| `pkg/widgetdsl/v2/spec/*` | Typed v2 model, validation, lowering | Future direction if we later want a typed `schedule.v2.dsl` |

## A.3 The current runtime extension points

### Module specs

`module.go` defines module names and `moduleSpec` at the top:

```go
const (
    UIModuleName            = "ui.dsl"
    DataModuleName          = "data.dsl"
    DataV2ModuleName        = "data.v2.dsl"
    ContextWindowModuleName = "context_window.dsl"
    CourseModuleName        = "course.dsl"
    CmsModuleName           = "cms.dsl"
)

type moduleSpec struct {
    name    string
    doc     string
    helpers map[string]string
    page    bool
    cell    bool
    action  bool
    recipes []string
}
```

Each `moduleSpec` can expose:

- generic node constructors: `text`, `element`, `component`, `fragment`;
- component helper functions from `helpers`;
- `page` for `ui.dsl`;
- `cell` helper object for `data.dsl`;
- `action` helper object for modules that need event/server actions;
- `recipes` helper object for composed widget trees.

### Helper maps

A helper map is a map from JS helper name to Widget IR component type:

```go
var dataHelpers = map[string]string{
    "dataTable": "DataTable",
}
```

At runtime, `install` iterates the map and exports a function for each entry:

```go
for name, componentType := range spec.helpers {
    setExport(exports, name, r.componentFactory(componentType))
}
```

`componentFactory` is intentionally generic:

```go
func (r *runtime) componentFactory(componentType string) func(goja.FunctionCall) goja.Value {
    return func(call goja.FunctionCall) goja.Value {
        props, childStart := propsAndChildStart(call.Arguments, 0)
        return r.vm.ToValue(r.buildComponent(componentType, props, call.Arguments[childStart:]))
    }
}
```

So adding a simple engine helper is usually just a map entry.

### Cell helpers

`data.dsl` has `cellObject()`:

```go
setExport(cell, "field", func(field string, options ...goja.Value) map[string]any { ... })
setExport(cell, "template", func(template string) map[string]any { ... })
setExport(cell, "actionButton", func(label goja.Value, action goja.Value, options ...goja.Value) map[string]any { ... })
```

This returns serializable `CellSpec` JSON consumed by the TS `renderCell` switch.
`MatrixGrid` now introduces two additional cell specs:

- `{ kind:"cycle", states:[...], glyphs?, styleSet? }`
- `{ kind:"value" }`

Those belong here.

### Action helpers

Most modules expose `actionObject()`:

```go
action.server(name, options?)   -> { kind:"server", name, ...options }
action.navigate(to, options?)   -> { kind:"navigate", to, ...options }
action.download(to, options?)   -> { kind:"download", to, ...options }
action.event(event, options?)   -> { kind:"event", event, ...options }
action.copy(value)              -> { kind:"copy", value }
```

The scheduling widgets use normal `ActionSpec`s; no new action kind is required.
The browser resolves payload path specs in `packages/rag-evaluation-site/src/widgets/actions.ts`.

### Recipes

Recipes are composite builders. They do not map one helper to one component; they
assemble a subtree. Existing examples live in `recipesObject` and functions like
`metricsRecipe`, `actionToolbarRecipe`, and `masterDetailTableRecipe`.

The pattern is:

```go
func (r *runtime) recipesObject(names []string) *goja.Object {
    recipes := r.vm.NewObject()
    for _, name := range names {
        switch name {
        case "metrics":
            setExport(recipes, name, r.metricsRecipe)
        }
    }
    return recipes
}
```

Scheduling presets should follow this pattern first. A direct
`schedule.availabilityMatrix(...)` alias can be added later if ergonomics demand
it, but `schedule.recipes.availabilityMatrix(...)` matches existing module style.

---

# Part B — The browser contracts the Go DSL must emit

The Go DSL output must match the TypeScript contracts in
`packages/rag-evaluation-site/src/widgets/ir/engines.ts` and the adapters under
`packages/rag-evaluation-site/src/components/molecules/*/*.widget.tsx`.

## B.1 Engine helper table

| Browser widget type | Browser module in manifest | Goja module | JS helper | Implementation site |
|---|---|---|---|---|
| `MatrixGrid` | `data.dsl` | `data.dsl` | `data.matrixGrid(props, ...children)` | add to `dataHelpers` |
| `SegmentedBar` | `ui.dsl` | `ui.dsl` | `ui.segmentedBar(props)` | add to `uiHelpers` |
| `MonthGrid` | `time.dsl` | `time.dsl` | `time.monthGrid(props)` | new `timeHelpers` + `TimeModuleName` |
| `TimeGrid` | `time.dsl` | `time.dsl` | `time.timeGrid(props)` | new `timeHelpers` + `TimeModuleName` |

## B.2 Cell helper table

| Cell helper | Emits | Used by |
|---|---|---|
| `data.cell.cycle(states, options?)` | `{ kind:"cycle", states, ...options }` | `MatrixGrid.cell` |
| `data.cell.value()` | `{ kind:"value" }` | `MatrixGrid.cell`, especially with `colorBy` |
| existing `data.cell.field/number/status/caption/template/...` | existing `CellSpec` | `DataTable` and `MatrixGrid` row/footer cells |

## B.3 Style helper table

| Helper | Emits | Used by |
|---|---|---|
| `ui.styleBy(styleSet, options?)` | `{ styleSet, field?, map?, fallbackStyleKey? }` | `MatrixGrid.colorBy` |

This helper is optional because authors can hand-write the JSON object, but it is
small and makes scripts easier to teach.

## B.4 Preset recipe table

| Goja module | Recipe | Emits | Based on TS preset |
|---|---|---|---|
| `schedule.dsl` | `recipes.availabilityMatrix({ poll, tallies?, editableResponseId?, styleSet? })` | one `MatrixGrid` node | `availabilityMatrix` |
| `schedule.dsl` | `recipes.pollResults({ poll, tallies, styleSet? })` | `Stack` of `Caption + SegmentedBar` | `pollResults` |
| `calendar.dsl` | `recipes.monthCalendar({ events, monthISO, styleSet? })` | one `MonthGrid` node | `monthCalendar` |
| `calendar.dsl` | `recipes.weekCalendar({ events, daysISO, styleSet?, hourStart?, hourEnd? })` | one `TimeGrid` node | `weekCalendar` |

Use one options object per recipe. Existing recipes are options-object based, and
it leaves room for default overrides without positional-argument churn.

---

# Part C — Desired JavaScript authoring experience

## C.1 Low-level engine authoring

This is what a power user or a recipe implementation should be able to write:

```js
const data = require("data.dsl");
const ui = require("ui.dsl");

const availabilityStyleSet = {
  id: "availability",
  styles: {
    yes: { fill: "var(--mac-green)", labelColor: "var(--mac-text-inv)" },
    ifneedbe: { fill: "var(--mac-amber)", labelColor: "var(--mac-text)" },
    no: { fill: "var(--mac-accent-2)", labelColor: "var(--mac-text-inv)" },
    unknown: { fill: "var(--mac-surface)", labelColor: "var(--mac-text-dim)" },
  },
  legend: [
    { id: "yes", label: "Yes", styleKey: "yes" },
    { id: "ifneedbe", label: "If need be", styleKey: "ifneedbe" },
    { id: "no", label: "No", styleKey: "no" },
  ],
};

const node = data.matrixGrid({
  rows: [
    { id: "alice", name: "Alice", cells: { s1: "yes", s2: "no" } },
    { id: "you", name: "You", cells: { s1: "unknown", s2: "ifneedbe" } },
  ],
  columns: [
    { id: "s1", header: data.text("Thu 14:00"), meta: { yes: 1, total: 2 } },
    { id: "s2", header: data.text("Fri 10:00"), meta: { yes: 0, total: 2 } },
  ],
  valueAt: { mapField: "cells" },
  cell: data.cell.cycle(["yes", "ifneedbe", "no", "unknown"], {
    glyphs: { yes: "✓", ifneedbe: "~", no: "✕", unknown: "·" },
  }),
  styleSet: availabilityStyleSet,
  rowHeader: data.cell.field("name"),
  getRowKey: { field: "id" },
  editableRowKey: "you",
  footer: { header: data.text("yes"), cell: data.cell.template("${yes}/${total}") },
  onCellAction: data.action.server("poll.toggleCell", {
    payload: {
      responseId: { kind: "path", path: "rowKey" },
      optionId: { kind: "path", path: "colId" },
      state: { kind: "path", path: "value" },
    },
  }),
});
```

The emitted node must be a plain object. No Go values, no functions, no hidden
handles should cross the Widget IR boundary.

## C.2 Time engine authoring

```js
const time = require("time.dsl");

const month = time.monthGrid({
  monthISO: "2026-07",
  todayISO: "2026-07-06",
  selectedDateISO: "2026-07-09",
  markers: {
    "2026-07-09": { count: 3, styleKey: "meeting" },
  },
  styleSet: eventStyleSet,
  onDaySelectAction: time.action.server("calendar.day", {
    payload: { dateISO: { kind: "path", path: "dateISO" } },
  }),
});

const week = time.timeGrid({
  days: ["2026-07-06", "2026-07-07"],
  blocks: [
    { id: "e1", dayISO: "2026-07-06", startISO: "2026-07-06T09:00", endISO: "2026-07-06T09:30", styleKey: "meeting", label: time.text("Standup") },
  ],
  styleSet: eventStyleSet,
  onBlockSelectAction: time.action.server("calendar.select", {
    payload: { eventId: { kind: "path", path: "blockId" } },
  }),
});
```

## C.3 Domain recipe authoring

Most app scripts should not hand-author a `MatrixGrid`. They should use recipes:

```js
const schedule = require("schedule.dsl");
const calendar = require("calendar.dsl");
const ui = require("ui.dsl");

const page = ui.page({
  id: "scheduling-demo",
  title: "Scheduling demo",
  sections: [
    schedule.recipes.availabilityMatrix({
      poll,
      tallies,
      editableResponseId: "you",
    }),
    schedule.recipes.pollResults({ poll, tallies }),
    calendar.recipes.monthCalendar({ events, monthISO: "2026-07" }),
    calendar.recipes.weekCalendar({ events, daysISO: ["2026-07-06", "2026-07-07"] }),
  ],
});
```

This is the right intern-facing mental model:

- **Engines** are precise and configurable.
- **Recipes** are domain vocabulary and defaults.
- **Pages** compose both.

---

# Part D — Implementation design

## D.1 New constants and helper maps

Add module constants near the existing module constants:

```go
const (
    UIModuleName            = "ui.dsl"
    DataModuleName          = "data.dsl"
    DataV2ModuleName        = "data.v2.dsl"
    TimeModuleName          = "time.dsl"
    ScheduleModuleName      = "schedule.dsl"
    CalendarModuleName      = "calendar.dsl"
    ContextWindowModuleName = "context_window.dsl"
    CourseModuleName        = "course.dsl"
    CmsModuleName           = "cms.dsl"
)
```

Add helper maps:

```go
var dataHelpers = map[string]string{
    "dataTable":  "DataTable",
    "matrixGrid": "MatrixGrid",
}

var uiHelpers = map[string]string{
    // existing helpers...
    "segmentedBar": "SegmentedBar",
}

var timeHelpers = map[string]string{
    "monthGrid": "MonthGrid",
    "timeGrid":  "TimeGrid",
}
```

Then add module specs:

```go
{
    name:    TimeModuleName,
    helpers: timeHelpers,
    action:  true,
    doc:     "time.dsl provides generic month-grid and week/day time-grid helpers.",
},
{
    name:    ScheduleModuleName,
    action:  true,
    recipes: []string{"availabilityMatrix", "pollResults"},
    doc:     "schedule.dsl provides scheduling presets built from generic Widget IR engines.",
},
{
    name:    CalendarModuleName,
    action:  true,
    recipes: []string{"monthCalendar", "weekCalendar"},
    doc:     "calendar.dsl provides calendar presets built from generic time-grid engines.",
},
```

`Register` and `init` already iterate `moduleSpecs`, so the new modules are
registered automatically once they are in that slice.

## D.2 Extend `cellObject`

Add two exports to `cellObject()`:

```go
setExport(cell, "cycle", func(states []string, options ...goja.Value) map[string]any {
    out := map[string]any{"kind": "cycle", "states": states}
    mergeOptions(out, exportOptions(options)) // glyphs, styleSet
    return out
})

setExport(cell, "value", func() map[string]any {
    return map[string]any{"kind": "value"}
})
```

`cycle` should accept an empty state list but tests should discourage it. The TS
adapter falls back safely, but a useful DSL error would be better if we want to
validate aggressively.

## D.3 Add `styleBy`

Because `StyleBySpec` is generic, put `styleBy` on `ui.dsl` rather than a domain
module:

```go
func (r *runtime) installStyleHelpers(exports *goja.Object) {
    setExport(exports, "styleBy", func(styleSet goja.Value, options ...goja.Value) map[string]any {
        out := map[string]any{"styleSet": styleSet.Export()}
        mergeOptions(out, exportOptions(options))
        return out
    })
}
```

Call it from `install` when `spec.name == UIModuleName`:

```go
if spec.name == UIModuleName {
    setExport(exports, "section", r.sectionVerb)
    r.installStyleHelpers(exports)
}
```

Use the helper like this:

```js
ui.styleBy(availabilityStyleSet, {
  field: "state",
  map: { unknown: "noData" },
  fallbackStyleKey: "unknown",
})
```

## D.4 Palette constants and helpers

The TS package has canonical palettes in `scheduling/palettes.ts`. The Go DSL
recipes need the same shapes. Do not import TypeScript; port the JSON shape to Go.

Recommended Go helpers:

```go
func availabilityStyleSet() map[string]any { ... }
func eventStyleSet() map[string]any { ... }
func availabilityStates() []string { return []string{"yes", "ifneedbe", "no", "unknown"} }
func availabilityGlyphs() map[string]any { return map[string]any{"yes": "✓", "ifneedbe": "~", "no": "✕", "unknown": "·"} }
```

These helpers should return fresh maps/slices so recipe calls cannot mutate a
shared global map in surprising ways.

Optionally expose palette helpers to JS:

```go
if spec.name == ScheduleModuleName {
    setExport(exports, "availabilityStyleSet", availabilityStyleSet)
    setExport(exports, "availabilityStates", availabilityStates)
    setExport(exports, "availabilityGlyphs", availabilityGlyphs)
}
if spec.name == CalendarModuleName {
    setExport(exports, "eventStyleSet", eventStyleSet)
}
```

This is not required for recipes, but it helps scripts that use low-level engine
helpers.

## D.5 Add recipe dispatch cases

Extend `recipesObject`:

```go
case "availabilityMatrix":
    setExport(recipes, name, r.availabilityMatrixRecipe)
case "pollResults":
    setExport(recipes, name, r.pollResultsRecipe)
case "monthCalendar":
    setExport(recipes, name, r.monthCalendarRecipe)
case "weekCalendar":
    setExport(recipes, name, r.weekCalendarRecipe)
```

## D.6 Recipe implementation pseudocode

### `availabilityMatrixRecipe`

Input shape:

```js
schedule.recipes.availabilityMatrix({
  poll,
  tallies,
  editableResponseId,
  styleSet,       // optional override
  actionName,     // optional, default "poll.toggleCell"
})
```

Pseudocode:

```go
func (r *runtime) availabilityMatrixRecipe(call goja.FunctionCall) goja.Value {
    options := firstObject(call.Arguments)
    poll, _ := options["poll"].(map[string]any)
    if poll == nil {
        panic(r.vm.NewGoError(fmt.Errorf("schedule.dsl recipes.availabilityMatrix requires { poll }")))
    }

    responses := anySlice(poll["responses"])
    pollOptions := anySlice(poll["options"])
    tallies := talliesByOptionID(anySlice(options["tallies"]))
    total := len(responses)

    columns := []any{}
    for _, raw := range pollOptions {
        option := raw.(map[string]any)
        id := stringFromMap(option, "id", "")
        slot, _ := option["slot"].(map[string]any)
        tally := tallies[id]
        columns = append(columns, map[string]any{
            "id": id,
            "header": textNode(formatSlot(slot) + bestSuffix(tally)),
            "meta": map[string]any{"yes": intFrom(tally["yes"]), "total": total},
        })
    }

    styleSet := valueOrDefault(options["styleSet"], availabilityStyleSet())
    props := map[string]any{
        "ariaLabel": poll["title"],
        "rows": responses,
        "columns": columns,
        "valueAt": map[string]any{"mapField": "cells"},
        "cell": map[string]any{"kind": "cycle", "states": availabilityStates(), "glyphs": availabilityGlyphs()},
        "styleSet": styleSet,
        "rowHeader": map[string]any{"kind": "field", "field": "name"},
        "getRowKey": map[string]any{"field": "id"},
        "footer": map[string]any{"header": textNode("yes"), "cell": map[string]any{"kind": "template", "template": "${yes}/${total}"}},
        "onCellAction": map[string]any{"kind": "server", "name": stringFromMap(options, "actionName", "poll.toggleCell"), "payload": map[string]any{
            "pollId": poll["id"],
            "responseId": pathPart("rowKey"),
            "optionId": pathPart("colId"),
            "state": pathPart("value"),
        }},
    }
    copyIfPresent(props, options, "editableResponseId") // but rename to editableRowKey
    if v := options["editableResponseId"]; v != nil { props["editableRowKey"] = v }

    return r.vm.ToValue(componentNode("MatrixGrid", props))
}
```

Important implementation details:

- Rename `editableResponseId` to `editableRowKey`; the browser component does not
  know scheduling vocabulary.
- Use payload path specs exactly as the TS adapter expects:
  `{ kind:"path", path:"rowKey" }`, `{ kind:"path", path:"colId" }`,
  `{ kind:"path", path:"value" }`.
- Do not compute tallies inside the recipe. Tallies are server/domain data. The
  recipe only places them into column `meta` so the footer can render them.

### `pollResultsRecipe`

Input:

```js
schedule.recipes.pollResults({ poll, tallies, styleSet })
```

Output shape:

```
Stack
├─ Stack
│  ├─ Caption("Jul 9 · 14:00 ★")
│  └─ SegmentedBar({ segments:[yes, maybe, no], styleSet, showCounts:true })
└─ ...
```

Pseudocode:

```go
func (r *runtime) pollResultsRecipe(call goja.FunctionCall) goja.Value {
    options := firstObject(call.Arguments)
    poll := requiredObject(options, "poll")
    tallies := talliesByOptionID(anySlice(options["tallies"]))
    styleSet := valueOrDefault(options["styleSet"], availabilityStyleSet())

    children := []any{}
    for _, raw := range anySlice(poll["options"]) {
        option := raw.(map[string]any)
        id := stringFromMap(option, "id", "")
        tally := tallies[id]
        segments := []any{
            map[string]any{"value": valueOrDefault(tally["yes"], 0), "styleKey": "yes", "label": textNode("yes")},
            map[string]any{"value": valueOrDefault(tally["ifneedbe"], 0), "styleKey": "ifneedbe", "label": textNode("maybe")},
            map[string]any{"value": valueOrDefault(tally["no"], 0), "styleKey": "no", "label": textNode("no")},
        }
        children = append(children,
            componentNode("Stack", map[string]any{"gap": "xs"},
                componentNode("Caption", map[string]any{}, textNode(formatSlot(slotOf(option))+bestSuffix(tally))),
                componentNode("SegmentedBar", map[string]any{"segments": segments, "styleSet": styleSet, "showCounts": true}),
            ),
        )
    }
    return r.vm.ToValue(componentNode("Stack", map[string]any{"gap": "md"}, children...))
}
```

### `monthCalendarRecipe`

Input:

```js
calendar.recipes.monthCalendar({ events, monthISO, styleSet })
```

Pseudocode:

```go
func (r *runtime) monthCalendarRecipe(call goja.FunctionCall) goja.Value {
    options := firstObject(call.Arguments)
    markers := map[string]any{}
    for _, raw := range anySlice(options["events"]) {
        event := raw.(map[string]any)
        date := first10(stringFromMap(event, "startISO", ""))
        existing := objectOrEmpty(markers[date])
        markers[date] = map[string]any{
            "count": intFrom(existing["count"]) + 1,
            "styleKey": valueOrDefault(existing["styleKey"], event["colorKey"]),
        }
    }
    props := map[string]any{
        "monthISO": stringFromMap(options, "monthISO", ""),
        "markers": markers,
        "styleSet": valueOrDefault(options["styleSet"], eventStyleSet()),
    }
    copyIfPresent(props, options, "selectedDateISO")
    copyIfPresent(props, options, "todayISO")
    return r.vm.ToValue(componentNode("MonthGrid", props))
}
```

### `weekCalendarRecipe`

Input:

```js
calendar.recipes.weekCalendar({ events, daysISO, styleSet, hourStart, hourEnd })
```

Pseudocode:

```go
func (r *runtime) weekCalendarRecipe(call goja.FunctionCall) goja.Value {
    options := firstObject(call.Arguments)
    blocks := []any{}
    for _, raw := range anySlice(options["events"]) {
        event := raw.(map[string]any)
        blocks = append(blocks, map[string]any{
            "id": event["id"],
            "dayISO": first10(stringFromMap(event, "startISO", "")),
            "startISO": event["startISO"],
            "endISO": event["endISO"],
            "styleKey": event["colorKey"],
            "label": textNode(stringFromMap(event, "title", "")),
        })
    }
    props := map[string]any{
        "days": anySlice(options["daysISO"]),
        "blocks": blocks,
        "styleSet": valueOrDefault(options["styleSet"], eventStyleSet()),
        "hourStart": valueOrDefault(options["hourStart"], 8),
        "hourEnd": valueOrDefault(options["hourEnd"], 18),
    }
    return r.vm.ToValue(componentNode("TimeGrid", props))
}
```

Do not pass `allDay` until the TS review finding is resolved. The current
`TimeGrid` TS component accepts `allDay` but drops all-day events. The Go recipe
should avoid reinforcing that contract until the frontend is fixed.

---

# Part E — TypeScript declaration design

`TypeScriptModule` emits `.d.ts` declarations from the same module specs used by
the runtime. That is good: helper map entries automatically appear as functions.
However, several new surfaces need explicit declaration edits.

## E.1 Automatic declarations

Adding these helper map entries automatically creates declarations like:

```ts
export function matrixGrid(props?: Props | WidgetChild, ...children: WidgetChild[]): WidgetNode;
export function segmentedBar(props?: Props | WidgetChild, ...children: WidgetChild[]): WidgetNode;
export function monthGrid(props?: Props | WidgetChild, ...children: WidgetChild[]): WidgetNode;
export function timeGrid(props?: Props | WidgetChild, ...children: WidgetChild[]): WidgetNode;
```

That is enough for the classic open-props style.

## E.2 Manual declarations to add

When `moduleSpec.cell` is true, `TypeScriptModule` appends the `cell` object.
Add declarations for the new cell helpers:

```ts
cycle(states: string[], options?: Props): CellSpec;
value(): CellSpec;
```

When `moduleSpec.name == UIModuleName`, add:

```ts
export function styleBy(styleSet: Props, options?: Props): Props;
```

For `time.dsl`, the generic helper map declarations are fine. For
`schedule.dsl` and `calendar.dsl`, recipes are currently declared generically as
`name(options: Props): WidgetNode;`, which is acceptable for phase 1. Later, add
specialized interfaces if authors need stronger help.

## E.3 TypeScript fixture tests

Add a new fixture test, similar to `TestDataV2TypeScriptFixture...`, that compiles
this script:

```ts
/// <reference path="./widgetdsl.d.ts" />
import * as ui from "ui.dsl";
import * as data from "data.dsl";
import * as time from "time.dsl";
import * as schedule from "schedule.dsl";
import * as calendar from "calendar.dsl";

const styleSet = { styles: { yes: { fill: "green" } }, legend: [] };

const grid = data.matrixGrid({
  rows: [{ id: "you", cells: { s1: "yes" } }],
  columns: [{ id: "s1", header: data.text("Thu") }],
  valueAt: { mapField: "cells" },
  cell: data.cell.cycle(["yes", "no"], { glyphs: { yes: "✓", no: "✕" } }),
  colorBy: ui.styleBy(styleSet),
});

const month = time.monthGrid({ monthISO: "2026-07" });
const week = time.timeGrid({ days: ["2026-07-06"], blocks: [], styleSet });
const pollNode = schedule.recipes.availabilityMatrix({ poll: { id: "p", title: "Poll", options: [], responses: [] } });
const calNode = calendar.recipes.monthCalendar({ events: [], monthISO: "2026-07" });

grid.kind; month.kind; week.kind; pollNode.kind; calNode.kind;
```

This fixture catches missing declarations before a user does.

---

# Part F — Runtime tests to add

## F.1 Module export test

Extend `TestSplitModulesExportExpectedHelpersAndOmitCrossDomainHelpers` or add a
new test:

```go
value, err := vm.RunString(`
  const ui = require("ui.dsl");
  const data = require("data.dsl");
  const time = require("time.dsl");
  const schedule = require("schedule.dsl");
  const calendar = require("calendar.dsl");
  ({
    matrixGrid: typeof data.matrixGrid,
    segmentedBar: typeof ui.segmentedBar,
    styleBy: typeof ui.styleBy,
    cellCycle: typeof data.cell.cycle,
    cellValue: typeof data.cell.value,
    monthGrid: typeof time.monthGrid,
    timeGrid: typeof time.timeGrid,
    scheduleRecipes: typeof schedule.recipes.availabilityMatrix,
    calendarRecipes: typeof calendar.recipes.weekCalendar,
  })
`)
```

Assert each is `"function"` except recipe containers as appropriate.

## F.2 Low-level helper shape test

Test that helpers emit the exact Widget IR type strings:

```js
const data = require("data.dsl");
const time = require("time.dsl");
JSON.stringify({
  grid: data.matrixGrid({ rows: [], columns: [] }),
  month: time.monthGrid({ monthISO: "2026-07" }),
  cycle: data.cell.cycle(["yes", "no"], { glyphs: { yes: "✓" } }),
  value: data.cell.value(),
});
```

Assertions:

- `grid.kind == "component"`, `grid.type == "MatrixGrid"`.
- `month.type == "MonthGrid"`.
- `cycle.kind == "cycle"`, `cycle.states` preserved.
- `value.kind == "value"`.

## F.3 Recipe shape tests

Add one test per recipe. Keep the assertions structural and stable; do not assert
on every CSS token.

### Availability matrix assertions

- result `type == "MatrixGrid"`;
- `props.valueAt.mapField == "cells"`;
- `props.cell.kind == "cycle"`;
- `props.onCellAction.kind == "server"`;
- `props.onCellAction.payload.responseId.kind == "path"`;
- footer cell template is `${yes}/${total}`;
- `columns[0].meta.yes` is populated from the tally.

### Poll results assertions

- result `type == "Stack"`;
- it has one child per poll option;
- each child contains a `Caption` and a `SegmentedBar`;
- each `SegmentedBar` has three segments: `yes`, `ifneedbe`, `no`.

### Month calendar assertions

- result `type == "MonthGrid"`;
- `props.monthISO` is preserved;
- events on the same date increment `markers[date].count`;
- first event color key wins unless an explicit policy is chosen.

### Week calendar assertions

- result `type == "TimeGrid"`;
- `props.days` comes from `daysISO`;
- each event becomes one block with `dayISO = startISO.slice(0,10)`;
- `label` is a text node, not a raw string, so it is a valid `RenderableValue`.

## F.4 TypeScript declaration tests

Add fragments to `typescript_test.go` so declarations prove the new public
surface exists:

```go
wantFragments := []string{
    "export function matrixGrid",
    "cycle(states: string[]",
    "value(): CellSpec;",
    "export function segmentedBar",
    "export function monthGrid",
    "export function timeGrid",
    "availabilityMatrix(options: Props): WidgetNode;",
    "weekCalendar(options: Props): WidgetNode;",
}
```

Prefer a separate test for the new modules instead of overloading the data-v2
fixture test, because data-v2 deliberately omits legacy option-bag helpers.

---

# Part G — Implementation phases

## Phase 1 — Engine helper parity

Files:

- `pkg/widgetdsl/module.go`
- `pkg/widgetdsl/typescript.go`
- `pkg/widgetdsl/module_test.go`

Work:

1. Add `matrixGrid` to `dataHelpers`.
2. Add `segmentedBar` to `uiHelpers`.
3. Add `TimeModuleName`, `timeHelpers`, and a `time.dsl` module spec.
4. Add `cell.cycle` and `cell.value`.
5. Add `ui.styleBy`.
6. Update TypeScript declarations.
7. Add runtime export and low-level shape tests.

Validation:

```bash
go test ./pkg/widgetdsl/... -count=1
```

At this point, hand-authored low-level Goja scripts should be able to produce all
four engine nodes.

## Phase 2 — Palette helpers and recipe modules

Files:

- `pkg/widgetdsl/module.go`
- optionally a new `pkg/widgetdsl/scheduling_recipes.go` if `module.go` gets too large
- `pkg/widgetdsl/module_test.go`

Work:

1. Add `ScheduleModuleName` and `CalendarModuleName`.
2. Add module specs with recipe names.
3. Add palette helper functions returning fresh maps.
4. Add `recipesObject` cases.
5. Implement `availabilityMatrixRecipe`, `pollResultsRecipe`, `monthCalendarRecipe`, and `weekCalendarRecipe`.
6. Add recipe shape tests.

Validation:

```bash
go test ./pkg/widgetdsl/... -count=1
```

At this point, domain scripts should be able to use scheduling/calendar presets.

## Phase 3 — TypeScript fixture coverage

Files:

- `pkg/widgetdsl/typescript.go`
- `pkg/widgetdsl/typescript_fixture_test.go`

Work:

1. Add declarations for the new cell/style helpers.
2. Ensure recipe names appear in generated declarations.
3. Add a compile fixture that imports `ui.dsl`, `data.dsl`, `time.dsl`,
   `schedule.dsl`, and `calendar.dsl`.

Validation:

```bash
go test ./pkg/widgetdsl/... -count=1
```

The fixture uses `packages/rag-evaluation-site/node_modules/.bin/tsc`; if the
compiler is missing, the existing test skips. In normal development, make sure
`pnpm install` has been run so this test actually executes.

## Phase 4 — Browser integration smoke

Files:

- possible fixture script under `ttmp/.../scripts` or Go testdata
- no required production file changes

Work:

1. Run a Goja script that emits a page containing all four recipes.
2. Serialize to JSON.
3. If a browser/dev harness exists, feed it to `WidgetRenderer` and confirm no
   `UnknownWidget` appears.
4. If no browser harness exists, at least assert all node `type` strings are in
   the TS registry list.

Validation:

```bash
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build-storybook
go test ./pkg/widgetdsl/... -count=1
```

---

# Part H — Decision records

## Decision: Add `time.dsl` instead of placing time widgets in `ui.dsl`

- **Context:** The TS manifests for `MonthGrid` and `TimeGrid` already use
  `time.dsl`, but Go does not yet have that module.
- **Options considered:** Move the TS widgets to `ui.dsl`; add them to
  `data.dsl`; add `time.dsl` in Go.
- **Decision:** Add `time.dsl` in Go.
- **Rationale:** Calendar/time geometry is generic but conceptually distinct from
  general layout/atoms and data tables. Matching the TS manifests avoids a
  cross-language inconsistency.
- **Consequences:** One new module must be registered and declared. Scripts must
  import `require("time.dsl")` for month/week engines.
- **Status:** accepted.

## Decision: Keep engine helpers low-level and open-props

- **Context:** `componentFactory` already makes helper additions trivial, and the
  TS Widget IR props are open JSON by design.
- **Options considered:** Build fully typed Go builders for every scheduling
  engine now; use classic open-props helpers; skip low-level helpers and only add
  recipes.
- **Decision:** Use classic open-props helpers for engines.
- **Rationale:** This matches current `ui.dsl`/`data.dsl` style and gets parity
  quickly. Recipes provide ergonomic domain-level APIs.
- **Consequences:** TypeScript declarations are intentionally broad. Stronger
  scheduling-specific typing can be added later as a v2-style layer.
- **Status:** accepted for phase 1.

## Decision: Implement domain presets as `recipes`

- **Context:** Scheduling presets are composite: they produce configured subtrees,
  not single component types.
- **Options considered:** Add direct module functions; add `recipes.*`; add new
  browser adapters for organism-level nodes.
- **Decision:** Use `schedule.recipes.*` and `calendar.recipes.*` first.
- **Rationale:** Existing DSL modules already use `recipes` for composed trees,
  and no new browser adapter is needed.
- **Consequences:** The API is a little more verbose, but consistent. Direct
  aliases can be added later.
- **Status:** accepted.

## Decision: Do not put scheduling into `data.v2.dsl` yet

- **Context:** `data.v2.dsl` is a typed/fluent experiment with its own handle
  model and validation/lowering pipeline.
- **Options considered:** Extend v2 with scheduling builders now; keep scheduling
  in classic modules; create `schedule.v2.dsl`.
- **Decision:** Keep scheduling in classic modules for this ticket.
- **Rationale:** The browser contracts already exist as open Widget IR props, and
  scheduling engines are not data-table collection builders. Mixing them into
  data-v2 would conflate two separate experiments.
- **Consequences:** Fewer typed compile-time guarantees for scheduling scripts in
  phase 1; much less implementation risk.
- **Status:** accepted for phase 1; revisit later.

## Decision: Recipes should accept one options object

- **Context:** TS presets sometimes take `(poll, options)`, but Go recipe helpers
  already use an options-object convention.
- **Options considered:** Positional args matching TS; one options object; support
  both.
- **Decision:** Use one options object.
- **Rationale:** It is consistent with existing Go recipes and easier to extend
  without breaking call sites.
- **Consequences:** The Go API is not a byte-for-byte mirror of the TS helper
  signature, but it emits the same IR.
- **Status:** accepted.

---

# Part I — Common mistakes and how to avoid them

## Mistake 1: Adding a helper but forgetting declarations

A helper map entry makes the runtime function exist, but users also need `.d.ts`
declarations. Add tests that inspect both runtime exports and `TypeScriptModule`.

## Mistake 2: Returning functions or Go-only values in Widget IR

The browser receives JSON-like data. Do not return Go structs, Go function values,
or Goja handles from recipes. Always lower to `map[string]any`, `[]any`, strings,
numbers, booleans, and nil.

## Mistake 3: Mutating shared palette maps

If `availabilityStyleSet()` returns a global map and one recipe mutates it, later
recipes can change. Return fresh maps/slices.

## Mistake 4: Confusing `module` with `type`

`module` is for authoring and manifests. `type` is the browser registry key. The
Go helper must emit the exact browser `type` string (`"MonthGrid"`, not
`"time.MonthGrid"`) unless the TS adapter also uses a namespaced type.

## Mistake 5: Computing product state in the recipe

Recipes may transform DTOs into IR props, but they should not own server/domain
truth. For example, do not compute final poll tallies from responses in the DSL
unless the product explicitly accepts client-side derived tallies. Prefer
server-provided `tallies`.

## Mistake 6: Reinforcing unresolved frontend contracts

The frontend review found that `TimeGrid` currently accepts but drops `allDay`
blocks. The Go recipe should avoid emitting `allDay` until that behavior is fixed
or explicitly documented.

---

# Part J — File-by-file implementation checklist

## `pkg/widgetdsl/module.go`

- [ ] Add module constants: `TimeModuleName`, `ScheduleModuleName`,
      `CalendarModuleName`.
- [ ] Add `matrixGrid` to `dataHelpers`.
- [ ] Add `segmentedBar` to `uiHelpers`.
- [ ] Add `timeHelpers` with `monthGrid` and `timeGrid`.
- [ ] Add module specs for `time.dsl`, `schedule.dsl`, and `calendar.dsl`.
- [ ] Add `cell.cycle` and `cell.value`.
- [ ] Add `ui.styleBy` helper.
- [ ] Add palette helper functions.
- [ ] Add recipe dispatch cases.
- [ ] Add recipe implementations.

## `pkg/widgetdsl/typescript.go`

- [ ] Add `cell.cycle` and `cell.value` declaration lines.
- [ ] Add `styleBy` declaration for `ui.dsl`.
- [ ] Ensure new module specs produce declarations.
- [ ] Optionally add more precise recipe declarations later.

## `pkg/widgetdsl/module_test.go`

- [ ] Assert new modules can be required.
- [ ] Assert helper exports exist.
- [ ] Assert low-level helpers emit correct Widget IR shape.
- [ ] Assert recipes emit structural node shapes.

## `pkg/widgetdsl/typescript_test.go`

- [ ] Assert new declaration fragments are present.
- [ ] Keep data-v2 tests isolated; do not add classic helpers to data-v2.

## `pkg/widgetdsl/typescript_fixture_test.go`

- [ ] Add a fixture that imports all new modules and compiles representative
      helper/recipe calls.

## Optional new file: `pkg/widgetdsl/scheduling_recipes.go`

If `module.go` becomes hard to navigate, move only scheduling/calendar recipe
helpers and palette builders to a new file in the same package. Keep helper maps
and module specs in `module.go` so the module registry remains easy to audit.

---

# Part K — Final validation runbook

Run these before handing the work back:

```bash
# Go runtime + declaration tests
go test ./pkg/widgetdsl/... -count=1

# Browser package still understands the emitted types
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build
pnpm --dir packages/rag-evaluation-site build-storybook

# Optional package smoke if public exports changed
pnpm --dir packages/rag-evaluation-site pack:smoke
pnpm --dir packages/rag-evaluation-site consumer:smoke
```

If a test fails, read the failure in terms of the data flow:

1. Did the module register?
2. Did the JS helper exist?
3. Did the helper emit the expected Widget IR type/props?
4. Did TypeScript declarations describe the helper?
5. Does the browser registry have an adapter for the emitted `type`?

---

# References

- `pkg/widgetdsl/module.go` — module names, helper maps, registration, generic
  constructors, cell/action helpers, and recipes.
- `pkg/widgetdsl/typescript.go` — generated `.d.ts` declarations for DSL modules.
- `pkg/widgetdsl/module_test.go` — runtime export and JSON-serializability tests.
- `pkg/widgetdsl/grammar.go` — existing intent-level `data.dsl` grammar.
- `pkg/widgetdsl/v2_builders.go` and `pkg/widgetdsl/v2/spec/*` — future typed
  builder pattern reference.
- `packages/rag-evaluation-site/src/widgets/ir/engines.ts` — browser-side IR
  contracts for `MatrixGrid`, `SegmentedBar`, `MonthGrid`, and `TimeGrid`.
- `packages/rag-evaluation-site/src/components/molecules/*/*.widget.tsx` —
  browser adapters that interpret the props emitted by Goja.
- `packages/rag-evaluation-site/src/widgets/presets/scheduling.ts` — TS preset
  shapes to mirror in Go recipes.
- `reference/02-scheduling-and-calendar-widgets-implementation-and-ir-handoff-for-dsl-wiring.md`
  — concise handoff for the same work.
- `reference/03-code-review-and-design-assessment-for-scheduling-widgets.md` —
  code review and cleanup findings that should inform the DSL implementation.
