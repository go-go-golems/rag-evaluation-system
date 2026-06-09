package widgetmanifest

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

func Discover(root string) (*Catalog, error) {
	root = normalizeRoot(root)
	absRoot, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	if stat, err := os.Stat(absRoot); err != nil {
		return nil, err
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("root %s is not a directory", absRoot)
	}

	catalog := &Catalog{Root: absRoot}
	modulePath := filepath.Join(absRoot, "schema", "dsl-modules.yaml")
	if err := readYAML(modulePath, &catalog.Modules); err != nil {
		return nil, fmt.Errorf("read module manifest %s: %w", modulePath, err)
	}

	err = filepath.WalkDir(absRoot, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			name := d.Name()
			if name == ".git" || name == "node_modules" || name == "dist" || name == "build" || name == "ttmp" {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".widget.yaml") {
			var manifest WidgetManifest
			if err := readYAML(path, &manifest); err != nil {
				return fmt.Errorf("read widget manifest %s: %w", path, err)
			}
			catalog.Widgets = append(catalog.Widgets, DiscoveredWidget{Path: path, Manifest: manifest})
			return nil
		}
		if strings.HasSuffix(path, ".recipe.yaml") {
			var manifest RecipeManifest
			if err := readYAML(path, &manifest); err != nil {
				return fmt.Errorf("read recipe manifest %s: %w", path, err)
			}
			catalog.Recipes = append(catalog.Recipes, DiscoveredRecipe{Path: path, Manifest: manifest})
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return catalog, nil
}

func normalizeRoot(root string) string {
	if root == "" {
		return "."
	}
	root = filepath.Clean(root)
	if strings.HasSuffix(root, string(filepath.Separator)+"...") {
		return strings.TrimSuffix(root, string(filepath.Separator)+"...")
	}
	if root == "..." || root == "./..." {
		return "."
	}
	return root
}

func readYAML(path string, target any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, target)
}
