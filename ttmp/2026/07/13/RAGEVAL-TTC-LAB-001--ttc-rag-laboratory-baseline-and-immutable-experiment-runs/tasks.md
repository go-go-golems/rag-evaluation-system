# Tasks

## TODO

- [x] Confirm the multi-repository workspace layout and preserve clean worktree boundaries <!-- t:p9zb -->
- [x] Rebuild and validate data/ttc-wordpress-rag.sqlite and record source/output fingerprints <!-- t:cicz -->
- [x] Complete the intern-oriented baseline and immutable-run design package <!-- t:8031 -->
- [x] Implement a Glazed TTC corpus import command with deterministic source-balanced snapshot manifests <!-- t:shi3 -->
- [x] Replace mutable document text inputs with content-addressed document revisions and corpus snapshots <!-- t:3ydv -->
- [x] Define canonical JSON normalization and SHA-256 fingerprints for corpus, chunking, embedding, index, retrieval, and evaluation plans <!-- t:26xz -->
- [x] Implement immutable chunk plans, chunk sets, and exact source-range validation for fixed, sentence, and Markdown-heading chunkers <!-- t:rggc -->
- [x] Implement immutable embedding sets using one real Ollama 768D baseline profile and an offline deterministic test provider <!-- t:lbwm -->
- [ ] Implement content-addressed BM25 artifacts plus exhaustive vector retrieval and document-collapsed RRF hybrid retrieval <!-- t:crkp -->
- [ ] Create TTC baseline evaluation dataset v1 with at least 20 intent-stratified queries and named relevance judgments <!-- t:blhh -->
- [ ] Implement append-only experiment run events, terminal summaries, per-query traces, metrics, latency, cost, and storage accounting <!-- t:e15p -->
- [ ] Expose corpus snapshot, experiment specification, run, trace, and comparison APIs under /api/v1 <!-- t:yjha -->
- [ ] Build the React/Redux/RTK Query laboratory UI for spec editing, parallel retrieval inspection, run progress, and run comparison <!-- t:5upr -->
- [ ] Add Storybook states for reusable experiment, metric, trace, and artifact components <!-- t:bbbc -->
- [ ] Add offline unit/integration tests and a bounded end-to-end TTC baseline smoke run <!-- t:93bb -->
- [ ] Update CLI, API, schema, TTC handbook, RAG laboratory, and operator documentation <!-- t:xl0t -->
- [ ] Run Go, TypeScript, Biome, and end-to-end validation and record the first reproducible baseline results <!-- t:k43l -->
- [x] Add ticket-local script validation guidance and isolate multi-main script package behavior <!-- t:3fi7 -->
- [x] Validate immutable BM25 artifact construction against the TTC 2,024-chunk corpus <!-- t:did2 -->
- [ ] Add BM25 retrieval fixture coverage for hydration, rank ordering, and immutable artifact reuse <!-- t:kdge -->
- [x] Run 20 scored candidate queries through immutable BM25, exhaustive vector, and RRF fusion <!-- t:nkqy -->
- [x] Persist candidate retrieval traces with artifact identifiers, timings, and source citations <!-- t:6ira -->
- [x] Parse named relevance judgments into a machine-readable provisional evaluation dataset <!-- t:qz76 -->
- [x] Score raw BM25, vector, and RRF traces for recall, rank, and citation coverage <!-- t:2b5i -->
- [ ] Prepare a human adjudication packet required to freeze evaluation dataset v1 <!-- t:phdg -->
- [ ] Design append-only experiment-run schema with immutable specification and artifact references <!-- t:l21v -->
- [ ] Implement append-only experiment creation, lifecycle events, and immutable terminal summaries <!-- t:1i9w -->
- [ ] Persist per-query retrieval traces, metrics, latency, cost, and storage accounting <!-- t:oah5 -->
- [ ] Create bounded end-to-end baseline run from immutable artifacts and candidate evaluation data <!-- t:qcqg -->
- [ ] Add /api/v1 experiment specification, run, trace, and comparison endpoints <!-- t:x2wy -->
- [ ] Add Redux Toolkit Query client models for immutable runs and trace inspection <!-- t:1pe2 -->
- [ ] Build laboratory UI run launcher, progress view, trace inspector, and result comparison <!-- t:j5zm -->
- [ ] Add Storybook states for artifact, metric, trace, and comparison UI components <!-- t:nb2g -->
- [ ] Document immutable retrieval, evaluation adjudication, experiment runs, API, UI, and operator workflow <!-- t:iawa -->
- [ ] Run and record Go, TypeScript, formatting, and bounded end-to-end validation <!-- t:vico -->
