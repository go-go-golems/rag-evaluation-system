---
Title: Investigation Diary
Ticket: RAG-TTC-V3-SWEEP
Status: active
Topics:
    - rag-eval
    - evaluation
    - workflow
    - chunking
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Chronological implementation and execution evidence for the Workflow V3 Umans batching and concurrency study.
LastUpdated: 2026-07-22T10:15:00-04:00
WhatFor: Preserve commands, decisions, failures, results, and review instructions for reproducing the study.
WhenToUse: When implementing, operating, reviewing, or resuming RAG-TTC-V3-SWEEP.
---

# Diary

## Goal

Capture the evidence-first design, implementation, no-cost control, and authorized real execution of the Workflow V3 Umans chunks-per-request and concurrency sweep.

## Step 1: Create the study ticket and freeze its experimental contract

I created a dedicated docmgr ticket rather than mixing this performance experiment into the broader Workflow V3 completion ticket. The design freezes batch sizes `1, 2, 4, 8`, concurrency limits `1, 2, 4`, an initial 16-chunk/60-request real smoke matrix, precise measurement contracts, privacy boundaries, and separate authorization for any larger run.

Repository inspection confirmed that `ragoperators` already owns deterministic multi-chunk planning and validated provider execution, while the existing Workflow V3 TTC adapter intentionally handles one chunk per task. The design therefore specifies a separate versioned sweep task that reuses the domain operator instead of modifying production task behavior or duplicating prompt construction.

### Prompt Context

**User prompt (verbatim):** "ok, create a new docmgr ticket for that study, create a design/implementation guide, and then run it."

**Assistant interpretation:** Create a structured study workspace, write an implementation-ready guide for the batching/concurrency experiment, implement its Workflow V3 harness and graphs, run the no-cost control, and run the real Umans sample only when its required authority is available.

**Inferred user intent:** Obtain precise, reproducible performance evidence and visualizations for choosing chunks-per-request and Umans concurrency, rather than relying on coarse diagnostic throughput.

### What I did

- Created ticket `RAG-TTC-V3-SWEEP` with topics `rag-eval,evaluation,workflow,chunking`.
- Added the primary design/implementation guide and this investigation diary.
- Added six tasks covering design, harness implementation, graph generation, fixture control, real smoke execution, and publication.
- Inspected `pkg/ragoperators/combined_batch.go`, `internal/workflowv3ttc/provider.go`, and `internal/workflowv3ttc/production_workflow.js` with line numbers.
- Inspected host environment availability without printing secret values.
- Wrote the implementation guide with architecture, exact matrices, measurement schema, decisions, pseudocode, phases, tests, privacy requirements, operator workflow, and hard acceptance criteria.

### Why

- A paid performance study needs an immutable experimental contract and bounded authority before provider calls begin.
- Existing v9 timestamps cannot separate queue wait from provider time, so a new precise measurement path is required.
- Keeping this work in a separate ticket makes its scripts, evidence, graphs, and operator decisions independently reviewable.

### What worked

- `docmgr ticket create-ticket`, both `docmgr doc add` calls, and all task additions succeeded.
- Existing batching code directly supports deterministic chunks-per-request factors.
- Existing Workflow V3 resources, attempts, budgets, and checked usage are suitable foundations.

### What didn't work

- Real execution cannot start in the current shell. The exact availability probe reported:
  - `RAG_GENERATOR_BASE_URL=unset`
  - `RAG_EMBEDDING_BASE_URL=unset`
  - `RAG_RERANKER_BASE_URL=unset`
  - `OPENAI_API_KEY=unset`
  - `ANTHROPIC_API_KEY=unset`
  - `GEMINI_API_KEY=unset`
  - `UMANS_API_KEY=unset`
  - `provider-config=absent`
- No credentials were printed or persisted.

### What I learned

- `CombinedPreparationPlan` already exposes stable batch membership and request text, so the study does not need a new domain batching algorithm.
- The current Workflow V3 provider requires exactly one planned batch from one chunk; changing it would alter production semantics. A dedicated batch contract is safer and more legible.
- The proposed 16-chunk smoke matrix costs exactly 60 generation requests before retries; the 32-chunk main matrix costs 180 requests per replicate.

### What was tricky to build

The main design edge is separating three times that old evidence conflated: durable queue wait, Workflow V3 task overhead, and provider wall time. Workflow V3 provides durable attempt boundaries, but provider start/end must be measured inside the RAG adapter with a monotonic elapsed duration. The design persists only bounded integer microseconds and usage, not request/response payloads. Another edge is avoiding hidden concurrency: every provider request must be one Workflow V3 attempt, and provider-internal generation concurrency must stay at one so the durable resource cap proves the active-call count.

### What warrants a second pair of eyes

- Whether 16 chunks provide enough variety for the smoke matrix without producing misleading batch-size conclusions.
- Whether batch size 8 fits the exact Umans context/rune limits for the frozen source subset.
- Whether connection warm-up should be measured as ordinary evidence or handled by a separately authorized warm-up request.
- The final numeric request/token/cost cap before real execution.

### What should be done in the future

- Implement the dedicated Workflow V3 sweep package, exporter, privacy scanner, and graph renderer.
- Run fixture control and inspect exact request arithmetic and peak concurrency.
- Configure host-local provider authority and obtain explicit numeric budget approval before the real smoke matrix.

### Code review instructions

- Start with the design decisions and experimental contract in `design-doc/01-workflow-v3-umans-batching-and-concurrency-study-design-and-implementation-guide.md`.
- Review existing batching authority in `pkg/ragoperators/combined_batch.go`, then compare the single-chunk adapter in `internal/workflowv3ttc/provider.go`.
- Validate ticket health with `docmgr doctor --ticket RAG-TTC-V3-SWEEP --stale-after 30`.

### Technical details

Commands executed:

```bash
docmgr ticket create-ticket --ticket RAG-TTC-V3-SWEEP \
  --title "Workflow V3 Umans Batching and Concurrency Study" \
  --topics rag-eval,evaluation,workflow,chunking

docmgr doc add --ticket RAG-TTC-V3-SWEEP --doc-type design-doc \
  --title "Workflow V3 Umans Batching and Concurrency Study Design and Implementation Guide"
docmgr doc add --ticket RAG-TTC-V3-SWEEP --doc-type reference \
  --title "Investigation Diary"
```

Exact request arithmetic:

```text
16 chunks: (16/1 + 16/2 + 16/4 + 16/8) × 2 concurrency levels = 60
32 chunks: (32/1 + 32/2 + 32/4 + 32/8) × 3 concurrency levels = 180
```

## Step 2: Implement and execute the no-cost Workflow V3 control matrix

I implemented a dedicated multi-chunk generation task that reuses `ragoperators.PlanCombinedPreparation` and `ExecuteCombinedPreparationBatch`, while preserving the existing single-chunk production task. The adapter now records an RFC3339 provider-start timestamp and monotonic provider elapsed microseconds, reports complete sorted budget usage, and validates every chunk's generated representations before completion.

A new command executed all 12 fixture cells over 16 chunks: batch sizes `1, 2, 4, 8` crossed with concurrency `1, 2, 4`. It performed exactly 90 Workflow V3 generation attempts, observed peak active calls equal to or below every cell's limit, wrote bounded canonical evidence and CSV, removed transient source-bearing runtimes, and rendered four graphs in SVG and PNG.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Implement and execute the safe portion of the requested study now, retaining the real run gate where external authority is unavailable.

**Inferred user intent:** Produce actual precise measurements and graphs using the new Workflow V3 execution path before spending on Umans.

### What I did

- Added `ChunkBatch`, `GeneratedBatch`, and bounded provider timing contracts.
- Added `OperatorProvider.GenerateBatch`, complete generation usage settlement, host module registration, exact bundle task identity, JavaScript task implementation, and a sweep-only Workflow V3 source.
- Added matrix planning with exact request arithmetic, deterministic counterbalancing, a maximum request ceiling, and concurrency rejection above four.
- Added `cmd/rag-ttc-v3-sweep` for isolated cell submission, durable execution, evidence export, and transient runtime cleanup.
- Added `scripts/01-render-sweep.py` and rendered throughput, makespan, and provider-latency graphs as SVG and PNG.
- Ran 12 fixture cells and asserted exactly 90 successful requests and `peakActive <= cell.concurrency <= 4`.

### Why

- The fixture control proves the experiment machinery without provider spend.
- One Workflow V3 attempt per request makes concurrency observable and enforceable.
- Provider timing inside the RAG adapter distinguishes provider wall time from Workflow V3 task makespan.
- Removing transient runtime databases after canonical export keeps ticket evidence compact and avoids checking source-bearing fixture artifacts into documentation.

### What worked

- Final run: `profile=fixtures cells=12 planned_requests=90`.
- Plan digest: `sha256:57d83a39233fb8c061273148ef47476c0c5cce78a4fabacae6bdcbdb8e713a65`.
- Observed peak active requests never exceeded 1, 2, or 4 in the respective cells.
- Canonical evidence is approximately 55 KB; the complete compact evidence and graph directory is approximately 516 KB.
- Focused package and command tests passed, and `git diff --check` passed.

### What didn't work

The first authoring attempt failed exactly with:

```text
execute workflow script: TypeError: output requires a name and workflow value at github.com/go-go-golems/scraper/pkg/gojamodules/workflow.(*planBuilder).object.func8 (native)
exit status 1
```

A map is a set-valued output, so I replaced `plan.output("generated", generated)` with `plan.outputSet("generated", generated)`.

The second attempt failed exactly with:

```text
migrate workflow v3 SQLite: invalid uri authority: ttmp
exit status 1
```

The SQLite opener requires a clean absolute database path. I canonicalized the output root with `filepath.Abs` before constructing per-cell paths.

The next attempt reached execution but failed with:

```text
cell {ChunksPerRequest:1 Concurrency:1 Replicate:1} status failed
```

SQLite attempt evidence showed:

```text
budget / BUDGET_USAGE_INVALID / task usage evidence was invalid
```

The batch adapter reported only nonzero usage plus requests, but Workflow V3 settlement requires every reserved dimension exactly once in sorted order. I added `completeGenerationUsage`, which emits `cost_microunits`, `input_tokens`, `output_tokens`, and `requests` in canonical order, including zeros. This also hardened the existing production adapter.

The first focused lint run then reported exactly:

```text
cmd/rag-ttc-v3-sweep/main.go:295:15: Error return value of `f.Close` is not checked (errcheck)
internal/workflowv3ttc/sweep_test.go:6:1: File is not properly formatted (gofmt)
internal/workflowv3ttc/runtime.go:111:10: ST1005: error strings should not be capitalized (staticcheck)
3 issues:
```

I made CSV flush/close errors explicit, formatted the new test, and lowercased the existing runtime error. The fresh lint run passed with `0 issues.`

Initial visual QA found confusing scientific notation on the makespan axis and overlapping provider-latency series with inconsistent legend wording/colors. I changed makespan to milliseconds, applied a stable palette, standardized `concurrency=N`, and added slight x offsets. The first follow-up vision request failed externally with `tls: bad record MAC`; a retry succeeded and found no publication-blocking defects.

### What I learned

- Workflow V3 correctly rejects partial usage evidence, which caught a production-adapter edge not exercised by the synthetic provider tests.
- A 16-chunk full fixture matrix has 90 requests, while the proposed real smoke matrix omits concurrency 2 and therefore has 60.
- At batch size 8 with only two requests, observed peak concurrency is necessarily two even when capacity is four; configured capacity and achievable parallelism must be plotted separately.

### What was tricky to build

The task output and attempt timeline are separate evidence streams. Provider spans are naturally associated with deterministic batch keys in generated output artifacts, while durable attempt intervals are associated with hashed node keys. The exporter retains both bounded arrays without inventing a false association. Aggregate distributions and active-count timelines remain exact. A future schema version may persist the batch key as bounded attempt metadata if per-attempt joins are required.

The first compact export left per-cell SQLite and content-addressed source artifacts in the ticket, totaling 4.8 MB and 227 files. Although the source was nonsensitive fixture text, that is the wrong evidence boundary. I changed the command to export attempt intervals and provider measurements first, close each store, and delete the transient cell runtime. The final ticket keeps only canonical evidence, CSV, and graphs.

### What warrants a second pair of eyes

- Validate whether the real provider reports all token dimensions for combined preparation.
- Review the conservative per-request budget maxima before real execution; actual smoke limits should be derived from the real profile and approved numerically.
- Confirm that the real batch-size-8 source subset remains under both rune and model context limits.
- Consider persisting a bounded batch-key correlation field in generic Workflow V3 attempt metadata rather than relying on two ordered evidence collections.

### What should be done in the future

- Add a host-only real profile to the command without accepting secret values as CLI arguments.
- Add preflight-only cost estimation and hard token/cost caps.
- Execute the 60-request smoke matrix after credentials and numeric authority are available.
- Extend graphs with real token rate, cost efficiency, percentiles, and generation/embedding overlap.

### Code review instructions

- Start at `internal/workflowv3ttc/sweep.go` for matrix invariants.
- Review `OperatorProvider.GenerateBatch` and `completeGenerationUsage` in `internal/workflowv3ttc/provider.go`.
- Review exact task/bundle wiring in `bundle.go`, `module.go`, `tasks.cjs`, and `sweep_workflow.js`.
- Run:

```bash
GOWORK=off go test ./internal/workflowv3ttc ./cmd/rag-ttc-v3-sweep -count=1
GOWORK=off go run ./cmd/rag-ttc-v3-sweep --output /tmp/rag-ttc-v3-sweep-fixture --chunks 16 --maximum-requests 90 --concurrency 1,2,4
python3 ttmp/2026/07/22/RAG-TTC-V3-SWEEP--workflow-v3-umans-batching-and-concurrency-study/scripts/01-render-sweep.py /tmp/rag-ttc-v3-sweep-fixture/evidence.json --output /tmp/rag-ttc-v3-sweep-fixture/graphs
```

### Technical details

Observed fixture control summary:

```text
batch=1 concurrency=1 requests=16 peak=1
batch=2 concurrency=1 requests=8  peak=1
batch=4 concurrency=1 requests=4  peak=1
batch=8 concurrency=1 requests=2  peak=1
batch=1 concurrency=2 requests=16 peak=2
batch=2 concurrency=2 requests=8  peak=2
batch=4 concurrency=2 requests=4  peak=2
batch=8 concurrency=2 requests=2  peak=2
batch=1 concurrency=4 requests=16 peak=4
batch=2 concurrency=4 requests=8  peak=4
batch=4 concurrency=4 requests=4  peak=4
batch=8 concurrency=4 requests=2  peak=2
```

## Step 3: Add real-profile preflight and expand canonical graphs

I extended the runner so the same Workflow V3 sweep can load an exact real provider profile and canonical researchctl execution specification without accepting secret values on the command line. The real path resolves profile and model digests, reconstructs the frozen source chunks from the verified corpus artifact, disables response-cache reuse across cells, pins provider-internal concurrency to one, and requires explicit execution plus request/token/cost ceilings before submission.

I reran the complete fixture control and expanded canonical evidence with settled usage dimensions and durable attempt intervals. The renderer now produces token-rate, cost-efficiency, and observed-concurrency timeline views in addition to the original plots. Because the fixture provider does not report tokens or cost, those panels explicitly display unavailable data rather than misleading zero-valued series.

### Prompt Context

**User prompt (verbatim):** "run and make graphs"

**Assistant interpretation:** Continue beyond the study design: execute all currently authorized cells, render precise graphs, and prepare the exact gated real-provider path rather than stopping at a plan.

**Inferred user intent:** Obtain actual Workflow V3 performance evidence and readable visualizations, ultimately from the bounded Umans smoke matrix.

### What I did

- Added real-provider profile/specification/corpus loading to `cmd/rag-ttc-v3-sweep`.
- Added explicit `--execute-real`, request, input-token, output-token, and cost authority checks.
- Added exact provider-profile and generation-model digest reporting to real dry-run output.
- Exported settled Workflow V3 usage and attempt intervals before deleting transient cell databases.
- Added token-rate, cost-efficiency, active-attempt timeline, and generation/embedding-overlap graphs.
- Added a downstream measured embedding task with an independent capacity of four, precise provider spans, representation/request counts, transactional token/request accounting, and no persisted vectors in study evidence.
- Re-ran the 12-cell fixture matrix: 90 generation requests and 192 embedding provider requests.
- Ran focused tests and lint after the changes.

### Why

- A real run must not silently fall back to fixtures or accept ambient unbounded spend.
- Reconstructing chunks from the canonical execution specification and digest-verified artifact keeps the real sample tied to the frozen TTC study.
- Missing fixture usage is not zero usage; explicit N/A plots prevent false conclusions.

### What worked

- Fixture rerun completed with `profile=fixtures cells=12 planned_requests=90`.
- The renderer produced 16 SVG/PNG image files plus its manifest.
- Every cell settled `embedding.requests` exactly to its measured provider-request count.
- Focused tests passed and lint reported `0 issues.`.
- Real-profile wiring compiles and reaches host-only provider resolution without making calls.

### What didn't work

The first real dry-run with the intended non-example path failed exactly with:

```text
RAG_PROVIDER_CONFIG_READ: open /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/experiments/real-provider-v2/provider-config.yaml: no such file or directory
exit status 1
```

Using the checked-in non-secret example configuration reached environment resolution and failed exactly with:

```text
RAG_PROVIDER_ENV_MISSING
exit status 1
```

No provider request was submitted.

Visual QA correctly rejected the initial zero token/cost lines as misleading and found the overlaid attempt bars insufficiently explicit about concurrency. I changed missing usage panels to centered `Data unavailable` annotations and changed the timeline to active-attempt step series with a visible hard-cap line. Follow-up QA required the zero-overlap panel to say explicitly that all fixture series coincide and that the result is not indicative of real-provider performance; the renderer now does so.

The first embedding-enabled run found zero generation/embedding overlap in all cells. Inspection of durable attempt timestamps proved embedding attempts began after the last generation attempt, despite independent dispatcher capacities. Reducing the generation map page size to one did not change that evidence. I retained the honest zero-overlap graph rather than fabricating overlap; resolving the chained-map publication behavior remains required before making a work-conserving real-run claim.

### What I learned

- The real provider configuration loader is the correct security boundary: it resolves environment references without exposing secret values.
- The current conservative task reservation implies an intentionally high preflight ceiling; it must be tightened from exact provider/profile constraints before asking for monetary authorization.
- Active-count series communicate enforcement more directly than overlapping Gantt bars for this matrix.

### What was tricky to build

The canonical researchctl specification stores the execution in `canonicalIdentity.domainConfig` and source artifact identity separately. The loader must verify both artifact size and SHA-256 before decoding the corpus, then execute only provider-free static pipeline stages to reproduce ordered chunks. It must not reuse diagnostic v9 source databases. The implementation follows that path and selects the first 16 chunks only after deterministic ID sorting.

### What warrants a second pair of eyes

- Review the real profile's per-request maximum token and cost reservations before authorization.
- Verify the chosen 16 chunks are representative enough for a smoke test while remaining deterministic.
- Confirm the provider profile's configured concurrency is at least four; Workflow V3 may request four but must respect any lower host cap.

### What should be done in the future

- Supply a host-only real provider configuration and its referenced environment values.
- Complete a credential-resolving dry-run and reduce conservative budget maxima to exact approved limits.
- Investigate why this sweep's chained embedding map starts only after generation completion and either fix the generic activation path or document a deliberately sequential study design before real execution.
- Run the 60-generation/128-embedding-request real matrix and render real token/cost graphs.

### Code review instructions

- Review `cmd/rag-ttc-v3-sweep/profile.go` for profile, specification, digest, and source boundaries.
- Review `main.go` for numeric execution gates and usage export.
- Review `scripts/01-render-sweep.py` for N/A semantics and active-attempt computation.
- Validate with:

```bash
GOWORK=off go test ./cmd/rag-ttc-v3-sweep ./internal/workflowv3ttc -count=1
GOWORK=off golangci-lint run ./cmd/rag-ttc-v3-sweep/... ./internal/workflowv3ttc/...
```

### Technical details

The bounded real smoke command shape is:

```bash
rag-ttc-v3-sweep \
  --profile real \
  --provider-config /host-only/provider-config.yaml \
  --specification /path/to/canonical-specification.json \
  --artifact-root /path/to/artifacts \
  --chunks 16 \
  --concurrency 1,4 \
  --maximum-requests 60
```

The first invocation is dry-run only. Actual submission additionally requires `--execute-real` and explicit numeric token/cost ceilings.

## Step 4: Validate and commit the fixture study tranche

Fresh focused tests, focused race tests, lint, command build, Python compilation, artifact cardinality/accounting checks, privacy scans, graph rendering, visual QA, docmgr doctor, and `git diff --check` passed. The repository pre-commit hook also passed Biome checks, repository lint/vet, and its Go test suite. I committed the implementation, compact fixture evidence, SVG graphs, design, and diary as `d68fe9a` (`feat(workflowv3): add TTC measurement sweep`). I removed an accidentally generated Python bytecode file before finalizing the commit and added a ticket-local `.gitignore` for `__pycache__/` and `*.pyc`.

The goal remains incomplete: the real provider environment is absent, no numeric real spend has been authorized, and no real Umans call has been submitted.

## Step 5: Complete canonical JSONL and expanded CSV custody

**Commit (code):** `15fb9ea` — "workflowv3: add canonical TTC measurement custody"

The active goal explicitly requires canonical JSON/JSONL and CSV. The first committed tranche had canonical aggregate JSON and CSV but no line-oriented evidence stream. I added deterministic `measurements.jsonl` export with one canonical record per generation attempt, embedding attempt, and provider batch measurement. The 12-cell fixture run emits exactly 270 lines: 90 generation attempts, 90 embedding task attempts, and 90 paired generation/embedding provider measurements.

I also expanded `cells.csv` with generation and embedding request counts, overlap, token dimensions, and cost microunits. Aggregate evidence now pins the Workflow V3 plan digest, registry generation, bundle digest, provider-profile digest, and generation-model digest. I reran the fixture study and regenerated all graphs from the same fresh evidence file so tables and rendered views cannot drift.

## Step 6: Resolve real Umans authority without making inference calls

I found the existing host-local Umans engine profile and secure Pi authentication entry, then constructed an untracked temporary provider configuration that combines `umans-flash` generation with the local `nomic-embed-text` service. The configuration references the existing Pinocchio profile stack; it does not copy or print the Umans credential. A non-inference `/v1/models/info` request confirmed that `umans-flash` is reachable and reports a 262,144-token context and completion ceiling. Ollama confirmed `nomic-embed-text:latest` is present.

The exact real Workflow V3 dry-run now succeeds through provider/profile/model resolution and digest-verified frozen-corpus materialization. It resolves provider profile digest `sha256:a90b04ff11b40cc83e8a79fa93b207bdea254822031fb9a369d6d5cc70a15086`, model digest `sha256:de4fd2ae4ab9b359e5ce7fee7b7773258f6eb9643cd61a682f54d76b7584a684`, 16 frozen chunks, 60 generation requests, and 128 embedding requests. No inference call was made.

The remaining gate is explicit numeric spend authority. Current deliberately conservative reservation maxima are 9,600,000 cost microunits, 983,040 input tokens, 983,040 output tokens, and 3,932,160 embedding tokens. These are hard reservation ceilings, not expected consumption. The real command still refuses submission without `--execute-real` and matching numeric authority.

## Step 7: Resolve chained-map overlap before real execution

I traced the zero-overlap result to the sweep composition rather than generic Workflow V3 activation. `ChainedExpansionCandidate` correctly exposes each contiguous successful upstream prefix, but the downstream `measure-embedding-batches` map requested pages of 16. With only 2–16 generated batches per cell, this deliberately withheld downstream materialization until all 16 were available or the upstream map became final. The earlier experiment changed the upstream generation page size, which could not remove that downstream threshold.

I changed only the downstream map page size from 16 to 1. The regenerated 12-cell fixture control now records generation/embedding overlap in 9 of 12 cells; the three zero cells remain honest timing outcomes for the extremely short fixture calls. Request arithmetic remains exactly 90 generation and 192 embedding requests, and peak active generation remains within each cell's admitted limit and the hard cap of four. The workflow plan digest changed to `sha256:9efd478fc91aef90acf30e4160542b1eebc492d245fb200fcdd8f9c5004a636b` as required for this semantic change.

Visual QA also found that the concurrency timeline showed only generation. I changed it to display solid generation and dotted embedding activity for every batch size, labeled the cap explicitly as the Umans generation cap, clamped elapsed time at zero, and kept the fixture-control qualification. The canonical JSON/JSONL/CSV and all eight SVG/PNG pairs were regenerated from the corrected run.

## Step 8: Attempt an explicit operator approval prompt

To avoid treating a generic continuation as spend authority, I attempted to collect the exact numeric authorization through the local `plz-confirm` operator gate. The prompt named all limits: 60 Umans generation requests, 128 local embedding requests, four active Umans calls, 983,040 input tokens, 983,040 output tokens, 3,932,160 embedding tokens, and 9,600,000 cost microunits ($9.60 if denominated in USD).

The approval service was not running, and the command failed before presenting any prompt:

```text
Error: create confirm request: post /api/requests: Post "http://localhost:3000/api/requests": dial tcp 127.0.0.1:3000: connect: connection refused
```

This is not approval. No provider inference request was made. The remaining valid path is an explicit conversational approval of the same numeric ceilings (or revised lower ceilings).

## Step 9: Start the operator service and receive a negative spend decision

I started the local `plz-confirm` backend on `127.0.0.1:3000` and repeated the exact approval prompt from Step 8. The structured response was `approved=false`. I treated that as a rejection, not as missing or ambiguous authority. No provider inference request was made.

The real matrix must remain blocked unless the operator later supplies a new explicit affirmative approval for the stated ceilings or replacement lower ceilings. The active goal therefore remains incomplete by design.

## Step 10: Replace the conservative cost ceiling with exact published pricing

### Prompt Context

**User prompt (verbatim):** (same active-goal continuation as Step 1)

**Assistant interpretation:** Continue with the next no-cost action that can reduce the authorization blocker without making real provider calls.

**Inferred user intent:** Finish the exact, reproducible study while preserving the explicit spend gate and avoiding repeated fixture work.

**Commit (code):** `3dbb512` — "workflowv3: pin exact Umans sweep authority"

### What I did

I investigated whether the rejected $9.60 ceiling could be tightened without making an inference call. The public Umans tariff states that Flash costs USD $0.15 per million input tokens, $1.00 per million output tokens, and $0.05 per million cache-read tokens, with no per-call minimum. The pinned host policy now caps each response at 8,192 tokens and records all tariff dimensions as integer microunits per million tokens. Cache writes use the conservative input rate of $0.15/M.

I added a strict host-side generation policy to `ragproviders`. It overlays a profile-resolved generator with an exact response-token cap and integer pricing, includes those non-secret values in the effective provider identity and settings fingerprint, and rejects incomplete policies. The sweep command now rejects any real authority whose effective concurrency, response cap, or tariff differs from the pinned study policy.

I reduced the per-attempt reservation from 160,000 to 10,650 cost microunits and from 16,384 to 8,192 output tokens. The cost is checked integer ceiling arithmetic: 16,384 input tokens at $0.15/M plus 8,192 output tokens at $1.00/M. The exact 60-request matrix ceiling is therefore 639,000 microunits (USD $0.639), 983,040 input tokens, and 491,520 output tokens. The embedding ceiling remains 3,932,160 tokens and 128 actual requests.

The exact real dry-run succeeded first with an untracked host configuration and then with the checked-in non-secret `provider-config.umans-flash.example.yaml`. The reproducible checked-in configuration has effective provider profile digest `sha256:cf092e0f1a389169a22e519beb8b81b7f059dd6bb66ec391ea934bf857e46702`; it verified all 16 frozen chunks and made no inference call. I regenerated the fixture control under workflow plan digest `sha256:eeada085d929a376eb7ab67abee1bc6da5d9e2a939bbab4a79bafee6526a7a17` and bundle digest `sha256:45fc68911c33c29b399e65ac18a893381b86df40b9d3a1cbd7275893678bee62`.

I presented a second structured approval prompt for the materially lower $0.639 ceiling. The response was again `approved=false`. I treated that as an explicit rejection and made no provider inference calls.

### Why

- Replace a generic conservative reservation with an exact profile- and tariff-backed ceiling before requesting authority again.
- Ensure real cost evidence is computed from explicit rates rather than treating unknown price as zero.

### What worked

- Strict policy validation, identity custody, exact dry-run, fixture regeneration, and focused tests all succeeded.
- The maximum monetary authority fell from $9.60 to $0.639 without narrowing the required 60-request matrix.

### What I learned

The original $9.60 figure was a generic task maximum rather than an exact profile-derived bound. Exact integer tariff custody both improves safety and makes provider cost available for real evidence instead of silently treating unknown price as zero.

### What was tricky to build

A cost rate of zero and an unknown cost rate are semantically different. The policy requires explicit positive input/output rates and carries them into Geppetto `ModelInfo`; this lets provider usage produce a real cost while preserving checked integer microunit admission. Applying the policy before fingerprinting was essential—the prior default response cap was assigned after identity recording.

### What didn't work

A temporary cross-registry profile overlay failed because an engine profile stack cannot reference a profile in a different YAML registry:

```text
RAG_PROVIDER_PROFILE_CHAIN: engine profile YAML registry validation failed: validation error (registry.profiles[ttc-umans-flash-priced].stack[0]): referenced profile "umans-flash" not found in registry "rag-ttc-v3-pricing"
```

I did not copy credentials into another registry. Instead, the host provider config now applies the non-secret execution policy after resolving the existing secure profile.

Two fresh visual-QA calls failed with the same external TLS error:

```text
pinocchio exited with code 1: inference failed: remote error: tls: bad record MAC
```

The renderer itself is unchanged and regenerated all eight SVG/PNG pairs; this external inspection failure must remain recorded rather than being called a pass.

### What warrants a second pair of eyes

- Confirm that cache-write tokens should retain the conservative standard-input rate until Umans publishes a distinct cache-write tariff.
- Review the new provider-policy identity fields and fail-closed real-authority comparison.
- Re-run visual QA when the external image-analysis service is healthy.

### What should be done in the future

Stop. The reduced exact spend request was explicitly rejected. Resume only after a new affirmative numeric authorization supersedes that decision.

### Code review instructions

Review `pkg/ragproviders/config.go` and `provider_set.go` for strict policy validation and non-secret identity custody; review `internal/workflowv3ttc/sweep.go`, `bundle.go`, and `sweep_workflow.js` for exact matching reservations; review `cmd/rag-ttc-v3-sweep/profile.go` and `main.go` for fail-closed authority and ceiling checks.

Validate with focused provider/sweep tests, the exact real dry-run without `--execute-real`, fixture custody assertions, Python compilation, `docmgr doctor`, and `git diff --check`.

### Technical details

- Input tariff: 150,000 microunits per million tokens.
- Output tariff: 1,000,000 microunits per million tokens.
- Cache-read tariff: 50,000 microunits per million tokens.
- Conservative cache-write tariff: 150,000 microunits per million tokens.
- Per request: ceil(16,384 × 0.15) + ceil(8,192 × 1.00) = 10,650 microunits.
- Matrix: 60 × 10,650 = 639,000 microunits.

## Step 11: Separate lease-attempt concurrency from provider-call concurrency

### Prompt Context

**User prompt (verbatim):** (same active-goal continuation as Step 1)

**Assistant interpretation:** Validate the regenerated fixture custody after pricing changes and fix any evidence defect rather than masking it.

**Inferred user intent:** Preserve precise, auditable distinctions between Workflow V3 control-plane timing and actual provider activity.

**Commit (code):** `3dbb512` — "workflowv3: pin exact Umans sweep authority"

### What I did

The fixture accounting assertion found one concurrency-1 cell with `peakActive=2`. Inspection showed that `peakActive` and `overlapMicros` were derived from lease-attempt timestamps, even though their labels and hard-cap graph implied provider calls. Attempt finalization can overlap the next admitted task's control-plane work without overlapping provider calls, so this was a semantic evidence defect rather than a provider-cap violation.

I versioned canonical sweep evidence to v2 and split the fields into `attemptPeakActive`, `providerPeakActive`, `attemptOverlapMicros`, and `providerOverlapMicros`. Provider values are now computed exclusively from the monotonic provider start plus elapsed spans retained in each measured batch. The CSV uses equally explicit columns. The graph renderer now consumes provider overlap and provider spans, labels elapsed time from the first provider span, and describes active provider spans rather than active attempts. Unit tests pin touching-interval ordering and provider-wall overlap arithmetic.

### Why

The objective requires both exact queue/attempt evidence and exact provider-wall/concurrency evidence. Conflating them could falsely report a hard-cap violation or claim generation/embedding overlap where only control-plane finalization overlapped.

### What worked

The regenerated fixture matrix reports provider peaks within every cell's configured limit: all concurrency-1 cells peak at one, concurrency-2 cells peak at two or less, and concurrency-4 cells peak at four or less. Provider-wall generation/embedding overlap is observed in seven of twelve cells, while attempt overlap is separately observed in nine of twelve.

Fresh visual QA accepted the provider overlap and provider concurrency timeline. It rejected the initial missing-usage wording because it said `Data unavailable` rather than the required explicit `N/A`; I changed both token-rate and cost-efficiency panels to `N/A`, regenerated them, and follow-up visual QA found no defects.

### What didn't work

The initial fixture validation failed on the old ambiguous assertion because a concurrency-1 cell had attempt peak two. Treating that value as provider concurrency would have been incorrect. I changed the evidence model rather than loosening the assertion.

### What I learned

Workflow attempt spans are control-plane evidence; provider spans are data-plane evidence. Both are useful, but only provider spans prove the maximum number of active Umans calls.

### What was tricky to build

Events ending and starting at the same timestamp must process the end first or the sweep would invent overlap. The shared interval reducer sorts negative deltas before positive deltas at equal timestamps and has a dedicated regression test.

### What warrants a second pair of eyes

- Verify that embedding batch spans are correctly described as provider-wall envelopes around their sequential per-representation requests.
- Confirm every report and graph uses provider fields for provider claims and attempt fields only for Workflow V3 control-plane analysis.

### What should be done in the future

Use the v2 evidence schema for the real run. Do not execute it until a new affirmative spend decision supersedes the rejection.

### Code review instructions

Review `cmd/rag-ttc-v3-sweep/main.go` and `main_test.go` for interval semantics and canonical fields, then review `scripts/01-render-sweep.py` for provider-only graph inputs and explicit N/A panels.

### Technical details

- Fixture provider peak: never above the cell limit or four.
- Cells with provider overlap: 7/12.
- Cells with attempt overlap: 9/12.
- Provider timeline: solid generation spans, dotted embedding span envelopes.
- Evidence schema: `rag-ttc-v3-sweep-evidence/v2`.

## Step 12: Correct the operator-decision record from rejection to timeout

### Prompt Context

**User prompt (verbatim):** (same active-goal continuation as Step 1)

**Assistant interpretation:** Reinspect the remaining authorization blocker and take the next evidence-backed action without making provider calls.

**Inferred user intent:** Ensure the spend gate and its audit trail are accurate before any real execution.

### What I did

I reopened both structured `plz-confirm` response artifacts rather than relying only on the parsed boolean. Each response contained `approved=false` together with `comment: AUTO_TIMEOUT`. The first $9.60 request and the later $0.639 request therefore expired after the server-side default timeout; neither was an explicit rejection.

I corrected the current design document and appended this correction to the diary and changelog. Earlier diary entries remain as chronological evidence of the mistaken interpretation, but this step supersedes their conclusion. No provider inference call was made.

### Why

An approval timeout and a rejection have different operator semantics. The real run remains blocked in both cases, but claiming an explicit negative decision when no decision occurred would make the custody record inaccurate.

### What worked

The raw response artifacts retained enough structured metadata to resolve the ambiguity exactly: `approved=false` plus `comment=AUTO_TIMEOUT`.

### What didn't work

The earlier parser selected only the boolean `approved` field and discarded `comment`, causing the timeout to be misclassified as a rejection. The `plz-confirm` command defaults both request expiration and wait timeout to 300 seconds unless explicitly overridden.

### What I learned

Operator-gate consumers must preserve and branch on the full terminal reason, not collapse approval, rejection, timeout, cancellation, and transport failure into one boolean.

### What was tricky to build

No code was needed, but the correction had to preserve append-only diary history while updating the current design statement. Rewriting prior diary entries would have erased the investigation path.

### What warrants a second pair of eyes

- Confirm future approval wrappers treat `AUTO_TIMEOUT` as missing authority rather than rejection.
- Ensure the real command remains gated until an affirmative response or explicit conversational approval is received.

### What should be done in the future

Request explicit conversational approval, or run `plz-confirm` with both `--timeout` and `--wait-timeout` long enough for an attended operator session. Do not make real calls meanwhile.

### Code review instructions

Review the raw response shape and this correction against Steps 9 and 10. No runtime code changed in this step.

### Technical details

- Original request result: `approved=false`, `comment=AUTO_TIMEOUT`.
- Reduced request result: `approved=false`, `comment=AUTO_TIMEOUT`.
- Current blocker: missing affirmative authority, not an explicit rejection.

## Step 13: Attempt an attended approval and stop without a decision

### Prompt Context

**User prompt (verbatim):** (same active-goal continuation as Step 1)

**Assistant interpretation:** Try one attended operator path for the exact $0.639 authority, then obey the objective's stop condition if no decision is available.

**Inferred user intent:** Advance the real study only through explicit operator authority and leave no ambiguous or dangling control process.

### What I did

I started the local `plz-confirm` backend, opened its Agent UI in a browser, and submitted a one-hour confirmation request with both `--timeout 3600` and `--wait-timeout 3600`. The visible request listed all exact ceilings and separate approve/reject buttons. Repeated status checks showed the request remained pending and its output artifact remained empty.

Because no operator decision was available, I terminated the pending confirmation process, stopped the local approval backend, and closed the browser. No provider inference call was made.

### Why

The prior five-minute requests had expired before an attended decision. This attempt verified that the complete UI path worked and that the exact prompt was visible, while preserving the rule that only a human affirmative response can authorize spend.

### What worked

- The backend, WebSocket, browser UI, request rendering, exact ceilings, and one-hour expiration all worked.
- The process remained safely blocked while approval was pending.
- Cleanup left no unattended approval process capable of completing later.

### What didn't work

No operator selected approve or reject during the attended window available to this session. The attempt produced no decision artifact.

### What I learned

A visible and technically functioning approval UI is not authorization. Pending state must remain fail-closed, and unattended requests should be canceled when the agent stops.

### What was tricky to build

The output file exists at request start because of shell redirection but remains zero bytes until completion. Status checks therefore had to combine file size with process liveness rather than treating file existence as a terminal response.

### What warrants a second pair of eyes

- Confirm the operator can use either an explicit conversational response or a separately attended Agent UI session next time.
- Ensure no future automation interprets generic goal continuation as spend authority.

### What should be done in the future

Stop until an operator explicitly authorizes or declines the exact bounded matrix. Do not launch real calls while authority is absent.

### Code review instructions

No repository runtime code changed. Review the chronological approval evidence in Steps 8 through 13 and the current blocker statement.

### Technical details

- Attended request: one-hour server and client wait timeouts.
- Terminal state: canceled locally while pending; no approval/rejection payload.
- Provider requests made: zero.

## Step 12: Abort the first authorized run at the cell timeout and preserve cumulative authority

The user explicitly authorized the bounded real matrix by replying `Ok, run it.` in direct response to the exact $0.639 ceiling and request/token/concurrency limits. I submitted the run with all matching command-line gates. The first concurrency-1, one-chunk batch attempt was admitted, but the command's fixture-oriented 30-second cell deadline expired before a provider response arrived. The command canceled the dispatcher and exited with the cell still recorded as running.

I did not retry. I inspected only compact control-plane columns in the transient database: one generation attempt was running, fifteen nodes were pending, one request and its token/cost maxima remained reserved, and no provider usage or response had been persisted. I exported a bounded non-sensitive aborted-attempt record, then deleted the transient SQLite/WAL/artifact tree because it contained source-bearing inputs. The submitted request is counted conservatively even though no response or usage was observed.

### Prompt Context

**User prompt (verbatim):** "Ok, run it."

**Assistant interpretation:** Execute the exact real 16-chunk matrix under the immediately preceding $0.639, 60-generation-request, 128-embedding-request, token, and concurrency ceilings.

**Inferred user intent:** Proceed from validated preflight to real provider evidence without weakening any privacy or budget guard.

**Commit (code):** `59097f6` — "workflowv3: preserve cumulative sweep authority"

### What I did

- Submitted the exact authorized command with `--execute-real` and every numeric gate.
- Inspected compact node, attempt, event, and budget columns after failure without reading source artifacts or provider bodies.
- Published `sources/real-attempt-001/evidence.json` and removed the source-bearing transient runtime.
- Added a real-only 30-minute default cell deadline; fixtures retain 30 seconds.
- Added a durable, fsync-before-call generation admission ledger with a cumulative request ceiling and prior-request offset.
- Made real output directories fail closed instead of deleting existing evidence.
- Corrected real run IDs so they use the `real-` prefix rather than `fixture-`.
- Added concurrency and persistence tests for the authority ledger and wired admission immediately before generation provider execution.

### Why

The first admitted request consumed request authority even without a response. A fresh 60-request matrix would therefore require a cumulative ceiling of 61, not a silent reset to 60. Persisting admission before each call prevents concurrent workers from exceeding the approved cumulative request count.

### What worked

- Cancellation stopped the command; no second provider request was submitted.
- Compact evidence and budget reservations made the consumed authority auditable.
- The sensitive transient tree was removed after bounded export.
- Focused tests, race tests, lint, and diff checks pass after the fix.

### What didn't work

The command failed exactly with:

```text
cell {ChunksPerRequest:1 Concurrency:1 Replicate:1} status running
exit status 1
```

The root cause was a hard-coded 30-second deadline shared with fast fixtures. It was not a provider-schema or semantic workflow failure. The old timeout branch also reported the last snapshot as merely `running`; the corrected branch now reports the configured duration and dispatcher result.

### What I learned

Authorization is cumulative across aborted invocations, not just successful evidence exports. A request admitted before cancellation must remain charged against future authority even when provider usage is unknown.

### What was tricky to build

The authority ledger must reserve before the provider call, serialize concurrent workers, atomically rename a fully synced file, and fsync the containing directory. Persisting after the call would leave a crash window where a submitted request was not counted. Existing real output must not be automatically removed because it may be the only record of consumed authority.

### What warrants a second pair of eyes

- Confirm the conservative decision to count the canceled request at its full token/cost reservation.
- Review the 30-minute per-cell default against expected Umans latency; the dispatcher still exits immediately on terminal workflow failure.
- Verify final custody includes the aborted-attempt record and cumulative generation-authority ledger.

### What should be done in the future

Obtain new explicit authority for one additional generation request and its reservation. The cumulative ceilings are 61 generation requests, 649,650 cost microunits (USD $0.64965), 999,424 input tokens, 499,712 output tokens, 128 embedding requests, 3,932,160 embedding tokens, and four concurrent Umans calls. Do not retry before approval.

### Code review instructions

Review `cmd/rag-ttc-v3-sweep/admission.go` first, then the real-only output/timeout path in `main.go`, admission wiring in `profile.go` and `internal/workflowv3ttc/provider.go`, and the compact aborted evidence JSON.

Validate with:

```text
GOWORK=off go test -race ./cmd/rag-ttc-v3-sweep ./internal/workflowv3ttc -count=1
GOWORK=off golangci-lint run ./cmd/rag-ttc-v3-sweep/... ./internal/workflowv3ttc/...
```

### Technical details

- Prior admitted requests: 1.
- New matrix requests still required: 60.
- Required cumulative generation ceiling: 61.
- Incremental maximum: 10,650 cost microunits, 16,384 input tokens, 8,192 output tokens.
- No provider response, usage, body, prompt, source text, or vector is present in the compact aborted evidence.

## Step 13: Exhaust cumulative authority after one real retry and preserve partial custody

The user replied `do it` directly to the cumulative 61-request ceiling. I ran the corrected 30-minute-timeout command with one prior request and exactly 60 newly available admissions. Seven cells completed. The final concurrency-4, one-chunk cell then failed because one earlier malformed/transient generation retry had consumed the only request that the final node needed. The durable authority ledger reached exactly 61 and correctly denied the next provider call before submission.

I stopped without extending authority. The failed cell had twelve succeeded generation nodes, three generation calls running when cancellation propagated, and one generation node failed with `RAG_TTC_GENERATION_REQUEST_CEILING`. The current implementation had already deleted each successful cell runtime after accumulating its evidence only in memory, so process exit lost the seven completed cell measurements. I recorded that loss explicitly rather than reconstructing or fabricating results, exported compact attempt custody, and deleted the remaining source-bearing runtime.

### Prompt Context

**User prompt (verbatim):** "do it"

**Assistant interpretation:** Execute the corrected real matrix under the immediately preceding cumulative 61-request authority.

**Inferred user intent:** Complete the real study now that the first aborted request is included in the cumulative budget.

**Commit (code):** `072ddd4` — "workflowv3: checkpoint partial real sweep evidence"

### What I did

- Executed with cumulative maximum 61 and prior count one.
- Verified the durable ledger stopped exactly at 61 admitted generation requests.
- Inspected only bounded control-plane status, failure-code, and budget columns.
- Added `sources/real-attempt-002/evidence.json`; removed all transient SQLite and source artifacts.
- Changed successful cells to atomically publish standalone canonical checkpoints before transient deletion.
- Added automatic privacy-safe failed-cell checkpoints and a test proving failure messages/provider bodies are excluded.
- Added a deferred real-runtime scrub for every normal return path.
- Added explicit retry-headroom authority; retries now contribute to cumulative request/token/cost maxima.
- Made final aggregate evidence atomic and include the final generation-authority state.

### Why

Real provider output validation can legitimately consume retry debt. A 60-node matrix with an exact 60-call ceiling cannot finish after even one malformed response. Retry headroom must be explicitly authorized and globally admitted, not hidden inside per-cell Workflow budgets.

### What worked

- The cumulative limiter denied request 62 before a provider call.
- Seven cells reached Workflow V3 success before the final cell failure.
- The terminal failure code was compact and non-sensitive.
- Privacy cleanup removed all transient databases and source-bearing artifacts.
- Focused tests, race tests, lint, a 12-cell checkpoint control, and diff checks pass.

### What didn't work

The run failed exactly with:

```text
cell {ChunksPerRequest:1 Concurrency:4 Replicate:1} status failed
exit status 1
```

The failed node recorded `budget / RAG_TTC_GENERATION_REQUEST_CEILING / retryable=false`. Complete measurements from the seven successful cells were unavailable after process exit because publication occurred only after the whole matrix. That custody defect is fixed for future runs, but the missing evidence cannot be recovered.

### What I learned

A matrix-wide aggregate written only at the end is not durable evidence. Each successful cell must publish independently before its source-bearing runtime is erased. Also, planned requests and authorized retry requests are different quantities and both must appear in the cumulative authority.

### What was tricky to build

The automatic failure checkpoint must retain failure codes and budget amounts while excluding arbitrary failure messages, which may contain provider bodies. Runtime cleanup must run after checkpoint publication on normal error paths without deleting durable cell evidence or the cumulative authority ledger.

### What warrants a second pair of eyes

- Review the proposed eight-request retry headroom against observed and diagnostic malformed-output rates.
- Verify successful cell checkpoints contain all fields needed to reconstruct final JSON/JSONL/CSV after a later cell failure.
- Confirm runtime cleanup behavior under SIGKILL remains an operational recovery concern; normal returns are scrubbed automatically.

### What should be done in the future

Obtain new cumulative authority before rerunning. With 61 requests already consumed, a fresh 60-request matrix plus eight retry requests requires a cumulative ceiling of 129 generation requests, 1,373,850 cost microunits (USD $1.37385), 2,113,536 input tokens, and 1,056,768 output tokens. Embedding ceilings remain 128 requests and 3,932,160 tokens; concurrency remains four. Do not retry without approval.

### Code review instructions

Review the per-cell checkpoint ordering and runtime-root defer in `cmd/rag-ttc-v3-sweep/main.go`, atomic writer in `admission.go`, and failure-redaction regression in `main_test.go`. Inspect both compact aborted-attempt evidence files and verify no transient databases remain.

### Technical details

- Cumulative admitted generation requests: 61/61.
- New calls in attempt 2: 60.
- Successful cells before failure: 7/8.
- Earlier-cell retry inferred from admission arithmetic: 1.
- Proposed fresh matrix allowance: 60 planned + 8 retries.
- Proposed new cumulative ceiling: 61 prior + 68 new = 129.

## Step 11: Per-cell durable operation custody

The sweep now exports the generic Workflow V3 operation ledger before closing and deleting each successful cell runtime. Each cell checkpoint and the aggregate evidence refer to the relative JSONL and manifest paths plus the canonical manifest. Timeout and terminal-failure paths use the same export before writing their privacy-safe failed-cell checkpoint. Failure custody no longer silently drops a budget, operation export, or checkpoint write error.

Fixture validation produced 12 cells and 282 durable provider-operation records. Every cell manifest and JSONL path existed under the sweep output. The default 30-second control also passed on a repeat. An earlier timeout was therefore retained as a transient investigation record rather than “fixed” by increasing the application default; the five-minute timeout was used only as an explicit diagnostic bound during one repeat.

**Validation:**

```text
GOWORK=off go test ./cmd/rag-ttc-v3-sweep ./internal/workflowv3ttc -count=1
GOWORK=off go run ./cmd/rag-ttc-v3-sweep --profile fixtures --chunks 16 --concurrency 1,2,4 --maximum-requests 90 --output /tmp/rag-ledger-default-check
GOWORK=off golangci-lint run ./cmd/rag-ttc-v3-sweep/... ./internal/workflowv3ttc/...
```

The fixture manifest count was 282, and all 12 referenced JSONL and manifest files existed.

## Step 12: Deterministic failed-cell operation reductions

Failed-cell checkpoints now include a closed, deterministic reduction of the exported operation rows: admitted/completed/incomplete counts, outcome counts, total provider elapsed microseconds, peak provider concurrency, generation/embedding overlap microseconds, and generation/embedding operation counts. The reduction reads no provider payload, free-form error message, URL, header, or credential.

A forced `--cell-timeout 1ns` fixture failure verified that the failure path writes an operation JSONL, a manifest, and an `operationReduction` object before exiting non-zero. A unit test verifies interval overlap and incomplete-operation accounting.

## Step 13: Generic researchctl operation-custody export builder

Added `researchctladapter.BuildOperationCustodyRunExport`. The RAG-owned adapter converts compact operation JSONL/manifest files and scalar counters into a strict `lab.RunExport`. It derives artifacts' digest and size from local files but persists only relative URIs and verified identities; timestamps and run identity remain explicit inputs. The builder rejects absolute/traversal URIs, duplicate artifacts/metrics, missing source identity, and non-terminal status. A unit test verifies the resulting bundle using researchctl's public artifact verifier.

This is intentionally an adapter primitive; the sweep CLI still needs an operator-facing export/import command that supplies a canonical specification and explicit custody identity.

## Step 14: Sweep-to-researchctl verified import

The sweep now optionally writes `researchctl-run-export.json` at its output root when all explicit custody identity flags are provided. It includes the aggregate evidence, each cell checkpoint, every operation JSONL and manifest, and four scalar-only metrics. It requires the export's canonical researchctl specification and explicit run/attempt/external-run/timestamp values; it does not generate identities or timestamps.

A 12-cell fixture sweep produced the bundle. `researchctl lab init` followed by `researchctl experiment import-run` staged and atomically imported 37 verified artifacts and four metrics into a fresh laboratory.

## Step 15: Fixture graph visual inspection

Visually inspected four rendered fixture figures: makespan, generation/embedding overlap, provider latency, and request timeline. Axes, markers, legends, and reference lines were legible; no clipping, blank regions, or sensitive content appeared. The vision review suggested a possible spelling issue in the overlap title, but the renderer source has the correct `embedding` spelling, so no speculative rendering change was made. These figures are suitable as fixture evidence only, not as real-provider conclusions.
