# Changelog

## 2026-07-16

- Initial workspace created


## 2026-07-16

Created the issue #25 research ticket and proposed page.shortcuts: typed page commands that reuse ActionSpec, dispatch through the app shell, preserve nested keyboard precedence, and include accessibility safety requirements.

### Related Files

- /home/manuel/code/wesen/go-go-golems/rag-evaluation-system/packages/rag-evaluation-site/src/app/App.tsx — Proposed page-owned listener and existing action dispatch integration point
- /home/manuel/code/wesen/go-go-golems/rag-evaluation-system/pkg/widgetdsl/v3.go — Proposed PageBuilder and typed shortcut builder integration points


## 2026-07-16

Implemented the complete page.shortcuts stack: typed specs and validation, v3 builder/declarations/help/goldens, React matching and safety guards, accessible help plus disable preference, embedded SPA refresh, and Upwork Triage Y/N/S consumer wiring.

### Related Files

- /home/manuel/code/wesen/claw-stuff/upwork/verbs/lib/pages.js — Triage consumer binds Y/N/S to visible-button actions
- /home/manuel/code/wesen/go-go-golems/rag-evaluation-system/packages/rag-evaluation-site/src/app/App.tsx — Page listener dispatch, preference, help, and ARIA integration
- /home/manuel/code/wesen/go-go-golems/rag-evaluation-system/pkg/widgetdsl/v3.go — Typed shortcut authoring and page IR lowering


## 2026-07-16

Ticket closed


## 2026-07-16

Expanded public docs with shortcut API concepts, runnable Y/N/S tutorial, React host integration and preference behavior, and Upwork operator keyboard instructions.

### Related Files

- /home/manuel/code/wesen/claw-stuff/upwork/README.md — Operator-facing Triage keyboard workflow
- /home/manuel/code/wesen/go-go-golems/rag-evaluation-system/packages/rag-evaluation-site/README.md — React host responsibilities and custom integration
- /home/manuel/code/wesen/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite/doc/02-widget-dsl-js-api-reference.md — Author-facing API and troubleshooting reference
- /home/manuel/code/wesen/go-go-golems/rag-evaluation-system/pkg/xgoja/providers/widgetsite/doc/04-widget-dsl-v3-examples.md — Runnable shortcut tutorial and host safety guidance


## 2026-07-16

Committed the complete upstream shortcut implementation and public documentation as 36bafb06a91a3374917c5b3e8e4dec53c2015ff7 (widgetdsl: add page keyboard shortcuts).

### Related Files

- /home/manuel/code/wesen/go-go-golems/rag-evaluation-system/packages/rag-evaluation-site/src/app/App.tsx — React shortcut runtime integration committed in 36bafb0
- /home/manuel/code/wesen/go-go-golems/rag-evaluation-system/pkg/widgetdsl/v3.go — Shortcut builder and page lowering committed in 36bafb0

## 2026-07-16

Committed the Upwork Triage consumer, rebuilt embedded assets, and operator documentation as 9905e09c1bf97d736eff7cddbdf61151675ff1e0 (upwork: add triage keyboard shortcuts).

### Related Files

- /home/manuel/code/wesen/claw-stuff/upwork/assets/public — Rebuilt shortcut-enabled SPA committed in 9905e09
- /home/manuel/code/wesen/claw-stuff/upwork/verbs/lib/pages.js — Y/N/S consumer bindings committed in 9905e09
