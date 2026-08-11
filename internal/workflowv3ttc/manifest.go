package workflowv3ttc

import (
	"fmt"
	"sort"
	"strings"

	"github.com/go-go-golems/scraper/pkg/workflowv3"
)

const RunManifestSchema = "rag-ttc-workflow-v3-run-manifest/v1"

type FrozenArtifactIdentity struct {
	ID            string `json:"id"`
	SchemaVersion string `json:"schemaVersion"`
	Digest        string `json:"digest"`
	SizeBytes     int64  `json:"sizeBytes"`
	ItemCount     int    `json:"itemCount"`
}

type FrozenIdentity struct {
	Kind   string `json:"kind"`
	ID     string `json:"id"`
	Digest string `json:"digest"`
}

type RunManifest struct {
	SchemaVersion      string                  `json:"schemaVersion"`
	RunID              string                  `json:"runId"`
	Corpus             FrozenArtifactIdentity  `json:"corpus"`
	EvaluationDataset  FrozenArtifactIdentity  `json:"evaluationDataset"`
	ProviderProfile    FrozenIdentity          `json:"providerProfile"`
	Models             []FrozenIdentity        `json:"models"`
	Prompts            []FrozenIdentity        `json:"prompts"`
	Schemas            []FrozenIdentity        `json:"schemas"`
	WorkflowPlan       workflowv3.WorkflowPlan `json:"workflowPlan"`
	WorkflowPlanDigest string                  `json:"workflowPlanDigest"`
	RegistryGeneration string                  `json:"registryGeneration"`
	BundleDigests      []string                `json:"bundleDigests"`
	RepositoryCommits  map[string]string       `json:"repositoryCommits"`
	RuntimeVersions    map[string]string       `json:"runtimeVersions"`
	ManifestDigest     string                  `json:"manifestDigest"`
}

func NewRunManifest(manifest RunManifest, plan workflowv3.WorkflowPlan, registry *workflowv3.SealedRegistry) (RunManifest, error) {
	if registry == nil || manifest.RunID == "" || manifest.Corpus.ItemCount < 1 || manifest.EvaluationDataset.ItemCount < 1 ||
		manifest.Corpus.SizeBytes < 1 || manifest.EvaluationDataset.SizeBytes < 1 || manifest.ProviderProfile.Digest == "" ||
		len(manifest.Models) == 0 || len(manifest.Prompts) == 0 || len(manifest.Schemas) == 0 || len(manifest.RepositoryCommits) < 2 {
		return RunManifest{}, fmt.Errorf("complete frozen TTC identity is required")
	}
	for _, digest := range append([]string{manifest.Corpus.Digest, manifest.EvaluationDataset.Digest, manifest.ProviderProfile.Digest, plan.Digest, registry.Generation()}, identityDigests(manifest)...) {
		if !validDigest(digest) {
			return RunManifest{}, fmt.Errorf("TTC identity digest is invalid")
		}
	}
	manifest.SchemaVersion = RunManifestSchema
	manifest.WorkflowPlan = plan
	manifest.WorkflowPlanDigest = plan.Digest
	manifest.RegistryGeneration = registry.Generation()
	bundles := map[string]struct{}{}
	for _, node := range plan.Nodes {
		bundles[node.Implementation.BundleDigest] = struct{}{}
	}
	for _, mapped := range plan.Maps {
		bundles[mapped.Implementation.BundleDigest] = struct{}{}
	}
	for _, reduced := range plan.Reductions {
		bundles[reduced.Implementation.BundleDigest] = struct{}{}
	}
	manifest.BundleDigests = manifest.BundleDigests[:0]
	for digest := range bundles {
		manifest.BundleDigests = append(manifest.BundleDigests, digest)
	}
	sort.Strings(manifest.BundleDigests)
	sort.Slice(manifest.Models, func(i, j int) bool { return manifest.Models[i].ID < manifest.Models[j].ID })
	sort.Slice(manifest.Prompts, func(i, j int) bool { return manifest.Prompts[i].ID < manifest.Prompts[j].ID })
	sort.Slice(manifest.Schemas, func(i, j int) bool { return manifest.Schemas[i].ID < manifest.Schemas[j].ID })
	manifest.ManifestDigest = ""
	digest, err := workflowv3.Digest(manifest)
	if err != nil {
		return RunManifest{}, err
	}
	manifest.ManifestDigest = digest
	return manifest, nil
}

func identityDigests(manifest RunManifest) []string {
	ret := make([]string, 0, len(manifest.Models)+len(manifest.Prompts)+len(manifest.Schemas))
	for _, identities := range [][]FrozenIdentity{manifest.Models, manifest.Prompts, manifest.Schemas} {
		for _, identity := range identities {
			ret = append(ret, identity.Digest)
		}
	}
	return ret
}

func validDigest(value string) bool {
	if !strings.HasPrefix(value, "sha256:") || len(value) != len("sha256:")+64 {
		return false
	}
	for _, r := range value[len("sha256:"):] {
		if !strings.ContainsRune("0123456789abcdef", r) {
			return false
		}
	}
	return true
}
