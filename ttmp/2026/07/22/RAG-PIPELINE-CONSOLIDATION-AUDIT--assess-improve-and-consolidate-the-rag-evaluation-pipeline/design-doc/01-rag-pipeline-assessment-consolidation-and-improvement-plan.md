---
Title: RAG pipeline assessment, consolidation, and improvement plan
Ticket: RAG-PIPELINE-CONSOLIDATION-AUDIT
Status: review
Topics:
    - rag-eval
    - evaluation
    - workflow
    - ttc
    - research
    - visualization
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/scraper/pkg/workflowv3/external_operation.go
      Note: Generic durable operation contract
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/scraper/pkg/workflowv3sqlite/external_operation_query.go
      Note: Canonical operation query and export source
    - Path: repo://cmd/rag-ttc-v3-sweep/main.go
      Note: Current execution reduction and publication boundary
    - Path: repo://pkg/researchctladapter/operation_custody.go
      Note: Current researchctl artifact and metric bridge
ExternalSources: []
Summary: Architecture and phased implementation plan for one reproducible RAG execution-to-publication pipeline spanning Workflow V3 operation evidence, RAG reductions, researchctl custody, graph/report generation, and publication gates.
LastUpdated: 2026-07-22T18:32:54.534884826-04:00
WhatFor: Guide consolidation of the current ticket-local measurement and reporting flow into maintained, tested product capabilities.
WhenToUse: Before changing run evidence schemas, RAG analytics, researchctl custody metrics, graph/report generation, or publication automation.
---


# RAG pipeline assessment, consolidation, and improvement plan

## Executive summary

The RAG evaluation pipeline now has a strong execution core and a weak analysis/publication edge. Scraper Workflow V3 durably records external operation admission and completion. The TTC sweep preserves per-cell operation ledgers before deleting source-bearing runtime state. Researchctl can verify and import immutable artifacts. The remaining work—metric reduction, graph generation, narrative reporting, ticket reconciliation, and publication—is ticket-local and manual.

This fragmentation caused the July 22 real qualification to be declared complete with a report that omitted retries, undercounted one cell's elapsed time, did not embed graphs, could not answer the requested TPS question, and disagreed with its own ticket state. The raw evidence was sufficient; the finalization pipeline did not force anyone to use it.

The proposed architecture introduces one versioned **RAG run analysis** layer between Workflow evidence and publication. It derives all request arithmetic, retry-safe boundaries, latency distributions, concurrency, occupancy, missingness, throughput, and quality joins from immutable evidence. A maintained command generates canonical analysis, graphs, report sections, researchctl metrics, and a publication receipt. Scraper remains domain-neutral; researchctl remains downstream immutable custody.

## 1. Problem statement

### 1.1 Current ownership is correct but incomplete

Current boundaries are sound:

- scraper owns generic workflow execution and operation evidence;
- rag-evaluation-system owns generation, embedding, chunking, evaluation, and domain interpretation;
- researchctl owns immutable scientific runs, artifacts, metrics, and comparison.

The missing component is a maintained reducer/finalizer in rag-evaluation-system. The current sweep has its own aggregate structs and ticket-local Python renderer. Researchctl receives only four coarse counters. Markdown is handwritten from whichever subset an agent notices.

### 1.2 Success and failure use different evidence paths

The sweep builds successful cell timing from successful output artifacts and successful attempts. Failed cells use operation ledgers. This creates inconsistent semantics:

```text
success path: output artifact -> successful attempts -> aggregate timing
failure path: operation ledger -> failure reduction
```

The real batch-2/concurrency-1 cell proves the problem. A failed call started before the first successful attempt and disappeared from published makespan while remaining present in operation custody.

### 1.3 Missingness is represented as zero

Most real generation operations did not report actual token/cost usage. Numeric zeros in aggregate usage cannot distinguish:

- reported zero;
- unavailable;
- not applicable;
- conservatively reserved but not observed.

Any token/s or cost graph based on those values would be misleading.

### 1.4 Publication is not a transaction

A completed run currently requires manual coordination of:

- cell evidence;
- aggregate evidence;
- researchctl export/import;
- renderer invocation;
- visual review;
- report writing;
- diary/changelog/task/index updates;
- doctor;
- reMarkable dry-run/upload/listing.

No single gate checks that these artifacts agree. The previous publication uploaded successfully while the final task remained unchecked and the index still described a pre-run state.

## 2. Target architecture

```text
Canonical execution specification + authority envelope
                         |
                         v
                 rag-eval run command
                         |
                         v
             scraper Workflow V3 runtime
       attempts, budgets, operation admissions/completions
                         |
                         v
          canonical per-cell operation custody
          JSONL + manifest + cell checkpoint
                         |
                         v
               RAG analysis reducer v1
       boundaries, request states, distributions,
       concurrency, occupancy, missingness, quality joins
                         |
             +-----------+-------------+
             |                         |
             v                         v
      analysis.json/CSV         researchctl export
             |                  artifacts + metrics
             v                         |
      graph/report generator           v
             |                  immutable lab/import
             v
       publication validator
       receipt + ticket checks
             |
             v
      reMarkable / Obsidian / web
```

### 2.1 Scraper responsibility

Scraper continues to provide:

- durable pre-effect admission;
- immutable completion;
- attempt and operation identity;
- closed outcomes and failures;
- bounded counters;
- canonical operation JSONL and manifests;
- partial/failed-run export.

It must not understand chunks, TTC, answer quality, token efficiency, graph definitions, or report prose.

### 2.2 RAG analysis responsibility

Add a package such as `internal/runanalysis` or `pkg/runanalysis` and a command such as `rag-eval analyze-run`.

It owns:

- operation-kind interpretation;
- planned/admitted/succeeded/failed/incomplete counts;
- timing boundary definitions;
- latency and queue distributions;
- concurrency and occupancy;
- overlap and throughput;
- usage availability;
- links to RAG quality metrics;
- graph/report model;
- schema migration for RAG analysis artifacts.

### 2.3 Researchctl responsibility

Researchctl stores:

- verified operation ledgers and manifests;
- verified canonical RAG analysis artifact;
- verified graph/report manifests;
- selected scalar metrics for cross-run queries;
- immutable run/specification/attempt identity.

It does not duplicate every operation row in laboratory tables and does not participate synchronously in provider calls.

## 3. Canonical analysis model

### 3.1 Top-level record

```go
type RunAnalysis struct {
    SchemaVersion string             `json:"schemaVersion"`
    Source        SourceIdentity     `json:"source"`
    Boundary      BoundaryDefinition `json:"boundary"`
    Totals        OperationTotals    `json:"totals"`
    Missingness   Missingness        `json:"missingness"`
    Cells         []CellAnalysis     `json:"cells"`
    Graphs        []GraphDefinition  `json:"graphs"`
}
```

`SourceIdentity` contains input artifact digests, provider profile/model digests, workflow plan digest, operation manifest digests, and analysis implementation version.

### 3.2 Request-state arithmetic

Every operation kind reports:

```go
type RequestCounts struct {
    Planned    int64 `json:"planned"`
    Admitted   int64 `json:"admitted"`
    Succeeded  int64 `json:"succeeded"`
    Failed     int64 `json:"failed"`
    Canceled   int64 `json:"canceled"`
    TimedOut   int64 `json:"timedOut"`
    Incomplete int64 `json:"incomplete"`
}
```

Validation requires:

```text
admitted = succeeded + failed + canceled + timedOut + incomplete
```

`planned` is independent and may differ because of retries.

### 3.3 Boundary definitions

```go
type BoundaryDefinition struct {
    Primary string `json:"primary"` // operation-inclusive
    Available []string `json:"available"`
}

type CellBoundaries struct {
    SuccessfulAttemptMicros int64
    AllAttemptMicros        int64
    OperationInclusiveMicros int64
    InvocationMicros        *int64
}
```

The primary performance denominator is operation-inclusive elapsed. The successful-attempt boundary remains available only for historical comparison and must never be labeled simply “makespan” without qualification.

### 3.4 Missingness-aware usage

```go
type ObservedCounter struct {
    Value        int64  `json:"value"`
    Availability string `json:"availability"` // reported, unavailable, n/a
    Accounting   string `json:"accounting"`   // actual, conservative, none
}
```

Reducers must never convert unavailable values to observed zero. Aggregate counters include reported and missing operation counts.

### 3.5 Distribution model

Store enough values for reproducibility or deterministic summaries with denominator:

```go
type Distribution struct {
    Count int64
    Min, P50, P95, Max float64
    Mean, Stddev, CV float64
    Unit string
    Censored int64
}
```

For small samples, reports explicitly state count and avoid confidence claims.

## 4. Analysis algorithms

### 4.1 Interval sweep

For provider concurrency and coverage:

```text
for each completed operation:
    start = provider_started_at
    end   = completed_at
    events += (start, +1), (end, -1)

sort by timestamp; process end before start on ties
integrate active count over time
```

Derive:

- union duration with any provider active;
- peak all-provider concurrency;
- peak by operation kind;
- configured-slot occupancy;
- generation/embedding overlap;
- idle intervals longer than a configured threshold.

### 4.2 Retry-safe elapsed

```text
start = min(
    all attempt starts,
    all operation admitted_at
)
end = max(
    all attempt finishes,
    all operation completed_at
)
inclusive_elapsed = end - start
```

Incomplete operations require an explicit as-of boundary; they do not receive invented completion times.

### 4.3 Speedup and throughput

```text
chunk_throughput = completed_chunks / operation_inclusive_elapsed
speedup(c2) = elapsed(batch, c1) / elapsed(batch, c2)
efficiency(c2) = speedup / 2
```

Every comparison records replicate count and whether either cell experienced failures.

### 4.4 Quality join

Performance recommendations require output quality. Join on immutable cell/run identity and calculate:

- valid structured-output rate;
- generated item count and duplication;
- RAG retrieval metrics;
- answer/citation quality;
- malformed-output and retry rates.

A configuration is Pareto-dominated if another configuration is faster, no more expensive, and no worse on declared quality metrics.

## 5. Graph contract

Graphs derive only from canonical analysis JSON. Each manifest entry contains:

```json
{
  "id": "generation-latency",
  "sourceDigest": "sha256:...",
  "metricDefinitions": ["generation.provider_latency_seconds"],
  "title": "...",
  "altText": "...",
  "files": ["generation-latency.svg", "generation-latency.png"]
}
```

Required figure families:

1. elapsed boundary comparison;
2. throughput and speedup;
3. latency distributions with failure markers;
4. occupancy, provider coverage, and retries with separate denominators;
5. per-cell timelines with common axes;
6. usage/missingness rather than false token/cost zeros;
7. quality/performance frontier when quality data exists.

Validation checks dimensions, nonblank rendering, referenced file existence, and manifest/source digest agreement. Publication runs retain manual visual QA.

## 6. Report contract

A complete report contains:

1. executive summary;
2. scope and evidence provenance;
3. matrix and execution order;
4. request-state arithmetic;
5. corrected result table;
6. every required graph embedded with interpretation;
7. retries and failure analysis;
8. missingness and unsupported claims;
9. quality results;
10. limitations and replication requirements;
11. operational recommendation;
12. reproduction commands and artifact references.

Numeric tables should be generated from analysis JSON. Reviewed prose may be hand-authored, but a validator extracts declared metrics or compares generated fragments to prevent drift.

## 7. Publication receipt

`rag-eval finalize-run` should emit:

```json
{
  "schemaVersion": "rag-eval-publication-receipt/v1",
  "sourceDigest": "sha256:...",
  "analysisDigest": "sha256:...",
  "graphManifestDigest": "sha256:...",
  "reportDigest": "sha256:...",
  "privacyScan": "passed",
  "ticketValidation": "passed",
  "createdAt": "..."
}
```

The command fails if:

- operation manifests do not verify;
- request arithmetic disagrees;
- report images are missing;
- graphs derive from another source digest;
- unavailable usage is graphed as zero;
- the report lacks limitations;
- ticket tasks/index/status disagree;
- doctor fails;
- privacy canaries match.

reMarkable/Obsidian publication consumes the receipt and performs dry-run first.

## 8. Decision records

### Decision: Keep raw operation evidence in scraper; put reductions in RAG

- **Context:** Operation persistence is generic, while chunk throughput and quality are domain-specific.
- **Options considered:** all analytics in scraper; all operations copied into researchctl; RAG reducer over scraper artifacts.
- **Decision:** RAG reducer over immutable scraper operation artifacts.
- **Rationale:** Preserves package boundaries and centralizes domain semantics.
- **Consequences:** RAG analysis schema needs explicit versioning and tests.
- **Status:** proposed.

### Decision: Use operation-inclusive elapsed as primary performance boundary

- **Context:** successful-attempt makespan omitted a failed provider call.
- **Options considered:** keep historical makespan; all-attempt span; operation-inclusive span.
- **Decision:** operation-inclusive is primary; retain named historical boundaries.
- **Rationale:** It includes durable failed work and matches observed effect time.
- **Consequences:** historical throughput values may change for retry cells; reports must name versions.
- **Status:** proposed.

### Decision: Do not infer token or cost values

- **Context:** provider usage is unavailable for most calls.
- **Options considered:** use reservations; treat zero as zero; model estimates; explicit missingness.
- **Decision:** explicit missingness; reservations remain authority only.
- **Rationale:** Estimated or reserved usage is not observed billing/performance evidence.
- **Consequences:** TPS/cost graphs remain N/A until adapters report usage.
- **Status:** accepted.

### Decision: Generate reports and graphs from one analysis record

- **Context:** independent scripts and prose drifted.
- **Options considered:** maintain manual scripts; notebook; deterministic command.
- **Decision:** deterministic command with manifest and optional reviewed prose overlay.
- **Rationale:** One source prevents inconsistent tables, graphs, and report claims.
- **Consequences:** graph/report templates become maintained code with regression tests.
- **Status:** proposed.

### Decision: Make publication fail closed

- **Context:** the old report uploaded despite stale ticket state and missing analysis.
- **Options considered:** checklist; reviewer reminder; executable receipt.
- **Decision:** executable publication receipt.
- **Rationale:** Existing prose checklists did not prevent premature completion.
- **Consequences:** publication has an additional explicit validation step.
- **Status:** proposed.

## 9. Alternatives rejected

### Keep ticket-local Python scripts

Useful for exploration, but rejected as the final pipeline because schemas, labels, and semantics drift per ticket.

### Put report generation in researchctl

Rejected because researchctl is domain-neutral and should not know TTC batch semantics. It can host generic report plumbing later, but RAG owns metric interpretation.

### Use only successful task outputs

Rejected because failures, cancellations, and incomplete calls are scientifically important and already durable in the operation ledger.

### Adopt OpenTelemetry as primary evidence

Rejected as primary custody because sampling/export/retention are weaker than Workflow's durable operation contract. OTel mirroring may complement the ledger for live diagnosis.

### Treat provider coverage as CPU utilization

Rejected because operation spans measure active client calls, not remote hardware utilization.

## 10. Implementation phases

### Phase 1: Golden reduction package

Files:

- new `internal/runanalysis/` package;
- fixtures copied or generated from bounded operation records;
- command tests.

Work:

1. Define analysis v1 schema.
2. Implement strict decode and manifest verification.
3. Implement request arithmetic and boundaries.
4. Implement interval/distribution/missingness reductions.
5. Add retry, incomplete, cancellation, and no-usage fixtures.
6. Use the July 22 real run as a golden regression without embedding sensitive data.

Exit criteria: generated summary matches this ticket's `summary.json` and catches the 21.472-second boundary undercount.

### Phase 2: Evidence v3

1. Add explicit operation totals and boundary fields to new sweep outputs.
2. Add usage availability.
3. Keep prior artifacts immutable; no compatibility adapter unless explicitly approved.
4. Version CSV headers and graph manifest.

Exit criteria: new aggregate evidence can answer planned/admitted/succeeded/failed/incomplete without reopening JSONL, while JSONL remains authoritative.

### Phase 3: Graph/report command

1. Move prototype graph semantics into maintained code/tooling.
2. Generate figures and manifest.
3. Generate numeric report sections.
4. Add HTML/md-view smoke validation for image loading and headings.
5. Add direct image QA checklist.

Exit criteria: one command produces the same figures and tables from canonical analysis.

### Phase 4: Researchctl metrics

1. Add richer scoped metrics to `OperationCustodyExportInput` construction.
2. Add analysis and graph manifests as verified artifacts.
3. Verify fresh import and re-export.
4. Add comparison queries.

Exit criteria: researchctl can rank runs by throughput, retry rate, latency, quality, and usage availability.

### Phase 5: Finalization and publication

1. Implement receipt validator.
2. Integrate docmgr task/index/doctor checks.
3. Integrate privacy canary scan.
4. Integrate reMarkable dry-run/upload/list verification.
5. Store receipt as a verified artifact.

Exit criteria: reproducing the old stale-ticket/shallow-report state causes finalization to fail.

### Phase 6: Replication support

1. Add replicate count and randomized/counterbalanced scheduling.
2. Persist execution order and provider/environment snapshots.
3. Add cell-level confidence intervals.
4. Join quality metrics.

Exit criteria: configuration decisions no longer rely on one sequential replicate.

## 11. Testing strategy

### Unit tests

- request arithmetic reconciliation;
- percentile interpolation;
- interval union and peak ties;
- retry-safe boundaries;
- incomplete operation handling;
- missing usage never emitted as observed zero;
- stable canonical serialization.

### Integration tests

- fixture sweep → operation export → analysis → graphs/report → researchctl import;
- failed cell retains prior successful cells and failure operations;
- graph manifest source digest changes when evidence changes;
- report references every required image;
- finalize fails on stale ticket task or index.

### Visual tests

- all images have nonzero dimensions;
- md-view reports zero image load failures;
- common timeline axes are asserted in graph metadata;
- manual visual review for publication runs.

### Privacy tests

Scan SQLite/WAL, operation artifacts, analysis, CSV, graph labels, Markdown, HTML, and publication bundle for source/provider/credential canaries.

## 12. Risks

### Schema proliferation

Mitigation: one RAG analysis schema, one manifest, and explicit version ownership. Do not add parallel ticket formats.

### False precision

Mitigation: include replicate counts, denominators, missingness, and operation boundary names with every metric.

### Heavy graph dependencies in Go builds

Mitigation: keep renderer as a packaged Python or web tool initially, but make canonical analysis and manifest contracts language-neutral. The critical consolidation is semantic, not language purity.

### Publication gate becomes environment-specific

Mitigation: separate core validation receipt from optional publishers. Core receipt is deterministic; publisher receipts append delivery evidence.

### Historical compatibility pressure

Mitigation: preserve old artifacts as immutable evidence and produce new versioned analysis. Do not add hidden adapters without explicit approval.

## 13. Open questions

1. Should `runanalysis` be public `pkg/` API or internal until the schema stabilizes?
2. Should graph generation remain Python, move to a web renderer, or use a Go plotting library?
3. Which quality metrics are mandatory before a performance recommendation is publishable?
4. Should provider transport failures gain bounded subcodes for tunnel, timeout, upstream status class, and malformed response?
5. Should researchctl gain a generic comparison/report query layer, leaving only domain metric definitions in RAG?
6. Should execution ordering support a seeded randomization mode with the seed in canonical specification?
7. What explicit authority model should cover replicated paid runs and retries?

## 14. Acceptance criteria

The consolidation is complete only when:

1. Retry-safe elapsed includes every durable operation admission/completion.
2. Request-state arithmetic is explicit and reconciled.
3. Token/cost missingness is represented structurally.
4. One canonical analysis drives CSV, graphs, report, and researchctl metrics.
5. Performance and quality can be joined by immutable cell identity.
6. Graphs are manifested, visually reviewed, and rendered in md-view without load errors.
7. Ticket state, doctor, report, and publication receipt agree.
8. Failed/partial runs remain publishable as failed evidence rather than being discarded.
9. No sensitive payload appears in analysis or publication artifacts.
10. A replicated study can produce cell-level uncertainty rather than single-run claims.

## 15. References

- [Full real-run performance audit](../analysis/01-ttc-real-run-performance-audit-and-pipeline-consolidation-assessment.md)
- [Investigation diary](../reference/01-investigation-diary.md)
- [Derived analysis source](../sources/derived-real-attempt-003/summary.json)
- [Graph manifest](../sources/derived-real-attempt-003/graphs/manifest.json)
- original ticket `RAG-TTC-V3-SWEEP`
- scraper ticket `SCRAPER-WORKFLOW-V3-EXTERNAL-OPERATIONS`
- `cmd/rag-ttc-v3-sweep/main.go`
- `pkg/researchctladapter/operation_custody.go`
- scraper `pkg/workflowv3/external_operation.go`
- scraper `pkg/workflowv3sqlite/external_operation_query.go`
