package ragworkflow

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/scraper/pkg/researchrunner"
)

var safeQueryScope = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,63}$`)

type DomainProjector struct{}

func (DomainProjector) Project(_ context.Context, input researchrunner.DomainProjectionInput) (researchrunner.DomainProjection, error) {
	output, ok := input.Outputs["result"]
	if !ok || output.SchemaVersion != ResultSchema {
		return researchrunner.DomainProjection{}, fmt.Errorf("RAG_WORKFLOW_PROJECTION_OUTPUT")
	}
	result, err := DecodeResult(output.Data)
	if err != nil {
		return researchrunner.DomainProjection{}, err
	}
	projection := researchrunner.DomainProjection{}
	for _, query := range result.Results {
		if !safeQueryScope.MatchString(query.QueryID) {
			return researchrunner.DomainProjection{}, fmt.Errorf("RAG_WORKFLOW_PROJECTION_QUERY")
		}
		queryDigest, _ := ragcontract.Digest(query.QueryID)
		for _, metric := range query.Metrics {
			metadata := map[string]any{"schemaVersion": "rag-workflow-metric-projection/v1", "queryIdDigest": queryDigest, "resultDigest": result.Digest}
			if len(metric.Metadata) > 0 {
				metadataDigest, digestErr := ragcontract.Digest(json.RawMessage(metric.Metadata))
				if digestErr != nil {
					return researchrunner.DomainProjection{}, digestErr
				}
				metadata["measureMetadataDigest"] = metadataDigest
			}
			metadataJSON, marshalErr := ragcontract.CanonicalJSON(metadata)
			if marshalErr != nil {
				return researchrunner.DomainProjection{}, marshalErr
			}
			projection.Metrics = append(projection.Metrics, researchrunner.Metric{Name: metric.Name, Scope: "rag.query." + query.QueryID, Value: metric.Value, NumericProjection: metric.Numeric, Unit: metric.Unit, Metadata: metadataJSON})
		}
		projectedTrace, sanitizeErr := privacySafeTrace(query.Trace)
		if sanitizeErr != nil {
			return researchrunner.DomainProjection{}, sanitizeErr
		}
		trace, marshalErr := ragcontract.CanonicalJSON(struct {
			SchemaVersion string                 `json:"schemaVersion"`
			QueryIDDigest string                 `json:"queryIdDigest"`
			ResultDigest  string                 `json:"resultDigest"`
			Trace         ragcontract.QueryTrace `json:"trace"`
		}{"rag-workflow-query-projection/v1", queryDigest, result.Digest, projectedTrace})
		if marshalErr != nil {
			return researchrunner.DomainProjection{}, marshalErr
		}
		projection.Traces = append(projection.Traces, researchrunner.Trace{Kind: "rag.query", Value: trace})
	}
	return projection, nil
}

func privacySafeTrace(trace ragcontract.QueryTrace) (ragcontract.QueryTrace, error) {
	body, err := ragcontract.CanonicalJSON(trace)
	if err != nil {
		return ragcontract.QueryTrace{}, err
	}
	var safe ragcontract.QueryTrace
	if err := json.Unmarshal(body, &safe); err != nil {
		return ragcontract.QueryTrace{}, err
	}
	safe.Query.Metadata = nil
	for channelIndex := range safe.Channels {
		for hitIndex := range safe.Channels[channelIndex].Hits {
			safe.Channels[channelIndex].Hits[hitIndex].Filter = nil
		}
	}
	for index := range safe.Failures {
		safe.Failures[index].Message = "redacted"
		safe.Failures[index].Details = nil
	}
	return safe, nil
}
