---
Title: Code Review and Design Assessment for Scheduling Widgets
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
    - Path: repo://packages/rag-evaluation-site/package.json
      Note: Reviewed public package export implications for scheduling and presets
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/MatrixGrid/MatrixGrid.tsx
      Note: Reviewed flagship generic grid engine and MatrixCellPayload contract
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/MatrixGrid/MatrixGrid.widget.tsx
      Note: Reviewed IR adapter value/cell/footer/action interpretation
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/TimeGrid/TimeGrid.tsx
      Note: Reviewed lane packing and all-day-event contract issue
    - Path: repo://packages/rag-evaluation-site/src/components/organisms/MeetingPollPanel/MeetingPollPanel.tsx
      Note: Reviewed readOnly/editable-row behavior and organism composition
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/engines.ts
      Note: Reviewed new serializable engine IR contracts
    - Path: repo://packages/rag-evaluation-site/src/widgets/presets/scheduling.ts
      Note: Reviewed scheduling/calendar preset node builders
ExternalSources: []
Summary: Quality assessment and intern-facing code review of the scheduling/calendar widget implementation, including design decisions, validation evidence, findings, and merge-readiness guidance.
LastUpdated: 2026-07-06T21:02:40.499848647-04:00
WhatFor: Use this as the reviewer handoff for the scheduling widget work and as a teaching document for interns learning the Widget IR/design-system architecture.
WhenToUse: Read before merging, extending, or wiring the scheduling/calendar widgets into pkg/widgetdsl or a production app surface.
---


# Code Review and Design Assessment for Scheduling Widgets

## Executive summary

The colleague's work is **substantially strong**. They delivered a coherent scheduling/calendar subsystem rather than a pile of one-off Doodle components: generic engines (`MatrixGrid`, `SegmentedBar`, `MonthGrid`, `TimeGrid`), small atoms (`CycleCell`, `DateTile`, `RatioBadge`), serializable Widget IR adapters, scheduling/calendar presets, and React-first organisms for the main product flows. The implementation also includes unusually good handoff documentation and a candid diary that records mistakes, fixes, and validation gates.

My recommendation is: **keep the architecture and most implementation choices, but do not merge/release without a short cleanup pass.** The code passes the important compile/build gates, but there are a few behavior and API-surface issues that an intern should learn to catch: `MeetingPollPanel.readOnly` does not actually make cells read-only, `TimeGrid` silently drops `allDay` events, and the new scheduling domain/preset APIs are built into `dist/` but are not exported from the package root or `package.json` subpath exports. Add focused tests for the pure algorithms and IR adapters before treating this as production-ready.

## Review scope and evidence read

I reviewed the ticket artifacts and the implementation described by the diary:

- `reference/01-implementation-diary.md` — chronological work log, including failures and self-correction.
- `reference/02-scheduling-and-calendar-widgets-implementation-and-ir-handoff-for-dsl-wiring.md` — Goja/DSL handoff guide.
- `design-doc/01-calendar-and-scheduling-widgets-analysis-design-and-implementation-guide.md` — original design and intern guide.
- The uncommitted TypeScript/React implementation under `packages/rag-evaluation-site/src/`.
- The Go DSL package enough to validate that the colleague's handoff correctly identifies `pkg/widgetdsl` as future work.

Important state note: the work is currently uncommitted in this workspace. `git status --short` shows modified existing files plus many untracked widget, scheduling, preset, and ticket files. Treat this as a review of a large working tree change, not as a review of a clean commit series.

## Validation performed during this review

All validation commands below passed in `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system`:

```bash
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build
pnpm --dir packages/rag-evaluation-site build-storybook
go test ./pkg/widgetdsl/...
pnpm --dir packages/rag-evaluation-site pack:smoke
pnpm --dir packages/rag-evaluation-site consumer:smoke
```

The green gates matter. They prove the IR split did not break TypeScript, the library package still builds, Storybook still compiles, and the Go widget DSL tests still pass even though the DSL wiring has not been implemented yet.

## Overall quality rating

**Rating: strong prototype / near mergeable design-system slice, not yet production-hard.**

- **Architecture:** A-/A. The generic-engine + IR-adapter + preset split is the right abstraction boundary for this codebase.
- **Implementation correctness:** B+. The primary flows compile and render, but there are a few real behavior bugs and missing edge-case tests.
- **Documentation/handoff quality:** A. The diary and DSL handoff are specific, actionable, and unusually honest about mistakes.
- **Test/validation quality:** B-. Build gates are excellent; focused unit tests are missing for the riskiest pure functions and adapter contracts.
- **API/package readiness:** B-. Internal exports are mostly wired, but public package export decisions need one pass.

## What they did well

### 1. They kept the central abstraction clean

The best design decision is the split between **space/time engines** and **domain meaning**. `MatrixGrid` owns rows, columns, sticky headers, footer slots, explicit-cell mode, and the `onAction` seam, but it does not know what availability means (`MatrixGrid.tsx:12-57`, `MatrixGrid.tsx:64-167`). That is exactly the right reusable component boundary.

The same pattern appears in the other engines:

- `MonthGrid` owns month geometry, UTC-safe date arithmetic, marker rendering, and selection callbacks (`MonthGrid.tsx:28-46`, `MonthGrid.tsx:84-183`).
- `TimeGrid` owns hour geometry and overlapping block lane-packing (`TimeGrid.tsx:82-134`, `TimeGrid.tsx:136-262`).
- `SegmentedBar` owns proportional segment display and palette lookup (`SegmentedBar.tsx:19-34`, `SegmentedBar.tsx:36-126`).

For an intern: this is the main lesson. A reusable design-system component should own mechanics, not product nouns.

### 2. The cell/action contracts are thoughtful

`MatrixGrid`'s `MatrixCellPayload` is a good seam: it gives a cell the row, column, resolved value, selection/editability flags, and a grid-owned `onAction` callback (`MatrixGrid.tsx:18-29`). The `onAction` override for `value` is a small but important detail: a `CycleCell` can emit the *next* value even though the grid only knows the current resolved value (`MatrixGrid.tsx:129-134`).

On the IR side, the adapter is similarly disciplined. `MatrixGrid.widget.tsx` maps serializable specs back into real React closures: value accessors (`MatrixGrid.widget.tsx:26-42`), cycle/value/data-table cells (`MatrixGrid.widget.tsx:44-97`), row headers, explicit cell matrices, footers, and action dispatch (`MatrixGrid.widget.tsx:99-154`). This is the exact adapter responsibility described in the original design.

### 3. The IR vocabulary is small and mostly orthogonal

The new engine IR types are compact and reusable (`widgets/ir/engines.ts:14-142`):

- `StyleBySpec` is a data representation of a color function.
- `CycleCellSpec` and `ValueCellSpec` extend the existing `CellSpec` idea only where needed.
- `MatrixValueSpec` handles the two real value-access patterns: map lookup and template interpolation.
- `MonthGridWidgetProps` and `TimeGridWidgetProps` stay close to their React engine props.

This is good IR design: serializable, boring, and close to the runtime adapter shape.

### 4. The IR split was worth doing

Moving the old monolithic `widgets/ir.ts` into `widgets/ir/{core,actions,cells,engines,props}.ts` improves maintainability. `core.ts` now owns JSON/node primitives and the `RagWidgetType` union (`core.ts:4-141`), while `props.ts` imports the engine prop types and includes them in `WidgetProps` (`props.ts:47-52`, `props.ts:608-690`).

The colleague also caught the build-system implication and updated the Vite lib entry. That was an important packaging detail, not just a TypeScript cleanup.

### 5. The stories cover the right review surfaces

The Storybook set is broad and practical:

- Atoms: all states, sizes, invalid/zero cases.
- Engines: default, selected, empty, custom renderer, interactive, bounds/leap-year cases.
- Organisms: read-only/respond/finalized/no-response/pending/interactive cases.
- WidgetRenderer stories: scheduling and calendar IR trees.

The stories are especially useful because this is UI work and several decisions are visual/interaction decisions that unit tests will not catch.

### 6. The diary and handoff document are high quality

The colleague recorded exact TypeScript errors, why fixes were chosen, and one major self-correction: they initially missed that the Goja DSL runtime exists in `pkg/widgetdsl`, then corrected the handoff (`reference/01-implementation-diary.md`, Step 9). That is good engineering hygiene. It gives the next engineer a truthful map instead of a polished but misleading story.

## Review findings to fix or consciously accept

### Finding 1 — `MeetingPollPanel.readOnly` does not disable editable cells

**Severity:** High before production use; Medium for Storybook-only prototype.

`MeetingPollPanel` computes `editing = currentResponseId != null && !readOnly` and uses that only to show/hide the submit row (`MeetingPollPanel.tsx:76-79`, `MeetingPollPanel.tsx:155-173`). But the `MatrixGrid` still receives `editableRowKey={currentResponseId}` unconditionally (`MeetingPollPanel.tsx:120-122`), and each `CycleCell` is read-only only when `!p.editable` (`MeetingPollPanel.tsx:122-131`).

That means a caller can pass `readOnly={true}` and `currentResponseId="you"`; the submit controls disappear, but the cells remain enabled and can still emit `onCellToggle` (`MeetingPollPanel.tsx:133-139`).

**Recommended fix:**

```tsx
const canEdit = currentResponseId != null && !readOnly;

<MatrixGrid
  editableRowKey={canEdit ? currentResponseId : undefined}
  onCell={canEdit ? handleCell : undefined}
/>
```

Also add a Storybook story or test that sets both `readOnly` and `currentResponseId` so this does not regress.

### Finding 2 — `TimeGrid` accepts `allDay` data but silently drops it

**Severity:** Medium/High.

`TimeGridBlock` includes `allDay?: boolean` (`TimeGrid.tsx:5-15`), `TimeGridBlockSpec` includes `allDay?: boolean` (`widgets/ir/engines.ts:115-124`), `CalendarEvent` includes `allDay?: boolean` (`scheduling/types.ts:64-74`), and `CalendarWeekPanel.toBlocks` passes `allDay` through (`CalendarWeekPanel.tsx:21-30`). But `TimeGrid` filters all-day blocks out of the timed layout and never renders them elsewhere (`TimeGrid.tsx:191-195`).

This is worse than not supporting all-day events because the type says they are accepted. Data can disappear without warning.

**Recommended fix options:**

1. Add an all-day row above the timed columns and render `allDay` blocks there.
2. If all-day support is out of scope, remove `allDay` from `TimeGridBlock`, `TimeGridBlockSpec`, and the week preset for now, or document that all-day blocks are intentionally ignored.

For an intern: never let the type contract imply support that the renderer silently discards.

### Finding 3 — Public package exports for scheduling/presets are incomplete

**Severity:** Medium.

The build output includes `dist/scheduling/*` and `dist/widgets/presets/scheduling.d.ts`, and the new organisms' public prop types refer to scheduling DTOs. However, the package root currently exports `cms`, `components`, `context`, `hooks/useWidgetPage`, and `widgets`, but not `scheduling` (`src/index.ts`). The `widgets/index.ts` barrel exports actions/cell renderers/default registry/IR/registry/renderer, but not `widgets/presets` or `widgets/styleBy`.

`dist/package.json` also has only these public subpath exports: `.`, `./ir`, `./app`, `./styles.css`, and `./theme.css`. So consumers can use the rendered components from the root, but the new scheduling DTOs and presets are not intentionally public API.

**Recommended fix:** decide one of these explicitly.

- If scheduling is public: add `export * from "./scheduling"` to `src/index.ts`, add a `widgets/presets/index.ts`, export presets from `widgets/index.ts`, and add package `exports` for `./scheduling` and maybe `./widgets/presets` if subpath imports are expected.
- If scheduling is internal/demo-only: document that and avoid exposing organisms whose public prop types depend on scheduling DTOs as stable API.

The current state passes `consumer:smoke`, but the smoke test does not prove the scheduling/preset public API shape.

### Finding 4 — The risky pure logic needs focused tests

**Severity:** Medium.

The green `typecheck`, `build`, and `build-storybook` gates are valuable, but they do not lock down algorithmic and contract behavior. The riskiest pure logic is currently exercised only indirectly through stories:

- `packColumn` overlap clustering and lane assignment (`TimeGrid.tsx:82-134`).
- Month matrix generation, UTC/date-bound handling, and min/max disabled logic (`MonthGrid.tsx:64-130`).
- `resolveStyleByVars` fallback chain (`styleBy.ts:11-25`).
- `MatrixGrid.widget.tsx` value resolution and footer meta rendering (`MatrixGrid.widget.tsx:26-42`, `MatrixGrid.widget.tsx:131-137`).
- Action contexts emitted by the four adapters (`MatrixGrid.widget.tsx:140-149`, `SegmentedBar.widget.tsx:25-35`, `MonthGrid.widget.tsx:33-52`, `TimeGrid.widget.tsx:23-43`).

**Recommended fix:** add small unit tests for these helpers. If private helpers need testing, either export a test-only pure helper from a nearby `*.logic.ts` file or test via the public component/adapters with React Testing Library.

### Finding 5 — Accessibility needs a second pass

**Severity:** Low/Medium.

The code makes a good-faith accessibility pass (`aria-label`, `aria-pressed`, semantic tables), but several details need review before production:

- `SegmentedBar` renders interactive segments as empty buttons with only `title` for labeling (`SegmentedBar.tsx:65-75`). Add `aria-label` based on label/style key/count.
- `MonthGrid` uses `role="grid"` and default day buttons with `role="gridcell"` (`MonthGrid.tsx:170-178`, `MonthGrid.tsx:195-207`). Confirm the final structure with an accessibility checker; grid patterns usually need rows and keyboard navigation semantics.
- `TimeGrid` has many hour-slot buttons (`TimeGrid.tsx:203-214`). That is acceptable for a prototype, but a production calendar may need keyboard shortcuts, current-time constraints, and more descriptive slot labels.

### Finding 6 — Visual/date polish remains prototype-level in a few spots

**Severity:** Low.

The major visual building blocks are good, but a few details should be reviewed in Storybook:

- `MeetingPollPanel.columnHeader` renders a `DateTile` and then renders the day number again as a separate `Text` node (`MeetingPollPanel.tsx:46-58`). That may be redundant visually.
- Date/time formatting is intentionally simple string slicing in multiple places (`MeetingPollPanel.tsx:18-23`, `PollResultsPanel.tsx:9-14`, `presets/scheduling.ts:22-28`). This is fine for deterministic stories, but production app code should centralize formatting and timezone decisions.
- `CalendarMonthPanel` uses inline styles for event swatches/buttons (`CalendarMonthPanel.tsx:68-92`). That is acceptable for a first slice, but CSS modules/tokens would make the anatomy easier to theme.

## Design decision review

### Decision: Generic engines instead of domain-specific Doodle widgets

- **Context:** The request was to build calendar/scheduling UI "using base widgets" and keep it intern-teachable.
- **Options considered:** One-off `PollGrid`/`BookingCalendar` components, or generic engines with domain presets.
- **Decision:** Generic engines plus domain presets.
- **Rationale:** The implementation proves reuse across poll grid, results, month calendar, week planner, and booking page without engine-specific domain branches.
- **Consequences:** Slightly more adapter/preset code up front; much better reuse and DSL fit.
- **Status:** Accepted.

### Decision: React-first organisms, IR adapters only for engines

- **Context:** The original guide says React first, Widget IR later.
- **Options considered:** Expose every organism to IR immediately, or expose only stable low-level engines first.
- **Decision:** Engine adapters now; organism adapters later.
- **Rationale:** The organisms are useful review surfaces, but their APIs may still move. Engine contracts are more stable and are the needed DSL primitive layer.
- **Consequences:** `MeetingPollPanel`, `PollResultsPanel`, `CalendarMonthPanel`, `CalendarWeekPanel`, and `BookingPagePanel` are not yet DSL-drivable as single nodes.
- **Status:** Accepted with follow-up.

### Decision: Add `time.dsl` on the TypeScript side before Go DSL wiring

- **Context:** Month/week calendar engines need a logical module; `pkg/widgetdsl` does not yet expose `time.dsl`.
- **Options considered:** Put `monthGrid`/`timeGrid` under `ui.dsl` or `data.dsl`, or introduce `time.dsl`.
- **Decision:** Introduce `time.dsl` in TS widget metadata and registry (`registry.ts:5-11`, `defaultRegistry.ts:121`, manifests for `MonthGrid`/`TimeGrid`).
- **Rationale:** Time/calendar engines are generic but distinct from data tables and general UI atoms.
- **Consequences:** Go DSL wiring must add the same module or the manifests and runtime will disagree. The handoff guide correctly calls this out.
- **Status:** Accepted on TS side; pending Go implementation.

### Decision: Keep bare engine `type` strings for now

- **Context:** The design guide discusses namespaced specialized types to avoid duplicate registry keys.
- **Options considered:** Bare `MatrixGrid`/`MonthGrid`/`TimeGrid` vs. namespaced `calendar/TimeGrid` etc.
- **Decision:** Keep bare engine type strings for generic engines (`core.ts:39-43`, `defaultRegistry.ts:119-121`).
- **Rationale:** These are base engines, not schedule/calendar specializations, and `createWidgetRegistry` will catch accidental duplicate types (`registry.ts:40-47`).
- **Consequences:** Future domain-specific adapters should use namespaced types if they wrap/specialize the same engine.
- **Status:** Accepted.

### Decision: Use `StyleBySpec` as the defunctionalized color function

- **Context:** The DSL cannot ship closures, but cells need value-driven styling.
- **Options considered:** Raw colors in props, hardcoded availability colors, or reusable `ContextStyleSet` lookup.
- **Decision:** Add `StyleBySpec` and interpret it with `resolveStyleByVars` (`engines.ts:14-21`, `styleBy.ts:11-25`).
- **Rationale:** It reuses the existing context palette contract and keeps the IR serializable.
- **Consequences:** Currently only `MatrixGrid` value cells consume it (`MatrixGrid.widget.tsx:70-90`). Other engines can adopt it later if needed.
- **Status:** Accepted as a seed abstraction.

## Merge-readiness checklist

Before merging this work, I would require:

1. Fix `MeetingPollPanel.readOnly` so it disables cells and action emission.
2. Decide and implement the `TimeGrid.allDay` contract: render all-day events or remove the prop until supported.
3. Decide the public API for `scheduling` and `widgets/presets`; export intentionally or document as internal.
4. Add unit tests for `packColumn`, month generation/bounds, `resolveStyleByVars`, and at least one adapter action context.
5. Run the same gates again: `typecheck`, `build`, `build-storybook`, `pack:smoke`, `consumer:smoke`, and `go test ./pkg/widgetdsl/...`.
6. Review Storybook visually on narrow widths for `MatrixGrid`, `TimeGrid`, and `BookingPagePanel`.

## Intern learning notes

### What to copy

- Copy the **engine/adapter/preset separation**. It is the most important design lesson in this ticket.
- Copy the `MatrixCellPayload` style of stable render payloads. A render-prop seam is how a reusable engine stays domain-blind.
- Copy the diary habit: record exact compiler errors and why the fix was chosen.
- Copy the validation discipline: typecheck, Storybook build, package build, package smoke, and relevant Go tests.

### What not to copy blindly

- Do not let optional props imply unsupported behavior (`allDay` is the example here).
- Do not assume a `readOnly` flag is wired just because the prop exists. Trace it to the actual interactivity switch.
- Do not rely on Storybook alone for pure logic; add focused tests for algorithms and IR interpreter behavior.
- Do not add public DTOs/presets without deciding the package export story.

### How to review this kind of work next time

Start with the contracts, not the visuals:

1. Read the exported prop types and IR types.
2. Confirm every typed feature is implemented or explicitly documented as unsupported.
3. Trace one user action from component callback → adapter dispatch context → `ActionSpec` payload.
4. Trace one rendered value from domain DTO → preset → IR node → adapter → React component.
5. Only then review Storybook visuals and CSS polish.

## Follow-up plan

### Phase 1 — Safety cleanup

- Patch `MeetingPollPanel.readOnly`.
- Patch or remove `TimeGrid.allDay` support.
- Add accessibility labels to `SegmentedBar` interactive segment buttons.
- Add a `readOnly + currentResponseId` story.

### Phase 2 — Contract tests

- Test `packColumn` overlap cases: non-overlap, exact-touch, nested overlap, clamped out-of-range, zero-length minimum.
- Test month generation: leap February, Sunday/Monday starts, min/max bounds, invalid month.
- Test `StyleBySpec`: direct key, mapped key, fallback key, fallback style, missing style.
- Test adapter action contexts for `MatrixGrid`, `MonthGrid`, and `TimeGrid`.

### Phase 3 — API/publication decisions

- Export `scheduling` and presets if the package should make these first-class.
- Add package subpath exports if consumers should import `@go-go-golems/rag-evaluation-site/scheduling` or `.../widgets/presets`.
- Extend `consumer:smoke` with a scheduling component/preset import.

### Phase 4 — DSL handoff

- Implement `reference/02` Part 5 in `pkg/widgetdsl`: helpers, `time.dsl`, `cell.cycle`, `cell.value`, optional `styleBy`, and recipes.
- Run `go test ./pkg/widgetdsl/...` after each slice.
- Keep the Go DSL TypeScript declarations in parity with the TS widget package.

## Final assessment

This is good work. The colleague made the right architectural bet and carried it far enough that the team can now review real screens, real IR trees, and a concrete Goja handoff. The remaining issues are the sort that appear when a large prototype crosses the line toward a reusable library: edge-case behavior, public API boundaries, accessibility details, and tests for pure logic.

For an intern, the most important lesson is not "write more code"; it is **make every abstraction earn its place**. Here, `MatrixGrid`, `MonthGrid`, `TimeGrid`, `SegmentedBar`, `StyleBySpec`, and the presets mostly earn their place because each separates a reusable mechanical concern from scheduling-specific meaning. The cleanup pass should preserve that structure while tightening the contracts.

## Related

- `reference/01-implementation-diary.md` — chronological build log.
- `reference/02-scheduling-and-calendar-widgets-implementation-and-ir-handoff-for-dsl-wiring.md` — DSL wiring handoff.
- `design-doc/01-calendar-and-scheduling-widgets-analysis-design-and-implementation-guide.md` — original intern guide.
