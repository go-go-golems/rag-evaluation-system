---
Title: Visual Evidence Capture Summary
DocType: reference
Ticket: WIDGETDSL-VISUAL-QUALITY
Status: active
Intent: short-term
Topics:
  - frontend
  - ui-dsl
  - design-system
---

# Widget DSL Visual Evidence Capture

- Widget URL: http://127.0.0.1:18791
- Storybook URL: http://127.0.0.1:6007
- Viewport: 1440x1100
- Wait: 800ms

## Captured targets

- **Generated Widget DSL action page** (widget-actions) — http://127.0.0.1:18791/pages/actions
  - overlay: `widget-actions.overlay.png`
  - preflight: `widget-actions.preflight.json`
- **Original RAG PipelinePage Storybook page** (storybook-pipeline-page) — http://127.0.0.1:6007/iframe.html?id=pages-pipelinepage--populated&viewMode=story
  - overlay: `storybook-pipeline-page.overlay.png`
  - preflight: `storybook-pipeline-page.preflight.json`
- **Original RetrievalResultsPanel Storybook organism** (storybook-retrieval-results) — http://127.0.0.1:6007/iframe.html?id=component-library-organisms-retrievalresultspanel--hybrid-results&viewMode=story
  - overlay: `storybook-retrieval-results.overlay.png`
  - preflight: `storybook-retrieval-results.preflight.json`
- **WidgetRenderer search workbench Storybook composition** (widget-renderer-storybook) — http://127.0.0.1:6007/iframe.html?id=widget-ir-renderer--search-workbench-composition&viewMode=story
  - overlay: `widget-renderer-storybook.overlay.png`
  - preflight: `widget-renderer-storybook.preflight.json`

## Notes for analysis

Use the overlay screenshots for layout-level review and the CSS JSON artifacts for computed style evidence. Missing optional probes are recorded in the preflight files and should guide selector refinement rather than be treated as failures.