---
Title: Investigation diary
Ticket: RAG-PIPELINE-CONSOLIDATION-AUDIT
Status: active
Topics:
    - rag-eval
    - evaluation
    - workflow
    - ttc
    - research
    - visualization
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://ttmp/2026/07/22/RAG-PIPELINE-CONSOLIDATION-AUDIT--assess-improve-and-consolidate-the-rag-evaluation-pipeline/analysis/01-ttc-real-run-performance-audit-and-pipeline-consolidation-assessment.md
      Note: Primary repaired deliverable
    - Path: repo://ttmp/2026/07/22/RAG-PIPELINE-CONSOLIDATION-AUDIT--assess-improve-and-consolidate-the-rag-evaluation-pipeline/design-doc/01-rag-pipeline-assessment-consolidation-and-improvement-plan.md
      Note: Consolidation design produced from the audit
    - Path: repo://ttmp/2026/07/22/RAG-PIPELINE-CONSOLIDATION-AUDIT--assess-improve-and-consolidate-the-rag-evaluation-pipeline/scripts/01-analyze-ttc-real-run.py
      Note: Implementation and visual QA chronology
ExternalSources: []
Summary: Chronological evidence, failures, decisions, and validation for the pipeline audit and consolidation work.
LastUpdated: 2026-07-22T18:32:55.201786557-04:00
WhatFor: Preserve a reproducible account of why the prior TTC result package was inadequate and how the pipeline is assessed and repaired.
WhenToUse: When reviewing, continuing, or validating the RAG pipeline assessment and consolidation work.
---


# Investigation diary

## Goal

Assess the previous TTC measurement work from transcript through published artifacts, replace unsupported or shallow conclusions with reproducible analysis, and produce a concrete plan to consolidate measurement, execution, custody, reporting, and publication across the RAG pipeline.

## Step 1: Establish the audit from transcript and repository evidence

I converted the exact 47 MB Pi session requested by the user into a normalized minitrace archive and queried its user turns, ticket operations, repository activity, and final publication sequence. This made it possible to compare what was requested—precise timings, graphs, a detailed report, and a retraceable diary—with what was actually committed and uploaded.

The first repository audit confirmed the user's concern. The final report was only 3.6 KB, contained one eight-row makespan table, did not embed or interpret the graph set, omitted latency distributions, retries, speedups, utilization, idle-time decomposition, jitter, and key measurement limitations, and left the ticket's publication task unchecked. The ticket index still said real-provider execution remained gated even though a real run had completed.

### Prompt Context

**User prompt (verbatim):**

> Look at the latest diary / ticket, to me it feels like half assed work and someone losing the plot. I got a report for graphs and timing that basically is empty, i'm quite disappointed, see what you can figure out and do.
>
> Use go-minitrace to analyse this transcripts: File: /home/manuel/.pi/agent/sessions/--home-manuel-workspaces-2026-06-30-benchmark-cpu-inference--/2026-07-19T00-24-40-023Z_019f77c2-c157-7f3b-8424-24adfb4af971.jsonl
>
> It was across multiple directories and worktrees, I think.
>
> Store scripts you write in the scripts/ folder of a new ticket that is about assessing and improving and consolidating the whole pipeline
>
> in rag-evaluation-system i guess

**Assistant interpretation:** Audit the prior agent session and its cross-repository output, create a new rag-evaluation-system ticket for the broader pipeline assessment, preserve every new script there, and actively repair the weak timing/reporting deliverable rather than merely criticizing it.

**Inferred user intent:** Recover a trustworthy technical result from work that was declared complete too early, identify why the process lost focus, and leave behind a consolidated, reproducible pipeline improvement package.

### What I did

- Read the go-minitrace transcript-analysis and transcript-friction skills, including schema, attribution, query, and interpretation references.
- Checked installed `go-minitrace convert pi`, `query run`, and schema help before using flags.
- Converted only the requested native session into `analysis/session-019f77c2-audit/archives/` without modifying the native JSONL.
- Saved source provenance, conversion output, validation output, SQL, and query results under the workspace audit directory.
- Queried 165 user turns plus session metadata and reconstructed the final real-run, graph, report, commit, and upload sequence.
- Verified external repository state and commits in `rag-evaluation-system`.
- Read the prior report, task list, ticket index, graph renderer, aggregate evidence, cell CSV, per-cell evidence, and operation ledgers.
- Created ticket `RAG-PIPELINE-CONSOLIDATION-AUDIT` with design, diary, and six explicit tasks.

### Why

- The transcript alone can show intent and commands, but repository state is required to verify claims and completion.
- A dedicated ticket keeps assessment scripts and derived evidence separate from the already-published run artifacts while preserving links to the original ticket.
- The new scope is broader than one report: the failure involved execution, instrumentation semantics, analysis, bookkeeping, and publication quality.

### What worked

- Conversion produced one quality-A Pi archive with 6,943 turns and 6,597 tool calls.
- The final transcript window precisely identified the relevant repository, ticket, run root, commits, and reMarkable path.
- Operation ledgers contain enough timestamps and outcomes to recover a substantially richer analysis without another paid run.
- The real run has 60 successful generation operations, four failed generation operations, and 128 successful embedding operations; those retry facts were omitted from the published report but are recoverable.

### What didn't work

- The initial workspace command combined `git status` with doc discovery at a non-repository parent and stopped at:

```text
fatal: not a git repository (or any of the parent directories): .git
```

- Two first-pass exploratory Python snippets assumed the wrong nesting in per-cell JSON and failed with:

```text
KeyError: 'chunksPerRequest'
```

and then:

```text
KeyError: 'operations'
```

  Reading one complete cell record established that the payload is nested under root `cell`, with parameters under `cell.cell` and operation paths under `cell.operations`; the corrected query then produced per-cell latency and coverage summaries.

### What I learned

- The prior session did substantial implementation and custody work, but the final reporting phase collapsed into a terse status note and a shallow 3.6 KB document.
- The report's statement of “exactly 60 successful generation attempts” is narrowly true but incomplete: four additional failed provider operations consumed 100.045 seconds of provider wall time before retries succeeded.
- The report did not update the original ticket's final task or stale index, so its “complete” narrative disagrees with docmgr state.
- The existing graph renderer still contains fixture-specific labels for token and cost figures when rendering a real run, and its latency graph uses one median point per cell rather than a distribution.

### What was tricky to build

- The session crossed the original benchmark workspace, the separate `rag-eval-ttc` worktree, `researchctl`, `scraper`, the Obsidian vault, and remote Mac services. Recorded cwd therefore could not identify the relevant repository by itself; exact tool-call paths and external Git verification were required.
- Aggregate evidence mixes planned generation requests, actual operation admissions, successful outputs, and usage counters. These must remain separate: 60 planned/successful outputs, 64 generation admissions, and four failed operations are all simultaneously correct.
- Provider-union coverage can exceed the measured cell span if failed operations began before the earliest retained successful-attempt boundary. The analysis must use a declared boundary and avoid presenting negative “idle” time as meaningful.

### What warrants a second pair of eyes

- Verify the intended semantics of `makespanMicros` and whether failed retry intervals should be included in the canonical cell boundary.
- Confirm whether generation provider concurrency is allowed to overlap timed-out/failed calls after cancellation; operation-ledger peaks and successful-batch peaks answer different questions.
- Review any derived utilization metric before treating it as scheduler utilization rather than provider-span occupancy.

### What should be done in the future

- Make publication completion an auditable gate: checked task, current index, report with graph inventory and limitations, `docmgr doctor`, upload verification, and diary step must agree.
- Preserve explicit planned/admitted/succeeded/failed request counters in the canonical aggregate schema so reports do not need to reconstruct retries from JSONL.
- Move robust latency, retry, occupancy, and boundary analysis into one maintained command rather than one-off plotting code.

### Code review instructions

- Start with the old report at `ttmp/2026/07/22/RAG-TTC-V3-SWEEP--workflow-v3-umans-batching-and-concurrency-study/analysis/02-authorized-real-umans-qualification-results.md`.
- Compare it with `sources/real-attempt-003/cells.csv`, `cells/*.json`, and `operations/*.jsonl`.
- Reproduce transcript findings from the saved SQL and JSON under `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/analysis/session-019f77c2-audit/`.

### Technical details

- Native session: `019f77c2-c157-7f3b-8424-24adfb4af971`.
- Relevant user request turn: 5233 created the study ticket; turn 5328 requested “run and make graphs”; turns 6406–6414 demanded performance insight and better persistent instrumentation; turns 6932–6942 created and uploaded the final shallow report.
- Verified final report commit: `7e7b41ba564ce1afd0d885239ea71b2a14078fef`.
- Real run root: `sources/real-attempt-003/` in the original `RAG-TTC-V3-SWEEP` ticket.

## Step 2: Reconstruct the performance result and design the consolidated pipeline

I implemented a reproducible ticket-local reducer over the retained real operation ledgers, generated corrected datasets and five graph families, and replaced the terse original result with a 36 KB analysis that embeds and interprets every figure. The key correction is an operation-inclusive timing boundary: one retry cell's published makespan omitted 21.472 seconds because the aggregate used successful attempts only.

I also wrote the consolidation design for turning this prototype into maintained RAG analysis, graph/report generation, researchctl metrics, and a fail-closed publication receipt. The report was opened through `md-view view`; browser inspection confirmed all five images loaded at full dimensions with no console errors.

### Prompt Context

**User prompt (verbatim):**

> Once you get all the graphs, and create a markdown doc that showcases the whle analysis and the graphs, open with md-view view ...
>
> you can use read tool to look at the images

**Assistant interpretation:** Finish the full visual analysis, create one Markdown showcase that embeds the complete result, open it in md-view, and directly inspect rendered images rather than assuming successful file generation means publication quality.

**Inferred user intent:** Receive a readable, visually verified technical deliverable rather than another list of artifact paths or an uninspected report.

**Commit (code):** `1207fb2128376778f9828dfff31e79c1e6ccae35` — "docs: audit TTC performance pipeline and rebuild report"; PNG graph assets were added in `e43578c8ed1e1d63e9a429d5d4e9f9e81e88a301` — "docs: add TTC audit graph renders".

### What I did

- Added `scripts/01-analyze-ttc-real-run.py` under this ticket as required.
- Parsed all eight cell checkpoints and all 192 operation JSONL rows.
- Derived explicit planned/admitted/succeeded/failed request counts, operation-inclusive elapsed boundaries, generation/embedding distributions, queue delay, occupancy, provider coverage, retries, and execution order.
- Wrote `summary.json`, `cells-derived.csv`, and `operations-derived.csv` under `sources/derived-real-attempt-003/`.
- Rendered five graph families in PNG and SVG plus a manifest.
- Ran visual QA on all five figures. The first pass found ambiguous stacked semantics, an unexplained mean marker, mixed denominators, clipped retry points, compressed near-100% comparisons, and inconsistent timeline axes.
- Revised the renderer to use grouped elapsed bars, labeled means/failures, separate occupancy/coverage/retry panels, data labels, and common timeline axes.
- Ran a second image QA pass, which found no publication blocker.
- Opened all five final PNG files directly with the read tool. This exposed one remaining overlap between provider-coverage labels and the subplot title; I moved near-100% labels inside bars and regenerated the figures.
- Wrote the full analysis report and the architecture/consolidation design.
- Ran `md-view view` and inspected the rendered page with Playwright. All five images completed with nonzero natural dimensions; there were no browser console errors.

### Why

- The existing report did not answer the user's actual questions about latency, underutilization, idle time, jitter, throughput, retries, or graph interpretation.
- Operation ledgers are the only evidence path that includes failed calls and therefore the correct source for retry-safe analysis.
- Direct image and rendered-page inspection catches defects that file existence and Python success cannot.
- The broader design addresses why the reporting failure happened rather than preserving another one-off script as permanent infrastructure.

### What worked

- The reducer reconciled 60 planned/successful generation calls, 64 admissions, four failures, and 128 embedding calls.
- It recovered 100.045 seconds of failed generation provider work.
- It proved that batch-2/concurrency-1's published 157.050-second makespan excluded 21.472 seconds; operation-inclusive elapsed is 178.522 seconds.
- It showed provider coverage of 98.99%–99.94%, refuting the hypothesis of broad scheduler pauses in completed cells.
- It identified batch-8/concurrency-2 as the fastest observed cell at 34.389 seconds and 0.4653 chunks/s while keeping the single-replicate limitation explicit.
- `md-view` loaded every image and all expected report headings.
- `docmgr doctor` passed cleanly for both the new audit ticket and repaired original sweep ticket.
- The reMarkable dry-run and real bundle upload succeeded; cloud listing verified `RAG Pipeline TTC Performance Audit and Consolidation Plan` under `/ai/2026/07/22/RAG-PIPELINE-CONSOLIDATION-AUDIT`.

### What didn't work

- First-round graphs were technically generated but not publication-safe. The image reviewer reported:
  - unclear stacked elapsed semantics;
  - unexplained mean triangles;
  - failed-call marker overlap;
  - conflated time-coverage and admission-failure denominators;
  - near-100% values visually compressed;
  - inconsistent timeline axes.
- After the second QA pass said the figures were publication-ready, direct read-tool inspection still found provider-coverage value labels overlapping the subplot title. The labels were moved inside the bars and the graph set was regenerated.
- Generation token/s and cost efficiency remain unavailable because most operations lack actual provider usage counters. The report explicitly refuses to infer them from reservation ceilings.
- The first commit attempt failed `git diff --cached --check` because Matplotlib emits multiline SVG path data with trailing spaces. I added deterministic SVG line normalization to the analysis script and regenerated all figures.
- The second commit attempt then found one extra final blank line in the docmgr-generated changelog and CRLF CSV output. I removed the changelog blank line, set the CSV writer's `lineterminator="\n"`, regenerated outputs, and the next diff check and commit passed.
- Repository `.gitignore` excludes `*.png`, so the first successful commit contained SVG but not the PNG files referenced by Markdown. I verified the ignore rule and force-added exactly the five reviewed PNG outputs in a separate focused commit.

### What I learned

- The current operation ledger is stronger than the aggregate evidence. The expensive run did preserve the evidence needed for a serious report; the previous finalization simply did not reduce it.
- `readCell` filters to successful attempts before deriving makespan. Operation-inclusive boundaries must become first-class aggregate fields.
- Provider call spans cover almost all elapsed time. Slow runtime is generation-provider latency and retry behavior, not large dispatcher gaps.
- Visual QA needs both an image-focused review and direct inspection in the actual Markdown renderer.

### What was tricky to build

- The operation-inclusive boundary must include admission timestamps as well as provider spans and attempt timestamps. Using provider union alone can still omit durable pre-provider queue time.
- Occupancy, provider coverage, and retry incidence have different denominators. Combining them on one unlabeled axis created a misleading graph even though each number was correct.
- A common x-axis makes timeline duration comparable, but it leaves blank right-hand regions for fast cells. That blank area is intentional comparative scale, not idle time; the report explains it.
- A single replicate permits descriptive distributions within a cell but not stable confidence intervals across configurations.

### What warrants a second pair of eyes

- Review whether operation admission or outer command invocation should be the canonical start boundary for future production evidence.
- Review the batch-8/concurrency-2 recommendation against output-quality evidence before adopting it globally.
- Confirm whether transport failure codes should gain bounded subcodes without weakening privacy.
- Review the proposed ownership of canonical analysis in rag-evaluation-system versus a generic researchctl reporting layer.

### What should be done in the future

- Port the prototype reducer into a tested maintained RAG analysis package and command.
- Add explicit usage availability and planned/admitted/succeeded/failed/incomplete counts to the next evidence schema.
- Add replicated randomized studies and quality joins before claiming a universal optimum.
- Add an executable publication receipt that blocks stale tickets and incomplete reports.

### Code review instructions

- Start with `analysis/01-ttc-real-run-performance-audit-and-pipeline-consolidation-assessment.md` and follow each graph to `sources/derived-real-attempt-003/graphs/`.
- Review the reducer's `load_cells`, `interval_union`, boundary calculation, and graph generation in `scripts/01-analyze-ttc-real-run.py`.
- Compare derived rows to original `RAG-TTC-V3-SWEEP/sources/real-attempt-003/operations/*.jsonl`.
- Re-run the reproduction command in the report, then open it with `md-view view`.

### Technical details

- Final report title: `TTC real-run performance audit and pipeline consolidation assessment`.
- md-view URL: `http://localhost:43883/render?file=/home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/22/RAG-PIPELINE-CONSOLIDATION-AUDIT--assess-improve-and-consolidate-the-rag-evaluation-pipeline/analysis/01-ttc-real-run-performance-audit-and-pipeline-consolidation-assessment.md`.
- Rendered images: five PNGs, all complete; natural dimensions ranged from 1980×990 to 2160×2700, with the three-panel figure at 2880×900.
- Browser console: zero errors and zero warnings.
