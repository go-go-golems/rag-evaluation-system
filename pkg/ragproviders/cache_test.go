package ragproviders

import (
	"os"
	"path/filepath"
	"sync"
	"testing"
)

func TestDiskCacheConcurrentAtomicReplacement(t *testing.T) {
	cache, err := NewDiskCache(CacheConfig{Directory: t.TempDir(), MaxEntryBytes: 64})
	if err != nil {
		t.Fatal(err)
	}
	values := [][]byte{[]byte("one"), []byte("two"), []byte("three")}
	var group sync.WaitGroup
	for _, value := range values {
		value := value
		group.Add(1)
		go func() {
			defer group.Done()
			for i := 0; i < 50; i++ {
				cache.Put("key", value)
			}
		}()
	}
	group.Wait()
	got, ok := cache.Get("key")
	if !ok {
		t.Fatal("Get() = miss")
	}
	valid := false
	for _, value := range values {
		if string(got) == string(value) {
			valid = true
		}
	}
	if !valid {
		t.Fatalf("Get() = %q, want complete writer value", got)
	}
	entries, err := os.ReadDir(filepath.Dir(cache.path("key")))
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if filepath.Ext(entry.Name()) == ".cache" {
			t.Fatalf("temporary cache file retained: %s", entry.Name())
		}
	}
}

func TestDiskCacheTreatsCorruptOrOversizedEntryAsMiss(t *testing.T) {
	cache, err := NewDiskCache(CacheConfig{Directory: t.TempDir(), MaxEntryBytes: 3})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(cache.path("bad"), []byte("toolong"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, ok := cache.Get("bad"); ok {
		t.Fatal("Get() = hit for oversized entry")
	}
}
