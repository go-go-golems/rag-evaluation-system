const { rag, retrieval } = require("./base");

const batchSize = 2;
const pipeline = rag.pipeline("ttc-flash-combined-speed", (p) =>
	p
		.corpus(rag.inputs.corpus("corpus"))
		.units(rag.units.identity())
		.chunks(rag.chunks.recursive({ maxRunes: 1200, overlapSpans: 0, levels: ["runes"] }))
		.representations(
			rag.representations.compose(
				rag.representations.raw("raw"),
				rag.representations.combinedSummaryQuestions({
					model: "generator-umans-flash",
					prompt: "ttc-combined-preparation-v2",
					outputSchema: "rag-combined-preparation/v2",
					batchSize,
					questionsPerChunk: 4,
					maxBatchRunes: 6000,
				}),
			),
		)
		.embedding(
			rag.embeddings.model("embedding-primary", {
				dimensions: 768,
				distance: "cosine",
				normalize: "l2",
				batchSize: 16,
			}),
		)
		.index(
			"representations",
			rag.indexes.bleveMulti({
				lexical: true,
				vector: { distance: "cosine", optimizeFor: "recall" },
			}),
		),
);

const study = rag.study("ttc-flash-combined-preparation-speed", (s) =>
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
			variants.add("flash-combined-batch-2", (entry) =>
				entry
					.selectRepresentations(["raw", "summary", "question"])
					.query(() => retrieval())
					.rerank(
						rag.rerank.crossEncoder({ model: "reranker-primary", candidates: 20, results: 5 }),
					)
					.generate(
						rag.generation.answer({
							model: "generator-umans",
							prompt: "ttc-grounded-answer-v1",
							citations: "required",
							citationFailurePolicy: "abstain",
							contextBudgetTokens: 6000,
						}),
					),
			),
		)
		.replicates(1)
		.metrics((metrics) =>
			metrics.latency(["query"]).tokenUsage().providerCost().storageBytes().failureRates(),
		)
		.tag("experiment", "preparation-speed-only")
		.tag("benchmarkClaim", "false")
		.tag("fixtureProviders", "false"),
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
