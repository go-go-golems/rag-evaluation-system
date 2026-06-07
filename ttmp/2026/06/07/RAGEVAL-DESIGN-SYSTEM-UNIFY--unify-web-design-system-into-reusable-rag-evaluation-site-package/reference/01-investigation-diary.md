---
Title: Investigation Diary
Ticket: RAGEVAL-DESIGN-SYSTEM-UNIFY
Status: active
Topics:
    - design-system
    - frontend-architecture
    - react
    - packaging
    - storybook
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--unify-web-design-system-into-reusable-rag-evaluation-site-package/scripts/01_capture_css_visual_diff_help.sh
      Note: Step 1 css-visual-diff help capture
    - Path: ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--unify-web-design-system-into-reusable-rag-evaluation-site-package/scripts/02_component_duplication_inventory.py
      Note: Step 1 duplication inventory script
    - Path: web/src/App.tsx
      Note: Direct package imports after compatibility barrel removal
    - Path: web/src/components/index.ts
      Note: Barrel no longer exposes package-owned compatibility APIs
    - Path: web/src/components/pages/SearchWorkbenchPage/SearchWorkbenchPage.tsx
      Note: Keeps app-owned molecules local while package primitives import directly
ExternalSources: []
Summary: Chronological diary for unifying duplicated web/src design-system components into packages/rag-evaluation-site with Storybook and css-visual-diff visual parity checks.
LastUpdated: 2026-06-07T00:00:00Z
WhatFor: Use to resume or review the design-system unification work.
WhenToUse: Before changing frontend package boundaries or deleting duplicated files.
---



# Diary

## Goal

Capture the design-system unification work: evidence, scripts, failures, commits, visual parity checks, and migration steps from duplicated `web/src` components toward package-canonical imports from `packages/rag-evaluation-site`.

## Step 1: Ticket Setup, css-visual-diff Help Capture, and Duplication Inventory

I created a dedicated docmgr ticket for making `packages/rag-evaluation-site` the canonical home of the reusable RAG design system and Widget IR renderer. Before changing code, I captured `css-visual-diff` help pages into the ticket and generated a machine-readable duplication inventory between `web/src` and the package.

The inventory confirms the unification premise: the package is already a copied-out subset of the `web/` design system. Most shared primitive/layout/molecule files are byte-for-byte identical, so the safest migration is to make `web` depend on the local package and redirect imports before deleting duplicated files.

### Prompt Context

**User prompt (verbatim):** "ok, create a new ticket to do the unification / move the design system . Create a detailed analysis and guide, and a list of detailed granular tasks , then work on them one by one, keeping a diary as you go, usig css-visual-diff to compare for visual parity (see css-visual-diff help --all on how to do that, and potential scripts to write, which you should store in the scripts/ folder of the ticket like all scripts), committing at appropriate intervals, keeping a diary.

Use storybook for the visual parity, don't hesitate to add stories if missing. 

Have the page properly depend on the package for development."

**Assistant interpretation:** Create a new docmgr ticket, write an analysis/implementation guide and granular task list, then begin migrating duplicated design-system code from `web/src` into `packages/rag-evaluation-site` with Storybook/css-visual-diff visual parity checks and frequent diary updates.

**Inferred user intent:** Reduce frontend duplication while preserving visual behavior, making the reusable package the canonical design-system/Widget IR surface and keeping `web/` focused on app-specific pages and backend-connected containers.

### What I did
- Committed the prior context-viewer ticket and `web/README.md` first as `2343bef docs: capture context viewer design analysis` to start this new work from a clean tree.
- Created ticket `RAGEVAL-DESIGN-SYSTEM-UNIFY`.
- Added design doc `design-doc/01-design-system-unification-analysis-and-implementation-guide.md`.
- Added this diary `reference/01-investigation-diary.md`.
- Added initial granular tasks for css-visual-diff setup, duplication inventory, workspace dependency, import migration, Storybook coverage, and visual parity checks.
- Wrote and ran `scripts/01_capture_css_visual_diff_help.sh`.
- Wrote and ran `scripts/02_component_duplication_inventory.py`.

### Why
- The user explicitly asked for ticketed work, scripts in the ticket, css-visual-diff, Storybook, diary, and appropriate commits.
- The duplication inventory gives objective evidence before changing imports or deleting files.

### What worked
- `css-visual-diff help --all` worked and focused help pages were captured under `sources/css-visual-diff-help/`.
- The duplication script generated `sources/01-duplication-inventory.md` and `.json`.
- Inventory result:
  - `components`: 72 identical shared files, 2 differing shared files, 120 web-only files.
  - `widgets`: 3 identical shared files, 2 differing shared files, 1 web-only file.
  - `hooks`: package and web each have one non-overlapping hook file.

### What didn't work
- No blocker yet.

### What I learned
- The package has no package-only component files; it is a strict subset of the `web` component library.
- Most duplicate files are identical, so a package-canonical migration can be incremental and low-risk.
- The package `WidgetRenderer.tsx` and `actions.ts` differ from `web`, and the package versions include newer behavior from PR review work. These should be treated as the current canonical versions.

### What was tricky to build
- The important distinction is between reusable presentational files and Storybook/app files. Story files are currently web-only but may need to move or be duplicated temporarily to test package components visually.
- `web` currently has no root workspace dependency on the local package, so development wiring must come before import migration.

### What warrants a second pair of eyes
- Whether Storybook should live in `web`, in the package, or in both during the migration.
- Whether domain-specific organisms like corpus/workflow panels should move to the package now or only after primitive/layout/widget duplication is removed.

### What should be done in the future
- Add workspace wiring and aliases so `web` can import package source during development.
- Add css-visual-diff scripts that compare Storybook pages before and after migration.
- Migrate imports in `web` from local primitive/layout/molecule/widget paths to package imports.

### Code review instructions
- Start with `sources/01-duplication-inventory.md`.
- Review `scripts/01_capture_css_visual_diff_help.sh` and `scripts/02_component_duplication_inventory.py`.
- Confirm the initial task list in `tasks.md` before code migration.

### Technical details
- css-visual-diff help capture: `sources/css-visual-diff-help/`.
- Duplication inventory: `sources/01-duplication-inventory.md` and `sources/01-duplication-inventory.json`.

## Step 2: Add Workspace Development Dependency

I wired the repository as a two-package pnpm workspace so the `web` application can depend on the local `@go-go-golems/rag-evaluation-site` package during development. This is the prerequisite for moving imports from duplicated local files to package-canonical files.

I also added Vite and TypeScript source aliases so development resolves package imports to `packages/rag-evaluation-site/src` rather than stale built artifacts in `dist`.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Begin the migration by making `web` properly depend on the local reusable package for development.

**Inferred user intent:** Ensure unification is real package consumption, not another copy or alias-only trick that hides duplication.

### What I did
- Added root `pnpm-workspace.yaml` with `web` and `packages/rag-evaluation-site`.
- Added `@go-go-golems/rag-evaluation-site: workspace:*` to `web/package.json`.
- Added Vite aliases in `web/vite.config.ts` for package root, `/app`, and `/ir` imports.
- Added TypeScript `paths` in `web/tsconfig.json` for the same package entrypoints.
- Ran `pnpm install` at the repo root.
- Added `/node_modules/` to `.gitignore` because root workspace install creates root `node_modules`.
- Ran package and web typechecks.

### Why
- `web` needs a real local workspace dependency before import migration.
- Source aliases keep local development fast and avoid requiring package rebuilds for every component change.

### What worked
- `pnpm install` succeeded with `Scope: all 2 workspace projects`.
- `pnpm --dir web typecheck` passed.
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.

### What didn't work
- Root `node_modules/` appeared as untracked because the repo only ignored `web/node_modules/`. I fixed this by adding `/node_modules/` to `.gitignore`.

### What I learned
- The root workspace lockfile is now the coordination point for both frontend packages.
- No code imports have been redirected yet, so this step is low-risk and validates dependency wiring only.

### What was tricky to build
- Vite/TypeScript need source aliases because the package `exports` point at built package files for distribution, while development should consume source.

### What warrants a second pair of eyes
- Confirm whether package subpath aliases should use an array form in Vite if future prefix matching causes an issue.

### What should be done in the future
- Begin import migration in small batches, starting with app-level imports and duplicated widgets.

### Code review instructions
- Review `pnpm-workspace.yaml`, `web/package.json`, `web/vite.config.ts`, and `web/tsconfig.json`.
- Validate with `pnpm install`, `pnpm --dir web typecheck`, and `pnpm --dir packages/rag-evaluation-site typecheck`.

## Step 3: Redirect Shared Barrels and Add Storybook/css-visual-diff Parity Harness

I redirected the main `web/src` component and widget barrels so app code that imports `../atoms`, `../foundation`, `../layout`, shared `../molecules`, or `../widgets` now receives implementations from `@go-go-golems/rag-evaluation-site`. This is the first import-migration slice: it changes consumers that go through the public layer barrels while intentionally leaving direct story imports and CSS reference imports in place for parity testing.

I added a Storybook parity story that renders representative local `web/src` components next to package components, then wrote a css-visual-diff script that builds Storybook, serves it, and compares the local and package sections. After removing intentionally different labels, css-visual-diff reports 0 changed pixels for the shared-component sample.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Start migrating duplicated shared components to package imports and use Storybook plus css-visual-diff to prove parity.

**Inferred user intent:** Make the package canonical without visual regressions and keep reusable visual evidence in the ticket.

### What I did
- Updated `web/src/components/atoms/index.ts` to re-export package atoms.
- Updated `web/src/components/foundation/index.ts` to re-export package foundation primitives and types.
- Updated `web/src/components/layout/index.ts` to re-export package layout primitives and types.
- Updated `web/src/components/molecules/index.ts` to keep local `CoveragePanel`/`QueryPresetList` while re-exporting package `AppNav`, `DataTable`, and `MetadataGrid`.
- Updated `web/src/widgets/index.ts` to re-export package Widget IR, renderer, action, and cell-renderer APIs.
- Added `web/src/components/PackageParity.stories.tsx`, rendering local and package versions of Button, Caption, AppNav, Panel, Stack, MetadataGrid, and DataTable side by side.
- Wrote `scripts/03_storybook_package_parity_compare.sh`.
- Ran the script; it built Storybook and generated css-visual-diff artifacts under `sources/visual-parity/package-parity/`.
- Removed the generated `storybook-static` directory from the ticket output after producing compare artifacts, because the script can regenerate it and it is large.
- Ran validation commands.

### Why
- Barrel redirection is a low-risk way to make many web consumers use the package while leaving local files available for direct parity stories.
- The parity story allows css-visual-diff to compare local and package components in the same Storybook runtime.

### What worked
- `pnpm --dir web typecheck` passed.
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `pnpm --dir web build` passed.
- `pnpm --dir web build-storybook` passed inside the parity script.
- `css-visual-diff compare` reported:
  - `changed_pixels: 0`
  - `changed_percent: 0`
  - no computed style diffs for selected props.

### What didn't work
- The first parity story failed typecheck because `AppNav` requires `onItemSelect`:
  - `Property 'onItemSelect' is missing ... but required in type 'AppNavProps'.`
- I fixed the story by passing `onItemSelect={() => {}}` to both local and package `AppNav` instances.
- The first visual comparison had `0.3385%` changed pixels because the local and package headings used different text (`Local web/src components` vs `Package components`). I changed both headings to `Shared components`, then reran css-visual-diff and got `0%` changed pixels.

### What I learned
- Local and package implementations are visually identical for the representative shared-component sample when rendered with the same props.
- Keeping local direct imports in the parity story is useful while migrating; direct story imports should not be deleted until package Storybook ownership is decided.

### What was tricky to build
- If the parity story labels differ, css-visual-diff correctly reports pixels changed even when components are identical. The comparison fixture must keep all non-component text equal.
- Storybook build output is too large to keep in the ticket; the durable artifacts are the compare screenshots/report and the script.

### What warrants a second pair of eyes
- The barrel strategy means many consumers now use the package, but direct imports from individual local component files remain in stories and CSS reference tooling. Reviewers should confirm whether to keep that as the temporary parity harness or move stories into the package next.

### What should be done in the future
- Audit remaining direct local imports.
- Move or rewrite stories to import package components where appropriate.
- Delete duplicated local component implementation files only after direct imports are gone or intentionally retained for parity.

### Code review instructions
- Review the five barrel files under `web/src/components/**/index.ts` and `web/src/widgets/index.ts`.
- Review `web/src/components/PackageParity.stories.tsx`.
- Re-run `scripts/03_storybook_package_parity_compare.sh` and inspect `sources/visual-parity/package-parity/compare-shared-components/01-compare.md`.

### Technical details
- css-visual-diff report: `sources/visual-parity/package-parity/compare-shared-components/01-compare.md`.
- Pixel diff result: `0 changed pixels / 0%`.

## Step 4: Move Shared Storybook Coverage into the Package and Add Missing Web View Stories

I moved the Storybook ownership boundary to match the new component ownership boundary. Shared atoms, foundation primitives, layout primitives, shared molecules, and the WidgetRenderer story now live inside `packages/rag-evaluation-site`, while `web` keeps app-specific stories and the temporary package-parity fixture used for migration review.

I also filled in Storybook gaps for web-owned views that were previously missing direct stories. Backend-connected view stories use a lightweight Storybook mock API/provider so they render without a live Go backend and without moving RTK Query containers into the reusable package.

### Prompt Context

**User prompt (verbatim):** "move all stories related to components now in the package into the package, yes. then add stories to the components that are missing some."

**Assistant interpretation:** Move shared component and WidgetRenderer stories from `web` into `packages/rag-evaluation-site`, add package Storybook wiring, and add missing stories for remaining web-owned components/views.

**Inferred user intent:** Make the package independently reviewable in Storybook and ensure app-specific components still have visual coverage after moving shared stories out of `web`.

### What I did
- Added package Storybook config under `packages/rag-evaluation-site/.storybook/`.
- Added `storybook` and `build-storybook` scripts plus Storybook dev dependencies to `packages/rag-evaluation-site/package.json`.
- Moved shared stories from `web/src/components/{atoms,foundation,layout,molecules}` to the matching package component folders.
- Moved `web/src/widgets/WidgetRenderer.stories.tsx` to `packages/rag-evaluation-site/src/widgets/WidgetRenderer.stories.tsx`.
- Excluded `*.stories.*` from `packages/rag-evaluation-site/tsconfig.build.json` so package declaration builds do not publish story declarations.
- Added `web/src/storybook/MockApiProvider.tsx` for app/container stories that need RTK Query and mocked `/api/v1` responses.
- Added missing web-owned stories for `CorpusExplorerView`, `EmbeddingsView`, `EvaluationView`, `PipelineView`, `SearchView`, `WorkflowsView`, `DslPreviewPage`, `SearchWorkbenchPage`, `MacButton`, `MacMenuBar`, and `MacWindow`.
- Updated the design guide with the accepted Storybook ownership decision.

### Why
- The reusable package should be reviewable without relying on web-local duplicate component implementations.
- Web containers should remain in `web`, but they still need Storybook coverage for visual review and regression discovery.

### What worked
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `pnpm --dir web typecheck` passed.
- `pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-evaluation-site-storybook` passed.
- `pnpm --dir web exec storybook build --output-dir /tmp/rag-web-storybook` passed.
- `pnpm --dir packages/rag-evaluation-site build` passed, confirming stories are excluded from the production package declaration build.

### What didn't work
- Running `pnpm --dir packages/rag-evaluation-site build-storybook -- --output-dir /tmp/rag-evaluation-site-storybook` failed with:
  - `error: too many arguments for 'build'. Expected 0 arguments but got 2.`
- I reran Storybook through `pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-evaluation-site-storybook`, which passed.
- The first `web` typecheck failed because `MacMenuBar.stories.tsx` had no default `args`, and `MockApiProvider.tsx` used `Array.prototype.at` despite the repo targeting `ES2020`. I added default args and replaced `.at(-1)` with `.pop()`.

### What I learned
- Package Storybook can build cleanly with the existing package CSS import (`src/styles.css`) and does not need web `index.css`.
- The app-level missing stories need mocked API responses rather than moving backend-connected containers into the package.

### What was tricky to build
- Storybook ownership is now split: package Storybook owns reusable components, while web Storybook owns backend-connected application containers. The tricky part was avoiding accidental package coupling to `web/src/services/api` while still giving `web` view stories realistic data.
- The build config also needed a separate guard: package typecheck may include stories, but package declaration build should exclude them to avoid shipping Storybook artifacts.

### What warrants a second pair of eyes
- The `MockApiProvider` response shapes should be reviewed against backend API contracts. They are intentionally small, but if they drift too far from production responses, view stories can become misleading.
- The temporary `PackageParity.stories.tsx` still imports web-local duplicate component files directly for comparison. That is useful during migration, but it should be deleted once local duplicates are removed.

### What should be done in the future
- Add css-visual-diff sweeps over package Storybook stories, not only the package parity fixture.
- Decide when to delete the local duplicated web component implementations now that package stories no longer depend on them.
- Consider moving `CoveragePanel` and `QueryPresetList` into the package once their ownership is decided.

### Code review instructions
- Start with `packages/rag-evaluation-site/.storybook/main.ts` and `packages/rag-evaluation-site/.storybook/preview.ts`.
- Review one moved story per layer, for example `packages/rag-evaluation-site/src/components/atoms/Button/Button.stories.tsx`, `packages/rag-evaluation-site/src/components/layout/Panel/Panel.stories.tsx`, and `packages/rag-evaluation-site/src/widgets/WidgetRenderer.stories.tsx`.
- Review `web/src/storybook/MockApiProvider.tsx` and the new web view stories.
- Validate with:
  - `pnpm --dir packages/rag-evaluation-site typecheck`
  - `pnpm --dir web typecheck`
  - `pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-evaluation-site-storybook`
  - `pnpm --dir web exec storybook build --output-dir /tmp/rag-web-storybook`
  - `pnpm --dir packages/rag-evaluation-site build`

### Technical details
- Package stories moved: atoms (6), foundation/overview (7 files plus CSS module), layout (8), shared molecules (3), WidgetRenderer (1).
- New web-owned stories added: corpus, embeddings, evaluation, pipeline, search, workflows, DSL preview, search workbench, and retro Mac components.

## Step 5: Remove Remaining Duplicated Web Shared Component Implementations

I removed the now-unneeded local `web/src` copies of package-owned atoms, foundation primitives, layout primitives, shared molecules, and WidgetRenderer implementation files. The web app now keeps only thin compatibility barrels plus web-owned components such as `CoveragePanel`, `QueryPresetList`, corpus/workflow organisms, pages, and views.

The one non-obvious migration surface was the standalone design-reference CSS extraction entry. It previously imported CSS modules from the web-local component copies so Vite would emit CSS matching server-rendered component class names. After deleting the copies, that entry imports the package source CSS modules instead, keeping the reference renderer aligned with the package-owned components.

### Prompt Context

**User prompt (verbatim):** "go ahead"

**Assistant interpretation:** Continue with the next unification step: delete remaining local duplicate component/widget implementations now that package Storybook ownership is in place.

**Inferred user intent:** Finish making `packages/rag-evaluation-site` the canonical implementation layer and leave `web` with app-specific code only.

### What I did
- Deleted web-local duplicate directories for package-owned atoms, foundation primitives, layout primitives, and shared molecules (`AppNav`, `DataTable`, `MetadataGrid`).
- Deleted web-local duplicate widget implementation files: `actions.ts`, `cellRenderers.tsx`, `ir.ts`, and `WidgetRenderer.tsx`.
- Deleted the temporary `web/src/components/PackageParity.stories.tsx` fixture because it depended on the local duplicate component files being removed.
- Updated `web/src/services/api.ts` and `web/src/storybook/MockApiProvider.tsx` to import Widget IR types/helpers from `@go-go-golems/rag-evaluation-site/ir`.
- Updated `web/src/reference-pages/css-entry.ts` to import package source CSS modules for the standalone design-reference CSS extraction path.
- Updated `web/src/components/molecules/CoveragePanel/CoveragePanel.tsx` to import package-owned `DataTable`, `MetadataGrid`, and `DataTableColumn` directly from the package.
- Fixed `web/vite.config.ts` package aliases so the `/app` and `/ir` subpath aliases resolve before the root package alias.

### Why
- Leaving duplicate implementations behind would keep ambiguity about which component version is canonical.
- Deleting duplicates after moving stories into the package ensures package Storybook remains the visual review surface for reusable components.

### What worked
- `pnpm --dir web typecheck` passed.
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `pnpm --dir web build` passed.
- `pnpm --dir web exec storybook build --output-dir /tmp/rag-web-storybook-dedup` passed.
- `pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-package-storybook-dedup` passed.
- `pnpm --dir packages/rag-evaluation-site build` passed.

### What didn't work
- The first web Storybook build after deletion failed because Vite resolved `@go-go-golems/rag-evaluation-site/ir` through the root package alias, producing `../packages/rag-evaluation-site/src/index.ts/ir`.
- I changed `web/vite.config.ts` from an object alias map to an ordered alias array with `/app` and `/ir` before the root alias.
- The first web typecheck failed after deleting local `DataTable`/`MetadataGrid` because `CoveragePanel` imported those components by relative sibling paths. I changed it to import them from `@go-go-golems/rag-evaluation-site`.

### What I learned
- Thin compatibility barrels are enough for existing app imports; the actual duplicated implementation directories can be removed once direct imports are gone.
- Vite alias ordering matters for package subpaths when the root package alias is also present.

### What was tricky to build
- The design-reference CSS extraction entry was easy to miss because it does not render components directly; it imports CSS module objects to force Vite to include CSS for server-rendered reference pages. Deleting local CSS modules required redirecting those imports to package source CSS modules rather than removing the entry.
- The temporary parity story was intentionally useful before deletion, but became a liability once local component files were removed.

### What warrants a second pair of eyes
- Review `web/src/reference-pages/css-entry.ts` to confirm importing package source CSS modules is acceptable for the design-reference build path.
- Review whether the thin `web/src/components/{atoms,foundation,layout}/index.ts` and `web/src/widgets/index.ts` barrels should remain as app compatibility shims or whether imports should be rewritten directly to package imports in a later cleanup.

### What should be done in the future
- Run a css-visual-diff sweep against package Storybook stories as the new canonical reusable-component visual baseline.
- Consider moving `CoveragePanel` and `QueryPresetList` to the package if they are reusable outside the app.

### Code review instructions
- Start with deleted paths under `web/src/components/{atoms,foundation,layout}` and `web/src/widgets` to verify only compatibility barrels remain.
- Review `web/vite.config.ts` alias ordering.
- Review `web/src/reference-pages/css-entry.ts` package CSS module imports.
- Validate with the commands listed in `What worked`.

### Technical details
- Remaining web shared-layer files are compatibility barrels plus web-only molecules.
- Package-owned implementation and Storybook files now live under `packages/rag-evaluation-site/src/components` and `packages/rag-evaluation-site/src/widgets`.

## Step 6: Capture Package Storybook Visual Sweep Baselines

I added a package Storybook visual sweep script and ran it across every package-owned story. The sweep builds the package Storybook, serves the static output locally, discovers stories from `index.json`, and runs `css-visual-diff compare` for each story against itself to produce a deterministic screenshot/report baseline for review.

This is not a before/after parity comparison; it is a canonical package baseline now that the duplicated web implementations have been removed. It gives us per-story screenshots, compare JSON, markdown reports, and an aggregate summary TSV for future regression sweeps.

### Prompt Context

**User prompt (verbatim):** "go ahead"

**Assistant interpretation:** Continue with the remaining visual QA task: run css-visual-diff across package-owned Storybook stories.

**Inferred user intent:** Establish a broader visual baseline for the package components after moving stories and deleting web-local duplicates.

### What I did
- Added `scripts/04_package_storybook_visual_sweep.sh`.
- The script builds `packages/rag-evaluation-site` Storybook to a ticket-local static directory, serves it on a local port, fetches `index.json`, and enumerates package stories.
- The script runs `css-visual-diff compare` for each story using `#storybook-root` as the selector on both sides.
- The script writes per-story artifacts under `sources/visual-parity/package-story-sweep/reports/<story-id>/` and an aggregate `summary.tsv`.
- Removed generated `storybook-static` after the run so the committed ticket contains durable artifacts, not the whole build output.

### Why
- After deleting web-local duplicates, the package Storybook is the canonical reusable-component visual surface.
- A full story sweep gives reviewers concrete screenshots for all exported package component/widget stories instead of only one parity fixture.

### What worked
- The sweep completed for 48 package stories.
- `summary.tsv` contains 49 lines including the header.
- All story self-comparisons reported `changed_pixels=0` and `changed_percent=0`.
- `docmgr doctor --ticket RAGEVAL-DESIGN-SYSTEM-UNIFY --stale-after 30` passed after generating the artifacts.

### What didn't work
- N/A for this step. The script completed successfully on the first run.

### What I learned
- The package Storybook index currently exposes 48 stories across atoms, foundation, layout, molecules, and Widget IR.
- The sweep output is about 8 MB after removing `storybook-static`, which is acceptable for a ticket visual baseline.

### What was tricky to build
- The css-visual-diff reports are Markdown files under the docmgr ticket tree, so the script prefixes each generated report with valid docmgr frontmatter and renames it to `01-compare.md` inside each story directory.
- The script intentionally removes raw `compare.md` files without frontmatter to keep `docmgr doctor` clean.

### What warrants a second pair of eyes
- Review whether committing all per-story PNG baselines is the desired long-term storage strategy, or whether future sweeps should keep only summary JSON plus selected screenshots.
- Spot-check a few generated screenshots in `sources/visual-parity/package-story-sweep/reports/` to confirm the stories render the intended states.

### What should be done in the future
- Convert this self-compare baseline into a true previous-vs-current regression workflow once there is a stored baseline from a prior commit/build.
- Optionally add a review-site index for browsing the generated screenshots more comfortably.

### Code review instructions
- Review `scripts/04_package_storybook_visual_sweep.sh` first.
- Inspect `sources/visual-parity/package-story-sweep/summary.tsv` for story count and changed-pixel results.
- Spot-check reports such as:
  - `sources/visual-parity/package-story-sweep/reports/design-system-atoms-button--variants/01-compare.md`
  - `sources/visual-parity/package-story-sweep/reports/component-library-molecules-datatable--retrieval-rows/01-compare.md`
  - `sources/visual-parity/package-story-sweep/reports/widget-ir-renderer--search-workbench-composition/01-compare.md`
- Re-run with `scripts/04_package_storybook_visual_sweep.sh`.

### Technical details
- Storybook port used by the script: `${CSSVD_PACKAGE_STORYBOOK_PORT:-18617}`.
- Viewport: `1440x1000`.
- Wait: `1500ms` for both URL captures.
- Selector: `#storybook-root`.

## Step 7: Convert Package Sweep Output into css-visual-diff Review Site Data

I converted the package Storybook sweep artifacts into the directory and manifest format consumed by `css-visual-diff serve`. The review site is now running locally and exposes all 48 package Storybook baseline cards through the interactive UI.

This step did not rerun browser comparisons; it reorganized the already-captured sweep artifacts into the review-site contract: a root `summary.json`, one `package-storybook` page, and per-story `artifacts/<story-id>/` directories with `compare.json`, `left_region.png`, `right_region.png`, `diff_only.png`, and `diff_comparison.png`.

### Prompt Context

**User prompt (verbatim):** "Continue. is there anything left on th eunification of the design system?"

**Assistant interpretation:** Finish wiring the css-visual-diff review-site serving flow and assess remaining design-system unification work.

**Inferred user intent:** Make the visual baseline easy to review interactively and understand whether the core unification task is complete or still has important follow-ups.

### What I did
- Added `scripts/05_prepare_package_storybook_review_site.py`.
- Generated `sources/visual-parity/package-story-sweep-review/summary.json` with 48 rows.
- Created serve-compatible artifact paths under `package-storybook/artifacts/<story-id>/` using symlinks to the existing sweep artifacts.
- Started `css-visual-diff serve` in tmux session `cssvd-package-review` on port `8097`.
- Verified `GET /api/manifest` reports 48 accepted rows.
- Fixed docmgr warnings by not creating `compare.md` symlinks in the serve-compatible artifact tree; the review site does not need Markdown reports.

### Why
- The existing sweep data was useful but inconvenient to browse manually.
- `css-visual-diff serve` gives an interactive review interface for scanning all package-owned component/widget stories and annotating any visual concerns.

### What worked
- `css-visual-diff serve` is running at `http://127.0.0.1:8097`.
- Manifest check passed:
  - `sectionCount: 48`
  - `classificationCounts.accepted: 48`
  - `policy.ok: true`
- `docmgr doctor --ticket RAGEVAL-DESIGN-SYSTEM-UNIFY --stale-after 30` passed.

### What didn't work
- The first generated review-site layout symlinked `compare.md` into every artifact directory. `docmgr doctor` warned about 48 unprefixed Markdown files.
- I removed the `compare.md` symlinks from the serve layout and regenerated the review-site directory, because the serve UI only requires JSON and PNG artifacts.

### What I learned
- `css-visual-diff serve` expects a stricter layout than raw `compare` output: `summary.json` plus `<page>/artifacts/<section>/...`.
- The review site can use symlinked artifacts, so we do not need to duplicate the PNGs already committed under `package-story-sweep/reports`.

### What was tricky to build
- The review UI rewrites absolute paths based on the `<page>/artifacts/<section>/<file>` pattern. The converter therefore needs to create that exact layout even though the source sweep used `reports/<story-id>/`.
- Keeping docmgr clean required avoiding extra Markdown files in deep artifact directories.

### What warrants a second pair of eyes
- The converter script is intentionally tailored to package Storybook baselines. Review whether it should become a generic css-visual-diff utility or remain ticket-local.
- Review whether symlinks in ticket artifacts are acceptable for long-term portability. If not, switch `link_or_copy` to always copy files.

### What should be done in the future
- If the review site annotations are valuable, export the review notes from the UI and attach them to the ticket.
- Build a true baseline-vs-current comparison script using this review-site layout as the reporting target.

### Code review instructions
- Open `http://127.0.0.1:8097` while `tmux attach -t cssvd-package-review` is running.
- Review `scripts/05_prepare_package_storybook_review_site.py`.
- Validate with:
  - `curl -fsS http://127.0.0.1:8097/api/manifest | jq '.sectionCount, .classificationCounts, .policy'`
  - `docmgr doctor --ticket RAGEVAL-DESIGN-SYSTEM-UNIFY --stale-after 30`

### Technical details
- Review site data root: `sources/visual-parity/package-story-sweep-review/`.
- tmux session: `cssvd-package-review`.
- URL: `http://127.0.0.1:8097`.

## Step 8: Make Review-Site Screenshots Human-Readable

The first review-site conversion technically served images, but the card images were hard to inspect because many Storybook screenshots were as wide as the viewport while their content was only a small control. I updated the converter so the review-site `left_region.png` and `right_region.png` files are whitespace-trimmed copies with a small white border, while the raw sweep artifacts remain unchanged for traceability.

I verified the fix in the browser: the first review card now shows readable side-by-side button screenshots instead of tiny wide strips. The labels still say `Prototype` and `React` because that is css-visual-diff review-site terminology, but in this baseline both sides are package Storybook captures.

### Prompt Context

**User prompt (verbatim):** "oh they're super small for some reason it looks like? can we make them not so wide?"

**Assistant interpretation:** Fix the review-site presentation so screenshot crops are not dominated by Storybook root whitespace.

**Inferred user intent:** Make the interactive review UI usable for human visual inspection, not just technically serveable.

### What I did
- Updated `scripts/05_prepare_package_storybook_review_site.py` to run ImageMagick `convert -trim +repage -bordercolor white -border 12` for review-site left/right region images.
- Regenerated `sources/visual-parity/package-story-sweep-review/`.
- Restarted `css-visual-diff serve` in tmux session `cssvd-package-review`.
- Opened the page with Playwright, clicked the first review card, and captured `cssvd-review-trimmed.png` for verification.

### Why
- The raw `#storybook-root` crops are correct for measurement but often too wide for review because Storybook root spans the viewport.
- The review site is for human inspection, so it should use tighter display crops.

### What worked
- The first card now shows clear, readable side-by-side images for the button variants.
- `docmgr doctor --ticket RAGEVAL-DESIGN-SYSTEM-UNIFY --stale-after 30` passed.

### What didn't work
- N/A after the trimming update.

### What I learned
- The review site labels are generic (`Prototype`/`React`) and cannot be interpreted literally for this self-baseline use case.
- For small Storybook controls, whitespace trimming is necessary for usable review cards.

### What was tricky to build
- I kept raw sweep artifacts unchanged and only trimmed the serve-layout copies, so future machine comparisons can still refer back to the original css-visual-diff outputs.

### What warrants a second pair of eyes
- Spot-check larger stories to ensure trimming does not remove meaningful whitespace that designers care about.

### What should be done in the future
- If true baseline-vs-current comparisons are added, use raw captures for diff computation and trimmed copies only for human review presentation.

### Code review instructions
- Review `trim_for_review` in `scripts/05_prepare_package_storybook_review_site.py`.
- Open `http://127.0.0.1:8097`, expand the first few cards, and confirm the images are readable.

### Technical details
- Requires ImageMagick `convert`; if unavailable, the converter falls back to copying the raw image.

## Step 9: Remove Web Compatibility Barrels and Close the Ticket

With the package now established as the canonical owner of shared UI primitives and the widget renderer, I removed the last thin web-side compatibility barrels. Web code now imports shared atoms, foundation primitives, layout primitives, shared molecules, and Widget IR APIs directly from `@go-go-golems/rag-evaluation-site`, while app-specific molecules continue to import from their concrete local directories.

This leaves the web app with a clearer ownership boundary: package-owned UI comes from the package, and web-owned application/container code remains under `web/src`. After typecheck and production build passed, I removed the placeholder task and closed the design-system unification ticket so the next work can return to context-viewer integration.

### Prompt Context

**User prompt (verbatim):** "clean up the compatibility barrels, then close the ticket, so we can move back to integrating the new context-viewer"

**Assistant interpretation:** Finish the optional import hygiene by deleting web re-export compatibility barrels, validate the app, and mark the docmgr ticket complete.

**Inferred user intent:** Clear the remaining design-system cleanup from the queue and unblock the next feature track.

**Commit (code):** bfe2dab4893d0d9bac03390ba34c01253cf2f09b — "refactor(web): remove shared compatibility barrels"

### What I did
- Rewrote web imports from local `atoms`, `foundation`, `layout`, `molecules`, and `widgets` barrels to direct `@go-go-golems/rag-evaluation-site` imports where the symbol is package-owned.
- Kept `CoveragePanel` and `QueryPresetList` as web-owned molecules and imported them from their concrete local directories.
- Deleted `web/src/components/atoms/index.ts`, `web/src/components/foundation/index.ts`, `web/src/components/layout/index.ts`, and `web/src/widgets/index.ts`.
- Updated `web/src/components/index.ts` and `web/src/components/molecules/index.ts` so they no longer re-export package-owned compatibility APIs.
- Ran `pnpm --dir web typecheck` and `pnpm --dir web build`.

### Why
- Compatibility barrels hid the ownership boundary after the duplicated implementations were removed.
- Direct package imports make it obvious which APIs are reusable package APIs and which APIs remain application-specific.

### What worked
- `pnpm --dir web typecheck` passed.
- `pnpm --dir web build` passed.
- The placeholder task was removed, leaving all ticket tasks completed.

### What didn't work
- The first validation run failed because `web/src/App.tsx` still imported `AppShell` and `AppNav` through now-deleted local barrels:
  - `src/App.tsx(9,26): error TS2307: Cannot find module './components/layout' or its corresponding type declarations.`
  - `src/App.tsx(10,10): error TS2305: Module '"./components/molecules"' has no exported member 'AppNav'.`
- Rewriting those imports to `@go-go-golems/rag-evaluation-site` fixed the failure.

### What I learned
- The web app had one direct `./components/layout` import outside the earlier relative-path cleanup patterns.
- `web/src/components/molecules/index.ts` should stay local-only unless `CoveragePanel` or `QueryPresetList` are promoted later.

### What was tricky to build
- The cleanup needed to distinguish package-owned molecules (`AppNav`, `DataTable`, `MetadataGrid`) from app-owned molecules (`CoveragePanel`, `QueryPresetList`). Rewriting every molecules import to the package would have been incorrect, so `SearchWorkbenchPage` now imports local app molecules from their concrete subdirectories.

### What warrants a second pair of eyes
- Review whether `CoveragePanel` and `QueryPresetList` should remain app-owned or become package components in a future pass.
- Review the larger import diff for readability; many files now import several symbols from the package directly.

### What should be done in the future
- Resume context-viewer integration using the now-stable package/web ownership boundary.

### Code review instructions
- Start with `web/src/App.tsx`, `web/src/components/pages/SearchWorkbenchPage/SearchWorkbenchPage.tsx`, and the deleted barrel files.
- Validate with `pnpm --dir web typecheck` and `pnpm --dir web build`.

### Technical details
- Package-owned import target: `@go-go-golems/rag-evaluation-site`.
- App-owned molecules retained locally: `web/src/components/molecules/CoveragePanel` and `web/src/components/molecules/QueryPresetList`.
