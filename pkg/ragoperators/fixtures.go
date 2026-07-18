package ragoperators

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"unicode/utf16"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

const FixtureSummaryModel = "fixture-summary-v1"
const FixtureQuestionModel = "fixture-question-v1"
const FixtureEmbeddingModel = "fixture-hash-32-v1"
const FixtureSummaryPrompt = "fixture-transcript-summary-v1"
const FixtureQuestionPrompt = "fixture-transcript-questions-v1"
const FixtureSummarySchema = "transcript-rag-summary/v1"

var fixtureWordPattern = regexp.MustCompile(`[a-z0-9_./-]{4,}`)

type FixtureProviders struct{ Resolver StaticManifestResolver }

func NewFixtureProviders() FixtureProviders {
	manifest := func(id, digest string, dimensions int) ragcontract.ModelManifest {
		return ragcontract.ModelManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.ModelManifestSchema, Digest: digest, Parents: []ragcontract.ParentDigest{}}, ModelID: id, ModelDigest: digest, Dimensions: dimensions, Tokenization: "fixture-utf16", Truncation: "none", Normalization: "none", ImplementationVersion: "fixture/v1", RequestParameters: json.RawMessage(`{}`)}
	}
	prompt := func(id, digest, output string) ragcontract.PromptManifest {
		return ragcontract.PromptManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.PromptManifestSchema, Digest: digest, Parents: []ragcontract.ParentDigest{}}, PromptID: id, TemplateDigest: digest, InputSchema: "text/plain", OutputSchema: output}
	}
	return FixtureProviders{Resolver: StaticManifestResolver{Models: map[string]ragcontract.ModelManifest{
		FixtureSummaryModel: manifest(FixtureSummaryModel, "sha256:"+strings.Repeat("1", 64), 0), FixtureQuestionModel: manifest(FixtureQuestionModel, "sha256:"+strings.Repeat("2", 64), 0), FixtureEmbeddingModel: manifest(FixtureEmbeddingModel, "sha256:"+strings.Repeat("3", 64), 32)}, Prompts: map[string]ragcontract.PromptManifest{
		FixtureSummaryPrompt: prompt(FixtureSummaryPrompt, "sha256:"+strings.Repeat("4", 64), FixtureSummarySchema), FixtureQuestionPrompt: prompt(FixtureQuestionPrompt, "sha256:"+strings.Repeat("5", 64), "fixture-questions/v1")}}}
}
func (p FixtureProviders) Generate(_ context.Context, request GenerationRequest) (GenerationResult, error) {
	switch request.Kind {
	case "representations.structured-summary":
		words := uniqueFixtureWords(request.Text, 8)
		summary := map[string]any{"schema": FixtureSummarySchema, "abstract": fixtureAbstract(request.Text), "decisions": fixtureSentences(request.Text, []string{"decid", "choose", "selected", "will use"}), "problems": fixtureSentences(request.Text, []string{"error", "fail", "problem", "issue", "broken"}), "actions": fixtureSentences(request.Text, []string{"add", "fix", "run", "update", "implement", "create"}), "artifacts": []string{}, "questions": []string{}, "keywords": words}
		raw, _ := json.Marshal(summary)
		return GenerationResult{Text: string(raw), FinishReason: "fixture"}, nil
	case "representations.synthetic-questions":
		questions := make([]string, request.Count)
		words := uniqueFixtureWords(request.Text, request.Count)
		for index := range questions {
			if index == 0 {
				questions[index] = "What happened in this transcript passage?"
			} else {
				word := "this passage"
				if index-1 < len(words) {
					word = words[index-1]
				}
				questions[index] = fmt.Sprintf("Which details are recorded for %s?", word)
			}
		}
		return GenerationResult{Questions: questions, FinishReason: "fixture"}, nil
	default:
		return GenerationResult{}, fmt.Errorf("RAG_FIXTURE_GENERATION_KIND: %s", request.Kind)
	}
}
func (p FixtureProviders) Embed(_ context.Context, _ string, texts []string) ([][]float64, Usage, error) {
	vectors := make([][]float64, len(texts))
	for index, text := range texts {
		vector := make([]float64, 32)
		for position, code := range utf16.Encode([]rune(strings.ToLower(text))) {
			vector[(int(code)*31+position*17)%len(vector)]++
		}
		vectors[index] = vector
	}
	return vectors, Usage{}, nil
}
func (p FixtureProviders) Validate(schema string, document json.RawMessage) error {
	if schema != FixtureSummarySchema {
		return fmt.Errorf("RAG_FIXTURE_SCHEMA: %s", schema)
	}
	var value struct {
		Schema                                                       string `json:"schema"`
		Abstract                                                     string `json:"abstract"`
		Decisions, Problems, Actions, Artifacts, Questions, Keywords []string
	}
	if err := json.Unmarshal(document, &value); err != nil {
		return err
	}
	if value.Schema != schema || value.Abstract == "" {
		return fmt.Errorf("RAG_FIXTURE_SUMMARY_INVALID")
	}
	return nil
}
func fixtureAbstract(text string) string {
	sentences := strings.FieldsFunc(text, func(r rune) bool { return r == '.' || r == '!' || r == '?' })
	parts := []string{}
	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence != "" {
			parts = append(parts, sentence)
			if len(parts) == 2 {
				break
			}
		}
	}
	if len(parts) == 0 {
		return strings.TrimSpace(text)
	}
	return strings.Join(parts, ". ") + "."
}
func fixtureSentences(text string, needles []string) []string {
	result := []string{}
	for _, sentence := range strings.FieldsFunc(text, func(r rune) bool { return r == '.' || r == '!' || r == '?' }) {
		lower := strings.ToLower(sentence)
		for _, needle := range needles {
			if strings.Contains(lower, needle) {
				result = append(result, strings.TrimSpace(sentence))
				break
			}
		}
		if len(result) == 3 {
			break
		}
	}
	return result
}
func uniqueFixtureWords(text string, limit int) []string {
	seen := map[string]bool{}
	result := []string{}
	for _, word := range fixtureWordPattern.FindAllString(strings.ToLower(text), -1) {
		if !seen[word] {
			seen[word] = true
			result = append(result, word)
			if len(result) == limit {
				break
			}
		}
	}
	sort.Strings(result)
	return result
}
