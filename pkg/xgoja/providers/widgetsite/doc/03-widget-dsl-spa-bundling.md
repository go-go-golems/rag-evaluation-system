---
Title: "Widget DSL SPA Bundling"
Slug: widget-dsl-spa-bundling
Short: "Serve React WidgetRenderer SPA assets for split Widget DSL pages."
Topics:
- xgoja
- widget-dsl
- widget-ir
- spa
- react
- assets
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

The split Widget DSL modules create Widget IR only. They do not render HTML and they do not include browser assets by themselves. A browser-facing xgoja application must also serve a React WidgetRenderer SPA and expose API routes that return Widget IR pages/actions.

## Runtime model

1. JavaScript imports one or more domain modules such as `ui.dsl`, `data.dsl`, `context_window.dsl`, and `course.dsl`.
2. API routes return Widget IR page JSON from `/api/widget/pages/{id}`.
3. The React SPA fetches those pages and renders them through the registry-backed WidgetRenderer.
4. Actions are dispatched to browser events, navigation, copy behavior, or `/api/widget/actions/{name}` depending on action kind.

## Build spec module selection

Select only the modules your scripts need:

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
```

A course/context app can also select:

```yaml
modules:
  - package: rag-widget-site
    name: context_window.dsl
    as: context_window.dsl
  - package: rag-widget-site
    name: course.dsl
    as: course.dsl
```

## Serve embedded SPA assets from a Go host

A Go application can mount the default embedded SPA and own its API routes:

```go
mux := http.NewServeMux()
mux.Handle("GET /api/widget/pages/{id}", pagesHandler)
mux.Handle("POST /api/widget/actions/{name}", actionsHandler)
mux.Handle("/", defaultspa.Handler())
```

The SPA fallback should be mounted after API routes.

## Serve copied SPA assets from xgoja

An xgoja-generated binary needs static assets. Build the React app, copy its output into the app asset tree, expose that tree through the host asset module, and mount it with `spaFromAssetsModule`:

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

The `excludePrefixes` option is important when mounting the SPA at `/`; otherwise the static fallback can answer API requests before dynamic routes see them.

## Data/context/course APIs

Use multiple modules when an endpoint needs multiple domains:

```js
const ui = require("ui.dsl")
const data = require("data.dsl")
const contextWindow = require("context_window.dsl")
const course = require("course.dsl")

app.get("/api/widget/pages/course", (_req, res) => {
  res.json(ui.page({
    id: "course",
    title: "Course",
    sections: [
      contextWindow.recipes.contextDiagram({ snapshot, view: "budget" }),
      course.recipes.courseStudio({
        sections,
        activeItemId: "slides",
        main: course.recipes.courseSlide({ slide, snapshot, index: 0, total: 1 })
      }),
      ui.panel({ title: "Rows" }, data.dataTable({ rows, getRowKey: "id", columns }))
    ]
  }))
})
```

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| `Cannot find module "ui.dsl"` | The module was not selected in the xgoja build spec. | Add a `modules:` entry for package `rag-widget-site`, name `ui.dsl`, alias `ui.dsl`. |
| `Cannot find module "data.dsl"` | The endpoint imports the data module but the build spec did not select it. | Add `data.dsl` to the selected modules. |
| Browser routes such as `/pages/demo` return `404`. | Static files are served without SPA fallback. | Use `defaultspa.Handler()` in a Go host, or `app.spaFromAssetsModule(...)` in an xgoja host. |
| API routes return `index.html`. | The root SPA static handler is catching `/api/...`. | Add `/api` to `excludePrefixes` or register API routes before the SPA fallback. |
| The npm package works in React but the xgoja binary has no UI. | DSL modules only create Widget IR; they do not include browser assets. | Bundle the default SPA assets or build a host frontend that imports `RagEvaluationSiteApp`. |

## See Also

- `widget-dsl-getting-started`
- `widget-dsl-js-api-reference`
- `tutorial-static-assets-http-server`
- `buildspec-reference`
