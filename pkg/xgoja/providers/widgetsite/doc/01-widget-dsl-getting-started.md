---
Title: "Widget DSL Getting Started"
Slug: widget-dsl-getting-started
Short: "Author React-rendered Widget IR pages in xgoja with the v3-first Widget DSL."
Topics:
- xgoja
- widget-dsl
- rag-widget-site
- widget-ir
- react
Commands:
- xgoja build
- xgoja doctor
- xgoja list-modules
- serve
Flags:
- --http-listen
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

This tutorial explains how to use the `rag-widget-site` xgoja provider to write JavaScript that produces Widget IR. Widget IR is JSON-compatible data rendered by the React `RagEvaluationSiteApp`; JavaScript authors describe page structure and data, while React owns CSS modules, accessibility, event handling, and component behavior.

## Modules

The provider exposes the new parallel `widget.dsl` v3 module and the legacy split modules. Prefer `widget.dsl` for new pages; keep the split modules selected for existing scripts until those scripts are ported.

| Module | Owns | Status |
| --- | --- | --- |
| `widget.dsl` | v3 namespaces: `raw`, `act`, `bind`, `ui`, `data`, `crm`, `cms`, `course`, `context`, `schedule`, `time`, and `style` | Preferred for new work |
| `ui.dsl` | page wrapper, text/element/component helpers, generic layout, primitive, foundation, and UI recipes | Legacy split module |
| `data.dsl` | legacy/current `DataTable`, `cell.*` helpers, and data recipes | Legacy split module |
| `data.v2.dsl` | typed/fluent data builders used by pre-v3 table/editor examples | Legacy split module |
| `context_window.dsl` | context-window diagrams, transcript, annotation, comment, and upload helpers | Legacy split module |
| `course.dsl` | course, slide, handout, course-studio helpers, and `contextStudioNavIcon` | Legacy split module |
| `cms.dsl` | media, asset, and article-management helpers | Legacy split module |

Select the modules you use in `xgoja.yaml`. A new v3 host usually selects only `widget.dsl` from this provider:

```yaml
packages:
  - id: rag-widget-site
    import: github.com/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite

modules:
  - package: rag-widget-site
    name: widget.dsl
    as: widget.dsl
```

A migration host may select `widget.dsl` alongside `ui.dsl`, `data.dsl`, `data.v2.dsl`, `context_window.dsl`, `course.dsl`, and `cms.dsl` while pages are ported. Remove legacy entries only after the scripts no longer import them. For local development, add a `replace` entry that points to the RAG repository root.

## Smallest `widget.dsl` v3 page

```js
const widget = require("widget.dsl")

const page = widget.page("Demo", (p) =>
  p.id("demo").section("Demo", (s) =>
    s.caption("Rendered by React from Widget IR").view(
      widget.ui.button("Refresh", widget.act.event("refresh"), { variant: "primary" })
    )
  )
)
```

`widget.page(...)` returns a page object with `schemaVersion`, `id`, `title`, and `root`. Strings and numbers used as children are normalized into text nodes.

## Legacy split-module smallest page

```js
const ui = require("ui.dsl")

const page = ui.page({
  id: "demo",
  title: "Demo",
  root: ui.panel({ title: "Demo" },
    ui.statusText({ status: "succeeded", icon: true }, "Rows: 2")
  )
})
```

`ui.page(...)` is still supported for existing split-module pages. For new pages, use `widget.page(...)` from `widget.dsl`.

## Data collection example

New v3 pages keep page composition and data widgets under one `widget.dsl` import. Build a schema with `widget.data`, then return the collection node from a section:

```js
const widget = require("widget.dsl")

const schema = widget.data.fields("Query", (f) =>
  f.key("id").primary("name").status("status")
).build()

const table = widget.data.collection("queries", rows, (c) =>
  c.schema(schema).table()
).toNode()

const page = widget.page("Queries", (p) =>
  p.section("Queries", (s) => s.view(table))
)
```

`data.v2.dsl` and the older direct `data.dsl` table form remain available for existing pages during migration:

```js
const ui = require("ui.dsl")
const data = require("data.dsl")

const rows = db.query("SELECT id, name, status FROM queries ORDER BY id")

return ui.page({
  id: "queries",
  title: "Queries",
  sections: [
    ui.panel({ title: "Queries" },
      data.dataTable({
        rows,
        getRowKey: "id",
        columns: [
          { id: "id", header: "ID", cell: data.cell.field("id") },
          { id: "name", header: "Name", cell: data.cell.field("name") },
          { id: "status", header: "Status", cell: data.cell.status("status") }
        ]
      })
    )
  ]
})
```

The table is still rendered by React. The JavaScript author supplies serializable rows and serializable cell specifications.

## Semantic page example

```js
const ui = require("ui.dsl")
const contextWindow = require("context_window.dsl")
const course = require("course.dsl")

const styleSet = contextWindow.paletteStyleSet({
  palette: "Dusty Magenta / Blue",
  entries: [
    { id: "system", label: "System", accent: "b", pattern: "checker" },
    { id: "retrieval", label: "Retrieved docs", accent: "a", pattern: "stipple" }
  ]
})

const snapshot = contextWindow.contextSnapshot({
  id: "ctx",
  title: "Context Window",
  limit: 16000,
  parts: [
    contextWindow.contextPart("system", "System", "system", 600),
    contextWindow.contextPart("retrieval", "Retrieved docs", "retrieval", 7000)
  ]
})

const slide = {
  id: "budget",
  number: "01",
  title: "Budget pressure",
  view: "budget",
  snapshotId: snapshot.id,
  notes: ["Retrieved documents dominate this example."]
}

return ui.page({
  id: "semantic",
  title: "Semantic context page",
  sections: [
    contextWindow.recipes.contextDiagram({ snapshot, styleSet, view: "budget" }),
    course.recipes.courseStudio({
      sections: [{ id: "course", label: "Course", items: [{ id: "slides", label: "Slides" }] }],
      activeItemId: "slides",
      main: course.recipes.courseSlide({ slide, snapshot, index: 0, total: 1 })
    })
  ]
})
```

## Serve from a jsverb

An xgoja route returns page JSON while the embedded SPA renders it in the browser. Keep the SPA fallback away from `/api` routes and use the same v3 module as the page code:

```js
__package__({ name: "sites", short: "WidgetRenderer sites" })

__verb__("demo", { name: "demo", output: "text", short: "Serve a WidgetRenderer demo" })
function demo() {
  const express = require("express")
  const assets = require("fs:assets")
  const widget = require("widget.dsl")

  const app = express.app()
  app.spaFromAssetsModule("/", assets, "/app/public", {
    excludePrefixes: ["/api", "/healthz", "/favicon.ico"]
  })

  app.get("/healthz", (_req, res) => res.json({ ok: true }))
  app.get("/api/widget/pages/demo", (_req, res) => {
    res.json(widget.page({ id: "demo", title: "Demo" }, (p) =>
      p.section("Demo", (s) => s.view(
        widget.ui.card({ title: "Demo" },
          widget.ui.caption("Rendered by React"),
          widget.ui.button("Refresh", widget.act.event("refresh"), { variant: "primary" })
        )
      ))
    ).toPage())
  })
}
```

## What to remember

- DSL constructors return JSON-compatible Widget IR; they do not return HTML strings or React elements.
- `widget.dsl` is the preferred module for new work and exposes all v3 namespaces from one import.
- The split modules still work for existing pages: `ui.dsl` owns old page helpers, `data.dsl` owns old `cell.*` helpers, `context_window.dsl` owns old context helpers, `course.dsl` owns old course helpers, and `cms.dsl` owns old CMS helpers.
- Import only the modules you use; there is no compatibility bucket module.

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| `Cannot find module "widget.dsl"` | The v3 module was not selected in `xgoja.yaml`. | Add a `modules:` entry for package `rag-widget-site`, name `widget.dsl`, alias `widget.dsl`. |
| `Cannot find module "ui.dsl"` | A legacy script imports the split UI module but the module was not selected. | Keep the legacy `ui.dsl` entry while migrating, or port the script to `widget.dsl`. |
| `xgoja build` tries to fetch the local provider from GitHub. | The provider package is local but the build spec has no `replace`. | Add `replace: ../../..` or another path to the RAG module root. |
| The browser route `/pages/demo` returns `404`. | The React app is served as static files without SPA fallback. | Use `app.spaFromAssetsModule("/", assets, "/app/public", { excludePrefixes: ["/api"] })`. |
| API routes return `index.html`. | The root SPA static handler is catching `/api/...`. | Add `/api` to `excludePrefixes`. |

## See Also

- `widget-dsl-v3-examples` — runnable v3 composition, scheduling, CRM, and action recipes.
- `widget-dsl-v3-api-reference` — descriptor-derived namespace inventory.
- `widget-dsl-js-api-reference` — legacy module details and migration reference.
- `widget-dsl-spa-bundling` — serve Widget IR from an xgoja application.
- `tutorial-http-serve-jsverbs`
- `tutorial-static-assets-http-server`
- `buildspec-reference`
