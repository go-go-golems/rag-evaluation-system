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
    - Path: repo://pkg/widgetdsl/v3_descriptors_test.go
      Note: Direct runtime and descriptor set-equality enforcement from commit f208624
    - Path: repo://ttmp/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER--widget-dsl-v3-full-feature-hard-cutover/design-doc/01-widget-dsl-v3-full-feature-analysis-design-and-intern-implementation-guide.md
      Note: Design synthesis produced by the investigation
    - Path: repo://ttmp/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER--widget-dsl-v3-full-feature-hard-cutover/scripts/01-inventory-widget-dsl.py
      Note: Repeatable evidence inventory script created during investigation
    - Path: repo://ttmp/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER--widget-dsl-v3-full-feature-hard-cutover/scripts/02-regenerate-v3-api-reference.go
      Note: Reproducible generated-prefix help updater from commit f208624
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

## Step 5: Enforce Direct Runtime and Descriptor Parity

The first production slice replaces the partial namespace list with a module descriptor that accounts for the root `page` export and every direct member of every namespace. A runtime parity test now boots the actual Goja installer and fails if runtime keys and descriptor keys differ in either direction.

The same descriptor now emits the root `page` TypeScript declaration and a substantially more useful API reference. A checked-in regeneration tool preserves the hand-authored help suffix while replacing only the descriptor-generated prefix.

### Prompt Context

**User prompt (verbatim):** (same as Step 4)

**Assistant interpretation:** Implement the first planned task as a tested, focused commit and document both success and failure evidence.

**Inferred user intent:** Establish trustworthy API inventory before broadening or reorganizing the unreleased DSL.

**Commit (code):** `f208624f137d9367bda7796da3e50a04a5c8be2a` — "widgetdsl: enforce direct v3 API descriptor parity"

### What I did

- Introduced `v3ModuleDescriptor`, root export descriptors, and direct namespace member descriptors.
- Described all current root/namespace keys, including empty reserved `style`.
- Moved the `page(...)` TypeScript declaration from a handwritten line to descriptor output.
- Added runtime parity tests using `installWidgetV3` and Goja object keys.
- Added a test requiring every semantic view to be a described direct member.
- Expanded generated API help with root exports and direct member inventories.
- Added `scripts/02-regenerate-v3-api-reference.go` and regenerated embedded help.
- Ran targeted tests, package/provider tests, full repository tests, and pre-commit hooks.

### Why

- Runtime, TypeScript, descriptors, help, and examples had drifted because no test compared actual installed keys with the public inventory.
- Modeling direct exports first creates a stable base for nested builders and intent/context contracts in Phase 1B.
- Generated help must be reproducible without overwriting its authored troubleshooting and cross-reference sections.

### What worked

- Targeted descriptor/help tests passed after regeneration.
- `go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1` passed.
- The commit hook passed its scoped package tests, golangci-lint, and glazed-lint checks.
- Runtime parity now checks unexpected runtime additions as well as missing described exports.

### What didn't work

- The first targeted run failed as expected because the embedded help snapshot was stale:

  ```text
  --- FAIL: TestWidgetV3EmbeddedAPIHelpMatchesDescriptorReference (0.00s)
      v3_descriptors_test.go:123: embedded API help descriptor reference is stale; regenerate ../xgoja/providers/widgetsite/doc/05-widget-dsl-v3-api-reference.md from WidgetV3APIReferenceMarkdown
  FAIL
  ```

  I added and ran the scoped regeneration tool, then reran the tests successfully.

- `go test ./... -count=1` is not currently a clean repository command. It failed in unrelated pre-existing surfaces:

  ```text
  ttmp/2026/06/04/XGOJA-WIDGETSITE--xgoja-widget-site-binary-design/scripts/01-current-xgoja-widgetsite-experiment/widgetprovider/provider.go:18:12: not enough arguments in call to widgetdsl.NewLoader
  have ()
  want (string)
  ```

  and:

  ```text
  validate_test.go:20: unexpected manifest error: {Severity:error Check:unknown_module Path:packages/rag-evaluation-site/src/components/atoms/ContentStatusBadge/ContentStatusBadge.widget.yaml Subject:cms.dsl Message:widget module is not declared in schema/dsl-modules.yaml}
  ```

  These failures were not caused by this slice. The repository's pre-commit command excludes historical ticket scripts and passed.

### What I learned

- The direct public surface is 13 root keys: `page` plus 12 namespaces.
- Current namespace kinds are not uniformly objects: `data.selection` is a callable function with a nested helper, while `data.cell`, domain intents, and `time.range` are objects. Nested parity belongs in Phase 1B.
- Help generation needs an explicit authored/generated boundary to avoid discarding prose.

### What was tricky to build

- A set-equality test is stronger than checking only described exports. Both runtime and descriptor names are sorted before comparison so source declaration order does not affect correctness.
- Semantic views are a documented subset of direct members. The generator builds a view lookup so detailed view lines replace, rather than duplicate, the generic member line.

### What warrants a second pair of eyes

- Review whether `v3MemberDescriptor.Kind` should become a typed enum before adding nested descriptors.
- Review the help generator's authored suffix marker (`## Using this reference`) as a long-term generation boundary.
- Confirm that direct property enumerability through `goja.Object.Keys()` matches the intended public-JavaScript definition.

### What should be done in the future

- Phase 1B should model nested intent/range/cell/selection members and builder methods.
- Add `.use(fragment)` consistently to all mutable builders and parity-test those builder keys.
- Decide whether the two unrelated full-repository test failures should receive a separate maintenance ticket.

### Code review instructions

- Start in `pkg/widgetdsl/v3_descriptors.go` at `widgetV3Module`.
- Review `TestWidgetV3DescriptorMatchesDirectRuntimeExports` for the enforcement contract.
- Run:

  ```bash
  go run ./ttmp/2026/07/12/WIDGETDSL-V3-FULL-FEATURE-CUTOVER--widget-dsl-v3-full-feature-hard-cutover/scripts/02-regenerate-v3-api-reference.go
  go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1
  ```

### Technical details

- Direct descriptor scope deliberately excludes nested builder methods and action-context schemas until Phase 1B.
- `widgetV3DescriptorTypeScriptLines()` now emits root function declarations and namespace constants.
- Generated help remains embedded by the existing widgetsite provider.

## Step 6: Describe Nested APIs and Make Every V3 Builder Composable

This slice extends parity below direct namespace keys. Nested cell/selection/range/intent namespaces and all 28 public builder types now have descriptor inventories, and runtime probes instantiate every builder through the public JavaScript API to compare its exact method set.

A shared builder constructor adds `.use(fragment)` uniformly. TypeScript builder interfaces inherit a generic `ComposableBuilder<T>`, so reusable policy fragments now work across pages, collections, domain views, fields, matrices, and actions rather than only pages and sections.

### Prompt Context

**User prompt (verbatim):** (same as Step 4)

**Assistant interpretation:** Continue Phase 1 task-by-task, commit the builder-composition behavior separately, then complete action-context documentation as a focused follow-up.

**Inferred user intent:** Make composition a language invariant and ensure nested public APIs are as discoverable and drift-resistant as root namespaces.

**Commit (code):** `7d78d447dc543f76c66094f73db9983984e28bed` — "widgetdsl: compose and describe all v3 builders"

**Commit (code):** `1a387adbe3ae82d11a25861892b8631af913c277` — "widgetdsl: document v3 action contexts"

### What I did

- Added descriptors for nine nested namespaces: data cells/selection, domain intents, and time ranges.
- Added descriptors for 28 public builder types and their exact methods.
- Added `newV3Builder(path)` and migrated every mutable v3 builder to it.
- Added `.use(fragment)` to field, collection, table, editor, matrix, UI action, CMS, course, context, schedule, time, and CRM builders.
- Added `ComposableBuilder<T>` to generated TypeScript declarations and made every v3 builder interface extend it.
- Added a public-API runtime probe that instantiates all builders and compares exact method sets.
- Added a fragment-composition behavior test.
- Documented 28 currently dispatched browser action contexts from collection, matrix, context, course, CMS, time, CRM, activity, and upload adapters.
- Added uniqueness/component-field validation for action-context descriptors.
- Regenerated embedded API help after both commits.
- Added follow-up task `tqve` for builder actions that source inspection revealed are serialized by Go but ignored by current React adapters.

### Why

- `.use(fragment)` was an explicit v3 grammar rule but existed only on page and section builders.
- Exact builder probes prevent TypeScript and runtime method sets from diverging silently.
- Action payload bindings are only safe when authors can discover which context fields adapters actually dispatch.

### What worked

- All 28 builder probes execute through the real public `widget` object.
- Nested namespace set equality passed for callable objects such as `data.selection` as well as ordinary namespace objects.
- A collection and nested table successfully composed two reusable fragments through `.use(...)`.
- Targeted Widget DSL/provider tests and both pre-commit test/lint suites passed for both commits.

### What didn't work

- The first formatting run failed because a TypeScript declaration edit accidentally joined two Go string literals:

  ```text
  pkg/widgetdsl/typescript.go:354:207: missing ',' in composite literal
  pkg/widgetdsl/typescript.go:354:419: string literal not terminated
  ```

  I separated `CrmActivityFeedBuilder` and `CrmIntentNamespace` into valid list entries and reran formatting.

- The first builder probe passed `pipeline.build()` and `crmFields.build()` into helpers that currently require hidden Goja builder references. It panicked in `crmRef` with a nil-pointer dereference. The probe was corrected to pass the builder handles directly. This exposes a declaration/runtime mismatch to revisit in Phase 1C or CRM completion.

- The first fragment assertion cast a typed `spec.JSONObject` export to `map[string]any`, then a Goja object lookup assumed a direct `props` shape. Both were brittle. The test now asserts observable fragment invocation count, while exact builder keys cover structural behavior.

### What I learned

- `.use` can be implemented once without changing lowering: it is author-time callback composition over the same builder object.
- Several public builder methods currently set props that their React adapters ignore: context diagram selection, Markdown editor actions, CMS publish/archive/preview actions, and course material deletion are notable examples.
- `CrmNamespace` TypeScript currently says built records are accepted where runtime relies on hidden builder references; the hard cutover can fix this contract rather than preserve it.
- Browser action context is defined by adapters, not by Go action builders, so adapter evidence must be part of API documentation.

### What was tricky to build

- Some builders are returned handles, while others exist only inside callbacks. The parity test probes all of them through realistic public calls instead of coupling to private Go constructors.
- `data.selection` is both callable and property-bearing. Resolving nested descriptor paths through `goja.Value.ToObject` handles functions and objects uniformly.
- Adding composition to inline builders required finding course material and Markdown editor builders that did not have named Go constructor methods.

### What warrants a second pair of eyes

- Review the 28 action-context descriptors against adapter source, especially CMS callbacks with shared dispatch helpers.
- Review whether all builders—including low-level ActionsBuilder—should remain public composable types.
- Review the CRM built-value mismatch before declarations are made fully descriptor-generated.
- Inspect task `tqve`; ignored action props are correctness gaps, not documentation-only debt.

### What should be done in the future

- Implement task `tqve` in the appropriate UI/domain phases and add browser interaction tests.
- Phase 1C should make builder declarations themselves descriptor-generated or parity-checked beyond the composition base interface.
- Add typed context interfaces generated from action-context descriptors after final field names stabilize.

### Code review instructions

- Start with `newV3Builder` in `pkg/widgetdsl/v3.go`, then inspect migrated builder constructors.
- Review `TestWidgetV3DescriptorMatchesBuilderRuntimeMethods` and its JavaScript probe.
- Compare `widgetV3Module.NestedNamespaces`, `Builders`, and `ActionContexts` to the generated help.
- Run:

  ```bash
  go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1
  ```

### Technical details

- Every described builder automatically receives `use`, and descriptor method lists append it through `v3Builder(...)`.
- `ComposableBuilder<TBuilder>` returns `TBuilder`, preserving fluent TypeScript inference.
- The runtime probe fails on both missing descriptor methods and unexpected JavaScript methods.

## Step 7: Turn Baseline Failures and Ignored Actions into Explicit Blockers

Before continuing declaration work, this step converted every known baseline failure and ignored-action category into explicit ticket tasks. The two existing `go test ./...` failures were then fixed immediately in separate commits, restoring the full repository test command as a meaningful release gate.

The action-adapter problem remains intentionally open as a parent audit plus component-focused wiring and CI-enforcement tasks. This prevents the implementation from declaring parity merely because Go successfully serializes an action prop that React never consumes.

### Prompt Context

**User prompt (verbatim):**

> how are we tackling this?  Issues: Full-repository tests expose two pre-existing failures, and several serialized actions are
>  ignored by React adapters.
>
> Make sure we have proper tasks to address that in the ticket.

**Assistant interpretation:** Explain the remediation strategy, add concrete ticket tasks for every known failure category, and ensure they are treated as release blockers rather than diary caveats.

**Inferred user intent:** Prevent known cross-layer correctness gaps and broken validation commands from being deferred or forgotten during the large v3 migration.

**Commit (code):** `2017908b6fa4cb151fe73564830c3912c5abb5bc` — "widgetdsl: update historical provider loader calls"

**Commit (code):** `217ad1399cc903fccb53843c2d0cd286f6054851` — "widgets: repair manifest catalog validation"

### What I did

- Added explicit tasks for both full-repository failures.
- Split ignored actions into an audit task, two component-wiring tasks, and a CI-enforcement task while retaining parent task `tqve`.
- Updated the historical experimental provider to call `NewLoader(widgetdsl.WidgetV3ModuleName)` for both exported aliases.
- Added transitional `cms.dsl` and `time.dsl` module ownership to the manifest catalog with hard-cutover intent in their descriptions.
- Migrated three stale manifests from `entry` to `adapter`, added schema versions, and completed missing context manifest metadata.
- Added the organisms folder to transitional `ui.dsl` ownership because `FormPanel` is registered there today.
- Reran manifest validation after every newly exposed error.
- Restored and verified `go test ./... -count=1`.

### Why

- Historical ticket scripts are Go packages under the module and therefore part of unscoped repository tests unless explicitly excluded.
- Widget manifests describe the current adapter catalog; their module catalog must remain internally valid during migration even though Phase 8 will collapse legacy ownership into `widget.dsl`.
- Ignored action props require cross-layer tests, not only individual bug fixes.

### What worked

- The historical provider package now compiles independently.
- `go test ./internal/widgetmanifest -count=1` passes.
- `go test ./... -count=1` now passes across production packages, historical scripts, ticket tooling, Widget DSL, and provider packages.
- The repairs remained in two focused code commits.

### What didn't work

- Adding `cms.dsl` exposed a stale manifest key:

  ```text
  unexpected manifest error: {Severity:error Check:required_field Path:packages/rag-evaluation-site/src/components/atoms/TextareaInput/TextareaInput.widget.yaml Subject:adapter.export Message:required widget manifest field is empty}
  ```

  Three manifests used obsolete `entry:` blocks. They were migrated to `adapter:`.

- The next validation run exposed missing schema versions, followed by missing helper/module/status/docs metadata in two context manifests.

- Once those were fixed, validation exposed undeclared `time.dsl`, then an incomplete `ui.dsl` folder root for `FormPanel`. Each catalog inconsistency was corrected before rerunning the full suite.

### What I learned

- The original two failures masked several additional manifest inconsistencies because the test stops on the first error.
- Manifest module ownership is still legacy-oriented. Transitional entries need explicit wording so they are not mistaken for the final hard-cutover architecture.
- The unscoped repository test is valuable precisely because it compiles tracked ticket scripts that can otherwise rot unnoticed.

### What was tricky to build

- Fixing only the first manifest error would have produced repeated CI failures. I compared all manifest module names against the schema and all `entry` versus `adapter` keys to repair the class of drift, then used the validator to expose remaining metadata/folder issues.
- The historical `rag.dsl` alias cannot be passed to `NewLoader` because it is not a registered widget module name. Both alias factories must use the real `widget.dsl` loader while the provider exposes different external names.

### What warrants a second pair of eyes

- Review whether historical ticket scripts should remain in `go test ./...` permanently or move behind a dedicated experiment module/build tag.
- Review transitional `cms.dsl` and `time.dsl` ownership before Phase 8 removes split-module identities.
- Review the action-audit task before wiring individual props; it should identify every ignored typed action, not only the four already observed.

### What should be done in the future

- Complete tasks `pi3h`, `pcgq`, `awor`, and `u11p`, then close parent task `tqve`.
- In Phase 8, replace transitional manifest module ownership with the final `widget.dsl` model rather than preserving these entries.
- Keep `go test ./... -count=1` as a Phase 9 hard release gate.

### Code review instructions

- Review `schema/dsl-modules.yaml` together with all component `*.widget.yaml` module values.
- Review the historical provider change to confirm both aliases load `WidgetV3ModuleName`.
- Validate with:

  ```bash
  go test ./internal/widgetmanifest -count=1
  go test ./... -count=1
  ```

### Technical details

- Closed tasks: `4sqg` and `sx1t`.
- Open ignored-action tasks: `tqve`, `pi3h`, `pcgq`, `awor`, and `u11p`.
- The catalog remains transitional until the Phase 8 provider/manifest hard cutover.

## Step 8: Finish the Descriptor and Typed-Spec Foundation

Phase 1C removed the obsolete `v2/spec` package identity and made TypeScript namespace/builder declarations exact-parity checked against descriptors. The typed kernel now has a release-neutral import path, while legacy data.v2 builders and v3 collections share it without a compatibility alias.

Declaration tests parse top-level TypeScript interface members, compare them as sets to direct, nested, and builder descriptors, and continue to compare runtime objects separately. This creates a three-way invariant across Goja runtime, descriptors, and TypeScript declarations.

### Prompt Context

**User prompt (verbatim):** "Continue. I'm going swimming, I want you to be done when I come back."

**Assistant interpretation:** Continue autonomously through the implementation queue, testing and committing each dependency-ordered slice without waiting for intermediate confirmation.

**Inferred user intent:** Return to a substantially completed, clean, reviewable implementation rather than another planning checkpoint.

**Commit (code):** `40703641d1e6a4f0cc22d76ed28157261b9aebb2` — "widgetdsl: promote shared typed spec kernel"

**Commit (code):** `4b61e61c27e261b20e1d921d092bc110e5e8f67f` — "widgetdsl: enforce descriptor declaration parity"

### What I did

- Moved `pkg/widgetdsl/v2/spec` to `pkg/widgetdsl/spec` with Git history preserved.
- Updated v2 legacy builders and v3 runtime imports to the neutral `widgetspec` alias.
- Added named TypeScript declarations for selection, CMS intents, and course intents.
- Added descriptor type names for all nested namespaces.
- Added a brace-aware TypeScript interface parser in tests.
- Enforced exact member parity for direct namespaces, nested namespaces, and all builder declarations.
- Ran targeted and full repository tests before both commits.

### Why

- Shipping v3 on a package named `v2/spec` would preserve misleading architecture and complicate future ownership.
- Runtime parity alone cannot prevent declarations from advertising missing methods or omitting installed methods.
- Named nested interfaces improve discoverability and enable descriptor parity checks.

### What worked

- `go test ./... -count=1` passed after the package move and declaration changes.
- Git recognized all six spec files as renames.
- The declaration parser handles one-line and multiline interfaces plus nested object types.
- Pre-commit tests and lint passed for both commits.

### What didn't work

- An existing TypeScript snapshot test expected the old inline selection intersection. It failed after introducing `SelectionNamespace`; the fixture was updated to assert the named interface contract.

### What I learned

- Exact method-name parity can be enforced without generating every complex TypeScript signature from metadata.
- CMS and course intent declarations previously used broad index signatures, masking runtime drift; explicit methods now match Goja objects.

### What was tricky to build

- Interface signatures contain nested object types, so splitting blindly on braces or semicolons produces false members. The parser tracks brace depth and only flushes declarations at the outer interface level.

### What warrants a second pair of eyes

- Review whether signature-level descriptor generation is worth the metadata duplication beyond exact member parity.
- Verify downstream repositories do not import the old internal Go package path.

### What should be done in the future

- Add signature metadata only where generated docs need parameter-level details.
- Continue with generic UI parity and action-adapter auditing.

### Code review instructions

- Review the package rename, then `TestWidgetV3TypeScriptNamespaceMembersMatchDescriptors` and `typescriptInterfaceMembers`.
- Validate with `go test ./... -count=1`.

### Technical details

- No compatibility package remains at `pkg/widgetdsl/v2/spec`.
- Runtime, descriptor, and declaration member sets are independently compared.

## Step 9: Wire Every Serialized V3 Action into React

This slice audited all action props assigned by the v3 Go runtime and established CI enforcement that each is declared in Widget IR and consumed by a React adapter. The known context, Markdown, CMS, and course gaps were fixed in the same vertical commit.

Context diagrams now dispatch selected part context, Markdown editing emits change and explicit Save actions, CMS article rows dispatch dedicated publish/archive/preview actions, and course material uploads can render deletable items with stable asset context.

### Prompt Context

**User prompt (verbatim):** (same as Step 8)

**Assistant interpretation:** Continue through release blockers without pausing, implementing cross-layer action behavior and tests.

**Inferred user intent:** Eliminate silent no-op APIs before adding more DSL surface.

**Commit (code):** `8e6c83178012c65c434056d345fb17a13b95f64b` — "widgets: wire all serialized v3 action props"

### What I did

- Added `onPartSelect` behavior to ContextDiagramPanel and adapter dispatch context.
- Added Markdown `onChangeAction` and `onSubmitAction` IR props and adapter behavior.
- Corrected Markdown lowering from ignored `value` to `defaultValue`.
- Added a Save toolbar action carrying current Markdown value.
- Added dedicated article publish/archive/preview IR props and adapter routing.
- Added a visible Preview row action.
- Added upload item/delete contracts and course material delete dispatch context.
- Added `TestWidgetV3SerializedActionPropsHaveIRAndAdapterConsumers`, auditing all action props extracted from v3 Go sources against IR declarations and adapter consumption.
- Updated the Markdown golden fixture and existing CMS test.

### Why

- Serializing an action prop without adapter consumption creates an API that looks valid in Goja and goldens but never executes in the browser.
- A source-driven audit prevents recurrence as new builders add actions.

### What worked

- The audit found 35 serialized action props and all now have IR declarations and adapter consumers.
- React typecheck, Biome, Widget DSL tests, full Go tests, and pre-commit checks passed.
- The final commit remained one cross-layer behavior slice.

### What didn't work

- Typecheck initially rejected unsupported `Inline align` and `Button variant=\"danger\"` props. The item row now uses `style={{ alignItems: \"center\" }}` and the supported default button variant.
- The first audit only read `props.ts`, but several engine action contracts live in `engines.ts`. It was corrected to scan every TypeScript file under `widgets/ir`.
- Existing tests and a golden expected Markdown `value`; they were updated after confirming the React adapter contract uses `defaultValue`.

### What I learned

- Several cross-layer action bugs were accompanied by prop-name drift, not only missing callbacks.
- A global adapter-consumption check is a useful floor; component interaction stories remain necessary to prove exact payload behavior.

### What was tricky to build

- CMS exposes both a generic row action and dedicated semantic actions. The adapter now selects the dedicated action first and falls back to the generic action.
- Course deletion required a real interactive representation; merely forwarding `onDeleteAction` would still leave no browser affordance.

### What warrants a second pair of eyes

- Review whether Markdown change actions should dispatch per keystroke or eventually support debounce/blur policy.
- Review upload deletion presentation and confirmation policy.
- Review the source-audit test for false positives when two components someday reuse an action prop name.

### What should be done in the future

- Add Storybook interaction tests for exact action contexts.
- Extend audit metadata to component-specific ownership if prop-name reuse becomes ambiguous.

### Code review instructions

- Begin with `pkg/widgetdsl/v3_action_adapter_test.go`.
- Trace each repaired action from `v3.go`, through `props.ts`, into its `.widget.tsx` adapter.
- Run `go test ./... -count=1` and `pnpm --dir packages/rag-evaluation-site typecheck`.

### Technical details

- Action audit threshold intentionally fails if extraction unexpectedly falls below 30 props.
- Markdown submit context is `{ value, componentType: \"MarkdownEditor\" }`.
- Upload delete context includes `assetId`, `asset`, `value`, and `componentType`.

## Step 10: Add the Generic Content Vocabulary

Phase 2 added typed v3 access to stable generic content and layout capabilities that previously required raw component construction. The helpers remain intent-oriented and reuse existing registered adapters rather than duplicating React behavior.

### Prompt Context

**User prompt (verbatim):** (same as Step 8)

**Assistant interpretation:** Continue filling approved generic parity after closing action correctness blockers.

**Inferred user intent:** Make ordinary application pages expressible without `widget.raw.component(...)`.

**Commit (code):** `e702cfcfed73144f4058c56acfc1452f867884c9` — "widgetdsl: add typed generic content helpers"

### What I did

- Added `ui.text`, `code`, `divider`, `scroll`, `tabs`, `summary`, `checkList`, `stepList`, `markdownArticle`, and `upload`.
- Added descriptor and TypeScript parity for all helpers.
- Added lowering tests covering every component target and Markdown source.
- Regenerated embedded API help.

### Why

- These are stable, domain-neutral concepts repeatedly used by first-party pages and raw v3 examples.

### What worked

- Descriptor/runtime/declaration parity tests passed immediately after all three surfaces were updated.
- Full Go and pre-commit test/lint suites passed.

### What didn't work

- N/A.

### What I learned

- Existing Widget adapters already cover most generic parity; the missing layer was a coherent v3 vocabulary.

### What was tricky to build

- Helpers accepting children use the existing component-factory normalization, while data-shaped helpers explicitly place ordered arrays into props.

### What warrants a second pair of eyes

- Review final naming of `summary`, `checkList`, and `upload` before public release.

### What should be done in the future

- Rewrite remaining raw examples to these typed helpers during migration.

### Code review instructions

- Review `v3UIObject`, `UINamespace`, and `TestWidgetV3GenericContentHelpers` together.
- Run `go test ./... -count=1`.

### Technical details

- All helpers lower to existing registry component names; no parallel React implementations were introduced.

## Step 11: Remove Raw Escapes from the Canonical V3 Corpus

The canonical v3 examples now use typed generic and domain APIs exclusively. A missing CRM single-field helper was added rather than leaving FieldRenderer as the final raw exception.

### Prompt Context

**User prompt (verbatim):** (same as Step 8)

**Assistant interpretation:** Continue migration until the parser-backed checker can enforce zero raw escapes in first-party v3 examples.

**Inferred user intent:** Make examples trustworthy release documentation for the final language.

**Commit (code):** `8c3c1d3618e46aac2ee78a05b7be6dcd1902da7a` — "widgetdsl: eliminate raw escapes from v3 examples"

### What I did

- Rewrote Markdown, CRM record, CRM board, activity, form, upload, and field-renderer examples to typed APIs.
- Added `crm.field(value, spec, options)` for standalone FieldRenderer use.
- Regenerated all affected goldens and API help.
- Ran the migration checker with `--fail-on-findings`.

### Why

- Raw examples normalize bypasses and hide missing public grammar.

### What worked

- Migration checker reports: `No legacy Widget DSL imports or raw component escape hatches found.`
- Full Go tests and pre-commit checks passed.

### What didn't work

- Biome reports existing example `page` variables as unused because the Goja example runner reads them by convention. These remain warnings, not failures.

### What I learned

- The only genuine remaining raw gap in the corpus was standalone CRM field rendering.

### What was tricky to build

- Golden updates had to follow semantic helper rewrites while preserving stable rendered IR.

### What warrants a second pair of eyes

- Review whether example-runner conventions should export page explicitly to eliminate Biome warnings.

### What should be done in the future

- Apply the same fail-on-findings gate workspace-wide before deleting legacy modules.

### Code review instructions

- Run the migration checker command from the design guide with `--fail-on-findings`.

### Technical details

- Raw findings fell from 11 to zero in the v3 golden corpus.

## Step 12: Deliver Interactive Collections, FormDialog, and Generic Activity

Phases 3–6 landed as one end-to-end interaction slice rather than disconnected APIs. The typed spec now shapes search and pagination, DataTable owns accessible keyboard mechanics, and FormDialog carries frozen row context through native form serialization and server validation.

### Prompt Context

**User prompt (verbatim):** "continue

do phases 3-8, i'm going swimming, so i want to be coming back to something that kicks ass"

**Assistant interpretation:** Autonomously implement and validate the queued interaction phases through hard cutover.

**Inferred user intent:** Return to a substantial, integrated implementation rather than another design-only checkpoint.

**Commit (code):** `ed372286eeebf89b7ac9b6d9e09848fec7fe69de` — "widgetdsl: add interactive collections and overlay forms"

### What I did

- Added structured navigation query/preserve/omit/replace semantics and overlay actions.
- Added FormDialog with opener restoration, native validation, `context.form`, field errors, and live-region notifications.
- Added collection search/pagination builders, page-size selection, stable lowering order, and typed specs.
- Added DataTable roving focus, Arrow/J/K navigation, Enter selection, row commands, editable-target guards, command help, and semantic tones.
- Promoted ActivityFeed to `data.activityFeed` and removed the CRM duplicate.
- Added example/golden 42 as the integrated interaction contract.

### Why

- The Upwork workflow needs URL-driven scale and keyboard-speed triage without custom DOM scripts.

### What worked

- Full Go tests, TypeScript, focused frontend checks, descriptor parity, action audit, and migration checks passed.

### What didn't work

- The first pre-commit attempt failed with `missing cases in switch of type spec.ActionKind: spec.ActionKindOpenOverlay, spec.ActionKindCloseOverlay` and unused CRM activity helpers.
- Biome initially rejected FormDialog's effect with `This hook does not specify its dependency on close`.

### What I learned

- Action result routing must exist both in the default dispatcher and the application-provided action handler.

### What was tricky to build

- Focus had to move by stable row key after React committed state; setting `tabIndex` alone does not move DOM focus.
- Structured navigation options had to survive typed Go action conversion instead of disappearing during lowering.

### What warrants a second pair of eyes

- Browser-test IME composition and nested editable controls around table commands.
- FormDialog field-error placement currently provides a summary rather than binding each message to its control.

### What should be done in the future

- Add Playwright coverage for repeated query values, Back/Forward, dialog validation, and S/R/T commands.

### Code review instructions

- Start with `pkg/widgetdsl/spec/types.go`, `pkg/widgetdsl/spec/lower.go`, `DataTable.tsx`, and `FormDialog.widget.tsx`.
- Run `go test ./... -count=1` and package `typecheck`/`test:focused`.

### Technical details

- The collection lowering order is search → arrangement → pagination → detail.
- Page-size changes dispatch `{page: 1, pageSize, value}`.

## Step 13: Hard-Cut First-Party Hosts to One Module

Phase 8 removed split modules from both production registration paths and converted the complete React adapter manifest catalog to `widget.dsl`. The xgoja widget-site example was rewritten around native v3 collections rather than receiving aliases.

### Prompt Context

**User prompt (verbatim):** (same as Step 12)

**Assistant interpretation:** Complete the first-party migration and make legacy module names unavailable in generated hosts.

**Inferred user intent:** Enforce one coherent language in production, examples, manifests, and documentation.

**Commit (code):** `a028a9ca5c053085ff6ad0c0f6d858c0536a4145` — "widgetdsl: hard cut over first-party hosts to widget.dsl"

**Commit (code):** `b439d5de97b1e6243e76024f89010de7a9aba49f` — "widgets: route server validation results to dialogs"

### What I did

- Restricted `Register`, `Registrar`, the xgoja provider, and public `NewLoader` to `widget.dsl`.
- Added negative tests for every legacy module name.
- Converted 85 widget manifests/adapters and the module catalog to sole `widget.dsl` ownership.
- Rewrote the xgoja widget-site host and xgoja module plan to native v3.
- Replaced split-module help and README guidance.
- Ran the parser-backed first-party migration gate with zero findings.
- Routed application-owned server responses into FormDialog and toast services.

### Why

- A hard cutover is only real when providers, manifests, examples, and docs all reject the old dialects.

### What worked

- `go test ./... -count=1`, frontend typecheck/focused checks, manifest validation, and first-party migration checks pass.
- Production module resolution exposes exactly `widget.dsl`.

### What didn't work

- The first commit attempt surfaced Biome's `useIterableCallbackReturn` errors in the rewritten xgoja example; callbacks were reformatted by the hook and the validated commit completed.
- Existing adapter non-null assertions and CSS `!important` declarations remain Biome warnings inherited from the adapter corpus.

### What I learned

- Provider registration, engine registration, and xgoja build-plan selection are three independent exposure surfaces and all need negative tests.

### What was tricky to build

- Historical legacy implementation tests still need private registration to preserve regression archaeology while public loaders reject those names. The test-only function is deliberately unexported and unreachable from production registrars.

### What warrants a second pair of eyes

- Decide whether to delete the retained private legacy implementation tests/code now or in a dedicated repository-size cleanup after release evidence is archived.
- Browser-smoke the rewritten generated xgoja host against embedded assets.

### What should be done in the future

- Phase 9 should regenerate host assets, run Playwright/Storybook, and publish the final review bundle.

### Code review instructions

- Start with `pkg/widgetdsl/module.go`, `pkg/xgoja/providers/widgetsite/provider.go`, `schema/dsl-modules.yaml`, and `examples/xgoja/widget-site/verbs/sites.js`.
- Verify old module resolution fails and `widget.dsl` succeeds.

### Technical details

- The public migration scan reports no legacy imports or raw escapes under `examples`, `packages`, and the root README.
