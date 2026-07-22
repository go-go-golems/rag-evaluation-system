---
Title: "Host a prepared RAG v2 product"
Slug: "rag-product-runtime"
Short: "Load a product plan, prepare immutable indexes, and serve bounded online requests."
Topics:
- rag
- product
- deployment
Commands:
- rag-product-server
Flags:
- plan
- corpus
- address
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: Application
---

The product target runs the same normalized operators as studies but uses host-owned request lifecycle and policy. It imports no researchctl package, prepares static corpus/index nodes once, and validates every request before provider execution.

## Start the reference host

```bash
rag-product-server --plan product-plan.json --corpus corpus-artifact.json --address 127.0.0.1:8780
curl -sS -H 'Content-Type: application/json' \
  -d '{"values":{"query":"What is RRF?"}}' http://127.0.0.1:8780/v1/query
```

The reference host wires deterministic fixture providers only. A production host supplies audited providers, exact manifests, endpoint policy and credentials explicitly.

## Choose explicit policy

Product plans declare timeout, maximum concurrency, citation behavior, failure behavior and trace custody. Provider failure may fail, abstain, or return hydrated retrieval-only evidence. No policy silently invokes another provider.

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| Startup rejects a corpus | Bound and materialized manifest digests differ | Rebuild the plan against the exact corpus envelope |
| Trace sink required | Plan uses artifact-backed trace policy | Supply a host `ArtifactSink` |
| Request rejected | Unknown/missing/type/length contract violation | Match the product request fields exactly |
| Provider operation fails | Capability or exact manifest is missing | Bind the declared provider; do not add fallback |

## See also

- `rag-v2-api-reference`
- `rag-study-workflow`
- `docs/guides/rag-product-runtime.md`
