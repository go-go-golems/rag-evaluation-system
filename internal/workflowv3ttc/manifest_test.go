package workflowv3ttc

import (
	"context"
	"testing"

	workflowmodule "github.com/go-go-golems/scraper/pkg/gojamodules/workflow"
)

func TestRunManifestFreezesPlanRegistryArtifactsAndRepositories(t *testing.T) {
	registry, err := Registry()
	if err != nil {
		t.Fatal(err)
	}
	catalog, _ := registry.Catalog()
	authored, err := workflowmodule.Author(context.Background(), ProductionWorkflowSource(), catalog, DescriptorModule())
	if err != nil {
		t.Fatal(err)
	}
	input := RunManifest{RunID: "ttc-v3-fixture", Corpus: FrozenArtifactIdentity{ID: "corpus", SchemaVersion: "corpus/v1", Digest: digestOf("a"), SizeBytes: 10, ItemCount: 1807}, EvaluationDataset: FrozenArtifactIdentity{ID: "evaluation", SchemaVersion: "evaluation/v1", Digest: digestOf("b"), SizeBytes: 10, ItemCount: 25}, ProviderProfile: FrozenIdentity{Kind: "provider-profile", ID: "fixture", Digest: digestOf("c")}, Models: []FrozenIdentity{{Kind: "model", ID: "generator", Digest: digestOf("d")}}, Prompts: []FrozenIdentity{{Kind: "prompt", ID: "generation", Digest: digestOf("e")}}, Schemas: []FrozenIdentity{{Kind: "schema", ID: "combined", Digest: digestOf("f")}}, RepositoryCommits: map[string]string{"rag": "be7178b", "scraper": "106e8d0"}, RuntimeVersions: map[string]string{"go": "fixture"}}
	first, err := NewRunManifest(input, authored.Plan, registry)
	if err != nil {
		t.Fatal(err)
	}
	second, err := NewRunManifest(input, authored.Plan, registry)
	if err != nil {
		t.Fatal(err)
	}
	if first.ManifestDigest == "" || first.ManifestDigest != second.ManifestDigest || first.WorkflowPlanDigest != authored.Plan.Digest || first.RegistryGeneration != registry.Generation() || len(first.BundleDigests) != 1 {
		t.Fatalf("manifest=%#v", first)
	}
}
