package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-go-golems/rag-evaluation-system/internal/services/experimentrun"
)

func (h *handler) experimentRuns() *experimentrun.Service { return experimentrun.NewService(h.queries) }

func (h *handler) handleLabCatalog(w http.ResponseWriter, r *http.Request) {
	type snapshot struct {
		ID            string `json:"id"`
		DocumentCount int    `json:"document_count"`
		CreatedAt     string `json:"created_at"`
	}
	type chunkSet struct {
		ID, CorpusSnapshotID, ChunkPlanID string
		ChunkCount                        int
		CreatedAt                         string
	}
	type embeddingSet struct {
		ID, ChunkSetID, EmbeddingPlanID string
		EmbeddingCount                  int
		CreatedAt                       string
	}
	type bm25Artifact struct {
		ID, ChunkSetID string
		ChunkCount     int
		CreatedAt      string
	}
	result := struct {
		Snapshots     []snapshot     `json:"snapshots"`
		ChunkSets     []chunkSet     `json:"chunk_sets"`
		EmbeddingSets []embeddingSet `json:"embedding_sets"`
		BM25Artifacts []bm25Artifact `json:"bm25_artifacts"`
	}{}
	rows, err := h.queries.DB().QueryContext(r.Context(), `SELECT id,document_count,created_at FROM corpus_snapshots ORDER BY created_at DESC,id`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "catalog_failed", err.Error())
		return
	}
	for rows.Next() {
		var item snapshot
		if err := rows.Scan(&item.ID, &item.DocumentCount, &item.CreatedAt); err != nil {
			_ = rows.Close()
			writeError(w, http.StatusInternalServerError, "catalog_failed", err.Error())
			return
		}
		result.Snapshots = append(result.Snapshots, item)
	}
	_ = rows.Close()
	rows, err = h.queries.DB().QueryContext(r.Context(), `SELECT id,corpus_snapshot_id,chunk_plan_id,chunk_count,created_at FROM chunk_sets ORDER BY created_at DESC,id`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "catalog_failed", err.Error())
		return
	}
	for rows.Next() {
		var item chunkSet
		if err := rows.Scan(&item.ID, &item.CorpusSnapshotID, &item.ChunkPlanID, &item.ChunkCount, &item.CreatedAt); err != nil {
			_ = rows.Close()
			writeError(w, http.StatusInternalServerError, "catalog_failed", err.Error())
			return
		}
		result.ChunkSets = append(result.ChunkSets, item)
	}
	_ = rows.Close()
	rows, err = h.queries.DB().QueryContext(r.Context(), `SELECT id,chunk_set_id,embedding_plan_id,embedding_count,created_at FROM embedding_sets ORDER BY created_at DESC,id`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "catalog_failed", err.Error())
		return
	}
	for rows.Next() {
		var item embeddingSet
		if err := rows.Scan(&item.ID, &item.ChunkSetID, &item.EmbeddingPlanID, &item.EmbeddingCount, &item.CreatedAt); err != nil {
			_ = rows.Close()
			writeError(w, http.StatusInternalServerError, "catalog_failed", err.Error())
			return
		}
		result.EmbeddingSets = append(result.EmbeddingSets, item)
	}
	_ = rows.Close()
	rows, err = h.queries.DB().QueryContext(r.Context(), `SELECT id,chunk_set_id,chunk_count,created_at FROM retrieval_artifacts WHERE kind='bm25' ORDER BY created_at DESC,id`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "catalog_failed", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var item bm25Artifact
		if err := rows.Scan(&item.ID, &item.ChunkSetID, &item.ChunkCount, &item.CreatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "catalog_failed", err.Error())
			return
		}
		result.BM25Artifacts = append(result.BM25Artifacts, item)
	}
	if err := rows.Err(); err != nil {
		writeError(w, http.StatusInternalServerError, "catalog_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *handler) handleListExperimentSpecifications(w http.ResponseWriter, r *http.Request) {
	items, err := h.experimentRuns().ListSpecifications(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
func (h *handler) handleGetExperimentSpecification(w http.ResponseWriter, r *http.Request) {
	item, err := h.experimentRuns().GetSpecification(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}
func (h *handler) handleCreateExperimentSpecification(w http.ResponseWriter, r *http.Request) {
	var input experimentrun.SpecificationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	item, reused, err := h.experimentRuns().CreateSpecification(r.Context(), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, "create_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"item": item, "reused": reused})
}
func (h *handler) handleCreateExperimentRun(w http.ResponseWriter, r *http.Request) {
	item, err := h.experimentRuns().CreateRun(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "create_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, item)
}
func (h *handler) handleListExperimentRuns(w http.ResponseWriter, r *http.Request) {
	items, err := h.experimentRuns().ListRuns(r.Context(), r.URL.Query().Get("specification_id"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
func (h *handler) handleGetExperimentRun(w http.ResponseWriter, r *http.Request) {
	item, err := h.experimentRuns().GetRun(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, item)
}
func (h *handler) handleListExperimentRunTraces(w http.ResponseWriter, r *http.Request) {
	items, err := h.experimentRuns().ListQueryTraces(r.Context(), r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "list_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}
func (h *handler) handleAppendExperimentRunEvent(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Type    string          `json:"type"`
		Payload json.RawMessage `json:"payload"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	item, err := h.experimentRuns().AppendEvent(r.Context(), r.PathValue("id"), input.Type, input.Payload)
	if err != nil {
		writeError(w, http.StatusBadRequest, "append_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, item)
}
func (h *handler) handleRecordExperimentQueryTrace(w http.ResponseWriter, r *http.Request) {
	var input experimentrun.QueryTraceInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if err := h.experimentRuns().RecordQueryTrace(r.Context(), r.PathValue("id"), input); err != nil {
		writeError(w, http.StatusBadRequest, "trace_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{"ok": true})
}
func (h *handler) handleCompleteExperimentRun(w http.ResponseWriter, r *http.Request) {
	var input experimentrun.SummaryInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	item, err := h.experimentRuns().CompleteRun(r.Context(), r.PathValue("id"), input)
	if err != nil {
		writeError(w, http.StatusBadRequest, "complete_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, item)
}
func (h *handler) handleCompareExperimentRuns(w http.ResponseWriter, r *http.Request) {
	leftID, rightID := r.URL.Query().Get("left"), r.URL.Query().Get("right")
	if leftID == "" || rightID == "" {
		writeError(w, http.StatusBadRequest, "missing_run", "left and right run IDs are required")
		return
	}
	service := h.experimentRuns()
	left, err := service.GetRun(r.Context(), leftID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", err.Error())
		return
	}
	right, err := service.GetRun(r.Context(), rightID)
	if err != nil {
		writeError(w, http.StatusNotFound, "not_found", err.Error())
		return
	}
	leftTraces, err := service.ListQueryTraces(r.Context(), leftID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "comparison_failed", err.Error())
		return
	}
	rightTraces, err := service.ListQueryTraces(r.Context(), rightID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "comparison_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"left": left, "right": right, "left_traces": leftTraces, "right_traces": rightTraces})
}
