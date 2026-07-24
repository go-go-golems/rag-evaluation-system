---
Title: Hard-cut operator worker pools in favor of Workflow V3 batch concurrency
Ticket: RAG-WORKFLOW-CONCURRENCY-HARD-CUT
Status: active
Topics:
    - rag
    - workflow
    - evaluation
    - embeddings
    - go
    - security
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-24T14:58:29.943891316-04:00
WhatFor: ""
WhenToUse: ""
---

# Hard-cut operator worker pools in favor of Workflow V3 batch concurrency

## Overview

This ticket performs the minimal concurrency cleanup before broader TTC scaling work resumes. Workflow V3 becomes the normal scheduler for generation and embedding batches; direct engine preparation becomes serial; `GenerationConcurrency`, operator worker pools, hidden worker defaults, duplicate progress machinery, and superseded coarse provider-task paths are hard-cut without compatibility modes.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- rag
- workflow
- evaluation
- embeddings
- go
- security

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
