package ragproviders

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"sync"

	"github.com/go-go-golems/rag-evaluation-system/pkg/ragoperators"
)

type DiskCache struct {
	directory     string
	maxEntryBytes int64
	mu            sync.RWMutex
}

var _ ragoperators.Cache = (*DiskCache)(nil)

func NewDiskCache(cfg CacheConfig) (*DiskCache, error) {
	if cfg.Directory == "" {
		return nil, os.ErrInvalid
	}
	maxBytes := cfg.MaxEntryBytes
	if maxBytes == 0 {
		maxBytes = 4 << 20
	}
	if maxBytes < 1 {
		return nil, os.ErrInvalid
	}
	if err := os.MkdirAll(cfg.Directory, 0o700); err != nil {
		return nil, err
	}
	return &DiskCache{directory: cfg.Directory, maxEntryBytes: maxBytes}, nil
}
func (c *DiskCache) path(key string) string {
	sum := sha256.Sum256([]byte(key))
	return filepath.Join(c.directory, hex.EncodeToString(sum[:])+".json")
}
func (c *DiskCache) Get(key string) ([]byte, bool) {
	if c == nil {
		return nil, false
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	data, err := os.ReadFile(c.path(key))
	if err != nil || int64(len(data)) > c.maxEntryBytes {
		return nil, false
	}
	return data, true
}
func (c *DiskCache) Put(key string, value []byte) {
	if c == nil || int64(len(value)) > c.maxEntryBytes {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	path := c.path(key)
	tmp, err := os.CreateTemp(c.directory, ".cache-*")
	if err != nil {
		return
	}
	tmpName := tmp.Name()
	defer func() { _ = tmp.Close(); _ = os.Remove(tmpName) }()
	if _, err := tmp.Write(value); err != nil {
		return
	}
	if err := tmp.Sync(); err != nil {
		return
	}
	if err := tmp.Close(); err != nil {
		return
	}
	_ = os.Rename(tmpName, path)
}
func (c *DiskCache) Close() error { return nil }
