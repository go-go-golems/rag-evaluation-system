# Part I: Styling, Look & Feel

## 1. Design Tokens & CSS Custom Properties

All visual constants live as CSS custom properties on `:root`, defined in
`web/src/styles/tokens.css` (production) and `web/src/styles/tokens.css` (Vite source).

### 1.1 Colour Tokens

| Token                | Value   | Usage                        |
|----------------------|---------|------------------------------|
| `--mac-bg`           | `#FFFFFF` | Page background            |
| `--mac-bg-dark`      | `#000000` | Dark mode / stripe           |
| `--mac-text`         | `#000000` | Primary text                |
| `--mac-text-dim`     | `#666666` | Secondary / de-emphasized text |
| `--mac-text-inv`     | `#FFFFFF` | Text on dark backgrounds    |
| `--mac-border`       | `#000000` | Border and rule colour      |
| `--mac-surface`      | `#FFFFFF` | Card / panel background     |
| `--mac-surface-2`    | `#EEEEEE` | Alternate row / zebra stripe|
| `--mac-accent`       | `#0000CC` | Links, primary actions       |
| `--mac-accent-2`     | `#CC0000` | Destructive / error accent  |
| `--mac-green`         | `#007722` | Success, "OK" status        |
| `--mac-amber`        | `#AA7700` | Warning / pending           |

These map one-to-one to CSS custom properties so every component consumes
variables, never hard-coded hex values.

### 1.2 Font Role Tokens

| Token | Computed Value | Role |
|---|---|---|
| `--rag-font-role-body` | 400 12px/1.45 Geneva, sans-serif | Body text |
| `--rag-font-role-compact` | 400 11px/1.4 Geneva | Dense metadata rows |
| `--rag-font-role-metadata` | 400 11px/1.35 Monaco, monospace | Labels, kv-pairs |
| `--rag-font-role-label` | 700 11px/1.2 Monaco | Section headings, tab labels |
| `--rag-font-role-metric` | 700 12px/1.25 Monaco | KPI numbers |
| `--rag-font-role-code` | 400 11px/1.45 Monaco | Code snippets, stack traces |

### 1.3 Typography Scale

There is no classical modular scale. The system uses four distinct "font
roles" that map to CSS custom properties:

- **Body** – 12px/1.45 Geneva – all readable paragraph text
- **Metadata** – 11px/1.35 Monaco – labels, data grid values
- **Label** – bold 11px Monaco – section headers, tab labels
- **Metric** – bold 12px Monaco – KPI numbers

### 1.4 Spacing and Layout

| Spacing Token | px | Common Usage |
|---|---|---|
| `xxs` | 2px | Tightest gaps |
| `xs`  | 4px | Icons to label padding |
| `sm`  | 8px  | Standard internal padding |
| `md`  | 16px | Section gaps |
| `lg`  | 24px | Page-level margins |

All spacing must be a multiple of the 4px grid.

### 1.5 Component Spacing Rules

- Every `Panel` has 12px internal padding.
- Adjacent panels use 8px gap when stacked, 16px when side-by-side.
- `DashboardGrid` applies `grid-gap: 8px`.
- Lists, Data Tables, and Metadata grids all use consistent 4px row padding.

## 2. Visual Style

The visual design references a "Classic Mac" monochrome aesthetic:

- **Background**: pure white (`#FFFFFF`) with black 1px borders on all panels.
- **Borders**: Solid 1px black; no box-shadow.
- **Buttons**: black text, 1px solid border, white background; hover darkens background to `#F5F5F5`.
- **Status indicators**: `green (#007722)`, `amber (#AA7700)`, `red (#CC0000)`.
- **Active / selected state**: Inverted (white text on black background).
- **Icons**: minimal, inline SVG only. No icon font; icons are shipped as inline SVG.

### Colour Palette Summary

```
White (#FFFFFF)    Page background
Black  (#000000)   Borders, text, stripes
Accent (#0000CC)   Primary actions
Green  (#007722)   Success, "OK" states
Amber  (#AA7700)   Warning
Red    (#CC0000)   Errors, destructive actions
```

### 3. CSS Architecture

```
src/styles/tokens.css   → Design tokens
src/styles/typography.css → Font roles
src/index.css          → Global resets, scrollbar styling
src/components/*/*.module.css → Component-scoped styles
```

Component CSS uses CSS Modules (`.module.css`). Class names are locally scoped
automatically by Vite. Tokens are referenced as `var(--mac-accent)`.

### 4. Accessibility

- All interactive elements must have visible focus indicators (`:focus-visible`).
- Contrast ratios must be at least 4.5:1 against the white background.
- Use `aria-label`, `role`, and `tabindex` on all interactive widgets.
- Status icons should include an `aria-label` and a `role="img"`.

---

## 5. Key Files Quick Reference

| File | Purpose |
|---|---|
| `web/src/styles/tokens.css` | All CSS custom property tokens |
| `web/src/index.css` | Global resets and scrollbar styling |
| `web/src/components/` | React components organized by tier (atoms, molecules, organisms, layout, pages) |
| `pkg/widgetdsl/` | Goja JavaScript module for authoring Widget IR |
| `pkg/widgetrunner/` | Runs trusted JS scripts that produce Widget IR pages |
| `pkg/widgetserver/` | HTTP API: GET /api/widget/pages/:id |
| `pkg/widgetschema/` | Widget IR JSON Schema, schema endpoint |

---

## Part I Summary

This reference is for contributors who need to:

1. **Modify styles** → edit `tokens.css` and component CSS modules
2. **Add a new widget** → extend `ir.ts`, `WidgetRenderer.tsx`, and the Goja DSL module
3. **Understand design decisions** → read the tokens, not pixel values

The CSS-only approach (no Tailwind utility classes, no CSS-in-JS) keeps the system predictable and auditable.

---

## Part 2 is in the companion document: `part-2-widgets-and-interaction.md`
