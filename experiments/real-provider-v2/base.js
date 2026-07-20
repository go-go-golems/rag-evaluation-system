const rag = require("rag");

function buildPipeline(name, maxRunes) {
	return rag.pipeline(name, (p) =>
		p
			.corpus(rag.inputs.corpus("corpus"))
			.units(rag.units.identity())
			.chunks(rag.chunks.recursive({ maxRunes, overlapSpans: 0, levels: ["runes"] }))
			.representations(
				rag.representations.compose(
					rag.representations.raw("raw"),
					rag.representations.combinedSummaryQuestions({
						model: "generator-umans-flash",
						prompt: "ttc-combined-preparation-v2",
						outputSchema: "rag-combined-preparation/v2",
						batchSize: 2,
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
}

const pipeline = buildPipeline("ttc-real-provider-v2", 1200);

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

module.exports = { rag, buildPipeline, pipeline, retrieval };
