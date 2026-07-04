---
Title: Intern onboarding — status, docs map, and how to start Phase 6
Ticket: RAGEVAL-CMS-WIDGETS
Status: active
Topics:
    - design-system
    - frontend
    - storybook
    - cms
DocType: playbook
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "One-page onboarding: what shipped (CMS widget set + cms.dsl, Phases 0–5), where every document/screenshot/test lives, and the exact steps to start Phase 6 (go-go-course integration)."
LastUpdated: 2026-07-04T13:45:00-04:00
WhatFor: "Get a new contributor productive on the CMS widget work in under an hour."
WhenToUse: "Day one on RAGEVAL-CMS-WIDGETS, or before picking up any P6.x task."
---

# Intern onboarding — status, docs map, and how to start Phase 6

## 1. What this project is (30 seconds)

We extended the strict "Classic Mac" design system in `rag-evaluation-system/packages/rag-evaluation-site` (npm: `@go-go-golems/rag-evaluation-site`) with a **CMS widget set** — media library, article management, markdown editing — and exposed it to JavaScript page authors through the goja **Widget DSL** (`cms.dsl`). The consumer is the sibling repo **`go-go-course`**: an xgoja binary whose whole UI is Widget IR JSON built in JS and rendered by the embedded React SPA. There is **no REST API / RTK Query** in this plan — mutations go through widget server actions (`POST /api/widget/actions/:name`) and native form posts, storage is files under `go-go-course/.../course/`.

## 2. Current status (done / not done)

Done (Phases 0–5, all in `rag-evaluation-system`, **uncommitted in the working tree** — check `git status` before anything):

- 21 React components + ~110 Storybook stories: atoms `MediaThumb Tag ContentStatusBadge MeterBar`, layout `TileGrid DialogShell`, molecules `AssetTile TagListInput Breadcrumbs Pagination SearchField EmptyState UploadQueueList MarkdownEditor`, organisms `ConfirmDialog MediaLibraryPanel AssetPickerDialog ArticleListPanel ArticleEditorPanel AssetDetailPanel CmsShell`; plus `gallery` ArticleBlock, `sanitizeUrl` in MarkdownArticle, `IconButton` size=large/variant=boxed.
- `cms/` types + fixtures (`src/cms/{types,fixtures}.ts`), exported from the npm barrel.
- **cms.dsl** on both sides: 14 widget adapters/manifests (TS, `cmsWidgetRegistry`) and 14 goja helpers + `recipes.mediaLibrary`/`recipes.articleList` (Go, `pkg/widgetdsl`), provider entry in `pkg/xgoja/providers/widgetsite`.
- `ActionSpec.confirm` — every action can carry a confirm prompt, handled centrally in `dispatchWidgetAction` with `${}` interpolation.
- `src/widgets/uploadSerialization.ts` — shared File→utf8/base64 serialization; MediaLibraryPanel uploads use the same `{files, fileNames, fileCount}` contract go-go-course already decodes.
- Validation green: `typecheck`, `biome`, `build-storybook`, `go test ./pkg/widgetdsl ./pkg/xgoja/providers/widgetsite`, Playwright story smoke.

Not done (Phase 6 = your work, tasks **P6.1–P6.6** in tasks.md, all in `go-go-course`):

- Release/pin bump loop, CmsAsset mapping in `course-material-service.js`, admin CMS page upgrade to `cms.recipes.mediaLibrary`, `confirm` on deletes, handout editor page, end-to-end smoke, SVG sanitization.

## 3. Where to find everything

Ticket workspace: `rag-evaluation-system/ttmp/2026/07/03/RAGEVAL-CMS-WIDGETS--…/`

| What | Where |
|---|---|
| **The design doc** (read first: §6.7–6.9 revised architecture, §8 flows, §9 phases, decisions D-1…D-9) | `design-doc/01-cms-widget-set-analysis-design-and-implementation-guide.md` |
| Goja DSL reference + cms.dsl notes | `reference/02-goja-widget-dsl-reference-and-cms-extension-notes.md` |
| Chronological diary (what failed, why decisions were made — Steps 6–8 cover the implementation) | `reference/01-investigation-diary.md` |
| Task list with checkboxes (`docmgr task list --ticket RAGEVAL-CMS-WIDGETS`) | `tasks.md` |
| Screenshots (React stories + IR-rendered CMS) | `sources/screenshots/` and `sources/screenshots/cms/` |
| Design-system rules (non-negotiable: layers, tokens, no radius/shadow, data-rag-*) | `packages/rag-evaluation-site/GUIDELINES.md` + repo `AGENTS.md` |
| JS API reference for page authors (includes cms.dsl section) | `pkg/xgoja/providers/widgetsite/doc/02-widget-dsl-js-api-reference.md` |
| go-go-course anatomy | `go-go-course/cmd/go-go-course/`: `server.js` (routes + widget actions), `lib/course-material-service.js` (file-backed storage), `lib/pages/admin-course-cms.js` (page to upgrade), `xgoja.yaml` (buildspec/pins), `webapp/` (SPA shell) |

Everything is also on the reMarkable as one PDF: `/ai/2026/07/03/RAGEVAL-CMS-WIDGETS`.

## 4. Get it running (15 minutes)

```bash
cd rag-evaluation-system
pnpm install
pnpm --dir packages/rag-evaluation-site storybook   # http://localhost:6007
# Browse: Design System/Atoms/MediaThumb · Component Library/Organisms/MediaLibraryPanel
#         Widget IR/Renderer/CMS (the same widgets rendered from JSON)
pnpm --dir packages/rag-evaluation-site typecheck && pnpm biome check packages/rag-evaluation-site/src
go test ./pkg/widgetdsl ./pkg/xgoja/providers/widgetsite
```

For go-go-course (Phase 6): `go-go-course/AGENT.md` has build/test commands; the server runs via the built binary (`./dist/go-go-course run server.js --http-listen 127.0.0.1:8787 --keep-alive`, use tmux). Admin pages require setting your display name to `admin_<yourname>` in Settings.

## 5. How to start Phase 6 (in order)

1. **P6.1 — get new widgets into go-go-course without publishing:** use the source alias first — `cd go-go-course/cmd/go-go-course/webapp && RAG_SITE_SRC=<path-to>/rag-evaluation-system/packages/rag-evaluation-site pnpm dev` (see `dev:rag-local` script; adjust the relative path). The Go side needs a `replace` directive in go.mod (or workspace go.work) pointing at rag-evaluation-system so `cms.dsl` is available. Real releases (npm publish + Go tag + pin bumps in `webapp/package.json` 0.1.16 and go.mod v0.1.2) come once things work.
2. **P6.2 — CmsAsset mapping:** extend `lib/course-material-service.js` `listCourseMaterial()` to also emit `mediaAssets: CmsAsset[]` (`id=file`, `kind`/`mime` by extension, `src=/course-assets/<file>`, `size` bytes, `status:"published"`, `tags:[]`). Shape reference: design-doc §6.2 / doc/02 cms.dsl section.
3. **P6.3 — upgrade the admin page:** in `lib/pages/admin-course-cms.js`, replace the media DataTable with `cms.recipes.mediaLibrary({ assets, onAssetSelect: ui.action.navigate("?asset=$assetId"), onFilesSelected: "admin-upload-course-material", … })`; add `confirm: "Delete ${file}?"` to the delete actions in `lib/pages/admin-common.js`.
4. **P6.4 — handout editor page:** `formPanel({ method:"post", formAction:"/settings/handout-body" }, formRow({ control: cms.markdownEditor({ name:"body", defaultValue: source }) }))` + a save route in `server.js` following the existing form-post→validate→`writeFileAtomic`→redirect `?status=saved` pattern.
5. **P6.5 — smoke the whole flow** in the running binary: upload image → thumbnail → reference in handout → edit with live preview → delete with confirm.
6. **P6.6 — SVG hardening:** run goja-text `sanitize` on SVG uploads before write (`validateSvgUpload` only checks for `<svg>` today).

## 6. Gotchas

- **Nothing from Phases 0–5 is committed yet** — coordinate before rebasing/stashing; there are also pre-existing uncommitted edits to `ContextDiagramPanel.stories/widget.tsx` that are not ours.
- `pnpm --dir web typecheck` fails with a **pre-existing** error (`ContextVisualizerPage.tsx:69` missing `styleSet`) — unrelated to this work.
- gopls may complain `go.work requires go >= 1.26.4` — editor noise; the shell toolchain builds fine.
- Follow GUIDELINES.md to the letter: correct layer, tokens only, no border-radius/box-shadow, `data-rag-*` on roots, stories for every state, `pnpm biome check --write` before review.
- Over IR, state lives in the URL (`?asset=`, `?page=`) and mutations in server actions — don't add client state to widgets beyond what the adapters already own (SearchField value, MarkdownEditor value).
- Keep the diary: append a step to `reference/01-investigation-diary.md` per work session (`/diary` format), check off tasks with `docmgr task check --ticket RAGEVAL-CMS-WIDGETS --id N`, and run `docmgr doctor --ticket RAGEVAL-CMS-WIDGETS` before handing off.
