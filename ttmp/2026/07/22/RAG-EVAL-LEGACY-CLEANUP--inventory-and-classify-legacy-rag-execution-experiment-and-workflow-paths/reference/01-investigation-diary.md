---
Title: Investigation diary
Ticket: RAG-EVAL-LEGACY-CLEANUP
Status: active
Topics:
    - rag
    - rag-eval
    - workflow
    - research
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-22T21:37:42.840060092-04:00
WhatFor: ""
WhenToUse: ""
---

# Investigation diary

## Goal

Inventory and classify legacy paths in RAG-eval under the `EXPERIMENT-PLATFORM-CONVERGENCE` umbrella, without deleting code during this investigation tranche.

## Step 1: Create the cleanup workstream

This step created a repository-local cleanup ticket with a ten-item task list. The three cleanup tickets separate repository ownership while sharing one classification vocabulary and one umbrella.

The stopping condition is explicit: identify canonical code, immediate safe removals, and deferred removals with named replacement gates; document and commit the result; do not execute deletion yet.

### Prompt Context

**User prompt (verbatim):** "Ok, create a separate ticket for the cleanup, and a detailed list of tasks. Then, start working on the cleanup, keep a detaile ddiary as you work.

Once you have identified what can immediately be removed, what should bre placed later on and how, stop. Make sure to keep a detailed document about the cleanup in each repositories respective ticket, and cross link all the files and tickets so that we don't forget we are working under a common umbrella down the road.

commit at appropriat eintervals."

**Assistant interpretation:** Create one cleanup ticket per repository, audit active and legacy paths, cross-link everything, commit at useful milestones, and stop after disposition planning.

**Inferred user intent:** Reduce architecture work before adding capabilities, while avoiding blind deletion and preserving a durable shared plan.

### What I did
- Created `RAG-EVAL-LEGACY-CLEANUP` and a detailed task list.
- Linked the workstream conceptually to `EXPERIMENT-PLATFORM-CONVERGENCE` and the two sibling cleanup tickets.
- Defined retain, remove-now, and remove-after-replacement classifications.

### Why
- Cleanup decisions need repository-local evidence and ownership.
- A shared umbrella prevents three independent cleanup efforts from producing incompatible boundaries.

### What worked
- Docmgr created the ticket, primary cleanup document, diary, tasks, index, and changelog.

### What didn't work
- N/A in ticket setup.

### What I learned
- The cleanup stage needs an explicit stop before deletion so classification can be reviewed.

### What was tricky to build
- “Legacy” cannot mean merely old. The ticket therefore requires active-reference and replacement evidence for every disposition.

### What warrants a second pair of eyes
- Review the eventual remove-now list before destructive work begins.

### What should be done in the future
- Complete the source inventory and disposition report, then pause for review.

### Code review instructions
- Start with `tasks.md`, then read the cleanup inventory document and later diary steps.

### Technical details
- Umbrella: `EXPERIMENT-PLATFORM-CONVERGENCE`.
- Siblings: `RESEARCHCTL-LEGACY-CLEANUP`, `SCRAPER-LEGACY-CLEANUP`, `RAG-EVAL-LEGACY-CLEANUP`.

## Step 2: Classify RAG execution generations and immediate removals

This step traced the canonical RAG v2 core and three execution generations: old intake workflow, direct rag-worker/ragengine, and TTC-specific Workflow V3 sweep. Unlike Scraper, RAG has a meaningful immediate deletion tranche because the TTC runner is isolated and its experiment evidence is already published.

The audit also discovered a repository-wide test failure unrelated to production semantics: several historical ticket scripts share one Go package and each declare `main`. This should be fixed as immediate cleanup while preserving the scripts as evidence.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Identify immediate and deferred RAG cleanup, preserve semantic oracles, document test failures, and stop before deleting.

**Inferred user intent:** Remove completed one-off experiment machinery early without losing the current behavior needed to build the clean replacement.

### What I did
- Mapped all `rag-eval` commands and imports of `workflowv3ttc`, `preparationworkflow`, old `workflow`, `ragengine`, and `researchctladapter`.
- Verified TTC non-test references are isolated to the standalone command and package.
- Traced `OperationCustody` symbols to the TTC sweep only.
- Confirmed `EchoRunner` is test-only while production registers `IntakeRunner`.
- Read the old chunks-schema compatibility migration and its call from `db.go`.
- Ran `GOWORK=off go test ./... -count=1`.

### Why
- RAG has both removable experiment scaffolding and active semantic/runtime code; they require different disposition.

### What worked
- All active production packages, including the TTC runner and current RAG worker, passed their package tests.
- Search produced clear immediate deletion boundaries.

### What didn't work
- The full test command failed because ticket-local scripts are discovered as one Go package:
  - `06-refresh-model-manifest-digests.go:18:6: main redeclared in this block`
  - `07-refresh-prompt-manifest-digests.go:18:6: main redeclared in this block`
  - `08-build-operation-custody-export.go:37:6: main redeclared in this block`
  - each conflicts with `02-build-researchctl-custody-spec.go:11:6`.
- The first documentation staging check also failed on prompt trailing whitespace and a blank changelog EOF; normalization fixed the second attempt.

### What I learned
- The TTC binary/package and post-hoc custody adapter can be removed immediately after extracting small regression fixtures.
- `ragengine`, `rag-worker`, preparation, and intake must remain until their named Workflow V3 replacements pass.
- Compatibility migration code contradicts the documented disposable-database hard cut and can be removed now.

### What was tricky to build
- `pkg/researchctladapter` mixes generic current behavior with TTC-specific custody. Classification had to occur file-by-file rather than deleting the package.

### What warrants a second pair of eyes
- Approve the ticket-script archival convention (`.go.txt` versus one directory per executable script).
- Confirm old developer databases may be discarded without exception.
- Select exact TTC ledger fixtures to preserve before code deletion.

### What should be done in the future
- Review the five-item immediate tranche, then execute it as a focused cleanup commit and run the full suite.

### Code review instructions
- Begin with the immediate removal table in the cleanup report.
- Inspect `cmd/rag-ttc-v3-sweep`, `internal/workflowv3ttc`, `operation_custody.go`, `echo_runner.go`, and `internal/db/migrations.go`.

### Technical details
- Baseline result: production packages passed; overall `go test ./...` failed only in historical ticket scripts.
- Setup commit: `500cd0924f08df57724f9ca7e51405a2e362e971`.

## Step 3: Validate, commit, and stop before deletion

The inventory document, detailed task state, related-file evidence, and cross-ticket navigation were validated and committed. This is the requested stopping point: disposition is complete, but no removal tranche has been executed.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Stop after classification and leave destructive tasks open for review.

**Inferred user intent:** Make deletion a deliberate reviewed follow-up rather than an uninterrupted audit-and-delete operation.

**Commit (code):** `292756d310af021656ec5bbfbea4652d0129521e` — cleanup inventory and classification documentation.

### What I did
- Checked the first seven inventory/classification tasks.
- Left review, immediate deletion, and deferred hard-cut tasks open.
- Ran docmgr doctor successfully.
- Cross-linked the cleanup report to key source files, the umbrella, sibling cleanup tickets, and replacement tickets.

### Why
- The remove-now list needs review before destructive edits.

### What worked
- Docmgr validation passed.
- Baseline test result: production packages passed; overall suite failed only because historical ticket scripts redeclare main.

### What didn't work
- No additional failure beyond those recorded in Step 2.

### What I learned
- Repository-local classifications can share one common umbrella without hiding different readiness levels.

### What was tricky to build
- The stopping point had to preserve enough evidence for a later agent to execute cleanup without rerunning the entire investigation. File references, deletion gates, exact commands, and open tasks provide that continuation state.

### What warrants a second pair of eyes
- Approve the remove-now table and any stated assumption about external users or disposable state.

### What should be done in the future
- After review, check the review task, execute only the immediate tranche, run validation, and commit it separately. Deferred paths remain until their named replacement tickets pass.

### Code review instructions
- Read the report's executive summary, immediate removal section, deferred removal section, and review checklist.
- Confirm tasks 8–10 remain open.

### Technical details
- Ticket: `RAG-EVAL-LEGACY-CLEANUP`.
- Stop condition reached: classification complete; deletion not started.

## Step 4: Execute the approved RAG-eval removal tranche

The approved immediate tranche removed completed TTC execution machinery and obsolete compatibility code while preserving historical evidence and active intake behavior.

### What I did
- Deleted `cmd/rag-ttc-v3-sweep` and `internal/workflowv3ttc`.
- Deleted TTC-only operation custody code/tests.
- Deleted `EchoRunner` and its tests.
- Removed the old chunks-schema upgrade and its upgrade test.
- Renamed four historical standalone scripts from `.go` to `.go.txt`.
- Preserved the active echo operation response as `IntakeOutput` and moved the shared test JSON helper into `intake_runner_test.go`.
- Committed as `8a33a613a0661c62b91eb98177efd3e2958a280d`.

### Validation and failure trace
- The first focused test exposed that `EchoOutput` and `mustJSON` were shared by active intake code/tests. Rather than restore the compatibility runner, I moved those two live semantics to the intake files; focused tests then passed.
- Full `GOWORK=off go test ./... -count=1` passed, including ticket directories.
- `make build` and sequential `make lint` passed. The first parallel lint attempt encountered the shared golangci-lint lock; sequential retry passed.
- The commit's pre-commit hook independently passed lint and tests.
- TTC analysis reproduced 8 cells and 192 operations. CSV and normalized JSON matched; all PNGs were byte-identical. SVG metadata timestamps and generated clip-path IDs differed, so raw SVG byte comparison is not a valid semantic check.

### Review guidance
Inspect the deletion commit and the small `EchoOutput` to `IntakeOutput` relocation. Verify that immutable TTC source evidence and the audit script remain present.
