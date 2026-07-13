---
Title: "Widget DSL JavaScript API Reference"
Slug: widget-dsl-js-api-reference
Short: "Conceptual reference for the single widget.dsl authoring language."
Topics:
- widget-dsl
- widget-ir
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
---

# Widget DSL JavaScript API Reference

`widget.dsl` exposes these namespaces:

- `widget.ui`: content, layout, forms, uploads, and FormDialog.
- `widget.data`: field schemas, collections, matrices, cells, selection, and activity feeds.
- `widget.context`, `course`, `cms`, `crm`, `schedule`, and `time`: stable domain views and intents.
- `widget.act`: server, navigation, download, event, copy, and overlay actions.
- `widget.bind`: field, path, map, template, context, and constant bindings.

All mutable builders support `.use(fragment)`. Browser behavior must be represented as `widget.act` data; Goja callbacks only configure author-time specs. See `widget-dsl-v3-api-reference` for generated method-level details.
