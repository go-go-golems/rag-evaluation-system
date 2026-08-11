# RAG v2 intern onboarding

This path gets an unfamiliar engineer from contracts to a validated change without learning deleted prototypes first.

## Day 1: map the architecture

Read in order:

1. `docs/guides/rag-v2-api-reference.md`;
2. `pkg/ragcontract/README.md`, `types.go`, `manifests.go`, `trace.go`;
3. `pkg/ragcompiler/registry.go`, `normalize.go`, `targets.go`;
4. `pkg/ragmodel/descriptors.go`, `builders.go`, `compile.go`;
5. `examples/rag-v2/README.md` and `experiments/rag-sol2/README.md`.

Draw this flow:

```text
pure JS/Go authoring
  → normalized pipeline IR
  → product plan or expanded study cells
  → exact artifacts/providers
  → native operators/engine
  → product response or generic researchctl observations
```

Be able to explain why representation, collapse and evidence identities differ.

## Day 2: follow one raw-BM25 query

Start at `examples/rag-v2/06-raw-study.js`, then follow:

- unit/chunk/raw representation in `pkg/ragoperators/prepare.go` and `represent.go`;
- Bleve index in `index.go`;
- retrieval/collapse/fusion/hydration in `rank.go`;
- evaluation in `evaluate.go`;
- orchestration in `pkg/ragengine/engine.go`;
- worker envelopes in `cmd/rag-worker` and `pkg/researchctladapter`.

Run:

```bash
rag-eval study validate examples/rag-v2/06-raw-study.js \
  --inputs examples/rag-v2/inputs-ttc-catalog.json \
  --ttc-database data/rag-eval.db

go test ./pkg/ragoperators ./pkg/ragengine ./pkg/researchctladapter -count=1
```

## Day 3: product target

Read `docs/guides/rag-product-runtime.md`, `pkg/ragengine/prepared.go`, `pkg/ragproduct/runtime.go`, and `qualification.go`. Verify that product dependencies contain no researchctl package and that qualification preserves pipeline bytes.

```bash
go test -race ./pkg/ragproduct ./pkg/ragengine ./cmd/rag-product-server -count=1
```

## First safe change

Choose a documentation fix, validation error improvement, or additional malformed-input test. Avoid changing operator behavior until you can identify:

- operator version and config defaults;
- input/output port kinds;
- manifest parents and production identity;
- ordering/tie-break rules;
- trace entries;
- product and study implications;
- required golden/property/fuzz/race coverage.

Behavior changes require a new operator version unless they correct a proven invariant violation documented in the ticket/diary.

## Repository boundaries

The workspace contains multiple repositories. Commit deliberately:

- researchctl generic SDK/lifecycle changes in researchctl;
- canonical RAG changes in rag-evaluation-system;
- no new runtime in the historical playground repository;
- ticket documentation under the RESEARCHCTL-014 workspace.

Never introduce a researchctl → RAG dependency. Never import `pkg/researchctladapter` from product code.

## Scientific vocabulary

- **smoke**: tiny engineering execution check;
- **candidate**: exploratory dataset/result, not frozen adjudication;
- **qualification**: exact deployment bindings wrapped for study;
- **replicate**: explicitly fresh scientific execution;
- **manifest digest**: domain artifact identity;
- **file digest**: generic custody identity;
- **source evidence**: hydrated source chunk/range, never generated text.

Do not call candidate matrices benchmarks without adjudication, holdout and provider freeze.

## Debugging checklist

1. Validate schema and canonical compilation.
2. Compare semantic IDs and exact bindings.
3. Check compiler/runtime registry parity.
4. Inspect manifest parent/production lineage.
5. Inspect channel hits before and after collapse.
6. Inspect fusion contributions and hydration.
7. Confirm relevance target identity and deduplication.
8. Check cancellation/resource budgets/provider manifests.
9. Scan traces/errors/artifacts for secret canaries.
10. Reproduce with focused unit test before broad acceptance.

## Definition of done

- focused tests and malformed cases;
- deterministic identities/goldens updated intentionally;
- race/fuzz/benchmark coverage when relevant;
- full active and standalone tests;
- vet and pinned lint;
- frontend/type declarations/help if affected;
- diary/changelog/tasks and absolute file relations;
- clean repository status and exact commit hashes.
