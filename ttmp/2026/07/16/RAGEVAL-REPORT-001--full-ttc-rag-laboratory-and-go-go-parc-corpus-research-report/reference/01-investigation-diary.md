---
Title: Investigation Diary
Ticket: RAGEVAL-REPORT-001
Status: active
Topics:
    - rag-eval
    - ttc
    - corpus
    - embeddings
    - search
    - evaluation
    - obsidian
    - reranking
    - workflow
    - intern-guide
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-16T19:15:55.11169759-04:00
WhatFor: ""
WhenToUse: ""
---

# Investigation Diary

## Goal

Produce a long-form, textbook-style project report for the Obsidian vault that explains the TTC RAG laboratory, the go-go-parc corpus investigation, the tools and scripts used, measured findings, implementation boundaries, and the next engineering steps. The report must be useful to a new intern and must be copied to the vault, uploaded to reMarkable, and committed and pushed in the vault repository.

## Context

The report was requested after a long sequence of TTC RAG, immutable experiment, evaluation-corpus, reranker, vault-ingestion, and fluent-JavaScript-playground work. The current working repository is `/home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system`. The primary evidence is in the `RAGEVAL-TTC-LAB-001` ticket and in repository documentation and scripts. The destination vault is `/home/manuel/code/wesen/go-go-golems/go-go-parc`.

## Quick Reference

### User request

> Ok, write a detailed project report for the obsidian vault as a deep dive technical analysis blog post using a textbook writing style (no analogies, see skill). Commit and push the bsidian vault when done (go-go-parc vault). about all the work we have been doing (consult older diaries, etc...), our findings, how it's built up, which tools and scripts we use, pseudocode / code snippets. It can be really really really long, this is a huge and important project.

### Evidence inventory

- TTC source: `data/ttc-wordpress-rag.sqlite`.
- Immutable catalog: `data/rag-eval.db`.
- Frozen snapshot: `sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409`.
- Chunk set: `sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392`.
- BM25 index: `sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691`.
- Embeddings: `sha256:2665c5249b8352ce6904fc00c934534dd179f3eeef0a6a75429a9034be0e03e0`.
- Candidate dataset: `candidate:ttc-expansion-v0`, 148 registered cards.
- Development metrics: vector Recall@1 `0.8750`, hybrid relevant recall@10 `0.8947`, mean latency `173 ms`.

### Diary entries

#### 2026-07-16 — Ticket and scope created

Created `RAGEVAL-REPORT-001` with topics for RAG evaluation, TTC, corpus work, embeddings, search, evaluation, Obsidian, reranking, workflow, and intern guidance. Added five tasks: inventory prior artifacts, collect and relate sources, write the report, copy it to the vault, and validate/upload/commit/push.

#### 2026-07-16 — Prior artifacts inventoried

Inspected the TTC ticket design documents, adjudication packets, 70/80/50 expansion batches, 240-card split protocol, source validation, candidate registration, and development results. Inspected the vault design and inclusion manifest, RAG DSL design/API, reranker design/results, Ollama tunnel playbook, and earlier diaries. Copied 22 local evidence documents into the report ticket's `sources/` directory so the report has an explicit local source set.

#### 2026-07-16 — Report drafted

Drafted a 40k+ character report covering architecture, corpus normalization, chunking, representations, BM25/vector/RRF retrieval, parent collapse, citation hydration, immutable runs, evaluation cards, reranking, the JavaScript builder target, the web UI, cost/latency accounting, the vault corpus, failure records, and a staged intern checklist. Included Mermaid diagrams, pseudocode, API contracts, actual artifact hashes, and measured development metrics. Explicitly marked the current results provisional and stated that the authoring agent is GPT-5.6 (Sol).

#### 2026-07-16 — Earlier implementation details retained

The TTC trace driver had previously required several repairs: support for comma-separated card paths, heading cards, indented `- id:` YAML cards, inline mappings, and inline quoted queries. A first invocation produced zero traces; another produced only 70. The final development invocation produced 150 traces. These events are recorded in the report as reproducibility lessons rather than hidden cleanup.

#### 2026-07-16 — Validation and delivery (pending)

Next actions are to run `docmgr doctor`, relate the report and sources to the ticket, upload the bundle with `remarquee`, copy the full report into `Research/2026/07/16/` in go-go-parc, inspect the vault diff, commit only the new article, and push the vault branch. If a command fails because of sandbox permissions, retry with an explicit escalation and record the exact error here.

## Usage Examples

To reproduce the local report inputs:

```sh
cd /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system
docmgr --root ttmp doctor --ticket RAGEVAL-REPORT-001
find ttmp/2026/07/16/RAGEVAL-REPORT-001--*/sources -type f | sort
```

To inspect the development result, read `sources/ttc-development-results.md` and the original ticket result document. To reproduce the local-model endpoint, follow `sources/ollama-tunnel-playbook.md`; do not assume a tunnel is alive merely because the model name is configured.

## Related

The full report is the sibling design document `design-doc/01-full-ttc-rag-laboratory-and-go-go-parc-corpus-research-report.md`. The source inventory is deliberately copied rather than summarized so a reviewer can inspect the evidence without searching the entire workspace. The TTC implementation diary remains authoritative for individual coding steps; this diary records the report-production process and its delivery state.
