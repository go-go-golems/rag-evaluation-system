---
Title: RAG intake Workflow V3 cutover acceptance audit
Ticket: RAG-INTAKE-WORKFLOW-V3-CUTOVER
Status: active
Topics:
    - rag
    - intake
    - workflow
    - scraper
DocType: analysis
Intent: long-term
Summary: Requirement-to-evidence audit for the completed RAG intake implementation slice; broader prerequisite cleanup remains open.
---

# RAG intake Workflow V3 cutover acceptance audit

## Result

The RAG-eval document-intake product is hard-cut to Scraper Workflow V3. The intake implementation, CLI, API, and frontend no longer import, route to, or project the old Scraper engine. This audit closes the implementation and acceptance tasks but does **not** close the ticket's final dependency-chain task or the durable goal; Scraper-wide cleanup remains.

## Requirement mapping

| Requirement | Evidence |
|---|---|
| Closed versioned request/task contracts | `pkg/ragintakeworkflow/types.go`, `package.go`, `task.cjs`; unknown fields and malformed identities fail closed in tests. |
| Deterministic Workflow V3 lowering | `plan.go`; deterministic plan test; explicit document nodes and bounded IDs. |
| Workflow custody | `application.go`, Workflow V3 SQLite/product/runtime; restart test reopens stores before dispatch. |
| Domain behavior | `runtime.go` reuses RAG preprocessing, chunking, enrichment, embedding, BM25, and publication services. |
| Provider authority/privacy | Exact authority digest; host-only credentials; `provider.embed/v1` descriptors; retry/operation-count/privacy tests; provider budget and usage evidence. |
| CLI hard cut | `cmd/rag-eval/cmds/intake`; submit, run-once, run-worker, status, ops, cancel; built-binary smoke. |
| API hard cut | `/api/v1/intake/runs`; bounded body and projections; legacy routes and manual retry deleted. |
| Frontend hard cut | `WorkflowsView.tsx`; old operation panels deleted; empty/running screenshots in this directory; production console clean after rendered submission. |
| Cancellation/retry/restart/artifacts/observations | Focused package/API tests and race tests; smoke and provider-operation tests. |
| Legacy deletion | `internal/workflow` deleted; `scripts/02-cutover-guards.sh` rejects old imports, routes, flags, retry UI, and panels. |

## Fresh commands

All passed on 2026-07-23/24:

```text
GOWORK=off go test ./... -count=1 -p=1
GOWORK=off go test -race ./pkg/ragintakeworkflow ./internal/api -count=1 -p=1
make lint
make logcopter-check
GOWORK=off go mod tidy                         # no module diff
pnpm --dir web typecheck
pnpm --dir web build
scripts/01-smoke-intake-workflow-v3.sh
scripts/02-cutover-guards.sh
```

Built-binary smoke result:

```text
PASS: built-binary intake run completed through Workflow V3 with 3 succeeded nodes and no provider operations
PASS: no legacy RAG intake lifecycle, route, flag, retry, or UI path remains
```

## Rendered inspection

- `01-workflow-v3-intake-empty-state.png`: visible ID inputs, disabled invalid submit, explicit empty states.
- `02-workflow-v3-intake-running-state.png`: selected running run, bounded identity/status/plan projection, cancel action, attempt table.
- Browser console after the repaired build: zero errors and zero warnings.

The first rendered submission exposed a premature observations request (404) and a null-attempt `.map` crash. The frontend now queries observations only for terminal runs and treats absent attempts as an empty list.

## Remaining dependency work

The final ticket task remains open until prerequisite convergence/cleanup tickets are closed and `SCRAPER-LEGACY-CLEANUP` is genuinely unblocked. No evidence in this audit authorizes deletion of unrelated Scraper product/site/runtime clusters.
