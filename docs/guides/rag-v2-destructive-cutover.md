# RAG v2 destructive cutover guide

The cutover replaced disposable prototypes; it did not preserve their contracts. This guide describes the final state and the checks that prevent old execution paths from returning.

## Final ownership

- Researchctl owns domain-neutral projects, specifications, runs, attempts, retries, observations, verified artifact custody, persistence and export/import.
- Rag-evaluation-system owns every RAG contract, compiler, typed JavaScript API, operator, engine, adapter, CLI, worker, trace, metric and product runtime.
- The legacy playground repository retains only historical ticket evidence in Git history; its executable playground was deleted.

Dependency direction is RAG → public researchctl laboratory SDK. Researchctl has no RAG vocabulary, package, command, schema, example or dependency. Product packages import neither researchctl nor the research adapter.

## Deleted categories

The destructive removal covered:

- prototype RAG DTO/builders, catalog/executor/persistence and worker packages;
- competing immutable retrieval implementations;
- old JavaScript lifecycle and descriptor APIs;
- domain-specific researchctl command, loader, adapter, schemas, help and examples;
- old worker protocol and trace schemas;
- writable RAG experiment HTTP lifecycle;
- local RAG specification/run/event/summary/query-trace database objects;
- playground generation/cache/index active pointers, duplicated hashing/evaluation and self-tests;
- stale generated declarations, UI commands and help text.

No converter, alias, dual reader, dual runner, environment feature flag, deprecated route, archived runtime package or warning-and-continue path remains.

## Database policy

RAG domain databases are disposable artifact catalogs. Fresh migration creates current sources/documents/chunks, immutable corpus/chunk/embedding/retrieval artifacts, evaluation datasets and representation sets. It creates no scientific run lifecycle objects; recreate development databases rather than copying prototype run rows.

Researchctl databases remain the sole run evidence store.

## Route policy

Rag-evaluation-system exposes the read-only domain artifact catalog at:

```text
GET /api/v1/artifacts/rag/catalog
```

There is no RAG experiment specification/run/event/completion/comparison route. The UI invokes the RAG-owned CLI workflow and displays domain artifact availability only.

## Absence enforcement

`pkg/ragcontract/cutover_test.go` walks active source/docs/examples/web trees and fails if a retired package, schema, protocol, command, route, table, lifecycle method or compatibility marker returns. Database migration tests inspect `sqlite_master`; API tests prove old lifecycle requests are unregistered; command/module tests prove old JavaScript exports are undefined.

The final acceptance script additionally searches both repositories, checks dependency graphs, generated TypeScript declarations and command help, then runs tests, race, vet, lint, security and frontend gates.

## Adding functionality after cutover

Do not restore an old name for convenience. Extend the current architecture:

- new RAG behavior → immutable native operator version;
- new authoring convenience → typed descriptor/fragment compiling to existing IR;
- new study UX → RAG CLI over `pkg/researchctladapter`;
- generic lifecycle capability → researchctl public laboratory SDK without RAG terms;
- online behavior → `pkg/ragproduct`/host policy without research lifecycle;
- historical data inspection → Git/ticket evidence, not an executable compatibility package.

Any proposal that requires schema detection, old field decoding, alias registration, dual persistence or silent fallback is a new architecture decision and violates the completed cutover by default.

## Verification

```bash
scripts/09-phase8-acceptance.sh
docmgr doctor --ticket RESEARCHCTL-014 --stale-after 30
```

A successful build is not enough. Completion requires clean negative searches, fresh-database inspection, current UI/docs, generated declarations, product dependency isolation, canonical export reconstruction and clean repositories.
