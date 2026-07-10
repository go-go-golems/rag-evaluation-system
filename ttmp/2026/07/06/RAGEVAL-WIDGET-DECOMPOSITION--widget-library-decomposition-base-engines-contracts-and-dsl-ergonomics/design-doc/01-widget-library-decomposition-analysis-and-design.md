---
Title: 'Widget Library Decomposition: Analysis and Design'
Ticket: RAGEVAL-WIDGET-DECOMPOSITION
Status: active
Topics:
    - design-system
    - widget-ir
    - ui-dsl
    - react
    - frontend-architecture
    - intern-guide
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/ContextStripDiagram/ContextStripDiagram.tsx
      Note: One of 5 context diagrams that re-implement a shared SegmentEngine (patternClass/styleName/formatTokens duplication)
    - Path: repo://packages/rag-evaluation-site/src/context/types.ts
      Note: Dormant ContextDiagramSegment/ContextDiagramView contract (scaffolded, never wired) — the segment engine seam
    - Path: repo://packages/rag-evaluation-site/src/widgets/cellRenderers.tsx
      Note: getPath/renderTemplate — one of 4 accessor mini-languages to unify into AccessorSpec
    - Path: repo://packages/rag-evaluation-site/src/widgets/registry.ts
      Note: RenderContext — target for ctx.actionHandler/renderFields helpers
    - Path: repo://pkg/widgetdsl/module.go
      Note: DSL runtime — cellObject/helper maps/install special-casing (Part 6 opportunities)
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-06T21:04:28.099834802-04:00
WhatFor: ""
WhenToUse: ""
---


# Widget Library Decomposition: Analysis and Design

> **Audience:** a new intern who will help refactor the RAG-Evaluation widget
> library. This document assumes you have read the companion guide,
> `reference/02-the-widget-system-a-new-intern-s-guide-to-the-architecture.md`,
> which teaches how the system works. This one is different: it looks at the code
> that *exists* and identifies, concretely and with file references, where the
> library repeats itself and how to fix that. Where the companion guide taught the
> architecture, this document is a survey of its condition and a plan for
> improving it.
>
> **It proposes changes; it does not make them.** Every finding cites real files
> and line numbers so you can open them and confirm. Read Parts 0–2 first — they
> establish the vocabulary and the single idea the rest depends on — then use
> Parts 3–8 as a catalog you can dip into.

---

## Part 0 — How to read this document, and a glossary

Before anything else: this document uses a handful of terms repeatedly, and if you
do not know them the rest will read as noise. Here they are, defined once.

- **The IR (intermediate representation).** A tree of plain JSON objects that
  *describes* a user interface without being one. "Intermediate" because it sits
  between the thing that wants to describe a UI and the thing that draws it. It is
  the pivot of the whole system.
- **A tier.** One of the three levels the system is built in. There are three, and
  they matter because the pattern we care about should appear in all three:
  1. **React components** — the actual `.tsx` files under
     `packages/rag-evaluation-site/src/components/`. These paint pixels.
  2. **The Widget IR** — the JSON representation and the code that renders it, under
     `.../src/widgets/`.
  3. **The Go/Goja DSL** — server-side JavaScript functions, implemented in Go under
     `pkg/widgetdsl/`, that *produce* the IR JSON. "DSL" means domain-specific
     language: a small vocabulary of functions for one job, here "build a UI."
- **An engine.** A generic component that owns *layout and interaction* — how units
  are arranged in space, how selection works, how keyboard navigation works — and
  knows nothing about what the units mean. `MatrixGrid` is an engine: it lays out
  rows and columns and knows nothing about votes or calendars.
- **A contract.** The fixed data shape an engine hands to each unit it lays out, so
  the engine can place a unit without knowing what the unit is. `MatrixCellPayload`
  is a contract: `{ row, col, value, selected, onAction }`.
- **A preset.** A small function that configures an engine with domain meaning.
  `availabilityMatrix(poll)` is a preset: it takes a poll and returns a configured
  `MatrixGrid`. The words "vote" and "slot" live here; the engine never hears them.
- **A seam.** The boundary between generic and specific — in practice, the contract.
  "Cutting the seam" means actually separating an engine from the domain code it is
  currently fused with, so the engine becomes reusable.
- **A code smell (or just "smell").** A surface symptom of a deeper design problem.
  The same forty lines copied into five files is a smell; the underlying problem is
  a missing engine.
- **Drift.** When the same fact is written by hand in several places, the copies
  fall out of sync over time. That divergence is drift.
- **A mirror.** One of those hand-maintained copies of a shared fact. Five mirrors
  of "the set of widget types" means the same list is retyped in five files.
- **To hoist.** To move a helper that is currently copied in many files up into one
  shared location that they all import. Purely mechanical; changes nothing visible.
- **A passthrough adapter.** An IR adapter (`X.widget.tsx`) that does almost nothing
  — it forwards its props straight to the React component with no translation. Many
  adapters are nearly this, which is why they are candidates for generation.
- **Codegen (code generation).** Writing a script that produces source code from a
  single description, so you edit the description instead of the generated files.
- **Behavior-preserving.** A change that does not alter what the user sees or can do.
  Every refactor proposed here is behavior-preserving; that is what makes them safe.
- **Opinionated.** A DSL is opinionated when it makes choices for the author rather
  than exposing every knob. This is a virtue here: fewer, better-chosen ways to do
  a thing. The goal is to stay opinionated *and* become more capable, not to expose
  everything.

With that vocabulary in hand, one sentence states the thesis of the entire
document. **One structural pattern — an engine, a contract, and a preset — should
recur at every tier; today it appears cleanly only in the newest code (the
scheduling widgets), and older code re-implements the same shapes by hand.** Parts
3 and 4 catalog where; Part 5 addresses a related problem (drift between mirrors);
Part 6 does the Go DSL; Part 7 sequences the work.

The findings below come from a systematic read-only review of five areas of the
codebase — atoms and foundation, molecules, organisms, the IR and renderer, and
`pkg/widgetdsl`. Where a number appears ("copied in 7 files"), it is a real count
you can reproduce with `grep`.

---

## Part 1 — The system as it is today

This part is a compressed recap of the architecture the companion guide teaches in
full. If Part 1 reads easily, you are ready for the analysis; if it does not, read
the companion guide first.

### 1.1 The layers and the file conventions

The React components are organized into five layers, from smallest to largest:
**foundation** (text and typography primitives like `Text` and `Caption`),
**atoms** (small self-contained controls like `Button`, `Tag`, `CycleCell`),
**layout** (structural pieces like `Panel`, `SplitPane`, `SidebarShell` that
arrange other things), **molecules** (reusable data-display components like
`DataTable` and `StepList`), and **organisms** (large feature panels like
`MediaLibraryPanel` and `MeetingPollPanel`). The rule is that a lower layer never
imports a higher one, so an atom can never depend on an organism. A component that
answers "how are regions arranged?" belongs in layout; one that answers "what
domain data is shown?" is a molecule or organism.

Every public component is a folder of six files, and you should recognize them all:
`X.tsx` is the React component; `X.module.css` is its private styling;
`X.stories.tsx` is its Storybook page (the visual test surface); `X.widget.tsx` is
its **IR adapter** (the code that turns a JSON node into this component);
`X.widget.yaml` is its **manifest** (a data description of the widget — Part 5
returns to this at length); and `index.ts` re-exports it. Components are strictly
*presentational*: they receive data as props and emit callbacks, and they never
reach into a global store, the router, or the backend. That restriction is what
makes them reusable across the real app, Storybook, and tests unchanged.

### 1.2 The Widget IR render pipeline

The path from a JSON node to a React component is a recursive tree walk. The
following diagram is the whole pipeline; the paragraph after it explains the one
step that matters most.

```
                         WidgetNode  (serializable JSON)
          ┌───────────────┬────────────────┬─────────────────────┐
          │ TextNode      │ ElementNode    │ ComponentNode        │
          │ {kind:"text"} │ {kind:element, │ {kind:"component",   │
          │               │  tag, attrs,   │  type, props,        │
          │               │  children}     │  children}           │
          └──────┬────────┴───────┬────────┴──────────┬──────────┘
                 │                 │                   │
 WidgetRenderer  ▼                 ▼                   ▼
  .tsx:13        node.text    createElement(tag)   registry.get(type)      ← flat lookup
                                                        │
                                                        ▼
                                     adapter.render(props, children, ctx, node)
                                                        │
                                                        ▼
                                          React tree (real components/**)

  RenderContext (registry.ts:13) — one per mount, threaded to every adapter:
    renderNode(node) · renderChildren(children) · renderValue(value)
    bindAction(action, ctx) · dispatchAction(action, ctx)
                                     │
                                     ▼  actions.ts:41 dispatchWidgetAction
      ActionSpec kind switch: copy · event(print/fullscreen/CustomEvent)
        · navigate(pushState) · download(<a download>) · server(POST /api/widget/actions/:name)
```

The step that governs everything is `registry.get(type)`. To render a component
node, the renderer takes its `type` string — `"Panel"`, `"MatrixGrid"` — and looks
it up in a flat table that maps type strings to **adapters**. The renderer does not
know what a Panel is; it knows how to look one up. Because that lookup is per-node
and flat, any widget can contain any other widget, and there is no coupling between
a parent's type and its children's types. The one capability an adapter receives is
`ctx`, a five-method `RenderContext`; it never sees the registry, the store, or the
router. That boundary — the renderer is a small, uniform interpreter and adapters
are sandboxed — is the strongest part of the current design, and none of the
proposals here touch it.

Key files: `widgets/WidgetRenderer.tsx`, `widgets/registry.ts`,
`widgets/defaultRegistry.ts`, `widgets/actions.ts`, `widgets/cellRenderers.tsx`;
the IR types in `widgets/ir/{core,actions,cells,engines,props}.ts`.

### 1.3 The pattern the IR already got right: functions as data

You cannot put a JavaScript function into JSON — functions do not survive
serialization. But a UI is full of functions: a click handler, a "how do I draw
this cell?" rule, a "what color is this value?" rule. The IR solves this with a
technique called **defunctionalization**: instead of the function, you store a
*data description* of what the function should do, as a tagged union (an object with
a `kind` field naming the variant), and you keep one interpreter that reads the
description and does the work. The system does this three times, once for each of
the three function-shaped things a UI needs:

| The function you want | Stored instead as | Defined in | Interpreted by |
|---|---|---|---|
| a click handler | `ActionSpec` (`kind`: navigate/download/server/event/copy) | `ir/actions.ts` | `actions.ts:41` `dispatchWidgetAction` |
| "how do I render this cell?" | `CellSpec` (`kind`: field/number/status/…/actionButton) | `ir/cells.ts` | `cellRenderers.tsx:9` `renderCell` |
| "what color is this value?" | `StyleBySpec` (value → styleKey → style) | `ir/engines.ts:14` | `styleBy.ts:11` `resolveStyleByVars` |

The payoff of doing this is that the set of things a UI can do or render becomes a
closed, inspectable vocabulary rather than arbitrary code. That property is the
foundation the whole IR rests on, and — importantly for Part 4 — it is a pattern
the codebase applies *inconsistently*: it has three clean tagged-union specs and
then, alongside them, several ad-hoc one-off versions of the same idea that were
never unified.

The most recent work (the scheduling ticket) added one more layer on top of these
primitives: a cluster of *engines* — `MatrixGrid`, `SegmentedBar`, `MonthGrid`,
`TimeGrid` — each with a payload contract in the `MatrixCellPayload` style. The
comment at `engines.ts:8-13` states outright that this cluster is the model the
rest of the library should converge on. This document is, in large part, the
argument for taking that comment seriously.

### 1.4 The Go/Goja DSL

The third tier lets a server-side script produce IR JSON. The Go program embeds a
JavaScript engine (Goja); a script calls `require("ui.dsl")` and gets an object of
helper functions; calling `ui.panel({title:"x"})` returns a plain object that is
exactly an IR node. The machinery is small: each DSL module (`ui.dsl`, `data.dsl`,
`context_window.dsl`, `course.dsl`, `cms.dsl`) is a `moduleSpec` (`module.go:23`)
carrying a `helpers` map from JS name to IR component type; `runtime.install`
(`module.go:236`) turns each map entry into a JS function via `componentFactory`
(`module.go:565`). Beyond the plain helpers there are builders for the specs from
1.3 (`cell.*` at `module.go:272`, `action.*` at `module.go:318`), composite
builders called *recipes* that emit a whole subtree at once (`recipesObject`
`module.go:495`, e.g. `masterDetailTableRecipe`), a higher-level "grammar" of
intent verbs (`grammar.go`), and a newer typed experiment (`v2/spec/`,
`v2_builders.go`). The TypeScript type declarations that script authors rely on are
generated from the same specs (`typescript.go`) and checked against the real
compiler in a test (`typescript_fixture_test.go`).

---

## Part 2 — The one idea: engine + contract + preset

Every proposal in this document is one refactoring lens applied over and over. The
lens is this: a widget that renders domain data usually contains *three separable
things fused into one file*, and separating them makes the reusable part reusable.
The three things are the engine, the contract, and the preset — defined in Part 0's
glossary and drawn here:

```
        ┌────────────────────────────────────────────────────────┐
        │  PRESET (domain)   availabilityMatrix(poll)             │
        │    knows: "a vote", "a slot", the availability palette  │
        │    does:  configures the engine with domain specs       │
        └───────────────▲────────────────────────────────────────┘
                        │ configures via serializable specs
        ┌───────────────┴────────────────────────────────────────┐
        │  ENGINE (generic)   MatrixGrid                          │
        │    knows: geometry / layout / selection ONLY            │
        │    talks to units through a stable CONTRACT             │
        └───────────────▲────────────────────────────────────────┘
                        │ { row, col, value, selected, onAction }
        ┌───────────────┴────────────────────────────────────────┐
        │  UNIT (swappable)   CycleCell, DayCell, EventBlock      │
        │    knows: how one unit looks + means                    │
        └────────────────────────────────────────────────────────┘
```

The reason to bother is economic. The engine holds the expensive, error-prone code
— layout math, selection state, keyboard handling — and you want to write and debug
that once. The unit is small and specific. The preset is a dozen lines. When a
fourth grid-shaped UI appears, you write a small preset and a small unit and inherit
a correct engine for free. The cost is a single layer of indirection, the contract,
and that indirection is the entire point: it is the line between "reusable" and
"domain-specific," drawn explicitly instead of left implicit.

When you hold this lens up to the library, three recurring problems come into
focus. They are the "three smells," and the rest of the document is organized
around them:

1. **Duplicated engine logic.** The same layout, list, or shell mechanics are
   re-implemented once per domain because the engine was never extracted. Examples:
   the context diagrams (five copies), the studio shells (two), the selectable rails
   (three), the collection panels (two), the list molecules (three).
2. **Ad-hoc spec shapes.** The same *idea* — reading a value by a path, expressing a
   selection, describing a list item, choosing a color — is expressed as a slightly
   different bespoke prop on each widget instead of one shared spec. There are at
   least four different value-accessors, fourteen different selection fields, and ten
   different list-item shapes.
3. **Drift between hand-maintained mirrors.** The set of widgets is described by hand
   in five different places (`RagWidgetType`, the `WidgetProps` union,
   `defaultRegistry`, the manifests, and the Go helper maps), and those copies have
   already diverged.

Parts 3 and 4 catalog smells 1 and 2, layer by layer. Part 5 addresses smell 3.
Part 6 applies the same thinking to the Go DSL.

---

## Part 3 — Duplicated engines, layer by layer

This part is the catalog for smell 1: reusable engines that exist implicitly,
spread across the widgets that need them, because the seam was never cut. Each entry
names the duplication with file references and describes the engine that would
replace it.

### 3.1 Atoms & foundation

Most atoms are already clean, self-contained primitives — `Button`, `TextInput`,
`CycleCell`, `DateTile`, `RatioBadge` need nothing. The work here is a few fused
engines and some bookkeeping.

The clearest fused engine is the **badge family**. Three components —
`ContentStatusBadge` (`atoms/ContentStatusBadge/ContentStatusBadge.tsx:10`),
`TranscriptRoleBadge` (`atoms/TranscriptRoleBadge/…:8`), and foundation `StatusText`
(`foundation/StatusText/StatusText.tsx:22`) — do exactly the same thing: they take
an enumerated value (a content status, a transcript role, a generic status), look
it up in a map to find a glyph and a CSS class, and render a small pill,
`<span data-x={value}>glyph label</span>`. The engine is "an enum-keyed pill";
each component wrote its own copy. Extract one base pill (a generalized `StatusText`
or a new `GlyphBadge`) that takes the `Record<Enum, glyph>` map as a prop, and the
three domain badges become thin wrappers that each pass their own map.

A second, smaller duplication: `RatioBadge` and `MeterBar` both draw a proportional
bar. `RatioBadge.tsx:52-56` and `MeterBar.tsx:21-32` each independently clamp a
number to `[0,1]`, set a fill element's width to that fraction, and add the
accessibility attributes that mark it as a progress bar ("progressbar a11y" — the
ARIA roles a screen reader needs). Extract a `ProportionTrack` primitive and have
both consume it.

`Caption` is a special case worth naming because it is already the pattern done
*wrong by omission*: `Caption` is essentially `Text` locked to the metadata size
with a restricted tone set, yet `Caption.tsx:14-21` keeps its own copy of the
tone-to-CSS-class map that also lives in `Text.tsx:36-45`. `Caption` should render
`Text` internally, or both should import one shared tone map, rather than maintain
two.

`ContextStudioNavIcon` fuses a generic thing with a specific one in a single file:
lines 19-37 are a reusable titled `<svg>` shell (a wrapper that handles the viewBox,
the accessible title, sizing) — call that an "SVG sprite," a reusable icon frame —
while `renderIcon` at lines 39-96 is a domain-specific set of icon paths (course,
slides, transcript, handout). Extract an `SvgIcon` base atom that takes a path or a
registry, and the domain icons become data.

Finally, two pieces of bookkeeping. Some pure, domain-blind atoms — `MediaThumb`,
`MeterBar`, `Tag` — are declared in their manifests as belonging to the `cms.dsl`
module even though nothing about them is CMS-specific; they should be reclassified
under the neutral `ui.dsl`. And `ContextStyleSwatch` has a `.widget.tsx` adapter but
no `.widget.yaml` manifest at all, which Part 5 shows is a symptom of the larger
drift problem. Separately, a purely mechanical cleanup: the string
`[styles.root, …].filter(Boolean).join(" ")` — the idiom for combining CSS class
names — is copy-pasted in about 27 components with no shared helper (conventionally
called `cx()`), and the `tone` and `size` option sets diverge from atom to atom
(`neutral/positive/warning/muted` in one, `accent/success/danger` in another;
`sm/md/lg` here, `compact/normal/large` there). A shared `cx()` and one canonical
`Tone` and `Size` union would remove that divergence.

### 3.2 Molecules — the flagship opportunity

The single largest instance of smell 1 in the library is the family of
context-window diagrams, and it is worth studying closely because it is the clearest
example of "an engine that was designed, half-built, and abandoned."

Five components draw the same thing. `ContextStripDiagram`,
`ContextGroupedStripDiagram`, `ContextStackDiagram`, `ContextTreemap`, and
`ContextBudgetBar` each take a set of "parts" (regions of a model's context window,
each with a token count and a style) and render them as a row or column of colored,
labeled, clickable segments. All five run the same loop, and the duplication is
literal, not merely similar:

- `patternClass(pattern)` — a helper that maps a fill-pattern name to a CSS class —
  is **byte-for-byte identical in seven files** (the five diagrams plus
  `TranscriptMessageCard.tsx:50` and `AnnotationNoteCard.tsx:21`).
- `styleName(styleSet, styleKey)` and `formatTokens(tokens)` are identical in all
  five diagrams.
- The block that makes a segment keyboard-accessible — the cluster of `role`,
  `tabIndex`, `aria-pressed`, `onClick`, and `onKeyDown` handlers wired to
  `handleContextPartKeyDown` — is duplicated in all five, differing only in whether
  the layout is horizontal or vertical.
- The one genuine difference between the diagrams is the *geometry*: a segment's size
  is expressed as a width percentage in the strip/grouped/budget variants, as a
  `minHeight` in the stack (`ContextStackDiagram.tsx:49`), and as `flexBasis`/
  `flexGrow` in the treemap (`ContextTreemap.tsx:67-71`).

Read that list and the diagnosis writes itself: there is one engine here — call it a
**segment engine** — and instead of extracting it, five files each copied it, with
the only real variation being a single geometry rule. The proposal is to build the
engine once and reduce each diagram to a preset. The engine would look like this
(and `SegmentedBar` at `SegmentedBar.tsx:36` is already about eighty percent of it):

```ts
// engine — owns layout, selection, keyboard nav, tooltip; knows nothing about "context"
interface SegmentEngineProps {
  segments: PositionedSegment[];
  styleSet: ContextStyleSet;
  layout: "strip" | "stack" | "treemap";   // the ONLY geometry difference between the 5 diagrams
  total?: number;
  selectedId?: string;
  orientation?: "horizontal" | "vertical";
  onSegmentSelect?: (id: string) => void;
  renderTooltip?: (seg: SegmentPayload) => ReactNode;
  markers?: SegmentMarker[];
}
// contract — what the engine hands each segment (the exact analogue of MatrixCellPayload)
interface SegmentPayload {
  id: string; label: ReactNode; styleKey: string;
  visualStyle: ContextVisualStyle;   // resolved once, by the engine
  value: number; fraction: number;   // fraction is layout-agnostic; the layout rule turns it into px/%
  selected: boolean; interactive: boolean;
  onSelect: () => void;
}
```

What makes this a *finish-the-refactor* rather than a new invention — that is, what
makes it completing work someone already started rather than adding a new
abstraction — is a striking pair of facts. First, `context/types.ts:80` already
defines a type called `ContextDiagramSegment` whose fields are *exactly* this
contract, and a grep shows it has **zero references anywhere in the codebase**:
someone declared the seam and never cut it. Second, `context/types.ts:78` already
defines `ContextDiagramView`, the enumeration `strip | stack | budget | treemap`,
and it *is* used elsewhere (by course slides and article blocks) — so the rest of
the system already speaks in terms of these layouts. The contract and the layout
enum were scaffolded for precisely this decomposition and then left unused when the
diagrams were written by hand instead.

Each diagram then collapses to a preset that converts a `ContextWindowSnapshot`
(the domain data) into `PositionedSegment[]`, picks a `layout`, and adds a caption
and tooltip. The conversion can reuse helpers that already exist —
`contextWindowIsHeadroomPart` and `contextWindowTokenTotal` at
`context/fixtures.ts:624-666`. Budget's only extra ingredient is setting `total` to
the snapshot's limit and adding an overflow marker. The single piece of genuine
domain logic left is `groupedParts` (`ContextGroupedStripDiagram.tsx:69`), which
groups parts before laying them out — a grouped strip is a strip whose cells are
themselves small strips.

> **A low-risk first step you can take before building the engine:** hoist
> `patternClass`, `styleName`, and `formatTokens` into `context/styles.ts` (next to
> the existing `resolveContextVisualStyle`). This deletes the seven- and five-way
> duplication immediately, changes nothing visible, and de-risks the larger engine
> extraction that follows.

The other molecule families are smaller versions of the same story:

- **`DataTable` is a special case of `MatrixGrid`.** Written `DataTable ⊂ MatrixGrid`
  ("⊂" is the subset symbol: everything `DataTable` does, `MatrixGrid` can do).
  `DataTable.tsx:22` renders cells through a `columns.cell(row)` render-prop and
  supports row selection; `MatrixGrid.tsx:64` renders cells through the richer
  `renderCell(payload)` contract. They are two grid engines whose adapters both lean
  on the same `renderCell`/`rowKey` helpers. `DataTable` could become a `MatrixGrid`
  preset, or at minimum share the cell contract instead of defining a parallel one.
- **The list family is one engine written three times.** `StepList.tsx:19`,
  `KeyPointList.tsx:27` (with its `normalizeItem` at line 17), and `CheckList.tsx:13`
  all take an `Array<Item | ReactNode>`, normalize each entry, render an
  ordered-or-unordered container, and optionally support a per-item marker and a
  selection. One `ItemList` engine with an item contract
  `{ id, marker, primary, secondary, meta, active, onSelect }` backs all three, and
  `DocumentListPanel.tsx:37`'s clickable list is the selectable variant of the same.
- **`MetadataGrid` and `KeyValueStrip` are one key/value engine.** They share an
  identical `{ key, value }` item shape (`MetadataGrid.tsx:4` vs `KeyValueStrip.tsx:4`)
  and differ only in layout (grid vs inline strip) and whether they offer a
  copy-to-clipboard affordance. One `KeyValueList` with a `layout` prop covers both.
- **The transcript cards share a style-resolution preamble.**
  `TranscriptMessageCard.tsx:58` and `AnnotationNoteCard.tsx:25` both begin by
  resolving a context style, computing a pattern class, and converting it to CSS
  variables. The cards themselves are genuinely domain-specific and should stay
  separate, but that shared five-line preamble belongs in a small hook
  (`useResolvedContextStyle`).

### 3.3 Organisms

Six families of near-duplicate organisms; each collapses to one engine plus a few
presets that supply domain defaults.

The cleanest of these — the one to do first — is the pair of **studio shells**.
`CmsShell` and `CourseStudioShell` are about ninety-five percent identical code.
Both build a header, wrap it in a `<SidebarShell sidebarWidth={188}>` with a
`<SidebarNav>` in the sidebar slot, and put the page content in a padding-controlled
div — and each ships its own near-identical `module.css`. They should be one
`StudioShell(sections, title, subtitle, headerSlot?, sidebarFooter?, contentPadding,
onNavigate)`, with `CmsShell` and `CourseStudioShell` reduced to presets that supply
the default navigation sections and copy (`cmsNav`, `courseStudioNav`).

The remaining five each follow the same shape — a repeated frame with a swappable
interior:

- **A `CollectionPanel` engine.** `MediaLibraryPanel` and `ArticleListPanel` share an
  identical frame: a `Panel` with a create action, a toolbar (`SearchField` + a
  filter + `Pagination`), an `EmptyState` when there is nothing to show, and a body.
  Only the body differs — a tile grid of assets versus a data table of articles.
  Factor the frame into a `CollectionPanel`; the bodies become a slot the presets
  fill.
- **A `MasterDetailShell`.** Five organisms are "a selector on the left, a detail view
  on the right": `HandoutDocumentShell`, `AssetDetailPanel`, `BookingPagePanel`,
  `ArticleEditorPanel`, `TranscriptWorkspacePanel`. But they are inconsistent — three
  use the `SplitPane` layout primitive and two hand-roll their own two-column CSS
  grid (`HandoutDocumentShell.module.css`, and the `withNotes`/`noNotes` classes in
  `TranscriptWorkspacePanel.module.css`). One `MasterDetailShell` with `left`/`right`
  slots and responsive collapse standardizes all five.
- **A `SelectableCardList` rail engine.** `AnnotationRailPanel`, `AnchoredCommentRail`,
  the molecule `DocumentListPanel`, and the event list inside `CalendarMonthPanel` all
  implement the same primitive: a header, then a list of `<button>`-wrapped cards, each
  calling `onSelect(id)` and reflecting a `selected` flag — and two of them even
  repeat the same inline "reset the button's default styling" CSS. One `role="listbox"`
  primitive parametrized by a card renderer replaces all four.
- **A shared context-diagram view renderer.** The `switch` that maps a
  `strip|stack|budget|treemap` choice to the matching diagram component is written
  three times, in `ContextDiagramPanel`, `CourseSlidePanel.tsx:renderSlideVisual`, and
  `ContextTurnPagerPanel`, each with its own copy of the legend-building helper. Extract
  one helper; the three organisms then differ only in their surrounding chrome — the
  frame around the content — which is a `Panel`, a `SlideShell`, and a pager respectively.
- **Bring the scheduling organisms up to the same bar.** `MeetingPollPanel`,
  `PollResultsPanel`, `CalendarMonthPanel`, `CalendarWeekPanel`, and `BookingPagePanel`
  are the newest organisms and the least standardized: they have no IR adapters yet,
  they each re-declare the same `MONTHS` array and `slotDate`/`formatSlot` date helpers
  (which belong in the shared `scheduling/` module), and some inline-style their event
  lists. Add the adapters, centralize the date helpers, and route the event list
  through the `SelectableCardList` above. While here, add the missing
  `RichArticle.widget.yaml`.

---

## Part 4 — Ad-hoc spec shapes, and how to unify them

This part is the catalog for smell 2. Where Part 3 found the same *component* logic
copied, Part 4 finds the same *idea* expressed as a slightly different one-off
type on each widget. The fix in every case is the same move the IR already made for
actions and cells (1.3): define one shared spec and one interpreter, and route every
site through it.

**1. One value-accessor instead of four.** An "accessor" is code that reads a value
out of an object given a path — for example, pulling `row.user.name` out of a row, or
filling `"${first} ${last}"` from a record. The codebase has *four* near-duplicate
implementations of this, and three of them re-implement the same `${field}`
string-template regular expression: `getPath`/`renderTemplate` in
`cellRenderers.tsx:109/131`; `lookupContext`/`interpolate` in `actions.ts:230/155`
(which additionally URL-encodes); `resolveValue` in `MatrixGrid.widget.tsx:26`; and
`styleBy.ts:17`, which reads only a single field and *silently cannot* follow a
dotted path the way the other three can — a latent bug waiting for someone to pass
`colorBy` a nested field. The *shape* of an accessor also recurs untyped across the
IR: `FieldCellSpec.field`, `RowKeySpec` (`{field}|{template}`), `MatrixValueSpec`
(`{mapField}|{template}`), `StyleBySpec.field`. Introduce one
`AccessorSpec = string | {field} | {path} | {template} | {mapField}` with a single
`resolveAccessor(spec, obj)` and one `interpolate()`, and route all five call sites
through it. This collapses four resolvers and three copies of the regex into one and
fixes the `styleBy` gap for free. This is the highest-leverage item in Part 4 —
"leverage" meaning benefit per unit of effort — because it unifies the most call
sites with the least behavioral risk.

**2. One selection spec instead of fourteen.** "Selection" — which item is currently
chosen — is expressed by about fourteen differently-named fields across the widgets:
`selectedKey`, `selectedCell`, `selectedPartId`, `selectedAnnotationId`,
`selectedCommentId`, `selectedItemId`, `selectedArticleId`, `selectedDocumentId`,
`selectedBlockId`, `selectedDateISO`, `selectedAssetIds[]`, and the parallel set
`activeItemId`/`activeId`/`activeAgendaItemId`. Only `MediaLibraryPanel` models the
single-versus-multiple distinction explicitly. Introduce one
`SelectionSpec { mode: "single" | "multi", keyField?, selected? }` reused across the
list and grid widgets.

**3. One list-item spec instead of ten.** At least ten item types are the same shape
apart from which optional fields they include: `AppNavItemSpec`,
`SidebarNavItemWidgetSpec`, `DocumentListItemWidgetSpec`, the `StepList` and
`KeyPointList` item types, `BreadcrumbItemWidgetSpec`, `TabListItemSpec`,
`SelectOptionSpec`, `KeyValueStripItem`, and `MetadataGridItemSpec`. Introduce one
generic base, `ListItemSpec<Extra> = { id, label: RenderableValue, icon?, badge?,
disabled? }`, that each specific item extends.

**4. Two `ctx` helpers that erase most adapter boilerplate.** Two lines of code
dominate the 81 IR adapters. The first is "render a node-or-string field over a list
of items" — `items.map(i => ({...i, label: ctx.renderValue(i.label)}))` — which
appears in nearly every list adapter. The second is "turn an action prop into a
callback, or `undefined` if it is absent" —
`props.onXAction ? (v) => ctx.dispatchAction(props.onXAction!, {value:v, componentType}) : undefined`
— which appears in 31 adapters, up to six times each. Add two methods to
`RenderContext`: `ctx.renderFields(obj, keys)` for the first and
`ctx.actionHandler(action, componentType, mapArgs?)` for the second. With those in
place, most **passthrough adapters** (adapters that only forward props) shrink to
almost nothing — and, as Part 5 explains, become mechanical enough to *generate*.

**5. Make action-prop naming consistent.** The convention is that any prop wiring an
action is named `onXAction` (36 props follow it), which lets tooling recognize action
props by their suffix. Two props break the convention: `DataTableWidgetProps.onRowSelect`
(`props.ts:430`) and `TabListWidgetProps.onChange` (`props.ts:568`). Rename them to
`onRowSelectAction`/`onChangeAction`, or adopt the rule "any prop whose type is
`ActionSpec` is an action regardless of name" — either makes the set machine-detectable.

Two smaller consistency fixes the review surfaced belong here as well. First,
`StyleBySpec` is only half-wired: the comment at `engines.ts:8-13` claims it unifies
coloring across `MatrixGrid`, `SegmentedBar`, `MonthGrid`, and `TimeGrid`, but in
fact only `MatrixGrid.colorBy` consumes it — the other three engines take a
pre-computed `styleKey` on each datum instead. Add `colorBy?: StyleBySpec` uniformly
to all four (keeping the per-datum `styleKey` as a fallback) so the design matches its
own documentation and the scheduling presets can stop precomputing keys. Second,
there are two template dialects living side by side: the structured, typed
`TemplateSpec`/`TemplatePartSpec` (`actions.ts:20`) and the older string form
`"${path}"` used by `TemplateCellSpec`, `RowKeySpec`'s `{template}`, and
`MatrixValueSpec`'s `{template}`. Fold the string cases onto the structured
interpreter (while keeping the string form as convenient sugar) so there is one
template engine, not two.

---

## Part 5 — The manifest as single source of truth

This is smell 3 — drift between hand-maintained mirrors — and it is the one place
where the IR tier and the DSL tier meet, which is why it gets its own part.

Recall from 1.1 that every widget has a manifest, `X.widget.yaml`. Open one, say
`MatrixGrid.widget.yaml`, and you find the widget described entirely as data: its
`type`, its `module` and `helper` (the DSL binding), its `props` type name, its
`reactComponent`, the list of `slots` (the props that hold renderable subtrees), and
the list of `actions` (the props that carry an `ActionSpec`). Every one of those
facts is *also* recorded, by hand, somewhere else. The `type` is a member of the
`RagWidgetType` union in `core.ts`. The `props` type is a member of the `WidgetProps`
union in `props.ts`. The registration is a hand-written import and array entry in
`defaultRegistry.ts`. The `module`/`helper` binding is a hand-written entry in a Go
map in `module.go`. The same widget is described in five files.

```
   RagWidgetType union        WidgetProps union         defaultRegistry.ts
   (core.ts, ~80-line literal) (props.ts, hand union)    (81 hand-written imports)
             ▲                        ▲                          ▲
             └──────────── all describe ────────────────────────┘
                         the same widget set
                                 │
              .widget.yaml manifests (79)      Go helper maps (module.go)
              type/module/helper/props/          helpers[name]=Type,
              slots[]/actions[]/reactComponent   cell/action/recipes
```

The problem is that **nothing reads the manifests.** A repository-wide search finds
no TypeScript that consumes `.widget.yaml` to produce anything; the Go tool at
`internal/widgetmanifest/` lints them but does not feed them into the runtime. So the
manifests are, in effect, hand-written prose that restates the adapters — and,
predictably, the copies have drifted. There are **81 adapters but only 79 manifests**:
`ContextStyleSwatch` and `RichArticle` have adapters with no manifest at all. And the
`helper:` names in the manifests (`matrixGrid`, `dataTable`) name Go functions that
the manifests have *no way to verify actually exist*, because the manifest validator
(`internal/widgetmanifest/validate.go`) only checks that a manifest agrees with its
React adapter, not that the Go DSL exports the promised helper. This is exactly the
mechanism by which the scheduling widgets ended up with manifests naming a `time.dsl`
module and `matrixGrid`/`segmentedBar`/`monthGrid`/`timeGrid` helpers that the Go
side never registered: the manifest promised them, and no check caught that the
promise was unkept.

The fix is to make the manifest the **single source of truth** — the one
authoritative description — and *generate* the mirrors from it instead of typing them
by hand. The manifest already carries everything a generator needs. From the set of
manifests, a build step could generate the `RagWidgetType` union, the `WidgetProps`
union, the `defaultRegistry` import list and per-module arrays, the passthrough
adapters (driven by `slots[]` → `renderFields` and `actions[]` → `actionHandler`,
using the two helpers from Part 4 #4, and falling back to a hand-written adapter only
where a widget is marked custom), the Go `helpers[name]=type` maps (or at least a Go
test that checks them against the manifests and reports a `helper_missing_in_go`
finding), and the typed `.d.ts` declarations (Part 6 #6). A single
`scripts/gen-widgets.mjs` plus a manifest-versus-adapter parity check would delete
hundreds of hand-maintained lines and make this class of drift impossible.

Note how Parts 4 and 5 reinforce each other: the `ctx.renderFields` and
`ctx.actionHandler` helpers proposed in Part 4 #4 are *exactly* the code a generated
adapter would call. Build those helpers first, and the code generator becomes a thin
templating step rather than a complicated one. That ordering dependency is why the
roadmap in Part 7 sequences the helpers before the codegen.

---

## Part 6 — The Go/Goja DSL: reaching parity, staying opinionated

The TypeScript IR has advanced past the Go DSL — several specs exist on the
TypeScript side with no DSL builder to produce them — and the DSL has some internal
seams of its own. These six changes bring it level and make it more data-driven,
without giving up the opinionated, curated feel that is its strength. ("Parity" here
means the DSL can express everything the IR can.)

1. **Close the IR gap — the highest-leverage, lowest-risk change, because it only
   adds.** The cell-builder object `cellObject` (`module.go:272`) stops at `constant`:
   there is no `cell.cycle`, no `cell.value`, no `styleBy` or `colorBy` builder. So a
   script author who wants an availability grid or a heatmap cannot express it through
   the friendly builders and must drop to the raw `component()` **escape hatch** — the
   low-level generic call you fall back to when the high-level helper is missing. Add
   these builders mirroring the existing `cellObject`/`actionObject`, reusing
   `buildPaletteStyleSet` (`module.go:393`) to produce the `ContextStyleSet` a
   `StyleBySpec` needs.
2. **Register the helpers the manifests already promise.** Add `matrixGrid`→`MatrixGrid`
   to `dataHelpers`, `segmentedBar`→`SegmentedBar` to `uiHelpers`, and a new `time.dsl`
   module with `monthGrid` and `timeGrid`. The frontend already knows `"time.dsl"`;
   only the Go side lags. This is the same reconciliation the scheduling handoff
   document calls out.
3. **Drive the helper maps from the manifests.** Rather than hand-typing the maps in
   #2 and risking the same drift again, generate or verify them from
   `internal/widgetmanifest.Discover`. This is the Go half of Part 5.
4. **Replace module-identity `if` chains with capability descriptors.**
   `runtime.install` (`module.go:236-269`) hard-codes special behavior per module with
   a chain of `if spec.name == …` — installing context-window style helpers for one
   module, the data grammar for another, the v2 builders for a third — and
   `typescript.go:73-116` mirrors the very same chain. Every new capability therefore
   requires editing two parallel switches by hand. Promote these to fields on the
   `moduleSpec` — an `installers []func` list, a `grammar bool`, a
   `recipes map[string]recipeFn` — so both the runtime and the declaration generator
   iterate over data instead of switching on names.
5. **Converge the classic grammar and the v2 model onto one lowering.** There are two
   near-identical compilers that turn a "collection" description into IR:
   `grammar.go:319-492` (the classic path) and `v2/spec/lower.go:87-251` (the typed
   v2 path), and the same duplication exists for record rows and action normalization.
   Make the classic option-bag verbs build a `v2spec.CollectionSpec` and call its
   `.ToNode()`, so there is one lowering ("lowering" = translating a higher-level
   description down to the concrete IR node form). This keeps the ergonomic classic
   surface that authors like, deletes the duplicate compiler, and lets the classic path
   inherit v2's validation diagnostics for free.
6. **Type the generated declarations.** Every generated helper is declared as
   `(props?: Props, …)` where `Props = Record<string, any>` (`typescript.go:26`), and
   every `cell`/`action` `options` argument is likewise untyped — so the rich,
   discriminated TypeScript unions (`CellSpec`, `StyleBySpec`, and the rest) are
   invisible to a script author, who gets no completion or checking. Emit real
   `import type { … } from "…"` lines and the actual union types (this is already done
   for the `data.v2.dsl` surface, just not the classic modules). This depends on #3;
   the runtime output stays open and unchanged — only the author-facing types tighten.

Two smaller items round this out. The recipe builders `mediaLibraryRecipe` and
`articleListRecipe` are each twenty-plus lines of `copyIfPresent` calls plus an
action-mapping loop, which a small declarative `recipeSpec{passthrough, actionMap,
component}` interpreted generically would collapse. And there is an asymmetry worth
resolving: `masterDetailTableRecipe` (`module.go:882`) accepts a `detail(row)` render
callback, while the v2 `masterDetail` derives its detail view automatically; a shared
"a render slot is a function that must return a validated node" convention would unify
the two.

---

## Part 7 — Prioritized roadmap

The work is ordered by *leverage-to-risk* — how much benefit a change delivers
relative to how likely it is to break something. Every item below is
behavior-preserving (it does not change what a user sees) and independently
shippable (it can land on its own without waiting for the others).

**Tier A — quick, low-risk hygiene. Do these first; each is an afternoon and none
can break rendering.**

- Hoist `patternClass`/`styleName`/`formatTokens` into `context/styles.ts` (Part 3.2).
- Add a shared `cx()` class-name helper; define canonical `Tone`/`Size` unions (Part 3.1).
- Add the two missing manifests (`ContextStyleSwatch`, `RichArticle`); reclassify the
  mis-namespaced pure atoms out of `cms.dsl` (Parts 3.1 and 5).
- Rename `onRowSelect`/`onChange` to end in `Action` (Part 4 #5).
- Add the `cell.cycle`/`cell.value`/`styleBy` DSL builders and register
  `matrixGrid`/`segmentedBar`/`monthGrid`/`timeGrid`/`time.dsl` (Part 6 #1–2).

**Tier B — shared specs and engine extractions. This is the core of the work: it is
where the duplication actually disappears.**

- Introduce `AccessorSpec` + `resolveAccessor` (Part 4 #1), then `SelectionSpec` and
  `ListItemSpec` (Part 4 #2–3).
- Add `ctx.actionHandler()` and `ctx.renderFields()` (Part 4 #4).
- Build the `SegmentEngine` and revive the `ContextDiagramSegment` contract, then
  migrate the five diagrams onto it (Part 3.2).
- Extract `ItemList`, `KeyValueList`, `GlyphBadge`, `ProportionTrack` (Parts 3.1–3.2).
- Extract `StudioShell`, `CollectionPanel`, `MasterDetailShell`, `SelectableCardList`,
  and the shared context-diagram view renderer (Part 3.3).

**Tier C — structural. Highest payoff, most design work, most worth doing carefully.**

- The manifest-as-source-of-truth code generator plus the parity lint (Part 5).
- Uniform `colorBy` and a single template engine (Part 4).
- The DSL changes: capability descriptors, classic↔v2 convergence, typed declarations
  (Part 6 #4–6).

**Sequencing.** Do Tier B #4 (the two `ctx` helpers) before Tier C (the code
generator), because the generator emits calls to exactly those helpers — building the
helpers first turns the generator into simple templating. The `SegmentEngine` work can
come before or after the accessor work, but the helper-hoist first step (Tier A) should
precede the engine so the engine starts from de-duplicated helpers.

---

## Part 8 — API & file reference index

Use this section as a lookup table while working. Every path is real and cited above.

### React components (`packages/rag-evaluation-site/src/components/`)
- Atoms: `atoms/{ContentStatusBadge,TranscriptRoleBadge,RatioBadge,MeterBar,ContextStudioNavIcon,ContextStyleSwatch,CycleCell,DateTile}/`
- Foundation: `foundation/{Text,Caption,StatusText}/`
- Molecules: `molecules/{ContextStripDiagram,ContextGroupedStripDiagram,ContextStackDiagram,ContextTreemap,ContextBudgetBar,ContextLegend,SegmentedBar,MatrixGrid,DataTable,StepList,KeyPointList,CheckList,MetadataGrid,KeyValueStrip,DocumentListPanel,TranscriptMessageCard,AnnotationNoteCard}/`
- Organisms: `organisms/{CmsShell,CourseStudioShell,MediaLibraryPanel,ArticleListPanel,HandoutDocumentShell,AssetDetailPanel,BookingPagePanel,ArticleEditorPanel,TranscriptWorkspacePanel,AnnotationRailPanel,AnchoredCommentRail,ContextDiagramPanel,ContextTurnPagerPanel,CourseSlidePanel,RichArticle}/`

### Widget IR (`packages/rag-evaluation-site/src/widgets/`)
- Files: `ir/{core,actions,cells,engines,props,index}.ts`, `WidgetRenderer.tsx`,
  `registry.ts`, `defaultRegistry.ts`, `cellRenderers.tsx`, `actions.ts`,
  `styleBy.ts`, `presets/scheduling.ts`.
- Key symbols: `RagWidgetType` (`core.ts`), `WidgetProps` union (`props.ts`),
  `ActionSpec` (`actions.ts`), `CellSpec`/`RowKeySpec` (`cells.ts`),
  `StyleBySpec`/`CycleCellSpec`/`ValueCellSpec`/`MatrixValueSpec` (`engines.ts`),
  `RenderContext` (`registry.ts`), `renderCell`/`rowKey` (`cellRenderers.tsx`),
  `MatrixCellPayload` (`molecules/MatrixGrid/MatrixGrid.tsx:18`).
- The dormant contract to revive: `ContextDiagramSegment`, `ContextDiagramView`
  (`context/types.ts:78,80`); the headroom helpers (`context/fixtures.ts:624-666`);
  `resolveContextVisualStyle` (`context/styles.ts:113`).

### Go/Goja DSL (`pkg/widgetdsl/`)
- `module.go` (`moduleSpec` :23, `moduleSpecs` :140, helper maps :35/82/86/109/126,
  `runtime.install` :236, `componentFactory` :565, `cellObject` :272,
  `actionObject` :318, `recipesObject` :495, `buildPaletteStyleSet` :393),
  `grammar.go` (`collectionVerb` :265, `recordVerb` :150), `v2_builders.go`,
  `v2/spec/{types,validate,lower}.go`, `typescript.go` (`TypeScriptModule` :14),
  `registrar.go`, `typescript_fixture_test.go`.
- Manifest tooling: `internal/widgetmanifest/{discover,validate,types}.go`.

## Open questions

These are genuine decisions that need an owner, not rhetorical questions.

1. **How far should codegen go?** Generate full adapters from the manifests, or only
   the mirrors (`RagWidgetType`, `WidgetProps`, `defaultRegistry`) plus a parity lint,
   and leave adapters hand-written? Recommendation: mirrors and lint first, adapters
   later once the `ctx` helpers exist.
2. **What becomes of `DataTable`?** Fold it into `MatrixGrid` as a preset, or keep it
   as a sibling that merely shares the cell contract? The two options differ in how
   much existing calling code has to change.
3. **Do the classic and v2 DSL paths converge or coexist?** Is `data.v2.dsl` meant to
   replace the classic grammar eventually, or live alongside it? Part 6 #5 assumes only
   the *backend* converges while both author-facing surfaces remain.
4. **Will `Tone`/`Size` unification change any pixels?** Collapsing divergent tone sets
   may map two currently-different tones to the same class somewhere; this needs a
   visual-regression pass before it ships.
5. **Who sequences the `time.dsl` reconciliation** relative to ticket
   `RAGEVAL-SCHEDULE-WIDGETS`, which the DSL colleague owns?

## References

- Companion guide: `reference/02-the-widget-system-a-new-intern-s-guide-to-the-architecture.md` — the textbook that teaches the architecture this document analyzes.
- Sibling ticket `RAGEVAL-SCHEDULE-WIDGETS` — the worked example of engine + contract + preset, and its DSL handoff guide.
- `packages/rag-evaluation-site/GUIDELINES.md` — the design-system rules.
- `reference/01-review-diary.md` — the method behind the five-area review these findings come from.
