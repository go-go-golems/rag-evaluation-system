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

Generated from widgetV3Module.

## `page`

`page(titleOrOptions: string | Record<string, any>, configure?: Fragment<PageBuilder>): PageBuilder`. Create a page builder.

## `raw` — RawNamespace

Low-level escape hatches for text, element, component, and fragment nodes.

Runtime factory: `v3RawObject`.

- `text` (function)
- `element` (function)
- `component` (function)
- `fragment` (function)

## `act` — ActionNamespace

Generic action builders.

Runtime factory: `actionObject`.

- `server` (function)
- `navigate` (function)
- `download` (function)
- `event` (function)
- `copy` (function)
- `openOverlay` (function)
- `closeOverlay` (function)

## `bind` — BindingNamespace

Accessor and constant binding builders.

Runtime factory: `bindingObject`.

- `field` (function)
- `path` (function)
- `map` (function)
- `template` (function)
- `context` (function)
- `const` (function)

## `app` — AppNamespace

Typed application shell and viewport ownership helpers.

Runtime factory: `v3AppObject`.

- `shell` (function)
- `none` (function)
- `rootOwned` (function)

## `ui` — UINamespace

Generic composition widgets.

Runtime factory: `v3UIObject`.

- `callout` (function)
- `stack` (function)
- `inline` (function)
- `splitPane` (function)
- `card` (function)
- `button` (function)
- `caption` (function)
- `badge` (function)
- `metadata` (function)
- `shareLink` (function)
- `form` (function)
- `formRow` (function)
- `textInput` (function)
- `textareaInput` (function)
- `selectInput` (function)
- `status` (function)
- `emptyState` (function)
- `text` (function)
- `code` (function)
- `divider` (function)
- `disclosure` (function)
- `scroll` (function)
- `tabs` (function)
- `summary` (function)
- `checkList` (function)
- `stepList` (function)
- `markdownArticle` (function)
- `upload` (function)
- `formDialog` (function)

## `data` — DataNamespace

Schema, collection, matrix, selection, cell, and item helpers.

Runtime factory: `v3DataObject`.

- `fields` (function)
- `collection` (function)
- `selection` (function)
- `item` (function)
- `matrix` (function)
- `activityFeed` (function)
- `cell` (object)

## `crm` — CrmNamespace

CRM field schemas, pipelines, records, activities, tasks, and actions.

Runtime factory: `v3CRMObject`.

- `fields` (function)
- `pipeline` (function)
- `pipelineBoard`: `pipelineBoard(pipeline: Record<string, any>, deals: Record<string, any>[], configure?: Fragment<CrmPipelineBoardBuilder>): WidgetNodeSpec` → `BoardEngine`. Render an opportunity pipeline board.
- `recordFields`: `recordFields(values: Record<string, JsonValue>, fields: Record<string, any>, configure?: Fragment<CrmRecordFieldsBuilder>): WidgetNodeSpec` → `RecordFieldList`. Render typed CRM fields.
- `field` (function)
- `tasksInbox` (function)
- `stat` (function)
- `funnel` (function)
- `intent` (object)

## `cms` — CmsNamespace

CMS media, queue, and markdown editor helpers.

Runtime factory: `v3CMSObject`.

- `shell` (function)
- `mediaLibrary` (function)
- `articleQueue` (function)
- `markdownEditor` (function)
- `intent` (object)

## `course` — CourseNamespace

Course shell, landing, slide, handout, metadata, agenda, and material helpers.

Runtime factory: `v3CourseObject`.

- `shell` (function)
- `landing` (function)
- `slideDeck` (function)
- `handouts` (function)
- `metadataForm` (function)
- `agendaEditor` (function)
- `materialUploads` (function)
- `intent` (object)

## `context` — ContextNamespace

Context style, diagram, transcript workspace, and intent helpers.

Runtime factory: `v3ContextObject`.

- `styleSet` (function)
- `palette` (function)
- `diagram` (function)
- `workspace` (function)
- `intent` (object)

## `schedule` — ScheduleNamespace

Availability poll, poll summary, booking picker, and schedule intent helpers.

Runtime factory: `v3ScheduleObject`.

- `availabilityPoll`: `availabilityPoll(poll: AvailabilityPoll, configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec` → `MatrixGrid`. Render respondent availability against poll options.
- `pollSummary`: `pollSummary(poll: AvailabilityPoll, tallies: PollTally[], configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec` → `MatrixGrid`. Render aggregate option tallies.
- `bookingPicker`: `bookingPicker(availability: Record<string, any>, configure?: Fragment<SchedulePollBuilder>): WidgetNodeSpec` → `MatrixGrid`. Render bookable resources by slot.
- `intent` (object)

## `time` — TimeNamespace

Month, week, formatting, range, and time intent helpers.

Runtime factory: `v3TimeObject`.

- `month`: `month(eventsOrMarkers: CalendarEvent[] | Record<string, any>, configure?: Fragment<TimeMonthBuilder>): WidgetNodeSpec` → `MonthGrid`. Render day markers for a month.
- `week`: `week(events: CalendarEvent[], configure?: Fragment<TimeWeekBuilder>): WidgetNodeSpec` → `TimeGrid`. Render event blocks for a week; allDay is intentionally omitted.
- `format` (function)
- `formatRange` (function)
- `slotLabel` (function)
- `range` (object)
- `intent` (object)

## `style` — Record<string, any>

Reserved style namespace.

Runtime factory: `NewObject`.

## Nested namespaces

- `data.cell`: [field status template cycle value]
- `data.selection`: [urlParam]
- `crm.intent`: [openDeal moveDeal updateField completeTask]
- `cms.intent`: [selectAsset openAsset uploadAssets selectArticle createArticle publishArticle archiveArticle previewArticle]
- `course.intent`: [navigate selectHandout downloadHandout printHandout previousSlide nextSlide presentSlide editAgenda uploadMaterial deleteMaterial]
- `context.intent`: [selectPart selectAnnotation]
- `schedule.intent`: [toggleAvailability submitResponse]
- `time.range`: [week]
- `time.intent`: [selectDay selectEvent]

## Composable builders

All builders below expose `use(fragment)` in addition to their listed methods.

- `PageBuilder`: [id title meta shell shortcuts root density breadcrumb section view validate toPage use]
- `ShortcutsBuilder`: [bind use]
- `AppShellBuilder`: [brand navigation content use]
- `NavigationBuilder`: [placement active width narrowMode ariaLabel section use]
- `NavigationItemsBuilder`: [item use]
- `ContentViewportBuilder`: [maxWidth padding scroll use]
- `SectionBuilder`: [caption anchor tone text view slot actions metric metadata use]
- `ActionsBuilder`: [add button use]
- `FormDialogBuilder`: [title body initialFocus submitLabel cancelLabel submit use]
- `FieldSetBuilder`: [key primary short prose count status date currency media url build validate use]
- `CollectionBuilder`: [id schema empty select search paginate table edit masterDetail validate toNode toIR use]
- `SearchBuilder`: [value query placeholder resultCount submit clear use]
- `PaginationBuilder`: [current size total sizes position onChange use]
- `TableBuilder`: [className rowSelect actionColumn sortable keyboard command styleWhen use]
- `TableKeyboardBuilder`: [mode selection vimAliases enterSelect use]
- `RowCommandBuilder`: [key label danger action use]
- `EditorBuilder`: [create submit submitPost reorder remove actions use]
- `MatrixBuilder`: [id columns column valueAt cell onCellAction toNode use]
- `SchedulePollBuilder`: [styleSet readOnly editableRow selectedCell onToggle ariaLabel use]
- `TimeMonthBuilder`: [styleSet selected today weekStartsOn onSelect use]
- `TimeWeekBuilder`: [styleSet range hours hourHeight viewportHeight now selected onSelect onSlotCreate use]
- `ContextStyleSetBuilder`: [style legend use]
- `ContextDiagramBuilder`: [styleSet palette view selected legend empty onSelect use]
- `ContextWorkspaceBuilder`: [selectedAnnotation showNotes styleSet message annotation empty onAnnotationSelect use]
- `CourseShellBuilder`: [active subtitle contentPadding main footer onNavigate use]
- `CourseLandingBuilder`: [activeAgenda onAgendaSelect onPrimary onSecondary use]
- `CourseSlideDeckBuilder`: [mode visualSide onPrevious onNext onPresent onFullscreen use]
- `CourseHandoutsBuilder`: [selected title empty onSelect onDownload onPrint use]
- `CourseMetadataFormBuilder`: [title onSubmit use]
- `CourseMaterialUploadsBuilder`: [accept onUpload onDelete use]
- `CmsShellBuilder`: [active subtitle contentPadding main header footer onNavigate use]
- `CmsMediaLibraryBuilder`: [selection selected query kindFilter page empty accept asset details toolbar onSelect onOpen onUpload use]
- `CmsArticleQueueBuilder`: [selected status query page empty row rowActions filters onSelect onCreate onRowAction onPublish onArchive onPreview use]
- `CmsMarkdownEditorBuilder`: [title placeholder onChange onSubmit use]
- `CrmFieldsBuilder`: [text longtext email phone url number currency percent date datetime boolean select multiselect tags user relation build validate use]
- `CrmPipelineBuilder`: [stage build validate use]
- `CrmPipelineBoardBuilder`: [summaries selected ariaLabel onMove onOpen use]
- `CrmRecordFieldsBuilder`: [mode refs onChange use]
- `ActivityFeedBuilder`: [groupByDay glyph glyphs styleSet onOpen onLoadMore use]

## Action contexts

- `page.shortcut` (`PageShortcut`): `[pageId shortcutId key componentType]`. Context dispatched when a page-level keyboard shortcut activates.
- `app.navigate` (`SidebarNav`): `[itemId value componentType]`. Context dispatched by typed application navigation.
- `table.rowSelect` (`DataTable`): `[row rowKey componentType]`. Context dispatched when a collection row is selected.
- `table.cellAction` (`DataTableCell`): `[row rowKey componentType]`. Context dispatched by an action-button cell.
- `matrix.cellAction` (`MatrixGrid`): `[row column rowKey colId value componentType]`. Context dispatched when a matrix cell is activated.
- `context.annotationSelect` (`TranscriptWorkspacePanel`): `[annotationId value componentType]`. Context dispatched when a transcript annotation is selected.
- `course.navigate` (`CourseStudioShell`): `[itemId item value componentType]`. Context dispatched from course navigation.
- `course.agendaSelect` (`CourseLessonPanel`): `[agendaItemId value componentType]`. Context dispatched when an agenda item is selected.
- `course.cta` (`CourseLessonPanel`): `[cta componentType]`. Context dispatched by a primary or secondary course call to action.
- `course.slideControl` (`CourseSlidePanel`): `[value componentType]`. Context dispatched by previous, next, present, and fullscreen slide controls.
- `course.handout` (`HandoutDocumentShell`): `[documentId document value componentType]`. Context dispatched when a handout is selected, downloaded, or printed.
- `upload.files` (`ContextUploadDropArea`): `[files fileNames fileCount componentType]`. Serialized file context dispatched by a generic upload area.
- `cms.asset` (`MediaLibraryPanel`): `[assetId value componentType]`. Context dispatched when an asset is selected or opened.
- `cms.assetQuery` (`MediaLibraryPanel`): `[query value componentType]`. Context dispatched when an asset query is submitted.
- `cms.assetKind` (`MediaLibraryPanel`): `[kind value componentType]`. Context dispatched when the asset kind filter changes.
- `cms.assetPage` (`MediaLibraryPanel`): `[page value componentType]`. Context dispatched when the asset page changes.
- `cms.assetUpload` (`MediaLibraryPanel`): `[files fileNames fileCount componentType]`. Serialized file context dispatched by the media library.
- `cms.article` (`ArticleListPanel`): `[articleId value componentType]`. Context dispatched when an article is selected.
- `cms.articleRowAction` (`ArticleListPanel`): `[articleId rowAction value componentType]`. Context dispatched by an article row action.
- `cms.articleStatus` (`ArticleListPanel`): `[status value componentType]`. Context dispatched when the article status filter changes.
- `cms.articleQuery` (`ArticleListPanel`): `[query value componentType]`. Context dispatched when an article query is submitted.
- `cms.articlePage` (`ArticleListPanel`): `[page value componentType]`. Context dispatched when the article page changes.
- `time.daySelect` (`MonthGrid`): `[dateISO value componentType]`. Context dispatched when a calendar day is selected.
- `time.blockSelect` (`TimeGrid`): `[blockId value componentType]`. Context dispatched when a time block is selected.
- `time.slotCreate` (`TimeGrid`): `[dayISO hour value componentType]`. Context dispatched when an empty time slot is activated.
- `crm.boardMove` (`BoardEngine`): `[cardId from to beforeId componentType]`. Context dispatched when a CRM card moves.
- `crm.boardOpen` (`BoardEngine`): `[cardId componentType]`. Context dispatched when a CRM card is opened.
- `crm.fieldChange` (`RecordFieldList`): `[key value componentType]`. Context dispatched when a CRM record field changes.
- `activity.open` (`ActivityFeed`): `[activityId componentType]`. Context dispatched when an activity is opened.
- `activity.loadMore` (`ActivityFeed`): `[componentType]`. Context dispatched when earlier activities are requested.

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
