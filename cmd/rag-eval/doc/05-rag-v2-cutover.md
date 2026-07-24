---
Title: "Understand the RAG v2 destructive cutover"
Slug: "rag-v2-cutover"
Short: "Review final ownership boundaries and checks that prevent prototype runtimes from returning."
Topics:
- rag
- architecture
- cutover
Commands:
- rag-eval study compile
- researchctl experiment run-plan
- rag-workflow-runner
- rag-product-server
Flags: []
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

RAG v2 replaced disposable prototypes without compatibility readers, aliases, dual runners, old database lifecycle, or dormant archived packages. Researchctl owns generic scientific lifecycle; Scraper Workflow V3 owns durable production execution; rag-evaluation-system owns RAG semantics and task packages; product execution owns online lifecycle. The removed `rag-worker`, `study run`, and preview command must not return.

## Preserve the boundary

- Add RAG semantics as native versioned operators.
- Add study UX only in rag-eval over the public generic laboratory SDK.
- Add online behavior only in the product package/host.
- Never add RAG imports to researchctl.
- Never make product binaries import the research adapter.
- Recreate disposable RAG databases instead of preserving old run rows.

## Verify absence

```bash
ttmp/2026/07/22/RAG-V2-EXECUTION-CUTOVER--hard-cut-rag-v2-execution-to-the-canonical-experiment-and-workflow-path/scripts/02-cutover-guards.sh
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
