package raglab

import (
	"context"
	"database/sql"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/pkg/errors"
)

// SQLiteCatalog exposes immutable artifact lineage from the rag-eval database.
type SQLiteCatalog struct{ queries *db.Queries }

var _ ArtifactCatalog = (*SQLiteCatalog)(nil)

func NewSQLiteCatalog(queries *db.Queries) *SQLiteCatalog { return &SQLiteCatalog{queries: queries} }

func (c *SQLiteCatalog) LookupArtifact(ctx context.Context, ref ArtifactRef) (ArtifactMetadata, error) {
	if c == nil || c.queries == nil {
		return ArtifactMetadata{}, errors.New("RAG_CATALOG_FAILURE: database queries are required")
	}
	if ref.ID == "" {
		return ArtifactMetadata{}, ErrArtifactNotFound
	}
	metadata := ArtifactMetadata{Ref: ref}
	var err error
	switch ref.Kind {
	case CorpusSnapshotArtifact:
		err = c.queries.DB().QueryRowContext(ctx, `SELECT id FROM corpus_snapshots WHERE id=?`, ref.ID).Scan(&metadata.Ref.ID)
	case ChunkSetArtifact:
		err = c.queries.DB().QueryRowContext(ctx, `SELECT id,corpus_snapshot_id FROM chunk_sets WHERE id=?`, ref.ID).Scan(&metadata.Ref.ID, &metadata.CorpusSnapshotID)
	case BM25IndexArtifact:
		err = c.queries.DB().QueryRowContext(ctx, `SELECT id,chunk_set_id FROM retrieval_artifacts WHERE id=? AND kind='bm25'`, ref.ID).Scan(&metadata.Ref.ID, &metadata.ChunkSetID)
	case EmbeddingSetArtifact:
		err = c.queries.DB().QueryRowContext(ctx, `SELECT es.id,es.chunk_set_id,ep.dimensions FROM embedding_sets es JOIN embedding_plans ep ON ep.id=es.embedding_plan_id WHERE es.id=?`, ref.ID).Scan(&metadata.Ref.ID, &metadata.ChunkSetID, &metadata.Dimensions)
	case EvaluationDatasetArtifact:
		err = c.queries.DB().QueryRowContext(ctx, `SELECT id,corpus_snapshot_id,status FROM evaluation_datasets WHERE id=?`, ref.ID).Scan(&metadata.Ref.ID, &metadata.CorpusSnapshotID, &metadata.Status)
	case RepresentationSetArtifact:
		err = c.queries.DB().QueryRowContext(ctx, `SELECT id,chunk_set_id FROM representation_sets WHERE id=?`, ref.ID).Scan(&metadata.Ref.ID, &metadata.ChunkSetID)
	default:
		return ArtifactMetadata{}, errors.Errorf("RAG_CATALOG_FAILURE: unsupported artifact kind %q", ref.Kind)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return ArtifactMetadata{}, ErrArtifactNotFound
	}
	if err != nil {
		return ArtifactMetadata{}, errors.Wrap(err, "look up immutable artifact")
	}
	return metadata, nil
}
