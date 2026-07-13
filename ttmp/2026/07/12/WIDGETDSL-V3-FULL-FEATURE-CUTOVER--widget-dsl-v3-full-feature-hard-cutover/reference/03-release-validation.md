---
Title: Widget DSL v3 Release Validation
Ticket: WIDGETDSL-V3-FULL-FEATURE-CUTOVER
Status: complete
Topics:
    - widget-dsl
    - release
    - xgoja
DocType: reference
Intent: ""
Owners: []
RelatedFiles:
    - Path: repo://examples/xgoja-widgetdsl-v3/Makefile
      Note: All 42 canonical pages generated-host gate
    - Path: repo://examples/xgoja/doodle-site/Makefile
      Note: Scheduling host integration gate
    - Path: repo://examples/xgoja/widget-site/verbs/sites.js
      Note: Interactive collection and overlay browser fixture
    - Path: repo://examples/xgoja/workshop-crm-site/Makefile
      Note: CRM host integration gate
    - Path: repo://scripts/test-widgetdsl-v3-sites.sh
      Note: Runs every generated-host release smoke suite
ExternalSources: []
Summary: ""
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: ""
WhenToUse: ""
---


# Widget DSL v3 Release Validation

## Result

The full-feature hard cutover is release-ready. Public providers expose only `widget.dsl`; canonical examples contain no legacy imports or raw component escape hatches; all generated xgoja reference hosts build and pass HTTP smoke suites; and browser validation confirms search, URL state, keyboard commands, FormDialog validation, successful submission, and focus restoration.

## Release gates

Passed on 2026-07-12/13:

- `go test ./... -count=1`
- `GOWORK=off go test ./... -count=1`
- `pnpm --dir packages/rag-evaluation-site typecheck`
- `pnpm --dir packages/rag-evaluation-site test:focused`
- `pnpm --dir packages/rag-evaluation-site build`
- `pnpm --dir packages/rag-evaluation-site build-storybook`
- `pnpm --dir packages/rag-evaluation-site pack:smoke`
- `pnpm --dir packages/rag-evaluation-site consumer:smoke`
- `make widgetdsl-sites-smoke`
- `go run ./cmd/widgetdsl-migration-checker --root . --fail-on-findings -- examples packages README.md pkg/widgetdsl/testdata/v3/examples`
- `docmgr doctor --ticket WIDGETDSL-V3-FULL-FEATURE-CUTOVER --stale-after 30`

## Generated-host coverage

`make widgetdsl-sites-smoke` validates four independently generated binaries:

1. `examples/xgoja/widget-site`: interactive DataTable, progressive search, deterministic server filtering, pagination metadata, keyboard command IR, FormDialog IR, action refresh, and 422 field errors.
2. `examples/xgoja/doodle-site`: scheduling tables, poll matrices, month/week calendar views, share links, and forms.
3. `examples/xgoja/workshop-crm-site`: funnel metrics, board engine, record fields, generic activity feed, and persisted deal movement.
4. `examples/xgoja-widgetdsl-v3`: every committed v3 example page (42 pages), schema/root integrity, and action endpoint dispatch.

Every host rebuilds or embeds the canonical current SPA. Generated asset folders are excluded from Biome source formatting so rebuilds remain byte-reproducible.

## Browser validation

The generated widget-site binary was exercised in Chromium at desktop and 390×844 viewports. Verified behavior:

- URL-backed `?q=Arbor` filtering returns one row and reports one result.
- Row activation and `t` open the note FormDialog.
- A whitespace-only note reaches the server, receives HTTP 422, keeps the dialog open, and renders both field and form errors.
- A valid note closes the dialog and restores focus to the originating row.
- Keyboard command status changes refresh the server-backed table.
- Narrow layouts retain usable table/search/pagination content and horizontally scroll the app navigation.

## Known non-blocking warnings

- Storybook reports a large preview chunk and plugin timing warnings.
- Biome reports intentional `${...}` binding strings and legacy CSS specificity/`!important` warnings in source checks.
- Browser devtools logs the intentional 422 validation response as a failed network resource during the negative-path scenario.

None of these warnings indicate a failed release gate.
