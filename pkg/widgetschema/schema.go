package widgetschema

const Version = "0.1.0"

var ComponentTypes = []string{
	"AppShell",
	"AppNav",
	"Button",
	"Caption",
	"CodeText",
	"ContextKindSwatch",
	"AnnotationBadge",
	"ContextLegend",
	"ContextBudgetBar",
	"ContextStripDiagram",
	"ContextStackDiagram",
	"ContextTreemap",
	"ContextDiagramPanel",
	"DashboardGrid",
	"DataTable",
	"Divider",
	"FormRow",
	"Inline",
	"MetadataGrid",
	"Panel",
	"ScrollRegion",
	"SectionBlock",
	"SelectInput",
	"SidebarShell",
	"SplitPane",
	"Stack",
	"StatusText",
	"TabList",
	"Text",
	"TextInput",
	"TranscriptRoleBadge",
	"TranscriptSessionHeader",
	"TranscriptMessageCard",
	"AnnotationNoteCard",
	"AnnotationRailPanel",
	"TranscriptReaderPanel",
	"TranscriptWorkspacePanel",
	"AnchoredCommentCard",
	"AnchoredCommentRail",
}

var CellKinds = []string{"field", "number", "status", "caption", "template", "link", "constant"}

func Summary() map[string]any {
	return map[string]any{
		"schemaVersion": Version,
		"components":    ComponentTypes,
		"cellKinds":     CellKinds,
		"jsonSchema":    JSONSchema(),
	}
}

func JSONSchema() map[string]any {
	return map[string]any{
		"$schema":              "https://json-schema.org/draft/2020-12/schema",
		"$id":                  "https://go-go-golems.github.io/rag-evaluation-site/widget-page.schema.json",
		"title":                "RAG Evaluation Widget Page",
		"type":                 "object",
		"additionalProperties": true,
		"required":             []string{"id", "title", "root"},
		"properties": map[string]any{
			"schemaVersion": map[string]any{"type": "string", "const": Version},
			"id":            map[string]any{"type": "string", "minLength": 1},
			"title":         map[string]any{"type": "string"},
			"meta":          map[string]any{"type": "object", "additionalProperties": true},
			"root":          map[string]any{"$ref": "#/$defs/widgetNode"},
		},
		"$defs": map[string]any{
			"widgetNode": map[string]any{
				"oneOf": []any{
					map[string]any{"$ref": "#/$defs/textNode"},
					map[string]any{"$ref": "#/$defs/elementNode"},
					map[string]any{"$ref": "#/$defs/componentNode"},
				},
			},
			"textNode": map[string]any{
				"type":                 "object",
				"required":             []string{"kind", "text"},
				"additionalProperties": true,
				"properties": map[string]any{
					"kind": map[string]any{"const": "text"},
					"text": map[string]any{"type": "string"},
				},
			},
			"elementNode": map[string]any{
				"type":                 "object",
				"required":             []string{"kind", "tag"},
				"additionalProperties": true,
				"properties": map[string]any{
					"kind":     map[string]any{"const": "element"},
					"tag":      map[string]any{"type": "string", "minLength": 1},
					"attrs":    map[string]any{"type": "object", "additionalProperties": true},
					"children": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/$defs/widgetNode"}},
				},
			},
			"componentNode": map[string]any{
				"type":                 "object",
				"required":             []string{"kind", "type"},
				"additionalProperties": true,
				"properties": map[string]any{
					"kind":     map[string]any{"const": "component"},
					"type":     map[string]any{"type": "string", "minLength": 1},
					"props":    map[string]any{"type": "object", "additionalProperties": true},
					"children": map[string]any{"type": "array", "items": map[string]any{"$ref": "#/$defs/widgetNode"}},
				},
			},
		},
	}
}
