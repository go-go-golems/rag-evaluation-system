# xgoja widget.dsl v3 example site

This xgoja app serves the committed `widget.dsl` v3 examples through the React WidgetRenderer SPA. It is the reference integration for a new host that selects `widget.dsl` from the `rag-widget-site` provider while leaving the legacy split modules unselected.

## Runtime module selection

`xgoja.yaml` selects the HTTP and asset provider modules needed to serve the preview, plus exactly one widget authoring module:

```yaml
runtime:
  modules:
    - provider: rag-widget-site
      name: widget.dsl
      as: widget.dsl
```

The JavaScript examples then use:

```js
const widget = require("widget.dsl")
```

See the ticket migration guide for side-by-side legacy/v3 migration notes:

- `ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/reference/06-widget-dsl-v3-integration-and-migration-guide.md`

## Build

```bash
pnpm --dir packages/rag-evaluation-site build:app
xgoja doctor -f examples/xgoja-widgetdsl-v3/xgoja.yaml
xgoja build -f examples/xgoja-widgetdsl-v3/xgoja.yaml \
  --output examples/xgoja-widgetdsl-v3/dist/widgetdsl-v3-examples
```

## Run

```bash
examples/xgoja-widgetdsl-v3/dist/widgetdsl-v3-examples \
  serve site start \
  --http-listen 127.0.0.1:8098
```

Open <http://127.0.0.1:8098/pages/index>.

## Useful checks

```bash
go test ./pkg/widgetdsl/... ./pkg/xgoja/providers/widgetsite/... -count=1
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build-storybook
python ttmp/2026/07/06/RAGEVAL-SCHEDULE-WIDGETS--calendar-scheduling-widgets-on-generic-base-engines/scripts/02-report-legacy-widget-dsl-usage.py \
  pkg/widgetdsl/testdata/v3/examples examples/xgoja-widgetdsl-v3/jsverbs
```

The migration checker should report no legacy split-module imports for this app. It may report deliberate `widget.raw.component(...)` examples; treat those as follow-up candidates for future v3 helper coverage.
