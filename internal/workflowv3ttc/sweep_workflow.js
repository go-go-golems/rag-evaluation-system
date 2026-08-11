const workflow = require("workflow");
const tasks = require("rag-ttc-v3-tasks");

module.exports = workflow.compile(
	workflow.define("rag-ttc-generation-sweep-cell-v1", (plan) => {
		plan.budget("generation", {
			limits: {
				requests: 100,
				input_tokens: 1000000,
				output_tokens: 1000000,
				cost_microunits: 100000000,
			},
			policyDigest: "sha256:5ec08d94b05f5d6f322f989c36bf3d7428e954856e5b60fdfa7734f52cbad1d6",
		});
		plan.budget("embedding", {
			limits: { embedding_tokens: 6553600, requests: 800 },
			policyDigest: "sha256:9a2800d3a51199c423c42e0faccba04d315f206590df4e87b02b44ea49b2c61a",
		});
		const batches = plan.inputSet("batches", {
			itemSchema: "rag-ttc-chunk-batch/v1",
			manifestSchema: "scraper-workflow-item-manifest/v1",
		});
		const generated = plan.map(
			"generate-batches",
			batches,
			(batch) => tasks.generateBatch({ batch }),
			(map) =>
				map
					.pageSize(1)
					.maxItems(100)
					.maxMaterializedAhead(16)
					.budget({
						account: "generation",
						reserve: {
							requests: 1,
							input_tokens: 16384,
							output_tokens: 8192,
							cost_microunits: 10650,
						},
						onExhausted: "fail-run",
					}),
		);
		const measured = plan.map(
			"measure-embedding-batches",
			generated,
			(generatedBatch) => tasks.embedBatch({ generated: generatedBatch }),
			(map) =>
				map
					.pageSize(1)
					.maxItems(100)
					.maxMaterializedAhead(16)
					.budget({
						account: "embedding",
						reserve: { embedding_tokens: 65536, requests: 8 },
						onExhausted: "fail-run",
					}),
		);
		plan.outputSet("measured", measured);
	}),
);
