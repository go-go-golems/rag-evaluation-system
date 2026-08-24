---
Title: Simplify rag evaluation site package and Widget DSL
Ticket: RAG-WIDGET-SYSTEM-SIMPLIFICATION-2026-07-26
Status: active
Topics:
    - architecture
    - frontend
    - widget-dsl
    - code-quality
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Architecture and code-quality audit of the published React package and Widget DSL, with a deletion-oriented hard-cutover plan for legacy modules, adapters, manifests, schemas, package exports, and migration tools."
LastUpdated: 2026-07-26T21:15:00-04:00
WhatFor: "Use this ticket to plan and execute simplification of rag-evaluation-site and widget.dsl without retaining unsupported legacy wrappers."
WhenToUse: "Read before changing Widget protocol, package exports, adapters, registries, DSL builders, manifests, or compatibility behavior."
---

# Simplify rag evaluation site package and Widget DSL

## Overview

This ticket maps the complete JavaScript-authoring-to-React-rendering path and identifies complexity that no longer earns its maintenance cost. The review covers the published package, Widget IR, adapters, registry, default host, action dispatch, Goja module registration, typed specs and lowering, TypeScript declarations, manifests, codegen/checking, migration tools, examples, tests, and known consumers.

The proposed architecture retains the useful boundaries while completing hard cutovers. The target has one DSL module, one versioned page protocol, one browser parser, one action executor, one supported adapter registry, and a narrow intentional npm surface.

## Primary deliverables

- [Simplification analysis and hard cutover guide](./design-doc/01-rag-evaluation-site-and-widget-dsl-simplification-analysis-and-hard-cutover-guide.md)
- [Investigation diary](./reference/01-investigation-diary.md)
- [Widget DSL specialist audit](./reference/02-widget-dsl-audit-subagent.md)
- [Consumer and compatibility inventory](./reference/03-consumer-and-compatibility-inventory.md)
- [Frontend specialist audit](./reference/04-frontend-audit-subagent.md)

## Principal findings

- The default host can confirm non-server actions twice and duplicates server-action execution.
- Legacy split DSL modules remain implemented and globally registered after the `widget.dsl` hard cutover.
- Several public slot APIs serialize inert markers rather than renderable output.
- Component contracts are repeated across TypeScript, adapters, registries, YAML manifests, Go builders, descriptors, declarations, help, stories, and goldens.
- The npm root exposes fixtures, story data, presets, internal registries, and hidden CSS side effects.
- Compatibility branches and migration tools have no defined removal point.
- Story coverage is broad, but behavior tests are insufficient for safe mass deletion.

## Status

Current status: **active**. Analysis and reMarkable delivery are complete; implementation decisions remain proposed.

## Topics

- architecture
- frontend
- widget-dsl
- code-quality

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for investigation and delivery history.
