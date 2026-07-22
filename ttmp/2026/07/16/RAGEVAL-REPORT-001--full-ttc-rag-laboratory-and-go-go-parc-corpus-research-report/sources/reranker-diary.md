---
Title: Implementation diary
Ticket: RAGEVAL-RERANK-001
Status: active
Topics:
    - rag
    - reranking
    - ttc
    - geppetto
    - ollama
    - rag-eval
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-15T16:55:30.081643451-04:00
WhatFor: ""
WhenToUse: ""
---

# Implementation diary

## Goal

Capture the evidence, decisions, experiments, and implementation steps for a
cross-encoder reranking stage in the immutable TTC RAG laboratory.

## Step 1: Establish the reranking boundary

The existing laboratory is a working raw-retrieval baseline, not a theoretical
prototype. A live JavaScript/Geppetto run completed 20 TTC evaluation cards
with weighted RRF and persisted every trace. The reranker may therefore be
evaluated as a narrow change to candidate ordering against known baseline
quality, latency, and storage behavior.

The selected first runtime is `llama-server` rather than Ollama. The Mac
already has BGE and Qwen reranker artifacts in Ollama storage, but the
laboratory needs a true score-per-query-document service. llama.cpp documents
that contract explicitly; Ollama's normal local API does not provide it.

### Prompt Context

**User prompt (verbatim):** "ok, continue with your tasks. and once you get to the reranker, create a new docmgr ticket and Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Finish the raw RAG work, then open a separate,
documented reranker workstream before implementation.

**Inferred user intent:** Make reranking an evidence-backed, reproducible
laboratory capability rather than an opaque model call.

### What I did

- Inspected the executor, run persistence, representation guards, and the live
  Mac model inventory.
- Read the current official llama.cpp server documentation for the native
  reranking endpoint and its required server flags.
- Created this ticket, its detailed design guide, granular tasks, and this
  diary before any reranker code is written.

### Why

Cross-encoder reranking changes rank semantics and latency. It needs a named
runtime and a persisted candidate-scoring record so an evaluator can distinguish
retrieval quality from reranker quality.

### What worked

- The baseline executor has a single clear insertion point after RRF fusion.
- The Mac inventory includes `qllama/bge-reranker-v2-m3:q4_k_m` and Qwen3
  reranker 4B/8B artifacts, so an initial local comparison is feasible.

### What didn't work

- No failure occurred in this research step. The prior tunnel listener was
  initially inaccessible from the sandbox, but operator-level tmux inspection
  confirmed the existing private tunnel was healthy.

### What I learned

llama.cpp exposes `/reranking` with aliases including `/v1/rerank`; it requires
a reranker model plus `--embedding --pooling rank` and accepts `query`,
`documents`, and `top_n`.

### What was tricky to build

The word “reranker” can refer either to a model artifact or to a scoring API.
The design treats the HTTP scoring contract as the integration boundary, which
prevents model-store details from leaking into immutable experiment identity.

### What warrants a second pair of eyes

- The candidate text budget and whether to collapse duplicate parent chunks
  before scoring materially affect both cost and relevance.

### What should be done in the future

- Execute Task 4's probe before adding Go code; capture the real response
  schema from the selected llama.cpp build.

### Code review instructions

- Start with the executor's fusion path in `pkg/raglab/executor.go`.
- Compare the proposed request contract with the official llama.cpp endpoint
  in the primary design document.

### Technical details

```text
BM25 + vector -> RRF candidates (N) -> cross encoder -> reordered K -> collapse -> citations -> metrics
```

## Step 2: Validate and publish the implementation package

The ticket now has a complete implementation foundation before any reranker
code is written. Its task list differentiates completed research decisions
from the first probe, Go contracts, executor integration, UI work, and the
TTC comparison matrix. This ordering ensures the response schema is observed
before a client is committed.

The document package passed docmgr validation after adding the new `reranking`
and `ollama` vocabulary terms. It was bundled with the ticket overview,
technical design, and diary and uploaded to the reMarkable destination below.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Store the design in a real ticket and deliver it
to the tablet, not merely report a recommendation in chat.

**Inferred user intent:** Give an intern a standalone, reviewable starting
point and let the project owner read it away from the workstation.

**Commit (docs):** `788c867` — "docs(ticket): design immutable RAG reranking stage"

### What I did

- Ran `docmgr --root ttmp doctor --ticket RAGEVAL-RERANK-001 --stale-after 30`.
- Added the `reranking` and `ollama` topic vocabulary entries required by the
  three ticket documents.
- Uploaded the bundle containing `index.md`, the primary design guide, and
  this diary with a two-level table of contents.

### Why

The ticket must be usable independently of this conversation. Validation and a
single tablet bundle make the durable documentation the handoff artifact.

### What worked

- Docmgr reported: `✅ All checks passed`.
- `remarquee` reported:
  `OK: uploaded RAGEVAL Reranking Stage Guide.pdf -> /ai/2026/07/15/RAGEVAL-RERANK-001`.

### What didn't work

- The first ticket scaffold was accidentally created using `docmgr --root .`,
  which placed it in a repository-level untracked `2026/` directory. It was
  recreated correctly under `ttmp/` before substantive docs were written; the
  empty mistaken scaffold was removed without touching project content.

### What I learned

This repository's active docmgr root is `ttmp/`, even though `--root .` is a
valid CLI argument. Explicitly naming the doc root prevents an untracked
parallel ticket tree.

### What was tricky to build

The design must distinguish evidence already observed from implementation
proposals. The guide labels the llama.cpp HTTP boundary as accepted, while the
exact decoder and collapse placement remain task-gated until the live probe.

### What warrants a second pair of eyes

- Before implementing Task 6, review the saved live probe response against the
  adapter's assumed index/score JSON structure.

### What should be done in the future

- Execute Task 4 in a dedicated tmux session and store its curl command and
  output artifact under `scripts/`.

### Code review instructions

- Read the primary design guide's Sections 4–9.
- Confirm the ticket task ordering retains the real-probe gate before adapter
  implementation.

### Technical details

```text
ticket docs -- docmgr doctor --> validated
           -- remarquee bundle --> /ai/2026/07/15/RAGEVAL-RERANK-001
```

## Step 3: Probe the real llama.cpp BGE reranker contract

The first reranker task was intentionally a probe rather than application code.
The Mac already contained the BGE reranker as an Ollama model blob and the
Ollama application bundles a compatible `llama-server` binary. This made it
possible to establish the HTTP contract without installing another package or
downloading a second model artifact.

The server runs only on the Mac loopback interface at port 8012. A dedicated
local tmux SSH tunnel forwards `127.0.0.1:18012` to that address. The bounded
three-document request confirms the exact result structure, negative score
range, `top_n` truncation behavior, and input-index mapping the Go adapter must
validate.

### Prompt Context

**User prompt (verbatim):** "go ahead."

**Assistant interpretation:** Continue from the reranker design into the next
task-gated implementation step without waiting for a separate confirmation.

**Inferred user intent:** Replace assumptions about the local reranking stack
with measured behavior before committing the adapter architecture.

### What I did

- Located `/Applications/Ollama.app/Contents/Resources/llama-server` on
  `mimimi-2.local`; its version is `1 (cb295bf59)` for Darwin arm64.
- Read the BGE Ollama manifest and used its local model blob with
  `--embedding --pooling rank --rerank --host 127.0.0.1 --port 8012`.
- Started the server using `nohup` because the Mac does not have `tmux`;
  verified `/health` returns `{"status":"ok"}`.
- Created local tmux session `rag-reranker-mimimi` for the private tunnel
  `127.0.0.1:18012 -> mimimi-2.local:127.0.0.1:8012`.
- Added the executable probe script and the observed contract/results record.
- Re-ran the script successfully: the complete request took `0.087942` seconds
  through the tunnel.

### Why

The adapter must decode the actual server, not a generic interpretation of an
OpenAI-style reranking API. Indexes refer to request document positions and
scores are negative finite values, both of which influence correctness checks.

### What worked

- `/v1/rerank` returned `{model, object, usage, results}`.
- Results contain `index` and `relevance_score`; the payroll candidate at index
  zero had the highest score.
- `top_n: 2` returned exactly two result rows, establishing that a complete
  rerank pass must request `top_n == len(documents)`.

### What didn't work

- The initial check did not find `llama-server` on `PATH`; the executable is
  bundled inside `Ollama.app`. The model is also stored as an Ollama blob rather
  than a named GGUF file. Both were resolved by using their explicit paths.

### What I learned

The server evaluates document pairs on separate slots. The API response in this
probe was score ordered, but request-document identity is carried only by the
returned index. A correct adapter maps index to candidate first and sorts only
after validation.

### What was tricky to build

The established local tunnel playbook uses tmux, but tmux is unavailable on the
Mac. The Mac server is therefore a bounded `nohup` process with an explicit log
at `/tmp/rag-reranker-llama-server.log`; the workstation tunnel remains tmux
managed and inspectable. This preserves the private binding without pretending
that remote process lifecycle is equivalent to the local tunnel lifecycle.

### What warrants a second pair of eyes

- Confirm that the adapter treats a short `top_n` response as invalid only for
  its complete-scoring call; a future API may intentionally support partial
  reranking with a separately declared policy.

### What should be done in the future

- Implement Task 5's pure contracts and fixtures from this recorded response.

### Code review instructions

- Run `scripts/01-probe-llamacpp-bge-reranker.sh` after checking the tunnel.
- Compare its output to `scripts/02-llamacpp-bge-reranker-probe-results.md`.

### Technical details

```text
documents[i] -> llama.cpp result.index i -> immutable candidate ID
                         + relevance_score (finite, descending is better)
```

## Step 4: Add immutable reranking policy and transport-neutral contracts

The first code change adds the vocabulary required to express a reranking
experiment without making a network call. `RerankingSpec` is part of the
retrieval policy and therefore participates in the canonical experiment
fingerprint. The `Reranker` interface and request/result types are pure Go;
they deliberately do not name llama.cpp, HTTP, Goja, or a database.

This isolates two classes of correctness. Builder validation establishes that
an experiment declares a coherent model and bounded candidate policy. A later
adapter establishes that a configured runtime can fulfill that policy. The
executor remains unchanged in this step, so no current raw retrieval behavior
is affected.

### Prompt Context

**User prompt (verbatim):** (same as Step 3)

**Assistant interpretation:** Continue from the observed reranking contract
into the next implementation task without expanding the scope prematurely.

**Inferred user intent:** Build the reranker as a typed, reproducible system
whose experiment identity remains meaningful across provider restarts.

**Commit (code):** `3764e20` — "feat: add immutable reranking policy types"

### What I did

- Added `RerankingKind`, `RerankingSpec`, and optional
  `RetrievalPlan.Reranking` in `pkg/raglab/types.go`.
- Added transport-neutral `RerankCandidate`, `RerankRequest`, `RerankResult`,
  `RerankerIdentity`, and `Reranker` interface in `pkg/raglab/reranker.go`.
- Added `RetrievalBuilder.RerankCrossEncoder(model, candidateCount, results)`.
- Added validation for kind, nonempty model, positive bounded counts, final
  result coverage, duplicate configuration, and canonical copying.
- Added tests proving stable fingerprints for the same policy and different
  fingerprints when the reranker model changes.

### Why

The model and candidate budget change the experiment. They must be serialized
before a runner is implemented; otherwise later runs could use different
rerankers while appearing to share an experiment specification.

### What worked

- `GOWORK=off go test ./pkg/raglab -count=1` passed.
- The new policy is absent by default, so existing raw plans retain their
  pre-reranking canonical form.

### What didn't work

- No implementation failure occurred in this step.

### What I learned

The explicit `results` value is a pre-collapse reranking window, while the
existing retrieval `results` value remains the final user-facing result count.
Validation requires the former to cover the latter. The future executor may
request scores for all candidates, then apply this declared window before
collapse.

### What was tricky to build

The experiment needs model identity but must not persist the endpoint. The
types keep `Model` in `RerankingSpec`, while URL, timeout, and credentials will
be adapter options. This mirrors the existing explicit query-embedding design.

## Step 5: Implement and validate the llama.cpp reranker adapter

The concrete adapter is intentionally small and strict. It posts the observed
`model`, `query`, `documents`, and `top_n` request shape to `/v1/rerank`, then
uses each returned `index` to recover the submitted candidate's durable ID.
The implementation does not trust server result order: it validates a complete
set of distinct in-range indexes, rejects missing or non-finite scores, and
sorts results by descending relevance score with deterministic retrieval-rank
tie breaking.

### Prompt Context

**User prompt (verbatim):** “go ahead.”

**Assistant interpretation:** Continue the already approved reranker task
sequence without pausing at the adapter boundary.

**Inferred user intent:** Produce a usable, reproducible runtime component,
not merely a wire-format sketch.

**Commit (code):** pending — adapter and tests are staged as this diary entry
is written.

### What I did

- Added `LlamaCPPReranker` and explicit `LlamaCPPRerankerOptions` in
  `pkg/raglab/reranker_llamacpp.go`.
- Required an operator-supplied base URL and model; kept them out of the
  immutable policy except for the policy's model identity.
- Enforced a configurable serialized-request limit and propagated the caller's
  context through `http.NewRequestWithContext`.
- Added strict response validation for result count, mandatory fields, index
  range/uniqueness, and finite scores.
- Added an `httptest` contract test that returns intentionally unsorted results
  and proves they hydrate as candidate IDs and deterministic ranks.
- Added regression tests for oversized requests and malformed response shapes.

### Why

The live probe established that llama.cpp returns candidate indexes and may
return negative scores. Mapping index to durable IDs before sorting prevents a
later executor from accidentally treating an ephemeral server array position
as an experiment artifact.

### What worked

- `GOCACHE=/tmp/rag-eval-go-build GOWORK=off go test ./pkg/raglab -count=1`
  passed.
- The request test verifies the exact `/v1/rerank` payload and the adapter
  ranks `-1` ahead of `-5`, matching the observed higher-is-better convention.

### What didn't work

- The first test invocation used Go's default build cache under
  `~/.cache/go-build`; this sandbox mounts it read-only. Re-running with a
  ticket-local-safe cache under `/tmp` resolved the environmental limitation.
  No product code failed.

### What I learned

Response ordering is not needed for correctness. The combination of submitted
candidate order, returned index, score validation, and deterministic sorting is
sufficient to construct a stable artifact.

### What was tricky to build

The adapter must distinguish three boundaries: the request's `TopN` requested
from llama.cpp, the server's result cardinality, and the final rank assigned
after local deterministic sorting. Keeping each boundary explicit makes a
short or malformed server response fail before it contaminates a run trace.

## Step 6: Expose the policy through the JavaScript fluent builder

The Goja module now provides the ergonomic authoring form specified in the
ticket while preserving the typed Go builder as the authority:

```js
.rerank(x => x
  .crossEncoder("qllama/bge-reranker-v2-m3:q4_k_m")
  .candidates(50)
  .results(10))
```

The temporary JavaScript configuration object only collects values. It then
calls `RetrievalBuilder.RerankCrossEncoder`, so duplicate and invalid policy
diagnostics remain centralized in `pkg/raglab` rather than reimplemented in
the binding.

### Prompt Context

**User prompt (verbatim):** “go ahead.”

**Assistant interpretation:** Continue through the approved API authoring
task after completing the adapter.

**Commit (code):** `51c9f89` — "feat: expose reranking through JavaScript DSL"

### What I did

- Added `retrieval.rerank(configure)` and its fluent `crossEncoder`,
  `candidates`, and `results` primitives.
- Added the `RerankingBuilder` TypeScript declaration.
- Added lower-camel `retrieval.reranking` projection in `toSpec()` and
  `toJSON()`.
- Extended the Goja module integration test to build, validate, persist, and
  inspect a reranking policy from JavaScript.

### Why

Experiments must declare reranking before persistence, not pass runtime model
choices as incidental JavaScript variables. The builder makes the policy part
of the immutable specification and therefore part of its fingerprint.

### What worked

- `GOCACHE=/tmp/rag-eval-go-build GOWORK=off go test ./pkg/gojamodules/rag -count=1`
  passed.

### What didn't work

- The first integration assertion found that the existing JavaScript
  `retrievalValue` projection omitted `Reranking`. The Go policy itself was
  correct; adding the lower-camel projection fixed the authoring contract.

### What I learned

Authoring, canonical Go state, and JavaScript JSON output are separate
boundaries. A feature is incomplete when any one of them omits a semantic
field, even if the fingerprint already includes it.

### What was tricky to build

The returned Goja maps must use `candidateCount`, while Go persistence uses
`candidate_count`. The explicit projection keeps JavaScript scripts consistent
with TypeScript without altering the stored immutable representation.

## Step 7: Execute reranking and persist both candidate orders

The executor now treats reranking as an explicit runtime capability. It forms
the ordinary raw or RRF candidate list, limits that list by the immutable
`candidateCount`, submits every candidate for scoring, and applies the
declared reranking result window before final document collapse and final
result truncation. The trace records the stable reranker identity, every
submitted candidate's pre-rerank rank and retrieval score, and every returned
scored result.

### Prompt Context

**User prompt (verbatim):** “do you know how to continue? you should only stop
if you're at a loss of what to do next”

**Assistant interpretation:** Resume a deterministic test-fixture correction
and continue the task sequence rather than treating ordinary validation as a
blocker.

**Commit (code):** pending — executor integration and regression coverage are
ready to commit.

### What I did

- Added the explicit `ExecutionOptions.Reranker` capability and fail-fast
  validation for a plan that declares reranking without supplying one.
- Used the policy candidate count as the RRF output limit when reranking is
  enabled, avoiding premature truncation to the final display count.
- Added `applyReranking`, which checks model identity, requires hydrated chunk
  text, preserves candidate IDs, requests scores for the complete submitted
  window, and rejects unknown returned IDs.
- Persisted `reranking.identity`, `reranking.candidates`, and
  `reranking.results` in each query trace.
- Added an executor regression test using a deterministic reverse reranker;
  it proves the first raw candidate becomes the second final result and that
  the complete before/after trace is recorded.

### Why

Reranking cannot be inferred after a run. The trace must retain candidate
order and original retrieval scores so a user can distinguish poor retrieval
from changed cross-encoder ranking.

### What worked

- `GOCACHE=/tmp/rag-eval-go-build GOWORK=off go test ./pkg/raglab -run
  TestExecutorReranksBoundedCandidatesAndPersistsBothOrders -count=1` passed.
- `GOWORK=off go test ./pkg/raglab -count=1` passed outside the sandbox. This
  full run also exercised the adapter's `httptest` HTTP contract.

### What didn't work

- The initial new fixture omitted required metric configuration, then omitted
  the relevance threshold required for a graded metric. Both errors were
  deterministic builder diagnostics; adding the explicit threshold fixed the
  test setup. They did not indicate an executor defect.
- The sandbox cannot bind the loopback listener used by `httptest`; the normal
  environment full suite passed.

### What I learned

For the current RRF implementation, duplicate document revisions are already
collapsed inside `FuseWeightedRRF`. Single-channel retrieval, however, retains
raw chunk candidates until reranking. The trace makes that distinction visible
and preserves the evidence needed for Task 11's collapse-order decision.

### What was tricky to build

The plan's `results` and the reranker policy's `results` have distinct
meanings. The latter is the post-score, pre-collapse window. The executor
scores the whole candidate set, records it, applies the policy window, then
collapses/truncates to the former final user-facing limit.

## Step 8: Render rank changes in the Evaluation trace inspector

The web trace inspector previously displayed only a truncated JSON block. It
now detects the optional reranking trace and renders a compact table for each
query card. The table intentionally shows every candidate that was submitted,
including a visible `truncated` post-rerank state when a candidate was scored
but did not survive the configured result window.

### What I did

- Added a narrow TypeScript projection for the optional trace payload without
  changing the generic API transport type.
- Added the cross-encoder identity and candidate transition table to
  `EvaluationPage`.
- Added responsive grid styling while retaining the complete JSON trace below
  it for unmodeled fields and forensic inspection.

### What worked

- `pnpm --dir web typecheck` passed.

### Why

An experimenter needs to see rank movement directly. A raw JSON blob does not
make it practical to distinguish a candidate that was never retrieved from a
candidate that was retrieved, scored, and then truncated.

## Step 9: Make the llama.cpp reranker an explicit JavaScript runtime capability

`lab.execute()` obtains operational capabilities from its laboratory handle.
The JavaScript `open` method now accepts a typed `reranker` object and creates
the existing strict llama.cpp adapter. This closes the final path between a
fluent experiment policy and the private Mac service while preserving the
separation between immutable experiment identity and mutable endpoint details.

```js
const lab = rag.open({
  database: "ttc-rag.sqlite",
  execution: "allowRuns",
  queryEmbed,
  reranker: {
    kind: "llama.cpp",
    baseURL: "http://127.0.0.1:18012",
    model: "qllama/bge-reranker-v2-m3:q4_k_m"
  }
});
```

### What I did

- Added `Reranker` to `raglab.OpenOptions` and defaulted it into execution in
  `Laboratory.Execute`.
- Added JavaScript validation for `reranker.kind === "llama.cpp"` and adapter
  construction from `baseURL`, `model`, and optional request limit.
- Added TypeScript declarations for `LlamaCPPRerankerOptions`.
- Added a Goja module test that captures the constructed capability and rejects
  an unsupported kind.
- Made `stringProperty` nil-safe after the test revealed that an omitted Goja
  optional property can be represented as nil rather than undefined.

### What worked

- Focused Goja/laboratory tests passed.
- `pnpm --dir web typecheck` passed.

### Why

The executor correctly rejects a reranking plan without a capability. Providing
this capability through `open` makes the authority explicit at the user-facing
boundary, rather than hiding an endpoint in environment state or in persisted
experiment JSON.

## Step 10: Run the first live frozen-TTC BGE reranker experiment

The local tunnels were re-established and the generated xgoja runtime was
rebuilt from the current branch. The JavaScript experiment then completed all
20 frozen TTC cards with BGE cross-encoder reranking. The result is an
append-only experiment run with complete candidate and score traces, rather
than an in-memory benchmark observation.

### What worked

- Run `run_76e8425d56b07b134915a749e05bb03f` succeeded with 20 traces.
- MRR was `0.9473684210526315`; mean relevant recall@10 was
  `0.8947368421052632`.
- Accumulated per-card execution time was 26,801 ms; wall-clock was 26,976 ms.
- SQLite inspection confirmed BGE identity and one-to-one candidate/result
  trace arrays for the sampled query cards.

### What did not work

- The previous local SSH tunnels had expired. The private services themselves
  were healthy; recreating `rag-reranker-mimimi` and `rag-ollama-mimimi` tmux
  tunnels restored them.
- The sandbox blocked the generated xgoja build's required Go toolchain proxy
  download. The permitted normal-environment build succeeded without source or
  dependency changes.

### What I learned

The current RRF implementation has already document-collapsed candidates, so
the first reranker does not always receive its configured maximum of 50 texts.
This is observable in the saved traces and is the empirical premise for the
still-pending collapse-order comparison.

### What warrants a second pair of eyes

- Review whether `RerankingSpec.Results` should remain a distinct pre-collapse
  window or be renamed before the JavaScript API is exposed.

### What should be done in the future

- Implement the llama.cpp adapter against the recorded index/score contract.

### Code review instructions

- Start at `pkg/raglab/reranker.go`, then inspect validation in
  `pkg/raglab/builder.go`.
- Run the focused raglab test command above.

### Technical details

```text
RerankingSpec(model, candidateCount, results) -> canonical retrieval JSON -> fingerprint
Reranker interface -> later llama.cpp adapter -> later executor
```
