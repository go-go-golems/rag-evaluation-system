---
Title: Authorized real Umans qualification results
Ticket: RAG-TTC-V3-SWEEP
Status: complete
DocType: analysis
---

# Authorized real Umans qualification results

## Result

The authorized Workflow V3 qualification completed all eight planned cells: 16 fixed TTC chunks, batch sizes 1/2/4/8, and generation concurrency 1/2. It planned 60 generation requests and recorded exactly 60 successful generation attempts plus 128 embedding requests. The cumulative authority ledger ended at 125 admissions: 61 prior admissions plus 64 new admissions, remaining below the explicitly authorized maximum of 129.

This is a bounded qualification, not a broad benchmark claim. It establishes that the real profile, Mac-backed embedding/reranking topology, durable external-operation ledger, per-cell custody, privacy boundary, and downstream researchctl import work together under explicit authority.

## Observed cell makespans

| Batch size | Concurrency | Makespan |
| --- | ---: | ---: |
| 1 | 1 | 208.607 s |
| 2 | 1 | 157.050 s |
| 4 | 1 | 92.202 s |
| 8 | 1 | 68.864 s |
| 8 | 2 | 34.389 s |
| 4 | 2 | 62.029 s |
| 2 | 2 | 77.990 s |
| 1 | 2 | 174.067 s |

Within this fixed sample, larger generation batches reduced elapsed cell time. Increasing concurrency from one to two also reduced elapsed time for every matched batch size. These are descriptive observations only: server scheduling, limited sample size, and missing provider usage fields prevent cost or token-efficiency conclusions.

## Accounting and evidence

The provider returned bounded nonzero generation usage only for two cells; the remaining cells retain explicit zero/unavailable fields rather than inferred values. The evidence therefore supports request and time analysis, but not a complete observed token/cost analysis. The authority ceiling remains a conservative authorization boundary, not a statement of actual billed usage.

Each cell exported one JSONL operation ledger and one manifest before source-bearing runtime deletion. The final output contains eight cell checkpoints, eight JSONL ledgers, eight manifests, aggregate evidence, CSV, measurement JSONL, and a durable generation-authority record. Retained runtime SQLite/WAL content was absent after completion.

A forbidden-string scan of compact retained evidence found no source canary, provider-body field, authorization/bearer material, URL, API-key pattern, private-key pattern, or vector field. A fresh researchctl import verified 25 artifacts and four scalar operation metrics: eight cells, 60 generation requests, 128 embedding requests, and 875,199,023 aggregate cell-makespan microseconds.

## Graph review

The real graph bundle is under `sources/real-attempt-003/graphs/`. It uses the title prefix `Workflow V3 real Umans qualification`. Initial visual review found makespan publication-readable and identified that the timeline had selected a nonexistent concurrency-4 series. The renderer now selects the highest observed concurrency level, was re-rendered for the real 1/2 matrix, and preserves the generation hard-cap reference line.

## Limits and next decision

No second matrix is authorized by this qualification. Before any broader performance statement or replicate, review the real graphs, investigate why the provider omitted usage fields for most cells, decide whether a larger sample is scientifically justified, and obtain a separate explicit authority envelope.

## References

- `sources/real-attempt-003/evidence.json`
- `sources/real-attempt-003/generation-authority.json`
- `sources/real-attempt-003/operations/`
- `sources/real-attempt-003/researchctl-run-export.json`
- `sources/real-attempt-003/graphs/manifest.json`
