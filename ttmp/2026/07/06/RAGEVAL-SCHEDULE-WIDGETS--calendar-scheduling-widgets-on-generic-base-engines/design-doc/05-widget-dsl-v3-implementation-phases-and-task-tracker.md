---
Title: Widget DSL v3 Implementation Phases and Task Tracker
Ticket: RAGEVAL-SCHEDULE-WIDGETS
Status: active
Topics:
    - ui-dsl
    - widget-ir
    - design-system
    - react
    - frontend-architecture
    - intern-guide
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://pkg/widgetdsl/module.go
      Note: Current runtime module registry and helper/recipe implementation to refactor behind widget.dsl
    - Path: repo://pkg/widgetdsl/typescript.go
      Note: Current declaration generation path to replace with descriptor-driven widget.dsl declarations
    - Path: repo://pkg/widgetdsl/v2/spec/types.go
      Note: Existing Go spec model informing the v3 backing specs
    - Path: repo://pkg/widgetdsl/v2_builders.go
      Note: Existing typed builder callback implementation to reuse for data namespace
    - Path: ws://go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js
      Note: Complex consumer fixture target for CMS/course/data integration
    - Path: ws://go-go-course/cmd/go-go-course/lib/pages/common.js
      Note: Current course shell wrapper targeted by course namespace
    - Path: ws://go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js
      Note: First real consumer fixture target for v3 port
    - Path: ws://go-go-course/cmd/go-go-course/lib/pages/handouts.js
      Note: Current handout page targeted by course namespace
ExternalSources: []
Summary: Detailed phase plan, task tracker, validation gates, diary protocol, and commit cadence for implementing the parallel clean Widget DSL v3 module.
LastUpdated: 2026-07-07T16:40:00-04:00
WhatFor: Use this document to track the step-by-step implementation of the new widget.dsl module while existing ui/data/cms/course/context modules remain available.
WhenToUse: Read before starting or resuming Widget DSL v3 implementation; update after every phase or meaningful subphase.
---


# Widget DSL v3 Implementation Phases and Task Tracker

This document turns the clean Widget DSL redesign into an implementation plan. The
new DSL should be implemented in parallel as `widget.dsl`; existing modules can stay
available for current scripts. The new module is not constrained by old API names,
but it should reuse current lowerers, Widget IR adapters, and React components where
that is the fastest reliable path.

The plan is intentionally phase-based. Each phase has tasks, acceptance criteria,
validation commands, and a suggested commit boundary. Update this document and the
implementation diary as work progresses.

---

## Working protocol

### Diary protocol

For every non-trivial step:

1. Implement the smallest coherent slice.
2. Run the validation gate for that slice.
3. Append a diary entry to `reference/01-implementation-diary.md` using the diary
   skill format.
4. Relate any newly modified files to the relevant design/tracker docs.
5. Update this tracker's status checkboxes if phase/task state changed.
6. Commit the focused code/docs changes with a clear message.

### Commit protocol

Commit at these boundaries:

- one commit for this tracker + diary setup;
- one commit for each phase that compiles/tests;
- one commit for documentation-only follow-up after a code phase, if the diary or
  tracker changes are substantial;
- never stage unrelated untracked ticket or widget work accidentally.

Before each commit:

```bash
git status --short
git diff -- <files-you-plan-to-stage>
git add <specific files>
git diff --cached --stat
git diff --cached --check
git commit -m "Widget DSL v3: <focused message>"
```

### Validation baseline

Use the smallest relevant gate while iterating, then the full gate before a phase is
done.

```bash
# Go DSL/runtime checks
go test ./pkg/widgetdsl/... -count=1

# Frontend checks when emitted IR or TS package contracts change
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build

# Documentation hygiene
docmgr validate frontmatter --doc <doc>
docmgr doctor --ticket RAGEVAL-SCHEDULE-WIDGETS --stale-after 30
```

---

## Phase 0 — Baseline inventory and scaffolding

**Goal:** Establish the exact current surface, add the tracking artifacts, and make
sure the implementation starts from evidence rather than memory.

**Status:** in progress.

### Tasks

- [x] Create this phase/task tracker document.
- [ ] Append a diary entry recording the start of Widget DSL v3 implementation.
- [ ] Relate this tracker to current runtime, TypeScript, v2 spec, and go-go-course
      consumer files.
- [ ] Generate an inventory of current module exports from `moduleSpecs`:
  - module name;
  - helper name;
  - component type;
  - recipe name;
  - special objects (`action`, `cell`, context style helpers);
  - whether each helper is generic, engine-level, or domain-level.
- [ ] Store the inventory as `ttmp/.../scripts/01-widget-dsl-export-inventory.*` or
      `reference/04-widget-dsl-current-export-inventory.md`.
- [ ] Identify the first go-go-course fixture to port. Default choice:
      `cmd/go-go-course/lib/pages/dsl-examples.js`, because it exercises all public
      module families.

### Acceptance criteria

- `docmgr doctor` passes.
- The tracker names all implementation phases.
- Current surface inventory exists and is related to the ticket.
- No production code is changed in this phase.

### Suggested commit

`Widget DSL v3: add implementation tracker`

---

## Phase 1 — Add the parallel `widget.dsl` root module skeleton

**Goal:** Add a new module without changing old modules or old script behavior.

**Status:** not started.

### Tasks

- [ ] Add `WidgetV3ModuleName = "widget.dsl"` or equivalent module constant.
- [ ] Register a new native module spec for `widget.dsl`.
- [ ] Export root namespace objects:
  - `page`
  - `ui`
  - `data`
  - `cms`
  - `course`
  - `context`
  - `schedule`
  - `time`
  - `bind`
  - `act`
  - `style`
  - `raw`
- [ ] Implement `raw.component(type, props?, children?)` as the explicit escape
      hatch over the existing `buildComponent` logic.
- [ ] Export `raw.text`, `raw.element`, and `raw.fragment` or decide that `ui.text`
      and `ui.fragment` are enough.
- [ ] Add a runtime test proving `require("widget.dsl")` works and old modules still
      load.
- [ ] Add a TypeScript declaration stub for `widget.dsl`.
- [ ] Add provider exposure if module selection is explicit in the xgoja provider.

### Acceptance criteria

- `go test ./pkg/widgetdsl/... -count=1` passes.
- Existing module tests still pass.
- A Goja snippet can require both `widget.dsl` and `ui.dsl` in the same runtime.
- `widget.dsl.raw.component("Panel", { title:"X" })` emits a valid component node.

### Suggested commit

`Widget DSL v3: add parallel widget.dsl module skeleton`

---

## Phase 2 — Core spec kernel: pages, builders, fragments, slots, bindings, actions

**Goal:** Build the reusable kernel once before adding domain-specific APIs.

**Status:** not started.

### Tasks

#### Page and node specs

- [ ] Define Go structs or internal representations for:
  - `WidgetV3PageSpec`
  - `WidgetV3NodeSpec`
  - `WidgetV3SectionSpec`
  - `WidgetV3SourceSpan`
- [ ] Implement lowerers from specs to current Widget IR map shape.
- [ ] Add validation helpers for required page ID/title/root invariants.

#### Builder callback machinery

- [ ] Add shared `applyBuilderCallback(builder, cb)` helper modeled after
      `researchctl` and `codesign`.
- [ ] Add `.use(fragment)` convention to page/section builders.
- [ ] Ensure fragments can return the builder or `undefined`.
- [ ] Ensure non-function fragment errors are clear and include the builder name.

#### Child and slot machinery

- [ ] Add child normalization that flattens arrays and drops `null`, `undefined`,
      and `false`.
- [ ] Add `SlotSpec` representation.
- [ ] Add `callSlot(slot, context, fallback)` helper.
- [ ] Add stable slot helper object `h` with initial helpers:
  - `text`
  - `caption`
  - `strong`
  - `stack`
  - `inline`
  - `card`
  - `button`
  - `badge`
  - `raw`

#### Bindings and actions

- [ ] Add `bind.field`, `bind.path`, `bind.map`, `bind.template`, `bind.context`,
      and `bind.const`.
- [ ] Add `act.server`, `act.navigate`, `act.download`, `act.event`, `act.copy`.
- [ ] Ensure action `confirm` and payload bindings lower to the existing
      `ActionSpec` shape.

#### TypeScript declarations

- [ ] Add core types:
  - `JsonValue`
  - `WidgetNodeSpec`
  - `WidgetPageSpec`
  - `Fragment<TBuilder>`
  - `Slot<TContext>`
  - `BindingSpec`
  - `ActionSpec`
  - `SlotHelpers`
- [ ] Add declaration tests for the core API.

### Acceptance criteria

- A TypeScript fixture using `page("Hello", p => p.section(...)).toPage()` compiles.
- Go tests show builder callbacks, `.use`, slots, bindings, and actions work.
- Emitted page IR renders through existing WidgetRenderer assumptions.

### Suggested commit

`Widget DSL v3: add core builder and spec kernel`

---

## Phase 3 — UI namespace: page composition and generic primitives

**Goal:** Make `widget.dsl` useful for simple static pages and shared layout.

**Status:** not started.

### Tasks

- [ ] Implement `page(titleOrOptions, builder)` root function.
- [ ] Implement `PageBuilder` methods:
  - `.id(id)`
  - `.title(title)`
  - `.meta(key, value)`
  - `.shell(shellSpec)` placeholder
  - `.density(value)`
  - `.breadcrumb(label, href)`
  - `.section(title, callback)`
  - `.view(nodeOrView)`
  - `.toPage()`
  - `.validate()`
- [ ] Implement `SectionBuilder` methods:
  - `.caption(text)`
  - `.anchor(id)`
  - `.tone(tone)`
  - `.actions(callback)`
  - `.text(value)`
  - `.view(nodeOrView)`
  - `.metric(label, value, options?)`
  - `.metadata(record)`
- [ ] Implement generic `ui` helpers backed by current components:
  - `ui.callout`
  - `ui.stack`
  - `ui.inline`
  - `ui.card`
  - `ui.button`
  - `ui.caption`
  - `ui.badge`
  - `ui.metadata`
  - `ui.form`
- [ ] Add examples/tests for:
  - smallest useful page;
  - page with actions;
  - fragments.

### Acceptance criteria

- Simple page examples from design-doc/04 execute in Goja and emit valid page IR.
- No old module imports are required for simple pages.
- Old `ui.dsl` behavior remains untouched.

### Suggested commit

`Widget DSL v3: implement page and ui composition`

---

## Phase 4 — Data namespace: fields, collections, records, matrices

**Goal:** Replace public `data.v2.dsl` with `widget.dsl.data` while reusing the good
v2 spec/lowering foundation.

**Status:** not started.

### Tasks

#### Fields and schema

- [ ] Add `data.fields<T>(callback)` builder.
- [ ] Support field builders:
  - `.key`
  - `.primary`
  - `.short`
  - `.prose`
  - `.count`
  - `.status`
  - `.date`
  - `.currency`
  - `.media`
  - `.url`
- [ ] Lower to existing or extended `v2/spec.SchemaSpec`.

#### Collections

- [ ] Add `data.collection<T>(rows, callback)`.
- [ ] Add collection builder methods:
  - `.id(name)`
  - `.schema(schema)`
  - `.empty(message)`
  - `.select(selection)`
  - `.table(callback)`
  - `.edit(callback)`
  - `.masterDetail(callback?)`
  - `.validate()`
  - `.toNode()`
- [ ] Reuse `v2/spec.CollectionSpec.Validate` and lowering where practical.
- [ ] Add `data.selection.urlParam`.
- [ ] Add editor actions: create, submit, reorder, remove.

#### Cells and matrices

- [ ] Add `data.cell.field`, `data.cell.status`, `data.cell.template`,
      `data.cell.cycle`, and `data.cell.value`.
- [ ] Add `data.matrix<T>(rows, callback)` engine-level helper.
- [ ] Lower matrix to `MatrixGrid` IR.

### Acceptance criteria

- Current go-go-course `dsl-examples` table/select/master-detail/action examples
  have v3 equivalents.
- Tests compare old `data.v2.dsl` output and new `widget.dsl.data` output for at
  least one table and one master-detail case.
- `go test ./pkg/widgetdsl/... -count=1` passes.

### Suggested commit

`Widget DSL v3: implement data collections and matrix engine`

---

## Phase 5 — CMS namespace: media library, article queue, markdown editor

**Goal:** Port current CMS authoring patterns to typed domain views.

**Status:** not started.

### Tasks

- [ ] Define DTO declarations for:
  - `CmsAsset`
  - `CmsArticleSummary`
  - `CmsUploadState`
- [ ] Add `cms.mediaLibrary(assets, callback)` over `MediaLibraryPanel`.
- [ ] Add media library builder methods:
  - `.selection(mode)`
  - `.selected(ids)`
  - `.query(value)`
  - `.kindFilter(value)`
  - `.page(page, pageCount)`
  - `.empty(message)`
  - `.accept(mimeList)`
  - `.asset(slot)`
  - `.details(slot)`
  - `.toolbar(callback)`
  - `.onSelect(action)`
  - `.onOpen(action)`
  - `.onUpload(action)`
- [ ] Add `cms.articleQueue(articles, callback)` over `ArticleListPanel`.
- [ ] Add article queue slots/actions for article row, row actions, filters,
      create/publish/archive/preview.
- [ ] Add `cms.markdownEditor(body, callback)` over current markdown editor/form
      components.
- [ ] Add `cms.intent.*` wrappers.
- [ ] Port the media-library section from `go-go-course` admin page as a fixture.

### Acceptance criteria

- V3 fixture emits a media library equivalent to current `cmsDsl.recipes.mediaLibrary`.
- Intent wrappers hide engine-level action context names from examples.
- Existing CMS renderer stories can render emitted IR without new React components.

### Suggested commit

`Widget DSL v3: implement CMS domain views`

---

## Phase 6 — Course namespace: shell, landing, slides, handouts, material admin

**Goal:** Replace current course shell/handout/slide component calls with typed
course domain views.

**Status:** not started.

### Tasks

- [ ] Define DTO declarations for:
  - `CourseDefinition`
  - `CourseNavSection`
  - `CourseNavItem`
  - `CourseSlide`
  - `SlideDeck`
  - `HandoutDocument`
  - `HandoutBundle`
  - `CourseMaterialIndex`
- [ ] Add `course.shell(definition, callback)` shell spec.
- [ ] Add `course.landing(definition, callback)`.
- [ ] Add `course.slideDeck(deck, callback)`.
- [ ] Add `course.handouts(bundle, callback)`.
- [ ] Add `course.metadataForm(metadata, callback)`.
- [ ] Add `course.agendaEditor(items, callback)` as a domain wrapper over
      `data.collection`.
- [ ] Add `course.materialUploads(material, callback)`.
- [ ] Add `course.intent.*` wrappers for navigation, handout select/download/print,
      slide navigation, agenda edit, and material upload/delete.
- [ ] Port current `courseShellPage` and `handouts.js` as fixtures.

### Acceptance criteria

- A v3 course shell page can replace `courseDsl.recipes.courseStudio` in a fixture.
- Handout fixture lowers to existing handout/rich article components.
- Course examples compile against generated TypeScript declarations.

### Suggested commit

`Widget DSL v3: implement course domain views`

---

## Phase 7 — Context namespace: diagrams and transcript workspace

**Goal:** Provide task-level context-analysis APIs over existing context-window
components.

**Status:** not started.

### Tasks

- [ ] Define DTO declarations for:
  - `ContextSnapshot`
  - `ContextPart`
  - `Transcript`
  - `TranscriptMessage`
  - `Annotation`
- [ ] Add `context.styleSet(callback)` and `context.palette(nameOrOptions)`.
- [ ] Add `context.diagram(snapshot, callback)`.
- [ ] Add `context.workspace(session, callback)`.
- [ ] Add slots for:
  - message;
  - annotation;
  - diagram legend;
  - empty state.
- [ ] Add `context.intent.selectPart` and `context.intent.selectAnnotation`.
- [ ] Port the context section from go-go-course DSL examples as a fixture.

### Acceptance criteria

- V3 context diagram emits `ContextDiagramPanel` or equivalent current IR.
- V3 context workspace emits transcript + diagram components without raw calls.
- Tests cover style-set builder lowering.

### Suggested commit

`Widget DSL v3: implement context workspace views`

---

## Phase 8 — Schedule and time namespaces

**Goal:** Bring the scheduling/calendar work into the same v3 grammar.

**Status:** not started.

### Tasks

- [ ] Define DTO declarations for:
  - `AvailabilityPoll`
  - `AvailabilityResponse`
  - `PollOption`
  - `PollTally`
  - `CalendarEvent`
  - `TimeRange`
- [ ] Add `schedule.availabilityPoll(poll, callback)`.
- [ ] Add `schedule.pollSummary(poll, tallies, callback)`.
- [ ] Add `schedule.bookingPicker(availability, callback)` if underlying frontend
      behavior is ready.
- [ ] Add `schedule.intent.toggleAvailability` and `schedule.intent.submitResponse`.
- [ ] Add `time.month(eventsOrMarkers, callback)`.
- [ ] Add `time.week(events, callback)`.
- [ ] Add `time.range.week`, `time.format`, `time.formatRange`, and `time.slotLabel`.
- [ ] Add `time.intent.selectDay` and `time.intent.selectEvent`.
- [ ] Ensure `TimeGrid.allDay` is not exposed until the frontend contract is fixed.

### Acceptance criteria

- Availability poll lowers to `MatrixGrid` and hides engine context names.
- Calendar week lowers to `TimeGrid` and hides block-conversion details.
- Tests cover read-only poll action behavior and all-day omission/decision.

### Suggested commit

`Widget DSL v3: implement schedule and time views`

---

## Phase 9 — Descriptor-driven TypeScript declarations and docs

**Goal:** Prevent runtime/declaration/doc drift.

**Status:** not started.

### Tasks

- [ ] Introduce namespace/view/builder descriptors as the source of truth.
- [ ] Generate `widget.dsl` TypeScript declarations from descriptors.
- [ ] Include slot context interfaces in generated declarations.
- [ ] Include DTO interfaces or importable DTO declaration modules.
- [ ] Add declaration fixture tests for:
  - simple page;
  - data collection;
  - CMS media library;
  - course shell/handouts;
  - context workspace;
  - schedule poll/time week.
- [ ] Generate or write a `widget.dsl` API reference doc.

### Acceptance criteria

- TypeScript fixture examples from design-doc/04 compile.
- Adding a new view descriptor adds runtime export and declaration output in one
  place.
- Existing docs can link to generated API reference.

### Suggested commit

`Widget DSL v3: generate declarations from descriptors`

---

## Phase 10 — Golden go-go-course fixture port

**Goal:** Prove the new DSL on real pages before broad migration.

**Status:** not started.

### Tasks

- [ ] Add a fixture directory under the ticket or under `pkg/widgetdsl/testdata` for
      v3 go-go-course examples.
- [ ] Port `dsl-examples.js` examples:
  - simple table;
  - selectable table;
  - master-detail editor;
  - row actions;
  - all-modules gallery.
- [ ] Port `admin-course-cms.js` in slices:
  - shell;
  - metadata forms;
  - agenda editor;
  - file upload section;
  - media library;
  - preview actions.
- [ ] Port `handouts.js` and one slide page.
- [ ] Execute each fixture in Goja and snapshot emitted Widget IR.
- [ ] Optionally compare current old-module output with new v3 output at the
      component-type/action-contract level.

### Acceptance criteria

- All fixtures execute in Goja without browser dependencies.
- Golden snapshots are stable.
- The complex Course CMS fixture contains no `raw.component` calls except documented
  experimental exceptions.

### Suggested commit

`Widget DSL v3: add go-go-course golden fixtures`

---

## Phase 11 — Integration and cutover guidance

**Goal:** Decide how hosts adopt `widget.dsl` while old modules remain available.

**Status:** not started.

### Tasks

- [ ] Update xgoja provider docs to explain parallel module availability.
- [ ] Add example xgoja config selecting `widget.dsl`.
- [ ] Add a migration guide from old module imports to `widget.dsl` namespaces.
- [ ] Add a lint/check script that can report legacy module usage in first-party
      scripts.
- [ ] Decide whether old modules remain indefinitely or become deprecated after v3
      fixtures cover first-party pages.

### Acceptance criteria

- New hosts can select only `widget.dsl` and run v3 examples.
- Existing hosts can keep old modules and run current scripts.
- Migration guidance is explicit and test-backed.

### Suggested commit

`Widget DSL v3: document integration and migration path`

---

## Current status snapshot

| Phase | Status | Next concrete action |
|---|---|---|
| Phase 0 — Baseline inventory and scaffolding | in progress | Append diary entry and commit tracker setup. |
| Phase 1 — `widget.dsl` skeleton | not started | Add root module and raw escape hatch. |
| Phase 2 — Core spec kernel | not started | Add page/node specs, builder callbacks, fragments, slots, bind/act. |
| Phase 3 — UI namespace | not started | Implement simple page examples. |
| Phase 4 — Data namespace | not started | Port data.v2 examples under `widget.dsl.data`. |
| Phase 5 — CMS namespace | not started | Port media library section. |
| Phase 6 — Course namespace | not started | Port course shell and handouts. |
| Phase 7 — Context namespace | not started | Port context diagram/workspace fixture. |
| Phase 8 — Schedule/time namespaces | not started | Add availability poll and calendar week views. |
| Phase 9 — Descriptor-driven declarations/docs | not started | Generate declarations from descriptors. |
| Phase 10 — Golden go-go-course fixtures | not started | Snapshot old/new page fixtures. |
| Phase 11 — Integration/cutover guidance | not started | Document provider and migration workflow. |
