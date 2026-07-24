const { rag, pipeline, query } = require("./pipeline");
const inputs = require("./inputs.json");

const study = rag.study("ttc-scripted-acceptance", (s) =>
	s
		.pipeline(pipeline)
		.dataset(
			rag.datasets.artifact("evaluation-dataset", {
				split: "acceptance",
				status: "candidate",
				relevanceTarget: "unit",
			}),
		)
		.variants((variants) =>
			variants.add("raw", (variant) =>
				variant.selectRepresentations(["raw"]).query((ctx) => query(ctx.factor("collapse"))),
			),
		)
		.factors((factors) => factors.enum("collapse", ["chunk", "unit"]))
		.replicates(3)
		.metrics((metrics) =>
			metrics
				.precisionAt([5])
				.recallAt([5, 10])
				.hitRateAt([5])
				.mrr()
				.ndcgAt([10])
				.latency(["query"])
				.storageBytes()
				.failureRates(),
		)
		.invariants((invariants) => invariants.require("source-hydrated-final-hit/v1"))
		.tag("evaluationStatus", "candidate"),
);

module.exports = study.compileStudy(inputs);
