---
Title: RAG Design System Reference Sheet — Styling/Design + Widget/Interaction Hierarchy
Ticket: RAGEVAL-DS-RAGEVAL-QWEN36-HIGH
Status: active
Topics:
    - frontend-architecture
    - design-system
    - react
    - storybook
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ttmp/2026/06/07/RAGEVAL-DS-RAGEVAL-QWEN36-HIGH--rag-design-system-reference-sheet-styling-design-widget-interaction-hierarchy/reference/01-investigation-diary.md
      Note: Investigation diary with evidence-gathering and document authoring steps
    - Path: web/src/index.css
      Note: 'Global CSS: Tailwind import'
    - Path: web/src/styles/tokens.css
      Note: All design tokens (Mac palette + font role aliases)
    - Path: web/src/widgets/WidgetRenderer.tsx
      Note: WidgetRenderer — maps IR nodes to React components
    - Path: web/src/widgets/cellRenderers.tsx
      Note: CellSpec → ReactNode conversion pipeline
    - Path: web/src/widgets/ir.ts
      Note: Widget IR TypeScript types and factory functions
ExternalSources: []
Summary: ""
LastUpdated: 2026-06-07T12:26:01.336050274-04:00
WhatFor: ""
WhenToUse: ""
---



# RAG Design System Reference Sheet — Styling/Design + Widget/Interaction Hierarchy

## Overview

<!-- Provide a brief overview of the ticket, its goals, and current status -->

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- frontend-architecture
- design-system
- react
- storybook

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
