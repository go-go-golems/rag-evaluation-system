---
Title: "Widget DSL v3 Examples"
Slug: widget-dsl-v3-examples
Short: "Build v3 Widget IR pages with composition, actions, keyboard shortcuts, scheduling, time, and CRM helpers."
Topics:
- xgoja
- widget-dsl
- widget-ir
- scheduling
- crm
- javascript
Commands:
- xgoja build
- xgoja help
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

`widget.dsl` lets a JavaScript host describe a page as serializable Widget IR while the React application owns rendering and interaction behavior. This tutorial uses the same builder-lambda style as the checked golden examples, so each snippet can be moved into a jsverb route after the host selects the `widget.dsl` runtime module.

## Start with a page and typed UI helpers

A v3 page is the composition boundary: `widget.page` creates the page wrapper, a section builder supplies labels and view slots, and domain helpers create the leaf Widget IR. Prefer a typed helper whenever one exists because it fixes the component contract in one place; use `widget.raw.component` only for a component without a typed v3 helper.

```js
const widget = require("widget.dsl")

const page = widget.page("Workshop overview", (p) =>
  p.id("workshop-overview")
   .section("Welcome", (s) =>
     s.caption("A server-produced page rendered by React.").view(
       widget.ui.card({ title: "AI engineering workshop" },
         widget.ui.caption("Two days, hands-on."),
         widget.ui.button("Request a workshop", widget.act.navigate("/pages/lead"), {
           variant: "primary",
         }),
       ),
     ),
   ),
)
```

A page builder returns a handle. Route handlers return `page.toPage()` when they need a plain page object, while page-aware xgoja helpers may accept the handle directly.

## Add page-level keyboard shortcuts

Page shortcuts are serializable command accelerators for workflows that are not naturally table-shaped. Each binding gives the command a stable id, a logical key, a visible help label, and an ordinary action. Reuse the same action object in the shortcut and visible button so keyboard and pointer activation cannot drift.

```js
const accept = widget.act.server("triage.accept")
const reject = widget.act.server("triage.reject")
const skip = widget.act.server("triage.skip")

const page = widget.page("Triage", (page) =>
  page
    .shortcuts((keys) =>
      keys
        .bind("accept", "y", accept, { label: "Yes" })
        .bind("reject", "n", reject, { label: "No" })
        .bind("skip", "s", skip, { label: "Skip" }),
    )
    .section("Current job", (section) =>
      section.view(widget.ui.inline(
        widget.ui.button("Yes", accept, { variant: "success" }),
        widget.ui.button("No", reject, { variant: "danger" }),
        widget.ui.button("Skip", skip),
      )),
    ),
)
```

Use separate modifier values rather than embedding them in the key string:

```js
page.shortcuts((keys) =>
  keys.bind("save", "s", widget.act.server("record.save"), {
    label: "Save",
    modifiers: ["Control"],
  }),
)
```

The default React app displays generated shortcut help and a persisted enable/disable control. It suppresses page commands in editable fields, dialogs, IME composition, repeated keydown events, and nested keyboard scopes. Keep table row commands on `table.command(...)`; page shortcuts should not replace component-owned keyboard interaction.

The executable version is `pkg/widgetdsl/testdata/v3/examples/43-page-shortcuts.js`, with stable output checked by its golden JSON fixture.

## Bind interaction context to an action

Bindings defer a value lookup until the user interacts with a rendered widget. This keeps host data serializable and prevents a page from guessing which grid cell or board card the user will select.

```js
const poll = {
  title: "Office hours",
  options: [{ id: "mon-9", label: "Monday 09:00" }],
  responses: [{ id: "ana", name: "Ana", availability: { "mon-9": "available" } }],
}

const pollView = widget.schedule.availabilityPoll(poll, (b) =>
  b.editableRow("ana").onToggle(
    widget.schedule.intent.toggleAvailability(
      widget.bind.context("row.id"),
      widget.bind.context("column.id"),
    ),
  ),
)
```

MatrixGrid provides `row`, a serializable `column`, `rowKey`, `colId`, and `value` to the action context. `schedule.intent.toggleAvailability` emits an event action, so the browser listener receives the resolved identifiers in `event.detail`. A server action instead sends resolved values to `/api/widget/actions/{name}` in `payload`.

## Render scheduling and time views

Scheduling helpers lower poll-shaped data to MatrixGrid, and time helpers lower calendar data to MonthGrid or TimeGrid. The helpers keep the page script focused on domain data and interaction intent instead of renderer-specific cell specifications.

```js
const week = widget.time.week([
  {
    id: "kickoff",
    title: "Workshop kickoff",
    startISO: "2026-07-14T09:00:00Z",
    endISO: "2026-07-14T10:00:00Z",
    styleKey: "event",
  },
], (b) =>
  b.range(widget.time.range.week("2026-07-14"))
   .hours(8, 18)
   .onSelect(widget.time.intent.selectEvent(widget.bind.context("block.id"))),
)
```

Use `widget.schedule.availabilityPoll`, `pollSummary`, and `bookingPicker` for table-like availability data. Use `widget.time.month` and `widget.time.week` for dates and timed blocks; TimeGrid deliberately does not model all-day blocks.

## Build a CRM pipeline and record form

CRM builders are opaque only while defining a field schema or pipeline. Deals, activities, task rows, and summaries remain plain JavaScript records, which lets a SQLite-backed jsverb load and return them without a separate client model.

```js
const fields = widget.crm.fields("Workshop opportunity", (f) =>
  f.text("organization", { label: "Organization", group: "Customer" })
   .email("buyerEmail", { label: "Buyer email", group: "Customer" })
   .currency("amount", { label: "Expected value", group: "Commercial", unit: "USD" }),
)

const pipeline = widget.crm.pipeline("Workshop sales", (p) =>
  p.stage("lead", "New lead", { colorKey: "lead" })
   .stage("proposal", "Proposal", { colorKey: "proposal" })
   .stage("won", "Won / scheduled", { colorKey: "won" }),
)

const pipelineView = widget.crm.pipelineBoard(pipeline, deals, (b) =>
  b.summaries(stageSummaries)
   .onMove(widget.crm.intent.moveDeal("${cardId}", "${to}"))
   .onOpen(widget.crm.intent.openDeal("${cardId}")),
)

const recordView = widget.crm.recordFields(deals[0].fields, fields, (b) =>
  b.mode("edit")
   .onChange(widget.crm.intent.updateField("deal-acme", "${key}", "${value}")),
)
```

BoardEngine supplies `cardId` for selection and `cardId`, `from`, `to`, and `beforeId` for moves. CRM intent placeholders become typed action payload bindings; they are not ordinary string interpolation. `widget.crm.funnel(pipeline, stageSummaries)` also accepts sparse summaries and displays missing stages as zero.

## Serve the page from xgoja

An xgoja host selects the provider module at build time and returns Widget IR from an API route. The SPA asset handler must exclude `/api` so it does not return `index.html` in place of page JSON.

```yaml
runtime:
  modules:
    - provider: rag-widget-site
      name: widget.dsl
      as: widget.dsl
```

```js
app.get("/api/widget/pages/pipeline", (_req, res) => {
  res.json(widget.page("Pipeline", (p) =>
    p.section("Opportunities", (s) => s.view(pipelineView)),
  ).toPage())
})
```

The workshop CRM vertical slice in `examples/xgoja/workshop-crm-site/` demonstrates lead creation, pipeline movement, availability selection, and workshop-run scheduling with SQLite persistence.

## Troubleshooting

Most authoring failures occur at the boundary between build-time module selection, server action contracts, and the SPA fallback. Match the symptom below to its boundary before changing the Widget IR itself.

| Problem | Cause | Solution |
| --- | --- | --- |
| `Cannot find module "widget.dsl"` | The build specification does not select the provider runtime module. | Add the `runtime.modules` entry shown above and rebuild the xgoja binary. |
| A grid action receives an accessor object rather than an id. | An accessor was placed in ordinary application data rather than an action `payload` or `detail`. | Use `widget.bind.context(...)` only inside an action contract. |
| An event listener gets no action values. | The listener reads `payload`, but browser event actions dispatch values in `CustomEvent.detail`. | Read `event.detail`; reserve `payload` for server actions. |
| `page.shortcuts is not a function` | The generated host uses an older `widget.dsl` provider. | Rebuild the xgoja host against the updated `rag-widget-site` provider. |
| A shortcut does not run while focus is in a table or form. | Component keyboard scopes and editable controls intentionally take precedence. | Use `table.command(...)` for row actions and keep page shortcuts outside editing contexts. |
| A single-character shortcut is active but no disable control is visible. | A custom host installed behavior without the default help/preference UI. | Use `RagEvaluationSiteApp` or provide equivalent `usePageShortcuts` help and preference integration. |
| A pipeline stage is empty or has an invalid width. | The host did not provide a summary for that stage. | Pass sparse summaries safely; `widget.crm.funnel` defaults missing counts to zero. |
| The browser receives `index.html` from an API route. | The SPA fallback catches `/api`. | Configure `spaFromAssetsModule` with `excludePrefixes: ["/api"]`. |

## See Also

These related entries separate host configuration, API discovery, and executable examples so a reader can move from a first page to a complete application without relying on ticket-local notes.

- `widget-dsl-getting-started` — select the provider and write a first page.
- `widget-dsl-v3-api-reference` — descriptor-derived namespace inventory.
- `widget-dsl-js-api-reference` — action contracts and legacy migration details.
- `widget-dsl-spa-bundling` — bundle provider help and browser assets.
- [`pkg/widgetdsl/testdata/v3/examples`](../../../../widgetdsl/testdata/v3/examples) — executable golden examples, including page shortcuts, scheduling, time, and CRM.
- [`examples/xgoja/workshop-crm-site`](../../../../../examples/xgoja/workshop-crm-site) — SQLite-backed CRM host.
