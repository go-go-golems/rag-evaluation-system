# Tasks

## TODO

- [ ] Fix MeetingPollPanel readOnly so it disables editable cells and action emission <!-- t:e3jq -->
- [ ] Decide TimeGrid allDay contract: render all-day row or remove allDay from types until supported <!-- t:efal -->
- [ ] Decide and wire public exports for scheduling DTOs and widget presets <!-- t:cl35 -->
- [ ] Add focused tests for TimeGrid lane packing, MonthGrid date bounds, StyleBySpec fallback, and adapter action contexts <!-- t:pek8 -->
- [ ] Implement Goja DSL engine helpers: data.matrixGrid, ui.segmentedBar, time.monthGrid, time.timeGrid, data.cell.cycle/value, and ui.styleBy <!-- t:mp00 -->
- [ ] Add schedule.dsl and calendar.dsl recipe modules for availabilityMatrix, pollResults, monthCalendar, and weekCalendar <!-- t:3ge8 -->
- [ ] Update widgetdsl TypeScript declarations and fixture tests for the new scheduling/time DSL modules <!-- t:1gte -->
- [ ] Add widgetdsl runtime tests for new module exports, helper IR shape, and scheduling/calendar recipe output <!-- t:4e8u -->
- [ ] Inventory and classify existing ui/data/context_window/course/cms DSL helpers into foundation helpers, domain views, engine helpers, recipes, and compatibility aliases <!-- t:nj3i -->
- [ ] Prototype shared widgetdsl composition kernel: builder callbacks, .use fragments, slot invocation, child normalization, and view descriptors <!-- t:beer -->
- [ ] Add builder-lambda domain views for cms.mediaLibrary/articleQueue, course.studio/slideDeck/handout, and context.workspace/diagram over existing panel components <!-- t:ld3w -->
- [x] V3 Phase 0: finish baseline export inventory, diary setup, and tracker commit <!-- t:p4o2 -->
- [x] V3 Phase 1: add parallel widget.dsl root module skeleton and raw escape hatch <!-- t:emfh -->
- [x] V3 Phase 2: implement shared page/spec builder kernel with fragments, slots, bind, and act <!-- t:19l1 -->
- [x] V3 Phase 3: implement widget.dsl ui/page composition namespace <!-- t:2zep -->
- [x] V3 Phase 4: implement widget.dsl data namespace over existing v2 specs and MatrixGrid engine <!-- t:fvsy -->
- [x] V3 Phase 5-8: implement cms/course/context/schedule/time domain namespaces with typed views and intents <!-- t:b1un -->
- [x] V3 Phase 9-11: generate declarations/docs, add go-go-course golden fixtures, and document integration/cutover <!-- t:90l0 -->
