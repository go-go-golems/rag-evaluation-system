const rag = require("rag");

const pipeline = rag.pipeline("ttc-real-provider-v2", (p) =>
	p
		.corpus(rag.inputs.corpus("corpus"))
		.units(rag.transcript.units.agentsViewRuns())
		.chunks(rag.chunks.recursive({ maxRunes: 1200, overlapSpans: 0, levels: ["runes"] }))
		.representations(
			rag.representations.compose(
				rag.representations.raw("raw"),
				rag.representations.structuredSummary("summary", {
					generator: rag.generation.structured("generator-primary", {
						prompt: "ttc-summary-v1",
						outputSchema: "transcript-rag-summary/v1",
					}),
				}),
				rag.representations.syntheticQuestions("question", {
					from: "summary",
					count: 4,
					model: "generator-primary",
					prompt: "ttc-questions-v1",
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

function retrieval() {
	return rag.queryPlan("ttc-real-retrieval", (q) =>
		q
			.channels(
				["raw", "summary", "question"].flatMap((kind) => [
					rag.retrieve.bm25(`${kind}.lexical`, {
						index: "representations",
						representation: kind,
						topK: 20,
					}),
					rag.retrieve.vector(`${kind}.vector`, {
						index: "representations",
						representation: kind,
						topK: 20,
					}),
				]),
			)
			.collapseChannels(
				rag.collapse.parent({ scope: "unit", representative: "scoreThenRepresentationId" }),
			)
			.fuse(rag.fusion.weightedRRF({ rankConstant: 60, weights: { "raw.vector": 2 } }))
			.collapseFinal(
				rag.collapse.parent({ scope: "unit", representative: "bestFusionContributionThenId" }),
			)
			.hydrate(rag.hydration.sourceEvidence({ selection: "bestContributionThenId" }))
			.results(5),
	);
}

module.exports = { rag, pipeline, retrieval };
