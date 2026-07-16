---
Title: go-go-parc First Snapshot Inclusion Manifest
Ticket: RAGEVAL-Vault-001
Type: Reference
Status: draft
Created: 2026-07-16
---

# Scope

The first go-go-parc snapshot is a curated, immutable input for RAG experiments. It is not a claim that every vault note is evaluation-grade. The initial extraction should contain approximately 100–200 Markdown notes, selected by explicit path rules and recorded as a sorted path list before parsing or chunking.

## Include roots

- `Projects/2026/07/` for current project notes and implementation reports.
- `Research/playbooks/` for operational procedures that can support factual queries.
- `Research/Institute/Guidelines/` for stable technical guidance.
- Selected historical project notes under `Projects/2026/05/` when they contain durable system decisions or reproducible findings.

## Exclude roots and file classes

- `.obsidian/`, `.git/`, `.pi/`, `.trash/`, and dependency trees such as `node_modules/`.
- Attachments and binary media; the first snapshot is Markdown-only.
- `ttmp/` ticket workspaces, raw transcript exports, generated logs, drafts, and private scratch notes.
- Duplicate exports whose canonical source is already included.

## Required manifest fields

Each included note records:

- vault-relative path and byte length;
- SHA-256 content hash;
- title, frontmatter, tags, and aliases as parsed from the original Markdown;
- heading tree with byte or line offsets;
- outbound wikilinks and their resolved target paths, when resolution succeeds;
- snapshot identifier and parser version.

The manifest is content-addressed. A change to note bytes, parser version, inclusion rules, or link resolution must create a new snapshot ID rather than mutating an existing one.

## Acceptance checks

1. The generated path list is sorted and contains 100–200 files.
2. Every path is within an include root and outside every exclusion rule.
3. Every file has a stable hash and parses without silently dropping frontmatter or headings.
4. The manifest and extraction artifacts are committed together; downstream experiments reference the snapshot ID.
5. A second run against unchanged bytes produces the same path-list and snapshot hash.

The companion inventory script is deliberately read-only and prints the candidate path list. It must be reviewed before its output is promoted to the immutable snapshot.
