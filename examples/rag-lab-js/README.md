# RAG laboratory JavaScript examples

These scripts use the `require("rag")` module in the generated
`rag-eval-js` binary. They intentionally use explicit immutable artifact IDs:
a script never selects "the latest" corpus, chunk set, embedding set, or
evaluation dataset.

Build the runtime from the repository root:

```bash
xgoja doctor -f cmd/rag-eval/xgoja.yaml
xgoja build -f cmd/rag-eval/xgoja.yaml --output cmd/rag-eval/dist/rag-eval-js
```

Run the pure, no-database authoring example:

```bash
cmd/rag-eval/dist/rag-eval-js run examples/rag-lab-js/01-plan-only.js
```

It prints a canonical specification and a stable fingerprint. It is safe to
run because `toSpec()` and `validate()` do not open or write a database.

To run `02-validate-persist-start.js`, replace every `REPLACE_WITH_...` ID
with artifacts from the same immutable lineage and set its `database` value.
The target database must already have the rag-eval migrations applied. The
script opens the database with `execution: "allowRuns"`, validates the
catalog, deduplicates the immutable specification, creates a distinct run,
and appends its `submitted` event. It does not yet execute retrieval; that is
the executor milestone.

For a read-only compatibility check, change `allowRuns` to `readOnly` and call
only `experiment.validate(lab)`. `persist` and `start` deliberately fail in
that mode.
