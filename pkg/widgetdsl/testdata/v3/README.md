# widget.dsl v3 examples

These examples are plain Goja JavaScript scripts that import `widget.dsl`, build a `page`, and render to serializable Widget IR JSON.

Run all examples and write rendered JSON:

```bash
go run ./cmd/widgetdsl-v3-examples \
  --examples pkg/widgetdsl/testdata/v3/examples \
  --out /tmp/widgetdsl-v3-rendered
```

Print rendered JSON to stdout:

```bash
go run ./cmd/widgetdsl-v3-examples --stdout
```

Refresh committed golden snapshots after intentional IR changes:

```bash
WIDGETDSL_UPDATE_GOLDEN=1 go test ./pkg/widgetdsl -run TestWidgetV3GoldenExamplesRenderStableIR -count=1
```

The golden snapshots live in `pkg/widgetdsl/testdata/v3/golden`. They are useful for reviewing the exact Widget IR that frontend hosts would serve from `/api/widget/pages/{id}`.

## Included examples

1. `01-simple-table.js` — basic `data.collection` table.
2. `02-selectable-table.js` — URL-backed row selection and row action.
3. `03-master-detail-editor.js` — master-detail collection with edit actions.
4. `04-row-actions.js` — action columns using CMS intents.
5. `05-all-modules-gallery.js` — UI, CMS, course, context, schedule, and time in one page.
6. `06-admin-course-cms.js` — course admin/CMS slice.
7. `07-handouts-and-slide.js` — handouts and slide deck page.
