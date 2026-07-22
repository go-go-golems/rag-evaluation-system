package workflowv3ttc

import (
	"context"
	"github.com/go-go-golems/scraper/pkg/workflowv3"
	"github.com/go-go-golems/scraper/pkg/workflowv3runtime"
	"path/filepath"
	"testing"
)

func TestGenerateTaskDirectContract(t *testing.T) {
	ctx := context.Background()
	b, _ := Bundle()
	r, _ := Registry()
	c, _ := r.Catalog()
	spec, _ := c.Lookup(GenerateKey)
	a, _ := workflowv3.NewFileArtifactStore(filepath.Join(t.TempDir(), "a"), 1<<20)
	text := "x"
	d, _ := workflowv3.Digest(text)
	body, _ := workflowv3.CanonicalJSON(Chunk{Key: "k", Text: text, TextDigest: d, CitationIDs: []string{"c"}, SourceDigest: d})
	ref, _ := a.Put(ctx, ChunkSchema, "application/json", body)
	p := &fixtureProvider{calls: map[string]int{}}
	m, _ := workflowv3runtime.NewTaskModuleRegistry(Module(p))
	_, err := workflowv3runtime.RunTask(ctx, workflowv3runtime.TaskRequest{RunID: "r", NodeKey: "n", Attempt: 1, Task: workflowv3.RegisteredTask{Spec: spec, Bundle: b}, Inputs: map[string]workflowv3.ArtifactRef{"chunk": ref}, Artifacts: a, Modules: m})
	if err != nil {
		t.Fatal(err)
	}
}
