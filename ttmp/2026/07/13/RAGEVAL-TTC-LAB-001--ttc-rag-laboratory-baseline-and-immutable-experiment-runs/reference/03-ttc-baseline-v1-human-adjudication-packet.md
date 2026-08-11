---
Title: TTC baseline v1 human adjudication packet
Ticket: RAGEVAL-TTC-LAB-001
Status: active
Topics: [rag, evaluation, ttc, adjudication, intern-guide]
DocType: reference
Intent: long-term
Summary: Review packet for converting the 20 source-validated TTC candidate cards into a human-approved immutable evaluation dataset.
LastUpdated: 2026-07-14T21:52:00-04:00
---

# TTC baseline v1 human adjudication packet

## Purpose and decision boundary

This packet is the authority gate between a useful candidate evaluation set and a benchmark that may be called `ttc-baseline-eval-v1`. The candidate cards were written from TTC source material, parsed into machine-readable named relevance judgments, and exercised against real immutable retrieval artifacts. They are not yet human-approved truth.

A TTC policy owner or editorial reviewer must make the decisions recorded here. The reviewer does not need to evaluate BM25, embedding, or RRF algorithms. Their responsibility is narrower: establish what TTC documentation authoritatively supports, identify which documents materially answer each information need, and record any source-precedence decision. The reviewer should work from the fixed source revisions in the corpus snapshot, not from live WordPress pages.

## Frozen inputs to inspect

| Input | Immutable identifier |
| --- | --- |
| TTC corpus snapshot | `sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409` |
| fixed chunk set | `sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392` |
| BM25 artifact | `sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691` |
| 768D embedding set | `sha256:2665c5249b8352ce6904fc00c934534dd179f3eeef0a6a75429a9034be0e03e0` |
| candidate card source | `reference/02-ttc-baseline-evaluation-dataset-v1-candidate-cards.md` |

The candidate metric report uses binary relevance at `2_SUBSTANTIAL` or above. It preserves the ordinal levels below for future nDCG and qualitative inspection.

| Stored level | Review meaning |
| --- | --- |
| `0_NOT_RELEVANT` | Does not answer the information need or is materially misleading. |
| `1_PARTIAL` | Topic-adjacent but misses a material condition or answer part. |
| `2_SUBSTANTIAL` | Materially useful answer with an important limitation, indirectness, or incomplete facet. |
| `3_AUTHORITATIVE` | Directly, fully, and appropriately answers the information need from TTC content. |

## Reviewer workflow

1. Open the candidate card and its named TTC source documents from the imported snapshot.
2. Confirm the required facets are actually present in the quoted source text, not inferred from general horticultural knowledge.
3. Approve, revise, or reject each named grade. Record an exact evidence excerpt and the immutable `document_revision_id` when approving a grade at `2_SUBSTANTIAL` or `3_AUTHORITATIVE`.
4. For each `0_NOT_RELEVANT` or `1_PARTIAL` near miss, record the decisive missing or contradictory facet.
5. Confirm that the card’s query is a distinct information need, not an accidental paraphrase of another card.
6. Sign the decision record. A changed query, source revision, evidence span, grade, or relevance threshold becomes a new dataset version; it does not alter v1 after publication.

### Required decision record

```yaml
card_id: ttc-eval-001
reviewer: "name and TTC role"
reviewed_at: "RFC3339 timestamp"
decision: approved | revised | rejected | withheld
relevance_threshold: 2_SUBSTANTIAL
judgments:
  - stable_document_id: wp:3699
    document_revision_id: sha256:...
    level: 3_AUTHORITATIVE
    evidence_quote: "..."
    source_start_runes: 0
    source_end_runes: 0
    rationale: "..."
policy_precedence: null
notes: "..."
```

## Card-by-card checklist

Mark one decision per row and add the full decision record above or in the future adjudication UI. The corresponding candidate-card section contains the source evidence and current proposed labels.

| Card | Information need | Proposed authoritative/substantial source IDs | Reviewer decision |
| --- | --- | --- | --- |
| `ttc-eval-001` | Thuja Green Giant attributes | `wp:3699` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-002` | constrained Blue Ice Cypress discovery | `wp:549614` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-003` | Blue Italian vs Italian Cypress comparison | `wp:15947`, `wp:3703` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-004` | compact Thuja taxonomy/dimensions | `wp:7347` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-005` | wet-soil Cypress discovery | `wp:3717` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-006` | planting-hole geometry | `wp:812290` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-007` | first-month watering | `wp:812290`, `wp:627148` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-008` | fruit-tree pruning time | `wp:9892`, `wp:28084` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-009` | yellow leaves after planting | `wp:4133`, `wp:627148` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-010` | ball-and-burlap planting | `wp:398454`, `wp:812290` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-011` | privacy screen versus hedge | `wp:405509`, `wp:405437` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-012` | acidic-soil explanation | `wp:15288`, `wp:418694` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-013` | shade under evergreens | `wp:19387`, `wp:9688` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-014` | well-drained-soil mechanism | `wp:224522` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-015` | hardiness-zone interpretation | `wp:4237`, `wp:4116` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-016` | destination/citrus shipping restrictions | `wp:76495` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-017` | preferred shipping-time semantics | `wp:76497` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-018` | returns and refund form | `wp:558351`, `wp:398600` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-019` | damaged-arrival guarantee | `wp:270766`, `wp:4140` | ☐ approve ☐ revise ☐ reject |
| `ttc-eval-020` | Bitcoin payment abstention | none; `wp:398551` is contextual only | ☐ approve abstention rule ☐ revise ☐ reject |

## Explicitly withheld cancellation-policy card

`ttc-eval-withheld-001` remains outside all v1 headline metrics. The snapshot contains an apparent conflict:

- `wp:398597`, *Cancellation policy*, says 20% and cancellation only before fulfillment.
- `wp:4128`, *What if I need to cancel my order?*, says 10% after one hour.

Do not resolve this by choosing the newer timestamp, the higher retrieval rank, or a plausible model answer. A TTC policy owner must supply a precedence decision and a rationale. After that decision, create a new immutable dataset version that records both source revisions and the policy decision; do not silently modify the candidate dataset.

## Completion criteria

The reviewer may authorize dataset publication only when all of the following are true:

- All 20 scored cards have a signed decision record.
- Every positive judgment names the resolved immutable document revision and an exact evidence slice.
- `ttc-eval-020` has an approved abstention expectation rather than a fabricated negative policy claim.
- The cancellation conflict is still withheld or has a recorded TTC policy-precedence decision.
- The canonical dataset manifest states `binaryRelevantAtOrAbove: 2_SUBSTANTIAL` and includes the four named levels.
- The resulting dataset receives a new SHA-256 ID and status `human-reviewed-published`.

Until then, APIs, runs, and reports must continue to display `candidate-source-validated-not-human-frozen`.
