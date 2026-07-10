---
Title: 'Full Widget DSL Redesign: Typed Builders, Slots, Fragments, and Domain Views'
Ticket: RAGEVAL-SCHEDULE-WIDGETS
Status: active
Topics:
    - ui-dsl
    - widget-ir
    - design-system
    - react
    - frontend-architecture
    - intern-guide
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/builders.go
      Note: .use(fragment) and builder callback pattern for reusable policies
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/callbacks.go
      Note: Trusted callback registry boundary that informs no-lambdas-in-browser rule
    - Path: abs:///home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/researchctl/builders.go
      Note: Scoped builder callback implementation pattern
    - Path: repo://pkg/widgetdsl/module.go
      Note: Current split-module runtime, component factories, actions/cells, and recipes being replaced by widget.dsl
    - Path: repo://pkg/widgetdsl/typescript.go
      Note: Current generated declaration surface motivating descriptor-driven TypeScript output
    - Path: repo://pkg/widgetdsl/v2/spec/types.go
      Note: Go backing spec model that informs the proposed PageSpec/ViewSpec/CollectionSpec design
    - Path: repo://pkg/widgetdsl/v2_builders.go
      Note: Existing typed builder mechanics retained under the new data.collection surface
    - Path: ws://go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js
      Note: In-the-wild Course CMS page used for complex TypeScript redesign examples
    - Path: ws://go-go-course/cmd/go-go-course/lib/pages/common.js
      Note: Current course shell wrapper around course.dsl recipe used as redesign target
    - Path: ws://go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js
      Note: In-the-wild mixed DSL usage and data.v2 examples used as redesign input
    - Path: ws://go-go-course/cmd/go-go-course/lib/pages/handouts.js
      Note: Current handout page using course.dsl shell components used as redesign target
ExternalSources: []
Summary: 'A clean-break redesign for the Goja Widget DSL across ui, data, cms, course, context, schedule, and time: one typed TypeScript-facing grammar backed by Go specs, builder lambdas, slots, fragments, bindings, and domain intents.'
LastUpdated: 2026-07-07T16:05:00-04:00
WhatFor: Use this as the target design when replacing the existing Widget DSL surface rather than incrementally extending ui.dsl/data.dsl/cms.dsl/course.dsl/context_window.dsl.
WhenToUse: Read before implementing a breaking Widget DSL v3 or writing new go-go-course-style pages against the redesigned API.
---


# Full Widget DSL Redesign: Typed Builders, Slots, Fragments, and Domain Views

> **Premise.** We do not need backward compatibility. That changes the design.
> We can stop treating the current modules as surfaces to preserve and instead
> design the DSL that should exist after learning from `go-go-course`,
> `researchctl`, the current `data.v2.dsl` experiment, CMS widgets, scheduling
> widgets, and the Widget IR renderer.
>
> **Result.** This document proposes a clean-break Widget DSL v3: a TypeScript-facing
> authoring API for Goja scripts, backed by typed Go specs. It uses builder lambdas
> for configuration, named slots for composition, fragments for reusable policy,
> bindings for serializable data access, and domain intents for actions. It replaces
> the old split between `ui.dsl`, `data.dsl`, `data.v2.dsl`, `cms.dsl`,
> `course.dsl`, and `context_window.dsl` with one coherent grammar.

---

## 1. Why a clean break is better than another extension

The existing DSL grew one layer at a time. The first layer was simple component
factories: `ui.panel(props, children)` lowers to a `Panel` Widget IR node. Then came
cell helpers, action helpers, context style helpers, hand-written recipes,
`data.v2.dsl` typed builders, and domain modules for course and CMS components.
Every layer solved a real problem, but the result is not one language. It is a set
of related dialects.

The `go-go-course` code shows both the power and the cost. Pages are real product
pages, not toy examples. They use `ui.dsl` for sections and forms, `course.dsl` for
course shells and handouts, `cms.dsl` for media libraries, `context_window.dsl` for
diagrams, and `data.v2.dsl` for editable collections. The output works. The authoring
experience is uneven:

- A course shell is a recipe option bag: `courseDsl.recipes.courseStudio({ ... })`.
- A handout page directly calls `courseDsl.handoutDocumentShell({ ... })`.
- A CMS media library uses `cmsDsl.recipes.mediaLibrary({ ... })`.
- A data editor uses a fluent builder: `dataV2.collection(...).schema(...).edit(...).masterDetail().toIR()`.
- Page structure uses `ui.page({ sections: [...] })`, but pages often use custom
  helper functions around it.
- Actions are transport-shaped: `ui.action.navigate(...)`, `ui.action.server(...)`,
  with context paths embedded as strings.

A new author has to learn where each style applies. A new implementer has to keep
TypeScript declarations, module helper maps, recipes, registry manifests, and action
contracts in sync. The clean-break redesign should remove that accidental complexity.

The replacement should not be a larger component factory. It should be a small
language for product screens.

---

## 2. The core contract: TypeScript syntax, Go-backed specs, JSON output

The DSL runs in Goja, but authors should write TypeScript-shaped code. The runtime
may execute transpiled JavaScript or plain JavaScript, but the public contract is a
`.d.ts` file generated from Go descriptors. That `.d.ts` is not decoration. It is the
primary teaching surface.

The backing model has three stages:

```text
TypeScript-facing DSL call
  -> Go builder/spec handles inside Goja
  -> validated Go spec tree
  -> Widget IR JSON page
  -> React WidgetRenderer
```

The browser never receives JavaScript functions. Lambdas are author-time tools. The
browser receives data.

The redesigned Go backing types should be explicit. The names below are illustrative,
but the split is important:

```go
type PageSpec struct {
    ID       string
    Title    string
    Meta     map[string]any
    Shell    *ShellSpec
    Sections []SectionSpec
}

type SectionSpec struct {
    ID       string
    Title    string
    Caption  string
    Actions  []ActionSpec
    Children []NodeSpec
}

type NodeSpec struct {
    Kind     NodeKind
    Type     string
    Props    map[string]any
    Text     string
    Children []NodeSpec
    Source   SourceSpan
}

type ViewSpec struct {
    Domain   string
    Name     string
    Data     any
    Options  map[string]any
    Slots    map[string]SlotSpec
    Actions  map[string]ActionSpec
}

type ActionSpec struct {
    Kind    ActionKind
    Name    string
    To      string
    Event   string
    Payload map[string]BindingSpec
    Confirm string
}

type BindingSpec struct {
    Kind  BindingKind
    Path  string
    Value any
}
```

The TypeScript declarations should mirror these concepts without exposing Go
implementation details:

```ts
export interface WidgetPageSpec {
  kind: "page";
  id: string;
  title: string;
  root: WidgetNodeSpec;
  meta?: JsonObject;
}

export type WidgetNodeSpec =
  | TextNodeSpec
  | ElementNodeSpec
  | ComponentNodeSpec;

export interface ComponentNodeSpec {
  kind: "component";
  type: string;
  props?: JsonObject;
  children?: WidgetNodeSpec[];
  source?: SourceSpan;
}

export interface ActionSpec {
  kind: "server" | "navigate" | "download" | "event" | "copy";
  name?: string;
  to?: string;
  event?: string;
  payload?: Record<string, BindingSpec | JsonValue>;
  confirm?: string;
}

export type BindingSpec =
  | { kind: "field"; path: string }
  | { kind: "context"; path: string }
  | { kind: "map"; field: string }
  | { kind: "template"; template: string }
  | { kind: "const"; value: JsonValue };
```

The clean break means `any` should become rare. It will still appear where arbitrary
JSON props are passed through to raw components, but domain views should have typed
options and typed slot contexts.

---

## 3. The new module shape

Instead of six unrelated CommonJS modules, expose one canonical module:

```ts
import {
  page,
  ui,
  data,
  cms,
  course,
  context,
  schedule,
  time,
  bind,
  act,
  style,
  raw,
} from "widget.dsl";
```

In Goja this can still be implemented as `require("widget.dsl")`. If the host wants
submodules for lazy loading, they can be properties of the root export. The author
should see one coherent namespace.

The namespaces have clear responsibilities:

| Namespace | Responsibility |
|---|---|
| `page` | Build pages, shells, sections, routing metadata, and document structure. |
| `ui` | Generic visual primitives and low-level composition helpers. |
| `data` | Records, collections, schemas, fields, tables, matrices, boards, and record actions. |
| `cms` | Editorial workflows: media libraries, article queues, markdown editing, upload/review/publish intents. |
| `course` | Course studios, lessons, slides, handouts, learner navigation, presenter actions. |
| `context` | Context snapshots, transcript workspaces, annotations, diagrams, analysis views. |
| `schedule` | Availability polls, poll summaries, booking pickers, availability intents. |
| `time` | Month/week/day calendar views, date ranges, formatting, event intents. |
| `bind` | Serializable data access specs. |
| `act` | Transport actions. Domain namespaces add intent wrappers over `act`. |
| `style` | Palettes, style sets, style-by-value specs. |
| `raw` | Explicit escape hatch to build component nodes by type. |

No backward compatibility means no `data.v2.dsl`, no `course.dsl`, no
`context_window.dsl`, and no `recipes` namespace in the public API. Recipes become
domain views. Component factories become engine helpers or raw escape hatches.

---

## 4. Lambdas: three kinds, one rule

The redesign uses lambdas heavily, but it uses them with one strict rule:

> A lambda may configure a Go-backed spec or produce Widget IR during authoring.
> A lambda must not be sent to the browser.

There are three allowed lambda kinds.

### 4.1 Builder lambdas

A builder lambda receives a scoped builder and mutates a Go-backed spec:

```ts
page("Course CMS", p => p
  .shell(course.shell(courseDef))
  .section("Media", s => s.view(cms.mediaLibrary(assets))));
```

Backing model:

```go
func page(title string, cb goja.Value) *goja.Object {
    spec := &PageSpec{Title: title}
    builder := pageBuilder(spec)
    applyBuilderCallback(builder, cb)
    validate(spec)
    return pageValue(spec)
}
```

This follows the successful `researchctl` pattern: scoped builders, fluent methods,
`.use(fragment)`, then validation.

### 4.2 Slot lambdas

A slot lambda renders one named subpart of a view:

```ts
cms.mediaLibrary(assets, m => m.asset(({ asset }, h) =>
  h.card({ tone: asset.status }, [
    h.image(asset.src, { alt: asset.title }),
    h.strong(asset.title),
    h.caption(asset.mime),
  ])));
```

The runtime calls this lambda while lowering the view. The returned node is embedded
in the generated IR.

### 4.3 Fragment lambdas

A fragment is a reusable builder policy:

```ts
const adminChrome: Fragment<PageBuilder> = p => p
  .meta("shell", "course-admin")
  .density("compact")
  .breadcrumb("Course CMS", "/pages/admin-course-cms");

page("Course CMS", p => p.use(adminChrome));
```

Fragments are ordinary functions. They are how we avoid inheritance and avoid
copy/paste.

---

## 5. The minimum TypeScript kernel

The public declarations should be generated from Go descriptor data, but the shape
should look like this:

```ts
export type JsonPrimitive = string | number | boolean | null;
export type JsonValue = JsonPrimitive | JsonValue[] | JsonObject;
export interface JsonObject { [key: string]: JsonValue; }

export type Fragment<TBuilder> = (builder: TBuilder) => void | TBuilder;
export type Slot<TContext> = (context: TContext, h: SlotHelpers) => WidgetChild;
export type WidgetChild = WidgetNodeSpec | string | number | boolean | null | undefined;

export interface SlotHelpers {
  text(value: unknown): WidgetNodeSpec;
  caption(value: unknown, options?: CaptionOptions): WidgetNodeSpec;
  strong(value: unknown): WidgetNodeSpec;
  stack(options: StackOptions, children: WidgetChild[]): WidgetNodeSpec;
  inline(options: InlineOptions, children: WidgetChild[]): WidgetNodeSpec;
  card(options: CardOptions, children: WidgetChild[]): WidgetNodeSpec;
  button(label: string, options?: ButtonOptions): WidgetNodeSpec;
  image(src: string, options?: ImageOptions): WidgetNodeSpec;
  badge(label: string, options?: BadgeOptions): WidgetNodeSpec;
  raw(type: string, props?: JsonObject, children?: WidgetChild[]): WidgetNodeSpec;
}

export interface Buildable<TSpec> {
  validate(): ValidationIssue[];
  toSpec(): TSpec;
}

export interface EmitWidget {
  toNode(): WidgetNodeSpec;
}
```

Every complex view builder should implement `validate()` in TypeScript declarations
because the Go builder has a validation method. Authors can test their DSL without
rendering a browser.

---

# Part A — Simple examples

These examples are intentionally written in TypeScript. They show the intended
shape of the `.d.ts` contract even though the runtime is Goja.

## 6. Example 1: the smallest useful page

```ts
import { page, ui } from "widget.dsl";
import type { WidgetPageSpec } from "widget.dsl";

export default page("Hello", p => p
  .section("Greeting", s => s
    .text("Hello from a Go-backed Widget DSL.")))
  .toPage();
```

Backing idea:

- `page(...)` creates a Go `PageSpec`.
- `.section(...)` appends a Go `SectionSpec`.
- `.text(...)` appends a Go `NodeSpec{Kind: Text}` or a `Text` component depending
  on final renderer convention.
- `.toPage()` validates and lowers to `WidgetPageSpec`.

No raw component type strings appear.

## 7. Example 2: page structure and actions

```ts
import { page, ui, act } from "widget.dsl";
import type { ActionSpec, WidgetPageSpec } from "widget.dsl";

const openSettings: ActionSpec = act.navigate("/pages/settings");
const uploadSession: ActionSpec = act.navigate("/pages/upload");

export default page("Workshop", p => p
  .section("Welcome", s => s
    .caption("A generated page backed by typed Go specs.")
    .actions(a => a
      .button("Settings", openSettings)
      .button("Upload", uploadSession))
    .view(ui.callout("Set your display name before uploading a session.", {
      tone: "info",
    }))))
  .toPage();
```

The action helpers return serializable `ActionSpec` values. They are not callbacks.
The browser can dispatch them because they are data.

## 8. Example 3: reusable fragments

```ts
import { page, ui } from "widget.dsl";
import type { Fragment, PageBuilder, SectionBuilder } from "widget.dsl";

const adminPage: Fragment<PageBuilder> = p => p
  .meta("shell", "admin")
  .density("compact")
  .breadcrumb("Admin", "/pages/admin-course-cms");

const mutedSection: Fragment<SectionBuilder> = s => s
  .tone("quiet")
  .caption("Generated from a shared section fragment.");

export default page("Admin", p => p
  .use(adminPage)
  .section("Status", s => s
    .use(mutedSection)
    .metric("Files", 42, { status: "ready" })))
  .toPage();
```

This is the `researchctl`/`codesign` `.use(fragment)` pattern applied to UI. The
fragment closes over whatever it needs, but it runs immediately. Only the resulting
spec is serialized.

---

# Part B — Data examples

## 9. Example 4: typed schema and table

```ts
import { page, data, act } from "widget.dsl";
import type { FieldSchema, CollectionSpec } from "widget.dsl";

interface SessionRow {
  sessionId: string;
  title: string;
  turnCount: number;
  status: "ready" | "running" | "failed";
  body: string;
}

const sessions: SessionRow[] = [
  { sessionId: "sess-intro", title: "Intro walkthrough", turnCount: 12, status: "ready", body: "Short transcript." },
  { sessionId: "sess-debug", title: "Debugging trace", turnCount: 28, status: "ready", body: "Long transcript." },
];

const sessionFields: FieldSchema<SessionRow> = data.fields<SessionRow>(f => f
  .key("sessionId", { label: "ID", width: "14ch" })
  .primary("title", { label: "Title", required: true, maxLength: 120 })
  .count("turnCount", { label: "Turns" })
  .status("status", { label: "Status" })
  .prose("body", { label: "Body", rows: 3 }));

export default page("Sessions", p => p
  .section("Uploaded sessions", s => s
    .view(data.collection(sessions, c => c
      .schema(sessionFields)
      .table(t => t
        .selectRow(act.navigate("/pages/sessions?selected=${row.sessionId}")))))))
  .toPage();
```

This replaces the current `data.v2.dsl` split with a normal `data.collection` view.
The Go backing type is still a `CollectionSpec`, but the author does not import a
separate module or remember `.toIR()`.

## 10. Example 5: editable master-detail collection

```ts
import { page, data, act, bind } from "widget.dsl";
import type { FieldSchema } from "widget.dsl";

interface AgendaItem {
  id: string;
  number: string;
  duration: string;
  title: string;
  description: string;
}

const agendaFields: FieldSchema<AgendaItem> = data.fields<AgendaItem>(f => f
  .key("id", { label: "ID", width: "18ch" })
  .short("number", { label: "Time", width: "8ch" })
  .short("duration", { label: "Duration", width: "8ch" })
  .primary("title", { label: "Title", required: true, maxLength: 160 })
  .prose("description", { label: "Description", rows: 4, maxLength: 800 }));

export function agendaEditor(agenda: AgendaItem[], selected: string | undefined) {
  return data.collection(agenda, c => c
    .id("agenda")
    .schema(agendaFields)
    .select(data.selection.urlParam("agenda", selected))
    .edit(e => e
      .submitPost("/settings/agenda-item")
      .create({ label: "New agenda item" })
      .reorder(act.server("admin-reorder-course-agenda"))
      .remove(act.server("admin-delete-agenda-item", {
        confirm: "Delete ${row.title}? This cannot be undone.",
        payload: { id: bind.context("row.id") },
      })))
    .masterDetail(md => md
      .detailTitle(({ row }) => row ? `Edit ${row.title}` : "Select an item")));
}

export default page("Agenda", p => p
  .section("Workshop agenda", s => s.view(agendaEditor(agenda, query.agenda))))
  .toPage();
```

The current `go-go-course` version uses `dataV2.collection(...).edit(...).masterDetail().toIR()`. The redesigned API keeps the good builder pattern and removes the separate module and terminal `.toIR()` call from normal page composition.

## 11. Example 6: a matrix engine for advanced authors

```ts
import { page, data, style, bind, act } from "widget.dsl";

interface VoteRow {
  id: string;
  name: string;
  cells: Record<string, "yes" | "ifneedbe" | "no" | "unknown">;
}

const availability = style.set("availability", s => s
  .entry("yes", "Yes", { fill: "var(--green)", labelColor: "white" })
  .entry("ifneedbe", "If need be", { fill: "var(--amber)", labelColor: "black" })
  .entry("no", "No", { fill: "var(--red)", labelColor: "white" })
  .entry("unknown", "Unknown", { fill: "var(--surface)", labelColor: "var(--text-muted)" }));

export default page("Raw matrix", p => p
  .section("Engine-level matrix", s => s
    .view(data.matrix<VoteRow>(responses, m => m
      .columns(options.map(o => ({ id: o.id, header: o.label })))
      .rowKey("id")
      .rowHeader(data.cell.field("name"))
      .valueAt(bind.map("cells"))
      .cell(data.cell.cycle(["yes", "ifneedbe", "no", "unknown"], {
        glyphs: { yes: "✓", ifneedbe: "~", no: "✕", unknown: "·" },
        styleSet: availability,
      }))
      .editableRow(currentResponseId)
      .onCell(act.server("poll.toggle", {
        payload: {
          row: bind.context("rowKey"),
          col: bind.context("colId"),
          value: bind.context("value"),
        },
      })))))
  .toPage();
```

Engine helpers are still available, but they are typed builders, not arbitrary prop
bags.

---

# Part C — Scheduling and calendar examples

## 12. Example 7: product-level availability poll

```ts
import { page, schedule, time, ui } from "widget.dsl";
import type { AvailabilityPoll, AvailabilityResponse, PollOption } from "widget.dsl/schedule";

const compactPoll = schedule.fragments.availabilityPoll(p => p
  .density("compact")
  .legend(false));

export default page("Find a time", p => p
  .section("Choose the times that work", s => s
    .view(schedule.availabilityPoll(poll, pollView => pollView
      .use(compactPoll)
      .currentResponse(currentUser.responseId)
      .readOnly(poll.closed)
      .option(({ option }, h) => h.stack({ gap: "xxs", align: "center" }, [
        h.strong(time.format(option.slot.startISO, "EEE")),
        h.caption(time.formatRange(option.slot.startISO, option.slot.endISO, "HH:mm")),
      ]))
      .participant(({ response }, h) => h.inline({ gap: "xs" }, [
        h.avatar(response.name),
        h.text(response.name),
      ]))
      .onChange(schedule.intent.toggleAvailability("poll.toggle"))
      .onSubmit(schedule.intent.submitResponse("poll.submit")))))
  .section("Best options", s => s
    .view(schedule.pollSummary(poll, tallies, summary => summary
      .order("best-first")
      .maxOptions(5)
      .label(({ option, tally }, h) => h.inline({ gap: "xs" }, [
        time.slotLabel(option.slot),
        h.badge(`${tally.yes} yes`),
      ]))))))
  .toPage();
```

This is the preferred level for product code. It does not expose `MatrixGrid` unless
the author asks for the engine.

## 13. Example 8: calendar week

```ts
import { page, time } from "widget.dsl";
import type { CalendarEvent } from "widget.dsl/time";

const events: CalendarEvent[] = loadEvents();

export default page("Calendar", p => p
  .section("Week", s => s
    .view(time.week(events, week => week
      .range(time.range.week("2026-07-06"))
      .hours(8, 18)
      .event(({ event }, h) => h.card({ tone: event.colorKey }, [
        h.strong(event.title),
        h.caption(time.formatRange(event.startISO, event.endISO, "HH:mm")),
      ]))
      .onSelect(time.intent.selectEvent("calendar.selectEvent")))))
  .toPage();
```

The view owns day generation, block conversion, lane packing, and selection context.
All-day events should not be exposed until the backing `TimeGrid` view has real
all-day rendering.

---

# Part D — CMS examples

## 14. Example 9: media library, replacing `cmsDsl.recipes.mediaLibrary`

The current `go-go-course` page builds a media library with an option-bag recipe and
then appends a separate selected-asset detail panel. The redesign makes the media
library an editorial view with slots and intents.

```ts
import { page, cms, ui, act } from "widget.dsl";
import type { CmsAsset } from "widget.dsl/cms";

const mediaAdmin = cms.fragments.mediaLibrary<CmsAsset>(m => m
  .selection("multi")
  .empty("No media files yet. Upload one to get started.")
  .accept("image/svg+xml,image/png,image/jpeg,image/webp,image/gif"));

export default page("Course CMS", p => p
  .section("Media library", s => s
    .caption("Files under the active site media directory.")
    .view(cms.mediaLibrary(material.mediaAssets, library => library
      .use(mediaAdmin)
      .selected(query.asset ? [query.asset] : [])
      .asset(({ asset, selected }, h) => h.card({ tone: selected ? "selected" : "default" }, [
        h.image(asset.src, { alt: asset.title, fit: "contain" }),
        h.strong(asset.filename),
        h.caption(`${asset.mime} · ${cms.formatBytes(asset.size)}`),
      ]))
      .details(({ asset }, h) => h.stack({ gap: "md" }, [
        h.metadata({
          File: asset.filename,
          "MIME type": asset.mime,
          Size: cms.formatBytes(asset.size),
          URL: asset.src,
        }),
        h.inline({ gap: "sm", wrap: true }, [
          h.button("Open", { action: act.navigate(asset.src) }),
          h.button("Download", { action: act.download(asset.src) }),
          h.button("Delete", {
            action: cms.intent.deleteAsset("admin-delete-course-material", {
              confirm: `Delete ${asset.filename}? This permanently removes the file.`,
            }),
          }),
        ]),
      ]))
      .onSelect(cms.intent.selectAsset("cms.selectAsset", {
        navigateTo: "/pages/admin-course-cms?asset=${asset.id}",
      }))
      .onOpen(cms.intent.openAsset("cms.openAsset"))
      .onUpload(cms.intent.uploadAssets("admin-upload-course-material", {
        payload: { kind: "media", overwrite: false },
      })))))
  .toPage();
```

This example demonstrates the design goal: CMS code names CMS tasks. The lowerer may
still emit `MediaLibraryPanel`, `AssetTile`, `MediaThumb`, `Button`, and
`MetadataGrid`, but the script is not an assembly of those parts.

## 15. Example 10: article queue

```ts
import { page, cms } from "widget.dsl";
import type { CmsArticleSummary } from "widget.dsl/cms";

export default page("Editorial queue", p => p
  .section("Needs review", s => s
    .view(cms.articleQueue(articles, q => q
      .statusFilter("needs-review")
      .search(query.q || "")
      .page(Number(query.page || 1), pageCount)
      .article(({ article }, h) => h.inline({ gap: "sm" }, [
        h.badge(article.status),
        h.stack({ gap: "xxs" }, [
          h.strong(article.title),
          h.caption(`${article.author} · ${article.updatedAt}`),
        ]),
      ]))
      .rowActions(actions => actions
        .action("Preview", cms.intent.previewArticle("cms.previewArticle"))
        .action("Publish", cms.intent.publishArticle("cms.publishArticle", {
          confirm: "Publish ${article.title}?",
        }))
        .action("Archive", cms.intent.archiveArticle("cms.archiveArticle", {
          confirm: "Archive ${article.title}?",
        })))
      .onSelect(cms.intent.selectArticle("cms.selectArticle", {
        navigateTo: "/pages/articles?article=${article.id}",
      })))))
  .toPage();
```

The old `articleListPanel` becomes an implementation detail. The public API is an
editorial queue.

## 16. Example 11: markdown editor

```ts
import { page, cms, act } from "widget.dsl";

export default page("Edit handout", p => p
  .section("Handout body", s => s
    .view(cms.markdownEditor(loaded.body, editor => editor
      .name("body")
      .fileName(file)
      .preview("split")
      .submitPost("/settings/handout-body")
      .hidden("file", file)
      .status(query.status === "saved" ? "success" : query.status === "error" ? "error" : "idle")
      .onCancel(act.navigate("/pages/admin-course-cms")))))
  .toPage();
```

This replaces hand-built `ui.formPanel(...)` plus `ui.textareaInput(...)` where the
product intent is specifically editing markdown with preview.

---

# Part E — Course examples

## 17. Example 12: course shell page, replacing `courseDsl.recipes.courseStudio`

```ts
import { page, course } from "widget.dsl";
import type { CourseDefinition, CourseNavItem } from "widget.dsl/course";

const courseChrome = (definition: CourseDefinition, user: User) =>
  course.fragments.shell(shell => shell
    .title(definition.title)
    .subtitle(definition.shellSubtitle)
    .sections(nav => nav
      .section("Course Material", s => s
        .item("course", "Course", { icon: "course" })
        .item("slides", "Slides", { icon: "slides" })
        .item("handouts", "Handouts", { icon: "handouts" }))
      .section("Settings", s => s.item("settings", "Settings", { icon: "settings" })))
    .onNavigate(course.intent.navigateItem("course.navigate", {
      navigateTo: "/pages/${item.id}",
    })));

export default page("Course", p => p
  .shell(course.shell(definition, shell => shell
    .use(courseChrome(definition, user))
    .active("course")))
  .section("Overview", s => s
    .view(course.landing(definition, landing => landing
      .primaryCta("View slides", course.intent.openSlides())
      .secondaryCta("View handouts", course.intent.openHandouts())))))
  .toPage();
```

The shell is not a recipe object anymore. It is a page shell builder with a typed
navigation model.

## 18. Example 13: slides and presenter mode

```ts
import { page, course, context } from "widget.dsl";
import type { CourseSlide } from "widget.dsl/course";

export default page(`Slides · ${selected.title}`, p => p
  .shell(course.shell(definition, s => s.active("slides")))
  .section("Slide deck", s => s
    .view(course.slideDeck(deck, deck => deck
      .selected(selected.id)
      .mode(query.presenter ? "presenter" : "reader")
      .visual(({ slide }, h) => slide.snapshot
        ? context.diagram(slide.snapshot, d => d.view(slide.view || "budget"))
        : h.callout("No context snapshot for this slide."))
      .notes(({ slide }, h) => h.markdown(slide.notes.join("\n")))
      .onPrevious(course.intent.previousSlide("course.previous"))
      .onNext(course.intent.nextSlide("course.next"))
      .onPresent(course.intent.presentSlide("course.present")))))
  .toPage();
```

The slide deck owns slide navigation and presenter actions. Context diagrams are
composed as views inside the slide visual slot.

## 19. Example 14: handouts

```ts
import { page, course, context } from "widget.dsl";
import type { HandoutBundle } from "widget.dsl/course";

export default page("Handouts", p => p
  .shell(course.shell(definition, s => s.active("handouts").contentPadding("none")))
  .section("Handouts", s => s
    .view(course.handouts(bundle, handouts => handouts
      .selected(documentId || bundle.docs[0]?.id)
      .style(context.palette("Signal Orange / Cyan"))
      .document(({ document }, h) => h.stack({ gap: "xs" }, [
        h.strong(document.title),
        h.caption(document.description),
      ]))
      .onSelect(course.intent.selectHandout("course.selectHandout", {
        navigateTo: "/pages/handouts?doc=${document.id}",
      }))
      .onDownload(course.intent.downloadHandout("course.downloadHandout", {
        to: "/api/handouts/${document.id}/download.md",
      }))
      .onPrint(course.intent.printHandout("course.printHandout", {
        navigateTo: "/pages/print-handout?doc=${document.id}",
      })))))
  .toPage();
```

This replaces direct use of `courseDsl.handoutDocumentShell({ ... })` with a typed
handout domain view.

---

# Part F — Context-window examples

## 20. Example 15: transcript workspace

```ts
import { page, context, course } from "widget.dsl";
import type { Transcript, ContextSnapshot, Annotation } from "widget.dsl/context";

export default page("Session analysis", p => p
  .shell(course.shell(definition, s => s.active("sessions")))
  .section("Context workspace", s => s
    .view(context.workspace(session, w => w
      .snapshot(snapshot)
      .transcript(transcript)
      .annotations(annotations)
      .layout("diagram-plus-transcript")
      .diagram(d => d
        .view("budget")
        .selected(query.part || undefined)
        .style(context.palette("Signal Orange / Cyan")))
      .message(({ message }, h) => h.transcriptMessage(message, {
        showTokenCount: true,
        compactTools: true,
      }))
      .annotation(({ annotation }, h) => h.annotationCard(annotation, {
        emphasis: annotation.severity,
      }))
      .onPartSelect(context.intent.selectPart("context.selectPart", {
        navigateTo: "/pages/session-visualize?part=${part.id}",
      }))
      .onAnnotationSelect(context.intent.selectAnnotation("context.selectAnnotation")))))
  .toPage();
```

Current `context_window.dsl` exposes many panel helpers. The redesigned API makes
the analysis workspace the primary concept and the panels implementation details.

## 21. Example 16: standalone context diagram

```ts
import { page, context } from "widget.dsl";

export default page("Context diagram", p => p
  .section("Budget", s => s
    .view(context.diagram(snapshot, d => d
      .view("treemap")
      .style(context.styleSet(style => style
        .entry("prompt", "Prompt", { accent: "a", pattern: "solid" })
        .entry("answer", "Answer", { accent: "b", pattern: "checker" }))
      .legend(true))))
  .toPage();
```

The style-set builder is typed and reusable; it replaces a mixture of style helper
functions and open objects.

---

# Part G — Complex page examples inspired by `go-go-course`

## 22. Example 17: full Course CMS page

This is the clean-break version of the current `admin-course-cms.js` shape. It uses
one grammar for page shell, metadata forms, agenda editing, course-material uploads,
media management, and preview actions.

```ts
import { page, ui, data, cms, course, act } from "widget.dsl";
import type {
  CmsAsset,
  CourseDefinition,
  CourseMaterialIndex,
  EditableCourseMetadata,
  AgendaItem,
} from "widget.dsl/domain";

interface AdminCourseCmsInput {
  user: User;
  definition: CourseDefinition;
  metadata: EditableCourseMetadata;
  material: CourseMaterialIndex;
  query: Record<string, string | undefined>;
}

const adminOnly = (input: AdminCourseCmsInput) =>
  course.fragments.adminShell(input.definition, shell => shell
    .active("admin-course-cms")
    .requireAdmin(input.user, {
      fallback: ui.callout("Set your display name to admin_<name> to reveal course editing controls.", {
        tone: "warning",
        action: act.navigate("/pages/settings"),
      }),
    }));

export function adminCourseCmsPage(input: AdminCourseCmsInput) {
  const { definition, metadata, material, query } = input;

  return page("Course CMS", p => p
    .use(adminOnly(input))

    .section("CMS building blocks", s => s
      .caption("A single admin workspace for metadata, ordered course lists, and file-backed assets.")
      .metadata({
        "Content root": material.contentRoot,
        "Metadata store": material.metadataPath,
        "Slide source": `${material.slideDir}/*.md`,
        "Handout source": `${material.handoutDir}/*.md`,
      }))

    .section("Landing page", s => s
      .view(course.metadataForm(metadata, form => form
        .submitPost("/settings/course-metadata")
        .statusFromQuery(query, "saved")
        .fields(fields => fields
          .text("kicker", { label: "Kicker", maxLength: 120 })
          .text("title", { label: "Title", required: true, maxLength: 140 })
          .textarea("tagline", { label: "Tagline", maxLength: 500, rows: 3 })
          .textarea("blurb", { label: "Blurb", maxLength: 1200, rows: 6 })
          .text("instructorName", { label: "Instructor name", maxLength: 120 })
          .textarea("instructorBio", { label: "Instructor bio", maxLength: 700, rows: 4 })))))

    .section("Agenda", s => s
      .view(course.agendaEditor(metadata.agenda, editor => editor
        .selected(query.agenda)
        .submitPost("/settings/agenda-item")
        .onReorder(course.intent.reorderAgenda("admin-reorder-course-agenda"))
        .onRemove(course.intent.deleteAgendaItem("admin-delete-agenda-item", {
          confirm: "Delete agenda item ${item.title}? This cannot be undone.",
        })))) )

    .section("Slides and handouts", s => s
      .view(course.materialUploads(material, uploads => uploads
        .slideUpload(course.intent.uploadMaterial("admin-upload-course-material", { kind: "slide" }))
        .handoutUpload(course.intent.uploadMaterial("admin-upload-course-material", { kind: "handout" }))
        .slidesTable(material.slides)
        .handoutsTable(material.handouts, table => table
          .editLink(file => `/pages/admin-handout-editor?file=${encodeURIComponent(file)}`)))))

    .section("Media library", s => s
      .view(cms.mediaLibrary(material.mediaAssets, library => library
        .selected(query.asset ? [query.asset] : [])
        .onUpload(cms.intent.uploadAssets("admin-upload-course-material", { kind: "media" }))
        .onSelect(cms.intent.selectAsset("cms.selectAsset", {
          navigateTo: "/pages/admin-course-cms?asset=${asset.id}",
        }))
        .details(({ asset }, h) => h.stack({ gap: "md" }, [
          h.metadata({ File: asset.filename, Size: cms.formatBytes(asset.size), URL: asset.src }),
          h.inline({ gap: "sm" }, [
            h.button("Open", { action: act.navigate(asset.src) }),
            h.button("Delete", { action: cms.intent.deleteAsset("admin-delete-course-material") }),
          ]),
        ])))))

    .section("Preview", s => s
      .actions(a => a
        .button("Open course page", act.navigate("/pages/course"))
        .button("Preview slides", act.navigate("/pages/slides"))
        .button("Preview handouts", act.navigate("/pages/handouts")))))
    .toPage();
}
```

This example is intentionally ambitious. It shows why the redesign should target
all DSLs, not schedule only. The page composes course, CMS, data, and UI concepts
without switching authoring styles.

## 23. Example 18: full workshop app router

The DSL should also support reusable page factories. A route handler can call a
page factory and return the serialized page.

```ts
import type { WidgetPageSpec } from "widget.dsl";
import { page, course, context, schedule } from "widget.dsl";

interface RouteContext {
  user: User;
  query: Record<string, string | undefined>;
  course: CourseDefinition;
  material: CourseMaterialIndex;
  session?: UploadedSession;
}

type PageFactory = (ctx: RouteContext) => WidgetPageSpec;

const pages: Record<string, PageFactory> = {
  course: ctx => course.landingPage(ctx.course, { user: ctx.user }),
  handouts: ctx => course.handoutsPage(ctx.material.handouts, { selected: ctx.query.doc }),
  "admin-course-cms": ctx => adminCourseCmsPage(ctx),
  "session-visualize": ctx => context.sessionWorkspacePage(ctx.session!, { selectedPart: ctx.query.part }),
  "schedule-demo": ctx => schedule.availabilityPollPage(loadPoll(), { user: ctx.user }),
};

export function buildWidgetPage(route: string, ctx: RouteContext): WidgetPageSpec {
  const factory = pages[route] || pages.course;
  return factory(ctx);
}
```

This is the level where Goja shines. Page factories use ordinary JavaScript control
flow, but every returned page is validated Go-backed Widget IR.

---

# Part H — Raw escape hatch and custom engines

A clean DSL still needs an escape hatch. The difference is that the escape hatch is
not named `component` and not taught as the normal path.

```ts
import { page, raw } from "widget.dsl";

export default page("Experimental", p => p
  .section("Prototype", s => s
    .view(raw.component("ExperimentalWidget", {
      featureFlag: "next-dashboard",
      payload: { a: 1 },
    }))))
  .toPage();
```

If a raw component becomes common, it should graduate into an engine or domain view:

```ts
// Temporary
raw.component("PipelineBoard", props)

// Engine-level
data.board(deals, b => b.columns(stages).card(...))

// Product-level
crm.pipeline(deals, p => p.stages(stages).onMove(crm.intent.moveDeal("deal.move")))
```

That graduation path is part of the design.

---

# Part I — Implementation design

## 24. One descriptor registry

The Go runtime should not hand-code TypeScript declarations independently from
runtime exports. Define descriptors once:

```go
type NamespaceDescriptor struct {
    Name    string
    Doc     string
    Views   []ViewDescriptor
    Specs   []SpecDescriptor
    Intents []IntentDescriptor
    Raw     []RawComponentDescriptor
}

type ViewDescriptor struct {
    Name       string
    DataType   string
    Builder    BuilderDescriptor
    Slots      []SlotDescriptor
    LowerTo    string // component type or lowering function id
    Doc        string
}

type BuilderDescriptor struct {
    GoType  reflect.Type
    Methods []BuilderMethodDescriptor
}
```

Use the descriptor for three outputs:

1. Install Goja runtime functions.
2. Generate TypeScript declarations.
3. Generate documentation and examples.

This avoids the drift visible in the old system.

## 25. Builder handles wrap Go specs

A builder method mutates a Go spec and returns the same builder:

```go
type AvailabilityPollSpec struct {
    Poll            PollDTO
    CurrentResponse string
    ReadOnly        bool
    Slots           map[string]SlotSpec
    OnChange        *ActionSpec
    OnSubmit        *ActionSpec
}

func (r *runtime) availabilityPoll(pollValue goja.Value, configure goja.Value) goja.Value {
    spec := &AvailabilityPollSpec{Poll: decodePoll(pollValue)}
    builder := r.availabilityPollBuilder(spec)
    r.applyOptionsOrBuilder(builder, spec, configure)
    r.validateAvailabilityPoll(spec)
    return r.vm.ToValue(r.lowerAvailabilityPoll(spec))
}
```

Options objects may be allowed for tiny views, but they should normalize into the
same spec as builder lambdas. There must not be two lowering paths.

## 26. Slot invocation is explicit

```go
type SlotSpec struct {
    Name     string
    Function goja.Value // author-time only
    Node     *NodeSpec
    Binding  *BindingSpec
}
```

When lowering a view, the Go code calls slots with a typed context map and a helper
object:

```go
func (r *runtime) callSlot(slot SlotSpec, ctx any, fallback func(any) NodeSpec) NodeSpec {
    if slot.Function != nil {
        fn, _ := goja.AssertFunction(slot.Function)
        value, err := fn(goja.Undefined(), r.vm.ToValue(ctx), r.slotHelpers())
        if err != nil { panic(err) }
        return r.mustNode(value)
    }
    return fallback(ctx)
}
```

The slot helper object should be stable across domains. That makes examples easier
and avoids requiring every slot to close over `ui`.

## 27. Validation happens before lowering

Every view spec validates at the domain level:

- A CMS media library cannot have `selection("single")` and multiple selected IDs.
- A course handout view cannot select a document that does not exist unless it has
  an empty-state fallback.
- A schedule availability poll cannot emit edit actions in read-only mode.
- A time week view cannot expose all-day slots until all-day rendering exists.
- A data collection cannot render a table without a schema.

Then the lowerer validates the generated Widget IR shape.

## 28. Browser action contexts are domain-shaped

The old engine context names leak into scripts: `rowKey`, `colId`, `assetId`,
`documentId`. In the new design, domain intents own context translation.

```ts
schedule.intent.toggleAvailability("poll.toggle")
```

means the author sees:

```ts
payload: {
  pollId: bind.context("poll.id"),
  responseId: bind.context("response.id"),
  optionId: bind.context("option.id"),
  value: bind.context("value"),
}
```

The lowerer may translate from `MatrixGrid` context internally. The product script
should not know that.

---

# Part J — Migration by replacement, not compatibility

No backward compatibility means we do not need aliases, but we still need an orderly
replacement plan.

1. **Build `widget.dsl` alongside the old modules.** Do not delete old modules until
   the new module can render real pages.
2. **Port `go-go-course` DSL examples first.** They exercise UI, data collections,
   CMS, course shells, and context diagrams in one place.
3. **Port the Course CMS page.** This is the best integration test because it mixes
   forms, uploads, editable collections, media library, selected detail, and course
   shell composition.
4. **Port handouts and slides.** These validate course domain views and context
   diagrams inside course slots.
5. **Port session visualization.** This validates context workspace views.
6. **Port scheduling demos.** This validates schedule/time views and matrix/calendar
   engines.
7. **Cut over provider declarations.** The xgoja provider should expose only
   `widget.dsl` for new hosts.
8. **Delete old public modules after all first-party pages move.** Internals may keep
   lowerers, but public imports should converge.

The important migration target is not API compatibility. It is example coverage. If
all real `go-go-course` pages can be expressed in the new grammar and produce the
same or better Widget IR, the redesign is ready.

---

# Part K — Acceptance criteria

The DSL redesign is successful when these are true:

1. A TypeScript declaration file generated from Go descriptors lets authors write
   all examples in this document with useful completion.
2. Every builder lambda mutates a typed Go spec and validates before lowering.
3. Every slot lambda runs only at authoring time and returns serializable Widget IR.
4. Every domain intent lowers to a transport `ActionSpec` with a documented context.
5. No example imports legacy `ui.dsl`, `data.dsl`, `data.v2.dsl`, `cms.dsl`,
   `course.dsl`, or `context_window.dsl` modules.
6. The Course CMS page, handouts page, slide page, session workspace page, and a
   scheduling poll page are implemented as golden examples.
7. Golden examples execute in Goja, emit JSON, and compare against stable Widget IR
   snapshots.
8. The browser renderer does not need to know whether a node came from a raw engine,
   a domain view, or a builder slot.

---

# Part L — What to implement first

The first slice should prove the full design without requiring every domain module.

1. Add `widget.dsl` root module.
2. Add `page(...)`, `ui` helpers, `act`, `bind`, and `raw`.
3. Add builder callback infrastructure and `.use(fragment)`.
4. Add `data.fields` and `data.collection` over the existing v2 `CollectionSpec`.
5. Add `cms.mediaLibrary` over the existing `MediaLibraryPanel` lowerer.
6. Add TypeScript declarations for those pieces.
7. Port `go-go-course` `/pages/dsl-examples-*` and the media-library part of
   `/pages/admin-course-cms` as test fixtures.

That slice hits every essential mechanism: page building, domain view, typed data
schema, slot lambdas, fragments, actions, bindings, validation, lowering, and
TypeScript generation.

---

## References

- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js` — current in-the-wild use of `ui.dsl`, `data.v2.dsl`, `context_window.dsl`, `cms.dsl`, and `course.dsl` on one demo page.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/lib/pages/admin-course-cms.js` — current Course CMS page combining metadata forms, agenda editor, uploads, media library, and asset details.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/lib/pages/common.js` — current course shell helper around `courseDsl.recipes.courseStudio`.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/lib/pages/handouts.js` — current handout page using `courseDsl.handoutDocumentShell` and `richArticle`.
- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/researchctl/builders.go` — scoped builder callback pattern.
- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/builders.go` — `.use(fragment)` builder pattern.
- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/callbacks.go` — registered callback-ID boundary for trusted server-side callbacks.
- `pkg/widgetdsl/module.go` — current split module registry, component factories, action/cell helpers, and hand-written recipes.
- `pkg/widgetdsl/typescript.go` — current generated `.d.ts` surface and evidence of why descriptors should generate runtime and TypeScript together.
- `pkg/widgetdsl/v2_builders.go` and `pkg/widgetdsl/v2/spec/types.go` — useful typed-builder and Go-spec foundation to retain, but under the new `widget.dsl` surface.
