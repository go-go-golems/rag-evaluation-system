# Changelog

## 2026-07-03

- Initial workspace created


## 2026-07-03

Investigated design system (3 line-anchored inventories), ran package Storybook on :6007, captured 14 story screenshots to sources/screenshots/, and wrote the full CMS widget-set design doc (gap analysis, 20 new components across atoms/layout/molecules/organisms, 7 decision records, 7-phase plan, backend/API sketch).

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/03/RAGEVAL-CMS-WIDGETS--cms-widget-set-atoms-molecules-organisms-for-images-articles-and-media-in-the-rag-evaluation-site-design-system/design-doc/01-cms-widget-set-analysis-design-and-implementation-guide.md — Primary deliverable
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/03/RAGEVAL-CMS-WIDGETS--cms-widget-set-atoms-molecules-organisms-for-images-articles-and-media-in-the-rag-evaluation-site-design-system/reference/01-investigation-diary.md — Chronological investigation record


## 2026-07-03

Validated ticket with docmgr doctor (all checks pass; added 'cms' topic to vocabulary) and uploaded the design doc + diary as a single ToC'd PDF bundle to reMarkable at /ai/2026/07/03/RAGEVAL-CMS-WIDGETS.


## 2026-07-03

Investigated and documented the goja Widget DSL (ui.dsl/data.dsl/context_window.dsl/course.dsl): registration paths, full helper/recipe surface, paletteStyleSet palettes, TS typing generation, 5 drift findings (typescript.go omissions, missing ContextGroupedStripDiagram helper, DSL dormant in-repo), and a concrete cms.dsl extension plan. New doc: reference/02-goja-widget-dsl-reference-and-cms-extension-notes.md; design-doc §4.6 cross-linked.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/03/RAGEVAL-CMS-WIDGETS--cms-widget-set-atoms-molecules-organisms-for-images-articles-and-media-in-the-rag-evaluation-site-design-system/reference/02-goja-widget-dsl-reference-and-cms-extension-notes.md — New goja DSL reference


## 2026-07-03

Expanded tasks.md into a full implementation breakdown: Phase 0 groundwork (sanitizeUrl, cms types/fixtures, gallery block), Phase 1 atoms (4), Phase 2 layout (2), Phase 3 molecules (8), Phase 4 organisms (7) — each task with API contract pointer, story-state requirements, and validation gates; Phase V validation (storybook build, style grep gates, playwright smoke); Phase 5 (backend/web) and Phase 6 (cms.dsl goja DSL) explicitly deferred until after the widgets land, per user direction.


## 2026-07-03

Implemented Phases 0-4: sanitizeUrl in MarkdownArticle, cms types/fixtures, gallery ArticleBlock; atoms MediaThumb/Tag/ContentStatusBadge/MeterBar; layout TileGrid/DialogShell; molecules AssetTile/TagListInput/Breadcrumbs/Pagination/SearchField/EmptyState/UploadQueueList/MarkdownEditor; organisms ConfirmDialog/MediaLibraryPanel/AssetPickerDialog/ArticleListPanel/ArticleEditorPanel/AssetDetailPanel/CmsShell; 103 new stories; storybook staticDirs + fixture SVG (fixes long-broken course-assets URLs); IconButton size=large for row actions (user feedback). Validation: typecheck/biome/build-storybook/package build pass; Playwright smoke over 116 stories: 0 errors, 0 missing data-rag roots. Known pre-existing failure: web typecheck (ContextVisualizerPage missing styleSet), verified via git stash.


## 2026-07-03

Feedback round: row-action icons were barely legible (bare 10-14px dim glyphs next to bordered badges). Added IconButton variant=boxed (1px border, full-contrast glyph, invert-on-hover, min 22x20 hit area) and applied it with size=large to ArticleListPanel and UploadQueueList row actions; new Boxed story; re-verified via screenshots.


## 2026-07-03

Studied go-go-course (the real CMS consumer: xgoja binary, IR pages from JS, implemented widget server actions, file-backed course/ storage) and revised the design doc to Widget-DSL-first delivery per user direction: dropped the Go-REST/RTK-Query Phase 5, added D-8 (accepted) and D-9 (ActionSpec.confirm), rewrote §6.7-6.9/§8/§9/§11/§12, replaced tasks D1/D2 with P5.1-P5.7 (IR enablement) and P6.1-P6.6 (go-go-course integration).


## 2026-07-03

Implemented Phase 5 (IR enablement): ActionSpec.confirm with interpolated window.confirm gate in dispatchWidgetAction (D-9); shared widgets/uploadSerialization.ts (extracted from ContextUploadDropArea.widget, reused by MediaLibraryPanel adapter); name/defaultValue form modes on MarkdownEditor + TagListInput + onQuerySubmit on Media/Article panels; 14 widget adapters + manifests under new cms.dsl module (cmsWidgetRegistry merged into defaultWidgetRegistry); 7 WidgetRenderer/CMS stories rendering the full CMS from JSON IR; Go cms.dsl (14 helpers + mediaLibrary/articleList recipes), rag-widget-site provider entry, boundary tests incl. confirm passthrough; drift fixed (typescript.go linkButton/actionButton/download, contextGroupedStripDiagram helper, doc/02 updated with cms.dsl section + confirm docs). Validation: tsc, biome, build-storybook, go test widgetdsl+widgetsite, Playwright smoke of all 7 IR stories (0 errors).


## 2026-07-04

Added intern onboarding playbook (playbook/01): status snapshot, docs map, run commands, ordered P6 start guide, gotchas.

