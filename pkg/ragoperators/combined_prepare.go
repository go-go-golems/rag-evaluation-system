package ragoperators

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

// combinedPreparationOperator prepares a deterministic batch of chunks with one
// structured provider call. It is deliberately a new operator: callers cannot
// accidentally reuse the independent summary/question cache or prepared bundle.
type combinedPreparationOperator struct{}

func (combinedPreparationOperator) Ref() ragcontract.OperatorRef {
	return ragcontract.OperatorRef{Kind: "representations.combined-summary-questions", Version: "v1"}
}

type combinedPreparationConfig struct {
	Model, Prompt, OutputSchema  string
	BatchSize, QuestionsPerChunk int
	MaxBatchRunes                int
}
type combinedBatchInput struct {
	ChunkID string `json:"chunkId"`
	Text    string `json:"text"`
}
type combinedBatch struct {
	chunks []Chunk
	text   string
}

func (combinedPreparationOperator) Execute(ctx context.Context, node ragcontract.Node, inputs map[string]any, env *Environment) (map[string]any, error) {
	chunks, ok := inputs["chunks"].([]Chunk)
	if !ok {
		return nil, fmt.Errorf("RAG_COMBINED_INPUT")
	}
	if env == nil || env.Generator == nil {
		return nil, fmt.Errorf("RAG_GENERATOR_UNAVAILABLE: combined preparation")
	}
	var cfg combinedPreparationConfig
	if err := decodeConfig(node.Config, &cfg); err != nil {
		return nil, err
	}
	if cfg.BatchSize < 1 || cfg.QuestionsPerChunk < 1 || cfg.MaxBatchRunes < 1 {
		return nil, fmt.Errorf("RAG_COMBINED_CONFIG")
	}
	model, err := resolveModel(env, cfg.Model)
	if err != nil {
		return nil, err
	}
	prompt, err := resolvePrompt(env, cfg.Prompt)
	if err != nil {
		return nil, err
	}
	if cfg.OutputSchema == "" {
		cfg.OutputSchema = prompt.OutputSchema
	}
	if cfg.OutputSchema == "" {
		return nil, fmt.Errorf("RAG_COMBINED_SCHEMA")
	}
	batches, err := makeCombinedBatches(chunks, cfg)
	if err != nil {
		return nil, err
	}
	results := make([][]Representation, len(batches))
	workers := env.GenerationConcurrency
	if workers < 1 {
		workers = defaultRepresentationGenerationWorkers
	}
	if workers > len(batches) {
		workers = len(batches)
	}
	workCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	jobs := make(chan int)
	var wg sync.WaitGroup
	var once sync.Once
	var firstErr error
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				batch := batches[index]
				parentDigest, digestErr := ragcontract.Digest(struct {
					Chunks []string
					Config json.RawMessage
				}{chunkDigests(batch.chunks), node.Config})
				if digestErr != nil {
					once.Do(func() { firstErr = digestErr; cancel() })
					return
				}
				result, generateErr := env.Generator.Generate(workCtx, GenerationRequest{Kind: "representations.combined-summary-questions", Model: model.ModelID, Prompt: prompt.PromptID, OutputSchema: cfg.OutputSchema, ParentID: parentDigest, Text: batch.text, Count: cfg.QuestionsPerChunk})
				if generateErr != nil {
					once.Do(func() { firstErr = generateErr; cancel() })
					return
				}
				representations, validateErr := validateCombinedBatch(batch.chunks, result.CombinedItems, cfg.QuestionsPerChunk, model, prompt)
				if validateErr != nil {
					once.Do(func() { firstErr = validateErr; cancel() })
					return
				}
				results[index] = representations
			}
		}()
	}
schedule:
	for i := range batches {
		select {
		case jobs <- i:
		case <-workCtx.Done():
			break schedule
		}
	}
	close(jobs)
	wg.Wait()
	if firstErr != nil {
		return nil, fmt.Errorf("RAG_COMBINED_PREPARATION: %w", firstErr)
	}
	out := []Representation{}
	for _, values := range results {
		out = append(out, values...)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Record.ID < out[j].Record.ID })
	parents := []ragcontract.ParentDigest{}
	if len(chunks) > 0 {
		parents = parent("chunk-set", chunks[0].ManifestDigest, ragcontract.ChunkSetManifestSchema)
	}
	return materializeRepresentations(node, parents, out)
}

func makeCombinedBatches(chunks []Chunk, cfg combinedPreparationConfig) ([]combinedBatch, error) {
	ordered := append([]Chunk(nil), chunks...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].Record.ID < ordered[j].Record.ID })
	batches := []combinedBatch{}
	for len(ordered) > 0 {
		inputs, used, runes := []combinedBatchInput{}, 0, 0
		for used < len(ordered) && len(inputs) < cfg.BatchSize {
			chunk := ordered[used]
			size := len([]rune(chunk.Text))
			if len(inputs) > 0 && runes+size > cfg.MaxBatchRunes {
				break
			}
			if size > cfg.MaxBatchRunes {
				return nil, fmt.Errorf("RAG_COMBINED_CHUNK_TOO_LARGE: %s", chunk.Record.ID)
			}
			inputs, runes, used = append(inputs, combinedBatchInput{ChunkID: chunk.Record.ID, Text: chunk.Text}), runes+size, used+1
		}
		data, err := json.Marshal(struct {
			Items []combinedBatchInput `json:"items"`
		}{inputs})
		if err != nil {
			return nil, err
		}
		batches = append(batches, combinedBatch{chunks: append([]Chunk(nil), ordered[:used]...), text: string(data)})
		ordered = ordered[used:]
	}
	return batches, nil
}
func chunkDigests(chunks []Chunk) []string {
	out := make([]string, len(chunks))
	for i, c := range chunks {
		out[i] = c.Record.TextDigest
	}
	return out
}
func validateCombinedBatch(chunks []Chunk, values []CombinedGenerationItem, questions int, model ragcontract.ModelManifest, prompt ragcontract.PromptManifest) ([]Representation, error) {
	byID := map[string]CombinedGenerationItem{}
	for _, value := range values {
		if value.ChunkID == "" || byID[value.ChunkID].ChunkID != "" {
			return nil, fmt.Errorf("RAG_COMBINED_RESPONSE_DUPLICATE")
		}
		byID[value.ChunkID] = value
	}
	if len(values) != len(chunks) {
		return nil, fmt.Errorf("RAG_COMBINED_RESPONSE_CARDINALITY: got %d want %d", len(values), len(chunks))
	}
	out := []Representation{}
	for _, chunk := range chunks {
		value, ok := byID[chunk.Record.ID]
		if !ok || value.Summary == "" || len(value.Questions) != questions {
			return nil, fmt.Errorf("RAG_COMBINED_RESPONSE_ITEM: %s", chunk.Record.ID)
		}
		out = append(out, newRepresentation("summary", "summary", chunk, value.Summary, "derived", 0, nil))
		for ordinal, question := range value.Questions {
			if question == "" {
				return nil, fmt.Errorf("RAG_COMBINED_RESPONSE_QUESTION: %s", chunk.Record.ID)
			}
			out = append(out, newRepresentation("question", "question", chunk, question, "derived", ordinal+1, nil))
		}
	}
	_ = model
	_ = prompt
	return out, nil
}
