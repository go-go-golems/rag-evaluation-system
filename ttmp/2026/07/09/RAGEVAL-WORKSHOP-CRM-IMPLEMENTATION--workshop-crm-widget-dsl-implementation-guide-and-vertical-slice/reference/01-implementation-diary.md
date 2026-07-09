---
Title: ""
Ticket: ""
Status: ""
Topics: []
DocType: ""
Intent: ""
Owners: []
RelatedFiles:
    - Path: repo://examples/xgoja/workshop-crm-site/verbs/lib/pages.js
      Note: CRM Widget DSL page composition
    - Path: repo://examples/xgoja/workshop-crm-site/verbs/lib/store.js
      Note: SQLite lead-to-workshop-run persistence
    - Path: repo://examples/xgoja/workshop-crm-site/verbs/workshop-crm.js
      Note: HTTP routes and xgoja host entrypoint (commit 0d81a70b594cfea9a1884d6cfc363c27c2fdb9d2)
    - Path: repo://packages/rag-evaluation-site/src/widgets/presets/crm.ts
      Note: Current CRM behavior investigated
    - Path: repo://pkg/widgetdsl/testdata/v3/examples/41-crm-workshop-pipeline.js
      Note: CRM fluent API golden fixture
    - Path: repo://pkg/widgetdsl/testdata/v3/golden/41-crm-workshop-pipeline.json
      Note: Expected CRM Widget IR
    - Path: repo://pkg/widgetdsl/v3.go
      Note: Current builder patterns investigated
    - Path: repo://pkg/widgetdsl/v3_crm.go
      Note: |-
        CRM namespace implementation (commit 196cb20800c7d3893daffe6aca37fa9682e0a251)
        CRM palette IR fix required by renderer
ExternalSources: []
Summary: ""
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---




# Diary

## Goal

Record the research, documentation, implementation, validation, and review path for the workshop CRM Widget DSL v3 vertical slice.

## Step 1: Establish the intern guide and vertical-slice boundary

I created a dedicated implementation ticket instead of extending the earlier workshop CRM screen-design ticket. The new ticket defines a buildable first slice: lead intake through scheduled workshop run. It is intentionally smaller than a complete business suite, but crosses every required architectural boundary: CRM records, Widget DSL, Widget IR, calendar and scheduling, xgoja HTTP, SQLite, React rendering, and browser validation.

The guide distinguishes proven current code from proposed work. CRM rendering is already implemented as Widget IR adapters and TypeScript presets. The missing implementation is the Go-side `widget.crm` namespace. The guide specifies that this namespace must emit existing IR and use opaque Goja builders only for field-schema and pipeline definitions.

### Prompt Context

**User prompt (verbatim):** "Ok, create a new docmgr ticket, and Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.

Ten implement it step by step, committing at appropriate intervals, keeping a detailed diary."

**Assistant interpretation:** Create a new ticket, document an intern-ready technical design and implementation guide, upload it to reMarkable, then implement the workshop CRM in small committed steps while maintaining a detailed diary.

**Inferred user intent:** Turn the workshop CRM proposal into a durable onboarding document and an auditable, working vertical slice rather than leaving it as an ASCII mockup.

**Commit (code):** N/A — research and documentation step; commit pending.

### What I did
- Created ticket `RAGEVAL-WORKSHOP-CRM-IMPLEMENTATION`.
- Added the primary intern guide and this diary.
- Inspected the existing Widget DSL v3 root installation, opaque `data.fields` builders, CRM DTOs, CRM IR contracts, React CRM presets/stories, time/schedule APIs, and Doodle xgoja host pattern.
- Defined the target as lead intake → pipeline → availability → workshop run rather than attempting billing, authentication, email, asset workflow, and renewal automation simultaneously.

### Why
- An intern needs a stable boundary before implementation. The previous design identified too many product modules for a safe first change.
- Existing CRM renderer support should be reused rather than duplicated through JSX or raw component escape hatches.

### What worked
- Existing CRM interfaces and adapters provide enough rendering capability for the first slice.
- Existing `widget.time`, `widget.schedule`, `widget.ui`, `widget.course`, and `widget.cms` APIs compose with CRM without a separate calendar module.
- `data.fields` establishes a concrete opaque-object pattern for definitions that need Go-side validation.

### What didn't work
- The current `widget.dsl` root has no `widget.crm` namespace. It exports `ui`, `data`, `cms`, `course`, `context`, `schedule`, and `time` only. This is the intended first implementation step.

### What I learned
- The existing TypeScript CRM presets are not an application backend. They are semantic IR compositions that should be treated as the parity reference for Go-side helpers.
- `TimeGrid` is a timed-block engine. Workshop delivery days must be represented as timed blocks until an all-day row is intentionally added.

### What was tricky to build
- The main risk is treating all CRM data as opaque DSL objects. That would entangle application persistence with DSL internals. The guide limits opaque objects to schemas and pipelines, while deals, contacts, tasks, activities, and workshop runs remain ordinary serializable DTOs.

### What warrants a second pair of eyes
- The proposed Go-side helpers overlap existing TypeScript preset semantics. Reviewers should insist on golden parity coverage before accepting a broad API.
- The xgoja example must preserve the Doodle boundary: routes wire requests, `store.js` owns SQLite, pages compose widgets, and calendar mapping derives display data.

### What should be done in the future
- Implement the contracts and demo phases defined in the guide.
- Re-evaluate whether TypeScript presets and Go-side helper mapping can share a generated schema once the behavior is stable.

### Code review instructions
- Begin with the guide design decisions and file layout.
- Compare `packages/rag-evaluation-site/src/widgets/presets/crm.ts` with every future `pkg/widgetdsl/v3.go` CRM helper.
- Validate docs with `docmgr doctor --ticket RAGEVAL-WORKSHOP-CRM-IMPLEMENTATION --stale-after 30`.

### Technical details
- Ticket ID: `RAGEVAL-WORKSHOP-CRM-IMPLEMENTATION`.
- Primary guide: `design-doc/01-intern-guide-workshop-crm-widget-dsl-vertical-slice.md`.
- First code target: add `setExport(exports, "crm", r.v3CRMObject())` in `pkg/widgetdsl/module.go`.

## Step 2: Implement the `widget.crm` v3 namespace

This step added the first executable CRM DSL surface. `widget.crm` is now installed alongside `ui`, `data`, `time`, and `schedule`; its helpers emit existing CRM Widget IR types instead of using raw component escape hatches. The first vertical slice covers pipeline definitions and boards, typed field schemas and field lists, activity feeds, task lists, metrics, funnel segments, and standard server/navigation action contracts.

The API uses opaque Goja builders only where definitions require validation and ordered mutation. `crm.fields(...)` and `crm.pipeline(...)` carry private Go references, expose chainable methods, and produce serializable snapshots through `build()`. Deals, activities, tasks, and record values remain ordinary JavaScript data supplied by the host application.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Implement the planned workshop CRM in committed, validated steps after documenting it.

**Inferred user intent:** Give the workshop CRM a real, idiomatic `widget.dsl` authoring surface before building a host application around it.

**Commit (code):** `196cb20800c7d3893daffe6aca37fa9682e0a251` — "Widget DSL: add CRM namespace"

### What I did
- Added `widget.crm` to `installWidgetV3` in `pkg/widgetdsl/module.go`.
- Added `pkg/widgetdsl/v3_crm.go` with CRM opaque builders, view helpers, and intent helpers.
- Added CRM declarations to `pkg/widgetdsl/typescript.go`.
- Added CRM namespace metadata to `pkg/widgetdsl/v3_descriptors.go` so generated API reference inventory includes it.
- Added `41-crm-workshop-pipeline.js` and its JSON golden fixture.
- Ran Go format, Widget DSL golden tests, package tests, and frontend typecheck.

### Why
- JavaScript authors need semantic CRM APIs that preserve the existing React/IR contracts.
- Pipeline and field definitions have rules that belong in Go-side builders; ordinary CRM data must remain host/persistence-owned JSON.

### What worked
- `go test ./pkg/widgetdsl/... -count=1` passed before commit.
- `go test ./pkg/widgetdsl -run TestWidgetV3GoldenExamplesRenderStableIR -count=1` passed before commit.
- `pnpm --dir packages/rag-evaluation-site typecheck` passed before commit.
- The pre-commit Go test and Go lint steps passed.

### What didn't work
- The pre-commit Biome lint inspected the new JavaScript golden fixture and warned about literal action templates such as `"${dealId}"` with `lint/suspicious/noTemplateCurlyInString`; this is expected DSL action-template syntax but Biome cannot infer it.
- The same hook transiently reported `pkg/widgetdsl/testdata/v3/golden/41-crm-workshop-pipeline.json:1:1 parse ... found the end of the file`, while the committed golden file is populated (`4706` bytes) and the normal Widget DSL golden test passes. This needs follow-up if the fixture lint scope remains enabled.

### What I learned
- The descriptor inventory in `v3_descriptors.go` is the source for exported namespace declarations and API-reference generation; adding runtime exports alone is incomplete.
- Existing CRM adapter action contexts use `cardId`, `from`, `to`, `beforeId`, `key`, and `value`; CRM intent helpers must not invent incompatible names.

### What was tricky to build
- Builder values must be accepted by later helpers without exposing Go internals. `crmPipelineFromValue` and `crmFieldsFromValue` first resolve the private builder reference, then fall back to a serializable built object. This lets callers pass either a fluent builder or a stored definition snapshot.
- Grouped CRM fields require conversion from a flat field schema to `RecordFieldList.sections`. The helper preserves first-seen group order and defaults ungrouped fields to `Details`.

### What warrants a second pair of eyes
- `v3CRMBoardProps` must stay semantically aligned with `packages/rag-evaluation-site/src/widgets/presets/crm.ts` as both construct `BoardEngine` props.
- The current `tasksInbox` helper intentionally provides a lightweight static panel; task completion action wiring should be added with a dedicated builder/test before relying on it for write workflows.

### What should be done in the future
- Add direct unit tests for invalid empty/duplicate field and stage definitions, in addition to the end-to-end golden fixture.
- Build the SQLite-backed `workshop-crm-site` host and validate actions in a browser.

### Code review instructions
- Start at `pkg/widgetdsl/v3_crm.go`, especially `v3CRMFieldsBuilder`, `v3CRMPipelineBuilder`, and `v3CRMPipelineBoard`.
- Compare generated `BoardEngine` properties against `packages/rag-evaluation-site/src/widgets/presets/crm.ts`.
- Run `go test ./pkg/widgetdsl/... -count=1` and `pnpm --dir packages/rag-evaluation-site typecheck`.

### Technical details
- New fixture: `pkg/widgetdsl/testdata/v3/examples/41-crm-workshop-pipeline.js`.
- New golden: `pkg/widgetdsl/testdata/v3/golden/41-crm-workshop-pipeline.json`.
- New runtime file: `pkg/widgetdsl/v3_crm.go`.

## Step 3: Build the SQLite-backed workshop CRM reference host

This step turns the CRM DSL into a runnable vertical slice at `examples/xgoja/workshop-crm-site`. The host uses the Doodle xgoja pattern—one route entrypoint, a SQLite store, and widget-composition pages—but models the workshop lifecycle instead: capture a lead, inspect it in the pipeline and record view, choose a delivery slot, and create a persisted workshop run.

While browser-testing the dashboard, the new funnel exposed a renderer-contract gap: `SegmentedBar` requires a `styleSet`, so rendering a `crm.funnel(...)` without one crashed in React. I fixed the CRM DSL to serialize the existing CRM stage and activity palettes as part of its generated IR, then rebuilt the fixture and the reference binary before retrying the interaction flow.

### Prompt Context

**User prompt (verbatim):** "go ahead, implement step by step, commit at appropriate intervals, keep a detailed diary as you work. (Read skills if necessary)"

**Assistant interpretation:** Continue the workshop CRM plan in small validated commits and record chronological implementation evidence.

**Inferred user intent:** Obtain a demonstrably working, maintainable CRM vertical slice rather than only a DSL API or design document.

**Commit (code):** `0d81a70b594cfea9a1884d6cfc363c27c2fdb9d2` — "Examples: add workshop CRM vertical slice"

### What I did
- Added `examples/xgoja/workshop-crm-site/`, including its xgoja v2 manifest, Makefile, embedded SPA assets, and generated runnable binary.
- Added `verbs/lib/store.js` with SQLite tables for organizations, contacts, deals, activities, availability options, and workshop runs plus deterministic seed data.
- Added `verbs/lib/pages.js` for command center, pipeline, lead intake, opportunity record/activity, availability, and workshop-run pages.
- Added form POST routes for lead creation and scheduling a selected availability option in `verbs/workshop-crm.js`.
- Added default CRM stage/activity `styleSet` payloads in `pkg/widgetdsl/v3_crm.go`, then refreshed the CRM golden fixture.
- Ran the raw-free migration checker, Go test suite, frontend typecheck, xgoja build, HTTP assertions, and a Playwright lead-to-run journey.

### Why
- The host proves the intended application seam: fluent DSL definitions and IR composition in JS, plain serializable data in SQLite, and existing React WidgetRenderer components in the browser.
- Palette data is a renderer requirement, not an optional visual detail for `SegmentedBar`; emitting it from `widget.crm` prevents all hosts from rediscovering the same runtime failure.

### What worked
- `go run ./cmd/widgetdsl-migration-checker -- examples/xgoja/workshop-crm-site/verbs examples/xgoja/workshop-crm-site/xgoja.v2.yaml` reported: `No legacy Widget DSL imports or raw component escape hatches found.`
- `go test ./pkg/widgetdsl/... -count=1` passed.
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `make -C examples/xgoja/workshop-crm-site sync-app` and `make -C examples/xgoja/workshop-crm-site build` produced `dist/workshop-crm-site`.
- Playwright successfully created `Orbit Analytics`, navigated to its opportunity, opened availability, submitted a delivery date, and reached `/pages/runs` with no browser console errors.
- The final `TimeGrid` API response covered `2026-09-07` through `2026-09-13` for the selected September run.

### What didn't work
- Initial xgoja build failed with `Error: parse lib/pages.js: syntax errors while collecting imports`; Biome pinpointed an invalid trailing comma/parenthesis in the fluent `opportunityPage` chain. I rewrote that function with explicit intermediate values and properly nested chained calls, after which `make ... build` passed.
- Initial browser navigation to `/pages/index` threw `TypeError: Cannot read properties of undefined (reading 'styles')`. The cause was `crm.funnel(...)` emitting `SegmentedBar` props without its required `styleSet`. Adding CRM palette serialization fixed the crash.
- A first workshop-run calendar render showed the current July week even for a September run. `widget.time.range.week` expects a date-like value; passing `startISO.slice(0, 10)` instead of the local timestamp made the generated `days` use the September week.
- The staged generated SPA assets cause Biome to report thousands of diagnostics against minified CSS/JS. This is pre-existing generated-asset lint behavior (the Doodle host also commits generated assets); hooks still ran the Go tests and linters successfully. Source-specific warnings remain for the intentional `"${dealId}"` DSL template and externally invoked `site()` verb.

### What I learned
- The TypeScript `stageStyleSet` and `activityStyleSet` contracts are necessary semantic parity data for Go CRM helpers, not merely Storybook styling.
- The xgoja SQLite file is relative to the binary process working directory. Running the binary from the repository root creates `workshop-crm.db` there; `make serve` runs from the example directory and keeps it beside the example. Operators should prefer `make -C examples/xgoja/workshop-crm-site serve`.

### What was tricky to build
- The host must not persist opaque Goja builders. `pages.js` keeps pipeline and field definitions in-memory as DSL builder values, while `store.js` returns only arrays/objects from SQLite. This preserves the intended definition/data separation.
- Calendar range selection initially appeared correct at the event level but wrong in the rendered week header because local timestamp parsing did not satisfy the range parser. Reducing the timestamp to ISO calendar date at the caller was the smallest explicit fix.

### What warrants a second pair of eyes
- Review whether copying the CRM palette maps into Go is acceptable short-term or whether the frontend/Go contract should share a generated declarative source before more palette variants are added.
- Review form write handling before treating this example as security guidance: it intentionally omits authentication, authorization, validation beyond required fields, and CSRF protections.
- The reference host offers explicit forms for durable writes; the `crm.intent.moveDeal` is emitted for board parity but does not yet have a generic xgoja action dispatcher.

### What should be done in the future
- Add committed Playwright smoke scripts for this host and the existing go-go-course host.
- Add unit coverage for duplicate/empty CRM builder definitions and a test that `crm.funnel` always serializes a usable style set.
- Add intentional action-dispatch plumbing only after agreeing on a server-action protocol; do not make BoardEngine drag/drop mutate SQLite implicitly.

### Code review instructions
- Start with `examples/xgoja/workshop-crm-site/verbs/workshop-crm.js` to inspect route and SPA boundaries.
- Review `verbs/lib/store.js` for persistence invariants and `verbs/lib/pages.js` for the lead-to-run composition.
- Review `pkg/widgetdsl/v3_crm.go` palette functions against `packages/rag-evaluation-site/src/crm/palettes.ts`.
- Validate with `make -C examples/xgoja/workshop-crm-site build`, then run `make -C examples/xgoja/workshop-crm-site serve` and visit `http://127.0.0.1:18794/pages/index`.

### Technical details
- Runtime server used for smoke validation: `127.0.0.1:18794`.
- Key routes: `/pages/index`, `/pages/pipeline`, `/pages/lead`, `/pages/opportunity?deal=<id>`, `/pages/availability?deal=<id>`, `/pages/runs`.
- Write routes: `POST /api/form/create-lead` and `POST /api/form/schedule-run?deal=<id>`.
