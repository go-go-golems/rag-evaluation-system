package ragintakeworkflow

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"
)

const (
	RequestSchema = "rag-intake-request/v1"
	ResultSchema  = "rag-intake-result/v1"
	SummarySchema = "rag-intake-summary/v1"
)

type Request struct {
	SchemaVersion           string   `json:"schemaVersion"`
	TargetDigest            string   `json:"targetDigest"`
	ProviderAuthorityDigest string   `json:"providerAuthorityDigest,omitempty"`
	DocumentIDs             []string `json:"documentIds"`
	SourceIDs               []string `json:"sourceIds,omitempty"`
	Strategy                string   `json:"strategy"`
	ChunkSize               int      `json:"chunkSize"`
	Overlap                 int      `json:"overlap"`
	SkipPreprocessing       bool     `json:"skipPreprocessing,omitempty"`
	ForcePreprocessing      bool     `json:"forcePreprocessing,omitempty"`
	PreprocessArtifactType  string   `json:"preprocessArtifactType"`
	PreprocessPromptVersion string   `json:"preprocessPromptVersion"`
	PreprocessProvider      string   `json:"preprocessProvider"`
	PreprocessModel         string   `json:"preprocessModel"`
	SkipChunkEnrichment     bool     `json:"skipChunkEnrichment,omitempty"`
	ForceChunkEnrichment    bool     `json:"forceChunkEnrichment,omitempty"`
	ChunkIDs                []string `json:"chunkIds,omitempty"`
	ChunkEnrichmentPrompt   string   `json:"chunkEnrichmentPrompt"`
	ChunkEnrichmentProvider string   `json:"chunkEnrichmentProvider"`
	ChunkEnrichmentModel    string   `json:"chunkEnrichmentModel"`
	SkipEmbeddings          bool     `json:"skipEmbeddings,omitempty"`
	ForceEmbeddings         bool     `json:"forceEmbeddings,omitempty"`
	ProfileRegistries       []string `json:"profileRegistries,omitempty"`
	Profile                 string   `json:"profile,omitempty"`
	BaseProfile             string   `json:"baseProfile,omitempty"`
	EmbeddingType           string   `json:"embeddingType,omitempty"`
	EmbeddingEngine         string   `json:"embeddingEngine,omitempty"`
	Dimensions              int      `json:"dimensions,omitempty"`
	BatchSize               int      `json:"batchSize"`
	EmbeddingLimit          int      `json:"embeddingLimit,omitempty"`
	CacheType               string   `json:"cacheType,omitempty"`
	SkipBM25                bool     `json:"skipBM25,omitempty"`
	ForceBM25               bool     `json:"forceBM25,omitempty"`
	IndexID                 string   `json:"indexId"`
	IndexLimit              int      `json:"indexLimit,omitempty"`
}

type Result struct {
	SchemaVersion string            `json:"schemaVersion"`
	Operation     string            `json:"operation"`
	DocumentID    string            `json:"documentId,omitempty"`
	ChunkID       string            `json:"chunkId,omitempty"`
	StrategyID    string            `json:"strategyId,omitempty"`
	IndexID       string            `json:"indexId,omitempty"`
	Counts        map[string]int    `json:"counts,omitempty"`
	Identity      map[string]string `json:"identity,omitempty"`
	SkippedFresh  bool              `json:"skippedFresh,omitempty"`
}

type Summary struct {
	SchemaVersion    string `json:"schemaVersion"`
	TargetDigest     string `json:"targetDigest"`
	DocumentCount    int    `json:"documentCount"`
	ChunkTargetCount int    `json:"chunkTargetCount"`
}

func TargetDigest(path string) string {
	sum := sha256.Sum256([]byte(path))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func NormalizeRequest(request Request, targetDigest string) (Request, error) {
	request.SchemaVersion = RequestSchema
	if request.TargetDigest == "" {
		request.TargetDigest = targetDigest
	}
	if request.TargetDigest != targetDigest || !strings.HasPrefix(targetDigest, "sha256:") {
		return Request{}, fmt.Errorf("RAG_INTAKE_TARGET")
	}
	request.DocumentIDs = normalizeIDs(request.DocumentIDs)
	request.SourceIDs = normalizeIDs(request.SourceIDs)
	request.ChunkIDs = normalizeIDs(request.ChunkIDs)
	request.ProfileRegistries = normalizeIDs(request.ProfileRegistries)
	if len(request.DocumentIDs) == 0 || len(request.DocumentIDs) > 10_000 || len(request.ChunkIDs) > 100_000 {
		return Request{}, fmt.Errorf("RAG_INTAKE_SELECTION")
	}
	if request.Strategy == "" {
		request.Strategy = "fixed"
	}
	if request.Strategy != "fixed" {
		return Request{}, fmt.Errorf("RAG_INTAKE_STRATEGY")
	}
	if request.ChunkSize == 0 {
		request.ChunkSize = 1200
	}
	if request.Overlap == 0 {
		request.Overlap = 150
	}
	if request.ChunkSize < 1 || request.ChunkSize > 1_000_000 || request.Overlap < 0 || request.Overlap >= request.ChunkSize {
		return Request{}, fmt.Errorf("RAG_INTAKE_CHUNK_POLICY")
	}
	if request.BatchSize == 0 {
		request.BatchSize = 16
	}
	if request.BatchSize < 1 || request.BatchSize > 10_000 {
		return Request{}, fmt.Errorf("RAG_INTAKE_BATCH")
	}
	if request.PreprocessArtifactType == "" {
		request.PreprocessArtifactType = "clean_text"
	}
	if request.PreprocessPromptVersion == "" {
		request.PreprocessPromptVersion = "v1"
	}
	if request.PreprocessProvider == "" {
		request.PreprocessProvider = "fake"
	}
	if request.PreprocessModel == "" {
		request.PreprocessModel = "fake-document-processor"
	}
	if request.ChunkEnrichmentProvider == "" {
		request.ChunkEnrichmentProvider = "fake"
	}
	if request.ChunkEnrichmentModel == "" {
		request.ChunkEnrichmentModel = "fake-chunk-enricher"
	}
	if request.ChunkEnrichmentPrompt == "" {
		request.ChunkEnrichmentPrompt = "v1"
	}
	if request.IndexID == "" {
		return Request{}, fmt.Errorf("RAG_INTAKE_INDEX_ID")
	}
	for _, value := range []string{request.PreprocessArtifactType, request.PreprocessPromptVersion, request.PreprocessProvider, request.PreprocessModel, request.ChunkEnrichmentProvider, request.ChunkEnrichmentModel, request.ChunkEnrichmentPrompt, request.IndexID, request.Profile, request.BaseProfile, request.EmbeddingType, request.EmbeddingEngine, request.CacheType} {
		if len(value) > 256 {
			return Request{}, fmt.Errorf("RAG_INTAKE_IDENTITY_LIMIT")
		}
	}
	return request, nil
}

func DecodeRequest(reader io.Reader, targetDigest string) (Request, error) {
	decoder := json.NewDecoder(io.LimitReader(reader, 4<<20))
	decoder.DisallowUnknownFields()
	var request Request
	if err := decoder.Decode(&request); err != nil {
		return Request{}, fmt.Errorf("RAG_INTAKE_REQUEST: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return Request{}, fmt.Errorf("RAG_INTAKE_REQUEST_TRAILING")
		}
		return Request{}, err
	}
	if request.SchemaVersion != RequestSchema {
		return Request{}, fmt.Errorf("RAG_INTAKE_REQUEST_SCHEMA")
	}
	normalized, err := NormalizeRequest(request, targetDigest)
	if err != nil {
		return Request{}, err
	}
	got, _ := json.Marshal(request)
	want, _ := json.Marshal(normalized)
	if !bytes.Equal(got, want) {
		return Request{}, fmt.Errorf("RAG_INTAKE_REQUEST_NONCANONICAL")
	}
	return request, nil
}

func normalizeIDs(values []string) []string {
	seen := map[string]bool{}
	ret := []string{}
	for _, raw := range values {
		for _, part := range strings.Split(raw, ",") {
			value := strings.TrimSpace(part)
			if value != "" && !seen[value] {
				seen[value] = true
				ret = append(ret, value)
			}
		}
	}
	sort.Strings(ret)
	return ret
}
