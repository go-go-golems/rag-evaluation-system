# TTC RAG laboratory baseline and immutable experiment runs

This ticket specifies the first bounded, reproducible TTC RAG laboratory slice in `rag-evaluation-system`. It combines a complete raw-text retrieval baseline with immutable experiment specifications, runs, artifacts, query traces, and terminal results.

Start with [index.md](index.md), then read the intern-oriented [design and implementation guide](design-doc/01-ttc-rag-laboratory-baseline-and-immutable-experiment-runs-design-and-implementation-guide.md). The guide contains the target architecture, schema, APIs, algorithms, pseudocode, diagrams, implementation phases, test strategy, and acceptance criteria. [tasks.md](tasks.md) is the ordered implementation backlog, while [reference/01-implementation-diary.md](reference/01-implementation-diary.md) records how this design package and the rebuilt TTC source database were produced.

Generated databases remain under the ignored repository `data/` directory. Any ticket-specific experiments or utilities added during implementation must be stored under `scripts/` in this ticket.

From the workspace root, use the explicit documentation root:

```bash
docmgr --root rag-evaluation-system/ttmp status --summary-only
docmgr --root rag-evaluation-system/ttmp doctor --ticket RAGEVAL-TTC-LAB-001 --stale-after 30
```
