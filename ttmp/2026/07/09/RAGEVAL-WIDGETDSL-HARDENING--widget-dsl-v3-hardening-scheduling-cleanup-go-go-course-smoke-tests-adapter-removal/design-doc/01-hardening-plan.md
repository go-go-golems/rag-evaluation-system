---
Title: Hardening Plan
Ticket: RAGEVAL-WIDGETDSL-HARDENING
Status: active
Topics:
    - ui-dsl
    - widget-ir
    - frontend-architecture
    - xgoja
DocType: design-doc
Intent: short-term
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/src/components/organisms/MeetingPollPanel/MeetingPollPanel.tsx
      Note: readOnly/editability bug to fix
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/TimeGrid/TimeGrid.tsx
      Note: allDay contract and lane packing test target
    - Path: repo://packages/rag-evaluation-site/src/index.ts
      Note: public scheduling export decision
    - Path: ws://go-go-course/cmd/go-go-course/lib/widget-dsl-v3-adapter.js
      Note: transitional adapter to shrink/remove
Summary: Hardening plan for the post-widget.dsl-v3 work: scheduling cleanup/tests, committed go-go-course browser smoke coverage, and staged adapter removal.
LastUpdated: 2026-07-09
WhatFor: Use this as the implementation checklist for the first follow-up ticket after the widget.dsl v3 migration.
WhenToUse: Before changing scheduling widgets, go-go-course smoke tests, or native widget.dsl page rewrites.
---

# Hardening Plan

## Scope

This ticket groups the three immediate follow-ups from the Widget DSL v3 workstream:

1. **Scheduling widget cleanup and tests** from `RAGEVAL-SCHEDULE-WIDGETS` review findings.
2. **Committed browser smoke coverage for `go-go-course`** so shell navigation, login, upload, and dynamic session pages are tested rather than manually rechecked.
3. **Incremental removal of the `go-go-course` v3 compatibility adapter** by rewriting page modules to native `widget.dsl` APIs.

## Work packages

### WP1 — Scheduling cleanup and focused tests

Fix correctness/API issues before expanding more DSL work:

- `MeetingPollPanel.readOnly` must disable editable cells and suppress cell actions.
- `TimeGrid.allDay` must have an explicit contract. Preferred first pass: remove/omit all-day support from the timed engine until an all-day row is implemented, and document the choice in stories/types.
- Decide public exports for scheduling DTOs and presets. Preferred first pass: export `./scheduling` from the package root because public organism props already mention scheduling DTOs.
- Add focused tests for:
  - `TimeGrid` lane packing;
  - `MonthGrid` bounds/date behavior;
  - `resolveStyleByVars` fallback;
  - Matrix/scheduling adapter action context if a suitable test harness exists.

### WP2 — go-go-course Playwright smoke

Add a committed script/test that validates behavior the manual session proved:

- open hotreload site;
- set display name to `admin_manuel`;
- upload a copied Pi session JSONL fixture;
- click sidebar items, including dynamic session transcript/visualize pages;
- fail on console errors, unknown widgets, render errors, or `[object Object]` output.

### WP3 — Native widget.dsl page rewrites

Reduce adapter debt in small slices:

1. rewrite DSL examples first;
2. rewrite handouts/slides;
3. rewrite admin CMS;
4. rewrite transcript/visualize pages;
5. remove adapter helpers as call sites disappear;
6. run `cmd/widgetdsl-migration-checker` after each slice.

## Acceptance criteria

- Scheduling cleanup passes `pnpm --dir packages/rag-evaluation-site typecheck` and targeted tests.
- Browser smoke is committed and documented with a reproducible command.
- Migration checker output improves from one central adapter raw escape hatch toward zero ordinary-page findings.
- Diary, changelog, and related-file links are updated for each coherent step.
