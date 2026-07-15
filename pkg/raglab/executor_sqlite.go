package raglab

import (
	"context"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	"github.com/go-go-golems/rag-evaluation-system/internal/services/immutableretrieval"
	"github.com/pkg/errors"
)

// SQLiteChannelRetriever executes raw immutable lexical and vector channels
// against the artifact IDs recorded in a laboratory database.
type SQLiteChannelRetriever struct{ queries *db.Queries }

var _ ChannelRetriever = (*SQLiteChannelRetriever)(nil)

func NewSQLiteChannelRetriever(queries *db.Queries) *SQLiteChannelRetriever {
	return &SQLiteChannelRetriever{queries: queries}
}

func (r *SQLiteChannelRetriever) BM25(ctx context.Context, artifactID, query string, limit int) ([]immutableretrieval.ChunkHit, error) {
	if r == nil || r.queries == nil {
		return nil, errors.New("RAG_EXECUTOR_REQUIRED: database queries are required")
	}
	return immutableretrieval.QueryBM25(ctx, r.queries, artifactID, query, limit)
}

func (r *SQLiteChannelRetriever) Vector(ctx context.Context, embeddingSetID string, vector []float32, limit int) ([]immutableretrieval.ChunkHit, error) {
	if r == nil || r.queries == nil {
		return nil, errors.New("RAG_EXECUTOR_REQUIRED: database queries are required")
	}
	return immutableretrieval.QueryVector(ctx, r.queries, embeddingSetID, vector, limit)
}
