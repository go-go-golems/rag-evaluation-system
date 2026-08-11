package ragoperators

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"sort"

	"github.com/blevesearch/bleve/v2"
	blevequery "github.com/blevesearch/bleve/v2/search/query"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type MultiIndex struct {
	Bleve                bleve.Index
	Representations      map[string]Representation
	Embeddings           map[string][]float64
	Manifest             ragcontract.IndexManifest
	Artifact             []byte
	EmbeddingModelDigest string
}

func (i *MultiIndex) Close() error {
	if i == nil || i.Bleve == nil {
		return nil
	}
	return i.Bleve.Close()
}

type indexOperator struct{ kind string }

func (o indexOperator) Ref() ragcontract.OperatorRef {
	kind := o.kind
	if kind == "" {
		kind = "index.bleve-multi"
	}
	return ragcontract.OperatorRef{Kind: kind, Version: "v1"}
}
func (o indexOperator) Execute(ctx context.Context, node ragcontract.Node, inputs map[string]any, _ *Environment) (map[string]any, error) {
	representations := []Representation{}
	embeddings := []Embedding{}
	for port, value := range inputs {
		if port == "embeddings" {
			var ok bool
			embeddings, ok = value.([]Embedding)
			if !ok {
				return nil, fmt.Errorf("RAG_INDEX_EMBEDDINGS")
			}
			continue
		}
		items, ok := value.([]Representation)
		if !ok {
			return nil, fmt.Errorf("RAG_INDEX_REPRESENTATIONS: %s", port)
		}
		representations = append(representations, items...)
	}
	sort.Slice(representations, func(i, j int) bool { return representations[i].Record.ID < representations[j].Record.ID })
	idx, err := bleve.NewMemOnly(bleve.NewIndexMapping())
	if err != nil {
		return nil, err
	}
	batch := idx.NewBatch()
	byID := map[string]Representation{}
	kinds := map[string]bool{}
	for _, rep := range representations {
		if err := ctx.Err(); err != nil {
			_ = idx.Close()
			return nil, err
		}
		if _, exists := byID[rep.Record.ID]; exists {
			_ = idx.Close()
			return nil, fmt.Errorf("RAG_INDEX_DUPLICATE: %s", rep.Record.ID)
		}
		byID[rep.Record.ID] = rep
		kinds[rep.Record.Kind] = true
		if err := batch.Index(rep.Record.ID, map[string]any{"text": rep.Text, "kind": rep.Record.Kind, "sourceId": rep.Record.Citation.SourceID}); err != nil {
			_ = idx.Close()
			return nil, err
		}
	}
	if err := idx.Batch(batch); err != nil {
		_ = idx.Close()
		return nil, err
	}
	vectors := map[string][]float64{}
	dimensions := 0
	embeddingModelDigest := ""
	for _, embedding := range embeddings {
		if _, ok := byID[embedding.Record.RepresentationID]; !ok {
			_ = idx.Close()
			return nil, fmt.Errorf("RAG_INDEX_EMBEDDING_PARENT: %s", embedding.Record.RepresentationID)
		}
		if dimensions == 0 {
			dimensions = len(embedding.Vector)
		}
		if len(embedding.Vector) != dimensions {
			_ = idx.Close()
			return nil, fmt.Errorf("RAG_INDEX_VECTOR_DIMENSIONS")
		}
		if embeddingModelDigest == "" {
			embeddingModelDigest = embedding.Record.ModelManifestDigest
		} else if embeddingModelDigest != embedding.Record.ModelManifestDigest {
			_ = idx.Close()
			return nil, fmt.Errorf("RAG_INDEX_EMBEDDING_MODEL")
		}
		vectors[embedding.Record.RepresentationID] = embedding.Vector
	}
	kindList := []string{}
	for kind := range kinds {
		kindList = append(kindList, kind)
	}
	sort.Strings(kindList)
	artifact, _ := json.Marshal(struct {
		SchemaVersion   string           `json:"schemaVersion"`
		Representations []Representation `json:"representations"`
	}{"rag-index-records/v1", representations})
	sum := sha256.Sum256(artifact)
	digest := "sha256:" + hex.EncodeToString(sum[:])
	parents := []ragcontract.ParentDigest{}
	seenParent := map[string]bool{}
	for _, representation := range representations {
		if representation.ManifestDigest != "" && !seenParent[representation.ManifestDigest] {
			seenParent[representation.ManifestDigest] = true
			parents = append(parents, ragcontract.ParentDigest{Role: fmt.Sprintf("representation-set.%03d", len(parents)), Digest: representation.ManifestDigest, SchemaVersion: ragcontract.RepresentationManifestSchema})
		}
	}
	for _, embedding := range embeddings {
		if embedding.ManifestDigest != "" && !seenParent[embedding.ManifestDigest] {
			seenParent[embedding.ManifestDigest] = true
			parents = append(parents, ragcontract.ParentDigest{Role: fmt.Sprintf("embedding-set.%03d", len(parents)), Digest: embedding.ManifestDigest, SchemaVersion: ragcontract.EmbeddingManifestSchema})
		}
	}
	sort.Slice(parents, func(i, j int) bool { return parents[i].Digest < parents[j].Digest })
	manifest := ragcontract.IndexManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.IndexManifestSchema, Digest: digest, Parents: parents, Production: &ragcontract.Production{Operator: o.Ref(), Config: node.Config}}, Engine: "bleve", EngineVersion: "v2.6.0", RepresentationKinds: kindList, VectorDimensions: dimensions, Distance: "cosine", DocumentCount: int64(len(representations)), ArtifactTreeDigest: digest}
	return map[string]any{"index": &MultiIndex{Bleve: idx, Representations: byID, Embeddings: vectors, Manifest: manifest, Artifact: artifact, EmbeddingModelDigest: embeddingModelDigest}}, nil
}

type retrieveOperator struct{ kind string }

func (o retrieveOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: o.kind, Version: "v1"}
}
func (o retrieveOperator) Execute(ctx context.Context, node ragcontract.Node, inputs map[string]any, env *Environment) (map[string]any, error) {
	index, ok := inputs["index"].(*MultiIndex)
	if !ok {
		return nil, fmt.Errorf("RAG_RETRIEVE_INDEX")
	}
	query, ok := inputs["query"].(Query)
	if !ok {
		return nil, fmt.Errorf("RAG_RETRIEVE_QUERY")
	}
	var config struct {
		Representation string          `json:"representation"`
		TopK           int             `json:"topK"`
		Filter         RetrievalFilter `json:"filter"`
	}
	if err := decodeConfig(node.Config, &config); err != nil {
		return nil, err
	}
	if len(config.Filter.DocumentIDs) > 0 || len(config.Filter.ContentTypes) > 0 || len(config.Filter.MetadataEquals) > 0 {
		return nil, fmt.Errorf("RAG_FILTER_UNSUPPORTED")
	}
	var hits []RankedRecord
	if o.kind == "retrieve.bm25" {
		match := bleve.NewMatchQuery(query.Text)
		kind := bleve.NewTermQuery(config.Representation)
		kind.SetField("kind")
		parts := []blevequery.Query{match, kind}
		if len(config.Filter.SourceIDs) > 0 {
			disjunction := []blevequery.Query{}
			for _, source := range config.Filter.SourceIDs {
				term := bleve.NewTermQuery(source)
				term.SetField("sourceId")
				disjunction = append(disjunction, term)
			}
			parts = append(parts, bleve.NewDisjunctionQuery(disjunction...))
		}
		request := bleve.NewSearchRequestOptions(bleve.NewConjunctionQuery(parts...), config.TopK, 0, false)
		result, err := index.Bleve.SearchInContext(ctx, request)
		if err != nil {
			return nil, err
		}
		for _, hit := range result.Hits {
			rep, exists := index.Representations[hit.ID]
			if exists {
				hits = append(hits, RankedRecord{Representation: rep, Score: hit.Score, Channel: node.ID})
			}
		}
	} else {
		if env.Embedder == nil {
			return nil, fmt.Errorf("RAG_EMBEDDER_UNAVAILABLE: query")
		}
		modelManifest, err := resolveModel(env, index.EmbeddingModelDigest)
		if err != nil {
			return nil, err
		}
		vectors, usage, err := env.Embedder.Embed(ctx, modelManifest.ModelID, []string{query.Text})
		if err != nil {
			return nil, err
		}
		if len(vectors) != 1 {
			return nil, fmt.Errorf("RAG_QUERY_VECTOR_COUNT")
		}
		env.Usage.EmbeddingTokens += usage.EmbeddingTokens
		for id, vector := range index.Embeddings {
			rep := index.Representations[id]
			if rep.Record.Kind != config.Representation || !filterSource(rep, config.Filter.SourceIDs) {
				continue
			}
			score, err := cosine(vectors[0], vector)
			if err != nil {
				return nil, err
			}
			hits = append(hits, RankedRecord{Representation: rep, Score: score, Channel: node.ID})
		}
		sort.Slice(hits, func(i, j int) bool {
			if hits[i].Score == hits[j].Score {
				return hits[i].Representation.Record.ID < hits[j].Representation.Record.ID
			}
			return hits[i].Score > hits[j].Score
		})
		if len(hits) > config.TopK {
			hits = hits[:config.TopK]
		}
	}
	sort.SliceStable(hits, func(i, j int) bool {
		if hits[i].Score == hits[j].Score {
			return hits[i].Representation.Record.ID < hits[j].Representation.Record.ID
		}
		return hits[i].Score > hits[j].Score
	})
	for i := range hits {
		hits[i].Rank = i + 1
	}
	if env.Trace != nil {
		trace := ragcontract.ChannelTrace{Name: node.ID, Operator: o.Ref(), Hits: []ragcontract.ChannelHit{}}
		for _, hit := range hits {
			trace.Hits = append(trace.Hits, ragcontract.ChannelHit{Rank: hit.Rank, Representation: ragcontract.RepresentationIdentity{ID: hit.Representation.Record.ID, Kind: hit.Representation.Record.Kind, ParentChunkID: hit.Representation.Record.ParentChunkID, ParentUnitID: hit.Representation.Record.ParentUnitID, ContentDigest: hit.Representation.Record.ContentDigest, EvidenceRole: hit.Representation.Record.EvidenceRole}, RawScore: hit.Score, Filter: node.Config})
		}
		env.Trace.Channels = append(env.Trace.Channels, trace)
	}
	return map[string]any{"hits": hits}, nil
}
func filterSource(rep Representation, ids []string) bool {
	if len(ids) == 0 {
		return true
	}
	for _, id := range ids {
		if rep.Record.Citation.SourceID == id {
			return true
		}
	}
	return false
}
func cosine(a, b []float64) (float64, error) {
	if len(a) != len(b) || len(a) == 0 {
		return 0, fmt.Errorf("RAG_VECTOR_DIMENSIONS")
	}
	dot, aa, bb := 0.0, 0.0, 0.0
	for i := range a {
		dot += a[i] * b[i]
		aa += a[i] * a[i]
		bb += b[i] * b[i]
	}
	if aa == 0 || bb == 0 {
		return 0, fmt.Errorf("RAG_VECTOR_ZERO")
	}
	return dot / (math.Sqrt(aa) * math.Sqrt(bb)), nil
}
