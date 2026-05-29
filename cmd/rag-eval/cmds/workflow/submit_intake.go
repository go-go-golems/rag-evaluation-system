package workflow

import (
	"encoding/json"
	"fmt"

	workflowservice "github.com/go-go-golems/rag-evaluation-system/internal/workflow"
	"github.com/spf13/cobra"
)

type submitIntakeOptions struct {
	req workflowservice.SubmitIntakeRequest
}

func newSubmitIntakeCommand() *cobra.Command {
	opts := &submitIntakeOptions{}
	cmd := &cobra.Command{
		Use:   "submit-intake",
		Short: "Submit a durable chunk/embed/BM25 intake workflow",
		Long: `Submit an intake workflow to the scraper engine database.

The submitted workflow wraps the existing rag-eval services: chunking remains
service-backed, embeddings use the normal provider resolver, and BM25 indexing
uses the existing search service. Direct chunk/embedding/search commands remain
available for debugging the same work outside workflow orchestration.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			result, err := workflowservice.SubmitIntakeWorkflow(cmd.Context(), opts.req)
			if err != nil {
				return err
			}
			b, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return err
			}
			fmt.Fprintln(cmd.OutOrStdout(), string(b))
			return nil
		},
	}
	addEngineDBFlag(cmd, &opts.req.EngineDB)
	cmd.Flags().StringVar(&opts.req.DBPath, "db", "data/rag-eval.db", "Path to the rag-eval SQLite database")
	cmd.Flags().StringVar(&opts.req.WorkflowID, "workflow-id", "", "Workflow ID; defaults to intake timestamp")
	cmd.Flags().StringVar(&opts.req.Name, "name", "rag-eval intake workflow", "Workflow display name")
	cmd.Flags().StringSliceVar(&opts.req.DocumentIDs, "document-ids", nil, "Document IDs to chunk, comma-separated or repeated")
	cmd.Flags().StringSliceVar(&opts.req.SourceIDs, "source-ids", nil, "Source IDs used for document selection and downstream filtering")
	cmd.Flags().IntVar(&opts.req.DocumentLimit, "document-limit", 0, "Maximum documents to select when --document-ids is omitted")
	cmd.Flags().StringVar(&opts.req.Strategy, "strategy", "fixed", "Chunking strategy")
	cmd.Flags().IntVar(&opts.req.ChunkSize, "chunk-size", 1200, "Chunk size")
	cmd.Flags().IntVar(&opts.req.Overlap, "overlap", 150, "Chunk overlap")
	cmd.Flags().StringSliceVar(&opts.req.ProfileRegistries, "profile-registries", nil, "Profile registry sources for embedding provider resolution")
	cmd.Flags().StringVar(&opts.req.Profile, "profile", "", "Embedding-capable profile to resolve")
	cmd.Flags().StringVar(&opts.req.BaseProfile, "base-profile", "", "Base profile to overlay direct embedding flags onto")
	cmd.Flags().StringVar(&opts.req.EmbeddingType, "embeddings-type", "ollama", "Embedding provider type: ollama or openai")
	cmd.Flags().StringVar(&opts.req.EmbeddingEngine, "embeddings-engine", "nomic-embed-text", "Embedding model/engine")
	cmd.Flags().IntVar(&opts.req.Dimensions, "embeddings-dimensions", 768, "Embedding dimensions")
	cmd.Flags().StringVar(&opts.req.APIKey, "api-key", "", "Provider API key")
	cmd.Flags().StringVar(&opts.req.BaseURL, "base-url", "", "Provider base URL")
	cmd.Flags().StringVar(&opts.req.CacheType, "cache-type", "none", "Embedding cache type: none, memory, or file")
	cmd.Flags().StringVar(&opts.req.CacheDirectory, "cache-directory", "state/embedding-cache", "Embedding cache directory")
	cmd.Flags().IntVar(&opts.req.BatchSize, "batch-size", 16, "Embedding batch size")
	cmd.Flags().IntVar(&opts.req.EmbeddingLimit, "embedding-limit", 0, "Maximum chunks to consider for embeddings")
	cmd.Flags().BoolVar(&opts.req.ForceEmbeddings, "force-embeddings", false, "Recompute embeddings even when fresh")
	cmd.Flags().BoolVar(&opts.req.SkipEmbeddings, "skip-embeddings", false, "Submit only chunking/BM25 ops without embedding op")
	cmd.Flags().StringVar(&opts.req.IndexRoot, "index-root", "data/indexes", "BM25 index root")
	cmd.Flags().StringVar(&opts.req.IndexID, "index-id", "", "BM25 index ID; defaults to bm25-<workflow-id>")
	cmd.Flags().IntVar(&opts.req.IndexLimit, "index-limit", 0, "Maximum chunks to index")
	cmd.Flags().BoolVar(&opts.req.ForceIndex, "force-index", false, "Replace an existing BM25 index")
	cmd.Flags().BoolVar(&opts.req.SkipBM25, "skip-bm25", false, "Submit only chunking/embedding ops without BM25 op")
	return cmd
}
