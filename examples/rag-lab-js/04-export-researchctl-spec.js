// Pure authoring example: this does not open SQLite, contact a model service,
// create a run, or write artifacts. The returned object is the
// rag-retrieval-spec/v1 domain payload consumed by the researchctl adapter.
const rag = require("rag");

const experiment = rag.experiment("ttc-raw-hybrid", (e) =>
	e
		.corpus("REPLACE_WITH_CORPUS_SNAPSHOT_ID")
		.chunks("REPLACE_WITH_CHUNK_SET_ID")
		.bm25("REPLACE_WITH_BM25_ARTIFACT_ID")
		.embeddings("REPLACE_WITH_EMBEDDING_SET_ID")
		.evaluation("REPLACE_WITH_EVALUATION_DATASET_ID")
		.tag("corpus", "ttc")
		.tag("status", "candidate")
		.representations((r) => r.rawChunks("raw"))
		.retrieval((r) =>
			r
				.channel("lexical", (c) => c.bm25().representation("raw").topK(50))
				.channel("semantic", (c) => c.vector().representation("raw").topK(50))
				.fuse((f) => f.rrf().rankConstant(60).weight("semantic", 1.25))
				.collapse("document")
				.results(10),
		)
		.metrics((m) =>
			m
				.relevanceAt(rag.grade("2_SUBSTANTIAL"))
				.precisionAt([1, 3, 10])
				.recallAt([10])
				.ndcgAt(10)
				.mrr(),
		),
);

const specification = experiment.exportSpecification({ datasetSplit: "development" });
console.log(JSON.stringify(specification, null, 2));
