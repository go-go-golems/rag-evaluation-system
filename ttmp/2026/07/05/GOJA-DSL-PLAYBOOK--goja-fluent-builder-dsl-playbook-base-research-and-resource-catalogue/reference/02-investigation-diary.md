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
