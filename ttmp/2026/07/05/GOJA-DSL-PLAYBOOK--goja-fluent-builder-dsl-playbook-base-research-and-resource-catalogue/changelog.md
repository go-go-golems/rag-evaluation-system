# Changelog

## 2026-07-05

- Initial workspace created


## 2026-07-05

Created base research ticket: DSL catalogue (design-doc), research logbook (reference), investigation diary. Inventoried 12 Goja DSLs + 3 cross-cutting pieces across goja-bleve, goja-dbus, go-minitrace, rag-evaluation-system widgetdsl, go-go-goja, goja-text, goja-git, goja-github-actions, goja-treesitter, glazed, go-emrichen. Classified 5 implementation patterns; identified goja-bleve typed-ref substrate as the runtime-typecheck model and goja-dbus as the composable-grammar model.

### Related Files

- /home/manuel/code/wesen/go-go-golems/goja-bleve/pkg/api_types.go — Typed-ref substrate to extract into a shared fluent package
- /home/manuel/code/wesen/go-go-golems/goja-dbus/pkg/dbusgoja/builders.go — Cleanest composable grammar reference


## 2026-07-05

Validated (docmgr doctor passes cleanly after adding fluent-builder + typescript vocabulary). Uploaded bundle 'Goja DSL Playbook — base research' (3 docs, toc-depth 2) to reMarkable at /ai/2026/07/05/GOJA-DSL-PLAYBOOK and verified listing. Ticket ready for senior-researcher handoff.


## 2026-07-05

Design-doc 02 added: insider self-assessment of the widgetdsl grammar. Five silent-failure modes empirically verified against the built binary (typo'd arrange/options/verbs, wrong markers, out-of-range enums all absorbed without error); root cause named as IR-as-API vs IR-as-output; v2 sketch keeps the grammar language on typed Go specs with configurator builders, strict option decoding, and error-panel rendering; eleven proposed catalogue additions including a silent-ignore failure taxonomy, a tree-shaped-data axis, dormant type sources (ir.ts + 80 widget manifests, zero consumers), and agent-authorship as a design constraint.

