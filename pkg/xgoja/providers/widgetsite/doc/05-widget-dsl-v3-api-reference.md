---
Title: "Widget DSL v3 API Reference"
Slug: widget-dsl-v3-api-reference
Short: "Descriptor-derived inventory of public widget.dsl v3 namespaces and domain views."
Topics:
- xgoja
- widget-dsl
- widget-ir
- javascript
IsTopLevel: false
IsTemplate: false
ShowPerDefault: true
SectionType: GeneralTopic
---

Generated from widgetV3NamespaceDescriptors.

## `raw` — RawNamespace

Low-level escape hatches for text, element, component, and fragment nodes.

Runtime factory: `v3RawObject`.

## `act` — ActionNamespace

Generic action builders.

Runtime factory: `actionObject`.

## `bind` — BindingNamespace

Accessor and constant binding builders.

Runtime factory: `bindingObject`.

## `ui` — UINamespace

Generic composition widgets.

Runtime factory: `v3UIObject`.

## `data` — DataNamespace

Schema, collection, matrix, selection, cell, and item helpers.

Runtime factory: `v3DataObject`.

## `crm` — CrmNamespace

CRM field schemas, pipelines, records, activities, tasks, and actions.

Runtime factory: `v3CRMObject`.

- `pipelineBoard`: `pipelineBoard(pipeline: Record<string, any>, deals: Record<string, any>[], configure?: Fragment<CrmPipelineBoardBuilder>): WidgetNodeSpec` → `BoardEngine`. Render an opportunity pipeline board.
- `recordFields`: `recordFields(values: Record<string, JsonValue>, fields: Record<string, any>, configure?: Fragment<CrmRecordFieldsBuilder>): WidgetNodeSpec` → `RecordFieldList`. Render typed CRM fields.
- `activityFeed`: `activityFeed(activities: Record<string, any>[], configure?: Fragment<CrmActivityFeedBuilder>): WidgetNodeSpec` → `ActivityFeed`. Render a CRM activity timeline.

## `cms` — CmsNamespace

CMS media, queue, and markdown editor helpers.

Runtime factory: `v3CMSObject`.

## `course` — CourseNamespace

Course shell, landing, slide, handout, metadata, agenda, and material helpers.

Runtime factory: `v3CourseObject`.

## `context` — ContextNamespace

Context style, diagram, transcript workspace, and intent helpers.

Runtime factory: `v3ContextObject`.

## `schedule` — ScheduleNamespace

Availability poll, poll summary, booking picker, and schedule intent helpers.

Runtime factory: `v3ScheduleObject`.

- `availabilityPoll`: `availabilityPoll(poll: AvailabilityPoll, configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec` → `MatrixGrid`. Render respondent availability against poll options.
- `pollSummary`: `pollSummary(poll: AvailabilityPoll, tallies: PollTally[], configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec` → `MatrixGrid`. Render aggregate option tallies.
- `bookingPicker`: `bookingPicker(availability: Record<string, any>, configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec` → `MatrixGrid`. Render bookable resources by slot.

## `time` — TimeNamespace

Month, week, formatting, range, and time intent helpers.

Runtime factory: `v3TimeObject`.

- `month`: `month(eventsOrMarkers: CalendarEvent[] | Record<string, any>, configure?: Fragment<TimeMonthBuilder>): WidgetNodeSpec` → `MonthGrid`. Render day markers for a month.
- `week`: `week(events: CalendarEvent[], configure?: Fragment<TimeWeekBuilder>): WidgetNodeSpec` → `TimeGrid`. Render event blocks for a week; allDay is intentionally omitted.

## `style` — Record<string, any>

Reserved style namespace.

Runtime factory: `NewObject`.


## Using this reference

The descriptor inventory names stable public namespaces and the domain views that lower into React components. It intentionally does not replace the generated TypeScript declarations: use the declarations for parameter-level completion, and use the examples tutorial for complete page and action compositions.

## Troubleshooting

The descriptor file is the source of truth for this snapshot. Use these checks when a runtime change, declaration, and help output appear to disagree.

| Problem | Cause | Solution |
| --- | --- | --- |
| A namespace is missing from this page. | The runtime, descriptor inventory, or help snapshot changed independently. | Update `pkg/widgetdsl/v3_descriptors.go`, regenerate this help body, and run the descriptor/help test. |
| A helper appears in JavaScript completion but not here. | The descriptor inventory intentionally summarizes namespaces and selected domain views. | Check the generated TypeScript declarations and add a descriptor view when the helper needs standalone discovery. |
| A component lacks a typed helper. | Its typed v3 surface has not been added yet. | Use `widget.raw.component(...)` only as a narrow escape hatch, then add the missing typed helper. |

## See Also

The following entries provide the parameter-level contracts and end-to-end code that this compact namespace inventory deliberately omits.

- `widget-dsl-v3-examples` — runnable composition, action, scheduling, and CRM patterns.
- `widget-dsl-js-api-reference` — detailed action contracts and legacy module reference.
- `widget-dsl-getting-started` — configure an xgoja host for v3.
