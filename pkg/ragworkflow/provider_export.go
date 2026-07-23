package ragworkflow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragworkflowops"
	"github.com/go-go-golems/scraper/pkg/researchrunner"
)

func WriteDeterministicProviderFixtures(ctx context.Context, directory string) (FixtureManifest, error) {
	if directory == "" {
		return FixtureManifest{}, fmt.Errorf("RAG_WORKFLOW_FIXTURE_DIRECTORY")
	}
	if err := os.MkdirAll(directory, 0755); err != nil {
		return FixtureManifest{}, err
	}
	services, err := NewDeterministicProviderServices()
	if err != nil {
		return FixtureManifest{}, err
	}
	providerPackage, err := NewProviderPackage(services, defaultProviderPolicy())
	if err != nil {
		return FixtureManifest{}, err
	}
	lowerer, err := NewProviderLowerer(providerPackage)
	if err != nil {
		return FixtureManifest{}, err
	}
	manifest := FixtureManifest{SchemaVersion: "rag-workflow-fixture-manifest/v1"}
	for index, id := range []string{"case-a", "case-b"} {
		fixture, fixtureErr := NewDeterministicProviderFixture()
		if fixtureErr != nil {
			return manifest, fixtureErr
		}
		if index == 1 {
			fixture.Corpus.Records[2].Text = "second irrelevant provider fixture"
		}
		corpusCanonical, _ := ragcontract.CanonicalJSON(fixture.Corpus)
		corpusDigest, _ := ragcontract.Digest(fixture.Corpus)
		size := int64(len(corpusCanonical) + 1)
		fixture.Execution.Bindings[0].Digest, fixture.Execution.Bindings[0].SizeBytes = corpusDigest, &size
		fixture.Execution.CellID = ""
		fixture.Execution.CellID, _ = ragcontract.Digest(fixture.Execution)
		lowered, lowerErr := lowerer.Lower(ctx, fixture.Execution)
		if lowerErr != nil {
			return manifest, lowerErr
		}
		bindings := map[string]researchrunner.InputBinding{"execution": {Role: "workflow-input", Kind: "rag-execution", ID: "execution-" + id}, "corpus": {Role: "workflow-input", Kind: "rag-corpus", ID: "corpus-" + id}, "queries": {Role: "workflow-input", Kind: "rag-query-set", ID: "queries-" + id}}
		domain, domainErr := BuildProviderRunnerExecution(lowered, bindings, providerPackage)
		if domainErr != nil {
			return manifest, domainErr
		}
		archive, archiveErr := BuildQueryArchive(fixture.Execution, fixture.Dataset)
		if archiveErr != nil {
			return manifest, archiveErr
		}
		parity, parityErr := buildProviderParity(ctx, fixture, services, providerPackage.authority.Digest)
		if parityErr != nil {
			return manifest, parityErr
		}
		caseDir := filepath.Join(directory, id)
		if mkdirErr := os.MkdirAll(caseDir, 0755); mkdirErr != nil {
			return manifest, mkdirErr
		}
		current := FixtureCase{ID: id}
		values := map[string]struct {
			schema string
			value  any
			target *FixtureFile
		}{"execution.json": {ragcontract.ExecutionSchemaVersion, fixture.Execution, &current.Execution}, "corpus.json": {CorpusSchema, fixture.Corpus, &current.Corpus}, "queries.json": {researchrunner.SetInputArchiveSchema, archive, &current.Queries}, "domain-config.json": {researchrunner.DomainSchemaVersion, domain, &current.DomainConfig}, "ragengine-parity.json": {"rag-workflow-parity/v1", parity, &current.Parity}}
		for name, value := range values {
			body, marshalErr := ragcontract.CanonicalJSON(value.value)
			if marshalErr != nil {
				return manifest, marshalErr
			}
			body = append(body, '\n')
			path := filepath.Join(caseDir, name)
			if writeErr := os.WriteFile(path, body, 0644); writeErr != nil {
				return manifest, writeErr
			}
			sum := sha256.Sum256(body)
			*value.target = FixtureFile{Path: filepath.ToSlash(filepath.Join(id, name)), SchemaVersion: value.schema, Digest: "sha256:" + hex.EncodeToString(sum[:]), SizeBytes: int64(len(body))}
		}
		manifest.Cases = append(manifest.Cases, current)
	}
	body, err := ragcontract.CanonicalJSON(manifest)
	if err != nil {
		return manifest, err
	}
	if err = os.WriteFile(filepath.Join(directory, "manifest.json"), append(body, '\n'), 0644); err != nil {
		return manifest, err
	}
	return manifest, nil
}

func defaultProviderPolicy() ragworkflowops.Policy {
	return ragworkflowops.Policy{MaxPerAttempt: 10_000, FinishTimeout: 5 * time.Second}
}

func buildProviderParity(ctx context.Context, fixture Fixture, services ProviderServices, fingerprint string) (ParityFixture, error) {
	result := ParityFixture{SchemaVersion: "rag-workflow-parity/v1"}
	result.ExecutionDigest, _ = ragcontract.Digest(fixture.Execution)
	for _, query := range fixture.Dataset.Queries {
		run, err := ragengine.New(nil).Execute(ctx, fixture.Execution, fixture.Corpus, ragoperators.EvaluationDataset{SchemaVersion: fixture.Dataset.SchemaVersion, Queries: []ragoperators.Query{query}}, nil, ragengine.Options{Manifests: services.Manifests, Schemas: services.Schemas, Generator: services.Generator, Embedder: services.Embedder, Reranker: services.Reranker, GenerationConcurrency: services.GenerationConcurrency, GenerationSettingsFingerprint: fingerprint, EmbeddingFingerprint: fingerprint, GeneratorFingerprint: fingerprint, RerankerFingerprint: fingerprint})
		if err != nil {
			return result, err
		}
		ragengine.SortMetrics(run.Metrics)
		result.Queries = append(result.Queries, ParityQuery{QueryID: query.ID, Results: run.Traces[0].Results, Metrics: workflowMetrics(run.Metrics)})
	}
	return result, nil
}
