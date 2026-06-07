---
title: "Investigation Diary"
ticket: RAGEVAL-DS-RAGEVAL-QWEN36-HIGH
topics: [frontend-architecture, design-system, react, storybook]
doc-type: reference
intent: long-term
---

# Investigation Diary — RAGEVAL-DS-RAGEVAL-QWEN36-HIGH

## Goal

Create two comprehensive reference documents for the RAG Evaluation System design system:
1. **Styling and Design Reference** — visual language, tokens, CSS architecture, component anatomy
2. **Widget Hierarchy and Interaction Reference** — composition model, information hierarchy,
   Widget IR, Goja DSL, WidgetRenderer, typical page patterns

## Step 1: Evidence Gathering

### What I did
- Read all three referenced articles:
  - `ARTICLE - RAG React Design System - From Prototype Dashboard to Structured Design System.md` (2026-06-02)
  - `ARTICLE - Building a Goja UI DSL from Scratch - Widget IR to xgoja.md` (2026-06-05)
  - `ARTICLE - WidgetRenderer Standalone Site - Goja Authored React Rendered UI.md` (2026-06-05)
- Created docmgr ticket `RAGEVAL-DS-RAGEVAL-QWEN36-HIGH`
- Explored the web directory: `web/src/components/` (all 45+ component files)
- Read key source files across all layers
- Examined Storybook stories for naming conventions and coverage patterns
- Read tokens.css, index.css, ir.ts, WidgetRenderer.tsx, cellRenderers.tsx

### Why
The reference documents need to be evidence-based, not speculative. Every claim about the design
system had to be anchored to actual source files and the three documented articles.

### What worked
- Reading all three articles gave a complete picture of the design system's history and architecture
- Exploring the file tree revealed the exact component structure (6 layers + corpus)
- Reading ir.ts + WidgetRenderer.tsx together clarified the IR ↔ React mapping
- Reading DataTable + StatusText + Caption + Panel CSS modules gave precise visual details

### What didn't work
- None. The evidence was straightforward to gather.

### What I learned
- The system has a very clear ownership discipline: each layer owns one type of concern
- StatusText uses 10 specific statuses with Unicode icons and CSS-class-based coloring
- Button has 2 variants, 2 sizes, and a selected state that looks identical to primary
- DashboardGrid has 3 named recipes, not generic props
- CellSpec replaces JavaScript callbacks in tables (data, not functions, cross the boundary)
- Widget IR has 3 node kinds: text, element, component — and "component" maps to named React components

### What was tricky
The main challenge was organizing the information into two clean documents without redundancy.
Both documents share the layer map and ownership rules, but each focuses on a different aspect.
The design decision was to reference across documents rather than duplicate.

### What should be done in the future
- Verify Storybook stories still match the component API (especially for any components added after
  the three source articles were written)
- Add Storybook screenshots to the reference documents if the user wants visual examples
- Consider whether the reference should be served as an interactive HTML page or Storybook
  embedded documentation

## Step 2: Writing the Documents

### What I did
- Created the ticket workspace with `docmgr ticket create-ticket`
- Added three documents:
  - `design-doc/01-styling-and-design-reference.md` — the styling/design document (26KB)
  - `design-doc/02-widget-hierarchy-and-interaction-reference.md` — the widget/interaction document (30KB)
  - `reference/01-investigation-diary.md` — this diary
- Wrote both design docs in full with all sections required by the design-doc guideline

### Why
The user requested two documents split by concern: styling/design on one side,
widget/interaction hierarchy on the other. Each document covers the full lifecycle
from high-level architecture down to specific props, CSS classes, and data attributes.

### What worked
- The layer map served as a common reference in both documents
- Quick reference tables at the end of each document made them useful as look-up
- The CSS variable map (Section 11 of the styling doc) is a compact single-page reference

### What I learned
- Both documents combined cover 56KB of technical reference material
- The widget hierarchy document includes the full Widget IR type system, cell specs,
  action specs, Goja DSL API, WidgetRenderer dispatch table, and 4 canonical page layouts
- The styling document includes the complete token map, component anatomy, CSS architecture rules,
  and a step-by-step guide for adding new components
