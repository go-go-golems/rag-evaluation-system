# devctl workflow for the xgoja Widget Site example

Run these commands from `examples/xgoja/widget-site`.

```bash
devctl plugins list
devctl plan
devctl up --force --timeout 10m
```

`devctl up` runs the plugin build phase first, which executes:

1. `make sync-app` to refresh `assets/public` from `packages/rag-evaluation-site/app-dist`.
2. `make build` to regenerate `dist/rag-widget-xgoja-site` with the local `go-go-goja` checkout.
3. A supervised `widget-site` service running `serve sites demo`.

Useful URLs after startup:

- App root: `http://127.0.0.1:18791/`
- Action demo: `http://127.0.0.1:18791/pages/actions`
- Health: `http://127.0.0.1:18791/healthz`
- Widget JSON: `http://127.0.0.1:18791/api/widget/pages/actions`

Useful commands:

```bash
devctl status
devctl logs --service widget-site --stderr --follow
devctl down
```

The plugin also exposes helper commands through devctl:

```bash
devctl sync-app --timeout 10m
devctl smoke --timeout 10m
devctl clean
```

If the preferred port `18791` is busy, `config.mutate` chooses a free port and
prints it in `devctl plan`. Use `devctl status` to see the actual health URL for
the running service.
