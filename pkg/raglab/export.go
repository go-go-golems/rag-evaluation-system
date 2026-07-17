package raglab

import (
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type ExportOptions struct {
	DatasetSplit string
}

// ExportSpecificationV1 converts the prototype authoring value into the pure
// researchctl RAG wire payload. Artifact IDs and prototype fingerprints are
// intentionally excluded: the researchctl adapter represents those as input
// references and provenance around this domain value.
func ExportSpecificationV1(input ExperimentSpecification, options ExportOptions) (PrototypeSpecification, error) {
	if strings.TrimSpace(options.DatasetSplit) == "" {
		return PrototypeSpecification{}, errors.New("RAG_DATASET_SPLIT_REQUIRED: export requires an explicit dataset split")
	}
	if len(input.Metrics.HitRateAt) > 0 || input.Metrics.MeanRelevantRecall > 0 || input.Metrics.Abstention {
		return PrototypeSpecification{}, errors.New("RAG_EXPORT_UNSUPPORTED: hit-rate, legacy mean-relevant-recall, and abstention metrics are not in rag-retrieval-spec/v1")
	}
	if input.Metrics.RelevanceAt == nil {
		return PrototypeSpecification{}, errors.New("RAG_RELEVANCE_GRADE_REQUIRED: export requires an explicit relevance threshold")
	}
	result := PrototypeSpecification{
		SchemaVersion: PrototypeSchemaVersion,
		Name:          input.Name,
		Dataset:       PrototypeDatasetSelection{Split: options.DatasetSplit},
		Metrics: PrototypeMetricPlan{
			RelevanceAt: input.Metrics.RelevanceAt.Name,
			RecallAt:    normalizedExportInts(input.Metrics.RecallAt),
			PrecisionAt: normalizedExportInts(input.Metrics.PrecisionAt),
			MRR:         input.Metrics.MRR,
		},
		Tags: cloneStringMap(input.Provenance.Tags),
	}
	if input.Metrics.NDCGAt > 0 {
		result.Metrics.NDCGAt = []int{input.Metrics.NDCGAt}
	}
	for _, representation := range input.Inputs.Representations {
		kind := string(representation.Kind)
		if representation.Kind == QuestionRepresentation {
			kind = "questions"
		}
		result.Representations = append(result.Representations, PrototypeRepresentation{Name: representation.Name, Kind: kind})
	}
	result.Retrieval = PrototypeRetrievalPlan{
		Filter:   exportFilter(input.Retrieval.Filter),
		Collapse: string(input.Retrieval.Collapse),
		Results:  input.Retrieval.Results,
	}
	for _, channel := range input.Retrieval.Channels {
		result.Retrieval.Channels = append(result.Retrieval.Channels, PrototypeChannel{
			Name: channel.Name, Backend: string(channel.Backend), Representation: channel.Representation,
			TopK: channel.TopK, Filter: exportFilter(channel.Filter),
		})
	}
	if input.Retrieval.Fusion != nil {
		result.Retrieval.Fusion = &PrototypeFusionPlan{
			Kind: input.Retrieval.Fusion.Kind, RankConstant: input.Retrieval.Fusion.RankConstant,
			Weights: cloneFloatMap(input.Retrieval.Fusion.Weights),
		}
	}
	if input.Retrieval.Reranking != nil {
		result.Retrieval.Reranking = &PrototypeRerankingPlan{
			Kind: string(input.Retrieval.Reranking.Kind), Model: input.Retrieval.Reranking.Model,
			CandidateCount: input.Retrieval.Reranking.CandidateCount, Results: input.Retrieval.Reranking.Results,
		}
	}
	return result, nil
}

func exportFilter(input FilterSpec) PrototypeFilterSpec {
	return PrototypeFilterSpec{
		SourceIDs: normalizedExportStrings(input.SourceIDs), DocumentIDs: normalizedExportStrings(input.DocumentIDs),
		ContentTypes: normalizedExportStrings(input.ContentTypes), MetadataEquals: cloneStringMap(input.MetadataEquals),
	}
}

func normalizedExportInts(input []int) []int {
	seen := map[int]bool{}
	for _, value := range input {
		seen[value] = true
	}
	result := make([]int, 0, len(seen))
	for value := range seen {
		result = append(result, value)
	}
	sort.Ints(result)
	if len(result) == 0 {
		return nil
	}
	return result
}

func normalizedExportStrings(input []string) []string {
	seen := map[string]bool{}
	for _, value := range input {
		seen[value] = true
	}
	result := make([]string, 0, len(seen))
	for value := range seen {
		result = append(result, value)
	}
	sort.Strings(result)
	if len(result) == 0 {
		return nil
	}
	return result
}

func cloneStringMap(input map[string]string) map[string]string {
	if len(input) == 0 {
		return nil
	}
	result := make(map[string]string, len(input))
	for key, value := range input {
		result[key] = value
	}
	return result
}

func cloneFloatMap(input map[string]float64) map[string]float64 {
	if len(input) == 0 {
		return nil
	}
	result := make(map[string]float64, len(input))
	for key, value := range input {
		result[key] = value
	}
	return result
}
