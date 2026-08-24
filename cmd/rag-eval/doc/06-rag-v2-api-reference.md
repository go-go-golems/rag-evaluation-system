---
Title: "RAG v2 API and contract map"
Slug: "rag-v2-api-reference"
Short: "Find canonical schemas, packages, targets, identities, policies, and evidence rules."
Topics:
- rag
- api
- contracts
Commands:
- rag-eval study validate
- rag-eval study compile
- rag-workflow-runner
- rag-product-server
Flags: []
IsTopLevel: true
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

RAG v2 has one canonical contract/compiler/operator stack with separate study and product lifecycle targets. This page provides the shortest route to the right package and contract.

## Contracts

- `rag-pipeline-ir/v2`: normalized operator DAG.
- `rag-product-plan/v2`: online bindings and policies.
- `rag-study/v2`: variants, factors, dataset and measures.
- `rag-pipeline-execution/v2`: one expanded cell.
- `scraper-workflow-execution/v2`: production Workflow runner input.
- `rag-workflow-study-bundle/v1`: compiled study bundle and plan manifest.
- `rag-query-trace/v2`: authoritative domain evidence and usage.
- `rag-product-qualification/v1`: exact deployment qualification.

## Packages

- contracts: `pkg/ragcontract`;
- normalization/targets: `pkg/ragcompiler`;
- pure authoring: `pkg/ragmodel`, `pkg/gojamodules/rag`;
- native behavior: `pkg/ragoperators`; `pkg/ragengine` is retained as a semantic harness, not a scheduler;
- study input resolution: `pkg/researchctladapter`;
- production lowering and tasks: `pkg/ragworkflow`, `pkg/ragworkflowops`, `cmd/rag-workflow-runner`;
- products: `pkg/ragproduct`, `cmd/rag-product-server`.

RAG-eval compiles immutable Workflow cases; Researchctl executes the generated experiment plan. No RAG command owns replicate scheduling, process retries, or run persistence.

## Evidence invariants

Representations are retrieval material. Collapse identity controls voting. Hydrated source chunks/ranges are evidence and citations. Collapse occurs per channel before fusion. Unit evaluation deduplicates mapped parent identities.

## Troubleshooting

| Problem | Cause | Solution |
|---|---|---|
| Unknown field/trailing JSON | Strict codec rejected noncanonical input | Emit exactly the current schema |
| Noncanonical pipeline | Graph/default/order differs from compiler output | Compile before execution |
| Manifest lineage error | Bound digest/parent/production mismatch | Resolve the exact immutable envelope |
| Unsupported operator | Compiler/runtime capability mismatch | Register a new native operator version |

## See also

- `rag-study-workflow`
- `rag-product-runtime`
- `rag-operator-authoring`
- `docs/guides/rag-v2-api-reference.md`
