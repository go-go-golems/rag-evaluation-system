---
Title: Research logbook — resource assessment
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
Summary: "Per-resource assessment of every file, README, and external page read while building the Goja DSL catalogue. Each entry records what was researched, what was looked for, why the resource was chosen, how it was found, what was useful, what was not useful, what is out of date or wrong, and what would need updating."
WhatFor: "Track which resources are useful, out of date, or need updating, so a senior researcher can prioritise the reflection and assessment pass without re-reading everything."
WhenToUse: "Consult before re-reading any source; update when a resource changes or is found to be wrong."
---

# Research logbook — resource assessment

## How to use this logbook

This logbook records every resource read while building the base research in `design-doc/01-goja-dsl-catalogue-and-base-research.md`. Each entry uses the fixed fields requested by the ticket:

- **What I was researching** — the question driving the read.
- **What I was looking for in this document** — the specific thing sought here.
- **Why I chose it** — why this resource over alternatives.
- **How I found the resource** — path from the starting request to this file.
- **What I found useful** — concrete takeaways.
- **What I didn't find useful** — noise, gaps, irrelevancies.
- **What is out of date / what was wrong** — staleness, errors, contradictions.
- **What would need updating** — concrete fixes for the senior researcher.

Resources are grouped by repository. A final **Cross-cutting** section covers shared infrastructure. Status legend: 🟢 useful & current · 🟡 useful but partial/stale · 🔴 out of date or wrong · ⚪ skeleton/empty.

---

## A. goja-bleve — `~/code/wesen/go-go-golems/goja-bleve`

### A1. `README.md` (root) — 🟢

- **What I was researching:** the public JS API surface, the batch lifecycle, and vector/KNN/hybrid scoring.
- **What I was looking for:** factory names, fluent chain examples, the `vectors` build-tag requirement, and how hybrid scoring differs from the rag-eval service.
- **Why I chose it:** root README is the canonical entry point for any goja-* module.
- **How I found it:** listed `~/code/wesen/go-go-golems/goja-bleve/` per the ticket, then read `README.md`.
- **What I found useful:** the minimal JS shape table, the field-builder fluent example, the batch single-use rule ("batch has already been executed"), the explicit contrast with `internal/services/search/hybrid.go`, the `-tags=vectors` + FAISS requirements.
- **What I didn't find useful:** some sections are phase-numbered ("Phase 7 core surface", "Phase 8 will expand…") which read as historical rather than current.
- **What is out of date / what was wrong:** "Phase 8 will expand it into full API documentation and golden declaration tests" — need to verify whether that landed. The phase language should be replaced with current-state language.
- **What would need updating:** drop phase framing; state current capabilities declaratively; add a one-line pointer to the typed-ref machinery in `pkg/api_types.go` (currently undocumented in the README).

### A2. `pkg/api_types.go` — 🟢

- **What I was researching:** how goja-bleve enforces runtime type safety on JS-facing handles.
- **What I was looking for:** the ref/tag model, the type-extraction helper, and how wrappers are constructed.
- **Why I chose it:** `grep` for `fieldBuilder` pointed here; it is the type system of the module.
- **How I found it:** `rg -n 'fieldBuilder' pkg/api_mapping.go` → `api_types.go` for the ref definitions.
- **What I found useful:** `refKind` enum (`refKindIndex`, `refKindFieldBuilder`, …), `refBase{api, kind, closed}`, the typed ref structs embedding `refBase`, `getTypedRef[T]` generic extractor (the key reusable primitive), `newWrapper(ref, kind)`.
- **What I didn't find useful:** nothing — this file is the gold.
- **What is out of date / what was wrong:** none observed.
- **What would need updating:** none in-file. For the playbook: this machinery should be extracted into a shared `fluent` package and documented as the canonical runtime-typecheck substrate.

### A3. `pkg/api_mapping.go` (fieldBuilder + installFieldOptions) — 🟢

- **What I was researching:** how a fluent field builder is wired (same-object chaining + `.build()` terminal).
- **What I was looking for:** the chain method pattern, where validation lives, and how `.build()` changes the wrapper kind.
- **Why I chose it:** direct implementation of the fluent pattern the playbook wants.
- **How I found it:** from A2, the mapping file is the natural next read.
- **What I found useful:** `fieldBuilder()` returns a wrapper; each type method (`text`, `keyword`, `number`, …) mutates `ref.mapping` and returns the same `obj`; `vector(dims)` returns `(obj, err)`; `installFieldOptions` adds `name`/`analyzer`/`store`/`index`/…; `.build()` returns a **new** `fieldMappingRef` wrapper or `(nil, err)` if no type set.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none in-file. Playbook should lift this exact shape as the reference example.

### A4. `pkg/module.go` (`hiddenRefKey`, `mustSet`, exports) — 🟢

- **What I was researching:** how the hidden Go reference is attached to JS objects and how methods are bound.
- **What I was looking for:** the non-enumerable key name, the `mustSet` helper, the export list.
- **Why I chose it:** needed the attachment mechanism for the typed-ref model.
- **How I found it:** referenced from `api_types.go`.
- **What I found useful:** `hiddenRefKey = "__bleve_ref"`, `mustSet(o, key, value)`, the `field` export at line 162.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none.

### A5. `docs/faiss-xgoja-playbook.md` — 🟢

- **What I was researching:** the FAISS build/link/runtime requirements for vector support.
- **What I was looking for:** exact `CGO_LDFLAGS`, the Bleve-compatible FAISS fork requirement, the loader-path fix.
- **Why I chose it:** vector search is part of the bleve DSL; the playbook must mention build constraints.
- **How I found it:** `ls docs/` after reading the README.
- **What I found useful:** the executive summary with the 5-step requirement list and the canonical `make test-vectors` command; the explicit `CGO_LDFLAGS` and `-ldflags "-r /usr/local/lib"`.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none observed (cannot verify FAISS fork currency without building).
- **What would need updating:** verify the FAISS fork URL and version are still the Bleve-compatible one.

### A6. `docs/quickstart.md`, `docs/README.md` — 🟡

- **What I was researching:** quickstart flow and docs index.
- **What I was looking for:** a runnable end-to-end example and a doc map.
- **Why I chose it:** completeness.
- **How I found it:** `ls docs/`.
- **What I found useful:** confirmed a quickstart exists.
- **What I didn't find useful:** did not read in full; deferred to README which is more complete.
- **What is out of date / what was wrong:** not assessed in depth.
- **What would need updating:** ensure quickstart matches the current factory names (`bleve.field().text().build()` etc.).

---

## B. goja-dbus — `~/code/wesen/go-go-golems/goja-dbus`

### B1. `README.md` (root) — 🟢

- **What I was researching:** the D-Bus module's JS API, typed values, and bus/method/signal builders.
- **What I was looking for:** factory names, the fluent chain shape, policy enforcement, runtime ownership.
- **Why I chose it:** root README.
- **How I found it:** listed `goja-dbus` per the ticket.
- **What I found useful:** the three runnable examples (GetId, typed values, signal subscription), the implemented-vs-deferred list, the explicit "all callbacks settle on the runtime owner" rule, the docmgr ticket pointer.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** "Deferred: … JavaScript-backed D-Bus service export" — confirm still deferred.
- **What would need updating:** add a note on the typed-value model (`{signature, value}`) since it is the cleanest example of typed helpers.

### B2. `pkg/dbusgoja/builders.go` — 🟢

- **What I was researching:** the fluent bus + method-call builder implementation.
- **What I was looking for:** the chain methods (timeout/policy/connect/destination/object/interface/method/in/out/call) and how they re-wrap.
- **Why I chose it:** the cleanest composable grammar in the ecosystem.
- **How I found it:** `rg -n 'obj.Set\("' pkg/dbusgoja/builders.go`.
- **What I found useful:** every builder method returns `goja.Value` and re-wraps the same builder; `call()` resolves to a Promise; `close()` is explicit; policy is enforced at `connect()`.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none. Playbook should use this as the "composable grammar" reference.

### B3. `pkg/dbusgoja/signals.go`, `typed_values.go`, `errors.go` — 🟢

- **What I was researching:** signal subscription builder, typed value helpers, error shape.
- **What I was looking for:** the signal chain (sender/path/interface/member/listen), the typed-value payload (`{signature, value}`), the error object shape.
- **Why I chose it:** completes the dbus API picture.
- **How I found it:** from B2.
- **What I found useful:** signal builder mirrors method builder; typed values carry signature for marshaling; `DBusError{name:"DBusError", code:"ERR_DBUS"}`.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none.

### B4. `ttmp/2026/06/15/GOJA-DBUS-DESIGN--goja-d-bus-module-intern-design-guide/` — 🟡 (not read in full)

- **What I was researching:** the design rationale and intern-facing implementation guide.
- **What I was looking for:** the runtime-ownership rule documentation.
- **Why I chose it:** referenced by the README.
- **How I found it:** README pointer.
- **What I found useful:** confirmed it exists and is the detailed design doc.
- **What I didn't find useful:** did not read in full this pass.
- **What is out of date / what was wrong:** not assessed.
- **What would need updating:** senior researcher should read this fully and extract the runtime-ownership rule for the playbook.

---

## C. go-minitrace — `…/go-minitrace`

### C1. `README.md` (root) — 🟢

- **What I was researching:** what go-minitrace is and how its JS DSL is exposed.
- **What I was looking for:** the pipeline framing (reduction, not reading), install, quick start.
- **Why I chose it:** root README.
- **How I found it:** listed in the working directory.
- **What I found useful:** the "reduction pipeline" framing, supported sources list, `convert`/`query`/`serve` commands.
- **What I didn't find useful:** the JS DSL is not described in the README — it focuses on the CLI.
- **What is out of date / what was wrong:** none.
- **What would need updating:** add a section on the `minitracejs` builder DSL (factory entry points and fluent example).

### C2. `pkg/minitracejs/module.go` — 🟢

- **What I was researching:** the JS module exports (factory entry points).
- **What I was looking for:** the `require("minitrace")` surface.
- **Why I chose it:** the module loader is the API surface.
- **How I found it:** `rg -n 'exports.Set'`.
- **What I found useful:** the factory list: `importer`, `db`, `sources`, `importPolicy`, `cache`, `limits`, `query`, `view`, `session`, plus `runtime` settings and `sql` helpers.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none.

### C3. `pkg/minitracejs/builders.go` — 🟢

- **What I was researching:** the source-set builder (Pattern B exemplar).
- **What I was looking for:** chain methods, error accumulation, `Validate()`/`Build()` terminals.
- **Why I chose it:** the simplest full builder to document.
- **How I found it:** from C2.
- **What I found useful:** `SourceSetBuilder{sources, last, errors}`, re-wrap pattern (`sourcesBuilderObject` returns fresh obj each call), `Validate() ValidationResult`, `Build() (SourceSet, error)`, dedup of file sources, empty-path/empty-content errors.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none. Reference example for Pattern B.

### C4. `pkg/minitracejs/db_builder.go` — 🟢

- **What I was researching:** the largest builder and its validation discipline.
- **What I was looking for:** `DBBuilder` fields, `ValidationResult`, `DBHandle`, error accumulation.
- **Why I chose it:** the most complete builder; shows the pattern at scale.
- **How I found it:** from C2.
- **What I found useful:** `DBBuilder` accumulates `errors []string`; `ValidationResult{Valid, Errors}`; `DBHandle` with cache diagnostics; `dbSource` struct.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none.

### C5. `pkg/minitracejs/import_builder.go` — 🟢

- **What I was researching:** a builder with multiple terminals (not just `Build`).
- **What I was looking for:** the multi-terminal pattern (Detect/Convert/Preview/Diagnostics/Save).
- **Why I chose it:** shows that "terminal" is not always `build()`.
- **How I found it:** from C2.
- **What I found useful:** `Content/File/Name/SourcePath/AutoDetect/Format/Strict/Into/SessionID/Overwrite` chain; `Detect/Convert/Converted/Preview/Diagnostics/Save` terminals returning varied types.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none.

### C6. `pkg/minitracejs/query_view_session.go` — 🟢

- **What I was researching:** the query recipe builder (a grammar of recipe kinds).
- **What I was looking for:** how a "grammar of kinds" is modeled.
- **Why I chose it:** closest thing to a composable operator grammar.
- **How I found it:** from C2.
- **What I found useful:** `QueryRecipeBuilder` with kind setters (SessionSummary/TurnRows/ToolRows/EventRows/TurnBlockRows/TokenUsageRows/TranscriptRows/TimelineRows) + grouping (BySession/ByTurn/ByRole/ByTool) + SessionID/IncludeTools + `Build()` returning `QueryRecipe{Name, SQL, Args, Description, Output}`.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none. Good model for "grammar of operators" in the playbook.

### C7. `pkg/minitracejs/typescript.go` — 🟡

- **What I was researching:** the TS descriptor for the module.
- **What I was looking for:** whether it emits precise types.
- **Why I chose it:** needed to assess compile-time type coverage.
- **How I found it:** `ls pkg/minitracejs/`.
- **What I found useful:** a descriptor exists.
- **What I didn't find useful:** minimal — did not emit per-builder named types in the snippet surveyed.
- **What is out of date / what was wrong:** none observed.
- **What would need updating:** expand to emit named builder types per the playbook rule.

---

## D. rag-evaluation-system widgetdsl — `…/rag-evaluation-system`

### D1. `pkg/widgetdsl/module.go` — 🟡 (large, mixed quality)

- **What I was researching:** the five-module registration and helper wiring.
- **What I was looking for:** module names, the `moduleSpecs` table, helper maps, page/cell/action wiring.
- **Why I chose it:** the central module file (1215 lines).
- **How I found it:** `rg -n 'moduleSpecsByName|DataModuleName'`.
- **What I found useful:** the five module names and their specs, the generic-primitive promotion list (breadcrumbs/emptyState/…), the `setExport` helper, the cell/action sub-objects.
- **What I didn't find useful:** the file mixes registration, helpers, and rendering; hard to follow. No typed builder objects.
- **What is out of date / what was wrong:** none factually; the design is the problem (Pattern C).
- **What would need updating:** this is the "not great" DSL — the playbook should define the target pattern and a migration sketch.

### D2. `pkg/widgetdsl/grammar.go` — 🟡

- **What I was researching:** the data grammar verbs (`f`, `schema`, `record`, `collection`, `urlParam`, `formPost`).
- **What I was looking for:** field roles, schema construction, validation path.
- **Why I chose it:** the most recent grammar addition (RAGEVAL-UI-GRAMMAR).
- **How I found it:** `rg -n 'installDataGrammar'`.
- **What I found useful:** the `fieldRoles` list (key/primary/short/prose/count/size/measure/date/status/tags/media/href), `gridableRoles`, `schemaCtor` preserving insertion order with `__ragSchema` tag, `record`/`collection` validation.
- **What I didn't find useful:** validation by `panic(r.vm.NewGoError(...))` — hostile to JS callers.
- **What is out of date / what was wrong:** the `__ragSchema` magic-string tag is a code smell; the role/type distinction is documented only inline.
- **What would need updating:** replace panic with `(value, error)` terminals; replace magic tag with a typed wrapper; document the role model in the README.

### D3. `pkg/widgetdsl/typescript.go` — 🔴 (weak by design)

- **What I was researching:** the TS declaration generation for the widget DSLs.
- **What I was looking for:** whether per-field/per-cell types are emitted.
- **Why I chose it:** the compile-time-type requirement hinges on this.
- **How I found it:** `cat pkg/widgetdsl/typescript.go`.
- **What I found useful:** it uses `tsgen/spec.Module` correctly.
- **What I didn't find useful:** everything is `Props = Record<string, any>` and `[key: string]: any`. This is the concrete evidence that the DSL fails the "compile-time types" requirement.
- **What is out of date / what was wrong:** not wrong per se — open-ended by design ("individual component props remain open-ended by design"). But that design choice is exactly what the playbook wants to reverse for type-checked builders.
- **What would need updating:** emit named builder types and precise `Param`/`Returns` for grammar verbs; keep open `Props` only for raw component helpers if retained.

### D4. `pkg/widgetdsl/module_test.go`, `grammar_test.go` — 🟢

- **What I was researching:** how the DSL is tested and what the examples look like.
- **What I was looking for:** runnable examples and the API shape in assertions.
- **Why I chose it:** tests are the most reliable API documentation.
- **How I found it:** `ls pkg/widgetdsl/`.
- **What I found useful:** `data.cell.field("title")`, `data.cell.status("status", {icon:true})`, `data.recipes.masterDetailTable`, `data.f.primary({required:true, maxLength:160})`.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none.

### D5. `ttmp/2026/06/02/RAGEVAL-UI-DSL--…/index.md` — 🟢

- **What I was researching:** the original widget DSL design history and related files.
- **What I was looking for:** design-doc list and related-file map.
- **Why I chose it:** the foundational ticket.
- **How I found it:** `docmgr ticket list` + `rg -l 'dsl' ttmp/`.
- **What I found useful:** the RelatedFiles map (server.go, dsl_handlers.go, module.go, design docs 02/03), confirming the module was originally at `internal/dsl/widgetdsl/` and moved to `pkg/widgetdsl/`.
- **What I didn't find useful:** some RelatedFiles paths reference the old `internal/dsl/widgetdsl/` location — may be stale.
- **What is out of date / what was wrong:** RelatedFiles paths may point to the old `internal/dsl/widgetdsl/module.go` (now `pkg/widgetdsl/module.go`).
- **What would need updating:** verify and update RelatedFiles to current `pkg/widgetdsl/` paths.

### D6. `ttmp/2026/07/04/RAGEVAL-UI-GRAMMAR--…/design-doc/02-…` — 🟡 (not read in full)

- **What I was researching:** the grammar-verbs API sketch (the design that produced `data.dsl` grammar).
- **What I was looking for:** the intended API and the module-reorganisation story.
- **Why I chose it:** the most recent design doc for the grammar.
- **How I found it:** `ls ttmp/2026/07/04/`.
- **What I found useful:** confirmed it exists; title references "Grammar verbs in data.dsl and ui.dsl, module reorganization, no new module".
- **What I didn't find useful:** not read in full this pass.
- **What is out of date / what was wrong:** not assessed.
- **What would need updating:** senior researcher should read design-doc 01 (analysis of alternatives) and 02 (API sketch) fully — they are the direct antecedent of the "not great" grammar.

### D7. `ttmp/2026/06/03/rag-eval-scripting-expansion--…/` — 🟡 (not read in full)

- **What I was researching:** the xgoja + goja-text jsverbs expansion.
- **What I was looking for:** how jsverbs integrate with the rag-eval DSLs.
- **Why I chose it:** referenced in the ticket list as a scripting expansion.
- **How I found it:** `rg -l 'goja|dsl|builder' ttmp/`.
- **What I found useful:** confirmed it exists with an intern guide + design doc + diary.
- **What I didn't find useful:** not read in full.
- **What is out of date / what was wrong:** not assessed.
- **What would need updating:** read fully; relevant to the jsverbs bridge and xgoja integration.

---

## E. go-go-goja — `…/go-go-goja`

### E1. `README.md` (root) — 🟢

- **What I was researching:** the go-go-goja host, module layout, and engine composition API.
- **What I was looking for:** folder layout, the canonical runtime composition steps, the help-tree pointers.
- **Why I chose it:** root README.
- **How I found it:** listed in the working directory.
- **What I found useful:** the explicit `engine.NewRuntimeFactoryBuilder() → Build() → NewRuntime(WithStartupContext, WithLifetimeContext) → Close()` flow, the two separate help trees (`goja-repl help` vs `xgoja help`), the planned-auth help pointers.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** "Legacy convenience wrappers (`engine.New()`, …) were removed" — confirm this is still true (it is, per `pkg/engine/`).
- **What would need updating:** none.

### E2. `pkg/engine/` (factory, options, runtime, middleware) — 🟢

- **What I was researching:** runtime ownership and module composition.
- **What I was looking for:** the builder/factory/runtime APIs and the startup vs lifetime context distinction.
- **Why I chose it:** every DSL host depends on this.
- **How I found it:** from E1.
- **What I found useful:** `factory.go`, `options.go`, `runtime.go`, `module_middleware.go`, `module_roots.go`, `module_specs.go` — the ownership model the playbook must respect.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none.

### E3. `modules/uidsl/` (components, node, render, tests) — 🟡

- **What I was researching:** the hyperscript element DSL.
- **What I was looking for:** the `ui.div(...)` API, class/style handling, render output.
- **Why I chose it:** a distinct DSL pattern (D) in the core repo.
- **How I found it:** `ls modules/`.
- **What I found useful:** `Element{Tag, Attrs, Children}`, conditional class arrays, style objects, `render()` to HTML, benchmarks.
- **What I didn't find useful:** no README section; documented only via tests.
- **What is out of date / what was wrong:** none.
- **What would need updating:** add a README section or a Glazed help page documenting the element DSL and its class/style semantics.

### E4. `modules/express/express.go`, `auth_builders.go` — 🟢

- **What I was researching:** the express route DSL and the planned-auth fluent builders.
- **What I was looking for:** the `express.user().required().mfaFresh()` chain and how specs are validated.
- **Why I chose it:** a fluent builder with a side-channel typed-lookup variant.
- **How I found it:** `ls modules/express/`.
- **What I found useful:** `builderStore` with `sync.Map[*goja.Object]*SecuritySpec`, `newUserBuilder`/`newResourceBuilder`, `authSpec` validating provenance, duration parsing with `(goja.Value, error)`.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none. Good contrast to goja-bleve's hidden-key approach (decision record in the design doc).

### E5. `pkg/tsgen/spec/types.go` — 🟢

- **What I was researching:** the TypeScript declaration spec model.
- **What I was looking for:** the type system primitives available for compile-time types.
- **Why I chose it:** the substrate for "compile-time types from generated declarations."
- **How I found it:** `ls pkg/tsgen/spec/`.
- **What I found useful:** `Module/Function/Param/TypeRef/Field` with full `TypeKind` (string/number/boolean/any/unknown/void/never/named/array/union/object) and `Bundle` for multi-module render.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none. Playbook must mandate named `TypeRef`s per builder.

### E6. `pkg/jsverbs/`, `pkg/jsdoc/` — 🟡

- **What I was researching:** the JS-command-to-glazed bridge and JSDoc→spec extraction.
- **What I was looking for:** how JS verbs compile to typed commands and how types can be extracted from JSDoc.
- **Why I chose it:** relevant to the "composable grammar that compiles to a typed core" idea.
- **How I found it:** `ls pkg/`.
- **What I found useful:** `jsverbs/command.go` (`VerbSpec` → glazed `Command`), `jsdoc/` (extract/batch/export/server/watch).
- **What I didn't find useful:** not read in full.
- **What is out of date / what was wrong:** not assessed.
- **What would need updating:** senior researcher should assess whether JSDoc extraction is a viable alternative authoring path for builder types.

### E7. `pkg/xgoja/` (provider system) — 🟡

- **What I was researching:** the xgoja provider packaging.
- **What I was looking for:** how a module becomes a provider package.
- **Why I chose it:** needed to understand `xgoja.yaml` wiring.
- **How I found it:** `ls pkg/xgoja/`.
- **What I found useful:** provider packages (e.g. `pkg/xgoja/providers/widgetsite/`).
- **What I didn't find useful:** not read in full.
- **What is out of date / what was wrong:** not assessed.
- **What would need updating:** document the provider packaging contract for the playbook.

---

## F. goja-text — `~/code/wesen/go-go-golems/goja-text`

### F1. `README.md` (root) — 🟢

- **What I was researching:** the markdown/sanitize/extract/template modules.
- **What I was looking for:** the module table, the markdown walk design, the template fluent chain, the xgoja build.
- **Why I chose it:** root README.
- **How I found it:** listed per the ticket.
- **What I found useful:** the module-purpose table, the `walk()` design rationale ("Go API stays small, JS owns the question"), the paired help entries per module, the `make build-xgoja` flow.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none. Strongest documentation example; playbook should cite the paired API-reference + user-guide pattern.

### F2. `xgoja.yaml` — 🟢

- **What I was researching:** the module/provider/help wiring config.
- **What I was looking for:** the YAML schema for packages/modules/commands/jsverbs/help.
- **Why I chose it:** the canonical example of an xgoja module config.
- **How I found it:** `cat xgoja.yaml`.
- **What I found useful:** `target.kind: xgoja`, `packages` (with `replace`), `modules` (with `as` + `config`), `commands` (eval/run/repl/jsverbs), `jsverbs` (embed), `help.sources`.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none. Reference config for the playbook's "how to ship a DSL" section.

---

## G. goja-git — `~/code/wesen/go-go-golems/goja-git`

### G1. `README.md` (root) — 🟢

- **What I was researching:** the Git DSL API.
- **What I was looking for:** factory names, options-object shapes, examples.
- **Why I chose it:** root README.
- **How I found it:** listed per the ticket.
- **What I found useful:** `init/open/status/add/commit/log/branch/checkout/tags/diff/filterRepo` with options objects.
- **What I didn't find useful:** it is Pattern E (imperative) — included for contrast, not as a model.
- **What is out of date / what was wrong:** "Go 1.25.5 or later" — verify currency.
- **What would need updating:** none for the playbook (contrast example).

---

## H. goja-github-actions — `~/code/wesen/go-go-golems/goja-github-actions`

### H1. `README.md` (root) — 🟢

- **What I was researching:** the `@actions/*` polyfills and `@goja-gha/ui` report DSL.
- **What I was looking for:** module list, workspace-first resolution, examples, token scopes.
- **Why I chose it:** root README.
- **How I found it:** listed per the ticket.
- **What I found useful:** the polyfill surface, `run`/`doctor` commands, workspace-first path resolution, the `permissions-audit` example, token-scope notes.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none observed.
- **What would need updating:** none. Polyfill pattern is out of scope for the fluent-builder playbook but documented for completeness.

---

## I. goja-treesitter — `~/code/wesen/go-go-golems/goja-treesitter`

### I1. `README.md` + `pkg/` — 🔴 (skeleton)

- **What I was researching:** a tree-sitter-backed DSL.
- **What I was looking for:** module exports, parse/query API.
- **Why I chose it:** listed per the ticket ("and others").
- **How I found it:** listed in `~/code/wesen/go-go-golems/`.
- **What I found useful:** confirmed it is a stub.
- **What I didn't find useful:** the README is only ASCII art; `pkg/` has only `doc.go` and `logcopter.go`.
- **What is out of date / what was wrong:** not wrong — just empty.
- **What would need updating:** either build the module or mark the repo as archived/placeholder so it is not mistaken for a working DSL.

---

## J. glazed — `…/glazed`

### J1. `README.md` (root) — 🟢

- **What I was researching:** the glazed CLI framework.
- **What I was looking for:** output formats, the command/field/middleware model.
- **Why I chose it:** the canonical Go-side builder ancestor.
- **How I found it:** listed in the working directory.
- **What I found useful:** the output-format examples, the field-flattening model.
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none.

### J2. `pkg/cmds/schema/schema.go`, `pkg/cmds/fields/definitions.go` — 🟢

- **What I was researching:** the `Schema`/`Section`/`Definition` functional-options model.
- **What I was looking for:** the typed-struct + functional-options pattern that should inspire the playbook.
- **Why I chose it:** the strongest Go-native pattern.
- **How I found it:** `rg -n 'type Schema struct'`.
- **What I found useful:** `Schema` over ordered map of `Section`; `Section` interface; `Definition` struct with `Option` funcs (`WithHelp`/`WithDefault`/`WithChoices`/`WithRequired`/`WithIsArgument`); `New(name, type, options...)`.
- **What I didn't find useful:** the TODO comment ("This is a pretty messy interface").
- **What is out of date / what was wrong:** the TODO suggests the author found it messy — senior researcher should note this when citing it as a model.
- **What would need updating:** none in-file. Playbook should cite the functional-options + typed-struct combination as the Go-side grammar, while noting the author's own TODO.

---

## K. go-emrichen — `~/code/wesen/go-go-golems/go-emrichen`

### K1. `README.md` — 🟢

- **What I was researching:** the tag-operator DSL.
- **What I was looking for:** composable operators, type safety, extensibility.
- **Why I chose it:** closest existing model for "composable grammar of operators extended with lambdas."
- **How I found it:** listed in `~/code/wesen/go-go-golems/`.
- **What I found useful:** `!Defaults`/`!If`/`!Var`/`!Format`/`!Loop` examples, "Strong Type Safety", "Extensible Design".
- **What I didn't find useful:** nothing.
- **What is out of date / what was wrong:** none.
- **What would need updating:** none. Conceptual model for the playbook's operator grammar.

### K2. `emrichen-spec.md`, `emrichen-in-practice.md` — 🟡 (not read in full)

- **What I was researching:** the operator spec and practical usage.
- **What I was looking for:** the full operator list and composition rules.
- **Why I chose it:** deeper than the README.
- **How I found it:** `ls` of the repo.
- **What I found useful:** confirmed both exist.
- **What I didn't find useful:** not read in full this pass.
- **What is out of date / what was wrong:** not assessed.
- **What would need updating:** senior researcher should read both and extract the operator-composition rules for the playbook's grammar section.

---

## L. External / web resources

No external web resources were fetched during this pass — all research was done from the local filesystem, which contains the full source of every repository. A web search was not needed because the DSLs are all in-house and the question is about *our* DSLs, not general goja usage.

If the senior researcher needs external context, candidates to fetch (not yet assessed):
- upstream goja `dop251/goja` and `goja_nodejs/require` READMEs — for the `require`/module-loader contract.
- Bleve mapping/query docs — for the domain API that goja-bleve wraps.

---

## M. Summary table

| Resource | Status | One-line verdict |
| --- | --- | --- |
| goja-bleve `README.md` | 🟢 | best fluent-builder example; drop phase framing |
| goja-bleve `api_types.go` / `api_mapping.go` / `module.go` | 🟢 | the typed-ref substrate to extract into a shared `fluent` package |
| goja-bleve `docs/faiss-xgoja-playbook.md` | 🟢 | model runbook for build constraints |
| goja-dbus `README.md` + `pkg/dbusgoja/*` | 🟢 | cleanest composable grammar; best lifecycle example |
| goja-dbus `GOJA-DBUS-DESIGN` ticket | 🟡 | read fully for the runtime-ownership rule |
| go-minitrace `README.md` | 🟢 | add a JS DSL section |
| go-minitrace `minitracejs/*` | 🟢 | Pattern B reference; `Validate()`/`Build()` discipline |
| go-minitrace `typescript.go` | 🟡 | expand to named builder types |
| widgetdsl `module.go` | 🟡 | the "not great" Pattern C; needs migration sketch |
| widgetdsl `grammar.go` | 🟡 | replace panic with `(v,error)`; replace magic tag |
| widgetdsl `typescript.go` | 🔴 | concrete evidence of weak compile-time types |
| widgetdsl tests | 🟢 | reliable API examples |
| RAGEVAL-UI-DSL ticket | 🟢 | RelatedFiles may be stale (old `internal/dsl/` paths) |
| RAGEVAL-UI-GRAMMAR design docs | 🟡 | read fully — antecedent of the grammar |
| rag-eval-scripting-expansion ticket | 🟡 | read fully for jsverbs/xgoja integration |
| go-go-goja `README.md` + `pkg/engine/` | 🟢 | runtime ownership model |
| go-go-goja `modules/uidsl/` | 🟡 | undocumented; add README/help |
| go-go-goja `modules/express/` | 🟢 | side-channel typed-lookup contrast |
| go-go-goja `pkg/tsgen/spec/types.go` | 🟢 | compile-time type substrate |
| go-go-goja `pkg/jsverbs/`, `pkg/jsdoc/` | 🟡 | assess as authoring paths |
| go-go-goja `pkg/xgoja/` | 🟡 | document provider packaging contract |
| goja-text `README.md` + `xgoja.yaml` | 🟢 | strongest docs; reference for paired help + config |
| goja-git `README.md` | 🟢 | Pattern E contrast example |
| goja-github-actions `README.md` | 🟢 | polyfill pattern; out of playbook scope |
| goja-treesitter `README.md` + `pkg/` | 🔴 | skeleton; mark as placeholder or build it |
| glazed `schema.go` / `definitions.go` | 🟢 | Go-side functional-options ancestor (note author TODO) |
| go-emrichen `README.md` | 🟢 | conceptual model for operator grammar |
| go-emrichen `emrichen-spec.md` / `-in-practice.md` | 🟡 | read fully for composition rules |

## N. Prioritised follow-ups for the senior researcher

1. **Read fully (deferred this pass):** RAGEVAL-UI-GRAMMAR design docs 01 & 02; rag-eval-scripting-expansion intern guide + design; GOJA-DBUS-DESIGN ticket; go-emrichen spec + practice; go-go-goja `jsverbs`/`jsdoc`/`xgoja` internals.
2. **Verify staleness:** RAGEVAL-UI-DSL RelatedFiles paths (`internal/dsl/widgetdsl/` → `pkg/widgetdsl/`); goja-git Go version; goja-bleve "Phase 8" claims; FAISS fork currency.
3. **Extract reusable primitives:** goja-bleve `refBase`/`refKind`/`getTypedRef[T]`/`newWrapper`/`mustSet` → candidate shared `fluent` package.
4. **Define the missing pieces:** validation discipline, TS emission rule, operator grammar + lambdas, lifecycle rule, runtime-ownership rule (see design-doc §6).
5. **Ratify or reject** the four decision records in design-doc §7.
