# Changelog

## 2026-05-28

- Initial workspace created


## 2026-05-28

Created end-to-end retrieval foundation design package: BM25 first, query-vector search next, hybrid retrieval after standalone validation, and smoke queries before benchmarks.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/28/RAGEVAL-004--end-to-end-search-retrieval-foundation/design-doc/01-end-to-end-search-retrieval-implementation-guide.md — Primary implementation guide
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/28/RAGEVAL-004--end-to-end-search-retrieval-foundation/reference/01-implementation-diary.md — Design workflow diary


## 2026-05-28

Uploaded RAGEVAL-004 retrieval foundation design bundle to reMarkable at /ai/2026/05/28/RAGEVAL-004.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/28/RAGEVAL-004--end-to-end-search-retrieval-foundation/design-doc/01-end-to-end-search-retrieval-implementation-guide.md — Uploaded as part of reMarkable bundle
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/05/28/RAGEVAL-004--end-to-end-search-retrieval-foundation/reference/01-implementation-diary.md — Uploaded as part of reMarkable bundle


## 2026-05-28

Implemented Phase 1 BM25 search service and CLI, built a TTC guides/articles index, and ran real lexical queries to validate retrieval plumbing.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/cmd/rag-eval/cmds/search/index.go — Index build CLI
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/cmd/rag-eval/cmds/search/query.go — Query CLI
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/search/bm25.go — BM25 build/query implementation


## 2026-05-28

Added BM25 HTTP endpoints for search index builds and lexical query execution, validated with a live curl smoke query.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/api/handlers.go — Search HTTP route registration and handlers


## 2026-05-28

Implemented query-vector search over stored embeddings with CLI/API adapters, fake-provider tests, source filtering, and a bounded live OpenAI smoke query.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/cmd/rag-eval/cmds/search/vector.go — Vector search CLI
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/db/search_queries.go — Embedding-with-context query helper
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/internal/services/search/vector.go — Vector retrieval service

