.DEFAULT_GOAL := smoke

REPO_ROOT := $(abspath ../../..)
XGOJA_ROOT := $(abspath ../../../../go-go-goja)
XGOJA := cd $(XGOJA_ROOT) && GOWORK=off go run ./cmd/xgoja
SPEC := $(CURDIR)/xgoja.v2.yaml
BIN := $(CURDIR)/dist/workshop-crm-site
HTTP_ADDR := 127.0.0.1:18794

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
	$(BIN) serve workshopcrm site --http-listen $(HTTP_ADDR)

serve-smoke: build
	@set -eu; \
	rm -f $(CURDIR)/workshop-crm.db; log=$$(mktemp); index=$$(mktemp); pipeline=$$(mktemp); opportunity=$$(mktemp); \
	$(BIN) serve workshopcrm site --http-listen $(HTTP_ADDR) >$$log 2>&1 & pid=$$!; \
	trap 'kill $$pid >/dev/null 2>&1 || true; wait $$pid >/dev/null 2>&1 || true; rm -f $$log $$index $$pipeline $$opportunity $(CURDIR)/workshop-crm.db' EXIT; \
	for i in $$(seq 1 100); do curl -fsS http://$(HTTP_ADDR)/healthz >/dev/null 2>&1 && break; kill -0 $$pid >/dev/null 2>&1 || { cat $$log; exit 1; }; sleep 0.1; done; \
	curl -fsS http://$(HTTP_ADDR)/api/widget/pages/index >$$index; \
	grep -q '"type":"SegmentedBar"' $$index; grep -q '"type":"StatTile"' $$index; \
	curl -fsS http://$(HTTP_ADDR)/api/widget/pages/pipeline >$$pipeline; grep -q '"type":"BoardEngine"' $$pipeline; \
	curl -fsS 'http://$(HTTP_ADDR)/api/widget/pages/opportunity?deal=1' >$$opportunity; \
	for type in RecordFieldList ActivityFeed; do grep -q "\"type\":\"$$type\"" $$opportunity; done; \
	curl -fsS -X POST http://$(HTTP_ADDR)/api/widget/actions/crm.deal.move -H 'Content-Type: application/json' -d '{"payload":{"dealId":1,"toStage":"proposal"}}' | grep -q '"refresh":true'

clean:
	rm -rf $(CURDIR)/dist $(CURDIR)/workshop-crm.db
