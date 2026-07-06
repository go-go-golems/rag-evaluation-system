---
Title: Goja DSL deep-dive — optional lambdas, typed IR, DTS parity, and tag operators
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
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/dts_parity_test.go
      Note: Generated DTS versus runtime export parity implementation
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/go-emrichen/pkg/emrichen/emrichen.go
      Note: TagFunc model
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/go-emrichen/pkg/emrichen/parser.go
      Note: Strict tag argument parser with unknown-key and required-key errors
    - Path: ../../../../../../../../../2026-06-30/benchmark-cpu-inference/researchctl/pkg/codesign/spec/types.go
      Note: Typed Go-side codesign RunSpec and ValidationResult
    - Path: ../../../../../../../../../2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/typescript.go
      Note: FragmentFn and precise builder declaration model
    - Path: ../../../../../../../../../2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/researchctl/builders.go
      Note: Evidence that researchctl entity lambdas are optional configurators
    - Path: ../../../../../../../../../2026-06-30/benchmark-cpu-inference/researchctl/pkg/research/spec/types.go
      Note: Typed Go-side ResearchProjectSpec and entity IR
    - Path: packages/rag-evaluation-site/src/widgets/ir.ts
      Note: Typed TS-side Widget IR source that widgetdsl does not yet consume
    - Path: pkg/widgetdsl/typescript.go
      Note: Current weak widgetdsl declaration generation
ExternalSources: []
Summary: 'Intern-facing deep dive answering four design questions for the Goja DSL playbook: whether researchctl lambdas are optional, where typed Go-side IR/specs exist, how geppetto''s DTS parity test works, and what the go-emrichen tag-operator model contributes to nested/composable DSL design.'
LastUpdated: 2026-07-05T14:10:00-04:00
WhatFor: Use when designing the next Goja fluent-builder DSL surface, especially to balance simple opinionated helpers with optional lambda configurators and typed specs.
WhenToUse: Read after the catalogue and widgetdsl self-assessments, before writing the final fluent-builder DSL playbook or refactoring widgetdsl.
---


# Goja DSL deep-dive — optional lambdas, typed IR, DTS parity, and tag operators

## Executive summary

This document answers four follow-up design questions for the Goja DSL playbook.

1. **Are the lambdas in `researchctl` optional?** Yes. The implementation accepts entity-builder callbacks as variadic optional arguments and explicitly returns without doing anything when the callback is missing, `undefined`, or `null` (`researchctl/builders.go:106-109`). The public TypeScript declaration also marks every `build` callback optional (`build?: ...`, `researchctl/module.go:27-29`). This means a useful, opinionated simple surface is possible: `project("X").goal("G").experiment("E")` works as a defaulted graph, and lambdas are an escape hatch for ids, statuses, links, evidence, and custom metadata. One caveat: a non-function callback is silently ignored today (`researchctl/builders.go:110-113`), which is not a good playbook rule.

2. **Do we have typed Go-side IR anywhere?** Yes, but not in the current `widgetdsl` authoring layer. `researchctl` and `codesign` both have typed Go spec/IR structs (`ResearchProjectSpec`, `GoalSpec`, `RunSpec`, `TopologySpec`, `WorkloadSpec`, `ValidationResult`) that builders mutate and terminals return. `rag-evaluation-system` has precise TypeScript-side Widget IR (`WidgetNode`, `ComponentNode`, `ActionSpec`, `WidgetProps`) and Go-side widget-manifest metadata, but `pkg/widgetdsl` still emits `map[string]any` and generates `Props = Record<string, any>`. So the answer is: typed Go-side IR exists in the research/codesign DSLs; widgetdsl has type sources, but they are not yet the Go-side authoring IR.

3. **How does geppetto's DTS parity test work?** It is a runtime/export-surface parity test, not a full TypeScript semantic checker. The test parses the generated `pkg/doc/types/geppetto.d.ts` with regexes to collect exported top-level names and object-namespace properties, instantiates a Goja runtime, calls `Object.keys(require("geppetto"))`, and compares the two sorted sets. It also checks selected namespaces (`consts`, `inferenceProfiles`, `schema`, `turnStores`). This catches drift where runtime exports and generated declarations disagree, which is exactly the kind of drift widgetdsl currently permits.

4. **What is go-emrichen's tag-operator model?** Go-emrichen treats YAML tags as composable operators over a YAML node tree. Each tag is a `TagFunc` from `(interpreter, yaml node)` to `(yaml node, error)`. Tags can evaluate nested child nodes by recursively calling `Process`, can introduce lexical variables with `!With` or loop variables with `!Loop`, can validate their mapping arguments with `ParseArgs`, and can be nested naturally because YAML trees can contain tagged nodes anywhere. This is not “lambdas” in the JavaScript-function sense; the closest analogue is “template subnodes evaluated under a scoped environment.” For Goja DSL design, the lesson is not to copy YAML syntax, but to copy the operator architecture: small typed operators, strict argument parsing, recursive composition, scoped sub-evaluation, and extension through registered custom operators.

The playbook conclusion: **lambdas should be optional configurators, not mandatory ceremony.** The simple surface should build a valid object with defaults. Lambdas and `.use(fragment)` should be available when the author needs identity, links, substructure, reusable fragments, or custom callbacks. The target architecture is “typed specs underneath, simple helper methods on top, optional lambda fragments for advanced composition.”

## Problem statement and scope

The previous catalogue and self-assessment established that `widgetdsl`'s current grammar is valuable as vocabulary but weak as a substrate. The user then asked four more precise questions:

- Are researchctl lambdas optional, and can we still have an opinionated simple surface?
- Do we already have typed Go-side IR somewhere?
- How exactly does geppetto's DTS parity test work?
- What is the go-emrichen tag-operator model, especially the idea that operators can nest and accept lambdas?

This document is written for an intern who needs to implement or review the next playbook. It is evidence-backed, but it is also explanatory: it maps concrete code to design rules.

Out of scope:

- It does not implement widgetdsl v2.
- It does not exhaustively catalogue every go-emrichen tag.
- It does not propose a final public API name-by-name.

## 1. Researchctl lambdas: optional by design, useful as an escape hatch

### 1.1 What the public API looks like

The README shows the canonical “rich” form:

```js
const { project } = require("researchctl");

module.exports = project("Example JS project")
  .goal("Choose a backend", g => g.id("GOAL-001").status("active").priority("P1"))
  .hypothesis("Simulation gives enough signal", h => h.id("H-001").status("open").priority("P1"))
  .experiment("Run simulation", e => e.id("EXP-001").status("planned").tests("H-001"));
```

The important subtlety is that this example uses lambdas, but the implementation does not require them.

The TypeScript declaration says every entity method takes an optional `build` callback:

```ts
interface ProjectBuilder {
  goal(title: string, build?: (g: any) => any): this;
  question(text: string, build?: (q: any) => any): this;
  hypothesis(claim: string, build?: (h: any) => any): this;
  // ...
}
```

Evidence: `pkg/gojamodules/researchctl/module.go:23-35`, especially line 29.

### 1.2 What the implementation does

Every entity method follows the same pattern. For example, `goal` creates a `GoalSpec` with default status and priority, applies an optional entity builder, appends the entity to the project, and returns the same project builder:

```go
set("goal", func(title string, cb ...goja.Value) (*goja.Object, error) {
    e := spec.GoalSpec{Title: title, Status: spec.StatusDraft, Priority: spec.PriorityP2}
    if err := m.applyEntityBuilder(m.goalBuilder(&e), cb...); err != nil {
        return nil, err
    }
    p.Goals = append(p.Goals, e)
    return obj, nil
})
```

Evidence: `pkg/gojamodules/researchctl/builders.go:15-22`.

The optionality is in `applyEntityBuilder`:

```go
func (m *moduleRuntime) applyEntityBuilder(builder *goja.Object, cb ...goja.Value) error {
    if len(cb) == 0 || goja.IsUndefined(cb[0]) || goja.IsNull(cb[0]) {
        return nil
    }
    fn, ok := goja.AssertFunction(cb[0])
    if !ok {
        return nil
    }
    _, err := fn(goja.Undefined(), builder)
    return err
}
```

Evidence: `pkg/gojamodules/researchctl/builders.go:106-115`.

This means all of these are accepted:

```js
project("Simple")
  .goal("Choose a backend")
  .hypothesis("Simulation is enough")
  .experiment("Run one simulation")
  .toSpec()
```

```js
project("Mixed")
  .goal("Choose a backend")
  .hypothesis("Simulation is enough", h => h.id("H-001"))
  .experiment("Run one simulation", e => e.tests("H-001"))
  .toSpec()
```

The lambdas are optional configurators. They are not the only way to make the object useful.

### 1.3 What the simple surface gives you

The simple no-lambda surface is useful because the builder sets defaults:

- `goal(title)` defaults to `status: draft`, `priority: P2` (`builders.go:15-16`).
- `question(text)` defaults to `status: draft`, `priority: P2` (`builders.go:23-24`).
- `hypothesis(claim)` defaults to `status: open`, `priority: P2`, `confidence: unknown` (`builders.go:31-32`).
- `workPackage(title)` and `experiment(title)` similarly default status and priority (`builders.go:39-49`).
- `source(title)` defaults to `status: unread` (`builders.go:55-56`).
- `evidence(summary)` defaults to `status: raw` (`builders.go:63-64`).
- `decision(title)` defaults to `status: proposed`, `confidence: unknown` (`builders.go:71-72`).

The typed Go spec defines these fields explicitly: `ResearchProjectSpec` contains typed slices for goals, questions, hypotheses, work packages, experiments, sources, evidence, decisions, reports, review rules, and views (`pkg/research/spec/types.go:56-73`). Individual types like `GoalSpec` carry fields such as `ID`, `Title`, `Status`, `Priority`, `Asks`, and `Tags` (`types.go:82-89`).

That gives us a very useful design principle:

> A good DSL should let the author write the shortest sentence that creates a valid default object. Lambdas should refine the object, not be required to make the object exist.

### 1.4 What lambdas add

The lambda configurator is useful when the author needs:

1. stable IDs (`g.id("GOAL-001")`),
2. status or priority overrides (`g.status("active").priority("P0")`),
3. graph references (`q.hypothesize("H-001")`, `e.tests("H-001")`),
4. tags and metadata,
5. extensible local vocabulary on sub-builders.

The sub-builders are small fluent objects. `goalBuilder` exposes `id`, `description`, `priority`, `status`, `asks`, and `tag` (`builders.go:126-134`). Other builders expose methods appropriate to their entity.

### 1.5 The flaw: non-function callbacks are silently ignored

The current optionality goes too far. `applyEntityBuilder` returns nil when the provided callback is not a function (`builders.go:110-113`). That means this likely does not fail:

```js
project("Oops").goal("Goal", { id: "GOAL-001" })
```

The object callback is ignored, so the goal gets no ID. For a playbook-quality DSL, this is a silent failure. The rule should be:

- missing / `undefined` / `null` callback: OK, use defaults;
- present but not a function: error with a precise message.

Codesign already does this better for builder callbacks. Its `applyBuilderCallback` returns `builder callback must be a function` if the value is present but not callable (`codesign/builders.go:365-374`).

### 1.6 Recommended pattern for the playbook

For widgetdsl v2 or any new Goja DSL, use a **two-tier API**:

```js
// Tier 1: simple, opinionated defaults.
data.collection("agenda", rows)
  .schema(agendaSchema)
  .edit()
  .table()
  .toIR()
```

```js
// Tier 2: optional lambda configurators.
data.collection("agenda", rows, c => c
  .schema(agendaSchema)
  .edit(e => e
    .select(data.selection.urlParam("agenda", query.agenda))
    .submit(data.formPost("/settings/agenda-item"))
    .actions(a => a
      .reorder(ui.action.server("admin-reorder-course-agenda"))
      .remove(ui.action.server("admin-delete-agenda-item"))))
  .arrange(a => a.masterDetail()))
  .toIR()
```

The first surface is not “less typed”; it is just defaulted. The second surface adds controlled extension points.

## 2. Typed Go-side IR/specs: where they exist, and where widgetdsl falls short

### 2.1 Researchctl has typed Go-side project IR

`researchctl`'s core representation is a typed Go spec, not a map. The root is:

```go
type ResearchProjectSpec struct {
    SchemaVersion int
    Kind          string
    Name          string
    Description   string
    Plugins       []PluginUseSpec
    Goals         []GoalSpec
    Questions     []QuestionSpec
    Hypotheses    []HypothesisSpec
    WorkPackages  []WorkPackageSpec
    Experiments   []ExperimentSpec
    Sources       []SourceSpec
    Evidence      []EvidenceSpec
    Decisions     []DecisionSpec
    Reports       []ReportSpec
    ReviewRules   []ReviewRuleSpec
    Views         []ViewSpec
    Metadata      JsonObject
}
```

Evidence: `pkg/research/spec/types.go:56-73`.

The builder mutates these typed structs directly. Its terminal returns `*spec.ResearchProjectSpec` (`researchctl/builders.go:13`) and its validation terminal returns `validate.Result` (`builders.go:14`).

Validation has a structured result type:

```go
type Issue struct {
    Severity Severity
    Code     string
    Path     string
    EntityID spec.ID
    Message  string
}

type Result struct {
    Issues []Issue
}
```

Evidence: `pkg/research/validate/result.go:17-27`.

This is a true typed Go-side IR/spec. It is also serializable to YAML/JSON, but the authoring API does not use arbitrary maps as its internal state.

### 2.2 Codesign has typed Go-side run IR/specs

`codesign` is an even better model for the playbook because it combines typed specs, optional fragments, validation, and precise TypeScript. Its typed root is:

```go
type RunSpec struct {
    SchemaVersion int
    Kind          string
    Name          string
    ExperimentID  string
    Backend       string
    Topology      TopologySpec
    Workload      WorkloadSpec
    Policy        PolicySpec
    Metrics       []MetricSpec
}
```

Evidence: `pkg/codesign/spec/types.go:21-31`.

It also defines typed nested specs: `TopologySpec`, `DeviceSpec`, `WorkloadSpec`, `StageSpec`, `PolicySpec`, `MetricSpec`, and `ValidationResult` (`types.go:33-92`). The Goja builder creates a `RunSpec` with defaults (`codesign/builders.go:12-20`), mutates that typed struct through fluent methods, and exposes `.validate()`, `.toSpec()`, and `.run()` terminals (`builders.go:76-80`).

The TypeScript side mirrors this much more precisely than widgetdsl:

- `RunSpec`, `TopologySpec`, `DeviceSpec`, `WorkloadSpec`, etc. (`typescript.go:3-12`).
- `FragmentFn<T>` (`typescript.go:27`).
- `RunSpecBuilder`, `TopologyBuilder`, `WorkloadBuilder`, `MetricsBuilder` with method signatures (`typescript.go:29-32`).
- `RunSpecLike = RunSpec | RunSpecBuilder | { toSpec(): RunSpec }` (`typescript.go:28`).

### 2.3 Widgetdsl has typed TS IR and manifests, but not typed Go authoring IR

`rag-evaluation-system` has strong type sources on the TypeScript/React side:

- `WidgetNode = TextNode | ElementNode | ComponentNode` (`packages/rag-evaluation-site/src/widgets/ir.ts:49`).
- `ComponentNode` has a widget `type`, typed props, and child nodes (`ir.ts:141-146`).
- `ActionSpec` is a discriminated union of navigate/download/server/event/copy actions (`ir.ts:150-155`).
- `WidgetProps` is a large union of all component prop interfaces (`ir.ts:841+`).
- Widget manifests name module, helper, props type, adapter, children, slots, and actions (`internal/widgetmanifest/types.go:10-22`).

So the system is not devoid of type information. The issue is that `pkg/widgetdsl` does not consume it for the JS authoring layer. The generated declarations still say:

```ts
export interface WidgetNode { kind: string; [key: string]: any; }
export interface WidgetAction { kind: string; [key: string]: any; }
export type Props = Record<string, any>;
```

Evidence: `pkg/widgetdsl/typescript.go:20-31`.

The grammar declarations are similarly loose:

```ts
export interface FieldSpec { role: string; [key: string]: any; }
export function record(values: Props, options: Props): WidgetNode;
export function collection(rows: Props[], options: Props): WidgetNode;
```

Evidence: `pkg/widgetdsl/typescript.go:78-100`.

### 2.4 Answer to “do we have typed Go-side IR somewhere?”

Yes, in the DSL ecosystem:

| System | Typed Go-side spec/IR? | Notes |
| --- | --- | --- |
| `researchctl` | Yes | `ResearchProjectSpec` and entity specs; validation result with severity/code/path. |
| `codesign` | Yes | `RunSpec`, topology/workload/policy/metric specs; best builder + DTS model. |
| `geppetto` | Yes, but more as typed Go refs/builders around domain objects | Hidden ref wrappers, generated declarations, parity test. |
| `goja-bleve` | Yes, as typed Go refs around Bleve mapping/query/index objects | Strong runtime handle substrate. |
| `widgetdsl` | Not in Go authoring layer | TS-side IR is typed; Go-side manifests exist; JS DSL emits maps. |

The playbook should distinguish **wire format** from **authoring IR**:

- It is fine for the terminal output to be JSON-like IR.
- It is not fine for the authoring API's internal state to be arbitrary maps.

A good target for widgetdsl v2 is:

```go
type WidgetNodeSpec struct { ... }        // Go mirror of TS WidgetNode

type CollectionSpec struct {
    Name        string
    Rows        []map[string]any
    Schema      *SchemaSpec
    Mode        CollectionMode
    Arrangement ArrangementSpec
    Selection   SelectionSpec
    Actions     []ActionBinding
    Issues      []ValidationIssue
}

func (c *CollectionSpec) Validate() ValidationResult
func (c *CollectionSpec) ToIR() WidgetNodeSpec
```

## 3. Geppetto DTS parity: export-surface synchronization test

### 3.1 The problem it solves

Generated `.d.ts` files can drift from runtime exports. Drift has two directions:

- Runtime has an export that the `.d.ts` does not mention. Authors cannot discover/typecheck it.
- `.d.ts` mentions an export that runtime no longer provides. Authors get compile-time success and runtime failure.

Widgetdsl currently risks this because declarations are hand-built from helper maps and raw strings. Geppetto demonstrates a concrete test to catch at least export-surface drift.

### 3.2 The runtime export surface

Geppetto installs exports in `installExports`:

```go
m.mustSet(exports, "version", "0.1.0")
m.installConsts(exports)
// inferenceProfiles.load/resolve/default
m.mustSet(exports, "engine", m.engineBuilder)
m.mustSet(exports, "embeddings", m.embeddingsBuilder)
m.mustSet(exports, "agent", m.agentBuilder)
m.installTurnStoresNamespace(exports)
m.mustSet(exports, "tool", m.toolBuilder)
m.mustSet(exports, "toolRegistry", m.toolRegistryBuilder)
m.installSchemaNamespace(exports)
```

Evidence: `geppetto/pkg/js/modules/geppetto/module.go:175-190`.

These are the names the runtime actually exposes via `require("geppetto")`.

### 3.3 The generated declaration source

Geppetto has `go:generate` directives:

```go
//go:generate go run ../../../../cmd/tools/gen-meta --schema ../../../spec/geppetto_codegen.yaml --section js-go
//go:generate go run ../../../../cmd/tools/gen-meta --schema ../../../spec/geppetto_codegen.yaml --section js-dts
```

Evidence: `geppetto/pkg/js/modules/geppetto/generate.go:1-4`.

The generated `.d.ts` begins with:

```ts
// Code generated by cmd/gen-meta from geppetto_codegen.yaml. DO NOT EDIT.

declare module "geppetto" {
    export const version: string;
    // ...
}
```

Evidence: `geppetto/pkg/js/modules/geppetto/spec/geppetto.d.ts.tmpl` and `pkg/doc/types/geppetto.d.ts`.

### 3.4 The parity test algorithm

The test is `TestGeneratedDTSMatchesRuntimeExportSurface` (`dts_parity_test.go:24-51`). It does four steps.

#### Step 1: parse exported names from the `.d.ts`

The test defines regexes:

```go
reExportConst    = regexp.MustCompile(`(?m)^\s*export const ([A-Za-z_][A-Za-z0-9_]*)\s*:\s*`)
reExportFunction = regexp.MustCompile(`(?m)^\s*export function ([A-Za-z_][A-Za-z0-9_]*)\s*\(`)
reObjectLevelProperty = regexp.MustCompile(`(?m)^\s{8}([A-Za-z_][A-Za-z0-9_]*)\s*(?:\(|:)`)
```

Evidence: `dts_parity_test.go:18-21`.

`parseDTSSurface` collects top-level `export const` and `export function` names into a set (`dts_parity_test.go:97-106`). For exported object literals, it finds the matching braces and extracts object-level properties (`dts_parity_test.go:108-127`). Helpers `findMatchingBrace`, `objectProperties`, and `sortedSet` make that deterministic (`dts_parity_test.go:135-165`).

#### Step 2: instantiate the runtime

The test creates a JS runtime:

```go
rt := newJSRuntime(t, Options{})
```

Evidence: `dts_parity_test.go:31`.

#### Step 3: inspect runtime exports using `Object.keys`

For top-level exports, it evaluates:

```js
Object.keys(require("geppetto")).sort()
```

For namespaces, it evaluates:

```js
Object.keys(require("geppetto").schema).sort()
```

Evidence: `runtimeObjectKeys` uses `mustEvalExprExport` with `Object.keys(...).sort()` (`dts_parity_test.go:62-65`).

#### Step 4: compare sets

`assertSameSet` sorts both sets, computes missing and extra names, and fails with a precise diff (`dts_parity_test.go:168-190+`). The test checks:

- all top-level `geppetto` exports (`dts_parity_test.go:31-37`),
- `consts`, `inferenceProfiles`, `schema`, and `turnStores` namespace exports (`dts_parity_test.go:39-50`).

### 3.5 What this test catches

It catches export surface drift:

- If runtime adds `foo` but `.d.ts` omits it: `extra: [foo]`.
- If `.d.ts` declares `foo` but runtime no longer exports it: `missing: [foo]`.
- If nested namespace `schema` gets a method mismatch: namespace-specific diff.

### 3.6 What it does not catch

It does **not** fully typecheck TypeScript semantics. It does not know whether:

- `engine().inference(settings)` has the correct parameter type,
- return types are precise,
- overloads are correct,
- generic constraints are correct,
- object properties nested deeper than the regex pattern match.

That is fine. It is a parity smoke test, not a TypeScript compiler test.

The playbook should recommend two layers:

1. **Export-surface parity**: geppetto-style runtime `Object.keys` vs generated declaration names.
2. **Type-level fixture tests**: small `.ts` files compiled with `tsc --noEmit` that use expected-good and expected-bad examples.

### 3.7 How widgetdsl could copy this

For widgetdsl, a parity test could:

1. Generate `.d.ts` from `moduleSpecs`, widget manifests, or a new schema source.
2. Start a Goja runtime and require each module (`ui.dsl`, `data.dsl`, `context_window.dsl`, `course.dsl`, `cms.dsl`).
3. Compare `Object.keys(require("ui.dsl"))` to the module's declared exported functions.
4. Compare nested namespaces: `action`, `cell`, `recipes`, `f`.
5. Fail if declarations and runtime diverge.

Pseudo-test:

```go
func TestWidgetDSLDeclarationsMatchRuntimeExports(t *testing.T) {
    for _, module := range []string{"ui.dsl", "data.dsl", "context_window.dsl", "course.dsl", "cms.dsl"} {
        expected := parseDTSModuleSurface(generatedDTS, module)
        actual := runtimeObjectKeys(vm, fmt.Sprintf(`require(%q)`, module))
        assertSameSet(t, module, expected.TopLevel, actual)
    }

    assertSameSet(t, "data.dsl.f", expected.Grouped["f"], runtimeObjectKeys(vm, `require("data.dsl").f`))
    assertSameSet(t, "data.dsl.cell", expected.Grouped["cell"], runtimeObjectKeys(vm, `require("data.dsl").cell`))
    assertSameSet(t, "ui.dsl.action", expected.Grouped["action"], runtimeObjectKeys(vm, `require("ui.dsl").action`))
}
```

This would immediately guard against forgetting to update declarations when adding helpers or moving generic helpers between modules.

## 4. Go-emrichen tag-operator model

### 4.1 The core idea

Go-emrichen is a YAML templating engine where **YAML tags are operators**. Examples:

```yaml
message: !If
  test: !Var is_admin
  then: !Format "Welcome, {name}!"
  else: !Format "Hello, {name}!"
```

```yaml
ports: !Loop
  over: !Var service.ports
  as: port
  template:
    name: !Format "http-{port}"
    containerPort: !Var port
```

The YAML tree is the AST. A tag on any node names the operator that should transform that node.

### 4.2 Tags are functions from YAML node to YAML node

The central type is:

```go
type TagFunc func(ei *Interpreter, node *yaml.Node) (*yaml.Node, error)
type TagFuncMap map[string]TagFunc
```

Evidence: `go-emrichen/pkg/emrichen/emrichen.go:41-42`.

The interpreter has an environment, additional tag handlers, and template function maps (`emrichen.go:19-22`). Built-in handlers are registered in `defaultHandlers`: `!Defaults`, `!All`, `!Any`, `!Concat`, `!Filter`, `!Group`, `!If`, `!Include`, `!Index`, `!Join`, `!Loop`, `!Lookup`, `!Merge`, `!Not`, `!Op`, `!Var`, `!With`, and others (`emrichen.go:44+`).

Custom tags can be added either through `WithAdditionalTags` (`emrichen.go:219-228`) or `RegisterTag` (`emrichen.go:298-306`). That is the extension mechanism.

### 4.3 The interpreter recursively evaluates the tree

`Process` is the evaluator. It reads `node.Tag`, splits comma-composed tags, normalizes tag names, reverses them, and applies each operator in sequence (`emrichen.go:333-352`). If a handler exists, it calls that handler (`emrichen.go:354-357`). If no handler exists, it recursively processes children based on node kind:

- sequence: process each child and build a new sequence (`emrichen.go:361-377`),
- mapping: process each value and build a new mapping (`emrichen.go:378-397`),
- scalar: return node unchanged (`emrichen.go:398-399`),
- document: process document content (`emrichen.go:402-406`).

This is why tags naturally nest. A tag handler can call `ei.Process` on a subnode; that subnode may itself contain tags; those tags can call more tags.

### 4.4 Strict argument parsing is part of the model

Many tags use `ParseArgs` to define accepted keys, required keys, and which values should be expanded before use. `ParsedVariable` has:

```go
type ParsedVariable struct {
    Name     string
    Expand   bool
    Required bool
}
```

Evidence: `parser.go:8-12`.

`ParseArgs` rejects unknown keys (`parser.go:61-64`), rejects missing required keys (`parser.go:80-85`), and processes values marked `Expand` through `ei.Process` (`parser.go:70-76`).

This is one of the most important lessons for widgetdsl: **operator arguments are schema-checked before evaluation**. The current widgetdsl options bags often silently absorb unknown keys; go-emrichen's tag handlers do not.

### 4.5 “Lambdas” in go-emrichen are template subnodes under lexical scope

Go-emrichen does not accept JavaScript lambdas. It accepts YAML subtrees that are evaluated later under a scoped environment. That is the functional analogue.

#### `!Loop`: template as a scoped lambda body

`!Loop` requires `over` and `template` (`loop.go:13-21`). It expands `over` first (`Expand: true` for the `over` variable), then iterates sequence or mapping nodes. For each item, it creates a local environment (`as`, optional `index_as`, optional `previous_as`) and processes the `template` node under that environment (`loop.go:70-111` for sequences; mapping support continues after line 112).

Conceptually:

```yaml
!Loop
  over: !Var items
  as: item
  template: <body using item>
```

is equivalent to:

```js
items.map(item => evaluate(body, { item }))
```

The YAML `template` node is the lambda body; the environment binding is the lambda parameter.

#### `!If`: branch nodes as lazy subexpressions

`!If` parses `test`, `then`, and `else` (`if.go:5-10`). It processes `test`, checks truthiness, and only processes the chosen branch (`if.go:15-30`). This matters because `then` and `else` can contain arbitrary nested tagged expressions.

#### `!With`: lexical scope plus template body

`!With` requires `vars` and `template` (`with.go:8-16`). It pushes variables, processes the template, then pops the environment (`with.go:22-28`). This is lexical scoping.

### 4.6 Operators can compose in two directions

Go-emrichen supports composition in two ways.

#### Nested composition

Any tag argument or template body can itself contain tagged nodes:

```yaml
result: !If
  test: !Op
    op: ">"
    a: !Var count
    b: 3
  then: !Format "large: {count}"
  else: !Format "small: {count}"
```

Here `!If` contains `!Op`, `!Var`, and `!Format`. Recursion handles the nesting.

#### Comma-tag pipeline composition

`Process` splits tags on commas and reverses/applies them in sequence (`emrichen.go:333-352`). The spec uses forms such as `!Include,Var` to include a file path that comes from a variable. This is a compact operator pipeline on one node.

The playbook can borrow the concept without copying the syntax. In JavaScript builder DSLs, the equivalent is `.use(fragment)` or chained operator builders.

### 4.7 What the tag-operator model contributes to Goja DSL design

The useful lessons are:

1. **Operators are first-class grammar nodes.** Each tag is a named operator with a contract.
2. **Operators validate their own argument shape.** `ParseArgs` rejects unknown and missing keys.
3. **Operators decide which arguments are evaluated eagerly.** `Expand: true` is a per-argument evaluation policy.
4. **Operators can accept template bodies.** `!Loop.template`, `!With.template`, and `!If.then/else` are deferred sub-expressions.
5. **Operators compose by nesting.** The tree shape makes composition natural.
6. **Operators can introduce scoped bindings.** `!Loop` and `!With` are not just macros; they evaluate subtrees under local variables.
7. **Custom operators are registered through a map.** `WithAdditionalTags` and `RegisterTag` are simple extension hooks.

### 4.8 What not to copy blindly

Do not copy these aspects into Goja DSLs without thought:

- YAML tags are terse but not discoverable through TypeScript autocomplete.
- The model is dynamically typed around `yaml.Node`.
- The “lambda” is a data subtree, not a JS function, so it does not directly map to callback ownership rules in Goja.
- Comma-tag composition is clever but can be hard to read.

For Goja, the better translation is:

```ts
type FragmentFn<T> = (builder: T) => void | T
```

plus typed operators/builders:

```js
data.collection(rows, c => c
  .filter(f => f.where(row => row.status === "active"))
  .group(g => g.by("status"))
  .arrange(a => a.table(t => t.columns("title", "status"))))
```

The Emrichen-inspired part is the operator algebra: filter, group, arrange, render, validate, each with a contract and local scopes.

## 5. Synthesis: answer the design questions as playbook rules

### Rule 1: Lambdas optional; simple path required

Every fluent builder that creates a meaningful domain object should support a defaulted, no-lambda form.

Bad:

```js
project("X").goal("G", g => g) // lambda required just to satisfy API
```

Good:

```js
project("X").goal("G")
```

Better:

```js
project("X").goal("G", g => g.id("GOAL-001"))
```

### Rule 2: Optional does not mean silently ignore wrong types

Researchctl's optional callback behavior is good for missing callbacks, but bad for wrong present callbacks. Use codesign's stricter pattern:

```go
if goja.IsUndefined(cb) || goja.IsNull(cb) { return nil }
fn, ok := goja.AssertFunction(cb)
if !ok { return fmt.Errorf("builder callback must be a function") }
```

Evidence: `codesign/builders.go:365-374`.

### Rule 3: Builders mutate typed Go specs, not maps

Use the `researchctl`/`codesign` pattern:

```go
spec := &DomainSpec{Defaults...}
return builder(spec)
```

Then provide terminals:

```js
.validate()  // structured issues
.toSpec()    // typed spec object
.toIR()      // for UI DSLs, serialized Widget IR
.run()       // for executable DSLs
```

### Rule 4: Declaration generation gets a parity test

Every module with generated declarations should have a geppetto-style parity test:

- parse generated `.d.ts` export names,
- instantiate runtime,
- compare `Object.keys(require(module))`,
- compare nested namespaces,
- optionally compile TypeScript fixtures for semantic checks.

### Rule 5: Tag-operator thinking belongs in grammar design

When a DSL starts to grow recipes, reframe them as operators:

- What is the input shape?
- What arguments are accepted?
- Which arguments are evaluated immediately?
- Which arguments are deferred bodies/configurators?
- What local bindings are in scope?
- What typed output does the operator produce?
- How does it compose with other operators?

This is the real transfer from go-emrichen to Goja builder DSLs.

## 6. Intern implementation guide: how to apply this to widgetdsl v2

### Phase 1: tighten existing option decoding

Before inventing new APIs, stop silent failures:

- Decode `record` and `collection` options into typed Go structs.
- Reject unknown keys.
- Validate enum values (`verb`, `arrange`, `role`) with suggestions.
- Return Go errors instead of panicking where possible.

This is inspired by go-emrichen `ParseArgs`, which rejects unknown keys and required-key omissions.

### Phase 2: introduce typed marker handles

Replace map markers with typed handles:

- `SchemaRef` for schemas,
- `UrlParamSelectionRef` for URL selection,
- `FormPostRef` for submit binding,
- `ActionRef` for actions.

Use the hidden-ref substrate from goja-bleve/geppetto or a shared `fluent` package. Passing `urlParam` where `formPost` is expected should be a typed error, not a silent no-op.

### Phase 3: define typed Go-side Widget intent specs

Do not start by mirroring every React prop. Start with high-level intent specs:

```go
type SchemaSpec struct { Fields []FieldSpec }
type FieldSpec struct { Name string; Role FieldRole; Editor EditorSpec; Summary SummarySpec }
type CollectionSpec struct { Rows []map[string]any; Schema *SchemaSpec; Mode Mode; Arrangement ArrangementSpec }
type SectionSpec struct { Title string; Children []WidgetIntent }
```

Then compile those to current Widget IR maps at the terminal.

### Phase 4: add simple methods first, optional lambdas second

Simple:

```js
data.collection("agenda", rows).schema(agendaSchema).edit().masterDetail().toIR()
```

Configurable:

```js
data.collection("agenda", rows, c => c
  .schema(agendaSchema)
  .edit(e => e.selectUrl("agenda", query.agenda).submitPost("/settings/agenda-item"))
  .arrange(a => a.masterDetail(md => md.summaryTable().detailForm()))
)
```

Fragments:

```js
const destructiveRows = a => a
  .reorder(ui.action.server("admin-reorder-course-agenda"))
  .remove(ui.action.server("admin-delete-agenda-item"));

data.collection("agenda", rows).edit(e => e.actions(destructiveRows))
```

### Phase 5: declaration parity and type fixtures

- Generate declarations from the same source as runtime exports.
- Add geppetto-style export parity tests.
- Add TypeScript examples that should compile.
- Add negative examples with `// @ts-expect-error` for invalid roles, invalid arrangements, wrong callback types, and wrong marker types.

## 7. Decision records

### Decision: Lambdas as optional refinement, not mandatory structure

- **Context:** Researchctl proves optional lambdas work; codesign proves lambda fragments work for complex specs.
- **Options considered:** Require lambdas everywhere; forbid lambdas for simplicity; make lambdas optional configurators.
- **Decision:** Make lambdas optional configurators.
- **Rationale:** Simple authoring remains approachable; advanced composition remains possible.
- **Consequences:** Builders need good defaults and strict callback type checks.
- **Status:** proposed.

### Decision: Typed specs before serialized IR

- **Context:** Widgetdsl currently emits maps; researchctl/codesign mutate typed Go specs.
- **Options considered:** Keep map IR as authoring API; create typed Go specs and compile to IR only at terminals.
- **Decision:** Use typed Go specs for authoring state.
- **Rationale:** Enables validation, typed markers, precise DTS, and better errors.
- **Consequences:** More Go code, but much safer DSL surface.
- **Status:** proposed.

### Decision: DTS parity required for generated declarations

- **Context:** Geppetto already catches runtime/declaration drift.
- **Options considered:** Trust generator; golden-file compare only; runtime export parity plus type fixtures.
- **Decision:** Use runtime export parity plus type fixtures.
- **Rationale:** Golden files catch text drift; parity catches runtime/API drift; fixtures catch semantic type drift.
- **Consequences:** Tests become part of every DSL module's API discipline.
- **Status:** proposed.

### Decision: Borrow go-emrichen's operator architecture, not its YAML syntax

- **Context:** The playbook needs composable grammar operators; go-emrichen demonstrates nested operators and scoped template bodies.
- **Options considered:** Copy YAML tags; ignore emrichen; translate its operator ideas into typed Goja builders.
- **Decision:** Translate operator ideas into typed Goja builders.
- **Rationale:** JavaScript authors need TypeScript autocomplete and Goja runtime safety, not YAML tag syntax.
- **Consequences:** Operators should have strict argument schemas, scoped configurators, and typed outputs.
- **Status:** proposed.

## 8. Reference map for interns

### Researchctl optional lambdas

- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/researchctl/builders.go:15-22` — `goal(title, cb...)` with default spec.
- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/researchctl/builders.go:106-115` — optional callback application.
- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/researchctl/module.go:23-35` — TS declaration with optional `build?:` callbacks.
- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/researchctl/module_test.go:24-35` — runnable JS builder example.

### Typed Go-side specs

- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/research/spec/types.go:56-90` — `ResearchProjectSpec` and `GoalSpec`.
- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/research/validate/result.go:17-27` — structured validation result.
- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/codesign/spec/types.go:21-92` — `RunSpec` and validation result.
- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/builders.go:12-80` — builder over typed `RunSpec`.
- `/home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/typescript.go:27-33` — `FragmentFn<T>` and builder interfaces.

### Widgetdsl type sources and gaps

- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/ir.ts:49-155` — typed TS Widget IR and actions.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/ir.ts:841+` — `WidgetProps` union.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/internal/widgetmanifest/types.go:10-22` — widget manifest metadata.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript.go:20-31` — current weak declarations.
- `/home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/typescript.go:78-100` — weak grammar declarations.

### Geppetto DTS parity

- `/home/manuel/code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/dts_parity_test.go:24-51` — test entrypoint.
- `/home/manuel/code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/dts_parity_test.go:97-132` — `.d.ts` surface parser.
- `/home/manuel/code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/dts_parity_test.go:168-190` — set comparison and failure reporting.
- `/home/manuel/code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/module.go:175-190` — runtime export installation.
- `/home/manuel/code/wesen/go-go-golems/geppetto/pkg/js/modules/geppetto/generate.go:1-4` — generation directives.

### Go-emrichen tag operators

- `/home/manuel/code/wesen/go-go-golems/go-emrichen/pkg/emrichen/emrichen.go:19-42` — interpreter fields and `TagFunc` type.
- `/home/manuel/code/wesen/go-go-golems/go-emrichen/pkg/emrichen/emrichen.go:44+` — default tag handler map.
- `/home/manuel/code/wesen/go-go-golems/go-emrichen/pkg/emrichen/emrichen.go:219-228` — `WithAdditionalTags`.
- `/home/manuel/code/wesen/go-go-golems/go-emrichen/pkg/emrichen/emrichen.go:298-306` — `RegisterTag`.
- `/home/manuel/code/wesen/go-go-golems/go-emrichen/pkg/emrichen/emrichen.go:333-420` — recursive `Process` evaluator.
- `/home/manuel/code/wesen/go-go-golems/go-emrichen/pkg/emrichen/parser.go:8-88` — `ParseArgs` strict argument parsing.
- `/home/manuel/code/wesen/go-go-golems/go-emrichen/pkg/emrichen/loop.go:8-111` — `!Loop` as scoped template-body evaluation.
- `/home/manuel/code/wesen/go-go-golems/go-emrichen/pkg/emrichen/if.go:5-30` — `!If` branch evaluation.
- `/home/manuel/code/wesen/go-go-golems/go-emrichen/pkg/emrichen/with.go:8-29` — `!With` lexical scope.
