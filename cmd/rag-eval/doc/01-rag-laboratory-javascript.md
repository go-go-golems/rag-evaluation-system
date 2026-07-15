---
Title: Run immutable RAG laboratory experiments from JavaScript
Slug: rag-laboratory-javascript
Short: Build, validate, execute, and inspect reproducible retrieval experiments with require("rag") and Geppetto embeddings.
Topics:
  - rag
  - javascript
  - evaluation
  - retrieval
Commands:
  - rag-eval
Flags: []
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

# Run immutable RAG laboratory experiments from JavaScript

The generated `rag-eval-js` runtime exposes a typed fluent authoring API as
`require("rag")`. A script describes immutable input artifacts and retrieval
policy. It never selects a mutable “latest” artifact. Validation checks that
the selected corpus, chunks, BM25 index, embedding set, and evaluation dataset
belong to one compatible lineage before a run can be created.

The execution boundary is deliberately explicit. Vector channels require a
synchronous `queryEmbed(query)` callback. Geppetto provides that callback via
a resolved inference-profile registry and an embedding provider. Credentials,
server endpoints, and the selected model remain operational configuration; the
persisted experiment stores artifact identities and retrieval policy instead.

## Build the JavaScript runtime

From the repository root, generate and inspect the runtime before using it:

```bash
xgoja doctor -f cmd/rag-eval/xgoja.yaml
xgoja build -f cmd/rag-eval/xgoja.yaml --output cmd/rag-eval/dist/rag-eval-js
cmd/rag-eval/dist/rag-eval-js help --all
```

The generated declaration file is `js/types/xgoja-modules.d.ts`.
It covers the `rag` module. Whole-runtime strict declaration generation remains
tracked separately until the Geppetto provider publishes its descriptor.

## Author a plan without side effects

Start with `examples/rag-lab-js/01-plan-only.js`. It calls `toSpec()` and
`validate()` only, so it does not open or modify a database. A complete plan
names all immutable inputs and a retrieval pipeline:

```javascript
const experiment = rag.experiment("ttc-vector", (e) => e
  .corpus("corpus-id")
  .chunks("chunk-set-id")
  .embeddings("embedding-set-id")
  .evaluation("evaluation-dataset-id")
  .retrieval((r) => r
    .channel("semantic", (c) => c.vector().topK(50))
    .collapse("document")
    .results(10))
  .metrics((m) => m
    .relevanceAt(rag.grade("2_SUBSTANTIAL"))
    .recallAt([10])
    .mrr()));
```

The stages are:

```text
JavaScript builder
        |
        v
canonical immutable specification -- validate catalog lineage --> persist/start
        |                                                            |
        |                                                    query embedding callback
        v                                                            v
stable fingerprint                                         retrieval and metrics
```

## Execute vector or hybrid retrieval with Geppetto

Copy `examples/rag-lab-js/03-execute-with-geppetto.js`; replace the explicit
artifact IDs, database path, and profile registry path. The essential wiring is:

```javascript
const gp = require("geppetto");
const rag = require("rag");

const settings = gp.inferenceProfiles.load("profiles.yaml").resolve("embeddings");
const embedder = gp.embeddings(settings);
const lab = rag.open({
  database: "data/rag-eval.db",
  execution: "allowRuns",
  queryEmbed: (query) => embedder.embed(query),
});
```

`lab.execute(experiment)` validates, persists/reuses the canonical
specification, creates a new run, retrieves each evaluation query, hydrates
original-source citations, and persists trace data and metrics. It returns the
run identifier, query count, metrics, timing, and completion timestamp.

The profile registry is an external YAML or SQLite source understood by
Geppetto. A named profile contains an `embeddings` block whose endpoint and
model describe the currently available embedding service. Do not put API keys
or host credentials into a JavaScript experiment or commit them to this
repository.

## Inspect and compare results

Start the web application with `rag-eval serve` and open its Evaluation page.
Select a run, then a query trace. The trace inspector displays the immutable
specification identifier and links directly to the exported canonical JSON at
`/api/v1/lab/specifications/{id}`. Use that JSON to establish exactly which
inputs and policy produced a result before comparing retrieval quality.

## Safety checks

- Use `execution: "readOnly"` for catalog compatibility checks. `persist`,
  `start`, and `execute` intentionally fail in that mode.
- Use `execution: "allowRuns"` only with a database containing the current
  rag-eval migrations.
- Keep all artifact IDs explicit. If validation reports lineage errors, create
  or select a compatible artifact set rather than overriding the check.
- Treat an embedding model or endpoint change as an operational change. For a
  fair comparison, use embeddings that match the persisted embedding artifact.

## See also

- `examples/rag-lab-js/README.md` for build and execution commands.
- `examples/rag-lab-js/03-execute-with-geppetto.js` for a complete copyable
  vector/hybrid experiment.
- `pkg/gojamodules/rag/typescript.go` for the generated public RAG API shape.
- `pkg/raglab` for the authoritative Go-side specification, validation, and
  executor implementation.
