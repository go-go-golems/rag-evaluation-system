package chunking

import (
	"strings"
	"testing"
)

func TestFixedSizeChunkerTerminatesAtEndWithOverlap(t *testing.T) {
	text := strings.Repeat("abcdefghij", 8) // 80 runes

	chunker := NewFixedSizeChunker(50, 10, "test")
	chunks, err := chunker.Chunk("test-doc", text)
	if err != nil {
		t.Fatalf("chunking failed: %v", err)
	}

	if got, want := len(chunks), 2; got != want {
		t.Fatalf("expected %d chunks, got %d: %#v", want, got, chunks)
	}
	if chunks[0].StartOffset != 0 || chunks[0].EndOffset != 50 {
		t.Fatalf("unexpected first chunk offsets: %#v", chunks[0])
	}
	if chunks[1].StartOffset != 40 || chunks[1].EndOffset != 80 {
		t.Fatalf("unexpected second chunk offsets: %#v", chunks[1])
	}
}

func TestFixedSizeChunkerRejectsOverlapAtLeastChunkSize(t *testing.T) {
	chunker := NewFixedSizeChunker(10, 10, "test")
	_, err := chunker.Chunk("test-doc", "Short text.")
	if err == nil {
		t.Fatal("expected overlap >= chunk size to return an error")
	}
}

func TestFixedSizeChunkerRejectsNegativeOverlap(t *testing.T) {
	chunker := NewFixedSizeChunker(10, -1, "test")
	_, err := chunker.Chunk("test-doc", "Short text.")
	if err == nil {
		t.Fatal("expected negative overlap to return an error")
	}
}

func TestSentenceChunkerRejectsOverlapAtLeastChunkSize(t *testing.T) {
	chunker := NewSentenceChunker(10, 10, "test")
	_, err := chunker.Chunk("test-doc", "Short text.")
	if err == nil {
		t.Fatal("expected overlap >= chunk size to return an error")
	}
}

func TestChunkTextAlwaysEqualsRecordedRuneRange(t *testing.T) {
	text := "  First sentence.  Second sentence.\n\n# Heading\nThird sentence.  "
	chunkers := []Chunker{
		NewFixedSizeChunker(20, 4, "fixed"),
		NewSentenceChunker(25, 4, "sentence"),
		NewMarkdownHeadingChunker(25, "markdown"),
	}
	for _, chunker := range chunkers {
		chunks, err := chunker.Chunk("test-doc", text)
		if err != nil {
			t.Fatalf("%s chunking failed: %v", chunker.Name(), err)
		}
		for _, chunk := range chunks {
			if actual := string([]rune(text)[chunk.StartOffset:chunk.EndOffset]); actual != chunk.Text {
				t.Fatalf("%s chunk %#v does not equal source range %q", chunker.Name(), chunk, actual)
			}
		}
	}
}
