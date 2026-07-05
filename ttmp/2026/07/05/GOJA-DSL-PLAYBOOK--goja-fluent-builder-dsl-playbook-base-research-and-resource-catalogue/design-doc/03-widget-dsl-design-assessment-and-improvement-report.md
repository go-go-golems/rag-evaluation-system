---
Title: Widget DSL design assessment and improvement report
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
      Note: Comparator for typed Go reference wrappers
    - Path: ../../../../../../../../../2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/builders.go
      Note: Comparator for lambda configurators and fragment composition
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js
      Note: First consumer proving the vocabulary but also showing current API shape
    - Path: packages/rag-evaluation-site/src/widgets/actions.ts
      Note: Evidence for string-template action interpolation and row-action context
    - Path: pkg/widgetdsl/grammar.go
      Note: Primary evidence for the Phase 0 widgetdsl grammar and its map/panic/string-switch limitations
    - Path: pkg/widgetdsl/module.go
      Note: Evidence for module organization, promoted helpers, and recipes that still bypass schemas+marks
    - Path: pkg/widgetdsl/typescript.go
      Note: Evidence for open-ended generated declarations and missing compile-time contract
ExternalSources: []
Summary: "Independent assessment of the RAGEVAL-UI-GRAMMAR work: what was strong, what remains weak, and what the Goja DSL playbook should add before using widgetdsl as a durable DSL direction."
LastUpdated: 2026-07-05T13:20:00-04:00
WhatFor: "Use as an external review of the existing self-assessment and current widgetdsl grammar design."
WhenToUse: "Read after design-doc 01 and the colleague-authored design-doc 02 self-assessment."
---


# Widget DSL design assessment and improvement report

## Executive summary

The colleague's RAGEVAL-UI-GRAMMAR work was a good **product/UI rescue** and a weak **DSL architecture endpoint**. The measured audit, the “components without intent” diagnosis, and the first vertical slice on the CMS agenda page are all valuable. They demonstrate that an intent-level vocabulary can reduce a pathological page from nested hand-authored boxes into a summary table plus one editor. That part should be preserved.

The problem is that the implementation stops at a shallow macro layer: `data.schema`, `data.record`, `data.collection`, and `ui.section` still return untyped Widget IR maps; validation is mostly panic-at-construction; TypeScript declarations are open-ended; arrangements and roles are strings; marks are not first-class; lambdas and fragment composition are absent; and the authored “grammar” has no typed intermediate representation, no reusable builder substrate, and no compile-time contract. In the taxonomy from design-doc 01, this is still **Pattern C: map IR + loose helper functions**, with a few better vocabulary words on top.

My recommendation is not to discard the work. Treat it as **Phase 0 evidence**: it proves which intentions the authoring layer needs (`section`, `schema`, `record`, `collection`, URL selection, native form-post binding), but it should not be promoted as the canonical Goja fluent-builder DSL model. The playbook should explicitly say: the current widgetdsl grammar is a useful macro prototype; the next design should re-express those concepts with typed Go builders, named schemas, typed marks/arrangements, composable lambda configurators, accumulated validation, and precise generated TypeScript declarations.

## Scope and evidence used

I read and compared these materials:

- The base catalogue: `design-doc/01-goja-dsl-catalogue-and-base-research.md` in this ticket.
- The Obsidian article: `/home/manuel/code/wesen/go-go-golems/go-go-parc/Projects/2026/07/05/ARTICLE - Widget DSL Grammar - Designing an Intent-Level UI Authoring Layer for a Widget IR System.md`.
- The RAGEVAL-UI-GRAMMAR ticket:
  - `design-doc/01-composable-ui-grammar-for-the-widget-dsls-analysis-of-the-cms-admin-page-and-brainstormed-alternatives.md`.
  - `design-doc/02-dsl-api-sketch-grammar-verbs-in-data-dsl-and-ui-dsl-module-reorganization-no-new-module.md`.
  - `reference/01-investigation-diary.md`.
- The current implementation in `rag-evaluation-system`:
  - `pkg/widgetdsl/grammar.go`.
  - `pkg/widgetdsl/module.go`.
  - `pkg/widgetdsl/typescript.go`.
  - `pkg/widgetdsl/grammar_test.go`.
  - `packages/rag-evaluation-site/src/widgets/ir.ts`.
  - `packages/rag-evaluation-site/src/widgets/actions.ts`.
  - `packages/rag-evaluation-site/src/components/molecules/DataTable/*`.
  - `packages/rag-evaluation-site/src/components/layout/{SectionBlock,FieldGrid}/*`.
- The first consumer in `go-go-course`:
  - `cmd/go-go-course/lib/pages/admin-course-cms.js`.
  - `cmd/go-go-course/server.js`.
  - `cmd/go-go-course/lib/course-metadata-service.js`.
- Stronger comparator DSLs from the catalogue:
  - `goja-bleve/pkg/api_types.go` and `api_mapping.go`.
  - `researchctl/pkg/gojamodules/codesign/builders.go` and `typescript.go`.
  - `researchctl/pkg/gojamodules/researchctl/builders.go` and `module.go`.

## Verdict: what they got right

### 1. The diagnostic method is strong

The RAGEVAL-UI-GRAMMAR docs are strongest where they stay empirical: the CMS page was measured by height, panel count, nesting depth, and form-row count; the resulting diagnosis is concrete. The docs correctly identify that the old DSL exposed component nouns (`panel`, `formRow`, `dataTable`) but not authoring intentions such as “edit this collection” or “show this record.”

That diagnosis should survive into the playbook. The best line of argument is: if a DSL only exposes components, authors will encode intent manually in page code, and every manually unrolled collection becomes visually and semantically inconsistent.

### 2. The “no grammar.dsl” decision is right

Design-doc 02 rejects a sixth `grammar.dsl` module and keeps the grammar in existing namespaces: `data.dsl` owns record/collection/schema; `ui.dsl` owns sectioning. The module maps in `pkg/widgetdsl/module.go` support this mechanically: `uiHelpers` collects generic helpers (`breadcrumbs`, `emptyState`, `fieldGrid`, `pagination`, `tag`, `tileGrid`, `uploadDropArea`, etc.) at lines 34-79, while `dataHelpers` is still small (`dataTable` only) at lines 81-83.

This is good API hygiene. A separate “good way” module would institutionalize two authoring styles. The target should be fewer, stronger module concepts, not another namespace.

### 3. The first vertical slice proves the intent vocabulary

The consumer page in `go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js` now uses `ui.section` for page structure and `dataDsl.collection(..., { verb: "edit", arrange: "master-detail" })` for agenda editing. The agenda schema lives at lines 95-103; the collection call is at lines 112-129. This is substantially more legible than hand-unrolled `agenda.map(...)` panels.

The server-side work also matches the UI grammar direction: per-record agenda saves are handled by `/settings/agenda-item` at `server.js` lines 362-379, and row-level delete/reorder actions use the DataTable action context at lines 461-495. The persistence helpers in `course-metadata-service.js` lines 314-354 correctly operate on the effective agenda rather than assuming an override file already exists.

### 4. The property-order fix is an important embedded-DSL lesson

`data.schema` captures `Object.Keys()` before exporting to Go maps (`pkg/widgetdsl/grammar.go` lines 77-93). This is exactly the kind of boundary detail the playbook should include: if field order is semantic, capture it before `Export()` destroys object-literal order.

### 5. The compile-to-existing-IR strategy was a sensible first step

`grammar.go` states the Phase 0 strategy explicitly: grammar calls compile to existing components (`SectionBlock`, `FieldGrid`, `FormPanel`, `FormRow`, `DataTable`, `Stack`, `Inline`, `Button`, `Caption`) at lines 3-12. That was the correct delivery strategy for a UI cleanup: minimal renderer churn, fast proof, and a real page exercising it.

## Where the work is not good enough

### 1. It is still a map macro layer, not a type-safe DSL

The implementation returns raw `map[string]any` Widget IR. `schemaCtor` returns a map tagged with `"__ragSchema": true` (`grammar.go` line 93). `schemaFields` then trusts a `map[string]any` shape and silently ignores malformed entries (`grammar.go` lines 96-105). `recordVerb` and `collectionVerb` panic if no usable schema is found (`grammar.go` lines 159-162 and 274-277).

This is better vocabulary, but it is not the fluent-builder model requested in the GOJA-DSL-PLAYBOOK ticket. Compare goja-bleve: it attaches a non-enumerable hidden Go reference (`api_types.go` lines 105-114), extracts typed refs with `getTypedRef[T]` (`api_types.go` lines 131-141), and returns wrapper objects with explicit kinds (`api_types.go` lines 143-148). The widget DSL has none of that. A random object with a `fields` array can masquerade as a schema once exported; the runtime cannot distinguish a real schema handle from a shape-compatible map.

**Improvement:** introduce typed Go builder handles for schemas, fields, collections, marks, and actions. The terminal should be `.build()` or `.toIR()`; raw map export should be the last step, not the authoring substrate.

### 2. Validation is immediate, partial, and panic-driven

The current code validates only a few boundary cases and does so with `panic(r.vm.NewGoError(...))`:

- `schema(fields)` requires an object and panics otherwise (`grammar.go` lines 79-82).
- schema fields must contain a `role` key, but the role value is not checked against the declared role set (`grammar.go` lines 85-89).
- `record` and `collection` panic when no schema fields are found (`grammar.go` lines 159-162 and 274-277).

There is no `validate()` result, no accumulated issue list, no path-addressed diagnostics, and no distinction between authoring errors, recoverable validation issues, and renderer-time missing data. This falls short of the playbook's goal and of better existing patterns. Codesign exposes `.validate()` and `.toSpec()` terminals on the builder (`codesign/builders.go` lines 76-79) and a `ValidationResult` type in its declarations (`codesign/typescript.go` lines 11-12). Go-minitrace, as catalogued in design-doc 01, has the same validation-result discipline.

**Improvement:** a widget grammar should accumulate issues and expose both `validate()` and `toIR()`/`build()` terminals. Misuse should return `error` from Goja-bound functions wherever possible, not panic.

### 3. TypeScript support actively hides mistakes

`pkg/widgetdsl/typescript.go` says the IR is “intentionally represented as JSON-like data” and keeps props open-ended (lines 10-13). The generated declarations define `WidgetNode { kind: string; [key: string]: any }`, `WidgetAction { kind: string; [key: string]: any }`, and `Props = Record<string, any>` (lines 20-31). The data grammar adds `FieldSpec { role: string; [key: string]: any }` and `record(values: Props, options: Props)` / `collection(rows: Props[], options: Props)` (lines 78-100).

That means TypeScript cannot catch:

- misspelled field roles or option keys,
- invalid `verb` / `arrange` values,
- missing `schema`,
- wrong action shape,
- `urlParam` without the current value,
- a mark expecting a `media` role on a schema that does not provide one.

By contrast, codesign's DTS models `FragmentFn<T>`, `RunSpecBuilder`, `TopologyBuilder`, `WorkloadBuilder`, and `MetricsBuilder` precisely (`codesign/typescript.go` lines 27-33). That is much closer to the requested “compile-time types from generated declarations.”

**Improvement:** the document should not merely say “TS declarations need improvement.” It should propose concrete named types: `FieldRole`, `FieldSpec<R>`, `Schema<TRecord>`, `CollectionOptions<TRecord, TSchema>`, `Arrangement<TRecord>`, `Mark<TRecord>`, `UrlSelection<TKey>`, `FormPostBinding`, and discriminated `ActionSpec` types.

### 4. Roles are overloaded and under-specified

The field role set is hard-coded at `grammar.go` lines 25-38: `key`, `primary`, `short`, `prose`, `count`, `size`, `measure`, `date`, `status`, `tags`, `media`, `href`. The comments say roles drive summary rendering, editor controls, and read-only views. That combines at least four separate axes:

1. identity (`key`),
2. semantic kind (`status`, `href`, `media`),
3. visual/editor density (`short`, `prose`),
4. aggregation/measurement (`count`, `size`, `measure`).

The code then uses role-specific hard-coded rules:

- `gridableRoles` controls field-grid batching (`grammar.go` lines 40-43 and 213-219).
- `prose` and `media` are elided from collection summary tables (`grammar.go` lines 321-324).
- numeric-ish roles become number cells (`grammar.go` lines 332-334).
- `status` becomes a status cell (`grammar.go` lines 334-335).
- only `prose` gets a `TextareaInput` (`grammar.go` lines 240-246).

This is a useful heuristic for the agenda page, but it is not enough for a reusable widget grammar. A field's storage type, display role, editor control, summary mark, validation rules, and layout preferences should be separable.

**Improvement:** split field declarations into typed facets, for example:

```js
f.string("title").primary().required().maxLength(160)
f.string("description").editor("textarea", { rows: 4 }).summary("elide")
f.enum("status", ["draft", "published"]).mark(data.marks.status())
f.number("tokens").measure({ unit: "tokens", scale: data.scale.linear({ domain: [0, limit] }) })
f.href("url").summary(data.cells.link({ label: "title" }))
```

The playbook can still allow shortcuts (`f.primary`, `f.prose`), but the underlying model should not collapse everything into a single `role` string.

### 5. “Arrangement” is a string switch, not a composable grammar

`data.collection` supports `arrange = "table" | "master-detail"` in practice (`grammar.go` lines 278-292). Everything else in the design-doc 02 vocabulary (`tiles`, `cards`, `disclosure`, domain marks, multi-view marks) remains aspirational. The current `collectionTable` and `collectionDetail` are not extension points; they are fixed compiler functions (`grammar.go` lines 319-367 and 413-472).

This is the biggest DSL-design gap. A grammar should let authors compose arrangement operators, not select one string from a closed switch. Codesign shows the pattern we should borrow: `runSpec(...).topology(fn).workload(fn).metrics(fn).use(fragment)` (`codesign/builders.go` lines 46-75), and sub-builders also support `.use(fragment)` (`codesign/builders.go` lines 172-177 and 356-360). The shared callback applicator is small and explicit (`codesign/builders.go` lines 365-374).

**Improvement:** turn arrangements and marks into typed builder fragments:

```js
const agendaList = data.arrangements.masterDetail(a => a
  .summary(data.marks.table(t => t.elide("description").rowSelect("agenda")))
  .detail(data.marks.recordForm(f => f.submit(data.formPost("/settings/agenda-item"))))
  .actions(actions => actions.reorder("admin-reorder-course-agenda").remove("admin-delete-agenda-item"))
)

data.collection("agenda", agenda)
  .schema(agendaSchema)
  .edit()
  .arrange(agendaList)
  .use(adminDefaults)
  .toIR()
```

This would preserve the vocabulary discovered by RAGEVAL-UI-GRAMMAR while moving it into the proven lambda-configurator + fragment composition model.

### 6. Actions and URL state are pragmatic but too stringly typed

The current selection model is `data.urlParam(param, value)` (`grammar.go` lines 50-56), and `collectionTable` creates a navigate target with a string template (`?${param}=${row.keyField}`) at lines 358-362. The action dispatcher interpolates `${...}` templates and URL-encodes values by default (`actions.ts` lines 71-80 and 115-130). Confirm prompts disable encoding (`actions.ts` lines 29-33), which fixed a real bug.

This works, but it spreads semantics across string templates, action context shapes, and server conventions. `DataTable.widget.tsx` dispatches `row`, `rowKey`, and `componentType` as context (lines 28-35), but neither the Go DSL nor TypeScript declarations make this context typed.

**Improvement:** model action context and bindings explicitly:

```ts
interface RowActionContext<T> { row: T; rowKey: string; componentType: "DataTable" }
interface UrlParamSelection<TKey> { kind: "urlParam"; param: string; current: TKey | null }
interface ServerRowAction<T> { kind: "server"; name: string; payload?: JsonObject | ((ctx: RowActionContext<T>) => JsonObject) }
```

In Goja, action builders should produce typed action handles rather than arbitrary maps. Templates can remain as a low-level escape hatch.

### 7. The “domain modules shrink to schemas + marks” claim is not yet implemented

Design-doc 02 has a good idea: `context_window.dsl`, `course.dsl`, and `cms.dsl` should export schemas and marks, while `data.dsl` owns the grammar sentence. The implementation does not do this yet. `cms.recipes.mediaLibrary` still compiles directly to `MediaLibraryPanel` (`module.go` lines 784-814), `articleList` still compiles directly to `ArticleListPanel` (`module.go` lines 816-843), and `masterDetailTable` remains an older recipe path (`module.go` lines 845-871). `context_window.dsl` still owns style helpers and bespoke snapshot helpers (`module.go` lines 339-421).

The article's “what remains” section correctly names this gap, but the catalogue document should make it more prominent. Until a domain recipe is re-expressed as `data.collection(..., { arrange: cms.mark.assetTiles })`, the architecture is not proven.

**Improvement:** the next document revision should include one fully worked migration of `cms.recipes.mediaLibrary` or `contextDiagram` into schemas + marks. Without that, “schemas + marks” is still a slogan.

### 8. The document overstates “grammar” and understates language design criteria

The Obsidian article is good narrative, but it is too celebratory for a crucial design decision. It describes what shipped and why it improved one page. It does not sufficiently evaluate whether the API is a durable language:

- What are the grammar's algebraic primitives?
- What composes with what?
- Which concepts are typed entities versus plain objects?
- What are the invariants?
- Which errors are caught at author time, runtime validation time, and render time?
- How do schemas and marks evolve across packages?
- How do generated types stay in sync with runtime exports?
- What is the migration path from recipes to grammar sentences?
- What makes this better than a handful of higher-level recipes?

These questions should be explicitly added to design-doc 01 or to a follow-up playbook chapter.

## What I would add to the base catalogue document

The catalogue currently does a good inventory and correctly classifies widgetdsl as Pattern C. To make it useful for the next decision, I would add a new section immediately after §4.4 or after §6:

### Proposed addition: “Widgetdsl self-assessment against the fluent-builder bar”

Add a table like this:

| Criterion | Current widgetdsl | Stronger model | Required next step |
| --- | --- | --- | --- |
| Runtime handles | `map[string]any`, `__ragSchema` tag | goja-bleve hidden typed refs | Typed schema/field/mark/action builders |
| Composition | strings: `verb`, `arrange` | codesign lambda configurators + `.use()` | Arrangement/mark fragments |
| Validation | panic + partial checks | `validate()` result + terminal errors | Accumulated `ValidationIssue` list |
| TypeScript | `Props = Record<string, any>` | codesign named builder interfaces | Precise declarations + parity test |
| Domain extension | recipes still direct to panels | schemas + marks | Rewrite one domain recipe as grammar sentence |
| Action context | string templates + duck-typed row context | typed callbacks/actions | typed `RowActionContext<T>` and bindings |
| IR boundary | authoring layer is the IR map | typed intent IR terminal to Widget IR | introduce `WidgetIntentSpec`/`toIR()` |

### Proposed addition: “Do not treat Phase 0 widgetdsl as the playbook pattern”

State this plainly:

> The current widgetdsl grammar is evidence for the vocabulary, not evidence for the implementation substrate. It should inspire the noun/verb set (`schema`, `record`, `collection`, `section`, selection, submit bindings), but future DSLs should not copy its raw map representation, `Props` declarations, or panic validation.

### Proposed addition: “Canonical target API sketch”

Include one aspirational but implementable API sketch using the stronger patterns:

```js
const data = require("data.dsl")
const ui = require("ui.dsl")

const agenda = data.schema("Agenda", s => s
  .key("id", { label: "ID" })
  .field("number", data.types.string()).role("time").summary(data.cells.caption())
  .field("duration", data.types.string()).summary(data.cells.field())
  .field("title", data.types.string()).primary().required().maxLength(160)
  .field("description", data.types.string()).editor(data.editors.textarea({ rows: 4 })).summary(data.summary.elide())
  .validate()
)

const agendaEditor = data.collection("agenda", rows)
  .schema(agenda)
  .select(data.selection.urlParam("agenda", query.agenda))
  .edit(c => c
    .arrange(data.arrangements.masterDetail(md => md
      .summary(data.marks.table())
      .detail(data.marks.recordForm().submit(data.formPost("/settings/agenda-item")))
    ))
    .actions(a => a
      .create("New agenda item")
      .reorder(ui.action.server("admin-reorder-course-agenda"))
      .remove(ui.action.server("admin-delete-agenda-item").confirm(ctx => `Delete ${ctx.row.title}?`))
    )
  )
  .validate()
  .toIR()

ui.section("Agenda").anchor("agenda").child(agendaEditor).toIR()
```

The exact method names can change. The important design points are typed builders, lambda configurators, fragment composition, validation terminals, and a final IR terminal.

### Proposed addition: “Negative examples and failure modes”

The document should show the anti-patterns explicitly:

- A raw component page with `agenda.map(item => ui.panel(...))`.
- The current Phase 0 grammar that improves readability but still uses `Props` and maps.
- The target typed builder version.

This would make the playbook teach the progression rather than merely list patterns.

## Proposed target architecture for widgetdsl v2

### Layer 1: typed intent builders in Go

Introduce Go structs for the authoring concepts:

```go
type SchemaBuilder struct {
    name string
    fields []FieldSpec
    issues []ValidationIssue
}

type CollectionBuilder struct {
    name string
    rows []map[string]any
    schema *SchemaRef
    mode CollectionMode
    arrangement ArrangementRef
    actions []ActionBinding
    selection SelectionBinding
    issues []ValidationIssue
}

type Arrangement interface {
    Validate(schema SchemaSpec) []ValidationIssue
    Compile(ctx CompileContext) (WidgetNode, error)
}
```

Expose those through Goja wrappers rather than exported maps. Use either the goja-bleve hidden-ref substrate or a shared `fluent` package derived from it.

### Layer 2: composable fragments and lambda configurators

Borrow codesign's model:

- Top-level factory returns a builder.
- Sub-builders are configured by lambdas.
- `.use(fragment)` applies reusable fragments.
- Runtime validates callbacks with `goja.AssertFunction`.
- TypeScript declarations expose `FragmentFn<T>` with concrete builder interfaces.

This is the missing piece in widgetdsl. The current design calls itself composable, but composition is mainly nesting maps and passing option objects.

### Layer 3: typed schema, mark, and arrangement contracts

Make marks first-class:

```ts
interface Mark<TRecord> {
  kind: string
  requiredRoles?: FieldRole[]
  compile(recordOrCollection: TRecord | TRecord[], schema: Schema<TRecord>): WidgetNode
}

interface Arrangement<TRecord> {
  validate(schema: Schema<TRecord>): ValidationResult
  compile(collection: CollectionSpec<TRecord>): WidgetNode
}
```

Then domain modules can export real marks:

```js
cms.marks.assetTiles()
contextWindow.marks.treemap()
course.marks.lessonCard()
```

A recipe becomes a wrapper over grammar terms, not a separate compiler.

### Layer 4: precise TypeScript declarations and parity tests

The widget DSL should stop emitting only `Props`. Keep low-level `component(type, props)` as an escape hatch, but the grammar should have named interfaces and unions:

```ts
type FieldRole = "key" | "primary" | "short" | "prose" | "count" | "size" | "measure" | "date" | "status" | "tags" | "media" | "href"
type CollectionVerb = "show" | "edit" | "pick" | "manage"
type BuiltInArrangement = "table" | "master-detail" | "disclosure" | "tiles"
interface SchemaBuilder<T> { field<K extends keyof T>(name: K, ...): this; validate(): ValidationResult; build(): Schema<T> }
interface CollectionBuilder<T> { schema(schema: Schema<T>): this; edit(fn?: FragmentFn<CollectionEditBuilder<T>>): this; validate(): ValidationResult; toIR(): WidgetNode }
```

Add a generated-DTS parity test similar in spirit to geppetto's parity test as described in design-doc 01, and treat the declarations as part of the API surface.

## Decision records

### Decision: classify current widgetdsl as a prototype, not the canonical pattern

- **Context:** RAGEVAL-UI-GRAMMAR improved a real UI and introduced useful vocabulary, but the implementation remains map-based and weakly typed.
- **Options considered:** Promote current widgetdsl as the playbook model; reject it entirely; use it as vocabulary evidence while replacing the substrate.
- **Decision:** Use it as vocabulary evidence while replacing the substrate.
- **Rationale:** The page outcome validates `section`/`record`/`collection`; the implementation fails the runtime/compile-time typechecking bar.
- **Consequences:** Future documents should be careful not to say “copy widgetdsl”; they should say “copy the intent concepts, not the representation.”
- **Status:** proposed.

### Decision: use lambda-configurator fragments for composition

- **Context:** The current `arrange` string switch is not extensible enough for domain marks, multi-view context diagrams, or reusable page conventions.
- **Options considered:** Add more strings; add more recipes; adopt lambda configurators and `.use(fragment)`.
- **Decision:** Adopt lambda configurators and `.use(fragment)` for widgetdsl v2.
- **Rationale:** Codesign already demonstrates this model with typed builders, runtime callback validation, and precise declarations.
- **Consequences:** More initial implementation work, but a much better grammar for reusable domain-specific arrangements.
- **Status:** proposed.

### Decision: separate field facets instead of one overloaded `role`

- **Context:** The current role string controls identity, semantic kind, summary rendering, editor choice, and layout density.
- **Options considered:** Keep the single role; add more roles; split storage type, semantic role, editor, summary mark, validation, and layout hints.
- **Decision:** Split facets underneath; keep shortcut role helpers for convenience.
- **Rationale:** Shortcuts are good for authors, but the core representation must support fields that do not fit one role axis.
- **Consequences:** The TypeScript and Go model get richer; migration can map existing `f.*` helpers to facet presets.
- **Status:** proposed.

### Decision: add validation terminals before expanding arrangements

- **Context:** More arrangements/marks without validation will make failures harder to diagnose.
- **Options considered:** Implement tiles/disclosure/subpages first; implement typed validation first.
- **Decision:** Add validation and typed declarations before broadening the grammar.
- **Rationale:** The current API already has enough surface to expose the core weaknesses; adding features now compounds them.
- **Consequences:** Short-term visible UI progress slows, but the design decision becomes safer.
- **Status:** proposed.

## Concrete recommendations for the next revision

1. **Revise design-doc 01 to include a harsh widgetdsl self-assessment.** Make clear that Pattern C is a prototype pattern, not a target pattern.
2. **Add a target widgetdsl v2 API sketch based on typed builders and lambdas.** The current catalogue lists better DSLs; it should synthesize them into a widget-specific target.
3. **Add a migration plan from Phase 0 helpers to typed builders.** Keep `data.schema({...})` as a compatibility facade that internally creates a `SchemaBuilder` and calls `.build()`.
4. **Write one domain mark proof.** Re-express `cms.recipes.mediaLibrary` or `context_window.recipes.contextDiagram` as `data.collection + domain mark`. This is the real test of “domain modules shrink to schemas + marks.”
5. **Replace open-ended grammar declarations.** Low-level component helpers can stay loose; grammar helpers should get named types.
6. **Add validation and golden tests.** Test `validate()` issue paths, not only expansion shape. Keep expansion snapshots for IR output.
7. **Document action context contracts.** DataTable row actions, confirm interpolation, URL selection, and server-action payload merging need a typed contract.
8. **Define the role/facet model.** Decide which pieces are storage type, semantic role, editor, summary mark, layout, validation, transform, and scale.
9. **Add explicit non-goals.** For example: this grammar should not become a full React state system; it should compile to serializable IR; low-level `component()` remains an escape hatch.
10. **Add a design-risk section.** The largest risks are overfitting to admin CMS, creating an underpowered role taxonomy, and building an untyped macro system that looks like a DSL but cannot be safely extended.

## Review checklist for this DSL decision

Before adopting the current widgetdsl grammar as a foundation, require yes/no answers to these questions:

- Can a schema be distinguished at runtime from a random object without relying on a magic string?
- Can TypeScript catch an invalid arrangement or missing schema?
- Can validation report all schema/mark/action problems without throwing on the first one?
- Can a domain module export a mark that declares which roles it needs?
- Can `cms.recipes.mediaLibrary` be rewritten as a one-line wrapper over `data.collection`?
- Can reusable fragments be shared across pages (`adminDefaults`, `pagedTable`, `destructiveRowActions`)?
- Can action payloads be typed against row context?
- Can the generated declarations be tested against runtime exports?
- Can the low-level Widget IR remain an escape hatch without infecting the high-level grammar with `any`?

If the answer to most of these is “no,” the grammar is not yet good enough for the playbook.

## References

- `rag-evaluation-system/ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/01-goja-dsl-catalogue-and-base-research.md` — base catalogue and Pattern C classification.
- `rag-evaluation-system/ttmp/2026/07/04/RAGEVAL-UI-GRAMMAR--composable-ui-design-system-grammar-for-the-widget-dsls-cms-admin-page-readability-overhaul-and-cross-dsl-section-list-form-primitives/design-doc/01-composable-ui-grammar-for-the-widget-dsls-analysis-of-the-cms-admin-page-and-brainstormed-alternatives.md` — original audit and brainstorm.
- `rag-evaluation-system/ttmp/2026/07/04/RAGEVAL-UI-GRAMMAR--composable-ui-design-system-grammar-for-the-widget-dsls-cms-admin-page-readability-overhaul-and-cross-dsl-section-list-form-primitives/design-doc/02-dsl-api-sketch-grammar-verbs-in-data-dsl-and-ui-dsl-module-reorganization-no-new-module.md` — API sketch that led to implementation.
- `/home/manuel/code/wesen/go-go-golems/go-go-parc/Projects/2026/07/05/ARTICLE - Widget DSL Grammar - Designing an Intent-Level UI Authoring Layer for a Widget IR System.md` — narrative article documenting the work.
- `rag-evaluation-system/pkg/widgetdsl/grammar.go` — current Phase 0 grammar implementation.
- `rag-evaluation-system/pkg/widgetdsl/typescript.go` — weak generated declarations.
- `rag-evaluation-system/pkg/widgetdsl/module.go` — module maps, recipes, helper promotion, and existing direct recipe compilers.
- `rag-evaluation-system/packages/rag-evaluation-site/src/widgets/ir.ts` — typed React-side Widget IR and DataTable/FormPanel prop contracts.
- `rag-evaluation-system/packages/rag-evaluation-site/src/widgets/actions.ts` — string-template action interpolation and confirm handling.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js` — first consumer of `ui.section` and `data.collection`.
- `/home/manuel/code/wesen/go-go-golems/goja-bleve/pkg/api_types.go` and `api_mapping.go` — typed-ref builder substrate comparator.
- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/builders.go` and `typescript.go` — lambda-configurator + fragment-composition comparator.
