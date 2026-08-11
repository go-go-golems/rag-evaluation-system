---
Title: "Author a native RAG operator"
Slug: "rag-operator-authoring"
Short: "Add immutable operator semantics across compiler, runtime, authoring, traces, and tests."
Topics:
- rag
- operators
- development
Commands: []
Flags: []
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: Tutorial
---

A RAG semantic is a versioned Go operator, not arbitrary JavaScript. Add its compiler definition, runtime implementation, typed descriptor, trace/lineage behavior, and tests together so authoring and execution cannot drift.

## Implement the vertical slice

1. Choose immutable `<namespace>.<operation>/<version>` identity.
2. Add typed ports, canonical config/defaults, capabilities and observations to the compiler registry.
3. Implement and register the native runtime operator.
4. Emit deterministic outputs, complete lineage, trace timing/usage and explicit failures.
5. Expose a typed Go descriptor and JavaScript factory if needed.
6. Add malformed, identity, lineage, cancellation, secret, parity and race/fuzz tests as applicable.

Behavior changes require a new operator version. Never add an alias or silent fallback.

## Validate

```bash
go test ./pkg/ragcompiler ./pkg/ragmodel ./pkg/ragoperators ./pkg/ragengine -count=1
go test -race ./pkg/ragoperators ./pkg/ragengine -count=1
xgoja doctor -f examples/xgoja/rag-v2/xgoja.yaml
```

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| Compiler accepts but runtime fails lookup | Registries disagree | Add the runtime factory and run registry parity tests |
| Identity changes between runs | Unsorted output or noncanonical config | Sort and canonicalize before hashing |
| Citation points to generated text | Hydration/evidence boundary was bypassed | Return representations to retrieval and hydrate source chunks separately |
| Prepared race | Static value is mutable or not concurrency-safe | Make it immutable/thread-safe or query-dependent |

## See also

- `rag-v2-api-reference`
- `docs/guides/rag-operator-authoring.md`
