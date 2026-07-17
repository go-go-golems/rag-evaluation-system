---
Title: "Author RAG specifications from JavaScript"
Slug: "rag-laboratory-javascript"
Short: "Export pure versioned retrieval specifications for researchctl execution."
Topics:
- rag
- javascript
- evaluation
- retrieval
Commands:
- rag-eval
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

The generated JavaScript runtime exposes a typed fluent authoring API through `require("rag")`. The module builds and structurally validates retrieval semantics, then exports `rag-retrieval-spec/v1` data. It has no database, provider, run, attempt, persistence, or terminal-lifecycle API.

Researchctl is the only supported native execution authority. It resolves immutable inputs, validates lineage, creates runs and attempts, supervises the worker, records observations, verifies artifacts, and owns export/import.

## Build the JavaScript runtime

Generate and inspect the runtime from the repository root:

```bash
xgoja doctor -f cmd/rag-eval/xgoja.yaml
xgoja build -f cmd/rag-eval/xgoja.yaml --output cmd/rag-eval/dist/rag-eval-js
cmd/rag-eval/dist/rag-eval-js help --all
```

The generated declaration file is `js/types/xgoja-modules.d.ts`.

## Export a pure specification

Start with `examples/rag-lab-js/01-plan-only.js`. A plan names retrieval semantics; researchctl input references carry corpus, chunk, index, embedding, and evaluation identities separately.

```javascript
const rag = require("rag");

const experiment = rag.experiment("ttc-vector", (e) => e
  .corpus("authoring-corpus-reference")
  .chunks("authoring-chunk-reference")
  .embeddings("authoring-embedding-reference")
  .evaluation("authoring-dataset-reference")
  .representations((r) => r.rawChunks("raw"))
  .retrieval((r) => r
    .channel("semantic", (c) => c.vector().representation("raw").topK(50))
    .collapse("document")
    .results(10))
  .metrics((m) => m
    .relevanceAt(rag.grade("2_SUBSTANTIAL"))
    .recallAt([10])
    .mrr()));

const report = experiment.validate();
if (!report.ok) throw new Error(JSON.stringify(report));
module.exports = experiment.exportSpecification({ datasetSplit: "development" });
```

The export path is:

```text
JavaScript builder → structural validation → rag-retrieval-spec/v1
                                                │
                       immutable inputs + researchctl run-rag
                                                │
                     supervised observation-only RAG worker
```

The module intentionally does not expose `rag.open`, `toSpec`, `toJSON`, `persist`, `start`, `run`, or `execute`.

## Execute through researchctl

Build the strict NDJSON worker:

```bash
go build -o .bin/rag-lab-worker ./cmd/rag-lab-worker
```

Then invoke researchctl with the exported program and explicit immutable inputs:

```bash
researchctl experiment run-rag experiment.js \
  --project project.yaml \
  --experiment-id EXP-RAG \
  --inputs inputs.json \
  --ttc-database data/rag-eval.db \
  --runner .bin/rag-lab-worker \
  --runner-arg=--db --runner-arg=data/rag-eval.db \
  --runner-has-embedder --timeout 10m
```

Embedding and reranking endpoints belong to worker/operator configuration. Never place API keys, bearer tokens, or host credentials in the JavaScript specification, trace payload, artifact metadata, or captured environment evidence.

## Unsupported features

Filters, generated representations, and parent-chunk collapse remain authorable for contract development but fail capability validation before retrieval. Missing vector embedders and rerankers also fail. The worker never substitutes another method silently.

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| `rag.open` is undefined | Prototype lifecycle authority was removed. | Export a specification and execute it with researchctl. |
| Export rejects the metric plan | The requested metric is outside `rag-retrieval-spec/v1`. | Use supported ranking measures or version the contract before adding semantics. |
| Native execution cannot resolve an input | The catalog identity or lineage is missing/incompatible. | Correct the researchctl input reference; do not query mutable “latest” state. |
| Vector execution reports a missing embedder | No compatible query provider was declared. | Configure the worker with a manifest-compatible embedder. |

## See Also

- `examples/rag-lab-js/README.md` — copyable pure-authoring workflow.
- `pkg/gojamodules/rag/typescript.go` — generated public API shape.
- `pkg/ragcontract/README.md` — observation and worker boundary.
- Researchctl help topic `rag-laboratory` — execution, inspection, export, and import.
