---
Title: Widget DSL Event Timelines and Cutover Task Plan
Ticket: GOJA-DSL-PLAYBOOK
Status: active
Topics:
    - goja
    - dsl
    - fluent-builder
    - go
    - typescript
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/course-pages.js
      Note: Backend page ID dispatcher
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js
      Note: Real master-detail agenda editor consumer
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js
      Note: Live demo pages for v2 simplest table selectable table and master-detail editor
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/server.js
      Note: Backend Widget page
    - Path: packages/rag-evaluation-site/src/app/App.tsx
      Note: Frontend page loading and action dispatch timeline
    - Path: packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.widget.tsx
      Note: DataTable Widget IR adapter and row action context
    - Path: packages/rag-evaluation-site/src/components/organisms/FormPanel/FormPanel.tsx
      Note: Native form submit behavior for detail editors
    - Path: packages/rag-evaluation-site/src/widgets/actions.ts
      Note: Browser action dispatcher for navigate/server/confirm behavior
    - Path: pkg/widgetdsl/grammar.go
      Note: Current data.collection
    - Path: pkg/widgetdsl/module.go
      Note: Registers experimental data.v2.dsl native module
    - Path: pkg/widgetdsl/v2/spec/doc.go
      Note: Package documentation for the v2 spec layer
    - Path: pkg/widgetdsl/v2/spec/lower.go
      Note: Initial lowering from typed v2 specs to existing Widget IR maps for pages
    - Path: pkg/widgetdsl/v2/spec/lower_test.go
      Note: Positive and negative tests for simple table
    - Path: pkg/widgetdsl/v2/spec/types.go
      Note: Initial typed v2 Widget DSL intent model for pages
    - Path: pkg/widgetdsl/v2/spec/validate.go
      Note: Initial typed v2 validation rules for specs
    - Path: pkg/widgetdsl/v2_builders.go
      Note: Initial data.v2.dsl Goja fluent builders with hidden typed refs
    - Path: pkg/widgetdsl/v2_builders_test.go
      Note: Goja runtime tests for v2 builder simple/selectable tables and strict callback errors
ExternalSources: []
Summary: 'Operational companion to the rag-evaluation-system DSL overhaul guide: simple-to-rich event timelines, HTTP/frontend/backend execution traces, and a phase/task tracker for the hard-cutover v2 implementation.'
LastUpdated: 2026-07-05T18:50:00-04:00
WhatFor: Use when implementing or reviewing Widget DSL v2 behavior. It explains what authors write, what Widget IR is produced, what HTTP requests happen, what React code runs, and what backend code handles each interaction.
WhenToUse: Read beside design-doc 05 before implementing table, selection, master-detail editor, form submit, row action, or richer collection examples.
---









# Widget DSL Event Timelines and Cutover Task Plan

## Executive Summary

This document is the operational companion to `05-rag-evaluation-system-dsl-overhaul-design-and-implementation-guide.md`. The design guide says what the hard-cutover v2 DSL should become. This companion shows what using the examples should feel like at runtime: the authoring code, the compiled Widget IR, the browser events, the HTTP requests, the React code path, and the backend handler path.

The first three examples form the foundation for the richer DSL: a read-only table, a URL-selectable table, and a master-detail editor. They intentionally use the current system as evidence, because these flows already exist in `pkg/widgetdsl/grammar.go`, `WidgetRenderer`, `DataTable`, `FormPanel`, and `go-go-course`. The v2 implementation should preserve the good runtime behavior while replacing the public authoring substrate with typed builders/specs.

## Problem Statement

The design guide contains API examples, but implementers also need an event-level model. Without that model, it is easy to make the DSL look good in isolation while breaking browser navigation, bookmarkable selection, form submission, action payload hydration, or page refresh behavior.

This document answers:

- What does an author write for each example?
- What Widget IR should the DSL emit?
- What HTTP requests happen?
- What frontend code executes?
- What backend code executes?
- What state is stored in the URL, in the page JSON, in the browser context object, and in the backend data store?
- Which steps become implementation tasks for the hard-cutover v2?

## Current Runtime Evidence Map

Key current files:

- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/grammar.go`
  - `collectionVerb`: current `data.collection(rows, options)` compiler.
  - `collectionTable`: compiles schema fields to `DataTable` props and `onRowSelect` actions.
  - `collectionDetail`: compiles selected rows to a detail `FormPanel`.
  - `recordVerb`: compiles one record to either `MetadataGrid` or `FormPanel`.
  - `formPost`: creates native form-submit settings.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/app/App.tsx`
  - derives `pageId` and query string from browser location.
  - fetches `/api/widget/pages/:id`.
  - dispatches server actions to `/api/widget/actions/:name`.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/hooks/useWidgetPage.ts`
  - fetches Widget page JSON and stores loading/error/page state.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx`
  - renders Widget IR through the registry.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/actions.ts`
  - dispatches `navigate`, `server`, `download`, `copy`, and `event` actions.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.widget.tsx`
  - bridges Widget IR `DataTable` props to the React `DataTable` component.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.tsx`
  - renders clickable/selectable rows.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/cellRenderers.tsx`
  - renders cells and dispatches action-button cells.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/organisms/FormPanel/FormPanel.widget.tsx`
  - bridges Widget IR form props to React `FormPanel`.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/organisms/FormPanel/FormPanel.tsx`
  - renders the actual native `<form>`.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/server.js`
  - `GET /api/widget/pages/:id`: page JSON endpoint.
  - `POST /settings/agenda-item`: native form submit for agenda item saves.
  - `POST /api/widget/actions/:name`: JSON server-action endpoint.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/lib/course-pages.js`
  - dispatches page IDs to page builders.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js`
  - first real master-detail editor consumer.

## Proposed Solution

Use this document as both a behavioral specification and an implementation tracker. Every hard-cutover v2 task should preserve or intentionally change the timelines below. If a behavior changes, update this document and the demos in the same commit.

The implementation should proceed from simple to rich:

1. Read-only table.
2. URL-selectable table.
3. Master-detail editor with native form submit.
4. Master-detail editor with row server actions.
5. Multi-view/rich collection.
6. Domain fragments/marks over the typed builders.
7. Demo pages and regression fixtures for each stage.

## Example 1: Simplest Table

### Authoring code

Target v2 API:

```js
data.collection("sessions", sessions)
  .schema(sessionSchema)
  .table()
  .toIR();
```

Equivalent current v1 idea:

```js
data.collection(sessions, {
  schema: sessionSchema,
  arrange: "table",
});
```

### Expected Widget IR

```js
{
  kind: "component",
  type: "DataTable",
  props: {
    rows: sessions,
    getRowKey: "sessionId",   // schema key field, or id fallback
    columns: [
      { id: "sessionId", header: "ID", cell: { kind: "caption", field: "sessionId", tone: "muted" } },
      { id: "title", header: "Title", cell: { kind: "field", field: "title" } },
      { id: "turnCount", header: "Turns", cell: { kind: "number", field: "turnCount" } },
      { id: "status", header: "Status", cell: { kind: "status", field: "status" } }
    ],
    emptyMessage: "No sessions."
  }
}
```

There is no `onRowSelect` and no mutation action.

### Timeline: initial page load

```text
Browser URL
  /pages/sessions

Frontend shell
  readPageIdFromLocation() -> "sessions"
  readSearchFromLocation() -> ""
  useWidgetPage("/api/widget/pages/sessions")

HTTP
  GET /api/widget/pages/sessions

Backend
  server.js GET /api/widget/pages/:id
  -> buildWidgetPage("sessions", query, context)
  -> course-pages.js dispatches to sessions page builder
  -> page builder calls DSL
  -> DSL emits WidgetPage JSON

HTTP
  200 application/json { id, title, root }

Frontend render
  useWidgetPage stores page
  App.tsx renderPage(page)
  WidgetRenderer.renderComponentNode(root)
  registry.get("DataTable")
  DataTable.widget.tsx maps props to <DataTable>
  DataTable.tsx renders <table>
  cellRenderers.tsx renders each cell
```

### Interaction behavior

Clicking a row does nothing. The row is not clickable because `DataTable.widget.tsx` passes `onRowSelect={undefined}`.

No request is sent. No browser history entry is added. No backend code runs.

## Example 2: Selectable Table via URL

### Authoring code

Target v2 API:

```js
data.collection("sessions", sessions)
  .schema(sessionSchema)
  .select(data.selection.urlParam("selected", query.selected))
  .table(t => t.onRowSelect(data.action.navigateToSelection()))
  .toIR();
```

More explicit target form:

```js
data.collection("sessions", sessions, c => c
  .schema(sessionSchema)
  .select(s => s.urlParam("selected", query.selected))
  .arrange(a => a.table(t => t
    .rowSelect(ui.action.navigate("/pages/sessions?selected=${row.sessionId}")))))
```

Current v1 equivalent:

```js
data.collection(sessions, {
  schema: sessionSchema,
  select: data.urlParam("selected", query.selected),
});
```

In current `grammar.go`, providing `select` automatically adds a navigate action:

```go
props["onRowSelect"] = map[string]any{
  "kind": "navigate",
  "to": fmt.Sprintf("?%s=${row.%s}", selParam, keyField),
}
```

### Expected Widget IR

```js
{
  kind: "component",
  type: "DataTable",
  props: {
    rows: sessions,
    getRowKey: "sessionId",
    selectedKey: query.selected,
    onRowSelect: {
      kind: "navigate",
      to: "/pages/sessions?selected=${row.sessionId}"
    },
    columns: [/* derived from schema */]
  }
}
```

### Timeline: initial load with no selection

```text
Browser URL
  /pages/sessions

Frontend fetch
  GET /api/widget/pages/sessions

Backend query
  query.selected === undefined

DSL output
  selectedKey omitted
  onRowSelect = navigate("/pages/sessions?selected=${row.sessionId}")

Frontend render
  DataTable rows are clickable because onRowSelect exists
  No row is highlighted because selectedKey is absent
```

### Timeline: user clicks row `s2`

```text
User event
  click <tr> for row { sessionId: "s2", ... }

React DataTable.tsx
  onClick={() => onRowSelect(row)}

DataTable.widget.tsx
  ctx.dispatchAction(rowSelectAction, {
    row,
    rowKey: rowKey(row, props.getRowKey),
    componentType: "DataTable"
  })

Action context
  {
    row: { sessionId: "s2", title: "Debugging", ... },
    rowKey: "s2",
    componentType: "DataTable"
  }

App.tsx
  handleAction(action, context)
  action.kind !== "server"
  -> dispatchWidgetAction(action, context)

actions.ts
  interpolate("/pages/sessions?selected=${row.sessionId}", context)
  -> "/pages/sessions?selected=s2"
  window.history.pushState({}, "", target)
  window.dispatchEvent(new PopStateEvent("popstate"))

App.tsx
  popstate listener increments locationVersion
  readPageIdFromLocation() -> "sessions"
  readSearchFromLocation() -> "?selected=s2"
  useWidgetPage("/api/widget/pages/sessions?selected=s2")

HTTP
  GET /api/widget/pages/sessions?selected=s2

Backend
  query.selected === "s2"
  DSL emits selectedKey: "s2"

Frontend render
  DataTable.tsx compares selectedKey === key
  row s2 gets selected CSS class
```

### Important behavior

Selection is URL state. A selectable table click does not POST to the backend. It changes the URL, triggers the SPA page loader, and causes the backend to rebuild Widget IR from query params.

This is desirable for bookmarkable/admin-safe selection.

## Example 3: Master-Detail Editor

### Authoring code

Target v2 API from the design guide:

```js
data.collection("agenda", agenda, c => c
  .schema(agendaSchema)
  .edit(e => e
    .selectUrl("agenda", query.agenda)
    .submitPost("/settings/agenda-item")
    .create({ label: "New agenda item" })
    .actions(a => a
      .reorder(ui.action.server("admin-reorder-course-agenda"))
      .remove(ui.action.server("admin-delete-agenda-item", s => s
        .payload(p => p.path("id", "row.id"))
        .confirm(c => c.text("Delete agenda item “").path("row.title").text("”?"))))))
  .masterDetail())
  .toIR();
```

Current real consumer:

```js
dataDsl.collection(agenda, {
  schema: agendaSchema(),
  verb: "edit",
  arrange: "master-detail",
  select: dataDsl.urlParam("agenda", query && query.agenda),
  submit: dataDsl.formPost("/settings/agenda-item"),
  reorder: ui.action.server("admin-reorder-course-agenda"),
  remove: {
    kind: "server",
    name: "admin-delete-agenda-item",
    confirm: "Delete agenda item “${row.title}”? This cannot be undone."
  },
  create: { label: "New agenda item" },
  empty: "No agenda items yet. Create one to populate the course page.",
  status,
  statusMessage,
})
```

This is in `go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js`, inside `agendaSection()`.

### Expected Widget IR shape

Master-detail is a composed tree, not a single magic component. The current compiler emits roughly:

```js
{
  kind: "component",
  type: "Stack",
  props: { gap: "md" },
  children: [
    // optional create button row
    {
      kind: "component",
      type: "Inline",
      props: { gap: "sm", justify: "end" },
      children: [
        {
          kind: "component",
          type: "Button",
          props: { action: { kind: "navigate", to: "?agenda=__new" } },
          children: [{ kind: "text", text: "New agenda item" }]
        }
      ]
    },

    // summary table
    {
      kind: "component",
      type: "DataTable",
      props: {
        rows: agenda,
        getRowKey: "id",
        selectedKey: query.agenda,
        onRowSelect: { kind: "navigate", to: "?agenda=${row.id}" },
        columns: [
          /* schema columns, excluding prose/media */,
          { id: "moveUp", cell: { kind: "actionButton", label: "↑", action: { kind: "server", name: "admin-reorder-course-agenda", payload: { direction: "up" } } } },
          { id: "moveDown", cell: { kind: "actionButton", label: "↓", action: { kind: "server", name: "admin-reorder-course-agenda", payload: { direction: "down" } } } },
          { id: "delete", cell: { kind: "actionButton", label: "Delete", action: { kind: "server", name: "admin-delete-agenda-item", confirm: "Delete agenda item “${row.title}”? This cannot be undone." } } }
        ]
      }
    },

    // detail area
    {
      kind: "component",
      type: "Stack",
      props: { gap: "sm" },
      children: [
        {
          kind: "component",
          type: "FormPanel",
          props: {
            title: "Edit: Cadrage + concepts fondamentaux",
            formAction: "/settings/agenda-item",
            method: "post",
            status,
            statusMessage
          },
          children: [
            { kind: "component", type: "FieldGrid", children: [/* short FormRows */] },
            { kind: "component", type: "FormRow", props: { control: { type: "TextareaInput", props: { name: "description", defaultValue: "..." } } } }
          ]
        },
        {
          kind: "component",
          type: "Inline",
          children: [
            { kind: "component", type: "Button", props: { action: { kind: "navigate", to: "?agenda=" } }, children: [{ kind: "text", text: "Close" }] }
          ]
        }
      ]
    }
  ]
}
```

### Timeline A: initial load with no selected row

```text
Browser URL
  /pages/admin-course-cms

Frontend fetch
  GET /api/widget/pages/admin-course-cms

Backend route
  server.js GET /api/widget/pages/:id
  -> buildWidgetPage("admin-course-cms", {}, context)
  -> course-pages.js
  -> adminCmsPage.buildAdminCourseCmsPage(query, context)

Backend page builder
  verifies admin display name through courseMaterial.isAdminUser(user)
  loads metadata/content/material
  calls agendaSection(content.agenda || [], query)
  query.agenda is undefined

DSL current compiler
  dataDsl.urlParam("agenda", undefined) -> { param: "agenda", value: "" }
  collectionVerb(...)
  keyField = schema key field "id"
  selParam = "agenda"
  selValue = ""
  collectionTable(...) emits DataTable with onRowSelect navigate("?agenda=${row.id}")
  collectionCreateButton(...) emits New button navigate("?agenda=__new")
  collectionDetail(...) sees selValue == "" and emits Caption "Select a row to open it here."

HTTP response
  200 WidgetPage JSON

Frontend render
  WidgetRenderer renders CourseStudioShell -> page main -> SectionBlock -> Stack
  DataTable rows are clickable
  No row is highlighted
  Detail area says: Select a row to open it here.
```

### Timeline B: user selects an existing row

```text
User event
  Click agenda table row with id "agenda-foundations"

DataTable.tsx
  <tr onClick> calls onRowSelect(row)

DataTable.widget.tsx
  dispatchAction({ kind: "navigate", to: "?agenda=${row.id}" }, {
    row,
    rowKey: "agenda-foundations",
    componentType: "DataTable"
  })

App.tsx / actions.ts
  navigate action interpolates target -> "?agenda=agenda-foundations"
  history.pushState({}, "", "?agenda=agenda-foundations")
  dispatches popstate

Frontend fetch
  GET /api/widget/pages/admin-course-cms?agenda=agenda-foundations

Backend page builder
  query.agenda === "agenda-foundations"
  agendaSection passes dataDsl.urlParam("agenda", query.agenda)

DSL compiler
  selectedKey = "agenda-foundations"
  collectionDetail searches rows for row.id == "agenda-foundations"
  found row becomes values for recordVerb
  recordVerb emits FormPanel with defaultValue for each field
  key field id is readOnly by default unless explicitly editable

Frontend render
  Row "agenda-foundations" is highlighted
  Detail form title becomes "Edit: <primary title>"
  Inputs are native form controls with names id/number/duration/title/description
```

### Timeline C: user clicks New agenda item

```text
User event
  Click "New agenda item" button

Button.widget.tsx
  ctx.bindAction(props.action, { componentType: "Button" })

actions.ts
  action.kind === "navigate"
  target "?agenda=__new"
  pushState + popstate

Frontend fetch
  GET /api/widget/pages/admin-course-cms?agenda=__new

Backend / DSL
  collectionDetail sees selValue == "__new"
  values = {}
  title = "New item"
  recordVerb emits editable FormPanel with empty defaults

Frontend render
  No existing row should be selected because collectionTable suppresses selectedKey for "__new"
  Detail form shows empty fields
```

### Timeline D: user submits the detail form

```text
User event
  Edit fields and click Save

React FormPanel.tsx
  renders a native <form action="/settings/agenda-item" method="post">
  no custom React submit handler is involved

Browser HTTP
  POST /settings/agenda-item
  Content-Type: application/x-www-form-urlencoded or multipart/form-data depending platform/form settings
  Body fields:
    id=agenda-foundations
    number=14h30
    duration=15+min
    title=Cadrage+updated
    description=...

Backend route
  server.js POST /settings/agenda-item
  requireAdminProfile(ctx, res)
  body = ctx.body || {}
  courseMetadata.upsertAgendaItem(currentCourseContent().agenda, body)
  res.redirect(`/pages/admin-course-cms?agenda=${saved.item.id}&status=agenda-item-saved`)

Browser navigation
  Receives 302 redirect
  Loads /pages/admin-course-cms?agenda=<saved-id>&status=agenda-item-saved
  SPA shell remains the page app; frontend fetches new Widget IR

Frontend fetch after redirect
  GET /api/widget/pages/admin-course-cms?agenda=<saved-id>&status=agenda-item-saved

Backend page builder
  formStatus(query, "agenda-item-saved") -> success
  statusMessage -> "Agenda item saved."
  recordVerb passes status/statusMessage to FormPanel

Frontend render
  selected row highlighted
  FormPanel status area announces Saved / Agenda item saved
```

### Timeline E: user clicks reorder arrow

```text
User event
  Click ↑ or ↓ action button cell

cellRenderers.tsx
  event.stopPropagation() prevents row selection
  dispatchAction(spec.action, {
    row,
    rowKey: rowKey(row, "file"), // current generic action-cell rowKey quirk
    componentType: "DataTableCell"
  })

ActionSpec
  {
    kind: "server",
    name: "admin-reorder-course-agenda",
    payload: { direction: "up" }
  }

App.tsx
  handleAction sees kind === "server"
  fetch(`${apiBase}/actions/${name}`, POST JSON)

HTTP
  POST /api/widget/actions/admin-reorder-course-agenda
  Content-Type: application/json
  Body:
  {
    "payload": { "direction": "up" },
    "context": {
      "row": { "id": "agenda-foundations", ... },
      "rowKey": "",
      "componentType": "DataTableCell"
    }
  }

Backend route
  server.js POST /api/widget/actions/:name
  actionName === "admin-reorder-course-agenda"
  requireAdminProfile
  payload = ctx.body.payload
  row = ctx.body.context.row
  courseMetadata.reorderAgendaItemById(currentCourseContent().agenda, payload.id || row.id, payload.direction)
  res.json({ ok: true, refresh: true, toast, data: { agenda } })

Frontend response handling
  App.tsx sees result.toast -> window.dispatchEvent(new CustomEvent("widget:toast", { detail: result }))
  result.refresh -> refresh()
  useWidgetPage refetches current URL

HTTP refresh
  GET /api/widget/pages/admin-course-cms?agenda=<current>&status=...

Frontend render
  Table order updates
  Current selected row remains selected if the URL still points to it
```

Note: `cellRenderers.tsx` currently computes action-cell `rowKey` with `rowKey(row, "file")`. That is acceptable for media/file tables but wrong for generic collections. V2 should pass the table's `getRowKey` into action cell rendering, or include a typed row context from the table adapter.

### Timeline F: user clicks Delete

```text
User event
  Click Delete action button cell

cellRenderers.tsx
  stopPropagation prevents selection
  dispatches server action with row context

actions.ts / App.tsx
  current app-level handler handles server actions directly; confirm behavior must remain centralized
  Current generic dispatchWidgetAction has confirm handling, but App.tsx bypasses it for server actions.

Expected v2 behavior
  confirm template is evaluated before POST
  if user cancels: no HTTP request
  if user confirms: POST server action

HTTP if confirmed
  POST /api/widget/actions/admin-delete-agenda-item
  Body:
  {
    "payload": {},
    "context": {
      "row": { "id": "agenda-foundations", "title": "..." },
      "componentType": "DataTableCell"
    }
  }

Backend route
  actionName === "admin-delete-agenda-item"
  id = payload.id || row.id
  courseMetadata.deleteAgendaItem(currentCourseContent().agenda, id)
  res.json({ ok: true, refresh: true, toast: `Deleted agenda item ${id}`, data: { agenda } })

Frontend refresh
  result.refresh -> useWidgetPage refreshes current URL

Potential edge
  If URL still has ?agenda=<deleted-id>, backend collectionDetail returns Caption "No row matches <id>."
  V2 can improve this by returning a patch/navigate result or by clearing selection after delete.
```

## Requests Summary for Master-Detail

Initial page:

```http
GET /pages/admin-course-cms
GET /api/widget/pages/admin-course-cms
```

Select row:

```http
# browser history change, then SPA JSON reload
GET /api/widget/pages/admin-course-cms?agenda=agenda-foundations
```

Create new:

```http
# browser history change, then SPA JSON reload
GET /api/widget/pages/admin-course-cms?agenda=__new
```

Save form:

```http
POST /settings/agenda-item
302 Location: /pages/admin-course-cms?agenda=<saved-id>&status=agenda-item-saved
GET /pages/admin-course-cms?agenda=<saved-id>&status=agenda-item-saved
GET /api/widget/pages/admin-course-cms?agenda=<saved-id>&status=agenda-item-saved
```

Reorder:

```http
POST /api/widget/actions/admin-reorder-course-agenda
GET /api/widget/pages/admin-course-cms?agenda=<current>
```

Delete:

```http
POST /api/widget/actions/admin-delete-agenda-item
GET /api/widget/pages/admin-course-cms?agenda=<possibly-deleted-current>
```

## Demo and Example Inventory

This inventory captures the current baseline before v2 demo pages are added. Treat these as evidence and regression references, not as final v2 examples.

### Existing live pages

- `/pages/sessions` — selectable `DataTable` built directly with `dataDsl.dataTable` in `go-go-course/cmd/go-go-course/lib/pages/sessions.js`. This is the closest current live page for the selectable-table timeline, but it is not authored through `data.collection`.
- `/pages/admin-course-cms` — agenda master-detail editor built through `dataDsl.collection(..., { verb: "edit", arrange: "master-detail" })` in `go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js`. This is the canonical current runtime fixture for the master-detail timeline.
- `adminMaterialTable(...)` in `go-go-course/cmd/go-go-course/lib/pages/admin-common.js` — direct `dataDsl.dataTable` helper used for slide/handout file tables with action/link cells.

### Existing Storybook / component examples

- `packages/rag-evaluation-site/src/widgets/WidgetRenderer.domain-registry.stories.tsx` — contains a WidgetRenderer-level `DataTable` example and validates registry wiring.
- `packages/rag-evaluation-site/src/widgets/WidgetRenderer.forms.stories.tsx` — contains `FormPanel` examples and form layout evidence.
- `packages/rag-evaluation-site/src/components/molecules/DataTable/*` stories/components — component-level table rendering behavior.
- `packages/rag-evaluation-site/src/components/organisms/FormPanel/*` stories/components — component-level form behavior.

### Missing demos to build

- A dedicated simplest-table live page authored through the new v2 `data.collection(...).table().toIR()` API.
- A dedicated selectable-table live page authored through the new v2 selection API, not direct `dataDsl.dataTable`.
- A dedicated master-detail demo page that is safe for repeated local testing and does not mutate course metadata unless explicitly intended.
- A row/server-action demo that covers confirm, cancel, POST, refresh, and toast behavior.

### Deprecated-example policy

When a v2 demo exists, any v1 option-bag example that teaches the same concept must be either removed from public docs or moved under a clearly named historical/legacy section. Small models should not see both forms as equally valid.

### Baseline validation commands

Recorded before implementing v2 code changes:

```bash
cd rag-evaluation-system
go test ./pkg/widgetdsl -count=1
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build
docmgr doctor --ticket GOJA-DSL-PLAYBOOK --stale-after 30
```

Observed result on 2026-07-05:

- `go test ./pkg/widgetdsl -count=1` passed.
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `pnpm --dir packages/rag-evaluation-site build` passed; build artifacts were generated under `dist` but did not leave tracked git changes.
- `docmgr doctor --ticket GOJA-DSL-PLAYBOOK --stale-after 30` passed.

## Design Decisions

### Decision 1: URL selection remains data, not component state

- **Context:** Tables and master-detail editors need bookmarkable selections and browser back/forward behavior.
- **Decision:** Selection state belongs in the URL for these examples.
- **Rationale:** The frontend can remain stateless with respect to selected rows; the backend rebuilds detail panes from query params.
- **Consequence:** Row selection is a navigate action, not a server action.
- **Status:** accepted for v2 unless a richer local-only interaction explicitly opts out.

### Decision 2: Form submit remains native for simple record saves

- **Context:** Master-detail editors currently use `FormPanel` with `action` and `method`, handled by normal browser form semantics.
- **Decision:** Keep native form submit as a blessed path for simple CRUD saves.
- **Rationale:** It is understandable, accessible, and works across reloads/redirects.
- **Consequence:** The backend save handler can redirect to the canonical selected URL and status state.
- **Status:** proposed for v2.

### Decision 3: Row action buttons use server actions with typed row context

- **Context:** Reorder/delete are not form submits; they need row context and JSON responses.
- **Decision:** Use `ActionSpec`/Action IR v2 data and typed context contracts.
- **Rationale:** The browser should POST data, not serialized closures.
- **Consequence:** V2 must fix generic row-key context and centralize confirm behavior.
- **Status:** proposed for v2.

### Decision 4: Demos are part of the implementation contract

- **Context:** The DSL is intended for humans and small models, so examples are not optional.
- **Decision:** Every implemented phase must add/update/remove demos and examples in the same phase.
- **Rationale:** Deprecated examples are more dangerous than missing examples because agents imitate them.
- **Consequence:** CI should eventually check that examples compile/run against the current public API.
- **Status:** accepted.

## Alternatives Considered

### Alternative: local React selection state

Rejected for these examples. It would make selection fast, but it would not survive reloads, browser back/forward, or server-rendered detail panes without duplicating selection state.

### Alternative: all mutations as server actions

Rejected as the only path. Server actions are right for row buttons and asynchronous actions, but native forms are simpler and more robust for basic record saves.

### Alternative: keep v1 option-bag examples in demos

Rejected for v2 demos. V1 examples may remain in a historical/legacy folder while porting, but the public demos for v2 must use only the hard-cutover API.

## Implementation Plan and Task Tracker

This section is intentionally detailed so progress can be tracked precisely in `tasks.md`, diary entries, and commits.

### Phase 0: Behavioral fixtures and demo scaffolding

Goal: preserve current behavior as evidence before replacing the public DSL.

Tasks:

- P0.1 Create this companion document with simple/selectable/master-detail timelines.
- P0.2 Add docmgr tasks for the whole cutover.
- P0.3 Create demo-page inventory: current pages, Storybook stories, and missing examples.
- P0.4 Add or identify demo pages for simplest table, selectable table, and master-detail editor.
- P0.5 Add request/interaction notes to each demo example.
- P0.6 Run current tests and record baseline failures.
- P0.7 Commit documentation and baseline demo scaffolding.

### Phase 1: Typed v2 spec model

Goal: create the typed Go model underneath the v2 DSL before exposing new JavaScript APIs.

Tasks:

- P1.1 Add `pkg/widgetdsl/v2/spec` package.
- P1.2 Define `PageSpec`, `NodeSpec`, `CollectionSpec`, `SchemaSpec`, `FieldSpec`, `ActionSpec`, `TemplateSpec`, `ValidationIssue`.
- P1.3 Implement validation for schema keys, field names, arrangements, section levels, action context paths, and template paths.
- P1.4 Implement conversion from typed specs to existing Widget IR nodes.
- P1.5 Add unit tests for valid simplest table, selectable table, and master-detail specs.
- P1.6 Add negative unit tests for v1 failure modes: typo'd arrangement, wrong marker kind, invalid section level, unknown field option, bad action path.
- P1.7 Commit typed spec model and tests.

### Phase 2: v2 builder substrate

Goal: expose typed/fluent builders in Goja without v1 compatibility shims.

Tasks:

- P2.1 Add hidden typed refs for schema, field, collection, action, template, arrangement, and page builders.
- P2.2 Implement strict callback handling: callback absent means defaults; callback present and non-function is an error.
- P2.3 Implement `.validate()` and `.toIR()` terminals.
- P2.4 Implement `data.schema(name).field(...).build()`.
- P2.5 Implement `data.collection(name, rows).schema(...).table().toIR()`.
- P2.6 Implement `select(s => s.urlParam(...))` and `table(t => t.rowSelect(...))`.
- P2.7 Implement `edit(...).masterDetail()` with native form submit.
- P2.8 Add Goja runtime tests for simple/selectable/master-detail authoring code.
- P2.9 Commit v2 builder substrate.

### Phase 3: Action IR v2 and context hydration

Goal: make all browser-visible behavior serializable and typed.

Tasks:

- P3.1 Define Action IR v2 payload/template structs and TypeScript types.
- P3.2 Implement template interpolation for path/literal/text parts.
- P3.3 Centralize confirm behavior so app-level server actions cannot bypass confirmation.
- P3.4 Fix DataTable action-cell context to use the table's `getRowKey`, not hard-coded `file`.
- P3.5 Add typed DataTable row and DataTable cell context contracts.
- P3.6 Update backend action handler examples to consume `payload` and `context` consistently.
- P3.7 Add tests for navigate, reorder, delete, cancel-confirm, and refresh.
- P3.8 Commit Action IR v2 foundation.

### Phase 4: Demo site examples

Goal: make the DSL teachable and manually testable.

Tasks:

- P4.1 Add `/pages/dsl-examples-table` demo.
- P4.2 Add `/pages/dsl-examples-selectable-table` demo.
- P4.3 Add `/pages/dsl-examples-master-detail` demo.
- P4.4 Add `/pages/dsl-examples-actions` demo for reorder/delete/server action refresh.
- P4.5 Add navigation entry or discoverability page for DSL examples.
- P4.6 Add demo README explaining expected HTTP requests and interactions.
- P4.7 Remove or clearly quarantine deprecated v1 examples.
- P4.8 Commit demos.

### Phase 5: TypeScript declarations and parity

Goal: make the API discoverable and impossible to miscall silently.

Tasks:

- P5.1 Generate precise `.d.ts` for v2 modules.
- P5.2 Add runtime export parity tests.
- P5.3 Add `tsc` positive fixtures for each example.
- P5.4 Add `@ts-expect-error` fixtures for removed v1 APIs and wrong handle types.
- P5.5 Remove `Props = Record<string, any>` from high-level v2 declarations.
- P5.6 Commit DTS/parity work.

### Phase 6: Port real pages by rewrite

Goal: prove the cutover on real consumers.

Tasks:

- P6.1 Rewrite the admin agenda editor to v2 master-detail API.
- P6.2 Rewrite media library usage as schema + marks + fragments.
- P6.3 Rewrite session browse table to v2 selectable table API.
- P6.4 Rewrite any remaining package examples using v1 collection/schema helpers.
- P6.5 Delete old public v1 exports from chosen v2 modules.
- P6.6 Commit real page port.

### Phase 7: CI and cleanup

Goal: make the hard cutover enforceable.

Tasks:

- P7.1 Add lint/test to reject v1 option-bag public APIs in v2 modules.
- P7.2 Add example-runner smoke test for demo pages.
- P7.3 Add docs/tests for `widget.unsafe` usage and production lint warning.
- P7.4 Run full Go and frontend test/build suite.
- P7.5 Update design docs, companion doc, and diary with final behavior.
- P7.6 Upload final bundle to reMarkable if requested.
- P7.7 Commit cleanup and final docs.

## Open Questions

1. Should `widget.unsafe` exist in production builds, or only in development/examples?
2. Should native form submit remain the primary record-save path, or should v2 also provide a fetch-based form action for no-navigation saves?
3. Should delete actions auto-clear URL selection after deleting the selected row?
4. Should demo pages live in `go-go-course`, Storybook, or both?
5. Should the v2 package live beside the existing `pkg/widgetdsl` initially, or replace it in place after baseline fixtures are captured?

## References

- `design-doc/05-rag-evaluation-system-dsl-overhaul-design-and-implementation-guide.md`
- `design-doc/03-widget-dsl-design-assessment-and-improvement-report.md`
- `design-doc/02-self-assessment-of-the-widgetdsl-grammar-what-pattern-c-actually-costs-and-what-the-playbook-should-add.md`
- `reference/02-investigation-diary.md`
