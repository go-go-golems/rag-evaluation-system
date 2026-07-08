package widgetdsl

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/engine"
)

func TestSplitModulesExportExpectedHelpersAndOmitCrossDomainHelpers(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const ui = require("ui.dsl");
		const data = require("data.dsl");
		const contextWindow = require("context_window.dsl");
		const course = require("course.dsl");
		({
			uiPage: typeof ui.page,
			uiPanel: typeof ui.panel,
			uiFormPanel: typeof ui.formPanel,
			uiTextareaInput: typeof ui.textareaInput,
			uiDataTable: typeof ui.dataTable,
			dataTable: typeof data.dataTable,
			dataCellField: typeof data.cell.field,
			dataPage: typeof data.page,
			contextDiagramPanel: typeof contextWindow.contextDiagramPanel,
			contextStyleSwatch: typeof contextWindow.contextStyleSwatch,
			contextVisualStyle: typeof contextWindow.visualStyle,
			contextStyleSet: typeof contextWindow.styleSet,
			contextPart: typeof contextWindow.contextPart,
			contextStudioNavIconFromContext: typeof contextWindow.contextStudioNavIcon,
			courseStudioNavIcon: typeof course.contextStudioNavIcon,
			courseStudioShell: typeof course.courseStudioShell,
		});
	`)
	if err != nil {
		t.Fatalf("require split modules: %v", err)
	}
	got := value.Export().(map[string]any)
	wantFunctions := []string{"uiPage", "uiPanel", "uiFormPanel", "uiTextareaInput", "dataTable", "dataCellField", "contextDiagramPanel", "contextStyleSwatch", "contextVisualStyle", "contextStyleSet", "contextPart", "courseStudioNavIcon", "courseStudioShell"}
	for _, name := range wantFunctions {
		if got[name] != "function" {
			t.Fatalf("%s export = %#v, want function (all: %#v)", name, got[name], got)
		}
	}
	wantUndefined := []string{"uiDataTable", "dataPage", "contextStudioNavIconFromContext"}
	for _, name := range wantUndefined {
		if got[name] != "undefined" {
			t.Fatalf("%s export = %#v, want undefined (all: %#v)", name, got[name], got)
		}
	}
}

func TestWidgetV3ModuleExportsRootNamespacesAndKeepsOldModulesAvailable(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const ui = require("ui.dsl");
		const node = widget.raw.component("Panel", { title: "V3" }, widget.raw.text("Hello"));
		({
			legacyPanel: typeof ui.panel,
			rawComponent: typeof widget.raw.component,
			rawText: typeof widget.raw.text,
			actServer: typeof widget.act.server,
			bindField: typeof widget.bind.field,
			pageFn: typeof widget.page,
			uiNamespace: typeof widget.ui,
			dataNamespace: typeof widget.data,
			cmsNamespace: typeof widget.cms,
			courseNamespace: typeof widget.course,
			contextNamespace: typeof widget.context,
			scheduleNamespace: typeof widget.schedule,
			timeNamespace: typeof widget.time,
			styleNamespace: typeof widget.style,
			node,
			binding: widget.bind.field("title"),
			action: widget.act.server("save", { payload: { id: 1 } }),
			dataSelection: typeof widget.data.selection,
			dataItem: typeof widget.data.item,
		});
	`)
	if err != nil {
		t.Fatalf("require widget.dsl: %v", err)
	}
	got := value.Export().(map[string]any)
	wantFunctions := []string{"legacyPanel", "rawComponent", "rawText", "actServer", "bindField", "pageFn", "dataSelection", "dataItem"}
	for _, name := range wantFunctions {
		if got[name] != "function" {
			t.Fatalf("%s export = %#v, want function (all: %#v)", name, got[name], got)
		}
	}
	wantObjects := []string{"uiNamespace", "dataNamespace", "cmsNamespace", "courseNamespace", "contextNamespace", "scheduleNamespace", "timeNamespace", "styleNamespace"}
	for _, name := range wantObjects {
		if got[name] != "object" {
			t.Fatalf("%s export = %#v, want object (all: %#v)", name, got[name], got)
		}
	}
	node := got["node"].(map[string]any)
	if node["kind"] != "component" || node["type"] != "Panel" {
		t.Fatalf("raw.component emitted %#v", node)
	}
	binding := got["binding"].(map[string]any)
	if binding["kind"] != "accessor" || binding["mode"] != "field" || binding["field"] != "title" {
		t.Fatalf("bind.field emitted %#v", binding)
	}
	action := got["action"].(map[string]any)
	if action["kind"] != "server" || action["name"] != "save" {
		t.Fatalf("act.server emitted %#v", action)
	}
}

func TestWidgetV3PageBuilderEmitsPageIR(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const quiet = s => s.tone("quiet");
		const builder = widget.page("Hello V3", p => p
			.id("hello-v3")
			.meta("source", "test")
			.use(p => p.title("Hello Widget V3"))
			.section("Intro", s => s
				.use(quiet)
				.caption("Builder callbacks lower to serializable Widget IR.")
				.anchor("intro")
				.text("Hello")
				.view(widget.raw.component("Caption", { tone: "muted" }, "World"))));
		({ issues: builder.validate(), page: builder.toPage() });
	`)
	if err != nil {
		t.Fatalf("build widget.dsl page: %v", err)
	}
	got := value.Export().(map[string]any)
	if issueCount := exportedSliceLen(got["issues"]); issueCount != 0 {
		t.Fatalf("unexpected validation issues: %#v", got["issues"])
	}
	page := got["page"].(map[string]any)
	if page["id"] != "hello-v3" || page["title"] != "Hello Widget V3" {
		t.Fatalf("unexpected page identity: %#v", page)
	}
	meta := page["meta"].(map[string]any)
	if meta["source"] != "test" {
		t.Fatalf("unexpected page meta: %#v", meta)
	}
	root := page["root"].(map[string]any)
	if root["type"] != "Stack" {
		t.Fatalf("unexpected root: %#v", root)
	}
	children := root["children"].([]any)
	if len(children) != 1 {
		t.Fatalf("root children = %#v, want one section", children)
	}
	section := children[0].(map[string]any)
	if section["type"] != "SectionBlock" {
		t.Fatalf("unexpected section: %#v", section)
	}
	props := section["props"].(map[string]any)
	if props["label"] != "Intro" || props["caption"] != "Builder callbacks lower to serializable Widget IR." || props["anchorId"] != "intro" || props["tone"] != "quiet" {
		t.Fatalf("unexpected section props: %#v", props)
	}
	sectionChildren := section["children"].([]any)
	if len(sectionChildren) != 2 {
		t.Fatalf("section children = %#v, want text + caption", sectionChildren)
	}
}

func TestWidgetV3UICompositionHelpersEmitPageIR(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const branded = p => p.density("compact").breadcrumb("Home", "/");
		const page = widget.page("UI Page", p => p
			.use(branded)
			.shell({ kind: "app" })
			.section("Overview", s => s
				.actions(a => a.button("Refresh", widget.act.event("refresh")))
				.caption("Built only with widget.dsl")
				.metric("Open", 12, { tone: "good" })
				.metadata({ Owner: "Ops", Status: "Ready" })
				.view(widget.ui.stack({ gap: "sm" },
					widget.ui.callout({ title: "Note" }, "Hello"),
					widget.ui.inline(widget.ui.badge("new"), widget.ui.caption("caption")),
					widget.ui.card({ title: "Card" }, widget.ui.button("Save", widget.act.server("save"))),
					widget.ui.form({ title: "Form" }, widget.raw.text("Body"))
				))
			)
		).toPage();
		page;
	`)
	if err != nil {
		t.Fatalf("build widget.dsl UI page: %v", err)
	}
	page := value.Export().(map[string]any)
	if page["shell"].(map[string]any)["kind"] != "app" {
		t.Fatalf("page shell = %#v", page["shell"])
	}
	root := page["root"].(map[string]any)
	props := root["props"].(map[string]any)
	if props["density"] != "compact" {
		t.Fatalf("root props = %#v, want compact density", props)
	}
	rootChildren := anySlice(root["children"])
	if rootChildren[0].(map[string]any)["type"] != "Breadcrumbs" {
		t.Fatalf("first root child = %#v, want Breadcrumbs", rootChildren[0])
	}
	section := rootChildren[1].(map[string]any)
	sectionProps := section["props"].(map[string]any)
	if len(anySlice(sectionProps["actions"])) != 1 {
		t.Fatalf("section actions = %#v, want one action", sectionProps["actions"])
	}
	children := anySlice(section["children"])
	if len(children) != 3 {
		t.Fatalf("section children = %#v, want metric + metadata + stack", children)
	}
	if children[0].(map[string]any)["type"] != "KeyValueStrip" || children[1].(map[string]any)["type"] != "MetadataGrid" || children[2].(map[string]any)["type"] != "Stack" {
		t.Fatalf("unexpected UI composition children: %#v", children)
	}
}

func TestWidgetV3SlotsAndChildNormalization(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const fallback = (ctx, h) => h.caption("fallback: " + ctx.label);
		const page = widget.page("Slots", p => p.section("Slot section", s => s
			.view(["A", null, undefined, false, widget.raw.fragment("B", false)])
			.slot({ label: "primary" }, (ctx, h) => h.stack(
				{ gap: "sm" },
				h.strong(ctx.label),
				false,
				h.badge("ready", { tone: "success" })
			), fallback)
			.slot({ label: "fallback" }, null, fallback)
		)).toPage();
		page;
	`)
	if err != nil {
		t.Fatalf("build widget.dsl page with slots: %v", err)
	}
	page := value.Export().(map[string]any)
	root := page["root"].(map[string]any)
	section := anySlice(root["children"])[0].(map[string]any)
	children := anySlice(section["children"])
	if len(children) != 4 {
		t.Fatalf("section children = %#v, want A, B, slot stack, fallback caption", children)
	}
	first := children[0].(map[string]any)
	if first["kind"] != "text" || first["text"] != "A" {
		t.Fatalf("first normalized child = %#v, want text A", first)
	}
	second := children[1].(map[string]any)
	if second["kind"] != "text" || second["text"] != "B" {
		t.Fatalf("second normalized child = %#v, want text B", second)
	}
	stack := children[2].(map[string]any)
	if stack["type"] != "Stack" {
		t.Fatalf("slot stack child = %#v, want Stack", stack)
	}
	stackChildren := anySlice(stack["children"])
	if len(stackChildren) != 2 {
		t.Fatalf("slot stack children = %#v, want strong + badge", stackChildren)
	}
	fallbackCaption := children[3].(map[string]any)
	if fallbackCaption["type"] != "Caption" {
		t.Fatalf("fallback child = %#v, want Caption", fallbackCaption)
	}
}

func TestWidgetV3ScheduleAndTimeViews(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const poll = {
			title: "Team sync availability",
			options: [
				{ id: "mon-9", label: "Mon 09:00", startISO: "2026-07-06T09:00:00Z", endISO: "2026-07-06T09:30:00Z" },
				{ id: "mon-10", label: "Mon 10:00", startISO: "2026-07-06T10:00:00Z", endISO: "2026-07-06T10:30:00Z" },
			],
			responses: [
				{ id: "ana", name: "Ana", availability: { "mon-9": "available", "mon-10": "maybe" } },
			],
		};
		const editablePoll = widget.schedule.availabilityPoll(poll, p => p
			.editableRow("ana")
			.selectedCell("ana", "mon-9")
			.onToggle(widget.schedule.intent.toggleAvailability(widget.bind.context("row.id"), widget.bind.context("column.id"), widget.bind.context("nextValue")))
		);
		const readOnlyPoll = widget.schedule.availabilityPoll(poll, p => p
			.onToggle(widget.schedule.intent.toggleAvailability("ana", "mon-9", "available"))
			.readOnly()
		);
		const summary = widget.schedule.pollSummary(poll, [{ id: "available", label: "Available", counts: { "mon-9": 1, "mon-10": 0 } }]);
		const booking = widget.schedule.bookingPicker({
			title: "Book room",
			resources: [{ id: "room-a", label: "Room A", availability: { "mon-9": "available" } }],
			slots: poll.options,
		}, b => b.onToggle(widget.schedule.intent.toggleAvailability("room-a", "mon-9", "selected")));
		const range = widget.time.range.week("2026-07-08");
		const month = widget.time.month([{ id: "ev1", title: "Launch", startISO: "2026-07-08T09:00:00Z", endISO: "2026-07-08T10:00:00Z", styleKey: "busy" }], m => m
			.selected("2026-07-08")
			.onSelect(widget.time.intent.selectDay(widget.bind.context("dayISO")))
		);
		const week = widget.time.week([
			{ id: "ev1", title: "Launch", startISO: "2026-07-08T09:00:00Z", endISO: "2026-07-08T10:00:00Z", styleKey: "busy", allDay: true },
		], w => w
			.range(range)
			.hours(7, 19)
			.selected("ev1")
			.onSelect(widget.time.intent.selectEvent(widget.bind.context("block.id")))
		);
		({ editablePoll, readOnlyPoll, summary, booking, range, month, week, slot: widget.time.slotLabel("2026-07-08T09:00:00Z", "2026-07-08T10:00:00Z") });
	`)
	if err != nil {
		t.Fatalf("build widget.dsl schedule/time views: %v", err)
	}
	got := value.Export().(map[string]any)
	editable := anyMap(got["editablePoll"])
	if editable["type"] != "MatrixGrid" {
		t.Fatalf("editable poll = %#v, want MatrixGrid", editable)
	}
	editableProps := anyMap(editable["props"])
	if editableProps["onCellAction"] == nil || editableProps["editableRowKey"] != "ana" || len(anySlice(editableProps["columns"])) != 2 {
		t.Fatalf("editable poll props = %#v", editableProps)
	}
	readOnlyProps := anyMap(anyMap(got["readOnlyPoll"])["props"])
	if readOnlyProps["onCellAction"] != nil || readOnlyProps["editableRowKey"] != nil {
		t.Fatalf("read-only poll leaked edit props = %#v", readOnlyProps)
	}
	if anyMap(got["summary"])["type"] != "MatrixGrid" || anyMap(got["booking"])["type"] != "MatrixGrid" {
		t.Fatalf("summary/booking = %#v / %#v", got["summary"], got["booking"])
	}
	rangeSpec := anyMap(got["range"])
	if rangeSpec["startISO"] != "2026-07-06" || rangeSpec["endISO"] != "2026-07-12" {
		t.Fatalf("range = %#v", rangeSpec)
	}
	monthProps := anyMap(anyMap(got["month"])["props"])
	if anyMap(monthProps["markers"])["2026-07-08"] == nil || monthProps["onDaySelectAction"] == nil {
		t.Fatalf("month props = %#v", monthProps)
	}
	week := anyMap(got["week"])
	if week["type"] != "TimeGrid" {
		t.Fatalf("week = %#v, want TimeGrid", week)
	}
	weekProps := anyMap(week["props"])
	blocks := anySlice(weekProps["blocks"])
	if weekProps["onBlockSelectAction"] == nil || len(blocks) != 1 || anyMap(blocks[0])["allDay"] != nil {
		t.Fatalf("week props = %#v", weekProps)
	}
	if got["slot"] != "09:00–10:00" {
		t.Fatalf("slot label = %#v", got["slot"])
	}
}

func TestWidgetV3ContextDomainViews(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const styleSet = widget.context.styleSet(s => s
			.style("prompt", { fill: "#fff", line: "#111" })
			.legend("prompt", "Prompt", { swatch: "solid" })
		);
		const palette = widget.context.palette("Dusty Magenta / Blue", [
			{ id: "answer", label: "Answer", accent: "a", pattern: "solid" },
		]);
		const snapshot = { id: "ctx", title: "Context", limit: 1000, parts: [{ id: "prompt", label: "Prompt", styleKey: "prompt", tokens: 10 }] };
		const diagram = widget.context.diagram(snapshot, d => d
			.styleSet(styleSet)
			.view("stack")
			.selected("prompt")
			.legend((ctx, h) => h.caption("Legend"))
			.empty((ctx, h) => h.caption("Empty"))
			.onSelect(widget.context.intent.selectPart(widget.bind.context("part.id")))
		);
		const workspace = widget.context.workspace({
			title: "Session",
			subtitle: "Review",
			messages: [{ id: "m1", role: "user", content: "Hello" }],
			annotations: [{ id: "a1", messageId: "m1", note: "Important" }],
			snapshot,
		}, w => w
			.selectedAnnotation("a1")
			.showNotes(false)
			.styleSet(palette)
			.message((message, h) => h.card({ title: message.role }, message.content))
			.annotation((annotation, h) => h.caption(annotation.note))
			.empty((ctx, h) => h.caption("No transcript"))
			.onAnnotationSelect(widget.context.intent.selectAnnotation(widget.bind.context("annotation.id")))
		);
		({ styleSet, palette, diagram, workspace });
	`)
	if err != nil {
		t.Fatalf("build widget.dsl context views: %v", err)
	}
	got := value.Export().(map[string]any)
	styleSet := anyMap(got["styleSet"])
	if len(anySlice(styleSet["legend"])) != 1 || anyMap(styleSet["styles"])["prompt"] == nil {
		t.Fatalf("styleSet = %#v", styleSet)
	}
	palette := anyMap(got["palette"])
	if len(anySlice(palette["legend"])) != 1 || anyMap(palette["styles"])["answer"] == nil {
		t.Fatalf("palette = %#v", palette)
	}
	diagram := anyMap(got["diagram"])
	if diagram["type"] != "ContextDiagramPanel" {
		t.Fatalf("diagram = %#v, want ContextDiagramPanel", diagram)
	}
	diagramProps := anyMap(diagram["props"])
	if diagramProps["selectedPartId"] != "prompt" || diagramProps["onPartSelectAction"] == nil || diagramProps["legendSlot"] == nil {
		t.Fatalf("diagram props = %#v", diagramProps)
	}
	workspace := anyMap(got["workspace"])
	if workspace["type"] != "TranscriptWorkspacePanel" {
		t.Fatalf("workspace = %#v, want TranscriptWorkspacePanel", workspace)
	}
	workspaceProps := anyMap(workspace["props"])
	if workspaceProps["selectedAnnotationId"] != "a1" || workspaceProps["showNotes"] != false || workspaceProps["onAnnotationSelectAction"] == nil {
		t.Fatalf("workspace props = %#v", workspaceProps)
	}
}

func TestWidgetV3CourseDomainViews(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const definition = { title: "Course", subtitle: "Week 1", sections: [{ id: "intro", label: "Intro", items: [] }] };
		const shell = widget.course.shell(definition, s => s
			.active("intro")
			.contentPadding("none")
			.footer(widget.ui.caption("Draft"))
			.main(widget.ui.card({ title: "Welcome" }, "Hello"))
			.onNavigate(widget.course.intent.navigate(widget.bind.context("item.id")))
		);
		const landing = widget.course.landing({ title: "Course", agenda: [] }, l => l
			.activeAgenda("day-1")
			.onAgendaSelect(widget.course.intent.editAgenda(widget.bind.context("agenda.id")))
		);
		const deck = widget.course.slideDeck({ slides: [{ id: "s1", title: "Slide" }], index: 0 }, d => d
			.mode("present")
			.visualSide("right")
			.onPrevious(widget.course.intent.previousSlide())
			.onNext(widget.course.intent.nextSlide())
			.onPresent(widget.course.intent.presentSlide())
		);
		const handouts = widget.course.handouts({ intro: "Read", docs: [{ id: "h1", title: "Handout" }] }, h => h
			.selected("h1")
			.title("Handouts")
			.onSelect(widget.course.intent.selectHandout(widget.bind.context("document.id")))
			.onDownload(widget.course.intent.downloadHandout(widget.bind.context("document.id")))
			.onPrint(widget.course.intent.printHandout(widget.bind.context("document.id")))
		);
		const metadata = widget.course.metadataForm({ Title: "Course" }, f => f
			.title("Metadata")
			.onSubmit(widget.act.server("save-metadata"))
		);
		const agenda = widget.course.agendaEditor([{ id: "a", title: "Intro" }], c => c
			.schema(widget.data.fields(f => f.key("id").primary("title")).build())
			.edit(e => e.submit("/agenda"))
		).toNode();
		const uploads = widget.course.materialUploads({ description: "Upload" }, u => u
			.accept(["application/pdf"])
			.onUpload(widget.course.intent.uploadMaterial())
			.onDelete(widget.course.intent.deleteMaterial(widget.bind.context("material.id")))
		);
		({ shell, landing, deck, handouts, metadata, agenda, uploads });
	`)
	if err != nil {
		t.Fatalf("build widget.dsl course views: %v", err)
	}
	got := value.Export().(map[string]any)
	for name, wantType := range map[string]string{
		"shell":    "CourseStudioShell",
		"landing":  "CourseLessonPanel",
		"deck":     "CourseSlidePanel",
		"handouts": "HandoutDocumentShell",
		"metadata": "FormPanel",
		"agenda":   "Stack",
		"uploads":  "ContextUploadDropArea",
	} {
		node := anyMap(got[name])
		if node["type"] != wantType {
			t.Fatalf("%s = %#v, want %s", name, node, wantType)
		}
	}
	shellProps := anyMap(anyMap(got["shell"])["props"])
	if shellProps["activeItemId"] != "intro" || shellProps["onNavigateAction"] == nil {
		t.Fatalf("shell props = %#v", shellProps)
	}
	handoutProps := anyMap(anyMap(got["handouts"])["props"])
	if handoutProps["selectedDocumentId"] != "h1" || handoutProps["onPrintAction"] == nil {
		t.Fatalf("handout props = %#v", handoutProps)
	}
}

func TestWidgetV3CMSDomainViews(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const assets = [{ id: "a", kind: "image", title: "Hero", filename: "hero.png", tags: [] }];
		const articles = [{ id: "p", title: "Post", status: "draft", tags: [] }];
		const media = widget.cms.mediaLibrary(assets, m => m
			.selection("multi")
			.selected(["a"])
			.query("hero")
			.kindFilter("image")
			.page(2, 4)
			.empty("No assets")
			.accept(["image/png"])
			.asset((asset, h) => h.card({ title: asset.title }))
			.details((asset, h) => h.caption(asset.filename))
			.toolbar((ctx, h) => h.button("Upload", widget.cms.intent.uploadAssets()))
			.onSelect(widget.cms.intent.selectAsset(widget.bind.context("asset.id")))
			.onOpen(widget.cms.intent.openAsset(widget.bind.context("asset.id")))
			.onUpload(widget.cms.intent.uploadAssets())
		);
		const queue = widget.cms.articleQueue(articles, q => q
			.selected("p")
			.status("draft")
			.query("post")
			.page(1, 3)
			.empty("No posts")
			.row((article, h) => h.caption(article.title))
			.rowActions((article, h) => h.button("Publish", widget.cms.intent.publishArticle(article.id)))
			.filters((ctx, h) => h.inline("Filters"))
			.onSelect(widget.cms.intent.selectArticle(widget.bind.context("article.id")))
			.onCreate(widget.cms.intent.createArticle())
			.onRowAction(widget.cms.intent.previewArticle(widget.bind.context("article.id")))
			.onPublish(widget.cms.intent.publishArticle(widget.bind.context("article.id")))
			.onArchive(widget.cms.intent.archiveArticle(widget.bind.context("article.id")))
			.onPreview(widget.cms.intent.previewArticle(widget.bind.context("article.id")))
		);
		const editor = widget.cms.markdownEditor("# Draft", e => e
			.title("Body")
			.placeholder("Write...")
			.onChange(widget.act.event("body-change"))
			.onSubmit(widget.act.server("save-body"))
		);
		({ media, queue, editor });
	`)
	if err != nil {
		t.Fatalf("build widget.dsl cms views: %v", err)
	}
	got := value.Export().(map[string]any)
	media := anyMap(got["media"])
	if media["type"] != "MediaLibraryPanel" {
		t.Fatalf("media = %#v, want MediaLibraryPanel", media)
	}
	mediaProps := anyMap(media["props"])
	if mediaProps["selectionMode"] != "multi" || mediaProps["query"] != "hero" || mediaProps["onFilesSelectedAction"] == nil {
		t.Fatalf("media props = %#v", mediaProps)
	}
	queue := anyMap(got["queue"])
	if queue["type"] != "ArticleListPanel" {
		t.Fatalf("queue = %#v, want ArticleListPanel", queue)
	}
	queueProps := anyMap(queue["props"])
	if queueProps["selectedArticleId"] != "p" || queueProps["statusFilter"] != "draft" || queueProps["onCreateAction"] == nil {
		t.Fatalf("queue props = %#v", queueProps)
	}
	editor := anyMap(got["editor"])
	if editor["type"] != "MarkdownEditor" {
		t.Fatalf("editor = %#v, want MarkdownEditor", editor)
	}
	editorProps := anyMap(editor["props"])
	if editorProps["value"] != "# Draft" || editorProps["title"] != "Body" || editorProps["onSubmitAction"] == nil {
		t.Fatalf("editor props = %#v", editorProps)
	}
}

func TestWidgetV3DataCollectionMatchesDataV2TableShape(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const v2 = require("data.v2.dsl");
		const rows = [{ id: "a", title: "Alpha", status: "draft" }];
		const oldSchema = v2.schema("articles")
			.field("id", v2.f.key().label("ID"))
			.field("title", v2.f.primary().label("Title"))
			.field("status", v2.f.status().label("Status"))
			.build();
		const oldNode = v2.collection("articles", rows)
			.schema(oldSchema)
			.empty("No articles")
			.select(v2.selection.urlParam("article", "a"))
			.table(t => t.actionColumn("open", "Open", "Open", v2.action.navigate("?article=${row.id}")))
			.toIR();
		const schema = widget.data.fields("articles", f => f
			.key("id", { label: "ID" })
			.primary("title", { label: "Title" })
			.status("status", { label: "Status" })
		).build();
		const node = widget.data.collection("articles", rows, c => c
			.schema(schema)
			.empty("No articles")
			.select(widget.data.selection.urlParam("article", "a"))
			.table(t => t.actionColumn("open", "Open", "Open", widget.act.navigate("?article=${row.id}")))
		).toNode();
		({ oldNode, node });
	`)
	if err != nil {
		t.Fatalf("build v3 data collection: %v", err)
	}
	got := value.Export().(map[string]any)
	oldNode := anyMap(got["oldNode"])
	node := anyMap(got["node"])
	if node["type"] != oldNode["type"] || node["type"] != "Stack" {
		t.Fatalf("node types old=%#v new=%#v", oldNode, node)
	}
	oldTable := anyMap(anySlice(oldNode["children"])[0])
	table := anyMap(anySlice(node["children"])[0])
	if table["type"] != "DataTable" || oldTable["type"] != "DataTable" {
		t.Fatalf("table nodes old=%#v new=%#v", oldTable, table)
	}
	props := anyMap(table["props"])
	oldProps := anyMap(oldTable["props"])
	if props["getRowKey"] != oldProps["getRowKey"] || props["selectedKey"] != oldProps["selectedKey"] || props["emptyMessage"] != oldProps["emptyMessage"] {
		t.Fatalf("table props old=%#v new=%#v", oldProps, props)
	}
	if len(anySlice(props["columns"])) != len(anySlice(oldProps["columns"])) {
		t.Fatalf("columns old=%#v new=%#v", oldProps["columns"], props["columns"])
	}
}

func TestWidgetV3DataMasterDetailEditAndMatrix(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const rows = [{ id: "a", title: "Alpha", status: "draft", cells: { mon: "yes" } }];
		const schema = widget.data.fields(f => f.key("id").primary("title").status("status")).build();
		const detail = widget.data.collection(rows, c => c
			.id("articles")
			.schema(schema)
			.select(widget.data.selection.urlParam("article", "a"))
			.edit(e => e
				.create({ label: "New article" })
				.submit("/articles")
				.reorder(widget.act.server("reorder", { payload: { id: widget.bind.context("row.id") } }))
				.remove(widget.act.server("remove", { confirm: "Remove?" }))
			)
			.masterDetail()
		).toNode();
		const matrix = widget.data.matrix(rows, m => m
			.id("availability")
			.column("mon", "Monday")
			.valueAt(widget.bind.map("cells"))
			.cell(widget.data.cell.cycle("availability"))
			.onCellAction(widget.act.server("set-cell"))
		).toNode();
		({ detail, matrix });
	`)
	if err != nil {
		t.Fatalf("build v3 data master-detail/matrix: %v", err)
	}
	got := value.Export().(map[string]any)
	detail := anyMap(got["detail"])
	if detail["type"] != "Stack" || len(anySlice(detail["children"])) < 2 {
		t.Fatalf("detail node = %#v, want Stack with table/detail", detail)
	}
	matrix := anyMap(got["matrix"])
	if matrix["type"] != "MatrixGrid" {
		t.Fatalf("matrix node = %#v, want MatrixGrid", matrix)
	}
	props := anyMap(matrix["props"])
	if len(anySlice(props["columns"])) != 1 || props["valueAt"] == nil || props["onCellAction"] == nil {
		t.Fatalf("matrix props = %#v", props)
	}
}

func TestWidgetV3AccessorsSelectionsItemsActionsAndValidation(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const widget = require("widget.dsl");
		const invalid = widget.page("Invalid", p => p.section("Broken", s => s.view({ kind: "component", props: {} })));
		({
			issues: invalid.validate(),
			field: widget.bind.field("title"),
			path: widget.bind.path("author.name"),
			map: widget.bind.map("cells"),
			template: widget.bind.template("${first} ${last}"),
			context: widget.bind.context("row.id"),
			constant: widget.bind.const(7),
			selection: widget.data.selection({ mode: "multi", keyField: "id", selected: ["a", "b"] }),
			singleSelection: widget.data.selection("single", { keyField: "id", selected: "a" }),
			item: widget.data.item("home", "Home", { href: "/", badge: "New", action: widget.act.navigate("/") }),
			action: widget.act.server("save", {
				confirm: "Save row?",
				payload: { id: widget.bind.context("row.id"), fixed: widget.bind.const(7) },
			}),
		});
	`)
	if err != nil {
		t.Fatalf("build widget.dsl core specs: %v", err)
	}
	got := value.Export().(map[string]any)
	issues := anySlice(got["issues"])
	if len(issues) != 1 {
		t.Fatalf("validation issues = %#v, want one component type issue", issues)
	}
	issue := issues[0].(map[string]any)
	if issue["code"] != "component_type_required" {
		t.Fatalf("validation issue = %#v, want component_type_required", issue)
	}
	field := got["field"].(map[string]any)
	if field["kind"] != "accessor" || field["mode"] != "field" || field["field"] != "title" {
		t.Fatalf("field accessor = %#v", field)
	}
	path := got["path"].(map[string]any)
	if path["kind"] != "accessor" || path["mode"] != "path" || path["path"] != "author.name" {
		t.Fatalf("path accessor = %#v", path)
	}
	mapAccessor := got["map"].(map[string]any)
	if mapAccessor["kind"] != "accessor" || mapAccessor["mode"] != "map" || mapAccessor["mapField"] != "cells" {
		t.Fatalf("map accessor = %#v", mapAccessor)
	}
	constant := got["constant"].(map[string]any)
	if constant["kind"] != "const" || constant["value"] != int64(7) {
		t.Fatalf("const binding = %#v", constant)
	}
	selection := got["selection"].(map[string]any)
	if selection["kind"] != "selection" || selection["mode"] != "multi" || selection["keyField"] != "id" {
		t.Fatalf("selection = %#v", selection)
	}
	singleSelection := got["singleSelection"].(map[string]any)
	if singleSelection["kind"] != "selection" || singleSelection["mode"] != "single" || singleSelection["selected"] != "a" {
		t.Fatalf("single selection = %#v", singleSelection)
	}
	item := got["item"].(map[string]any)
	if item["kind"] != "listItem" || item["id"] != "home" || item["href"] != "/" {
		t.Fatalf("list item = %#v", item)
	}
	action := got["action"].(map[string]any)
	payload := action["payload"].(map[string]any)
	payloadID := payload["id"].(map[string]any)
	if action["kind"] != "server" || action["confirm"] != "Save row?" || payloadID["kind"] != "accessor" || payloadID["mode"] != "context" {
		t.Fatalf("action with payload bindings = %#v", action)
	}
}

func TestOldBucketModulesAreAbsent(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		function canRequire(name) {
			try { require(name); return true; } catch (error) { return false; }
		}
		({ rag: canRequire("rag.dsl") });
	`)
	if err != nil {
		t.Fatalf("check old modules: %v", err)
	}
	got := value.Export().(map[string]any)
	if got["rag"] != false {
		t.Fatalf("old bucket module rag.dsl should be absent, got %#v", got)
	}
}

func TestBuildsWidgetIRAcrossSplitModules(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const ui = require("ui.dsl");
		const data = require("data.dsl");
		const contextWindow = require("context_window.dsl");
		const course = require("course.dsl");
		const styleSet = contextWindow.styleSet({ legend: [contextWindow.legendItem("prompt", "Prompt")], styles: { prompt: contextWindow.visualStyle({ pattern: "checker", fill: "#dde6f2", line: "#4f74a8" }) } });
		const snapshot = { id: "ctx", title: "Window", limit: 1000, parts: [contextWindow.contextPart("p", "Prompt", "prompt", 300)] };
		const slide = { id: "s1", number: "01", title: "Window", view: "budget", snapshotId: "ctx", notes: ["Budget"] };
		const page = ui.page({
			id: "split",
			title: "Split modules",
			sections: [
				ui.formPanel({ title: "Settings", method: "post", formAction: "/settings", submitLabel: "Save" },
					ui.formRow({
						label: "Display name",
						required: true,
						hint: "Shown on uploads",
						control: ui.textInput({ name: "displayName", defaultValue: "Manuel", readOnly: false, maxLength: 40 })
					})
				),
				ui.panel({ title: "Table" }, data.dataTable({
					rows: [{ id: "a", title: "Alpha", status: "done" }],
					getRowKey: "id",
					columns: [
						{ id: "title", header: "Title", cell: data.cell.field("title") },
						{ id: "status", header: "Status", cell: data.cell.status("status", { icon: true }) }
					]
				})),
				contextWindow.contextDiagramPanel({ snapshot, styleSet, initialView: "budget" }),
				course.courseStudioShell({
					sections: [{ id: "course", label: "Course", items: [{ id: "slides", label: "Slides", icon: course.contextStudioNavIcon({ id: "slides" }) }] }],
					activeItemId: "slides",
					title: "Studio"
				}, course.courseSlidePanel({ slide, snapshot, index: 0, total: 1 }))
			]
		});
		JSON.stringify(page);
	`)
	if err != nil {
		t.Fatalf("build split-module page: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(value.String()), &decoded); err != nil {
		t.Fatalf("split-module page is not JSON: %v\n%s", err, value.String())
	}
	assertString(t, decoded, "id", "split")
	root := decoded["root"].(map[string]any)
	children := root["children"].([]any)
	assertString(t, children[0].(map[string]any), "type", "FormPanel")
	assertString(t, children[1].(map[string]any), "type", "Panel")
	assertString(t, children[2].(map[string]any), "type", "ContextDiagramPanel")
	assertString(t, children[3].(map[string]any), "type", "CourseStudioShell")
}

func TestSplitModuleRecipesAreJSONSerializable(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const ui = require("ui.dsl");
		const data = require("data.dsl");
		const contextWindow = require("context_window.dsl");
		const course = require("course.dsl");
		const rows = [{ id: 1, name: "Alpha", status: "running" }];
		const styleSet = contextWindow.styleSet({ legend: [contextWindow.legendItem("prompt", "Prompt")], styles: { prompt: contextWindow.visualStyle({ pattern: "checker", fill: "#dde6f2", line: "#4f74a8" }) } });
		const snapshot = { id: "ctx", title: "Window", limit: 1000, parts: [contextWindow.contextPart("p", "Prompt", "prompt", 300)] };
		const transcript = { title: "Session", messages: [{ id: "m1", role: "user", text: "hello" }], annotations: [] };
		const slide = { id: "s1", number: "01", title: "Window", view: "budget", snapshotId: "ctx", notes: ["Budget"] };
		const bundle = { intro: "Docs", docs: [{ id: "d1", title: "Guide", file: "guide.md", format: "markdown", description: "Guide", body: "# Guide" }] };
		const sections = [{ id: "course", label: "Course", items: [{ id: "slides", label: "Slides" }] }];
		const page = ui.page({ id: "recipes", title: "Recipes", sections: [
			ui.recipes.metrics({ items: [{ label: "Total", value: rows.length, status: "ready" }] }),
			ui.recipes.actionToolbar({ title: "Controls", actions: [{ label: "Add", action: ui.action.server("add-query") }] }),
			data.recipes.masterDetailTable({
				rows,
				columns: [{ id: "name", header: "Name", cell: data.cell.field("name") }],
				selectedKey: 1,
				detail: row => ui.panel({ title: "Selected" }, row.name)
			}),
			contextWindow.recipes.contextDiagram({ snapshot, styleSet, view: "budget" }),
			contextWindow.recipes.annotatedTranscript({ transcript, onAnnotationSelect: contextWindow.action.server("select-annotation") }),
			course.recipes.courseStudio({ sections, activeItemId: "slides", main: course.recipes.courseSlide({ slide, snapshot, index: 0, total: 1 }) }),
			course.recipes.handout({ bundle, selectedDocumentId: "d1", onSelect: course.action.server("select-doc") })
		]});
		JSON.stringify(page);
	`)
	if err != nil {
		t.Fatalf("build split-module recipes: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(value.String()), &decoded); err != nil {
		t.Fatalf("split-module recipe page is not JSON: %v\n%s", err, value.String())
	}
	root := decoded["root"].(map[string]any)
	children := root["children"].([]any)
	if len(children) != 7 {
		t.Fatalf("recipe children len = %d, want 7: %#v", len(children), children)
	}
	assertString(t, children[2].(map[string]any), "type", "DashboardGrid")
	assertString(t, children[3].(map[string]any), "type", "ContextDiagramPanel")
	assertString(t, children[5].(map[string]any), "type", "CourseStudioShell")
}

func TestContextWindowStyleSetHelpersBuildExpectedShape(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const contextWindow = require("context_window.dsl");
		const styleSet = contextWindow.paletteStyleSet({
			palette: "Signal Orange / Cyan",
			entries: [
				{ id: "prompt", label: "Prompt", accent: "b", pattern: "checker" },
				{ id: "evidence", label: "Evidence", accent: "a", pattern: "stipple" },
				{ id: "answer", label: "Answer", accent: "a", pattern: "solid" },
			]
		});
		const snapshot = contextWindow.contextSnapshot({
			id: "ctx",
			title: "Window",
			limit: 1000,
			parts: [contextWindow.contextPart("p", "Prompt", "prompt", 300)]
		});
		JSON.stringify({ styleSet, snapshot });
	`)
	if err != nil {
		t.Fatalf("build style helpers: %v", err)
	}
	var decoded map[string]any
	if err := json.Unmarshal([]byte(value.String()), &decoded); err != nil {
		t.Fatalf("style helper output is not JSON: %v", err)
	}
	styleSet := decoded["styleSet"].(map[string]any)
	styles := styleSet["styles"].(map[string]any)
	if _, ok := styles["prompt"].(map[string]any); !ok {
		t.Fatalf("styleSet.styles.prompt missing: %#v", styleSet)
	}
	snapshot := decoded["snapshot"].(map[string]any)
	parts := snapshot["parts"].([]any)
	part := parts[0].(map[string]any)
	assertString(t, part, "styleKey", "prompt")
	if _, hasKind := part["kind"]; hasKind {
		t.Fatalf("contextPart emitted forbidden kind field: %#v", part)
	}
}

func TestContextDiagramRecipeRequiresStyleSet(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	_, err := vm.RunString(`
		const contextWindow = require("context_window.dsl");
		const snapshot = { id: "ctx", title: "Window", limit: 1000, parts: [] };
		contextWindow.recipes.contextDiagram({ snapshot, view: "budget" });
	`)
	if err == nil {
		t.Fatalf("contextDiagram recipe without styleSet should fail")
	}
	if !strings.Contains(err.Error(), "requires styleSet") {
		t.Fatalf("error = %v, want useful styleSet message", err)
	}
}

func TestEngineRegistrarRegistersSplitModulesOnly(t *testing.T) {
	factory, err := engine.NewRuntimeFactoryBuilder().WithModules(NewRegistrar()).Build()
	if err != nil {
		t.Fatalf("build runtime factory: %v", err)
	}
	rt, err := factory.NewRuntime()
	if err != nil {
		t.Fatalf("create runtime: %v", err)
	}
	defer func() { _ = rt.Close(context.Background()) }()

	value, err := rt.VM.RunString(`
		function canRequire(name) {
			try { require(name); return true; } catch (error) { return false; }
		}
		const ui = require("ui.dsl");
		const data = require("data.dsl");
		const contextWindow = require("context_window.dsl");
		const course = require("course.dsl");
		({
			uiPanel: typeof ui.panel,
			dataTable: typeof data.dataTable,
			contextDiagramPanel: typeof contextWindow.contextDiagramPanel,
			courseStudioShell: typeof course.courseStudioShell,
			widget: canRequire("widget.dsl"),
			rag: canRequire("rag.dsl"),
		});
	`)
	if err != nil {
		t.Fatalf("require split modules through engine registrar: %v", err)
	}
	got := value.Export().(map[string]any)
	for _, name := range []string{"uiPanel", "dataTable", "contextDiagramPanel", "courseStudioShell"} {
		if got[name] != "function" {
			t.Fatalf("%s export = %#v, want function (all: %#v)", name, got[name], got)
		}
	}
	if got["widget"] != true || got["rag"] != false {
		t.Fatalf("widget.dsl should be present and rag.dsl should be absent from engine registrar, got %#v", got)
	}
}

func anyMap(value any) map[string]any {
	if m, ok := value.(map[string]any); ok {
		return m
	}
	out := map[string]any{}
	for key, value := range exportAnyMap(value) {
		out[key] = value
	}
	return out
}

func exportAnyMap(value any) map[string]any {
	out := map[string]any{}
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Map {
		return out
	}
	for _, key := range rv.MapKeys() {
		if key.Kind() == reflect.String {
			out[key.String()] = rv.MapIndex(key).Interface()
		}
	}
	return out
}

func exportedSliceLen(value any) int {
	switch v := value.(type) {
	case []any:
		return len(v)
	case []map[string]any:
		return len(v)
	default:
		return -1
	}
}

func assertString(t *testing.T, m map[string]any, key, want string) {
	t.Helper()
	if got, _ := m[key].(string); got != want {
		t.Fatalf("%s = %#v, want %q (map=%#v)", key, m[key], want, m)
	}
}

func TestCmsModuleExportsHelpersRecipesAndBoundaries(t *testing.T) {
	vm := goja.New()
	reg := require.NewRegistry()
	Register(reg)
	reg.Enable(vm)

	value, err := vm.RunString(`
		const ui = require("ui.dsl");
		const contextWindow = require("context_window.dsl");
		const cms = require("cms.dsl");
		const library = cms.recipes.mediaLibrary({
			assets: [{ id: "a", kind: "image", title: "A", filename: "a.png", mime: "image/png", size: 10, src: "/course-assets/a.png", tags: [], status: "published", createdAt: "", updatedAt: "" }],
			selectedAssetIds: ["a"],
			onFilesSelected: "admin-upload-course-material",
			onAssetSelect: cms.action.navigate("?asset=$assetId"),
		});
		const list = cms.recipes.articleList({
			articles: [{ id: "x", slug: "x", title: "X", status: "draft", tags: [], updatedAt: "" }],
			onRowAction: { kind: "event", event: "row-action", confirm: "Really?" },
		});
		({
			cmsMediaLibraryPanel: typeof cms.mediaLibraryPanel,
			cmsTag: typeof cms.tag,
			cmsMediaThumb: typeof cms.mediaThumb,
			cmsMarkdownEditor: typeof cms.markdownEditor,
			cmsActionServer: typeof cms.action.server,
			cmsPage: typeof cms.page,
			cmsCell: typeof cms.cell,
			uiMediaLibraryPanel: typeof ui.mediaLibraryPanel,
			contextGroupedStrip: typeof contextWindow.contextGroupedStripDiagram,
			libraryType: library.type,
			libraryUploadKind: library.props.onFilesSelectedAction.kind,
			libraryUploadName: library.props.onFilesSelectedAction.name,
			librarySelectKind: library.props.onAssetSelectAction.kind,
			listType: list.type,
			listRowActionConfirm: list.props.onRowActionAction.confirm,
		});
	`)
	if err != nil {
		t.Fatalf("require cms.dsl: %v", err)
	}
	got := value.Export().(map[string]any)
	wantFunctions := []string{"cmsMediaLibraryPanel", "cmsTag", "cmsMediaThumb", "cmsMarkdownEditor", "cmsActionServer", "contextGroupedStrip"}
	for _, name := range wantFunctions {
		if got[name] != "function" {
			t.Fatalf("%s = %#v, want function (all: %#v)", name, got[name], got)
		}
	}
	wantUndefined := []string{"cmsPage", "cmsCell", "uiMediaLibraryPanel"}
	for _, name := range wantUndefined {
		if got[name] != "undefined" {
			t.Fatalf("%s = %#v, want undefined (all: %#v)", name, got[name], got)
		}
	}
	if got["libraryType"] != "MediaLibraryPanel" {
		t.Fatalf("libraryType = %#v, want MediaLibraryPanel", got["libraryType"])
	}
	if got["libraryUploadKind"] != "server" || got["libraryUploadName"] != "admin-upload-course-material" {
		t.Fatalf("upload action = %#v/%#v, want server/admin-upload-course-material", got["libraryUploadKind"], got["libraryUploadName"])
	}
	if got["librarySelectKind"] != "navigate" {
		t.Fatalf("select action kind = %#v, want navigate", got["librarySelectKind"])
	}
	if got["listType"] != "ArticleListPanel" {
		t.Fatalf("listType = %#v, want ArticleListPanel", got["listType"])
	}
	if got["listRowActionConfirm"] != "Really?" {
		t.Fatalf("confirm passthrough = %#v, want Really?", got["listRowActionConfirm"])
	}
}
