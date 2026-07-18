package api

import "net/http"

func (h *handler) handleRAGArtifactCatalog(w http.ResponseWriter, r *http.Request) {
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
	defer func() { _ = rows.Close() }()
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
