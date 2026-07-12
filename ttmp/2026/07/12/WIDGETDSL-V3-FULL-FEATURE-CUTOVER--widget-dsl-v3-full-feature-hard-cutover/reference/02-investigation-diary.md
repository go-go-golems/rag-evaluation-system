---
Title: Investigation Diary
Ticket: WIDGETDSL-V3-FULL-FEATURE-CUTOVER
Status: active
Topics:
    - widget-dsl
    - ui-dsl
    - widget-ir
    - goja
    - xgoja
    - react
    - design-system
    - frontend-architecture
    - typescript
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://ttmp/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER--widget-dsl-v3-full-feature-hard-cutover/design-doc/01-widget-dsl-v3-full-feature-analysis-design-and-intern-implementation-guide.md
      Note: Design synthesis produced by the investigation
    - Path: repo://ttmp/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER--widget-dsl-v3-full-feature-hard-cutover/scripts/01-inventory-widget-dsl.py
      Note: Repeatable evidence inventory script created during investigation
    - Path: repo://ttmp/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER--widget-dsl-v3-full-feature-hard-cutover/sources/01-generated-runtime-inventory.md
      Note: Generated helper namespace and registry counts
    - Path: repo://ttmp/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER--widget-dsl-v3-full-feature-hard-cutover/sources/02-v3-example-migration-check.txt
      Note: Raw escape-hatch findings in golden examples
    - Path: repo://ttmp/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER--widget-dsl-v3-full-feature-hard-cutover/sources/03-v3-host-migration-check.txt
      Note: Clean example-host migration check
ExternalSources: []
Summary: Chronological evidence and design diary for the Widget DSL v3 full-feature hard-cutover research and intern implementation package.
LastUpdated: 2026-07-12T19:45:00-04:00
WhatFor: Preserve the commands, evidence, failures, decisions, and review instructions behind the v3 parity and hard-cutover design.
WhenToUse: Before resuming this ticket, reviewing its recommendations, or implementing any phase of the cutover.
---


# Diary

## Goal

Capture how the Widget DSL v3 full-feature hard-cutover design was assembled from repository evidence, which assumptions changed, what remains risky, and how a new engineer should validate and continue the work.

## Step 1: Create the Ticket and Inventory the Runtime Surfaces

This step created a dedicated docmgr workspace inside `rag-evaluation-system/ttmp` and established an evidence baseline before recommending API changes. The investigation treated the legacy modules, v3 runtime, typed Go specs, provider package, React registry, Widget IR, examples, and migration tooling as one system rather than auditing only JavaScript helper names.

A generated inventory script was added under the ticket so future work can refresh helper and registry counts. Parser-backed migration checks then identified where v3 examples still rely on raw component construction and distinguished stale examples from genuine missing APIs.

### Prompt Context

**User prompt (verbatim):**

> Ok, creat ea detailed docmgr ticket in rag-evaluation-system/ttmp itself, for bringing full feature set to the new DSL, including the new designs above, and keep a consistent and clean and composable and also opinionated DSL (so, reorganize things if necessary, we can do hard cutover since v3 is not yet released publicly).
>
> Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.

**Assistant interpretation:** Create a new repository-local ticket containing an exhaustive architecture/design/implementation package for completing Widget DSL v3, freely reorganizing unreleased APIs, integrating the previously designed collection interactions, documenting the system for an intern, validating with docmgr, and publishing a bundled PDF to reMarkable.

**Inferred user intent:** Establish a durable technical blueprint that can guide a new engineer through a coordinated, opinionated hard cutover instead of allowing v3 to accumulate another inconsistent layer of compatibility helpers.

### What I did

- Read repository `AGENTS.md` and package design-system guidelines.
- Created ticket `WIDGETDSL-V3-FULL-FEATURE-CUTOVER` with design, parity reference, diary, tasks, changelog, and standard ticket directories.
- Inspected legacy helper maps and recipes in `pkg/widgetdsl/module.go`.
- Inspected v3 namespaces in `v3.go`, `v3_crm.go`, descriptors, TypeScript declarations, typed collection specs, validation, and lowering.
- Inspected provider registration and tests.
- Inspected Widget IR actions/cells/props/engines, WidgetRenderer, App action flow, and the default registry.
- Inspected DataTable, Pagination, SearchField, DialogShell, TagListInput, and ActivityFeed.
- Added `scripts/01-inventory-widget-dsl.py` and generated `sources/01-generated-runtime-inventory.md`.
- Ran migration checks over golden examples and example hosts.

### Why

- “Full feature” cannot be designed safely from one application request; it requires a complete view of all existing layers and migration constraints.
- Helper counts and raw uses provide concrete parity evidence.
- The hard cutover must account for xgoja provider selection and generated declarations, not only runtime calls.

### What worked

- `docmgr status --summary-only` confirmed the repository-local ticket root and vocabulary.
- Existing vocabulary already contained all required topics.
- The inventory script found 87 legacy direct component helpers and 87 registered adapters across registry groups.
- The v3 host migration check returned: `No legacy Widget DSL imports or raw component escape hatches found.`
- The example migration check produced 11 actionable raw findings.
- Existing v3 descriptors, typed collection specs, migration checker, golden corpus, and provider tests provide strong foundations.

### What didn't work

- An assumed provider path did not exist:

  `ENOENT: no such file or directory, access '/home/manuel/code/wesen/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite/widgetsite.go'`

  The actual provider entrypoint is `pkg/xgoja/providers/widgetsite/provider.go`.

- An assumed ActivityFeed organism path did not exist:

  `ENOENT: no such file or directory, access '/home/manuel/code/wesen/go-go-golems/rag-evaluation-system/packages/rag-evaluation-site/src/components/organisms/ActivityFeed/ActivityFeed.tsx'`

  ActivityFeed is correctly a molecule under `components/molecules/ActivityFeed`.

- The first inventory parser looked for the time namespace variable `time`; the implementation uses `timeObj`, so the generated report initially showed zero time exports. The script was corrected to parse `timeObj` and regenerated, yielding seven exports.

### What I learned

- React registry breadth and DSL quality are separate: every component is renderable, but not every component should become a top-level public helper.
- ActivityFeed is explicitly domain-blind and registered in the data registry even though its v3 helper lives under CRM.
- The golden example corpus has drifted behind runtime improvements; raw-use counts include both missing APIs and stale examples.
- Typed `v2/spec` is already the implementation core of v3 collections, so final cutover should rename/reorganize rather than replace it.
- Provider tests currently enshrine split-module availability and must change late in the migration, not at the beginning.

### What was tricky to build

- Regex-based source inventory is intentionally evidence support, not a Go parser. The script handles known map/function shapes and must not become the authoritative API source. The design therefore recommends descriptor/runtime parity tests in Go as the durable replacement.
- “Legacy capability missing from v3” required classification: some are genuinely absent, some are represented semantically, some should stay internal, and some examples are simply stale.

### What warrants a second pair of eyes

- Verify the generated helper counts if module maps are reorganized before implementation begins.
- Review whether all first-party repositories selecting the provider are included in future migration scans.
- Review the disposition of generic app shell/sidebar and semantic list helpers; these have the greatest risk of either underexposure or API duplication.

### What should be done in the future

- Replace the ticket-local regex inventory with descriptor-backed Go parity checks during Phase 1.
- Add workspace-wide migration scans before deleting provider modules.
- Keep raw-use findings classified rather than enforcing a blind zero rule on experimental fixtures.

### Code review instructions

- Start with `sources/01-generated-runtime-inventory.md` and compare it to `module.go`, `v3.go`, `v3_crm.go`, and `defaultRegistry.ts`.
- Review migration outputs in `sources/02-v3-example-migration-check.txt` and `sources/03-v3-host-migration-check.txt`.
- Run:

  ```bash
  ttmp/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER--widget-dsl-v3-full-feature-hard-cutover/scripts/01-inventory-widget-dsl.py
  go run ./cmd/widgetdsl-migration-checker --root . -- pkg/widgetdsl/testdata/v3/examples
  ```

### Technical details

- Legacy direct component helpers: 87.
- Default registry adapters: 87.
- V3 golden examples: 41.
- V3 example raw findings: 11.
- V3 example host findings: 0.

## Step 2: Synthesize the Hard-Cutover Language and Implementation Design

This step converted the inventory into a complete target language and implementation sequence. The design deliberately avoids both extremes: it does not mirror every React component as a public factory, and it does not restrict v3 to a few domain recipes that force ordinary applications back to `widget.raw`.

The resulting architecture has stable tiers: generic page/UI/data grammar, typed engines and collection shaping, domain views and intents, and an explicit raw escape hatch. It integrates keyboard commands, overlay dialogs, progressive search, pagination/page size, conditional styling, activity timelines, descriptor parity, browser state ownership, and final provider deletion into one phased plan.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Turn evidence into an intern-ready hard-cutover specification with concrete APIs, typed specs, renderer boundaries, phases, tests, risks, and definition of done.

**Inferred user intent:** Make future implementation converge on one clean language rather than solving each application request with disconnected component helpers.

### What I did

- Wrote `reference/01-legacy-to-v3-feature-parity-inventory.md` with per-family dispositions.
- Wrote `design-doc/01-widget-dsl-v3-full-feature-analysis-design-and-intern-implementation-guide.md`.
- Defined public grammar tiers and namespace ownership.
- Designed complete collection shaping and arrangement specs.
- Integrated keyboard row commands, FormDialog/overlays, progressive search, pagination, page-size control, structured navigation, action result policy, semantic styles, and ActivityFeed promotion.
- Added Mermaid architecture, state-flow, and migration diagrams.
- Added decision records, nine implementation phases, file-level maps, tests, risks, alternatives, open questions, and definition of done.

### Why

- The implementation crosses Goja, Go specs, TypeScript declarations, Widget IR, React adapters, browser state, server actions, provider packaging, examples, and first-party hosts.
- Interns need not only an API sketch but also a map explaining why each layer exists and where behavior belongs.
- A hard cutover is safest when deletion criteria and migration gates are defined before code changes start.

### What worked

- The existing composition-grammar and v3 migration documents gave consistent principles: one module, intent-level views, scoped builders, named slots, bindings, intents, and author-time lambdas.
- The current collection lowering provides a concrete place to integrate search/pagination around arrangements.
- Existing Pagination, SearchField, DialogShell, ActivityFeed, and DataTable components reduce new React work.
- Existing `App.tsx` already centralizes server-action refresh/toast behavior and URL-based page refetching.

### What didn't work

- A simple “add all missing legacy names to `widget.ui`” strategy was rejected because it would reproduce the component-catalog design.
- Keeping ActivityFeed solely under CRM was rejected because source comments and registry ownership identify it as generic.
- Preserving split modules through aliases was rejected because the user explicitly authorized hard cutover and v3 is unreleased.

### What I learned

- Full feature parity is best expressed as an explicit disposition matrix, not a helper-count equality target.
- Search and pagination are collection shaping; keyboard commands are arrangement mechanics; dialogs are renderer services; this separation keeps CollectionSpec composable.
- Descriptor completeness is not documentation polish—it is the mechanism that prevents runtime/declaration/example drift.
- React component stories and browser interaction tests remain necessary even when golden IR is stable.

### What was tricky to build

- CollectionSpec can become too broad. The design uses nested `CollectionShapingSpec` and arrangement-specific specs so search/pagination can apply across tables/cards/master-detail while keyboard behavior stays table-specific.
- Transient dialog/focus state cannot live in Goja or URL state. The design introduces renderer services while preserving API-free presentational components.
- “Opinionated but flexible” required a layered escape strategy: domain views first, generic grammar second, typed engines third, and raw only last.

### What warrants a second pair of eyes

- Review the proposed final namespace names before code makes them expensive to change.
- Review whether `style` should become a real namespace now or remain scoped under domains/data.
- Review whether all new complex views need dedicated typed Go specs immediately or can temporarily normalize into shared specs.
- Review the 30–45 engineer-day estimate and phase boundaries against staffing.

### What should be done in the future

- Implement phases in order, beginning with descriptor parity and baseline tests.
- Record each vertical slice in this diary with exact commands and failures.
- Do not delete legacy provider modules until migration scans and browser tests are clean.

### Code review instructions

- Read the design in this order: Executive Summary, Intern Orientation, Language Design Principles, Target Public API, Decision Records, Implementation Plan.
- Use the parity inventory alongside the design when deciding whether to promote or hide a React adapter.
- Compare every API sketch against current `v3.go`, typed specs, declarations, and component contracts.
- Validate Mermaid and PDF rendering before publication.

### Technical details

- Primary guide: `design-doc/01-widget-dsl-v3-full-feature-analysis-design-and-intern-implementation-guide.md`.
- Parity reference: `reference/01-legacy-to-v3-feature-parity-inventory.md`.
- Proposed implementation: nine phases, approximately 30–45 engineer days.
- Hard-cutover exit: one provider module, no legacy imports/adapters, no unexplained raw uses, complete runtime/declaration/docs/example/browser parity.

## Step 3: Validate and Deliver the Research Bundle

This step validated both the documentation workspace and the current code baseline, then delivered a single PDF bundle to reMarkable. No production source files were changed; the baseline tests demonstrate that the design was written against a currently passing runtime and React package rather than an already-broken branch.

The upload used a dry run first and a new remote directory/name, so it did not overwrite an earlier document or risk annotations. The remote listing confirmed the bundle after upload.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Validate the research package and publish it to the requested reading device.

**Inferred user intent:** Make the design immediately reviewable away from the development machine and leave evidence that the ticket is internally consistent.

### What I did

- Validated frontmatter for index, design, parity inventory, and diary.
- Ran targeted Go tests for Widget DSL, migration checker, typed specs, and provider.
- Ran the React package TypeScript check.
- Related key implementation/evidence files to focused documents.
- Updated ticket tasks and changelog.
- Ran `docmgr doctor` successfully.
- Verified remarquee status/account.
- Dry-ran, uploaded, and remotely listed the bundled PDF.

### Why

- Ticket metadata and related-file integrity are part of the deliverable.
- Targeted baseline tests catch incorrect claims about current APIs.
- Dry-run plus remote listing makes delivery reproducible and non-destructive.

### What worked

- `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1` passed.
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `docmgr doctor --ticket WIDGETDSL-V3-FULL-FEATURE-CUTOVER --stale-after 30` reported all checks passed.
- `remarquee status` returned `remarquee: ok`.
- Account verification returned `user=wesen@ruinwesen.com sync_version=1.5`.
- Initial upload returned `OK: uploaded Widget DSL v3 Full Feature Hard Cutover.pdf -> /ai/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER`.
- After recording delivery in this diary, the complete final bundle was uploaded non-destructively as `Widget DSL v3 Full Feature Hard Cutover v2`.

### What didn't work

- N/A. The validation and upload commands completed successfully.

### What I learned

- The current Widget DSL/provider baseline is green despite the documented API parity gaps.
- The 55 KB primary guide and supporting references render successfully through the bundle pipeline.
- Keeping generated inventory in the PDF makes the design's numerical claims independently reviewable.

### What was tricky to build

- The bundle contains Mermaid code blocks. The default PDF pipeline preserves them as readable source rather than rendered diagrams; the prose and labels therefore remain sufficient without requiring a Mermaid rendering dependency.

### What warrants a second pair of eyes

- Review the PDF on-device for page breaks around wide API tables and Mermaid blocks.
- Confirm whether future revisions should use a new bundle name or intentionally replace this version.

### What should be done in the future

- When implementation starts, add code/test commits and phase-specific diary steps rather than rewriting this research history.
- Upload revised bundles under versioned names unless overwrite is explicitly requested.

### Code review instructions

- Run `docmgr doctor --ticket WIDGETDSL-V3-FULL-FEATURE-CUTOVER --stale-after 30`.
- Run the targeted Go and TypeScript commands listed above.
- Verify the remote path with:

  ```bash
  remarquee cloud ls /ai/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER --long --non-interactive
  ```

### Technical details

- Remote directory: `/ai/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER`.
- Final bundle: `Widget DSL v3 Full Feature Hard Cutover v2.pdf`.
- Bundle contents: ticket index, primary guide, parity inventory, diary, and generated runtime inventory.

## Step 4: Convert the Research Ticket into an Implementation Work Queue

Implementation began by preserving the validated research package as its own baseline commit and translating the nine-phase design into trackable docmgr tasks. This keeps architectural evidence separate from production changes and gives each future commit a concrete task boundary.

The first engineering slice is descriptor parity rather than user-facing UI. Complete API inventory and runtime parity checks reduce the risk of adding more helpers while declarations, docs, and runtime exports continue to drift.

### Prompt Context

**User prompt (verbatim):** "Implement task by task, commit at aprpopriate intervals, keep a detailed diary as you go."

**Assistant interpretation:** Execute the complete hard-cutover plan incrementally, use focused commits after passing tests, and maintain the ticket diary/changelog/tasks throughout implementation.

**Inferred user intent:** Obtain a reviewable implementation history where each architectural phase can be validated, reverted, or continued independently.

### What I did

- Re-read the ticket diary and current implementation plan before editing code.
- Checked branch, status, recent history, ticket tasks, and docmgr health.
- Confirmed the repository is on `main` and the ticket directory is the only untracked work.
- Added twelve implementation tasks: research baseline, three Phase 1 slices, Phases 2–8, and final release validation.
- Selected direct namespace descriptor parity as the first production-code slice.

### Why

- Committing the research baseline independently avoids mixing 11,000 words of design evidence with runtime code.
- Smaller Phase 1 tasks make descriptor modeling, builder composition, and package reorganization separately reviewable.
- The task list now mirrors the design guide and can drive diary/changelog updates.

### What worked

- `docmgr doctor --ticket WIDGETDSL-V3-FULL-FEATURE-CUTOVER --stale-after 30` still passed before implementation.
- Git status showed no tracked modifications that could be accidentally included.
- Existing descriptor tests identify a narrow starting seam.

### What didn't work

- N/A. Planning and task creation completed successfully.

### What I learned

- The ticket had seven completed research tasks but no implementation tasks, so reopening work required an explicit second task set.
- The current descriptor captures all top-level namespace names but only eight semantic view methods; direct runtime namespace exports are not parity-tested.

### What was tricky to build

- A 30–45 day design cannot safely map to one task per broad phase alone. Phase 1 was split into direct namespace parity, nested builders/contexts, and generation/package reorganization so commits remain coherent.

### What warrants a second pair of eyes

- Confirm that the implementation-task granularity remains useful as later phases reveal smaller vertical slices.
- Review commits on `main` carefully because no feature branch was requested or present.

### What should be done in the future

- Add subtasks when a broad phase exceeds one independently testable vertical slice.
- Check one task only after tests and the associated code commit exist.

### Code review instructions

- Review `tasks.md` to see the implementation sequence.
- Compare task wording to the nine phases in the primary design guide.
- Verify `git status --short --untracked-files=all -- ttmp/2026/07/12` contains only this ticket before the baseline commit.

### Technical details

- Branch: `main`.
- First implementation task: `dmpg`, direct namespace descriptor/runtime parity.
- Production baseline command: `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1`.
