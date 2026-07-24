---
Title: RAG v2 execution cutover acceptance audit
Ticket: RAG-V2-EXECUTION-CUTOVER
Status: active
Topics:
  - rag
  - workflow
  - research
  - evaluation
DocType: analysis
Intent: long-term
Summary: Requirement-to-evidence audit for the destructive cutover from direct RAG execution to Researchctl plans and Scraper Workflow V3.
---

# RAG v2 execution cutover acceptance audit

## Objective restatement

Completion requires all of the following to be true at once:

1. Pure RAG JavaScript remains the authoring surface and canonical RAG v2 compilation remains intact.
2. RAG-eval compiles semantic cells but does not own generic cases, process execution, ordering, concurrency, retry, resume, or scientific persistence.
3. Researchctl plans are the sole scientific experiment lifecycle.
4. Scraper Workflow V3 is the sole production executor for canonical RAG study plans.
5. Direct `ragengine` study execution, the direct worker, preview lifecycle, old preparation Workflow, TTC-specific runner, and post-hoc direct-runner custody are unavailable as ordinary production paths.
6. Provider-free and provider-backed semantics, artifacts, metrics, identities, failures, usage, privacy, cancellation, retry, restart, and replay remain proved before deletion.
7. Generated inputs are canonical, bounded, lineage-bound, immutable, provider-authority-bound where required, and executable across repositories.
8. Active CLI help, guides, examples, frontend text, package graph, and guards teach only the canonical path.
9. The separately gated document-intake product path remains supported but cannot be mistaken for RAG v2 study execution.
10. The detailed diary, deletion manifest, tests, smoke evidence, docs, tasks, changelog, and repository validation are complete.

## Prompt-to-artifact checklist

| Requirement | Concrete evidence | Coverage assessment |
|---|---|---|
| Canonical authoring/core retained | `pkg/ragmodel`, `pkg/gojamodules/rag`, `pkg/ragcompiler`, and `pkg/ragcontract` remain; full package tests pass | Covered |
| Compile-only RAG ownership | `cmd/rag-eval/cmds/study/command.go` registers only validate/explain/compile; compile creates `StudyWorkflowCase` values | Covered |
| Researchctl owns lifecycle | Generated plan carries cases/factors/replicates but no explicit ordering, concurrency, fail-fast, retry, or execution loop; commit `68b1084` | Covered |
| Workflow V3 sole study executor | Every generated domain config is `scraper-workflow-execution/v2`; built smoke advertises `scraper-workflow-runner/v1` | Covered |
| Direct worker removed | `cmd/rag-worker` absent; package graph and guard verify absence | Covered |
| Preview lifecycle removed | `cmd/rag-eval/cmds/preview`, command registration, help page, and preview example absent | Covered |
| RAG-owned run loop removed | `study run`, `researchctladapter.ExecuteSpecification`, worker capability/process runner, run/progress adapter absent | Covered |
| Old preparation Workflow removed | `internal/preparationworkflow` and `cmd/rag-preparation-smoke` absent | Covered |
| TTC-specific runner removed | Predecessor cleanup commit `8a33a61`; guards find no active `workflowv3ttc` or `rag-ttc-v3-sweep` | Covered |
| `ragengine` not a production study scheduler | No command imports it; guard allows only Workflow semantic use, provider option construction, and online product runtime | Covered |
| Domain input lineage | `researchctladapter.LoadDomainArtifacts`, strict staged digest checks, `ragoperators.ValidateInputArtifacts`, adapter tests | Covered |
| Canonical bundle contract | `StudyBundleSchema = rag-workflow-study-bundle/v1`; canonical JSON per case plus manifest | Covered |
| Output path boundary | `filepath.Abs`/`Rel` containment and `RAG_WORKFLOW_STUDY_OUTPUT_BOUNDARY` test | Covered |
| Immutable publication | exclusive-create `writeImmutableStudyFile`; same-byte idempotency and changed-byte conflict tests; commits `f72f164`, `e119086` | Covered |
| Bounded materialization | query archive enforces 1..10,000 ordered items and Workflow map policies retain bounded paging/materialization | Covered |
| Provider-free lowering | `NewLowerer`, `BuildRunnerExecution`, stable bundle tests, built smoke | Covered |
| Provider authority | `NewProviderLowerer`, `BuildProviderRunnerExecution`, provider bundle test checks package and authority digest | Covered |
| No provider secrets in plan | Provider test rejects `provider-config` path and fixture response content in plan; config remains runner argument | Covered |
| Semantic parity | Predecessor lowering and provider tickets compare Workflow output with direct `ragengine` fixtures; design parity ledger updated | Covered |
| Provider failures/retries/cost/privacy | `RAG-GEPPETTO-WORKFLOW-OPERATIONS` acceptance audit and commits `4b93559`, `d4a957c`, `f6e083d`, `51ff7f4` | Covered |
| Crash/restart/cancel/replay | Predecessor built smoke and canonical observation evidence; current compiled-plan smoke proves resume | Covered |
| Cross-repository execution | `scripts/01-smoke-compiled-study.sh` builds RAG-eval, RAG runner, and Researchctl; validates then runs generated plan | Covered |
| Replicate ownership | Smoke compiles one RAG case requesting two replicates; Researchctl reports `executed=2`, then `resumed=2` | Covered |
| Verified outputs and metrics | Smoke checks succeeded attempts, verified artifacts, `rag.mrr`, and zero provider operations | Covered |
| Deletion inventory | `analysis/01-execution-path-deletion-manifest.json`: 14 unique, classified entries, no unresolved classification | Covered |
| Negative guards | `scripts/02-cutover-guards.sh`: packages, imports, command help, active docs, manifest, and `go list` | Covered |
| CLI clarity | root has `intake`, not ambiguous `workflow` or `preview`; study has compile, not run | Covered |
| Deferred intake behavior | `cmd/rag-eval/cmds/intake` and `internal/workflow` retained; help identifies product-local intake and separate V3 study route | Covered; separate product/API deletion gate remains intentionally outside this ticket |
| Active documentation | README, embedded help, guides, TTC guide, rag-sol2 README, and Evaluation page updated | Covered |
| Frontend | `pnpm --dir web typecheck` and build passed; embedded `internal/web/dist/index.html` regenerated | Covered |
| Detailed diary | Steps 1–7 record prompts, decisions, exact failures, commits, validation, review instructions, and follow-ups | Covered |

## Deletion evidence

The active package graph contains none of:

```text
cmd/rag-worker
cmd/rag-preparation-smoke
cmd/rag-eval/cmds/preview
cmd/rag-eval/cmds/workflow
internal/preparationworkflow
pkg/researchctladapter/run.go
pkg/researchctladapter/progress.go
```

The former top-level `workflow` name was removed rather than retained as an alias. Product-local document intake now uses `rag-eval intake`; no compatibility alias exists.

The eight remaining direct imports of Scraper's older engine/workflow packages are confined to the separately gated document-intake implementation and API:

```text
cmd/rag-eval/cmds/intake/ops.go
cmd/rag-eval/cmds/intake/status.go
internal/api/workflow_artifact_handlers.go
internal/workflow/engine.go
internal/workflow/intake_runner.go
internal/workflow/intake_runner_test.go
internal/workflow/submit.go
internal/workflow/submit_test.go
```

They do not participate in RAG v2 study compilation or execution. Deleting them here would violate the explicit frontend/API replacement gate in `RAG-EVAL-LEGACY-CLEANUP`.

## Fresh acceptance evidence

### Compiled-plan smoke

`analysis/02-compiled-study-smoke.txt` records:

```json
{"bundleSchema":"rag-workflow-study-bundle/v1","executed":2,"externalOperations":0,"metric":"rag.mrr","resumed":2,"status":"succeeded","validatedCases":1}
```

The script additionally checks every bundle file digest and size, Researchctl plan validation, verified exported artifacts, exact query scope, terminal success, and same-plan resume.

### Validation already completed during implementation

- Sequential `GOWORK=off go test ./... -count=1 -p=1`: passed.
- Full configured golangci-lint and Glazed vet path: passed.
- Pre-commit lint and package/internal tests: passed on commits `8925004`, `68b1084`, `f72f164`, and `e119086`.
- `pnpm --dir web typecheck`: passed.
- `pnpm --dir web build`: passed with only the existing bundle-size warning.
- `GOWORK=off go generate ./...`: passed.
- Built RAG-eval, `rag-workflow-runner`, and Researchctl: passed in smoke.
- Researchctl `validate-plan`, `run-plan`, and resume: passed.
- Cutover deletion/import/help/docs guards: passed.

## Final completion validation

Fresh validation after the final exclusive-create change:

- `GOWORK=off go test ./... -count=1 -p=1`: passed.
- `GOWORK=off go test -race ./pkg/ragworkflow ./pkg/researchctladapter ./cmd/rag-eval/cmds/study ./internal/workflow ./internal/api -count=1 -p=1`: passed.
- `GOWORK=off go mod tidy`: no `go.mod` or `go.sum` change.
- `GOWORK=off go generate ./...`: passed with no generated drift.
- `GOWORK=off go build ./...`: passed.
- Full configured golangci-lint and Glazed vet invocation: passed with zero issues.
- `pnpm --dir web typecheck` and `pnpm --dir web build`: passed; only the existing Vite chunk-size warning remained.
- `scripts/01-smoke-compiled-study.sh`: passed with two executions and two resumed replicates.
- `scripts/02-cutover-guards.sh`: passed with 14 classified manifest entries.
- Researchctl experiment plan/JS/service/lab/process-runner tests: passed.
- Scraper research runner, Workflow V3, runtime, and SQLite tests: passed.
- `docmgr doctor --ticket RAG-V2-EXECUTION-CUTOVER --stale-after 30`: passed after ticket closure.

No objective requirement remains intentionally deferred. The retained document-intake implementation is not a RAG v2 study executor and is explicitly governed by a separate accepted product/API deletion gate.
