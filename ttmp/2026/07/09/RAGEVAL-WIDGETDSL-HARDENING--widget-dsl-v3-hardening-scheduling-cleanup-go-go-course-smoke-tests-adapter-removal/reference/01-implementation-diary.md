---
Title: Implementation Diary
Ticket: RAGEVAL-WIDGETDSL-HARDENING
Status: active
Topics:
    - ui-dsl
    - widget-ir
    - frontend-architecture
    - xgoja
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/scripts/focused-checks.mjs
      Note: focused no-runner regression checks
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/MonthGrid/MonthGrid.logic.ts
      Note: focused month-boundary test target
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/TimeGrid/TimeGrid.logic.ts
      Note: focused lane-packing test target
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/TimeGrid/TimeGrid.tsx
      Note: TimeGrid allDay contract cleanup and logic extraction
    - Path: repo://packages/rag-evaluation-site/src/components/organisms/MeetingPollPanel/MeetingPollPanel.tsx
      Note: readOnly editability/action fix
    - Path: repo://packages/rag-evaluation-site/src/widgets/styleBy.logic.ts
      Note: focused StyleBy fallback test target
ExternalSources: []
Summary: Chronological diary for widget.dsl v3 hardening follow-up work.
LastUpdated: 2026-07-09T00:00:00Z
WhatFor: Record implementation steps, validation, failures, and review instructions for the hardening ticket.
WhenToUse: Read before resuming scheduling cleanup, go-go-course smoke coverage, or adapter removal.
---


# Diary

## Goal

Capture the follow-up hardening work after the `widget.dsl` v3 implementation and `go-go-course` host migration: scheduling correctness cleanup, browser smoke coverage, and adapter removal.

## Step 1: Ticket setup and first-slice planning

Created this ticket to group the three immediate follow-ups from the previous Widget DSL v3 workstream. The first implementation slice is intentionally small: fix the scheduling correctness/API findings before starting broader browser-test and adapter-removal work.

### Prompt Context

**User prompt (verbatim):** "do one ticket for 1, 2, 3.

another ticket for 4 

then work the first ticket"

**Assistant interpretation:** Create one docmgr ticket for scheduling cleanup, go-go-course browser smoke coverage, and adapter removal; create a second ticket for porting the Doodle example to widget.dsl v3; then begin work on the first ticket.

**Inferred user intent:** Split follow-up work into two manageable ticket tracks and start with the hardening work that protects the v3 migration before doing another example port.

**Commit (code):** N/A — setup in progress.

### What I did
- Created `RAGEVAL-WIDGETDSL-HARDENING`.
- Created `DOODLE-WIDGETDSL-V3`.
- Added design plans, task lists, and diaries for both tickets.
- Chose WP1 scheduling cleanup as the first implementation slice.

### Why
- The previous ticket completed broad v3 implementation, but its diaries and review docs still point to correctness/test gaps and adapter debt.
- Separating the Doodle port keeps the demo migration from blocking hardening work.

### What worked
- Ticket creation and plan scaffolding completed cleanly.

### What didn't work
- N/A.

### What I learned
- The immediate implementation work should start with small correctness fixes rather than another broad migration.

### What was tricky to build
- Scope control: item 3 (adapter removal) is large, so the ticket plan explicitly stages it after correctness fixes and browser smoke coverage.

### What warrants a second pair of eyes
- Whether `TimeGrid.allDay` should be implemented now or removed from the contract until a later calendar-row design exists.

### What should be done in the future
- Add a second diary step after the first code changes land, with exact validation commands and failures.

### Code review instructions
- Start with `design-doc/01-hardening-plan.md` and `tasks.md`.
- For the first code slice, review `MeetingPollPanel`, `TimeGrid`, package exports, and focused tests.

### Technical details
- Related source-ticket context: `RAGEVAL-SCHEDULE-WIDGETS` review findings and diary Steps 28–29.

## Step 2: Scheduling cleanup and focused logic checks

Implemented the first hardening slice from WP1. The changes close the concrete scheduling review findings that were cheapest and safest to fix now: read-only meeting polls no longer leave editable cells active, `TimeGrid` no longer advertises unsupported all-day block support, scheduling DTOs and presets are intentionally exported, and the riskiest pure logic now has a focused Node-based check script.

This step deliberately did not add a full React adapter/action-context test harness. The package currently has no Vitest/RTL setup, so I kept the first test slice dependency-free and extracted pure logic into `.logic.ts` modules that can be checked with Node 22's type-stripping support.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Start executing the first hardening ticket, beginning with the scheduling cleanup/test pass.

**Inferred user intent:** Stabilize known correctness/API gaps before moving on to browser smoke tests and adapter removal.

**Commit (code):** N/A — changes not committed yet.

### What I did
- Fixed `MeetingPollPanel` so `readOnly` clears `editableRowKey` and removes `onCell` action emission.
- Removed unsupported `allDay` from `TimeGridBlock` and `TimeGridBlockSpec`; `CalendarWeekPanel` now filters all-day `CalendarEvent`s before converting to timed blocks.
- Added public exports for scheduling DTOs and presets:
  - root exports `./scheduling`;
  - `widgets/index.ts` exports `./presets`;
  - package subpath exports for `./scheduling` and `./widgets/presets`;
  - Vite library entries for both subpaths.
- Extracted pure logic modules:
  - `TimeGrid.logic.ts` for `timeParts` and `packTimeGridColumn`;
  - `MonthGrid.logic.ts` for month parsing/shifting and calendar-cell generation;
  - `styleBy.logic.ts` for fallback/style selection before CSS-var conversion.
- Added `scripts/focused-checks.mjs` and package script `test:focused`.
- Updated hardening task state for completed WP1 items.

### Why
- These were explicit review findings in the source scheduling ticket.
- `readOnly` that only hides submit controls is a correctness bug: users could still toggle cells.
- A type surface that accepts `allDay` while silently dropping it is misleading; timed `TimeGrid` should not claim all-day support until an all-day row exists.
- Public organism prop types already reference scheduling DTOs, so scheduling needs an intentional public export path.
- Focused pure checks give quick regression coverage without introducing a new frontend test framework.

### What worked
- `pnpm --dir packages/rag-evaluation-site test:focused` passed.
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `pnpm --dir packages/rag-evaluation-site build` passed.
- `pnpm --dir packages/rag-evaluation-site pack:smoke` passed.
- `pnpm --dir packages/rag-evaluation-site consumer:smoke` passed.
- Manual subpath import checks confirmed `dist/scheduling.js` exports `AVAILABILITY_STATES` and `dist/widgets/presets.js` exports `availabilityMatrix`.

### What didn't work
- I did not add adapter/action-context tests in this slice because the package does not currently have a lightweight React/adapter test harness. I split that into a remaining task instead of adding a broad test dependency as part of a cleanup patch.

### What I learned
- Node 22's `--experimental-strip-types` is enough for small dependency-free TypeScript logic checks, as long as the tested modules avoid JSX/CSS imports.
- Extracting geometry/style selection into `.logic.ts` files improves testability without changing rendered component behavior.

### What was tricky to build
- `MonthGrid` originally mixed calendar-cell generation with JSX rendering. The extraction had to preserve marker and `onSelect` attachment in the component while moving only deterministic date/bounds logic to the testable module.
- `StyleBySpec` resolution returns CSS vars in the public helper, but fallback behavior is easier to test one layer earlier. I added `resolveStyleByStyle(...)` so the fallback chain can be asserted without importing context barrels or CSS-adjacent React code.

### What warrants a second pair of eyes
- The `CalendarWeekPanel` all-day policy: all-day domain events remain valid `CalendarEvent`s, but the week panel now omits them from the timed grid. A future all-day row would be better for production calendar use.
- Whether dependency-free Node checks are enough for this package, or whether the next frontend hardening pass should introduce Vitest/RTL for adapter/action-context tests.

### What should be done in the future
- Add focused adapter/action-context tests once a renderer harness exists.
- Implement a real all-day row in `TimeGrid`/`CalendarWeekPanel` if all-day events become product requirements.
- Continue with WP2: committed go-go-course Playwright smoke coverage.

### Code review instructions
- Review `MeetingPollPanel.tsx` first: `canEdit`, `editableRowKey`, and `onCell` should all agree.
- Review `TimeGrid.tsx` and `TimeGrid.logic.ts` together; component behavior should be the same for timed blocks.
- Review `MonthGrid.tsx` and `MonthGrid.logic.ts` together; generated cell dates/bounds should match previous behavior.
- Review public exports in `src/index.ts`, `src/widgets/index.ts`, `src/widgets/presets/index.ts`, `package.json`, and `vite.config.ts`.
- Validate with:
  - `pnpm --dir packages/rag-evaluation-site test:focused`
  - `pnpm --dir packages/rag-evaluation-site typecheck`
  - `pnpm --dir packages/rag-evaluation-site build`
  - `pnpm --dir packages/rag-evaluation-site pack:smoke`
  - `pnpm --dir packages/rag-evaluation-site consumer:smoke`

### Technical details
- `test:focused` runs `node --experimental-strip-types scripts/focused-checks.mjs`.
- The focused checks cover:
  - overlapping/back-to-back `TimeGrid` lane packing;
  - February 2026 month-cell generation, adjacent-month cells, min/max disabled states, today/selected flags, and month shifting;
  - `StyleBySpec` direct, mapped, fallback-key, and fallback-style selection.
