const { rag, pipeline, retrieval } = require("./base");

const product = rag.product("ttc-real-provider-assistant", (p) =>
	p
		.pipeline(pipeline)
		.query(retrieval())
		.rerank(rag.rerank.crossEncoder({ model: "reranker-primary", candidates: 20, results: 5 }))
		.generate(
			rag.generation.answer({
				model: "generator-primary",
				prompt: "ttc-grounded-answer-v1",
				citations: "required",
				contextBudgetTokens: 6000,
			}),
		)
		.request((r) => r.field("query", "string", { required: true, maxLength: 4096 }))
		.response((r) => r.answer("markdown").citations("source").includeTraceId(true))
		.runtime((r) => r.timeoutMs(60000).maxConcurrent(1).onProviderFailure("fail")),
);

module.exports = product.compileProduct({
	inputs: {
		corpus: {
			role: "corpus",
			kind: "manifest-envelope",
			digest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			schemaVersion: "rag-corpus-snapshot-manifest/v2",
		},
	},
});
