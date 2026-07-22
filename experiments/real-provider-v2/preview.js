const { rag, buildPipeline, retrieval } = require("./base");

const study = rag.study("ttc-real-provider-v2-preview", (s) =>
	s
		.pipeline(buildPipeline("ttc-real-provider-v2-preview", 100000))
		.dataset(
			rag.datasets.artifact("evaluation-dataset", {
				split: "preview",
				status: "candidate",
				relevanceTarget: "unit",
			}),
		)
		.variants((variants) =>
			variants.add("all-real-preview", (entry) =>
				entry.selectRepresentations(["raw", "summary", "question"]).query(() => retrieval()),
			),
		)
		.replicates(1)
		.metrics((metrics) =>
			metrics.latency(["query"]).tokenUsage().providerCost().storageBytes().failureRates(),
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
