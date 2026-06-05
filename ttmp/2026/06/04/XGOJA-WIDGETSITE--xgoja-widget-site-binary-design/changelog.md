# Changelog

## 2026-06-04

- Initial workspace created


## 2026-06-04

Created XGOJA-WIDGETSITE design package: captured xgoja help, inspected providers/assets/db/runtime code, ran upstream and ticket-local experiments, and wrote design guide, research logbook, and diary

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/04/XGOJA-WIDGETSITE--xgoja-widget-site-binary-design/design-doc/01-xgoja-widget-site-binary-analysis-and-implementation-guide.md — Primary design guide
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/04/XGOJA-WIDGETSITE--xgoja-widget-site-binary-design/reference/01-research-logbook.md — Research logbook
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/04/XGOJA-WIDGETSITE--xgoja-widget-site-binary-design/scripts/01-current-xgoja-widgetsite-experiment/xgoja.yaml — Combined xgoja experiment


## 2026-06-04

Uploaded XGOJA-WIDGETSITE design bundle to reMarkable after fixing Mermaid diagram syntax

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/04/XGOJA-WIDGETSITE--xgoja-widget-site-binary-design/design-doc/01-xgoja-widget-site-binary-analysis-and-implementation-guide.md — Corrected Mermaid diagram and uploaded bundle
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/04/XGOJA-WIDGETSITE--xgoja-widget-site-binary-design/reference/02-diary.md — Recorded validation and upload step


## 2026-06-04

Implemented first xgoja WidgetRenderer slice: stable widgetsite provider, generated binary example, and curl smoke for health/static/Widget IR endpoints (code commit 70f30b1f53a5ed6530b1780431ee2b3bceddcd91).

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/examples/xgoja/widget-site/xgoja.yaml — Example build spec exercises local package replace and embedded assets
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/pkg/xgoja/providers/widgetsite/provider.go — Provider registers widget.dsl and rag.dsl


## 2026-06-04

Embedded the packaged React WidgetRenderer app in the xgoja example and validated it with curl plus Pi Playwright smoke (0 console errors, code commit 36cd6ea98fbc24d9db6d8f1b077c80099b0a3431).

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/examples/xgoja/widget-site/assets/public/index.html — React app entry served from /static/
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/examples/xgoja/widget-site/verbs/sites.js — API page aliases and favicon route for clean browser smoke


## 2026-06-05

Added provider-bundled Glazed help docs for widget.dsl getting started and JavaScript API reference, and verified they render in the generated xgoja binary help command (code commit 5cdcf5b42fec3df59b35462bd06943e219508b6d).

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/pkg/xgoja/providers/widgetsite/doc/01-widget-dsl-getting-started.md — Getting started tutorial
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/pkg/xgoja/providers/widgetsite/doc/02-widget-dsl-js-api-reference.md — JavaScript API reference

