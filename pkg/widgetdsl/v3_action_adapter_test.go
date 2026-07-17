package widgetdsl

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestWidgetV3SerializedActionPropsHaveIRAndAdapterConsumers(t *testing.T) {
	actionPropPattern := regexp.MustCompile(`props\["([A-Za-z0-9]+Action)"\]`)
	serialized := map[string]struct{}{}
	for _, path := range []string{"v3.go", "v3_crm.go"} {
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatalf("read %s: %v", path, err)
		}
		for _, match := range actionPropPattern.FindAllStringSubmatch(string(data), -1) {
			serialized[match[1]] = struct{}{}
		}
	}

	root := filepath.Join("..", "..", "packages", "rag-evaluation-site", "src")
	irSources := strings.Builder{}
	err := filepath.WalkDir(filepath.Join(root, "widgets", "ir"), func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".ts") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		irSources.Write(data)
		irSources.WriteByte('\n')
		return nil
	})
	if err != nil {
		t.Fatalf("read Widget IR declarations: %v", err)
	}
	adapterSources := strings.Builder{}
	err = filepath.WalkDir(filepath.Join(root, "components"), func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.HasSuffix(path, ".widget.tsx") {
			return nil
		}
		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		adapterSources.Write(data)
		adapterSources.WriteByte('\n')
		return nil
	})
	if err != nil {
		t.Fatalf("read Widget adapters: %v", err)
	}

	for actionProp := range serialized {
		if !strings.Contains(irSources.String(), actionProp+"?") {
			t.Errorf("v3 serializes %s but Widget IR props do not declare it", actionProp)
		}
		if !strings.Contains(adapterSources.String(), "props."+actionProp) {
			t.Errorf("v3 serializes %s but no React Widget adapter consumes props.%s", actionProp, actionProp)
		}
	}
	if len(serialized) < 30 {
		t.Fatalf("action audit found only %d serialized props; source extraction likely drifted", len(serialized))
	}
}
