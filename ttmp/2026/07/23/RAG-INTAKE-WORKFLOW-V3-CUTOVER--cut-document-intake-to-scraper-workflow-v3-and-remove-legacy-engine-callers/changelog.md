# Changelog

## 2026-07-23

- Initial workspace created


## 2026-07-23

Step 1: create intake Workflow V3 cutover ticket, guide, tasks, and detailed implementation diary

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/ttmp/2026/07/23/RAG-INTAKE-WORKFLOW-V3-CUTOVER--cut-document-intake-to-scraper-workflow-v3-and-remove-legacy-engine-callers/design-doc/01-rag-intake-workflow-v3-cutover-design-and-implementation-guide.md — Accepted implementation route


## 2026-07-23

Hard-cut RAG intake CLI, API, frontend, provider operations, and persistence to Workflow V3; delete the old engine path and pass focused acceptance. The final dependency-chain task remains open.

### Related Files

- cmd/rag-eval/cmds/intake/worker.go — Built-binary worker and run-once execution boundary
- internal/api/workflow_artifact_handlers.go — Canonical bounded intake API
- pkg/ragintakeworkflow/application.go — Canonical intake Workflow V3 product application
- web/src/components/workflows/WorkflowsView.tsx — Rendered Workflow V3 intake UI

