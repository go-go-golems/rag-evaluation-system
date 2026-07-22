---
Title: Workflow V3 Umans Batching and Concurrency Study Design and Implementation Guide
Ticket: RAG-TTC-V3-SWEEP
Status: active
Topics:
    - rag-eval
    - evaluation
    - workflow
    - chunking
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://cmd/rag-ttc-v3-sweep/main.go
      Note: Executable fixture study and canonical exporter
    - Path: repo://internal/workflowv3ttc/provider.go
      Note: Precise batch provider adapter and checked usage
    - Path: repo://internal/workflowv3ttc/sweep.go
      Note: Exact matrix planner and batch materialization
    - Path: repo://internal/workflowv3ttc/sweep_workflow.js
      Note: Workflow V3 sweep graph and budget
    - Path: repo://pkg/ragoperators/combined_batch.go
      Note: Existing deterministic domain batching and provider-call authority
    - Path: repo://ttmp/2026/07/22/RAG-TTC-V3-SWEEP--workflow-v3-umans-batching-and-concurrency-study/scripts/01-render-sweep.py
      Note: Deterministic graph renderer
    - Path: repo://ttmp/2026/07/22/RAG-TTC-V3-SWEEP--workflow-v3-umans-batching-and-concurrency-study/sources/fixture-control/evidence.json
      Note: Executed 12-cell fixture evidence
ExternalSources: []
Summary: Exact, bounded Workflow V3 experiment for measuring how chunks per Umans request and request concurrency affect latency, throughput, tokens, cost, reliability, and downstream embedding work.
LastUpdated: 2026-07-22T10:15:00-04:00
WhatFor: Implement, operate, review, and reproduce the Umans batching and concurrency sweep without leaking provider authority or confusing queue time with provider time.
WhenToUse: Before changing TTC generation batch size or concurrency, running a paid sample, interpreting its graphs, or reproducing the study.
---


# Workflow V3 Umans Batching and Concurrency Study Design and Implementation Guide

## Executive summary

This study measures two independent factors in TTC combined preparation: the number of source chunks placed into one Umans LLM request and the maximum number of concurrent Umans requests. The accepted factor levels are batch sizes `1, 2, 4, 8` and concurrency limits `1, 2, 4`. Four is a hard ceiling imposed by the operator's provider plan, not merely a chart label. Every cell uses the same ordered source-chunk prefix, exact model, prompt, output schema, provider profile, task bundle, registry generation, policy, and workflow plan identities.

The study runs in two stages. A no-cost fixture control first proves batching, admission, accounting, output validation, canonical evidence, graph generation, and privacy. A real Umans smoke run follows only when host-local credentials are available and a finite request/token/cost budget is explicitly authorized. The initial real matrix uses 16 chunks, four batch sizes, and concurrency 1 and 4: 60 requests in total. The larger 32-chunk, three-concurrency matrix costs 180 requests per replicate and is not automatically authorized by the smoke run.

Workflow V3 owns durable attempts, admission, retries, budget reservations, fencing, timestamps, and immutable output references. `ragoperators` remains the domain authority for deterministic batch planning, provider calls, response validation, and usage. The measurement exporter reads bounded Workflow V3 attempt evidence and bounded provider spans into canonical JSON and CSV. A deterministic Python renderer creates SVG and PNG plots from that evidence; charts are derived views, never primary evidence.

## Problem statement and scope

The prior diagnostic v9 run proves that provider work occurred but cannot answer precise performance questions. Its durable operation creation timestamp precedes queue admission, and successful lease rows are removed, so `updated_at - created_at` combines queue wait, dependency wait, and provider execution. It also lacks complete token usage. It supports coarse completion throughput, not request latency or token throughput.

The study must answer:

1. Does batching 2, 4, or 8 chunks per request increase chunks/second or reduce cost/chunk?
2. How does batch size affect request latency, first completion, tail latency, token rate, malformed output, and representation completeness?
3. How much throughput is gained at concurrency 2 and 4?
4. Does the dispatcher remain work-conserving while never exceeding four active Umans requests?
5. How much downstream embedding work is induced by each generation cell?
6. Are output identity, citations, accounting, and privacy preserved across all cells?

The study does not select a production configuration solely from speed. Any candidate must also pass response-schema validation, exact chunk coverage, citation lineage, retry limits, cost limits, and prepared-corpus validation.

## Current-state architecture and evidence

### Domain batching already exists

`pkg/ragoperators/combined_batch.go:11-32` defines immutable `CombinedPreparationPlan` and `CombinedPreparationBatch` contracts. `PlanCombinedPreparation` validates `BatchSize`, `QuestionsPerChunk`, and `MaxBatchRunes`, then deterministically partitions sorted chunks (`combined_batch.go:42-60`). `ExecuteCombinedPreparationBatch` resolves the exact model, prompt, and output schema before making one provider request (`combined_batch.go:63-80`). The sweep must reuse this authority rather than creating a second prompt or batching implementation.

### The current Workflow V3 adapter is single-chunk

`internal/workflowv3ttc/provider.go:46-76` accepts one `Chunk`, requires planning to produce exactly one batch, executes that batch, and converts the environment usage delta into bounded Workflow V3 usage. This is correct for the production-cardinality preparation graph but cannot represent the batching factor. The sweep therefore needs a distinct batch-item contract and task kind rather than changing existing production semantics.

### Usage accounting exists but timing needs an explicit provider span

`internal/workflowv3ttc/provider.go:115-149` converts input tokens, output tokens, embedding tokens, and provider cost to checked integer dimensions. Workflow V3 attempts already retain lease-time `StartedAt` and terminal timestamps. The sweep adds monotonic provider-call start/end measurement inside the RAG-owned adapter, returning only duration and bounded usage—not prompts, responses, headers, URLs, or credentials. This separates queue wait and task overhead from provider wall time.

### Workflow admission and budgets already exist

`internal/workflowv3ttc/production_workflow.js:6-21` defines separate generation, embedding, and evaluation budgets. Its generation and embedding maps attach independent resources and transactional budgets (`production_workflow.js:53-92`). The sweep uses the same admission model but pins an effective generation capacity per cell to 1, 2, or 4.

## Experimental contract

### Factors

| Factor | Levels | Meaning |
|---|---:|---|
| `chunks_per_request` | 1, 2, 4, 8 | Maximum source chunks in one deterministic combined-preparation request |
| `generation_concurrency` | 1, 2, 4 | Maximum active Umans provider calls for the cell |
| `replicate` | 1 initially | Fresh run and artifact namespace; additional replicates require separate authorization |

### Controlled identities

Each cell pins and records:

- ordered source-manifest digest and selected chunk keys;
- generation node digest, model digest, prompt digest, and output-schema digest;
- provider profile digest and non-secret endpoint capability digest;
- task kind/version, bundle digest, entrypoint, ABI, registry generation;
- requested/effective resource claims;
- retry and budget policy digests;
- Go binary/repository revision and host measurement-clock metadata;
- embedding profile/model digest when downstream embedding is enabled.

Credentials, endpoint URLs, headers, prompt bodies, source text, provider response bodies, and vectors are forbidden in the ledger, workflow database, CSV, graphs, and report.

### Stages and exact request counts

#### Fixture control

The fixture control runs all batch sizes and concurrency levels over deterministic nonsensitive chunks. It validates mechanics, not Umans performance.

#### Real smoke matrix

With 16 chunks:

| Batch size | Requests per concurrency level |
|---:|---:|
| 1 | 16 |
| 2 | 8 |
| 4 | 4 |
| 8 | 2 |

Concurrency levels 1 and 4 therefore require `(16 + 8 + 4 + 2) × 2 = 60` generation requests. A cell may retry only according to its pinned retry policy; retry reservations count against the hard request budget.

#### Main matrix

With 32 chunks and concurrency 1, 2, and 4, one replicate requires `(32 + 16 + 8 + 4) × 3 = 180` generation requests. It is a separate operator decision after smoke evidence and estimated cost are reviewed.

### Randomization and warm-up

The ordered chunk subset is frozen before execution. Batch membership is deterministic. Cell execution order is deterministically counterbalanced so model/provider drift is not perfectly confounded with increasing batch size. A separately labeled warm-up request may be used only if authorized and excluded from estimates; it still counts against budget and remains ledger evidence.

## Measurement model

For each attempt, capture:

- run, cell, batch, attempt, and immutable identity keys;
- chunk count and aggregate input rune count;
- enqueued, leased, task-start, provider-start, provider-end, and committed timestamps/durations;
- input tokens, output tokens, total tokens, and cost microunits;
- output representation count and per-kind counts;
- status, typed failure code/class, retryability, and attempt number;
- requested/effective generation capacity;
- embedding batch count, embedding tokens, and embedding duration when enabled.

Provider timing uses `time.Now` only to calculate elapsed monotonic duration in-process. Durable UTC timestamps establish ordering. The result contains integer microseconds so canonical JSON does not depend on floating-point duration encoding.

Derived per-cell metrics include:

- request latency p50/p95/p99;
- queue wait and commit overhead p50/p95/p99;
- successful requests/minute;
- chunks/second and representations/second;
- input, output, and total tokens/second;
- cost/request, cost/chunk, and cost/representation;
- malformed-output and retry rates;
- observed peak active provider calls;
- makespan, time to first completion, and time to last completion;
- generation/embedding overlap and utilization.

Percentiles use a documented deterministic nearest-rank method. Failed attempts appear in timeline/failure views but are excluded from successful-latency distributions; both counts are always displayed.

## Canonical evidence and graphs

The exporter writes:

```text
artifacts/<study-id>/
  manifest.json
  attempts.jsonl
  cells.csv
  timeline.csv
  privacy-scan.json
  graphs/
    throughput-by-batch-and-concurrency.svg
    latency-by-batch-and-concurrency.svg
    token-rate-by-batch-and-concurrency.svg
    cost-efficiency.svg
    request-timeline.svg
    generation-embedding-overlap.svg
```

Canonical JSON/JSONL and CSV are primary evidence. SVG is deterministic where practical; PNG is a convenience rendering. Every graph embeds the study manifest digest, generation timestamp, and explicit units. Missing values are shown as missing rather than zero.

## Proposed implementation

### Packages

Create a RAG-owned package such as `internal/workflowv3ttcsweep` with:

- `contracts.go`: canonical study, cell, batch item, bounded task result, and evidence schemas;
- `plan.go`: matrix validation, deterministic cell ordering, request-count and budget calculation;
- `provider.go`: adapter over `ragoperators.PlanCombinedPreparation` and `ExecuteCombinedPreparationBatch`, with precise bounded spans;
- `bundle.go`, `workflow.js`, `tasks.cjs`: exact Workflow V3 bundle and graph;
- `runner.go`: fixture/real host wiring and sequential cell execution;
- `export.go`: canonical attempt and summary export;
- `privacy.go`: forbidden-canary and credential-pattern scan.

Place the operator script in the ticket `scripts/` directory, numbered to retain run order. Place generated study evidence under ticket `sources/` or an external artifact root with compact ticket references.

### Workflow shape

```text
frozen chunk manifest
  -> deterministic batch manifest for one cell
  -> lazy map(generate batch)
       resource: umans.generation.remote, capacity = cell concurrency
       budget: requests + tokens + cost
       retry: malformed/transient only, bounded
  -> bounded validation/reduction
  -> optional embedding map
       resource: rag.embedding.local/remote, independently admitted
  -> cell evidence root
```

Cells run in isolated run IDs and budgets. They do not share provider responses or generated-representation caches because cache hits would invalidate throughput comparisons. Source materialization may be shared by immutable digest.

### API sketch

```go
type Cell struct {
    ChunksPerRequest     int
    GenerationConcurrency int
    Replicate            int
}

type GenerationMeasurement struct {
    ProviderElapsedMicros int64
    InputTokens           int64
    OutputTokens          int64
    CostMicrounits        int64
    ChunkCount            int
    RepresentationCount   int
    PeakActiveObserved    int
}

func PlanMatrix(spec StudySpec) ([]Cell, RequestBudget, error)
func RunCell(ctx context.Context, host HostAuthorities, cell Cell) (EvidenceRef, error)
func ExportStudy(ctx context.Context, store Store, refs []EvidenceRef) (Manifest, error)
func RenderGraphs(manifestPath string) error
```

## Decision records

### Decision: Use a dedicated batch task instead of altering production single-chunk behavior

- **Context:** Current TTC Workflow V3 task identity consumes exactly one chunk, while the experiment needs one provider request to consume multiple chunks.
- **Options considered:** Change the production task contract; group workflow items before the existing task; create a dedicated versioned sweep task.
- **Decision:** Create a dedicated versioned sweep batch task that calls the existing `ragoperators` planner/executor.
- **Rationale:** This preserves existing production semantics and retry evidence while avoiding duplicate domain logic.
- **Consequences:** The sweep has an additional exact bundle/task identity, but results remain directly traceable to the production domain operator.
- **Status:** accepted

### Decision: Enforce concurrency in Workflow V3 admission

- **Context:** Provider-internal fan-out would hide actual concurrency from durable resource projections.
- **Options considered:** Provider SDK concurrency; goroutines inside a task; one request per Workflow V3 attempt with resource admission.
- **Decision:** Set provider-internal concurrency to one and represent every request as one admitted Workflow V3 attempt.
- **Rationale:** The store can prove that active calls never exceed four, and cancellation/fencing/accounting remain per request.
- **Consequences:** Runtime overhead is measurable and can be separated from provider wall time.
- **Status:** accepted

### Decision: Separate fixture control, smoke authorization, and main authorization

- **Context:** The main matrix can make 180 requests per replicate and the provider has monetary and plan limits.
- **Options considered:** Run the full matrix immediately; run only fixture data; stage increasingly expensive runs.
- **Decision:** Require fixture success, then separately authorize a 60-request real smoke matrix, then separately authorize main replicates.
- **Rationale:** This minimizes spend and prevents invalid broad runs.
- **Consequences:** The ticket may remain blocked after local completion until credentials and numeric cost authority exist.
- **Status:** accepted

### Decision: Keep canonical evidence independent of plotting libraries

- **Context:** Charts are useful but renderer versions and binary image encoding can vary.
- **Options considered:** Store only images; store a notebook; store canonical tables and derive images.
- **Decision:** JSONL/CSV are primary; graphs are reproducible derived artifacts.
- **Rationale:** Reviewers can recompute every statistic and use alternate visualization tools.
- **Consequences:** Export schema and renderer both require tests.
- **Status:** accepted

## Failure, restart, and cancellation behavior

- A retry is a new immutable attempt and consumes request budget if the provider call began.
- Configuration, profile, prompt, schema, model, and identity mismatches fail permanently before provider admission.
- Malformed provider output is typed, bounded, and retryable only within the pinned policy.
- Cancellation stops new leases, cancels active HTTP calls, and leaves durable attempt evidence.
- Lease loss fences stale completion and artifact publication.
- Restart reopens the same run and never duplicates successful batch keys.
- A provider span with no committed result remains infrastructure evidence; reconciliation must not count it as a successful sample.
- Any observed peak concurrency above the effective cell limit invalidates the study.

## Privacy and security

The workflow database and evidence are scanned for source canaries, environment values, credential prefixes, Authorization headers, endpoint URLs, prompt text, provider bodies, and vectors. Only source artifact references, compact identities, aggregate rune counts, usage integers, timing integers, typed failures, and output content digests may persist. Graph labels use stable cell IDs rather than provider URLs.

Real-provider configuration is loaded from a host-only YAML file or environment. The file is never copied into the ticket. The command reports only a non-secret capability/profile digest. The script exits before submission if required authority or an exact budget ceiling is absent.

## Implementation phases

1. **Contract and matrix planner:** validate factor levels, request arithmetic, counterbalanced order, budget ceilings, and canonical study identity.
2. **Batch provider adapter:** reuse `ragoperators`, add precise spans, bounded usage, output validation, and fixture tests.
3. **Workflow V3 bundle:** implement batch lazy map, hard resource capacity, retries, budgets, reduction, and evidence root.
4. **Exporter and renderer:** produce canonical JSONL/CSV and deterministic graphs; test percentile and concurrency calculations.
5. **Fixture control:** execute matrix, inject known delays/failures, prove exact timings, peak concurrency, restart, privacy, and deterministic exports.
6. **Real smoke preflight:** verify profile identities, dry-run request count, token/cost maximums, source subset, credentials, and operator authority.
7. **Real smoke:** run exactly the authorized matrix, analyze evidence, and stop.
8. **Optional main study:** only after a new explicit decision based on observed smoke cost/reliability.
9. **Publication:** attach immutable evidence to researchctl, complete report, validate docmgr, and upload a bundled PDF to reMarkable.

## Validation strategy

Unit tests must cover matrix arithmetic, invalid levels, batch partitioning, context/rune bounds, percentile rules, cost overflow, canonical encoding, and privacy redaction. Integration tests must cover resource capacities 1/2/4, observed peak admission, independent embedding refill, retries, cancellation, stale completion, close/reopen, cache-disabled behavior, exact chunk coverage, and deterministic evidence across completion orders.

Before real submission:

```text
- all fixture tests and race tests pass;
- request count equals the dry-run count;
- requested/effective generation capacity never exceeds 4;
- hard request/token/cost budgets exist;
- exact profile/model/prompt/schema digests resolve;
- provider credentials resolve without being printed;
- privacy scan passes;
- graph renderer reproduces fixture outputs;
- operator confirms the numeric real-run ceiling.
```

## Operator runbook

1. Run the matrix command with `--profile fixtures --dry-run`.
2. Run the complete fixture control and render graphs.
3. Review `manifest.json`, `cells.csv`, privacy scan, and fixture graphs.
4. Configure the host-only Umans profile.
5. Run `--profile real --dry-run`; record its exact request/token/cost limits.
6. Obtain explicit operator approval for those numeric limits.
7. Run the real smoke once. Do not auto-repeat a failed whole matrix.
8. Validate every cell, inspect cost and malformed-output evidence, and publish immutable artifacts.
9. Decide separately whether to run the 180-request main replicate.

## Risks, alternatives, and open questions

- Batch size 8 may exceed a configured rune/token limit for some chunks. Planning must reject or deterministically split; it must not silently change the factor.
- Server-side scheduling may vary over time. Counterbalancing reduces but cannot eliminate drift; replicates are needed before strong conclusions.
- Provider usage may omit token fields. Such cells can still report request/chunk throughput but must display token throughput as unavailable.
- Shared caches or HTTP connection warm-up can confound results. Provider-result caches are disabled; connection behavior and warm-up policy are recorded.
- The exact monetary ceiling for the real smoke is unresolved. Concurrency approval alone is not cost approval.
- Real provider credentials and endpoints are currently unavailable in this process.

## Hard acceptance criteria

The fixture control is complete only when all 12 cells execute with exact batch membership, no active-call count above four, canonical measurements, deterministic charts, restart/fencing/cancellation tests, and clean privacy scans. The real smoke is complete only when exactly 60 planned generation requests plus only policy-permitted retries remain within an explicitly approved budget, all results are validated, precise provider timing and usage are present or explicitly missing, graphs are generated from immutable evidence, and researchctl custody is verified.

## References

- `pkg/ragoperators/combined_batch.go`: deterministic combined preparation batching and provider execution.
- `pkg/ragoperators/types.go`: provider environment, usage, and domain records.
- `internal/workflowv3ttc/provider.go`: existing exact-profile Workflow V3 provider adapter and checked usage conversion.
- `internal/workflowv3ttc/production_workflow.js`: production resource, budget, map, gate, and reduction patterns.
- `internal/workflowv3ttc/contracts.go`: compact TTC Workflow V3 contracts.
- `internal/workflowv3ttc/manifest.go`: exact run identity conventions.
- Scraper `pkg/workflowv3/types.go`: canonical attempt timestamps and resource identities.
- Scraper `pkg/workflowv3sqlite`: durable attempts, admission, budgets, fencing, and projections.
