package chunking

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Chunk represents a single text chunk extracted from a document
type Chunk struct {
	ID          string `json:"id"`
	DocumentID  string `json:"document_id"`
	StrategyID  string `json:"strategy_id"`
	ChunkIndex  int    `json:"chunk_index"`
	Text        string `json:"text"`
	TokenCount  int    `json:"token_count"`
	StartOffset int    `json:"start_offset"`
	EndOffset   int    `json:"end_offset"`
}

// Chunker is the interface for chunking strategies
type Chunker interface {
	// Name returns the human-readable name of the chunking strategy
	Name() string
	// Chunk splits document text into chunks
	Chunk(documentID, text string) ([]Chunk, error)
}

// FixedSizeChunker splits text into fixed-size chunks with optional overlap
type FixedSizeChunker struct {
	ChunkSize  int // in characters
	Overlap    int // overlap in characters
	StrategyID string
}

func NewFixedSizeChunker(chunkSize, overlap int, strategyID string) *FixedSizeChunker {
	return &FixedSizeChunker{
		ChunkSize:  chunkSize,
		Overlap:    overlap,
		StrategyID: strategyID,
	}
}

func (c *FixedSizeChunker) Name() string {
	return fmt.Sprintf("fixed-%d-%d", c.ChunkSize, c.Overlap)
}

func (c *FixedSizeChunker) Chunk(documentID, text string) ([]Chunk, error) {
	if c.ChunkSize <= 0 {
		return nil, fmt.Errorf("chunk size must be positive, got %d", c.ChunkSize)
	}
	if c.Overlap < 0 {
		return nil, fmt.Errorf("overlap must be non-negative, got %d", c.Overlap)
	}
	if c.Overlap >= c.ChunkSize {
		return nil, fmt.Errorf("overlap must be smaller than chunk size: overlap=%d chunk_size=%d", c.Overlap, c.ChunkSize)
	}

	var chunks []Chunk
	runes := []rune(text)
	totalRunes := len(runes)
	start := 0
	index := 0

	for start < totalRunes {
		end := start + c.ChunkSize
		if end > totalRunes {
			end = totalRunes
		}

		chunkText := string(runes[start:end])
		chunkID := generateChunkID(documentID, c.StrategyID, index)

		chunks = append(chunks, Chunk{
			ID:          chunkID,
			DocumentID:  documentID,
			StrategyID:  c.StrategyID,
			ChunkIndex:  index,
			Text:        chunkText,
			TokenCount:  estimateTokens(chunkText),
			StartOffset: start,
			EndOffset:   end,
		})

		index++

		// If we've reached the end of the text, we're done.
		if end >= totalRunes {
			break
		}

		// Advance start, respecting overlap but ensuring forward progress
		nextStart := end - c.Overlap
		if nextStart <= start {
			// overlap >= chunkSize would cause no progress; advance by at least 1
			nextStart = start + 1
		}
		if nextStart < 0 {
			nextStart = 0
		}
		start = nextStart
	}

	return chunks, nil
}

// SentenceChunker splits text on sentence boundaries
type SentenceChunker struct {
	ChunkSize  int // target size in characters
	Overlap    int // overlap in characters
	StrategyID string
}

func NewSentenceChunker(chunkSize, overlap int, strategyID string) *SentenceChunker {
	return &SentenceChunker{
		ChunkSize:  chunkSize,
		Overlap:    overlap,
		StrategyID: strategyID,
	}
}

func (c *SentenceChunker) Name() string {
	return fmt.Sprintf("sentence-%d-%d", c.ChunkSize, c.Overlap)
}

func (c *SentenceChunker) Chunk(documentID, text string) ([]Chunk, error) {
	if c.ChunkSize <= 0 {
		return nil, fmt.Errorf("chunk size must be positive, got %d", c.ChunkSize)
	}
	if c.Overlap < 0 {
		return nil, fmt.Errorf("overlap must be non-negative, got %d", c.Overlap)
	}
	if c.Overlap >= c.ChunkSize {
		return nil, fmt.Errorf("overlap must be smaller than chunk size: overlap=%d chunk_size=%d", c.Overlap, c.ChunkSize)
	}

	// Split into sentences
	sentences := splitSentences(text)

	var chunks []Chunk
	currentText := ""
	currentStart := 0
	chunkStartOffset := 0
	index := 0

	for _, sent := range sentences {
		if utf8.RuneCountInString(currentText)+utf8.RuneCountInString(sent) > c.ChunkSize && currentText != "" {
			// Emit current chunk
			chunkID := generateChunkID(documentID, c.StrategyID, index)
			chunks = append(chunks, Chunk{
				ID:         chunkID,
				DocumentID: documentID,
				StrategyID: c.StrategyID,
				ChunkIndex: index,
				// Preserve the exact source slice. Trimming here makes the
				// recorded rune offsets disagree with chunk text.
				Text:        currentText,
				TokenCount:  estimateTokens(currentText),
				StartOffset: chunkStartOffset,
				EndOffset:   currentStart,
			})
			index++

			// Overlap: keep last N characters of current chunk
			runes := []rune(currentText)
			if c.Overlap > 0 && len(runes) > c.Overlap {
				currentText = string(runes[len(runes)-c.Overlap:])
				chunkStartOffset = currentStart - utf8.RuneCountInString(currentText)
			} else {
				currentText = ""
				chunkStartOffset = currentStart
			}
		}

		if currentText == "" {
			chunkStartOffset = currentStart
		}
		currentText += sent
		currentStart += utf8.RuneCountInString(sent)
	}

	// Emit last chunk
	if currentText != "" {
		chunkID := generateChunkID(documentID, c.StrategyID, index)
		chunks = append(chunks, Chunk{
			ID:          chunkID,
			DocumentID:  documentID,
			StrategyID:  c.StrategyID,
			ChunkIndex:  index,
			Text:        currentText,
			TokenCount:  estimateTokens(currentText),
			StartOffset: chunkStartOffset,
			EndOffset:   currentStart,
		})
	}

	return chunks, nil
}

// MarkdownHeadingChunker splits text on markdown headings
type MarkdownHeadingChunker struct {
	MaxChunkSize int // if a section exceeds this, sub-chunk with fixed size
	StrategyID   string
}

func NewMarkdownHeadingChunker(maxChunkSize int, strategyID string) *MarkdownHeadingChunker {
	return &MarkdownHeadingChunker{
		MaxChunkSize: maxChunkSize,
		StrategyID:   strategyID,
	}
}

func (c *MarkdownHeadingChunker) Name() string {
	return "markdown-heading"
}

func (c *MarkdownHeadingChunker) Chunk(documentID, text string) ([]Chunk, error) {
	sections := splitMarkdownSections(text)

	var chunks []Chunk
	index := 0
	offset := 0

	for _, section := range sections {
		if c.MaxChunkSize > 0 && utf8.RuneCountInString(section) > c.MaxChunkSize {
			// Sub-chunk large sections. Keep the fallback overlap below the
			// chunk size invariant enforced by FixedSizeChunker.
			overlap := 100
			if overlap >= c.MaxChunkSize {
				overlap = c.MaxChunkSize / 5
			}
			subChunker := NewFixedSizeChunker(c.MaxChunkSize, overlap, c.StrategyID)
			subChunks, err := subChunker.Chunk(documentID, section)
			if err != nil {
				return nil, err
			}
			for i := range subChunks {
				subChunks[i].ChunkIndex = index
				subChunks[i].ID = generateChunkID(documentID, c.StrategyID, index)
				subChunks[i].StrategyID = c.StrategyID
				subChunks[i].StartOffset += offset
				subChunks[i].EndOffset += offset
				chunks = append(chunks, subChunks[i])
				index++
			}
		} else {
			chunkID := generateChunkID(documentID, c.StrategyID, index)
			chunks = append(chunks, Chunk{
				ID:          chunkID,
				DocumentID:  documentID,
				StrategyID:  c.StrategyID,
				ChunkIndex:  index,
				Text:        section,
				TokenCount:  estimateTokens(section),
				StartOffset: offset,
				EndOffset:   offset + utf8.RuneCountInString(section),
			})
			index++
		}
		offset += utf8.RuneCountInString(section)
	}

	return chunks, nil
}

// NewChunkerFromType creates a chunker by strategy type name
func NewChunkerFromType(strategyType string, chunkSize, overlap int, strategyID string) (Chunker, error) {
	switch strategyType {
	case "fixed":
		return NewFixedSizeChunker(chunkSize, overlap, strategyID), nil
	case "sentence":
		return NewSentenceChunker(chunkSize, overlap, strategyID), nil
	case "markdown-heading":
		return NewMarkdownHeadingChunker(chunkSize*2, strategyID), nil // max = 2x chunk size
	default:
		return nil, fmt.Errorf("unknown chunking strategy: %s", strategyType)
	}
}

// --- Helpers ---

func generateChunkID(documentID, strategyID string, index int) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%d", documentID, strategyID, index)))
	return "chk-" + hex.EncodeToString(h[:])[:16]
}

// estimateTokens approximates token count (rough: 1 token ≈ 4 chars for English)
func estimateTokens(text string) int {
	return utf8.RuneCountInString(text) / 4
}

// splitSentences splits text into sentences
func splitSentences(text string) []string {
	var sentences []string
	var current strings.Builder

	for _, r := range text {
		current.WriteRune(r)
		if r == '.' || r == '!' || r == '?' || r == '\n' {
			// Keep whitespace-only fragments as source material. Dropping them
			// changes subsequent rune offsets and breaks citation evidence.
			sentences = append(sentences, current.String())
			current.Reset()
		}
	}

	if current.Len() > 0 {
		sentences = append(sentences, current.String())
	}

	return sentences
}

// splitMarkdownSections splits text on markdown headings
func splitMarkdownSections(text string) []string {
	var sections []string
	var current strings.Builder

	// SplitAfter preserves whether each line actually ended in a newline.
	// Reconstructing lines by unconditionally adding one changes source ranges.
	lines := strings.SplitAfter(text, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "#") && current.Len() > 0 {
			sections = append(sections, current.String())
			current.Reset()
		}
		current.WriteString(line)
	}

	if current.Len() > 0 {
		sections = append(sections, current.String())
	}

	return sections
}
