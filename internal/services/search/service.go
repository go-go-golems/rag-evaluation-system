package search

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
)

const (
	DefaultIndexRoot    = "data/indexes"
	DefaultLimit        = 10
	DefaultPreviewRunes = 220
	MaxLimit            = 100
)

// Service provides retrieval operations over derived indexes and canonical
// SQLite chunk/embedding state.
type Service struct {
	queries   *db.Queries
	indexRoot string
}

// NewService creates a search service. If indexRoot is empty, data/indexes is used.
func NewService(queries *db.Queries, indexRoot string) *Service {
	if indexRoot == "" {
		indexRoot = DefaultIndexRoot
	}
	return &Service{queries: queries, indexRoot: indexRoot}
}

type BuildIndexRequest struct {
	IndexID     string
	StrategyID  string
	SourceIDs   []string
	DocumentIDs []string
	Force       bool
	Limit       int
}

type BuildIndexResult struct {
	IndexID       string   `json:"index_id"`
	StrategyID    string   `json:"strategy_id"`
	SourceIDs     []string `json:"source_ids,omitempty"`
	DocumentIDs   []string `json:"document_ids,omitempty"`
	IndexPath     string   `json:"index_path"`
	ChunkCount    int      `json:"chunk_count"`
	DocumentCount int      `json:"document_count"`
	Rebuilt       bool     `json:"rebuilt"`
}

type QueryRequest struct {
	IndexID      string
	Query        string
	Limit        int
	PreviewRunes int
}

type RetrievalResult struct {
	Rank       int                           `json:"rank"`
	ChunkID    string                        `json:"chunk_id"`
	DocumentID string                        `json:"document_id"`
	SourceID   string                        `json:"source_id"`
	Title      string                        `json:"title"`
	URL        string                        `json:"url,omitempty"`
	StrategyID string                        `json:"strategy_id"`
	ChunkIndex int                           `json:"chunk_index"`
	Score      float64                       `json:"score"`
	Retriever  string                        `json:"retriever"`
	Preview    string                        `json:"preview"`
	Components map[string]RetrievalComponent `json:"components,omitempty"`
}

type RetrievalComponent struct {
	Rank  int     `json:"rank"`
	Score float64 `json:"score"`
}

type QueryResult struct {
	Query     string            `json:"query"`
	IndexID   string            `json:"index_id,omitempty"`
	Retriever string            `json:"retriever"`
	Items     []RetrievalResult `json:"items"`
}

func (s *Service) indexPath(indexID string) string {
	return filepath.Join(s.indexRoot, "bm25", indexID)
}

func normalizeSourceIDs(sourceIDs []string) []string {
	seen := map[string]bool{}
	var normalized []string
	for _, sourceID := range sourceIDs {
		for _, part := range strings.Split(sourceID, ",") {
			trimmed := strings.TrimSpace(part)
			if trimmed == "" || seen[trimmed] {
				continue
			}
			seen[trimmed] = true
			normalized = append(normalized, trimmed)
		}
	}
	sort.Strings(normalized)
	return normalized
}

func deriveIndexID(strategyID string, sourceIDs []string) string {
	parts := append([]string{strategyID}, sourceIDs...)
	sum := sha256.Sum256([]byte(strings.Join(parts, "\x00")))
	return fmt.Sprintf("bm25-%s-%s", slug(strategyID), hex.EncodeToString(sum[:])[:12])
}

func slug(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastDash := false
	for _, r := range s {
		ok := (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9')
		if ok {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteRune('-')
			lastDash = true
		}
	}
	return strings.Trim(b.String(), "-")
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}

func preview(text string, maxRunes int) string {
	if maxRunes <= 0 || text == "" {
		return ""
	}
	runes := []rune(strings.Join(strings.Fields(text), " "))
	if len(runes) <= maxRunes {
		return string(runes)
	}
	return string(runes[:maxRunes]) + "…"
}
