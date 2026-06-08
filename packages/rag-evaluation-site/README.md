# @go-go-golems/rag-evaluation-site

React Widget IR renderer and default RAG evaluation site shell.

## Component ownership

This package is the canonical reusable UI layer. Keep package components API-free and Storybook-first:

- `components/foundation`: text, captions, status, code, dividers, accessibility helpers.
- `components/atoms`: basic controls and semantic markers.
- `components/layout`: generic structure primitives such as shells, panes, stacks, panels, and grids.
- `components/molecules`: reusable data-display/content patterns with no backend hooks.
- `components/organisms`: feature panels with DTO-shaped props and callbacks.
- `widgets`: JSON-compatible Widget IR types and the React `WidgetRenderer`.

Storybook ownership follows the same boundary. Use `Design System/Layout/<Name>` for layout primitives, `Design System/Foundation/<Name>` for foundation, and `Component Library/Molecules/<Name>` / `Component Library/Organisms/<Name>` for reusable product components. Add or update stories before exposing a new public component.

When adding a component that should later be available from the Widget IR/Goja DSL, stabilize the React API first, keep props JSON-shaped where practical, add `data-rag-*` identity attributes, and only then extend `src/widgets/ir.ts` and `WidgetRenderer.tsx` with semantic component types.

## Install

```bash
pnpm add @go-go-golems/rag-evaluation-site react react-dom
```

## Renderer usage

```tsx
import { WidgetRenderer, useWidgetPage } from '@go-go-golems/rag-evaluation-site';
import '@go-go-golems/rag-evaluation-site/styles.css';

export function Page() {
  const { page, loading, error } = useWidgetPage('/api/widget/pages/demo');
  if (loading) return <div>Loading…</div>;
  if (error || !page) return <div>Failed to load page</div>;
  return <WidgetRenderer node={page.root} />;
}
```

## Default app usage

```tsx
import { RagEvaluationSiteApp } from '@go-go-golems/rag-evaluation-site/app';
import '@go-go-golems/rag-evaluation-site/styles.css';

export function App() {
  return <RagEvaluationSiteApp apiBase="/api/widget" defaultPageId="index" />;
}
```

## Build outputs

- `dist/`: npm package artifact.
- `app-dist/`: default Vite app build for Go embedding.
