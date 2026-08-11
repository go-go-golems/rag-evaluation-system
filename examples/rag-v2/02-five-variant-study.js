const { rag, pipeline, retrievalFor } = require("./common");
const inputs = require("./inputs.json");
const variants = {
	raw: ["raw"],
	summary: ["summary"],
	rawSummary: ["raw", "summary"],
	rawQuestion: ["raw", "question"],
	all: ["raw", "summary", "question"],
};

const study = rag.study("representation-study", (s) =>
	s
		.pipeline(pipeline)
		.dataset(
			rag.datasets.artifact("judgments", {
				split: "smoke",
				status: "candidate",
				relevanceTarget: "unit",
			}),
		)
		.variants((v) => {
			for (const [name, kinds] of Object.entries(variants)) {
				v.add(name, (x) =>
					x
						.selectRepresentations(kinds)
						.query((ctx) => retrievalFor(kinds, ctx.factor("collapse"))),
				);
			}
		})
		.factors((f) => f.enum("collapse", ["chunk", "unit"]))
		.replicates(3)
		.metrics((m) =>
			m
				.precisionAt([5])
				.recallAt([5])
				.hitRateAt([5])
				.mrr()
				.ndcgAt([5])
				.latency(["query"])
				.tokenUsage()
				.providerCost()
				.storageBytes()
				.failureRates(),
		)
		.invariants((i) =>
			i.require("derived-is-not-source-evidence/v1").require("source-hydrated-final-hit/v1"),
		)
		.tag("evaluationStatus", "candidate"),
);

module.exports = study.compileStudy(inputs);
