package widgetdsl

import "fmt"

type v3ModuleDescriptor struct {
	Exports          []v3ExportDescriptor
	Namespaces       []v3NamespaceDescriptor
	NestedNamespaces []v3NestedNamespaceDescriptor
	Builders         []v3BuilderDescriptor
	ActionContexts   []v3ActionContextDescriptor
}

type v3ExportDescriptor struct {
	Name        string
	Signature   string
	Description string
}

type v3NamespaceDescriptor struct {
	ExportName     string
	TypeName       string
	Description    string
	RuntimeFactory string
	Members        []v3MemberDescriptor
	Views          []v3ViewDescriptor
}

type v3MemberDescriptor struct {
	Name string
	Kind string
}

type v3NestedNamespaceDescriptor struct {
	Path     string
	TypeName string
	Members  []v3MemberDescriptor
}

type v3BuilderDescriptor struct {
	TypeName string
	Methods  []string
}

type v3ActionContextDescriptor struct {
	Name        string
	Component   string
	Fields      []string
	Description string
}

type v3ViewDescriptor struct {
	Name        string
	Signature   string
	Component   string
	Description string
}

func v3Members(functions []string, objects ...string) []v3MemberDescriptor {
	members := make([]v3MemberDescriptor, 0, len(functions)+len(objects))
	for _, name := range functions {
		members = append(members, v3MemberDescriptor{Name: name, Kind: "function"})
	}
	for _, name := range objects {
		members = append(members, v3MemberDescriptor{Name: name, Kind: "object"})
	}
	return members
}

func v3Builder(typeName string, methods ...string) v3BuilderDescriptor {
	return v3BuilderDescriptor{TypeName: typeName, Methods: append(methods, "use")}
}

// widgetV3Module describes every public export and composable builder installed
// by widget.dsl. Runtime parity tests execute the actual Goja factories so this
// inventory cannot silently drift from JavaScript behavior.
var widgetV3Module = v3ModuleDescriptor{
	Exports: []v3ExportDescriptor{
		{
			Name:        "page",
			Signature:   "page(titleOrOptions: string | Record<string, any>, configure?: Fragment<PageBuilder>): PageBuilder",
			Description: "Create a page builder.",
		},
	},
	Namespaces: []v3NamespaceDescriptor{
		{
			ExportName:     "raw",
			TypeName:       "RawNamespace",
			Description:    "Low-level escape hatches for text, element, component, and fragment nodes.",
			RuntimeFactory: "v3RawObject",
			Members:        v3Members([]string{"text", "element", "component", "fragment"}),
		},
		{
			ExportName:     "act",
			TypeName:       "ActionNamespace",
			Description:    "Generic action builders.",
			RuntimeFactory: "actionObject",
			Members:        v3Members([]string{"server", "navigate", "download", "event", "copy", "openOverlay", "closeOverlay"}),
		},
		{
			ExportName:     "bind",
			TypeName:       "BindingNamespace",
			Description:    "Accessor and constant binding builders.",
			RuntimeFactory: "bindingObject",
			Members:        v3Members([]string{"field", "path", "map", "template", "context", "const"}),
		},
		{
			ExportName:     "app",
			TypeName:       "AppNamespace",
			Description:    "Typed application shell and viewport ownership helpers.",
			RuntimeFactory: "v3AppObject",
			Members:        v3Members([]string{"shell", "none", "rootOwned"}),
		},
		{
			ExportName:     "ui",
			TypeName:       "UINamespace",
			Description:    "Generic composition widgets.",
			RuntimeFactory: "v3UIObject",
			Members: v3Members([]string{
				"callout", "stack", "inline", "splitPane", "card", "button", "caption", "badge", "metadata",
				"shareLink", "form", "formRow", "textInput", "textareaInput", "selectInput", "status", "emptyState",
				"text", "code", "divider", "disclosure", "scroll", "tabs", "summary", "checkList", "stepList", "markdownArticle", "upload", "formDialog",
			}),
		},
		{
			ExportName:     "data",
			TypeName:       "DataNamespace",
			Description:    "Schema, collection, matrix, selection, cell, and item helpers.",
			RuntimeFactory: "v3DataObject",
			Members:        v3Members([]string{"fields", "collection", "selection", "item", "matrix", "activityFeed"}, "cell"),
		},
		{
			ExportName:     "crm",
			TypeName:       "CrmNamespace",
			Description:    "CRM field schemas, pipelines, records, activities, tasks, and actions.",
			RuntimeFactory: "v3CRMObject",
			Members:        v3Members([]string{"fields", "pipeline", "pipelineBoard", "recordFields", "field", "tasksInbox", "stat", "funnel"}, "intent"),
			Views: []v3ViewDescriptor{
				{Name: "pipelineBoard", Signature: "pipelineBoard(pipeline: Record<string, any>, deals: Record<string, any>[], configure?: Fragment<CrmPipelineBoardBuilder>): WidgetNodeSpec", Component: "BoardEngine", Description: "Render an opportunity pipeline board."},
				{Name: "recordFields", Signature: "recordFields(values: Record<string, JsonValue>, fields: Record<string, any>, configure?: Fragment<CrmRecordFieldsBuilder>): WidgetNodeSpec", Component: "RecordFieldList", Description: "Render typed CRM fields."},
			},
		},
		{
			ExportName:     "cms",
			TypeName:       "CmsNamespace",
			Description:    "CMS media, queue, and markdown editor helpers.",
			RuntimeFactory: "v3CMSObject",
			Members:        v3Members([]string{"shell", "mediaLibrary", "articleQueue", "markdownEditor"}, "intent"),
		},
		{
			ExportName:     "course",
			TypeName:       "CourseNamespace",
			Description:    "Course shell, landing, slide, handout, metadata, agenda, and material helpers.",
			RuntimeFactory: "v3CourseObject",
			Members:        v3Members([]string{"shell", "landing", "slideDeck", "handouts", "metadataForm", "agendaEditor", "materialUploads"}, "intent"),
		},
		{
			ExportName:     "context",
			TypeName:       "ContextNamespace",
			Description:    "Context style, diagram, transcript workspace, and intent helpers.",
			RuntimeFactory: "v3ContextObject",
			Members:        v3Members([]string{"styleSet", "palette", "diagram", "workspace"}, "intent"),
		},
		{
			ExportName:     "schedule",
			TypeName:       "ScheduleNamespace",
			Description:    "Availability poll, poll summary, booking picker, and schedule intent helpers.",
			RuntimeFactory: "v3ScheduleObject",
			Members:        v3Members([]string{"availabilityPoll", "pollSummary", "bookingPicker"}, "intent"),
			Views: []v3ViewDescriptor{
				{Name: "availabilityPoll", Signature: "availabilityPoll(poll: AvailabilityPoll, configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec", Component: "MatrixGrid", Description: "Render respondent availability against poll options."},
				{Name: "pollSummary", Signature: "pollSummary(poll: AvailabilityPoll, tallies: PollTally[], configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec", Component: "MatrixGrid", Description: "Render aggregate option tallies."},
				{Name: "bookingPicker", Signature: "bookingPicker(availability: Record<string, any>, configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec", Component: "MatrixGrid", Description: "Render bookable resources by slot."},
			},
		},
		{
			ExportName:     "time",
			TypeName:       "TimeNamespace",
			Description:    "Month, week, formatting, range, and time intent helpers.",
			RuntimeFactory: "v3TimeObject",
			Members:        v3Members([]string{"month", "week", "format", "formatRange", "slotLabel"}, "range", "intent"),
			Views: []v3ViewDescriptor{
				{Name: "month", Signature: "month(eventsOrMarkers: CalendarEvent[] | Record<string, any>, configure?: Fragment<TimeMonthBuilder>): WidgetNodeSpec", Component: "MonthGrid", Description: "Render day markers for a month."},
				{Name: "week", Signature: "week(events: CalendarEvent[], configure?: Fragment<TimeWeekBuilder>): WidgetNodeSpec", Component: "TimeGrid", Description: "Render event blocks for a week; allDay is intentionally omitted."},
			},
		},
		{
			ExportName:     "style",
			TypeName:       "Record<string, any>",
			Description:    "Reserved style namespace.",
			RuntimeFactory: "NewObject",
		},
	},
	NestedNamespaces: []v3NestedNamespaceDescriptor{
		{Path: "data.cell", TypeName: "CellNamespace", Members: v3Members([]string{"field", "status", "template", "cycle", "value"})},
		{Path: "data.selection", TypeName: "SelectionNamespace", Members: v3Members([]string{"urlParam"})},
		{Path: "crm.intent", TypeName: "CrmIntentNamespace", Members: v3Members([]string{"openDeal", "moveDeal", "updateField", "completeTask"})},
		{Path: "cms.intent", TypeName: "CmsIntentNamespace", Members: v3Members([]string{"selectAsset", "openAsset", "uploadAssets", "selectArticle", "createArticle", "publishArticle", "archiveArticle", "previewArticle"})},
		{Path: "course.intent", TypeName: "CourseIntentNamespace", Members: v3Members([]string{"navigate", "selectHandout", "downloadHandout", "printHandout", "previousSlide", "nextSlide", "presentSlide", "editAgenda", "uploadMaterial", "deleteMaterial"})},
		{Path: "context.intent", TypeName: "ContextIntentNamespace", Members: v3Members([]string{"selectPart", "selectAnnotation"})},
		{Path: "schedule.intent", TypeName: "ScheduleIntentNamespace", Members: v3Members([]string{"toggleAvailability", "submitResponse"})},
		{Path: "time.range", TypeName: "TimeRangeNamespace", Members: v3Members([]string{"week"})},
		{Path: "time.intent", TypeName: "TimeIntentNamespace", Members: v3Members([]string{"selectDay", "selectEvent"})},
	},
	Builders: []v3BuilderDescriptor{
		v3Builder("PageBuilder", "id", "title", "meta", "shell", "root", "density", "breadcrumb", "section", "view", "validate", "toPage"),
		v3Builder("AppShellBuilder", "brand", "navigation", "content"),
		v3Builder("NavigationBuilder", "placement", "active", "width", "narrowMode", "ariaLabel", "section"),
		v3Builder("NavigationItemsBuilder", "item"),
		v3Builder("ContentViewportBuilder", "maxWidth", "padding", "scroll"),
		v3Builder("SectionBuilder", "caption", "anchor", "tone", "text", "view", "slot", "actions", "metric", "metadata"),
		v3Builder("ActionsBuilder", "add", "button"),
		v3Builder("FormDialogBuilder", "title", "body", "initialFocus", "submitLabel", "cancelLabel", "submit"),
		v3Builder("FieldSetBuilder", "key", "primary", "short", "prose", "count", "status", "date", "currency", "media", "url", "build", "validate"),
		v3Builder("CollectionBuilder", "id", "schema", "empty", "select", "search", "paginate", "table", "edit", "masterDetail", "validate", "toNode", "toIR"),
		v3Builder("SearchBuilder", "value", "query", "placeholder", "resultCount", "submit", "clear"),
		v3Builder("PaginationBuilder", "current", "size", "total", "sizes", "position", "onChange"),
		v3Builder("TableBuilder", "className", "rowSelect", "actionColumn", "keyboard", "command", "styleWhen"),
		v3Builder("TableKeyboardBuilder", "mode", "selection", "vimAliases", "enterSelect"),
		v3Builder("RowCommandBuilder", "key", "label", "danger", "action"),
		v3Builder("EditorBuilder", "create", "submit", "submitPost", "reorder", "remove", "actions"),
		v3Builder("MatrixBuilder", "id", "columns", "column", "valueAt", "cell", "onCellAction", "toNode"),
		v3Builder("SchedulePollBuilder", "styleSet", "readOnly", "editableRow", "selectedCell", "onToggle", "ariaLabel"),
		v3Builder("TimeMonthBuilder", "styleSet", "selected", "today", "weekStartsOn", "onSelect"),
		v3Builder("TimeWeekBuilder", "styleSet", "range", "hours", "hourHeight", "viewportHeight", "now", "selected", "onSelect", "onSlotCreate"),
		v3Builder("ContextStyleSetBuilder", "style", "legend"),
		v3Builder("ContextDiagramBuilder", "styleSet", "palette", "view", "selected", "legend", "empty", "onSelect"),
		v3Builder("ContextWorkspaceBuilder", "selectedAnnotation", "showNotes", "styleSet", "message", "annotation", "empty", "onAnnotationSelect"),
		v3Builder("CourseShellBuilder", "active", "subtitle", "contentPadding", "main", "footer", "onNavigate"),
		v3Builder("CourseLandingBuilder", "activeAgenda", "onAgendaSelect", "onPrimary", "onSecondary"),
		v3Builder("CourseSlideDeckBuilder", "mode", "visualSide", "onPrevious", "onNext", "onPresent", "onFullscreen"),
		v3Builder("CourseHandoutsBuilder", "selected", "title", "empty", "onSelect", "onDownload", "onPrint"),
		v3Builder("CourseMetadataFormBuilder", "title", "onSubmit"),
		v3Builder("CourseMaterialUploadsBuilder", "accept", "onUpload", "onDelete"),
		v3Builder("CmsShellBuilder", "active", "subtitle", "contentPadding", "main", "header", "footer", "onNavigate"),
		v3Builder("CmsMediaLibraryBuilder", "selection", "selected", "query", "kindFilter", "page", "empty", "accept", "asset", "details", "toolbar", "onSelect", "onOpen", "onUpload"),
		v3Builder("CmsArticleQueueBuilder", "selected", "status", "query", "page", "empty", "row", "rowActions", "filters", "onSelect", "onCreate", "onRowAction", "onPublish", "onArchive", "onPreview"),
		v3Builder("CmsMarkdownEditorBuilder", "title", "placeholder", "onChange", "onSubmit"),
		v3Builder("CrmFieldsBuilder", "text", "longtext", "email", "phone", "url", "number", "currency", "percent", "date", "datetime", "boolean", "select", "multiselect", "tags", "user", "relation", "build", "validate"),
		v3Builder("CrmPipelineBuilder", "stage", "build", "validate"),
		v3Builder("CrmPipelineBoardBuilder", "summaries", "selected", "ariaLabel", "onMove", "onOpen"),
		v3Builder("CrmRecordFieldsBuilder", "mode", "refs", "onChange"),
		v3Builder("ActivityFeedBuilder", "groupByDay", "glyph", "glyphs", "styleSet", "onOpen", "onLoadMore"),
	},
	ActionContexts: []v3ActionContextDescriptor{
		{Name: "app.navigate", Component: "SidebarNav", Fields: []string{"itemId", "value", "componentType"}, Description: "Context dispatched by typed application navigation."},
		{Name: "table.rowSelect", Component: "DataTable", Fields: []string{"row", "rowKey", "componentType"}, Description: "Context dispatched when a collection row is selected."},
		{Name: "table.cellAction", Component: "DataTableCell", Fields: []string{"row", "rowKey", "componentType"}, Description: "Context dispatched by an action-button cell."},
		{Name: "matrix.cellAction", Component: "MatrixGrid", Fields: []string{"row", "column", "rowKey", "colId", "value", "componentType"}, Description: "Context dispatched when a matrix cell is activated."},
		{Name: "context.annotationSelect", Component: "TranscriptWorkspacePanel", Fields: []string{"annotationId", "value", "componentType"}, Description: "Context dispatched when a transcript annotation is selected."},
		{Name: "course.navigate", Component: "CourseStudioShell", Fields: []string{"itemId", "item", "value", "componentType"}, Description: "Context dispatched from course navigation."},
		{Name: "course.agendaSelect", Component: "CourseLessonPanel", Fields: []string{"agendaItemId", "value", "componentType"}, Description: "Context dispatched when an agenda item is selected."},
		{Name: "course.cta", Component: "CourseLessonPanel", Fields: []string{"cta", "componentType"}, Description: "Context dispatched by a primary or secondary course call to action."},
		{Name: "course.slideControl", Component: "CourseSlidePanel", Fields: []string{"value", "componentType"}, Description: "Context dispatched by previous, next, present, and fullscreen slide controls."},
		{Name: "course.handout", Component: "HandoutDocumentShell", Fields: []string{"documentId", "document", "value", "componentType"}, Description: "Context dispatched when a handout is selected, downloaded, or printed."},
		{Name: "upload.files", Component: "ContextUploadDropArea", Fields: []string{"files", "fileNames", "fileCount", "componentType"}, Description: "Serialized file context dispatched by a generic upload area."},
		{Name: "cms.asset", Component: "MediaLibraryPanel", Fields: []string{"assetId", "value", "componentType"}, Description: "Context dispatched when an asset is selected or opened."},
		{Name: "cms.assetQuery", Component: "MediaLibraryPanel", Fields: []string{"query", "value", "componentType"}, Description: "Context dispatched when an asset query is submitted."},
		{Name: "cms.assetKind", Component: "MediaLibraryPanel", Fields: []string{"kind", "value", "componentType"}, Description: "Context dispatched when the asset kind filter changes."},
		{Name: "cms.assetPage", Component: "MediaLibraryPanel", Fields: []string{"page", "value", "componentType"}, Description: "Context dispatched when the asset page changes."},
		{Name: "cms.assetUpload", Component: "MediaLibraryPanel", Fields: []string{"files", "fileNames", "fileCount", "componentType"}, Description: "Serialized file context dispatched by the media library."},
		{Name: "cms.article", Component: "ArticleListPanel", Fields: []string{"articleId", "value", "componentType"}, Description: "Context dispatched when an article is selected."},
		{Name: "cms.articleRowAction", Component: "ArticleListPanel", Fields: []string{"articleId", "rowAction", "value", "componentType"}, Description: "Context dispatched by an article row action."},
		{Name: "cms.articleStatus", Component: "ArticleListPanel", Fields: []string{"status", "value", "componentType"}, Description: "Context dispatched when the article status filter changes."},
		{Name: "cms.articleQuery", Component: "ArticleListPanel", Fields: []string{"query", "value", "componentType"}, Description: "Context dispatched when an article query is submitted."},
		{Name: "cms.articlePage", Component: "ArticleListPanel", Fields: []string{"page", "value", "componentType"}, Description: "Context dispatched when the article page changes."},
		{Name: "time.daySelect", Component: "MonthGrid", Fields: []string{"dateISO", "value", "componentType"}, Description: "Context dispatched when a calendar day is selected."},
		{Name: "time.blockSelect", Component: "TimeGrid", Fields: []string{"blockId", "value", "componentType"}, Description: "Context dispatched when a time block is selected."},
		{Name: "time.slotCreate", Component: "TimeGrid", Fields: []string{"dayISO", "hour", "value", "componentType"}, Description: "Context dispatched when an empty time slot is activated."},
		{Name: "crm.boardMove", Component: "BoardEngine", Fields: []string{"cardId", "from", "to", "beforeId", "componentType"}, Description: "Context dispatched when a CRM card moves."},
		{Name: "crm.boardOpen", Component: "BoardEngine", Fields: []string{"cardId", "componentType"}, Description: "Context dispatched when a CRM card is opened."},
		{Name: "crm.fieldChange", Component: "RecordFieldList", Fields: []string{"key", "value", "componentType"}, Description: "Context dispatched when a CRM record field changes."},
		{Name: "activity.open", Component: "ActivityFeed", Fields: []string{"activityId", "componentType"}, Description: "Context dispatched when an activity is opened."},
		{Name: "activity.loadMore", Component: "ActivityFeed", Fields: []string{"componentType"}, Description: "Context dispatched when earlier activities are requested."},
	},
}

func widgetV3DescriptorTypeScriptLines() []string {
	lines := make([]string, 0, len(widgetV3Module.Exports)+len(widgetV3Module.Namespaces))
	for _, export := range widgetV3Module.Exports {
		lines = append(lines, fmt.Sprintf("export function %s;", export.Signature))
	}
	for _, namespace := range widgetV3Module.Namespaces {
		lines = append(lines, fmt.Sprintf("export const %s: %s;", namespace.ExportName, namespace.TypeName))
	}
	return lines
}

func WidgetV3APIReferenceMarkdown() string {
	out := "# widget.dsl API Reference\n\nGenerated from widgetV3Module.\n\n"
	for _, export := range widgetV3Module.Exports {
		out += fmt.Sprintf("## `%s`\n\n`%s`. %s\n\n", export.Name, export.Signature, export.Description)
	}
	for _, namespace := range widgetV3Module.Namespaces {
		out += fmt.Sprintf("## `%s` — %s\n\n%s\n\n", namespace.ExportName, namespace.TypeName, namespace.Description)
		if namespace.RuntimeFactory != "" {
			out += fmt.Sprintf("Runtime factory: `%s`.\n\n", namespace.RuntimeFactory)
		}
		views := make(map[string]v3ViewDescriptor, len(namespace.Views))
		for _, view := range namespace.Views {
			views[view.Name] = view
		}
		for _, member := range namespace.Members {
			if view, ok := views[member.Name]; ok {
				out += fmt.Sprintf("- `%s`: `%s` → `%s`. %s\n", view.Name, view.Signature, view.Component, view.Description)
				continue
			}
			out += fmt.Sprintf("- `%s` (%s)\n", member.Name, member.Kind)
		}
		if len(namespace.Members) > 0 {
			out += "\n"
		}
	}
	out += "## Nested namespaces\n\n"
	for _, namespace := range widgetV3Module.NestedNamespaces {
		out += fmt.Sprintf("- `%s`: %v\n", namespace.Path, v3DescriptorMemberNamesForMarkdown(namespace.Members))
	}
	out += "\n## Composable builders\n\nAll builders below expose `use(fragment)` in addition to their listed methods.\n\n"
	for _, builder := range widgetV3Module.Builders {
		out += fmt.Sprintf("- `%s`: %v\n", builder.TypeName, builder.Methods)
	}
	out += "\n## Action contexts\n\n"
	for _, context := range widgetV3Module.ActionContexts {
		out += fmt.Sprintf("- `%s` (`%s`): `%v`. %s\n", context.Name, context.Component, context.Fields, context.Description)
	}
	out += "\n"
	return out
}

func v3DescriptorMemberNamesForMarkdown(members []v3MemberDescriptor) []string {
	names := make([]string, 0, len(members))
	for _, member := range members {
		names = append(names, member.Name)
	}
	return names
}
