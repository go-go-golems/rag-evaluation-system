package widgetmanifest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Validate(catalog *Catalog) []Finding {
	var findings []Finding
	findings = append(findings, validateWidgets(catalog)...)
	findings = append(findings, validateRecipes(catalog)...)
	if len(findings) == 0 {
		findings = append(findings, Finding{Severity: "ok", Check: "summary", Message: "all discovered widget and recipe manifests passed implemented checks"})
	}
	return findings
}

func validateWidgets(catalog *Catalog) []Finding {
	var findings []Finding
	byType := map[string]string{}
	byModuleHelper := map[string]string{}

	for _, widget := range catalog.Widgets {
		m := widget.Manifest
		path := rel(catalog.Root, widget.Path)
		required := map[string]string{
			"schemaVersion":  m.SchemaVersion,
			"type":           m.Type,
			"module":         m.Module,
			"helper":         m.Helper,
			"props":          m.Props,
			"adapter.path":   m.Adapter.Path,
			"adapter.export": m.Adapter.Export,
			"reactComponent": m.ReactComponent,
			"status":         m.Status,
			"docs":           m.Docs,
		}
		for field, value := range required {
			if strings.TrimSpace(value) == "" {
				findings = append(findings, Finding{Severity: "error", Check: "required_field", Path: path, Subject: field, Message: "required widget manifest field is empty"})
			}
		}
		if previous, ok := byType[m.Type]; ok && m.Type != "" {
			findings = append(findings, Finding{Severity: "error", Check: "duplicate_widget_type", Path: path, Subject: m.Type, Message: fmt.Sprintf("also declared in %s", previous)})
		} else if m.Type != "" {
			byType[m.Type] = path
		}
		helperKey := m.Module + ":" + m.Helper
		if previous, ok := byModuleHelper[helperKey]; ok && m.Helper != "" {
			findings = append(findings, Finding{Severity: "error", Check: "duplicate_module_helper", Path: path, Subject: helperKey, Message: fmt.Sprintf("also declared in %s", previous)})
		} else if m.Helper != "" {
			byModuleHelper[helperKey] = path
		}
		moduleSpec, ok := catalog.Modules.Modules[m.Module]
		if !ok {
			findings = append(findings, Finding{Severity: "error", Check: "unknown_module", Path: path, Subject: m.Module, Message: "widget module is not declared in schema/dsl-modules.yaml"})
		} else if !pathMatchesAnyRoot(catalog.Root, widget.Path, moduleSpec.FolderRoots) {
			findings = append(findings, Finding{Severity: "error", Check: "module_folder_mismatch", Path: path, Subject: m.Module, Message: "widget path is not under any folderRoots declared for its module"})
		}
		adapterPath := widget.AdapterPath()
		if adapterPath != "" {
			if _, err := os.Stat(adapterPath); err != nil {
				findings = append(findings, Finding{Severity: "warning", Check: "adapter_exists", Path: path, Subject: rel(catalog.Root, adapterPath), Message: "adapter file is missing; expected until Phase 2 creates colocated adapters"})
			} else if m.Adapter.Export != "" {
				b, err := os.ReadFile(adapterPath)
				if err != nil {
					findings = append(findings, Finding{Severity: "error", Check: "adapter_export", Path: path, Subject: m.Adapter.Export, Message: err.Error()})
				} else if !strings.Contains(string(b), m.Adapter.Export) {
					findings = append(findings, Finding{Severity: "error", Check: "adapter_export", Path: path, Subject: m.Adapter.Export, Message: "adapter file does not contain expected export symbol"})
				}
			}
		}
		if len(m.Slots) > 0 || len(m.Actions) > 0 {
			findings = append(findings, Finding{Severity: "warning", Check: "schema_path_validation", Path: path, Subject: m.Type, Message: "slot/action path validation against TypeScript/protobuf schema is not implemented yet"})
		}
	}
	return findings
}

func validateRecipes(catalog *Catalog) []Finding {
	var findings []Finding
	byModuleRecipe := map[string]string{}
	for _, recipe := range catalog.Recipes {
		m := recipe.Manifest
		path := rel(catalog.Root, recipe.Path)
		required := map[string]string{
			"schemaVersion":         m.SchemaVersion,
			"name":                  m.Name,
			"module":                m.Module,
			"input":                 m.Input,
			"output":                m.Output,
			"implementation.path":   m.Implementation.Path,
			"implementation.symbol": m.Implementation.Symbol,
			"status":                m.Status,
			"docs":                  m.Docs,
		}
		for field, value := range required {
			if strings.TrimSpace(value) == "" {
				findings = append(findings, Finding{Severity: "error", Check: "required_field", Path: path, Subject: field, Message: "required recipe manifest field is empty"})
			}
		}
		if _, ok := catalog.Modules.Modules[m.Module]; !ok {
			findings = append(findings, Finding{Severity: "error", Check: "unknown_module", Path: path, Subject: m.Module, Message: "recipe module is not declared in schema/dsl-modules.yaml"})
		}
		key := m.Module + ":" + m.Name
		if previous, ok := byModuleRecipe[key]; ok && m.Name != "" {
			findings = append(findings, Finding{Severity: "error", Check: "duplicate_module_recipe", Path: path, Subject: key, Message: fmt.Sprintf("also declared in %s", previous)})
		} else if m.Name != "" {
			byModuleRecipe[key] = path
		}
		if m.Implementation.Path != "" {
			implPath := filepath.Join(catalog.Root, m.Implementation.Path)
			if _, err := os.Stat(implPath); err != nil {
				findings = append(findings, Finding{Severity: "error", Check: "recipe_implementation", Path: path, Subject: m.Implementation.Path, Message: "implementation file does not exist"})
			}
		}
	}
	return findings
}

func pathMatchesAnyRoot(root, path string, roots []string) bool {
	if len(roots) == 0 {
		return true
	}
	relPath := filepath.ToSlash(rel(root, path))
	for _, candidate := range roots {
		candidate = strings.TrimSuffix(filepath.ToSlash(candidate), "/")
		if strings.HasPrefix(relPath, candidate+"/") || relPath == candidate {
			return true
		}
	}
	return false
}

func rel(root, path string) string {
	r, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(r)
}
