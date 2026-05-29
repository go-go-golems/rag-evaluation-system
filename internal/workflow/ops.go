package workflow

const (
	OperationEcho          = "echo"
	OperationChunkDocument = "chunk_document"
)

type IntakeOpInput struct {
	Operation string `json:"operation"`
	DBPath    string `json:"db_path,omitempty"`

	Payload map[string]any `json:"payload,omitempty"`

	DocumentID   string `json:"document_id,omitempty"`
	Strategy     string `json:"strategy,omitempty"`
	ChunkSize    int    `json:"chunk_size,omitempty"`
	Overlap      int    `json:"overlap,omitempty"`
	StrategyName string `json:"strategy_name,omitempty"`
	Description  string `json:"description,omitempty"`
}

type ChunkDocumentOutput struct {
	DocumentID string `json:"document_id"`
	StrategyID string `json:"strategy_id"`
	ChunkCount int    `json:"chunk_count"`
}
