# Changelog

## 2026-07-06

- Initial workspace created


## 2026-07-06

Created ticket + intern-facing analysis/design/implementation guide for calendar/scheduling widgets built on generic base engines (MatrixGrid, SegmentedBar, MonthGrid, TimeGrid, CycleCell) with a three-layer architecture (base component -> IR adapter -> DSL preset) and a swappable cell contract. Started implementation diary.


## 2026-07-06

Uploaded intern guide to reMarkable. Implemented first base-widget slice (React-first, no IR adapters yet): atoms DateTile/RatioBadge/CycleCell and molecule engines SegmentedBar/MatrixGrid, each with Storybook stories. MatrixGrid demonstrates both injection modes (renderCell+CycleCell poll; explicit cells matrix). typecheck + build-storybook pass.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/atoms/CycleCell/CycleCell.tsx — First cell-contract atom (n-state toggle) consuming the ContextStyleSet palette
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/MatrixGrid/MatrixGrid.tsx — Flagship generic grid engine + CellRenderPayload contract (the swappable-cell seam)
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/SegmentedBar/SegmentedBar.tsx — Proportional-bar engine generalizing ContextBudgetBar


## 2026-07-06

Fixed DateTile sm sizing (text overflow); added MonthGrid calendar engine (day-cell contract, UTC-safe date math, markers/heat via ContextStyleSet, min/max bounds, weekStartsOn); reformatted implementation diary to the diary-skill Step format; documented StyleBySpec. typecheck + build-storybook pass.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/MonthGrid/MonthGrid.tsx — New calendar-month base engine


## 2026-07-06

Added TimeGrid week/day calendar engine (lane-packed overlapping blocks, sticky headers/gutter, now indicator, slot-create, Mode-A renderBlock). Completes the generic base-engine layer (MatrixGrid, SegmentedBar, MonthGrid, TimeGrid + atoms). typecheck + build-storybook pass.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/TimeGrid/TimeGrid.tsx — New week/day calendar base engine


## 2026-07-06

Wired first engine to the DSL: added scheduling domain module (types/palettes/fixtures); StyleBySpec + CycleCellSpec + MatrixGridWidgetProps in ir.ts; MatrixGrid IR adapter + manifest; registered matrixGridWidget; availabilityMatrix preset; Widget IR/Renderer/Scheduling stories (poll renders from serialized IR). typecheck + build-storybook pass.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/MatrixGrid/MatrixGrid.widget.tsx — MatrixGrid IR adapter (spec -> React lambdas)


## 2026-07-06

Added MeetingPollPanel organism (first full Doodle screen): title + meta + deadline line + availability MatrixGrid with editable You row + RatioBadge/star tally footer + submit row. Presentational (DTO in, callbacks out); stories hold interactive state. typecheck + build-storybook pass.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/organisms/MeetingPollPanel/MeetingPollPanel.tsx — Participant-facing meeting poll organism


## 2026-07-06

Split ir.ts into ir/ modules (core/actions/cells/engines/props); fixed vite lib entry. Added SegmentedBar/MonthGrid/TimeGrid IR adapters+manifests (timeWidgetRegistry, time.dsl module). Wired StyleBySpec via ValueCellSpec + colorBy + resolveStyleByVars. Added pollResults/monthCalendar/weekCalendar presets + IR renderer stories (scheduling+calendar). Added organisms PollResultsPanel, CalendarMonthPanel, CalendarWeekPanel, BookingPagePanel + booking fixtures. typecheck + build-storybook + library build all pass.


## 2026-07-06

Wrote DSL handoff guide (reference/02) documenting every implemented widget, its IR contract, emitted actions, and a file/symbol-level plan for wiring MatrixGrid/SegmentedBar/MonthGrid/TimeGrid + presets into pkg/widgetdsl (helper maps, cell.cycle/value builders, time.dsl reconciliation, recipes, TS codegen). Corrected earlier error: the Goja DSL runtime DOES exist in-repo at pkg/widgetdsl.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module.go — Target of the DSL wiring plan


## 2026-07-06

Added intern-facing code review/design assessment (reference/03) for scheduling widgets; validation passed (typecheck, build, build-storybook, go test ./pkg/widgetdsl/..., pack:smoke, consumer:smoke); logged follow-up tasks for readOnly behavior, all-day TimeGrid contract, package exports, and focused tests.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/03-code-review-and-design-assessment-for-scheduling-widgets.md — New review/design assessment deliverable
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/tasks.md — Follow-up cleanup/test tasks from review


## 2026-07-06

Uploaded the code review/design assessment PDF to reMarkable at /ai/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS/RAGEVAL Schedule Widgets Review.pdf after dry-run.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/03-code-review-and-design-assessment-for-scheduling-widgets.md — Uploaded review report


## 2026-07-06

Added design-doc/02 Goja DSL layer implementation guide for scheduling widgets: module architecture, target JS APIs, implementation phases, recipe pseudocode, TypeScript declaration plan, tests, decisions, and intern checklist.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/02-goja-dsl-layer-design-and-implementation-guide-for-scheduling-widgets.md — New Goja DSL design/implementation guide
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/tasks.md — Added DSL implementation follow-up tasks


## 2026-07-06

Uploaded Goja DSL layer design guide to reMarkable at /ai/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS/RAGEVAL Schedule Goja DSL Guide.pdf after dry-run.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/02-goja-dsl-layer-design-and-implementation-guide-for-scheduling-widgets.md — Uploaded Goja DSL guide


## 2026-07-07

Added design-doc/03 composition-grammar DSL redesign proposal: named slots, domain views, shared bindings, intent actions, backwards application plan, examples, decisions, and migration/testing strategy.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/03-a-composition-grammar-for-the-widget-dsl.md — New composition-first DSL redesign proposal


## 2026-07-07

Expanded design-doc/03 after investigating researchctl: added lambda-inspired builder-callback form, .use(fragment) reuse model, callback-ID boundary, revised recommendation, and implementation slice for availabilityPoll builders.

### Related Files

- /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/builders.go — Source of .use(fragment) builder callback pattern
- /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/callbacks.go — Source of callback ID and trusted runtime boundary discussion
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/03-a-composition-grammar-for-the-widget-dsl.md — Updated with researchctl-inspired lambda DSL ideas


## 2026-07-07

Re-uploaded updated Widget DSL Composition Grammar PDF to reMarkable after adding researchctl lambda-inspired addendum (forced overwrite of prior PDF).

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/03-a-composition-grammar-for-the-widget-dsl.md — Updated and re-uploaded composition grammar document


## 2026-07-07

Expanded design-doc/03 to make the DSL redesign explicitly target existing ui/data/context_window/course/cms modules, including module responsibilities, target APIs, examples, shared kernel, migration roadmap, and acceptance criteria.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module.go — Existing ui/data/context_window/course/cms helper and recipe surfaces targeted by the refactor
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/03-a-composition-grammar-for-the-widget-dsl.md — Broadened DSL redesign beyond scheduling to existing modules


## 2026-07-07

Added follow-up tasks for broad DSL refactor inventory, shared composition kernel prototype, and cms/course/context_window builder-lambda domain views.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/tasks.md — New broad DSL refactor tasks


## 2026-07-07

Added design-doc/04 full clean-break Widget DSL redesign: single widget.dsl root module, Go-backed typed specs, TypeScript-facing declarations, builder lambdas, slots, fragments, bindings, intents, and examples from simple pages to go-go-course-scale CMS/course/context/schedule pages.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js — Current complex consumer page that shaped examples
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js — Current mixed DSL example page that shaped examples
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/04-full-widget-dsl-redesign-typed-builders-slots-fragments-and-domain-views.md — New full DSL redesign document


## 2026-07-07

Uploaded Full Widget DSL Redesign document to reMarkable at /ai/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS/RAGEVAL Full Widget DSL Redesign.pdf after dry-run.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/04-full-widget-dsl-redesign-typed-builders-slots-fragments-and-domain-views.md — Uploaded full DSL redesign document


## 2026-07-07

Added Widget DSL v3 phase/task tracker and diary Step 10 to guide step-by-step implementation with validation gates and commit boundaries.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — New v3 phase and task tracker
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 10 records tracker setup and working protocol
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/tasks.md — Added high-level v3 phase tracking tasks

## 2026-07-07

Committed Widget DSL v3 tracker and ticket docs in c99b32d33af17cd9aa86ff38cd999ff1aeb57533; updated diary Step 10 with commit reference.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — Committed tracker content
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Step 10 commit reference

## 2026-07-07

Completed Widget DSL v3 Phase 0 inventory: generated current module/helper/recipe export inventory, related source/generator files, and marked Phase 0 complete in the tracker.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — Phase 0 marked complete
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/04-widget-dsl-current-export-inventory.md — Generated current export inventory
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/scripts/01-widget-dsl-export-inventory.py — Inventory generator
