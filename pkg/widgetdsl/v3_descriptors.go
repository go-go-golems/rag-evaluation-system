package widgetdsl

import "fmt"

type v3ModuleDescriptor struct {
	Exports    []v3ExportDescriptor
	Namespaces []v3NamespaceDescriptor
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

// widgetV3Module describes every direct public export installed by widget.dsl.
// Nested builder methods, intent methods, and action-context contracts are
// intentionally modeled separately from this direct-runtime parity inventory.
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
			Members:        v3Members([]string{"server", "navigate", "download", "event", "copy"}),
		},
		{
			ExportName:     "bind",
			TypeName:       "BindingNamespace",
			Description:    "Accessor and constant binding builders.",
			RuntimeFactory: "bindingObject",
			Members:        v3Members([]string{"field", "path", "map", "template", "context", "const"}),
		},
		{
			ExportName:     "ui",
			TypeName:       "UINamespace",
			Description:    "Generic composition widgets.",
			RuntimeFactory: "v3UIObject",
			Members: v3Members([]string{
				"callout", "stack", "inline", "splitPane", "card", "button", "caption", "badge", "metadata",
				"shareLink", "form", "formRow", "textInput", "textareaInput", "selectInput", "status", "emptyState",
			}),
		},
		{
			ExportName:     "data",
			TypeName:       "DataNamespace",
			Description:    "Schema, collection, matrix, selection, cell, and item helpers.",
			RuntimeFactory: "v3DataObject",
			Members:        v3Members([]string{"fields", "collection", "selection", "item", "matrix"}, "cell"),
		},
		{
			ExportName:     "crm",
			TypeName:       "CrmNamespace",
			Description:    "CRM field schemas, pipelines, records, activities, tasks, and actions.",
			RuntimeFactory: "v3CRMObject",
			Members:        v3Members([]string{"fields", "pipeline", "pipelineBoard", "recordFields", "activityFeed", "tasksInbox", "stat", "funnel"}, "intent"),
			Views: []v3ViewDescriptor{
				{Name: "pipelineBoard", Signature: "pipelineBoard(pipeline: Record<string, any>, deals: Record<string, any>[], configure?: Fragment<CrmPipelineBoardBuilder>): WidgetNodeSpec", Component: "BoardEngine", Description: "Render an opportunity pipeline board."},
				{Name: "recordFields", Signature: "recordFields(values: Record<string, JsonValue>, fields: Record<string, any>, configure?: Fragment<CrmRecordFieldsBuilder>): WidgetNodeSpec", Component: "RecordFieldList", Description: "Render typed CRM fields."},
				{Name: "activityFeed", Signature: "activityFeed(activities: Record<string, any>[], configure?: Fragment<CrmActivityFeedBuilder>): WidgetNodeSpec", Component: "ActivityFeed", Description: "Render a CRM activity timeline."},
			},
		},
		{
			ExportName:     "cms",
			TypeName:       "CmsNamespace",
			Description:    "CMS media, queue, and markdown editor helpers.",
			RuntimeFactory: "v3CMSObject",
			Members:        v3Members([]string{"mediaLibrary", "articleQueue", "markdownEditor"}, "intent"),
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
	return out
}
