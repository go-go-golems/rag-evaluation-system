# TTC chunk-level evidence-label proposal

This is a **review draft**, derived read-only from `data/rag-eval.db` on
2026-07-16. It does not alter the candidate dataset, judgments, corpus, or
immutable runs. The current dataset contains 20 scored cards plus one
withheld conflict card—not 25–35 existing cards—so this draft covers all 21
available cards. Expansion cards need to be authored and source-validated
before they can receive defensible chunk labels.

## Scope and notation

All proposed chunks are in immutable chunk set
`sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392`
(2024 chunks; nominal 1200-rune windows with 150-rune overlap). `doc / idx /
[start,end)` identifies a chunk; the document revision is resolved from its
stable `wp:*` source ID in the database. Product documents repeat their
structured detail blocks many times. For those rows, the first listed index is
a canonical *proposed* evidence occurrence, not a claim that identical later
occurrences are non-relevant.

`A` means authoritative evidence; `S` means substantial supporting evidence;
`M` is a deliberately misleading or partial document/chunk candidate. Exact
immutable chunk IDs should be materialized alongside these labels only when a
new reviewed dataset is frozen, since the current source dataset is still a
candidate draft.

| Card | Proposed evidence label(s) | Possible misleading / partial evidence | Confidence and caveat |
|---|---|---|---|
| ttc-eval-001 | **A:** `wp:3699 / 46 / [48300,49500)` — product details have 20–40, 6–12, zone 5–9, full sun to partial shade. | **M:** `wp:3701` Leyland Cypress product-detail chunks; privacy-screen language can swamp entity match. | High for facts; repeated structured blocks need deduplication at judgment compile time. |
| ttc-eval-002 | **A:** `wp:549614 / 11 / [11550,12750)` — 15–25, 6–8, zone 6–9, full sun, very drought resistant; companion privacy-screen wording appears in later repeated chunks. | **M:** `wp:3709` Carolina Sapphire (similar climate/drought); `wp:552438` Silver Smoke (dimensions but no privacy category). | Medium-high: one chunk carries facts but category membership is metadata elsewhere in document. |
| ttc-eval-003 | **S:** `wp:15947 / 8 / [8400,9600)` Blue Italian facts; **S:** `wp:3703 / 53 / [55650,56850)` Italian facts. Together form the answer. | **M:** `wp:3701` Leyland Cypress. | High; this is intentionally a two-document retrieval target, never a single-chunk authority. |
| ttc-eval-004 | **A:** `wp:7347 / 35 / [36750,37950)` — Thuja occidentalis Danica, 1–2 height and width. | **M:** `wp:26028` Thuja Can Can detail chunk (taxonomy match, wrong dimensions). | High. |
| ttc-eval-005 | **A:** `wp:3717 / 26 / [27300,28500)` — Bald Cypress 50–70 height; **S:** `wp:3717 / 6 / [6300,7500)` — Tolerates Wet Soil metadata. | **M:** `wp:10069` Red Star White Cypress, wet soil but wrong height. | Medium-high: conjunction spans two chunks, so label both rather than pretending one is complete. |
| ttc-eval-006 | **A:** `wp:812290 / 0 / [0,1200)` — first-month schedule and planting-hole geometry/root-collar guidance. | **M:** `wp:4131 / 0` arrival FAQ directs to guide without measurements. | High; this compact guide is an unusually clean single-chunk label. |
| ttc-eval-007 | **A:** `wp:812290 / 0 / [0,1200)` — soak twice weekly for first month. **S:** `wp:627148` chunk containing root-zone soaking principles. | **M:** any general watering chunk that lacks the initial-month cadence. | High for authority; exact substantial chunk should be selected after a reviewer reads the long article. |
| ttc-eval-008 | **A:** `wp:9892 / 3 / [3150,4350)` — explicitly avoids fall; adjacent/preceding pruning section supplies post-last-frost timing. **S:** `wp:28084` late-winter pruning chunk. | **M:** `wp:751617` young-tree pruning basics. | Medium: answer facets straddle a chunk boundary; preserve both 2 and 3 during adjudication. |
| ttc-eval-009 | **A:** `wp:4133 / 0 / [0,390)` — yellow leaves normally signal excess water/poor drainage. **S:** `wp:627148` wet-soil/root-oxygen explanation. | **M:** `wp:4134 / 0` curling/brown leaves describes a materially different diagnosis. | High. |
| ttc-eval-010 | **A:** `wp:398454 / 0 / [0,1200)` B&B identity and procedure; **A:** `/ 1 / [1050,2250)` removal/planting steps; **S:** `/ 5–6 / [5250,7500)` hydrated root ball and depth. | **M:** `wp:405431` bare-root planting chunks. | High for topic, medium for exact set: procedure requires a multi-chunk evidence group. |
| ttc-eval-011 | **A:** `wp:405509 / 0 / [0,1200)` and `/ 4 / [4200,5400)` — screen form/spacing and trim control. **S:** `wp:405437` hedge chunks. | **M:** `wp:8017` screening-selection article. | Medium-high; inspect chunk 1–4 as a contiguous group before freeze. |
| ttc-eval-012 | **A:** `wp:15288 / 0–2 / [0,3300)` — acidic soil, pH and nutrient/iron consequence. **S:** `wp:418694` practical acid-planting chunk. | **M:** `wp:224522` well-drained soil explanation. | High conceptually; multi-chunk only because article introduces and explains in sequence. |
| ttc-eval-013 | **A:** `wp:19387 / 2–4 / [2100,5400)` — deciduous/building contrast then dense, permanent evergreen shade. **S:** `wp:9688` shade categories. | **M:** `wp:4355` shade-tree purchasing/selection content. | High; group is needed for the comparison’s two sides. |
| ttc-eval-014 | **A:** `wp:224522 / 1–2 / [1050,3300)` — pores, air, root oxygen and saturated-soil gas-exchange failure. | **M:** `wp:812290 / 0`, practical damp-not-soaked guidance without mechanism. | High. |
| ttc-eval-015 | **A:** `wp:4237 / 1 / [1050,2250)` temperature basis; **A:** `/ 4 / [4200,5400)` number plus a/b subzone explanation. **S:** `wp:4116` hardiness-zone FAQ. | **M:** product chunks that merely list a hardiness zone. | High; retain two authoritative chunks for distinct facets. |
| ttc-eval-016 | **A:** `wp:76495 / 0 / [0,244)` — California exclusion and Florida Citrus restriction. | **M:** `wp:76497 / 0`, shipping-date semantics. | High. |
| ttc-eval-017 | **A:** `wp:76497 / 0 / [0,700)` — preferred time means estimated ship date, not arrival. | **M:** `wp:456943` delivery/shipped-order support chunks. | High. |
| ttc-eval-018 | **A:** `wp:558351 / 0–1 / [0,2250)` — seven-day return, buyer shipping, store credit/no cash-refund policy. **S:** `wp:398600` abbreviated delivery/returns text. | **M:** `wp:398593` warranty article. | High, subject to policy-owner review already required by candidate card. |
| ttc-eval-019 | **A:** `wp:270766 / 0–1 / [0,1315)` — arrival guarantee, 30 days, picture/description and remedy. **S:** `wp:4140` routing FAQ. | **M:** generic warranty or return-policy chunks. | High. |
| ttc-eval-020 | **No positive chunk.** **Partial/context only:** `wp:398551 / 0 / [0,1200)` enumerates accepted methods but cannot prove Bitcoin is refused. | Retrieval of payment-method list can induce overconfident negative answers. | High abstention label: grade refusal to invent a rejection, not a document result. |
| ttc-eval-withheld-001 | **Conflict evidence:** `wp:398597` cancellation-policy chunk(s), 20%/pre-fulfilment; `wp:4128` cancellation FAQ chunk(s), 10%/one-hour. | Each is misleading if presented as a single current answer without policy precedence. | High conflict detection, deliberately **not scoreable** until owner adjudication. |

## Proposed review rules before implementation

1. Label contiguous evidence groups where an answer necessarily spans chunks;
   do not force a false single-chunk ground truth.
2. Promote structured product-detail chunks only when they contain every
   requested facet; otherwise retain a paired metadata/text chunk group.
3. Mark the repeated copies in oversized product pages as duplicate evidence,
   not independently relevant documents. This makes chunk metrics sensitive to
   evidence selection rather than extraction repetition.
4. Preserve `ttc-eval-020` as a no-positive control and the withheld card as a
   conflict-control; neither should be silently converted into ordinary
   positive retrieval labels.

