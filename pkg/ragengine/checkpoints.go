package ragengine

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragcontract"
	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

const queryCheckpointSchemaVersion = "rag-query-checkpoint/v1"

type QueryCheckpointIdentity struct {
	SchemaVersion          string `json:"schemaVersion"`
	ExecutionDigest        string `json:"executionDigest"`
	QueryID                string `json:"queryId"`
	QueryTextDigest        string `json:"queryTextDigest"`
	PreparedCorpusDigest   string `json:"preparedCorpusDigest"`
	GeneratorFingerprint   string `json:"generatorFingerprint"`
	RerankerFingerprint    string `json:"rerankerFingerprint"`
	EvaluationPolicyDigest string `json:"evaluationPolicyDigest"`
}

type QueryCheckpoint struct {
	SchemaVersion string                  `json:"schemaVersion"`
	Identity      QueryCheckpointIdentity `json:"identity"`
	Trace         ragcontract.QueryTrace  `json:"trace"`
	Metrics       []ragoperators.Metric   `json:"metrics"`
	Artifacts     []ragoperators.Artifact `json:"artifacts"`
	Answer        *ragoperators.Answer    `json:"answer,omitempty"`
}

type QueryCheckpointStore interface {
	Get(context.Context, QueryCheckpointIdentity) (QueryCheckpoint, bool, error)
	Put(context.Context, QueryCheckpoint) error
}

// FileQueryCheckpointStore uses content-addressed keys and atomic file
// replacement. It is intentionally host-local; researchctl remains the
// authoritative custodian for the attempt observations emitted from a resume.
type FileQueryCheckpointStore struct{ directory string }

func NewFileQueryCheckpointStore(directory string) (*FileQueryCheckpointStore, error) {
	if directory == "" {
		return nil, os.ErrInvalid
	}
	if err := os.MkdirAll(directory, 0o700); err != nil {
		return nil, err
	}
	return &FileQueryCheckpointStore{directory: directory}, nil
}

func (s *FileQueryCheckpointStore) path(identity QueryCheckpointIdentity) (string, error) {
	canonical, err := ragcontract.CanonicalJSON(identity)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(canonical)
	return filepath.Join(s.directory, hex.EncodeToString(sum[:])+".json"), nil
}

func (s *FileQueryCheckpointStore) Get(ctx context.Context, identity QueryCheckpointIdentity) (QueryCheckpoint, bool, error) {
	if err := ctx.Err(); err != nil {
		return QueryCheckpoint{}, false, err
	}
	path, err := s.path(identity)
	if err != nil {
		return QueryCheckpoint{}, false, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return QueryCheckpoint{}, false, nil
	}
	if err != nil {
		return QueryCheckpoint{}, false, err
	}
	var checkpoint QueryCheckpoint
	if err := json.Unmarshal(data, &checkpoint); err != nil {
		return QueryCheckpoint{}, false, nil
	}
	if checkpoint.SchemaVersion != queryCheckpointSchemaVersion || !checkpointIdentityEqual(checkpoint.Identity, identity) {
		return QueryCheckpoint{}, false, nil
	}
	return checkpoint, true, nil
}

func (s *FileQueryCheckpointStore) Put(ctx context.Context, checkpoint QueryCheckpoint) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if checkpoint.SchemaVersion != queryCheckpointSchemaVersion {
		return fmt.Errorf("RAG_QUERY_CHECKPOINT_SCHEMA")
	}
	path, err := s.path(checkpoint.Identity)
	if err != nil {
		return err
	}
	data, err := ragcontract.CanonicalJSON(checkpoint)
	if err != nil {
		return err
	}
	tmp, err := os.CreateTemp(s.directory, ".checkpoint-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer func() { _ = tmp.Close(); _ = os.Remove(tmpName) }()
	if _, err := tmp.Write(data); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	return os.Rename(tmpName, path)
}

func checkpointIdentityEqual(left, right QueryCheckpointIdentity) bool {
	leftJSON, leftErr := ragcontract.CanonicalJSON(left)
	rightJSON, rightErr := ragcontract.CanonicalJSON(right)
	return leftErr == nil && rightErr == nil && string(leftJSON) == string(rightJSON)
}
