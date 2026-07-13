package widgetdsl

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dop251/goja"
	"github.com/dop251/goja_nodejs/require"
	"github.com/go-go-golems/go-go-goja/pkg/tsgen/render"
	"github.com/go-go-golems/go-go-goja/pkg/tsgen/spec"
)

func TestWidgetV3TypeScriptFixtureCompilesExamples(t *testing.T) {
	repoRoot := findRepoRoot(t)
	tscPath := os.Getenv("WIDGETDSL_TSC")
	if tscPath == "" {
		tscPath = filepath.Join(repoRoot, "packages", "rag-evaluation-site", "node_modules", ".bin", "tsc")
	}
	if _, err := os.Stat(tscPath); err != nil {
		t.Skipf("TypeScript compiler not installed at %s; run pnpm install first", tscPath)
	}

	dts, err := render.Bundle(&spec.Bundle{Modules: []*spec.Module{TypeScriptModule(WidgetV3ModuleName)}})
	if err != nil {
		t.Fatalf("render widget.dsl DTS: %v", err)
	}

	tmp := t.TempDir()
	dtsPath := filepath.Join(tmp, "widgetdsl.d.ts")
	fixturePath := filepath.Join(tmp, "widget-v3-fixture.ts")
	if err := os.WriteFile(dtsPath, []byte(dts), 0o644); err != nil {
		t.Fatalf("write DTS fixture: %v", err)
	}
	if err := os.WriteFile(fixturePath, []byte(widgetV3TypeScriptFixture), 0o644); err != nil {
		t.Fatalf("write TS fixture: %v", err)
	}

	cmd := exec.Command(tscPath,
		"--strict",
		"--noEmit",
		"--target", "ES2022",
		"--module", "NodeNext",
		"--moduleResolution", "NodeNext",
		"--skipLibCheck",
		fixturePath,
	)
	cmd.Dir = tmp
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("widget.dsl typescript fixture failed: %v\n%s", err, output)
	}
}

func TestDataV2TypeScriptFixtureCompilesPositiveAndExpectedNegativeExamples(t *testing.T) {
	repoRoot := findRepoRoot(t)
	tscPath := os.Getenv("WIDGETDSL_TSC")
	if tscPath == "" {
		tscPath = filepath.Join(repoRoot, "packages", "rag-evaluation-site", "node_modules", ".bin", "tsc")
	}
	if _, err := os.Stat(tscPath); err != nil {
		t.Skipf("TypeScript compiler not installed at %s; run pnpm install first", tscPath)
	}

	dts, err := render.Bundle(&spec.Bundle{Modules: []*spec.Module{TypeScriptModule(DataV2ModuleName)}})
	if err != nil {
		t.Fatalf("render data.v2.dsl DTS: %v", err)
	}

	tmp := t.TempDir()
	dtsPath := filepath.Join(tmp, "widgetdsl.d.ts")
	fixturePath := filepath.Join(tmp, "data-v2-fixture.ts")
	if err := os.WriteFile(dtsPath, []byte(dts), 0o644); err != nil {
		t.Fatalf("write DTS fixture: %v", err)
	}
	if err := os.WriteFile(fixturePath, []byte(dataV2TypeScriptFixture), 0o644); err != nil {
		t.Fatalf("write TS fixture: %v", err)
	}

	cmd := exec.Command(tscPath,
		"--strict",
		"--noEmit",
		"--target", "ES2022",
		"--module", "NodeNext",
		"--moduleResolution", "NodeNext",
		"--skipLibCheck",
		fixturePath,
	)
	cmd.Dir = tmp
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("typescript fixture failed: %v\n%s", err, output)
	}
}

func findRepoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	for dir := wd; ; dir = filepath.Dir(dir) {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatalf("could not find repo root from %s", wd)
		}
	}
}

const widgetV3TypeScriptFixture = `/// <reference path="./widgetdsl.d.ts" />
import * as widget from "widget.dsl";

const page = widget.page("Dashboard", p => p
  .section("Overview", s => s
    .view(widget.ui.card({ title: "Hello" }, "World"))));
page.toPage();

const collection = widget.data.collection("rows", [{ id: "a", title: "Alpha" }], c => c
  .schema(widget.data.fields("rows", f => f.key("id").primary("title")).build())
  .table(t => t.rowSelect(widget.act.event("row.select"))));
collection.toNode();

widget.cms.mediaLibrary([{ id: "asset-1", title: "Hero" }], m => m
  .asset((asset, h) => h.card({ title: asset.title }, asset.id))
  .onSelect(widget.cms.intent.selectAsset(widget.bind.context("asset.id"))));

widget.course.shell({ title: "Course", sections: [{ id: "intro", items: [{ id: "start", label: "Start" }] }] }, c => c
  .active("start")
  .onNavigate(widget.course.intent.navigate(widget.bind.context("item.id"))));

widget.course.handouts({ docs: [{ id: "h1", title: "Handout" }] }, h => h
  .selected("h1")
  .onDownload(widget.course.intent.downloadHandout(widget.bind.context("doc.id"))));

widget.context.workspace({
  title: "Session",
  messages: [{ id: "m1", role: "user", content: "Hello" }],
  annotations: [{ id: "a1", messageId: "m1", note: "Important" }],
}, w => w
  .message((message, h) => h.caption(message.content))
  .annotation((annotation, h) => h.caption(annotation.note))
  .onAnnotationSelect(widget.context.intent.selectAnnotation(widget.bind.context("annotation.id"))));

const poll = {
  title: "Availability",
  options: [{ id: "mon-9", label: "Mon 9", startISO: "2026-07-08T09:00:00Z", endISO: "2026-07-08T10:00:00Z" }],
  responses: [{ id: "ana", name: "Ana", availability: { "mon-9": "available" } }],
};
widget.schedule.availabilityPoll(poll, p => p.readOnly());
widget.time.week([{ id: "ev1", title: "Launch", startISO: "2026-07-08T09:00:00Z", endISO: "2026-07-08T10:00:00Z" }], w => w
  .range(widget.time.range.week("2026-07-08"))
  .onSelect(widget.time.intent.selectEvent(widget.bind.context("block.id"))));
`

const dataV2TypeScriptFixture = `/// <reference path="./widgetdsl.d.ts" />
import * as data from "data.v2.dsl";

const rows = [
  { id: "sess-intro", title: "Intro", turns: 12, status: "ready" },
  { id: "sess-debug", title: "Debugging", turns: 28, status: "ready" },
];

const schema = data.schema("sessions")
  .field("id", data.f.key().label("ID").readOnly())
  .field("title", data.f.primary().label("Title").required().maxLength(120))
  .field("turns", data.f.count().label("Turns").width("6rem"))
  .field("status", data.f.status().label("Status"))
  .build();

const tableNode = data.collection("sessions", rows)
  .schema(schema)
  .empty("No sessions yet.")
  .select((selection) => selection.urlParam("selected", "sess-intro"))
  .table((table) => table
    .className("sessions-table")
    .rowSelect(data.action.navigate("?selected=${row.id}"))
    .actionColumn("open", "Open", "Open", data.action.navigate("/sessions/${row.id}"), { maxWidth: "8ch" }))
  .toIR();

tableNode.kind;

const editorNode = data.collection("agenda", rows)
  .schema(schema)
  .edit((editor) => editor
    .selectUrl("agenda", "sess-intro")
    .submitPost("/settings/demo")
    .create({ label: "New item" })
    .actions((actions) => actions
      .reorder(data.action.server("reorder").payloadPath("id", "row.id"))
      .remove(data.action.server("delete").confirm("Delete ${row.title}?"))))
  .masterDetail()
  .toIR();

editorNode.kind;

// @ts-expect-error data.v2.dsl does not expose the v1 option-bag dataTable helper.
data.dataTable({ rows });

// @ts-expect-error schema() requires a schema name.
data.schema();

// @ts-expect-error schema.field() requires a typed FieldHandle, not a raw object.
data.schema("bad").field("id", { kind: "string", semantic: "key" });

// @ts-expect-error collection rows must be object records.
data.collection("bad", [1, 2, 3]);

// @ts-expect-error rowSelect() requires a typed ActionHandle, not raw action JSON.
data.collection("bad", rows).table((table) => table.rowSelect({ kind: "navigate", to: "/x" }));
` + "\n"

func newRegisteredRuntime(t *testing.T) *goja.Runtime {
	t.Helper()
	vm := goja.New()
	reg := require.NewRegistry()
	registerLegacyModulesForTests(reg)
	reg.Enable(vm)
	return vm
}

func TestDataV2RuntimeExportsMatchDeclaredPublicSurface(t *testing.T) {
	dts := strings.Join(TypeScriptModule(DataV2ModuleName).RawDTS, "\n")
	for _, declared := range []string{
		"export const f: FieldFactory;",
		"export function schema(name: string): SchemaBuilder;",
		"export function collection(name: string, rows: Record<string, JsonValue>[]): CollectionBuilder;",
		"export const selection:",
		"export const action: ActionFactory;",
	} {
		if !strings.Contains(dts, declared) {
			t.Fatalf("DTS missing declared public surface %q", declared)
		}
	}

	vm := newRegisteredRuntime(t)
	value, err := vm.RunString(`
		const data = require("data.v2.dsl");
		({
			f: typeof data.f,
			schema: typeof data.schema,
			collection: typeof data.collection,
			selection: typeof data.selection,
			action: typeof data.action,
			dataTable: typeof data.dataTable,
			cell: typeof data.cell,
		});
	`)
	if err != nil {
		t.Fatalf("inspect data.v2.dsl exports: %v", err)
	}
	got := value.Export().(map[string]any)
	for _, name := range []string{"f", "selection", "action"} {
		if got[name] != "object" {
			t.Fatalf("data.v2.dsl export %s = %#v, want object (all: %#v)", name, got[name], got)
		}
	}
	for _, name := range []string{"schema", "collection"} {
		if got[name] != "function" {
			t.Fatalf("data.v2.dsl export %s = %#v, want function (all: %#v)", name, got[name], got)
		}
	}
	for _, name := range []string{"dataTable", "cell"} {
		if got[name] != "undefined" {
			t.Fatalf("data.v2.dsl legacy export %s = %#v, want undefined (all: %#v)", name, got[name], got)
		}
	}
}
