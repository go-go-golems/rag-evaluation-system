package ragoperators

import (
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
)

type CorpusArtifact struct {
	Manifest ragcontract.CorpusManifest `json:"manifest"`
	Corpus   Corpus                     `json:"corpus"`
}

type EvaluationArtifact struct {
	Manifest ragcontract.EvaluationDatasetManifest `json:"manifest"`
	Dataset  EvaluationDataset                     `json:"dataset"`
}

func NewCorpusArtifact(corpus Corpus, sourceNamespace string) CorpusArtifact {
	digest, _ := ragcontract.Digest(corpus)
	return CorpusArtifact{Manifest: ragcontract.CorpusManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.CorpusManifestSchema, Digest: digest, Parents: []ragcontract.ParentDigest{}}, SourceNamespace: sourceNamespace, RecordSchema: "rag-source-record/v2", Ordering: "session-ordinal-id", RecordCount: int64(len(corpus.Records))}, Corpus: corpus}
}

func NewEvaluationArtifact(dataset EvaluationDataset, datasetID, split, status, relevanceTarget, corpusDigest string) EvaluationArtifact {
	digest, _ := ragcontract.Digest(dataset)
	return EvaluationArtifact{Manifest: ragcontract.EvaluationDatasetManifest{ManifestBase: ragcontract.ManifestBase{SchemaVersion: ragcontract.EvaluationManifestSchema, Digest: digest, Parents: []ragcontract.ParentDigest{{Role: "corpus", Digest: corpusDigest, SchemaVersion: ragcontract.CorpusManifestSchema}}, Production: &ragcontract.Production{Operator: ragcontract.OperatorRef{Kind: "evaluation.import", Version: "v1"}, Config: []byte(`{}`)}}, DatasetID: datasetID, Split: split, Status: status, RelevanceTarget: relevanceTarget, QueryCount: int64(len(dataset.Queries)), GradeSchema: []byte(`{"type":"number"}`)}, Dataset: dataset}
}

func ValidateInputArtifacts(execution ragcontract.PipelineExecution, corpus CorpusArtifact, evaluation EvaluationArtifact) error {
	if err := ragcontract.ValidateManifestBase(corpus.Manifest.ManifestBase, ragcontract.CorpusManifestSchema, false); err != nil {
		return err
	}
	if err := ragcontract.ValidateManifestBase(evaluation.Manifest.ManifestBase, ragcontract.EvaluationManifestSchema, true); err != nil {
		return err
	}
	corpusDigest, _ := ragcontract.Digest(corpus.Corpus)
	if corpusDigest != corpus.Manifest.Digest || int64(len(corpus.Corpus.Records)) != corpus.Manifest.RecordCount {
		return fmt.Errorf("RAG_INPUT_CORPUS_DIGEST")
	}
	evaluationDigest, _ := ragcontract.Digest(evaluation.Dataset)
	if evaluationDigest != evaluation.Manifest.Digest || int64(len(evaluation.Dataset.Queries)) != evaluation.Manifest.QueryCount {
		return fmt.Errorf("RAG_INPUT_EVALUATION_DIGEST")
	}
	boundCorpus := ""
	for _, binding := range execution.Bindings {
		if binding.Role == "corpus" {
			boundCorpus = binding.Digest
		}
	}
	if boundCorpus != corpus.Manifest.Digest {
		return fmt.Errorf("RAG_INPUT_CORPUS_BINDING: got %s want %s", boundCorpus, corpus.Manifest.Digest)
	}
	if execution.Dataset.ManifestDigest != evaluation.Manifest.Digest {
		return fmt.Errorf("RAG_INPUT_EVALUATION_BINDING: got %s want %s", execution.Dataset.ManifestDigest, evaluation.Manifest.Digest)
	}
	if len(evaluation.Manifest.Parents) != 1 || evaluation.Manifest.Parents[0].Digest != corpus.Manifest.Digest {
		return fmt.Errorf("RAG_INPUT_EVALUATION_LINEAGE")
	}
	if evaluation.Manifest.Split != execution.Dataset.Split || evaluation.Manifest.Status != execution.Dataset.Status || evaluation.Manifest.RelevanceTarget != execution.Dataset.RelevanceTarget {
		return fmt.Errorf("RAG_INPUT_EVALUATION_POLICY")
	}
	return nil
}
