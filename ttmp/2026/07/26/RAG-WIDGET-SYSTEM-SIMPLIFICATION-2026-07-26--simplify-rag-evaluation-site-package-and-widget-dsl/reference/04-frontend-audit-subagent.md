---
Title: Frontend package specialist audit
Ticket: RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26
Status: active
Topics:
    - architecture
    - frontend
    - code-quality
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Independent review of React package architecture, adapters, registries, manifests, exports, actions, tests, and simplification opportunities."
LastUpdated: 2026-07-26T21:20:00-04:00
WhatFor: "Supporting frontend evidence for the primary simplification report."
WhenToUse: "Use when reviewing package API, adapter, registry, test, and publishing cleanup."
---

# Frontend package audit

## Review
### Scope
The requested root `plan.md` and `progress.md` did not exist. The active ticket documents were present but still template-only, so this audit is based on verified source, consumers, stories, build outputs, and tests.
### Architecture map for a new intern
1. **React design system:** `src/components/{foundation,atoms,layout,molecules,organisms}`.
2. **Transport contract:** `src/widgets/ir/*` defines JSON-compatible nodes, props, cells, and actions.
3. **Translation layer:** 90 `*.widget.tsx` adapters convert transport props/actions into React props/callbacks.
4. **Dispatch:** `defaultRegistry.ts` registers adapters; `WidgetRenderer.tsx` walks nodes and invokes them.
5. **Host:** `useWidgetPage.ts` fetches pages; `app/App.tsx` owns routing, shells, shortcuts, server actions, and toasts.
6. **Authoring:** Go `widget.dsl` builders lower semantic specs into the React Widget IR.
7. **Packaging:** Vite produces the npm root, `/ir`, `/app`, `/scheduling`, and `/widgets/presets` entrypoints plus the embedded SPA.
The component layering is generally coherent and package code does not import Redux, router, backend services, or `web` internals.
### Findings
- **High — actions have two execution paths and confirmations run twice in the default app.**
  `app/App.tsx:87-95` confirms every action and then calls `dispatchWidgetAction` for non-server actions. That function confirms again at `widgets/actions.ts:44-56`. Server fetching/result-event/toast behavior is also duplicated between `app/App.tsx:96-128` and `widgets/actions.ts:114-138`. Consolidate this into one executor accepting an injectable server transport/API base.
- **High — the completed Widget DSL hard cutover still compiles thousands of lines of legacy archaeology.**
  Production registration exposes only `widget.dsl` (`pkg/widgetdsl/module.go:203-225`), but the same file retains six obsolete split-module constants, helper maps, recipes, loaders, and test registration (`module.go:14-21,36-186,227-237`). The directly implicated legacy implementation/test files total about 4,184 lines. Extract v3’s still-shared spec/ref helpers under neutral names, then delete split-module installation, `data.v2.dsl` declarations, and archaeology tests.
- **High — widget metadata has multiple drifting sources of truth.**
  A component is repeated in `RagWidgetType` (`widgets/ir/core.ts:22-110`), its props union, adapter, `defaultRegistry.ts`, YAML manifest, Go builder/lowering code, and the stale Go schema catalog. The 85 manifests contain roughly 1,395 lines; repository use is discovery/list/check only (`internal/widgetmanifest/discover.go:42-55`, `cmd/widget-codegen/main.go:86-140,195-223`), not runtime or generation. Five adapters have no manifest at all. Either make manifests generate the other artifacts or delete them; the current lint-only duplicate should not remain.
- **High — the advertised IR types provide less safety than their size suggests.**
  `ComponentNode.type` is `RagWidgetType | string`, props are an unrelated `WidgetProps` union, and `BaseWidgetProps` permits every key (`widgets/ir/core.ts:112-126`). Consequently component type and prop shape are not correlated despite the 709-line generated `props.d.ts`. Choose one honest contract: a correlated component map, or a deliberately open JSON node with adapter-local transport types.
- **Medium — 35 low-level adapters appear retained only for legacy/raw authoring.**
  Comparing registered adapters with production v3 lowering found 35 types not emitted by typed v3 paths, including `AppShell`, `AppNav`, `DashboardGrid`, individual context diagrams/cards, transcript leaf widgets, `SlideShell`, `DocumentListPanel`, and `AssetTile`. They account for approximately 740 adapter lines plus 33 manifests/515 lines. Typed v3 exposes semantic namespaces (`pkg/widgetdsl/v3_descriptors.go:107-150`), while `raw.component` remains an explicit escape hatch (`v3_descriptors.go:84-91`). Remove `raw.component`, migrate the static demo at `internal/api/dsl_handlers.go:12-66`, then prune adapters no typed authoring path can emit. Keep the underlying React components.
- **Medium — registry metadata and public partial registries are unused complexity.**
  Every one of the 90 adapters repeats `module: "widget.dsl"`, but registry lookup uses only `type` (`widgets/registry.ts:18-46`). `node` is also passed to every adapter but unused. The six partial registries are consumed only to build `defaultWidgetRegistry` (`defaultRegistry.ts:93-205`); repository consumers use the default registry. Delete `WidgetModule`, `module`, unused `node`, `entries`, `mergeWidgetRegistries`, and partial-registry exports; construct one default registry array.
- **Medium — the npm root is an indiscriminate public API and has hidden global effects.**
  `src/index.ts:1-10` automatically imports global CSS and star-exports components, hooks, fixtures, story palettes, registries, actions, and presets. The built root exposes 225 runtime symbols, while `web` imports 34. `context/index.ts:1-5`, `cms/index.ts:1-2`, and `scheduling/index.ts:1-3` publish fixture/story data. The dry-run package contains 377 files. Remove the implicit stylesheet import, move fixtures to Storybook/web, and expose only deliberate component, renderer, IR, and app APIs. `web` needs one explicit stylesheet import and two fixture imports relocated.
- **Medium — TypeScript domain presets duplicate the canonical Go authoring layer.**
  `widgets/presets/scheduling.ts:36-143` implements scheduling-to-IR recipes already represented by `widget.schedule`/`widget.time`; all four functions are used only by package stories. The 348-line CRM preset file is likewise story-only in repository code. Delete the `/widgets/presets` export and render stories from Go-generated golden IR or direct component stories.
- **Medium — several hand-authored React “twins” have no repository consumer beyond Storybook.**
  `CalendarMonthPanel`, `CalendarWeekPanel`, `BookingPagePanel`, `MeetingPollPanel`, `PollResultsPanel`, and `RecordShell` are exported publicly (`components/organisms/index.ts:7-9,21-23`) while Go DSL uses generic engines/compositions. `RecordShell.tsx:29-34` explicitly calls itself the hand-authored twin. Delete these surfaces unless an identified external React consumer requires them; repository migration is story deletion/replacement only.
- **Medium — `FormDialog` bypasses the documented React-first architecture.**
  Its directory contains only CSS, manifest, and adapter. The adapter itself owns the full dialog component, global event subscriptions, state, focus, form serialization, and rendering (`FormDialog.widget.tsx:20-131`). It has no component story despite being emitted by a golden DSL example. Extract a tested presentational `FormDialog` or keep the capability host-owned; do not leave a stateful organism hidden in an adapter.
- **Medium — schema/version artifacts have drifted.**
  Unreferenced `pkg/widgetschema/schema.go:3-65` advertises an incomplete 0.1.0 component list, while current lowering emits 0.2.0 (`pkg/widgetdsl/spec/lower.go:5-23`). Go TypeScript declarations require `schemaVersion` in one contract and make it optional in v3 (`pkg/widgetdsl/typescript.go:20-31,137-155`), while React omits it entirely (`hooks/useWidgetPage.ts:55-62`). Delete `pkg/widgetschema` and establish one page-envelope version policy.
- **Medium — compatibility tokens are permanent in practice but temporary in documentation.**
  `theme.css:18-40` describes `--mac-*` as a migration bridge, yet 98 CSS files contain 887 `--mac-*` references and guidelines permit them. Decide explicitly: make `--mac-*` canonical or perform a mechanical hard cutover to `--rag-*`. The ambiguous “bridge forever” state is worse than either choice.
- **High residual risk — behavioral coverage is inadequate for safe deletion.**
  The package has 126 stories but zero `*.test.*`/`*.spec.*` source files. `scripts/focused-checks.mjs:1-139` covers only calendar packing, month cells, style lookup, and shortcut logic—not renderer traversal, adapters, actions, host routing, dialogs, uploads, or server result behavior. Add narrow characterization tests before cutting the catalog.
### Prioritized simplification plan
1. **P0: correctness and easy deletion**
   - Unify action execution and fix double confirmation.
   - Flatten the default registry; remove singleton `module` metadata, partial registries, and unused registry methods.
   - Add renderer/action/FormDialog characterization tests.
2. **P1: finish the hard cutover**
   - Remove legacy split Go modules and v2 public declarations/tests.
   - Remove `widget.raw.component`.
   - Migrate `internal/api/dsl_handlers.go` to the current page-shell/typed vocabulary.
   - Prune the 35 legacy/raw-only adapters and their IR prop/manifests.
3. **P2: remove duplicate product APIs**
   - Delete TypeScript CRM/scheduling presets and their package subpath.
   - Delete story-only React organism twins unless an external consumer is identified.
   - Delete lint-only widget manifests, or first promote them to the sole generator input.
4. **P3: narrow publishing**
   - Stop root-exporting fixtures, story palettes, internal action helpers, and partial registries.
   - Make stylesheet loading explicit.
   - Emit declarations only for supported entrypoints.
   - Delete stale `pkg/widgetschema`.
5. **P4: token cleanup**
   - Choose one canonical token namespace and execute the decision mechanically.
### Migration estimate
- **In-repository React app:** 34 imported root symbols; most core primitives remain. Expected edits are an explicit CSS import, relocation of two context fixtures, and migration of one static IR demo.
- **Widget DSL tests/examples:** 44 golden examples currently contain 35 component types. Preserve these semantic outputs; deletions should target adapters outside this verified typed set after checking variable-based lowering.
- **Storybook:** highest churn. Many low-level IR stories must become direct React stories or semantic Go-golden fixtures.
- **Published npm/external repositories:** unknown and therefore the largest risk. Treat removal as a major release and scan sibling repositories before deleting root exports or raw component support.
- **No blocker to planning**, but action correctness and missing characterization tests block a safe mass deletion.
```acceptance-report
{
  "criteriaSatisfied": [
    {
      "id": "criterion-1",
      "status": "satisfied",
      "evidence": "Concrete severity-ranked findings cite package and DSL file/line evidence, with repository consumer, story, adapter, manifest, and golden-example migration estimates."
    }
  ],
  "changedFiles": [],
  "testsAddedOrUpdated": [],
  "commandsRun": [
    {
      "command": "pnpm --dir packages/rag-evaluation-site typecheck",
      "result": "passed",
      "summary": "TypeScript typecheck completed successfully."
    },
    {
      "command": "pnpm --dir packages/rag-evaluation-site test:focused",
      "result": "passed",
      "summary": "Focused utility checks passed."
    },
    {
      "command": "go test ./pkg/widgetdsl ./internal/widgetmanifest -count=1",
      "result": "passed",
      "summary": "Widget DSL and manifest tests passed."
    },
    {
      "command": "pnpm --dir packages/rag-evaluation-site pack:smoke",
      "result": "passed",
      "summary": "npm dry-run produced a 377-file, 668.3 kB unpacked package."
    }
  ],
  "validationOutput": [
    "Inventory: 90 widget adapters, 85 widget manifests, 126 stories, and zero package source test/spec files.",
    "Built root exports 225 runtime symbols; web imports 34 root symbols.",
    "No staged files were present; the existing ticket workspace remained untracked."
  ],
  "residualRisks": [
    "External npm and sibling-repository consumers were not exhaustively available from this repository.",
    "Storybook states are not automated behavioral assertions.",
    "The 35-adapter pruning candidate was derived from repository v3 lowering/literal usage and must be confirmed against persisted or externally authored Widget IR."
  ],
  "noStagedFiles": true,
  "diffSummary": "Read-only audit; no repository files changed.",
  "reviewFindings": [
    "high: packages/rag-evaluation-site/src/app/App.tsx:87-95 - non-server actions with confirmations are confirmed twice.",
    "high: pkg/widgetdsl/module.go:14-237 - production exposes only widget.dsl but legacy split modules remain compiled and test-maintained.",
    "high: packages/rag-evaluation-site/src/widgets/ir/core.ts:112-126 - component type and prop union are uncorrelated and open-ended.",
    "medium: packages/rag-evaluation-site/src/index.ts:1-10 - root import applies global CSS and exports fixtures/internal machinery.",
    "medium: packages/rag-evaluation-site/src/widgets/registry.ts:5-50 - singleton module metadata, adapter node argument, entries, and merge APIs are unnecessary.",
    "high: packages/rag-evaluation-site/scripts/focused-checks.mjs:1-139 - no behavioral tests cover renderer, actions, adapters, or host."
  ],
  "manualNotes": "The latest instruction prohibited further tool calls, so the completed audit is returned here for persistence to the authoritative frontend-audit-subagent.md path."
}
```
<summary>
This turn: Completed the frontend architecture and hard-cutover audit with prioritized evidence and migration estimates.
Session so far: Verified package structure, consumers, exports, adapters, manifests, stories, DSL code, build artifact, and focused validations.
Issues: Double action confirmation, retained legacy DSL machinery, duplicated catalogs/presets, oversized public API, and missing behavioral tests.
Next steps: Persist this report, then execute P0 characterization and action/registry simplification before catalog deletion.
</summary>
