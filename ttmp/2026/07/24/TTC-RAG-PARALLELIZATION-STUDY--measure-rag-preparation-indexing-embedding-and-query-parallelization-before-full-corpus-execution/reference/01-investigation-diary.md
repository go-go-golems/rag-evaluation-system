---
Title: Investigation diary
Ticket: TTC-RAG-PARALLELIZATION-STUDY
Status: on-hold
Topics:
    - ttc
    - rag
    - evaluation
    - research
    - workflow
    - embeddings
    - chunking
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/scraper/pkg/workflowv3runtime/dispatcher.go
      Note: Workflow capacity behavior reviewed during investigation
    - Path: repo://pkg/ragoperators/embedding_batch.go
      Note: Existing durable-ready embedding batch boundary
    - Path: repo://pkg/ragoperators/represent.go
      Note: Current in-process generation worker pool evidence
    - Path: repo://pkg/ragworkflow/lower.go
      Note: Primary evidence for current preparation granularity
ExternalSources: []
Summary: Chronological investigation record for the TTC RAG preparation, indexing, embedding, and query parallelization study.
LastUpdated: 2026-07-24T15:01:12.200075188-04:00
WhatFor: Preserve evidence, decisions, failures, commands, and review guidance for the medium-corpus scaling study.
WhenToUse: Read before resuming implementation, changing concurrency contracts, selecting study cells, or admitting provider budget.
---



# Investigation diary

## Goal

Design and execute a bounded scaling study that measures preparation, representation generation, embedding, indexing, and query parallelization on a medium immutable TTC dataset. The study must produce safe product and full-corpus capacity recommendations without compromising semantic output, provider accounting, custody, privacy, resume, or statistical validity.

## Step 1: Establish the scaling-study architecture and ticket

This step created a dedicated ticket instead of mixing performance work into the broader real-provider retrieval-quality study. It mapped current concurrency behavior across RAG-eval, Workflow V3, and Researchctl, then defined a staged experimental design and an implementation path for durable preparation maps.

The central finding is that current query work is durably parallelized through a Workflow V3 map, while static preparation is lowered as a sequential chain of coarse tasks. Representation generation has an in-process worker pool, and deterministic generation/embedding batch plan/execute APIs already exist, but embedding and index batches are not independently leased Workflow work. Measuring preparation capacity without changing this graph would produce misleading results because capacity above one has no additional ready static nodes to schedule.

### Prompt Context

**User prompt (verbatim):** "I want to do a first suite of experiments regarding the parallelization speed of preparing, indexing, embedding and querying a dataset, with a reasonably small size of data (but big enough to get insights), before we start the full corpus. We already did quite a bit of work on improving parallelization in the past, which led to a massive refactor across all 3 repositories.

Create a new ticket for that and Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Create a new docmgr ticket and exhaustive intern guide for implementing and running a medium-corpus TTC RAG concurrency/scaling study across preparation, generation, embedding, indexing, and querying, then validate and upload it to reMarkable.

**Inferred user intent:** Convert the previous cross-repository parallelization refactor into measured evidence and safe product/full-corpus capacity defaults before incurring the cost and operational risk of complete TTC processing.

### What I did

- Inspected RAG real-provider source, Workflow lowering, task bundles, operator worker pools, deterministic batch APIs, engine options, and canonical contracts.
- Inspected Workflow V3 dispatcher, map scheduling, resource capacities, SQLite leasing, reduction, and tests.
- Inspected Researchctl plan concurrency, replicates, ordering, durable plan provenance, and run-level analysis.
- Created ticket `TTC-RAG-PARALLELIZATION-STUDY`.
- Added twelve tasks covering frozen inputs, capacity contracts, durable maps, telemetry, readable source, fixture acceptance, real stage studies, analysis, product defaults, cross-repository guards, and closure.
- Wrote the architecture and implementation guide with diagrams, pseudocode, APIs, file references, study factors, phased matrix, measurements, decisions, risks, and acceptance criteria.

### Why

- Full-corpus execution is too expensive and slow for discovering basic concurrency saturation.
- A smaller but sufficiently batched corpus can reveal provider, CPU, disk, SQLite, index, and query bottlenecks.
- Parallelism exists at several layers and must be isolated experimentally rather than represented by one generic concurrency factor.
- Current static Workflow granularity is insufficient for durable preparation scaling claims.

### What worked

- The previous refactor already provides deterministic batch planning and one-batch execution APIs for combined representation generation and embeddings.
- Workflow V3 already has bounded maps, resource-class capacity, immediate slot refill, durable retries, operations, reductions, and cancellation.
- Researchctl v0.0.3 now provides strict plans, durable plan provenance, run-level samples, deterministic ordering, resume, and analysis.
- Existing query lowering provides a concrete pattern for preparation batch maps.

### What didn't work

- The first broad RAG `rg` search produced no terminal output because the combined output handling was dominated by other parallel commands. I followed with focused file reads and narrower searches.
- No provider experiment was attempted. The current deliverable is the implementation-ready guide and ticket.

### What I learned

- `cpu.rag.prepare` capacity cannot accelerate the current sequential static chain.
- `GenerationConcurrency` currently changes in-task worker count through provider services and is not a closed canonical study factor.
- `combined_batch.go` and `embedding_batch.go` already split deterministic planning from one-batch execution, making them appropriate durable map boundaries.
- Query map materialization is fixed at eight today; runtime capacity is intentionally absent from `PipelineExecution`.
- Index parallelism requires a separate Bleve design proof; concurrent mutation must not be assumed safe or deterministic.

### What was tricky to build

The study design had to distinguish concurrency that changes scientific semantics from concurrency that changes only runtime scheduling. Batch size changes provider request boundaries, cache identity, retry scope, and preparation identity. Resource capacity should change timing and queueing but not semantic output identity. Researchctl `maxConcurrent` changes how many complete runs contend for the host and therefore must remain one during isolated scaling measurement.

Another difficulty was avoiding a huge Cartesian product. Generation, embedding, indexing, and querying occur in dependent stages, so most cross-stage combinations are wasteful. The guide uses frozen upstream artifacts and stage-separated studies, followed by a small integrated confirmation matrix.

### What warrants a second pair of eyes

- Whether the proposed capacity profile belongs in `ragcontract`, runner-only contracts, or a generic Workflow host profile.
- Whether durable generation/embedding maps should hard-cut the in-task worker pool or coexist temporarily for direct-engine execution.
- Which Bleve segment/shard strategy preserves exact retrieval semantics.
- Exact medium corpus work-unit targets after the cost-free sizing pass.
- Whether local SQLite contention needs further hardening before capacities above eight.
- Which stage metrics can be projected entirely from authoritative attempts and operations.

### What should be done in the future

- Implement the provider-free sizing pass and freeze nested TTC subsets.
- Define the closed capacity profile and delayed fixture maps.
- Prove output equivalence and active high-water marks before real provider calls.
- Run staged generation, embedding, indexing, query, and integrated studies.

### Code review instructions

- Start with Sections 5–9 of the guide to understand current versus target parallelism.
- Review Sections 12–17 for dataset sizing, factors, phases, replicates, and metrics.
- Review `pkg/ragworkflow/lower.go` and verify that static preparation is currently sequential.
- Review `pkg/ragoperators/combined_batch.go` and `embedding_batch.go` as the target durable batch APIs.
- Review Workflow V3 `dispatcher.go`, map integration tests, SQLite expansion, and reduction before implementing new maps.

### Technical details

Primary evidence commands included:

```text
rg -n "GenerationConcurrency|ExecuteCombinedPreparationBatch|PlanEmbeddingBatches" pkg -S
rg -n "Capacities|ResourceClass|Dispatcher" scraper/pkg/workflowv3* -S
rg -n "MaxConcurrent|ordering|replicates|sampleUnit" researchctl/pkg -S
```

Ticket root:

`ttmp/2026/07/24/TTC-RAG-PARALLELIZATION-STUDY--measure-rag-preparation-indexing-embedding-and-query-parallelization-before-full-corpus-execution/`

## Step 2: Validate and publish the study guide

This step validated ticket metadata and published the design guide, diary, and task list as one reMarkable bundle. The ticket remains active because implementation and experiment execution have not started.

The dry run confirmed the exact source files, PDF name, and destination. The real upload completed successfully.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Validate the ticket deliverables and upload them for offline review.

**Inferred user intent:** Make the complete scaling-study plan available before approving implementation or provider spend.

### What I did

- Related key RAG, Workflow V3, and Researchctl evidence files to the guide and diary.
- Updated the ticket changelog.
- Validated both document frontmatter blocks.
- Ran `git diff --check`.
- Ran `docmgr doctor --ticket TTC-RAG-PARALLELIZATION-STUDY --stale-after 30 --details`.
- Dry-ran and uploaded the bundled PDF.

### Why

- Cross-repository relations make architecture claims traceable.
- Validation prevents stale paths and malformed metadata from reaching the review copy.
- A bundled PDF keeps architecture, chronology, and implementation tasks together.

### What worked

- Frontmatter validation passed.
- Ticket doctor passed cleanly.
- Upload returned: `OK: uploaded TTC RAG Parallelization Study Guide.pdf -> /ai/2026/07/24/TTC-RAG-PARALLELIZATION-STUDY`.

### What didn't work

- N/A.

### What I learned

- The guide and tasks can be reviewed independently of any provider credentials or generated study artifacts.

### What was tricky to build

The bundle contains several Mermaid diagrams and dense technical tables. The dry run was necessary to verify that the long guide, diary, and task list were assembled under one stable name before upload.

### What warrants a second pair of eyes

- Diagram and table readability on the reMarkable device.
- The proposed durable-map scope before implementation begins.

### What should be done in the future

- Record design-review decisions in this diary.
- Re-upload only when replacement is intentional because forced replacement can remove annotations.

### Code review instructions

- Run targeted ticket doctor.
- Confirm all twelve implementation/execution tasks remain open.
- Review the PDF destination recorded below.

### Technical details

Remote document:

`/ai/2026/07/24/TTC-RAG-PARALLELIZATION-STUDY/TTC RAG Parallelization Study Guide.pdf`

## Step 3: Put the broad scaling study on hold

This step paused the broad parallelization study before implementation. A narrower prerequisite ticket, `RAG-WORKFLOW-CONCURRENCY-HARD-CUT`, now owns removal of duplicate scheduling and configuration paths.

The original tasks remain visible for future review, but each is marked on hold or superseded. This preserves the study design without implying that its additional runtime controls should be implemented now.

### Prompt Context

**User prompt (verbatim):** "Create a new design ticket for the minimal cleanup scope, and mark th eprevious document as "on hold for now".

Update the tasks accordingly too.."

**Assistant interpretation:** Pause the broad scaling design, move the minimal simplifying cutover into a separate design ticket, and make task status explicit.

**Inferred user intent:** Simplify concurrency ownership before adding experimental capacity and telemetry features.

### What I did

- Changed the ticket index, design guide, and diary status to `on-hold`.
- Added an explicit hold notice at the top of the design guide.
- Marked every task as on hold or superseded with the prerequisite named.
- Created and linked the focused hard-cut design ticket.

### Why

- The broader document contains useful future research design but also machinery that should not precede cleanup.
- Explicit hold markers prevent implementation from starting from stale scope.

### What worked

- Docmgr now recognizes the explicit `on-hold` status.
- Stable task IDs were preserved while descriptions were updated.

### What didn't work

- N/A.

### What I learned

- The scaling design remains useful as a future measurement plan, but its runtime architecture must be reviewed after the hard cut because several proposed controls may no longer be necessary.

### What was tricky to build

The task list had to preserve deferred work without presenting it as immediately actionable. Stable task IDs were retained, while descriptions now state the exact prerequisite or superseding ticket.

### What warrants a second pair of eyes

- Whether any broader scaling task should be deleted rather than resumed after the hard cut.
- Whether index/query scaling should eventually move to separate tickets.

### What should be done in the future

- Resume only after `RAG-WORKFLOW-CONCURRENCY-HARD-CUT` passes acceptance and the broader design is reviewed.

### Code review instructions

- Confirm the design guide begins with the on-hold notice.
- Run `docmgr task list --ticket TTC-RAG-PARALLELIZATION-STUDY` and verify every task is explicitly paused or superseded.

### Technical details

Prerequisite ticket: `RAG-WORKFLOW-CONCURRENCY-HARD-CUT`.
