---
Title: TTC scripted experiment acceptance design and implementation guide
Ticket: TTC-SCRIPTED-EXPERIMENT-ACCEPTANCE
Status: active
Topics:
    - ttc
    - rag
    - rag-eval
    - evaluation
    - scripting
    - intern-guide
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://experiments/real-provider-v2/study-flash-combined-speed.js
      Note: Existing scripted RAG experiment reference
    - Path: repo://experiments/ttc-scripted/analysis.js
      Note: Maintained Researchctl analysis replacing ticket-local Python
    - Path: repo://experiments/ttc-scripted/study.js
      Note: Canonical workload-only study replacing the deleted custom sweep
    - Path: repo://ttmp/2026/07/22/RAG-PIPELINE-CONSOLIDATION-AUDIT--assess-improve-and-consolidate-the-rag-evaluation-pipeline/analysis/01-ttc-real-run-performance-audit-and-pipeline-consolidation-assessment.md
      Note: Retry-aware acceptance evidence
ExternalSources: []
Summary: Final acceptance design proving the converged platform can express, run, analyze, and publish TTC without bespoke Go or Python orchestration.
LastUpdated: 2026-07-22T23:15:00-04:00
WhatFor: Provide an end-to-end goal and executable acceptance contract for every framework ticket.
WhenToUse: Use when reviewing whether the convergence program is actually complete; implement only after predecessor tickets pass.
---




# TTC scripted experiment acceptance design and implementation guide

## Program context

This is the final child of **EXPERIMENT-PLATFORM-CONVERGENCE**. It depends on every sibling: `RESEARCHCTL-EXPERIMENT-PLANS`, `SCRAPER-WORKFLOW-V3-PRODUCT-CUTOVER`, `EXPERIMENT-PLATFORM-SCRAPER-RUNNER`, `SCRAPER-WORKFLOW-OBSERVATIONS`, `RAG-V2-WORKFLOW-LOWERING`, `RAG-GEPPETTO-WORKFLOW-OPERATIONS`, `RESEARCHCTL-EXPERIMENT-ANALYSIS`, and `RAG-V2-EXECUTION-CUTOVER`. All guides are in the three repositories' `ttmp/2026/07/22/` directories.

## Executive summary

TTC is the acceptance workload, not a framework. It should contain corpus/evaluation inputs, reusable RAG fragments, factor values, hypotheses, metric selection, and report interpretation. It must not contain a custom Go sweep runner, matrix scheduler, custody importer, timing reducer, or graphing script.

The study compares RAG preparation batching and concurrency while measuring both performance and quality over multiple replicates. It proves that failed provider operations are included, actual usage remains distinct from reservations, runs can resume, and a checked-in JS analysis regenerates the report.

## Historical lesson

The previous TTC path used `cmd/rag-ttc-v3-sweep/`, `internal/workflowv3ttc/`, custom JSON/CSV, a post-hoc Researchctl adapter, and ticket-local Python. One retained run contained 60 successful and four failed generation operations; a published timing boundary omitted 21.472 seconds of failed provider work in one cell. The new acceptance contract is designed to make that class of error impossible by construction.

## Intended source layout

```text
ttc/
  experiment.js
  pipeline.js
  analysis.js
  inputs.json
  profiles.example.yaml
  README.md
  interpretation.md
```

No Go package is permitted unless a genuinely reusable capability is first moved into RAG, Scraper, or Researchctl.

## Experiment script sketch

```javascript
const research = require("researchctl");
const rag = require("rag");
const { ttcPipeline } = require("./pipeline");

const cases = rag.cases(ttcPipeline, {
  chunksPerRequest: [1, 2, 4, 8],
  concurrency: [1, 2, 4],
});

module.exports = research.experimentPlan("ttc-rag-preparation", plan =>
  plan
    .experiment("EXP-TTC-PREPARATION")
    .cases(cases)
    .replicates(3)
    .ordering({ strategy: "randomized", seed: 20260722 })
    .execution({ maxConcurrent: 1, failFast: false })
);
```

Execution concurrency starts at one because the factor under study includes workflow/provider concurrency. Parallel case execution would confound provider load unless a later design deliberately blocks and balances it.

## Pipeline script sketch

```javascript
const rag = require("rag");

module.exports.ttcPipeline = rag.pipeline("ttc-preparation", p => p
  .corpus(rag.inputs.corpus("corpus"))
  .units(rag.units.identity())
  .chunks(rag.chunks.recursive({ maxRunes: 800 }))
  .representations(rag.representations.combined({
    outputs: ["summary", "questions"],
    chunksPerRequest: rag.factor("chunksPerRequest"),
  }))
  .embeddings(rag.embeddings.profile("ttc-embedding"))
  .index("hybrid", rag.indexes.bleveMulti({ lexical: true, vector: true }))
  .query(rag.queries.evaluationDataset())
  .measure(m => m
    .recallAt([5, 10])
    .mrr()
    .ndcgAt([10])
    .latency(["prepare", "query"])
    .tokenUsage()
    .providerCost()
    .failureRates()));
```

The exact fluent API may differ. The acceptance requirement is declarative composition over reusable primitives.

## Required observations

Performance:

- retry-aware workflow elapsed;
- provider operation union and sum durations;
- generation and embedding latency distributions;
- chunk and query throughput;
- operation success/failure/cancellation counts;
- retries;
- actual token usage and cost where available;
- cache hits and provider contacts.

Quality:

- output schema validity;
- summary and question production completeness;
- retrieval recall@k, MRR, and nDCG;
- optional answer correctness, faithfulness, and citation coverage;
- invariants proving factor changes did not silently change corpus or prompts.

Experimental integrity:

- three replicates per case;
- randomized order with recorded seed;
- profile, corpus, dataset, pipeline, prompt, and task-catalog fingerprints;
- explicit missing measurements;
- provider/environment metadata sufficient to interpret drift.

## Analysis script sketch

```javascript
module.exports = analysis.define("ttc-report", ({ runs, stats, charts, report }) => {
  const complete = runs.requireReplicates(3);
  const grouped = complete.groupBy(["chunksPerRequest", "concurrency"]);

  const table = grouped.summarize({
    elapsed: stats.meanCI("workflow.elapsed"),
    retries: stats.sum("workflow.retries"),
    throughput: stats.meanCI("rag.chunks_per_second"),
    quality: stats.meanCI("rag.mrr"),
    usageCoverage: stats.coverage("rag.token_usage"),
  });

  report.table("Primary results", table);
  report.chart("Latency and quality", charts.pareto(table, {
    x: "elapsed.mean", y: "quality.mean", color: "concurrency",
  }));
  report.limitations(complete.missingness());
});
```

## Execution diagram

```text
experiment.js -> Researchctl plan -> randomized case/replicate schedule
                                      |
                                      v
                             Scraper Workflow V3
                                      |
                         RAG preparation/query tasks
                                      |
                           Geppetto external operations
                                      |
                  workflow observations + RAG observations
                                      |
                         Researchctl analysis.js
                                      |
                         tables + charts + report
```

## Decisions

### Decision: replicated, quality-aware acceptance

- **Context:** A single timing replicate cannot distinguish configuration effects from run-order/provider variation.
- **Decision:** At least three replicates and quality metrics are required.
- **Consequences:** The final real-provider run requires explicit cost authorization.
- **Status:** accepted.

### Decision: TTC contains no core lifecycle logic

- **Decision:** Any missing generic feature blocks TTC and is implemented in the owning framework ticket.
- **Rationale:** This makes TTC an architectural fitness function.
- **Status:** accepted.

## Implementation phases

1. Create fixture-backed TTC script and compile it without provider contact.
2. Execute a two-case deterministic smoke through the complete platform.
3. Add profile-backed generation and embedding with operation evidence.
4. Add full factors and three replicates.
5. Run quality and performance analysis.
6. Regenerate report from immutable records in a clean environment.
7. Remove any temporary TTC-specific helper that duplicates framework logic.
8. Conduct an architecture review against every negative acceptance rule.

## Negative acceptance rules

The ticket fails if it introduces:

- a TTC-specific Go runner;
- a TTC matrix-expansion package;
- direct writes to Researchctl SQLite;
- a post-hoc run importer for new runs;
- ticket-local Python/R analysis;
- a second operation ledger;
- a second makespan definition;
- hidden retries;
- inferred actual usage from configured maxima;
- hand-edited report tables.

## Test and review checklist

- Compile and validate scripts.
- Confirm expected case and replicate counts before execution.
- Interrupt and resume midway.
- Inject a provider failure and verify inclusive timing.
- Verify failed operations survive Researchctl publication.
- Compare analysis regeneration digests.
- Inspect all charts at document size.
- Verify single-replicate language never appears.
- Ensure the TTC source remains small and readable by an intern.

## Completion criteria

A clean checkout can execute one command to run or resume the TTC experiment and another to regenerate the report. TTC-specific source is primarily JS and data. Quality and performance share run identity. Failed operations affect standard metrics automatically. No legacy TTC runner remains.

## Technology primer: what the TTC experiment is testing

Batching and concurrency change different parts of execution. Batching changes how many chunks are combined into one generation request. Larger batches may reduce request overhead, but they can increase prompt size, output complexity, validation failures, and retry cost. Concurrency changes how many independent tasks or provider operations can proceed at once. Higher concurrency may reduce wall time until provider limits, local scheduling, rate limits, or resource contention dominate.

The experiment must therefore measure more than makespan. A configuration that finishes quickly but omits summaries, changes retrieval quality, retries expensive batches, or loses usage data is not automatically better.

```text
chunks_per_request affects:
  request count, prompt size, work lost per failed request, output validation

concurrency affects:
  overlap, provider pressure, queueing, rate limits, memory, ordering

both may affect:
  elapsed time, throughput, retries, cost, and downstream RAG quality
```

## Inputs and invariants

The corpus and evaluation dataset are immutable Researchctl input artifacts. The corpus manifest identifies every source record. The evaluation dataset identifies queries and judgments. Every case uses the same digests. Provider profile fingerprints, prompt versions, RAG pipeline digest, Workflow task catalog digest, and analysis source digest are recorded.

Invariants should fail the run or analysis when violated:

- Every expected chunk has the requested representation outputs.
- Embedding cardinality equals representation cardinality.
- Every index entry traces to a source chunk.
- The evaluation query set is unchanged across cases.
- Prompt and model fingerprints are identical unless deliberately factors.
- Requested measures have explicit availability coverage.

## Experimental design

A full factorial design with four batch values and three concurrency values produces twelve cases. Three replicates produce thirty-six runs. Run order is randomized with a recorded seed, but case execution concurrency remains one initially so simultaneous cases do not create uncontrolled provider load.

```text
Factors:
  chunksPerRequest = 1, 2, 4, 8
  concurrency      = 1, 2, 4

Cases: 4 x 3 = 12
Replicates per case: 3
Total runs: 36
```

If provider behavior drifts strongly over time, a blocked design can distribute each replicate round across all cases. Researchctl's ordering policy should record the exact order so time trends can be inspected.

## Reading the existing TTC implementation

`cmd/rag-ttc-v3-sweep/main.go` demonstrates the responsibilities that must disappear from TTC: matrix execution, command flags, provider setup, output files, aggregate metrics, and custody integration. `internal/workflowv3ttc/sweep.go` demonstrates domain-local case expansion. `internal/workflowv3ttc/module.go` contains workflow task behavior that should move into reusable RAG task packages.

The existing `experiments/real-provider-v2/` JS scripts show that much of the desired authoring surface already exists. The new TTC script should build on canonical RAG v2 values and generic Researchctl plans rather than translating the old sweep command line into JavaScript.

## Worked run trace with retry

```text
case batch=4 concurrency=2 replicate=1
  Researchctl creates run R
  Scraper creates workflow W
  preparation partitions 16 chunks into 4 generation batches
  batch 3 provider operation O1 times out after 20 s
  Workflow retries batch 3 as operation O2
  O2 succeeds after 8 s
  all representations and embeddings publish
  retrieval evaluation emits MRR and recall
  Scraper observations include O1 and O2
  Researchctl closes R with succeeded status and retry_count=1
```

The analysis includes the full workflow elapsed value and failed-operation counts. It may additionally show successful-provider-only latency, but that is not substituted for the run boundary.

## Report interpretation rules

The report separates observed facts from conclusions. A table can say that a case had the lowest mean elapsed among completed runs. A recommendation requires replicate completeness, uncertainty, quality equivalence or improvement, acceptable retry/cost behavior, and a clearly bounded provider/environment context.

With only three replicates, intervals may remain wide. The report should state that limitation. It must not call variation among provider calls “benchmark jitter” or treat dozens of within-run calls as dozens of independent runs.

## Expected report figures

1. Run elapsed by batch and concurrency with replicate points and intervals.
2. Throughput and speedup relative to a named baseline.
3. Generation and embedding operation latency distributions, labeled as within-run operations.
4. Retry incidence and failed-operation duration.
5. Quality/performance Pareto plot.
6. Token/cost coverage and values where actual usage exists.
7. Provider-operation timelines for selected diagnostic runs.

Every figure includes denominator, unit, boundary, sample count, and caption. Charts are generated from the checked-in analysis script, not manually edited.

## Operator runbook

Before real-provider execution:

```text
1. Compile and explain the experiment plan.
2. Verify 12 cases and 36 desired runs.
3. Verify corpus, evaluation, prompt, model, and task-catalog digests.
4. Run a fixture-backed two-case smoke.
5. Estimate request/token/cost authority.
6. Obtain explicit authorization.
7. Run with resumable Researchctl custody.
8. Monitor failures without reading secret/provider payloads.
9. Generate analysis only after run completeness review.
```

If interrupted, rerun the same plan. Resume should identify terminal replicates and continue missing work. Do not create another output directory and reconcile files by hand.

## First-week implementation route

The intern first compiles the script against fixture profiles and prints the plan. Next they run two cases with one replicate through the full platform. They inject one deterministic provider failure and inspect operation-inclusive metrics. Then they add quality analysis and verify identical corpus/query digests. Full matrix execution is the last step, not the first.

## Common mistakes

- Running cases concurrently while concurrency is itself a factor confounds provider load.
- Comparing configurations with different prompt/model fingerprints invalidates attribution.
- Declaring fastest from one replicate overstates evidence.
- Ignoring failed operations rewards fragile configurations.
- Treating missing token usage as zero produces false cost rankings.
- Keeping a small TTC helper in Go because it seems convenient starts another bespoke framework.
- Generating figures outside Researchctl loses source provenance.

## Acceptance artifact inventory

A completed study should produce a Researchctl plan artifact, canonical specification per case, run exports for every replicate, selected-attempt records, Workflow V3 observation sets, RAG result/evaluation artifacts, an analysis source manifest, derived tables, chart specifications, rendered figures, and a Markdown/PDF report. These artifacts form one digest-linked tree. The TTC directory itself does not duplicate them in a custom layout.

A reviewer should select any plotted point and follow it backward:

```text
figure point -> table row -> reducer inputs -> run IDs -> selected attempts
             -> workflow observation source digest -> operation records
             -> RAG artifacts -> corpus/evaluation/profile fingerprints
```

If that traversal requires an undocumented script or mutable database query, acceptance fails.

## Review exercise

Before the full run, generate a fixture report and choose one elapsed point, one retry count, and one MRR value. Reconstruct each from raw canonical exports. Then delete derived artifacts and regenerate them. Matching digests prove that publication is derived rather than hand-maintained.

## Intern onboarding checklist

The implementer should explain both factors, calculate case/run counts, locate all immutable input digests, compile the experiment without contact, inject and trace a failed operation, reproduce one analysis table by hand, and show that the TTC directory contains only scripts, data, and interpretation—not lifecycle machinery.

## References

- Program: Researchctl `EXPERIMENT-PLATFORM-CONVERGENCE`.
- Previous audit: `RAG-PIPELINE-CONSOLIDATION-AUDIT`.
- Superseded implementation candidates: `cmd/rag-ttc-v3-sweep/`, `internal/workflowv3ttc/`.
- Depends on every child ticket listed above.
