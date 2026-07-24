package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/chunkenrichment"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/documentprocessing"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragintakeworkflow"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

func (h *handler) intakeApplication(r *http.Request) (*ragintakeworkflow.Application, error) {
	config := ragintakeworkflow.DefaultConfig(h.intakeConfig.DatabasePath)
	config.WorkflowDatabase = h.intakeConfig.WorkflowDB
	config.ArtifactRoot = h.intakeConfig.WorkflowArtifactRoot
	config.Runtime.IndexRoot = h.intakeConfig.IndexRoot
	return ragintakeworkflow.Open(r.Context(), config)
}
func (h *handler) handleListIntakeRuns(w http.ResponseWriter, r *http.Request) {
	app, err := h.intakeApplication(r)
	if err != nil {
		writeError(w, 500, "intake_open_failed", "intake workflow unavailable")
		return
	}
	defer func() { _ = app.Close() }()
	runs, err := app.ListRuns(r.Context(), r.URL.Query().Get("status"), intQueryDefault(r, "limit", 50))
	if err != nil {
		writeError(w, 500, "intake_query_failed", "intake runs unavailable")
		return
	}
	writeJSON(w, 200, map[string]any{"runs": runs})
}
func (h *handler) handleGetIntakeRun(w http.ResponseWriter, r *http.Request) {
	app, err := h.intakeApplication(r)
	if err != nil {
		writeError(w, 500, "intake_open_failed", "intake workflow unavailable")
		return
	}
	defer func() { _ = app.Close() }()
	view, err := app.Show(r.Context(), workflowv3.RunID(r.PathValue("id")))
	if err != nil {
		writeError(w, 404, "intake_run_not_found", "intake run not found")
		return
	}
	writeJSON(w, 200, view)
}
func (h *handler) handleIntakeObservations(w http.ResponseWriter, r *http.Request) {
	app, err := h.intakeApplication(r)
	if err != nil {
		writeError(w, 500, "intake_open_failed", "intake workflow unavailable")
		return
	}
	defer func() { _ = app.Close() }()
	value, err := app.Observations(r.Context(), workflowv3.RunID(r.PathValue("id")))
	if err != nil {
		writeError(w, 404, "intake_observations_not_found", "intake observations not found")
		return
	}
	writeJSON(w, 200, value)
}
func (h *handler) handleCancelIntakeRun(w http.ResponseWriter, r *http.Request) {
	app, err := h.intakeApplication(r)
	if err != nil {
		writeError(w, 500, "intake_open_failed", "intake workflow unavailable")
		return
	}
	defer func() { _ = app.Close() }()
	value, err := app.Cancel(r.Context(), workflowv3.RunID(r.PathValue("id")))
	if err != nil {
		writeError(w, 400, "intake_cancel_failed", "intake cancellation failed")
		return
	}
	writeJSON(w, 200, value)
}

type intakeSubmitRequest struct {
	RunID               string   `json:"run_id"`
	DocumentIDs         []string `json:"document_ids"`
	SourceIDs           []string `json:"source_ids"`
	DocumentLimit       int      `json:"document_limit"`
	Strategy            string   `json:"strategy"`
	ChunkSize           int      `json:"chunk_size"`
	Overlap             int      `json:"overlap"`
	SkipPreprocessing   bool     `json:"skip_preprocessing"`
	SkipChunkEnrichment bool     `json:"skip_chunk_enrichment"`
	SkipEmbeddings      bool     `json:"skip_embeddings"`
	SkipBM25            bool     `json:"skip_bm25"`
	Profile             string   `json:"profile"`
	EmbeddingType       string   `json:"embeddings_type"`
	EmbeddingEngine     string   `json:"embeddings_engine"`
	Dimensions          int      `json:"embeddings_dimensions"`
	BatchSize           int      `json:"batch_size"`
	IndexID             string   `json:"index_id"`
	ForceIndex          bool     `json:"force_index"`
}

func (h *handler) handleSubmitIntakeRun(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	var input intakeSubmitRequest
	if err := decoder.Decode(&input); err != nil {
		writeError(w, 400, "invalid_intake_request", "invalid intake request")
		return
	}
	var trailing any
	if err := decoder.Decode(&trailing); err != io.EOF {
		writeError(w, 400, "invalid_intake_request", "invalid intake request")
		return
	}
	if input.RunID == "" {
		input.RunID = "intake-" + time.Now().UTC().Format("20060102T150405.000000000")
	}
	if input.IndexID == "" {
		input.IndexID = "bm25-" + input.RunID
	}
	config := ragintakeworkflow.DefaultConfig(h.intakeConfig.DatabasePath)
	config.WorkflowDatabase = h.intakeConfig.WorkflowDB
	config.ArtifactRoot = h.intakeConfig.WorkflowArtifactRoot
	config.Runtime.IndexRoot = h.intakeConfig.IndexRoot
	request, err := ragintakeworkflow.PrepareRequest(r.Context(), config.Runtime, ragintakeworkflow.Request{Strategy: input.Strategy, ChunkSize: input.ChunkSize, Overlap: input.Overlap, SkipPreprocessing: input.SkipPreprocessing, SkipChunkEnrichment: input.SkipChunkEnrichment, SkipEmbeddings: input.SkipEmbeddings, SkipBM25: input.SkipBM25, Profile: input.Profile, EmbeddingType: input.EmbeddingType, EmbeddingEngine: input.EmbeddingEngine, Dimensions: input.Dimensions, BatchSize: input.BatchSize, IndexID: input.IndexID, ForceBM25: input.ForceIndex}, ragintakeworkflow.Selection{DocumentIDs: input.DocumentIDs, SourceIDs: input.SourceIDs, DocumentLimit: input.DocumentLimit, ChunksPerDocumentToEnrich: 1})
	if err != nil {
		writeError(w, 400, "intake_prepare_failed", fmt.Sprintf("intake request rejected: %s", err))
		return
	}
	app, err := ragintakeworkflow.Open(r.Context(), config)
	if err != nil {
		writeError(w, 500, "intake_open_failed", "intake workflow unavailable")
		return
	}
	defer func() { _ = app.Close() }()
	submission, err := app.SubmitRequest(r.Context(), request, workflowv3.RunID(input.RunID))
	if err != nil {
		writeError(w, 500, "intake_submit_failed", "intake submission failed")
		return
	}
	writeJSON(w, 201, map[string]any{"submission": submission, "document_ids": request.DocumentIDs, "strategy_id": fmt.Sprintf("%s-%d-%d", request.Strategy, request.ChunkSize, request.Overlap)})
}

func (h *handler) handleDocumentProcessingIdentities(w http.ResponseWriter, r *http.Request) {
	result, err := h.queries.ListDocumentProcessingIdentities()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
		return
	}
	if result == nil {
		result = []db.DocumentProcessingIdentity{}
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": result})
}
func (h *handler) handleDocumentProcessingCoverage(w http.ResponseWriter, r *http.Request) {
	service := documentprocessing.NewService(h.queries)
	result, err := service.Coverage(r.Context(), documentprocessing.CoverageRequest{ArtifactType: r.URL.Query().Get("artifact_type"), PromptVersion: r.URL.Query().Get("prompt_version"), Provider: r.URL.Query().Get("provider"), Model: r.URL.Query().Get("model")})
	if err != nil {
		writeError(w, 400, "coverage_failed", err.Error())
		return
	}
	writeJSON(w, 200, result)
}
func (h *handler) handleChunkEnrichmentIdentities(w http.ResponseWriter, r *http.Request) {
	result, err := h.queries.ListChunkEnrichmentIdentities()
	if err != nil {
		writeError(w, 500, "query_failed", err.Error())
		return
	}
	if result == nil {
		result = []db.ChunkEnrichmentIdentity{}
	}
	writeJSON(w, 200, map[string]any{"items": result})
}
func (h *handler) handleChunkEnrichmentCoverage(w http.ResponseWriter, r *http.Request) {
	service := chunkenrichment.NewService(h.queries)
	result, err := service.Coverage(r.Context(), chunkenrichment.CoverageRequest{StrategyID: r.URL.Query().Get("strategy_id"), PromptVersion: r.URL.Query().Get("prompt_version")})
	if err != nil {
		writeError(w, 400, "coverage_failed", err.Error())
		return
	}
	writeJSON(w, 200, result)
}
func (h *handler) handleDocumentProcessingArtifacts(w http.ResponseWriter, r *http.Request) {
	items, err := h.queries.ListDocumentProcessingArtifacts(r.PathValue("id"))
	if err != nil {
		writeError(w, 500, "query_failed", err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"document_id": r.PathValue("id"), "items": items})
}
func (h *handler) handleChunkEnrichments(w http.ResponseWriter, r *http.Request) {
	items, err := h.queries.ListChunkEnrichments(r.PathValue("id"), r.URL.Query().Get("strategy_id"), r.URL.Query().Get("prompt_version"))
	if err != nil {
		writeError(w, 500, "query_failed", err.Error())
		return
	}
	writeJSON(w, 200, map[string]any{"chunk_id": r.PathValue("id"), "items": items})
}
