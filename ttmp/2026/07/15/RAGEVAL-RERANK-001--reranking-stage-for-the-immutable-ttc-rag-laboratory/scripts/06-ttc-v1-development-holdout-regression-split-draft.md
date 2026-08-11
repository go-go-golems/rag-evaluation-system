---
Title: TTC v1 development, holdout, and regression split draft
Ticket: RAGEVAL-RERANK-001
Status: draft
DocType: script-notes
LastUpdated: 2026-07-16
---

# TTC v1 development, holdout, and regression split draft

## Decision

Treat the current 20 cards as a **small, source-grouped pilot**, not as a
statistically decisive benchmark.  The allocation below is deterministic and
uses no retrieval-run outcomes.  It is designed to prevent the most direct
document-evidence leakage while preserving the currently known difficult
behaviours.

| Partition | Role | Existing card IDs | Count |
| --- | --- | --- | ---: |
| Development | Select retrieval, candidate-budget, collapse, and reranker configurations; inspect traces. | `001`, `002`, `004`, `005`, `006`, `007`, `010` | 7 |
| Holdout | Compare a pre-registered, frozen shortlist of configurations. Do not inspect per-card outputs until selection is frozen. | `003`, `008`, `009`, `011`, `012`, `013`, `014`, `015` | 8 |
| Regression | Must remain passing after a configuration is selected; use to catch known policy, support, and abstention regressions, not to tune. | `016`, `017`, `018`, `019`, `020` | 5 |

All IDs above are shorthand for `ttc-eval-NNN`. `ttc-eval-withheld-001` is
outside every partition until an owner records a policy-precedence decision.

## Why this allocation

The cards are grouped before splitting by their named positive and material
near-miss evidence documents, rather than randomly by query text. This avoids
one query revealing the exact answer document or wording for another query.

| Evidence/query family | Cards | Partition | Rationale |
| --- | --- | --- | --- |
| Constrained catalog facts and discovery | `001`, `002`, `004`, `005` | Development | Provides product metadata, taxonomy, and conjunctional discovery while tuning candidate construction. Each has a different primary product, so this is diversity within the development slice. |
| General planting guide / root-ball handling | `006`, `007`, `010` | Development | All depend materially on `wp:812290`; keeping the family together prevents guide-text leakage. |
| Cypress comparison | `003` | Holdout | Tests multi-document complementary retrieval without exposing this exact comparison during tuning. |
| Seasonal diagnosis and care | `008`, `009` | Holdout | Tests temporal qualification and adversarial opposite-condition retrieval. |
| Editorial explanatory concepts | `011`, `012`, `013`, `014`, `015` | Holdout | Tests explanatory rather than catalog matching. `012` and `014` are together because `wp:224522` is a named near miss for the former and authoritative for the latter. |
| Shipping, returns, guarantee, and abstention | `016`, `017`, `018`, `019`, `020` | Regression | These capture high-value operational failure modes: policy precision, date semantics, evidence requirements, and calibrated no-answer behavior. |

The partition therefore includes every currently scored intent family, but it
does **not** claim that every stratum is represented in every partition. With
only 20 cards, enforcing that would either leak source families or give a
misleading impression of sampling power.

## Non-negotiable evaluation protocol

1. Freeze corpus snapshot, chunks, evaluation judgments, embedding artifact,
   candidate text/truncation policy, and metric definitions for a comparison.
2. Use only Development cards to choose weights, channel top-K, candidate
   budget, collapse placement, reranker, reranker prompt/input form, and any
   query rewriting policy. Record every attempted configuration and decision.
3. Before opening Holdout results, write a short selection record naming the
   one configuration (or a bounded predeclared shortlist) and the tie-breaker:
   primary metric, secondary metric, latency cap, and abstention rule.
4. Execute Holdout once for that decision. Report aggregate and per-intent
   results, confidence intervals or bootstrap intervals where practical, and
   every query trace. A Holdout result cannot be used to retune the selected
   configuration.
5. Run Regression after Holdout as a release gate. It may diagnose a failure,
   but cannot be silently optimized against. A fix selected using its results
   starts a new evaluation cycle with a newly frozen holdout version.
6. Preserve `020` separately in aggregate reporting: standard retrieval MRR
   excludes it because it has no relevant document; report abstention accuracy,
   false-supported-answer rate, and whether irrelevant citations were emitted.

## Leakage risks and controls

| Risk | Existing example | Control |
| --- | --- | --- |
| Same authoritative document in multiple cards | `006`, `007`, and `010` use `wp:812290`. | Keep the full evidence family in Development; future cards sharing a positive or material near-miss document go in the same partition. |
| A near miss becomes another card's answer | `wp:224522` is a `012` near miss and `014` authority. | Group by every judged document at level 1–3, not positives alone. |
| Same policy domain, distinct FAQ wording | `016` and `017` are shipping FAQs. | Keep both in Regression; do not tune shipping-specific prompts on either. |
| Query paraphrase leakage | first-month watering paraphrases could merely test memorized wording. | Deduplicate by information need and evidence family before assignment; paraphrases stay in the same partition and count as one family in headline summaries. |
| Corpus/annotation leakage | a revised document or label changes only one partition. | Version all document revisions and judgments together; any changed evidence span, grade, or query creates a new dataset manifest and reruns all partitions. |
| Repeated holdout peeking | examining individual misses drives informal tuning. | Restrict detailed Holdout traces until pre-registration; log access and invalidate the holdout after a tuning-informing review. |

## Protocol for the next 100+ cards

Create a new immutable dataset version rather than extending this split in
place. For each proposed card, capture intent, answerability, required facets,
source kind, named positive/partial/negative documents, exact evidence spans,
and an `evidence_family_id` computed from all judged documents at levels 1–3.

1. Human-adjudicate and freeze cards before retrieval experiments; candidate
   pools can find possible judgments but cannot define truth from rank.
2. Form indivisible groups by shared `evidence_family_id`; also group direct
   paraphrases, the same product entity, and policy pages governed by one
   precedence decision.
3. Stratify groups by: intent class, answerability, source kind, cardinality
   (one-document versus multi-document), query difficulty, and adversarial
   near-miss presence. Assign groups deterministically using a recorded seed.
4. Target approximately 60% Development, 20% Holdout, and 20% Regression only
   after each important stratum has at least five independent evidence groups.
   Until then, report the stratum as exploratory rather than pretending the
   split is balanced.
5. Keep a final untouched confirmation set if model selection will be reported
   externally. Regression is a safety suite, not an unbiased confirmation set.

## Current conclusion

The present set can support disciplined reranker and collapse-order exploration
on Development and a single guarded sanity check on Holdout. It cannot support
a durable claim that one reranker is generally superior. The next expansion
should add independent evidence families first—especially product near-neighbor
discovery, support-policy ambiguity, multi-chunk answers, and unanswerable
queries—then create `ttc-baseline-eval-v2` with this grouping rule applied
before any ranking runs.
