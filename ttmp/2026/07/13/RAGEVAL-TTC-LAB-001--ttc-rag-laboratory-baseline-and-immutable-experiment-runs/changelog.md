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
