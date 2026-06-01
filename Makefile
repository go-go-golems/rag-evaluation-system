.PHONY: build run test clean web-build web-dev

BINARY := rag-eval
DB_PATH ?= data/rag-eval.db
ADDRESS ?= 127.0.0.1:8772

build:
	go build -o $(BINARY) ./cmd/rag-eval

run: build
	RAG_EVAL_DB=$(DB_PATH) RAG_EVAL_ADDRESS=$(ADDRESS) ./$(BINARY)

test:
	go test ./... -count=1

clean:
	rm -f $(BINARY)
	rm -rf data/

web-install:
	cd web && pnpm install

web-dev:
	cd web && pnpm dev

web-build:
	cd web && pnpm build


# Build Go binary with embedded SPA
build-full: web-build
	go build -o $(BINARY) ./cmd/rag-eval

lint:
	golangci-lint run -v
