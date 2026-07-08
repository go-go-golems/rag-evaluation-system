---
Title: Review Diary
Ticket: RAGEVAL-WIDGET-DECOMPOSITION
Status: active
Topics:
    - design-system
    - widget-ir
    - ui-dsl
    - react
    - frontend-architecture
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-06T21:04:27.528164927-04:00
WhatFor: ""
WhenToUse: ""
---

# Review Diary

## Goal

Review the existing widget library (`packages/rag-evaluation-site`) and its IR /
Goja-DSL layers, then produce an intern-facing analysis/design document on how to
**decompose domain widgets into reusable base engines + stable cell/segment
contracts + thin domain presets** — the pattern proven in the
`RAGEVAL-SCHEDULE-WIDGETS` ticket. **This is a review/design ticket only: no code
changes; the deliverable is the report.**

## Step 1: Setup and parallel review fan-out

Created the ticket and, per the "review only" constraint, launched five read-only
Explore agents to survey the layers concurrently rather than editing anything.

### Constraint (explicit)
- **Do not touch code.** Read, study, and write the report only. All agents are
  read-only; the only writes are the docmgr docs in this ticket + the reMarkable
  upload.

### Review fan-out (5 concurrent read-only agents)
1. Atoms + foundation — inventory, pure-vs-domain, decomposition smells.
2. Molecules — esp. the context-diagram family + DataTable; shared engine/segment
   contract opportunities.
3. Organisms — shells/panels/rails; near-duplicate families; engine-vs-preset.
4. Widget IR + rendering — the spec vocabulary, adapter boilerplate,
   manifest-as-source-of-truth, unification opportunities.
5. Goja DSL (`pkg/widgetdsl`) — module/helper/cell/action/recipe mechanics, v2
   spec, TS codegen; elegance/versatility opportunities.

### Next
- Synthesize the five reports into `design-doc/01-widget-library-decomposition-analysis-and-design.md`
  (intern-facing: prose + bullets + pseudocode + diagrams + API/file refs).
- Upload to reMarkable.

## Step 2: Five reviews in; synthesis + upload

All five read-only reviews completed and converged strongly. Wrote the intern
design doc and uploaded it. No code touched.

### Headline findings per layer
- **Atoms/foundation:** badge family (`ContentStatusBadge`/`TranscriptRoleBadge`/
  `StatusText`) = one enum→glyph engine ×3; `RatioBadge`⟂`MeterBar` fill;
  `Caption` is a preset of `Text`; `ContextStudioNavIcon` = SVG sprite + registry;
  pure atoms mis-namespaced to `cms.dsl`; `ContextStyleSwatch` missing manifest.
- **Molecules (flagship):** the 5 context diagrams re-implement one engine —
  `patternClass` duplicated in **7** files, `styleName`/`formatTokens` in 5, the
  ARIA/keyboard segment block in all 5. The segment contract
  (`context/types.ts:80 ContextDiagramSegment`) was scaffolded and never wired →
  finishing an abandoned refactor. Plus `DataTable ⊂ MatrixGrid`,
  `StepList/KeyPointList/CheckList → ItemList`, `MetadataGrid/KeyValueStrip → KeyValueList`.
- **Organisms:** `StudioShell` (CmsShell≈CourseStudioShell, ~95%), `CollectionPanel`,
  `MasterDetailShell` (5 organisms), `SelectableCardList` rail, shared
  `ContextDiagramView` (triplicated), scheduling parity gaps.
- **IR/rendering:** 4 accessor mini-languages + 3 `${}` regexes → one `AccessorSpec`;
  ~14 selection fields → `SelectionSpec`; ~10 item types → `ListItemSpec`; ~30×
  adapter boilerplate → `ctx.actionHandler`/`renderFields`; `.widget.yaml`
  manifests read by nothing (81 adapters vs 79 manifests) → manifest-as-SoT codegen.
- **Goja DSL:** no `cell.cycle`/`value`/`styleBy` builders; scheduling helpers +
  `time.dsl` unregistered; `if spec.name==` special-casing → capability descriptors;
  classic grammar duplicates v2 lowering; untyped `.d.ts` props.

### Synthesis
- The five reviews converge on one lens (**engine + contract + preset**) and three
  smells (duplicated engine logic, ad-hoc spec shapes, hand-maintained-mirror
  drift). Structured the doc as: pattern → per-layer catalog → cross-cutting IR
  specs → manifest-as-source-of-truth → DSL → prioritized A/B/C roadmap.
- Uploaded to reMarkable: `/ai/2026/07/06/RAGEVAL-WIDGET-DECOMPOSITION` → `OK: uploaded`.

### Method note (for reproducibility)
- Findings were gathered by five concurrent read-only Explore agents (one per
  layer). Counts (e.g. "7 files") are grep-reproducible against the cited paths.

## Related

- `design-doc/01-widget-library-decomposition-analysis-and-design.md`
- Sibling ticket `RAGEVAL-SCHEDULE-WIDGETS` (the decomposition pattern this generalizes).
