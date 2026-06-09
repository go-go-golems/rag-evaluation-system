---
Title: Color Palette Theming for Context Window Diagrams
Ticket: CTX-COLOR-PALETTE
Status: active
Topics:
    - frontend
    - theming
    - context-window
    - visualization
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: packages/rag-evaluation-site/src/components/atoms/ContextStyleSwatch/ContextStyleSwatch.tsx
      Note: Generic visual-style swatch replacing ContextKindSwatch
    - Path: packages/rag-evaluation-site/src/components/molecules/ContextBudgetBar/ContextBudgetBar.module.css
      Note: Budget bar generic pattern CSS
    - Path: packages/rag-evaluation-site/src/components/molecules/ContextStackDiagram/ContextStackDiagram.module.css
      Note: Stack diagram generic pattern CSS
    - Path: packages/rag-evaluation-site/src/components/molecules/ContextStripDiagram/ContextStripDiagram.tsx
      Note: Strip diagram styleKey + styleSet rendering
    - Path: packages/rag-evaluation-site/src/components/molecules/ContextTreemap/ContextTreemap.module.css
      Note: Treemap generic pattern CSS
    - Path: packages/rag-evaluation-site/src/components/organisms/ContextDiagramPanel/ContextDiagramPanel.tsx
      Note: Panel passes required styleSet to diagrams and legend
    - Path: packages/rag-evaluation-site/src/context/fixtures.ts
      Note: StyleKey fixtures and example style sets
    - Path: packages/rag-evaluation-site/src/context/styles.ts
      Note: ContextStyleSet model and palette helpers
    - Path: packages/rag-evaluation-site/src/context/types.ts
      Note: Core styleKey and ContextStyleSet-related types
ExternalSources: []
Summary: ""
LastUpdated: 2026-06-09T16:38:36.060938613-04:00
WhatFor: ""
WhenToUse: ""
---













# Color Palette Theming for Context Window Diagrams

## Overview

<!-- Provide a brief overview of the ticket, its goals, and current status -->

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- frontend
- theming
- context-window
- visualization

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
