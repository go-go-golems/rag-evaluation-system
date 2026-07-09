---
Title: Implementation Diary
Ticket: DOODLE-WIDGETDSL-V3
Status: active
Topics:
    - widget-dsl
    - xgoja
    - sqlite
    - doodle
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Chronological diary for porting the Doodle scheduling example to widget.dsl v3.
LastUpdated: 2026-07-09
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
