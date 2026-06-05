---
Title: Investigation Diary
DocType: reference
Ticket: WIDGETDSL-VISUAL-QUALITY
Status: active
Intent: long-term
Topics:
  - frontend
  - widget-dsl
  - visual-design
Created: 2026-06-05T12:10:00-04:00
LastUpdated: 2026-06-05T12:10:00-04:00
---

# Diary

## Goal

Capture the evidence-gathering and design process for improving the visual quality of Widget DSL output so DSL-authored sites can become rich, polished web experiences with little author code.

## Step 1: Create the Visual Quality Ticket and Define the Investigation

This step created a dedicated docmgr workspace for visual-quality research instead of mixing the work into the existing xgoja implementation ticket. The request is not just to tweak one demo page; it asks for a technical analysis, a clear design and implementation guide, visual comparisons, reusable scripts, and a reMarkable-ready deliverable.

The ticket will collect all scripts under its `scripts/` directory, keep screenshots and extracted evidence under ticket-local folders, and use the diary to preserve each command, failure, and conclusion. The intended reader is a new intern who needs enough architecture, API, visual-system, and workflow context to make concrete improvements without rediscovering the whole Widget DSL stack.

### Prompt Context

**User prompt (verbatim):** "ngl, look at the screenshots, this looks a bit like ass compared to the original rag site itself (see storybook, which you can launch and screenshot too, potentially use the tool `css-visual-diff help --all` to write some custom snreenshot / comparison code and jsverbs). Save all scripts in a enew ticket which is about improving the visual look of the widget dsl, and creating solid websites. The point of the DSL is creating rich web experiences that look and work great, without writing a lt of code. 

Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**User prompt (verbatim):** "keep a frequent and detailed diary as you work so we can retrace your steps."

**Assistant interpretation:** Create a new docmgr ticket focused on Widget DSL visual quality, gather visual evidence from the current DSL/demo and original RAG Storybook components, use screenshot/comparison tooling where useful, save scripts in the ticket, write a detailed intern-facing implementation guide, and upload the final bundle to reMarkable.

**Inferred user intent:** Preserve the low-code value of the DSL while raising the resulting site quality so DSL-authored pages look and behave like first-class RAG web experiences instead of crude generated layouts.

**Commit (code):** N/A — documentation/research setup only.

### What I did

- Created docmgr ticket `WIDGETDSL-VISUAL-QUALITY`.
- Created design document `design-doc/01-widget-dsl-visual-quality-analysis-and-implementation-guide.md`.
- Created diary document `reference/01-investigation-diary.md`.
- Added initial tasks for screenshot evidence, scripts, gap analysis, intern guide writing, validation, and reMarkable upload.
- Removed an accidental empty `2026/06/02` ticket stub after noticing docmgr created the real ticket under `2026/06/05`.

### Why

- The work needs isolated artifacts because it includes research scripts, visual evidence, design guidance, and delivery steps.
- A new intern needs a stable entry point: not only screenshots, but also the surrounding architecture and implementation sequence.
- Keeping a detailed diary from the beginning makes later review and continuation possible.

### What worked

- `docmgr ticket create-ticket` created the ticket workspace under `ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design`.
- `docmgr doc add` created the design doc and diary doc.
- The initial task list was added successfully.

### What didn't work

- I initially ran `mkdir -p` against a `2026/06/02/...` path from muscle memory. The real ticket was created under `2026/06/05/...`.
- Fix: removed the empty accidental stub and continued with the actual ticket path.

### What I learned

- For new docmgr tickets, trust the path printed by `docmgr ticket create-ticket` rather than inferring the date path manually.
- This ticket should be treated as a design/research deliverable first, not an implementation patch, because the requested output is an intern-ready guide and visual-quality plan.

### What was tricky to build

- The task spans documentation, visual inspection, browser tooling, and architecture planning. The main risk is collecting screenshots without connecting them to concrete file-level causes.
- To avoid that, the next steps will pair screenshots with file references in the renderer, component library, demo jsverb, and app shell.

### What warrants a second pair of eyes

- Review whether the proposed visual-quality work should result in changes to Widget IR schema, renderer defaults, example authoring conventions, or all three.
- Review whether the intern guide is specific enough to support implementation without hand-holding.

### What should be done in the future

- Gather current widget-site screenshots and original RAG/Storybook screenshots.
- Use `css-visual-diff help --all` and create ticket-local screenshot/comparison scripts.
- Write the final guide and upload it to reMarkable.

### Code review instructions

- Review the ticket setup at `ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design`.
- Confirm all future helper scripts live under that ticket's `scripts/` directory.
- Validate ticket hygiene with `docmgr doctor --ticket WIDGETDSL-VISUAL-QUALITY --stale-after 30`.

### Technical details

- Ticket id: `WIDGETDSL-VISUAL-QUALITY`
- Ticket path: `ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design`
- Primary design doc: `design-doc/01-widget-dsl-visual-quality-analysis-and-implementation-guide.md`
- Diary doc: `reference/01-investigation-diary.md`
