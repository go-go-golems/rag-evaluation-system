---
Title: "Widget DSL Getting Started"
Slug: widget-dsl-getting-started
Short: "Author React-rendered Widget IR pages in xgoja with split Widget DSL modules."
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

The provider exposes clean-break domain modules:

| Module | Owns |
| --- | --- |
| `ui.dsl` | page wrapper, text/element/component helpers, generic layout, primitive, foundation, and UI recipes |
| `data.dsl` | legacy/current `DataTable`, `cell.*` helpers, and data recipes |
| `data.v2.dsl` | experimental hard-cutover typed/fluent data builders for new table/editor examples |
| `context_window.dsl` | context-window diagrams, transcript, annotation, comment, and upload helpers |
| `course.dsl` | course, slide, handout, course-studio helpers, and `contextStudioNavIcon` |

Select the modules you use in `xgoja.yaml`:

```yaml
packages:
  - id: rag-widget-site
    import: github.com/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite

modules:
  - package: rag-widget-site
    name: ui.dsl
    as: ui.dsl
  - package: rag-widget-site
    name: data.dsl
    as: data.dsl
  - package: rag-widget-site
    name: data.v2.dsl
    as: data.v2.dsl
  - package: rag-widget-site
    name: context_window.dsl
    as: context_window.dsl
  - package: rag-widget-site
    name: course.dsl
    as: course.dsl
```

For local development, add a `replace` entry that points to the RAG repository root.

## Smallest page

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

`ui.page(...)` returns a page object with `schemaVersion`, `id`, `title`, and `root`. Strings and numbers used as children are normalized into text nodes.

## Data table example

For new hard-cutover examples, prefer `data.v2.dsl`:

```js
const data = require("data.v2.dsl")

const schema = data.schema("Query")
  .field("id", data.f.key().label("ID"))
  .field("name", data.f.primary().label("Name"))
  .field("status", data.f.status().label("Status"))
  .build()

return data.collection("queries", rows)
  .schema(schema)
  .table()
  .toIR()
```

The older direct `data.dsl` table form is still supported for existing pages:

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

```js
__package__({ name: "sites", short: "WidgetRenderer sites" })

__verb__("demo", { name: "demo", output: "text", short: "Serve a WidgetRenderer demo" })
function demo() {
  const express = require("express")
  const assets = require("fs:assets")
  const ui = require("ui.dsl")

  const app = express.app()
  app.spaFromAssetsModule("/", assets, "/app/public", {
    excludePrefixes: ["/api", "/healthz", "/favicon.ico"]
  })

  app.get("/healthz", (_req, res) => res.json({ ok: true }))
  app.get("/api/widget/pages/demo", (_req, res) => {
    res.json(ui.page({
      id: "demo",
      title: "Demo",
      root: ui.panel({ title: "Demo" },
        ui.caption({ tone: "muted" }, "Rendered by React"),
        ui.button({ variant: "primary" }, "Refresh")
      )
    }))
  })
}
```

## What to remember

- DSL constructors return JSON-compatible Widget IR; they do not return HTML strings or React elements.
- `ui.dsl` owns `page(...)`.
- `data.dsl` owns `cell.*` helpers.
- `context_window.dsl` owns transcript, annotation, anchored-comment, context diagram, and upload helpers.
- `course.dsl` owns course, slide, handout, and course-studio helpers.
- Import only the domain modules you use; there is no compatibility bucket module.

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| `Cannot find module "ui.dsl"` | The module was not selected in `xgoja.yaml`. | Add a `modules:` entry for package `rag-widget-site`, name `ui.dsl`, alias `ui.dsl`. |
| `xgoja build` tries to fetch the local provider from GitHub. | The provider package is local but the build spec has no `replace`. | Add `replace: ../../..` or another path to the RAG module root. |
| The browser route `/pages/demo` returns `404`. | The React app is served as static files without SPA fallback. | Use `app.spaFromAssetsModule("/", assets, "/app/public", { excludePrefixes: ["/api"] })`. |
| API routes return `index.html`. | The root SPA static handler is catching `/api/...`. | Add `/api` to `excludePrefixes`. |

## See Also

- `widget-dsl-js-api-reference`
- `widget-dsl-spa-bundling`
- `tutorial-http-serve-jsverbs`
- `tutorial-static-assets-http-server`
- `buildspec-reference`
