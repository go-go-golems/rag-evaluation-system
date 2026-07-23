package geppetto

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	geppettoengine "github.com/go-go-golems/geppetto/pkg/inference/engine"
	enginefactory "github.com/go-go-golems/geppetto/pkg/inference/engine/factory"
	"github.com/go-go-golems/geppetto/pkg/steps/ai/settings"
	"github.com/go-go-golems/geppetto/pkg/turns"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflowops"
)

type PromptTextResolver interface{ PromptText(string) (string, error) }
type SchemaResolver interface {
	Raw(string) (json.RawMessage, error)
}
type Generator struct {
	settings *settings.InferenceSettings
	prompts  PromptTextResolver
	schemas  SchemaResolver
}

var _ ragoperators.TextGenerator = (*Generator)(nil)

func NewGenerator(base *settings.InferenceSettings, prompts PromptTextResolver, schemas SchemaResolver) (*Generator, error) {
	if base == nil || base.Chat == nil || base.API == nil || prompts == nil || schemas == nil {
		return nil, fmt.Errorf("RAG_GEPPETTO_GENERATOR_CONFIG")
	}
	return &Generator{settings: base.Clone(), prompts: prompts, schemas: schemas}, nil
}

type preparedGeneration struct {
	engine  geppettoengine.Engine
	turn    *turns.Turn
	request ragoperators.GenerationRequest
}

func (g *Generator) Generate(ctx context.Context, request ragoperators.GenerationRequest) (ragoperators.GenerationResult, error) {
	prepared, err := g.PrepareGeneration(request)
	if err != nil {
		return ragoperators.GenerationResult{}, err
	}
	return prepared.Execute(ctx)
}

func (g *Generator) PrepareGeneration(request ragoperators.GenerationRequest) (ragworkflowops.PreparedGeneration, error) {
	if g == nil || g.settings == nil {
		return nil, fmt.Errorf("RAG_GEPPETTO_GENERATOR_UNAVAILABLE")
	}
	prompt, err := g.prompts.PromptText(request.Prompt)
	if err != nil {
		return nil, err
	}
	schemaRaw, err := g.schemas.Raw(request.OutputSchema)
	if err != nil {
		return nil, err
	}
	var schema map[string]any
	if err := json.Unmarshal(schemaRaw, &schema); err != nil {
		return nil, fmt.Errorf("RAG_GENERATOR_SCHEMA: %w", err)
	}
	ss := g.settings.Clone()
	if ss.Chat == nil {
		return nil, fmt.Errorf("RAG_GENERATOR_CHAT_SETTINGS")
	}
	model := request.Model
	ss.Chat.Engine = &model
	temperature := float64(0)
	ss.Chat.Temperature = &temperature
	ss.Chat.StructuredOutputMode = settings.StructuredOutputModeJSONSchema
	ss.Chat.StructuredOutputSchema = string(schemaRaw)
	ss.Chat.StructuredOutputName = strings.ReplaceAll(request.OutputSchema, "/", "_")
	ss.Chat.StructuredOutputStrict = boolPtr(true)
	ss.Chat.StructuredOutputRequireValid = true
	engine, err := enginefactory.NewEngineFromSettings(ss)
	if err != nil {
		return nil, fmt.Errorf("RAG_GEPPETTO_GENERATOR_ENGINE: %w", err)
	}
	turn := turns.NewTurnBuilder().WithUserPrompt(buildUserPrompt(prompt, request)).Build()
	return &preparedGeneration{engine: engine, turn: turn, request: request}, nil
}

func (p *preparedGeneration) Execute(ctx context.Context) (ragoperators.GenerationResult, error) {
	out, inference, err := geppettoengine.RunInferenceWithResult(ctx, p.engine, p.turn)
	if err != nil {
		return ragoperators.GenerationResult{}, classifyProviderError(err)
	}
	text := lastAssistantText(out)
	result := ragoperators.GenerationResult{Text: text, FinishReason: "completed"}
	if inference != nil {
		if inference.Usage != nil {
			result.InputTokens = int64(inference.Usage.InputTokens)
			result.OutputTokens = int64(inference.Usage.OutputTokens)
		}
		if inference.StopReason != "" {
			result.FinishReason = inference.StopReason
		}
		if inference.Cost != nil {
			cost := *inference.Cost
			result.Cost = &cost
		}
	}
	if text == "" {
		return result, ragworkflowops.ProviderSucceededWithInvalidResult("RAG_GEPPETTO_GENERATOR_EMPTY")
	}
	switch p.request.Kind {
	case "representations.synthetic-questions":
		var value struct {
			Questions []string `json:"questions"`
		}
		if err := json.Unmarshal([]byte(text), &value); err != nil {
			return result, ragworkflowops.ProviderSucceededWithInvalidResult("RAG_GENERATOR_QUESTIONS_JSON")
		}
		result.Questions = value.Questions
	case "generate.answer":
		var value struct {
			Answer           string   `json:"answer"`
			CitationChunkIDs []string `json:"citationChunkIds"`
			Abstained        bool     `json:"abstained"`
		}
		if err := json.Unmarshal([]byte(text), &value); err != nil {
			return result, ragworkflowops.ProviderSucceededWithInvalidResult("RAG_GENERATOR_ANSWER_JSON")
		}
		result.Text, result.CitationChunkIDs, result.Abstained = value.Answer, value.CitationChunkIDs, value.Abstained
	case "representations.combined-summary-questions":
		var value struct {
			Items []ragoperators.CombinedGenerationItem `json:"items"`
		}
		if err := json.Unmarshal([]byte(text), &value); err != nil {
			return result, ragworkflowops.ProviderSucceededWithInvalidResult("RAG_GENERATOR_COMBINED_JSON")
		}
		result.CombinedItems = value.Items
	}
	return result, nil
}
func buildUserPrompt(template string, request ragoperators.GenerationRequest) string {
	switch request.Kind {
	case "representations.structured-summary":
		return renderSummaryPrompt(template, request.Text)
	case "representations.synthetic-questions":
		return renderQuestionsPrompt(template, request.Text)
	case "generate.answer":
		return renderAnswerPrompt(template, request.Text, request.Evidence)
	case "representations.combined-summary-questions":
		return renderLabeledTextPrompt(template, "CHUNKS JSON", request.Text)
	default:
		return renderTextPrompt(template, request.Text)
	}
}

func renderSummaryPrompt(template, source string) string {
	return renderLabeledTextPrompt(template, "SOURCE TEXT", source)
}

func renderQuestionsPrompt(template, source string) string {
	return renderLabeledTextPrompt(template, "SOURCE TEXT", source)
}

func renderTextPrompt(template, text string) string {
	return renderLabeledTextPrompt(template, "TEXT", text)
}

func renderLabeledTextPrompt(template, label, text string) string {
	return fmt.Sprintf("%s\n\n%s:\n%s", template, label, text)
}

func renderAnswerPrompt(template, question string, evidence []ragoperators.Evidence) string {
	var b strings.Builder
	b.WriteString(template)
	b.WriteString("\n\nQUESTION:\n")
	b.WriteString(question)
	b.WriteString("\n\nEVIDENCE:\n")
	for _, item := range evidence {
		fmt.Fprintf(&b, "[%s]\n%s\n\n", item.Chunk.Record.ID, item.Chunk.Text)
	}
	return b.String()
}
func lastAssistantText(t *turns.Turn) string {
	if t == nil {
		return ""
	}
	for i := len(t.Blocks) - 1; i >= 0; i-- {
		if t.Blocks[i].Kind == turns.BlockKindLLMText {
			if value, ok := t.Blocks[i].Payload[turns.PayloadKeyText].(string); ok {
				return value
			}
		}
	}
	return ""
}
func boolPtr(value bool) *bool { return &value }
