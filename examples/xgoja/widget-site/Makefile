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

doctor:
	$(XGOJA) doctor -f $(CURDIR)/xgoja.yaml

list-modules:
	$(XGOJA) list-modules -f $(CURDIR)/xgoja.yaml

build:
	$(XGOJA) build -f $(CURDIR)/xgoja.yaml --output $(BIN) --xgoja-replace $(XGOJA_ROOT) --keep-work

serve-smoke: build
	@set -eu; \
	log=$$(mktemp); page=$$(mktemp); action=$$(mktemp); invalid=$$(mktemp); \
	$(BIN) serve sites demo --http-listen $(HTTP_ADDR) >$$log 2>&1 & pid=$$!; \
	trap 'kill $$pid >/dev/null 2>&1 || true; wait $$pid >/dev/null 2>&1 || true; rm -f $$log $$page $$action $$invalid' EXIT; \
	for i in $$(seq 1 100); do curl -fsS http://$(HTTP_ADDR)/healthz >/dev/null 2>&1 && break; kill -0 $$pid >/dev/null 2>&1 || { cat $$log; exit 1; }; sleep 0.1; done; \
	curl -fsS http://$(HTTP_ADDR)/healthz | grep -q '"module":"widget.dsl"'; \
	curl -fsS http://$(HTTP_ADDR)/api/widget/pages/demo >$$page; \
	for type in DataTable SearchField Pagination FormDialog; do grep -q "\"type\":\"$$type\"" $$page; done; \
	curl -fsS -X POST http://$(HTTP_ADDR)/api/widget/actions/cycle-status -H 'Content-Type: application/json' -d '{"payload":{"id":2},"context":{"componentType":"DataTable"}}' >$$action; \
	grep -q '"refresh":true' $$action; \
	curl -sS -o $$invalid -w '%{http_code}' -X POST http://$(HTTP_ADDR)/api/widget/actions/save-note -H 'Content-Type: application/json' -d '{"payload":{"id":2,"notes":""},"context":{"overlayId":"edit-note"}}' | grep -q 422; \
	grep -q '"fieldErrors"' $$invalid

clean:
	rm -rf $(CURDIR)/dist
