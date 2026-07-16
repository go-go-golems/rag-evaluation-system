---
Title: TTC Evaluation Expansion Y v0 - 80 Proposed Cards
Ticket: RAGEVAL-TTC-LAB-001
Status: draft
Topics:
    - rag
    - rag-eval
    - ttc
    - evaluation
    - search
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Eighty additional unlabelled TTC candidate cards covering comparisons, climate, diagnostics, procedures, and controls.
LastUpdated: 2026-07-16T00:00:00Z
WhatFor: Expand source-grounded query coverage toward the 240-card benchmark target.
WhenToUse: Use for source validation, evidence-family grouping, and adjudication planning.
---

# TTC evaluation expansion y v0 — 80 proposed cards

Status: **authoring-only draft**. These are candidate information needs, not
judgments: expected source IDs came from the snapshot inventory; each requires
source/evidence validation, competing-source review, and human adjudication.
No card may enter a registered or frozen evaluation set from this document
alone. `leakage_family` is a split constraint: keep every member of a family
on one side of a development/holdout split.

## Multi-document comparisons (y-001–y-016)
```yaml
- {id: ttc-y-001, query: "Compare planting apple/pear trees with peach/nectarine trees.", intent: cross-guide-comparison, expected_source_ids: [wp:9210, wp:9208, wp:418609, wp:398553], answerability: expected-answerable, difficulty: hard, leakage_family: compare-fruit-planting}
- {id: ttc-y-002, query: "Compare evergreen hedges and deciduous hedges for planting and upkeep.", intent: cross-guide-comparison, expected_source_ids: [wp:405437, wp:405441], answerability: expected-answerable, difficulty: hard, leakage_family: compare-hedge-types}
- {id: ttc-y-003, query: "How do ball-and-burlap and bare-root planting procedures differ?", intent: cross-guide-comparison, expected_source_ids: [wp:398454, wp:405431], answerability: expected-answerable, difficulty: hard, leakage_family: compare-root-form}
- {id: ttc-y-004, query: "Compare planting Japanese maples with planting citrus trees.", intent: cross-guide-comparison, expected_source_ids: [wp:418603, wp:418611], answerability: expected-answerable, difficulty: hard, leakage_family: compare-species-planting}
- {id: ttc-y-005, query: "Compare Blue Ice and Carolina Sapphire Arizona Cypress for site requirements.", intent: product-comparison, expected_source_ids: [wp:549614, wp:3709], answerability: expected-answerable, difficulty: hard, leakage_family: compare-arizona-cypress}
- {id: ttc-y-006, query: "Compare Bald Cypress and Serbian Spruce for wet or dry site tolerance.", intent: product-comparison, expected_source_ids: [wp:3717, wp:779897], answerability: expected-answerable, difficulty: hard, leakage_family: compare-moisture-tolerance}
- {id: ttc-y-007, query: "How do Thuja Green Giant and Leyland Cypress differ as privacy screens?", intent: product-comparison, expected_source_ids: [wp:3699, wp:3701, wp:9468], answerability: expected-answerable, difficulty: hard, leakage_family: compare-privacy-evergreen}
- {id: ttc-y-008, query: "Compare Black Republican and Black Tartarian cherry trees.", intent: product-comparison, expected_source_ids: [wp:647077, wp:814274], answerability: expected-answerable, difficulty: hard, leakage_family: compare-cherry}
- {id: ttc-y-009, query: "Compare Scots Pine and Dwarf Mountain Pine for mature form.", intent: product-comparison, expected_source_ids: [wp:802558, wp:742129], answerability: expected-answerable, difficulty: hard, leakage_family: compare-pine-form}
- {id: ttc-y-010, query: "Compare Cityline Vienna and Berry White hydrangeas.", intent: product-comparison, expected_source_ids: [wp:710014, wp:708559, wp:308336], answerability: expected-answerable, difficulty: hard, leakage_family: compare-hydrangea}
- {id: ttc-y-011, query: "Compare Nandina varieties Compact and Tuscan Flame.", intent: product-comparison, expected_source_ids: [wp:31506, wp:644492], answerability: expected-answerable, difficulty: hard, leakage_family: compare-nandina}
- {id: ttc-y-012, query: "Compare Rose Spirea and White Gold Spirea for garden use.", intent: product-comparison, expected_source_ids: [wp:646816, wp:545357], answerability: expected-answerable, difficulty: hard, leakage_family: compare-spirea}
- {id: ttc-y-013, query: "Compare serviceberry products Rainbow Pillar and Serviceberry.", intent: product-comparison, expected_source_ids: [wp:779818, wp:9413], answerability: expected-answerable, difficulty: hard, leakage_family: compare-serviceberry}
- {id: ttc-y-014, query: "Compare hydrangea planting guidance with rose planting guidance.", intent: cross-guide-comparison, expected_source_ids: [wp:418605, wp:418613], answerability: expected-answerable, difficulty: hard, leakage_family: compare-shrub-planting}
- {id: ttc-y-015, query: "How do screening/privacy trees differ from shade trees as a selection goal?", intent: editorial-comparison, expected_source_ids: [wp:8017, wp:4355], answerability: expected-answerable, difficulty: hard, leakage_family: compare-selection-goals}
- {id: ttc-y-016, query: "Compare acidic-soil needs with well-drained-soil needs.", intent: editorial-comparison, expected_source_ids: [wp:15288, wp:224522], answerability: expected-answerable, difficulty: hard, leakage_family: compare-soil-concepts}
```

## Climate, zone, and geography (y-017–y-032)
```yaml
- {id: ttc-y-017, query: "Which trees are recommended for Arkansas conditions?", intent: regional-selection, expected_source_ids: [wp:4407], answerability: expected-answerable, difficulty: medium, leakage_family: geography-arkansas}
- {id: ttc-y-018, query: "Which trees are suitable for North Carolina?", intent: regional-selection, expected_source_ids: [wp:4719], answerability: expected-answerable, difficulty: medium, leakage_family: geography-north-carolina}
- {id: ttc-y-019, query: "How should USDA zone numbers guide plant selection?", intent: hardiness-explanation, expected_source_ids: [wp:4237, wp:4116], answerability: expected-answerable, difficulty: medium, leakage_family: climate-zone-interpretation}
- {id: ttc-y-020, query: "What does the letter in a zone such as 7a or 7b mean?", intent: hardiness-explanation, expected_source_ids: [wp:4237], answerability: expected-answerable, difficulty: medium, leakage_family: climate-zone-subzone}
- {id: ttc-y-021, query: "Which product is a cold-climate flowering tree?", intent: constrained-discovery, expected_source_ids: [wp:56433, wp:9410, wp:708570], answerability: expected-answerable, difficulty: hard, leakage_family: climate-cold-flowering}
- {id: ttc-y-022, query: "Which trees or shrubs are drought resistant?", intent: climate-selection, expected_source_ids: [wp:24391, wp:549614, wp:15947], answerability: expected-answerable, difficulty: medium, leakage_family: climate-drought}
- {id: ttc-y-023, query: "Which TTC material addresses trees for coastal or wet places?", intent: climate-selection, expected_source_ids: [wp:635865, wp:3717], answerability: expected-answerable, difficulty: medium, leakage_family: climate-wet}
- {id: ttc-y-024, query: "How does evergreen shade affect plant choices in a garden?", intent: microclimate-explanation, expected_source_ids: [wp:19387, wp:9688], answerability: expected-answerable, difficulty: medium, leakage_family: climate-evergreen-shade}
- {id: ttc-y-025, query: "What trees can provide shade quickly in a hot sunny yard?", intent: climate-selection, expected_source_ids: [wp:4355, wp:5111], answerability: expected-answerable, difficulty: hard, leakage_family: climate-fast-shade}
- {id: ttc-y-026, query: "Which plants are suitable for dry gardens?", intent: climate-selection, expected_source_ids: [wp:595960, wp:24391], answerability: expected-answerable, difficulty: medium, leakage_family: climate-dry-garden}
- {id: ttc-y-027, query: "What should gardeners consider when winter temperatures fluctuate?", intent: climate-care, expected_source_ids: [wp:4237], answerability: expected-answerable, difficulty: hard, leakage_family: climate-winter-temperature}
- {id: ttc-y-028, query: "Which content helps select plants for shade versus full sun?", intent: climate-selection, expected_source_ids: [wp:19387, wp:9688], answerability: expected-answerable, difficulty: medium, leakage_family: climate-light-selection}
- {id: ttc-y-029, query: "What makes a soil site unsuitable because it stays saturated?", intent: climate-soil, expected_source_ids: [wp:224522, wp:635865], answerability: expected-answerable, difficulty: medium, leakage_family: climate-drainage}
- {id: ttc-y-030, query: "Which cypress is catalogued for zones 6 through 9?", intent: constrained-discovery, expected_source_ids: [wp:549614, wp:3709, wp:552438], answerability: expected-answerable, difficulty: hard, leakage_family: climate-cypress-zone}
- {id: ttc-y-031, query: "Which catalogued product has winter color as a feature?", intent: constrained-discovery, expected_source_ids: [wp:549614, wp:779897, wp:762106], answerability: expected-answerable, difficulty: hard, leakage_family: climate-winter-color}
- {id: ttc-y-032, query: "Can TTC documentation tell me the current weather in Raleigh?", intent: weather-current, expected_source_ids: [], answerability: unanswerable-in-declared-corpus, difficulty: easy, leakage_family: abstention-current-weather}
```

## Diagnostic care (y-033–y-048)
```yaml
- {id: ttc-y-033, query: "My leaves are yellow after heavy rain. What cause should I investigate?", intent: plant-diagnosis, expected_source_ids: [wp:4133, wp:224522], answerability: expected-answerable, difficulty: medium, leakage_family: diagnosis-yellow-overwater}
- {id: ttc-y-034, query: "My tree has brown curled leaves. Is too much water the likely issue?", intent: plant-diagnosis, expected_source_ids: [wp:4134, wp:4133], answerability: expected-answerable, difficulty: hard, leakage_family: diagnosis-leaf-contrast}
- {id: ttc-y-035, query: "What should I do when a newly arrived tree looks sick?", intent: arrival-diagnosis, expected_source_ids: [wp:4142, wp:4137, wp:270766], answerability: expected-answerable, difficulty: medium, leakage_family: diagnosis-sick-arrival}
- {id: ttc-y-036, query: "My tree died after planting. Which guarantee or support path applies?", intent: post-planting-support, expected_source_ids: [wp:4140, wp:270766], answerability: expected-answerable, difficulty: medium, leakage_family: diagnosis-tree-dead}
- {id: ttc-y-037, query: "How can I distinguish a pest or fungus concern from a watering issue?", intent: plant-diagnosis, expected_source_ids: [wp:4135, wp:627148, wp:4133], answerability: expected-answerable, difficulty: hard, leakage_family: diagnosis-pest-vs-water}
- {id: ttc-y-038, query: "What is the root-health consequence of poor drainage?", intent: plant-diagnosis, expected_source_ids: [wp:224522], answerability: expected-answerable, difficulty: medium, leakage_family: diagnosis-root-oxygen}
- {id: ttc-y-039, query: "How should a plant be watered outdoors to avoid root problems?", intent: preventative-care, expected_source_ids: [wp:627148, wp:812290], answerability: expected-answerable, difficulty: medium, leakage_family: diagnosis-preventative-water}
- {id: ttc-y-040, query: "My tree is leaning after planting. Should I stake it?", intent: planting-diagnosis, expected_source_ids: [wp:4132, wp:812290], answerability: expected-answerable, difficulty: medium, leakage_family: diagnosis-staking}
- {id: ttc-y-041, query: "Can pruning in fall harm a fruit tree?", intent: seasonal-diagnosis, expected_source_ids: [wp:9892, wp:28084], answerability: expected-answerable, difficulty: medium, leakage_family: diagnosis-fruit-pruning-season}
- {id: ttc-y-042, query: "What should I check before transplanting an established evergreen?", intent: transplant-diagnosis, expected_source_ids: [wp:752143], answerability: expected-answerable, difficulty: hard, leakage_family: diagnosis-evergreen-transplant}
- {id: ttc-y-043, query: "Why might an acid-loving plant show nutrient-related trouble?", intent: soil-diagnosis, expected_source_ids: [wp:15288, wp:418694], answerability: expected-answerable, difficulty: hard, leakage_family: diagnosis-acid-soil}
- {id: ttc-y-044, query: "What can I do if Japanese beetles are damaging a plant?", intent: pest-diagnosis, expected_source_ids: [wp:9612], answerability: expected-answerable, difficulty: easy, leakage_family: diagnosis-japanese-beetle}
- {id: ttc-y-045, query: "Is a wilting tree necessarily dead on arrival?", intent: arrival-diagnosis, expected_source_ids: [wp:4137, wp:270766], answerability: expected-answerable, difficulty: hard, leakage_family: diagnosis-wilting-arrival}
- {id: ttc-y-046, query: "What help is available if my issue is not one of the listed symptoms?", intent: support-escalation, expected_source_ids: [wp:4136, wp:4144], answerability: expected-answerable, difficulty: easy, leakage_family: diagnosis-escalation}
- {id: ttc-y-047, query: "Can TTC diagnose a virus in my tree from a photograph?", intent: diagnosis-service, expected_source_ids: [], answerability: unanswerable-in-declared-corpus, difficulty: medium, leakage_family: abstention-diagnosis-service}
- {id: ttc-y-048, query: "What pesticide concentration should I use for an unidentified fungus?", intent: treatment-specificity, expected_source_ids: [wp:4135], answerability: unknown-requires-source-review, difficulty: hard, leakage_family: ambiguity-unidentified-treatment}
```

## Chunk-sensitive procedures (y-049–y-064)
```yaml
- {id: ttc-y-049, query: "List the ordered steps for putting a ball-and-burlap tree into its planting hole.", intent: ordered-procedure, expected_source_ids: [wp:398454], answerability: expected-answerable, difficulty: hard, leakage_family: procedure-ball-burlap-order}
- {id: ttc-y-050, query: "What should be removed from a ball-and-burlap root ball before backfilling?", intent: procedure-detail, expected_source_ids: [wp:398454], answerability: expected-answerable, difficulty: medium, leakage_family: procedure-ball-burlap-material}
- {id: ttc-y-051, query: "How do root-ball width and hole depth relate in general planting?", intent: procedure-geometry, expected_source_ids: [wp:812290], answerability: expected-answerable, difficulty: medium, leakage_family: procedure-hole-geometry}
- {id: ttc-y-052, query: "What is the first-month watering sequence after planting?", intent: ordered-procedure, expected_source_ids: [wp:812290], answerability: expected-answerable, difficulty: medium, leakage_family: procedure-first-month-water}
- {id: ttc-y-053, query: "What are the steps for planting a privacy screen in a straight line?", intent: ordered-procedure, expected_source_ids: [wp:405509], answerability: expected-answerable, difficulty: hard, leakage_family: procedure-privacy-screen}
- {id: ttc-y-054, query: "What does the evergreen-hedge guide say about preparing and spacing a hedge?", intent: ordered-procedure, expected_source_ids: [wp:405437], answerability: expected-answerable, difficulty: hard, leakage_family: procedure-evergreen-hedge}
- {id: ttc-y-055, query: "What is the procedure for planting a bare-root tree?", intent: ordered-procedure, expected_source_ids: [wp:405431], answerability: expected-answerable, difficulty: hard, leakage_family: procedure-bare-root}
- {id: ttc-y-056, query: "What follow-up care is recommended after planting a citrus tree?", intent: ordered-procedure, expected_source_ids: [wp:418611], answerability: expected-answerable, difficulty: hard, leakage_family: procedure-citrus-followup}
- {id: ttc-y-057, query: "What steps are recommended to plant apple and pear trees?", intent: ordered-procedure, expected_source_ids: [wp:418609, wp:9210], answerability: expected-answerable, difficulty: hard, leakage_family: procedure-apple-pear}
- {id: ttc-y-058, query: "How should peach and nectarine trees be planted?", intent: ordered-procedure, expected_source_ids: [wp:398553, wp:9208], answerability: expected-answerable, difficulty: hard, leakage_family: procedure-peach-nectarine}
- {id: ttc-y-059, query: "What should I do before and after planting a Japanese maple?", intent: ordered-procedure, expected_source_ids: [wp:418603, wp:9194], answerability: expected-answerable, difficulty: hard, leakage_family: procedure-japanese-maple}
- {id: ttc-y-060, query: "What sequence is recommended for planting roses?", intent: ordered-procedure, expected_source_ids: [wp:418613], answerability: expected-answerable, difficulty: hard, leakage_family: procedure-rose}
- {id: ttc-y-061, query: "How should hydrangeas be planted and watered initially?", intent: ordered-procedure, expected_source_ids: [wp:418605], answerability: expected-answerable, difficulty: hard, leakage_family: procedure-hydrangea}
- {id: ttc-y-062, query: "What should be done when staking a newly planted tree?", intent: ordered-procedure, expected_source_ids: [wp:812290, wp:4132], answerability: expected-answerable, difficulty: medium, leakage_family: procedure-staking}
- {id: ttc-y-063, query: "What are the safe steps for removing a tree stump?", intent: ordered-procedure, expected_source_ids: [wp:9832], answerability: expected-answerable, difficulty: hard, leakage_family: procedure-stump-removal}
- {id: ttc-y-064, query: "Can I follow the corpus for instructions to kill a healthy neighbor's tree?", intent: harmful-procedure, expected_source_ids: [wp:9551], answerability: ambiguous-requires-policy-review, difficulty: hard, leakage_family: ambiguity-harmful-tree-removal}
```

## Transaction workflows and abstention/ambiguity (y-065–y-080)
```yaml
- {id: ttc-y-065, query: "What is the workflow after an order is delivered but an item is missing?", intent: transaction-workflow, expected_source_ids: [wp:456943, wp:4138], answerability: expected-answerable, difficulty: hard, leakage_family: transaction-wrong-item}
- {id: ttc-y-066, query: "How do I start a return and what happens after it is received?", intent: transaction-workflow, expected_source_ids: [wp:4230, wp:558351, wp:398600], answerability: expected-answerable, difficulty: hard, leakage_family: transaction-return}
- {id: ttc-y-067, query: "What information is needed to ask about an existing order?", intent: transaction-workflow, expected_source_ids: [wp:398563, wp:558358], answerability: expected-answerable, difficulty: medium, leakage_family: transaction-existing-order}
- {id: ttc-y-068, query: "How does a customer use a support ticket for a delivery issue?", intent: transaction-workflow, expected_source_ids: [wp:52190, wp:456943], answerability: expected-answerable, difficulty: medium, leakage_family: transaction-support-ticket}
- {id: ttc-y-069, query: "Can I choose a delivery day rather than a shipping date?", intent: transaction-workflow, expected_source_ids: [wp:76497, wp:4125], answerability: expected-answerable, difficulty: medium, leakage_family: transaction-shipping-date}
- {id: ttc-y-070, query: "What does an order's pre-sale state mean for when it will ship?", intent: transaction-workflow, expected_source_ids: [wp:76498, wp:4125], answerability: expected-answerable, difficulty: medium, leakage_family: transaction-presale-ship}
- {id: ttc-y-071, query: "What account page or process helps track an order?", intent: transaction-workflow, expected_source_ids: [wp:558358, wp:4124], answerability: expected-answerable, difficulty: medium, leakage_family: transaction-account-status}
- {id: ttc-y-072, query: "Can I receive a printed TTC catalog by mail?", intent: catalog-policy, expected_source_ids: [wp:4143], answerability: expected-answerable, difficulty: easy, leakage_family: transaction-catalog}
- {id: ttc-y-073, query: "What warranty route applies after the plant arrives?", intent: warranty-workflow, expected_source_ids: [wp:398593, wp:270766], answerability: expected-answerable, difficulty: hard, leakage_family: transaction-warranty-vs-guarantee}
- {id: ttc-y-074, query: "Can I cancel an order and receive a specific fee quote?", intent: policy-conflict-detection, expected_source_ids: [wp:398597, wp:4128, wp:4129], answerability: conflicting-requires-policy-owner, difficulty: hard, leakage_family: conflict-cancellation}
- {id: ttc-y-075, query: "Does TTC guarantee delivery by next Tuesday?", intent: delivery-guarantee, expected_source_ids: [wp:76497, wp:4125], answerability: unanswerable-in-declared-corpus, difficulty: medium, leakage_family: abstention-delivery-guarantee}
- {id: ttc-y-076, query: "Can I pick a precise carrier and tracking number before purchase?", intent: shipping-detail, expected_source_ids: [wp:76492, wp:4126], answerability: unknown-requires-source-review, difficulty: hard, leakage_family: ambiguity-carrier-tracking}
- {id: ttc-y-077, query: "Does TTC ship plants to Canada?", intent: shipping-eligibility, expected_source_ids: [wp:76495], answerability: unanswerable-in-declared-corpus, difficulty: medium, leakage_family: abstention-international-shipping}
- {id: ttc-y-078, query: "Can I get a cash refund for a plant that I merely dislike?", intent: returns-policy, expected_source_ids: [wp:558351, wp:398600], answerability: expected-answerable, difficulty: medium, leakage_family: transaction-return}
- {id: ttc-y-079, query: "Will a product that is out of stock be reservable for me?", intent: inventory-policy, expected_source_ids: [wp:398551, wp:76498], answerability: expected-answerable, difficulty: medium, leakage_family: transaction-out-of-stock}
- {id: ttc-y-080, query: "What is the current sale price of Thuja Green Giant today?", intent: current-price, expected_source_ids: [wp:3699], answerability: unanswerable-in-frozen-snapshot, difficulty: medium, leakage_family: abstention-current-price}
```

## Mandatory qualification gate

For every candidate, verify the actual text supports its premise; resolve
revision IDs and exact evidence spans; identify partial and adversarial sources;
then obtain human adjudication. The `unanswerable`, `unknown`, `ambiguous`, and
`conflicting` values are test-behavior proposals, not factual claims about
TTC's live business.
