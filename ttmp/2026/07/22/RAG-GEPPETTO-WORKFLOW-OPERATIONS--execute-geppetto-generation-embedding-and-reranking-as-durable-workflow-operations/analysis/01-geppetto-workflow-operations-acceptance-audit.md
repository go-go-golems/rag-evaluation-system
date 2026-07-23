---
Title: Geppetto Workflow Operations Acceptance Audit
Ticket: RAG-GEPPETTO-WORKFLOW-OPERATIONS
Status: active
Topics:
    - rag
    - workflow
    - geppetto
    - research
DocType: analysis
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/scraper/pkg/workflowv3runtime/engine.go
      Note: |-
        Renewable lease lifecycle used by long provider calls
        Renewable Workflow lease authority
    - Path: pkg/ragworkflow/provider_package.go
      Note: |-
        Versioned provider package, authority, budgets, and environment wiring
        Versioned provider package and canonical authority
    - Path: pkg/ragworkflow/provider_real_acceptance_test.go
      Note: Authorized real-provider acceptance
    - Path: pkg/ragworkflowops/operations.go
      Note: |-
        Durable provider operation decorators and safe accounting
        Durable operation lifecycle and accounting
    - Path: repo://cmd/rag-workflow-runner/main.go
      Note: Production provider runner composition
ExternalSources: []
Summary: Requirement-to-evidence audit for production Geppetto-backed RAG Workflow V3 operations.
LastUpdated: 2026-07-23T22:10:00-04:00
WhatFor: Verify provider operation identity, custody, lifecycle, privacy, parity, and cross-repository acceptance before closing the ticket.
WhenToUse: Review implementation completeness or reproduce provider Workflow acceptance.
---


# Geppetto Workflow Operations Acceptance Audit

## Result

The implementation satisfies the ticket's production contract. Canonical RAG v2 remains the authoring and semantic authority. The distinct `rag-v2-geppetto@1.0.0` package lowers and executes every provider-required attachment point through thin RAG-owned Geppetto adapters. Scraper remains the sole authority for leases, attempts, retries, external-operation records, budgets, cancellation, artifacts, and observations. Researchctl remains the sole cross-run authority.

Deterministic and authorized real-provider acceptance both passed. Geppetto source code is unchanged. The only generic infrastructure change is fenced Workflow V3 lease renewal, required to prevent duplicate provider contacts during long calls.

## Contract and identity evidence

| Requirement | Evidence |
|---|---|
| Exact versioned operator registry | `pkg/ragworkflow/registry.go`; structured, synthetic, combined, embedding, rerank, and answer variants execute in `TestProviderPackageBindsAuthorityAndExecutesDurableOperations`. |
| Distinct provider catalog | `ProviderPackageName=rag-v2-geppetto`, version `1.0.0`; provider-free lowering rejects provider-required operators. |
| Canonical provider authority | `ProviderAuthority` records schema, profile, fixture marker, capabilities, sorted manifest digests, safe provider identities, settings fingerprints, concurrency, response-token limits, and pricing policy; digest is embedded as `provider-authority.json`. |
| Strict package/host match | `NewProviderLowerer`, `BuildProviderRunnerExecution`, and provider task factory bind the same authority and preparation fingerprint. Unknown/malformed authority fails validation. |
| Exact scientific identity | Provider authority participates in preparation fingerprints; canonical execution, corpus, dataset manifest, factor identity, and query boundaries remain unchanged. |
| No compatibility shim | Provider-free and provider-backed packages are separate closed catalogs; neither adapts the other at runtime. |

## Provider operation custody

Three descriptors are production contracts: `provider.generate/v1`, `provider.embed/v1`, and `provider.rerank/v1`.

Each decorator performs strict preflight, computes a bounded semantic correlation digest, obtains Workflow admission, invokes exactly one underlying provider method, validates post-contact domain evidence, and persists completion through a detached bounded context. The ledger retains provider start time, elapsed microseconds, outcome, safe failure class/code, reservation, cardinality measures, actual/conservative accounting mode, usage counters, authority digest, node, attempt, and ordinal. It never stores request text or response bodies.

Deterministic full-pipeline counts are exact:

- structured: generation 5, embedding 3, rerank 2;
- combined: generation 4, embedding 3, rerank 2;
- synthetic: generation 8, embedding 3, rerank 2.

The two-case/two-replicate Researchctl acceptance exports exactly ten completed operations per selected successful run. The injected process crash occurs after admission and before completion; the abandoned subordinate database retains exactly one incomplete operation rather than inventing a terminal outcome.

## Retry, replay, cache, and idempotency

A typed transport failure is mapped to a retryable safe Workflow failure. The next Workflow attempt creates a new operation ID and preserves the semantic correlation digest. The integration test proves one failed contact plus a succeeding retry. No failed operation is overwritten.

Provider preparation cache identities include operator config, parent digest, model manifest, prompt manifest, output schema fingerprint, and effective provider settings. The second identical Workflow proves that cache hits produce no generation operation and settle request usage as zero. Query-time contacts remain visible. Cache reuse therefore cannot masquerade as a provider request.

The correlation digest is the portable semantic idempotency/replay key. Physical provider idempotency is not claimed where a Geppetto/provider API does not expose it. A process crash after contact can therefore produce an uncertain open operation and a later distinct retry; the evidence preserves that uncertainty.

## Failure, timeout, cancellation, and lease behavior

| Condition | Evidence |
|---|---|
| Preflight rejection | No admission and no provider contact. |
| Budget/admission rejection | No provider contact; permanent admission failure. |
| Transport/rate-limit/provider-5xx | Safe allow-listed code, failed operation, retryable task policy. |
| Caller cancellation | Canceled operation; completion uses detached bounded context and conservative reservation. |
| Deadline | Timed-out operation with `PROVIDER_TIMEOUT`. |
| Provider succeeds but domain result is malformed | Provider operation remains succeeded; task result is permanent malformed output. |
| Completion persistence fails | Domain success is not published. |
| Long provider call | Scraper renews the exact fenced lease at half-life. Authorized acceptance uses a two-second lease and completes each node in one attempt despite longer generation calls. |
| Process crash after admission | One durable incomplete operation; Researchctl attempt retry creates a separate subordinate Workflow. |
| Experiment timeout | Researchctl failure kind `timeout`; subordinate Workflow reaches durable `canceled`. |

Scraper renewal requires matching run, node, token, cancel epoch, running status, and an unexpired lease. Cancellation blocks renewal. Store, runtime heartbeat, focused race, full Scraper, and real-provider tests cover this lifecycle.

## Usage, cost, limits, and batching

- Every contact reserves one request before provider invocation; task and external-operation request accounting are reconciled transactionally.
- Generation records input tokens, output tokens, cache token fields when supplied, and normalized `cost_microunits`.
- Embedding records input/output cardinality, embedding tokens when supplied, dimensions through domain validation, and cost when supplied.
- Rerank records candidate/output cardinality, input tokens, cost when supplied, and complete ordering evidence in the RAG trace.
- The Workflow plan has a finite provider request account and finite per-task claims. `MaxPerAttempt` is bounded and overflow-checked.
- Operator configs retain finite batch sizes, concurrency, timeout, context, response-token, map-item, materialization, and artifact limits.
- Request budget is the enforced Workflow budget dimension. Token and cost are authoritative bounded operation counters; model manifests retain response-token and pricing authority. The implementation does not claim an unenforced aggregate currency ceiling.

## Privacy audit

Safe portable evidence contains only bounded IDs, digests, integers, timestamps, operation kinds, and allow-listed failure codes. Tests inject secret canaries into prompts and provider errors. Neither completion records nor task failures contain them. Generic Scraper observation and Researchctl exports exclude provider request/response bodies, credentials, headers, evidence text, locators, and arbitrary error strings.

`03-boundary-privacy-guards.sh` additionally proves:

- Scraper and Researchctl do not depend on RAG-eval or Geppetto;
- RAG provider code does not import Scraper's legacy workflow engine;
- no local Scraper replacement exists;
- generic operation contracts expose no raw provider payload fields;
- canonical fixtures contain no secret-like values;
- provider fixtures regenerate byte-identically.

## Semantic parity

For deterministic fixtures, direct `ragengine` and Workflow execution match exact retrieval result traces, metrics, and answers. All structured/synthetic/combined variants pass. The built-binary matrix independently compares each selected result with checked-in parity evidence.

The authorized acceptance binds one bounded source and query from `data/ttc-wordpress-rag.sqlite`, then uses real Geppetto generation through `umans-flash`, real Ollama embeddings, and a real llama.cpp cross-encoder. Recursive chunking produces six preparation batches, so the run proves ten contacts (generation 7, embedding 2, rerank 1), actual generation token/cost evidence, rerank tokens, embedding cardinality, one-attempt lease behavior, and grounded citations.

## Cross-repository acceptance

Fresh built binaries executed two cases × two replicates through Researchctl and `scraper-workflow-execution/v2`:

- scheduled 4, executed 4, failed 0;
- Researchctl attempt counts `[1,1,1,2]` after one injected process crash;
- five subordinate Workflow databases including the crashed run;
- ten completed operations per selected run;
- one incomplete operation in the crash-window run;
- 24 RAG metrics, five traces, four verified artifacts per selected run;
- canonical observation reprojection exactly equal after restart;
- immutable resume: executed 0, resumed 4;
- timeout maps to Researchctl failure and durable Workflow cancellation.

Evidence: `sources/smoke/01-summary.json`.

## Validation evidence

Passed:

- RAG focused tests, full pre-commit package/internal suite, full lint, and module tidiness;
- provider fixture byte regeneration;
- focused race tests for provider operations and Scraper renewal;
- full Researchctl suite;
- full Scraper suite except one known concurrent timing occurrence, followed by focused success and two sequential full `workflowv3runtime` successes;
- Geppetto events, observability, embeddings, rerank, factory, and llama.cpp suites;
- built commands and embedded help;
- deterministic matrix smoke, boundary/privacy guards, and authorized real-provider acceptance;
- clean Geppetto, Scraper, and Researchctl trees.

The concurrent Scraper suite occurrence was `TestLazyMapOutputDigestIsIndependentOfConcurrency: Condition never satisfied` at 90 seconds. The same focused test passed in 18.39 seconds and the package passed twice sequentially in 66.41 and 57.63 seconds. This is retained as provenance, not omitted.

## Requirement audit

- Structured, synthetic, combined, embedding, reranking, and answer attachments: proven.
- Thin Geppetto adapters and unchanged Geppetto: proven.
- Strict identity and schema/version rejection: proven.
- Operation custody, request authority, retries, replay identity, cache distinction, and crash uncertainty: proven.
- Cancellation, timeout, failure taxonomy, malformed results, bounded completion, and lease lifecycle: proven.
- Usage/cost/cardinality evidence and finite request budgets: proven.
- Privacy/redaction and generic boundary neutrality: proven.
- Provider-free non-regression and semantic parity: proven.
- Multi-case/replicate, resume, restart, cancellation, and observation export: proven.
- Explicitly authorized real-provider acceptance: proven.

No implementation requirement remains deferred.
