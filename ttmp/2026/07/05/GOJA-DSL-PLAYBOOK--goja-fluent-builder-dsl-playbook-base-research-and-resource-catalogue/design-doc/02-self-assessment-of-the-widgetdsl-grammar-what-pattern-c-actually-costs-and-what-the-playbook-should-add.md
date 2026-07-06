---
Title: Self-assessment of the widgetdsl grammar — what Pattern C actually costs, and what the playbook should add
Ticket: GOJA-DSL-PLAYBOOK
Status: active
Topics:
    - goja
    - dsl
    - fluent-builder
    - go
    - typescript
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/goja-bleve/pkg/api_types.go
      Note: getTypedRef substrate proposed for widgetdsl markers
    - Path: ../../../../../../../../../2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/builders.go
      Note: The typed-spec + lambda-configurator architecture the v2 sketch adopts
    - Path: ../../../../../../../../../2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/typescript.go
      Note: Precise this-typed builder declarations — the .d.ts target
    - Path: packages/rag-evaluation-site/src/widgets/ir.ts
      Note: The precise TS types that already exist and are not consumed by declaration generation
    - Path: pkg/widgetdsl/grammar.go
      Note: The code under assessment — written by the author of this document
    - Path: pkg/widgetdsl/typescript.go
      Note: Where three existing type sources get discarded into Props = Record<string, any>
ExternalSources: []
Summary: 'The author of pkg/widgetdsl/grammar.go assesses that work against the catalogue in design-doc 01: five empirically verified silent-failure modes, a root-cause analysis (the IR wire format was conflated with the authoring API), a concrete v2 sketch that keeps the grammar''s language design but moves it onto typed Go specs with lambda-configurator builders and strict option decoding, and eleven proposed additions to the catalogue — including a failure-mode taxonomy where silent-ignore ranks worse than panic, a tree-shaped-data axis the pattern taxonomy currently lacks, and the dormant type sources (ir.ts, widget manifests) that make precise .d.ts generation cheaper than the catalogue assumes.'
LastUpdated: 2026-07-05T13:00:00-04:00
WhatFor: Feed the senior-researcher playbook pass with an insider post-mortem of the weakest pattern and concrete additions the catalogue is missing.
WhenToUse: Read after design-doc 01, before writing the playbook; also the reference for any widgetdsl v2 migration.
---


# Self-assessment of the widgetdsl grammar — what Pattern C actually costs, and what the playbook should add

## 0. Position of this document

I wrote `pkg/widgetdsl/grammar.go` (the `data.dsl` verbs: `f.*`, `schema`, `record`, `collection`, plus `ui.section`) during RAGEVAL-UI-GRAMMAR, on top of the pre-existing map-IR module machinery in `module.go`. Design-doc 01 of this ticket classifies that work as Pattern C and calls it "the not-great one." Having now read the catalogue and the reference implementations it points to (goja-bleve's typed refs, codesign's lambda configurators and typed specs, geppetto's DTS parity test, go-minitrace's validation terminals), I agree with the classification, and this document explains *why* the assessment is correct from the inside: what the grammar got right at the language level, where the substrate fails concretely (with empirically verified failure cases, not hypotheticals), what a v2 on a sound substrate looks like, and what the catalogue itself should add before the playbook pass.

The short version: **the grammar's language design is worth keeping; its implementation substrate is not.** Roles, verbs, arrangements, URL-backed selection, and elision rules survived contact with a real page and measurably fixed it. But I built that language on `map[string]any` because the house idiom was map-IR, and the honest accounting below shows that idiom fails in exactly the ways that matter most for a DSL whose authors are increasingly LLM agents: silently.

## 1. What the grammar got right (and should survive any migration)

Fairness first, because the v2 sketch in §5 keeps all of this:

1. **Intent-level verbs over component nouns.** `collection(rows, {verb, arrange, select…})` replaced a 40-line hand-rolled panel-per-record function. The consuming admin page went from 5,611 px / 21 nested panels to 3,496 px / 5 panels / zero nesting. The *vocabulary* was the right fix; that conclusion stands regardless of substrate.
2. **Field roles rather than types.** `key/primary/short/prose/status/measure…` deciding summary rendering, editor control, and elision is a sound scale-like abstraction, directly analogous to what makes `paletteStyleSet` work in `context_window.dsl`.
3. **URL-backed selection with explicit value threading.** `urlParam(name, value)` respects the no-client-state architecture and the fact that the Go DSL never sees the request. The *concept* is right even though the marker object carrying it is untyped (§3.2).
4. **Compile-to-existing-IR.** Shipping the grammar without touching the renderer was the correct sequencing decision; it is also fully compatible with a typed substrate (§4).
5. **Order preservation across the boundary.** `schema()` iterating `Object.Keys()` before export, because `Export()` to `map[string]any` destroys insertion order, is a real contribution the playbook should absorb as a hazard rule.
6. **Grammar expansion tests.** `grammar_test.go` snapshots verb expansion, which is more than most of the map-IR family has.

## 2. The root cause: IR-as-API instead of IR-as-output

The catalogue's Pattern C description lists symptoms (untyped maps, panic validation, weak declarations). The underlying cause deserves to be named because it is the load-bearing lesson for the playbook:

**Widget IR is a wire format, and widgetdsl uses the wire format as the authoring API.** Every helper returns the serialized artifact directly. There is no intermediate typed representation on the Go side at all — the maps that helpers build are the exact JSON that ships to the renderer.

This one decision produces every downstream weakness mechanically. There is nothing for `getTypedRef[T]` to extract because there are no refs. Validation has nowhere to accumulate because there is no builder carrying state. TypeScript declarations cannot narrow because the API's real signature *is* "any JSON-shaped map." Unknown options cannot be rejected because a map merge has no notion of known keys.

Contrast with codesign (Pattern F), which is the same architecture done right: builder methods accumulate into **typed Go spec structs** (`codesignspec.RunSpec`), methods return `(*goja.Object, error)` so misuse throws a JS exception with a precise message, `.validate()` returns structured `{code, path, message}` issues, and serialization to JSON happens only at the terminals. Builders in, wire format out.

The irony specific to widgetdsl: the typed spec layer **already exists** — on the wrong side of the wire. `packages/rag-evaluation-site/src/widgets/ir.ts` defines precise per-widget interfaces (`SectionBlockWidgetProps`, `DataTableColumnSpec`, `ActionSpec` as a discriminated union), and 80+ `*.widget.yaml` manifests name each widget's props type, helper, module, and slots. Nothing on the Go side consumes either. The system has three sources of type truth (ir.ts, manifests, Go helper maps) and the declaration generator ignores all of them, emitting `Props = Record<string, any>`.

## 3. Empirical failure-mode audit

Claims about DSL ergonomics are cheap, so I ran the failure cases against the built go-go-course binary (which embeds the current grammar). All five are **silent** — no error, no warning, wrong output:

| # | Author writes | What happens | Verified result |
|---|---|---|---|
| 1 | `arrange: "masterdetail"` (typo for `master-detail`) | String compared with `==`; anything else falls through to table | Summary table renders, detail editor silently absent |
| 2 | `data.f.primary({ maxLenght: 99 })` | Option merged into the field spec, never read | Typo kept in the map; no length limit applied |
| 3 | `submit: data.urlParam("x","a")` (wrong marker) | `copyIfPresent` finds no `formAction` key | FormPanel renders with no form action; save button does nothing |
| 4 | `verb: "browse"` | Unknown verb treated as `"show"` | Read-only rendering, no error |
| 5 | `ui.section("T", { level: 7 })` | Passed through to the widget; `styles.level7` is undefined | Renders with default label styling |

Add the two failure classes already on the record: the *panic* path (`record()` without a schema panics with a Go error — in a request-scoped page builder this is a 500 for the whole page), and the **confirm-encoding bug** from RAGEVAL-UI-GRAMMAR Step 3, where the `${row.title}` interpolation mini-language URL-encoded human-facing confirm text (`Delete “Smoke%20item”?`). That bug is Pattern C in miniature: an untyped embedded template language whose destination semantics (URL vs. prose) were invisible to every layer until a value with a space flowed through at runtime.

The taxonomy this suggests — and which the catalogue's gap analysis currently lacks — has *silent-ignore* as its own category, ranked **worse than panic**:

1. **Silent-ignore** (cases 1–5): wrong output, no signal. The author discovers the defect visually, or never.
2. **Panic/throw-at-build**: hostile but at least loud and located.
3. **Accumulated validation at terminals**: the target (`{code, path, message}` issues à la codesign/minitrace).

Silent-ignore matters double for this codebase because the page builders are substantially authored by LLM agents. An agent's only feedback channels are the `.d.ts` (which says `any`, so its editor/typecheck pass catches nothing) and runtime errors (which don't fire). Cases 1–5 are precisely the defects agents produce — plausible near-miss strings — and the current substrate absorbs all of them without a sound.

## 4. Ranked weaknesses, with the fix each one implies

1. **Silent option and enum failures** (§3). Fix: decode every options bag into a typed Go struct with unknown-key rejection, and validate enums with suggestion messages (`unknown arrange "masterdetail"; did you mean "master-detail"?`). This is independent of any builder migration and is the cheapest, highest-value change available.
2. **Untyped markers.** `urlParam`/`formPost`/`schema` results are shape-tagged maps (`__ragSchema` magic string, duck-typed `{param,value}`). Fix: hidden-key typed refs (bleve's `attachRef`/`getTypedRef[T]`) so passing the wrong marker is a typed error naming both kinds.
3. **Declaration generation discards existing types.** `typescript.go` emits `RawDTS` strings with `Props` for all ~70 helpers, while `tsgen/spec` supports full named `TypeRef`s and ir.ts already contains the interfaces. Fix: make the widget manifests load-bearing — generate per-helper parameter types (and Go-side option decoders) from a single schema source, and add a geppetto-style **DTS parity test**. Note that during RAGEVAL-UI-GRAMMAR I added `section`/`f`/`record`/`collection` declarations by hand-appending strings; nothing would have caught it if I had forgotten, which is exactly the drift the parity test exists to prevent.
4. **Panic as the only loud path, in a request-scoped DSL.** Page builders run per request; a panic is a 500. A UI DSL has a better option unavailable to config DSLs: **render the validation failure**. Accumulate issues during building; if the terminal `page()` finds errors, emit an error-panel IR node listing `{path, message}` instead of failing the request. Authors (and agents) see the defect in the page itself.
5. **Options-bag ergonomics at the complex end.** One giant nested literal for `collection` means no discoverability, no partial reuse, and deep-nesting errors. Codesign's lambda configurators + `.use(fragment)` solve exactly this. The v2 sketch below adopts them for grammar-level objects — while §5.3 argues plain props bags should *stay* at the leaf-widget level.
6. **Two authoring layers with implicit compilation.** Raw component helpers and grammar verbs coexist with no stated rule for when each applies and no shared types between them. A migration must either type both or explicitly bless the raw layer as the escape hatch.

## 5. What v2 looks like: the same grammar on a typed substrate

### 5.1 Typed specs mirroring ir.ts, builders in front

```go
// Go side: typed spec structs (the missing mirror of ir.ts)
type FieldSpec struct {
    Name, Role, Label, Width, Placeholder, Hint string
    Required, ReadOnly bool
    MaxLength, Rows int
}
type CollectionSpec struct {
    Schema  *SchemaSpec
    Verb    Verb        // typed enum, validated at parse
    Arrange Arrangement
    Select  *URLParamSpec
    Submit  *FormPostSpec
    Reorder, Remove *ActionSpec
    issues  []Issue     // accumulated, not thrown
}
func (c *CollectionSpec) Compile() (WidgetNode, []Issue)  // serialize to IR here, and only here
```

Builders wrap these structs behind hidden-key refs; methods return `(obj, error)`; `Compile()`/`page()` are the terminals. The IR stays byte-identical — the renderer, the existing pages, and the phase-α strategy are untouched. This is the codesign architecture with ir.ts as the spec vocabulary.

### 5.2 Authoring surface: lambda configurators for grammar objects

```js
const agenda = data.collection(rows, c => c
  .schema(agendaSchema)
  .edit(e => e
    .select("agenda", query.agenda)
    .submit("/settings/agenda-item")
    .reorder(a => a.server("admin-reorder-course-agenda"))
    .remove(a => a.server("admin-delete-agenda-item")
                   .confirm(t => t.text("Delete agenda item "), t.value("row.title"), t.text("?")))
    .create("New agenda item"))
  .arrange("master-detail"))
```

Details worth noting: `edit(…)` replaces the stringly `verb:` option with a method whose existence the `.d.ts` can state (`this`-typed, as codesign's declarations do); the confirm template becomes structured parts instead of a `${…}` string, which makes the URL-vs-prose encoding decision a property of the part, killing the §3 bug class at the root; `.use(fragment)` gives page files reusable column sets and action groups. Every callsite error — wrong method, wrong argument type, unknown enum — is now either a TS-visible signature mismatch or a thrown Go error with a path.

### 5.3 What should *not* become fluent

The playbook should resist totalizing the builder pattern, and widgetdsl is the evidence case: it builds **trees**. `panel(props, ...children)` nesting is the natural notation for trees — hyperscript (Pattern D) is not a lesser pattern here but the correct one for structure. Fluent chains are the correct notation for *stateful or multi-step configuration* (schemas, collections, actions, queries). The v2 rule: **hyperscript nesting for structure, builders for specs, strict decoding for leaf props bags.** A `caption({tone: "muted"}, "text")` call does not need a builder; it needs its options decoded against a generated type with unknown-key rejection.

### 5.4 Migration path (incremental, no flag-day)

1. Strict option decoding + enum validation inside the existing helpers (fix #1; no API change; existing pages keep working or reveal latent typos — both outcomes are wins).
2. Typed refs for the existing markers (`schema`, `urlParam`, `formPost`) and grammar returns (fix #2; the JS surface is unchanged).
3. Manifest-driven codegen for declarations + decoders, DTS parity test (fix #3).
4. Builder/configurator surface for `collection`/`record`/`action` as the *blessed* API, options-bag forms kept as deprecated aliases (fix #5); grammar callsites are few — one consumer page today — so this is cheap *now* and gets more expensive with every new page.
5. Error-panel rendering for accumulated issues (fix #4), which requires one new widget and a host-side convention.

Steps 1–3 are independent of the playbook's final builder shape and could land immediately; step 4 should wait for the playbook so widgetdsl v2 *is* the playbook's reference implementation rather than a sixth ad-hoc pattern.

## 6. Proposed additions to the catalogue (design-doc 01)

Concrete gaps in the current document, in rough priority order:

1. **Failure-mode taxonomy with silent-ignore as a first-class category.** §6's gap table has "validation" rows but no axis for *unknown-input handling*. The §3 audit shows silent-ignore is Pattern C's real cost; panic is almost incidental. Each catalogued DSL should get a "what happens on a typo'd option / wrong enum / wrong handle" row — bleve and codesign pass, express partially, widgetdsl fails five ways.
2. **A tree-shaped-data axis in the pattern taxonomy.** Every Pattern A/B/F exemplar builds flat-ish configs (mappings, run specs, queries). widgetdsl builds trees, and Pattern D exists precisely because trees want nesting notation. The playbook needs an explicit answer for children — otherwise it will prescribe fluent chains for structure and produce something worse than what exists. Proposed rule: D for structure, F for specs (§5.3).
3. **Split Pattern C into wire-format and authoring-API.** "Map IR" conflates two decisions: JSON-serializable output (fine, shared by codesign's `toSpec()`) and maps-as-API (the defect). The playbook rule is then stateable in four words: *builders in, IR out.*
4. **The dormant type sources.** The catalogue says widgetdsl has "weak compile-time types" but not that the precise types already exist twice (ir.ts interfaces, 80+ `*.widget.yaml` manifests with `props:` type names) with zero programmatic consumers. This changes the cost estimate for fixing declaration generation from "author everything" to "wire up what exists," and single-source-of-truth codegen (manifest → Go decoder + TS decl) deserves to be a playbook pattern.
5. **Embedded mini-languages as typed surface.** The `${row.x}` interpolation language inside ActionSpecs is an unmentioned DSL-within-the-DSL; its encode bug is documented evidence of what untyped template strings cost. The playbook should require template parts or typed template objects wherever a string crosses between URL and human-text destinations.
6. **Agent authorship as a design constraint.** These DSLs are written by LLM agents at least as often as by humans. That inverts some classic tradeoffs: precise `.d.ts` and rich error messages are the *only* feedback channels an agent gets, discoverability-by-chaining matters more than brevity, and silent failure is maximally expensive because the agent's visual check is weak. Worth a paragraph in the playbook's goals.
7. **Request-scoped error policy.** The validation decision record (accumulate + `(value, error)`) is right for config DSLs but incomplete for UI DSLs, which have a rendering channel: a page builder's accumulated issues can render as an error panel rather than failing the request. Add as a variant to the validation decision.
8. **The property-order hazard.** goja preserves insertion order via `Object.Keys()` but not through `Export()` to Go maps. Any DSL whose semantics depend on object-literal order (schemas, column sets) must capture order at the boundary. Belongs in the playbook's hazards section; currently only discoverable by reading `grammar.go` comments.
9. **Per-DSL scoring rubric.** The catalogue describes; the playbook pass would benefit from scores. Suggested axes: typo safety, wrong-handle safety, error quality (message + path), `.d.ts` precision, tree ergonomics, lambda support, lifecycle discipline, cost-to-add-a-module. Even coarse 0/1/2 scores would make the §7 decision records checkable.
10. **Answers to open questions 3 and 9 from the widgetdsl experience.** Q3 (eager vs. lazy validation): eager *shape* checks at each call (unknown keys, enum membership — these need no cross-field context), lazy *semantic* validation at terminals (cross-field, cardinality). Q9 (migrate `data.dsl`?): yes, and specifically *now* — the grammar verbs have exactly one consumer page, so the API can still change cheaply; the raw component helpers can follow later or remain the documented escape hatch.
11. **One survey correction.** §4.4 lists `ui.section("Title", { collapsible: true }, …)` — there is no `collapsible` option; the real options are `level`, `anchor`, `caption`, `actions`, `rule`, `density`, `divider`. Fittingly, the current implementation would accept the wrong example silently, which is case 2 of §3.

## 7. Verdict

The RAGEVAL-UI-GRAMMAR work answered the *language* question well — what should page authors be able to say — and answered the *substrate* question by defaulting to the house idiom without examining it. Reading this catalogue makes the cost of that default measurable: five silent failure modes in the grammar's own surface, a shipped encoding bug in its template mini-language, and declarations that give the system's primary authors (agents) nothing to check against. The pattern the playbook should bless — typed Go specs, configurator builders with `(obj, error)` methods, strict option decoding, accumulated validation with a rendering channel, manifest-driven codegen with a parity test — already exists piecewise across codesign, bleve, minitrace, and geppetto. widgetdsl should become its first migration, starting with the three steps (§5.4, 1–3) that require no API decision at all.
