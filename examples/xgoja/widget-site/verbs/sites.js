__package__({ name: "sites", short: "WidgetRenderer xgoja sites" })

__verb__("demo", {
  name: "demo",
  output: "text",
  short: "Serve a tiny Widget IR site backed by express, db, assets, and widget.dsl",
  tags: ["http", "widget", "db"]
})
function demo() {
  const express = require("express")
  const assets = require("fs:assets")
  const db = require("db")
  const rag = require("widget.dsl")

  db.configure("sqlite3", ":memory:")
  db.exec("CREATE TABLE queries (id INTEGER PRIMARY KEY, name TEXT, status TEXT)")
  db.exec("INSERT INTO queries (name, status) VALUES (?, ?)", "Fast Growing Trees", "succeeded")
  db.exec("INSERT INTO queries (name, status) VALUES (?, ?)", "Arborvitae Spacing", "running")

  const app = express.app()
  app.spaFromAssetsModule("/", assets, "/app/public", { excludePrefixes: ["/api", "/healthz", "/favicon.ico"] })
  app.get("/favicon.ico", (_req, res) => res.status(204).end())
  app.get("/healthz", (_req, res) => res.json({ ok: true, site: "rag-widget-xgoja-site" }))
  app.get("/api/widget/schema", (_req, res) => res.json({ schemaVersion: "0.1.0", components: ["Panel", "StatusText", "DataTable"] }))

  function widgetPage(id) {
    const rows = db.query("SELECT id, name, status FROM queries ORDER BY id")
    return {
      schemaVersion: "0.1.0",
      id,
      title: "xgoja widget demo",
      root: rag.panel({ title: "xgoja widget demo" },
        rag.statusText({ status: "succeeded", icon: true }, "Rows: " + rows.length),
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
  }

  app.get("/api/widget/pages/index", (_req, res) => res.json(widgetPage("index")))
  app.get("/api/widget/pages/demo", (_req, res) => res.json(widgetPage("demo")))
}
