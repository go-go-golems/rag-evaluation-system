---
Title: "Durable Geppetto Workflow operations"
Slug: rag-geppetto-workflow-operations
Short: "Run provider-backed RAG v2 pipelines with Workflow V3 operation custody."
Topics:
- rag
- workflow-v3
- geppetto
- providers
Commands:
- rag-workflow-runner
Flags:
- provider-config
- provider-fixture
- max-provider-operations-per-attempt
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: Application
---

Provider-backed RAG uses the exact `rag-v2-geppetto@1.0.0` task package. RAG-eval owns request and response semantics; Geppetto owns provider APIs; Scraper Workflow V3 owns admission, leases, attempts, budgets, durable external-operation records, cancellation, and observation export. The canonical RAG v2 execution remains the scientific input.

## Provider attachment points

The package executes these exact operator versions:

- `representations.structured-summary/v1`
- `representations.synthetic-questions/v1`
- `representations.combined-summary-questions/v1`
- `embed.model/v1` when the model is not the fixture embedding model
- `rerank.cross-encoder/v1`
- `generate.answer/v1`

Unknown operators, versions, manifests, schemas, provider profiles, and authority digests fail closed. Provider-free plans continue to use `rag-v2-provider-free@1.0.0`; the two catalogs are not interchangeable.

## Start a provider runner

```bash
rag-workflow-runner \
  --provider-config experiments/real-provider-v2/provider-config.yaml \
  --state-root state/rag-provider-workflows \
  --artifact-root state/rag-provider-artifacts \
  --lease-duration 30s \
  --max-provider-operations-per-attempt 10000
```

`--provider-config` is host-only configuration. Plans contain canonical profile, model, prompt, settings, and manifest digests, never credentials. The runner resolves Geppetto profiles and environment references locally and rejects a plan whose embedded authority differs from the host authority.

For deterministic acceptance only, use `--provider-fixture`. `--provider-fixture-delay` deliberately delays each fixture contact so crash-window tests can kill a process after durable admission. Never use fixture mode for scientific provider claims.

## Generate deterministic fixtures

```bash
rag-workflow-provider-fixture --output /tmp/rag-workflow-provider
```

The generator writes two cases, canonical domain configs, set-input archives, and direct `ragengine` parity files. Regeneration is byte-stable. The checked-in fixtures live under `examples/rag-workflow-provider/`.

## Custody and replay semantics

Every actual provider contact has one `provider.generate/v1`, `provider.embed/v1`, or `provider.rerank/v1` record. Admission occurs immediately before contact. Completion records contain only bounded counters, safe failure class/code, timing, authority digest, and a request correlation digest. Raw prompts, evidence text, responses, URLs, headers, credentials, and provider bodies are excluded.

A Workflow retry creates a distinct operation. Correlation digests remain stable for the same semantic request, while operation IDs and attempt numbers preserve physical-contact identity. Cache hits create no provider operation and report zero request usage. A crash after admission can leave an incomplete operation; exported manifests preserve that uncertainty rather than claiming success or failure.

Provider request authority is reserved before work. The default Workflow budget dimension is `requests`; provider token and cost counters remain authoritative operation evidence. Provider model manifests also bound response tokens and pricing behavior. Cancellation and timeout close observable contacts using detached, bounded completion persistence. If completion persistence itself fails, the task does not publish successful domain output.

## Authorized real-provider acceptance

The repository includes an opt-in acceptance test. It requires explicit authorization and reachable embedding/reranking endpoints:

```bash
RAG_WORKFLOW_REAL_PROVIDER_ACCEPT=1 \
RAG_EMBEDDING_BASE_URL=http://127.0.0.1:11434 \
RAG_RERANKER_BASE_URL=http://127.0.0.1:18012 \
ttmp/2026/07/22/RAG-GEPPETTO-WORKFLOW-OPERATIONS--execute-geppetto-generation-embedding-and-reranking-as-durable-workflow-operations/scripts/02-authorized-real-provider-acceptance.sh
```

The test uses a fresh provider cache, a two-second Workflow lease, real Geppetto generation, Ollama embeddings, and a llama.cpp cross-encoder. It proves lease renewal, all three operation kinds, usage/cost capture, grounded answer citations, and one-attempt completion.

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| `RAG_WORKFLOW_PROVIDER_AUTHORITY` | Host provider identity differs or is malformed | Regenerate the plan with the same canonical host manifests and profile settings. |
| `RAG_PROVIDER_OPERATION_ADMISSION` | Request authority or budget is exhausted | Increase an explicitly reviewed limit or reduce batching/case scope; no provider contact occurred. |
| `PROVIDER_TIMEOUT` or `PROVIDER_TRANSPORT` | Geppetto contact failed | Inspect safe operation class/code and retry evidence; raw provider bodies are intentionally absent. |
| `BUDGET_USAGE_INVALID` | Task and operation accounting disagree | Treat this as a custody defect; do not bypass settlement checks. |
| Repeated `lease_lost` attempts | Worker lacks lease renewal support or uses an old Scraper build | Use the pinned Scraper version in `go.mod` or newer. |
| Real acceptance is skipped | Explicit authorization is absent | Set `RAG_WORKFLOW_REAL_PROVIDER_ACCEPT=1` only after reviewing provider cost and endpoint configuration. |

## See Also

- `rag-workflow-v3-lowering`
- `rag-real-provider-operator-playbook`
- `rag-v2-api-reference`
