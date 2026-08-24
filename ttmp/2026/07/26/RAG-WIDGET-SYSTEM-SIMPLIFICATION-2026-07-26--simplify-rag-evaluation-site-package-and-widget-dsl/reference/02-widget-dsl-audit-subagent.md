---
Title: Widget DSL specialist audit
Ticket: RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26
Status: active
Topics:
    - architecture
    - widget-dsl
    - code-quality
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Independent review of Widget DSL registration, builders, compatibility code, false APIs, consumers, and hard-cutover opportunities."
LastUpdated: 2026-07-26T21:20:00-04:00
WhatFor: "Supporting evidence for the primary simplification report."
WhenToUse: "Use when reviewing DSL-specific cleanup claims and migration risks."
---

# Widget DSL deep audit

## Review

- **Correct:** The production xgoja provider is already narrowed to one selected module: `pkg/xgoja/providers/widgetsite/provider.go:12-33` registers only `widget.dsl`, supplies its TypeScript descriptor, and embeds help. `pkg/widgetdsl/registrar.go:8-19` similarly registers only `widget.dsl` for engine runtimes. The provider and engine boundary tests verify that split modules cannot be required (`pkg/xgoja/providers/widgetsite/provider_test.go:14-49`, `pkg/widgetdsl/module_test.go:1096-1150`).
- **Correct:** The browser boundary is conceptually sound: JavaScript callbacks configure server-side builders, emitted actions/bindings remain serializable data, the React registry renders known component nodes, and the app owns network/action behavior (`pkg/xgoja/providers/widgetsite/doc/02-widget-dsl-js-api-reference.md:15-27`, `packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx:57-112`).
- **Blocker (high):** The claimed hard cutover is bypassable through the global native-module registry. `pkg/widgetdsl/module.go:14-21` still exports every legacy module name, `moduleSpecs` still defines them at `:141-187`, and `init()` globally registers all seven at `:249-252`. A runtime using go-go-goja's default-registry module selection can therefore select `ui.dsl`, `data.dsl`, `data.v2.dsl`, `context_window.dsl`, `course.dsl`, or `cms.dsl` despite `NewLoader` rejecting them. A runtime probe in this audit confirmed `modules.GetModule(name) != nil` for all seven names. This contradicts the hard-cutover comments at `module.go:218-229` and the provider docs at `doc/01-widget-dsl-getting-started.md:18-24,39`.
- **Blocker (high):** Several advertised slot APIs are inert. Context/CMS builders serialize a slot callback as only `{kind:"slot", registered:true}` (`pkg/widgetdsl/v3.go:652-695,946-1021`), losing the function and its output. The emitted golden proves this exact marker (`pkg/widgetdsl/testdata/v3/examples/38-context-empty-state.js:3-11`, `pkg/widgetdsl/testdata/v3/golden/38-context-empty-state.json:8-25`). Frontend contracts have no `emptySlot`, `legendSlot`, `messageSlot`, `assetSlot`, `rowSlot`, etc. (`packages/rag-evaluation-site/src/widgets/ir/props.ts:187-197,257-265,870-909`), and the adapters do not read them (`MediaLibraryPanel.widget.tsx:22-85`, `TranscriptWorkspacePanel.widget.tsx:5-29`). Thus the documented example's “No context parts yet” never reaches the browser. These methods should not remain public merely because descriptor parity tests see their names.
- **Note (high):** The real cross-repository host still depends on a broad compatibility facade and the raw component escape hatch. `/home/manuel/code/wesen/go-go-golems/go-go-course/cmd/go-go-course/server.js:9-11` loads `widget.dsl` and immediately recreates `ui`, `dataV2`, `contextWindow`, `courseDsl`, and `cmsDsl`. Its adapter uses `widget.raw.component` at `lib/widget-dsl-v3-adapter.js:1-14`, reimplements actions and a page envelope at `:16-72`, and recreates v2 field/collection builders at `:133-250`. The default migration-checker run found this raw-component call. Generated declarations remain stale too: `go-go-course/cmd/go-go-course/types/xgoja-modules.d.ts:3-15,732-740` still declares split modules even though the generated runtime selects only `widget.dsl`. This facade is the principal prerequisite for deleting compatibility code in this repository.
- **Note (high):** There are three page-version vocabularies and none is enforced by the browser. Public v3 emits `0.1.0` (`pkg/widgetdsl/v3.go:66-76` and all 44 goldens); the typed `spec.PageSpec` lowerer defaults to `0.2.0` (`pkg/widgetdsl/spec/lower.go:5-10`); the sibling adapter emits `widget.ir/v1` (`go-go-course/.../widget-dsl-v3-adapter.js:56-68`). Meanwhile `WidgetPageResponse` omits `schemaVersion` entirely and casts unvalidated JSON (`packages/rag-evaluation-site/src/hooks/useWidgetPage.ts:55-62,100-110`). The supposedly shared spec is still explicitly described as “v2 work” (`pkg/widgetdsl/spec/doc.go:1-4`). Either make `pkg/widgetschema.Version` the sole accepted value and validate it, or remove the version field until it gates real behavior; do not preserve three inert labels.
- **Note (medium):** Production v3 is structurally coupled to the legacy v2 implementation. V3 fields attach `v2Ref` handles and return `v2SchemaValue` (`pkg/widgetdsl/v3.go:1233-1287`); v3 collections use `v2Rows`, `attachV2Ref`, `mustV2Ref`, and `v3SelectionToV2` (`:1295-1337,2320-2333`). Those helpers live in the 403-line `v2_builders.go`, so deleting `data.v2.dsl` currently cannot delete its implementation. Extract one neutral opaque `schemaHandle` and a row conversion helper into v3/shared code; then remove v2 constructors and tests.
- **Note (medium):** The public API contains several migration-era aliases/overloads rather than one idiom. `data.collection` accepts both `(rows, callback)` and `(name, rows, callback)` (`pkg/widgetdsl/v3.go:1295-1318`); `EditorBuilder.submitPost` duplicates `submit` (`:1537-1544`); collection `toIR` duplicates `toNode`; `ActionsBuilder.button` duplicates `add`; every builder gets `.use`; and the empty reserved `style` namespace is exported (`pkg/widgetdsl/module.go:299-313`). The descriptors preserve all of these (`pkg/widgetdsl/v3_descriptors.go:204-243`). For a hard cutover, choose rows-first plus `.id(name)`, `submit`, `toNode`, `add`, and delete empty `style`; migrate first-party examples in one pass rather than maintaining aliases.
- **Note (medium):** `raw` is broader than demonstrated need and undermines both type safety and registry validation. It exposes text, arbitrary HTML, arbitrary component names, and fragment (`pkg/widgetdsl/v3.go:1909-1917`), while normal child coercion and arrays already cover text/fragments. The only tracked v3 raw example uses `raw.element`; no tracked repository example uses `raw.component`, but the sibling adapter does. `raw.component` permits any string and only fails later as `UnknownWidget` (`packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx:102-111`). After the sibling cutover, delete `raw.component`, `raw.text`, and `raw.fragment`; retain only a deliberately named HTML-element helper if arbitrary semantic HTML remains a requirement.
- **Note (medium):** The migration checker is expensive temporary infrastructure but is neither complete nor wired into CI. It hardcodes selected sibling/repository paths (`pkg/widgetdsl/migrationcheck/checker.go:52-84`), matches only exact call text `raw.component` or `widget.raw.component` (`:225-255`), and misses aliases/destructuring such as `const make = widget.raw.component; make(...)`. No active `.github` workflow invokes it. Its command alone loads 101 Go dependency packages, primarily for tree-sitter. Keep it only through the final sibling migration, make `--fail-on-findings` mandatory during that short period, then delete `cmd/widgetdsl-migration-checker` and `pkg/widgetdsl/migrationcheck` rather than promoting a migration tool into permanent architecture.
- **Note (medium):** The example tooling duplicates an evaluator and one tool actively emits legacy page structure. Both `cmd/widgetdsl-v3-examples/main.go:49-69` and `cmd/widgetdsl-v3-preview/main.go:141-169` independently create a Goja runtime, register the module, wrap source text, and export a `page` variable. Preview additionally hand-builds raw IR and injects legacy `meta.shell/navItems/activeNavItemId/maxWidth` (`cmd/widgetdsl-v3-preview/main.go:114-138,189-200`). The React app still contains fallback normalization for exactly this metadata plus a special root `CourseStudioShell` path (`packages/rag-evaluation-site/src/app/App.tsx:367-470`). Prefer the generated xgoja example host as the single preview/smoke path; delete both Go commands, or at minimum share one evaluator and make preview use typed `page.shell(widget.app.shell(...))`.
- **Note (medium):** Documentation/API generation has become a second schema rather than a simple reference. Runtime installation is hand-written (`pkg/widgetdsl/module.go:299-313` and ~2,800 lines of v3 builders), TypeScript is a 536-line hand-built string list (`pkg/widgetdsl/typescript.go:137-454`), and `v3_descriptors.go` repeats every namespace/builder/action-context in 336 more lines (`:73-327`). Another 378 lines of tests introspect names and parse TypeScript strings (`pkg/widgetdsl/v3_descriptors_test.go:14-192`) before snapshotting a 316-line generated Markdown document. These tests guarantee three lists drift together, not that props work—as the inert slot APIs demonstrate. Simplest cutover: treat the generated `.d.ts` as the machine-readable public inventory, retain focused runtime behavior tests, and replace the generated method-list help page with short conceptual docs plus a link to declarations.
- **Note (medium):** Legacy tests keep removed behavior executable and can hide accidental dependencies. `grammar.go` (537 lines), `v2_builders.go` (403), and large portions of `module.go` remain solely for split modules; tests repeatedly call `registerLegacyModulesForTests`, including a test explicitly named “KeepsOldModulesAvailable” (`pkg/widgetdsl/module_test.go:63-110`) and a v2-v3 output equivalence test (`:719-760`). There are also full split-module grammar/v2 test files and legacy TypeScript fixture tests. Once the sibling adapter is gone, archive history in Git rather than executable tests: retain negative tests proving legacy names are absent and v3 behavior/goldens only.
- **Note (low):** The provider itself has a needless one-argument loader closure (`pkg/xgoja/providers/widgetsite/provider.go:14-18`) even though it registers one constant module. Inline `NewModuleFactory` and consider replacing `NewLoader(moduleName)` with `NewLoader()`; the string parameter only exists because the implementation still models many modules.
- **Note (low):** Action dispatch has two server-action implementations. `RagEvaluationSiteApp` POSTs and refreshes its hook at `packages/rag-evaluation-site/src/app/App.tsx:87-130`; the generic dispatcher POSTs to a fixed URL and refreshes via `popstate` at `widgets/actions.ts:44-52,114-135`. This is outside the Go DSL core but is part of its runtime bridge. Inject a server transport/base URL and refresh callback into the central dispatcher instead of maintaining two response/toast paths.

## Intern map: public API and request flow

### Build-time/public API

1. A generated host selects provider `rag-widget-site`, module `widget.dsl`, alias `widget.dsl` in `xgoja.yaml` (for example `examples/xgoja/widget-site/xgoja.yaml:38-58`).
2. `widgetsite.Register` adds one provider module and help source (`pkg/xgoja/providers/widgetsite/provider.go:12-33`). The engine-native alternative is `widgetdsl.NewRegistrar()` (`pkg/widgetdsl/registrar.go:8-19`).
3. The host asks the provider for its module factory. `widgetdsl.NewLoader` accepts only `widget.dsl` (`pkg/widgetdsl/module.go:203-215`), and the loader creates a per-Goja `runtime` and installs exports (`:239-246,299-313`).
4. JavaScript calls `require("widget.dsl")`. The root API is:
   - `page`: page builder and final `.toPage()` envelope;
   - `app`: typed shell ownership/navigation;
   - `ui`: generic components/forms/layout;
   - `data`: fields, collections, matrices, selection/cells/activity;
   - `crm`, `cms`, `course`, `context`, `schedule`, `time`: domain helpers;
   - `act`: serializable actions;
   - `bind`: late-bound interaction values;
   - currently `raw` and empty `style`, both simplification candidates.
5. Builders mutate Go structs or `map[string]any`. Collection builders use `widgetdsl/spec.CollectionSpec`, validate selectively, and lower through `CollectionSpec.ToNode`; page builders lower with `v3PageToIR`. Returned builder handles are Goja objects, not JSON; route code must call `.toPage()`/`.toNode()`.

### HTTP/render/action flow

1. A jsverb route returns JSON from `page.toPage()` at `/api/widget/pages/{id}` (example route pattern documented in `pkg/xgoja/providers/widgetsite/doc/04-widget-dsl-v3-examples.md:157-177`).
2. `useWidgetPage` fetches and stores that JSON (`packages/rag-evaluation-site/src/hooks/useWidgetPage.ts:76-125`). There is currently no runtime schema validation.
3. `RagEvaluationSiteApp` resolves typed `page.shell`, or legacy metadata/root fallbacks, then invokes `WidgetRenderer` (`packages/rag-evaluation-site/src/app/App.tsx:367-470`).
4. `WidgetRenderer` switches on node kind; component nodes resolve by string `type` in `defaultWidgetRegistry`, then an adapter maps JSON props/actions to a React component (`packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx:71-112`).
5. An interaction supplies component context (`row`, `assetId`, `page`, etc.) to an `ActionSpec`. `bind` accessors are resolved in the browser. Non-server actions navigate/copy/download/dispatch events locally; server actions POST `{payload, context}` to `/api/widget/actions/{name}` (`packages/rag-evaluation-site/src/app/App.tsx:87-127`, `widgets/actions.ts:44-119`). A successful `{refresh:true}` fetches the page again.

## Prioritized hard-cutover plan

### P0 — make “single module” true

1. Migrate `go-go-course` pages from `createWidgetDslV3Adapters` to native namespaces; regenerate its embedded jsverbs/runtime and declarations. This is a cross-repository atomic change.
2. Remove the global `init()` registration and every split module constant/spec/loader. Make `NewLoader()` parameterless or private to the provider. Add a regression assertion that `modules.GetModule(oldName) == nil`, not only that one custom registry cannot require it.
3. Delete `grammar.go`, legacy recipes/helpers, v2 constructors, legacy TypeScript branches, and their tests. Extract only neutral v3 schema handles/row conversion first.

**Migration/test consequences:** update all sibling page scripts and browser tests; regenerate `go-go-course/internal/xgojaruntime/xgoja_embed`, generated runtime plans, and `.d.ts`; run all four site smoke suites plus go-go-course tests. In this repository, delete split-module/v2 fixtures and change all v3 tests to call `Register`, never `registerLegacyModulesForTests`.

### P0 — remove false public behavior

1. Delete unsupported dynamic domain slot methods (`context.diagram.legend/empty`, workspace message/annotation/empty, CMS asset/details/row/rowActions/filters) unless frontend components gain an actual serializable slot contract. Server callbacks cannot be deferred to browser-only item context, so deletion is simpler and more honest.
2. Replace example 38 with a supported component-level empty message/state or remove it. Add an end-to-end render assertion if slots are intentionally implemented instead.

**Migration/test consequences:** update TypeScript, descriptors/help, example 38 and its golden; search downstream source for the removed methods. Current tracked first-party use is only example 38, but downstream consumers remain a release risk.

### P1 — one transport schema and one v3 idiom

1. Choose one canonical page version. Prefer importing `pkg/widgetschema.Version` into the DSL and rejecting overrides; alternatively remove `schemaVersion` until the React host validates it. Remove dormant `0.2.0` and `widget.ir/v1` paths.
2. Standardize `data.collection(rows, configure)` with `.id(name)` and `data.fields(configure)` with an optional `.id/name` method if necessary. Remove `submitPost`, `toIR`, `button` alias, empty `style`, and redundant constructor overloads.
3. After downstream raw use reaches zero, remove `raw.component`, `raw.text`, and `raw.fragment`; decide explicitly whether arbitrary `raw.element` belongs in the supported language.

**Migration/test consequences:** mechanically rewrite 44 JS examples and the three in-repo generated-host verbs, regenerate 44 goldens and declarations, then run Go tests, TypeScript fixture compilation, frontend typecheck, and generated-host smoke. A schema-version change should add frontend rejection/normalization tests rather than only golden churn.

### P2 — delete migration and inventory machinery

1. Run the checker once with `--fail-on-findings` across both repositories, complete the migration, then delete the checker/tree-sitter dependencies.
2. Delete the generated API inventory descriptor layer or reduce it to documentation metadata not duplicated from runtime/DTS. Keep focused export and behavior tests; do not parse handwritten TypeScript in Go tests.
3. Use `examples/xgoja-widgetdsl-v3` as the single preview/smoke implementation. Delete `cmd/widgetdsl-v3-examples` and `cmd/widgetdsl-v3-preview`, or share one evaluator if a standalone renderer remains necessary.
4. Remove React legacy shell metadata and root-CourseStudio special cases after all page producers emit typed `page.shell`.
5. Consolidate server action dispatch behind one configurable transport.

**Migration/test consequences:** update README/help command references, release scripts, and generated assets; retain one provider runtime smoke, golden execution tests, a browser page/action smoke, and explicit absence tests for removed APIs.

## Scale and residual risks

- The audited tracked surface is about **18,050 lines**, including **44 JS examples and 44 goldens**. Core production complexity is concentrated in `module.go` (1,339 lines), `v3.go` (2,439), `typescript.go` (536), `v3_descriptors.go` (336), `spec` lowering/validation (1,500+), and legacy grammar/v2 files (940).
- All targeted Go tests pass, so the main risk is not current compile failure; it is a green test suite preserving dead or false contracts.
- The source repository had no root `plan.md` or `progress.md` at the requested paths; the ticket-local design document was only an empty template. Findings therefore derive from code, tests, generated artifacts, docs, and runtime probes.
- Downstream npm consumers or generated xgoja hosts not present locally may use overloads/raw/slot APIs. Before deletion, publish a one-time API usage search or major-version note; do not add compatibility shims back into runtime code.

```acceptance-report
{
  "criteriaSatisfied": [
    {
      "id": "criterion-1",
      "status": "satisfied",
      "evidence": "Review findings cite concrete file/line evidence with high/medium/low severity, include the request/runtime flow, prioritized simplifications, migration/test consequences, and residual risks."
    }
  ],
  "changedFiles": [
    "/home/manuel/code/wesen/claw-stuff/.pi-subagents/artifacts/outputs/637b8c50/ttmp/2026/07/26/RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26--simplify-rag-evaluation-site-package-and-widget-dsl/reference/widget-dsl-audit-subagent.md"
  ],
  "testsAddedOrUpdated": [],
  "commandsRun": [
    {
      "command": "go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... ./cmd/widgetdsl-migration-checker -count=1",
      "result": "passed",
      "summary": "Widget DSL, spec, migration checker, provider, and help packages passed."
    },
    {
      "command": "go run ./cmd/widgetdsl-migration-checker --json",
      "result": "passed",
      "summary": "Reported one remaining raw-component finding in the sibling go-go-course adapter."
    },
    {
      "command": "go run ./cmd/widgetdsl-migration-checker --json cmd pkg examples",
      "result": "passed",
      "summary": "Reported no findings in explicitly scanned in-repository JS/TS paths."
    },
    {
      "command": "runtime probe of go-go-goja modules.GetModule after importing pkg/widgetdsl",
      "result": "passed",
      "summary": "Confirmed all six legacy split names and widget.dsl remain globally registered."
    }
  ],
  "validationOutput": [
    "Relevant Go test packages all returned ok.",
    "Default migration scan found go-go-course/cmd/go-go-course/lib/widget-dsl-v3-adapter.js:13 raw.component.",
    "Global module registry probe returned true for ui.dsl, data.dsl, data.v2.dsl, widget.dsl, context_window.dsl, course.dsl, and cms.dsl.",
    "Tracked audit surface count: 18,050 lines; 44 v3 examples; 44 goldens."
  ],
  "residualRisks": [
    "Downstream hosts not present locally may use compatibility overloads, raw helpers, or inert slot APIs.",
    "The sibling go-go-course repository must be migrated and regenerated atomically before deleting legacy runtime code.",
    "Frontend does not validate schemaVersion, so changing version labels alone provides no compatibility protection.",
    "No browser smoke was run during this read-only audit; slot non-functionality was established from emitted goldens and absent frontend props/adapter reads."
  ],
  "noStagedFiles": true,
  "diffSummary": "Read-only source audit; only the requested external Markdown report was written.",
  "reviewFindings": [
    "high: pkg/widgetdsl/module.go:249-252 - init globally registers every legacy split module despite the claimed single-module hard cutover.",
    "high: pkg/widgetdsl/v3.go:652-695,946-1021 - public domain slot methods serialize inert markers that frontend props/adapters ignore.",
    "high: go-go-course/cmd/go-go-course/lib/widget-dsl-v3-adapter.js:1-250 - the main downstream host still recreates legacy APIs through raw.component and duplicate builders.",
    "high: pkg/widgetdsl/v3.go:70; pkg/widgetdsl/spec/lower.go:9; go-go-course widget-dsl-v3-adapter.js:61 - three unenforced schema-version labels coexist.",
    "medium: pkg/widgetdsl/v3.go:1233-1337 - v3 schema/collection builders still depend on v2Ref machinery.",
    "medium: pkg/widgetdsl/v3_descriptors.go and pkg/widgetdsl/typescript.go - runtime, DTS, descriptor, and generated help repeat the public API without catching behavioral defects.",
    "medium: cmd/widgetdsl-v3-preview/main.go:114-200 - preview duplicates evaluation/raw builders and emits legacy shell metadata.",
    "medium: pkg/widgetdsl/migrationcheck/checker.go:52-84,225-255 - a non-CI migration tool has hardcoded scope and incomplete raw alias detection.",
    "low: pkg/xgoja/providers/widgetsite/provider.go:14-18 - single-module provider retains a multi-module loader wrapper."
  ],
  "manualNotes": "Root plan.md and progress.md were absent; the ticket-local design document was an empty template. No repository source file was edited."
}
```
