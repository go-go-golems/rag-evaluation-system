---
Title: "Widget DSL Getting Started"
Slug: widget-dsl-getting-started
Short: "Author React-rendered Widget IR pages in xgoja with widget.dsl."
Topics:
- xgoja
- widget-dsl
- widget-ir
Commands:
- xgoja build
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
---

# Widget DSL Getting Started

Select exactly one module in `xgoja.yaml`:

```yaml
modules:
  - package: rag-widget-site
    name: widget.dsl
    as: widget.dsl
```

Author a page through typed namespaces:

```js
const widget = require("widget.dsl");
const jobs = widget.data.collection(rows, c => c
  .schema(fields.build())
  .search(s => s.query("q").submit(searchAction))
  .paginate(p => p.current(page).size(pageSize).total(total).sizes(20, 50, 100).onChange(pageAction))
  .table(t => t.keyboard(k => k.mode("rows")).rowSelect(openAction)));
const page = widget.page("Jobs", p => p.section("Queue", s => s.view(jobs)));
```

Use `widget.act`, `widget.bind`, and intent-level namespace helpers. Split modules and raw component construction are not part of the public provider.
