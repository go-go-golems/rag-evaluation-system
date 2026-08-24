---
Title: RAG legacy deferred hard-cut closure audit
Ticket: RAG-EVAL-LEGACY-CLEANUP
Status: complete
Topics:
    - rag
    - workflow
    - workflow
DocType: analysis
Intent: long-term
Summary: Evidence that every deferred RAG execution and intake lifecycle now has a canonical replacement and has been deleted.
---

# RAG legacy deferred hard-cut closure audit

## Result

Every deferred production lifecycle named by the original disposition has crossed its replacement gate and is deleted. `pkg/ragengine` remains only as RAG-owned semantics used by task implementations/tests; it is not a scheduler or command path.

| Deferred path | Replacement evidence | Deletion evidence |
|---|---|---|
| direct `rag-worker` | RAG lowering, Geppetto operations, generic Scraper runner, Researchctl plans | `cmd/rag-worker` absent; study guard rejects its return |
| `internal/preparationworkflow` | `rag-v2-provider-free@1.0.0` and `rag-v2-geppetto@1.0.0` task packages | directory absent; full tests pass |
| old intake workflow | `pkg/ragintakeworkflow`, Workflow V3 CLI/API/frontend, restart/provider/privacy acceptance | `internal/workflow` and old Scraper imports absent |
| RAG study run loop | immutable study bundle + pure Researchctl plan | study exposes only validate/explain/compile |
| adapter lifecycle | Researchctl generic process runner and Scraper runner | `run.go`/`progress.go` absent |
| TTC bespoke lifecycle | workload-only `experiments/ttc-scripted` | sweep/internal package absent; exact inventory guard passes |

## Fresh evidence

- Full RAG tests, race tests, lint/vet, build, typecheck/build, module tidiness, generation, intake smoke, study smoke, TTC six-run acceptance, and deletion guards passed in the predecessor tickets.
- `rg` finds no active old Scraper engine import, direct worker command, preparation package, or execution adapter.
- `docmgr doctor` passes after replacing historical related-file pointers with canonical current files.

No compatibility aliases, dormant copies, TODO placeholders, or duplicate lifecycle owners remain.
