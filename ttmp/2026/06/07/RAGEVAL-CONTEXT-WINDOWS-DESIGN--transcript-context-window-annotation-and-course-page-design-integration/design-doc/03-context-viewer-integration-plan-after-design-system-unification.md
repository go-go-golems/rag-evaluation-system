---
Title: Context Viewer Integration Plan After Design-System Unification
Ticket: RAGEVAL-CONTEXT-WINDOWS-DESIGN
Status: active
Topics:
    - design-system
    - frontend-architecture
    - react
    - rag
    - ui-dsl
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: packages/rag-evaluation-site/src/index.ts
      Note: Package export surface for upcoming context DTOs/components
    - Path: packages/rag-evaluation-site/src/widgets/ir.ts
      Note: Future Widget IR extension point after React components stabilize
    - Path: pkg/widgetdsl/module.go
      Note: Future Goja recipe extension point
    - Path: ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/sources/03-context-viewer-design-iteration/diagrams.jsx
      Note: Prototype diagram renderer source
    - Path: ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/sources/03-context-viewer-design-iteration/patterns.jsx
      Note: Prototype kind vocabulary and diagram style source
    - Path: ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/sources/03-context-viewer-design-iteration/screens2.jsx
      Note: Prototype transcript/handout source
    - Path: ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/sources/03-context-viewer-design-iteration/screens3.jsx
      Note: Prototype comments/slides source
ExternalSources: []
Summary: Concrete implementation plan for integrating the context-viewer prototype after packages/rag-evaluation-site became the canonical reusable UI layer.
LastUpdated: 2026-06-07T18:11:51.274262921-04:00
WhatFor: Use as the handoff plan before coding context-window, transcript, annotation, and course surfaces.
WhenToUse: After reading the context-viewer analysis and before adding production components or Widget IR recipes.
---


# Context Viewer Integration Plan After Design-System Unification

## Executive Summary

The design-system unification ticket is complete: reusable foundation, atoms, layout primitives, shared molecules, and Widget IR rendering now belong to `packages/rag-evaluation-site`, while `web` owns app routing, API containers, backend-connected feature panels, and app-only molecules. That decision resolves the main open question from the context-viewer analysis: **new reusable context-viewer vocabulary should start in the package**, and the web app should consume it directly.

The next integration should therefore be Storybook-first and package-first:

1. promote context-viewer data contracts and reusable visual components into `packages/rag-evaluation-site`;
2. add package Storybook stories for static, realistic transcript/context-window/course states;
3. compose app-specific pages/containers under `web` only after the reusable package pieces are stable;
4. expose high-level Widget IR / Goja recipes last, once the React component vocabulary is proven.

This avoids copying the prototype's CDN/global/inline-style implementation while preserving its product architecture: context diagrams, transcript annotation, anchored comments, course/slide/handout reading surfaces, and the slightly more readable typography direction.

## Problem Statement

The context-viewer prototype is valuable but not production-ready architecture. It uses global JSX scripts, `window` exports, inline state, broad `.mac-*` classes, and duplicated primitive components. The current repository now has a clearer split:

```text
packages/rag-evaluation-site
  reusable tokens/theme bridge
  foundation/atoms/layout/shared molecules
  Widget IR types + WidgetRenderer
  package-owned Storybook visual baselines

web
  app routing and navigation
  services/API/store/backend containers
  app-specific organisms/pages
  local-only molecules such as CoveragePanel and QueryPresetList
```

The integration problem is to translate prototype concepts into this split without regressing the package boundary or creating another duplicated design-system layer.

## Proposed Solution

### Target boundary

Use this ownership rule for all new context-viewer code:

| Concern | Landing zone | Rationale |
|---|---|---|
| Context-window DTOs and visual component props | `packages/rag-evaluation-site/src/context/*` or `src/components/context/*` | Shared by Storybook, web app, and future Widget IR recipes. |
| Diagram primitives/molecules (`KindSwatch`, legend, budget bar, strip/stack/treemap) | `packages/rag-evaluation-site/src/components/context/*` | Reusable product vocabulary, not backend-specific. |
| Transcript message and annotation presentational pieces | `packages/rag-evaluation-site/src/components/transcript/*` | Reusable visual vocabulary for static stories and future DSL recipes. |
| Course/slide/handout presentational panels | Package first if generic; web only if tied to app routing/API | Course content is product-level but can be fixture/story-driven initially. |
| Upload/parsing workflow and local JSON validation UI | `web` page/container unless generalized later | It is workflow/application behavior, not core component vocabulary. |
| API calls, persistence, document/session loading | `web/src/services`, `web/src/components/pages`, containers | Backend ownership remains in the app. |
| Widget IR node/recipe support | package `src/widgets/*`, Go `pkg/widgetdsl/*` after React components stabilize | Preserve “Goja authors data; React renders UI.” |

### Atomic layer split

The implementation should preserve the recovered ownership rule from `design-doc/01-design-integration-research-and-implementation-plan.md`: tokens/foundation/atoms/layout stay generic, molecules own reusable data-display patterns, organisms own feature panels with DTO-shaped props, pages compose organisms, and containers own API/state side effects.

For the context-viewer work, that means the first component pass should not create one `ContextViewer` mega-component. Split it this way:

| Layer | Context-viewer additions | Notes |
|---|---|---|
| tokens/theme | typography comparison or scoped readable-mode variables | Review separately; do not silently change global defaults. |
| foundation | `SectionLabel` only if `Caption` cannot cover it | Prefer extending existing `Caption`/`Text` before adding primitives. |
| atoms | `Chip` / `AnnotationBadge` / `ContextKindSwatch` | Small, reusable semantic markers; no layout or data fetching. |
| layout | no new layout primitives initially | Use `Panel`, `Stack`, `Inline`, `DashboardGrid`, `ScrollRegion`, `TabList`, `FormRow`. |
| molecules | `ContextLegend`, `ContextBudgetBar`, `ContextStripDiagram`, `ContextStackDiagram`, `ContextTreemap`, `TranscriptMessageCard`, `TranscriptTokenBar`, `AnnotationNoteCard`, `CourseStepNav` | Reusable visual/data-display units with typed DTO props. |
| organisms | `ContextWindowVisualizerPanel`, `TranscriptReaderPanel`, `AnnotationRailPanel`, `AnchoredCommentRail`, `CourseLessonPanel`, `CourseSlidePanel` | Feature panels that compose molecules and own controlled selection props. |
| pages | `ContextVisualizerPage`, `TranscriptAnnotationPage`, `ContextCoursePage` | Static/storyable package demos or web pages composed from organisms. |
| containers | upload parser, API loaders, persisted annotations, route state | `web`-owned unless generalized later. |
| Widget IR / Goja | `contextWindowStudio(...)`, `annotatedTranscript(...)` recipes | Add after React molecules/organisms stabilize. |

### Component promotion table

| Prototype source | Prototype symbol | Production target | First Storybook states |
|---|---|---|---|
| `patterns.jsx` | `KIND`, `TONE`, `OUTLINE`, `resolveKind` | package context model + style helpers | all context kinds, unknown kind fallback |
| `patterns.jsx` | `KindRect`, `Legend` | `ContextKindSwatch`, `ContextLegend` | pattern/tone/outline modes |
| `diagrams.jsx` | `BudgetBar` | `ContextBudgetBar` | under budget, near limit, over budget |
| `diagrams.jsx` | `StripDiagram` | `ContextStripDiagram` | dense segments, selected segment, overflow labels |
| `diagrams.jsx` | `StackDiagram` | `ContextStackDiagram` | grouped by kind, many chunks |
| `diagrams.jsx` | `Treemap` | `ContextTreemap` | mixed token weights, selected segment |
| `diagrams.jsx` | `DiagramRenderer` | `ContextDiagramPanel` or `ContextDiagramRenderer` | switch strip/stack/budget/treemap |
| `screens2.jsx` | `Transcript` | `TranscriptReaderPanel`, `TranscriptMessageCard`, `AnnotationRailPanel` | no annotation, selected annotation, tool/result-heavy transcript |
| `screens3.jsx` | `Comments` | `AnchoredCommentRail`, `AnchoredCommentCard` | rail, sticky, popover-style states |
| `screens.jsx` | `Visualize` | package organism plus web page/container | static fixture first; upload later in web |
| `screens.jsx` / `data2.jsx` | landing/slides | `ContextCoursePage` / `ContextSlideViewer` candidates | lesson overview, active slide, exercise slide |
| `screens2.jsx` | `Handout`, `Markdown` | defer or use existing markdown renderer if present | handout reading only after markdown strategy is decided |
| `ds.jsx` | `Chip`, `SectionLabel` | package atom/foundation only if needed by multiple components | neutral/accent/success/danger; section headers |
| `ds.jsx` | `MenuBar`, `MacWindow` | defer; page-specific chrome, not a global primitive yet | only if studio shell needs it |

### Data contracts to define first

Keep DTOs JSON-compatible so they can later pass through Widget IR and Goja recipes.

```ts
export type ContextPartKind =
  | 'system'
  | 'instruction'
  | 'conversation'
  | 'retrieval'
  | 'tool'
  | 'result'
  | 'annotation'
  | 'course'
  | 'active'
  | 'other';

export interface ContextWindowPart {
  id: string;
  label: string;
  kind: ContextPartKind;
  tokens: number;
  sourceId?: string;
  startToken?: number;
  endToken?: number;
  metadata?: Record<string, string | number | boolean | null>;
}

export interface ContextWindowSnapshot {
  id: string;
  title: string;
  limit: number;
  parts: ContextWindowPart[];
  selectedPartId?: string;
}

export interface TranscriptMessage {
  id: string;
  role: 'system' | 'user' | 'assistant' | 'tool' | 'developer' | 'other';
  text: string;
  tokens?: number;
  timestamp?: string;
  metadata?: Record<string, string | number | boolean | null>;
  annotationIds?: string[];
}

export interface TranscriptAnnotation {
  id: string;
  targetMessageId: string;
  kind: ContextPartKind;
  label: string;
  note?: string;
  confidence?: number;
}
```

## Design Decisions

### Decision 1: Package-first reusable visual vocabulary

**Status:** accepted

The completed design-system unification makes the package the canonical reusable UI layer. Context diagrams, transcript cards, annotation rails, and shared course panels should not start in `web` and later be copied into the package. They should start in the package with fixtures and stories.

### Decision 2: Web app only after static visual states are stable

**Status:** accepted

The web app should consume package components once the package Storybook proves the visual and interaction states. Backend/API containers can then be thin adapters from real data to package DTOs.

### Decision 3: Widget IR recipes are a later phase

**Status:** accepted

The Widget DSL should expose semantic recipes such as `contextWindowStudio(...)` and `annotatedTranscript(...)`, not low-level SVG or duplicated HTML. Add recipes after the React components and DTOs have stabilized.

### Decision 4: Typography adjustment should be explicit and reviewable

**Status:** proposed

The prototype's larger body/compact/metric/code roles are desirable for transcript/course readability, but they affect the whole system if made default. First implementation should either:

- add a package Storybook typography comparison story, or
- introduce a scoped `readable` page mode/theme class for transcript/course surfaces.

Do not silently change all tokens in the same commit as new components.

## Alternatives Considered

### Copy prototype screens directly into `web`

Rejected. This would reintroduce global CSS, inline styles, and duplicated primitives immediately after removing design-system duplication.

### Implement Widget IR first

Rejected for the first pass. The recovered DSL docs repeatedly show that low-level DSL-first rendering creates visual quality and maintenance problems. React components must come first.

### Keep all context-viewer concepts app-local

Rejected for diagram/transcript primitives. The same visuals are likely needed by Storybook baselines, web pages, standalone package apps, and future Goja recipes.

## Implementation Plan

### Phase 0 — Ticket closure and branch hygiene

- Confirm `RAGEVAL-DESIGN-SYSTEM-UNIFY` is complete. Done: ticket status is `complete`.
- Continue on `task/add-ui-dsl` unless a fresh branch is preferred.
- Keep context-viewer changes in small commits: contracts, first molecules, first organisms/pages, app wiring, Widget IR.

### Phase 1 — Package DTOs and fixtures

Files to add:

```text
packages/rag-evaluation-site/src/context/types.ts
packages/rag-evaluation-site/src/context/fixtures.ts
packages/rag-evaluation-site/src/context/index.ts
packages/rag-evaluation-site/src/context/README.md        # optional if types need prose
```

Tasks:

1. Define `ContextWindowSnapshot`, `ContextWindowPart`, `TranscriptMessage`, `TranscriptAnnotation`, and lightweight course/slide DTOs.
2. Port a small normalized subset of `sources/03-context-viewer-design-iteration/data.jsx` and `data2.jsx` into package fixtures.
3. Export context DTOs from `packages/rag-evaluation-site/src/index.ts`.
4. Validate with `pnpm --dir packages/rag-evaluation-site typecheck`.

### Phase 2 — Context-window diagram package components

Files to add under a package-owned namespace, for example:

```text
packages/rag-evaluation-site/src/components/context/ContextLegend/*
packages/rag-evaluation-site/src/components/context/ContextBudgetBar/*
packages/rag-evaluation-site/src/components/context/ContextStripDiagram/*
packages/rag-evaluation-site/src/components/context/ContextStackDiagram/*
packages/rag-evaluation-site/src/components/context/ContextTreemap/*
packages/rag-evaluation-site/src/components/context/ContextDiagramPanel/*
```

Tasks:

1. Port pattern/kind vocabulary from `patterns.jsx` into typed helpers.
2. Reimplement the diagram renderers as typed React components with CSS Modules where possible.
3. Add stories for each renderer and a combined diagram panel.
4. Run package typecheck/build and Storybook build.
5. Capture package Storybook visual sweep with the existing ticket script pattern.

### Phase 3 — Transcript and annotation package components

Files to add under package namespaces, for example:

```text
packages/rag-evaluation-site/src/components/transcript/TranscriptMessageCard/*
packages/rag-evaluation-site/src/components/transcript/TranscriptReaderPanel/*
packages/rag-evaluation-site/src/components/transcript/AnnotationRailPanel/*
packages/rag-evaluation-site/src/components/transcript/AnchoredCommentCard/*
```

Tasks:

1. Port the annotated transcript interaction model from `screens2.jsx` and comment states from `screens3.jsx`.
2. Keep state controlled via props: `selectedAnnotationId`, `onAnnotationSelect`, `selectedMessageId`, etc.
3. Add Storybook states for selected/unselected annotations, long text, tool-heavy transcript, and empty states.

### Phase 4 — Web page composition

Files likely to add:

```text
web/src/components/pages/ContextViewerPage/*
web/src/components/pages/TranscriptWorkbenchPage/*
web/src/components/pages/ContextCoursePage/*
```

Tasks:

1. Compose package components into app pages using static fixtures first.
2. Add navigation entries only after pages render cleanly in stories.
3. Add stories for web pages using package fixtures.
4. Wire backend/API data only after the static UI is validated.

### Phase 5 — Widget IR / Goja recipes

Files likely to modify:

```text
packages/rag-evaluation-site/src/widgets/ir.ts
packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx
pkg/widgetdsl/module.go
pkg/widgetdsl/*_test.go
```

Tasks:

1. Add semantic Widget IR component types only for stable package components.
2. Add Goja recipes such as `contextWindowStudio` and `annotatedTranscript`.
3. Add Go tests for JSON shape and custom key/selection behavior.
4. Validate with package typecheck and targeted Go tests.

## Immediate Next Coding Step

Start with **Phase 1**:

1. create package context DTOs and fixtures;
2. export them;
3. add a minimal Storybook fixture preview if useful;
4. validate with package typecheck.

This is the smallest safe step because it does not commit to SVG rendering details yet, but it gives every later component a stable contract.

## Validation Commands

Use these after each phase:

```bash
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build
pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-package-storybook-context
pnpm --dir web typecheck
pnpm --dir web build
```

For Widget IR phases:

```bash
go test ./pkg/widgetdsl ./pkg/defaultspa -count=1
```

Avoid the broader Go package set until the unrelated `../scraper/pkg/js/runtime/databases.go` missing-module issue is resolved.

## Open Questions

1. Should readable typography be global or scoped to transcript/course contexts?
2. Should course/slides/handout be product pages in `web`, documentation/demo pages in the package, or both?
3. Which backend source should feed the first real context-window data: minitrace sessions, RAG evaluation runs, uploaded JSON, or new APIs?
4. Do annotations/comments need persistence in the first implementation, or are static/storybook states sufficient for this pass?

## References

- `design-doc/02-context-viewer-design-iteration-analysis.md`
- `sources/index.md`
- `sources/04-context-viewer-source-map.md`
- `sources/03-context-viewer-design-iteration/patterns.jsx`
- `sources/03-context-viewer-design-iteration/diagrams.jsx`
- `sources/03-context-viewer-design-iteration/screens2.jsx`
- `sources/03-context-viewer-design-iteration/screens3.jsx`
- `packages/rag-evaluation-site/src/index.ts`
- `packages/rag-evaluation-site/src/widgets/ir.ts`
- `pkg/widgetdsl/module.go`
