package preview

import (
	"encoding/json"
	"testing"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

func TestPreviewCellSelection(t *testing.T) {
	factors, err := parseFactors([]string{"collapse=unit"})
	if err != nil {
		t.Fatal(err)
	}
	cells := []ragcontract.ExpandedCell{{VariantID: "raw", Factors: []ragcontract.FactorSelection{{FactorID: "collapse", ValueID: "chunk", Value: json.RawMessage(`"chunk"`)}}}, {VariantID: "all", Factors: []ragcontract.FactorSelection{{FactorID: "collapse", ValueID: "unit", Value: json.RawMessage(`"unit"`)}}}}
	cell, err := selectCell(cells, "all", factors)
	if err != nil {
		t.Fatal(err)
	}
	if cell.VariantID != "all" {
		t.Fatalf("cell=%#v", cell)
	}
	if _, err := parseFactors([]string{"broken"}); err == nil {
		t.Fatal("accepted malformed factor")
	}
}
