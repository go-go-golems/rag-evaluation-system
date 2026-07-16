---
Title: TTC RAG laboratory baseline and immutable experiment runs
Ticket: RAGEVAL-TTC-LAB-001
Status: active
Topics:
    - rag
    - rag-eval
    - ttc
    - corpus
    - chunking
    - embeddings
    - search
    - evaluation
    - workflow
    - web
    - frontend
    - intern-guide
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://2026/07/15/RAGEVAL-RERANK-001--reranking-stage-for-the-immutable-ttc-rag-laboratory/scripts/02-ttc-eval-v2-proposed-stratified-cards.md
      Note: 72-card stratified expansion draft awaiting source validation
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/04-ttc-evaluation-expansion-v0-70-proposed-cards.md
      Note: 70 additional source-discovery cards for the 240-card target
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/04-ttc-v2-support-faq-adjudication-batch-01.md
      Note: First concrete twelve-card direct-source adjudication packet
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/05-ttc-evaluation-expansion-y-v0-80-proposed-cards.md
      Note: 80 additional comparison, climate, diagnostic, procedure, transaction, and abstention candidates
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/05-ttc-v2-proposed-split-and-chunk-label-reconciliation.md
      Note: Reconciles reranker split and chunk-label drafts into canonical TTC v2 decisions
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/06-ttc-expansion-audit-and-50-card-source-grounded-draft.md
      Note: 50 additional SQLite-grounded candidate cards and coverage audit
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/06-ttc-v2-240-card-partition-and-leakage-audit-protocol.md
      Note: 240-card partition and union-find leakage audit contract
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/07-ttc-expansion-source-validation-result.md
      Note: 142 unique source IDs resolved across 120 candidate cards
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/08-ttc-expansion-candidate-dataset-registration.md
      Note: 148-card candidate dataset registration and withheld conflict explanation
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/08-audit-ttc-evaluation-cards.py
      Note: Read-only candidate-card parser and future partition/leakage audit
    - Path: repo://ttmp/2026/07/15/RAGEVAL-RERANK-001--reranking-stage-for-the-immutable-ttc-rag-laboratory/scripts/06-ttc-v1-development-holdout-regression-split-draft.md
      Note: Leakage-aware pilot split and expansion protocol
    - Path: repo://ttmp/2026/07/15/RAGEVAL-RERANK-001--reranking-stage-for-the-immutable-ttc-rag-laboratory/scripts/07-ttc-chunk-level-evidence-label-proposal.md
      Note: Proposed chunk-level evidence labels for pilot cards
ExternalSources: []
Summary: Build a bounded TTC retrieval baseline and introduce immutable, content-addressed experiment specifications, runs, artifacts, traces, and comparisons.
LastUpdated: 2026-07-14T16:10:00-04:00
WhatFor: Track the implementation of the first reproducible TTC RAG laboratory slice in the maintained RAG Evaluation System.
WhenToUse: Start here when implementing or reviewing TTC corpus snapshots, baseline retrieval, experiment runs, evaluation, or laboratory UI work.
---








# TTC RAG laboratory baseline and immutable experiment runs

## Goal

Build a web-testable TTC baseline across fixed, sentence, and Markdown chunking; BM25, vector, and hybrid retrieval; and at least 20 queries with named relevance judgments. Store corpus membership, configurations, artifacts, query traces, metrics, costs, timings, and terminal run results under deterministic or append-only identities so no experiment can silently change after execution.

## Rebuilt source artifact

`data/ttc-wordpress-rag.sqlite` was rebuilt on 2026-07-13 from the local TTC MySQL corpus and passed the repository validator and SQLite integrity check.

```text
documents:      3,258
products:       2,600
non-products:     658
output sha256:  c55953ee0d9289577062ac11001c25f63c0286ace45dbc6b4b056c11b0ea6db4
```

## Documents

- [Design and implementation guide](design-doc/01-ttc-rag-laboratory-baseline-and-immutable-experiment-runs-design-and-implementation-guide.md)
- [Evaluation dataset authoring and adjudication protocol](design-doc/02-evaluation-dataset-authoring-and-adjudication-protocol.md)
- [TTC baseline evaluation dataset v1 candidate cards](reference/02-ttc-baseline-evaluation-dataset-v1-candidate-cards.md)
- [Implementation diary](reference/01-implementation-diary.md)
- [Tasks](tasks.md)
- [Changelog](changelog.md)

## Delivery stages

1. **Baseline:** Import an explicit 200-document TTC snapshot, create three exact chunk sets, compute one real 768D embedding set, run BM25/vector/RRF, and judge at least 20 queries.
2. **Immutable runs:** Fingerprint every semantic component, publish indexes content-addressably, append run events, insert terminal summaries once, and expose trace/comparison APIs and UI.

Generated data remains under ignored `data/` paths. Ticket documentation and any ticket-specific scripts are committed under this ticket workspace.
