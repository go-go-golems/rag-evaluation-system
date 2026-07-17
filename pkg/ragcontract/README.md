# RAG researchctl wire contract

`pkg/ragcontract` contains dependency-free DTOs for the coordinated researchctl integration. It intentionally has no database, filesystem, Goja, provider, run, attempt, timestamp, or terminal-lifecycle API.

The authoritative schema/version is `rag-retrieval-spec/v1`. The matching researchctl implementation and JSON Schemas live in `researchctl/pkg/rag/spec`. Cross-repository golden tests prevent field/tag drift.

The existing `pkg/raglab.Laboratory` and `internal/services/experimentrun` paths remain prototype persistence APIs during migration. Native researchctl execution must use `raglab.ObservationExecutor`, whose observer emits domain events, query traces, metrics, and artifacts without creating or completing a run. Researchctl owns the enclosing run and attempt through its `ObservationSink`.

Unsupported execution behavior is an error. In particular, filters remain authorable but are rejected before any event or retrieval call until every channel can apply and trace them. Summary/question representations and parent-chunk collapse remain unsupported by the current executor.

## Observation contract

`ObservationExecutor` accepts resolved corpus, index, embedding, and evaluation inputs and reports only public DTOs:

- `Event` for domain progress without lifecycle authority.
- `QueryTrace` using payload schema `rag-query-trace/v1`.
- `Metric` with canonical JSON value and optional numeric projection.
- `Artifact` with immutable content identity.

The researchctl adapter routes trace payloads under core kind `rag-query-trace.v1`; schema version and routing kind are intentionally different strings. Cancellation retains accepted partial evidence but does not turn the attempt into success.

`cmd/rag-lab-worker` exposes these observations over `researchctl-rag-runner-stdio/v1`. It writes NDJSON protocol frames only to stdout and diagnostics to stderr. Provider secrets must enter through operator-controlled environment or secret stores and must not appear in specifications, traces, artifact metadata, or captured environment JSON.
