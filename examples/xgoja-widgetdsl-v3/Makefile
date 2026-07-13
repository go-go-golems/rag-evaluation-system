.DEFAULT_GOAL := smoke

REPO_ROOT := $(abspath ../..)
XGOJA_ROOT := $(abspath ../../../go-go-goja)
XGOJA := cd $(XGOJA_ROOT) && GOWORK=off go run ./cmd/xgoja
SPEC := $(CURDIR)/xgoja.yaml
BIN := $(CURDIR)/dist/widgetdsl-v3-examples
HTTP_ADDR := 127.0.0.1:18795

.PHONY: smoke doctor list-modules build serve serve-smoke clean

smoke: doctor list-modules serve-smoke

doctor:
	$(XGOJA) doctor -f $(SPEC)

list-modules:
	$(XGOJA) list-modules -f $(SPEC)

build:
	pnpm --dir $(REPO_ROOT)/packages/rag-evaluation-site build:app
	$(XGOJA) build -f $(SPEC) --output $(BIN) --xgoja-replace $(XGOJA_ROOT) --keep-work

serve: build
	$(BIN) serve site start --http-listen $(HTTP_ADDR)

serve-smoke: build
	@set -eu; \
	log=$$(mktemp); page=$$(mktemp); \
	$(BIN) serve site start --http-listen $(HTTP_ADDR) >$$log 2>&1 & pid=$$!; \
	trap 'kill $$pid >/dev/null 2>&1 || true; wait $$pid >/dev/null 2>&1 || true; rm -f $$log $$page' EXIT; \
	for i in $$(seq 1 100); do curl -fsS http://$(HTTP_ADDR)/healthz >/dev/null 2>&1 && break; kill -0 $$pid >/dev/null 2>&1 || { cat $$log; exit 1; }; sleep 0.1; done; \
	curl -fsS http://$(HTTP_ADDR)/api/widget/pages/index | grep -q 'widget.dsl v3 examples'; \
	for source in $(REPO_ROOT)/pkg/widgetdsl/testdata/v3/examples/*.js; do \
		id=$$(basename $$source .js); curl -fsS http://$(HTTP_ADDR)/api/widget/pages/$$id >$$page; \
		grep -q '"schemaVersion"' $$page || { echo "missing schemaVersion for $$id"; exit 1; }; \
		grep -q '"kind":"component"' $$page || { echo "missing component root for $$id"; exit 1; }; \
	done; \
	curl -fsS -X POST http://$(HTTP_ADDR)/api/widget/actions/release-check -H 'Content-Type: application/json' -d '{"payload":{"source":"smoke"}}' | grep -q 'Preview action accepted'

clean:
	rm -rf $(CURDIR)/dist $(CURDIR)/js/types
