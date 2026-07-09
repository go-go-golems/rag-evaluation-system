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

