const rag = require("rag");

const preparation = rag.fragment("transcript-preparation", (p) =>
	p
		.units(rag.transcript.units.agentsViewRuns())
		.chunks(rag.chunks.recursive({ maxRunes: 1200, overlapSpans: 0 })),
);

const representations = rag.representations.compose(
	rag.representations.raw("raw"),
	rag.representations.structuredSummary("summary", {
		generator: rag.generation.structured("summary-qwen-v1", {
			prompt: "transcript-summary-v1",
			outputSchema: "transcript-rag-summary/v1",
		}),
	}),
	rag.representations.syntheticQuestions("question", { from: "summary", count: 4 }),
);

const pipeline = rag.pipeline("transcript-rag", (p) =>
	p
		.corpus(rag.inputs.corpus("corpus"))
		.use(preparation)
		.representations(representations)
		.embedding(rag.embeddings.model("nomic-embed-v1", { distance: "cosine", normalize: "l2" }))
		.index(
			"representations",
			rag.indexes.bleveMulti({
				lexical: true,
				vector: { distance: "cosine", optimizeFor: "recall" },
			}),
		),
);

function retrievalFor(kinds, collapse) {
	return rag.queryPlan(`retrieve-${kinds.join("-")}`, (q) =>
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
				rag.collapse.parent({ scope: collapse, representative: "bestFusionContributionThenId" }),
			)
			.hydrate(rag.hydration.sourceEvidence({ selection: "bestContributionThenId" }))
			.results(5),
	);
}

module.exports = { rag, pipeline, retrievalFor };
