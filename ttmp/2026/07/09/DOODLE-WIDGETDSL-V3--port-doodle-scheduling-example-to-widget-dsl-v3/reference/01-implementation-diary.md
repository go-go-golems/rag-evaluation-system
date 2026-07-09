---
Title: Implementation Diary
Ticket: DOODLE-WIDGETDSL-V3
Status: active
Topics:
    - ui-dsl
    - xgoja
    - sqlite
    - web
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://examples/xgoja/doodle-site/verbs/doodle.js
      Note: |-
        ported Doodle pages from legacy split modules to widget.dsl v3
        uses typed ui and schedule DSL helpers without raw escape hatches
        Added calendar visualization over Doodle poll slots and votes
        Small jsverb entrypoint after module split
    - Path: repo://examples/xgoja/doodle-site/verbs/lib/calendar.js
      Note: Calendar marker/detail/selection module
    - Path: repo://examples/xgoja/doodle-site/verbs/lib/pages.js
      Note: Doodle Widget DSL page composition module
    - Path: repo://examples/xgoja/doodle-site/verbs/lib/store.js
      Note: SQLite schema, seed, query, and mutation module
    - Path: repo://examples/xgoja/doodle-site/xgoja.v2.yaml
      Note: selected widget.dsl instead of ui.dsl/data.dsl
    - Path: repo://pkg/widgetdsl/typescript.go
      Note: declared new v3 UI wrappers
    - Path: repo://pkg/widgetdsl/v3.go
      Note: |-
        added typed v3 UI wrappers used by Doodle
        Typed splitPane and TimeGrid viewport DSL helpers
    - Path: repo://ttmp/2026/07/09/DOODLE-WIDGETDSL-V3--port-doodle-scheduling-example-to-widget-dsl-v3/tasks.md
      Note: completed Doodle v3 port checklist
ExternalSources: []
Summary: Chronological diary for porting the Doodle scheduling example to widget.dsl v3.
LastUpdated: 2026-07-09T00:00:00Z
WhatFor: Record implementation steps, validation, failures, and review instructions for the Doodle v3 port.
WhenToUse: Read before editing examples/xgoja/doodle-site.
---





# Diary

## Goal

Capture the future Doodle example migration from legacy split widget modules to `widget.dsl` v3.

## Step 1: Ticket setup

Created a separate ticket for the Doodle v3 port so the product-shaped example migration does not blur with the immediate hardening work for scheduling widgets and `go-go-course`.

### Prompt Context

**User prompt (verbatim):** (see `RAGEVAL-WIDGETDSL-HARDENING` diary Step 1)

**Assistant interpretation:** Create a dedicated ticket for item 4: port the Doodle scheduling example to `widget.dsl` v3.

**Inferred user intent:** Keep the Doodle demo migration independently trackable and resumable.

**Commit (code):** N/A — setup only.

### What I did
- Created `DOODLE-WIDGETDSL-V3`.
- Added a port plan, task list, and this diary.

### Why
- `DOODLE-1` is a working legacy-module demo; the v3 port is valuable but should follow the hardening ticket.

### What worked
- Ticket setup completed cleanly.

### What didn't work
- N/A.

### What I learned
- The Doodle port should reuse the already verified SQLite/xgoja flow and change the widget authoring layer, not redesign the whole app.

### What was tricky to build
- N/A for setup.

### What warrants a second pair of eyes
- Whether the v3 port should keep native forms or introduce v3 action-driven submit flows.

### What should be done in the future
- Port `examples/xgoja/doodle-site` after the hardening ticket's first slice is stable.

### Code review instructions
- Start with `design-doc/01-doodle-widget-dsl-v3-port-plan.md`.

### Technical details
- Source baseline: `DOODLE-1` and `examples/xgoja/doodle-site`.

## Step 2: Port Doodle from split modules to `widget.dsl` v3

Ported the Doodle xgoja example so it now selects and imports only `widget.dsl` from the `rag-widget-site` provider. The SQLite schema, planned-route Express handlers, native form POST flow, and browser behavior remain the same; the page authoring layer now uses v3 page/section/action/data helpers with small raw escape hatches for form controls that do not yet have typed v3 wrappers.

The browser validation covered the real flow rather than only API JSON: open the create page, create a new poll, arrive on the new poll page, submit availability, and verify that the availability grid and tallies update without console errors.

### Prompt Context

**User prompt (verbatim):** "Update the doodle"

**Assistant interpretation:** Work the Doodle v3 ticket by updating the existing Doodle scheduling example to use `widget.dsl` v3.

**Inferred user intent:** Bring the product-shaped Doodle demo in line with the new widget DSL module policy after the hardening ticket setup.

**Commit (code):** N/A — changes not committed yet.

### What I did
- Updated `examples/xgoja/doodle-site/xgoja.v2.yaml` so `rag-widget-site` selects only `widget.dsl`.
- Rewrote `examples/xgoja/doodle-site/verbs/doodle.js` to import `const widget = require("widget.dsl")` instead of `ui.dsl` and `data.dsl`.
- Rebuilt pages with v3 constructs:
  - `widget.page(...)` with page metadata via builder calls;
  - `widget.act.navigate(...)` for navigation;
  - `widget.data.fields(...)` + `widget.data.collection(...).table().toNode()` for index, availability, and result tables;
  - `widget.ui.form(...)`, `widget.ui.button(...)`, `widget.ui.inline(...)`, section metrics, section metadata, and captions.
- Kept native form POST handlers and SQLite persistence unchanged.
- Used `widget.raw.component(...)` for form rows/inputs, status text, and empty state because those controls still lack dedicated v3 typed wrappers.
- Checked all Doodle v3 tasks.

### Why
- `DOODLE-1` proved the Doodle app as a legacy split-module example. This ticket's job is to prove the same app can run with the v3 module selection.
- Keeping the database and HTTP flow unchanged isolates the migration to the widget authoring layer.

### What worked
- `make build` succeeded for `examples/xgoja/doodle-site`.
- `make list-modules` shows only `rag-widget-site | widget.dsl | widget.dsl` for widget modules.
- API smoke returned root `Stack` pages for `index`, `create`, and `poll`.
- Browser flow succeeded:
  - opened `/pages/create`;
  - filled and submitted a new poll;
  - redirected to `/pages/poll?poll=2`;
  - submitted availability as `Rosalind`;
  - saw the availability grid and results table update.
- Console check after the new-server run reported no errors/warnings.
- Response JSON did not contain `[object Object]`.

### What didn't work
- The first browser attempt hit an old Doodle server still listening on `127.0.0.1:18793` from the earlier legacy run. The new server failed with `bind: address already in use`; I killed the old process and restarted the freshly built binary before validating.
- The migration checker reports raw-component escape hatches for `StatusText`, `EmptyState`, `FormRow`, `TextInput`, `TextareaInput`, and `SelectInput`. These are expected transitional findings, not legacy split-module imports.

### What I learned
- The Doodle app can be migrated cleanly without changing persistence or route semantics; the v3 port is mainly page-construction syntax plus table schemas.
- The v3 `data.collection` path is a good replacement for legacy `data.dataTable` when the table schema can be declared near the page.
- A real browser flow is still necessary because native form submit behavior is not proven by API page JSON alone.

### What was tricky to build
- The old example used `ui.formRow`, `ui.textInput`, `ui.textareaInput`, and `ui.selectInput`; v3 has `widget.ui.form` but not yet typed helpers for every form atom. I isolated those remaining gaps in tiny raw helper functions so they are easy to replace later.
- The availability grid has dynamic columns from database rows. I built its schema dynamically with `widget.data.fields(...)`, adding one field per option, then used a v3 collection table.

### What warrants a second pair of eyes
- Whether v3 should grow typed `formRow`, `textInput`, `textareaInput`, `selectInput`, `statusText`, and `emptyState` helpers so examples like Doodle can become raw-free.
- Whether the Doodle availability grid should switch to a typed `schedule.availabilityPoll` view after slot labels become structured date/time DTOs instead of free-form labels.

### What should be done in the future
- Add typed v3 wrappers for the remaining raw form/status/empty-state helpers if they recur in examples.
- Optionally close `DOODLE-WIDGETDSL-V3` after review, since the checklist is complete.

### Code review instructions
- Start with `examples/xgoja/doodle-site/xgoja.v2.yaml` to verify module selection.
- Review `examples/xgoja/doodle-site/verbs/doodle.js` helpers first, then `indexPage`, `createPage`, and `pollPage`.
- Validate with:
  - `go run ./cmd/widgetdsl-migration-checker -- examples/xgoja/doodle-site/verbs`
  - `cd examples/xgoja/doodle-site && make list-modules`
  - `cd examples/xgoja/doodle-site && make build`
  - run `./dist/doodle-site serve doodle site --http-listen 127.0.0.1:18793` and smoke create/vote in the browser.

### Technical details
- Post-port checker result: zero `legacy-module-import` findings; six expected `raw-component-escape-hatch` findings in local raw helper functions.
- Browser smoke used a clean `doodle.db` created after restarting the fresh binary.

## Step 3: Replace raw Doodle helpers with typed v3 UI and schedule DSLs

Tightened the Doodle v3 port so it no longer relies on `widget.raw.component(...)`. I added the missing typed UI wrappers to `widget.dsl` itself, then rewrote the Doodle page code to use those wrappers plus the schedule DSL's availability matrix and poll-summary components.

The Doodle example now demonstrates the intended new surface more clearly: `widget.ui.formRow`, `widget.ui.textInput`, `widget.ui.textareaInput`, `widget.ui.selectInput`, `widget.ui.status`, `widget.ui.emptyState`, `widget.schedule.availabilityPoll`, and `widget.schedule.pollSummary`. The migration checker now reports zero legacy imports and zero raw escape hatches for the Doodle source.

### Prompt Context

**User prompt (verbatim):** "use the new DSLs and components."

**Assistant interpretation:** Remove the Doodle example's remaining raw-component compatibility helpers and use the typed v3 DSL/component helpers instead.

**Inferred user intent:** Make Doodle a clean, idiomatic v3 demo rather than a host that merely selects `widget.dsl` while still hand-emitting raw components.

**Commit (code):** N/A — changes will amend the Doodle v3 commit.

### What I did
- Added typed v3 UI helpers in `pkg/widgetdsl/v3.go`:
  - `ui.formRow(...)`
  - `ui.textInput(...)`
  - `ui.textareaInput(...)`
  - `ui.selectInput(...)`
  - `ui.status(...)`
  - `ui.emptyState(...)`
- Updated widget.dsl TypeScript declarations for the new UI helpers.
- Replaced Doodle raw helper functions with typed `widget.ui.*` calls.
- Replaced the hand-built availability table with `widget.schedule.availabilityPoll(...)`.
- Added a schedule summary grid via `widget.schedule.pollSummary(...)` and kept a compact result table for score/best-slot information.
- Rebuilt the Doodle binary and reran browser validation.

### Why
- Raw escape hatches are useful for migration but are not the desired final demo surface.
- Doodle is now a better example for new hosts because it uses v3 UI, data, schedule, and action namespaces directly.

### What worked
- `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1` passed.
- `cd examples/xgoja/doodle-site && make build` passed.
- Migration checker now reports: `No legacy Widget DSL imports or raw component escape hatches found.`
- Browser smoke passed again: create poll, submit availability, availability matrix updates, summary grid updates, score/best table updates, and no new console errors/warnings.

### What didn't work
- My first replacement still left Doodle using schedule-neutral data tables for the availability grid. I changed that to `widget.schedule.availabilityPoll(...)` and converted database votes (`yes`/`maybe`/`no`) into schedule states (`available`/`maybe`/`unavailable`) at the view boundary.

### What I learned
- `schedule.availabilityPoll` is usable even when persisted data uses product-specific values, as long as the page layer normalizes them into the schedule view contract.
- The missing typed form/status wrappers were the real reason Doodle needed raw calls; adding those helpers made the example simpler and cleaned up the checker output.

### What was tricky to build
- Doodle stores free-form slot labels and vote values rather than structured scheduling DTOs. The solution was to keep persistence unchanged and normalize only the rendered `availabilityPoll` object: option IDs become strings, vote values become `available`/`maybe`/`unavailable`, and the native form can keep posting `yes`/`maybe`/`no`.

### What warrants a second pair of eyes
- Whether the UI helper signatures are the right long-term shape, especially `ui.status(status, value, options?)` and `ui.emptyState(title, description?, options?)`.
- Whether `schedule.availabilityPoll` should support custom state labels/glyphs so Doodle can display `yes/no` while keeping the generic schedule contract.

### What should be done in the future
- Add runnable v3 example fixtures for the new UI helpers if they are not already covered by Doodle.
- Consider promoting the score/best-slot result table into a schedule-specific result view.

### Code review instructions
- Start with the new `v3UIObject` exports and helper functions in `pkg/widgetdsl/v3.go`.
- Review the Doodle `pollPage` normalization to `availabilityPoll` and `summaryTallies`.
- Validate with:
  - `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1`
  - `go run ./cmd/widgetdsl-migration-checker -- examples/xgoja/doodle-site/verbs examples/xgoja/doodle-site/xgoja.v2.yaml`
  - `cd examples/xgoja/doodle-site && make build`
  - browser smoke on `127.0.0.1:18793`.

### Technical details
- Checker result after this step: zero findings.
- Doodle uses native form POST for writes; schedule widgets are read-only render views in this example.

## Step 4: Add calendar widgets to the Doodle poll page

Extended the Doodle poll page so it demonstrates the calendar/time side of `widget.dsl` v3 in addition to the schedule availability matrix. The page now derives calendar events from offered poll slots and renders them through `widget.time.month(...)` and `widget.time.week(...)` when slot labels can be parsed as dates and times.

The calendar blocks include the participant names grouped by vote value, so the week view shows both the offered time and who answered yes, maybe, or no. This keeps SQLite storage unchanged while adding a second visualization over the same poll data.

### Prompt Context

**User prompt (verbatim):** "ok, let's visualize the offered times and who signaled availability using the calendar widgets as well, to demonstrate their use"

**Assistant interpretation:** Add a Doodle poll-page section that uses calendar widgets to show parseable offered slots and participant availability.

**Inferred user intent:** Make Doodle demonstrate both schedule availability widgets and calendar/time widgets in one real app.

**Commit (code):** N/A — pending commit.

### What I did
- Added slot-label parsing for labels like `Thu Jul 9 · 19:00` and `2026-07-09 19:00`.
- Added `calendarEventsForPoll(...)` to convert options, participants, votes, and tallies into calendar event DTOs.
- Added a `Calendar view` section to the poll page.
- Rendered parseable slots through `widget.time.month(...)` and `widget.time.week(...)`.
- Kept the existing `widget.schedule.availabilityPoll(...)`, `widget.schedule.pollSummary(...)`, and result table.
- Rebuilt the Doodle binary.

### Why
- Doodle already demonstrated availability-oriented schedule widgets.
- The app did not yet demonstrate the time/calendar widgets on the same real data.
- Calendar rendering is useful only if it shows participant availability, not only empty offered slots.

### What worked
- The migration checker still reports no legacy imports and no raw escape hatches.
- `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1` passed.
- `cd examples/xgoja/doodle-site && make build` passed.
- API validation confirmed the poll page now emits both `MonthGrid` and `TimeGrid` nodes.
- The generated `TimeGrid` blocks include participant names in labels and vote groups in `meta`.
- Browser smoke loaded `/pages/poll?poll=1` with no console warnings or errors.

### What didn't work
- N/A.

### What I learned
- The current Doodle slot model is intentionally free-form, so calendar rendering needs a parseable-label adapter at the page boundary.
- `widget.time.week(...)` can be demonstrated without changing persistence if the page layer derives `startISO`, `endISO`, `dayISO`, and `label` fields from existing options and votes.

### What was tricky to build
- Doodle does not store structured start/end timestamps. The solution was to keep storage unchanged and add conservative parsing for the slot formats already shown in the create-page placeholder. Unparseable labels fall back to an empty-state message rather than producing incorrect calendar events.
- The week calendar needs a useful hour range. I derive a small range from the parsed slot hours so evening polls do not appear in a large mostly empty 8–18 grid.

### What warrants a second pair of eyes
- Whether free-form slot labels should remain acceptable for the demo, or whether Doodle should eventually store structured timestamps.
- Whether the calendar event labels are too long when many participants respond.

### What should be done in the future
- Consider storing optional `start_iso` / `end_iso` columns if Doodle becomes more than a demo.
- Consider adding a dedicated compact event detail renderer if participant lists become large.

### Code review instructions
- Start in `examples/xgoja/doodle-site/verbs/doodle.js` at `parseSlotLabel`, `calendarEventsForPoll`, and the `Calendar view` section inside `pollPage`.
- Validate with:
  - `go run ./cmd/widgetdsl-migration-checker -- examples/xgoja/doodle-site/verbs examples/xgoja/doodle-site/xgoja.v2.yaml`
  - `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1`
  - `cd examples/xgoja/doodle-site && make build`
  - load `/pages/poll?poll=1` and confirm the Calendar view renders without console errors.

### Technical details
- Calendar event style keys are chosen from availability counts: `focus` when at least one yes exists, `personal` when only maybe exists, and `meeting` otherwise.
- Slot labels are parsed with the poll creation year as fallback for month/day labels.

## Step 5: Refactor Doodle modules and tighten calendar UX

The calendar visualization exposed two issues in the first implementation. The month calendar was too cramped because full participant strings were passed as marker labels, and the page source had grown too large for easy review. I changed the month data to compact markers only, moved details into a right-hand column, added day/slot selection through URL-backed navigation, and split the Doodle jsverb into smaller CommonJS modules.

The split keeps the jsverb entrypoint small: `doodle.js` now wires modules and HTTP routes, `store.js` owns SQLite access, `pages.js` owns page composition, `calendar.js` owns slot parsing and calendar view models, and `widget-helpers.js` owns local v3 helper aliases.

### Prompt Context

**User prompt (verbatim):** "the view for /tmp/pi-clipboard-682eb06d-a7d2-41e7-822b-4f598bc1ba1b.png is way too cramped, and shuld in fact not allow passing strings through the API but only \"markers\", is that makes sense. Look at the image with the read tool: /tmp/pi-clipboard-07438e5a-06cd-499c-9494-cc8f2ecaf9c9.png

What we could do is maybe, when you select the day, then on a right hand column it shows the details of the day.

Can you do side by side horizontal columns with you widget.dsl primitives?"

**Assistant interpretation:** Fix the cramped calendar month view by using marker-only month cells and show selected-day details beside the month calendar, using Widget DSL layout primitives.

**Inferred user intent:** Keep the calendar demo readable and prove that `widget.dsl` can express a two-column detail layout.

**Commit (code):** N/A — pending commit.

### What I did
- Added `widget.ui.splitPane(left, right, options)` to `widget.dsl` v3.
- Added `viewportHeight(height)` to the `widget.time.week(...)` builder and passed style through the TimeGrid widget adapter so the week view can scroll inside a bounded viewport.
- Changed Doodle month rendering to pass `markers` only, not event strings, into `widget.time.month(...)`.
- Added selected-day state via `?day=YYYY-MM-DD` and selected-slot state via `?slot=slot-id`.
- Added a right-hand details pane beside the month view with participant names grouped by yes/maybe/no.
- Made TimeGrid block clicks navigate to `?slot=...`, which selects the corresponding week block and focuses the month/detail view on that block's day.
- Split `verbs/doodle.js` into:
  - `verbs/doodle.js`
  - `verbs/lib/store.js`
  - `verbs/lib/pages.js`
  - `verbs/lib/calendar.js`
  - `verbs/lib/widget-helpers.js`

### Why
- Month cells should stay compact. The month widget should receive marker data, while detailed text belongs outside the month grid.
- The Doodle source had become a mixed file containing schema, routes, page composition, calendar parsing, and widget helper aliases. Splitting by responsibility makes future edits safer.
- The week view needed a bounded viewport so users can scroll across the full 8:00-22:00 range.

### What worked
- `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1` passed.
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- Migration checker still reports no legacy imports and no raw escape hatches.
- `cd examples/xgoja/doodle-site && make build` passed.
- API smoke for `/api/widget/pages/poll?poll=1&slot=slot-3` confirmed `SplitPane`, selected `MonthGrid`, selected `TimeGrid`, and bounded TimeGrid style output.
- Browser smoke loaded `/pages/poll?poll=1&slot=slot-3` with no console warnings or errors.

### What didn't work
- TypeScript declaration editing initially missed commas in generated declaration string slices. `go test` caught the Go syntax error, and I fixed the missing commas before continuing.

### What I learned
- jsverbs supports relative CommonJS `require(...)`, so Doodle can be split into local modules without changing the xgoja source spec.
- `SplitPane` was already registered in the React widget registry, but v3 lacked a typed helper. Adding `widget.ui.splitPane` was the missing DSL surface.
- TimeGrid accepted a `style` prop at the IR level but its widget adapter was not passing it through. Passing `style` through made a DSL-level viewport helper possible.

### What was tricky to build
- The month and week widgets have different selection contexts. Month clicks provide `dateISO`, while TimeGrid block clicks provide `blockId`. I kept the page stateless by encoding both selections in URL query parameters and deriving the selected day from the selected block when present.
- The month widget should not receive long participant strings as marker labels. The solution was to emit only `{ count, styleKey }` markers and render participant details in the side pane.

### What warrants a second pair of eyes
- Whether `widget.time.week(...).viewportHeight(...)` should set `maxHeight` directly or expose a more general scroll-region/layout helper.
- Whether `widget.ui.splitPane` should accept a richer ratio set for application pages.
- Whether Doodle should store structured slot timestamps instead of parsing labels once the demo becomes more than an example.

### What should be done in the future
- Add a small fixture/golden for `widget.ui.splitPane`.
- Add a small fixture/golden for `widget.time.week(...).viewportHeight(...)`.
- Consider promoting calendar slot parsing into a reusable helper only if another app needs the same free-form label behavior.

### Code review instructions
- Start at `examples/xgoja/doodle-site/verbs/doodle.js` to see the new entrypoint.
- Review `examples/xgoja/doodle-site/verbs/lib/store.js` for persistence boundaries.
- Review `examples/xgoja/doodle-site/verbs/lib/pages.js` for page composition and schedule view models.
- Review `examples/xgoja/doodle-site/verbs/lib/calendar.js` for marker-only month data and selection behavior.
- Review `pkg/widgetdsl/v3.go` for `ui.splitPane` and `time.week().viewportHeight`.
- Validate with:
  - `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1`
  - `pnpm --dir packages/rag-evaluation-site typecheck`
  - `go run ./cmd/widgetdsl-migration-checker -- examples/xgoja/doodle-site/verbs examples/xgoja/doodle-site/xgoja.v2.yaml`
  - `cd examples/xgoja/doodle-site && make build`
  - load `/pages/poll?poll=1&slot=slot-3`.

### Technical details
- Month selection uses `act.navigate('/pages/poll?poll=<id>&day=$dateISO')`.
- Week block selection uses `act.navigate('/pages/poll?poll=<id>&slot=$blockId')`.
- If `slot` is present, it wins over `day` and determines the selected day.
- Month markers are now compact objects: `{ count, styleKey }`.
