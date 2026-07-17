package ragcontract

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
)

func DecodePipeline(reader io.Reader) (PipelineIR, error) {
	var value PipelineIR
	if err := decodeStrict(reader, &value); err != nil {
		return value, err
	}
	return value, RequireValid(ValidatePipeline(value))
}
func DecodeProduct(reader io.Reader) (ProductPlan, error) {
	var value ProductPlan
	if err := decodeStrict(reader, &value); err != nil {
		return value, err
	}
	if value.SchemaVersion != ProductSchemaVersion {
		return value, fmt.Errorf("RAG_V2_PRODUCT_SCHEMA: expected %s", ProductSchemaVersion)
	}
	return value, RequireValid(ValidatePipeline(value.Pipeline))
}
func DecodeStudy(reader io.Reader) (Study, error) {
	var value Study
	if err := decodeStrict(reader, &value); err != nil {
		return value, err
	}
	if value.SchemaVersion != StudySchemaVersion {
		return value, fmt.Errorf("RAG_V2_STUDY_SCHEMA: expected %s", StudySchemaVersion)
	}
	if len(value.Variants) == 0 {
		return value, fmt.Errorf("RAG_V2_STUDY_VARIANTS: at least one variant is required")
	}
	return value, nil
}
func DecodeExecution(reader io.Reader) (PipelineExecution, error) {
	var value PipelineExecution
	if err := decodeStrict(reader, &value); err != nil {
		return value, err
	}
	if value.SchemaVersion != ExecutionSchemaVersion {
		return value, fmt.Errorf("RAG_V2_EXECUTION_SCHEMA: expected %s", ExecutionSchemaVersion)
	}
	return value, RequireValid(ValidatePipeline(value.Pipeline))
}
func DecodeTrace(reader io.Reader) (QueryTrace, error) {
	var value QueryTrace
	if err := decodeStrict(reader, &value); err != nil {
		return value, err
	}
	if value.SchemaVersion != TraceSchemaVersion {
		return value, fmt.Errorf("RAG_V2_TRACE_SCHEMA: expected %s", TraceSchemaVersion)
	}
	if value.Query.ID == "" || value.Query.TextDigest == "" {
		return value, fmt.Errorf("RAG_V2_TRACE_QUERY: id and text digest are required")
	}
	return value, nil
}

func decodeStrict(reader io.Reader, target any) error {
	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	if err := decoder.Decode(target); err != nil {
		return fmt.Errorf("RAG_V2_DECODE: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		if err == nil {
			return fmt.Errorf("RAG_V2_DECODE: trailing JSON value")
		}
		return fmt.Errorf("RAG_V2_DECODE: trailing content: %w", err)
	}
	return nil
}

func CanonicalJSON(value any) ([]byte, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("RAG_V2_CANONICAL_JSON: %w", err)
	}
	var normalized any
	decoder := json.NewDecoder(bytes.NewReader(encoded))
	decoder.UseNumber()
	if err := decoder.Decode(&normalized); err != nil {
		return nil, err
	}
	return json.Marshal(normalized)
}

func CanonicalRaw(value json.RawMessage, fallback string) (json.RawMessage, error) {
	if len(value) == 0 {
		value = json.RawMessage(fallback)
	}
	var normalized any
	decoder := json.NewDecoder(bytes.NewReader(value))
	decoder.UseNumber()
	if err := decoder.Decode(&normalized); err != nil {
		return nil, err
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("trailing JSON content")
	}
	encoded, err := json.Marshal(normalized)
	return json.RawMessage(encoded), err
}

func Digest(value any) (string, error) {
	canonical, err := CanonicalJSON(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
