const workflow = require("workflow");
const tasks = require("rag-ttc-v3-tasks");

module.exports = workflow.compile(
	workflow.define("rag-ttc-preparation-v1", (plan) => {
		plan.budget("generation", {
			limits: { requests: 2000, input_tokens: 8000, output_tokens: 4000, cost_microunits: 2000 },
			policyDigest: "sha256:1111111111111111111111111111111111111111111111111111111111111111",
		});
		plan.budget("embedding", {
			limits: { embedding_tokens: 6000 },
			policyDigest: "sha256:2222222222222222222222222222222222222222222222222222222222222222",
		});
		const chunks = plan.inputSet("chunks", {
			itemSchema: "rag-ttc-chunk/v1",
			manifestSchema: "scraper-workflow-item-manifest/v1",
		});
		const generated = plan.map(
			"generate-representations",
			chunks,
			(chunk) => tasks.generate({ chunk }),
			(map) =>
				map
					.pageSize(64)
					.maxItems(2000)
					.maxMaterializedAhead(128)
					.budget({
						account: "generation",
						reserve: { requests: 1, input_tokens: 4, output_tokens: 2, cost_microunits: 1 },
						onExhausted: "fail-run",
					}),
		);
		const embedded = plan.map(
			"embed-representations",
			generated,
			(generatedItem) => tasks.embed({ generated: generatedItem }),
			(map) =>
				map
					.pageSize(64)
					.maxItems(2000)
					.maxMaterializedAhead(128)
					.budget({
						account: "embedding",
						reserve: { embedding_tokens: 3 },
						onExhausted: "fail-run",
					}),
		);
		const shard = plan.reduce(
			"merge-prepared",
			embedded,
			(partition) => tasks.merge({ partition }),
			(reduce) => reduce.fanIn(32).maxLevels(4),
		);
		plan.output("prepared", shard);
	}),
);
