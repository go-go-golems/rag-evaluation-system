# Real-provider RAG v2 candidate

This directory is the explicit candidate-only real-provider profile for RESEARCHCTL-015.

- `study.js` compiles the candidate study with `fixtureProviders: false`.
- `product.js` compiles the corresponding product plan.
- `preview.js` is a bounded one-large-chunk diagnostic study; it exercises the full retrieval trace without sending thousands of generation requests.
- `manifests/` and `schemas/` are public, immutable provider identity inputs.
- `provider-config.example.yaml` is operational-only; copy it outside this directory and set the endpoint environment references before use.

The placeholder corpus and evaluation envelope digests in `inputs.json`, `study.js`, and `product.js` must be replaced with the custody-verified candidate artifacts before a run. This is not a benchmark claim.

## Durable preparation smoke

`cmd/rag-preparation-smoke` is the narrow operator check for the scraper-backed
combined-preparation workflow. It sends exactly one small chunk to the configured
real generator, persists its scraper workflow in the supplied SQLite database,
and reports only workflow counts and IDs. It does not index a corpus or make a
benchmark claim.

Provide a host-only provider configuration with the normal model, prompt, schema,
cache, and profile settings; do not commit credentials or endpoint secrets. For
example:

```sh
GOWORK=off go run ./cmd/rag-preparation-smoke \
  --provider-config /secure/path/providers.yaml \
  --state-db /secure/path/rag-preparation-smoke.sqlite
```

A successful run has one `rag-preparation/combined-batch/v1` operation and one
local finalizer in `succeeded` state. Re-running the same input attaches to the
same immutable workflow identity rather than issuing another provider request.
