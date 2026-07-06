---
Title: Investigation diary
Ticket: GOJA-DSL-PLAYBOOK
Status: active
Topics:
    - goja
    - dsl
    - fluent-builder
    - go
    - typescript
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/internal/xgojaruntime/xgoja_embed/jsverbs/minitrace_viz_site/lib/pages/admin-course-cms.js
      Note: Step 25 generated mirror refreshed after xgoja build (commit 34bc642)
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/pages/admin-common.js
      Note: Step 27 admin material tables data.v2.dsl migration (commit 0bf3820)
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js
      Note: Step 25 real admin agenda editor data.v2.dsl port (commit 34bc642)
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js
      Note: Step 24 cross-module DSL gallery examples (commit 30ab8c3)
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/pages/sessions.js
      Note: Step 27 session browse data.v2.dsl migration (commit 0bf3820)
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/webapp/package.json
      Note: Step 21 test script and Playwright dependency (commit 06aa1c9)
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/webapp/playwright.config.ts
      Note: Step 21 hotreload-host Playwright server setup (commit 06aa1c9)
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/webapp/tests/dsl-examples.spec.ts
      Note: |-
        Step 21 durable browser/action smoke tests (commit 06aa1c9)
        Step 24 module gallery browser test (commit 30ab8c3)
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/xgoja.package.yaml
      Note: Step 28 legacy data.dsl runtime removal (commit 49dc871)
    - Path: pkg/widgetdsl/typescript.go
      Note: Step 22 precise data.v2.dsl TypeScript declarations (commit dcd5156)
    - Path: pkg/widgetdsl/typescript_fixture_test.go
      Note: Step 23 runtime export parity and TypeScript positive/negative fixtures (commit cee7525)
    - Path: pkg/widgetdsl/typescript_test.go
      Note: Step 22 declaration shape and legacy-API absence tests (commit dcd5156)
    - Path: pkg/widgetdsl/v2/spec/lower.go
      Note: Step 26 lowering for explicit table action columns and table rowSelect (commit 166e8dc)
    - Path: pkg/widgetdsl/v2_builders.go
      Note: Step 26 v2 table action API and rowSelect split (commit 166e8dc)
    - Path: pkg/xgoja/providers/widgetsite/provider_test.go
      Note: Step 22 provider TypeScript descriptor coverage for data.v2.dsl (commit dcd5156)
    - Path: ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/05-rag-evaluation-system-dsl-overhaul-design-and-implementation-guide.md
      Note: Step 5 records the DSL overhaul design guide
ExternalSources: []
Summary: Chronological investigation diary for the Goja DSL catalogue and base research. Records what was tried, what worked, what failed, and what to do next.
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: Record the research journey so a senior researcher can resume without re-reading every source.
WhenToUse: Read before resuming work on this ticket.
---








# Investigation diary

## Goal

Catalogue every Goja-based DSL across go-go-goja, go-minitrace, rag-evaluation-system, and the goja-* family; classify the implementation patterns; and produce a resource logbook so a senior researcher can write the opinionated fluent-builder playbook.

## Step 1: Scope the request and inventory the repositories

The user asked for a docmgr ticket + playbook research on Goja DSLs using the fluent builder pattern (Go implementation, strict runtime typecheck, validation, compile-time types from generated declarations, composable grammar extensible with lambdas). The first task was to research the existing DSLs in `rag-evaluation-system/` (especially `data.*`, which is "not really all that great"), plus the primitives in `./go-go-goja` and `./go-minitrace`, plus `~/code/wesen/go-go-golems/goja-text` and `goja-bleve` "and others." Deliverables: a base research catalogue of DSLs (API, docs, examples), a research logbook tracking which resources are useful/out-of-date/need-updating, stored in the ticket and uploaded to reMarkable.

### Prompt Context

**User prompt (verbatim):** "Create a new docmgr about creating a great document / playbook on how to create goja DSLs that use the fluent builder pattern where the implementation of the DSL is in go, to allow for strict typechecking and validation at runtime (and at compile type from generated types), by building up a core composable 'grammar' of operators and functions that can then be extended with lambdas, so that it is possible to express complex compositions while being opinionated. First, research the DSLs in rag-evaluation-system/ (the latest one, data.* is not really all that great, but also from the primitives in ./go-go-goja and ./go-minitrace and ~/code/wesen/go-go-golems/goja-text and ~/code/wasen/go-go-golems/goja-bleve and others. Go through the existing resources extensively and compile a list of the different DSLs we have created over time, their API, where they are documented, examples, basically a base research from which we can then pass this off to a senior researcher to do a proper documentation and reflexion and assessment. Create a research logbook document in the ticket that allows us to keep track of which resources are useful / out of date / need updating. For all the documents and external resources you read, list: what you were researching, what you were looking for in this document in particular, why you chose it, how you found the resource itself, what you found useful in the document, what you didn't find useful, what is out of date / what was wrong, what would need updating. Store in ticket and upload to remarkable."

**Assistant interpretation:** Create a docmgr ticket that inventories all our Goja DSLs, classifies their implementation patterns, and produces a per-resource logbook (with the eight requested fields) to hand off to a senior researcher for the playbook. Upload to reMarkable.

**Inferred user intent:** Establish the evidence base before anyone writes the playbook, so the playbook is grounded in what we actually have (including the weak `data.dsl`) rather than invented from scratch.

### What I did
- Loaded the `ticket-research-docmgr-remarkable`, `docmgr`, and `diary` skill references (writing-style, deliverable-checklist, docmgr reference).
- Inventoried the working directory (`go-go-goja`, `go-minitrace`, `rag-evaluation-system`, `glazed`) and `~/code/wesen/go-go-golems/` (found `goja-bleve`, `goja-text`, `goja-git`, `goja-github-actions`, `goja-treesitter`, `goja-dbus`, `go-emrichen`, `glazed`).
- Read each repo's `README.md` and key `pkg/` source to map APIs.

### Why
- The ticket is explicitly a base-research handoff; I needed the full inventory before writing anything.

### What worked
- The `data.*` DSL was located at `rag-evaluation-system/pkg/widgetdsl/` — it is actually a family of five modules (`ui.dsl`, `data.dsl`, `context_window.dsl`, `course.dsl`, `cms.dsl`) registered in `module.go:139-178`. The grammar verbs live in `grammar.go`.
- The typed-ref substrate was located in `goja-bleve/pkg/api_types.go` (`refBase`, `refKind`, `getTypedRef[T]`, `newWrapper`).
- The `tsgen/spec` TypeScript model was located in `go-go-goja/pkg/tsgen/spec/types.go` — this is the compile-time-type substrate.
- `goja-treesitter` is a stub (`pkg/` has only `doc.go` + `logcopter.go`; README is ASCII art).

### What didn't work
- Initial `rg` for `data.` found no top-level `data` package — the DSL is registered under the `widgetdsl` package with `DataModuleName = "data.dsl"`. Resolved by grepping `DataModuleName`.
- `goja-treesitter` README is ASCII art only — no API to catalogue; flagged as skeleton in the logbook.

### What I learned
- The ecosystem has **five distinct implementation patterns** (typed-ref fluent, plain builder struct, map-IR helpers, hyperscript element, imperative object) plus two cross-cutting Go-native patterns (glazed functional-options, go-emrichen tag-operators).
- Only goja-bleve (and partially goja-dbus) realises the ticket's goal. `widgetdsl` (Pattern C) is the "not great" one — concrete evidence: `Props = Record<string, any>` in `typescript.go` and `panic(vm.NewGoError(...))` validation in `grammar.go`.
- goja-dbus's `pkg/dbusgoja/builders.go` is the cleanest composable grammar (bus → destination → object → interface → method → out → call).

### What was tricky to build
- Distinguishing Pattern A (hidden-key typed refs, goja-bleve) from the express side-channel `sync.Map` lookup (Pattern A/B hybrid). Both achieve type-safe handles but via different identity mechanisms. This matters for the playbook's "runtime typecheck substrate" decision record.
- The widgetdsl TypeScript generation is *intentionally* open-ended (`Props = Record<string, any>`) per a code comment ("individual component props remain open-ended by design"). So calling it "wrong" is too strong — it is a deliberate design choice that the playbook wants to reverse for type-checked builders. Captured this nuance in the logbook.

### What warrants a second pair of eyes
- The pattern classification (5 patterns + 2 cross-cutting). A reviewer might split Pattern A and the express hybrid differently.
- The four decision records in design-doc §7 are all `proposed` — they need ratification or rejection by the senior researcher.
- The claim that goja-bleve's machinery should be extracted into a shared `fluent` package has not been validated against go-go-goja's existing `pkg/engine` module-middleware layer (possible overlap).

### What should be done in the future
- Read the deferred resources (RAGEVAL-UI-GRAMMAR design docs 01 & 02; rag-eval-scripting-expansion; GOJA-DBUS-DESIGN ticket; go-emrichen spec/practice; go-go-goja jsverbs/jsdoc/xgoja internals) — see logbook §N.
- Verify staleness: RAGEVAL-UI-DSL RelatedFiles paths; goja-git Go version; goja-bleve "Phase 8" claims; FAISS fork currency.

### Code review instructions
- Start with `design-doc/01-goja-dsl-catalogue-and-base-research.md` — the inventory and pattern taxonomy.
- Cross-check against `reference/01-research-logbook-resource-assessment.md` — every claim should trace to a logbook entry.
- Validate by spot-checking the evidence anchors (file:line) in the design doc.

### Technical details
- Ticket created: `GOJA-DSL-PLAYBOOK` at `ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--…/`.
- Three docs created: design-doc (catalogue), reference (logbook), reference (this diary).
- Repos surveyed: 11 (goja-bleve, goja-dbus, go-minitrace, rag-evaluation-system, go-go-goja, goja-text, goja-git, goja-github-actions, goja-treesitter, glazed, go-emrichen).

## Step 2: Author the catalogue and logbook

With the inventory complete, I wrote the design-doc catalogue (`design-doc/01-goja-dsl-catalogue-and-base-research.md`) covering all 12 DSLs + 3 cross-cutting pieces with APIs, docs, examples, evidence anchors, pattern classification, gap analysis, decision records, open questions, and references. Then I wrote the research logbook (`reference/01-research-logbook-resource-assessment.md`) with the eight requested fields per resource and a summary table.

### Prompt Context
**User prompt (verbatim):** (see Step 1)

### What I did
- Wrote the catalogue (44 KB) with 9 sections: exec summary, problem/scope, pattern taxonomy, per-DSL inventory (12 entries), cross-cutting infrastructure (engine/xgoja/tsgen/jsverbs/jsdoc), gap analysis, decision records (4, all `proposed`), open questions (10), references.
- Wrote the logbook (37 KB) with per-resource entries grouped by repo (A–K), an external-resources note, a summary table, and prioritised follow-ups.
- Anchored every major claim to `file:line` evidence.

### Why
- The user asked for "a base research from which we can then pass this off to a senior researcher." Evidence-anchored catalogue + logbook is exactly that handoff.

### What worked
- The pattern taxonomy made the "why data.* is not great" concrete: Pattern C (map-IR + loose helpers + panic validation + open-ended TS) vs Pattern A (typed refs + `(v,error)` terminals + named types).
- The logbook's eight-field format surfaced concrete staleness (RAGEVAL-UI-DSL RelatedFiles may point to old `internal/dsl/widgetdsl/` paths; goja-bleve "Phase 8" claims).

### What didn't work
- Nothing failed; this was a writing step.

### What I learned
- go-minitrace's `import_builder.go` shows that "terminal" is not always `build()` — it has Detect/Convert/Preview/Diagnostics/Save. The playbook must allow multiple terminals.
- glazed's `schema.go` has a TODO ("This is a pretty messy interface") from the author — worth noting when citing glazed as the canonical model.

### What was tricky to build
- Keeping the logbook's eight fields consistent across ~35 resource entries without becoming repetitive. Resolved by using a fixed sub-heading per field and a status emoji legend.
- Deciding whether to read every deferred resource (RAGEVAL-UI-GRAMMAR design docs, etc.) in this pass or flag them for the senior researcher. Chose to flag them (logbook §N) because the user's intent is a handoff, not a final playbook.

### What warrants a second pair of eyes
- Completeness: did I miss any Goja DSL? Candidates to verify: `cozodb-goja`, `goja-dbus` (covered), any DSL inside `geppetto`/`dmeta`/`discord-bot`. The user said "and others" — a senior researcher should confirm the inventory is exhaustive.
- The gap analysis (design-doc §6) is the bridge to the playbook; a reviewer should confirm the missing pieces are the right ones.

### What should be done in the future
- Validate completeness against other goja-* repos not surveyed (`cozodb-goja`, etc.).
- Run `docmgr doctor` and upload to reMarkable (Step 3).

### Code review instructions
- Read design-doc §3 (pattern taxonomy) and §4 (inventory) first.
- Then read logbook §M (summary table) for the at-a-glance status.

### Technical details
- Design doc: `design-doc/01-goja-dsl-catalogue-and-base-research.md` (44 KB).
- Logbook: `reference/01-research-logbook-resource-assessment.md` (37 KB).

## Step 3: Bookkeeping, validation, and reMarkable upload

Relate the key source files to the ticket, update the changelog and tasks, run `docmgr doctor`, add missing vocabulary, and upload the bundle to reMarkable.

### Prompt Context
**User prompt (verbatim):** (see Step 1) — "Store in ticket and upload to remarkable."

### What I did
- Fixed YAML frontmatter in all three docs: `Topics` had to be a YAML list, not a comma string (validator rejected `cannot unmarshal !!str into []string`).
- Relate failed first with a relative `--doc` path (`ttmp/2026/...`); docmgr prepends the docs root, so the path doubled to `ttmp/ttmp/...`. Resolved by using the full absolute path for `--doc`.
- Related 7 key source files to the design-doc via `docmgr doc relate --doc <abs> --file-note "<abs>:reason"` (goja-bleve api_types/api_mapping, goja-dbus builders, minitracejs builders, widgetdsl grammar/typescript, tsgen spec).
- Added 9 tasks; checked off 1–7 after authoring/relating.
- `docmgr doctor` flagged two unknown topics (`fluent-builder`, `typescript`); added both to vocabulary via `docmgr vocab add`. Re-ran doctor: `All checks passed`.
- `remarquee status` OK; `remarquee cloud account` verified `user=wesen@ruinwesen.com sync_version=1.5`.
- Dry-run bundle upload confirmed layout + pandoc + upload plan.
- Real bundle upload: `OK: uploaded Goja DSL Playbook — base research.pdf -> /ai/2026/07/05/GOJA-DSL-PLAYBOOK`.
- Verified listing: `remarquee cloud ls /ai/2026/07/05/GOJA-DSL-PLAYBOOK --long` shows `[f] Goja DSL Playbook — base research`.
- Checked off tasks 8 and 9; final changelog entry recorded.

### Why
- The deliverable-checklist requires bookkeeping, validation, and reMarkable delivery. The user explicitly asked to "store in ticket and upload to remarkable."

### What worked
- The absolute-path rule for `--doc` resolved the relate failure once identified.
- Vocabulary add cleared the only doctor warning.
- Bundle upload succeeded on first real attempt after a clean dry-run.

### What didn't work
- `docmgr meta update --doc <relative>` failed with "no such file or directory" — same path-resolution quirk as `doc relate`. Skipped (the Summary is already in the frontmatter written by hand); not worth another round-trip.
- Initial frontmatter used comma-string `Topics` (copied from the create-ticket output style) which the validator rejects. Fixed in all three files.

### What I learned
- docmgr's `--doc` path resolution differs between subcommands: `validate frontmatter` and `doc list` accept paths relative to the docs root (no `ttmp/` prefix), while `doc relate --doc` and `meta update --doc` want the full absolute filesystem path. Worth a note in the docmgr skill.
- The `Topics` field must be a YAML list, not a comma string, even though `create-ticket` accepts comma-separated `--topics`.

### What was tricky to build
- The path-resolution inconsistency across docmgr subcommands cost two failed commands before I switched to absolute paths. The docs root is `rag-evaluation-system/ttmp` (not the workspace root), so "absolute from workspace root" in the skill docs is ambiguous when the docs root differs from the workspace root.

### What warrants a second pair of eyes
- The bundle PDF should be opened on reMarkable to confirm the ToC (toc-depth 2) renders all three docs and the code blocks are legible.
- Confirm the ticket should stay `active` (not closed) since it is a handoff to a senior researcher — docmgr prompted "Consider closing the ticket" but closing would be premature.

### What should be done in the future
- Senior researcher reads the bundle, ratifies or rejects the four decision records in design-doc §7, reads the deferred resources (logbook §N), and writes the opinionated fluent-builder playbook.
- Consider extracting goja-bleve's typed-ref machinery into a shared `fluent` package in go-go-goja (open question 1).

### Code review instructions
- Verify `docmgr doctor --ticket GOJA-DSL-PLAYBOOK --stale-after 30` prints `All checks passed`.
- Verify `remarquee cloud ls /ai/2026/07/05/GOJA-DSL-PLAYBOOK --long --non-interactive` lists the bundle.
- Read `design-doc/01-...md` §3 (pattern taxonomy) and §6 (gap analysis) for the handoff crux.

### Technical details
- Bundle: design-doc (44 KB) + logbook (37 KB) + diary, `--toc-depth 2`, name "Goja DSL Playbook — base research".
- Remote path: `/ai/2026/07/05/GOJA-DSL-PLAYBOOK`.
- Vocabulary added: `fluent-builder`, `typescript`.

## Step 4: Extend the catalogue with geppetto, discord-bot, researchctl, codesign

The user asked to also check `geppetto` and `discord-bot` ("those are interesting too"), and then to look at the `researchctl` and `codesign` DSLs in `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference`. This pass surfaced **four additional DSLs**, two of which (researchctl, codesign) are the strongest realisations of the ticket's full goal and materially change the gap analysis.

### Prompt Context
**User prompt (verbatim):** "check geppetto and discord-bot as well, those are interesting too" followed by "Look also at the researchctl and codesign DSLs in /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference"

### What I did
- Surveyed `geppetto/pkg/js/modules/geppetto/` (30+ `api_*.go` files): found the same hidden-key typed-ref substrate as goja-bleve (`__geppetto_ref`), but with **clone-on-each-step** immutable builders, and a **DTS parity test** enforcing the generated `geppetto.d.ts` against the runtime export surface.
- Surveyed `discord-bot`: found the `defineBot` registration DSL and the `require("ui")` Go-side typed-builder DSL surfaced via **Goja Proxy traps** (returns typed builders, wrong-parent calls fail loudly, raw JS objects rejected).
- Surveyed `researchctl/pkg/gojamodules/researchctl/`: found the **lambda-configurator** pattern (`project(name).goal(title, g => g.id(...).status(...))`) with `.toSpec()`/`.validate()` terminals and a precise `ProjectBuilder` TS interface.
- Surveyed `researchctl/pkg/gojamodules/codesign/`: found the **strongest model** — `runSpec(name)` fluent builder + `.topology(fn)`/`.workload(fn)`/`.metrics(fn)` lambda configurators + `.use(fragment)` `FragmentFn<T>` composition + JS-lambda runtime callbacks (`jsDevice`/`policyCallback`) validated with `goja.AssertFunction` + `modules.TypeScriptDeclarer` emitting precise `RunSpecBuilder`/`TopologyBuilder`/`MetricsBuilder`/`RunSpecLike` interfaces.
- Updated the design-doc: exec summary (12→16 DSLs, five→six+ patterns), pattern taxonomy (added Pattern A′ clone-on-each-step, Pattern F lambda-configurator+`.use()`, Pattern G Proxy traps), four new inventory entries (§4.13–4.16), gap analysis (composable-grammar and lambdas are no longer the main gap — they are realised), a new decision record (hidden-key vs Proxy traps), and three new open questions.
- Updated the logbook: four new resource sections (O/P/Q/R) with the eight-field assessment, plus four new summary-table rows.

### Why
- The user explicitly widened the scope; these DSLs were missing from the first pass and two of them are the best existing models for the playbook.

### What worked
- codesign turned out to realise **all seven** of the ticket's requirements (fluent, Go, typecheck, validation, compile-time types, composable grammar, lambdas) — the gap analysis flipped from "composable grammar is missing" to "composable grammar is realised; the remaining gap is extraction + standardisation."
- geppetto's DTS parity test is the missing compile-time-type enforcement piece — a test that guarantees declarations match runtime.

### What didn't work
- Could not find the discord-bot `ui` Proxy-trap Go implementation in `pkg/` via grep (the `ui` module is registered via the xgoja provider and may live under `internal/jsdiscord/` or be generated); documented it from the API-reference doc and the design rules there rather than the Go source. Flagged for the senior researcher to locate the Proxy-trap implementation.

### What I learned
- The lambda-configurator pattern (`g => g.id(...).status(...)` applied via `applyBuilderCallback`) is the answer to the ticket's "extended with lambdas" requirement — it already exists in production code.
- `.use(fragment)` with `FragmentFn<T>` is the composable-grammar mechanism — reusable builder lambdas applied to any builder of type `T`.
- There are now **two typed-builder mechanisms** (hidden-key refs vs Proxy traps) and **two mutation models** (same-object vs clone-on-each-step); the playbook must pick defaults.
- `modules.NativeModule` + `modules.TypeScriptDeclarer` is the canonical go-go-goja module registration interface for emitting typed declarations.

### What was tricky to build
- Integrating four new DSLs without rewriting the design doc. Resolved by appending §4.13–4.16 and editing the exec summary / pattern taxonomy / gap analysis / decision records / references in place. Had to re-balance the exec summary's "strongest model" claims: goja-bleve is still the typecheck substrate, but researchctl/codesign are now the composable-grammar + lambda models, and geppetto is the type-parity-test model.
- The codesign entry needed a 7-point evidence checklist to substantiate the claim that it realises the full goal — without that, the claim would be unanchored.

### What warrants a second pair of eyes
- The claim that codesign realises all seven requirements — a reviewer should verify each ✅ against `typescript.go` and `builders.go`.
- The Pattern A′ vs Pattern A distinction (clone-on-each-step vs same-object mutation) — a reviewer should confirm geppetto always clones and goja-bleve always mutates.
- The Proxy-trap characterisation of discord-bot `ui` is from the API-reference doc, not the Go source — needs source confirmation.

### What should be done in the future
- Locate and read the discord-bot `ui` Proxy-trap Go implementation.
- Confirm whether `cozodb-goja`, `dmeta`, `discord-bot`'s `database`/`timer` modules contain further DSLs.
- The senior researcher can now point at codesign as the reference implementation when writing the playbook.

### Code review instructions
- Read design-doc §4.15 (researchctl) and §4.16 (codesign) for the new models.
- Cross-check the gap analysis table (§6) — the "composable grammar" and "lambdas" rows now say "Realised."
- Read the new decision record (hidden-key vs Proxy traps) and open questions 11–13.

### Technical details
- Design doc grew from ~44 KB to ~62 KB; logbook from ~37 KB to ~52 KB.
- New patterns: A′ (clone-on-each-step), F (lambda-configurator + `.use()`), G (Proxy traps).
- New decision record: typed-builder mechanism (hidden-key vs Proxy traps).

## Step 5: Write the rag-evaluation-system DSL overhaul guide

The user asked for a full-system overhaul guide after the catalogue, self-assessment, independent assessment, and deep-dive had established that the current Widget DSL vocabulary is useful but the map-IR implementation substrate is not. I reread all four design documents and then investigated the current `rag-evaluation-system` renderer, Widget IR, action dispatcher, widget manifests, module exports, recipes, and the real `go-go-course` action/page host before writing the new intern-facing guide.

The resulting document reframes the target architecture around typed Go-side intent specs, hybrid hyperscript-plus-builder authoring, optional lambda configurators, stable server-action registration, serializable action templates, strict option decoding, declaration parity, and incremental migration of all five DSL modules.

### Prompt Context

**User prompt (verbatim):** "so read all 4 documents again, and do any other research needed. 

We want to overhaul the rag-evaluation-system DSLs (all of them).

What we want is a powerful yet opinionatedly simple set of APIs to build user interfaces for the rag-evaluation-system UI toolkit.

We can revisit the IR format, including the react widget renderer, and would like to benefit from the best design possible given all we have.

Furthermore, if we for example allow lambdas to be used as ways to say, handle callbacks or other compositional needs (as mentioned, codesign is pretty interesting), we need to make sure that the renderer <> jsruntime setup works properly, since actions will have to be serialized and then deserialized and dispatched properly.

Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.


Brainstorm different ideas / choices, make a lot of small to complex examples of the different APIs so I can jduge which forms I like the most, doing both opinionated and programs leveraging the more flexible side of things of the API. I want to be able to build rich and consistent UIs for a course / cms / general webpage system."

**Assistant interpretation:** Produce the concrete widgetdsl v2 overhaul guide for every module, including renderer/IR/action runtime considerations, many API alternatives/examples, implementation phases, and reMarkable upload.

**Inferred user intent:** Decide the next durable API architecture before implementing it, with enough examples to judge taste and enough technical guidance for an intern to execute safely.

### What I did
- Reread design-docs 01–04 in the GOJA-DSL-PLAYBOOK ticket.
- Inspected current module exports in `pkg/widgetdsl/module.go`, grammar in `pkg/widgetdsl/grammar.go`, declarations in `pkg/widgetdsl/typescript.go`, tests, manifests, React `ir.ts`, `WidgetRenderer.tsx`, `actions.ts`, registry/adapters, and `RagEvaluationSiteApp` action handling.
- Inspected real consumer/server behavior in `go-go-course/cmd/go-go-course/server.js`, `lib/course-pages.js`, and `lib/pages/admin-course-cms.js`.
- Created `design-doc/05-rag-evaluation-system-dsl-overhaul-design-and-implementation-guide.md`.

### Why
- The earlier documents were research and assessment. This step turns them into a concrete architecture and implementation guide for overhauling the DSLs, including the action serialization problem that determines whether lambdas can be used safely.

### What worked
- The current renderer/action architecture is already close to the right shape: Widget IR is JSON, adapters dispatch typed-ish contexts, `RagEvaluationSiteApp` posts server actions to `/api/widget/actions/:name`, and `go-go-course` already dispatches those names.
- The best target is a hybrid: hyperscript/tree helpers for structure, typed fluent builders for schemas/collections/actions/marks, and stable server action names for runtime callbacks.

### What didn't work
- A first manifest-summary script failed because two widget YAML manifests had missing `module` values, causing Python to compare `None` and `str` while sorting. I reran the script after stringifying missing values.

### What I learned
- The Widget IR/React side is much more typed than the Goja declarations suggest; the missing layer is typed Go-side authoring specs.
- The action boundary must remain data-only. Lambdas are safe as builder configurators; runtime callbacks need a server-side registry and stable serialized action names/templates.

### What was tricky to build
- The hardest design constraint was reconciling codesign-style lambdas with a browser renderer that only receives JSON. The solution is to split lambdas into builder-time configurators versus registered runtime handlers. Payload mapping and confirm text should compile to serializable template descriptors, not closures.

### What warrants a second pair of eyes
- The proposed `ActionSpec` v2 / payload-template design needs review by whoever owns `RagEvaluationSiteApp` and the xgoja/express server-action host.
- The module migration order (`data.collection` agenda first, then CMS media library, then context/course marks) should be checked against product priorities.

### What should be done in the future
- Create a follow-up implementation ticket (`RAGEVAL-DSL-V2` or similar) and implement the phases in the guide.
- Decide which API examples Manuel prefers before freezing method names.

### Code review instructions
- Start with design-doc 05 sections 4–8: target architecture, public API model, action/runtime design, and implementation plan.
- Validate the evidence anchors by checking `pkg/widgetdsl/module.go`, `pkg/widgetdsl/grammar.go`, `packages/rag-evaluation-site/src/widgets/actions.ts`, and `go-go-course/cmd/go-go-course/server.js`.

### Technical details
- New doc path: `ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/05-rag-evaluation-system-dsl-overhaul-design-and-implementation-guide.md`.
- The doc includes ASCII architecture diagrams, pseudocode, TypeScript/Go API sketches, API alternatives, simple/rich examples for course/CMS/general pages, phased implementation tasks, tests, and decision records.

## Step 6: Revise the DSL overhaul guide for hard cutover instead of compatibility

The user clarified that the new `rag-evaluation-system` DSL design does not need backwards compatibility. That changes the optimization target: the document should not recommend v1 compatibility facades, strict-decoding old option bags, or introducing builders alongside current APIs. Instead, the v2 design should cut over cleanly, delete old public APIs, and use the current DSL only as vocabulary evidence.

I updated the guide to make hard cutover explicit in the frontmatter, reading guide, executive summary, architecture rules, Action IR section, schema recommendation, brainstormed API choices, implementation plan, decision records, risks, and next-ticket proposal.

### Prompt Context

**User prompt (verbatim):** "Update."

**Assistant interpretation:** Apply the prior hard-cutover recommendations to the existing `design-doc/05-rag-evaluation-system-dsl-overhaul-design-and-implementation-guide.md` document.

**Inferred user intent:** Remove compatibility-oriented ambiguity from the guide so the future implementation maximizes v2 simplicity and type safety.

### What I did
- Changed the document framing from compatibility migration to hard-cutover replacement.
- Added a new `No backwards compatibility requirement` section.
- Rewrote the implementation plan as hard-cutover phases.
- Replaced the compatibility-facades decision record with a `no compatibility shims` decision.
- Reframed old option-bag APIs as rejected baselines or unsafe escape hatches.
- Updated risks and the next implementation ticket proposal.

### Why
- Backwards compatibility would preserve exactly the v1 complexity the redesign is meant to remove: raw maps, option bags, marker maps, string switches, weak declarations, and recipes that bypass typed specs.

### What worked
- The document now states that old DSL APIs are research data / vocabulary evidence, not public contracts.
- The new implementation plan prioritizes typed v2 specs/builders, Action IR v2, `widget.unsafe`, a real page rewrite, and removal of v1 public exports.

### What didn't work
- N/A. This was a documentation update only.

### What I learned
- Hard cutover simplifies both the implementation plan and the API story: fewer transitional states, fewer aliases, and stronger tests that assert old APIs are absent.

### What was tricky to build
- Keeping the brainstorming section without making rejected v1-shaped APIs look equally viable. I left the strict-props/options-bag approach as an explicit rejected baseline rather than deleting it, because it is still useful for comparison and model trials.

### What warrants a second pair of eyes
- Review whether the proposed `widget.unsafe` namespace should exist at all, or whether the hard cutover should forbid raw Widget IR in authored pages entirely.
- Review whether the chosen schema API should be the hybrid `.field(name, data.f.primary()...)` shape or an even smaller fluent-only shape.

### What should be done in the future
- If the hard-cutover direction is accepted, rename the follow-up ticket proposal to `RAGEVAL-DSL-V2-CUTOVER` and implement the new phases.

### Code review instructions
- Start with design-doc 05 sections 3.4, 8, 10 Decision 6, and 12.
- Search the guide for compatibility language and confirm it only appears as rejected context.

### Technical details
- Updated document: `ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/05-rag-evaluation-system-dsl-overhaul-design-and-implementation-guide.md`.

## Step 7: Add event timelines and a precise cutover task tracker

The user asked to store the chat-level walkthroughs as a separate ticket document and to turn the hard-cutover direction into phases and detailed tasks that can be tracked during implementation. I created a companion design document that explains the simple table, selectable table, and master-detail editor from authoring API down to Widget IR, React execution, browser history, HTTP requests, and backend handlers.

This step also turns the broad implementation plan into docmgr tasks. The companion document now doubles as a behavioral specification: when we implement each phase, the demos and code should either preserve these event timelines or explicitly update the document in the same commit.

### Prompt Context

**User prompt (verbatim):** "Ok, store all of these in a separate document in the ticket, as a companion to the design document. 

Then, create detailed phases and for each phase detailed tasks (so we can precisely keep track of our progress), and then implement task by task, keeping a detailed diary as you work, and committing at regular intervals. Also build demo sites as you go as examples (make sure the examples are properly updated / removed if deprecated)."

**Assistant interpretation:** Create a new GOJA-DSL-PLAYBOOK companion document for the event timelines and establish a precise implementation tracker, then proceed with incremental tracked work and commits.

**Inferred user intent:** Make the DSL cutover execution observable and reviewable, with no vague implementation phases and no stale examples left behind.

### What I did
- Added `design-doc/06-widget-dsl-event-timelines-and-cutover-task-plan.md`.
- Documented simple table, selectable table, and master-detail editor timelines.
- Included HTTP requests, frontend code paths, backend code paths, Widget IR shapes, and current behavior caveats.
- Added P0-P7 implementation phases with detailed tasks.
- Added docmgr tasks 11-40 corresponding to the phase tracker.
- Checked off tasks 11 and 12.
- Related the new document to the key runtime/frontend/backend files.
- Updated the changelog and ran `docmgr doctor` successfully.

### Why
- The DSL API shape is not enough; implementation needs event-level semantics so we do not accidentally break URL selection, native form behavior, server action refresh, or demos.

### What worked
- The new companion document provides a concrete operational contract for the first three examples and a precise tracker for the hard-cutover work.
- `docmgr doctor --ticket GOJA-DSL-PLAYBOOK --stale-after 30` passed after adding the document and relations.

### What didn't work
- N/A. This was a documentation and planning step.

### What I learned
- The current master-detail path is already a useful behavioral fixture: URL selection, native form submit, server-action row buttons, and refresh all exist today. The v2 work can focus on typed authoring and stronger contracts rather than inventing entirely new runtime semantics.

### What was tricky to build
- The master-detail delete path revealed a subtle action-dispatch caveat: generic `dispatchWidgetAction` handles `confirm`, but `App.tsx` handles server actions directly. V2 should centralize confirmation so server actions cannot bypass it. The action-cell row key also currently uses `rowKey(row, "file")`, which is wrong for generic tables and should be fixed in the Action IR/context phase.

### What warrants a second pair of eyes
- Confirm whether native form submit should remain a blessed v2 save path, or whether a fetch-based form action should also be introduced.
- Review whether delete should clear URL selection after removing the selected row.
- Review where demo pages should live: `go-go-course`, Storybook, or both.

### What should be done in the future
- Start P0.3/P0.4 by inventorying existing examples and adding/identifying live demo pages for the three foundational flows.

### Code review instructions
- Start with `design-doc/06-widget-dsl-event-timelines-and-cutover-task-plan.md`.
- Verify that the timelines match the referenced files: `grammar.go`, `App.tsx`, `actions.ts`, `DataTable.widget.tsx`, `FormPanel.tsx`, and `go-go-course/server.js`.
- Validate with `docmgr doctor --ticket GOJA-DSL-PLAYBOOK --stale-after 30`.

### Technical details
- New document: `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/06-widget-dsl-event-timelines-and-cutover-task-plan.md`.
- Doctor command: `cd rag-evaluation-system && docmgr doctor --ticket GOJA-DSL-PLAYBOOK --stale-after 30`.

## Step 8: Inventory current demos and mark the first baseline tasks complete

After creating the operational companion document, I started Phase 0 by inventorying existing live pages and Storybook examples. The key finding is that the repository already has strong behavioral evidence but not clean v2 demos: `/pages/sessions` demonstrates selectable table behavior through direct `dataDsl.dataTable`, while `/pages/admin-course-cms` demonstrates master-detail behavior through the current `dataDsl.collection` option-bag grammar.

I updated the companion document with this inventory and a deprecated-example policy: once a v2 demo exists, equivalent v1 option-bag examples must be removed from public docs or moved to a clearly historical section so agents do not imitate deprecated APIs.

### Prompt Context

**User prompt (verbatim):** (same as Step 7)

**Assistant interpretation:** Continue executing the newly created P0 tasks by identifying existing examples and demo gaps.

**Inferred user intent:** Ensure implementation starts with a clear baseline and that examples will stay synchronized with the API cutover.

### What I did
- Searched current `go-go-course` page files and WidgetRenderer stories.
- Identified `/pages/sessions` as current selectable-table evidence.
- Identified `/pages/admin-course-cms` as current master-detail evidence.
- Identified Storybook examples that cover `DataTable` and `FormPanel` behavior.
- Added a `Demo and Example Inventory` section to design-doc 06.
- Checked off tasks 13 and 14.
- Updated the changelog and ran `docmgr doctor` successfully.

### Why
- The next implementation tasks need to know which examples already exist, which are only indirect evidence, and which demos must be built from scratch.

### What worked
- `rg` quickly showed that `dataDsl.collection` has only one real consumer: the admin CMS agenda editor. This makes the hard-cutover real-page rewrite tractable.

### What didn't work
- There is no dedicated simplest-table live demo yet. The current examples demonstrate table behavior, but not the exact future v2 authoring path.

### What I learned
- `/pages/sessions` is a useful selectable-table runtime fixture even though it bypasses the grammar collection API. It should become a v2 demo or be paired with one.

### What was tricky to build
- Distinguishing behavioral fixtures from public examples: current pages are valuable for runtime semantics, but their authoring forms should not all remain public teaching material after v2 lands.

### What warrants a second pair of eyes
- Confirm whether the first dedicated DSL demos should be added to `go-go-course` live pages, Storybook, or both.

### What should be done in the future
- Record baseline validation commands and then add the first dedicated demo page/safe example route.

### Code review instructions
- Review the `Demo and Example Inventory` section in design-doc 06.
- Check the related files `/pages/sessions`, `/pages/admin-course-cms`, and existing WidgetRenderer stories.

### Technical details
- Commands used included `rg --files ../go-go-course/cmd/go-go-course/lib/pages`, `rg --files packages/rag-evaluation-site/src/widgets | rg 'stories\\.tsx$'`, and `rg -n "dataDsl\\.collection|data\\.collection\\(" ...`.

## Step 9: Record baseline validation before code implementation

Before changing implementation code, I ran the core validation commands for the current Widget DSL and React package. This gives the v2 cutover a known-good baseline: if later changes fail, we can distinguish new regressions from pre-existing failures.

The baseline passed for `pkg/widgetdsl` Go tests, the React package TypeScript typecheck, the React package production build, and docmgr doctor. The production build generated `dist` output internally but did not leave tracked git changes.

### Prompt Context

**User prompt (verbatim):** (same as Step 7)

**Assistant interpretation:** Continue the task-by-task implementation workflow by completing P0.5: record baseline commands and failures before implementation.

**Inferred user intent:** Make the cutover safe to review and bisect by validating before code changes begin.

### What I did
- Ran `go test ./pkg/widgetdsl -count=1`.
- Ran `pnpm --dir packages/rag-evaluation-site typecheck`.
- Ran `pnpm --dir packages/rag-evaluation-site build`.
- Ran `docmgr doctor --ticket GOJA-DSL-PLAYBOOK --stale-after 30`.
- Added the baseline command/results section to design-doc 06.
- Checked off task 15.
- Updated the changelog.

### Why
- A hard cutover should proceed with frequent validation and small commits. Establishing a clean baseline prevents chasing unrelated issues later.

### What worked
- All baseline commands passed.

### What didn't work
- N/A.

### What I learned
- The current frontend package builds cleanly after the documentation-only changes, so implementation can start from a stable state.

### What was tricky to build
- No implementation was changed in this step. The only operational detail is to avoid committing generated build artifacts; `git status` stayed clean before the doc updates.

### What warrants a second pair of eyes
- Confirm whether additional baseline checks should include Storybook build, full repo `go test ./...`, or an end-to-end browser smoke test before P1 begins.

### What should be done in the future
- Begin P1 typed v2 spec package or add first dedicated demo page, depending on whether we want model-facing examples before or after the v2 substrate exists.

### Code review instructions
- Review the `Baseline validation commands` section in design-doc 06.
- Re-run the commands listed there if validating the baseline.

### Technical details
- Commands passed exactly as recorded in design-doc 06.

## Step 10: Start P1 with the typed v2 spec package

I began implementation with the smallest code-bearing P1 task: create the typed v2 spec package. This package does not expose Goja builders yet; it defines the intermediate intent model that future builders will produce and future lowering code will convert into current Widget IR nodes.

The new types cover the major concepts from the event timelines: pages, nodes, sections, schemas, fields, collections, URL selection, arrangements, native submit specs, server/browser actions, payload templates, text/path templates, server result policy, and validation issues.

### Prompt Context

**User prompt (verbatim):** (same as Step 7)

**Assistant interpretation:** Move from planning into task-by-task implementation, starting with P1.1 typed v2 spec package.

**Inferred user intent:** Establish a typed substrate before implementing fluent Goja APIs so v2 does not repeat the v1 map-IR authoring problem.

**Commit (code):** pending — "Add widget DSL v2 spec skeleton"

### What I did
- Added `pkg/widgetdsl/v2/spec/doc.go`.
- Added `pkg/widgetdsl/v2/spec/types.go`.
- Defined initial typed structs/enums for page/node/schema/field/collection/action/template/validation concepts.
- Ran `gofmt -w pkg/widgetdsl/v2/spec/*.go`.
- Ran `go test ./pkg/widgetdsl/... -count=1`.
- Checked off task 16.
- Related the new files to design-doc 06 and updated changelog.

### Why
- The hard-cutover design requires typed intent specs underneath the public DSL. Starting here makes later builder work concrete and reviewable.

### What worked
- The new package compiles, and existing widgetdsl tests still pass.

### What didn't work
- N/A. This task intentionally added type skeletons only; validation/lowering are later tasks.

### What I learned
- The event timelines map cleanly to a compact set of typed concepts. The master-detail editor needs `SelectionSpec`, `SubmitSpec`, collection-level actions, and `TemplateSpec` from the beginning.

### What was tricky to build
- Avoiding premature implementation: the package needs enough structure to guide v2 but should not hard-code every renderer detail before validation/lowering tasks. I kept lowering out of this step.

### What warrants a second pair of eyes
- Review whether `JSONValue any` should be tightened immediately or kept permissive until lowering/validation code defines the exact JSON boundary.
- Review whether `FieldKind` and `FieldSemantic` names are the right split of former v1 field roles.

### What should be done in the future
- Implement P1.2 validation rules against these structs.
- Implement P1.3 lowering into current Widget IR maps.

### Code review instructions
- Start with `pkg/widgetdsl/v2/spec/types.go`.
- Validate with `go test ./pkg/widgetdsl/... -count=1`.

### Technical details
- Test command result: `ok github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl`; `pkg/widgetdsl/v2/spec` has no tests yet.

## Step 11: Add initial validation rules for typed v2 specs

I implemented the first validation layer for the v2 spec package. The validators check the invariants that the old map-IR DSL could miss or fail late: required page IDs, node kinds, section levels, schema fields, duplicate field names, multiple key fields, collection modes, arrangement kinds, selection params, action targets, server action names, payload field names, and template paths.

This is still a foundation step: the validators exist, but focused unit tests and lowering are separate tasks. The goal here was to make the diagnostic type concrete and establish where author-facing errors will come from.

### Prompt Context

**User prompt (verbatim):** (same as Step 7)

**Assistant interpretation:** Continue P1 by implementing validation rules for the new typed spec substrate.

**Inferred user intent:** Replace v1's typo-prone option maps with explicit validation that catches DSL misuse early and precisely.

**Commit (code):** pending — "Add widget DSL v2 spec validation"

### What I did
- Added `pkg/widgetdsl/v2/spec/validate.go`.
- Implemented `Validate()` methods for page, node, section, schema, field, collection, selection, action, payload template, template spec, and template value.
- Added `HasErrors` helper.
- Added diagnostic helpers for error/warning issues.
- Ran `gofmt -w pkg/widgetdsl/v2/spec/*.go`.
- Ran `go test ./pkg/widgetdsl/... -count=1`.
- Checked off task 17.
- Updated doc relations and changelog.

### Why
- Validation is where v2 earns its keep over v1. The known v1 failure modes are mostly missing/weak validation boundaries.

### What worked
- Existing package tests still pass, and the new validation package compiles cleanly.

### What didn't work
- No validation unit tests exist yet. That is intentional for this small commit, but P1.4 should add positive/negative fixtures before more layers depend on these rules.

### What I learned
- The validation rules expose a design detail: zero-value enums currently validate as errors. Builders should set explicit defaults before validation rather than relying on implicit zero values.

### What was tricky to build
- Keeping validation strict without making it too renderer-specific. For example, `ActionKindCopy` is currently permissive because its exact value/context semantics need to be finalized in Action IR v2.

### What warrants a second pair of eyes
- Review whether `SchemaSpec` should require exactly one key field, or allow zero and let lowering fall back to `id`.
- Review whether section levels should remain restricted to 1-3 as current `ui.section` documentation implies.

### What should be done in the future
- Add validation tests for typo'd arrangements, duplicate fields, invalid section levels, invalid server actions, and bad template paths.
- Implement lowering after the validation behavior is tested.

### Code review instructions
- Start with `pkg/widgetdsl/v2/spec/validate.go`.
- Validate with `go test ./pkg/widgetdsl/... -count=1`.

### Technical details
- Test command passed: `go test ./pkg/widgetdsl/... -count=1`.

## Step 12: Lower typed v2 specs to current Widget IR

I implemented the first lowering pass from typed v2 specs into the current JSON-like Widget IR maps. This makes the typed model executable in principle: pages lower to Widget page objects, nodes lower to text/element/component nodes, sections lower to `SectionBlock`, collections lower to composed `Stack` + `DataTable` + optional detail form trees, and actions lower to serializable action objects.

The lowering code intentionally targets the existing React renderer and current component names. This keeps the hard-cutover focused on replacing the authoring substrate first, not rewriting the renderer at the same time.

### Prompt Context

**User prompt (verbatim):** (same as Step 7)

**Assistant interpretation:** Continue P1 by implementing conversion from typed v2 specs to existing Widget IR nodes.

**Inferred user intent:** Make the v2 typed substrate practical by connecting it to the current renderer instead of designing types in isolation.

**Commit (code):** pending — "Lower widget DSL v2 specs to IR"

### What I did
- Added `pkg/widgetdsl/v2/spec/lower.go`.
- Implemented `PageSpec.ToWidgetPage()`.
- Implemented `NodeSpec.ToWidgetNode()`.
- Implemented `SectionSpec.ToNode()`.
- Implemented `CollectionSpec.ToNode()` for table and master-detail arrangements.
- Implemented table column derivation, create button lowering, selected detail form lowering, action column lowering, and action/template lowering.
- Ran `gofmt -w pkg/widgetdsl/v2/spec/*.go`.
- Ran `go test ./pkg/widgetdsl/... -count=1`.
- Checked off task 18.
- Updated doc relations and changelog.

### Why
- Lowering lets v2 builders target typed specs while still using the existing WidgetRenderer and components.

### What worked
- The package compiles and existing widgetdsl tests still pass.

### What didn't work
- First compile attempt failed because I accidentally referenced a non-existent `FieldSemanticDate` and then briefly duplicated `FieldSemanticShort` while fixing it:
  - Command: `go test ./pkg/widgetdsl/... -count=1`
  - Error: `pkg/widgetdsl/v2/spec/lower.go:386:106: duplicate case FieldSemanticShort ...`
- I fixed the switch to use the existing semantic constants and reran tests successfully.

### What I learned
- Date is currently modeled as a `FieldKind`, not a `FieldSemantic`. That split should stay clear in later builder APIs.

### What was tricky to build
- Action lowering has to bridge v2's typed template idea with the current frontend's string interpolation. I lowered confirm templates to the legacy string form for now and left path payloads as typed descriptors for future Action IR v2 consumers.

### What warrants a second pair of eyes
- Review whether collection lowering should generate full-page URLs or relative query-string URLs for row selection. Current v1 uses `?param=value`; the target v2 examples sometimes show full `/pages/...` URLs.
- Review whether path payload descriptors should be disallowed until Action IR v2 hydration exists.

### What should be done in the future
- Add focused P1.4 tests for simple table, selectable table, and master-detail lowering outputs.

### Code review instructions
- Start with `pkg/widgetdsl/v2/spec/lower.go`.
- Validate with `go test ./pkg/widgetdsl/... -count=1`.

### Technical details
- Successful command: `cd rag-evaluation-system && go test ./pkg/widgetdsl/... -count=1`.

## Step 13: Add v2 spec tests for the foundational examples

I added tests for the exact foundational examples from the companion document: simplest table, selectable table, and master-detail editor. These tests validate the typed specs and inspect the lowered Widget IR shape so later builder work has a concrete target.

The tests also include a negative invalid-arrangement case to lock in the hard-cutover validation behavior: typo'd arrangement names should fail as diagnostics, not silently degrade or emit malformed maps.

### Prompt Context

**User prompt (verbatim):** (same as Step 7)

**Assistant interpretation:** Complete P1.4 by adding tests for the v2 spec examples before moving on to Goja builders.

**Inferred user intent:** Keep the implementation measurable against the examples and prevent regressions as the v2 API is built.

**Commit (code):** pending — "Test widget DSL v2 spec examples"

### What I did
- Added `pkg/widgetdsl/v2/spec/lower_test.go`.
- Tested simplest table lowering.
- Tested selectable table lowering and URL navigate action generation.
- Tested master-detail editor lowering with create/table/detail form children.
- Tested validation rejection for an invalid arrangement value.
- Ran `gofmt -w pkg/widgetdsl/v2/spec/*.go`.
- Ran `go test ./pkg/widgetdsl/... -count=1`.
- Checked off task 19.
- Updated doc relations and changelog.

### Why
- These tests turn the chat/document examples into executable fixtures.

### What worked
- Both `pkg/widgetdsl` and `pkg/widgetdsl/v2/spec` tests pass.

### What didn't work
- N/A.

### What I learned
- The typed spec can already express the first three example timelines without needing v2 Goja builders yet.

### What was tricky to build
- The test helpers need to inspect `JSONObject` / `[]JSONValue` assertions. This is acceptable for low-level lowering tests, but higher-level builder tests should compare more semantic outputs or snapshots to avoid brittle map assertions.

### What warrants a second pair of eyes
- Review whether the master-detail test should assert more of the generated form rows and delete action confirm template, or whether that belongs in Action IR v2 tests.

### What should be done in the future
- Start P2 builder substrate so JavaScript authors can produce these tested typed specs.

### Code review instructions
- Start with `pkg/widgetdsl/v2/spec/lower_test.go`.
- Validate with `go test ./pkg/widgetdsl/... -count=1`.

### Technical details
- Successful command: `cd rag-evaluation-system && go test ./pkg/widgetdsl/... -count=1`.

## Step 14: Implement initial Goja v2 builders for simple and selectable tables

I implemented the first JavaScript-facing builder layer as an experimental `data.v2.dsl` module. It exposes the typed v2 substrate through fluent field/schema/collection APIs, uses non-enumerable Goja object properties as hidden typed refs, and rejects present non-function callbacks instead of silently ignoring them.

This step covers the first two runtime examples from the companion document: simplest table and selectable table. The builders produce the typed `spec.CollectionSpec` model and then lower it to the existing Widget IR, so the current React renderer can still render the output.

### Prompt Context

**User prompt (verbatim):** "continue."

**Assistant interpretation:** Continue the task-by-task implementation plan from P1 into P2, starting with Goja builders.

**Inferred user intent:** Keep making concrete implementation progress while preserving the diary/task/commit workflow.

**Commit (code):** pending — "Add widget DSL v2 Goja table builders"

### What I did
- Added `pkg/widgetdsl/v2_builders.go`.
- Registered an experimental `data.v2.dsl` native module in `pkg/widgetdsl/module.go`.
- Implemented hidden typed refs for field, schema, selection, action, and collection builder handles.
- Implemented `data.f.key/primary/short/prose/count/status` field builders.
- Implemented `data.schema(name).field(...).build()`.
- Implemented `data.collection(name, rows).schema(...).table().toIR()`.
- Implemented `select(s => s.urlParam(...))` and `table(t => t.rowSelect(...))` callback paths.
- Implemented strict callback handling for `table(callback)`.
- Added `pkg/widgetdsl/v2_builders_test.go` runtime tests.
- Updated lowering so explicit row-select actions override the default URL-selection navigate action.
- Ran `go test ./pkg/widgetdsl/... -count=1`.
- Checked off tasks 20, 21, and 22.

### Why
- P1 proved typed specs can model the examples. P2 begins making those specs usable from JavaScript without reintroducing v1 option bags.

### What worked
- The runtime tests verify simple table output, selectable table output, and strict rejection of a non-function callback.
- Existing widgetdsl and v2 spec tests still pass.

### What didn't work
- The first selectable table test failed because collection lowering always preferred the default selection navigate action over the explicit `.rowSelect(...)` action:
  - Command: `go test ./pkg/widgetdsl/... -count=1`
  - Failure: `onRowSelect.to = "?selected=${row.sessionId}"`
- I fixed `CollectionSpec.tableNode` so an explicit `Actions.Open`/row-select action wins when present, while selection still contributes `selectedKey`.

### What I learned
- Selection state and row-click behavior should be separate in the typed model. Selection tells the renderer which row is selected; row-select action tells it what to do on click. The default action can be derived only when no explicit action is supplied.

### What was tricky to build
- Hidden refs need to be attached to Goja objects without becoming public authoring API. I used a non-enumerable, non-writable, non-configurable internal property. That is sufficient for tests, though it is not a security boundary.

### What warrants a second pair of eyes
- Review whether the module should be named `data.v2.dsl` during implementation or whether it should replace `data.dsl` behind a branch before demos are public.
- Review whether hidden ref extraction should include stronger type assertions or source diagnostics.

### What should be done in the future
- Implement P2.4 master-detail editor builder with `edit(...).masterDetail()` and native form submit.
- Add demos once builder coverage reaches the examples we want to expose.

### Code review instructions
- Start with `pkg/widgetdsl/v2_builders.go`, then `pkg/widgetdsl/v2_builders_test.go`.
- Check `pkg/widgetdsl/v2/spec/lower.go` for row-select precedence.
- Validate with `go test ./pkg/widgetdsl/... -count=1`.

### Technical details
- Successful command: `cd rag-evaluation-system && go test ./pkg/widgetdsl/... -count=1`.

## Step 15: Add v2 master-detail editor builder and runtime test

I finished the first pass of P2 by adding master-detail editor authoring to `data.v2.dsl`. The v2 collection builder can now express all three foundational examples from the companion document: simple table, selectable table, and master-detail editor with URL selection and native form submit.

The master-detail API is still intentionally small, but it covers the essential shape: `.edit(e => e.selectUrl(...).submitPost(...).create(...).actions(...)).masterDetail().toIR()`. The implementation lowers to the same current Widget IR component tree used by the existing renderer: create button, `DataTable`, detail `FormPanel`, and close button.

### Prompt Context

**User prompt (verbatim):** (same as Step 14)

**Assistant interpretation:** Continue P2 implementation by adding the master-detail editor builder and corresponding runtime tests.

**Inferred user intent:** Get the v2 DSL to cover the concrete examples before moving on to demos or Action IR v2.

**Commit (code):** pending — "Add widget DSL v2 master-detail builder"

### What I did
- Added `collection.edit(callback)`.
- Added `collection.masterDetail()`.
- Added editor-builder methods: `selectUrl`, `submitPost`, `create`, `reorder`, `remove`, and `actions`.
- Added collection action-builder methods: `reorder` and `remove`.
- Added a Goja runtime test for the master-detail authoring example.
- Ran `go test ./pkg/widgetdsl/... -count=1`.
- Checked off tasks 23 and 24.
- Updated changelog and docmgr doctor.

### Why
- The master-detail editor is the most important current real-world grammar behavior. It needs to be expressible in v2 before demo pages or real page rewrites start.

### What worked
- The new runtime test verifies selected row state, form action, and detail form title in the lowered IR.
- All widgetdsl tests pass.

### What didn't work
- N/A.

### What I learned
- The current renderer contract is enough for the first v2 authoring layer. We can defer renderer changes until Action IR v2/context hydration work.

### What was tricky to build
- The editor builder has two ways to set actions: direct `.reorder/.remove` and nested `.actions(a => ...)`. The nested form matches the design-guide examples, while direct methods keep tests and simple scripts compact. This may need pruning after small-model trials.

### What warrants a second pair of eyes
- Decide whether both direct editor action methods and nested `.actions(...)` should remain in the final public API.
- Review whether `.create({ label })` should accept only an object, only a string, or both.

### What should be done in the future
- Start P3 Action IR v2 work, especially confirm centralization and DataTable action-cell row-key context.
- Start P4 demos now that the foundational v2 authoring examples run in Goja tests.

### Code review instructions
- Review `pkg/widgetdsl/v2_builders.go` editor builder methods.
- Review `TestDataV2BuilderBuildsMasterDetailEditor` in `pkg/widgetdsl/v2_builders_test.go`.
- Validate with `go test ./pkg/widgetdsl/... -count=1`.

### Technical details
- Successful command: `cd rag-evaluation-system && go test ./pkg/widgetdsl/... -count=1`.

## Step 16: Add live go-go-course demo pages for the foundational v2 examples

With the foundational `data.v2.dsl` builders in place, I added live demo pages to `go-go-course`. These pages are deliberately small and map directly to the companion document timelines: simplest table, selectable table, and master-detail editor. They are routed as normal Widget pages and use the current React renderer through Widget IR produced by the new v2 builders.

The master-detail demo uses a safe demo form route, `/settings/dsl-demo-agenda-item`, which redirects back with status but does not mutate course metadata. That keeps the example usable for repeated local testing.

### Prompt Context

**User prompt (verbatim):** (same as Step 14)

**Assistant interpretation:** Continue implementation by adding demo sites/examples now that the basic v2 builders are usable.

**Inferred user intent:** Ensure examples evolve alongside the implementation and provide concrete pages for manual/browser validation.

**Commit (code):** f82f20a99780902cb776e022ad6d1a3b3c2ee9a7 — "Add Widget DSL v2 demo pages"

### What I did
- Added `go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js`.
- Added `/pages/dsl-examples-table` for the simplest table example.
- Added `/pages/dsl-examples-selectable-table` for URL-backed selection.
- Added `/pages/dsl-examples-master-detail` for URL-backed master-detail plus native form submit.
- Wired page dispatch in `go-go-course/cmd/go-go-course/lib/course-pages.js`.
- Required and passed `data.v2.dsl` from `go-go-course/cmd/go-go-course/server.js`.
- Added a safe `POST /settings/dsl-demo-agenda-item` redirect route.
- Added a navigation entry labelled `DSL examples`.
- Ran `cd go-go-course && go test ./...`.
- Checked off tasks 29, 30, and 31.

### Why
- The user explicitly requested demo sites as we go and stale/deprecated examples to be managed. These live pages make the first v2 examples manually testable before deeper Action IR work.

### What worked
- The go-go-course Go test suite passes.
- The demo pages were committed in the go-go-course repo.

### What didn't work
- I initially added one doc relation with a malformed duplicated path while updating docmgr relations. I removed that bad frontmatter entry manually from design-doc 06.

### What I learned
- The v2 builders can be used directly from existing xgoja page files once `data.v2.dsl` is required and passed through the environment.

### What was tricky to build
- The master-detail demo needed to exercise native form submit without modifying real course data. The safe redirect-only route solves that for now, though a richer demo action page will still be needed for server action behavior.

### What warrants a second pair of eyes
- Review whether adding `DSL examples` to the main course material navigation is the right discoverability choice, or whether it should live in an internal/dev section.
- Manually run the server and click through the three demo pages in a browser.

### What should be done in the future
- Add the row/server-action demo page for P4.4 after P3 Action IR fixes land.
- Add a demo README for P4.5.

### Code review instructions
- Start with `go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js`.
- Review `go-go-course/cmd/go-go-course/server.js` for the `data.v2.dsl` require and safe demo form route.
- Validate with `cd go-go-course && go test ./...`; then manually browse `/pages/dsl-examples-table`, `/pages/dsl-examples-selectable-table`, and `/pages/dsl-examples-master-detail`.

### Technical details
- Successful command: `cd go-go-course && go test ./...`.

## Step 17: Start P3 Action IR foundations in the frontend

I started P3 by adding the frontend pieces needed by typed Action IR v2. The TypeScript IR now has data shapes for template prompts and payload templates; the browser action dispatcher can hydrate payload templates from action context; and app-level direct navigation paths now route through `dispatchWidgetAction` so confirmation handling is not bypassed when actions originate outside a normal `WidgetRenderer` adapter.

I also fixed the DataTable action-cell context issue found during the timeline write-up. Action button cells now receive the table's configured row-key spec instead of hard-coding `file`, so generic collection actions can send the correct `rowKey` for agenda/session rows.

### Prompt Context

**User prompt (verbatim):** "go ahead."

**Assistant interpretation:** Proceed with the next implementation phase, starting P3 Action IR v2 and context fixes.

**Inferred user intent:** Continue the planned cutover implementation without waiting for further confirmation.

**Commit (code):** pending — "Add frontend Action IR v2 foundations"

### What I did
- Added `TemplateSpec`, `TemplatePartSpec`, and `PayloadTemplateSpec` types to `packages/rag-evaluation-site/src/widgets/ir.ts`.
- Allowed server-action payloads to be either `JsonObject` or `PayloadTemplateSpec`.
- Added payload-template hydration in `packages/rag-evaluation-site/src/widgets/actions.ts`.
- Added typed template rendering for confirm prompts.
- Updated `App.tsx` server-action POSTs to send hydrated payloads.
- Updated direct AppNav/CourseStudioShell action dispatch to route through `dispatchWidgetAction` for centralized confirmation behavior.
- Updated `cellRenderers.tsx` so action-button cells compute `rowKey` from the table row-key spec.
- Updated `DataTable.widget.tsx` to pass `props.getRowKey` into `renderCell`.
- Ran `pnpm --dir packages/rag-evaluation-site typecheck`.
- Ran `go test ./pkg/widgetdsl/... -count=1`.
- Checked off tasks 25, 26, and 27.

### Why
- The v2 DSL must serialize browser-visible behavior as data, and row/server actions need reliable context. These changes make typed payload templates and correct row identity available before adding richer server-action demos.

### What worked
- TypeScript typecheck passed.
- Existing Go tests passed.

### What didn't work
- N/A. No compile errors in this step.

### What I learned
- Confirmation was already handled for normal WidgetRenderer-originated actions, but app-level direct action calls in navigation shells could bypass that central dispatcher. Routing those direct paths through `dispatchWidgetAction` closes the gap.

### What was tricky to build
- Payload templates need to support both the future explicit `kind: "payloadTemplate"` form and the current lowered object form where individual fields may be `{ kind: "path", path: "row.id" }`. I made `resolveActionPayload` hydrate template parts in either shape.

### What warrants a second pair of eyes
- Review whether confirm templates should be allowed as both strings and typed templates long-term, or whether strings should move to `widget.unsafe` after hard cutover.
- Review whether `resolveActionPayload` should reject non-JSON context values instead of stringifying them.

### What should be done in the future
- Add P3.4 tests for navigate, server payload hydration, confirm cancel, and refresh. The package currently has no Vitest/Jest setup, so this may require adding a small test harness or testing via Storybook/browser automation.
- Add the row/server-action demo page now that row key context and payload hydration are in place.

### Code review instructions
- Start with `packages/rag-evaluation-site/src/widgets/actions.ts` and `packages/rag-evaluation-site/src/widgets/ir.ts`.
- Then review `cellRenderers.tsx` and `DataTable.widget.tsx` for row context changes.
- Validate with `pnpm --dir packages/rag-evaluation-site typecheck` and `go test ./pkg/widgetdsl/... -count=1`.

### Technical details
- Successful commands:
  - `cd rag-evaluation-system && pnpm --dir packages/rag-evaluation-site typecheck`
  - `cd rag-evaluation-system && go test ./pkg/widgetdsl/... -count=1`

## Step 18: Add row/server-action demo and v2 action helper methods

I added the fourth live demo page: a row/server-action master-detail example. It uses the v2 collection builder with reorder and delete action columns, and the backend demo handlers return `refresh` and `toast` without changing persistent course data.

To make the demo expressive, I also added small v2 action-builder helpers for confirmation and payload fields. The delete action can now be authored as `dataV2.action.server("...").confirm("Delete demo agenda item “${row.title}”?")`, which lowers to the current confirm string behavior while keeping the authoring API fluent.

### Prompt Context

**User prompt (verbatim):** (same as Step 17)

**Assistant interpretation:** Continue P3/P4 implementation by adding a row/server-action demo once payload hydration and row context fixes are available.

**Inferred user intent:** Keep adding runnable examples as implementation proceeds so behavior can be validated manually and stale examples do not accumulate.

**Commit (code):** f7bac955bcbbbc1221636462deea65efbeec4fba — "Add Widget DSL v2 row action demo"

### What I did
- Added `confirm`, `payloadPath`, and `payload` methods to v2 action builder handles.
- Extended the v2 master-detail runtime test to assert the lowered delete confirm prompt.
- Added `/pages/dsl-examples-actions` in `go-go-course`.
- Added demo server action handlers for `dsl-demo-reorder-agenda` and `dsl-demo-delete-agenda`.
- Added the demo page route to `course-pages.js`.
- Updated the DSL examples navigation.
- Ran combined validation commands.
- Checked off task 32.

### Why
- Row actions are where Widget DSL behavior crosses from rendering into backend effects. The demo makes the row context, payload, refresh, and toast flow visible.

### What worked
- `go test ./pkg/widgetdsl/... -count=1` passed.
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `cd go-go-course && go test ./...` passed.
- The `go-go-course` row-action demo was committed separately.

### What didn't work
- N/A.

### What I learned
- The demo can exercise server-action dispatch safely without persistent mutation by returning synthetic refresh/toast results.

### What was tricky to build
- The confirm helper currently lowers to the existing string interpolation model. That is pragmatic for the current renderer, but the long-term v2 API should decide whether confirm strings are final or just a bridge until fully typed templates are exposed.

### What warrants a second pair of eyes
- Review whether `payload(name, value)` should exist as a chain method or whether payload configuration should only happen through a dedicated payload builder callback.
- Manually verify confirm dialog behavior and row context by clicking Delete/arrow buttons in `/pages/dsl-examples-actions`.

### What should be done in the future
- Add formal tests for action dispatch/hydration/confirm behavior. The frontend package currently lacks a JS test harness, so this likely needs a small test setup or Playwright-based smoke test.
- Add the demo README for P4.5.

### Code review instructions
- Review `pkg/widgetdsl/v2_builders.go` action helper methods.
- Review `go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js` and `server.js` demo action handlers.
- Validate with:
  - `cd rag-evaluation-system && go test ./pkg/widgetdsl/... -count=1`
  - `cd rag-evaluation-system && pnpm --dir packages/rag-evaluation-site typecheck`
  - `cd go-go-course && go test ./...`

### Technical details
- Successful commands:
  - `go test ./pkg/widgetdsl/... -count=1`
  - `pnpm --dir packages/rag-evaluation-site typecheck`
  - `go test ./...` in `go-go-course`

## Step 19: Clean up deprecated examples and perform real browser validation

The user asked whether I had tested the demos for real, including screenshots. At that point I had run unit/type/build checks and added live pages, but I had not yet browser-tested them. I corrected that in this step: I built the generated `go-go-course` runtime, ran it through the hotreload host, exercised the Widget page/action/form endpoints with `curl`, captured Playwright screenshots for all four v2 demo pages, and stored the screenshots under the ticket artifacts.

I also cleaned up public examples so new readers see `data.v2.dsl` first. The public provider docs now list `data.v2.dsl`, show v2 fluent examples for table/selectable/master-detail/actions, and explicitly label the old `data.dsl` option-bag grammar as legacy/current-runtime behavior rather than the recommended direction.

### Prompt Context

**User prompt (verbatim):** "cleanup deprecated examples. Did you already test it for real? like up to taking screenshots or so?"

**Assistant interpretation:** Remove ambiguity around deprecated examples and perform actual browser/demo validation with screenshots, not just unit tests.

**Inferred user intent:** Ensure the demos are actually usable in the app and that public documentation does not train humans or agents on deprecated APIs.

**Commit (code):** pending — "Document Widget DSL v2 docs cleanup and browser validation"

### What I did
- Added `data.v2.dsl` to the widgetsite provider module list.
- Added `data.v2.dsl` to `go-go-course/cmd/go-go-course/xgoja.package.yaml`.
- Rebuilt generated `go-go-course` runtime package with `make build`.
- Updated public widgetsite docs:
  - `01-widget-dsl-getting-started.md`
  - `02-widget-dsl-js-api-reference.md`
- Marked v1 `data.dsl` option-bag examples as legacy/current-runtime examples.
- Added v2 examples for simplest table, selectable table, master-detail editor, and row actions.
- Started local server with `go run ./hotreload-host -listen 127.0.0.1:8787`.
- Captured four screenshots with Playwright CLI and copied them to ticket artifacts.
- Ran API checks for page JSON, selected key, form title, demo server action, and native form redirect.
- Updated design-doc 06 with the cleanup and browser-validation evidence.
- Checked off task 33.
- Ran `docmgr doctor` successfully.

### Why
- Tests alone did not prove that the generated xgoja runtime could import `data.v2.dsl`, that the live demo pages rendered in a browser, or that public docs had stopped teaching deprecated option-bag forms as the main path.

### What worked
- After exposing `data.v2.dsl` through the provider and package YAML, `make build` succeeded.
- The hotreload host served the demo pages at `127.0.0.1:8787`.
- Playwright screenshots showed the intended simplest table, selectable table, master-detail editor, and row-action pages.
- API checks returned the expected selected key, form title, reorder action result, and native form redirect.

### What didn't work
- Initial `make build` failed because `server.js` imported `data.v2.dsl`, but `xgoja.package.yaml` did not select that module:
  - Command: `cd go-go-course && make build`
  - Error: `Error: server.js imports unknown bare specifier "data.v2.dsl"`
- I fixed this by adding `data.v2.dsl` to both the widgetsite provider module list and the go-go-course xgoja package selected modules.
- I also tried the stale Makefile/comment command shape with `--http-listen`, but the generated runtime rejected it:
  - `go-go-course: unknown flag: --http-listen`
- The working local server command was `go run ./hotreload-host -listen 127.0.0.1:8787` from `go-go-course/cmd/go-go-course`.

### What I learned
- Adding a native module to `widgetdsl.Register` is not enough for generated xgoja binaries; the provider must advertise it and each generated app must select it in `xgoja.package.yaml`.
- The current README/Makefile command comments around `--http-listen` appear stale for the generated binary path; hotreload-host is the reliable local validation path here.

### What was tricky to build
- The browser tooling harness was initially unavailable because a Playwright MCP profile was locked by a running process. I used the Playwright CLI for screenshots instead and stored the resulting PNGs as ticket artifacts.
- I also had to distinguish cleanup of public examples from deletion of historical research docs: old RAGEVAL-UI-GRAMMAR ticket docs intentionally keep v1 examples as historical evidence, while public provider docs now direct new work to v2.

### What warrants a second pair of eyes
- Review whether `data.v2.dsl` should be listed as experimental in public provider docs or whether it should replace `data.dsl` naming before broader publication.
- Review whether the stale `--http-listen` comments/Makefile help should be cleaned up in a separate go-go-course task.

### What should be done in the future
- Add a durable automated browser smoke test rather than relying on ad-hoc Playwright CLI screenshots.
- Finish P3.4 formal action tests once a JS/Playwright test harness decision is made.

### Code review instructions
- Review public docs in `pkg/xgoja/providers/widgetsite/doc/` for v2-first examples and legacy labeling.
- Review `pkg/xgoja/providers/widgetsite/provider.go` and `go-go-course/cmd/go-go-course/xgoja.package.yaml` for `data.v2.dsl` exposure.
- Inspect screenshot artifacts under `ttmp/.../artifacts/dsl-demo-screenshots/`.
- Reproduce with `go run ./hotreload-host -listen 127.0.0.1:8787` and Playwright screenshot commands listed in design-doc 06.

### Technical details
- Screenshot files:
  - `artifacts/dsl-demo-screenshots/table.png`
  - `artifacts/dsl-demo-screenshots/selectable.png`
  - `artifacts/dsl-demo-screenshots/master-detail.png`
  - `artifacts/dsl-demo-screenshots/actions.png`

## Step 20: Fix form submit visibility and clarify New demo agenda behavior

The user reported that form submit buttons looked like tiny black squares and asked whether “New demo agenda item” was supposed to work. I reproduced the problem in screenshots: the button box was rendered, but the `Save` label was missing. The root cause was `FormPanel.widget.tsx` passing `null` for an absent `submitLabel`; React default parameters only apply to `undefined`, so `FormPanel` received `submitLabel = null` and rendered an empty primary button.

I fixed the label fallback inside `FormPanel`, made button sizing more robust, and clarified the New demo flow. The New button is supposed to open a blank demo form; it is not a persistent create operation. The demo now says that explicitly, and v2 lowering makes the key field editable in `__new` mode so the blank form is less confusing.

### Prompt Context

**User prompt (verbatim):** "the submit button on the forms is like, a tiny black square that's barely visible. 

Also not sure if the "New demo agenda item" is uspposed to work or not"

**Assistant interpretation:** Fix the visual defect on form submit buttons and clarify or repair the New demo item behavior.

**Inferred user intent:** Make the demos usable and self-explanatory before treating them as validation/examples.

**Commit (code):** 724f3d57 — "Fix Widget DSL demo form controls"

### What I did
- Fixed `FormPanel` to use `submitLabel ?? "Save"`, so absent labels still render visibly even when the widget adapter passes `null`.
- Changed form submit buttons to use normal button size.
- Added minimum inline/block sizes and inline-flex centering to button CSS.
- Changed v2 key field defaults so existing rows remain read-only through lowering logic, but new `__new` master-detail forms can edit the key field.
- Added a regression test for editable key fields in new master-detail mode.
- Added explanatory captions to the master-detail/action demos when `agenda=__new`.
- Rebuilt the rag frontend package and go-go-course embedded SPA assets from the local rag source.
- Re-ran browser screenshots and verified the Save button is visible and the new-item caption appears.

### Why
- Empty primary buttons make the form look broken. The New demo behavior also needed to be explicit because the demo route intentionally does not persist new agenda rows.

### What worked
- Visual QA confirmed that the button is now labelled “Save” and the new-item screenshot shows a blank editable form plus explanatory caption.
- `go test ./pkg/widgetdsl/... -count=1` passed.
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `pnpm --dir packages/rag-evaluation-site build` passed.
- `cd go-go-course && go test ./...` passed.

### What didn't work
- The first screenshot after only CSS changes still showed a black square. That revealed the true bug was not only size/contrast; the label was null because default props did not apply to null.

### What I learned
- Widget adapter `ctx.renderValue(undefined)` returns `null`, so React component defaults should use nullish fallbacks internally for optional renderable labels.

### What was tricky to build
- New-item mode needed both UI clarification and model behavior. If the key field is always read-only, a blank new form is confusing. The lowering code now treats `__new` as an editable-key context unless a field explicitly sets readOnly.

### What warrants a second pair of eyes
- Check whether all optional renderable props in other widgets should also use nullish fallbacks rather than relying on default parameters.
- Decide whether the demo should eventually persist rows in memory for the session, or remain a non-persistent interaction demo.

### What should be done in the future
- Add automated browser smoke tests for the Save label and New-item behavior.

### Code review instructions
- Start with `FormPanel.tsx` and `Button.module.css`.
- Then review `pkg/widgetdsl/v2/spec/lower.go` and `lower_test.go` for new-item key editability.
- Finally review `go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js` for the explanatory caption.

### Technical details
- Updated screenshots:
  - `artifacts/dsl-demo-screenshots/master-detail.png`
  - `artifacts/dsl-demo-screenshots/master-detail-new.png`

## Step 21: Add durable browser smoke tests for Widget DSL demos and actions

I turned the prior manual browser validation into a repeatable Playwright test harness in `go-go-course`. The new smoke suite starts the `hotreload-host`, loads the four Widget DSL v2 demo pages, verifies selection URL behavior, native form save redirects, new-item editability, row action POST/refresh behavior, and both cancel and accept branches of the delete confirmation.

This closes the earlier gap where screenshots and `curl` checks proved the demos once, but did not provide a durable regression test for future renderer or DSL changes. The test lives with the generated course SPA so it validates the real embedded app path rather than only backend JSON fixtures.

### Prompt Context

**User prompt (verbatim):** "continue, build all remaining tasks, and create examples for all DSL areas, including DSLs, CMS, UI, context, etc."

**Assistant interpretation:** Continue the remaining GOJA-DSL-PLAYBOOK implementation tasks, starting with the unfinished action/browser validation work, and then expand coverage/examples across the other DSL module families.

**Inferred user intent:** Move from design and one-off demos toward a complete, validated implementation with durable examples and tests across the full DSL surface.

**Commit (code):** 06aa1c9 — "Add Widget DSL demo browser smoke tests"

### What I did
- Added `@playwright/test` to `go-go-course/cmd/go-go-course/webapp` and a `test:dsl-examples` script.
- Added `playwright.config.ts` that starts `cd .. && go run ./hotreload-host -listen 127.0.0.1:18787` and waits for `/api/widget/health`.
- Added `tests/dsl-examples.spec.ts` covering:
  - simple table render/no-navigation behavior,
  - selectable table row click and selected query param,
  - master-detail selection and native Save redirect,
  - `__new` item behavior and editable `id` field,
  - reorder action POST and refresh result,
  - delete confirmation cancel with no POST,
  - delete confirmation accept with server POST and refresh result.
- Ran `pnpm test:dsl-examples`; fixed selector assumptions; reran until all 7 tests passed.
- Ran `pnpm typecheck` for the webapp.
- Committed the test harness in `go-go-course`.
- Checked off task 28 and updated the ticket changelog/file relationships.

### Why
- P3.4 required action tests for navigate/reorder/delete/confirm cancel/refresh. The existing manual screenshots were useful evidence but not sufficient as a regression harness.
- The test should exercise the same app shell and hotreload-host route path a human uses, because previous failures involved integration boundaries such as module selection, embedded assets, and browser behavior.

### What worked
- Final command succeeded:
  - `cd go-go-course/cmd/go-go-course/webapp && pnpm test:dsl-examples`
  - Result: `7 passed (9.9s)`
- TypeScript validation succeeded:
  - `cd go-go-course/cmd/go-go-course/webapp && pnpm typecheck`
- The confirm-cancel test confirmed that dismissing the browser dialog prevents the delete POST.
- The confirm-accept test confirmed that accepting the dialog posts to `/api/widget/actions/dsl-demo-delete-agenda` and returns `{ ok: true, refresh: true }`.

### What didn't work
- The first Playwright run failed 5 of 7 tests because I assumed several section titles were semantic headings. The renderer exposes them as plain text, not heading roles.
- Example failure:
  - Command: `cd go-go-course/cmd/go-go-course/webapp && pnpm test:dsl-examples`
  - Error: `Locator: getByRole('heading', { name: 'Simplest table' }) Expected: visible Timeout: 5000ms Error: element(s) not found`
- I fixed the tests by using visible text locators for plain-text section titles and direct field selectors where the current form renderer does not create labelled textbox roles.
- A second run had one remaining failure because `getByText('New item ID')` did not match text split across elements. I changed that assertion to `getByText('New item', { exact: false })`.

### What I learned
- The current renderer does not always emit accessible heading/label semantics for section titles and form fields. Tests can still validate behavior, but this is also evidence that UI accessibility should be improved separately.
- Playwright webServer can reliably start the hotreload-host, which is better than relying on the stale generated binary flags that previously failed with `unknown flag: --http-listen`.

### What was tricky to build
- The browser-visible page content is produced by the Widget renderer, not by static HTML. That meant the test assertions needed to follow actual rendered accessibility output rather than the conceptual component names in the DSL examples.
- Testing confirm cancel required proving a negative: after dismissing the dialog, no POST should happen. I used a request counter and a short wait after dismissal; this is sufficient for the current synchronous action path but should be revisited if actions become queued/asynchronous.

### What warrants a second pair of eyes
- The tests use text locators for section titles because they are not headings. A reviewer should decide whether the renderer should emit better semantic headings and labels, then update tests accordingly.
- The negative assertion for confirm cancel uses a 300ms wait. It is pragmatic, but not as strong as an explicit action dispatcher test would be.

### What should be done in the future
- Add renderer accessibility improvements for headings and form label associations.
- Consider adding lower-level frontend tests for `dispatchWidgetAction` once the frontend package has a stable JS test harness.
- Continue with P5 TypeScript declaration generation/parity and the broader DSL family examples requested by the user.

### Code review instructions
- Start with `go-go-course/cmd/go-go-course/webapp/tests/dsl-examples.spec.ts`.
- Review `go-go-course/cmd/go-go-course/webapp/playwright.config.ts` for the `hotreload-host` command and health URL.
- Validate with:
  - `cd go-go-course/cmd/go-go-course/webapp && pnpm test:dsl-examples`
  - `cd go-go-course/cmd/go-go-course/webapp && pnpm typecheck`

### Technical details
- Test command: `pnpm test:dsl-examples`
- Passing result: `7 passed (9.9s)`
- Code commit: `06aa1c9ca4e49f986475007c2221767b9aaadc60`
- Ticket task checked: `[28] P3.4 Add action tests for navigate, reorder, delete, confirm cancel, and refresh`

## Step 22: Add precise TypeScript declarations for data.v2.dsl

I implemented the first TypeScript hard-cutover task by giving `data.v2.dsl` an explicit typed/fluent declaration surface instead of inheriting only the loose Widget IR helper declarations. The generated declaration now names opaque handles and builder interfaces for fields, schemas, selections, actions, tables, editors, and collections.

The important behavior change is that the v2 declaration surface teaches `data.collection(name, rows).schema(schema).select(...).edit(...).table(...).toIR()` and does not expose the legacy `dataTable`, `cell`, `record`, or option-bag `collection(rows, options)` grammar. This keeps the public type surface aligned with the hard-cutover direction.

### Prompt Context

**User prompt (verbatim):** (same as Step 21)

**Assistant interpretation:** Continue the remaining GOJA-DSL-PLAYBOOK implementation tasks after browser smoke validation, starting with P5 TypeScript declaration work.

**Inferred user intent:** Make the new DSL discoverable and type-checkable so humans and agents are guided toward v2 fluent builders instead of legacy option bags.

**Commit (code):** dcd5156 — "Add Widget DSL v2 TypeScript declarations"

### What I did
- Added `dataV2TypeScriptLines()` in `pkg/widgetdsl/typescript.go`.
- Declared opaque-ish branded handle interfaces for `FieldHandle`, `SchemaHandle`, `SelectionHandle`, and `ActionHandle`.
- Declared fluent interfaces for `FieldBuilder`, `SchemaBuilder`, `SelectionBuilder`, `ActionBuilder`, `TableBuilder`, `EditorBuilder`, `CollectionActionsBuilder`, and `CollectionBuilder`.
- Declared `f`, `schema`, `collection`, `selection`, and `action` exports for `data.v2.dsl`.
- Added `pkg/widgetdsl/typescript_test.go` with tests that assert the v2 builder declarations exist and legacy option-bag declarations are absent from `data.v2.dsl`.
- Updated the widgetsite provider test so `data.v2.dsl` is included in TypeScript descriptor/export coverage.
- Ran `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite -count=1`.
- Checked off task 34.

### Why
- Without precise declarations, v2 remains mostly a runtime experiment. The ticket goal explicitly includes compile-time help from generated types, and P5.1 called for precise v2 TypeScript declarations.
- The declarations also serve as API documentation for small models: the types constrain the sequence of builder handles and make the absence of legacy v1 API visible.

### What worked
- The declaration unit tests passed and guard against accidentally exposing the old `data.dsl` option-bag helpers through `data.v2.dsl`.
- Provider registration still resolves `data.v2.dsl` and now tests its TypeScript descriptor in the widgetsite package.

### What didn't work
- N/A for the implementation step. I did not yet add external TypeScript positive/negative fixture compilation; that belongs to P5.2.

### What I learned
- The current `TypeScriptModule` generator is line-oriented `RawDTS`, which is simple enough for the v2 declaration but has no structural parity test yet. The next step should compile fixture TypeScript against the emitted module declarations.

### What was tricky to build
- The runtime uses hidden Goja refs, while TypeScript needs a way to represent handles that users cannot freely construct. I used a `unique symbol` brand in interfaces to make handles nominal-ish in the declaration surface.
- The v2 module still inherits generic `text`, `element`, `component`, and `fragment` helpers from the shared module installer. The declaration keeps those generic helpers while adding the v2-specific builder surface; a future hard cutover may decide whether `data.v2.dsl` should be purely data-only.

### What warrants a second pair of eyes
- Review whether branded handles should be stricter or whether the current exported interfaces are sufficient for go-go-goja generated declarations.
- Review whether `payload(name, value: JsonValue)` is too narrow, since the runtime currently accepts `value.Export()` for arbitrary Goja values.
- Confirm whether `select(selection?: SelectionHandle | SelectionCallback | null)` should also type the callback as accepting and returning `SelectionBuilder`; runtime accepts the builder if the callback returns nothing after mutating it.

### What should be done in the future
- Complete P5.2: add runtime export parity plus TypeScript positive/negative fixtures.
- Consider a generated declaration snapshot test once the xgoja declaration bundling path is selected.

### Code review instructions
- Start with `pkg/widgetdsl/typescript.go`, especially `dataV2TypeScriptLines()`.
- Then review `pkg/widgetdsl/typescript_test.go` for the desired v2 API surface and legacy exclusions.
- Validate with:
  - `cd rag-evaluation-system && go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite -count=1`

### Technical details
- Successful command: `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite -count=1`
- Code commit: `dcd5156e0633c80669e45103fdb5b24e51042b75`
- Ticket task checked: `[34] P5.1 Generate precise v2 TypeScript declarations`

## Step 23: Add runtime export parity and TypeScript positive/negative fixtures

I completed the second TypeScript validation task by adding a Go test that renders the generated `data.v2.dsl` declaration bundle, writes a real TypeScript fixture, and compiles it with the repository's frontend TypeScript compiler. The fixture includes positive fluent-builder examples and negative examples guarded by `@ts-expect-error`, so the test fails both when valid v2 code stops compiling and when legacy/raw option-bag misuse accidentally becomes accepted.

I also added a runtime export parity check for the declared v2 public surface. The test verifies that `f`, `schema`, `collection`, `selection`, and `action` exist at runtime while legacy v1 exports such as `dataTable` and `cell` remain absent from `data.v2.dsl`.

### Prompt Context

**User prompt (verbatim):** (same as Step 21)

**Assistant interpretation:** Continue P5 after adding declarations by validating that the runtime export surface and TypeScript compiler behavior match the intended hard-cutover API.

**Inferred user intent:** Prevent divergence between generated `.d.ts` files, Goja runtime exports, and public examples before expanding the DSL surface further.

**Commit (code):** cee7525 — "Test Widget DSL v2 TypeScript fixtures"

### What I did
- Added `pkg/widgetdsl/typescript_fixture_test.go`.
- Rendered the `data.v2.dsl` declaration through `tsgen/render.Bundle` in a test.
- Wrote a temp `widgetdsl.d.ts` plus `data-v2-fixture.ts`.
- Compiled the fixture with `packages/rag-evaluation-site/node_modules/.bin/tsc` under `--strict --noEmit`.
- Included positive examples for schema fields, selectable table, master-detail edit, action payload paths, confirm, and `toIR()`.
- Included expected-negative examples for:
  - `data.dataTable(...)`,
  - missing `schema(name)`,
  - raw object passed to `schema.field`,
  - non-object collection rows,
  - raw JSON action passed to `table.rowSelect`.
- Added a runtime export parity test for top-level `data.v2.dsl` exports and legacy-export absence.
- Ran `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite -count=1`.
- Checked off task 35.

### Why
- Declaration strings alone are easy to accidentally make syntactically invalid or structurally too loose. Compiling real TypeScript is the cheapest way to prove the declarations guide valid code and reject invalid code.
- Runtime export parity prevents public declarations from documenting functions that the Goja module does not actually export.

### What worked
- The fixture compiled successfully with TypeScript, including the `@ts-expect-error` assertions.
- The runtime parity test confirmed the v2 module exports the typed/fluent surface and does not expose v1 `dataTable`/`cell` helpers.

### What didn't work
- N/A. The first fixture/test design passed after adding the local TypeScript compiler lookup.

### What I learned
- A Go test can exercise the TypeScript declaration pipeline without committing generated `.d.ts` files: render the module declaration into a temp dir, write fixtures, then invoke `tsc`.
- `@ts-expect-error` is useful here because it turns negative examples into durable assertions rather than comments.

### What was tricky to build
- The TypeScript compiler path is not at the repository root; it currently lives under `packages/rag-evaluation-site/node_modules/.bin/tsc`. The test locates the repo root and uses that compiler, with `WIDGETDSL_TSC` as an override.
- The test skips if the compiler is not installed, which keeps Go-only environments from failing before `pnpm install`, but means CI should ensure frontend dependencies are installed for full validation.

### What warrants a second pair of eyes
- Decide whether the TypeScript fixture should live as a committed `.ts` file under `testdata/` instead of an embedded string in Go.
- Decide whether the test should fail hard instead of skipping when TypeScript is unavailable once CI has a guaranteed Node setup.

### What should be done in the future
- Add the TypeScript fixture command to CI once the repo's Node dependency installation step is standardized.
- Continue with P6 migrations: rewrite real admin/session examples to v2 APIs.

### Code review instructions
- Review `pkg/widgetdsl/typescript_fixture_test.go` first; it defines the parity and compile contract.
- Validate with:
  - `cd rag-evaluation-system && go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite -count=1`

### Technical details
- Successful command: `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite -count=1`
- Code commit: `cee7525306500256690955deb340660bc0492c0d`
- Ticket task checked: `[35] P5.2 Add runtime export parity and TypeScript positive/negative fixtures`

## Step 24: Add a cross-module Widget DSL gallery demo

I added the first broader example requested by the user: a single live gallery page that imports and renders every public Widget DSL family currently used by `go-go-course`. The page now covers `ui.dsl`, `data.v2.dsl`, `context_window.dsl`, `cms.dsl`, and `course.dsl` in one place instead of only demonstrating the data-table family.

This is not a replacement for deeper per-domain examples, but it creates a discoverable browser-visible contract for the full module set and extends the Playwright smoke suite so regressions in any of these public DSL modules are caught when the demo page fails to render.

### Prompt Context

**User prompt (verbatim):** (same as Step 21)

**Assistant interpretation:** Start fulfilling the expanded request to make examples for all DSL areas, not only the new data v2 builders.

**Inferred user intent:** Ensure every DSL family has at least one current, runnable example so future humans and agents do not have to infer usage from old v1 or internal code.

**Commit (code):** 30ab8c3 — "Add Widget DSL module gallery demo"

### What I did
- Extended `go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js` with `/pages/dsl-examples-modules` content.
- Added navigation tab text `All DSL modules` to the DSL examples page.
- Added a `ui.dsl` primitive panel with status text and buttons.
- Added a `data.v2.dsl` typed collection table.
- Added a `context_window.dsl` `contextDiagramPanel` with a small snapshot and style set.
- Added a `cms.dsl` media-library recipe example using an inline data-URI demo asset to avoid missing asset requests.
- Added a `course.dsl` mini course shell and slide panel.
- Added routing for `dsl-examples-modules` in `course-pages.js`.
- Extended the Playwright smoke test with a module-gallery test that asserts all five DSL-family sections render.
- Ran `pnpm test:dsl-examples` and `go test ./...` in `go-go-course`.
- Added and checked off task 41.

### Why
- The existing demo pages validated `data.v2.dsl`, but the user's expanded request explicitly called for examples across DSLs, CMS, UI, context, and related modules.
- A single gallery page is the smallest useful cross-module fixture: it proves the modules can be imported together and gives readers a concrete starting point for each family.

### What worked
- Final Playwright result: `8 passed (4.0s)`.
- `cd go-go-course && go test ./...` passed.
- The module gallery renders through the real hotreload-host path and is covered by the same browser smoke harness as the data v2 demos.

### What didn't work
- The first version used `/course-assets/hero.png` for the CMS asset example. The page rendered, but the server logged a 404 for that image request during Playwright validation.
- I fixed that by using a tiny inline data-URI image for the demo asset so the example stays self-contained.

### What I learned
- Cross-module examples should avoid external or file-backed assets unless the test fixture also creates those assets. Self-contained examples are much more reliable for smoke tests.
- The existing DSL examples navigation can naturally scale from data-v2 pages to broader module pages without adding a separate docs page yet.

### What was tricky to build
- The `course.dsl` example is rendered inside the existing course shell, which can create a nested-shell layout. I kept it as a mini shell example because the goal is API coverage, not final page layout polish.
- The `cms.dsl` recipe wants realistic asset metadata. Using a data URI keeps the browser path clean but may not represent real upload/media serving behavior.

### What warrants a second pair of eyes
- Review whether the cross-module gallery should stay as one page or split into separate `/pages/dsl-examples-ui`, `/pages/dsl-examples-cms`, `/pages/dsl-examples-course`, and `/pages/dsl-examples-context-window` pages as examples grow.
- Review the nested `courseStudioShell` visual layout in the module gallery.

### What should be done in the future
- Add deeper domain-specific examples for CMS editing flows, course slide/handout flows, and context-window transcript/annotation flows.
- Keep the module gallery current as public DSL exports change.

### Code review instructions
- Start with `go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js`, especially `moduleGallerySection()`.
- Review `go-go-course/cmd/go-go-course/webapp/tests/dsl-examples.spec.ts` for the new browser smoke assertion.
- Validate with:
  - `cd go-go-course/cmd/go-go-course/webapp && pnpm test:dsl-examples`
  - `cd go-go-course && go test ./...`

### Technical details
- Successful browser command: `pnpm test:dsl-examples` → `8 passed (4.0s)`
- Successful Go command: `go test ./...` in `go-go-course`
- Code commit: `30ab8c3c8257e3cc8bf0e852c886397385b2a8ec`
- Ticket task added and checked: `[41] P4.6 Add cross-module demo page covering ui.dsl, data.v2.dsl, cms.dsl, course.dsl, and context_window.dsl`

## Step 25: Port the real admin agenda editor to data.v2.dsl

With the `origin/main` merge completed, I resumed the hard-cutover sequence on the first real page still using the previous `data.dsl` grammar. The admin Course CMS agenda editor now uses the typed/fluent `data.v2.dsl` collection API instead of the v1 option-bag `dataDsl.collection(...)` helper, while preserving the merged active-site wording and the existing server endpoints for save, reorder, and delete.

This step intentionally leaves the media library and session browse migrations for the next P6 task. Keeping the commit narrow made it possible to validate the real page route, generated xgoja mirror, and smoke suite without mixing agenda migration risk with CMS recipe or session-table changes.

### Prompt Context

**User prompt (verbatim):** "go ahead, keep detailed diary as you work, commit at appropriate intervals."

**Assistant interpretation:** Start the next post-merge cutover task, make focused commits, validate as work lands, and keep the GOJA-DSL-PLAYBOOK diary current.

**Inferred user intent:** Continue the Widget DSL v2 hard cutover now that the go-go-course merge blocker is resolved, with enough documentation to review or resume safely.

**Commit (code):** 34bc642 — "Port admin agenda editor to data.v2 DSL"

### What I did
- Confirmed `go-go-course` was on `task/improve-rag-evaluation-system`, clean, and fully merged with `origin/main`.
- Inspected the current admin page and verified it still used v1 `dataDsl.schema`, `dataDsl.f.*`, `dataDsl.collection`, `dataDsl.urlParam`, and `dataDsl.formPost` for the agenda editor.
- Rewrote `agendaSchema()` in `cmd/go-go-course/lib/pages/admin-course-cms.js` to use `dataV2.schema("AgendaItem")` with ordered fluent `.field(...)` calls.
- Replaced `dataDsl.collection(agenda, { ... })` with `dataV2.collection("agenda", agenda).schema(...).edit(...).masterDetail().toIR()`.
- Preserved the existing `/settings/agenda-item`, `admin-reorder-course-agenda`, and `admin-delete-agenda-item` contracts.
- Moved the save/error/idle message outside the v2 collection as a `ui.caption(...)` because the current v2 collection builder does not yet expose FormPanel status props.
- Ran `make -C cmd/go-go-course build` so the generated xgoja mirror under `internal/xgojaruntime/.../admin-course-cms.js` matched the source page.
- Ran broad and targeted validation, then committed the code change as `34bc642`.
- Checked ticket task `[36] P6.1 Rewrite admin agenda editor to v2 master-detail API` and updated changelog/relations.

### Why
- P6.1 is the first real-page proof that `data.v2.dsl` can replace the previous map/option-bag grammar outside dedicated demo pages.
- The merge resolved content-directory conflicts but deliberately kept the older admin page grammar; this step moves the agenda editor forward without disturbing media and session work that belong to P6.2.

### What worked
- `node --check` passed for the edited page and all source JS files.
- `GOWORK=off go test ./...` passed in `go-go-course`.
- `pnpm --dir cmd/go-go-course/webapp typecheck` passed.
- `GOWORK=off make -C cmd/go-go-course build XGOJA_LOCAL_DIR=/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-goja` passed and refreshed the generated mirror.
- The existing smoke suite passed on alternate port `18789`: `11 passed, 0 failed`.
- A targeted admin Widget IR smoke passed on port `18791` after setting the profile display name to `admin_p6`; it verified the page id, Agenda section, DataTable, FormPanel, New agenda item button, delete/reorder actions, success message, and selected edit title.

### What didn't work
- My first targeted admin smoke used the wrong profile endpoint, `POST /api/user/profile`, and failed with curl exit 22 because the server returned 404:

```text
method=POST path=/api/user/profile ... status=404
```

- I inspected `server.js`, found the correct endpoint is `PUT /api/user-session`, and reran the targeted smoke successfully.

### What I learned
- The real admin page can consume `data.v2.dsl` without server changes because the v2 action cells still dispatch row context, and the existing handlers already read `ctx.body.context.row` for agenda delete/reorder.
- The current v2 builder surface is sufficient for the agenda master-detail interaction, but it does not yet expose every v1 convenience (`empty`, FormPanel `status/statusMessage`, field `hint/placeholder`). For P6.1 I preserved status messaging with an outer `ui.caption` and kept the empty message as a conditional caption.

### What was tricky to build
- The biggest sharp edge was preserving behavior while moving from v1's option bag to v2's narrower typed builder surface. The source v1 call configured schema, selection, submit, create, reorder, delete, empty-state text, and status text in one object. The v2 chain cleanly covers schema/selection/submit/create/actions/layout, but status and empty text needed to be represented alongside the collection rather than inside it.
- The generated xgoja mirror is easy to forget because the source JS lives under `cmd/go-go-course/lib/...`, while runtime embeds a synchronized copy under `cmd/go-go-course/internal/xgojaruntime/...`. Running `make -C cmd/go-go-course build` after source edits kept those in sync.

### What warrants a second pair of eyes
- Review whether `data.v2.dsl` should grow first-class `.empty(...)`, `.status(...)`, `.placeholder(...)`, and `.hint(...)` builder methods before more real-page migrations, or whether status/empty presentation should remain explicit `ui.*` composition around typed collections.
- Review the visual placement of the outer agenda status caption versus the old FormPanel-integrated status message.
- Review the generated xgoja mirror to confirm it matches the source page and contains no stale v1 agenda grammar.

### What should be done in the future
- Continue with P6.2: migrate the remaining real media-library/session-browse v1 usages to v2 APIs.
- Consider adding a permanent admin-course-cms API smoke test so regressions in the real admin page are caught without manually starting the server.
- Consider making `cmd/go-go-course/test-fixtures/smoke-test.sh` accept `PORT` from the environment, as already noted by the merge diary.

### Code review instructions
- Start with `go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js`, specifically `agendaSchema()` and `agendaSection()`.
- Then compare the generated mirror at `go-go-course/cmd/go-go-course/internal/xgojaruntime/xgoja_embed/jsverbs/minitrace_viz_site/lib/pages/admin-course-cms.js`.
- Validate with:
  - `cd go-go-course && GOWORK=off go test ./...`
  - `cd go-go-course && find cmd/go-go-course/lib -name '*.js' -print0 | xargs -0 -n1 node --check`
  - `cd go-go-course && pnpm --dir cmd/go-go-course/webapp typecheck`
  - `cd go-go-course && GOWORK=off make -C cmd/go-go-course build XGOJA_LOCAL_DIR=/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-goja`
  - `cd go-go-course/cmd/go-go-course && cp test-fixtures/smoke-test.sh /tmp/gocourse-smoke-18789.sh && perl -0pi -e 's/PORT=18787/PORT=18789/' /tmp/gocourse-smoke-18789.sh && bash /tmp/gocourse-smoke-18789.sh`

### Technical details
- Code commit: `34bc642e452311939f17f0dc086acd61c9811967`
- Successful targeted admin smoke used:

```bash
cd go-go-course/cmd/go-go-course
./dist/go-go-course serve site start --http-listen 127.0.0.1:18791
curl -sf -c /tmp/gocourse-admin-p6-cookie.txt -b /tmp/gocourse-admin-p6-cookie.txt \
  -H 'Content-Type: application/json' \
  -X PUT http://127.0.0.1:18791/api/user-session \
  -d '{"displayName":"admin_p6"}'
curl -sf -c /tmp/gocourse-admin-p6-cookie.txt -b /tmp/gocourse-admin-p6-cookie.txt \
  'http://127.0.0.1:18791/api/widget/pages/admin-course-cms?agenda=agenda-regenerate&status=agenda-item-saved'
```
- Targeted smoke assertions passed for page id, Agenda section, DataTable, FormPanel, New agenda item, delete action, reorder action, save status, and selected edit title.
- Ticket task checked: `[36] P6.1 Rewrite admin agenda editor to v2 master-detail API`

## Step 26: Extend data.v2.dsl for real table action columns

Before migrating the remaining real pages, I found that the existing v2 table surface was still too narrow for `go-go-course` material/session tables. Demo pages only needed row selection, reorder, and remove, but the real pages needed explicit per-row action columns such as Open, Edit, and Delete while still preserving URL-backed row selection.

I added the smallest v2 API extension needed for P6: collection empty text, table class names, explicit action columns, and a separate table row-select action. The row-select separation matters because the original builder stored `rowSelect(...)` in `Actions.Open`, which could both create an unintended Open column and override selection behavior.

### Prompt Context

**User prompt (verbatim):** (same as Step 25)

**Assistant interpretation:** Finish all remaining P6 work, including any v2 API support required to migrate real pages safely.

**Inferred user intent:** Complete the hard-cutover page migration rather than stopping at the agenda editor proof-of-concept.

**Commit (code):** 166e8dc — "Extend data v2 table action API"

### What I did
- Added `TableSpec` and `TableActionColumnSpec` to `pkg/widgetdsl/v2/spec/types.go`.
- Added validation for explicit table action columns and table-level row-select actions.
- Extended lowering so `CollectionSpec.Table.ClassName` lowers to DataTable `className`, `CollectionSpec.Empty` lowers to `emptyMessage`, and table action columns lower to `actionButton` cells.
- Split `.table(t => t.rowSelect(...))` away from `CollectionActions.Open` by storing it in `CollectionSpec.Table.RowSelect`.
- Added Goja builder methods:
  - `collection.empty(message)`
  - `table.className(className)`
  - `table.actionColumn(id, header, label, action, { maxWidth })`
  - corrected `table.rowSelect(action)` to set `Table.RowSelect`
- Updated `data.v2.dsl` TypeScript declarations and the TS fixture to compile the new surface.
- Added a Go runtime test covering explicit Open/Delete action columns, empty text, className, and action lowering.

### Why
- `adminMaterialTable` needs Open/Edit/Delete columns; v2 had no general action-column mechanism.
- `sessions.js` needs row click to select the row and a separate transcript action column; reusing one `Open` action for both was semantically wrong.

### What worked
- The new `actionColumn` API lets real pages express link-button-like behavior as typed action data without reintroducing v1 `data.cell.*` helpers.
- The TypeScript fixture still compiles and now covers `.empty`, `.className`, and `.actionColumn`.
- `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite -count=1` passed.

### What didn't work
- First formatting/test run failed because I introduced an unterminated DTS string literal in `pkg/widgetdsl/typescript.go`:

```text
pkg/widgetdsl/typescript.go:200:3: string literal not terminated
```

- I fixed the declaration line to close the string and reran `gofmt` plus the widgetdsl/provider tests successfully.

### What I learned
- `rowSelect` and `Open` are not the same concept. A row can be selected for details while also having an Open Transcript or Open File button; conflating them produces duplicate or wrong action behavior.
- A typed table API can still support pragmatic row action columns without falling back to raw v1 column/cell specs.

### What was tricky to build
- The tricky bit was noticing the design bug before migrating pages: the initial extension would have let sessions add a transcript action column, but `rowSelect` would also have generated an automatic Open column and changed the row-click behavior. Splitting row selection into `TableSpec.RowSelect` fixed the underlying representation instead of papering over it in page code.
- The public TypeScript declaration and Goja runtime needed to evolve together; the DTS fixture made the string-literal typo visible immediately.

### What warrants a second pair of eyes
- Review whether `actionColumn(id, header, label, action, options)` is the right final API name/shape before P7 locks the surface with lint rules.
- Review whether `CollectionActions.Open` is still needed now that table rowSelect and explicit table action columns cover the real cases.

### What should be done in the future
- Add docs examples for explicit table action columns in the public v2 API reference during P7 documentation cleanup.
- Consider richer action-column options later (variant/size/disabled) if real pages need them.

### Code review instructions
- Start in `pkg/widgetdsl/v2/spec/types.go`, `validate.go`, and `lower.go` to review the typed representation.
- Then review `pkg/widgetdsl/v2_builders.go` and `pkg/widgetdsl/typescript.go` for the JS/DTS API.
- Validate with:
  - `cd rag-evaluation-system && go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite -count=1`

### Technical details
- Successful command: `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite -count=1`
- Code commit: `166e8dcb2d2e5a9c052663d0134092d21f175673`

## Step 27: Migrate remaining go-go-course tables to data.v2.dsl

With table action columns available, I migrated the remaining real `go-go-course` page code away from `data.dsl`. The shared admin material table now uses a typed `CourseMaterial` schema and v2 explicit action columns, and the sessions browse page now uses a typed selectable v2 table with a separate Open Transcript action column.

This completes the real-page part of P6: the agenda editor, material tables, and sessions browse table are no longer authored with `dataDsl.dataTable`, `dataDsl.cell.*`, or the v1 `dataDsl.collection` option bag.

### Prompt Context

**User prompt (verbatim):** (same as Step 25)

**Assistant interpretation:** Apply the new v2 table API to the remaining real course pages and validate source plus generated runtime code.

**Inferred user intent:** Remove old table grammar from `go-go-course` real pages, not just demos.

**Commit (code):** 0bf3820 — "Migrate remaining course tables to data.v2 DSL"

### What I did
- Changed `cmd/go-go-course/lib/pages/admin-common.js` from `{ ui, dataDsl, contextWindow }` to `{ ui, dataV2, contextWindow }`.
- Replaced `adminMaterialTable`'s raw `dataDsl.dataTable({ columns: dataDsl.cell.* })` construction with:
  - `dataV2.schema("CourseMaterial")`
  - `dataV2.collection(...).schema(...).empty(...).edit(...).table(...).toIR()`
  - explicit Open/Edit action columns via `table.actionColumn(...)`
  - delete via the v2 remove action path, preserving `admin-delete-course-material`
- Changed `cmd/go-go-course/lib/pages/sessions.js` from `dataDsl.dataTable` to a v2 selectable `Session` collection.
- Added an explicit transcript action column while keeping row click selection via `/pages/sessions?scope=...&selected=${row.sessionId}`.
- Removed `dataDsl` from `cmd/go-go-course/server.js` dependency wiring because no source page consumed it anymore.
- Rebuilt xgoja output so generated mirrors under `internal/xgojaruntime/...` matched source.

### Why
- P6.2 required migrating the remaining media/material and session-browse table examples to v2 APIs.
- The hard-cutover goal is to make real page code use typed/fluent v2 builders rather than preserving v1 option bags indefinitely.

### What worked
- `node --check` passed for source JS.
- `GOWORK=off go test ./...` passed in `go-go-course`.
- `pnpm --dir cmd/go-go-course/webapp typecheck` passed.
- `GOWORK=off make -C cmd/go-go-course build XGOJA_LOCAL_DIR=/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-goja` passed and refreshed generated mirrors.
- Generic smoke passed on port `18792`: `11 passed, 0 failed`.
- Targeted P6 API smoke passed on port `18793` for:
  - sessions page id/table/selected-query action/transcript action column;
  - admin material page id/table class/delete action/open action;
  - admin Course CMS page id/agenda delete/material table/media library.
- A grep over source and generated course runtime found no remaining `dataDsl`, `data.dsl`, `dataTable(`, or `cell.` usage in go-go-course page/server code.

### What didn't work
- N/A for the page migration after Step 26 fixed the rowSelect/action-column split.

### What I learned
- Route templates should avoid substituting a whole URL as `${row.href}` because the frontend navigate interpolation URL-encodes variable values. The material-table migration therefore uses route-specific templates such as `/pages/slides?slide=${row.id}`, `/pages/handouts?doc=${row.id}`, and `/course-assets/${row.file}`.
- The existing server handlers already accept row context, so material delete and session selection did not require backend changes.

### What was tricky to build
- `adminMaterialTable` is shared across slide, handout, and media rows. The open target had to be selected from homogeneous row kind rather than using the raw `href` field as a whole URL template.
- Handouts have an extra Edit action while slides/media do not. The new explicit table action columns made this conditional column straightforward without reintroducing v1 cell helpers.

### What warrants a second pair of eyes
- Review route templates in `adminMaterialTable`, especially the assumption that each rendered table contains homogeneous `kind` rows.
- Review whether the sessions table should keep an explicit Open Transcript action column or rely solely on the selected-session panel buttons.

### What should be done in the future
- Add a permanent targeted API/browser smoke for sessions and admin material tables during P7.
- Update public docs so examples match the new `actionColumn` approach.

### Code review instructions
- Start with `go-go-course/cmd/go-go-course/lib/pages/admin-common.js` and `go-go-course/cmd/go-go-course/lib/pages/sessions.js`.
- Compare generated mirrors under `go-go-course/cmd/go-go-course/internal/xgojaruntime/xgoja_embed/jsverbs/minitrace_viz_site/lib/pages/`.
- Validate with:
  - `cd go-go-course && GOWORK=off go test ./...`
  - `cd go-go-course && find cmd/go-go-course/lib -name '*.js' -print0 | xargs -0 -n1 node --check && node --check cmd/go-go-course/server.js`
  - `cd go-go-course && pnpm --dir cmd/go-go-course/webapp typecheck`
  - `cd go-go-course && GOWORK=off make -C cmd/go-go-course build XGOJA_LOCAL_DIR=/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-goja`

### Technical details
- Code commit: `0bf382064c421a77b533c7f1725d03404eba9318`
- Successful smoke command: copied `cmd/go-go-course/test-fixtures/smoke-test.sh` to `/tmp/gocourse-smoke-18792.sh`, changed `PORT=18787` to `PORT=18792`, and ran it from `cmd/go-go-course`.
- Successful targeted page smoke used port `18793` and `PUT /api/user-session` with `displayName=admin_p6` before fetching admin pages.

## Step 28: Remove legacy data.dsl from the course runtime selection

After the page migration, `data.dsl` was no longer imported by `go-go-course` source or generated page code. I removed the legacy module from the course xgoja package specs so the built course runtime selects `data.v2.dsl` without also carrying the old public `data.dsl` table/cell module.

This is the P6 cleanup step for the course runtime: real course pages no longer author against v1, and the xgoja runtime no longer selects the legacy data module for this app.

### Prompt Context

**User prompt (verbatim):** (same as Step 25)

**Assistant interpretation:** Complete the public/runtime cleanup portion of P6 after real pages have stopped importing the legacy data module.

**Inferred user intent:** Make the hard cutover visible in runtime packaging, not only in source page style.

**Commit (code):** 49dc871 — "Remove legacy data DSL from course runtime"

### What I did
- Removed the `rag-widget-site` selected module entry for `data.dsl` from `cmd/go-go-course/xgoja.package.yaml`.
- Replaced the stale `data.dsl` entry in `cmd/go-go-course/xgoja.yaml` with `data.v2.dsl`.
- Rebuilt xgoja output, which removed the generated `data.dsl` module selection from `cmd/go-go-course/internal/xgojaruntime/xgoja_runtime.gen.go`.
- Verified with `rg` that `go-go-course` source, package specs, and generated site code no longer mention `data.dsl`, `dataDsl`, `dataTable(`, or `cell.`.

### Why
- P6.3 required deleting old public v1 exports from the selected v2/course runtime path. Since the shared provider still offers `data.dsl` for historical docs/tests, the safe app-level cleanup is to stop selecting it in the course app package.

### What worked
- Rebuild after module removal passed.
- `GOWORK=off go test ./...` passed.
- JS source syntax checks passed.
- `pnpm --dir cmd/go-go-course/webapp typecheck` passed.
- Generic smoke passed again on port `18794`: `11 passed, 0 failed`.

### What didn't work
- N/A.

### What I learned
- There were two xgoja specs to update: `xgoja.package.yaml` and the older `xgoja.yaml`. They were not identical: `xgoja.package.yaml` already selected both `data.dsl` and `data.v2.dsl`, while `xgoja.yaml` still selected only `data.dsl`.

### What was tricky to build
- The provider-level `data.dsl` module still exists for legacy tests/docs in `rag-evaluation-system`; deleting it globally would be a larger compatibility decision. The correct P6 scope was app-level removal from `go-go-course` selected modules after confirming no source import remained.

### What warrants a second pair of eyes
- Review whether P7 should add a generated-runtime assertion that `go-go-course` does not select `data.dsl`.
- Review whether the legacy provider/docs should be renamed from "current" to "legacy-only" now that course app no longer uses it.

### What should be done in the future
- P7 should add CI/lint checks rejecting `require("data.dsl")`, `dataDsl`, `dataTable(`, and `cell.*` in go-go-course page code.
- P7 should update public docs and final handoff material to state that `go-go-course` has completed the course-app v2 data cutover.

### Code review instructions
- Start with `go-go-course/cmd/go-go-course/xgoja.package.yaml` and `go-go-course/cmd/go-go-course/xgoja.yaml`.
- Then inspect `go-go-course/cmd/go-go-course/internal/xgojaruntime/xgoja_runtime.gen.go` to confirm the generated module list dropped `data.dsl`.
- Validate with:
  - `cd go-go-course && rg -n 'data\.dsl|dataDsl|dataTable\(|cell\.' cmd/go-go-course/lib cmd/go-go-course/server.js cmd/go-go-course/xgoja.package.yaml cmd/go-go-course/xgoja.yaml cmd/go-go-course/internal/xgojaruntime/xgoja_embed/jsverbs/minitrace_viz_site -S`
  - `cd go-go-course && GOWORK=off make -C cmd/go-go-course build XGOJA_LOCAL_DIR=/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-goja`

### Technical details
- Code commit: `49dc8716aa44db7ecb663be9545a9d6a5c5b1313`
- Successful final smoke command: copied `cmd/go-go-course/test-fixtures/smoke-test.sh` to `/tmp/gocourse-smoke-18794.sh`, changed `PORT=18787` to `PORT=18794`, and ran it from `cmd/go-go-course`.
- Ticket tasks checked: `[37] P6.2 Rewrite media library and session browse examples to v2 APIs`, `[38] P6.3 Delete old public v1 exports from v2 modules`.
