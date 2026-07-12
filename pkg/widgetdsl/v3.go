package widgetdsl

import (
	"fmt"
	"strings"
	"time"

	"github.com/dop251/goja"
	widgetspec "github.com/go-go-golems/rag-evaluation-system/pkg/widgetdsl/spec"
)

type v3PageSpec struct {
	SchemaVersion string
	ID            string
	Title         string
	Meta          map[string]any
	Shell         any
	Density       string
	Breadcrumbs   []map[string]any
	Sections      []v3SectionSpec
}

type v3SectionSpec struct {
	Title    any
	Caption  string
	AnchorID string
	Tone     string
	Actions  []any
	Children []v3NodeSpec
}

type v3NodeSpec struct {
	Kind   string
	IR     map[string]any
	Source *v3SourceSpan
}

type v3SourceSpan struct {
	File   string
	Line   int
	Column int
}

type v3SlotSpec struct {
	Function goja.Value
	Fallback goja.Value
}

type v3SelectionSpec struct {
	Mode     string
	KeyField string
	Selected any
}

type v3ListItemSpec struct {
	ID       string
	Label    any
	Icon     any
	Badge    any
	Disabled bool
	Extra    map[string]any
}

func (r *runtime) v3Page(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl page(titleOrOptions, configure?) requires a title string or options object")))
	}
	spec := &v3PageSpec{SchemaVersion: "0.1.0", ID: "page", Title: "Page", Meta: map[string]any{}}
	first := call.Arguments[0]
	if isPlainObject(first) && !looksLikeWidgetNodeExport(first) {
		options := exportObject(first)
		spec.ID = stringFromMap(options, "id", spec.ID)
		spec.Title = stringFromMap(options, "title", spec.Title)
		spec.SchemaVersion = stringFromMap(options, "schemaVersion", spec.SchemaVersion)
		if meta, ok := options["meta"].(map[string]any); ok {
			spec.Meta = meta
		}
	} else {
		title := strings.TrimSpace(first.String())
		if title == "" {
			panic(r.vm.NewGoError(fmt.Errorf("widget.dsl page title must not be empty")))
		}
		spec.Title = title
		spec.ID = slugID(title)
	}
	builder := r.v3PageBuilder(spec)
	if len(call.Arguments) > 1 {
		r.applyV3BuilderCallback(builder, call.Arguments[1], "page")
	}
	return builder
}

func (r *runtime) v3ScheduleObject() *goja.Object {
	schedule := r.vm.NewObject()
	setExport(schedule, "availabilityPoll", r.v3ScheduleAvailabilityPoll)
	setExport(schedule, "pollSummary", r.v3SchedulePollSummary)
	setExport(schedule, "bookingPicker", r.v3ScheduleBookingPicker)
	setExport(schedule, "intent", r.v3ScheduleIntentObject())
	return schedule
}

func (r *runtime) v3ScheduleIntentObject() *goja.Object {
	intent := r.vm.NewObject()
	setExport(intent, "toggleAvailability", func(rowID goja.Value, optionID goja.Value, value ...goja.Value) map[string]any {
		payload := map[string]any{"responseId": rowID.Export(), "optionId": optionID.Export()}
		if len(value) > 0 {
			payload["value"] = value[0].Export()
		}
		return map[string]any{"kind": "event", "event": "schedule.availability.toggle", "detail": payload}
	})
	setExport(intent, "submitResponse", func(response goja.Value) map[string]any {
		return map[string]any{"kind": "event", "event": "schedule.availability.submit", "detail": map[string]any{"response": response.Export()}}
	})
	return intent
}

func (r *runtime) v3ScheduleAvailabilityPoll(poll goja.Value, cb ...goja.Value) map[string]any {
	p := exportObject(poll)
	props := map[string]any{
		"rows":         anySlice(valueOrDefault(p["responses"], p["rows"])),
		"columns":      schedulePollColumns(p),
		"valueAt":      map[string]any{"mapField": stringFromMap(p, "valuesField", "availability")},
		"cell":         map[string]any{"kind": "cycle", "states": []any{"available", "maybe", "unavailable"}},
		"rowHeader":    map[string]any{"kind": "field", "field": stringFromMap(p, "rowLabelField", "name")},
		"getRowKey":    map[string]any{"field": stringFromMap(p, "rowKeyField", "id")},
		"stickyHeader": true,
		"ariaLabel":    stringFromMap(p, "title", "Availability poll"),
	}
	copyIfPresent(props, p, "styleSet")
	builder := r.v3SchedulePollBuilder(props)
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "schedule.availabilityPoll")
	}
	if props["readOnly"] == true {
		delete(props, "onCellAction")
		delete(props, "editableRowKey")
	}
	delete(props, "readOnly")
	return componentNode("MatrixGrid", props)
}

func (r *runtime) v3SchedulePollSummary(poll goja.Value, tallies goja.Value, cb ...goja.Value) map[string]any {
	p := exportObject(poll)
	props := map[string]any{
		"rows":         anySlice(tallies.Export()),
		"columns":      schedulePollColumns(p),
		"valueAt":      map[string]any{"mapField": stringFromMap(p, "tallyField", "counts")},
		"cell":         map[string]any{"kind": "value"},
		"rowHeader":    map[string]any{"kind": "field", "field": stringFromMap(p, "summaryLabelField", "label")},
		"getRowKey":    map[string]any{"field": stringFromMap(p, "summaryKeyField", "id")},
		"stickyHeader": true,
		"ariaLabel":    stringFromMap(p, "summaryTitle", "Availability poll summary"),
	}
	builder := r.v3SchedulePollBuilder(props)
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "schedule.pollSummary")
	}
	return componentNode("MatrixGrid", props)
}

func (r *runtime) v3ScheduleBookingPicker(availability goja.Value, cb ...goja.Value) map[string]any {
	options := exportObject(availability)
	props := map[string]any{
		"rows":         anySlice(valueOrDefault(options["resources"], options["rows"])),
		"columns":      scheduleOptionColumns(anySlice(valueOrDefault(options["slots"], options["options"]))),
		"valueAt":      map[string]any{"mapField": stringFromMap(options, "valuesField", "availability")},
		"cell":         map[string]any{"kind": "cycle", "states": []any{"available", "selected", "unavailable"}},
		"rowHeader":    map[string]any{"kind": "field", "field": stringFromMap(options, "rowLabelField", "label")},
		"getRowKey":    map[string]any{"field": stringFromMap(options, "rowKeyField", "id")},
		"stickyHeader": true,
		"ariaLabel":    stringFromMap(options, "title", "Booking picker"),
	}
	builder := r.v3SchedulePollBuilder(props)
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "schedule.bookingPicker")
	}
	return componentNode("MatrixGrid", props)
}

func (r *runtime) v3SchedulePollBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("schedule.poll")
	setExport(obj, "styleSet", func(styleSet goja.Value) *goja.Object { props["styleSet"] = styleSet.Export(); return obj })
	setExport(obj, "readOnly", func(readOnly ...bool) *goja.Object { props["readOnly"] = len(readOnly) == 0 || readOnly[0]; return obj })
	setExport(obj, "editableRow", func(rowKey string) *goja.Object { props["editableRowKey"] = rowKey; return obj })
	setExport(obj, "selectedCell", func(rowKey string, colID string) *goja.Object {
		props["selectedCell"] = map[string]any{"rowKey": rowKey, "colId": colID}
		return obj
	})
	setExport(obj, "onToggle", func(action goja.Value) *goja.Object { props["onCellAction"] = action.Export(); return obj })
	setExport(obj, "ariaLabel", func(label string) *goja.Object { props["ariaLabel"] = label; return obj })
	return obj
}

func schedulePollColumns(poll map[string]any) []any {
	return scheduleOptionColumns(anySlice(valueOrDefault(poll["options"], poll["columns"])))
}

func scheduleOptionColumns(options []any) []any {
	columns := []any{}
	for _, raw := range options {
		option, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		id := stringFromMap(option, "id", stringFromMap(option, "dateISO", stringFromMap(option, "startISO", "")))
		if id == "" {
			continue
		}
		header := valueOrDefault(option["label"], valueOrDefault(option["title"], id))
		columns = append(columns, map[string]any{"id": id, "header": header, "meta": option})
	}
	return columns
}

func (r *runtime) v3TimeObject() *goja.Object {
	timeObj := r.vm.NewObject()
	setExport(timeObj, "month", r.v3TimeMonth)
	setExport(timeObj, "week", r.v3TimeWeek)
	setExport(timeObj, "format", v3TimeFormat)
	setExport(timeObj, "formatRange", v3TimeFormatRange)
	setExport(timeObj, "slotLabel", v3TimeSlotLabel)
	setExport(timeObj, "range", r.v3TimeRangeObject())
	setExport(timeObj, "intent", r.v3TimeIntentObject())
	return timeObj
}

func (r *runtime) v3TimeRangeObject() *goja.Object {
	rangeObj := r.vm.NewObject()
	setExport(rangeObj, "week", v3TimeRangeWeek)
	return rangeObj
}

func (r *runtime) v3TimeIntentObject() *goja.Object {
	intent := r.vm.NewObject()
	setExport(intent, "selectDay", func(dayISO goja.Value) map[string]any {
		return map[string]any{"kind": "event", "event": "time.day.select", "detail": map[string]any{"dayISO": dayISO.Export()}}
	})
	setExport(intent, "selectEvent", func(eventID goja.Value) map[string]any {
		return map[string]any{"kind": "event", "event": "time.event.select", "detail": map[string]any{"eventId": eventID.Export()}}
	})
	return intent
}

func (r *runtime) v3TimeMonth(eventsOrMarkers goja.Value, cb ...goja.Value) map[string]any {
	input := eventsOrMarkers.Export()
	props := map[string]any{"monthISO": currentMonthISO(), "markers": map[string]any{}, "showHeader": true}
	if options, ok := input.(map[string]any); ok {
		props["monthISO"] = stringFromMap(options, "monthISO", stringFromMap(options, "month", currentMonthISO()))
		if markers := options["markers"]; markers != nil {
			props["markers"] = markers
		} else {
			props["markers"] = monthMarkersFromEvents(anySlice(options["events"]))
		}
		copyIfPresent(props, options, "styleSet")
		copyIfPresent(props, options, "selectedDateISO")
		copyIfPresent(props, options, "todayISO")
		copyIfPresent(props, options, "minDateISO")
		copyIfPresent(props, options, "maxDateISO")
		copyIfPresent(props, options, "weekStartsOn")
	} else {
		props["markers"] = monthMarkersFromEvents(anySlice(input))
	}
	builder := r.v3TimeMonthBuilder(props)
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "time.month")
	}
	return componentNode("MonthGrid", props)
}

func (r *runtime) v3TimeMonthBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("time.month")
	setExport(obj, "styleSet", func(styleSet goja.Value) *goja.Object { props["styleSet"] = styleSet.Export(); return obj })
	setExport(obj, "selected", func(dayISO string) *goja.Object { props["selectedDateISO"] = dayISO; return obj })
	setExport(obj, "today", func(dayISO string) *goja.Object { props["todayISO"] = dayISO; return obj })
	setExport(obj, "weekStartsOn", func(day int) *goja.Object { props["weekStartsOn"] = day; return obj })
	setExport(obj, "onSelect", func(action goja.Value) *goja.Object { props["onDaySelectAction"] = action.Export(); return obj })
	return obj
}

func (r *runtime) v3TimeWeek(events goja.Value, cb ...goja.Value) map[string]any {
	items := anySlice(events.Export())
	props := map[string]any{
		"days":      weekDaysFromEvents(items),
		"blocks":    timeBlocksFromEvents(items),
		"styleSet":  r.v3ContextStyleSet(),
		"hourStart": 8,
		"hourEnd":   18,
	}
	builder := r.v3TimeWeekBuilder(props)
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "time.week")
	}
	return componentNode("TimeGrid", props)
}

func (r *runtime) v3TimeWeekBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("time.week")
	setExport(obj, "styleSet", func(styleSet goja.Value) *goja.Object { props["styleSet"] = styleSet.Export(); return obj })
	setExport(obj, "range", func(rangeSpec goja.Value) *goja.Object {
		props["days"] = weekDaysFromRange(exportObject(rangeSpec))
		return obj
	})
	setExport(obj, "hours", func(start int, end int) *goja.Object { props["hourStart"] = start; props["hourEnd"] = end; return obj })
	setExport(obj, "hourHeight", func(height int) *goja.Object { props["hourHeight"] = height; return obj })
	setExport(obj, "viewportHeight", func(height int) *goja.Object {
		props["style"] = map[string]any{"maxHeight": fmt.Sprintf("%dpx", height)}
		return obj
	})
	setExport(obj, "now", func(nowISO string) *goja.Object { props["nowISO"] = nowISO; return obj })
	setExport(obj, "selected", func(id string) *goja.Object { props["selectedBlockId"] = id; return obj })
	setExport(obj, "onSelect", func(action goja.Value) *goja.Object { props["onBlockSelectAction"] = action.Export(); return obj })
	setExport(obj, "onSlotCreate", func(action goja.Value) *goja.Object { props["onSlotCreateAction"] = action.Export(); return obj })
	return obj
}

func monthMarkersFromEvents(events []any) map[string]any {
	markers := map[string]any{}
	for _, raw := range events {
		event, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		day := eventDayISO(event)
		if day == "" {
			continue
		}
		marker := map[string]any{"count": 1, "label": valueOrDefault(event["label"], event["title"]), "styleKey": stringFromMap(event, "styleKey", "event")}
		if existing, ok := markers[day].(map[string]any); ok {
			if count, ok := existing["count"].(int); ok {
				existing["count"] = count + 1
			} else {
				existing["count"] = numberFromMap(existing, "count", 1) + 1
			}
			markers[day] = existing
		} else {
			markers[day] = marker
		}
	}
	return markers
}

func timeBlocksFromEvents(events []any) []any {
	blocks := []any{}
	for _, raw := range events {
		event, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		start := stringFromMap(event, "startISO", stringFromMap(event, "start", ""))
		end := stringFromMap(event, "endISO", stringFromMap(event, "end", ""))
		if start == "" || end == "" {
			continue
		}
		block := map[string]any{
			"id":       stringFromMap(event, "id", start),
			"dayISO":   eventDayISO(event),
			"startISO": start,
			"endISO":   end,
			"styleKey": stringFromMap(event, "styleKey", "event"),
			"label":    valueOrDefault(event["label"], valueOrDefault(event["title"], stringFromMap(event, "id", "Event"))),
		}
		if meta, ok := event["meta"].(map[string]any); ok {
			block["meta"] = meta
		}
		// Deliberately omit allDay until the frontend contract is fixed.
		blocks = append(blocks, block)
	}
	return blocks
}

func weekDaysFromEvents(events []any) []any {
	if len(events) == 0 {
		return weekDays(time.Now())
	}
	for _, raw := range events {
		if event, ok := raw.(map[string]any); ok {
			if parsed, ok := parseISODateLike(eventDayISO(event)); ok {
				return weekDays(parsed)
			}
		}
	}
	return weekDays(time.Now())
}

func weekDaysFromRange(rangeSpec map[string]any) []any {
	if start, ok := parseISODateLike(stringFromMap(rangeSpec, "startISO", stringFromMap(rangeSpec, "start", ""))); ok {
		return weekDays(start)
	}
	return weekDays(time.Now())
}

func weekDays(anchor time.Time) []any {
	start := anchor
	for start.Weekday() != time.Monday {
		start = start.AddDate(0, 0, -1)
	}
	days := []any{}
	for i := 0; i < 7; i++ {
		days = append(days, start.AddDate(0, 0, i).Format("2006-01-02"))
	}
	return days
}

func eventDayISO(event map[string]any) string {
	if day := stringFromMap(event, "dayISO", stringFromMap(event, "dateISO", "")); day != "" {
		return day
	}
	start := stringFromMap(event, "startISO", stringFromMap(event, "start", ""))
	if len(start) >= 10 {
		return start[:10]
	}
	return ""
}

func currentMonthISO() string { return time.Now().Format("2006-01") }

func v3TimeRangeWeek(anchorISO ...string) map[string]any {
	anchor := time.Now()
	if len(anchorISO) > 0 {
		if parsed, ok := parseISODateLike(anchorISO[0]); ok {
			anchor = parsed
		}
	}
	days := weekDays(anchor)
	return map[string]any{"kind": "week", "startISO": days[0], "endISO": days[6], "days": days}
}

func v3TimeFormat(iso string, layout ...string) string {
	parsed, ok := parseISODateLike(iso)
	if !ok {
		return iso
	}
	if len(layout) > 0 && layout[0] != "" {
		return parsed.Format(layout[0])
	}
	return parsed.Format("Jan 2, 2006")
}

func v3TimeFormatRange(startISO string, endISO string, layout ...string) string {
	return v3TimeFormat(startISO, layout...) + " – " + v3TimeFormat(endISO, layout...)
}

func v3TimeSlotLabel(startISO string, endISO string) string {
	return v3TimeFormat(startISO, "15:04") + "–" + v3TimeFormat(endISO, "15:04")
}

func parseISODateLike(value string) (time.Time, bool) {
	if value == "" {
		return time.Time{}, false
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04", "2006-01-02"} {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, true
		}
	}
	return time.Time{}, false
}

func (r *runtime) v3ContextObject() *goja.Object {
	context := r.vm.NewObject()
	setExport(context, "styleSet", r.v3ContextStyleSet)
	setExport(context, "palette", r.v3ContextPalette)
	setExport(context, "diagram", r.v3ContextDiagram)
	setExport(context, "workspace", r.v3ContextWorkspace)
	setExport(context, "intent", r.v3ContextIntentObject())
	return context
}

func (r *runtime) v3ContextIntentObject() *goja.Object {
	intent := r.vm.NewObject()
	setExport(intent, "selectPart", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "event", "event": "context.part.select", "detail": map[string]any{"partId": id.Export()}}
	})
	setExport(intent, "selectAnnotation", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "event", "event": "context.annotation.select", "detail": map[string]any{"annotationId": id.Export()}}
	})
	return intent
}

func (r *runtime) v3ContextStyleSet(args ...goja.Value) map[string]any {
	styleSet := map[string]any{"legend": []any{}, "styles": map[string]any{}}
	if len(args) > 0 && isPlainObject(args[0]) {
		styleSet = exportObject(args[0])
	}
	if len(args) > 0 {
		if fn, ok := goja.AssertFunction(args[len(args)-1]); ok {
			builder := r.v3ContextStyleSetBuilder(styleSet)
			if _, err := fn(goja.Undefined(), builder); err != nil {
				panic(err)
			}
		}
	}
	if _, ok := styleSet["legend"]; !ok {
		styleSet["legend"] = []any{}
	}
	if _, ok := styleSet["styles"]; !ok {
		styleSet["styles"] = map[string]any{}
	}
	return styleSet
}

func (r *runtime) v3ContextStyleSetBuilder(styleSet map[string]any) *goja.Object {
	obj := r.newV3Builder("context.styleSet")
	setExport(obj, "style", func(id string, options goja.Value) *goja.Object {
		styles, _ := styleSet["styles"].(map[string]any)
		if styles == nil {
			styles = map[string]any{}
			styleSet["styles"] = styles
		}
		styles[id] = exportObject(options)
		return obj
	})
	setExport(obj, "legend", func(id string, label string, options ...goja.Value) *goja.Object {
		legend := anySlice(styleSet["legend"])
		item := map[string]any{"id": id, "label": label}
		mergeOptions(item, exportOptions(options))
		styleSet["legend"] = append(legend, item)
		return obj
	})
	return obj
}

func (r *runtime) v3ContextPalette(nameOrOptions goja.Value, entries ...goja.Value) map[string]any {
	options := map[string]any{}
	if isPlainObject(nameOrOptions) {
		options = exportObject(nameOrOptions)
	} else if nameOrOptions != nil && !goja.IsUndefined(nameOrOptions) && !goja.IsNull(nameOrOptions) {
		options["palette"] = nameOrOptions.String()
	}
	if len(entries) > 0 {
		options["entries"] = entries[0].Export()
	}
	return buildPaletteStyleSet(options)
}

func (r *runtime) v3ContextDiagram(snapshot goja.Value, cb ...goja.Value) map[string]any {
	props := map[string]any{"snapshot": valueOrDefault(snapshot.Export(), map[string]any{"id": "empty", "title": "Context", "limit": 0, "parts": []any{}})}
	builder := r.v3ContextDiagramBuilder(props)
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "context.diagram")
	}
	if props["styleSet"] == nil {
		props["styleSet"] = r.v3ContextStyleSet()
	}
	return componentNode("ContextDiagramPanel", props)
}

func (r *runtime) v3ContextDiagramBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("context.diagram")
	setExport(obj, "styleSet", func(styleSet goja.Value) *goja.Object { props["styleSet"] = styleSet.Export(); return obj })
	setExport(obj, "palette", func(nameOrOptions goja.Value, entries ...goja.Value) *goja.Object {
		props["styleSet"] = r.v3ContextPalette(nameOrOptions, entries...)
		return obj
	})
	setExport(obj, "view", func(view string) *goja.Object { props["initialView"] = view; return obj })
	setExport(obj, "selected", func(id string) *goja.Object { props["selectedPartId"] = id; return obj })
	setExport(obj, "legend", func(slot goja.Value) *goja.Object { props["legendSlot"] = r.v3SlotRef(slot); return obj })
	setExport(obj, "empty", func(slot goja.Value) *goja.Object { props["emptySlot"] = r.v3SlotRef(slot); return obj })
	setExport(obj, "onSelect", func(action goja.Value) *goja.Object { props["onPartSelectAction"] = action.Export(); return obj })
	return obj
}

func (r *runtime) v3ContextWorkspace(session goja.Value, cb ...goja.Value) map[string]any {
	s := exportObject(session)
	props := map[string]any{
		"title":       valueOrDefault(s["title"], "Transcript"),
		"subtitle":    s["subtitle"],
		"messages":    anySlice(s["messages"]),
		"annotations": anySlice(s["annotations"]),
		"showNotes":   true,
	}
	if snapshot := s["snapshot"]; snapshot != nil {
		props["snapshot"] = snapshot
	}
	builder := r.v3ContextWorkspaceBuilder(props)
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "context.workspace")
	}
	return componentNode("TranscriptWorkspacePanel", props)
}

func (r *runtime) v3ContextWorkspaceBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("context.workspace")
	setExport(obj, "selectedAnnotation", func(id string) *goja.Object { props["selectedAnnotationId"] = id; return obj })
	setExport(obj, "showNotes", func(show bool) *goja.Object { props["showNotes"] = show; return obj })
	setExport(obj, "styleSet", func(styleSet goja.Value) *goja.Object { props["styleSet"] = styleSet.Export(); return obj })
	setExport(obj, "message", func(slot goja.Value) *goja.Object { props["messageSlot"] = r.v3SlotRef(slot); return obj })
	setExport(obj, "annotation", func(slot goja.Value) *goja.Object { props["annotationSlot"] = r.v3SlotRef(slot); return obj })
	setExport(obj, "empty", func(slot goja.Value) *goja.Object { props["emptySlot"] = r.v3SlotRef(slot); return obj })
	setExport(obj, "onAnnotationSelect", func(action goja.Value) *goja.Object { props["onAnnotationSelectAction"] = action.Export(); return obj })
	return obj
}

func (r *runtime) v3CourseObject() *goja.Object {
	course := r.vm.NewObject()
	setExport(course, "shell", r.v3CourseShell)
	setExport(course, "landing", r.v3CourseLanding)
	setExport(course, "slideDeck", r.v3CourseSlideDeck)
	setExport(course, "handouts", r.v3CourseHandouts)
	setExport(course, "metadataForm", r.v3CourseMetadataForm)
	setExport(course, "agendaEditor", r.v3CourseAgendaEditor)
	setExport(course, "materialUploads", r.v3CourseMaterialUploads)
	setExport(course, "intent", r.v3CourseIntentObject())
	return course
}

func (r *runtime) v3CourseIntentObject() *goja.Object {
	intent := r.vm.NewObject()
	setExport(intent, "navigate", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "navigate", "to": "?item=" + v3URLTemplateValue(id)}
	})
	setExport(intent, "selectHandout", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "event", "event": "course.handout.select", "detail": map[string]any{"documentId": id.Export()}}
	})
	setExport(intent, "downloadHandout", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "download", "to": "/handouts/" + v3URLTemplateValue(id)}
	})
	setExport(intent, "printHandout", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "event", "event": "course.handout.print", "detail": map[string]any{"documentId": id.Export()}}
	})
	setExport(intent, "previousSlide", func() map[string]any { return map[string]any{"kind": "event", "event": "course.slide.previous"} })
	setExport(intent, "nextSlide", func() map[string]any { return map[string]any{"kind": "event", "event": "course.slide.next"} })
	setExport(intent, "presentSlide", func() map[string]any { return map[string]any{"kind": "event", "event": "course.slide.present"} })
	setExport(intent, "editAgenda", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "server", "name": "course.agenda.edit", "payload": map[string]any{"id": id.Export()}}
	})
	setExport(intent, "uploadMaterial", func() map[string]any { return map[string]any{"kind": "server", "name": "course.material.upload"} })
	setExport(intent, "deleteMaterial", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "server", "name": "course.material.delete", "payload": map[string]any{"id": id.Export()}}
	})
	return intent
}

func (r *runtime) v3CourseShell(definition goja.Value, cb ...goja.Value) map[string]any {
	def := exportObject(definition)
	props := map[string]any{"sections": anySlice(valueOrDefault(def["sections"], []any{})), "title": valueOrDefault(def["title"], "Course")}
	copyIfPresent(props, def, "subtitle")
	builder := r.v3CourseShellBuilder(props)
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "course.shell")
	}
	children := []any{}
	if main, ok := widgetNodeFromAny(props["main"]); ok {
		children = append(children, main)
		delete(props, "main")
	}
	return componentNode("CourseStudioShell", props, children...)
}

func (r *runtime) v3CourseShellBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("course.shell")
	setExport(obj, "active", func(id string) *goja.Object { props["activeItemId"] = id; return obj })
	setExport(obj, "subtitle", func(value goja.Value) *goja.Object { props["subtitle"] = r.v3Renderable(value); return obj })
	setExport(obj, "contentPadding", func(value string) *goja.Object { props["contentPadding"] = value; return obj })
	setExport(obj, "main", func(node goja.Value) *goja.Object { props["main"] = r.v3Renderable(node); return obj })
	setExport(obj, "footer", func(node goja.Value) *goja.Object { props["sidebarFooter"] = r.v3Renderable(node); return obj })
	setExport(obj, "onNavigate", func(action goja.Value) *goja.Object { props["onNavigateAction"] = action.Export(); return obj })
	return obj
}

func (r *runtime) v3CourseLanding(definition goja.Value, cb ...goja.Value) map[string]any {
	def := exportObject(definition)
	props := map[string]any{"course": def}
	if len(cb) > 0 {
		r.applyV3BuilderCallback(r.v3CourseLandingBuilder(props), cb[0], "course.landing")
	}
	return componentNode("CourseLessonPanel", props)
}

func (r *runtime) v3CourseLandingBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("course.landing")
	setExport(obj, "activeAgenda", func(id string) *goja.Object { props["activeAgendaItemId"] = id; return obj })
	setExport(obj, "onAgendaSelect", func(action goja.Value) *goja.Object { props["onAgendaItemSelectAction"] = action.Export(); return obj })
	setExport(obj, "onPrimary", func(action goja.Value) *goja.Object { props["onPrimaryCtaAction"] = action.Export(); return obj })
	setExport(obj, "onSecondary", func(action goja.Value) *goja.Object { props["onSecondaryCtaAction"] = action.Export(); return obj })
	return obj
}

func (r *runtime) v3CourseSlideDeck(deck goja.Value, cb ...goja.Value) map[string]any {
	d := exportObject(deck)
	slides := anySlice(valueOrDefault(d["slides"], []any{}))
	index := 0
	if v, ok := d["index"].(int64); ok {
		index = int(v)
	}
	slide := map[string]any{}
	if len(slides) > index {
		if m, ok := slides[index].(map[string]any); ok {
			slide = m
		}
	}
	props := map[string]any{"slide": valueOrDefault(d["slide"], slide), "snapshot": valueOrDefault(d["snapshot"], map[string]any{"id": "empty", "title": "Context", "limit": 0, "parts": []any{}}), "index": index, "total": len(slides)}
	if len(cb) > 0 {
		r.applyV3BuilderCallback(r.v3CourseSlideBuilder(props), cb[0], "course.slideDeck")
	}
	return componentNode("CourseSlidePanel", props)
}

func (r *runtime) v3CourseSlideBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("course.slideDeck")
	setExport(obj, "mode", func(mode string) *goja.Object { props["mode"] = mode; return obj })
	setExport(obj, "visualSide", func(side string) *goja.Object { props["visualSide"] = side; return obj })
	setExport(obj, "onPrevious", func(action goja.Value) *goja.Object { props["onPreviousAction"] = action.Export(); return obj })
	setExport(obj, "onNext", func(action goja.Value) *goja.Object { props["onNextAction"] = action.Export(); return obj })
	setExport(obj, "onPresent", func(action goja.Value) *goja.Object { props["onPresentAction"] = action.Export(); return obj })
	setExport(obj, "onFullscreen", func(action goja.Value) *goja.Object { props["onFullscreenAction"] = action.Export(); return obj })
	return obj
}

func (r *runtime) v3CourseHandouts(bundle goja.Value, cb ...goja.Value) map[string]any {
	b := exportObject(bundle)
	props := map[string]any{"intro": valueOrDefault(b["intro"], "Handout"), "documents": anySlice(valueOrDefault(b["docs"], b["documents"]))}
	builder := r.v3CourseHandoutsBuilder(props)
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "course.handouts")
	}
	return componentNode("HandoutDocumentShell", props)
}

func (r *runtime) v3CourseHandoutsBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("course.handouts")
	setExport(obj, "selected", func(id string) *goja.Object { props["selectedDocumentId"] = id; return obj })
	setExport(obj, "title", func(title goja.Value) *goja.Object { props["title"] = r.v3Renderable(title); return obj })
	setExport(obj, "empty", func(message goja.Value) *goja.Object { props["emptyMessage"] = r.v3Renderable(message); return obj })
	setExport(obj, "onSelect", func(action goja.Value) *goja.Object { props["onDocumentSelectAction"] = action.Export(); return obj })
	setExport(obj, "onDownload", func(action goja.Value) *goja.Object { props["onDownloadAction"] = action.Export(); return obj })
	setExport(obj, "onPrint", func(action goja.Value) *goja.Object { props["onPrintAction"] = action.Export(); return obj })
	return obj
}

func (r *runtime) v3CourseMetadataForm(metadata goja.Value, cb ...goja.Value) map[string]any {
	props := map[string]any{"title": "Course metadata"}
	children := []any{r.v3MetadataNode(exportObject(metadata))}
	if len(cb) > 0 {
		r.applyV3BuilderCallback(r.v3CourseFormBuilder(props), cb[0], "course.metadataForm")
	}
	return componentNode("FormPanel", props, children...)
}

func (r *runtime) v3CourseFormBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("course.metadataForm")
	setExport(obj, "title", func(title string) *goja.Object { props["title"] = title; return obj })
	setExport(obj, "onSubmit", func(action goja.Value) *goja.Object { props["onSubmitAction"] = action.Export(); return obj })
	return obj
}

func (r *runtime) v3CourseAgendaEditor(items goja.Value, cb ...goja.Value) *goja.Object {
	args := append([]goja.Value{items}, cb...)
	return r.v3Collection(args...)
}

func (r *runtime) v3CourseMaterialUploads(material goja.Value, cb ...goja.Value) map[string]any {
	props := exportObject(material)
	if _, ok := props["title"]; !ok {
		props["title"] = "Course materials"
	}
	builder := r.newV3Builder("course.materialUploads")
	setExport(builder, "accept", func(list goja.Value) *goja.Object { props["accept"] = anySlice(list.Export()); return builder })
	setExport(builder, "onUpload", func(action goja.Value) *goja.Object { props["onFilesSelectedAction"] = action.Export(); return builder })
	setExport(builder, "onDelete", func(action goja.Value) *goja.Object { props["onDeleteAction"] = action.Export(); return builder })
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "course.materialUploads")
	}
	return componentNode("ContextUploadDropArea", props)
}

func (r *runtime) v3CMSObject() *goja.Object {
	cms := r.vm.NewObject()
	setExport(cms, "mediaLibrary", r.v3CMSMediaLibrary)
	setExport(cms, "articleQueue", r.v3CMSArticleQueue)
	setExport(cms, "markdownEditor", r.v3CMSMarkdownEditor)
	setExport(cms, "intent", r.v3CMSIntentObject())
	return cms
}

func (r *runtime) v3CMSIntentObject() *goja.Object {
	intent := r.vm.NewObject()
	setExport(intent, "selectAsset", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "event", "event": "cms.asset.select", "detail": map[string]any{"assetId": id.Export()}}
	})
	setExport(intent, "openAsset", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "event", "event": "cms.asset.open", "detail": map[string]any{"assetId": id.Export()}}
	})
	setExport(intent, "uploadAssets", func() map[string]any { return map[string]any{"kind": "server", "name": "cms.assets.upload"} })
	setExport(intent, "selectArticle", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "event", "event": "cms.article.select", "detail": map[string]any{"articleId": id.Export()}}
	})
	setExport(intent, "createArticle", func() map[string]any { return map[string]any{"kind": "event", "event": "cms.article.create"} })
	setExport(intent, "publishArticle", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "server", "name": "cms.article.publish", "payload": map[string]any{"articleId": id.Export()}}
	})
	setExport(intent, "archiveArticle", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "server", "name": "cms.article.archive", "payload": map[string]any{"articleId": id.Export()}}
	})
	setExport(intent, "previewArticle", func(id goja.Value) map[string]any {
		return map[string]any{"kind": "navigate", "to": "?article=" + v3URLTemplateValue(id) + "&preview=1"}
	})
	return intent
}

func (r *runtime) v3CMSMediaLibrary(assets goja.Value, cb ...goja.Value) map[string]any {
	props := map[string]any{"assets": anySlice(assets.Export())}
	builder := r.v3CMSMediaLibraryBuilder(props)
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "cms.mediaLibrary")
	}
	return componentNode("MediaLibraryPanel", props)
}

func (r *runtime) v3CMSMediaLibraryBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("cms.mediaLibrary")
	setExport(obj, "selection", func(mode string) *goja.Object { props["selectionMode"] = mode; return obj })
	setExport(obj, "selected", func(ids goja.Value) *goja.Object { props["selectedAssetIds"] = anySlice(ids.Export()); return obj })
	setExport(obj, "query", func(value string) *goja.Object { props["query"] = value; return obj })
	setExport(obj, "kindFilter", func(value string) *goja.Object { props["kindFilter"] = value; return obj })
	setExport(obj, "page", func(page int, pageCount int) *goja.Object {
		props["page"] = page
		props["pageCount"] = pageCount
		return obj
	})
	setExport(obj, "empty", func(message string) *goja.Object { props["emptyMessage"] = message; return obj })
	setExport(obj, "accept", func(mimeList goja.Value) *goja.Object { props["accept"] = anySlice(mimeList.Export()); return obj })
	setExport(obj, "asset", func(slot goja.Value) *goja.Object { props["assetSlot"] = r.v3SlotRef(slot); return obj })
	setExport(obj, "details", func(slot goja.Value) *goja.Object { props["detailsSlot"] = r.v3SlotRef(slot); return obj })
	setExport(obj, "toolbar", func(cb goja.Value) *goja.Object {
		props["toolbar"] = r.callV3Slot(v3SlotSpec{Function: cb}, props)
		return obj
	})
	setExport(obj, "onSelect", func(action goja.Value) *goja.Object { props["onAssetSelectAction"] = action.Export(); return obj })
	setExport(obj, "onOpen", func(action goja.Value) *goja.Object { props["onAssetOpenAction"] = action.Export(); return obj })
	setExport(obj, "onUpload", func(action goja.Value) *goja.Object { props["onFilesSelectedAction"] = action.Export(); return obj })
	return obj
}

func (r *runtime) v3CMSArticleQueue(articles goja.Value, cb ...goja.Value) map[string]any {
	props := map[string]any{"articles": anySlice(articles.Export())}
	builder := r.v3CMSArticleQueueBuilder(props)
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "cms.articleQueue")
	}
	return componentNode("ArticleListPanel", props)
}

func (r *runtime) v3CMSArticleQueueBuilder(props map[string]any) *goja.Object {
	obj := r.newV3Builder("cms.articleQueue")
	setExport(obj, "selected", func(id string) *goja.Object { props["selectedArticleId"] = id; return obj })
	setExport(obj, "status", func(status string) *goja.Object { props["statusFilter"] = status; return obj })
	setExport(obj, "query", func(query string) *goja.Object { props["query"] = query; return obj })
	setExport(obj, "page", func(page int, pageCount int) *goja.Object {
		props["page"] = page
		props["pageCount"] = pageCount
		return obj
	})
	setExport(obj, "empty", func(message string) *goja.Object { props["emptyMessage"] = message; return obj })
	setExport(obj, "row", func(slot goja.Value) *goja.Object { props["rowSlot"] = r.v3SlotRef(slot); return obj })
	setExport(obj, "rowActions", func(slot goja.Value) *goja.Object { props["rowActionsSlot"] = r.v3SlotRef(slot); return obj })
	setExport(obj, "filters", func(slot goja.Value) *goja.Object { props["filtersSlot"] = r.v3SlotRef(slot); return obj })
	setExport(obj, "onSelect", func(action goja.Value) *goja.Object { props["onArticleSelectAction"] = action.Export(); return obj })
	setExport(obj, "onCreate", func(action goja.Value) *goja.Object { props["onCreateAction"] = action.Export(); return obj })
	setExport(obj, "onRowAction", func(action goja.Value) *goja.Object { props["onRowActionAction"] = action.Export(); return obj })
	setExport(obj, "onPublish", func(action goja.Value) *goja.Object { props["onPublishAction"] = action.Export(); return obj })
	setExport(obj, "onArchive", func(action goja.Value) *goja.Object { props["onArchiveAction"] = action.Export(); return obj })
	setExport(obj, "onPreview", func(action goja.Value) *goja.Object { props["onPreviewAction"] = action.Export(); return obj })
	return obj
}

func (r *runtime) v3CMSMarkdownEditor(body goja.Value, cb ...goja.Value) map[string]any {
	props := map[string]any{"value": body.Export()}
	if len(cb) > 0 && !goja.IsUndefined(cb[0]) && !goja.IsNull(cb[0]) {
		builder := r.newV3Builder("cms.markdownEditor")
		setExport(builder, "title", func(title string) *goja.Object { props["title"] = title; return builder })
		setExport(builder, "placeholder", func(placeholder string) *goja.Object { props["placeholder"] = placeholder; return builder })
		setExport(builder, "onChange", func(action goja.Value) *goja.Object { props["onChangeAction"] = action.Export(); return builder })
		setExport(builder, "onSubmit", func(action goja.Value) *goja.Object { props["onSubmitAction"] = action.Export(); return builder })
		r.applyV3BuilderCallback(builder, cb[0], "cms.markdownEditor")
	}
	return componentNode("MarkdownEditor", props)
}

func (r *runtime) v3SlotRef(slot goja.Value) map[string]any {
	if slot == nil || goja.IsUndefined(slot) || goja.IsNull(slot) {
		return nil
	}
	return map[string]any{"kind": "slot", "registered": true}
}

func (r *runtime) v3UIObject() *goja.Object {
	ui := r.vm.NewObject()
	setExport(ui, "callout", r.v3ComponentFactory("Panel", map[string]any{"tone": "callout"}))
	setExport(ui, "stack", r.v3ComponentFactory("Stack", nil))
	setExport(ui, "inline", r.v3ComponentFactory("Inline", nil))
	setExport(ui, "splitPane", r.v3UISplitPane)
	setExport(ui, "card", r.v3ComponentFactory("Panel", nil))
	setExport(ui, "button", r.v3UIButton)
	setExport(ui, "caption", r.v3ComponentFactory("Caption", nil))
	setExport(ui, "badge", r.v3ComponentFactory("Tag", nil))
	setExport(ui, "metadata", r.v3UIMetadata)
	setExport(ui, "shareLink", r.v3UIShareLink)
	setExport(ui, "form", r.v3ComponentFactory("FormPanel", nil))
	setExport(ui, "formRow", r.v3UIFormRow)
	setExport(ui, "textInput", r.v3ComponentFactory("TextInput", map[string]any{"readOnly": false}))
	setExport(ui, "textareaInput", r.v3ComponentFactory("TextareaInput", map[string]any{"readOnly": false}))
	setExport(ui, "selectInput", r.v3ComponentFactory("SelectInput", nil))
	setExport(ui, "status", r.v3UIStatus)
	setExport(ui, "emptyState", r.v3UIEmptyState)
	return ui
}

func (r *runtime) v3DataObject() *goja.Object {
	data := r.vm.NewObject()
	setExport(data, "fields", r.v3Fields)
	setExport(data, "collection", r.v3Collection)
	setExport(data, "selection", r.v3Selection)
	if selection := data.Get("selection").ToObject(r.vm); selection != nil {
		setExport(selection, "urlParam", func(param string, value goja.Value) map[string]any {
			return map[string]any{"kind": "urlParam", "param": param, "value": stringifyValue(value)}
		})
	}
	setExport(data, "item", r.v3ListItem)
	setExport(data, "cell", r.v3CellObject())
	setExport(data, "matrix", r.v3Matrix)
	return data
}

func (r *runtime) v3ComponentFactory(componentType string, defaults map[string]any) func(goja.FunctionCall) goja.Value {
	return func(call goja.FunctionCall) goja.Value {
		props, childStart := propsAndChildStart(call.Arguments, 0)
		if len(defaults) > 0 {
			merged := map[string]any{}
			for key, value := range defaults {
				merged[key] = value
			}
			for key, value := range props {
				merged[key] = value
			}
			props = merged
		}
		return r.vm.ToValue(r.v3BuildComponent(componentType, props, call.Arguments[childStart:]))
	}
}

func (r *runtime) v3UIButton(label goja.Value, action goja.Value, options ...goja.Value) map[string]any {
	props := exportOptions(options)
	if action != nil && !goja.IsUndefined(action) && !goja.IsNull(action) {
		props["action"] = action.Export()
	}
	return componentNode("Button", props, r.v3NodeSpecsToIR(r.v3ExportChild(label))...)
}

func (r *runtime) v3UISplitPane(left goja.Value, right goja.Value, options ...goja.Value) map[string]any {
	props := exportOptions(options)
	props["left"] = r.v3Renderable(left)
	props["right"] = r.v3Renderable(right)
	return componentNode("SplitPane", props)
}

func (r *runtime) v3UIFormRow(label goja.Value, control goja.Value, options ...goja.Value) map[string]any {
	props := exportOptions(options)
	props["label"] = r.v3Renderable(label)
	props["control"] = control.Export()
	return componentNode("FormRow", props)
}

func (r *runtime) v3UIStatus(status goja.Value, value goja.Value, options ...goja.Value) map[string]any {
	props := exportOptions(options)
	props["status"] = status.String()
	if _, ok := props["icon"]; !ok {
		props["icon"] = true
	}
	return componentNode("StatusText", props, r.v3NodeSpecsToIR(r.v3ExportChild(value))...)
}

func (r *runtime) v3UIEmptyState(title goja.Value, description goja.Value, options ...goja.Value) map[string]any {
	props := exportOptions(options)
	props["title"] = r.v3Renderable(title)
	if description != nil && !goja.IsUndefined(description) && !goja.IsNull(description) {
		props["description"] = r.v3Renderable(description)
	}
	return componentNode("EmptyState", props)
}

func (r *runtime) v3UIMetadata(record goja.Value, options ...goja.Value) map[string]any {
	props := exportOptions(options)
	props["items"] = v3MetadataItems(exportObject(record))
	return componentNode("MetadataGrid", props)
}

func (r *runtime) v3UIShareLink(href goja.Value, options ...goja.Value) map[string]any {
	props := exportOptions(options)
	url := href.String()
	props["href"] = url
	if _, ok := props["copyAction"]; !ok {
		props["copyAction"] = map[string]any{"kind": "copy", "value": url}
	}
	return componentNode("ShareLink", props)
}

func (r *runtime) v3ActionsBuilder(actions *[]any) *goja.Object {
	obj := r.newV3Builder("actions")
	setExport(obj, "add", func(label goja.Value, action goja.Value, options ...goja.Value) *goja.Object {
		item := exportOptions(options)
		item["label"] = r.v3Renderable(label)
		if action != nil && !goja.IsUndefined(action) && !goja.IsNull(action) {
			item["action"] = action.Export()
		}
		*actions = append(*actions, item)
		return obj
	})
	setExport(obj, "button", obj.Get("add"))
	return obj
}

func (r *runtime) v3Fields(args ...goja.Value) *goja.Object {
	name := "fields"
	var cb goja.Value
	if len(args) > 0 {
		if fn, ok := goja.AssertFunction(args[0]); ok {
			_ = fn
			cb = args[0]
		} else if strings.TrimSpace(args[0].String()) != "" {
			name = args[0].String()
		}
	}
	if len(args) > 1 {
		cb = args[1]
	}
	schema := &widgetspec.SchemaSpec{Name: name}
	builder := r.v3FieldsBuilder(schema)
	if cb != nil && !goja.IsUndefined(cb) && !goja.IsNull(cb) {
		r.applyV3BuilderCallback(builder, cb, "data.fields")
	}
	return builder
}

func (r *runtime) v3FieldsBuilder(schema *widgetspec.SchemaSpec) *goja.Object {
	obj := r.newV3Builder("data.fields")
	r.attachV2Ref(obj, &v2Ref{kind: "schemaBuilder", schema: schema})
	addField := func(name string, field widgetspec.FieldSpec, options ...goja.Value) *goja.Object {
		if strings.TrimSpace(name) == "" {
			panic(r.vm.NewGoError(fmt.Errorf("widget.dsl data.fields field name must not be empty")))
		}
		field.Name = name
		opts := exportOptions(options)
		field.Label = stringFromMap(opts, "label", field.Label)
		field.Layout.Width = stringFromMap(opts, "width", field.Layout.Width)
		if required, ok := opts["required"].(bool); ok {
			field.Validation.Required = required
		}
		schema.Fields = append(schema.Fields, field)
		return obj
	}
	setExport(obj, "key", func(name string, options ...goja.Value) *goja.Object {
		return addField(name, v3Field(widgetspec.FieldKindString, widgetspec.FieldSemanticKey, "caption", widgetspec.EditorControlText), options...)
	})
	setExport(obj, "primary", func(name string, options ...goja.Value) *goja.Object {
		return addField(name, v3Field(widgetspec.FieldKindString, widgetspec.FieldSemanticPrimary, "field", widgetspec.EditorControlText), options...)
	})
	setExport(obj, "short", func(name string, options ...goja.Value) *goja.Object {
		return addField(name, v3Field(widgetspec.FieldKindString, widgetspec.FieldSemanticShort, "field", widgetspec.EditorControlText), options...)
	})
	setExport(obj, "prose", func(name string, options ...goja.Value) *goja.Object {
		field := v3Field(widgetspec.FieldKindString, widgetspec.FieldSemanticProse, "", widgetspec.EditorControlTextarea)
		field.Editor.Rows = 4
		field.Summary.Elide = true
		return addField(name, field, options...)
	})
	setExport(obj, "count", func(name string, options ...goja.Value) *goja.Object {
		return addField(name, v3Field(widgetspec.FieldKindNumber, widgetspec.FieldSemanticCount, "number", widgetspec.EditorControlText), options...)
	})
	setExport(obj, "status", func(name string, options ...goja.Value) *goja.Object {
		return addField(name, v3Field(widgetspec.FieldKindString, widgetspec.FieldSemanticStatus, "status", widgetspec.EditorControlText), options...)
	})
	setExport(obj, "date", func(name string, options ...goja.Value) *goja.Object {
		return addField(name, v3Field(widgetspec.FieldKindDate, widgetspec.FieldSemanticShort, "field", widgetspec.EditorControlText), options...)
	})
	setExport(obj, "currency", func(name string, options ...goja.Value) *goja.Object {
		return addField(name, v3Field(widgetspec.FieldKindNumber, widgetspec.FieldSemanticMeasure, "number", widgetspec.EditorControlText), options...)
	})
	setExport(obj, "media", func(name string, options ...goja.Value) *goja.Object {
		return addField(name, v3Field(widgetspec.FieldKindMedia, widgetspec.FieldSemanticShort, "field", widgetspec.EditorControlText), options...)
	})
	setExport(obj, "url", func(name string, options ...goja.Value) *goja.Object {
		return addField(name, v3Field(widgetspec.FieldKindURL, widgetspec.FieldSemanticShort, "link", widgetspec.EditorControlText), options...)
	})
	setExport(obj, "build", func() *goja.Object { built := *schema; return r.v2SchemaValue(&built) })
	setExport(obj, "validate", func() []map[string]any { return validationIssuesForJS(schema.Validate("fields")) })
	return obj
}

func v3Field(kind widgetspec.FieldKind, semantic widgetspec.FieldSemantic, cellKind string, control widgetspec.EditorControl) widgetspec.FieldSpec {
	return widgetspec.FieldSpec{Kind: kind, Semantic: semantic, Editor: widgetspec.EditorSpec{Control: control}, Summary: widgetspec.SummarySpec{CellKind: cellKind}}
}

func (r *runtime) v3Collection(args ...goja.Value) *goja.Object {
	if len(args) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl data.collection(rows, configure?) requires rows")))
	}
	name := "collection"
	rowsArg := args[0]
	var cb goja.Value
	if len(args) > 1 {
		if _, ok := goja.AssertFunction(args[1]); ok {
			cb = args[1]
		} else if strings.TrimSpace(args[0].String()) != "" {
			name = args[0].String()
			rowsArg = args[1]
			if len(args) > 2 {
				cb = args[2]
			}
		}
	}
	collection := &widgetspec.CollectionSpec{Name: name, Rows: v2Rows(rowsArg.Export()), Mode: widgetspec.CollectionModeShow, Arrangement: widgetspec.ArrangementSpec{Kind: widgetspec.ArrangementKindTable}}
	builder := r.v3CollectionBuilder(collection)
	if cb != nil && !goja.IsUndefined(cb) && !goja.IsNull(cb) {
		r.applyV3BuilderCallback(builder, cb, "data.collection")
	}
	return builder
}

func (r *runtime) v3CollectionBuilder(collection *widgetspec.CollectionSpec) *goja.Object {
	obj := r.newV3Builder("data.collection")
	r.attachV2Ref(obj, &v2Ref{kind: "collectionBuilder", collection: collection})
	setExport(obj, "id", func(name string) *goja.Object {
		if strings.TrimSpace(name) != "" {
			collection.Name = name
		}
		return obj
	})
	setExport(obj, "schema", func(schemaValue goja.Value) *goja.Object {
		collection.Schema = *r.mustV2Ref(schemaValue, "schema").schema
		return obj
	})
	setExport(obj, "empty", func(message string) *goja.Object { collection.Empty = message; return obj })
	setExport(obj, "select", func(selectionValue goja.Value) *goja.Object {
		collection.Selection = v3SelectionToV2(selectionValue.Export())
		return obj
	})
	setExport(obj, "table", func(args ...goja.Value) *goja.Object {
		collection.Arrangement = widgetspec.ArrangementSpec{Kind: widgetspec.ArrangementKindTable}
		if len(args) > 0 && !goja.IsUndefined(args[0]) && !goja.IsNull(args[0]) {
			r.applyV3BuilderCallback(r.v3TableBuilder(collection), args[0], "data.collection.table")
		}
		return obj
	})
	setExport(obj, "edit", func(args ...goja.Value) *goja.Object {
		collection.Mode = widgetspec.CollectionModeEdit
		if len(args) > 0 && !goja.IsUndefined(args[0]) && !goja.IsNull(args[0]) {
			r.applyV3BuilderCallback(r.v3EditorBuilder(collection), args[0], "data.collection.edit")
		}
		return obj
	})
	setExport(obj, "masterDetail", func(args ...goja.Value) *goja.Object {
		collection.Arrangement = widgetspec.ArrangementSpec{Kind: widgetspec.ArrangementKindMasterDetail}
		return obj
	})
	setExport(obj, "validate", func() []map[string]any { return validationIssuesForJS(collection.Validate("collection")) })
	setExport(obj, "toNode", func() any { return collection.ToNode().ToWidgetNode() })
	setExport(obj, "toIR", func() any { return collection.ToNode().ToWidgetNode() })
	return obj
}

func (r *runtime) v3TableBuilder(collection *widgetspec.CollectionSpec) *goja.Object {
	obj := r.newV3Builder("data.collection.table")
	setExport(obj, "className", func(className string) *goja.Object { collection.Table.ClassName = className; return obj })
	setExport(obj, "rowSelect", func(actionValue goja.Value) *goja.Object {
		action := v3ActionFromAny(actionValue.Export())
		collection.Table.RowSelect = &action
		return obj
	})
	setExport(obj, "actionColumn", func(id string, header string, label string, actionValue goja.Value, options ...goja.Value) *goja.Object {
		action := v3ActionFromAny(actionValue.Export())
		column := widgetspec.TableActionColumnSpec{ID: id, Header: header, Label: label, Action: action}
		column.MaxWidth = stringFromMap(exportOptions(options), "maxWidth", column.MaxWidth)
		collection.Table.ActionColumns = append(collection.Table.ActionColumns, column)
		return obj
	})
	return obj
}

func (r *runtime) v3EditorBuilder(collection *widgetspec.CollectionSpec) *goja.Object {
	obj := r.newV3Builder("data.collection.edit")
	setExport(obj, "create", func(value goja.Value) *goja.Object {
		label := "New item"
		if !goja.IsUndefined(value) && !goja.IsNull(value) {
			if isPlainObject(value) {
				label = stringFromMap(exportObject(value), "label", label)
			} else {
				label = value.String()
			}
		}
		collection.Actions.Create = &widgetspec.CreateActionSpec{Label: label}
		return obj
	})
	setExport(obj, "submit", func(formAction string) *goja.Object {
		collection.Actions.Submit = &widgetspec.SubmitSpec{FormAction: formAction, Method: "post"}
		return obj
	})
	setExport(obj, "submitPost", func(formAction string) *goja.Object {
		collection.Actions.Submit = &widgetspec.SubmitSpec{FormAction: formAction, Method: "post"}
		return obj
	})
	setExport(obj, "reorder", func(actionValue goja.Value) *goja.Object {
		action := v3ActionFromAny(actionValue.Export())
		collection.Actions.Reorder = &action
		return obj
	})
	setExport(obj, "remove", func(actionValue goja.Value) *goja.Object {
		action := v3ActionFromAny(actionValue.Export())
		collection.Actions.Remove = &action
		return obj
	})
	setExport(obj, "actions", func(cb goja.Value) *goja.Object {
		r.applyV3BuilderCallback(obj, cb, "data.collection.edit.actions")
		return obj
	})
	return obj
}

func (r *runtime) v3Selection(modeOrOptions goja.Value, options ...goja.Value) map[string]any {
	spec := v3SelectionSpec{Mode: "single"}
	if isPlainObject(modeOrOptions) {
		opts := exportObject(modeOrOptions)
		spec.Mode = stringFromMap(opts, "mode", spec.Mode)
		spec.KeyField = stringFromMap(opts, "keyField", spec.KeyField)
		spec.Selected = opts["selected"]
	} else if modeOrOptions != nil && !goja.IsUndefined(modeOrOptions) && !goja.IsNull(modeOrOptions) {
		spec.Mode = strings.TrimSpace(modeOrOptions.String())
		opts := exportOptions(options)
		spec.KeyField = stringFromMap(opts, "keyField", spec.KeyField)
		spec.Selected = opts["selected"]
	}
	if spec.Mode != "single" && spec.Mode != "multi" {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl data.selection mode must be single or multi")))
	}
	out := map[string]any{"kind": "selection", "mode": spec.Mode}
	if spec.KeyField != "" {
		out["keyField"] = spec.KeyField
	}
	if spec.Selected != nil {
		out["selected"] = spec.Selected
	}
	return out
}

func (r *runtime) v3CellObject() *goja.Object {
	cell := r.vm.NewObject()
	setExport(cell, "field", func(field string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "field", "field": field}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(cell, "status", func(field string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "status", "field": field}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(cell, "template", func(template string) map[string]any { return map[string]any{"kind": "template", "template": template} })
	setExport(cell, "cycle", func(field string, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "cycle", "field": field}
		mergeOptions(out, exportOptions(options))
		return out
	})
	setExport(cell, "value", func(value goja.Value, options ...goja.Value) map[string]any {
		out := map[string]any{"kind": "constant", "value": value.Export()}
		mergeOptions(out, exportOptions(options))
		return out
	})
	return cell
}

func (r *runtime) v3Matrix(rows goja.Value, cb ...goja.Value) *goja.Object {
	spec := map[string]any{"rows": rows.Export(), "columns": []any{}}
	builder := r.v3MatrixBuilder(spec)
	if len(cb) > 0 {
		r.applyV3BuilderCallback(builder, cb[0], "data.matrix")
	}
	return builder
}

func (r *runtime) v3MatrixBuilder(spec map[string]any) *goja.Object {
	obj := r.newV3Builder("data.matrix")
	setExport(obj, "id", func(id string) *goja.Object { spec["id"] = id; return obj })
	setExport(obj, "columns", func(columns goja.Value) *goja.Object { spec["columns"] = columns.Export(); return obj })
	setExport(obj, "column", func(id string, label goja.Value, options ...goja.Value) *goja.Object {
		column := exportOptions(options)
		column["id"] = id
		column["header"] = r.v3Renderable(label)
		spec["columns"] = append(anySlice(spec["columns"]), column)
		return obj
	})
	setExport(obj, "valueAt", func(accessor goja.Value) *goja.Object { spec["valueAt"] = accessor.Export(); return obj })
	setExport(obj, "cell", func(cell goja.Value) *goja.Object { spec["cell"] = cell.Export(); return obj })
	setExport(obj, "onCellAction", func(action goja.Value) *goja.Object { spec["onCellAction"] = action.Export(); return obj })
	setExport(obj, "toNode", func() map[string]any { return componentNode("MatrixGrid", spec) })
	return obj
}

func (r *runtime) v3ListItem(id string, label goja.Value, options ...goja.Value) map[string]any {
	if strings.TrimSpace(id) == "" {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl data.item id must not be empty")))
	}
	spec := v3ListItemSpec{ID: id, Label: r.v3Renderable(label), Extra: exportOptions(options)}
	out := map[string]any{"kind": "listItem", "id": spec.ID, "label": spec.Label}
	for key, value := range spec.Extra {
		out[key] = value
	}
	if spec.Icon != nil {
		out["icon"] = spec.Icon
	}
	if spec.Badge != nil {
		out["badge"] = spec.Badge
	}
	if spec.Disabled {
		out["disabled"] = true
	}
	return out
}

func (r *runtime) v3PageBuilder(spec *v3PageSpec) *goja.Object {
	obj := r.newV3Builder("page")
	setExport(obj, "id", func(id string) *goja.Object {
		if strings.TrimSpace(id) != "" {
			spec.ID = id
		}
		return obj
	})
	setExport(obj, "title", func(title string) *goja.Object {
		if strings.TrimSpace(title) != "" {
			spec.Title = title
		}
		return obj
	})
	setExport(obj, "meta", func(key string, value goja.Value) *goja.Object {
		if spec.Meta == nil {
			spec.Meta = map[string]any{}
		}
		spec.Meta[key] = value.Export()
		return obj
	})
	setExport(obj, "shell", func(shell goja.Value) *goja.Object {
		spec.Shell = shell.Export()
		return obj
	})
	setExport(obj, "density", func(density string) *goja.Object {
		spec.Density = density
		return obj
	})
	setExport(obj, "breadcrumb", func(label goja.Value, href ...string) *goja.Object {
		item := map[string]any{"label": r.v3Renderable(label)}
		if len(href) > 0 && strings.TrimSpace(href[0]) != "" {
			item["href"] = href[0]
		}
		spec.Breadcrumbs = append(spec.Breadcrumbs, item)
		return obj
	})
	setExport(obj, "section", func(title goja.Value, cb ...goja.Value) *goja.Object {
		section := v3SectionSpec{Title: r.v3RenderableTitle(title)}
		sectionBuilder := r.v3SectionBuilder(&section)
		if len(cb) > 0 {
			r.applyV3BuilderCallback(sectionBuilder, cb[0], "section")
		}
		spec.Sections = append(spec.Sections, section)
		return obj
	})
	setExport(obj, "view", func(value goja.Value) *goja.Object {
		section := v3SectionSpec{Title: "Content", Children: r.v3ExportChild(value)}
		spec.Sections = append(spec.Sections, section)
		return obj
	})
	setExport(obj, "validate", func() []map[string]any {
		return v3PageValidationIssues(spec)
	})
	setExport(obj, "toPage", func() map[string]any {
		issues := v3PageValidationIssues(spec)
		if len(issues) > 0 {
			panic(r.vm.NewGoError(fmt.Errorf("widget.dsl page is invalid: %s", issues[0]["message"])))
		}
		return r.v3PageToIR(spec)
	})
	return obj
}

func (r *runtime) v3SectionBuilder(spec *v3SectionSpec) *goja.Object {
	obj := r.newV3Builder("section")
	setExport(obj, "caption", func(caption string) *goja.Object {
		spec.Caption = caption
		return obj
	})
	setExport(obj, "anchor", func(anchor string) *goja.Object {
		spec.AnchorID = anchor
		return obj
	})
	setExport(obj, "tone", func(tone string) *goja.Object {
		spec.Tone = tone
		return obj
	})
	setExport(obj, "text", func(value goja.Value) *goja.Object {
		spec.Children = append(spec.Children, r.v3TextNode(value))
		return obj
	})
	setExport(obj, "view", func(value goja.Value) *goja.Object {
		spec.Children = append(spec.Children, r.v3ExportChild(value)...)
		return obj
	})
	setExport(obj, "slot", func(context goja.Value, slot goja.Value, fallback ...goja.Value) *goja.Object {
		var fallbackSlot goja.Value
		if len(fallback) > 0 {
			fallbackSlot = fallback[0]
		}
		nodes := r.callV3Slot(v3SlotSpec{Function: slot, Fallback: fallbackSlot}, context.Export())
		spec.Children = append(spec.Children, nodes...)
		return obj
	})
	setExport(obj, "actions", func(cb goja.Value) *goja.Object {
		actions := r.v3ActionsBuilder(&spec.Actions)
		r.applyV3BuilderCallback(actions, cb, "section.actions")
		return obj
	})
	setExport(obj, "metric", func(label goja.Value, value goja.Value, options ...goja.Value) *goja.Object {
		props := exportOptions(options)
		props["key"] = r.v3Renderable(label)
		props["label"] = r.v3Renderable(label)
		props["value"] = r.v3Renderable(value)
		spec.Children = append(spec.Children, v3NodeSpecFromIR(componentNode("KeyValueStrip", map[string]any{"items": []any{props}})))
		return obj
	})
	setExport(obj, "metadata", func(record goja.Value) *goja.Object {
		spec.Children = append(spec.Children, v3NodeSpecFromIR(r.v3MetadataNode(exportObject(record))))
		return obj
	})
	return obj
}

func (r *runtime) newV3Builder(path string) *goja.Object {
	builder := r.vm.NewObject()
	setExport(builder, "use", func(fragment goja.Value) *goja.Object {
		r.applyV3BuilderCallback(builder, fragment, path+".use")
		return builder
	})
	return builder
}

func (r *runtime) applyV3BuilderCallback(builder *goja.Object, cb goja.Value, name string) {
	if cb == nil || goja.IsUndefined(cb) || goja.IsNull(cb) {
		return
	}
	fn, ok := goja.AssertFunction(cb)
	if !ok {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl %s callback must be a function", name)))
	}
	if _, err := fn(goja.Undefined(), builder); err != nil {
		panic(err)
	}
}

func (r *runtime) callV3Slot(slot v3SlotSpec, ctx any) []v3NodeSpec {
	return r.callV3SlotFunction(slot.Function, ctx, func(any) []v3NodeSpec {
		return r.callV3SlotFunction(slot.Fallback, ctx, nil)
	})
}

func (r *runtime) callV3SlotFunction(slot goja.Value, ctx any, fallback func(any) []v3NodeSpec) []v3NodeSpec {
	if slot == nil || goja.IsUndefined(slot) || goja.IsNull(slot) {
		if fallback == nil {
			return nil
		}
		return fallback(ctx)
	}
	fn, ok := goja.AssertFunction(slot)
	if !ok {
		panic(r.vm.NewGoError(fmt.Errorf("widget.dsl slot must be a function")))
	}
	value, err := fn(goja.Undefined(), r.vm.ToValue(ctx), r.v3SlotHelpers())
	if err != nil {
		panic(err)
	}
	if isV3EmptySlotResult(value) && fallback != nil {
		return fallback(ctx)
	}
	return r.v3ExportChild(value)
}

func (r *runtime) v3SlotHelpers() *goja.Object {
	h := r.vm.NewObject()
	setExport(h, "text", func(value goja.Value) map[string]any {
		return r.v3TextNode(value).toIR()
	})
	setExport(h, "caption", func(value goja.Value, options ...goja.Value) map[string]any {
		props := exportOptions(options)
		return componentNode("Caption", props, r.v3NodeSpecsToIR(r.v3ExportChild(value))...)
	})
	setExport(h, "strong", func(call goja.FunctionCall) goja.Value {
		return r.vm.ToValue(r.v3BuildElement("strong", map[string]any{}, call.Arguments))
	})
	setExport(h, "stack", func(call goja.FunctionCall) goja.Value {
		props, childStart := propsAndChildStart(call.Arguments, 0)
		return r.vm.ToValue(r.v3BuildComponent("Stack", props, call.Arguments[childStart:]))
	})
	setExport(h, "inline", func(call goja.FunctionCall) goja.Value {
		props, childStart := propsAndChildStart(call.Arguments, 0)
		return r.vm.ToValue(r.v3BuildComponent("Inline", props, call.Arguments[childStart:]))
	})
	setExport(h, "card", func(call goja.FunctionCall) goja.Value {
		props, childStart := propsAndChildStart(call.Arguments, 0)
		return r.vm.ToValue(r.v3BuildComponent("Panel", props, call.Arguments[childStart:]))
	})
	setExport(h, "button", func(label goja.Value, action goja.Value, options ...goja.Value) map[string]any {
		props := exportOptions(options)
		if action != nil && !goja.IsUndefined(action) && !goja.IsNull(action) {
			props["action"] = action.Export()
		}
		return componentNode("Button", props, r.v3NodeSpecsToIR(r.v3ExportChild(label))...)
	})
	setExport(h, "badge", func(value goja.Value, options ...goja.Value) map[string]any {
		props := exportOptions(options)
		return componentNode("Tag", props, r.v3NodeSpecsToIR(r.v3ExportChild(value))...)
	})
	setExport(h, "raw", r.v3RawObject())
	return h
}

func (r *runtime) v3RawObject() *goja.Object {
	raw := r.vm.NewObject()
	setExport(raw, "text", func(value goja.Value) map[string]any {
		return r.v3TextNode(value).toIR()
	})
	setExport(raw, "element", r.v3Element)
	setExport(raw, "component", r.v3Component)
	setExport(raw, "fragment", r.v3Fragment)
	return raw
}

func (r *runtime) v3Element(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("widget DSL element(tag, attrs?, ...children) requires a tag")))
	}
	tag := strings.TrimSpace(call.Arguments[0].String())
	if tag == "" {
		panic(r.vm.NewGoError(fmt.Errorf("widget DSL element tag must not be empty")))
	}
	attrs := map[string]any{}
	childStart := 1
	if len(call.Arguments) > 1 && isPlainObject(call.Arguments[1]) && !looksLikeWidgetNodeExport(call.Arguments[1]) {
		attrs = exportObject(call.Arguments[1])
		childStart = 2
	}
	return r.vm.ToValue(r.v3BuildElement(tag, attrs, call.Arguments[childStart:]))
}

func (r *runtime) v3Component(call goja.FunctionCall) goja.Value {
	if len(call.Arguments) == 0 {
		panic(r.vm.NewGoError(fmt.Errorf("widget DSL component(type, props?, ...children) requires a type")))
	}
	componentType := strings.TrimSpace(call.Arguments[0].String())
	if componentType == "" {
		panic(r.vm.NewGoError(fmt.Errorf("widget DSL component type must not be empty")))
	}
	props, childStart := propsAndChildStart(call.Arguments, 1)
	return r.vm.ToValue(r.v3BuildComponent(componentType, props, call.Arguments[childStart:]))
}

func (r *runtime) v3Fragment(call goja.FunctionCall) goja.Value {
	return r.vm.ToValue(r.v3NodeSpecsToIR(r.v3ExportChildren(call.Arguments)))
}

func (r *runtime) v3RenderableTitle(value goja.Value) any {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return "Section"
	}
	if exported, ok := value.Export().(bool); ok && !exported {
		return "Section"
	}
	if _, ok := value.Export().(string); ok {
		return value.String()
	}
	return r.v3Renderable(value)
}

func (r *runtime) v3Renderable(value goja.Value) any {
	nodes := r.v3ExportChild(value)
	if len(nodes) == 0 {
		return nil
	}
	if len(nodes) == 1 {
		return nodes[0].toIR()
	}
	return r.v3NodeSpecsToIR(nodes)
}

func (r *runtime) v3PageToIR(spec *v3PageSpec) map[string]any {
	children := make([]any, 0, len(spec.Sections)+1)
	if len(spec.Breadcrumbs) > 0 {
		children = append(children, componentNode("Breadcrumbs", map[string]any{"items": spec.Breadcrumbs}))
	}
	for _, section := range spec.Sections {
		children = append(children, r.v3SectionToNode(section))
	}
	rootProps := map[string]any{"gap": "lg"}
	if spec.Density != "" {
		rootProps["density"] = spec.Density
	}
	out := map[string]any{
		"schemaVersion": spec.SchemaVersion,
		"id":            spec.ID,
		"title":         spec.Title,
		"root":          componentNode("Stack", rootProps, children...),
	}
	if len(spec.Meta) > 0 {
		out["meta"] = spec.Meta
	}
	if spec.Shell != nil {
		out["shell"] = spec.Shell
	}
	return out
}

func v3SectionActionsNode(actions []any) map[string]any {
	buttons := make([]any, 0, len(actions))
	for _, raw := range actions {
		item, ok := toStringAnyMap(raw)
		if !ok {
			continue
		}
		props := map[string]any{}
		copyIfPresent(props, item, "variant")
		copyIfPresent(props, item, "size")
		copyIfPresent(props, item, "disabled")
		copyIfPresent(props, item, "selected")
		copyIfPresent(props, item, "action")
		buttons = append(buttons, componentNode("Button", props, v3RenderableToNode(item["label"])))
	}
	return componentNode("Inline", map[string]any{"gap": "sm", "justify": "end"}, buttons...)
}

func v3RenderableToNode(value any) map[string]any {
	if node, ok := widgetNodeFromAny(value); ok {
		return node
	}
	return map[string]any{"kind": "text", "text": fmt.Sprint(value)}
}

func (r *runtime) v3SectionToNode(spec v3SectionSpec) map[string]any {
	props := map[string]any{"label": spec.Title, "level": 1, "rule": true, "density": "flush"}
	if spec.Caption != "" {
		props["caption"] = spec.Caption
	}
	if spec.AnchorID != "" {
		props["anchorId"] = spec.AnchorID
	}
	if spec.Tone != "" {
		props["tone"] = spec.Tone
	}
	if len(spec.Actions) > 0 {
		props["actions"] = v3SectionActionsNode(spec.Actions)
	}
	return componentNode("SectionBlock", props, r.v3NodeSpecsToIR(spec.Children)...)
}

func (r *runtime) v3ExportChildren(values []goja.Value) []v3NodeSpec {
	out := []v3NodeSpec{}
	for _, value := range values {
		out = append(out, r.v3ExportChild(value)...)
	}
	return out
}

func (r *runtime) v3ExportChild(value goja.Value) []v3NodeSpec {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return nil
	}
	if exported, ok := value.Export().(bool); ok && !exported {
		return nil
	}
	if isArrayLike(value) {
		obj := value.ToObject(r.vm)
		length := int(obj.Get("length").ToInteger())
		out := []v3NodeSpec{}
		for i := 0; i < length; i++ {
			out = append(out, r.v3ExportChild(obj.Get(fmt.Sprintf("%d", i)))...)
		}
		return out
	}
	if isWidgetNode(r.vm, value) {
		node, ok := widgetNodeFromAny(value.Export())
		if !ok {
			return []v3NodeSpec{r.v3TextNode(value)}
		}
		return []v3NodeSpec{v3NodeSpecFromIR(node)}
	}
	return []v3NodeSpec{r.v3TextNode(value)}
}

func (r *runtime) v3TextNode(value goja.Value) v3NodeSpec {
	return v3NodeSpecFromIR(map[string]any{"kind": "text", "text": stringifyValue(value)})
}

func (r *runtime) v3BuildElement(tag string, attrs map[string]any, childValues []goja.Value) map[string]any {
	out := map[string]any{"kind": "element", "tag": tag}
	if len(attrs) > 0 {
		out["attrs"] = attrs
	}
	children := r.v3NodeSpecsToIR(r.v3ExportChildren(childValues))
	if len(children) > 0 {
		out["children"] = children
	}
	return out
}

func (r *runtime) v3BuildComponent(componentType string, props map[string]any, childValues []goja.Value) map[string]any {
	out := map[string]any{"kind": "component", "type": componentType}
	if len(props) > 0 {
		out["props"] = props
	}
	children := r.v3NodeSpecsToIR(r.v3ExportChildren(childValues))
	if len(children) > 0 {
		out["children"] = children
	}
	return out
}

func (r *runtime) v3NodeSpecsToIR(nodes []v3NodeSpec) []any {
	out := make([]any, 0, len(nodes))
	for _, node := range nodes {
		out = append(out, node.toIR())
	}
	return out
}

func v3NodeSpecFromIR(ir map[string]any) v3NodeSpec {
	return v3NodeSpec{Kind: stringFromMap(ir, "kind", ""), IR: ir}
}

func (n v3NodeSpec) toIR() map[string]any {
	out := map[string]any{}
	for k, v := range n.IR {
		out[k] = v
	}
	if n.Source != nil {
		out["source"] = map[string]any{"file": n.Source.File, "line": n.Source.Line, "column": n.Source.Column}
	}
	return out
}

func isV3EmptySlotResult(value goja.Value) bool {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return true
	}
	if exported, ok := value.Export().(bool); ok && !exported {
		return true
	}
	return false
}

func v3PageValidationIssues(spec *v3PageSpec) []map[string]any {
	issues := []map[string]any{}
	if strings.TrimSpace(spec.ID) == "" {
		issues = append(issues, v3ValidationIssue("page_id_required", "page.id", "page id is required"))
	}
	if strings.TrimSpace(spec.Title) == "" {
		issues = append(issues, v3ValidationIssue("page_title_required", "page.title", "page title is required"))
	}
	for sectionIndex, section := range spec.Sections {
		sectionPath := fmt.Sprintf("page.sections[%d]", sectionIndex)
		if section.Title == nil {
			issues = append(issues, v3ValidationIssue("section_title_required", sectionPath+".title", "section title is required"))
		}
		for childIndex, child := range section.Children {
			issues = append(issues, v3NodeValidationIssues(child, fmt.Sprintf("%s.children[%d]", sectionPath, childIndex))...)
		}
	}
	return issues
}

func v3NodeValidationIssues(node v3NodeSpec, path string) []map[string]any {
	issues := []map[string]any{}
	switch node.Kind {
	case "text":
		if _, ok := node.IR["text"]; !ok {
			issues = append(issues, v3ValidationIssue("text_value_required", path+".text", "text node requires a text value"))
		}
	case "element":
		if strings.TrimSpace(stringFromMap(node.IR, "tag", "")) == "" {
			issues = append(issues, v3ValidationIssue("element_tag_required", path+".tag", "element node requires a tag"))
		}
	case "component":
		if strings.TrimSpace(stringFromMap(node.IR, "type", "")) == "" {
			issues = append(issues, v3ValidationIssue("component_type_required", path+".type", "component node requires a type"))
		}
	default:
		issues = append(issues, v3ValidationIssue("node_kind_invalid", path+".kind", "node kind must be text, element, or component"))
	}
	for childIndex, child := range anySlice(node.IR["children"]) {
		childPath := fmt.Sprintf("%s.children[%d]", path, childIndex)
		childNode, ok := widgetNodeFromAny(child)
		if !ok {
			issues = append(issues, v3ValidationIssue("node_child_invalid", childPath, "node child must be a widget node"))
			continue
		}
		issues = append(issues, v3NodeValidationIssues(v3NodeSpecFromIR(childNode), childPath)...)
	}
	return issues
}

func (r *runtime) v3MetadataNode(record map[string]any) map[string]any {
	return componentNode("MetadataGrid", map[string]any{"items": v3MetadataItems(record)})
}

func v3MetadataItems(record map[string]any) []any {
	items := make([]any, 0, len(record))
	for key, value := range record {
		items = append(items, map[string]any{"key": key, "label": key, "value": value})
	}
	return items
}

func v3ValidationIssue(code string, path string, message string) map[string]any {
	return map[string]any{"severity": "error", "code": code, "path": path, "message": message}
}

func v3SelectionToV2(value any) *widgetspec.SelectionSpec {
	m, ok := value.(map[string]any)
	if !ok || m == nil {
		return nil
	}
	kind, _ := m["kind"].(string)
	if kind == "urlParam" {
		return &widgetspec.SelectionSpec{Kind: widgetspec.SelectionKindURLParam, Param: stringFromMap(m, "param", "id"), Value: stringFromMap(m, "value", "")}
	}
	if selected, ok := m["selected"].(string); ok && selected != "" {
		return &widgetspec.SelectionSpec{Kind: widgetspec.SelectionKindURLParam, Param: stringFromMap(m, "keyField", "id"), Value: selected}
	}
	return nil
}

func v3ActionFromAny(value any) widgetspec.ActionSpec {
	m, _ := value.(map[string]any)
	kind, _ := m["kind"].(string)
	action := widgetspec.ActionSpec{Kind: widgetspec.ActionKindEvent, Event: kind}
	switch kind {
	case "server":
		action.Kind = widgetspec.ActionKindServer
		action.Name = stringFromMap(m, "name", "")
	case "navigate":
		action.Kind = widgetspec.ActionKindNavigate
		action.To = stringFromMap(m, "to", "")
	case "download":
		action.Kind = widgetspec.ActionKindDownload
		action.To = stringFromMap(m, "to", "")
	case "copy":
		action.Kind = widgetspec.ActionKindCopy
		if v, ok := m["value"].(string); ok {
			action.Payload.Fields = append(action.Payload.Fields, widgetspec.PayloadFieldSpec{Name: "value", Value: widgetspec.TemplateValue{Kind: widgetspec.TemplateValueLiteral, Value: v}})
		}
	default:
		action.Kind = widgetspec.ActionKindEvent
		action.Event = stringFromMap(m, "event", kind)
	}
	if confirm, ok := m["confirm"].(string); ok && confirm != "" {
		action.Confirm = &widgetspec.TemplateSpec{Parts: []widgetspec.TemplateValue{{Kind: widgetspec.TemplateValueText, Text: confirm}}}
	}
	if payload, ok := m["payload"].(map[string]any); ok {
		for name, raw := range payload {
			action.Payload.Fields = append(action.Payload.Fields, widgetspec.PayloadFieldSpec{Name: name, Value: v3TemplateValueFromAny(raw)})
		}
	}
	return action
}

func v3TemplateValueFromAny(value any) widgetspec.TemplateValue {
	if m, ok := value.(map[string]any); ok {
		if kind, _ := m["kind"].(string); kind == "accessor" {
			return widgetspec.TemplateValue{Kind: widgetspec.TemplateValuePath, Path: stringFromMap(m, "path", stringFromMap(m, "field", stringFromMap(m, "mapField", "")))}
		}
		if kind, _ := m["kind"].(string); kind == "const" {
			return widgetspec.TemplateValue{Kind: widgetspec.TemplateValueLiteral, Value: m["value"]}
		}
	}
	return widgetspec.TemplateValue{Kind: widgetspec.TemplateValueLiteral, Value: value}
}

func v3AccessorSpec(mode string, valueKey string, value string) map[string]any {
	out := map[string]any{"kind": "accessor", "mode": mode}
	if strings.TrimSpace(value) != "" {
		out[valueKey] = value
	}
	return out
}

func v3URLTemplateValue(value goja.Value) string {
	if value == nil || goja.IsUndefined(value) || goja.IsNull(value) {
		return ""
	}
	if m, ok := value.Export().(map[string]any); ok {
		if kind, _ := m["kind"].(string); kind == "accessor" {
			path := stringFromMap(m, "path", stringFromMap(m, "field", stringFromMap(m, "mapField", stringFromMap(m, "template", ""))))
			if strings.TrimSpace(path) != "" {
				return "${" + path + "}"
			}
		}
		if kind, _ := m["kind"].(string); kind == "const" {
			return fmt.Sprint(m["value"])
		}
	}
	return value.String()
}

func slugID(s string) string {
	lower := strings.ToLower(strings.TrimSpace(s))
	var b strings.Builder
	lastDash := false
	for _, r := range lower {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	if out == "" {
		return "page"
	}
	return out
}
