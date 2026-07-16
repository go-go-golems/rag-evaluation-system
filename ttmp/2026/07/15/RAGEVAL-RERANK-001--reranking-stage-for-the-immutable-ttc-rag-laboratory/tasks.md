# Tasks

## TODO

- [x] 1. Map the current raw retrieval executor, experiment persistence, and live TTC baseline; record the reranking insertion point.
- [x] 2. Select the initial runtime contract: llama.cpp native reranking endpoint behind a private Mac tunnel, not Ollama scoring emulation.
- [x] 3. Write the intern-ready design, API specification, evaluation matrix, and operational playbook; upload the package to reMarkable.
- [x] 4. Run a no-code `llama-server` probe with `qllama/bge-reranker-v2-m3:q4_k_m`; capture request/response shape, index semantics, latency, and model metadata in `scripts/`.
- [x] 5. Define Go domain types for `Reranker`, `RerankRequest`, `RerankResult`, and an immutable `RerankingSpec`; add contract tests without modifying the executor.
- [x] 6. Add the llama.cpp HTTP adapter with strict response validation, request size limits, context cancellation, and operator-configured endpoint/model identity.
- [x] 7. Extend experiment specifications and JavaScript builders with an optional `.rerank(...)` stage; no backwards-compatibility adapter is required because this API is unmerged.
- [x] 8. Apply reranking after candidate-channel fusion and before collapse/citation hydration; persist candidate and scored order in query traces.
- [x] 9. Add web trace visualization for pre-rerank rank, score, post-rerank rank, truncation, and reranker timing.
- [ ] 10. Run the TTC matrix: raw vector, weighted RRF, RRF+BGE rerank, and one Qwen reranker comparison; record quality, latency, local cost scope, and storage.
- [ ] 11. Decide whether reranking occurs before or after parent/document collapse using measured duplicate/citation behavior; document the decision.
- [ ] 12. Validate full Go/web/generated-JS paths and update operator playbooks and docs.
