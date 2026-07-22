---
Title: Typed fluent JavaScript RAG laboratory module
Ticket: RAGEVAL-RAG-DSL-001
Status: active
Topics:
    - rag
    - rag-eval
    - dsl
    - fluent-builder
    - goja
    - xgoja
    - javascript
    - typescript
    - intern-guide
    - playground
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://rag-evaluation-system/docs/guides/ttc-rag-laboratory.md
      Note: Operational context and inspection guide for the immutable laboratory
    - Path: repo://rag-evaluation-system/ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/design-doc/01-ttc-rag-laboratory-baseline-and-immutable-experiment-runs-design-and-implementation-guide.md
      Note: Parent design that establishes artifact and append-only run invariants
ExternalSources: []
Summary: Normative design ticket for a typed Go-backed require("rag") JavaScript module that authors reproducible RAG laboratory specifications and starts append-only runs.
LastUpdated: 2026-07-14T22:09:39.511605621-04:00
WhatFor: Establish the API, implementation boundary, delivery plan, and validation criteria for the reusable RAG laboratory DSL.
WhenToUse: Read before adding a JavaScript RAG primitive, changing immutable experiment specifications, or packaging rag-eval functionality into xgoja.
---


# Typed fluent JavaScript RAG laboratory module

## Overview

This ticket specifies a new `require("rag")` module for `rag-eval-js`. It is
not a generic vector-store wrapper and it does not replace the Go HTTP API.
It is an opinionated, typed authoring language for the RAG laboratory: a
script names immutable inputs, composes retrieval and representation policies,
selects evaluation metrics, and produces a canonical experiment
specification. An explicit execution operation persists that specification and
creates a new append-only run.

The ticket intentionally follows the builder/lambda/fragment conventions of
the Widget DSL and researchctl while retaining a different execution model:
RAG lambdas configure typed Go builders; they are never persisted as source
code or arbitrary callbacks. The persisted artifact is canonical JSON whose
fingerprint is the experiment identity.

The work is a hard cutover for this new module. It neither wraps nor evolves
the current ad-hoc `db` JavaScript verbs into a durable public API.

## Key Links

- [Normative API specification](reference/01-rag-laboratory-javascript-module-api-specification.md)
- [Design and implementation guide](design-doc/01-typed-fluent-rag-module-design-and-implementation-guide.md)
- [Implementation diary](reference/02-implementation-diary.md)
- [Existing TTC laboratory guide](../../2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/design-doc/01-ttc-rag-laboratory-baseline-and-immutable-experiment-runs-design-and-implementation-guide.md)

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- rag
- rag-eval
- dsl
- fluent-builder
- goja
- xgoja
- javascript
- typescript
- intern-guide
- playground

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
