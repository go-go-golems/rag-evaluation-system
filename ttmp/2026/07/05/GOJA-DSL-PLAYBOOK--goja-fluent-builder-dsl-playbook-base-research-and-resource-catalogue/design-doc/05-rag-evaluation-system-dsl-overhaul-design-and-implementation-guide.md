---
Title: Rag Evaluation System DSL Overhaul Design and Implementation Guide
Ticket: GOJA-DSL-PLAYBOOK
Status: active
Topics:
    - goja
    - dsl
    - fluent-builder
    - go
    - typescript
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js
      Note: Real course/CMS consumer and first target hard-cutover rewrite page
    - Path: ../../../../../../../go-go-course/cmd/go-go-course/server.js
      Note: Real xgoja host serving pages and dispatching /api/widget/actions/:name
    - Path: packages/rag-evaluation-site/src/app/App.tsx
      Note: App-level server-action POST and refresh behavior
    - Path: packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx
      Note: Renderer registry/adaptor execution model that v2 IR must preserve
    - Path: packages/rag-evaluation-site/src/widgets/actions.ts
      Note: Current browser action serialization/interpolation/dispatch boundary
    - Path: packages/rag-evaluation-site/src/widgets/ir.ts
      Note: Typed React-side Widget IR and action/prop contracts
    - Path: pkg/widgetdsl/grammar.go
      Note: Current data/ui grammar vocabulary and Pattern C limitations
    - Path: pkg/widgetdsl/module.go
      Note: Current split module/helper/action/recipe implementation and map-IR helpers
    - Path: pkg/widgetdsl/typescript.go
      Note: Current weak TypeScript declaration generator to replace
ExternalSources: []
Summary: 'Intern-facing hard-cutover overhaul guide for all rag-evaluation-system Widget DSLs: current architecture, target typed intent IR, fluent/lambda API, action serialization/runtime dispatch design, replacement phases, tests, and many simple-to-rich examples for course/CMS/general pages.'
LastUpdated: 2026-07-05T17:55:00-04:00
WhatFor: Use as the implementation guide when replacing the current map-IR Widget DSL family with a hard-cutover v2 typed builder system that intentionally does not preserve backwards compatibility.
WhenToUse: Read before implementing the hard-cutover widgetdsl v2, changing the React Widget IR renderer contract, or adding new course/CMS/context-window DSL APIs.
---


# Rag Evaluation System DSL Overhaul Design and Implementation Guide

## 0. Reading guide for the intern

This document is the handoff for a future overhaul of **all** `rag-evaluation-system` Widget DSLs:

- `ui.dsl`
- `data.dsl`
- `context_window.dsl`
- `course.dsl`
- `cms.dsl`
- the shared cell/action/page/recipe helpers that those modules expose
- the React `WidgetRenderer` and Widget IR shape when the DSL needs a better serialized target
- the browser renderer ⇄ Goja/xgoja server-action contract

Read it as a design and implementation guide for a **hard-cutover v2**, not as a backwards-compatible migration plan. The document still shows competing API shapes so Manuel can decide which forms feel best, but the implementation guidance now assumes we can delete old public APIs and choose the simplest durable language rather than preserve v1 call sites.

The four companion documents that this guide synthesizes are:

1. `design-doc/01-goja-dsl-catalogue-and-base-research.md` — catalogue of Goja DSL patterns: typed refs, builder structs, map IR, hyperscript, lambda configurators, DTS parity, Proxy builders.
2. `design-doc/02-self-assessment-of-the-widgetdsl-grammar-what-pattern-c-actually-costs-and-what-the-playbook-should-add.md` — insider assessment: the current grammar vocabulary is good; the map substrate is not.
3. `design-doc/03-widget-dsl-design-assessment-and-improvement-report.md` — independent assessment: use the current widget grammar as vocabulary evidence, not as canonical architecture.
4. `design-doc/04-goja-dsl-deep-dive-optional-lambdas-typed-ir-dts-parity-and-tag-operators.md` — deep dive on optional lambdas, typed specs, Geppetto DTS parity, and go-emrichen operator composition.

The short version is:

> Keep the intent vocabulary. Replace the substrate. Do a hard cutover. Use typed Go-side intent specs and fluent builders. Keep output serializable. Treat lambdas as builder-time configurators or registered server handlers, never as functions serialized to the browser.

## 1. Executive summary

The current Widget DSL family works because it is simple: each helper returns JSON-like Widget IR maps, and the React `WidgetRenderer` renders those maps. That simplicity enabled a useful first grammar (`data.schema`, `data.collection`, `ui.section`) and a real page improvement in `go-go-course`. It also created the exact failure modes that become dangerous as the system grows: silent option typos, weak TypeScript declarations, magic marker objects, stringly arrangements/actions, and runtime panics inside request-scoped page builders.

The overhaul should make the DSLs **powerful but opinionatedly simple**:

- Simple pages should be one or two readable sentences.
- Rich pages should use optional lambda configurators and reusable fragments.
- All high-level DSL concepts should accumulate into typed Go-side specs.
- The terminal output should still be serializable Widget IR.
- The React renderer should remain generic, but the IR can evolve to encode actions, bindings, diagnostics, and intent-compiled widgets more precisely.
- JavaScript lambdas should configure builders on the server/runtime side; any browser-visible action must compile to data.
- Server callbacks should be registered under stable names and invoked by the host action endpoint, not smuggled through JSON as closures.

Recommended target architecture:

```text
┌──────────────────────────────────────────────────────────────────────┐
│ JS authoring DSLs                                                     │
│ ui.dsl, data.dsl, context_window.dsl, course.dsl, cms.dsl             │
│                                                                      │
│ Opinionated simple helpers + optional lambdas + .use(fragment)        │
└──────────────────────────────────────────────────────────────────────┘
                         │ builder calls mutate typed specs
                         ▼
┌──────────────────────────────────────────────────────────────────────┐
│ Go-side typed intent specs                                            │
│ PageSpec, NodeSpec, SchemaSpec, CollectionSpec, ActionSpec, MarkSpec  │
│ ValidationIssue[] with path/code/message/severity                     │
└──────────────────────────────────────────────────────────────────────┘
                         │ .validate(), .toIR(), page()
                         ▼
┌──────────────────────────────────────────────────────────────────────┐
│ Serialized Widget IR v2                                               │
│ JSON-compatible nodes + typed action descriptors + optional issues    │
└──────────────────────────────────────────────────────────────────────┘
                         │ HTTP JSON
                         ▼
┌──────────────────────────────────────────────────────────────────────┐
│ React renderer                                                        │
│ WidgetRenderer + adapters + central action dispatcher                 │
└──────────────────────────────────────────────────────────────────────┘
                         │ server action POST {payload, context}
                         ▼
┌──────────────────────────────────────────────────────────────────────┐
│ Host/xgoja runtime                                                    │
│ /api/widget/actions/:name dispatches stable registered handlers       │
└──────────────────────────────────────────────────────────────────────┘
```

The main design decision is to separate three concepts that are currently conflated:

1. **Authoring API** — what page authors write in JavaScript.
2. **Intent spec** — typed Go-side state that can validate and compile.
3. **Widget IR** — serializable JSON sent to React.

The current code uses Widget IR as the authoring API. The hard-cutover overhaul should use typed intent specs as the authoring substrate and Widget IR as the terminal output. No v2 public API should require authors to construct raw Widget IR maps for high-level concepts.

## 2. Current system map

### 2.1 The modules and what they expose today

The split modules are declared in `pkg/widgetdsl/module.go:14-20`:

```go
const (
    UIModuleName            = "ui.dsl"
    DataModuleName          = "data.dsl"
    ContextWindowModuleName = "context_window.dsl"
    CourseModuleName        = "course.dsl"
    CmsModuleName           = "cms.dsl"
)
```

The helper maps are hard-coded in `pkg/widgetdsl/module.go`:

- `ui.dsl` has 31 generic layout/foundation helpers such as `appShell`, `button`, `caption`, `fieldGrid`, `formPanel`, `sectionBlock`, `stack`, `textInput`, and `textareaInput` (`module.go:34-79`).
- `data.dsl` has `dataTable` and a `cell` namespace (`module.go:81-83`, `module.go:264-312`).
- `context_window.dsl` has context, transcript, annotation, and upload helpers (`module.go:85-106`).
- `course.dsl` has course, slide, handout, and shell helpers (`module.go:125-137`).
- `cms.dsl` has media/article/admin helpers (`module.go:108-123`).
- All modules except `data.dsl` receive the shared `action` namespace if `action: true` is set; in practice all five module specs enable it (`module.go:139-177`).

The install path is generic: `runtime.install` always exports `text`, `element`, `component`, and `fragment`; optionally exports `page`, helper factories, `cell`, `action`, context style helpers, data grammar, `ui.section`, and recipes (`module.go:231-260`).

This is flexible, but it is also why the modules are shallow. A helper is mostly just a name-to-component-type mapping.

### 2.2 The authoring substrate today: maps, not builders

The low-level helper path is:

```go
func (r *runtime) componentFactory(componentType string) func(goja.FunctionCall) goja.Value {
    return func(call goja.FunctionCall) goja.Value {
        props, childStart := propsAndChildStart(call.Arguments, 0)
        return r.vm.ToValue(r.buildComponent(componentType, props, call.Arguments[childStart:]))
    }
}
```

Evidence: `pkg/widgetdsl/module.go:565-570`.

`buildComponent` returns a map:

```go
func (r *runtime) buildComponent(componentType string, props map[string]any, childValues []goja.Value) map[string]any {
    out := map[string]any{"kind": "component", "type": componentType}
    if len(props) > 0 { out["props"] = props }
    children := r.exportChildren(childValues)
    if len(children) > 0 { out["children"] = children }
    return out
}
```

Evidence: `pkg/widgetdsl/module.go:916-925`.

That means `ui.button({ variant: "primary" }, "Save")` returns a plain JSON-ish object immediately. There is no typed `ButtonSpec` or `ButtonBuilder` on the Go side.

The grammar layer is the same. `data.schema` returns:

```go
map[string]any{"__ragSchema": true, "fields": fields}
```

Evidence: `pkg/widgetdsl/grammar.go:77-94`. `data.collection` receives that map, extracts fields with `schemaFields`, and then compiles directly to `DataTable`, `FormPanel`, `FieldGrid`, `Stack`, and `SectionBlock` maps (`grammar.go:265-472`).

### 2.3 The React-side Widget IR is already much more typed

The TypeScript side has a real IR definition:

- `WidgetNode = TextNode | ElementNode | ComponentNode` (`packages/rag-evaluation-site/src/widgets/ir.ts:37-49`).
- `RagWidgetType` enumerates concrete widget types (`ir.ts:51-128`).
- `ComponentNode` carries `type`, `props`, and `children` (`ir.ts:141-146`).
- `ActionSpec` is a discriminated union of `navigate`, `download`, `server`, `event`, and `copy` (`ir.ts:150-194`).
- `DataTableWidgetProps`, `DataTableColumnSpec`, and `CellSpec` are strongly typed (`ir.ts:575-666`).
- `WidgetProps` is a union of all public widget prop interfaces (`ir.ts:841+`).

So the renderer side has a stronger type vocabulary than the Goja authoring side. That is a key opportunity: the overhaul does not need to invent all types from scratch. It should move type knowledge left, toward the authoring runtime.

### 2.4 The current TypeScript declarations throw most of that away

The generated `.d.ts` in `pkg/widgetdsl/typescript.go` intentionally emits broad types:

```ts
export interface WidgetNode { kind: string; [key: string]: any; }
export interface WidgetAction { kind: string; [key: string]: any; }
export type Props = Record<string, any>;
```

Evidence: `pkg/widgetdsl/typescript.go:18-31`.

The data grammar declarations are also loose:

```ts
export interface FieldSpec { role: string; [key: string]: any; }
export interface Schema { fields: FieldSpec[]; [key: string]: any; }
export function record(values: Props, options: Props): WidgetNode;
export function collection(rows: Props[], options: Props): WidgetNode;
```

Evidence: `pkg/widgetdsl/typescript.go:76-100`.

This is the root of the agent-authoring problem: LLM agents and human authors get autocomplete, but the types say almost everything is valid.

### 2.5 The renderer/action path today

The renderer path is simple:

```text
WidgetPageResponse JSON
  → RagEvaluationSiteApp.useWidgetPage(...)
  → renderPage(...)
  → WidgetRenderer
  → registry.get(node.type)
  → adapter.render(props, children, ctx, node)
```

Important files:

- `packages/rag-evaluation-site/src/app/App.tsx:43-80` fetches pages and intercepts server actions.
- `packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx:14-71` renders text/element/component nodes.
- `packages/rag-evaluation-site/src/widgets/defaultRegistry.ts:91-173` groups adapters by module and merges them.
- `packages/rag-evaluation-site/src/widgets/registry.ts:7-24` defines `RenderContext`, including `bindAction` and `dispatchAction`.

The action path is:

```text
Widget adapter creates context
  → ctx.dispatchAction(action, context)
  → RagEvaluationSiteApp.handleAction(action, context)
  → if server: POST /api/widget/actions/:name { payload, context }
  → host endpoint dispatches by name
```

The browser action implementation is in `packages/rag-evaluation-site/src/widgets/actions.ts:20-105`:

- `copy` writes to clipboard.
- `event` fires a `CustomEvent`, with special `print` and `fullscreen` cases.
- `navigate` interpolates the target and uses `history.pushState`.
- `download` interpolates the target and clicks a temporary anchor.
- `server` posts `{payload, context}` to `/api/widget/actions/:name` if no app-level handler intercepts it.

The app-level handler in `App.tsx` performs the same server POST using the configured `apiBase` (`App.tsx:61-80`). This is the correct general shape. The missing piece is that the action descriptors and contexts are not typed or validated well enough.

### 2.6 The real host currently dispatches string-named actions

The real `go-go-course` consumer serves Widget pages at `GET /api/widget/pages/:id` and server actions at `POST /api/widget/actions/:name`.

Evidence:

- `go-go-course/cmd/go-go-course/server.js:124-154` serves pages by calling `buildWidgetPage(...)`.
- `server.js:422-526` dispatches action names such as `admin-upload-course-material`, `admin-delete-course-material`, `admin-delete-agenda-item`, `admin-reorder-course-agenda`, and `upload-session`.
- `lib/course-pages.js:1-96` wires environment dependencies into page builders.
- `lib/pages/admin-course-cms.js:96-130` is the first real consumer of `data.schema` and `data.collection`.

The existing runtime/action architecture is good enough to preserve, but not typed enough to scale.

## 3. What must change

### 3.1 Keep the vocabulary, replace the substrate

The RAGEVAL-UI-GRAMMAR vocabulary is worth keeping:

- `ui.section` for flat document structure.
- `data.schema` and `data.f.*` for field intent.
- `data.record` for one object.
- `data.collection` for tables/master-detail/etc.
- URL-backed selection via `data.urlParam`.
- Native form posting via `data.formPost`.
- Server actions via `ui.action.server(...)`.

The self-assessment states the key thesis: **the grammar's language design is worth keeping; its implementation substrate is not.** This guide adopts that thesis.

### 3.2 Stop using Widget IR as the authoring API

Widget IR should be a **wire format**. It is allowed to be JSON-like, because it must cross HTTP and be rendered by React. But the authoring API should be typed builders and specs.

Bad target:

```js
// Looks concise, but every option is a bag and every marker is a map.
data.collection(rows, {
  schema,
  verb: "edit",
  arrange: "master-detail",
  submit: data.formPost("/settings/agenda-item"),
});
```

Better target:

```js
data.collection("agenda", rows)
  .schema(agendaSchema)
  .edit()
  .masterDetail()
  .submitPost("/settings/agenda-item")
  .toIR();
```

Best flexible target:

```js
data.collection("agenda", rows, c => c
  .schema(agendaSchema)
  .edit(e => e
    .select(data.selection.urlParam("agenda", query.agenda))
    .submit(data.form.post("/settings/agenda-item"))
    .actions(adminRowActions))
  .arrange(a => a.masterDetail(md => md
    .summary(s => s.table(t => t.elide("description")))
    .detail(d => d.form()))))
  .toIR();
```

The JSON IR appears only at `.toIR()` or `ui.page(...)`.

### 3.3 Treat lambdas carefully

There are two different meanings of “lambda”:

1. **Builder-time configurator lambda** — runs in the Goja runtime while building a spec. Safe. Not serialized.
2. **Runtime action/callback lambda** — runs later in response to a browser event. Cannot be serialized as JSON. Must be registered on the server/runtime side and referenced by a stable action name or handler id.

Builder-time lambdas are recommended:

```js
data.collection("agenda", rows, c => c.schema(agendaSchema).edit())
```

Runtime callbacks require a registry:

```js
ui.actions.register("admin.deleteAgendaItem", actionContextSchema, async (ctx) => {
  const id = ctx.row.id;
  await courseMetadata.deleteAgendaItem(id);
  return ui.result.refresh(`Deleted ${id}`);
});

ui.action.server("admin.deleteAgendaItem")
```

Do **not** put a raw JavaScript function into the Widget IR. React cannot deserialize it, and even if it could, it would run in the wrong environment.

### 3.4 No backwards compatibility requirement

This guide now assumes a **hard cutover**. The old DSL is research data and vocabulary evidence, not a public contract v2 must preserve. That changes the design pressure: choose the smallest coherent v2 language, not the smoothest path from v1.

Rules for the hard cutover:

- No public `data.collection(rows, { ... })` options-bag grammar.
- No public raw schema marker objects such as `__ragSchema`.
- No `Props = Record<string, any>` for high-level grammar APIs.
- No string `verb` / `arrange` switches as the blessed API. Use methods such as `.edit()` and `.masterDetail()`.
- No raw action maps as the blessed API. Use typed `ActionBuilder` handles and serializable templates.
- No recipes that directly emit giant component maps. Recipes are named fragments over typed builders.
- Raw Widget IR remains available only through an explicit escape hatch such as `widget.unsafe`, and production page code should treat that as a design smell.

The implementation may include temporary local scripts to port known pages, but it should not include runtime shims whose only purpose is to make old v1 page code keep running.

## 4. Target architecture

### 4.1 Packages and responsibilities

Recommended new Go package structure inside `pkg/widgetdsl`:

```text
pkg/widgetdsl/
  module.go                # thin native module registration; no map-IR business logic
  typescript.go            # calls schema/declaration generator; no hand-written any soup
  fluent/
    refs.go                # hidden-key typed refs, getTypedRef[T], mustSet
    callbacks.go           # optional configurator application, strict AssertFunction
    validation.go          # ValidationIssue, ValidationResult, issue paths
  spec/
    widget.go              # WidgetNodeSpec / WidgetPageSpec / serialization
    actions.go             # ActionSpecV2, PayloadTemplate, ContextSchema
    fields.go              # SchemaSpec, FieldSpec, facets, roles
    collections.go         # CollectionSpec, ArrangementSpec, ViewSpec
    forms.go               # FormSpec, FormPostSpec, field editors
    marks.go               # MarkSpec interfaces + built-ins
    modules.go             # module helper registry loaded from manifests
  builders/
    ui.go                  # PageBuilder, SectionBuilder, structure helpers
    data.go                # SchemaBuilder, CollectionBuilder, RecordBuilder
    actions.go             # ActionBuilder, ActionRegistryBuilder
    context_window.go      # style sets, context schemas, marks
    course.go              # course schemas/fragments/layouts
    cms.go                 # cms schemas/fragments/layouts
  compile/
    to_ir.go               # typed intent spec -> Widget IR v1/v2
    diagnostics.go         # validation issues -> ErrorCallout/Panel nodes
  codegen/
    manifests.go           # load *.widget.yaml + ir.ts/generated type metadata
    dts.go                 # precise .d.ts generation
    parity_test.go         # geppetto-style runtime export parity
```

This keeps module registration small and makes the real language model testable in Go without a Goja runtime.

### 4.2 Typed refs as the default Goja substrate

Adopt the goja-bleve/geppetto hidden-key pattern as the default. Every builder handle should contain a typed Go ref hidden on the JS object.

Pseudocode:

```go
type RefKind string

const (
    KindPageBuilder       RefKind = "PageBuilder"
    KindSchemaBuilder     RefKind = "SchemaBuilder"
    KindSchema            RefKind = "Schema"
    KindCollectionBuilder RefKind = "CollectionBuilder"
    KindAction            RefKind = "Action"
    KindArrangement       RefKind = "Arrangement"
)

type RefBase struct {
    kind RefKind
    path string
    issues *[]ValidationIssue
}

type SchemaRef struct {
    RefBase
    Spec *spec.SchemaSpec
}

func AttachRef(vm *goja.Runtime, obj *goja.Object, ref any) {
    // Define non-enumerable hidden property, like goja-bleve/geppetto.
}

func GetRef[T any](vm *goja.Runtime, value goja.Value, want RefKind) (*T, error) {
    // Extract hidden ref; verify kind; verify Go type; return rich JS error.
}
```

Use hidden refs for high-level handles: schema, fields, actions, arrangements, marks, fragments, builders. Low-level `WidgetNode` can still be serialized maps because it is already terminal-ish.

### 4.3 Validation model

Use three layers of validation:

1. **Eager shape validation at method boundaries.** Wrong callback type, wrong handle kind, invalid enum literal, unknown option key. These should return Go errors immediately.
2. **Accumulated semantic validation on specs.** Missing schema, arrangement needs a primary field, table column refers to unknown field, action payload references unknown context path. These should accumulate as `ValidationIssue`s.
3. **Terminal behavior.** `.validate()` returns issues. `.toIR()` either returns IR plus embedded diagnostics, or returns an error depending on caller policy. `ui.page()` should optionally render validation issues as an error panel in development.

Suggested issue shape:

```go
type Severity string
const (
    SeverityError Severity = "error"
    SeverityWarning Severity = "warning"
    SeverityInfo Severity = "info"
)

type ValidationIssue struct {
    Severity Severity `json:"severity"`
    Code     string   `json:"code"`
    Path     string   `json:"path"`       // e.g. $.sections[2].collection.arrangement
    Message  string   `json:"message"`
    Hint     string   `json:"hint,omitempty"`
}

type ValidationResult struct {
    Valid  bool              `json:"valid"`
    Issues []ValidationIssue `json:"issues"`
}
```

Development-mode terminal policy:

```js
ui.page({ id: "admin", root: myPage }).toIR({ onValidation: "render" })
```

Could render:

```text
┌───────────────────────────────────────────┐
│ Widget DSL validation failed              │
│ error $.agenda.arrange: unknown arrange…  │
│ hint: did you mean "master-detail"?       │
└───────────────────────────────────────────┘
```

Production-mode terminal policy:

```js
ui.page(...).toIR({ onValidation: "throw" })
```

### 4.4 Widget IR v2: minimal but useful evolution

Do not redesign the entire renderer immediately. Keep these stable:

```ts
type WidgetNode = TextNode | ElementNode | ComponentNode;
interface ComponentNode { kind: "component"; type: string; props?: object; children?: WidgetNode[]; }
```

For v2, add metadata and diagnostics as first-class fields. They can remain optional at the renderer boundary, but the v2 DSL does not need to preserve v1 page-authoring compatibility:

```ts
interface WidgetPageV2 {
  schemaVersion: "0.2.0";
  id: string;
  title?: string;
  meta?: JsonObject;
  root: WidgetNode;
  diagnostics?: ValidationIssue[];
}

interface ComponentNode {
  kind: "component";
  type: RagWidgetType | string;
  props?: WidgetProps;
  children?: WidgetNode[];
  key?: string;
  source?: SourceSpan;      // optional debug provenance
  diagnostics?: ValidationIssue[];
}

interface SourceSpan {
  module?: string;
  helper?: string;
  path?: string;
}
```

The renderer can ignore `source` and `diagnostics` until there is a debug overlay.

### 4.5 Action IR v2: data, not closures

Current `ActionSpec` is useful evidence, but v2 should replace the public action authoring API with typed templates and context requirements. The serialized IR can stay close to the existing union, but authors should not write raw action maps or string interpolation objects as the primary form.

Proposed shape:

```ts
type WidgetAction =
  | NavigateAction
  | DownloadAction
  | ServerAction
  | EventAction
  | CopyAction;

interface ActionBase<C = unknown> {
  confirm?: TemplateSpec<C>;
  context?: ContextContract<C>;
  disabledWhen?: PredicateSpec<C>;
}

interface ServerAction<C = unknown, P = JsonObject> extends ActionBase<C> {
  kind: "server";
  name: string;
  payload?: P | PayloadTemplate<C, P>;
  result?: ServerResultPolicy;
}

interface PayloadTemplate<C, P> {
  kind: "payloadTemplate";
  fields: Record<string, TemplateValue<C>>;
}

interface TemplateValue<C> {
  kind: "path" | "literal" | "format";
  path?: string;       // e.g. "row.id"
  value?: JsonValue;
}
```

The v2 public form is typed. The old raw object/string-template form should move to an explicitly unsafe namespace if it is kept at all:

```js
widget.unsafe.action({
  kind: "server",
  name: "admin-delete-agenda-item",
  confirm: "Delete agenda item “${row.title}”?",
});
```

Blessed v2 form:

```js
ui.action.server("admin-delete-agenda-item", a => a
  .payload(p => p
    .field("id", ctx => ctx.row.id)              // builder-time produces path descriptor
    .field("kind", "agendaItem"))
  .confirm(c => c.text("Delete agenda item “").path("row.title").text("”?"))
  .expects(ctx => ctx.row({ id: "string", title: "string" }))
  .result(r => r.refresh().toastFrom("toast")))
```

Important: `ctx => ctx.row.id` in the example must not be serialized as a JS function. There are two ways to implement it:

- Prefer explicit path builders: `.fieldPath("id", "row.id")`.
- Optionally expose a Proxy-based path recorder so `ctx => ctx.row.id` records the path without running later. That is advanced; keep the explicit form first.

Recommended initial API:

```js
ui.action.server("admin-delete-agenda-item", a => a
  .payload(p => p.path("id", "row.id"))
  .confirm(c => c.text("Delete agenda item “").value("row.title").text("”?")))
```

### 4.6 Server action registry: where real JS callbacks belong

For runtime callbacks, add a host-side action registry. The browser only receives a stable name. The server owns the function.

```js
const actions = require("ui.dsl").actions;

actions.register("admin-delete-agenda-item", h => h
  .context(ctx => ctx.row({ id: "string", title: "string" }))
  .payload(p => p.field("id", "string").optional("kind", "string"))
  .handle(async ({ payload, context, services }) => {
    const id = payload.id || context.row.id;
    const deleted = services.courseMetadata.deleteAgendaItem(id);
    return actions.result().refresh().toast(`Deleted agenda item ${deleted.deleted}`).data({ agenda: deleted.agenda });
  }));
```

The page references that action by name:

```js
ui.action.server("admin-delete-agenda-item", a => a
  .payload(p => p.path("id", "row.id"))
  .confirm(c => c.text("Delete “").value("row.title").text("”?")))
```

Server dispatch pseudocode:

```go
func HandleWidgetAction(req ActionRequest) ActionResult {
    spec, handler := registry.Lookup(req.Name)
    if handler == nil { return NotFound(req.Name) }

    issues := spec.Validate(req.Payload, req.Context)
    if issues.HasErrors() { return BadRequest(issues) }

    // Must run on the Goja runtime owner if handler is a goja function.
    return runtimeOwner.Run(func(vm *goja.Runtime) (ActionResult, error) {
        return handler.Invoke(vm, req.Payload, req.Context)
    })
}
```

This mirrors the current `go-go-course/server.js:422-526` string dispatch, but turns it into a typed registry with explicit contracts.

## 5. Recommended public API model

### 5.1 API principles

1. **One blessed short path, one explicit unsafe hatch.** The common case should be obvious; high-level code should not drop to maps. Raw Widget IR belongs only under `widget.unsafe`.
2. **Lambdas configure builders; they do not cross the browser boundary.** All browser-visible behavior compiles to `ActionSpec` data.
3. **Hyperscript for trees, fluent builders for specs.** `ui.stack(...children)` is good. `data.collection(...).schema(...).edit(...)` is good. Do not force every tree into chained `.child()` calls.
4. **No v1 public aliases.** Do not preserve old option-bag forms, raw marker maps, or duplicated domain helper exports unless model trials prove a specific form is better as a v2 design.
5. **Domain modules export schemas, marks, layouts, and fragments.** Direct one-off recipes should be rewritten as named fragments over typed builders.
6. **Every helper has precise `.d.ts` and parity tests.** Runtime export and TypeScript declarations must stay synchronized.

### 5.2 Minimal authoring style: opinionated defaults

Course landing page:

```js
const ui = require("ui.dsl");
const course = require("course.dsl");

module.exports = ui.page("course", "Course", p => p
  .courseShell(s => s.active("course"))
  .main(
    ui.section("Welcome", ui.stack(
      ui.caption({ tone: "muted" }, "Workshop · GenAI"),
      ui.textBlock({ size: "title" }, "Context Window Engineering"),
      ui.inline(
        ui.button.primary("Open slides", ui.action.navigate("/pages/slides")),
        ui.button.secondary("Handouts", ui.action.navigate("/pages/handouts"))
      )
    ))
  )
).toIR();
```

This shape assumes convenience overloads such as `ui.page(id, title, fn)` and `ui.button.primary(label, action)`. It is pleasant, but it adds more overloads and must be carefully typed.

Equivalent lower-ceremony current-style tree, with strict props:

```js
module.exports = ui.page({ id: "course", title: "Course" },
  ui.section("Welcome",
    ui.stack({ gap: "md" },
      ui.caption({ tone: "muted" }, "Workshop · GenAI"),
      ui.textBlock({ size: "title" }, "Context Window Engineering"),
      ui.inline({ gap: "sm" },
        ui.button({ variant: "primary", action: ui.action.navigate("/pages/slides") }, "Open slides"),
        ui.button({ variant: "secondary", action: ui.action.navigate("/pages/handouts") }, "Handouts")))));
```

Recommendation: support the second form first, then add selected sugar such as `ui.button.primary` only if authors really want it.

### 5.3 Data schema examples

#### Option A: current-like strict schema object

```js
const agendaSchema = data.schema("Agenda", {
  id: data.f.key({ label: "ID", hint: "Stable anchor", width: "18ch" }),
  number: data.f.short({ label: "Time", width: "8ch" }),
  duration: data.f.short({ width: "8ch" }),
  title: data.f.primary({ required: true, maxLength: 160 }),
  description: data.f.prose({ rows: 4, maxLength: 800 }),
});
```

Pros:

- Very close to current code.
- Object literal preserves visual field order if captured via `Object.Keys()`.
- Easy for humans and agents.

Cons:

- Field names are not easy to type-narrow in generated DTS.
- Deep option bags can still become typo-prone unless decoders reject unknown keys.

#### Option B: fluent schema builder

```js
const agendaSchema = data.schema("Agenda", s => s
  .key("id", f => f.label("ID").hint("Stable anchor").width("18ch"))
  .short("number", f => f.label("Time").width("8ch"))
  .short("duration", f => f.width("8ch"))
  .primary("title", f => f.required().maxLength(160))
  .prose("description", f => f.rows(4).maxLength(800))
).build();
```

Pros:

- Great autocomplete.
- Every field method can validate immediately.
- Supports `.use(fragment)` naturally.

Cons:

- More verbose.
- Some people prefer schemas to look like data.

#### Option C: hybrid recommended shape

```js
const agendaSchema = data.schema("Agenda")
  .fields({
    id: data.f.key().label("ID").hint("Stable anchor").width("18ch"),
    number: data.f.short().label("Time").width("8ch"),
    duration: data.f.short().width("8ch"),
    title: data.f.primary().required().maxLength(160),
    description: data.f.prose().rows(4).maxLength(800),
  })
  .build();
```

Pros:

- Still reads like a schema object.
- Field specs are typed handles, not maps.
- Order is captured by the `.fields(...)` boundary.
- Allows both preset helpers and fluent facets.

Recommendation: implement Option C as the durable target. Do **not** keep Option A as a public compatibility facade. If an object-literal schema shape is later desired, add it only if small-model trials prove it is better than the fluent field-handle shape; do not preserve it merely because v1 used it.

### 5.4 Field facet examples

The current `role` string combines storage type, semantic role, summary rendering, editor control, validation, and layout. V2 should split those facets internally while preserving shortcuts.

Shortcut field:

```js
data.f.primary().required().maxLength(160)
```

Expanded field:

```js
data.field("title")
  .type(data.type.string())
  .identity("primary")
  .summary(data.summary.text({ width: "28ch" }))
  .editor(data.editor.text({ maxLength: 160 }))
  .validate(v => v.required())
```

Status field:

```js
data.f.status()
  .choices("draft", "published", "archived")
  .summary(data.summary.status({ icon: true }))
  .editor(data.editor.select())
```

Media field:

```js
data.field("cover")
  .type(data.type.assetRef({ kinds: ["image"] }))
  .summary(data.summary.mediaThumb({ aspect: "wide" }))
  .editor(cms.editor.assetPicker({ accept: "image/*" }))
```

Metric field:

```js
data.field("tokens")
  .type(data.type.number())
  .semantic("measure")
  .summary(data.summary.meter({ limitField: "tokenLimit", tone: "accent" }))
```

### 5.5 Collection examples: simple to rich

#### Simplest table

```js
data.collection("sessions", sessions)
  .schema(sessionSchema)
  .table()
  .toIR();
```

Default behavior:

- choose key field from schema key or `id`,
- elide prose/media fields,
- render status/count/measure with default cells,
- show empty state if no rows,
- no mutation actions.

#### Selectable table via URL

```js
data.collection("sessions", sessions)
  .schema(sessionSchema)
  .select(data.selection.urlParam("selected", query.selected))
  .table(t => t.onRowSelect(data.action.navigateToSelection()))
  .toIR();
```

More explicit:

```js
data.collection("sessions", sessions, c => c
  .schema(sessionSchema)
  .select(s => s.urlParam("selected", query.selected))
  .arrange(a => a.table(t => t
    .rowSelect(ui.action.navigate("/pages/sessions?selected=${row.sessionId}")))))
```

#### Master-detail editor

```js
data.collection("agenda", agenda, c => c
  .schema(agendaSchema)
  .edit(e => e
    .selectUrl("agenda", query.agenda)
    .submitPost("/settings/agenda-item")
    .create({ label: "New agenda item" })
    .actions(a => a
      .reorder(ui.action.server("admin-reorder-course-agenda"))
      .remove(ui.action.server("admin-delete-agenda-item", s => s
        .payload(p => p.path("id", "row.id"))
        .confirm(c => c.text("Delete agenda item “").value("row.title").text("”?"))))))
  .masterDetail())
  .toIR();
```

#### Rich multi-view collection

```js
data.collection("media", assets, c => c
  .schema(cms.schemas.asset())
  .state(s => s
    .query(query.q)
    .page(Number(query.page || 1))
    .selected(query.asset))
  .actions(a => a
    .select(ui.action.navigate("/pages/admin-course-cms?asset=${assetId}"))
    .open(ui.action.navigate("/course-assets/${assetId}"))
    .upload(ui.action.server("admin-upload-course-material", x => x.payload({ kind: "media" }))))
  .views(v => v
    .tiles(cms.marks.assetTiles(t => t.minTileWidth(180)))
    .table(cms.marks.assetTable())
    .detail(cms.marks.assetInspector()))
  .defaultView("tiles"))
  .toIR();
```

This example points toward replacing `cms.recipes.mediaLibrary` with a domain schema + mark set.

### 5.6 UI structure examples

Keep hyperscript nesting for trees.

Simple section:

```js
ui.section("Agenda",
  data.collection("agenda", agenda).schema(agendaSchema).table())
```

Strict props section:

```js
ui.section("Agenda", {
  level: 2,
  anchor: "agenda",
  caption: "Ordered workshop agenda shown on the Course page.",
  density: "flush",
}, agendaEditor)
```

Fluent section option, if wanted:

```js
ui.section("Agenda", s => s
  .level(2)
  .anchor("agenda")
  .caption("Ordered workshop agenda shown on the Course page.")
  .child(agendaEditor))
```

Recommendation: keep the strict props form as primary; add fluent section only if it unlocks validation/error paths that props cannot.

### 5.7 Context-window DSL examples

The context-window DSL should stop being mostly direct component factories and become a vocabulary for context-specific data, marks, and layouts.

Current-ish style:

```js
const styleSet = contextWindow.paletteStyleSet({
  palette: "Signal Orange / Cyan",
  entries: [
    { id: "prompt", label: "Prompt", accent: "b", pattern: "checker" },
    { id: "evidence", label: "Evidence", accent: "a", pattern: "stipple" },
    { id: "answer", label: "Answer", accent: "a", pattern: "solid" },
  ],
});

contextWindow.contextDiagramPanel({ snapshot, styleSet, initialView: "budget" });
```

V2 style-set builder:

```js
const styleSet = contextWindow.styles.palette("Signal Orange / Cyan", p => p
  .entry("prompt", "Prompt", e => e.accent("b").pattern("checker"))
  .entry("evidence", "Evidence", e => e.accent("a").pattern("stipple"))
  .entry("answer", "Answer", e => e.accent("a").solid())
  .entry("free", "Headroom", e => e.accent("grid").hidden()));
```

V2 diagram mark:

```js
contextWindow.snapshot("rag-window", s => s
  .title("RAG answer window")
  .limit(32000)
  .part("prompt", "Prompt", 1400)
  .part("evidence", "Evidence", 9200)
  .part("answer", "Draft", 1800)
  .headroom("free"))
  .render(contextWindow.marks.diagram(d => d.view("budget").styleSet(styleSet)))
```

Domain mark inside general data grammar:

```js
data.record("context", snapshot)
  .schema(contextWindow.schemas.snapshot())
  .show(s => s.mark(contextWindow.marks.budgetDiagram({ styleSet })))
```

### 5.8 Course DSL examples

Course module should provide course-specific schemas, marks, and shell conventions.

Simple course studio page:

```js
course.page("course", p => p
  .active("course")
  .section("main", "Course", s => s
    .lesson(course.lesson(definition, l => l
      .primaryCta("Open slides", ui.action.navigate("/pages/slides"))
      .secondaryCta("Handouts", ui.action.navigate("/pages/handouts"))))))
```

Keep component escape hatch:

```js
course.courseStudioShell({ sections, activeItemId: "slides", title: "Studio" },
  course.courseSlidePanel({ slide, snapshot, index: 0, total: 1 }))
```

Schema + collection approach for agenda:

```js
const agenda = data.collection("agenda", content.agenda)
  .schema(course.schemas.agendaItem())
  .edit(course.fragments.agendaEditor({ query, submit: "/settings/agenda-item" }))
  .toIR();
```

Course fragments:

```js
course.fragments.agendaEditor = ({ query, submit }) => c => c
  .selectUrl("agenda", query.agenda)
  .submitPost(submit)
  .masterDetail()
  .actions(a => a
    .reorder(ui.action.server("admin-reorder-course-agenda"))
    .remove(ui.action.server("admin-delete-agenda-item", x => x
      .payload(p => p.path("id", "row.id"))
      .confirm(c => c.text("Delete “").value("row.title").text("”?")))))
```

### 5.9 CMS DSL examples

The CMS DSL should export typed schemas and marks before it exports large recipes.

Asset schema:

```js
const assetSchema = cms.schemas.asset(a => a
  .id("id")
  .file("filename")
  .kind("kind")
  .url("src")
  .size("size")
  .updatedAt("updatedAt"));
```

Simple media library:

```js
cms.mediaLibrary(material.mediaAssets)
  .selected(query.asset)
  .upload(ui.action.server("admin-upload-course-material", a => a.payload({ kind: "media" })))
  .toIR();
```

Composable media library through `data.collection`:

```js
data.collection("media", material.mediaAssets, c => c
  .schema(cms.schemas.asset())
  .selectUrl("asset", query.asset)
  .arrange(a => a.tiles(cms.marks.assetTileGrid()))
  .detail(cms.marks.assetDetailPanel(d => d
    .open(ui.action.navigate("/course-assets/${assetId}"))
    .download(ui.action.download("/course-assets/${assetId}"))
    .delete(ui.action.server("admin-delete-course-material", x => x
      .payload({ kind: "media" })
      .payloadPath("file", "asset.filename")
      .confirm(c => c.text("Delete ").value("asset.filename").text("?"))))))
  .toIR();
```

Article management:

```js
data.collection("articles", articles, c => c
  .schema(cms.schemas.articleSummary())
  .state(s => s.query(query.q).filter("status", query.status || "all"))
  .arrange(a => a.table(cms.marks.articleTable(t => t
    .rowActions(actions => actions
      .edit(ui.action.navigate("/pages/admin-article?article=${row.id}"))
      .publish(ui.action.server("admin-publish-article"))
      .archive(ui.action.server("admin-archive-article"))))))
  .toolbar(t => t
    .search("q")
    .statusFilter("status")
    .create(ui.action.navigate("/pages/admin-article?new=1")))
  .toIR();
```

### 5.10 General webpage examples

Static-ish documentation page:

```js
ui.page("about", "About", p => p
  .root(ui.stack({ gap: "lg" },
    ui.section("What this is",
      ui.markdownArticle({ source: "# RAG Evaluation\n\nA toolkit for inspecting retrieval." })),
    ui.section("Next steps",
      ui.checkList({ items: [
        { id: "upload", label: "Upload a session", checked: true },
        { id: "inspect", label: "Inspect context windows", checked: false },
      ] })))))
```

Dashboard page:

```js
ui.page("dashboard", "Dashboard", p => p
  .root(ui.stack({ gap: "lg" },
    ui.recipes.metrics({ items: [
      { label: "Sessions", value: stats.sessions, status: "ready" },
      { label: "Uploads", value: stats.uploads, status: "ready" },
    ]}),
    data.collection("recentSessions", sessions)
      .schema(app.schemas.session())
      .table(t => t.compact().open(ui.action.navigate("/pages/sessions/${row.sessionId}"))))))
```

## 6. Brainstormed API choices

This section remains useful for API taste discussions, but under the hard-cutover assumption only Choices 2–5 are plausible v2 directions. Choice 1 is retained as an explicit rejected baseline: it shows the kind of v1-shaped option-bag API we should not preserve merely for compatibility.

### Choice 1: strict props + terminal only (rejected baseline)

```js
const page = ui.page({
  id: "admin",
  title: "Admin",
  sections: [
    ui.section("Agenda", { level: 2 },
      data.collection(agenda, {
        schema: agendaSchema,
        verb: "edit",
        arrange: "master-detail",
      }))
  ]
});
```

This is the least disruptive option, but that is no longer the goal. It remains useful only as a baseline for small-model trials or as a description of v1. Do not implement this as the v2 public API.

Why rejected for hard cutover:

- preserves the v1 option-bag shape;
- keeps `verb` and `arrange` as strings instead of methods;
- makes reusable fragments awkward;
- keeps complex arrangements as nested JSON;
- encourages the exact raw-map authoring style v2 is meant to remove.

### Choice 2: fluent all the way

```js
const page = ui.page("admin")
  .title("Admin")
  .section("Agenda", s => s
    .level(2)
    .child(data.collection("agenda", agenda)
      .schema(agendaSchema)
      .edit()
      .masterDetail()))
  .toIR();
```

Use when:

- you want maximal autocomplete and explicit terminals,
- the DSL is mostly specs rather than tree structure.

Weakness:

- trees can become verbose,
- harder to skim than natural nesting,
- adds many builder types.

### Choice 3: hybrid recommended model

```js
const page = ui.page({ id: "admin", title: "Admin" },
  ui.section("Agenda", { level: 2 },
    data.collection("agenda", agenda, c => c
      .schema(agendaSchema)
      .edit(course.fragments.agendaEditor({ query }))
      .masterDetail())));
```

Use when:

- structure is naturally nested,
- complex specs use builders,
- simple props stay simple but validated.

This is the recommended model.

### Choice 4: recipe-first domain APIs

```js
course.adminAgendaEditor({
  agenda,
  query,
  submit: "/settings/agenda-item",
  reorderAction: "admin-reorder-course-agenda",
  deleteAction: "admin-delete-agenda-item",
});
```

Use when:

- there is one blessed product pattern,
- authoring speed matters more than composability,
- domain defaults are strong.

Weakness:

- can hide too much,
- risks becoming another pile of one-off recipes,
- harder to compose across domains.

Recommended compromise:

```js
course.recipes.adminAgendaEditor(options)
// internally implemented as:
data.collection("agenda", options.agenda).schema(course.schemas.agendaItem()).use(...)
```

Recipes should be wrappers over typed grammar terms, not separate compiler paths.

### Choice 5: declarative operator objects

```js
data.collection({
  id: "agenda",
  rows: agenda,
  schema: agendaSchema,
  operators: [
    data.op.edit(),
    data.op.selectUrl("agenda", query.agenda),
    data.op.masterDetail(),
    data.op.serverActions({ reorder: "admin-reorder-course-agenda" }),
  ],
});
```

This borrows from go-emrichen's tag-operator model. It is powerful and serializable, but less idiomatic for JavaScript authors than fluent builders.

Recommendation: keep this as an internal representation. Public API should be builder/configurator methods that compile to operator specs.

## 7. Renderer and runtime action design in detail

### 7.1 The invariant

A Widget IR page is JSON. Therefore:

- It may contain strings, numbers, booleans, null, arrays, and objects.
- It may contain action descriptors.
- It may contain payload/context templates.
- It must not contain `goja.Value`, Go pointers, JavaScript functions, DOM nodes, Promises, or closures.

Builder lambdas are consumed before serialization:

```text
JS lambda c => c.schema(...)
  runs inside Goja while building page
  mutates typed Go spec
  disappears before HTTP JSON
```

Action handlers are registered outside the IR:

```text
JS handler async ({payload, context}) => ...
  stored in host action registry
  invoked later by /api/widget/actions/:name
  referenced in IR only by stable name
```

### 7.2 Action context contracts

Every adapter that dispatches an action should declare the context it sends. Today this is implicit in adapter code:

- `DataTable` row selection sends `{ row, rowKey, componentType: "DataTable" }` (`DataTable.widget.tsx:22-28`).
- `DataTable` action cells send `{ row, rowKey, componentType: "DataTableCell" }` (`cellRenderers.tsx:55-69`).
- `MediaLibraryPanel` sends `{ assetId }`, `{ query }`, `{ kind }`, `{ page }`, or file upload context (`MediaLibraryPanel.widget.tsx:19-83`).
- `ArticleListPanel` sends `{ articleId }`, `{ rowAction }`, `{ status }`, `{ query }`, `{ page }` (`ArticleListPanel.widget.tsx:16-55`).
- `SearchField` sends `{ query, value, componentType: "SearchField" }` (`SearchField.widget.tsx:28-38`).
- `Button` sends `{ componentType: "Button" }` (`Button.widget.tsx:8-17`).

V2 should make these contracts data:

```ts
interface WidgetActionContextMap {
  Button: { componentType: "Button" };
  DataTable: { componentType: "DataTable"; row: JsonObject; rowKey: string };
  DataTableCell: { componentType: "DataTableCell"; row: JsonObject; rowKey: string };
  MediaLibraryPanelAsset: { componentType: "MediaLibraryPanel"; assetId: string; value: string };
  MediaLibraryPanelUpload: { componentType: "MediaLibraryPanel"; files: SerializedUploadFile[]; fileNames: string[]; fileCount: number };
  ArticleListPanelRowAction: { componentType: "ArticleListPanel"; articleId: string; rowAction: "edit" | "publish" | "archive" | "delete"; value: string };
}
```

At minimum, encode the context contract in manifests:

```yaml
actions:
  - name: onRowSelect
    context: DataTableRowContext
  - name: onFilesSelectedAction
    context: UploadFilesContext
```

This lets the DSL declarations say:

```ts
type ServerAction<C> = { kind: "server"; name: string; payload?: JsonObject; confirm?: TemplateSpec<C> };
interface DataTableBuilder<T> {
  onRowSelect(action: ServerAction<RowContext<T>> | NavigateAction<RowContext<T>>): this;
}
```

### 7.3 Payload templates instead of function serialization

Bad idea:

```js
ui.action.server("delete", ctx => ({ id: ctx.row.id })) // cannot serialize
```

Good explicit form:

```js
ui.action.server("delete", a => a.payload(p => p.path("id", "row.id")))
```

Good structured confirm form:

```js
ui.action.server("delete", a => a
  .confirm(c => c.text("Delete “").path("row.title").text("”?")))
```

Compiled IR:

```json
{
  "kind": "server",
  "name": "delete",
  "payloadTemplate": {
    "id": { "kind": "path", "path": "row.id" }
  },
  "confirm": {
    "kind": "template",
    "parts": [
      { "kind": "text", "text": "Delete “" },
      { "kind": "path", "path": "row.title", "encode": false },
      { "kind": "text", "text": "”?" }
    ]
  }
}
```

Renderer dispatch algorithm:

```ts
function dispatchWidgetAction(action, context) {
  const hydrated = hydrateAction(action, context);
  if (hydrated.confirm && !window.confirm(hydrated.confirm)) return;
  switch (hydrated.kind) {
    case "server": post(hydrated.name, { payload: hydrated.payload, context }); break;
    case "navigate": navigate(hydrated.to); break;
    // ...
  }
}
```

This resolves the current confirm/URL encoding split (`actions.ts:107-134`) by making encoding a property of template parts.

### 7.4 Server result policy

The current `ServerActionResult` shape is:

```ts
interface ServerActionResult {
  ok: boolean;
  refresh?: boolean;
  toast?: string;
  patch?: JsonObject;
  data?: JsonObject;
}
```

Evidence: `packages/rag-evaluation-site/src/widgets/actions.ts:9-15`.

Keep this, but add policy in the action descriptor so authors can be explicit:

```js
ui.action.server("save", a => a
  .payloadForm()
  .result(r => r.refresh().toast()))
```

Action result handling can remain centralized in `RagEvaluationSiteApp.handleAction`.

### 7.5 Runtime ownership rule

If handlers are Goja functions, they must run on the runtime owner. This is not optional. The catalogue identified goja-dbus and geppetto as evidence that callbacks/Promises must settle on the owning runtime. For this system:

- Page building already occurs in the host runtime.
- Action handlers registered as JS functions must be invoked through the same runtime owner/queue.
- Do not invoke a captured `goja.Callable` from arbitrary HTTP goroutines.
- If the host is JavaScript `express` in xgoja, preserve the xgoja event-loop semantics rather than inventing Go goroutine dispatch.

For intern implementation, first support explicit string-named server actions using the host's existing route. Add function registration only after runtime-owner tests exist.

## 8. Hard-cutover implementation plan

The old plan was incremental: strict-decode v1, add typed refs under compatibility facades, then introduce builders beside existing helpers. That is no longer the right optimization target. Since backwards compatibility is not required, v2 should be implemented as a clean replacement and known consumers should be rewritten.

### Phase 0: preserve v1 only as evidence, not as API

Goal: capture enough behavior to avoid losing product insight, then stop treating v1 as a contract.

Tasks:

1. Snapshot the current `go-go-course` admin CMS page and relevant example pages as **behavioral fixtures**, not compatibility tests.
2. Record the useful concepts discovered by v1: `section`, schema/field intent, record editing, collection editing, URL-backed selection, native form posts, row actions, upload actions.
3. Add negative tests that assert the known v1 failure modes are rejected by v2: typo'd arrangement, wrong marker kind, invalid section level, unknown field option, unknown action context path.
4. Delete or quarantine tests whose only purpose is to keep v1 map shapes working.

Validation commands:

```bash
go test ./pkg/widgetdsl -count=1
pnpm --dir packages/rag-evaluation-site build
```

### Phase 1: define the v2 spec model first

Goal: design the typed core before exposing JavaScript methods.

Implement typed packages for:

```go
type PageSpec struct { ... }
type SectionSpec struct { ... }
type SchemaSpec struct { Fields []FieldSpec }
type FieldSpec struct { Name string; Facets FieldFacets }
type CollectionSpec struct { Name string; Rows []JsonObject; Schema *SchemaSpec; Mode CollectionMode; Arrangement ArrangementSpec }
type ActionSpec struct { Kind ActionKind; Name string; Payload PayloadTemplate; Confirm TemplateSpec }
type ValidationIssue struct { Severity, Code, Path, Message, Hint string }
```

Hard-cutover rules:

- Do not model `LegacyRole` except in one-off migration scripts.
- Do not model v1 `map[string]any` marker compatibility.
- Do not make `CollectionOptions` mirror the v1 options object.
- Do model the concepts v1 proved useful: selection, submit binding, create/reorder/remove, table/master-detail, field facets.

### Phase 2: implement the v2 builders only

Goal: expose one blessed authoring shape, not v1 plus v2.

Minimum public API:

```js
const agendaSchema = data.schema("Agenda")
  .field("id", data.f.key().label("ID").readOnly())
  .field("number", data.f.short().label("Time").width("8ch"))
  .field("title", data.f.primary().required().maxLength(160))
  .field("description", data.f.prose().rows(4))
  .build();

const agendaEditor = data.collection("agenda", agenda)
  .schema(agendaSchema)
  .edit(e => e
    .selectUrl("agenda", query.agenda)
    .submitPost("/settings/agenda-item")
    .create("New agenda item")
    .reorder("admin-reorder-course-agenda")
    .remove("admin-delete-agenda-item", a => a
      .confirm(c => c.text("Delete “").path("row.title").text("”?"))))
  .masterDetail()
  .toIR();
```

Builder implementation requirements:

- Hidden typed refs for schemas, fields, actions, arrangements, marks, and builders.
- Strict `goja.AssertFunction` for any present callback.
- Methods return `(value, error)` where misuse can occur; no panics for author errors.
- `.validate()` and `.toIR()` terminals on complex builders.
- No public v1 `data.schema({ ... })`, `data.collection(rows, { ... })`, `data.urlParam(...)` marker maps, or `data.formPost(...)` marker maps.

### Phase 3: implement Action IR v2 and renderer hydration

Goal: make server/browser behavior serializable from the start.

Implement:

- `TemplateSpec` parts: literal text, context path, formatted value.
- `PayloadTemplate` fields: literal value, context path, object composition.
- Renderer hydration before dispatch.
- Typed action contexts for key adapters: Button, DataTable row, DataTable cell, UploadDropArea, MediaLibraryPanel, ArticleListPanel, SearchField.
- Server-action result policy: refresh, toast, patch/data.

Blessed authoring form:

```js
ui.action.server("admin-delete-agenda-item")
  .payload(p => p.path("id", "row.id"))
  .confirm(c => c.text("Delete “").path("row.title").text("”?"))
```

Explicit escape hatch:

```js
widget.unsafe.action(rawActionObject)
```

Tests should fail if raw action maps are accepted by high-level v2 APIs.

### Phase 4: implement `widget.unsafe` as the only raw escape hatch

Goal: keep emergency access to low-level Widget IR without polluting the blessed DSL.

Expose raw constructors only under a deliberately named namespace:

```js
const unsafe = require("widget.unsafe");
unsafe.component("DataTable", rawProps, ...children)
unsafe.node(rawJson)
unsafe.action(rawActionJson)
```

Rules:

- `ui.dsl`, `data.dsl`, `course.dsl`, `cms.dsl`, and `context_window.dsl` should not expose generic `component(type, props)` as a normal helper.
- Lint production page code that imports `widget.unsafe`.
- Document `widget.unsafe` as a temporary extension point and a signal that the DSL is missing a concept.

### Phase 5: port one real page by rewriting it

Goal: prove the cutover with a real consumer, not with adapters.

Rewrite `go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js` against v2. Do not adapt old calls. The target page should contain no raw collection option bags and no direct giant component maps for agenda/media/article surfaces.

Acceptance criteria:

- Agenda editor uses `data.collection(...).schema(...).edit(...).masterDetail().toIR()`.
- Delete/reorder actions use Action IR v2 payload/confirm templates.
- Media library is either a v2 CMS fragment or explicit `data.collection + cms.schemas.asset + cms.marks.assetTiles`.
- The page visually matches or improves the v1 grammar page.
- Server action POSTs still deliver the payload/context expected by `go-go-course` handlers.

### Phase 6: convert domain recipes into fragments, marks, and schemas

Goal: remove independent recipe compilers.

Implement first:

1. `cms.schemas.asset()` and `cms.marks.assetTiles()`.
2. `cms.fragments.mediaLibrary(options)` as a wrapper over `data.collection`.
3. `course.schemas.agendaItem()` and `course.fragments.agendaEditor(options)`.
4. `contextWindow.snapshot(...)`, `contextWindow.styles.palette(...)`, and `contextWindow.marks.diagram(...)`.

Recipe rule:

```text
A recipe is a named fragment over typed builders.
A recipe is not a separate path that emits a large untyped component tree.
```

### Phase 7: generate precise TypeScript declarations and enforce parity

Goal: make the v2 API discoverable and machine-checkable.

Implement:

- `FragmentFn<T>` and concrete builder interfaces.
- Named field/action/collection/context types.
- Runtime export parity tests like Geppetto.
- TypeScript fixtures with positive examples and `@ts-expect-error` negative examples.

Negative fixtures should include:

```ts
// @ts-expect-error v1 options-bag collection is gone
data.collection(rows, { schema, arrange: "master-detail" });

// @ts-expect-error raw action maps are not high-level actions
ui.button({ action: { kind: "server", name: "delete" } }, "Delete");

// @ts-expect-error invalid section level
ui.section("Agenda").level(7);
```

### Phase 8: delete old public exports and add cutover lint rules

Goal: make the hard cutover enforceable.

Tests/lints should fail if these are public in v2 modules:

- `data.collection(rows, options)` v1 shape,
- v1 `data.record(values, options)` shape if not redesigned as a builder,
- `data.f.primary({ ... })` option-bag field shape,
- `data.urlParam` and `data.formPost` marker-map helpers,
- generic `component` outside `widget.unsafe`,
- raw action maps outside `widget.unsafe`,
- `Props = Record<string, any>` in generated high-level declarations.

### Phase 9: add diagnostics overlay and author feedback

Goal: make v2 pleasant for humans and agents.

- Add optional `diagnostics` to `WidgetPage` and `ComponentNode`.
- Add `WidgetDiagnosticsPanel` or reuse `ErrorCallout`.
- In development mode, render validation issues at the top of the page.
- In production mode, throw on validation errors before serving the page unless explicitly configured otherwise.

## 9. Testing strategy

### 9.1 Go unit tests

- Builder method validation.
- Typed ref extraction and wrong-handle errors.
- Hard-cutover rejection of removed v1 option-bag APIs and raw marker maps.
- Collection compiler output.
- Recipe wrappers compile through grammar terms.
- Server action registry validation.

### 9.2 Goja integration tests

Run real JS snippets:

```go
value, err := vm.RunString(`
  const data = require("data.dsl");
  const ui = require("ui.dsl");
  const schema = data.schema("Agenda").fields({ title: data.f.primary().required() }).build();
  data.collection("agenda", [{ title: "Intro" }]).schema(schema).table().toIR();
`)
```

Assert JSON output and validation result.

### 9.3 TypeScript declaration tests

- Export parity: runtime keys vs generated declaration names.
- `tsc --noEmit` positive examples.
- `@ts-expect-error` negative examples.
- Generated declaration golden file for human review.

### 9.4 React renderer tests

- Renderer ignores v2 metadata when not needed.
- Action templates hydrate payload from row context.
- Confirm template does not URL-encode human-facing text.
- Navigate/download templates do encode URL values.
- Server action POST includes hydrated payload and original context.

### 9.5 End-to-end smoke tests

Use `go-go-course` and the example widget site:

- Fetch `/api/widget/pages/admin-course-cms`.
- Click agenda row: URL selection changes.
- Submit form post: redirect/status works.
- Reorder/delete: server action receives `{payload.direction, context.row}`.
- Upload: serialized files arrive in `context.files`.
- Media select/open/delete actions work.

## 10. Decision records

### Decision 1: hybrid API model

- **Context:** UI pages are trees; data schemas/collections/actions are specs.
- **Options:** strict props only; fluent everything; hybrid hyperscript + fluent specs.
- **Decision:** Use hybrid hyperscript for structure and fluent builders for high-level specs.
- **Rationale:** This preserves readable tree authoring while giving complex objects type safety and fragments.
- **Status:** proposed.

### Decision 2: typed Go-side intent specs before Widget IR

- **Context:** Current helpers emit maps directly; researchctl/codesign prove typed specs work.
- **Options:** keep maps; add validators around maps; introduce typed specs and compile to IR.
- **Decision:** Introduce typed specs and compile to IR at terminals.
- **Rationale:** Enables validation, precise DTS, wrong-handle errors, and renderer-independent design.
- **Status:** proposed.

### Decision 3: action lambdas do not serialize

- **Context:** The user wants lambdas for callbacks/composition, but renderer and server run in different environments.
- **Options:** forbid lambdas; serialize functions somehow; use lambdas only for builders and stable registries for runtime handlers.
- **Decision:** Builder lambdas configure specs; runtime callbacks are registered under stable names; IR carries action data.
- **Rationale:** JSON pages remain portable and safe; server behavior remains testable.
- **Status:** proposed.

### Decision 4: server actions use typed context and payload templates

- **Context:** Current action context is implicit and string-template based.
- **Options:** keep interpolation only; add JS callback payload functions; add serializable template specs.
- **Decision:** Add serializable payload/context template specs and manifest-declared context contracts.
- **Rationale:** Solves row/file/asset actions without serializing closures; types can guide authors.
- **Status:** proposed.

### Decision 5: recipes must compile through grammar terms

- **Context:** Current domain recipes directly emit component maps.
- **Options:** keep recipes as independent compilers; delete recipes; make recipes wrappers over schemas/marks/fragments.
- **Decision:** Recipes become wrappers over grammar terms.
- **Rationale:** Preserves opinionated simplicity without multiplying untyped pathways.
- **Status:** proposed.

### Decision 6: no compatibility shims

- **Context:** Backwards compatibility is not required, and preserving v1 map/option APIs would make the final v2 language larger and less coherent.
- **Options:** compatibility facades; compatibility forever; hard cutover with page rewrites and migration scripts only.
- **Decision:** Hard cutover. Do not implement runtime compatibility shims for v1 public APIs.
- **Rationale:** The old DSL is vocabulary evidence, not a contract. Simplicity and strong typing in v2 are more valuable than preserving old call shapes.
- **Consequences:** Known consumers must be rewritten promptly. Migration help should be scripts/docs, not permanent runtime aliases.
- **Status:** proposed.

## 11. Risks and mitigations

### Risk: overbuilding a framework

Mitigation: implement in vertical slices. The first slice is agenda collection v2. Do not build every mark/layout before one real page proves the model.

### Risk: fluent API becomes too verbose

Mitigation: choose one compact blessed builder style, add domain fragments for common patterns, and use small-model trials to justify any sugar. Do not keep v1 option bags as an anti-verbosity escape hatch.

### Risk: action templates become a second programming language

Mitigation: keep templates small: literal, path, format, object. For real logic, use named server handlers.

### Risk: TypeScript generation becomes a separate source of truth

Mitigation: generate declarations from module specs, manifests, and typed builder registrations; add runtime parity tests.

### Risk: hard cutover breaks current consumers

Mitigation: enumerate known consumers, rewrite them in the same implementation branch, and keep migration scripts/docs for humans. Do not make the renderer contract more complicated just to keep v1 authoring code alive. The serialized Widget IR can remain structurally familiar, but the public authoring API should cut over.

### Risk: runtime callback dispatch violates Goja ownership

Mitigation: do not implement anonymous callback action handlers first. Start with explicit string names and existing host dispatch. Add function registry only with runtime-owner tests.

## 12. Concrete next implementation ticket proposal

Create a new ticket, e.g. `RAGEVAL-DSL-V2-CUTOVER`, with this scope:

1. Define the v2 typed spec model: page, section, schema, field, collection, action, template, mark, arrangement, validation issue.
2. Build `pkg/widgetdsl/fluent` hidden-ref + validation substrate.
3. Implement the v2 builder API only: schema, fields, collection, record/form, section/page, action/template, fragments/marks.
4. Implement Action IR v2 hydration in the React renderer and typed action contexts for the main adapters.
5. Add `widget.unsafe` as the only raw Widget IR escape hatch.
6. Generate precise TypeScript declarations and add runtime export parity plus `tsc` fixtures.
7. Rewrite the `go-go-course` admin CMS page against v2 as the proof page.
8. Re-express `cms` media library and `course` agenda helpers as schemas/marks/fragments.
9. Remove old public v1 exports from the selected v2 modules.

Acceptance criteria:

- `data.collection(rows, { ... })` is not a valid public v2 API.
- `data.schema({ ... })` marker-map style is not a valid public v2 API.
- High-level declarations do not use `Props = Record<string, any>`.
- Raw action maps are rejected outside `widget.unsafe`.
- Typo'd arrangement/mode/section-level/action-context paths fail with precise validation errors.
- Agenda page can be authored in the new API with no raw component maps in the collection body.
- Server actions still POST data-only JSON and receive typed/hydrated payloads.
- Generated `.d.ts` catches invalid arrangements, wrong callback values, wrong handle types, and removed v1 APIs.

## 13. Reference map

### Current rag-evaluation-system files

- `pkg/widgetdsl/module.go` — module names, helper maps, generic component factories, recipes, action/cell helpers.
- `pkg/widgetdsl/grammar.go` — current data/ui grammar verbs: field roles, schema, record, collection, urlParam, formPost, section.
- `pkg/widgetdsl/typescript.go` — current weak declaration generation.
- `pkg/widgetdsl/grammar_test.go` — current grammar expansion tests.
- `pkg/widgetdsl/module_test.go` — split module and recipe tests.
- `internal/widgetmanifest/types.go` — widget manifest model naming module/helper/props/actions.
- `packages/rag-evaluation-site/src/widgets/ir.ts` — React-side Widget IR, prop types, action union.
- `packages/rag-evaluation-site/src/widgets/actions.ts` — current central action dispatch and interpolation.
- `packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx` — renderer dispatch through registry/adapters.
- `packages/rag-evaluation-site/src/widgets/defaultRegistry.ts` — module registry grouping.
- `packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.widget.tsx` — row action context.
- `packages/rag-evaluation-site/src/widgets/cellRenderers.tsx` — cell action context and row template rendering.
- `packages/rag-evaluation-site/src/components/organisms/MediaLibraryPanel/MediaLibraryPanel.widget.tsx` — asset/query/page/upload action contexts.
- `packages/rag-evaluation-site/src/components/organisms/ArticleListPanel/ArticleListPanel.widget.tsx` — article action contexts.
- `packages/rag-evaluation-site/src/app/App.tsx` — app-level server action POST and refresh behavior.

### Current consumer files

- `go-go-course/cmd/go-go-course/server.js` — Widget page endpoint and server action dispatch.
- `go-go-course/cmd/go-go-course/lib/course-pages.js` — page builder routing and dependency wiring.
- `go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js` — real admin/CMS page using current grammar.
- `go-go-course/cmd/go-go-course/lib/pages/admin-common.js` — DataTable/cell/action patterns for material tables.

### Comparator patterns

- `goja-bleve/pkg/api_types.go` — hidden typed refs and `getTypedRef[T]` substrate.
- `geppetto/pkg/js/modules/geppetto/dts_parity_test.go` — runtime export vs generated `.d.ts` parity.
- `researchctl/pkg/gojamodules/codesign/builders.go` — lambda configurators, `.use(fragment)`, strict callback assertion.
- `researchctl/pkg/gojamodules/codesign/typescript.go` — precise builder interfaces and `FragmentFn<T>`.
- `researchctl/pkg/research/spec/types.go` and `pkg/codesign/spec/types.go` — typed Go-side specs.
- `go-emrichen/pkg/emrichen/parser.go` and `emrichen.go` — strict operator argument parsing and recursive operator composition.

### Companion GOJA-DSL-PLAYBOOK documents

- `design-doc/01-goja-dsl-catalogue-and-base-research.md`
- `design-doc/02-self-assessment-of-the-widgetdsl-grammar-what-pattern-c-actually-costs-and-what-the-playbook-should-add.md`
- `design-doc/03-widget-dsl-design-assessment-and-improvement-report.md`
- `design-doc/04-goja-dsl-deep-dive-optional-lambdas-typed-ir-dts-parity-and-tag-operators.md`
