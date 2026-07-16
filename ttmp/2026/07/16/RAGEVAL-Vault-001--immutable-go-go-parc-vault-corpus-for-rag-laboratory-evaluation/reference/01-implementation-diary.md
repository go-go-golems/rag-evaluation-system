---
Title: Implementation diary
Ticket: RAGEVAL-Vault-001
Status: active
Topics:
    - rag
    - evaluation
    - obsidian
    - sqlite
    - chunking
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: ""
LastUpdated: 2026-07-16T15:55:12.371612601-04:00
WhatFor: ""
WhenToUse: ""
---

# Implementation diary

## Goal

Create a second immutable RAG evaluation corpus from the go-go-parc Obsidian
vault, with headings, links, and citations preserved as first-class evidence.

## Context

TTC has established the laboratory and yielded a promising BGE reranker result,
but it is product/FAQ data. The vault supplies long technical documents and is
therefore the next corpus needed to test chunking and citation decisions.

## Step 1: Ticket creation and initial corpus reconnaissance

### What I did

- Created `RAGEVAL-Vault-001` using `docmgr --root ttmp`.
- Counted 1,037 Markdown files and approximately 253 MB in go-go-parc.
- Inspected `publish-vault` as the source of parsing, ignore, and vault
  navigation conventions.
- Wrote the initial design, task plan, and dispatched a bounded read-only
  inventory/evaluation-seed subagent.

### What worked

- The vault has enough durable projects, reports, research notes, and
  playbooks for a curated corpus that differs materially from TTC.

### What did not work

- A first search did not find SQLite ingestion in `publish-vault`; its useful
  role is filesystem vault parsing/publishing, not immutable experiment storage.

### What was tricky to build

The vault contains durable reports and temporary material side by side. The
design therefore chooses an explicit inclusion manifest rather than assuming
all Markdown is benchmark data.

## 2026-07-16: first deterministic candidate inventory

I added the first-snapshot inclusion manifest and a read-only inventory script.
The selected roots currently yield 108 Markdown candidates: 96 current project
notes, 8 playbooks, and 4 institute guidelines. The sorted candidate path list
hash is `sha256:1d81752cc302bfbefa9dc010477283313dc63609f3d98398202f85bf55436b95`.
This is within the intended 100–200-note range and excludes ticket workspaces,
attachments, generated material, and private scratch trees.

The inventory is still a candidate set rather than an immutable corpus. The
next implementation task is to persist the path list, note hashes, parser
version, and extracted heading/link metadata under one snapshot ID.

### Code review instructions

Start with the design document's source selection and heading-chunking
decisions, then compare the referenced `publish-vault` parsing code with the
RAG laboratory's immutable artifact boundaries.

## Quick Reference

<!-- Provide copy/paste-ready content, API contracts, or quick-look tables -->

## Usage Examples

<!-- Show how to use this reference in practice -->

## Related

<!-- Link to related documents or resources -->
