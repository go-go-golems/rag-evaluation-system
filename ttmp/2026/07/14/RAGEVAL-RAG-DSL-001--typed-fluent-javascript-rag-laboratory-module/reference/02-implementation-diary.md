---
Title: Implementation diary
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
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Chronological record of the contract-first design work for the typed RAG laboratory JavaScript module.
LastUpdated: 2026-07-14T22:09:50.352697004-04:00
WhatFor: Preserve decisions, evidence, and the next implementation steps for the ticket.
WhenToUse: Read before resuming implementation or reviewing why a public API decision was made.
---

# Implementation diary

## 2026-07-14 — Step 1: establish scope and evidence

**Goal.** Convert the exploratory JavaScript-playground discussion into a
durable module contract without pretending that the module already exists.

**Evidence examined.**

- `cmd/rag-eval/xgoja.yaml` selects generic `db`, `fs`, `markdown`,
  `geppetto`, and related modules for the generated `rag-eval-js` runtime.
- `cmd/rag-eval/jsverbs/database.js` and `explorer.js` demonstrate useful but
  untyped, SQL-oriented exploration verbs. They are not an experiment API.
- `docs/howtos/how-to-write-rag-eval-js-scripts.md` documents the current
  generic runtime and its build boundary.
- `RAGEVAL-TTC-LAB-001` defines immutable corpus/artifact/specification/run
  identities and already owns the current laboratory persistence model.
- The Widget DSL and researchctl use nested configurator lambdas and `.use()`
  fragments. The transcript prototype demonstrates channel retrieval, RRF,
  parent collapse, and source hydration.

**Decision.** The new module is named `rag`, exposes `rag.open(...)`, and
compiles authoring-time builder operations into the canonical immutable
experiment specification. It has no hidden database mutation during
`.toSpec()` or `.validate()`.

**Result.** Created this ticket, its task list, a normative API reference, and
an intern-oriented design/implementation guide. No application code or
experiments were written in this design step, therefore no ticket-local script
was needed.

**Next.** Confirm the concrete persisted specification structure in the Go
service, then implement the pure Go domain builder and test it before adding a
goja adapter.

## 2026-07-14 — Step 2: validate and publish the design package

`docmgr validate frontmatter` succeeded for the API reference and design guide;
`docmgr doctor --ticket RAGEVAL-RAG-DSL-001 --stale-after 30` reported all
checks passed. A restricted-network upload could not resolve the reMarkable
cloud host. Retrying the exact bundle with external network permission
succeeded at `/ai/2026/07/14/RAGEVAL-RAG-DSL-001/RAG Laboratory JavaScript
Module Design.pdf`. The bundle contains the ticket index, specification, guide,
and diary. No code or ticket-local experiment script was written in this
documentation phase.
