# xgoja widget.dsl v3 example site

This xgoja app serves the committed `widget.dsl` v3 examples through the React WidgetRenderer SPA.

Build:

```bash
xgoja doctor -f examples/xgoja-widgetdsl-v3/xgoja.yaml
xgoja build -f examples/xgoja-widgetdsl-v3/xgoja.yaml \
  --output examples/xgoja-widgetdsl-v3/dist/widgetdsl-v3-examples
```

Run:

```bash
examples/xgoja-widgetdsl-v3/dist/widgetdsl-v3-examples \
  serve site start \
  --http-listen 127.0.0.1:8098
```

Open <http://127.0.0.1:8098/pages/index>.
