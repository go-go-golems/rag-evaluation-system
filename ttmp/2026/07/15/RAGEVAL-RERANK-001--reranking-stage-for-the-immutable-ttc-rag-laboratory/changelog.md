# Changelog

## 2026-07-15

- Initial workspace created

- Completed current-state analysis and selected a llama.cpp-native cross-encoder
  boundary, preserving Ollama for embeddings and keeping all reranking
  provenance in immutable experiment/run records.

- Validated the ticket with docmgr and uploaded the combined index, design, and
  diary bundle to `/ai/2026/07/15/RAGEVAL-RERANK-001` on reMarkable.

- Completed the live llama.cpp BGE reranker probe through the private Mac
  tunnel; stored the executable probe and exact response/decoder contract.

- Added immutable cross-encoder reranking policy and pure transport-neutral Go
  contracts with fingerprint and validation regression coverage (commit 3764e20).

- Added and contract-tested the strict llama.cpp `/v1/rerank` adapter: explicit
  endpoint/model identity, request-size limits, context-aware transport, and
  durable-ID hydration from validated result indexes.

- Exposed immutable reranking policy in `require("rag")` with the fluent
  `rerank(x => x.crossEncoder(...).candidates(...).results(...))` API,
  TypeScript declarations, and JavaScript-spec projection (commit 51c9f89).
