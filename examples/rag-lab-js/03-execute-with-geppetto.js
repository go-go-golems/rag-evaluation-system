// Execute an immutable vector + lexical RRF experiment through Geppetto.
// Replace every REPLACE_WITH_... value before running this script.
const gp = require("geppetto");
const rag = require("rag");

const database = "data/rag-eval.db";
const profileRegistry = "REPLACE_WITH_PROFILE_REGISTRY.yaml";
const profileName = "REPLACE_WITH_EMBEDDING_PROFILE_NAME";

// The profile registry keeps credentials, base URL, and model configuration
// outside the immutable experiment. embed() synchronously returns number[].
const settings = gp.inferenceProfiles.load(profileRegistry).resolve(profileName);
const embedder = gp.embeddings(settings);
const lab = rag.open({
	database,
	execution: "allowRuns",
	queryEmbed: (query) => embedder.embed(query),
});

try {
	const experiment = rag.experiment("ttc-raw-vector-rrf", (e) =>
		e
			.corpus("REPLACE_WITH_CORPUS_SNAPSHOT_ID")
			.chunks("REPLACE_WITH_CHUNK_SET_ID")
			.bm25("REPLACE_WITH_BM25_ARTIFACT_ID")
			.embeddings("REPLACE_WITH_EMBEDDING_SET_ID")
			.evaluation("REPLACE_WITH_EVALUATION_DATASET_ID")
			.retrieval((r) =>
				r
					.channel("lexical", (c) => c.bm25().topK(50))
					.channel("semantic", (c) => c.vector().topK(50))
					.fuse((f) => f.rrf().rankConstant(60))
					.collapse("document")
					.results(10),
			)
			.metrics((m) => m.relevanceAt(rag.grade("2_SUBSTANTIAL")).recallAt([10]).mrr()),
	);

	const report = experiment.validate(lab);
	if (!report.ok) {
		throw new Error(`Refusing incompatible immutable inputs: ${JSON.stringify(report, null, 2)}`);
	}

	const result = lab.execute(experiment);
	console.log(JSON.stringify({ embeddingModel: embedder.model(), result }, null, 2));
} finally {
	lab.close();
}
