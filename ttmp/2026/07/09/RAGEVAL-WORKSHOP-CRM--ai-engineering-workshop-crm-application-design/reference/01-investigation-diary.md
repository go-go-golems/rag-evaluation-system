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
      Note: Investigation source for CRM helpers
    - Path: repo://pkg/widgetdsl/v3.go
      Note: Investigation source for time/schedule/ui helper availability
ExternalSources: []
Summary: ""
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---


# Diary

## Goal

Capture the investigation and design pass for an AI engineering workshop CRM built with the CRM, calendar/time, scheduling, and generic Widget DSL surfaces.

## Step 1: Map existing widget capabilities and design workshop CRM screens

I created a new docmgr ticket for the workshop CRM design and inspected the existing CRM, calendar/time, scheduling, and generic UI widgets. The result is a screen map for a workshop business CRM rather than a generic sales CRM: it connects inbound leads, opportunity stages, contact records, workshop delivery dates, facilitator schedules, materials, tasks, and post-workshop follow-up.

The main technical finding is that CRM rendering support already exists in the React/Widget IR layer, especially through TypeScript presets and reusable engines. The Go `widget.dsl` v3 root namespace does not yet expose a first-class `widget.crm` object, so a production implementation should add thin Go-side helpers that emit the existing IR nodes.

### Prompt Context

**User prompt (verbatim):** "ok, can you create with the CRM part of the DSL, along with the calendar part (and of course the other ones), a full CRM system for our workshops, we are giving these AI engineering workshops with a friend.

Create a new docmgr ticket, and then after investigating the widgets and all that, make a set of screens for the application and output them here in the chat as ASCII screenshots + skektch of JS"

**Assistant interpretation:** Create a ticket, inspect current widget/DSL capabilities, then design a workshop-focused CRM app using CRM, calendar/time, scheduling, UI, data, course, and CMS widgets; provide screen sketches and JavaScript DSL sketches.

**Inferred user intent:** Quickly turn the existing widget platform into an application concept for running an AI engineering workshop business with a co-facilitator.

**Commit (code):** N/A — documentation/design only.

### What I did
- Created ticket `RAGEVAL-WORKSHOP-CRM`.
- Added design doc `design-doc/01-workshop-crm-screens-and-widget-dsl-plan.md`.
- Added this investigation diary.
- Inspected CRM domain DTOs, CRM presets, CRM Storybook coverage, Widget IR engine props, calendar/time DSL declarations, scheduling DSL declarations, and generic UI helpers.
- Added and checked three ticket tasks covering capability confirmation, screen design, and JS sketches.

### Why
- The workshop CRM should reuse generic engines instead of inventing workshop-specific widgets too early.
- The current widget ecosystem already has most of the required primitives: pipeline, fields, activities, stat tiles, month/week calendar, availability/booking, forms, tables, course shells, media libraries, panels, split panes, and share links.

### What worked
- `docmgr ticket create-ticket --ticket RAGEVAL-WORKSHOP-CRM --title "AI engineering workshop CRM application design" --topics widget-dsl,crm,calendar,workshops,design` created the workspace.
- `docmgr doc add` created the design doc and diary.
- Repository inspection found CRM presets in `packages/rag-evaluation-site/src/widgets/presets/crm.ts` and CRM stories in `packages/rag-evaluation-site/src/widgets/WidgetRenderer.crm.stories.tsx`.
- Time and scheduling helpers are exposed in `pkg/widgetdsl/v3.go` and `pkg/widgetdsl/typescript.go`.

### What didn't work
- There is no current `widget.crm` namespace in Go `widget.dsl` v3. The CRM capabilities are currently present as React/TypeScript presets and IR components, so the JS sketch uses a proposed namespace that should be implemented as a thin wrapper.

### What I learned
- The CRM field system is data-driven through `FieldDef`, `FieldSpec`, and `RecordFieldList`; workshop-specific fields should be schema data.
- Calendar/time support is already good enough for workshop delivery and prep scheduling, including bounded week views.
- Scheduling widgets can cover client availability polling and booking without a separate calendar module.

### What was tricky to build
- The phrase “CRM part of the DSL” maps to existing CRM presets and IR widgets, but not yet to Go-side `widget.dsl` helpers. I treated this as an implementation gap and designed the sketches around a future `widget.crm` namespace instead of relying on raw component escape hatches.
- A workshop CRM crosses several domains: sales, delivery, course materials, scheduling, tasks, and post-workshop success. The screen plan keeps the domain model small and lets record fields capture the custom details.

### What warrants a second pair of eyes
- Whether `widget.crm` should be implemented directly in Go `pkg/widgetdsl/v3.go` or generated from the TypeScript preset layer.
- Whether `WorkshopRun` deserves first-class widget helpers or should remain a record type rendered by CRM fields plus calendar/course widgets.

### What should be done in the future
- Implement `widget.crm.*` helpers in Go DSL v3.
- Add golden fixtures for `widget.crm.pipelineBoard`, `widget.crm.record`, `widget.crm.dashboard`, and `widget.crm.tasksInbox`.
- Build a SQLite-backed xgoja demo app for the workshop CRM.

### Code review instructions
- Start with `packages/rag-evaluation-site/src/widgets/presets/crm.ts` for existing CRM composition patterns.
- Review `packages/rag-evaluation-site/src/widgets/ir/engines.ts` for the serializable CRM contracts.
- Review `pkg/widgetdsl/v3.go` and `pkg/widgetdsl/typescript.go` before adding `widget.crm` helpers.
- Validate future changes with `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1` and frontend typecheck/Storybook.

### Technical details
- Current CRM IR widgets: `BoardEngine`, `RecordFieldList`, `FieldRenderer`, `ActivityFeed`, `StatTile`, `SegmentedBar`.
- Current time widgets: `MonthGrid`, `TimeGrid` via `widget.time.month(...)` and `widget.time.week(...)`.
- Current scheduling widgets: `MatrixGrid`-based availability/poll widgets and booking picker via `widget.schedule.*`.
