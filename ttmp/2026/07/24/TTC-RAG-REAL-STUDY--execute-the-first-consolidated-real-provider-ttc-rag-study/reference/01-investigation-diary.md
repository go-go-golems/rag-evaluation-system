---
Title: Investigation diary
Ticket: TTC-RAG-REAL-STUDY
Status: active
Topics:
    - ttc
    - rag
    - evaluation
    - research
    - workflow
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/code/wesen/go-go-golems/go-go-parc/Projects/2026/07/24/ARTICLE - RAG Experiment JavaScript - Language API and End-to-End Execution.md
      Note: Supplemental source supplied by user
    - Path: repo://experiments/real-provider-v2/base.js
      Note: Evidence reviewed during study architecture design
    - Path: repo://experiments/real-provider-v2/study-full.js
      Note: Evidence reviewed for full provider path
ExternalSources:
    - /home/manuel/code/wesen/go-go-golems/go-go-parc/Projects/2026/07/24/ARTICLE - RAG Experiment JavaScript - Language API and End-to-End Execution.md
Summary: Chronological planning and research record for the consolidated real-provider TTC RAG study.
LastUpdated: 2026-07-24T07:20:00-04:00
WhatFor: Preserve why the study is structured in phases, which evidence shaped it, and how an engineer should continue implementation.
WhenToUse: Read before resuming TTC-RAG-REAL-STUDY or changing its architecture, inputs, phases, budgets, or claim boundaries.
---


# Investigation diary

## Goal

Design and execute the first consolidated real-provider TTC RAG study from readable human source through immutable Researchctl, Workflow V3, provider-operation, and analysis custody. This diary records the study's architecture decisions, source evidence, implementation sequence, failures, validation, and final claims.

## Step 1: Create the real-study workstream and implementation guide

This step created a dedicated ticket for the first real TTC study rather than extending the completed architecture-acceptance tickets. The ticket separates scientific execution work from platform convergence and from the deferred Researchctl PR review.

The initial design turns the current readable but fragmented examples into one intern-oriented project plan. It defines immutable inputs, four retrieval cases, a linked answer study, phased provider qualification, evidence inspection, and explicit claim limits.

### Prompt Context

**User prompt (verbatim):** "ok, create a new ticket for the study. Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.

. We'll handle the PR review later."

**Assistant interpretation:** Create a repository-local docmgr ticket and an exhaustive onboarding/design/implementation guide for the real TTC RAG study, but do not fix PR #3 in this tranche. Validate the ticket and upload the documentation bundle to reMarkable.

**Inferred user intent:** Establish a durable, reviewable plan that a new intern can follow to consolidate the scripts and safely execute the real provider-backed study without repeating architecture mistakes or making premature benchmark claims.

### What I did

- Ran `docmgr status --summary-only` and inspected the repository vocabulary.
- Created ticket `TTC-RAG-REAL-STUDY` with topics `ttc,rag,evaluation,research,workflow,intern-guide`.
- Added ten tasks covering Researchctl prerequisites, source/generated separation, immutable inputs, source consolidation, P0–P3 execution, acceptance, and closure.
- Added the primary design document and investigation diary.
- Wrote the intern guide with system ownership, vocabulary, architecture diagrams, APIs, pseudocode, file references, decision records, phased implementation, tests, risks, open questions, and review checklists.

### Why

- The architecture is ready, but the human study source remains fragmented and the real-provider examples still contain placeholders.
- Provider-backed execution requires stronger cost, privacy, identity, and claim gates than a fixture acceptance run.
- A new engineer needs a dependency-ordered explanation before editing study source or interpreting generated plans.

### What worked

- The existing vocabulary already contained every required topic.
- Docmgr created the complete ticket workspace, design document, diary, tasks, index, and changelog without manual directory setup.
- Existing source files provided concrete APIs for every major stage from corpus binding through cited answer generation.

### What didn't work

- N/A. This step was documentation and planning only. The deferred PR #3 findings remain an explicit prerequisite rather than an attempted change.

### What I learned

- The first real study should compare retrieval incrementally rather than run one opaque all-real configuration.
- Preparation reuse is necessary to isolate retrieval differences and control provider cost.
- Answer generation should consume frozen hydrated evidence in a linked study so prompt/model changes do not rerun retrieval.
- The generated `researchctl-plan.js` is valid custody but unsuitable as the primary review surface.

### What was tricky to build

The design had to preserve three different meanings of execution. A RAG cell defines domain semantics, a Researchctl run defines the scientific sample, and Workflow V3 task attempts define operational retries. Mixing these levels would produce incorrect replicate counts and ambiguous provider accounting. The guide therefore introduces the levels before presenting commands or analysis.

The guide also had to separate real-provider qualification from benchmark validity. A successful full run proves execution and evidence custody; it does not prove that the dataset, replicates, holdout, or provider freeze support a benchmark claim.

### What warrants a second pair of eyes

- Whether `rag.units.identity()` matches the fixed-truth relevance target.
- Whether the initial multi-representation fusion should retain `raw.vector: 2` or use equal weights.
- Whether three P2 replicates are justified after immutable preparation reuse.
- Whether the current projected metric names cover every proposed analysis reducer.
- Whether the generated-plan loader design belongs in PR #3 or a separate focused change.

### What should be done in the future

- Address the three Researchctl PR #3 provenance/validation findings before P0.
- Resolve real corpus/evaluation manifests and approved provider ceilings.
- Implement the project structure and execute P0 only after prerequisite review.

### Code review instructions

- Start with the executive summary and Sections 4–7 of the design guide.
- Review the four-case matrix in Sections 9–10 and the P0–P3 gates in Sections 17–20.
- Compare API sketches with `experiments/real-provider-v2/base.js`, `study-full.js`, and `product.js`.
- Verify that no code or provider call is part of this documentation tranche.

### Technical details

Ticket root:

`ttmp/2026/07/24/TTC-RAG-REAL-STUDY--execute-the-first-consolidated-real-provider-ttc-rag-study/`

Primary guide:

`design-doc/01-real-provider-ttc-rag-study-architecture-and-implementation-guide.md`

## Step 2: Incorporate the RAG JavaScript language and execution article

This step incorporated a comprehensive vault article supplied by the user. The article explains the Goja-backed RAG authoring language, complete execution path, APIs, data binding, Workflow lowering, provider operations, Researchctl scheduling, analysis, and scientific failure modes.

The design guide uses the article as onboarding context while keeping repository contracts and source files authoritative. It adds the article to external sources and translates its broad language reference into study-specific decisions, phase gates, and implementation tasks.

### Prompt Context

**User prompt (verbatim):** "you can potentialy use /home/manuel/code/wesen/go-go-golems/go-go-parc/Projects/2026/07/24/ARTICLE - RAG Experiment JavaScript - Language API and End-to-End Execution.md as information too"

**Assistant interpretation:** Read the full article and use its language/API/execution explanations as source material for the intern guide.

**Inferred user intent:** Avoid duplicating incomplete investigation and ensure the study guide reflects the most complete current explanation of how authored RAG JavaScript reaches real execution and analysis.

### What I did

- Read the article completely in two passes because the first tool response reached the 50 KB output limit.
- Incorporated its explanations of authoring callbacks, hidden Go descriptors, pipelines, query plans, studies, products, immutable inputs, lowering, runner protocol, Workflow execution, external operations, and analysis.
- Added the article as an explicit external source in the design guide and diary.
- Used its proposed readable experiment structure as the basis for `experiments/ttc-real/`.
- Converted its general progression into study-specific P0–P3 phases and acceptance criteria.

### Why

- The article is the strongest existing onboarding source for the RAG JavaScript language.
- It records important scientific rules: generated text is not evidence, collapse precedes fusion, query metrics aggregate within runs, and generated plans are custody rather than author source.
- Reusing that evidence allows the ticket guide to focus on the concrete real-study design.

### What worked

- The article matched the completed convergence architecture and current repository files.
- Its code examples map directly to the proposed study project.
- It identified a TypeScript declaration parity gap for `combinedSummaryQuestions`, which is now an implementation prerequisite.

### What didn't work

- The first `read` result was truncated at 50 KB. I continued from line 1498 and read the remaining content before drafting conclusions.

### What I learned

- The language guide already provides the correct human review order: README, pipeline, query, study, analysis, inputs, generated manifest, then generated plans.
- The real-provider candidate separates retrieval-only and full rerank/answer paths, which supports the proposed P0/P1/P2/P3 decomposition.
- The generic analysis engine already encodes the correct `sampleUnit: run` behavior.

### What was tricky to build

The external article describes both current behavior and suggested future UX. The design guide distinguishes observed APIs from proposed consolidation. In particular, manifest-backed generated execution loading remains a proposal because the current generator still emits embedded canonical domain configurations.

### What warrants a second pair of eyes

- Verify every API listed in Section 30 against current TypeScript declarations and native module exports.
- Confirm which article statements are current contract versus recommended cleanup before implementation begins.

### What should be done in the future

- Keep the article and ticket guide cross-referenced as APIs evolve.
- Update the article when the generated-plan separation is implemented.

### Code review instructions

- Compare the design guide's API section with the article's Sections 3–19.
- Compare the execution sections with the article's Sections 24–35.
- Compare the phase and review checklists with the article's Sections 37 and 45.

### Technical details

External source:

`/home/manuel/code/wesen/go-go-golems/go-go-parc/Projects/2026/07/24/ARTICLE - RAG Experiment JavaScript - Language API and End-to-End Execution.md`

## Step 3: Validate and publish the intern guide to reMarkable

This step validated the ticket metadata and delivered the guide, diary, and task list as one reMarkable PDF with a table of contents. The ticket remains active because this tranche produced the implementation design rather than executing the study.

The dry run showed the exact input files, PDF name, and remote directory. The real upload then completed successfully without requiring manual authentication recovery.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Validate the complete ticket and upload the documentation bundle to the requested reMarkable destination.

**Inferred user intent:** Make the technical guide available for offline review before implementation or provider spend begins.

### What I did

- Ran `docmgr validate frontmatter` for the design guide and diary; both passed.
- Ran `git diff --check`; it passed.
- Ran targeted `docmgr doctor --ticket TTC-RAG-REAL-STUDY --stale-after 30 --details`.
- Corrected one stale relation from nonexistent `pkg/ragworkflow/plan.go` to current `pkg/ragworkflow/lower.go`.
- Reran doctor; all checks passed.
- Ran a reMarkable bundle dry run for the design guide, diary, and tasks.
- Uploaded `TTC RAG Real Study Intern Guide.pdf` to `/ai/2026/07/24/TTC-RAG-REAL-STUDY`.

### Why

- Frontmatter and relation validation ensure the guide remains discoverable and does not point reviewers to nonexistent current files.
- One bundled PDF gives the intern the architecture, chronological context, and implementation checklist together.

### What worked

- Frontmatter validation passed on the first run.
- The reMarkable dry run showed the intended files and destination.
- The upload returned: `OK: uploaded TTC RAG Real Study Intern Guide.pdf -> /ai/2026/07/24/TTC-RAG-REAL-STUDY`.

### What didn't work

- Initial doctor reported `missing_related_file` for `repo://pkg/ragworkflow/plan.go`. The lowering implementation is actually `pkg/ragworkflow/lower.go`. I removed the invalid relation, added the current file, and doctor passed.

### What I learned

- Related-file metadata catches filename assumptions that prose-only references may not expose.
- The study ticket can remain implementation-ready while every execution task stays open.

### What was tricky to build

The upload bundle includes a 58 KB design guide with multiple Mermaid diagrams, JavaScript examples, tables, and nested checklists. A dry run was necessary to verify bundle composition and naming before rendering/upload.

### What warrants a second pair of eyes

- Review the generated PDF's diagram and table readability on the device.
- Review the open decisions before any task is checked or provider budget approved.

### What should be done in the future

- After design review, record accepted decisions and update the guide before implementation.
- Re-upload with `--force` only if replacing the existing PDF is explicitly intended, because replacement removes annotations.

### Code review instructions

- Run `docmgr doctor --ticket TTC-RAG-REAL-STUDY --details`.
- Review `git diff --check` and confirm all implementation tasks remain open.
- Verify the successful upload output recorded above; no routine cloud listing is required.

### Technical details

Remote destination:

`/ai/2026/07/24/TTC-RAG-REAL-STUDY/TTC RAG Real Study Intern Guide.pdf`
