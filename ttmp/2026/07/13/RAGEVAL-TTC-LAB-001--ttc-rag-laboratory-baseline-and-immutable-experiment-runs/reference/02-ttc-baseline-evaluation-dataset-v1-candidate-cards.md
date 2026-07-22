---
Title: TTC baseline evaluation dataset v1 candidate cards
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
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://data/ttc-wordpress-rag.sqlite
      Note: Rebuilt source artifact validated by the candidate-card script
    - Path: repo://docs/guides/ttc-data-handbook.md
      Note: Meaning of the source export document kinds and searchable fields
    - Path: repo://ttmp/2026/06/02/RAGEVAL-TTC-SQLITE-EXPORT--export-ttc-wordpress-data-to-sqlite-for-rag-querying/scripts/07-export-ttc-wordpress-to-sqlite.py
      Note: Rebuilds the exact rich SQLite export used for card validation
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/design-doc/02-evaluation-dataset-authoring-and-adjudication-protocol.md
      Note: Defines the review and freezing process for this draft
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/01-validate-ttc-baseline-evaluation-cards.sh
      Note: Read-only source ID, evidence, product predicate, and negative-query validator
ExternalSources: []
Summary: Source-grounded, independently proposed, and mechanically validated draft query cards for the first TTC retrieval evaluation dataset.
LastUpdated: 2026-07-14T16:25:00-04:00
WhatFor: Supply a reviewable authoring draft for ttc-baseline-eval-v1 before source revisions and human adjudication freeze immutable truth.
WhenToUse: Use to implement the dataset-draft loader, conduct human adjudication, select the 200-document snapshot, and build the first retrieval test run.
---


# TTC baseline evaluation dataset v1 candidate cards

## Goal

This is a draft, not a frozen evaluation dataset. It contains 20 candidate information needs with provisional relevance labels, source evidence, and intentional near misses. They were independently proposed from product/fact, care/policy, and editorial corpus reviews, then checked against the rebuilt TTC SQLite export by `scripts/01-validate-ttc-baseline-evaluation-cards.sh`.

The draft deliberately avoids using retrieval rankings to define truth. It uses stable source `doc_id` values only because document revisions are not implemented yet. During implementation, resolve each source ID against `ttc-baseline-v1`, attach exact document revision IDs and rune-range evidence, perform human adjudication, then compile the immutable `ttc-baseline-eval-v1` dataset.

## Context

### Source artifact and validation status

```text
source database:  data/ttc-wordpress-rag.sqlite
source sha256:    c55953ee0d9289577062ac11001c25f63c0286ace45dbc6b4b056c11b0ea6db4
validation date:  2026-07-14
validation state: source IDs, source kinds, evidence anchors, product predicates,
                  and the Bitcoin negative FTS check passed
```

The validator is read-only and does not prove that a label is semantically perfect. It proves that the proposed cards reference real documents of the expected kind, that key evidence phrases still occur in source text, that selected constrained product queries are unique in the exported catalog, and that the explicit Bitcoin negative query has no FTS match.

### Relevance vocabulary

| Level | Rank | Meaning | Binary relevant? |
|---|---:|---|---|
| `0_NOT_RELEVANT` | 0 | Does not satisfy the information need. | no |
| `1_PARTIAL` | 1 | Related but misses a material condition. | no |
| `2_SUBSTANTIAL` | 2 | Materially answers the primary need. | yes |
| `3_AUTHORITATIVE` | 3 | Direct, complete, ideal corpus evidence. | yes |

The frozen dataset must declare `binaryRelevantAtOrAbove: 2_SUBSTANTIAL`. Levels in this document are provisional until human review.

### Candidate-card schema

```yaml
id: ttc-eval-001
status: source-validated-draft
query: "Natural-language information need"
intent: product-attributes
required_facets: [facet-one, facet-two]
provisional_judgments:
  - document_id: wp:123
    level: 3_AUTHORITATIVE
    rationale: "Why the source meets all required facets."
    evidence_anchor: "Stable phrase checked by the validator or review query."
  - document_id: wp:456
    level: 1_PARTIAL
    rationale: "The material missing from an otherwise plausible result."
review_required: [resolve-document-revision, verify-evidence-range, human-adjudication]
```

## Quick reference

### Coverage and status

| Intent stratum | Cards | Notes |
|---|---:|---|
| exact product attributes and constrained discovery | 5 | includes catalog facts, taxonomy, climate, and attribute conjunctions |
| product comparison | 1 | deliberately requires two complementary product pages |
| planting, watering, pruning, and maintenance | 5 | guide/post/FAQ distinctions plus hard negatives |
| editorial explanatory retrieval | 4 | privacy screens, acidic soil, shade, drainage, hardiness |
| FAQ and order policy | 4 | destination restrictions, shipping date, returns, guarantee |
| negative/unanswerable | 1 | no Bitcoin evidence in source FTS |
| withheld due source contradiction | 1 | cancellation policy; excluded from the 20 scored cards |

The 20 scored cards are listed below. The withheld cancellation card remains useful as a future contradiction-retrieval scenario but must not affect v1 metrics until a policy owner establishes source precedence.

## Candidate cards

### Product facts, discovery, comparison, and taxonomy

#### `ttc-eval-001` — Thuja Green Giant catalog attributes

```yaml
query: "What mature height, width, hardiness zone, and sun exposure are catalogued for Thuja Green Giant?"
intent: exact-product-attributes
required_facets: [height, width, hardiness-zone, sunlight]
```

- `3_AUTHORITATIVE` — `wp:3699`, *Thuja Green Giant*: product details record height `20-40`, width `6-12`, zone `5-9`, and `Full Sun to Partial Shade`.
- `0_NOT_RELEVANT` near miss — `wp:3701`, *Leyland Cypress*: privacy-screen overlap, but it is a different product.
- Human review: verify the future importer declares the product-details field family authoritative for these facets.

#### `ttc-eval-002` — constrained Blue Ice Arizona Cypress discovery

```yaml
query: "Find the privacy-tree cypress that is full sun, very drought resistant, hardy in zones 6–9, and matures 15–25 feet tall by 6–8 feet wide."
intent: constrained-product-discovery
required_facets: [privacy-tree, cypress, full-sun, very-drought-resistant, zone-6-9, height-15-25, width-6-8]
```

- `3_AUTHORITATIVE` — `wp:549614`, *Blue Ice Arizona Cypress*: satisfies every declared facet.
- `1_PARTIAL` — `wp:3709`, *Carolina Sapphire Arizona Cypress*: overlaps zone, sun, and drought attributes, but has `25-30` by `10-15` dimensions.
- `1_PARTIAL` — `wp:552438`, *Silver Smoke Arizona Cypress*: has the specified dimensions but lacks the Privacy Trees category.
- Validation: the product predicate in the validator returns exactly one product.

#### `ttc-eval-003` — Blue Italian versus Italian Cypress comparison

```yaml
query: "Compare Blue Italian Cypress with Italian Cypress: mature height, mature width, sunlight, and drought tolerance."
intent: product-comparison
required_facets: [blue-italian-cypress, italian-cypress, height, width, sunlight, drought-tolerance]
```

- `2_SUBSTANTIAL` — `wp:15947`, *Blue Italian Cypress*: supplies one required half of the comparison (height `60-80`, width `4-6`, Full Sun, Very Drought Resistant).
- `2_SUBSTANTIAL` — `wp:3703`, *Italian Cypress*: supplies the other half (height `40-50`, width `5`, Full Sun to Partial Shade, Good Drought Tolerance).
- `0_NOT_RELEVANT` near miss — `wp:3701`, *Leyland Cypress*: plausible cypress screen result, but not either named entity.
- Adjudication note: no single source is expected to be `3_AUTHORITATIVE`; the evaluation unit is document retrieval, not answer synthesis.

#### `ttc-eval-004` — compact Thuja taxonomy and dimensions

```yaml
query: "Which catalogued Thuja or Arborvitae product has both mature height and width of 1–2 feet?"
intent: taxonomy-navigation
required_facets: [thuja-or-arborvitae, height-1-2, width-1-2]
```

- `3_AUTHORITATIVE` — `wp:7347`, *Danica Globe Thuja Arborvitae*: in Arborvitae and Thuja categories; details are `1-2` by `1-2`.
- `1_PARTIAL` — `wp:26028`, *Thuja Can Can*: category match but dimensions are `6-10` by `3-5`.
- Validation: exact category/dimension predicate returns one product.

#### `ttc-eval-005` — wet-soil Cypress discovery

```yaml
query: "Which Cypress Trees product is marked Tolerates Wet Soil and has a mature height of 50–70 feet?"
intent: constrained-product-discovery
required_facets: [cypress-tree, tolerates-wet-soil, height-50-70]
```

- `3_AUTHORITATIVE` — `wp:3717`, *Bald Cypress Tree*.
- `1_PARTIAL` — `wp:10069`, *Red Star White Cypress*: marked Tolerates Wet Soil but is only `6-20` feet tall.
- Validation: exact category/attribute/height predicate returns one product.

### Planting, watering, pruning, and maintenance

#### `ttc-eval-006` — planting-hole geometry

```yaml
query: "How deep and how wide should I dig the hole for a newly planted tree?"
intent: planting-instructions
required_facets: [hole-width, hole-depth-or-root-collar-position]
```

- `3_AUTHORITATIVE` — `wp:812290`, *General Planting Guide*: “3 times as wide as the root ball but not quite as deep,” with root collar at original ground height.
- `1_PARTIAL` — `wp:4131`, *My tree just arrived. What is the proper way to plant it?*: directs the reader to a guide but does not provide measurements.

#### `ttc-eval-007` — first-month watering schedule

```yaml
query: "How often should I water a newly planted tree during its first month?"
intent: planting-aftercare
required_facets: [first-month, watering-frequency]
```

- `3_AUTHORITATIVE` — `wp:812290`, *General Planting Guide*: soak twice per week for the first month, then reduce to once weekly.
- `2_SUBSTANTIAL` — `wp:627148`, *How to Water Your Plants – Part 1, Basic Principles and Watering Outdoors*: root-zone soaking guidance but no TTC first-month cadence.
- The same query appears only once in the dataset; do not include another watered-first-month paraphrase as an independent test.

#### `ttc-eval-008` — fruit-tree pruning time

```yaml
query: "When should I prune my fruit tree, and when should I avoid pruning it?"
intent: seasonal-maintenance
required_facets: [fruit-tree, recommended-pruning-time, avoid-time-or-condition]
```

- `3_AUTHORITATIVE` — `wp:9892`, *How and When To Prune Fruit Trees*: after last frost and before spring growth; avoid fall absent serious disease.
- `2_SUBSTANTIAL` — `wp:28084`, *Useful Tips for Pruning Fruit Trees*: late-winter/annual guidance without the equally important fall qualification.
- `1_PARTIAL` — `wp:751617`, *Pruning Young Trees for Future Strength and Long Life – Part 1, The Basics*: related pruning advice, but not the requested fruit-tree seasonal rule.

#### `ttc-eval-009` — yellow leaves after planting

```yaml
query: "Why are my tree's leaves turning yellow after planting?"
intent: plant-problem-diagnosis
required_facets: [yellow-leaves, cause, post-planting]
```

- `3_AUTHORITATIVE` — `wp:4133`, *Why are the leaves on my tree turning yellow?*: usually excess water, rain, low sites, or insufficient drainage.
- `2_SUBSTANTIAL` — `wp:627148`, *How to Water Your Plants – Part 1*: explains wet soil and root oxygen, but does not identify yellow leaves directly.
- `0_NOT_RELEVANT` near miss — `wp:4134`, *Why are the leaves curling and turning brown?*: opposite watering condition; this is an adversarial diagnostic negative.

#### `ttc-eval-010` — ball-and-burlap procedure

```yaml
query: "How should I handle and plant a ball-and-burlap tree after delivery?"
intent: care-guide
required_facets: [ball-and-burlap, safe-handling, planting-steps]
```

- `3_AUTHORITATIVE` — `wp:398454`, *How To Plant Ball and Burlap Trees*: handling by ropes/strings, root-ball sizing, burlap removal, filling, watering, and mulch.
- `2_SUBSTANTIAL` — `wp:812290`, *General Planting Guide*: shared root-ball and watering process, without B&B-specific handling.
- `0_NOT_RELEVANT` near miss — `wp:405431`, *How to Plant Bare Root Trees*: right broad domain, wrong root treatment.

### Editorial explanatory retrieval

#### `ttc-eval-011` — privacy screen versus hedge

```yaml
query: "What is the planting difference between a privacy screen and a hedge?"
intent: care-guide-comparison
required_facets: [privacy-screen, hedge, spacing-or-growth-form, trimming]
```

- `3_AUTHORITATIVE` — `wp:405509`, *How To Plant a Privacy Screen*: differentiates natural-growing, wider-spaced screens from trimmed hedges.
- `2_SUBSTANTIAL` — `wp:405437`, *How to Plant Evergreen Hedges*: covers hedge planting and trimming but not the complete screen distinction.
- `1_PARTIAL` — `wp:8017`, *Screening and Privacy Trees*: useful selection context, not planting-method comparison.

#### `ttc-eval-012` — acidic soil meaning

```yaml
query: "What does it mean when a plant needs acidic soil, and why does it matter?"
intent: explanatory-taxonomy
required_facets: [acidic-soil, pH, plant-consequence]
```

- `3_AUTHORITATIVE` — `wp:15288`, *What Does “Needs Acidic Soil” Mean?*: pH below 7, acid-loving plants, iron uptake, and chlorosis context.
- `2_SUBSTANTIAL` — `wp:418694`, *How To Plant Rhododendrons, Azaleas and Camellias*: practical acid-loving plant guidance but not the general explanation.
- `0_NOT_RELEVANT` near miss — `wp:224522`, *What is ‘Well-Drained Soil’?*: related soil vocabulary, wrong property.

#### `ttc-eval-013` — shade under evergreens

```yaml
query: "Why is shade under evergreen trees harder for plants than shade from a building?"
intent: explanatory-taxonomy
required_facets: [evergreen-shade, building-shade, causal-difference]
```

- `3_AUTHORITATIVE` — `wp:19387`, *Understanding Shade in the Garden*: distinguishes building, deciduous, evergreen, and deep dry shade; evergreen shade is dense and permanent.
- `2_SUBSTANTIAL` — `wp:9688`, *Shade Loving Plants*: useful shade categories but no evergreen-versus-building explanation.
- `0_NOT_RELEVANT` near miss — `wp:4355`, *Best Shade Trees*: selection content, not shade-condition analysis.

#### `ttc-eval-014` — well-drained soil mechanism

```yaml
query: "Why do plants need well-drained soil—what problem does poor drainage cause?"
intent: explanatory-care-fact
required_facets: [well-drained-soil, root-oxygen-or-pore-air, poor-drainage-consequence]
```

- `3_AUTHORITATIVE` — `wp:224522`, *What is ‘Well-Drained Soil’?*: explains pore air and root oxygen; saturated soil prevents gas exchange.
- `1_PARTIAL` — `wp:812290`, *General Planting Guide*: advises damp rather than soaked soil but does not explain drainage mechanism.

#### `ttc-eval-015` — hardiness-zone interpretation

```yaml
query: "What does a USDA hardiness zone tell me, and which part of ‘8b’ matters most for plant selection?"
intent: explanatory-taxonomy
required_facets: [hardiness-zone, temperature-basis, number-versus-letter]
```

- `3_AUTHORITATIVE` — `wp:4237`, *Plant Hardiness Zone Map*: temperature basis and zone explanation.
- `2_SUBSTANTIAL` — `wp:4116`, *What is my Hardiness Zone?*: directs to the map and explains that the number matters more than the letter in `8b`.

### FAQ, order policy, and support

#### `ttc-eval-016` — destination and citrus shipping restrictions

```yaml
query: "Does The Tree Center ship to California, and can Citrus Trees be shipped to Florida?"
intent: shipping-eligibility
required_facets: [california, citrus, florida, restriction]
```

- `3_AUTHORITATIVE` — `wp:76495`, *Where do you ship to?*: explicitly excludes California and Citrus Trees to Florida.
- `0_NOT_RELEVANT` near miss — `wp:76497`, *Can I arrange a shipping date?*: shipping-related but no destination restriction.

#### `ttc-eval-017` — preferred shipping-time semantics

```yaml
query: "If I choose a preferred shipping time at checkout, is that my delivery date?"
intent: order-scheduling
required_facets: [preferred-shipping-time, ship-date, not-arrival-date]
```

- `3_AUTHORITATIVE` — `wp:76497`, *Can I arrange a shipping date?*: it is the estimated shipping date, not the arrival date.
- `1_PARTIAL` — `wp:456943`, *Delivered and/or Shipped Orders*: shipping/delivery support flow but not chosen-date semantics.

#### `ttc-eval-018` — returns and refund form

```yaml
query: "If I change my mind, how long do I have to return an order and will I get cash back?"
intent: returns-policy
required_facets: [return-timing, refund-form, cash-refund]
```

- `3_AUTHORITATIVE` — `wp:558351`, *Returns & Refunds*: return within seven days of delivery, buyer-paid return shipping, store credit after inspection, no cash refunds.
- `2_SUBSTANTIAL` — `wp:398600`, *Delivery and Returns*: seven-day and store-credit content, with abbreviated wording that requires human policy review.
- `0_NOT_RELEVANT` near miss — `wp:398593`, *Item warranty*: separate warranty process.

#### `ttc-eval-019` — damaged-arrival guarantee

```yaml
query: "My plant arrived damaged or dead. What guarantee applies and what should I provide?"
intent: support-policy
required_facets: [arrival-condition, guarantee, required-evidence, remedy]
```

- `3_AUTHORITATIVE` — `wp:270766`, *Arrive & Thrive Guarantee™*: healthy living arrival, 30-day satisfaction period, store credit, picture and description.
- `2_SUBSTANTIAL` — `wp:4140`, *My tree died.*: routes to the guarantee but omits terms and evidence requirements.

### Explicitly unanswerable control

#### `ttc-eval-020` — Bitcoin payment method

```yaml
query: "Can I pay for a tree with Bitcoin?"
intent: payment-method
answerability: unanswerable-in-declared-corpus
required_facets: [bitcoin-payment-confirmation]
```

- No `2_SUBSTANTIAL` or `3_AUTHORITATIVE` document is expected.
- `1_PARTIAL` — `wp:398551`, *Order by phone*: lists several payment methods but does not make an exhaustive statement that rejects Bitcoin.
- Validation: `documents_fts MATCH 'bitcoin'` returns zero rows in the rebuilt source database.
- Review requirement: the desired answer behavior is calibrated abstention—state that TTC documentation does not confirm Bitcoin payment—rather than asserting that TTC never accepts it.

## Withheld card: policy conflict

### `ttc-eval-withheld-001` — cancellation fee

```yaml
query: "What cancellation fee applies if I cancel an order?"
intent: policy-conflict-detection
status: WITHHELD_PENDING_POLICY_ADJUDICATION
```

- `wp:398597`, *Cancellation policy* (modified 2023), says a 20% fee and cancellation only before fulfillment.
- `wp:4128`, *What if I need to cancel my order?* (modified 2019), describes a 10% fee after one hour.

Do not choose a winner merely by timestamp or let an LLM select a plausible answer. A TTC policy owner must establish precedence. Preserve the conflicting sources and decision rationale in a later dataset version. This card is excluded from all v1 headline metrics.

## Usage examples

### Run source validation

```bash
bash ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/01-validate-ttc-baseline-evaluation-cards.sh
```

Expected terminal state:

```text
PASS: all draft-card document IDs resolve to expected kinds
PASS: all required source evidence phrases are present
PASS: Blue Ice Arizona Cypress constrained discovery is unique
PASS: Danica Globe Thuja dimensions and taxonomy identify one product
PASS: Bald Cypress wet-soil height constraint identifies one product
PASS: Bitcoin has no corpus FTS hit for the explicitly unanswerable card
PASS: TTC baseline evaluation-card source validation completed
```

### Resolve and freeze after implementation

```text
draft cards + ttc-baseline-v1 snapshot
    -> resolve stable source IDs to document revision IDs
    -> validate exact evidence rune slices
    -> create candidate pools from source, SQL, BM25, vector, and hybrid results
    -> blind human adjudication
    -> canonical JSON compilation
    -> sha256 evaluationDatasetId for ttc-baseline-eval-v1
```

No card in this document is eligible for immutable experiment metrics until this sequence completes.

## Related

- `design-doc/02-evaluation-dataset-authoring-and-adjudication-protocol.md` defines the authoring, review, level, pooling, validation, and versioning rules.
- `design-doc/01-ttc-rag-laboratory-baseline-and-immutable-experiment-runs-design-and-implementation-guide.md` defines the broader snapshot, retrieval, run, and UI architecture.
- `scripts/01-validate-ttc-baseline-evaluation-cards.sh` is the read-only source validation used for this draft.

## Quick Reference

<!-- Provide copy/paste-ready content, API contracts, or quick-look tables -->

## Usage Examples

<!-- Show how to use this reference in practice -->

## Related

<!-- Link to related documents or resources -->
