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

## 2026-07-14

Executed and scored the first 20-card immutable TTC retrieval comparison; added trace and candidate-score scripts plus isolated standalone ticket scripts from package tests.

## 2026-07-14

Added append-only experiment specifications, run events, query traces, terminal summaries, and offline lifecycle coverage; corrected RRF score initialization nondeterminism.

## 2026-07-14

Exposed immutable experiment laboratory APIs and React trace-inspection UI; imported the real 20-card baseline run; added human adjudication packet and operator guide.

## 2026-07-16

Reconciled reranker expansion drafts into TTC v2 policy and added twelve-card direct-source adjudication batch; v2 remains unfrozen pending remaining family labels.

### Related Files

- ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/04-ttc-v2-support-faq-adjudication-batch-01.md — Concrete direct-source evidence packet

## 2026-07-16

Added draft 240-card TTC v2 partition and evidence-family leakage-audit protocol; labels remain unfrozen.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/01-implementation-diary.md — Step 18 records the design decision and review guidance
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/06-ttc-v2-240-card-partition-and-leakage-audit-protocol.md — Exact 144/48/48 allocation, abstention distribution, and validator contract


## 2026-07-16

Added 70-card expansion draft and 50-card SQLite-grounded audit draft. These 120 candidates are explicitly unlabelled; source validation and evidence-family adjudication remain required before any freeze.

### Related Files

- ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/04-ttc-evaluation-expansion-v0-70-proposed-cards.md — 70-card structured authoring queue
- ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/06-ttc-expansion-audit-and-50-card-source-grounded-draft.md — 50 source-grounded candidate cards and coverage audit


## 2026-07-16

Validated the 70-card and 50-card expansion drafts against the rebuilt TTC SQLite export: 142 unique source IDs, 142 resolved, none missing. This is identity validation only; labels and evidence spans remain open.

### Related Files

- ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/07-validate-expansion-source-ids.py — Read-only source-ID validator


## 2026-07-16

Added the second 80-card expansion batch and upgraded the draft audit utility to parse both YAML-list formats, including explicit unanswerable controls. The expansion queue now contains 200 records beyond the pilot; source validation resolves 173 unique IDs.

### Related Files

- ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/05-ttc-evaluation-expansion-y-v0-80-proposed-cards.md — Second expansion batch


## 2026-07-16

Registered candidate:ttc-expansion-v0 with 148 snapshot-compatible cards; excluded two policy-conflict cards whose source is intentionally absent from the current snapshot. Emitted and committed the canonical manifest.

### Related Files

- ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/09-register-ttc-expansion-candidate-dataset.py — Immutable candidate registrar


## 2026-07-16

Generalized immutable retrieval trace driver to parse multiple expansion card files and inline YAML card records. Development run command is ready; execution awaits Ollama endpoint at 127.0.0.1:11435.

### Related Files

- ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/04-run-immutable-retrieval-traces.go — Multi-format expansion trace runner


## 2026-07-16

Ran the 148-card expansion development trace through Mac Ollama: vector MRR 0.9174, hybrid MRR 0.9005, BM25 MRR 0.8221; hybrid/vector Recall@10 0.9722. Mean latency 173 ms, P95 230 ms. Metrics remain provisional because source IDs are not human-adjudicated.

### Related Files

- ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/09-ttc-expansion-development-run-results.md — Development run result record

