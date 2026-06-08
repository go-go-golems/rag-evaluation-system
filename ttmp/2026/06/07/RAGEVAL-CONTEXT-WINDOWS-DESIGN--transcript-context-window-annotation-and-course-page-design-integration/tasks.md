# Tasks

## TODO

- [x] Add tasks here

- [x] Use go-minitrace to discover prior design-system and web playbooks from Pi sessions
- [x] Extract concrete frontend design-system conventions for web/ contribution
- [x] Map requirements for transcript, context-window, annotation, and course pages
- [x] Write phased integration plan with related files and validation steps
- [x] Add package context DTOs and fixtures
- [x] Define context component taxonomy and export style helpers for context kinds
- [x] Add ContextKindSwatch atom with Storybook states
- [x] Add ContextLegend molecule with pattern/tone/outline stories
- [x] Add ContextBudgetBar molecule with under/near/over-budget stories
- [x] Add ContextStripDiagram molecule with dense and selected segment stories
- [x] Add ContextStackDiagram molecule with grouped context-window stories
- [x] Add ContextTreemap molecule with proportional token stories
- [x] Compose ContextDiagramPanel organism for switching diagram views
- [x] Run package Storybook build and capture visual sweep for context diagram components
- [x] Ask for human Storybook feedback on context diagram visuals before transcript work
- [x] Add TranscriptMessageCard atom/molecule split and stories
- [x] Add AnnotationBadge/AnnotationNoteCard molecules and stories
- [x] Add TranscriptReaderPanel and AnnotationRailPanel organisms with controlled selection stories
- [x] Add AnchoredCommentCard/Rail components and stories from prototype comment states
- [x] Add CourseStepNav/CourseLessonPanel/CourseSlidePanel components and stories
- [x] Add generic landing-page vocabulary for course surfaces (SectionBlock, SplitPane, KeyValueStrip, CheckList, StepList, PersonSummary)
- [x] Add generic slide-shell vocabulary and refactor CourseSlidePanel to support visual-left or visual-right compositions
- [x] Add global course sidebar shell for Course/Slides/Visualize/Upload/Transcript/Comments/Handout navigation
- [ ] Compose web ContextVisualizerPage from package components using fixtures
- [ ] Compose web TranscriptAnnotationPage from package components using fixtures
- [ ] Compose web ContextCoursePage from package components using fixtures
- [ ] Wire optional web navigation entries and route/page state after page stories pass
- [ ] Add Widget IR component types for stable context components
- [ ] Add Goja widget.dsl recipes for contextWindowStudio and annotatedTranscript
- [ ] Run final package/web typecheck, builds, Storybook builds, visual sweeps, and ticket closeout

## Phase 2: Prototype-to-Design Visual Comparison

- [x] Generate standalone HTML pages from prototype JSX (screens.jsx, screens2.jsx, screens3.jsx)
- [x] Add data-* annotations to prototype widgets for css-visual-diff section matching
- [x] Create css-visual-diff jsverb capture workflow for prototype pages/widgets
- [x] Capture prototype baseline screenshots with css-visual-diff browser runtime
- [ ] Run css-visual-diff comparison: prototype vs package components
- [ ] Document visual parity gaps and tuning priorities
