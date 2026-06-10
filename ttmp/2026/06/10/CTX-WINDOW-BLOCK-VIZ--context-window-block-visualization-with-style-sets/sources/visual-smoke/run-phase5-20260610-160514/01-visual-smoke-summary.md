---
Title: ClubMed Visual Smoke Summary
DocType: reference
Ticket: CTX-WINDOW-BLOCK-VIZ
Status: active
Intent: short-term
Topics:
  - frontend
  - visual-testing
---

# ClubMed Visual Smoke Summary

- App URL: http://127.0.0.1:18789
- Storybook URL: http://127.0.0.1:6009
- Viewport: 1440x1100
- Wait: 1000ms
- Missing required probes: none

## Captured targets

- **ClubMed course landing** (clubmed-course) — http://127.0.0.1:18789/pages/course
  - overlay: `clubmed-course.overlay.png`
  - preflight: `clubmed-course.preflight.json`
- **ClubMed upload workflow** (clubmed-upload) — http://127.0.0.1:18789/pages/upload
  - overlay: `clubmed-upload.overlay.png`
  - preflight: `clubmed-upload.preflight.json`
- **ClubMed restored slides page** (clubmed-slides) — http://127.0.0.1:18789/pages/slides
  - overlay: `clubmed-slides.overlay.png`
  - preflight: `clubmed-slides.preflight.json`
- **ClubMed restored handouts page** (clubmed-handouts) — http://127.0.0.1:18789/pages/handouts
  - overlay: `clubmed-handouts.overlay.png`
  - preflight: `clubmed-handouts.preflight.json`
- **Storybook context diagram custom legend** (storybook-context-diagram) — http://127.0.0.1:6009/iframe.html?id=widget-ir-renderer-context-diagrams--custom-three-label-widget-ir&viewMode=story
  - overlay: `storybook-context-diagram.overlay.png`
  - preflight: `storybook-context-diagram.preflight.json`
- **Storybook transcript with notes rail** (storybook-transcript-notes) — http://127.0.0.1:6009/iframe.html?id=widget-ir-renderer-transcript-and-notes--annotated-transcript-with-notes-rail&viewMode=story
  - overlay: `storybook-transcript-notes.overlay.png`
  - preflight: `storybook-transcript-notes.preflight.json`
- **Storybook course and handout registry surface** (storybook-course-handouts) — http://127.0.0.1:6009/iframe.html?id=widget-ir-renderer-domain-registry-coverage--course-handout-registry-surface&viewMode=story
  - overlay: `storybook-course-handouts.overlay.png`
  - preflight: `storybook-course-handouts.preflight.json`