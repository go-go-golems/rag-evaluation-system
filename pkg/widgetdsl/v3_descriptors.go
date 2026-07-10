package widgetdsl

import "fmt"

type v3NamespaceDescriptor struct {
	ExportName     string
	TypeName       string
	Description    string
	RuntimeFactory string
	Views          []v3ViewDescriptor
}

type v3ViewDescriptor struct {
	Name        string
	Signature   string
	Component   string
	Description string
}

// widgetV3NamespaceDescriptors is the first descriptor-backed source of truth
// for the public widget.dsl namespace surface. Phase 9 intentionally starts with
// declarations and API reference generation so future namespace work has a
// single inventory to extend instead of hand-editing export lists in tests/docs.
var widgetV3NamespaceDescriptors = []v3NamespaceDescriptor{
	{ExportName: "raw", TypeName: "RawNamespace", Description: "Low-level escape hatches for text, element, component, and fragment nodes.", RuntimeFactory: "v3RawObject"},
	{ExportName: "act", TypeName: "ActionNamespace", Description: "Generic action builders.", RuntimeFactory: "actionObject"},
	{ExportName: "bind", TypeName: "BindingNamespace", Description: "Accessor and constant binding builders.", RuntimeFactory: "bindingObject"},
	{ExportName: "ui", TypeName: "UINamespace", Description: "Generic composition widgets.", RuntimeFactory: "v3UIObject"},
	{ExportName: "data", TypeName: "DataNamespace", Description: "Schema, collection, matrix, selection, cell, and item helpers.", RuntimeFactory: "v3DataObject"},
	{ExportName: "crm", TypeName: "CrmNamespace", Description: "CRM field schemas, pipelines, records, activities, tasks, and actions.", RuntimeFactory: "v3CRMObject",
		Views: []v3ViewDescriptor{
			{Name: "pipelineBoard", Signature: "pipelineBoard(pipeline: Record<string, any>, deals: Record<string, any>[], configure?: Fragment<CrmPipelineBoardBuilder>): WidgetNodeSpec", Component: "BoardEngine", Description: "Render an opportunity pipeline board."},
			{Name: "recordFields", Signature: "recordFields(values: Record<string, JsonValue>, fields: Record<string, any>, configure?: Fragment<CrmRecordFieldsBuilder>): WidgetNodeSpec", Component: "RecordFieldList", Description: "Render typed CRM fields."},
			{Name: "activityFeed", Signature: "activityFeed(activities: Record<string, any>[], configure?: Fragment<CrmActivityFeedBuilder>): WidgetNodeSpec", Component: "ActivityFeed", Description: "Render a CRM activity timeline."},
		},
	},
	{ExportName: "cms", TypeName: "CmsNamespace", Description: "CMS media, queue, and markdown editor helpers.", RuntimeFactory: "v3CMSObject"},
	{ExportName: "course", TypeName: "CourseNamespace", Description: "Course shell, landing, slide, handout, metadata, agenda, and material helpers.", RuntimeFactory: "v3CourseObject"},
	{ExportName: "context", TypeName: "ContextNamespace", Description: "Context style, diagram, transcript workspace, and intent helpers.", RuntimeFactory: "v3ContextObject"},
	{
		ExportName: "schedule", TypeName: "ScheduleNamespace", Description: "Availability poll, poll summary, booking picker, and schedule intent helpers.", RuntimeFactory: "v3ScheduleObject",
		Views: []v3ViewDescriptor{
			{Name: "availabilityPoll", Signature: "availabilityPoll(poll: AvailabilityPoll, configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec", Component: "MatrixGrid", Description: "Render respondent availability against poll options."},
			{Name: "pollSummary", Signature: "pollSummary(poll: AvailabilityPoll, tallies: PollTally[], configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec", Component: "MatrixGrid", Description: "Render aggregate option tallies."},
			{Name: "bookingPicker", Signature: "bookingPicker(availability: Record<string, any>, configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec", Component: "MatrixGrid", Description: "Render bookable resources by slot."},
		},
	},
	{
		ExportName: "time", TypeName: "TimeNamespace", Description: "Month, week, formatting, range, and time intent helpers.", RuntimeFactory: "v3TimeObject",
		Views: []v3ViewDescriptor{
			{Name: "month", Signature: "month(eventsOrMarkers: CalendarEvent[] | Record<string, any>, configure?: Fragment<TimeMonthBuilder>): WidgetNodeSpec", Component: "MonthGrid", Description: "Render day markers for a month."},
			{Name: "week", Signature: "week(events: CalendarEvent[], configure?: Fragment<TimeWeekBuilder>): WidgetNodeSpec", Component: "TimeGrid", Description: "Render event blocks for a week; allDay is intentionally omitted."},
		},
	},
	{ExportName: "style", TypeName: "Record<string, any>", Description: "Reserved style namespace.", RuntimeFactory: "NewObject"},
}

func widgetV3DescriptorTypeScriptLines() []string {
	lines := make([]string, 0, len(widgetV3NamespaceDescriptors))
	for _, namespace := range widgetV3NamespaceDescriptors {
		lines = append(lines, fmt.Sprintf("export const %s: %s;", namespace.ExportName, namespace.TypeName))
	}
	return lines
}

func WidgetV3APIReferenceMarkdown() string {
	out := "# widget.dsl API Reference\n\nGenerated from widgetV3NamespaceDescriptors.\n\n"
	for _, namespace := range widgetV3NamespaceDescriptors {
		out += fmt.Sprintf("## `%s` — %s\n\n%s\n\n", namespace.ExportName, namespace.TypeName, namespace.Description)
		if namespace.RuntimeFactory != "" {
			out += fmt.Sprintf("Runtime factory: `%s`.\n\n", namespace.RuntimeFactory)
		}
		for _, view := range namespace.Views {
			out += fmt.Sprintf("- `%s`: `%s` → `%s`. %s\n", view.Name, view.Signature, view.Component, view.Description)
		}
		if len(namespace.Views) > 0 {
			out += "\n"
		}
	}
	return out
}
