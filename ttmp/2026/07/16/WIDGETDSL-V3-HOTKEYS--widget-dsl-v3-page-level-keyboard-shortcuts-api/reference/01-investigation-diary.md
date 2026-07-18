---
Title: Investigation diary
Ticket: WIDGETDSL-V3-HOTKEYS
Status: active
Topics:
    - widget-dsl
    - ui-dsl
    - widget-ir
    - frontend
    - react
    - xgoja
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/code/wesen/claw-stuff/upwork/verbs/lib/pages.js
      Note: First consumer wiring for Triage Y/N/S actions
    - Path: repo://packages/rag-evaluation-site/README.md
      Note: Custom host integration and preference behavior documentation
    - Path: repo://packages/rag-evaluation-site/src/app/App.tsx
      Note: Primary React host evidence inspected during API design
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/KeyboardShortcutHelp/KeyboardShortcutHelp.tsx
      Note: Accessible native-dialog shortcut help implementation
    - Path: repo://packages/rag-evaluation-site/src/hooks/pageShortcuts.logic.ts
      Note: Pure canonicalization and matching implementation with focused checks
    - Path: repo://pkg/widgetdsl/v3.go
      Note: Primary Go runtime evidence inspected during API design
    - Path: repo://pkg/xgoja/providers/widgetsite/doc/02-widget-dsl-js-api-reference.md
      Note: Conceptual shortcut API and safety reference added after implementation
ExternalSources:
    - https://github.com/go-go-golems/rag-evaluation-system/issues/25
    - https://developer.mozilla.org/en-US/docs/Web/Accessibility/ARIA/Reference/Attributes/aria-keyshortcuts
    - https://www.w3.org/WAI/WCAG21/Understanding/character-key-shortcuts.html
Summary: 'Chronological investigation of the Widget DSL v3 page shortcut API design prompted by issue #25.'
LastUpdated: 2026-07-16T18:25:00-04:00
WhatFor: Use this diary to resume implementation with the original prompt, evidence, decisions, and validation plan intact.
WhenToUse: Read before implementing or revising WIDGETDSL-V3-HOTKEYS.
---




# Diary

## Goal

Capture the evidence and reasoning used to design a safe, serializable page-level keyboard shortcut
API for Widget DSL v3, with issue #25's Upwork Triage page as the first consumer.

## Step 1: Establish the issue, architecture, and API recommendation

I created a dedicated docmgr ticket and traced keyboard handling from the Go authoring model through
Widget IR lowering to the React application shell. The investigation showed that the existing
keyboard contract is intentionally owned by `DataTable`, while the application—not the leaf widget
renderer—owns page lifecycle and complete server action dispatch.

I proposed `page.shortcuts(...)` with stable command IDs, logical keys, labels, modifiers, and
existing `ActionSpec` values. The design includes event precedence, editable/IME/dialog guards,
duplicate validation, accessibility discoverability, and the single-character disable requirement
needed before broad release.

### Prompt Context

**User prompt (verbatim):** "Look at https://github.com/go-go-golems/rag-evaluation-system/issues/25 and create a new docmgr ticket. What kind of API would you design"

**Assistant interpretation:** Review issue #25 and the relevant implementation, create a structured
docmgr workspace, and propose a concrete upstream API for non-table page keyboard actions.

**Inferred user intent:** Unblock keyboard-first operation of the generated Upwork Triage page with
an API that fits Widget DSL v3 rather than adding a one-off consumer workaround.

### What I did

- Retrieved issue #25 with:
  `gh issue view 25 --json number,title,body,labels,author,url,comments`.
- Verified docmgr configuration and vocabulary with:
  `docmgr status --summary-only` and `docmgr vocab list`.
- Created ticket `WIDGETDSL-V3-HOTKEYS`, its design document, and this diary.
- Added six implementation/review tasks covering contract approval, Go/DSL work, React runtime,
  accessibility, tests, and consumer wiring.
- Read the repository's `AGENTS.md` instructions and the complete package design-system guidelines.
- Inspected `PageSpec`, `TableKeyboardSpec`, `RowCommandSpec`, lowering, the active v3 page builder,
  v3 descriptors, generated TypeScript declarations, `WidgetPageResponse`, `RagEvaluationSiteApp`,
  `WidgetRenderer`, central action dispatch, and the React `DataTable` handler.
- Captured line-numbered evidence with `nl -ba ... | awk ...` for the design document.
- Searched MDN and W3C sources for `aria-keyshortcuts`, character-key shortcut requirements, and
  `KeyboardEvent` behavior.
- Wrote `design-doc/01-page-level-keyboard-shortcuts-api-design.md` with API/IR sketches, five
  decision records, runtime pseudocode, phases, test matrix, risks, and open questions.

### Why

- A page card cannot correctly use table row keyboard semantics.
- The shortcut must remain serializable across Goja/xgoja, JSON transport, and React.
- Reusing `ActionSpec` preserves confirmation, server endpoint, toast, overlay, navigation, and
  refresh behavior.
- Page-level single-character commands introduce focus, assistive-technology, and accidental
  mutation risks that the API contract must address before implementation.

### What worked

- Issue #25 clearly identified both the immediate Triage use case and the current table-only gap.
- Existing action dispatch is already generic enough for shortcut activation; no action kind is
  needed.
- The page envelope and app shell provide a natural ownership boundary.
- Existing DataTable behavior provides a useful precedence model because consumed keys call
  `preventDefault()`.
- The existing topic vocabulary already contained every topic needed for the ticket.

### What didn't work

- I initially tried to read a monolithic IR file:
  `packages/rag-evaluation-site/src/widgets/ir.ts`.
- The exact tool error was:
  `ENOENT: no such file or directory, access '/home/manuel/code/wesen/go-go-golems/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/ir.ts'`.
- The IR is a directory with `index.ts`, `core.ts`, `actions.ts`, `props.ts`, and other modules. I
  corrected the path and inspected those files instead.

### What I learned

- There are currently two page representations: shared `widgetspec.PageSpec` and the active
  `v3PageSpec` in `v3.go`. A complete implementation must update both or first consolidate them.
- `WidgetRenderer` receives only a root node; it cannot own page shortcuts without broadening its
  contract. `RagEvaluationSiteApp` receives the page and owns the complete action pipeline.
- A button-only hotkey is too narrow because page commands may navigate or open overlays without a
  visual button.
- Stable command IDs and labels are needed even though key plus action is executable; they support
  diagnostics, help, analytics, accessibility, and future remapping.
- WCAG 2.1 criterion 2.1.4 means global `Y/N/S` cannot ship responsibly without a disable,
  remapping, or focus-only mechanism.

### What was tricky to build

- **Separating component and page ownership:** DataTable must continue to win when it consumes a
  key, but a card needs page accelerators. The proposed bubble-phase listener exits on
  `defaultPrevented` and on nested `data-rag-keyboard-scope` ownership.
- **Representing chords without a parser:** A compact `Ctrl+Y` string looks convenient but creates
  alias and platform ambiguity. The design uses `key` plus explicit modifiers.
- **Balancing the immediate request with accessibility:** `Y/N/S` is the desired operator UX, but
  unrestricted single-character shortcuts can interfere with speech input and assistive
  technology. The design treats help and a persisted disable control as release requirements,
  rather than silently allowing shortcuts in text/editing contexts.
- **Preserving one action path:** Shortcut dispatch must call the host's existing action callback,
  not `fetch` independently, so confirmation and refresh behavior remain identical to a click.

### What warrants a second pair of eyes

- Whether `page.shortcuts` is preferable to `page.keyboard` in the public vocabulary.
- Whether unmodified letter matching should ignore Shift, and how Caps Lock should behave.
- Whether the built-in disable/help UI is required in the first implementation or can be supplied
  by each host.
- How `shell.kind === "none"` and `root-owned` pages expose shortcut help without app chrome.
- Whether browser-reserved chord validation should be warning-only across platforms.
- Whether implementing both page models is acceptable or should trigger prior consolidation.

### What should be done in the future

- Review and accept the proposed API and IR envelope.
- Implement tasks `5wff`, `abw4`, `bd9v`, and `l7k3` in the ticket.
- Release the upstream package/generated host and then update the Upwork Triage consumer under task
  `9i0h`.
- Consider a later command-registry design that reuses shortcut IDs for menus and a command palette.

### Code review instructions

- Start with the API sketch and decision records in
  `design-doc/01-page-level-keyboard-shortcuts-api-design.md`.
- Verify the ownership argument against `pkg/widgetdsl/v3.go` (`v3PageBuilder`, `v3PageToIR`) and
  `packages/rag-evaluation-site/src/app/App.tsx` (`RagEvaluationSiteApp`, `handleAction`).
- Verify precedence against `DataTable.tsx` (`handleKeyDown`).
- During implementation run:
  - `make logcopter-check`
  - `go test ./...`
  - `pnpm biome check --write .`
  - `pnpm --dir packages/rag-evaluation-site typecheck`
  - `pnpm --dir packages/rag-evaluation-site test:focused`
  - generated host smoke tests used by PR #24.

### Technical details

Recommended authoring shape:

```js
page.shortcuts((keys) =>
  keys
    .bind("accept", "y", act.server("triage.accept"), { label: "Yes" })
    .bind("reject", "n", act.server("triage.reject"), { label: "No" })
    .bind("skip", "s", act.server("triage.skip"), { label: "Skip" }),
);
```

Recommended page transport shape:

```json
{
  "shortcuts": {
    "bindings": [
      {
        "id": "accept",
        "key": "y",
        "label": "Yes",
        "action": { "kind": "server", "name": "triage.accept" },
        "preventDefault": true,
        "allowRepeat": false
      }
    ]
  }
}
```

Recommended action context:

```json
{
  "componentType": "PageShortcut",
  "pageId": "triage",
  "shortcutId": "accept",
  "key": "y"
}
```

## Step 2: Implement the complete shortcut stack and Triage consumer

I implemented the approved API from typed Go authoring through browser runtime and the first real
consumer. Widget DSL pages now serialize validated shortcut bindings, generated declarations and
API help describe them, and the React host matches shortcuts through the existing action pipeline.

The browser runtime gives nested composite widgets precedence, blocks shortcuts while editing,
composing, or using a modal, suppresses key repeat by default, and exposes a persisted user toggle
plus generated help. The Upwork Triage page now binds `Y`, `N`, and `S` to the exact action objects
used by its visible buttons.

### Prompt Context

**User prompt (verbatim):** "alright, i like it. Build the entire ticket."

**Assistant interpretation:** Treat the proposed API as accepted and implement every ticket layer,
including contracts, runtime behavior, accessibility, tests, generated artifacts, documentation,
and the Upwork consumer.

**Inferred user intent:** Finish the keyboard workflow end to end so the Triage operator can use
single-key decisions without a one-off browser hack or divergence from Widget DSL v3.

**Commit (code):** 36bafb06a91a3374917c5b3e8e4dec53c2015ff7 — "widgetdsl: add page keyboard shortcuts"

**Commit (consumer):** 9905e09c1bf97d736eff7cddbdf61151675ff1e0 — "upwork: add triage keyboard shortcuts"

### What I did

- Added `PageShortcutSpec`, modifiers, page storage, canonical chords, validation, and lowering in
  `pkg/widgetdsl/spec`.
- Added `page.shortcuts` and `ShortcutsBuilder.bind` to the active v3 runtime.
- Preserved arbitrary action options while converting v3 actions so copy/event/overlay/navigation
  details survive shortcut lowering.
- Updated descriptors, action-context documentation, generated TypeScript declarations, fixture
  compilation, and embedded API help.
- Added v3 runtime tests, spec tests, a golden example, and generated golden IR.
- Added browser page types, pure chord matching, `usePageShortcuts`, page-level ARIA metadata, and
  nested keyboard-scope precedence.
- Added `KeyboardShortcutHelp` as a story-covered molecule using package atoms/foundation text and a
  native modal dialog.
- Added a browser-local “Enable page keyboard shortcuts” preference; disabling it removes both the
  listener and `aria-keyshortcuts` annotation.
- Marked keyboard-enabled `DataTable` instances as nested keyboard scopes.
- Rebuilt and synchronized `pkg/defaultspa/dist`.
- Updated `/home/manuel/code/wesen/claw-stuff/upwork/verbs/lib/pages.js` so Triage reuses three
  action objects for visible buttons and page shortcut bindings.

### Why

- The accepted design required one serializable behavior contract rather than button-local event
  handlers.
- Existing action dispatch must remain the only route to confirmations, server POSTs, toasts, and
  refresh.
- Single-character shortcuts require a user-facing disable mechanism and discoverable labels.
- Embedded SPA and generated reference artifacts must match source for CI and generated hosts.

### What worked

- Go spec, runtime, descriptor, golden, and full repository tests pass.
- Package TypeScript checking, focused shortcut checks, package build, app build, and Storybook build
  pass.
- The new Storybook story was included in the successful static build.
- The Upwork consumer passes `node --check` and `git diff --check`.
- `make logcopter-check` and generated-file checks pass.

### What didn't work

- The first Go formatting/test command failed because `pkg/widgetdsl/typescript.go` had a malformed
  edit:
  `pkg/widgetdsl/typescript.go:426:28: expected operand, found '{'` and
  `pkg/widgetdsl/typescript.go:452:2: expected ';', found 'return'`.
  I repaired the declaration slice and reran `gofmt`.
- The first shortcut runtime test imported the wrong require package and failed with:
  `no required module provides package github.com/go-go-golems/go-go-goja/modules/require`.
  I switched it to the repository's existing `github.com/dop251/goja_nodejs/require` package.
- Descriptor tests correctly caught stale generated help:
  `embedded API help descriptor reference is stale; regenerate ../xgoja/providers/widgetsite/doc/05-widget-dsl-v3-api-reference.md from WidgetV3APIReferenceMarkdown`.
  I regenerated the body and removed the generator's extra terminal blank line.
- Running `pnpm biome check --write .` touched 138 files and exposed unrelated baseline diagnostics
  (`Found 53 errors`, `Found 232 warnings`). I reverted all 88 unrelated tracked modifications,
  retained only ticket files, then ran Biome narrowly over the changed frontend files; that check
  passes.
- Playwright visual inspection could not start because its shared profile was occupied:
  `Browser is already in use for /home/manuel/.cache/ms-playwright/mcp-chrome-profile`.
  Storybook's production build still completed successfully, but no interactive screenshot was
  captured in this step.
- The first commit attempt passed formatting, frontend typecheck, and Go tests but failed when the
  hook encountered another concurrent lint process: `Error: parallel golangci-lint is running`.
  No lint process remained when checked, so I retried normally; all pre-commit hooks then passed and
  commit `36bafb06a91a3374917c5b3e8e4dec53c2015ff7` was created.

### What I learned

- The v3 action converter must preserve unknown action options, not only navigation options, because
  shortcuts can dispatch every existing action kind.
- Bubble-phase page handling plus `defaultPrevented` is necessary but not sufficient for composite
  widgets; an explicit `data-rag-keyboard-scope` prevents unrecognized page commands from stealing
  future component keys.
- A native `<dialog>` gives the help surface focus containment and reliable modal suppression with
  less custom accessibility machinery.
- `aria-keyshortcuts` must disappear when shortcuts are disabled; discoverability metadata should
  reflect active behavior.
- Repository-wide Biome currently includes unrelated baseline failures, so ticket validation needs
  both the mandated broad attempt and a clean targeted check.

### What was tricky to build

- **Validation parity across two page models:** shared `PageSpec` and active `v3PageSpec` both had to
  carry shortcuts. The v3 validator reuses shared shortcut diagnostics and rejects only errors, so
  the accessibility warning does not prevent `toPage()`.
- **Exact chord matching:** the implementation canonicalizes modifier ordering and ASCII letter
  case while still requiring exact modifier state. This keeps `Control+s` distinct from `Meta+s`
  and prevents shifted keystrokes from silently matching unmodified bindings.
- **Modal and focus safety:** checking only the event target was insufficient because focus can
  remain on the opener briefly. The hook also detects any active modal before matching.
- **Preference behavior:** the user toggle needed to suppress execution and remove ARIA metadata,
  while tolerating unavailable `localStorage` through an in-memory fallback.
- **Generated artifacts:** descriptor help, golden JSON, and embedded SPA assets all needed explicit
  regeneration after source changes.
- **External consumer state:** the Upwork repository already contains many unrelated modifications.
  I edited only `upwork/verbs/lib/pages.js` and did not stage, reset, or alter its other work.

### What warrants a second pair of eyes

- Browser-level behavior of Shift/Caps Lock across non-US keyboard layouts.
- Whether every future composite keyboard widget consistently adds `data-rag-keyboard-scope`.
- Placement of the fixed shortcut-help trigger in highly customized root-owned applications.
- The custom host's eventual upgrade/rebuild path after the upstream package is merged and released.
- Interactive focus/visual behavior of the help dialog, because Playwright was unavailable during
  this step.

### What should be done in the future

- Merge/release the upstream branch and rebuild the Upwork generated binary against that release.
- Consider per-command remapping using the stable shortcut IDs if disabling alone is insufficient.
- Consider consolidating `PageSpec` and `v3PageSpec` in a separate refactor.

### Code review instructions

- Start at `pkg/widgetdsl/spec/types.go` (`PageShortcutSpec`) and
  `pkg/widgetdsl/v3.go` (`v3ShortcutsBuilder`, `v3PageToIR`).
- Review browser matching in
  `packages/rag-evaluation-site/src/hooks/pageShortcuts.logic.ts` and listener ownership in
  `usePageShortcuts.ts`.
- Review host integration and preference behavior in
  `packages/rag-evaluation-site/src/app/App.tsx`.
- Review help UI and Storybook coverage under
  `components/molecules/KeyboardShortcutHelp/`.
- Review the consumer diff in
  `/home/manuel/code/wesen/claw-stuff/upwork/verbs/lib/pages.js` around `triagePage`.
- Validate with:
  - `make logcopter-check`
  - `go test ./...`
  - `pnpm --dir packages/rag-evaluation-site typecheck`
  - `pnpm --dir packages/rag-evaluation-site test:focused`
  - `pnpm --dir packages/rag-evaluation-site build`
  - `pnpm --dir packages/rag-evaluation-site build-storybook`
  - `node --check /home/manuel/code/wesen/claw-stuff/upwork/verbs/lib/pages.js`

### Technical details

Runtime dispatch context:

```json
{
  "componentType": "PageShortcut",
  "pageId": "triage",
  "shortcutId": "accept",
  "key": "y"
}
```

Ignored event conditions are `defaultPrevented`, IME composition, disabled preference, editable or
modal target, any active modal, nested keyboard scope, unmatched exact modifiers, and disallowed
repeat. Matching bindings call the same `handleAction` callback used by rendered widgets.

## Step 3: Extend public and consumer documentation

I updated the documentation surfaces that readers use outside the ticket workspace. The conceptual
API reference now explains the shortcut contract and safety model, while the runnable v3 tutorial
shows the same action-reuse pattern implemented by the Upwork Triage page.

I also documented the React host responsibilities and browser preference in the npm package README,
and added concrete Triage controls to the Upwork operator README. These changes close the gap
between generated method inventory, authoring guidance, host integration, and end-user operation.

### Prompt Context

**User prompt (verbatim):** "do it"

**Assistant interpretation:** Apply the previously recommended documentation updates to the v3
examples, conceptual API reference, package README, and Upwork README.

**Inferred user intent:** Ensure the new feature is discoverable and usable without reading the
implementation ticket or source diffs.

### What I did

- Read the Glazed help authoring guidance with:
  `glaze help how-to-write-good-documentation-pages` and
  `glaze help writing-help-entries`.
- Expanded `02-widget-dsl-js-api-reference.md` into a complete conceptual reference with valid
  `SectionType`, page shortcut API, options table, action context, safety semantics,
  troubleshooting, and See Also links.
- Added a runnable `Y/N/S` example, modifier example, host behavior, failure modes, and golden
  fixture reference to `04-widget-dsl-v3-examples.md`.
- Added page-envelope, runtime ownership, accessibility preference, custom-host integration, and
  Storybook guidance to `packages/rag-evaluation-site/README.md`.
- Added separate Job table and Triage keyboard sections to
  `/home/manuel/code/wesen/claw-stuff/upwork/README.md`.

### Why

- The generated API reference lists methods but does not teach safe composition or explain runtime
  behavior.
- Custom React hosts need to know that rendering `WidgetRenderer` alone does not install page
  shortcuts.
- Operators need the actual Triage keys and disable-control location in the application README.
- Glazed help entries need troubleshooting and cross-references to remain discoverable from CLI
  help.

### What worked

- Existing help frontmatter and related pages provided stable slugs for cross-references.
- The executable golden example supplied a source-backed code sample for the tutorial.
- Package and consumer documentation could describe the implemented behavior without introducing
  new contracts.

### What didn't work

- I mistakenly ran docmgr frontmatter validation on the embedded Glazed help pages. It failed with
  `schema validation failed: missing required fields: Ticket, DocType` because Glazed help
  frontmatter intentionally uses `Slug`, `Short`, and `SectionType` instead of docmgr's ticket
  schema. I kept the Glazed schema, validated loading through the provider/widgetdsl Go tests, and
  used `docmgr doctor` only for the ticket workspace.

### What I learned

- The conceptual API page was previously too short to explain even the existing action/binding
  boundary; documenting shortcuts was a useful opportunity to make that distinction explicit.
- Host documentation must mention both execution and discoverability because `aria-keyshortcuts`
  and a listener alone do not satisfy the operator-facing requirement.
- The Upwork application has two keyboard vocabularies whose overlapping `S` key is safe only
  because table scopes take precedence over page commands.

### What was tricky to build

- The four documents target different readers: DSL authors, tutorial readers, React embedders, and
  Upwork operators. Repeating the same long reference would make each less useful, so each update
  emphasizes its reader's boundary and links to the others.
- The Glazed help pages must omit a manually rendered top-level title and include operational
  troubleshooting, while ordinary package and consumer READMEs retain conventional Markdown
  headings.

### What warrants a second pair of eyes

- Whether the browser preference key should be treated as stable public API or documented as an
  implementation detail.
- Whether the Upwork README should include screenshots after the upstream host is rebuilt.

### What should be done in the future

- Add release notes and bump the package version when publication is scheduled.
- Rebuild the Upwork binary after the upstream package release and verify the documented keys in a
  live browser.

### Code review instructions

- Review the conceptual contract in `02-widget-dsl-js-api-reference.md` first.
- Compare the tutorial snippet in `04-widget-dsl-v3-examples.md` with
  `pkg/widgetdsl/testdata/v3/examples/43-page-shortcuts.js`.
- Verify package host claims against `usePageShortcuts.ts` and `App.tsx`.
- Verify Upwork key descriptions against `triagePage` in `upwork/verbs/lib/pages.js`.
- Validate embedded help with `go test ./pkg/xgoja/providers/widgetsite/... ./pkg/widgetdsl/...`.

### Technical details

The documentation now distinguishes:

1. `page.shortcuts(...)` for page commands;
2. `table.command(...)` for focused row commands;
3. `usePageShortcuts` plus help/preference UI for custom React hosts;
4. visible buttons as the primary interaction surface.
