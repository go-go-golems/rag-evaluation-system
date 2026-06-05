.DEFAULT_GOAL := smoke

REPO_ROOT := $(abspath ../../..)
XGOJA_ROOT := $(abspath ../../../../go-go-goja)
XGOJA := cd $(XGOJA_ROOT) && GOWORK=off go run ./cmd/xgoja
BIN := $(CURDIR)/dist/rag-widget-xgoja-site
HTTP_ADDR := 127.0.0.1:18791

.PHONY: smoke sync-app doctor list-modules build serve-smoke clean

smoke: doctor list-modules build serve-smoke

sync-app:
	rm -rf $(CURDIR)/assets/public
	mkdir -p $(CURDIR)/assets/public
	cp -R $(REPO_ROOT)/packages/rag-evaluation-site/app-dist/index.html $(REPO_ROOT)/packages/rag-evaluation-site/app-dist/assets $(CURDIR)/assets/public/
	find $(CURDIR)/assets/public -name '*.map' -delete
	python3 -c 'from pathlib import Path; root=Path("$(CURDIR)/assets/public/assets"); [p.write_text(p.read_text().replace("\\n//# sourceMappingURL=" + p.name + ".map", "").replace("//# sourceMappingURL=" + p.name + ".map", "")) for p in root.glob("*.js")]'

doctor:
	$(XGOJA) doctor -f $(CURDIR)/xgoja.yaml

list-modules:
	$(XGOJA) list-modules -f $(CURDIR)/xgoja.yaml

build:
	$(XGOJA) build -f $(CURDIR)/xgoja.yaml --output $(BIN) --xgoja-replace $(XGOJA_ROOT) --keep-work

serve-smoke: build
	@set -eu; \
	log=$$(mktemp); \
	$(BIN) serve sites demo --http-listen $(HTTP_ADDR) >$$log 2>&1 & \
	pid=$$!; \
	trap 'kill $$pid >/dev/null 2>&1 || true; wait $$pid >/dev/null 2>&1 || true; rm -f $$log /tmp/rag-widget-xgoja-page.json /tmp/rag-widget-xgoja-spa.html /tmp/rag-widget-xgoja-health.json /tmp/rag-widget-xgoja-root.html /tmp/rag-widget-xgoja-asset.js' EXIT; \
	for i in $$(seq 1 100); do \
		if curl -fsS http://$(HTTP_ADDR)/healthz >/tmp/rag-widget-xgoja-health.json 2>/dev/null; then break; fi; \
		if ! kill -0 $$pid >/dev/null 2>&1; then cat $$log; exit 1; fi; \
		sleep 0.1; \
	done; \
	grep -q '"ok":true' /tmp/rag-widget-xgoja-health.json; \
	grep -q '"site":"rag-widget-xgoja-site"' /tmp/rag-widget-xgoja-health.json; \
	curl -fsS -o /dev/null http://$(HTTP_ADDR)/favicon.ico; \
	curl -fsS http://$(HTTP_ADDR)/api/widget/pages/demo >/tmp/rag-widget-xgoja-page.json; \
	grep -q '"schemaVersion":"0.1.0"' /tmp/rag-widget-xgoja-page.json; \
	grep -q '"type":"DataTable"' /tmp/rag-widget-xgoja-page.json; \
	grep -q 'Fast Growing Trees' /tmp/rag-widget-xgoja-page.json; \
	curl -fsS http://$(HTTP_ADDR)/ >/tmp/rag-widget-xgoja-root.html; \
	grep -q 'RAG Evaluation Site' /tmp/rag-widget-xgoja-root.html; \
	grep -q '/assets/index-' /tmp/rag-widget-xgoja-root.html; \
	curl -fsS http://$(HTTP_ADDR)/pages/demo >/tmp/rag-widget-xgoja-spa.html; \
	grep -q 'RAG Evaluation Site' /tmp/rag-widget-xgoja-spa.html; \
	curl -fsS http://$(HTTP_ADDR)/assets/$$(basename $$(find $(CURDIR)/assets/public/assets -name 'index-*.js' | head -1)) >/tmp/rag-widget-xgoja-asset.js; \
	grep -q 'RagEvaluationSiteApp' /tmp/rag-widget-xgoja-asset.js; \
	kill $$pid >/dev/null 2>&1 || true; \
	wait $$pid >/dev/null 2>&1 || true

clean:
	rm -rf $(CURDIR)/dist
