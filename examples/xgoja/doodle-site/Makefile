.DEFAULT_GOAL := smoke

REPO_ROOT := $(abspath ../../..)
XGOJA_ROOT := $(abspath ../../../../go-go-goja)
XGOJA := cd $(XGOJA_ROOT) && GOWORK=off go run ./cmd/xgoja
SPEC := $(CURDIR)/xgoja.v2.yaml
BIN := $(CURDIR)/dist/doodle-site
HTTP_ADDR := 127.0.0.1:18793

.PHONY: smoke sync-app doctor list-modules build serve serve-smoke clean

smoke: doctor list-modules serve-smoke

# Refresh the embedded SPA from the canonical app-dist build.
sync-app:
	pnpm --dir $(REPO_ROOT)/packages/rag-evaluation-site build:app
	rm -rf $(CURDIR)/assets/public
	mkdir -p $(CURDIR)/assets/public
	cp -R $(REPO_ROOT)/packages/rag-evaluation-site/app-dist/index.html $(REPO_ROOT)/packages/rag-evaluation-site/app-dist/assets $(CURDIR)/assets/public/
	find $(CURDIR)/assets/public -name '*.map' -delete

doctor:
	$(XGOJA) doctor -f $(SPEC)

list-modules:
	$(XGOJA) list-modules -f $(SPEC)

build: sync-app
	$(XGOJA) build -f $(SPEC) --output $(BIN) --xgoja-replace $(XGOJA_ROOT) --keep-work

serve: build
	$(BIN) serve doodle site --http-listen $(HTTP_ADDR)

serve-smoke: build
	@set -eu; \
	rm -f $(CURDIR)/doodle.db; log=$$(mktemp); index=$$(mktemp); poll=$$(mktemp); \
	$(BIN) serve doodle site --http-listen $(HTTP_ADDR) >$$log 2>&1 & pid=$$!; \
	trap 'kill $$pid >/dev/null 2>&1 || true; wait $$pid >/dev/null 2>&1 || true; rm -f $$log $$index $$poll $(CURDIR)/doodle.db' EXIT; \
	for i in $$(seq 1 100); do curl -fsS http://$(HTTP_ADDR)/healthz >/dev/null 2>&1 && break; kill -0 $$pid >/dev/null 2>&1 || { cat $$log; exit 1; }; sleep 0.1; done; \
	curl -fsS http://$(HTTP_ADDR)/api/widget/pages/index >$$index; \
	grep -q '"type":"DataTable"' $$index; grep -q '"type":"StatusText"' $$index; \
	curl -fsS 'http://$(HTTP_ADDR)/api/widget/pages/poll?poll=1&view=week' >$$poll; \
	for type in MatrixGrid MonthGrid TimeGrid FormPanel ShareLink; do grep -q "\"type\":\"$$type\"" $$poll; done

clean:
	rm -rf $(CURDIR)/dist $(CURDIR)/doodle.db
