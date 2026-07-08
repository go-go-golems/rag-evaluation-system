---
Title: Implementation diary
Ticket: DOODLE-1
Status: active
Topics:
    - xgoja
    - widget-dsl
    - sqlite
    - doodle
DocType: diary
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-07T13:26:04.973177092-04:00
WhatFor: ""
WhenToUse: ""
---

# DOODLE-1 — Implementation diary

Goal: build a small **Doodle-style scheduling site** (create a poll with proposed
time slots; participants vote yes/no/maybe; a results grid tallies the best slot),
using:

- **SQLite** as the datastore (via the `db` goja module),
- the **rag Widget DSL** (`ui.dsl` / `data.dsl` from the `rag-widget-site` xgoja
  provider) to produce Widget IR that the React `RagEvaluationSiteApp` renders,
- **xgoja** to compile a custom goja binary that serves the site over HTTP.

## Step 1 — Environment discovery

- `xgoja` is installed at `~/.local/bin/xgoja`; `rag` is *not* a standalone binary —
  the "rag DSL" is the **Widget DSL** shipped by the `rag-widget-site` provider in
  `rag-evaluation-system/pkg/xgoja/providers/widgetsite`.
- Found a complete reference: `examples/xgoja/widget-site/` (spec + `verbs/sites.js`
  + prebuilt SPA under `assets/public/`). This is the template to adapt.
- Widget DSL docs read: `01-widget-dsl-getting-started.md`, `02-...js-api-reference.md`.
  Key helpers: `ui.page`, `ui.panel`, `ui.textInput/textareaInput/selectInput/formRow`,
  `ui.button`, `ui.metadataGrid`, `ui.recipes.metrics`, `data.dataTable`, `data.cell.*`,
  `data.recipes.masterDetailTable`, `ui.action.{navigate,server}`.

## Step 2 — Toolchain reality-check (important gotchas)

1. **The example spec is a legacy xgoja spec.** The locally-built `xgoja` is
   **v2-only**; `xgoja doctor -f xgoja.yaml` fails with "appears to be a legacy
   xgoja spec". Fix: `xgoja migrate-spec` → `schema: xgoja/v2`. Captured the v2
   shape (providers / runtime.modules / sources / commands / artifacts,
   `workspace.mode: auto`).

2. **go-go-goja resolves to the LOCAL workspace copy** (via `go.work`), confirmed by
   `xgoja doctor` (`resolution_kind: workspace`, `go-work`). That copy has the
   **new Express API**, where the old two-arg `app.get(path, handler)` was *removed*
   and now **panics**:
   `app.<verb>(pattern, handler) was removed; use app.<verb>(pattern).public().handle(handler)`.
   → The example `sites.js` (two-arg form) would panic; I must use the planned-route
   API. Request shape (from `modules/express` tests):
   - route: `app.get(p).public().handle((ctx, res) => …)`
   - path params: `ctx.params.x`
   - query: `ctx.request.query.x`
   - body (parsed JSON): `ctx.body`
   - response: `res.json(...)`, `res.status(n).json(...)`, `res.type(...).send(...)`.

3. **Query strings ARE forwarded by the SPA.** `App.tsx` builds the fetch URL as
   `/api/widget/pages/${pageId}${window.location.search}` and `ui.action.navigate`
   does `pushState(params, "", target)`. So `navigate("/pages/poll?poll=3")` →
   fetch `/api/widget/pages/poll?poll=3`, readable via `ctx.request.query.poll`.
   Interpolation supports `${row.slug}` / `$value` against the action context.

## Step 3 — Design (Doodle site)

SQLite schema: `polls`, `options` (time slots), `participants`, `votes`
(value ∈ yes|no|maybe, PK participant+option). File-backed DB for persistence.

Pages (Widget IR at `/api/widget/pages/<id>`):
- `index` — metrics + table of polls; "New poll" navigates to `create`.
- `create` — form (title/description/location/options) → server action `create-poll`.
- `poll?poll=<id>` — metadata + results grid (participants × slots + tallies +
  best-slot highlight) + "add availability" form → server action `cast-vote`.

Next: scaffold `examples/xgoja/doodle-site/`, write the jsverb with the planned-route
API, build, run, verify with curl + Playwright screenshot.

## Step 4 — Build & implement

Scaffolded `examples/xgoja/doodle-site/`:
- `xgoja.v2.yaml` — providers host+http+rag-widget-site, modules `express`/`fs:assets`/
  `db` (sqlite3, `file:doodle.db?_foreign_keys=on`, file-backed) / `ui.dsl` / `data.dsl`.
- `verbs/doodle.js` — package `doodle`, verb `site`. SQLite schema
  (`polls`/`options`/`participants`/`votes`) with `CREATE TABLE IF NOT EXISTS` +
  first-run seed. Pages `index` / `create` / `poll`. Native `<form>` POST handlers
  `/api/form/create-poll` and `/api/form/cast-vote` → `res.redirect(303, ...)`.
- `Makefile` (doctor/build/serve/sync-app), copied SPA into `assets/public/`.

`xgoja doctor` → all ok. `xgoja build` → 81 MB binary, exit 0.

Gotchas fixed while writing:
- `ui.formRow` reads `control` from **props**, not children — moved control into the
  props object.
- `ui.textInput` / `ui.textareaInput` default to **`readOnly: true`** — passed
  `readOnly: false` on every editable input, otherwise the fields render locked.

## Step 5 — Verify (curl)

`serve doodle site --http-listen 127.0.0.1:18793`. Verified: healthz ok; index page
renders seeded poll; `poll?poll=1` shows grid + best-slot; `POST /api/form/create-poll`
(urlencoded) → `303 → /pages/poll?poll=2`; `POST /api/form/cast-vote?poll=1` → 303 and
the new participant appears. SQLite `doodle.db` persisted across a server restart
(data survived) — confirms real file-backed persistence.

## Step 6 — Verify (Playwright) + the stale-SPA fix

The **embedded SPA bundle was stale**. The bundle I first copied from
`examples/xgoja/widget-site/assets` (`index-BWguF5a5.js`) predates several widgets:
the browser showed `Unsupported cell` for a `data.cell.actionButton` column and
`Unknown widget: FormPanel` for the create form.

Fix: rebuilt the **current** SPA from source — `packages/rag-evaluation-site` via its
own `build:app` script (`vite build --config vite.app.config.ts`) → `app-dist/`
(`index-jD_m6sTj.js`). **No React/TS source was modified**; only the build output was
regenerated. Synced `app-dist` into `assets/public/` (stripped source maps) and
rebuilt the binary. Also swapped the index page's per-row `actionButton` cell for
plain nav buttons (reads cleanly; works regardless of bundle age).

Browser flow (screenshots in `various/`):
- `/pages/create` → FormPanel renders with editable inputs; filled + submitted →
  redirected to `/pages/poll?poll=3` (new "Sprint retro" poll).
- Poll page: metadata grid, empty availability grid, results-by-slot table, and a
  vote form with a `<select>` per slot.
- Cast a vote (Rosalind: yes/no/yes) → grid populates, tallies update, the two
  top slots flip to `✔ succeeded` (best), score = yes*2 + maybe.
- Index: live metrics (3 polls / 5 responses / 10 slots), polls table, "Open a poll"
  nav buttons; clicking one client-navigates to `/pages/poll?poll=1`.

**Result: complete Doodle-style scheduling site — SQLite + rag Widget DSL + xgoja —
verified end-to-end in a real browser.**

### How to run
```
cd rag-evaluation-system/examples/xgoja/doodle-site
make serve            # builds + serves on 127.0.0.1:18793
# or: ./dist/doodle-site serve doodle site --http-listen 127.0.0.1:18793
```
