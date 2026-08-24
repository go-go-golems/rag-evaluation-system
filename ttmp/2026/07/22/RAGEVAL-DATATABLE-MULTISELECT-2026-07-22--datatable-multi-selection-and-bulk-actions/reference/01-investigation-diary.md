---
Title: Investigation diary
Ticket: RAGEVAL-DATATABLE-MULTISELECT-2026-07-22
Status: active
Topics:
    - react
    - widget-dsl
    - design-system
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Chronological evidence for the DataTable multi-selection design and delivery."
LastUpdated: 2026-07-22T16:51:51.989568093-04:00
WhatFor: "Record how the reusable React, Widget IR, and widget.dsl boundaries were mapped before implementation."
WhenToUse: "Read before beginning or reviewing the DataTable multi-selection implementation."
---

# Diary

## Goal

Capture the evidence-backed design work for accessible DataTable multi-selection, including the reason focus movement must not clear the checked row set.

## Step 1: Map the table, renderer, and DSL ownership boundaries

The repository already has the beginnings of multi-selection in the generic Widget DSL data namespace, but not an end-to-end table capability. The investigation separated what is actually rendered by the reusable React molecule from what is merely serializable in Goja, preventing a plan that would add DSL syntax without visible behavior.

The proposed design is checkbox-first with Shift range enhancement and one aggregate bulk-action bar. Focus remains a one-row navigation state, while selection is a controlled set of keys; therefore ordinary keyboard movement preserves the selection rather than clearing it.

### Prompt Context

**User prompt (verbatim):**

<pre>Create a ticket in rag-evaluation-system for the multiselect (both with shift click option with multiple selected bars, or with checkboxes). In the first case, moving with the keyboard would clear th e multiselect I guess? Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.


Make sure to follow the rag-evaluation-system widget-dsl and react conventions / GUIDELINES</pre>

**Assistant interpretation:** Create and deliver an intern-facing, evidence-backed implementation design for reusable DataTable multi-selection, covering checkboxes, Shift ranges, keyboard semantics, Widget IR, and widget.dsl.

**Inferred user intent:** Establish a safe shared foundation for future bulk table operations without creating a one-off Tracker behavior or losing selection during keyboard navigation.

### What I did
- Read `AGENTS.md` and `packages/rag-evaluation-site/GUIDELINES.md` before design work.
- Inspected `DataTable.tsx`, its Widget adapter, Widget IR props/actions, typed collection lowering, grammar lowering, v3 selection serialization, styles, stories, and existing CheckboxRow atom.
- Created ticket `RAGEVAL-DATATABLE-MULTISELECT-2026-07-22` and wrote the detailed design guide.

### Why
- A multi-row UI crosses the reusable React molecule, JSON-compatible Widget IR, Goja authoring surface, and server-action context. Each boundary needs one stable contract.

### What worked
- Concrete code evidence showed the current table has a single `selectedKey`, while `widget.data.selection({ mode: "multi" })` is presently only a validated generic serialized object.
- The package guidelines supplied a clear implementation order: React component and stories first, Widget IR second, DSL lowering last.

### What didn't work
- No implementation command failed. The initial assumption that a generic multi-selection DSL value implied DataTable support was disproved by the adapter and lowering inspection: neither consumes the value.

### What I learned
- Existing `followFocus` is suitable for single master-detail navigation but cannot represent a persistent multi-selection. A separate focus key, selection set, and range anchor are required.

### What was tricky to build
- The difficult part is semantic, not rendering checkboxes: the current row callback means “select/open this one row,” whereas bulk selection must mean “change this set.” The design resolves the ambiguity with a discriminated React API and a keys-only Widget action context.

### What warrants a second pair of eyes
- Review the proposed Widget DSL builder name and whether a consumer needs cross-page selection. The latter would require an explicit server-side token rather than extending the in-memory visible-row set silently.

### What should be done in the future
- Implement phases in the design document; do not ship product-specific bulk behavior before reusable component, IR, and DSL tests pass.

### Code review instructions
- Begin with `DataTable.tsx:29-208`, then trace `DataTable.widget.tsx:6-70`, `props.ts:460-485`, and `spec/lower.go:235-276`.
- Validate eventual work with Biome, package typecheck/tests/Storybook, and `go test` for widgetdsl packages.

### Technical details

```text
focusKey: one roving keyboard target
selectedKeys: controlled set for bulk operations
anchorKey: local range-selection origin

Arrow movement -> focus only
Space / checkbox -> mutate selection
Shift-click / Shift+Arrow -> inclusive visible range
```

## Step 2: Validate ticket documentation and deliver it to reMarkable

The documentation was checked against the ticket workspace conventions, then rendered as a ToC-enabled bundle for reMarkable. The ticket is intentionally a design deliverable rather than an implementation change, so no production source files were edited.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Finish the ticket bookkeeping and deliver the requested technical guide to reMarkable.

**Inferred user intent:** Make the design available as a durable, reviewable reference for an intern before code work begins.

### What I did
- Related the primary guide to the React table, Widget adapter/IR, typed lowering, and v3 runtime files.
- Ran `docmgr doctor` after resolving ticket vocabulary.
- Ran a bundle upload dry-run, uploaded the document bundle, and verified the remote listing.

### Why
- The ticket must be findable from the source files it describes and usable outside the repository as a review/reference artifact.

### What worked
- The ticket doctor completed without findings, and the bundled PDF upload succeeded.

### What didn't work
- The first real `remarquee upload bundle` failed after dry-run because the diary quoted blank lines as literal `\\n`, which Pandoc passed through to LaTeX: `Error producing PDF. ! Undefined control sequence. l.1166 remarkable.\\n` (exit status 43). Replacing those literal escapes with a `<pre>` block preserved the verbatim prompt and produced a valid PDF.

### What I learned
- The current multi-selection representation is a useful future integration point, but it needs a table-specific action/keyboard contract rather than a thin visual adapter. ReMarkable/Pandoc rendering also treats literal backslash-newline sequences in prose as TeX syntax, so verbatim multi-line prompts need actual newlines or safe markup.

### What was tricky to build
- The guide has to distinguish observed code from proposed API. File references identify present behavior; all new public contracts are marked as proposed decision records and phased work.

### What warrants a second pair of eyes
- UX/accessibility review should approve the final bulk-action placement and whether Shift ranges union or replace selections.

### What should be done in the future
- Assign an implementation owner and turn the phase checklist into code tasks once the open product decisions are resolved.

### Code review instructions
- Review the decision records and API sketches first, then compare each proposed layer to its referenced source file.
- Open the reMarkable bundle and use the ToC to inspect the current-state, API, and test-plan sections.

### Technical details

```text
Ticket bundle destination:
/ai/2026/07/22/RAGEVAL-DATATABLE-MULTISELECT-2026-07-22
```

## Step 3: Record accepted product decisions and implement shared multi-selection

The product owner accepted every recommended interaction default and selected Archive and Tag as the first consumer actions. The reusable implementation now adds controlled multi-selection to the React molecule, a single toolbar above the table, Widget IR/action-context wiring, and a typed v3 `table.multiSelect` lowering surface.

This changes the shared system only. A product consumer must still provide its deliberate enter/exit bulk-mode control and server-side Archive/Tag handlers; the shared package dispatches selected keys but never performs a domain mutation itself.

### Prompt Context

**User prompt (verbatim):** "6. exclusive. Record decisions, the implement"

**Assistant interpretation:** Accept exclusive bulk/detail modes, record the confirmed product decisions, and implement the shared feature.

**Inferred user intent:** Move from design to a reusable, predictable multi-select foundation for future Archive and Tag workflows.

### What I did
- Added controlled `multiSelection`, checkbox/range interaction, keyboard toggling, and a bulk toolbar to `DataTable`.
- Added serializable Widget IR props and adapter contexts containing `selectedRowKeys`, count, reason, and bulk action ID.
- Added `TableSpec.MultiSelection`, lowering, `table.multiSelect`, generated TypeScript declaration text, descriptor documentation, and lowering coverage.

### Why
- The React-first package guidelines require stable component behavior before Widget IR and DSL authoring are made available.

### What worked
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `go test ./pkg/widgetdsl/...` passed after updating the generated Widget DSL API reference.

### What didn't work
- `go test ./pkg/widgetdsl/...` initially failed with `embedded API help descriptor reference is stale; regenerate ../xgoja/providers/widgetsite/doc/05-widget-dsl-v3-api-reference.md from WidgetV3APIReferenceMarkdown`. Updating that generated reference for `TableBuilder.multiSelect` resolved the failure.

### What I learned
- The shared table can provide Archive/Tag-ready key context without importing any product backend, preserving the design-system API-free boundary.

### What was tricky to build
- Shift state is available from a click event rather than React's checkbox change event. The checkbox click handler owns range/toggle selection while the controlled change handler remains intentionally inert, preventing a duplicate selection transition.

### What warrants a second pair of eyes
- Review the public builder signature and whether an actual product requires selection persistence across pagination; this v1 intentionally normalizes selection to visible rows.

### What should be done in the future
- Publish a compatible `@go-go-golems/rag-evaluation-site` release and wire Upwork Triage’s explicit bulk mode plus Archive/Tag server handlers against it.

### Code review instructions
- Begin with `DataTable.tsx` for interaction semantics, then `DataTable.widget.tsx`, `props.ts`, and `pkg/widgetdsl/spec/lower.go`.
- Run package typecheck and `go test ./pkg/widgetdsl/...`.

### Technical details

```js
collection.table((table) => table.multiSelect(["job-a"], {
  actions: [{ id: "archive", label: "Archive", action: act.server("archive-jobs") }],
}));
```
