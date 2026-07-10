---
Title: Design Diary
Ticket: RAGEVAL-CRM-WIDGETS
Status: active
Topics:
    - design
    - design-system
    - widget-ir
    - react
    - frontend-architecture
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-07T14:38:29.605658499-04:00
WhatFor: ""
WhenToUse: ""
---

# Design Diary

## Goal

Produce an intern-facing design/implementation guide that lets a new UX intern
design a full CRM widget kit (records, pipeline board, field system, timeline,
dashboard) using the established **engine + contract + preset** pattern — the same
approach used for the scheduling/Doodle widgets. Deliverable is the guide, stored in
the ticket and uploaded to reMarkable, to hand off in the morning.

## Step 1: Ticket + design guide

Created ticket `RAGEVAL-CRM-WIDGETS` and wrote
`design-doc/01-designing-the-crm-widget-kit-analysis-anatomy-and-implementation-guide.md`
in the textbook style (per the textbook-authoring skill): foundational-first, prose
that develops each idea, concrete pseudocode/ASCII/diagrams, no analogies, terms
defined on first use.

### Structure of the guide
- **Parts 0–2:** how a widget reaches the screen (the three forms + data-not-functions
  constraint), and the engine/contract/preset pattern the intern's designs must follow.
- **Part 3:** the CRM domain model (`Contact`/`Company`/`Deal`/`Pipeline`/`Stage`/
  `Activity`/`Task`/`FieldDef`) as a `src/crm/types.ts` sketch — the key modeling
  choices being the `fields` bag (custom fields as data) and the uniform `Activity`
  stream (one timeline for all event kinds).
- **Part 4 (the heart):** the field system — `FieldSpec` + a `FieldRenderer` engine
  through a `FieldRenderPayload` contract (the CRM analogue of `CellSpec`/
  `MatrixCellPayload`), with a read/edit appearance table per `FieldType`, and the
  reuse insight that `date` edit mode IS `MonthGrid` and `select` colors are the
  `ContextStyleSet` palette.
- **Part 5:** the engine catalog — new (`BoardEngine` [signature kanban],
  `RecordShell`, `ActivityFeed`, `FieldRenderer`/`RecordFieldList`, `StatTile`,
  `FilterBar`) vs reused (`MatrixGrid`, `ItemList`/`CollectionPanel`, `SegmentedBar`,
  `DashboardGrid`, `MonthGrid`). ASCII + prop sketches for the new ones.
- **Part 6:** the four core screens as compositions (pipeline board, record page,
  deal table, dashboard) with ASCII anatomy + YAML component trees.
- **Parts 7–9:** IR/DSL wiring (a `crm.dsl` module, `FieldSpec` in `ir/engines.ts`),
  backend actions (`deal.move`/`field.update`/`activity.log`/…), and build order +
  file/reference index pointing at the scheduling kit as the worked example.

### Key design decisions captured in the guide
- The field system is the first thing to design; everything renders fields.
- `BoardEngine` is the flagship new engine (kanban), the CRM parallel to `MatrixGrid`.
- Records carry a `fields` bag so custom fields need no code.
- Flagged the main engineering reconciliation: the CRM `FieldSpec` vs the DSL's
  existing `record`/field-role grammar (`grammar.go`) — recommend converging.

### Next
- Upload the guide to reMarkable.

## Step 2: Implementing the kit (engineering follow-through)

Moved from the design guide to actually building the widget kit in
`packages/rag-evaluation-site`, mirroring the committed scheduling kit
(`src/scheduling` + `MatrixGrid`/`MonthGrid`/`TimeGrid`/`SegmentedBar` engines).
Working on branch `task/improve-rag-evaluation-system`. Build order follows Part 9
of the design doc. Committing per milestone; typecheck (`tsc --noEmit`) is the
gate and lefthook re-runs typecheck+biome on every commit.

### M1 — CRM domain module (`src/crm/`) ✅  commit 6a568c2
Pure-data DTOs (`types.ts`): `Contact`/`Company`/`Deal`/`Pipeline`/`Stage`/
`Activity`/`Task`/`FieldDef`/`FieldValue`/`FieldType`, mirroring
`src/scheduling/types.ts`. `palettes.ts` = three `ContextStyleSet`s (stages,
activity kinds, tags) + `ACTIVITY_GLYPHS`. `fixtures.ts` = sample contacts,
Sales pipeline + deals, activities, tasks, users, contact/deal field defs, and
server-computed `StageSummary`s. No React/IR — the shared vocabulary layer.

### M2 — Field system (the heart) ✅  commit a073bba
The defunctionalized field renderer, the CRM analogue of `CellSpec`/`MatrixGrid`:
- `FieldRenderer` molecule — switches on `FieldType × mode`. Read: mailto/tel/url
  links, formatted currency/percent/date, colored select pills (via
  `ContextStyleSet`), relation/user chips (resolved through a `resolveRef`),
  boolean check. Edit: TextInput/TextareaInput/SelectInput/number+prefix/native
  date input/comma-entry for tags. `FieldRenderPayload` is the domain-blind
  contract handed to each field.
- `RecordFieldList` molecule — arranges label+control rows into sections, hands
  each field a payload. Owns arrangement only.
- IR: `FieldSpec` + `FieldRendererWidgetProps` + `RecordFieldListWidgetProps` in
  `ir/engines.ts`; `BoardEngine`/`ActivityFeed`/`StatTile`/`MetricRow` prop
  interfaces staged there too. New `crm.dsl` `WidgetModule`, `crmWidgetRegistry`
  merged into the default registry. Adapters (`*.widget.tsx`) + manifests
  (`*.widget.yaml`).

Design decision: `FieldType`/`FieldValue`/`FieldOption` live once in
`crm/types.ts` and are imported by both the React engine and the IR (the IR
already imports `ContextStyleSet` from `../../context`, so importing from a
sibling domain module is consistent). One source of truth, no drift.

### Review feedback mid-M3 (important course-correction)
Three notes from the user reshaped the decomposition — captured because they
override the design doc's forward-looking plan in favor of *what the code
already does*:
1. **Extensive per-widget stories are required** (GUIDELINES §Storybook:
   default/empty/overflow/selected/error/interactive). Added component-level
   `.stories.tsx` for every widget, not just the WidgetRenderer IR stories.
2. **No rounded corners — retro style.** The existing components (Tag, CycleCell,
   MatrixGrid) use zero `border-radius`. Stripped all radius from the new CSS
   (square cards, pills, avatars).
3. **Follow the existing decomposition.** Two corrections:
   - `DealCard` is a **molecule**, not an atom (it composes title/subtitle/meta
     content slots, like `TranscriptMessageCard`).
   - Do **not** mint a `crm.dsl` widget module/registry. The scheduling kit
     created no `schedule.dsl` registry — its generic engines live in
     `data.dsl`/`time.dsl` and *all* domain logic is in `presets/scheduling.ts`.
     So the generic engines (`FieldRenderer`, `RecordFieldList`, `BoardEngine`,
     `ActivityFeed`, `StatTile`) register under the generic `dataWidgetRegistry`
     (`module: data.dsl`); CRM lives only in `src/crm/` + `presets/crm.ts`.
     Also dropped the invented `MetricRow` — several `StatTile`s compose via the
     existing `TileGrid`/`DashboardGrid` layout instead.

### M3 — BoardEngine kanban + DealCard ✅  commit f4f03fd (+ 4ff3e3c fixes)
`BoardEngine` molecule: columns of cards, HTML5 drag-between-columns, drop
targets, per-column scroll, selection; `BoardCardPayload` contract; domain-blind
(the CRM parallel to `MatrixGrid`). `DealCard` molecule: the swappable card unit
(title/amount/meta, stage accent bar, won/lost badge, selected/dragging states).
IR adapter renders each card from a `BoardCardSpec` (title/subtitle/meta
CellSpecs) into a `DealCard`. `pipelineBoard(pipeline, deals)` preset →
stages-as-columns, deals-as-cards, `deal.move` server action.

### M4 — ActivityFeed timeline + RecordShell page ✅  commit 74094b1
`ActivityFeed` molecule: reverse-chronological stream grouped by day with a
connective spine; per-kind glyph+color from the shared `ContextStyleSet`;
load-more + clickable rows (CRM analogue of the transcript message list).
`RecordShell` organism: header + left field column + right timeline/related
column, composing the existing `Panel`/`SplitPane`/`Stack` (React-first, not
IR-registered). `contactRecord` preset emits the record page as an IR
composition of already-registered widgets — no bespoke node, mirroring how
`pollResults` composes `Stack`+`SegmentedBar`.

### M5 — StatTile + dashboard + tasks inbox ✅  commit 885e1fc
`StatTile` molecule: labeled number + delta/trend arrow + inline `MeterBar`.
`crmDashboard` (TileGrid of StatTiles + `pipelineFunnel` SegmentedBar + recent
activity) and `tasksInbox` (tasks with an inline `FieldRenderer` boolean
"mark done" → `task.complete`) presets. All four core screens now exist.

### Verification (Storybook + Playwright)
Booted Storybook (`-p 6007`) and drove it with Playwright. Screenshots in
`various/screenshots/`. All render correctly, retro square style throughout, no
render errors (only a favicon 404):
- **pipeline board** — stage columns with colored accent bars, cards, summed
  headers, horizontal scroll.
- **contact record page** — mailto/tel links, owner/company avatar chips, a
  green select pill, percent, tag pills, and the day-grouped activity spine with
  colored kind glyphs.
- **dashboard** — StatTiles with colored trend arrows + MeterBars, the
  stage-colored funnel with legend.
- **field-type table** — every `FieldType` in read *and* edit mode (the design's
  Part 4 table, made real: currency `$` prefix, percent `%` suffix, native
  date/datetime pickers, colored select/tag pills, relation/user chips).

Every commit passed lefthook (typecheck + biome lint/format).

### Known follow-ups (not blockers)
- Board card amount shows `12,000` (no `$`) and owner shows the raw id
  (`u-you`); a richer preset could resolve owner names + currency-format.
- The `contactRecord` IR preset repeats the "Details" label (Panel title +
  first section group); cosmetic.
- Reconcile `FieldSpec` with the Go DSL `record`/field-role grammar
  (`grammar.go`) — the main engineering design question flagged in the doc.

## Related

- `design-doc/01-designing-the-crm-widget-kit-analysis-anatomy-and-implementation-guide.md`
- Sibling ticket `RAGEVAL-SCHEDULE-WIDGETS` (the worked example the CRM kit mirrors).
- Sibling ticket `RAGEVAL-WIDGET-DECOMPOSITION` (the shared engines the CRM kit consumes).
