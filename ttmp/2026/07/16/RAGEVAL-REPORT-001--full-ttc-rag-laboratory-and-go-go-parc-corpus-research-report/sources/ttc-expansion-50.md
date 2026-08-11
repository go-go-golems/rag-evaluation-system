---
Title: TTC Expansion Audit and 50-Card Source-Grounded Draft
Ticket: RAGEVAL-TTC-LAB-001
Type: Reference
Status: draft
Created: 2026-07-16
---

# TTC expansion audit and 50-card source-grounded draft

## Status and method

This is a source-discovery queue, not an evaluation manifest. Candidate source
IDs were read from `data/ttc-wordpress-rag.sqlite` on 2026-07-16. They are not
relevance labels, frozen revision IDs, chunk IDs, or approved answers. Each
card still needs source/revision inspection, chunk evidence, review, an
evidence-family ID, and a partition before it can enter the immutable dataset.

The v2 documentation refers to a 72-card proposed pool, but no separately
recoverable individual-card list exists in the current ticket tree. The visible
materials are the 20-card pilot and a 12-card direct-FAQ adjudication batch.
Thus “72” is currently a planning count, not auditable card evidence. Adding
the 50 proposals below would reach a 122-card *candidate* queue, still 78
short of the stated 200-card candidate target and much further short of 200
adjudicated cards.

## Missing intent strata

| Stratum | Gap exposed by current visible queue |
|---|---|
| Product attribute discrimination | Insufficient numeric, botanical-name, fielded metadata, and near-neighbor product cases. |
| Constraint comparison | Insufficient multi-facet and multi-document answer requirements. |
| State/climate suitability | Insufficient USDA zone, state page, and geographic query phrasing. |
| Planting method | Insufficient separation of bare-root, evergreen, hedge, bamboo, and ball-and-burlap instructions. |
| Diagnostic care | Insufficient separation of transit, plant-care, disease/pest, and warranty evidence. |
| Policy freshness/conflict | Needs explicit known-overlap controls rather than accidental scoring contamination. |
| Transaction workflow | Needs order lifecycle/account/return cases independent from horticulture. |
| Negative/ambiguous control | Needs corpus-wide verified abstention and deliberately ambiguous near-miss cases. |
| Chunk-sensitive evidence | Needs questions whose correct answer lies in a specific guide section rather than anywhere in a correct document. |

## Proposed cards — no final labels

| ID | Intent / candidate query | Candidate stable `wp:*` source IDs |
|---|---|---|
| x-001 | Accepted payment methods: “Which payment methods can I use for an order?” | `wp:4120` |
| x-002 | Order status: “How can a customer check the current status of an order?” | `wp:4124`, `wp:16771`, `wp:398559` |
| x-003 | Shipment timing: “When should an order normally ship after it is placed?” | `wp:4125` |
| x-004 | Shipping method: “How are live plants shipped to customers?” | `wp:76492`, `wp:4235` |
| x-005 | Shipping destinations: “Which locations can this store ship plants to?” | `wp:76495` |
| x-006 | Company location: “Where is The Tree Center located?” | `wp:76336` |
| x-007 | Cancellation fee: “Is there a charge for cancelling an order?” | `wp:4129`, `wp:398597` |
| x-008 | Return workflow: “How does a customer start a return?” | `wp:4230`, `wp:558351` |
| x-009 | Warranty/guarantee: “What guarantee applies if an ordered tree does not survive?” | `wp:4140`, `wp:4148`, `wp:270766`, `wp:398593` |
| x-010 | Phone ordering: “Can I place an order by phone instead of online?” | `wp:398551` |
| x-011 | Wilting on arrival: “What should I do if my tree arrives wilted?” | `wp:4137` |
| x-012 | Wrong size received: “What should I do if the delivered tree is not the size I ordered?” | `wp:4139` |
| x-013 | Yellowing leaves: “Why might leaves be turning yellow after planting?” | `wp:4133` |
| x-014 | Curling brown leaves: “What causes leaves to curl and turn brown?” | `wp:4134` |
| x-015 | Pest/fungus: “Where should I look for guidance when I suspect a pest or fungal problem?” | `wp:4135` |
| x-016 | Sick tree: “What support is available when a tree looks sick?” | `wp:4142` |
| x-017 | Not-shipped order: “Which support route applies to an order that has not shipped?” | `wp:456943`, `wp:456973` |
| x-018 | Delivery/return policy: “Where is the delivery and returns policy explained?” | `wp:398600`, `wp:558329` |
| x-019 | Bare-root planting: “How should I plant a bare-root tree?” | `wp:405431`, `wp:9171` |
| x-020 | Ball-and-burlap planting: “What planting method is recommended for ball-and-burlap trees?” | `wp:398454`, `wp:9173` |
| x-021 | Evergreen planting: “How do I plant an evergreen tree?” | `wp:405433`, `wp:9175` |
| x-022 | Garden shrub planting: “What are the planting instructions for a garden shrub?” | `wp:405435`, `wp:9182` |
| x-023 | Evergreen hedge: “How should an evergreen hedge be planted?” | `wp:405437`, `wp:9187` |
| x-024 | Deciduous hedge: “What differs when planting a deciduous hedge?” | `wp:405441`, `wp:9189` |
| x-025 | Bamboo planting: “How should bamboo trees be planted?” | `wp:398536`, `wp:9204` |
| x-026 | Hydrangea planting: “What instructions apply when planting hydrangeas?” | `wp:418605`, `wp:9202` |
| x-027 | Japanese maple planting: “How should a Japanese maple be planted?” | `wp:418603`, `wp:9194` |
| x-028 | Privacy screen: “How do I plant a privacy screen?” | `wp:405509`, `wp:9195` |
| x-029 | Green Giant zones: “What USDA hardiness zones suit Thuja Green Giant?” | `wp:3699` |
| x-030 | Green Giant dimensions: “How tall and wide can Thuja Green Giant become?” | `wp:3699` |
| x-031 | Drought juniper: “Which narrow blue juniper is rated very drought resistant?” | `wp:3704` |
| x-032 | Sky Pencil width: “How wide does Sky Pencil Holly grow at maturity?” | `wp:3708` |
| x-033 | Full-sun Arizona Cypress: “Which Carolina Sapphire listing specifies full sun?” | `wp:3709` |
| x-034 | Meyer Lemon zones: “Which hardiness zones are listed for Meyer Lemon?” | `wp:3767` |
| x-035 | Meyer Lemon height: “How tall can a Meyer Lemon Tree grow?” | `wp:3767` |
| x-036 | Honeycrisp zones: “Which zones are listed for Honeycrisp apple trees?” | `wp:3787` |
| x-037 | Arbequina drought tolerance: “Is Arbequina Olive described as drought resistant?” | `wp:3814` |
| x-038 | Chicago Hardy Fig zones: “Which zone range is listed for Chicago Hardy Fig?” | `wp:3826` |
| x-039 | Arabica light conditions: “What light conditions are listed for Arabica Coffee?” | `wp:3838` |
| x-040 | Rainbow Eucalyptus height: “What mature height is listed for Rainbow Eucalyptus?” | `wp:3742` |
| x-041 | Narrow privacy screen: “Which listed evergreen is narrow and tolerates partial shade?” | `wp:3699`, `wp:3701`, `wp:3707`, `wp:3708` |
| x-042 | Cold-zone fruit: “Which is suitable for colder zones, Honeycrisp or Meyer Lemon?” | `wp:3787`, `wp:3767` |
| x-043 | Full-sun drought constraint: “Which listed tree combines full sun with very high drought tolerance?” | `wp:3704`, `wp:3709`, `wp:3814` |
| x-044 | Compact Japanese maple: “Which listed Japanese maple has the smallest mature height range?” | `wp:3734`, `wp:3743`, `wp:3745` |
| x-045 | Fig/lime winter tolerance: “Which is listed for colder zones, Chicago Hardy Fig or Persian Lime?” | `wp:3826`, `wp:3777` |
| x-046 | California trees: “Where can a customer find trees offered for California?” | `wp:4436`, `wp:4233` |
| x-047 | Zone determination: “How can I determine a plant hardiness zone before choosing a tree?” | `wp:4116`, `wp:4237` |
| x-048 | Wholesale: “Where is information for wholesale purchasing published?” | `wp:14867`, `wp:4123` |
| x-049 | Cancellation conflict: “What cancellation policy applies after an order is placed?” | `wp:4128`, `wp:4129`, `wp:398597` |
| x-050 | Unanswerable control: “Does the catalog guarantee same-day drone delivery to Alaska?” | _none; requires corpus-wide absence review_ |

## Adjudication and partition guardrails

1. Deduplicate by information need, not wording, against the recovered 72 and
   pilot; shipment/cancellation/planting have close neighbors.
2. Treat legacy guide/page pairs (for example `wp:405431`/`wp:9171`) as one
   evidence family, not independent train/holdout cases.
3. Keep `x-049` out of tuning until policy precedence is explicitly reviewed.
4. Do not assign authoritative sources merely because they appear in this
   table; validate every document revision and exact chunk span.
5. Promote no `x-*` ID directly. Accepted cards need immutable IDs, revision
   IDs, chunk judgments, reviewer metadata, family grouping, and partition.
