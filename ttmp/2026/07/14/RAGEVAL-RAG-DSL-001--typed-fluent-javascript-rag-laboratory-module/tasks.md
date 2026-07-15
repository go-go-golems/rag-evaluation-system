# Tasks

## TODO

- [x] 1. Reconcile the canonical experiment-specification JSON schema with the API contract; add explicit schema-version handling if it is absent.
- [x] 2. Define pure Go `raglab` domain specs, artifact references, validation errors, and fluent builders without importing goja.
- [x] 3. Implement compatibility validation for snapshot, chunk, embedding, BM25, and evaluation-dataset references.
- [x] 4. Implement retrieval-channel, RRF, collapse, result-count, representation, and metric builders with deterministic canonical output.
- [x] 5. Add pure Go service tests for valid plans, invalid combinations, error paths, canonicalisation, and stable fingerprints.
- [x] 6. Implement the `require("rag")` NativeModule adapter with lower-camel JavaScript codecs and thrown validation errors.
- [x] 7. Add JavaScript runtime tests for lambdas, reusable fragments, `.toSpec()`, diagnostics, and explicit execution.
- [ ] 8. Implement an xgoja provider package, TypeScript declaration descriptor, and `cmd/rag-eval/xgoja.yaml` module selection.
- [ ] 9. Add copy/paste examples under `examples/rag-lab-js/` and an xgoja doctor/build/declaration smoke test.
- [x] 10. Connect `lab.persist()` and `lab.start()` to immutable experiment-specification and run services; do not bypass their append-only rules.
- [ ] 11. Implement the first execution adapter for lexical, vector, and RRF retrieval, trace persistence, and terminal summaries.
- [ ] 12. Add operator documentation, help pages, generated declarations, and a web-UI link from an inspected run to its exported spec.
- [ ] 13. Validate with Go unit/integration tests, xgoja doctor, declaration generation, binary build, example scripts, and the normal web build.
- [ ] 14. Review the public API after one TTC study and decide whether summary/question representation generation needs a second module or a v1 extension.
