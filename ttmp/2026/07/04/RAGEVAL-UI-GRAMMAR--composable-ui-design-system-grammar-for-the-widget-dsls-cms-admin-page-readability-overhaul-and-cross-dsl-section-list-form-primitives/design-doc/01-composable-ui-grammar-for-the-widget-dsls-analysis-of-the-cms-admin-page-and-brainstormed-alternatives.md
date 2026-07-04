---
Title: Composable UI grammar for the Widget DSLs — analysis of the CMS admin page and brainstormed alternatives
Ticket: RAGEVAL-UI-GRAMMAR
Status: active
Topics:
    - cms
    - design-system
    - frontend
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js
      Note: The audited page — every problem in §2 traces to a construct here
    - Path: packages/rag-evaluation-site/GUIDELINES.md
      Note: Constraints any new primitive from §4F must obey
    - Path: packages/rag-evaluation-site/src/widgets/ir.ts
      Note: IR the grammar would compile to
    - Path: pkg/widgetdsl/module.go
      Note: The DSL vocabulary whose gaps §3 diagnoses; recipes are the proto-grammar
ExternalSources: []
Summary: 'Section-by-section audit of the go-go-course admin CMS page (5,611 px of nested boxes), a diagnosis of why pages authored in the Widget DSLs read poorly at scale, and a brainstorm — explicitly no implementation — of alternatives: flat sections, in-page/sub-page navigation, collection-editing patterns for long lists, and a composable ''UI design-system grammar'' (grammar-of-graphics style) that would apply across all five DSL modules.'
LastUpdated: 2026-07-04T15:00:00-04:00
WhatFor: Give a new intern everything needed to understand the current page-authoring model, see concretely why it degrades on content-heavy admin pages, and evaluate the brainstormed directions before any implementation decision.
WhenToUse: Read before proposing or reviewing any change to the widget DSLs' sectioning, list, or form primitives; pair with the RAGEVAL-CMS-WIDGETS design doc for the underlying architecture.
---


# Composable UI grammar for the Widget DSLs — analysis of the CMS admin page and brainstormed alternatives

## Executive summary

The go-go-course admin CMS page (`/pages/admin-course-cms`) works — every flow on it was smoke-tested end to end in RAGEVAL-CMS-WIDGETS Phase 6 — but it does not read well. It is 5,611 px tall, contains 21 bordered panels (up to two levels of nesting), and edits an eight-item agenda as eight stacked, nearly identical sub-boxes of forty total form rows. That is not a bug in any single component. It is the predictable output of the current authoring model, in which **the bordered `Panel` is the only sectioning device, the expanded `FormRow` stack is the only record editor, and there is no collection primitive at all** — every list of things is hand-unrolled into repeated boxes by the page author.

This ticket is a **brainstorm, not an implementation plan** (explicit user instruction). It audits the page section by section with measured numbers and screenshots, diagnoses the underlying grammar gap, and lays out alternatives on three levels:

1. **Local fixes** — flat title-rule sections instead of nested borders; summary-first collection editing (tables, master-detail, disclosure rows); sticky action bars; sub-page or anchor navigation.
2. **New design-system primitives** the local fixes would need (Section, Disclosure, EditableList, TocRail, FieldGrid, StickyActionBar).
3. **The core concept: a composable "UI design-system grammar"** — a grammar-of-graphics-style layer where page authors declare *data + schema + presentation intent* (`show | edit | pick | arrange`) and the grammar compiles to Widget IR, instead of hand-assembling boxes. Today's recipes (`masterDetailTable`, `mediaLibrary`, `articleList`) are early, hard-coded special cases of exactly this idea.

Everything here applies across **all five DSL modules** (`ui.dsl`, `data.dsl`, `context_window.dsl`, `course.dsl`, `cms.dsl`), not just the CMS page — the user was explicit about this. The CMS page is simply the best specimen because it composes the most primitives on one screen.

## 1. Context for the newcomer: how pages are made today

Read the RAGEVAL-CMS-WIDGETS design doc first (`rag-evaluation-system/ttmp/2026/07/03/RAGEVAL-CMS-WIDGETS--…/design-doc/01-…`, especially §2–§6) — it explains the whole pipeline in depth. The 60-second version:

- **Widget IR** is a JSON tree (`WidgetNode = text | element | component`) defined in `packages/rag-evaluation-site/src/widgets/ir.ts`. A React `WidgetRenderer` (same package) renders it via a registry of ~70 widget adapters. Actions are declarative `ActionSpec`s (`navigate | server | download | event | copy`, each with optional `confirm`).
- **Pages are authored in JavaScript inside go-go-course** (`cmd/go-go-course/lib/pages/*.js`), running in a goja (Go) runtime. Authors call helper functions from five DSL modules implemented in `rag-evaluation-system/pkg/widgetdsl/module.go`:
  - `ui.dsl` — layout & primitives: `panel`, `stack`, `inline`, `formPanel`, `formRow`, `textInput`, `button`, `metadataGrid`, `action.*` …
  - `data.dsl` — `dataTable`, `cell.*` helpers, `recipes.masterDetailTable`, `recipes.actionToolbar` …
  - `context_window.dsl`, `course.dsl` — domain widgets (diagrams, transcripts, course shells).
  - `cms.dsl` — the new media/article/editor widgets and `recipes.mediaLibrary` / `recipes.articleList`.
- Every helper returns an IR node; the page builder composes them and the server returns the tree from `/api/widget/pages/:id`. State lives in the URL (`?asset=`, `?status=`), mutations in server actions (`POST /api/widget/actions/:name`) and native form posts.
- The design system underneath ("Classic Mac": 1px black borders, no radius, no shadow, 10–13 px tokenized type) is documented in `packages/rag-evaluation-site/GUIDELINES.md`.

The key structural fact for this ticket: **the DSLs expose components, not intent.** `ui.panel(...)` gives you a bordered box with a black title bar; `ui.formRow(...)` gives you one label+control row. There is nothing in between "a single widget" and "a whole recipe" — no notion of *section*, *collection*, or *record* that the system could lay out for you.

### Where to reproduce everything

```bash
cd go-go-course && devctl up --profile prod        # or: devctl up  (hot-reload dev)
# browser: http://127.0.0.1:8787/pages/settings → display name "admin_<you>"
#          http://127.0.0.1:8787/pages/admin-course-cms
```

Page source: `go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js` (+ `admin-common.js`, `admin-handout-editor.js`). Screenshots in this ticket under `sources/screenshots/`.

## 2. The specimen: section-by-section audit of `/pages/admin-course-cms`

Measured on 2026-07-04 against the live page (DOM inspection via Playwright). Total document height **5,611 px** (≈7 full screens at 1080p); **8 top-level panels**, **21 panels total**.

| # | Section (panel title) | Height | Inner panels | Form rows | Inputs | Notes |
|---|---|---|---|---|---|---|
| 0 | CMS building blocks | 164 px | 0 | 0 | 0 | intro caption + metadataGrid |
| 1 | Main page metadata | 703 px | 0 | 10 | 7 text + 3 textarea | one formPanel |
| 2 | Learning outcomes | 832 px | 0 | 7 | 7 textareas | one textarea per outcome, all expanded |
| 3 | **Agenda** | **2,158 px** | **8** | **40** | 32 text + 8 textarea | one nested panel per agenda item |
| 4 | Slides and handouts | 531 px | 2 | 0 | 0 | two upload drop-areas, each its own panel |
| 5 | Current course material | 393 px | 2 | 0 | 0 | 2 tables, 23 buttons |
| 6 | Media library | 493 px | 1 | 0 | 0 | cms.recipes.mediaLibrary + optional detail panel |
| 7 | Preview | 79 px | 0 | 0 | 0 | three navigation buttons |

Screenshots: `gg-course-admin-cms-full.png` (whole page), `section-agenda.png`, `section-outcomes.png`, `section-material-tables.png`, `section-media-library.png`.

### What the screenshots show, concretely

- **Boxes inside boxes inside boxes.** The Agenda screenshot is the money shot: a black-title-bar panel ("AGENDA") containing eight more black-title-bar panels ("AGENDA ITEM 1"…"AGENDA ITEM 8 · NEW"), each containing label/input rows and a mostly-empty 5-row textarea. Border weight is identical at every level, so nesting communicates nothing — the eye cannot tell section from item from field group. The title bar, the system's *strongest* visual element (inverted black), is spent on repetitive labels like "Agenda item 4".
- **Vertical waste on empty affordances.** Each of the seven outcome textareas and eight agenda descriptions reserves 3–6 rows whether or not there is content; the trailing "new row" placeholders (Outcome 6–7, Agenda item 8) cost ~600 px for pure potential input. ~40 % of the page is blank input area.
- **The list is invisible.** The agenda is a *schedule* — `14h30 → 15 min → title` per row — but nowhere can you see it as a schedule. You see it one 260-px box at a time. Reordering means clicking "Move up" in one box and re-scanning to verify. Compare: the same data as a 9-row table would be ~250 px total and instantly scannable.
- **No wayfinding.** The sidebar navigates *between* pages; nothing navigates *within* this one. Finding "Media library" means scrolling past 5,000 px of forms. The status line of a form you just submitted may be off-screen when the page reloads.
- **Inconsistent section devices.** Sections 0/7 are thin panels used as prose containers; 1–3 are formPanels; 4–6 are panels wrapping other panels. Same visual element, four different jobs.
- **Good bones nonetheless.** Tables (section 5) and the media tile grid (section 6) *are* the right density and scan well — evidence that the failure is not the design language but the lack of collection/section grammar in the authoring layer.

### Why the page author wrote it this way (and would again)

Look at `admin-course-cms.js`: the author had only these moves available —

```js
ui.panel({title}, ...children)          // the ONLY sectioning device
ui.formPanel({...}, ...formRows)        // the ONLY record editor (always fully expanded)
ui.formRow({label, control})           // one field
...agenda.map((item, i) => agendaEditorRow(item, i))   // lists = hand-unrolled boxes
```

`agendaEditorRow` is 40 lines of manual panel+rows per item because nothing else exists. The nesting problem, the repetition problem, and the length problem are all *downstream of the vocabulary*, which is exactly why the fix should be grammatical, not cosmetic.

## 3. Diagnosis: the missing grammar

Generalizing across the DSL pages (this repo's course/slides/handouts/sessions pages have the same shapes):

1. **One sectioning device.** `Panel` is simultaneously: page section, sub-item card, tool container, prose frame. There is no lighter-weight `Section` (title + rule + content, no box), no heading hierarchy, no anchors. The package actually has `SectionBlock` (layout layer) — the DSL never exposes a sectioning *policy*, so authors reach for `panel` every time.
2. **No collection primitive.** Every "N of the same thing" is `array.map(handRolledBox)`. The system cannot summarize, paginate, collapse, or reorder a collection because it never *knows* something is a collection — it only sees the unrolled output.
3. **No progressive disclosure.** Everything renders expanded. There is no accordion/disclosure, no "edit one at a time", no dialog-based editing (DialogShell exists in the package but has no DSL story for "summary row → open editor").
4. **No in-page navigation.** No TOC rail, no anchor links, no sub-page splitting convention. `CourseStudioShell`'s sidebar shows this is solvable — it is just not available *within* a page.
5. **Forms are monolithic.** `formPanel` couples: the box, the title, the method/action, the submit row, the status line. You cannot get "a form laid out as a two-column grid inside a flat section" without rebuilding it from elements.
6. **Density is uniform.** `density: compact|condensed|comfortable` exists on panels/tables but there is no page-level or section-level density policy, so mixed content (prose captions vs. 40-row forms) all renders at the same rhythm.
7. **Recipes are the escape hatch — and the proof.** `data.recipes.masterDetailTable`, `cms.recipes.mediaLibrary`, `cms.recipes.articleList` each encode *intent* ("browse rows and inspect one", "manage a media collection") and expand to a correct composition. They are point solutions; the grammar below is their generalization.

## 4. Brainstormed alternatives (no implementation — evaluation material)

Grouped from least to most ambitious. These are not mutually exclusive; the likely end state is E (grammar) built out of B+C (primitives + patterns), with A as the immediate styling posture.

### A. Flat sections: title + rule, boxes only for tools

Replace "panel = section" with a `Section` presentation: an uppercase 11-px heading, a 1-px rule, content, generous top margin; optional anchor id; heading levels 1–3 with decreasing weight. Panels remain for genuinely boxed things (an upload drop-zone, a selected-asset card, a dialog). The same page becomes:

```
COURSE CMS                          ← page header (exists today)
─ Main page metadata ────────────── ← Section level 1 (no box)
  [two-column field grid]
─ Learning outcomes ───────────────
  [collection editor]
─ Agenda ──────────────────────────
  [collection editor]
…
```

- Pros: kills the nesting ambiguity outright; roughly halves vertical chrome; matches how the handout/print pages already read.
- Cons: pure restyling — does nothing for the 40-row agenda; needs a rule for when a box is still correct (proposal: *interactive tools and selected-item cards get boxes; document structure never does*).

### B. Wayfinding: sub-pages, anchor rail, or tabs

Three sketches for "this page does too many jobs":

1. **Split into sub-pages** (strongest): `admin-course-cms` becomes a small hub or gains sidebar children — Metadata / Outcomes & Agenda / Files / Media. Each sub-page is one job, one form, one save. The sidebar (`CourseStudioShell` nav sections) already supports this; the handout editor added in Phase 6 (`admin-handout-editor?file=…`) is the existing precedent: list page → focused editor page, state in the URL.
2. **Anchor/TOC rail**: keep one page, add a slim right-hand rail of section links with scroll-position highlighting (new `TocRail` molecule). Cheapest wayfinding; does not reduce page weight.
3. **TabList** (exists in layout layer): tabs across the top of the admin area. Middle ground; hides sections rather than shortening them, and printed/long-scroll review is lost.

Open question for the intern to evaluate: does the admin area stay one URL (rail/tabs) or become several (sub-pages)? Sub-pages compose best with the URL-state philosophy and with form-post redirects (`?status=saved` returns you to a short page, not a 5,600-px one).

### C. Collection editing patterns for long lists

The agenda/outcomes problem generalized: *n records, each with a small schema; user needs scan, edit, add, remove, reorder.* Candidate patterns, each with the summary→detail split made explicit:

1. **Editable table** — one row per record, inline text inputs for short fields, description elided; whole table posts as one form (`agenda_3_title` names already exist). Best for ≤ ~20 records with mostly short fields. Reorder via row drag or per-row up/down icon buttons (boxed IconButtons from Phase 4 fit).
2. **Master-detail editing** — summary table (read-only, dense, scannable) + one detail form for the selected record (`?agenda=agenda-break`), exactly the `masterDetailTable` recipe plus write. One expanded editor at a time; the page stays ~900 px regardless of n. Precedent: media library + asset detail panel from Phase 6.
3. **Disclosure rows / accordion** — each record renders as a one-line summary (`14h30 · 15 min · Cadrage + concepts`) that expands in place to the full form. New `Disclosure` primitive (native `<details>` styled to guidelines is plausible). Keeps single-page editing; risks re-creating the wall of boxes when everything is expanded.
4. **Summary + dialog editing** — table rows with an Edit action opening `DialogShell` (exists; has widget adapter) containing the record form; form-posts from a dialog need a design decision (native post inside dialog vs. server action).
5. **Windowed lists** — "show first 5 · N more" for read-mostly lists (session lists, transcript blocks in the other DSLs). Pagination molecule already exists; what is missing is a *policy* attachable to any collection.

Recommendation to evaluate first: **2 (master-detail)** for agenda, **1 (editable table)** for outcomes — both need zero new interaction machinery (URL state + form posts), only the section/table-input primitives.

### D. Form ergonomics

- `FieldGrid`: two-column label/control grid for short fields (When/Where/Format are three 44-px rows today; one grid row fits all three). `metadataGrid` is the read-only sibling — the pair should rhyme.
- Auto-sizing / collapsed-until-focused textareas; character-count only near limit.
- Sticky action bar (`Save · Reset · status`) pinned to viewport bottom while its form is dirty — status is currently off-screen after scroll.
- Dirty-state marker on section headings (`Agenda •`), pairing with per-section saves.

### E. The core concept: a composable UI design-system grammar

> **Update (2026-07-04):** design-doc 02 develops this section into a concrete API sketch and answers the module question — there is no `grammar.dsl`; the verbs land in `data.dsl` (data grammar) and `ui.dsl` (structure grammar), with the domain modules reduced to schemas + marks.

**Grammar of graphics analogy.** GoG decomposes a chart into orthogonal layers — data → transforms → scales → marks → facets → theme — so a bar chart is not a `BarChart` widget but a *sentence*: `data + x:category + y:sum + mark:bar`. The proposal is the same decomposition for UI pages, across all five DSLs:

| GoG layer | UI grammar layer | Examples |
|---|---|---|
| data | **data** | records array, single record, file collection, snapshot |
| scales/aesthetics | **schema** | fields with type + role: `{id: key, number: time, title: primary, description: prose}` |
| statistical transform | **shaping** | sort, filter, page, group, summarize (first line, count) |
| mark | **arrangement** | table · list · cards · tiles · form · detail · timeline · strip |
| facet | **composition** | section, page/sub-page, split (master-detail), dialog, wizard, rail |
| — (interaction) | **verbs & bindings** | show · edit · pick · arrange(reorder) · confirm; bound to ActionSpecs / form posts |
| theme | **tokens/density** | existing theme.css tokens + density policy per section |

**Authoring sketch (pseudocode — shape, not API):**

```js
const g = require("grammar.dsl");                    // hypothetical layer above the 5 DSLs

g.page("admin-course-cms", { nav: "sections" },      // composition: sub-page hub or rail
  g.section("Main page metadata",
    g.record(metadata, {                             // data: one record
      schema: { kicker: g.f.short(), title: g.f.short({required: true}),
                tagline: g.f.prose(), blurb: g.f.prose({rows: "auto"}) },
      verb: "edit",                                  // → form
      arrange: "field-grid",                         // → 2-col grid, sticky actions
      submit: g.formPost("/settings/course-metadata"),
    })),
  g.section("Agenda",
    g.collection(agenda, {                           // data: records
      schema: { number: g.f.short({label: "Time"}), duration: g.f.short(),
                title: g.f.primary(), description: g.f.prose() },
      verb: "edit",
      arrange: "master-detail",                      // summary table + one editor
      select: g.urlParam("agenda"),                  // state in URL, per house rules
      reorder: g.action.server("admin-reorder-course-agenda"),
      remove: { action: g.action.server("…"), confirm: "Delete ${row.title}?" },
    })),
  g.section("Media library",
    g.collection(mediaAssets, { verb: "pick+manage", arrange: "tiles",
      open: g.action.navigate("/course-assets/${assetId}"), … })),  // ≈ today's cms.recipes.mediaLibrary
)
```

The grammar node compiles to plain Widget IR (component nodes the renderer already knows), so **the React side needs nothing new to start** — arrangements map onto DataTable/TileGrid/FormPanel/DialogShell/SplitPane compositions the way recipes already do. New primitives from §C/§D slot in as better compilation targets when they exist.

**Why a grammar and not more recipes?** Recipes are frozen sentences; the grammar is productive. `mediaLibrary` = `collection + tiles + pick/manage`; `masterDetailTable` = `collection + table + detail(show)`; the agenda editor the page *couldn't* express = `collection + table + detail(edit) + reorder`. With ~7 orthogonal axes the authoring surface covers those and every neighbor, instead of one Go function per combination in `pkg/widgetdsl/module.go`.

**Cross-DSL reach (user directive: "these improvements would apply across all the DSLs").** The grammar layer is domain-neutral; domain DSLs contribute *schemas and marks*, not layouts: `context_window.dsl` offers snapshot data + strip/stack/treemap arrangements; `course.dsl` offers lesson/handout records; `data.dsl`'s cells become field-role renderers; `cms.dsl` offers asset/article schemas. One `g.collection(sessions, {arrange: "table", …})` should look and behave identically whether the records are sessions, slides, assets, or annotations — that is the design-system payoff.

**Implementation-shape options (to brainstorm later, NOT now):**
- *Go-side grammar* — a `grammar.dsl` goja module compiling to IR in Go (like recipes today). Pros: no renderer changes, versioned with widgetdsl. Cons: layout logic in Go strings.
- *IR-level grammar* — new high-level IR node types (`CollectionEditor`, `Section`) with a TS-side interpreter. Pros: React owns layout where it belongs. Cons: IR schema growth, adapter complexity.
- *Hybrid (likely)* — grammar compiles in Go to a small set of new structural widgets (Section, EditableList, TocRail…) + existing components; the structural widgets are the only TS additions.

### F. New design-system primitives implied (wishlist, with layer placement)

| Primitive | Layer | Job | Guidelines notes |
|---|---|---|---|
| `Section` | layout | title + 1px rule + content, levels 1–3, anchor id | no box; uppercase label token |
| `TocRail` | molecule | in-page anchor nav with active highlight | 1px left rule, invert active |
| `Disclosure` / `Accordion` | layout | summary row ⇄ expanded content | `<details>`-based; ▸/▾ glyphs |
| `EditableList` / `Repeater` | organism | schema-driven collection editor (table/master-detail/disclosure modes) | composes DataTable + FormPanel |
| `FieldGrid` | layout | n-column label/control grid | sibling of metadataGrid |
| `StickyActionBar` | molecule | pinned submit/status for dirty forms | bottom border-top 1px, surface bg |
| `Toolbar` | molecule | standard section-level action strip | exists ad hoc in panels today |

All must obey GUIDELINES.md (tokens only, no radius/shadow, `data-rag-*`, stories per state).

## 5. What "good" looks like (acceptance sketch for a future implementation ticket)

- Admin CMS content reachable in ≤ 2 screens per view (sub-pages) or with a persistent rail (single page); agenda scannable as a table in < 300 px.
- No panel-inside-panel anywhere on the page; boxes only on interactive tools.
- The same collection grammar demonstrably renders agenda items, media assets, and uploaded sessions with only schema/verb changes.
- Storybook: every new primitive with all states; a "grammar gallery" story rendering one collection through each arrangement.

## 6. References

- Page under audit: `go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js`, `admin-common.js`, `admin-handout-editor.js`; server: `server.js`.
- DSL implementation: `rag-evaluation-system/pkg/widgetdsl/module.go` (helper maps, `actionObject`, recipes at `mediaLibraryRecipe`/`articleListRecipe`/`masterDetailTableRecipe`), TS declarations in `pkg/widgetdsl/typescript.go`.
- IR + renderer: `packages/rag-evaluation-site/src/widgets/{ir.ts,WidgetRenderer.tsx,actions.ts,defaultRegistry.ts}`.
- Design rules: `packages/rag-evaluation-site/GUIDELINES.md`; existing structural pieces: `SectionBlock`, `TabList`, `DialogShell`, `SplitPane`, `Pagination`, `TileGrid`.
- JS page-author API reference: `pkg/xgoja/providers/widgetsite/doc/02-widget-dsl-js-api-reference.md`.
- Background ticket: RAGEVAL-CMS-WIDGETS (design doc §6.7–6.9 for the go-go-course consumer architecture; diary Steps 6–10).
- Grammar of graphics: Wilkinson, *The Grammar of Graphics*; Wickham, *A Layered Grammar of Graphics* (ggplot2) — the layering/orthogonality argument being borrowed.
- Evidence: `sources/screenshots/{gg-course-admin-cms-full,section-agenda,section-outcomes,section-material-tables,section-media-library}.png`.

## 7. Open questions

1. One admin page with a rail, or sub-pages per concern? (Interacts with form-post redirects and printability.)
2. Where does the grammar compile — Go, TS, or hybrid (§4E)? Who owns layout decisions long-term?
3. How do grammar verbs bind to mutations — keep native form posts for records + server actions for item operations, or unify on server actions?
4. Reorder UX without client state: up/down buttons post-and-refresh today; is drag-and-drop worth the client-state exception?
5. Does `Disclosure` (`<details>`) violate the "state in URL" rule, or is transient open/closed state acceptable like SearchField's input value?
6. Which existing pages migrate first as the proving ground (admin CMS vs. sessions list vs. slides admin)?
