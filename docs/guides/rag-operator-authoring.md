# Authoring a native RAG operator

A new semantic enters RAG through a versioned Go operator, compiler definition, typed descriptor, tests, and documentation. Arbitrary JavaScript executors and unregistered nodes are intentionally unsupported.

## 1. Choose an immutable identity

Use `<namespace>.<operation>/<version>`, for example `retrieve.vector/v1`. Pick a new version whenever observable behavior changes: defaults, ordering, scoring, truncation, normalization, lineage, failure behavior, or trace shape.

Never reuse a version for a behavior change and never add aliases for an old prototype name.

## 2. Add the compiler definition

In `pkg/ragcompiler/registry.go`, declare:

- execution phase;
- typed input and output ports;
- required/optional config fields and defaults;
- capabilities/resources;
- expected observations.

Normalization must produce canonical config and deterministic node identity. Config decoding rejects unknown fields. Validate unsafe combinations before execution.

## 3. Implement the runtime operator

Implement `ragoperators.Operator`:

```go
type Operator interface {
    Ref() ragcontract.OperatorRef
    Execute(context.Context, ragcontract.Node, map[string]any, *Environment) (map[string]any, error)
}
```

Register it once in `pkg/ragoperators/registry.go`. Duplicate registration must fail. Runtime and compiler registries have a parity test; a compiler-only operator is invalid.

Required properties:

1. honor context cancellation around every expensive/provider operation;
2. decode canonical config strictly;
3. validate runtime value types and dimensions;
4. sort all order-insensitive outputs deterministically;
5. return explicit errors—never another algorithm as fallback;
6. emit complete parent/production manifest lineage;
7. update canonical traces with identities, policies, duration and usage;
8. keep generated text separate from source evidence;
9. exclude credentials and endpoint secrets from every artifact/error/trace.

## 4. Materialization operators

Units, chunks, representations, embeddings and indexes produce immutable records and manifests. Hash record bytes first, then build manifest metadata; do not create circular self-digests. Parent roles and production config are ordered and canonical.

Static operators must be safe to retain in `ragengine.Prepared`. Any retained object must support concurrent reads and deterministic close. A mutable/query-local object must transitively depend on the query input instead.

## 5. Retrieval operators

Retrievers return ranked representation records with stable tie-breaking. They do not hydrate source text. Channel collapse removes repeated parent votes before fusion. Fusion records each channel/rank/weight contribution. Hydration maps winners to exact chunks/source ranges. Rerankers preserve first-stage scores and identities in trace.

## 6. Provider operators

Resolve model/prompt manifests before provider invocation. Validate model dimensions, tokenization/truncation, request parameters, output schema, result cardinality and citation IDs. Cache keys include semantic inputs and manifest/config identity. Record measured tokens/cost; zero means an explicit zero-cost fixture, not missing telemetry.

Endpoint and credential handling belongs to host/provider adapters. Operators receive interfaces, not environment lookups or raw secret values.

## 7. Expose a typed descriptor

Add a descriptor factory in `pkg/ragmodel`. If JavaScript-facing, expose the same typed factory in `pkg/gojamodules/rag` and update precise TypeScript declarations. Configurator functions execute immediately and must leave only Go-owned data.

Do not expose raw node construction, provider callbacks, SQL/index handles, or lifecycle methods.

## 8. Tests

Minimum coverage:

- normalization/default and malformed config tests;
- compiler/runtime registry parity;
- exact deterministic ordering and identity;
- unknown field, missing input and wrong runtime type failures;
- lineage and manifest validation;
- cancellation/provider failure;
- trace completeness and secret canary scan;
- property/fuzz coverage for ranges, graph/config decoding, or cardinality as appropriate;
- product/study parity if usable in both targets;
- race test if state can be prepared/shared;
- benchmark if algorithmic or storage cost is material.

Run:

```bash
go test ./pkg/ragcompiler ./pkg/ragmodel ./pkg/ragoperators ./pkg/ragengine -count=1
go test -race ./pkg/ragoperators ./pkg/ragengine -count=1
go vet ./...
GOWORK=off go run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.12.2 run ./...
xgoja doctor -f examples/xgoja/rag-v2/xgoja.yaml
```

## Review checklist

- [ ] New immutable operator version, no alias.
- [ ] Compiler definition and runtime factory agree.
- [ ] Canonical config/defaults and stable ordering.
- [ ] Explicit capability/failure behavior.
- [ ] Complete lineage and evidence separation.
- [ ] Cancellation, budgets and resource ownership.
- [ ] Trace/usage/cost and secret safety.
- [ ] Typed Go/JavaScript authoring surface.
- [ ] Unit, negative, parity, race/fuzz/benchmark coverage as applicable.
