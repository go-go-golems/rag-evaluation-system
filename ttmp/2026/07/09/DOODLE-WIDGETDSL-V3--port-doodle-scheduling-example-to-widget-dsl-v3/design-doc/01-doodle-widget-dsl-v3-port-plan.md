---
Title: Doodle widget.dsl v3 Port Plan
Ticket: DOODLE-WIDGETDSL-V3
Status: active
Topics:
    - widget-dsl
    - xgoja
    - sqlite
    - doodle
DocType: design-doc
Intent: short-term
Owners: []
RelatedFiles:
    - Path: repo://examples/xgoja/doodle-site/verbs/doodle.js
      Note: existing Doodle app uses legacy ui.dsl/data.dsl and should be ported to widget.dsl
    - Path: repo://examples/xgoja/doodle-site/xgoja.v2.yaml
      Note: module selection should move to widget.dsl for rag-widget-site
Summary: Plan to port the Doodle-style scheduling example from legacy split widget modules to widget.dsl v3.
LastUpdated: 2026-07-09
WhatFor: Use when resuming the Doodle example migration after the hardening ticket.
WhenToUse: Before editing examples/xgoja/doodle-site.
---

# Doodle widget.dsl v3 Port Plan

## Goal

Port the existing Doodle-style SQLite + xgoja demo from legacy `ui.dsl` / `data.dsl` imports to the new `widget.dsl` v3 module.

## Current state

The `DOODLE-1` ticket produced a working browser-verified example under `examples/xgoja/doodle-site`. It uses:

- SQLite tables for polls, options, participants, and votes;
- planned-route Express handlers;
- native form POST + 303 redirects;
- legacy `ui.dsl` and `data.dsl` helpers.

## Desired state

- `xgoja.v2.yaml` selects only `widget.dsl` from `rag-widget-site`.
- `verbs/doodle.js` imports `const widget = require("widget.dsl")`.
- Pages use v3 page/section/data/action helpers.
- Scheduling/time helpers are used where they improve the poll/results/calendar rendering.
- Browser validation covers create poll, cast vote, index metrics, and poll results.

## Suggested implementation order

1. Change module selection and import `widget.dsl`.
2. Port the index page to `widget.page`, `widget.ui`, and `widget.data.collection`.
3. Port create and poll pages, keeping native forms where they are the simplest reliable submission path.
4. Replace hand-authored result tables with v3 data/schedule helpers where available.
5. Run migration checker and browser smoke.
6. Update the original `DOODLE-1` diary/changelog with a cross-reference if needed.
