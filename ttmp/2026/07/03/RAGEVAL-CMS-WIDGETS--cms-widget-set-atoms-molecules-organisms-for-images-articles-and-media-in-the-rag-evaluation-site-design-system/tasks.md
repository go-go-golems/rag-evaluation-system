# Tasks

Conventions for every implementation task below (from `packages/rag-evaluation-site/GUIDELINES.md` and design-doc/01 §6):

- Component folder = `Component.tsx` + `Component.module.css` + `Component.stories.tsx` + `index.ts`, exported from the layer barrel (`src/components/<layer>/index.ts`).
- Tokens only (`--mac-*`, `--rag-font-role-*`); no border-radius, no box-shadow, no raw font literals; selection = invert to `--mac-bg-dark`/`--mac-text-inv`.
- `data-rag-atom|layout|molecule|organism="<Name>"` on the root, plus state hooks (`data-state`, `data-active`, …).
- Story titles: `Design System/Atoms|Layout/<Name>` or `Component Library/Molecules|Organisms/<Name>`; cover default/populated, empty, dense/overflow, selected/active, disabled, error where applicable.
- Gate each task with `pnpm --dir packages/rag-evaluation-site typecheck` and `pnpm biome check --write .`.

API contracts for all components: design-doc/01 §6.3–§6.6. DTOs: §6.2.

## Done — investigation & deliverables

- [x] Create ticket workspace, design doc, and diary
- [x] Read AGENTS.md + package GUIDELINES.md (design-system rules)
- [x] Run package Storybook (port 6007) and enumerate all stories via index.json
- [x] Capture screenshots of 14 representative stories into sources/screenshots/
- [x] Inventory foundation/atoms/layout APIs (line-anchored)
- [x] Inventory molecules/organisms APIs + content DTOs (line-anchored)
- [x] Inventory Widget IR, web data layer, and Go backend surface
- [x] Write design doc: analysis, gap analysis, CMS widget set design, decision records, phased plan
- [x] Investigate + document the goja Widget DSL (reference/02) with cms.dsl extension notes
- [x] Maintain investigation diary; relate files; changelog; docmgr doctor; reMarkable upload

## Phase 0 — groundwork

- [x] P0.1 Add `sanitizeUrl()` to `src/components/molecules/MarkdownArticle/MarkdownArticle.tsx` for link `href` (:29) and image `src` (:100–104): allow http/https/mailto/relative, reject `javascript:`/`data:`; add `rel="noopener noreferrer"` on external links; add a story with hostile fixture (`javascript:` link, `data:text/html` image) proving they render inert
- [x] P0.2 Create `src/cms/types.ts` with `CmsContentStatus`, `CmsAsset`, `CmsArticleSummary`, `CmsArticleDetail`, `ArticleGalleryBlock` (design-doc §6.2); re-export from `src/index.ts` alongside `./context`
- [x] P0.3 Create `src/cms/fixtures.ts` mirroring `src/context/fixtures.ts`: ~12 assets (mixed image/file, one broken `src`, long filenames), ~8 article summaries (all four statuses), 1 article detail with markdown+image blocks, upload-queue items in every `UploadItemStatus`
- [x] P0.4 Add `gallery` member to the `ArticleBlock` union in `src/context/types.ts` (:170–192) and render it in `RichArticle` (grid of figures, 2–4 columns) + RichArticle story `GalleryBlock`

## Phase 1 — atoms (`src/components/atoms/`)

- [x] P1.1 `MediaThumb`: aspect square|wide|natural, fit cover|contain, frame bordered|none, loading checkerboard (reuse `.pattern_*` technique from ContextStyleSwatch.module.css), broken/empty fallback glyph + Caption, `selected` accent outline, `data-state` hook; stories: Default, Loading, Broken, Empty, Contain, Wide, Selected, DenseGridSample
- [x] P1.2 `Tag`: bordered mono chip, `selected` inversion, `onRemove` × via IconButton, disabled; stories: Default, Selected, Removable, Disabled, OverflowRow (12+ tags wrapping)
- [x] P1.3 `ContentStatusBadge`: draft|published|scheduled|archived → dim/green/accent/dim+line-through, glyphs ◌ ● ◔ ▣, `--rag-font-role-label` uppercase bordered; stories: AllStatuses, NoIcon, InTableRow
- [x] P1.4 `MeterBar`: value 0..1, tones accent|success|danger, optional metric label, 10px track `1px solid --mac-border`; stories: Progress (0/25/50/100), Tones, WithLabel
- [x] P1.5 Export all four from `src/components/atoms/index.ts`; typecheck + biome pass

## Phase 2 — layout (`src/components/layout/`)

- [x] P2.1 `TileGrid`: `repeat(auto-fill, minmax(var(--rag-tile-min), 1fr))`, minTileWidth prop → inline CSS var (SidebarShell pattern), gap sm|md; stories: Default, DenseFiftyTiles, NarrowContainer
- [x] P2.2 `DialogShell`: native `<dialog>` + `showModal()/close()` effect, `cancel`/`close` → onClose, Panel-style black title bar, footer slot, sizes sm|md|lg (420/640/920), flat `::backdrop` dim; stories: Static (rendered inline `open` without showModal, for visual diff), Interactive (useState open/close), Sizes, WithFooterActions
- [x] P2.3 Export both from `src/components/layout/index.ts`; typecheck + biome pass

## Phase 3 — molecules (`src/components/molecules/`)

- [x] P3.1 `AssetTile`: MediaThumb + truncated compact title + `PNG · 214 KB` metadata line, selection outline + inverted title row (DocumentListPanel active pattern), `onSelect`/`onOpen`, `footerSlot`, `data-rag-asset-id`; stories: Default, Selected, BrokenImage, FileKind, LongTitle, WithStatusFooter
- [x] P3.2 `TagListInput`: Tag row + borderless mini TextInput committing on Enter/comma, native `<datalist>` suggestions, onAdd/onRemove, disabled; stories: Default, Empty, ManyTags, WithSuggestions, Disabled, Interactive
- [x] P3.3 `Breadcrumbs`: mono metadata font, `/` separators, last item bold/unclickable, `<nav aria-label>`; stories: Default, Deep (6 levels), SingleItem, Interactive
- [x] P3.4 `Pagination`: compact ‹prev/next› Buttons + `page 3 / 9` Caption + optional `1–24 of 210`; stories: Default, FirstPage, LastPage, WithTotals, Interactive
- [x] P3.5 `SearchField`: TextInput wrap with ⌕ glyph + clear IconButton when non-empty, Enter → onSubmit; stories: Empty, WithValue, Disabled, Interactive
- [x] P3.6 `EmptyState`: glyph + title + mono hint + action slot, optional dashed `framed` variant; stories: Default, WithAction, Framed, InsidePanel
- [x] P3.7 `UploadQueueList`: rows of format glyph + CodeText filename + MeterBar + status + cancel/retry IconButtons, error caption in danger tone; stories: Mixed (all five UploadItemStatus), AllUploading, WithErrors, Empty, Interactive (simulated progress)
- [x] P3.8 `MarkdownEditor`: toolbar (B, code, link, H2, list, image via `wrapSelection` helper) + TextareaInput + counter, `onInsertAsset` hook; stories: Default, WithContent, MaxLength, Disabled, ToolbarSlot, Interactive
- [x] P3.9 Export all eight from `src/components/molecules/index.ts`; typecheck + biome pass

## Phase 4 — organisms (`src/components/organisms/`)

- [x] P4.1 `ConfirmDialog`: DialogShell sm + message + Cancel/confirm (danger-toned label, no new button variant); stories: Static, Destructive, Interactive
- [x] P4.2 `MediaLibraryPanel`: Panel(title MEDIA, kind TabList) + SearchField/Pagination toolbar + optional FileDropZone strip + UploadQueueList + ScrollRegion>TileGrid>AssetTile + EmptyState fallback (anatomy design-doc §6.6); stories: Populated, Empty, Uploading, MultiSelect, PickerMode, DenseFiftyAssets, Interactive
- [x] P4.3 `AssetPickerDialog`: DialogShell lg wrapping MediaLibraryPanel selectionMode=single, footer Cancel/Use-asset (primary disabled until selection); stories: Static, Interactive
- [x] P4.4 `ArticleListPanel`: DataTable<CmsArticleSummary> (title+slug CodeText caption, ContentStatusBadge cell, 2 Tags +n, author, mono updatedAt, IconButton row actions) + SearchField/status SelectInput/Pagination toolbar + ConfirmDialog for delete/archive; stories: Populated, Empty, Filtered, RowSelected, Overflow (40 rows), Interactive
- [x] P4.5 `ArticleEditorPanel`: FormPanel(ARTICLE, status vocab) + FormRows (title/slug/status/tags/cover MediaThumb+Choose…) + SplitPane(MarkdownEditor | ScrollRegion>MarkdownArticle live preview); stories: Draft, Saving, Success, Error, PreviewHidden, Interactive (local state round-trip incl. insert-image flow with AssetPickerDialog)
- [x] P4.6 `AssetDetailPanel`: DocumentPreviewToolbar + large MediaThumb natural/contain + MetadataGrid (copyable) + FormPanel (title/alt/tags/status) + used-in list + Delete via ConfirmDialog; stories: Image, File, BrokenImage, WithUsage, Interactive
- [x] P4.7 `CmsShell`: SidebarShell(188) + header + SidebarNav with exported `cmsNavSections` (Content: Articles/Media; Organize: Tags/Archive) + footer slot; stories: Default, MediaActive, WithFooter, Interactive
- [x] P4.8 Export all seven from `src/components/organisms/index.ts`; keep experimental ones out of `src/index.ts` npm barrel until stories stabilize (design-doc §11 risk); typecheck + biome pass

## Phase V — validation & wrap-up

- [x] PV.1 `pnpm --dir packages/rag-evaluation-site build-storybook` passes with all new stories
- [x] PV.2 Style gates: `rg -n "border-radius|box-shadow" packages/rag-evaluation-site/src/components` → no new hits; `rg -n "font-size" src/components/{atoms,molecules,organisms}` → only tokenized fallbacks
- [x] PV.3 Playwright smoke: iterate index.json for the new story titles, screenshot each `iframe.html?id=…`, assert zero console errors and `data-rag-*` roots present; save captures to this ticket's `sources/screenshots/`
- [x] PV.4 Update diary + changelog per phase; re-upload reMarkable bundle after Phase 4

## Phase 5 — IR enablement in this repo (revised per D-8: Widget-DSL-first, no REST/RTK; design-doc §6.7–6.9)

- [x] P5.1 `ActionSpec.confirm?: string` in `ir.ts` + central handling in `dispatchWidgetAction` (D-9); WidgetRenderer story covering a confirmed destructive action
- [x] P5.2 Extract `serializeUploadFile` from `ContextUploadDropArea.widget.tsx` into shared `src/widgets/uploadSerialization.ts`
- [x] P5.3 Uncontrolled/form modes: `name`/`defaultValue` on `MarkdownEditor` (adapter-internal live preview + synced named textarea) and `TagListInput` (hidden joined input)
- [x] P5.4 Widget adapters + manifests for MediaThumb, Tag, ContentStatusBadge, MeterBar, TileGrid, AssetTile, Breadcrumbs, Pagination, SearchField, EmptyState, MarkdownEditor, MediaLibraryPanel, ArticleListPanel, CmsShell (dialog/controlled surfaces skipped in IR v1); callbacks → `on*Action`; navigate-based selection/paging
- [x] P5.5 `cms.dsl` TS registry: `RagWidgetType` additions, `cmsWidgetRegistry`, `"cms.dsl"` in WidgetModule union, `WidgetRenderer.cms.stories.tsx`
- [x] P5.6 `cms.dsl` Go: helpers + `mediaLibrary`/`articleList` recipes in `pkg/widgetdsl`, provider entry in `widgetsite/provider.go`, typings, module_test boundary tests
- [x] P5.7 Fix DSL drift while there: typescript.go omissions (`cell.linkButton`, `cell.actionButton`, `action.download`), doc/02 reference, add `contextGroupedStripDiagram` helper

## Phase 6 — go-go-course integration (in ../go-go-course)

- [ ] P6.1 Release loop: npm publish rag-evaluation-site + Go module tag; bump `webapp/package.json` (0.1.16) and `xgoja.yaml`/`go.mod` (v0.1.2); rebuild webapp → assets → binary (iterate first via `dev:rag-local` source alias)
- [ ] P6.2 `course-material-service.js`: emit CmsAsset-shaped media entries (id/kind/src/mime/size); optional sidecar metadata JSON for tags/status/alt
- [ ] P6.3 `admin-course-cms.js`: media table → `cms.recipes.mediaLibrary` with thumbnails + navigate-selection; `confirm` on delete actions; contentStatusBadge in material tables
- [ ] P6.4 Handout editor page: formPanel + named markdownEditor (live preview), `/settings/handout-body` save route (form-post + redirect pattern)
- [ ] P6.5 End-to-end smoke in the running binary: upload image → thumbnail appears → reference from handout → edit with live preview → delete with confirm
- [ ] P6.6 Harden SVG uploads: run goja-text `sanitize` before write (or CSP on /course-assets) — validateSvgUpload only checks for an <svg> element today
