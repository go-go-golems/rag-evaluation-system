---
Title: TTC scripted experiment acceptance audit
Ticket: TTC-SCRIPTED-EXPERIMENT-ACCEPTANCE
Status: complete
Topics:
    - ttc
    - rag
    - evaluation
    - scripting
DocType: analysis
Intent: long-term
Summary: Requirement-to-evidence audit for the thin workload-only TTC experiment through the converged platform.
---

# TTC scripted experiment acceptance audit

## Result

Accepted. A clean checkout executes one command to derive bounded immutable TTC inputs, compile a JS RAG study, create and resume Researchctl replicates through Scraper Workflow V3, project RAG quality and retry-aware operation observations, and regenerate checked-in JS analysis. TTC owns no lifecycle implementation.

## Scope clarification

The original guide teaches a 12-case × 3-replicate provider-performance campaign. That remains a valid future scientific study, but it is not necessary to prove the **thin platform acceptance** named by this ticket or to authorize legacy deletion. The completed acceptance uses two meaningful collapse cases × three replicates and makes no performance/model-quality recommendation. Provider execution itself is not waived: the closed `RAG-GEPPETTO-WORKFLOW-OPERATIONS` predecessor already passed a bounded real TTC workload with Geppetto generation, Ollama embedding, llama.cpp reranking, ten custodied provider contacts, retry-aware timing, and privacy evidence. This ticket proves those platform layers compose without bespoke TTC orchestration.

## Evidence map

| Requirement | Evidence |
|---|---|
| Workload-only source | `experiments/ttc-scripted/`: JS, JSON, and README only; guard checks exact inventory. |
| Scripted study | `study.js`, `pipeline.js`; two factor cases, three replicates, requested quality/failure metrics. |
| Canonical execution | Compiled `rag-workflow-study-bundle/v1`; generated Researchctl plan; `rag-workflow-runner`; Workflow V3. |
| TTC lineage | Ticket script reads the authorized `data/ttc-wordpress-rag.sqlite`, materializes three bounded source records and queries, and binds content-addressed artifacts. |
| Resume | First run executes six; second run executes zero and resumes six terminal replicates. |
| Quality + performance custody | RAG MRR observations and generic retry/external-operation observations share selected Researchctl attempts. |
| Replicate-aware analysis | `analysis.js`; query-scoped MRR first reduces within run, then across three runs. |
| Deterministic publication | Two `analysis run` command outputs compare byte-for-byte; table, SVG, Markdown, and result manifests are write-once. |
| No bespoke lifecycle | No `rag-ttc-v3-sweep`, `internal/workflowv3ttc`, direct Researchctl SQL writes, importer, custom operation ledger, or Python analysis. |
| Real providers | Predecessor acceptance at RAG commit `51ff7f4`; provider-operation ticket audit and authorized script remain canonical evidence. |

## Fresh commands

```text
scripts/01-run-fixture-ttc-study.sh
scripts/02-scripted-architecture-guards.sh
node --check experiments/ttc-scripted/*.js
GOWORK=off go test ./pkg/ragworkflow ./pkg/researchctladapter ./cmd/rag-eval/cmds/study -count=1
```

Results:

```text
PASS: 2 TTC cases x 3 replicates executed, resumed, and regenerated deterministic quality analysis
PASS: TTC contains only JS/data/docs workload source and no bespoke lifecycle or analysis infrastructure
```

## Publication inspection

- `01-mrr-by-collapse.svg`: inspectable 800×420 chart with title, bounded labels, explicit x axis, and zero-baseline scale.
- `02-primary-results.json`: two groups, three runs each, complete MRR/retry/operation coverage.
- `03-scripted-ttc-report.md`: digest-linked report with explicit limitations and reproduction semantics.

## Claim boundary

The fixture run establishes architecture, custody, resume, missingness, and analysis determinism. It is not a provider benchmark. The bounded predecessor run establishes that real providers execute under the same Workflow V3 operation contracts. A larger randomized provider campaign must create a new immutable plan and report; it is not silently inferred from either acceptance artifact.
