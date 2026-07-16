---
Title: TTC Expansion Candidate Dataset Registration
Ticket: RAGEVAL-TTC-LAB-001
Status: active
Topics:
    - ttc
    - rag-eval
    - evaluation
    - database
    - corpus
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/09-register-ttc-expansion-candidate-dataset.py
      Note: Snapshot-bound candidate dataset registrar
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/09-ttc-expansion-v0-manifest.json
      Note: Reproducible canonical manifest emitted by the registrar
ExternalSources: []
Summary: Registers the source-grounded expansion queue as an immutable candidate dataset for laboratory runs.
LastUpdated: 2026-07-16T00:00:00Z
WhatFor: Run experiments against the 148 snapshot-compatible expansion cards without changing v1.
WhenToUse: Use the candidate dataset for exploratory development runs only; do not open holdout claims from it.
---

# Registration result

The ticket-local registrar created dataset `candidate:ttc-expansion-v0` in the
local `data/rag-eval.db` catalog using the immutable TTC snapshot
`sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409`.
The dataset contains **148 cards** from the 150 structured expansion cards.

Two known policy-conflict cards were excluded because their source
`wp:398597` is intentionally absent from the current snapshot:

- `ttc-expand-070`
- `ttc-y-074`

Those cards remain in the authoring queue and must not be silently converted
into ordinary relevance judgments. The emitted manifest is committed beside
the registrar; the SQLite database is a local generated artifact and remains
ignored by Git.

## Reproduction

```sh
PYTHONDONTWRITEBYTECODE=1 python3 scripts/09-register-ttc-expansion-candidate-dataset.py \
  --db data/rag-eval.db \
  --dataset-id candidate:ttc-expansion-v0 \
  --exclude-card ttc-expand-070 \
  --exclude-card ttc-y-074 \
  --manifest-out scripts/09-ttc-expansion-v0-manifest.json \
  reference/04-ttc-evaluation-expansion-v0-70-proposed-cards.md \
  reference/05-ttc-evaluation-expansion-y-v0-80-proposed-cards.md
```

This is an immutable **candidate** dataset. It is suitable for development
retrieval experiments, but it is not a reviewer-frozen holdout or regression
benchmark. Evidence-family assignment, exact chunk spans, and adjudicated
grades remain required before publication.
