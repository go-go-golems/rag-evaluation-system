package ragworkflow

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragengine"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
	"github.com/go-go-golems/scraper/pkg/researchrunner"
)

type FixtureFile struct {
	Path          string `json:"path"`
	SchemaVersion string `json:"schemaVersion"`
	Digest        string `json:"digest"`
	SizeBytes     int64  `json:"sizeBytes"`
}
type FixtureCase struct {
	ID           string      `json:"id"`
	Execution    FixtureFile `json:"execution"`
	Corpus       FixtureFile `json:"corpus"`
	Queries      FixtureFile `json:"queries"`
	DomainConfig FixtureFile `json:"domainConfig"`
	Parity       FixtureFile `json:"parity"`
}
type FixtureManifest struct {
	SchemaVersion string        `json:"schemaVersion"`
	Cases         []FixtureCase `json:"cases"`
}
type ParityQuery struct {
	QueryID string                    `json:"queryId"`
	Results []ragcontract.ResultTrace `json:"results"`
	Metrics []Metric                  `json:"metrics"`
}
type ParityFixture struct {
	SchemaVersion   string        `json:"schemaVersion"`
	ExecutionDigest string        `json:"executionDigest"`
	Queries         []ParityQuery `json:"queries"`
}

func WriteProviderFreeFixtures(ctx context.Context, directory string) (FixtureManifest, error) {
	if directory == "" {
		return FixtureManifest{}, fmt.Errorf("RAG_WORKFLOW_FIXTURE_DIRECTORY")
	}
	if err := os.MkdirAll(directory, 0755); err != nil {
		return FixtureManifest{}, err
	}
	manifest := FixtureManifest{SchemaVersion: "rag-workflow-fixture-manifest/v1"}
	for index, id := range []string{"case-a", "case-b"} {
		fixture, err := NewProviderFreeFixture(true)
		if err != nil {
			return manifest, err
		}
		if index == 1 {
			fixture.Corpus.Records[2].Text = "second irrelevant fixture"
		}
		corpusCanonical, _ := ragcontract.CanonicalJSON(fixture.Corpus)
		corpusDigest, _ := ragcontract.Digest(fixture.Corpus)
		fixture.Execution.Bindings[0].Digest = corpusDigest
		size := int64(len(corpusCanonical) + 1)
		fixture.Execution.Bindings[0].SizeBytes = &size
		fixture.Execution.CellID = ""
		fixture.Execution.CellID, _ = ragcontract.Digest(fixture.Execution)
		lowered, err := NewLowerer().Lower(ctx, fixture.Execution)
		if err != nil {
			return manifest, err
		}
		bindings := map[string]researchrunner.InputBinding{
			"execution": {Role: "workflow-input", Kind: "rag-execution", ID: "execution-" + id},
			"corpus":    {Role: "workflow-input", Kind: "rag-corpus", ID: "corpus-" + id},
			"queries":   {Role: "workflow-input", Kind: "rag-query-set", ID: "queries-" + id},
		}
		domain, err := BuildRunnerExecution(lowered, bindings)
		if err != nil {
			return manifest, err
		}
		archive, err := BuildQueryArchive(fixture.Execution, fixture.Dataset)
		if err != nil {
			return manifest, err
		}
		parity, err := buildParity(ctx, fixture)
		if err != nil {
			return manifest, err
		}
		caseDir := filepath.Join(directory, id)
		if err := os.MkdirAll(caseDir, 0755); err != nil {
			return manifest, err
		}
		current := FixtureCase{ID: id}
		for name, value := range map[string]struct {
			schema string
			value  any
			target *FixtureFile
		}{
			"execution.json":        {ragcontract.ExecutionSchemaVersion, fixture.Execution, &current.Execution},
			"corpus.json":           {CorpusSchema, fixture.Corpus, &current.Corpus},
			"queries.json":          {researchrunner.SetInputArchiveSchema, archive, &current.Queries},
			"domain-config.json":    {researchrunner.DomainSchemaVersion, domain, &current.DomainConfig},
			"ragengine-parity.json": {"rag-workflow-parity/v1", parity, &current.Parity},
		} {
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
	if err := os.WriteFile(filepath.Join(directory, "manifest.json"), append(body, '\n'), 0644); err != nil {
		return manifest, err
	}
	return manifest, nil
}

func buildParity(ctx context.Context, fixture Fixture) (ParityFixture, error) {
	env, err := providerFreeEnvironment(fixture.Execution)
	if err != nil {
		return ParityFixture{}, err
	}
	result := ParityFixture{SchemaVersion: "rag-workflow-parity/v1"}
	result.ExecutionDigest, _ = ragcontract.Digest(fixture.Execution)
	for _, query := range fixture.Dataset.Queries {
		run, runErr := ragengine.New(nil).Execute(ctx, fixture.Execution, fixture.Corpus, ragoperators.EvaluationDataset{SchemaVersion: fixture.Dataset.SchemaVersion, Queries: []ragoperators.Query{query}}, nil, ragengine.Options{Manifests: env.Manifests, Embedder: env.Embedder, EmbeddingFingerprint: "fixture-embedding/v1"})
		if runErr != nil {
			return result, runErr
		}
		ragengine.SortMetrics(run.Metrics)
		result.Queries = append(result.Queries, ParityQuery{QueryID: query.ID, Results: run.Traces[0].Results, Metrics: workflowMetrics(run.Metrics)})
	}
	return result, nil
}
