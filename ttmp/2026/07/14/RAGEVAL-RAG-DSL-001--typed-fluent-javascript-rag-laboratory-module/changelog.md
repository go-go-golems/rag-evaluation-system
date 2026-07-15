# Changelog

## 2026-07-14

- Initial workspace created

## 2026-07-14

Defined the proposed require("rag") public contract, fluent builders, canonical immutable-specification boundary, and an intern-oriented implementation plan before application code is written.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/14/RAGEVAL-RAG-DSL-001--typed-fluent-javascript-rag-laboratory-module/design-doc/01-typed-fluent-rag-module-design-and-implementation-guide.md — Implementation guide and decision records
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/14/RAGEVAL-RAG-DSL-001--typed-fluent-javascript-rag-laboratory-module/reference/01-rag-laboratory-javascript-module-api-specification.md — Normative API contract

## 2026-07-14

Step 3: Extracted the shared immutable experiment specification schema and fingerprint contract (commit 95c153e).

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/internal/experimentspec/specification.go — New schema authority
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/internal/services/experimentrun/service.go — Uses schema authority

## 2026-07-14

Step 4: Added the pure typed RAG experiment builder, deterministic validation, and fingerprint tests (commit 31a3c93).

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/pkg/raglab/builder.go — Fluent typed builder
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/pkg/raglab/builder_test.go — Regression coverage

## 2026-07-14

Step 5: Added immutable artifact-catalog lineage validation and evaluation/representation catalog schema (commit 3b6dc55).

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/internal/db/db.go — Immutable catalog schema
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/pkg/raglab/catalog.go — Compatibility rules

## 2026-07-15

Implemented native JavaScript RAG module, explicit laboratory persistence/start, runtime tests, and TypeScript descriptor (commit c46485e).

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/pkg/gojamodules/rag/module.go — Public require(rag) API
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/pkg/raglab/laboratory.go — Append-only persistence/start boundary

