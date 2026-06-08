---
Title: Phased implementation task list
Ticket: RAGEVAL-WIDGET-IR-SEMANTIC-COMPONENTS
Status: active
Topics:
    - frontend
    - widget-ir
    - design-system
    - rag-evaluation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: packages/rag-evaluation-site/src/widgets/WidgetRenderer.stories.tsx
      Note: Phase checklist Storybook validation point
    - Path: packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx
      Note: Phase checklist renderer extension point
    - Path: packages/rag-evaluation-site/src/widgets/ir.ts
      Note: Phase checklist TypeScript IR extension point
    - Path: pkg/widgetdsl/module.go
      Note: Phase checklist Goja helper and recipe extension point
    - Path: pkg/widgetdsl/module_test.go
      Note: Phase checklist Go JSON serialization test point
    - Path: pkg/widgetschema/schema.go
      Note: Phase checklist schema synchronization point
ExternalSources: []
Summary: Granular phase-by-phase checklist for extending the existing Widget IR, WidgetRenderer, widget.dsl helpers, schemas, stories, and recipes to cover the expanded design-system package.
LastUpdated: 2026-06-08T00:00:00-04:00
WhatFor: Use as the implementation task board for the semantic Widget IR component support work.
WhenToUse: Before starting a phase, when slicing commits/PRs, or when reviewing whether WidgetRenderer and widget.dsl support are complete.
---


# Phased Implementation Task List

## Goal

This document turns the semantic Widget IR implementation guide into an executable task plan. It is intentionally granular so an intern can pick up one phase, understand exactly which files to touch, what tests to add, and what “done” means.

## Context

This work extends the **existing** Widget IR system. Do not create a second renderer, a second DSL module, or a Go-side HTML renderer.

The existing extension points are:

- `packages/rag-evaluation-site/src/widgets/ir.ts` — TypeScript Widget IR types and prop contracts.
- `packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx` — React renderer dispatch.
- `packages/rag-evaluation-site/src/widgets/WidgetRenderer.stories.tsx` — Storybook renderer coverage.
- `pkg/widgetdsl/module.go` — Goja `widget.dsl` / `rag.dsl` helpers and recipes.
- `pkg/widgetdsl/module_test.go` — Goja DSL JSON-serialization tests.
- `pkg/widgetschema/schema.go` — server schema component list.
- `pkg/widgetserver/server_test.go` — schema/action/page endpoint tests.
- `examples/xgoja/widget-site/verbs/sites.js` — end-to-end example site.

The recurring implementation loop is:

1. Add TypeScript IR props.
2. Add a `WidgetRenderer` case.
3. Add a WidgetRenderer Storybook story that uses IR, not JSX.
4. Add Goja helper or recipe support.
5. Update schema.
6. Add Go tests.
7. Run validation.
8. Commit one coherent phase.

---

## Phase 0: Baseline, inventory, and guardrails

### Purpose

Establish a clean baseline before changing code. This phase protects existing behavior and makes later failures easier to diagnose.

### Tasks

- [ ] Confirm working tree state.
  - [ ] Run `git status --short`.
  - [ ] Note unrelated untracked dirs and do not stage them.
  - [ ] Confirm current branch is `task/add-ui-dsl`.
- [ ] Read the core guidance.
  - [ ] Read `packages/rag-evaluation-site/GUIDELINES.md`.
  - [ ] Read `design-doc/01-semantic-widget-ir-component-support-implementation-guide.md`.
  - [ ] Read this task list.
- [ ] Run baseline package validation.
  - [ ] `pnpm --dir packages/rag-evaluation-site typecheck`
  - [ ] `pnpm --dir packages/rag-evaluation-site build`
  - [ ] `pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-package-storybook-widget-ir-baseline`
- [ ] Run targeted Go validation.
  - [ ] `go test ./pkg/widgetdsl ./pkg/widgetrunner ./pkg/widgetserver ./pkg/widgetschema -count=1`
  - [ ] If workspace-wide `go test ./...` fails because of unrelated module issues, record the exact error in the diary instead of blocking this ticket.
- [ ] Create/update diary entry.
  - [ ] Record baseline commands and outcomes.
  - [ ] Record current commit hash.
- [ ] Decide first implementation slice.
  - [ ] Confirm whether `SlideShell` belongs in Phase 1 or Phase 4.
  - [ ] Confirm `Text` helper name should be `textBlock` in Goja to avoid colliding with `rag.text()`.

### Files touched

- Usually only ticket docs/diary.

### Done criteria

- Baseline commands have been run or failures documented.
- No product code changed yet.
- Diary has a baseline entry.

---

## Phase 1: Extend existing WidgetRenderer for foundation, atoms, and layout

### Purpose

Prove the extension pattern on low-risk components before tackling DTO-heavy organisms. This phase extends the existing TypeScript IR and React renderer in-place.

### Components in scope

Foundation:

- `Text`
- `CodeText`
- `Divider`

Atoms:

- `ContextKindSwatch`
- `AnnotationBadge`
- `TranscriptRoleBadge`

Layout:

- `SectionBlock`
- `SplitPane`
- `SidebarShell`
- Optional after review: `SlideShell`

### TypeScript IR tasks

- [ ] Edit `packages/rag-evaluation-site/src/widgets/ir.ts`.
- [ ] Import required prop-related types.
  - [ ] Import `ContextPartKind` and `TranscriptRole` from `../context` if needed.
  - [ ] Import foundation prop unions if they are exported; otherwise define narrow JSON-friendly unions in `ir.ts`.
- [ ] Add component names to `RagWidgetType`.
  - [ ] `Text`
  - [ ] `CodeText`
  - [ ] `Divider`
  - [ ] `ContextKindSwatch`
  - [ ] `AnnotationBadge`
  - [ ] `TranscriptRoleBadge`
  - [ ] `SectionBlock`
  - [ ] `SplitPane`
  - [ ] `SidebarShell`
- [ ] Add prop interfaces.
  - [ ] `TextWidgetProps`
  - [ ] `CodeTextWidgetProps`
  - [ ] `DividerWidgetProps`
  - [ ] `ContextKindSwatchWidgetProps`
  - [ ] `AnnotationBadgeWidgetProps`
  - [ ] `TranscriptRoleBadgeWidgetProps`
  - [ ] `SectionBlockWidgetProps`
  - [ ] `SplitPaneWidgetProps`
  - [ ] `SidebarShellWidgetProps`
- [ ] For slot props, use `WidgetNode`.
  - [ ] `SplitPaneWidgetProps.left: WidgetNode`
  - [ ] `SplitPaneWidgetProps.right: WidgetNode`
  - [ ] `SidebarShellWidgetProps.sidebar?: WidgetNode`
  - [ ] `SidebarShellWidgetProps.header?: WidgetNode`
  - [ ] `SidebarShellWidgetProps.footer?: WidgetNode`
- [ ] Add new prop interfaces to `WidgetProps` union.
- [ ] Keep existing interfaces and component names unchanged.

### WidgetRenderer tasks

- [ ] Edit `packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx`.
- [ ] Add imports.
  - [ ] `Text`, `CodeText`, `Divider` from foundation.
  - [ ] `ContextKindSwatch`, `AnnotationBadge`, `TranscriptRoleBadge` from atoms.
  - [ ] `SectionBlock`, `SplitPane`, `SidebarShell` from layout.
- [ ] Add helper for node-valued props.

```ts
function renderNodeProp(node: WidgetNode | undefined, onAction?: WidgetActionHandler): ReactNode | undefined {
  return node ? renderWidgetNode(node, onAction) : undefined;
}
```

- [ ] Add switch cases in layer order.
- [ ] Implement render helpers.
  - [ ] `renderText`
  - [ ] `renderCodeText`
  - [ ] `renderDivider`
  - [ ] `renderContextKindSwatch`
  - [ ] `renderAnnotationBadge`
  - [ ] `renderTranscriptRoleBadge`
  - [ ] `renderSectionBlock`
  - [ ] `renderSplitPane`
  - [ ] `renderSidebarShell`
- [ ] Ensure `className` still passes through where supported.
- [ ] Do not pass arbitrary unsupported props via spread unless the existing renderer already does so for that component family.

### Storybook tasks

- [ ] Edit `packages/rag-evaluation-site/src/widgets/WidgetRenderer.stories.tsx`.
- [ ] Add story `FoundationAtomsAndLayout`.
- [ ] The story must use IR helpers (`component`, `text`) rather than direct JSX.
- [ ] Cover:
  - [ ] `Text` sizes/tones/weights.
  - [ ] `CodeText` with a path or ID.
  - [ ] `Divider` horizontal.
  - [ ] `ContextKindSwatch` pattern/tone examples.
  - [ ] `AnnotationBadge` selected/unselected.
  - [ ] `TranscriptRoleBadge` with assistant/tool roles.
  - [ ] `SectionBlock` with label/caption/divider.
  - [ ] `SplitPane` with left/right `WidgetNode` props.
  - [ ] `SidebarShell` with sidebar/header/footer slots.

### Validation tasks

- [ ] Run `pnpm --dir packages/rag-evaluation-site typecheck`.
- [ ] Run `pnpm --dir packages/rag-evaluation-site build`.
- [ ] Run Storybook build to `/tmp/rag-package-storybook-widget-ir-phase-1`.
- [ ] Browser sanity check the new story if Storybook is running locally.
- [ ] Update diary with commands, errors, and fixes.
- [ ] Commit Phase 1 TypeScript renderer work.

### Done criteria

- New Phase 1 components render via the existing `WidgetRenderer`.
- Existing WidgetRenderer stories still work.
- No Goja helpers yet unless explicitly included in the same phase slice.

---

## Phase 2: Add direct Goja helpers and schema support for Phase 1 components

### Purpose

Let Goja authors produce the Phase 1 component nodes through `require("widget.dsl")` without using raw `rag.component("Type", ...)` calls.

### Goja helper tasks

- [ ] Edit `pkg/widgetdsl/module.go`.
- [ ] Add helper names to `componentNames`.
  - [ ] `textBlock`
  - [ ] `codeText`
  - [ ] `divider`
  - [ ] `contextKindSwatch`
  - [ ] `annotationBadge`
  - [ ] `transcriptRoleBadge`
  - [ ] `sectionBlock`
  - [ ] `splitPane`
  - [ ] `sidebarShell`
- [ ] Add helper-to-component mapping in `componentTypes`.
  - [ ] `textBlock -> Text`
  - [ ] `codeText -> CodeText`
  - [ ] `divider -> Divider`
  - [ ] etc.
- [ ] Do not rename or overload existing `text`.
- [ ] Verify `componentFactory` handles node-valued props for slots because props are exported as maps.

### Schema tasks

- [ ] Edit `pkg/widgetschema/schema.go`.
- [ ] Add Phase 1 component types to `ComponentTypes`.
- [ ] Keep existing schema permissive unless strict prop schemas are part of a later task.

### Go test tasks

- [ ] Edit `pkg/widgetdsl/module_test.go`.
- [ ] Extend `TestRequireWidgetDSLExportsHelpers` or add a new test for new helper exports.
- [ ] Add `TestFoundationAtomLayoutHelpersAreJSONSerializable`.
- [ ] Test a script that builds a page with:
  - [ ] `rag.textBlock(...)`
  - [ ] `rag.contextKindSwatch(...)`
  - [ ] `rag.sectionBlock(...)`
  - [ ] `rag.splitPane({ left: ..., right: ... })`
  - [ ] `rag.sidebarShell({ sidebar: ..., header: ... }, ...)`
- [ ] Assert JSON round-trip.
- [ ] Assert key component `type` values appear.
- [ ] Edit `pkg/widgetserver/server_test.go` if needed to assert `/api/widget/schema` includes one or more new component names.

### Validation tasks

- [ ] `go test ./pkg/widgetdsl ./pkg/widgetschema ./pkg/widgetserver -count=1`
- [ ] `pnpm --dir packages/rag-evaluation-site typecheck`
- [ ] Update diary.
- [ ] Commit Phase 2 Go/schema work.

### Done criteria

- Goja scripts can create Phase 1 nodes with named helpers.
- `/api/widget/schema` lists the new component types.
- Existing Goja recipe tests still pass.

---

## Phase 3: Context diagram Widget IR support

### Purpose

Expose context-window visualization components through Widget IR using the existing JSON-compatible context DTOs.

### Components in scope

Molecules:

- `ContextLegend`
- `ContextBudgetBar`
- `ContextStripDiagram`
- `ContextStackDiagram`
- `ContextTreemap`

Organism:

- `ContextDiagramPanel`

### TypeScript IR tasks

- [ ] Import context DTOs from `../context`.
  - [ ] `ContextWindowSnapshot`
  - [ ] `ContextDiagramView`
  - [ ] `ContextDiagramStyle`
  - [ ] `ContextPartKind`
- [ ] Add component type names to `RagWidgetType`.
- [ ] Add props interfaces.
  - [ ] `ContextLegendWidgetProps`
  - [ ] `ContextBudgetBarWidgetProps`
  - [ ] `ContextStripDiagramWidgetProps`
  - [ ] `ContextStackDiagramWidgetProps`
  - [ ] `ContextTreemapWidgetProps`
  - [ ] `ContextDiagramPanelWidgetProps`
- [ ] Use `snapshot: ContextWindowSnapshot` for diagram props.
- [ ] Use `mode?: ContextDiagramStyle` where supported.
- [ ] Use `selectedPartId?: string` where supported.
- [ ] Add optional `onPartSelectAction?: ActionSpec` only if the React component supports selection callbacks; otherwise defer.

### WidgetRenderer tasks

- [ ] Import context diagram components from molecules/organisms.
- [ ] Add switch cases.
- [ ] Implement direct prop mapping.
- [ ] If `onPartSelectAction` is added, bind context as:

```ts
{ componentType: 'ContextDiagramPanel', partId, value: partId }
```

### Storybook tasks

- [ ] Add `ContextDiagramComponents` story under `Widget IR/Renderer`.
- [ ] Use `contextWindowFixture` or equivalent fixture data from `src/context/fixtures.ts`.
- [ ] Render a `DashboardGrid` of low-level diagram components.
- [ ] Render one `ContextDiagramPanel` organism.
- [ ] Include at least one selected segment.

### Goja/schema/test tasks

- [ ] Add direct helpers:
  - [ ] `contextLegend`
  - [ ] `contextBudgetBar`
  - [ ] `contextStripDiagram`
  - [ ] `contextStackDiagram`
  - [ ] `contextTreemap`
  - [ ] `contextDiagramPanel`
- [ ] Add component types to `pkg/widgetschema/schema.go`.
- [ ] Add Go JSON serialization test for a `ContextDiagramPanel` node.

### Done criteria

- A Widget IR page can render context diagrams using real package components.
- Goja can author the same nodes without raw component strings.

---

## Phase 4: Transcript, annotation, and anchored comment Widget IR support

### Purpose

Expose the transcript UI that was recently polished: title-bar message cards, session header, optional notes rail, and anchored comments.

### Components in scope

Molecules:

- `TranscriptSessionHeader`
- `TranscriptMessageCard`
- `AnnotationNoteCard`
- `AnchoredCommentCard`

Organisms:

- `TranscriptReaderPanel`
- `AnnotationRailPanel`
- `TranscriptWorkspacePanel`
- `AnchoredCommentRail`

### TypeScript IR tasks

- [ ] Import DTOs from `../context`.
  - [ ] `TranscriptMessage`
  - [ ] `TranscriptAnnotation`
  - [ ] `AnchoredComment`
- [ ] Add component type names to `RagWidgetType`.
- [ ] Add props interfaces.
  - [ ] `TranscriptSessionHeaderWidgetProps`
  - [ ] `TranscriptMessageCardWidgetProps`
  - [ ] `AnnotationNoteCardWidgetProps`
  - [ ] `TranscriptReaderPanelWidgetProps`
  - [ ] `AnnotationRailPanelWidgetProps`
  - [ ] `TranscriptWorkspacePanelWidgetProps`
  - [ ] `AnchoredCommentCardWidgetProps`
  - [ ] `AnchoredCommentRailWidgetProps`
- [ ] Add action props.
  - [ ] `onAnnotationSelectAction?: ActionSpec`
  - [ ] `onDismissAction?: ActionSpec`
- [ ] Keep message/annotation/comment DTOs as props; do not flatten them into DOM-like fields.

### WidgetRenderer tasks

- [ ] Import transcript and comment components.
- [ ] Add switch cases.
- [ ] Implement low-level card render helpers.
- [ ] Implement high-level panel helpers.
- [ ] Bind `onAnnotationSelectAction` with context:

```ts
{ componentType: 'TranscriptWorkspacePanel', annotationId, value: annotationId }
```

- [ ] Bind `onDismissAction` with context:

```ts
{ componentType: 'AnchoredCommentRail', commentId, value: commentId }
```

### Storybook tasks

- [ ] Add `TranscriptWithoutNotes` story.
- [ ] Add `TranscriptWithNotes` story.
- [ ] Add `AnchoredComments` story.
- [ ] Add a story with a custom `onAction` logger if practical.
- [ ] Ensure the no-notes story does not render the annotation rail.
- [ ] Ensure the with-notes story renders title-bar note cards.

### Goja/schema/test tasks

- [ ] Add direct helpers.
  - [ ] `transcriptSessionHeader`
  - [ ] `transcriptMessageCard`
  - [ ] `annotationNoteCard`
  - [ ] `transcriptReaderPanel`
  - [ ] `annotationRailPanel`
  - [ ] `transcriptWorkspacePanel`
  - [ ] `anchoredCommentCard`
  - [ ] `anchoredCommentRail`
- [ ] Add schema component types.
- [ ] Add Go test for JSON serialization of direct transcript helpers.
- [ ] Add Go test for action specs on transcript selection.

### Done criteria

- Widget IR can render transcript without notes and with notes.
- Action callback props serialize through Goja and bind in React.
- Existing non-transcript renderer stories still pass.

---

## Phase 5: Course, slide, sidebar, handout, and document Widget IR support

### Purpose

Expose the course/studio/handout UI components as semantic Widget IR nodes.

### Components in scope

Layout/molecules:

- `SlideShell` if not done in Phase 1
- `SidebarNav`
- `CourseStepNav`
- `FigureBlock`
- `KeyPointList`
- `KeyValueStrip`
- `CheckList`
- `StepList`
- `PersonSummary`
- `MarkdownArticle`
- `DocumentListPanel`
- `DocumentPreviewToolbar`

Organisms:

- `CourseLessonPanel`
- `CourseSlidePanel`
- `CourseStudioShell`
- `HandoutDocumentShell`

### TypeScript IR tasks

- [ ] Import DTOs:
  - [ ] `ContextCourse`
  - [ ] `ContextSlide`
  - [ ] `ContextWindowSnapshot`
  - [ ] `ContextHandoutDocument`
  - [ ] `ContextHandoutBundle`
- [ ] Add component names to `RagWidgetType`.
- [ ] Add prop interfaces for all in-scope components.
- [ ] Add action props for selection/navigation.
  - [ ] `onSelectAction?: ActionSpec` for document selection.
  - [ ] `onNavSelectAction?: ActionSpec` if shell/nav components support controlled selection.
- [ ] Use `WidgetNode` for visual/content slots in `FigureBlock` or `SlideShell` if needed.

### WidgetRenderer tasks

- [ ] Import in-scope components.
- [ ] Add switch cases.
- [ ] Implement slot rendering helpers where needed.
- [ ] Bind document selection action context:

```ts
{ componentType: 'HandoutDocumentShell', documentId, value: documentId }
```

- [ ] Bind nav selection action context:

```ts
{ componentType: 'CourseStudioShell', itemId, value: itemId }
```

### Storybook tasks

- [ ] Add `CourseStudioSlide` Widget IR story.
- [ ] Add `CourseLesson` Widget IR story.
- [ ] Add `HandoutDocumentShell` Widget IR story.
- [ ] Add dense/sidebar/document states if useful.

### Goja/schema/test tasks

- [ ] Add direct helpers for in-scope components.
- [ ] Add schema component types.
- [ ] Add Go JSON serialization tests for:
  - [ ] course lesson
  - [ ] course slide
  - [ ] handout document shell
  - [ ] markdown article

### Done criteria

- Widget IR can render a course slide inside the course studio shell.
- Widget IR can render a handout/document reader.
- Selection actions are represented as `ActionSpec` props where supported.

---

## Phase 6: Semantic recipes for product-level authoring

### Purpose

Reduce raw IR boilerplate by adding high-level recipes that expand into the direct component nodes from Phases 3-5.

### Recipes in scope

- `rag.recipes.contextDiagram(...)`
- `rag.recipes.annotatedTranscript(...)`
- `rag.recipes.courseStudio(...)`
- `rag.recipes.courseLesson(...)`
- `rag.recipes.courseSlide(...)`
- `rag.recipes.handout(...)`

### Design tasks

- [ ] For each recipe, write the expected input shape in a comment or doc test.
- [ ] Decide fallback/default behavior.
  - [ ] Missing transcript title.
  - [ ] Empty annotations.
  - [ ] Missing selected handout document.
  - [ ] Missing course nav sections.
- [ ] Use existing `normalizeActionSpec` for action-like options.
- [ ] Use helper functions to avoid repeated map literals.

### Go implementation tasks

- [ ] Edit `pkg/widgetdsl/module.go`.
- [ ] Add recipe exports in `install`.
- [ ] Add recipe functions.
- [ ] Add helpers if needed:
  - [ ] `componentNode`
  - [ ] `widgetNodeFromOption`
  - [ ] `mapFromOption`
  - [ ] maybe `stringValueOrMapFallback`
- [ ] Keep recipes pure: no fetch, no database access, no mutation.

### Test tasks

- [ ] Add `TestContextTranscriptCourseRecipesAreJSONSerializable`.
- [ ] Test each recipe at least once.
- [ ] Assert output component types.
- [ ] Assert action normalization produces `{ kind: "server", name: "..." }` for string actions.
- [ ] Assert recipes still JSON.stringify cleanly.

### Example tasks

- [ ] Add a recipe-based page to `examples/xgoja/widget-site/verbs/sites.js` or a new example script.
- [ ] Include at least one transcript or handout recipe.
- [ ] Include at least one action that refreshes state.

### Done criteria

- Authors can build high-level context/transcript/course/handout pages with recipe calls.
- Recipes produce only plain Widget IR.
- Tests prove JSON round-trip and structure.

---

## Phase 7: Schema, docs, and examples polish

### Purpose

Make the expanded component vocabulary discoverable and teachable.

### Schema tasks

- [ ] Ensure `pkg/widgetschema/schema.go` component list matches TypeScript renderer support.
- [ ] Add test that schema includes representative new components:
  - [ ] `ContextDiagramPanel`
  - [ ] `TranscriptWorkspacePanel`
  - [ ] `CourseStudioShell`
  - [ ] `HandoutDocumentShell`
- [ ] Decide whether to keep schema version `0.1.0` or bump.
- [ ] If version changes, update `rag.page()` default and tests consistently.

### Package docs tasks

- [ ] Update `packages/rag-evaluation-site/README.md`.
  - [ ] Add expanded Widget IR section.
  - [ ] Mention direct nodes vs recipes.
  - [ ] Include one small code example.
- [ ] Update `packages/rag-evaluation-site/GUIDELINES.md` only if new rules emerged.
- [ ] Update any bundled widget DSL docs if present.

### Example docs/tasks

- [ ] Update `examples/xgoja/widget-site/README.devctl.md` if example pages changed.
- [ ] Add curl command for new page if relevant.
- [ ] Add screenshot or smoke instructions if expected by repo conventions.

### Done criteria

- `/api/widget/schema` is accurate.
- Package README explains the new authoring model.
- Examples demonstrate at least one new semantic recipe.

---

## Phase 8: End-to-end validation and visual review

### Purpose

Validate the full path from Goja recipe to server endpoint to React rendering.

### Validation tasks

- [ ] Run TypeScript validation.
  - [ ] `pnpm --dir packages/rag-evaluation-site typecheck`
  - [ ] `pnpm --dir packages/rag-evaluation-site build`
  - [ ] `pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-package-storybook-widget-ir-complete`
- [ ] Run Go validation.
  - [ ] `go test ./pkg/widgetdsl ./pkg/widgetrunner ./pkg/widgetserver ./pkg/widgetschema -count=1`
- [ ] Run example smoke if available.
  - [ ] `make -C examples/xgoja/widget-site smoke`
- [ ] Start local example server if needed.
- [ ] Open browser to new recipe page.
- [ ] Capture screenshots of key pages/stories.
- [ ] Check console for errors.
- [ ] Check network for `/api/widget/pages/...` and `/api/widget/actions/...` behavior.

### Visual review tasks

- [ ] Compare Widget IR rendered transcript against package transcript Storybook component.
- [ ] Compare Widget IR context diagram against package context diagram stories.
- [ ] Compare course/handout rendered pages against package component stories.
- [ ] Record visual gaps in ticket diary.
- [ ] Do not tune component visuals inside this phase unless the issue is caused by incorrect renderer prop mapping.

### Done criteria

- Full path works from Goja to React.
- No missing renderer cases for planned components.
- Unknown widget fallback still works for invalid nodes.
- Action callbacks work for at least transcript selection and handout document selection.

---

## Phase 9: Final docs, changelog, and handoff

### Purpose

Make the implementation easy to review, resume, and maintain.

### Documentation tasks

- [ ] Update the primary design doc if implementation decisions changed.
- [ ] Update this phased task list with actual completion status.
- [ ] Update diary after every phase.
- [ ] Add changelog entries per phase.
- [ ] Relate modified files to the most relevant ticket docs using `docmgr doc relate`.

### Git tasks

- [ ] Review `git diff --stat`.
- [ ] Review important diffs manually.
- [ ] Stage only relevant files.
- [ ] Avoid staging unrelated untracked web dirs.
- [ ] Commit with focused messages per phase.
- [ ] Push branch.

### Final validation tasks

- [ ] `docmgr doctor --ticket RAGEVAL-WIDGET-IR-SEMANTIC-COMPONENTS --stale-after 30`
- [ ] Re-run key code validation from Phase 8.
- [ ] If asked, upload updated guide/task bundle to reMarkable.

### Done criteria

- All phases are either complete or explicitly deferred with rationale.
- Ticket docs match implementation reality.
- Review instructions identify files and commands.

---

## Cross-phase review checklist

Use this checklist for every component before marking it complete.

- [ ] React component is stable and story-covered.
- [ ] Props are JSON-compatible or converted through `WidgetNode` / `ActionSpec`.
- [ ] `RagWidgetType` includes the component string.
- [ ] `WidgetProps` includes the prop interface.
- [ ] `WidgetRenderer` imports the actual React component.
- [ ] `WidgetRenderer` has an explicit switch case.
- [ ] Renderer helper maps props intentionally, not blindly.
- [ ] Any callback prop uses `ActionSpec` and useful action context.
- [ ] WidgetRenderer story exercises the node.
- [ ] Goja helper exists if author-facing direct construction is desired.
- [ ] Recipe exists if the component is normally used as part of a larger product composition.
- [ ] Schema component list is updated.
- [ ] Go JSON serialization test covers helper or recipe.
- [ ] Existing Widget IR stories and recipes still work.

---

## Recommended commit slices

1. `feat(widget-ir): add foundation atom layout renderer nodes`
2. `feat(widget-dsl): add foundation atom layout helpers`
3. `feat(widget-ir): add context diagram nodes`
4. `feat(widget-ir): add transcript annotation nodes`
5. `feat(widget-ir): add course and handout nodes`
6. `feat(widget-dsl): add semantic context recipes`
7. `docs(widget-ir): document expanded semantic components`

Each slice should include its matching tests and diary update.
