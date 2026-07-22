---
Title: Assess, improve, and consolidate the RAG evaluation pipeline
Ticket: RAG-PIPELINE-CONSOLIDATION-AUDIT
Status: review
Topics:
    - rag-eval
    - evaluation
    - workflow
    - ttc
    - research
    - visualization
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Audit and repair of the July 22 TTC result package, plus a design for one maintained execution-to-publication pipeline.
LastUpdated: 2026-07-22T18:49:00-04:00
WhatFor: Find the corrected real-run analysis, reproducible graph script, evidence, diary, and pipeline consolidation plan.
WhenToUse: When reviewing the TTC result, planning RAG measurement changes, or implementing report/publication automation.
---

# Assess, improve, and consolidate the RAG evaluation pipeline

## Overview

This ticket reconstructs the July 22 real TTC performance result from durable Workflow V3 operation evidence after the original published report proved materially incomplete. It includes corrected elapsed boundaries, retry and latency analysis, five visually reviewed graph families, a complete Markdown report, and an architecture plan for consolidating execution, reduction, researchctl custody, reporting, and publication.

## Key links

- [Full performance audit and graph showcase](analysis/01-ttc-real-run-performance-audit-and-pipeline-consolidation-assessment.md)
- [Pipeline consolidation design](design-doc/01-rag-pipeline-assessment-consolidation-and-improvement-plan.md)
- [Investigation diary](reference/01-investigation-diary.md)
- [Reproducible analysis script](scripts/01-analyze-ttc-real-run.py)
- [Derived summary](sources/derived-real-attempt-003/summary.json)
- [Graph manifest](sources/derived-real-attempt-003/graphs/manifest.json)

## Main finding

The original aggregate makespan used successful attempts only. Batch-2/concurrency-1 omitted 21.472 seconds of failed provider work. The corrected operation-inclusive analysis shows that batch-8/concurrency-2 was fastest in this single replicate, provider spans covered 98.99%–99.94% of cell elapsed time, and broad scheduler pauses were not the cause of slowness.

## Status

Current status: **review**. Analysis, graphs, design, and md-view inspection are complete. The next work is implementation of the proposed maintained analysis/finalization pipeline and a separately authorized replicated study.

## Tasks

See [tasks.md](./tasks.md).

## Changelog

See [changelog.md](./changelog.md).

## Structure

- `analysis/` — full result and architecture assessment
- `design-doc/` — implementation design
- `reference/` — chronological investigation diary
- `scripts/` — reproducible ticket-local analysis scripts
- `sources/` — derived bounded datasets and figures
