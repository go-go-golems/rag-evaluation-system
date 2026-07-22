// Package workflowv3ttc owns the RAG-specific contracts and host authority used
// to run the TTC preparation workload on scraper Workflow V3.
package workflowv3ttc

import (
	"context"
	"fmt"
	"math"
	"sort"
	"strings"

	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

const (
	ModuleAlias        = "rag:ttc"
	ResourceGeneration = "rag.generation.remote"
	ResourceEmbedding  = "rag.embedding.local"
	ResourceLocal      = "rag.local"

	ChunkSchema     = "rag-ttc-chunk/v1"
	GeneratedSchema = "rag-ttc-generated/v1"
	EmbeddedSchema  = "rag-ttc-embedded/v1"
	ShardSchema     = "rag-ttc-prepared-shard/v1"
)

type Chunk struct {
	Key          string   `json:"key"`
	Text         string   `json:"text"`
	TextDigest   string   `json:"textDigest"`
	CitationIDs  []string `json:"citationIds"`
	SourceDigest string   `json:"sourceDigest"`
}

type Generated struct {
	Key                   string   `json:"key"`
	TextDigest            string   `json:"textDigest"`
	Representation        string   `json:"representation"`
	CitationIDs           []string `json:"citationIds"`
	ProviderProfileDigest string   `json:"providerProfileDigest"`
	ModelDigest           string   `json:"modelDigest"`
}

type Embedded struct {
	Generated              Generated `json:"generated"`
	Vector                 []float64 `json:"vector"`
	EmbeddingProfileDigest string    `json:"embeddingProfileDigest"`
}

type PreparedShard struct {
	SchemaVersion string     `json:"schemaVersion"`
	FirstKey      string     `json:"firstKey"`
	LastKey       string     `json:"lastKey"`
	Items         []Embedded `json:"items"`
	Digest        string     `json:"digest"`
}

type Usage struct {
	Dimension string `json:"dimension"`
	Units     int64  `json:"units"`
}

type Result[T any] struct {
	Value T       `json:"value"`
	Usage []Usage `json:"usage"`
}

type Provider interface {
	Generate(context.Context, Chunk) (Result[Generated], error)
	Embed(context.Context, Generated) (Result[Embedded], error)
}

type Failure struct {
	Class     string
	Code      string
	Retryable bool
}

func (f *Failure) Error() string { return f.Class + "/" + f.Code }

func validateChunk(chunk Chunk) error {
	if strings.TrimSpace(chunk.Key) == "" || strings.TrimSpace(chunk.Text) == "" {
		return fmt.Errorf("chunk key and text are required")
	}
	digest, err := workflowv3.Digest(chunk.Text)
	if err != nil || digest != chunk.TextDigest {
		return fmt.Errorf("chunk text digest does not match")
	}
	if len(chunk.CitationIDs) == 0 {
		return fmt.Errorf("chunk citations are required")
	}
	for index := 1; index < len(chunk.CitationIDs); index++ {
		if chunk.CitationIDs[index] <= chunk.CitationIDs[index-1] {
			return fmt.Errorf("chunk citations must be unique and sorted")
		}
	}
	return nil
}

func validateGenerated(chunk Chunk, generated Generated) error {
	if generated.Key != chunk.Key || generated.TextDigest != chunk.TextDigest ||
		generated.ProviderProfileDigest == "" || generated.ModelDigest == "" ||
		strings.TrimSpace(generated.Representation) == "" {
		return fmt.Errorf("generated representation identity is invalid")
	}
	if strings.Join(generated.CitationIDs, "\x00") != strings.Join(chunk.CitationIDs, "\x00") {
		return fmt.Errorf("generated citations do not match source")
	}
	return nil
}

func validateEmbedded(generated Generated, embedded Embedded) error {
	if embedded.Generated.Key != generated.Key || embedded.Generated.TextDigest != generated.TextDigest ||
		embedded.EmbeddingProfileDigest == "" || len(embedded.Vector) == 0 {
		return fmt.Errorf("embedded representation identity is invalid")
	}
	for _, value := range embedded.Vector {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return fmt.Errorf("embedding contains non-finite values")
		}
	}
	return nil
}

func mergeShards(items []Embedded) (PreparedShard, error) {
	if len(items) == 0 {
		return PreparedShard{}, fmt.Errorf("prepared shard is empty")
	}
	sort.Slice(items, func(i, j int) bool { return items[i].Generated.Key < items[j].Generated.Key })
	for index := 1; index < len(items); index++ {
		if items[index].Generated.Key == items[index-1].Generated.Key {
			return PreparedShard{}, fmt.Errorf("duplicate prepared key")
		}
	}
	shard := PreparedShard{SchemaVersion: ShardSchema, FirstKey: items[0].Generated.Key, LastKey: items[len(items)-1].Generated.Key, Items: items}
	digest, err := workflowv3.Digest(struct {
		SchemaVersion string     `json:"schemaVersion"`
		FirstKey      string     `json:"firstKey"`
		LastKey       string     `json:"lastKey"`
		Items         []Embedded `json:"items"`
	}{shard.SchemaVersion, shard.FirstKey, shard.LastKey, shard.Items})
	if err != nil {
		return PreparedShard{}, err
	}
	shard.Digest = digest
	return shard, nil
}
