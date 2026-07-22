---
Title: RAG v2 execution hard-cut design and implementation guide
Ticket: RAG-V2-EXECUTION-CUTOVER
Status: active
Topics:
    - rag
    - rag-eval
    - workflow
    - scripting
    - intern-guide
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://cmd/rag-eval/cmds/study/command.go
      Note: Current domain-owned experiment loop
    - Path: repo://cmd/rag-ttc-v3-sweep/main.go
      Note: TTC-specific orchestration to delete
    - Path: repo://cmd/rag-worker/main.go
      Note: Current direct RAG runner candidate for cutover
    - Path: repo://internal/workflowv3ttc/module.go
      Note: TTC-specific Workflow V3 task package to replace
ExternalSources: []
Summary: Deletion-oriented guide for making the Researchctl-to-Workflow-V3 RAG path authoritative and removing competing runners and schemas.
LastUpdated: 2026-07-22T23:15:00-04:00
WhatFor: Prevent permanent dual execution paths after RAG workflow lowering reaches parity.
WhenToUse: Use only after deterministic and provider-backed Workflow V3 RAG slices pass their acceptance tests.
---


# RAG v2 execution hard-cut design and implementation guide

## Program context

This ticket belongs to **EXPERIMENT-PLATFORM-CONVERGENCE** and follows `RESEARCHCTL-EXPERIMENT-PLANS`, `SCRAPER-WORKFLOW-V3-PRODUCT-CUTOVER`, `EXPERIMENT-PLATFORM-SCRAPER-RUNNER`, `SCRAPER-WORKFLOW-OBSERVATIONS`, `RAG-V2-WORKFLOW-LOWERING`, `RAG-GEPPETTO-WORKFLOW-OPERATIONS`, and `RESEARCHCTL-EXPERIMENT-ANALYSIS`. It enables `TTC-SCRIPTED-EXPERIMENT-ACCEPTANCE`. All sibling guides live under the repositories' `ttmp/2026/07/22/` trees.

## Executive summary

RAG-eval currently contains a coherent greenfield v2 core plus multiple application, service, workflow, experiment, and TTC execution paths. Keeping the direct `ragengine` worker, RAG-owned experiment loops, old preparation workflows, and Workflow V3 lowering indefinitely would recreate the ambiguity this program is intended to remove.

This ticket performs a hard cut. The surviving path is: pure RAG JS authoring, canonical RAG v2 compilation, Researchctl experiment planning, Scraper Workflow V3 execution, RAG task packages, and Researchctl analysis. There are no compatibility shims.

## Surviving architecture

```text
require("rag") authoring
   -> pkg/ragmodel
   -> pkg/ragcompiler
   -> pkg/ragcontract
   -> pkg/ragworkflow lowering
   -> Scraper Workflow V3
   -> pkg/ragoperators task implementations
   -> RAG metrics/artifacts
   -> Researchctl laboratory and analysis
```

Likely retained packages:

- `pkg/ragcontract`
- `pkg/ragmodel`
- `pkg/ragcompiler`
- `pkg/gojamodules/rag`
- `pkg/ragoperators`
- RAG provider abstractions/adapters
- new `pkg/ragworkflow`
- selected product/query APIs that consume canonical plans

## Candidate superseded paths

The final list must come from dependency analysis, but inspect:

- `cmd/rag-ttc-v3-sweep/`;
- `internal/workflowv3ttc/`;
- generic replicate loops in `cmd/rag-eval/cmds/study/`;
- direct execution in `cmd/rag-worker/` after the Scraper runner replaces it;
- `internal/preparationworkflow/` where superseded by RAG Workflow V3 tasks;
- older `internal/workflow/` paths;
- custody adapters needed only for post-hoc imports;
- duplicate artifact/result schemas;
- stale examples and experiments using deleted execution commands.

`pkg/ragengine` requires a deliberate decision. It may remain temporarily as a deterministic semantic harness, but it must not remain a second production scheduler. Prefer extracting reusable operator evaluation from it and deleting its orchestration role.

## Cutover method

```pseudo
inventory = goListAllPackagesAndCommands()
for candidate in inventory:
    classify(candidate):
        canonical-core
        canonical-product
        test-fixture-only
        superseded-delete
        unresolved

assert no unresolved entries
run parity suites on canonical path
change commands/docs to canonical path
remove superseded packages in one branch
run go list/go test/static import scans
add negative guards for forbidden package paths
```

Maintain a machine-readable deletion manifest in this ticket listing package path, reason, replacement, and proof. Do not keep an empty forwarding package.

## CLI after cutover

RAG-eval should own domain authoring and inspection:

```text
rag-eval pipeline validate <script>
rag-eval pipeline explain <script>
rag-eval pipeline compile <script>
rag-eval product compile <script>
rag-eval study compile <script>
```

Experiment execution belongs to Researchctl. Workflow execution belongs to Scraper. Convenience commands may invoke those tools, but must be thin and must not implement their lifecycle again.

## Study boundary change

RAG studies may continue describing domain variants, requested RAG measures, and applicable factors. Generic replicate count, ordering, execution concurrency, and resume move to Researchctl's `ExperimentPlan`.

```text
RAG study compiler -> list of canonical domain cases
Researchctl plan   -> factors/cases + replicates + schedule
```

If factor expansion is domain-specific—for example an override modifies a RAG operator—RAG performs the semantic expansion and returns cases. Researchctl decides how and how often to run them.

## Decisions

### Decision: one production execution backend

- **Context:** Dual backends produce different retry, durability, and timing semantics.
- **Decision:** Workflow V3 is the sole production executor for canonical RAG plans.
- **Consequences:** Direct engine orchestration is deleted or reduced to a pure test harness.
- **Status:** accepted.

### Decision: deletion in the same ticket

- **Decision:** A replacement is not complete while the superseded path remains available as an ordinary command.
- **Consequences:** Documentation, examples, CI, and generated binaries must switch atomically.
- **Status:** accepted.

## Implementation phases

1. Produce package/command/schema inventory and deletion manifest.
2. Establish semantic parity goldens for deterministic pipelines.
3. Establish provider-operation parity and improved failure evidence.
4. Move generic study lifecycle to Researchctl plans.
5. Switch default commands and examples.
6. Remove TTC-specific and direct-runner orchestration.
7. Remove superseded packages and schemas.
8. Add import-boundary and forbidden-command tests.
9. Rewrite onboarding around the canonical path.

## Validation

- `go list ./...` contains no superseded packages.
- repository-wide `rg` finds no deleted schema/version names outside migration history in ticket docs.
- deterministic RAG metric goldens pass.
- crash/restart, retry, cancellation, and artifact lineage tests pass through Workflow V3.
- JS authoring and TypeScript parity tests pass.
- a Researchctl fixture plan executes the canonical RAG worker path.
- no command performs a post-hoc Researchctl import for newly executed runs.

## Intern guidance

This is not the ticket for inventing new architecture. Every deletion requires a replacement proved by preceding tickets. Start with an import graph and command inventory. Delete leaf packages first only after their callers have moved. Avoid renaming obsolete packages into vague `legacy` directories; Git history is the archive.

## Completion criteria

There is one documented way to author and execute RAG studies. Generic experiment lifecycle is absent from RAG. Production RAG scheduling is absent from `ragengine`. TTC-specific execution packages are gone. New engineers cannot accidentally choose the old path.

## Technology primer: semantic engine versus production executor

`ragengine` currently provides both domain semantics and execution sequencing. Those responsibilities can be separated. Domain semantics include how an operator transforms inputs, validates results, and emits RAG artifacts or metrics. Production execution includes durable scheduling, retries, leases, isolation, and restart. Workflow V3 should own the latter, while reusable operator code preserves the former.

This distinction prevents the cutover from becoming a rewrite of RAG algorithms. The implementation should extract or call semantic functions from Workflow V3 tasks, compare results with the existing engine, and then remove the engine's production orchestration role.

## How to build the deletion inventory

Use `go list -deps`, repository search, command help, examples, and CI workflows. A package with no imports may still be referenced by scripts or documentation. For each candidate, record:

```yaml
path: internal/workflowv3ttc
classification: superseded-delete
replacement: pkg/ragworkflow plus generic Scraper runner
proof:
  - deterministic lowering parity test
  - real-provider operation test
remove_with: RAG-V2-EXECUTION-CUTOVER
```

Group candidates by behavior, not directory name. The repository also contains frontend/widget and product-server areas unrelated to this program; do not delete them merely because they are not on the experiment path.

## Current-to-target command trace

Current study execution:

```text
rag-eval study run
  -> load RAG JS
  -> expand cells
  -> loop cells and replicates
  -> shell out to researchctl execute-spec
  -> rag-worker
  -> ragengine
```

Target execution:

```text
researchctl experiment run-plan
  -> load experiment JS and RAG descriptors
  -> schedule cases/replicates
  -> scraper-workflow-runner
  -> Workflow V3
  -> RAG task catalog
```

After cutover, `rag-eval study compile` may still compile domain studies to canonical cases. It must not create runs or own replicate scheduling.

## Parity ledger

Create a table for every retained operator family:

| Capability | Direct fixture | Workflow fixture | Domain artifacts equal | Metrics equal | Failure semantics reviewed |
|---|---|---|---|---|---|
| raw chunking | yes | yes | pending | pending | n/a |
| embeddings | yes | yes | pending | pending | pending |
| hybrid retrieval | yes | yes | pending | pending | n/a |
| generation | yes | yes | pending | pending | pending |

A passing happy-path output is insufficient for provider operators. Compare cancellation, timeout, malformed response, retry, usage, and cache behavior.

## Hard-cut mechanics

Make the canonical command path the default before deletion and run all smoke tests. Then delete complete package groups rather than leaving forwarding wrappers. Remove command registration, flags, docs, examples, CI invocations, generated artifacts, and tests tied only to old behavior. Run `go mod tidy` only after imports stabilize so removed dependencies become visible.

Add a repository guard test or CI grep for forbidden imports and commands:

```bash
if rg 'internal/workflowv3ttc|cmd/rag-ttc-v3-sweep' --glob '*.go' --glob '!ttmp/**'; then
  echo 'superseded execution path remains' >&2
  exit 1
fi
```

Historical ticket documentation may retain names. Guards should target production and active example directories.

## Data migration stance

No compatibility is required, but existing research evidence remains valuable as immutable historical artifacts. Do not migrate old workflow databases into the new runtime merely to preserve operability. Keep published reports and run exports as historical evidence. New executions use new canonical schemas and commands.

## Worked review scenario

Suppose `cmd/rag-worker` is deleted, but `pkg/researchctladapter/run.go` still constructs its command line. Compilation may not catch this if the string remains. CLI integration tests and repository search should. Similarly, deleting `internal/workflowv3ttc` may leave TTC-specific schema names in active analysis code. The deletion manifest ensures each reference has a disposition.

## First-week route

The intern begins by producing, not deleting, the inventory. Next they select one provider-free pipeline and execute it through both paths. They record semantic differences. Only when predecessor tickets resolve those differences do they switch one active example to the new path. Deletion starts after a reviewer accepts the parity ledger.

## Common mistakes

- Deleting semantic helper code before Workflow tasks reuse or replace it.
- Keeping a direct runner “for debugging” that becomes an unofficial production path.
- Moving old code into a `legacy` directory instead of removing it.
- Treating historical evidence as executable state that must be migrated.
- Forgetting scripts and docs that continue teaching old commands.
- Claiming parity based only on success cases.

## File-level implementation checklist

For each command under `cmd/`, record whether it remains, changes to compile-only behavior, or is deleted. For each package under `internal/workflow*`, record its replacement. For `pkg/researchctladapter`, separate generic canonical specification construction from TTC/post-hoc import code. For examples and experiments, identify the new invocation and expected output.

Review build surfaces beyond Go imports:

- Makefile and CI command invocations;
- embedded help pages;
- shell scripts and devctl pipelines;
- generated xgoja provider configuration;
- README examples;
- tests that launch binaries by string;
- container and release packaging.

After deletion, run help snapshots and smoke commands from a clean checkout. A repository can compile while still teaching a command that no longer exists.

## Review exercise

Select `cmd/rag-ttc-v3-sweep` and trace every package, schema, output file, and external command it uses. Assign each behavior to Researchctl, Scraper, RAG, Geppetto, TTC data, or deletion. Any behavior with two proposed owners is an unresolved architecture issue.

## Intern onboarding checklist

The engineer should draw both command traces, classify every candidate package, run one dual-backend parity fixture, identify non-Go references to old commands, explain what historical artifacts remain, and show the CI guard that prevents reintroduction.

## References

- Program: Researchctl `EXPERIMENT-PLATFORM-CONVERGENCE`.
- Canonical core: `pkg/ragcontract/`, `pkg/ragmodel/`, `pkg/ragcompiler/`, `pkg/ragoperators/`, `pkg/gojamodules/rag/`.
- Candidate removal paths listed above.
- Depends on `RAG-V2-WORKFLOW-LOWERING`, `RAG-GEPPETTO-WORKFLOW-OPERATIONS`, and `RESEARCHCTL-EXPERIMENT-ANALYSIS`.
