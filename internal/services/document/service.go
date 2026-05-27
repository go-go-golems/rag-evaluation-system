package document

import (
	"context"
	"fmt"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
)

// Service owns document read behavior shared by CLI, HTTP, and future workflows.
type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

type ListRequest struct {
	Limit  int
	Offset int
}

func (s *Service) List(ctx context.Context, req ListRequest) ([]db.Document, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if req.Limit <= 0 {
		req.Limit = 50
	}
	if req.Offset < 0 {
		req.Offset = 0
	}
	return s.queries.ListDocuments(req.Limit, req.Offset)
}

func (s *Service) Get(ctx context.Context, id string) (*db.Document, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if id == "" {
		return nil, fmt.Errorf("document id is required")
	}
	return s.queries.GetDocument(id)
}

type ChunksRequest struct {
	DocumentID string
}

func (s *Service) Chunks(ctx context.Context, req ChunksRequest) ([]db.Chunk, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	if req.DocumentID == "" {
		return nil, fmt.Errorf("document id is required")
	}
	return s.queries.ListChunks(req.DocumentID)
}
