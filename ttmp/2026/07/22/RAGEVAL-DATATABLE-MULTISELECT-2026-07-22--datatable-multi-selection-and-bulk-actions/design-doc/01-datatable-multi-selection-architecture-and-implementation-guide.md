---
Title: DataTable multi-selection architecture and implementation guide
Ticket: RAGEVAL-DATATABLE-MULTISELECT-2026-07-22
Status: active
Topics:
    - react
    - widget-dsl
    - design-system
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.tsx
      Note: |-
        Reusable React table selection, focus, commands, and markup baseline
        Implemented controlled checkbox/range/keyboard multi-selection and toolbar
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.widget.tsx
      Note: |-
        Widget IR-to-React adapter and action context boundary
        Implemented Widget action dispatch for selected key context
    - Path: repo://packages/rag-evaluation-site/src/widgets/ir/props.ts
      Note: JSON-compatible DataTable Widget IR contract
    - Path: repo://pkg/widgetdsl/spec/lower.go
      Note: Typed collection lowering from selection/table spec to Widget IR
    - Path: repo://pkg/widgetdsl/spec/types.go
      Note: Typed multi-selection and bulk-action model
    - Path: repo://pkg/widgetdsl/v3.go
      Note: |-
        Existing generic multi-selection serialization behavior
        Implemented table.multiSelect builder
ExternalSources: []
Summary: Intern-facing design for accessible multi-row DataTable selection, bulk actions, Widget IR, and widget.dsl lowering.
LastUpdated: 2026-07-22T16:51:51.344584385-04:00
WhatFor: Plan and implement reusable DataTable multi-selection without conflating focus, selection, or application-side bulk effects.
WhenToUse: Use before adding checkbox or range selection to a Widget DSL table, or before a product introduces bulk table actions.
---



# DataTable multi-selection architecture and implementation guide

## Executive summary

`DataTable` is the reusable React molecule behind Widget DSL tables. Today it has one selected key, one focused row, and a row callback. That makes master-detail selection clear, but it cannot express a checked set, a range anchor, a select-all control, or an action whose payload represents more than one row. The current generic Widget DSL `data.selection({ mode: "multi" })` helper is a serializable value only; it is not connected to `DataTable` lowering or rendering.

This ticket proposes a deliberate two-layer feature. First, stabilize a controlled React `DataTable` multi-select API with a visible checkbox column, a single bulk-action bar, and independent focus/selection state. Second, expose the same semantic contract in Widget IR and `widget.dsl`, with explicit selected row keys and an action context containing `selectedRowKeys`. The recommended interaction includes checkbox toggles and Shift-click range extension. Keyboard arrow movement **must not clear selection**: focus is navigation state, while selection is an explicit set. Space toggles the focused row; Shift+Arrow extends a range; Escape clears the set. This gives keyboard users an equivalent workflow without accidental destructive selection changes.

The proposal intentionally does not add product-specific bulk behaviors to the design system. A consuming product supplies a named action and validates IDs server-side. The reusable package only renders rows, collects intent, and dispatches a JSON-compatible action context.

## Problem statement and scope

A user needs to select several table rows before performing a single bulk operation. Common pointer paths are clicking individual checkboxes and Shift-clicking an inclusive row range. A compact bulk-action bar should appear once at least one row is selected, show the count, and expose caller-supplied actions. The table must work with a mouse, touch input, and keyboard.

The scope is deliberately broader than checkbox visuals:

- reusable React component behavior and accessible markup;
- Widget IR props and adapter dispatch context;
- declarative `widget.dsl` authoring and typed Go lowering;
- generated TypeScript/API documentation and tests;
- stories and test cases that make keyboard rules reviewable.

It excludes persistence of selection across pagination, server-side bulk-operation semantics, application authorization, and a product-specific bulk toolbar. Those are consumers' responsibilities. The first version operates on the currently rendered `rows`; filtering or paging resets unavailable keys from the controlled selection supplied by the caller.

## Terms and invariants

- **Focus** is the one roving-tabindex target that receives keyboard events. It is transient client interaction state.
- **Selection** is the set of row keys marked for bulk operation. It is controlled input/output data, not an implication of focus.
- **Anchor** is the key from which a Shift range extends. It is transient client interaction state and is not sent to a server.
- **Visible row set** is the ordered `rows` array currently rendered. A range is inclusive and cannot cross a filtered/paginated boundary.
- **Bulk action** is an existing Widget `ActionSpec` dispatched once with an ordered, deduplicated list of visible selected keys.

The component must preserve these invariants:

1. Every item in `selectedKeys` is a current row key after the component reconciles props.
2. Changing focus with ArrowDown, ArrowUp, `j`, or `k` never changes the selection set by itself.
3. A plain checkbox click toggles one key. Shift-click selects the inclusive range from the anchor to the target; it does not depend on OS-specific Ctrl/Command behavior.
4. The header checkbox selects all current visible rows or clears all current visible selections. Its indeterminate state is derived, never stored.
5. A row click in multi mode does not invoke the legacy `onRowSelect` single-row behavior. Detail navigation and bulk selection remain separate intentional modes.
6. Action payloads contain keys, not row records. The server must re-authorize and re-read records by key.

## Current-state architecture

### React molecule

`packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.tsx` is a generic presentational table. Its public props use `selectedKey?: string | null`, `onRowSelect?: (row) => void`, and optional keyboard commands (lines 29-40). Internally, it maintains one `focusedKey` and restores DOM focus after navigation (lines 64-94). A `followFocus` mode calls `onRowSelect` while moving (line 93). Rows compare `selectedKey === key`, expose one `aria-selected`, and call the single row callback on click (lines 164-198).

This is an important boundary: the table is package-owned and API-free. It cannot reach an RTK store, route, or backend to own selection. The design-system guidelines require the correct molecule layer, CSS Modules, foundation/theme tokens, `data-rag-*` identity attributes, and stories before Widget IR support. See `packages/rag-evaluation-site/GUIDELINES.md` sections “Molecules” and “Widget IR / WidgetRenderer rules.”

### Widget IR and renderer adapter

`packages/rag-evaluation-site/src/widgets/ir/props.ts:460-485` mirrors the single-select surface as `selectedKey`, `onRowSelect`, keyboard fields, commands, style rules, and empty message. `DataTable.widget.tsx` maps those props into React and dispatches row action context `{ row, rowKey, componentType: "DataTable" }` (lines 41-66). It currently has no selected-key list, selection-change callback, bulk action declaration, or selection context.

`widgets/actions.ts` is the transport boundary. `dispatchWidgetAction` resolves payload templates and posts `{ payload, context }` for server actions. The context is JSON-shaped; extending it with `selectedRowKeys: string[]` is safe if the adapter constructs it locally and the action resolver treats it like the existing `rowKey` context.

### Go specification and Widget DSL lowering

The typed collection model in `pkg/widgetdsl/spec/types.go:223-322` has a URL-param-only `SelectionSpec` containing one `Value`, and its table model has a singular `RowSelect` plus row commands. In `pkg/widgetdsl/spec/lower.go:235-276`, a selected collection lowers to one `selectedKey` and one `onRowSelect` navigation/action. This is why adding only React checkboxes would leave declarative collection authors without a usable contract.

There is also a low-level legacy-style grammar path in `pkg/widgetdsl/grammar.go:327-365`; it emits the same singular `selectedKey` and `onRowSelect`. The v3 runtime does validate and serialize `widget.data.selection({ mode: "multi", keyField, selected })` in `pkg/widgetdsl/v3.go:1540-1569`, and `module_test.go` verifies that standalone serialization. It is not consumed by `CollectionSpec` or the DataTable adapter. The design must either connect this representation or replace it with a table-specific configuration; it must not pretend existing `data.selection` enables a table.

### Existing visual/test conventions

`DataTable.stories.tsx` has one populated single-selected story. `DataTable.module.css` contains the selected, focused, sortable, command-help, and semantic-tone anatomy. `CheckboxRow.tsx` is an existing package atom whose input semantics can inform checkbox handling, but it is a label-row control and should not be forced into a table-cell layout. Widget registry coverage is demonstrated in `WidgetRenderer.domain-registry.stories.tsx`.

## Gap analysis

| Needed capability | Current behavior | Gap |
|---|---|---|
| Several selected rows | one `selectedKey` | requires `selectedKeys` and callback |
| Range selection | no anchor or modifiers | requires ordered-key range helper |
| Visible selection affordance | selected row color only | requires header/row checkbox cells and count |
| Bulk operation dispatch | `onCommand(command, row)` | requires action context for a key list |
| Keyboard parity | arrows may call `onRowSelect` | focus and selection must be decoupled |
| Declarative authoring | URL-param single selection | requires explicit Widget IR/DSL contract |
| Review surface | one story | requires selection, overflow, disabled, and keyboard stories |

## Proposed architecture

### 1. React-first controlled API

Add an explicit, opt-in `selection` prop rather than overloading `selectedKey`. Keep the existing single-selection props unchanged to avoid silently changing master-detail tables.

```ts
export interface DataTableMultiSelection {
  mode: "multi";
  selectedKeys: readonly string[];
  onSelectionChange: (nextKeys: string[], reason: DataTableSelectionReason) => void;
  bulkActions?: readonly DataTableBulkAction[];
  ariaLabel?: string;
}

type DataTableSelectionReason =
  | "toggle"
  | "range"
  | "selectAll"
  | "clearAll"
  | "keyboardToggle"
  | "keyboardRange"
  | "clear";

interface DataTableBulkAction {
  id: string;
  label: ReactNode;
  danger?: boolean;
  disabled?: boolean;
  onInvoke: (selectedKeys: readonly string[]) => void;
}
```

`DataTableProps` accepts either legacy single selection or the multi-selection object, never both. At runtime and in TypeScript, reject/diagnose contradictory props. The component derives a `Set` from `selectedKeys`, filters it against current `keys`, and emits a stable document-order array. It owns only `focusedKey` and `anchorKey` locally.

Render a leading checkbox header and one checkbox cell per row only in multi mode. The header checkbox receives a visually hidden label, has `checked` when every visible row is selected, and uses the DOM `indeterminate` property when some but not all visible rows are selected. A single `BulkActionBar` region above the table appears when the derived selection is nonempty. It says “N selected,” offers caller-defined actions, and includes “Clear selection.” This is one bar for the aggregate selection, not a separate bar per selected row.

### 2. Pointer and keyboard state machine

Checkboxes provide the discoverable baseline. Shift-click is an enhancement, not the only way to multiselect. Do not require Ctrl/Command click because it conflicts with browser/platform conventions and harms touch discoverability.

```text
plain checkbox click(key):
  next = toggle(selected, key)
  anchor = key
  emit(next, "toggle")

shift checkbox click(key):
  start = anchor if anchor is visible else focusedKey if visible else key
  next = union(selected, orderedInclusiveRange(keys, start, key))
  anchor = start
  emit(next, "range")

ArrowDown / ArrowUp:
  focusedKey = adjacent visible key
  // selection unchanged

Space:
  toggle(focusedKey)
  anchor = focusedKey
  emit(next, "keyboardToggle")

Shift+ArrowDown / Shift+ArrowUp:
  focusedKey = adjacent visible key
  next = union(selected, range(anchor or oldFocus, focusedKey))
  emit(next, "keyboardRange")

Escape:
  emit([], "clear")
```

A normal row click in multi mode should move focus but not toggle. This prevents accidental selection during inspection. The checkbox is the explicit pointer control. If product research later demonstrates that row-click toggle is preferred, add it as a named option and test it; do not make an implicit behavior change.

**Answer to the keyboard question:** no, moving with the keyboard should not clear a multi-selection. The current `followFocus` contract is fundamentally a single-selection/master-detail convenience. Multi mode must disable/ignore `followFocus`, retain selection across movement, and use Space/Shift+Arrow for selection changes. This makes focus visibly distinct from checked rows and avoids turning a review sweep into a destructive selection reset.

### 3. Widget IR contract

Once the React API and stories are stable, add JSON-compatible props to `DataTableWidgetProps`:

```ts
selection?: {
  mode: "multi";
  selectedKeys: string[];
  onChange?: ActionSpec;
};
bulkActions?: Array<{
  id: string;
  label: RenderableValue;
  danger?: boolean;
  disabled?: boolean;
  action: ActionSpec;
}>;
```

The Widget adapter translates selection state into the React controlled API. On any change, it dispatches the optional `onChange` action with:

```json
{
  "selectedRowKeys": ["job-12", "job-15"],
  "selectedCount": 2,
  "selectionReason": "range",
  "componentType": "DataTable"
}
```

A bulk action dispatches the same context plus `bulkActionId`. Do not include every selected row object: it duplicates data, enlarges requests, and tempts server handlers to trust stale client records. A product can use `widget.bind.context("selectedRowKeys")` in the action payload. Server handlers must validate IDs, authorization, and permitted state transitions afresh.

### 4. widget.dsl / typed lowering contract

Add a table-specific selection declaration to the v3 collection builder and matching typed spec. The exact builder name should follow existing collection API patterns, but a concrete target is:

```js
collection.table(table => table
  .multiSelect({
    selected: selectedJobIds,
    onChange: act.server("triage-selection-changed", {
      payload: { jobIds: bind.context("selectedRowKeys") },
    }),
    actions: [
      { id: "archive", label: "Archive", danger: true,
        action: act.server("archive-jobs", {
          payload: { jobIds: bind.context("selectedRowKeys") },
        }) },
    ],
  })
)
```

The Go `TableSpec` gains a `MultiSelection *MultiSelectionSpec`; lowering must emit `selection` and `bulkActions` only when present. Validate these invariants before IR emission:

- `selected` contains strings and no duplicates;
- every bulk action has nonempty id, label, and valid action;
- `MultiSelection` cannot coexist with a URL-param `SelectionSpec` or `followFocus` keyboard mode;
- duplicate action IDs are rejected;
- no key is assumed trustworthy merely because it appears in the selection list.

Generated TypeScript declarations, v3 descriptor metadata, Go tests, and API reference docs must be updated from the same source of truth. The current standalone `data.selection` helper may remain generic, but the guide recommends not making it the table API until the collection builder consumes it directly; a table-specific structure has clearer validity rules and avoids accidental ambiguity with other selectable controls.

### Data flow diagram

```text
JS verb / typed CollectionSpec
          |
          | lower to JSON-compatible props
          v
Widget IR DataTable { selection, bulkActions }
          |
          | DataTable.widget.tsx maps props + dispatches actions
          v
React DataTable molecule
  focusKey + anchorKey (local)       selectedKeys (controlled)
          |                                  |
          +--------- pointer/keyboard -------+
                                             v
                              onSelectionChange / bulk action
                                             |
                                             v
Widget action context { selectedRowKeys, selectedCount, reason }
                                             |
                                             v
product server action validates IDs and executes domain operation
```

## Decision records

### Decision: checkbox-first interaction with Shift range enhancement

- **Context:** The request allows either Shift-click selection with multiple selected bars or checkboxes. A reusable table must work on pointer, touch, and keyboard devices.
- **Options considered:** range-only row click; checkbox-only; checkbox plus Shift-click range.
- **Decision:** Render checkbox cells and a select-all header checkbox; additionally support Shift-click/Shift+Arrow range selection.
- **Rationale:** Checkboxes make persistent selection visible and discoverable. Range extension accelerates large contiguous selection without becoming the only mechanism.
- **Consequences:** The table gets one structural column and must precisely implement indeterminate/header state and modifier behavior.
- **Status:** proposed.

### Decision: focus and selection are separate state machines

- **Context:** Existing `followFocus` calls single-row selection while navigating. That would overwrite or blur a set selection.
- **Options considered:** clear selection on every move; replace selection with focused row; preserve selection and add explicit toggle keys.
- **Decision:** Preserve selection on ordinary movement; use Space and Shift+Arrow to mutate it.
- **Rationale:** It avoids accidental loss, matches roving-focus accessibility patterns, and answers the requested keyboard concern directly.
- **Consequences:** Multi mode cannot use `followFocus`; visual styling must distinguish focus ring from selected fill.
- **Status:** proposed.

### Decision: keys-only action context

- **Context:** Widget actions already serialize context/payload to browser/server boundaries.
- **Options considered:** pass full row objects; pass keys and count; mutate a global selection store.
- **Decision:** Pass ordered `selectedRowKeys`, count, reason, and action ID only.
- **Rationale:** The API stays JSON-compatible, bounded, composable, and does not make stale browser records authoritative.
- **Consequences:** Product handlers perform a keyed lookup/authorization step; examples and docs must state this safety requirement.
- **Status:** proposed.

### Decision: React API and stories before Widget DSL

- **Context:** package guidelines explicitly require React-first stabilization before Widget IR support.
- **Options considered:** add raw IR props first; build an app-specific bulk control; stabilize the molecule first.
- **Decision:** Implement/test/stories for `DataTable` first, then adapter/IR, then Go/DSL lowering.
- **Rationale:** It prevents JSON DSL design from encoding an untested DOM interaction model.
- **Consequences:** The feature lands in phases and Widget DSL support follows a stable component API.
- **Status:** proposed.

## Confirmed product decisions (2026-07-22)

The product owner confirmed the recommended defaults. These are now accepted constraints for the implementation:

- Shift ranges **union** with independently checked rows.
- The aggregate bulk-action bar appears **above** the table.
- Row content does not toggle selection; explicit checkboxes and Space do.
- Selection is limited to the **currently visible** filtered/page rows.
- The first consumer actions are **Archive** and **Tag**.
- Bulk selection and master-detail navigation are **exclusive modes**. A table is either in normal single-row/detail mode or in bulk-select mode; it does not make one click mean both “open” and “toggle.”

### Decision: exclusive detail and bulk modes

- **Context:** A single-row master-detail callback and a checked row set assign incompatible meanings to ordinary row interaction.
- **Options considered:** combine detail opening and checked-set mutation in one mode; use explicit bulk mode and preserve normal detail mode.
- **Decision:** Use explicit, mutually exclusive modes.
- **Rationale:** It keeps row click, Enter, focus, Space, and bulk actions unambiguous for pointer and keyboard users.
- **Consequences:** A consumer needs a deliberate “Select” entry/exit affordance; it gains predictable Archive/Tag action context.
- **Status:** accepted.

## Implementation plan

### Phase 0: confirm product requirements and name the public API

1. Confirm whether all candidate uses want client-controlled selected keys or URL persistence. This guide assumes controlled in-memory selection.
2. Confirm whether a selected row can still open detail. The proposed mode uses explicit row/detail actions rather than conflicting row click.
3. Add acceptance criteria for keyboard and screen-reader behavior before writing component code.

### Phase 1: React DataTable molecule

1. Extend `DataTableProps` with a discriminated selection union; retain `selectedKey` behavior for legacy callers.
2. Factor pure helpers in `DataTable.tsx` or a colocated testable module:
   - `normalizeSelection(keys, selectedKeys)`;
   - `toggleKey`;
   - `rangeKeys(orderedKeys, start, end)`;
   - `headerCheckboxState`.
3. Add `anchorKey` and keep `focusedKey` independent.
4. Render a checkbox column only in multi mode. Use a ref/effect for native `input.indeterminate`.
5. Render one `data-rag-component="DataTableBulkActions"` region with count and caller actions; retain the root `data-rag-component="DataTable"` attribute.
6. Update `DataTable.module.css` using existing `--mac-*`/`--rag-*` tokens. Do not hardcode new typography/color literals; avoid inline layout styles.
7. Do not import CheckboxRow if its label-row anatomy is unsuitable; use a native checkbox with accessible labels and local CSS anatomy, or extract a purpose-built atom only if it will be reused.

### Phase 2: package review surface and tests

1. Extend `DataTable.stories.tsx` with: default multi-select, partial/all selection, long/overflow labels, disabled bulk action, empty table, and row tone combined with checked/focus state.
2. Add interaction tests for mouse checkbox toggling, Shift-click range, select all/clear, Space toggle, Shift+Arrow range, Escape clear, focus restoration after rows change, and no selection change on plain arrow navigation.
3. Validate screen-reader markup: header/row checkbox names, `aria-checked` semantics, count announcement via a polite status region if the final accessibility review requires it, and only one tab stop through rows.
4. Run package typecheck, test suite, Storybook/build, and Biome per `AGENTS.md`.

### Phase 3: Widget IR and React adapter

1. Extend `DataTableWidgetProps` in `src/widgets/ir/props.ts` with the narrow JSON contract.
2. Map it in `DataTable.widget.tsx`; construct contexts uniformly for change and bulk action dispatch.
3. Update WidgetRenderer/domain-registry stories with a concrete DataTable multi-selection IR node and an action spy assertion.
4. Add action resolver tests proving `widget.bind.context("selectedRowKeys")` is preserved as an array, not stringified.

### Phase 4: typed spec and widget.dsl

1. Add `MultiSelectionSpec` and bulk action spec to `pkg/widgetdsl/spec/types.go`.
2. Add validation in `spec/validate.go`; update `spec/lower.go` to emit the exact IR shape.
3. Add v3 builder method, descriptors, TypeScript declaration generation, documentation, and generated API tests in the same change.
4. Decide whether the legacy grammar also receives authoring support. If retained, make it emit the identical IR shape; do not create two runtime contracts.
5. Add Go tests for valid lowering and every incompatible configuration.

### Phase 5: first consumer and release discipline

1. Add a minimal demo/fixture or a real product consumer only after generic validation passes.
2. Keep domain bulk handlers outside `rag-evaluation-site`; validate keys server-side.
3. Release the package only through the repository’s GitHub Actions Trusted Publishing flow if a consumer needs a new published version. Do not run local `npm publish`.

## Test and validation matrix

| Layer | Test cases | Evidence target |
|---|---|---|
| Pure selection helpers | duplicate input, missing keys, inclusive range order, anchor removal | unit tests beside DataTable |
| React interaction | checkbox, Shift-click, select all, Space, Shift+Arrow, Escape, arrows preserve set | DataTable tests/stories |
| Accessibility | labels, indeterminate header, focus ring versus selected visuals, disabled action | interaction/a11y review |
| Widget adapter | selected keys/context and bulk action dispatch | `DataTable.widget`/WidgetRenderer tests |
| IR typing | JSON-safe prop shape, no full rows in selection context | TypeScript typecheck |
| Typed lowering | builder/spec validation and golden IR | `pkg/widgetdsl/spec/*_test.go` |
| V3 DSL | generated declarations/descriptors and JS builder result | `pkg/widgetdsl/module_test.go` |
| Repository | formatting/typecheck/tests | commands mandated by `AGENTS.md` |

Suggested commands after implementation:

```bash
pnpm biome check .
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site test
pnpm --dir packages/rag-evaluation-site build-storybook
go test ./pkg/widgetdsl/... ./pkg/xgoja/...
```

Use the actual package scripts from `package.json` before copying the exact test/build commands into CI documentation.

## Risks, alternatives, and open questions

### Risks

- **Semantic collision:** callers might pass both a selected detail key and a selected set. A discriminated API and typed validation must make this impossible.
- **Stale bulk selection:** sorting, filtering, and paging change visible rows. Normalize against current `rows` and do not claim cross-page selection.
- **Accessibility regression:** checkbox tables need careful focus/label/indeterminate behavior. Build the interaction model before DSL work.
- **Shortcut conflict:** row commands currently receive ordinary keys. Space and Escape must be reserved in multi mode and must not trigger a command by accident.
- **Payload trust:** keys from the browser are intent, not authorization. Bulk handlers must re-query and validate every key.

### Alternatives considered

1. **Shift-click only:** rejected as undiscoverable and weak for touch/keyboard users.
2. **Checkboxes only:** viable but slower for contiguous ranges; Shift range is low-complexity once anchor state exists.
3. **Use `selectedKey: string[]`:** rejected because it breaks TypeScript/API meaning and makes legacy single-selection ambiguous.
4. **Keep selection in a global WidgetRenderer store:** rejected because the design system is API-free and controlled props keep product ownership explicit.
5. **Start from `data.selection({ mode: "multi" })`:** deferred. It is generic and presently unconsumed; table-specific options make action and keyboard incompatibilities explicit.

### Open questions requiring an owner before merge

1. Should `onSelectionChange` be an action only, or should Widget pages be able to retain a client-local selected set without a server round trip? The first version should not invent global client state.
2. Does any existing consumer require selection to survive server pagination? If yes, define explicit server-side selection tokens rather than silently retaining hidden keys.
3. Should Shift-click range replace the selected set or union with it? This guide recommends union because it is least surprising after independently toggled rows, but UX review should confirm.
4. What bulk action visual form best matches the design system: an inline table toolbar or a surrounding `Panel` action region? Prototype both in Storybook before freezing the molecule API.

## Reference map

- `AGENTS.md` — repository formatting, package validation, and release constraints.
- `packages/rag-evaluation-site/GUIDELINES.md` — mandatory package layering, styling, stories, and React-first Widget IR rules.
- `packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.tsx:29-208` — current single selection/focus/keyboard/rendering behavior.
- `packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.module.css` — current tokenized table anatomy and visual states.
- `packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.stories.tsx` — package review entry point.
- `packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.widget.tsx:6-70` — Widget adapter and action context boundary.
- `packages/rag-evaluation-site/src/widgets/ir/props.ts:460-485` — serializable DataTable Widget IR props.
- `packages/rag-evaluation-site/src/widgets/actions.ts` — action dispatch and JSON payload/context transport.
- `packages/rag-evaluation-site/src/components/atoms/CheckboxRow/CheckboxRow.tsx` — existing checkbox atom reference, not necessarily the table implementation.
- `packages/rag-evaluation-site/src/widgets/WidgetRenderer.domain-registry.stories.tsx` — registry/story convention for Widget IR.
- `pkg/widgetdsl/spec/types.go:223-322` — collection, selection, table keyboard, and command typed model.
- `pkg/widgetdsl/spec/lower.go:235-276` and `pkg/widgetdsl/spec/lower_test.go` — collection-to-DataTable lowering and regression coverage.
- `pkg/widgetdsl/grammar.go:327-365` — older grammar path emitting singular selection.
- `pkg/widgetdsl/v3.go:1540-1569` and `pkg/widgetdsl/module_test.go` — existing generic multi-selection serialization.
