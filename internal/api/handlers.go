package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	chunkcore "github.com/go-go-golems/rag-evaluation-system/internal/chunking"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	chunkservice "github.com/go-go-golems/rag-evaluation-system/internal/services/chunking"
	documentservice "github.com/go-go-golems/rag-evaluation-system/internal/services/document"
	embeddingservice "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	sourceservice "github.com/go-go-golems/rag-evaluation-system/internal/services/source"
)

// RegisterHandlers wires all API routes into the given mux
func RegisterHandlers(mux *http.ServeMux, database *sql.DB) {
	queries := db.NewQueries(database)
	h := &handler{queries: queries}

	// Health check
	mux.HandleFunc("GET /api/v1/health", h.handleHealth)

	// Sources
	mux.HandleFunc("GET /api/v1/sources", h.handleListSources)
	mux.HandleFunc("POST /api/v1/sources", h.handleCreateSource)

	// Documents
	mux.HandleFunc("GET /api/v1/documents", h.handleListDocuments)
	mux.HandleFunc("GET /api/v1/documents/{id}", h.handleGetDocument)
	mux.HandleFunc("GET /api/v1/documents/{id}/chunks", h.handleListChunks)

	// Source scan (ingest files from a directory)
	mux.HandleFunc("POST /api/v1/sources/{id}/scan", h.handleScanSource)

	// Chunking
	mux.HandleFunc("POST /api/v1/documents/{id}/chunk", h.handleChunkDocument)
	mux.HandleFunc("GET /api/v1/chunking-strategies", h.handleListChunkingStrategies)

	// Embeddings
	mux.HandleFunc("POST /api/v1/embeddings/compute", h.handleComputeEmbeddings)
	mux.HandleFunc("POST /api/v1/embeddings/coverage", h.handleEmbeddingCoverage)
	mux.HandleFunc("POST /api/v1/embeddings/similarity", h.handleEmbeddingSimilarity)
}

type handler struct {
	queries *db.Queries
}

func (h *handler) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":  "ok",
		"service": "rag-evaluation-system",
	})
}

// --- Sources ---

func (h *handler) handleListSources(w http.ResponseWriter, r *http.Request) {
	sources, err := h.queries.ListSources()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items": sources,
	})
}

func (h *handler) handleCreateSource(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ID     string                 `json:"id"`
		Name   string                 `json:"name"`
		Type   string                 `json:"type"`
		Config map[string]interface{} `json:"config"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	service := sourceservice.NewService(h.queries)
	result, err := service.Create(r.Context(), sourceservice.CreateRequest{
		ID:     req.ID,
		Name:   req.Name,
		Type:   req.Type,
		Config: req.Config,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "create_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, result)
}

func (h *handler) handleScanSource(w http.ResponseWriter, r *http.Request) {
	sourceID := r.PathValue("id")

	var req struct {
		Dir string `json:"dir"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	service := sourceservice.NewService(h.queries)
	result, err := service.Scan(r.Context(), sourceservice.ScanRequest{SourceID: sourceID, Dir: req.Dir})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "scan_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// --- Documents ---

func (h *handler) handleListDocuments(w http.ResponseWriter, r *http.Request) {
	service := documentservice.NewService(h.queries)
	docs, err := service.List(r.Context(), documentservice.ListRequest{Limit: 50, Offset: 0})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items": docs,
		"page":  map[string]interface{}{"limit": 50, "offset": 0},
	})
}

func (h *handler) handleGetDocument(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	service := documentservice.NewService(h.queries)
	doc, err := service.Get(r.Context(), id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
		return
	}
	if doc == nil {
		writeError(w, http.StatusNotFound, "not_found", "document not found")
		return
	}

	result := map[string]interface{}{
		"id":           doc.ID,
		"source_id":    doc.SourceID,
		"external_id":  doc.ExternalID,
		"title":        doc.Title,
		"author":       doc.Author,
		"url":          doc.URL,
		"content_type": doc.ContentType,
		"word_count":   doc.WordCount,
		"language":     doc.Language,
		"status":       doc.Status,
		"created_at":   doc.CreatedAt,
		"updated_at":   doc.UpdatedAt,
	}

	writeJSON(w, http.StatusOK, result)
}

// --- Chunks ---

func (h *handler) handleListChunks(w http.ResponseWriter, r *http.Request) {
	docID := r.PathValue("id")

	service := documentservice.NewService(h.queries)
	chunks, err := service.Chunks(r.Context(), documentservice.ChunksRequest{DocumentID: docID})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items":       chunks,
		"document_id": docID,
		"chunk_count": len(chunks),
	})
}

// --- Chunking ---

func (h *handler) handleChunkDocument(w http.ResponseWriter, r *http.Request) {
	docID := r.PathValue("id")

	var req struct {
		Strategy     string `json:"strategy"`
		ChunkSize    int    `json:"chunk_size"`
		Overlap      int    `json:"overlap"`
		Emit         string `json:"emit"`
		PreviewRunes int    `json:"preview_runes"`
		Limit        int    `json:"limit"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if req.Strategy == "" {
		req.Strategy = "fixed"
	}
	if req.ChunkSize == 0 {
		req.ChunkSize = 500
	}

	content, err := h.queries.GetDocumentContent(docID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
		return
	}
	if content == "" {
		writeError(w, http.StatusNotFound, "no_content", "document has no content")
		return
	}

	service := chunkservice.NewService(h.queries)
	result, err := service.Apply(r.Context(), chunkservice.ApplyRequest{
		DocumentID:  docID,
		Strategy:    req.Strategy,
		ChunkSize:   req.ChunkSize,
		Overlap:     req.Overlap,
		Description: "HTTP-created chunking strategy",
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "chunk_failed", err.Error())
		return
	}

	if req.Emit == "" {
		req.Emit = "summary"
	}
	if req.PreviewRunes == 0 {
		req.PreviewRunes = 120
	}
	if req.Limit == 0 {
		req.Limit = 50
	}

	response := map[string]interface{}{
		"document_id": result.DocumentID,
		"strategy_id": result.StrategyID,
		"chunk_count": result.ChunkCount,
	}

	switch req.Emit {
	case "summary", "none":
		// no chunk text in response
	case "preview":
		response["chunks"] = previewChunks(result.Chunks, req.PreviewRunes, req.Limit)
		response["emitted_count"] = minInt(len(result.Chunks), req.Limit)
	case "full":
		if req.Limit > 0 && len(result.Chunks) > req.Limit {
			response["chunks"] = result.Chunks[:req.Limit]
			response["emitted_count"] = req.Limit
		} else {
			response["chunks"] = result.Chunks
			response["emitted_count"] = len(result.Chunks)
		}
	default:
		writeError(w, http.StatusBadRequest, "invalid_emit", "emit must be summary, preview, full, or none")
		return
	}

	writeJSON(w, http.StatusOK, response)
}

// --- Embeddings ---

func (h *handler) handleComputeEmbeddings(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StrategyID        string   `json:"strategy_id"`
		SourceIDs         []string `json:"source_ids"`
		ProfileRegistries []string `json:"profile_registries"`
		Profile           string   `json:"profile"`
		BaseProfile       string   `json:"base_profile"`
		EmbeddingType     string   `json:"embeddings_type"`
		EmbeddingEngine   string   `json:"embeddings_engine"`
		Dimensions        int      `json:"embeddings_dimensions"`
		APIKey            string   `json:"api_key"`
		BaseURL           string   `json:"base_url"`
		CacheType         string   `json:"cache_type"`
		CacheDirectory    string   `json:"cache_directory"`
		BatchSize         int      `json:"batch_size"`
		Limit             int      `json:"limit"`
		Force             bool     `json:"force"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if req.EmbeddingType == "" {
		req.EmbeddingType = "ollama"
	}
	if req.EmbeddingEngine == "" {
		req.EmbeddingEngine = "nomic-embed-text"
	}
	if req.Dimensions == 0 {
		req.Dimensions = 768
	}
	if req.CacheType == "" {
		req.CacheType = "none"
	}
	if req.CacheDirectory == "" {
		req.CacheDirectory = "state/embedding-cache"
	}

	resolved, err := embeddingservice.ResolveProvider(r.Context(), embeddingservice.ProviderConfig{
		ProfileRegistries: req.ProfileRegistries,
		Profile:           req.Profile,
		BaseProfile:       req.BaseProfile,
		Type:              req.EmbeddingType,
		Engine:            req.EmbeddingEngine,
		Dimensions:        req.Dimensions,
		APIKey:            req.APIKey,
		BaseURL:           req.BaseURL,
		CacheType:         req.CacheType,
		CacheDirectory:    req.CacheDirectory,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "provider_failed", err.Error())
		return
	}
	if resolved.Close != nil {
		defer func() { _ = resolved.Close() }()
	}

	service := embeddingservice.NewService(h.queries)
	result, err := service.Compute(r.Context(), embeddingservice.ComputeRequest{
		StrategyID:   req.StrategyID,
		SourceIDs:    req.SourceIDs,
		Provider:     resolved.Provider,
		ProviderType: resolved.ProviderType,
		BatchSize:    req.BatchSize,
		Limit:        req.Limit,
		Force:        req.Force,
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "compute_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"strategy_id":       result.StrategyID,
		"source_ids":        result.SourceIDs,
		"provider_type":     result.ProviderType,
		"model":             result.Model,
		"dimensions":        result.Dimensions,
		"effective_profile": resolved.EffectiveProfile,
		"considered":        result.Considered,
		"computed":          result.Computed,
		"skipped_fresh":     result.SkippedFresh,
	})
}

func (h *handler) handleEmbeddingCoverage(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StrategyID   string `json:"strategy_id"`
		ProviderType string `json:"provider_type"`
		Model        string `json:"model"`
		Dimensions   int    `json:"dimensions"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}

	service := embeddingservice.NewService(h.queries)
	result, err := service.Coverage(r.Context(), embeddingservice.CoverageRequest{
		StrategyID:   req.StrategyID,
		ProviderType: req.ProviderType,
		Model:        req.Model,
		Dimensions:   req.Dimensions,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "coverage_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *handler) handleEmbeddingSimilarity(w http.ResponseWriter, r *http.Request) {
	var req struct {
		StrategyID     string `json:"strategy_id"`
		ProviderType   string `json:"provider_type"`
		Model          string `json:"model"`
		Dimensions     int    `json:"dimensions"`
		ChunkIDA       string `json:"chunk_id_a"`
		ChunkIDB       string `json:"chunk_id_b"`
		Limit          int    `json:"limit"`
		CandidateLimit int    `json:"candidate_limit"`
		PreviewRunes   int    `json:"preview_runes"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if req.Limit == 0 {
		req.Limit = 20
	}
	if req.CandidateLimit == 0 {
		req.CandidateLimit = 200
	}
	if req.PreviewRunes == 0 {
		req.PreviewRunes = 120
	}

	service := embeddingservice.NewService(h.queries)
	result, err := service.Similarity(r.Context(), embeddingservice.SimilarityRequest{
		StrategyID:     req.StrategyID,
		ProviderType:   req.ProviderType,
		Model:          req.Model,
		Dimensions:     req.Dimensions,
		ChunkIDA:       req.ChunkIDA,
		ChunkIDB:       req.ChunkIDB,
		Limit:          req.Limit,
		CandidateLimit: req.CandidateLimit,
		PreviewRunes:   req.PreviewRunes,
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "similarity_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, result)
}

// --- Chunking strategies ---

func (h *handler) handleListChunkingStrategies(w http.ResponseWriter, r *http.Request) {
	rows, err := h.queries.DB().QueryContext(r.Context(), `
		SELECT id, name, type, description, created_at
		FROM chunking_strategies ORDER BY created_at DESC
	`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
		return
	}
	defer rows.Close()

	type strategy struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Type        string `json:"type"`
		Description string `json:"description"`
		CreatedAt   string `json:"created_at"`
	}

	var strategies []strategy
	for rows.Next() {
		var s strategy
		if err := rows.Scan(&s.ID, &s.Name, &s.Type, &s.Description, &s.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "scan_failed", err.Error())
			return
		}
		strategies = append(strategies, s)
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"items": strategies,
	})
}

// --- Helpers ---

type chunkPreview struct {
	ID          string `json:"id"`
	DocumentID  string `json:"document_id"`
	StrategyID  string `json:"strategy_id"`
	ChunkIndex  int    `json:"chunk_index"`
	TextPreview string `json:"text_preview"`
	TokenCount  int    `json:"token_count"`
	StartOffset int    `json:"start_offset"`
	EndOffset   int    `json:"end_offset"`
}

func previewChunks(chunks []chunkcore.Chunk, previewRunes, limit int) []chunkPreview {
	if limit > 0 && len(chunks) > limit {
		chunks = chunks[:limit]
	}
	previews := make([]chunkPreview, 0, len(chunks))
	for _, ch := range chunks {
		previews = append(previews, chunkPreview{
			ID:          ch.ID,
			DocumentID:  ch.DocumentID,
			StrategyID:  ch.StrategyID,
			ChunkIndex:  ch.ChunkIndex,
			TextPreview: truncateRunes(ch.Text, previewRunes),
			TokenCount:  ch.TokenCount,
			StartOffset: ch.StartOffset,
			EndOffset:   ch.EndOffset,
		})
	}
	return previews
}

func truncateRunes(s string, maxRunes int) string {
	if maxRunes <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= maxRunes {
		return s
	}
	return string(runes[:maxRunes]) + "..."
}

func minInt(a, b int) int {
	if b <= 0 || a < b {
		return a
	}
	return b
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error":   code,
		"message": message,
	})
}
