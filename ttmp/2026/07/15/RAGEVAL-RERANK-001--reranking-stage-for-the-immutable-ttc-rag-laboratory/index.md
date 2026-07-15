---
Title: Reranking stage for the immutable TTC RAG laboratory
Ticket: RAGEVAL-RERANK-001
Status: active
Topics:
    - rag
    - reranking
    - ttc
    - geppetto
    - ollama
    - rag-eval
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Design and implementation ticket for inserting a reproducible cross-encoder reranking stage after lexical/vector candidate retrieval and before final citation hydration.
LastUpdated: 2026-07-15T16:55:29.28654433-04:00
WhatFor: Establish the durable API, runtime contract, experiment provenance, and rollout plan for cross-encoder reranking in the TTC RAG laboratory.
WhenToUse: Read before implementing a reranker client, changing trace schemas, starting llama.cpp on the Mac, or comparing reranked and baseline runs.
---

# Reranking stage for the immutable TTC RAG laboratory

## Overview

This ticket introduces a bounded reranking stage for the existing immutable
TTC RAG laboratory. Candidate retrieval remains lexical/vector/RRF and keeps
its current artifact lineage. A configured cross-encoder receives only the
top candidate texts, returns a score per candidate, and supplies a new final
order before document collapse, citation hydration, and metric calculation.

The implementation is intentionally separate from representation generation.
It must first make reranker identity, candidate-window size, endpoint policy,
latency, and score output observable in the immutable run record. The first
target runtime is `llama-server` on the Mac, accessed through a private SSH
loopback tunnel; Ollama remains the embedding provider and is not treated as a
native reranker API.

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- [Design and implementation guide](design-doc/01-reranker-stage-analysis-design-and-implementation-guide.md)
- [Implementation diary](reference/01-implementation-diary.md)
- [llama.cpp server reranking API](https://github.com/ggml-org/llama.cpp/blob/master/tools/server/README.md)

## Status

Current status: **active**

## Topics

- rag
- reranking
- ttc
- geppetto
- ollama
- rag-eval

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
