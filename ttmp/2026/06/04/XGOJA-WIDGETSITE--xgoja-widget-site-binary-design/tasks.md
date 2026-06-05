# Tasks

## Documentation and research deliverable

- [x] Create docmgr ticket `XGOJA-WIDGETSITE`.
- [x] Capture `xgoja help --all`.
- [x] Capture `xgoja help tutorial-http-serve-jsverbs`.
- [x] Capture `xgoja help tutorial-generated-runtime-package`.
- [x] Capture `xgoja help tutorial-static-assets-http-server`.
- [x] Inspect xgoja provider, host-service, asset, generated go.mod, HTTP, Express, and database source code.
- [x] Inspect RAG WidgetRenderer `widget.dsl` and frontend package sources relevant to xgoja integration.
- [x] Run upstream embedded-assets and HTTP-jsverbs example smoke targets.
- [x] Create ticket-local combined xgoja experiment under `scripts/`.
- [x] Record experiment validation and build failure outputs under `sources/`.
- [x] Write intern-facing design/implementation guide.
- [x] Write research logbook with usefulness/staleness/update notes for each resource.
- [x] Keep diary while working.
- [x] Relate key files to ticket documents.
- [x] Run `docmgr doctor --ticket XGOJA-WIDGETSITE --stale-after 30`.
- [x] Upload design bundle to reMarkable.
- [x] Commit ticket materials.

## Future implementation tasks

- [x] Add real provider package `pkg/xgoja/providers/widgetsite` registering `widget.dsl` and `rag.dsl`.
- [x] Add provider unit tests for module registration and runtime `require("widget.dsl")` behavior.
- [x] Add `examples/xgoja/widget-site` with tiny static webapp, jsverb, xgoja.yaml, and Makefile smoke.
- [x] Build generated binary with `--xgoja-replace` and `packages[].replace` for local development.
- [x] Add curl smoke for `/healthz`, `/static/`, and `/api/widget/pages/demo`.
- [x] Replace tiny webapp with `packages/rag-evaluation-site/app-dist` once the xgoja path is stable.
- [x] Add Playwright smoke for the generated xgoja binary with embedded WidgetRenderer app.
- [x] Decide whether to implement preconfigured database host-service contribution or use `db.configure()` only for demo.
- [ ] Consider adding Express `spaFromAssetsModule` helper for React client-side route fallback.
- [ ] Reconcile or restore the missing `examples/xgoja/14-generated-runtime-package` referenced by help.
- [ ] Implement Express app.spaFromAssetsModule helper in go-go-goja for API-safe React SPA fallback from embedded assets.
- [ ] Use app.spaFromAssetsModule in examples/xgoja/widget-site so /pages/demo renders the React app while /api/widget remains dynamic.
- [ ] Implement preconfigured database host-provider config in go-go-goja using driverName/dataSourceName and disabled configure().
- [ ] Use preconfigured db config in examples/xgoja/widget-site and remove JavaScript db.configure() from the jsverb.
- [ ] Add smoke coverage proving db.configure() is disabled when the xgoja db module is preconfigured.
- [x] Add provider-bundled Glazed help entries for widget.dsl getting started and JavaScript API reference.
