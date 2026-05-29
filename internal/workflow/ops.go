package workflow

const (
	OperationEcho              = "echo"
	OperationChunkDocument     = "chunk_document"
	OperationComputeEmbeddings = "compute_embeddings"
	OperationBuildBM25         = "build_bm25"
)

type IntakeOpInput struct {
	Operation string `json:"operation"`
	DBPath    string `json:"db_path,omitempty"`

	Payload map[string]any `json:"payload,omitempty"`

	DocumentID   string `json:"document_id,omitempty"`
	Strategy     string `json:"strategy,omitempty"`
	StrategyID   string `json:"strategy_id,omitempty"`
	ChunkSize    int    `json:"chunk_size,omitempty"`
	Overlap      int    `json:"overlap,omitempty"`
	StrategyName string `json:"strategy_name,omitempty"`
	Description  string `json:"description,omitempty"`

	SourceIDs   []string `json:"source_ids,omitempty"`
	DocumentIDs []string `json:"document_ids,omitempty"`

	ProfileRegistries []string `json:"profile_registries,omitempty"`
	Profile           string   `json:"profile,omitempty"`
	BaseProfile       string   `json:"base_profile,omitempty"`
	EmbeddingType     string   `json:"embeddings_type,omitempty"`
	EmbeddingEngine   string   `json:"embeddings_engine,omitempty"`
	Dimensions        int      `json:"embeddings_dimensions,omitempty"`
	APIKey            string   `json:"api_key,omitempty"`
	BaseURL           string   `json:"base_url,omitempty"`
	CacheType         string   `json:"cache_type,omitempty"`
	CacheDirectory    string   `json:"cache_directory,omitempty"`
	BatchSize         int      `json:"batch_size,omitempty"`
	Limit             int      `json:"limit,omitempty"`
	Force             bool     `json:"force,omitempty"`

	IndexRoot string `json:"index_root,omitempty"`
	IndexID   string `json:"index_id,omitempty"`
}

type ChunkDocumentOutput struct {
	DocumentID string `json:"document_id"`
	StrategyID string `json:"strategy_id"`
	ChunkCount int    `json:"chunk_count"`
}

type ComputeEmbeddingsOutput struct {
	StrategyID   string   `json:"strategy_id"`
	SourceIDs    []string `json:"source_ids,omitempty"`
	DocumentIDs  []string `json:"document_ids,omitempty"`
	ProviderType string   `json:"provider_type"`
	Model        string   `json:"model"`
	Dimensions   int      `json:"dimensions"`
	Considered   int      `json:"considered"`
	Computed     int      `json:"computed"`
	SkippedFresh int      `json:"skipped_fresh"`
}

type BuildBM25Output struct {
	IndexID       string   `json:"index_id"`
	StrategyID    string   `json:"strategy_id"`
	SourceIDs     []string `json:"source_ids,omitempty"`
	DocumentIDs   []string `json:"document_ids,omitempty"`
	IndexPath     string   `json:"index_path"`
	ChunkCount    int      `json:"chunk_count"`
	DocumentCount int      `json:"document_count"`
	Rebuilt       bool     `json:"rebuilt"`
}
