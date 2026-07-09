---
Title: ""
Ticket: ""
Status: ""
Topics: []
DocType: ""
Intent: ""
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/src/crm/types.ts
      Note: CRM domain DTOs and field schema
    - Path: repo://packages/rag-evaluation-site/src/widgets/WidgetRenderer.crm.stories.tsx
      Note: CRM renderer story coverage
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/engines.ts
      Note: CRM/calendar/time Widget IR contracts
    - Path: repo://packages/rag-evaluation-site/src/widgets/presets/crm.ts
      Note: Existing CRM preset composition layer
    - Path: repo://pkg/widgetdsl/typescript.go
      Note: TypeScript declarations for current widget.dsl v3 helpers
    - Path: repo://pkg/widgetdsl/v3.go
      Note: widget.dsl v3 namespaces for time/schedule/ui and proposed CRM extension point
ExternalSources: []
Summary: ""
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---


# Workshop CRM Screens and Widget DSL Plan

## Executive summary

This ticket designs a CRM application for a two-person AI engineering workshop business. The product combines the existing CRM widget stack with calendar/time widgets, scheduling widgets, generic UI primitives, tables, forms, course widgets, and CMS/media widgets.

The current codebase already has the rendering engines required for a credible application: pipeline boards, deal cards, typed record fields, activity feeds, stat tiles, segmented funnels, month grids, week time grids, availability polls, forms, split panes, panels, data tables, course shells, and media/library surfaces. The main gap is namespace exposure: the React/IR side contains CRM presets in `packages/rag-evaluation-site/src/widgets/presets/crm.ts`, while the Go `widget.dsl` v3 root namespace does not yet expose a first-class `widget.crm` object. The JS sketches in this document assume a thin `widget.crm.*` namespace that emits the existing IR components and follows the same style as `widget.time.*`, `widget.schedule.*`, and `widget.ui.*`.

## Evidence from the current widget system

- CRM domain DTOs exist in `packages/rag-evaluation-site/src/crm/types.ts`: `Contact`, `Company`, `Deal`, `Pipeline`, `Stage`, `Activity`, `Task`, and custom field definitions.
- CRM presets exist in `packages/rag-evaluation-site/src/widgets/presets/crm.ts`: `pipelineBoard`, `contactRecord`, `recordFieldList`, `activityFeed`, `statTile`, `pipelineFunnel`, `crmDashboard`, and `tasksInbox`.
- CRM Storybook coverage exists in `packages/rag-evaluation-site/src/widgets/WidgetRenderer.crm.stories.tsx` for pipeline board, record fields, contact record page, dashboard, and tasks inbox.
- Engine-level IR props for CRM live in `packages/rag-evaluation-site/src/widgets/ir/engines.ts`: `BoardEngineWidgetProps`, `RecordFieldListWidgetProps`, `FieldRendererWidgetProps`, `ActivityFeedWidgetProps`, and `StatTileWidgetProps`.
- Calendar/time DSL support exists in `pkg/widgetdsl/v3.go` and `pkg/widgetdsl/typescript.go`: `widget.time.month(...)`, `widget.time.week(...)`, `widget.time.range.week(...)`, and `widget.time.intent.*`.
- Scheduling DSL support exists in `pkg/widgetdsl/v3.go`: `widget.schedule.availabilityPoll(...)`, `widget.schedule.pollSummary(...)`, and `widget.schedule.bookingPicker(...)`.
- Generic UI support exists in `widget.ui`: split panes, stacks, inline layout, panels/cards, forms, form rows, inputs, status, metadata, empty state, buttons, and share links.

## Product model

Core records:

1. `Organization`: company, university, agency, conference, or enterprise lead.
2. `Contact`: sponsor, buyer, technical champion, training coordinator, participant, procurement contact.
3. `Opportunity`: a workshop sales deal with stage, expected value, requested format, close date, and next action.
4. `WorkshopRun`: a scheduled delivery engagement after an opportunity is won.
5. `Session`: an individual training block inside a workshop run.
6. `Task`: follow-up, proposal, invoice, prep, logistics, post-workshop feedback, or renewal task.
7. `Activity`: notes, emails, calls, meetings, stage changes, field changes, proposal sent, feedback received.
8. `Asset`: slides, handouts, exercises, datasets, environment setup guide, recording, invoice/proposal PDF.

Important relationships:

- Organization has many contacts.
- Opportunity belongs to one organization and has many contacts.
- WorkshopRun is created from a won opportunity.
- WorkshopRun has sessions, participants, assets, tasks, and feedback.
- Activities can attach to any major record.

## Screen set

### 1. Command center

Purpose: morning operating view. It shows pipeline health, upcoming workshop commitments, overdue tasks, and recent activity.

Widgets:

- `widget.crm.dashboard(...)` or explicit `StatTile`/`SegmentedBar` composition.
- `widget.time.week(...)` for delivery/prep schedule.
- `widget.crm.tasksInbox(...)` for next actions.
- `widget.crm.activityFeed(...)` for recent changes.
- `widget.ui.splitPane(...)` and `widget.ui.stack(...)` for layout.

### 2. Workshop sales pipeline

Purpose: manage opportunities from inbound lead through delivered workshop.

Pipeline stages:

1. New lead
2. Discovery scheduled
3. Proposal draft
4. Proposal sent
5. Contracting
6. Won / Scheduled
7. Delivered
8. Expansion

Widgets:

- `widget.crm.pipelineBoard(pipeline, opportunities, { summaries })`.
- `BoardEngine` for drag/drop movement.
- `DealCard` for each opportunity.
- `SegmentedBar` for funnel health.

### 3. Organization / contact record

Purpose: one place to understand a buyer, the internal champion, stakeholder map, fields, history, and related opportunities.

Widgets:

- `widget.crm.recordFieldList(...)` for typed custom fields.
- `widget.crm.activityFeed(...)` for timeline.
- `widget.data.table(...)` for related opportunities or participants.
- `widget.ui.splitPane(...)` for detail + activity.

### 4. Workshop calendar

Purpose: schedule discovery calls, proposal deadlines, prep blocks, workshop delivery blocks, office hours, and post-workshop review.

Widgets:

- `widget.time.month(events, ...)` for month overview with compact markers.
- `widget.time.week(events, ...)` for detailed week blocks.
- `widget.schedule.availabilityPoll(...)` for finding dates with a client.
- `widget.schedule.bookingPicker(...)` for booking discovery calls or prep sessions.

### 5. Workshop run room

Purpose: delivery operations for a won engagement. It tracks agenda, participants, assets, logistics, facilitator split, and completion status.

Widgets:

- `widget.course.shell(...)` or course/handout widgets for agenda and materials.
- `widget.cms.mediaLibrary(...)` for slides, handouts, exercise repos, recordings.
- `widget.time.week(...)` for session schedule.
- `widget.ui.checkList` / `widget.crm.tasksInbox(...)` for logistics.
- `widget.ui.shareLink(...)` for participant onboarding URL.

### 6. Lead intake and proposal builder

Purpose: capture a new inquiry and turn it into a proposal/deal.

Widgets:

- `widget.ui.form(...)`, `widget.ui.formRow(...)`, `widget.ui.textInput(...)`, `widget.ui.textareaInput(...)`, and `widget.ui.selectInput(...)`.
- `widget.ui.status(...)` for save/submission state.
- `widget.crm.recordFieldList(..., { mode: "edit" })` once a draft record exists.

### 7. Post-workshop success / renewal

Purpose: close the loop after delivery. It shows feedback, outcomes, artifacts sent, testimonial candidates, expansion opportunities, and next recommended workshop.

Widgets:

- `widget.data.table(...)` for participant feedback rows.
- `widget.crm.activityFeed(...)` for post-workshop timeline.
- `widget.crm.pipelineBoard(...)` filtered to expansion opportunities.
- `widget.cms.mediaLibrary(...)` for deliverables.

## Proposed `widget.crm` namespace

This is a thin namespace over existing IR components and TypeScript presets. It should be added to `widget.dsl` v3 so Goja authors can write CRM pages without importing React-side presets.

Proposed helpers:

- `widget.crm.pipelineBoard(pipeline, deals, options?)`
- `widget.crm.dashboard(pipeline, summaries, options?)`
- `widget.crm.record(record, fieldDefs, options?)`
- `widget.crm.recordFields(values, fieldDefs, options?)`
- `widget.crm.activityFeed(activities, options?)`
- `widget.crm.tasksInbox(tasks, options?)`
- `widget.crm.funnel(pipeline, summaries, options?)`
- `widget.crm.intent.openDeal(id)`
- `widget.crm.intent.moveDeal(id, toStage)`
- `widget.crm.intent.updateField(recordId, key, value)`
- `widget.crm.intent.completeTask(taskId)`

## Implementation plan

1. Add a `widget.crm` namespace in `pkg/widgetdsl/v3.go` and TypeScript declarations in `pkg/widgetdsl/typescript.go`.
2. Implement helpers by emitting existing IR nodes: `BoardEngine`, `RecordFieldList`, `ActivityFeed`, `StatTile`, `SegmentedBar`, `Panel`, `Stack`, and `SplitPane`.
3. Add golden fixtures for the CRM helpers.
4. Add Storybook fixture pages for the workshop CRM screens.
5. Build an xgoja example app similar to the Doodle site, backed by SQLite.
6. Add browser smoke tests for dashboard, pipeline drag/select actions, calendar selection, and intake form submission.

## Risks and open questions

- The current CRM presets are TypeScript-side presets. The JS sketches require either a new Go-side `widget.crm` namespace or local JS helper functions that emit the same IR.
- The existing CRM model is intentionally generic. Workshop-specific fields should be data (`FieldDef`) rather than new widgets.
- Calendar month markers must remain compact. Detailed participant/client names belong in side panels or records.
- Some screens need real persistence and server actions before they become interactive beyond client-side dispatch.
