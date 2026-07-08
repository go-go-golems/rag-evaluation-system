---
title: Widget DSL v3 API Reference
doc_type: reference
topics:
  - widget-dsl
  - generated-api-reference
---

# widget.dsl API Reference

Generated from `pkg/widgetdsl/v3_descriptors.go` (`widgetV3NamespaceDescriptors`). The descriptor inventory is now the source used by declaration tests and this reference for namespace exports.

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
