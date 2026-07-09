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

## 2026-07-07

Committed Widget DSL v3 Phase 0 export inventory in 4ff1ae57d55f478addb71679718cf6b4e19bbb03; updated diary Step 11 with commit reference.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Step 11 commit reference
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/04-widget-dsl-current-export-inventory.md — Committed Phase 0 inventory

## 2026-07-07

Completed Widget DSL v3 Phase 1: added parallel widget.dsl module skeleton with raw/act/bind namespaces, provider exposure, TypeScript stub, runtime/provider tests, and tracker/diary updates.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module.go — widget.dsl module skeleton, raw namespace, bind namespace
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module_test.go — widget.dsl runtime coexistence tests
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript.go — widget.dsl TypeScript declaration stub
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/xgoja/providers/widgetsite/provider.go — widget.dsl provider exposure
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — Phase 1 marked complete
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 12 records Phase 1

## 2026-07-07

Committed Widget DSL v3 Phase 1 module skeleton in c5b50e83fb528d13128b7b062237a9b6c9fcdbf7; pre-commit ran Go tests and lint successfully.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module.go — Committed widget.dsl skeleton
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/reference/01-implementation-diary.md — Step 12 commit reference

## 2026-07-07

Updated Widget DSL v3 tracker to explicitly incorporate RAGEVAL-WIDGET-DECOMPOSITION review constraints: engine/contract/preset, AccessorSpec/SelectionSpec/ListItemSpec direction, context segment-engine target, descriptor/manifests, and opaque TS types.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — Updated with decomposition review constraints
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 13 records the planning correction
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-WIDGET-DECOMPOSITION--widget-library-decomposition-base-engines-contracts-and-dsl-ergonomics/design-doc/01-widget-library-decomposition-analysis-and-design.md — Review input for v3 constraints

## 2026-07-08

Started Widget DSL v3 Phase 2: added page/section builder kernel, builder callbacks, .use fragments, page validation/lowering, TypeScript declarations, runtime tests, and tracker/diary updates.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module.go — widget.dsl now installs page function
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module_test.go — Page builder runtime test
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript.go — Page/section declaration updates
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Initial v3 page/section builder kernel
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — Phase 2 marked in progress with completed subtasks
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 14 records Phase 2 start

## 2026-07-08

Committed Widget DSL v3 Phase 2 initial page builder kernel in 8f8abfebba039ef2fd057507f8d2c32d55fc3690; pre-commit ran Go tests and lint successfully.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Committed page/section builder kernel
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Step 14 commit reference

## 2026-07-08

Continued Widget DSL v3 Phase 2: added node/source specs, v3-only child normalization, slot specs/calls, slot helpers, TypeScript slot declarations, and runtime slot tests.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module.go — widget.dsl now installs v3 raw helpers
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module_test.go — Slot and child-normalization test
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript.go — Slot and SlotHelpers TypeScript declarations
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Node specs, slot specs, callV3Slot, v3 raw helpers, and slot helpers
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — Phase 2 subtasks marked complete
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 15 records slot/child-normalization slice

## 2026-07-08

Committed Widget DSL v3 Phase 2 slot/node kernel in f6aa36fdb3d5a7788072599803acea6e0216a328; pre-commit ran Go tests and lint successfully.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Committed slot and node kernel
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Step 15 commit reference

## 2026-07-08

Finished Widget DSL v3 Phase 2: AccessorSpec-aligned bind helpers, selection/list-item specs, action confirm/payload passthrough, recursive node validation, final TS declarations, tracker completion, and tests.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module.go — AccessorSpec bind helpers and v3 data namespace installation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module_test.go — Final Phase 2 runtime coverage
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript.go — Final Phase 2 declaration surface
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Selection/list item specs and node validation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — Phase 2 marked complete
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 16 records Phase 2 completion

## 2026-07-08

Committed final Widget DSL v3 Phase 2 completion in e3d9d955a1bc007a177d12f6fe51f09ed6c7ec71; pre-commit ran Go tests and lint successfully.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Committed final Phase 2 kernel
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Step 16 commit reference

## 2026-07-08

Completed Widget DSL v3 Phase 3: concrete ui namespace, page shell/density/breadcrumbs, section actions/metrics/metadata, runtime tests, TypeScript declarations, and tracker update.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module.go — widget.dsl installs v3 ui namespace
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module_test.go — Phase 3 UI runtime coverage
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript.go — Phase 3 declarations
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Phase 3 UI composition implementation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — Phase 3 marked complete
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 17 records Phase 3

## 2026-07-08

Committed Widget DSL v3 Phase 3 UI composition in 501a0714821bf5e150b1dc1525c47d1d54ef55ea; pre-commit ran Go tests and lint successfully.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Committed Phase 3 UI composition
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Step 17 commit reference

## 2026-07-08

Completed Widget DSL v3 Phase 4: data.fields/schema builder, data.collection over v2 specs/lowering, selection.urlParam, editor actions, cell helpers, MatrixGrid helper, declarations, and tests.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module_test.go — Phase 4 runtime tests
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript.go — Phase 4 TypeScript declarations
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Phase 4 data implementation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — Phase 4 marked complete
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 18 records Phase 4

## 2026-07-08

Committed Widget DSL v3 Phase 4 data namespace in ee18d7374a995e2ec4b31714ff4034a6434dce1d; pre-commit ran Go tests and lint successfully.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Committed Phase 4 data namespace
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Step 18 commit reference

## 2026-07-08

Completed Widget DSL v3 Phase 5: cms.mediaLibrary, cms.articleQueue, cms.markdownEditor, CMS intent wrappers, slot placeholders, declarations, tests, and tracker updates.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module.go — widget.dsl cms namespace installation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module_test.go — Phase 5 CMS tests
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript.go — Phase 5 declarations
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Phase 5 CMS implementation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — Phase 5 marked complete
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 19 records Phase 5

## 2026-07-08

Committed Widget DSL v3 Phase 5 CMS namespace in df1975484c1508cf8152e217957cff4b78fd8d35; pre-commit ran Go tests and lint successfully.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Committed Phase 5 CMS namespace
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Step 19 commit reference

## 2026-07-08

Completed Widget DSL v3 Phase 6: course shell, landing, slide deck, handouts, metadata form, agenda editor, material uploads, course intents, declarations, tests, and tracker updates.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module.go — widget.dsl course namespace installation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module_test.go — Phase 6 course tests
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript.go — Phase 6 declarations
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Phase 6 course implementation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — Phase 6 marked complete
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 20 records Phase 6

## 2026-07-08

Committed Widget DSL v3 Phase 6 course namespace in e13eccf8041f6d74ff7a5eabb273057c38d650df; pre-commit ran Go tests and lint successfully.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Committed Phase 6 course namespace
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Step 20 commit reference

## 2026-07-08

Completed Widget DSL v3 Phase 7: context styleSet/palette, diagram/workspace helpers, context intents, slot placeholders, declarations, tests, and tracker updates.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module.go — widget.dsl context namespace installation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module_test.go — Phase 7 context tests
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript.go — Phase 7 declarations
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Phase 7 context implementation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md — Phase 7 marked complete
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 21 records Phase 7

## 2026-07-08

Committed Widget DSL v3 Phase 7 context namespace in 48683f4f1d5b349efdb3c963edf96b0591e0e02f; pre-commit ran Go tests and lint successfully.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Committed Phase 7 context namespace
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Step 21 commit reference

## 2026-07-08

Completed Widget DSL v3 Phases 8 and 9: schedule/time namespace helpers, schedule/time declarations and tests, descriptor-backed namespace exports, descriptor tests, TypeScript fixture examples, and API reference doc.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module_test.go — Phase 8 runtime tests
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript_fixture_test.go — Phase 9 TypeScript fixture examples
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Phase 8 schedule/time runtime implementation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3_descriptors.go — Phase 9 descriptor inventory
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 22 records Phases 8 and 9
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/05-widget-dsl-v3-api-reference.md — Generated-style API reference

## 2026-07-08

Committed Widget DSL v3 Phases 8 and 9 in 0e5c85d68137036d1c2f2a07eaaa733e3c30b1f0; pre-commit ran Go tests and lint successfully.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Committed Phase 8 schedule/time runtime helpers
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3_descriptors.go — Committed Phase 9 descriptor inventory
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Step 22 commit reference

## 2026-07-08

Completed Widget DSL v3 Phase 10: runnable example scripts, renderer CLI, golden Widget IR snapshots, stability tests, README, and typed-map normalization for embedding data collection nodes in v3 pages.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/cmd/widgetdsl-v3-examples/main.go — Example renderer CLI
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/module.go — Typed map normalization for v3 embedding
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/testdata/v3/examples — Runnable Phase 10 widget.dsl example scripts
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3_examples_test.go — Golden example stability tests
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Diary Step 23 records Phase 10

## 2026-07-08

Committed Widget DSL v3 Phase 10 golden examples in 07beb87fcb8d6af03a6da55c137e3aa4c416942d; pre-commit ran Biome, broader Go tests, and lint successfully.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/testdata/v3/examples — Committed example scripts
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/testdata/v3/golden — Committed golden snapshots
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/01-implementation-diary.md — Step 23 commit reference

## 2026-07-08

Widget DSL v3 preview gallery hardened, browser regressions fixed, and Storybook regression stories added (commits 071dbb0, 57b701d, a926537).

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/examples/xgoja-widgetdsl-v3/jsverbs/server.js — xgoja preview gallery server
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/WidgetRenderer.v3-regressions.stories.tsx — New regression story suite
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Runtime lowering fixes


## 2026-07-08

Phase 11 complete: documented widget.dsl v3 provider integration, migration/cutover workflow, and legacy usage checker (commit e7c28b7).

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/xgoja/providers/widgetsite/doc/01-widget-dsl-getting-started.md — Provider docs
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/06-widget-dsl-v3-integration-and-migration-guide.md — Migration guide
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/scripts/02-report-legacy-widget-dsl-usage.py — Migration checker


## 2026-07-08

Replaced Python migration checker with Go/tree-sitter CLI and reusable scanner package (commit 5872998).

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/cmd/widgetdsl-migration-checker/main.go — CLI entrypoint
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/migrationcheck/checker.go — Parser-backed scanner
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/migrationcheck/checker_test.go — Scanner tests

