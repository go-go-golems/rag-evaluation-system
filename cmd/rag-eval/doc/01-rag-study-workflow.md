---
Title: "Compile and run RAG v2 studies"
Slug: "rag-study-workflow"
Short: "Validate, explain, compile, and execute pure RAG studies through researchctl."
Topics:
- rag
- studies
- evaluation
- researchctl
Commands:
- rag-eval study validate
- rag-eval study explain
- rag-eval study compile
- rag-eval study run
Flags:
- inputs
- artifact-root
- ttc-database
- experiment-id
- worker-command
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

RAG studies keep semantic authoring in JavaScript while Go owns compilation and execution. `rag-eval` loads a pure `rag-study/v2` export, resolves immutable domain inputs, expands stable cells, wraps each cell in researchctl's generic execution identity, and invokes the generic laboratory command. Researchctl never imports or decodes a RAG package.

## Author a study

Export the result of `study.compileStudy(...)`; do not call a lifecycle method from JavaScript. `examples/rag-v2/06-raw-study.js` is the smallest native BM25 example.

```javascript
const rag = require("rag");
// Build pipeline, query plan, and study with Go-backed builders.
module.exports = study.compileStudy({ inputs: placeholderBindings });
```

The `--inputs` document replaces placeholder bindings before cell expansion. A file reference points to a v2 manifest envelope. A TTC alias names a catalog object that the RAG-owned adapter materializes under researchctl artifact custody.

```json
{
  "inputs": {
    "corpus": {
      "catalog": {"namespace": "rag-eval-ttc", "id": "sha256:..."}
    },
    "evaluation-dataset": {
      "catalog": {"namespace": "rag-eval-ttc", "id": "candidate:..."}
    }
  }
}
```

## Validate and explain before allocating runs

Validation resolves manifests and deeply expands every cell. Explanation adds the ordered variants, factors, and registered operator IDs.

```bash
rag-eval study validate study.js --inputs inputs.json --ttc-database rag-eval.db
rag-eval study explain study.js --inputs inputs.json --ttc-database rag-eval.db
```

Failures here occur before researchctl allocates a run. Missing aliases, digest mismatches, malformed manifests, unsafe collapse, and unsupported operators never degrade silently.

## Compile generic specifications

Compilation writes one canonical researchctl specification per stable cell. The opaque `domainConfig` remains `rag-pipeline-execution/v2`; researchctl validates only its generic envelope.

```bash
rag-eval study compile study.js \
  --inputs inputs.json \
  --ttc-database rag-eval.db \
  --spec-output-dir ./compiled
```

## Execute cells and replicates

Initialize the laboratory, then let the RAG CLI invoke researchctl's generic external-runner command. Capability probing happens before run allocation, and the worker revalidates domain/version, canonical execution identity, manifests, and lineage.

```bash
researchctl lab init --project project.yaml
rag-eval study run study.js \
  --inputs inputs.json \
  --ttc-database rag-eval.db \
  --project project.yaml \
  --experiment-id EXP-RAG \
  --researchctl-command researchctl \
  --worker-command rag-worker
```

Researchctl owns run/attempt IDs, retry and timeout policy, timestamps, artifact custody, observation ordering, required-measure checks, and terminal summaries. The worker owns only RAG execution and observations.

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| `RAG_CATALOG_RESOLVER_REQUIRED` | An input uses a catalog alias without a TTC database. | Pass `--ttc-database` or replace the alias with a manifest-envelope URI. |
| `RAG_WORKER_CAPABILITY` | The executable did not advertise the exact generic protocol/runner identity. | Build the current `cmd/rag-worker` and check `--worker-command`. |
| `RAG_WORKER_INPUT_LINEAGE` | Corpus/evaluation manifests do not match the execution bindings or each other. | Regenerate envelopes from the same immutable corpus snapshot. |
| Researchctl reports missing required measures | The worker failed before evaluating all requested measures. | Inspect the preserved attempt events, artifacts, traces, and terminal payload. |

## See Also

- `rag-eval help rag-preview-workflow`
- `examples/rag-v2/06-raw-study.js`
- `examples/rag-v2/inputs-ttc-catalog.json`
