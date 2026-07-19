# Real-provider RAG v2 candidate

This directory is the explicit candidate-only real-provider profile for RESEARCHCTL-015.

- `study.js` compiles the candidate study with `fixtureProviders: false`.
- `product.js` compiles the corresponding product plan.
- `manifests/` and `schemas/` are public, immutable provider identity inputs.
- `provider-config.example.yaml` is operational-only; copy it outside this directory and set the endpoint environment references before use.

The placeholder corpus and evaluation envelope digests in `inputs.json`, `study.js`, and `product.js` must be replaced with the custody-verified candidate artifacts before a run. This is not a benchmark claim.
