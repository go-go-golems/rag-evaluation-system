# `@go-go-golems/rag-evaluation-site` Design-System Guidelines

This package is the canonical reusable UI layer for the RAG Evaluation System. Treat it as a strict design system, not as a scratchpad for one-off page CSS.

Read this file before adding or changing components in `packages/rag-evaluation-site`.

## Non-negotiable rules

1. **Use the package layers.** Put reusable structure in `components/layout`, reusable data/content patterns in `components/molecules`, and domain panels with DTO-shaped props in `components/organisms`.
2. **Use existing foundation typography.** Prefer `Text`, `Caption`, `CodeText`, `StatusText`, and the `--rag-font-role-*` tokens before writing font declarations.
3. **Do not invent ad-hoc typography.** No random `font-size: 13px`, `font-weight: 600`, `letter-spacing: .16em`, or `line-height` unless it is being added deliberately to a token/foundation primitive.
4. **Do not create CSS-in-JS systems.** No generic `Box`, no arbitrary style props, no components that accept `padding`, `gap`, `color`, `fontSize`, and `display` as an unbounded styling API.
5. **Keep components API-free.** Package components must not import RTK Query hooks, app stores, router state, backend services, or web-only containers.
6. **Storybook is required.** New public components need package-owned stories before they are considered part of the design system.
7. **React first, Widget IR later.** Stabilize the React component API and visual states before adding semantic `WidgetRenderer`/Goja DSL support.
8. **Use `data-rag-*` identity attributes.** These are required for visual-diff extraction, prototype parity checks, and future IR/runtime inspection.
9. **Prototype JSX is reference material, not production code.** Copy concepts, anatomy, and token usage; do not copy global classes, inline styles, CDN assumptions, or app-local state.

## Package layer ownership

```text
src/theme.css
  -> foundation primitives
  -> atoms
  -> layout primitives
  -> molecules
  -> organisms
  -> widgets / WidgetRenderer
```

### `src/theme.css`

Owns raw tokens and compatibility bridge tokens:

- `--rag-color-*`
- `--mac-*`
- `--font-body`, `--font-mono`
- `--rag-font-role-body`
- `--rag-font-role-compact`
- `--rag-font-role-metadata`
- `--rag-font-role-label`
- `--rag-font-role-metric`
- `--rag-font-role-code`

If a visual choice repeats, add or adjust a token/foundation primitive instead of duplicating literals across CSS modules.

### Foundation

Foundation owns text and semantic display roles:

- `Text` — body/compact/metadata/label/metric text with tone/weight/alignment.
- `Caption` — compact metadata/caption copy.
- `CodeText` — code-like IDs, paths, model names, compact technical values.
- `StatusText` — status-to-tone mapping.
- `Divider` — separators.
- `VisuallyHidden` — accessibility-only text.

Use these roles instead of local typography whenever possible.

### Atoms

Atoms own small controls and semantic markers. They may use foundation typography and theme tokens, but they should not own page layout.

Examples: `Button`, `TextInput`, `SelectInput`, `CheckboxRow`, `IconButton`, `ContextStyleSwatch`, `AnnotationBadge`.

### Layout

Layout owns generic structure and spacing recipes. Layout components must not know domain nouns like “course agenda”, “transcript message”, or “context part”.

Examples:

- `AppShell`
- `SidebarShell`
- `Panel`
- `Stack`
- `Inline`
- `DashboardGrid`
- `SplitPane`
- `SectionBlock`
- `SlideShell`
- `ScrollRegion`
- `TabList`
- `FormRow`

If a component answers “where do regions go?” it belongs in layout. If it answers “what domain data is shown?” it does not.

### Molecules

Molecules own reusable data-display or content-display patterns. They should have typed props, no app services, and no route/store assumptions.

Examples:

- `DataTable`
- `MetadataGrid`
- `SidebarNav`
- `KeyValueStrip`
- `CheckList`
- `StepList`
- `KeyPointList`
- `FigureBlock`
- `PersonSummary`
- context diagram molecules
- transcript/annotation cards

If a pattern is useful outside one page but still carries content semantics, it is probably a molecule.

### Organisms

Organisms compose atoms/layout/molecules into feature panels or product-level shells with DTO-shaped props and callbacks.

Examples:

- `ContextDiagramPanel`
- `TranscriptReaderPanel`
- `AnnotationRailPanel`
- `AnchoredCommentRail`
- `CourseLessonPanel`
- `CourseSlidePanel`
- `CourseStudioShell`

Organisms may be domain-specific, but they still must be presentational. They should accept data and callbacks; containers decide where data comes from.

### Web-owned code

Keep the following outside this package unless explicitly split into API-free presentational props:

- RTK Query hooks and services
- app store access
- route state and navigation adapters
- upload/parsing workflows tied to backend/runtime behavior
- backend-connected containers/views
- app-specific pages that are not reusable package surfaces

## Storybook conventions

Use package Storybook as the canonical review surface for reusable package components.

Required title prefixes:

```text
Design System/Foundation/<Primitive>
Design System/Atoms/<Atom>
Design System/Layout/<Primitive>
Component Library/Molecules/<Component>
Component Library/Organisms/<Component>
Widget IR/Renderer
```

Required story states where applicable:

- default/populated
- empty
- overflow/dense
- selected/active
- disabled
- error/warning
- alternate layout direction, e.g. visual left/right

Do not hide reusable stories in `web`. `web` Storybook is for app containers, backend-connected views, and integration/page compositions.

## Typography rules

The current visual language is intentionally compact and tokenized. Use these roles:

| Purpose | Use |
|---|---|
| Body prose | `Text size="body"` or `font: var(--rag-font-role-body)` |
| Readable article prose | `font: var(--rag-font-role-readable-body)` inside document/article renderers |
| Article/display heading | `font: var(--rag-font-role-display)` or `font: var(--rag-font-role-heading)` in document/article renderers |
| Dense labels / navigation rows | `Text size="compact"` or `font: var(--rag-font-role-compact)` |
| Metadata / captions / side notes | `Caption` or `font: var(--rag-font-role-metadata)` |
| Uppercase section labels | `Text size="label"` or `font: var(--rag-font-role-label)` plus tokenized uppercase behavior |
| Compact numeric/metric values | `Text size="metric"` or `font: var(--rag-font-role-metric)` |
| Code-like technical text | `CodeText` or `font: var(--rag-font-role-code)` |

Do not set component-local typography unless the role does not exist yet. If a role is missing, add the role deliberately and document it.

### Sidebar typography example

The imported context-viewer prototype defines sidebar rows and groups with the same font roles:

```css
.mac-navitem { font: var(--rag-font-role-compact); }
.mac-navgroup { font: var(--rag-font-role-label); letter-spacing: 0.5px; }
.mac-caption { font: var(--rag-font-role-metadata); }
```

Package sidebar components should follow this token language, not custom body-font weights or arbitrary letter-spacing.

## CSS module rules

CSS Modules are allowed and expected, but they own only local anatomy.

Allowed in CSS modules:

- component layout anatomy
- state selectors such as active/selected/disabled
- dynamic visualization geometry wrappers
- overflow/clamp behavior
- local borders/backgrounds using theme tokens

Avoid in CSS modules:

- raw colors when a token exists
- arbitrary font literals
- global selectors
- utility classes
- domain-specific layout that should be a layout primitive
- broad styling props that turn components into style systems

Inline styles are allowed only for truly dynamic geometry, CSS variable plumbing, or Storybook-only swatches.

## Prototype-source usage

The imported prototype under the context ticket is valuable for:

- product vocabulary
- screen hierarchy
- spacing and rhythm cues
- `data-rag-*` extraction targets
- visual parity comparison
- token/foundation usage examples

It is not production architecture. Do not copy:

- global `.mac-*` classes directly into package components
- inline JSX styles as permanent implementation
- `window` exports
- CDN/Babel assumptions
- screen-local state into reusable components

When porting from prototype source:

1. Identify the generic layer: layout, molecule, organism, or app container.
2. Replace global classes with CSS Modules and package tokens.
3. Replace ad-hoc typography with foundation roles.
4. Add Storybook stories for the extracted component.
5. Add `data-rag-*` identity attributes.
6. Validate typecheck/build/Storybook before using it in `web`.

## Widget IR / WidgetRenderer rules

`src/widgets/ir.ts` and `src/widgets/WidgetRenderer.tsx` expose JSON-compatible semantic UI nodes. Do not extend them with unstable visual fragments.

Add Widget IR support only when:

1. The React component is stable and story-covered.
2. Its props are mostly JSON-compatible.
3. There is a clear semantic reason for Goja/IR authors to use it.
4. The renderer can preserve existing actions (`copy`, `event`, `navigate`, `server`) and typed value normalization.
5. There are WidgetRenderer stories or tests that cover the new node.

Prefer high-level semantic recipes over low-level DOM recreation. For example, a future `CourseStudioShell` or `ContextDiagramPanel` node is better than emitting many anonymous `div` nodes that mimic CSS module internals.

## Review checklist

Before committing package UI changes:

- [ ] Correct layer chosen.
- [ ] No backend/router/store imports.
- [ ] Typography uses foundation roles or `--rag-font-role-*` tokens.
- [ ] Colors, borders, and surfaces use theme tokens.
- [ ] CSS module owns local anatomy only.
- [ ] Public component has folder layout: `Component.tsx`, `Component.module.css`, `Component.stories.tsx`, `index.ts`.
- [ ] Storybook title uses the correct prefix.
- [ ] `data-rag-*` identity attribute is present.
- [ ] `pnpm --dir packages/rag-evaluation-site typecheck` passes.
- [ ] Build/Storybook build runs for larger changes.
- [ ] Diary/changelog updated when working under a ticket.

## Key references

- `src/theme.css`
- `src/components/foundation/*`
- `src/widgets/ir.ts`
- `src/widgets/WidgetRenderer.tsx`
- `../web/README.md`
- `../../ttmp/2026/06/07/RAGEVAL-DESIGN-SYSTEM-UNIFY--unify-web-design-system-into-reusable-rag-evaluation-site-package/design-doc/01-design-system-unification-analysis-and-implementation-guide.md`
- `../../ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/design-doc/03-context-viewer-integration-plan-after-design-system-unification.md`
- `../../ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/sources/03-context-viewer-design-iteration/`
