---
Title: DSL API sketch — grammar verbs in data.dsl and ui.dsl, module reorganization, no new module
Ticket: RAGEVAL-UI-GRAMMAR
Status: active
Topics:
    - cms
    - design-system
    - frontend
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Answers 'what is grammar.dsl?' — it should not exist. The grammar verbs live in the existing modules: data.dsl grows record/collection/schema (the data grammar), ui.dsl grows section/subpages/disclosure (the structure grammar), and the domain modules shrink to schemas + marks. Includes the cross-page audit that motivates this, a concrete API sketch with three worked page rewrites, a module reorganization table, and a compilation/migration story."
WhatFor: "Turn the doc-01 brainstorm into a reviewable API surface proposal before any implementation."
WhenToUse: "When evaluating or estimating the DSL grammar work; read doc 01 first for the diagnosis."
LastUpdated: 2026-07-04T15:45:00-04:00
---

# DSL API sketch — grammar verbs in data.dsl and ui.dsl, module reorganization, no new module

**Status: sketch for review, not a commitment.** Follows design-doc 01. Working assumptions taken from the ticket owner's answers: sub-page navigation, master-detail for long collections, hybrid compilation (Go grammar → a few new structural widgets + existing components), button-based reorder.

## 0. The headline answer: there is no `grammar.dsl`

Doc 01 used `grammar.dsl` as a placeholder. Sketching it against real pages shows a separate module would be a mistake:

- A sixth module would *add* incoherence — authors would import `ui.dsl` for a panel, `data.dsl` for a table, and `grammar.dsl` for "the good way", with two ways to do everything and no pressure to converge.
- The grammar's two halves already have natural homes. Declaring **what data looks like and how records are shown/edited** is `data.dsl`'s stated job ("data-display helpers and data recipes"). Declaring **how a page is structured and navigated** is `ui.dsl`'s job. The grammar is just those two jobs done at the level of intent instead of components.
- The recipes prove the modules can carry intent already (`data.recipes.masterDetailTable`, `cms.recipes.mediaLibrary`); the grammar generalizes recipes *in place*.

So the proposal is: **grow `data.dsl` into the data grammar, grow `ui.dsl` into the structure grammar, and shrink the domain modules to schemas + marks.** Module names, `require()` ids, and the provider wiring stay exactly as they are.

## 1. Evidence: the audit generalizes (task 2)

Same method as doc 01 §2 (DOM metrics on the live binary, 2026-07-04):

| Page | Height | Panels (top/total/depth) | Form rows | Verdict |
|---|---|---|---|---|
| `/pages/admin-course-cms` | 5,611 px | 8 / 21 / 2 | 57 | pathological (doc 01) |
| `/pages/admin-course-material` | 1,567 px | 5 / 8 / 1 | 0 | same shape, milder: 3 upload panels + boxed tables |
| `/pages/sessions` | 800 px | 2 / 2 / 0 | 0 | **reads well** — `masterDetailTable` recipe |
| `/pages/handouts?doc=…` | 5,043 px | 0 / 0 / – | 0 | **reads well** — long flat document, `HandoutDocumentShell` |
| `/pages/course` | 1,048 px | 0 | 0 | **reads well** — `CourseStudioShell` |
| `/pages/upload`, `/pages/settings` | ≤ 800 px | 1–2 / 0 | ≤ 5 | fine — single-job pages |

The correlation is exact: every page that reads well is either **one recipe/shell** (intent expressed once, layout owned by the system) or **a flat document**. Every page that degrades is hand-assembled panels around collections. Length itself is not the problem (the 5,043 px handout is fine); *boxed repetition without summarization* is. This is why the fix belongs in the DSL layer, and why "intent in, layout owned by the system" is the design principle.

## 2. Current module anatomy (what we are reorganizing)

From `pkg/widgetdsl/module.go` (`moduleSpec`, helper maps, `install`):

| Module | Exports today | Incoherences |
|---|---|---|
| `ui.dsl` | `page` + 30 component factories (layout, form, primitives) + `action` + recipes `metrics`, `actionToolbar` | has `sectionBlock` nobody uses; `formPanel` monolith; no section/nav/collection concepts |
| `data.dsl` | `dataTable` + `cell.*` + `action` + recipe `masterDetailTable` | the *only* module with a cell/field-role concept — and it's trapped inside tables |
| `context_window.dsl` | 20 domain widgets + style helpers + `action` + 2 recipes | `contextUploadDropArea` (generic file drop) stranded here — admin pages import the whole module for it |
| `course.dsl` | 11 domain widgets + `action` + 3 recipes | `markdownArticle`, `richArticle` are generic content marks, not course-specific |
| `cms.dsl` | 14 widgets + `action` + 2 recipes | `tag`, `meterBar`, `pagination`, `searchField`, `emptyState`, `tileGrid`, `breadcrumbs` are generic primitives that landed here because the ticket that made them was a CMS ticket |
| all five | `text`, `element`, `component`, `fragment`, `action` duplicated | harmless escape hatches, but `action` in 5 places blurs ownership |

Two structural facts worth keeping: every helper is a pure factory `props → IR node` (cheap, predictable), and recipes are plain Go functions that expand options → IR (the compilation pattern the grammar reuses).

## 3. Target module layout

```
ui.dsl      = STRUCTURE GRAMMAR + primitives
              page, subpages, section, toolbar, disclosure, fieldGrid, dialog…
              (everything that answers: how is this page organized?)
data.dsl    = DATA GRAMMAR + table marks
              schema/f.*, record, collection, cell, dataTable
              (everything that answers: how are these records shown/edited?)
context_window.dsl, course.dsl, cms.dsl
            = SCHEMAS + MARKS only
              domain widgets usable as arrangements/marks, field-role presets,
              domain schemas; their recipes become sugar over the grammar
```

### Moves and deprecations (aliases kept one release; helper maps make this a 10-line diff each)

| Symbol | From → To | Note |
|---|---|---|
| `tag`, `meterBar`, `pagination`, `searchField`, `emptyState`, `tileGrid`, `breadcrumbs` | `cms.dsl` → `ui.dsl` | generic primitives; cms keeps re-exports marked deprecated in the TS declarations |
| `contextUploadDropArea` | `context_window.dsl` → `ui.dsl` as `uploadDropArea` | generic; old name stays as alias |
| `markdownArticle`, `richArticle` | `course.dsl` → `ui.dsl` (or a future `content.dsl` — see open Q1) | generic content marks used by CMS editor preview |
| `action` | all five → canonical in `ui.dsl` | others keep it (cheap), docs say "import from ui.dsl" |
| `data.recipes.masterDetailTable` | stays | reimplemented as `data.collection(...)` sugar |
| `cms.recipes.mediaLibrary` / `articleList` | stay | reimplemented as grammar sentences with cms marks (§5.3) |

## 4. The data grammar: `data.dsl` additions

### 4.1 Schema and field roles

A schema names each field's **role** — role drives rendering (summary column vs. prose editor), input control, and elision. Roles, not types: `time` and `duration` are both strings; what matters is *how they behave in a summary row vs. an editor*.

```js
const data = require("data.dsl");
const f = data.f;                      // field-role helpers

const agendaSchema = data.schema({
  id:          f.key({ hint: "Stable internal anchor. Leave blank for a generated ID." }),
  number:      f.short({ label: "Time", width: "6ch", placeholder: "14h30" }),
  duration:    f.short({ width: "8ch", placeholder: "15 min" }),
  title:       f.primary({ required: true, maxLength: 160 }),   // the scannable column
  description: f.prose({ rows: "auto", maxLength: 800 }),       // elided in summaries
});
```

Role set (initial): `key`, `primary`, `short`, `prose`, `count`, `size`, `date`, `status`, `tags`, `media`, `href`. Each maps to: a summary renderer (today's `cell.*` become these — `cell.caption` ≈ `f.short` summary form), an editor control, and an elision rule. `data.cell` stays as the low-level escape hatch.

### 4.2 `data.collection(rows, opts)` — the workhorse

```js
data.collection(agenda, {
  schema: agendaSchema,
  title: "Agenda",
  verb: "edit",                        // "show" | "edit" | "pick" | "manage"
  arrange: "master-detail",            // "table" | "list" | "tiles" | "cards" | "master-detail" | "disclosure"
  select: data.urlParam("agenda"),     // selection state lives in the URL (house rule)
  submit: data.formPost("/settings/course-agenda"),   // record edits post natively
  reorder: ui.action.server("admin-reorder-course-agenda"),   // up/down buttons (assumption Q5a)
  remove:  { action: ui.action.server("admin-delete-agenda-item"),
             confirm: "Delete ${row.title}?" },
  create: true,                        // “New item” affordance
  page: { size: 20, param: "p" },      // windowing policy, optional
  empty: "No agenda items yet.",
});
```

Semantics: `verb` picks the interaction contract (`show` = read-only, `edit` = record editing per arrangement, `pick` = selection returns via action, `manage` = pick + create/remove/upload). `arrange` picks the summary presentation; for `master-detail` the detail editor is derived from the schema (prose fields expand, shorts become a field grid). Everything compiles to IR the renderer already has — `DataTable` + `FormPanel` + `Panel`/`Section` compositions — exactly like `masterDetailTableRecipe` does today, just parameterized by schema + verb instead of hand-fed columns.

`data.record(record, opts)` is the n=1 case (one form): `verb: "edit"`, `arrange: "field-grid" | "rows"`, same `submit`. The metadata form in doc 01 §4E's sketch is this.

### 4.3 Marks: how domain modules plug in

An `arrange` value can also be a **mark** — a component contract for rendering one record or the whole collection. Domain modules export theirs:

```js
data.collection(mediaAssets, {
  schema: cms.schemas.asset,           // exported by cms.dsl (matches CmsAsset)
  verb: "manage",
  arrange: cms.marks.assetTiles,       // tile grid of AssetTile — today's MediaLibraryPanel body
  select: data.urlParam("asset"),
  open:   ui.action.navigate("/course-assets/${assetId}"),
  upload: { action: ui.action.server("admin-upload-course-material",
            { payload: { kind: "media" } }), accept: ".svg,.png,.jpg,.jpeg,.webp,.gif" },
  remove: { action: ui.action.server("admin-delete-course-material"),
            confirm: "Delete ${row.filename}?" },
});
```

A mark declares which schema roles it consumes (`assetTiles` needs `media` + `primary` + `size` + `status`). `context_window.dsl` exposes `marks.strip/stack/treemap` for snapshot data the same way. This is the grammar-of-graphics split landing in module structure: **data.dsl owns the sentence, domain modules own vocabulary.**

## 5. The structure grammar: `ui.dsl` additions

### 5.1 `ui.section(title, opts, ...children)` — flat sectioning

Title + 1 px rule + content, levels 1–3, anchor id, optional toolbar; **no box**. Compiles initially to the existing `SectionBlock`; later to a dedicated `Section` widget. Deprecation posture: `panel` remains for tools/cards; GUIDELINES gains the rule *"document structure uses section; interactive tools use panel."*

### 5.2 `ui.subpages(...)` — the wayfinding answer (assumption Q2a)

```js
ui.subpages("admin-course-cms", {
  nav: "sidebar",                      // renders as child items of the shell nav
  pages: [
    { id: "metadata", title: "Metadata",  build: (q, ctx) => [metadataRecord(q)] },
    { id: "lists",    title: "Outcomes & Agenda", build: … },
    { id: "files",    title: "Slides & Handouts", build: … },
    { id: "media",    title: "Media",    build: … },
  ],
})
```

Compiles to N registered page ids (`admin-course-cms/metadata`, …) plus nav wiring — a convention over the existing `courseShellPage`/`parseWidgetPageId` machinery, so go-go-course's dispatcher keeps working. Single-page alternatives (`ui.tocRail`) stay in the §C toolbox but sub-pages are the default posture.

### 5.3 Worked rewrites (the acceptance test for the API)

**Agenda page (today: 8 nested panels, 40 form rows, 2,158 px):**

```js
ui.section("Agenda", { level: 1 },
  data.collection(agenda, { schema: agendaSchema, verb: "edit",
    arrange: "master-detail", select: data.urlParam("agenda"),
    submit: data.formPost("/settings/course-agenda"),
    reorder: ui.action.server("admin-reorder-course-agenda"), create: true }))
// → ~250 px summary table + one ~400 px editor for the selected row
```

**Sessions page (today: masterDetailTable recipe — must stay expressible):**

```js
data.collection(sessions, { schema: sessionSchema, verb: "show",
  arrange: "master-detail", select: data.urlParam("session"),
  open: ui.action.navigate("/pages/session-transcript--${row.id}") })
```

**Media library (today: cms.recipes.mediaLibrary):** §4.3 above. The recipe survives as a one-line wrapper calling that sentence — proof the grammar is a superset, and the migration path for all ten existing recipes.

## 6. Compilation and delivery (hybrid, assumption Q4a)

1. **Phase α — no renderer changes:** `section`, `record`, `collection` compile in Go to existing IR (SectionBlock, DataTable, FormPanel, TileGrid, Panel). Ships value immediately; layout quality capped by existing widgets.
2. **Phase β — structural widgets:** add `Section`, `FieldGrid`, `Disclosure`, `StickyActionBar`, `TocRail` to the package (doc 01 §4F wishlist) and retarget compilation. IR grows only these leaf-ish widgets — no interpreter on the TS side, adapters stay dumb.
3. TS declarations (`pkg/widgetdsl/typescript.go`) grow `schema/f/record/collection/section/subpages` signatures; the JS API reference (`doc/02-widget-dsl-js-api-reference.md`) gets a "grammar" chapter with the three worked examples.
4. Tests follow the existing pattern (`module_test.go`): expansion snapshots per verb × arrangement, plus boundary tests that domain marks only consume declared roles.

## 7. Open questions (beyond doc 01's)

1. Do `markdownArticle`/`richArticle` go to `ui.dsl` or justify a small `content.dsl` (prose/document marks)? Leaning ui.dsl — five modules is enough.
2. Schema sharing with the server: go-go-course builds rows in JS; should schemas also validate action payloads server-side (one schema, two uses), or stay presentation-only? (Presentation-only is the smaller first bite.)
3. `verb: "edit"` + `arrange: "table"` (inline editable table) needs a `TableInput`-ish control story — defer to phase β?
4. How do grammar nodes nest — can a `collection` detail contain another `collection` (agenda item → sub-steps)? Propose: yes structurally, but lint depth > 2.
5. Naming bikeshed to settle before building: `f.*` vs `field.*`, `arrange` vs `as`, `verb` vs `mode`.

## 8. References

- Module anatomy: `pkg/widgetdsl/module.go` (`moduleSpec` L22, helper maps L34–123, `install` L217, recipes L467+); declarations `pkg/widgetdsl/typescript.go`.
- Audit method + raw numbers: doc 01 §2 and this doc §1; screenshots in `sources/screenshots/`.
- Pages used as rewrite targets: `go-go-course/cmd/go-go-course/lib/pages/{admin-course-cms,sessions,admin-common}.js`.
