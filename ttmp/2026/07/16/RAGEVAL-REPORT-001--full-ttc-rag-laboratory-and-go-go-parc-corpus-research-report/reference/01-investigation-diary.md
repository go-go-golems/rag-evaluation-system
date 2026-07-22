---
Title: Investigation Diary
Ticket: RAGEVAL-REPORT-001
Status: active
Topics:
    - rag-eval
    - ttc
    - corpus
    - embeddings
    - search
    - evaluation
    - obsidian
    - reranking
    - workflow
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/.codex/skills/diary/SKILL.md
      Note: Strict diary format and required sections used for this retroactive rewrite
    - Path: abs:///home/manuel/.codex/skills/diary/references/diary.md
      Note: Diary rationale and worked example consulted before restructuring
ExternalSources: []
Summary: Chronological report-production diary with commands, failures, decisions, and review instructions.
LastUpdated: 2026-07-16T20:30:00-04:00
WhatFor: Reconstruct how the full TTC RAG and go-go-parc research report was produced.
WhenToUse: Use this diary when reviewing the report or continuing the corpus and laboratory work.
---


# Diary

## Goal

This diary captures the production of the full TTC RAG laboratory and go-go-parc corpus report. It records the investigation inputs, report drafting, validation, reMarkable upload, vault copy, Git delivery, and the later correction to this diary's format. It is intentionally retroactive: the report was completed before the diary was normalized to the strict diary-skill structure.

## Step 1: Create the report ticket and define the evidence boundary

The first step established a dedicated docmgr ticket so that the report, source inventory, diary, tasks, and changelog would be reviewable as one unit. The scope covered the TTC RAG laboratory, immutable experiment runs, evaluation-card authoring, retrieval findings, reranking research, fluent JavaScript API work, the Ollama tunnel procedure, and the parallel investigation of the go-go-parc Obsidian corpus.

The report was deliberately scoped as a research artifact rather than a code change. The evidence boundary is the current `rag-evaluation-system` worktree plus copied local source documents from earlier tickets. Generated databases and traces remain local artifacts; committed scripts and manifests explain how to reproduce them.

### Prompt Context

**User prompt (verbatim):** "Ok,  write a detailed project report for the obsidian vault as a deep dive technical analysis blog post using a textbook writing style (no analogies, see skill).      Commit and push the bsidian vault when done (go-go-parc vault).    about all the work we have been doing (consult older diaries, etc...), our findings, how it's built up, which tools and scripts we use, pseudocode / code snippets. It can be really really really long, this is a huge and important project."

**Assistant interpretation:** Create a very long, textbook-style technical report from the accumulated TTC, vault, reranker, experiment, and diary work; publish the complete article in go-go-parc and commit and push it.

**Inferred user intent:** Preserve the project's accumulated reasoning and measurements in a durable, searchable Obsidian document that a colleague or intern can use without reconstructing the entire chat history.

### What I did

- Created ticket `RAGEVAL-REPORT-001` under `ttmp/2026/07/16/`.
- Added the report design document and a diary reference document.
- Added tasks for inventory, source collection, report writing, vault copy, and delivery.
- Identified the TTC ticket `RAGEVAL-TTC-LAB-001` as the primary implementation source.

### Why

- A separate ticket prevents the report from being mixed with the still-active TTC implementation tasks.
- A local source inventory makes claims auditable and makes the reMarkable bundle self-contained.

### What worked

- `docmgr status --ticket RAGEVAL-REPORT-001` showed an active ticket with two documents.
- The ticket structure was accepted by docmgr without stale-document findings.

### What didn't work

- The first inspection assumed the ticket was rooted at the workspace-level `ttmp`; the actual ticket is under `rag-evaluation-system/ttmp`. The corrected absolute path is `/home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/16/RAGEVAL-REPORT-001--full-ttc-rag-laboratory-and-go-go-parc-corpus-research-report`.

### What I learned

- In this workspace, `docmgr --root ttmp` must be run from the `rag-evaluation-system` worktree, not from its parent workspace.
- The report ticket can safely refer to earlier ticket documents by copying them into `sources/` while keeping the original tickets authoritative.

### What was tricky to build

- The worktree contains several repositories and multiple `ttmp` directories. Using relative paths from the parent workspace produced a false “file not found” diagnosis. I resolved this by checking `pwd`, locating the ticket with `find`, and then using absolute paths for report operations.

### What warrants a second pair of eyes

- Confirm that the source inventory remains a faithful copy of the cited ticket documents and that no generated or private vault files were included.

### What should be done in the future

- Keep report tickets inside the repository whose code and docmgr root they document.
- Add a ticket README note with the absolute repository root when a workspace contains multiple worktrees.

### Code review instructions

- Start at `ttmp/2026/07/16/RAGEVAL-REPORT-001--full-ttc-rag-laboratory-and-go-go-parc-corpus-research-report/README.md` and `tasks.md`.
- Validate with `docmgr --root ttmp status --ticket RAGEVAL-REPORT-001` from `rag-evaluation-system`.

### Technical details

```sh
cd /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system
docmgr --root ttmp status --ticket RAGEVAL-REPORT-001
```

## Step 2: Inventory prior TTC, reranker, vault, and experiment artifacts

The second step read the accumulated design documents, implementation diaries, evaluation-card packets, experiment results, API sketches, and operational playbooks. The goal was to distinguish implemented behavior from planned behavior and to carry forward measured values rather than merely repeating architectural aspirations.

The inventory established the report's core evidence: the TTC snapshot and artifact hashes, the 148-card candidate dataset, the development retrieval metrics, the 240-card split plan, the BGE reranker status, the fluent JavaScript builder target, the Ollama tunnel procedure, and the vault inclusion boundary.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Read older diaries and documents before writing, so the article represents the complete project rather than only the latest code state.

**Inferred user intent:** Avoid losing the history of decisions, failed experiments, and operational knowledge that gives the measurements their meaning.

### What I did

- Read the TTC baseline design, adjudication protocol, v2 corpus design, 240-card split protocol, source validation, candidate registration, and development-results documents.
- Read the TTC implementation diary and the reranker design/results and diary.
- Read the vault design, inclusion manifest, and vault diary.
- Read the RAG DSL design/API and the Ollama tunnel playbook.
- Copied 22 primary local documents into the report ticket's `sources/` directory.
- Related the central implementation files with `docmgr doc relate`.

### Why

- The report needs citations to local evidence that a future reviewer can inspect without searching the full workspace.
- The artifact hashes and measured metrics must be copied exactly from the result documents.

### What worked

- The source inventory includes 22 files covering architecture, data, evaluation, reranking, vault ingestion, DSL design, operations, and diaries.
- The report can cite both the original TTC ticket and a stable local source copy.

### What didn't work

- Several initial copy attempts used shorthand paths that did not exist. I corrected them by listing the ticket tree and copying from the exact `reference/` and `design-doc/` paths.
- A prior development trace invocation produced zero traces and another produced only 70 because the parser and input paths did not match the authored card formats. Those failures are retained as report evidence rather than omitted.

### What I learned

- Source validation and card registration are different operations: a card can be syntactically valid and still refer to a source absent from the frozen snapshot.
- The development result is provisional because model-authored cards still require human adjudication.

### What was tricky to build

- It was necessary to avoid presenting planned parent collapse, citation hydration, summary representations, and reranking as already complete. I separated “implemented,” “measured,” and “next” states in the report and labeled the development result provisional.

### What warrants a second pair of eyes

- Verify every metric in the report against `09-ttc-expansion-development-run-results.md`.
- Review whether the vault inclusion boundary excludes generated artifacts, duplicate exports, and private material as intended.

### What should be done in the future

- Add automated source-count and metric cross-checks to report generation so copied numbers cannot drift.
- Keep a manifest of every report source with its content hash.

### Code review instructions

- Review `sources/` first, then compare claims in the report with `sources/ttc-development-results.md`, `sources/ttc-240-partition.md`, and `sources/vault-inclusion-manifest.md`.
- Run `find ttmp/2026/07/16/RAGEVAL-REPORT-001--*/sources -type f | sort`.

### Technical details

The copied source set contains 22 Markdown/Go files. Current artifact identities recorded in the report are:

```text
snapshot  sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409
chunks    sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392
bm25      sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691
vectors   sha256:2665c5249b8352ce6904fc00c934534dd179f3eeef0a6a75429a9034be0e03e0
dataset   candidate:ttc-expansion-v0
```

## Step 3: Draft the textbook technical report

The third step replaced the report stub with a complete technical article. It explains the five-layer laboratory architecture, normalized source and representation models, chunking strategies, BM25/vector/RRF retrieval, parent collapse, citation hydration, immutable experiment identity, evaluation-card construction, reranking, the web UI, JavaScript API targets, cost accounting, failure records, and a staged intern checklist.

The article includes Mermaid diagrams, pseudocode, Go and TypeScript contracts, shell commands, measured tables, and explicit references. It names GPT-5.6 (Sol) as the authoring agent because that was requested context for the report. The report is intentionally careful about status: vector and hybrid measurements are current development findings, while human adjudication, full representation comparisons, parent collapse, citation hydration, and BGE reranking remain follow-up work.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Write the article in a textbook style with enough technical depth for an intern to implement the remaining work.

**Inferred user intent:** Make the project's architecture and reasoning teachable, not merely summarize the latest benchmark.

### What I did

- Replaced the report placeholders with 675 lines of prose, tables, diagrams, pseudocode, API references, code snippets, findings, alternatives, failure records, metric definitions, and review guidance.
- Included the TTC counts, content hashes, candidate-card count, 240-card plan, development metrics, latency, and cost observations.
- Added the JavaScript fluent-builder target while marking unimplemented methods as an API target.
- Added a direct explanation of the filter-leak risk and the required trace provenance.
- Added a go-go-parc corpus section covering project reports, design docs, diaries, playbooks, code, and inclusion rules.

### Why

- A new intern needs the data model and invariants before editing retrieval code.
- Explicit failure records prevent future runs from repeating silent zero-card or partial-card experiments.

### What worked

- The final report is 44,333 bytes and passes Markdown/docmgr validation.
- The report includes 15+ local references through the copied source inventory and explicit reference list.

### What didn't work

- The first large `apply_patch` failed with `apply_patch verification failed: invalid hunk at line 319, Expected update hunk to start with a @@ context marker, got: 'type RepresentationSpec struct {'`. The cause was an unprefixed line inside a fenced code block. I split the report into smaller patches and added every code line with the required patch prefix.

### What I learned

- Large documentation patches are safer when divided by conceptual sections; code fences make malformed patch hunks easy to miss.
- A report should state the denominator and provisional status beside every measured table.

### What was tricky to build

- The report needed to be long without turning planned work into false claims. I used separate sections for current findings, implementation plan, alternatives, open questions, and appendices, and kept “API target” labels on future JavaScript methods.

### What warrants a second pair of eyes

- Review the RRF formula, metric definitions, and the interpretation that hybrid improves relevant recall@10 but not every top-rank metric.
- Review the proposed Go and TypeScript interfaces before implementation; they are design contracts, not yet exported API guarantees.

### What should be done in the future

- Add citations from the report to exact line ranges in generated result artifacts where possible.
- Generate the report's metric table directly from the immutable run JSON after the evaluation schema stabilizes.

### Code review instructions

- Start with `design-doc/01-full-ttc-rag-laboratory-and-go-go-parc-corpus-research-report.md`.
- Validate with `wc -c -l` and `docmgr --root ttmp doctor --ticket RAGEVAL-REPORT-001`.

### Technical details

The development table included in the report is:

```text
BM25       R@1 .7847  R@3 .8403  R@10 .8889  MRR .8221  relevant-R@10 .7442
Vector     R@1 .8750  R@3 .9653  R@10 .9722  MRR .9174  relevant-R@10 .8588
Hybrid RRF R@1 .8542  R@3 .9375  R@10 .9722  MRR .9005  relevant-R@10 .8947
latency mean 173 ms, p50 175 ms, p95 230 ms, min 72 ms, max 465 ms
```

## Step 4: Copy, validate, and upload the report

The fourth step delivered the report as an Obsidian note and a reMarkable bundle. The vault destination is dated under `Research/2026/07/16/`, and the copied file was compared by SHA-256 against the ticket report before the vault commit. The reMarkable bundle includes the report, diary, source inventory, scripts, and ticket metadata.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Publish the complete report to the specified Obsidian vault and make it available on the reMarkable tablet.

**Inferred user intent:** Make the research accessible both in the searchable vault and for tablet review without manually reconstructing the document.

### What I did

- Checked that the destination article did not already exist.
- Copied the report to `Research/2026/07/16/ARTICLE - Full TTC RAG Laboratory and go-go-parc Corpus Research Report.md`.
- Compared ticket and vault SHA-256 values.
- Ran `docmgr --root ttmp doctor --ticket RAGEVAL-REPORT-001`.
- Uploaded the ticket bundle with `remarquee upload bundle` to `/ai/2026/07/16/RAGEVAL-REPORT-001`.
- Checked the report ticket tasks for source inventory, report drafting, vault copy, and delivery.

### Why

- Dated Research paths preserve the vault's chronological organization.
- A checksum comparison proves that the vault article is the complete ticket report rather than a shortened summary.

### What worked

- Ticket and vault SHA-256 matched: `abcb73617b44f75affbd2c17f55184356ebdbcd389d1bbf8a5c001a4c0398bf2`.
- `docmgr doctor` reported `All checks passed`.
- ReMarkable returned: `OK: uploaded TTC-RAG-LAB-FULL-REPORT.pdf -> /ai/2026/07/16/RAGEVAL-REPORT-001`.

### What didn't work

- A checksum command initially used a path relative to the wrong working directory and returned `sha256sum: ttmp/2026/07/16/...: No such file or directory`. I reran it with both absolute paths and obtained the matching hashes.

### What I learned

- External vault operations should use absolute paths and an existence check before copying.
- The vault contained unrelated untracked files (`.pi/npm/`, `.pi/settings.json`, `.ttmp.yaml`); staging only the new article avoids claiming ownership of those changes.

### What was tricky to build

- The report ticket and vault are separate Git repositories with different writable boundaries. The copy, checksum, staging, commit, and push operations had to be performed with explicit external authorization while preserving unrelated vault state.

### What warrants a second pair of eyes

- Confirm the reMarkable bundle renders all source documents and that the long Markdown note opens correctly in Obsidian.
- Confirm the vault commit contains exactly one new file.

### What should be done in the future

- Keep using a checksum comparison whenever a complete ticket document is mirrored into the vault.
- Add a small delivery script only if repeated report publication justifies it; do not introduce a script for this one-off operation.

### Code review instructions

- Check the vault note path and the ticket report path.
- Run `git -C /home/manuel/code/wesen/go-go-golems/go-go-parc show --stat --oneline 8697995`.
- Verify the reMarkable remote directory `/ai/2026/07/16/RAGEVAL-REPORT-001`.

### Technical details

```sh
remarquee upload bundle \
  ttmp/2026/07/16/RAGEVAL-REPORT-001--full-ttc-rag-laboratory-and-go-go-parc-corpus-research-report \
  --remote-dir /ai/2026/07/16/RAGEVAL-REPORT-001 \
  --name TTC-RAG-LAB-FULL-REPORT --force --non-interactive
```

## Step 5: Commit and push the vault article

The fifth step committed only the new article to go-go-parc and pushed `main` to origin. The report ticket itself was also committed in the rag-evaluation-system worktree so its source inventory, diary, and metadata are durable. The vault still contains unrelated untracked files, which were intentionally left untouched.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Complete the Git delivery and leave both the report ticket and the vault article in a reviewable state.

**Inferred user intent:** Ensure the report is not merely present locally but available to collaborators through the repository remote.

### What I did

- Committed the report ticket as `13f55dc` with message `docs(rag-lab): add full TTC and vault research report`.
- Staged only the new vault article.
- Committed the vault article as `8697995` with message `Research report: TTC RAG laboratory and vault corpus`.
- Pushed `main` to `ssh://git@github.com/go-go-golems/go-go-parc`.
- Confirmed all five report-ticket tasks are checked and `docmgr doctor` passes.

### Why

- Separate commits keep the implementation repository's research package independent from the Obsidian vault mirror.
- Narrow staging protects unrelated user changes in the vault.

### What worked

- Vault push succeeded: `839eed4..8697995  main -> main`.
- The vault commit contains one file and 675 insertions.
- The report ticket status is clean, with all tasks complete and no doctor findings.

### What didn't work

- No commit or push failure occurred. The only caveat is that the vault working tree still reports unrelated untracked files; they were not part of this task.

### What I learned

- A complete research delivery has three independently verifiable states: ticket docs committed, vault note checksum-matched, and vault commit pushed.
- Marking docmgr tasks complete after the delivery checks makes the ticket's status immediately legible to the next contributor.

### What was tricky to build

- It was important not to run `git add .` in the vault. The working tree contained unrelated files, so the article was staged by exact path and the resulting cached diff was inspected before commit.

### What warrants a second pair of eyes

- Confirm origin's `main` contains commit `8697995` and that no vault automation reformats the Markdown after push.
- Review the report for any claim that should be downgraded from “implemented” to “planned” as the codebase evolves.

### What should be done in the future

- Close `RAGEVAL-REPORT-001` only after a human reviews the report and confirms the reMarkable rendering.
- Keep the report immutable as a historical snapshot; add dated follow-up notes for new measurements rather than rewriting findings silently.

### Code review instructions

- In the rag worktree: `git show --stat 13f55dc` and `docmgr --root ttmp doctor --ticket RAGEVAL-REPORT-001`.
- In the vault: `git -C /home/manuel/code/wesen/go-go-golems/go-go-parc show --stat 8697995` and `git -C /home/manuel/code/wesen/go-go-golems/go-go-parc status --short`.
- Reopen the uploaded bundle at `/ai/2026/07/16/RAGEVAL-REPORT-001` on reMarkable.

### Technical details

```text
rag-evaluation-system commit 13f55dc
go-go-parc commit          8697995
vault remote                ssh://git@github.com/go-go-golems/go-go-parc
vault branch                main
report ticket tasks         5/5 complete
docmgr doctor               All checks passed
```

## Step 6: Normalize this diary to the strict diary-skill format

The final step was prompted by a review request after delivery. The previous diary contained useful facts but used `Context`, `Quick Reference`, and free-form dated entries rather than the strict numbered-step structure. This update preserves the facts and rewrites them into numbered steps with prose introductions, prompt context, exact failures, tricky-build notes, second-pair-of-eyes items, future work, review commands, and technical details.

### Prompt Context

**User prompt (verbatim):** "Read the diary format in the diary skill, update the document retroactively to have the proper format."

**Assistant interpretation:** Read the authoritative diary skill instructions and rewrite the existing report diary so every historical step follows the required structure.

**Inferred user intent:** Make the diary useful for code review and future continuation, not merely a chronological list of outcomes.

### What I did

- Read `/home/manuel/.codex/skills/diary/SKILL.md` completely.
- Read `/home/manuel/.codex/skills/diary/references/diary.md` for the detailed format and worked example.
- Replaced the prior diary with six numbered steps.
- Added the original report prompt verbatim in Step 1 and the diary-format correction prompt verbatim in this step.
- Added exact commands, hashes, commit IDs, failures, lessons, tricky-build explanations, review instructions, and future work to each step.

### Why

- The diary skill requires a strict per-step structure and explicit prompt context.
- Retroactive formatting is appropriate because the implementation and delivery already happened; the missing structure is documentation debt, not a reason to rewrite the report history.

### What worked

- The rewritten file begins with `# Diary`, has sequential `Step 1` through `Step 6` headings, and includes all required subsections.
- Historical failures remain visible, including the malformed large patch, parser/card-count issues, path mistakes, and the sandboxed IPv6 test limitation.

### What didn't work

- The prior diary did not satisfy the strict format because it used free-form dated entries. No code behavior was affected; the correction is documentation-only.

### What I learned

- A diary is a review artifact: each step must explain not only what happened, but why it happened, what was difficult, what remains risky, and how to validate it.
- The exact user prompt should be recorded once and referenced in later steps to avoid noisy repetition.

### What was tricky to build

- Reconstructing the chronology required separating report production from prior TTC implementation work. I kept the TTC implementation diary as the source of truth for code changes and used this diary for the report's inventory, authoring, delivery, and documentation correction.

### What warrants a second pair of eyes

- Check that every step's command examples point to existing files and that no historical failure was accidentally converted into a success claim.
- Confirm that the diary's `LastUpdated` metadata is changed when future entries are added.

### What should be done in the future

- Use this exact step template for any follow-up report or implementation diary.
- Add a new numbered step rather than rewriting completed steps when a reviewer finds a new issue.

### Code review instructions

- Start at this file and scan each `## Step N` for all required headings.
- Run `docmgr --root ttmp doctor --ticket RAGEVAL-REPORT-001`.
- Inspect the staged diff before committing the diary-only update.

### Technical details

The required per-step headings are:

```text
Prompt Context
What I did
Why
What worked
What didn't work
What I learned
What was tricky to build
What warrants a second pair of eyes
What should be done in the future
Code review instructions
Technical details
```

