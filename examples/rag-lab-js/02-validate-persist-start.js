// Persist and submit an immutable experiment. Edit all IDs and database first.
// The script performs catalog lineage validation before it writes anything.
const rag = require("rag");

const database = "data/rag-eval.db";
const lab = rag.open({ database, execution: "allowRuns" });

try {
	const experiment = rag.experiment("ttc-raw-hybrid", (e) => e
		.corpus("REPLACE_WITH_CORPUS_SNAPSHOT_ID")
		.chunks("REPLACE_WITH_CHUNK_SET_ID")
		.bm25("REPLACE_WITH_BM25_ARTIFACT_ID")
		.embeddings("REPLACE_WITH_EMBEDDING_SET_ID")
		.evaluation("REPLACE_WITH_EVALUATION_DATASET_ID")
		.retrieval((r) => r
			.channel("lexical", (c) => c.bm25().topK(50))
			.channel("semantic", (c) => c.vector().topK(50))
			.fuse((f) => f.rrf().rankConstant(60))
			.collapse("document")
			.results(10))
		.metrics((m) => m.relevanceAt(rag.grade("2_SUBSTANTIAL")).recallAt([10]).mrr()));

	const report = experiment.validate(lab);
	if (!report.ok) {
		throw new Error(`Refusing to persist incompatible experiment: ${JSON.stringify(report, null, 2)}`);
	}

	const specification = lab.persist(experiment);
	const run = lab.start(experiment);
	console.log(JSON.stringify({ specification, run }, null, 2));
} finally {
	lab.close();
}
