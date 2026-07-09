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
    - Path: repo://pkg/widgetdsl/v3.go
      Note: Current builder patterns investigated
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
