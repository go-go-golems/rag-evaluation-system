package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-go-golems/glazed/pkg/cli"
	"github.com/go-go-golems/glazed/pkg/cmds"
	"github.com/go-go-golems/glazed/pkg/cmds/fields"
	"github.com/go-go-golems/glazed/pkg/cmds/schema"
	"github.com/go-go-golems/glazed/pkg/cmds/values"
	workflowservice "github.com/go-go-golems/rag-evaluation-system/internal/workflow"
	"github.com/spf13/cobra"
	"io"
	"reflect"
)

type submitIntakeCommand struct{ *cmds.CommandDescription }

var _ cmds.WriterCommand = (*submitIntakeCommand)(nil)

type submitField struct {
	name, target string
	typ          fields.Type
	def          any
	help         string
}

var submitFields = []submitField{
	{"db", "DBPath", fields.TypeString, "data/rag-eval.db", "Path to the rag-eval SQLite database"},
	{"workflow-id", "WorkflowID", fields.TypeString, "", "Workflow ID; defaults to intake timestamp"},
	{"name", "Name", fields.TypeString, "rag-eval intake workflow", "Workflow display name"},
	{"document-ids", "DocumentIDs", fields.TypeStringList, []string{}, "Document IDs to chunk, comma-separated or repeated"},
	{"source-ids", "SourceIDs", fields.TypeStringList, []string{}, "Source IDs used for document selection and downstream filtering"},
	{"document-limit", "DocumentLimit", fields.TypeInteger, 0, "Maximum documents to select when --document-ids is omitted"},
	{"strategy", "Strategy", fields.TypeString, "fixed", "Chunking strategy"},
	{"chunk-size", "ChunkSize", fields.TypeInteger, 1200, "Chunk size"},
	{"overlap", "Overlap", fields.TypeInteger, 150, "Chunk overlap"},
	{"profile-registries", "ProfileRegistries", fields.TypeStringList, []string{}, "Profile registry sources for embedding provider resolution"},
	{"profile", "Profile", fields.TypeString, "", "Embedding-capable profile to resolve"},
	{"base-profile", "BaseProfile", fields.TypeString, "", "Base profile to overlay direct embedding flags onto"},
	{"embeddings-type", "EmbeddingType", fields.TypeString, "ollama", "Embedding provider type: ollama or openai"},
	{"embeddings-engine", "EmbeddingEngine", fields.TypeString, "nomic-embed-text", "Embedding model/engine"},
	{"embeddings-dimensions", "Dimensions", fields.TypeInteger, 768, "Embedding dimensions"},
	{"api-key", "APIKey", fields.TypeString, "", "Provider API key"},
	{"base-url", "BaseURL", fields.TypeString, "", "Provider base URL"},
	{"cache-type", "CacheType", fields.TypeString, "none", "Embedding cache type: none, memory, or file"},
	{"cache-directory", "CacheDirectory", fields.TypeString, "state/embedding-cache", "Embedding cache directory"},
	{"batch-size", "BatchSize", fields.TypeInteger, 16, "Embedding batch size"},
	{"embedding-limit", "EmbeddingLimit", fields.TypeInteger, 0, "Maximum chunks to consider for embeddings"},
	{"force-embeddings", "ForceEmbeddings", fields.TypeBool, false, "Recompute embeddings even when fresh"},
	{"skip-embeddings", "SkipEmbeddings", fields.TypeBool, false, "Submit only chunking/BM25 ops without embedding op"},
	{"index-root", "IndexRoot", fields.TypeString, "data/indexes", "BM25 index root"},
	{"index-id", "IndexID", fields.TypeString, "", "BM25 index ID; defaults to bm25-<workflow-id>"},
	{"index-limit", "IndexLimit", fields.TypeInteger, 0, "Maximum chunks to index"},
	{"force-index", "ForceIndex", fields.TypeBool, false, "Replace an existing BM25 index"},
	{"skip-bm25", "SkipBM25", fields.TypeBool, false, "Submit only chunking/embedding ops without BM25 op"},
	{"skip-preprocessing", "SkipPreprocessing", fields.TypeBool, true, "Skip document preprocessing artifact ops; set false to include fake preprocessing ops"},
	{"preprocess-artifact-type", "PreprocessArtifactType", fields.TypeString, "clean_text", "Document preprocessing artifact type"},
	{"preprocess-prompt-version", "PreprocessPromptVersion", fields.TypeString, "v1", "Document preprocessing prompt version"},
	{"preprocess-provider", "PreprocessDocumentProvider", fields.TypeString, "fake", "Document preprocessing provider; currently only fake"},
	{"preprocess-model", "PreprocessDocumentModel", fields.TypeString, "fake-document-processor", "Document preprocessing model identity"},
	{"force-preprocessing", "ForcePreprocessing", fields.TypeBool, false, "Recompute document preprocessing artifacts even when fresh"},
	{"skip-chunk-enrichment", "SkipChunkEnrichment", fields.TypeBool, true, "Skip chunk enrichment ops; set false to enrich existing first chunks per selected document"},
	{"chunks-per-document-to-enrich", "ChunksPerDocumentToEnrich", fields.TypeInteger, 1, "Maximum existing chunks per selected document to enrich"},
	{"chunk-enrichment-prompt", "ChunkEnrichmentPrompt", fields.TypeString, "v1", "Chunk enrichment prompt version"},
	{"chunk-enrichment-provider", "ChunkEnrichmentProvider", fields.TypeString, "fake", "Chunk enrichment provider; currently only fake"},
	{"chunk-enrichment-model", "ChunkEnrichmentModel", fields.TypeString, "fake-chunk-enricher", "Chunk enrichment model identity"},
	{"force-chunk-enrichment", "ForceChunkEnrichment", fields.TypeBool, false, "Recompute chunk enrichments even when fresh"}}

func newSubmitIntakeCommand() *cobra.Command {
	c, e := newSubmitIntakeGlazeCommand()
	cobra.CheckErr(e)
	r, e := cli.BuildCobraCommandFromCommand(c, cli.WithParserConfig(cli.CobraParserConfig{AppName: "rag-eval", ShortHelpSections: []string{schema.DefaultSlug}}))
	cobra.CheckErr(e)
	return r
}
func newSubmitIntakeGlazeCommand() (*submitIntakeCommand, error) {
	f := make([]*fields.Definition, 0, len(submitFields)+1)
	f = append(f, fields.New("engine-db", fields.TypeString, fields.WithDefault("state/rag-eval-workflows.db"), fields.WithHelp("Path to the scraper workflow engine SQLite database")))
	for _, x := range submitFields {
		f = append(f, fields.New(x.name, x.typ, fields.WithDefault(x.def), fields.WithHelp(x.help)))
	}
	return &submitIntakeCommand{CommandDescription: cmds.NewCommandDescription("submit-intake", cmds.WithShort("Submit a durable chunk/embed/BM25 intake workflow"), cmds.WithFlags(f...))}, nil
}
func (c *submitIntakeCommand) RunIntoWriter(ctx context.Context, v *values.Values, w io.Writer) error {
	data := v.GetDataMap()
	req := workflowservice.SubmitIntakeRequest{}
	rv := reflect.ValueOf(&req).Elem()
	for _, x := range submitFields {
		value, ok := data[x.name]
		if !ok {
			continue
		}
		field := rv.FieldByName(x.target)
		if !field.IsValid() {
			return fmt.Errorf("unknown submit request field %s", x.target)
		}
		field.Set(reflect.ValueOf(value))
	}
	if value, ok := data["engine-db"].(string); ok {
		req.EngineDB = value
	}
	result, e := workflowservice.SubmitIntakeWorkflow(ctx, req)
	if e != nil {
		return e
	}
	return json.NewEncoder(w).Encode(result)
}
