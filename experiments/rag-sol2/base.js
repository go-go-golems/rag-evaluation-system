const rag = require("rag");

const pipeline = rag.pipeline("rag-sol2-parity-fixture", (p) =>
	p
		.corpus(rag.inputs.corpus("corpus"))
		.units(rag.units.identity())
		.chunks(rag.chunks.recursive({ maxRunes: 800, overlapSpans: 0, levels: ["runes"] }))
		.representations(
			rag.representations.compose(
				rag.representations.raw("raw"),
				rag.representations.structuredSummary("summary", {
					generator: rag.generation.structured("fixture-summary-v1", {
						prompt: "fixture-transcript-summary-v1",
						outputSchema: "transcript-rag-summary/v1",
					}),
				}),
				rag.representations.syntheticQuestions("question", {
					from: "summary",
					count: 2,
					model: "fixture-question-v1",
					prompt: "fixture-transcript-questions-v1",
				}),
			),
		)
		.embedding(
			rag.embeddings.model("fixture-hash-32-v1", {
				dimensions: 32,
				distance: "cosine",
				normalize: "l2",
				batchSize: 64,
			}),
		)
		.index(
			"representations",
			rag.indexes.bleveMulti({ lexical: true, vector: { distance: "cosine" } }),
		),
);

function retrievalFor(kinds, collapse) {
	return rag.queryPlan(`parity-${kinds.join("-")}`, (q) =>
		q
			.channels(
				kinds.flatMap((kind) => [
					rag.retrieve.bm25(`${kind}.lexical`, {
						index: "representations",
						representation: kind,
						topK: 30,
					}),
					rag.retrieve.vector(`${kind}.vector`, {
						index: "representations",
						representation: kind,
						topK: 30,
					}),
				]),
			)
			.collapseChannels(
				rag.collapse.parent({ scope: collapse, representative: "scoreThenRepresentationId" }),
			)
			.fuse(rag.fusion.weightedRRF({ rankConstant: 60, weights: { "raw.vector": 2 } }))
			.collapseFinal(
				rag.collapse.parent({
					scope: collapse,
					representative: "bestFusionContributionThenId",
				}),
			)
			.hydrate(rag.hydration.sourceEvidence({ selection: "bestContributionThenId" }))
			.results(10),
	);
}

module.exports = { rag, pipeline, retrievalFor };
