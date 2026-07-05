---
Title: Goja DSL catalogue and base research
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
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/discord-bot/pkg/doc/topics/discord-js-bot-api-reference.md
      Note: defineBot + Proxy-trap ui typed builders (Pattern G)
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/api_engine_builder.go
      Note: Clone-on-each-step immutable fluent builder (Pattern A prime)
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/dts_parity_test.go
      Note: Canonical DTS parity test enforcing generated .d.ts against runtime exports
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/goja-bleve/pkg/api_mapping.go
      Note: fieldBuilder fluent chain + .build() terminal — the reference fluent builder
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/goja-bleve/pkg/api_types.go
      Note: Typed-ref substrate (refBase
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/goja-dbus/pkg/dbusgoja/builders.go
      Note: Cleanest composable grammar (bus/destination/object/interface/method/out/call)
    - Path: ../../../../../../../../../2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/builders.go
      Note: Strongest model - runSpec fluent builder + use fragment composition + lambda configurators
    - Path: ../../../../../../../../../2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/typescript.go
      Note: Precise RunSpecBuilder/TopologyBuilder/MetricsBuilder TS interfaces via TypeScriptDeclarer
    - Path: ../../../../../../../../../2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/researchctl/module.go
      Note: Lambda-configurator project builder (Pattern F)
    - Path: ../../../../../../../go-go-goja/pkg/tsgen/spec/types.go
      Note: TypeScript declaration spec model — compile-time type substrate
    - Path: ../../../../../../../go-minitrace/pkg/minitracejs/builders.go
      Note: Pattern B reference (SourceSetBuilder
    - Path: pkg/widgetdsl/grammar.go
      Note: The not-great data.dsl grammar verbs (f/schema/record/collection
    - Path: pkg/widgetdsl/typescript.go
      Note: Concrete evidence of weak compile-time types
ExternalSources: []
Summary: Base research catalogue of every Goja-based DSL across go-go-goja, go-minitrace, rag-evaluation-system, and the goja-* family, with their APIs, documentation locations, examples, and the implementation patterns they use. Intended as the evidence base a senior researcher turns into a fluent-builder DSL playbook.
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: Catalogue every Goja DSL we have built, map the implementation patterns, and identify which ones are worth promoting into the opinionated fluent-builder playbook.
WhenToUse: Read this before writing or reviewing the playbook; it is the single inventory of what exists today.
---



# Goja DSL catalogue and base research

## 0. How to read this document

This is a **base research** document, not the final playbook. Its job is to inventory every Goja-based DSL across the relevant repositories, capture each DSL's public API, where it is documented, a runnable example, and — critically — the **implementation pattern** behind it. A separate senior-researcher pass will use this catalogue to reflect on the patterns, assess them, and write the opinionated fluent-builder playbook.

The companion documents in this ticket are:

- `reference/01-research-logbook-resource-assessment.md` — per-resource assessment (useful / out of date / needs updating) for every file, README, and external page read during this research.
- `reference/02-investigation-diary.md` — chronological investigation diary.

Repositories surveyed (all on the local filesystem):

| Repository | Path | Role |
| --- | --- | --- |
| go-go-goja | `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-goja` | goja runtime + native-module host + `engine` builder + `uidsl` + `express` route DSL |
| go-minitrace | `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-minitrace` | transcript analysis toolkit with the `minitracejs` builder DSL |
| rag-evaluation-system | `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system` | the `widgetdsl` family (`ui.dsl`, `data.dsl`, `context_window.dsl`, `course.dsl`, `cms.dsl`) — the "latest" DSLs, where `data.dsl` is explicitly flagged as "not really all that great" |
| goja-bleve | `~/code/wesen/go-go-golems/goja-bleve` | the most sophisticated fluent-builder DSL (typed Go refs) |
| goja-text | `~/code/wesen/go-go-golems/goja-text` | markdown / sanitize / extract / template modules |
| goja-git | `~/code/wesen/go-go-golems/goja-git` | imperative Git DSL |
| goja-github-actions | `~/code/wesen/go-go-golems/goja-github-actions` | `@actions/*` polyfills + `@goja-gha/ui` report DSL |
| goja-dbus | `~/code/wesen/go-go-golems/goja-dbus` | typed value helpers + fluent D-Bus bus/method/signal builders |
| goja-treesitter | `~/code/wesen/go-go-golems/goja-treesitter` | skeleton/stub (no module yet) |
| geppetto | `~/code/wesen/go-go-golems/geppetto` | LLM runtime core with a large `require("geppetto")` DSL (engine/agent/tool/schema builders + DTS parity test) |
| discord-bot | `~/code/wesen/go-go-golems/discord-bot` | `require("discord")` defineBot framework + `require("ui")` Go-side typed-builder DSL via Proxy traps |
| researchctl | `~/workspaces/2026-06-30/benchmark-cpu-inference/researchctl` | `require("researchctl")` claim/evidence/decision graph builder — **lambda-configurator** pattern |
| codesign (in researchctl) | `…/researchctl/pkg/gojamodules/codesign` | `require("codesign")` CPU/GPU run builder — **lambda-configurator + `.use()` fragment composition** |
| glazed | `…/glazed` | Go-native CLI command builder (the canonical functional-options pattern that inspired the goja DSLs) |
| go-emrichen | `~/code/wesen/go-go-golems/go-emrichen` | YAML tag-operator DSL (composable operators, strong typing) |
| go-emrichen | `~/code/wesen/go-go-golems/go-emrichen` | YAML tag-operator DSL (composable operators, strong typing) |

## 1. Executive summary

Across the surveyed repositories there are **at least 16 distinct Goja DSLs** plus three pieces of cross-cutting infrastructure (`engine`, `xgoja`/`tsgen`, `glazed`). They fall into **six implementation patterns**. The ticket's goal — strict runtime typechecking, validation, composable grammar, compile-time types, and lambdas as extensibility — is realised most completely by **researchctl** and **codesign** (both added in a follow-up pass after the user flagged `geppetto` and `discord-bot`, which in turn led to `researchctl`/`codesign`). `goja-bleve` and `goja-dbus` remain the strongest *typecheck* models; `researchctl`/`codesign` are the strongest *composable-grammar + lambda* models.

The "latest" DSLs in `rag-evaluation-system/pkg/widgetdsl` (`data.dsl` in particular) use a **map-IR + loose helper function** pattern that is easy to author but gives up type safety, returns untyped `map[string]any`, validates by `panic(vm.NewGoError(...))`, and emits open-ended TypeScript (`Props = Record<string, any>`). This matches the user's assessment that `data.*` is "not really all that great."

The strongest existing models for the playbook, by concern:

- **Runtime typecheck substrate** — **goja-bleve**: typed `refBase` + `refKind` enum (`pkg/api_types.go:18-130`), generic `getTypedRef[T]` extractor (`pkg/api_types.go:131`), `newWrapper(ref, kind)` attaching non-enumerable `__bleve_ref` (`pkg/api_types.go:143`, `pkg/module.go:16`), same-object chaining + `.build()` terminal returning a new wrapper kind (`pkg/api_mapping.go:137-200`), terminal validation returning `(value, error)`. **geppetto** reuses this exact substrate (`hiddenRefKey = "__geppetto_ref"`) but with **clone-on-each-step** immutable builders — a notable variant.
- **Composable grammar + lambdas** — **researchctl** and **codesign**: fluent `project(name).goal(title, g => g.id(...).status(...))` and `runSpec(name).topology(fn).workload(fn).metrics(fn)` where each `fn` is a **lambda configurator** applied to a sub-builder. codesign adds `.use(fragment)` with `FragmentFn<T>` reusable builder lambdas — the cleanest realisation of "composable grammar of operators extended with lambdas." Both emit **precise TypeScript `interface`s** (not `any`) and implement `modules.TypeScriptDeclarer`. geppetto goes further with a **DTS parity test** (`dts_parity_test.go`) that asserts the generated `.d.ts` matches the runtime export surface.
- **Typed value helpers + fluent grammar** — **goja-dbus**: `dbus.u32()`/`variant()`/`dict()`/`struct()` + `bus.destination().object().interface().method().out().call()` resolving to Promises, with strict policy enforcement and all callbacks on the runtime owner.
- **Validation discipline** — **go-minitrace**: `Validate()` → `ValidationResult{Valid, Errors}` + `Build()` → `(value, error)`, with accumulated errors. researchctl/codesign adopt the same `.validate()`/`.toSpec()` terminals.
- **Proxy-trap typed builders** — **discord-bot `require("ui")`**: Go-side fluent builders surfaced through Goja Proxy traps; returns **typed builders**, not plain JS objects; wrong-parent calls fail loudly; raw JS objects rejected where builders are expected.
- **Conceptual operator grammar** — **go-emrichen**: tag operators (`!Defaults`/`!If`/`!Loop`) as composable, type-aware operators.

The playbook should treat: goja-bleve's typed-ref machinery as the runtime-typecheck substrate; researchctl/codesign's lambda-configurator + `.use()` fragment composition as the composable-grammar model; geppetto's DTS parity test as the compile-time-type enforcement discipline; go-minitrace's `Validate()`/`Build()` as the validation rule; and discord-bot's `ui` Proxy-trap builders as an alternative typed-builder mechanism worth comparing against the hidden-key approach.

## 2. Problem statement and scope

**Problem.** We have built many Goja DSLs over time with inconsistent patterns. Some are fluent and type-safe (goja-bleve, goja-dbus), some are pragmatic builder structs (go-minitrace), and some are loose map-IR helpers with no type safety (widgetdsl `data.dsl`). There is no playbook that says: *when you build a new Goja DSL, here is the opinionated way to do it — fluent builders, Go-side implementation, strict runtime typechecking, validation, compile-time types from generated declarations, and a composable grammar extensible with lambdas.*

**Scope of this document.** Catalogue and evidence base only. It does **not** prescribe the final playbook API. It inventories what exists, classifies the patterns, and flags what is useful, out of date, or missing — so a senior researcher can do the reflection and assessment pass.

**Out of scope.** Implementing the playbook, refactoring existing DSLs, or choosing a winner. Those are follow-ups.

## 3. Pattern taxonomy

Five distinct implementation patterns appear across the surveyed DSLs. Every DSL maps to exactly one primary pattern (a few borrow a second).

### Pattern A — Typed Go reference + fluent terminal builder (goja-bleve, goja-dbus)

- Every JS-facing handle wraps a typed Go struct via a non-enumerable hidden key (`__bleve_ref`).
- A `refKind` enum tags each handle so the Go boundary can reject wrong-type handles.
- A generic `getTypedRef[T]` does runtime typecheck + extraction.
- Fluent methods mutate the wrapped Go struct and return the **same** JS object (chainable).
- A `.build()` terminal returns a **new** wrapper of a different kind (builder → built artifact).
- Validation happens in terminals, returning `(value, error)`.
- Strongest type safety; closest to the ticket's goal.

### Pattern B — Plain Go builder struct + re-wrap (go-minitrace `minitracejs`)

- A plain Go struct (`DBBuilder`, `SourceSetBuilder`, `QueryRecipeBuilder`) accumulates state and an `errors []string` slice.
- Each chain method mutates the struct and returns a **fresh** `*goja.Object` that re-wraps the same struct pointer.
- `.Validate()` returns a `ValidationResult{Valid, Errors}`; `.Build()` returns `(value, error)`.
- Errors are accumulated and reported at terminals rather than thrown immediately.
- Pragmatic, easy to test in pure Go, but loses the "same object" identity that Pattern A uses for typed-ref lookup.

### Pattern C — Map IR + loose helper functions (rag-evaluation-system `widgetdsl`)

- Helpers return plain `map[string]any` (JSON-serializable IR).
- No typed builder objects; composition is by nesting maps.
- Validation by `panic(r.vm.NewGoError(...))`.
- Option merging via `mergeOptions`/`exportOptions`.
- TypeScript generation emits open-ended `Props = Record<string, any>` — **weak compile-time types**.
- This is the pattern the user wants to move away from.

### Pattern D — Hyperscript / element-function DSL (go-go-goja `uidsl`)

- `ui.div(attrs, ...children)` returns a Go `Element{Tag, Attrs, Children}` struct.
- `class` accepts arrays with falsy filtering; `style` accepts objects.
- `render()` walks the tree to HTML.
- Good for static structure, not designed for fluent builder composition or typechecked terminals.

### Pattern E — Imperative object API (goja-git, goja-github-actions polyfills)

- `git.open({Dir})` returns a repo object; `repo.add({Paths})`, `repo.commit({Message, Author})`.
- Options objects, no chaining, no builders.
- Mimics popular JS libraries (`simple-git`, `@actions/core`).
- Not a builder pattern; included for completeness and contrast.

### Pattern A′ — Typed-ref + clone-on-each-step immutable builder (geppetto)

- Same hidden-key typed-ref substrate as Pattern A (`__geppetto_ref`, `attachRef`, `mustSet`, `requireXxxRef` type-asserting extractors).
- But each fluent method returns a **new** wrapper built from a **cloned** ref (`cloneFor(m)`), not the same mutated object. Builders are immutable per step.
- `.build()` terminal materialises the Go artifact from the accumulated ref.
- **Geppetto also enforces a DTS parity test** (`TestGeneratedDTSMatchesRuntimeExportSurface`) that asserts the generated `geppetto.d.ts` matches the runtime export surface — the strongest compile-time-type discipline found.
- Validation is mixed: some `panic(m.vm.NewTypeError(...))`, some `(value, error)`.

### Pattern F — Fluent builder + lambda configurators + fragment composition (researchctl, codesign)

- A top-level factory (`project(name)`, `runSpec(name)`) returns a builder.
- Each collection method takes `(title, build?)` where `build` is a **JS lambda configurator** applied to a fresh sub-builder: `project("X").goal("G", g => g.id("GOAL-001").status("active").priority("P0"))`.
- Sub-builders expose the same fluent same-object chaining (`g.id(...).status(...).priority(...)`).
- **codesign adds `.use(fragment)` composition**: a `FragmentFn<T>` is a reusable `(b: T) => void` lambda applied to a builder, so reusable topology/workload/metric fragments can be shared across runs.
- Lambdas also cross the Go/JS boundary as **runtime callbacks**: `jsDevice(id, callback, config?)`, `policyCallback(id, fn)`, `callback(id, fn)` metrics — validated with `goja.AssertFunction`.
- Terminals: `.validate()` → `ValidationResult`, `.toSpec()` → typed spec, `.run(options?)` (codesign).
- Precise TypeScript `interface`s emitted via `modules.TypeScriptDeclarer` (e.g. `RunSpecBuilder`, `TopologyBuilder`, `MetricsBuilder`, `RunSpecLike = RunSpec | RunSpecBuilder | { toSpec(): RunSpec }`).
- **This pattern is the closest existing realisation of the ticket's full goal** (fluent + Go + typecheck + validation + compile-time types + composable grammar + lambdas).

### Pattern G — Go-side typed builders via Proxy traps (discord-bot `require("ui")`)

- Go-side fluent builders surfaced to JS through **Goja Proxy traps** (not hidden keys).
- Returns **typed builders**, not plain JS objects; `.build()` terminal returns a typed Discord payload or the host's `normalizedResponse` fast path.
- Wrong-parent calls fail loudly (`ui.message().field(...)` is an error telling you to use `ui.embed()`); raw JS objects rejected where builders are expected.
- Coexists with a registration-style bot DSL (`defineBot(({command, event, component, modal, autocomplete, configure}) => {...})`).
- An alternative typed-builder mechanism worth comparing against the hidden-key approach (Pattern A/A′).

### Cross-cutting — Go-native functional-options builder (glazed, go-emrichen)

These are **not** Goja DSLs, but they are the conceptual ancestors of the fluent-builder pattern and the canonical Go-side grammar.

- glazed: `schema.NewSchema(WithSections(...))`, `fields.New(name, type, WithHelp(...), WithRequired(true))` — functional options over `Definition` structs, `Section` interface, `Definitions` collection.
- go-emrichen: YAML tag operators (`!Defaults`, `!If`, `!Var`, `!Format`, `!Loop`) — each tag is a composable operator with type-aware validation. Closest existing thing to "a core composable grammar of operators extended with lambdas."

## 4. The DSL inventory

Each entry below gives: module name, repository/path, public API, documentation location, a runnable example, the implementation pattern, and evidence anchors.

---

### 4.1 goja-bleve — fluent mapping/query/search builders (Pattern A, gold standard)

**Module name.** `require("bleve")`

**Repository.** `~/code/wesen/go-go-golems/goja-bleve`

**Public API (factories).**
```js
const bleve = require("bleve")
bleve.mapping()           // index mapping builder
bleve.docMapping()        // document mapping builder
bleve.field()             // field mapping builder (terminal: .build())
bleve.memory()            // in-memory index builder
bleve.indexMapping()
bleve.matchAll() / bleve.match(q) / bleve.matchNone()
bleve.search()            // search request builder
```

**Fluent example (mapping + vector + hybrid search).**
```js
const embedding = bleve.field()
  .vector(4)
  .similarity("cosine")
  .optimizedFor("recall")
  .build()

const request = bleve.search()
  .query(bleve.match("privacy").field("text"))
  .knn("embedding", [1,0,0,0], 10, 1.0)
  .score("rrf")
  .scoreRankConstant(60)
  .scoreWindowSize(50)
  .build()
```

**Field builder fluent chain.**
```js
bleve.field().text().store(true).includeTermVectors(true).build()
bleve.field().keyword().store(true).build()
bleve.field().number().store(true).build()
bleve.field().datetime().store(true).build()
bleve.field().boolean()
bleve.field().geoPoint()
bleve.field().ip()
bleve.field().disabled()
bleve.field().vector(4)            // returns (obj, err), dims validated
bleve.field().vectorBase64(4)
```

**Batch lifecycle (single-use after execute).**
```js
const batch = index.newBatch()
batch.index("id", doc)
batch.delete("id")
batch.execute()   // after this, mutation/reset throws "batch has already been executed"
```

**Implementation pattern.** Pattern A. The typed-ref machinery is the key contribution:
- `refBase{api, kind, closed}` + typed ref structs embedding it (`pkg/api_types.go:18-130`).
- `getTypedRef[T]` generic extractor enforces types at the Go boundary (`pkg/api_types.go:131-141`).
- `newWrapper(ref, kind)` attaches non-enumerable `__bleve_ref` (`pkg/api_types.go:143-148`, `pkg/module.go:16`).
- `mustSet(obj, key, value)` helper for method binding (`pkg/module.go:187`).
- Fluent methods return the same `*goja.Object`; `.build()` returns a new wrapper of a different kind (`pkg/api_mapping.go:137-200`).
- Build tag `vectors` gates vector/KNN support with clear non-vector errors.

**Documentation.**
- `README.md` (root) — Phase 7 surface, batch lifecycle, vector/KNN, hybrid scoring, provider integration.
- `docs/quickstart.md`, `docs/faiss-xgoja-playbook.md` (FAISS build/link/runbook), `docs/README.md`.
- Glazed help pages bundled in the xgoja binary (provider-shipped).

**Evidence.** `pkg/api_types.go:18-148`, `pkg/api_mapping.go:137-200`, `pkg/module.go:16,162,187`, `pkg/mapping_test.go:17-19,79-94`.

---

### 4.2 goja-dbus — typed values + fluent bus/method/signal builders (Pattern A)

**Module name.** `require("dbus")`

**Repository.** `~/code/wesen/go-go-golems/goja-dbus` (JS layer in `pkg/dbusgoja/`, Go core in `pkg/dbuscore/`).

**Typed value helpers.**
```js
dbus.u32(42)
dbus.i32(-7)
dbus.path("/com/example/App1")
dbus.signature("s")
dbus.variant("s", "hello")
dbus.array("as", ["default", "Open"])
dbus.dict("a{sv}", { urgency: dbus.variant("u", dbus.u32(1)) })
dbus.struct("(su)", ["count", dbus.u32(7)])
```

**Fluent bus + method-call builder (resolves to Promise).**
```js
const bus = await dbus.session().timeout(2000).connect()
try {
  const id = await bus
    .destination("org.freedesktop.DBus")
    .object("/org/freedesktop/DBus")
    .interface("org.freedesktop.DBus")
    .method("GetId")
    .out("s")
    .call()
  console.log("bus id:", id)
} finally {
  await bus.close()
}
```

**Fluent signal subscription builder.**
```js
const sub = await bus
  .signals()
  .interface("org.freedesktop.DBus.Properties")
  .member("PropertiesChanged")
  .listen(emitter)   // EventEmitter-based delivery
await sub.close()
```

**Implementation pattern.** Pattern A. Each builder method returns `goja.Value` and re-wraps the same builder. Strict policy enforcement (default-denied system bus). All callbacks/Promise settlement happen on the go-go-goja runtime owner. Typed values carry `{signature, value}` so the Go side can marshal D-Bus signatures correctly.

**Documentation.**
- `README.md` (root) — implemented vs deferred, three runnable examples.
- docmgr ticket at `ttmp/2026/06/15/GOJA-DBUS-DESIGN--goja-d-bus-module-intern-design-guide/` — intern-facing design and implementation guide.
- xgoja provider with `getting-started`, `user-guide`, `api-reference` Glazed help pages.

**Evidence.** `pkg/dbusgoja/builders.go:28-148` (timeout/policy/connect/close/destination/object/interface/method/in/out/call), `pkg/dbusgoja/signals.go:28-89` (sender/path/interface/member/listen/close), `pkg/dbusgoja/typed_values.go:16-17`, `pkg/dbusgoja/errors.go:7-8`.

---

### 4.3 go-minitrace `minitracejs` — pragmatic builder structs (Pattern B)

**Module name.** `require("minitrace")` (the `minitracejs` package).

**Repository.** `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-minitrace`

**Public API (factory entry points).**
```js
const m = require("minitrace")
m.importer()        // import pipeline builder
m.db()              // DuckDB builder (sources, cache, query options)
m.sources()         // source set builder
m.importPolicy()    // conversion policy builder
m.cache()           // cache policy builder
m.limits()          // query limits builder
m.query()           // query recipe builder
m.view()            // view plan builder
m.session()         // session builder
m.runtime          // runtime settings object
m.sql               // SQL helpers: .string(), .stringIn(), .like()
```

**Fluent example (sources + db).**
```js
const sources = m.sources()
  .File("./a.json")
  .Archive("./b.minitrace.json")
  .Dir("./sessions/")
  .Glob("./output/active/*/*.minitrace.json")
  .RuntimeArchives()
  .Name("my-set")

const validation = sources.Validate()   // { valid, errors }
const set = sources.Build()              // (SourceSet, error)
```

**Query recipe builder (grammar of recipe kinds).**
```js
m.query()
  .SessionSummary()
  .SessionID("abc")
  .ByTurn()
  .Build()
// kinds: SessionSummary, TurnRows, ToolRows, EventRows, TurnBlockRows,
//        TokenUsageRows, TranscriptRows, TimelineRows
// grouping: BySession / ByTurn / ByRole / ByTool
```

**Import builder (multi-terminal: Detect/Convert/Preview/Diagnostics/Save).**
```js
m.importer()
  .File("./claude.json")
  .Into("./out")
  .SessionID("abc")
  .Strict(true)
  .Detect()        // (map, error)
  .Convert()       // (obj, error)
  .Preview()       // (map, error)
  .Diagnostics()   // []map
  .Save()          // (map, error)
```

**Implementation pattern.** Pattern B. Plain Go builder structs (`DBBuilder`, `SourceSetBuilder`, `QueryRecipeBuilder`, `ImportBuilder`, `ViewPlanBuilder`, `SessionBuilder`, `ImportPolicyBuilder`, `CachePolicyBuilder`, `QueryLimitsBuilder`) each accumulate state + `errors []string`. Chain methods re-wrap the same struct pointer into a fresh `*goja.Object`. `Validate()` returns `ValidationResult{Valid, Errors}`; `Build()`/terminals return `(value, error)`. SQL helpers (`SQLString`, `SQLStringIn`, `SQLLike`) provide safe escaping.

**Documentation.**
- `README.md` (root) — pipeline framing, install, quick start.
- Glazed structured commands under `pkg/minitracecmd/` (preset/recipe system).
- TypeScript descriptor at `pkg/minitracejs/typescript.go` (minimal).

**Evidence.** `pkg/minitracejs/module.go:29-69` (factory exports), `pkg/minitracejs/builders.go:1-90` (`SourceSetBuilder`, `sourcesBuilderObject`), `pkg/minitracejs/db_builder.go:1-60` (`DBBuilder`, `dbSource`, `ValidationResult`, `DBHandle`), `pkg/minitracejs/import_builder.go:150-218` (Content/File/Name/Format/Strict/Into/Detect/Convert/Preview/Diagnostics/Save), `pkg/minitracejs/query_view_session.go:1-80` (`QueryRecipeBuilder`, recipe kinds + grouping).

---

### 4.4 rag-evaluation-system `widgetdsl` — the map-IR family (Pattern C, the "not great" one)

**Module names.** Five split modules, registered in `pkg/widgetdsl/module.go:139-178`:
- `ui.dsl` (`UIModuleName`)
- `data.dsl` (`DataModuleName`) — the one explicitly flagged as not great
- `context_window.dsl` (`ContextWindowModuleName`)
- `course.dsl` (`CourseModuleName`)
- `cms.dsl` (`CmsModuleName`)

**Repository.** `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system`

**`data.dsl` grammar (the grammar verbs added in RAGEVAL-UI-GRAMMAR).**
```js
const data = require("data.dsl")

// field roles (not types): key, primary, short, prose, count, size, measure,
// date, status, tags, media, href
const schema = data.schema({
  id:    data.f.key({ label: "ID" }),
  title: data.f.primary({ required: true, maxLength: 160 }),
  body:  data.f.prose(),
  count: data.f.count(),
  status: data.f.status(),
})

data.record(values, { schema })
data.collection(rows, { schema, arrangement })
data.urlParam("id", value)
data.formPost("/api/save")
```

**`data.dsl` helpers (table cells, recipes).**
```js
data.dataTable({ columns: [
  { id: "title",  header: "Title",  cell: data.cell.field("title") },
  { id: "status", header: "Status", cell: data.cell.status("status", { icon: true }) },
]})
data.cell.number(field) / .caption(field) / .template(tpl) / .link(...) / .linkButton(...)
data.cell.actionButton(label, action) / .constant(value)
data.recipes.masterDetailTable(...)
data.page({ sections: [...] })
```

**`ui.dsl` helpers (component factory + structure grammar).**
```js
const ui = require("ui.dsl")
ui.panel({ title: "x" }, ui.text("hello"))
ui.section("Title", { collapsible: true }, ui.stack(...))
ui.appShell / ui.appNav / ui.button / ui.breadcrumbs / ui.emptyState
ui.fieldGrid / ui.markdownArticle / ui.meterBar / ui.pagination / ui.tag
// ... ~35 promoted generic primitives
```

**`action` sub-grammar.**
```js
data.action.server("saveRecord") / .navigate(to) / .download(to)
data.action.event("copy") / .copy(value)
```

**Implementation pattern.** Pattern C. Helpers return plain `map[string]any`. `schemaCtor` preserves field order via goja's insertion-ordered object keys and tags the result with `__ragSchema: true` (`pkg/widgetdsl/grammar.go:81-103`). `record`/`collection` validate by reading that tag and `panic(r.vm.NewGoError(...))` on misuse (`grammar.go:161, 276`). Field roles are an enumerated string list (`fieldRoles`, `grammar.go:18-33`); `gridableRoles` is a separate allow-set. `mergeOptions`/`exportOptions` flatten variadic option objects. TypeScript generation emits open-ended `Props = Record<string, any>` and `WidgetNode { kind: string; [key: string]: any }` — **no per-field type narrowing** (`pkg/widgetdsl/typescript.go:1-80`).

**Why it is "not great" (evidence-backed).**
1. **No typed handles.** Everything is `map[string]any`; the Go boundary cannot reject a wrong-shape cell spec until render time.
2. **Validation by panic.** `panic(r.vm.NewGoError(...))` is the only error path (`grammar.go:81, 161, 276`). No accumulated errors, no `Validate()` step.
3. **Weak compile-time types.** The `.d.ts` uses `Props = Record<string, any>` and `[key: string]: any` everywhere, so TypeScript cannot catch wrong cell/field options.
4. **Schema tag is a magic string** (`__ragSchema`) checked by type-assertion, not a typed wrapper.
5. **Two parallel authoring layers** (raw component helpers + grammar verbs) with implicit compilation rules — easy to get wrong.

**Documentation.**
- `pkg/widgetdsl/module.go` (1215 lines) — module specs, helpers, page/cell/action wiring.
- `pkg/widgetdsl/grammar.go` (537 lines) — the data grammar verbs and field roles.
- `pkg/widgetdsl/typescript.go` — declaration generation.
- docmgr tickets:
  - `RAGEVAL-UI-DSL` (`ttmp/2026/06/02/RAGEVAL-UI-DSL--…`) — original widget DSL design (3 design docs).
  - `RAGEVAL-UI-GRAMMAR` (`ttmp/2026/07/04/RAGEVAL-UI-GRAMMAR--…`) — the grammar-verbs addition that introduced `f.*`, `schema`, `record`, `collection`. Design-doc 02 is the API sketch; design-doc 01 is the analysis of alternatives.
  - `RAGEVAL-CMS-WIDGETS`, `RAGEVAL-CONTEXT-WINDOWS-DESIGN`, `RAGEVAL-WIDGET-IR-SEMANTIC-COMPONENTS`, `DESIGN-REF-001`, `CTX-WINDOW-BLOCK-VIZ`, `CTX-COLOR-PALETTE` — related widget-system tickets.
- `ttmp/2026/06/03/rag-eval-scripting-expansion--…` — xgoja + goja-text jsverbs expansion (intern guide + design).

**Evidence.** `pkg/widgetdsl/module.go:15-16,139-178,232-260`, `pkg/widgetdsl/grammar.go:1-120,161,276`, `pkg/widgetdsl/typescript.go:1-80`, `pkg/widgetdsl/module_test.go:22-164`, `pkg/widgetdsl/grammar_test.go:39-44`.

---

### 4.5 go-go-goja `uidsl` — hyperscript element DSL (Pattern D)

**Module name.** `require("ui")` (package `go-go-goja/modules/uidsl`).

**Repository.** `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-goja`

**Public API.**
```js
const ui = require("ui")
ui.div({ class: ["base", false, "extra"], style: { color: "red" }, id: "n1" }, "child")
ui.span / ui.a / ui.p / ui.ul / ui.li / ui.table / ui.tr / ui.td / ...
ui.text("hello")
ui.fragment(child1, child2)
ui.render(node)   // → HTML string
```

**Implementation pattern.** Pattern D. `ui.div(attrs, ...children)` returns a Go `Element{Tag, Attrs, Children}` (`modules/uidsl/node.go`, `components.go`). `class` accepts arrays with falsy filtering; `style` accepts objects rendered to CSS (`render_attrs_test.go`, `attrs_compat_test.go`). `render.go` walks the tree to HTML. Benchmarks exist for attribute-heavy nodes (`attrs_bench_test.go`).

**Documentation.** Source-internal only (tests + benchmarks). No README section; documented through test cases.

**Evidence.** `modules/uidsl/module.go`, `modules/uidsl/components.go`, `modules/uidsl/node.go`, `modules/uidsl/render.go`, `modules/uidsl/attrs_compat_test.go:1-40`.

---

### 4.6 go-go-goja `express` — fluent HTTP route + auth builder (Pattern A/B hybrid)

**Module name.** `require("express")` (package `go-go-goja/modules/express`).

**Repository.** `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-goja`

**Public API (planned-auth fluent builders).**
```js
const express = require("express")
express.user().required().mfaFresh("5m")        // → auth spec
express.resource("file")                        // → resource spec
// routes carry .auth(spec) which validates the spec came from express.user()
```

**Implementation pattern.** Hybrid. A `builderStore` holds `sync.Map[*goja.Object]*gojahttp.SecuritySpec` and `*gojahttp.ResourceSpec` (`modules/express/auth_builders.go:8-12`). Builders return the same JS object (Pattern A flavour) but typed lookup is via the side-channel map rather than a hidden key. `.auth(value)` validates the value was produced by `express.user()` (`auth_builders.go:43-58`). Duration parsing validates and returns `(goja.Value, error)`.

**Documentation.** Glazed help tree in `go-go-goja`:
- `goja-repl help express_auth_user-guide` / `express_auth-examples`
- `xgoja help go-planned-auth-api`, `express-auth-host-integration-guide`, `hostauth-config-reference`, `auth-stores-reference`
- Deploy runbooks: `goja-repl help deploying-an-express-auth-host`, `xgoja help auth-host-production-runbook`

**Evidence.** `modules/express/express.go:1-70` (registrar), `modules/express/auth_builders.go:8-60` (user/resource builders + spec validation), `modules/express/typescript.go`.

---

### 4.7 goja-text — markdown / sanitize / extract / template (Pattern B-ish)

**Module name.** `require("markdown")`, `require("sanitize")`, `require("extract")`, `require("template")`.

**Repository.** `~/code/wesen/go-go-golems/goja-text`

**Public API.**
```js
const markdown = require("markdown")
markdown.parse("# Hello")            // → Go-backed AST (PascalCase fields)
markdown.walk(ast, (node, ctx) => {}) // document-specific query
markdown.textContent(ast)
markdown.render(ast)

const sanitize = require("sanitize")
sanitize.json.sanitize(input)        // repair + fix metadata
sanitize.yaml.sanitize(input)

const extract = require("extract")
extract.all(input) / .markdown(input) / .frontmatter(input) / .tags(input)

const template = require("template")
template.text().Parse(input).Execute(data)   // fluent on template object
template.html().Parse(input).Execute(data)   // Glazed/Sprig helpers
```

**Implementation pattern.** Functional with Go-backed option builders. The markdown AST is exposed as Go structs read via PascalCase field names. `walk()` is the central design choice: small Go API, JS owns the query. `template` uses a fluent `.Parse().Execute()` chain.

**Documentation.** Strong — paired Glazed help entries per module (API reference + user guide): `goja-text-markdown-user-guide`/`-api-reference`, `sanitize-…`, `extract-…`, `template-…` plus `goja-text-template-writing-documentation`. `README.md` is detailed. `xgoja.yaml` at repo root wires the provider + modules + jsverbs + help.

**Evidence.** `README.md` (root), `xgoja.yaml`, `pkg/xgoja/providers/text/`.

---

### 4.8 goja-git — imperative Git DSL (Pattern E)

**Module name.** `require("git")` (via `git` global or require).

**Repository.** `~/code/wesen/go-go-golems/goja-git`

**Public API.**
```js
const git = require("git")
const repo = git.init({ Dir, DefaultBranch: "main", Bare: false })
const repo = git.open({ Dir })
repo.status()                       // → [{ path, staging, worktree }]
repo.add({ Paths: [...] } | { All: true })
const hash = repo.commit({ Message, Author: { Name, Email }, Amend: false })
repo.log({ Ref: "HEAD", Depth: 10 })
repo.branch / repo.checkout / repo.tags / repo.diff
git.filterRepo({ ... })             // history rewrite / subdir extraction
```

**Implementation pattern.** Pattern E. Imperative options objects, no chaining, returns plain JS objects/arrays. Mirrors `simple-git` / `isomorphic-git`.

**Documentation.** `README.md` (root, detailed), `TEST-RESULTS.md`, `examples/`.

**Evidence.** `README.md`, `pkg/`, `examples/`, `filterrepo/`.

---

### 4.9 goja-github-actions — `@actions/*` polyfills + `@goja-gha/ui` report DSL (Pattern E + DSL)

**Module names.** `@actions/core`, `@actions/github`, `@actions/io`, `@actions/exec`, `@goja-gha/ui`.

**Repository.** `~/code/wesen/go-go-golems/goja-github-actions`

**Public API.**
```js
const core = require("@actions/core")
core.getInput("name"); core.setOutput("k", v); core.setFailed("...")
core.summary.addRaw("...").write()
const io = require("@actions/io")     // cp/mv/rm/mkdirP, workspace-first
const exec = require("@actions/exec") // promise-based, stdout/stderr capture
const ui = require("@goja-gha/ui")    // report DSL for terminal summaries
```

**Implementation pattern.** Polyfill (Pattern E) for the `@actions/*` family, workspace-first path resolution. `@goja-gha/ui` is a small report DSL. Sync + async exports supported.

**Documentation.** `README.md` (root, very detailed — settings, doctor, examples, token scopes), `action.yml` (composite action), `examples/`, `integration/`.

**Evidence.** `README.md`, `pkg/`, `examples/`, `lib/`, `action.yml`.

---

### 4.10 goja-treesitter — skeleton (no DSL yet)

**Repository.** `~/code/wesen/go-go-golems/goja-treesitter`

**Status.** Stub only. `pkg/` contains `doc.go` and `logcopter.go`; no module exports, no `require()` target. README is ASCII art with no content. This is a placeholder for a future tree-sitter-backed DSL; nothing to catalogue yet.

**Evidence.** `pkg/doc.go`, `pkg/logcopter.go`, `README.md`.

---

### 4.11 glazed — Go-native CLI command builder (cross-cutting, the ancestor)

**Repository.** `…/glazed`

**What it is.** Not a Goja DSL, but the canonical Go-side fluent-builder grammar that informed the goja DSLs. Structured CLI commands over a `Schema` of `Section`s of `Field` `Definition`s, with a middlewares pipeline and multiple output formatters.

**Public API (Go).**
```go
schema.NewSchema(schema.WithSections(section1, section2))
fields.New("name", fields.TypeString,
    fields.WithHelp("..."),
    fields.WithDefault(def),
    fields.WithChoices("a","b"),
    fields.WithRequired(true),
    fields.WithIsArgument(true))
// Section interface: AddFields, GetDefinitions, GetName/Slug/Description/Prefix
// Middlewares pipeline + formatters (table, json, yaml, csv, markdown)
```

**Implementation pattern.** Functional options over typed structs (`Definition`, `Schema`, `Section` interface, `Definitions` collection). Strong typing, ordered maps for field order, cloneable. This is the model to lift into the goja playbook: typed structs + functional options + collection types.

**Documentation.** `README.md` (root, extensive), Glazed help system (`pkg/help/`), `pkg/cmds/schema/schema.go`, `pkg/cmds/fields/definitions.go`.

**Evidence.** `pkg/cmds/schema/schema.go:35-80` (`Schema`, `Section`, `NewSchema`, `WithSections`), `pkg/cmds/fields/definitions.go:20-90` (`Definition`, `Option` funcs, `New`).

---

### 4.12 go-emrichen — YAML tag-operator DSL (cross-cutting, the "composable grammar" model)

**Repository.** `~/code/wesen/go-go-golems/go-emrichen`

**What it is.** A YAML template engine where each tag (`!Defaults`, `!If`, `!Var`, `!Format`, `!Loop`, …) is a composable operator with type-aware validation. This is the closest existing thing to the ticket's "core composable grammar of operators extended with lambdas."

**Example.**
```yaml
!Defaults
isAdmin: true
ports: [80, 443]
---
accessLevel: !If
  test: !Var isAdmin
  then: "Full Access"
  else: "Restricted Access"
services: !Loop
  over: !Var ports
  as: port
  body: { port: !Var port }
```

**Implementation pattern.** Tag-dispatched operators, each implementing an evaluate interface. Strong type safety, type-aware operations, programmatically extensible (add custom tags/operators). Composable by nesting.

**Documentation.** `README.md`, `emrichen-spec.md`, `emrichen-in-practice.md`, `pkg/`, `test-data/`, `prompto/`.

**Evidence.** `README.md`, `emrichen-spec.md`, `emrichen-in-practice.md`.

---

### 4.13 geppetto — LLM runtime DSL with typed refs + DTS parity test (Pattern A′)

**Module name.** `require("geppetto")`

**Repository.** `~/code/wesen/go-go-golems/geppetto` (JS layer in `pkg/js/modules/geppetto/`).

**Public API (top-level exports, `module.go:176-190`).**
```js
const gp = require("geppetto")
gp.version
gp.inferenceProfiles          // profile registry namespace
gp.engine                      // engine builder: engine().inference(settings).build()
gp.embeddings                  // embeddings builder
gp.agent                       // agent builder: agent().name(...).engine(...).tools(...).build()
gp.tool                        // tool builder
gp.toolRegistry                // tool registry builder
gp.schema                      // JSON-schema builder namespace
gp.turnStores / gp.sessions / gp.consts / gp.events / ...
```

**Fluent examples.**
```js
// Engine builder (clone-on-each-step)
const eng = gp.engine().inference(settings).build()

// Schema builder (JSON schema)
const s = gp.schema.object()
  .description("a person")
  .property("name", gp.schema.string().description("full name"))
  .build()

// Agent builder
const agent = gp.agent().name("research").engine(eng).tools(reg).build()
```

**Implementation pattern.** Pattern A′. Same hidden-key substrate as goja-bleve (`hiddenRefKey = "__geppetto_ref"`, `attachRef`, `mustSet`, `requireInferenceSettingsRef`/type-asserting extractors). The key difference is **clone-on-each-step**: each fluent method clones the ref (`cloneFor(m)`) and returns a **new** wrapper, so builders are immutable per step — a safer variant for concurrent reuse than goja-bleve's same-object mutation. Validation is mixed (`panic(m.vm.NewTypeError(...))` for misuse, `(value, error)` for some terminals).

**The DTS parity test (the key contribution).** `pkg/js/modules/geppetto/dts_parity_test.go` defines `TestGeneratedDTSMatchesRuntimeExportSurface`, which parses the generated `pkg/doc/types/geppetto.d.ts` and asserts that (a) top-level exports and (b) grouped namespaces (`consts`, `inferenceProfiles`, `schema`, `turnStores`) match the runtime `require("geppetto")` export surface exactly. The `.d.ts` is generated by `cmd/gen-meta` from `geppetto_codegen.yaml`. **This is the compile-time-type enforcement discipline the playbook wants, with a test guaranteeing the declarations stay in sync.**

**Documentation.** `README.md` (root — runtime model, building blocks, profile registries, multimodal inputs), `pkg/doc/topics/` (`06-inference-engines.md`, `08-turns.md`), `pkg/doc/types/geppetto.d.ts` (generated declarations), extensive `_test.go` files.

**Evidence.** `pkg/js/modules/geppetto/module.go:28,73,176-190`, `api_engine_builder.go:18-70` (`engineBuilder`, clone-for, `inference`/`build`), `api_schema_builders.go:1-60` (`schema.string()/object()/enum()`, `.description`/`.property`), `api_agent.go:1-55` (`agentBuilderRef`), `dts_parity_test.go:1-60`, `pkg/doc/types/geppetto.d.ts:1-70`.

---

### 4.14 discord-bot — `defineBot` framework + Proxy-trap `ui` builders (Pattern G + registration DSL)

**Module names.** `require("discord")`, `require("ui")`, `require("timer")`, `require("database")`.

**Repository.** `~/code/wesen/go-go-golems/discord-bot`.

**`defineBot` registration DSL.**
```js
const { defineBot } = require("discord")
module.exports = defineBot(({ command, event, component, modal, autocomplete, configure }) => {
  configure({ name: "ping", description: "...", category: "examples" })
  command("ping", { description: "..." }, async (ctx) => { return { content: "pong" } })
  command("echo", { options: { text: { type: "string", required: true } } }, async (ctx) => { ... })
  subcommand(rootName, name, spec?, handler)
  userCommand(name, handler) / messageCommand(name, handler)
  event("ready", async (ctx) => { ... })
  component("ping:panel", async (ctx) => { ... })
  modal("feedback:submit", async (ctx) => { ... })
  autocomplete("search", "query", async (ctx) => { ... })
})
```

**`require("ui")` — Go-side typed builders via Proxy traps.**
```js
const ui = require("ui")
return ui.message()
  .content("Search results")
  .embed(ui.embed("Results").description("Found 3 items"))
  .row(ui.button("search:next", "Next", "primary"))
  .build()

ui.card(title?) / ui.select(customId) / ui.userSelect(customId) / ui.roleSelect(customId)
ui.channelSelect(customId) / ui.mentionableSelect(customId)
ui.form(customId, title)        // modal form builder
ui.row(...components) / ui.pager(prevId, nextId, controls?) / ui.actions(defs)
ui.confirm(confirmId, cancelId, options?) / ui.ok(content) / ui.error(content)
ui.emptyResults(query?) / ui.flow(namespace, options?)
```

**Implementation pattern.** Pattern G + registration DSL. `defineBot(builderFn)` is a registration DSL: the host calls the builder with a destructured set of registration helpers. `require("ui")` is the typed-builder DSL — Go-side fluent builders surfaced through **Goja Proxy traps**, returning **typed builders** (not plain JS objects) with a `.build()` terminal. Design rules: wrong-parent calls fail loudly (`ui.message().field(...)` errors with a hint to use `ui.embed()`); raw JS objects rejected where builders are expected; `.followUp()` vs update-in-place is explicit. xgoja provider at `pkg/xgoja/provider/`.

**Documentation.** `README.md` (root — install, quick start, env vars), `pkg/doc/topics/discord-js-bot-api-reference.md` (the full API reference with the `ui` builder table and design rules), `pkg/doc/tutorials/building-and-running-discord-js-bots.md`, `examples/discord-bots/` (ping, poker, knowledge-base, support, moderation).

**Evidence.** `pkg/framework/framework.go:1-90` (functional-options embedding API), `pkg/doc/topics/discord-js-bot-api-reference.md:40-160` (`ui` builder table + Proxy-trap description + design rules), `pkg/xgoja/provider/provider.go:36-42`, `examples/discord-bots/ping/index.js` (full defineBot example).

---

### 4.15 researchctl — claim/evidence/decision graph builder (Pattern F, lambda-configurator)

**Module name.** `require("researchctl")`

**Repository.** `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl` (JS layer in `pkg/gojamodules/researchctl/`).

**Public API.**
```js
const { project } = require("researchctl")

module.exports = project("Example JS project")
  .describe("...")
  .goal("Choose a backend", g => g.id("GOAL-001").status("active").priority("P1"))
  .question("...", q => q.id("Q-001").hypothesize("H-001").sources("S-001"))
  .hypothesis("Simulation gives enough signal", h => h.id("H-001").status("open").priority("P1").confidence("unknown"))
  .workPackage("...", w => w.id("WP-001").owner("...").dependsOn("WP-000"))
  .experiment("Run simulation", e => e.id("EXP-001").status("planned").priority("P1").tests("H-001"))
  .source("...", s => s.id("S-001"))
  .evidence("...", e => e.id("E-001"))
  .decision("...", d => d.id("D-001").confidence("high").reversalCondition("..."))
  .report("...", r => r.id("R-001"))
  .reviewRule("...", r => r.id("RR-001"))
  .view("...", v => v.id("V-001"))
  .toSpec()       // → ResearchProjectSpec
  .validate()     // → ValidationResult
```

**Entity sub-builders (lambda-configurator target).** Each `build` callback receives a fresh sub-builder: `id`, `description`, `priority`, `status`, `asks`, `tag`, `hypothesize`, `sources`, `testedBy`, `evidence`, `decision`, `confidence`, `reversalCondition`, `kind`, `owner`, `team`, `dependsOn`, `blocks`, `tests`, etc. (`builders.go:128-168`).

**Implementation pattern.** Pattern F. `project(name)` returns a `ProjectBuilder`; each collection method takes an optional lambda configurator applied to a fresh sub-builder. Same-object chaining within sub-builders. Terminals `.toSpec()` / `.validate()`. The TS declaration (`module.go:29`) emits a precise `ProjectBuilder` interface (with typed method signatures including the `build?: (g: any) => any` callback params) — note the callback param is still `any`-typed, which is the one weak spot. `fromSpec(v)` round-trips a spec back into a builder.

**Documentation.** `README.md` (root — graph model, YAML format, JS grammar format with the runnable example above), `pkg/research/` (graph model), `pkg/codesign/` (the codesign package), `examples/codesign/`, `pkg/gojamodules/researchctl/module_test.go`.

**Evidence.** `pkg/gojamodules/researchctl/module.go:29,57-90` (`project`, `fromSpec`, `validate`, `valueToSpec` with `toSpec()` duck-typing), `builders.go:11,128-168` (entity sub-builders), `module_test.go:28-33` (runnable example), `README.md`.

---

### 4.16 codesign — CPU/GPU run builder with `.use()` fragment composition (Pattern F, the strongest model)

**Module name.** `require("codesign")`

**Repository.** `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/`.

**Public API.**
```js
const codesign = require("codesign")
codesign.runSpec(name)            // → RunSpecBuilder
codesign.compareMetric(...) / reduceValues(...) / validateRun(...)
codesign.toSpec(...) / run(...) / manifest(...) / writeArtifacts(...)
codesign.configHash(...) / summarize(...) / context(...)
codesign.toYaml(...) / loadYaml(...) / writeYamlArtifact(...)
codesign.registerSource(...) / registerCallbackSource(...)
```

**`runSpec` fluent builder + lambda sub-builders.**
```js
const run = codesign.runSpec("my-run")
  .experiment("EXP-001")
  .backend("cpu-sim")
  .policy("min_finish_time")            // or .policy(type, config?)
  .policyCallback("my-policy", (task, candidateIds, state, scores) => "dev0")  // JS lambda
  .topology(t => t                       // lambda configurator
    .cpu("cpu0", { speed: 3.0, lanes: 4 })
    .accelerator("gpu0", { speed: 1.5 })
    .jsDevice("est", (phase, task, state, fallback) => estimate(task), { ... })  // JS lambda
  )
  .workload(w => w.fixed(1))            // lambda configurator
  .metrics(m => m                       // lambda configurator
    .latencyP95()
    .requestCount()
    .callback("custom", (events) => ({ value: events.length, unit: "events" }))  // JS lambda
  )
  .use(myReusableFragment)              // fragment composition: FragmentFn<RunSpecBuilder>
  .validate()                            // → ValidationResult
  .toSpec()                              // → RunSpec
  .run({ ... })                         // → RunResult
```

**Sub-builders.** `TopologyBuilder` (`.device/.cpu/.accelerator/.bandwidthDevice/.gpu/.appleMSeries/.nvidiaBlackwell/.disaggregatedMemory/.jsDevice/.use`), `WorkloadBuilder`, `MetricsBuilder` (`.metric/.callback/.latencyP95/.requestCount/.tasksByDevice/.taskTimeByDevice/.use`).

**Implementation pattern.** Pattern F — the strongest realisation of the ticket's full goal. `runSpec(name)` returns a `RunSpecBuilder`; `.topology(fn)`/`.workload(fn)`/`.metrics(fn)` apply a **lambda configurator** to a fresh sub-builder via `applyBuilderCallback`. `.use(fragment)` composes reusable `FragmentFn<T>` builder lambdas. **Lambdas cross the Go/JS boundary as runtime callbacks** (`jsDevice`, `policyCallback`, `callback` metrics), validated with `goja.AssertFunction`. Registers via `modules.Register(&module{})` and implements both `modules.NativeModule` and `modules.TypeScriptDeclarer`. **Precise TypeScript `interface`s** in `typescript.go`: `RunSpecBuilder`, `TopologyBuilder`, `MetricsBuilder`, `RunSpecLike = RunSpec | RunSpecBuilder | { toSpec(): RunSpec }` — typed unions, not `any`. Terminals: `.validate()` → `ValidationResult`, `.toSpec()` → `RunSpec`, `.run(options?)` → `RunResult`.

**Why this is the strongest model (evidence-backed).**
1. **Fluent builder** ✅ — `runSpec(name).experiment(id).backend(name).policy(type).build()`.
2. **Go implementation** ✅ — all in `pkg/gojamodules/codesign/`.
3. **Strict runtime typecheck** ✅ — typed builders, `RunSpecLike` union, `goja.AssertFunction` validation of callbacks.
4. **Validation at runtime** ✅ — `.validate()` → `ValidationResult`; errors returned.
5. **Compile-time types from generated declarations** ✅ — `TypeScriptDeclarer` interface, precise `interface` types (not `any`) in `typescript.go`.
6. **Composable grammar** ✅ — `.use(fragment)` with `FragmentFn<T>` reusable lambdas.
7. **Extended with lambdas** ✅ — `topology(fn)`/`workload(fn)`/`metrics(fn)` configurators + `jsDevice`/`policyCallback`/`callback` runtime callbacks.

**Documentation.** `pkg/gojamodules/codesign/typescript.go` (full TS declarations), `module.go` (exports + `TypeScriptDeclarer`), `builders.go` (runSpec/topology/workload/metrics builders), `module_test.go` (runnable examples), `pkg/codesign/` (Go domain: `spec`, `devices`, `devicefamilies`, `simulator`, `policies`, `metrics`, `sweeps`, `workloads`, `compare`, `registry`).

**Evidence.** `pkg/gojamodules/codesign/module.go:18-55` (`NativeModule` + `TypeScriptDeclarer`, exports), `builders.go:12-60` (`runSpec`/`runSpecBuilder`, `.experiment/.backend/.policy/.policyCallback/.topology/.workload`), `typescript.go:28-33` (`RunSpecBuilder`/`TopologyBuilder`/`MetricsBuilder`/`RunSpecLike` interfaces), `module_test.go:75,102,190` (runnable `.policy("min_finish_time")` examples).

## 5. Cross-cutting infrastructure

### 5.1 `engine` — runtime ownership and composition (go-go-goja)

The explicit runtime composition API that every DSL host uses:

1. `engine.NewRuntimeFactoryBuilder(...)` → add `WithModules(...)`, `WithRuntimeInitializers(...)` → `Build()` → immutable factory.
2. `factory.NewRuntime(engine.WithStartupContext(ctx), engine.WithLifetimeContext(ctx))` → owned runtime.
3. `rt.Close(ctx)` — explicit cleanup.

`WithStartupContext` controls construction + initializers; `WithLifetimeContext` controls runtime-owned resources + cancellation. Legacy `engine.New()`/`NewWithOptions()`/`Open()` were removed.

**Evidence.** `go-go-goja/pkg/engine/factory.go`, `options.go`, `runtime.go`, `module_middleware.go`, `module_roots.go`, `module_specs.go`.

**Relevance to playbook.** The fluent-builder DSL must respect runtime ownership: callbacks, `goja.Value` creation, Promise settlement, and EventEmitter delivery happen on the runtime owner (goja-dbus enforces this; it is the correctness model).

### 5.2 `xgoja` + `tsgen` — codegen and compile-time types (go-go-goja)

**xgoja** generates a host binary from an `xgoja.yaml` (see `goja-text/xgoja.yaml`): declares `packages`, `modules` (with `as` aliases and `config`), `commands` (eval/run/repl/jsverbs), `jsverbs` (embedded JS verb trees), and `help` sources. `go generate` runs `go tool xgoja build --work-dir . --dry-run` to refresh checked-in scaffold.

**tsgen/spec** is the TypeScript declaration model (`pkg/tsgen/spec/types.go`):
- `Module{Name, Description, Functions, RawDTS}`
- `Function{Name, Description, Params, Returns}`
- `Param{Name, Type, Optional, Variadic, Description}`
- `TypeRef{Kind, Name, Item, Union, Fields}` — full TS type system (string/number/boolean/any/unknown/void/never/named/array/union/object).
- `Field{Name, Type, Optional}`.

Each module implements `TypeScriptModule() *spec.Module`. A `Bundle` renders all modules into one `.d.ts`.

**Relevance to playbook.** This is the **compile-time types from generated declarations** substrate. The playbook must specify that every fluent builder and terminal emits a precise `TypeRef` (not `Record<string, any>`). goja-bleve's typed refs map cleanly onto named `TypeRef`s; widgetdsl's loose maps do not.

**Evidence.** `go-go-goja/pkg/tsgen/spec/types.go:1-80`, `pkg/tsgen/render/`, `pkg/tsgen/validate/`, `pkg/tsscript/`, `pkg/jsdoc/` (jsdoc extraction → spec), `goja-text/xgoja.yaml`.

### 5.3 `jsverbs` — JavaScript commands bridged to glazed (go-go-goja)

JavaScript verbs that compile to glazed `Command`s. `VerbSpec` → `Command`/`WriterCommand` wrapping glazed `CommandDescription` + `Schema`/`fields`. This is the bridge that lets JS-authored commands use the glazed parameter/middleware/output system.

**Evidence.** `go-go-goja/pkg/jsverbs/command.go:1-50`, `binding.go`, `model.go`, `runtime.go`, `scan.go`.

**Relevance.** Shows how a JS-facing DSL can compile down to the glazed typed-command grammar — relevant for the "composable grammar that compiles to a typed core" idea.

### 5.4 `jsdoc` — JSDoc → TypeScript spec extraction (go-go-goja)

Tooling (`pkg/jsdoc/` with `extract`, `batch`, `export`, `server`, `watch`, `model`) that extracts TypeScript specs from JSDoc comments on Go-exported functions. This is an alternative authoring path for the compile-time types: annotate the Go binding, generate the `.d.ts`.

**Evidence.** `go-go-goja/pkg/jsdoc/`.

## 6. Gap analysis — what is missing for the playbook

Measured against the ticket's goal (fluent builder, Go implementation, strict runtime typecheck, validation, compile-time types from generated declarations, composable grammar extensible with lambdas, opinionated):

| Requirement | Best existing model | Gap |
| --- | --- | --- |
| Fluent builder pattern | goja-bleve, goja-dbus, geppetto, researchctl, codesign | No shared builder library; each DSL reimplements `mustSet`/`newWrapper`/`getTypedRef` (or `attachRef`/`applyBuilderCallback`) |
| Strict runtime typechecking | goja-bleve `getTypedRef[T]` + `refKind`; codesign typed builders + `RunSpecLike` union | Still not extracted into a reusable `fluent` package; discord-bot `ui` uses Proxy traps instead — two mechanisms |
| Validation at runtime | go-minitrace `Validate()`/`Build()`; researchctl/codesign `.validate()`/`.toSpec()`; goja-bleve terminal `(v,error)` | Three converged styles now; geppetto/widgetdsl still panic — needs a rule |
| Compile-time types | codesign `TypeScriptDeclarer` with precise `interface`s; geppetto DTS parity test | **No longer the main gap.** codesign/geppetto show the way; widgetdsl still emits `Record<string, any>` |
| Composable grammar of operators | codesign `.use(fragment)` + `FragmentFn<T>`; go-emrichen tag operators | codesign already realises this; needs extraction into the core grammar |
| Lambdas as extensibility | researchctl/codesign lambda configurators `g => g.id(...)`; codesign `jsDevice`/`policyCallback` runtime callbacks | **Realised.** The pattern exists end-to-end; needs documenting as the canonical extensibility model |
| Opinionated single way | — | **The main remaining gap.** The mechanisms now exist in code; the playbook must pick one substrate (hidden-key vs Proxy traps) and one composition model (`.use()` fragments + lambda configurators) |
| Runtime ownership correctness | goja-dbus (callbacks on owner); geppetto (`runtimeowner.RuntimeOwner` in Options) | Documented only in dbus ticket; geppetto reinforces it; not generalized |
| Lifecycle (close/cleanup) | goja-bleve batch single-use, goja-dbus `bus.close()`, geppetto turn stores | Pattern exists but undocumented as a rule |

**Concrete missing pieces for the playbook to define:**
1. A reusable `fluent` package: `refBase`/`refKind`/`newWrapper`/`getTypedRef[T]` extracted from goja-bleve.
2. A validation discipline: accumulate errors (go-minitrace style) and surface them at `.Validate()`/`.Build()` terminals returning `(value, error)`.
3. A TypeScript emission rule: every builder and terminal emits a named `TypeRef`, not `any`.
4. A composable-operator grammar: borrow go-emrichen's tag-operator model so operators can nest and accept lambdas (`walk`-style callbacks as first-class operators).
5. A lifecycle rule: builders that own resources expose `.close()`/single-use semantics.
6. A runtime-ownership rule: all callbacks/Promises settle on the runtime owner.

## 7. Decision records (for the senior researcher to ratify or reject)

### Decision: runtime typecheck substrate

- **Context:** Three substrates exist — hidden-key typed refs (goja-bleve), side-channel `sync.Map` (express), and none (widgetdsl).
- **Options considered:** (a) hidden-key typed refs + `getTypedRef[T]`; (b) `sync.Map` side-channel; (c) untyped maps.
- **Decision:** *proposed* — adopt goja-bleve's hidden-key + `getTypedRef[T]` as the substrate for the playbook.
- **Rationale:** it is the only one that gives both JS-chainability and Go-side type enforcement with a single source of truth. Side-channel maps work but duplicate identity; untyped maps abandon typecheck.
- **Consequences:** enables strict runtime typecheck; requires extracting and sharing the machinery; builders must not escape their wrapper.
- **Status:** proposed.

### Decision: validation discipline

- **Context:** go-minitrace accumulates + `Validate()`; goja-bleve validates at terminals with `(v,error)`; widgetdsl panics.
- **Options considered:** (a) accumulate + `Validate()` terminal; (b) validate-at-terminal `(v,error)`; (c) panic.
- **Decision:** *proposed* — combine: accumulate errors during chaining, validate at `.Build()` returning `(value, error)`, never panic.
- **Rationale:** panic is hostile to JS callers; accumulation lets one `Build()` report all errors.
- **Consequences:** builders carry an `errors []string`; terminals must return `(T, error)`.
- **Status:** proposed.

### Decision: compile-time type emission

- **Context:** `tsgen/spec` is capable of named types, but widgetdsl emits `any`.
- **Options considered:** (a) every builder emits named `TypeRef` + precise `Param`/`Returns`; (b) open `Props`; (c) JSDoc extraction.
- **Decision:** *proposed* — require (a); JSDoc extraction (c) may augment but not replace.
- **Rationale:** the whole point of "compile-time types from generated declarations" is that the `.d.ts` actually narrows.
- **Consequences:** each module's `TypeScriptModule()` must model builder chains as named types; more authoring discipline.
- **Status:** proposed.

### Decision: composable operator grammar + lambdas

- **Context:** The follow-up pass found that researchctl and codesign **already realise** this: codesign's `.use(fragment)` with `FragmentFn<T>` reusable builder lambdas, and both DSLs' lambda-configurator pattern `g => g.id(...).status(...)` applied via `applyBuilderCallback`. go-emrichen remains the conceptual tag-operator model.
- **Options considered:** (a) adopt codesign's `.use()` + lambda-configurator pattern as the canonical composition model; (b) adopt emrichen-style tag operators in JS; (c) build a new `operator` core from scratch.
- **Decision:** *proposed* — adopt (a): the playbook's composable grammar is the lambda-configurator + `.use(fragment)` pattern proven by researchctl/codesign, not a new filter/map/reduce core.
- **Rationale:** it is already in production code, already emits precise TS types, and already crosses the Go/JS boundary safely with `goja.AssertFunction`. Option (c) would reinvent what codesign already ships.
- **Consequences:** the playbook must specify `FragmentFn<T>` and `applyBuilderCallback` as core primitives; lambda param typing in TS is still `any` in researchctl — the playbook should narrow this (codesign already narrows via `interface`s).
- **Status:** proposed (strengthened — there is now a proven implementation to point at).

## 8. Open questions for the senior researcher

1. Should the playbook extract goja-bleve's typed-ref machinery into a shared `fluent` package in go-go-goja, or keep it per-module? (Tradeoff: shared dep vs. duplication.)
2. How should builder chains be modeled in TypeScript — as a builder interface that narrows per method, or as a single type with overloaded signatures? (Affects `.d.ts` readability.)
3. Should validation be eager (at each chain step) or lazy (at `.Build()`)? go-minitrace is lazy; goja-bleve is eager on some terminals.
4. What is the canonical lifecycle rule — `close()` vs. finalizer vs. context cancellation? goja-dbus and goja-bleve batch differ.
5. How do lambdas cross the Go/JS boundary safely — `goja.FunctionCall` directly, or a typed `Lambda` wrapper with input/output `TypeRef`s? (codesign validates with `goja.AssertFunction`; researchctl passes the sub-builder as the lambda's single arg.)
6. Should the playbook be prescriptive (one pattern) or present a small decision tree (Pattern A for type-heavy, Pattern B for pipeline-heavy)?
7. Is `go-emrichen`'s tag-operator model in scope for the JS DSL, or only conceptual inspiration?
8. What is the testing pattern? goja-bleve has golden declaration tests; go-minitrace has builder unit tests; **geppetto has a DTS parity test**; no shared convention.
9. Should the widgetdsl `data.dsl` be migrated to the new pattern, or left as-is and the playbook only governs new DSLs?
10. How does the playbook interact with the `engine` ownership model — does every builder carry a runtime owner ref? (geppetto passes `runtimeowner.RuntimeOwner` in `Options`.)
11. **Hidden-key typed refs (goja-bleve/geppetto) vs Proxy-trap builders (discord-bot `ui`) — which is the canonical typed-builder mechanism?** Proxy traps give natural JS ergonomics but are harder to introspect from Go; hidden keys are simpler but require `getTypedRef` discipline.
12. **Clone-on-each-step (geppetto) vs same-object mutation (goja-bleve) — which is the canonical builder mutation model?** Geppetto's immutable-per-step is safer for reuse; goja-bleve's same-object is less allocating.
13. Should the `FragmentFn<T>` + `.use()` composition (codesign) be lifted into a shared `fragments` package, or kept per-DSL?

### Decision: typed-builder mechanism (hidden-key vs Proxy traps)

- **Context:** Two mechanisms now coexist — hidden-key typed refs (`__bleve_ref`/`__geppetto_ref` + `getTypedRef[T]`, goja-bleve/geppetto) and Goja Proxy traps (discord-bot `ui`).
- **Options considered:** (a) hidden-key + `getTypedRef[T]` as canonical; (b) Proxy traps as canonical; (c) support both.
- **Decision:** *proposed* — adopt (a) hidden-key + `getTypedRef[T]` as the default; document Proxy traps as an optional ergonomics layer for builder-heavy UIs.
- **Rationale:** hidden keys are simpler, introspectable from Go, and already shared by two DSLs; Proxy traps are powerful but add complexity and are harder to type-declare. Defaulting to one reduces cognitive load.
- **Consequences:** the playbook must specify `attachRef`/`getTypedRef[T]`/`mustSet` as the core; Proxy traps become an advanced section.
- **Status:** proposed.

## 9. References (key files)

### goja-bleve
- `~/code/wesen/go-go-golems/goja-bleve/pkg/api_types.go` — `refBase`, `refKind`, `getTypedRef[T]`, `newWrapper`
- `~/code/wesen/go-go-golems/goja-bleve/pkg/api_mapping.go` — `fieldBuilder`, `installFieldOptions`
- `~/code/wesen/go-go-golems/goja-bleve/pkg/module.go` — `hiddenRefKey`, `mustSet`, exports
- `~/code/wesen/go-go-golems/goja-bleve/README.md`, `docs/faiss-xgoja-playbook.md`

### goja-dbus
- `~/code/wesen/go-go-golems/goja-dbus/pkg/dbusgoja/builders.go` — bus/method fluent builders
- `~/code/wesen/go-go-golems/goja-dbus/pkg/dbusgoja/signals.go` — signal subscription builder
- `~/code/wesen/go-go-golems/goja-dbus/pkg/dbusgoja/typed_values.go` — typed value helpers
- `~/code/wesen/go-go-golems/goja-dbus/README.md`

### go-minitrace
- `go-minitrace/pkg/minitracejs/module.go` — factory exports
- `go-minitrace/pkg/minitracejs/builders.go` — `SourceSetBuilder`
- `go-minitrace/pkg/minitracejs/db_builder.go` — `DBBuilder`, `ValidationResult`, `DBHandle`
- `go-minitrace/pkg/minitracejs/import_builder.go` — `ImportBuilder` + terminals
- `go-minitrace/pkg/minitracejs/query_view_session.go` — `QueryRecipeBuilder`

### rag-evaluation-system widgetdsl
- `rag-evaluation-system/pkg/widgetdsl/module.go` — module specs + helpers
- `rag-evaluation-system/pkg/widgetdsl/grammar.go` — data grammar verbs
- `rag-evaluation-system/pkg/widgetdsl/typescript.go` — TS declaration gen
- `rag-evaluation-system/pkg/widgetdsl/module_test.go`, `grammar_test.go`

### go-go-goja
- `go-go-goja/modules/uidsl/` — hyperscript element DSL
- `go-go-goja/modules/express/express.go`, `auth_builders.go` — route + auth builders
- `go-go-goja/pkg/engine/` — runtime ownership
- `go-go-goja/pkg/tsgen/spec/types.go` — TS spec model
- `go-go-goja/pkg/jsverbs/` — JS command → glazed bridge
- `go-go-goja/pkg/jsdoc/` — JSDoc → spec extraction
- `go-go-goja/README.md`

### goja-text
- `~/code/wesen/go-go-golems/goja-text/README.md`, `xgoja.yaml`

### geppetto
- `~/code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/module.go` — exports, `hiddenRefKey`, `Options` (runtime owner)
- `~/code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/api_engine_builder.go` — clone-on-each-step engine builder
- `~/code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/api_schema_builders.go` — JSON-schema builder
- `~/code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/api_agent.go` — agent builder
- `~/code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/dts_parity_test.go` — **DTS parity test** (compile-time-type enforcement)
- `~/code/wesen/go-go-golems/geppetto/pkg/doc/types/geppetto.d.ts` — generated declarations
- `~/code/wesen/go-go-golems/geppetto/README.md`

### discord-bot
- `~/code/wesen/go-go-golems/discord-bot/pkg/framework/framework.go` — functional-options embedding API
- `~/code/wensen/go-go-golems/discord-bot/pkg/doc/topics/discord-js-bot-api-reference.md` — `defineBot` + `ui` builder reference
- `~/code/wensen/go-go-golems/discord-bot/pkg/doc/tutorials/building-and-running-discord-js-bots.md`
- `~/code/wensen/go-go-golems/discord-bot/pkg/xgoja/provider/provider.go` — xgoja provider
- `~/code/wensen/go-go-golems/discord-bot/examples/discord-bots/ping/index.js` — full defineBot example

### researchctl + codesign
- `~/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/README.md` — graph model + JS grammar format
- `~/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/researchctl/module.go` — `project()`, `fromSpec`, `validate`, TS `ProjectBuilder` interface
- `~/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/researchctl/builders.go` — entity sub-builders
- `~/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/module.go` — `NativeModule` + `TypeScriptDeclarer`, exports
- `~/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/builders.go` — `runSpec`/topology/workload/metrics builders, `.use()` fragments
- `~/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/typescript.go` — precise `RunSpecBuilder`/`TopologyBuilder`/`MetricsBuilder`/`RunSpecLike` interfaces
- `~/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/module_test.go` — runnable examples

### glazed + go-emrichen
- `glazed/pkg/cmds/schema/schema.go`, `glazed/pkg/cmds/fields/definitions.go`
- `~/code/wesen/go-go-golems/go-emrichen/README.md`, `emrichen-spec.md`, `emrichen-in-practice.md`

### Related docmgr tickets
- `RAGEVAL-UI-DSL`, `RAGEVAL-UI-GRAMMAR`, `RAGEVAL-CMS-WIDGETS`, `RAGEVAL-CONTEXT-WINDOWS-DESIGN`
- `RAGEVAL-WIDGET-IR-SEMANTIC-COMPONENTS`, `DESIGN-REF-001`, `CTX-WINDOW-BLOCK-VIZ`, `CTX-COLOR-PALETTE`
- `rag-eval-scripting-expansion` (xgoja + goja-text jsverbs)
- `GOJA-DBUS-DESIGN` (in goja-dbus repo ttmp)
