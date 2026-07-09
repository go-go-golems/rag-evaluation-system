---
Title: ""
Ticket: ""
Status: ""
Topics: []
DocType: ""
Intent: ""
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/src/widgets/presets/crm.ts
      Note: Current CRM behavior investigated
    - Path: repo://pkg/widgetdsl/testdata/v3/examples/41-crm-workshop-pipeline.js
      Note: CRM fluent API golden fixture
    - Path: repo://pkg/widgetdsl/testdata/v3/golden/41-crm-workshop-pipeline.json
      Note: Expected CRM Widget IR
    - Path: repo://pkg/widgetdsl/v3.go
      Note: Current builder patterns investigated
    - Path: repo://pkg/widgetdsl/v3_crm.go
      Note: CRM namespace implementation (commit 196cb20800c7d3893daffe6aca37fa9682e0a251)
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
