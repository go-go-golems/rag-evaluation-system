---
Title: TTC v2 240-card partition and leakage-audit protocol
Ticket: RAGEVAL-TTC-LAB-001
Status: draft
Topics: []
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/design-doc/03-ttc-evaluation-corpus-v2-foundation-and-adjudication-protocol.md
      Note: Defines the preceding 110-card staged adjudication milestone
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/05-ttc-v2-proposed-split-and-chunk-label-reconciliation.md
      Note: Establishes existing v2 split and evidence reconciliation decisions
ExternalSources: []
Summary: Exact 240-card partition, evidence-family grouping, stratum floors, and machine-checkable audit contract for the TTC evaluation corpus.
LastUpdated: 2026-07-16T00:00:00Z
WhatFor: ""
WhenToUse: ""
---


# TTC v2 240-card partition and leakage-audit protocol

## Decision

The next publishable TTC benchmark target is **240 adjudicated cards**, not an
incremental extension of the 20-card pilot or the provisional 110-card
foundation. It is a draft authoring and audit contract; it freezes no labels,
queries, source revisions, or partitions.

| Partition | Cards | Permitted use |
| --- | ---: | --- |
| Development | 144 | Configuration selection, error analysis, candidate-budget and collapse-policy tuning. |
| Holdout | 48 | One pre-registered comparison of the selected configuration; no tuning after inspection. |
| Regression | 48 | Release gate for known failures; diagnosis only, never a tuning pool. |
| Total | **240** | All cards are assigned exactly once. |

There are **30 abstention cards**. Abstention is an answerability label, not a
fourth partition: 18 are Development, 6 Holdout, and 6 Regression. This keeps
answerability visible in every evaluation role without turning a negative set
into a side benchmark that never exercises release gating.

## Required stratum matrix

Each row contains 30 cards, allocated 18/6/6 across Development/Holdout/
Regression. A card has exactly one primary stratum; secondary tags may record
multi-hop, long-context, or policy-conflict characteristics.

| Primary stratum | Dev | Holdout | Regression | Total | Minimum independent evidence families |
| --- | ---: | ---: | ---: | ---: | ---: |
| Exact product attributes and taxonomy | 18 | 6 | 6 | 30 | 24 |
| Constrained product discovery / near-neighbor discrimination | 18 | 6 | 6 | 30 | 24 |
| Multi-document product comparison | 18 | 6 | 6 | 30 | 24 |
| Planting and maintenance procedure | 18 | 6 | 6 | 30 | 24 |
| Diagnosis and adversarial symptom distinction | 18 | 6 | 6 | 30 | 24 |
| Editorial explanatory retrieval | 18 | 6 | 6 | 30 | 24 |
| Commerce, shipping, return, guarantee, and support policy | 18 | 6 | 6 | 30 | 24 |
| Unanswerable / calibrated abstention | 18 | 6 | 6 | 30 | 30 |
| **Total** | **144** | **48** | **48** | **240** | — |

The family floor means that no primary stratum may contain more than two cards
from one evidence family, and those cards must be in the same partition. For
the abstention stratum, each card must have its own independently documented
absence review; there is no shared "no result" family.

## Evidence-family assignment rules

An evidence family is the indivisible unit of assignment. Construct it before
partitioning using union-find over candidate cards. Union two cards when any
one condition holds:

1. They share a document revision, source document, or chunk at relevance
   grade `1_PARTIAL`, `2_SUBSTANTIAL`, or `3_AUTHORITATIVE`.
2. They ask for the same product entity, product variant, policy rule,
   customer workflow, or editorial article, even if their initially named
   evidence differs.
3. They are paraphrases, answer the same information need, or differ only in
   numerical/seasonal values that are stated in the same evidence passage.
4. They rely on the same source-precedence decision, including a conflict
   resolution or an exhaustive-corpus absence review.
5. A reviewer declares a shared latent answer after inspecting the immutable
   snapshot. This declaration must be recorded with a rationale.

The transitive closure is intentional: if card A shares a substantial document
with B and B shares a misleading chunk with C, all three remain together. A
family may be larger than two cards; it receives one partition, and it counts
once against that partition's independent-family coverage.

Do **not** union cards merely because both are `product`, `FAQ`, or from the
same broad taxonomy. Those are stratification attributes, not evidence
leakage. Conversely, do not split a multi-document comparison: its complete
set of supporting documents and chunks is one family.

## Partition procedure

1. Freeze an authoring snapshot and schema version. Adjudicate candidate
   evidence before assignment; retrieval ranks may expand a judgment pool but
   never establish truth.
2. Build family components using the rules above. Reject any card missing an
   immutable document revision, exact evidence span, or documented absence
   review.
3. Assign each component a primary stratum and difficulty tags. Use a recorded
   deterministic seed to stratify components into 60%/20%/20% card targets.
   Optimise card totals subject to the invariant that components never split.
4. Human-review the resulting matrix before exposing any Holdout query/result
   trace. The reviewer checks semantic family leakage, not just IDs.
5. Write the canonical manifest, hash it, and lock partition assignment. A
   changed query, source revision, evidence span, relevance grade, family edge,
   or partition creates a new dataset version.

If a component makes exact totals impossible, preserve component integrity and
record the variance. A partition may differ by at most two cards from its
target only with an audit exception naming the component and its size. Do not
solve this by moving a connected card across partitions.

## Minimum coverage and quality floors

Before v2 can be called evaluation-ready, all conditions below must hold:

- Every primary stratum meets its exact 18/6/6 card allocation unless an
  approved component-size exception is recorded.
- Every stratum has at least 24 independent families, except abstention (30).
- Each partition contains at least five independent families per answerable
  stratum and six independent abstention reviews. A six-card slice may only
  fail this where a documented two-card family makes five impossible.
- At least 25% of every answerable stratum has a named adversarial or
  materially misleading chunk/document; at least 25% requires two or more
  authoritative/substantial chunks or documents. These tags can overlap.
- Every positive card names one authoritative evidence group; every requested
  facet maps to an authoritative or substantial evidence span. Multi-document
  cards must name all required components.
- Every abstention card records corpus scope, absence-review method, reviewed
  sources/queries, and the permitted response: no unsupported factual denial.
- Policy-conflict cards remain `WITHHELD` until owner precedence exists. They
  do not fill a policy or regression quota merely because they are difficult.

## Machine-checkable manifest and audit checklist

The future canonical JSON manifest must include, at minimum:

```json
{
  "dataset_id": "ttc-baseline-eval-v2",
  "schema_version": "rag-eval-dataset/v2",
  "corpus_snapshot_id": "sha256:...",
  "cards": [{
    "id": "ttc-v2-001",
    "partition": "development",
    "primary_stratum": "product-attributes",
    "answerability": "answerable",
    "evidence_family_id": "ef-...",
    "source_document_ids": ["wp:..."],
    "document_revision_ids": ["sha256:..."],
    "evidence_groups": [{"grade": "3_AUTHORITATIVE", "chunk_ids": ["sha256:..."], "facets": ["..."]}],
    "misleading_document_revision_ids": [],
    "review": {"reviewer": "...", "reviewed_at": "..."}
  }]
}
```

The validator must fail with a non-zero exit if any check fails:

```text
[ ] schema and dataset identifiers are present; 240 unique card IDs exist
[ ] partition counts are development=144, holdout=48, regression=48
[ ] each primary-stratum/partition cell is 18/6/6 and each row totals 30
[ ] answerability=abstention count is 30 and partition counts are 18/6/6
[ ] every card has exactly one partition, one primary stratum, and one family
[ ] every family appears in exactly one partition
[ ] union-find recomputation from all grade 1–3 document/chunk IDs produces no cross-partition component
[ ] no card has a direct duplicate query or declared paraphrase in another partition
[ ] all revision IDs and chunk IDs belong to corpus_snapshot_id and each chunk belongs to its stated revision
[ ] every positive required facet has at least one evidence group at grade 2 or 3
[ ] each answerable card has >=1 grade-3 evidence group; each multi-document card has all required components
[ ] each abstention card has zero positive evidence groups, a scope-bound absence review, and abstention response policy
[ ] no WITHHELD/PENDING_POLICY_PRECEDENCE card is scored or counted toward quota
[ ] per-stratum family floors and five-family-per-partition floors are met, or named exceptions validate
[ ] adversarial and multi-evidence tag floors are met
[ ] reviewer identity, review time, rationale, source spans, and canonical manifest SHA-256 are present
```

The audit output should emit a stable `families.csv` (family, cards,
documents, chunks, partition, reason edges) and `split-audit.json` containing
the count matrix, exceptions, and manifest hash. A reviewer must be able to
trace every family edge to a source field or written semantic-family rationale.

## Evaluation conduct after audit

Development results may be inspected freely. Before running Holdout, record
the candidate configurations, primary metric, secondary metric, latency cap,
abstention metric, and tie-breaker. Open the Holdout aggregate once; per-card
inspection that informs a change invalidates it as a selection holdout.

Regression is a release gate. It reports known failure preservation separately
from model-selection quality. Any fix informed by Holdout or Regression begins
a new dataset/evaluation cycle rather than silently reusing a contaminated
partition.

## Relationship to the staged foundation

The documented 110-card (75/20/15) proposal remains an operational
adjudication milestone. This 240-card contract is the first benchmark-sized
target: batches may be authored and reviewed toward it, but no intermediate
batch should be described as the final v2 benchmark or alter the frozen pilot.
