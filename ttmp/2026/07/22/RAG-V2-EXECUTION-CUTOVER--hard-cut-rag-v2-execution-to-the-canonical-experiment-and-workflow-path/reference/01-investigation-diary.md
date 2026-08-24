---
Title: Investigation diary
Ticket: RAG-V2-EXECUTION-CUTOVER
Status: active
Topics:
    - rag
    - rag-eval
    - workflow
    - scripting
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://cmd/rag-eval/cmds/intake/root.go
      Note: Explicit separation of product intake from scientific study execution (commit 8925004)
    - Path: repo://cmd/rag-eval/cmds/study/command.go
      Note: Compile-only RAG study command and provider authority selection (commit 8925004)
    - Path: repo://cmd/rag-eval/main.go
      Note: Canonical command registration after direct runner and preview removal (commit 8925004)
    - Path: repo://pkg/ragworkflow/study_export.go
      Note: Immutable Workflow case bundle and Researchctl plan compiler (commit 8925004)
    - Path: repo://pkg/ragworkflow/study_export_test.go
      Note: Provider-free/provider-backed identity and boundary proofs (commit 8925004)
    - Path: repo://pkg/researchctladapter/adapter.go
      Note: Retained immutable domain input resolver after lifecycle deletion (commit 8925004)
    - Path: repo://ttmp/2026/07/22/RAG-V2-EXECUTION-CUTOVER--hard-cut-rag-v2-execution-to-the-canonical-experiment-and-workflow-path/analysis/03-cutover-acceptance-audit.md
      Note: Final objective checklist and validation evidence
    - Path: repo://ttmp/2026/07/22/RAG-V2-EXECUTION-CUTOVER--hard-cut-rag-v2-execution-to-the-canonical-experiment-and-workflow-path/scripts/01-smoke-compiled-study.sh
      Note: Built-binary compiler to Researchctl to Workflow V3 acceptance
    - Path: repo://ttmp/2026/07/22/RAG-V2-EXECUTION-CUTOVER--hard-cut-rag-v2-execution-to-the-canonical-experiment-and-workflow-path/scripts/02-cutover-guards.sh
      Note: Deletion, import, help, docs, and manifest guards
    - Path: repo://web/src/components/pages/EvaluationPage/EvaluationPage.tsx
      Note: Frontend instructions switched to compile plus run-plan (commit 8925004)
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-22T19:25:22.618456064-04:00
WhatFor: ""
WhenToUse: ""
---



# Investigation diary

## Goal

Record the evidence and decisions used to create the intern-facing design for **RAG-V2-EXECUTION-CUTOVER** and place it in the wider experiment-platform convergence program.

## Step 1: Establish the ticket and implementation guide

This step created the ticket as a focused child of the multi-repository convergence program. The guide explains the existing systems, proposed contracts, pseudocode, diagrams, implementation phases, tests, deletion rules, and neighboring tickets so an intern can begin without reconstructing the architecture from Git history.

No production code was changed. The architecture is proposed work; current files are referenced as evidence and implementation entry points.

### Prompt Context

**User prompt (verbatim):** "Create all the tickets with docmgr, then for each ticket, Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.


That way we will have everything. Have all the tickets reference each other so that if someone stumbles upon one ticket, they'll know what this is all about."

**Assistant interpretation:** Create the umbrella and all directed tickets, write self-contained technical guides, cross-link them, validate docmgr state, and publish the complete set to reMarkable.

**Inferred user intent:** Preserve a coherent program plan before more ad-hoc implementation occurs, while making every workstream independently discoverable and approachable by a new engineer.

### What I did
- Created `RAG-V2-EXECUTION-CUTOVER` with docmgr in the rag-eval repository.
- Added a long-term design/implementation guide and this diary.
- Mapped the ticket to all nine sibling/program tickets.
- Grounded the guide in current repository packages and public contracts.

### Why
- The convergence crosses repository boundaries and needs explicit ownership and sequencing.
- Standalone guides reduce the chance that an implementer recreates generic functionality in a workload package.

### What worked
- Existing greenfield contracts provided concrete API and file anchors.
- The ticket scope could be expressed as one independently testable capability: Hard-cut to the sole canonical RAG execution path.

### What didn't work
- N/A. No implementation or runtime validation was attempted in this documentation-only step.

### What I learned
- The repositories are closer to the target internally than the prior TTC path suggested; the largest deficit is the canonical integration route.

### What was tricky to build
- The guides must be self-contained without duplicating ownership. This was handled by giving each ticket one primary capability, linking the complete program map, and making predecessor/successor relationships explicit.

### What warrants a second pair of eyes
- Public schema names, package deletion lists, and the exact boundary between RAG factor expansion and Researchctl scheduling should be reviewed before implementation.

### What should be done in the future
- Review and accept the design, then implement phases in dependency order with the ticket's acceptance tests.

### Code review instructions
- Begin with the ticket's `design-doc/` guide and its “Program context” section.
- Validate documentation with `docmgr doctor --ticket RAG-V2-EXECUTION-CUTOVER --stale-after 30`.

### Technical details
- Program umbrella: `EXPERIMENT-PLATFORM-CONVERGENCE`.
- Ticket responsibility: Hard-cut to the sole canonical RAG execution path.

## Step 2: Rewrite the briefs as textbook-style intern guides

The first drafts were structurally correct but too compressed. They named responsibilities and implementation phases without spending enough time teaching the foundational concepts, tracing real code, or showing how an intern would validate each claim. This revision expanded every guide into a teaching document rather than leaving it as an architecture summary.

The revised guides now introduce the domain vocabulary before APIs, walk through current source packages, show concrete identity and failure traces, provide pseudocode and data structures, explain common mistakes, and end with onboarding exercises. The ticket-research workflow remains responsible for evidence, bookkeeping, validation, and reMarkable publication; the textbook-authoring style governs the prose.

### Prompt Context

**User prompt (verbatim):** "these guides are much too concise / terse. Read the skill on writing guides for interns, and update them."

**User prompt (verbatim):** "write in a textbook authoring style as well"

**User prompt (verbatim):** "it's ticket-research-docmgr-remarkable skill"

**Assistant interpretation:** Replace concise RFC-style briefs with long-form, foundational, intern-oriented teaching guides while continuing to follow the ticket research, docmgr, diary, validation, and publication workflow.

**Inferred user intent:** Ensure the ticket set can genuinely onboard implementers without relying on the prior conversation or requiring them to rediscover the systems.

### What I did
- Loaded the full-blown intern research/design and textbook-authoring guidance, while retaining the pinned ticket-research workflow as the delivery contract.
- Expanded conceptual foundations, code tours, worked traces, API references, implementation routes, debugging guidance, common mistakes, and onboarding checklists.
- Preserved program navigation in every ticket.

### Why
- A list of recommendations does not teach the invariants that an implementer must preserve.
- The program spans three repositories and several identity/lifecycle layers; foundational explanation is required before implementation advice.

### What worked
- Existing canonical contracts and integration tests provide concrete teaching anchors.
- Textbook-style worked flows make responsibility boundaries easier to evaluate than responsibility tables alone.

### What didn't work
- The initial drafts were approximately architecture-brief length and did not meet the requested exhaustive intern-guide standard.

### What I learned
- Each guide needs both a local implementation path and the complete cross-program context; either one alone is insufficient.

### What was tricky to build
- Expansion had to add depth without duplicating ownership or introducing analogies. The solution was to teach using exact identities, state transitions, source files, schemas, and failure windows.

### What warrants a second pair of eyes
- Review whether each guide provides enough real source-level orientation for an intern and whether proposed package names should be accepted before implementation.

### What should be done in the future
- Review the expanded guides, then implement in dependency order. Add experiments to each ticket as its implementation starts.

### Code review instructions
- Read the technology-primer and guided-source-tour sections before reviewing proposed APIs.
- Check that all bullet points are complete technical statements and diagrams preserve actual system boundaries.

### Technical details
- Style: foundational prose, concrete examples, no analogies, precise diagrams, worked traces, and executable validation routes.

## Step 3: Validate and publish the program guide bundle

The ticket passed docmgr validation and its guide was included in the ordered ten-document program bundle. The dry run confirmed document ordering and destination; the real upload completed successfully.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Complete ticket bookkeeping and deliver all guides as one navigable PDF.

**Inferred user intent:** Make the complete program available as a durable reading package rather than scattered repository files.

### What I did
- Ran `docmgr doctor --ticket <ticket> --stale-after 30`.
- Performed a reMarkable bundle dry run with all ten guides in dependency order.
- Uploaded `Scriptable Experiment Platform Intern Guides` to `/ai/2026/07/22/EXPERIMENT-PLATFORM-CONVERGENCE`.

### Why
- A single bundle provides a table of contents across repository boundaries.

### What worked
- All ticket doctor checks passed and the upload reported success.

### What didn't work
- N/A.

### What I learned
- Cross-repository tickets remain locally owned while one ordered publication can present the program coherently.

### What was tricky to build
- Document order needed to follow dependency order rather than repository or alphabetical order.

### What warrants a second pair of eyes
- Review the PDF at normal reMarkable size for code-block and diagram readability.

### What should be done in the future
- Re-upload with `--force` only after accepted content revisions, because overwrite removes annotations.

### Code review instructions
- Verify ticket frontmatter, program links, design guide, tasks, and changelog.

### Technical details
- Remote destination: `/ai/2026/07/22/EXPERIMENT-PLATFORM-CONVERGENCE`.

## Step 4: Re-audit the live execution graph and establish the deletion manifest

Implementation resumed after the lowering, generic Scraper runner, canonical observations, and durable Geppetto operation tickets had completed. I re-read this ticket and the earlier legacy-cleanup disposition before deleting anything, then traced commands, imports, documentation, and package ownership against the current tree rather than relying on the original July 22 inventory.

The resulting machine-readable manifest distinguishes the direct RAG study executor that must now disappear from the older document-intake product path whose API/frontend replacement remains separately gated. This avoids both unsafe product deletion and the opposite error of keeping `rag-worker` merely because unrelated intake code still uses Scraper's old engine.

### Prompt Context

**User prompt (verbatim):** "RAG-V2-EXECUTION-CUTOVER and keep a detailed diary as you work ."

**Assistant interpretation:** Complete the accepted hard cut to Researchctl plans and Scraper Workflow V3, deleting superseded direct RAG execution while recording detailed evidence throughout implementation.

**Inferred user intent:** Leave one unambiguous production path for RAG v2 studies, preserve independently gated product behavior, and make the destructive cutover reviewable and reproducible.

### What I did
- Read the complete ticket guide, existing diary, tasks, diary skill, and `RAG-EVAL-LEGACY-CLEANUP` disposition.
- Enumerated command binaries, `internal/workflow*` packages, `ragengine` imports, Researchctl adapter imports, active docs, examples, and direct Scraper legacy imports.
- Confirmed that predecessor work has already removed `cmd/rag-ttc-v3-sweep` and `internal/workflowv3ttc` and has supplied deterministic/provider Workflow parity and generic runner acceptance.
- Created `analysis/01-execution-path-deletion-manifest.json` with an explicit classification, action, replacement, and required proof for every execution-path cluster.

### Why
- The guide requires a machine-readable deletion manifest and forbids unresolved classifications.
- The current tree has two superficially similar but semantically distinct old paths: direct RAG v2 study execution, which is superseded, and document intake/API processing, whose product migration is still independently gated.

### What worked
- Repository searches localized direct study execution to `study run`, `preview`, `cmd/rag-worker`, `internal/preparationworkflow`, and the run/progress portions of `pkg/researchctladapter`.
- Existing `pkg/ragworkflow` tests and generated Researchctl fixtures provide the required replacement evidence for provider-free and provider-backed execution.

### What didn't work
- An attempted read of `cmd/rag-eval/cmds/study/runner.go` failed with `ENOENT`; the execution loop is in `command.go`, so inspection continued there.
- The initial broad search mixed historical documentation, domain-level `replicates` fields, and live runner calls. I separated active Go imports, active operator docs, and historical ticket evidence before classifying paths.

### What I learned
- `ragengine` is still imported by Workflow tasks and product query runtime as a domain sequencing/semantic component, but only `cmd/rag-worker` exposes it as a competing production scheduler.
- The old `internal/workflow` package is not part of RAG v2 study execution; it drives document intake and APIs. Its deletion still requires the accepted product/API replacement gate and should not be smuggled into this cutover.
- Current `study compile` still emits direct-worker Researchctl specifications, so removing `study run` alone would not complete the cutover; compile output must target Workflow runner cases and Researchctl planning.

### What was tricky to build
- Directory names alone are misleading: both `internal/workflow` and `pkg/ragworkflow` contain workflow behavior, but only the latter owns canonical RAG v2 lowering. I classified by command trace, imported engine generation, custody owner, and product behavior rather than names.

### What warrants a second pair of eyes
- Review the `canonical-product-deferred` entries to confirm that document intake is correctly outside the RAG v2 study cutover and remains governed by the explicit API/frontend migration gate.
- Review whether `pkg/ragengine` production-import guards should allow only `pkg/ragworkflow`, `pkg/ragproduct`, provider construction, and tests.

### What should be done in the future
- Replace `study compile` direct-worker specifications with immutable Workflow input bundles plus a Researchctl experiment plan.
- Delete the direct runner, preview execution, old preparation workflow, and RAG-owned process lifecycle after the new compiler acceptance passes.

### Code review instructions
- Start with `analysis/01-execution-path-deletion-manifest.json`, then compare each entry with `cmd/rag-eval/cmds/study/command.go` and the earlier legacy-cleanup replacement matrix.
- Reproduce the live graph with `rg -l 'pkg/ragengine|pkg/researchctladapter|internal/workflow' --glob '*.go' --glob '!ttmp/**'`.

### Technical details
- Live direct path: `rag-eval study run -> researchctladapter.ExecuteSpecification -> researchctl execute-spec -> rag-worker -> ragengine`.
- Canonical replacement already proven by predecessor tickets: `researchctl experiment run-plan -> scraper-workflow-runner -> pkg/ragworkflow -> Workflow V3 -> ragoperators`.

## Step 5: Compile studies to Workflow plans and delete direct execution

This step performed the hard cut itself. `rag-eval study compile` now writes bounded immutable Workflow inputs plus a pure Researchctl experiment plan, and no RAG command starts scientific execution. The old direct worker, preview lifecycle, old preparation Workflow, RAG-owned process retry/progress adapter, and ambiguous top-level `workflow` command are gone.

The remaining document-intake product command was deliberately renamed to `intake`. It still has a separate product/API migration gate, but its help now explicitly says it is product-local and directs scientific studies to compile plus Researchctl/Workflow V3. This preserves supported intake behavior without presenting a second RAG study executor.

### Prompt Context

**User prompt (verbatim):** (same as Step 4)

**Assistant interpretation:** Implement the accepted destructive cutover and preserve detailed evidence, including failures encountered during the migration.

**Inferred user intent:** Make the canonical Researchctl-to-Workflow-V3 route the only ordinary RAG v2 study execution path.

**Commit (code):** `8925004` — "feat: hard cut RAG studies to Workflow V3"

### What I did
- Added `pkg/ragworkflow/study_export.go` with `rag-workflow-study-bundle/v1`, strict artifact-root containment, canonical per-cell execution/corpus/query/domain files, provider authority binding, digests, byte counts, and deterministic Researchctl plan rendering.
- Added stable provider-free and provider-backed bundle tests, output-boundary rejection, provider-required rejection, and no-secret/no-direct-runner assertions.
- Reworked `rag-eval study` to expose only `validate`, `explain`, and `compile`; compile resolves immutable envelopes, validates lineage, lowers each semantic cell, and emits Researchctl cases with the study's replicate count.
- Deleted `cmd/rag-worker`, `cmd/rag-preparation-smoke`, `cmd/rag-eval/cmds/preview`, `internal/preparationworkflow`, and direct run/progress/parity code in `pkg/researchctladapter`.
- Reduced `pkg/researchctladapter` to immutable input/catalog resolution and RAG semantic expansion, with new staged-input digest and round-trip tests.
- Renamed the separately gated product-local `rag-eval workflow` command to `rag-eval intake` and updated its help so it cannot be mistaken for canonical study execution.
- Removed the preview example and stale operator help, rewrote active README/help/guides/experiment/UI instructions, regenerated the embedded web build, and updated runnable-example coverage.
- Added `scripts/01-smoke-compiled-study.sh`, `scripts/02-cutover-guards.sh`, the deletion manifest, and captured a successful built-binary smoke in `analysis/02-compiled-study-smoke.txt`.

### Why
- Removing only `study run` would leave `study compile` producing a direct-worker specification and would therefore preserve the obsolete architecture in generated artifacts.
- Provider identity must be bound before the plan is immutable, while provider secrets remain host-only runner configuration.
- The old preparation package used Scraper's prior Workflow engine and duplicated execution custody already provided by Workflow V3.

### What worked
- A generated one-cell/two-replicate plan passed Researchctl `validate-plan`, executed twice through the built `rag-workflow-runner`, emitted verified artifacts and `rag.mrr`, recorded zero provider operations for both provider-free runs, and resumed as `executed=0,resumed=2`.
- Byte-stability tests proved identical plans and domain configs across independent roots.
- Provider-bundle tests proved exact `rag-v2-geppetto` authority inclusion without embedding provider configuration paths or fixture response content.
- Sequential full tests, full lint, web typecheck/build, cutover guards, and the pre-commit suite passed.

### What didn't work
- The first attempt to compile against `data/ttc-wordpress-rag.sqlite` failed with `Error: no such table: corpus_snapshot_documents`; that database is the raw TTC corpus, not the imported immutable RAG catalog expected by `NewTTCCatalog`. The permanent smoke therefore generates strict immutable envelopes and tests catalog-independent compilation.
- The first cross-repository validation invocation attempted `go run` on Researchctl outside the active module and failed with `directory .../researchctl/cmd/researchctl outside main module or its selected dependencies`. Building Researchctl from its own repository fixed the module boundary.
- The first smoke passed unsupported `--output` to the Glazed writer command and failed with `Error: unknown flag: --output`. The compiler already emits JSON to stdout, so the flag and stale help text were removed.
- The initial broad package test encountered the previously observed concurrent synthetic-provider timing failure: `expected: "succeeded" actual: "failed"` with an open query embedding operation. Focused reruns and the later sequential full suite passed; no cutover code touched provider runtime scheduling.
- The first full test after deleting `05-preview.js` failed in `TestRunnableExamples` with `open ../../../examples/rag-v2/05-preview.js: no such file or directory`. The active example list was corrected to exercise `06-raw-study.js`.
- Focused lint initially reported `cmd/rag-eval/cmds/study/command.go:28:6: type compiled is unused`; the obsolete direct-spec output DTO was removed.
- `make generate` failed with `No rule to make target 'generate'`; the repository's applicable command is `GOWORK=off go generate ./...`, which passed.

### What I learned
- A compiler output contract is part of the execution architecture: retaining direct-worker specifications would have kept a second backend even after removing the command that executed them.
- Input artifact file digests and domain manifest digests are distinct identities and both must remain intact through staging and Workflow bundle generation.
- The provider authority digest belongs in `scraper-workflow-execution/v2`; the provider configuration path and credentials do not.
- RAG's `replicates` field remains a domain request carried into the Researchctl plan, but only Researchctl expands and schedules those replicates.

### What was tricky to build
- Researchctl resolves input URIs relative to its artifact root, so allowing an arbitrary output directory would create plans that validate but cannot run. `WriteStudyBundle` canonicalizes both paths and rejects any relative path beginning with `..`.
- Every expanded execution binds the same immutable corpus/evaluation manifests but has distinct semantic factors and cell identity. The exporter writes exact execution bytes per cell and carries factor IDs into the generic case without reinterpreting RAG values.
- Provider-free and provider-backed task catalogs have different package identities. The exporter selects the lowerer and execution builder from the supplied provider package rather than mutating a provider-free plan after lowering.

### What warrants a second pair of eyes
- Review `renderStudyPlan` against Researchctl's JS builder contract, especially measure projection and factor IDs.
- Review whether the current default `maxConcurrent: 1` should remain conservative or become an explicit compile option; Researchctl still owns the policy either way.
- Confirm the renamed `intake` command adequately separates the deferred product path from scientific RAG v2 execution.

### What should be done in the future
- Complete the separately gated Workflow V3 migration of document-intake APIs/frontend before deleting `internal/workflow` and old Scraper engine imports.
- Use `TTC-SCRIPTED-EXPERIMENT-ACCEPTANCE` for the full scripted TTC study, not to restore a TTC-specific runner.

### Code review instructions
- Start at `cmd/rag-eval/cmds/study/command.go:RunIntoWriter`, then inspect `pkg/ragworkflow/study_export.go:WriteStudyBundle` and `renderStudyPlan`.
- Run `scripts/01-smoke-compiled-study.sh`, followed by `scripts/02-cutover-guards.sh`.
- Confirm `rag-eval study --help` lacks `run`, root help lacks `preview` and `workflow`, and root help contains `intake`.

### Technical details
- Bundle schema: `rag-workflow-study-bundle/v1`.
- Execution domain: `scraper-workflow` / `scraper-workflow-execution/v2`.
- Production runner: `scraper-workflow-runner` / `v1` through `cmd/rag-workflow-runner`.
- Code delta: 839 insertions and 3,124 deletions across 48 files.

## Step 6: Remove the final RAG-owned scheduling default

Reviewing ownership after the first successful smoke exposed one subtle leftover: the generated JavaScript plan explicitly selected blocked ordering, single-case concurrency, and fail-fast behavior. Those are valid Researchctl fields, but choosing them inside the RAG compiler would make RAG-eval a generic scheduling-policy owner again.

The renderer now returns only the cases, factors, and replicate requests. Researchctl normalizes its own ordering and execution defaults. A fresh two-replicate built-binary smoke passed after this correction.

### Prompt Context

**User prompt (verbatim):** (same as Step 4)

**Assistant interpretation:** Continue auditing the hard cut for hidden lifecycle ownership after the main implementation.

**Inferred user intent:** Ensure the cutover is architectural, not merely a command rename.

**Commit (code):** `68b1084` — "fix: leave experiment scheduling policy to Researchctl"

### What I did
- Removed `.ordering({strategy: "blocked"}).execution({maxConcurrent: 1, failFast: false})` from generated plans.
- Re-ran bundle tests, the one-case/two-replicate Researchctl smoke, lint, and the pre-commit suite.

### Why
- Ordering, execution concurrency, and fail-fast semantics are generic experiment lifecycle and belong to Researchctl.

### What worked
- Researchctl supplied canonical defaults, scheduled both replicates, and resumed both without RAG-selected scheduling policy.

### What didn't work
- N/A.

### What I learned
- Ownership leaks can exist in generated configuration even when no runtime loop remains in RAG code.

### What was tricky to build
- Replicate count is a legitimate domain study request that must cross into Researchctl, while ordering and concurrency are not RAG semantics. The final renderer preserves the former and omits the latter.

### What warrants a second pair of eyes
- Verify that all remaining fields in `renderStudyPlan` are identity or domain-case data rather than generic execution policy.

### What should be done in the future
- Add scheduling policy only in Researchctl-authored project/plan overlays when experiments require non-default behavior.

### Code review instructions
- Review the final callback in `renderStudyPlan` and run `scripts/01-smoke-compiled-study.sh`.

### Technical details
- Fresh smoke: one case, two replicates, `executed=2`, then `resumed=2`.

## Step 7: Fence compiled bundle custody against overwrite

A final artifact-custody review found that canonical bytes and digests were correct, but recompiling into an occupied output directory used `os.WriteFile` and could overwrite prior bytes. That would violate the immutable laboratory contract even though Researchctl later verifies its inputs.

The exporter now treats every case file, generated plan, and bundle manifest as immutable: identical recompilation is idempotent, while any byte difference at an existing path fails closed.

### Prompt Context

**User prompt (verbatim):** (same as Step 4)

**Assistant interpretation:** Continue the completion audit and fix any custody weakness before closing the cutover.

**Inferred user intent:** Preserve immutable run identity and prevent a compiler rerun from mutating evidence already referenced by a plan.

**Commits (code):** `f72f164` — "fix: preserve immutable compiled study custody"; `e119086` — "fix: exclusively create compiled study artifacts"

### What I did
- Added `writeImmutableStudyFile` for case files, plan JavaScript, and manifest JSON, using exclusive creation so a competing writer cannot replace an existing path.
- Added tests proving same-byte recompilation succeeds and changed corpus bytes at the same output path fail with `RAG_WORKFLOW_STUDY_CONFLICT`.
- Re-ran focused tests, focused lint, and the full pre-commit package suite.

### Why
- Digest metadata does not itself prevent mutation of the file behind a URI before Researchctl consumes it.

### What worked
- Idempotent compile produces the same bundle without rewriting semantics; conflicting compile fails before publishing a replacement manifest.

### What didn't work
- N/A.

### What I learned
- Immutable identity requires both canonical digest construction and write-once path behavior.

### What was tricky to build
- A rerun is expected during resume and should not fail merely because files exist. Exact byte comparison distinguishes safe idempotency from an identity collision.

### What warrants a second pair of eyes
- Review whether future concurrent compilers should wait and re-read a file another compiler is still writing; current exclusive creation is mutation-safe and may conservatively fail such a race.

### What should be done in the future
- If concurrent compilation to one directory becomes supported, replace read/check/write with an exclusive-create or content-addressed publication primitive.

### Code review instructions
- Review `writeImmutableStudyFile` and `TestWriteStudyBundleIsStableAndWorkflowNative`.

### Technical details
- Stable rerun: existing bytes equal canonical bytes, success.
- Conflict: existing bytes differ, `RAG_WORKFLOW_STUDY_CONFLICT`.

## Step 8: Complete the cross-repository acceptance and requirement audit

The final step mapped every ticket requirement to concrete source, test, command, and artifact evidence. Validation was rerun after the last custody change rather than relying on pre-commit history. The resulting audit distinguishes implemented completion from the separately gated document-intake migration.

Fresh evidence covers full tests, race detection, builds, lint, module tidiness, generation, frontend type/build checks, generated-plan execution and resume, deletion guards, and focused downstream Researchctl/Scraper suites.

### Prompt Context

**User prompt (verbatim):** (same as Step 4)

**Assistant interpretation:** Prove the current tree satisfies the complete cutover objective before closing its tasks and goal.

**Inferred user intent:** Avoid declaring completion based on implementation effort or proxy green checks.

### What I did
- Wrote `analysis/03-cutover-acceptance-audit.md` with an objective restatement and prompt-to-artifact checklist.
- Updated the design parity ledger from pending to evidence-backed provider-free/provider-backed status.
- Re-ran full sequential Go tests, affected-package race tests, all builds, full lint and Glazed vet, module tidy comparison, generation, web typecheck/build, both ticket scripts, and downstream suites.
- Re-inspected active strings/imports and documented the exact eight remaining old-Scraper imports as separately gated document-intake code.

### Why
- A successful smoke proves one path, not the absence of competing commands, stale instructions, mutable custody, or uncategorized callers.

### What worked
- All RAG final commands passed.
- The built generated plan executed two Researchctl-owned replicates and then resumed both.
- Scraper research runner/Workflow runtime/SQLite packages passed, including the 75-second runtime suite.
- Researchctl experiment plan, JS, service, laboratory, and process-runner packages passed.

### What didn't work
- The first downstream Researchctl command named a nonexistent package and failed with `pattern ./pkg/experiments/...: lstat ./pkg/experiments/: no such file or directory`. Inspecting `pkg/` showed the correct packages are `experimentplan`, `experimentplanjs`, and `experimentservice`; the corrected command passed along with `lab` and `lab/processrunner`.
- Vite continued to report its pre-existing warning that the 506.42 kB JavaScript chunk exceeds 500 kB. Type checking and production build succeeded; the cutover did not increase or redesign frontend bundling.
- The first final `docmgr doctor` reported unknown audit topics (`acceptance`, `researchctl`) and three missing related-file entries for paths intentionally deleted by the cutover. I changed the topics to repository vocabulary (`evaluation`, `research`) and removed the stale related-file entries; historical path names remain in prose and the deletion manifest.

### What I learned
- The old Scraper import count fell from the earlier 15 downstream files to eight after deleting direct execution and renaming the remaining product command. Every survivor belongs to document intake or its API, not study execution.
- Final completion evidence needs both positive execution and negative discoverability/import/package checks.

### What was tricky to build
- The completion audit had to avoid using predecessor parity as proof of the new compiler while also avoiding rerunning irrelevant provider internals. The new compiler received its own bundle tests and cross-repository smoke; predecessor evidence is used only for semantic/provider behavior deliberately reused unchanged.

### What warrants a second pair of eyes
- Review the requirement mapping in `analysis/03-cutover-acceptance-audit.md`, especially the classification of retained intake files and `ragengine` semantic-harness imports.

### What should be done in the future
- Execute `TTC-SCRIPTED-EXPERIMENT-ACCEPTANCE` as the next program ticket.
- Migrate document intake under its separate V3 product/API gate before deleting the final eight old-Scraper imports.

### Code review instructions
- Read the acceptance audit, run both scripts, and compare `go list ./...` plus root/study help against the deletion manifest.
- Inspect commits `8925004`, `68b1084`, `f72f164`, and `e119086` in order.

### Technical details
- Full Go tests: passed sequentially with `-p=1`.
- Race packages: `ragworkflow`, `researchctladapter`, study command, retained intake workflow, and API.
- Downstream Scraper runtime suite: 75.195 seconds.
- Smoke result: `executed=2`, `resumed=2`, `failed=0`, verified `rag.mrr`, zero provider operations.
