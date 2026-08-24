---
Title: Investigation diary
Ticket: RAG-V2-WORKFLOW-LOWERING
Status: active
Topics:
    - rag
    - rag-eval
    - workflow
    - scraper
    - javascript
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/internal/labsqlite/query.go
      Note: Authoritative ordinal trace export
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/lab/processrunner/resolver.go
      Note: Verified full-selector multi-input resolution
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/scraper/pkg/researchrunner/runner.go
      Note: Generic bounded set staging package linking and domain projection
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/scraper/pkg/researchrunner/types.go
      Note: Domain-neutral projector and set archive contracts
    - Path: repo://cmd/rag-workflow-runner/main.go
      Note: RAG-owned production process runner composition root
    - Path: repo://examples/rag-workflow/researchctl-plan.js
      Note: Two-case two-replicate generic Researchctl acceptance plan
    - Path: repo://pkg/ragworkflow/failure_test.go
      Note: Permanent failure mapping and canary redaction test
    - Path: repo://pkg/ragworkflow/lower.go
      Note: Deterministic canonical RAG execution to Workflow V3 lowering
    - Path: repo://pkg/ragworkflow/package.go
      Note: Versioned task package and immutable implementation bundle
    - Path: repo://pkg/ragworkflow/projector.go
      Note: Bounded privacy-safe domain observation projection
    - Path: repo://pkg/ragworkflow/registry.go
      Note: Closed exact operator/version lowering and provider attachment registry
    - Path: repo://pkg/ragworkflow/runtime.go
      Note: Provider-free preparation query reduction and publication runtime
    - Path: repo://pkg/ragworkflow/runtime_test.go
      Note: Semantic parity preparation-boundary restart and privacy tests
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-22T19:25:19.952904556-04:00
WhatFor: ""
WhenToUse: ""
---


# Investigation diary

## Goal

Record the evidence and decisions used to create the intern-facing design for **RAG-V2-WORKFLOW-LOWERING** and place it in the wider experiment-platform convergence program.

## Step 1: Establish the ticket and implementation guide

This step created the ticket as a focused child of the multi-repository convergence program. The guide explains the existing systems, proposed contracts, pseudocode, diagrams, implementation phases, tests, deletion rules, and neighboring tickets so an intern can begin without reconstructing the architecture from Git history.

No production code was changed. The architecture is proposed work; current files are referenced as evidence and implementation entry points.

### Prompt Context

**User prompt (verbatim):** "Create all the tickets with docmgr, then for each ticket, Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.


That way we will have everything. Have all the tickets reference each other so that if someone stumbles upon one ticket, they'll know what this is all about."

**Assistant interpretation:** Create the umbrella and all directed tickets, write self-contained technical guides, cross-link them, validate docmgr state, and publish the complete set to reMarkable.

**Inferred user intent:** Preserve a coherent program plan before more ad-hoc implementation occurs, while making every workstream independently discoverable and approachable by a new engineer.

### What I did
- Created `RAG-V2-WORKFLOW-LOWERING` with docmgr in the rag-eval repository.
- Added a long-term design/implementation guide and this diary.
- Mapped the ticket to all nine sibling/program tickets.
- Grounded the guide in current repository packages and public contracts.

### Why
- The convergence crosses repository boundaries and needs explicit ownership and sequencing.
- Standalone guides reduce the chance that an implementer recreates generic functionality in a workload package.

### What worked
- Existing greenfield contracts provided concrete API and file anchors.
- The ticket scope could be expressed as one independently testable capability: RAG v2 task catalog and compiler to Workflow V3.

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
- Validate documentation with `docmgr doctor --ticket RAG-V2-WORKFLOW-LOWERING --stale-after 30`.

### Technical details
- Program umbrella: `EXPERIMENT-PLATFORM-CONVERGENCE`.
- Ticket responsibility: RAG v2 task catalog and compiler to Workflow V3.

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

## Step 4: Backfill the production lowering, task package, runner, and first end-to-end acceptance

This step records the implementation work performed after the design was accepted and before the diary request was repeated. The work established the first complete provider-free route: canonical RAG v2 execution → deterministic Workflow V3 plan → a closed RAG task package → Scraper execution and observations → Researchctl immutable multi-case/multi-replicate runs. It also added restart projection, crash recovery, timeout cancellation, fixture generation, and semantic parity evidence.

The implementation deliberately kept RAG authoring and semantics in rag-eval, generic execution and observations in Scraper, and experiment scheduling/resume in Researchctl. It did not add provider behavior or a second scheduler, retry loop, artifact store, status system, or telemetry authority.

### Prompt Context

**User prompt (verbatim):** "keep a detailed diary as you work btwm and backfill if necessary, read the skill for the format if you forgot."

**Assistant interpretation:** Resume strict implementation-diary maintenance immediately, backfill all substantial work and failures since implementation began, and continue recording future steps in the required format.

**Inferred user intent:** Preserve enough chronological implementation evidence that another engineer can review, reproduce, or continue this multi-repository cutover without reconstructing decisions from diffs and terminal history.

**Commit (code):** `6c2640a` — "feat: support domain task packages and set inputs" (initial Scraper integration slice; subsequent changes described here were still uncommitted when this step was written)

### What I did
- Inspected the accepted lowering guide, tasks, canonical `ragcontract.PipelineExecution`, compiler targets, operator registry, `ragengine`, JavaScript authoring, preparation flow, deterministic fixtures, and Scraper Workflow V3 APIs.
- Extended Scraper's task-package boundary so a domain-owned package can contribute trusted runtime module factories without creating an import cycle.
- Extended `researchrunner.Config` with domain packages and a domain observation projector while retaining generic `scraper-workflow-execution/v2` framing.
- Added bounded `researchctl-set-input-archive/v1` materialization and a Workflow V3 item manifest for query map inputs; verified schema, digest, size, item count, item bytes, duplicate keys, traversal, and staging limits.
- Added strict validation and bounded projection for domain metrics and traces. Scraper still emits canonical workflow observations; RAG contributes only RAG measurements/traces.
- Implemented `pkg/ragworkflow`: versioned schemas; a closed operator/version/config registry; deterministic IR and Workflow V3 plans; exact task catalog/package identities; package/module registration; staged preparation; raw representations; fixture embeddings; lexical/vector index construction; query map; rank/fusion/evaluation; bounded reduction; result publication; preparation/index fingerprints; and canonical RAG observations.
- Exported `ragengine` prepared-value codecs so the new backend reuses domain semantics rather than copying private encodings.
- Added `cmd/rag-workflow-runner`, `cmd/rag-workflow-fixture`, and `cmd/rag-workflow-inspect` as built-binary composition roots.
- Added deterministic two-case fixtures and generated a Researchctl JavaScript plan with two replicates per case.
- Added Workflow execution parity tests against `ragengine`, including closing and reopening the SQLite database before observation/result verification.
- Added a permanent-failure test proving malformed corpus input maps to non-retryable `RAG_WORKFLOW_TASK_FAILED` and does not expose a payload canary.
- Fixed two generic Researchctl defects exposed by the richer workload: process inputs now permit repeated roles when `(role, kind, id)` differs, and exported traces are read in authoritative ordinal order rather than lexical kind order.
- Ran a built-binary smoke through Researchctl. The successful run proved four scheduled/executed runs, one injected runner crash and retry, five subordinate Workflow databases, exact `ragengine` result/metric parity for two queries in each case, 24 metrics, five traces, four verified artifacts per selected attempt, byte-equivalent restart reprojection, four-run immutable resume, and durable Workflow cancellation after a Researchctl timeout.

### Why
- RAG cases require several files under one logical input role; the old role-only uniqueness check could not represent that neutral use case.
- Domain task packages must be linked by the domain runner binary while Scraper retains generic execution ownership.
- Query sets must be bounded and materialized as authoritative Workflow map manifests rather than smuggled through an unbounded task payload.
- Existing `ragengine` outputs are the semantic oracle for the provider-free cutover.
- A fresh-process observation projection is necessary to prove that observation artifacts are derived from authoritative Workflow records and not process memory.

### What worked
- Focused Scraper tests and lint passed before the first Scraper commit.
- Deterministic lowering produced equal IR, plan, catalog, package, and digests on repeated compilation.
- Provider-backed and unknown operator variants failed closed.
- Provider-free Workflow execution matched `ragengine` retrieval results and evaluation metrics.
- The third full cross-repository smoke completed with the following summary:
  - `scheduled=4`, `executed=4`, `failed=0`;
  - Researchctl attempt counts `[1,1,1,2]` after one process crash;
  - five subordinate Workflow databases;
  - two cases × two queries with matching `ragengine` results and `rag.mrr`;
  - 24 metrics, five traces, and four verified artifacts per successful run;
  - resume `executed=0`, `resumed=4`;
  - timeout mapped to Researchctl failure kind `timeout` and durable Workflow status `canceled`;
  - freshly projected observations equaled the exported canonical observation artifact.
- Changing the JavaScript task wrapper from `return task.failure(...)` to `throw task.failure(...)` made Scraper's typed failure classifier receive the intended class/code/retryability contract.

### What didn't work
- First cross-repository command:
  - command: built `rag-workflow-runner` and `researchctl`, initialized the lab, then ran the generated four-run plan;
  - error: `PROCESS_RUNNER_INPUT_DUPLICATE: inputs[1] duplicates role "workflow-input"`;
  - cause: `VerifiedInputResolver` incorrectly treated role as globally unique even though `ArtifactRef` identity includes role, kind, and ID;
  - fix: reject only duplicate `(role, kind, id)` selectors and add positive/negative resolver coverage.
- Second cross-repository command:
  - error: `validate exported run: attempts[0].traces[0]: invalid ordinal, kind, or value` for all four runs;
  - evidence: SQLite held ordinal traces `workflow.*` at 1–3 and `rag.query` at 4–5, but export queried `ORDER BY trace_kind,ordinal`, placing `rag.query` first with ordinal 4;
  - fix: query traces by `ORDER BY ordinal` and add a regression test that writes lexical order `z-*` then `a-*` and validates preserved ordinals.
- First smoke-script acceptance pass executed all workloads but the final Python verifier failed before evaluating evidence:
  - exact error: `SyntaxError: '(' was never closed` at the restart-reprojection load;
  - fix: close the `json.load(open(...))` expression and rerun the entire smoke from clean temporary state.
- First permanent-failure unit test failed at `require.True(t, errors.As(err, &failure))`:
  - cause: `task.cjs` returned a failure object, so it was handled as an invalid normal result rather than a typed thrown task failure;
  - fix: use `throw task.failure(...)`, rerun the focused test, regenerate all package-identity-dependent fixtures, and plan to rerun acceptance.
- The first generated Researchctl plan was created ad hoc. This was corrected by adding a checked-in deterministic generator and comparing generated output byte-for-byte in the smoke.

### What I learned
- Researchctl's exported-run validator made an implicit persistence-order bug observable only when a domain projector added trace kinds that sort before canonical Workflow kinds.
- Generic input-role cardinality must not be inferred from a display/selection role; hard identity uses the complete selector.
- Scraper's task failure helper is an exception value. Returning it violates the task ABI even though the JavaScript looks superficially plausible.
- Task package script bytes participate in package/catalog/domain identities, so even a one-word script correction requires fixture and plan regeneration.
- Preparation identity and query dataset identity must be explicit at the artifact boundary; relying only on the outer Researchctl artifact digest is insufficient for domain semantic custody. Work began immediately after this diary backfill to place the declared dataset manifest digest into every query item.

### What was tricky to build
- The package composition direction is easy to invert. Scraper cannot import rag-eval, and rag-eval cannot ask Scraper's default binary to know RAG modules. The solution was a generic package/module contribution interface in Scraper and a RAG-owned runner binary that links the package.
- Workflow map nodes consume item manifests, while Researchctl supplies verified files. The runner therefore needs a bounded, strict archive-to-manifest staging step. It must verify all outer and inner identities before submission and must not put query text into events or generic observations.
- Semantic parity includes preparation and query semantics but excludes wall-clock values. The acceptance compares result ordering, result traces, and evaluation metrics, while canonical Workflow timing remains independent execution evidence.
- Crash injection had to occur after `workflow.submitted` so the subordinate run existed durably, but before terminal frames reached Researchctl. This proved that one Researchctl retry creates a second Workflow run without mutating or silently resuming the crashed attempt's run.
- Cancellation had to block all RAG resources without changing the plan. Passing only `blocked.resource=1` overrides defaults, leaves every RAG node waiting, and gives the timeout path time to persist a canceled run and observation.

### What warrants a second pair of eyes
- Review `pkg/ragworkflow/lower.go` and its closed registry for accidental acceptance of semantically equivalent but noncanonical configs.
- Review preparation/index fingerprints and dataset-manifest propagation for complete stale-reuse rejection.
- Review the set-input archive limits and temp-file cleanup for denial-of-service and partial-materialization behavior.
- Review domain projector bounds and redaction guarantees; query traces intentionally include query IDs and retrieval identities but not query text or source corpus text.
- Review crash semantics: the intentionally orphaned first Workflow remains a separate durable run, while Researchctl correctly selects only the successful retry attempt.
- Review the Researchctl input-selector and trace-order fixes as generic contract corrections, not RAG-specific exceptions.

### What should be done in the future
- Finish dataset-manifest binding through query archive items and add rejection tests for stale/mismatched dataset identity.
- Regenerate fixtures after identity-affecting changes and rerun the full smoke.
- Complete package docs, help, examples, ticket tasks/changelog/file relations, and the acceptance report.
- Run full tests, races, builds, lint, module-tidy/generated guards, frontend checks where applicable, downstream acceptance, privacy/deletion guards, and `docmgr doctor`.
- Remove the temporary absolute `go.mod` replacement before final commits and choose an auditable cross-repository dependency state.
- Inspect all diffs, make focused commits in Scraper, Researchctl, and rag-eval, and ensure all trees are clean.

### Code review instructions
- Start in rag-eval at `pkg/ragworkflow/lower.go`, `package.go`, `runtime.go`, `projector.go`, and `runtime_test.go`.
- Then inspect Scraper at `pkg/researchrunner/runner.go`, `types.go`, and the task-package/runtime composition changes.
- Inspect Researchctl at `pkg/lab/processrunner/resolver.go` and `internal/labsqlite/query.go`.
- Run focused tests:
  - `GOWORK=off go test ./pkg/ragworkflow -count=1` in rag-eval;
  - `GOWORK=off go test ./pkg/researchrunner ./pkg/workflowv3product ./pkg/workflowv3observations -count=1` in Scraper;
  - `GOWORK=off go test ./pkg/lab/processrunner ./internal/labsqlite -count=1` in Researchctl.
- Run built-binary acceptance:
  - `bash ttmp/2026/07/22/RAG-V2-WORKFLOW-LOWERING--*/scripts/02-smoke-rag-workflow-lowering.sh`.

### Technical details
- Domain contract: `scraper-workflow-execution/v2`.
- Canonical Workflow observation contract: `scraper-workflow-observations/v1`.
- RAG package identity: `rag-v2-provider-free@1.0.0`; module alias: `rag:workflow`.
- RAG result schema: `rag-workflow-result/v1`; task identities use exact kind/version pairs plus immutable implementation digests.
- Research fixture plan: two cases, two replicates, blocked ordering, concurrency two.
- Temporary successful smoke work directory retained for evidence during this step: `/tmp/tmp.0lit5Tg6wk`.

## Step 5: Close identity gaps, harden privacy, and run repository-wide acceptance

This step completed the semantic and operational hardening discovered during review of the first passing cross-repository slice. It replaced prefix-based operator acceptance with an explicit versioned lowering registry, bound every query artifact to the canonical dataset manifest, proved restart precisely at the preparation/query boundary, bounded aggregate domain projection bytes, and removed potentially sensitive metadata from Researchctl observation frames.

It also converted the successful interactive smoke into reproducible ticket scripts and durable evidence, ran full validation across all three repositories, fixed generic defects and lint blockers exposed by the new workload, and created focused commits. The implementation now depends on a published Scraper pseudo-version rather than a machine-local `replace` directive.

### Prompt Context

**User prompt (verbatim):** (the diary requirement is recorded verbatim in Step 4; this step continues the active `RAG-V2-WORKFLOW-LOWERING` goal)

**Assistant interpretation:** Continue implementation and verification without repeating completed work, maintain the diary, and do not close the durable goal until every requirement has fresh evidence.

**Inferred user intent:** Finish the cross-repository lowering as a reviewable production cutover with reproducible evidence, not merely a locally passing prototype.

**Commit (code):**
- Scraper `c497bc1` — "feat: project bounded domain observations"
- Scraper `391c2d0` — "test: keep domain projection fixture neutral"
- Researchctl `8d68226` — "fix: preserve generic input and trace identities"
- RAG-eval `5ec4f16` — "chore: resolve repository lint findings"
- RAG-eval `f7cd54a` — "feat: lower RAG v2 into Workflow V3"
- RAG-eval `33bcffa` — "test: add RAG Workflow cross-repository acceptance"
- RAG-eval `698928c` — "fix: preserve canonical generated fixture bytes"

### What I did
- Added `OperatorRegistry` with exact `(kind, version)` entries, static/query phases, task identities, and explicit provider-required attachment points. Removed prefix-based support decisions.
- Changed fixture construction to compute real corpus and dataset semantic digests before computing execution `CellID`.
- Changed query archive items from bare queries to `rag-workflow-query/v1` envelopes containing the exact dataset manifest digest. Archive creation and query execution reject mismatches.
- Changed the parity restart test to execute every preparation task, close SQLite, reopen it, and only then execute query map, reduction, and publication.
- Added a permanent task-failure test and corrected `task.cjs` to throw, rather than return, `task.failure`.
- Added aggregate byte accounting for domain outputs and projected metrics/traces in Scraper's generic runner.
- Added privacy-safe RAG projection: query metadata and per-hit filters are omitted; failure messages/details are redacted; metric metadata is replaced by a semantic digest. Added a canary assertion over the complete projected frame set.
- Added `rag-workflow-inspect` for fresh-process observation reprojection.
- Added deterministic fixture/plan generation, built-binary crash/retry/resume/timeout smoke, and boundary/privacy/deletion guards.
- Added operator-facing embedded help, README guidance, acceptance evidence, and a review route.
- Pushed the Scraper task branch so RAG-eval could replace the temporary absolute module replacement with `github.com/go-go-golems/scraper v0.0.5-0.20260723191048-c497bc1245cc`.
- Ran full tests, focused race tests, builds, full lint, module-tidy checks, command/help smoke, cross-repository smoke, and guards.

### Why
- A prefix such as `retrieve.*` is not a closed lowering registry. A future compiler operator could otherwise become accepted without a reviewed lowering.
- Outer file custody alone does not prove that map items belong to the dataset named by the canonical execution.
- Restart before any work does not test reusable preparation artifacts. Restart after preparation proves the intended boundary.
- Per-frame limits still allow an unbounded aggregate projection. Aggregate limits are required before invoking a domain projector or emitting its frames.
- Canonical RAG traces may contain arbitrary query metadata, filters, failure messages, details, or measure metadata. Those values belong in verified domain artifacts, not generic observation frames.

### What worked
- `GOWORK=off go test ./... -count=1` passed in RAG-eval and Researchctl on the first full runs.
- Focused race suites passed in all three repositories.
- Full builds and module-tidy no-diff checks passed in all three repositories.
- Full Scraper and Researchctl lint passed. Full RAG lint passed after the repository findings below were corrected.
- Final cross-repository smoke passed with four selected runs, attempt counts `[1,1,1,2]`, five Workflow databases, exact provider-free parity, zero provider operations, complete resume, durable timeout cancellation, verified artifacts, and restart reprojection equality.
- Final boundary/privacy/deletion guard passed.
- A direct-module dependency replaced the temporary local path and `GOWORK=off` remained green.

### What didn't work
- Running all three full Go suites simultaneously caused one Scraper timing failure:
  - test: `TestHTTPSnapshotRetriesAndReopensWithoutPersistingRequestSecrets`;
  - symptom: expected two attempts but observed a third `lease_lost` attempt after the test lease exceeded two seconds under machine contention;
  - exact assertion summary included `should have 2 item(s), but has 3`;
  - triage: the focused test passed five consecutive runs; the complete Scraper suite then passed sequentially in 97.202 seconds.
- Running three `golangci-lint` processes simultaneously caused RAG lint to exit with `Error: parallel golangci-lint is running`.
- The subsequent sequential full RAG lint exposed six existing findings: four unchecked `Close` calls, a non-exhaustive workflow-status switch, and `func max` shadowing a predeclared identifier. I corrected these with explicit ignored close results, explicit pending/running cases, and `maxInt`; full lint passed.
- The first command/help smoke redirected only stdout. Go's `flag` package printed usage to stderr, so the grep files were empty and the command exited 1. Redirecting `2>&1` fixed the test; all help checks passed.
- The first boundary guard found RAG-specific names in Scraper's generic projector unit test (`rag.mrr`, `rag.query`). The implementation was generic, but the fixture vocabulary violated the boundary. I renamed it to `domain.score` and `domain.item`, committed, pushed, and reran the guard successfully.
- The fixture acceptance commit triggered the repository's Biome pre-commit formatter. Biome expanded canonical one-line JSON and JavaScript after their byte digests had been generated, so `scripts/03-boundary-privacy-deletion-guards.sh` exited 1. I added `!examples/rag-workflow` to Biome's generated-file exclusions, regenerated all fixtures and the plan, confirmed Biome checked zero generated files, and committed the canonical bytes.
- Immediately after pushing Scraper, ordinary `go get` failed while the new pseudo-version had not reached the public checksum service: `reading https://sum.golang.org/lookup/...: 500 Internal Server Error`. Retrying with `GOPROXY=direct GONOSUMDB='github.com/go-go-golems/*'` fetched the exact commit and produced stable module hashes.

### What I learned
- Generated scientific fixture bytes must be excluded from formatters when manifests custody those exact bytes.
- The operator registry needs to represent provider-required entries, not merely omit them, so the subsequent operations phase has explicit attachment points and callers receive a deliberate provider-required error.
- Domain artifacts and generic observations have different privacy boundaries. The full result artifact preserves semantic evidence; Researchctl projections carry bounded identities, measurements, and digests.
- Concurrent repository-wide validation can invalidate short lease assumptions. A focused repeated test plus sequential full rerun distinguishes scheduler contention from a deterministic defect.
- A pushed feature-branch pseudo-version is preferable to committing a local replacement in a cross-repository integration.

### What was tricky to build
- Dataset identity is recursive: the dataset digest must be known before execution `CellID`, while fixture files and outer file digests are created after execution. The solution computes the semantic corpus/dataset digests first, updates bindings, computes `CellID`, then writes byte-custodied files and their separate SHA-256 file digests.
- Privacy-safe projection had to retain metric lineage without emitting raw metadata. The projector computes a canonical metadata digest and includes the result/query digests, while keeping raw metadata only in the verified result artifact.
- Preparation restart needed deterministic control over task order. Calling `Engine.RunOne` exactly `len(staticNodes)+1` times executes corpus load plus every static node; the test asserts only successful `prepare-*` attempts exist before closing the store.
- The task-package script is embedded in implementation identity. Correcting `return task.failure` to `throw task.failure` changed package/catalog/plan digests, requiring complete fixture and Researchctl-plan regeneration.

### What warrants a second pair of eyes
- Confirm the exact provider-required registry is the intended handoff surface for `RAG-GEPPETTO-WORKFLOW-OPERATIONS`.
- Confirm that query IDs are safe experiment identifiers; the projector allows only 64 characters from `[A-Za-z0-9_-]` and includes them in metric scope while hashing them in metadata/traces.
- Review `PreparedBundle.Indexes` reconstruction against multi-index implementation changes; exact manifest, record digest, and size are checked.
- Review the intentional lack of cross-run preparation import. Reuse in this phase is durable within one Workflow across restart; any future cross-run reuse must add a canonical prepared input and preserve the same fingerprint checks.
- Review the generated-file Biome exclusion to ensure future generated RAG fixtures remain byte-custodied.

### What should be done in the future
- Implement real provider operations only in `RAG-GEPPETTO-WORKFLOW-OPERATIONS` using the explicit provider-required registry entries and Scraper external-operation custody.
- Use TTC only after the provider operation phase is accepted.
- If cross-run preparation sharing is desired, design it as an explicit immutable artifact input rather than a hidden cache lookup.

### Code review instructions
- Begin with `pkg/ragworkflow/registry.go`, `lower.go`, `runtime.go`, `projector.go`, and `runtime_test.go`.
- Review `examples/rag-workflow/provider-free/manifest.json` together with the generator and Biome exclusion.
- Review the generic changes at Scraper commits `ebf9a85`, `c497bc1`, `391c2d0` and Researchctl commit `8d68226`.
- Read `analysis/01-rag-workflow-lowering-acceptance-and-boundary-audit.md` for the requirement-to-evidence map.
- Validate with:
  - `GOWORK=off go test ./... -count=1`;
  - `GOWORK=off golangci-lint run ./...`;
  - `bash .../scripts/02-smoke-rag-workflow-lowering.sh`;
  - `bash .../scripts/03-boundary-privacy-deletion-guards.sh`.

### Technical details
- RAG task package: `rag-v2-provider-free@1.0.0`.
- Scraper dependency: `v0.0.5-0.20260723191048-c497bc1245cc`.
- Query item schema: `rag-workflow-query/v1`.
- Result schema: `rag-workflow-result/v1`.
- Canonical observation schema: `scraper-workflow-observations/v1`.
- Final smoke evidence: `sources/smoke/01-summary.json`.

## Step 6: Publish the acceptance record and close the ticket

This step converted the final implementation and validation evidence into the permanent operator help, acceptance audit, design outcome, ticket metadata, file relations, tasks, and changelog. The ticket was closed only after every task was checked and `docmgr doctor` passed without warnings.

The documentation commit captures the complete review route and known validation failures rather than presenting only the green endpoint. Ticket closure then changed the ticket status from active to complete and refreshed its changelog/index metadata.

### Prompt Context

**User prompt (verbatim):** (same active-goal continuation and diary requirement as Steps 4–5)

**Assistant interpretation:** Complete documentation/bookkeeping and close only after evidence-backed validation.

**Inferred user intent:** Leave a durable, discoverable handoff whose status accurately reflects production readiness.

**Commit (code):** `189c146` — "docs: complete RAG Workflow lowering ticket"

### What I did
- Added embedded operator help and README usage.
- Added the requirement-to-evidence acceptance and boundary audit.
- Marked the accepted design outcome and linked production files.
- Checked all ticket tasks with docmgr.
- Related 13 material RAG, Scraper, and Researchctl files to the diary with absolute file notes.
- Added cross-repository commit hashes and acceptance evidence to the changelog.
- Corrected acceptance-report topics to the repository vocabulary after doctor warnings.
- Ran `docmgr doctor --ticket RAG-V2-WORKFLOW-LOWERING --stale-after 30` until all checks passed.
- Closed the ticket with `docmgr ticket close --ticket RAG-V2-WORKFLOW-LOWERING`.

### Why
- Production code without operator contracts and reproduction commands is not a complete cutover.
- Ticket status must follow, not precede, validation and evidence.

### What worked
- All ticket tasks were complete.
- Final doctor output was `✅ All checks passed`.
- Ticket closure reported `Status: active → complete` and updated the changelog.

### What didn't work
- The first doctor run warned that `researchctl`, `semantic-parity`, and `workflow-v3` were not registered topic vocabulary values. I replaced them with existing precise topics: `rag`, `workflow`, `scraper`, and `evaluation`; the second doctor run passed.

### What I learned
- Acceptance documents are validated against repository-local vocabulary, even when cross-repository terms are technically accurate.

### What was tricky to build
- The audit had to distinguish evidence captured before the final docs-only changes from checks that needed fresh reruns. Generated contracts, code, dependency state, and guards were final before closure; documentation changes did not alter runtime behavior.

### What warrants a second pair of eyes
- Review the acceptance matrix against the durable goal language and confirm no provider-execution requirement was accidentally moved into this provider-free phase.

### What should be done in the future
- Begin `RAG-GEPPETTO-WORKFLOW-OPERATIONS` from the explicit provider-required attachment points.

### Code review instructions
- Read the acceptance audit first, then follow its six-step review route.
- Run doctor and the two ticket scripts before merging.

### Technical details
- Documentation commit: `189c1466bde629964df6d19890261bd6b326fede`.
- Ticket state: complete.
- Doctor state: all checks passed.

## Step 7: Audit and prove non-empty factor identity preservation

The final requirement audit noticed that factor preservation was implemented by copying canonical execution factors into the result, but the provider-free fixture used an empty factor list. This step added explicit evidence rather than relying on code inspection.

The restart/parity test now injects a non-empty factor selection, recomputes the canonical execution `CellID`, executes the complete preparation/restart/query/publication flow, and asserts exact factor equality in the terminal result.

### Prompt Context

**User prompt (verbatim):** (same active-goal continuation as Steps 4–6)

**Assistant interpretation:** Audit every explicit requirement against concrete test evidence before marking the durable goal complete.

**Inferred user intent:** Prevent a superficially complete implementation from omitting identity fields that happen to be empty in fixtures.

**Commit (code):** `2cb0c9e` — "test: preserve RAG factor identity through publication"

### What I did
- Added a `retrieval-profile=hybrid` canonical factor with structured value to the runtime parity test.
- Recomputed `CellID` after factor insertion.
- Asserted `execution.Factors == workflowResult.Factors` after preparation-boundary restart and publication.
- Reran focused tests, race tests, lint, and the repository pre-commit suite.
- Updated the acceptance audit with the factor evidence and commit.

### Why
- The goal explicitly requires exact factor identities; empty-slice fixtures cannot prove non-empty preservation.

### What worked
- Focused package tests, `-race`, and lint passed.
- The pre-commit package/internal test and lint suites passed.

### What didn't work
- N/A.

### What I learned
- Final requirement mapping should distinguish “field exists in code” from “non-empty value crosses the entire execution boundary under test.”

### What was tricky to build
- Changing factors changes canonical execution identity. The test must clear and recompute `CellID`; otherwise the lowerer correctly rejects the stale execution before the factor assertion is reached.

### What warrants a second pair of eyes
- Confirm future factor values remain canonical JSON and continue participating in `CellID` and result digest.

### What should be done in the future
- N/A.

### Code review instructions
- Review the setup and terminal assertions in `TestWorkflowExecutionMatchesRAGEngineAndSurvivesRestart`.
- Run `GOWORK=off go test -race ./pkg/ragworkflow -count=1`.

### Technical details
- Factor ID: `retrieval-profile`.
- Value ID: `hybrid`.
- Structured value: `{"channels":2}`.
