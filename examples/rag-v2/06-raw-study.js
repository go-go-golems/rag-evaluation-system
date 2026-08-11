const rag = require("rag");

const pipeline = rag.pipeline("raw-bm25", (p) =>
	p
		.corpus(rag.inputs.corpus("corpus"))
		.units(rag.units.identity())
		.chunks(rag.chunks.recursive({ maxRunes: 800, overlapSpans: 120 }))
		.representations(rag.representations.raw("raw"))
		.index("representations", rag.indexes.bleveMulti({ lexical: true })),
);
const query = rag.queryPlan("raw-query", (q) =>
	q
		.channels([
			rag.retrieve.bm25("raw.lexical", {
				index: "representations",
				representation: "raw",
				topK: 10,
			}),
		])
		.collapseChannels(
			rag.collapse.parent({ scope: "unit", representative: "scoreThenRepresentationId" }),
		)
		.fuse(rag.fusion.weightedRRF({ rankConstant: 60 }))
		.collapseFinal(
			rag.collapse.parent({
				scope: "unit",
				representative: "bestFusionContributionThenId",
			}),
		)
		.hydrate(rag.hydration.sourceEvidence({ selection: "bestContributionThenId" }))
		.results(10),
);
const study = rag.study("raw-bm25-study", (s) =>
	s
		.pipeline(pipeline)
		.dataset(
			rag.datasets.artifact("evaluation-dataset", {
				split: "smoke",
				status: "candidate",
				relevanceTarget: "unit",
			}),
		)
		.variants((variants) =>
			variants.add("raw", (variant) => variant.selectRepresentations(["raw"]).query(() => query)),
		)
		.replicates(1)
		.metrics((metrics) => metrics.mrr().recallAt([10]).latency(["query"]).storageBytes())
		.invariants((invariants) => invariants.require("source-hydrated-final-hit/v1")),
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
