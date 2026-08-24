package ragworkflow

import (
	"bytes"
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func DecodeResult(body []byte) (Result, error) {
	var result Result
	if err := strictJSON(bytes.TrimSpace(body), &result); err != nil {
		return Result{}, fmt.Errorf("RAG_WORKFLOW_RESULT_DECODE: %w", err)
	}
	if err := ValidateResult(result); err != nil {
		return Result{}, err
	}
	return result, nil
}

func ValidateResult(result Result) error {
	want, err := resultDigest(result)
	if err != nil || result.SchemaVersion != ResultSchema || result.Digest != want || result.ExecutionDigest == "" || result.CellID == "" || result.VariantID == "" || result.PreparationFingerprint == "" || len(result.Results) == 0 {
		return fmt.Errorf("RAG_WORKFLOW_RESULT_INVALID")
	}
	previous := ""
	for _, query := range result.Results {
		if query.QueryID == "" || query.QueryID <= previous || query.Trace.SchemaVersion != ragcontract.TraceSchemaVersion || query.Trace.Query.ID != query.QueryID {
			return fmt.Errorf("RAG_WORKFLOW_RESULT_QUERY")
		}
		previous = query.QueryID
	}
	return nil
}
