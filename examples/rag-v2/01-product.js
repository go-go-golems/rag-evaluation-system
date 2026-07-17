const { rag, pipeline, retrievalFor } = require("./common");
const inputs = require("./inputs.json");

const product = rag.product("transcript-assistant", (p) =>
	p
		.pipeline(pipeline)
		.query(retrievalFor(["raw", "summary", "question"], "unit"))
		.rerank(rag.rerank.crossEncoder({ model: "bge-reranker-v2", candidates: 20, results: 5 }))
		.generate(
			rag.generation.answer({
				model: "qwen-answer-v1",
				prompt: "grounded-answer-v2",
				citations: "required",
				contextBudgetTokens: 6000,
			}),
		)
		.request((r) => r.field("query", "string", { required: true, maxLength: 4096 }))
		.response((r) => r.answer("markdown").citations("source").includeTraceId(true))
		.runtime((r) => r.timeoutMs(15000).maxConcurrent(16).onProviderFailure("fail")),
);

module.exports = product.compileProduct(inputs);
