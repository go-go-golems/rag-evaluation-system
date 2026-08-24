package intake

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragintakeworkflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/spf13/cobra"
)

type submitIntakeCommand struct{ *cmds.CommandDescription }

var _ cmds.WriterCommand = (*submitIntakeCommand)(nil)

type submitSettings struct {
	DB                      string   `glazed:"db"`
	WorkflowDB              string   `glazed:"workflow-db"`
	ArtifactRoot            string   `glazed:"artifact-root"`
	IndexRoot               string   `glazed:"index-root"`
	RunID                   string   `glazed:"run-id"`
	DocumentIDs             []string `glazed:"document-ids"`
	SourceIDs               []string `glazed:"source-ids"`
	DocumentLimit           int      `glazed:"document-limit"`
	Strategy                string   `glazed:"strategy"`
	ChunkSize               int      `glazed:"chunk-size"`
	Overlap                 int      `glazed:"overlap"`
	ProfileRegistries       []string `glazed:"profile-registries"`
	Profile                 string   `glazed:"profile"`
	BaseProfile             string   `glazed:"base-profile"`
	EmbeddingType           string   `glazed:"embeddings-type"`
	EmbeddingEngine         string   `glazed:"embeddings-engine"`
	Dimensions              int      `glazed:"embeddings-dimensions"`
	CacheType               string   `glazed:"cache-type"`
	ProviderAuthorityDigest string   `glazed:"provider-authority-digest"`
	BatchSize               int      `glazed:"batch-size"`
	EmbeddingLimit          int      `glazed:"embedding-limit"`
	ForceEmbeddings         bool     `glazed:"force-embeddings"`
	SkipEmbeddings          bool     `glazed:"skip-embeddings"`
	IndexID                 string   `glazed:"index-id"`
	IndexLimit              int      `glazed:"index-limit"`
	ForceIndex              bool     `glazed:"force-index"`
	SkipBM25                bool     `glazed:"skip-bm25"`
	SkipPreprocessing       bool     `glazed:"skip-preprocessing"`
	PreprocessArtifactType  string   `glazed:"preprocess-artifact-type"`
	PreprocessPromptVersion string   `glazed:"preprocess-prompt-version"`
	PreprocessProvider      string   `glazed:"preprocess-provider"`
	PreprocessModel         string   `glazed:"preprocess-model"`
	ForcePreprocessing      bool     `glazed:"force-preprocessing"`
	SkipChunkEnrichment     bool     `glazed:"skip-chunk-enrichment"`
	ChunksPerDocument       int      `glazed:"chunks-per-document-to-enrich"`
	ChunkEnrichmentPrompt   string   `glazed:"chunk-enrichment-prompt"`
	ChunkEnrichmentProvider string   `glazed:"chunk-enrichment-provider"`
	ChunkEnrichmentModel    string   `glazed:"chunk-enrichment-model"`
	ForceChunkEnrichment    bool     `glazed:"force-chunk-enrichment"`
}

func newSubmitIntakeCommand() *cobra.Command {
	command, err := newSubmitIntakeGlazeCommand()
	cobra.CheckErr(err)
	result, err := cli.BuildCobraCommandFromCommand(command, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
	cobra.CheckErr(err)
	return result
}
func newSubmitIntakeGlazeCommand() (*submitIntakeCommand, error) {
	f := []*fields.Definition{
		fields.New("db", fields.TypeString, fields.WithDefault("data/rag-eval.db"), fields.WithHelp("RAG domain SQLite database")), fields.New("workflow-db", fields.TypeString, fields.WithDefault("state/rag-eval-intake-v3.db"), fields.WithHelp("Workflow V3 SQLite database")), fields.New("artifact-root", fields.TypeString, fields.WithDefault("state/rag-eval-intake-v3-artifacts"), fields.WithHelp("Workflow V3 artifact root")), fields.New("index-root", fields.TypeString, fields.WithDefault("data/indexes"), fields.WithHelp("Host BM25 index root")), fields.New("run-id", fields.TypeString, fields.WithHelp("Immutable Workflow V3 run ID")),
		fields.New("document-ids", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("Document IDs")), fields.New("source-ids", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("Source IDs used for selection")), fields.New("document-limit", fields.TypeInteger, fields.WithDefault(0), fields.WithHelp("Selection limit")),
		fields.New("strategy", fields.TypeString, fields.WithDefault("fixed"), fields.WithHelp("Chunk strategy")), fields.New("chunk-size", fields.TypeInteger, fields.WithDefault(1200), fields.WithHelp("Chunk size")), fields.New("overlap", fields.TypeInteger, fields.WithDefault(150), fields.WithHelp("Chunk overlap")),
		fields.New("profile-registries", fields.TypeStringList, fields.WithDefault([]string{}), fields.WithHelp("Provider profile registries")), fields.New("profile", fields.TypeString, fields.WithHelp("Embedding profile")), fields.New("base-profile", fields.TypeString, fields.WithHelp("Embedding base profile")), fields.New("embeddings-type", fields.TypeString, fields.WithDefault("ollama"), fields.WithHelp("Embedding provider type")), fields.New("embeddings-engine", fields.TypeString, fields.WithDefault("nomic-embed-text"), fields.WithHelp("Embedding model")), fields.New("embeddings-dimensions", fields.TypeInteger, fields.WithDefault(768), fields.WithHelp("Embedding dimensions")), fields.New("cache-type", fields.TypeString, fields.WithDefault("none"), fields.WithHelp("Cache identity")), fields.New("provider-authority-digest", fields.TypeString, fields.WithHelp("Exact host provider authority digest")), fields.New("batch-size", fields.TypeInteger, fields.WithDefault(16), fields.WithHelp("Embedding batch size")), fields.New("embedding-limit", fields.TypeInteger, fields.WithDefault(0), fields.WithHelp("Embedding limit")), fields.New("force-embeddings", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Recompute embeddings")), fields.New("skip-embeddings", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Skip embeddings")),
		fields.New("index-id", fields.TypeString, fields.WithHelp("BM25 index ID")), fields.New("index-limit", fields.TypeInteger, fields.WithDefault(0), fields.WithHelp("Index limit")), fields.New("force-index", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Replace index")), fields.New("skip-bm25", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Skip BM25")),
		fields.New("skip-preprocessing", fields.TypeBool, fields.WithDefault(true), fields.WithHelp("Skip preprocessing")), fields.New("preprocess-artifact-type", fields.TypeString, fields.WithDefault("clean_text"), fields.WithHelp("Preprocess artifact")), fields.New("preprocess-prompt-version", fields.TypeString, fields.WithDefault("v1"), fields.WithHelp("Preprocess prompt")), fields.New("preprocess-provider", fields.TypeString, fields.WithDefault("fake"), fields.WithHelp("Preprocess provider identity")), fields.New("preprocess-model", fields.TypeString, fields.WithDefault("fake-document-processor"), fields.WithHelp("Preprocess model")), fields.New("force-preprocessing", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Recompute preprocessing")),
		fields.New("skip-chunk-enrichment", fields.TypeBool, fields.WithDefault(true), fields.WithHelp("Skip enrichment")), fields.New("chunks-per-document-to-enrich", fields.TypeInteger, fields.WithDefault(1), fields.WithHelp("Existing chunks per document")), fields.New("chunk-enrichment-prompt", fields.TypeString, fields.WithDefault("v1"), fields.WithHelp("Enrichment prompt")), fields.New("chunk-enrichment-provider", fields.TypeString, fields.WithDefault("fake"), fields.WithHelp("Enrichment provider")), fields.New("chunk-enrichment-model", fields.TypeString, fields.WithDefault("fake-chunk-enricher"), fields.WithHelp("Enrichment model")), fields.New("force-chunk-enrichment", fields.TypeBool, fields.WithDefault(false), fields.WithHelp("Recompute enrichment")),
	}
	return &submitIntakeCommand{CommandDescription: cmds.NewCommandDescription("submit", cmds.WithShort("Submit a Workflow V3 document intake run"), cmds.WithFlags(f...))}, nil
}
func (c *submitIntakeCommand) RunIntoWriter(ctx context.Context, v *values.Values, w io.Writer) error {
	s := &submitSettings{}
	if err := v.DecodeSectionInto(schema.DefaultSlug, s); err != nil {
		return err
	}
	if s.RunID == "" {
		s.RunID = "intake-" + time.Now().UTC().Format("20060102T150405.000000000")
	}
	if s.IndexID == "" {
		s.IndexID = "bm25-" + s.RunID
	}
	config := ragintakeworkflow.DefaultConfig(s.DB)
	config.WorkflowDatabase = s.WorkflowDB
	config.ArtifactRoot = s.ArtifactRoot
	config.Runtime.IndexRoot = s.IndexRoot
	config.Runtime.ProviderAuthorityDigest = s.ProviderAuthorityDigest
	request := ragintakeworkflow.Request{Strategy: s.Strategy, ChunkSize: s.ChunkSize, Overlap: s.Overlap, SkipPreprocessing: s.SkipPreprocessing, ForcePreprocessing: s.ForcePreprocessing, PreprocessArtifactType: s.PreprocessArtifactType, PreprocessPromptVersion: s.PreprocessPromptVersion, PreprocessProvider: s.PreprocessProvider, PreprocessModel: s.PreprocessModel, SkipChunkEnrichment: s.SkipChunkEnrichment, ForceChunkEnrichment: s.ForceChunkEnrichment, ChunkEnrichmentPrompt: s.ChunkEnrichmentPrompt, ChunkEnrichmentProvider: s.ChunkEnrichmentProvider, ChunkEnrichmentModel: s.ChunkEnrichmentModel, SkipEmbeddings: s.SkipEmbeddings, ForceEmbeddings: s.ForceEmbeddings, ProfileRegistries: s.ProfileRegistries, Profile: s.Profile, BaseProfile: s.BaseProfile, EmbeddingType: s.EmbeddingType, EmbeddingEngine: s.EmbeddingEngine, Dimensions: s.Dimensions, CacheType: s.CacheType, BatchSize: s.BatchSize, EmbeddingLimit: s.EmbeddingLimit, SkipBM25: s.SkipBM25, ForceBM25: s.ForceIndex, IndexID: s.IndexID, IndexLimit: s.IndexLimit}
	request, err := ragintakeworkflow.PrepareRequest(ctx, config.Runtime, request, ragintakeworkflow.Selection{DocumentIDs: s.DocumentIDs, SourceIDs: s.SourceIDs, DocumentLimit: s.DocumentLimit, ChunksPerDocumentToEnrich: s.ChunksPerDocument})
	if err != nil {
		return err
	}
	app, err := ragintakeworkflow.Open(ctx, config)
	if err != nil {
		return err
	}
	defer func() { _ = app.Close() }()
	submission, err := app.SubmitRequest(ctx, request, workflowv3.RunID(s.RunID))
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(map[string]any{"submission": submission, "documentIds": request.DocumentIDs, "strategyId": fmt.Sprintf("%s-%d-%d", s.Strategy, s.ChunkSize, s.Overlap)})
}
