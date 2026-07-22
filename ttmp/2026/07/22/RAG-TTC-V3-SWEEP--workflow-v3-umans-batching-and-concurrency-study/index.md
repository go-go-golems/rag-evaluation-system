---
Title: Workflow V3 Umans Batching and Concurrency Study
Ticket: RAG-TTC-V3-SWEEP
Status: active
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
Summary: Design, implementation, fixture evidence, and eventual bounded real-provider results for the Workflow V3 Umans chunks-per-request and concurrency sweep.
LastUpdated: 2026-07-22T09:58:57.045540453-04:00
WhatFor: Choose TTC generation batching and concurrency from precise, bounded, reproducible evidence.
WhenToUse: When implementing, operating, reviewing, or publishing the Umans performance sweep.
---

# Workflow V3 Umans Batching and Concurrency Study

## Overview

Measure the interaction between chunks per Umans generation request (`1, 2, 4, 8`) and Workflow V3 concurrency (`1, 2, 4`). The no-cost fixture control is implemented and executed. Real-provider execution remains gated on host-local credentials and an explicit numeric request/token/cost ceiling.

## Key Links

- [Design and implementation guide](design-doc/01-workflow-v3-umans-batching-and-concurrency-study-design-and-implementation-guide.md)
- [Investigation diary](reference/01-investigation-diary.md)
- [Fixture evidence](sources/fixture-control/evidence.json)
- [Fixture graphs](sources/fixture-control/graphs/manifest.json)

## Status

Current status: **active**

## Topics

- rag-eval
- evaluation
- workflow
- chunking

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design-doc/ - Architecture and implementation guide
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
