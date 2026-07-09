---
Title: Widget DSL v3 Integration and Migration Guide
Ticket: RAGEVAL-SCHEDULE-WIDGETS
Status: active
Topics:
    - ui-dsl
    - widget-ir
    - design-system
    - react
    - frontend-architecture
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://examples/xgoja-widgetdsl-v3/README.md
      Note: Build and validation instructions for the v3 preview host
    - Path: repo://examples/xgoja-widgetdsl-v3/jsverbs/server.js
      Note: Query-aware preview server pattern for returning Widget IR pages from JavaScript.
    - Path: repo://examples/xgoja-widgetdsl-v3/xgoja.yaml
      Note: |-
        Reference xgoja build spec that selects the parallel widget.dsl module and serves the v3 example gallery.
        Reference widget.dsl-only rag-widget-site module selection
    - Path: repo://pkg/xgoja/providers/widgetsite/doc/01-widget-dsl-getting-started.md
      Note: Embedded provider getting-started documentation updated with widget.dsl v3 guidance.
    - Path: repo://pkg/xgoja/providers/widgetsite/doc/03-widget-dsl-spa-bundling.md
      Note: |-
        Embedded provider SPA bundling documentation updated with widget.dsl v3 module selection.
        Embedded SPA bundling docs with v3 module-selection guidance
    - Path: repo://ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/scripts/02-report-legacy-widget-dsl-usage.py
      Note: |-
        Migration helper that reports legacy split-module imports and raw component escape hatches.
        Checker referenced by the migration workflow
ExternalSources: []
Summary: Integration and migration guide for adopting the parallel widget.dsl v3 module while keeping legacy split modules available for existing hosts.
LastUpdated: 2026-07-08T20:20:00-04:00
WhatFor: Use this when wiring a new xgoja host, migrating first-party pages from split modules, or deciding whether a script is ready for widget.dsl v3.
WhenToUse: Read before changing xgoja runtime module selections or porting ui.dsl/data.dsl/context_window.dsl/course.dsl pages to widget.dsl.
---


# Widget DSL v3 Integration and Migration Guide

## Goal

`widget.dsl` v3 is now available as a parallel module through the `rag-widget-site` xgoja provider. New hosts should prefer `widget.dsl`; existing hosts can keep the legacy split modules (`ui.dsl`, `data.dsl`, `data.v2.dsl`, `context_window.dsl`, `course.dsl`, `cms.dsl`) until their pages are intentionally ported.

The cutover is deliberately incremental. The provider exposes both families, the React renderer can consume Widget IR from either family, and the v3 example gallery proves the new module in both Goja golden tests and a browser-backed xgoja preview.

## Current module families

| Family | Modules | Intended use |
| --- | --- | --- |
| Widget DSL v3 | `widget.dsl` | New pages, new hosts, example-gallery work, and migration targets. Exposes namespaces such as `widget.ui`, `widget.data`, `widget.course`, `widget.context`, `widget.schedule`, and `widget.time`. |
| Legacy split DSL | `ui.dsl`, `data.dsl`, `data.v2.dsl`, `context_window.dsl`, `course.dsl`, `cms.dsl` | Existing first-party pages and external scripts that already run. Keep these selected until the page is ported and validated. |
| Removed bucket DSL | `rag.dsl` | Not exposed. Do not add compatibility shims for it. |

## xgoja module selection

### New v3 host

Select `widget.dsl` from `rag-widget-site`. Browser hosts still need the HTTP and asset modules that serve the SPA, but they should select only `widget.dsl` from the widget provider.

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

The reference implementation is `examples/xgoja-widgetdsl-v3/xgoja.yaml`. It also selects:

- `go-go-goja-host/fs` as `fs:assets` for embedded SPA/example files;
- `go-go-goja-http/express` as `express` for HTTP routes;
- `widget.dsl` as the only widget authoring module.

### Existing host during migration

Keep legacy modules selected while migrating page-by-page. It is safe to select `widget.dsl` alongside the old modules when a binary needs both old and new pages.

```yaml
runtime:
  modules:
    - provider: rag-widget-site
      name: ui.dsl
      as: ui.dsl
    - provider: rag-widget-site
      name: data.dsl
      as: data.dsl
    - provider: rag-widget-site
      name: data.v2.dsl
      as: data.v2.dsl
    - provider: rag-widget-site
      name: context_window.dsl
      as: context_window.dsl
    - provider: rag-widget-site
      name: course.dsl
      as: course.dsl
    - provider: rag-widget-site
      name: cms.dsl
      as: cms.dsl
    - provider: rag-widget-site
      name: widget.dsl
      as: widget.dsl
```

Remove the legacy module entries only after every script in that host stops importing them.

## Authoring pattern in v3

Old split modules usually look like this:

```js
const ui = require("ui.dsl")
const data = require("data.dsl")
const course = require("course.dsl")

return ui.page({
  id: "handouts",
  title: "Handouts",
  sections: [
    course.recipes.handout({
      documents,
      selectedDocumentId,
      onSelect: { kind: "navigate", to: "?item=${document.id}" },
    }),
  ],
})
```

V3 uses one module and namespaces:

```js
const widget = require("widget.dsl")

function renderPage(query) {
  const selected = query.item || "overview"
  return widget.page("Handouts", (p) => {
    p.id("handouts")
    p.section("Course handouts", (s) => {
      s.view(
        widget.course.handouts({ documents, selectedDocumentId: selected }, (h) =>
          h.onSelect(widget.act.navigate("?item=${document.id}")),
        ),
      )
    })
  })
}
```

Prefer builder callbacks, fragments, typed helpers, and namespace-level intents. Use `widget.raw.component(...)` only when the v3 namespace intentionally lacks a helper and record the exception in the migration notes.

## Migration mapping

| Old pattern | V3 replacement | Notes |
| --- | --- | --- |
| `require("ui.dsl")` | `const widget = require("widget.dsl")` and `widget.ui.*` / `widget.page(...)` | `widget.page(title, builder)` is the preferred page entrypoint. |
| `ui.page({ id, title, sections })` | `widget.page(title, (p) => { p.id(id); p.section(...) })` | Keep page IDs stable so `/pages/{id}` routes do not change. |
| `ui.panel`, `ui.caption`, `ui.button`, `ui.markdown`, `ui.section` | `widget.ui.panel`, `widget.ui.caption`, `widget.ui.button`, `widget.ui.markdown`, page/section builders | Use generic UI helpers for layout and small affordances. |
| `data.dsl` table/field/cell helpers | `widget.data.schema`, `widget.data.collection`, `widget.data.cell`, `widget.data.matrix` | V3 reuses v2 data contracts where possible. |
| `context_window.dsl` diagrams/workspaces | `widget.context.diagram`, `widget.context.workspace`, `widget.context.styleSet`, `widget.context.intent.*` | Treat current panels as lowering details; author against context concepts. |
| `course.dsl` course studio, slides, handouts | `widget.course.shell`, `widget.course.slideDeck`, `widget.course.handouts`, `widget.course.intent.*` | Use action contexts such as `item.id` and `document.id` for navigation templates. |
| `cms.dsl` media/article helpers | `widget.cms.mediaLibrary`, `widget.cms.articleQueue`, `widget.cms.markdownEditor`, `widget.cms.intent.*` | Prefer CMS intents over hand-built action objects. |
| Hand-built action objects | `widget.act.navigate`, `widget.act.event`, `widget.act.copy`, domain intents | Prefer typed action helpers so frontend interpolation receives the expected context. |
| String paths into row/document fields | `widget.bind.context("document.id")`, `widget.bind.field("status")` | Use accessors rather than concatenating object values. |
| `raw.component(...)` for existing widgets | V3 domain helper, or a documented `widget.raw.component(...)` exception | Raw escape hatches should decrease over time. |

## Migration workflow

1. **Inventory imports.** Run the migration checker against the host scripts.
   ```bash
   python ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/scripts/02-report-legacy-widget-dsl-usage.py path/to/scripts
   ```
2. **Add `widget.dsl` beside old modules.** Do not remove old module entries yet.
3. **Port one route/page at a time.** Keep route IDs and page IDs stable.
4. **Prefer v3 examples as fixtures.** Compare against `pkg/widgetdsl/testdata/v3/examples` for course, CMS, data, context, schedule, and time patterns.
5. **Regenerate or inspect Widget IR.** For first-party v3 examples:
   ```bash
   go run ./cmd/widgetdsl-v3-examples --out pkg/widgetdsl/testdata/v3/golden
   go test ./pkg/widgetdsl/... -count=1
   ```
6. **Validate in React.** Use Storybook regression stories and/or the xgoja preview gallery.
   ```bash
   pnpm --dir packages/rag-evaluation-site typecheck
   pnpm --dir packages/rag-evaluation-site build-storybook
   xgoja build -f examples/xgoja-widgetdsl-v3/xgoja.yaml --output examples/xgoja-widgetdsl-v3/dist/widgetdsl-v3-examples
   examples/xgoja-widgetdsl-v3/dist/widgetdsl-v3-examples serve site start --http-listen 127.0.0.1:8098
   ```
7. **Remove legacy module selections only when the checker is clean.** If a host still imports `ui.dsl`, `data.dsl`, `data.v2.dsl`, `context_window.dsl`, `course.dsl`, or `cms.dsl`, keep those runtime module entries.

## Cutover policy

- Legacy split modules remain supported for current scripts.
- New first-party widget pages should start on `widget.dsl` unless a v3 helper is missing.
- Missing v3 helpers should be added to `widget.dsl` rather than papered over with long-lived raw component calls.
- Deprecation of old modules should wait until first-party course/CMS/handout pages have v3 fixtures and the migration checker reports no legacy imports for those hosts.
- Do not add a broad backwards-compatibility bucket module; explicit parallel modules make host dependencies visible.

## Validation checklist

A migrated page is ready when all of these are true:

- the page imports `widget.dsl` and no legacy split modules;
- the xgoja runtime selects `widget.dsl`;
- the page returns JSON-compatible Widget IR without browser globals;
- no browser-visible `[object Object]` text or URLs appear;
- action templates interpolate from supplied action contexts (`item.id`, `document.id`, row fields, etc.);
- `go test ./pkg/widgetdsl/... -count=1` passes if the page is part of committed examples;
- `pnpm --dir packages/rag-evaluation-site typecheck` passes after renderer contract changes;
- Storybook or xgoja preview covers any new component/contract shape.

## Related

- `reference/05-widget-dsl-v3-api-reference.md` — namespace inventory for `widget.dsl`.
- `design-doc/05-widget-dsl-v3-implementation-phases-and-task-tracker.md` — phase tracker and acceptance criteria.
- `examples/xgoja-widgetdsl-v3/README.md` — build/run instructions for the v3 preview site.
- `pkg/xgoja/providers/widgetsite/doc/01-widget-dsl-getting-started.md` — embedded xgoja provider tutorial.
- `packages/rag-evaluation-site/src/widgets/WidgetRenderer.v3-regressions.stories.tsx` — renderer regression stories for browser contract failures found during preview validation.
