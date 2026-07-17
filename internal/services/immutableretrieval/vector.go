package immutableretrieval

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
	legacyembedding "github.com/go-go-golems/rag-evaluation-system/internal/services/embedding"
	"github.com/pkg/errors"
)

// QueryVector exhaustively scores every vector in an immutable embedding set.
// No candidate limit is applied before scoring.
func QueryVector(ctx context.Context, q *db.Queries, embeddingSetID string, queryVector []float32, limit int) ([]ChunkHit, error) {
	if q == nil || embeddingSetID == "" || len(queryVector) == 0 {
		return nil, errors.New("database queries, embedding set ID, and query vector are required")
	}
	if limit <= 0 {
		limit = 10
	}
	rows, err := q.DB().QueryContext(ctx, `SELECT ie.chunk_id, ie.vector FROM immutable_embeddings ie WHERE ie.embedding_set_id = ? ORDER BY ie.chunk_id`, embeddingSetID)
	if err != nil {
		return nil, errors.Wrap(err, "load immutable embedding set")
	}
	defer func() { _ = rows.Close() }()
	var hits []ChunkHit
	for rows.Next() {
		var chunkID string
		var blob []byte
		if err := rows.Scan(&chunkID, &blob); err != nil {
			return nil, errors.Wrap(err, "scan immutable embedding")
		}
		vector, err := legacyembedding.DecodeFloat32Vector(blob)
		if err != nil {
			return nil, errors.Wrapf(err, "decode immutable embedding %s", chunkID)
		}
		score, err := legacyembedding.CosineSimilarity(queryVector, vector)
		if err != nil {
			return nil, fmt.Errorf("score immutable embedding %s: %w", chunkID, err)
		}
		hits = append(hits, ChunkHit{ChunkID: chunkID, Score: score, Channel: "vector"})
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "iterate immutable embeddings")
	}
	sort.SliceStable(hits, func(i, j int) bool {
		if hits[i].Score == hits[j].Score {
			return hits[i].ChunkID < hits[j].ChunkID
		}
		return hits[i].Score > hits[j].Score
	})
	if len(hits) > limit {
		hits = hits[:limit]
	}
	for i := range hits {
		hits[i].Rank = i + 1
	}
	return hydrate(ctx, q, hits, nil)
}

// CollapseDocuments keeps the first (therefore highest-ranked) evidence chunk
// per immutable document revision within one retrieval channel.
func CollapseDocuments(hits []ChunkHit) []ChunkHit {
	seen := map[string]struct{}{}
	result := make([]ChunkHit, 0, len(hits))
	for _, hit := range hits {
		if _, ok := seen[hit.DocumentRevisionID]; ok {
			continue
		}
		seen[hit.DocumentRevisionID] = struct{}{}
		result = append(result, hit)
	}
	for i := range result {
		result[i].Rank = i + 1
	}
	return result
}

type RRFComponent struct {
	Rank           int     `json:"rank"`
	Score          float64 `json:"score"`
	WinningChunkID string  `json:"winning_chunk_id"`
	Contribution   float64 `json:"contribution"`
}
type FusedHit struct {
	ChunkHit
	Components map[string]RRFComponent `json:"components"`
}

// FuseRRF collapses each named channel by document revision before fusing.
func FuseRRF(channels map[string][]ChunkHit, rankConstant, limit int) []FusedHit {
	return FuseWeightedRRF(channels, rankConstant, nil, limit)
}

// FuseWeightedRRF applies reciprocal-rank fusion after deterministic
// document-level collapse. Missing channel weights default to one. Supplying a
// non-positive weight deliberately contributes no score; plan validation is
// responsible for rejecting such weights before execution.
func FuseWeightedRRF(channels map[string][]ChunkHit, rankConstant int, weights map[string]float64, limit int) []FusedHit {
	if rankConstant <= 0 {
		rankConstant = 60
	}
	if limit <= 0 {
		limit = 10
	}
	fused := map[string]*FusedHit{}
	winningContribution := map[string]float64{}
	winningChannel := map[string]string{}
	for name, hits := range channels {
		weight := 1.0
		if configured, ok := weights[name]; ok {
			weight = configured
		}
		for _, hit := range CollapseDocuments(hits) {
			contribution := weight / float64(rankConstant+hit.Rank)
			item := fused[hit.DocumentRevisionID]
			if item == nil {
				candidate := FusedHit{ChunkHit: hit, Components: map[string]RRFComponent{}}
				candidate.Score = 0
				item = &candidate
				fused[hit.DocumentRevisionID] = item
				winningContribution[hit.DocumentRevisionID] = contribution
				winningChannel[hit.DocumentRevisionID] = name
			} else if contribution > winningContribution[hit.DocumentRevisionID] ||
				(contribution == winningContribution[hit.DocumentRevisionID] && (name < winningChannel[hit.DocumentRevisionID] ||
					(name == winningChannel[hit.DocumentRevisionID] && hit.ChunkID < item.ChunkID))) {
				fusedScore := item.Score
				item.ChunkHit = hit
				item.Score = fusedScore
				winningContribution[hit.DocumentRevisionID] = contribution
				winningChannel[hit.DocumentRevisionID] = name
			}
			item.Score += contribution
			item.Components[name] = RRFComponent{hit.Rank, hit.Score, hit.ChunkID, contribution}
		}
	}
	result := make([]FusedHit, 0, len(fused))
	for _, hit := range fused {
		result = append(result, *hit)
	}
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Score == result[j].Score {
			return result[i].DocumentRevisionID < result[j].DocumentRevisionID
		}
		return result[i].Score > result[j].Score
	})
	if len(result) > limit {
		result = result[:limit]
	}
	for i := range result {
		result[i].Rank = i + 1
	}
	return result
}
