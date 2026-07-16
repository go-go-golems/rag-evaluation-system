---
Title: "Widget DSL JavaScript API Reference"
Slug: widget-dsl-js-api-reference
Short: "Conceptual reference for the single widget.dsl authoring language."
Topics:
- widget-dsl
- widget-ir
- javascript
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

`widget.dsl` is the typed authoring language for serializable pages, components, actions, and interaction metadata. Goja callbacks configure builders on the server; they never become browser callbacks. This boundary lets generated xgoja hosts describe rich behavior while the React application retains rendering, accessibility, and network ownership.

## Public namespaces

The module groups helpers by intent so page scripts can use stable semantic APIs instead of recreating renderer internals. All mutable builders support `.use(fragment)` for reusable configuration.

- `widget.ui`: content, layout, forms, uploads, and FormDialog.
- `widget.data`: field schemas, collections, matrices, cells, selection, and activity feeds.
- `widget.context`, `course`, `cms`, `crm`, `schedule`, and `time`: stable domain views and intents.
- `widget.act`: server, navigation, download, event, copy, and overlay actions.
- `widget.bind`: field, path, map, template, context, and constant bindings.

Prefer typed helpers whenever they exist. Use `widget.raw` only as a narrow migration escape hatch for a component that does not yet have a typed v3 surface.

## Page keyboard shortcuts

Page shortcuts bind a logical browser key to an ordinary `widget.act` action. They belong to the page envelope rather than a button, so the same API can activate server mutations, navigation, events, downloads, copy actions, or overlays. Visible controls remain the primary interaction; shortcuts are accelerators.

```js
const accept = widget.act.server("triage.accept")
const reject = widget.act.server("triage.reject")
const skip = widget.act.server("triage.skip")

const page = widget.page("Triage", (page) =>
  page
    .shortcuts((keys) =>
      keys
        .bind("accept", "y", accept, { label: "Yes" })
        .bind("reject", "n", reject, { label: "No" })
        .bind("skip", "s", skip, { label: "Skip" }),
    )
    .section("Current job", (section) =>
      section.view(widget.ui.inline(
        widget.ui.button("Yes", accept),
        widget.ui.button("No", reject),
        widget.ui.button("Skip", skip),
      )),
    ),
)
```

Each binding requires:

| Field | Meaning |
| --- | --- |
| `id` | Stable command identity for diagnostics, help, tests, and future remapping. |
| `key` | A logical `KeyboardEvent.key` value such as `y`, `Enter`, or `Escape`. |
| `action` | Any serializable action returned by `widget.act`. |
| `label` | Human-readable command name displayed by shortcut help. |
| `modifiers` | Optional array containing `Alt`, `Control`, `Meta`, or `Shift`. |
| `preventDefault` | Whether a match prevents the browser default; defaults to `true`. |
| `allowRepeat` | Whether holding the key may dispatch repeatedly; defaults to `false`. |

Represent modifiers separately instead of writing `Ctrl+S` in `key`:

```js
page.shortcuts((keys) =>
  keys.bind("save", "s", widget.act.server("record.save"), {
    label: "Save",
    modifiers: ["Control"],
  }),
)
```

The React host dispatches a matched command with this action context:

```json
{
  "componentType": "PageShortcut",
  "pageId": "triage",
  "shortcutId": "accept",
  "key": "y"
}
```

The runtime ignores shortcuts while the user types in editable controls, uses an IME, interacts with a modal, or operates a nested keyboard scope such as `DataTable`. Component handlers win when they consume an event. Repeated keydown events are ignored unless the binding explicitly enables repeat.

Unmodified character shortcuts require a user-facing disable, remap, or focus-only mechanism. `RagEvaluationSiteApp` supplies generated shortcut help and stores the enable/disable preference in the browser. Hosts that render pages without the default app must integrate the exported `usePageShortcuts` hook and provide equivalent discoverability and preference controls.

## Action and binding boundary

Actions describe what the browser should do, while bindings describe where runtime values come from. Keep both as data so the server output remains JSON-compatible and testable.

- Server actions POST resolved `payload` values to `/api/widget/actions/{name}`.
- Event actions dispatch resolved values through `CustomEvent.detail`.
- Navigation actions preserve or replace query state according to their options.
- Copy, download, and overlay actions execute through the same central dispatcher.

A page shortcut reuses this exact pipeline. Do not create a second fetch call or raw `onKeyDown` callback for shortcut behavior.

## Troubleshooting

Shortcut failures usually come from invalid chord representation, focus safety rules, or a host that renders only the root node. Match the symptom before changing the action itself.

| Problem | Cause | Solution |
| --- | --- | --- |
| `page.shortcuts is not a function` | The generated host predates the page-shortcut API. | Rebuild the xgoja host with the updated `rag-widget-site` provider. |
| Validation reports a duplicate chord. | Two bindings normalize to the same key and modifiers. | Give each page command a unique chord; letter case alone is not distinct. |
| `Control+y` is rejected as a key. | Modifiers were serialized into `key`. | Use `key: "y"` with `modifiers: ["Control"]`. |
| A shortcut does not run while typing or in a dialog. | The host intentionally suppresses page commands in editable and modal contexts. | Keep editor/dialog commands component-scoped; do not weaken the page safety guard. |
| A table-focused shortcut does not run. | `DataTable` owns its nested keyboard scope. | Use `table.command(...)` for row commands or move focus out of the table for page commands. |
| The page shows no shortcut help. | The host renders `WidgetRenderer` directly and never installs page behavior. | Use `RagEvaluationSiteApp` or integrate `usePageShortcuts` and a help/preference surface. |
| A held key dispatches only once. | `allowRepeat` defaults to `false` to protect action endpoints. | Enable repeat only for an idempotent command that is safe to repeat. |

## See Also

These entries provide runnable examples, generated method inventory, and host integration details.

- `widget-dsl-v3-examples` — complete page, action, shortcut, scheduling, and CRM examples.
- `widget-dsl-v3-api-reference` — descriptor-generated namespace, builder, and action-context inventory.
- `widget-dsl-getting-started` — provider selection and first-page setup.
- `widget-dsl-spa-bundling` — package and embed the React host with a generated xgoja binary.
- [`pkg/widgetdsl/testdata/v3/examples/43-page-shortcuts.js`](../../../../widgetdsl/testdata/v3/examples/43-page-shortcuts.js) — executable golden shortcut example.
