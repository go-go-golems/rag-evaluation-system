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
