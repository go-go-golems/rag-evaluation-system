---
Title: Widget DSL v3 Full Feature Hard Cutover
Ticket: WIDGETDSL-V3-FULL-FEATURE-CUTOVER
Status: active
Topics:
    - widget-dsl
    - ui-dsl
    - widget-ir
    - goja
    - xgoja
    - react
    - design-system
    - frontend-architecture
    - typescript
    - intern-guide
DocType: index
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://packages/rag-evaluation-site/src/components/atoms/TextareaInput/TextareaInput.widget.yaml
      Note: Stale entry manifest migrated to adapter schema
    - Path: repo://packages/rag-evaluation-site/src/components/molecules/ContextGroupedStripDiagram/ContextGroupedStripDiagram.widget.yaml
      Note: Context manifest metadata completed
    - Path: repo://packages/rag-evaluation-site/src/components/organisms/ContextTurnPagerPanel/ContextTurnPagerPanel.widget.yaml
      Note: Context manifest metadata completed
    - Path: repo://pkg/xgoja/providers/widgetsite/doc/05-widget-dsl-v3-api-reference.md
      Note: Descriptor-generated direct API inventory updated by commit f208624
    - Path: repo://schema/dsl-modules.yaml
      Note: Transitional manifest module catalog repaired in commit 217ad13
    - Path: repo://ttmp/2026/06/04/XGOJA-WIDGETSITE--xgoja-widget-site-binary-design/scripts/01-current-xgoja-widgetsite-experiment/widgetprovider/provider.go
      Note: Historical provider loader drift fixed in commit 2017908
ExternalSources: []
Summary: Research and implementation design workspace for completing Widget DSL v3, migrating all first-party consumers, and removing the unreleased legacy split-module surface through a coordinated hard cutover.
LastUpdated: 2026-07-12T19:50:00-04:00
WhatFor: Track the full-feature parity analysis, target language, implementation phases, migration evidence, and reMarkable delivery for Widget DSL v3.
WhenToUse: Use before implementing or reviewing v3 parity, namespace reorganization, first-party migration, or legacy module deletion.
---



# Widget DSL v3 Full Feature Hard Cutover

## Overview

This ticket defines the hard-cutover path from the current split Widget DSL modules and partial v3 surface to one complete, typed, composable, and opinionated `widget.dsl` language. It covers generic UI/content parity, collections, keyboard commands, dialogs, progressive search, pagination, activity timelines, domain views, descriptors, declarations, examples, browser state, provider migration, and release validation.

V3 has not been released publicly, so the plan deliberately permits API renaming and namespace reorganization. Compatibility aliases are not part of the final design.

## Key Links

- [Primary intern implementation guide](design-doc/01-widget-dsl-v3-full-feature-analysis-design-and-intern-implementation-guide.md)
- [Legacy-to-v3 parity inventory](reference/01-legacy-to-v3-feature-parity-inventory.md)
- [Investigation diary](reference/02-investigation-diary.md)
- [Generated runtime inventory](sources/01-generated-runtime-inventory.md)
- [V3 example migration findings](sources/02-v3-example-migration-check.txt)

## Status

Current status: **active**

## Topics

- widget-dsl
- ui-dsl
- widget-ir
- goja
- xgoja
- react
- design-system
- frontend-architecture
- typescript
- intern-guide

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design-doc/ - Architecture, language, and implementation design
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
