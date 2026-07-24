---
Title: Investigation diary
Ticket: RAG-INTAKE-WORKFLOW-V3-CUTOVER
Status: active
Topics:
    - rag
    - rag-eval
    - intake
    - scraper
    - workflow
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/scraper/pkg/workflowv3product/service.go
      Note: Canonical V3 product submission/read/cancel/worker APIs
    - Path: repo://cmd/rag-eval/cmds/intake/root.go
      Note: Current product-local intake CLI boundary
    - Path: repo://internal/api/workflow_artifact_handlers.go
      Note: Current old engine-view API and submission surface
    - Path: repo://internal/workflow/intake_runner.go
      Note: Current RAG service dispatch and provider resolution semantics
    - Path: repo://internal/workflow/submit.go
      Note: Current old-engine intake plan construction and selection behavior
ExternalSources: []
Summary: Chronological implementation and validation record for the RAG intake Workflow V3 hard cut and its path to Scraper legacy cleanup.
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---


# Investigation diary

## Goal

Record every design, implementation, failure, validation, and deletion step required to migrate RAG document intake to Workflow V3 and continue through `SCRAPER-LEGACY-CLEANUP`.

## Step 1: Create the missing ticket and establish the migration boundary

The new ticket fills the explicit gap identified after the RAG v2 study execution cutover. I first re-audited the live callers rather than treating every remaining old-Scraper import as stale. The surviving code implements active document preprocessing, chunking, enrichment, embeddings, BM25 indexing, command workers, API projections, and frontend views.

The implementation guide therefore defines a hard cut with behavior preservation, not immediate deletion. It assigns durable execution and read models to Workflow V3 while keeping RAG request schemas, task semantics, providers, and database services domain-owned.

### Prompt Context

**User prompt (verbatim):** "create the missing ticket and add a design / implementation guide (doesnt need to be super elaborate), then work your way all the way up to SCRAPER-LEGACY-CLEANUP, keeping a detailed diary as you work. You have permission to hit external sources when testing."

**Assistant interpretation:** Create and document the missing intake cutover, implement it fully, then complete every required dependency and cleanup through Scraper's legacy-engine deletion, using external services for authorized acceptance and maintaining detailed evidence.

**Inferred user intent:** Remove the final cross-repository blockers rather than merely inventorying them, while preserving product behavior and leaving a reviewable migration record.

### What I did
- Created `RAG-INTAKE-WORKFLOW-V3-CUTOVER` with five implementation/acceptance tasks.
- Added a concise but implementation-ready design guide covering current/target flows, schemas, task identities, plan construction, provider custody, APIs, deletion gates, sequencing, and validation.
- Re-read the completed RAG execution cutover evidence and current Scraper legacy/Workflow V3 task status.
- Inspected all remaining RAG intake command and `internal/workflow` files, old API routes, frontend consumers, and Scraper Workflow V3 product APIs.

### Why
- RAG-eval remains an external consumer of Scraper's old engine, so `SCRAPER-LEGACY-CLEANUP` cannot safely delete it.
- The intake path mutates domain databases and contacts providers; a blind mechanical conversion would risk duplicated writes, secrets in plans, or incorrect retry semantics.

### What worked
- Scraper already exposes `workflowv3product.Application` for submission, list/show/cancel, worker operation, canonical observations, restart-safe SQLite, and domain-owned task packages.
- Existing RAG services isolate chunking, preprocessing, enrichment, embedding, and BM25 semantics from scheduler persistence, making them reusable behind a V3 task module.

### What didn't work
- N/A. This first step was inventory and design; no runtime implementation was attempted.

### What I learned
- The earlier 15 downstream legacy files are now eight direct old-Scraper importers after the study cutover.
- Manual per-operation retry is an old-engine UI/API concept and should not be copied; V3 retry policy is immutable plan identity, and reruns should be explicit new runs.
- The broader Scraper legacy ticket still has three unfinished Workflow V3 slices, including stronger isolation and a real TTC workload, in addition to product/site/API migration.

### What was tricky to build
- The ticket must separate RAG intake completion from Scraper-wide deletion. Completing the RAG package removes one external blocker but does not justify deleting Scraper root commands, site infrastructure, JavaScript runtime, metrics, events, and UI before their own dispositions pass.

### What warrants a second pair of eyes
- Review explicit-node versus lazy-map plan construction for database-mutating document work.
- Review whether V3 intake APIs should expose raw snapshots, canonical observations, or a smaller stable RAG projection.
- Review provider-operation custody for the existing embedding service, which currently performs provider calls internally.

### What should be done in the future
- Characterize exact existing database/output behavior in focused tests.
- Implement `pkg/ragintakeworkflow` and migrate one deterministic vertical slice before deleting old code.

### Code review instructions
- Start with the design guide's current/target traces.
- Compare `internal/workflow/submit.go`, `internal/workflow/intake_runner.go`, and Scraper `pkg/workflowv3product/service.go`.
- Validate ticket structure with `docmgr doctor --ticket RAG-INTAKE-WORKFLOW-V3-CUTOVER --stale-after 30`.

### Technical details
- New package identity: `rag-intake-v1@1.0.0`.
- Proposed request schema: `rag-intake-request/v1`.
- Remaining direct old-Scraper importers: eight files in intake commands, API, and `internal/workflow`.
