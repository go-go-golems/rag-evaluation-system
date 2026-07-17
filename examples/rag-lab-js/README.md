# Pure RAG JavaScript authoring examples

The `rag` module is a lifecycle-free frontend for `rag-retrieval-spec/v1`. It builds and validates retrieval semantics, then exports plain data for researchctl. It does not open SQLite, contact providers, create runs, persist traces, or choose terminal state.

## Examples

- `01-plan-only.js` demonstrates fragments, raw BM25/vector channels, weighted RRF, metrics, candidate status, structural validation, and pure export.
- `04-export-researchctl-spec.js` is the cross-repository golden-style export used by the integration contract.

The older database-backed `rag.open`, `persist`, `start`, and direct execution examples were removed. Researchctl is the only supported native lifecycle path.

## Run an export

Build the xgoja runtime and execute the program using its JavaScript command, or load it through researchctl's RAG specification loader. The observable output is a `rag-retrieval-spec/v1` value; input artifact IDs are supplied separately to researchctl.

```bash
researchctl experiment run-rag examples/rag-lab-js/01-plan-only.js \
  --project /path/to/project.yaml \
  --experiment-id EXP-RAG \
  --inputs /path/to/inputs.json \
  --ttc-database /path/to/rag-eval.db \
  --runner /path/to/rag-lab-worker \
  --runner-arg=--db --runner-arg=/path/to/rag-eval.db
```

Vector and reranking capabilities are worker configuration. Keep credentials out of authored specifications, manifests, traces, and recorded environment evidence.
