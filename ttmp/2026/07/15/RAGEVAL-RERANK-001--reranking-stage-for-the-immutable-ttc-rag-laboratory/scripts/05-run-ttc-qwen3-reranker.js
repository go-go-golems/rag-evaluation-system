// Runs the frozen TTC weighted-RRF baseline with the private llama.cpp Qwen3
// reranker. Prerequisite: tunnels at 127.0.0.1:11435 and 127.0.0.1:18013.
const gp = require("geppetto");
const rag = require("rag");

const settings = gp.inferenceProfiles
	.load("ttmp/2026/07/14/RAGEVAL-RAG-DSL-001--typed-fluent-javascript-rag-laboratory-module/scripts/04-mimimi-ollama-embeddings-profile.yaml")
	.resolve("ttc-mimimi-nomic-embed");
const embedder = gp.embeddings(settings);
const rerankerModel = "dengcao/Qwen3-Reranker-8B:q4_k_m";
const lab = rag.open({
	database: "data/rag-eval.db",
	execution: "allowRuns",
	queryEmbed: (query) => embedder.embed(query),
	reranker: { kind: "llama.cpp", baseURL: "http://127.0.0.1:18013", model: rerankerModel },
});

try {
	console.log("qwen: construct experiment");
	const experiment = rag.experiment("ttc-js-geppetto-weighted-rrf-qwen3-rerank-v1", (e) => e
		.corpus("sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409")
		.chunks("sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392")
		.bm25("sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691")
		.embeddings("sha256:2665c5249b8352ce6904fc00c934534dd179f3eeef0a6a75429a9034be0e03e0")
		.evaluation("candidate:ttc-baseline-v1")
		.note("TTC weighted RRF with explicit private llama.cpp Qwen3 reranking; endpoint remains a runtime capability.")
		.tag("provider", "ollama/nomic-embed-text/768")
		.tag("reranker", rerankerModel)
		.retrieval((r) => r
			.channel("lexical", (c) => c.bm25().topK(50))
			.channel("semantic", (c) => c.vector().topK(50))
			.fuse((f) => f.rrf().rankConstant(60).weight("semantic", 2))
			.rerank((x) => x.crossEncoder(rerankerModel).candidates(50).results(10))
			.collapse("document")
			.results(10))
		.metrics((m) => m
			.relevanceAt(rag.grade("2_SUBSTANTIAL"))
			.recallAt([10])
			.mrr()));

	console.log("qwen: validate experiment");
	const report = experiment.validate(lab);
	if (!report.ok) {
		throw new Error(`Refusing incompatible immutable inputs: ${JSON.stringify(report, null, 2)}`);
	}

	console.log("qwen: execute experiment");
	const result = lab.execute(experiment);
	console.log("qwen: execute returned");
	console.log(JSON.stringify({ embeddingModel: embedder.model(), rerankerModel, result }, null, 2));
} finally {
	lab.close();
}
