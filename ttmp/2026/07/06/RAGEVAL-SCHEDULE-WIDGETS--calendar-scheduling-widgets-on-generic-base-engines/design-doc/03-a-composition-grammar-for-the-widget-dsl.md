---
Title: A Composition Grammar for the Widget DSL
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
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/cmd/researchctl/doc/codesign-js-user-guide.md
      Note: Documentation explaining loops, fragments, callbacks, provenance, and trusted runtime boundaries
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/examples/jsverbs/codesign.js
      Note: Examples of builder fragments, JavaScript callbacks, sweep loops, and runtime callback registration
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/examples/jsverbs/research.js
      Note: Examples of scoped entity-builder lambdas and reusable fragments in the researchctl JavaScript API
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/builders.go
      Note: Go implementation of builder callbacks and .use(fragment) patterns
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/callbacks.go
      Note: Runtime callback registry pattern and callback ID boundary used as contrast for Widget IR
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/researchctl/builders.go
      Note: Go implementation of scoped builder callbacks for project graph entities
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/engines.ts
      Note: Generic engine contracts that domain views lower into
    - Path: repo://packages/rag-evaluation-site/src/widgets/presets/scheduling.ts
      Note: Scheduling preset layer that motivates product-level DSL views
    - Path: repo://pkg/widgetdsl/module.go
      Note: Current Goja module/helper/action/recipe runtime that the redesigned API would extend
    - Path: repo://pkg/widgetdsl/typescript.go
      Note: Generated declaration surface that must describe slots, bindings, domain views, and intent helpers
    - Path: repo://ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/02-goja-dsl-layer-design-and-implementation-guide-for-scheduling-widgets.md
      Note: Prior implementation-oriented DSL guide that this API redesign builds on
    - Path: repo://ttmp/2026/07/06/RAGEVAL-WIDGET-DECOMPOSITION--widget-library-decomposition-base-engines-contracts-and-dsl-ergonomics/design-doc/02-redesigning-the-widget-dsl-a-composition-first-opinionated-javascript-api.md
      Note: Colleague composition-first DSL proposal used as a contrast point
ExternalSources: []
Summary: 'A textbook-style redesign proposal for the Goja Widget DSL API: composition through named slots, opinionated presets, data bindings, and intent-level actions, applied backward to the existing modules.'
LastUpdated: 2026-07-07T15:35:00-04:00
WhatFor: Use this when rethinking the JavaScript API exposed by pkg/widgetdsl modules, especially after the scheduling/calendar engines have shown the value of generic engines plus domain presets.
WhenToUse: Read before designing new DSL helpers, recipe APIs, action APIs, or backwards-compatible migration paths for ui.dsl, data.dsl, time.dsl, schedule.dsl, calendar.dsl, and existing domain modules.
---



# A Composition Grammar for the Widget DSL

> **What this document is.** This is a second design document for the Goja Widget
> DSL API. The previous DSL guide explains how to wire the new scheduling/calendar
> widgets into the existing module system. This document asks a different question:
> if we could improve the JavaScript API itself, using what the scheduling work has
> taught us, what should the authoring language become?
>
> **The answer in one sentence:** the DSL should become a small composition grammar
> built from pages, sections, domain views, named slots, data bindings, and
> intent-level actions, with raw component construction kept as an explicit escape
> hatch rather than the everyday style.

---

## 1. The purpose of the DSL

The Widget DSL is not a renderer. It does not mount React components, compute CSS,
or reconcile browser state. It runs inside Goja, receives ordinary JavaScript values,
and returns Widget IR: plain JSON-like data that the browser can interpret later.
That one fact should shape the API. A good DSL should make it easy to write the
screen the product wants while still making the generated IR unsurprising.

The scheduling/calendar work clarifies the shape we want. The frontend now has
small generic engines such as `MatrixGrid`, `SegmentedBar`, `MonthGrid`, and
`TimeGrid`, and it also has domain presets such as availability polls and calendar
views. The important lesson is not merely that we need four more helpers. The
important lesson is that authors should not spend their time filling raw prop bags
for engines. They should describe product concepts, override a few named pieces, and
let the DSL lower that description to the generic engines.

That gives the DSL a precise job:

1. It should give script authors a vocabulary close to the product: poll, week,
   month, record, collection, summary, action.
2. It should expose composition seams deliberately: not every prop is an extension
   point, but every important visual unit has a named slot.
3. It should make data flow visible: values come from bindings, actions carry
   context, and defaults are part of the contract.
4. It should still compile to the same Widget IR the browser already renders.

The rest of this document proposes that API.

---

## 2. What changed after the scheduling work

Before the scheduling widgets, the DSL could get away with a flat helper model.
`ui.panel(...)` and `data.dataTable(...)` are useful because panels and tables are
already familiar shapes. The author passes props and children, and the browser does
what the component says.

Scheduling breaks that simple model. An availability poll is not just a grid. It is
a matrix with participants, time slots, a cycling cell state, a footer tally, an
editable row, a style palette, and an action payload that must carry row, column,
and value context. A week calendar is not just a time grid. It is a range of days,
a set of event blocks, a lane-packing interpretation, optional selection, and a
small vocabulary for labels and colors.

If the DSL only exposes low-level helpers, authors have to rebuild that knowledge
in every script:

```js
// This is possible, but it is not the API we should encourage.
data.matrixGrid({
  rows: poll.responses,
  columns: poll.options.map(option => ({
    id: option.id,
    header: timeLabel(option.slot),
    meta: tallyByOption[option.id],
  })),
  valueAt: { mapField: "cells" },
  cell: data.cell.cycle(["yes", "ifneedbe", "no", "unknown"], { glyphs }),
  styleSet: availabilityStyleSet,
  editableRowKey: currentResponseId,
  onCellAction: data.action.server("poll.toggle", {
    payload: {
      responseId: { kind: "path", path: "rowKey" },
      optionId: { kind: "path", path: "colId" },
      value: { kind: "path", path: "value" },
    },
  }),
});
```

This is a good internal lowering target. It is a poor everyday API. The script
contains too many implementation details: the matrix engine, the map-field accessor,
the state order, the glyph set, the payload path names, and the rename from
`responseId` to `rowKey`. None of those are the product concept. They are the cost
of lowering the product concept into Widget IR.

The DSL should let the same screen be written at the level of intent:

```js
schedule.availabilityPoll(poll, {
  currentResponse: "you",
  onChange: schedule.intent.toggleAvailability("poll.toggle"),
});
```

Then, when the author needs customization, the API should expose named seams:

```js
schedule.availabilityPoll(poll, {
  currentResponse: "you",
  slots: {
    participant: ({ response }) => ui.person(response.name, { subtitle: response.role }),
    option: ({ option }) => time.slotLabel(option.slot, { style: "weekday-time" }),
    footer: ({ tally }) => ui.caption(`${tally.yes}/${tally.total} yes`),
  },
  onChange: schedule.intent.toggleAvailability("poll.toggle"),
});
```

This is the central proposal: **composition should happen through named slots on
opinionated views.**

---

## 3. The design in one page

The redesigned API has five layers. Each layer is small, and each layer has a
specific responsibility.

| Layer | Authoring concept | Example | Responsibility |
|---|---|---|---|
| 1 | Page structure | `ui.page`, `ui.section`, `ui.stack` | Arrange screens and reading order. |
| 2 | Domain views | `schedule.availabilityPoll`, `time.week`, `data.collection` | Choose the engine and defaults for a product concept. |
| 3 | Named slots | `slots.participant`, `slots.option`, `slots.event` | Let authors customize stable subparts without rewriting the engine. |
| 4 | Data bindings | `bind.field`, `bind.map`, `bind.template`, `bind.context` | Describe how values are read from data and action context. |
| 5 | Intent actions | `schedule.intent.toggleAvailability`, `act.server` | Describe what happens and which contextual values are sent. |

The raw Widget IR layer remains available, but it is no longer the center of the
language:

```js
ui.raw("MatrixGrid", props, children);  // explicit escape hatch
```

The everyday code should read like this:

```js
const ui = require("ui.dsl");
const schedule = require("schedule.dsl");
const time = require("time.dsl");

module.exports = ui.page({ title: "Team scheduling" }, [
  ui.section("Find a time", [
    schedule.availabilityPoll(poll, {
      currentResponse: currentUser.responseId,
      onChange: schedule.intent.toggleAvailability("poll.toggle"),
    }),
  ]),

  ui.section("Best options", [
    schedule.pollSummary(poll, tallies, { order: "best-first" }),
  ]),

  ui.section("Calendar", [
    time.week(events, {
      range: time.weekOf("2026-07-06"),
      onSelect: time.intent.selectEvent("calendar.select"),
    }),
  ]),
]);
```

The important property of this example is not that it is short. It is that the
screen is composed from product-level nouns, and every noun still lowers to ordinary
Widget IR. An intern reading the script can identify the screen sections and the
user actions without knowing the `MatrixGrid` prop contract.

---

## 4. The first principle: compose with named slots, not arbitrary prop bags

A prop bag is open-ended. It is useful for a low-level component helper, but it does
not teach the author where customization is safe. If a recipe accepts every prop the
engine accepts, then the recipe is not really a recipe; it is a slightly renamed raw
component.

A named slot is different. A slot says: this subpart is intentionally replaceable,
and this is the context you receive when replacing it. It makes composition explicit
and bounded.

For an availability poll, the slots might be:

| Slot | Context | Default |
|---|---|---|
| `participant` | `{ response, rowIndex }` | respondent name text |
| `option` | `{ option, columnIndex }` | formatted date/time label |
| `cell` | `{ response, option, value, editable }` | cycle cell with availability glyphs |
| `footer` | `{ option, tally, total }` | yes/total caption |
| `empty` | `{ poll }` | muted empty-state caption |

For a week calendar, the slots might be:

| Slot | Context | Default |
|---|---|---|
| `dayHeader` | `{ dayISO, index }` | weekday/date label |
| `event` | `{ event, block, overlaps }` | compact event card |
| `nowMarker` | `{ nowISO }` | themed horizontal marker |
| `emptyDay` | `{ dayISO }` | no-op or subtle empty label |

A slot function runs in Goja at authoring time. It returns a Widget IR node. The
recipe calls the slot while lowering product data into engine props. This is already
proven by the existing master-detail recipe, which calls a JavaScript `detail(row)`
function and embeds the returned node. The redesign generalizes that pattern, but it
does not make every internal prop a slot. The domain view chooses the few seams that
matter.

The API shape is simple:

```js
schedule.availabilityPoll(poll, {
  slots: {
    participant: ({ response }) => ui.person(response.name),
    option: ({ option }) => time.slotLabel(option.slot),
    footer: ({ tally }) => ui.caption(`${tally.yes}/${tally.total}`),
  },
});
```

A slot may also be a spec when a full function would be noise:

```js
schedule.availabilityPoll(poll, {
  slots: {
    participant: data.cell.field("name"),
    cell: data.cell.cycle(schedule.availability.states()),
  },
});
```

That dual form is useful because some customization is structural and some is just
a rendering spec. The slot contract accepts both, validates both, and lowers both to
ordinary IR.

---

## 5. The second principle: domain views own engine defaults

An engine is a reusable arrangement primitive. A domain view is a product opinion
about how that engine should be configured. The DSL should expose both, but most
scripts should start with the domain view.

For scheduling, this means:

| Domain view | Engine underneath | Defaults owned by the view |
|---|---|---|
| `schedule.availabilityPoll` | `MatrixGrid` | availability states, glyphs, palette, row/column context, footer shape, toggle payload |
| `schedule.pollSummary` | `Stack` + `SegmentedBar` | yes/maybe/no segments, best-first ordering, count labels |
| `schedule.bookingPicker` | `MonthGrid` + `TimeGrid` or `MatrixGrid` | selectable slot semantics, disabled slots, booking action payload |
| `time.month` | `MonthGrid` | month bounds, marker aggregation, today/selected conventions |
| `time.week` | `TimeGrid` | day list, hour bounds, event block conversion, event selection payload |

The lower-level engines are still available:

```js
data.matrix({ rows, columns, cell, valueAt });
time.monthGrid({ monthISO, markers });
time.timeGrid({ days, blocks });
```

But a domain view is the normal entry point:

```js
time.week(events, {
  range: time.weekOf("2026-07-06"),
  slots: {
    event: ({ event }) => ui.card({ density: "compact" }, [
      ui.strong(event.title),
      ui.caption(time.rangeLabel(event.startISO, event.endISO)),
    ]),
  },
});
```

This distinction keeps the DSL opinionated. The author can customize labels and
cards, but they do not reimplement lane packing, day columns, event-to-block
conversion, or action context names in every script.

---

## 6. The third principle: bindings should replace stringly prop conventions

The current system already uses defunctionalized specs for cells, actions, and
styles. The redesign should make that idea systematic. Anywhere a view needs to
read from a datum or context object, authors should use a binding helper rather than
inventing a prop shape.

The proposed binding namespace:

```js
const bind = require("bind.dsl");

bind.field("name")                 // { kind: "field", path: "name" }
bind.path("owner.name")            // { kind: "path", path: "owner.name" }
bind.map("cells")                  // { kind: "map", field: "cells" }
bind.template("${first} ${last}")  // { kind: "template", template: "${first} ${last}" }
bind.context("rowKey")             // { kind: "context", path: "rowKey" }
bind.const("yes")                  // { kind: "const", value: "yes" }
```

The exact module name can be debated. It could be `data.bind`, `ui.bind`, or a
shared object installed into every DSL module. What matters is that there is one
vocabulary.

The binding helpers make action payloads readable:

```js
schedule.intent.toggleAvailability("poll.toggle", {
  pollId: bind.const(poll.id),
  responseId: bind.context("rowKey"),
  optionId: bind.context("colId"),
  state: bind.context("value"),
});
```

They also make engine props less ad hoc:

```js
data.matrix(rows, {
  columns,
  valueAt: bind.map("cells"),
  rowKey: bind.field("id"),
  colorBy: style.by(availabilityStyleSet, { value: bind.context("value") }),
});
```

A binding is not a JavaScript function. It is a serializable description of how the
browser or adapter should read a value later. That distinction is important. Slot
functions run during DSL lowering; bindings survive into the IR and are interpreted
when an event fires or a cell renders.

---

## 7. The fourth principle: actions should start from intent

The existing `action.server`, `action.navigate`, `action.event`, and related helpers
are necessary. They describe transport: where an interaction goes. They do not
describe product intent. Scheduling scripts should not have to remember that a
matrix cell change carries `rowKey`, `colId`, and `value`. That is the action context
of the engine, not the product vocabulary of the poll.

The redesigned DSL should keep transport helpers but add intent wrappers in domain
modules:

```js
schedule.intent.toggleAvailability("poll.toggle")
schedule.intent.submitResponse("poll.submit")
time.intent.selectDay("calendar.day")
time.intent.selectEvent("calendar.event")
data.intent.selectRecord("record.select")
```

An intent wrapper returns an ordinary action spec. It only fills the payload mapping
for the known view context.

For example:

```js
schedule.intent.toggleAvailability("poll.toggle")
```

lowers to:

```js
act.server("poll.toggle", {
  payload: {
    pollId: bind.context("poll.id"),
    responseId: bind.context("response.id"),
    optionId: bind.context("option.id"),
    state: bind.context("value"),
  },
});
```

The context names are the domain view's contract. The engine adapter may still use
`rowKey` and `colId` internally, but the domain view translates them before the
action leaves the view. That is the point: the script author should not have to know
engine-local names unless they are using the raw engine helper.

The escape hatch remains:

```js
act.server("poll.toggle", {
  payload: {
    responseId: bind.context("rowKey"),
    optionId: bind.context("colId"),
    state: bind.context("value"),
  },
});
```

Use that form inside engine-level scripts. Use intent wrappers in product-level
scripts.

---

## 8. A complete example: availability poll

This section walks through one screen twice: first as product-level DSL, then as the
lowering shape. The goal is to make the design concrete enough to implement.

### 8.1 Product-level script

```js
const ui = require("ui.dsl");
const schedule = require("schedule.dsl");
const time = require("time.dsl");

module.exports = ui.page({ id: "poll", title: poll.title }, [
  ui.section("Choose the times that work", [
    schedule.availabilityPoll(poll, {
      currentResponse: currentUser.responseId,
      readOnly: poll.closed,
      slots: {
        option: ({ option }) => ui.stack({ gap: "xxs", align: "center" }, [
          ui.strong(time.format(option.slot.startISO, "EEE")),
          ui.caption(time.formatRange(option.slot.startISO, option.slot.endISO, "HH:mm")),
        ]),
        participant: ({ response }) => ui.inline({ gap: "xs" }, [
          ui.avatar(response.name),
          ui.text(response.name),
        ]),
        empty: () => ui.callout("No responses yet. Share the poll link to collect availability."),
      },
      onChange: schedule.intent.toggleAvailability("poll.toggle"),
      onSubmit: schedule.intent.submitResponse("poll.submit"),
    }),
  ]),

  ui.section("Current best options", [
    schedule.pollSummary(poll, tallies, {
      order: "best-first",
      maxOptions: 5,
      slots: {
        label: ({ option, tally }) => ui.inline({}, [
          time.slotLabel(option.slot),
          ui.badge(`${tally.yes} yes`),
        ]),
      },
    }),
  ]),
]);
```

This script uses three kinds of composition:

- Page composition: `page` contains `section`; `section` contains views.
- View customization: `availabilityPoll` exposes `option`, `participant`, and
  `empty` slots.
- Intent binding: `schedule.intent.toggleAvailability` provides the correct action
  payload for this view.

The script does not mention `MatrixGrid`, `valueAt.mapField`, `editableRowKey`, or
`rowKey`. Those are lowering details.

### 8.2 Lowered engine shape

The domain view lowers to something equivalent to:

```js
data.matrix({
  rows: poll.responses,
  columns: poll.options.map(option => ({
    id: option.id,
    header: slots.option({ option }),
    meta: tallyByOption[option.id],
  })),
  valueAt: bind.map("cells"),
  rowKey: bind.field("id"),
  rowHeader: slots.participant || data.cell.field("name"),
  cell: slots.cell || data.cell.cycle(schedule.availability.states(), {
    glyphs: schedule.availability.glyphs(),
    styleSet: schedule.availability.styleSet(),
  }),
  editableRow: currentResponse,
  readOnly,
  footer: slots.footer || (({ tally }) => ui.caption(`${tally.yes}/${tally.total}`)),
  onCell: lowerIntent(onChange),
});
```

The lowered shape is not necessarily public API. It is a useful implementation
model: a domain view gathers data, applies defaults, calls slots, maps intent
actions, and emits an engine node.

### 8.3 Why this is better than a raw engine helper

The raw helper is still useful for experiments. The domain view is better for real
scripts because it owns three kinds of knowledge:

- It owns scheduling defaults such as state order, glyphs, and palette.
- It owns context translation from engine names to product names.
- It owns the stable extension seams, so authors customize visible pieces without
  depending on every engine prop.

That is the balance we want: opinionated outside, composable at the seams.

---

## 9. A complete example: booking picker

An availability poll collects preferences from many participants. A booking picker
chooses one available time from a service schedule. The same engines can power both,
but the DSL should not force the author to describe them the same way.

```js
const booking = require("booking.dsl");
const ui = require("ui.dsl");
const time = require("time.dsl");

module.exports = ui.page({ title: "Book office hours" }, [
  booking.picker(availability, {
    durationMinutes: 30,
    timezone: "Europe/Berlin",
    selectedSlot: draft.slotId,
    slots: {
      day: ({ day, count }) => ui.stack({ align: "center" }, [
        ui.strong(time.format(day, "d")),
        ui.caption(`${count} slots`),
      ]),
      slot: ({ slot, selected, disabled }) => ui.button({
        variant: selected ? "primary" : "secondary",
        disabled,
      }, time.formatRange(slot.startISO, slot.endISO, "HH:mm")),
    },
    onSelect: booking.intent.selectSlot("booking.selectSlot"),
    onConfirm: booking.intent.confirm("booking.confirm"),
  }),
]);
```

This example is intentionally not named `timeGrid`. The product concept is booking.
The view may lower to a `MonthGrid` plus a `TimeGrid`, or to a `MatrixGrid` if the
product later chooses a compact layout. The script should not change merely because
we choose a different engine internally.

This illustrates a rule for domain modules:

> Name the API after the user's task when the component represents a task. Name it
> after the engine only when the author is deliberately choosing an arrangement
> primitive.

---

## 10. A complete example: calendar week

Calendar week views need different slots from scheduling polls. Event rendering is
usually the main customization point; the day grid and lane packing should remain an
engine concern.

```js
const ui = require("ui.dsl");
const time = require("time.dsl");

ui.page({ title: "Calendar" }, [
  time.week(events, {
    range: time.range.week("2026-07-06"),
    hours: [8, 18],
    style: time.eventStyles({ palette: "default" }),
    slots: {
      dayHeader: ({ dayISO }) => ui.stack({ gap: "xxs", align: "center" }, [
        ui.caption(time.format(dayISO, "EEE")),
        ui.strong(time.format(dayISO, "d")),
      ]),
      event: ({ event }) => ui.card({ tone: event.colorKey }, [
        ui.strong(event.title),
        ui.caption(time.formatRange(event.startISO, event.endISO, "HH:mm")),
      ]),
    },
    onSelect: time.intent.selectEvent("calendar.selectEvent"),
  }),
]);
```

The lower-level engine helper remains:

```js
time.timeGrid({ days, blocks, hourStart: 8, hourEnd: 18 });
```

But most scripts should use `time.week`. The domain view turns `events` into blocks,
sets default hour bounds, filters all-day events according to the decided contract,
and maps event selection into an action payload.

The current frontend review found that `TimeGrid.allDay` exists in types but is not
rendered. The redesigned API should not expose all-day customization until the view
can honor it. When all-day rendering is implemented, it should appear as an explicit
slot:

```js
time.week(events, {
  slots: {
    allDayEvent: ({ event }) => ui.badge(event.title),
  },
});
```

That gives the contract a visible authoring surface only when the behavior exists.

---

## 11. Applying the style backward to existing modules

The phrase "apply it backward" matters. We are not designing a new DSL for only the
scheduling ticket. We are using the scheduling work to clarify the authoring style
that older modules should converge toward.

### 11.1 `ui.dsl`

`ui.dsl` should own page composition and generic visual structure:

```js
ui.page({ title }, [ ...sections ])
ui.section("Title", [ ...children ])
ui.stack({ gap }, [ ...children ])
ui.inline({ gap, align }, [ ...children ])
ui.card({ tone }, [ ...children ])
ui.when(condition, node)
ui.list(items, itemSlot, options)
ui.raw(type, props, children)
```

Existing component helpers can stay. The change is that the recommended style uses
array children and named layout helpers instead of deeply nested variadic calls.
Both can coexist; the runtime can flatten arrays so old and new styles interoperate.

Before:

```js
ui.panel({ title: "Summary" },
  ui.text("Revenue"),
  ui.caption("Last 30 days"));
```

After:

```js
ui.section("Summary", [
  ui.metric("Revenue", revenue, { caption: "Last 30 days" }),
]);
```

The new style is not just shorter. It gives the author a semantic unit (`section`,
`metric`) rather than a generic panel with arbitrary children.

### 11.2 `data.dsl`

`data.dsl` should move from table-first to collection-first. A collection view owns
selection, empty state, bulk actions, row rendering, and detail rendering as slots.

```js
data.collection(records, {
  key: bind.field("id"),
  schema: [
    data.field.text("name", { label: "Name" }),
    data.field.status("stage", { label: "Stage", palette: dealStagePalette }),
    data.field.currency("amount", { label: "Amount" }),
  ],
  view: "table",
  slots: {
    rowActions: ({ record }) => ui.inline({}, [
      ui.button("Open", { onClick: data.intent.openRecord("record.open") }),
      ui.button("Archive", { onClick: data.intent.archiveRecord("record.archive") }),
    ]),
    detail: ({ record }) => data.record(record, { mode: "read" }),
  },
});
```

The existing `data.dataTable` helper remains as an engine-level escape hatch. The
collection API becomes the normal path because it expresses what most product
screens are doing: showing records, fields, selection, and actions.

### 11.3 `data.v2.dsl`

The v2 fluent builder introduced a valuable idea: validate an intermediate model
before lowering it to IR. The redesign should keep that idea but not require authors
to learn a separate fluent dialect.

The new collection API can lower into the same internal model as `data.v2.dsl`:

```
data.collection(...)  ─┐
                       ├─> CollectionSpec -> validate -> lower -> Widget IR
data.v2.collection(...) ┘
```

During migration, `data.v2.dsl` remains available. Eventually it can become either
an internal implementation detail or a compatibility facade over the same spec.

### 11.4 `course.dsl`, `cms.dsl`, and other domain modules

Domain modules should expose task-level views with named slots:

```js
course.lessonOutline(course, {
  slots: {
    lesson: ({ lesson }) => ui.card({}, [ui.strong(lesson.title), ui.caption(lesson.duration)]),
  },
  onSelect: course.intent.selectLesson("course.selectLesson"),
});

cms.mediaLibrary(assets, {
  view: "grid",
  slots: {
    asset: ({ asset }) => ui.card({}, [ui.image(asset.thumbnail), ui.text(asset.title)]),
  },
  onOpen: cms.intent.openAsset("cms.openAsset"),
});
```

This is the same pattern as scheduling: domain view, named slots, intent actions,
engine lowering.

---

## 12. Proposed module surface

This section is a compact reference for the redesigned API. It is deliberately
small. The goal is not to list every eventual helper; it is to show the shape.

### 12.1 `ui.dsl`

```ts
page(options, children): WidgetNode
section(titleOrOptions, children): WidgetNode
stack(options, children): WidgetNode
inline(options, children): WidgetNode
card(options, children): WidgetNode
text(value, options?): WidgetNode
caption(value, options?): WidgetNode
strong(value, options?): WidgetNode
badge(value, options?): WidgetNode
button(labelOrOptions, optionsOrChildren?, children?): WidgetNode
when(condition, node): WidgetNode | null
map(items, slot): WidgetNode[]
raw(type, props?, children?): WidgetNode
```

### 12.2 `bind.dsl` or shared `bind`

```ts
field(path): BindingSpec
path(path): BindingSpec
map(field): BindingSpec
template(template): BindingSpec
context(path): BindingSpec
const(value): BindingSpec
```

### 12.3 `act.dsl` or shared `act`

```ts
server(name, options?): ActionSpec
navigate(to, options?): ActionSpec
event(name, options?): ActionSpec
copy(value, options?): ActionSpec
payload(shape): PayloadSpec
```

The old `data.action.*` exports can alias this namespace.

### 12.4 `data.dsl`

```ts
collection(records, options): WidgetNode
record(record, options): WidgetNode
matrix(rows, options): WidgetNode
field.text(path, options?): FieldSpec
field.status(path, options?): FieldSpec
field.currency(path, options?): FieldSpec
cell.field(path, options?): CellSpec
cell.template(template, options?): CellSpec
cell.cycle(states, options?): CellSpec
cell.value(options?): CellSpec
intent.selectRecord(actionName, options?): ActionSpec
intent.updateField(actionName, options?): ActionSpec
```

### 12.5 `time.dsl`

```ts
month(eventsOrMarkers, options): WidgetNode
week(events, options): WidgetNode
day(events, options): WidgetNode
slotLabel(slot, options?): WidgetNode
format(iso, format): string
formatRange(startISO, endISO, format?): string
range.week(anchorISO, options?): RangeSpec
intent.selectDay(actionName, options?): ActionSpec
intent.selectEvent(actionName, options?): ActionSpec
```

### 12.6 `schedule.dsl`

```ts
availabilityPoll(poll, options): WidgetNode
pollSummary(poll, tallies, options?): WidgetNode
bookingPicker(availability, options): WidgetNode
availability.states(): string[]
availability.glyphs(): Record<string, string>
availability.styleSet(options?): ContextStyleSet
intent.toggleAvailability(actionName, options?): ActionSpec
intent.submitResponse(actionName, options?): ActionSpec
```

The exact TypeScript names can be refined, but the grammar should remain stable:
view functions, slots, bindings, intents, raw escape hatch.

---

## 13. Implementation model

The implementation does not require a rewrite. It requires a small runtime kernel
that recipes can share.

### 13.1 Normalize children

Existing helpers accept variadic children. The new examples use array children
because arrays read better for document-like composition. Support both:

```go
func normalizeChildren(args []goja.Value) []any {
    // flatten arrays, discard null/undefined/false, preserve valid nodes/strings
}
```

This one helper unlocks `ui.section("Title", [a, b, c])` without breaking
`ui.section("Title", a, b, c)`.

### 13.2 Normalize slots

A slot can be missing, a function, a spec, a literal node, or a string. The runtime
should normalize it once:

```go
type Slot struct {
    Kind string // function, spec, node, none
    Value goja.Value
}

func (r *runtime) callSlot(slot Slot, context map[string]any, fallback func(map[string]any) any) any {
    switch slot.Kind {
    case "function":
        return r.callJSFunction(slot.Value, context)
    case "spec":
        return renderSpec(slot.Value.Export(), context)
    case "node":
        return slot.Value.Export()
    case "none":
        return fallback(context)
    }
}
```

The first implementation can be simpler: support only function and fallback for
slots in domain views, and only later add spec/node slot forms. The public contract
should still name the slot concept early.

### 13.3 Normalize bindings

Binding helpers are plain map constructors. Evaluation belongs in adapters and
action payload resolution, not in the Go DSL runtime unless a recipe needs to read
source data while lowering.

```go
func binding(kind string, fields map[string]any) map[string]any {
    out := map[string]any{"kind": kind}
    for k, v := range fields { out[k] = v }
    return out
}
```

The important implementation detail is consistency: every binding should use one
shape, and action payloads should reuse it.

### 13.4 Domain views lower to engine helpers

`availabilityPoll` can be implemented as a Go function that builds a `MatrixGrid`
node. It does not need a new browser adapter.

```
schedule.availabilityPoll
  -> parse poll/options
  -> choose defaults
  -> call slots for headers/participants/footer if provided
  -> build MatrixGrid props
  -> return { kind:"component", type:"MatrixGrid", props, children:[] }
```

That is the same lowering strategy as the previous DSL guide. The difference is the
public API: scripts call `schedule.availabilityPoll`, not
`schedule.recipes.availabilityMatrix` and not `data.matrixGrid` unless they want the
engine.

### 13.5 Intent wrappers lower to action specs

Intent wrappers are also just map constructors:

```go
func toggleAvailabilityIntent(actionName string, options map[string]any) map[string]any {
    payload := map[string]any{
        "pollId": binding("context", map[string]any{"path": "poll.id"}),
        "responseId": binding("context", map[string]any{"path": "response.id"}),
        "optionId": binding("context", map[string]any{"path": "option.id"}),
        "state": binding("context", map[string]any{"path": "value"}),
    }
    mergePayloadOverrides(payload, options["payload"])
    return serverAction(actionName, map[string]any{"payload": payload})
}
```

The domain view is responsible for ensuring the context object passed to the action
resolver contains `poll`, `response`, `option`, and `value`.

---

## 14. Migration plan

### Phase 1: Add aliases and kernel helpers

- Add `ui.raw` as an alias of `component`.
- Add child-array flattening to generic helpers.
- Add `ui.when` and `ui.map`.
- Add `bind` helpers under `data.bind` first, or as a shared object on each module.
- Add `act` aliases for existing action helpers while keeping `data.action` working.

This phase is safe because it does not change existing output.

### Phase 2: Add scheduling/time domain views

- Add `time.month`, `time.week`, and `time.day` on top of `MonthGrid`/`TimeGrid`.
- Add `schedule.availabilityPoll`, `schedule.pollSummary`, and
  `schedule.bookingPicker` on top of `MatrixGrid`, `SegmentedBar`, and time views.
- Add named slots for the visible subparts only.
- Add intent wrappers for scheduling and calendar actions.

This phase proves the pattern on the newest widgets.

### Phase 3: Backfill data and existing domain modules

- Add `data.collection` as the recommended API over tables, records, selection,
  detail slots, and actions.
- Keep `data.dataTable` as the table engine helper.
- Make `data.v2.dsl` and classic `data.collection` lower through the same internal
  collection spec.
- Update `course.dsl` and `cms.dsl` recipes to use named slots and intent wrappers.

This phase applies the style backward.

### Phase 4: Generate declarations from descriptors

Each domain view should have a descriptor:

```go
type ViewSpec struct {
    Module  string
    Name    string
    Engine  string
    Slots   []SlotSpec
    Intents []IntentSpec
    Docs    string
}
```

The runtime can use the descriptor to install helpers; the TypeScript generator can
use it to emit declarations; documentation can use it to list slots and context
shapes. This is how we prevent the DSL from drifting again.

### Phase 5: Deprecate but do not remove old styles quickly

- Keep `component` as an alias of `raw` with documentation that says `raw` is the
  preferred name for new scripts.
- Keep `data.action` as an alias of `act`.
- Keep `schedule.recipes.availabilityMatrix` as an alias of
  `schedule.availabilityPoll` or as an engine-level preset, depending on final naming.
- Keep `data.v2.dsl` until callers have moved.

The migration should be boring. The API can become better without making existing
scripts fail.

---

## 15. Decision records

### Decision: Prefer named slots over unrestricted unit renderers

- **Context:** A colleague's composition-first proposal centers on engine verbs that
  accept unit renderers such as `card`, `cell`, and `item`.
- **Options considered:** Accept arbitrary unit renderers for every engine; expose
  only prop bags; expose domain views with named slots.
- **Decision:** Use named slots on domain views as the primary API, while allowing
  lower-level engine helpers to accept unit renderers where they are genuinely the
  engine contract.
- **Rationale:** Named slots preserve composition but keep domain APIs opinionated.
  They tell authors which subparts are stable extension points and prevent every
  engine prop from becoming public domain-view API.
- **Consequences:** Domain view descriptors need slot names and context shapes. The
  API is slightly less free-form, but easier to document and safer to evolve.
- **Status:** proposed.

### Decision: Name product APIs after tasks, not engines

- **Context:** `MatrixGrid` and `TimeGrid` are engine names. Product scripts are
  usually about availability polls, booking pickers, and calendars.
- **Options considered:** Expose only engine helpers; expose recipes under
  `recipes.*`; expose task-level domain views.
- **Decision:** Use task-level names such as `schedule.availabilityPoll` and
  `booking.picker` for everyday scripts; keep engine helpers for advanced use.
- **Rationale:** Product names survive engine changes. If a booking picker changes
  from a matrix layout to a month/week layout, the script should not need to change.
- **Consequences:** The DSL needs domain modules and intent wrappers, not just
  generic helpers.
- **Status:** proposed.

### Decision: Add one binding vocabulary

- **Context:** Current IR uses several small ad hoc shapes for field access,
  template access, map access, and action payload paths.
- **Options considered:** Keep each helper's local convention; add a shared binding
  vocabulary; use JavaScript functions for all access.
- **Decision:** Add shared `bind.*` helpers for serializable access specs.
- **Rationale:** Bindings make data flow visible and serializable. JavaScript
  functions are useful for author-time slot rendering, but action payloads and
  browser-time rendering need data descriptions, not closures.
- **Consequences:** Adapters and action resolution should converge on the binding
  shape over time.
- **Status:** proposed.

### Decision: Add intent wrappers but keep transport actions

- **Context:** `action.server` is flexible but requires authors to know engine
  context names.
- **Options considered:** Only use transport actions; only use domain intents; use
  domain intents that lower to transport actions.
- **Decision:** Use domain intents as the recommended layer and keep transport
  actions as the escape hatch.
- **Rationale:** Intent wrappers make product scripts clearer and reduce duplicated
  payload mapping. Transport helpers are still necessary for advanced and new cases.
- **Consequences:** Each domain view must define its action context contract.
- **Status:** proposed.

### Decision: Introduce `raw` without removing `component`

- **Context:** `component(type, props, children)` is necessary but currently feels
  like a normal path.
- **Options considered:** Keep only `component`; rename to `raw` and break scripts;
  add `raw` as alias and document `component` as legacy.
- **Decision:** Add `raw` as the preferred escape-hatch name and keep `component`.
- **Rationale:** The name `raw` communicates intent without a breaking change.
- **Consequences:** Documentation and examples should use `raw` only when escaping
  the curated DSL.
- **Status:** proposed.

---

## 16. Alternatives considered

### Alternative 1: Only add the missing low-level helpers

This is the smallest change: add `matrixGrid`, `timeGrid`, `cell.cycle`, and the
recipe names described in the previous guide. We should still do that as an
implementation phase, but it is not enough as an API vision. It gives authors the
power to express scheduling widgets, but it still asks them to understand every
engine prop and action context.

### Alternative 2: Make everything a fluent builder

A fluent API can read well for schemas:

```js
data.collection(records).field("name").field("stage").table().toIR();
```

The problem is that screens are trees, and trees compose more naturally as nested
values than as chains. Fluent builders also create a second authoring model next to
ordinary JavaScript composition. The v2 builder can remain useful internally, but the
recommended API should be function-and-object based.

### Alternative 3: Put arbitrary unit renderers on every engine

This is powerful and aligns with the existing `detail(row)` prototype. It is also
easy to overuse. If every engine exposes a fully arbitrary renderer for every part,
then product views become thin pass-throughs and stability moves back to the author.
Named slots keep the good part — author-time composition with JS functions — while
limiting the extension surface.

### Alternative 4: Generate all APIs directly from React manifests

Manifest generation is important for avoiding drift, but manifests describe engines
and components. They do not know product tasks such as "toggle availability" or
"confirm booking" unless we enrich them with domain descriptors. The right answer is
not pure generation from React components; it is descriptors for domain views that
can generate runtime exports, TypeScript declarations, and docs.

---

## 17. Testing and validation

The API redesign should be tested at three levels.

### 17.1 Runtime shape tests

Each domain view should have a test that runs a Goja script and asserts the emitted
IR shape:

```js
const schedule = require("schedule.dsl");
const node = schedule.availabilityPoll(poll, {
  currentResponse: "you",
  onChange: schedule.intent.toggleAvailability("poll.toggle"),
});
JSON.stringify(node);
```

Assertions:

- The returned node is a component node.
- The underlying engine type is `MatrixGrid`.
- Defaults are present: availability states, glyphs, style set, row header, footer.
- The action payload uses the domain context mapping.
- Provided slots are called exactly once per relevant datum.

### 17.2 Declaration tests

The generated declarations should prove that authors can discover slot names and
intent helpers:

```ts
schedule.availabilityPoll(poll, {
  slots: {
    option: ({ option }) => time.slotLabel(option.slot),
    participant: ({ response }) => ui.text(response.name),
  },
  onChange: schedule.intent.toggleAvailability("poll.toggle"),
});
```

The first declaration pass can type slot context as broad `Props`. A later pass
should emit named context interfaces from descriptors.

### 17.3 Golden examples

Keep a small set of example scripts under testdata:

- `availability-poll.js`
- `booking-picker.js`
- `calendar-week.js`
- `record-collection.js`
- `media-library.js`

Each example should execute through Goja, emit JSON, and compare against a stable
golden shape. These examples become living documentation for the grammar.

---

## 18. What an intern should implement first

If an intern starts from this document, the first useful slice is not the whole
redesign. It is the narrowest slice that proves the grammar.

1. Add `ui.raw` as an alias to `component` and update examples to use `raw` only
   for escape hatches.
2. Add child-array flattening to `buildComponent`/child normalization.
3. Add `data.bind` helpers for `field`, `map`, `template`, `context`, and `const`.
4. Add `schedule.intent.toggleAvailability` as a wrapper around `action.server`.
5. Add `schedule.availabilityPoll(poll, options)` with three slots:
   `participant`, `option`, and `footer`.
6. Add tests proving the slot functions are called and the result lowers to a
   `MatrixGrid` node.
7. Add TypeScript declarations for the new view and slots.

That slice is small enough to review and large enough to validate the design. If it
feels good in examples, apply the same pattern to `time.week` and `data.collection`.

---


## 19. Addendum: what the `researchctl` lambda style changes

After drafting the first version of this document, I inspected the JavaScript DSLs
in `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl`. They
are useful because they use Goja lambdas in two distinct ways that the Widget DSL
can learn from.

The first pattern is a **scoped builder callback**. A research project is built as
a chain, but the details of each entity live inside a lambda that receives a narrow
builder:

```js
researchctl.project("Lifecycle example")
  .question("Which scheduler policy should ship?", q => q
    .id("Q-SCHED")
    .status("active")
    .priority("P1")
    .hypothesize("H-MIN-FINISH"))
  .hypothesis("Min-finish-time scheduling reduces p95 latency", h => h
    .id("H-MIN-FINISH")
    .status("open")
    .priority("P1")
    .confidence("medium")
    .testedBy("EXP-SCHED-SWEEP"));
```

The Go implementation is direct. `projectBuilder` creates a `GoalSpec`,
`QuestionSpec`, or `ExperimentSpec`, creates the matching builder object, calls the
JavaScript callback with that builder, then appends the finished spec. The lambda is
not serialized. It is an author-time editing function over a scoped builder. That
matters for the Widget DSL because it gives us a third API shape between raw prop
bags and arbitrary runtime closures.

The second pattern is a **fragment callback**. `codesign` builders expose `.use`, and
ordinary JavaScript functions become reusable conventions:

```js
const baseTopology = r => r
  .experiment("EXP-WORKLOADS")
  .backend("cpu-sim")
  .topology(t => t.cpu("cpu0").gpu("gpu0"))
  .policy("min_finish_time")
  .metrics(m => m.requestCount().latencyP95());

codesign.runSpec("fixed workload")
  .use(baseTopology)
  .workload(w => w.fixed({ count: 4, interarrivalNs: 10 }).stage("infer", {
    computeUnits: 100,
    supportedDevices: ["cpu0", "gpu0"],
  }));
```

This is stronger than copy/paste and lighter than inheritance. A fragment is just a
function from builder to builder. It can be named, tested, combined, and applied in
more than one place.

The third pattern is a **runtime callback registered by ID**. `codesign` allows a
JavaScript callback to participate in simulation:

```js
codesign.runSpec("callback prototype")
  .topology(t => t.jsDevice("js0", (phase, task, state, fallback) => ({
    startNs: fallback.startNs,
    durationNs: 1,
    finishNs: fallback.startNs + 1,
    score: fallback.startNs + 1,
  }), { speed: 100 }))
  .policyCallback("choose-js", () => "js0")
  .metrics(m => m.callback("completed", events => ({
    value: events.filter(e => e.eventType === "request_completed").length,
    unit: "requests",
  })));
```

This third pattern should inspire the Widget DSL carefully. `codesign` can keep the
callback in the same Goja runtime because the simulator executes immediately inside
that runtime. Widget IR is different: it is shipped to a browser. A Goja lambda
cannot cross that boundary. The Widget DSL should therefore use lambdas primarily as
**author-time macros** that lower to serializable IR, and only use callback IDs for
server-side extension points that are explicitly registered and invoked on the
server.

These observations suggest a second API variant that is worth adding to the design:
a **builder-lambda form** for views.

### 19.1 The builder-lambda form

The original proposal in this document uses option objects:

```js
schedule.availabilityPoll(poll, {
  currentResponse: "you",
  slots: {
    option: ({ option }) => time.slotLabel(option.slot),
    participant: ({ response }) => ui.person(response.name),
  },
  onChange: schedule.intent.toggleAvailability("poll.toggle"),
});
```

The `researchctl` style suggests an alternative that may read better for complex
views:

```js
schedule.availabilityPoll(poll, p => p
  .currentResponse("you")
  .readOnly(poll.closed)
  .option(({ option }, h) => h.stack({ gap: "xxs", align: "center" }, [
    h.strong(time.format(option.slot.startISO, "EEE")),
    h.caption(time.formatRange(option.slot.startISO, option.slot.endISO, "HH:mm")),
  ]))
  .participant(({ response }, h) => h.inline({ gap: "xs" }, [
    h.avatar(response.name),
    h.text(response.name),
  ]))
  .empty((_, h) => h.callout("No responses yet."))
  .onChange(schedule.intent.toggleAvailability("poll.toggle"))
  .onSubmit(schedule.intent.submitResponse("poll.submit")));
```

This form keeps the named-slot design, but it moves the customization from one
large options object into a scoped builder. The builder exposes only the operations
that are valid for an availability poll. You cannot accidentally set a raw
`MatrixGrid` prop unless the poll builder intentionally exposes an escape hatch.

The implementation shape mirrors `researchctl`:

```
schedule.availabilityPoll(poll, callback)
  -> create AvailabilityPollSpec with defaults
  -> create AvailabilityPollBuilder bound to that spec
  -> call callback(builder)
  -> validate the completed spec
  -> lower spec to MatrixGrid Widget IR
```

The lambda again is not serialized. It is called while the DSL script runs. The
resulting spec is serializable.

### 19.2 Why builder lambdas are different from arbitrary slot lambdas

A slot lambda returns a node for one subpart. A builder lambda configures the whole
view. The distinction matters because the builder can validate the complete view
before lowering.

```js
schedule.availabilityPoll(poll, p => p
  .currentResponse("you")
  .option(optionHeader)
  .participant(participantCell)
  .onChange(schedule.intent.toggleAvailability("poll.toggle")));
```

The builder can check these rules:

- `currentResponse` must refer to an existing response when the poll is editable.
- `option` must return a renderable node or string.
- `onChange` must be absent when `readOnly(true)` is set, or the builder must lower
  it to a disabled/no-op interaction.
- The view must not expose `allDay` event slots until the underlying calendar view
  supports all-day rendering.

An options object can also be validated, but the builder form creates natural
places for domain-specific verbs. For example, a booking picker can expose
`.duration(30)`, `.timezone("Europe/Berlin")`, `.disablePast()`, and
`.confirmWith(...)`. Those names read as a product grammar rather than as fields in
a large configuration object.

### 19.3 Fragments as reusable view policies

The `.use(fragment)` pattern from `codesign` is the most directly applicable idea.
The Widget DSL needs reusable visual and behavioral policies that are smaller than a
full component and more structured than spreading an object.

```js
const compactScheduling = p => p
  .density("compact")
  .showLegend(false)
  .option(({ option }, h) => h.caption(time.shortSlot(option.slot)));

const editableForCurrentUser = user => p => p
  .currentResponse(user.responseId)
  .onChange(schedule.intent.toggleAvailability("poll.toggle"));

schedule.availabilityPoll(poll, p => p
  .use(compactScheduling)
  .use(editableForCurrentUser(currentUser))
  .participant(({ response }, h) => h.person(response.name)));
```

This is a clean way to apply a house style backward across modules. Instead of
inventing inheritance or theme-specific recipe variants, write fragments:

```js
const quietSection = s => s.tone("quiet").density("comfortable");
const reviewActions = v => v
  .action("Approve", data.intent.approve("review.approve"))
  .action("Reject", data.intent.reject("review.reject"));

ui.page("Review", p => p
  .section("Pending items", s => s
    .use(quietSection)
    .collection(items, c => c
      .use(reviewActions)
      .view("table")
      .detail(({ record }) => data.record(record)))));
```

Fragments are especially attractive because they use ordinary JavaScript lexical
scope. A fragment can close over `currentUser`, `featureFlags`, `palette`, or a
route prefix, but the closure is consumed during lowering. The final Widget IR
contains only the decisions made by that closure.

### 19.4 A page builder variant

The same style can apply to page composition. The earlier document examples used
array children:

```js
ui.page({ title: "Team scheduling" }, [
  ui.section("Find a time", [
    schedule.availabilityPoll(poll, { currentResponse: "you" }),
  ]),
]);
```

The lambda-inspired version reads as a document builder:

```js
ui.page("Team scheduling", page => page
  .section("Find a time", section => section
    .view(schedule.availabilityPoll(poll, pollView => pollView
      .currentResponse("you")
      .onChange(schedule.intent.toggleAvailability("poll.toggle")))))
  .section("Calendar", section => section
    .view(time.week(events, week => week
      .range(time.range.week("2026-07-06"))
      .hours(8, 18)
      .event(({ event }, h) => h.card({ tone: event.colorKey }, [
        h.strong(event.title),
        h.caption(time.formatRange(event.startISO, event.endISO)),
      ]))))));
```

This form is more verbose for small pages, but it scales well when sections have
local defaults, repeated policies, and nested views. It also lets the page builder
provide a scoped helper object `h` to slot callbacks, so slot functions do not need
to import every small UI helper directly.

A practical API can support both forms:

```js
ui.page(options, childrenArray)
ui.page(titleOrOptions, builderCallback)

schedule.availabilityPoll(poll, optionsObject)
schedule.availabilityPoll(poll, builderCallback)
```

The implementation should normalize both into the same internal `ViewSpec`. The
choice becomes a style decision: use options objects for small views and builder
lambdas for views with several customizations or reusable fragments.

### 19.5 Callback IDs for server-side computed views

The `codesign` runtime callback pattern raises one more possibility: server-side
computed view extensions. Most widget lambdas should lower immediately, but there
are cases where the server may need to call a registered function later, before
sending a refreshed widget tree.

A safe version would be explicit:

```js
schedule.registerComputer("poll.rankOptions", ({ poll, tallies }) =>
  tallies.slice().sort((a, b) => b.yes - a.yes).map(t => t.optionId));

schedule.pollSummary(poll, summary => summary
  .orderBy(schedule.computed("poll.rankOptions")));
```

The IR would not contain the JavaScript function. It would contain a registered
callback ID plus provenance metadata, and only a trusted server runtime would invoke
it. This is similar to `codesign.registerCallbackSource`, which exists because a
YAML artifact can record `callbackId` but cannot contain the JavaScript function
body.

This is not a phase-one feature. It is a useful boundary to name now:

- **Builder lambdas** run immediately and lower to IR.
- **Slot lambdas** run immediately while building a domain view and lower to nodes.
- **Registered computed callbacks** may run later on the server, but only by ID and
  only in a trusted runtime.
- **Browser interactions** use serializable action specs, not Goja lambdas.

That boundary lets us leverage JavaScript as a real language without pretending that
closures can be shipped to React.

### 19.6 Revised recommendation

The earlier recommendation was: use domain views, named slots, shared bindings, and
intent actions. After studying `researchctl`, I would revise it to this:

> The Widget DSL should support two equivalent authoring syntaxes over one internal
> view spec: an object form for compact cases and a builder-lambda form for complex
> composition. The builder-lambda form should support `.use(fragment)` everywhere a
> view or section has reusable policy, and it should treat lambdas as author-time
> macros that lower to serializable Widget IR.

That gives us the best parts of both designs. The object form remains simple and
familiar:

```js
schedule.availabilityPoll(poll, {
  currentResponse: "you",
  onChange: schedule.intent.toggleAvailability("poll.toggle"),
});
```

The builder form unlocks reusable fragments and scoped validation:

```js
schedule.availabilityPoll(poll, p => p
  .use(compactScheduling)
  .use(editableForCurrentUser(currentUser))
  .footer(({ tally }, h) => h.caption(`${tally.yes}/${tally.total}`)));
```

The internal representation is the same in both cases. That is the key constraint.
We should not create another split like classic `data.dsl` versus `data.v2.dsl`.
We should create one `AvailabilityPollSpec`, one validator, one lowering path, and
two syntactic ways to fill it.

### 19.7 Implementation slice inspired by `researchctl`

A minimal implementation can copy the proven `researchctl` mechanics:

1. Define an `AvailabilityPollSpec` Go struct that contains `currentResponse`,
   `readOnly`, slot values, action specs, and display options.
2. Implement `availabilityPollBuilder(spec)` returning a Goja object whose methods
   mutate the spec and return the same builder object.
3. Implement `.use(fragment)` on the builder by calling the JavaScript fragment with
   the builder, exactly like `codesign` does.
4. Let `schedule.availabilityPoll(poll, arg)` accept either an options object or a
   builder callback.
5. Normalize both forms into the same spec.
6. Validate the spec.
7. Lower the spec to the existing `MatrixGrid` IR.
8. Add TypeScript declarations with a reusable fragment type:

```ts
type Fragment<T> = (builder: T) => void | T;

interface AvailabilityPollBuilder {
  currentResponse(id: string): this;
  readOnly(value?: boolean): this;
  option(slot: SlotFn<OptionContext>): this;
  participant(slot: SlotFn<ResponseContext>): this;
  footer(slot: SlotFn<TallyContext>): this;
  onChange(action: ActionSpec): this;
  onSubmit(action: ActionSpec): this;
  use(fragment: Fragment<AvailabilityPollBuilder>): this;
}

export function availabilityPoll(
  poll: AvailabilityPollDTO,
  configure?: AvailabilityPollOptions | Fragment<AvailabilityPollBuilder>
): WidgetNode;
```

This slice would prove the lambda design without forcing a rewrite of the whole DSL.
If it works well, apply the same builder pattern to `time.week`, `data.collection`,
and `ui.page`.


## 20. Refactoring the existing DSLs, not only scheduling

The schedule examples are a proving ground, not the scope boundary. The same
redesign should target the existing `ui.dsl`, `data.dsl`, `context_window.dsl`,
`course.dsl`, and `cms.dsl` modules. Those modules already show the symptoms the
new grammar is meant to fix: generic helpers are duplicated across modules, domain
recipes are hand-coded as prop-copying functions, and the authoring style changes
from one module to another.

The current module table makes this visible. `ui.dsl` owns generic layout and
primitive helpers. `cms.dsl` also exports generic helpers such as `breadcrumbs`,
`emptyState`, `meterBar`, `pagination`, `searchField`, `tag`, and `tileGrid`.
`course.dsl` exports `markdownArticle` and `richArticle`, which are also generic
reading primitives. The domain modules need these components, but exporting them as
parallel domain helpers blurs the line between a foundation layer and a product
view layer. A script author has to ask: should a CMS page use `cms.breadcrumbs` or
`ui.breadcrumbs`? Is `course.richArticle` a course concept or just the article
engine?

The refactor should answer that question by making module responsibilities explicit:

| Module | Future responsibility | Compatibility posture |
|---|---|---|
| `ui.dsl` | Page builders, layout, typography, cards, buttons, raw component escape hatch, shared fragment mechanics. | Keep all existing helpers; add builder forms and document `ui.*` as the canonical home for generic components. |
| `data.dsl` | Collections, records, fields, matrices, selection, bindings, and data intents. | Keep `dataTable` and existing `cell.*`; make `collection` and `record` the recommended APIs. |
| `context_window.dsl` | Context snapshots, diagrams, transcript workspaces, annotations, and context-window intents. | Keep current panel helpers; add task-level views such as `context.workspace` and `context.transcript`. |
| `course.dsl` | Course studio, lessons, slide decks, handouts, navigation, learner/presenter intents. | Keep current panel helpers and recipes; add builder-lambda views over them. |
| `cms.dsl` | Media libraries, article queues, editorial workflows, assets, upload/review/publish intents. | Keep current helpers and recipes; add semantic views and move generic visuals to `ui.dsl` in examples. |

This is a redesign of the DSL surface, not a redesign of the React component
library. Most new APIs can lower to the same component types already used by the
hand-written recipes: `MediaLibraryPanel`, `ArticleListPanel`, `CourseStudioShell`,
`CourseSlidePanel`, `HandoutDocumentShell`, `TranscriptWorkspacePanel`, and
`ContextDiagramPanel`.

### 20.1 A shared composition kernel

The first implementation target should be a shared kernel inside `pkg/widgetdsl`,
not a schedule-specific helper. The kernel supplies the behaviors every module uses:

```go
type Fragment[T any] func(T) T

type ViewDescriptor struct {
    Module  string
    Name    string
    Engine  string
    Slots   []SlotDescriptor
    Intents []IntentDescriptor
}
```

In Go this will not be generic TypeScript-style code; it will be a small set of
runtime helpers:

- `applyBuilderCallback(builder, cb)` calls a JavaScript lambda with a scoped
  builder, following the `researchctl` and `codesign` pattern.
- `applyFragment(builder, fragment)` powers `.use(fragment)` everywhere.
- `normalizeChildren` accepts variadic children and array children.
- `normalizeSlot` accepts a JavaScript function, a node, a string, or a spec.
- `callSlot` invokes an author-time slot and validates that the return value can be
  embedded in Widget IR.
- `normalizeAction` keeps existing transport actions working while allowing domain
  intent wrappers.
- `installViewDescriptor` installs runtime exports and gives the TypeScript
  generator enough information to describe the view.

Once this kernel exists, schedule, CMS, course, and context-window views are just
view descriptors plus lowering functions. That is the point of the refactor. The
composition style should not be reimplemented per module.

### 20.2 `ui.dsl`: from component bag to document builder

`ui.dsl` should become the canonical document-composition layer. Existing helpers
such as `panel`, `stack`, `inline`, `button`, `caption`, and `sectionBlock` remain,
but examples should shift toward page and section builders.

Current style:

```js
const ui = require("ui.dsl");

module.exports = ui.page({
  id: "cms-review",
  title: "CMS Review",
  sections: [
    ui.panel({ title: "Needs review" },
      ui.caption("12 articles waiting for editorial sign-off")),
  ],
});
```

Target style:

```js
const ui = require("ui.dsl");

module.exports = ui.page("CMS Review", page => page
  .section("Needs review", section => section
    .metric("Articles", 12, { caption: "waiting for editorial sign-off" })
    .action("Open queue", ui.intent.navigate("/cms/review"))));
```

The builder form gives `ui.dsl` a natural place for reusable page-level policies:

```js
const compactAdminPage = page => page
  .density("compact")
  .chrome("admin")
  .breadcrumbs(true);

ui.page("Media", page => page
  .use(compactAdminPage)
  .section("Library", section => section.view(cms.mediaLibrary(assets))));
```

The compatibility rule is simple: keep the old helper names, but stop teaching them
as the primary composition model. New docs should use `page`, `section`, fragments,
and named views.

### 20.3 `cms.dsl`: from panels to editorial tasks

`cms.dsl` currently exposes both panel-level helpers and recipes:

- `mediaLibraryPanel` and `recipes.mediaLibrary` lower to `MediaLibraryPanel`.
- `articleListPanel` and `recipes.articleList` lower to `ArticleListPanel`.
- Several generic helpers overlap with `ui.dsl`.

The target API should name editorial tasks:

```js
const cms = require("cms.dsl");
const ui = require("ui.dsl");

cms.mediaLibrary(assets, library => library
  .selection("multi", selectedAssetIds)
  .query(query)
  .kindFilter(kind)
  .pagination(page, pageCount)
  .asset(({ asset }, h) => h.card({ tone: asset.status }, [
    h.image(asset.thumbnailUrl),
    h.strong(asset.title),
    h.caption(asset.kind),
  ]))
  .toolbar(toolbar => toolbar
    .action("Upload", cms.intent.upload("cms.upload"))
    .action("Delete", cms.intent.deleteAssets("cms.deleteAssets")))
  .onOpen(cms.intent.openAsset("cms.openAsset")));
```

For articles:

```js
cms.articleQueue(articles, queue => queue
  .statusFilter("needs-review")
  .search(query)
  .article(({ article }, h) => h.inline({ gap: "sm" }, [
    h.status(article.status),
    h.stack({}, [h.strong(article.title), h.caption(article.author)]),
  ]))
  .rowActions(actions => actions
    .action("Preview", cms.intent.previewArticle("cms.previewArticle"))
    .action("Publish", cms.intent.publishArticle("cms.publishArticle")))
  .onSelect(cms.intent.selectArticle("cms.selectArticle")));
```

The old recipe remains as an alias:

```js
cms.recipes.mediaLibrary(options)  // compatibility
cms.mediaLibrary(options.assets, builderOrOptions)  // preferred
```

The new view owns CMS vocabulary: asset selection, upload, query submit, page
change, article preview, publish, archive, and status filtering. It lowers those to
the existing `MediaLibraryPanel` and `ArticleListPanel` action props.

### 20.4 `course.dsl`: from shell components to teaching flows

`course.dsl` should describe learning and presentation tasks rather than shell
components. The current `courseStudio`, `courseSlide`, and `handout` recipes are
good lowering targets, but their public API should expose lesson, slide, handout,
and navigation concepts.

Target course studio:

```js
const course = require("course.dsl");

course.studio(courseSpec, studio => studio
  .active(activeLessonId)
  .nav(nav => nav
    .section("Foundations", section => section
      .lesson("intro")
      .lesson("widget-ir")
      .lesson("actions")))
  .lesson(({ lesson }, h) => h.stack({ gap: "md" }, [
    h.richArticle(lesson.article),
    course.checkpoint(lesson.checkpoint),
  ]))
  .onNavigate(course.intent.navigateLesson("course.navigate")));
```

Target slide deck:

```js
course.slideDeck(deck, deckView => deckView
  .index(currentIndex)
  .mode("presenter")
  .visual(({ slide }, h) => context.diagram(slide.snapshot))
  .notes(({ slide }, h) => h.markdown(slide.speakerNotes))
  .onNext(course.intent.nextSlide("course.next"))
  .onPrevious(course.intent.previousSlide("course.previous"))
  .onPresent(course.intent.present("course.present")));
```

Target handout:

```js
course.handout(bundle, handout => handout
  .selected(selectedDocumentId)
  .document(({ document }, h) => h.card({}, [
    h.strong(document.title),
    h.caption(document.description),
  ]))
  .onSelect(course.intent.selectDocument("course.selectDocument"))
  .onDownload(course.intent.downloadDocument("course.downloadDocument")));
```

This gives course scripts the same composition grammar as schedule scripts:
semantic view, scoped builder, named slots, and intent actions. The existing panel
helpers remain as engine helpers for advanced users.

### 20.5 `context_window.dsl`: from visualization panels to analysis workspaces

The context-window module is already domain-specific, but it still exposes many raw
panel helpers. The redesign should introduce task-level views over those helpers:

```js
const context = require("context_window.dsl");

context.workspace(session, workspace => workspace
  .snapshot(snapshot)
  .transcript(transcript)
  .annotations(annotations)
  .layout("diagram-plus-transcript")
  .diagram(diagram => diagram
    .view("stack")
    .selected(selectedPartId)
    .legend(true))
  .message(({ message }, h) => h.transcriptMessage(message, {
    showTokenCount: true,
  }))
  .annotation(({ annotation }, h) => h.annotationCard(annotation))
  .onAnnotationSelect(context.intent.selectAnnotation("context.selectAnnotation"))
  .onPartSelect(context.intent.selectPart("context.selectPart")));
```

The lower-level APIs still matter:

```js
context.contextDiagramPanel({ snapshot, styleSet });
context.transcriptReaderPanel({ messages });
```

But the recommended API should express the full analysis workspace. This is where
named slots help: message rendering and annotation rendering are safe seams;
internal context-budget layout is not something every script should rebuild.

### 20.6 `data.dsl`: the bridge between domain modules and engines

`data.dsl` should become the shared data-view vocabulary used by CMS, course, and
schedule when they display collections. The target is not to remove tables. The
target is to make tables one view mode of a collection.

```js
data.collection(articles, collection => collection
  .key("id")
  .fields(fields => fields
    .text("title", { label: "Title" })
    .status("status", { palette: cms.articleStatusPalette() })
    .date("updatedAt", { label: "Updated" }))
  .view("table")
  .selection("single", selectedArticleId)
  .detail(({ record }) => cms.articlePreview(record))
  .onSelect(data.intent.selectRecord("cms.selectArticle")));
```

CMS can then decide whether `cms.articleQueue` lowers directly to `ArticleListPanel`
or to `data.collection` plus CMS slots. The important design point is that the
collection grammar becomes a shared middle layer rather than another isolated DSL.

### 20.7 Refactor roadmap for all modules

The roadmap should be shared across modules so we do not create one new style for
schedule and leave the old modules behind.

1. **Inventory current exports.** Generate a table from `moduleSpecs`: module,
   helper name, component type, recipe name, and whether the helper is generic or
   domain-specific.
2. **Classify helpers.** Move documentation ownership of generic helpers to
   `ui.dsl`; keep old domain aliases as compatibility exports.
3. **Add the composition kernel.** Implement child normalization, builder callback
   application, `.use(fragment)`, slot invocation, and descriptor metadata once.
4. **Add builder forms to `ui.dsl`.** Prove the style on `page`, `section`, and a
   few layout helpers before touching domain modules.
5. **Wrap existing recipes as domain views.** Add `cms.mediaLibrary`,
   `cms.articleQueue`, `course.studio`, `course.slideDeck`, `course.handout`,
   `context.workspace`, and `context.diagram` as preferred APIs over existing panel
   components.
6. **Add intent namespaces.** Add `cms.intent.*`, `course.intent.*`,
   `context.intent.*`, and keep existing `action.*` as transport-level primitives.
7. **Unify declarations.** Emit TypeScript declarations from view descriptors,
   including builder interfaces and fragment types.
8. **Write golden examples.** Every module should have one object-form example and
   one builder-lambda example that execute through Goja and compare emitted IR.
9. **Update documentation.** Teach semantic views first, engine helpers second,
   raw escape hatch last.
10. **Deprecate only in docs at first.** Do not remove old helpers until scripts have
    migrated and tests cover the new surface.

### 20.8 Concrete acceptance criteria

The broader refactor is successful when these scripts are possible and tested:

```js
ui.page("Editorial dashboard", page => page
  .use(adminChrome)
  .section("Media", s => s.view(cms.mediaLibrary(assets, mediaDefaults)))
  .section("Articles", s => s.view(cms.articleQueue(articles, reviewQueue))));
```

```js
course.studio(courseSpec, studio => studio
  .use(courseChrome)
  .active(activeLessonId)
  .lesson(lessonRenderer)
  .onNavigate(course.intent.navigateLesson("course.navigate")));
```

```js
context.workspace(session, workspace => workspace
  .snapshot(snapshot)
  .transcript(transcript)
  .annotation(annotationRenderer)
  .onAnnotationSelect(context.intent.selectAnnotation("context.select")));
```

```js
schedule.availabilityPoll(poll, pollView => pollView
  .use(compactScheduling)
  .currentResponse(currentUser.responseId)
  .onChange(schedule.intent.toggleAvailability("poll.toggle")));
```

These examples should all lower to existing Widget IR component nodes. If they
require new React components before the DSL can be tested, the design has drifted
from its purpose. The refactor should first improve the authoring layer over what
already exists, then expose new engines as they arrive.

## 21. Closing model

The redesigned DSL should feel less like a bag of component constructors and more
like a small language for product screens. The language has pages and sections for
structure, domain views for product concepts, slots for controlled composition,
bindings for data flow, and intents for interactions. Engines still exist, and raw
IR construction still exists, but they are no longer the first thing an author sees.

The scheduling work is the right place to make this shift because it contains the
whole problem in miniature. An availability poll is generic enough to reuse a matrix
engine, specific enough to need product defaults, interactive enough to need action
context, and visual enough to need customization seams. If the DSL can express that
cleanly, the same grammar can carry older tables, records, media libraries, course
views, and future widgets.

The final test is simple: a new intern should be able to read a script and answer
three questions without opening the React code:

- What product view is this screen showing?
- Which parts did the script customize?
- What user intents can leave the screen?

If the answer is visible in the JavaScript, the DSL is doing its job.

---

## References

- `design-doc/02-goja-dsl-layer-design-and-implementation-guide-for-scheduling-widgets.md` — implementation-level guide for wiring scheduling/calendar engines into the existing Goja module system.
- `reference/02-scheduling-and-calendar-widgets-implementation-and-ir-handoff-for-dsl-wiring.md` — concise IR handoff for scheduling/calendar widgets.
- `reference/03-code-review-and-design-assessment-for-scheduling-widgets.md` — review findings, especially `readOnly`, `allDay`, exports, and focused tests.
- `RAGEVAL-WIDGET-DECOMPOSITION/design-doc/02-redesigning-the-widget-dsl-a-composition-first-opinionated-javascript-api.md` — colleague proposal centered on engine verbs and unit renderers; this document keeps the composition goal but narrows it through named slots and domain views.
- `pkg/widgetdsl/module.go` — current module specs, helper maps, `componentFactory`, `cellObject`, `actionObject`, recipes, and the existing `detail(row)` slot prototype.
- `pkg/widgetdsl/typescript.go` — generated TypeScript declaration surface that must evolve with any API redesign.
- `packages/rag-evaluation-site/src/widgets/ir/engines.ts` — engine contracts the redesigned DSL lowers into.
- `packages/rag-evaluation-site/src/widgets/presets/scheduling.ts` — TypeScript scheduling presets that inspired the domain-view layer.
