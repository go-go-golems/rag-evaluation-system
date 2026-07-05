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
Summary: "Chronological investigation diary for the Goja DSL catalogue and base research. Records what was tried, what worked, what failed, and what to do next."
WhatFor: "Record the research journey so a senior researcher can resume without re-reading every source."
WhenToUse: "Read before resuming work on this ticket."
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
