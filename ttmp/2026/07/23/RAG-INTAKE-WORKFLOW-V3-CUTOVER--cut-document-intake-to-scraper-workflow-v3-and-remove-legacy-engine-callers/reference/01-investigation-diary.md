---
Title: Investigation diary
Ticket: RAG-INTAKE-WORKFLOW-V3-CUTOVER
Status: active
Topics:
    - rag
    - rag-eval
    - intake
    - scraper
    - workflow
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/scraper/pkg/workflowv3product/service.go
      Note: Canonical V3 product submission/read/cancel/worker APIs
    - Path: repo://cmd/rag-eval/cmds/intake/root.go
      Note: Current product-local intake CLI boundary
    - Path: repo://internal/api/workflow_artifact_handlers.go
      Note: Current old engine-view API and submission surface
    - Path: repo://pkg/ragintakeworkflow/application.go
      Note: Canonical Workflow V3 intake application and submission
    - Path: repo://pkg/ragintakeworkflow/runtime.go
      Note: RAG task execution and provider operation custody
ExternalSources: []
Summary: Chronological implementation and validation record for the RAG intake Workflow V3 hard cut and its path to Scraper legacy cleanup.
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---




# Investigation diary

## Goal

Record every design, implementation, failure, validation, and deletion step required to migrate RAG document intake to Workflow V3 and continue through `SCRAPER-LEGACY-CLEANUP`.

## Step 1: Create the missing ticket and establish the migration boundary

The new ticket fills the explicit gap identified after the RAG v2 study execution cutover. I first re-audited the live callers rather than treating every remaining old-Scraper import as stale. The surviving code implements active document preprocessing, chunking, enrichment, embeddings, BM25 indexing, command workers, API projections, and frontend views.

The implementation guide therefore defines a hard cut with behavior preservation, not immediate deletion. It assigns durable execution and read models to Workflow V3 while keeping RAG request schemas, task semantics, providers, and database services domain-owned.

### Prompt Context

**User prompt (verbatim):** "create the missing ticket and add a design / implementation guide (doesnt need to be super elaborate), then work your way all the way up to SCRAPER-LEGACY-CLEANUP, keeping a detailed diary as you work. You have permission to hit external sources when testing."

**Assistant interpretation:** Create and document the missing intake cutover, implement it fully, then complete every required dependency and cleanup through Scraper's legacy-engine deletion, using external services for authorized acceptance and maintaining detailed evidence.

**Inferred user intent:** Remove the final cross-repository blockers rather than merely inventorying them, while preserving product behavior and leaving a reviewable migration record.

### What I did
- Created `RAG-INTAKE-WORKFLOW-V3-CUTOVER` with five implementation/acceptance tasks.
- Added a concise but implementation-ready design guide covering current/target flows, schemas, task identities, plan construction, provider custody, APIs, deletion gates, sequencing, and validation.
- Re-read the completed RAG execution cutover evidence and current Scraper legacy/Workflow V3 task status.
- Inspected all remaining RAG intake command and `internal/workflow` files, old API routes, frontend consumers, and Scraper Workflow V3 product APIs.

### Why
- RAG-eval remains an external consumer of Scraper's old engine, so `SCRAPER-LEGACY-CLEANUP` cannot safely delete it.
- The intake path mutates domain databases and contacts providers; a blind mechanical conversion would risk duplicated writes, secrets in plans, or incorrect retry semantics.

### What worked
- Scraper already exposes `workflowv3product.Application` for submission, list/show/cancel, worker operation, canonical observations, restart-safe SQLite, and domain-owned task packages.
- Existing RAG services isolate chunking, preprocessing, enrichment, embedding, and BM25 semantics from scheduler persistence, making them reusable behind a V3 task module.

### What didn't work
- N/A. This first step was inventory and design; no runtime implementation was attempted.

### What I learned
- The earlier 15 downstream legacy files are now eight direct old-Scraper importers after the study cutover.
- Manual per-operation retry is an old-engine UI/API concept and should not be copied; V3 retry policy is immutable plan identity, and reruns should be explicit new runs.
- The broader Scraper legacy ticket still has three unfinished Workflow V3 slices, including stronger isolation and a real TTC workload, in addition to product/site/API migration.

### What was tricky to build
- The ticket must separate RAG intake completion from Scraper-wide deletion. Completing the RAG package removes one external blocker but does not justify deleting Scraper root commands, site infrastructure, JavaScript runtime, metrics, events, and UI before their own dispositions pass.

### What warrants a second pair of eyes
- Review explicit-node versus lazy-map plan construction for database-mutating document work.
- Review whether V3 intake APIs should expose raw snapshots, canonical observations, or a smaller stable RAG projection.
- Review provider-operation custody for the existing embedding service, which currently performs provider calls internally.

### What should be done in the future
- Characterize exact existing database/output behavior in focused tests.
- Implement `pkg/ragintakeworkflow` and migrate one deterministic vertical slice before deleting old code.

### Code review instructions
- Start with the design guide's current/target traces.
- Compare `internal/workflow/submit.go`, `internal/workflow/intake_runner.go`, and Scraper `pkg/workflowv3product/service.go`.
- Validate ticket structure with `docmgr doctor --ticket RAG-INTAKE-WORKFLOW-V3-CUTOVER --stale-after 30`.

### Technical details
- New package identity: `rag-intake-v1@1.0.0`.
- Proposed request schema: `rag-intake-request/v1`.
- Remaining direct old-Scraper importers: eight files in intake commands, API, and `internal/workflow`.

## Step 2: Implement and hard-cut the intake product path

### What I did
- Added `pkg/ragintakeworkflow` with closed `rag-intake-request/v1`, result, and summary contracts; deterministic explicit-node plan compilation; the `rag-intake-v1@1.0.0` package; and trusted `rag:intake` task execution.
- Reused domain services for preprocessing, chunking, enrichment, embeddings, BM25 indexes, and publication while assigning leases, attempts, retries, budgets, artifacts, cancellation, observations, and SQLite custody to Scraper Workflow V3.
- Migrated `rag-eval intake` submit, run-once/run-worker, status, operations, and cancel commands to the Workflow V3 application.
- Replaced `/api/v1/workflows` with `/api/v1/intake/runs`, including submit, list, show, observations, and cancel. Removed manual operation retry.
- Replaced the legacy-engine frontend model and panels with a bounded Workflow V3 intake run view.
- Deleted `internal/workflow` and all legacy workflow UI panels in the same commit.
- Strengthened non-fake embedding execution with exact host provider-authority binding, provider request budgets, task usage evidence, durable `provider.embed/v1` external-operation custody, retries, and artifact secrecy tests.

### Key decisions
- Plans carry provider type/model identity and an authority digest, never endpoint credentials or secret values. The worker supplies host provider configuration.
- Fake embeddings do not reserve provider budget or create external-operation records; real providers do both.
- Database writes remain idempotent effects behind deterministic document/chunk identities. Restart tests reopen both databases before dispatch.
- API submission defaults to `skip_embeddings=true` when no embedding provider type is supplied; this avoids an implicit, unauthoritative provider choice.
- Manual operation retry was deleted rather than emulated because Workflow V3 retry policy belongs to immutable node policy.

### Commands and evidence
- `GOWORK=off go test ./pkg/ragintakeworkflow ./internal/api ./cmd/rag-eval/cmds/intake -count=1`
- `GOWORK=off go test ./... -count=1 -p=1`
- `./.bin/golangci-lint run ./pkg/ragintakeworkflow/... ./internal/api/... ./cmd/rag-eval/...` -> `0 issues.`
- `pnpm --dir web typecheck`
- `pnpm --dir web build`
- Pre-commit ran repository Go tests, configured lint, Glazed vet, Biome format, and Biome lint.
- Implementation commit: `a6bf3a8` (`feat: cut RAG intake to Workflow V3`).

### Failures and resolutions
- Initial full intake test failed with `expected: "succeeded" actual: "failed"`. Attempt inspection exposed `BUDGET_USAGE_INVALID: task usage evidence was invalid`. The plan reserved a provider request even for fake embeddings, while the task reported no usage. I made provider budget declaration conditional on a real authority and added sorted `ctx.usage.report` evidence for non-fake providers.
- The first commit attempt failed Glazed vet with `use Glazed config/env middleware or an explicit command field instead of os.Getenv in CLI code`. I restored an explicit host-only `--api-key` command field instead of bypassing the project CLI contract.
- The frontend hook reported pre-existing/fixable Biome advice around template construction and type-only imports. These were warnings, not acceptance failures; formatting, lint, typecheck, and build remained successful.

### What warrants a second pair of eyes
- Review `pkg/ragintakeworkflow/runtime.go` around external-operation descriptors and usage evidence.
- Review `internal/api/workflow_artifact_handlers.go` for the intentionally bounded projection and generic error bodies.
- Review whether a later product change should source provider secrets from Glazed's environment/config middleware rather than the preserved explicit CLI field.

### Code review instructions
1. Read `pkg/ragintakeworkflow/types.go`, `plan.go`, `package.go`, then `runtime.go`.
2. Run `GOWORK=off go test ./pkg/ragintakeworkflow -count=1` and inspect restart, cancellation, provider retry, privacy, and deterministic-plan tests.
3. Compare deleted `internal/workflow` in commit parent `a6bf3a8^` with the new CLI/API surfaces.
4. Search for forbidden old imports/routes with the ticket deletion guard once added.

## Step 3: Exercise the built binary and repair run-once semantics

The first real CLI smoke found a defect that package tests did not expose. `Dispatcher.DispatchOnce` is intentionally a lease-only deterministic hook. The initial command encoded the returned lease directly and failed with `Error: json: unsupported type: func()` because a lease contains executable task data. Worse, closing the application immediately left the claimed attempt running until lease expiry.

I changed `rag-eval intake run-once` to execute the claimed lease synchronously through `Engine.ExecuteLease` and emit only bounded identifiers (`runId`, `nodeKey`, `attempt`, `dispatched`). A fresh built-binary smoke then ingested one source document, submitted one immutable run, executed chunk/BM25/publish in three separate process invocations, reached `succeeded`, produced the BM25 files, and reported zero provider operations.

I added ticket scripts:
- `scripts/01-smoke-intake-workflow-v3.sh`
- `scripts/02-cutover-guards.sh`

Additional script failures were evidence-preserving test corrections:
- `Error: unknown flag: --run-id` on `intake ops`; the canonical flag is `--workflow-id`.
- The first operation assertion looked at `.externalOperations`; the bounded ops response nests it under `.run.operations.externalOperations`.
- The first index check assumed `smoke-index/index.json`; the BM25 service uses its own physical file layout, so the acceptance now checks for a non-empty host index tree instead of coupling to an internal filename.
- The first route guard matched the unrelated documentation slug `rag-geppetto-workflow-operations`; it now rejects only concrete legacy API route prefixes.

Final script output:
```text
PASS: built-binary intake run completed through Workflow V3 with 3 succeeded nodes and no provider operations
PASS: no legacy RAG intake lifecycle, route, flag, retry, or UI path remains
```

## Step 4: Render the migrated frontend and complete focused acceptance

Per operator request, the local production server runs in tmux session `rag-intake-ui`, with live output available through `tmux capture-pane -pt rag-intake-ui:0.0 -S -100` and a tee log at `/tmp/rag-intake-ui/server.log`.

I inspected the screenshots directly with the image-capable `read` tool. The first image showed that global reset styles made both ID inputs visually disappear and that empty run/list states were ambiguous. I added explicit borders, padding, background, and font inheritance; disabled submission until at least one document/source ID is present; and added loading/error/empty copy.

A rendered API submission then exposed two runtime defects in the browser console:
```text
404 /api/v1/intake/runs/<id>/observations
TypeError: Cannot read properties of null (reading 'map')
```
The observations query was enabled before the run response existed, and a no-attempt snapshot encoded a nil Go slice as JSON null. I now query observations only for terminal runs and render `(attempts ?? [])`. A fresh production build in a new browser tab showed the empty and selected-running states with zero console errors/warnings. Durable screenshots are in `analysis/01-workflow-v3-intake-empty-state.png` and `analysis/02-workflow-v3-intake-running-state.png`.

Final focused validation passed full sequential Go tests, affected race tests, configured lint and Glazed vet, frontend typecheck/build, built-binary smoke, and deletion guards. `make logcopter-check` initially failed with:
```text
logcopter-gen: generated file is not current: .../pkg/ragintakeworkflow/logcopter.go
```
`GOWORK=off go generate ./...` produced the required generated area catalog; the subsequent check passed. `GOWORK=off go mod tidy` produced no module diff.

## Final dependency-chain closure

Researchctl legacy cleanup and Scraper Workflow V3, external-operation, and legacy-cleanup tickets are now complete. Scraper commit `b3df00e` deleted the old engine/workflow/site/API/frontend stack after its guard found zero RAG downstream imports. I reran the built RAG intake binary smoke and deletion guard after that hard cut: the three-node chunk/BM25/publish run succeeded with no fake-provider operation, and no legacy lifecycle/route/flag/manual retry/UI path remained. A fresh full sequential RAG suite also passed, including `pkg/ragintakeworkflow` and `pkg/ragworkflow`. This satisfies the final prerequisite task and proves intake no longer blocks Scraper deletion.
