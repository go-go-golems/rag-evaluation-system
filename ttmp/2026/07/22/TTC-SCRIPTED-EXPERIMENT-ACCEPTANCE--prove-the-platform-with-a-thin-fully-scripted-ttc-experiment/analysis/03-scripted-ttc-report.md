---
Title: TTC scripted acceptance generated report
Ticket: TTC-SCRIPTED-EXPERIMENT-ACCEPTANCE
Status: complete
Topics:
    - ttc
    - rag
    - evaluation
DocType: analysis
Intent: long-term
Summary: Generated deterministic Researchctl report from the six-run thin TTC acceptance.
---

# TTC scripted acceptance study

- Dataset: `sha256:f158efeb7ec7f2fe21ae2e0bb40fe7c42becaf40188953bdbacaf6fc2260f6cb`
- Analysis source: `sha256:60f0b197758b342476033d1191c5ae08af2bb1cc96ab66aea78b3e15d8cb2464`
- Runs: 6
- Groups: 2

## Completeness and results

- `{"collapse":"chunk","failedOperations":{"coverage":1,"missing":0,"n":3,"unit":"count","value":0},"failedRuns":0,"mrr":{"confidence":0.95,"coverage":1,"high":0,"low":0,"mean":0,"missing":0,"n":3,"sampleUnit":"run","unit":"ratio"},"retries":{"coverage":1,"missing":0,"n":3,"unit":"count","value":0},"runs":3,"runsFailed":0,"runsRequested":3,"runsSucceeded":3}`
- `{"collapse":"unit","failedOperations":{"coverage":1,"missing":0,"n":3,"unit":"count","value":0},"failedRuns":0,"mrr":{"confidence":0.95,"coverage":1,"high":0,"low":0,"mean":0,"missing":0,"n":3,"sampleUnit":"run","unit":"ratio"},"retries":{"coverage":1,"missing":0,"n":3,"unit":"count","value":0},"runs":3,"runsFailed":0,"runsRequested":3,"runsSucceeded":3}`

## Interpretation

This thin TTC acceptance checks that quality and retry-aware execution evidence share immutable run identity. Three replicates support workflow-level completeness checks; this bounded fixture is not a provider-performance recommendation.

## Reproduction

Regenerate with the exact checked-in analysis source and immutable laboratory database recorded by the invoking command. Missing measurements are not coerced to zero; confidence intervals use runs as the sampling unit.
