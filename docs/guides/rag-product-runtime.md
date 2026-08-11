# RAG product runtime v2

`pkg/ragproduct` executes a compiled `rag-product-plan/v2` under a host-owned request lifecycle. It has no researchctl dependency and never creates studies, runs, attempts, metrics, or researchctl artifacts during online requests.

## Lifecycle

1. Strictly decode a product plan with `ragproduct.Load`.
2. Resolve and verify the corpus manifest and every declared model binding.
3. Construct provider interfaces explicitly in the host.
4. Call `ragproduct.New`; query-independent unit, chunk, representation, embedding, and index nodes are prepared once.
5. Call `Runtime.Execute` concurrently, up to `runtime.maxConcurrent`.
6. Call `Runtime.Close` after draining requests; prepared indexes close exactly once.

The normalized graph is identical to study execution. Only lifecycle, request validation, trace policy, and failure policy differ.

## Request and response

Requests have a host request ID and values governed by the plan's typed field contract:

```json
{"id":"request-42","values":{"query":"What is weighted reciprocal rank fusion?"}}
```

Unknown fields, missing required fields, type mismatches, and rune-length violations fail before provider execution. Responses contain ranked collapse identities, source chunk citations, optional answer text, abstention/failure state, and policy-controlled trace/trace ID. Query text is represented in traces only by its digest.

## Policies

Accepted trace policies are:

- `authoritative`: return the complete canonical trace (which contains identities/digests, not source or query text);
- `metadata-only`: return operators, timing, usage, and redacted failure metadata;
- `artifact-backed`: write the canonical trace to the host's `ArtifactSink` and return only its trace ID;
- `none`: return no trace.

Accepted provider failure policies are:

- `fail`: fail the request;
- `abstain`: return an explicit empty abstention;
- `retrieval-only`: return hydrated retrieval evidence plus a redacted failure marker.

No implicit provider fallback exists. Adding an alternative provider requires a future explicit, identity-bearing branch contract.

Citations are `required`, `source`, or `none`. Required/source modes return only hydrated source chunk identities and source ranges; generated representations are never citations.

## Reference HTTP host

`cmd/rag-product-server` uses only `net/http`:

```bash
go run ./cmd/rag-product-server \
  --address 127.0.0.1:8780 \
  --plan product-plan.json \
  --corpus corpus-artifact.json

curl -sS http://127.0.0.1:8780/healthz
curl -sS -H 'Content-Type: application/json' \
  -d '{"values":{"query":"What is RRF?"}}' \
  http://127.0.0.1:8780/v1/query
```

The reference host intentionally wires only deterministic fixture providers currently shipped by the repository. A production host must supply audited provider implementations, endpoint allow policy, credentials, and exact manifests explicitly; missing capabilities fail startup or the request.

## Qualification

`ragproduct.Qualify` emits `rag-product-qualification/v1`. It freezes the product semantic ID, exact model bindings, citation/runtime policies, and a `rag-study/v2` containing the byte-equivalent normalized pipeline and exact corpus binding. It performs no researchctl call. The resulting study can be executed later through the ordinary RAG-owned research adapter.

## Validation

```bash
go test ./pkg/ragproduct ./cmd/rag-product-server -count=1
go test -race ./pkg/ragproduct ./pkg/ragengine -count=1
go test ./pkg/ragproduct -run TestProductLatencyProfile -v -count=1
go test ./pkg/ragproduct -bench BenchmarkPreparedProductQuery -benchmem -count=3
! go list -deps ./pkg/ragproduct ./cmd/rag-product-server | grep researchctl
```

Performance values from fixture corpora are engineering smoke measurements, not production SLOs or model-quality evidence.
