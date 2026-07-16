---
Title: TTC evaluation corpus v2 foundation and adjudication protocol
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
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Protocol for expanding the TTC pilot into a 110-card, source-grouped, immutable evaluation foundation."
LastUpdated: 2026-07-16T16:09:56.653338305-04:00
WhatFor: "Define development, holdout, regression, chunk-evidence, and adjudication boundaries before immutable TTC dataset v2 creation."
WhenToUse: "Read before authoring TTC cards, tuning BGE retrieval, or using a TTC result for model selection."
---

# TTC evaluation corpus v2 foundation and adjudication protocol

## Executive Summary

The existing 20-card TTC pilot proved the laboratory but cannot select a
chunking, fusion, or reranking policy reliably. TTC v2 is a staged 110-card
foundation: 75 development cards, 20 untouched holdout cards, and 15 regression
cards. The later target is a 185-card benchmark (120/40/25). Every split is
grouped by evidence family, so documents with overlapping answers do not leak
between tuning and evaluation.

## Problem Statement

The pilot has only 19 answerable cards and document-level judgments. Its
existing cards often share evidence, and proposed new cards are not yet
adjudicated at revision or chunk level. Tuning BGE candidate budgets or
collapse placement on all of them would overfit to repeated product/FAQ text.

## Proposed Solution

### Partition contract

```text
frozen TTC corpus snapshot
  -> evidence-family grouping
  -> 75 development cards: choose configurations
  -> 20 holdout cards: one frozen-configuration evaluation
  -> 15 regression cards: known failures retained across changes
```

The current pilot supplies initial grouped cards; the proposed 72-card pool is
an authoring queue, not an evaluation dataset. Every candidate becomes scored
only after an adjudicator supplies source revision IDs, authoritative and
substantial evidence, a rationale, and negative/ambiguous status where
applicable.

### Card schema

```json
{
  "id": "ttc-v2-001",
  "query": "...",
  "family": "product-attribute",
  "partition": "development",
  "authoritative_document_revision_ids": ["sha256:..."],
  "authoritative_chunk_ids": ["sha256:..."],
  "substantial_chunk_ids": ["sha256:..."],
  "misleading_chunk_ids": ["sha256:..."],
  "rationale": "...",
  "expected_abstention": false
}
```

For an unanswerable card, authoritative IDs are empty and
`expected_abstention` is true. Never infer absence from a search result; it
requires a corpus review.

### Adjudication workflow

1. Select a source/evidence family from the manifest.
2. Write a natural query without embedding the target answer verbatim.
3. Inspect the frozen document revision and chunks.
4. Label authoritative, substantial, and misleading evidence.
5. Independently review the label and lock the card's partition.
6. Add the card to immutable evaluation dataset v2 only after review.

### Metrics

Report document MRR/recall@10 plus chunk MRR/nDCG. Break results down by card
family, answerability, and candidate count. The BGE collapse-order experiment
is tuned only on development; the holdout is evaluated after the selected
configuration is frozen.

## Design Decisions

- **110 cards before 185 cards:** accepted. It gives an actionable, manually
  validated intermediate benchmark rather than delaying experiments until all
  labels are complete.
- **Regression is disjoint:** accepted. Regression cards are not a subset of
  development, preventing repeated tuning against known failures.
- **Chunk labels for all v2 cards:** proposed. They are necessary to measure
  reranking and chunking rather than merely correct-document retrieval.

## Alternatives Considered

- Score the 72 proposed cards immediately: rejected because their evidence and
  negative claims are not adjudicated.
- Random card split: rejected because product/FAQ evidence families overlap.
- Tune on holdout: rejected because it destroys the only credible model-choice
  estimate in a small corpus.

## Implementation Plan

1. Reconcile existing pilot, split, chunk-label, and proposed-card drafts.
2. Create an evidence-family manifest and allocate 75/20/15 cards.
3. Adjudicate cards in small batches, with separate reviewer validation.
4. Register frozen immutable dataset v2 and validate every referenced artifact.
5. Run BGE development experiments, freeze a configuration, and score holdout.
6. Expand to 120/40/25 only after the v2 protocol is operating cleanly.

## Risks and Open Questions

- Which reviewer performs independent final label validation?
- How should multi-document comparison questions be graded at chunk level?
- What abstention behavior is available to a retrieval-only experiment?
- How many negative controls can be validated without creating synthetic,
  unrepresentative queries?

## References

- `reference/02-ttc-baseline-evaluation-dataset-v1-candidate-cards.md`
- `scripts/06-ttc-v1-development-holdout-regression-split-draft.md`
- `scripts/07-ttc-chunk-level-evidence-label-proposal.md`
- `pkg/raglab/executor.go`
