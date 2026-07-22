---
Title: RAG laboratory JavaScript module API specification
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
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://rag-evaluation-system/cmd/rag-eval/xgoja.yaml
      Note: Defines the current generated rag-eval-js runtime into which require("rag") will be packaged
    - Path: repo://rag-evaluation-system/docs/howtos/how-to-write-rag-eval-js-scripts.md
      Note: Documents the generic JavaScript runtime the typed API will complement
    - Path: repo://rag-evaluation-system/internal/services/experimentrun/service.go
      Note: Owns immutable experiment specifications and append-only run persistence consumed by the API
    - Path: repo://rag-evaluation-system/internal/services/immutableretrieval/vector.go
      Note: Defines vector retrieval RRF collapse and hydration behavior the DSL must express
ExternalSources: []
Summary: Normative proposed v1 contract for the typed Go-backed require("rag") module used to author and execute immutable RAG laboratory experiments.
LastUpdated: 2026-07-14T22:09:50.341418734-04:00
WhatFor: Give JavaScript authors and Go implementers one unambiguous contract for RAG laboratory scripts, generated TypeScript declarations, validation, and execution semantics.
WhenToUse: Read before writing a rag-eval-js experiment, implementing a rag module export, or reviewing a public API change.
---


# RAG laboratory JavaScript module API specification

## 1. Status, normative language, and purpose

This is the **proposed v1 public contract** for `require("rag")`. It is the
source of truth for implementation and generated TypeScript declarations. The
module is not implemented at the time of writing. Examples below are intended
syntax tests: once the module exists, every compact example should execute in
an integration test or live under `examples/rag-lab-js/`.

The key words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** are normative.

The module makes a RAG laboratory script reproducible. It does four things:

1. resolves and validates immutable artifact references;
2. constructs a typed experiment specification from readable fluent calls;
3. produces canonical JSON and a deterministic specification fingerprint;
4. explicitly persists a specification or starts a new append-only run.

It does not silently ingest data, generate embeddings, call an LLM, modify an
index, or mutate an existing run merely because a builder method is called.
Those operations are future execution capabilities and must remain explicit.

## 2. Module entry point and mental model

```js
const rag = require("rag");
const lab = rag.open({ database: "data/rag-eval.db" });
```

`rag` is a native Goja module selected by the `rag-eval-js` xgoja
configuration. `rag.open()` returns a `Laboratory`, a small handle bound to a
database and optional execution policy. It is not a global singleton. A script
may open more than one laboratory, which makes test databases and comparison
tools straightforward.

```text
JavaScript authoring calls
        │
        ▼
typed Go builder graph ── validate ──► canonical ExperimentSpecification
        │                                      │
        │                                      ├──► fingerprint (SHA-256)
        │                                      └──► JSON / YAML inspection
        ▼
explicit Laboratory effect
  persist(spec) or start(spec)
        │
        ▼
immutable specification row + new append-only run row
```

The diagram is intentionally asymmetric. Authoring is pure. Persistence and
execution are effects. That separation is what lets an author inspect the
exact object before consuming compute or producing a result that another
experiment must explain.

## 3. The public surface at a glance

| Export | Returns | Purpose |
|---|---|---|
| `rag.open(options)` | `Laboratory` | Open an explicit laboratory context. |
| `rag.fragment(name, configure)` | `Fragment` | Reusable authoring-time builder fragment. |
| `rag.experiment(name, configure)` | `Experiment` | Create one experiment builder. |
| `rag.study(name, configure)` | `Study` | Create a named collection of related experiments. |
| `rag.artifact(kind, id)` | `ArtifactRef` | Create a typed opaque immutable-artifact reference. |
| `rag.grade(name)` | `RelevanceGrade` | Choose a named relevance grade for metric thresholding. |
| `rag.version` | string | Public module-contract version, initially `"v1"`. |

Every builder method returns its receiver unless the method name begins with
`to`, `build`, `validate`, `inspect`, `persist`, `start`, or `compare`. That
keeps linear scripts readable and lets TypeScript preserve the exact builder
type after every call.

### 3.1 `rag.open(options)`

```ts
interface OpenOptions {
  database: string;                 // required, SQLite database path
  execution?: "readOnly" | "allowRuns";
  queryEmbed?: (query: string) => number[];
}
```

`database` is required deliberately. The generic `db` module permits a
default database path for ad-hoc inspection; an experiment API should not
accidentally write into whichever database happened to be configured earlier
in the process. The module validates that a database path is non-empty. It
does not open or mutate a database until an operation needs artifact lookup,
`persist`, or `start`.

`execution` defaults to `"readOnly"`. In that mode, `persist()` and `start()`
throw `RAG_EXECUTION_DISABLED`. A generated project-analysis runtime SHOULD
select only this mode. An operator runtime may select `"allowRuns"`.

`queryEmbed` is optional only for lexical plans. A vector channel requires a
synchronous callback returning a non-empty finite `number[]`; it is called for
each evaluation query at execution time. The callback is runtime capability,
not immutable plan data, so it can construct a Geppetto provider without
placing hostnames, credentials, or model-server transport into a fingerprint.

```js
const lab = rag.open({
  database: "data/rag-eval.db",
  execution: "allowRuns",
  defaultTopK: 10,
});
```

### 3.2 `rag.artifact(kind, id)`

```ts
type ArtifactKind =
  | "corpusSnapshot"
  | "chunkSet"
  | "embeddingSet"
  | "bm25Index"
  | "evaluationDataset";

interface ArtifactRef<K extends ArtifactKind = ArtifactKind> {
  readonly kind: K;
  readonly id: string;
  toJSON(): { kind: K; id: string };
}
```

Artifact references are opaque named values, not strings spread through an
experiment. The `id` normally has the form `sha256:<hex>`. A runtime call does
not claim that a supplied artifact exists; `validate()` and any effectful
method resolve it against the laboratory database.

```js
const source = rag.artifact(
  "corpusSnapshot",
  "sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409",
);
```

For convenience, `ExperimentBuilder` has corresponding typed methods such as
`.corpus(id)` and `.embeddings(id)`. `rag.artifact()` is most useful in shared
fragments, helper functions, tests, and tooling that operates on a generic
artifact kind.

### 3.3 `rag.grade(name)`

The baseline evaluation contract stores an ordinal grade and gives it a clear
human-readable name. V1 exports the following grades:

| Name | Ordinal | Meaning |
|---|---:|---|
| `"0_FAIL"` | 0 | Does not support the query or is materially wrong. |
| `"1_PARTIAL"` | 1 | Contains a related fragment but misses material evidence. |
| `"2_SUBSTANTIAL"` | 2 | Provides material support for a useful answer. |
| `"3_AUTHORITATIVE"` | 3 | Direct, complete, authoritative support. |

```js
metrics.relevanceAt(rag.grade("2_SUBSTANTIAL"));
```

`rag.grade()` throws `RAG_UNKNOWN_GRADE` for an unrecognised name. A numerical
grade is intentionally not accepted in scripts: the name documents the
threshold at its point of use. The canonical spec stores both the stable name
and its ordinal, so a reader does not need an external lookup table to
interpret old results.

## 4. Core experiment builder

### 4.1 Constructor and terminal operations

```ts
rag.experiment(
  name: string,
  configure?: (experiment: ExperimentBuilder) => void,
): Experiment;

interface Experiment {
  readonly name: string;
  use(fragment: Fragment): this;
  corpus(snapshot: string | ArtifactRef<"corpusSnapshot">): this;
  chunks(chunkSet: string | ArtifactRef<"chunkSet">): this;
  bm25(index: string | ArtifactRef<"bm25Index">): this;
  embeddings(set: string | ArtifactRef<"embeddingSet">): this;
  evaluation(dataset: string | ArtifactRef<"evaluationDataset">): this;
  representations(configure: (r: RepresentationBuilder) => void): this;
  retrieval(configure: (r: RetrievalBuilder) => void): this;
  metrics(configure: (m: MetricsBuilder) => void): this;
  note(text: string): this;
  tag(name: string, value?: string): this;
  validate(lab?: Laboratory): ValidationReport;
  toSpec(): ExperimentSpecification;
  toJSON(): ExperimentSpecification;
}
```

The optional configurator follows the Widget DSL and researchctl pattern: it
receives the typed builder, configures it synchronously, and returns nothing.
Returning the builder is harmless but ignored. The callback is an
**authoring-time callback**. V1 MUST NOT store it, invoke it later, or use it
to transform individual documents or search results at execution time.

`.toSpec()` performs local structural validation and returns plain JSON-safe
data. It does not need a database. It throws an aggregate `RagValidationError`
when required structure is missing or contradictory. `.validate(lab)` adds
artifact existence and compatibility validation when a laboratory is supplied;
it returns a report so a user interface or a CLI can display multiple issues.

### 4.2 Minimal lexical experiment

```js
const rag = require("rag");
const lab = rag.open({ database: "data/rag-eval.db" });

const experiment = rag.experiment("ttc-bm25-baseline", (e) =>
  e
    .corpus("sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409")
    .chunks("sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392")
    .bm25("sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691")
    .evaluation("candidate:ttc-baseline-v1")
    .retrieval((r) => r.channel("lexical", (c) => c.bm25().topK(50)).results(10))
    .metrics((m) => m.relevanceAt(rag.grade("2_SUBSTANTIAL")).recallAt([1, 3, 10]).mrr()),
);

const report = experiment.validate(lab);
if (!report.ok) throw new Error(JSON.stringify(report.issues, null, 2));

const spec = experiment.toSpec();
console.log(spec.fingerprint);
```

The build chooses a channel candidate depth of 50 but asks for ten final
ranked results. The two numbers have different meanings and the builder keeps
them separate.

### 4.3 Reusable input fragment

`Fragment` is a deterministic, authoring-only unit of composition. It is
inspired by `widget.dsl` `.use(fragment)` and researchctl component bundles.
It gives a shared corpus foundation a meaningful name without forcing a
factory function or duplicating artifact IDs.

```ts
rag.fragment(
  name: string,
  configure: (experiment: ExperimentBuilder) => void,
): Fragment;
```

```js
const ttcRawBaseline = rag.fragment("ttc-raw-2026-07-14", (e) =>
  e
    .corpus("sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409")
    .chunks("sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392")
    .bm25("sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691")
    .embeddings("sha256:2665c5249b8352ce6904fc00c934534dd179f3eeef0a6a75429a9034be0e03e0")
    .evaluation("candidate:ttc-baseline-v1"),
);

const vector = rag.experiment("ttc-vector", (e) =>
  e.use(ttcRawBaseline).retrieval((r) => r.channel("semantic", (c) => c.vector().topK(50))),
);
```

`use()` MUST apply the fragment immediately and record its name in the spec's
provenance. It MUST reject a fragment that sets a single-valued field to a
different value already set by the experiment. This avoids silent ordering
semantics such as "the last `chunks()` wins." Re-applying the same fragment is
idempotent.

## 5. Retrieval plan API

### 5.1 Shape

```ts
interface RetrievalBuilder {
  channel(name: string, configure: (channel: ChannelBuilder) => void): this;
  filter(configure: (filter: FilterBuilder) => void): this;
  fuse(configure: (fusion: FusionBuilder) => void): this;
  collapse(scope: "none" | "parentChunk" | "document"): this;
  results(count: number): this;
}

interface ChannelBuilder {
  bm25(): this;
  vector(): this;
  representation(name: string): this;
  topK(count: number): this;
  filter(configure: (filter: FilterBuilder) => void): this;
}
```

An experiment has one retrieval plan. A plan contains named candidate
channels, optional filters, a fusion algorithm, an optional collapse policy,
and a final result count. Channel names are unique, lower-case identifiers.
They become trace component names, so names such as `"lexical"`, `"semantic"`,
`"raw"`, `"summaries"`, and `"questions"` are preferred.

Each channel selects one backend (`bm25` or `vector`) and one representation.
The default representation is `"raw"`. The implementation rejects a second
backend selection, missing `.topK()`, duplicate channel names, a `topK` lower
than `results`, and vector retrieval without an embedding artifact.

### 5.2 Pure vector retrieval

```js
e.retrieval((r) =>
  r
    .channel("semantic", (c) => c.vector().representation("raw").topK(50))
    .collapse("document")
    .results(10),
);
```

With one channel, `.fuse()` is optional. The output preserves the vector score
as `score`, records the channel rank in `components.semantic.rank`, and stores
the selected representation and source/parent identity in each trace result.

### 5.3 Hybrid retrieval with reciprocal rank fusion

```js
e.retrieval((r) =>
  r
    .channel("lexical", (c) => c.bm25().representation("raw").topK(50))
    .channel("semantic", (c) => c.vector().representation("raw").topK(50))
    .fuse((f) => f.rrf().rankConstant(60))
    .collapse("document")
    .results(10),
);
```

RRF uses the deterministic score:

```text
fusedScore(item) = Σ channels where item appears [ weight(channel) / (rank + rankConstant) ]
```

`rank` is one-based. The v1 default weight is 1 and the default `rankConstant`
is 60. Ties MUST be resolved by original source ID, then parent chunk ID, then
representation item ID; Go map iteration order MUST NOT affect a ranking.

```ts
interface FusionBuilder {
  rrf(): this;
  rankConstant(value: number): this;
  weight(channel: string, value: number): this;
}
```

`weight()` is only valid after `.rrf()` and only for an existing channel. It
must be a finite positive number. The builder stores explicit non-default
weights, keeping the default implicit in canonical JSON.

### 5.4 Filters

```ts
interface FilterBuilder {
  sourceIds(ids: string[]): this;
  documentIds(ids: string[]): this;
  contentTypes(types: string[]): this;
  metadataEquals(key: string, value: string): this;
}
```

Filters are restrictive predicates ANDed together. Repeated calls of the same
type union values and canonicalise them in lexical order; repeated metadata
keys require the same value or throw. V1 deliberately has no arbitrary SQL
predicate, JavaScript callback, or raw filter string. Those forms would make a
spec impossible to validate, unsafe to run, and difficult to compare.

```js
e.retrieval((r) =>
  r
    .filter((f) => f.contentTypes(["product", "article"]).metadataEquals("locale", "en_US"))
    .channel("lexical", (c) =>
      c.bm25().topK(100).filter((f) => f.sourceIds(["ttc-wordpress"])),
    )
    .results(10),
);
```

Plan-level filters apply to all channels. Channel filters further restrict only
that channel. The execution adapter MUST push filters into the retrieval
backend before candidate selection when the backend can support it; filtering
after a top-K retrieval changes recall and is therefore semantically wrong.

### 5.5 Collapse and hydration

Collapse controls the identity that competes in the final ranking:

| Scope | Deduplication key | Intended use |
|---|---|---|
| `"none"` | representation item | diagnose chunk/representation behavior |
| `"parentChunk"` | original source chunk | compare raw, summary, and question representations fairly |
| `"document"` | original source document | answer-oriented retrieval and concise inspection |

The executor always hydrates final results to original-source evidence. A
summary or generated question may retrieve a parent chunk, but returned
citations must name the original document and show the raw parent chunk text,
with representation provenance retained separately.

```js
e.retrieval((r) =>
  r
    .channel("summary-vector", (c) => c.vector().representation("summary").topK(80))
    .collapse("parentChunk")
    .results(10),
);
```

This is not an optional presentational feature. Parent collapse and evidence
hydration prevent a synthetic representation from receiving multiple ranks for
the same evidence and prevent a UI from presenting generated text as if it
were a source citation.

## 6. Representation API

### 6.1 V1 representation model

The first implementation should support raw chunks only, because the raw TTC
baseline already has compatible chunk, BM25, and embedding artifacts. The
public contract reserves structured representation declarations so that
summary/question experiments do not require redesign later.

```ts
interface RepresentationBuilder {
  rawChunks(name?: "raw"): this;
  summaries(name: string, configure: (s: SummaryRepresentationBuilder) => void): this;
  questions(name: string, configure: (q: QuestionRepresentationBuilder) => void): this;
}

interface SummaryRepresentationBuilder {
  artifact(set: string): this;
  parent("sourceChunk"): this;
}

interface QuestionRepresentationBuilder {
  artifact(set: string): this;
  parent("sourceChunk"): this;
}
```

`rawChunks()` is implicitly present in every experiment. Calling it explicitly
is harmless and is useful for a self-documenting script. `summaries` and
`questions` describe **already materialised immutable representation sets**.
They do not run a provider. A later module or explicit command may generate
such a set and return its artifact identity.

```js
e.representations((r) =>
  r
    .rawChunks()
    .summaries("summary", (s) =>
      s.artifact("sha256:summary-representation-set").parent("sourceChunk"),
    )
    .questions("question", (q) =>
      q.artifact("sha256:question-representation-set").parent("sourceChunk"),
    ),
);
```

The implementation MUST validate that each declared representation set has the
same corpus snapshot and compatible chunk set as the experiment. It MUST
reject a summary/question representation with no explicit parent relation.

### 6.2 Multi-representation hybrid example

```js
const allRepresentations = rag.experiment("ttc-all-representations", (e) =>
  e
    .use(ttcRawBaseline)
    .representations((r) =>
      r
        .summaries("summary", (s) => s.artifact("sha256:summary-set").parent("sourceChunk"))
        .questions("question", (q) => q.artifact("sha256:question-set").parent("sourceChunk")),
    )
    .retrieval((r) =>
      r
        .channel("raw-lexical", (c) => c.bm25().representation("raw").topK(50))
        .channel("raw-semantic", (c) => c.vector().representation("raw").topK(50))
        .channel("summary-semantic", (c) => c.vector().representation("summary").topK(50))
        .channel("question-semantic", (c) => c.vector().representation("question").topK(50))
        .fuse((f) => f.rrf().rankConstant(60).weight("raw-semantic", 1.25))
        .collapse("parentChunk")
        .results(10),
    )
    .metrics((m) => m.relevanceAt(rag.grade("2_SUBSTANTIAL")).recallAt([1, 3, 10]).mrr()),
);
```

The script is a plan, not a claim that the named summary artifacts exist.
`validate(lab)` checks availability before an expensive run starts.

## 7. Evaluation and metrics API

```ts
interface MetricsBuilder {
  relevanceAt(grade: RelevanceGrade): this;
  precisionAt(cutoffs: number[]): this;
  recallAt(cutoffs: number[]): this;
  hitRateAt(cutoffs: number[]): this;
  ndcgAt(cutoff: number): this;
  mrr(): this;
  meanRelevantRecallAt(cutoff: number): this;
  abstention(): this;
}
```

At least one metric MUST be selected. If a graded retrieval metric is selected,
`relevanceAt()` MUST be selected first. The builder de-duplicates and sorts
cutoffs. Thus `.recallAt([10, 1, 3, 3])` canonicalises to `[1, 3, 10]`.

```js
e.metrics((m) =>
  m
    .relevanceAt(rag.grade("2_SUBSTANTIAL"))
    .precisionAt([1, 3, 10])
    .recallAt([1, 3, 10])
    .hitRateAt([1, 3, 10])
    .ndcgAt(10)
    .mrr()
    .meanRelevantRecallAt(10)
    .abstention(),
);
```

`abstention()` records coverage/abstention handling for answerability-aware
evaluation; v1's retrieval executor reports it even when it cannot yet perform
answer generation. It is an honest measurement flag, not a permission to
discard hard queries.

## 8. Canonical specification, inspection, and validation

### 8.1 Canonical JSON contract

`.toSpec()` returns a plain object which is JSON serializable without special
handling. Its persisted immutable-contract schema version is
`rag-eval-experiment-spec/v1`. The JavaScript projection uses lower-camel keys;
the storage manifest uses the equivalent snake-case Go JSON fields.

```ts
interface ExperimentSpecification {
  schemaVersion: "rag-eval-experiment-spec/v1";
  fingerprint: string;                  // sha256:<hex>, computed from canonical payload
  name: string;
  provenance: { fragments: string[]; notes: string[]; tags: Record<string, string> };
  inputs: {
    corpusSnapshot: ArtifactRefJson;
    chunkSet: ArtifactRefJson;
    bm25Index?: ArtifactRefJson;
    embeddingSet?: ArtifactRefJson;
    evaluationDataset: ArtifactRefJson;
    representations: RepresentationSpec[];
  };
  retrieval: RetrievalPlan;
  metrics: MetricsPlan;
}
```

`fingerprint` MUST be calculated from the canonical payload **without** the
fingerprint field. Canonicalisation uses sorted object keys, deterministic
array ordering where order has no semantic meaning, UTF-8 JSON, and SHA-256.
The experiment name is part of the payload; two differently named experiments
are intentionally different specimens even when their retrieval plan matches.

An abbreviated shape is:

```json
{
  "schemaVersion": "rag-eval-experiment-spec/v1",
  "fingerprint": "sha256:...",
  "name": "ttc-vector-and-rrf",
  "inputs": {
    "corpusSnapshot": { "kind": "corpusSnapshot", "id": "sha256:..." },
    "chunkSet": { "kind": "chunkSet", "id": "sha256:..." },
    "embeddingSet": { "kind": "embeddingSet", "id": "sha256:..." },
    "evaluationDataset": { "kind": "evaluationDataset", "id": "candidate:ttc-baseline-v1" },
    "representations": [{ "name": "raw", "kind": "rawChunks" }]
  },
  "retrieval": {
    "channels": [{ "name": "semantic", "backend": "vector", "representation": "raw", "topK": 50 }],
    "collapse": "document",
    "results": 10
  },
  "metrics": { "relevanceAt": { "name": "2_SUBSTANTIAL", "ordinal": 2 }, "mrr": true }
}
```

### 8.2 Validation report

```ts
interface ValidationReport {
  ok: boolean;
  fingerprint?: string;
  issues: ValidationIssue[];
}

interface ValidationIssue {
  code: string;
  path: string;       // JSONPath-like, e.g. $.retrieval.channels[1].topK
  message: string;
  severity: "error" | "warning";
}
```

Structural errors are thrown by `.toSpec()` because no usable spec exists.
`.validate(lab)` returns all local and database-backed diagnostics and does not
start execution. The standard `RagValidationError` has `name`, `code`, and
`issues` fields; the Goja adapter presents it as a normal thrown JS `Error`.

```js
const report = experiment.validate(lab);
for (const issue of report.issues) {
  console.log(`${issue.severity} ${issue.code} at ${issue.path}: ${issue.message}`);
}
```

Important error codes include:

| Code | Meaning |
|---|---|
| `RAG_MISSING_INPUT` | Required immutable input was not selected. |
| `RAG_DUPLICATE_CHANNEL` | Channel names must be unique. |
| `RAG_CHANNEL_BACKEND_CONFLICT` | A channel selected both BM25 and vector. |
| `RAG_MISSING_EMBEDDINGS` | Vector retrieval lacks a compatible embedding set. |
| `RAG_MISSING_BM25` | BM25 retrieval lacks a compatible lexical index. |
| `RAG_UNKNOWN_REPRESENTATION` | Channel references an undeclared representation. |
| `RAG_INCOMPATIBLE_ARTIFACT` | Artifact snapshots/chunk sets/dimensions do not match. |
| `RAG_INVALID_CUTOFF` | A metric cutoff is non-positive or exceeds final results. |
| `RAG_CONFLICTING_FRAGMENT` | Two fragments chose different single-valued inputs. |
| `RAG_EXECUTION_DISABLED` | A read-only laboratory attempted a side effect. |
| `RAG_SCHEMA_UNSUPPORTED` | A persisted spec cannot be executed by this module version. |

## 9. Persisting and executing a run

```ts
interface Laboratory {
  inspect(ref: ArtifactRef | { kind: ArtifactKind; id: string }): ArtifactInspection;
  persist(experiment: Experiment | ExperimentSpecification): PersistedSpecification;
  start(experiment: Experiment | ExperimentSpecification, options?: StartOptions): RunHandle;
  execute(experiment: Experiment): ExecutionResult;
  compare(leftRunId: string, rightRunId: string): RunComparison;
}
```

### 9.1 Persist is explicit and idempotent by fingerprint

```js
const lab = rag.open({ database: "data/rag-eval.db", execution: "allowRuns" });
const persisted = lab.persist(experiment);

console.log(persisted.id);          // canonical specification id / fingerprint
console.log(persisted.created);     // false when the same fingerprint existed
```

Persisting the same canonical spec MAY return an existing immutable
specification identity, because immutable content is deduplicated. It MUST NOT
overwrite the stored specification. Persisting is not a run and performs no
retrieval.

### 9.2 Start creates a distinct append-only run

```js
const run = lab.start(experiment, {
  label: "baseline rerun after filter fix",
  requestedBy: "manuel",
});

console.log(run.id);
console.log(run.status); // "queued" or "running"
```

Every successful `start` creates a new run ID even when the experiment
fingerprint is identical. It attaches a specification identity, records run
events and query traces, and eventually writes one terminal summary. It MUST
NOT edit prior runs, trace rows, or summaries.

`start()` is synchronous only through durable submission. It does not block
until all queries finish. `run.status()` and `run.awaitTerminal()` are deferred
until the laboratory's workflow integration is available; they are not part of
the v1 initial module contract. Operators use the existing run-inspection API
or web UI for progress.

### 9.3 Inspection and comparison

### 9.3 Execute runs a frozen evaluation manifest synchronously

```js
const lab = rag.open({
  database: "data/rag-eval.db",
  execution: "allowRuns",
  queryEmbed: query => embedder.embed(query),
});
const result = lab.execute(experiment);
```

`execute()` loads the immutable evaluation manifest named by the experiment,
creates one append-only run, and records its query traces and terminal summary.
It rejects vector plans that lack `queryEmbed`; it never infers a provider from
the selected embedding-set artifact.

```js
const embedding = lab.inspect(rag.artifact("embeddingSet", "sha256:..."));
console.log(embedding.compatibility);

const comparison = lab.compare("run_abc", "run_def");
console.log(comparison.metricDeltas);
```

`inspect` is read-only and returns artifact metadata, dimensions, source
identity, bytes, and compatibility parents. `compare` wraps the existing
immutable run comparison service; it does not recompute a run.

## 10. Study builder

A study groups comparable experiments and emits a serializable study plan. It
does not change how each experiment is fingerprinted or executed.

```ts
interface Study {
  readonly name: string;
  use(fragment: Fragment): this;
  case(name: string, configure: (experiment: ExperimentBuilder) => void): this;
  toPlan(): StudyPlan;
}
```

```js
const study = rag.study("ttc-representation-comparison", (s) =>
  s
    .use(ttcRawBaseline)
    .case("raw-only", (e) =>
      e.retrieval((r) => r.channel("raw", (c) => c.vector().topK(50)).results(10)),
    )
    .case("raw-plus-summaries", (e) =>
      e
        .representations((r) =>
          r.summaries("summary", (x) => x.artifact("sha256:summary-set").parent("sourceChunk")),
        )
        .retrieval((r) =>
          r
            .channel("raw", (c) => c.vector().representation("raw").topK(50))
            .channel("summary", (c) => c.vector().representation("summary").topK(50))
            .fuse((f) => f.rrf())
            .collapse("parentChunk")
            .results(10),
        ),
    ),
);

const plan = study.toPlan();
```

`case` begins with the study fragments but otherwise builds a normal
experiment. Case names are unique. V1 does not provide `study.startAll()`;
starting many compute-bearing runs deserves an explicit operator command with
budget and concurrency fields rather than a short hidden loop in JavaScript.

## 11. TypeScript declaration sketch

The provider MUST ship declarations generated from the same module contract.
This sketch communicates ergonomics; Go's `tsgen/spec` representation is the
implementation source.

```ts
declare module "rag" {
  export function open(options: OpenOptions): Laboratory;
  export function fragment(name: string, configure: Configure<ExperimentBuilder>): Fragment;
  export function experiment(name: string, configure?: Configure<ExperimentBuilder>): Experiment;
  export function study(name: string, configure?: Configure<StudyBuilder>): Study;
  export function artifact<K extends ArtifactKind>(kind: K, id: string): ArtifactRef<K>;
  export function grade(name: RelevanceGradeName): RelevanceGrade;
  export const version: "v1";

  type Configure<T> = (builder: T) => void;
  type RelevanceGradeName = "0_FAIL" | "1_PARTIAL" | "2_SUBSTANTIAL" | "3_AUTHORITATIVE";
  // remaining interfaces follow the API sections above
}
```

Two type systems protect the contract. The generated declaration gives script
authors completion and early diagnostics. The Go builder validates every call
at runtime; JavaScript remains dynamic, and scripts can be run without TypeScript.

## 12. Conventions and non-goals

- JavaScript property names and option keys MUST be lower camel case.
- IDs, hashes, channels, fragment names, and case names MUST be non-empty.
- Configurator lambdas MUST execute synchronously and MUST NOT be persisted.
- Builders MUST not accept raw SQL, unbounded JavaScript predicates, Go
  objects, functions as data, or provider credentials.
- `.toSpec()` and `.validate()` MUST be free of writes, provider calls, and
  embedding work.
- `lab.persist()`, `lab.start()`, and `lab.execute()` are explicit public effects.
- A module configuration MUST make side-effect authority opt-in.
- Any source/representation result returned after collapse MUST include
  original-source citation fields and representation provenance.
- V1 intentionally omits custom chunker lambdas, custom scorers, rerankers,
  prompt construction, and answer generation. Those need their own
  immutable-artifact and evaluation contracts before they become public DSL
  operations.

## 13. Related implementation references

- `cmd/rag-eval/xgoja.yaml` — current generated runtime module selection.
- `cmd/rag-eval/jsverbs/database.js` and `explorer.js` — current generic JS
  exploration primitives that the new module complements rather than replaces.
- `internal/services/experimentrun/service.go` — immutable specification and
  append-only run persistence boundary.
- `internal/services/immutableretrieval/bm25.go` and `vector.go` — baseline
  lexical/vector retrieval and deterministic RRF behavior.
- `docs/guides/ttc-rag-laboratory.md` — laboratory operator context.
- `RAGEVAL-TTC-LAB-001` design guide — corpus, artifact, evaluation, and UI
  contracts that this DSL must consume.
- `GOJA-DSL-PLAYBOOK` — builder/lambda/fragment and generated-declaration
  conventions used as API ergonomics evidence.
