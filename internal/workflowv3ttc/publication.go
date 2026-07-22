package workflowv3ttc

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

const (
	ValidationReceiptSchema   = "rag-ttc-validation-receipt/v1"
	PublicationDecisionSchema = "rag-ttc-publication-decision/v1"
	PublicationReceiptSchema  = "rag-ttc-publication-receipt/v1"
)

type ValidationReceipt struct {
	SchemaVersion       string `json:"schemaVersion"`
	ShardDigest         string `json:"shardDigest"`
	ItemCount           int    `json:"itemCount"`
	RepresentationCount int    `json:"representationCount"`
	EmbeddingCount      int    `json:"embeddingCount"`
}

type PublicationDecision struct {
	SchemaVersion string `json:"schemaVersion"`
	Approved      bool   `json:"approved"`
	ShardDigest   string `json:"shardDigest"`
	PolicyDigest  string `json:"policyDigest"`
}

type PublicationReceipt struct {
	SchemaVersion  string                           `json:"schemaVersion"`
	ShardDigest    string                           `json:"shardDigest"`
	PreparedDigest string                           `json:"preparedDigest"`
	Identity       ragengine.PreparedCorpusIdentity `json:"identity"`
	ItemCount      int                              `json:"itemCount"`
}

type PublicationConfig struct {
	Store              ragengine.PreparedCorpusStore
	Engine             *ragengine.Engine
	Pipeline           ragcontract.PipelineIR
	Corpus             ragoperators.Corpus
	Options            ragengine.Options
	Identity           ragengine.PreparedCorpusIdentity
	ChunksOutputKey    string
	RawOutputKey       string
	DerivedOutputKey   string
	MergedOutputKey    string
	EmbeddingOutputKey string
	PolicyDigest       string
}

type PublicationService interface {
	Validate(PreparedShard) (ValidationReceipt, error)
	Publish(context.Context, PreparedShard, PublicationDecision) (PublicationReceipt, error)
}

type PublicationAuthority struct{ config PublicationConfig }

func NewPublicationAuthority(config PublicationConfig) (*PublicationAuthority, error) {
	if config.Store == nil || config.Engine == nil || config.Identity.SchemaVersion != "rag-prepared-corpus-manifest/v1" ||
		config.ChunksOutputKey == "" || config.RawOutputKey == "" || config.DerivedOutputKey == "" || config.MergedOutputKey == "" || config.EmbeddingOutputKey == "" || config.PolicyDigest == "" {
		return nil, fmt.Errorf("complete prepared publication configuration is required")
	}
	return &PublicationAuthority{config: config}, nil
}

func (a *PublicationAuthority) Validate(shard PreparedShard) (ValidationReceipt, error) {
	if shard.SchemaVersion != ShardSchema || shard.Digest == "" || len(shard.Items) == 0 || shard.FirstKey != shard.Items[0].Generated.Key || shard.LastKey != shard.Items[len(shard.Items)-1].Generated.Key {
		return ValidationReceipt{}, fmt.Errorf("RAG_TTC_PUBLICATION_SHARD")
	}
	representationCount, embeddingCount := 0, 0
	previous := ""
	for _, item := range shard.Items {
		if item.Generated.Key <= previous || validateEmbedded(item.Generated, item) != nil {
			return ValidationReceipt{}, fmt.Errorf("RAG_TTC_PUBLICATION_ITEM")
		}
		previous = item.Generated.Key
		representationCount += len(item.Representations)
		embeddingCount += len(item.Embeddings)
	}
	return ValidationReceipt{SchemaVersion: ValidationReceiptSchema, ShardDigest: shard.Digest, ItemCount: len(shard.Items), RepresentationCount: representationCount, EmbeddingCount: embeddingCount}, nil
}

func (a *PublicationAuthority) Publish(ctx context.Context, shard PreparedShard, decision PublicationDecision) (PublicationReceipt, error) {
	receipt, err := a.Validate(shard)
	if err != nil {
		return PublicationReceipt{}, err
	}
	if decision.SchemaVersion != PublicationDecisionSchema || !decision.Approved || decision.ShardDigest != shard.Digest || decision.PolicyDigest != a.config.PolicyDigest {
		return PublicationReceipt{}, fmt.Errorf("RAG_TTC_PUBLICATION_DECISION")
	}
	chunks := make([]ragoperators.Chunk, 0, len(shard.Items))
	raw := []ragoperators.Representation{}
	derived := []ragoperators.Representation{}
	merged := []ragoperators.Representation{}
	embeddings := []ragoperators.Embedding{}
	for _, item := range shard.Items {
		chunks = append(chunks, item.Generated.Chunk)
		raw = append(raw, item.RawRepresentations...)
		derived = append(derived, item.Generated.Representations...)
		merged = append(merged, item.Representations...)
		embeddings = append(embeddings, item.Embeddings...)
	}
	sort.Slice(chunks, func(i, j int) bool { return chunks[i].Record.ID < chunks[j].Record.ID })
	sort.Slice(raw, func(i, j int) bool { return raw[i].Record.ID < raw[j].Record.ID })
	sort.Slice(derived, func(i, j int) bool { return derived[i].Record.ID < derived[j].Record.ID })
	sort.Slice(merged, func(i, j int) bool { return merged[i].Record.ID < merged[j].Record.ID })
	sort.Slice(embeddings, func(i, j int) bool {
		return embeddings[i].Record.RepresentationID < embeddings[j].Record.RepresentationID
	})
	values := map[string]any{a.config.ChunksOutputKey: chunks, a.config.RawOutputKey: raw, a.config.DerivedOutputKey: derived, a.config.MergedOutputKey: merged, a.config.EmbeddingOutputKey: embeddings}
	preparedDigest, err := ragengine.PublishPreparedCorpus(ctx, ragengine.PreparedCorpusPublication{Store: a.config.Store, Engine: a.config.Engine, Pipeline: a.config.Pipeline, Corpus: a.config.Corpus, Options: a.config.Options, Identity: a.config.Identity, Values: values})
	if err != nil {
		return PublicationReceipt{}, err
	}
	return PublicationReceipt{SchemaVersion: PublicationReceiptSchema, ShardDigest: receipt.ShardDigest, PreparedDigest: preparedDigest, Identity: a.config.Identity, ItemCount: receipt.ItemCount}, nil
}
