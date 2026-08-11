package researchctladapter_test

import (
	"context"
	"path/filepath"
	"testing"

	studycmd "github.com/go-go-golems/rag-evaluation-system/cmd/rag-eval/cmds/study"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/researchctladapter"
)

func TestRagSol2FiveVariantsByTwoCollapseModes(t *testing.T) {
	study, pathErr := studycmd.LoadStudy(filepath.Join("..", "..", "experiments", "rag-sol2", "study.js"))
	if pathErr != nil {
		t.Fatal(pathErr)
	}
	root := t.TempDir()
	corpus := ragoperators.NewCorpusArtifact(ragoperators.Corpus{SchemaVersion: "rag-source-record-set/v2", Records: []ragoperators.SourceRecord{{ID: "u1", SessionID: "s", Ordinal: 1, Role: "user", Text: "Explain vector fusion."}, {ID: "a2", SessionID: "s", Ordinal: 2, Role: "assistant", Text: "The build failed. We decided to fix vector mapping."}, {ID: "a3", SessionID: "s", Ordinal: 3, Role: "assistant", Text: "Then run verification."}}}, "fixture")
	dataset := ragoperators.NewEvaluationArtifact(ragoperators.EvaluationDataset{SchemaVersion: "rag-evaluation-queries/v2", Queries: []ragoperators.Query{{ID: "q", Text: "Which vector mapping decision was made?"}}}, "fixture", "candidate", "candidate", "unit", corpus.Manifest.Digest)
	corpusInput, err := researchctladapter.StageEnvelope(researchctladapter.InputReference{Role: "corpus"}, corpus, root)
	if err != nil {
		t.Fatal(err)
	}
	evaluationInput, err := researchctladapter.StageEnvelope(researchctladapter.InputReference{Role: "evaluation-dataset"}, dataset, root)
	if err != nil {
		t.Fatal(err)
	}
	resolved := researchctladapter.ResolvedInputs{ByRole: map[string]researchctladapter.ResolvedInput{"corpus": corpusInput, "evaluation-dataset": evaluationInput}}
	_, cells, err := researchctladapter.Expand(study, resolved)
	if err != nil {
		t.Fatal(err)
	}
	if len(cells) != 10 {
		t.Fatalf("cells=%d", len(cells))
	}
	ids := map[string]bool{}
	variants := map[string]int{}
	fixtures := ragoperators.NewFixtureProviders()
	for _, cell := range cells {
		if ids[cell.ID] {
			t.Fatalf("duplicate cell %s", cell.ID)
		}
		ids[cell.ID] = true
		variants[cell.VariantID]++
		result, err := ragengine.New(nil).Execute(context.Background(), cell.Execution, corpus.Corpus, dataset.Dataset, ragengine.NopObserver{}, ragengine.Options{Manifests: fixtures.Resolver, Schemas: fixtures, Generator: fixtures, Embedder: fixtures, Cache: ragoperators.NewMemoryCache()})
		if err != nil {
			t.Fatalf("cell %s: %v", cell.ID, err)
		}
		if len(result.Traces) != 1 {
			t.Fatalf("traces=%d", len(result.Traces))
		}
		trace := result.Traces[0]
		kinds := variantKinds(cell.VariantID)
		if len(trace.Channels) != len(kinds)*2 {
			t.Fatalf("%s channels=%d", cell.VariantID, len(trace.Channels))
		}
		for _, collapse := range trace.Collapses {
			seen := map[string]bool{}
			for _, group := range collapse.Groups {
				if seen[group.Key.ID] {
					t.Fatalf("duplicate collapse vote %s/%s", collapse.Channel, group.Key.ID)
				}
				seen[group.Key.ID] = true
			}
		}
		if trace.Fusion == nil {
			t.Fatal("missing fusion trace")
		}
		for _, fused := range trace.Fusion.Results {
			seen := map[string]bool{}
			for _, contribution := range fused.Contributions {
				if seen[contribution.Channel] {
					t.Fatalf("duplicate channel contribution %s", contribution.Channel)
				}
				seen[contribution.Channel] = true
			}
		}
		for _, hit := range trace.Results {
			if hit.Evidence.Citation.SourceID == "" || len(hit.MatchedRepresentations) == 0 {
				t.Fatalf("unhydrated result=%#v", hit)
			}
		}
		for _, cost := range trace.Usage.ProviderCost {
			if cost != 0 {
				t.Fatalf("fixture cost=%f", cost)
			}
		}
	}
	for _, id := range []string{"raw", "summary", "raw-summary", "raw-question", "all"} {
		if variants[id] != 2 {
			t.Fatalf("variant %s cells=%d", id, variants[id])
		}
	}
}
func variantKinds(id string) []string {
	switch id {
	case "raw":
		return []string{"raw"}
	case "summary":
		return []string{"summary"}
	case "raw-summary":
		return []string{"raw", "summary"}
	case "raw-question":
		return []string{"raw", "question"}
	case "all":
		return []string{"raw", "summary", "question"}
	default:
		return nil
	}
}
