const { rag, pipeline, retrieval } = require("./base");

const study = rag.study("ttc-real-provider-v2-candidate", (s) =>
	s
		.pipeline(pipeline)
		.dataset(
			rag.datasets.artifact("evaluation-dataset", {
				split: "candidate",
				status: "candidate",
				relevanceTarget: "unit",
			}),
		)
		.variants((variants) =>
			variants.add("all-real", (entry) =>
				entry.selectRepresentations(["raw", "summary", "question"]).query(() => retrieval()),
			),
		)
		.replicates(1)
		.metrics((metrics) =>
			metrics
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
		.invariants((invariants) =>
			invariants
				.require("derived-is-not-source-evidence/v1")
				.require("one-vote-per-collapse-key-per-channel/v1")
				.require("source-hydrated-final-hit/v1"),
		)
		.tag("evaluationStatus", "candidate")
		.tag("fixtureProviders", "false")
		.tag("benchmarkClaim", "false"),
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
