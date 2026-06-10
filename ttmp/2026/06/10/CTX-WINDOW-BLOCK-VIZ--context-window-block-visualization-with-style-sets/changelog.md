# Changelog

## 2026-06-10

- Initial workspace created


## 2026-06-10

Created intern-facing design guide and diary for uploaded-session context-window block visualization; no production code changed.

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/10/CTX-WINDOW-BLOCK-VIZ--context-window-block-visualization-with-style-sets/design-doc/01-design-and-implementation-guide.md — Design and implementation guide
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/10/CTX-WINDOW-BLOCK-VIZ--context-window-block-visualization-with-style-sets/reference/01-investigation-diary.md — Chronological investigation diary


## 2026-06-10

Validated CTX-WINDOW-BLOCK-VIZ and uploaded the design guide bundle to reMarkable at /ai/2026/06/10/CTX-WINDOW-BLOCK-VIZ.

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/10/CTX-WINDOW-BLOCK-VIZ--context-window-block-visualization-with-style-sets/reference/01-investigation-diary.md — Records validation and reMarkable upload evidence


## 2026-06-10

Removed live/course workspace surface from active navigation and Storybook, keeping Upload / Visualize / Transcript; commits 560864d and 2edac8e.

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/organisms/CourseStudioShell/courseStudioNav.ts — Default nav scope
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/widgets/WidgetRenderer.course-handout.stories.tsx — Session workspace stories


## 2026-06-10

Updated design/package documentation to mark live/page-through session views out of current scope and describe the active upload → visualize/transcript flow.

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/README.md — Package README scope clarification
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/10/CTX-WINDOW-BLOCK-VIZ--context-window-block-visualization-with-style-sets/design-doc/01-design-and-implementation-guide.md — Scope clarification


## 2026-06-10

Added phased cleanup and widget-improvement task plan before starting further work.

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/10/CTX-WINDOW-BLOCK-VIZ--context-window-block-visualization-with-style-sets/tasks.md — Phased implementation task plan


## 2026-06-10

Expanded the design guide with concrete go-minitrace schema and minitracejs DSL row mappings before code cleanup.

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/10/CTX-WINDOW-BLOCK-VIZ--context-window-block-visualization-with-style-sets/design-doc/01-design-and-implementation-guide.md — go-minitrace mapping update
- /home/manuel/workspaces/2026-06-07/club-meetup-site/go-minitrace/pkg/minitracejs/query_view_session.go — DSL row output source


## 2026-06-10

Restored course/slides/handout workspace surfaces after narrowing cleanup scope to only LiteLLM live service; commit 5cd3fb1.

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/organisms/CourseStudioShell/courseStudioNav.ts — Restored default course studio nav
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/widgets/WidgetRenderer.course-handout.stories.tsx — Restored course/handout Storybook surface


## 2026-06-10

Step 6: removed JS LiteLLM live-service integration, fixed WidgetRenderer child keys found by Playwright smoke, and captured Storybook/ClubMed visual evidence

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx — Fixed child key warning found by Playwright smoke
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/10/CTX-WINDOW-BLOCK-VIZ--context-window-block-visualization-with-style-sets/scripts/01-capture-clubmed-visual-smoke.js — Reusable css-visual-diff smoke capture script
- /home/manuel/workspaces/2026-06-07/club-meetup-site/ClubMedMeetup/minitrace-viz/lib/course-pages.js — Removed live-litellm nav/page while preserving course/slides/handouts
- /home/manuel/workspaces/2026-06-07/club-meetup-site/ClubMedMeetup/minitrace-viz/server.js — Removed direct LiteLLM live routes and actions


## 2026-06-10

Step 6 committed: visual-smoke/key fix in RAG repo (ace0505) and LiteLLM live-service removal in ClubMed repo (49fd9fa)

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx — Committed keyed child rendering fix (ace0505)
- /home/manuel/workspaces/2026-06-07/club-meetup-site/ClubMedMeetup/minitrace-viz/lib/litellm-live-service.js — Deleted JS LiteLLM live service (49fd9fa)


## 2026-06-10

Step 7: implemented local minitrace-viz uploaded-session context block mapping with user/agent/thinking/tool/file style keys (ClubMed commit 58ee4d3)

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/ClubMedMeetup/minitrace-viz/lib/course-session-data.js — Maps timeline/tool/file rows into ContextWindowPart blocks
- /home/manuel/workspaces/2026-06-07/club-meetup-site/ClubMedMeetup/minitrace-viz/lib/timeline-data.js — Adds original-upload fallback rows when normalized minitrace rows are empty
- /home/manuel/workspaces/2026-06-07/club-meetup-site/ClubMedMeetup/minitrace-viz/test-fixtures/smoke-test.sh — Asserts uploaded-session style keys in smoke test


## 2026-06-10

Step 8: closed Phase 3 and Phase 4 as no-change decisions after package Storybook/typecheck and focused DSL/schema Go tests

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/GUIDELINES.md — Reviewed before deciding no upstream component changes were needed
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/pkg/widgetdsl/module.go — Existing DSL helpers were sufficient; no changes made


## 2026-06-10

Step 9: completed Phase 5 final validation with uploaded-session API check, Playwright smoke, css-visual-diff evidence, and doc updates

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/10/CTX-WINDOW-BLOCK-VIZ--context-window-block-visualization-with-style-sets/design-doc/01-design-and-implementation-guide.md — Implementation outcome and final validation summary
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/10/CTX-WINDOW-BLOCK-VIZ--context-window-block-visualization-with-style-sets/sources/visual-smoke/run-phase5-20260610-160514/01-visual-smoke-summary.md — Final css-visual-diff evidence summary


## 2026-06-10

Step 10: removed noisy derived annotations, preserving only model-change and explicit compaction annotations (ClubMed commit d38d877)

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/ClubMedMeetup/minitrace-viz/lib/course-session-data.js — Derived annotation filtering and context-note cleanup


## 2026-06-10

Step 11: added ContextGroupedStripDiagram, ContextTurnPagerPanel, Widget IR stories, and fast CSS strip tooltips (commit 982904b)

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/ContextGroupedStripDiagram/ContextGroupedStripDiagram.tsx — Grouped strip visualization
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/ContextStripDiagram/ContextStripDiagram.module.css — Fast CSS tooltip behavior
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/organisms/ContextTurnPagerPanel/ContextTurnPagerPanel.tsx — Turn pager panel


## 2026-06-10

Step 12: changed ContextTurnPagerPanel to default to turn-only context-window slices instead of selection-only paging (commit 59dd5e9)

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/organisms/ContextTurnPagerPanel/ContextTurnPagerPanel.tsx — Turn-only snapshot filtering and headroom recomputation

