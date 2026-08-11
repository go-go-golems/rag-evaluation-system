// Runs a fresh TTC hybrid observation through the generated JavaScript API.
// Prerequisite: the documented rag-ollama-mimimi tunnel serves nomic-embed-text
// at 127.0.0.1:11435. The operational profile is deliberately separate from
// this immutable experiment specification.
const gp = require("geppetto");
const rag = require("rag");

const settings = gp.inferenceProfiles
	.load("ttmp/2026/07/14/RAGEVAL-RAG-DSL-001--typed-fluent-javascript-rag-laboratory-module/scripts/04-mimimi-ollama-embeddings-profile.yaml")
	.resolve("ttc-mimimi-nomic-embed");
const embedder = gp.embeddings(settings);
const lab = rag.open({
	database: "data/rag-eval.db",
	execution: "allowRuns",
	queryEmbed: (query) => embedder.embed(query),
});

try {
	const experiment = rag.experiment("ttc-js-geppetto-weighted-rrf-v1", (e) => e
		.corpus("sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409")
		.chunks("sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392")
		.bm25("sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691")
		.embeddings("sha256:2665c5249b8352ce6904fc00c934534dd179f3eeef0a6a75429a9034be0e03e0")
		.evaluation("candidate:ttc-baseline-v1")
		.note("JavaScript RAG module validation using explicit Geppetto/Ollama query embeddings through the mimimi loopback tunnel.")
		.tag("provider", "ollama/nomic-embed-text/768")
		.tag("runtime", "rag-eval-js")
		.retrieval((r) => r
			.channel("lexical", (c) => c.bm25().topK(50))
			.channel("semantic", (c) => c.vector().topK(50))
			.fuse((f) => f.rrf().rankConstant(60).weight("semantic", 2))
			.collapse("document")
			.results(10))
		.metrics((m) => m
			.relevanceAt(rag.grade("2_SUBSTANTIAL"))
			.recallAt([10])
			.mrr()));

	const report = experiment.validate(lab);
	if (!report.ok) {
		throw new Error(`Refusing incompatible immutable inputs: ${JSON.stringify(report, null, 2)}`);
	}

	const result = lab.execute(experiment);
	console.log(JSON.stringify({ embeddingModel: embedder.model(), result }, null, 2));
} finally {
	lab.close();
}
