---
Title: Workflow V3 Umans Batching and Concurrency Study
Ticket: RAG-TTC-V3-SWEEP
Status: complete
Topics:
    - rag-eval
    - evaluation
    - workflow
    - chunking
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Completed bounded real-provider batching/concurrency qualification with durable operation custody; the corrected analysis is maintained in the pipeline consolidation audit ticket.
LastUpdated: 2026-07-22T18:50:00-04:00
WhatFor: Find the original execution artifacts, design, diary, and corrected performance analysis for the Workflow V3 Umans sweep.
WhenToUse: When reviewing the July 22 real qualification or planning a separately authorized replicated follow-up.
---

# Workflow V3 Umans Batching and Concurrency Study

## Overview

This ticket designed and executed the bounded Workflow V3 Umans batching/concurrency qualification. The final real run processed 16 fixed chunks across batch sizes 1/2/4/8 and generation concurrency 1/2, retained per-cell operation custody, passed privacy checks, and was imported into researchctl.

A later audit found that the initial result note was too shallow and that successful-attempt makespan omitted 21.472 seconds in one retry cell. Use the corrected operation-ledger analysis below for performance conclusions.

## Key links

- [Design and implementation guide](design-doc/01-workflow-v3-umans-batching-and-concurrency-study-design-and-implementation-guide.md)
- [Investigation diary](reference/01-investigation-diary.md)
- [Initial real result note](analysis/02-authorized-real-umans-qualification-results.md)
- [Corrected full performance audit and graphs](../../RAG-PIPELINE-CONSOLIDATION-AUDIT--assess-improve-and-consolidate-the-rag-evaluation-pipeline/analysis/01-ttc-real-run-performance-audit-and-pipeline-consolidation-assessment.md)
- [Real evidence](sources/real-attempt-003/evidence.json)
- [Real operation ledgers](sources/real-attempt-003/operations/)
- [Original real graphs](sources/real-attempt-003/graphs/manifest.json)

## Status

Current status: **complete** for the bounded single-replicate qualification. No second real run is authorized. The corrected analysis recommends batch-8/concurrency-2 for this exact shape, subject to replicated quality-aware validation.

## Tasks

See [tasks.md](./tasks.md).

## Changelog

See [changelog.md](./changelog.md).
