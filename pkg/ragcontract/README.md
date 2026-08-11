# Canonical RAG v2 contracts

`pkg/ragcontract` is the sole canonical wire-contract package for the greenfield RAG pipeline. It owns:

- `rag-pipeline-ir/v2`;
- `rag-product-plan/v2`;
- `rag-study/v2`;
- `rag-pipeline-execution/v2`;
- `rag-query-trace/v2`;
- immutable corpus, unit, chunk, representation, embedding, index, evaluation, model, and prompt manifests.

The package is data-only and dependency-light. It imports no Goja, provider, retrieval engine, SQLite, filesystem runtime, or researchctl package. Strict decoders reject unknown fields and trailing values. Display metadata is structurally separate and excluded by compiler semantic-identity helpers.

Operator identifiers use `<namespace>.<operation>/<version>`, for example `fusion.weighted-rrf/v1`. An operator identifier names registered semantics; it is not a schema name or JavaScript factory.

Generated representations, collapse identities, and source-evidence identities are distinct. A trace may report a matched generated question, a winning unit, and a hydrated source chunk without treating generated text as evidence.

The disposable v1 laboratory and worker were removed after the one-time parity extraction. No compatibility DTO, adapter, runner, or schema remains in the active source tree.
