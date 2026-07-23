---
Title: RAG evaluation legacy cleanup inventory and disposition plan
Ticket: RAG-EVAL-LEGACY-CLEANUP
Status: active
Topics:
    - rag
    - rag-eval
    - workflow
    - research
    - intern-guide
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://cmd/rag-ttc-v3-sweep/main.go
      Note: Immediately removable standalone TTC orchestration
    - Path: repo://cmd/rag-worker/main.go
      Note: Deferred direct RAG execution path and parity oracle
    - Path: repo://internal/db/migrations.go
      Note: Old disposable-database compatibility upgrade
    - Path: repo://internal/workflow/echo_runner.go
      Note: Test-only Phase 0 compatibility runner
    - Path: repo://internal/workflowv3ttc/module.go
      Note: Immediately removable TTC-specific Workflow V3 task package
    - Path: repo://pkg/researchctladapter/operation_custody.go
      Note: TTC-only post-hoc custody adapter
ExternalSources: []
Summary: Evidence-backed disposition of TTC-specific execution, old workflow intake, direct RAG execution, compatibility migrations, and canonical RAG v2 packages.
LastUpdated: 2026-07-22T23:45:00-04:00
WhatFor: Approve immediate RAG cleanup and gate deferred execution cutovers on Workflow V3 replacement work.
WhenToUse: Read before removing TTC code or beginning RAG-V2-WORKFLOW-LOWERING.
---


# RAG evaluation legacy cleanup inventory and disposition plan

## Executive summary

RAG-eval contains a strong canonical RAG v2 core alongside three execution generations: an old Scraper-engine intake workflow, direct `ragengine` execution through `rag-worker`, and a TTC-specific Workflow V3 sweep. The TTC path is isolated enough for immediate removal after preserving its regression evidence. The old intake workflow and direct RAG engine remain active product paths and must wait for named replacements.

The immediate cleanup tranche should remove the standalone `rag-ttc-v3-sweep` binary, `internal/workflowv3ttc`, and the TTC-only operation-custody adapter; remove a test-only Phase 0 `EchoRunner`; remove old-development-database compatibility migration code; and prevent historical ticket scripts from being compiled as one Go package. The last issue currently causes `GOWORK=off go test ./...` to fail because four ticket scripts each declare `main` in the same directory.

Deferred cleanup should preserve `pkg/ragengine` as a semantic oracle until RAG Workflow V3 parity passes, preserve `rag-worker` until the generic Scraper runner executes lowered RAG plans, preserve `internal/preparationworkflow` until RAG tasks replace it, and preserve the old intake workflow until product intake has an accepted V3 replacement.

This report stops at classification. No production or historical evidence files have been deleted.

## Program navigation

- Umbrella: `EXPERIMENT-PLATFORM-CONVERGENCE`.
- Cleanup siblings: `RESEARCHCTL-LEGACY-CLEANUP`, `SCRAPER-LEGACY-CLEANUP`.
- Replacements: `RAG-V2-WORKFLOW-LOWERING`, `RAG-GEPPETTO-WORKFLOW-OPERATIONS`, `RAG-V2-EXECUTION-CUTOVER`.
- Final acceptance: `TTC-SCRIPTED-EXPERIMENT-ACCEPTANCE`.

## Method and baseline

The audit inspected command registration, execution-package imports, legacy markers, TTC references, researchctl adapter symbols, database migrations, ticket scripts, and the full test suite.

```bash
cd /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system
GOWORK=off go test ./... -count=1
```

Most active packages passed, including `cmd/rag-ttc-v3-sweep`, `cmd/rag-worker`, `internal/workflowv3ttc`, `pkg/ragengine`, and `pkg/researchctladapter`. The overall command failed in:

```text
ttmp/2026/07/22/RAG-TTC-V3-SWEEP--.../scripts
```

with exact errors including:

```text
06-refresh-model-manifest-digests.go:18:6: main redeclared in this block
02-build-researchctl-custody-spec.go:11:6: other declaration of main
07-refresh-prompt-manifest-digests.go:18:6: main redeclared in this block
08-build-operation-custody-export.go:37:6: main redeclared in this block
```

Ticket scripts are evidence, not module packages. They need an exclusion convention.

## Canonical RAG v2 core to retain

Retain:

- `pkg/ragcontract`: sole data-only wire contract.
- `pkg/ragmodel`: Go-backed authoring model.
- `pkg/gojamodules/rag`: pure JavaScript authoring.
- `pkg/ragcompiler`: normalization and domain compilation.
- `pkg/ragoperators`: domain semantics, artifacts, and metrics.
- `pkg/ragproviders` and `pkg/ragproviders/geppetto`: provider boundary.
- `pkg/ragproduct`: product/query runtime until separately superseded.
- immutable corpus/chunk/embedding/search services that implement current domain artifacts.

Retain `pkg/ragengine` temporarily as a semantic parity oracle, not as the final production scheduler.

## Immediate removal tranche

### 1. TTC-specific sweep binary and Workflow V3 package

**Where.** `cmd/rag-ttc-v3-sweep` (about 1,223 Go lines) and `internal/workflowv3ttc` (about 2,468 lines).

**Evidence.** Non-test repository references are confined to the standalone command and the package itself. The command is not registered under `cmd/rag-eval`. Its completed experiment, operation ledgers, corrected analysis, and design lessons are preserved in docmgr tickets and run artifacts.

**Remove.** Delete both code trees and command-specific tests. Preserve canonical regression fixtures needed by the new workflow observation and TTC acceptance tickets, preferably as small `testdata` records rather than the entire runner.

**Replacement.** `RAG-V2-WORKFLOW-LOWERING` supplies reusable RAG tasks; Researchctl supplies matrix/replicates; `TTC-SCRIPTED-EXPERIMENT-ACCEPTANCE` supplies the future workload script.

### 2. TTC-only operation custody adapter

`pkg/researchctladapter/operation_custody.go` builds post-hoc Researchctl run exports exclusively for the old TTC sweep. Its only production caller is `cmd/rag-ttc-v3-sweep/main.go`.

Delete the file and test with the sweep. Keep generic specification wrapping, worker capability checks, progress events, and run execution adapters until their deferred replacement.

### 3. Phase 0 `EchoRunner`

`internal/workflow/echo_runner.go` declares itself a Phase 0 compatibility runner. Search found it used only by `echo_runner_test.go`; production `internal/workflow/engine.go:49-50` registers `IntakeRunner`, not `EchoRunner`.

Delete both echo files. They no longer characterize production behavior.

### 4. Old development-database compatibility migration

`internal/db/migrations.go:8-63` detects an old chunks schema, creates a `legacy` strategy, disables foreign keys, rebuilds the table, and preserves old development rows. Current RAG v2 documentation explicitly says RAG databases are disposable and should be recreated instead of preserving prototype lifecycle/state.

Delete `ensureChunksStrategyID`, `tableHasColumn` if no longer used, its invocation at `internal/db/db.go:71`, and migration compatibility tests. Fresh migration tests remain mandatory.

### 5. Historical ticket Go scripts in the module package graph

Do not delete historical evidence. Rename independent `.go` scripts to a non-build extension such as `.go.txt`, place each in its own package directory, or establish a repository-wide ticket-script exclusion convention. The preferred cleanup is `.go.txt` for archival one-off programs and executable ticket-local scripts only when they have independent module/package structure.

The deletion tranche must make `go test ./...` pass without weakening production tests.

## Remove later after named replacement

### 1. Direct `rag-worker` / `ragengine` production execution

`cmd/rag-worker/main.go` decodes Researchctl domain specifications and directly calls `ragengine`. Study and preview commands use `pkg/researchctladapter.ExecuteSpecification` to launch it.

**Replacement.** `EXPERIMENT-PLATFORM-SCRAPER-RUNNER` plus `RAG-V2-WORKFLOW-LOWERING`.

**Deletion gate.** Deterministic domain artifacts and metrics match across direct and Workflow V3 execution; provider failure semantics are better or equal.

**Later action.** Remove direct scheduling from `ragengine` or reduce it to a pure semantic test harness. Replace `rag-worker` with a RAG task package executed through the generic Scraper runner.

### 2. `internal/preparationworkflow`

This package uses the old Scraper `pkg/workflow` engine and is called by `rag-worker`, `rag-preparation-smoke`, and TTC materialization. Keep it until chunk/representation/embedding preparation lowers into Workflow V3 tasks and prepared-artifact reuse passes parity.

### 3. Old intake workflow

`internal/workflow` and `cmd/rag-eval/cmds/workflow` provide active product intake: submit, worker, status, operations, and API artifact access. `internal/workflow/engine.go` registers the real `IntakeRunner`. This is not safe to delete now.

**Replacement.** Scraper V3 product cutover plus RAG task packages for intake operations.

**Deletion gate.** Product/API requirements are explicitly reimplemented or rejected, frontend/API callers migrate, and no active command imports old Scraper engine packages.

### 4. Generic Researchctl adapter lifecycle

Keep `pkg/researchctladapter/adapter.go`, `run.go`, `progress.go`, and `ttc.go` for now. Study and preview commands use them. Later:

- Researchctl experiment plans replace RAG-owned replicate loops.
- Generic Scraper runner replaces direct `ExecuteSpecification` to `rag-worker`.
- TTC catalog input resolution may remain a valid domain input adapter until TTC data becomes ordinary immutable input artifacts.

### 5. RAG study execution loop

`cmd/rag-eval/cmds/study/command.go` expands cells and loops replicates. Retain compile/validate/explain behavior. Move generic execution lifecycle to Researchctl plans, then remove `run` orchestration or make it a thin compiler/delegator.

## Replacement matrix

| Legacy path | Remove | Replacement | Gate |
|---|---|---|---|
| TTC sweep command/package | now | scripted TTC + generic framework | preserve regression fixture |
| operation custody adapter | now | native runner custody | sweep removed |
| EchoRunner | now | none | production never uses it |
| old DB upgrade | now | recreate disposable DB | fresh migration tests pass |
| ticket script package collision | now | archival script convention | `go test ./...` passes |
| direct rag-worker execution | later | Scraper runner + RAG lowering | parity suite passes |
| preparationworkflow | later | RAG V3 tasks | preparation parity/reuse passes |
| intake workflow | later | V3 product/task packages | product/API cutover passes |
| study run loop | later | Researchctl experiment plans | plan scheduling/resume passes |

## Proposed immediate deletion validation

```bash
GOWORK=off go test ./... -count=1
rg -n 'workflowv3ttc|rag-ttc-v3-sweep|BuildOperationCustodyRunExport|EchoRunner' \
  --glob '*.go' --glob '!ttmp/**' .
rg -n 'ensureChunksStrategyID|tableHasColumn' internal/db
```

Expected result: no active source references to removed symbols; complete test suite passes, including ticket tree package discovery.

## Deferred cutover sequence

1. Remove the immediate tranche in a focused commit.
2. Build deterministic RAG Workflow V3 lowering.
3. Add Geppetto-backed durable operation tasks.
4. Run direct-versus-workflow semantic parity fixtures.
5. Switch Researchctl execution to the generic Scraper runner.
6. Move study scheduling to Researchctl plans.
7. Port product intake requirements.
8. Delete direct worker, old preparation/intake workflow, and obsolete adapter lifecycle.
9. Add cutover guards for removed commands/packages.
10. Implement TTC only as a script and input/report package.

## Risks

The TTC package contains useful failure and operation-ledger fixtures. Extract those before deletion. The old DB migration may protect developer databases that someone still values, but project documentation declares them disposable and compatibility is explicitly out of scope; communicate the hard cut. `ragengine` must not be deleted before it serves its parity role.

## Review checklist

- Accept the five immediate cleanup items.
- Confirm historical TTC evidence remains in tickets/artifacts.
- Confirm developer databases may be recreated.
- Accept Workflow V3 lowering and generic runner as gates for direct execution removal.
- Accept V3 product/task work as the gate for intake removal.
- Stop here before destructive edits, as requested.

## References

- `cmd/rag-ttc-v3-sweep/`
- `internal/workflowv3ttc/`
- `pkg/researchctladapter/operation_custody.go`
- `internal/workflow/echo_runner.go`
- `internal/db/migrations.go:8-86`
- `internal/db/db.go:71`
- `cmd/rag-worker/main.go`
- `internal/preparationworkflow/`
- `internal/workflow/`
- `cmd/rag-eval/cmds/study/command.go`
- `docs/guides/rag-v2-destructive-cutover.md`
- `RAG-V2-WORKFLOW-LOWERING`
- `RAG-GEPPETTO-WORKFLOW-OPERATIONS`
- `RAG-V2-EXECUTION-CUTOVER`
- `TTC-SCRIPTED-EXPERIMENT-ACCEPTANCE`
