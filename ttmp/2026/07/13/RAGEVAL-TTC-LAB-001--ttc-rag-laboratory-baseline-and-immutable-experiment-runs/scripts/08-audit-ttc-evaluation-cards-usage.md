---
Title: TTC Evaluation-Card Audit Utility Usage
Ticket: RAGEVAL-TTC-LAB-001
Status: active
Topics:
    - rag-eval
    - ttc
    - evaluation
    - scripting
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Usage guide for the read-only TTC candidate-card and future partition audit utility.
LastUpdated: 2026-07-16T00:00:00Z
WhatFor: Validate draft card syntax before evidence adjudication.
WhenToUse: Run after changing candidate-card files or adding a machine-readable split manifest.
---

# TTC evaluation-card audit usage

`08-audit-ttc-evaluation-cards.py` is read-only. It checks draft card IDs,
parseable queries, and declared source IDs. When given a future JSON metadata
manifest, it additionally checks that every draft card has one partition and
evidence family, and that no family or declared source document crosses a
partition.

It does not validate source truth, change judgments, assign partitions, or
freeze labels.

```bash
python3 ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/08-audit-ttc-evaluation-cards.py \
  --cards ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/02-ttc-baseline-evaluation-dataset-v1-candidate-cards.md \
  --cards ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/04-ttc-evaluation-expansion-v0-70-proposed-cards.md
```

When an authoring manifest exists, add `--metadata path/to/manifest.json`.
The manifest must have a top-level `cards` array. Every record needs `id`,
`partition`, `evidence_family_id`, and `source_document_ids` (or the draft
alias `expected_source_ids`). Valid partitions are `development`, `holdout`,
and `regression`.

The script deliberately stops short of the full 240-card quota and evidence
span checks in the canonical reference. Those require reviewer-approved
revision/chunk labels and belong in the future dataset-manifest validator.
