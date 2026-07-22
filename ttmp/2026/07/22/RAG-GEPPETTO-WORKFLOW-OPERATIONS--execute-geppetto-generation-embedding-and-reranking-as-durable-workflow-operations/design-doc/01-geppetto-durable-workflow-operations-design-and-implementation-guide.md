---
Title: Geppetto durable workflow operations design and implementation guide
Ticket: RAG-GEPPETTO-WORKFLOW-OPERATIONS
Status: active
Topics:
    - rag
    - rag-eval
    - workflow
    - geppetto
    - embeddings
    - intern-guide
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/code/wesen/go-go-golems/geppetto/pkg/embeddings/embeddings.go
      Note: Embedding provider interface
    - Path: abs:///home/manuel/code/wesen/go-go-golems/geppetto/pkg/engineprofiles/registry.go
      Note: Geppetto profile resolution
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/scraper/pkg/workflowv3/external_operation.go
      Note: Durable operation API
    - Path: repo://pkg/ragproviders/geppetto/generation.go
      Note: Current RAG-to-Geppetto generation adapter
ExternalSources: []
Summary: Thin-adapter design for recording Geppetto generation, embedding, and reranking as durable Workflow V3 external operations.
LastUpdated: 2026-07-22T23:15:00-04:00
WhatFor: Integrate real providers while preserving authority, retry, timing, usage, privacy, and cancellation semantics.
WhenToUse: Use after deterministic RAG workflow lowering works and before adding real-provider experiment workloads.
---


# Geppetto durable workflow operations design and implementation guide

## Program context

This ticket is part of **EXPERIMENT-PLATFORM-CONVERGENCE** with `RESEARCHCTL-EXPERIMENT-PLANS`, `SCRAPER-WORKFLOW-V3-PRODUCT-CUTOVER`, `EXPERIMENT-PLATFORM-SCRAPER-RUNNER`, `SCRAPER-WORKFLOW-OBSERVATIONS`, `RAG-V2-WORKFLOW-LOWERING`, `RESEARCHCTL-EXPERIMENT-ANALYSIS`, `RAG-V2-EXECUTION-CUTOVER`, and `TTC-SCRIPTED-EXPERIMENT-ACCEPTANCE`. The umbrella explains system ownership. This ticket does not change Geppetto.

## Executive summary

Geppetto owns model profiles, inference, embeddings, reranking, provider events, usage, and credentials. Scraper owns external-operation admission, authority, timing, outcomes, and durable effect records. RAG owns the semantic request and result. This ticket adds thin RAG-owned task adapters that preserve all three responsibilities.

Every provider contact follows a strict begin/call/finish sequence. Failed and canceled calls are closed explicitly. Provider payloads and arbitrary errors never enter the workflow ledger. Usage is marked actual only when Geppetto/provider evidence supplies it; otherwise accounting remains conservative or unavailable.

## Existing APIs to understand

Geppetto, read-only for this program:

- `geppetto/pkg/engineprofiles/`: profile stack and resolved settings.
- `geppetto/pkg/inference/engine/` and `runner/`: inference execution.
- `geppetto/pkg/embeddings/`: provider and batch interfaces.
- `geppetto/pkg/rerank/`: reranker interfaces and usage.
- `geppetto/pkg/events/` and `pkg/observability/`: provider call and usage events.
- `geppetto/pkg/js/modules/geppetto/`: existing JS surface.

Scraper:

- `scraper/pkg/workflowv3/external_operation.go`: descriptors, tickets, completions, counters.
- `scraper/pkg/workflowv3runtime/`: task context and recorder injection.

RAG:

- `pkg/ragproviders/geppetto/`: current provider adapter patterns.
- `pkg/ragoperators/`: generation, embedding, reranking request/result types.
- `pkg/ragcontract/manifests.go`: provider and model identity.

## Operation descriptors

Use separate immutable descriptors:

```text
provider.generate/v1
provider.embed/v1
provider.rerank/v1
```

Authorized counters include only bounded integers with explicit roles:

```text
requests          reservation, usage
input_tokens      reservation, usage
output_tokens     reservation, usage
embedding_tokens  reservation, usage
input_items       measure
output_items      measure
cost_microunits   reservation, usage
```

Model/profile fingerprints belong in safe metadata or task configuration, not secret-bearing provider settings.

## Adapter sequence

```pseudo
function generateTask(ctx, request):
    resolved = hostProfiles.resolve(request.profileRef)
    descriptor = registry.require("provider.generate/v1")
    reservation = estimateBoundedAuthority(request, resolved)

    ticket = ctx.externalOperations.begin({
        descriptorDigest: descriptor.digest,
        reservation: reservation,
        measures: {input_items: request.items.length},
        safeMetadata: modelFingerprint(resolved),
    })

    started = nowUTC()
    try:
        result = geppetto.generate(ctx, resolved, request.turns)
        usage, mode = normalizeUsage(result.events)
        ctx.externalOperations.finish(ticket, {
            providerStartedAt: started,
            elapsedMicros: elapsed(started),
            outcome: "succeeded",
            accountingMode: mode,
            counters: usage,
        })
        return ragResult(result)
    catch err:
        failure = classifyAllowListed(err)
        ctx.externalOperations.finish(ticket, {
            providerStartedAt: started,
            elapsedMicros: elapsed(started),
            outcome: outcomeFromContext(ctx, err),
            failure: failure,
            accountingMode: "conservative",
        })
        throw typedTaskError(failure)
```

`FinishExternalOperation` must execute even when result decoding fails after the provider responds. Use a guarded state machine rather than a casual `defer` that cannot report the right outcome.

## Retry ownership

- Geppetto transport-internal retries, if any, must be visible through events or disabled when Workflow V3 owns retries.
- Workflow V3 retries the logical task according to explicit policy.
- Every provider contact receives its own operation ID.
- An operation is never overwritten by its retry.
- Researchctl retries an entire runner attempt only when the runner/workflow boundary fails.

Prefer one visible retry owner. Hidden provider SDK retries make timing and request authority inaccurate.

## Usage normalization

```text
provider reports usage -> accountingMode=actual
provider completed but reports no usage -> accountingMode=none or conservative
call failed after admission -> preserve reservation; usage only if observed
cache hit with no provider contact -> no external operation; emit cache metric
```

Never convert maximum token settings into actual usage. Cost requires an explicit price table version and provider/model identity; otherwise omit it.

## Privacy and credentials

Profiles are resolved in the worker host. Canonical plans carry only profile references and non-secret capability requirements. Artifacts and events exclude prompts, source text, generated text unless the RAG output artifact explicitly requires them, raw HTTP headers, API keys, URLs containing credentials, and arbitrary provider errors.

## Decisions

### Decision: adapters live in RAG-eval

- **Context:** Geppetto cannot depend on Scraper, and Scraper cannot know RAG request semantics.
- **Decision:** RAG-owned tasks adapt Geppetto calls to Workflow V3 operations.
- **Consequences:** No Geppetto change; adapters must track stable public APIs.
- **Status:** accepted.

### Decision: no hidden retry

- **Decision:** Disable or surface transport retries so every provider contact is represented.
- **Consequences:** Provider-specific behavior needs capability tests.
- **Status:** proposed pending provider audit.

## Implementation phases

1. Define descriptors and safe failure taxonomy.
2. Implement fake-provider generation adapter and all terminal outcomes.
3. Add usage normalization from Geppetto events.
4. Add embedding batch adapter and cache-hit distinction.
5. Add reranking adapter.
6. Add profile capability validation and fingerprints.
7. Add real-provider smoke tests behind explicit authorization.
8. Connect operation observations to Researchctl analysis.

## Test strategy

- success with actual usage;
- success without usage;
- timeout before response;
- cancellation;
- transport failure followed by Workflow retry;
- malformed provider result;
- embedding batch cardinality mismatch;
- cache hit without provider operation;
- rerank usage and ordering;
- secret canaries in credentials, prompt, and provider errors;
- process crash after begin and before finish, yielding recoverable unknown outcome;
- authority exhaustion before provider contact.

## Intern guidance

Begin with a fake Geppetto-facing interface and the operation state machine. Do not call a real provider until every failure branch closes the operation correctly. Keep raw provider result conversion separate from operation accounting. Compare `events.Usage` fields to the counters individually; do not assume all providers populate every field.

## Completion criteria

Generation, embedding, and reranking tasks execute through Geppetto, every actual provider contact is represented exactly once in the durable ledger, failed contacts remain visible, usage provenance is honest, secret-canary tests pass, and Geppetto itself is unchanged.

## Technology primer: provider call versus durable operation

A provider call is a transport interaction performed by Geppetto. A durable operation is Scraper's record that a side effect was authorized, started, and completed with a bounded outcome. The two are related but not identical. A cache hit may produce a RAG result without a provider call. A provider SDK may retry transport internally and perform multiple calls unless configured otherwise. A process may crash after the provider receives a request but before the workflow records completion.

The adapter makes these distinctions explicit. It begins an operation immediately before provider contact, uses Geppetto for the call, and closes the operation for every observable terminal path. It emits separate cache and domain metrics when no provider contact occurs.

## Profiles and host configuration

Geppetto profiles combine model selection, provider configuration, credentials, and inference settings. Canonical workflow plans should carry a stable profile reference and required capabilities, not resolved secrets. The worker host loads profile sources and computes a safe fingerprint of behavior-affecting non-secret settings.

```text
plan: profileRef=ttc-generation
host: profile registry + secret stores
resolved: provider, model, temperature, limits, endpoint policy, credentials
recorded safe identity: profile slug + model + settings fingerprint
never recorded: API key, bearer token, raw secret source
```

If a profile changes, its fingerprint changes. That prevents two runs with materially different model settings from appearing identical.

## Operation state machine

```text
not-started
    |
    | BeginExternalOperation succeeds
    v
admitted -----------------------------+
    |                                  |
    | provider returns                 | process dies
    v                                  v
finishing                         durable unknown/open
    |
    +-> succeeded
    +-> failed
    +-> canceled
    +-> timed-out
    +-> unknown
```

`BeginExternalOperation` can fail before provider contact, for example because budget authority is exhausted. In that case there is no provider operation and the task fails with an admission error. Once begin succeeds, the adapter owns the obligation to finish or allow recovery to classify the open operation as unknown.

## Error classification

Raw provider errors may contain prompts, URLs, headers, account identifiers, or arbitrary server bodies. The ledger stores an allow-listed class and code:

| Raw condition | Class | Code |
|---|---|---|
| context deadline | timeout | `PROVIDER_TIMEOUT` |
| caller cancellation | canceled | `PROVIDER_CANCELED` |
| HTTP 429 | throttled | `PROVIDER_RATE_LIMITED` |
| transport reset | transport | `PROVIDER_TRANSPORT` |
| schema decode failure | invalid-response | `PROVIDER_INVALID_RESPONSE` |
| unknown error | unknown | `PROVIDER_UNKNOWN` |

Detailed diagnostics may go to protected operational logs under a separate retention policy, but not to portable experiment evidence.

## Generation walkthrough

The generation task first validates the RAG request and computes an upper-bound reservation from configured maximum input/output tokens and request count. It resolves the profile, begins an operation, invokes Geppetto, consumes provider events, validates structured output, and publishes a RAG representation artifact.

Usage extraction should prefer terminal provider events. If streaming emits intermediate usage updates, normalize them according to Geppetto's event contract rather than summing cumulative counters accidentally. Preserve cached-token fields when the standard counter vocabulary supports them; otherwise retain them in a versioned domain metric.

## Embedding walkthrough

An embedding batch has two cardinalities: input texts and returned vectors. The adapter records one provider request, input item count, reported embedding tokens, model dimensions, and elapsed time. It validates vector count and dimensions before closing the task successfully. If the provider call succeeds but result validation fails, the external operation is still a succeeded provider interaction while the workflow task fails with an invalid domain result. Those outcomes must not be collapsed.

```text
provider operation outcome: succeeded
RAG task outcome: failed (dimension mismatch)
```

This distinction is essential for cost and reliability analysis.

## Reranking walkthrough

Reranking sends a query and candidate documents and receives scored ordering. The adapter records provider contact and usage, while the RAG artifact records candidate identities and scores. Source text is not copied into operation metadata. Validate that every returned index refers to an input candidate and that duplicate/missing candidates follow explicit policy.

## Crash-window analysis

List each durable boundary and ask what a restart sees:

1. Crash before begin: no provider contact and no operation.
2. Crash after begin but before request: admitted operation may be recovered as unknown, although no contact occurred.
3. Crash after request but before response: provider effect may have occurred; operation is unknown.
4. Crash after response but before finish: usage may be lost unless events are durably captured; operation is unknown/conservative.
5. Crash after finish but before task completion: operation is closed; task may retry and create another operation.

An idempotency key derived from run/node/attempt/operation identity can reduce duplicate provider work when supported, but it does not replace durable recording.

## First-week route

Implement a fake provider that can produce every terminal condition and usage shape. Build a table-driven adapter state-machine test. Add secret canaries to prompts and errors. Then integrate one Geppetto fixture engine. Only after all branches pass should you run an explicitly authorized real-provider smoke.

Useful source tests include Geppetto's JS/event usage tests and RAG provider adapter tests. Run them without changing Geppetto:

```bash
cd /home/manuel/code/wesen/go-go-golems/geppetto
go test ./pkg/events/... ./pkg/observability/... ./pkg/embeddings/... ./pkg/rerank/... -count=1
cd /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system
go test ./pkg/ragproviders/geppetto/... ./pkg/ragoperators/... -count=1
```

## Common mistakes

- Finishing the operation only on success leaves failures invisible.
- Marking the provider operation failed when provider output succeeded but RAG validation failed corrupts provider reliability metrics.
- Counting configured maximum tokens as actual usage invents data.
- Recording cache hits as provider operations overstates requests.
- Allowing an SDK to retry invisibly defeats request authority.
- Hashing secrets into a profile fingerprint can leak information through offline guessing.

## Adapter API reference

Keep provider-facing interfaces narrow enough to fake:

```go
type Generator interface {
    Generate(context.Context, GenerationRequest) (GenerationResponse, error)
}

type UsageEvidence struct {
    InputTokens, OutputTokens, CachedTokens int64
    CostMicrounits                          int64
    Availability                           string // actual, partial, unavailable
}
```

The production implementation wraps Geppetto. The operation decorator wraps this narrow interface and receives a clock plus `ExternalOperationRecorder`. Tests can therefore force timing and every failure window. Domain result validation remains outside the provider interface so a successful call with invalid RAG output can be represented correctly.

Operation completion errors are serious. If the provider succeeds but `FinishExternalOperation` fails, the task must not publish a successful domain artifact and ignore custody failure. Return a typed infrastructure error so Workflow policy can stop or retry while the open operation remains diagnosable.

## Review exercise

Write a table with provider outcome, usage availability, RAG validation outcome, operation outcome, and task outcome for at least eight combinations. If two rows cannot be represented without lying about one field, revise the adapter contract before implementation.

## Intern onboarding checklist

The engineer should resolve a test profile, identify safe and secret fields, draw the operation state machine, classify each fake error, explain provider-operation versus task outcome, and demonstrate a failed first provider contact followed by a successful Workflow V3 retry with two durable operations.

## References

- Program: `EXPERIMENT-PLATFORM-CONVERGENCE`.
- Geppetto: `pkg/engineprofiles/`, `pkg/inference/`, `pkg/embeddings/`, `pkg/rerank/`, `pkg/events/`.
- Scraper: `pkg/workflowv3/external_operation.go`.
- RAG: `pkg/ragproviders/geppetto/`, `pkg/ragoperators/`.
- Depends on `RAG-V2-WORKFLOW-LOWERING` and `SCRAPER-WORKFLOW-OBSERVATIONS`.
