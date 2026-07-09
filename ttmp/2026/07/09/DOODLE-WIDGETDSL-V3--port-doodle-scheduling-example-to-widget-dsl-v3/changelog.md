# Changelog

## 2026-07-09

- Initial workspace created


## 2026-07-09

Ported examples/xgoja/doodle-site from legacy ui.dsl/data.dsl to widget.dsl v3, rebuilt the xgoja binary, and browser-smoked create/vote flows.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/examples/xgoja/doodle-site/verbs/doodle.js — v3 page/data/form migration
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/examples/xgoja/doodle-site/xgoja.v2.yaml — widget.dsl module selection


## 2026-07-09

Removed Doodle raw component escape hatches by adding typed widget.dsl v3 UI helpers and using schedule.availabilityPoll/pollSummary for poll views.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/examples/xgoja/doodle-site/verbs/doodle.js — raw-free v3 Doodle pages
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — new typed UI helpers


## 2026-07-09

Add calendar-widget visualization to the Doodle poll page so parseable offered slots render through widget.time.month and widget.time.week with participant availability labels.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/examples/xgoja/doodle-site/verbs/doodle.js — Builds calendar events from poll options and votes for MonthGrid/TimeGrid rendering


## 2026-07-09

Refactor Doodle jsverb into store/page/calendar/helper modules and improve calendar UX with marker-only month data, split-pane day details, selected day/slot navigation, and a scrollable week viewport.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/examples/xgoja/doodle-site/verbs/doodle.js — Now a small entrypoint that wires modules and HTTP routes
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/examples/xgoja/doodle-site/verbs/lib/calendar.js — Owns slot parsing, calendar markers, selected-day details, and calendar widget composition
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/examples/xgoja/doodle-site/verbs/lib/pages.js — Owns page composition and poll view models
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — Adds splitPane and week viewport-height DSL helpers


## 2026-07-09

Add a reusable ShareLink design-system molecule, Widget IR adapter, widget.dsl helper, and stories; use it for the Doodle poll share URL and refresh embedded assets.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/examples/xgoja/doodle-site/verbs/lib/pages.js — Doodle uses ShareLink for poll URLs
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/ShareLink/ShareLink.stories.tsx — Component Storybook coverage
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/ShareLink/ShareLink.tsx — Reusable share-link molecule
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/WidgetRenderer.share-link.stories.tsx — WidgetRenderer story coverage
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3.go — widget.ui.shareLink DSL helper


## 2026-07-09

Refine ShareLink to render as an inline link with a small icon-only clipboard button instead of an intrusive panel/button treatment.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/examples/xgoja/doodle-site/assets/public/index.html — Refreshed embedded app shell assets
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/ShareLink/ShareLink.module.css — Inline, non-intrusive ShareLink styling
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/ShareLink/ShareLink.tsx — Icon-only ShareLink copy control

