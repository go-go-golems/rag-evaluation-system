---
Title: TTC v2 Proposed Split and Chunk-Label Reconciliation
Ticket: RAGEVAL-TTC-LAB-001
Type: Reference
Status: review
Created: 2026-07-16
---

# Reconciliation result

The reranker investigation produced two useful drafts. They are treated as
input to the TTC laboratory ticket, not as frozen truth:

- `RAGEVAL-RERANK-001/scripts/06-ttc-v1-development-holdout-regression-split-draft.md`
  defines the leakage-aware allocation protocol for the existing 20-card pilot
  and the expansion rule for 100+ cards.
- `RAGEVAL-RERANK-001/scripts/07-ttc-chunk-level-evidence-label-proposal.md`
  proposes authoritative, substantial, misleading, and abstention labels for
  all 20 scored cards plus the withheld conflict card against immutable chunk
  set `sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392`.

## Canonical v2 decisions

1. The pilot remains a diagnostic set. It is not enlarged by reusing its
   judgments under a new name.
2. The first v2 foundation targets 110 cards: 75 development, 20 untouched
   holdout, and 15 regression cards. The partitions are assigned to complete
   evidence families, including positive, substantial, and misleading
   documents.
3. Cards v2-001, v2-002, v2-004, v2-005, v2-007, v2-008, v2-010, v2-011,
   v2-014, v2-016, v2-018, and v2-019 have a direct-source adjudication packet
   in `04-ttc-v2-support-faq-adjudication-batch-01.md`.
4. Multi-chunk answers retain contiguous evidence groups. A card is not
   considered adjudicated until every requested facet has an evidence span or
   an explicit `expected_abstention` decision.
5. The withheld cancellation-policy conflict remains outside scoreable data
   until a policy owner resolves precedence; it is a regression/control case,
   not a tuning card.

## Required next artifact

The next implementation step is a machine-readable `ttc-baseline-eval-v2`
manifest containing all 110 cards, evidence-family IDs, source revision IDs,
chunk IDs, labels, reviewer metadata, and a content hash. Until that manifest
exists and passes snapshot compatibility checks, the v2 holdout is not opened.
