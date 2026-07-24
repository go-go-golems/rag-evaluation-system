const rag = require("rag");

const pipeline = rag.pipeline("ttc-scripted-raw", (p) =>
	p
		.corpus(rag.inputs.corpus("corpus"))
		.units(rag.units.identity())
		.chunks(rag.chunks.recursive({ maxRunes: 800, overlapSpans: 120 }))
		.representations(rag.representations.raw("raw"))
		.index("representations", rag.indexes.bleveMulti({ lexical: true })),
);

function query(collapse) {
	return rag.queryPlan(`ttc-query-${collapse}`, (q) =>
		q
			.channels([
				rag.retrieve.bm25("raw.lexical", {
					index: "representations",
					representation: "raw",
					topK: 10,
				}),
			])
			.collapseChannels(
				rag.collapse.parent({ scope: collapse, representative: "scoreThenRepresentationId" }),
			)
			.fuse(rag.fusion.weightedRRF({ rankConstant: 60 }))
			.collapseFinal(
				rag.collapse.parent({ scope: collapse, representative: "bestFusionContributionThenId" }),
			)
			.hydrate(rag.hydration.sourceEvidence({ selection: "bestContributionThenId" }))
			.results(10),
	);
}

module.exports = { rag, pipeline, query };
