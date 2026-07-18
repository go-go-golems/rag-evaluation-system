# Canonical RAG v2 API reference

This is the final active API map. Historical ticket documents explain earlier prototypes; they are not supported APIs.

## Package ownership

| Package | Owns | Must not own |
|---|---|---|
| `pkg/ragcontract` | Wire DTOs, schema/version constants, strict codecs, manifest and trace records | Goja, providers, filesystems, databases, researchctl |
| `pkg/ragcompiler` | Operator definitions, defaults, normalization, validation, study expansion, semantic identity | Provider execution or lifecycle |
| `pkg/ragmodel` | Pure typed Go authoring model used by JavaScript | I/O, credentials, workers, indexes |
| `pkg/gojamodules/rag` | `require("rag")` factories and immediate configurators | Retained callbacks or execution |
| `pkg/ragoperators` | Versioned native operator implementations and domain artifacts | Research run allocation/persistence |
| `pkg/ragengine` | Canonical graph execution and prepared static state | Product HTTP or researchctl lifecycle |
| `pkg/researchctladapter` | RAG input resolution and generic researchctl submission | Generic laboratory persistence internals |
| `pkg/ragproduct` | Online request validation, prepared execution, policies, response and qualification | Researchctl dependencies |

Dependency direction is product/adapter → compiler/contracts/engine/operators. Researchctl is imported only by `pkg/researchctladapter`.

## Wire contracts

- `rag-pipeline-ir/v2`: normalized operator DAG.
- `rag-product-plan/v2`: online bindings, request/response, citation and runtime policies.
- `rag-study/v2`: variants, factors, exact artifacts, dataset, measures and replicates.
- `rag-pipeline-execution/v2`: one expanded study cell.
- `rag-query-trace/v2`: query digest, operator/channel/collapse/fusion/hydration/reranking/generation/evaluation evidence.
- `rag-product-qualification/v1`: exact product/model/prompt bindings plus equivalent study.
- `rag-preview-request/v1`: pure preview selection, never lifecycle execution.

Operators have independent immutable IDs `<namespace>.<operation>/<version>`, beginning at `/v1`. Operator versions are not contract schema versions.

## Strict decoding and identity

```go
pipeline, err := ragcontract.DecodePipeline(reader)
product, err := ragcontract.DecodeProduct(reader)
study, err := ragcontract.DecodeStudy(reader)
execution, err := ragcontract.DecodeExecution(reader)
trace, err := ragcontract.DecodeTrace(reader)
```

All codecs reject unknown fields and trailing JSON values. Compilation applies defaults, recursively resolves factors, expands recipes, validates ports/configs/bindings, topologically sorts nodes, and canonicalizes JSON. Identity APIs:

```go
id, err := ragcompiler.ProductSemanticIdentity(product)
id, err := ragcompiler.StudySemanticIdentity(study)
id, err := ragcontract.Digest(value)
```

Display metadata is excluded where documented. Changing pipeline semantics, bindings, policies, factors, dataset, or measures changes the corresponding identity.

## Authoring entry points

```go
pipeline := ragmodel.NewPipeline(name, configure)
query := ragmodel.NewQueryPlan(name, configure)
product := ragmodel.NewProduct(name, configure)
study := ragmodel.NewStudy(name, configure)

productPlan, err := ragmodel.CompileProduct(product, options)
studyPlan, cells, err := ragmodel.CompileStudy(study, options)
```

JavaScript uses the same Go-owned values through `require("rag")`:

```javascript
const rag = require("rag");
const p = rag.pipeline("name", p => /* immediate configuration */);
const product = rag.product("service", p => /* compose values */);
module.exports = product.compileProduct(inputs);
```

No configurator, JavaScript function, runtime object, secret, or database handle enters IR.

## Study execution

`pkg/researchctladapter` resolves catalog aliases to verified RAG envelopes, wraps opaque canonical domain config in public `researchctl/pkg/lab` contracts, checks worker capabilities, and invokes the generic process runner. The worker advertises only:

- protocol `researchctl-runner-stdio/v1`;
- runner `rag-worker/v2`;
- domain `rag-pipeline/v2`;
- trace `rag-query-trace/v2`.

Researchctl verifies generic file custody and observation framing. The adapter/worker verify RAG manifest identity, dataset policy, lineage, and canonical execution identity.

## Product execution

```go
plan, err := ragproduct.Load(reader)
runtime, err := ragproduct.New(ctx, plan, ragproduct.Bindings{ /* exact artifacts/providers */ })
defer runtime.Close()
response, err := runtime.Execute(ctx, ragproduct.Request{
    ID: "request-42",
    Values: map[string]any{"query": "..."},
})
```

`Runtime.Close` drains bounded concurrent requests and closes prepared indexes once. Failure policies are `fail`, `abstain`, and `retrieval-only`; trace policies are `authoritative`, `metadata-only`, `artifact-backed`, and `none`. There is no implicit fallback.

## Evidence rules

- Representations are retrieval material, not evidence.
- Collapse identity is distinct from representation and cited chunk identity.
- Collapse happens within each channel before weighted fusion.
- One collapse key contributes at most one vote per channel.
- Hydration selects exact source chunks and ranges.
- Only hydrated source chunks may become citations.
- Unit-target evaluation maps chunks through parent unit identity and deduplicates targets.

## Errors

Stable error prefixes identify the failing boundary:

- `RAG_V2_*`: authoring/compiler/contract errors;
- `RAG_INPUT_*`: input envelope and lineage errors;
- `RAG_ENGINE_*` / `RAG_RUNTIME_*`: execution errors;
- `RAG_PRODUCT_*`: online validation, binding, policy or host errors;
- `RAG_WORKER_*`: process worker negotiation/execution errors.

Errors never trigger silent downgrade. Product-facing provider details are redacted.
