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
RelatedFiles:
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js
      Note: First grammar consumer — sections + master-detail agenda (Step 3
    - Path: packages/rag-evaluation-site/src/components/layout/FieldGrid/FieldGrid.tsx
      Note: New n-column field grid compile target (Step 3
    - Path: pkg/widgetdsl/grammar.go
      Note: The grammar implementation — schema/f/record/collection/section (Step 3
    - Path: pkg/widgetdsl/grammar_test.go
      Note: Expansion tests for the grammar (Step 3)
ExternalSources: []
Summary: 'Chronological diary for the UI-grammar brainstorm ticket: how the CMS admin page was audited, what was measured, and how the design doc was structured.'
LastUpdated: 2026-07-04T14:24:43.477438741-04:00
WhatFor: Trace the evidence and decisions behind design-doc 01.
WhenToUse: When picking up any RAGEVAL-UI-GRAMMAR task or reviewing the audit method.
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

## Step 3: Grammar implemented (phase α+) and the admin CMS page rebuilt on it

Implemented design-doc 02 end to end in one vertical slice: the React compile targets (SectionBlock levels/rule/actions/flush + new FieldGrid), the Go grammar verbs in `pkg/widgetdsl/grammar.go` (data.dsl `schema`/`f.*`/`record`/`collection`/`urlParam`/`formPost`, ui.dsl `section`), the module reorganization (generic primitives promoted to ui.dsl with deprecated aliases left behind), and the first consumer: go-go-course's admin Course CMS page rewritten on the grammar. Measured result on the live page: **5,611 px → 3,496 px, 21 panels → 5, nesting depth 2 → 0**; the agenda went from eight stacked ~260 px item-panels to a ~250 px summary table plus one schema-derived editor selected via `?agenda=`.

### Prompt Context

**User prompt (verbatim):** "ok, implement. commit at appropriate intervals, keep a detailed diary (see skill if you forgot the format)"

**Assistant interpretation:** Implement the doc-02 sketch (grammar verbs in data.dsl/ui.dsl, module reorg, structural widgets, go-go-course adoption) with checkpoint commits and diary steps.

**Inferred user intent:** Land the grammar as working code with the admin CMS page as proof, not another document.

**Commits (code):** rag-evaluation-system `51aca3c` ("feat(layout): SectionBlock levels/rule/actions/flush density; new FieldGrid layout"), `4a70c56` ("feat(widgetdsl): data/ui grammar verbs — schema, f roles, record, collection, section"), `9c91539` ("fix(widgets): don't URL-encode interpolated values in confirm prompts"); go-go-course `319c545` ("feat(cms): rebuild admin Course CMS on the data/ui grammar")

### What I did
- **React targets (51aca3c):** SectionBlock gained `level 1–3` (heading→label→metadata font roles), `rule` (1px under-label border), `actions` slot (baseline-aligned right), `anchorId` over IR, `density:"flush"`; new `FieldGrid` layout (`repeat(var(--rag-field-columns), minmax(0,1fr))`, collapses at 720px); both wired through ir.ts/defaultRegistry/barrels with stories.
- **Go grammar (4a70c56), `pkg/widgetdsl/grammar.go`:** 12 field roles (`key primary short prose count size measure date status tags media href`); order-preserving `schema()` (goja `Object.Keys()` — a plain map export would scramble columns); `record()` compiling verb `edit` → FormPanel with consecutive gridable fields batched into FieldGrid (2 cols, 3 when ≥3) and prose as stacked textareas, verb `show` → MetadataGrid; `collection()` compiling summary DataTable (prose/media elided, key→muted caption cell, status→StatusText cell, numerics→number cells) + optional master-detail editor + `create`/`reorder`/`remove` bindings + URL-param selection (`onRowSelect: navigate("?param=${row.key}")`, `"__new"` sentinel for the empty editor); `ui.section()` → flat SectionBlock. Five grammar tests in `grammar_test.go`; TS declarations and a JS API reference grammar chapter.
- **Module reorg (same commit):** `tag meterBar pagination searchField emptyState tileGrid breadcrumbs fieldGrid markdownArticle richArticle uploadDropArea` now in uiHelpers; old module-local names untouched (aliases by construction since helper maps are per-module).
- **go-go-course adoption (319c545):** admin-course-cms.js rebuilt — six `ui.section` blocks; agenda = `data.collection({verb:"edit", arrange:"master-detail", select: urlParam("agenda", query.agenda), submit: formPost("/settings/agenda-item"), reorder, remove-with-confirm, create})`; metadata form batched with fieldGrid; uploads via promoted `ui.uploadDropArea`; media library re-wrapped in a section. Server: `/settings/agenda-item` (upsert via new `upsertAgendaItem`), `admin-delete-agenda-item`, id-based reorder path in the existing reorder action; `course-pages` exposes `currentCourseContent()`.
- **Smoke (Playwright on the devctl prod build):** row click → `?agenda=agenda-regenerate` + populated editor with read-only key; title edit + save → redirect `?status=agenda-item-saved`, table updated; ↓ reorder swapped rows 1/2; "New agenda item" → `__new` empty editor → save → generated id `agenda-18h00-smoke-item`; delete with confirm; cancel path leaves row. Screenshot `sources/screenshots/gg-course-admin-cms-grammar.png`; deleted `course/course-metadata.json` afterwards to restore repo defaults.

### Why
- Doc 02 phase α with the two structural widgets pulled forward (they were cheap and the flat-section look is the visible payoff); single-page-with-anchors instead of the sub-page split so the grammar demonstration wasn't entangled with navigation churn.

### What worked
- The compile-to-existing-IR strategy: the whole grammar shipped without touching WidgetRenderer, and the first real page worked on the second build (the only iteration was a test expectation).
- `Object.Keys()` for schema order; `?param=value` relative navigation via `history.pushState` keeps path and swaps query exactly as needed.

### What didn't work
- **Confirm prompts URL-encoded their interpolations** — the delete dialog read `Delete agenda item “Smoke%20item”?`. `interpolate()` was written for navigate targets and encoded everything; Phase 6 never noticed because filenames can't contain spaces. Fixed with an `{ encode: false }` mode used only by the confirm gate (9c91539) and re-verified: the dialog now renders unicode/spaces raw.
- First test run: `section()` wrapped a string title into a `{kind:"text"}` node via `exportRenderable`, breaking `label` prop equality; strings now pass through as strings.

### What I learned
- goja preserves JS property insertion order through `Object.Keys()` but not through `Export()` to `map[string]any` — any DSL API that depends on field order must capture it before export.
- The elision rule (prose/media skipped from summary tables) is what makes master-detail dramatic: the agenda summary is 4 visible columns of an 5-field schema and instantly scannable.

### What was tricky to build
- **Selection state without client state:** the grammar needs the *current* URL param value at page-build time, but the Go DSL never sees the request — so `urlParam(name, value)` takes both, and the page author feeds `query.agenda` through. Symptom if you forget: table renders but nothing selects. Documented in the API reference example.
- **Per-record saves against a whole-list store:** the agenda lives as one array in a JSON override file that may not exist (defaults come from courseDefinition()). The new service functions take the *effective* agenda as input (`currentCourseContent().agenda`) rather than reading the override, so first-ever saves work; `upsertAgendaItem` matches by slugified id and appends when unmatched, which is also how `__new` creation falls out for free.

### What warrants a second pair of eyes
- `grammar.go` `collectionDetail`: row lookup compares `anyToString(row[key])` — numeric ids will match their string forms, but composite/duplicate keys silently pick the first match.
- The reorder contract (grammar sends `payload.direction` + `context.row`; server resolves index by id) — concurrent edits between render and click can move the wrong neighbor; acceptable for a single admin, worth a version check someday.
- Promoted-alias policy: cms.dsl/course.dsl/context_window.dsl still export the moved helpers with no deprecation warning at runtime — only docs say so.

### What should be done in the future
- Sub-page split (ui.subpages) — task 3 still open; arrange "tiles"/marks contract (doc 02 §4.3) and `data.scale.*` promotion (§4.4) unimplemented; outcomes still a padded-rows form (could become `collection` with a string schema); Storybook grammar-gallery story; selected-row visual check in DataTable.

### Code review instructions
- Read `pkg/widgetdsl/grammar.go` top-to-bottom (430 lines) with `grammar_test.go` beside it; then the page diff `go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js` and the three service functions in `course-metadata-service.js`.
- Validate: `go test ./pkg/widgetdsl/`; `pnpm --dir packages/rag-evaluation-site typecheck && pnpm --dir packages/rag-evaluation-site build-storybook`; `cd go-go-course && devctl up --profile prod` → `/pages/admin-course-cms` as `admin_<you>` → click an agenda row, edit, save, reorder, delete (note: exercising saves creates `course/course-metadata.json` — delete it to restore defaults).
- Storybook: "Design System/Layout/SectionBlock" (Levels/WithRule/WithActions/Flush) and "Design System/Layout/FieldGrid".

### Technical details
- Metrics before → after (same DOM-metrics snippet as Step 1): height 5,611→3,496 px; panels 21→5 (all top-level tools: 2 FormPanels, 2 upload areas' inner chrome, asset detail when selected); FieldGrids 3; agenda table 7 rows.
- Biome full-package check reports 20 findings with or without these changes (pre-existing a11y/story lint in untouched files); new files are clean.
