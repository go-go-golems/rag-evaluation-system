---
Title: CMS widget set analysis, design, and implementation guide
Ticket: RAGEVAL-CMS-WIDGETS
Status: active
Topics:
    - design-system
    - frontend
    - storybook
    - cms
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/course-material-service.js
      Note: 'File-backed CMS core: upload validation'
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/pages/admin-common.js
      Note: Upload panels + material tables (delete without confirm — D-9 motivation)
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js
      Note: Existing IR-authored admin CMS page the new widgets upgrade
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/server.js
      Note: Widget pages + implemented /api/widget/actions/:name + form-post routes + /course-assets
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/webapp/package.json
      Note: npm pin 0.1.16 + dev:rag-local source-alias iteration path
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/xgoja.yaml
      Note: 'Consumer buildspec: rag-widget-site provider (v0.1.2)'
    - Path: internal/api/handlers.go
      Note: Go route table media/articles routes join
    - Path: internal/db/db.go
      Note: SQLite schema; no media/article tables today
    - Path: packages/rag-evaluation-site/GUIDELINES.md
      Note: Design-system constitution all proposals conform to
    - Path: packages/rag-evaluation-site/src/components/molecules/FileDropZone/FileDropZone.tsx
      Note: Existing upload primitive MediaLibraryPanel reuses
    - Path: packages/rag-evaluation-site/src/components/molecules/MarkdownArticle/MarkdownArticle.tsx
      Note: Hand-rolled markdown renderer; sanitization gap (D-5)
    - Path: packages/rag-evaluation-site/src/components/organisms/CourseStudioShell/CourseStudioShell.tsx
      Note: Shell template for CmsShell (D-7)
    - Path: packages/rag-evaluation-site/src/components/organisms/FormPanel/FormPanel.tsx
      Note: Form chrome + status vocabulary ArticleEditorPanel builds on
    - Path: packages/rag-evaluation-site/src/components/organisms/HandoutDocumentShell/HandoutDocumentShell.tsx
      Note: Read-side CMS shell prototype (list + toolbar + article)
    - Path: packages/rag-evaluation-site/src/components/organisms/RichArticle/RichArticle.tsx
      Note: Block-based article renderer; canonical CMS content path
    - Path: packages/rag-evaluation-site/src/context/types.ts
      Note: ArticleBlock union and ContextHandoutDocument DTOs the CMS model extends
    - Path: packages/rag-evaluation-site/src/theme.css
      Note: Token inventory (colors
    - Path: packages/rag-evaluation-site/src/widgets/actions.ts
      Note: Server action posts to unimplemented /api/widget/actions (D-6)
    - Path: packages/rag-evaluation-site/src/widgets/ir.ts
      Note: Widget IR node and action model for phase 6
    - Path: web/src/services/api.ts
      Note: RTK Query surface the CMS endpoints extend
ExternalSources: []
Summary: 'Intern-ready analysis of the rag-evaluation-site design system and a full design + phased implementation guide for a CMS widget set (images, articles, media library, editing) in the strict Classic-Mac visual language. Revised 2026-07-03: delivery is Widget-DSL-first into go-go-course (widget server actions + file-backed course/ storage), replacing the original Go-REST/RTK-Query plan (decisions D-8/D-9).'
LastUpdated: 2026-07-03T17:45:00-04:00
WhatFor: Blueprint for adding CMS atoms/molecules/organisms (media thumbnails, tags, dialogs, media library, article editor/list) plus the web containers and Go backend endpoints they need.
WhenToUse: Read before building any CMS/content-management UI in this repo, or as an onboarding tour of the design-system architecture.
---



# CMS widget set — analysis, design, and implementation guide

## 1. Executive summary

The `rag-evaluation-system` repo already contains a strict, publishable design system (`packages/rag-evaluation-site`, npm: `@go-go-golems/rag-evaluation-site`) with a five-layer architecture (foundation → atoms → layout → molecules → organisms → Widget IR), a hard-edged "Classic Mac / terminal" visual language (1px black borders, zero border-radius, zero shadows, 10–13px tokenized typography), and two Storybooks (package on port 6007, app on 6006). Its course/handout components — `MarkdownArticle`, `RichArticle`, `HandoutDocumentShell`, `DocumentListPanel`, `DocumentPreviewToolbar`, `FileDropZone` — already form the **read side** of a CMS.

What is missing for a real CMS (images, articles, media) is the **write side** and the **media primitives**:

1. No image/media atom — there is no thumbnail component with loading/broken/placeholder states (the `HandoutDocumentShell` "Illustrated" story visibly renders a broken `<img>` today).
2. No taxonomy atoms (tags/chips), no content-lifecycle status (draft/published/scheduled/archived).
3. No dialog/modal primitive, no pagination, no breadcrumbs, no empty-state molecule, no upload-progress list.
4. No editing organisms (markdown editor with preview, article metadata form, asset picker).
5. No backend media storage (SQLite has no blobs/uploads/media table), no upload endpoint, and the front-end `server` action posts to `/api/widget/actions/{name}` which the Go server does not implement.
6. `MarkdownArticle` passes `href`/`src` through unsanitized — acceptable for trusted fixtures, not for CMS authors.

This document proposes a CMS widget set of **4 new atoms** (`MediaThumb`, `Tag`, `ContentStatusBadge`, `MeterBar`), **2 new layout primitives** (`TileGrid`, `DialogShell`), **8 new molecules** (`AssetTile`, `TagListInput`, `Breadcrumbs`, `Pagination`, `SearchField`, `EmptyState`, `UploadQueueList`, `MarkdownEditor`), and **6 new organisms** (`MediaLibraryPanel`, `AssetPickerDialog`, `ArticleListPanel`, `ArticleEditorPanel`, `AssetDetailPanel`, `CmsShell`), plus one extension to the `ArticleBlock` union (`gallery`). All are presentational (API-free), tokenized, story-covered, and `data-rag-*`-tagged, per the package guidelines.

> **Revision (2026-07-03, after Phases 0–4 landed):** the delivery target changed. The CMS consumer is **`go-go-course`** — an xgoja-generated binary that embeds the WidgetRenderer SPA and authors all pages as Widget IR from JavaScript. It **already implements** the widget server-action endpoint (`/api/widget/actions/:name`), file-backed storage under `course/{slides,handouts,media}`, base64 file-upload serialization with magic-byte validation, and an admin CMS page composed from IR. Consequently the original Phase 5 (Go REST API + RTK Query + `web/` containers) is **dropped**; integration flows through the Widget DSL instead. Sections 6.7–6.9, the flows in §8, and the plan in §9 are rewritten below; decision D-6 is superseded by D-8/D-9. Sections describing `web/`'s RTK Query layer (§4.7) remain as accurate background about this repo, but are no longer on the CMS critical path.

## 2. Problem statement and scope

**Problem.** We want to author and manage content — articles with embedded images and diagrams, plus a library of uploaded media — inside this product, using widgets that are indistinguishable in style from the existing design system. Today the system can *render* article-like content beautifully but cannot *manage* it: there is no way to upload an image, browse assets, tag content, edit an article, or move it through a draft→published lifecycle.

**In scope:**

- A complete component inventory and architecture tour (for onboarding).
- Gap analysis: what a CMS needs vs. what exists.
- Full API design (TypeScript props + DTOs) for the new atoms/molecules/organisms.
- Pseudocode for the key flows (upload, article editing, asset picking).
- Backend API + storage sketch (Go, SQLite) and web container wiring (RTK Query).
- Phased, file-level implementation plan and test strategy.

**Out of scope (explicitly):**

- Multi-user auth/permissions, versioning/rollback, localization, publishing pipelines to external sites.
- WYSIWYG contenteditable editing — the editor is markdown-first with live preview (see Decision D-1).
- Actually implementing the components (this is the guide; implementation is follow-up tickets).

## 3. How to read this document (intern orientation)

If you are new to the repo, read in this order:

1. Section 4 (current state) with the referenced files open — every claim has a `path:line` anchor.
2. Section 5 (gaps) to see *why* each proposed piece exists.
3. Section 6 (proposed design) for the APIs; Section 8 for runtime flows.
4. Section 9 when you start implementing — it is ordered by dependency.

**Glossary** (terms used throughout):

- **Design system package**: `packages/rag-evaluation-site`, published to npm, no app dependencies. "Package" alone means this.
- **Web app**: `web/`, the Vite/React/RTK-Query application that consumes the package and talks to the Go backend.
- **Layer**: one of foundation / atoms / layout / molecules / organisms / widgets, in strict import order (never import upward).
- **Token**: a CSS custom property in `src/theme.css` (`--rag-color-*`, `--mac-*`, `--rag-font-role-*`).
- **Widget IR**: a JSON tree (`WidgetNode`) describing UI semantically; rendered by `WidgetRenderer` through a registry. Produced by hand, by the Go backend, or by the Goja JS DSL.
- **DTO**: plain JSON-compatible data shape passed as props (e.g. `ContextHandoutDocument`). Organisms take DTOs + callbacks; containers decide where data comes from.
- **CMS**: here, a content-management surface for *articles* (block-based documents) and *assets* (uploaded images/files), with taxonomy (tags) and lifecycle status.

## 4. Current-state architecture (evidence-based)

### 4.1 Repository layout

```
rag-evaluation-system/
├── packages/rag-evaluation-site/   # design-system package (npm, API-free)
│   ├── src/theme.css               # all tokens
│   ├── src/styles.css              # package stylesheet entry
│   ├── src/components/{foundation,atoms,layout,molecules,organisms}/
│   ├── src/widgets/                # Widget IR: ir.ts, WidgetRenderer.tsx, registries
│   ├── src/context/                # domain types + palette/style-set model
│   ├── src/app/                    # default SPA shell
│   └── .storybook/                 # package Storybook (port 6007)
├── web/                            # app shell: RTK Query, store, pages (port 6006 storybook)
├── cmd/rag-eval/                   # Go CLI: source|chunk|document|embedding|search|workflow|serve
├── internal/api/                   # HTTP handlers (net/http, Go 1.22 patterns)
├── internal/db/                    # SQLite schema + migrations
├── pkg/widgetdsl/                  # Go→JS Widget IR DSL (ui.dsl, data.dsl, context_window.dsl, course.dsl)
├── pkg/defaultspa/, internal/web/  # embedded SPAs (go:embed)
└── ttmp/                           # docmgr ticket workspaces
```

Workspace: `pnpm-workspace.yaml` lists `web` and `packages/rag-evaluation-site`. Biome v2 owns formatting (tabs, double quotes, 100 cols — `AGENTS.md`).

### 4.2 The layer model and its rules

`packages/rag-evaluation-site/GUIDELINES.md` is the constitution. The non-negotiables (GUIDELINES.md:7–17):

- Correct layer per component; **never import upward**.
- Typography only via `Text`/`Caption`/`CodeText`/`StatusText` or `--rag-font-role-*` tokens; no ad-hoc `font-size`/`font-weight`.
- Colors only via `--mac-*`/`--rag-*` tokens; **CSS Modules only** — no CSS-in-JS, no `Box`, no style-prop systems.
- Package components are **API-free**: no RTK Query, no store, no router (GUIDELINES.md:13, 125–132).
- **Storybook required** for every public component; **`data-rag-*` identity attributes** required (GUIDELINES.md:14, 16).
- React first, Widget IR later (GUIDELINES.md:15): stabilize the React API before adding IR nodes.

Layer flow (GUIDELINES.md:21–29):

```
theme.css ──▶ foundation ──▶ atoms ──▶ layout ──▶ molecules ──▶ organisms ──▶ widgets/WidgetRenderer
 (tokens)     (Text, ...)   (Button)  (Panel)    (DataTable)   (panels)      (JSON IR)
```

Layer ownership, paraphrased from GUIDELINES.md:47–132:

- **Foundation** (`src/components/foundation/`): `Text` (sizes body/compact/metadata/label/metric; tones primary/muted/inverse/accent/success/warning/danger/inherit), `Caption`, `CodeText`, `StatusText` (statuses pending/ready/running/succeeded/done/partial/warning/failed/error/canceled — StatusText.tsx:4–14), `Divider`, `VisuallyHidden`.
- **Atoms** (`src/components/atoms/`): `Button` (variants default|primary, sizes normal|compact, `selected` — Button.tsx:7–12), `TextInput`, `TextareaInput`, `SelectInput`, `CheckboxRow`, `IconButton`, `ErrorCallout`, `UploadGlyph`, `ContextStudioNavIcon`, plus context-domain atoms (`ContextStyleSwatch`, `AnnotationBadge`, `TranscriptRoleBadge`).
- **Layout** (`src/components/layout/`): `AppShell`, `SidebarShell` (grid `var(--rag-sidebar-width) minmax(0,1fr)`), `Panel` (dark uppercase mono title bar — Panel.module.css), `Stack` (gaps xs2/sm4/md8/lg12/xl16), `Inline`, `DashboardGrid` (recipes searchWorkbench|corpusExplorer|twoColumn), `SplitPane` (ratios balanced|leftNarrow|rightNarrow|course|sidebar), `SectionBlock`, `SlideShell`, `ScrollRegion`, `TabList`, `FormRow` (label/control/hint/error/success/counter, orientations inline|stacked — FormRow.tsx:6–16).
- **Molecules** (`src/components/molecules/`): data/content display patterns with typed props — `DataTable<T>`, `MetadataGrid`, `KeyValueStrip`, `SidebarNav`, `AppNav`, `CheckList`, `StepList`, `KeyPointList`, `FigureBlock`, `PersonSummary`, `MarkdownArticle`, `DocumentListPanel`, `DocumentPreviewToolbar`, `FileDropZone`, plus context/transcript molecules.
- **Organisms** (`src/components/organisms/`): DTO-shaped feature panels — `HandoutDocumentShell`, `RichArticle`, `CourseStudioShell`, `CourseLessonPanel`, `CourseSlidePanel`, `FormPanel`, `TranscriptReaderPanel`, `ContextDiagramPanel`, etc.

### 4.3 Tokens and typography (`src/theme.css`, 54 lines)

Canonical tokens (theme.css:3–16): `--rag-font-sans` (Inter stack), `--rag-font-mono` (SFMono/Consolas), and 11 colors — bg `#f6f7f8`, surface `#ffffff`, surface-muted `#f0f2f4`, text `#1d232a`, text-muted `#68727d`, border `#d8dde3`, **border-strong `#000000`**, accent `#2457d6`, success `#14883e`, warning `#a56a00`, danger `#bd2c2c`.

A compatibility bridge (theme.css:18–40) maps the legacy `--mac-*` names onto these: `--mac-border` → border-strong (pure black), `--mac-bg-dark` → `#000000` (the inverted title-bar/selection color), `--mac-accent-2` → danger, `--mac-green`/`--mac-amber` → success/warning. Component CSS overwhelmingly uses the `--mac-*` names.

Nine font roles (theme.css:45–53) are the entire typographic vocabulary:

| Role | Value | Use |
|---|---|---|
| `--rag-font-role-body` | 400 12px/1.45 sans | UI prose |
| `--rag-font-role-compact` | 400 11px/1.4 sans | dense rows, nav |
| `--rag-font-role-metadata` | 400 11px/1.35 mono | captions, meta, buttons |
| `--rag-font-role-label` | 700 11px/1.2 mono | uppercase section labels |
| `--rag-font-role-metric` | 700 12px/1.25 mono | numeric values |
| `--rag-font-role-code` | 400 11px/1.45 mono | ids, paths, statuses |
| `--rag-font-role-display` | 700 22px/1.15 sans | article H1 |
| `--rag-font-role-heading` | 700 16px/1.25 sans | article H2 / section labels |
| `--rag-font-role-readable-body` | 400 13px/1.55 sans | article prose |

**The visual signature** (observed across all module CSS and confirmed in screenshots under `sources/screenshots/`): every border is `1px solid var(--mac-border)` (2px only for the SlideShell header rule); **no border-radius and no box-shadow exist anywhere in the package**; focus is `outline: 1px solid var(--mac-accent)`; selection/active states invert to `--mac-bg-dark` background with `--mac-text-inv` text; panel title bars are black with white bold 11px uppercase mono text. New CMS components must reproduce exactly this language.

### 4.4 Component conventions

Uniform folder layout per public component: `Component.tsx`, `Component.module.css`, `Component.stories.tsx`, `index.ts`, and usually `Component.widget.tsx` + `Component.widget.yaml` (the Widget IR adapter + manifest). Conventions observed across all layers:

- Props interfaces `extend` the matching DOM attributes type and spread `...rest`; `Omit<…, "title">` where a prop collides.
- className merge idiom: `[styles.root, …conditional, className ?? ""].filter(Boolean).join(" ")`.
- Identity attributes are namespaced per layer: `data-rag-foundation` / `data-rag-atom` / `data-rag-layout` / `data-rag-molecule` / `data-rag-organism`, value = PascalCase name. Known inconsistencies to *not* copy: `DataTable`, `MetadataGrid`, `AppNav` use `data-rag-component`; most foundation primitives emit none.
- Secondary state hooks as `data-*`: `data-active`, `data-disabled`, `data-state`, `data-orientation`, etc.
- Storybook titles (GUIDELINES.md:140–147): `Design System/{Foundation|Atoms|Layout}/<Name>`, `Component Library/{Molecules|Organisms}/<Name>`, `Widget IR/Renderer`. Required story states (GUIDELINES.md:149–157): default/populated, empty, overflow/dense, selected/active, disabled, error/warning, alternate direction.
- Package Storybook: `.storybook/main.ts` globs `../src/**/*.stories.@(ts|tsx)` and sets readable CSS-module class names (`[name]_[local]`); `preview.ts` imports `../src/styles.css`, layout "padded".

### 4.5 The proto-CMS: existing content components

These five components are the seed of the CMS and set the patterns everything new must follow.

**`MarkdownArticle`** (molecule, `molecules/MarkdownArticle/MarkdownArticle.tsx`). Props: `{ source: string }` (:4–6). A ~260-line hand-rolled markdown renderer — **no external library, no sanitization**. It handles fenced code (with `data-language`), `#/##/###` headings mapped to display/heading/label font roles, blockquotes, GFM tables, task lists (☑/☐), ordered/unordered lists, `---` rules, standalone images (`![alt](src "title")` → `<figure>` + lazy `<img>` + `<figcaption>`, :52–56, 96–114), and inline `**bold**`, `` `code` ``, `[text](url)` (:10–39). Inline `href` (:29) and image `src` pass through unescaped. CSS caps width at 720px and maps headings/prose to the article font roles (MarkdownArticle.module.css:18–91).

**`RichArticle`** (organism, `organisms/RichArticle/RichArticle.tsx`). Props: `{ blocks: ArticleBlock[], styleSet?: ContextStyleSet }` (:8–11). Block-based, not parsed. The block union (`src/context/types.ts:170–192`):

```ts
export interface ArticleMarkdownBlock      { kind: "markdown";       id: string; source: string }
export interface ArticleContextWindowBlock { kind: "context-window"; id: string; snapshot: ContextWindowSnapshot; view?: ContextDiagramView; caption?: string }
export interface ArticleImageBlock         { kind: "image";          id: string; src: string; alt: string; caption?: string }
export type ArticleBlock = ArticleMarkdownBlock | ArticleContextWindowBlock | ArticleImageBlock;
```

`markdown` blocks delegate to `MarkdownArticle`; `context-window` blocks render an inline `ContextDiagramPanel`; `image` blocks render a borderless centered figure. Each emits `data-rag-article-block="<kind>"`.

**`HandoutDocumentShell`** (organism, `organisms/HandoutDocumentShell/HandoutDocumentShell.tsx:8–20`). A two-column reading shell: left `DocumentListPanel`, right `DocumentPreviewToolbar` + content. Content precedence (:75–83): `blocks` → `RichArticle`; else format contains "markdown" → `MarkdownArticle(body)`; else description placeholder. Document DTO `ContextHandoutDocument` (types.ts:194–205): `{ id, title, file, format, size?, description, body, blocks?, downloadHref?, printHref? }`. Grid `268px minmax(0,1fr)`, collapsing at 900px.

**`DocumentListPanel`** (molecule, `molecules/DocumentListPanel/DocumentListPanel.tsx:6–27`). `role="listbox"` of documents with icon-by-format glyphs (pdf `▤`, json `{ }`, markdown `¶`, else `□`), selection inverts to black, footer "⤓ Download all (.zip)" primary button.

**`FileDropZone`** (molecule, `molecules/FileDropZone/FileDropZone.tsx:14–25`). Full drag/drop + click + keyboard file input with `onFilesSelected(files: File[])`, `accept`, `multiple`, `disabled`, `active`; dashed 1px border; default `UploadGlyph` icon. Wrapped by the `ContextUploadDropArea` organism for the context-JSON use case.

Also relevant: `DataTable<T>` (render-prop columns, selected-row inversion — molecules/DataTable/DataTable.tsx:4–20), `FormPanel` (organism; form chrome with status idle|saving|success|error, aria-live status region — organisms/FormPanel/FormPanel.tsx:7–20), `FormRow`, `MetadataGrid` (dl/dt/dd with per-item copy), `SidebarNav` (sectioned nav), and `CourseStudioShell` (SidebarShell + SidebarNav + header, sidebarWidth 188 — the template for `CmsShell`).

### 4.6 Widget IR and the action model (`src/widgets/`)

The package ships a JSON UI representation so the Go backend (or Goja scripts) can author pages:

- `WidgetNode = TextNode | ElementNode | ComponentNode` (ir.ts:40–122); `component` nodes name one of ~65 `RagWidgetType` strings and carry JSON props.
- `ActionSpec` (ir.ts:126–161): `navigate | download | server | event | copy`. Default dispatch (actions.ts:24–96): `copy` → clipboard; `event` → `CustomEvent` on `window` (special-cases print/fullscreen); `navigate` → pushState; `download` → temp anchor; **`server` → `POST /api/widget/actions/{name}`** with `{payload, context}` — note the Go backend does not implement this route yet.
- Registries (registry.ts, defaultRegistry.ts) group adapters by module: `ui.dsl` (30 primitives), `data.dsl` (DataTable), `context_window.dsl`, `course.dsl` (MarkdownArticle, RichArticle, HandoutDocumentShell, DocumentListPanel, DocumentPreviewToolbar, Course*). `WidgetRenderer` renders unknown types as an `ErrorCallout` fallback.
- `useWidgetPage(url)` (hooks/useWidgetPage.ts:23–71) fetches `{ id, title, root }`; the web app instead uses RTK Query `getDslPage` → `GET /api/v1/dsl/pages/{id}`.
- Go side: `pkg/widgetdsl` mirrors the four modules as goja `require()` modules (`ui.dsl`, `data.dsl`, `context_window.dsl`, `course.dsl`) with camelCase helpers, `action.*`/`cell.*` builders, palette style helpers, and 8 macro-style recipes; the `rag-widget-site` xgoja provider (`pkg/xgoja/providers/widgetsite/provider.go`) exposes them plus generated TS typings to external generated binaries. The DSL is currently dormant in-repo (not selected in `cmd/rag-eval/xgoja.yaml`, unused by jsverbs); the only served page is hardcoded in `internal/api/dsl_handlers.go:5–88`. Full line-anchored reference and `cms.dsl` extension notes: `reference/02-goja-widget-dsl-reference-and-cms-extension-notes.md` in this ticket.

### 4.7 Web app data layer (`web/src/`)

- `App.tsx:11–19,27–97`: no router — a `useState` view switcher (search, corpus, workflows, pipeline, embeddings, evaluation, dsl) inside the package `AppShell`+`AppNav`. Cross-view navigation via window `CustomEvent`s (`rag:navigate-to-chunk`, `rag:navigate-to-workflows`).
- `store/index.ts:5–10`: one reducer — `ragApi` (RTK Query). No custom slices.
- `services/api.ts:303–494`: `baseQuery: fetchBaseQuery({ baseUrl: "/api/v1" })`; endpoints for sources, documents, chunks, embeddings, search (bm25/vector/hybrid), corpus, workflows, artifacts, and `getDslPage`. Tag types: Sources, Documents, Chunks, Strategies, Embeddings, Corpus, Workflows, Artifacts (:306–315).
- `web/src/storybook/MockApiProvider.tsx` monkey-patches `fetch` with canned `/api/v1` payloads so web-level stories run without a backend — the pattern to reuse for CMS container stories.

### 4.8 Go backend

- `cmd/rag-eval` (Cobra) with `serve` starting `net/http` + `ServeMux` (serve/server.go:19–63); all routes in `internal/api/handlers.go:38–89` under `/api/v1` (health, dsl/pages/{id}, sources, documents, chunks, chunking-strategies, embeddings, search, workflows, artifacts, corpus). SPA embedded via `internal/web/spa.go` (`go:embed all:dist`).
- Storage: SQLite (`internal/db/db.go:70–230`): `sources`, `documents` (text-only content: `raw_content`, `content_text`, `content_html`, `metadata_json`, `status` — db.go:81–98), `chunks`, `chunk_embeddings`, `chunk_enrichments`, `document_processing_artifacts`, `search_indexes`, `eval_*`. **No media/blob/upload/article tables.**

## 5. Gap analysis

What a minimal-but-real CMS needs, versus what exists. "Have" cites the reusable piece; "Gap" is what must be built.

| CMS capability | Have | Gap |
|---|---|---|
| Render an article (markdown, images, diagrams) | `RichArticle` + `ArticleBlock` union; `MarkdownArticle` | `gallery` block kind; sanitization of `href`/`src` (MarkdownArticle.tsx:29, 100–104) |
| Read-side document shell | `HandoutDocumentShell`, `DocumentListPanel`, `DocumentPreviewToolbar` | Generalize labels ("Handout" default title) — props already allow it |
| Image/media primitive | none — raw `<img>` inside article figures only | `MediaThumb` atom: aspect-ratio frame, loading/broken/placeholder states, flat bordered style |
| Media library (browse/search/upload assets) | `FileDropZone`, `Panel`, `DashboardGrid` | `TileGrid` layout, `AssetTile`, `MediaLibraryPanel`, `UploadQueueList`, `Pagination`, `SearchField`, `EmptyState` |
| Pick an asset from within a form | nothing modal exists in the package | `DialogShell` layout primitive + `AssetPickerDialog`; `ImagePickerField` control pattern |
| Taxonomy (tags) | `AnnotationBadge` is context-styleSet-specific | Generic `Tag` atom + `TagListInput` molecule |
| Content lifecycle | `StatusText` statuses are workflow-shaped (running/failed/…) (StatusText.tsx:4–14) | `ContentStatusBadge` atom mapping draft/published/scheduled/archived onto existing tones |
| Article listing/management | `DataTable<T>` with selection | `ArticleListPanel` organism (columns, status cells, row actions, pagination, empty state) |
| Article editing | `FormPanel`, `FormRow`, `TextareaInput` | `MarkdownEditor` molecule (toolbar + textarea + counter), `ArticleEditorPanel` organism (metadata + editor + live preview via `SplitPane`) |
| Upload progress | nothing | `MeterBar` atom + `UploadQueueList` molecule |
| CMS navigation shell | `CourseStudioShell` (course-specific nav) | `CmsShell` organism (same anatomy, CMS sections) — or generalize CourseStudioShell |
| Backend storage/APIs for media & articles | `documents` table (text only); no upload routes | `media_assets` + `articles` tables; `POST /api/v1/media` (multipart), CRUD routes; static file serving for blobs |
| Server actions from widgets | `ActionSpec kind:"server"` posts to `/api/widget/actions/{name}` (actions.ts:82–95) | Implement the route in Go (or route CMS mutations through RTK Query instead — see D-6) |

Also worth fixing while in the area (consistency debt, not blockers): `data-rag-component` outliers (DataTable/MetadataGrid/AppNav), `Caption` danger tone using `--mac-red` while everything else uses `--mac-accent-2`, and `DataTable`'s raw `font-family: var(--font-mono); font-size: 11px` instead of a font role.

## 6. Proposed design

### 6.1 Design principles (how new components stay "in style")

1. **Flat and hairline**: every frame is `1px solid var(--mac-border)`; no radius, no shadow, ever. Dashed borders only for drop targets (matching `FileDropZone`).
2. **Inversion is selection**: selected/active states use `--mac-bg-dark` + `--mac-text-inv`, exactly like `Button[selected]`, `DocumentListPanel` active rows, and `DataTable` selected rows.
3. **Mono for chrome, sans for content**: labels/meta/status in mono roles (`metadata`, `label`, `code`); article/reading text in `readable-body`; UI prose in `body`.
4. **Tokens only**: any new color need (e.g. image placeholder checker) is built from existing tokens or added deliberately to `theme.css` with a comment.
5. **DTO props + callbacks**: organisms accept plain JSON-compatible data and `on*` callbacks; containers in `web/` own fetching/mutation.
6. **Stories are the spec**: each component ships stories for default, empty, dense/overflow, selected, disabled, and error states with the correct title prefix.
7. **`data-rag-*` everywhere**, layer-correct, plus block/state attributes (`data-rag-asset-id`, `data-state`, `data-active`).

### 6.2 CMS domain model (new DTOs, in `src/context/types.ts` or a new `src/cms/types.ts`)

```ts
export type CmsContentStatus = "draft" | "published" | "scheduled" | "archived";

export interface CmsAsset {
	id: string;
	kind: "image" | "file";
	title: string;
	filename: string;            // "hero-diagram.png"
	mime: string;                // "image/png"
	size: number;                // bytes
	src: string;                 // canonical URL (served by backend)
	thumbSrc?: string;           // optional downscaled URL
	width?: number;              // images only
	height?: number;
	alt?: string;
	tags: string[];
	status: CmsContentStatus;
	createdAt: string;           // ISO 8601
	updatedAt: string;
}

export interface CmsArticleSummary {
	id: string;
	slug: string;
	title: string;
	status: CmsContentStatus;
	author?: string;
	tags: string[];
	excerpt?: string;
	updatedAt: string;
}

export interface CmsArticleDetail extends CmsArticleSummary {
	blocks: ArticleBlock[];      // canonical content (see D-1)
	coverAssetId?: string;
}

// ArticleBlock union gains one member (types.ts:170–192):
export interface ArticleGalleryBlock {
	kind: "gallery";
	id: string;
	images: Array<{ src: string; alt: string; caption?: string }>;
	columns?: 2 | 3 | 4;         // default 3
}
```

All fields are JSON-compatible on purpose — these DTOs flow unchanged from Go handlers through RTK Query into organism props, and later into Widget IR nodes.

### 6.3 New atoms (`src/components/atoms/`)

**`MediaThumb`** — the missing image primitive. Anatomy: a fixed-aspect frame that letterboxes the image and owns loading/error/empty presentation.

```ts
export interface MediaThumbProps extends Omit<HTMLAttributes<HTMLDivElement>, "onError"> {
	src?: string;                 // absent → placeholder state
	alt?: string;
	aspect?: "square" | "wide" | "natural";   // default "square" (1:1); wide = 16/10
	fit?: "cover" | "contain";                // default "cover"
	frame?: "bordered" | "none";              // default "bordered" (1px --mac-border)
	fallbackGlyph?: ReactNode;                // default "▨"
	selected?: boolean;
}
```

Behavior pseudocode:

```
state = src ? "loading" : "empty"
render frame div [data-rag-atom=MediaThumb] [data-state=state] [data-active=selected]
  if src: <img loading="lazy" onLoad→state=loaded onError→state=broken>
  if state in {empty, broken}: centered fallbackGlyph + Caption ("no image" / filename)
  while loading: checkerboard background built from repeating-linear-gradient of
                 --mac-surface-2 / --mac-surface (same technique as .pattern_* classes
                 in ContextStyleSwatch.module.css — reuse, don't invent)
```

CSS notes: `image-rendering: auto` (unlike the pixelated swatches — photos must not pixelate); selected state adds `outline: 1px solid var(--mac-accent)` inside the black frame, mirroring `ContextStyleSwatch.selected`.

**`Tag`** — taxonomy chip. Deliberately distinct from `AnnotationBadge` (which requires a `ContextVisualStyle`).

```ts
export interface TagProps extends HTMLAttributes<HTMLSpanElement> {
	label: string;
	selected?: boolean;
	onRemove?: () => void;        // renders a compact "×" IconButton when present
	disabled?: boolean;
}
```

Style: inline-flex, `1px solid var(--mac-border)`, `font: var(--rag-font-role-metadata)`, padding `1px 5px`, selected inverts to black/white. Interactive removal reuses `IconButton`.

**`ContentStatusBadge`** — content lifecycle marker (see D-2 for why this is not a `StatusText` extension).

```ts
export type ContentStatus = "draft" | "published" | "scheduled" | "archived";
export interface ContentStatusBadgeProps extends HTMLAttributes<HTMLSpanElement> {
	status: ContentStatus;
	icon?: boolean;               // default true: ◌ draft, ● published, ◔ scheduled, ▣ archived
}
```

Tone mapping (all existing tokens): draft → `--mac-text-dim`; published → `--mac-green`; scheduled → `--mac-accent`; archived → `--mac-text-dim` + line-through (the `canceled` treatment from StatusText.module.css). Font: `--rag-font-role-label`, uppercase, bordered like `TranscriptRoleBadge`.

**`MeterBar`** — determinate progress for uploads.

```ts
export interface MeterBarProps extends HTMLAttributes<HTMLDivElement> {
	value: number;                // 0..1
	tone?: "accent" | "success" | "danger";   // default "accent"
	label?: ReactNode;            // optional right-aligned metric text
}
```

Anatomy: outer track `1px solid var(--mac-border)` on `--mac-surface`, inner fill solid tone color, height 10px, no animation except width transition. Label uses `--rag-font-role-metric`. This is visually a simplified sibling of `ContextBudgetBar` without the styleSet machinery.

### 6.4 New layout primitives (`src/components/layout/`)

**`TileGrid`** — generic tile layout for asset grids (layout must not know "asset"):

```ts
export interface TileGridProps extends HTMLAttributes<HTMLDivElement> {
	minTileWidth?: number;        // default 160 (px)
	gap?: "sm" | "md";            // 4 | 8 px, matching Stack scale
	children?: ReactNode;
}
```

CSS: `display: grid; grid-template-columns: repeat(auto-fill, minmax(var(--rag-tile-min), 1fr))`, min width plumbed via inline CSS var like `SidebarShell` does with `--rag-sidebar-width` (SidebarShell.tsx:38 pattern).

**`DialogShell`** — the package's first modal primitive (see D-3).

```ts
export interface DialogShellProps extends Omit<HTMLAttributes<HTMLDialogElement>, "title"> {
	open: boolean;
	onClose: () => void;
	title: ReactNode;             // dark Panel-style title bar
	actions?: ReactNode;          // title-bar right slot (e.g. close IconButton is implicit)
	footer?: ReactNode;           // bottom action row
	size?: "sm" | "md" | "lg";    // 420 | 640 | 920 px max-width
	children?: ReactNode;
}
```

Implementation: native `<dialog>` element driven by a `useEffect` calling `showModal()/close()`; `onClose` wired to the `cancel` and `close` events (Esc works for free). Styling: `border: 1px solid var(--mac-border); background: var(--mac-surface); padding: 0;` — the header reuses the exact Panel header recipe (black bar, white bold 11px uppercase mono — Panel.module.css `.header`); `::backdrop { background: rgb(0 0 0 / 0.25); }` — flat dimming, no blur. Emits `data-rag-layout="DialogShell"`.

### 6.5 New molecules (`src/components/molecules/`)

**`AssetTile`** — one asset in a grid. Composes `MediaThumb` + text rows.

```ts
export interface AssetTileProps extends Omit<HTMLAttributes<HTMLDivElement>, "onSelect"> {
	asset: CmsAsset;
	selected?: boolean;
	onSelect?: (id: string) => void;
	onOpen?: (id: string) => void;      // double-click / Enter
	footerSlot?: ReactNode;             // e.g. status badge, tag count
}
```

Anatomy: `role="option"`-style button wrapping `MediaThumb` (square), title in `--rag-font-role-compact` (truncated), meta line `PNG · 214 KB` in `--rag-font-role-metadata` muted. Selection = 1px accent outline + inverted title row, consistent with `DocumentListPanel`'s active item. `data-rag-molecule="AssetTile"`, `data-rag-asset-id`.

**`TagListInput`** — view + edit tags.

```ts
export interface TagListInputProps {
	tags: string[];
	onAdd?: (tag: string) => void;
	onRemove?: (tag: string) => void;
	suggestions?: string[];             // datalist-style completion
	placeholder?: string;               // default "add tag…"
	disabled?: boolean;
	className?: string;
}
```

Anatomy: `Inline` wrap of `Tag`(with `onRemove`) + a borderless mini `TextInput` that commits on Enter/comma. No dropdown UI in v1 — a native `<datalist>` keeps it dependency-free.

**`Breadcrumbs`** — folder/collection path.

```ts
export interface BreadcrumbItem { id: string; label: ReactNode }
export interface BreadcrumbsProps {
	items: BreadcrumbItem[];
	onNavigate?: (id: string) => void;  // last item is current, unclickable
	className?: string;
}
```

Mono `--rag-font-role-metadata`, `/` separators, current item bold. `<nav aria-label="Breadcrumbs">`.

**`Pagination`** — compact pager.

```ts
export interface PaginationProps {
	page: number;                        // 1-based
	pageCount: number;
	onPageChange: (page: number) => void;
	pageSize?: number; totalItems?: number;  // optional "1–24 of 210" caption
	className?: string;
}
```

Anatomy: `‹ prev` / `next ›` compact `Button`s + `Caption` "page 3 / 9". No numbered buttons in v1 (dense style favors terseness).

**`SearchField`** — the standard filter box.

```ts
export interface SearchFieldProps extends Omit<InputHTMLAttributes<HTMLInputElement>, "onChange"> {
	value: string;
	onValueChange: (value: string) => void;
	onSubmit?: (value: string) => void;
	onClear?: () => void;               // shows "×" IconButton when value non-empty
}
```

Wraps `TextInput` with a leading `⌕` glyph and trailing clear `IconButton`; Enter triggers `onSubmit`.

**`EmptyState`** — standardizes the ad-hoc `emptyMessage` patterns.

```ts
export interface EmptyStateProps extends HTMLAttributes<HTMLDivElement> {
	glyph?: ReactNode;                  // default "□"
	title: ReactNode;
	hint?: ReactNode;
	action?: ReactNode;                 // usually a Button
}
```

Centered, muted, mono hint; dashed 1px border variant off by default (`framed?: boolean`) so it can sit inside panels or stand alone.

**`UploadQueueList`** — in-flight uploads.

```ts
export type UploadItemStatus = "queued" | "uploading" | "done" | "error" | "canceled";
export interface UploadQueueItem {
	id: string; filename: string; size: number;
	progress: number;                   // 0..1
	status: UploadItemStatus;
	error?: string;
}
export interface UploadQueueListProps {
	items: UploadQueueItem[];
	onCancel?: (id: string) => void;
	onRetry?: (id: string) => void;
	onDismiss?: (id: string) => void;
	className?: string;
}
```

Row anatomy: format glyph + `CodeText` filename + `MeterBar` (tone: error → danger, done → success) + `StatusText`-style status + cancel/retry `IconButton`s. Errors render the message in an `ErrorCallout`-toned caption.

**`MarkdownEditor`** — plain-text editing with affordances (no contenteditable; see D-1).

```ts
export interface MarkdownEditorProps {
	value: string;
	onValueChange: (value: string) => void;
	onInsertAsset?: () => void;         // container opens AssetPickerDialog, then inserts markdown
	minRows?: number;                   // default 16
	maxLength?: number;
	disabled?: boolean;
	toolbarSlot?: ReactNode;            // extra buttons
	className?: string;
}
```

Anatomy: toolbar (`Inline` of compact `Button`s: **B**, `code`, link, H2, list, image — each wraps the current selection with markdown syntax via a `wrapSelection(textarea, before, after)` helper) above a `TextareaInput` (mono? No — body font, matching TextareaInput) with character counter in a `FormRow`-style counter slot. Toolbar pseudocode:

```
wrapSelection(el, before, after):
  { start, end } = el.selectionRange
  next = value[0:start] + before + value[start:end] + after + value[end:]
  onValueChange(next); restore selection at start+len(before)
insertImage(): onInsertAsset?.()   // container inserts "![alt](src)" at cursor after picking
```

### 6.6 New organisms (`src/components/organisms/`)

**`MediaLibraryPanel`** — the asset browser.

```ts
export interface MediaLibraryPanelProps extends Omit<HTMLAttributes<HTMLDivElement>, "onSelect"> {
	assets: CmsAsset[];
	selectedAssetIds?: string[];
	selectionMode?: "none" | "single" | "multi";     // default "single"
	onAssetSelect?: (id: string) => void;
	onAssetOpen?: (id: string) => void;
	// toolbar state (controlled by container):
	query?: string; onQueryChange?: (q: string) => void;
	kindFilter?: "all" | "image" | "file"; onKindFilterChange?: (k: string) => void;
	page?: number; pageCount?: number; onPageChange?: (p: number) => void;
	// upload:
	onFilesSelected?: (files: File[]) => void;        // renders FileDropZone strip when set
	uploads?: UploadQueueItem[];
	onUploadCancel?/onUploadRetry?: (id: string) => void;
	emptyMessage?: ReactNode;
	title?: ReactNode;                                 // default "Media"
}
```

Anatomy (all existing + new pieces):

```
┌─ Panel [title=MEDIA, actions=kind TabList] ──────────────────────┐
│ Inline: SearchField ····························· Pagination     │
│ FileDropZone (compact strip, when onFilesSelected)               │
│ UploadQueueList (when uploads non-empty)                         │
│ ScrollRegion                                                     │
│   TileGrid minTileWidth=160                                      │
│     AssetTile × n   (footerSlot = ContentStatusBadge)            │
│   — or EmptyState (glyph ▨, "No assets yet", action: Upload)     │
└──────────────────────────────────────────────────────────────────┘
```

**`AssetPickerDialog`** — `DialogShell size="lg"` wrapping a `MediaLibraryPanel` in `selectionMode="single"` picker mode, with footer `Cancel` / `Use asset` (primary, disabled until selection). Props: the MediaLibraryPanel surface plus `{ open, onClose, onConfirm(asset: CmsAsset) }`.

**`ArticleListPanel`** — content management table.

```ts
export interface ArticleListPanelProps {
	articles: CmsArticleSummary[];
	selectedArticleId?: string;
	onArticleSelect?: (id: string) => void;
	onCreate?: () => void;                             // "＋ New article" panel action
	onRowAction?: (id: string, action: "edit" | "publish" | "archive" | "delete") => void;
	statusFilter?: CmsContentStatus | "all"; onStatusFilterChange?: (s: string) => void;
	query?: string; onQueryChange?: (q: string) => void;
	page?: number; pageCount?: number; onPageChange?: (p: number) => void;
	emptyMessage?: ReactNode;
}
```

Built on `DataTable<CmsArticleSummary>` with columns: title (+ slug as `CodeText` caption), `ContentStatusBadge`, tags (first 2 `Tag`s + "+n"), author, updatedAt (mono), and a row-actions cell of `IconButton`s. Toolbar mirrors MediaLibraryPanel (SearchField + status `SelectInput` + Pagination). Destructive actions (delete/archive) confirm via `DialogShell` (`ConfirmDialog` pattern below).

**`ArticleEditorPanel`** — the write surface.

```ts
export interface ArticleEditorDraft {
	title: string; slug: string; status: CmsContentStatus;
	tags: string[]; excerpt?: string; coverAssetId?: string;
	body: string;                                      // markdown source (v1: single markdown block)
}
export interface ArticleEditorPanelProps {
	draft: ArticleEditorDraft;
	onDraftChange: (draft: ArticleEditorDraft) => void;
	onSave?: () => void; onPublish?: () => void;
	status?: "idle" | "saving" | "success" | "error";  // FormPanel vocabulary (FormPanel.tsx:7)
	statusMessage?: ReactNode;
	onPickCoverAsset?: () => void;                     // container opens AssetPickerDialog
	onInsertAsset?: () => void;                        // ditto, inserts into body at cursor
	preview?: "live" | "hidden";                       // default "live"
	tagSuggestions?: string[];
}
```

Anatomy:

```
┌─ FormPanel [title=ARTICLE, status, submit=Save, actions=Publish] ─┐
│ FormRow label=TITLE   control=TextInput                           │
│ FormRow label=SLUG    control=TextInput (CodeText-styled)         │
│ FormRow label=STATUS  control=SelectInput(draft/published/…)      │
│ FormRow label=TAGS    control=TagListInput                        │
│ FormRow label=COVER   control=MediaThumb(sm) + "Choose…" Button   │
│ SplitPane ratio=balanced divider                                  │
│   left:  MarkdownEditor(body)                                     │
│   right: ScrollRegion > MarkdownArticle(source=body)  ← live      │
└───────────────────────────────────────────────────────────────────┘
```

The live preview is literally the production renderer — zero drift between editing and reading (this is the payoff of markdown-first, D-1).

**`AssetDetailPanel`** — inspect/edit a single asset: large `MediaThumb (aspect="natural", fit="contain")`, `MetadataGrid` (id, mime, dimensions, size, createdAt — with copy buttons), `FormPanel` for title/alt/tags/status, usage list (`DocumentListPanel`-style "used in n articles"), and Delete (confirming) / Download actions in the toolbar (`DocumentPreviewToolbar` reused verbatim).

**`CmsShell`** — navigation chrome. Same anatomy as `CourseStudioShell` (SidebarShell width 188 + header + SidebarNav + footer — CourseStudioShell.tsx:20–59) with CMS sections:

```ts
export const cmsNavSections: SidebarNavSection[] = [
	{ id: "content", label: "Content", items: [
		{ id: "articles", label: "Articles", icon: <ContextStudioNavIcon id="handout" /> },
		{ id: "media",    label: "Media",    icon: <ContextStudioNavIcon id="upload" /> },
	]},
	{ id: "organize", label: "Organize", items: [
		{ id: "tags",     label: "Tags" },
		{ id: "trash",    label: "Archive" },
	]},
];
```

Decision D-7 covers generalizing `CourseStudioShell` instead of duplicating it.

**`ConfirmDialog`** (small organism or molecule): `DialogShell size="sm"` + message + `Cancel`/destructive-primary confirm; the destructive confirm button uses default variant with danger-toned label (`--mac-accent-2`) — do **not** invent a red button variant without a token decision.

### 6.7 The consumer: `go-go-course` (revised — replaces the RTK Query plan)

`go-go-course` (sibling repo, module `github.com/go-go-golems/go-go-course`) is an **xgoja-generated binary** (`xgoja.yaml` name `minitrace-viz`) whose entire UI is Widget IR authored in JavaScript. Its pipeline:

```
rag-evaluation-system (this repo)                         go-go-course
┌────────────────────────────────┐   npm publish   ┌─────────────────────────────────────┐
│ React components + .widget.tsx │ ───────────────▶│ webapp/ (15-line shell around       │
│ adapters + cms.dsl (Go)        │                  │  RagEvaluationSiteApp               │
└────────────────────────────────┘   Go module tag  │  apiBase="/api/widget") → vite build│
                                  ────────────────▶│ assets/ → embedded in xgoja binary  │
                                                    │ server.js + lib/pages/*.js author  │
                                                    │  IR via ui/data/context_window/     │
                                                    │  course .dsl (xgoja.yaml providers) │
                                                    └─────────────────────────────────────┘
```

Key facts (all verified in `go-go-course/cmd/go-go-course/`):

- **Buildspec** (`xgoja.yaml`): providers include `rag-widget-site` (this repo's `pkg/xgoja/providers/widgetsite`, currently pinned `v0.1.2`); runtime selects `ui.dsl`, `data.dsl`, `context_window.dsl`, `course.dsl`; the built SPA is embedded as asset module `minitrace-widget-spa` and served via `app.spaFromAssetsModule("/", …)` (server.js:619–625). `webapp/package.json` pins `@go-go-golems/rag-evaluation-site: 0.1.16` (npm).
- **Pages**: `GET /api/widget/pages/:id` builds IR through `createCoursePages(...)` → `lib/pages/*.js` (course, slides, handouts, sessions, settings, upload, **admin-course-cms**) (server.js:128–153).
- **Server actions exist**: `POST /api/widget/actions/:name` (server.js:383–456) handles `upload-session`, `admin-upload-course-material`, `admin-delete-course-material`, `admin-reorder-course-agenda`, returning the `ServerActionResult` shape `{ ok, refresh, toast, data }` that `actions.ts:82–95` consumes. The "missing endpoint" flagged in §4.6 exists here, in JS.
- **File uploads over IR already work**: `ContextUploadDropArea.widget.tsx` serializes `File`s into the action context (`{name, size, type, encoding: "utf8"|"base64", text|base64}`), and `lib/course-material-service.js` decodes base64, validates magic bytes (PNG/JPEG/GIF/WebP) or `<svg>` presence, sanitizes filenames (no separators/`..`, charset allowlist, extension allowlist per kind), and writes atomically.
- **Storage is files, not SQL**: `course/slides/*.md`, `course/handouts/*.md`, `course/media/*` (served at `/course-assets/`, server.js:615), `course/course-metadata.json` for landing-page metadata/outcomes/agenda. Loaders (`slide-loader.js`, `handout-loader.js`, `course-content-loader.js`) parse Markdown (frontmatter + `ArticleBlock` extraction) at request time.
- **Form-heavy editing uses native form posts**: `formPanel({ method: "post", formAction: "/settings/course-metadata" })` → server saves → `res.redirect("/pages/admin-course-cms?status=saved")` → the rebuilt page shows `status`/`statusMessage` (admin-course-cms.js:81–141, server.js:318–358). Interactivity across requests is **query-parameter driven**, not client state.
- **Admin gate**: display name `admin_<name>` (`courseMaterial.isAdminUser`), enforced per action/route (server.js:92–109) — demo-grade by design.

**What this replaces:** no RTK Query endpoints, no Go REST API in this repo, no `web/` containers for the CMS. The `web/` app remains the RAG-evaluation dashboard; the CMS ships to go-go-course through npm + the Go module.

### 6.8 What the CMS components need for IR consumption (revised)

The go-go-course study surfaces concrete requirements the pure-React API (§6.3–6.6) doesn't yet satisfy:

1. **Widget adapters for every new component** (`X.widget.tsx` + `X.widget.yaml`), with callback props mapped to `on*Action?: ActionSpec` (the `HandoutDocumentShell` convention). Registered in a new `cmsWidgetRegistry` under module `cms.dsl`.
2. **File uploads**: `MediaLibraryPanel`'s adapter must reuse the `ContextUploadDropArea.widget.tsx` serialization (extract `serializeUploadFile` into a shared `widgets/uploadSerialization.ts`); the action context carries `{files, fileNames, fileCount}` exactly as `course-material-service.js` expects. Upload *progress* (`UploadQueueItem`) is unavailable over a single JSON POST — over IR, `uploads` stays server-fed or omitted; `UploadQueueList` remains most useful in future streaming hosts.
3. **Form participation (uncontrolled mode)**: pages like admin-course-cms post native forms. `MarkdownEditor` and `TagListInput` need `name` + `defaultValue` support so they can sit inside `formPanel({method,formAction})`:
   - `MarkdownEditor`: `name?: string; defaultValue?: string` — the widget adapter wraps it in a local controlled state (toolbar + live `MarkdownArticle` preview both work purely client-side) while a synced hidden/named textarea joins the form post. Live preview over IR is therefore possible without any server round-trip.
   - `TagListInput`: `name?: string` emits `<input type="hidden" name value={tags.join(",")}>`; server splits on commas.
   - `SearchField`/`SelectInput` filters: adapters render inside a GET `formPanel` or map `onValueChange` to `action.navigate("?query=$value")`.
4. **Navigation-driven selection/paging**: adapters map `onAssetSelect`/`onPageChange`/`onArticleSelect` to interpolated `action.navigate` specs (`"?asset=${assetId}"`, `"?page=${page}"`), matching go-go-course's existing `?status=`/`?slide=`/`?doc=` pattern. The JS page reads `query` and passes `selectedAssetIds`/`page` back into props — state lives in the URL.
5. **Destructive confirmation over IR**: today `data.cell.actionButton("Delete", action.server("admin-delete-course-material"))` fires immediately (admin-common.js). Rather than replicating `ConfirmDialog`'s client state in IR, extend `ActionSpec` with `confirm?: string` handled centrally in `dispatchWidgetAction` (show a `DialogShell`-styled confirm — or `window.confirm` in v1 — before dispatching). One flat mechanism covers every action site (D-9).
6. **`CmsAsset` mapping in go-go-course**: `courseMaterial.listCourseMaterial()` media rows are `{file, size, modifiedAt, href}`; extend the service to emit `CmsAsset`-shaped objects (`id=file`, `kind` by extension, `src=/course-assets/<file>`, `mime` by extension, byte size) so `cms.mediaLibraryPanel({assets})` renders `AssetTile` thumbnails directly. `status`/`tags` default to `"published"`/`[]` until the file-backed store grows a sidecar metadata JSON (same pattern as `course-metadata.json`).

### 6.9 `cms.dsl` and the go-go-course integration (revised — now the delivery path)

Per reference/02 §5, plus what §6.7 adds:

1. **Go side (this repo)**: `cmsHelpers` map + `moduleSpec` in `pkg/widgetdsl/module.go` (module `cms.dsl`, `action: true`), recipes `mediaLibrary` and `articleList` normalizing `on*` strings/specs via `normalizeActionSpec`; `TypeScriptModule(CmsModuleName)`; provider entry in `pkg/xgoja/providers/widgetsite/provider.go`; doc page update; `module_test.go` boundary tests. Fix the known drift while there (`cell.linkButton`/`cell.actionButton`/`action.download` typings, `contextGroupedStripDiagram` helper).
2. **TS side (this repo)**: widget adapters per §6.8, `"cms.dsl"` added to the `WidgetModule` union (registry.ts:5), `cmsWidgetRegistry` merged into `defaultWidgetRegistry`, `WidgetRenderer.cms.stories.tsx` under `Widget IR/Renderer/CMS`.
3. **Release**: publish npm `@go-go-golems/rag-evaluation-site` (currently 0.1.16 in go-go-course) and tag the Go module (currently v0.1.2 there); bump both pins in go-go-course; rebuild `webapp` → `assets/` → xgoja binary.
4. **go-go-course side**: upgrade `lib/pages/admin-course-cms.js` — media table → `cms.recipes.mediaLibrary({assets})` with thumbnails; deletes gain `confirm`; handout authoring page using `formPanel` + named `markdownEditor` with client-side live preview; `contentStatusBadge` in material tables. Target authoring experience:

```js
const cms = require("cms.dsl")
// media library with thumbnails, selection via query param, upload via server action
cms.recipes.mediaLibrary({
  assets: courseMaterial.listCourseMaterial({...}).mediaAssets,   // CmsAsset[]
  selectedAssetIds: query.asset ? [query.asset] : [],
  onAssetSelect: ui.action.navigate("?asset=$assetId"),
  onFilesSelected: ui.action.server("admin-upload-course-material", { payload: { kind: "media" } }),
  onAssetDelete: { ...ui.action.server("admin-delete-course-material"), confirm: "Delete this media file?" },
})
```

## 7. Decision records

### Decision D-1: Markdown-first content model, stored as `ArticleBlock[]`

- **Context:** Articles could be stored as (a) one markdown string, (b) a block array, or (c) HTML. Two renderers already exist: string-based `MarkdownArticle` and block-based `RichArticle`; `HandoutDocumentShell` already prefers `blocks` over `body` (HandoutDocumentShell.tsx:75–83).
- **Options considered:** raw markdown string only; `ArticleBlock[]` with markdown blocks for prose; full custom block editor (Notion-style); HTML/contenteditable WYSIWYG.
- **Decision:** Store `ArticleBlock[]`. The v1 editor edits a single `markdown` block as text (with toolbar helpers); images/diagrams/galleries are separate typed blocks. Editing UI can grow toward per-block editing without a storage migration.
- **Rationale:** Aligns with the existing renderer precedence; keeps the live preview trivially exact (`MarkdownArticle` renders the same source being edited); avoids a WYSIWYG dependency that would fight the strict CSS rules; typed blocks are what Widget IR and the Goja DSL can address semantically (GUIDELINES.md:268 prefers semantic recipes).
- **Consequences:** The v1 editor is not WYSIWYG; embedded standalone images inside markdown remain possible (parser already supports them) but gallery/diagram blocks require block-level editing in a later phase. Must validate `blocks_json` server-side.
- **Status:** proposed

### Decision D-2: New `ContentStatusBadge` atom instead of extending `RagStatus`

- **Context:** `StatusText` owns a workflow status vocabulary (pending/running/failed/… — StatusText.tsx:4–14) used by evaluation/workflow UI. CMS needs draft/published/scheduled/archived.
- **Options considered:** extend `RagStatus` with the four content statuses; a generic `Badge` atom with free-form tone; a dedicated `ContentStatusBadge`.
- **Decision:** Dedicated `ContentStatusBadge` atom mapping content lifecycle onto existing tone tokens.
- **Rationale:** The two vocabularies serve different domains and will evolve independently; mixing them turns `RagStatus` into a junk drawer and weakens the strong status→tone contract that workflow views rely on. A free-form `Badge` violates rule 4 of GUIDELINES.md (unbounded styling APIs).
- **Consequences:** Slight duplication of the icon+tone pattern (acceptable; ~40 lines); Widget IR later gets a `ContentStatusBadge` node with a closed enum.
- **Status:** proposed

### Decision D-3: Native `<dialog>`-based `DialogShell` layout primitive

- **Context:** The package has no modal anything; a CMS needs asset picking and destructive confirmation. Guidelines forbid app-owned chrome inside the package but layout primitives answer "where do regions go" (GUIDELINES.md:85).
- **Options considered:** no modals (inline panels + view switching only); portal + div overlay implementation; native `<dialog>` element.
- **Decision:** Native `<dialog>` with `showModal()`, styled flat; `DialogShell` lives in layout.
- **Rationale:** Native gets focus trapping, Esc handling, and `::backdrop` for free with zero dependencies; baseline browser support is long since fine for this product; a flat black-bordered dialog with a Panel-style title bar sits naturally in the visual language. Inline-only was seriously considered (very much in the terminal spirit) but asset-picking from within a form genuinely needs an overlay to preserve form state.
- **Consequences:** Storybook stories must render the dialog opened via play/args (open prop) — add a `Static` story rendering the shell inline (no `showModal`) for visual diffing; jsdom tests need a `HTMLDialogElement.showModal` shim.
- **Status:** proposed

### Decision D-4: `MediaThumb` is an atom with owned load/error/empty states

- **Context:** Article figures render raw `<img>` today; broken URLs render browser default broken-image chrome (visible in `HandoutDocumentShell` "Illustrated" story screenshot).
- **Options considered:** keep raw `<img>` and style per-consumer; a molecule combining image+caption; an atom owning only the framed image states (caption stays with `FigureBlock`).
- **Decision:** Atom owning frame/aspect/fit/loading/broken/empty; `FigureBlock` continues to own captions/legends; `RichArticle`'s image block adopts `MediaThumb` internally in a later, opt-in pass.
- **Rationale:** Small controls and semantic markers are the atom charter (GUIDELINES.md:62–64); every CMS surface (tiles, detail, picker, cover field) needs identical image behavior; captioning is already solved one layer up.
- **Consequences:** Two image paths exist until `RichArticle`/`MarkdownArticle` adopt it internally; adoption must preserve print rules (`@media print` bounds in MarkdownArticle.module.css:166–175).
- **Status:** proposed

### Decision D-5: Sanitize `MarkdownArticle` URLs before any CMS exposure

- **Context:** `renderInline` emits `<a href>` unvalidated (MarkdownArticle.tsx:29); image `src` likewise (:100–104). Fixtures are trusted; CMS authors are not necessarily.
- **Options considered:** swap to a markdown library with sanitizer (marked+DOMPurify); keep the hand-rolled parser and add a URL allowlist; render untrusted content in a sandboxed iframe.
- **Decision:** Keep the parser; add `sanitizeUrl(url)` allowing only `http:`, `https:`, `mailto:`, and relative paths (reject `javascript:`, `data:` except `data:image/` if ever needed); apply to both `href` and `src`; add `rel="noopener noreferrer"` on external links.
- **Rationale:** The parser is deliberately minimal and dependency-free (a package-level virtue); the actual injection surface is URLs since the parser never emits raw HTML (all text goes through JSX text nodes, which React escapes). A library swap would change rendering subtly across existing course content.
- **Consequences:** ~20 lines + tests; document that raw-HTML markdown is unsupported by design.
- **Status:** proposed

### Decision D-6: CMS mutations flow through RTK Query containers, not widget `server` actions

- **Context:** The Widget IR action model has a `server` kind posting to `/api/widget/actions/{name}` (actions.ts:82–95), unimplemented in Go. CMS needs many mutations.
- **Decision (original):** Containers + RTK Query for the CMS application surface; server actions stay a Widget-IR-era concern.
- **Status:** **superseded by D-8** — the premise ("the server-action endpoint doesn't exist") was wrong for the actual consumer: go-go-course implements `/api/widget/actions/:name` in JS (server.js:383–456) and already runs file uploads, deletes, and reorders through it.

### Decision D-8: Widget-DSL-first delivery via go-go-course; no REST/RTK layer

- **Context:** The CMS consumer is `go-go-course`, an xgoja binary whose whole UI is Widget IR authored in JS; it embeds the WidgetRenderer SPA (npm 0.1.16) and pulls the Go DSL modules from this repo (v0.1.2). It already has server actions, file-backed storage under `course/`, upload serialization/validation, form-post + redirect editing, and an admin CMS page.
- **Options considered:** (a) original plan — Go REST API + RTK Query containers in this repo's `web/`; (b) Widget-DSL-first — ship the components as widget adapters + `cms.dsl` helpers and integrate in go-go-course's JS pages; (c) both.
- **Decision:** (b). The React components stay presentational; their widget adapters plus `cms.dsl` are the integration surface; go-go-course's existing services (`course-material-service.js`, `course-metadata-service.js`) are the backend.
- **Rationale:** The mutation path, storage, validation, auth-gate, and page-composition patterns all exist and are proven in go-go-course; building a parallel REST/SQL stack in this repo would duplicate them for a consumer that wouldn't use it. IR-authoring is the product's declared direction (widget manifests, DSL modules, `RagEvaluationSiteApp`).
- **Consequences:** Every CMS component needs a widget adapter with `on*Action` props plus IR-compatibility features (form `name`/`defaultValue`, navigate-based selection, shared upload serialization — §6.8). Client-only affordances (upload progress bars, controlled live-preview editing) must either move inside widget adapters (local state is fine there) or wait. Release loop is npm publish + Go tag + double pin-bump in go-go-course. The `web/` app is untouched.
- **Status:** accepted (per user direction, 2026-07-03)

### Decision D-9: `ActionSpec.confirm` for destructive IR actions

- **Context:** IR pages can't hold dialog state; go-go-course's delete buttons (`cell.actionButton("Delete", action.server(...))`, admin-common.js) fire with no confirmation today. `ConfirmDialog` (P4.1) only helps stateful React hosts.
- **Options considered:** interactive confirm state in IR (server round-trip to render a confirm page); a `ConfirmDialog` widget node with its own open-state protocol; a `confirm?: string` field on `ActionSpec` handled centrally in `dispatchWidgetAction`.
- **Decision:** `confirm?: string` on `ActionSpec`; when present, `dispatchWidgetAction` shows a confirm (v1: `window.confirm`; v2: `DialogShell`-styled confirm rendered by the SPA shell) before dispatching. Applies to every action kind.
- **Rationale:** One flat mechanism at the dispatch layer covers all current and future action sites (table cells, buttons, tiles) with a single JSON field — no state protocol, no per-widget wiring.
- **Consequences:** `ir.ts` ActionSpec fields grow by one optional member; `actions.ts` gains ~6 lines; `pkg/widgetdsl` action builders pass `confirm` through `options` already (mergeOptions). v2 (styled dialog) needs a mount point in `RagEvaluationSiteApp`.
- **Status:** proposed

### Decision D-7: `CmsShell` starts as its own organism; generalize with `CourseStudioShell` only if they converge

- **Context:** `CourseStudioShell` is 90% of the needed chrome but has course defaults baked in (`title="Context Window Engineering"`, `subtitle="Course studio"` — CourseStudioShell.tsx:7–18) and ships `courseStudioNavSections`.
- **Options considered:** rename/generalize `CourseStudioShell` to `StudioShell` now; parametrize it further; new thin `CmsShell` that composes the same layout pieces (`SidebarShell` + `SidebarNav`).
- **Decision:** New `CmsShell` composed from the same layout/molecule primitives (~60 lines), plus exported `cmsNavSections`.
- **Rationale:** Renaming a published component breaks npm consumers (package is public, v0.1.16); the shells may diverge (CMS wants breadcrumbs/user slot in the header). Composition is cheap because the heavy lifting already lives in `SidebarShell`/`SidebarNav`.
- **Consequences:** Two similar shells; if a third appears, extract `StudioShell` then (rule of three).
- **Status:** proposed

## 8. Key flows (pseudocode — revised for the widget-DSL path)

### 8.1 Upload flow (IR / go-go-course)

```
user drops files on MediaLibraryPanel's FileDropZone (rendered from cms.mediaLibraryPanel IR)
  → MediaLibraryPanel.widget adapter serializes each File (shared serializeUploadFile:
      text files → {encoding:"utf8", text}; binaries → {encoding:"base64", base64})
  → dispatchAction(onFilesSelectedAction) → POST /api/widget/actions/admin-upload-course-material
      body: { payload: {kind:"media"}, context: {files, fileNames, fileCount} }
  → server.js: requireAdminProfile → course-material-service.storeCourseMaterialUpload
      sanitize filename → decode base64 → magic-byte/SVG validation → atomic write to course/media/
  → responds { ok:true, refresh:true, toast:"Uploaded media <file>" }
  → actions.ts sees result.refresh → popstate → page rebuilds → new AssetTile appears
     (no per-file progress over a single JSON POST; UploadQueueList reserved for streaming hosts)
```

### 8.2 Article/handout editing flow (IR / go-go-course)

```
admin page IR: formPanel({ method:"post", formAction:"/settings/handout-body" },
    formRow({ control: markdownEditor({ name:"body", defaultValue: handout.source }) }), …)
  MarkdownEditor.widget wraps the editor in local client state:
    toolbar wraps selection, right pane renders MarkdownArticle(source=localValue) live —
    no server round-trip; a synced named textarea carries the value in the form post
  Save (native submit) → server validates via course-content-loader.parseCourseMarkdownFile
    → writeFileAtomic(course/handouts/<file>) → redirect "?status=saved"
    → rebuilt page renders formPanel status="success"
  insert-image: v1 = paste /course-assets/<file> path from the media library page;
    v2 = AssetPickerDialog inside the widget adapter (client state), inserting at cursor
```

### 8.3 Asset browse / destructive action flow (IR / go-go-course)

```
media library page: cms.recipes.mediaLibrary({ assets, selectedAssetIds:[query.asset], … })
  tile click → action.navigate("?asset=$assetId") → page rebuild with detail section
  Delete → action.server("admin-delete-course-material") with confirm:"Delete <file>?" (D-9)
    → dispatchWidgetAction shows confirm → POST action → deleteCourseMaterial (unlink)
    → { ok, refresh:true, toast:"Deleted media <file>" } → page rebuilds
usage list: server greps course/handouts + course/slides sources for the /course-assets/<file>
  reference and passes usedBy into the page (file-backed analogue of the SQL scan)
```

### 8.4 Component/data topology (revised)

```
go-go-course (xgoja binary "minitrace-viz")
┌───────────────────────────────────────────────────────────────────────────┐
│ server.js (express over goja)                                             │
│  GET /api/widget/pages/:id  ── lib/pages/*.js build IR via ui/data/…/cms  │
│  POST /api/widget/actions/:name ── course-material / metadata services    │
│  /course-assets/* ← course/media   ·   course/{slides,handouts}/*.md      │
│  SPA (embedded assets) = webapp shell around RagEvaluationSiteApp         │
└──────────────▲──────────────────────────────▲─────────────────────────────┘
               │ IR JSON pages                │ ActionSpec dispatch (+confirm, D-9)
┌──────────────┴──────────────────────────────┴─────────────────────────────┐
│ @go-go-golems/rag-evaluation-site (npm, this repo)                        │
│  WidgetRenderer + defaultRegistry (+ cmsWidgetRegistry)                   │
│  MediaLibraryPanel ─ TileGrid > AssetTile > MediaThumb/ContentStatusBadge │
│  ArticleListPanel ─ DataTable + Tag + boxed IconButtons                   │
│  formPanel + markdownEditor(name) ─ live preview inside the adapter      │
└────────────────────────────────────────────────────────────────────────────┘
  Go module (this repo): pkg/widgetdsl cms.dsl helpers → rag-widget-site provider
```

## 9. Implementation phases (file-level)

Each component = 4 files (`X.tsx`, `X.module.css`, `X.stories.tsx`, `index.ts`) in its layer directory + a barrel export in the layer `index.ts`. Run `pnpm --dir packages/rag-evaluation-site typecheck` and `pnpm biome check --write .` after every step; Storybook is the review surface.

**Phase 0 — groundwork (small, independent):**
1. `sanitizeUrl` in `MarkdownArticle.tsx` (+ unit-ish story fixtures showing stripped `javascript:` links) — D-5.
2. Add `CmsAsset`/`CmsArticleSummary`/`CmsArticleDetail`/`ArticleGalleryBlock` types (new `src/cms/types.ts`; re-export from `src/index.ts` like `src/context`).
3. Optional consistency sweep: `data-rag-component` → layer-correct attributes on DataTable/MetadataGrid/AppNav (check visual-diff tooling first — these attrs are extraction targets, GUIDELINES.md:16).

**Phase 1 — atoms:** `MediaThumb`, `Tag`, `ContentStatusBadge`, `MeterBar`. Stories: `Design System/Atoms/<Name>` with default/empty/broken (MediaThumb), selected/removable/disabled (Tag), all four statuses (ContentStatusBadge), 0/50/100/error (MeterBar).

**Phase 2 — layout:** `TileGrid`, `DialogShell`. Stories: `Design System/Layout/<Name>`; DialogShell gets `Static` (inline, for visual diff) + `Interactive` (open/close with useState) stories.

**Phase 3 — molecules:** `AssetTile`, `TagListInput`, `Breadcrumbs`, `Pagination`, `SearchField`, `EmptyState`, `UploadQueueList`, `MarkdownEditor`. Stories: `Component Library/Molecules/<Name>` covering the GUIDELINES state matrix (dense/overflow: 50-asset tile, 12-tag input, long filenames).

**Phase 4 — organisms:** `MediaLibraryPanel`, `AssetPickerDialog`, `ArticleListPanel`, `ArticleEditorPanel`, `AssetDetailPanel`, `CmsShell`, `ConfirmDialog`. Stories with rich fixtures (create `src/cms/fixtures.ts` mirroring `src/context/fixtures.ts`) including an `Interactive` story per organism wiring local useState, like `HandoutDocumentShell`'s.

**Phase 5 — IR enablement in this repo (revised; replaces the old Go-API/RTK phase):**
1. Cross-cutting: `ActionSpec.confirm` in `ir.ts` + `dispatchWidgetAction` (D-9); extract `serializeUploadFile` from `ContextUploadDropArea.widget.tsx` into `src/widgets/uploadSerialization.ts`.
2. Uncontrolled/form modes: `name`/`defaultValue` on `MarkdownEditor` and `TagListInput` (§6.8.3); verify inputs already forward `name` (TextInput/SelectInput/TextareaInput do — passthrough props).
3. Widget adapters + manifests (`.widget.tsx`/`.widget.yaml`) for: MediaThumb, Tag, ContentStatusBadge, MeterBar, TileGrid, AssetTile, Breadcrumbs, Pagination, SearchField, EmptyState, MarkdownEditor, MediaLibraryPanel, ArticleListPanel, CmsShell (skip DialogShell/ConfirmDialog/AssetPickerDialog/ArticleEditorPanel/AssetDetailPanel in IR v1 — dialog/controlled-state surfaces; the adapter-internal live preview in MarkdownEditor covers the editing need).
4. `RagWidgetType` additions (ir.ts:54–115), `cmsWidgetRegistry` + `"cms.dsl"` module (registry.ts:5, defaultRegistry.ts), `WidgetRenderer.cms.stories.tsx` under `Widget IR/Renderer/CMS`.
5. Go: `cms.dsl` in `pkg/widgetdsl` (helpers + `mediaLibrary`/`articleList` recipes), provider entry in `widgetsite/provider.go`, typings, tests; fix the reference/02 §4 drift (typescript.go omissions, `contextGroupedStripDiagram` helper).

**Phase 6 — go-go-course integration (in the go-go-course repo):**
1. Release: publish npm `@go-go-golems/rag-evaluation-site` and tag `github.com/go-go-golems/rag-evaluation-system`; bump `webapp/package.json` (from 0.1.16) and `xgoja.yaml`/`go.mod` (from v0.1.2); rebuild webapp → assets → binary.
2. `course-material-service.js`: emit `CmsAsset`-shaped media entries (§6.8.6); optional sidecar metadata JSON for tags/status/alt.
3. `lib/pages/admin-course-cms.js`: media DataTable → `cms.recipes.mediaLibrary` (thumbnails, upload, navigate-selection); `confirm` on delete actions; `contentStatusBadge` in material tables.
4. New handout-editing page: `formPanel` + named `markdownEditor` with adapter-internal live preview; `/settings/handout-body` save route following the existing form-post + redirect pattern.
5. Smoke: upload an image via the drop area, see the thumbnail, reference it from a handout, edit the handout with live preview, delete with confirm — end-to-end in the running binary.

Suggested tickets: one per phase; phases 1–3 parallelize well across people once phase 0 lands. Phase 5 items 1–2 are prerequisites for item 3.

## 10. Testing and validation strategy

1. **Typecheck:** `pnpm --dir packages/rag-evaluation-site typecheck` (tsc, no emit) — gates every phase.
2. **Format/lint:** `pnpm biome check .` (tabs, double quotes, import order).
3. **Storybook as spec:** every new component's stories must cover the GUIDELINES.md:149–157 state matrix; `pnpm --dir packages/rag-evaluation-site build-storybook` must pass (CI builds a static Storybook — see `Dockerfile.storybook-static`).
4. **Visual style audit (manual but scripted):** grep-gates that fail review if violated — `rg -n "border-radius|box-shadow" packages/rag-evaluation-site/src/components` must return nothing new; `rg -n "font-size" …/components/{atoms,molecules,organisms}` only in tokenized fallbacks.
5. **Playwright story smoke:** iterate `http://localhost:6007/index.json` entries for the new titles and screenshot `iframe.html?id=…` (the exact loop used for this ticket's evidence — see `sources/screenshots/`); assert no console errors and that `data-rag-*` roots exist via `document.querySelector('[data-rag-organism="MediaLibraryPanel"]')`.
6. **Sanitizer tests:** fixture stories + (if a test runner is introduced) unit tests for `sanitizeUrl` covering `javascript:`, `data:text/html`, protocol-relative, and relative URLs.
7. **Go:** `go test ./...`; handler tests for multipart upload (happy path, oversize, disallowed mime), blob serving headers, and article CRUD round-trips including `blocks_json` validation.
8. **Integration smoke:** `rag-eval serve` + `pnpm --dir web dev`, upload an image, insert it into an article, publish, and view — the flow in §8.1–8.2 end-to-end.

## 11. Risks, alternatives, open questions

**Risks:**

- *Scope creep in the editor.* Markdown-first is a deliberate ceiling; pressure for WYSIWYG will come. Mitigation: D-1 documents the upgrade path (per-block editing) that doesn't invalidate storage.
- *SVG uploads are an XSS vector* (scriptable). go-go-course's validator only checks for an `<svg>` element (course-material-service.js `validateSvgUpload`) — scripts inside pass. Mitigate with goja-text's `sanitize` module (already in its runtime) before write, or serve `/course-assets/*.svg` with a restrictive CSP.
- *Published-package API stability.* Everything added to `src/index.ts` ships to npm consumers at the next version bump — and go-go-course pins exact versions (npm 0.1.16, Go v0.1.2), so nothing reaches it until both pins move. That protects the consumer but makes the **release loop the main integration cost**: npm publish + Go tag + two pin bumps + webapp rebuild per iteration. Use `webapp`'s `dev:rag-local`/`build:rag-local` (`RAG_SITE_SRC=…` vite alias) to iterate against the package source before publishing.
- *IR-mode capability gaps.* Upload progress, drag-reorder, and modal pickers don't translate to JSON round-trips; the plan confines them to adapter-internal client state (live preview) or defers them. Watch for pressure to smuggle app state into widgets — that's the GUIDELINES rule-5 boundary.
- *Broken-asset references* after delete: v1 shows `MediaThumb` broken state rather than blocking deletes; the usage list (8.3) gives authors visibility.

**Alternatives considered (system-level):** adopting an off-the-shelf headless CMS (Payload/Strapi) and only building read-side widgets — rejected because the product's value is the integrated, strictly-styled studio experience and the Widget IR path; a generic `Badge`/`Chip` styling API — rejected as a GUIDELINES rule-4 violation.

**Open questions:**

1. Folders/collections for media — flat + tags (current design) vs. hierarchical? `Breadcrumbs` is included so hierarchy can land later without new primitives.
2. Should `RichArticle` image blocks adopt `MediaThumb` immediately (visual change to existing course handouts) or behind a prop? (Proposed: behind `frame` defaulting to current borderless look.)
3. Does `documents` (RAG corpus) ever unify with `articles` (CMS)? They share shape but different lifecycles; keeping them separate avoids coupling eval pipelines to authoring. Revisit if articles should be searchable via the existing BM25/vector endpoints — that would be a compelling unification.
4. Dark theme: tokens are light-only today; the CMS adds no new blockers but also shouldn't hardcode anything that would (all proposals are token-pure).

## 12. References

**Design system package (all under `packages/rag-evaluation-site/`):**
- `GUIDELINES.md` — layer rules, typography, Storybook conventions (the constitution).
- `src/theme.css` — all tokens (colors :3–40, font roles :45–53).
- `src/components/foundation/…` — Text, Caption, CodeText, StatusText (status vocab :4–14), Divider, VisuallyHidden.
- `src/components/atoms/…` — Button (Button.tsx:7–12), inputs, IconButton, UploadGlyph, ContextStudioNavIcon (icon ids :4–12).
- `src/components/layout/…` — Panel, SidebarShell (:4–11), SplitPane, FormRow (:6–16), TabList, ScrollRegion, SectionBlock, SlideShell.
- `src/components/molecules/MarkdownArticle/MarkdownArticle.tsx` — parser (:10–39 inline, :52–56 images), props (:4–6).
- `src/components/organisms/RichArticle/RichArticle.tsx` (:8–11) and `src/context/types.ts` — ArticleBlock union (:170–192), ContextHandoutDocument (:194–205).
- `src/components/organisms/HandoutDocumentShell/HandoutDocumentShell.tsx` — props (:8–20), content precedence (:75–83).
- `src/components/molecules/{DocumentListPanel,DocumentPreviewToolbar,FileDropZone,DataTable}/…` — read-side + upload building blocks.
- `src/components/organisms/{FormPanel,CourseStudioShell}/…` — form chrome (:7–20), shell template (:7–59).
- `src/widgets/ir.ts` (nodes :40–122, actions :126–161), `src/widgets/actions.ts` (dispatch :24–96), `src/widgets/registry.ts`, `defaultRegistry.ts`, `WidgetRenderer.tsx`.
- `.storybook/main.ts`, `.storybook/preview.ts` — Storybook config (port 6007).

**Web app:** `web/src/App.tsx` (views :11–19), `web/src/services/api.ts` (base :303–305, endpoints :316–494), `web/src/store/index.ts`, `web/src/storybook/MockApiProvider.tsx`.

**Go backend:** `internal/api/handlers.go` (routes :38–89), `internal/api/dsl_handlers.go` (:5–88), `internal/db/db.go` (schema :70–230; documents :81–98), `cmd/rag-eval/cmds/serve/server.go` (:19–63), `pkg/widgetdsl/module.go` (:14–19), `pkg/defaultspa/spa.go`, `internal/web/spa.go`.

**Ticket evidence:** `sources/screenshots/*.png` — 14 full-page story captures (foundation palette/typography, Button, Panel, DataTable, MarkdownArticle, FigureBlock, DocumentListPanel, DocumentPreviewToolbar, FileDropZone, RichArticle, HandoutDocumentShell, CourseStudioShell, FormPanel); `reference/01-investigation-diary.md` — how this evidence was gathered.

**The consumer — go-go-course (sibling repo, `go-go-course/cmd/go-go-course/`):**
- `xgoja.yaml` — buildspec: `rag-widget-site` provider (Go pin v0.1.2), DSL module selection, embedded SPA assets.
- `server.js` — `GET /api/widget/pages/:id` (:128–153), **`POST /api/widget/actions/:name`** (:383–456: upload-session, admin-upload/delete-course-material, admin-reorder-course-agenda), form-post routes `/settings/*` (:305–358), `/course-assets` static (:615), SPA mount (:619–625).
- `lib/course-material-service.js` — file-backed CMS core: kind targets (slides/handouts/media), filename sanitization, base64 decode, magic-byte/SVG validation, atomic writes, `listCourseMaterial`.
- `lib/course-metadata-service.js` — `course/course-metadata.json` (metadata/outcomes/agenda) save/reorder.
- `lib/pages/admin-course-cms.js` — the existing IR-authored admin CMS page (formPanels + upload panels + material tables).
- `lib/pages/admin-common.js` — `adminUploadPanel` (contextUploadDropArea + server action) and material DataTables (delete without confirm — D-9 motivation).
- `webapp/src/main.tsx` + `webapp/package.json` — 15-line SPA shell around `RagEvaluationSiteApp apiBase="/api/widget"`; npm pin 0.1.16; `dev:rag-local` source-alias script.
- `packages/rag-evaluation-site/src/components/organisms/ContextUploadDropArea/ContextUploadDropArea.widget.tsx` (this repo) — the file-serialization pattern (`SerializedUploadFile`, utf8/base64) that MediaLibraryPanel's adapter must reuse.

**Prior tickets:** `ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--…/design-doc/01-design-system-unification-analysis-and-implementation-guide.md`; `ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--…/design-doc/03-context-viewer-integration-plan-after-design-system-unification.md`.
