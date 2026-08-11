---
Title: "Understand the RAG v2 destructive cutover"
Slug: "rag-v2-cutover"
Short: "Review final ownership boundaries and checks that prevent prototype runtimes from returning."
Topics:
- rag
- architecture
- cutover
Commands:
- rag-eval study run
- rag-product-server
Flags: []
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

RAG v2 replaced disposable prototypes without compatibility readers, aliases, dual runners, old database lifecycle, or dormant archived packages. Researchctl owns generic scientific lifecycle; rag-evaluation-system owns RAG semantics; product execution owns online lifecycle.

## Preserve the boundary

- Add RAG semantics as native versioned operators.
- Add study UX only in rag-eval over the public generic laboratory SDK.
- Add online behavior only in the product package/host.
- Never add RAG imports to researchctl.
- Never make product binaries import the research adapter.
- Recreate disposable RAG databases instead of preserving old run rows.

## Verify absence

```bash
scripts/09-phase8-acceptance.sh
```

The gate checks package/dependency trees, commands, routes, database objects, generated declarations, frontend text, tests, race, fuzz, security and canonical reconstruction.

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| An old name seems convenient | A caller still depends on removed prototype behavior | Rewrite the caller against canonical v2; do not add an alias |
| RAG server needs run history | Lifecycle ownership is being crossed | Query/export through researchctl's generic laboratory |
| Product host wants study submission | Online and scientific lifecycle are being mixed | Emit qualification data, then submit separately through rag-eval |

## See also

- `rag-v2-api-reference`
- `rag-product-runtime`
- `rag-study-workflow`
- `docs/guides/rag-v2-destructive-cutover.md`
