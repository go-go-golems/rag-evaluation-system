---
Title: "Widget DSL SPA Bundling"
Slug: widget-dsl-spa-bundling
Short: "Bundle the widget.dsl provider and embedded React renderer into xgoja hosts."
Topics:
- xgoja
- widget-dsl
- spa
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
---

# Widget DSL SPA Bundling

A generated host selects `rag-widget-site/widget.dsl` and embeds the `@go-go-golems/rag-evaluation-site` assets. The provider exposes no split modules.

```yaml
modules:
  - package: rag-widget-site
    name: widget.dsl
    as: widget.dsl
```

Serve `/api/widget/pages/{id}` as Widget Page JSON and `/api/widget/actions/{name}` as action endpoints. Server actions may return `refresh`, `toast`, `error`, and `fieldErrors`; the renderer routes these to navigation refresh, the live region, and active FormDialog respectively. Regenerate embedded assets after changing adapters and run the provider smoke test before release.
