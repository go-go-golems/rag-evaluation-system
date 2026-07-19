package ragoperators

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"sync"
	"sync/atomic"
	"time"

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

type representationConfig struct {
	Name, From, Model, Prompt, OutputSchema string
	Count                                   int
	Parameters                              json.RawMessage
	SeedPolicy                              *ragcontract.SeedPolicy
}

type GenerationCacheIdentityV2 struct {
	SchemaVersion                string
	Operator                     ragcontract.OperatorRef
	CanonicalOperatorConfig      json.RawMessage
	ParentDigest                 string
	ModelManifestDigest          string
	PromptManifestDigest         string
	OutputSchemaFingerprint      string
	EffectiveSettingsFingerprint string
}

type generationCacheEnvelopeV2 struct {
	SchemaVersion string                    `json:"schemaVersion"`
	Identity      GenerationCacheIdentityV2 `json:"identity"`
	Value         GenerationResult          `json:"value"`
}

type representationWork struct {
	chunk         Chunk
	parentDigest  string
	parentReprIDs []string
	request       GenerationRequest
	cacheIdentity GenerationCacheIdentityV2
	cacheKey      string
}

type representationResult struct {
	reps  []Representation
	usage GenerationResult
	err   error
}

const defaultRepresentationGenerationWorkers = 1

type schemaRawResolver interface {
	Raw(string) (json.RawMessage, error)
}

func outputSchemaFingerprint(env *Environment, name string) string {
	if raw, ok := env.Schemas.(schemaRawResolver); ok {
		if schema, err := raw.Raw(name); err == nil {
			if digest, err := ragcontract.Digest(json.RawMessage(schema)); err == nil {
				return digest
			}
		}
	}
	// Fixture validators and custom hosts may not expose schema bytes. Their
	// schema name is still part of the identity, but production file registries
	// provide the stronger content digest above.
	digest, _ := ragcontract.Digest(struct{ Name string }{name})
	return digest
}

func cacheIdentityEqual(left, right GenerationCacheIdentityV2) bool {
	leftJSON, leftErr := ragcontract.CanonicalJSON(left)
	rightJSON, rightErr := ragcontract.CanonicalJSON(right)
	return leftErr == nil && rightErr == nil && string(leftJSON) == string(rightJSON)
}

type representationProgress struct {
	Phase       string `json:"phase"`
	Completed   int64  `json:"completed"`
	Total       int64  `json:"total"`
	CacheHits   int64  `json:"cacheHits"`
	CacheMisses int64  `json:"cacheMisses"`
	Active      int64  `json:"active"`
	Failed      int64  `json:"failed"`
}

func emitRepresentationProgress(ctx context.Context, env *Environment, progress representationProgress) error {
	if env == nil || env.EmitEvent == nil {
		return nil
	}
	payload, err := json.Marshal(progress)
	if err != nil {
		return err
	}
	return env.EmitEvent(ctx, Event{Type: "rag.phase.progress/v1", Payload: payload})
}

func (o representationOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: o.kind, Version: "v1"}
}
func (o representationOperator) Execute(ctx context.Context, node ragcontract.Node, inputs map[string]any, env *Environment) (map[string]any, error) {
	chunks, ok := inputs["chunks"].([]Chunk)
	if !ok {
		return nil, fmt.Errorf("RAG_REPRESENTATION_INPUT: chunks")
	}
	var config representationConfig
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
	// Raw representations need no generation; process sequentially.
	if o.kind == "representations.raw" {
		for _, chunk := range chunks {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			out = append(out, newRepresentation(config.Name, "raw", chunk, chunk.Text, "source", 0, nil))
		}
		sort.Slice(out, func(i, j int) bool { return out[i].Record.ID < out[j].Record.ID })
		parents := []ragcontract.ParentDigest{}
		if len(chunks) > 0 {
			parents = parent("chunk-set", chunks[0].ManifestDigest, ragcontract.ChunkSetManifestSchema)
		}
		return materializeRepresentations(node, parents, out)
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
	outputSchema := config.OutputSchema
	if outputSchema == "" {
		outputSchema = promptManifest.OutputSchema
	}
	// Build per-chunk work items. Provider calls are scheduled through a
	// bounded worker pool below; work construction itself is deterministic.
	items := make([]representationWork, 0, len(chunks))
	for _, chunk := range chunks {
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
		request := GenerationRequest{Kind: o.kind, Model: modelManifest.ModelID, Prompt: promptManifest.PromptID, OutputSchema: outputSchema, ParentID: parentID, Text: generationText, Count: count}
		identity := GenerationCacheIdentityV2{
			SchemaVersion: "rag-generation-cache-identity/v2",
			Operator:      o.Ref(), CanonicalOperatorConfig: node.Config,
			ParentDigest: parentDigest, ModelManifestDigest: modelManifest.Digest,
			PromptManifestDigest:         promptManifest.Digest,
			OutputSchemaFingerprint:      outputSchemaFingerprint(env, outputSchema),
			EffectiveSettingsFingerprint: env.GenerationSettingsFingerprint,
		}
		key, _ := ragcontract.Digest(identity)
		items = append(items, representationWork{
			chunk:         chunk,
			parentDigest:  parentDigest,
			parentReprIDs: parentRepresentationIDs,
			request:       request,
			cacheIdentity: identity,
			cacheKey:      key,
		})
	}
	// Resolve cache hits before starting provider work. This keeps hits cheap and
	// leaves only genuine misses for the bounded provider worker pool.
	results := make([]representationResult, len(items))
	missIndexes := make([]int, 0, len(items))
	var completed, cacheHits, failures atomic.Int64
	for i, item := range items {
		var cached GenerationResult
		if env.Cache != nil {
			var envelope generationCacheEnvelopeV2
			if raw, found := env.Cache.Get(item.cacheKey); found && json.Unmarshal(raw, &envelope) == nil && envelope.SchemaVersion == "rag-generation-cache-envelope/v2" && cacheIdentityEqual(envelope.Identity, item.cacheIdentity) {
				cached = envelope.Value
				reps, err := buildRepresentations(o, config, modelManifest, promptManifest, item, cached, "hit", env)
				if err != nil {
					return nil, err
				}
				results[i] = representationResult{reps: reps, usage: cached}
				completed.Add(1)
				cacheHits.Add(1)
				continue
			}
		}
		missIndexes = append(missIndexes, i)
	}

	workCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	var active atomic.Int64
	var progressMu sync.Mutex
	var progressErr error
	emitProgress := func() {
		if err := emitRepresentationProgress(ctx, env, representationProgress{
			Phase: o.kind, Completed: completed.Load(), Total: int64(len(items)),
			CacheHits: cacheHits.Load(), CacheMisses: int64(len(missIndexes)),
			Active: active.Load(), Failed: failures.Load(),
		}); err != nil {
			progressMu.Lock()
			if progressErr == nil {
				progressErr = err
			}
			progressMu.Unlock()
		}
	}
	progressDone := make(chan struct{})
	var progressWG sync.WaitGroup
	if env.EmitEvent != nil {
		emitProgress()
		progressWG.Add(1)
		go func() {
			defer progressWG.Done()
			ticker := time.NewTicker(5 * time.Second)
			defer ticker.Stop()
			for {
				select {
				case <-ticker.C:
					emitProgress()
				case <-progressDone:
					return
				}
			}
		}()
	}
	jobs := make(chan int)
	var wg sync.WaitGroup
	var firstErr error
	var errOnce sync.Once
	workerCount := env.GenerationConcurrency
	if workerCount <= 0 {
		workerCount = defaultRepresentationGenerationWorkers
	}
	if workerCount > len(missIndexes) {
		workerCount = len(missIndexes)
	}
	for worker := 0; worker < workerCount; worker++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				if workCtx.Err() != nil {
					return
				}
				item := items[index]
				active.Add(1)
				value, err := env.Generator.Generate(workCtx, item.request)
				active.Add(-1)
				if err != nil {
					failures.Add(1)
					errOnce.Do(func() {
						firstErr = fmt.Errorf("RAG_GENERATION_FAILED: %s: %w", item.chunk.Record.ID, err)
						cancel()
					})
					return
				}
				if env.Cache != nil {
					raw, _ := json.Marshal(generationCacheEnvelopeV2{SchemaVersion: "rag-generation-cache-envelope/v2", Identity: item.cacheIdentity, Value: value})
					env.Cache.Put(item.cacheKey, raw)
				}
				reps, err := buildRepresentations(o, config, modelManifest, promptManifest, item, value, "miss", env)
				if err != nil {
					failures.Add(1)
					errOnce.Do(func() {
						firstErr = err
						cancel()
					})
					return
				}
				results[index] = representationResult{reps: reps, usage: value}
				completed.Add(1)
			}
		}()
	}
schedule:
	for _, index := range missIndexes {
		select {
		case jobs <- index:
		case <-workCtx.Done():
			break schedule
		}
	}
	close(jobs)
	wg.Wait()
	if env.EmitEvent != nil {
		close(progressDone)
		progressWG.Wait()
		emitProgress()
	}
	progressMu.Lock()
	err := progressErr
	progressMu.Unlock()
	if err != nil {
		return nil, err
	}
	if firstErr != nil {
		return nil, firstErr
	}
	for _, result := range results {
		if result.err != nil {
			return nil, result.err
		}
		addGenerationUsage(env, config.Model, result.usage)
		out = append(out, result.reps...)
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

// buildRepresentations converts a GenerationResult into Representation records
// for a single chunk. It handles schema validation for summaries and question
// count validation for synthetic questions. This is the per-chunk logic that
// was previously inline in the sequential loop; it is now called from both the
// cache-hit path and the concurrent generation goroutines.
func buildRepresentations(o representationOperator, config representationConfig, modelManifest ragcontract.ModelManifest, promptManifest ragcontract.PromptManifest, w representationWork, result GenerationResult, cacheOutcome string, env *Environment) ([]Representation, error) {
	if o.kind == "representations.structured-summary" {
		if result.Text == "" || !json.Valid([]byte(result.Text)) {
			return nil, fmt.Errorf("RAG_STRUCTURED_OUTPUT_INVALID: %s", w.chunk.Record.ID)
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
		derivation := derivation(o.Ref(), modelManifest.Digest, promptManifest.Digest, w.parentDigest, w.parentReprIDs, sourceRecordIDs(w.chunk), result.Text, result, cacheOutcome)
		return []Representation{newRepresentation(config.Name, "summary", w.chunk, result.Text, "derived", 0, derivation)}, nil
	}
	if len(result.Questions) != w.request.Count {
		return nil, fmt.Errorf("RAG_QUESTION_COUNT_MISMATCH: got %d want %d", len(result.Questions), w.request.Count)
	}
	reps := make([]Representation, 0, len(result.Questions))
	for ordinal, question := range result.Questions {
		if question == "" {
			return nil, fmt.Errorf("RAG_QUESTION_EMPTY")
		}
		derivation := derivation(o.Ref(), modelManifest.Digest, promptManifest.Digest, w.parentDigest, w.parentReprIDs, sourceRecordIDs(w.chunk), question, result, cacheOutcome)
		reps = append(reps, newRepresentation(config.Name, "question", w.chunk, question, "derived", ordinal, derivation))
	}
	return reps, nil
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
