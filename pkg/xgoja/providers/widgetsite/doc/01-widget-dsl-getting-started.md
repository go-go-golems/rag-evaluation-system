---
Title: "Widget DSL Getting Started"
Slug: widget-dsl-getting-started
Short: "Author React-rendered Widget IR pages in xgoja with require(\"widget.dsl\")."
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

This tutorial explains how to use the `rag-widget-site` xgoja provider to write JavaScript that produces Widget IR. Widget IR is JSON-compatible data. It is not HTML, and it is not a Go-side imitation of the React component tree. The generated xgoja binary runs JavaScript, the JavaScript returns Widget IR, and the React `RagEvaluationSiteApp` renders that IR with the package's real component library.

The boundary matters. JavaScript authors describe a page as data; React owns rendering, CSS modules, event handling, table behavior, and accessibility details. This separation keeps the xgoja script small while preserving the frontend system that already exists.

## The smallest useful page

A Widget DSL script starts by requiring the module selected in `xgoja.yaml`:

```js
const rag = require("widget.dsl")
```

The `rag` object contains constructors for Widget IR nodes. A constructor returns a plain JavaScript object that can be serialized as JSON. For example, this code creates a panel with one status line:

```js
const page = {
  schemaVersion: "0.1.0",
  id: "demo",
  title: "Demo",
  root: rag.panel({ title: "Demo" },
    rag.statusText({ status: "succeeded", icon: true }, "Rows: 2")
  )
}
```

The resulting object has the shape expected by the React application:

```json
{
  "schemaVersion": "0.1.0",
  "id": "demo",
  "title": "Demo",
  "root": {
    "kind": "component",
    "type": "Panel",
    "props": { "title": "Demo" },
    "children": [
      {
        "kind": "component",
        "type": "StatusText",
        "props": { "status": "succeeded", "icon": true },
        "children": [{ "kind": "text", "text": "Rows: 2" }]
      }
    ]
  }
}
```

The page wrapper contains metadata. The `root` field is the Widget IR tree. `rag.panel(...)` creates a `Panel` component node, and the string child is normalized into a text node. You do not need to write the text-node object manually unless you want to.

## Select the provider in xgoja.yaml

A generated binary can only require modules that are selected by its build specification. The widget provider is selected as a normal xgoja package:

```yaml
packages:
  - id: rag-widget-site
    import: github.com/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite

modules:
  - package: rag-widget-site
    name: widget.dsl
    as: widget.dsl
```

For local development, add a `replace` entry that points to the RAG repository root:

```yaml
packages:
  - id: rag-widget-site
    import: github.com/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite
    replace: ../../..
```

The module name and alias are usually the same. The provider also exposes `rag.dsl` as an alias for RAG-oriented scripts:

```yaml
modules:
  - package: rag-widget-site
    name: rag.dsl
    as: rag.dsl
```

Use one alias consistently inside a script. The two modules expose the same helpers.

## Serve a WidgetRenderer page from a jsverb

The HTTP provider supplies `require("express")`. The host provider can expose embedded assets as `require("fs:assets")`. A WidgetRenderer site usually combines those modules with `widget.dsl`:

```js
__package__({ name: "sites", short: "WidgetRenderer sites" })

__verb__("demo", {
  name: "demo",
  output: "text",
  short: "Serve a WidgetRenderer demo"
})
function demo() {
  const express = require("express")
  const assets = require("fs:assets")
  const rag = require("widget.dsl")

  const app = express.app()
  app.spaFromAssetsModule("/", assets, "/app/public", {
    excludePrefixes: ["/api", "/healthz", "/favicon.ico"]
  })

  app.get("/healthz", (_req, res) => res.json({ ok: true }))
  app.get("/api/widget/pages/demo", (_req, res) => {
    res.json({
      schemaVersion: "0.1.0",
      id: "demo",
      title: "Demo",
      root: rag.panel({ title: "Demo" },
        rag.caption({ tone: "muted" }, "Rendered by React"),
        rag.button({ variant: "primary" }, "Refresh")
      )
    })
  })
}
```

The SPA route serves the React application. The API route serves Widget IR. The `excludePrefixes` option is important when the SPA is mounted at `/`; without it, the static handler would answer API requests before the dynamic routes see them.

## Use data from a database

The Widget DSL does not know where data came from. It accepts plain JavaScript objects. That means database rows, HTTP results, computed summaries, and constants all enter the renderer through the same Widget IR boundary.

```js
const rows = db.query("SELECT id, name, status FROM queries ORDER BY id")

return {
  schemaVersion: "0.1.0",
  id: "queries",
  title: "Queries",
  root: rag.panel({ title: "Queries" },
    rag.dataTable({
      rows,
      getRowKey: "id",
      columns: [
        { id: "id", header: "ID", cell: rag.cell.field("id") },
        { id: "name", header: "Name", cell: rag.cell.field("name") },
        { id: "status", header: "Status", cell: rag.cell.status("status") }
      ]
    })
  )
}
```

The table is still rendered by React. The JavaScript author only supplies serializable rows and serializable cell specifications. Avoid function-valued cell renderers; they cannot cross the JSON boundary and they are not part of Widget IR.

## Use recipes for common dashboards

Low-level helpers are useful when you need full control, but common RAG dashboard pages should use semantic recipes where possible. Recipes return ordinary Widget IR, so they remain JSON-compatible and render through the same React components.

```js
const rows = db.query("SELECT id, name, status FROM queries ORDER BY id")

return rag.page({
  id: "actions",
  title: "Actions",
  sections: [
    rag.recipes.metrics({ items: [
      { label: "Total", value: rows.length, status: "ready" },
      { label: "Running", value: rows.filter(row => row.status === "running").length, status: "running" }
    ]}),
    rag.recipes.actionToolbar({
      title: "Queue controls",
      actions: [
        { label: "Add query", variant: "primary", action: "add-query" },
        { label: "Reset", action: rag.action.server("reset-demo") }
      ]
    })
  ]
})
```

The current recipe set includes dashboard recipes (`metrics`, `actionToolbar`, `masterDetailTable`) and semantic context-viewer recipes (`contextDiagram`, `annotatedTranscript`, `courseStudio`, `courseSlide`, `handout`). Use recipes to make scripts read like page intent, then drop down to component helpers such as `panel`, `dataTable`, `transcriptWorkspacePanel`, `courseSlidePanel`, and `handoutDocumentShell` for custom sections.

For example, a jsverb can expose a context-window teaching page without building React-specific markup:

```js
const snapshot = {
  id: "ctx",
  title: "Context Window",
  limit: 16000,
  parts: [
    { id: "system", label: "System", kind: "system", tokens: 600 },
    { id: "retrieval", label: "Retrieved docs", kind: "retrieval", tokens: 7000 }
  ]
}

const slide = {
  id: "budget",
  number: "01",
  title: "Budget pressure",
  view: "budget",
  snapshotId: snapshot.id,
  notes: ["Retrieved documents dominate this example."]
}

return rag.page({
  id: "semantic",
  title: "Semantic context page",
  sections: [
    rag.recipes.contextDiagram({ snapshot, view: "budget" }),
    rag.recipes.courseStudio({
      sections: [{ id: "course", label: "Course", items: [{ id: "slides", label: "Slides" }] }],
      activeItemId: "slides",
      main: rag.recipes.courseSlide({ slide, snapshot, index: 0, total: 1 })
    })
  ]
})
```

## Build and run

Run the normal xgoja checks before building:

```bash
xgoja doctor -f xgoja.yaml
xgoja list-modules -f xgoja.yaml
xgoja build -f xgoja.yaml --xgoja-replace /path/to/go-go-goja
```

Then start the jsverb-backed site:

```bash
dist/rag-widget-xgoja-site serve sites demo --http-listen 127.0.0.1:18793
```

Open a frontend route handled by the SPA fallback:

```text
http://127.0.0.1:18793/pages/demo
```

The bundled widget-site example also exposes semantic pages for the expanded component library:

```text
http://127.0.0.1:18793/pages/semantic
http://127.0.0.1:18793/pages/transcripts
http://127.0.0.1:18793/pages/slides
http://127.0.0.1:18793/pages/handouts
http://127.0.0.1:18793/pages/course-examples
```

The browser asks the React app to render `demo`. The app calls `/api/widget/pages/demo`, receives Widget IR, and renders the page.

## What to remember

- Widget DSL constructors return JSON-compatible Widget IR. They do not return HTML strings or React elements.
- The page endpoint returns a page object with `schemaVersion`, `id`, `title`, and `root`.
- Strings and numbers used as children are converted to text nodes automatically.
- Tables use serializable `rag.cell.*` specifications instead of JavaScript render functions.
- The `widget.dsl` and `rag.dsl` modules expose the same helpers; choose one alias per script.
- The React app fetches `/api/widget/pages/{id}` and renders the returned `root` field.

## Troubleshooting

| Problem | Cause | Solution |
| --- | --- | --- |
| `Cannot find module "widget.dsl"` | The module was not selected in `xgoja.yaml`. | Add a `modules:` entry for package `rag-widget-site`, name `widget.dsl`, alias `widget.dsl`. |
| `xgoja build` tries to fetch the local provider from GitHub. | The provider package is local but the build spec has no `replace`. | Add `replace: ../../..` or another path to the RAG module root. |
| The browser route `/pages/demo` returns `404`. | The React app is served as static files without SPA fallback. | Use `app.spaFromAssetsModule("/", assets, "/app/public", { excludePrefixes: ["/api"] })`. |
| API routes return `index.html`. | The root SPA static handler is catching `/api/...`. | Add `/api` to `excludePrefixes`. |
| A status icon renders as `?`. | The status is not in the renderer's status vocabulary. | Use statuses such as `succeeded`, `running`, `warning`, `failed`, or `error`. |

## See Also

- `widget-dsl-js-api-reference`
- `tutorial-http-serve-jsverbs`
- `tutorial-static-assets-http-server`
- `buildspec-reference`
