package geppetto

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-go-golems/geppetto/pkg/steps/ai/settings"
	"github.com/go-go-golems/geppetto/pkg/steps/ai/types"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

type staticPromptResolver map[string]string

func (r staticPromptResolver) PromptText(name string) (string, error) {
	value, ok := r[name]
	if !ok {
		return "", errors.New("prompt not found")
	}
	return value, nil
}

type staticSchemaResolver map[string]json.RawMessage

func (r staticSchemaResolver) Raw(name string) (json.RawMessage, error) {
	value, ok := r[name]
	if !ok {
		return nil, errors.New("schema not found")
	}
	return value, nil
}

func TestGeneratorConsumesStructuredStreamIntoCompleteResult(t *testing.T) {
	var request struct {
		Model          string  `json:"model"`
		Stream         bool    `json:"stream"`
		Temperature    float64 `json:"temperature"`
		ResponseFormat struct {
			Type string `json:"type"`
		} `json:"response_format"`
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("request path = %q, want /v1/chat/completions", r.URL.Path)
		}
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("decode request: %v", err)
		}
		w.Header().Set("Content-Type", "text/event-stream")
		_, _ = io.WriteString(w, "data: {\"model\":\"qwen3:8b\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"{\\\"summary\\\":\\\"\"}}]}\n\n")
		_, _ = io.WriteString(w, "data: {\"model\":\"qwen3:8b\",\"choices\":[{\"index\":0,\"delta\":{\"content\":\"ok\\\"}\"},\"finish_reason\":\"stop\"}],\"usage\":{\"prompt_tokens\":7,\"completion_tokens\":3}}\n\n")
		_, _ = io.WriteString(w, "data: [DONE]\n\n")
		w.(http.Flusher).Flush()
	}))
	defer server.Close()

	generator := newTestGenerator(t, server.URL+"/v1")
	result, err := generator.Generate(context.Background(), ragoperators.GenerationRequest{
		Kind:         "representations.structured-summary",
		Model:        "qwen3:8b",
		Prompt:       "summary",
		OutputSchema: "summary/v1",
		Text:         "Payroll adjustments correct wages.",
	})
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if !request.Stream {
		t.Error("request stream = false, want the configured streamed transport")
	}
	if request.Model != "qwen3:8b" || request.Temperature != 0 || request.ResponseFormat.Type != "json_schema" {
		t.Errorf("request = %#v, want qwen model, zero temperature, and JSON Schema output", request)
	}
	if result.Text != `{"summary":"ok"}` {
		t.Errorf("result text = %q", result.Text)
	}
	if result.InputTokens != 7 || result.OutputTokens != 3 || result.FinishReason != "stop" {
		t.Errorf("result = %#v, want final stream usage and finish reason", result)
	}
}

func TestGeneratorLiveOllamaRequestKinds(t *testing.T) {
	if os.Getenv("RAG_GENERATOR_LIVE_TEST") != "1" {
		t.Skip("set RAG_GENERATOR_LIVE_TEST=1 and RAG_GENERATOR_LIVE_BASE_URL to run against a live Ollama service")
	}
	baseURL := strings.TrimRight(os.Getenv("RAG_GENERATOR_LIVE_BASE_URL"), "/")
	if baseURL == "" {
		t.Fatal("RAG_GENERATOR_LIVE_BASE_URL is required when RAG_GENERATOR_LIVE_TEST=1")
	}
	model := os.Getenv("RAG_GENERATOR_LIVE_MODEL")
	if model == "" {
		model = "qwen3:8b"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()
	tests := []struct {
		name    string
		request ragoperators.GenerationRequest
		check   func(t *testing.T, result ragoperators.GenerationResult)
	}{
		{
			name:    "structured summary",
			request: ragoperators.GenerationRequest{Kind: "representations.structured-summary", Model: model, Prompt: "summary", OutputSchema: "summary/v1", Text: "A payroll adjustment corrects wages or deductions."},
			check: func(t *testing.T, result ragoperators.GenerationResult) {
				var output struct {
					Summary string `json:"summary"`
				}
				if err := json.Unmarshal([]byte(result.Text), &output); err != nil || strings.TrimSpace(output.Summary) == "" {
					t.Fatalf("text = %q, JSON summary error = %v", result.Text, err)
				}
			},
		},
		{
			name:    "synthetic questions",
			request: ragoperators.GenerationRequest{Kind: "representations.synthetic-questions", Model: model, Prompt: "questions", OutputSchema: "questions/v1", Text: "A payroll adjustment corrects wages or deductions.", Count: 4},
			check: func(t *testing.T, result ragoperators.GenerationResult) {
				if len(result.Questions) != 4 {
					t.Fatalf("questions = %#v, want exactly four", result.Questions)
				}
			},
		},
		{
			name:    "grounded answer",
			request: ragoperators.GenerationRequest{Kind: "generate.answer", Model: model, Prompt: "answer", OutputSchema: "answer/v1", Text: "What does a payroll adjustment correct?", Evidence: []ragoperators.Evidence{{Chunk: ragoperators.Chunk{Record: ragcontract.ChunkRecord{ID: "chunk-1"}, Text: "A payroll adjustment corrects wages or deductions."}}}},
			check: func(t *testing.T, result ragoperators.GenerationResult) {
				if !result.Abstained && strings.TrimSpace(result.Text) == "" {
					t.Fatal("answer is empty without abstention")
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := newTestGenerator(t, baseURL).Generate(ctx, test.request)
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if result.FinishReason == "" {
				t.Fatalf("finish reason is empty: %#v", result)
			}
			test.check(t, result)
		})
	}
}

func TestGeneratorHonorsCancellationAndRedactsProviderError(t *testing.T) {
	t.Run("cancellation", func(t *testing.T) {
		started := make(chan struct{})
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/event-stream")
			w.WriteHeader(http.StatusOK)
			w.(http.Flusher).Flush()
			close(started)
			<-r.Context().Done()
		}))
		defer server.Close()

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		result := make(chan error, 1)
		go func() {
			_, err := newTestGenerator(t, server.URL+"/v1").Generate(ctx, ragoperators.GenerationRequest{Kind: "representations.structured-summary", Model: "qwen3:8b", Prompt: "summary", OutputSchema: "summary/v1", Text: "test"})
			result <- err
		}()
		<-started
		cancel()
		select {
		case err := <-result:
			if !errors.Is(err, context.Canceled) {
				t.Fatalf("Generate() error = %v, want context cancellation", err)
			}
		case <-time.After(3 * time.Second):
			t.Fatal("Generate() did not return after cancellation")
		}
	})

	t.Run("provider error", func(t *testing.T) {
		const providerBody = "provider internal detail that must not escape"
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			http.Error(w, providerBody, http.StatusBadGateway)
		}))
		defer server.Close()

		_, err := newTestGenerator(t, server.URL+"/v1").Generate(context.Background(), ragoperators.GenerationRequest{Kind: "representations.structured-summary", Model: "qwen3:8b", Prompt: "summary", OutputSchema: "summary/v1", Text: "test"})
		if err == nil || !strings.Contains(err.Error(), "RAG_GEPPETTO_GENERATOR_PROVIDER") || strings.Contains(err.Error(), providerBody) {
			t.Fatalf("Generate() error = %v; provider response body must be redacted", err)
		}
	})
}

func newTestGenerator(t *testing.T, baseURL string) *Generator {
	t.Helper()
	in, err := settings.NewInferenceSettings()
	if err != nil {
		t.Fatal(err)
	}
	apiType := types.ApiTypeOpenAI
	model := "qwen3:8b"
	in.Chat.ApiType = &apiType
	in.Chat.Engine = &model
	in.Chat.Stream = true
	in.API.BaseUrls["openai-base-url"] = baseURL
	in.API.APIKeys["openai-api-key"] = "test-key"
	in.API.AllowHTTP["openai"] = true
	in.API.AllowLocalNetworks["openai"] = true
	generator, err := NewGenerator(in, staticPromptResolver{
		"summary":   "Return JSON only. Summarize the source text.",
		"questions": "Return exactly four JSON questions based on the source text.",
		"answer":    "Return a JSON answer grounded only in the supplied evidence, citing chunk IDs or abstaining.",
	}, staticSchemaResolver{
		"summary/v1":   json.RawMessage(`{"type":"object","properties":{"summary":{"type":"string"}},"required":["summary"],"additionalProperties":false}`),
		"questions/v1": json.RawMessage(`{"type":"object","properties":{"questions":{"type":"array","minItems":4,"maxItems":4,"items":{"type":"string"}}},"required":["questions"],"additionalProperties":false}`),
		"answer/v1":    json.RawMessage(`{"type":"object","properties":{"answer":{"type":"string"},"citationChunkIds":{"type":"array","items":{"type":"string"}},"abstained":{"type":"boolean"}},"required":["answer","citationChunkIds","abstained"],"additionalProperties":false}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	return generator
}
