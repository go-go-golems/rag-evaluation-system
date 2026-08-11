const { rag, pipeline, retrievalFor } = require("./base");
const variants = [
	require("./variants/raw"),
	require("./variants/summary"),
	require("./variants/raw-summary"),
	require("./variants/raw-question"),
	require("./variants/all"),
];

const study = rag.study("rag-sol2-candidate-parity", (s) =>
	s
		.pipeline(pipeline)
		.dataset(
			rag.datasets.artifact("evaluation-dataset", {
				split: "candidate",
				status: "candidate",
				relevanceTarget: "unit",
			}),
		)
		.variants((builder) => {
			for (const variant of variants) {
				builder.add(variant.id, (entry) =>
					entry
						.selectRepresentations(variant.kinds)
						.query((context) => retrievalFor(variant.kinds, context.factor("collapse"))),
				);
			}
		})
		.factors((factors) => factors.enum("collapse", ["chunk", "unit"]))
		.replicates(1)
		.metrics((metrics) =>
			metrics
				.precisionAt([10])
				.recallAt([10])
				.hitRateAt([10])
				.mrr()
				.ndcgAt([10])
				.latency(["query"])
				.tokenUsage()
				.providerCost()
				.storageBytes()
				.failureRates(),
		)
		.invariants((invariants) =>
			invariants
				.require("derived-is-not-source-evidence/v1")
				.require("one-vote-per-collapse-key-per-channel/v1")
				.require("source-hydrated-final-hit/v1"),
		)
		.tag("evaluationStatus", "candidate")
		.tag("parityOracle", "rag-sol2-one-time"),
);

module.exports = study.compileStudy({
	inputs: {
		corpus: {
			role: "corpus",
			kind: "manifest-envelope",
			digest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			schemaVersion: "rag-corpus-snapshot-manifest/v2",
		},
		"evaluation-dataset": {
			role: "evaluation-dataset",
			kind: "manifest-envelope",
			digest: "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
			schemaVersion: "rag-evaluation-dataset-manifest/v2",
		},
	},
});
