---
Title: TTC Expansion Development Retrieval Run Results
Ticket: RAGEVAL-TTC-LAB-001
Status: review
Topics:
    - ttc
    - rag-eval
    - evaluation
    - embeddings
    - search
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/10-score-expansion-candidate-traces.py
      Note: Reproducible provisional scorer
ExternalSources: []
Summary: First 148-card BM25/vector/RRF development run using the Mac Ollama embedding tunnel.
LastUpdated: 2026-07-16T00:00:00Z
WhatFor: Establish a development signal and latency baseline before adjudication.
WhenToUse: Use for engineering diagnosis only; do not report as a frozen benchmark.
---

# Run configuration

- Dataset: `candidate:ttc-expansion-v0` (148 cards; two policy conflicts excluded)
- Corpus snapshot: `sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409`
- BM25 artifact: `sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691`
- Embedding set: `sha256:2665c5249b8352ce6904fc00c934534dd179f3eeef0a6a75429a9034be0e03e0`
- Query embedding endpoint: Mac Ollama through `127.0.0.1:11435`
- Candidate limit: 50 per channel; RRF constant 60; fused result limit 10

## Provisional metrics

| Method | Answerable cards | Recall@1 | Recall@3 | Recall@10 | MRR | Relevant recall@10 |
|---|---:|---:|---:|---:|---:|---:|
| BM25 | 144 | 0.7847 | 0.8403 | 0.8889 | 0.8221 | 0.7442 |
| Vector | 144 | 0.8750 | 0.9653 | 0.9722 | 0.9174 | 0.8588 |
| RRF hybrid | 144 | 0.8542 | 0.9375 | 0.9722 | 0.9005 | 0.8947 |

Four cards are unanswerable controls and are excluded from answerable retrieval
metrics. The run still records their retrieved candidates for later abstention
and false-support analysis.

## Latency and cost

| Statistic | Total query latency |
|---|---:|
| Mean | 173 ms |
| P50 | 175 ms |
| P95 | 230 ms |
| Range | 72–465 ms |

Recorded embedding cost is zero in billed currency because the run used the
user-owned Ollama service. Hardware amortization and energy are not estimated.

## Interpretation boundary

These values are **not evaluation claims**. The relevant document IDs came from
the authoring draft and have not received human evidence-span adjudication.
The numbers are useful for identifying engineering regressions, comparing
retrieval implementation choices, and selecting candidate configurations on
development data. They must be recomputed after source/chunk judgments are
reviewed and must not be used to open the holdout partition.

## Reproduction

```sh
PYTHONDONTWRITEBYTECODE=1 python3 scripts/10-score-expansion-candidate-traces.py \
  --manifest scripts/09-ttc-expansion-v0-manifest.json \
  --traces data/artifacts/traces/ttc-expansion-v0.json \
  --out data/artifacts/metrics/ttc-expansion-v0-candidate.json
```
