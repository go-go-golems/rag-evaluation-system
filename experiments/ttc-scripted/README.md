# Scripted TTC acceptance

This directory is the workload-only acceptance source for the converged experiment platform. It contains no scheduler, worker, SQLite importer, operation ledger, or analysis program outside the maintained framework layers.

- `project.js` declares the Researchctl project and experiment.
- `study.js` declares two RAG cases and three replicates.
- `pipeline.js` contains reusable RAG composition.
- `analysis.js` declares missing-aware, replicate-aware Researchctl analysis.
- `inputs.json` is a compile-only placeholder; the acceptance script resolves immutable TTC artifacts from the authorized SQLite source.

Run the complete fixture-backed study from a clean checkout:

```bash
ttmp/2026/07/22/TTC-SCRIPTED-EXPERIMENT-ACCEPTANCE--*/scripts/01-run-fixture-ttc-study.sh
```

Set `RAG_TTC_DATABASE` to select another authorized immutable TTC SQLite source. Set `KEEP_TTC_SCRIPTED_WORK=1` to preserve the generated bundle, Workflow V3 stores, Researchctl laboratory, analysis table, SVG, and Markdown report for inspection.

The script builds binaries, derives a bounded three-document/three-query input from TTC, compiles a `rag-workflow-study-bundle/v1`, executes two cases × three replicates through Researchctl and Workflow V3, resumes without duplicate runs, and regenerates byte-identical analysis. Fixture providers make this command cost-free and deterministic. Real provider acceptance remains an explicit authorized predecessor gate under `RAG-GEPPETTO-WORKFLOW-OPERATIONS`; it does not change this study's scientific claims.

This is architecture acceptance, not a provider-performance recommendation. The bounded data and three replicates prove lineage, custody, quality-measure availability, resume, and reproducible analysis—not general model quality.
