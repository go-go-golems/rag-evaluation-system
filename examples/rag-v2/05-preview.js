const { rag, pipeline, retrievalFor } = require("./common");
const inputs = require("./inputs.json");
const study = rag.study("preview", (s) =>
	s
		.pipeline(pipeline)
		.dataset(
			rag.datasets.artifact("judgments", {
				split: "smoke",
				status: "candidate",
				relevanceTarget: "unit",
			}),
		)
		.variants((v) =>
			v.add("all", (x) =>
				x
					.selectRepresentations(["raw", "summary", "question"])
					.query((ctx) => retrievalFor(["raw", "summary", "question"], ctx.factor("collapse"))),
			),
		)
		.factors((f) => f.enum("collapse", ["unit"]))
		.metrics((m) => m.mrr()),
);
module.exports = rag.preview(study, {
	...inputs,
	variant: "all",
	factors: { collapse: "unit" },
	query: "Where was the decision recorded?",
	trace: "full",
});
