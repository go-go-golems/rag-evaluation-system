package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-go-golems/rag-evaluation-system/internal/chunking"
	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/ingest"
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

	if req.ID == "" || req.Name == "" || req.Type == "" {
		writeError(w, http.StatusBadRequest, "missing_fields", "id, name, and type are required")
		return
	}

	configJSON, _ := json.Marshal(req.Config)

	if err := h.queries.InsertSource(req.ID, req.Name, req.Type, string(configJSON)); err != nil {
		writeError(w, http.StatusInternalServerError, "insert_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"id":   req.ID,
		"name": req.Name,
	})
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
	if req.Dir == "" {
		writeError(w, http.StatusBadRequest, "missing_dir", "dir is required")
		return
	}

	scanner := ingest.NewScanner(h.queries)
	docIDs, err := scanner.ScanDir(sourceID, req.Dir)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "scan_failed", err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"source_id":      sourceID,
		"documents":      docIDs,
		"document_count": len(docIDs),
	})
}

// --- Documents ---

func (h *handler) handleListDocuments(w http.ResponseWriter, r *http.Request) {
	docs, err := h.queries.ListDocuments(50, 0)
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

	doc, err := h.queries.GetDocument(id)
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

	chunks, err := h.queries.ListChunks(docID)
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
		Strategy  string `json:"strategy"`
		ChunkSize int    `json:"chunk_size"`
		Overlap   int    `json:"overlap"`
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

	strategyID := fmt.Sprintf("%s-%d-%d", req.Strategy, req.ChunkSize, req.Overlap)

	chunker, err := chunking.NewChunkerFromType(req.Strategy, req.ChunkSize, req.Overlap, strategyID)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid_strategy", err.Error())
		return
	}

	chunks, err := chunker.Chunk(docID, content)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "chunk_failed", err.Error())
		return
	}

	// Store chunks
	for _, ch := range chunks {
		boundariesJSON, _ := json.Marshal(map[string]interface{}{"strategy_id": strategyID})
		h.queries.InsertChunk(ch.ID, ch.DocumentID, ch.ChunkIndex, ch.Text,
			ch.TokenCount, ch.StartOffset, ch.EndOffset, string(boundariesJSON))
	}

	h.queries.UpdateDocumentStatus(docID, "chunked")

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"document_id": docID,
		"strategy":    strategyID,
		"chunk_count": len(chunks),
		"chunks":      chunks,
	})
}

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
