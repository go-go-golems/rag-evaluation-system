package widgetmanifest

import "path/filepath"

type AdapterSpec struct {
	Path   string `yaml:"path"`
	Export string `yaml:"export"`
}

type WidgetManifest struct {
	SchemaVersion  string      `yaml:"schemaVersion"`
	Type           string      `yaml:"type"`
	Module         string      `yaml:"module"`
	Helper         string      `yaml:"helper"`
	Props          string      `yaml:"props"`
	Adapter        AdapterSpec `yaml:"adapter"`
	ReactComponent string      `yaml:"reactComponent"`
	Children       bool        `yaml:"children"`
	Slots          []string    `yaml:"slots"`
	Actions        []string    `yaml:"actions"`
	Status         string      `yaml:"status"`
	Docs           string      `yaml:"docs"`
}

type RecipeImplementation struct {
	Path   string `yaml:"path"`
	Symbol string `yaml:"symbol"`
}

type RecipeManifest struct {
	SchemaVersion  string               `yaml:"schemaVersion"`
	Name           string               `yaml:"name"`
	Module         string               `yaml:"module"`
	Input          string               `yaml:"input"`
	Output         string               `yaml:"output"`
	Implementation RecipeImplementation `yaml:"implementation"`
	Status         string               `yaml:"status"`
	Docs           string               `yaml:"docs"`
}

type ModuleManifest struct {
	SchemaVersion string                `yaml:"schemaVersion"`
	Modules       map[string]ModuleSpec `yaml:"modules"`
}

type ModuleSpec struct {
	Description     string   `yaml:"description"`
	PackageID       string   `yaml:"packageId"`
	DependsOn       []string `yaml:"dependsOn"`
	FolderRoots     []string `yaml:"folderRoots"`
	OwnsPageHelper  bool     `yaml:"ownsPageHelper"`
	OwnsCellHelpers bool     `yaml:"ownsCellHelpers"`
}

type DiscoveredWidget struct {
	Path     string
	Manifest WidgetManifest
}

func (w DiscoveredWidget) AdapterPath() string {
	if w.Manifest.Adapter.Path == "" {
		return ""
	}
	return filepath.Clean(filepath.Join(filepath.Dir(w.Path), w.Manifest.Adapter.Path))
}

type DiscoveredRecipe struct {
	Path     string
	Manifest RecipeManifest
}

type Catalog struct {
	Root    string
	Modules ModuleManifest
	Widgets []DiscoveredWidget
	Recipes []DiscoveredRecipe
}

type Finding struct {
	Severity string
	Check    string
	Path     string
	Subject  string
	Message  string
}

func (f Finding) IsError() bool { return f.Severity == "error" }
