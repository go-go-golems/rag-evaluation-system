package workflowv3ttc

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/internal/preparationworkflow"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

type MaterializedInputs struct {
	Chunks           workflowv3.ArtifactRef
	Queries          workflowv3.ArtifactRef
	ChunkCount       int
	QueryCount       int
	ValidCitationIDs map[string]struct{}
}

func MaterializeInputs(ctx context.Context, artifacts workflowv3.ArtifactStore, engine *ragengine.Engine, execution ragcontract.PipelineExecution, corpus ragoperators.Corpus, dataset ragoperators.EvaluationDataset, options ragengine.Options, corpusDigest, datasetDigest string) (MaterializedInputs, error) {
	if artifacts == nil || engine == nil || !validDigest(corpusDigest) || !validDigest(datasetDigest) || len(dataset.Queries) == 0 {
		return MaterializedInputs{}, fmt.Errorf("complete TTC materialization inputs are required")
	}
	mapping, err := preparationworkflow.DeriveCanonicalMapping(execution.Pipeline)
	if err != nil {
		return MaterializedInputs{}, err
	}
	inputs, err := engine.StaticInputs(ctx, execution.Pipeline, corpus, options, mapping.CombinedNode.ID)
	if err != nil {
		return MaterializedInputs{}, err
	}
	chunks, ok := inputs["chunks"].([]ragoperators.Chunk)
	if !ok || len(chunks) == 0 {
		return MaterializedInputs{}, fmt.Errorf("RAG_TTC_CHUNKS_EMPTY")
	}
	sort.Slice(chunks, func(i, j int) bool { return chunks[i].Record.ID < chunks[j].Record.ID })
	chunkItems := make([]workflowv3.ManifestItem, len(chunks))
	validCitations := make(map[string]struct{}, len(chunks))
	for index, chunk := range chunks {
		if chunk.Record.ID == "" || (index > 0 && chunk.Record.ID == chunks[index-1].Record.ID) {
			return MaterializedInputs{}, fmt.Errorf("RAG_TTC_CHUNK_IDENTITY")
		}
		envelope := Chunk{Key: chunk.Record.ID, Chunk: chunk, CitationIDs: []string{chunk.Record.ID}, SourceDigest: corpusDigest}
		body, err := workflowv3.CanonicalJSON(envelope)
		if err != nil {
			return MaterializedInputs{}, err
		}
		ref, err := artifacts.Put(ctx, ChunkSchema, "application/json", body)
		if err != nil {
			return MaterializedInputs{}, err
		}
		chunkItems[index] = workflowv3.ManifestItem{Key: chunk.Record.ID, Value: ref}
		validCitations[chunk.Record.ID] = struct{}{}
	}
	chunkManifestRef, err := putManifest(ctx, artifacts, ChunkSchema, chunkItems)
	if err != nil {
		return MaterializedInputs{}, err
	}
	queries := append([]ragoperators.Query(nil), dataset.Queries...)
	sort.Slice(queries, func(i, j int) bool { return queries[i].ID < queries[j].ID })
	queryItems := make([]workflowv3.ManifestItem, len(queries))
	for index, query := range queries {
		if query.ID == "" || query.Text == "" || (index > 0 && query.ID == queries[index-1].ID) {
			return MaterializedInputs{}, fmt.Errorf("RAG_TTC_QUERY_IDENTITY")
		}
		envelope := QueryEnvelope{SchemaVersion: QuerySchema, DatasetDigest: datasetDigest, Query: query}
		body, err := workflowv3.CanonicalJSON(envelope)
		if err != nil {
			return MaterializedInputs{}, err
		}
		ref, err := artifacts.Put(ctx, QuerySchema, "application/json", body)
		if err != nil {
			return MaterializedInputs{}, err
		}
		queryItems[index] = workflowv3.ManifestItem{Key: query.ID, Value: ref}
	}
	queryManifestRef, err := putManifest(ctx, artifacts, QuerySchema, queryItems)
	if err != nil {
		return MaterializedInputs{}, err
	}
	return MaterializedInputs{Chunks: chunkManifestRef, Queries: queryManifestRef, ChunkCount: len(chunks), QueryCount: len(queries), ValidCitationIDs: validCitations}, nil
}

func putManifest(ctx context.Context, artifacts workflowv3.ArtifactStore, schema string, items []workflowv3.ManifestItem) (workflowv3.ArtifactRef, error) {
	manifest, err := workflowv3.NewItemManifest(schema, items)
	if err != nil {
		return workflowv3.ArtifactRef{}, err
	}
	body, err := workflowv3.EncodeItemManifest(manifest)
	if err != nil {
		return workflowv3.ArtifactRef{}, err
	}
	return artifacts.Put(ctx, workflowv3.ItemManifestSchemaV1, "application/json", body)
}
