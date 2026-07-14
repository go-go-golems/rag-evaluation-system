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
RelatedFiles: []
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
