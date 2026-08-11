package ragmodel

import (
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcompiler"
)

func ValidateProduct(value *Product) error {
	if value == nil || value.Pipeline == nil || value.Query == nil {
		return fmt.Errorf("RAG_V2_PRODUCT_INCOMPLETE: pipeline and query are required")
	}
	ir, err := BuildIR(value.Pipeline, value.Query, normalizedReranker(value.Reranker), value.Generator)
	if err != nil {
		return err
	}
	_, err = ragcompiler.Normalize(ir, nil)
	return err
}

func ValidateStudy(value *Study) error {
	if value == nil || value.Pipeline == nil || value.Dataset.Role == "" || len(value.Variants) == 0 {
		return fmt.Errorf("RAG_V2_STUDY_INCOMPLETE: pipeline, dataset, and variants are required")
	}
	if value.Replicates < 1 {
		return fmt.Errorf("RAG_V2_STUDY_REPLICATES: replicates must be positive")
	}
	for _, variant := range value.Variants {
		if variant.ID == "" || variant.Query == nil {
			return fmt.Errorf("RAG_V2_VARIANT_INCOMPLETE: id and query are required")
		}
		if _, err := BuildIR(value.Pipeline, variant.Query, nil, nil); err != nil {
			return fmt.Errorf("variant %s: %w", variant.ID, err)
		}
	}
	return nil
}
