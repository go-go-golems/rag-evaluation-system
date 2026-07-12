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
RelatedFiles: []
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
