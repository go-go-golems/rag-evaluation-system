---
Title: "Run a one-query RAG preview"
Slug: "rag-preview-workflow"
Short: "Select one study cell and execute one query in a scratch or existing laboratory."
Topics:
- rag
- preview
- traces
- researchctl
Commands:
- rag-eval preview
Flags:
- query
- variant
- factor
- project
- inputs
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Application
---

A preview reuses normal study compilation and worker execution for one query. It does not add a JavaScript `.run()` API or a second lifecycle: `rag-eval` creates a one-query candidate dataset, selects one expanded cell, and submits the same generic specification used by study runs.

## Run in a scratch laboratory

Omit `--project` to create and initialize a new scratch project. The command reports its path so the run, verified artifacts, and trace remain inspectable after completion.

```bash
rag-eval preview examples/rag-v2/06-raw-study.js \
  --inputs examples/rag-v2/inputs-ttc-catalog.json \
  --ttc-database data/rag-eval.db \
  --query "Why did reciprocal rank fusion select this source?" \
  --variant raw \
  --researchctl-command researchctl \
  --worker-command rag-worker
```

The generated dataset is explicitly `candidate` with split `preview`. It is diagnostic evidence, not an adjudicated benchmark.

## Select a factorized cell

Use one `--factor key=value` flag per selection. Omitted variant/factors choose the first stable cell in compiler order.

```bash
rag-eval preview study.js \
  --inputs inputs.json \
  --query "What changed?" \
  --variant all \
  --factor collapse=unit
```

## Use an existing project

Pass an existing initialized project and experiment when preview evidence belongs with ongoing work.

```bash
rag-eval preview study.js \
  --inputs inputs.json \
  --query "What changed?" \
  --project project.yaml \
  --experiment-id EXP-RAG-PREVIEW
```

Researchctl still owns retries, timeout, artifact verification, and terminal state. A failed preview preserves already-recorded events and artifacts in its failed attempt.

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| `RAG_PREVIEW_CELL_NOT_FOUND` | Variant or factor values do not identify an expanded cell. | Run `rag-eval study explain` and copy the exact IDs. |
| Scratch initialization fails | The researchctl executable is missing or incompatible. | Pass the current binary with `--researchctl-command`. |
| Preview succeeds but relevance is zero | Preview queries have no adjudicated relevance labels by default. | Treat the trace as diagnostic; use a study dataset for quality claims. |
| Worker exits on cancellation | The attempt timeout or caller context cancelled execution. | Inspect the preserved partial attempt, then increase `--timeout` only when justified. |

## See Also

- `rag-eval help rag-study-workflow`
- `rag-eval study explain --help`
- `researchctl experiment execute-spec --help`
