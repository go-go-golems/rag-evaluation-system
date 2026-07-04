---
Title: Investigation diary
Ticket: RAGEVAL-UI-GRAMMAR
Status: active
Topics:
    - cms
    - design-system
    - frontend
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Chronological diary for the UI-grammar brainstorm ticket: how the CMS admin page was audited, what was measured, and how the design doc was structured."
LastUpdated: 2026-07-04T14:24:43.477438741-04:00
WhatFor: "Trace the evidence and decisions behind design-doc 01."
WhenToUse: "When picking up any RAGEVAL-UI-GRAMMAR task or reviewing the audit method."
---

# Diary

## Goal

Capture how the CMS-admin-page audit and the cross-DSL UI-grammar brainstorm (design-doc 01) were produced, so the numbers and screenshots can be regenerated and the method reapplied to other pages.

## Step 1: Audit the admin CMS page and write the brainstorm design doc

Created this ticket while RAGEVAL-CMS-WIDGETS Phase 6 was landing: with the media library, confirm-deletes, and handout editor freshly smoke-tested, the admin page was fully functional and visibly hard to read — the right moment to audit it. Measured the live page with DOM inspection, captured whole-page and per-section screenshots, and wrote design-doc 01: a section-by-section audit, a diagnosis (the DSLs expose components, not intent — panel is the only sectioning device, there is no collection primitive), and brainstormed alternatives up to a grammar-of-graphics-style composable UI grammar spanning all five DSL modules. Explicitly no implementation.

### Prompt Context

**User prompt (verbatim):** "Can you look at the course CMS page and its many sections individually too, and make a new ticket to make the UI easier to read, easier to manage, how to deal with these long content lists, etc... This CMS UI is good in the sense that it composed many existing smaller primitives, but it could need better primitives and an overhault. Create a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable." — followed by: "for exampel, the nested boxes for the sections are confusing, and I think that sections wuld be better served with a title, and then tables / forms. Also dynamically expanding lists or splitting things into more subscetions through the navigation. Don't implement, but in the ticket, brainstorm different alternatives. The core concept is to have DSL languages that allow to do a composoble 'UI design system grammar' that allows combining primitives and richer widgets in a composable manner (think 'grammar of graphics')" — and: "these improvements would apply across all the DSLs, btw."

**Assistant interpretation:** New ticket; per-section UI audit of /pages/admin-course-cms; intern-grade analysis document that brainstorms alternatives (flat sections, navigation splits, expanding lists) without implementing, centered on a composable grammar-of-graphics-style DSL layer applicable to all five widget DSL modules; upload to reMarkable.

**Inferred user intent:** Build the case and the conceptual frame for the next generation of the widget DSLs — an intent-level authoring grammar — before spending any implementation effort.

### What I did
- Ran the freshly built go-go-course binary via the new devctl prod profile; walked `/pages/admin-course-cms` as `admin_manuel`.
- Measured via Playwright `evaluate`: document height 5,611 px; 8 top-level / 21 total panels; per-section height, nesting depth, form-row/input/table/button counts (table in design-doc §2). Agenda: 2,158 px, 8 nested item panels, 40 form rows, 16 reorder buttons.
- Tagged top-level panels with `data-shot` attributes to take stable per-section element screenshots: `section-agenda.png`, `section-outcomes.png`, `section-material-tables.png`, `section-media-library.png`, plus the full page; stored under `sources/screenshots/`.
- Read `admin-course-cms.js`/`admin-common.js` against the screenshots to attribute each visual problem to the authoring construct that produced it.
- Wrote design-doc 01 (audit → diagnosis → alternatives A–F → acceptance sketch → references → open questions) and seeded tasks 2–7.

### Why
- The user asked for evaluation material, not code: alternatives must be argued from measured evidence and mapped to the vocabulary gap, otherwise the "overhaul" degenerates into restyling panels.

### What worked
- The token-cheap audit loop: one `evaluate` call returns all section metrics; `data-shot` tagging makes element screenshots deterministic. Reapply this to the other DSL pages (task 2).

### What didn't work
- `data-rag-*` values are PascalCase component names (`data-rag-layout="Panel"`), not kebab-case — first selector probes returned nothing. Documented in the CMS ticket diary too (Step 9).

### What I learned
- The failure is grammatical, not stylistic: tables and tile grids on the same page read fine; only the hand-unrolled collections and panel-as-section usages degrade. Recipes (`masterDetailTable`, `mediaLibrary`, `articleList`) already prove the intent-layer idea works — they are just frozen sentences instead of a productive grammar.

### What was tricky to build
- Keeping the doc a brainstorm: every alternative invites an API design. Resolved by confining API material to explicitly-labeled pseudocode sketches and pushing all commitment into open questions + tasks.

### What warrants a second pair of eyes
- The grammar table in §4E (GoG layer ↔ UI layer mapping) — the axes chosen (schema/shaping/arrangement/composition/verbs) determine the whole future API; challenge them before anything gets built.
- §4C's recommendation (master-detail for agenda, editable table for outcomes) is an opinion formed from one page.

### What should be done in the future
- Tasks 2–7 in tasks.md, starting with auditing the other DSL pages to confirm the diagnosis generalizes.

### Code review instructions
- Read design-doc 01 top to bottom with `sources/screenshots/section-agenda.png` open beside §2.
- Regenerate numbers: `devctl up --profile prod` in go-go-course, browse as admin, run the §2 evaluate snippet from this ticket's design doc method (documented in diary Step 1 above).

### Technical details
- Metrics collection: `[...document.querySelectorAll('main [data-rag-layout="Panel"]')].filter(p => !p.parentElement.closest('[data-rag-layout="Panel"]'))` → per-panel `getBoundingClientRect().height`, descendant counts for FormRow/textarea/input/table/button, nesting depth via `closest` walking.

## Step 2: Cross-page audit + DSL API sketch (design-doc 02) — no grammar.dsl, grow data.dsl/ui.dsl

Answered the user's follow-up: audited the remaining DSL pages to confirm the diagnosis generalizes (task 2), then wrote design-doc 02 — the concrete API sketch (task 5). The central call: `grammar.dsl` should not exist. The data grammar (`schema`/`f.*`/`record`/`collection`) belongs in `data.dsl`, the structure grammar (`section`/`subpages`/`disclosure`) in `ui.dsl`, and the three domain modules shrink to schemas + marks, with stranded generic primitives (tag, pagination, searchField, tileGrid, uploadDropArea, markdownArticle…) migrating to ui.dsl under deprecation aliases.

### Prompt Context

**User prompt (verbatim):** "also do 1b once you are well informe enough to do a sketch of the nupdate DSL apis (what's grammar.dsl? can't we put this into ui.dsl or somehow reorganize the other existing DSLs, so we have something coherent?)"

**Assistant interpretation:** First get informed (audit the other pages — form answer 1a), then sketch the updated DSL APIs (1b), treating "grammar.dsl" as an open question — prefer folding the grammar into the existing modules / reorganizing them coherently over adding a sixth module.

**Inferred user intent:** A reviewable API surface proposal that makes the whole DSL family coherent, not a bolt-on module.

### What I did
- Audited `/pages/{course,sessions,admin-course-material,handouts?doc,upload}` with the doc-01 DOM-metrics method (server still up via devctl prod). Result table in doc 02 §1.
- Read `pkg/widgetdsl/module.go` end to end for the module anatomy: `moduleSpec` {helpers, page, cell, action, recipes}, five helper maps, shared exports duplicated across modules.
- Wrote design-doc 02: audit evidence → current-anatomy table with incoherences → target layout (data.dsl = data grammar, ui.dsl = structure grammar, domain modules = schemas + marks) → moves/deprecations table → API sketch (`data.schema`/`f.*` roles, `data.collection` with verb/arrange/select/submit/reorder/remove, `data.record`, marks contract, `ui.section`, `ui.subpages`) → three worked page rewrites (agenda, sessions, media library) → hybrid compilation phases α/β → open questions.
- Cross-linked doc 01 §4E to doc 02; checked tasks 2 and 5.

### Why
- The audit had to precede the sketch: it produced the key law — every page that reads well is one recipe/shell or a flat document; every degraded page is hand-assembled panels around collections. Length is innocent (handout doc: 5,043 px, 0 panels, reads fine); boxed repetition is guilty. That law justifies "intent in, layout owned by the system" as the API principle.

### What worked
- Sketching against three real pages immediately killed the separate-module idea: the media library sentence needed cms marks + data verbs + ui structure in one expression, which reads naturally as `data.collection(…, {arrange: cms.marks.assetTiles})` and would read terribly split across a sixth module.

### What didn't work
- N/A (analysis/writing step; no code changed).

### What I learned
- `data.dsl` is nearly empty today (`dataTable` + `cell`) — it is the natural landing zone for the data grammar precisely because `cell.*` is already a proto-role system trapped inside tables.
- The helper-map architecture makes the reorganization mechanically trivial (move a line between maps, keep an alias); the cost is entirely in docs/declarations/deprecation discipline.

### What was tricky to build
- Keeping the sketch honest about working assumptions: sub-pages, master-detail, hybrid compilation and button-reorder were the user's quick-form defaults, not decisions — each is labeled with its assumption in doc 02 so a changed answer invalidates a paragraph, not the document.

### What warrants a second pair of eyes
- The field-role set (`key/primary/short/prose/…`) and the verb set (`show/edit/pick/manage`) — these are the grammar's axes; wrong axes here are expensive later (doc 02 §7 naming/nesting questions too).
- Whether recipes-as-wrappers really covers all ten existing recipes (checked mentally for masterDetailTable/mediaLibrary/articleList only).

### What should be done in the future
- Remaining tasks 3, 4, 6, 7 — wayfinding decision is assumed (sub-pages) but not ratified; Storybook prototype of Section/FieldGrid; agenda pattern decision; then the implementation-plan ticket.

### Code review instructions
- Read doc 02 §0 and §5.3 first (the thesis and its acceptance test), then §3–§4 for the surface.
- Verify the anatomy claims against `pkg/widgetdsl/module.go` L22–163 and `install()` L217–241.

### Technical details
- Audit metrics snippet documented in Step 1; pages measured this step: course 1048px/0 panels, sessions 800px/2, admin-course-material 1567px/8 (depth 1), handouts?doc 5043px/0, upload 800px/1.
