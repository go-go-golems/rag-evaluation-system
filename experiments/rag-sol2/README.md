# rag-sol2 one-time parity extraction

This directory is the sole retained study form of the former `rag-sol2` runtime. It contains pure `require("rag")` authoring, immutable catalog inputs, and a compact candidate result summary. It does not import the deleted provider, cache, index, hashing, evaluation, or lifecycle libraries.

## Scope

The deterministic oracle was executed once before deletion from `rag-sol2`:

- 2 agents-view units;
- 4 recursive chunks;
- 4 raw, 4 summary, and 8 question representations;
- exact graded fixture metrics: precision@3 `2/3`, recall@3 `1`, hit-rate@3 `1`, MRR `1`, nDCG@3 `0.9828422279067397`.

The canonical fixture is `pkg/ragoperators/testdata/rag-sol2-parity-v1.json`. Canonical tests additionally prove:

- tool records do not split an assistant run;
- UTF-8 byte/rune ranges reconstruct only source bytes;
- every representation has explicit chunk/unit/source lineage;
- question multiplicity exists before channel-local collapse;
- one collapse key contributes at most once per channel;
- fusion has at most one contribution per channel;
- final results hydrate source chunks and citations;
- fixture provider cost is exactly zero;
- all five variants × chunk/unit collapse produce ten stable cells.

## Discrepancy classification

There are **zero blocking discrepancies** for the retained semantics.

Intentional corrections:

1. Canonical units omit rag-sol2's synthetic `\n\n` separator. Every canonical byte must belong to a source record; display separators are not evidence.
2. Canonical IDs are Go-owned manifest/record identities. JavaScript keys and duplicated SHA/canonical-stringify code were not preserved.
3. Recursive parity uses explicit rune-level configuration. Exact source reconstruction is authoritative; old paragraph-boundary bytes containing synthetic separators are not preserved.
4. Generated text comes from explicitly named fixture providers. It proves schemas, multiplicity, lineage, caching, channels, and costs—not model quality or real-model text equality.
5. Canonical retrieval collapses each channel before weighted fusion and hydrates source chunks explicitly. Generated text is never cited as evidence.
6. Unit-target evaluation deduplicates hydrated unit IDs. The old chunk-collapse evaluator could count multiple chunks from one relevant unit more than once; this multiplicity bias is intentionally removed.
7. The candidate matrix includes BM25 and deterministic hash-vector channels. Its numbers are not comparable to the earlier frozen TTC BM25/vector/model baselines.

## Candidate matrix

`candidate-result.json` records one researchctl run for each of ten cells over `candidate:ttc-expansion-v1` (148 queries). Every specification and run was labeled `candidate`; the summary explicitly sets `benchmarkClaim: false` and `fixtureProviders: true`.

The matrix is execution/parity evidence only. It does not freeze judgments, models, prompts, provider servers, tokenizers, or holdouts and must not be described as an adjudicated benchmark.

Execution evidence measured approximately 2m49s wall time, 1.62 GB maximum RSS, and 869 MB of researchctl project storage for all ten cells. Static preparation/indexing was shared across queries within each cell, but not across cells.

## Commands

```bash
rag-eval study validate experiments/rag-sol2/study.js \
  --inputs experiments/rag-sol2/inputs.json \
  --ttc-database data/rag-eval.db

rag-eval study run experiments/rag-sol2/study.js \
  --inputs experiments/rag-sol2/inputs.json \
  --ttc-database data/rag-eval.db \
  --project project.yaml \
  --experiment-id EXP-RAG-PARITY \
  --researchctl-command researchctl \
  --worker-command rag-worker
```
