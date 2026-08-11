# RAG v2 authoring examples

These scripts run only the pure `require("rag")` authoring compiler. They do not open files, call models, build indexes, access a database, or create laboratory runs.

- `01-product.js` compiles `rag-product-plan/v2`.
- `02-five-variant-study.js` compiles five variants crossed with chunk/unit collapse.
- `03-fragment.js` demonstrates immediate fragment configuration and hidden references.
- `04-explain.js` emits a pure explanation DTO.
- `05-preview.js` emits `rag-preview-request/v1`; a later CLI/worker executes it.

Use the xgoja-generated runner or an embedding that registers `pkg/gojamodules/rag`:

```bash
xgoja doctor -f examples/xgoja/rag-v2/xgoja.yaml
xgoja gen-dts -f examples/xgoja/rag-v2/xgoja.yaml --out /tmp/rag-v2.d.ts
```

The fixture digests are structural examples, not claims that files with those identities exist.
