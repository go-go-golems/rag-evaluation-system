.PHONY: build run test clean web-build web-dev widgetdsl-sites-smoke bleve-knn-experiment docker-lint lint lintmax golangci-lint-install glazed-lint-build glazed-lint gosec govulncheck logcopter-generate logcopter-check goreleaser tag-major tag-minor tag-patch release bump-go-go-golems install

BINARY := rag-eval
MODULE := github.com/go-go-golems/rag-evaluation-system
CMD_DIR := ./cmd/$(BINARY)
DB_PATH ?= data/rag-eval.db
ADDRESS ?= 127.0.0.1:8772
FAISS_LIB_DIR ?= /usr/local/lib

GORELEASER_ARGS ?= --skip=sign --snapshot --clean
GORELEASER_TARGET ?= --single-target
GOLANGCI_LINT_VERSION ?= $(shell cat .golangci-lint-version)
GOLANGCI_LINT_BIN ?= $(CURDIR)/.bin/golangci-lint
GLAZED_LINT_BIN ?= /tmp/glazed-lint
GLAZED_LINT_PKG ?= github.com/go-go-golems/glazed/cmd/tools/glazed-lint
GLAZED_LINT_TOOL_VERSION ?= v1.3.5
GLAZED_LINT_FLAGS ?= -glazedclilint.allow-paths=pkg/analysis/,pkg/cli/,pkg/cmds/fields/,pkg/cmds/logging/,pkg/cmds/sources/,pkg/help/

# TODO: switch this back to ./... once github.com/go-go-golems/scraper no longer imports
# the removed github.com/go-go-golems/go-go-goja/engine package. Keep ttmp out of
# default checks because ticket-local experiments are not repository packages.
GO_PACKAGES ?= \
	./pkg/... \
	./internal/chunking \
	./internal/db \
	./internal/ingest \
	./internal/services/... \
	./internal/web
# Application commands are a mandatory Glazed-lint scope: direct Cobra/pflag
# command definitions must fail CI.
GLAZED_LINT_DIRS ?= $(GO_PACKAGES) ./cmd/rag-eval/...

build:
	GOWORK=off go generate ./...
	GOWORK=off go build $(GO_PACKAGES)

build-bin:
	GOWORK=off go build -o $(BINARY) $(CMD_DIR)

run: build-bin
	RAG_EVAL_DB=$(DB_PATH) RAG_EVAL_ADDRESS=$(ADDRESS) ./$(BINARY)

test:
	GOWORK=off go test $(GO_PACKAGES) -count=1

bleve-knn-experiment:
	GOWORK=off CGO_LDFLAGS="-L$(FAISS_LIB_DIR) -lfaiss_c -lfaiss -lstdc++ -lm" \
		go run -tags=vectors -ldflags "-r $(FAISS_LIB_DIR)" ./cmd/experiments/bleve-knn/

clean:
	rm -f $(BINARY)
	rm -rf data/

web-install:
	cd web && pnpm install

web-dev:
	cd web && pnpm dev

web-build:
	cd web && pnpm build

widgetdsl-sites-smoke:
	./scripts/test-widgetdsl-v3-sites.sh

# Build Go binary with embedded SPA.
build-full: web-build build-bin

docker-lint:
	docker run --rm -v $(shell pwd):/app -w /app golangci/golangci-lint:$(GOLANGCI_LINT_VERSION) golangci-lint run -v $(GO_PACKAGES)

golangci-lint-install:
	mkdir -p $(dir $(GOLANGCI_LINT_BIN))
	GOBIN=$(dir $(GOLANGCI_LINT_BIN)) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)

glazed-lint-build:
	@echo "Building glazed-lint from pinned tool module..."
	@echo "Installing $(GLAZED_LINT_PKG)@$(GLAZED_LINT_TOOL_VERSION)"
	@GOBIN=$(dir $(GLAZED_LINT_BIN)) GOWORK=off go install $(GLAZED_LINT_PKG)@$(GLAZED_LINT_TOOL_VERSION)

glazed-lint: glazed-lint-build
	GOWORK=off go vet -vettool=$(GLAZED_LINT_BIN) $(GLAZED_LINT_FLAGS) $(GLAZED_LINT_DIRS)

lint: golangci-lint-install glazed-lint-build
	$(GOLANGCI_LINT_BIN) run -v $(GO_PACKAGES)
	GOWORK=off go vet -vettool=$(GLAZED_LINT_BIN) $(GLAZED_LINT_FLAGS) $(GLAZED_LINT_DIRS)

lintmax: golangci-lint-install glazed-lint-build
	$(GOLANGCI_LINT_BIN) run -v --max-same-issues=100 $(GO_PACKAGES)
	GOWORK=off go vet -vettool=$(GLAZED_LINT_BIN) $(GLAZED_LINT_FLAGS) $(GLAZED_LINT_DIRS)

gosec:
	GOWORK=off go install github.com/securego/gosec/v2/cmd/gosec@latest
	gosec -exclude-generated -exclude=G101,G304,G301,G306 -exclude-dir=.history $(GO_PACKAGES)

govulncheck:
	GOWORK=off go install golang.org/x/vuln/cmd/govulncheck@latest
	govulncheck $(GO_PACKAGES)

logcopter-generate:
	GOWORK=off go generate ./...

logcopter-check:
	GOWORK=off go tool logcopter-gen -area-prefix go-go-golems.rag-evaluation-system -strip-prefix $(MODULE) -check ./pkg/... ./internal/... ./cmd/...

goreleaser:
	GOWORK=off goreleaser release $(GORELEASER_ARGS) $(GORELEASER_TARGET)

tag-major:
	git tag $$(svu major)

tag-minor:
	git tag $$(svu minor)

tag-patch:
	git tag $$(svu patch)

release:
	git push origin --tags
	GOWORK=off GOPROXY=proxy.golang.org go list -m $(MODULE)@$$(svu current)

bump-go-go-golems:
	@deps="$$(awk '/^require[[:space:]]+github\.com\/go-go-golems\// { print $$2 } /^[[:space:]]*github\.com\/go-go-golems\// { print $$1 }' go.mod | sort -u)"; \
	if [ -z "$$deps" ]; then \
		echo "No github.com/go-go-golems dependencies in go.mod"; \
	else \
		echo "Bumping go-go-golems dependencies:"; \
		echo "$$deps"; \
		for dep in $$deps; do GOWORK=off go get "$${dep}@latest"; done; \
	fi
	GOWORK=off go mod tidy

install:
	GOWORK=off go build -o ./dist/$(BINARY) $(CMD_DIR) && \
		cp ./dist/$(BINARY) $$(which $(BINARY))
