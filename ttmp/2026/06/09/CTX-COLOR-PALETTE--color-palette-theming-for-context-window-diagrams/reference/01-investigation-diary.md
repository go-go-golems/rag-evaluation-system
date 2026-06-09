---
title: "Investigation Diary"
doc_type: reference
status: active
intent: long-term
topics: [frontend, theming, context-window, visualization]
ticket: CTX-COLOR-PALETTE
created: "2026-06-09"
owners: []
---

# Investigation Diary

## Goal

Investigate the context window visualization subsystem to design a color palette theming layer, and produce a detailed implementation guide for a new intern.

## Step 1: Evidence Gathering and Architecture Mapping

### Prompt Context

**User prompt (verbatim):** "Create a new ticket in 2026-05-27--rag-evaluation-system/ttmp/ (using `docmgr --root 2026-05-27--rag-evaluation-system/ttmp ...`) to add color palette theming support to the context_window part of @2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/ so that we can have more legible visualization of context window instead of just blue and different halftone. ~/Downloads/macos1_subtle_color_palettes.json I like /tmp/pi-clipboard-7a252d34-c9d4-4c4d-a7f7-ecd7f2f932b0.png signal orange, slate coral, cobalt / send, dusty magenta the most because the different colors are complementary so that we can nicely codify the diagrams. Create a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet points and pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and then upload to remarkable. Import the json into the sources/ folder of the ticket"

**Assistant interpretation:** Create a docmgr ticket, do deep analysis of the context window visualization code, write an intern-ready design doc with full prose explanations and pseudocode, store everything in the ticket, and upload to reMarkable.

**Inferred user intent:** The user wants a comprehensive, self-contained design document that a new team member could pick up and implement without needing to ask questions. The document should explain what context windows are, how the visualization system works, and exactly what code to change — with concrete color values from the preferred palettes.

### What I did

- Read every source file in the context window visualization subsystem:
  - `src/context/types.ts` — 15 `ContextPartKind` values, 3 diagram styles
  - `src/context/kinds.ts` — Visual spec system with 3 registries (pattern/tone/outline)
  - `src/context/fixtures.ts` — Example snapshot data (4 snapshots)
  - `src/components/molecules/ContextStripDiagram/` — Horizontal strip view
  - `src/components/molecules/ContextStackDiagram/` — Vertical layer view
  - `src/components/molecules/ContextBudgetBar/` — Budget bar with headroom
  - `src/components/molecules/ContextTreemap/` — Area-proportional treemap
  - `src/components/molecules/ContextLegend/` — Kind swatch legend
  - `src/components/atoms/ContextKindSwatch/` — Individual swatch atom
  - `src/components/organisms/ContextDiagramPanel/` — View orchestrator
  - `src/theme.css` — CSS custom properties
  - All 4 CSS modules — Pattern-based `.kind_*` rules
- Read the `macos1_subtle_color_palettes.json` (16 palettes, 225 lines)
- Analyzed the screenshot showing the palette sheet
- Identified the 3-tier color mapping strategy (accent / tint / neutral)
- Pre-computed tint values for all 4 preferred palettes

### Why

The user needs to understand every layer of the system to design a theming approach that fits cleanly into the existing architecture. The `kinds.ts` visual spec system already has a mode-switching mechanism (`pattern` / `tone` / `outline`), making `'color'` a natural fourth mode.

### What worked

- The existing `getContextKindSpec(kind, mode)` function is perfectly designed for extension — just add a new `colorOverrides` registry and a `mode === 'color'` branch.
- The `data-mode` attribute pattern already on all diagram roots means CSS can use attribute selectors for mode-specific rules.
- CSS custom properties enable palette switching without React changes.

### What didn't work

- Initially tried to use inline styles for fills, but realized this breaks palette switching (would need React re-renders). CSS custom properties are the right approach.
- The pre-computed hex values for tints are tedious to calculate by hand. Used `color-mix()` mental model and verified with hex arithmetic.

### What I learned

- The system has exactly 4 diagram molecules, 1 atom, 1 organism, and 1 spec registry — a small, well-scoped surface area.
- The 15 `ContextPartKind` values can be reduced to 3 visual tiers based on semantic importance.
- The Dusty Magenta / Blue palette naturally maps to the human-vs-machine duality in context windows.

### What was tricky to build

- **Color-to-kind mapping:** 15 kinds × 2-3 accent colors means most kinds get tints, not pure accents. The three-tier system (accent / tint / neutral) was the key insight that makes this work.
- **Tint computation:** Getting the right opacity level for each tier so that kinds are distinguishable but not garish. 30% for retrieval, 20% for tool, 25% for result — these differences matter visually.
- **Colored halftone vs flat fills:** Initially designed color mode as flat solid fills (replacing patterns entirely). User requested keeping the halftone patterns and coloring the pattern lines instead. This required splitting the CSS tokens into `--rag-context-line-<kind>` (pattern marks) and `--rag-context-fill-<kind>` (background behind pattern), and re-creating each CSS pattern with the colored tokens.
- **Provider pattern for CSS vars:** How to switch palettes without React re-renders. Solution: `ContextPaletteProvider` sets CSS custom properties as inline `style` on a wrapper `<div>`. Descendants consume via `var()`. No React context needed — pure CSS cascade.

### What warrants a second pair of eyes

- The WCAG contrast ratios for label-on-fill in color mode — verify with an automated checker.
- The `color-mix()` fallback strategy — confirm that all target browsers support it or that the pre-computed values work as a fallback.

### What should be done in the future

- Phase 2: Palette selector component
- Phase 3: Dark mode palettes
- Chromatic snapshot tests for color mode

### Code review instructions

- Start with `design-doc/01-color-palette-theming-design-implementation-guide.md`
- Verify the three-tier mapping makes sense for the user's mental model
- Check the pre-computed hex values in the palette swatch tables
- Walk through the pseudocode in Steps 1–6 and verify it matches the existing code patterns

### Technical details

- Total files to modify: 9 (types, kinds, theme, 4 CSS modules, swatch CSS, panel TSX)
- Total estimated lines: ~430
- No new dependencies needed
- No changes to data flow, widget registry, or layout components

## Step 2: Reframe Around Configurable Legends and Style Sets

The second design pass changed the core abstraction. Instead of making the existing hardcoded `ContextPartKind` list prettier, the revised design treats kinds as only a default preset and introduces a caller-defined `ContextStyleSet` for custom legend labels, flexible palette size, and per-label visual styles.

This was prompted by the requirement that a caller should be able to say: "I have three labels for this context diagram, here is the thing to render, and here is what style each label has." That requirement means the legend, segment style, and data semantics must be separated.

### Prompt Context

**User prompt (verbatim):** "Ok, so we actually also want to make the legend labels and even the size of the palette flexible and configurable. So basically we'd say \" i have 3 labels for this context diagram, and here's the thing to render, and here's what style each has\", which I think has some pretty interesting impact on the design. 

So far, this has been done by your little brother, so I want you to take a critical look at it, especially in light of the above instructions to redo how legends for context window widgets work, and then make a new reviewed updated version (in a separate document) that points out what is good, what is problematic, what should be changed, how it should be changed, giving proper explanations so that they can do better next time. Upload to remarkable as well, write in a similar style."

**Assistant interpretation:** Create a separate critical review document that preserves the good parts of the previous design but revises the architecture around caller-configurable legend labels, arbitrary palette size, and per-label styles.

**Inferred user intent:** The user wants the implementation guide to teach a better abstraction: data segments should refer to caller-defined style keys, legends should not be hardcoded to `ContextPartKind`, and Storybook should prove that the same layout can render with custom labels and style sets.

### What I did

- Inspected the Widget IR contract in `src/widgets/ir.ts`.
- Confirmed that `ContextLegendWidgetProps`, `ContextDiagramPanelWidgetProps`, and individual diagram widget props are currently tied to `ContextPartKind` and `ContextDiagramStyle`.
- Added `design-doc/02-reviewed-design-configurable-legends-and-style-palettes.md`.
- Critiqued the previous design constructively:
  - good: colored halftone, CSS variables, Storybook showcase direction;
  - problematic: `kind` does too much, legends derive from hardcoded enums, CSS is hardcoded as `.kind_*`, palette size is assumed to match built-in vocabulary.
- Proposed `styleKey`, `ContextVisualStyle`, `ContextLegendItemSpec`, and `ContextStyleSet`.
- Proposed replacing `.kind_*` CSS with generic `.pattern_*` classes and per-segment CSS variables.
- Proposed Widget IR changes so server-sent widgets can carry arbitrary label/style definitions.

### Why

The new requirement invalidates a purely hardcoded-kind color palette design. A custom three-label diagram should not require editing global TypeScript unions, CSS classes, or visual registries. The component library should provide a generic visual grammar and treat built-in context kinds as one preset.

### What worked

- The existing Widget IR is already JSON-like, so `ContextStyleSet` can be made JSON-compatible and carried through widget props.
- The diagram components already centralize rendering around snapshot parts, so adding `styleKey` resolution is conceptually straightforward.
- The existing pattern CSS can be preserved if rewritten as generic `.pattern_*` renderers.

### What didn't work

- The previous `ContextPaletteProvider` idea alone does not solve the legend problem. It can scope theme tokens, but it cannot express "these are my three legend labels and styles" unless paired with a per-diagram `styleSet`.

### What I learned

- `ContextPartKind` is currently a semantic type, visual style key, legend key, and CSS class suffix all at once.
- The better API is `snapshot.parts[].styleKey` plus `styleSet.styles[styleKey]` and `styleSet.legend`.
- Storybook should include a custom diagram that uses no built-in `kind` fields at all; that is the strongest test of the new design.

### What was tricky to build

- Distinguishing global theme switching from per-diagram style configuration. A provider helps global CSS tokens cascade, but the required product behavior is local: each diagram can define its own labels and style count.
- Keeping backward compatibility without preserving the bad abstraction as the primary API. The solution is to keep `ContextPartKind` as a built-in style preset and a legacy path, while making `items + styles` the primary legend API.

### What warrants a second pair of eyes

- The exact naming of `styleKey` vs `legendKey` vs `categoryId`. `styleKey` is implementation-friendly, but `legendKey` might better express caller intent.
- Whether `ContextStyleSet.legend.label` should be plain string or `RenderableValue` in Widget IR. The doc recommends Widget IR support for richer labels, but React component props can start with string.

### What should be done in the future

- Implement the new generic style layer before adding more palette presets.
- Add Storybook tests for a three-label diagram with no `kind` fields.
- Add Widget IR examples with `styleSet` embedded in JSON.

### Code review instructions

- Start with `design-doc/02-reviewed-design-configurable-legends-and-style-palettes.md`.
- Review `src/widgets/ir.ts` first, especially the context diagram widget prop interfaces.
- Then review `ContextLegend.tsx`, `ContextKindSwatch.tsx`, and the four diagram molecules to plan the migration away from `.kind_*` CSS.

### Technical details

- New primary contract: `ContextStyleSet`.
- New part lookup field: `ContextWindowPart.styleKey`.
- New rendering primitive: generic `.pattern_*` CSS classes plus `--ctx-fill`, `--ctx-line`, `--ctx-stroke`, and `--ctx-label` variables.
- Legacy path: `ContextPartKind` maps into a built-in default `ContextStyleSet`.

## Step 3: Hard-Cutover Update — No Backwards Compatibility

This update removes the migration-wrapper framing from the reviewed design. The implementation should not preserve `ContextPartKind`, `ContextKindSwatch`, `legendKinds`, `legendMode`, `.kind_*` CSS, or fallback-to-kind behavior. The only current consumers are `2026-05-27--rag-evaluation-system/` and `ClubMedMeetup/`, so both should be updated to the new `styleKey` + `ContextStyleSet` API in the same branch.

### Prompt Context

**User prompt (verbatim):** "we don't need backwards compatibility, we only have 2026-05-27--rag-evaluation-system/ and ClubMedMeetup/ using the system for now, and those can be updated. So no backwards wrappers or so, just hardcutover. Update the doc and create detailed tasks in the ticket."

**Assistant interpretation:** Revise the reviewed design so it explicitly requires a hard cutover instead of preserving legacy wrappers or migration paths, and replace the ticket task list with detailed implementation tasks.

**Inferred user intent:** The user wants a cleaner implementation with one API, not a dual-mode compatibility layer. They want the docs and tasks to prevent the implementer from adding shims that preserve the wrong abstraction.

### What I did

- Updated `design-doc/02-reviewed-design-configurable-legends-and-style-palettes.md` to add a hard-cutover directive.
- Replaced optional/compatibility language with required `styleKey: string` and required `styleSet: ContextStyleSet`.
- Removed recommendations to preserve `ContextKindSwatch`, `legendKinds`, `legendMode`, diagram `mode`, `.kind_*` CSS, and fallback-to-old-kind resolution.
- Replaced `tasks.md` with a detailed implementation checklist covering architecture, style model, CSS, diagram components, legend/swatch, Widget IR, Storybook, and validation.

### Why

Maintaining compatibility would keep the wrong abstraction alive. Since the only consumers are known and can be updated, the cleanest implementation is a single hard-cutover to caller-defined style sets.

### What worked

- The reviewed design already had the right destination (`styleKey` + `ContextStyleSet`); the update mainly removed transition scaffolding.
- The task list now directly encodes forbidden old APIs and validation searches.

### What didn't work

- The previous reviewed document still contained several legacy/migration phrases (`kind` fallback, `ContextKindSwatch` wrapper, `legendKinds`, and `.kind_*` migration layer). Those had to be removed to avoid implementer confusion.

### What I learned

- The acceptance criteria need to state negative requirements explicitly: no wrappers, no legacy props, no `.kind_*` selectors, no fallback resolution.

### What was tricky to build

- Keeping `fallbackStyle` language without implying backwards compatibility. The doc now frames it as production resilience only, not as a way to silently support old `kind` data.

### What warrants a second pair of eyes

- Whether any non-diagram components still need `ContextPartKind` for annotations/badges. If yes, either cut those over too or separate their domain model from diagram styling explicitly.

### What should be done in the future

- During implementation, run a repo-wide search for `ContextPartKind`, `ContextKindSwatch`, `legendKinds`, `legendMode`, `.kind_`, and `kind:` before review.

### Code review instructions

- Start with `tasks.md`; it is now the authoritative implementation checklist.
- Then review `design-doc/02-reviewed-design-configurable-legends-and-style-palettes.md`, especially the hard-cutover directive and acceptance criteria.

### Technical details

- Updated document: `design-doc/02-reviewed-design-configurable-legends-and-style-palettes.md`
- Updated task list: `tasks.md`

## Step 4: Frontend Core StyleSet Cutover

This implementation slice replaced the frontend context diagram rendering path from hardcoded `kind` classes to required `styleKey` plus `ContextStyleSet`. It also replaced `ContextKindSwatch` with `ContextStyleSwatch` and updated the molecule/organism Storybook stories so reviewers can change styles live.

### Prompt Context

**User prompt (verbatim):** "Now implement these tasks one by one, commit at appropriate intervals, keep a frequent detailed diary as you work, add lots of testing stories in the storybook and the IR Widget transformer, and carry things all the way through to the goja DSL as well."

**Assistant interpretation:** Begin the implementation, commit in focused chunks, and keep a detailed diary. Start with the frontend hard-cutover before moving to Widget IR and Goja DSL.

**Inferred user intent:** The user wants the design turned into working code across React components, Storybook, Widget IR, and Goja DSL, with reviewable commits and progress notes.

**Commit (code):** pending — frontend core cutover

### What I did

- Rewrote `src/context/types.ts` to remove `ContextPartKind` from context diagram parts and require `styleKey`.
- Added `src/context/styles.ts` with `ContextVisualStyle`, `ContextStyleSet`, palette definitions, style resolution, CSS var conversion, and palette-derived style-set helpers.
- Converted fixtures in `src/context/fixtures.ts` from `kind` to `styleKey` and exported default/preferred style sets.
- Deleted `ContextKindSwatch` and created `ContextStyleSwatch`.
- Replaced the four diagram molecule renderers to require `styleSet`, resolve each part by `styleKey`, and render generic pattern classes with CSS variables.
- Replaced `.kind_*` CSS with `.pattern_*` CSS in strip, stack, budget, treemap, and swatch styles.
- Replaced `ContextLegend` with required `items + styles` props.
- Updated molecule Storybook stories, including live style switchers for `ContextStripDiagram` and `ContextDiagramPanel`.
- Updated `CourseSlidePanel`, `FigureBlock` stories, `SlideShell` stories, and `CourseStudioShell` stories to pass style sets.
- Ran `pnpm typecheck`; initial run found expected old API references; after fixes, typecheck passes.

### Why

This is the foundation needed before Widget IR and Goja DSL can be updated. The React components must have a stable hard-cutover API first; otherwise the downstream DSL work would target a moving contract.

### What worked

- The generic pattern CSS approach worked well: all diagrams now use the same style variables (`--ctx-fill`, `--ctx-line`, `--ctx-stroke`, `--ctx-label`) and pattern class names.
- TypeScript caught many old references immediately, especially Storybook stories and annotation components.
- The live Storybook style switcher is a useful acceptance test because the same snapshot can be rendered under multiple style sets without changing data.

### What didn't work

- A simple mechanical `kind` → `styleKey` replacement was not enough. Components such as `AnnotationBadge`, `AnnotationNoteCard`, `CourseSlidePanel`, and `FigureBlock` needed explicit style-set or visual-style handling.
- Two story files referenced `contextDefaultStyleSet` without importing it after automated replacements; typecheck caught this.

### What I learned

- The hard-cutover removes a lot of branching: after `ContextStyleSet` is required, components no longer need `mode`, `legendMode`, `legendKinds`, or kind lookups.
- Molecule stories are important because they catch lower-level missing props before organism stories do.

### What was tricky to build

- Keeping annotation UI compiling without reintroducing `ContextPartKind`. I updated it to use `ContextVisualStyle` directly, with transcript note cards resolving `annotation.styleKey` through a style set.
- Making Storybook interactive without adding a production palette selector. The stories use local `useState` buttons to switch style sets, which proves the styling model without adding app-level UI.

### What warrants a second pair of eyes

- Whether annotation/transcript panels should require a `styleSet` prop instead of using `contextDefaultStyleSet` as a temporary default.
- Whether the default style-set helper should live in `fixtures.ts` or a dedicated preset file.

### What should be done in the future

- Update Widget IR stories and Goja DSL helpers to emit required `styleSet`.
- Run `build-storybook` after Widget IR stories are migrated.
- Search for remaining `kind`, `ContextKindSwatch`, `legendKinds`, `legendMode`, and `.kind_` references.

### Code review instructions

- Start with `src/context/types.ts` and `src/context/styles.ts`.
- Then review `ContextStyleSwatch`, `ContextLegend`, and one diagram molecule (strip) before reviewing the other diagrams.
- Validate with `pnpm --dir packages/rag-evaluation-site typecheck`.

### Technical details

- Validation run: `pnpm typecheck` in `packages/rag-evaluation-site`.
- Final result: typecheck passed after story/component fixes.

## Step 5: Widget IR and Goja DSL StyleSet Cutover

This implementation slice carried the hard-cutover contract through Widget IR stories, the Go widget DSL module, the xgoja provider documentation, and schema metadata. The JavaScript-facing `context_window.dsl` module now has helpers for building visual styles, legend items, style sets, context parts, snapshots, and palette-derived style sets.

### Prompt Context

**User prompt (verbatim):** "create it, then continue working"

**Assistant interpretation:** Create the Goja DSL design note, update tasks, then continue implementing the cutover through the DSL layer.

**Inferred user intent:** The user wants the DSL implementation to follow a documented contract and keep progressing beyond frontend-only changes.

**Commit (code):** pending — Widget IR and Goja DSL cutover

### What I did

- Created `design-doc/03-goja-dsl-styleset-cutover.md` defining the Goja DSL contract.
- Updated `tasks.md` with a Goja DSL hard-cutover checklist and marked completed frontend/DSL tasks.
- Updated `WidgetRenderer.context-diagrams.stories.tsx` so all context diagram IR examples include `styleSet` and use `styleKey` parts.
- Updated `WidgetRenderer.foundation-atoms.stories.tsx`, `WidgetRenderer.layout-recipes.stories.tsx`, and `WidgetRenderer.course-handout.stories.tsx` to use `ContextStyleSwatch` and `AnnotationBadge.visualStyle`.
- Removed `ContextKindSwatch` from the Go widget schema component list and replaced it with `ContextStyleSwatch`.
- Updated `pkg/widgetdsl/module.go`:
  - removed `contextKindSwatch`,
  - added `contextStyleSwatch`,
  - added `visualStyle`, `legendItem`, `styleSet`, `contextPart`, `contextSnapshot`, and `paletteStyleSet`,
  - updated `recipes.contextDiagram` to require `styleSet` or `palette + entries`.
- Updated `pkg/widgetdsl/module_test.go` for style-set helpers, JSON shape, and missing-styleSet errors.
- Updated xgoja provider docs to show `styleKey` + `styleSet` instead of `kind`.

### Why

The React hard cutover is only useful if server/Goja-generated IR can express the same contract. The DSL must make the correct API easy and must fail when callers omit required styling information.

### What worked

- The DSL's generic component factory already copies props through, so most component helper behavior needed no special-case code.
- Adding focused JSON-constructor helpers in `context_window.dsl` kept the Go implementation compact.
- The forbidden API search found only documentation references after the code cutover, which were straightforward to update.

### What didn't work

- The existing `recipes.contextDiagram` silently constructed a panel without style information. That had to become a deliberate error to enforce the new contract.

### What I learned

- The Goja DSL layer is mostly a JSON builder; this style-set cutover fits its existing helper pattern well.
- Tests at the DSL level are valuable because they assert the emitted JSON shape, not just React type correctness.

### What was tricky to build

- Ensuring the Go DSL can construct useful palette-derived style sets without pulling in a heavy color library. I used CSS `color-mix(...)` strings, which are JSON-compatible and resolved by the browser.
- Keeping the missing-styleSet behavior strict while still allowing a convenience path via `palette + entries`.

### What warrants a second pair of eyes

- Whether `paletteStyleSet` should support all 16 imported JSON palettes or only the four preferred palettes for now.
- Whether `contextSnapshot` should validate that every part styleKey exists in the styleSet. Currently that validation is left to React/style resolution and recipe requirements.

### What should be done in the future

- Add a full app/package build after committing this slice.
- Audit `ClubMedMeetup/` and any backend API examples for runtime pages that still emit old context diagram props.

### Code review instructions

- Review `pkg/widgetdsl/module.go` helper exports and `contextDiagramRecipe` first.
- Then review `pkg/widgetdsl/module_test.go` to understand the expected JS-facing API.
- Validate with `go test ./pkg/widgetdsl ./pkg/widgetschema ./internal/widgetmanifest -count=1`, `pnpm --dir packages/rag-evaluation-site typecheck`, and `pnpm --dir packages/rag-evaluation-site build-storybook`.

### Technical details

Validation passed:

- `pnpm typecheck` in `packages/rag-evaluation-site`
- `go test ./pkg/widgetdsl ./pkg/widgetschema ./internal/widgetmanifest -count=1`
- `pnpm build-storybook` in `packages/rag-evaluation-site`
- forbidden API search for `ContextKindSwatch`, `ContextPartKind`, `legendKinds`, `legendMode`, `.kind_`, and context-diagram `kind:` data
