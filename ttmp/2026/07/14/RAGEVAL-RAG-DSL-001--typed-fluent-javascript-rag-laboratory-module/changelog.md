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


## 2026-07-15

Migrated rag-eval-js to xgoja/v2, packaged the RAG provider, generated declarations, and added runnable RAG scripts (commit 7b74539).

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/examples/rag-lab-js/README.md — Operator example instructions
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/pkg/xgoja/providers/rag/provider.go — Generated runtime provider


## 2026-07-15

Implemented raw lexical/vector/weighted-RRF execution primitives, immutable card loading, and ran a 20-card raw BM25 TTC observation (commits 2106b99, 99b2476).

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/pkg/raglab/executor.go — Append-only execution adapter
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/14/RAGEVAL-RAG-DSL-001--typed-fluent-javascript-rag-laboratory-module/scripts/02-run-ttc-raw-bm25-experiment.go — Executed TTC observation


## 2026-07-15

Validated all Go packages plus xgoja v2 plan, declaration generation, generated binary, and RAG plan-only example after executor work.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/xgoja.yaml — Final generated runtime validation
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/pkg/raglab/executor_test.go — Executor regression coverage


## 2026-07-15

Verified the reachable Mac-hosted Ollama route and added a private SSH-loopback tunnel operator playbook; documented the unresolved mimimi.local alias, live mimimi-2.local model inventory, tmux lifecycle, Geppetto configuration, and vector/RRF preflight.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/14/RAGEVAL-RAG-DSL-001--typed-fluent-javascript-rag-laboratory-module/reference/02-implementation-diary.md — Step 11 operational investigation diary
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/14/RAGEVAL-RAG-DSL-001--typed-fluent-javascript-rag-laboratory-module/reference/03-mimimi-ollama-tunnel-operator-playbook.md — Reusable verified tunnel procedure


## 2026-07-15

Step 12: added and executed fresh explicit-provider TTC vector-only and semantic-weighted RRF runs (commit 5fb76c0); both runs succeeded with append-only traces and recorded candidate quality, latency, cost scope, and shared-artifact storage evidence.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/14/RAGEVAL-RAG-DSL-001--typed-fluent-javascript-rag-laboratory-module/reference/02-implementation-diary.md — Step 12 durable results and review instructions
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/14/RAGEVAL-RAG-DSL-001--typed-fluent-javascript-rag-laboratory-module/scripts/03-run-ttc-vector-and-weighted-rrf-experiments.go — Executable observation definition


## 2026-07-15

Step 13: exposed synchronous lab.execute() to JavaScript with an explicit finite-number queryEmbed callback and immutable evaluation-card loading (commit cb04190).

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/pkg/gojamodules/rag/module.go — JavaScript execution and query embedding adapter
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/pkg/raglab/laboratory.go — Atomic public run orchestration


## 2026-07-15

Completed operator documentation and immutable-specification inspection for
the JavaScript RAG laboratory (commit 4cd09ff).

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/doc/01-rag-laboratory-javascript.md — Embedded operator tutorial
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/examples/rag-lab-js/03-execute-with-geppetto.js — Geppetto-backed execution example
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/internal/services/experimentrun/service.go — Canonical manifest API response


## 2026-07-15

Completed the generated runtime, full repository, web production-build, and
live Geppetto/Ollama JavaScript execution validation (commit be3623a).

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/14/RAGEVAL-RAG-DSL-001--typed-fluent-javascript-rag-laboratory-module/scripts/05-run-rag-lab-js-geppetto.js — Frozen TTC live JS observation
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/xgoja.yaml — Generated runtime validation plan
