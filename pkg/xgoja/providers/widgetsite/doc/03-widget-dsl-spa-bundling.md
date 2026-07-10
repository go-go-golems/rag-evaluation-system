---
Title: "Widget DSL SPA Bundling"
Slug: widget-dsl-spa-bundling
Short: "Serve React WidgetRenderer SPA assets for v3-first Widget DSL applications."
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

Widget DSL modules create Widget IR only. They do not render HTML and they do not include browser assets by themselves. A browser-facing xgoja application must also serve a React WidgetRenderer SPA and expose API routes that return Widget IR pages/actions. New hosts use `widget.dsl`; legacy split modules remain available only while a host is being migrated.

## Runtime model

1. JavaScript imports `widget.dsl` for new pages, or one or more legacy split modules such as `ui.dsl`, `data.dsl`, `context_window.dsl`, and `course.dsl` while existing pages are being migrated.
2. API routes return Widget IR page JSON from `/api/widget/pages/{id}`.
3. The React SPA fetches those pages and renders them through the registry-backed WidgetRenderer.
4. Actions are dispatched to browser events, navigation, copy behavior, or `/api/widget/actions/{name}` depending on action kind.

## Build spec module selection

Select only the modules your scripts need. A new v3 app usually selects `widget.dsl` as the only module from `rag-widget-site`:

```yaml
providers:
  - id: rag-widget-site
    import: github.com/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite
    register: Register

runtime:
  modules:
    - provider: rag-widget-site
      name: widget.dsl
      as: widget.dsl
```

A migration app can temporarily select `widget.dsl` beside the legacy split modules:

```yaml
runtime:
  modules:
    - provider: rag-widget-site
      name: widget.dsl
      as: widget.dsl
    - provider: rag-widget-site
      name: ui.dsl
      as: ui.dsl
    - provider: rag-widget-site
      name: data.dsl
      as: data.dsl
    - provider: rag-widget-site
      name: context_window.dsl
      as: context_window.dsl
    - provider: rag-widget-site
      name: course.dsl
      as: course.dsl
    - provider: rag-widget-site
      name: cms.dsl
      as: cms.dsl
```

Remove legacy module entries once the host scripts no longer import them.

## Include provider help in an xgoja application

The provider ships Glazed help entries with the generated application; adding the help source makes topics such as `widget-dsl-v3-examples` available through that application's `help` command. This is independent of runtime module selection, so include both when the host should teach its own DSL.

```yaml
sources:
  - id: widget-dsl
    kind: help
    from:
      provider:
        provider: rag-widget-site
        source: widget-dsl
```

The Doodle and workshop CRM examples use this source alongside their `widget.dsl` runtime module.

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
  const widget = require("widget.dsl")

  const app = express.app()
  app.spaFromAssetsModule("/", assets, "/app/public", {
    excludePrefixes: ["/api", "/healthz", "/favicon.ico"]
  })

  app.get("/healthz", (_req, res) => res.json({ ok: true }))
  app.get("/api/widget/pages/demo", (_req, res) => {
    res.json(widget.page("Demo", (p) =>
      p.id("demo").section("Demo", (s) =>
        s.caption("Rendered by React").view(
          widget.ui.button("Refresh", widget.act.event("refresh"), { variant: "primary" })
        )
      )
    ))
  })
}
```

The `excludePrefixes` option is important when mounting the SPA at `/`; otherwise the static fallback can answer API requests before dynamic routes see them.

## Data/context/course APIs

In v3, use one module and compose across namespaces:

```js
const widget = require("widget.dsl")

app.get("/api/widget/pages/course", (_req, res) => {
  res.json(widget.page("Course", (p) => {
    p.id("course")
    p.section("Context", (s) =>
      s.view(widget.context.diagram(snapshot, (d) => d.view("budget")))
    )
    p.section("Studio", (s) =>
      s.view(widget.course.shell({ title: "Course", sections }, (shell) =>
        shell.active("slides").main(widget.course.slideDeck({ slides, snapshot, index: 0 }))
      ))
    )
    p.section("Rows", (s) =>
      s.view(widget.data.collection("rows", rows, (c) => c.schema(schema).table()).toNode())
    )
  }))
})
```

Legacy endpoints can keep importing `ui.dsl`, `data.dsl`, `context_window.dsl`, `course.dsl`, and `cms.dsl` until they are ported.

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| `Cannot find module "widget.dsl"` | The v3 module was not selected in the xgoja build spec. | Add a runtime module entry for package/provider `rag-widget-site`, name `widget.dsl`, alias `widget.dsl`. |
| `Cannot find module "ui.dsl"` | A legacy endpoint imports the split UI module but the build spec did not select it. | Keep `ui.dsl` selected while migrating, or port the endpoint to `widget.dsl`. |
| `Cannot find module "data.dsl"` | A legacy endpoint imports the data module but the build spec did not select it. | Keep `data.dsl` selected while migrating, or port the endpoint to `widget.data`. |
| Browser routes such as `/pages/demo` return `404`. | Static files are served without SPA fallback. | Use `defaultspa.Handler()` in a Go host, or `app.spaFromAssetsModule(...)` in an xgoja host. |
| API routes return `index.html`. | The root SPA static handler is catching `/api/...`. | Add `/api` to `excludePrefixes` or register API routes before the SPA fallback. |
| The npm package works in React but the xgoja binary has no UI. | DSL modules only create Widget IR; they do not include browser assets. | Bundle the default SPA assets or build a host frontend that imports `RagEvaluationSiteApp`. |

## See Also

- `widget-dsl-getting-started` — select the provider module and build a first v3 page.
- `widget-dsl-v3-examples` — composition, scheduling, CRM, and action recipes.
- `widget-dsl-v3-api-reference` — descriptor-derived v3 namespace inventory.
- `widget-dsl-js-api-reference` — action contracts and legacy module details.
- `tutorial-static-assets-http-server`
- `buildspec-reference`
