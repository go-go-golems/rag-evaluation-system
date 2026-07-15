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
