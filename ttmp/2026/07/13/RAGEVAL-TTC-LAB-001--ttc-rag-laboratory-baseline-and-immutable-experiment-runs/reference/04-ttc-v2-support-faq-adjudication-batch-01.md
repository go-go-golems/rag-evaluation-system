---
Title: TTC v2 Support FAQ Adjudication Batch 01
Ticket: RAGEVAL-TTC-LAB-001
Type: Reference
Status: review
Created: 2026-07-16
---

# Purpose

This is the first concrete adjudication packet for the TTC v2 evaluation corpus. It covers twelve direct support/FAQ cards whose authoritative source is unambiguous in the immutable TTC snapshot. The packet is intentionally separate from the dataset manifest: a reviewer can inspect the source revision and approve or amend the labels before a frozen dataset is registered.

## Accepted evidence labels

- `authoritative`: the source directly answers the question and is the preferred citation.
- `substantial`: the source contains useful supporting information but is not sufficient as the sole answer.
- `misleading`: lexical overlap is likely, but the source would lead an answer in the wrong direction.
- `abstention_expected`: the corpus does not contain enough evidence for a responsible answer.

## Direct-source batch

| Card | Intent | Authoritative source | Revision ID | Proposed label | Adjudication note |
|---|---|---|---|---|---|
| v2-001 | bulk order discount | `wp:4123` | `sha256:bda0a3b59fb2a55965edc1b92df6c009dcdeb2253d5a65e87a1b3533916d1d68` | authoritative | Direct bulk-discount FAQ. |
| v2-002 | gift ordering | `wp:4122` | `sha256:e26aee041ba6aebb8bf443c5d81eac1f47446bbd525a1d02ebd824b72d618d3a` | authoritative | Direct gift FAQ. |
| v2-004 | shipment notification | `wp:4126` | `sha256:210acf374156b634e02ccdb85ecba1081f6f99b2aa05c09dcb3696ec140f4272` | authoritative | Direct shipped-email/tracking FAQ. |
| v2-005 | shipping charges | `wp:4121` | `sha256:2fab60d9032faedee44c2b77431c1ed7fe7df3bb2f582b6ca1dc70966e5b2d2f` | authoritative | Direct shipping-charge FAQ. |
| v2-007 | edit an order | `wp:69806` | `sha256:cfdc238136ce566c138e409d4a24ee76f6519024703a2507c83773ec17d44021` | authoritative | Direct edit-order FAQ. |
| v2-008 | delayed shipment | `wp:4127` | `sha256:6b3b6dc66eabfd529ef3d30680f7455a5db446eeb4ca01143f2e384686d27331` | authoritative | Direct delay-shipment FAQ. |
| v2-010 | presale timing | `wp:76498` | `sha256:460b77df17e9205855ac78c6eb4d6e29914af3cce9015ff43862080ecd76a21c` | authoritative | Direct presale FAQ. |
| v2-011 | wrong plant received | `wp:4138` | `sha256:b9e496dcadf8b9b928dc9dc57a43c083bef706ca785f9bfc5826bdf454c12abd` | authoritative | Direct wrong-item FAQ. |
| v2-014 | gallon size versus height | `wp:76490` | `sha256:a774c1a5834869ce25889593486f9c7c2dea3e732308a57ff4c3e3abd5ab6253` | authoritative | Direct container-size FAQ. |
| v2-016 | printed catalog availability | `wp:4143` | `sha256:21d083061f75b8b98aa0eeadd26d44abb99926d2ac0a6b000727fff6244d6c9e` | authoritative | Direct catalog FAQ. |
| v2-018 | request more help | `wp:4144` | `sha256:52503485bb1a1189686da4b7db1a821729f477ffae3333d6012dc46f3f837888` | authoritative | Direct customer-service FAQ. |
| v2-019 | USDA extension office | `wp:4146` | `sha256:d47cbbab544af94f36e93270b0d8fec2b7b7892db862ea4d8087f111e98331a4` | authoritative | Direct USDA-extension FAQ. |

## Review procedure

1. Resolve each revision ID in the immutable snapshot and read the complete source, not only the retrieved chunk.
2. Confirm that the proposed source answers the exact question and record any required scope or date qualification.
3. Add at least one negative or misleading candidate for each card when lexical neighbors exist.
4. Record reviewer, timestamp, and decision in the eventual v2 manifest. Do not edit v1 in place.

Cards v2-003, v2-006, v2-009, v2-012, v2-013, v2-015, and v2-017 remain in the pooled-adjudication queue because their answers require comparing multiple source revisions. Card v2-020 is reserved for an explicit expected-abstention case.
