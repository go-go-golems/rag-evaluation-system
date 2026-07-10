# Changelog

## 2026-07-09

- Initial workspace created


## 2026-07-09

Step 1: created an intern-ready workshop CRM vertical-slice guide with evidence-backed API, architecture, implementation, and validation plan.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/09/RAGEVAL-WORKSHOP-CRM-IMPLEMENTATION--workshop-crm-widget-dsl-implementation-guide-and-vertical-slice/design-doc/01-intern-guide-workshop-crm-widget-dsl-vertical-slice.md — Primary implementation guide
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/09/RAGEVAL-WORKSHOP-CRM-IMPLEMENTATION--workshop-crm-widget-dsl-implementation-guide-and-vertical-slice/reference/01-implementation-diary.md — Chronological implementation record


## 2026-07-09

Step 2: added widget.crm opaque field/pipeline builders, CRM view helpers, intents, declarations, descriptors, and a golden fixture (commit 196cb20800c7d3893daffe6aca37fa9682e0a251).

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript.go — CRM API declarations
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3_crm.go — CRM runtime implementation


## 2026-07-09

Step 3: added and browser-validated SQLite xgoja workshop CRM reference host, plus CRM palette style sets required by funnel/activity rendering (commit 0d81a70b594cfea9a1884d6cfc363c27c2fdb9d2).

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/examples/xgoja/workshop-crm-site/verbs/lib/store.js — SQLite application data
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3_crm.go — CRM Widget IR palettes


## 2026-07-09

Step 4: fixed pipeline card opening and durable drag/drop stage changes using BoardEngine cardId/to context and a CRM action route (commit 9b70f4af07fb89c2ef536348e02b0adbbdd5e478).

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/examples/xgoja/workshop-crm-site/verbs/workshop-crm.js — Stage-move endpoint
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3_crm.go — CRM action payload contract


## 2026-07-09

Step 5: addressed four Widget DSL action-contract review findings and fixed logcopter/govulncheck CI failures with generated logger output and Go 1.26.5 (commit 8984e12e44ebbae7373c595af3dcc2927ff85d45).

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/go.mod — Standard-library security remediation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/actions.ts — Action binding/event remediation


## 2026-07-10

Step 6: exported CRM presets publicly and defaulted sparse CRM funnel counts to zero with runtime regression coverage (commit 9126ccaa418f15270694c5ad9cbd50fd400f062c).

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/presets/index.ts — Publish CRM presets
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v3_crm.go — Prevent NaN SegmentedBar widths


## 2026-07-10

Step 7: published v3-first embedded Glazed help cookbook/API inventory, CRM and action-contract guidance, and README links (commit b49d497332eb6b29b84456f2602b25a619e187d1).

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/xgoja/providers/widgetsite/doc/04-widget-dsl-v3-examples.md — Runnable authoring documentation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/xgoja/providers/widgetsite/doc/05-widget-dsl-v3-api-reference.md — Installed public v3 API inventory

