---
Title: PaletteProvider implementation plan
Ticket: RAGEVAL-PALETTE-PROVIDER
Status: complete
Topics:
    - design-system
    - frontend
    - storybook
    - theming
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/.storybook/preview.ts
      Note: |-
        Adds Storybook global palette toolbar/decorator.
        Storybook global palette toolbar/decorator
    - Path: /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/context/PaletteProvider.tsx
      Note: |-
        Implements provider, CSS variable mapping, and palette hooks.
        Provider implementation and palette CSS variable mapping
    - Path: /home/manuel/workspaces/2026-07-03/improve-rag-evaluation-system/rag-evaluation-system/packages/rag-evaluation-site/src/widgets/WidgetRenderer.transcript-notes.stories.tsx
      Note: |-
        Transcript stories now use provider fallback style sets.
        Transcript stories consume provider fallback
    - Path: repo://packages/rag-evaluation-site/src/components/layout/Panel/Panel.module.css
      Note: Panel heading title styling
ExternalSources: []
Summary: Implementation plan and final behavior for global palette switching in rag-evaluation-site Storybook.
LastUpdated: 0001-01-01T00:00:00Z
WhatFor: Use when reviewing or extending PaletteProvider, Storybook palette switching, or transcript palette behavior.
WhenToUse: Before adding more palette-aware components or changing storybook theming.
---


# PaletteProvider implementation plan

## Executive summary

The package previously had palette data and per-widget `ContextStyleSet` helpers, but it did not have a global provider that could switch the whole design-system palette at once. Storybook transcript stories manually built `styleSet` props, so role accents could change, but global CSS variables such as `--mac-surface`, `--mac-text`, and `--mac-border` stayed on the default theme.

This change adds a `PaletteProvider` that maps the existing palette definitions onto CSS variables and exposes context/window plus transcript style sets through React context. Storybook wraps every story in the provider and adds a global toolbar control so designers can switch palettes across the whole rendered story surface.

## Problem statement

The transcript view appeared not to switch themes in Storybook because the palette controls only changed explicit `styleSet` props in selected stories. Component surfaces still depended on global `--rag-*` / `--mac-*` tokens from `theme.css`.

The desired behavior is one palette selector that updates both:

- design-system chrome and surfaces;
- context/transcript semantic style sets when widgets do not receive an explicit `styleSet` prop.

## Implemented solution

### Provider

`src/context/PaletteProvider.tsx` defines:

- `PaletteProvider`
- `usePalette()`
- `useContextStyleSet(explicit?)`
- `useTranscriptStyleSet(explicit?)`
- `paletteDefinitions`
- `paletteCssVars(palette)`

The provider maps `PaletteDefinition.colors` onto canonical CSS variables:

```ts
--rag-color-bg
--rag-color-surface
--rag-color-surface-muted
--rag-color-text
--rag-color-text-muted
--rag-color-border
--rag-color-border-strong
--rag-color-accent
--rag-color-success
--rag-color-warning
--rag-color-danger
```

It also sets compatibility variables that are not fully derived from canonical tokens:

```ts
--mac-bg-dark
--mac-stripe
--mac-text-inv
```

### Storybook wiring

`.storybook/preview.ts` now registers a global `palette` toolbar and wraps stories in `PaletteProvider`:

```ts
globalTypes: {
  palette: {
    defaultValue: "Dusty Magenta / Blue",
    toolbar: { icon: "paintbrush", items: contextPaletteOptions, dynamicTitle: true },
  },
}
```

### Transcript fallback

`TranscriptMessageCard` and `AnnotationNoteCard` now call `useTranscriptStyleSet(styleSet)` so explicit `styleSet` props still win, but Storybook/global provider palettes drive transcript styling when no explicit prop is passed.

The transcript WidgetRenderer stories were simplified to stop constructing local palette-specific style sets. They now rely on the global provider fallback, so the Storybook toolbar changes the transcript surface and transcript semantic colors together.

### Panel heading refinement

`Panel` remains a lightweight, non-boxy section-like layout. Its heading now:

- uses `--rag-font-role-panel-heading`;
- uses primary accent color via `--mac-accent`;
- has no underline/bottom border.

## Validation

Successful commands:

```bash
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build-storybook
```

## Follow-ups

- Consider using `useContextStyleSet` in context-window widgets when their `styleSet` prop is optional.
- Decide whether app runtime should expose a user-facing palette switcher, or whether this is Storybook/design-review-only for now.
