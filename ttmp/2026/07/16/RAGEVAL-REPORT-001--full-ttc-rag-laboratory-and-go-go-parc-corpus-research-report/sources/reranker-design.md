---
Title: Reranker stage analysis design and implementation guide
Ticket: RAGEVAL-RERANK-001
Status: active
Topics:
    - rag
    - reranking
    - ttc
    - geppetto
    - ollama
    - rag-eval
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://internal/services/experimentrun/service.go
      Note: Append-only trace and terminal summary persistence
    - Path: repo://pkg/gojamodules/rag/module.go
      Note: Explicit JS runtime capability pattern
    - Path: repo://pkg/raglab/executor.go
      Note: Current candidate retrieval, RRF fusion, trace, and metric insertion point
    - Path: repo://pkg/raglab/types.go
      Note: Immutable retrieval and experiment specification types
ExternalSources: []
Summary: Intern-ready design for a reproducible cross-encoder reranking stage after TTC candidate retrieval and before final result hydration.
LastUpdated: 2026-07-15T16:55:29.71504488-04:00
WhatFor: Define the first reranker architecture, durable experiment contract, implementation phases, and validation plan.
WhenToUse: Use before adding reranking code, operating llama.cpp, modifying query traces, or interpreting reranked TTC results.
---


# Reranker stage analysis design and implementation guide

## Executive Summary

The TTC laboratory retrieves candidates with BM25, vector search, or weighted
reciprocal-rank fusion (RRF), then evaluates the final ranked chunk list. It
currently has no semantic scoring stage that reads a query and candidate text
together. This ticket adds an optional cross-encoder reranker between candidate
selection and final result formation.

The initial adapter calls a local `llama-server` reranking endpoint. The server
is run on the Mac with a dedicated reranker model and reached from the worktree
through a private SSH loopback tunnel. The experiment specification records a
stable reranker identity and candidate-window policy; operational endpoint and
credentials remain outside immutable identity. The trace preserves both the
pre-rerank candidates and post-rerank scores.

## 1. Scope and non-goals

In scope:

- rerank up to a bounded number of already retrieved chunk candidates;
- persist reranker configuration identity, per-query candidate order, scores,
  timing, and failures in append-only experiment records;
- compare reranked runs with the existing 20-card TTC vector and weighted-RRF
  observations;
- expose enough trace data for a user to understand why rank changed.

Out of scope for the first implementation:

- generation of summaries/questions or parent mappings;
- changing corpus embeddings, BM25 indexes, or relevance labels;
- an Ollama-specific synthetic scoring protocol;
- a generic remote multi-tenant reranking service;
- backwards-compatible public APIs. This RAG API is unmerged and may receive
  the necessary v1 shape directly.

## 2. Current system and insertion point

`pkg/raglab/executor.go` implements one query-card loop. It retrieves channel
hits, fuses multiple channels using weighted RRF, optionally collapses by
document, truncates to requested result count, computes relevance metrics, and
stores an immutable trace. Lines 190–241 are the critical execution path:
query embedding is shared across vector channels, channel results are recorded,
and RRF produces `trace.Results`.

```text
evaluation card
  | query embedding (only for vector)
  +--> BM25 channel -----------+
  +--> vector channel ---------+--> weighted RRF --> result list --> metrics
                                                           |
                                                   proposed reranker stage
```

`executionTrace` currently stores channels, fusion rows, final results, and
embedding/retrieval/fusion/total milliseconds (`executor.go:164–180`). The
experiment service writes that encoded trace once per card and rejects writes
after a terminal summary (`internal/services/experimentrun/service.go:227–280`).
The reranker must enrich this trace before `RecordQueryTrace`, not create a
mutable side table that can drift from the run.

The existing raw experiment has a proved baseline. The 20-card JS/Geppetto
weighted-RRF run `run_20b25df32dc874af1265a9e6ccf87570` completed with MRR
0.820175 and relevant recall@10 0.815789. It uses frozen TTC artifact IDs and
a 768D `nomic-embed-text` query embedder. Reranking must use that baseline,
not rebuild the corpus, as its first comparison.

## 3. Runtime API reference

The official llama.cpp HTTP server documentation exposes `POST /reranking` and
aliases `/rerank`, `/v1/rerank`, and `/v1/reranking`. It requires a reranker
model and `--embedding --pooling rank`; the request contains `query`, a
`documents` string array, and optional `top_n`. Use `/v1/rerank` initially.
The documentation calls the endpoint subject to change, so Task 4 captures the
actual response of the installed version before the adapter is frozen.

Example server lifecycle on the Mac (model path/model download mechanism is an
operator choice, not application code):

```bash
llama-server --model /path/to/bge-reranker-v2-m3-q4_k_m.gguf \
  --embedding --pooling rank --rerank --host 127.0.0.1 --port 8012

# local workstation tunnel, managed in tmux
ssh -N -L 127.0.0.1:18012:127.0.0.1:8012 mimimi-2.local
```

Probe contract:

```json
{
  "model": "qllama/bge-reranker-v2-m3:q4_k_m",
  "query": "How does TTC calculate a payroll adjustment?",
  "documents": ["candidate chunk 0", "candidate chunk 1"],
  "top_n": 2
}
```

The adapter must treat the response index as an index into the submitted array,
not an implicit chunk rank. It reconstructs chunk identity from its original
candidate array and rejects out-of-range, duplicate, non-finite, or incomplete
results.

## 4. Proposed architecture

### Go contracts

```go
type RerankCandidate struct {
    ID string                 // immutable chunk ID; never an array position
    Text string               // hydrated candidate text, size-limited
    OriginalRank int
    RetrievalScore float64    // RRF score when available
}

type RerankRequest struct {
    Query string
    Candidates []RerankCandidate
    TopN int
}

type RerankResult struct {
    CandidateID string
    Index int
    Score float64
    Rank int
}

type Reranker interface {
    Rerank(ctx context.Context, request RerankRequest) ([]RerankResult, error)
    Identity() RerankerIdentity
}

type RerankingSpec struct {
    Kind string `json:"kind"`           // "crossEncoder"
    Model string `json:"model"`         // stable configured identity
    CandidateCount int `json:"candidate_count"`
    Results int `json:"results"`
}
```

`RerankingSpec` belongs in `ExperimentSpecification.Retrieval`, after fusion
policy. The endpoint URL, bearer token, and timeouts are `RerankerOptions`
passed to `Laboratory`/`Executor`, exactly like the current explicit query
embedder: capability is runtime configuration, not a property inferred from an
embedding artifact.

### Execution algorithm

```text
retrieve each channel -> fuse RRF -> take candidateCount
       -> hydrate bounded chunk text
       -> reranker.Rerank(query, candidates)
       -> validate one score per submitted candidate
       -> sort by score descending; stable tie-break by original rank then ID
       -> apply collapse policy
       -> trim requested result count
       -> calculate metrics and persist trace
```

Pseudocode:

```go
candidates := fused[:min(spec.Reranking.CandidateCount, len(fused))]
inputs := hydrateAndLimit(candidates, maxCharsPerCandidate)
scored := reranker.Rerank(ctx, RerankRequest{Query: query, Candidates: inputs, TopN: len(inputs)})
ordered := validateAndOrder(inputs, scored) // strict one-to-one invariant
results := collapse(ordered, spec.Retrieval.Collapse)
results = results[:min(spec.Retrieval.Results, len(results))]
trace.Reranking = RerankingTrace{Candidates: inputs, Scores: scored, Milliseconds: elapsed}
```

## 5. Data and trace design

Extend `executionTrace` rather than overwriting `Channels`, `Fusion`, or
`Results`:

```json
"reranking": {
  "identity": {"kind":"llama.cpp", "model":"bge-reranker-v2-m3-q4_k_m"},
  "candidateCount": 50,
  "submittedCount": 50,
  "returnedCount": 50,
  "milliseconds": 420,
  "items": [
    {"chunkId":"...", "originalRank":4, "score":0.817, "rerankedRank":1}
  ]
}
```

Do not persist candidate text twice in the trace unless necessary for debugging;
the immutable chunk ID plus hydrated citation remains the canonical source.
Record text truncation count and character limits because cross-encoder scores
are conditioned on exactly that input. Include latency in a new
`RerankingMilliseconds` field and retain current timing fields for comparison.

## 6. Decisions

### Decision: use a native llama.cpp reranking endpoint

- **Context:** Ollama stores useful reranker artifacts, but the laboratory needs
  a documented score-per-document HTTP contract.
- **Options considered:** Ollama generation/prompt scoring emulation; LM Studio;
  llama.cpp native reranking endpoint; hosted reranking API.
- **Decision:** Start with `llama-server` `/v1/rerank` over a private tunnel.
- **Rationale:** Official llama.cpp documentation defines query/documents/top_n
  and the required reranker flags. It can run on the existing Apple Silicon
  machine and does not introduce billed remote-provider dependency.
- **Consequences:** A small HTTP adapter and model-server operator playbook are
  required; exact response shape must be probed before coding.
- **Status:** accepted.

### Decision: rerank after fusion and before final collapse

- **Context:** RRF combines independent evidence; a cross encoder should score
  a bounded shared candidate set. Early reranking per channel duplicates calls
  and distorts fusion scores.
- **Options considered:** per-channel rerank then fuse; fuse then rerank; collapse
  then rerank.
- **Decision:** Fuse first, rerank chunk candidates, then apply final collapse.
- **Rationale:** It preserves the baseline retrieval pipeline and gives the
  cross encoder fine-grained chunk context. Parent/document collapse remains a
  measurable later decision.
- **Consequences:** Candidate count must be bounded and duplicate chunks may be
  scored; trace data makes that trade-off visible.
- **Status:** proposed; confirm with Task 11.

### Decision: no automatic fallback

- **Context:** A silent fallback to RRF would create a run whose declared
  reranker policy differs from its observed result.
- **Options considered:** fall back on transport error; fail run; omit reranker
  from experiment identity.
- **Decision:** A declared reranker failure fails the append-only run.
- **Rationale:** Reproducibility is more important than availability for an
  evaluation laboratory.
- **Consequences:** UI and operator docs must make failure states understandable.
- **Status:** accepted.

## 7. JavaScript builder API

The fluent builder mirrors the Go types and leaves network settings out:

```javascript
const experiment = rag.experiment("ttc-rrf-bge-rerank", (e) => e
  .corpus(snapshot).chunks(chunks).bm25(bm25).embeddings(vectors).evaluation(cards)
  .retrieval((r) => r
    .channel("lexical", (c) => c.bm25().topK(50))
    .channel("semantic", (c) => c.vector().topK(50))
    .fuse((f) => f.rrf().rankConstant(60).weight("semantic", 2))
    .rerank((x) => x.crossEncoder("bge-reranker-v2-m3-q4_k_m").candidates(50).results(20))
    .collapse("document").results(10)));
```

The runtime opens with an explicit reranker capability, analogous to
`queryEmbed`. A future `require("rag")` option may accept a synchronous
`rerank(query, candidates)` callback, but the first implementation should
provide a typed Go `llama.cpp` adapter. That lets Go enforce cancellation,
payload limits, JSON validation, and source-citation identity.

## 8. File-level implementation guide

1. Add pure types and canonical validation in `pkg/raglab/types.go` and
   `pkg/raglab/builder.go`. Validate positive `candidateCount`, nonempty model,
   and `results <= candidateCount`.
2. Add `pkg/raglab/reranker.go` for the interface and deterministic ordering
   tests. Do not import HTTP or Goja there.
3. Add `pkg/raglab/reranker_llamacpp.go` for the HTTP adapter. Use contexts,
   bounded request body, no implicit retry, and `github.com/pkg/errors`.
4. Extend `pkg/raglab/executor.go`: fuse, form candidates, hydrate text,
   rerank, attach trace detail, collapse, and only then compute metrics.
5. Extend `pkg/gojamodules/rag/module.go` and `typescript.go` with one
   `.rerank()` builder object. Update native-module tests for callback/type
   errors and canonical JSON.
6. Extend `web/src/services/api.ts` and `EvaluationPage` to render a score and
   before/after rank table. Use Bootstrap-compatible existing page components.
7. Store probes and comparative scripts under this ticket's `scripts/` with
   numbered names; do not hide experimental commands in shell history.

## 9. Validation and experiment matrix

Unit tests:

- score response index mapping, duplicate/missing indexes, NaN/Inf score,
  stable ties, cancellation, body-size rejection, and zero candidates;
- canonical spec/fingerprint changes when model or candidate count changes;
- no reranker configured preserves the present baseline result exactly.

Integration tests:

- an `httptest` llama.cpp response fixture;
- append-only trace and terminal summary behavior on reranker success/failure;
- generated `rag-eval-js` builder and TypeScript declaration smoke test.

TTC study table:

| Run | Retrieval | Reranker | Primary measures |
| --- | --- | --- | --- |
| A | raw vector | none | MRR, recall@10, latency |
| B | weighted RRF | none | existing baseline |
| C | weighted RRF | BGE v2 m3 | delta quality, rerank ms, candidate budget |
| D | weighted RRF | Qwen3 4B or 8B | model-quality/latency trade-off |

All runs use the same evaluation cards and base immutable artifacts. Record
local cost as hardware/energy unestimated unless a provider charges; record
model bytes separately from shared corpus index storage.

## 10. Operational playbook outline

1. Verify the Mac over SSH and identify the exact reranker artifact/version.
2. Start `llama-server` in a dedicated tmux session with loopback binding,
   `--embedding --pooling rank --rerank`, and a chosen port.
3. Start a separate local SSH loopback tunnel in tmux.
4. Run a curl probe containing two or three known candidate strings; save raw
   response and model/server version in `scripts/` output notes.
5. Configure the Go adapter endpoint locally; never commit SSH hostnames,
   credentials, or tokens into experiment specs.
6. Run one bounded TTC comparison before a complete matrix. Inspect durable
   trace records and stop the server/tunnel cleanly after the observation.

## 11. Risks and open questions

- llama.cpp says its reranking route may change; probe before stabilizing the
  decoder and retain an adapter-specific test fixture.
- Cross encoders have context limits. Candidate text must be truncated with a
  named policy, and truncation must be observable.
- Candidate chunks from the same document may dominate before collapse. Measure
  chunk-first and document-first alternatives rather than assuming either.
- The Mac has BGE and Qwen artifacts in Ollama storage, but llama.cpp may need
  compatible GGUF paths or downloads. Model availability is an operator
  preflight, not a code fallback.
- Do not claim a quality improvement from one 20-card set; report per-card
  differences and regressions alongside aggregate MRR/recall.

## References

- `pkg/raglab/executor.go:69–243` — current execution, fusion, timing, trace,
  and relevance calculations.
- `pkg/raglab/types.go:129–210` — representation, channel, retrieval, metric,
  and experiment spec data model.
- `internal/services/experimentrun/service.go:227–316` — immutable trace and
  terminal-summary persistence.
- `pkg/gojamodules/rag/module.go:84–100` — explicit JS runtime capability
  pattern for query embedding.
- `ttmp/2026/07/14/RAGEVAL-RAG-DSL-001--typed-fluent-javascript-rag-laboratory-module/reference/03-mimimi-ollama-tunnel-operator-playbook.md` — proven private Mac tunnel pattern.
- [llama.cpp HTTP server reranking API](https://github.com/ggml-org/llama.cpp/blob/master/tools/server/README.md) — endpoint aliases, required flags, and request fields.
