# Changelog

## 2026-06-05

- Initial workspace created


## 2026-06-05

Added an example-local devctl workflow for rebuilding, launching, supervising, and smoking examples/xgoja/widget-site while keeping plugin stdout protocol-clean.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/examples/xgoja/widget-site/.devctl.yaml — Devctl plugin configuration
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/examples/xgoja/widget-site/devctl/widget-site.py — Devctl plugin implementation


## 2026-06-05

Captured visual evidence with css-visual-diff, including current Widget DSL action page, Storybook reference pages, overlay screenshots, computed CSS artifacts, and the root-cause finding that standalone package theme tokens do not define the --mac variables used by copied components.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/scripts/03-capture-visual-evidence.js — Reusable evidence capture verb
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-02/README.md — Successful visual evidence run summary


## 2026-06-05

Wrote the intern-facing Widget DSL visual quality analysis and implementation guide, covering architecture, evidence, root causes, design decisions, phased fixes, testing, and API/file references.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/design-doc/01-widget-dsl-visual-quality-analysis-and-implementation-guide.md — Primary guide


## 2026-06-05

Validated WIDGETDSL-VISUAL-QUALITY and uploaded the guide bundle to reMarkable at /ai/2026/06/05/WIDGETDSL-VISUAL-QUALITY.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/design-doc/01-widget-dsl-visual-quality-analysis-and-implementation-guide.md — Uploaded primary guide
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/reference/01-investigation-diary.md — Uploaded investigation diary
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-02/01-visual-evidence-summary.md — Uploaded visual evidence summary


## 2026-06-05

Implemented Phase 1 token bridge in the standalone rag-evaluation-site theme, rebuilt/synced xgoja widget-site assets, passed package and generated-site smokes, and captured post-bridge visual evidence showing panels/buttons now receive component chrome.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/theme.css — Token bridge implementation
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-token-bridge-02/01-visual-evidence-summary.md — Post-fix visual evidence


## 2026-06-05

Implemented Phase 2 default shell/page chrome for the standalone Widget DSL app, rebuilt/synced xgoja assets, passed generated-site and browser action smokes, and captured post-shell visual evidence showing root background and padding now render correctly.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/app/App.tsx — Default shell implementation
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/app/app.css — Default shell styling
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-shell-01/01-visual-evidence-summary.md — Post-shell visual evidence


## 2026-06-05

Implemented Phase 3 semantic Widget DSL recipes, refactored the xgoja action dashboard to use them, documented the API, and validated with Go tests, xgoja smoke, devctl smoke, generated help checks, and browser action smoke.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/examples/xgoja/widget-site/verbs/sites.js — Recipe-backed showcase
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/pkg/widgetdsl/module.go — Recipe implementation

