module.exports = {
	schemaVersion: "researchctl-analysis-spec/v1",
	id: "ttc-scripted-report",
	title: "TTC scripted acceptance study",
	groupBy: ["collapse"],
	reducers: [
		{ name: "runs", kind: "count" },
		{ name: "failedRuns", kind: "count-failed" },
		{
			name: "mrr",
			kind: "mean-ci",
			metric: "rag.mrr",
			scopePrefix: "rag.query.",
			withinRun: "mean",
			confidence: 0.95,
		},
		{ name: "retries", kind: "sum", metric: "workflow.retries", requiredUnit: "count" },
		{
			name: "failedOperations",
			kind: "sum",
			metric: "workflow.external_operations.failed",
			requiredUnit: "count",
		},
	],
	charts: [{ id: "mrr", title: "MRR by collapse scope", x: "collapse", y: "mrr.mean" }],
	interpretation:
		"This thin TTC acceptance checks that quality and retry-aware execution evidence share immutable run identity. Three replicates support workflow-level completeness checks; this bounded fixture is not a provider-performance recommendation.",
};
