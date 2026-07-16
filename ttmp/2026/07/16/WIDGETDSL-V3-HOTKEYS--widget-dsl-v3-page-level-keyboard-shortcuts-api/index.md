---
Title: Widget DSL v3 page-level keyboard shortcuts API
Ticket: WIDGETDSL-V3-HOTKEYS
Status: complete
Topics:
    - widget-dsl
    - ui-dsl
    - widget-ir
    - frontend
    - react
    - xgoja
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources:
    - https://github.com/go-go-golems/rag-evaluation-system/issues/25
Summary: Design and implementation plan for safe page-level keyboard shortcuts in Widget DSL v3.
LastUpdated: 2026-07-16T19:21:27.283029498-04:00
WhatFor: 'Track the upstream keyboard shortcut contract requested by issue #25''s Upwork Triage use case.'
WhenToUse: Use when reviewing, implementing, testing, or consuming page-level Widget DSL shortcuts.
---


# Widget DSL v3 page-level keyboard shortcuts API

## Overview

Issue [#25](https://github.com/go-go-golems/rag-evaluation-system/issues/25) identified that Widget
DSL v3 only supports table-scoped keyboard behavior. This ticket designs a serializable page-level
shortcut API for card and workflow pages such as Upwork Triage.

The current proposal is `page.shortcuts((keys) => keys.bind(...))`. Each binding carries a stable
ID, logical key, label, optional modifiers, and an existing `ActionSpec`; the React page shell owns
safe event matching and dispatch.

## Key documents

- [Page-level keyboard shortcuts API design](design-doc/01-page-level-keyboard-shortcuts-api-design.md)
- [Investigation diary](reference/01-investigation-diary.md)
- [Tasks](tasks.md)
- [Changelog](changelog.md)

## Status

Current status: **active**. The API is proposed and ready for review; implementation tasks remain
open.

## Topics

- widget-dsl
- ui-dsl
- widget-ir
- frontend
- react
- xgoja
