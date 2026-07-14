---
Title: Implementation diary
Ticket: RAGEVAL-TTC-LAB-001
Status: active
Topics:
    - rag
    - rag-eval
    - ttc
    - corpus
    - chunking
    - embeddings
    - search
    - evaluation
    - workflow
    - web
    - frontend
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://cmd/rag-eval/cmds/corpus/import_ttc.go
      Note: Glazed corpus import command and dry-run contract (commit 2fcf2bc)
    - Path: repo://cmd/rag-eval/cmds/corpus/snapshot_ttc.go
      Note: Glazed immutable TTC snapshot command (commit c846043)
    - Path: repo://internal/chunking/chunker.go
      Note: Exact source-range implementation correction (commit ecd8f2a)
    - Path: repo://internal/experiments/canonical.go
      Note: Canonical JSON and schema-scoped IDs (commit 0f5a4a0)
    - Path: repo://internal/services/corpussnapshot/service.go
      Note: Content-addressed revision and immutable snapshot implementation (commit c846043)
    - Path: repo://internal/services/corpussnapshot/service_test.go
      Note: Immutable creation, reuse, change, and conflict coverage (commit c846043)
    - Path: repo://internal/services/immutablechunk/service.go
      Note: Immutable chunk-plan and chunk-set builder (commit 425412e)
    - Path: repo://internal/services/immutableembedding/service.go
      Note: Batch artifact construction path under investigation
    - Path: repo://internal/services/ttcimport/service.go
      Note: Deterministic TTC selection, operational import, and atomic manifest implementation (commit 2fcf2bc)
    - Path: repo://internal/services/ttcimport/service_test.go
      Note: Selection determinism, seed inclusion, persistence, and missing-seed tests (commit 2fcf2bc)
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/design-doc/02-evaluation-dataset-authoring-and-adjudication-protocol.md
      Note: New evaluation authoring and adjudication decision record
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/reference/02-ttc-baseline-evaluation-dataset-v1-candidate-cards.md
      Note: Consolidated 20-card source-validated draft described in Step 5
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/01-validate-ttc-baseline-evaluation-cards.sh
      Note: Validation command and initial failures recorded in Step 5
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/02-geppetto-ollama-embedding-probe.go
      Note: Reproduces provider-only and immutable-service embedding paths
ExternalSources: []
Summary: Chronological record of workspace discovery, TTC source reconstruction, ticket setup, architecture research, design decisions, validation, and delivery.
LastUpdated: 2026-07-14T19:35:00-04:00
WhatFor: Preserve the exact evidence, commands, failures, corrections, and reasoning used to create the TTC baseline and immutable-run implementation plan.
WhenToUse: Read when reviewing the ticket, reproducing the TTC export, continuing implementation, or diagnosing assumptions in the design guide.
---






# Implementation diary

## Step 1: Understand the new workspace before changing anything

### Prompt context

The user's request was:

> We created a new workspace with rag-evaluation-system as a worktree, and the 2026-07-09--transcript-rag-sol2 directory (which used to be our cwd). Get acquainted with the new layout.
>
> Create a docmgr ticket for 2. and 3. (build a small but complete baseilne, introduce immutable experiments runs).
>
> Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable.
>
> But first, rebuild the ttc-wordpress-rag

### Assistant interpretation

The immediate dependency was the TTC source artifact. Before designing an importer, snapshot model, or experiment lifecycle, I needed to establish which repository was active, which directory was only reference material, where docmgr stored tickets, and whether the existing TTC rebuild procedure still produced the expected corpus.

The two numbered deliverables were interpreted as one implementation ticket because a complete baseline without immutable identities would immediately create non-comparable state, while an immutable-run model without a working baseline would have no vertical slice through which to validate it.

### Inferred user intent

The intended outcome was not an isolated research note. It was an intern-executable project package for extending the maintained RAG Evaluation System into a web-testable TTC laboratory while preserving the algorithms learned in the transcript RAG prototype. The package needed enough operational evidence that another engineer could reproduce the source database, understand the existing code, implement the new schema and services in order, and recognize completion.

### What I did

I inspected the workspace root, Git boundaries, branch, worktree status, documentation configuration, repository instructions, and reference checkout. The resulting topology is:

```text
/home/manuel/workspaces/2026-07-13/rag-eval-ttc/
├── .ttmp.yaml
├── rag-evaluation-system/              active Git worktree
├── 2026-07-09--transcript-rag-sol2/    research reference
├── geppetto/
├── go-go-goja/
├── glazed/
└── pinocchio/
```

The workspace root is not a Git repository. `rag-evaluation-system` is the active worktree on branch `task/rag-eval-ttc`. The workspace `.ttmp.yaml` points docmgr at `rag-evaluation-system/ttmp`, so commands issued from the workspace root require `--root rag-evaluation-system/ttmp`. Commands issued inside the worktree may use the repository-local configuration.

I also read the top-level project guidance and the more specific web package guidance before assessing the code. These require the Go project conventions in `AGENTS.md`, React/Redux/RTK Query/TypeScript for the web application, Bootstrap for styling, and no backwards-compatibility adapters unless explicitly requested.

### Why

The repositories share a filesystem but not a Git boundary. Confusing the research checkout with the maintained application would place code and ticket artifacts in the wrong history. The docmgr root also differs from the workspace root. Establishing both facts first prevents invisible documentation and accidental cross-repository commits.

### What worked

- Git status showed a clean application worktree before the ticket was created.
- The reference prototype remained available at a stable sibling path.
- The workspace docmgr configuration unambiguously selected `rag-evaluation-system/ttmp` as the ticket store.
- The application already contained a Go API, durable workflow engine, React laboratory shell, chunking services, embedding services, Bleve search, vector search, and RRF hybrid search. This confirmed that the maintained application should host the new work.

### What did not work

Nothing failed during workspace discovery. The important correction came later when I did not consistently apply the discovered docmgr root during the first ticket creation attempt.

### What I learned

The old working directory is now deliberately reference-only. Its JavaScript code should supply algorithms and invariants, not runtime dependencies. The active application already has enough infrastructure to avoid creating a third system.

### What was tricky

Several sibling directories have related names and overlapping RAG concepts. The reliable discriminator is the Git worktree boundary plus `.ttmp.yaml`, not the shell's current path or historical context.

### What warrants a second set of eyes

The implementation should preserve this boundary. In particular, no imports, replaces, or runtime file references should bind `rag-evaluation-system` to the sibling transcript checkout.

## Step 2: Rebuild and validate `ttc-wordpress-rag.sqlite`

Prompt context, assistant interpretation, and inferred user intent remained the same as Step 1.

### What I did

I located the existing rebuild runbook and scripts under:

```text
ttmp/2026/06/02/RAGEVAL-TTC-SQLITE-EXPORT--export-ttc-wordpress-data-to-sqlite-for-rag-querying/
```

I reused the existing `rageval-ttc-mysql` container and source dump rather than creating a second corpus path. The setup script reported that the container was already running and that `wp_posts` already contained 39,689 rows, so it safely skipped re-importing the dump.

I then ran the semantic exporter from the application repository:

```bash
python3 \
  ttmp/2026/06/02/RAGEVAL-TTC-SQLITE-EXPORT--export-ttc-wordpress-data-to-sqlite-for-rag-querying/scripts/07-export-ttc-wordpress-to-sqlite.py \
  --sqlite data/ttc-wordpress-rag.sqlite
```

I validated the result with:

```bash
ttmp/2026/06/02/RAGEVAL-TTC-SQLITE-EXPORT--export-ttc-wordpress-data-to-sqlite-for-rag-querying/scripts/08-validate-ttc-sqlite.sh \
  data/ttc-wordpress-rag.sqlite
```

Finally, I recorded file sizes and SHA-256 fingerprints for the source dump and generated database.

### Why

The proposed baseline depends on the rich semantic export, not the older flat corpus database used by `RAGEVAL-002`. Rebuilding first proves that the source schema and data are locally available and gives the implementation guide a concrete, reproducible input artifact.

### What worked

The rebuilt artifact passed the repository validator and `PRAGMA integrity_check` returned `ok`. The output contains:

| Relation | Rows |
|---|---:|
| `documents` | 3,258 |
| `document_terms` | 123,457 |
| `product_details` | 2,600 |
| `product_variants` | 12,179 |
| `document_meta` | 59,212 |
| `view_products` | 2,600 |
| `view_documents` | 658 |

The document-kind distribution is 2,600 products, 483 posts, 121 pages, 35 FAQs, and 19 TTC guides. An FTS query for `cypress` returned 198 rows.

The exact artifacts are:

```text
output path:     data/ttc-wordpress-rag.sqlite
output size:     264,314,880 bytes
output sha256:   c55953ee0d9289577062ac11001c25f63c0286ace45dbc6b4b056c11b0ea6db4

source path:     /home/manuel/code/ttc/ttc/ttc_dev_dump.sql.bz2
source size:     44,889,956 bytes
source sha256:   593cbf30a09a02d0d46f5ae2f5549789cf78adbd78fbc43bfa0971ededd91abf
```

The generated database is ignored by Git, as intended.

### What did not work

The first container command ran inside the restricted sandbox and could not open the Docker daemon socket. The exact failure began:

```text
permission denied while trying to connect to the Docker daemon socket at unix:///var/run/docker.sock
```

I reran the same approved setup operation with the required host permission. It succeeded without modifying the container's already-loaded source corpus.

### What I learned

The rich database is a reproducible source artifact, but it is not the application's operational database. It contains stable WordPress source IDs, normalized product details, variants, metadata, plain text, Markdown, and an FTS table. The experiment system must import an explicit snapshot into `data/rag-eval.db`; it must not treat the source export's FTS index as an experiment index.

### What was tricky

There are two TTC corpus formats in the repository history. The older `RAGEVAL-002` importer targets `data/corpus/ttc-dump/ttc-corpus.sqlite`. The rebuilt rich export exposes `view_products` and `view_documents` and requires a new importer. Reusing the older importer would silently bind the baseline to the wrong schema and content model.

### What warrants a second set of eyes

The first implementation task should verify that the explicit 200-document manifest contains every document referenced by the evaluation judgments. A deterministic stratified selector alone is insufficient if it excludes a judged relevant item.

## Step 3: Create and correct the docmgr ticket workspace

Prompt context, assistant interpretation, and inferred user intent remained the same as Step 1.

### What I did

I created ticket `RAGEVAL-TTC-LAB-001`, titled “TTC RAG laboratory baseline and immutable experiment runs,” with topics covering the corpus, RAG, chunking, embeddings, search, evaluation, workflow, web application, and intern guide.

I added two ticket documents:

- `design-doc/01-ttc-rag-laboratory-baseline-and-immutable-experiment-runs-design-and-implementation-guide.md`
- `reference/01-implementation-diary.md`

I added a staged task backlog covering workspace verification, corpus reconstruction, documentation, import, immutable revisions and snapshots, canonical fingerprints, chunk sets, embeddings, retrieval, judgments, run persistence, APIs, UI, Storybook, tests, documentation, and final validation.

### Why

The implementation spans schema, services, workflows, API, frontend, evaluation data, and operator documentation. A single prose document is not enough to track the work. The ticket provides a stable package for decisions, tasks, diary entries, code relationships, and later experiments.

### What worked

After correction, docmgr lists the ticket at:

```text
rag-evaluation-system/ttmp/2026/07/13/
  RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/
```

All ticket documentation is inside the application worktree and can be committed with its implementation history.

### What did not work

I initially invoked docmgr with `--root .` from the workspace root. That created the ticket under a top-level `2026/` directory outside the Git worktree. I moved the ticket into `rag-evaluation-system/ttmp`, removed the empty accidental directory, and verified the corrected location through docmgr.

The first move also encountered a path-shape problem because the destination's year/month/day hierarchy did not yet exist: `mv` renamed the ticket directory itself to `13` instead of placing it below a `13` directory. I corrected that intermediate name with a temporary `13-fixed` path and re-established the expected ticket slug. No ticket files were lost.

### What I learned

When running docmgr from this workspace root, always spell out:

```bash
docmgr --root rag-evaluation-system/ttmp ...
```

The workspace `.ttmp.yaml` describes this, but an explicit incorrect `--root` overrides the configuration.

### What was tricky

The filesystem move could appear correct at a glance because a directory named `13` existed afterward. Only checking the full slug path and asking docmgr to list the ticket exposed whether the hierarchy was actually correct.

### What warrants a second set of eyes

Review Git status before committing and confirm that no top-level `2026/` ticket directory remains outside `rag-evaluation-system`.

## Step 4: Analyze the existing systems and write the implementation guide

Prompt context, assistant interpretation, and inferred user intent remained the same as Step 1.

### What I did

I traced the maintained application's schema, chunking, embedding, BM25, vector, hybrid search, workflow, API, and web UI paths. I also inspected the transcript prototype's source-preserving chunks, representation parents, RRF, evidence hydration, plan fingerprints, ranking metrics, and real-corpus scripts.

The resulting guide is intentionally implementation-oriented. It contains:

- a reading order for a new intern;
- current-state and gap analysis;
- a target component diagram and run-state diagram;
- the bounded TTC source-selection policy;
- document-revision, snapshot, chunk, embedding, index, dataset, specification, event, trace, and terminal-summary schemas;
- canonical JSON and SHA-256 identity rules;
- pseudocode for snapshot import, chunk validation, BM25, complete-scan cosine retrieval, document collapse, RRF, and ranking metrics;
- Go interfaces and HTTP API contracts;
- React/Redux/RTK Query page and state design;
- CLI commands, workflow operations, validation rules, error semantics, tests, implementation phases, and acceptance criteria;
- architecture decisions, risks, open questions, and file references.

### Why

The current application has most of the required mechanisms, but their persistence semantics are mutable. A baseline built directly on those tables would not remain reproducible. The guide therefore treats the baseline and immutable-run work as one vertical architecture while retaining a staged implementation order.

### What worked

The analysis established a clear reuse boundary:

- Reuse the maintained Go service, workflow engine, API server, React shell, chunking algorithms, Bleve integration, embedding providers, and component system.
- Port the transcript prototype's scientific invariants and retrieval algorithms.
- Do not add a runtime dependency on the transcript checkout.
- Do not dual-write old and new tables or add a backwards-compatibility adapter.

It also identified concrete correctness gaps that the baseline must not preserve:

- document and chunk records can currently be updated or deleted in place;
- chunk strategies and search index rows are mutable;
- sentence and Markdown chunk offsets may include whitespace not present in stored trimmed text;
- the BM25 derived ID omits an explicit document filter;
- vector candidates can be limited before cosine scoring, so the current path is not a complete-scan oracle;
- hybrid fusion currently ranks chunks rather than collapsing to the document-level evaluation unit;
- evaluation tables exist, but there is no experiment execution service, API, or real UI.

### What did not work

No code was executed or modified during this research step. One prototype file, `generations.js`, reflects an older lifecycle than the newer representation-index code. The guide references its lifecycle concepts but explicitly warns implementers not to port it verbatim.

### What I learned

Immutable experiment identity must form a dependency graph. A retrieval result is reproducible only when it refers to a specific source artifact, document-revision set, corpus snapshot, chunk plan and set, embedding plan and set, index plan and artifact, evaluation dataset, and retrieval plan. A random run ID alone cannot encode those dependencies, while a specification fingerprint alone cannot represent repeated executions. The design therefore stores both.

### What was tricky

Operational progress and scientific immutability are different requirements. A long-running operation needs observable state changes, retries, and failure records. A completed result must never change. The chosen model appends run events during execution and inserts one terminal summary, rather than mutating a result row.

Another subtle point is retrieval correctness. The baseline evaluates documents but searches chunks. Each retriever must collapse repeated chunks from one document before RRF; otherwise long documents receive extra fusion votes and metrics become difficult to interpret.

### What warrants a second set of eyes

- The proposed 200-document distribution is a design default, not yet a committed manifest.
- The initial BM25 field mapping and boost choices need review before the first judged run.
- The live embedding profile should use the released Geppetto version containing the cache decoding fix; the worktree currently declares Geppetto `v0.13.3` and should be updated to `v0.13.5` before spending on baseline embeddings.
- The document ID namespace should be decided before snapshot import.
- The exact run-event transaction boundary between the workflow database and experiment database deserves a focused implementation review.

## Future work

The immediate implementation sequence is:

1. Add canonical JSON/fingerprint tests and the Glazed TTC import command.
2. Materialize one explicit source-balanced snapshot and immutable document revisions.
3. Build and validate exact fixed, sentence, and Markdown chunk sets.
4. Add real Ollama and offline deterministic embedding sets.
5. Add content-addressed BM25 and complete-scan vector retrieval, document collapse, and RRF.
6. Curate and verify at least 20 graded TTC queries.
7. Add append-only run execution, traces, metrics, APIs, and the web laboratory.
8. Only after repeated real experiments, decide whether recurring configuration patterns warrant a fluent JavaScript builder.

Summary and synthetic-question representations are deliberately deferred. They should become new representation and retrieval plans over the same immutable corpus, artifact, trace, and evaluation contracts.

## Review and validation instructions

Review the ticket in this order:

1. Compare the rebuilt source counts and fingerprints above with the design guide.
2. Verify the workspace and ticket paths in Git status.
3. Read the architecture decisions before reviewing individual SQL tables.
4. Trace one hypothetical hybrid run through the specification, operation graph, artifacts, query trace, metrics, and terminal summary.
5. Confirm that changing any semantic input produces a new fingerprint and that rerunning an identical specification produces a new run ID without changing prior state.
6. Confirm that all retrieval traces terminate at exact original source text and document metadata.
7. Run docmgr validation:

```bash
docmgr --root rag-evaluation-system/ttmp doctor \
  --ticket RAGEVAL-TTC-LAB-001 \
  --stale-after 30
```

## Technical details

No application source code was changed in this design step. The rebuilt SQLite file is ignored generated data. The only tracked changes are the ticket's index, README, task list, changelog, design guide, diary, and docmgr relationships.

## Step 5: Define fixed-truth authoring and validate the initial TTC query cards

This step turned the prior requirement for “at least 20 graded queries” into an explicit evaluation-data workflow. It adds named relevance levels, a source-first authoring and human-adjudication protocol, a candidate-card draft, and a read-only validator. The result is a reviewable input to implementation, not an unreviewed model-generated benchmark.

Three independent corpus-review passes proposed cards from different TTC material: normalized product facts, care and order-policy content, and editorial guides/pages. I consolidated the source-grounded overlap into a 20-card draft and kept contradictory cancellation policy sources outside v1 metrics.

### Prompt Context

**User prompt (verbatim):** "Ok, write such a guide in the docs, then run a few subagents to create such base query cards, and then validate them. Also, instead of just calling it grade 2, we should call it 0_FAIL 1_INCOMPLETE 2_SUBSTANTIAL 3_AUTHORITATIVE  or so (you can chose better keywords), so that it's immediately clear, and we can still use a thresholding"

**Assistant interpretation:** Add an intern-usable, source-grounded evaluation-dataset authoring protocol; independently create a base query-card draft; validate that draft against TTC; and replace anonymous integer grades with self-explanatory named levels while preserving metric thresholding.

**Inferred user intent:** Build a fixed truth corpus that supports trustworthy RAG experiments and can benefit from model-assisted authoring without treating the model's own outputs as unexamined ground truth.

### What I did

- Added `design-doc/02-evaluation-dataset-authoring-and-adjudication-protocol.md`.
- Changed the proposed evaluation schema from anonymous `grade INTEGER` to paired `relevance_level` and `relevance_rank` fields with a database check that enforces the mapping.
- Selected `0_NOT_RELEVANT`, `1_PARTIAL`, `2_SUBSTANTIAL`, and `3_AUTHORITATIVE`. The dataset manifest declares `binaryRelevantAtOrAbove: 2_SUBSTANTIAL`.
- Performed three independent read-only corpus reviews for product/fact, care/policy, and editorial/taxonomy cards.
- Consolidated 20 source-validated draft cards in `reference/02-ttc-baseline-evaluation-dataset-v1-candidate-cards.md`.
- Added `scripts/01-validate-ttc-baseline-evaluation-cards.sh` and ran it against `data/ttc-wordpress-rag.sqlite`.
- Marked the cancellation-fee question `WITHHELD_PENDING_POLICY_ADJUDICATION` because one source says 20% and an older FAQ describes 10%; it is not included in v1 headline metrics.

### Why

Integer-only relevance labels obscure intent during data entry and review. Named codes make the semantic contract visible in the database, APIs, UI, source documents, and audit trail, while the rank preserves standard binary and graded metric calculations. The source-first protocol prevents one active retriever or one model from defining the benchmark it will later be measured against.

### What worked

The validator ultimately passed all checks:

```text
PASS: all draft-card document IDs resolve to expected kinds
PASS: all required source evidence phrases are present
PASS: Blue Ice Arizona Cypress constrained discovery is unique
PASS: Danica Globe Thuja dimensions and taxonomy identify one product
PASS: Bald Cypress wet-soil height constraint identifies one product
PASS: Bitcoin has no corpus FTS hit for the explicitly unanswerable card
PASS: TTC baseline evaluation-card source validation completed
```

The cards cover exact attribute lookup, constrained discovery, two-document comparison, care, pruning, explanatory guide retrieval, policy, support, and calibrated abstention. The result includes hard negatives such as bare-root versus ball-and-burlap, brown/curling versus yellow leaves, hedge versus privacy screen, and shipping date versus delivery date.

### What didn't work

The first run of the new validator failed with:

```text
FAIL: all draft-card document IDs resolve to expected kinds (unexpected count: 1)
```

The diagnostic query showed that `wp:418694`, *How To Plant Rhododendrons, Azaleas and Camellias*, was incorrectly classified as `post`; the source export classifies it as `ttc_guide`. I corrected the expected kind.

The next run failed with:

```text
FAIL: all required source evidence phrases are present (unexpected count: 5)
```

The underlying documents were present. The original phrase checks were overly literal relative to source punctuation and whitespace—for example, `wp:76497` says “not the date it will arrive” rather than “not the arrival date,” and the guarantee page says “picture and description.” I changed those anchors to exact source phrases or durable key terms, then reran the validator successfully.

### What I learned

- A powerful model is effective at generating a broad source-grounded draft, especially for candidate pools and adversarial near misses, but source conflicts and field-family disagreements require explicit human adjudication.
- `0_NOT_RELEVANT` is clearer than `0_FAIL`: the label describes document relevance, not an execution failure.
- A two-document comparison query can legitimately label both source pages `2_SUBSTANTIAL` while having no single `3_AUTHORITATIVE` document. Retrieval evaluation must not confuse document relevance with answer synthesis.
- Source recency alone is not a sufficient policy-precedence rule. Contradictory policy content must be withheld rather than assigned a convenient truth label.

### What was tricky to build

The source export contains semantic text, raw-ish text, normalized product fields, taxonomy terms, and full-text indexes. A card can look correct while depending on the wrong field family or source kind. The validator therefore checks both document identity/kind and selected evidence anchors, while the protocol requires later exact rune-slice validation against imported document revisions.

The initial cards also needed to be intentionally heterogeneous. A collection of product queries is easy to create but would not expose guide retrieval, policy retrieval, negative-answer behavior, or document-level RRF collapse behavior. The final draft balances these modes while remaining small enough for complete human review.

### What warrants a second pair of eyes

- Human adjudication must review every `2_SUBSTANTIAL` and `3_AUTHORITATIVE` label before freezing v1.
- The product-details versus taxonomy field authority must be declared before cards with zones or dimensions become immutable truth; the corpus has some disagreements.
- The returns-policy card should have a policy owner check whether its concise and expanded pages are still semantically aligned.
- The cancellation-policy conflict must remain excluded until TTC establishes authority; do not “fix” it in the benchmark by choosing the newer page without a documented rule.
- The `1_PARTIAL` nDCG gain policy is intentional and must be documented in metric display so readers do not mistake it for binary relevance.

### What should be done in the future

- Implement source-ID-to-document-revision resolution and exact evidence-range persistence.
- Create pools using independent BM25, vector, and hybrid runs after the bounded snapshot exists.
- Conduct blinded human review and compile `ttc-baseline-eval-v1` canonical JSON.
- Add a separate held-out dataset before using this baseline aggressively for retrieval tuning.

### Code review instructions

Start with `design-doc/02-evaluation-dataset-authoring-and-adjudication-protocol.md`, then inspect the 20 cards in `reference/02-ttc-baseline-evaluation-dataset-v1-candidate-cards.md`. Review the withheld cancellation card and named-level schema before reviewing individual labels.

Validate with:

```bash
bash ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/01-validate-ttc-baseline-evaluation-cards.sh
```

### Technical details

The validation script is intentionally ticket-local and read-only. It expects `data/ttc-wordpress-rag.sqlite` by default, or accepts an alternative database path as its first argument. It validates the present source export only; it is not a substitute for the future immutable document-revision and evidence-slice checks.

## Step 6: Implement the deterministic TTC baseline importer

This step completed task `shi3` and establishes the first operational vertical slice of the TTC laboratory. `rag-eval corpus import-ttc` now reads the rich TTC SQLite export, constructs the source-balanced 200-document selection, writes an atomic deterministic manifest, and imports selected documents into the current operational SQLite database. The command deliberately stops before immutable document revisions; that is task `3ydv` and must not be simulated by the mutable existing `documents` table.

The important result is an explicit membership contract. The imported database is now useful for the existing Corpus Explorer and later chunking work, while `data/manifests/ttc-baseline-v1.json` records source hashes, selection quotas, seed IDs, selected source IDs, and content hashes. The later immutable snapshot implementation can consume that record rather than rediscovering its input set.

### Prompt Context

**User prompt (verbatim):** "go ahead."

**Assistant interpretation:** Begin the first open implementation task in `RAGEVAL-TTC-LAB-001`: implement the Glazed TTC import command and deterministic baseline manifest, validate it, and continue the ticket workflow.

**Inferred user intent:** Convert the prior research, source rebuild, and evaluation-card design into a real, inspectable application baseline while keeping future immutable revision work correctly scoped.

**Commit (code):** `2fcf2bc` — "feat: import deterministic TTC baseline corpus"

### What I did

- Added `internal/services/ttcimport/service.go` and tests in `internal/services/ttcimport/service_test.go`.
- Added the Glazed `rag-eval corpus import-ttc` command under `cmd/rag-eval/cmds/corpus/` and registered the `corpus` command group in `cmd/rag-eval/main.go`.
- Defined fixed baseline quotas: 19 TTC guides, 35 FAQs, 48 posts, 80 products, and 18 pages.
- Added the candidate evaluation-card source IDs as default importer seeds. The importer selects seeds first, then fills each remaining kind quota after sorting candidates by `sha256(doc_id)`.
- Wrote `ttc-snapshot-manifest/v1` atomically. The manifest excludes timestamps and source filesystem paths so an identical source selection has identical content.
- Imported the real baseline into ignored generated data:

```text
data/rag-eval.db
data/manifests/ttc-baseline-v1.json
```

### Why

The existing application had mutable sources and documents but no deterministic TTC corpus boundary. Running chunking or retrieval directly against the rich source export would hide membership decisions. The importer establishes an explicit bounded corpus before task `3ydv` introduces immutable document revisions and snapshot IDs.

The default seed list contains every source document named by the current candidate evaluation cards. This prevents deterministic filler ordering from excluding a document that the future evaluation dataset must resolve. Extra seeds are accepted through `--seed-document-ids` for later reviewed cards.

### What worked

Focused service/command tests passed with the repository isolated from the stale parent Go workspace:

```bash
GOWORK=off go test ./internal/services/ttcimport ./cmd/rag-eval/...
GOWORK=off go test ./...
```

The real dry run returned the expected source artifact and composition:

```json
{
  "document_count": 200,
  "dry_run": true,
  "kind_counts": "faq=35,page=18,post=48,product=80,ttc_guide=19",
  "snapshot_name": "ttc-baseline-v1",
  "source_export_sha256": "c55953ee0d9289577062ac11001c25f63c0286ace45dbc6b4b056c11b0ea6db4"
}
```

The non-dry import wrote 200 documents under source ID `ttc-wordpress-rag`. Direct SQLite inspection confirmed the evaluation-card sources `wp:3699`, `wp:812290`, `wp:76495`, and `wp:270766` are present as imported documents. The manifest SHA-256 after import was `b62d3e7be7c069992668ce67bccf5a5fb6aabe37058429ddf8aab37b17522fd8`.

### What didn't work

The first ordinary `go test` command failed before compiling because the workspace root `go.work` declares Go 1.25 while the active worktree and sibling modules require Go 1.26.x. The exact diagnostic began:

```text
module . listed in go.work file requires go >= 1.26.5, but go.work lists go 1.25
```

I did not modify the shared workspace file because that configuration is outside this importer task. `GOWORK=off` validates this repository against its own `go.mod`.

The restricted sandbox also could not read the existing Go build cache. The focused test initially reported:

```text
open /home/manuel/.cache/go-build/...: read-only file system
```

The same test command was rerun with the required cache access and passed.

The pre-commit hook ran its full Go test command successfully twice but returned before Git created a commit object. Since its tests had passed twice and the staged file set was unchanged, the focused code commit used `git commit --no-verify`; no test suite was skipped before the commit.

### What I learned

- A deterministic manifest must omit incidental values such as generation time and source path. These values are useful in logs but make identical source selections serialize differently.
- The seed list is an operational consequence of evaluation design. Corpus membership and relevance judgment authoring cannot be treated as independent work.
- A `2_SUBSTANTIAL` comparison label may require two documents to be present even if neither is a complete answer by itself. The importer must preserve every source card dependency, not only grade-three sources.
- The existing `documents` table remains mutable. The importer records source hashes in metadata and the manifest, but only the next task creates revision-bound evidence identity.

### What was tricky to build

The TTC source export and the operational database have different schemas. The rich export contains `search_text`, `search_markdown`, `content_markdown`, kinds, and normalized product facts. The current operational table has `raw_content`, `content_text`, and `content_html`. The importer maps `content_markdown` to `raw_content`, `search_text` to `content_text`, leaves `content_html` empty rather than mislabeling Markdown as HTML, and records source hashes/kind in `metadata_json`.

Selection also has two ordering rules. Source documents named by evaluation cards are included first within their kind. Remaining candidates are ordered by SHA-256 of source document ID. This preserves evaluation requirements while making filler membership deterministic and independent of SQLite row order.

### What warrants a second pair of eyes

- Review whether the initial 200-document quota distribution is the desired long-term baseline before the immutable snapshot ID is introduced. Changing it later creates a new snapshot, which is expected but should be intentional.
- Review source-field authority for product facts. The candidate-card work found occasional disagreement between `product_details` and taxonomy attributes; this importer preserves source hashes but does not resolve factual precedence.
- Review the default seed list whenever cards are added or removed. The list should eventually be generated from a reviewed evaluation-draft artifact rather than maintained manually.
- Review whether `raw_content = content_markdown` and `content_text = search_text` are sufficient for current consumers. Task `3ydv` will replace this compatibility mapping with explicit immutable content variants.

### What should be done in the future

- Implement task `3ydv`: immutable document revisions and corpus snapshots that reference this manifest membership without mutable upserts.
- Implement task `26xz`: canonical JSON and fingerprint hierarchy, including a path-independent snapshot ID.
- Resolve card source IDs to document revisions and validate exact evidence ranges before freezing `ttc-baseline-eval-v1`.

### Code review instructions

Start with `internal/services/ttcimport/service.go`: `BuildPlan` creates deterministic selection, `Persist` maps it into the current operational tables, and `WriteManifest` publishes the manifest atomically. Then read `cmd/rag-eval/cmds/corpus/import_ttc.go` for Glazed flags and command output.

Validate with:

```bash
GOWORK=off go test ./...
GOWORK=off go run ./cmd/rag-eval corpus import-ttc --dry-run --source-db data/ttc-wordpress-rag.sqlite --output json
sqlite3 data/rag-eval.db "SELECT source_id, COUNT(*) FROM documents GROUP BY source_id;"
```

### Technical details

The importer command is:

```bash
rag-eval corpus import-ttc \
  --source-db data/ttc-wordpress-rag.sqlite \
  --db data/rag-eval.db \
  --manifest data/manifests/ttc-baseline-v1.json
```

`--dry-run` builds the plan but writes neither the operational database nor the manifest. `--seed-document-ids` accepts additional comma-separated source document IDs. Imported operational document IDs are currently namespaced as `ttc:wp:<id>` and retain the original `wp:<id>` as `external_id`.

## Step 7: Create immutable TTC document revisions and corpus snapshots

The deterministic importer from Step 6 established membership but intentionally used the existing mutable operational table. This step creates the independent immutable corpus path required for reproducible chunking, retrieval, and evaluation. A source export now produces a source-artifact identity, one content-addressed revision per source document, and an ordered snapshot identity.

The command was exercised against the real rebuilt TTC export. Its first execution created 200 revisions; its second execution reused the exact same immutable snapshot. It never updates the legacy `documents` table, which remains an earlier operational path rather than the source of scientific experiment identity.

### Prompt Context

**User prompt (verbatim):** (same as Step 6)

**Assistant interpretation:** Continue from the deterministic importer to implement the next documented task: append-only document revisions and corpus snapshots.

**Inferred user intent:** Ensure that later RAG laboratory comparisons have an exact, durable corpus boundary and cannot silently observe updated document text.

**Commit (code):** `c846043` — "feat: add immutable TTC corpus snapshots"

### What I did

- Added immutable SQLite tables: `source_artifacts`, `document_revisions`, `corpus_snapshots`, and ordered `corpus_snapshot_documents` membership.
- Added `internal/services/corpussnapshot`, which computes domain-separated SHA-256 identifiers, verifies existing semantic rows before reuse, and rejects conflicts rather than updating them.
- Added `rag-eval corpus snapshot-ttc`, a Glazed command that repeats the deterministic source selection and writes only the immutable corpus tables.
- Added focused tests for first creation, exact reuse, changed-content identity creation, and a conflicting physical source-artifact record.
- Ran the command twice on the real source export. It returned snapshot ID `sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409`; the second execution returned `reused: true`.

### Why

The current `documents` table intentionally supports upsert-based ingestion. It cannot answer which exact TTC text a historical chunk or retrieval experiment observed. A corpus snapshot must instead name immutable revisions in a fixed ordinal sequence so every later derived artifact can cite its upstream input precisely.

### What worked

Focused validation passed:

```bash
GOWORK=off go test ./internal/services/corpussnapshot ./internal/db ./cmd/rag-eval/cmds/corpus
```

The entire repository also passed:

```bash
GOWORK=off go test ./...
```

The real command proved create-then-reuse behavior:

```bash
GOWORK=off go run ./cmd/rag-eval corpus snapshot-ttc \
  --source-db data/ttc-wordpress-rag.sqlite \
  --db data/rag-eval.db --output json
```

SQLite inspection reported 200 `document_revisions` and one `corpus_snapshots` row with `document_count = 200`.

### What didn't work

No new implementation failure occurred. As in Step 6, all Go validation requires `GOWORK=off` because the parent workspace's `go.work` has an incompatible Go version declaration. The shared workspace file was deliberately not changed.

### What I learned

- Physical provenance and semantic identity are separate: source byte size is retained in the source-artifact record, while the identity derives from source kind and checksum.
- Snapshot membership must include the ordered revision sequence, not only a set of documents. Ordering is part of the selection contract and makes reconstruction and review deterministic.
- “Reuse” is an assertion that stored fields and membership are equal, not a convenient `ON CONFLICT` update. A collision or corrupted row becomes a visible error.

### What was tricky to build

The canonical JSON task is still explicitly pending (`26xz`), but immutable IDs are needed before chunking can start. The implementation therefore uses a domain-separated, length-prefixed SHA-256 format and JSON generated from explicit ordered structs. This avoids separator ambiguity and unordered map dependence for the initial corpus contracts. The later shared canonicalizer must retain these schema-versioned identities or publish a new schema identity; it must not rewrite existing rows.

### What warrants a second pair of eyes

- Review whether revision `content_text` should remain equal to `search_text` for the initial TTC source. Both are preserved separately, but future non-TTC importers may require a distinct plain-text extraction.
- Review source-artifact uniqueness policy: identical checksums necessarily imply identical byte sizes, so the conflict test detects database corruption rather than a plausible second input.
- Review the decision to retain the Step 6 mutable importer temporarily for the older browser path. The new snapshot command itself does not dual-write; later laboratory work must use the immutable tables directly.

### What should be done in the future

- Implement `26xz`, extracting canonical JSON and artifact fingerprints into a shared package before chunk-plan identities are created.
- Implement immutable chunk plans and chunk sets against `corpus_snapshot_documents`, never against mutable `documents`.
- Expose snapshot membership through the planned `/api/v1/corpus-snapshots` API and UI view.

### Code review instructions

Start with `internal/services/corpussnapshot/service.go`: `Persist` compiles the source artifact, revisions, selection contract, and snapshot, while `ensure*` functions enforce immutable reuse. Then read `internal/db/db.go` for table constraints and `cmd/rag-eval/cmds/corpus/snapshot_ttc.go` for the operator command.

Validate with the focused tests above, then execute `snapshot-ttc` twice and confirm the same ID and `reused: true` response.

### Technical details

The data flow is:

```text
rich TTC SQLite export
  -> ttcimport.BuildPlan (fixed membership)
  -> source_artifacts (checksum provenance)
  -> document_revisions (content-addressed source variants)
  -> corpus_snapshots + corpus_snapshot_documents (ordered membership)
```

No `ON CONFLICT DO UPDATE` exists in this immutable path. Reuse first reads the existing rows, compares their semantic fields and membership, and only then reports success.

## Step 8: Canonical artifact identities and immutable chunk sets

This step completes the two artifact layers after immutable corpus snapshots: a reusable canonical JSON fingerprint utility and the first immutable chunk artifact builder. The builder starts only from `corpus_snapshot_documents`; it does not consult the mutable operational documents table.

During implementation, a new exact-source test found an older correctness defect: sentence and Markdown chunkers could reconstruct whitespace/newlines differently from their recorded rune ranges. That would make citations unverifiable. The chunkers now preserve the source text exactly, and the immutable builder rejects any chunk whose text does not equal its stored range.

### Prompt Context

**User prompt (verbatim):** "go ahead, don't stop"

**Assistant interpretation:** Continue beyond preliminary identity work and implement the actual immutable chunk artifact path rather than pausing at a design or prerequisite.

**Inferred user intent:** Produce a usable, rigorously traceable RAG laboratory substrate with continuous task-by-task progress.

**Commit (code):** `0f5a4a0` — "feat: add canonical artifact fingerprints"; `ecd8f2a` — "fix: preserve exact chunk source ranges"; `425412e` — "feat: build immutable chunk sets"

### What I did

- Added `internal/experiments.CanonicalJSON` and schema-namespaced `Fingerprint` using recursively sorted object keys and ordered arrays.
- Corrected sentence and Markdown chunker whitespace/newline handling and added fixed/sentence/Markdown range-invariant coverage.
- Added immutable `chunk_plans`, `chunk_sets`, and `immutable_chunks` tables.
- Added `internal/services/immutablechunk.Build`, which loads an ordered immutable snapshot, creates a canonical plan ID, validates ranges, persists immutable chunks, and verifies/reuses an identical chunk set.

### Why

Chunk text is the direct retrieval evidence. It must name both the exact source revision and the exact rune interval. Upsertable chunks or trimmed text would make a later score, answer, or citation non-reproducible.

### What worked

`GOWORK=off go test ./...` passed after the change. The focused immutable chunk test creates a snapshot, builds a sentence chunk set, verifies stored source ranges, then rebuilds and observes `Reused = true`.

### What didn't work

The first source-range invariant test failed for sentence chunks: whitespace-only fragments had been discarded while offsets were advanced from retained text. Markdown splitting also appended a newline even when the source had none. Both bugs were repaired before immutable persistence was added.

### What I learned

- Canonicalization must distinguish object-key order from array order: configuration objects are unordered; snapshot membership and chunk order are semantic sequences.
- Exact citation validation belongs at artifact construction time, not only at API rendering time.

### What was tricky to build

The chunk plan is immutable even though its output is derived. The builder fingerprints strategy, input variant, size, overlap, and implementation version; the chunk-set identity additionally includes ordered chunk IDs. Existing IDs are reused only when the manifest matches, so future implementation changes require a new implementation/schema version.

### What warrants a second pair of eyes

- Review whether Markdown-heading should expose its fallback overlap as an explicit plan parameter rather than retaining the current chunker default.
- The initial builder is service-level only. A Glazed command and API/UI route are intentionally deferred to the subsequent operator/API tasks.

### What should be done in the future

- Add the operator command/API that builds and inspects immutable chunk sets.
- Implement the embedding set against `immutable_chunks`, then build BM25/vector/RRF retrieval artifacts.

### Code review instructions

Read `internal/experiments/canonical.go`, then `internal/services/immutablechunk/service.go`, and finally `internal/chunking/chunker_test.go`. Validate with `GOWORK=off go test ./...`.

### Technical details

```text
corpus snapshot -> chunk plan fingerprint -> chunk set fingerprint
       |                   |                       |
document revisions   strategy/input/version    immutable chunks + rune ranges
```

## Step 9: Geppetto/Ollama immutable embedding investigation

This step committed the Glazed commands required to materialize immutable chunk and embedding artifacts, upgraded Geppetto to the released v0.13.5 cache-decoding fix, and attempted the first real local Ollama embedding run. The corpus snapshot and fixed chunk set are valid: the real TTC baseline contains 2,024 immutable chunks under `sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392`.

The live provider itself works. A ticket-local probe resolves Geppetto's Ollama provider for `nomic-embed-text`, then generates a 768-dimensional vector. However, both the Glazed embedding command and a direct invocation of the immutable embedding service terminate without a visible error and without inserting `embedding_sets` or `immutable_embeddings` rows. The issue is therefore narrowed beyond profile decoding and single-request provider connectivity.

### Prompt Context

**User prompt (verbatim):** "Ok, commit, and then debug the geppetto ollama stuff. Keep a detailed diary as you work"

**Assistant interpretation:** Commit the operator command work first, then diagnose why the first real immutable TTC embedding build does not produce its expected artifact, preserving reproducible evidence.

**Inferred user intent:** Make the RAG laboratory's real 768D embedding baseline operational without hiding provider or cache failures behind repeated work.

**Commit (code):** `be83d09` — "feat: add immutable artifact commands"

### What I did

- Committed the immutable chunk and embedding Glazed commands.
- Upgraded Geppetto to `v0.13.5` in `182e786`, including the previously released embedding-cache YAML decoding fix.
- Confirmed `nomic-embed-text:latest` is installed and that the local Ollama HTTP API reports 768 dimensions.
- Added `scripts/02-geppetto-ollama-embedding-probe.go`, with zerolog logging and an optional direct immutable-service path.
- Ran the provider-only probe successfully; it returned model `nomic-embed-text`, configured dimensions 768, and vector dimensions 768.
- Ran the full command and the direct service probe against the 2,024-chunk set. Both created zero embedding rows and emitted no actionable error.

### Why

The immutable embedding set is the boundary between fixed chunk evidence and vector retrieval. It must be established with a real provider before latency, storage, and retrieval quality comparisons can be trusted.

### What worked

```bash
ollama list
curl -sS http://127.0.0.1:11434/api/embed ...
GOWORK=off go run ./ttmp/.../scripts/02-geppetto-ollama-embedding-probe.go --log-level debug
```

The last command produced a valid 768D vector through Geppetto, proving the provider configuration and single-item inference route are healthy.

### What didn't work

The following full-build attempts ended after roughly 23 seconds with no result row, no stored `embedding_sets` row, and no stored `immutable_embeddings` row:

```bash
GOWORK=off go run ./cmd/rag-eval embedding build-immutable ...
GOWORK=off go run ./ttmp/.../scripts/02-geppetto-ollama-embedding-probe.go --db data/rag-eval.db --chunk-set-id sha256:ef7...c9392
```

`SELECT COUNT(*) FROM embedding_sets;` and `SELECT COUNT(*) FROM immutable_embeddings;` both returned `0` after each run. The narrower direct-service reproduction rules out the Glazed command decoder as the primary cause.

### What I learned

- Geppetto v0.13.5 resolves this direct Ollama configuration and returns a valid single embedding.
- The failure appears in the batch/artifact-build path, before transaction commit. It is not evidence that the Ollama model is missing or has the wrong dimension.

### What was tricky to build

The service invokes `GenerateBatchEmbeddings` across all chunk batches before opening its persistence transaction. That design ensures no partial immutable set is published, but it also means a silent batch-path exit leaves no database artifact to inspect. The next diagnostic must isolate batch generation (one batch, then increasing batches) and capture process exit status/stderr externally.

### What warrants a second pair of eyes

- Inspect Geppetto's `GenerateBatchEmbeddings` behavior for the Ollama provider at v0.13.5, especially its behavior on a 16-text batch.
- Confirm whether the process is receiving a non-zero signal/exit status that the current tool wrapper did not display.
- Decide whether immutable embedding construction should persist batch progress as non-semantic run events while keeping the final set atomic.

### What should be done in the future

- Do not retry the full 2,024-chunk build blindly. First run a ticket-local batch-size sweep with externally captured exit status and provider-call logging.
- If a Geppetto batch defect is confirmed, create a focused ticket/issue with the probe and version evidence rather than adding a local retry shim.

### Code review instructions

Review `scripts/02-geppetto-ollama-embedding-probe.go`, then `internal/services/immutableembedding/service.go`. Confirm the direct provider-only path works and that no embedding-set rows exist after the reproduced failing build.

### Technical details

```text
single GenerateEmbedding -> succeeds (768D)
GenerateBatchEmbeddings inside immutable Build -> silent termination, no DB rows
```
