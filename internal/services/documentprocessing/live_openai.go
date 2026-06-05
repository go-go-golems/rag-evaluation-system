package documentprocessing

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	embeddingservice "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
)

type OpenAIResponsesProvider struct {
	Profile           string
	ProfileRegistries []string
	ProviderName      string
	ModelName         string
	MaxOutputTokens   int
}

func NewOpenAIResponsesProvider(ctx context.Context, profile string, registries []string) (*OpenAIResponsesProvider, error) {
	settings, closeFn, effective, err := embeddingservice.ResolveInferenceSettings(ctx, embeddingservice.ProviderConfig{Profile: profile, ProfileRegistries: registries})
	if err != nil {
		return nil, err
	}
	if closeFn != nil {
		defer func() { _ = closeFn() }()
	}
	if settings == nil || settings.Chat == nil || settings.Chat.Engine == nil || strings.TrimSpace(*settings.Chat.Engine) == "" {
		return nil, fmt.Errorf("profile %q does not define chat engine", profile)
	}
	if settings.API == nil || strings.TrimSpace(settings.API.APIKeys["openai-api-key"]) == "" {
		return nil, fmt.Errorf("profile %q does not define openai-api-key", profile)
	}
	return &OpenAIResponsesProvider{Profile: effective, ProviderName: "openai-responses", ModelName: *settings.Chat.Engine, MaxOutputTokens: 700}, nil
}

func (p *OpenAIResponsesProvider) Identity() ProviderIdentity {
	return ProviderIdentity{Provider: p.ProviderName, Model: p.ModelName}
}

func (p *OpenAIResponsesProvider) ProcessDocument(ctx context.Context, req ProviderRequest) (*ProviderResult, error) {
	prompt := fmt.Sprintf(`You are preprocessing a document for a RAG evaluation corpus.
Return a concise cleaned representation preserving factual content. Do not invent facts.

Document ID: %s
Title: %s
Artifact type: %s
Prompt version: %s

Content:
%s`, req.DocumentID, req.Title, req.ArtifactType, req.PromptVersion, limitString(req.Content, 4000))
	text, err := p.call(ctx, prompt)
	if err != nil {
		return nil, err
	}
	return &ProviderResult{OutputText: strings.TrimSpace(text), OutputJSON: map[string]any{"profile": p.Profile, "provider": p.ProviderName, "model": p.ModelName}}, nil
}

func (p *OpenAIResponsesProvider) call(ctx context.Context, prompt string) (string, error) {
	settings, closeFn, _, err := embeddingservice.ResolveInferenceSettings(ctx, embeddingservice.ProviderConfig{Profile: p.Profile, ProfileRegistries: p.ProfileRegistries})
	if err != nil {
		return "", err
	}
	if closeFn != nil {
		defer func() { _ = closeFn() }()
	}
	apiKey := settings.API.APIKeys["openai-api-key"]
	baseURL := "https://api.openai.com/v1"
	if settings.API != nil && strings.TrimSpace(settings.API.BaseUrls["openai-base-url"]) != "" {
		baseURL = strings.TrimRight(settings.API.BaseUrls["openai-base-url"], "/")
	}
	maxTokens := p.MaxOutputTokens
	if maxTokens <= 0 {
		maxTokens = 700
	}
	body := map[string]any{"model": p.ModelName, "input": prompt, "max_output_tokens": maxTokens}
	if settings.Inference != nil && settings.Inference.ReasoningEffort != nil {
		body["reasoning"] = map[string]any{"effort": *settings.Inference.ReasoningEffort}
	}
	payload, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, baseURL+"/responses", bytes.NewReader(payload))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 90 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("openai responses request failed: status=%s body=%s", resp.Status, limitString(string(data), 800))
	}
	var parsed struct {
		OutputText string `json:"output_text"`
		Output     []struct {
			Content []struct {
				Type string `json:"type"`
				Text string `json:"text"`
			} `json:"content"`
		} `json:"output"`
	}
	if err := json.Unmarshal(data, &parsed); err != nil {
		return "", err
	}
	if strings.TrimSpace(parsed.OutputText) != "" {
		return parsed.OutputText, nil
	}
	var b strings.Builder
	for _, item := range parsed.Output {
		for _, content := range item.Content {
			if strings.TrimSpace(content.Text) != "" {
				if b.Len() > 0 {
					b.WriteString("\n")
				}
				b.WriteString(content.Text)
			}
		}
	}
	if strings.TrimSpace(b.String()) == "" {
		return "", fmt.Errorf("openai responses returned empty text")
	}
	return b.String(), nil
}

func limitString(s string, limit int) string {
	if limit <= 0 || len(s) <= limit {
		return s
	}
	return s[:limit]
}
