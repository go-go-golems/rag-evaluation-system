---
Title: TTC Expansion Source-ID Validation Result
Ticket: RAGEVAL-TTC-LAB-001
Status: review
Topics:
    - ttc
    - rag-eval
    - corpus
    - evaluation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/07-validate-expansion-source-ids.py
      Note: Read-only validator used for this result
ExternalSources: []
Summary: Mechanical identity validation for the 120-card TTC expansion drafts.
LastUpdated: 2026-07-16T00:00:00Z
WhatFor: Distinguish source-grounded candidates from unvalidated authoring ideas before evidence adjudication.
WhenToUse: Run after changing either expansion draft and before building evidence pools.
---

# Result

The validator scanned the 70-card structured expansion draft and the 50-card
SQLite-grounded audit draft. It found **142 unique `wp:*` source IDs**, and all
142 resolved in `data/ttc-wordpress-rag.sqlite`:

```text
files=2 unique_source_ids=142
resolved=142 missing=0
```

The validator checks only document identity and reports document kind/title.
It does not prove that a source answers a query, that a candidate is relevant,
that a phrase is in the required evidence span, or that two cards are
independent. Those checks remain mandatory before partition assignment and
human adjudication.

## Reproduction

```sh
python3 ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/07-validate-expansion-source-ids.py \
  --db data/ttc-wordpress-rag.sqlite \
  ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/04-ttc-evaluation-expansion-v0-70-proposed-cards.md \
  ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/06-ttc-expansion-audit-and-50-card-source-grounded-draft.md
```

The 120 cards are therefore source-ID-grounded candidates, not 120 frozen
judgments. The next gate is full-text inspection, exact revision/chunk
evidence, evidence-family union, and reviewer adjudication.
