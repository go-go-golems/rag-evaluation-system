const rag = require("rag");
console.log("before-open");
const lab = rag.open({
	database: "data/rag-eval.db",
	execution: "allowRuns",
	reranker: {
		kind: "llama.cpp",
		baseURL: "http://127.0.0.1:18013",
		model: "dengcao/Qwen3-Reranker-8B:q4_k_m",
	},
});
console.log("after-open");
lab.close();
console.log("after-close");
