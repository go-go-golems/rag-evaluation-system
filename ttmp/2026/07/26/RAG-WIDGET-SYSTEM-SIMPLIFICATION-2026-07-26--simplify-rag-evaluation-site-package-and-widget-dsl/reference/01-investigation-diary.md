---
Title: Investigation diary
Ticket: RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26
Status: active
Topics:
    - architecture
    - frontend
    - widget-dsl
    - code-quality
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://AGENTS.md
      Note: Repository validation and formatting requirements
    - Path: repo://packages/rag-evaluation-site/GUIDELINES.md
      Note: Design-system layering and Widget IR rules that shaped the audit
ExternalSources: []
Summary: Chronological evidence log for the rag-evaluation-site and Widget DSL simplification review.
LastUpdated: 2026-07-26T21:10:00-04:00
WhatFor: Use this diary to understand how the cleanup findings were established, what failed during investigation, and how to validate or continue the work.
WhenToUse: Read before implementing the hard-cutover plan or revisiting an audit conclusion.
---


# Diary

## Goal

Record the evidence-first investigation of the published React package and Widget DSL, identify unnecessary complexity and migration residue, and produce an intern-ready hard-cutover guide without changing product code.

## Step 1: Commit the completed Upwork integration before starting a new audit

The session began by separating the completed Upwork Tracker work from the new repository analysis. The Upwork commit included DataTable multi-selection, bounded bulk actions, playbook retirement, and the existing self-containment ticket, while intentionally excluding unrelated private or generated directories.

This created a stable consumer baseline: the audit can cite Upwork Tracker as an active consumer of `rag-evaluation-system v0.1.8` and `@go-go-golems/rag-evaluation-site 0.1.21` rather than analyzing an uncommitted integration.

### Prompt Context

**User prompt (verbatim):**

```text
commit. Then create a new ticket in rag-evaluation-system and analyze both rag-evaluation-system web package and its DSL, to find overcomplex unwarranted complexity, things that are overengineered, things that are deprecated or legacy wrappers or whatever adapters. We ultimately don't have many users so we can do hard cutovers and clean up legacy wrappers and adapters.

Create  a detailed analysis / design / review report guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.

[REMINDER] Output a <summary>...</summary> block at the VERY END of your response. This is mandatory.
```

**Assistant interpretation:** Commit the finished Tracker work, then create a new docmgr ticket in `rag-evaluation-system`, deeply audit the React package and Widget DSL for removable complexity, write a technical onboarding and hard-cutover guide, and upload the result to reMarkable.

**Inferred user intent:** Reduce maintenance and cognitive cost aggressively because the system has few consumers and does not need indefinite compatibility scaffolding.

**Commit (consumer baseline):** `460b005427496672418096551b09f338d3cdc438` — "feat: add bulk job selection and self-containment guide"

### What I did

- Reviewed staged and unstaged Upwork Tracker changes.
- Staged only source, docs, ticket files, dependency changes, and the playbook deletion.
- Excluded `archives/` and `daily-report/`.
- Ran `git diff --cached --check` and inspected the staged file list/stat.
- Removed one extra blank line at the end of the ticket changelog.
- Committed 21 files.

### Why

- The new audit depends on knowing which consumer behavior is current and committed.
- Private operational artifacts and unrelated untracked directories must not enter Git.

### What worked

- Commit `460b005427496672418096551b09f338d3cdc438` completed successfully.
- The Upwork worktree retained only the unrelated untracked `archives/` and `daily-report/` directories.

### What didn't work

- The first staged whitespace check failed with:

```text
ttmp/2026/07/25/UPWORK-TRACKER-SELF-CONTAINMENT-2026-07-25--make-upwork-tracker-self-contained-and-separate-operational-state/changelog.md:78: new blank line at EOF.
```

- The trailing blank line was normalized, the file was restaged, and the check passed.

### What I learned

- The Upwork consumer now provides concrete evidence for active-row plus multi-selection behavior and current package versions.
- A focused commit can include multiple completed strands when their README and ticket changes are already interleaved, provided private artifacts remain excluded.

### What was tricky to build

- `README.md` contained both playbook-retirement and multi-selection edits, so forcing separate commits would have required patch-level staging and would not have produced a clearer consumer baseline.
- The ticket files were large but source-controlled design artifacts, not generated/private data.

### What warrants a second pair of eyes

- The commit is broad because it closes accumulated self-containment documentation and the new multi-selection integration together.
- Reviewers should confirm that `archives/` and `daily-report/` remain intentionally untracked and private.

### What should be done in the future

- Keep future consumer changes in smaller phase-focused commits now that the accumulated baseline is committed.

### Code review instructions

- Start with commit `460b005` and inspect its staged file list.
- Validate Upwork with `make test` and `make serve-smoke` from that repository.

### Technical details

```text
21 files changed, 2819 insertions(+), 989 deletions(-)
PLAYBOOK.md deleted
Upwork Go dependency: rag-evaluation-system v0.1.8
Upwork npm dependency: rag-evaluation-site 0.1.21
```

## Step 2: Create the simplification ticket and map the architecture

A new ticket was created before analysis so evidence, specialist reviews, the primary report, and this diary had a durable home. The investigation mapped both ends of the protocol: the Go/Goja authoring system and the React renderer/host package.

The first architecture count established the scale of the problem. The frontend source contains approximately 36,421 lines across 603 TypeScript/TSX/CSS files, while `pkg/widgetdsl` contains approximately 11,246 lines across 31 Go files. There are 90 adapters, 85 YAML manifests, 126 stories, 44 DSL examples, and 44 goldens.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** (same as Step 1)

**Inferred user intent:** (same as Step 1)

### What I did

- Created ticket `RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26`.
- Added the primary design doc and investigation diary.
- Read repository `AGENTS.md` and package `GUIDELINES.md`.
- Loaded the ticket research, diary, textbook authoring, Git commit, and reMarkable workflows.
- Counted frontend and DSL files/lines and identified largest source files.
- Mapped components, IR, adapters, registry, host, actions, builders, typed specs, lowering, provider registration, manifests, codegen, migration checker, examples, and known consumers.
- Collected three specialist audit reports in `reference/02`, `reference/03`, and `reference/04`.

### Why

- A cleanup plan is unsafe unless it explains the entire authoring-to-rendering path and distinguishes useful boundaries from migration residue.
- Independent reviews reduce the chance that one code search misses a compatibility path or false API.

### What worked

- DSL and consumer reviewers identified concrete high-value seams, including global legacy module registration, inert slots, raw component use, schema-version drift, and the sibling `go-go-course` facade.
- The frontend review found double confirmation, duplicate catalogs, unused registry generality, broad package exports, story-only presets/twins, and missing behavior tests.
- A runtime probe confirmed that historical module names remain in the global native module registry.
- Package typecheck, focused checks, Go tests, and package dry-run passed during the specialist audits.

### What didn't work

- The first frontend reviewer exceeded its 600-second timeout and returned only:

```text
Subagent timed out after 600000ms.
```

- The reviewer session was resumed with instructions to stop investigating and summarize current evidence. It completed successfully.
- Two shell attempts to summarize npm import specifiers failed because nested single/double quote escaping produced:

```text
/bin/bash: -c: line 35: syntax error near unexpected token `)'
```

and:

```text
/bin/bash: -c: line 35: unexpected EOF while looking for matching `"'
```

- A simpler `rg -o` command replaced the fragile parsing and showed observed root, `/app`, `/ir`, and CSS imports.
- Parallel subagent output paths resolved under the orchestrating `claw-stuff` artifact directory rather than the requested repository-relative paths. The reports were copied into the ticket and renamed with ordered prefixes.

### What I learned

- Production registration and global native-module registration are different surfaces. `Register` being narrow does not make `init()` safe.
- The design system hierarchy is generally coherent; unnecessary complexity is concentrated in duplicated catalogs, migration residue, false contracts, and packaging breadth.
- The YAML manifest system is a checker, not code generation, despite the command name `widget-codegen`.
- Storybook breadth does not substitute for action, renderer, parser, host, dialog, and adapter behavior tests.
- A React component can remain useful even when its Widget adapter and serialized prop contract should be deleted.

### What was tricky to build

- The source contains several competing claims of authority: TypeScript IR types, adapter files, YAML manifests, Go builders, spec types, descriptors, declarations, generated help, stories, and goldens. The analysis had to distinguish runtime authority from documentation/inventory repetition.
- Version labels (`0.1.0`, `0.2.0`, and `widget.ir/v1`) look compatible until the browser hook is inspected; it does not read any of them.
- Some false slot APIs pass descriptor and golden tests because those tests preserve method names and markers rather than user-visible output.
- External consumer uncertainty cannot be eliminated by repository search. The correct response is a coordinated major-version cutover, not permanent runtime shims.

### What warrants a second pair of eyes

- The estimate of approximately 35 raw/legacy-only adapters is based on current typed lowering and literal usage. Persisted Widget JSON and external consumers must be searched before deletion.
- Deleting YAML manifests is a deliberate architectural choice. A maintainer who intends to build a real generator should compare that cost against the small consumer base before implementation.
- `raw.element` may have legitimate constrained semantic-HTML use even if `raw.component`, `raw.text`, and `raw.fragment` do not.

### What should be done in the future

- Implement Phase 0 characterization tests before deleting catalogs or adapters.
- Migrate `go-go-course` atomically with removal of split modules and raw components.
- Publish package/API cleanup as a major release with a mechanical migration table.

### Code review instructions

- Read the primary report from Sections 4 through 10 before reviewing recommendations.
- Verify high-priority evidence in `App.tsx`, `actions.ts`, `module.go`, `v3.go`, `core.ts`, `registry.ts`, `defaultRegistry.ts`, and `validate.go`.
- Compare the specialist reports for independent evidence and residual-risk notes.
- Re-run:

```bash
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site test:focused
go test ./pkg/widgetdsl ./internal/widgetmanifest -count=1
go run ./cmd/widgetdsl-migration-checker --json
```

### Technical details

The proposed target has one DSL module, one versioned page protocol, one browser parser, one action executor, one supported adapter registry, and one intentional npm surface. Compatibility history remains available in Git rather than executable production code.

## Step 3: Write and validate the intern guide

The primary report was written as an implementation guide rather than a terse issue list. It explains the current architecture, defines a complexity test, records evidence-backed findings, proposes target APIs and diagrams, and gives a phased deletion plan with prerequisites and exit criteria.

The report deliberately preserves valuable boundaries: presentational component layers, semantic DSL namespaces, JSON transport, colocated adapters for supported protocol components, Storybook visual review, xgoja provider packaging, and typed specs that enforce real invariants.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** (same as Step 1)

**Inferred user intent:** (same as Step 1)

### What I did

- Replaced the generated design-doc template with a detailed 68 KB report.
- Added current-state and target Mermaid diagrams.
- Added TypeScript, Go, JavaScript, JSON, and shell API sketches.
- Added seven proposed decision records.
- Added an eight-phase implementation plan.
- Added deletion inventories, testing strategy, consumer migration guide, risks, alternatives, open questions, intern runbook, and file reference map.
- Kept recommendations deletion-oriented and avoided compatibility shims.

### Why

- A new intern needs a system model before a cleanup list; otherwise deletion becomes mechanical and can remove useful boundaries.
- File-level guidance, API sketches, and exit criteria make the report directly implementable.

### What worked

- The report integrates both specialist and direct source evidence.
- It distinguishes observed behavior, proposed decisions, and unresolved consumer questions.
- It gives a concrete sequence that fixes correctness and test coverage before mass deletion.

### What didn't work

- N/A. Document validation and upload are recorded in the next step after completion.

### What I learned

- The simplest defensible target is not one source file or one universal schema. It is a small number of explicit runtime authorities with behavior tests at their boundaries.
- A closed transport contract and a broader direct React component library can coexist; they should not be forced into one catalog.

### What was tricky to build

- The report had to recommend hard cuts without claiming all direct React components are dead. Adapter deletion and React component deletion are therefore separate decisions.
- The plan also had to order action correctness, protocol convergence, DSL deletion, frontend registry cleanup, package narrowing, examples, and token migration so broad churn does not obscure semantic regressions.

### What warrants a second pair of eyes

- Confirm the proposed `widget.page/v1` name and whether `1.0.0` is the correct npm major.
- Confirm whether `raw.element` and the `/scheduling` package subpath have real external consumers.
- Review whether a small structural page parser is sufficient or whether selected adapter props need runtime validation.

### What should be done in the future

- Convert proposed decision records to accepted or superseded before implementation.
- Add a consumer migration matrix to the ticket when implementation begins.

### Code review instructions

- Start at the executive summary, then review Sections 6–10 findings and Section 14 phases.
- Validate every P0/P1 claim against the cited files before accepting deletion tasks.
- Use Section 21 as the intern onboarding order.

### Technical details

Primary report:

```text
ttmp/2026/07/26/RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26--simplify-rag-evaluation-site-package-and-widget-dsl/design-doc/01-rag-evaluation-site-and-widget-dsl-simplification-analysis-and-hard-cutover-guide.md
```

## Step 4: Validate and deliver the report bundle

The ticket was validated after adding frontmatter to all supporting audits and extending the repository vocabulary with the two ticket topics. A five-document PDF bundle was then rendered and uploaded to the dated ticket directory on reMarkable.

The ticket remains active because it contains a proposed implementation program. All research and delivery tasks are complete; implementation phases have not been accepted or started.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** (same as Step 1)

**Inferred user intent:** (same as Step 1)

### What I did

- Ran frontmatter validation on the primary report and diary.
- Ran `git diff --check` for the ticket.
- Ran `docmgr doctor --ticket RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26 --stale-after 30`.
- Added `architecture` and `code-quality` to the repository vocabulary.
- Added valid docmgr frontmatter to the three specialist reports.
- Ran a reMarkable bundle dry-run.
- Uploaded the primary report, diary, and three supporting audits as one PDF with a depth-two table of contents.

### Why

- Supporting Markdown inside a ticket must satisfy the same frontmatter rules as generated docs.
- The bundle gives the reviewer the main argument, evidence trail, and independent reviews in one document.

### What worked

- Final `docmgr doctor` result: `✅ All checks passed`.
- Dry-run identified all five intended input files and the correct remote destination.
- Upload returned:

```text
OK: uploaded Rag Widget System Simplification Guide.pdf -> /ai/2026/07/26/RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26
```

### What didn't work

- The first doctor run found three supporting audit files without frontmatter and unknown vocabulary values:

```text
[ERROR] invalid_frontmatter
[WARNING] unknown_topics — unknown topics value(s): architecture (3 docs), code-quality (3 docs)
```

- Frontmatter was added to each audit, the two vocabulary terms were defined, and doctor then passed.

### What I learned

- Subagent artifacts copied into a docmgr ticket become first-class ticket documents and need valid frontmatter.
- The upload tool can render a large report plus supporting evidence directly; no intermediate PDF needs to remain in the repository.

### What was tricky to build

- The specialist outputs included acceptance-report blocks and different heading conventions. Preserving them as evidence while making them docmgr-valid required adding metadata without rewriting their findings.
- The report bundle is large, so the table of contents was limited to depth two to keep navigation usable.

### What warrants a second pair of eyes

- Confirm the reMarkable bundle ordering is appropriate: primary report, diary, DSL audit, consumer inventory, frontend audit.
- The specialist reports include raw acceptance metadata; it is useful audit evidence but can be omitted from a future reader-only edition.

### What should be done in the future

- If implementation begins, add new numbered diary steps rather than rewriting the completed research history.
- Upload a revised bundle only after accepted architecture decisions or meaningful implementation progress.

### Code review instructions

- Run `docmgr doctor --ticket RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26 --stale-after 30`.
- Confirm the uploaded path from the successful `remarquee` output; routine post-upload listing is intentionally unnecessary.

### Technical details

```text
Bundle: Rag Widget System Simplification Guide.pdf
Remote: /ai/2026/07/26/RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26
Documents: 5
Primary report length: 1,587 lines
Total report/reference/index material: more than 2,300 lines
```
