---
Title: Implementation Diary
Ticket: RAGEVAL-003
Status: active
Topics:
    - rag
    - playground
    - corpus
    - embeddings
    - search
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ttmp/2026/05/28/RAGEVAL-003--user-centered-corpus-explorer-and-rag-pipeline-visualization/design-doc/01-corpus-explorer-and-pipeline-visualization-implementation-guide.md
      Note: Design guide authored in Step 1
ExternalSources: []
Summary: Chronological diary for the Corpus Explorer and pipeline visualization design ticket.
LastUpdated: 2026-05-28T00:00:00Z
WhatFor: Track design decisions, evidence, and validation for the website improvement planning work.
WhenToUse: Read before implementing RAGEVAL-003 or updating the Corpus Explorer design.
---


# Diary

## Goal

This diary records the design workflow for adding a user-centered Corpus Explorer and RAG pipeline visualization to the website.

## Step 1: Ticket Creation and Intern-Ready Design Guide

I created a new docmgr ticket for improving the website from the perspective of a user who wants to learn, experiment with, and validate the RAG pipeline. I inspected the current frontend and API surfaces, then wrote a detailed implementation guide that explains the system, user intents, proposed information architecture, API contracts, frontend slices, validation plan, and design constraints.

The guide deliberately leaves visual freedom for a designer while being strict about technical clarity. The user should be able to see what was ingested, what each document became, where chunk boundaries are, which chunks have embeddings, and what can be tested next.

### Prompt Context

**User prompt (verbatim):** "ok, create a new docmgr ticket for adding this to the website / improving the website, starting from a user intent perspective as someone who wants to learn and experiment and validate the RAG pipeline here. 

Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.


Let enough freedom for a great designer to improve the site, and focus on like user intent and value"

**Assistant interpretation:** Create a new ticket for a user-centered website improvement, write an intern-ready design/implementation guide for a Corpus Explorer and pipeline visualization, and upload the package to reMarkable.

**Inferred user intent:** Turn the backend corpus/chunk/embedding pipeline into an understandable visual product surface rather than just CLI/API functionality.

### What I did

- Created ticket `RAGEVAL-003` titled `User-Centered Corpus Explorer and RAG Pipeline Visualization`.
- Added tasks for current-state analysis, user-intent design, implementation guide writing, and reMarkable upload.
- Inspected current frontend files:
  - `web/src/App.tsx`
  - `web/src/components/retro/MacMenuBar.tsx`
  - `web/src/components/pipeline/PipelineView.tsx`
  - `web/src/components/embeddings/EmbeddingsView.tsx`
  - `web/src/services/api.ts`
- Inspected current API/backend files:
  - `internal/api/handlers.go`
  - `internal/db/queries.go`
- Reused current corpus context from `RAGEVAL-002`, especially the dump-to-SQLite design guide and imported TTC corpus state.
- Wrote `design-doc/01-corpus-explorer-and-pipeline-visualization-implementation-guide.md`.

### Why

The system now has real data and useful backend operations, but users need a website surface that explains the pipeline. The design must start from user intent because the goal is not to expose every table; the goal is to help users understand, experiment with, and validate each transformation stage.

### What worked

- The current code already has clear places to extend:
  - add `CorpusExplorerView` under `web/src/components/corpus/`;
  - add a `corpus` menu item in `MacMenuBar`;
  - add RTK Query bindings in `web/src/services/api.ts`;
  - add corpus-specific HTTP endpoints in `internal/api/handlers.go`;
  - add a shared corpus service under `internal/services/corpus/`.
- The current embedding coverage and source-aware compute work provides the first useful coverage data model for the Corpus Explorer.

### What didn't work

- No website implementation was done in this step. This ticket is a design and implementation guide package, not a code implementation ticket.
- The current generic document API is insufficient for the proposed UI because it lacks source filtering, document-level chunk counts, content text, metadata, and chunk-level embedding status in one payload.

### What I learned

- The first useful visualization is not a chart-heavy dashboard. It is a Corpus Explorer that makes source, document, chunk, and embedding identities visible.
- The best product value is helping a user answer validation questions: what was ingested, what changed, how it was chunked, what is embedded, and what can be tested next.
- A designer should have freedom over layout and visual style, but exact IDs and embedding identity tuples must remain visible and copyable.

### What was tricky to build

The main challenge was balancing design freedom with technical specificity. If the document over-prescribes layout, it limits good design work. If it only describes user experience goals, it is not enough for an intern to implement. The guide therefore specifies required data, endpoints, files, validation commands, and component responsibilities, while leaving visual composition open.

### What warrants a second pair of eyes

- Review whether the proposed corpus API should be a separate route group or folded into the existing document endpoints.
- Review whether the first UI slice should be read-only or include bounded compute actions.
- Review whether product metadata should be copied into `documents.content_text` before the Corpus Explorer highlights missing searchable fields.

### What should be done in the future

- Implement the backend corpus service and endpoints.
- Implement the first read-only Corpus Explorer frontend slice.
- Add source-level embedding coverage to the Embedding Inspector or cross-link it from Corpus Explorer.
- Add search-result-to-corpus navigation after BM25 search exists.

### Code review instructions

- Start with the design guide:
  - `ttmp/2026/05/28/RAGEVAL-003--user-centered-corpus-explorer-and-rag-pipeline-visualization/design-doc/01-corpus-explorer-and-pipeline-visualization-implementation-guide.md`
- Then inspect the current files referenced by the guide:
  - `web/src/App.tsx`
  - `web/src/components/retro/MacMenuBar.tsx`
  - `web/src/services/api.ts`
  - `web/src/components/pipeline/PipelineView.tsx`
  - `web/src/components/embeddings/EmbeddingsView.tsx`
  - `internal/api/handlers.go`
  - `internal/db/queries.go`

### Technical details

- Ticket: `RAGEVAL-003`.
- Primary design document: `design-doc/01-corpus-explorer-and-pipeline-visualization-implementation-guide.md`.
- Proposed first endpoint group: `/api/v1/corpus/*`.
- Proposed first frontend view: `web/src/components/corpus/CorpusExplorerView.tsx`.

---

## Step 2: reMarkable Upload and Ticket Completion

I uploaded the RAGEVAL-003 design package to reMarkable as a bundled PDF. The bundle includes the intern-ready Corpus Explorer implementation guide and this diary.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Complete the requested documentation delivery by publishing the design package to reMarkable.

**Inferred user intent:** Make the website design plan reviewable away from the codebase, with a readable PDF bundle and table of contents.

### What I did

- Uploaded the design document and diary with `remarquee upload bundle`.
- Used remote directory `/ai/2026/05/28/RAGEVAL-003`.
- Marked task 4 complete.
- Updated the changelog with the upload location.

### Why

The user explicitly requested upload to reMarkable. A bundled PDF keeps the design guide and diary together for review.

### What worked

- Upload succeeded:
  - `OK: uploaded RAGEVAL-003 Corpus Explorer Design.pdf -> /ai/2026/05/28/RAGEVAL-003`

### What didn't work

- N/A.

### What I learned

- The design package is ready for review and implementation planning.

### What was tricky to build

N/A for this upload step. The only operational constraint was to avoid redundant reMarkable status/list calls and use the direct non-interactive upload command.

### What warrants a second pair of eyes

- Review the design guide before implementation begins, especially the proposed `/api/v1/corpus/*` endpoint boundaries and first read-only UI slice.

### What should be done in the future

- Implement the first read-only Corpus Explorer slice under a follow-up implementation task.
- Add frontend coverage wiring after the corpus endpoints exist.

### Code review instructions

- Review the uploaded PDF or the source markdown:
  - `design-doc/01-corpus-explorer-and-pipeline-visualization-implementation-guide.md`

### Technical details

- reMarkable path: `/ai/2026/05/28/RAGEVAL-003`.
- Uploaded bundle name: `RAGEVAL-003 Corpus Explorer Design.pdf`.
