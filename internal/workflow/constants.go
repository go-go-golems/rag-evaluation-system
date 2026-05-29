package workflow

const (
	// IntakeRunnerKind is the scraper runner kind used by rag-eval intake workflow ops.
	IntakeRunnerKind = "rag-eval/intake"

	QueueCPU       = "rag-eval:cpu"
	QueueLLM       = "rag-eval:llm"
	QueueEmbedding = "rag-eval:embedding"
	QueueIndex     = "rag-eval:index"
)
