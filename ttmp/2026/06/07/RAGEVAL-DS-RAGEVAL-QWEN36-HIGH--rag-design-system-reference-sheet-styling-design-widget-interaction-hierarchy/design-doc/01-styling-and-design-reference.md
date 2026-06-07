---
title: "Styling and Design Reference"
ticket: RAGEVAL-DS-RAGEVAL-QWEN36-HIGH
topics:
  - frontend-architecture
  - design-system
  - css
  - tokens
  - storybook
what_for: >
  Comprehensive reference for designers and frontend engineers contributing to the
  RAG Evaluation System's visual language, token system, CSS architecture, and
  component styling discipline.
---

# RAG Evaluation System — Styling and Design Reference

> **Audience:** External designers, frontend interns, and React engineers who need to
> understand the visual language, CSS architecture, and styling discipline of the
> RAG Evaluation System dashboard.
>
> **Status:** Active · **Target version:** `web/` (RAG Evaluation System, Go 2025)

## 1. Executive Summary

The RAG Evaluation System uses a **retro monochrome Mac OS 1-inspired design language**
delivered through a layered React component architecture with CSS Modules. Every visual
decision has a single owner at a specific layer. The system is organized into five
layers:

```
tokens → foundation → atoms → layout primitives → molecules → organisms → pages
```

Global CSS (`index.css`) now contains only Tailwind import, token import, body defaults,
and scrollbar styles. All component anatomy lives in local CSS Modules. Storybook
stories provide the primary review surface for every public component.

## 2. Visual Identity

### 2.1 The Retro Monochrome Language

The dashboard looks like a high-density operational interface from the early Macintosh
era: black panel headers, thin black borders, compact monospace labels, sparse accent
colors, and square surfaces with no border radius.

| Property | Value | Used for |
|---|---|---|
| Background | `#FFFFFF` | Page background, panel surfaces |
| Background dark | `#000000` | Panel headers, body background, active states |
| Text | `#000000` | Primary text |
| Text dim | `#666666` | Muted/disabled text |
| Text inverse | `#FFFFFF` | Text on dark backgrounds |
| Border | `#000000` | All 1px borders |
| Surface | `#FFFFFF` | Card panels |
| Surface 2 | `#EEEEEE` | Secondary fills, scrollbar tracks |
| Accent (blue) | `#0000CC` | Running status, link-like states |
| Accent 2 (red) | `#CC0000` | Failed/error status |
| Green | `#007722` | Succeeded/done status |
| Amber | `#AA7700` | Partial/warning status |
| Border radius | `0px` (square) | All surfaces |

**Rule:** The system intentionally uses zero border radius. Square corners are part of
the retro visual language. Adding radius would break the design identity.

### 2.2 Typography Roles

Typography is organized into named roles, not size ranges. Each role maps to a complete
`font` shorthand consuming CSS variables.

```css
/* --rag-font-role-body:   400 12px/1.45 var(--font-body);       */
/* --rag-font-role-compact: 400 11px/1.4 var(--font-body);       */
/* --rag-font-role-metadata: 400 11px/1.35 var(--font-mono);    */
/* --rag-font-role-label:     700 11px/1.2 var(--font-mono);    */
/* --rag-font-role-metric:    700 12px/1.25 var(--font-mono);   */
/* --rag-font-role-code:      400 11px/1.45 var(--font-mono);   */
```

| Role | Weight/Size/Line | Font Family | Used for |
|---|---|---|---|
| `--rag-font-role-body` | 400 / 12px / 1.45 | `--font-body` (sans-serif) | Readable dashboard copy, placeholder text |
| `--rag-font-role-compact` | 400 / 11px / 1.4 | `--font-body` | Dense panels, constrained vertical space |
| `--rag-font-role-metadata` | 400 / 11px / 1.35 | `--font-mono` | Metadata values, table headers |
| `--rag-font-role-label` | 700 / 11px / 1.2 | `--font-mono` | Uppercase categorical labels |
| `--rag-font-role-metric` | 700 / 12px / 1.25 | `--font-mono` | Metric values, counts |
| `--rag-font-role-code` | 400 / 11px / 1.45 | `--font-mono` | Code identifiers, chunk IDs, model names |

**Rule:** Choose the role by *semantic function* (is it a label? metadata? body text?),
not by visual similarity. If you find yourself matching two elements by eye, they
should share a role.

### 2.3 Status Vocabulary

The system uses a fixed set of 10 status values with associated colors and Unicode icons:

| Status | Color | Icon | Used for |
|---|---|---|---|
| `pending` | dim | ◌ | Queued, waiting |
| `ready` | dim | ◌ | Ready to execute |
| `running` | blue (accent) | ● | Currently executing |
| `succeeded` | green | ✔ | Completed successfully |
| `done` | green | ● | Fully completed |
| `partial` | amber | ◐ | Partial completion |
| `warning` | amber | ⚠ | Non-fatal issue |
| `failed` | red (accent 2) | ✘ | Operation failed |
| `error` | red (accent 2) | ✘ | Fatal error |
| `canceled` | dim, strikethrough | ⊘ | Operation was canceled |

**Important:** The icon for `succeeded` is `✔` (checkmark). Using `success` instead of
`succeeded` will render `?` because the icon lookup is case-insensitive but exact-match
only on the listed values.

## 3. CSS Architecture

### 3.1 Layer Ownership Rules

Every CSS decision must have a single owner. This is the core rule:

```
tokens.css        → color, font, border, spacing, semantic aliases
foundation/       → text roles, code text, captions, status text, dividers, hidden text
atoms/            → buttons, icon buttons, inputs, selects, checkbox rows, error callouts
layout/           → panels, stacks, inline groups, dashboard grids, shells, tabs, form rows, scroll regions
molecules/        → reusable data/display patterns (DataTable, MetadataGrid, AppNav...)
organisms/        → feature panels with DTO-shaped props (SearchControlsPanel, etc.)
pages/            → page composition and page-local styling
```

**What does NOT belong at each layer:**

| Layer | Must NOT contain |
|---|---|
| `tokens.css` | Component selectors, feature-specific anatomy |
| Foundation | Grids, panels, API behavior, component-specific styling |
| Atoms | Page state, data fetching, table rendering, panel layout |
| Layout | Domain formatting, API behavior, feature-specific rendering |
| Molecules | RTK Query hooks, page orchestration, container logic |

### 3.2 The Global Stylesheet

`web/src/index.css` is now intentionally minimal:

```css
@import "tailwindcss";
@import "./styles/tokens.css";

* { box-sizing: border-box; }

body {
  font-family: var(--font-body);
  font-size: 12px;
  line-height: 1.45;
  color: var(--mac-text);
  background: var(--mac-bg);
  margin: 0;
  padding: 0;
  -webkit-font-smoothing: none;
}

/* Scrollbar styling (retro thick style) */
::-webkit-scrollbar { width: 16px; height: 16px; }
::-webkit-scrollbar-track { background: var(--mac-surface); border: 1px solid #CCCCCC; }
::-webkit-scrollbar-thumb { background: var(--mac-surface-2); border: 1px solid var(--mac-border); }
```

**Rule:** If you need to add global CSS, ask whether it belongs to a component layer
first. Only add to `index.css` if it affects the entire application (body defaults,
scrollbars, Tailwind imports).

### 3.3 CSS Modules

Every component (except foundation and atoms which sometimes share global-ish styling)
uses CSS Modules. The pattern is:

```
WidgetName/
  WidgetName.tsx          ← component with className prop
  WidgetName.module.css   ← local styles
  WidgetName.stories.tsx  ← Storybook story
  index.ts                ← re-export
```

The component receives `className` as a prop and merges it:

```tsx
// Pattern — always merge className
<div className={[styles.root, className ?? ''].filter(Boolean).join(' ')}
     data-rag-atom="Button"
     {...rest}
>
```

**Rule:** CSS Modules own local anatomy. They do not redefine generic panel headers,
table cells, status tones, form controls, or typography roles. Those belong to
foundation or appropriate layers.

### 3.4 The Safe Deletion Rule

Never delete a global class before migrating all consumers:

```
1. Identify the class or class family to remove
2. Scan for active consumers:   rg -n '<class-name>' web/src
3. Migrate active consumers to the component owner
4. Scan again:                 rg -n '<class-name>' web/src
5. Review false positives (e.g., .selected vs .selectable in CSS Modules)
6. Only then remove from index.css
7. Run typecheck + Storybook + visual smoke
```

## 4. Design Tokens

### 4.1 Token File

**Location:** `web/src/styles/tokens.css`

The token file defines two families:

| Family | Purpose | Example |
|---|---|---|
| `--mac-*` | Original Mac OS 1 color palette, font families | `--mac-bg`, `--mac-border`, `--font-mono` |
| `--rag-font-role-*` | Semantic typography aliases | `--rag-font-role-metadata`, `--rag-font-role-label` |

**Rule:** Foundation components and molecules consume `--rag-font-role-*` names.
Atoms may consume either `--mac-*` or `--rag-*` names depending on their layer.

### 4.2 Standalone Package Token Bridge

When components are extracted to the standalone npm package
(`@go-go-golems/rag-evaluation-site`), a token bridge maps the original `--mac-*`
and `--font-*` names onto package-local `--rag-*` tokens:

```css
:root {
  --rag-color-bg: #f6f7f8;
  --rag-color-surface: #ffffff;
  --rag-color-text: #1d232a;
  --rag-color-border-strong: #000000;

  --mac-bg: var(--rag-color-bg);
  --mac-bg-dark: #000000;
  --mac-text: var(--rag-color-text);
  --mac-text-inv: #ffffff;
  --mac-border: var(--rag-color-border-strong);
  --mac-surface: var(--rag-color-surface);
  --font-mono: var(--rag-font-mono);
  --rag-font-role-metadata: 400 11px/1.35 var(--font-mono);
}
```

**Rule:** The canonical tokens for the standalone package are `--rag-color-*` and
`--rag-font-*`. The `--mac-*` aliases exist for backward compatibility with components
that have not been migrated.

## 5. Component Anatomy

### 5.1 Foundation Primitives

Foundation owns text roles, accessibility, and semantic text decoration.

| Component | Element | Props (key) | data attribute |
|---|---|---|---|
| `Text` | `<p>` | `size: "normal" \| "compact"` | — |
| `CodeText` | `<code>` | — | — |
| `StatusText` | `<span>` | `status: RagStatus`, `icon: boolean` | — |
| `Caption` | `<span>` | `tone: CaptionTone`, `transform: CaptionTransform`, `truncate: boolean` | `data-rag-foundation="Caption"` |
| `Divider` | `<hr>` | — | — |
| `VisuallyHidden` | `<span>` | — | — |

**StatusText details:**
- Consumes `--rag-font-role-code` (via CSS variable fallback)
- Each status value maps to a CSS class (e.g., `.succeeded { color: var(--mac-green); }`)
- The icon is a Unicode character prepended before the label: `✔ succeeded`

**Caption tones:**

| Tone | Color |
|---|---|
| `muted` | Default (inherits) |
| `accent` | Blue (accent) |
| `success` | Green |
| `warning` | Amber |
| `danger` | Red (accent 2) |
| `inherit` | Inherit from parent |

**Caption transforms:**

| Transform | Effect |
|---|---|
| `none` | Normal casing |
| `uppercase` | `text-transform: uppercase` |

### 5.2 Atoms

Atoms own the appearance and state of native HTML controls.

| Component | Element | Props (key) | data attribute |
|---|---|---|---|
| `Button` | `<button>` | `variant: "default" \| "primary"`, `size: "normal" \| "compact"`, `selected: boolean`, `disabled: boolean` | `data-rag-atom="Button"` |
| `IconButton` | `<button>` | `label: string`, `icon: string` | `data-rag-atom="IconButton"` |
| `TextInput` | `<input>` | All native input props | `data-rag-atom="TextInput"` |
| `SelectInput` | `<select>` | All native select props | `data-rag-atom="SelectInput"` |
| `CheckboxRow` | `<label>` | `checked: boolean`, `onChange`, children (label) | `data-rag-atom="CheckboxRow"` |
| `ErrorCallout` | `<div>` | `children` (error message) | — |

**Button states and appearance:**

| State | Background | Text Color | Border |
|---|---|---|---|
| Default (idle) | `var(--mac-surface)` | `var(--mac-text)` | 1px `var(--mac-border)` |
| Default (hover) | `var(--mac-surface-2)` | `var(--mac-text)` | 1px `var(--mac-border)` |
| Default (active) | `var(--mac-bg-dark)` | `var(--mac-text-inv)` | 1px `var(--mac-border)` |
| Primary | `var(--mac-bg-dark)` | `var(--mac-text-inv)` | 1px `var(--mac-border)` |
| Selected | Same as primary | Same as primary | 1px `var(--mac-border)` |
| Disabled | `var(--mac-surface)` | `var(--mac-text-dim)` | 1px `var(--mac-border)` |

**Button sizing:**

| Size | Padding | Font |
|---|---|---|
| Normal | `2px 12px` | `--rag-font-role-metadata` (11px mono) |
| Compact | `1px 6px` | `--rag-font-role-metadata` (11px mono) |

### 5.3 Layout Primitives

Layout owns structure: panels, stacks, grids, shells, scrolling regions, and form rows.

| Component | Element | Props (key) | data attribute |
|---|---|---|---|
| `Panel` | `<section>` | `title: ReactNode`, `actions: ReactNode`, `density: "normal" \| "condensed"`, `fill: boolean` | `data-rag-layout="Panel"` |
| `Stack` | `<div>` | `gap: StackGap`, `align: StackAlign` | — |
| `Inline` | `<div>` | `gap: InlineGap`, `justify: InlineJustify`, `wrap: boolean` | — |
| `DashboardGrid` | `<div>` | `recipe: DashboardGridRecipe` | `data-rag-layout="DashboardGrid"`, `data-rag-layout-recipe={recipe}` |
| `AppShell` | `<div>` | `header: ReactNode`, `sidebar: ReactNode` | — |
| `ScrollRegion` | `<div>` | `axis: "y" \| "x" \| "both"` | — |
| `TabList` | `<div>` | `items`, `activeId`, `ariaLabel`, `onChange` | — |
| `FormRow` | `<div>` | `label: ReactNode`, `control: ReactNode`, `hint: ReactNode` | — |

**Panel anatomy:**

```
<section class="root">
  <div class="header">              ← dark background, uppercase title, optional actions
    <span>{title}</span>
    {actions}
  </div>
  <div class="body|condensed">      ← padding: 8px (normal) or 4px 6px (condensed)
    {children}
  </div>
</section>
```

| Property | Normal | Condensed |
|---|---|---|
| Body padding | `8px` | `4px 6px` |

**Stack and Inline spacing values:**

```
"xs" → 4px
"sm" → 8px
"md" → 12px
"lg" → 16px
"xl" → 24px
```

**DashboardGrid recipes:**

| Recipe | Layout | Used by |
|---|---|---|
| `twoColumn` (default) | Two-column grid | Most panels, detail views |
| `searchWorkbench` | Search-specific grid | SearchWorkbenchPage |
| `corpusExplorer` | Corpus-specific grid | CorpusExplorerView |

**FormRow layout:**

```
<div class="form-row">
  <dt class="label">{label}</dt>
  <dd class="control">{control}</dd>
  {hint && <dd class="hint">{hint}</dd>}
</div>
```

### 5.4 Molecules

Molecules own reusable data-display patterns that span multiple features.

| Component | Props (key) | data attribute |
|---|---|---|
| `AppNav` | `brand`, `items`, `activeItemId`, `onItemSelect` | `data-rag-component="AppNav"` |
| `DataTable` | `columns`, `rows`, `getRowKey`, `selectedKey`, `onRowSelect`, `emptyMessage` | `data-rag-component="DataTable"` |
| `MetadataGrid` | `items: [{key, value, copyValue?}]`, `density: "normal" \| "compact"` | `data-rag-component="MetadataGrid"` |
| `CoveragePanel` | `coverageData` | — |
| `QueryPresetList` | `presets` | — |

**MetadataGrid anatomy:**

```
<dl class="root|compact">
  <div class="row">
    <dt class="key">{key}</dt>
    <dd class="value">
      {value}
      {copyValue && <button class="copyButton" onClick={() => copy(copyValue)}>⧉</button>}
    </dd>
  </div>
</dl>
```

| Property | Normal | Compact |
|---|---|---|
| Row gap | default | reduced |

### 5.5 Organisms

Organisms are feature-specific panels. They accept DTO-shaped props and callbacks.

| Organism | Purpose | Key Dependencies |
|---|---|---|
| `SearchControlsPanel` | Search query composition, retriever selection, source filter | `Panel`, `Stack`, `TextInput`, `Button`, `Caption`, `FormRow`, `ScrollRegion`, `CheckboxRow` |
| `RetrievalResultsPanel` | Dense table of retrieval results | `Panel`, `DataTable`, `DataTableCellSpec` |
| `ResultInspectorPanel` | Selected result detail view | `Panel`, `MetadataGrid`, `Caption`, `StatusText`, `Divider` |
| `QueueHealthPanel` | Workflow queue status overview | `Panel`, `StatusText`, `DataTable` |
| `WorkflowListPanel` | Workflow list with status indicators | `Panel`, `DataTable`, `StatusText` |
| `WorkflowSummaryPanel` | Workflow aggregate statistics | `Panel`, `StatusText`, `DataTable` |
| `WorkflowOpGraphPanel` | Workflow operation dependency graph | `Panel` |
| `WorkflowOpGroupsPanel` | Grouped workflow operations | `Panel` |
| `WorkflowOpInspectorPanel` | Single operation detail with failures | `Panel`, `ErrorCallout`, `MetadataGrid`, `Caption`, `Button` |
| `WorkflowOpResultPanel` | Operation result rendering | `Panel` |

**Rule:** Organisms accept DTO-shaped data and callback props. They should not own RTK
Query hooks or cross-view navigation. That belongs in the container/page layer.

### 5.6 Page Boundaries

Pages own composition of organisms and page-local state.

| Page | Path | Purpose |
|---|---|---|
| `SearchWorkbenchPage` | `pages/SearchWorkbenchPage` | Search query + results + inspector |
| `PipelinePage` | `pages/PipelinePage` | Pipeline overview and status |
| `EvaluationPage` | `pages/EvaluationPage` | Evaluation results and comparisons |
| `DslPreviewPage` | `pages/DslPreviewPage` | Widget DSL preview |

### 5.7 Corpus Domain Widgets

Corpus-specific widgets remain in `corpus/` but adopt shared primitives and atoms.

| Widget | Purpose |
|---|---|
| `SourcePanel` | Source metadata and configuration |
| `DocumentBrowser` | Document list with search/filter |
| `DocumentInspector` | Single document detail with tabs and metadata |
| `IdentityBar` | Document/source identity row |
| `ArtifactIdentityBar` | Artifact-specific identity row |
| `ChunkTimelineBar` | Chunk timeline visualization |

## 6. Data Attributes for Tooling

The system uses explicit `data-*` attributes for visual diff tooling, Playwright smoke
tests, and accessibility scanning:

| Attribute | Example Value | Set on | Purpose |
|---|---|---|---|
| `data-rag-page` | `"RagEvaluationSiteApp"` | App root | Identify the app container |
| `data-rag-shell` | `"default"` | Shell wrapper | Identify shell presence |
| `data-rag-layout` | `"Panel"`, `"DashboardGrid"` | Layout primitives | Identify layout structure |
| `data-rag-layout-recipe` | `"searchWorkbench"` | DashboardGrid | Identify grid recipe |
| `data-rag-component` | `"DataTable"`, `"MetadataGrid"`, `"AppNav"` | Molecules | Identify molecule type |
| `data-rag-atom` | `"Button"`, `"TextInput"`, `"SelectInput"`, `"CheckboxRow"` | Atoms | Identify atom type |
| `data-rag-foundation` | `"Caption"` | Foundation | Identify foundation element |
| `data-page-id` | `"demo"`, `"smoke"` | App root | Identify current page |

**Rule:** Visual diff tooling should target explicit semantic attributes, not incidental
HTML structure. Use selectors like `[data-rag-component="DataTable"]` rather than
`table:first-of-type`.

## 7. Storybook Coverage

Storybook is the primary review surface for every reusable component.

### 7.1 Story Location and Naming

Stories live alongside their components:

```
WidgetName/
  WidgetName.tsx
  WidgetName.module.css
  WidgetName.stories.tsx    ← always present for public components
  index.ts
```

Storybook title convention:

```
Design System/Foundation/Overview    ← foundation overview
Design System/Atoms/Button           ← atom stories
Design System/Layout/Panel           ← layout stories
Component Library/Molecules/DataTable ← molecule stories
```

### 7.2 Required Stories

Every public primitive, atom, molecule, organism, and page boundary should have at least
one story. The story shows:

1. Default state
2. Key variants (e.g., Button variants, StatusText statuses)
3. Edge cases (e.g., disabled states, empty states, dense data)
4. Composition examples (e.g., Panel with actions + MetadataGrid)

### 7.3 Foundation Overview Story

The Foundation story (`Design System/Foundation/Overview`) provides a single-page
reference for the entire design language:

- **Colors** — swatch grid showing all token colors with CSS variable names
- **Typography** — specimen of each font role
- **Status Tones** — all 10 statuses with icons
- **Spacing** — visual spacing scale
- **Borders and Radii** — square corners, 1px borders, surface fills
- **Accessibility** — VisuallyHidden example

### 7.4 Storybook as Review Surface

A component without a story is not yet a stable design-system component. Stories
demonstrate state contracts, visual regressions, and composition patterns. When adding
or modifying a component, add or update its story simultaneously.

## 8. Adding a New Component

Follow this checklist when adding a new component:

1. **Decide the layer** using the ownership table in Section 3.1
2. **Create the directory** following the standard file pattern
3. **Write the component** with `className` prop merging and `data-rag-*` attributes
4. **Write the CSS Module** — own local anatomy only, reuse tokens from parent layers
5. **Write the Storybook story** — show default state, key variants, and edge cases
6. **Export from the layer index** — add to the layer's `index.ts` barrel export
7. **Run Storybook** — verify the story renders correctly
8. **Add visual QA** — if the component is visible in the real app, add it to the
   visual smoke screenshot set

### 8.1 Naming Convention

- Component files: PascalCase (e.g., `StatusText.tsx`)
- CSS modules: camelCase (e.g., `StatusText.module.css`)
- CSS classes: camelCase within modules (e.g., `root`, `header`, `condensed`)
- Data attributes: kebab-case (e.g., `data-rag-atom="StatusText"`)
- Storybook titles: PascalCase with path separators (e.g., `Design System/Atoms/StatusText`)

### 8.2 Styling a New Atom

An atom should follow the Button pattern:

```tsx
// NewAtom.tsx
import styles from './NewAtom.module.css';

export type NewAtomVariant = 'default' | 'highlight';

export interface NewAtomProps extends HTMLAttributes<HTMLDivElement> {
  variant?: NewAtomVariant;
  children?: ReactNode;
}

export function NewAtom({ variant = 'default', className, children, ...rest }: NewAtomProps) {
  return (
    <div
      className={[styles.root, variant === 'highlight' ? styles.highlight : '', className ?? ''].filter(Boolean).join(' ')}
      data-rag-atom="NewAtom"
      {...rest}
    >
      {children}
    </div>
  );
}
```

```css
/* NewAtom.module.css */
.root {
  font: var(--rag-font-role-metadata);
  padding: 2px 8px;
  border: 1px solid var(--mac-border);
  background: var(--mac-surface);
}

.highlight {
  background: var(--mac-bg-dark);
  color: var(--mac-text-inv);
}
```

### 8.3 Styling a New Layout Primitive

A layout primitive should follow the Panel pattern:

```css
/* NewLayout.module.css */
.root {
  border: 1px solid var(--mac-border);
  background: var(--mac-surface);
  overflow: hidden;
}

.header {
  background: var(--mac-bg-dark);
  color: var(--mac-text-inv);
  font: var(--rag-font-role-label);
  padding: 3px 8px;
  letter-spacing: 0.5px;
  text-transform: uppercase;
  border-bottom: 1px solid var(--mac-border);
}

.body {
  padding: 8px;
}
```

## 9. Common Pitfalls

### 9.1 Re-defining Primitives in Molecules

**Bad:** Defining `.panel-header` inside `WorkflowOpInspectorPanel.module.css` because
the panel header looks slightly different.

**Good:** Use the shared `Panel` component with `title` and `actions` props. If the
header needs variant styling, add a variant prop to `Panel` itself, not a CSS Module
duplicate.

### 9.2 Mixing Layer Concerns in a Component

**Bad:** A search controls panel that sets its own button colors, border styles, and
font sizes inline in JSX.

**Good:** The panel composes `Button`, `TextInput`, `Caption`, and `FormRow`. The
button's colors come from the Button atom, not from the panel's CSS.

### 9.3 Using Inline Styles for Visual Decisions

**Bad:** `<div style={{ background: '#000000', color: '#FFFFFF' }}>` in a component.

**Good:** Use tokens and CSS Modules. If the style is dynamic (e.g., based on props),
use variant classes or style props that map to tokens: `style={{ color: var(--mac-text)
}}`.

### 9.4 Breaking the Token Bridge in Standalone Mode

When working on the standalone npm package, remember that components may consume
`--mac-*` variables that only exist because of the token bridge. If a component looks
broken in the standalone package, check computed CSS first — the variable may be missing
from the bridge mapping.

## 10. File Reference

| File | Purpose |
|---|---|
| `web/src/styles/tokens.css` | All design tokens (Mac palette + font role aliases) |
| `web/src/index.css` | Global CSS: Tailwind import, body defaults, scrollbars |
| `web/src/components/atoms/` | Button, IconButton, TextInput, SelectInput, CheckboxRow, ErrorCallout |
| `web/src/components/foundation/` | Text, CodeText, StatusText, Caption, Divider, VisuallyHidden |
| `web/src/components/layout/` | Panel, Stack, Inline, DashboardGrid, AppShell, ScrollRegion, TabList, FormRow |
| `web/src/components/molecules/` | AppNav, DataTable, MetadataGrid, CoveragePanel, QueryPresetList |
| `web/src/components/organisms/` | Feature panels (SearchControlsPanel, RetrievalResultsPanel, etc.) |
| `web/src/components/pages/` | SearchWorkbenchPage, PipelinePage, EvaluationPage, DslPreviewPage |
| `web/src/components/corpus/` | Corpus domain widgets (SourcePanel, DocumentInspector, etc.) |
| `web/src/widgets/` | Widget IR types, WidgetRenderer, cellRenderers, actions |

## 11. Quick Reference: CSS Variable Map

```
--mac-bg           → #FFFFFF          (page/panel background)
--mac-bg-dark      → #000000          (panel headers, active states)
--mac-text         → #000000          (primary text)
--mac-text-dim     → #666666          (muted text, disabled)
--mac-text-inv     → #FFFFFF          (text on dark backgrounds)
--mac-border       → #000000          (all borders)
--mac-surface      → #FFFFFF          (card panels)
--mac-surface-2    → #EEEEEE          (secondary fills)
--mac-accent       → #0000CC          (blue - running state)
--mac-accent-2     → #CC0000          (red - failed/error state)
--mac-green        → #007722          (success state)
--mac-amber        → #AA7700          (warning/partial state)
--font-body        → sans-serif stack
--font-mono        → Monaco, Courier New, monospace
--rag-font-role-body         → 400 12px/1.45 var(--font-body)
--rag-font-role-compact      → 400 11px/1.4 var(--font-body)
--rag-font-role-metadata     → 400 11px/1.35 var(--font-mono)
--rag-font-role-label        → 700 11px/1.2 var(--font-mono)
--rag-font-role-metric       → 700 12px/1.25 var(--font-mono)
--rag-font-role-code         → 400 11px/1.45 var(--font-mono)
```
