# Changelog

## 2026-07-06

- Initial workspace created


## 2026-07-06

Reviewed the full widget library (atoms/foundation, molecules, organisms, Widget IR/renderer, Go/Goja pkg/widgetdsl) via five read-only agents and wrote an intern-facing analysis/design doc: the engine+contract+preset lens, three recurring smells (duplicated engines, ad-hoc spec shapes, mirror drift), per-layer decomposition catalog, cross-cutting IR spec unification (AccessorSpec/SelectionSpec/ListItemSpec/ctx helpers), manifest-as-source-of-truth codegen, DSL parity+elegance, and an A/B/C roadmap. Uploaded to reMarkable. No code changed (review-only).


## 2026-07-06

Rewrote the architecture material as a Norvig-style intern textbook (reference/02): 8 chapters teaching the widget system from foundations — three forms of a UI, the IR tree + renderer, adapters + RenderContext, defunctionalized specs (how to send a function through JSON), the engine+contract+preset pattern (MatrixGrid trace), the Go/Goja DSL, manifests + drift, and a lens for reading the library critically. Uploaded to reMarkable as v2.


## 2026-07-07

Expanded the analysis/design doc (design-doc/01) into the textbook style: added a Part 0 glossary defining every recurring term (IR, tier, engine/contract/preset, seam, smell, drift, mirror, hoist, passthrough adapter, codegen, behavior-preserving, opinionated, leverage-to-risk, escape hatch, lowering, parity), and rewrote every Part from terse shorthand into developed prose that explains each finding and defines jargon on first use. Uploaded to reMarkable as v2 (expanded).


## 2026-07-07

Added design-doc/02: a textbook-style redesign of the Go/Goja Widget DSL toward a composition-first, opinionated JS API. Diagnoses the five coexisting authoring styles and their friction, then proposes engine verbs (matrix/board/timeline/record/list/calendar) each taking a data set + a unit-renderer function the runtime calls and serializes (generalizing the existing detail(row) pattern at module.go:882) — exposing the engine/contract/preset pattern as the DSL. Keeps a named raw() escape hatch, unifies the spec builders (at/style.by/cell.cycle-value/field/select), and lays out a non-breaking backwards migration + manifest-driven generation. Uploaded to reMarkable next to design-doc 01 and v2.


## 2026-07-07

Added design-doc/03: a by-example tour of every DSL (ui/data/data.v2/context_window/course/cms + proposed engine verbs & crm), 21 escalating TypeScript examples each paired with its types and Go counterpart, with a running 'map-watch' flagging every untyped Record<string,any>/map[string]any. Culminates in a proposal for branded/opaque TS spec+node types (builder-only, distinct per kind) and Go sealed builder returns + a goja boundary brand (generalizing v2's __widgetdsl_v2_ref) so bad API usage fails to compile or is rejected at the boundary. Uploaded to reMarkable next to design-docs 01/01v2/02.

