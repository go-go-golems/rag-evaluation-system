# rag-evaluation-system

A workflow-driven RAG (Retrieval-Augmented Generation) evaluation toolkit — document indexing, retrieval search, and an interactive Widget IR frontend, all backed by SQLite.

## Overview

This project provides:

- **Document ingestion** — source loading, chunking with multiple strategies, and embedding generation.
- **Dual-index retrieval** — BM25 via Bleve, vector search via Geppetto/Pinocchio, and hybrid RRF fusion.
- **Workflow orchestration** — durable, scraper-backed intake pipelines that coordinate chunking, embedding, and indexing.
- **Interactive playground** — a React-based WidgetRenderer frontend that renders structured Widget IR, with semantic panels for context windows, transcripts, annotations, slides, and handouts.

## Quick start

```bash
# Build the CLI and run the HTTP server
make build-full
./rag-eval serve --db data/rag-eval.db
```

Open `http://127.0.0.1:8772` in your browser. The server serves the React UI and provides REST API endpoints at `/api/v1/*`.

## CLI reference

The `rag-eval` binary is a Glazed Cobra application. Every subcommand produces structured table output by default and supports JSON, CSV, YAML, and other formats.

### Chunking

```bash
rag-eval chunk apply --doc-id <ID> --strategy fixed --chunk-size 500 --overlap 50
rag-eval chunk strategies --db data/rag-eval.db
```

### Document management

```bash
rag-eval source add --name <name> --type pdf <path>
rag-eval document list --db data/rag-eval.db
```

### Embeddings

```bash
rag-eval embedding apply --doc-id <ID> --profile openai-embedding-small \
  --profile-registries ~/.config/pinocchio/profiles.yaml
```

### Search

```bash
# BM25
rag-eval search query --index-id bm25-my-index --query "search text"

# Vector
rag-eval search vector --query "search text" --strategy-id fixed-1200-150 \
  --profile openai-embedding-small

# Hybrid
rag-eval search hybrid --query "search text" --index-id bm25-my-index \
  --strategy-id fixed-1200-150 --profile openai-embedding-small
```

### Workflows

Submit a complete intake pipeline (chunk → embed → BM25):

```bash
rag-eval workflow submit-intake \
  --document-ids doc-1,doc-2 \
  --strategy fixed --chunk-size 1200 --overlap 150 \
  --embeddings-type openai --profile openai-embedding-small
```

Run the local scheduler:

```bash
rag-eval workflow run-worker
rag-eval workflow status
```

### HTTP server

```bash
rag-eval serve --address 127.0.0.1:8772 --log-level debug
```

The server exposes:

| Route | Purpose |
| --- | --- |
| `/` | React WidgetRenderer SPA |
| `/api/v1/*` | REST API endpoints |
| `/api/v1/health` | Health check |

## Architecture

```
cmd/rag-eval/         Glazed Cobra CLI commands
  cmds/chunk/         Document chunking
  cmds/document/      Document CRUD
  cmds/embedding/     Embedding generation
  cmds/search/        BM25, vector, hybrid search
  cmds/serve/         HTTP server + embedded SPA
  cmds/workflow/      Scraper-backed durable workflows
  cmds/source/        Source document management

internal/services/    Business logic
  chunking/           Chunking strategies
  embedding/          Provider resolution (Geppetto/Pinocchio)
  search/             Bleve BM25 + vector retrieval
  corpus/             Corpus metadata
  document/           Document lifecycle
  documentprocessing/ Preprocessing artifacts

pkg/                  Reusable libraries
  defaultspa/         Embedded WidgetRenderer SPA (Go embed.FS)
  widgetdsl/          Goja Widget IR DSL
  widgetschema/       Widget IR schema validation
  xgoja/providers/widgetsite/  xgoja Widget DSL provider + help docs

packages/rag-evaluation-site/  React component library
  components/        Foundation atoms, molecules, organisms
  widgets/           Widget IR types and WidgetRenderer
```

## Researchctl laboratory integration

The authoritative cross-purpose execution path exports a pure `rag-retrieval-spec/v1` value and emits observations through `pkg/raglab.ObservationExecutor`. Researchctl owns runs, attempts, retries, timestamps, terminal state, SQLite, artifact custody, import, and export. This repository owns retrieval, fusion, reranking, relevance, and query-trace semantics.

```bash
# Inspect a lifecycle-free JavaScript specification export.
go run ./cmd/rag-eval js run examples/rag-lab-js/04-export-researchctl-spec.js

# Build the strict NDJSON worker used by researchctl experiment run-rag.
go build -o .bin/rag-lab-worker ./cmd/rag-lab-worker
```

The worker speaks `researchctl-rag-runner-stdio/v1`, accepts one request on stdin, writes protocol frames only to stdout, and sends diagnostics to stderr. It opens the TTC SQLite catalog with WAL-aware `mode=ro` and query-only access. Filters, generated representations, parent-chunk collapse, missing embedders, and missing rerankers fail explicitly before observations.

See [`pkg/ragcontract/README.md`](pkg/ragcontract/README.md) for the dependency-free boundary and `examples/rag-lab-js/04-export-researchctl-spec.js` for pure authoring. Candidate TTC datasets and metrics remain provisional development evidence until human adjudication and holdout freeze.

## Widget IR frontend

The frontend uses a Widget IR (Interchange Representation) — a JSON-compatible data format that describes pages as trees of component nodes. Host applications expose Widget IR at routes such as `/api/widget/pages/{id}`, and the React `WidgetRenderer` renders it in the browser.

### Programming models

| Model | How the frontend is bundled |
| --- | --- |
| **Go embedded** | `pkg/defaultspa.Handler()` serves the prebuilt SPA via `embed.FS`; the host Go app owns the Widget IR API routes |
| **xgoja embedded** | Build the React app, copy `app-dist` into `assets/`, mount with `fs:assets` and `app.spaFromAssetsModule(...)`; the xgoja/Express app owns the Widget IR API routes |
| **Development proxy** | Proxy to the Vite dev server during frontend changes |
| **API-only** | Host app renders Widget IR in its own React shell |

See the provider documentation for full bundling instructions:

- [Widget DSL Getting Started](pkg/xgoja/providers/widgetsite/doc/01-widget-dsl-getting-started.md)
- [Widget DSL v3 Examples](pkg/xgoja/providers/widgetsite/doc/04-widget-dsl-v3-examples.md)
- [Widget DSL v3 API Reference](pkg/xgoja/providers/widgetsite/doc/05-widget-dsl-v3-api-reference.md)
- [Widget DSL JavaScript API Reference and Migration Notes](pkg/xgoja/providers/widgetsite/doc/02-widget-dsl-js-api-reference.md)
- [Bundling the WidgetRenderer SPA](pkg/xgoja/providers/widgetsite/doc/03-widget-dsl-spa-bundling.md)

### Generated-host integration suites

Run `make widgetdsl-sites-smoke` to rebuild the canonical SPA and validate every xgoja reference host. The suite covers the interactive collection/search/pagination/dialog host, the scheduling poll, the workshop CRM, and all committed Widget DSL v3 golden examples through generated binaries and real HTTP routes. Each site also exposes `make smoke` for focused debugging.

### Semantic recipes

New hosts import the single v3 module and compose its domain namespaces:

```js
const widget = require("widget.dsl")

const page = widget.page("Workshop", (p) =>
  p.section("Schedule", (s) =>
    s.view(widget.time.week(events, (w) =>
      w.range(widget.time.range.week("2026-07-14")),
    )),
  ),
)
```

The Goja `rag-widget-site` provider exposes only `widget.dsl`. Context, course, CMS, CRM, scheduling, and time capabilities live under typed namespaces on that module; split helper modules are no longer registered.

## Build

```bash
# Build Go binary with embedded SPA
make build-full

# Build Go binary only
make build-bin

# Run tests
make test

# Run linters
make lint
```

## GoReleaser

Releases are built with GoReleaser:

```bash
make goreleaser
```

The config targets Linux amd64 with CGO enabled (required for `github.com/mattn/go-sqlite3`). Release builds produce:

```
dist/rag-evaluation-system_0.0.0-snapshot_linux_amd64.tar.gz
```

## License

MIT
