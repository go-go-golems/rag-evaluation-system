---
Title: Page-level keyboard shortcuts API design
Ticket: WIDGETDSL-V3-HOTKEYS
Status: active
Topics:
    - widget-dsl
    - ui-dsl
    - widget-ir
    - frontend
    - react
    - xgoja
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/src/app/App.tsx
      Note: Page lifecycle and complete local/server action dispatch ownership
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.tsx
      Note: Existing component-scoped keyboard handling and precedence model
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/KeyboardShortcutHelp/KeyboardShortcutHelp.tsx
      Note: Implemented discoverability and persisted-enable UI surface
    - Path: repo://packages/rag-evaluation-site/src/hooks/usePageShortcuts.ts
      Note: Implemented page-owned listener lifecycle and event safety guards
    - Path: repo://packages/rag-evaluation-site/src/hooks/useWidgetPage.ts
      Note: Browser WidgetPageResponse transport contract
    - Path: repo://pkg/widgetdsl/spec/types.go
      Note: Shared typed page, table keyboard, row-command, and action contracts that establish the new shortcut model boundary
    - Path: repo://pkg/widgetdsl/v3.go
      Note: Active v3 page builder, page lowering, validation, and current table keyboard builders
ExternalSources:
    - https://github.com/go-go-golems/rag-evaluation-system/issues/25
    - https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Reference/Attributes/aria-keyshortcuts
    - https://www.w3.org/WAI/WCAG21/Understanding/character-key-shortcuts.html
    - https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent
Summary: Design for serializable, page-owned widget.dsl shortcuts that dispatch existing actions safely and accessibly.
LastUpdated: 2026-07-16T18:25:00-04:00
WhatFor: Use this document to review and implement page-level keyboard shortcuts in Widget DSL v3 and the React host.
WhenToUse: Use before changing the page builder, Widget IR page envelope, React app shell, or generated xgoja hosts for keyboard shortcuts.
---



# Page-level keyboard shortcuts API design

## Executive summary

Issue #25 exposes a real abstraction gap: Widget DSL v3 can describe keyboard behavior for a
`DataTable`, but cannot bind a page-level key to an ordinary serializable action. The Upwork Triage
page is a card with Yes, No, and Skip actions, so table row navigation is the wrong primitive.

I recommend a **page-owned command API** named `page.shortcuts(...)`, not a button-only `hotkey`
property and not a generic DOM keyboard-event API. Each binding has a stable ID, a logical
`KeyboardEvent.key`, optional modifiers, a human-readable label, and an existing `ActionSpec`. The
page envelope transports those bindings; the React application shell owns one `keydown` listener
and dispatches matches through the same action handler used by buttons and navigation.

The proposed authoring API is:

```js
const page = widget.page("Triage", (page) =>
  page
    .shortcuts((keys) =>
      keys
        .bind("accept", "y", act.server("triage.accept"), { label: "Yes" })
        .bind("reject", "n", act.server("triage.reject"), { label: "No" })
        .bind("skip", "s", act.server("triage.skip"), { label: "Skip" }),
    )
    .view(triageCard),
);
```

This design deliberately keeps shortcuts out of individual button props. Shortcuts represent page
commands: a command may navigate, call the server, open an overlay, or have no visual button. The
visible Triage buttons remain the primary operable controls; keyboard bindings are accelerators.

## Implementation outcome

The design was accepted and implemented on 2026-07-16. The shipped implementation includes typed
Go specs and validation, `page.shortcuts` and `ShortcutsBuilder.bind`, generated declarations and
reference help, page-owned React matching, nested keyboard-scope precedence, an accessible help
molecule, a persisted disable preference, focused and golden tests, refreshed embedded SPA assets,
and `Y/N/S` wiring in the Upwork Triage consumer.

The implementation follows the proposed transport and action context. It also generalized the v3
action converter to preserve arbitrary action options, ensuring shortcut activation is behaviorally
identical to visible controls for every existing action kind.

## Problem statement and scope

### Required outcome

A generated xgoja page must be able to declare `Y`, `N`, and `S` shortcuts that invoke the same
serializable actions as the visible Triage controls. The implementation must:

1. work for pages whose root is not a table;
2. preserve the existing action pipeline, including confirmations and server refresh behavior;
3. avoid firing while the user types, uses an IME, or operates a modal;
4. define deterministic precedence with component-owned keyboard behavior;
5. expose labels and canonical key data for help and accessibility;
6. reject ambiguous duplicate bindings before transport;
7. remain JSON-compatible and describable in generated TypeScript declarations.

### In scope

- Widget DSL v3 page builder surface;
- typed Go shortcut representation and page lowering;
- page response TypeScript contract;
- React page-level event ownership and action dispatch;
- validation, descriptors, declarations, tests, and discoverability requirements;
- a migration example for the Upwork Triage consumer.

### Out of scope for the first implementation

- arbitrary JavaScript callback handlers in browser IR;
- key sequences such as `g g`;
- context-sensitive predicates serialized as executable code;
- replacing DataTable's row focus and command behavior;
- OS-wide shortcuts;
- a general command palette, although the contract should support one later.

## Current-state architecture

### The page is the correct transport boundary

The shared typed authoring model already treats `PageSpec` as the envelope for ID, title, metadata,
shell, root, and diagnostics (`pkg/widgetdsl/spec/types.go:11-20`). Its lowerer writes those fields at
the top level of the browser payload (`pkg/widgetdsl/spec/lower.go:7-25`). The active v3 builder uses
a parallel `v3PageSpec`; `v3PageToIR` emits the same top-level page envelope
(`pkg/widgetdsl/v3.go:1909-1939`).

A shortcut that should work regardless of root component therefore belongs beside `shell` and
`root`, not inside a `Button` or synthetic widget node.

### Existing keyboard behavior is table-local

`TableSpec` owns `TableKeyboardSpec` and row commands (`pkg/widgetdsl/spec/types.go:266-296`). Its
lowerer places `keyboard` and `commands` in `DataTable` props
(`pkg/widgetdsl/spec/lower.go:222-230`). The React `DataTable` listens on the focused table row,
handles ArrowUp/ArrowDown, optional `j`/`k`, Enter, and row commands, and calls `preventDefault()`
for a match (`packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.tsx:96-120`).

That is a sound component-scoped design and should remain unchanged. Reusing it for a card would
require inventing fake rows or moving table assumptions into the page.

### The application shell owns page lifecycle and server action dispatch

`WidgetPageResponse` currently contains page ID, title, shell, root, and metadata
(`packages/rag-evaluation-site/src/hooks/useWidgetPage.ts:39-45`). `RagEvaluationSiteApp` fetches
that page and owns `handleAction`, including confirmation, non-server dispatch, server POST,
result events, toasts, and refresh (`packages/rag-evaluation-site/src/app/App.tsx:67-110`).

`WidgetRenderer` only receives one node, a registry, and an action callback
(`packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx:14-27`). It cannot see page-level
metadata. Therefore page shortcuts should be installed by `RagEvaluationSiteApp` (or a page-level
hook it calls), not by `WidgetRenderer` and not by every component adapter.

### The action contract is already reusable

Every shortcut can carry the existing `ActionSpec` union. The application can call its existing
`onAction` handler with a shortcut context, and `dispatchWidgetAction` already delegates to that
handler before handling local actions itself (`packages/rag-evaluation-site/src/widgets/actions.ts:44-55`).
No new action kind is required.

## Gap analysis

| Capability | Current state | Required change |
|---|---|---|
| Table row navigation | Supported by `TableKeyboardSpec` | Keep unchanged |
| Table row commands | Supported by `RowCommandSpec` | Keep component-local |
| Card/page commands | Not representable | Add page shortcut bindings |
| Serializable action dispatch | Already supported | Reuse `ActionSpec` |
| Page runtime listener | None | Add one page-owned hook/controller |
| Editable/IME/repeat guards | Table has an editable-target guard only | Define shared page guards |
| Duplicate key validation | None at page level | Normalize chords and reject duplicates |
| Shortcut discoverability | Table caption lists row commands | Expose page labels/chords and a help surface |
| Single-character disable/remap | None | Add host-level disable mechanism before general release |

## Proposed public API

### Authoring surface

Add `shortcuts(configure)` to `PageBuilder` and a dedicated `ShortcutBuilder`:

```ts
export type ShortcutModifier = "Alt" | "Control" | "Meta" | "Shift";

export interface ShortcutOptions {
  label: string;
  modifiers?: ShortcutModifier[];
  preventDefault?: boolean;
  allowRepeat?: boolean;
}

export interface ShortcutsBuilder extends ComposableBuilder<ShortcutsBuilder> {
  bind(
    id: string,
    key: string,
    action: ActionSpec,
    options: ShortcutOptions,
  ): this;
}

export interface PageBuilder extends ComposableBuilder<PageBuilder> {
  // existing methods...
  shortcuts(configure: Fragment<ShortcutsBuilder>): this;
}
```

The four positional concepts are intentionally distinct:

- `id` is stable command identity for diagnostics, analytics, testing, and a future command palette;
- `key` is the logical browser key (`"y"`, `"Enter"`, `"Escape"`), not a parsed display string;
- `action` is existing serializable behavior;
- `options.label` is required human-readable intent, while modifiers and runtime behavior are options.

For modifier shortcuts:

```js
page.shortcuts((keys) =>
  keys.bind("save", "s", act.server("record.save"), {
    label: "Save",
    modifiers: ["Control"],
  }),
);
```

The first version should not accept `allowInEditable`. Page shortcuts must never fire from
`input`, `textarea`, `select`, contenteditable regions, or active dialogs. Making that safety rule
configurable invites destructive action dispatch while typing.

### Typed Go model

Add a page-level value object:

```go
type ShortcutModifier string

const (
    ShortcutModifierAlt     ShortcutModifier = "Alt"
    ShortcutModifierControl ShortcutModifier = "Control"
    ShortcutModifierMeta    ShortcutModifier = "Meta"
    ShortcutModifierShift   ShortcutModifier = "Shift"
)

type PageShortcutSpec struct {
    ID             string
    Key            string
    Modifiers      []ShortcutModifier
    Label          string
    Action         ActionSpec
    PreventDefault bool
    AllowRepeat    bool
}

type PageSpec struct {
    // existing fields
    Shortcuts []PageShortcutSpec
}
```

The active `v3PageSpec` should also carry `[]widgetspec.PageShortcutSpec` until the parallel page
models are consolidated. This avoids introducing another untyped `map[string]any` island.

Defaults applied by the builder:

- `PreventDefault: true`;
- `AllowRepeat: false`;
- no modifiers;
- case-insensitive matching for a single printable key;
- exact `KeyboardEvent.key` matching for named keys.

### Widget IR page envelope

Use an object envelope rather than a bare array so policy and help metadata can evolve without
changing the top-level field's shape:

```json
{
  "schemaVersion": "0.1.0",
  "id": "triage",
  "title": "Triage",
  "shortcuts": {
    "bindings": [
      {
        "id": "accept",
        "key": "y",
        "modifiers": [],
        "label": "Yes",
        "action": { "kind": "server", "name": "triage.accept" },
        "preventDefault": true,
        "allowRepeat": false
      }
    ]
  },
  "root": { "kind": "component", "type": "Stack" }
}
```

Corresponding browser types:

```ts
export interface PageShortcutSpec {
  id: string;
  key: string;
  modifiers?: ShortcutModifier[];
  label: string;
  action: ActionSpec;
  preventDefault?: boolean;
  allowRepeat?: boolean;
}

export interface PageShortcutsSpec {
  bindings: PageShortcutSpec[];
}

export interface WidgetPageResponse {
  // existing fields
  shortcuts?: PageShortcutsSpec;
}
```

### Action context

A shortcut dispatches this stable context:

```json
{
  "componentType": "PageShortcut",
  "pageId": "triage",
  "shortcutId": "accept",
  "key": "y"
}
```

Register it in the v3 action-context descriptors as `page.shortcut`. Payload bindings can then
refer to these fields without inventing a separate handler protocol.

## Matching and runtime semantics

### Canonical chord

Validation and runtime matching must use the same canonical representation:

```text
canonical = sorted(unique(modifiers)) + normalize(key)
normalize(single printable ASCII letter) = lowercase letter
normalize(named key) = exact KeyboardEvent.key spelling
```

Examples:

- `y` and `Y` conflict;
- `Control+s` and `Control+S` conflict;
- `Control+s` and `Meta+s` do not conflict;
- duplicate modifiers are invalid rather than silently retained.

Do not parse strings such as `"Ctrl+S"`. Parsing introduces aliases (`Cmd`, `Command`, `⌘`),
platform ambiguity, escaping problems for `+`, and inconsistent generated declarations.

### Event ownership and precedence

Install one `keydown` listener in the page runtime during the bubble phase. A component-local
handler receives the event first. The page handler exits when `event.defaultPrevented` is true, so
DataTable and future component scopes retain precedence.

The page handler also exits when:

1. `event.isComposing` is true;
2. `event.repeat` is true and the binding does not allow repeat;
3. the target is an input, textarea, select, contenteditable element, or inside an active dialog;
4. a closer keyboard scope marks the event as owned;
5. browser shortcut modifiers do not exactly equal the binding modifiers;
6. shortcuts are disabled by the user's host preference.

A match calls `preventDefault()` when configured, then dispatches the existing action once.

```ts
function onPageKeyDown(event: KeyboardEvent) {
  if (event.defaultPrevented || event.isComposing) return;
  if (isEditableTarget(event.target) || isInsideActiveDialog(event.target)) return;
  if (isOwnedByNestedKeyboardScope(event.target)) return;

  const binding = matchShortcut(event, page.shortcuts?.bindings ?? []);
  if (!binding || (event.repeat && !binding.allowRepeat)) return;

  if (binding.preventDefault !== false) event.preventDefault();
  onAction(binding.action, {
    componentType: "PageShortcut",
    pageId: page.id,
    shortcutId: binding.id,
    key: event.key,
  });
}
```

The listener must be removed when the page changes or unmounts. Its effect dependencies should be
page ID, shortcut bindings, and the stable action callback, preventing stale actions after route
transitions.

### Nested component behavior

`DataTable` already prevents default for keys it consumes. As a hardening step, keyboard-owning
components should also expose `data-rag-keyboard-scope="DataTable"`; the page guard should ignore
events originating inside that scope. This prevents an unrecognized page command from stealing a
key that a focused composite widget may reserve later.

The Triage card has no nested keyboard scope, so `Y`, `N`, and `S` reach the page handler.

## Accessibility and safety requirements

`aria-keyshortcuts` communicates an implemented accelerator but does not implement it. MDN also
recommends making shortcuts discoverable through menus, tooltips, a cheat sheet, or equivalent.
The page root should expose a canonical `aria-keyshortcuts` value, and the host should provide a
visible shortcut-help affordance generated from the required labels.

WCAG 2.1 success criterion 2.1.4 requires single-character shortcuts to be turn-off-able,
remappable, or active only while the relevant component has focus. Because Triage requests global
single characters, general release must include a user-facing disable mechanism. A practical first
host implementation is:

1. a “Keyboard shortcuts” help control in app chrome;
2. a help dialog listing label and chord;
3. an “Enable single-key shortcuts” toggle persisted per browser;
4. page shortcuts disabled while a dialog is open;
5. visible buttons remain fully operable without shortcuts.

Remapping can be deferred, because a disable mechanism satisfies the criterion. The IR contract's
stable IDs make future per-command remapping possible without changing action identity.

Destructive shortcuts must continue through `confirmWidgetAction`; shortcut dispatch must not
bypass the current action handler. `allowRepeat` should default false to prevent repeated server
mutations when a key is held down.

## Validation rules

Page validation should emit errors for:

1. blank IDs, keys, or labels;
2. duplicate shortcut IDs;
3. duplicate canonical chords;
4. unsupported or duplicate modifier names;
5. missing/invalid actions;
6. key strings containing multiple keys or a serialized chord such as `Ctrl+Y`;
7. bindings that use modifier-only keys.

It should emit warnings for:

1. unmodified single-character shortcuts, reminding authors about disable/focus requirements;
2. browser/OS-reserved or risky chords where known;
3. `allowRepeat` on server or destructive actions.

Validation paths should be stable, for example `page.shortcuts.bindings[1].key`, so generated-host
errors point to a specific declaration.

## Decision records

### Decision: use page commands rather than button hotkeys

- **Context:** Triage has buttons, but future shortcuts may navigate, open overlays, or invoke
  commands without a one-to-one button.
- **Options considered:** `ui.button(..., { hotkey: "y" })`, `page.shortcuts(...)`, or both.
- **Decision:** Make `page.shortcuts(...)` the canonical behavior API.
- **Rationale:** The page is the lifecycle and transport boundary, while actions are already
  independent of visual controls. A button-only API couples behavior to one renderer.
- **Consequences:** Buttons may need a separate presentation convention for showing key hints, but
  command identity and dispatch stay centralized.
- **Status:** accepted.

### Decision: transport structured key plus modifiers

- **Context:** Shortcut strings have aliases and platform-dependent notation.
- **Options considered:** parse `"Ctrl+Y"`, store a DOM `code`, or store `key` plus modifiers.
- **Decision:** Store logical `KeyboardEvent.key` and an explicit modifier list.
- **Rationale:** It is serializable, localizable by browser keyboard layout, easy to validate, and
  maps directly to DOM events.
- **Consequences:** Physical-position shortcuts are not supported initially; display formatting is
  a renderer concern.
- **Status:** accepted.

### Decision: React app owns one bubble-phase listener

- **Context:** `WidgetRenderer` receives a node, while `RagEvaluationSiteApp` receives the page and
  owns complete action dispatch.
- **Options considered:** listener per button, listener in `WidgetRenderer`, or listener in the app
  page controller.
- **Decision:** Install one page listener through an app-level hook/controller.
- **Rationale:** It aligns lifecycle, payload visibility, and server action ownership, and allows
  component handlers to win through `defaultPrevented`.
- **Consequences:** Embedders that render pages without `RagEvaluationSiteApp` need an exported
  reusable hook/controller, not hidden app-only logic.
- **Status:** accepted.

### Decision: do not permit page shortcuts while editing or in overlays

- **Context:** Single-letter actions can mutate server state and conflict with text entry or modal
  controls.
- **Options considered:** per-binding `allowInEditable`, host-wide suppression, or unconditional
  suppression.
- **Decision:** Unconditionally suppress in editable targets and active dialogs for v1.
- **Rationale:** Safety outweighs niche editor shortcuts; editor-specific commands belong to the
  editor component's own keyboard scope.
- **Consequences:** A future editor command API must be component-scoped rather than weakening page
  invariants.
- **Status:** accepted.

### Decision: stable IDs and labels are mandatory

- **Context:** A key and action alone are enough to execute but insufficient for diagnostics,
  accessibility, help UI, analytics, and remapping.
- **Options considered:** infer IDs/labels, make them optional, or require them.
- **Decision:** Require both a stable ID and a human-readable label.
- **Rationale:** This keeps the serialized contract self-describing and supports future host UX.
- **Consequences:** Authoring is slightly more verbose but avoids retrofitting identity later.
- **Status:** accepted.

## Alternatives considered

### Button `hotkey` option

```js
ui.button("Yes", action, { hotkey: "y" })
```

This is concise for Triage but does not define page-level conflict resolution, cannot represent
non-button commands, and would install behavior from leaf renderers. It is acceptable later as a
presentation hint that references a page binding ID, but not as the behavior source of truth.

### Raw `page.onKeyDown(callback)`

This cannot cross JSON transport, bypasses typed actions and xgoja host boundaries, and would move
browser execution into an unsafe callback model. Reject it.

### Reuse `DataTable` commands

Table commands require focused row context and dispatch a row with the action. Triage has no row.
Generalizing them would weaken the table contract and still not solve page ownership. Reject it.

### A generic global command registry before page shortcuts

A full command registry could unify menus, palettes, keyboard shortcuts, and analytics. It is a
reasonable future direction, but larger than issue #25. Stable shortcut IDs, labels, and actions
allow this proposal to evolve into such a registry without blocking the immediate Triage use case.

## Implementation plan

### Phase 1: typed contract and lowering

1. Add `PageShortcutSpec` and modifier constants in `pkg/widgetdsl/spec/types.go`.
2. Add shortcut validation in `pkg/widgetdsl/spec/validate.go` and unit tests.
3. Lower the page envelope in `pkg/widgetdsl/spec/lower.go`.
4. Add `Shortcuts []widgetspec.PageShortcutSpec` to `v3PageSpec` and emit the envelope in
   `v3PageToIR`.
5. Add `page.shortcuts` and `ShortcutsBuilder.bind` in `pkg/widgetdsl/v3.go`.

### Phase 2: declaration and parity surfaces

1. Add `ShortcutsBuilder` and shortcut types to `pkg/widgetdsl/typescript.go`.
2. Add `shortcuts` and `bind` to `pkg/widgetdsl/v3_descriptors.go`.
3. Add `page.shortcut` action-context fields.
4. Extend descriptor parity, TypeScript fixture, generated help, and migration-checker expectations.
5. Add a golden v3 example using Yes/No/Skip.

### Phase 3: React runtime

1. Extend `WidgetPageResponse` in `useWidgetPage.ts`.
2. Add pure canonicalization/matching helpers in a page-shortcuts logic module.
3. Add an exported `usePageShortcuts` hook or `PageShortcutController` that accepts bindings,
   page ID, enabled state, and an action handler.
4. Call it from `RagEvaluationSiteApp`, preserving the current `handleAction` path.
5. Mark composite keyboard owners such as `DataTable` with `data-rag-keyboard-scope`.
6. Add `aria-keyshortcuts` to the page root and a package-owned help/toggle surface following the
   design-system layer rules.

### Phase 4: tests and consumer migration

1. Go unit tests: lowering, defaults, action preservation, validation, duplicate IDs/chords.
2. Runtime builder tests: exact IR output and fluent callback composition.
3. Descriptor and TypeScript declaration parity tests.
4. Pure TypeScript focused checks: case normalization, modifiers, repeats, composition, editable
   targets, nested scopes, dialogs, and `defaultPrevented`.
5. Storybook/integration story: buttons plus `Y/N/S`, help text, disabled shortcut preference.
6. Generated-host test: shortcut server action reaches the existing endpoint and refresh policy.
7. Update the Upwork Triage page only after the upstream package and generated host include the new
   contract.

## Test matrix

| Scenario | Expected result |
|---|---|
| Press `y` on Triage page | Accept action dispatches once |
| Press uppercase `Y` without a Shift modifier (for example, Caps Lock) | Same single-letter binding matches |
| Press `y` in a text input | No page action |
| Press `y` during IME composition | No page action |
| Hold `y` | One action unless `allowRepeat` is true |
| DataTable consumes a command | Page shortcut does not run |
| Focus is in an open dialog | Page shortcut does not run |
| Two bindings normalize to `y` | Page validation error |
| `Control+s` and `Meta+s` | Distinct bindings |
| Server shortcut has confirmation | Existing confirmation appears |
| User disables single-key shortcuts | No unmodified character shortcut runs |
| Route changes | Old page bindings are removed |

## Risks and mitigations

- **Accidental mutations:** suppress repeats/editable targets/dialogs and keep confirmation in the
  central action path.
- **Shortcut conflicts:** validate canonical chords and give nested component scopes precedence.
- **Accessibility regression:** require labels, expose help, annotate root, and provide disable UI.
- **Stale action closure after navigation:** install and clean up through an effect keyed by page.
- **Cross-layout key behavior:** use `event.key` intentionally; document that physical key positions
  are not the contract.
- **Parallel page models:** implement both shared `PageSpec` and active `v3PageSpec`, then consider a
  separate consolidation ticket rather than hiding the divergence.
- **Embedding outside the stock app:** export the reusable page-shortcut controller/hook from the
  package so alternate hosts can preserve the same semantics.

## Open questions

1. Should the first release include remapping, or is a persisted disable toggle sufficient? This
   design recommends disable first.
2. Where should shortcut help live for `shell.kind === "none"` and `root-owned` pages? A small
   package-owned overlay trigger may be needed because app navigation chrome is absent.
3. Should `Shift+y` match an unmodified `y` binding because uppercase `Y` implies Shift on many
   layouts? This design recommends normalizing letter case while requiring declared non-Shift
   modifiers exactly; tests must lock the behavior.
4. Should known reserved browser shortcuts be errors or warnings? This design recommends warnings
   because reservations vary by browser and operating system.
5. Should a later button option reference a shortcut ID solely to render a key hint and
   `aria-keyshortcuts` on the control? That can be added without changing command execution.

## References

### Repository evidence

- `pkg/widgetdsl/spec/types.go:11-20` — shared page envelope.
- `pkg/widgetdsl/spec/types.go:266-296` — table-specific keyboard and row command types.
- `pkg/widgetdsl/spec/lower.go:7-25` — shared page lowering.
- `pkg/widgetdsl/spec/lower.go:222-230` — table keyboard lowering.
- `pkg/widgetdsl/v3.go:1626-1699` — active page builder.
- `pkg/widgetdsl/v3.go:1909-1939` — active v3 page transport.
- `pkg/widgetdsl/v3.go:2144-2173` — active v3 page validation.
- `packages/rag-evaluation-site/src/hooks/useWidgetPage.ts:39-45` — browser page response.
- `packages/rag-evaluation-site/src/app/App.tsx:67-110` — host action dispatch and refresh.
- `packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.tsx:96-120` — existing
  component-local keyboard precedence.
- `packages/rag-evaluation-site/src/widgets/actions.ts:44-55` — reusable action dispatch entrypoint.

### External references

- GitHub issue #25: <https://github.com/go-go-golems/rag-evaluation-system/issues/25>
- MDN `aria-keyshortcuts`:
  <https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Reference/Attributes/aria-keyshortcuts>
- W3C WCAG 2.1, Understanding 2.1.4 Character Key Shortcuts:
  <https://www.w3.org/WAI/WCAG21/Understanding/character-key-shortcuts.html>
- MDN `KeyboardEvent`: <https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent>
