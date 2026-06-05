__package__({ name: "sites", short: "WidgetRenderer xgoja sites" })

__verb__("demo", {
  name: "demo",
  output: "text",
  short: "Serve an interactive Widget IR site backed by express, db, assets, and widget.dsl",
  tags: ["http", "widget", "db", "actions"]
})
function demo() {
  const express = require("express")
  const assets = require("fs:assets")
  const db = require("db")
  const rag = require("widget.dsl")

  db.configure("sqlite3", ":memory:")
  db.exec("CREATE TABLE queries (id INTEGER PRIMARY KEY, name TEXT, status TEXT, priority INTEGER, owner TEXT, notes TEXT)")
  db.exec("INSERT INTO queries (name, status, priority, owner, notes) VALUES (?, ?, ?, ?, ?)", "Fast Growing Trees", "succeeded", 3, "botany", "Seeded from the xgoja demo")
  db.exec("INSERT INTO queries (name, status, priority, owner, notes) VALUES (?, ?, ?, ?, ?)", "Arborvitae Spacing", "running", 2, "landscape", "Needs another pass")
  db.exec("INSERT INTO queries (name, status, priority, owner, notes) VALUES (?, ?, ?, ?, ?)", "Broken Import Example", "failed", 1, "ingest", "Retry after source cleanup")

  const appState = { selectedId: 1, audit: ["Demo initialized"] }
  const app = express.app()
  app.spaFromAssetsModule("/", assets, "/app/public", { excludePrefixes: ["/api", "/healthz", "/favicon.ico"] })
  app.get("/favicon.ico", (_req, res) => res.status(204).end())
  app.get("/healthz", (_req, res) => res.json({ ok: true, site: "rag-widget-xgoja-site" }))
  app.get("/api/widget/schema", (_req, res) => res.json({ schemaVersion: "0.1.0", components: ["Panel", "StatusText", "DataTable", "Button", "MetadataGrid"] }))

  function allRows() {
    return db.query("SELECT id, name, status, priority, owner, notes FROM queries ORDER BY priority DESC, id ASC")
  }

  function selectedRow() {
    const rows = db.query("SELECT id, name, status, priority, owner, notes FROM queries WHERE id = ?", appState.selectedId)
    return rows.length > 0 ? rows[0] : null
  }

  function counts() {
    const rows = db.query("SELECT status, COUNT(*) AS count FROM queries GROUP BY status ORDER BY status")
    const out = { total: 0, succeeded: 0, running: 0, failed: 0, pending: 0 }
    rows.forEach(row => {
      const count = Number(row.count || 0)
      out.total += count
      out[row.status] = count
    })
    return out
  }

  function nextStatus(status) {
    if (status === "pending") return "running"
    if (status === "running") return "succeeded"
    if (status === "succeeded") return "failed"
    return "pending"
  }

  function pageSummary(id) {
    const c = counts()
    return rag.recipes.metrics({ items: [
      { label: "Total queries", value: c.total, status: "ready" },
      { label: "Succeeded", value: c.succeeded || 0, status: "succeeded" },
      { label: "Running", value: c.running || 0, status: "running" },
      { label: "Failed", value: c.failed || 0, status: "failed" }
    ]})
  }

  function toolbar() {
    return rag.recipes.actionToolbar({
      title: "Queue controls",
      caption: "Actions mutate in-memory SQLite state and ask the React app to refresh.",
      actions: [
        { label: "Add query", variant: "primary", action: "add-query", payload: { owner: "research" } },
        { label: "Retry failed", action: "bulk-retry-failed" },
        { label: "Reset demo", action: "reset-demo" }
      ]
    })
  }

  function queryColumns() {
    return [
      { id: "priority", header: "Priority", align: "right", cell: rag.cell.number("priority") },
      { id: "name", header: "Query", cell: rag.cell.field("name") },
      { id: "owner", header: "Owner", cell: rag.cell.caption("owner", { tone: "muted" }) },
      { id: "status", header: "Status", cell: rag.cell.status("status", { icon: true }) },
      { id: "notes", header: "Notes", cell: rag.cell.caption("notes", { tone: "muted" }) }
    ]
  }

  function selectedPanel(row) {
    if (!row) {
      return rag.panel({ title: "Selected query" }, rag.caption({ tone: "muted" }, "Select a row to inspect it."))
    }
    return rag.panel({ title: "Selected query" },
      rag.metadataGrid({ density: "compact", items: [
        { key: "ID", value: String(row.id), copyValue: String(row.id) },
        { key: "Name", value: row.name, copyValue: row.name },
        { key: "Owner", value: row.owner },
        { key: "Priority", value: String(row.priority) },
        { key: "Status", value: rag.statusText({ status: row.status, icon: true }, row.status) },
        { key: "Notes", value: row.notes }
      ]}),
      rag.inline({ gap: "sm", wrap: true },
        rag.button({ variant: "primary", action: { kind: "server", name: "cycle-status", payload: { id: row.id } } }, "Cycle status"),
        rag.button({ variant: "secondary", action: { kind: "server", name: "bump-priority", payload: { id: row.id, delta: 1 } } }, "Priority +1"),
        rag.button({ variant: "secondary", action: { kind: "server", name: "bump-priority", payload: { id: row.id, delta: -1 } } }, "Priority -1"),
        rag.button({ variant: "danger", action: { kind: "server", name: "archive-query", payload: { id: row.id } } }, "Archive")
      )
    )
  }

  function auditPanel() {
    return rag.panel({ title: "Action audit trail", density: "condensed" },
      rag.stack({ gap: "xs" },
        appState.audit.slice(-6).reverse().map(entry => rag.caption({ tone: "muted" }, entry))
      )
    )
  }

  function widgetPage(id) {
    const rows = allRows()
    const selected = selectedRow()
    return rag.page({
      schemaVersion: "0.1.0",
      id,
      title: "xgoja widget actions demo",
      sections: [
        rag.panel({ title: "xgoja widget actions demo" },
          rag.statusText({ status: "succeeded", icon: true }, "Rows: " + rows.length),
          rag.caption({ tone: "muted" }, "This page is authored by a jsverb. Buttons and table row selection call /api/widget/actions/{name} and refresh the React app.")
        ),
        pageSummary(id),
        toolbar(),
        rag.recipes.masterDetailTable({
          title: "Query queue",
          rows,
          columns: queryColumns(),
          selectedKey: appState.selectedId,
          onRowSelect: "select-query",
          detail: () => selectedPanel(selected)
        }),
        auditPanel()
      ]
    })
  }

  function payload(req) {
    return (req.body && req.body.payload) || {}
  }

  function actionContext(req) {
    return (req.body && req.body.context) || {}
  }

  function remember(message) {
    appState.audit.push(new Date().toISOString().slice(11, 19) + " " + message)
  }

  function actionResult(message) {
    remember(message)
    return { ok: true, refresh: true, toast: message, data: { message } }
  }

  app.get("/api/widget/pages/index", (_req, res) => res.json(widgetPage("index")))
  app.get("/api/widget/pages/demo", (_req, res) => res.json(widgetPage("demo")))
  app.get("/api/widget/pages/actions", (_req, res) => res.json(widgetPage("actions")))

  app.post("/api/widget/actions/select-query", (req, res) => {
    const ctx = actionContext(req)
    const row = ctx.row || {}
    const nextId = Number(row.id || ctx.rowKey || payload(req).id || appState.selectedId)
    if (nextId > 0) appState.selectedId = nextId
    res.json(actionResult("Selected query #" + appState.selectedId))
  })

  app.post("/api/widget/actions/cycle-status", (req, res) => {
    const id = Number(payload(req).id || appState.selectedId)
    const rows = db.query("SELECT status FROM queries WHERE id = ?", id)
    if (rows.length === 0) return res.status(404).json({ ok: false, error: "query not found" })
    const status = nextStatus(rows[0].status)
    db.exec("UPDATE queries SET status = ? WHERE id = ?", status, id)
    appState.selectedId = id
    res.json(actionResult("Query #" + id + " status -> " + status))
  })

  app.post("/api/widget/actions/bump-priority", (req, res) => {
    const data = payload(req)
    const id = Number(data.id || appState.selectedId)
    const delta = Number(data.delta || 1)
    db.exec("UPDATE queries SET priority = MAX(0, priority + ?) WHERE id = ?", delta, id)
    appState.selectedId = id
    res.json(actionResult("Query #" + id + " priority changed by " + delta))
  })

  app.post("/api/widget/actions/archive-query", (req, res) => {
    const id = Number(payload(req).id || appState.selectedId)
    db.exec("DELETE FROM queries WHERE id = ?", id)
    const rows = allRows()
    appState.selectedId = rows.length > 0 ? Number(rows[0].id) : 0
    res.json(actionResult("Archived query #" + id))
  })

  app.post("/api/widget/actions/add-query", (req, res) => {
    const owner = String(payload(req).owner || "research")
    const next = db.query("SELECT COALESCE(MAX(id), 0) + 1 AS id FROM queries")[0].id
    db.exec("INSERT INTO queries (name, status, priority, owner, notes) VALUES (?, ?, ?, ?, ?)", "Follow-up Query " + next, "pending", 1, owner, "Created by a Widget IR action")
    appState.selectedId = Number(next)
    res.json(actionResult("Added query #" + next))
  })

  app.post("/api/widget/actions/bulk-retry-failed", (_req, res) => {
    db.exec("UPDATE queries SET status = 'running', notes = 'Bulk retry requested from Widget IR action' WHERE status = 'failed'")
    res.json(actionResult("Retried failed queries"))
  })

  app.post("/api/widget/actions/reset-demo", (_req, res) => {
    db.exec("DELETE FROM queries")
    db.exec("INSERT INTO queries (id, name, status, priority, owner, notes) VALUES (?, ?, ?, ?, ?, ?)", 1, "Fast Growing Trees", "succeeded", 3, "botany", "Seeded from the xgoja demo")
    db.exec("INSERT INTO queries (id, name, status, priority, owner, notes) VALUES (?, ?, ?, ?, ?, ?)", 2, "Arborvitae Spacing", "running", 2, "landscape", "Needs another pass")
    db.exec("INSERT INTO queries (id, name, status, priority, owner, notes) VALUES (?, ?, ?, ?, ?, ?)", 3, "Broken Import Example", "failed", 1, "ingest", "Retry after source cleanup")
    appState.selectedId = 1
    appState.audit = ["Demo reset"]
    res.json({ ok: true, refresh: true, toast: "Demo reset", data: { message: "Demo reset" } })
  })
}
