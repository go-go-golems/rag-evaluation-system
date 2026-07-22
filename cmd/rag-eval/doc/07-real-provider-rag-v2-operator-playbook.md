---
Title: "Operate a real-provider RAG v2 run"
Slug: "rag-real-provider-operator-playbook"
Short: "Configure, preflight, preview, and execute a real-provider RAG v2 study without fixture fallback."
Topics:
- rag
- providers
- researchctl
- experiments
Commands:
- rag-eval study validate
- rag-eval preview
- rag-eval study run
- rag-worker
Flags:
- provider-profile
- provider-config
- worker-arg
- timeout
- ttc-database
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Application
---

A real-provider run is a custody-sensitive experiment, not a provider smoke test. This playbook keeps provider construction in the host, keeps credentials and endpoints out of JavaScript, and makes the fixture profile impossible to select accidentally.

## 1. Prepare the host

Build the current worker and researchctl binaries from their repositories. Keep the provider host configuration outside the RAG repository because it contains operational paths and environment references.

```bash
cd /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system
GOWORK=off go build -o /tmp/rag-worker ./cmd/rag-worker

cd /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl
GOWORK=off go build -o /tmp/researchctl ./cmd/researchctl
```

Copy `experiments/real-provider-v2/provider-config.example.yaml` to a temporary host directory. Set endpoint variables in the shell or service manager; never put credentials or endpoint values in the JavaScript study, manifest, or committed config.

```bash
export RAG_EMBEDDING_BASE_URL=http://127.0.0.1:11434
export RAG_GENERATOR_BASE_URL=http://127.0.0.1:11434/v1
export RAG_RERANKER_BASE_URL=http://127.0.0.1:18012
```

Prefer the Mac-hosted generator and embedder when reachable. The selected reranker is a llama.cpp `/v1/rerank` service. A typical SSH tunnel is:

```bash
ssh -fN -L 18012:127.0.0.1:8012 user@mac-host
```

Verify the service before constructing a real ProviderSet. Do not substitute fixtures when a real endpoint is unavailable.

## 2. Validate immutable inputs

Use the TTC database resolver for catalog bindings. Do not replace catalog references with ad-hoc text for a candidate claim.

```bash
GOWORK=off go run ./cmd/rag-eval study validate \
  experiments/real-provider-v2/study.js \
  --inputs experiments/real-provider-v2/inputs.json \
  --ttc-database data/rag-eval.db
```

The study must be candidate-labeled and must retain `fixtureProviders: false`. Placeholder digests in an example are not custody evidence; the run must bind resolved envelopes with verified SHA-256 identities.

## 3. Preflight the real worker

The worker requires an explicit profile. `fixtures` is only for tests. `real` requires the host configuration and performs an in-process execution-requirement check before provider requests.

```bash
/tmp/rag-worker \
  --capabilities \
  --provider-profile real \
  --provider-config /tmp/researchctl-015-real-provider-host/providers.yaml
```

The capability response may contain public profile and manifest identity, but must not contain endpoint URLs, credentials, or provider response bodies. A missing generator, embedder, reranker, schema validator, or persistent cache is a preflight failure.

## 4. Run a bounded preview first

A preview exercises the compiler, worker, provider adapters, traces, and researchctl custody with one query. It is diagnostic evidence and never a benchmark result.

```bash
GOWORK=off go run ./cmd/rag-eval preview \
  experiments/real-provider-v2/preview.js \
  --inputs /tmp/verified-preview-inputs.json \
  --query 'What information is present in the source?' \
  --researchctl-command /tmp/researchctl \
  --worker-command /tmp/rag-worker \
  --worker-arg --provider-profile --worker-arg real \
  --worker-arg --provider-config \
  --worker-arg /tmp/researchctl-015-real-provider-host/providers.yaml \
  --timeout 15m
```

Inspect every query-trace operator, channel, collapse, hydration result, usage value, provider/model identity, and cost map. Unknown local cost must be absent, not numeric zero.

## 5. Execute a candidate study

Initialize the project laboratory before a direct study run. The laboratory's artifact root must match the root used to stage input envelopes.

```bash
/tmp/researchctl lab init \
  --project /path/to/project.yaml \
  --database /tmp/rag-candidate.db

GOWORK=off go run ./cmd/rag-eval study run \
  experiments/real-provider-v2/study.js \
  --project /path/to/project.yaml \
  --experiment-id EXP-RAG-V2-REAL \
  --inputs /path/to/custody-verified-inputs.json \
  --artifact-root /tmp/artifacts \
  --database /tmp/rag-candidate.db \
  --ttc-database data/rag-eval.db \
  --researchctl-command /tmp/researchctl \
  --worker-command /tmp/rag-worker \
  --worker-arg --provider-profile --worker-arg real \
  --worker-arg --provider-config \
  --worker-arg /tmp/researchctl-015-real-provider-host/providers.yaml \
  --timeout 30m \
  --max-attempts 1 \
  --spec-output-dir /tmp/rag-candidate-specs
```

Record the run ID, attempt ID, specification ID, input digests, artifact count, trace count, metrics, usage, and terminal status. Recompile with identical authoring and inputs and compare specification bytes before interpreting metrics.

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| `provider-profile is required` | The worker no longer has an implicit fixture default. | Pass `--provider-profile real` or explicitly use `fixtures` only in a test. |
| `RAG_INPUT_EVALUATION_POLICY` | The evaluation envelope is not candidate-labeled. | Use a custody-verified candidate envelope with `split: candidate` and `status: candidate`. |
| `RAG_MODEL_MANIFEST_MISSING` after zero units | The unit operator does not match the corpus record shape. | Use the corpus-compatible unit operator and inspect the trace's unit count. |
| `RAG_OUTPUT_SCHEMA_MISSING` | A generation node did not resolve its prompt output schema. | Check the prompt manifest and the canonical generation request. |
| `lstat .../artifacts: no such file or directory` | Input staging and the laboratory's resolved artifact root differ. | Initialize the laboratory and stage inputs under the same artifact root. |
| Reranker connection refused | The Mac-hosted llama.cpp service or SSH tunnel is down. | Start the service/tunnel; do not switch to fixtures or claim P4.3 evidence. |
| Quality values are all zero in a preview | Preview queries have no adjudicated relevance IDs. | Treat them as plumbing evidence, not ranking quality. |

## See Also

- `rag-study-workflow` — compile and execute studies.
- `rag-preview-workflow` — run one diagnostic query.
- `rag-product-runtime` — operate the product host separately from researchctl.
- `rag-v2-cutover` — understand retired fixture and lifecycle paths.
