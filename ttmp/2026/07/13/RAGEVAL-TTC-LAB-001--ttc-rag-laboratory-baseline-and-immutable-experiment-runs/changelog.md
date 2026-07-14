# Changelog

## 2026-07-13

- Initial workspace created


## 2026-07-14

Added the source-first evaluation dataset authoring/adjudication protocol, named relevance levels, 20 source-validated TTC candidate cards, and a read-only validator; withheld conflicting cancellation-policy evidence pending owner adjudication.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/design-doc/02-evaluation-dataset-authoring-and-adjudication-protocol.md — Defines fixed-truth authoring and immutable versioning
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/02-ttc-baseline-evaluation-dataset-v1-candidate-cards.md — Records candidate cards and provisional labels
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/01-validate-ttc-baseline-evaluation-cards.sh — Validates source grounding before adjudication

## 2026-07-14

Completed shi3: added and validated the deterministic TTC baseline importer, 200-document manifest, source-card seed inclusion, and Glazed corpus import command (commit 2fcf2bc).

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/cmds/corpus/import_ttc.go — Operator-facing command
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/internal/services/ttcimport/service.go — Deterministic importer implementation

## 2026-07-14

Completed 3ydv: added append-only TTC source artifacts, document revisions, ordered corpus snapshots, and the snapshot-ttc command (commit c846043).

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/internal/db/db.go — Immutable corpus schema
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/internal/services/corpussnapshot/service.go — Immutable corpus persistence

## 2026-07-14

Completed 26xz and rggc: added canonical artifact fingerprints, exact source-range chunking, immutable chunk plans, chunk sets, and persisted immutable chunks (commits 0f5a4a0, ecd8f2a, 425412e).

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/internal/services/immutablechunk/service.go — Chunk artifact builder

## 2026-07-14

Recorded Geppetto/Ollama embedding investigation: single 768D provider calls succeed, while the batch immutable-set path reproduces a silent no-artifact exit; added ticket probe script.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/02-geppetto-ollama-embedding-probe.go — Reproduction probe

## 2026-07-14

Measured real TTC payload behavior: CPU-only Ollama takes over two minutes for one 1,200-character chunk; added --request-timeout-seconds to surface bounded live-run failures.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/cmds/embedding/build_immutable.go — Command-level provider deadline
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/02-geppetto-ollama-embedding-probe.go — Real payload probe

## 2026-07-14

Audited and completed the retroactive embedding diary: corrected foreground-wrapper interpretation, recorded payload/runtime evidence, and documented timeout/long-job operation.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/01-implementation-diary.md — Step 10 operational diary audit

## 2026-07-14

Completed lbwm: built and verified the real Mac-backed Ollama 768D immutable embedding set (2,024 vectors; set sha256:2665c524...e03e0) alongside offline deterministic coverage.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/internal/services/immutableembedding/service.go — Immutable embedding artifact persistence
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/01-implementation-diary.md — Mac-backed execution evidence
