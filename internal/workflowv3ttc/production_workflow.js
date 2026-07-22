const workflow = require("workflow");
const tasks = require("rag-ttc-v3-tasks");

module.exports = workflow.compile(
	workflow.define("rag-ttc-production-v1", (plan) => {
		plan.budget("generation", {
			limits: { requests: 1, input_tokens: 2048, output_tokens: 2048, cost_microunits: 20000 },
			policyDigest: "sha256:1111111111111111111111111111111111111111111111111111111111111111",
		});
		plan.budget("embedding", {
			limits: { embedding_tokens: 4096 },
			policyDigest: "sha256:2222222222222222222222222222222222222222222222222222222222222222",
		});
		plan.gate("approve-generation-spend", {
			schema: "rag-ttc-budget-decision/v1",
			timeoutMs: 86400000,
			requiredRole: "rag.ttc.budget-approver",
			onReject: "fail-run",
			onExpire: "fail-run",
		});
		plan.gate("approve-embedding-spend", {
			schema: "rag-ttc-budget-decision/v1",
			timeoutMs: 86400000,
			requiredRole: "rag.ttc.budget-approver",
			onReject: "fail-run",
			onExpire: "fail-run",
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
						reserve: {
							requests: 1,
							input_tokens: 2048,
							output_tokens: 2048,
							cost_microunits: 20000,
						},
						onExhausted: "require-approval",
						approvalGate: "approve-generation-spend",
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
						reserve: { embedding_tokens: 4096 },
						onExhausted: "require-approval",
						approvalGate: "approve-embedding-spend",
					}),
		);
		const shard = plan.reduce(
			"merge-prepared",
			embedded,
			(partition) => tasks.merge({ partition }),
			(reduce) => reduce.fanIn(32).maxLevels(4),
		);
		const validation = plan.task("validate-publication", tasks.validatePublication({ shard }));
		const decision = plan.gate(
			"approve-publication",
			{
				schema: "rag-ttc-publication-decision/v1",
				timeoutMs: 86400000,
				requiredRole: "rag.ttc.publisher",
				onReject: "fail-run",
				onExpire: "fail-run",
			},
			(gate) => gate.after(validation),
		);
		const publication = plan.task("publish-prepared", tasks.publish({ shard, decision }));
		plan.output("publication", publication.output("publication"));
	}),
);
