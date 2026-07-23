---
Title: RAG Workflow lowering acceptance and boundary audit
Ticket: RAG-V2-WORKFLOW-LOWERING
Status: complete
Topics:
  - rag
  - workflow
  - scraper
  - evaluation
DocType: analysis
Intent: long-term
Summary: Acceptance evidence for deterministic provider-free RAG v2 lowering into Scraper Workflow V3 and Researchctl.
---

# RAG Workflow lowering acceptance and boundary audit

## Executive result

The provider-free RAG v2 lowering is accepted. Canonical `rag-pipeline-execution/v2` remains the semantic authority. `pkg/ragworkflow` lowers that contract through a closed exact operator registry into generic Workflow V3 plans and executes those plans with the versioned `rag-v2-provider-free@1.0.0` task package. Scraper and Researchctl remain domain neutral.

The built-binary acceptance scheduled two cases with two replicates each. All four selected runs succeeded, one injected runner crash produced a separate Researchctl attempt and subordinate Workflow database, complete resume executed no new work, and a forced timeout durably canceled its subordinate Workflow. Provider-free Workflow results matched `ragengine` retrieval ordering and MRR values.

## Ownership audit

| Authority | Evidence | Result |
|---|---|---|
| RAG authoring and canonical execution | `ragmodel`, `ragcompiler`, `ragcontract.PipelineExecution`; lowerer accepts only normalized execution plus exact `CellID` | preserved |
| RAG operator lowering | `pkg/ragworkflow/registry.go`; exact kind/version entries and explicit provider-required attachments | closed |
| Workflow lifecycle | generic Scraper `workflowv3*` packages; RAG package contributes only tasks/modules/projector | no duplicate scheduler/store/retry/status |
| Experiment lifecycle | Researchctl `Laboratory.Execute`, immutable `(specification_id, replicate_index)`, process runner | preserved |
| Provider execution | provider-backed registry entries return `RAG_WORKFLOW_PROVIDER_REQUIRED`; provider-free operation evidence has zero records | deferred, not faked |

Boundary guard: `scripts/03-boundary-privacy-deletion-guards.sh` rejects Researchctl/provider/legacy imports from `pkg/ragworkflow`, RAG/TTC semantics in affected generic Scraper/Researchctl packages, unfinished markers, local module replacements, privacy canaries, stale generated fixtures, and returned TTC custody paths.

## Contracts and deterministic identity

The lowering validates the canonical pipeline, exact execution `CellID`, corpus binding schema/digest, dataset manifest digest, operator kind/version, fixture embedding model/config, package/catalog identity, and Workflow plan digest. It rejects unknown fields and trailing JSON.

The preparation fingerprint includes:

- canonical pipeline digest, including operator versions and normalized configs;
- immutable corpus semantic digest;
- deterministic embedding implementation identity.

It excludes query text and evaluation labels. Every `rag-workflow-query/v1` item repeats the execution's dataset manifest digest. Archive construction and query execution both reject stale dataset identity. Completed preparation nodes and artifacts survive SQLite restart; the focused parity test closes the store after all preparation tasks and resumes query map/reduction/publication in a new store/runtime instance.

Query maps admit at most 10,000 items, page by eight, and materialize at most eight ahead. Reduction fan-in is 16 with at most four levels. Research runner archives, output bodies, projected observation counts, and aggregate projected bytes are bounded.

Generated fixture and plan bytes are reproducible. `biome.json` excludes `examples/rag-workflow` because formatter rewriting would invalidate byte digests recorded in the fixture manifest. The generator/guard compares all generated files byte-for-byte.

## Provider-free semantic coverage

The production task package executes:

1. corpus load;
2. immutable unit preparation;
3. recursive/identity chunking;
4. raw representation;
5. deterministic fixture embedding;
6. lexical/vector multi-index build and index evidence capture;
7. index reconstruction with exact manifest/artifact digest validation;
8. BM25/vector retrieval;
9. collapse and weighted reciprocal-rank fusion;
10. source-evidence hydration;
11. relevance evaluation;
12. bounded result reduction and canonical publication.

Durable preparation values use the existing `ragengine` closed codec. Live indexes and clients are never serialized; indexes are deterministically rebuilt and compared with recorded manifest, record digest, and size.

## Parity and lifecycle evidence

Fresh smoke summary: `sources/smoke/01-summary.json`.

| Check | Evidence |
|---|---|
| Matrix | 4 scheduled, 4 executed, 0 failed |
| Independent durability | 5 subordinate Workflow SQLite databases after one injected runner crash |
| Research retries | selected attempt counts `[1,1,1,2]` |
| Provider-free parity | 2 cases × 2 queries; exact result ordering/traces and `rag.mrr` values matched `ragengine` |
| Canonical observations | 24 metrics and 5 traces per selected run |
| Artifact custody | 4 verified artifacts per selected run |
| Provider custody | 0 external-operation records and verified empty operation artifact/manifest |
| Restart | fresh-process observation reprojection equaled exported observation artifact |
| Preparation restart | unit test closes SQLite after preparation and resumes query/reduction/publication |
| Resume | second plan invocation executed 0 and resumed 4 immutable runs |
| Timeout/cancellation | Researchctl terminal kind `timeout`; subordinate Workflow terminal status `canceled`; canceled observation reprojected after restart |
| Failure mapping | malformed corpus maps to permanent `RAG_WORKFLOW_TASK_FAILED`; injected process crash retries only at Researchctl attempt boundary |
| Privacy | payload canary absent from typed task error and domain projection; trace metadata/filter/failure message/details redacted or digested |

## Validation record

The following completed successfully after final implementation changes:

```text
# RAG-eval
GOWORK=off go test ./... -count=1
GOWORK=off go test -race ./pkg/ragworkflow -count=1
GOWORK=off go build ./...
GOWORK=off golangci-lint run ./...
GOWORK=off go mod tidy                      # no unexplained diff

# Scraper
GOWORK=off go test ./... -count=1
GOWORK=off go test -race ./pkg/researchrunner ./pkg/workflowv3product ./pkg/workflowv3observations -count=1
GOWORK=off go build ./...
GOWORK=off golangci-lint run ./...
GOWORK=off go mod tidy                      # no diff

# Researchctl
GOWORK=off go test ./... -count=1
GOWORK=off go test -race ./pkg/lab/processrunner ./internal/labsqlite -count=1
GOWORK=off go build ./...
GOWORK=off golangci-lint run ./...
GOWORK=off go mod tidy                      # no diff

# Cross-repository
bash scripts/02-smoke-rag-workflow-lowering.sh
bash scripts/03-boundary-privacy-deletion-guards.sh
```

The first parallel Scraper full-test run had one lease timing failure: `TestHTTPSnapshotRetriesAndReopensWithoutPersistingRequestSecrets` observed an additional `lease_lost` attempt and expected two attempts. The focused test passed five consecutive reruns, and a subsequent sequential full Scraper suite passed. This was environmental scheduler contention during three simultaneous full repository suites, not a product invariant failure.

The first RAG full-lint invocation failed because another parallel `golangci-lint` held the process lock. The sequential invocation exposed six pre-existing findings (unchecked close errors, a non-exhaustive status switch, and a helper named `max`). Those findings were corrected in focused commit `5ec4f16`, and full lint then passed.

The first command/help smoke redirected only stdout, while Go's standard `flag` package writes usage to stderr. The corrected smoke redirects both streams and verifies the new help topic and all three binary flag surfaces.

## Focused commits

### Scraper

- `ebf9a85` — `feat: support domain task packages and set inputs`
- `c497bc1` — `feat: project bounded domain observations`
- `391c2d0` — `test: keep domain projection fixture neutral`

### Researchctl

- `8d68226` — `fix: preserve generic input and trace identities`

### RAG-eval

- `5ec4f16` — `chore: resolve repository lint findings`
- `f7cd54a` — `feat: lower RAG v2 into Workflow V3`
- `33bcffa` — `test: add RAG Workflow cross-repository acceptance`
- `698928c` — `fix: preserve canonical generated fixture bytes`

## Review route

1. Read `pkg/ragworkflow/registry.go` and `lower.go` for accepted semantics and generated graph shape.
2. Read `package.go`, `task.cjs`, and `runtime.go` for the task ABI and durable preparation/query boundary.
3. Read `projector.go` and its privacy assertions in `runtime_test.go`.
4. Read Scraper `pkg/researchrunner` changes for bounded set input and domain-neutral projection.
5. Read Researchctl resolver and trace query changes for generic identity/ordinal correctness.
6. Run the fixture generator, smoke, and guards from this ticket.

## Explicit non-goals

This phase does not execute Geppetto or any real provider. The registry keeps exact provider-required attachment points for structured representation, reranking, and generation. Their implementation belongs to `RAG-GEPPETTO-WORKFLOW-OPERATIONS`. TTC remains a later acceptance workload and is not imported by this package.
