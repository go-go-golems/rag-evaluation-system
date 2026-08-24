---
Title: Investigation diary
Ticket: RAG-GEPPETTO-WORKFLOW-OPERATIONS
Status: active
Topics:
    - rag
    - rag-eval
    - workflow
    - geppetto
    - embeddings
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-22T19:25:21.485468617-04:00
WhatFor: ""
WhenToUse: ""
---

# Investigation diary

## Goal

Record the evidence and decisions used to create the intern-facing design for **RAG-GEPPETTO-WORKFLOW-OPERATIONS** and place it in the wider experiment-platform convergence program.

## Step 1: Establish the ticket and implementation guide

This step created the ticket as a focused child of the multi-repository convergence program. The guide explains the existing systems, proposed contracts, pseudocode, diagrams, implementation phases, tests, deletion rules, and neighboring tickets so an intern can begin without reconstructing the architecture from Git history.

No production code was changed. The architecture is proposed work; current files are referenced as evidence and implementation entry points.

### Prompt Context

**User prompt (verbatim):** "Create all the tickets with docmgr, then for each ticket, Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.


That way we will have everything. Have all the tickets reference each other so that if someone stumbles upon one ticket, they'll know what this is all about."

**Assistant interpretation:** Create the umbrella and all directed tickets, write self-contained technical guides, cross-link them, validate docmgr state, and publish the complete set to reMarkable.

**Inferred user intent:** Preserve a coherent program plan before more ad-hoc implementation occurs, while making every workstream independently discoverable and approachable by a new engineer.

### What I did
- Created `RAG-GEPPETTO-WORKFLOW-OPERATIONS` with docmgr in the rag-eval repository.
- Added a long-term design/implementation guide and this diary.
- Mapped the ticket to all nine sibling/program tickets.
- Grounded the guide in current repository packages and public contracts.

### Why
- The convergence crosses repository boundaries and needs explicit ownership and sequencing.
- Standalone guides reduce the chance that an implementer recreates generic functionality in a workload package.

### What worked
- Existing greenfield contracts provided concrete API and file anchors.
- The ticket scope could be expressed as one independently testable capability: Geppetto generation, embedding, and reranking operation adapters.

### What didn't work
- N/A. No implementation or runtime validation was attempted in this documentation-only step.

### What I learned
- The repositories are closer to the target internally than the prior TTC path suggested; the largest deficit is the canonical integration route.

### What was tricky to build
- The guides must be self-contained without duplicating ownership. This was handled by giving each ticket one primary capability, linking the complete program map, and making predecessor/successor relationships explicit.

### What warrants a second pair of eyes
- Public schema names, package deletion lists, and the exact boundary between RAG factor expansion and Researchctl scheduling should be reviewed before implementation.

### What should be done in the future
- Review and accept the design, then implement phases in dependency order with the ticket's acceptance tests.

### Code review instructions
- Begin with the ticket's `design-doc/` guide and its “Program context” section.
- Validate documentation with `docmgr doctor --ticket RAG-GEPPETTO-WORKFLOW-OPERATIONS --stale-after 30`.

### Technical details
- Program umbrella: `EXPERIMENT-PLATFORM-CONVERGENCE`.
- Ticket responsibility: Geppetto generation, embedding, and reranking operation adapters.

## Step 2: Rewrite the briefs as textbook-style intern guides

The first drafts were structurally correct but too compressed. They named responsibilities and implementation phases without spending enough time teaching the foundational concepts, tracing real code, or showing how an intern would validate each claim. This revision expanded every guide into a teaching document rather than leaving it as an architecture summary.

The revised guides now introduce the domain vocabulary before APIs, walk through current source packages, show concrete identity and failure traces, provide pseudocode and data structures, explain common mistakes, and end with onboarding exercises. The ticket-research workflow remains responsible for evidence, bookkeeping, validation, and reMarkable publication; the textbook-authoring style governs the prose.

### Prompt Context

**User prompt (verbatim):** "these guides are much too concise / terse. Read the skill on writing guides for interns, and update them."

**User prompt (verbatim):** "write in a textbook authoring style as well"

**User prompt (verbatim):** "it's ticket-research-docmgr-remarkable skill"

**Assistant interpretation:** Replace concise RFC-style briefs with long-form, foundational, intern-oriented teaching guides while continuing to follow the ticket research, docmgr, diary, validation, and publication workflow.

**Inferred user intent:** Ensure the ticket set can genuinely onboard implementers without relying on the prior conversation or requiring them to rediscover the systems.

### What I did
- Loaded the full-blown intern research/design and textbook-authoring guidance, while retaining the pinned ticket-research workflow as the delivery contract.
- Expanded conceptual foundations, code tours, worked traces, API references, implementation routes, debugging guidance, common mistakes, and onboarding checklists.
- Preserved program navigation in every ticket.

### Why
- A list of recommendations does not teach the invariants that an implementer must preserve.
- The program spans three repositories and several identity/lifecycle layers; foundational explanation is required before implementation advice.

### What worked
- Existing canonical contracts and integration tests provide concrete teaching anchors.
- Textbook-style worked flows make responsibility boundaries easier to evaluate than responsibility tables alone.

### What didn't work
- The initial drafts were approximately architecture-brief length and did not meet the requested exhaustive intern-guide standard.

### What I learned
- Each guide needs both a local implementation path and the complete cross-program context; either one alone is insufficient.

### What was tricky to build
- Expansion had to add depth without duplicating ownership or introducing analogies. The solution was to teach using exact identities, state transitions, source files, schemas, and failure windows.

### What warrants a second pair of eyes
- Review whether each guide provides enough real source-level orientation for an intern and whether proposed package names should be accepted before implementation.

### What should be done in the future
- Review the expanded guides, then implement in dependency order. Add experiments to each ticket as its implementation starts.

### Code review instructions
- Read the technology-primer and guided-source-tour sections before reviewing proposed APIs.
- Check that all bullet points are complete technical statements and diagrams preserve actual system boundaries.

### Technical details
- Style: foundational prose, concrete examples, no analogies, precise diagrams, worked traces, and executable validation routes.

## Step 3: Validate and publish the program guide bundle

The ticket passed docmgr validation and its guide was included in the ordered ten-document program bundle. The dry run confirmed document ordering and destination; the real upload completed successfully.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Complete ticket bookkeeping and deliver all guides as one navigable PDF.

**Inferred user intent:** Make the complete program available as a durable reading package rather than scattered repository files.

### What I did
- Ran `docmgr doctor --ticket <ticket> --stale-after 30`.
- Performed a reMarkable bundle dry run with all ten guides in dependency order.
- Uploaded `Scriptable Experiment Platform Intern Guides` to `/ai/2026/07/22/EXPERIMENT-PLATFORM-CONVERGENCE`.

### Why
- A single bundle provides a table of contents across repository boundaries.

### What worked
- All ticket doctor checks passed and the upload reported success.

### What didn't work
- N/A.

### What I learned
- Cross-repository tickets remain locally owned while one ordered publication can present the program coherently.

### What was tricky to build
- Document order needed to follow dependency order rather than repository or alphabetical order.

### What warrants a second pair of eyes
- Review the PDF at normal reMarkable size for code-block and diagram readability.

### What should be done in the future
- Re-upload with `--force` only after accepted content revisions, because overwrite removes annotations.

### Code review instructions
- Verify ticket frontmatter, program links, design guide, tasks, and changelog.

### Technical details
- Remote destination: `/ai/2026/07/22/EXPERIMENT-PLATFORM-CONVERGENCE`.

## Step 4: Establish the durable provider-operation state machine

This step implemented the first production slice for the active operations phase: a RAG-owned decorator layer that wraps generation, embedding, and reranking provider interfaces with Scraper Workflow V3 external-operation custody. The slice is deliberately independent of concrete Geppetto construction so every state transition, counter, failure class, and privacy boundary can be tested before linking host profiles or making a real provider call.

The implementation creates exact immutable descriptors for `provider.generate/v1`, `provider.embed/v1`, and `provider.rerank/v1`; admits each call before provider contact; records only bounded counters and safe correlation identities; finishes successful, failed, canceled, and timed-out contacts; detaches completion persistence from caller cancellation; and distinguishes a successful provider interaction followed by invalid domain observation from a failed provider interaction.

### Prompt Context

**User prompt (verbatim):** "RAG-GEPPETTO-WORKFLOW-OPERATIONS"

**Assistant interpretation:** Implement the next convergence phase completely, beginning with the lowest-risk provider-operation custody contract and deterministic tests.

**Inferred user intent:** Add real Geppetto-backed RAG operations without compromising the Workflow V3 lifecycle, operation ledger, provider privacy, or the accepted provider-free path.

**Commit (code):** `eef610d` — "feat: add durable provider operation decorators"

### What I did
- Re-read the accepted operations guide, task list, existing diary, completed lowering package, RAG provider adapters, Geppetto provider interfaces, and Scraper external-operation contracts/store implementation.
- Added `pkg/ragworkflowops` with:
  - exact versioned generation, embedding, and reranking descriptors;
  - canonical authority digests and safe correlation digests that omit prompt/source/query text;
  - bounded measure, reservation, and usage counter vocabularies;
  - configurable per-attempt admission and finish timeout;
  - generation, embedding, and reranking decorators;
  - admission-before-contact and completion-after-contact sequencing;
  - UTC provider timing and elapsed microseconds;
  - actual, conservative, and unavailable accounting modes;
  - allow-listed cancellation, timeout, and transport failures;
  - a typed `ProviderResultError` for successful provider calls whose domain result cannot be accepted.
- Added fake-recorder/provider tests for actual token/cost usage, cardinality, secret-bearing payload exclusion, cancellation with detached completion context, conservative reservation accounting, admission failure before contact, and provider-success/domain-validation-failure separation.
- Ran focused tests and focused lint, then committed through the repository pre-commit hook, which passed package/internal tests and lint.

### Why
- Provider contact must never occur before durable authority is admitted.
- Cancellation must not prevent the adapter from recording the outcome of an already admitted call.
- Provider reliability and RAG result validity are different observations. A valid provider response with unusable dimensions/schema/cost evidence must not be falsely recorded as a transport failure.
- RAG request payloads can contain prompts, corpus text, query text, evidence, credentials, or arbitrary provider errors; none belong in the portable operation ledger.

### What worked
- All focused state-machine tests passed on the first run.
- Generation evidence recorded one request, exact input/output token counts, and deterministic cost microunits.
- A canceled call recorded `canceled/PROVIDER_CANCELED` with conservative accounting when a request reservation existed, and `FinishExternalOperation` received a live detached context.
- Admission failure prevented provider invocation entirely.
- Invalid post-call cost evidence left the operation outcome `succeeded` while returning a typed domain-result error to the task.
- The full pre-commit package/internal suite and lint passed.

### What didn't work
- N/A. No command or test failed in this slice.

### What I learned
- Scraper validates reservations against attempt-owned Workflow budget reservations, so production lowering must add matching budget envelopes before enabling non-empty provider reservations.
- `FinishExternalOperation` validates the completion ticket rather than the live lease, specifically allowing post-cancellation outcome persistence; the adapter must still use a bounded context detached from caller cancellation.
- Existing Geppetto embedding interfaces do not expose usage, while generation does and Geppetto reranking does but the current RAG reranker adapter drops it. Reranker result/usage propagation needs a deliberate contract update in the next slice.

### What was tricky to build
- External-operation counter descriptors require counters, roles, and runtime values to be strictly sorted. The decorator builds one closed vocabulary and sorts every emitted set before store validation.
- Cost is provider-reported floating currency in the existing RAG interface but ledger authority is integer-only. The decorator converts to integer microunits with finite, non-negative, overflow-safe validation.
- A canceled caller context cannot be reused for completion persistence. The implementation uses `context.WithoutCancel` plus a required bounded timeout; it does not create an unbounded background write.
- A post-call validation error needs two truths: operation succeeded, task failed. `ProviderResultError` preserves that distinction without placing arbitrary error text in the ledger.

### What warrants a second pair of eyes
- Review whether currency-to-microunit rounding should remain nearest-integer or use a provider/pricing-specific exact rational conversion before real billing evidence is accepted.
- Review reservation sizing once canonical operator configs and host profile limits are connected to Workflow budgets.
- Review the generic fallback classification of untyped provider errors as `transport/PROVIDER_TRANSPORT`; concrete Geppetto adapters should return typed safe classifications for rate limits and provider 5xx responses.
- Confirm safe correlation fields are sufficient for scientific linkage without hashing secret-bearing payloads.

### What should be done in the future
- Propagate rich rerank usage/cost through the RAG adapter contract.
- Make Geppetto adapters return typed provider-contact outcomes for invalid responses and allow-listed failures.
- Add the Geppetto task package/module, host provider loading, matching plan budgets, and operation descriptors to `pkg/ragworkflow`.
- Prove the decorators against Scraper's real SQLite recorder, restart, retry, crash-window, and exported operation manifest.

### Code review instructions
- Start at `pkg/ragworkflowops/operations.go`, especially `NewDescriptors`, `executeOperation`, and `completionFor`.
- Review `operations_test.go` as the state-transition table.
- Validate with `GOWORK=off go test ./pkg/ragworkflowops -count=1` and `GOWORK=off golangci-lint run ./pkg/ragworkflowops/...`.

### Technical details
- Operation kinds: `provider.generate/v1`, `provider.embed/v1`, `provider.rerank/v1`.
- Safe counters: requests, input/output/embedding tokens, input/output items, and cost microunits.
- Completion writes use `context.WithoutCancel` constrained by `Policy.FinishTimeout`.
- Provider payloads and raw errors are absent from operation descriptors, specs, and completions.

## Step 5: Preserve reranking usage and cost through the RAG contract

The provider-operation decorator could initially custody reranking cardinality but not provider-reported token usage or cost because the existing RAG `Reranker` interface returned only scores. Geppetto's rerank response already distinguishes absent usage/cost from explicit zero values, so this step widened the internal RAG result contract to preserve that evidence instead of fabricating or dropping it.

The native rerank operator now accumulates provider-reported input tokens and cost into the same RAG usage model used by generation. The Geppetto adapter maps rich response evidence into `RerankResult`, and the Workflow operation decorator emits exact integer token and cost-microunit counters.

### Prompt Context

**User prompt (verbatim):** (same active goal as Step 4)

**Assistant interpretation:** Continue with the next concrete gap discovered by the operation audit: honest reranking usage custody.

**Inferred user intent:** Ensure every supported Geppetto operation has scientifically honest usage and cost evidence before production Workflow integration.

**Commit (code):** `c983985` — "feat: preserve rerank provider usage"

### What I did
- Changed `ragoperators.Reranker` to return `RerankResult` containing scores, optional input-token usage, and optional cost.
- Updated the native rerank operator to accumulate reported usage and preserve nil-versus-zero cost semantics.
- Updated the Geppetto rerank adapter to retain `Response.Usage.InputTokens` and `Response.Cost` after validating complete score coverage.
- Updated the concurrency-limiting adapter without changing its admission behavior.
- Extended durable Workflow operation projection to emit rerank `input_tokens` and `cost_microunits` counters.
- Updated all fake rerankers and tests, including native operator usage accumulation, adapter propagation, and operation-ledger counters.
- Ran focused operator/provider/operation tests and lint, then committed through the full package/internal pre-commit tests and lint.

### Why
- Dropping available usage makes provider cost and token measurements incomplete.
- Treating nil cost as zero would falsely assert that an unpriced call was free.
- Scraper operation counters require bounded integers, so provider floating cost must be validated and converted at the RAG adapter boundary.

### What worked
- Geppetto adapter tests now prove 41 input tokens and non-zero cost survive response mapping.
- Native operator tests prove rerank usage reaches `Environment.Usage` under the configured model identity.
- Workflow decorator tests prove the ledger receives sorted request, output-item, input-token, and cost-microunit counters.
- Focused and pre-commit validation passed.

### What didn't work
- The first compile after the interface change found two expected migration defects:
  - `NewReranker` had been mechanically changed to return `ragoperators.RerankResult{}` even though its constructor returns `*Reranker`;
  - the native operator test fake still implemented the old `[]RerankScore` signature.
- Exact compiler errors included `cannot use ragoperators.RerankResult{} ... as *Reranker value` and `fakeReranker does not implement Reranker (wrong type for method Rerank)`. I corrected the constructor return, migrated the fake, reran the affected packages, and then added explicit usage assertions.

### What I learned
- Geppetto's rerank API already has the correct nil-versus-zero evidence model; the loss occurred solely in the RAG adapter contract.
- Interface migration is preferable to side-channel usage because scores and their usage belong to one provider response and one operation completion.

### What was tricky to build
- Usage must be accumulated under the canonical configured model reference while Geppetto validates and reports the exact model. The adapter rejects model mismatch before exposing the result, and the operator attributes accepted usage to its immutable config.
- Operation completion must convert only present cost. A nil pointer emits no counter; an explicit zero pointer emits a valid zero semantic value, though zero counters are omitted from storage under the current compact counter representation.

### What warrants a second pair of eyes
- Review whether total rerank tokens should be retained separately when providers report both input and total values; the current RAG measure vocabulary records input tokens only.
- Review the cost currency assumption before combining costs from heterogeneous providers.

### What should be done in the future
- Add typed Geppetto failure/contact outcomes.
- Connect operation-decorated services to the provider Workflow task package and host authority identity.

### Code review instructions
- Start with `pkg/ragoperators/types.go` and `rank.go`, then inspect `pkg/ragproviders/geppetto/reranker.go` and `pkg/ragworkflowops/operations.go`.
- Run `GOWORK=off go test ./pkg/ragoperators ./pkg/ragproviders/... ./pkg/ragworkflowops -count=1`.

### Technical details
- Rich result: `ragoperators.RerankResult`.
- Preserved evidence: score list, input tokens, optional provider cost.
- Workflow cost unit: integer microunits.

## Step 6: Bind provider authority into a production Workflow task package

This step connected the operation decorators to the accepted RAG Workflow backend. It added a separate `rag-v2-geppetto@1.0.0` task package whose immutable bundle contains a privacy-safe provider authority document, whose trusted task module exposes the exact operation descriptors, and whose runtime environment decorates the host's generator, embedder, and reranker for every task attempt.

The provider lowerer now accepts the exact provider-required operator entries that the provider-free lowerer continues to reject. Provider authority also participates in the preparation fingerprint, preventing prepared representations, embeddings, and indexes from being reused across materially different profile/model/settings identities.

### Prompt Context

**User prompt (verbatim):** (same active goal as Step 4)

**Assistant interpretation:** Connect the tested operation state machine to the actual RAG Workflow task package without weakening provider-free behavior or generic Scraper boundaries.

**Inferred user intent:** Make provider-backed plans executable through the production Workflow V3 engine with immutable host authority and durable operation evidence.

**Commit (code):** `c6e332f` — "feat: add Geppetto Workflow task package"

### What I did
- Refactored the RAG task runtime to receive an environment factory and preparation implementation identity while retaining the provider-free factory unchanged.
- Generalized task-bundle construction so provider-free and Geppetto packages share task semantics but have separate package identities and bundle bytes.
- Added strict `ProviderAuthority` and `WorkflowProviderIdentity` contracts containing safe profile/model/settings/pricing/concurrency identity while excluding profile source paths and credentials.
- Added authority canonicalization, ordering, cardinality, SHA-256, concurrency, pricing, and self-digest validation.
- Added `ProviderServicesFromSet` to adapt the existing host-loaded `ragproviders.ProviderSet` without moving host configuration into plans.
- Added `ProviderPackage`, operation descriptors, task module factory, provider environment creation, and immutable `provider-authority.json` bundle content.
- Added `NewProviderLowerer` and `BuildProviderRunnerExecution`; the default lowerer still fails provider-required operators closed.
- Changed preparation fingerprints to include either `fixture-embedding/v1` or the exact provider authority digest.
- Added an SQLite-backed integration test that:
  - proves provider-free rejection and provider lowerer acceptance;
  - executes structured generation and embedding through real Workflow tasks;
  - records complete generation/embedding operations with the exact authority digest;
  - proves no rerank operation occurs when the pipeline has no rerank node;
  - compares Workflow retrieval results with direct `ragengine` execution.
- Ran focused tests/lint and the full package/internal pre-commit suite.

### Why
- Operation descriptors alone do not bind host behavior into a plan. Embedding `provider-authority.json` in the bundle changes package/catalog/plan identity whenever provider behavior-affecting identity changes.
- Host credentials and endpoint details must remain worker-local, while model/prompt/settings/pricing/concurrency identities must remain auditable.
- Provider-backed preparation cannot reuse artifacts created under the deterministic fixture implementation or a different provider profile.

### What worked
- Provider-free regression tests remained green.
- The provider package integration produced succeeded `provider.generate/v1` and `provider.embed/v1` ledger rows with the exact authority digest.
- Direct `ragengine` and Workflow execution produced equal first-query retrieval results.
- The task module's descriptors were resolved by Scraper's ordinary module registry and persisted by the ordinary SQLite external-operation recorder; no RAG-specific Scraper changes were needed.
- Full pre-commit tests and lint passed.

### What didn't work
- The first provider integration run completed Workflow execution but failed during direct parity setup with `operator representations.structured-summary/v1: RAG_GENERATOR_UNAVAILABLE: representations.structured-summary/v1`.
- Cause: the direct `ragengine.Options` supplied manifests and embedding but omitted generator, schemas, reranker, and generation fingerprint. This was a parity-harness defect, not a Workflow defect. I supplied the complete environment through `ragengine.Options` and reran both provider and provider-free integration tests successfully.

### What I learned
- Scraper's task package catalog does not directly include operation descriptor contents, so authority must also participate in immutable bundle bytes. The embedded authority file supplies that binding.
- Preparation identity must be injected into the task runtime as well as the lowerer; otherwise the compile-time and runtime fingerprints diverge.
- The existing `ProviderSet.CapabilityDescriptor` includes profile source information. The Workflow authority intentionally reconstructs a safe subset and excludes that path.

### What was tricky to build
- The same task JavaScript and Go runtime must serve provider-free and provider-backed packages without allowing both module factories in one selected package set. Separate package selection with a shared module alias is safe because one run selects exactly one package.
- Authority slices and provider identities require canonical sorting before digesting. Validation also rejects a caller that supplies equivalent but noncanonical ordering.
- The provider package must keep shared caches outside operation decoration: cache hits happen in native RAG operators before decorated provider contact, so they correctly create no operation row.

### What warrants a second pair of eyes
- Review whether every field in `WorkflowProviderIdentity` affects provider behavior or authority; remove non-semantic fields and add missing semantic fields before schema freeze.
- Confirm `provider-authority.json` is sufficient package/catalog binding under future Scraper bundle changes.
- Review default `MaxPerAttempt=10,000`; production lowering should derive tighter per-task bounds from batch plans.
- Review pricing currency/version identity, which is not yet explicit beyond configured microunit rates.

### What should be done in the future
- Add canonical Workflow budgets and reservations derived from provider configs and map/reduction cardinalities.
- Add provider-backed fixtures covering synthetic/combined generation, answer generation, and reranking.
- Add typed Geppetto error/contact outcomes, malformed-response distinction, retries, crash windows, and exported operation-manifest acceptance.
- Add a provider-configured runner command and Researchctl multi-case plan.

### Code review instructions
- Start at `pkg/ragworkflow/provider_package.go`, then review the environment-factory changes in `runtime.go`, lowerer mode in `lower.go`, and the SQLite integration in `provider_package_test.go`.
- Run `GOWORK=off go test ./pkg/ragworkflow -count=1`.

### Technical details
- Provider package: `rag-v2-geppetto@1.0.0`.
- Authority schema: `rag-workflow-provider-authority/v1`.
- Bundle identity file: `provider-authority.json`.
- Preparation reuse identity: provider authority digest.

## Step 7: Separate preflight, provider contact, and domain validation

This step removed a lifecycle ambiguity in the first decorator design. Provider adapters previously performed request/profile/schema validation inside the decorated `Generate`, `Embed`, or `Rerank` call, after operation admission. A deterministic configuration error could therefore create an admitted operation even though no provider contact occurred. The adapters now expose optional prepared-call interfaces: preflight completes before admission, provider contact executes after admission, and post-contact result validation reports a successful operation plus a failed task when appropriate.

The same change introduced a closed, safe Geppetto failure taxonomy. Raw provider errors remain available only as wrapped in-process causes for `errors.Is`; portable task errors and ledger completions expose allow-listed codes such as `PROVIDER_RATE_LIMITED`, `PROVIDER_SERVER_ERROR`, `PROVIDER_TIMEOUT`, and `PROVIDER_CANCELED`.

### Prompt Context

**User prompt (verbatim):** (same active goal as Step 4)

**Assistant interpretation:** Continue hardening provider lifecycle truth before adding broader workload acceptance.

**Inferred user intent:** Ensure the operation ledger represents actual provider contact and never confuses configuration failures, transport failures, and post-response domain failures.

**Commit (code):** `6451850` — "feat: classify Geppetto operation outcomes"

### What I did
- Added optional generation, embedding, and reranking preflight/prepared-call interfaces to `pkg/ragworkflowops`.
- Changed decorators to execute preflight before `BeginExternalOperation` and only execute prepared provider contact after admission.
- Refactored the Geppetto generation adapter into `PrepareGeneration` and `preparedGeneration.Execute`.
- Refactored embedding and reranking adapters similarly, cloning bounded request inputs during preflight.
- Marked empty/malformed generation responses, embedding cardinality/dimension/non-finite results, and rerank model/cardinality/identity/score failures as provider-success/domain-result failures.
- Replaced arbitrary result-error construction with a closed `RAG_*` code validator; arbitrary text maps to `RAG_PROVIDER_RESULT_INVALID`.
- Added `ProviderCallError` with validated Workflow failure class, code, outcome, safe `Error()`, and raw-cause `Unwrap()`.
- Added Geppetto classification for cancellation, timeout, 429, provider 5xx, and transport failures without copying raw bodies into evidence.
- Added tests proving preflight failure creates no operation/contact, arbitrary result text is redacted, and classified rate limits produce only safe ledger taxonomy.
- Updated existing adapter tests to assert safe provider codes and `errors.Is` cancellation.
- Ran focused and full pre-commit package/internal tests and lint.

### Why
- A normal preflight rejection is not an external effect and must not appear as provider contact.
- Provider response validation happens after a billable/reliability-relevant contact; recording it as transport failure would corrupt provider success rates.
- Raw provider errors may contain response bodies, URLs, account identifiers, prompts, or credentials and cannot cross the portable evidence boundary.

### What worked
- Preflight rejection test proved zero operation specs and zero provider calls.
- Provider-result code validation replaced a secret canary with `RAG_PROVIDER_RESULT_INVALID`.
- Classified rate-limit test recorded `rate-limit/PROVIDER_RATE_LIMITED` while returning only the safe code.
- Existing live-shape HTTP generation tests preserved cancellation via `errors.Is` and proved a 502 body was absent from the returned error.
- Full pre-commit validation passed.

### What didn't work
- The first focused run after introducing typed result errors failed five existing adapter assertions because `ProviderResultError.Error()` initially returned only the generic `RAG_PROVIDER_RESULT_INVALID`, hiding stable adapter codes such as `RAG_GEPPETTO_EMBED_DIMENSIONS`.
- Exact symptoms included `error=RAG_PROVIDER_RESULT_INVALID want=RAG_GEPPETTO_EMBED_COUNT` and rerank tests expecting stable `RAG_GEPPETTO_RERANK` errors.
- I replaced arbitrary-error construction with a closed code-only constructor. Stable safe RAG codes remain visible; unapproved text is reduced to the generic code. Focused tests then passed.
- After provider classification, two tests still expected historical strings (`RAG_GEPPETTO_GENERATOR_PROVIDER` and literal `context canceled`). I updated them to the accepted safe taxonomy (`PROVIDER_SERVER_ERROR`) and `errors.Is(context.Canceled)` plus `PROVIDER_CANCELED`.

### What I learned
- Stable error detail and privacy are compatible when the error contract accepts only validated codes rather than arbitrary wrapped strings.
- Geppetto currently exposes some HTTP status information only through error text. Classification may inspect that text in-process, but only a closed class/code leaves the adapter.
- Prepared-call interfaces can be optional, preserving simple deterministic fixture providers while giving real adapters an exact admission boundary.

### What was tricky to build
- Post-contact validation needs provider usage in the operation completion even though the task returns an error. Prepared generation returns the partially populated usage result alongside `ProviderResultError`; the decorator derives counters before handling the error and records a succeeded operation.
- Cancellation needs both a safe outward code and `errors.Is` compatibility. `ProviderCallError` exposes only its code from `Error()` but unwraps the original context error in-process.
- Request preflight must copy slices so caller mutation cannot change the prepared effect after authority admission.

### What warrants a second pair of eyes
- Review string-based HTTP status classification; a future Geppetto typed transport error should replace this local parsing.
- Confirm all post-response validation branches return `ProviderResultError`, especially any future provider adapters.
- Review whether preflight can perform hidden network I/O through engine/profile construction; current Geppetto engine factory is expected to be local configuration only.

### What should be done in the future
- Add direct concrete-Geppetto decorator tests proving malformed HTTP responses produce succeeded operation rows with preserved usage where available.
- Add Workflow retry tests yielding separate durable operations per contacted attempt.
- Add budget reservations and authority exhaustion before contact.

### Code review instructions
- Review the prepared-call interfaces and `completionFor` in `pkg/ragworkflowops/operations.go`.
- Review the split adapters in `pkg/ragproviders/geppetto/{generation,embeddings,reranker}.go` and safe classifier in `errors.go`.
- Run `GOWORK=off go test ./pkg/ragworkflowops ./pkg/ragproviders/geppetto -count=1`.

### Technical details
- Safe result-code pattern: `^RAG_[A-Z0-9_]{3,63}$`.
- Failure classes: canceled, timeout, rate-limit, provider-5xx, transport.
- Preflight happens before operation admission; prepared execution happens after admission.

## Step 8: Compile provider reservations into Workflow budgets

This step connected external-operation reservations to Workflow V3's authoritative budget system. Provider operations now reserve one request by default before contact, provider task specifications declare immutable maxima, provider IR nodes/maps carry matching claims, and the run plan contains an authority-digest-bound provider account large enough for every bounded materialization allowed by the plan.

The implementation derives all claims from the same operation policy used by the runtime decorator. It checks multiplication and aggregation overflow rather than saturating or relying on hidden assumptions.

### Prompt Context

**User prompt (verbatim):** (same active goal as Step 4)

**Assistant interpretation:** Add the missing authority layer so provider operation admission is backed by Workflow-owned budgets rather than descriptors alone.

**Inferred user intent:** Make provider calls impossible when their bounded request/token/cost authority has not been reserved by the current task attempt.

**Commit (code):** `4d28123` — "feat: reserve provider operation budgets"

### What I did
- Added default per-call request reservations for generation, embedding, and reranking.
- Added provider task budget maxima to the immutable task bundle:
  - representation generation: up to descriptor `MaxPerAttempt`;
  - embedding batches: up to descriptor `MaxPerAttempt`;
  - one query item: at most one generation, embedding, and rerank contact.
- Added exact IR budget claims to provider representation/embedding nodes and query maps.
- Added a run-level `provider` budget account whose limits cover static provider nodes plus every query item allowed by `MaxItems`.
- Bound the budget policy digest to the exact provider authority digest.
- Added checked integer multiplication/summation and explicit errors for invalid counts or overflow; no panic/saturation path remains.
- Extended the SQLite integration test to prove:
  - account limit `requests=30,200` for the fixture graph;
  - query-map claim `requests=3`;
  - each durable provider operation has reservation `requests=1`.
- Reran focused tests/lint and full pre-commit validation.

### Why
- Scraper rejects external-operation reservations not already reserved by the active Workflow attempt. Descriptors authorize vocabulary and per-attempt cardinality; Workflow budgets authorize scalar consumption.
- A package maximum without an IR claim is not runtime authority.
- Query maps can materialize many task attempts, so the account maximum must incorporate bounded `MaxItems` rather than only static graph nodes.

### What worked
- The real SQLite provider integration succeeded with reservation enforcement enabled.
- Budget claims and operation allocations matched exactly.
- Provider-free plans retained no provider budget account or reservation behavior.
- Full pre-commit tests and lint passed.

### What didn't work
- N/A. The initial budget integration passed focused tests.

### What I learned
- Scraper has three distinct authority levels: descriptor counter roles/max calls, task-spec budget maximum, and run/attempt reservations. Provider integration must satisfy all three.
- Query-map budget accounts must be sized from map cardinality bounds, not observed fixture size, because admission occurs against the immutable plan.

### What was tricky to build
- One query task can contact three distinct provider operation kinds, but all reserve the same `requests` dimension. The task claim must sum those reservations while descriptor maxima remain per kind.
- Static representation and embedding tasks can perform batched contacts whose exact count is data-dependent. The immutable bound uses `MaxPerAttempt`; operation planning will tighten this when exact batch manifests are available.
- Multiplying task claims by `MaxItems` can overflow. The implementation rejects overflow before plan compilation.

### What warrants a second pair of eyes
- Review whether the default maximum of 10,000 contacts per preparation attempt is too permissive; fixture/real plan generation should set a tighter explicit value.
- Add output-token and cost reservations once profile pricing and maximum response tokens can be mapped to each task's exact possible calls.
- Confirm `fail-run` is the desired exhaustion policy for noninteractive Researchctl workloads.

### What should be done in the future
- Derive tighter generation batch counts and token/cost reservations from canonical operator plans and provider authority.
- Add budget-exhaustion acceptance proving provider contact count remains zero.
- Export budget usage alongside operation evidence in Researchctl acceptance.

### Code review instructions
- Review `taskBudgetMaximums`, `operationBudgetClaim`, and `applyBudgets` in `pkg/ragworkflow/provider_package.go`.
- Review task maximum plumbing in `package.go` and lowerer application in `lower.go`.
- Run `GOWORK=off go test ./pkg/ragworkflow -run TestProviderPackage -count=1`.

### Technical details
- Account: `provider`.
- Default reservation: `requests=1` per actual contact.
- Exhaustion policy: `fail-run`.
- Policy digest: provider authority digest.

## Step 9: Execute generation, embedding, reranking, and answers in one durable run

This step expanded the SQLite provider fixture from preparation-only contacts to the complete query-time provider path. The fixture pipeline now performs structured representation generation, preparation embeddings, query embeddings, cross-encoder reranking, and grounded answer generation for two queries, with direct `ragengine` parity for retrieval results and answers.

The expanded test exposed and fixed a real integration defect: the Workflow task runtime passed only manifests and embedding into `ragengine.Options`. Provider-free retrieval happened to work, but reranking and generation services were silently absent at query time. The runtime now passes the complete environment and exact provider fingerprints.

### Prompt Context

**User prompt (verbatim):** (same active goal as Step 4)

**Assistant interpretation:** Extend deterministic acceptance through all currently wired provider operation types and use failures to correct the production runtime.

**Inferred user intent:** Prove provider operation custody in realistic preparation and query phases, not only isolated decorators.

**Commit (code):** `dbd77c7` — "test: cover generation embedding and rerank operations"

### What I did
- Added a complete deterministic provider fixture that supplies:
  - structured summaries;
  - 32-dimensional embeddings;
  - complete rerank scores with input-token/cost evidence;
  - grounded answers with citation, token, and cost evidence.
- Added exact fixture model/prompt manifests for reranking and answer generation.
- Extended the canonical fixture pipeline with rerank and answer nodes and changed retrieval to the generated `summary` representation.
- Executed two query items through one Workflow run and asserted exact operation counts:
  - five generation contacts (three preparation summaries plus two answers);
  - three embedding contacts (one preparation batch plus two query embeddings);
  - two reranking contacts.
- Asserted every operation is succeeded, completed, authority-bound, and request-reserved.
- Compared Workflow retrieval results and answers with direct `ragengine` output.
- Fixed `taskRuntime.query` to pass schemas, generator, embedder, reranker, cache, concurrency, and exact generation/rerank/embedding fingerprints into `ragengine.Options`.
- Reran focused and full pre-commit validation.

### Why
- Provider operations span both preparation and query phases; testing only preparation cannot prove reranking, answer usage, query budgets, or query-time service composition.
- Direct semantic parity is required in addition to operation custody.

### What worked
- The final fixture run produced ten completed operations with exact expected kind counts and safe usage/cost counters.
- Query task claims covered exactly one embedding, rerank, and generation reservation per query.
- Direct and Workflow retrieval results and answers matched.
- Provider-free and all affected provider/operator tests remained green.

### What didn't work
- The first expanded direct parity run failed with `operator generate.answer/v1: RAG_ANSWER_FAILED: fixture answer evidence required` because retrieval still filtered for `raw` representations after the pipeline changed to generated `summary` representations. I changed BM25/vector retrieval config to `summary`; direct execution then succeeded.
- Workflow execution still failed after the query embedding operation and before reranking. Durable operation evidence showed a succeeded query embed but no rerank admission, while direct execution succeeded.
- Root cause: `taskRuntime.query` constructed `ragengine.Options` with only `Prepared`, `Manifests`, `Embedder`, and a hard-coded fixture embedding fingerprint. It omitted schemas, generator, reranker, cache, concurrency, and provider identities. I passed the complete environment and replaced the hard-coded fingerprint with the task runtime's preparation/provider identity. The run then succeeded with all ten operations.

### What I learned
- Provider-free tests can conceal missing query-time service wiring because vector retrieval only requires embedding.
- Durable operation ordering is useful failure localization: a succeeded query embedding followed by no rerank admission narrowed the defect to the in-process path between those contacts.
- Retrieval representation filters are semantic identity and must change with representation operators in fixture variants.

### What was tricky to build
- Rerank and answer nodes are dynamic and must preserve evidence output for evaluation while answer remains a side output. The pipeline output points to reranked evidence; `ragengine` separately collects answer values.
- Query budget claims aggregate three operation kinds into one request dimension, while operation counts remain distinct by descriptor.
- Provider fingerprints must be identical in lowerer preparation identity, task runtime, and direct parity options.

### What warrants a second pair of eyes
- Review the full `ragengine.Options` mapping in `taskRuntime.query` for future fields; a constructor/helper may prevent another silent omission.
- Confirm operation count expectations remain stable if batching or query-plan semantics change.
- Review whether answers should become a named terminal Workflow output in addition to their inclusion in `rag-workflow-result/v1`.

### What should be done in the future
- Add deterministic synthetic-question and combined-summary-question Workflow variants.
- Add retry, malformed response, cancellation, restart, and open-operation crash-window tests.
- Add built-binary runner/Researchctl acceptance and exported operation manifest checks.

### Code review instructions
- Review the complete provider fixture in `provider_package_test.go` and the corrected options mapping in `runtime.go`.
- Run `GOWORK=off go test ./pkg/ragworkflow -run TestProviderPackage -count=1`.

### Technical details
- Final operation counts: generate=5, embed=3, rerank=2.
- Query count: 2.
- Semantic parity: result traces and answer artifacts.

## Step 10: Productize deterministic fixtures and the provider runner

**User prompt (verbatim):** (same active goal as Step 4)

**Assistant interpretation:** Make provider-backed plans reproducible and executable through the production process runner, not only an in-process test.

**Commits:** `c0eaca1` — deterministic provider Workflow fixtures; `4e5bc61` — every provider attachment point.

### Work
- Added `NewDeterministicProviderServices`, a fixture-marked canonical provider authority, and deterministic answer/rerank evidence.
- Added `rag-workflow-provider-fixture`, two canonical cases, set-input archives, domain configs, and byte-stable direct parity files.
- Added `--provider-config`, `--provider-fixture`, and bounded provider-operation policy flags to `rag-workflow-runner`.
- Executed structured, combined-summary-question, synthetic-question, embedding, rerank, and answer attachment points through SQLite Workflow V3.
- Added exact operation-count and direct semantic-parity assertions.

### Failures and fixes
- Lefthook's Biome formatter rewrote canonical JSON after generation. Regeneration then differed. Added `examples/rag-workflow-provider` to the custody exclusion and amended the fixture commit.
- Synthetic-question execution initially failed with `RAG_INDEX_EMBEDDING_PARENT` because embeddings came from synthetic representations while the index still consumed summaries. Rebound both embed and index representation inputs to the synthetic node.
- Synthetic lowering changed the account request total from 30,200 to 30,300. Replaced a single hard-coded expectation with explicit per-variant budget evidence.

### Review
Regenerate fixtures and compare recursively. Review `provider_fixture.go`, `provider_export.go`, runner host configuration, and the three integration subtests.

## Step 11: Preserve retry taxonomy and cache-hit accounting

**Commits:** `4b93559` — safe retry classification; `d4a957c` — replay/cache/crash custody.

### Work
- Replaced opaque native Go errors at the JS task boundary with an allow-listed result envelope.
- Mapped typed provider transport/rate-limit/timeout failures to retryable Workflow failures and malformed provider results to permanent `malformed-output` failures.
- Proved one failed provider contact followed by one successful Workflow retry, with distinct operation IDs and stable request correlation.
- Added decorator request/usage snapshots and task-side zero request reporting so a cache hit releases its reservation without inventing an operation.
- Proved a second identical Workflow uses cached preparation generation, creates no preparation-generation operations, and still executes query contacts.
- Added delayed deterministic providers for crash-window acceptance.

### Failures and fixes
- The first retry integration remained terminal because `task.cjs` converted every native error to non-retryable `RAG_WORKFLOW_TASK_FAILED`. The native result envelope fixed this without exposing raw errors.
- The first cache-hit integration failed `BUDGET_USAGE_INVALID`: a reserved provider task with no operation reported no request dimension. The task now reports authoritative decorator usage, including zero on success.
- Reporting failed-contact usage as actual conflicted with conservative failed-operation settlement. Failed tasks now omit explicit usage and let Scraper conservatively settle the reservation; successful tasks report exact operation-derived usage.
- I explored reserving token and cost dimensions by default. This exposed a multi-dimension settlement mismatch and was not retained. Production currently reserves bounded requests; token and cost remain exact operation counters, while model response-token and pricing policy remain host-authority inputs. This limitation is documented rather than hidden.

### Review
Run `TestProviderWorkflowRetriesFailedContactAsDistinctOperation` and `TestProviderWorkflowCacheHitsDoNotCreateProviderOperations`. Confirm secret canaries cannot cross the safe failure envelope.

## Step 12: Fix long-running provider lease custody in Scraper

**Commit (Scraper):** `981d0a0` — "fix: renew active Workflow V3 leases"

### Work
- Added fenced `Store.RenewLease`: renewal requires the exact run, node, lease token, cancel epoch, running status, and an unexpired current lease.
- Extended `watchLease` to renew at half-life while retaining cancellation checks.
- Added renewal/cancellation store tests and a runtime heartbeat test that holds authority across more than three original lease periods.
- Pushed the Scraper commit and updated RAG-eval to pseudo-version `v0.0.5-0.20260723213944-981d0a0d84af`.

### Failure and fix
- The first authorized provider run used a five-second lease. Real generation exceeded it, producing 18 lease-lost attempts and repeated provider contacts. This was a production lifecycle defect, not a test timeout. A first rerun accidentally used the cached Scraper module because the RAG `go.work` does not include Scraper; a temporary workspace including local Scraper proved the renewal fix. After publishing the commit and updating `go.mod`, `GOWORK=off` acceptance passed with a two-second lease and one attempt per node.
- `go get` initially failed because `sum.golang.org` returned HTTP 500 for the new pseudo-version. Retried through `GOPROXY=direct` with narrowly scoped `GONOSUMDB=github.com/go-go-golems/scraper`.
- A parallel full Scraper suite twice timed out in the known `TestLazyMapOutputDigestIsIndependentOfConcurrency`; the focused test and two sequential package reruns passed (18.39s, then 66.41s and 57.63s). The rest of the full suite passed.

### Review
Audit renewal fencing in `workflowv3sqlite/store.go` and fail-closed behavior in `workflowv3runtime/engine.go`. Verify cancellation prevents renewal.

## Step 13: Authorized real-provider acceptance

**Commit:** `f6e083d` — authorized real-provider Workflow acceptance.

### Configuration and authorization
- Explicit authorization: `RAG_WORKFLOW_REAL_PROVIDER_ACCEPT=1`.
- Generation: Geppetto profile `umans-flash` from the configured Pinocchio profile registry.
- Embedding: local Ollama `nomic-embed-text:latest` at `127.0.0.1:11434`.
- Reranking: authorized SSH tunnel `127.0.0.1:18012 -> mimimi-2.local:127.0.0.1:8012`, llama.cpp reranker health verified.
- Provider cache copied to a fresh temporary host-config root for each acceptance run.

### Work and evidence
- Added an opt-in real-provider test with one bounded source and query loaded from the canonical TTC WordPress SQLite database.
- Recursive chunking produced six combined-preparation batches. The final run executed ten real contacts: seven generation, two embedding, and one rerank.
- Asserted exact authority, one attempt per node under a two-second renewable lease, succeeded operation outcomes, generation input/output token and cost counters, rerank input tokens, embedding output cardinality, and grounded answer citations.
- Final TTC acceptance passed in 41.93 seconds; log is `sources/smoke/02-real-provider-acceptance.txt`.

### Failure and fix
- The first post-heartbeat real run failed after a succeeded generation operation. A direct diagnostic exposed `RAG_COMBINED_RESPONSE_QUESTION_COUNT ... got 4 want 2`; the canonical `ttc-combined-preparation-v2` prompt requires four questions. Updated the real smoke fixture to `questionsPerChunk: 4`. Direct engine execution and durable Workflow execution then passed.

### Review
Run the opt-in script only with reviewed endpoints and cost authority. Confirm `FixtureProviders=false`, a fresh cache, the canonical TTC database binding, ten durable operations, and non-empty citations.

## Step 14: Cross-repository crash, resume, and privacy acceptance

### Work
- Extended the built-binary Researchctl smoke to two cases × two replicates.
- The first runner process now polls SQLite and is killed only after an incomplete provider operation exists. The abandoned subordinate Workflow preserves exactly one open operation; the Researchctl process retry creates a new Workflow and succeeds.
- Proved immutable resume (`executed=0`, `resumed=4`), exact attempt counts `[1,1,1,2]`, five subordinate Workflow databases, ten completed operations per selected successful run, direct RAG parity, canonical observation reprojection, timeout-to-canceled propagation, and one crash-window open operation.
- Added boundary/privacy guards for dependency direction, local replacements, unfinished markers, secret-like fixture material, raw provider payload fields, generic dependency graphs, and canonical regeneration.
- Ran prescribed Geppetto event/observability/embedding/rerank suites without modifying Geppetto.

### Evidence
- `sources/smoke/01-summary.json`
- `sources/smoke/02-real-provider-acceptance.txt`
- `sources/smoke/03-boundary-privacy-guards.txt`

### Review
Inspect the crash wrapper in `scripts/01-smoke-provider-workflow.sh`; it must kill only after admission and before completion. Confirm the abandoned database has one operation without a completion row and selected successful runs each export ten completed operations.

## Step 15: Final validation, documentation, and closure

### Final commits
- RAG `d4a957c` — replay, cache, crash custody, and delayed fixtures.
- RAG `0bd34d8` — generated logging metadata and race-safe acceptance deadline.
- RAG `b9415b2` — help, report, scripts, evidence, diary, and task completion.
- Scraper `981d0a0` — fenced renewable leases (published and pinned).

### Validation
- `GOWORK=off go test ./... -count=1` passed in RAG-eval.
- Focused RAG provider race suites passed after increasing only the test polling deadline from 20s to 60s; production timeouts were unchanged.
- `make lint`, Glazed lint, `make web-build`, `make logcopter-check`, module tidiness, all provider command builds, runner/fixture help, and embedded help lookup passed.
- Full Researchctl tests and lint passed.
- Full Scraper tests passed except one concurrent occurrence of the previously documented map timing test; focused and two sequential package reruns passed. Full Scraper lint passed.
- Geppetto event, observability, embedding, rerank, factory, and llama.cpp suites passed; Geppetto remained clean.
- Deterministic matrix smoke, authorized real-provider acceptance, fixture regeneration, and boundary/privacy guards passed from permanent ticket scripts.
- `docmgr validate frontmatter` passed. Initial `docmgr doctor` warned that new report topics `experiments` and `workflow-v3` were outside repository vocabulary; replaced them with existing `research` and `workflow`. Doctor then passed.
- All four ticket tasks are checked. Ticket closed with `docmgr ticket close`; final doctor passed.

### Final failure provenance
- Running Scraper and Researchctl `golangci-lint` concurrently produced `parallel golangci-lint is running` for Researchctl. Reran sequentially; both passed.
- The initial RAG provider race suite exceeded a 20-second test polling deadline in the synthetic variant. Increased the test-only deadline to 60 seconds; the race suite passed in 11.38 seconds on rerun.
- The first build/help smoke grepped for GNU-style `--provider-config`, while Go `flag` renders `-provider-config`; corrected the smoke assertion. The command itself was correct.
- `make logcopter-check` found missing generated package metadata for the two new packages. Ran `make logcopter-generate`, reviewed the two generated files, and reran the check successfully.

### Requirement audit
The acceptance report maps every goal clause to fresh source, ledger, test, command, or cross-repository evidence. There are no deferred implementation requirements. Request counts are enforced Workflow budgets; token/cost values are exact bounded operation evidence and are not misrepresented as an aggregate currency ceiling.
