---
Title: go-go-parc corpus analysis design and implementation guide
Ticket: RAGEVAL-Vault-001
Status: active
Topics:
    - rag
    - evaluation
    - obsidian
    - sqlite
    - chunking
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/code/wesen/go-go-golems/go-go-parc/index.md
      Note: Vault corpus entry point
    - Path: abs:///home/manuel/code/wesen/go-go-golems/publish-vault/internal/parser/parser.go
      Note: Markdown frontmatter and wikilink parsing reference
    - Path: abs:///home/manuel/code/wesen/go-go-golems/publish-vault/internal/vault/vault.go
      Note: Vault discovery and ignore semantics
    - Path: repo://pkg/raglab/executor.go
      Note: Experiment execution and trace boundary
    - Path: repo://pkg/raglab/types.go
      Note: Immutable corpus and retrieval specification types
ExternalSources: []
Summary: Intern-facing design for importing a curated immutable snapshot of the go-go-parc Obsidian vault as a second RAG laboratory corpus.
LastUpdated: 2026-07-16T15:55:11.73275812-04:00
WhatFor: Define corpus boundaries, ingestion semantics, chunking, evaluation data, and experiment protocol for go-go-parc.
WhenToUse: Read before implementing vault ingestion, authoring corpus labels, or interpreting vault RAG experiments.
---


# go-go-parc corpus analysis design and implementation guide

## Executive Summary

`go-go-parc` is a technical Obsidian vault with approximately 1,037 Markdown
files and substantially different retrieval properties from the TTC WordPress
corpus. TTC tests product and FAQ retrieval. The vault will test structured
technical prose, headings, frontmatter, wikilinks, file references, repeated
project history, and long-form design explanations. It must be imported as a
separate immutable corpus snapshot, never appended to TTC.

The proposed first slice is a curated 100–200 note snapshot. It preserves
vault-relative paths, frontmatter, heading paths, resolved wikilinks, and note
hashes. It produces heading-aware chunks, BM25 and embedding artifacts, and a
hand-adjudicated evaluation dataset. The existing RAG laboratory then runs the
same vector, BM25, RRF, and reranker experiments against a new corpus shape.

## Problem Statement

The laboratory currently has one useful but narrow corpus. Its TTC evaluation
cards are document-level and its pages are mostly product/FAQ text. Results
there cannot establish that a chunking or retrieval policy works for project
documentation. The vault contains the material needed for a second benchmark,
but naïvely walking every Markdown file would introduce transient tickets,
generated files, duplicated reports, and unstable references.

The implementation must answer four questions:

- Which notes form a reproducible corpus snapshot?
- Which vault semantics become searchable metadata rather than flattened text?
- How are authoritative chunks and acceptable supporting evidence labeled?
- Which experiment comparisons isolate chunking, fusion, and reranking effects?

## Proposed Solution

### Architecture and data flow

```text
go-go-parc filesystem
  -> inclusion policy + ignore rules
  -> canonical note manifest (path, hash, frontmatter, links)
  -> immutable corpus snapshot
  -> heading-aware immutable chunks
  -> BM25 / embedding artifacts
  -> hand-adjudicated evaluation dataset
  -> immutable experiment specification and append-only run traces
```

The importer belongs in the RAG laboratory repository, because that repository
owns immutable corpus identities and evaluation runs. `publish-vault` is an
important reference for vault walking, ignore behavior, Markdown parsing, and
wikilink handling; it is not the source of experiment persistence.

### Source selection

Start with a manifest rather than an unrestricted recursive import.

- Include: `Projects/`, `Research/`, `Technical Reports/`, `RFCs/`, and
  selected durable playbooks.
- Exclude: `.obsidian/`, `node_modules/`, attachments/binaries, generated
  exports, `ttmp/_guidelines/`, and transient diary material unless explicitly
  sampled as a separate temporal-history experiment.
- Store each note's vault-relative path and SHA-256 content hash. A changed
  note creates a new snapshot; it never mutates a prior artifact.

### Note and chunk representation

Each document should retain the following metadata:

```json
{
  "source_kind": "obsidian-note",
  "vault_path": "Projects/2026/.../ARTICLE - ...md",
  "title": "...",
  "frontmatter": {"tags": ["rag", "evaluation"]},
  "outgoing_links": ["[[Geppetto]]"],
  "heading_path": ["3. Evaluation", "3.2 Metrics"]
}
```

Chunk at Markdown heading boundaries first. Split only sections exceeding the
configured size limit; prepend a compact title/heading breadcrumb to each
subchunk. Retain parent note ID, heading path, chunk ordinal, and rune range.
This permits parent-note collapse without losing a precise citation.

### Evaluation dataset protocol

The first 40–60 cards should be created manually from the frozen note
snapshot. A card contains query text, query family, authoritative note/chunk
IDs, acceptable substantial evidence, optional misleading near-neighbors, and
an adjudication rationale. Keep development, holdout, and regression cards
separate by evidence family: notes that repeat the same decision or report
must not cross partitions.

Suggested families:

- exact API/file lookup;
- architecture rationale;
- project-history/change questions;
- cross-note synthesis;
- provider/configuration boundary questions;
- negative or ambiguous queries;
- citation-location questions.

### Retrieval and reranking protocol

Run identical immutable configurations across raw heading chunks, BM25, vector,
weighted RRF, and RRF plus BGE/Qwen reranking. Add two vault-specific axes:

- heading-aware chunks versus fixed windows;
- raw evidence chunks versus parent-note collapse after reranking.

Record document and chunk nDCG/MRR/recall, candidate counts, reranker input
length, timing, storage, and citation correctness.

## API and File References

- `pkg/raglab/types.go` — immutable artifacts, representations, retrieval plan.
- `pkg/raglab/builder.go` — validation and canonical experiment fingerprint.
- `pkg/raglab/executor.go` — channel retrieval, RRF, reranking, traces.
- `internal/services/immutableretrieval/` — immutable BM25/vector retrieval.
- `internal/vault/vault.go` in `publish-vault` — note loading and vault rules.
- `internal/parser/parser.go` in `publish-vault` — Markdown/frontmatter/link
  parsing reference.

## Pseudocode

```text
for note in selectedVaultNotes(manifest):
  raw = read(note.path)
  parsed = parseMarkdown(raw)
  document = canonicalDocument(path, hash(raw), parsed.frontmatter, parsed.links)
  for section in splitByHeadings(parsed.body):
    for chunk in boundedChunks(section, maxTokens):
      writeImmutableChunk(document, headingBreadcrumb(section), chunk)
buildBM25(chunkSet)
buildEmbeddings(chunkSet, embeddingProfile)
registerEvaluationCards(snapshot, adjudicatedCards)
```

## Design Decisions

### Decision: curated snapshot before full-vault import

**Status:** proposed. A 100–200 note slice gives stable semantics and feasible
adjudication. Full-vault indexing follows only after the importer and labels
are proven.

### Decision: heading-first chunking

**Status:** proposed. Headings are durable author intent in this corpus. Fixed
windows remain an experimental baseline, not the default import semantics.

### Decision: source-grouped evaluation splits

**Status:** accepted. Repeated reports and project diaries create leakage if
related evidence appears in both development and holdout sets.

## Design Decisions

<!-- Document key design decisions and rationale -->

## Alternatives Considered

<!-- List alternative approaches that were considered and why they were rejected -->

## Implementation Plan

1. Inventory candidate roots and write a versioned inclusion manifest.
2. Implement a read-only vault discovery adapter and snapshot manifest test.
3. Implement Markdown/frontmatter/heading/wikilink extraction with fixtures.
4. Add heading-aware immutable chunking and fixed-window baseline mode.
5. Build first snapshot and retrieval artifacts.
6. Author and blind-validate 40–60 seed cards; freeze development/holdout.
7. Run baseline and reranker matrix; inspect traces and citations.
8. Decide whether summary representations or a reusable vault DSL are justified.

## Risks and Open Questions

- Wikilink resolution can be ambiguous where note names collide; store both raw
  target and resolution status.
- The vault includes historical reports that paraphrase one another; source
  grouping is required for credible evaluation.
- Code references may require repository-aware citation hydration later; first
  phase treats them as note text and metadata only.
- Do not claim general RAG quality from a small curated slice; keep TTC and
  vault metrics separate and report uncertainty.

## References

- `/home/manuel/code/wesen/go-go-golems/go-go-parc/index.md`
- `/home/manuel/code/wesen/go-go-golems/publish-vault/internal/vault/vault.go`
- `/home/manuel/code/wesen/go-go-golems/publish-vault/internal/parser/parser.go`
- `../RAGEVAL-RERANK-001--reranking-stage-for-the-immutable-ttc-rag-laboratory/`
