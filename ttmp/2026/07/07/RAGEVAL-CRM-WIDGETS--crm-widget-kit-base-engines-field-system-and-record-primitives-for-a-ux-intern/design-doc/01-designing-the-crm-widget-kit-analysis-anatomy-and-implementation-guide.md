---
Title: 'Designing the CRM Widget Kit: Analysis, Anatomy, and Implementation Guide'
Ticket: RAGEVAL-CRM-WIDGETS
Status: active
Topics:
    - design
    - design-system
    - widget-ir
    - react
    - frontend-architecture
    - intern-guide
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/MatrixGrid/MatrixGrid.tsx
      Note: Reference engine + MatrixCellPayload contract the CRM BoardCardPayload/FieldRenderPayload imitate
    - Path: repo://packages/rag-evaluation-site/src/scheduling/types.ts
      Note: Reference domain module the src/crm/types.ts imitates
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/engines.ts
      Note: Where the new FieldSpec + CRM engine *WidgetProps are added
    - Path: repo://packages/rag-evaluation-site/src/widgets/presets/scheduling.ts
      Note: Reference presets the crm.ts presets imitate
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-07T14:38:29.895850562-04:00
WhatFor: ""
WhenToUse: ""
---


# Designing the CRM Widget Kit: Analysis, Anatomy, and Implementation Guide

> **Who this is for.** You are a UX intern who will design a set of CRM widgets —
> the building blocks of a customer-relationship-management product: contact and
> company records, a deal pipeline, activity timelines, tasks, dashboards. You do
> not need to have written a line of this codebase yet. This document teaches you
> the one pattern the widget library is built on, gives you the CRM-specific pieces
> that pattern needs, and then walks you through designing each screen as an
> arrangement of those pieces, with drawings you can copy.
>
> **What "designing a widget" means here.** It is not only pixels. In this system a
> widget has an *anatomy* (what parts it is made of), a *prop shape* (the data it
> takes), and a *composition* (how it nests other widgets). Your design output for
> each widget is those three things — a labeled sketch, a prop table, and a
> component tree — and this guide shows you how to produce them. The engineers then
> implement against your design.
>
> **The most important thing to absorb.** We recently built a set of scheduling
> widgets (a Doodle-style poll, a calendar) by *not* building them as one-off
> screens. We built a few small generic **engines** and then assembled the product
> screens on top of them. Your CRM kit should be designed the same way. Sections 1
> and 2 teach that idea; the rest of the document applies it to CRM.

---

## Part 0 — How to read this document

Read Parts 1 and 2 in order — they are short and everything else depends on them.
Part 1 explains how a widget becomes something on screen. Part 2 explains the
engine/contract/preset pattern that your designs should follow. After that:

- **Part 3** is the CRM data model — the nouns your widgets display.
- **Part 4** is the field system, which is the heart of any CRM and the one piece
  most worth getting right.
- **Part 5** is the catalog of engines you will assemble from — some already exist,
  some are new and are yours to design.
- **Part 6** designs the actual screens (record page, pipeline board, timeline,
  dashboard) as compositions, with ASCII mockups and component trees.
- **Parts 7–9** are for the engineers who implement your designs: how a widget
  reaches the JSON layer and the server. Skim them so you know what is downstream
  of your work, but they are not where you will spend your time.

Two companion documents are worth having open. `RAGEVAL-SCHEDULE-WIDGETS`'s intern
guide is the worked example of this exact process for scheduling; steal its
structure freely. And `reference/02-the-widget-system-a-new-intern-s-guide-...`
teaches the rendering machinery in full if Part 1 here moves too fast.

---

## Part 1 — How a widget reaches the screen (the short version)

A widget in this system exists in three forms, and knowing this shapes how you
design. First it is a **React component** — the `.tsx` file an engineer writes.
Second it can be described as a **Widget IR node** — a plain JSON object of the form
`{ kind: "component", type: "ContactRecord", props: {...}, children: [...] }` — so
that a description of a UI can be sent over a network or produced by a script
without shipping React. Third, that JSON can be *produced* by a small server-side
scripting language (a DSL) so an analyst can assemble a screen without writing
React at all.

For your purposes the consequence is simple and concrete: **a widget's props must
be plain data.** No functions, because functions cannot be turned into JSON. When a
widget needs "what happens when I click this" or "how do I render this cell," the
system does not pass a function; it passes a small tagged data object describing the
intent, and one interpreter on the other side carries it out. A click becomes an
`ActionSpec` like `{ kind: "server", name: "deal.move" }`. A cell's rendering
becomes a `CellSpec` like `{ kind: "status", field: "stage" }`. You will design
using these data descriptions rather than callbacks, and Part 4 gives you the CRM
one — the `FieldSpec`.

The rendering path itself is a recursive tree walk you do not need to implement but
should be able to picture:

```
  Widget IR node                         React on screen
  { type: "PipelineBoard", props }  ──►  registry.get("PipelineBoard")
                                         → adapter turns props/specs into a real
                                           <PipelineBoard> React component
                                         → which renders its children the same way
```

The registry lookup is flat: any widget type can contain any other, so your
compositions are unconstrained. A dashboard can contain a board can contain a card
can contain a field. That freedom is what makes "assemble screens from small
pieces" practical.

---

## Part 2 — The pattern your designs must follow: engine + contract + preset

Here is the idea in one paragraph, then the drawing. A widget that shows domain data
tends to fuse three separable jobs into one component: arranging things in space
(layout, selection, drag, keyboard), knowing what a thing *means* (a deal, a
contact, a stage), and wiring the two together. When you keep these fused, every new
screen re-implements the arranging. When you separate them, you write the hard
arranging code once — as an **engine** — and every screen becomes a thin
configuration of it — a **preset** — connected through a fixed data shape — a
**contract**.

```
        ┌─────────────────────────────────────────────────────────────┐
        │  PRESET (domain)    pipelineBoard(pipeline)                   │
        │    speaks:  "deal", "stage", "owner", "amount"                │
        │    does:    configures an engine with CRM data + specs        │
        └───────────────▲─────────────────────────────────────────────┘
                        │ configures with plain-data specs
        ┌───────────────┴─────────────────────────────────────────────┐
        │  ENGINE (generic)   BoardEngine  (columns of draggable cards) │
        │    knows:   columns, drag-between-columns, selection ONLY     │
        │    knows nothing about deals or stages                        │
        └───────────────▲─────────────────────────────────────────────┘
                        │ CONTRACT: { card, columnId, selected, onMove } │
        ┌───────────────┴─────────────────────────────────────────────┐
        │  CARD / CELL / FIELD (swappable)   DealCard, ContactCard      │
        │    knows:   how one unit looks + means                        │
        └─────────────────────────────────────────────────────────────┘
```

The three terms, defined once:

- An **engine** is a generic component that owns arrangement and interaction and is
  blind to meaning. `BoardEngine` lays out columns of cards and handles dragging a
  card from one column to another; it has never heard the word "deal."
- A **contract** is the fixed data shape the engine hands each unit so it can place a
  unit it knows nothing about. For a board, that is roughly
  `{ card, columnId, selected, onMove }`. Any card component that accepts this shape
  is a valid card.
- A **preset** is a small function that configures an engine with domain vocabulary.
  `pipelineBoard(pipeline)` turns a pipeline into a configured `BoardEngine` whose
  columns are stages and whose cards are deals. The words "deal" and "stage" live
  here and nowhere in the engine.

Why this matters for a CRM specifically: a CRM is a handful of arrangement patterns
(a board, a record page, a timeline, a table, a list) reused across many object
types (contacts, companies, deals, tickets, custom objects). If you design the
arrangement patterns as engines, then supporting a new object type — even one a
customer invents — is writing a preset, not a screen. That is the difference between
a CRM you can extend and one you rewrite. Keep the drawing above in mind for every
widget you design: ask "what is the engine here, and what is the domain preset?"

**Key points:**

- Design every widget as an engine (arrangement, meaning-blind), a contract (the
  data shape for one unit), and a preset (domain configuration).
- The engine is where the expensive, reusable interaction lives; write it once.
- New object types become new presets, not new screens. This is what makes the CRM
  extensible.

---

## Part 3 — The CRM domain model (the nouns)

Before designing screens, name the data. In this codebase, domain data lives in its
own module of plain TypeScript types with no UI in it — the scheduling widgets have
`src/scheduling/types.ts`; the CRM kit gets `src/crm/types.ts`. Keeping the data
types separate from the widgets is what lets many widgets share one definition of
"what a Deal is." Here is the model, abbreviated to the fields that drive UI:

```ts
// src/crm/types.ts  (design sketch — plain data, no React, no IR)

type Id = string;

interface FieldValue { /* string | number | boolean | Id | Id[] | null | {…} */ }

interface Contact {
  id: Id;
  name: string;
  avatarUrl?: string;
  title?: string;
  companyId?: Id;
  fields: Record<string, FieldValue>;   // email, phone, custom fields — see Part 4
  ownerId?: Id;
  tags?: string[];
  updatedAtISO: string;
}

interface Company {
  id: Id; name: string; domain?: string; logoUrl?: string;
  fields: Record<string, FieldValue>;
  ownerId?: Id; tags?: string[];
}

interface Deal {
  id: Id; title: string;
  amount?: number; currency?: string;
  stageId: Id;                          // which pipeline column it sits in
  pipelineId: Id;
  contactIds?: Id[]; companyId?: Id;
  ownerId?: Id; closeDateISO?: string;
  fields: Record<string, FieldValue>;
  status: "open" | "won" | "lost";
}

interface Stage { id: Id; name: string; order: number; colorKey: string; probability?: number; }
interface Pipeline { id: Id; name: string; stages: Stage[]; }

type ActivityKind = "note" | "email" | "call" | "meeting" | "task" | "stage_change" | "field_change";
interface Activity {
  id: Id; kind: ActivityKind;
  actor: { id: Id; name: string; avatarUrl?: string };
  atISO: string;
  subjectId: Id;                        // the record this activity is on
  title: string; body?: string;
  meta?: Record<string, unknown>;       // e.g. { from: "Lead", to: "Qualified" } for stage_change
}

interface Task {
  id: Id; title: string; dueISO?: string;
  status: "open" | "done";
  assigneeId?: Id; relatedId?: Id; priority?: "low" | "med" | "high";
}

// The field *definitions* — the schema a workspace configures. Part 4 is all about this.
interface FieldDef {
  key: string;
  label: string;
  type: FieldType;                      // "text" | "email" | "currency" | "select" | "relation" | …
  options?: { value: string; label: string; colorKey?: string }[];  // for select
  relatedObject?: "contact" | "company" | "deal" | "user";           // for relation
  required?: boolean; readOnly?: boolean;
}
```

Two modeling choices are worth flagging because they shape every widget downstream.
First, records carry a `fields: Record<string, FieldValue>` bag *in addition to* a
few first-class columns (`name`, `amount`). This is deliberate: a CRM's whole selling
point is that customers add their own fields, so the field set is data, not code.
Second, everything that happens to a record is an `Activity` with a `kind`. A stage
change, an email, a note — one stream, many kinds. That uniformity is what lets a
single timeline engine render the entire history of a record. Both choices mirror how
the scheduling model kept computed data (`SlotTally`) separate from raw data
(`MeetingPoll`).

Alongside `types.ts`, the module carries `src/crm/palettes.ts` (the `ContextStyleSet`
color maps — one for stage colors, one for activity kinds, one for tags) and
`src/crm/fixtures.ts` (sample contacts, a sample pipeline with deals, a sample
timeline) so every widget has realistic data to design against.

---

## Part 4 — The field system (the heart of the CRM)

Everything in a CRM is ultimately a **typed field** shown in one of two modes:
reading or editing. An email is a `mailto:` link when you read it and a validated
text box when you edit it. A stage is a colored pill when you read it and a dropdown
when you edit it. A currency amount is a right-aligned formatted number when you read
it and a number input when you edit it. Get this one abstraction right and contacts,
companies, deals, and any custom object all render for free; get it wrong and you
will special-case fields forever. This is the single most important thing you will
design.

The pattern is the same defunctionalization the IR already uses for cells and
actions (Part 1): a field's rendering is described by data, a `FieldSpec`, and one
interpreter turns it into the right control for its mode. You are, in effect,
designing CRM's version of the existing `CellSpec`.

```ts
type FieldType =
  | "text" | "longtext" | "email" | "phone" | "url"
  | "number" | "currency" | "percent" | "date" | "datetime"
  | "boolean" | "select" | "multiselect" | "tags"
  | "relation" | "user" | "address";

interface FieldSpec {
  key: string;                 // which key in record.fields
  type: FieldType;
  label?: RenderableValue;
  options?: SelectOption[];    // for select/multiselect
  relatedObject?: string;      // for relation/user
  readOnly?: boolean;
  styleSet?: ContextStyleSet;  // colors for select/tag values (reuse the palette contract)
}
```

The **contract** — the fixed shape the engine hands each field — is what makes fields
swappable and is the exact analogue of `MatrixCellPayload`:

```ts
interface FieldRenderPayload {
  spec: FieldSpec;
  value: FieldValue;
  mode: "read" | "edit";
  invalid?: boolean;
  onChange: (next: FieldValue) => void;    // edit mode reports changes back
  onCommit: () => void;                     // blur / enter
}
```

A single `FieldRenderer` component consumes this payload and switches on
`spec.type × mode` to produce the correct control. In read mode a `currency` renders
formatted text; in edit mode it renders a number input. This table is the core of
your field design work — fill in the read and edit appearance for each type:

| `FieldType` | Read appearance | Edit appearance |
|---|---|---|
| `text` / `longtext` | plain / wrapped text | text input / textarea |
| `email` | `mailto:` link + copy | validated input |
| `phone` | `tel:` link | input with format mask |
| `url` | truncated link (opens new tab) | input |
| `currency` | right-aligned `$1,200` | number input + currency prefix |
| `percent` | `62%` (+ optional bar) | number input 0–100 |
| `date` / `datetime` | `Jul 9, 2026` (relative on hover) | date picker (reuse `MonthGrid`!) |
| `boolean` | check / dash | toggle |
| `select` | colored pill (via `styleSet`) | dropdown of colored options |
| `multiselect` / `tags` | row of pills | pill multi-select |
| `relation` | avatar/name chip → record | search-and-pick (async) |
| `user` | avatar + name | user picker |
| `address` | formatted lines / map link | grouped inputs |

Notice two reuse opportunities that fall out immediately. A `date` field's edit mode
*is* the `MonthGrid` engine from the scheduling kit — you do not design a new date
picker, you present the existing one. And a `select`/`tag` field's colors go through
the same `ContextStyleSet` palette contract the context diagrams and scheduling
widgets use — one coloring mechanism across the whole product.

The `FieldRenderer` is an **engine** in the Part 2 sense: it owns "how is a typed
value shown and edited" and knows nothing about which object the field belongs to.
Above it sits a small **`RecordFieldList`** engine that lays out many fields (label +
control rows, grouped into sections) and hands each one a `FieldRenderPayload`. Above
*that* sit the record presets (`contactRecord`, `dealRecord`) that supply the field
definitions and values. Three layers, same as always.

**Key points:**

- Every CRM value is a typed field with a read mode and an edit mode; design the two
  appearances per type in the table above.
- A field's rendering is a `FieldSpec` (data), interpreted by one `FieldRenderer`
  engine through a `FieldRenderPayload` contract — the CRM analogue of
  `CellSpec`/`MatrixCellPayload`.
- Reuse what exists: `date` edit mode is `MonthGrid`; `select`/`tag` colors are the
  `ContextStyleSet` palette.
- Custom fields are data (`FieldDef`), so a workspace adding a field adds no code.

---

## Part 5 — The engine catalog

Some engines already exist and you should reuse them; some are new to CRM and are
yours to design. The table separates them; the subsections design the new ones.

| Need in the CRM | Engine | Status | Contract |
|---|---|---|---|
| Pipeline kanban | **`BoardEngine`** | NEW | `{ card, columnId, selected, onMove }` |
| Record page (header + fields + related + timeline) | **`RecordShell`** | NEW (specializes `MasterDetailShell`) | slots |
| Typed field view/edit | **`FieldRenderer`** + **`RecordFieldList`** | NEW | `FieldRenderPayload` |
| Activity history | **`ActivityFeed`** | NEW | `{ activity, isLast, onOpen }` |
| KPI numbers | **`StatTile`** / `MetricRow` | NEW (small) | — |
| Deal / contact table | `MatrixGrid` (or `DataTable`) | exists | `MatrixCellPayload` |
| Related list, task list, saved views | `ItemList` / `CollectionPanel` | exists\* | list-item contract |
| Segment / filter builder | **`FilterBar`** | NEW | query spec |
| Funnel / stage breakdown | `SegmentedBar` (stack layout) | exists | segment contract |
| Dashboard layout | `DashboardGrid` | exists | — |
| Date pick (in fields & filters) | `MonthGrid` | exists | day contract |

\* `ItemList`/`CollectionPanel` are proposed in the decomposition ticket
(`RAGEVAL-WIDGET-DECOMPOSITION`); coordinate so the CRM kit consumes them rather than
re-inventing a list.

### 5.1 `BoardEngine` — the signature new engine (kanban pipeline)

This is the CRM equivalent of what `MatrixGrid` was for scheduling: the flagship
generic engine that proves the pattern. A board is columns of cards where a card can
be dragged from one column to another; the engine owns the columns, the drag, the
drop targets, per-column scrolling, and selection, and knows nothing about deals.

```
  Pipeline: Sales                                                 [+ Deal]  [⋯]
  ┌── Lead ──────┐ ┌── Qualified ──┐ ┌── Proposal ───┐ ┌── Won ───────┐
  │ $12k         │ │ $40k          │ │ $88k          │ │ $30k         │  ← column header: name + Σ amount
  │ 3 deals      │ │ 2 deals       │ │ 4 deals       │ │ 1 deal       │
  ├──────────────┤ ├───────────────┤ ├───────────────┤ ├──────────────┤
  │ ┌──────────┐ │ │ ┌───────────┐ │ │ ┌───────────┐ │ │ ┌──────────┐ │
  │ │ Acme      │ │ │ │ Globex    │ │ │ │ Initech   │ │ │ │ Umbrella │ │  ← DealCard (the swappable unit)
  │ │ $8,000    │ │ │ │ $25,000   │ │ │ │ $40,000   │ │ │ │ $30,000  │ │
  │ │ 🧑 Dana ●  │ │ │ │ 🧑 Lee  ●● │ │ │ │ 🧑 Priya  │ │ │ │ ✓ won    │ │
  │ └──────────┘ │ │ └───────────┘ │ │ └───────────┘ │ │ └──────────┘ │
  │ ┌──────────┐ │ │ ┌───────────┐ │ │  …            │ │              │
  │ │ Wayne Co │ │ │ │ Stark Ind │ │ │               │ │              │
  │ └──────────┘ │ │ └───────────┘ │ │               │ │              │
  └──────────────┘ └───────────────┘ └───────────────┘ └──────────────┘
     drag a card between columns → onMove({cardId, from, to, beforeId})
```

Base props and contract (design sketch):

```ts
interface BoardEngineProps<Card> {
  columns: { id: string; header: ReactNode; footer?: ReactNode; accent?: string }[];
  cards: Card[];
  columnOf: (card: Card) => string;          // which column a card is in
  getCardId: (card: Card) => string;
  renderCard: (p: BoardCardPayload<Card>) => ReactNode;   // the swappable unit
  selectedCardId?: string;
  onMove?: (m: { cardId: string; from: string; to: string; beforeId?: string }) => void;
  onCardSelect?: (cardId: string) => void;
}
interface BoardCardPayload<Card> {
  card: Card; columnId: string; selected: boolean;
  dragging: boolean; onSelect: () => void;
}
```

The `pipelineBoard(pipeline, deals)` **preset** supplies `columns` from the stages
(header = stage name + summed amount), `columnOf = deal => deal.stageId`, a
`DealCard` for `renderCard`, and an `onMove` that emits a `deal.move` server action.
A different preset, `contactsByStatusBoard`, could reuse the identical engine for a
lead-status board. The engine is written once.

> **Design note for you:** the `DealCard` is where your visual design energy goes —
> it is small, high-frequency, and read at a glance. Design its anatomy explicitly
> (title, amount, owner avatar, due indicator, tag dots) and its states (default,
> selected, dragging, won/lost). The board engine is mostly invisible; the card is
> the product.

### 5.2 `RecordShell` — the record page

Every object (contact, company, deal) opens to a record page with the same anatomy:
a header (identity + key actions), a left column of fields, and a right column of
related lists and a timeline. This is a specialization of the `MasterDetailShell`
pattern; design it once and every object type is a preset.

```
┌───────────────────────────────────────────────────────────────────────────┐
│ 🧑  Dana Whitmore                                   [ Edit ] [ Log ▾ ] [⋯] │  ← header: avatar, name, title, actions
│     VP Sales · Acme Corp · 🏷 enterprise                                    │
├───────────────────────────┬───────────────────────────────────────────────┤
│  DETAILS                   │  ACTIVITY                          [+ Note]    │
│  Email    dana@acme.com 📋 │  ● Today   Email sent "Q3 proposal"           │  ← ActivityFeed (Part 5.3)
│  Phone    +1 555 0142      │  ○ Jul 6   Stage → Qualified                   │
│  Owner    🧑 You            │  ○ Jul 2   Call · 12m                          │
│  Stage    ▧ Qualified      │  ○ Jun 28  Note "Wants annual billing"        │
│  ─────────────────         │  ────────────────────────────                 │
│  CUSTOM                     │  DEALS (2)                        [+ Deal]    │  ← related list = CollectionPanel
│  Segment  Mid-Market       │  • Acme renewal      $8,000   Qualified       │
│  NPS      9                 │  • Acme expansion    $25,000  Proposal        │
│                            │  COMPANIES (1)                                 │
│                            │  • Acme Corp                                    │
└───────────────────────────┴───────────────────────────────────────────────┘
```

Component tree (this is the "component structure" half of your design deliverable,
expressed as YAML):

```yaml
RecordShell:                       # NEW engine (MasterDetailShell specialization)
  header:
    RecordHeader:                  # avatar + name + subtitle + action bar
      identity: { avatar, name, subtitle }
      actions: [Button(Edit), Menu(Log), IconButton(more)]
  left:
    Stack:
      - SectionBlock(label="Details"):
          RecordFieldList:         # NEW engine — lays out FieldRenderer rows
            fields: [email, phone, owner, stage]   # each a FieldSpec (Part 4)
            mode: read
      - SectionBlock(label="Custom"):
          RecordFieldList: { fields: [segment, nps], mode: read }
  right:
    Stack:
      - Panel(title="Activity", actions=[Button(+Note)]):
          ActivityFeed: { activities }             # Part 5.3
      - CollectionPanel(title="Deals"):            # related list, existing engine
          body: ItemList / MatrixGrid
      - CollectionPanel(title="Companies"): ...
```

The presets `contactRecord(contact)`, `companyRecord(company)`, `dealRecord(deal)`
differ only in which `FieldSpec`s they list and which related lists they include. The
shell, the field list, and the timeline are shared.

### 5.3 `ActivityFeed` — the timeline

A record's history is one stream of `Activity` items of different `kind`s, rendered
in reverse-chronological order with a connective spine. The engine owns the spine,
the grouping by day, and the "load more"; each item kind is a swappable renderer
chosen by `activity.kind`.

```
  ● Today
  │  ✉  Email sent · "Q3 proposal"                     10:24  🧑 You
  │     ▸ 2 attachments
  ○ Jul 6
  │  ▧  Stage changed  Lead → Qualified                 09:12  🧑 You
  │  ☎  Call · 12 min · "left voicemail"                08:40  🧑 Lee
  ○ Jun 28
  │  📝 Note · "Wants annual billing, decision by Q3"    15:03  🧑 You
  [ load earlier ]
```

The contract is `{ activity, isLast, onOpen }`; the preset maps `Activity.kind` to a
small renderer (`EmailActivity`, `CallActivity`, `NoteActivity`, `StageChangeActivity`,
`TaskActivity`). Design each kind's one-line anatomy — icon, primary text, metadata,
timestamp, actor — and the shared spine. This is the CRM analogue of the transcript
message list: one feed engine, many message kinds.

### 5.4 Smaller pieces

- **`StatTile` / `MetricRow`** — a labeled number with an optional delta and
  sparkline, for dashboards ("Open pipeline $2.1M ▲12%"). Reuse `ProportionTrack`/
  `SegmentedBar` for the inline bar. Design the tile anatomy (label, value, delta,
  trend) and a `MetricRow` that lays out several.
- **`FilterBar` / segment builder** — a row of filter chips ("Owner = me", "Stage =
  Proposal", "Amount > $10k") that compiles to a query spec the server runs.
  Design the chip, the add-filter menu, and the saved-segment affordance. The engine
  owns chip layout and the query spec; presets supply the field list to filter on.
- **`PipelineFunnel`** — stage counts as a funnel; this is `SegmentedBar` in stack
  layout with stage colors, no new engine.

---

## Part 6 — Designing the screens (compositions)

With the engines named, each product screen is an arrangement. For each screen your
deliverable is the ASCII anatomy plus the YAML component tree, as in Part 5.2. Here
are the four core screens and how they decompose; design them in this order.

**1. Pipeline board** — `pipelineBoard` preset over `BoardEngine`, wrapped in a
`Panel` with a pipeline switcher and a `[+ Deal]` action. The one screen that most
sells a CRM; start here.

**2. Contact / Company / Deal record** — `RecordShell` with the three-layer
composition in Part 5.2. Design `contactRecord` first; company and deal are the same
shell with different field lists.

**3. Deal table (list view)** — `MatrixGrid`/`DataTable` with a `FilterBar` above and
a `CollectionPanel` frame (search, saved views, pagination). Columns are `FieldSpec`s
rendered in read mode, so the table and the record page share field rendering.

**4. CRM dashboard** — `DashboardGrid` of `StatTile`s and a `PipelineFunnel` and a
"my tasks" `ItemList` and a "recent activity" `ActivityFeed`. Pure composition of
existing engines; no new primitives.

A fifth, the **tasks inbox**, is an `ItemList` of `Task`s grouped by due date with a
`FieldRenderer`-driven inline "mark done" — worth designing because it exercises
inline field editing outside a record page.

For every screen, apply the Part 2 question: *what is the engine, and what is the
preset?* If a screen seems to need a new engine, check the Part 5 catalog first —
most CRM screens are the same five engines in different arrangements.

---

## Part 7 — How your designs become IR and DSL (for the engineers)

This part is downstream of your work; skim it so you know the shape of what you are
feeding. Each engine you design gets, on the engineering side, a JSON prop interface
and a builder. The field system adds one new spec to the IR — `FieldSpec` — beside
the existing `CellSpec` and `ActionSpec`, defined in
`src/widgets/ir/engines.ts` (the file that already holds the scheduling engine
specs). Each new engine (`BoardEngine`, `RecordShell`, `ActivityFeed`,
`FieldRenderer`, `RecordFieldList`, `StatTile`, `FilterBar`) gets a `*WidgetProps`
interface there, an entry in the `RagWidgetType` union, an IR adapter
(`X.widget.tsx`) that maps the JSON to the React component, and a `.widget.yaml`
manifest.

On the Go/Goja DSL side (`pkg/widgetdsl/`), the CRM widgets become a `crm.dsl` module
whose helper map exposes `board`, `recordShell`, `activityFeed`, `field`,
`recordFieldList`, `statTile`, and `filterBar`, plus composite **recipes** for the
presets (`contactRecord`, `pipelineBoard`) that emit whole configured subtrees — the
Go analogue of the TypeScript presets. The important prior art to hand the engineers:
the DSL *already* has a `record`/`collection` grammar with field roles
(`grammar.go` — roles like `key/primary/status/date/tags/measure`) that overlaps
heavily with the CRM field types; the CRM field system should be reconciled with it,
not built parallel to it. That reconciliation is the main engineering design question
this kit raises.

---

## Part 8 — The backend (for the engineers)

Interactions emit `ActionSpec`s; a server action dispatcher receives them at
`POST /api/widget/actions/<name>` and returns `{ ok, refresh?, patch?, toast? }`.
The CRM actions your designs imply:

| Action | Fired by | Payload |
|---|---|---|
| `deal.move` | dragging a card | `{ dealId, fromStage, toStage, beforeId }` |
| `field.update` | editing any field | `{ object, recordId, key, value }` |
| `record.create` | `[+ Deal]` / `[+ Contact]` | `{ object, fields }` |
| `activity.log` | `[+ Note]`, log call/email | `{ subjectId, kind, title, body }` |
| `task.complete` | inbox checkbox | `{ taskId }` |
| `segment.query` | `FilterBar` | `{ object, filters }` → `{ data }` |

The server owns the parts your widgets must not: computing stage sums and pipeline
metrics, validating field types, enforcing permissions per field, running saved-segment
queries, and appending to the activity stream on every mutation (so a `field.update`
also produces a `field_change` activity that the timeline then shows). Records are
loaded by REST (`GET /api/crm/contacts/:id` returning the record plus its field
definitions); mutations go through the action envelope so the widget can refresh or
patch in place.

---

## Part 9 — Build order, and API & file reference

**Design in this order** (each unblocks the next):

1. The **field system** — the `FieldSpec` type and the read/edit table in Part 4.
   Everything else renders fields, so design this first and get it reviewed.
2. The **`DealCard`** and **`BoardEngine`** — the pipeline board is the signature
   screen; the card is your highest-visibility artifact.
3. The **`RecordShell`** + **`ActivityFeed`** — the record page and its timeline.
4. The **dashboard** and **tasks inbox** — pure compositions that validate the kit.

**Files the engineers will create** (mirroring the scheduling kit):

```
src/crm/types.ts                 domain DTOs (Part 3)
src/crm/palettes.ts              stage/activity/tag ContextStyleSets
src/crm/fixtures.ts              sample records, pipeline, timeline (for your Storybook)
src/components/atoms/{DealCard,StatTile,...}/         card + tile atoms
src/components/molecules/{FieldRenderer,RecordFieldList,ActivityFeed,BoardEngine,FilterBar}/   engines
src/components/organisms/{RecordShell,PipelineBoardPanel,CrmDashboard,TasksInbox}/             screens
src/widgets/ir/engines.ts        + FieldSpec + the new *WidgetProps
src/widgets/presets/crm.ts       contactRecord/companyRecord/dealRecord/pipelineBoard presets
pkg/widgetdsl/module.go          + crm.dsl module, helper map, recipes
```

**Reference — the pattern to copy.** The scheduling kit is the worked example of this
entire process end to end:

- `src/components/molecules/MatrixGrid/MatrixGrid.tsx` — the reference engine + the
  `MatrixCellPayload` contract your `BoardCardPayload`/`FieldRenderPayload` imitate.
- `src/widgets/presets/scheduling.ts` — the reference presets your `crm.ts` presets
  imitate.
- `src/scheduling/{types,palettes,fixtures}.ts` — the reference domain module your
  `src/crm/` imitates.
- `src/widgets/ir/engines.ts` — where the new `FieldSpec` and engine props go.
- Ticket `RAGEVAL-SCHEDULE-WIDGETS` — the intern guide and DSL handoff you are
  parallel to.
- Ticket `RAGEVAL-WIDGET-DECOMPOSITION` — the `ItemList`/`CollectionPanel`/
  `MasterDetailShell` engines you will consume; coordinate so you reuse them.

## Open questions

1. **Field system vs. the DSL `record` grammar.** The Go DSL already has a field-role
   grammar (`grammar.go`). Do we reconcile the CRM `FieldSpec` with it (one field
   model) or keep them separate? Strong recommendation: reconcile.
2. **Custom objects.** Should `RecordShell` be generic over object type from day one
   (so a customer-defined "Property" object works), or start with the three built-ins
   (contact/company/deal) and generalize later?
3. **Inline edit vs. edit mode.** Do records edit field-by-field inline (click a value
   to edit it) or flip the whole page into an edit mode? This changes the
   `FieldRenderPayload` and the `RecordShell` header design.
4. **Board scale.** Pipelines can have thousands of deals; does `BoardEngine` need
   virtualized columns in v1, or is a per-column cap acceptable to start?

## References

- `RAGEVAL-SCHEDULE-WIDGETS` — the scheduling widget kit: intern guide, presets, and
  DSL handoff. The worked example of engine + contract + preset.
- `RAGEVAL-WIDGET-DECOMPOSITION` — the decomposition analysis and the shared engines
  (`ItemList`, `CollectionPanel`, `MasterDetailShell`) the CRM kit should consume.
- `reference/02-the-widget-system-a-new-intern-s-guide-to-the-architecture.md` — the
  rendering machinery in full.
- `packages/rag-evaluation-site/GUIDELINES.md` — the design-system rules (typography
  tokens, `data-rag-*`, layer ownership) your components must follow.
