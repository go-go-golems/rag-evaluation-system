# TTC BGE reranker run — 2026-07-16

## Provenance

- Script: `03-run-ttc-bge-reranker.js`
- Embedding model: `nomic-embed-text`, 768 dimensions, via private Ollama SSH
  loopback tunnel at `127.0.0.1:11435`.
- Reranker: `qllama/bge-reranker-v2-m3:q4_k_m`, via private llama.cpp SSH
  loopback tunnel at `127.0.0.1:18012`.
- Experiment run: `run_76e8425d56b07b134915a749e05bb03f`
- Query cards: 20 total; 19 answerable.

## Immutable recorded result

```json
{
  "meanReciprocalRank": 0.9473684210526315,
  "meanRelevantRecallAtResults": 0.8947368421052632,
  "totalMilliseconds": 26801,
  "wallClockMilliseconds": 26976
}
```

The run stored exactly 20 immutable query traces. Every inspected trace had a
`reranking` object whose identity named the BGE model and whose candidate and
result arrays had matching cardinality. Candidate counts varied (for example
4, 11, and 3) because the current RRF implementation has already collapsed
duplicate document revisions before it emits the bounded candidate window.

## Baseline comparison

The earlier JavaScript weighted-RRF run
`run_20b25df32dc874af1265a9e6ccf87570` recorded MRR `0.8201754385964911` and
mean relevant recall@10 `0.8157894736842105`. On the same frozen evaluation
dataset, BGE reranking improved MRR by about 0.1272 and recall by about 0.0789.
This is one controlled local observation, not a general model-quality claim.
