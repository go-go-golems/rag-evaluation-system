---
Title: Consumer and compatibility inventory
Ticket: RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26
Status: active
Topics:
    - architecture
    - frontend
    - widget-dsl
    - code-quality
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Independent inventory of current consumers, exports, compatibility obligations, adapters, examples, and safe deletion seams."
LastUpdated: 2026-07-26T21:20:00-04:00
WhatFor: "Supporting consumer evidence for hard-cutover decisions."
WhenToUse: "Use before deleting exports, adapters, raw authoring, or legacy shell behavior."
---

# Code Context

## Files Retrieved
1. `packages/rag-evaluation-site/package.json` (lines 1-67) - published package/version, export map, CSS side effects, and consumer smoke contract.
2. `packages/rag-evaluation-site/src/index.ts` (lines 1-10) - very broad root barrel.
3. `packages/rag-evaluation-site/src/widgets/ir/core.ts` (lines 1-134) - unversioned JSON IR node contract and open-ended component type.
4. `packages/rag-evaluation-site/src/widgets/registry.ts` (lines 1-47) - registry is now explicitly single-module `"widget.dsl"`.
5. `packages/rag-evaluation-site/src/app/App.tsx` (lines 213-272) - browser shell and retained legacy CourseStudioShell rendering path.
6. `packages/rag-evaluation-site/src/widgets/defaultRegistry.ts` (lines 1-207) - central adapter aggregation (90 colocated `*.widget.tsx` adapters repository-wide).
7. `web/tsconfig.json` (lines 8-14) - monorepo aliases package root, `/app`, and `/ir` directly to source.
8. `web/src/services/api.ts` (lines 1-3) and `web/src/storybook/MockApiProvider.tsx` (lines 1-8) - actual `/ir` subpath consumers.
9. `pkg/widgetdsl/module.go` (lines 13-242) - seven historical module specs, production-only `widget.dsl`, and test-only legacy registration.
10. `pkg/widgetdsl/registrar.go` (lines 1-21) - public Go registrar installs only `widget.dsl`.
11. `pkg/widgetdsl/migrationcheck/checker.go` (lines 15-79, 182-263) - sibling-aware compatibility scanner and the exact legacy/raw patterns it enforces.
12. `pkg/xgoja/providers/widgetsite/provider.go` (lines 1-35) - externally consumable xgoja provider exports one module only.
13. `cmd/widgetdsl-v3-examples/main.go` (lines 1-65) and `cmd/widgetdsl-v3-preview/main.go` (lines 1-62) - only direct in-repo runtime consumers of public `widgetdsl.Register`.
14. `pkg/widgetdsl/spec/lower.go` (lines 460-484) - remaining typed-template-to-legacy-string lowering bridge.

## Key Code

- **Published JS API:** package `@go-go-golems/rag-evaluation-site` is currently `0.1.21`; exports root, `/ir`, `/app`, `/scheduling`, `/widgets/presets`, and two CSS paths (`package.json:1-55`). Root eagerly re-exports CMS, all components, context, hooks, scheduling, and widgets (`src/index.ts:1-10`).
- **IR is structurally permissive:** `WidgetNode = TextNode | ElementNode | ComponentNode`, while `ComponentNode.type` is `RagWidgetType | string` and props are open (`core.ts:7-25, 105-123`). There is no schema/version discriminator; compatibility is therefore encoded in adapter names/props and goldens rather than an explicit protocol version.
- **Runtime convergence already happened:** `WidgetModule` is literally `"widget.dsl"` (`registry.ts:5`), production `NewLoader` rejects every other module (`module.go:202-208`), `Register` installs only `widget.dsl` (`module.go:219-225`), and the xgoja provider exposes only that module (`provider.go:12-30`).
- **Legacy surface is archaeology, not production:** split modules and helpers remain in `module.go:13-185`, but only `registerLegacyModulesForTests` can install them (`module.go:227-237`). Numerous old tests call that private helper. This is the clearest deletion seam.
- **Explicit frontend compatibility bridge:** app rendering detects a root `CourseStudioShell` only when `page.shell` is absent, then uses a special renderer (`App.tsx:233-252`). Typed/root-owned pages take the ordinary registry route. This is a real payload compatibility obligation until sibling pages are confirmed migrated.
- **Explicit Go compatibility bridge:** typed `TemplateSpec` is still lowered through `LegacyTemplateString()` (`spec/lower.go:460,482-484`). Removing it changes serialized action confirmation strings.

## Architecture

`widget.dsl` JS builders in Go produce unversioned Widget IR JSON. The xgoja `rag-widget-site` provider exposes the builder and generated TypeScript help to generated binaries. The React package receives the JSON, looks up component-type adapters in `defaultWidgetRegistry`, and dispatches actions. The repository web app consumes the React package source through TS/Vite aliases; sibling `go-go-course` is the important external host (its checked-out `go.mod:13` pins this module at `v0.1.5`).

### Consumer evidence

- **React package, in repository:** 35 source files under `web` import the package; 32 import the root. `/ir` is used by exactly two files (`web/src/services/api.ts:2`; `web/src/storybook/MockApiProvider.tsx:6-7`). No `web/src` file imports `/app`. The `/app` contract is exercised by the package's generated consumer smoke script (`scripts/consumer-smoke.mjs:82-84`) and documented in README, not by the app source.
- **React package, known external:** repository docs repeatedly identify sibling `go-go-course` as the deployed consumer and historical npm pins. This scout directly confirmed sibling `../go-go-course/go.mod:13` pins the Go module `v0.1.5`; npm pin was not present at the probed `webapp/package.json` path, so its current frontend packaging must be checked before a release break.
- **Go package, in repository:** direct imports are provider, two commands, and tests; no product backend imports `pkg/widgetdsl` directly. External hosts normally consume it indirectly through `pkg/xgoja/providers/widgetsite`.
- **Stories/tests:** 15 WidgetRenderer story files and 90 adapters form a broad executable compatibility matrix; 44 v3 JS examples plus JSON goldens lock emitted IR. They are internal verification, not separate consumers.
- **Migration checker result (current workspace):** default scan included sibling `go-go-course` and found exactly one issue: `../go-go-course/cmd/go-go-course/lib/widget-dsl-v3-adapter.js:13`, a `widget.raw.component(...)` escape hatch. No legacy split-module imports or legacy shell-metadata findings were reported.
- **History:** commit `a028a9c` is explicitly `widgetdsl: hard cut over first-party hosts to widget.dsl`; later `62d75d4` validates host cutover. The current code matches that intent.

### High-value cleanup candidates

1. **High / safe with retained v3 tests:** delete legacy split-module implementation (`ui.dsl`, `data.dsl`, `data.v2.dsl`, `context_window.dsl`, `course.dsl`, `cms.dsl`), `newLegacyLoaderForTests`, and tests that only exercise those modules (`module.go:13-185, 210-217, 227-237`). Production already rejects/does not register them. Preserve reusable lowering/spec code used by v3.
2. **High / safe after replacing one sibling escape:** delete or sharply narrow `widget.raw.component`; the migration checker identifies only one actual first-party use at sibling `widget-dsl-v3-adapter.js:13`. This escape defeats component/prop inventory guarantees.
3. **Medium / likely safe after consumer check:** remove `/app` as a public package subpath if standalone embedding is no longer desired. It has no app-source consumer here, but the package smoke test and README deliberately promise it, so this is a semver hard cut rather than dead private code.
4. **Medium / likely safe:** reconsider `/scheduling` and `/widgets/presets` subpath exports. No repository import uses these package specifiers; root already re-exports scheduling. Audit sibling lockfiles/source before removal.
5. **Medium / do not hard-cut blindly:** remove `legacyCourseShellNode` special rendering only after running the sibling pages and proving all responses carry typed `page.shell`. The scanner checks some shell metadata but does not prove response shape.
6. **Medium / simplify contract:** introduce a single explicit IR schema/version (or declare intentional versionlessness) before deleting adapter props. Today open `string` component types and `[key:string]: unknown` props (`core.ts:105-123`) make static compatibility auditing incomplete.
7. **Low:** narrow the giant root barrel to intentional public entry points. This affects 32 in-repo imports, so first codemod them to stable subpaths; it is not presently a safe immediate cut.

## Residual Risks

- Repository search cannot enumerate npm users or all consumers of the public Go module; npm download/dependent data and GitHub code search were not available.
- Sibling `go-go-course` is present and scanned, but its current npm package location/pin was not established; only its Go module pin was confirmed.
- The migration checker catches imports, raw component calls, and selected legacy metadata, but not dynamic module names, arbitrary JSON payload producers, or all old prop shapes.
- Git history shows intentional cutover but does not eliminate consumers pinned to older releases.

## Start Here

Open `pkg/widgetdsl/module.go:13-242` first: it cleanly shows the production `widget.dsl` boundary versus the large test-only legacy implementation, making it the highest-confidence simplification target. For frontend compatibility, next inspect `packages/rag-evaluation-site/src/app/App.tsx:213-252` and the sibling page payloads.

```acceptance-report
{
  "criteriaSatisfied": [
    {
      "id": "criterion-1",
      "status": "satisfied",
      "evidence": "Concrete consumer, compatibility, severity-ranked cleanup, and residual-risk findings cite current repository paths and line ranges."
    }
  ],
  "changedFiles": [
    "/home/manuel/code/wesen/claw-stuff/.pi-subagents/artifacts/outputs/637b8c50/ttmp/2026/07/26/RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26--simplify-rag-evaluation-site-package-and-widget-dsl/reference/consumer-and-compatibility-inventory.md"
  ],
  "testsAddedOrUpdated": [],
  "commandsRun": [
    {
      "command": "go run ./cmd/widgetdsl-migration-checker --json",
      "result": "passed",
      "summary": "Scanned default repository and sibling sources; found one raw.component escape in go-go-course and no legacy module imports."
    },
    {
      "command": "repository rg/find consumer and adapter counts",
      "result": "passed",
      "summary": "Counted 35 web import files, 32 root-import files, 2 /ir consumers, 0 /app source consumers, 90 adapters, and 15 WidgetRenderer stories."
    },
    {
      "command": "git log/show inspection for packages/rag-evaluation-site and pkg/widgetdsl",
      "result": "passed",
      "summary": "Confirmed explicit first-party hard-cutover commit a028a9c and subsequent host validation history."
    }
  ],
  "validationOutput": [
    "Current production registration and provider expose only widget.dsl.",
    "Default migration scan reports only ../go-go-course/cmd/go-go-course/lib/widget-dsl-v3-adapter.js:13 raw.component."
  ],
  "residualRisks": [
    "Unknown npm and external Go-module consumers cannot be proven absent from repository-only evidence.",
    "Current sibling npm pin/location was not established; only go.mod v0.1.5 was confirmed.",
    "Unversioned permissive IR means static search cannot prove all payload compatibility."
  ],
  "noStagedFiles": true,
  "diffSummary": "Read-only code investigation; only the requested external markdown artifact was written.",
  "reviewFindings": [
    "high: pkg/widgetdsl/module.go:13-237 - production already hard-cuts to widget.dsl, leaving a large split-module implementation reachable only by private test registration.",
    "high: ../go-go-course/cmd/go-go-course/lib/widget-dsl-v3-adapter.js:13 - sole scanner-detected first-party raw component escape blocks complete adapter contract enforcement.",
    "medium: packages/rag-evaluation-site/src/app/App.tsx:233-252 - legacy CourseStudioShell payload path remains a real compatibility branch.",
    "medium: packages/rag-evaluation-site/src/widgets/ir/core.ts:105-123 - unversioned/open component and props types make breaking-change detection incomplete."
  ],
  "manualNotes": "Safest immediate hard cut is test-only legacy split DSL removal; public npm export removal still requires sibling/external consumer confirmation."
}
```
