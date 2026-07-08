---
Title: Implementation Diary
Ticket: RAGEVAL-SCHEDULE-WIDGETS
Status: active
Topics:
    - ui-dsl
    - widget-ir
    - design-system
    - react
    - frontend-architecture
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/src/components/atoms/DateTile/DateTile.module.css
      Note: sm-size overflow fix (weekday hidden, per-size day scaling)
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/MatrixGrid/MatrixGrid.tsx
      Note: Flagship grid engine + MatrixCellPayload (relaxed generic bound)
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/MonthGrid/MonthGrid.tsx
      Note: Calendar-month engine + MonthGridDayPayload day-cell contract
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/TimeGrid/TimeGrid.tsx
      Note: Week/day calendar engine + packColumn lane packing + TimeGridBlockPayload contract
    - Path: repo://packages/rag-evaluation-site/src/components/organisms/BookingPagePanel/BookingPagePanel.tsx
      Note: 1:1 booking organism (SplitPane MonthGrid + slot tiles)
    - Path: repo://packages/rag-evaluation-site/src/components/organisms/MeetingPollPanel/MeetingPollPanel.tsx
      Note: First Doodle organism — participant poll screen composed from base engines
    - Path: repo://packages/rag-evaluation-site/src/components/organisms/PollResultsPanel/PollResultsPanel.tsx
      Note: Organizer results organism (SegmentedBar list + finalize)
    - Path: repo://packages/rag-evaluation-site/src/scheduling/types.ts
      Note: Scheduling domain DTOs (pure data, mirrors src/context)
    - Path: repo://packages/rag-evaluation-site/src/widgets/WidgetRenderer.scheduling.stories.tsx
      Note: IR renderer stories proving the poll renders from serialized nodes
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/index.ts
      Note: |-
        ir.ts split into core/actions/cells/engines/props modules
        IR barrel created by the split from the former ir.ts file
    - Path: repo://packages/rag-evaluation-site/src/widgets/presets/scheduling.ts
      Note: |-
        availabilityMatrix schedule.dsl preset (engine -> IR node)
        pollResults/monthCalendar/weekCalendar presets
    - Path: repo://packages/rag-evaluation-site/src/widgets/styleBy.ts
      Note: resolveStyleByVars — StyleBySpec consumption (value->styleKey->style)
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-06T18:48:57.302545525-04:00
WhatFor: ""
WhenToUse: ""
---








# Diary

## Goal

Capture the implementation journey for the calendar/scheduling widget system in
`packages/rag-evaluation-site`: a set of **generic base engines** (grid, month
calendar, time grid, segmented bar, cycle cell) plus, later, domain presets in a
`schedule.dsl` / `calendar.dsl` module. This diary records what was built, the
gotchas hit, and how to continue. Read it before resuming.

## Step 1: Ticket, intern guide, and reMarkable upload

Started the ticket by first reading the existing widget system end to end so the
proposal would match reality rather than invent a parallel mechanism. Wrote a
long intern-facing analysis/design/implementation guide and uploaded it to
reMarkable. No product code changed in this step — it is the design foundation
the later steps build on.

### Prompt Context

**User prompt (verbatim):** "alright, create a new docmgr ticket to implement the calendar / schedule stuff above, but using base widgets. Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.\n\n\nThen implement a few base widgets and then stop to let me check their looks in storybook"

**Assistant interpretation:** Create a docmgr ticket; write a thorough, technical, intern-oriented guide (prose + bullets + pseudocode + diagrams + API/file references) describing the widget system and the calendar/scheduling design; store it in the ticket and upload to reMarkable; then implement a few base widgets and stop for a Storybook review.

**Inferred user intent:** Get a durable, teachable design artifact plus a small, reviewable first slice of real components — building the "right" reusable base rather than one-off scheduling widgets.

**Commit (code):** N/A — no VCS in this workspace (`git repo: false`); changes tracked via docmgr only.

### What I did
- Read `GUIDELINES.md`; `src/widgets/{ir.ts,actions.ts,cellRenderers.tsx,registry.ts,WidgetRenderer.tsx,defaultRegistry.ts}`; `src/context/types.ts`; and the `StepList` / `CourseLessonPanel` folders for the six-file convention.
- Created ticket `RAGEVAL-SCHEDULE-WIDGETS` and this diary.
- Wrote `design-doc/01-...-analysis-design-and-implementation-guide.md` (Parts A–E: the world as it is / the idea / the widgets / the backend / build order, plus API- and file-reference tables).
- Related six source files to the design doc; updated the changelog.
- Uploaded the guide to reMarkable: `remarquee upload bundle ... --remote-dir /ai/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS` → `OK: uploaded`.

### Why
- The codebase already encodes "lambda-as-data" (`ActionSpec`, `CellSpec`, accessor/template specs). Building on that — rather than inventing a mechanism — is what lets the new engines be driven identically from React, the IR, and the DSL.

### What worked
- The guide grounded every claim in real identifiers/paths, so it doubles as a lookup table for the coding steps.
- reMarkable upload succeeded on the first try.

### What didn't work
- N/A for this step (no code compiled yet).

### What I learned
- The renderer is a flat `registry.get(node.type)` over a merged registry, which is exactly why cross-DSL composition is free — but `createWidgetRegistry` throws on duplicate `type`, so type strings are a global namespace.

### What was tricky to build
- Deciding the generic/domain split up front: which widgets are engines (own space/time) vs presets (own meaning). Resolved with the litmus test "how are cells arranged?" = engine; "what does a vote mean?" = preset.

### What warrants a second pair of eyes
- The backend section (server actions + REST + server-computed tallies) is a proposal, not yet validated against `packages/web`.

### What should be done in the future
- Confirm product-flow priority (group poll / 1:1 booking / full calendar); decide on namespaced `type` strings before the first adapter.

### Code review instructions
- Start at `design-doc/01-...-implementation-guide.md`; cross-check its API Reference table against `src/widgets/*`.

### Technical details
- reMarkable target: `/ai/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS/RAGEVAL Calendar Scheduling Widgets Intern Guide.pdf`.

## Step 2: First base-widget slice (atoms + two engines)

Implemented the first reviewable slice React-first (no IR adapters yet, per
"React first, Widget IR later"): three atoms and the two flagship engines, each
with Storybook stories. The decisive artifact is `MatrixGrid`, which proves the
whole thesis by hosting swappable cells through a stable cell contract.

### Prompt Context

**User prompt (verbatim):** (see Step 1 — the "then implement a few base widgets and then stop" clause.)

**Assistant interpretation:** Build a small, coherent set of base widgets with stories good enough to judge looks in Storybook, then stop.

**Inferred user intent:** Validate the base-engine approach visually before going wide.

**Commit (code):** N/A — no VCS in this workspace.

### What I did
- Atoms: `DateTile`, `RatioBadge`, `CycleCell` (the first cell-contract citizen).
- Molecule engines: `SegmentedBar` (generalizes `ContextBudgetBar`), `MatrixGrid` (flagship, defines `MatrixCellPayload`).
- Six-file set minus the `.widget.*` IR pair for each; exported from the `atoms` and `molecules` barrels.
- Stories including the two decisive `MatrixGrid` demos: Mode A (`renderCell` + `CycleCell`, interactive poll with editable "You" row + `RatioBadge` tally footer) and Mode B (explicit `cells` matrix as a plan-comparison grid).
- Verified `typecheck` and `build-storybook`.

### Why
- A working `MatrixGrid` with both injection modes is the minimum that demonstrates "generic engine + swappable cells" is real and not just a diagram.

### What worked
- `typecheck` passed and `build-storybook` completed after the fixes below.

### What didn't work
- `noUncheckedIndexedAccess` is on, so `states[(idx+1)%len]` is `string | undefined`:
  `src/components/atoms/CycleCell/CycleCell.tsx(30,2): error TS2322: Type 'string | undefined' is not assignable to type 'string'.` Fixed with `?? current`.
- The generic bound `MatrixGrid<Row extends MatrixRow>` rejected a plain interface:
  `error TS2344: Type 'Respondent' does not satisfy the constraint 'MatrixRow'. Index signature for type 'string' is missing in type 'Respondent'.` (Interfaces lack an implicit index signature.) Relaxed to `<Row = MatrixRow>` with `as MatrixRow` casts in the two default accessors.
- Storybook `Meta` demanded required props on `render`-only stories:
  `error TS2322: ... Property 'args' is missing`. Fixed by adding `args: { rows: [], columns: [] }` to the `MatrixGrid` meta.

### What I learned
- `noUncheckedIndexedAccess` must be assumed for every future engine (`MonthGrid`, `TimeGrid` do heavy indexed access).
- A generic constraint of `Record<string, unknown>` is hostile to interface-typed callers; leaving `Row` unconstrained with internal casts is friendlier.

### What was tricky to build
- The cell-contract seam: the grid must resolve `value` and expose a grid-owned `onAction` while letting the cell override the emitted value (`onAction({ value: next })`). Cause: without an override channel the editable `CycleCell` could not report its next state. Solution: `onAction(extra?)` merges `extra.value` over the resolved value before calling `onCell`.

### What warrants a second pair of eyes
- `MatrixGrid` sticky header + sticky row-header z-index/`overflow` interplay inside `ScrollRegion`-style containers — verify no clipping on narrow viewports.

### What should be done in the future
- Add `MonthGrid`, `TimeGrid`; then IR adapters + `StyleBySpec`; then the `scheduling/` domain module and presets.

### Code review instructions
- Start at `src/components/molecules/MatrixGrid/MatrixGrid.tsx` (`MatrixCellPayload`, `onAction`), then `src/components/atoms/CycleCell/CycleCell.tsx`.
- Validate: `pnpm --dir packages/rag-evaluation-site typecheck` and `... build-storybook`; in `pnpm storybook`, open *Molecules/MatrixGrid → PollWithCycleCells* and *ModeBExplicitCells*.

### Technical details
- Cell contract: `{ row, col, value, rowKey, rowIndex, colIndex, selected, editable, onAction }`.
- Story-local palettes (per GUIDELINES) instead of a premature `scheduling/` module.

## Step 3: DateTile sizing fix + MonthGrid engine + diary reformat

Fixed the small `DateTile` (overflowing text), added the `MonthGrid` calendar
engine (the base behind date pickers, poll slot-pickers, and activity heatmaps),
and rewrote this diary to follow the diary skill's strict Step format.

### Prompt Context

**User prompt (verbatim):** "the small datetile is too small and the text looks out.\n\nNow continue. \n\nWhat' StyleBySpec?\n\nAlso, read the diary skill and fix the formatting you use, you have to follow the skill guidelines (and retroactively fix if necessary)"

**Assistant interpretation:** The `sm` `DateTile` clips its text — fix it. Continue implementing base widgets. Explain what `StyleBySpec` is. Read the diary skill and reformat this diary (retroactively) to match it.

**Inferred user intent:** A visually correct small tile, forward progress on the engines, a clear definition of the `StyleBySpec` concept from the guide, and a diary that actually conforms to the team's diary standard.

**Commit (code):** N/A — no VCS in this workspace.

### What I did
- `DateTile`: root `overflow: hidden`; `sm` enlarged to 40×40 with the weekday row hidden (two lines only) so 11–12px mono glyphs never clip; per-size `day` font scaling (`sm` 1.15em, `md` 1.35em, `lg` 1.75em); `white-space: nowrap` on all three lines.
- Added `MonthGrid` (`.tsx` + `.module.css` + `.stories.tsx` + `index.ts`) with a `MonthGridDayPayload` day-cell contract mirroring `MatrixCellPayload`, UTC-safe date math, `weekStartsOn`, min/max bounds, markers (count or `styleKey` heat via `ContextStyleSet`), header prev/next, and a default day renderer + Mode-A `renderDay` override.
- Exported `MonthGrid` from the molecules barrel; verified `typecheck` + `build-storybook`.
- Rewrote this diary to the strict skill format (Goal + Step N with Prompt Context and all required subsections).

### Why
- The `sm` tile stacked three 11–12px lines into a 38px box → vertical overflow. Fewer lines + clipping + scaled numbers fixes it without arbitrary font literals beyond geometric `em` scaling.
- `MonthGrid` is the next engine both product flows (poll slot-picking, booking day selection) and the generic activity-heatmap use case all need.

### What worked
- `typecheck` clean; `build-storybook` completed successfully after fixing a barrel duplicate (below).

### What didn't work
- While adding the `MonthGrid` export I initially also re-added `PersonSummary` / `SearchField` to the molecules barrel, which already export later in the file (would be duplicate exports). Caught on read-back and reverted to adding only `MonthGrid`.

### What I learned
- Building the month matrix as a flat `weekCount * 7` array seeded at `Date.UTC(year, month, 1 - leading)` lets JS normalize adjacent-month days for free — no special-casing of month boundaries or leap years (verified via a `2028-02` story).

### What was tricky to build
- Avoiding timezone drift: a bare `YYYY-MM-DD` parsed as local time can shift a day. Cause: `new Date("2026-07-01")` is UTC midnight but `getDate()` is local. Solution: do all arithmetic with `Date.UTC(...)` and read back with `getUTC*`, and accept an explicit `todayISO` prop instead of calling `new Date()` (also keeps stories deterministic).

### What warrants a second pair of eyes
- `MonthGrid` marker contrast when a `styleKey` heat fill sits under the day number — confirm `--ctx-label` gives readable text on the darkest fills.
- Adjacent-month days are rendered muted and disabled; confirm that is the desired affordance vs. click-to-navigate.

### What should be done in the future
- Add `TimeGrid` (week/day engine).
- Add `StyleBySpec` to `ir.ts` (see definition below) when the IR-adapter milestone starts; it is an adapter-layer concept, so it is deliberately absent from these React-first base components.
- Create `scheduling/types.ts` + `fixtures.ts`; build presets (`availabilityMatrix`, `monthCalendar`, ...) and organism panels.

### Code review instructions
- `DateTile`: `src/components/atoms/DateTile/DateTile.module.css` (size blocks + `.sm .weekday { display:none }`).
- `MonthGrid`: `src/components/molecules/MonthGrid/MonthGrid.tsx` (`parseMonth`/`isoDate`/`shiftMonth`, `MonthGridDayPayload`).
- Validate: `pnpm --dir packages/rag-evaluation-site typecheck` and `... build-storybook`; in `pnpm storybook` open *Atoms/DateTile → Sizes* and *Molecules/MonthGrid*.

### Technical details — what `StyleBySpec` is
`StyleBySpec` is the **defunctionalized color function** for the IR/DSL layer — the
"how do I color this from data?" spec, analogous to how `ActionSpec` is the
event-handler spec and `CellSpec` is the render spec. A base engine cannot ship a
`(value) => color` closure across the Goja/JSON boundary, so instead the IR carries
data that an interpreter turns into a `ContextVisualStyle`:

```ts
interface StyleBySpec {
  field?: string;                  // which datum field to key on (default: the cell value)
  styleSet: ContextStyleSet;       // styleKey -> ContextVisualStyle (existing palette contract)
  map?: Record<string, string>;    // optional value -> styleKey remap
  fallbackStyleKey?: string;
}
```

Resolution at render time: take the value (or `row[field]`), optionally remap it
through `map`, look the resulting `styleKey` up in `styleSet.styles`, fall back to
`fallbackStyleKey`/`fallbackStyle`, and apply via `contextVisualStyleToCssVars`.
Adding it **once** makes `MatrixGrid`, `MonthGrid`, `SegmentedBar`, and the existing
context diagrams all recolorable purely from serialized data — a poll cell colored
by availability state, a calendar day colored by event category, a bar segment
colored by bucket, all with the same spec. It belongs on the IR adapters, not the
React base components (which already take real `ContextStyleSet` props directly), so
it lands with the adapter/`.widget.tsx` milestone, not now.

## Step 4: TimeGrid engine (base-engine layer complete)

Added `TimeGrid`, the week/day calendar engine — hour-ruled day columns with
absolutely-positioned event blocks, lane-packed so overlapping events sit side by
side. This completes the generic base-engine layer from the guide (`MatrixGrid`,
`SegmentedBar`, `MonthGrid`, `TimeGrid`) plus the supporting atoms. Everything is
still React-first; no IR adapters yet.

### Prompt Context

**User prompt (verbatim):** "alright, go ahead. continue keeping a detailed diary."

**Assistant interpretation:** Continue the build-order implementation (next engine = `TimeGrid`) and keep the diary in the strict skill format as I go.

**Inferred user intent:** Keep making forward progress on the reusable engines while maintaining a reviewable, standards-conformant record.

**Commit (code):** N/A — no VCS in this workspace.

### What I did
- Added `TimeGrid` (`.tsx` + `.module.css` + `.stories.tsx` + `index.ts`) with a `TimeGridBlockPayload` block contract mirroring the other engines.
- Implemented `packColumn` — overlap-cluster + lane assignment so concurrent blocks share a cluster lane count and line up.
- Wall-clock positioning parsed directly from the ISO `HH:MM` substring (no `Date`/timezone conversion in the base engine — that is a preset concern).
- Features: `hourStart`/`hourEnd`/`hourHeight`, sticky headers + gutter, per-hour clickable slots (`onSlotCreate`), `nowISO` indicator, `selectedBlockId`, and Mode-A `renderBlock`.
- Exported from the molecules barrel; verified `typecheck` + `build-storybook`.

### Why
- `TimeGrid` is the last engine both the full-calendar surface and the day/week planner need; finishing it closes the engine layer so the next milestone (IR adapters + presets) has a stable base.

### What worked
- `typecheck` clean and `build-storybook` succeeded first try after the gutter-alignment refinement.

### What didn't work
- No compile failures this step. One layout refinement mid-write: initial gutter rendered `hours + 1` fixed-height labels (total height overshot the body by one row); switched to `hourLabels.slice(0, -1)` positioned absolutely at `i * hourHeight` so labels sit on the gridlines.

### What I learned
- Parsing the ISO time-of-day as a plain string keeps the engine deterministic and timezone-free; converting UTC→local is exactly the kind of "meaning" that belongs in a preset, not the geometry engine — the same base/preset boundary the whole ticket is built on.

### What was tricky to build
- Side-by-side layout of overlapping events. Cause: a naive "one width for the whole day" wastes space and a per-block overlap count doesn't give stable widths. Symptom: blocks either overlap or shrink inconsistently. Solution: group blocks into overlap *clusters* (a new cluster starts when a block begins at/after the running cluster end), assign lanes within a cluster via a free-lane scan, and give every member the cluster's lane count so widths are uniform within the cluster. See `packColumn`; exercised by the overlapping trio (`b4`/`b5`/`b6`) in the *Week* story.

### What warrants a second pair of eyes
- `packColumn` cluster/lane logic on edge cases: zero-length blocks (min-clamped to 15 min), blocks starting before `hourStart` or ending after `hourEnd` (clamped), and back-to-back blocks that touch exactly (should not share a cluster).
- Sticky corner/gutter/header z-index stacking inside the scroll container.

### What should be done in the future
- IR-adapter milestone: add `StyleBySpec` to `ir.ts`; write `.widget.tsx` + `.widget.yaml` for each engine; add `timeWidgetRegistry` / `scheduleWidgetRegistry`; register in `defaultRegistry.ts` (mind the duplicate-`type` global namespace — use namespaced types like `calendar/TimeGrid`).
- `scheduling/types.ts` + `fixtures.ts`; presets (`availabilityMatrix`, `pollResults`, `bookingPage`, `monthCalendar`, `weekCalendar`); organism panels.

### Code review instructions
- Start at `src/components/molecules/TimeGrid/TimeGrid.tsx`: `packColumn` (lane packing) then the column render (block `top`/`height`/`left`/`width` from packed percentages).
- Validate: `pnpm --dir packages/rag-evaluation-site typecheck` and `... build-storybook`; in `pnpm storybook` open *Molecules/TimeGrid → Week* (overlap packing), *WithNowIndicator*, and *Interactive* (block-select + slot-create).

### Technical details
- Block position: `topPct = (start - rangeStart) / rangeMinutes * 100`; `heightPct` likewise; `left = lane/lanes * 100%`, `width = 1/lanes * 100%`.
- Base-engine layer now complete: atoms `DateTile` / `RatioBadge` / `CycleCell`; engines `MatrixGrid` / `SegmentedBar` / `MonthGrid` / `TimeGrid`. Next layer is IR adapters + DSL presets.

## Step 5: Scheduling domain module + first IR adapter + preset + WidgetRenderer story

Wired the first engine to the DSL end to end. Created the pure-data `scheduling`
domain module (DTOs + palettes + fixtures), added `StyleBySpec` / `CycleCellSpec`
/ `MatrixGridWidgetProps` to the IR, wrote the `MatrixGrid` IR adapter + manifest,
registered it, wrote the `availabilityMatrix` `schedule.dsl` preset, and added a
`Widget IR/Renderer/Scheduling` story that renders the Doodle poll from a
serialized node tree. This closes the loop the guide describes: base engine →
adapter → preset → registry → composed IR tree.

### Prompt Context

**User prompt (verbatim):** "ok, so go ahead. continue keeping a diary"

**Assistant interpretation:** Proceed with the domain-module + first-adapter milestone (per the plan discussed in prior turns) and keep the diary in the strict skill format.

**Inferred user intent:** Move from React-only engines to a DSL-drivable widget, proving the whole layering with a real IR story, while keeping the record continuous.

**Commit (code):** N/A — no VCS in this workspace.

### What I did
- `src/scheduling/`: `types.ts` (DTOs), `palettes.ts` (`availabilityStyleSet` / `eventStyleSet` — de-duplicated out of the stories), `fixtures.ts` (`sampleTeamSyncPoll`, `sampleTeamSyncTallies`, `sampleWeekEvents`), `index.ts`.
- `src/widgets/ir.ts`: added `StyleBySpec`, `CycleCellSpec`, `MatrixCellSpec`, `MatrixColumnWidgetSpec`, `MatrixValueSpec`, `MatrixGridWidgetProps`; registered `"MatrixGrid"` in `RagWidgetType` and the `WidgetProps` union.
- `MatrixGrid.widget.tsx` adapter (+ `MatrixGrid.widget.yaml` manifest): interprets the value-accessor spec, the cell spec (cycle → `CycleCell`, else DataTable-style `CellSpec`), Mode-B `cells`, the footer CellSpec (evaluated against each column's `meta`), and `onCellAction`.
- Registered `matrixGridWidget` in `dataWidgetRegistry` (`defaultRegistry.ts`).
- `src/widgets/presets/scheduling.ts`: `availabilityMatrix(poll, { tallies, editableResponseId })` → `WidgetNode`.
- `src/widgets/WidgetRenderer.scheduling.stories.tsx`: `AvailabilityPoll` (preset → renderer), `PollInPanel` (cross-widget IR composition), `HandAuthoredMatrix` (raw Mode-B IR).
- Verified `typecheck` + `build-storybook`.

### Why
- An IR story for a bare engine is verbose JSON; pairing the adapter with the `availabilityMatrix` preset makes the story a realistic one-liner and demonstrates the "opinionated skin over a generic engine" payoff.

### What worked
- `build-storybook` succeeded; the poll now renders both as a React story (`Molecules/MatrixGrid`) and from serialized IR (`Widget IR/Renderer/Scheduling`).

### What didn't work
- Type friction at the `MatrixRow` ↔ `JsonObject` boundary:
  - `MatrixGrid.widget.tsx(64): error TS2345: Argument of type 'MatrixRow' is not assignable to parameter of type 'JsonObject'. ... Type 'unknown' is not assignable to type 'JsonValue'.` — `renderCell` wants `JsonObject`; `MatrixRow`'s index type is `unknown`. Fixed with `row as unknown as JsonObject` at the two `renderCell` call sites.
  - `MatrixGrid.widget.tsx(113): error TS2345: ... property 'value' ... Type 'unknown' is not assignable to type 'string | number | boolean | null | undefined'.` — `WidgetActionContext.value` is a narrow union; spreading the whole `{ value: unknown }` payload broke it. Fixed by building the context explicitly and casting `value`.
- Initial `HandAuthoredMatrix` used `cell: { kind: "field", field: "name" }`, which renders the row name in every column (the non-cycle cell path reads a row *field*, not the resolved (row,col) value). Switched it to Mode-B `cells` (a `text` matrix) — the honest raw-IR shape.

### What I learned
- The non-cycle cell path is row-field-oriented (inherited from `DataTable`'s `CellSpec`), so per-cell *value* rendering wants either the `cycle` spec or Mode-B `cells`. A future `{ kind: "value" }` CellSpec that renders `payload.value` would make hand-authored value grids ergonomic.
- Keeping `scheduling/` free of IR/React means the preset (which needs both `component()` and the DTOs) lives in `src/widgets/presets/`, not in the domain module — the same separation `context/` keeps.

### What was tricky to build
- The footer contract. The base `MatrixGrid` footer renders per column, but IR can't compute a tally. Cause: the tally is data the preset already has. Solution: the preset stuffs `{ yes, total }` into each column's `meta`, and the adapter evaluates the footer `CellSpec` (`"${yes}/${total}"`) against `col.meta` as the synthetic row. So "compute server-side, pass as meta, render as template" — no computation in the widget.

### What warrants a second pair of eyes
- `resolveValue` accessor (`mapField` / `template` / default) — confirm the `template` `${colId}` substitution and missing-map handling.
- The `onCellAction` server payload uses `{ kind: "path", path: "rowKey" | "colId" | "value" }` template parts resolved by `resolveActionPayload`; verify the dispatch context keys line up (`rowKey`/`colId`/`value`).
- `rows: poll.responses as unknown as JsonObject[]` cast in the preset — acceptable because responses are JSON-shaped, but worth a glance.

### What should be done in the future
- Adapters + manifests for `SegmentedBar`, `MonthGrid`, `TimeGrid` (+ `calendar/`-namespaced types), and their `WidgetRenderer` stories.
- Add a `{ kind: "value" }` CellSpec (or a `colorBy` application path) so `StyleBySpec` is actually consumed by the renderer, not just declared.
- Presets: `pollResults`, `bookingPage`, `monthCalendar`, `weekCalendar`; then organism panels (`MeetingPollPanel`, ...).

### Code review instructions
- Start at `src/widgets/presets/scheduling.ts` (`availabilityMatrix`) → `src/components/molecules/MatrixGrid/MatrixGrid.widget.tsx` (spec interpretation) → `src/widgets/WidgetRenderer.scheduling.stories.tsx`.
- Validate: `pnpm --dir packages/rag-evaluation-site typecheck` and `... build-storybook`; in `pnpm storybook` open *Widget IR/Renderer/Scheduling → AvailabilityPoll* and compare with *Molecules/MatrixGrid → PollWithCycleCells*.

### Technical details
- Domain module mirrors `src/context/`: `types` + `palettes` + `fixtures`, no React/IR.
- Registry: `matrixGridWidget` (`module: data.dsl`, `type: "MatrixGrid"`) added to `dataWidgetRegistry`; merged via `defaultWidgetRegistry`. Namespaced types (`calendar/…`, `schedule/…`) deferred until a second DSL specializes an engine.

## Step 6: MeetingPollPanel organism (first full Doodle screen)

Composed the base engines and atoms into `MeetingPollPanel` — the
participant-facing Doodle screen (title + meta + deadline/response line + the
availability grid with an editable "You" row + tally footer + a submit row).
First organism in the ticket; presentational (DTO in, callbacks out) per the
GUIDELINES organism contract, with the interactive state held by the story.

### Prompt Context

**User prompt (verbatim):** "Go ahead. you chose."

**Assistant interpretation:** Pick the most valuable next step myself and proceed (keeping the diary). I chose the organism over the remaining engine adapters because it yields a complete, reviewable screen that exercises the whole stack.

**Inferred user intent:** Keep momentum with something worth looking at; trust my judgment on ordering.

**Commit (code):** N/A — no VCS in this workspace.

### What I did
- Added `components/organisms/MeetingPollPanel/` (`.tsx` + `.stories.tsx` + `index.ts`); exported from the organisms barrel.
- Props are DTO-shaped: `poll: MeetingPoll`, `tallies?`, `currentResponseId`, `styleSet`, `readOnly`, `draftName`/`draftComment`, and callbacks `onCellToggle` / `onNameChange` / `onCommentChange` / `onSubmit`.
- Composition: `Panel` → `Stack` → `KeyValueStrip` (location/organizer) + a deadline/response `Caption` + `MatrixGrid<ParticipantResponse>` (columns from options with a `DateTile` header, `CycleCell` cells, `RatioBadge`+★ footer) + a submit `Inline` (name/comment `TextInput` + primary `Button`).
- Stories: `ReadOnly`, `Respond` (stateful — cycles the "You" row and captures draft name/comment + submit), `Finalized`, `NoResponses`.
- Verified `typecheck` + `build-storybook`.

### Why
- The organism turns the abstract "engine + preset" story into a concrete product screen, and validates that the base engines compose cleanly with the existing layout/atoms (Panel/Stack/Inline/KeyValueStrip/Button/TextInput) without new primitives.

### What worked
- `typecheck` and `build-storybook` passed first try — the base engines and `scheduling` DTOs slotted in without adjustment, which is the signal the layering is holding.

### What didn't work
- No failures this step.

### What I learned
- Building the organism surfaced no missing atoms — the guide's proposed `AvatarStack` / `PollDeadlineBanner` molecules turned out to be optional; a `Caption` line and (future) `KeyValueStrip` cover the header cheaply. Worth trimming them from the must-build list.

### What was tricky to build
- Keeping the organism presentational while the "You" row is interactive. Cause: a poll response edit is stateful, but organisms must stay data-in/callbacks-out. Solution: the panel is fully controlled — it renders `poll.responses` and emits `onCellToggle({ responseId, optionId, state })`; the *story* owns the `useState` and rewrites the response. Matches how `TranscriptReaderPanel` etc. are structured.

### What warrants a second pair of eyes
- The editable-row affordance: only the `currentResponseId` row is interactive (others `readOnly`); confirm that reads clearly vs. greying-out non-editable rows.
- `MatrixGrid` horizontal scroll inside `Panel` at narrow widths (the panel caps at 560px in stories).

### What should be done in the future
- Engine adapters for `SegmentedBar` / `MonthGrid` / `TimeGrid` (+ `calendar/` types + IR stories).
- Sibling organisms: `PollResultsPanel` (SegmentedBar list + finalize), `BookingPagePanel` (MonthGrid + slots), `CalendarWeekPanel` (TimeGrid).
- A `MeetingPollPanel.widget.tsx` adapter once the React API is settled (React-first, so deferred).

### Code review instructions
- Start at `src/components/organisms/MeetingPollPanel/MeetingPollPanel.tsx` (props contract + the `MatrixGrid` composition), then the `Respond` story for the controlled-state pattern.
- Validate: `pnpm --dir packages/rag-evaluation-site typecheck` and `... build-storybook`; in `pnpm storybook` open *Organisms/MeetingPollPanel → Respond*.

### Technical details
- The panel defaults `styleSet` to `availabilityStyleSet` and pulls `AVAILABILITY_STATES` / `AVAILABILITY_GLYPHS` from the `scheduling` module — no palette duplication.
- `onCell` from `MatrixGrid` is adapted to the domain callback `onCellToggle({ responseId, optionId, state })`.

## Step 7: Split ir.ts into an ir/ module directory

`ir.ts` had grown to ~1255 lines and was about to grow further (booking/results
props). Split it into a directory of focused modules so future engine work
touches a small file, keeping every `from "../ir"` import working via folder
resolution.

### Prompt Context

**User prompt (verbatim):** "Work I think it makes sense to split ir.ts into separate modules, it'sg etting huge."

**Assistant interpretation:** Break the monolithic `ir.ts` into cohesive sub-modules, without breaking the ~30 files that import from it.

**Inferred user intent:** Keep the IR layer maintainable as the widget set grows; stop piling everything into one file.

**Commit (code):** N/A — no VCS in this workspace.

### What I did
- Converted `src/widgets/ir.ts` → `src/widgets/ir/` with: `core.ts` (JSON + node model + `RagWidgetType` + `BaseWidgetProps` + constructors), `actions.ts` (`ActionSpec` + template specs), `cells.ts` (`CellSpec`/`RowKeySpec`/`DataTableColumnSpec`), `engines.ts` (the scheduling engine specs — `StyleBySpec`, cycle/value cell, Matrix/Segmented/Month/Time props), `props.ts` (all other widget props + the `WidgetProps` union), and `index.ts` re-exporting all.
- Did the slice with a one-off script (`scratchpad/split-ir.mjs`) that cut by line ranges and prepended per-module imports; adjusted `../` → `../../` for the deeper path.
- Fixed the vite lib entry `src/widgets/ir.ts` → `src/widgets/ir/index.ts`.

### Why
- Folder-with-index keeps `import … from "../ir"` / `"../../../widgets/ir"` resolving unchanged, so it is a pure internal reorg with zero call-site edits. Type-only cross-module cycles (e.g. `core` ↔ `props` via `ComponentNode.props: WidgetProps`) are fine for `tsc`.

### What worked
- After two small fixes, `typecheck`, `build-storybook`, and the library `build` all pass.

### What didn't work
- `core.ts(141): error TS1005: '}' expected` — the slice range was off by one at EOF and clipped `isWidgetNode`'s closing brace. Appended `}`.
- `props.ts(45): error TS6196: 'CellSpec' is declared but never used` — `props` imported `CellSpec` it no longer uses (it moved to `cells.ts`). Removed it.
- `[UNRESOLVED_ENTRY] Cannot resolve entry module src/widgets/ir.ts` — the vite build had an explicit lib entry pointing at the old file. Repointed to `ir/index.ts`.

### What I learned
- The library build has hard-coded lib entries in `vite.config.ts`; any file move of an entry module must update it. `typecheck`/`build-storybook` alone would not have caught this — only the packaging `build` did.

### What was tricky to build
- Ordering the modules to avoid value-level cycles while allowing type-level ones. Cause: `ComponentNode` (in `core`) references `WidgetProps` (in `props`), and `BaseWidgetProps` (core) references `ActionSpec` (actions). Solution: keep all cross-references `import type` only, so they erase at compile time and the runtime module graph stays acyclic (constructors in `core` don't import from `props`).

### What warrants a second pair of eyes
- `props.ts` is still ~830 lines; a follow-up could split it by layer (foundation/layout vs molecules vs cms). Left as-is to bound this change.

### What should be done in the future
- Optionally split `props.ts` further; move `engines.ts` specs alongside future calendar-specific specs.

### Code review instructions
- Start at `src/widgets/ir/index.ts` then skim each module; confirm `defaultRegistry.ts` and `cellRenderers.tsx` still import cleanly.
- Validate: `pnpm --dir packages/rag-evaluation-site typecheck && … build-storybook && … build`.

### Technical details
- Split script retained at `ttmp/.../scratchpad/split-ir.mjs` (one-off; not part of the package).

## Step 8: Complete the engine adapters, presets, StyleBySpec, and Doodle/calendar/booking organisms

Finished the remaining scope in one pass: IR adapters for the other three
engines, the `StyleBySpec` consumption path, the calendar/results presets, IR
renderer stories for all of them, and four more organisms so the group-poll,
full-calendar, and 1:1-booking flows all have complete screens.

### Prompt Context

**User prompt (verbatim):** "continue, do it all."

**Assistant interpretation:** Complete everything outstanding — the SegmentedBar/MonthGrid/TimeGrid adapters + IR stories, wire `StyleBySpec` into a real renderer path, build the remaining presets, and the organizer/calendar/booking organisms.

**Inferred user intent:** Reach a feature-complete state across all three product flows rather than stopping at the poll.

**Commit (code):** N/A — no VCS in this workspace.

### What I did
- **Adapters + manifests** for `SegmentedBar` (`ui.dsl`), `MonthGrid`, `TimeGrid` (`time.dsl`); registered `segmentedBarWidget` in `uiWidgetRegistry` and a new `timeWidgetRegistry`, merged into `defaultWidgetRegistry`. Added `"time.dsl"` to `WidgetModule`.
- **StyleBySpec consumption**: new `ValueCellSpec` (`{ kind: "value" }`), a `resolveStyleByVars` util (`src/widgets/styleBy.ts`), and `colorBy` wiring in the `MatrixGrid` adapter so value cells tint from data → styleKey → `ContextVisualStyle`.
- **Presets**: `pollResults`, `monthCalendar`, `weekCalendar` in `presets/scheduling.ts`.
- **IR stories**: extended `WidgetRenderer.scheduling` (`PollResults`, `ColorByHeatmap`) and added `WidgetRenderer.calendar` (`MonthDensity`, `Week`).
- **Booking fixtures** (`sampleBookingType` / `sampleBookableDays` / `sampleBookableSlots`).
- **Organisms** (React-first) + stories + barrels: `PollResultsPanel` (SegmentedBar list + Pick/Finalize/Remind), `CalendarMonthPanel` (MonthGrid + selected-day event list), `CalendarWeekPanel` (TimeGrid), `BookingPagePanel` (SplitPane: host + MonthGrid | timezone + slot tiles + confirm).
- Verified `typecheck`, `build-storybook`, and the library `build`.

### Why
- Closing all three flows demonstrates the base/preset/contract layering scales past the poll: the same four engines back a poll grid, a results view, a month calendar, a week planner, and a booking page with no engine changes.

### What worked
- All three build gates pass. The organisms again composed from existing engines + layout/atoms with no new primitives.

### What didn't work
- `MatrixGrid.widget.tsx` type friction at the `MatrixRow` ↔ `JsonObject` boundary (same class as Step 5): `renderCell` needs `JsonObject`, and `WidgetActionContext.value` is a narrow union — fixed with `as unknown as JsonObject` casts and an explicit `value` cast in `onCell`.
- Left an unused `const timezone` in `BookingPagePanel`; removed after spotting it (would trip `noUnusedLocals`).

### What I learned
- `StyleBySpec` only becomes useful with a renderer that *asks* for it: `ValueCellSpec` + `colorBy` is the minimum pair that makes it observable (a heatmap), which is why declaring it in Step 5 without a consumer was a real gap.
- Namespacing (`calendar/…`) was not needed yet — no second DSL specializes an engine — so I kept bare types (`MonthGrid`, `TimeGrid`) and one `timeWidgetRegistry`.

### What was tricky to build
- `TimeGrid` block lane packing already existed; the adapter just had to pass serializable blocks through `ctx.renderValue` for labels while preserving `meta`. The subtlety was the `days` union (`string | { dayISO, header }`) — the adapter maps only the object form's header through `renderValue`, leaving bare strings untouched.

### What warrants a second pair of eyes
- `BookingPagePanel` `SplitPane ratio="sidebar"` behavior at narrow widths; slot `TileGrid` wrapping.
- `resolveStyleByVars` fallback chain (map → styleKey → fallbackStyleKey → fallbackStyle) and the `colorBy` wrapper's explicit `background: var(--ctx-fill)` (the css-var helper sets the var but the wrapper must consume it).
- `CalendarWeekPanel` sizing (`maxHeight: 480`) and header `DateTile` alignment.

### What should be done in the future
- `MeetingPollPanel` / `PollResultsPanel` / calendar / booking `.widget.tsx` adapters (organism-level IR) once the React APIs settle.
- `SchedulingStudioShell` (SidebarShell app frame) to tie the flows into one navigable app.
- Backend wiring in `packages/web` (the server actions/REST from the guide).

### Code review instructions
- Adapters: `src/components/molecules/{SegmentedBar,MonthGrid,TimeGrid}/*.widget.tsx`; registry in `src/widgets/defaultRegistry.ts`.
- StyleBySpec: `src/widgets/styleBy.ts` + the `colorBy` branch in `MatrixGrid.widget.tsx`; see `Widget IR/Renderer/Scheduling → ColorByHeatmap`.
- Organisms: `src/components/organisms/{PollResultsPanel,CalendarMonthPanel,CalendarWeekPanel,BookingPagePanel}/`.
- Validate: `pnpm --dir packages/rag-evaluation-site typecheck && … build-storybook && … build`; in `pnpm storybook` open the four new *Organisms* and the *Widget IR/Renderer/{Scheduling,Calendar}* stories.

### Technical details
- Registries now: `uiWidgetRegistry` (+SegmentedBar), `dataWidgetRegistry` (MatrixGrid), `timeWidgetRegistry` (MonthGrid, TimeGrid), merged in `defaultWidgetRegistry`.
- Full engine set is now DSL-drivable; presets cover poll / results / month / week; organisms cover all three product flows.

## Step 9: Correct the "no DSL runtime" error; write the DSL handoff guide

Two things: I corrected a wrong claim I had made (that the Goja DSL and backend
"don't exist in this repo") — they partly do, in `pkg/widgetdsl/` — and then
wrote a detailed, file/symbol-level handoff guide for the colleague who will wire
the widgets into the DSL.

### Prompt Context

**User prompt (verbatim):** "Ok, create a detailed guide of all the widgets you implemented and how they work and of their ir, we will hand off the implementation of the DSL to our goja expert colleague. Store the guide / report in the ticket, be really detailed, indicate code files and locations, symbols, etc..."

**Assistant interpretation:** Produce a thorough reference doc covering every implemented widget, its IR contract, emitted actions, and exactly where/how to wire it into the Go DSL, stored in the ticket.

**Inferred user intent:** A self-contained handoff so a Goja engineer can implement the DSL side without reading all the TS.

**Commit (code):** N/A — no VCS in this workspace.

### What I did
- Investigated the actual repo layout after the user's pointed questions ("which directory / cwd?" and "what about `pkg/widgetdsl/`?"). Found this repo is a Go project too (`cmd/ pkg/ internal/ web/ go.mod`) and that `pkg/widgetdsl/` is the go-go-goja Widget IR DSL runtime.
- Read `pkg/widgetdsl/{module.go,registrar.go,typescript.go}` and `v2/spec/validate.go` to understand the helper-map → `componentFactory` mechanism, `cellObject`/`actionObject`/`recipes`, TS-decl codegen, and validation scope.
- Wrote `reference/02-scheduling-and-calendar-widgets-implementation-and-ir-handoff-for-dsl-wiring.md`: file index (TS + Go), the new IR specs, a per-widget reference (props tables + action contexts + example IR JSON), the presets, and a step-by-step DSL wiring plan with exact `module.go` edit sites.

### Why
- The handoff must be precise enough to implement against; grounding every claim in real symbols/paths (and the `~Lxxx` line hints) makes it actionable.

### What worked
- The DSL mechanism turned out to be highly regular (helper map auto-generates the function), so the wiring plan reduces mostly to map entries + a `time.dsl` decision + cell builders + recipes.

### What didn't work / what I got wrong
- **My earlier answer was wrong.** I claimed "no backend, no Goja DSL in this repo" after grepping only `packages/`. The repo root has the full Go tree, and `pkg/widgetdsl/` is the DSL runtime. Root cause: I searched a subdirectory, not the repo root / session primary dir. Corrected in the changelog and the guide.

### What I learned
- `cwd` was `…/rag-evaluation-system` (the repo root) the whole time; my `ls packages` gave a false "TS-only" impression. Lesson: check the repo root and `go.mod` before asserting what a repo contains.
- Component `type` is open-ended in both the Go DSL and TS renderer (only node/action/cell *kinds* are validated), so exposing an engine needs no type-schema registration — just a helper-map entry.

### What was tricky to build
- Pinning the exact Go edit sites without over-reading: traced `install` → `componentFactory` → `buildComponent` and `cellObject`/`recipesObject` to cite `module.go` line ranges the colleague can jump to.

### What warrants a second pair of eyes
- The `time.dsl` reconciliation (add a Go module vs. relocate the two widgets) — documented as the one inconsistency my manifests introduced.
- Whether the preset recipes belong on new `schedule.dsl`/`calendar.dsl` modules or as recipes on existing modules.

### What should be done in the future
- Implement Part 5 of the guide in `pkg/widgetdsl`; survey `web/` for the widget-action dispatcher / REST handlers (not found in a quick grep).

### Code review instructions
- Read `reference/02-…-handoff-for-dsl-wiring.md`; cross-check its Part 1 file index against `src/widgets/ir/engines.ts` and `pkg/widgetdsl/module.go`.

### Technical details
- DSL mechanism: `runtime.install` (`module.go` ~L234) → `componentFactory` (~L565); cell builders in `cellObject` (~L272); composite builders as `recipes` (`recipesObject` ~L495, e.g. `masterDetailTableRecipe` ~L853). TS decls generated in `typescript.go` `TypeScriptModule`.

## Related

- `design-doc/01-calendar-and-scheduling-widgets-analysis-design-and-implementation-guide.md`
- `reference/02-scheduling-and-calendar-widgets-implementation-and-ir-handoff-for-dsl-wiring.md`
- `packages/rag-evaluation-site/GUIDELINES.md`
- `pkg/widgetdsl/module.go` (the DSL runtime to edit)
- IR modules at `packages/rag-evaluation-site/src/widgets/ir/`; engine adapters under `.../src/components/molecules/{MatrixGrid,SegmentedBar,MonthGrid,TimeGrid}`; presets at `.../src/widgets/presets/scheduling.ts`; organisms at `.../src/components/organisms/{MeetingPollPanel,PollResultsPanel,CalendarMonthPanel,CalendarWeekPanel,BookingPagePanel}`; domain module at `.../src/scheduling`

## Step 10: Widget DSL v3 phase tracker and working protocol

Started the Widget DSL v3 implementation track by turning the clean-break design into an explicit phase/task tracker. This step does not change runtime behavior; it establishes the control surface for the work: phases, task lists, validation gates, diary expectations, and commit boundaries.

The key decision in this step is to implement the new design in parallel as `widget.dsl` while leaving the existing modules available. That lets current `go-go-course` and other scripts keep running, while new fixtures and new pages can move to the cleaner TypeScript-facing builder/slot/fragment API.

### Prompt Context

**User prompt (verbatim):** "alright, create a detailed set of phases and for each phase a detailed list of tasks, so that we can keep track of where we are at precisely, and then implement step by step, keeping a detailed d iary as you work, and  committing at appropriate intervals."

**Assistant interpretation:** Create a precise phase/task tracker for the full Widget DSL v3 redesign, then proceed implementation incrementally with diary entries and focused commits.

**Inferred user intent:** The user wants the broad DSL redesign converted from architectural prose into an executable project plan with progress tracking, validation gates, and disciplined implementation hygiene.

**Commit (docs):** c99b32d33af17cd9aa86ff38cd999ff1aeb57533 — "Docs: record scheduling DSL design plan"

**Commit (code):** N/A — planning/documentation step; no runtime code changed.

### What I did
- Created `design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md`.
- Broke the Widget DSL v3 work into phases 0–11, from baseline inventory through `widget.dsl` skeleton, core builder kernel, UI/data/CMS/course/context/schedule namespaces, descriptor-generated declarations, go-go-course fixtures, and integration guidance.
- Added per-phase tasks, acceptance criteria, validation commands, and suggested commit messages.
- Related the tracker to the current runtime files and the go-go-course pages that should become real v3 fixtures.
- Added high-level docmgr tasks for the v3 phase sequence.

### Why
- The prior design documents explained the target API, but implementation needed a concrete progress tracker.
- The work is large enough that ambiguous “next steps” would quickly become hard to review.
- A phase tracker gives future sessions a stable checklist and makes it easier to stop after coherent commits.

### What worked
- The existing docmgr ticket already had the right context and design documents, so the tracker could reference the current files and decisions directly.
- The plan can preserve current scripts while implementing `widget.dsl` as a parallel v3 API.
- The phase boundaries map cleanly onto reviewable commits.

### What didn't work
- The repository working tree already contains many untracked and modified files from prior widget/ticket work. This means commits must stage specific files only; broad `git add .` would be unsafe.

### What I learned
- The right implementation unit is not “replace ui.dsl” but “add `widget.dsl` root module, then port real pages as fixtures.”
- `go-go-course` should be treated as the proving ground because it uses UI, data, CMS, course, and context modules together.
- The current `data.v2.dsl` builder mechanics are useful implementation material but should not remain a separate public authoring dialect in v3.

### What was tricky to build
- The tracker had to balance two truths: the new API should not be constrained by backward compatibility, but the old modules should remain available while v3 matures. The resolution is a parallel module with shared internals and fixture-based porting.
- The task list needed enough detail to guide implementation without pretending that every method name is already final. I kept phase-level acceptance criteria stable and left details like exact descriptor structs to the implementation phases.

### What warrants a second pair of eyes
- Phase boundaries: especially whether descriptor-driven TypeScript generation should move earlier than Phase 9.
- The first fixture choice: `go-go-course` DSL examples are broad, while the Course CMS page is more representative and more complex.
- Whether high-level docmgr phase tasks should be further decomposed into one task per namespace once Phase 1 begins.

### What should be done in the future
- Finish Phase 0 by generating the current DSL export inventory.
- Commit the tracker/diary/task setup as its own focused documentation commit.
- Begin Phase 1 with only the `widget.dsl` skeleton and `raw.component` escape hatch.

### Code review instructions
- Start with `design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md`.
- Check that each phase has tasks, acceptance criteria, and a validation/commit boundary.
- Verify the tracker does not require deleting old modules before `widget.dsl` is proven.
- Validate with `docmgr doctor --ticket RAGEVAL-SCHEDULE-WIDGETS --stale-after 30`.

### Technical details
- Phase 0 is planning/inventory only.
- Phase 1 is the first runtime code phase: add `widget.dsl` in parallel and prove old modules still load.
- The commit protocol requires staging specific files because the working tree contains many unrelated untracked files.

## Step 11: Phase 0 export inventory

Completed the Phase 0 inventory slice by generating a current export map for the existing split Widget DSL modules. This is the bridge between the redesign and implementation: it names what old capabilities exist today, classifies helpers as generic foundation helpers, engine helpers, domain helpers, or domain-exported generic aliases, and makes the first v3 port target explicit.

No runtime behavior changed. The only executable artifact is a small ticket-local inventory generator so the current export table can be refreshed as the old modules evolve.

### Prompt Context

**User prompt (verbatim):** (same as Step 10)

**Assistant interpretation:** Continue from the tracker by executing Phase 0: produce a precise current DSL surface inventory before starting runtime changes.

**Inferred user intent:** The user wants implementation to proceed step by step with exact progress tracking, not jump straight into code without a baseline.

**Commit (docs):** 4ff1ae57d55f478addb71679718cf6b4e19bbb03 — "Docs: complete Widget DSL v3 inventory"

**Commit (code):** N/A — documentation/script-only Phase 0 slice; no runtime code changed.

### What I did
- Added `scripts/01-widget-dsl-export-inventory.py` under the ticket workspace.
- Generated `reference/04-widget-dsl-current-export-inventory.md` from `pkg/widgetdsl/module.go`.
- Inventoried helper counts, recipes, component type mappings, and helper classifications for `ui.dsl`, `data.dsl`, `context_window.dsl`, `course.dsl`, `cms.dsl`, and noted `data.v2.dsl` as separately installed typed builders.
- Updated `design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md` to mark Phase 0 complete.
- Checked docmgr task `p4o2`.

### Why
- A clean redesign still needs to preserve current capability. The inventory shows what capability exists without forcing v3 to preserve current names.
- Classification helps decide whether a helper should become `widget.dsl.ui`, a typed engine builder, a domain view, or an implementation detail.
- The generator avoids stale hand-written lists if the current module maps change before v3 lands.

### What worked
- The helper maps and module specs in `module.go` are simple enough for a lightweight script to parse reliably for planning purposes.
- The inventory confirms the main design hypothesis: generic UI helpers are mixed into domain modules as aliases, while domain recipes and engine helpers need a cleaner v3 home.

### What didn't work
- My first run of the script used the wrong repository-root parent count and failed with:
  `FileNotFoundError: [Errno 2] No such file or directory: '/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/pkg/widgetdsl/module.go'`
- I fixed it by changing the script root from `parents[7]` to `parents[6]`.

### What I learned
- Current `ui.dsl` has 41 helper-map helpers; `cms.dsl` and `course.dsl` still export generic aliases that v3 should centralize under `ui`.
- `context_window.dsl` contains both real engines (`contextTreemap`, `contextStripDiagram`, etc.) and task-level panels (`transcriptWorkspacePanel`), which supports the v3 split between engine helpers and domain views.
- `data.v2.dsl` remains the strongest implementation precedent for typed builders, but the public v3 namespace should be `widget.dsl.data`.

### What was tricky to build
- The generator needed to be useful without becoming a fragile Go parser. I kept it deliberately narrow: parse simple map literals and recipe arrays, then apply a small hand-maintained classification table.
- Classification is partly subjective. For example, `contextUploadDropArea` is domain-exported today but semantically generic enough to become a `ui` upload helper in v3. The inventory labels current state and v3 implications rather than claiming final truth.

### What warrants a second pair of eyes
- Helper classification, especially context-window engines vs domain panels.
- Whether `uploadDropArea` should be treated as purely generic `ui` in v3 or remain available through context/CMS convenience builders.
- Whether the Phase 1 skeleton should include descriptor structs immediately or defer descriptors until Phase 2.

### What should be done in the future
- Start Phase 1 by adding `widget.dsl` as a parallel module skeleton.
- Keep `dsl-examples.js` as the first real fixture target.
- Use the inventory as a checklist when deciding which old capabilities need v3 equivalents.

### Code review instructions
- Review `scripts/01-widget-dsl-export-inventory.py` first; confirm it only reads `pkg/widgetdsl/module.go` and writes the ticket reference doc.
- Review `reference/04-widget-dsl-current-export-inventory.md` for generated output and classification sanity.
- Review `design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md` to confirm Phase 0 status is now complete.
- Validate with `python3 ttmp/.../scripts/01-widget-dsl-export-inventory.py` and `docmgr doctor --ticket RAGEVAL-SCHEDULE-WIDGETS --stale-after 30`.

### Technical details
- The script uses `Path(__file__).resolve().parents[6]` to find the repository root from the ticket script directory.
- It parses `map[string]string` helper maps and `recipes: []string{...}` arrays from `module.go`.
- It writes a full Markdown reference document with docmgr frontmatter.

## Step 12: Phase 1 parallel `widget.dsl` skeleton

Implemented the first runtime slice of Widget DSL v3: a new `widget.dsl` module that exists beside the current split modules. This is deliberately small. It does not implement the final builder API yet; it proves the module can be registered, selected by the xgoja provider, required in Goja, and used for the initial explicit raw escape hatch plus shared `act` and `bind` namespaces.

The important behavior is coexistence. `require("widget.dsl")` now succeeds, and `require("ui.dsl")`, `require("data.dsl")`, and the other current modules still work. This lets v3 grow through fixtures without breaking existing scripts.

### Prompt Context

**User prompt (verbatim):** (same as Step 10)

**Assistant interpretation:** Continue step-by-step implementation after planning and Phase 0 inventory, starting with the smallest code slice: the parallel `widget.dsl` module skeleton.

**Inferred user intent:** The user wants the redesign to become executable incrementally, with commits at safe boundaries and a detailed diary of what changed and what failed.

**Commit (code):** pending — Phase 1 code slice to be committed after this diary update.

### What I did
- Added `WidgetV3ModuleName = "widget.dsl"` and registered a new module spec.
- Added a special install path for `widget.dsl` that exports root namespaces: `raw`, `act`, `bind`, `page`, `ui`, `data`, `cms`, `course`, `context`, `schedule`, `time`, and `style`.
- Implemented `raw.text`, `raw.element`, `raw.component`, and `raw.fragment` using the existing low-level constructors.
- Reused the existing action helper object as `widget.act`.
- Added an initial binding helper object as `widget.bind` with `field`, `path`, `map`, `template`, `context`, and `const` constructors.
- Added a TypeScript declaration stub for `widget.dsl` root namespaces.
- Exposed `widget.dsl` from the `rag-widget-site` xgoja provider.
- Added tests for runtime export shape, raw component output, binding/action output, TypeScript declaration fragments, and provider exposure.
- Updated the Phase 1 tracker status to complete.

### Why
- The full redesign needs a safe parallel entry point before any builder or domain work can land.
- `raw.component` gives v3 an immediate escape hatch and a simple testable output path.
- `act` and `bind` establish the two serializable helper namespaces that later builder phases will reuse.

### What worked
- The existing `moduleSpecs` registration loop made adding a new module straightforward.
- The existing low-level constructors and `actionObject` could be reused directly.
- Provider registration followed the same pattern as the current split modules.
- Targeted tests passed after updating old assumptions that `widget.dsl` was absent.

### What didn't work
- The first test run failed because an existing test still expected `widget.dsl` to be absent from the engine registrar:
  `--- FAIL: TestEngineRegistrarRegistersSplitModulesOnly`
  `old bucket modules should be absent from engine registrar, got map[string]interface {}{..., "widget":true}`
- This was a correct failure: `widget.dsl` is no longer an old bucket module; it is the new parallel v3 module. I updated the test expectation to require `widget.dsl` and still reject `rag.dsl`.

### What I learned
- The old test name and assertion encoded the historical decision that `widget.dsl` should not exist. Phase 1 intentionally supersedes that decision while preserving the absence of `rag.dsl`.
- It is useful to make `widget.dsl` special-cased in `install` for now. Later phases can replace the manual namespace setup with descriptor-driven installation.
- The first TypeScript declaration can be a stub for namespace shape; precise builders belong in later phases.

### What was tricky to build
- `widget.dsl` uses a different shape from the old split modules: old modules export helper functions at the top level, while v3 exports namespaces. The runtime install path therefore needed to branch early instead of flowing through the existing helper-map installation path.
- The `bind.const` helper uses a property name that is a JavaScript keyword. It is valid as an object property and can be called as `bind.const(...)`, but it warrants attention in TypeScript fixture tests later.

### What warrants a second pair of eyes
- Whether `raw` should expose `element` and `fragment` long-term or whether those should move to `ui` once the UI namespace exists.
- Whether `bind.path` and `bind.context` should both use `path`, or whether `bind.field` should keep `field` for closer compatibility with current IR accessors.
- Whether the v3 TypeScript stub should remain hand-written until Phase 9 or move to descriptors earlier.

### What should be done in the future
- Start Phase 2 by adding real page/node specs, builder callback helpers, `.use(fragment)`, and slot invocation.
- Add TypeScript fixture compilation for `widget.dsl` once builder APIs exist.
- Keep old modules working until v3 fixtures cover real go-go-course pages.

### Code review instructions
- Start with `pkg/widgetdsl/module.go`: review `WidgetV3ModuleName`, `installWidgetV3`, `rawObject`, and `bindingObject`.
- Then review `pkg/widgetdsl/module_test.go` for coexistence expectations.
- Review `pkg/widgetdsl/typescript.go` and `pkg/widgetdsl/typescript_test.go` for the initial declaration stub.
- Review `pkg/xgoja/providers/widgetsite/provider.go` and provider tests for xgoja exposure.
- Validate with `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1`.

### Technical details
- Validation command run:
  `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1`
- Final result:
  `ok github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl`
  `ok github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl/v2/spec`
  `ok github.com/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite`
