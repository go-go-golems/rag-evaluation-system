# Changelog

## 2026-07-18

- Initial workspace created


## 2026-07-18

Step 1: Created native Glazed migration inventory and task plan; direct command lint baseline established.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/main.go — Current adapter-based registration


## 2026-07-18

Step 2: Verified source leaves are already native Glazed; no rewrite required.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/cmds/source/create.go — Typed Glazed source create command
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/cmds/source/list.go — Typed Glazed source list command
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/cmds/source/scan.go — Typed Glazed source scan command


## 2026-07-18

Step 3: Verified corpus leaves are already native Glazed and pass command-directory Glazed lint.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/cmds/corpus/import_ttc.go — Typed Glazed corpus import command
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/cmds/corpus/snapshot_ttc.go — Typed Glazed corpus snapshot command


## 2026-07-18

Step 2: Ported chunk enrichment, document preprocessing, and all workflow leaves to native Glazed (ae9ba32, 153089e, 61362e7); only study and preview remain.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/cmds/workflow/submit_intake.go — Largest native Glazed schema port


## 2026-07-18

Step 3: Ported study and preview, deleted the Glazed-Cobra adapter, and made cmd/rag-eval mandatory Glazed lint scope; make lint/test/logcopter-check pass.

### Related Files

- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/Makefile — Mandatory command-directory Glazed lint
- /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system/cmd/rag-eval/main.go — Native command-group registration

