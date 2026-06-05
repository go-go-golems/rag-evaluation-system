package db

import (
	"database/sql"
	"fmt"
)

type DocumentProcessingArtifact struct {
	DocumentID    string `json:"document_id"`
	SourceID      string `json:"source_id,omitempty"`
	ArtifactType  string `json:"artifact_type"`
	PromptVersion string `json:"prompt_version"`
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	InputHash     string `json:"input_hash"`
	OutputText    string `json:"output_text,omitempty"`
	OutputJSON    string `json:"output_json,omitempty"`
	Status        string `json:"status"`
	ErrorCode     string `json:"error_code,omitempty"`
	ErrorMessage  string `json:"error_message,omitempty"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

type DocumentProcessingCoverage struct {
	SourceID      string `json:"source_id"`
	DocumentCount int    `json:"document_count"`
	ArtifactCount int    `json:"artifact_count"`
	FreshCount    int    `json:"fresh_count"`
	FailedCount   int    `json:"failed_count"`
	MissingCount  int    `json:"missing_count"`
}

func (q *Queries) UpsertDocumentProcessingArtifact(a DocumentProcessingArtifact) error {
	if a.OutputJSON == "" {
		a.OutputJSON = "{}"
	}
	_, err := q.db.Exec(`
		INSERT INTO document_processing_artifacts (
			document_id, artifact_type, prompt_version, provider, model, input_hash,
			output_text, output_json, status, error_code, error_message
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(document_id, artifact_type, prompt_version, provider, model) DO UPDATE SET
			input_hash = excluded.input_hash,
			output_text = excluded.output_text,
			output_json = excluded.output_json,
			status = excluded.status,
			error_code = excluded.error_code,
			error_message = excluded.error_message,
			updated_at = datetime('now')
	`, a.DocumentID, a.ArtifactType, a.PromptVersion, a.Provider, a.Model, a.InputHash,
		a.OutputText, a.OutputJSON, a.Status, a.ErrorCode, a.ErrorMessage)
	if err != nil {
		return fmt.Errorf("upsert document processing artifact: %w", err)
	}
	return nil
}

func (q *Queries) GetDocumentProcessingArtifact(documentID, artifactType, promptVersion, provider, model string) (*DocumentProcessingArtifact, bool, error) {
	var a DocumentProcessingArtifact
	var outputText, outputJSON, errorCode, errorMessage sql.NullString
	err := q.db.QueryRow(`
		SELECT document_id, artifact_type, prompt_version, provider, model, input_hash,
		       output_text, output_json, status, error_code, error_message, created_at, updated_at
		FROM document_processing_artifacts
		WHERE document_id = ? AND artifact_type = ? AND prompt_version = ? AND provider = ? AND model = ?
	`, documentID, artifactType, promptVersion, provider, model).Scan(
		&a.DocumentID, &a.ArtifactType, &a.PromptVersion, &a.Provider, &a.Model, &a.InputHash,
		&outputText, &outputJSON, &a.Status, &errorCode, &errorMessage, &a.CreatedAt, &a.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("get document processing artifact: %w", err)
	}
	if outputText.Valid {
		a.OutputText = outputText.String
	}
	if outputJSON.Valid {
		a.OutputJSON = outputJSON.String
	}
	if errorCode.Valid {
		a.ErrorCode = errorCode.String
	}
	if errorMessage.Valid {
		a.ErrorMessage = errorMessage.String
	}
	return &a, true, nil
}

func (q *Queries) IsDocumentProcessingArtifactFresh(documentID, artifactType, promptVersion, provider, model, inputHash string) (bool, error) {
	a, ok, err := q.GetDocumentProcessingArtifact(documentID, artifactType, promptVersion, provider, model)
	if err != nil || !ok {
		return false, err
	}
	return a.Status == "succeeded" && a.InputHash == inputHash, nil
}

func (q *Queries) ListDocumentProcessingArtifacts(documentID string) ([]DocumentProcessingArtifact, error) {
	rows, err := q.db.Query(`
		SELECT document_id, artifact_type, prompt_version, provider, model, input_hash,
		       COALESCE(output_text, ''), COALESCE(output_json, '{}'), status,
		       COALESCE(error_code, ''), COALESCE(error_message, ''), created_at, updated_at
		FROM document_processing_artifacts
		WHERE document_id = ?
		ORDER BY artifact_type, prompt_version, provider, model
	`, documentID)
	if err != nil {
		return nil, fmt.Errorf("list document processing artifacts: %w", err)
	}
	defer func() { _ = rows.Close() }()
	ret := []DocumentProcessingArtifact{}
	for rows.Next() {
		var a DocumentProcessingArtifact
		if err := rows.Scan(&a.DocumentID, &a.ArtifactType, &a.PromptVersion, &a.Provider, &a.Model, &a.InputHash, &a.OutputText, &a.OutputJSON, &a.Status, &a.ErrorCode, &a.ErrorMessage, &a.CreatedAt, &a.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan document processing artifact: %w", err)
		}
		ret = append(ret, a)
	}
	return ret, rows.Err()
}

type DocumentProcessingIdentity struct {
	ArtifactType  string `json:"artifact_type"`
	PromptVersion string `json:"prompt_version"`
	Provider      string `json:"provider"`
	Model         string `json:"model"`
	ArtifactCount int    `json:"artifact_count"`
}

func (q *Queries) ListDocumentProcessingIdentities() ([]DocumentProcessingIdentity, error) {
	rows, err := q.db.Query(`
		SELECT artifact_type, prompt_version, provider, model, COUNT(*)
		FROM document_processing_artifacts
		GROUP BY artifact_type, prompt_version, provider, model
		ORDER BY artifact_type, prompt_version, provider, model
	`)
	if err != nil {
		return nil, fmt.Errorf("list document processing identities: %w", err)
	}
	defer func() { _ = rows.Close() }()
	ret := []DocumentProcessingIdentity{}
	for rows.Next() {
		var id DocumentProcessingIdentity
		if err := rows.Scan(&id.ArtifactType, &id.PromptVersion, &id.Provider, &id.Model, &id.ArtifactCount); err != nil {
			return nil, fmt.Errorf("scan document processing identity: %w", err)
		}
		ret = append(ret, id)
	}
	return ret, rows.Err()
}

func (q *Queries) ListDocumentProcessingCoverage(artifactType, promptVersion, provider, model string) ([]DocumentProcessingCoverage, error) {
	rows, err := q.db.Query(`
		SELECT d.source_id,
		       COUNT(d.id) AS document_count,
		       COUNT(a.document_id) AS artifact_count,
		       SUM(CASE WHEN a.status = 'succeeded' THEN 1 ELSE 0 END) AS fresh_count,
		       SUM(CASE WHEN a.status = 'failed' THEN 1 ELSE 0 END) AS failed_count
		FROM documents d
		LEFT JOIN document_processing_artifacts a ON a.document_id = d.id
		  AND a.artifact_type = ? AND a.prompt_version = ? AND a.provider = ? AND a.model = ?
		GROUP BY d.source_id
		ORDER BY d.source_id
	`, artifactType, promptVersion, provider, model)
	if err != nil {
		return nil, fmt.Errorf("list document processing coverage: %w", err)
	}
	defer func() { _ = rows.Close() }()
	ret := []DocumentProcessingCoverage{}
	for rows.Next() {
		var c DocumentProcessingCoverage
		if err := rows.Scan(&c.SourceID, &c.DocumentCount, &c.ArtifactCount, &c.FreshCount, &c.FailedCount); err != nil {
			return nil, fmt.Errorf("scan document processing coverage: %w", err)
		}
		c.MissingCount = c.DocumentCount - c.FreshCount
		ret = append(ret, c)
	}
	return ret, rows.Err()
}
