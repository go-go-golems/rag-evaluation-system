// Build and inspect a canonical RAG experiment without opening a database.
// Replace the example IDs before using this plan with a real laboratory.
const rag = require("rag");

const ttcInputs = rag.fragment("ttc-inputs", (e) => e
	.corpus("REPLACE_WITH_CORPUS_SNAPSHOT_ID")
	.chunks("REPLACE_WITH_CHUNK_SET_ID")
	.bm25("REPLACE_WITH_BM25_ARTIFACT_ID")
	.embeddings("REPLACE_WITH_EMBEDDING_SET_ID")
	.evaluation("REPLACE_WITH_EVALUATION_DATASET_ID"));

const experiment = rag.experiment("ttc-raw-hybrid", (e) => e
	.use(ttcInputs)
	.note("Raw immutable chunks, BM25 plus vector retrieval, and weighted RRF.")
	.tag("corpus", "ttc")
	.tag("representation", "raw")
	.representations((r) => r.rawChunks("raw"))
	.retrieval((r) => r
		.channel("lexical", (c) => c.bm25().representation("raw").topK(50))
		.channel("semantic", (c) => c.vector().representation("raw").topK(50))
		.fuse((f) => f.rrf().rankConstant(60).weight("semantic", 1.25))
		.collapse("document")
		.results(10))
	.metrics((m) => m
		.relevanceAt(rag.grade("2_SUBSTANTIAL"))
		.precisionAt([1, 3, 10])
		.recallAt([10])
		.ndcgAt(10)
		.mrr()));

const report = experiment.validate();
if (!report.ok) {
	throw new Error(JSON.stringify(report, null, 2));
}

console.log(JSON.stringify(experiment.toSpec(), null, 2));
