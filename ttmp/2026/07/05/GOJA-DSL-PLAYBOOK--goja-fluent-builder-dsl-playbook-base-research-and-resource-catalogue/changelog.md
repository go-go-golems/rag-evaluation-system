# Changelog

## 2026-07-05

- Initial workspace created


## 2026-07-05

Created base research ticket: DSL catalogue (design-doc), research logbook (reference), investigation diary. Inventoried 12 Goja DSLs + 3 cross-cutting pieces across goja-bleve, goja-dbus, go-minitrace, rag-evaluation-system widgetdsl, go-go-goja, goja-text, goja-git, goja-github-actions, goja-treesitter, glazed, go-emrichen. Classified 5 implementation patterns; identified goja-bleve typed-ref substrate as the runtime-typecheck model and goja-dbus as the composable-grammar model.

### Related Files

- /home/manuel/code/wesen/go-go-golems/goja-bleve/pkg/api_types.go — Typed-ref substrate to extract into a shared fluent package
- /home/manuel/code/wesen/go-go-golems/goja-dbus/pkg/dbusgoja/builders.go — Cleanest composable grammar reference


## 2026-07-05

Validated (docmgr doctor passes cleanly after adding fluent-builder + typescript vocabulary). Uploaded bundle 'Goja DSL Playbook — base research' (3 docs, toc-depth 2) to reMarkable at /ai/2026/07/05/GOJA-DSL-PLAYBOOK and verified listing. Ticket ready for senior-researcher handoff.


## 2026-07-05

Design-doc 02 added: insider self-assessment of the widgetdsl grammar. Five silent-failure modes empirically verified against the built binary (typo'd arrange/options/verbs, wrong markers, out-of-range enums all absorbed without error); root cause named as IR-as-API vs IR-as-output; v2 sketch keeps the grammar language on typed Go specs with configurator builders, strict option decoding, and error-panel rendering; eleven proposed catalogue additions including a silent-ignore failure taxonomy, a tree-shaped-data axis, dormant type sources (ir.ts + 80 widget manifests, zero consumers), and agent-authorship as a design constraint.


## 2026-07-05

Extended catalogue with 4 new DSLs: geppetto (typed-ref + clone-on-each-step + DTS parity test), discord-bot (defineBot + Proxy-trap ui builders), researchctl (lambda-configurator project builder), codesign (strongest model: runSpec + .use() fragment composition + lambda configurators + precise TS interfaces via TypeScriptDeclarer). Added Patterns A-prime, F, G. Gap analysis updated: composable-grammar and lambdas are now realised; remaining gap is extraction + standardisation. Related 6 new key files.

### Related Files

- /home/manuel/workspaces/2026-06-30/benchmark-cpu-inference/researchctl/pkg/gojamodules/codesign/builders.go — Strongest composable-grammar model - use fragment composition


## 2026-07-05

Added widgetdsl grammar self-assessment report: the RAGEVAL-UI-GRAMMAR work is valuable as UI/vocabulary evidence but remains Pattern C; recommends typed builders, validation terminals, precise TypeScript, and lambda-configurable marks/arrangements.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/02-self-assessment-of-the-widgetdsl-grammar-what-pattern-c-actually-costs-and-what-the-playbook-should-add.md — New assessment report


## 2026-07-05

Added intern-facing deep dive on optional researchctl lambdas, typed Go-side IR/specs, geppetto DTS parity testing, and go-emrichen tag-operator composition; uploaded as a separate reMarkable document.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/04-goja-dsl-deep-dive-optional-lambdas-typed-ir-dts-parity-and-tag-operators.md — New deep-dive document


## 2026-07-05

Design-doc 05 written (into the placeholder): the v2 overhaul design for all five DSLs, synthesizing docs 01-04 — typed intent specs + two-tier builders (defaulted chains, optional lambda configurators, .use fragments), strict decoding, accumulated validation rendered as ValidationIssues nodes, marks contract shrinking domain modules to schemas+marks, callback action kind making JS lambdas safe across the renderer/runtime boundary (deterministic ids, fail-closed declaration-site hashes, 409/rebuild/retry staleness protocol), manifest-driven codegen with DTS parity, and a 5-phase file-level implementation plan.


## 2026-07-05

Added full rag-evaluation-system DSL overhaul guide covering all Widget DSL modules, typed intent IR, renderer/runtime action serialization, API alternatives, examples, migration phases, and testing strategy.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/05-rag-evaluation-system-dsl-overhaul-design-and-implementation-guide.md — New design and implementation guide


## 2026-07-05

Revised the rag-evaluation-system DSL overhaul guide for a no-backwards-compatibility hard cutover: v1 APIs are now vocabulary evidence, compatibility facades are rejected, and the implementation plan targets clean v2 specs/builders, Action IR v2, widget.unsafe, page rewrites, and removal of old public exports.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/05-rag-evaluation-system-dsl-overhaul-design-and-implementation-guide.md — Hard-cutover revision


## 2026-07-05

Added the Widget DSL event-timeline companion document and hard-cutover phase/task tracker, covering simple table, selectable table, master-detail editor, HTTP requests, frontend execution, backend handlers, and precise implementation phases.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/05-rag-evaluation-system-dsl-overhaul-design-and-implementation-guide.md — Parent design guide being operationalized
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/06-widget-dsl-event-timelines-and-cutover-task-plan.md — New companion operational spec and task tracker


## 2026-07-05

Inventoried current live pages, Storybook examples, missing v2 demos, and deprecated-example policy for the Widget DSL cutover companion document.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/06-widget-dsl-event-timelines-and-cutover-task-plan.md — Demo and example inventory


## 2026-07-05

Recorded baseline validation commands for the Widget DSL cutover: widgetdsl Go tests, rag-evaluation-site typecheck/build, and docmgr doctor all pass before v2 implementation.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/06-widget-dsl-event-timelines-and-cutover-task-plan.md — Baseline validation section


## 2026-07-05

Started P1 typed v2 spec implementation by adding pkg/widgetdsl/v2/spec with PageSpec, NodeSpec, SchemaSpec, FieldSpec, CollectionSpec, ActionSpec, TemplateSpec, and ValidationIssue skeletons.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v2/spec/types.go — Typed v2 spec model
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/ttmp/2026/07/05/GOJA-DSL-PLAYBOOK--goja-fluent-builder-dsl-playbook-base-research-and-resource-catalogue/design-doc/06-widget-dsl-event-timelines-and-cutover-task-plan.md — Task tracker updated for P1.1


## 2026-07-05

Implemented initial P1.2 validation rules for the v2 Widget DSL spec model, including schema/field uniqueness, collection mode/arrangement, URL selection, action shape, payload fields, template parts, and diagnostic helpers.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v2/spec/validate.go — Validation implementation


## 2026-07-05

Implemented initial P1.3 lowering from typed v2 specs to current Widget IR shapes, covering pages, nodes, sections, simple/selectable tables, master-detail form trees, and serializable actions.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v2/spec/lower.go — Typed spec to Widget IR lowering


## 2026-07-05

Added P1.4 v2 spec tests for simplest table, selectable table, master-detail editor lowering, and invalid arrangement validation.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v2/spec/lower_test.go — Spec validation/lowering tests


## 2026-07-05

Implemented P2.1-P2.3 initial data.v2.dsl Goja builder substrate: hidden typed refs, strict callback errors, field/schema builders, simple/selectable table collection builders, and runtime tests.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v2_builders.go — V2 Goja builder implementation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v2_builders_test.go — V2 builder tests


## 2026-07-05

Completed P2.4-P2.5 by adding data.v2.dsl master-detail editor builders with edit/selectUrl/submitPost/create/actions/masterDetail and runtime tests for the full simple/selectable/master-detail authoring examples.

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v2_builders.go — Master-detail editor builder implementation
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/pkg/widgetdsl/v2_builders_test.go — Master-detail Goja runtime test


## 2026-07-05

Added live go-go-course demo pages for P4.1-P4.3: simplest table, selectable table, and master-detail editor using data.v2.dsl (go-go-course commit f82f20a99780902cb776e022ad6d1a3b3c2ee9a7).

### Related Files

- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/lib/course-pages.js — Routes demo page IDs to demo builder
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/lib/pages/dsl-examples.js — Live v2 DSL demo pages
- /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/go-go-course/cmd/go-go-course/server.js — Requires data.v2.dsl and adds safe demo form redirect route

