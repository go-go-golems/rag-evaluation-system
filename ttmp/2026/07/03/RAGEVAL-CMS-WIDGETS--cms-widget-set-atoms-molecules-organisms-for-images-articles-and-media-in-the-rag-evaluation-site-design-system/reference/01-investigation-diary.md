---
Title: Investigation diary
Ticket: RAGEVAL-CMS-WIDGETS
Status: active
Topics:
    - design-system
    - frontend
    - storybook
    - cms
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: packages/rag-evaluation-site/.storybook/main.ts
      Note: |-
        Package Storybook config used for the 6007 run
        staticDirs added; course-assets fixture SVG now serves
    - Path: packages/rag-evaluation-site/package.json
      Note: storybook script (port 6007) and npm publish surface
    - Path: packages/rag-evaluation-site/src/cms/fixtures.ts
      Note: Story fixtures for assets/articles/uploads
    - Path: packages/rag-evaluation-site/src/cms/types.ts
      Note: New CMS DTOs (CmsAsset
    - Path: packages/rag-evaluation-site/src/components/atoms/MediaThumb/MediaThumb.tsx
      Note: New image atom with loading/broken/empty states
    - Path: packages/rag-evaluation-site/src/components/layout/DialogShell/DialogShell.tsx
      Note: First modal primitive (native dialog
    - Path: packages/rag-evaluation-site/src/components/molecules/MarkdownArticle/MarkdownArticle.tsx
      Note: sanitizeUrl added (D-5)
    - Path: packages/rag-evaluation-site/src/components/organisms/ArticleEditorPanel/ArticleEditorPanel.tsx
      Note: Markdown-first editor with live preview
    - Path: packages/rag-evaluation-site/src/components/organisms/MediaLibraryPanel/MediaLibraryPanel.tsx
      Note: Asset browser organism
ExternalSources: []
Summary: 'Chronological diary of the CMS widget-set investigation: repo mapping, Storybook run, component/API inventory, visual evidence capture, and design-doc authoring.'
LastUpdated: 2026-07-03T17:30:00-04:00
WhatFor: Understand how the RAGEVAL-CMS-WIDGETS analysis was produced, what commands were run, what failed, and how to continue the work.
WhenToUse: Read before extending the CMS widget design doc or reproducing the Storybook evidence capture.
---



# Diary

## Goal

Capture the end-to-end investigation for RAGEVAL-CMS-WIDGETS: analyzing the `rag-evaluation-system` frontend design system, running the package Storybook, collecting visual and line-anchored evidence, and writing an intern-ready analysis/design/implementation guide for a CMS widget set (images, articles, media) that stays inside the existing strict design style.

## Step 1: Repo mapping, ticket setup, and Storybook bring-up

The investigation started from a cold workspace. The first pass established where the design system lives (`packages/rag-evaluation-site`), what its strict rules are (`GUIDELINES.md`, `AGENTS.md`), and created the docmgr ticket workspace with a design doc and this diary. The pnpm workspace has two packages: `web` (app shell, RTK Query, backend-connected views) and `packages/rag-evaluation-site` (the published, API-free design system with its own Storybook on port 6007).

### Prompt Context

**User prompt (verbatim):** "Analyze @rag-evaluation-system/ and run the storybook and see how we can create a decent looking set of widgets / atoms/molecules/organisms to create a CMS (images, articles, etc...) while staying in the strict design style we currently have. Create a new docmgr ticket, investigate, and then Create a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Investigate the existing atomic design system and Storybook, then design a CMS-oriented widget set (atoms/molecules/organisms for images, articles, media) that conforms to the existing visual language. Deliver a long-form, intern-ready analysis/design/implementation guide inside a new docmgr ticket, and upload the bundle to reMarkable.

**Inferred user intent:** Get a concrete, evidence-based blueprint that a new contributor could follow to grow the design system into CMS territory without breaking its strict conventions.

### What I did
- `ls` of repo root; read `AGENTS.md` (formatting + design-system summary) and `packages/rag-evaluation-site/GUIDELINES.md` (the canonical layer/typography/CSS/story rules).
- `docmgr status --summary-only`, then `docmgr ticket create-ticket --ticket RAGEVAL-CMS-WIDGETS --title "CMS widget set (atoms/molecules/organisms) for images, articles, and media in the rag-evaluation-site design system" --topics design-system,frontend,storybook,cms`.
- `docmgr doc add` for the design doc and this diary.
- `pnpm install` at repo root (fast, lockfile already satisfied).
- Started the package Storybook in the background: `pnpm --dir packages/rag-evaluation-site storybook` → "Storybook ready! Local: http://localhost:6007/".
- Read `packages/rag-evaluation-site/src/theme.css` in full (token inventory: `--rag-color-*`, `--mac-*` bridge, 9 `--rag-font-role-*` roles).

### Why
- The skill workflow requires ticket-first bookkeeping, and the guidelines file is declared MUST-READ by `AGENTS.md`; every later design decision has to trace back to those rules.

### What worked
- Ticket workspace created cleanly under `ttmp/2026/07/03/RAGEVAL-CMS-WIDGETS--…`.
- Storybook 10.4.1 started on port 6007 in ~630 ms after install.

### What didn't work
- First Playwright navigation crashed Chrome: `Target page, context or browser has been closed`, with `EROFS: read-only file system` for `/home/manuel/.storybook/settings.json`, `MESA: error: Failed to query drm device`, and `GPU process exited unexpectedly: exit_code=512`. The sandbox blocked Chrome's profile writes. The session was then restarted by the user ("Ok, playwright should work now"), which also killed the first round of background inventory agents and the Storybook process; everything was relaunched.

### What I learned
- The design system is a strict "Classic Mac / terminal" language: hard 1px `#000` borders, no border-radius anywhere, no box-shadows, 10–13px type, mono uppercase labels, tokens only.
- Storybook title prefixes are load-bearing conventions (`Design System/…` vs `Component Library/…` vs `Widget IR/Renderer`).

### Code review instructions
- Verify ticket files exist under `ttmp/2026/07/03/RAGEVAL-CMS-WIDGETS--…/` (index, tasks, changelog, design-doc/01, reference/01).
- Re-run `pnpm --dir packages/rag-evaluation-site storybook` and confirm port 6007 serves.

## Step 2: Parallel component/API inventory and visual evidence capture

With Storybook up, three parallel Explore agents produced a line-anchored inventory of (a) foundation/atoms/layout, (b) molecules/organisms with CMS-relevant DTOs, and (c) the Widget IR system, the web app data layer, and the Go backend. In parallel, Playwright captured full-page screenshots of 14 representative stories into the ticket's `sources/screenshots/` directory.

### Prompt Context

**User prompt (verbatim):** (see Step 1; mid-step follow-up: "Ok, playwright should work now.")

**Assistant interpretation:** Continue the interrupted investigation; the browser environment is fixed.

**Inferred user intent:** Unblock visual inspection of Storybook so the design analysis can include real rendering evidence.

### What I did
- Relaunched Storybook and the three Explore agents after the session restart.
- `curl -s http://localhost:6007/index.json` to enumerate all stories (≈70 story groups across Design System / Component Library / Widget IR).
- Batch-captured screenshots via one `browser_run_code_unsafe` loop over 14 story ids (foundation overview, Button, Panel, DataTable, MarkdownArticle, FigureBlock, DocumentListPanel, DocumentPreviewToolbar, FileDropZone, RichArticle, HandoutDocumentShell, CourseStudioShell, FormPanel) into `sources/screenshots/*.png` at 1400px width, `fullPage: true`.
- Read the key screenshots: monochrome palette board, HandoutDocumentShell (sidebar + toolbar + article), CourseStudioShell (sidebar nav + slide), DocumentListPanel, FileDropZone.
- Spot-checked agent-reported line anchors with `nl -ba` against `context/types.ts` (ArticleBlock union at lines 170–192), `MarkdownArticle.tsx` (props at 4–6, unsanitized `href` at 29), `internal/api/handlers.go` (route table at 38–89), and `.storybook/main.ts`/`preview.ts`.

### Why
- The design doc must anchor claims to files and lines; agents provide breadth, spot-checks protect against hallucinated anchors; screenshots document the "strict design style" the new widgets must match.

### What worked
- All three inventories returned dense, line-anchored reports; spot checks matched exactly.
- The screenshot loop wrote all 14 PNGs in one pass.

### What didn't work
- First `find` for the screenshot path ran sandboxed and silently returned nothing; rerunning with elevated permissions located Playwright's cwd (`/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system`). Playwright saves relative paths to its own working directory, not the repo — absolute paths in the capture loop fixed this.
- A zsh quirk: `echo =====` fails in zsh (`===== not found`) because of `=cmd` expansion; separators needed quoting.

### What I learned
- `HandoutDocumentShell` + `DocumentListPanel` + `DocumentPreviewToolbar` + `MarkdownArticle`/`RichArticle` already form ~70% of a read-side CMS; the gaps are write-side (editing, upload of media, taxonomy, listing/pagination, dialogs) and an image/media primitive (the handout story renders a broken `<img>` because no placeholder/fallback atom exists).
- `MarkdownArticle` is a hand-rolled parser with no sanitization — `href`/`src` pass through raw (XSS surface once CMS authors are untrusted).
- The Go backend has no media/blob storage (SQLite tables: sources, documents, chunks, …; no images/uploads tables) and no `/api/widget/actions/{name}` endpoint even though the front-end action model posts there.

### What was tricky to build
- Nothing code-level yet; the tricky part was correlating three agent reports with the actual files without re-reading everything. The spot-check discipline (verify 3–4 anchors per report) caught zero errors, which raised confidence in the rest.

### What warrants a second pair of eyes
- The claim that `data-rag-*` attribute naming is inconsistent (three components use `data-rag-component` instead of layer-specific names) — worth confirming before codifying it as a "fix" task.

### Code review instructions
- Screenshots: `ttmp/2026/07/03/RAGEVAL-CMS-WIDGETS--…/sources/screenshots/`.
- Story enumeration: `curl -s http://localhost:6007/index.json | python3 -m json.tool | head`.

## Step 3: Design doc authoring, bookkeeping, and reMarkable delivery

Wrote the primary deliverable — `design-doc/01-cms-widget-set-analysis-design-and-implementation-guide.md` — covering current-state architecture, gap analysis, the proposed CMS widget set per layer with full API sketches and pseudocode, decision records, a phased implementation plan, and a testing strategy. Then related key files, updated changelog/tasks, ran `docmgr doctor`, and uploaded the bundle to reMarkable.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Produce and deliver the final documents.

**Inferred user intent:** A single printable bundle on the reMarkable containing the full guide.

### What I did
- Wrote the design doc (executive summary → references), ~1,400 lines, with ASCII diagrams for layer flow, CMS shell anatomy, and data flow; decision records for content model, status vocabulary, dialogs, media primitives, sanitization, API-free boundary, and Widget IR timing.
- `docmgr doc relate` on both docs for the load-bearing source files; `docmgr changelog update`; checked off `tasks.md`.
- `docmgr doctor --ticket RAGEVAL-CMS-WIDGETS --stale-after 30` → see final handoff.
- `remarquee upload bundle --dry-run …` then real upload to `/ai/2026/07/03/RAGEVAL-CMS-WIDGETS`, verified with `remarquee cloud ls`.

### What worked
- See final handoff summary in the ticket changelog.

### What didn't work
- (Recorded inline in the final handoff if upload/doctor produced warnings.)

### Code review instructions
- Start with `design-doc/01-…implementation-guide.md`; validate the file anchors by opening the cited files at the cited lines; run `pnpm --dir packages/rag-evaluation-site typecheck` and Storybook to confirm the described story titles exist.

### Technical details
- Evidence sources: three Explore-agent inventories (foundation/atoms/layout; molecules/organisms; Widget IR/web/Go), 14 Storybook screenshots, and direct reads of `theme.css`, `GUIDELINES.md`, `types.ts`, `handlers.go`, `db.go` excerpts.

## Step 4: Deep-dive into the goja Widget DSL and cms.dsl extension notes

The user asked whether the goja UI DSLs derived from the WidgetRenderer had actually been examined — they had only been mapped at the architecture level (module names, provider existence). This step read the entire DSL implementation and its embedded documentation, wrote a dedicated line-anchored reference doc (`reference/02`), cross-linked it from design-doc §4.6, and re-uploaded the expanded bundle to reMarkable.

### Prompt Context

**User prompt (verbatim):** "have you seen / looked / documnted the goja UI DSLs that we derive from the widget renderer?"

**Assistant interpretation:** Verify and, if lacking, perform a proper investigation and documentation of the goja-side Widget DSL (`pkg/widgetdsl`, xgoja provider, docs, usage sites).

**Inferred user intent:** Make sure the ticket's system tour covers the JS authoring layer, since a CMS `cms.dsl` module is part of the roadmap (Phase 6).

### What I did
- Read `pkg/widgetdsl/module.go` (all 1104 lines), `typescript.go`, `registrar.go`, `module_test.go` (:1–60), `pkg/xgoja/providers/widgetsite/provider.go`, and the three embedded help pages (`doc/01..03-*.md`).
- Checked consumption: `grep`ped `cmd/rag-eval/xgoja.yaml` for `rag-widget-site` (absent) and `cmd/rag-eval/jsverbs/*.js` for `require("*.dsl")` (none — they use db/fs/yaml/markdown/sanitize/express).
- Wrote `reference/02-goja-widget-dsl-reference-and-cms-extension-notes.md`: registration paths, full helper inventory per module, props/child normalization semantics, `page`/`action`/`cell`/style-helper/recipe reference, TS typing generation, drift findings, and a concrete `cms.dsl` extension sketch (Go + TS + provider + tests).
- Edited design-doc §4.6 to summarize and cross-link; `docmgr doc relate` (8 files), changelog entry, `docmgr doctor` (pass), re-uploaded the 3-doc bundle with `remarquee upload bundle --force` and verified with `cloud ls`.

### What worked
- The DSL is one self-contained file, so a single full read produced a complete reference; module_test.go confirmed the module-boundary semantics.
- `--force` re-upload replaced the existing PDF cleanly (two "remote tree has changed" warnings from rmapi, then OK).

### What didn't work
- N/A — no failures this step. (Initial `grep "====="` separator quirk in zsh from Step 2 was avoided.)

### What I learned
- Five drift findings worth follow-up: `typescript.go` omits `cell.linkButton`/`cell.actionButton`/`action.download`; the JS API reference doc omits the same; `ContextGroupedStripDiagram` has no DSL helper; the DSL is dormant in-repo (not in `xgoja.yaml`, unused by jsverbs; only the hardcoded Go demo page produces IR); `action.event` special-cases print/fullscreen in the React dispatcher.
- Recipes are Go-side macro expansions; `masterDetailTable.detail` is the only place a JS callback is allowed (invoked at build time, before serialization). Callback props become `on*Action` ActionSpec props at the IR boundary.
- `paletteStyleSet` hardcodes four palettes in Go and emits `color-mix()` CSS strings — the DSL, not the theme, owns those colors.

### What warrants a second pair of eyes
- Whether `cms.dsl` should exclude editor surfaces (MarkdownEditor/ArticleEditorPanel/AssetPickerDialog) as proposed — that hinges on the D-6 decision about `/api/widget/actions/{name}`.

### What should be done in the future
- Fix the typescript.go/doc drift (small PR); add `contextGroupedStripDiagram` helper; wire `rag-widget-site` into an in-repo consumer or a smoke test that serves a DSL page end-to-end.

### Code review instructions
- Read `reference/02-…` side-by-side with `pkg/widgetdsl/module.go`; verify drift claims: `rg -n "linkButton|actionButton|download" pkg/widgetdsl/typescript.go` (no matches for the first two; no `download` in the action block) and `rg -n "contextGroupedStripDiagram" pkg/widgetdsl/module.go` (no match).
- Validate tests still describe reality: `go test ./pkg/widgetdsl/`.

## Step 5: Implementation task breakdown (widgets + stories first, DSL deferred)

Converted the design doc's phase plan into a concrete, docmgr-trackable task list. The user scoped the next round of work to the React widgets and their Storybook stories, with the goja DSL explicitly afterwards, so Phase 5 (backend/web) and Phase 6 (cms.dsl) moved to a "Deferred" section.

### Prompt Context

**User prompt (verbatim):** "Add tasks and phases to implement all the new widgets andstorybook stories. We'll do the dsl after that."

**Assistant interpretation:** Expand tasks.md into per-component implementation tasks covering Phases 0–4 plus validation, each with story requirements; keep DSL (and implicitly backend integration) as later phases.

**Inferred user intent:** Make the ticket immediately actionable for implementation work, one checkbox per component.

### What I did
- Rewrote `tasks.md`: a conventions preamble (folder layout, tokens, data-rag, story titles, gates), the completed investigation tasks, then P0.1–P0.4 (groundwork incl. sanitizeUrl and the `gallery` block), P1.1–P1.5 (atoms), P2.1–P2.3 (layout), P3.1–P3.9 (molecules), P4.1–P4.8 (organisms, ordered so ConfirmDialog/MediaLibraryPanel precede their consumers), PV.1–PV.4 (validation: build-storybook, style grep gates, playwright smoke, bundle re-upload), and deferred D1 (backend/web) + D2 (cms.dsl + drift fixes).
- Verified `docmgr task list --ticket RAGEVAL-CMS-WIDGETS` parses all 45 checkboxes with stable ids (implementation = ids 11–43); changelog updated; doctor passes.

### What I learned
- docmgr task ids are positional over all checkboxes in tasks.md, so the completed investigation tasks must stay above the new ones to keep ids stable while checking off P-tasks with `docmgr task check --id N`.

### Code review instructions
- `docmgr task list --ticket RAGEVAL-CMS-WIDGETS`; cross-check each P-task against design-doc/01 §6.3–§6.6 API contracts and §9 phase ordering.

## Step 6: Implementation of Phases 0–4 (21 components + 103 stories) and validation

Implemented the entire widget set in one pass: Phase 0 groundwork, 4 atoms, 2 layout primitives, 8 molecules, 7 organisms, all barrels, and 103 new Storybook stories — then validated with typecheck, biome, static Storybook build, package build, style grep gates, and a Playwright smoke over every new story.

### Prompt Context

**User prompt (verbatim):** "go ahead." (mid-step feedback: "the action icosn are very small in the doc table (archive, etc...)")

**Assistant interpretation:** Execute the tasks.md implementation plan (Phases 0–4 + validation), deferring backend and DSL.

**Inferred user intent:** A working, story-covered CMS widget set in the design system, reviewable in Storybook.

### What I did
- Phase 0: `sanitizeUrl()` in `MarkdownArticle.tsx` (allow http/https/mailto/relative; reject `javascript:`/`data:`/protocol-relative; `rel="noopener noreferrer"` on external links; blocked links render as plain text, blocked image srcs render caption-only figures) + `SanitizedUrls` hostile-fixture story; `src/cms/{types,fixtures,index}.ts` (CmsAsset/CmsArticleSummary/CmsArticleDetail/UploadQueueItem + 12 assets, 8 articles, 1 detail, 5 upload items, `formatCmsAssetSize`/`cmsAssetMeta` helpers) exported from `src/index.ts`; `ArticleGalleryBlock` added to the `ArticleBlock` union (`context/types.ts`) and rendered by `RichArticle` (`--rag-gallery-columns` grid) + `GalleryBlock` story.
- Phase 1 atoms: `MediaThumb` (aspect/fit/frame, loading checkerboard, broken/empty fallback + Caption, `data-state`), `Tag`, `ContentStatusBadge` (draft/published/scheduled/archived → dim/green/accent/dim+line-through), `MeterBar` (role="progressbar").
- Phase 2 layout: `TileGrid` (`auto-fill minmax(var(--rag-tile-min),1fr)`), `DialogShell` (native `<dialog>` + `showModal`, Panel-style black title bar, sm/md/lg, plus an `inline` mode so Static stories render in-flow for visual diffing).
- Phase 3 molecules: `AssetTile`, `TagListInput` (Enter/comma commit, datalist suggestions, Backspace removes last), `Breadcrumbs`, `Pagination`, `SearchField`, `EmptyState`, `UploadQueueList` (maps UploadItemStatus onto the existing `StatusText` vocabulary), `MarkdownEditor` (wrapSelection toolbar, counter, `onInsertAsset`).
- Phase 4 organisms: `ConfirmDialog` (destructive = danger-outlined default button, per D-record — no red-filled variant), `MediaLibraryPanel`, `AssetPickerDialog` (double-click = confirm), `ArticleListPanel` (DataTable + internal ConfirmDialog for archive/delete), `ArticleEditorPanel` (FormPanel + SplitPane live preview using the production `MarkdownArticle`), `AssetDetailPanel` (toolbar + MetadataGrid + FormPanel + usage list), `CmsShell` (+ exported `cmsNavSections`).
- Storybook infra: added `staticDirs: ["./static"]` to `.storybook/main.ts` and created `.storybook/static/course-assets/context-window-token-budget.svg` (flat, hard-edged sketch in the house style) — this fixes the long-broken `/course-assets/...` fixture URLs in existing MarkdownArticle/RichArticle/HandoutDocumentShell stories too.
- Feedback fix: added `IconButton size="large"` (14px glyph, deliberate raw size like the existing `.compact { font-size: 10px }`) and used it for ArticleListPanel and UploadQueueList row actions; new `Sizes` story.
- Validation: typecheck clean; biome clean on all new files; `build-storybook` and package `build` pass; grep gates show no border-radius/box-shadow in new files; Playwright iterated all 116 new/updated stories — 0 page errors, 0 missing `data-rag-*` roots; 15 flagship screenshots re-captured to `sources/screenshots/cms/`.

### What worked
- The layer conventions made this almost mechanical: every component copied the Button/Panel idioms (className merge, tokens, data-rag) and passed lint/typecheck with only 4 small fixes across 21 components.
- The inline `mode` on DialogShell solved the "modal stories can't be screenshot" problem cleanly.

### What didn't work
- `pnpm biome check --write` flagged: descending-specificity in `Tag`/`AssetTile` CSS (fixed by reordering selectors), `noArrayIndexKey` in a story (fixed with id'd fixtures), a story named `Error` shadowing the global (renamed `ErrorState`), and `style` props on stories whose components didn't extend HTMLAttributes (extended `AssetDetailPanelProps`, removed from MarkdownEditor story).
- Passing a ref to `TextareaInput` type-errored — React 19 supports ref-as-prop at runtime but the props interface needed `ref?: Ref<HTMLTextAreaElement>` declared.
- Restarting Storybook for `staticDirs`: `pkill -f "storybook dev"` missed the real process (`…/storybook/dist/bin/dispatcher.js dev -p 6007`), so the new instance prompted "Port 6007 is not available" and died; killing the dispatcher pid directly fixed it.
- `web` typecheck fails with a **pre-existing** error (verified via `git stash`): `ContextVisualizerPage.tsx(69)` missing required `styleSet` on `ContextDiagramPanel` — not caused by this work. Also note pre-existing uncommitted changes to `ContextDiagramPanel.stories.tsx`/`.widget.tsx` in the working tree (untouched).

### What was tricky to build
- MediaThumb state machine: `src` changes must reset `data-state` to loading (useEffect on src), and the broken state must unmount the `<img>` so the browser's broken-image chrome never shows.
- MarkdownEditor selection restore after toolbar wraps needs `requestAnimationFrame` before `setSelectionRange` because React re-renders the textarea value first.

### What warrants a second pair of eyes
- Sanitizer policy (D-5): protocol-relative `//` URLs are rejected — confirm that's wanted for internal CDN use cases.
- `ArticleListPanel` owns its archive/delete `ConfirmDialog` internally; if containers want to intercept, the confirm would need lifting.
- New organisms are exported from the npm barrel via `components/organisms/index.ts` — decide whether to hold them out of the next npm release until the API settles (design-doc §11 risk).

### What should be done in the future
- D1 (backend + web integration) and D2 (`cms.dsl`) per tasks.md; consider adopting `MediaThumb` inside `RichArticle`/`MarkdownArticle` image blocks behind a prop.

### Code review instructions
- Start at `packages/rag-evaluation-site/src/cms/types.ts`, then one component per layer (suggest `MediaThumb`, `DialogShell`, `UploadQueueList`, `ArticleEditorPanel`) against GUIDELINES.md.
- Validate: `pnpm --dir packages/rag-evaluation-site typecheck && pnpm biome check packages/rag-evaluation-site/src && pnpm --dir packages/rag-evaluation-site build-storybook`; browse `http://localhost:6007` → "Design System/Atoms/MediaThumb", "Component Library/Organisms/CmsShell" etc.; screenshots in `sources/screenshots/cms/`.

## Step 7: go-go-course study — delivery pivots to Widget-DSL-first (design doc revised)

The user redirected integration: no REST API or RTK Query; all product work happens in `./go-go-course`, which consumes the design system through the goja bindings. Studying that repo overturned a load-bearing assumption in the design doc — the widget server-action endpoint I had marked "unimplemented" exists there in JS, along with a working file-backed CMS. The design doc was revised in place (revision note, §6.7–6.9, D-6 superseded by D-8/D-9, §8 flows, §9 phases, §11 risks, §12 references) and re-uploaded.

### Prompt Context

**User prompt (verbatim):** "let's keep it at widget DSL, no need for API and RTK-query, we are doing all the work in ./go-go-course which uses the goja bindings. study that and see if that changesyour plans and you have design doc to update. If so, update it, and upload a copy to remarkable"

**Assistant interpretation:** Study go-go-course as the real CMS consumer, revise the integration plan to be Widget-DSL-only, update the design doc, re-upload to reMarkable.

**Inferred user intent:** Keep this repo's contribution at the widget/DSL layer; the application (storage, actions, pages) lives in go-go-course.

### What I did
- Read `go-go-course/cmd/go-go-course/`: `xgoja.yaml` (buildspec: `rag-widget-site` provider pinned v0.1.2; ui/data/context_window/course.dsl selected; embedded SPA assets), `server.js` (all 629 lines), `lib/pages/admin-course-cms.js` (all 251), `lib/pages/admin-common.js`, `lib/course-material-service.js` (upload/delete/list core), `webapp/src/main.tsx` + `webapp/package.json` (15-line shell around `RagEvaluationSiteApp apiBase="/api/widget"`, npm pin 0.1.16), plus `ContextUploadDropArea.widget.tsx` back in this repo (the file-serialization pattern).
- Rewrote design-doc sections: exec-summary revision note; §6.7 (consumer architecture + pipeline diagram), §6.8 (IR-consumption requirements: shared upload serialization, form `name`/`defaultValue` modes, navigate-based selection, `ActionSpec.confirm`, CmsAsset mapping), §6.9 (cms.dsl as the delivery path + target authoring example); D-6 → superseded, new D-8 (Widget-DSL-first, accepted) and D-9 (`ActionSpec.confirm`, proposed); §8 flows (upload via base64 action context; handout editing via form-post + adapter-internal live preview; navigate-driven selection); §9 Phases 5–6 replaced; §11 risks (SVG validator only checks for `<svg>`; release-loop pin-bump cost; IR capability gaps); §12 references (8 go-go-course files).
- Replaced tasks.md D1/D2 with P5.1–P5.7 (IR enablement here) and P6.1–P6.6 (go-go-course integration).

### What I learned (key findings)
- `POST /api/widget/actions/:name` is implemented in go-go-course (server.js:383–456): `upload-session`, `admin-upload-course-material`, `admin-delete-course-material`, `admin-reorder-course-agenda`, returning `{ok, refresh, toast, data}` — exactly the `ServerActionResult` shape actions.ts consumes. My "missing endpoint" claim was true for this repo's Go server but false for the product.
- File uploads already cross the JSON boundary: `ContextUploadDropArea.widget.tsx` serializes Files as utf8/base64 into the action context; `course-material-service.js` decodes, validates magic bytes (PNG/JPEG/GIF/WebP) or `<svg>`, sanitizes filenames, writes atomically to `course/{slides,handouts,media}`.
- Form-heavy editing is native form posts (`formPanel({method:"post", formAction})` → save → redirect `?status=saved`); cross-request state is query-parameter driven. This means my controlled-only `MarkdownEditor`/`TagListInput` need `name`/`defaultValue` modes, and selection/paging adapters should emit `action.navigate("?asset=$id")`.
- Deletes fire with **no confirmation** (`cell.actionButton("Delete", …)`) — motivated D-9 (`ActionSpec.confirm` handled centrally in `dispatchWidgetAction`) instead of trying to express dialog state in IR.
- Version pins gate delivery: npm 0.1.16 + Go v0.1.2 in go-go-course; `webapp` has `dev:rag-local`/`build:rag-local` scripts (`RAG_SITE_SRC` vite alias) for iterating against package source pre-publish.

### What warrants a second pair of eyes
- D-9's confirm-in-dispatch design (v1 `window.confirm`) — confirm UX vs. the flat aesthetic; v2 proposes a DialogShell-styled confirm mounted in `RagEvaluationSiteApp`.
- The claim that upload *progress* is out of scope over IR (single JSON POST) — if large media uploads matter, the express binding would need multipart/streaming first.

### Code review instructions
- Diff design-doc/01 around the "Revision (2026-07-03…)" note; verify go-go-course line anchors: `nl -ba go-go-course/cmd/go-go-course/server.js | sed -n '383,456p'` and `sed -n '1,60p' …/lib/pages/admin-common.js`.

## Step 8: Phase 5 implemented — cms.dsl, 14 widget adapters, ActionSpec.confirm, drift fixes

Implemented the whole IR-enablement phase in one pass: the cross-cutting action-confirm mechanism, shared upload serialization, form-mode props, 14 widget adapters + manifests under a new `cms.dsl` module on both the TS and Go sides, IR stories, and the previously documented DSL drift fixes.

### Prompt Context

**User prompt (verbatim):** "go ahead"

**Assistant interpretation:** Execute tasks P5.1–P5.7 (IR enablement in this repo).

**Inferred user intent:** Make the CMS widgets consumable from go-go-course's JS pages via the goja DSL.

### What I did
- **P5.1** `ActionSpecBase.confirm?: string` in `ir.ts` (all five action kinds extend it); `dispatchWidgetAction` gates on `window.confirm(interpolate(confirm, context))` *before* both custom handlers and built-in dispatch; `ConfirmedDestructiveAction` IR story.
- **P5.2** `src/widgets/uploadSerialization.ts` (SerializedUploadFile + serializeUploadFile(s), utf8/base64 heuristics) extracted from `ContextUploadDropArea.widget.tsx`, which now re-exports the type for compatibility.
- **P5.3** `name` on `MarkdownEditor` (forwarded to its textarea) and `TagListInput` (hidden `name=value` input with joined tags); `onQuerySubmit` added to `MediaLibraryPanel`/`ArticleListPanel` (wired to SearchField Enter).
- **P5.4/5.5** 14 `.widget.tsx` + `.widget.yaml` pairs (module `cms.dsl`, status experimental): atoms MediaThumb/Tag/ContentStatusBadge/MeterBar, layout TileGrid, molecules AssetTile/Breadcrumbs/Pagination/SearchField/EmptyState/MarkdownEditor, organisms MediaLibraryPanel/ArticleListPanel/CmsShell. Stateful adapters: SearchField (local value, dispatch on Enter) and MarkdownEditor (local value + live `MarkdownArticle` preview in a SplitPane; `preview:"hidden"` opts out). MediaLibraryPanel adapter reuses `serializeUploadFiles` for `onFilesSelectedAction` — same `{files, fileNames, fileCount}` contract as ContextUploadDropArea. New IR prop interfaces + `RagWidgetType`/`WidgetProps` additions in ir.ts; `"cms.dsl"` in the `WidgetModule` union; `cmsWidgetRegistry` merged into `defaultWidgetRegistry`; `WidgetRenderer.cms.stories.tsx` (7 stories: AtomsGallery, MediaLibraryFromIr, ArticleListFromIr, MarkdownEditorWithLivePreview, ConfirmedDestructiveAction, CmsShellFromIr, BreadcrumbsPaginationSearch).
- **P5.6** Go: `CmsModuleName`, `cmsHelpers` (14), moduleSpec (action + recipes), `mediaLibraryRecipe`/`articleListRecipe` (normalizeActionSpec on all `on*` options, so bare strings become server actions and explicit specs — including `confirm` — pass through); provider entry in `widgetsite/provider.go`; `TestCmsModuleExportsHelpersRecipesAndBoundaries` (helper presence, ui/cms boundary, recipe expansion, confirm passthrough).
- **P5.7** typescript.go: `cell.linkButton`, `cell.actionButton`, `action.download` declared; `contextGroupedStripDiagram` helper added to context_window.dsl (React adapter already existed); doc/02 API reference updated (cms.dsl section with recipes + CmsAsset/CmsArticleSummary shapes, action.download, confirm docs, cell additions, grouped-strip helper).

### What worked
- Everything passed on the first full run: tsc, biome (after auto-fixes), build-storybook, `go test ./pkg/widgetdsl ./pkg/xgoja/providers/widgetsite`, and a Playwright smoke of all 7 IR stories (0 page errors). The media library screenshot renders the complete CMS panel purely from JSON IR.

### What didn't work
- gopls emits `go.work requires go >= 1.26.4 (running go 1.25.5)` diagnostics (the workspace go.work spans go-go-course's go 1.26.4) — editor-only noise; the shell Go toolchain builds and tests fine.

### What was tricky to build
- `EmptyState`'s React `action` slot prop collides with `BaseWidgetProps.action` (ActionSpec) — the widget prop is named `actionSlot` instead.
- Widget adapters can't use hooks in `render()` directly; SearchField/MarkdownEditor render a `*WidgetHost` inner component that owns `useState`.

### What warrants a second pair of eyes
- `confirm` fires before custom `onAction` handlers too (design choice: it's part of the action contract). Storybook's default logger therefore also triggers the prompt — intended, but worth confirming.
- MediaLibraryPanel adapter passes a no-op `onQueryChange` when only `onQuerySubmitAction` is set (the React panel requires onQueryChange to show the search box); acceptable shim or should the panel accept submit-only?

### What should be done in the future
- Phase 6 (go-go-course): bump pins, CmsAsset mapping in course-material-service, admin page upgrade, handout editor page, end-to-end smoke (tasks P6.1–P6.6).

### Code review instructions
- Start: `src/widgets/ir.ts` (ActionSpecBase + cms props), `src/widgets/uploadSerialization.ts`, `MediaLibraryPanel.widget.tsx`, `pkg/widgetdsl/module.go` (cmsHelpers + recipes).
- Validate: `pnpm --dir packages/rag-evaluation-site typecheck && pnpm biome check packages/rag-evaluation-site/src`; `go test ./pkg/widgetdsl ./pkg/xgoja/providers/widgetsite`; Storybook → "Widget IR/Renderer/CMS"; screenshots in `sources/screenshots/cms/widget-ir-*`.
