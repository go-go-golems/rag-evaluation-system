# Changelog

## 2026-06-07

- Initial workspace created


## 2026-06-07

Initialized design-system unification ticket, captured css-visual-diff help, generated duplication inventory, and wrote the analysis/implementation guide with granular task plan.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--unify-web-design-system-into-reusable-rag-evaluation-site-package/design-doc/01-design-system-unification-analysis-and-implementation-guide.md — Unification guide and task plan
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--unify-web-design-system-into-reusable-rag-evaluation-site-package/reference/01-investigation-diary.md — Step 1 diary
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--unify-web-design-system-into-reusable-rag-evaluation-site-package/sources/01-duplication-inventory.md — Duplication evidence


## 2026-06-07

Added pnpm workspace wiring so web depends on the local rag-evaluation-site package for development; typecheck passes for both frontend packages.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/pnpm-workspace.yaml — New workspace file
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/package.json — Workspace dependency
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/tsconfig.json — Package source path mapping
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/vite.config.ts — Package source aliasing


## 2026-06-07

Redirected shared web barrels/widgets to package exports, added Storybook package parity story, and captured css-visual-diff evidence with 0 changed pixels.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--unify-web-design-system-into-reusable-rag-evaluation-site-package/sources/visual-parity/package-parity/compare-shared-components/01-compare.md — Visual parity report
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/src/components/PackageParity.stories.tsx — New visual parity Storybook fixture


## 2026-06-07

Moved reusable component and WidgetRenderer stories into the package, added package Storybook wiring, and added missing web-owned view/retro stories with a mock API provider.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/.storybook/main.ts — Package Storybook configuration
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/src/storybook/MockApiProvider.tsx — Mocked Storybook API provider for app stories


## 2026-06-07

Removed duplicated web-local shared component and WidgetRenderer implementations, leaving package-owned implementations plus thin web compatibility barrels.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/src/reference-pages/css-entry.ts — Package CSS module imports for reference rendering
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/vite.config.ts — Alias fix for package subpaths


## 2026-06-07

Added and ran a package Storybook css-visual-diff sweep across 48 package-owned stories; all self-comparisons reported 0 changed pixels.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--unify-web-design-system-into-reusable-rag-evaluation-site-package/scripts/04_package_storybook_visual_sweep.sh — Visual sweep automation
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--unify-web-design-system-into-reusable-rag-evaluation-site-package/sources/visual-parity/package-story-sweep/summary.tsv — Visual sweep summary


## 2026-06-07

Prepared css-visual-diff review-site data for the package Storybook sweep and served it on localhost:8097.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--unify-web-design-system-into-reusable-rag-evaluation-site-package/scripts/05_prepare_package_storybook_review_site.py — Review-site data converter
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--unify-web-design-system-into-reusable-rag-evaluation-site-package/sources/visual-parity/package-story-sweep-review/summary.json — Review-site manifest


## 2026-06-07

Trimmed review-site left/right screenshots so package Storybook cards are human-readable instead of viewport-wide strips.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--unify-web-design-system-into-reusable-rag-evaluation-site-package/scripts/05_prepare_package_storybook_review_site.py — Trims review-site display images


## 2026-06-07

Removed web compatibility barrels and rewrote shared UI imports to @go-go-golems/rag-evaluation-site (commit bfe2dab).

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/src/App.tsx — Imports AppShell/AppNav directly from package
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/src/components/index.ts — Stops re-exporting deleted compatibility barrels
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/src/components/pages/SearchWorkbenchPage/SearchWorkbenchPage.tsx — Separates package primitives from app-owned molecules


## 2026-06-07

Ticket complete: shared design-system ownership moved to package, Storybook/visual review artifacts captured, compatibility barrels removed.

