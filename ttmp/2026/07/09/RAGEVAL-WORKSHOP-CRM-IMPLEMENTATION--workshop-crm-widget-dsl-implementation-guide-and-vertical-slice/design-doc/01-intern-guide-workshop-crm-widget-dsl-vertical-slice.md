---
Title: ""
Ticket: ""
Status: ""
Topics: []
DocType: ""
Intent: ""
Owners: []
RelatedFiles:
    - Path: repo://examples/xgoja/doodle-site/verbs/doodle.js
      Note: Reference xgoja route entrypoint
    - Path: repo://examples/xgoja/doodle-site/verbs/lib/store.js
      Note: Reference SQLite persistence boundary
    - Path: repo://packages/rag-evaluation-site/src/crm/types.ts
      Note: Generic CRM DTOs and field model
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/engines.ts
      Note: CRM and time Widget IR contracts
    - Path: repo://packages/rag-evaluation-site/src/widgets/presets/crm.ts
      Note: Existing CRM-to-IR parity reference
    - Path: repo://pkg/widgetdsl/module.go
      Note: Widget DSL v3 root namespace installation
    - Path: repo://pkg/widgetdsl/typescript.go
      Note: Typed JavaScript declaration reference
    - Path: repo://pkg/widgetdsl/v3.go
      Note: Goja v3 builder and helper implementation reference
ExternalSources: []
Summary: ""
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---


# Intern Guide: Workshop CRM Widget DSL Vertical Slice

## 1. Executive summary

This guide explains how to add a workshop-focused customer relationship management application to the Widget DSL v3 system. The application is for two facilitators delivering AI engineering workshops. Its first usable workflow is deliberately narrow: capture a lead, inspect and move its opportunity through a pipeline, coordinate dates with a client, create a scheduled workshop run, and see the run in a calendar.

The design does not create a second frontend architecture. JavaScript authors use `require("widget.dsl")` to create a serializable Widget IR tree. The xgoja HTTP host executes that JavaScript and returns JSON. The React application fetches the JSON and renders it through `WidgetRenderer` and the default widget registry. The proposed `widget.crm` namespace is a Go-side authoring convenience layer over existing CRM IR components. It must emit `BoardEngine`, `RecordFieldList`, `ActivityFeed`, `StatTile`, and `SegmentedBar`; it must not return JSX or import the React package.

This document is written for an intern who needs to understand both the architectural boundary and the first implementation sequence. It specifies concrete API shapes, data contracts, file boundaries, validation, and review criteria.

## 2. Scope

### In scope

1. A `widget.crm` v3 namespace with typed/fluent helpers.
2. CRM schema and pipeline definition builders stored as opaque Goja-backed objects.
3. CRM views that emit existing Widget IR nodes.
4. A new SQLite-backed xgoja example application, `workshop-crm-site`.
5. Pages for dashboard, pipeline, opportunity record, client availability, workshop run, and lead intake.
6. Golden fixtures, WidgetRenderer stories, Go tests, and browser smoke coverage.

### Out of scope

- Authentication, multi-tenant authorization, invoicing, email delivery, payment collection, or a real proposal PDF generator.
- Replacing the existing TypeScript CRM presets immediately.
- A separate `calendar.dsl` module. Calendar support remains under `widget.time` and `widget.schedule`.
- A generic no-code CRM builder.

## 3. Terms

- **Widget DSL**: Goja JavaScript APIs exposed by the Go `pkg/widgetdsl` module.
- **Widget IR**: JSON-compatible semantic tree with component `type`, `props`, and `children`.
- **Builder**: fluent JavaScript object whose methods mutate a Go-owned specification and return the same builder.
- **Opaque object**: a Goja object carrying a private Go reference. JavaScript can call its documented methods but cannot depend on its internal representation.
- **CRM record**: a customer-facing business object such as an organization, contact, opportunity, or workshop run.
- **Workshop run**: a delivery engagement created after an opportunity is won or scheduled.

## 4. Current architecture, with evidence

### 4.1 JavaScript to rendered UI

```text
JavaScript verb
  require("widget.dsl")
          |
          v
widget.page + widget.crm/time/schedule/ui builders
          |
          v
Widget IR JSON
          |
          v
xgoja HTTP response: /api/widget/pages/<page>
          |
          v
React app fetches Widget IR
          |
          v
WidgetRenderer + defaultWidgetRegistry
          |
          v
registered React widget adapters
          |
          v
DOM/UI
```

The v3 module is installed in `pkg/widgetdsl/module.go` in `installWidgetV3`. Today it exports `page`, `raw`, `act`, `bind`, `ui`, `data`, `cms`, `course`, `context`, `schedule`, `time`, and `style`. Add `crm` in that one installation function.

`pkg/widgetdsl/v3.go` owns the Goja-facing v3 implementations. `v3Page` creates a Go-owned `v3PageSpec`, passes an opaque builder to a callback, then `toPage()` serializes it. `v3Fields` and `v3Collection` demonstrate opaque references through `attachV2Ref` and `mustV2Ref`. The CRM namespace should follow that style for definitions that require validation and ordered mutation.

The React registry is assembled in `packages/rag-evaluation-site/src/widgets/defaultRegistry.ts`. Any IR type emitted by `widget.crm` must already have a registered adapter. The current CRM adapters are `BoardEngine`, `RecordFieldList`, `FieldRenderer`, `ActivityFeed`, `StatTile`, and `SegmentedBar`.

### 4.2 Existing CRM domain and IR

`packages/rag-evaluation-site/src/crm/types.ts` defines the current domain nouns:

- `Contact`, `Company`, `Deal`, `Pipeline`, and `Stage`.
- `Activity` and `Task` for operational history.
- `FieldDef`, `FieldOption`, and `FieldValue` for workspace-defined fields.

`packages/rag-evaluation-site/src/widgets/ir/engines.ts` defines serializable renderer contracts:

```text
BoardEngineWidgetProps
  columns + cards + columnField + card cell spec + actions

RecordFieldListWidgetProps
  values + typed FieldSpec sections + read/edit mode + field action

ActivityFeedWidgetProps
  activities + glyph map + style set + open/load-more actions

StatTileWidgetProps
  label + value + metric trend/progress/tone

MonthGridWidgetProps / TimeGridWidgetProps
  month markers or week blocks + selection + actions
```

`packages/rag-evaluation-site/src/widgets/presets/crm.ts` is the immediate implementation reference. It already translates CRM DTOs into Widget IR through `pipelineBoard`, `recordFieldList`, `activityFeed`, `contactRecord`, `statTile`, `pipelineFunnel`, `crmDashboard`, and `tasksInbox`. The first Go `widget.crm` implementation should preserve these semantics, then the project can decide whether to consolidate the duplicate mapping later.

### 4.3 Existing time, schedule, and generic UI APIs

`pkg/widgetdsl/v3.go` exposes:

```js
widget.time.month(eventsOrMarkers, builder)
widget.time.week(events, builder)
widget.time.range.week(anchorISO)
widget.schedule.availabilityPoll(poll, builder)
widget.schedule.pollSummary(poll, tallies, builder)
widget.schedule.bookingPicker(availability, builder)
widget.ui.splitPane(left, right, options)
widget.ui.form(options, ...children)
widget.ui.formRow(label, control, options)
widget.ui.textInput(options)
widget.ui.selectInput(options)
widget.ui.status(status, value, options)
widget.ui.shareLink(href, options)
```

`packages/rag-evaluation-site/src/components/molecules/TimeGrid/TimeGrid.widget.tsx` is evidence that `viewportHeight` reaches the React component through the IR `style` prop. `packages/rag-evaluation-site/src/components/molecules/MonthGrid/MonthGrid.widget.tsx` is evidence that month selection and navigation are actions, not client-local application state.

## 5. Desired workshop CRM domain

The existing generic CRM DTOs are sufficient for organizations, contacts, opportunities, activities, and tasks. Add application-local SQLite records for workshop delivery:

```text
Organization 1 ──── * Contact
Organization 1 ──── * Opportunity
Opportunity  1 ──── 0..1 WorkshopRun
WorkshopRun  1 ──── * WorkshopSession
WorkshopRun  1 ──── * Task
WorkshopRun  1 ──── * Activity
WorkshopRun  1 ──── * AssetLink
```

Suggested `WorkshopRun` shape:

```js
{
  id: "run-acme-july",
  opportunityId: "deal-acme",
  title: "Acme Robotics — Agentic Coding Bootcamp",
  status: "scheduled", // draft | scheduled | delivered | canceled
  startISO: "2026-07-23T09:00",
  endISO: "2026-07-24T17:00",
  timeZone: "Europe/Amsterdam",
  format: "onsite-2d",
  location: "Acme Robotics, Delft",
  facilitatorIds: ["manuel", "friend"],
  participantCount: 24,
  setupPath: "/pages/workshop?run=run-acme-july"
}
```

Do not add a widget specifically for `WorkshopRun` in the first slice. Render it through CRM fields, time events, tasks, course content, CMS assets, and generic UI composition. This keeps widget ownership separate from business data.

## 6. Proposed `widget.crm` API

### 6.1 Namespace

```js
const widget = require("widget.dsl");
const crm = widget.crm;
```

The first public surface is:

```js
crm.fields(name?, configure?)
crm.pipeline(nameOrOptions, configure?)
crm.pipelineBoard(pipeline, deals, configure?)
crm.recordFields(values, schema, configure?)
crm.activityFeed(activities, configure?)
crm.tasksInbox(tasks, configure?)
crm.stat(label, value, options?)
crm.funnel(pipeline, summaries, options?)
crm.intent.openDeal(id)
crm.intent.moveDeal(id, toStage)
crm.intent.updateField(recordId, key, value)
crm.intent.completeTask(taskId)
```

Do not initially add `crm.dashboard`. A dashboard is a page-level composition of `crm.stat`, `crm.funnel`, `crm.activityFeed`, `crm.tasksInbox`, `widget.time.week`, and generic layout. This avoids hiding product-specific dashboard choices behind a prematurely generic helper.

### 6.2 Opaque definition builders

Use opaque builders only for definitions with invariants.

```js
const opportunityFields = crm.fields("Workshop opportunity", (f) =>
  f.text("organization", { label: "Organization", group: "Customer", required: true })
   .relation("buyerId", "contact", { label: "Buyer", group: "Customer" })
   .currency("amount", { label: "Expected value", unit: "USD", group: "Commercial" })
   .date("targetDate", { label: "Target delivery", group: "Workshop" })
   .select("format", {
     label: "Format",
     group: "Workshop",
     options: [
       { value: "onsite-1d", label: "1-day onsite", colorKey: "onsite" },
       { value: "onsite-2d", label: "2-day onsite", colorKey: "onsite" },
       { value: "remote", label: "Remote", colorKey: "remote" }
     ]
   })
   .tags("topics", { label: "Topics", group: "Workshop" })
);

const pipeline = crm.pipeline("AI engineering workshops", (p) =>
  p.stage("lead", "New lead", { colorKey: "lead", probability: 0.05 })
   .stage("discovery", "Discovery", { colorKey: "discovery", probability: 0.2 })
   .stage("proposal", "Proposal", { colorKey: "proposal", probability: 0.45 })
   .stage("contracting", "Contracting", { colorKey: "contracting", probability: 0.75 })
   .stage("won", "Won / scheduled", { colorKey: "won", probability: 1 })
);
```

Expected builder behavior:

- Mutator methods return the same object (`this`).
- `build()` returns a plain serializable snapshot.
- `validate()` returns structured DSL validation issues.
- `pipelineBoard` and `recordFields` accept the opaque object or a built snapshot.
- A builder must reject blank IDs, duplicate field keys, duplicate stage IDs, and non-monotonic stage order.

### 6.3 Render helpers

```js
const board = crm.pipelineBoard(pipeline, deals, (b) =>
  b.summaries(stageSummaries)
   .selected(selectedDealId)
   .onMove(crm.intent.moveDeal("${dealId}", "${toStage}"))
   .onOpen(crm.intent.openDeal("${dealId}"))
);

const fields = crm.recordFields(opportunity.fields, opportunityFields, (r) =>
  r.mode("edit")
   .refs(referenceMap)
   .onChange(crm.intent.updateField(opportunity.id, "${key}", "${value}"))
);

const history = crm.activityFeed(activities, (a) =>
  a.groupByDay(true)
   .onOpen({ kind: "event", event: "activity.open" })
);
```

`pipelineBoard` must emit `BoardEngine`. `recordFields` must emit `RecordFieldList`. `activityFeed` must emit `ActivityFeed`. The helper must never use `widget.raw.component(...)` internally.

### 6.4 Intent helpers

Intent helpers reduce repetitive, inconsistent action maps while retaining explicit server contracts.

```js
crm.intent.moveDeal("${dealId}", "${toStage}")
// {
//   kind: "server",
//   name: "crm.deal.move",
//   payload: { dealId: "${dealId}", toStage: "${toStage}" }
// }

crm.intent.updateField("deal-acme", "${key}", "${value}")
// {
//   kind: "server",
//   name: "crm.field.update",
//   payload: { recordId: "deal-acme", key: "${key}", value: "${value}" }
// }
```

The xgoja host receives this action through the standard WidgetRenderer action dispatch path. The application route owns persistence and authorization; the design-system component remains API-free.

## 7. The vertical slice

### 7.1 User journey

```text
Inbound lead
  -> Lead intake form
  -> Opportunity in "New lead"
  -> Drag or action moves to "Discovery"
  -> Client availability poll selects dates
  -> Opportunity moved to "Won / scheduled"
  -> Server creates WorkshopRun + calendar events + logistics tasks
  -> Workshop run page shows agenda, schedule, assets, and setup share link
```

### 7.2 Page routes

```text
/pages/dashboard
/pages/pipeline
/pages/opportunity?deal=<id>
/pages/calendar?day=<YYYY-MM-DD>&event=<id>
/pages/workshop?run=<id>
/pages/new-lead
/api/widget/pages/<page>
```

Selection remains URL-backed. A selected opportunity, day, calendar event, or workshop run must survive reload and be linkable.

### 7.3 Pseudocode: page route composition

```js
function calendarPage(req) {
  const selectedDay = req.query.day || todayISO();
  const selectedEvent = req.query.event || null;
  const events = store.workshopCalendarEvents(selectedDay);

  const month = widget.time.month(events, (m) =>
    m.selected(selectedDay)
     .today(todayISO())
     .onSelect(widget.time.intent.selectDay("${dayISO}"))
  );

  const week = widget.time.week(events, (w) =>
    w.range(widget.time.range.week(selectedDay))
     .hours(8, 20)
     .viewportHeight(420)
     .selected(selectedEvent)
     .onSelect(widget.time.intent.selectEvent("${eventId}"))
  );

  return widget.page("Workshop calendar", (p) =>
    p.id("calendar")
     .view(widget.ui.splitPane(month, week, {
       ratio: "leftNarrow", gutter: "lg", divider: true
     }))
     .toPage()
  );
}
```

### 7.4 Pseudocode: server-side opportunity conversion

```text
function scheduleWonOpportunity(dealID, selectedDates):
  transaction:
    deal = loadDeal(dealID)
    require deal.stage == "won"

    run = insertWorkshopRun(
      opportunityID = deal.id,
      title = deal.organization + " — " + deal.workshopTitle,
      start = selectedDates.start,
      end = selectedDates.end,
      status = "scheduled"
    )

    insertCalendarEvent(run.id, run.start, run.end, "workshop")
    insertTask(run.id, "Confirm room and equipment", due = run.start - 14d)
    insertTask(run.id, "Send participant setup link", due = run.start - 10d)
    insertActivity(deal.id, "stage_change", "Workshop scheduled")
  commit
```

## 8. File layout

```text
pkg/widgetdsl/
  v3.go                         # widget.crm Goja namespace and builders
  typescript.go                 # widget.crm declarations
  v3_crm_test.go                # namespace/unit tests
  testdata/v3/examples/
    41-crm-pipeline.js
    42-crm-record.js
  testdata/v3/golden/
    41-crm-pipeline.json
    42-crm-record.json

packages/rag-evaluation-site/src/
  widgets/ir/engines.ts          # existing CRM engine contracts
  widgets/defaultRegistry.ts     # existing adapters must remain registered
  widgets/presets/crm.ts         # semantic parity reference
  widgets/WidgetRenderer.crm.stories.tsx
  widgets/WidgetRenderer.workshop-crm.stories.tsx
  crm/types.ts                   # existing generic DTOs

examples/xgoja/workshop-crm-site/
  xgoja.v2.yaml
  Makefile
  verbs/workshop-crm.js          # route registration only
  verbs/lib/store.js             # SQLite schema, queries, writes, seeds
  verbs/lib/pages.js             # page builders only
  verbs/lib/fixtures.js          # deterministic seed records
  verbs/lib/calendar.js          # CRM/workshop records -> CalendarEvent DTOs
  assets/public/                 # copied app build
  dist/workshop-crm-site         # rebuilt generated binary
```

Keep the Doodle modularization pattern: routes in the entrypoint, persistence in `store.js`, page composition in `pages.js`, and derived calendar mapping in `calendar.js`.

## 9. Decisions

### Decision: `widget.crm` emits existing IR widgets

- **Context:** CRM rendering components already exist and are registered in React.
- **Options considered:** Create bespoke workshop widgets; expose raw component constructors; add semantic CRM helpers over current IR.
- **Decision:** Add semantic helpers that emit current IR widgets.
- **Rationale:** It gives JavaScript authors a stable domain vocabulary without duplicating React components or bypassing validation.
- **Consequences:** Go mapping logic initially overlaps `widgets/presets/crm.ts`; parity tests must guard against drift.
- **Status:** accepted.

### Decision: opaque builders only for schemas and pipelines

- **Context:** Field schemas and pipelines require order and validation; CRM record data must stay serializable and persistence-owned.
- **Options considered:** All-plain JS objects; opaque objects for everything; builders only for definitions.
- **Decision:** Builders are opaque only for definitions.
- **Rationale:** This matches existing `data.fields` / `data.collection` practice and avoids accidentally making application data runtime-owned.
- **Consequences:** `crm.pipelineBoard` accepts either a builder or built snapshot; data DTOs remain easy to persist and test.
- **Status:** accepted.

### Decision: build an xgoja example before a production backend

- **Context:** The system needs an interaction proof across JS DSL, IR, React, actions, SQLite, and generated binary packaging.
- **Options considered:** Add only unit tests; build production web integration first; build a standalone xgoja example.
- **Decision:** Add `examples/xgoja/workshop-crm-site` as the first full vertical slice.
- **Rationale:** Doodle already demonstrates the build and embedded-app model, while a focused example keeps business persistence isolated.
- **Consequences:** The example must have deterministic data, browser smoke tests, and no imports from app/backend containers.
- **Status:** accepted.

## 10. Implementation phases

### Phase 1: CRM DSL contracts

1. Add `widget.crm` installation and TypeScript declarations.
2. Implement field-schema and pipeline opaque builders.
3. Implement `pipelineBoard`, `recordFields`, `activityFeed`, `tasksInbox`, `stat`, and `funnel`.
4. Add intent helpers and validation errors.
5. Add golden fixtures and Go tests.

### Phase 2: frontend review surface

1. Add a WidgetRenderer workshop CRM Storybook file.
2. Render pipeline, opportunity record, dashboard composition, and calendar composition.
3. Ensure existing components remain API-free and story-covered.

### Phase 3: xgoja example

1. Scaffold `examples/xgoja/workshop-crm-site` from the Doodle structure.
2. Add SQLite schema/migrations and deterministic seed data.
3. Implement read routes first.
4. Add write routes: create lead, move deal, update field, submit availability, schedule run, complete task.
5. Build and embed the current frontend app.

### Phase 4: verification

1. Run Go unit/golden tests.
2. Run package typecheck and focused frontend tests.
3. Run migration checker to ensure example is raw-free.
4. Build the xgoja binary.
5. Browser-smoke the full lead-to-workshop-run journey.

## 11. Test strategy

```text
Layer                    Primary check
-----------------------------------------------------------------
Goja builders            Go unit tests in pkg/widgetdsl
IR shape                 JSON goldens in pkg/widgetdsl/testdata/v3/golden
React adapter            WidgetRenderer Storybook stories
Package compilation      pnpm typecheck + focused tests + app build
Example build            make -C examples/xgoja/workshop-crm-site build
Server behavior          curl/API assertions against /api/widget pages
Browser interaction      Playwright: navigation, action routes, no console errors
Migration hygiene        widgetdsl-migration-checker over the example verbs
```

Minimum tests for each CRM helper:

- success golden
- builder chaining
- callback omitted
- invalid duplicate stage/field error
- action map emitted exactly
- no `widget.raw.component(...)` in the example

## 12. Risks and review points

1. **TS/Go semantic drift:** CRM presets currently live in TypeScript. Review exact IR mapping and add golden parity cases.
2. **Action payload contract:** Drag/drop and field-edit contexts must use names exposed by current adapters: `cardId`, `from`, `to`, `beforeId`, `key`, and `value`.
3. **Calendar interpretation:** Time zones and all-day delivery spans require explicit application conventions. The current `TimeGrid` supports timed blocks; do not force all-day blocks into it.
4. **Generated assets:** New frontend widgets require rebuilding `app-dist`, synchronizing assets, and rebuilding the xgoja binary.
5. **Demo boundaries:** The xgoja example is a reference host, not the production system of record.

## 13. References

- `pkg/widgetdsl/module.go` — v3 root namespace installation.
- `pkg/widgetdsl/v3.go` — Goja v3 builders, time/schedule/UI APIs, and opaque data builders.
- `pkg/widgetdsl/typescript.go` — generated TypeScript declaration surface.
- `packages/rag-evaluation-site/src/widgets/ir/engines.ts` — CRM and calendar IR contracts.
- `packages/rag-evaluation-site/src/widgets/presets/crm.ts` — existing CRM-to-IR semantic reference.
- `packages/rag-evaluation-site/src/widgets/WidgetRenderer.crm.stories.tsx` — current CRM rendering proof.
- `packages/rag-evaluation-site/src/crm/types.ts` — generic CRM data model.
- `examples/xgoja/doodle-site/verbs/` — modular SQLite-backed xgoja host pattern.
