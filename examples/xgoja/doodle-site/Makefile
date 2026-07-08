.DEFAULT_GOAL := build

REPO_ROOT := $(abspath ../../..)
XGOJA_ROOT := $(abspath ../../../../go-go-goja)
XGOJA := cd $(XGOJA_ROOT) && GOWORK=off go run ./cmd/xgoja
SPEC := $(CURDIR)/xgoja.v2.yaml
BIN := $(CURDIR)/dist/doodle-site
HTTP_ADDR := 127.0.0.1:18793

.PHONY: sync-app doctor list-modules build serve clean

# Refresh the embedded SPA from the canonical app-dist build.
sync-app:
	rm -rf $(CURDIR)/assets/public
	mkdir -p $(CURDIR)/assets/public
	cp -R $(REPO_ROOT)/packages/rag-evaluation-site/app-dist/index.html $(REPO_ROOT)/packages/rag-evaluation-site/app-dist/assets $(CURDIR)/assets/public/
	find $(CURDIR)/assets/public -name '*.map' -delete

doctor:
	$(XGOJA) doctor -f $(SPEC)

list-modules:
	$(XGOJA) list-modules -f $(SPEC)

build:
	$(XGOJA) build -f $(SPEC) --output $(BIN) --xgoja-replace $(XGOJA_ROOT) --keep-work

serve: build
	$(BIN) serve doodle site --http-listen $(HTTP_ADDR)

clean:
	rm -rf $(CURDIR)/dist $(CURDIR)/doodle.db
