---
Title: TTC Evaluation Expansion v0 - 70 Proposed Cards
Ticket: RAGEVAL-TTC-LAB-001
Status: draft
Topics:
    - rag
    - rag-eval
    - ttc
    - corpus
    - evaluation
    - search
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Source-discovery draft of 70 additional TTC evaluation cards awaiting validation and adjudication.
LastUpdated: 2026-07-16T00:00:00Z
WhatFor: Expand the TTC evaluation queue beyond the pilot without inventing relevance judgments.
WhenToUse: Use for source validation, evidence-family grouping, and human adjudication planning.
---

# TTC evaluation expansion v0 — 70 proposed cards

Status: authoring draft; **not source-validated, adjudicated, registered, or
scoreable**. This supplies a deliberately varied queue toward a 200+ card
program. `expected_source_ids` comes only from the imported snapshot's title
inventory, not from retrieval results. Every card requires source validation,
exact evidence-range review, competing-source review, and human adjudication
before it may become an evaluation judgment.

`leakage_family` groups semantic near-duplicates. Split a family across
development/holdout only as a unit; do not let a paraphrase of the same fact
appear in both.

## Support and policy (14)

```yaml
- id: ttc-expand-001
  query: "How can I check where my order is in the fulfillment process?"
  intent: order-status; expected_source_ids: [wp:4124]; answerability: expected-answerable; difficulty: easy; leakage_family: support-order-status
- id: ttc-expand-002
  query: "What tells me that my order has actually shipped?"
  intent: shipment-confirmation; expected_source_ids: [wp:4126]; answerability: expected-answerable; difficulty: easy; leakage_family: support-shipment-confirmation
- id: ttc-expand-003
  query: "Can I change an order after I place it?"
  intent: order-modification; expected_source_ids: [wp:69806]; answerability: expected-answerable; difficulty: easy; leakage_family: support-order-change
- id: ttc-expand-004
  query: "Can I postpone an order instead of having it sent now?"
  intent: order-scheduling; expected_source_ids: [wp:4127, wp:76497]; answerability: expected-answerable; difficulty: medium; leakage_family: support-shipping-timing
- id: ttc-expand-005
  query: "What should I do if the box contains a different plant from the one I ordered?"
  intent: wrong-item-support; expected_source_ids: [wp:4138, wp:456943]; answerability: expected-answerable; difficulty: medium; leakage_family: support-delivery-problem
- id: ttc-expand-006
  query: "My plants arrived limp. Is that expected and what should I do?"
  intent: arrival-condition; expected_source_ids: [wp:4137, wp:270766]; answerability: expected-answerable; difficulty: medium; leakage_family: support-arrival-condition
- id: ttc-expand-007
  query: "What does pre-sale mean for a tree listing?"
  intent: inventory-status; expected_source_ids: [wp:76498]; answerability: expected-answerable; difficulty: easy; leakage_family: support-presale
- id: ttc-expand-008
  query: "Are shipping and handling charged separately, and how much do they cost?"
  intent: shipping-cost; expected_source_ids: [wp:4121]; answerability: expected-answerable; difficulty: easy; leakage_family: support-shipping-cost
- id: ttc-expand-009
  query: "Which payment methods are accepted by The Tree Center?"
  intent: payment-methods; expected_source_ids: [wp:4120, wp:398551]; answerability: expected-answerable; difficulty: easy; leakage_family: support-payments
- id: ttc-expand-010
  query: "Can I send a tree as a gift?"
  intent: gifting; expected_source_ids: [wp:4122]; answerability: expected-answerable; difficulty: easy; leakage_family: support-gifts
- id: ttc-expand-011
  query: "Does the store offer discounts for bulk or wholesale purchases?"
  intent: wholesale-policy; expected_source_ids: [wp:4123]; answerability: expected-answerable; difficulty: easy; leakage_family: support-wholesale
- id: ttc-expand-012
  query: "Where is The Tree Center located, and can I visit to pick up plants?"
  intent: business-location; expected_source_ids: [wp:76336, wp:398551]; answerability: expected-answerable; difficulty: medium; leakage_family: support-location-pickup
- id: ttc-expand-013
  query: "I need help not covered by the FAQ. What support route should I use?"
  intent: support-escalation; expected_source_ids: [wp:4136, wp:4144, wp:52190]; answerability: expected-answerable; difficulty: medium; leakage_family: support-escalation
- id: ttc-expand-014
  query: "Can I order trees using Ethereum?"
  intent: payment-method; expected_source_ids: [wp:4120, wp:398551]; answerability: unanswerable-in-declared-corpus; difficulty: medium; leakage_family: abstention-unsupported-payment
```

## Planting and care (14)

```yaml
- id: ttc-expand-015
  query: "When should a newly planted tree be staked, and when is staking unnecessary?"
  intent: planting-aftercare; expected_source_ids: [wp:4132, wp:812290]; answerability: expected-answerable; difficulty: medium; leakage_family: care-staking
- id: ttc-expand-016
  query: "How should I plant a bare-root tree after it arrives?"
  intent: planting-instructions; expected_source_ids: [wp:405431, wp:812290]; answerability: expected-answerable; difficulty: medium; leakage_family: care-root-treatment
- id: ttc-expand-017
  query: "How should I plant an evergreen hedge?"
  intent: planting-instructions; expected_source_ids: [wp:405437, wp:405509]; answerability: expected-answerable; difficulty: medium; leakage_family: care-hedge-planting
- id: ttc-expand-018
  query: "What is the correct planting process for Japanese maples?"
  intent: species-planting-guide; expected_source_ids: [wp:418603, wp:9194]; answerability: expected-answerable; difficulty: medium; leakage_family: care-japanese-maple
- id: ttc-expand-019
  query: "How do I plant citrus trees, and are there special establishment steps?"
  intent: species-planting-guide; expected_source_ids: [wp:418611]; answerability: expected-answerable; difficulty: medium; leakage_family: care-citrus
- id: ttc-expand-020
  query: "How should I plant roses?"
  intent: species-planting-guide; expected_source_ids: [wp:418613]; answerability: expected-answerable; difficulty: easy; leakage_family: care-roses
- id: ttc-expand-021
  query: "What planting advice applies specifically to hydrangeas?"
  intent: species-planting-guide; expected_source_ids: [wp:418605]; answerability: expected-answerable; difficulty: easy; leakage_family: care-hydrangea
- id: ttc-expand-022
  query: "How do I plant rhododendrons, azaleas, or camellias?"
  intent: species-planting-guide; expected_source_ids: [wp:418694]; answerability: expected-answerable; difficulty: medium; leakage_family: care-acid-lovers
- id: ttc-expand-023
  query: "When is the right time to prune Leyland Cypress trees?"
  intent: seasonal-maintenance; expected_source_ids: [wp:10577]; answerability: expected-answerable; difficulty: medium; leakage_family: care-leyland-pruning
- id: ttc-expand-024
  query: "How should roses be pruned?"
  intent: seasonal-maintenance; expected_source_ids: [wp:9373, wp:205483]; answerability: expected-answerable; difficulty: medium; leakage_family: care-rose-pruning
- id: ttc-expand-025
  query: "What are the basic principles for watering plants outdoors?"
  intent: watering-guidance; expected_source_ids: [wp:627148]; answerability: expected-answerable; difficulty: easy; leakage_family: care-watering-principles
- id: ttc-expand-026
  query: "How can I transplant an established evergreen tree?"
  intent: transplanting; expected_source_ids: [wp:752143]; answerability: expected-answerable; difficulty: hard; leakage_family: care-evergreen-transplant
- id: ttc-expand-027
  query: "What should I do if I think my tree has a pest or fungus?"
  intent: plant-problem-diagnosis; expected_source_ids: [wp:4135, wp:4146]; answerability: expected-answerable; difficulty: medium; leakage_family: care-pest-diagnosis
- id: ttc-expand-028
  query: "What does it mean if a tree's leaves curl and turn brown?"
  intent: plant-problem-diagnosis; expected_source_ids: [wp:4134]; answerability: expected-answerable; difficulty: medium; leakage_family: care-leaf-symptoms
```

## Editorial concepts and discovery (14)

```yaml
- id: ttc-expand-029
  query: "What makes a tree deciduous?"
  intent: explanatory-taxonomy; expected_source_ids: [wp:9811]; answerability: expected-answerable; difficulty: easy; leakage_family: editorial-deciduous
- id: ttc-expand-030
  query: "How do I choose a privacy arborvitae for my yard?"
  intent: product-discovery; expected_source_ids: [wp:640339, wp:9468, wp:3699]; answerability: expected-answerable; difficulty: medium; leakage_family: editorial-arborvitae-privacy
- id: ttc-expand-031
  query: "Which trees are good choices when I want shade quickly?"
  intent: product-discovery; expected_source_ids: [wp:4355, wp:5111]; answerability: expected-answerable; difficulty: medium; leakage_family: editorial-fast-shade
- id: ttc-expand-032
  query: "Which shrubs are appropriate beside water or in wet ground?"
  intent: plant-selection; expected_source_ids: [wp:635865]; answerability: expected-answerable; difficulty: medium; leakage_family: editorial-wet-soil-selection
- id: ttc-expand-033
  query: "What kinds of trees and shrubs tolerate drought?"
  intent: plant-selection; expected_source_ids: [wp:24391]; answerability: expected-answerable; difficulty: medium; leakage_family: editorial-drought-selection
- id: ttc-expand-034
  query: "How can I make tree fall color brighter?"
  intent: seasonal-care; expected_source_ids: [wp:27314]; answerability: expected-answerable; difficulty: medium; leakage_family: editorial-fall-color
- id: ttc-expand-035
  query: "What garden chores in fall help make spring better?"
  intent: seasonal-care; expected_source_ids: [wp:322338]; answerability: expected-answerable; difficulty: medium; leakage_family: editorial-fall-cleanup
- id: ttc-expand-036
  query: "How should I approach planning a new garden?"
  intent: garden-design; expected_source_ids: [wp:57915, wp:639755]; answerability: expected-answerable; difficulty: medium; leakage_family: editorial-garden-planning
- id: ttc-expand-037
  query: "How can I remove or control bamboo?"
  intent: invasive-plant-control; expected_source_ids: [wp:9823, wp:398536]; answerability: expected-answerable; difficulty: medium; leakage_family: editorial-bamboo-control
- id: ttc-expand-038
  query: "How do I get rid of Japanese beetles?"
  intent: pest-control; expected_source_ids: [wp:9612]; answerability: expected-answerable; difficulty: easy; leakage_family: editorial-japanese-beetles
- id: ttc-expand-039
  query: "What attracts carpenter ants around trees, and how can they be managed?"
  intent: pest-control; expected_source_ids: [wp:9358]; answerability: expected-answerable; difficulty: medium; leakage_family: editorial-carpenter-ants
- id: ttc-expand-040
  query: "How large should I expect a mature tree to become?"
  intent: plant-selection; expected_source_ids: [wp:695361]; answerability: expected-answerable; difficulty: medium; leakage_family: editorial-mature-size
- id: ttc-expand-041
  query: "What are the basic goals and methods of tree trimming?"
  intent: maintenance-explanation; expected_source_ids: [wp:385682, wp:751617]; answerability: expected-answerable; difficulty: medium; leakage_family: editorial-tree-trimming
- id: ttc-expand-042
  query: "Which plants are suited to shade gardens?"
  intent: plant-selection; expected_source_ids: [wp:9688, wp:19387]; answerability: expected-answerable; difficulty: medium; leakage_family: editorial-shade-selection
```

## Product attributes and comparison (18)

```yaml
- id: ttc-expand-043
  query: "What mature size and hardiness zone are listed for American Hornbeam?"
  intent: exact-product-attributes; expected_source_ids: [wp:9410]; answerability: expected-answerable; difficulty: easy; leakage_family: product-american-hornbeam
- id: ttc-expand-044
  query: "What sun and growing-zone requirements are catalogued for Buddha's Hand Citron?"
  intent: exact-product-attributes; expected_source_ids: [wp:18605]; answerability: expected-answerable; difficulty: medium; leakage_family: product-buddhas-hand
- id: ttc-expand-045
  query: "How do Blue Italian Cypress and Italian Cypress differ in mature height?"
  intent: product-comparison; expected_source_ids: [wp:15947, wp:3703]; answerability: expected-answerable; difficulty: medium; leakage_family: product-italian-cypress-comparison
- id: ttc-expand-046
  query: "Which catalogued Cypress is a compact option for a small garden?"
  intent: constrained-product-discovery; expected_source_ids: [wp:10069, wp:552438, wp:549614]; answerability: expected-answerable; difficulty: hard; leakage_family: product-cypress-compact
- id: ttc-expand-047
  query: "What mature size and sun exposure are listed for Burning Bush?"
  intent: exact-product-attributes; expected_source_ids: [wp:3852]; answerability: expected-answerable; difficulty: easy; leakage_family: product-burning-bush
- id: ttc-expand-048
  query: "Which attributes distinguish Danica Globe Thuja Arborvitae from Pancake Arborvitae?"
  intent: product-comparison; expected_source_ids: [wp:7347, wp:68247]; answerability: expected-answerable; difficulty: hard; leakage_family: product-compact-arborvitae
- id: ttc-expand-049
  query: "What size, zone, and light conditions are listed for Fireglow Japanese Maple?"
  intent: exact-product-attributes; expected_source_ids: [wp:374240]; answerability: expected-answerable; difficulty: medium; leakage_family: product-fireglow-maple
- id: ttc-expand-050
  query: "Which product is the Fool Proof Gardenia, and what conditions does it require?"
  intent: exact-product-attributes; expected_source_ids: [wp:547120]; answerability: expected-answerable; difficulty: medium; leakage_family: product-fool-proof-gardenia
- id: ttc-expand-051
  query: "Compare Kleim's Hardy Gardenia with Dwarf Radicans Gardenia for mature size."
  intent: product-comparison; expected_source_ids: [wp:15243, wp:3881]; answerability: expected-answerable; difficulty: hard; leakage_family: product-gardenia-comparison
- id: ttc-expand-052
  query: "What characteristics are listed for Hachiya Japanese Persimmon Tree?"
  intent: exact-product-attributes; expected_source_ids: [wp:66135]; answerability: expected-answerable; difficulty: medium; leakage_family: product-hachiya-persimmon
- id: ttc-expand-053
  query: "Which catalogued Juniper is Emerald Sentinel, and what mature form does it have?"
  intent: exact-product-attributes; expected_source_ids: [wp:531103]; answerability: expected-answerable; difficulty: medium; leakage_family: product-emerald-juniper
- id: ttc-expand-054
  query: "What growing conditions are listed for Hardy Pampas Grass?"
  intent: exact-product-attributes; expected_source_ids: [wp:708509]; answerability: expected-answerable; difficulty: medium; leakage_family: product-pampas-grass
- id: ttc-expand-055
  query: "What makes Lunar Magic Crape Myrtle distinct in the catalog?"
  intent: product-identity; expected_source_ids: [wp:505379]; answerability: expected-answerable; difficulty: medium; leakage_family: product-crape-myrtle-identity
- id: ttc-expand-056
  query: "Compare Pink Pearl Black Diamond and Dynamite Crape Myrtle by flower color and mature size."
  intent: product-comparison; expected_source_ids: [wp:432203, wp:3749]; answerability: expected-answerable; difficulty: hard; leakage_family: product-crape-myrtle-comparison
- id: ttc-expand-057
  query: "What type of site is suitable for Miss Helen American Holly?"
  intent: exact-product-attributes; expected_source_ids: [wp:529598]; answerability: expected-answerable; difficulty: medium; leakage_family: product-miss-helen-holly
- id: ttc-expand-058
  query: "What mature size and hardiness are listed for Rainbow Pillar Serviceberry?"
  intent: exact-product-attributes; expected_source_ids: [wp:779818]; answerability: expected-answerable; difficulty: medium; leakage_family: product-rainbow-serviceberry
- id: ttc-expand-059
  query: "What conditions are listed for Serbian Spruce in tree form?"
  intent: exact-product-attributes; expected_source_ids: [wp:779897]; answerability: expected-answerable; difficulty: medium; leakage_family: product-serbian-spruce
- id: ttc-expand-060
  query: "Which product is appropriate if I want a Thuja privacy screen?"
  intent: constrained-product-discovery; expected_source_ids: [wp:3699, wp:26028, wp:7347, wp:9468]; answerability: expected-answerable; difficulty: hard; leakage_family: product-thuja-privacy
```

## More product, controls, and abstention (10)

```yaml
- id: ttc-expand-061
  query: "What mature attributes are catalogued for Tasty Red Urban Apple Tree?"
  intent: exact-product-attributes; expected_source_ids: [wp:10232]; answerability: expected-answerable; difficulty: medium; leakage_family: product-tasty-red-apple
- id: ttc-expand-062
  query: "How does Wolf River Apple Tree compare with Tasty Red Urban Apple Tree in size?"
  intent: product-comparison; expected_source_ids: [wp:803787, wp:10232]; answerability: expected-answerable; difficulty: hard; leakage_family: product-apple-comparison
- id: ttc-expand-063
  query: "What conditions are listed for Winter Chocolate Heather?"
  intent: exact-product-attributes; expected_source_ids: [wp:762106]; answerability: expected-answerable; difficulty: medium; leakage_family: product-winter-heather
- id: ttc-expand-064
  query: "What makes ZZ Plant suitable or unsuitable for indoor conditions?"
  intent: exact-product-attributes; expected_source_ids: [wp:682105]; answerability: expected-answerable; difficulty: medium; leakage_family: product-zz-plant
- id: ttc-expand-065
  query: "Which plant in this corpus is recommended for a dry garden: Chinese Astilbe or a wet-soil shrub?"
  intent: cross-document-comparison; expected_source_ids: [wp:595960, wp:635865]; answerability: expected-answerable; difficulty: hard; leakage_family: cross-dry-vs-wet
- id: ttc-expand-066
  query: "Can I buy a blueberry bush from The Tree Center?"
  intent: catalog-availability; expected_source_ids: []; answerability: unanswerable-in-declared-corpus; difficulty: medium; leakage_family: abstention-unsupported-catalog-item
- id: ttc-expand-067
  query: "Does The Tree Center offer installation services at my home?"
  intent: service-availability; expected_source_ids: []; answerability: unanswerable-in-declared-corpus; difficulty: medium; leakage_family: abstention-unsupported-service
- id: ttc-expand-068
  query: "Can I pay in Canadian dollars?"
  intent: payment-method; expected_source_ids: [wp:4120, wp:398551]; answerability: unanswerable-in-declared-corpus; difficulty: medium; leakage_family: abstention-unsupported-payment
- id: ttc-expand-069
  query: "What fertilizer should I use for a Thuja Green Giant in its third year?"
  intent: product-specific-care; expected_source_ids: [wp:3699, wp:650812]; answerability: unknown-requires-source-review; difficulty: hard; leakage_family: unknown-product-care
- id: ttc-expand-070
  query: "What is the cancellation fee if I cancel after two hours?"
  intent: policy-conflict-detection; expected_source_ids: [wp:398597, wp:4128, wp:4129]; answerability: conflicting-requires-policy-owner; difficulty: hard; leakage_family: policy-cancellation-conflict
```

## Required next gate

For each selected card: (1) inspect each expected source and any competing
source, (2) record exact revision and rune/chunk ranges, (3) add named
relevance judgments and adversarial near misses, (4) validate against the
declared snapshot, and (5) have a reviewer adjudicate before dataset
registration. Cards `014`, `066`–`068`, and `070` are controls: preserve the
requested abstention/conflict behavior rather than fabricating a positive
truth label.
