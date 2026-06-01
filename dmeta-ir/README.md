# RAG Evaluation DMETA IR

This directory contains the first hand-authored DMETA-style catalog for the RAG Evaluation dashboard. It is documentation-only for now: no React scaffolds are generated from it, and no repository validator is required for the current scope. It starts with the Search Workbench vertical slice and intentionally keeps the layers separate:

- `core-model/` names target-neutral RAG semantics.
- `interactions/` names visible obligations and user actions without choosing React, CSS, or Web layout.
- `meta-design-systems/web/` lowers those obligations into Web dashboard templates and coarse layout choices.
- `instantiations/` selects the first dashboard slice to generate/promote.

Every major YAML object should include `summary`, `description`, `why`, examples where useful, and explicit non-goals. Comments are only navigation aids; durable explanation belongs in fields. Future validation or scaffold generation belongs to later phases, not this documentation pass.
