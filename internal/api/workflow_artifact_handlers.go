package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-go-golems/rag-evaluation-system/internal/services/chunkenrichment"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/documentprocessing"
	"github.com/go-go-golems/rag-evaluation-system/internal/workflow"
	"github.com/go-go-golems/scraper/pkg/engine/model"
	"github.com/go-go-golems/scraper/pkg/services/engineview"
)

func (h *handler) handleListWorkflows(w http.ResponseWriter, r *http.Request) {
	service := engineview.NewService(h.engineDB)
	result, err := service.ListWorkflows(r.Context(), engineview.ListWorkflowsOptions{
		Site:   "rag-eval",
		Status: model.WorkflowStatus(r.URL.Query().Get("status")),
		Limit:  intQueryDefault(r, "limit", 50),
		Offset: intQueryDefault(r, "offset", 0),
	})
	if err != nil {
		writeError(w, http.StatusInternalServerError, "workflow_query_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *handler) handleGetWorkflow(w http.ResponseWriter, r *http.Request) {
	service := engineview.NewService(h.engineDB)
	result, err := service.Workflow(r.Context(), model.WorkflowID(r.PathValue("id")))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "workflow_query_failed", err.Error())
		return
	}
	if result == nil {
		writeError(w, http.StatusNotFound, "not_found", "workflow not found")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *handler) handleWorkflowOps(w http.ResponseWriter, r *http.Request) {
	service := engineview.NewService(h.engineDB)
	result, err := service.WorkflowOps(r.Context(), model.WorkflowID(r.PathValue("id")))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "workflow_ops_query_failed", err.Error())
		return
	}
	if result == nil {
		result = []engineview.WorkflowOp{}
	}
	// Summarize: group by operation+status, return sample + counts
	type opGroup struct {
		Operation string `json:"operation"`
		Queue     string `json:"queue"`
		Status    string `json:"status"`
		Count     int    `json:"count"`
		Sample    *engineview.WorkflowOp `json:"sample,omitempty"`
	}
	groups := map[string]*opGroup{}
	order := []string{}
	for i := range result {
		input := result[i].Op.Input
		var opInput map[string]any
		_ = json.Unmarshal(input, &opInput)
		operation, _ := opInput["operation"].(string)
		key := operation + "|" + string(result[i].Status)
		if _, ok := groups[key]; !ok {
			groups[key] = &opGroup{
				Operation: operation,
				Queue:     string(result[i].Op.Queue),
				Status:    string(result[i].Status),
				Sample:    &result[i],
			}
			order = append(order, key)
		}
		groups[key].Count++
	}
	summary := make([]opGroup, 0, len(groups))
	for _, key := range order {
		summary = append(summary, *groups[key])
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"workflow_id": r.PathValue("id"),
		"total":       len(result),
		"groups":      summary,
	})
}

func (h *handler) handleGetOpResult(w http.ResponseWriter, r *http.Request) {
	service := engineview.NewService(h.engineDB)
	result, found, err := service.GetOpResult(r.Context(), model.WorkflowID(r.PathValue("id")), model.OpID(r.PathValue("opId")))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
		return
	}
	if !found || result == nil {
		writeError(w, http.StatusNotFound, "not_found", "op result not found")
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *handler) handleRetryOp(w http.ResponseWriter, r *http.Request) {
	service := engineview.NewService(h.engineDB)
	if err := service.RetryOp(r.Context(), model.WorkflowID(r.PathValue("id")), model.OpID(r.PathValue("opId"))); err != nil {
		writeError(w, http.StatusBadRequest, "retry_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "retried"})
}

func (h *handler) handleCancelWorkflow(w http.ResponseWriter, r *http.Request) {
	service := engineview.NewService(h.engineDB)
	if err := service.CancelWorkflow(r.Context(), model.WorkflowID(r.PathValue("id"))); err != nil {
		writeError(w, http.StatusBadRequest, "cancel_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "canceled"})
}

func (h *handler) handleSubmitIntake(w http.ResponseWriter, r *http.Request) {
	var req workflow.SubmitIntakeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_json", err.Error())
		return
	}
	if req.DBPath == "" {
		req.DBPath = "data/rag-eval.db"
	}
	if req.EngineDB == "" {
		req.EngineDB = h.engineDB
	}

	result, err := workflow.SubmitIntakeWorkflow(r.Context(), req)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "submit_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusCreated, result)
}

func (h *handler) handleListQueues(w http.ResponseWriter, r *http.Request) {
	service := engineview.NewService(h.engineDB)
	queues, err := service.ListQueues(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
		return
	}
	if queues == nil {
		queues = []engineview.QueueStatus{}
	}
	writeJSON(w, http.StatusOK, map[string]interface{}{"queues": queues})
}

func (h *handler) handleDocumentProcessingCoverage(w http.ResponseWriter, r *http.Request) {
	service := documentprocessing.NewService(h.queries)
	result, err := service.Coverage(r.Context(), documentprocessing.CoverageRequest{
		ArtifactType:  r.URL.Query().Get("artifact_type"),
		PromptVersion: r.URL.Query().Get("prompt_version"),
		Provider:      r.URL.Query().Get("provider"),
		Model:         r.URL.Query().Get("model"),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "coverage_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *handler) handleChunkEnrichmentCoverage(w http.ResponseWriter, r *http.Request) {
	service := chunkenrichment.NewService(h.queries)
	result, err := service.Coverage(r.Context(), chunkenrichment.CoverageRequest{
		StrategyID:    r.URL.Query().Get("strategy_id"),
		PromptVersion: r.URL.Query().Get("prompt_version"),
	})
	if err != nil {
		writeError(w, http.StatusBadRequest, "coverage_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, result)
}

func (h *handler) handleDocumentProcessingArtifacts(w http.ResponseWriter, r *http.Request) {
	items, err := h.queries.ListDocumentProcessingArtifacts(r.PathValue("id"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"document_id": r.PathValue("id"), "items": items})
}

func (h *handler) handleChunkEnrichments(w http.ResponseWriter, r *http.Request) {
	items, err := h.queries.ListChunkEnrichments(r.PathValue("id"), r.URL.Query().Get("strategy_id"), r.URL.Query().Get("prompt_version"))
	if err != nil {
		writeError(w, http.StatusInternalServerError, "query_failed", err.Error())
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"chunk_id": r.PathValue("id"), "items": items})
}
