# RAG researchctl wire contract

`pkg/ragcontract` contains dependency-free DTOs for the coordinated researchctl integration. It intentionally has no database, filesystem, Goja, provider, run, attempt, timestamp, or terminal-lifecycle API.

The authoritative schema/version is `rag-retrieval-spec/v1`. The matching researchctl implementation and JSON Schemas live in `researchctl/pkg/rag/spec`. Cross-repository golden tests prevent field/tag drift.

The existing `pkg/raglab.Laboratory` and `internal/services/experimentrun` paths remain prototype persistence APIs during migration. Native researchctl execution must use `raglab.ObservationExecutor`, whose observer emits domain events, query traces, metrics, and artifacts without creating or completing a run. Researchctl owns the enclosing run and attempt through its `ObservationSink`.

Unsupported execution behavior is an error. In particular, filters remain authorable but are rejected before any event or retrieval call until every channel can apply and trace them. Summary/question representations and parent-chunk collapse remain unsupported by the current executor.
