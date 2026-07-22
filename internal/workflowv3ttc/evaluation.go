package workflowv3ttc

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

const (
	QuerySchema         = "rag-ttc-query/v1"
	QueryEvidenceSchema = "rag-ttc-query-evidence/v1"
)

type QueryEnvelope struct {
	SchemaVersion string             `json:"schemaVersion"`
	DatasetDigest string             `json:"datasetDigest"`
	Query         ragoperators.Query `json:"query"`
}

type MetricEvidence struct {
	Name    string   `json:"name"`
	Unit    string   `json:"unit"`
	Value   []byte   `json:"value"`
	Numeric *float64 `json:"numeric,omitempty"`
}

type QueryEvidence struct {
	SchemaVersion    string           `json:"schemaVersion"`
	QueryID          string           `json:"queryId"`
	DatasetDigest    string           `json:"datasetDigest"`
	CitationChunkIDs []string         `json:"citationChunkIds"`
	Abstained        bool             `json:"abstained"`
	InputTokens      int64            `json:"inputTokens"`
	OutputTokens     int64            `json:"outputTokens"`
	Metrics          []MetricEvidence `json:"metrics"`
	FailureCodes     []string         `json:"failureCodes,omitempty"`
	Usage            []Usage          `json:"usage"`
}

type EvaluationConfig struct {
	Store            ragengine.PreparedCorpusStore
	Engine           *ragengine.Engine
	Execution        ragcontract.PipelineExecution
	Corpus           ragoperators.Corpus
	Options          ragengine.Options
	Identity         ragengine.PreparedCorpusIdentity
	DatasetDigest    string
	ValidCitationIDs map[string]struct{}
}

type EvaluationService interface {
	Evaluate(context.Context, PublicationReceipt, QueryEnvelope) (QueryEvidence, error)
}

type EvaluationAuthority struct{ config EvaluationConfig }

func NewEvaluationAuthority(config EvaluationConfig) (*EvaluationAuthority, error) {
	if config.Store == nil || config.Engine == nil || config.Execution.SchemaVersion != ragcontract.ExecutionSchemaVersion || config.DatasetDigest == "" || len(config.ValidCitationIDs) == 0 {
		return nil, fmt.Errorf("complete immutable TTC evaluation configuration is required")
	}
	return &EvaluationAuthority{config: config}, nil
}

func (a *EvaluationAuthority) Evaluate(ctx context.Context, publication PublicationReceipt, query QueryEnvelope) (QueryEvidence, error) {
	if publication.SchemaVersion != PublicationReceiptSchema || publication.Identity != a.config.Identity || publication.PreparedDigest == "" ||
		query.SchemaVersion != QuerySchema || query.DatasetDigest != a.config.DatasetDigest || query.Query.ID == "" || query.Query.Text == "" {
		return QueryEvidence{}, fmt.Errorf("RAG_TTC_EVALUATION_IDENTITY")
	}
	prepared, found, err := a.config.Store.Open(ctx, a.config.Engine, a.config.Execution.Pipeline, a.config.Corpus, a.config.Options, a.config.Identity)
	if err != nil || !found {
		return QueryEvidence{}, fmt.Errorf("RAG_TTC_EVALUATION_REOPEN")
	}
	defer func() { _ = prepared.Close() }()
	options := a.config.Options
	options.Prepared = prepared
	result, err := a.config.Engine.Execute(ctx, a.config.Execution, a.config.Corpus, ragoperators.EvaluationDataset{SchemaVersion: "rag-evaluation-dataset/v1", Queries: []ragoperators.Query{query.Query}}, nil, options)
	if err != nil {
		return QueryEvidence{}, err
	}
	evidence := QueryEvidence{SchemaVersion: QueryEvidenceSchema, QueryID: query.Query.ID, DatasetDigest: query.DatasetDigest}
	if len(result.Answers) > 0 {
		answer := result.Answers[0]
		evidence.CitationChunkIDs = append([]string(nil), answer.CitationChunkIDs...)
		sort.Strings(evidence.CitationChunkIDs)
		for index, citation := range evidence.CitationChunkIDs {
			if _, ok := a.config.ValidCitationIDs[citation]; !ok || (index > 0 && citation == evidence.CitationChunkIDs[index-1]) {
				return QueryEvidence{}, fmt.Errorf("RAG_TTC_CITATION_INVALID")
			}
		}
		evidence.Abstained, evidence.InputTokens, evidence.OutputTokens = answer.Abstained, answer.InputTokens, answer.OutputTokens
	}
	costMicrounits, embeddingTokens := int64(0), int64(0)
	if len(result.Metrics) > 64 {
		return QueryEvidence{}, fmt.Errorf("RAG_TTC_METRICS_BOUNDED")
	}
	for _, metric := range result.Metrics {
		if metric.Name == "" || len(metric.Value) > 64<<10 {
			return QueryEvidence{}, fmt.Errorf("RAG_TTC_METRICS_BOUNDED")
		}
		evidence.Metrics = append(evidence.Metrics, MetricEvidence{Name: metric.Name, Unit: metric.Unit, Value: append([]byte(nil), metric.Value...), Numeric: metric.Numeric})
		switch metric.Name {
		case "rag.token-usage":
			var usage map[string]int64
			if err := json.Unmarshal(metric.Value, &usage); err == nil {
				embeddingTokens += usage["embedding"]
			}
		case "rag.provider-cost":
			var costs map[string]float64
			if err := json.Unmarshal(metric.Value, &costs); err == nil {
				for _, cost := range costs {
					costMicrounits += int64(math.Round(cost * 1_000_000))
				}
			}
		}
	}
	for _, failure := range result.Failures {
		evidence.FailureCodes = append(evidence.FailureCodes, failure.Code)
	}
	evidence.Usage = []Usage{{Dimension: "cost_microunits", Units: costMicrounits}, {Dimension: "embedding_tokens", Units: embeddingTokens}, {Dimension: "input_tokens", Units: evidence.InputTokens}, {Dimension: "output_tokens", Units: evidence.OutputTokens}, {Dimension: "requests", Units: 2}}
	return evidence, nil
}
