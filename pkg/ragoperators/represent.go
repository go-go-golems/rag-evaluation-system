package ragoperators

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"sync"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type MemoryCache struct {
	mu     sync.RWMutex
	values map[string][]byte
}

func NewMemoryCache() *MemoryCache { return &MemoryCache{values: map[string][]byte{}} }
func (c *MemoryCache) Get(k string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.values[k]
	return append([]byte(nil), v...), ok
}
func (c *MemoryCache) Put(k string, v []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.values[k] = append([]byte(nil), v...)
}

type representationOperator struct{ kind string }

func (o representationOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: o.kind, Version: "v1"}
}
func (o representationOperator) Execute(ctx context.Context, node ragcontract.Node, inputs map[string]any, env *Environment) (map[string]any, error) {
	chunks, ok := inputs["chunks"].([]Chunk)
	if !ok {
		return nil, fmt.Errorf("RAG_REPRESENTATION_INPUT: chunks")
	}
	var config struct {
		Name, From, Model, Prompt, OutputSchema string
		Count                                   int
		Parameters                              json.RawMessage
		SeedPolicy                              *ragcontract.SeedPolicy
	}
	if err := decodeConfig(node.Config, &config); err != nil {
		return nil, err
	}
	if config.Name == "" {
		config.Name = representationKind(o.kind)
	}
	sourceByChunk := map[string]Representation{}
	if source, ok := inputs["source"].([]Representation); ok {
		for _, representation := range source {
			sourceByChunk[representation.Record.ParentChunkID] = representation
		}
	}
	var modelManifest ragcontract.ModelManifest
	var promptManifest ragcontract.PromptManifest
	if o.kind != "representations.raw" {
		var err error
		modelManifest, err = resolveModel(env, config.Model)
		if err != nil {
			return nil, err
		}
		promptManifest, err = resolvePrompt(env, config.Prompt)
		if err != nil {
			return nil, err
		}
	}
	out := []Representation{}
	for _, chunk := range chunks {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if o.kind == "representations.raw" {
			out = append(out, newRepresentation(config.Name, "raw", chunk, chunk.Text, "source", 0, nil))
			continue
		}
		if env.Generator == nil {
			return nil, fmt.Errorf("RAG_GENERATOR_UNAVAILABLE: %s", o.Ref().ID())
		}
		count := 1
		if o.kind == "representations.synthetic-questions" {
			count = config.Count
			if count <= 0 {
				return nil, fmt.Errorf("RAG_QUESTION_COUNT")
			}
		}
		generationText := chunk.Text
		parentDigest := chunk.Record.TextDigest
		parentRepresentationIDs := []string{}
		if config.From != "" {
			source, exists := sourceByChunk[chunk.Record.ID]
			if !exists {
				return nil, fmt.Errorf("RAG_REPRESENTATION_SOURCE_MISSING: %s/%s", config.From, chunk.Record.ID)
			}
			generationText = source.Text
			parentDigest = source.Record.ContentDigest
			parentRepresentationIDs = []string{source.Record.ID}
		}
		parentID := chunk.Record.ID
		if len(parentRepresentationIDs) > 0 {
			parentID = parentRepresentationIDs[0]
		}
		outputSchema := config.OutputSchema
		if outputSchema == "" {
			outputSchema = promptManifest.OutputSchema
		}
		request := GenerationRequest{Kind: o.kind, Model: modelManifest.ModelID, Prompt: promptManifest.PromptID, OutputSchema: outputSchema, ParentID: parentID, Text: generationText, Count: count}
		key, _ := ragcontract.Digest(struct {
			Ref    ragcontract.OperatorRef
			Config json.RawMessage
			Parent string
		}{o.Ref(), node.Config, parentDigest})
		var result GenerationResult
		cacheOutcome := "miss"
		if env.Cache != nil {
			if raw, found := env.Cache.Get(key); found {
				if json.Unmarshal(raw, &result) == nil {
					cacheOutcome = "hit"
				}
			}
		}
		if cacheOutcome == "miss" {
			value, err := env.Generator.Generate(ctx, request)
			if err != nil {
				return nil, fmt.Errorf("RAG_GENERATION_FAILED: %s: %w", chunk.Record.ID, err)
			}
			result = value
			if env.Cache != nil {
				raw, _ := json.Marshal(result)
				env.Cache.Put(key, raw)
			}
		}
		addGenerationUsage(env, config.Model, result)
		if o.kind == "representations.structured-summary" {
			if result.Text == "" || !json.Valid([]byte(result.Text)) {
				return nil, fmt.Errorf("RAG_STRUCTURED_OUTPUT_INVALID: %s", chunk.Record.ID)
			}
			if config.OutputSchema == "" || promptManifest.OutputSchema != config.OutputSchema {
				return nil, fmt.Errorf("RAG_STRUCTURED_OUTPUT_SCHEMA: configured schema does not match prompt manifest")
			}
			if env.Schemas == nil {
				return nil, fmt.Errorf("RAG_OUTPUT_SCHEMA_VALIDATOR_UNAVAILABLE")
			}
			if err := env.Schemas.Validate(config.OutputSchema, json.RawMessage(result.Text)); err != nil {
				return nil, fmt.Errorf("RAG_STRUCTURED_OUTPUT_SCHEMA: %w", err)
			}
			derivation := derivation(o.Ref(), modelManifest.Digest, promptManifest.Digest, parentDigest, parentRepresentationIDs, sourceRecordIDs(chunk), result.Text, result, cacheOutcome)
			out = append(out, newRepresentation(config.Name, "summary", chunk, result.Text, "derived", 0, derivation))
		} else {
			if len(result.Questions) != count {
				return nil, fmt.Errorf("RAG_QUESTION_COUNT_MISMATCH: got %d want %d", len(result.Questions), count)
			}
			for ordinal, question := range result.Questions {
				if question == "" {
					return nil, fmt.Errorf("RAG_QUESTION_EMPTY")
				}
				derivation := derivation(o.Ref(), modelManifest.Digest, promptManifest.Digest, parentDigest, parentRepresentationIDs, sourceRecordIDs(chunk), question, result, cacheOutcome)
				out = append(out, newRepresentation(config.Name, "question", chunk, question, "derived", ordinal, derivation))
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Record.ID < out[j].Record.ID })
	parents := []ragcontract.ParentDigest{}
	if len(chunks) > 0 {
		parents = parent("chunk-set", chunks[0].ManifestDigest, ragcontract.ChunkSetManifestSchema)
	}
	seenParent := map[string]bool{}
	for _, value := range parents {
		seenParent[value.Digest] = true
	}
	for _, source := range sourceByChunk {
		if source.ManifestDigest != "" && !seenParent[source.ManifestDigest] {
			seenParent[source.ManifestDigest] = true
			parents = append(parents, ragcontract.ParentDigest{Role: "representation-set", Digest: source.ManifestDigest, SchemaVersion: ragcontract.RepresentationManifestSchema})
		}
	}
	sort.Slice(parents, func(i, j int) bool { return parents[i].Digest < parents[j].Digest })
	return materializeRepresentations(node, parents, out)
}
func representationKind(operator string) string {
	switch operator {
	case "representations.structured-summary":
		return "summary"
	case "representations.synthetic-questions":
		return "question"
	default:
		return "raw"
	}
}
func newRepresentation(name, kind string, chunk Chunk, text, role string, ordinal int, d *ragcontract.DerivationRef) Representation {
	digest, _ := ragcontract.Digest(text)
	idDigest, _ := ragcontract.Digest(struct {
		Name, Kind, Parent, Digest string
		Ordinal                    int
	}{name, kind, chunk.Record.ID, digest, ordinal})
	return Representation{Record: ragcontract.RepresentationRecord{ID: "representation:" + idDigest[7:23], Kind: kind, ParentChunkID: chunk.Record.ID, ParentUnitID: chunk.Record.ParentUnitID, ContentDigest: digest, EvidenceRole: role, Derivation: d, Citation: chunk.Record.Citation}, Text: text}
}
func sourceRecordIDs(chunk Chunk) []string {
	seen := map[string]bool{}
	ids := []string{}
	for _, sourceRange := range chunk.Ranges {
		if !seen[sourceRange.SourceID] {
			seen[sourceRange.SourceID] = true
			ids = append(ids, sourceRange.SourceID)
		}
	}
	if len(ids) == 0 && chunk.Record.Citation.SourceID != "" {
		ids = append(ids, chunk.Record.Citation.SourceID)
	}
	return ids
}
func derivation(ref ragcontract.OperatorRef, model, prompt, parent string, parentRepresentationIDs, sourceIDs []string, text string, result GenerationResult, cache string) *ragcontract.DerivationRef {
	output, _ := ragcontract.Digest(text)
	return &ragcontract.DerivationRef{Operator: ref, ModelManifestDigest: model, PromptManifestDigest: prompt, Parameters: json.RawMessage(`{}`), ParentDigest: parent, ParentRepresentationIDs: parentRepresentationIDs, SourceRecordIDs: sourceIDs, OutputDigest: output, InputTokens: result.InputTokens, OutputTokens: result.OutputTokens, CacheOutcome: cache}
}

type mergeOperator struct{}

func (mergeOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: "representations.merge", Version: "v1"}
}
func (mergeOperator) Execute(_ context.Context, node ragcontract.Node, inputs map[string]any, _ *Environment) (map[string]any, error) {
	out := []Representation{}
	for _, value := range inputs {
		items, ok := value.([]Representation)
		if !ok {
			return nil, fmt.Errorf("RAG_REPRESENTATION_MERGE_INPUT")
		}
		out = append(out, items...)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Record.ID < out[j].Record.ID })
	for i := 1; i < len(out); i++ {
		if out[i].Record.ID == out[i-1].Record.ID {
			return nil, fmt.Errorf("RAG_REPRESENTATION_DUPLICATE: %s", out[i].Record.ID)
		}
	}
	parents := []ragcontract.ParentDigest{}
	seen := map[string]bool{}
	for _, record := range out {
		if record.ManifestDigest != "" && !seen[record.ManifestDigest] {
			seen[record.ManifestDigest] = true
			parents = append(parents, ragcontract.ParentDigest{Role: fmt.Sprintf("representation-set.%03d", len(parents)), Digest: record.ManifestDigest, SchemaVersion: ragcontract.RepresentationManifestSchema})
		}
	}
	sort.Slice(parents, func(i, j int) bool { return parents[i].Digest < parents[j].Digest })
	return materializeRepresentations(node, parents, out)
}

func materializeRepresentations(node ragcontract.Node, parents []ragcontract.ParentDigest, records []Representation) (map[string]any, error) {
	data, digest := materializationData(records)
	kinds, roles := []string{}, []string{}
	kindSeen, roleSeen := map[string]bool{}, map[string]bool{}
	for _, record := range records {
		if !kindSeen[record.Record.Kind] {
			kindSeen[record.Record.Kind] = true
			kinds = append(kinds, record.Record.Kind)
		}
		if !roleSeen[record.Record.EvidenceRole] {
			roleSeen[record.Record.EvidenceRole] = true
			roles = append(roles, record.Record.EvidenceRole)
		}
	}
	sort.Strings(kinds)
	sort.Strings(roles)
	manifest := ragcontract.RepresentationSetManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.RepresentationManifestSchema, Digest: digest, Parents: parents, Production: &ragcontract.Production{Operator: node.Operator, Config: node.Config}}, RepresentationCount: int64(len(records)), Kinds: kinds, EvidenceRoles: roles}
	for index := range records {
		records[index].ManifestDigest = digest
	}
	artifact := materializedArtifact("representation-set", "rag-representation-set", node.ID+".json", ragcontract.RepresentationManifestSchema, data, manifest)
	return map[string]any{"representations": records, "artifact": artifact, "manifest": manifest}, nil
}

type embeddingOperator struct{}

func (embeddingOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: "embed.model", Version: "v1"}
}
func (embeddingOperator) Execute(ctx context.Context, node ragcontract.Node, inputs map[string]any, env *Environment) (map[string]any, error) {
	items, ok := inputs["representations"].([]Representation)
	if !ok {
		return nil, fmt.Errorf("RAG_EMBED_INPUT")
	}
	if env.Embedder == nil {
		return nil, fmt.Errorf("RAG_EMBEDDER_UNAVAILABLE")
	}
	var config struct {
		Model      string
		Dimensions int
		Normalize  string
	}
	if err := decodeConfig(node.Config, &config); err != nil {
		return nil, err
	}
	modelManifest, err := resolveModel(env, config.Model)
	if err != nil {
		return nil, err
	}
	texts := make([]string, len(items))
	for i := range items {
		texts[i] = items[i].Text
	}
	vectors, usage, err := env.Embedder.Embed(ctx, modelManifest.ModelID, texts)
	if err != nil {
		return nil, fmt.Errorf("RAG_EMBED_FAILED: %w", err)
	}
	if len(vectors) != len(items) {
		return nil, fmt.Errorf("RAG_EMBED_COUNT: got %d want %d", len(vectors), len(items))
	}
	result := make([]Embedding, len(items))
	for i, vector := range vectors {
		if config.Dimensions > 0 && len(vector) != config.Dimensions {
			return nil, fmt.Errorf("RAG_EMBED_DIMENSIONS: item %d got %d want %d", i, len(vector), config.Dimensions)
		}
		norm := 0.0
		for _, v := range vector {
			if math.IsNaN(v) || math.IsInf(v, 0) {
				return nil, fmt.Errorf("RAG_EMBED_NONFINITE")
			}
			norm += v * v
		}
		if config.Normalize == "l2" {
			if norm == 0 {
				return nil, fmt.Errorf("RAG_EMBED_ZERO_VECTOR")
			}
			norm = math.Sqrt(norm)
			for j := range vector {
				vector[j] /= norm
			}
		}
		digest, _ := ragcontract.Digest(vector)
		result[i] = Embedding{Record: ragcontract.EmbeddingRecord{RepresentationID: items[i].Record.ID, ModelManifestDigest: modelManifest.Digest, Dimensions: len(vector), VectorDigest: digest, Position: int64(i)}, Vector: vector}
	}
	env.Usage.EmbeddingTokens += usage.EmbeddingTokens
	data, digest := materializationData(result)
	representationDigest := ""
	if len(items) > 0 {
		representationDigest = items[0].ManifestDigest
	}
	dimensions := config.Dimensions
	if dimensions == 0 && len(result) > 0 {
		dimensions = len(result[0].Vector)
	}
	manifest := ragcontract.EmbeddingSetManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.EmbeddingManifestSchema, Digest: digest, Parents: parent("representation-set", representationDigest, ragcontract.RepresentationManifestSchema), Production: &ragcontract.Production{Operator: (embeddingOperator{}).Ref(), Config: node.Config}}, VectorCount: int64(len(result)), Dimensions: dimensions, Distance: "cosine", Normalization: config.Normalize, ModelManifestDigest: modelManifest.Digest}
	for index := range result {
		result[index].ManifestDigest = digest
	}
	artifact := materializedArtifact("embedding-set", "rag-embedding-set", node.ID+".json", ragcontract.EmbeddingManifestSchema, data, manifest)
	return map[string]any{"embeddings": result, "artifact": artifact, "manifest": manifest}, nil
}
