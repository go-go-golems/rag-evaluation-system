---
Title: Investigation diary
Ticket: RAG-WORKFLOW-CONCURRENCY-HARD-CUT
Status: active
Topics:
    - rag
    - workflow
    - evaluation
    - embeddings
    - go
    - security
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://pkg/ragoperators/represent.go
      Note: Additional representation worker pool and progress machinery inventory
    - Path: repo://pkg/ragproviders/provider_set.go
      Note: GenerationConcurrency propagation and provider wrapper construction evidence
    - Path: repo://pkg/ragworkflow/lower.go
      Note: Current Workflow preparation granularity evidence
    - Path: repo://pkg/ragworkflow/provider_package.go
      Note: Workflow provider-service concurrency propagation evidence
ExternalSources: []
Summary: Chronological record for removing operator-local generation concurrency and cutting Workflow preparation over to deterministic batch maps.
LastUpdated: 2026-07-24T18:58:00-04:00
WhatFor: Preserve the cleanup rationale, evidence, commands, risks, and review instructions.
WhenToUse: Read before implementing or reviewing the concurrency hard cut.
---


# Investigation diary

## Goal

Remove duplicate RAG preparation scheduling with the smallest coherent cutover: Workflow V3 owns normal generation and embedding concurrency, the direct engine is serial, and provider limits remain only as validated host-global safety ceilings.

## Step 1: Separate simplification from broader scaling features

The previous parallelization-study guide intentionally covered a broad measurement program and several possible runtime improvements. Review showed that only a smaller subset genuinely simplifies the current implementation. This step created a dedicated hard-cut ticket for that subset and placed the broader study on hold.

The cleanup centers on deleting `GenerationConcurrency`, representation worker pools, hidden defaults, duplicate progress machinery, and superseded coarse provider-task paths. It deliberately excludes index sharding, query decomposition, new capacity-profile schemas, and full scaling experiments.

### Prompt Context

**User prompt (verbatim):** "Create a new design ticket for the minimal cleanup scope, and mark th eprevious document as "on hold for now".

Update the tasks accordingly too.."

**Assistant interpretation:** Create a focused design ticket containing only the concurrency changes that reduce implementation complexity, pause the broader parallelization design, and align both tickets' task lists with that decision.

**Inferred user intent:** Prevent the scaling initiative from accumulating new controls before duplicate scheduling and configuration paths are removed.

### What I did

- Rechecked all `GenerationConcurrency` references and current provider limiter behavior.
- Identified the distinction between Workflow dispatcher capacity and provider host-global safety ceilings.
- Created `RAG-WORKFLOW-CONCURRENCY-HARD-CUT`.
- Added an implementation design with current state, target architecture, deletion plan, phase gates, invariants, tests, decisions, and file references.
- Added focused tasks for generation/embedding maps, serial parity, worker-pool deletion, option removal, provider-ceiling validation, coarse-path deletion, capacity-default consolidation, and acceptance.
- Added `on-hold` to the docmgr status vocabulary for explicit paused work.

### Why

- A separate ticket prevents measurement features from obscuring the smaller cleanup boundary.
- The cleanup must hard-cut old paths rather than introduce mode flags.
- Provider safety limits cannot be deleted until their cross-workflow host-global responsibility is replaced or proven unnecessary.

### What worked

- Existing `PlanCombinedPreparation`, `ExecuteCombinedPreparationBatch`, `PlanEmbeddingBatches`, and `ExecuteEmbeddingBatch` provide the shared batch boundary.
- Workflow V3 already provides bounded maps, resource capacity, leases, attempts, retry, cancellation, and external operations.
- A serial direct engine provides a simpler semantic parity oracle.

### What didn't work

- N/A. No implementation code was changed in this step.

### What I learned

- `GenerationConcurrency` is propagated through engine, operator, provider, Workflow, fixture, export, and test contracts.
- Real providers already have separate generator, embedder, and reranker semaphores in `pkg/ragproviders/limit.go`.
- Those semaphores may be global across workflows while dispatcher capacity may be local, so their safety role must be retained and narrowed rather than assumed redundant.

### What was tricky to build

The main design difficulty was defining a hard cut that removes normal duplicate scheduling without weakening the provider's host-wide safety boundary. The solution is to let Workflow own scheduling, validate its capacity against `MaxInFlight`, and retain the provider wrapper only as a ceiling that does not normally contend.

A second difficulty was avoiding feature creep. Durable generation and embedding maps are required to remove worker pools. Index maps, query decomposition, artifact-graph redesign, and a global resource governor are not required for that deletion and remain outside this ticket.

### What warrants a second pair of eyes

- Whether provider wrapper instances are in fact shared across all relevant workflows in each runner mode.
- The exact startup location for validating Workflow capacity against provider `MaxInFlight`.
- Whether every representation operator can use the same serial plan/execute/finalize structure.
- Which coarse Workflow task keys can be deleted atomically after map parity.

### What should be done in the future

- Implement the serial operator baseline first.
- Lower generation and embedding maps.
- Execute the hard cut and acceptance suite.
- Reassess the paused TTC scaling study after the simpler runtime is merged.

### Code review instructions

- Start with the design's Sections 2–9 for current ownership and deletions.
- Verify every `GenerationConcurrency` reference listed by `rg` is covered by a task.
- Compare `pkg/ragproviders/limit.go` with Workflow dispatcher scoping before changing provider ceilings.
- Reject compatibility modes or fallback worker pools.

### Technical details

Evidence commands:

```text
rg -n "GenerationConcurrency|generationConcurrency" pkg cmd internal experiments -S
rg -n "newLimitedGenerator|newLimitedEmbedder|newLimitedReranker|providerConcurrencyLimit" pkg/ragproviders -S
rg -n "PlanCombinedPreparation|ExecuteCombinedPreparationBatch|PlanEmbeddingBatches|ExecuteEmbeddingBatch" pkg -S
```
