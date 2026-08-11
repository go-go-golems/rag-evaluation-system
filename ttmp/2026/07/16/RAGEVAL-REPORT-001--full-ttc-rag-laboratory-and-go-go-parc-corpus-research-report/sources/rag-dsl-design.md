---
Title: Typed fluent RAG module design and implementation guide
Ticket: RAGEVAL-RAG-DSL-001
Status: active
Topics:
    - rag
    - rag-eval
    - dsl
    - fluent-builder
    - goja
    - xgoja
    - javascript
    - typescript
    - intern-guide
    - playground
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://rag-evaluation-system/cmd/rag-eval/xgoja.yaml
      Note: Runtime packaging point for the future xgoja provider
    - Path: repo://rag-evaluation-system/internal/db/db.go
      Note: Defines the SQLite immutable experiment and append-only run schema constraints
    - Path: repo://rag-evaluation-system/internal/services/experimentrun/service.go
      Note: Persistence service that the module must reuse rather than duplicate
    - Path: repo://rag-evaluation-system/internal/services/immutableretrieval/bm25.go
      Note: Existing lexical retrieval service to adapt for execution
    - Path: repo://rag-evaluation-system/internal/services/immutableretrieval/vector.go
      Note: Existing vector RRF collapse and source hydration service to adapt
ExternalSources: []
Summary: Intern-oriented design and delivery guide for implementing the typed Go-backed require("rag") module, xgoja provider, canonical experiment specifications, and explicit immutable-run execution boundary.
LastUpdated: 2026-07-14T22:09:50.311752843-04:00
WhatFor: Implement the RAG laboratory JavaScript module safely, consistently, and without creating a second experiment persistence model.
WhenToUse: Read before beginning the pure Go builder, Goja adapter, xgoja provider, runtime tests, or execution integration.
---


# Typed fluent RAG module design and implementation guide

## 0. Executive summary

This ticket introduces a reusable JavaScript RAG laboratory module with the
entry point `require("rag")`. Its purpose is to make RAG experimentation fast
to author without making it informal, untraceable, or detached from the
existing immutable experiment system.

The implementation consists of three deliberately separated layers:

1. a pure Go domain package that holds typed builders, compatibility checks,
   canonical serialisation, and fingerprints;
2. a thin Goja NativeModule adapter that maps lower-camel JavaScript calls to
   those builders and turns validation problems into normal thrown JS errors;
3. an xgoja provider package that exposes the module to generated
   `rag-eval-js` binaries and generates a matching TypeScript declaration.

The domain package produces the **same canonical experiment specification**
that `internal/services/experimentrun` persists and associates with an
append-only experiment run. There is one identity model, one database schema,
one run history, and one web UI inspector. JavaScript is a second authoring
surface, not an alternate RAG engine.

The recommended first vertical slice is intentionally narrow: raw-chunk
retrieval over already-materialised TTC artifacts, BM25 and vector channels,
RRF fusion, parent/document collapse, retrieval metrics, `.toSpec()`,
`.validate(lab)`, `lab.persist()`, and `lab.start()`. Summary and synthetic
question representations use the same plan model later, after their artifact
production and cache identities are sound.

## 1. Reading order for a new intern

Read in this order. It avoids mistaking the old generic JavaScript playground
for the intended RAG API.

1. Read [the normative API reference](../reference/01-rag-laboratory-javascript-module-api-specification.md).
   Treat it as the desired user-visible behavior and example suite.
2. Read `docs/guides/ttc-rag-laboratory.md` for current operator terminology
   and artifact/run inspection workflows.
3. Read `internal/services/experimentrun/service.go` and its tests. This is
   the persistence authority for immutable specifications, runs, events,
   summaries, and query traces.
4. Read `internal/services/immutableretrieval/bm25.go`, `vector.go`, and
   their tests. These own baseline retrieval, deterministic sorting, RRF, and
   source hydration behavior.
5. Read `internal/db/db.go` to see the SQLite tables and append-only trigger
   constraints. Do not introduce a parallel JavaScript-owned table.
6. Read `cmd/rag-eval/xgoja.yaml` and
   `docs/howtos/how-to-write-rag-eval-js-scripts.md` to understand how native
   modules become available in `rag-eval-js`.
7. Read the Widget DSL's current builder documentation and its use of nested
   builder lambdas and `.use(fragment)`. Its ergonomics are useful precedent,
   but its browser-IR target is not the RAG module's target.
8. Read the transcript research prototype's `rag.js` and `bleve-rag.js` under
   `../2026-07-09--transcript-rag-sol2/.../scripts/playground/verbs/lib/`.
   Extract retrieval invariants (channel-local retrieval, RRF, collapse,
   hydration); do not import its transcript-specific code as a dependency.

## 2. The problem to solve

The existing `rag-eval-js` runtime is useful for exploration but it exposes
generic primitives: SQLite access through `db`, filesystem operations through
`fs`, text parsing through goja-text, and inference through `geppetto`.
`cmd/rag-eval/jsverbs/database.js` and `explorer.js` demonstrate this model.
An author can query tables or construct an HTTP helper quickly, but the runtime
does not know whether a query used a compatible embedding set, whether RRF was
applied deterministically, or which experiment identity should own the result.

The laboratory now has immutable artifact identities and append-only runs. A
useful RAG playground must retain these guarantees rather than make a second,
mutable collection of results. At the same time, authors should be able to
express a comparison at the level they actually reason about:

```js
rag.experiment("vector-vs-rrf", (e) =>
  e.use(ttcInputs).retrieval((r) => r.channel("semantic", (c) => c.vector().topK(50)))
)
```

The implementation challenge is therefore not vector-search syntax. It is the
boundary between convenient JavaScript composition and typed reproducible Go
state.

## 3. Scope and non-goals

### 3.1 In scope

- A native `rag` module for the generated `rag-eval-js` runtime.
- Typed Go builders and JSON-safe experiment/study specifications.
- Immutable artifact-reference compatibility validation.
- Raw representation, lexical/vector channels, RRF, filters, collapse, and
  source/citation hydration as spec-level concepts.
- Metric threshold and metric selection configuration.
- Explicit `persist` and append-only `start` calls.
- xgoja provider registration, TypeScript declarations, examples, help, and
  runtime integration coverage.
- The first execution adapter that maps the validated spec to existing
  immutable retrieval and experiment-run services.

### 3.2 Explicitly out of scope

- Replacing the Go HTTP API or the React laboratory UI.
- A general-purpose JavaScript vector database.
- Direct SQL in a public DSL primitive.
- A generic callback called for every chunk/query/result.
- Provider credentials, network configuration, or direct Geppetto calls in a
  retrieval spec.
- Summary/question generation during `build`, validation, or run submission.
- Backward-compatible wrappers around pre-existing generic jsverbs.
- A stable high-level agent framework.

The no-callback rule is important. A JS callback can be useful at authoring
time, but persisting an arbitrary closure would not preserve its code,
dependencies, capability boundary, or deterministic behavior. Named,
versioned Go extension points can be added later where customization proves
necessary.

## 4. Current architecture and target architecture

### 4.1 Current state

```text
JavaScript script
   │ require("db"), require("geppetto"), require("markdown")
   ▼
generic xgoja modules / jsverbs
   │ SQL or ad-hoc calls
   ▼
SQLite + service code + optional external provider
```

The current path is appropriate for one-off corpus diagnosis. It has no
first-class typed experiment object at the authoring boundary.

### 4.2 Target state

```text
                    ┌────────────────────────────────────┐
                    │ `require("rag")` JavaScript surface │
                    │ experiment / fragment / study        │
                    └────────────────────────────────────┘
                                      │ lower-camel codecs
                                      ▼
 ┌────────────────────────────────────────────────────────────────────┐
 │ pkg/raglab (pure Go)                                                │
 │ typed builders → ValidationIssue[] → canonical spec → fingerprint  │
 └────────────────────────────────────────────────────────────────────┘
                  │ read-only artifact lookup           │ explicit effect
                  ▼                                     ▼
 ┌───────────────────────────────┐      ┌───────────────────────────────┐
 │ immutable artifact services   │      │ experimentrun service          │
 │ snapshots, chunks, embeddings │      │ spec + append-only run/events  │
 └───────────────────────────────┘      └───────────────────────────────┘
                  │                                     │
                  └───────────────┬─────────────────────┘
                                  ▼
                     immutable retrieval executor
                   BM25 / vector / RRF / collapse / traces
                                  │
                                  ▼
                    existing API + React laboratory UI
```

The pure package is the center. It cannot import Goja. This enables ordinary
Go tests and allows Go CLI, web handlers, or a future YAML loader to use the
same validator and canonicaliser.

## 5. Key design decisions

### Decision: compile to the existing immutable specification

- **Context:** The TTC laboratory already persists content-addressed
  specifications and append-only runs.
- **Options considered:** Create JavaScript-specific tables; serialise a script
  and execute it later; compile to the current specification.
- **Decision:** Compile to the current canonical `ExperimentSpecification`.
- **Rationale:** One query trace, metric, artifact, and UI inspection model is
  easier to trust and compare.
- **Consequences:** The Go persisted schema becomes the compatibility boundary.
  The initial implementation may need to make its JSON schema explicit before
  it exposes the module.
- **Status:** accepted.

### Decision: use configurator lambdas, never persisted behavior lambdas

- **Context:** Widget DSL and researchctl make nested configuration readable.
- **Options considered:** plain nested option objects; persistent JS callbacks;
  authoring-time configuration lambdas.
- **Decision:** Use synchronous lambdas only to mutate a typed builder during
  script evaluation.
- **Rationale:** This preserves familiar fluent ergonomics without embedding
  opaque executable behavior in an experiment identity.
- **Consequences:** Custom chunking/scoring must be a named, versioned Go
  extension rather than `map(chunk => ...)` in v1.
- **Status:** accepted.

### Decision: explicit effect boundary

- **Context:** A RAG script can accidentally cause expensive provider work or
  produce misleading experiment records.
- **Options considered:** automatic execution on `.build()`; a separate CLI;
  explicit methods on the laboratory handle.
- **Decision:** `.toSpec()` and `.validate()` are pure; `lab.persist()` and
  `lab.start()` are explicit effects, guarded by `execution: "allowRuns"`.
- **Rationale:** A user can inspect the exact plan/fingerprint before compute.
- **Consequences:** A small amount of extra syntax buys reliable dry-run and
  read-only tooling.
- **Status:** accepted.

### Decision: represent artifacts as typed opaque references

- **Context:** Hash strings are easy to mix up and difficult to validate at
  the point a plan is authored.
- **Options considered:** bare strings everywhere; Go-backed mutable model
  objects; typed references plus convenience string overloads.
- **Decision:** Provide `ArtifactRef` with `kind` and `id`; accept string IDs
  only in semantically typed builder positions.
- **Rationale:** JavaScript is concise while the Go builder can check intent
  and compatibility.
- **Consequences:** `rag.artifact()` is available for reusable helpers; a
  reference must be JSON-safe and cannot carry an open database handle.
- **Status:** accepted.

### Decision: one execution-capable xgoja runtime, safe by default

- **Context:** The existing generic xgoja runtime includes database and
  network-capable modules, but project loaders may not have that authority.
- **Options considered:** always make `rag` side-effectful; split it into two
  module names; one module with opt-in execution configuration.
- **Decision:** Use one `rag` module whose `Laboratory` is read-only by
  default. Expose it only in runtimes whose xgoja config intentionally selects
  it; execution is enabled explicitly.
- **Rationale:** The authoring API remains coherent and capability review stays
  in generated-binary configuration.
- **Consequences:** Tests must prove read-only rejection and allowed execution.
- **Status:** accepted.

## 6. Domain model and Go package layout

Create these packages; names are suggested but the separation is mandatory.

```text
pkg/raglab/
  artifacts.go        # ArtifactKind, ArtifactRef, compatibility metadata
  grades.go           # named ordinal relevance grades
  specification.go    # JSON-safe ExperimentSpecification + StudyPlan
  builder.go          # ExperimentBuilder, fragment application, local checks
  retrieval.go        # channels, filters, RRF, collapse, deterministic sorting spec
  metrics.go          # metrics builder and cutoff canonicalisation
  canonical.go        # canonical JSON and SHA-256 fingerprint
  validation.go       # Issue, Report, Error
  laboratory.go       # narrow service interfaces, not database/sql
  executor.go         # spec-to-existing-service orchestration
  *_test.go

pkg/gojamodules/rag/
  module.go           # NativeModule, Loader, exports
  codecs.go           # Goja values/options to domain builder calls
  errors.go           # RagValidationError conversion
  typescript.go       # TypeScriptDeclarer contract
  module_test.go      # require("rag") integration tests

pkg/xgoja/providers/rag/
  rag.go              # provider Register and module descriptor
  rag_test.go         # registry/module/dts smoke tests

examples/rag-lab-js/
  01-bm25-baseline.js
  02-vector-baseline.js
  03-hybrid-rrf.js
  04-study-representations.js
```

If `internal/services/experimentrun` exposes an adequate interface already,
`pkg/raglab/laboratory.go` should depend on a small interface rather than on
the concrete database service. For example:

```go
type ArtifactCatalog interface {
    Inspect(ctx context.Context, ref ArtifactRef) (ArtifactMetadata, error)
}

type ExperimentStore interface {
    CreateSpecification(ctx context.Context, spec ExperimentSpecification) (PersistedSpecification, error)
    CreateRun(ctx context.Context, specificationID string, options RunOptions) (RunHandle, error)
}

type RetrievalExecutor interface {
    Submit(ctx context.Context, runID string, spec ExperimentSpecification) error
}
```

Every concrete Go type that implements one of these interfaces MUST carry a
compile-time assertion such as:

```go
var _ ArtifactCatalog = (*SQLiteArtifactCatalog)(nil)
```

The pure builder itself should not need `context.Context`; external lookup,
persistence, and execution methods do.

## 7. Data flow in detail

### 7.1 Build phase

```text
JS e.retrieval(lambda) ──► Goja invokes lambda synchronously
                                │
                                ▼
                       RetrievalBuilder mutates typed fields
                                │
JS experiment.toSpec() ─────────┤
                                ▼
                    normalize: validate local invariants
                    canonicalize: sort unordered collections
                    hash: SHA-256 canonical payload
                                │
                                ▼
                      JSON-safe ExperimentSpecification
```

Pseudocode:

```go
func (b *ExperimentBuilder) ToSpec() (ExperimentSpecification, error) {
    issues := b.ValidateStructural()
    if issues.HasErrors() {
        return ExperimentSpecification{}, NewValidationError(issues)
    }

    spec := b.normalizedSpecification()
    payload := canonicalJSON(spec.withoutFingerprint())
    spec.Fingerprint = sha256ID(payload)
    return spec, nil
}
```

Normalization must not look at the clock, a random source, map iteration
order, network state, or database rows. `notes`, `tags`, and fragment names
need an explicit canonical policy. Recommended policy: fragment names are
ordered in application order because provenance order can be explanatory;
tags are key-sorted; cutoff and unordered ID collections are key-sorted and
deduplicated; channels retain authoring order because RRF component display may
use it but execution tie breaking must never depend on it.

### 7.2 Validation phase

There are two levels of validation.

| Level | Where | Examples | Effect free? |
|---|---|---|---|
| Structural | `ExperimentBuilder.ToSpec()` | missing inputs, duplicate channel, impossible cutoff | yes |
| Compatibility | `Experiment.Validate(lab)` / `persist` / `start` | artifact exists, snapshot matches, vector dimensions match | read-only lookup |

Pseudocode:

```go
func (lab *Laboratory) Validate(ctx context.Context, spec ExperimentSpecification) ValidationReport {
    report := validateStructuralSpec(spec)
    if report.HasErrors() { return report }

    snapshot := lab.catalog.Inspect(ctx, spec.Inputs.CorpusSnapshot)
    chunks := lab.catalog.Inspect(ctx, spec.Inputs.ChunkSet)
    report.AddAll(checkChunkSetSnapshot(chunks, snapshot))

    for _, channel := range spec.Retrieval.Channels {
        if channel.Backend == Vector {
            embeddings := lab.catalog.Inspect(ctx, *spec.Inputs.EmbeddingSet)
            report.AddAll(checkEmbeddingCompatibility(embeddings, chunks, channel.Representation))
        }
    }
    return report
}
```

The validator should collect as many actionable diagnostics as possible. It
should not submit a run when an error exists. Warnings (for example a declared
representation that no channel selects) can be displayed and stored in run
submission metadata but must not change the fingerprint.

### 7.3 Persistence/execution phase

```text
spec ──► lab.persist(spec)
             │ compatibility validation
             ▼
        experiment_specs INSERT OR find same immutable fingerprint
             │
spec ──► lab.start(spec)
             │ persist/fetch immutable spec
             ▼
        experiment_runs INSERT (new run id)
             │ append submitted event
             ▼
        durable executor receives (run id, spec id)
             │
             ├─ per query: embed/query → channels → fuse → collapse → hydrate
             ├─ append query trace
             ├─ append events
             └─ append terminal summary
```

`start()` must create the run before scheduling the durable job, so an enqueue
failure is recorded as a terminal failed run or a clear failed submission
event. It must not leave a silent untracked side effect.

## 8. Retrieval execution requirements

The domain spec is only useful if execution preserves its stated semantics.
The initial executor should adapt existing immutable retrieval code rather than
duplicate vector mathematics in the Goja adapter.

For each evaluation query:

```text
for each retrieval channel:
    derive query text / embedding exactly once when reusable
    execute backend with the channel's filter and topK
    attach channel name, backend score, one-based rank, representation identity

if more than one channel:
    fuse candidates with deterministic RRF
else:
    preserve the one channel score and rank

collapse using declared scope
hydrate each survivor to original document + source chunk + citation metadata
truncate to final result count
grade ranked sources against fixed evaluation card
append immutable trace record

aggregate metrics / latency / storage / provider-cost counters
append terminal summary
```

Important implementation constraints:

- Filters must be applied before the top-K cutoff whenever a backend supports
  it. Filtering the first unfiltered K candidates is a recall bug.
- Vector query embeddings should be computed once per query/model identity and
  cache hits/misses must become trace metrics, not invisible behavior.
- RRF starts each fused candidate at zero. Never copy a raw channel score into
  an RRF score. This was already identified as a deterministic correctness
  issue in the baseline work.
- Original representations and raw sources must not share a collision-prone
  identifier. Use a typed retrieval item identity plus parent/source IDs.
- Collapse follows fusion unless the specification explicitly gains
  channel-local collapse later. V1 has one documented order: retrieve → fuse →
  collapse → hydrate → truncate.
- Trace records must carry enough information to reconstruct why an item
  ranked: channel ranks/scores, fusion components, representation, parent chunk,
  original source document, and final citation.

## 9. Goja adapter design

The Goja package does conversion and no domain policy. It should look like a
native module, not a second builder implementation.

```go
type module struct {
    factory func(OpenOptions) (*raglab.Laboratory, error)
}

var _ modules.NativeModule = (*module)(nil)
var _ modules.TypeScriptDeclarer = (*module)(nil)

func (m *module) Name() string { return "rag" }

func (m *module) Loader(vm *goja.Runtime, moduleObj *goja.Object) {
    exports := moduleObj.Get("exports").(*goja.Object)
    _ = exports.Set("open", m.open(vm))
    _ = exports.Set("experiment", m.experiment(vm))
    _ = exports.Set("fragment", m.fragment(vm))
    _ = exports.Set("study", m.study(vm))
    _ = exports.Set("artifact", m.artifact(vm))
    _ = exports.Set("grade", m.grade(vm))
    _ = exports.Set("version", "v1")
}
```

Adapter responsibilities:

- decode plain JS option objects with lower-camel names;
- reject `undefined`, functions, or objects of the wrong shape at an exported
  API boundary with a clear `TypeError`/validation error;
- invoke configurator functions synchronously with a wrapped Go builder;
- trap Go panics and translate known domain validation failures to a JS error
  carrying `code` and `issues`;
- convert `ExperimentSpecification`, validation reports, inspections, and run
  handles to ordinary lower-camel plain objects;
- avoid giving JS a reflected `database/sql.DB`, service, or mutable Go struct.

The adapter must explicitly define callback failure semantics. If a lambda
throws, the outer builder call throws the same JavaScript exception and the
partially configured builder is discarded. It must not persist partial state,
wrap an error into a generic string, or call the callback again.

### Example adapter bridge pseudocode

```go
func configureExperiment(vm *goja.Runtime, b *raglab.ExperimentBuilder, fn goja.Callable) error {
    wrapped := newExperimentObject(vm, b)
    if _, err := fn(goja.Undefined(), wrapped); err != nil {
        return errors.Wrap(err, "experiment configurator failed")
    }
    return nil
}

func exportError(vm *goja.Runtime, err error) {
    var validation *raglab.ValidationError
    if errors.As(err, &validation) {
        obj := vm.NewTypeError(validation.Error())
        _ = obj.Set("code", validation.Code())
        _ = obj.Set("issues", exportIssues(vm, validation.Issues))
        panic(obj)
    }
    panic(vm.NewGoError(err))
}
```

Use `github.com/pkg/errors` for Go error wrapping, following repository
guidelines. Do not use a panic to express ordinary domain control flow outside
the normal Goja exported-function mechanism.

## 10. xgoja provider and TypeScript packaging

`modules.Register()` is not enough for a generated xgoja binary. The provider
must register `rag` with `providerapi.ProviderRegistry`, supply a module loader,
and expose the same TypeScript descriptor used by module-level declarations.

```text
pkg/gojamodules/rag ── registers native module ──► direct Go embedding
          │
          ▼
pkg/xgoja/providers/rag ── Register() ──► xgoja provider registry
          │                                      │
          ├── TypeScript module descriptor        ├── `xgoja gen-dts`
          └── module factory                      └── generated rag-eval-js
```

Add a provider import to `cmd/rag-eval/xgoja.yaml` and select it in
`runtime.modules` with `as: rag`. Do not make the module implicitly available
to unrelated runtime binaries. The generated runtime should document whether
the default operation is `readOnly` and how an operator obtains an
execution-capable laboratory.

The declaration generator must cover function/method overloads enough that
these fail early in TypeScript:

```ts
rag.grade("2");                         // invalid grade spelling
e.retrieval((r) => r.channel("x", (c) => c.topK("10"))); // wrong argument type
```

Runtime validation remains authoritative.

## 11. Implementation plan and checkpoints

Implement in this order. Commit at every independently reviewable checkpoint.

### Checkpoint A — freeze the interchange schema

1. Inspect actual structs/JSON in `internal/services/experimentrun` and
   `internal/api/experiment_handlers.go`.
2. Define or extract an explicit `rag-eval.experiment/v1` schema.
3. Write fixtures for the current TTC lexical, vector, and RRF baseline.
4. Decide whether fingerprint field is stored or derived at read time.
5. Commit schema/fixtures/tests alone.

Acceptance: canonical fixture JSON is stable and matches a persisted baseline
specification.

### Checkpoint B — pure Go builder

1. Add `pkg/raglab` structs and builders.
2. Implement fragments, inputs, channels, RRF, filters, collapse, metrics,
   named grades, canonicalisation, and fingerprinting.
3. Add validation reports and structural error tests.
4. Keep all tests free of Goja/SQLite where possible.
5. Commit domain package and tests.

Acceptance: constructing the same semantic plan in different harmless input
orders has the same fingerprint; incompatible configurations have stable issue
codes and paths.

### Checkpoint C — artifact compatibility and persistence adapter

1. Define narrow catalog/store interfaces.
2. Adapt immutable artifact and experiment-run services.
3. Implement validate, persist, and read-only enforcement.
4. Test against a temporary SQLite database with realistic artifact metadata.
5. Commit service boundary and tests.

Acceptance: a valid raw TTC plan persists idempotently and `start` creates a
new run without mutating prior rows.

### Checkpoint D — Goja module

1. Implement `pkg/gojamodules/rag`.
2. Write runtime integration tests using `engine.New()` and
   `require("rag")`.
3. Add simple script fixtures for every API-reference quick example.
4. Ensure builder callback errors remain JavaScript errors.
5. Commit module and tests.

Acceptance: `toSpec()` JavaScript output matches Go fixture JSON byte for byte
after canonical formatting.

### Checkpoint E — provider, declarations, and examples

1. Implement provider registration and declaration descriptor.
2. Update xgoja config and generated-binary packaging.
3. Add `examples/rag-lab-js` scripts.
4. Run xgoja doctor, declaration generation, and binary/script smoke tests.
5. Commit packaging/docs/examples.

Acceptance: a generated `rag-eval-js` can execute
`examples/rag-lab-js/01-bm25-baseline.js` and `xgoja gen-dts` produces a
declaration containing `declare module "rag"`.

### Checkpoint F — execution vertical slice

1. Add a spec-to-retrieval executor that calls immutable retrieval services.
2. Write append-only events, traces, and summaries.
3. Smoke test raw BM25, vector, and RRF against the existing 20-card TTC
candidate dataset.
4. Link run/spec inspection from the web UI and document operator steps.
5. Commit the complete vertical slice.

Acceptance: a script can start a raw baseline experiment, the UI can inspect
its immutable spec and trace, and rerunning creates a different run ID.

## 12. Testing matrix

| Layer | Test | Purpose |
|---|---|---|
| Canonicalisation | same semantic plan, reordered inputs | stable fingerprint |
| Builder | missing/duplicate/conflicting settings | quality of diagnostics |
| Compatibility | mismatched snapshot/chunk/dimension | prevent invalid runs |
| Retrieval | BM25/vector/RRF/collapse fixtures | preserve algorithm semantics |
| Persistence | persist twice/start twice | immutable idempotence vs new runs |
| Goja | `require("rag")` + lambda error | public runtime contract |
| xgoja provider | registry resolution, DTS descriptor | generated-binary contract |
| Generated binary | `xgoja doctor`, `gen-dts`, example run | packaging is real |
| Web/API | inspect started run/spec | observable result, not hidden work |

Commands, adjusted for the repository's actual generated paths:

```bash
GOWORK=off go test ./pkg/raglab/... ./pkg/gojamodules/rag/... ./pkg/xgoja/providers/rag/... -count=1
GOWORK=off go test ./... -count=1
GOWORK=off go build ./...
xgoja doctor -f cmd/rag-eval/xgoja.yaml
xgoja gen-dts -f cmd/rag-eval/xgoja.yaml --out /tmp/rag-eval-js.d.ts
pnpm --dir web typecheck
pnpm --dir web build
```

When an execution test needs a local server, use tmux and inspect it with
`capture-pane`; kill a stale port with `lsof-who -p <port> -k`, per the
repository instructions. Do not use an interactive server as the only proof of
correctness.

## 13. Risks and mitigation

| Risk | Mitigation |
|---|---|
| A DSL creates another mutable experiment format | Compile only to existing canonical spec; no JS tables. |
| Fluent API grows into a vague wrapper | Keep v1 opinionated; unsupported operations require decision records. |
| Lambdas make plans irreproducible | Only authoring-time configurators; never persist functions. |
| Different callers create different hashes for equivalent plans | One canonicalisation package with golden fixtures. |
| Side effects happen in safe analysis tooling | Default read-only lab and xgoja module selection. |
| RRF/collapse semantics drift from prototype | Unit fixtures and trace-level assertions, not just aggregate metrics. |
| Generated declaration differs from runtime | Treat the API reference examples as both Goja and DTS tests. |
| Artifact IDs are valid syntax but incompatible | Catalog-based validation before persistence or execution. |

## 14. Open questions to resolve during Checkpoint A

1. Which exact Go type should become the schema authority: an extracted public
   package type or the existing service type with a dedicated JSON view?
2. Should persisted `candidate:ttc-baseline-v1` evaluation-data identifiers be
   migrated to a hash-form identity before the DSL's first release, or does the
   artifact catalog need a typed non-hash stable-ID allowance? The API reference
   allows both but the decision must be explicit.
3. Does `lab.start()` submit through the existing workflow service immediately,
   or should v1 ship `persist` only while the executor is completed? Do not
   publish a method that reports successful work without durable scheduling.
4. Do channel-level filters have an exact backend implementation for Bleve and
   the vector scan? If not, omit the public filter method until semantics are
   correct.
5. Should a study plan itself get a fingerprint and stored identity, or remain
   an exported convenience object until batch execution has budgets and
   concurrency controls?

## 15. Reference file map

| File | Why it matters |
|---|---|
| `cmd/rag-eval/xgoja.yaml` | Existing runtime and provider/module selection. |
| `cmd/rag-eval/jsverbs/database.js` | Baseline generic DB JavaScript surface. |
| `cmd/rag-eval/jsverbs/explorer.js` | Current corpus exploration verbs. |
| `docs/howtos/how-to-write-rag-eval-js-scripts.md` | Operator-facing xgoja/jsverb explanation. |
| `internal/services/experimentrun/service.go` | Immutable specification/run persistence authority. |
| `internal/experimentspec/specification.go` | Shared schema, manifest, normalization, and fingerprint contract. |
| `pkg/raglab/types.go` | Typed artifact, retrieval, representation, metric, validation, and specification model. |
| `pkg/raglab/builder.go` | Pure fluent Go builder and deterministic structural validation. |
| `internal/services/immutableretrieval/bm25.go` | Lexical retrieval implementation. |
| `internal/services/immutableretrieval/vector.go` | Vector/RRF/collapse baseline behavior. |
| `internal/db/db.go` | Schema and append-only database constraints. |
| `pkg/raglab/catalog.go` | Read-only artifact lineage validation contract. |
| `pkg/raglab/catalog_sqlite.go` | SQLite implementation of immutable artifact lookup. |
| `internal/api/experiment_handlers.go` | Existing laboratory API transport layer. |
| `docs/guides/ttc-rag-laboratory.md` | Current operator guide. |
| `ttmp/.../RAGEVAL-TTC-LAB-001/...` | Baseline artifacts, evaluation protocol, and immutable-run decisions. |
| `ttmp/.../GOJA-DSL-PLAYBOOK/...` | Widget/researchctl-fluent builder design precedent. |

## 16. Completion definition

The ticket is ready for review when all of the following are true:

- API reference examples run against a generated `rag-eval-js` binary.
- Go and TypeScript declarations describe the same supported calls.
- A script can produce an inspectable canonical raw TTC specification without
  write or provider authority.
- A permitted script can persist that specification and start a new append-only
  run through the normal service path.
- The run records retriever channels, fusion, collapse, original-source
  citations, metrics, latency, and relevant cost/storage counters in traces
  and summary data.
- Repeat execution reuses the specification fingerprint but creates a distinct
  run ID; no earlier run is changed.
- Documentation, examples, tests, xgoja provider packaging, and web inspection
  links are committed together.
