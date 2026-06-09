package ingest

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-go-golems/rag-evaluation-system/internal/db"
)

// Scanner discovers documents in a filesystem directory
type Scanner struct {
	db *db.Queries
}

// NewScanner creates a scanner that stores discovered documents
func NewScanner(queries *db.Queries) *Scanner {
	return &Scanner{db: queries}
}

// ScanDir walks a directory tree and ingests all text-readable files as documents
func (s *Scanner) ScanDir(sourceID, dirPath string) ([]string, error) {
	var docIDs []string

	root, err := os.OpenRoot(dirPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open scan root %s: %w", dirPath, err)
	}
	defer func() { _ = root.Close() }()

	err = filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// Skip hidden directories
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip hidden files, symlinks, and non-text files.
		if strings.HasPrefix(d.Name(), ".") || d.Type()&fs.ModeSymlink != 0 {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		if !isTextFile(ext) {
			return nil
		}

		// Generate stable ID from path hash.
		relPath, err := filepath.Rel(dirPath, path)
		if err != nil {
			return fmt.Errorf("failed to resolve relative path for %s: %w", path, err)
		}

		content, err := root.ReadFile(relPath)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", path, err)
		}

		docID := "doc-" + hashString(sourceID + ":" + relPath)[:16]
		title := extractTitle(content, d.Name())
		wordCount := countWords(string(content))

		err = s.db.InsertDocument(
			docID,
			sourceID,
			relPath,
			title,
			"", // author
			"", // url
			detectContentType(ext),
			string(content),
			string(content),
			"", // content_html
			wordCount,
			"en",
			"extracted",
		)
		if err != nil {
			return fmt.Errorf("failed to insert document %s: %w", path, err)
		}

		docIDs = append(docIDs, docID)
		return nil
	})

	return docIDs, err
}

func isTextFile(ext string) bool {
	textExts := map[string]bool{
		".txt": true, ".md": true, ".markdown": true, ".rst": true,
		".go": true, ".py": true, ".js": true, ".ts": true, ".tsx": true,
		".java": true, ".c": true, ".cpp": true, ".h": true, ".hpp": true,
		".rs": true, ".rb": true, ".sh": true, ".bash": true, ".zsh": true,
		".yaml": true, ".yml": true, ".json": true, ".toml": true, ".ini": true,
		".csv": true, ".html": true, ".htm": true, ".xml": true, ".svg": true,
		".sql": true, ".css": true, ".scss": true, ".less": true,
		".dockerfile": true, ".makefile": true,
	}
	return textExts[ext]
}

func detectContentType(ext string) string {
	switch ext {
	case ".md", ".markdown":
		return "markdown"
	case ".html", ".htm":
		return "html"
	case ".json":
		return "json"
	case ".yaml", ".yml":
		return "yaml"
	case ".csv":
		return "csv"
	default:
		return "text"
	}
}

func extractTitle(content []byte, filename string) string {
	text := string(content)
	lines := strings.Split(text, "\n")

	// Try to find a markdown heading
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# ") {
			return strings.TrimPrefix(trimmed, "# ")
		}
	}

	// Fall back to filename
	name := filepath.Base(filename)
	ext := filepath.Ext(name)
	return strings.TrimSuffix(name, ext)
}

func countWords(text string) int {
	return len(strings.Fields(text))
}

func hashString(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}
