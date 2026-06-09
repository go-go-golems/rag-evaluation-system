# @go-go-golems/rag-evaluation-site

React Widget IR renderer and default RAG evaluation site shell.

## Component ownership

This package is the canonical reusable UI layer. Read [`GUIDELINES.md`](./GUIDELINES.md) before changing package UI. Keep package components API-free and Storybook-first:

- `components/foundation`: text, captions, status, code, dividers, accessibility helpers.
- `components/atoms`: basic controls and semantic markers.
- `components/layout`: generic structure primitives such as shells, panes, stacks, panels, and grids.
- `components/molecules`: reusable data-display/content patterns with no backend hooks.
- `components/organisms`: feature panels with DTO-shaped props and callbacks.
- `widgets`: JSON-compatible Widget IR types and the React `WidgetRenderer`.

Storybook ownership follows the same boundary. Use `Design System/Layout/<Name>` for layout primitives, `Design System/Foundation/<Name>` for foundation, and `Component Library/Molecules/<Name>` / `Component Library/Organisms/<Name>` for reusable product components. Add or update stories before exposing a new public component.

When adding a component that should later be available from the Widget IR/Goja DSL, stabilize the React API first, keep props JSON-shaped where practical, add `data-rag-*` identity attributes, and only then extend `src/widgets/ir.ts` and `WidgetRenderer.tsx` with semantic component types.

## Widget IR coverage

`WidgetRenderer` supports both the original dashboard vocabulary and the newer semantic context-viewer vocabulary:

- foundation/atoms/layout nodes such as `Text`, `CodeText`, `ContextKindSwatch`, `SectionBlock`, `SplitPane`, `SlideShell`, and `SidebarShell`;
- context diagram nodes such as `ContextBudgetBar`, `ContextStripDiagram`, `ContextTreemap`, and `ContextDiagramPanel`;
- transcript/comment nodes such as `TranscriptWorkspacePanel`, `AnnotationRailPanel`, and `AnchoredCommentRail`;
- course/handout nodes such as `CourseStudioShell`, `CourseSlidePanel`, `MarkdownArticle`, and `HandoutDocumentShell`.

Review package Storybook under `Widget IR/Renderer/...` for IR-authored examples grouped by Foundation and Atoms, Layout Recipes, Context Diagrams, Transcript and Notes, and Course Studio.

Goja authors can use direct helpers from split modules such as `require("ui.dsl")`, `require("data.dsl")`, `require("context_window.dsl")`, and `require("course.dsl")`. Semantic recipes live under their owning domains, for example `contextWindow.recipes.contextDiagram`, `contextWindow.recipes.annotatedTranscript`, `course.recipes.courseStudio`, `course.recipes.courseSlide`, and `course.recipes.handout`.

## Install

```bash
pnpm add @go-go-golems/rag-evaluation-site react react-dom
```

## Publishing

The package is published from the generated `dist/` artifact, not directly from source. Build and smoke-test the artifact before publishing:

```bash
pnpm --dir packages/rag-evaluation-site build
pnpm --dir packages/rag-evaluation-site consumer:smoke
```

GitHub Actions workflow `.github/workflows/publish-npm.yml` publishes with npm Trusted Publishing. The workflow uses GitHub OIDC (`id-token: write`) and does not require `NODE_AUTH_TOKEN` after the npm package has been bootstrapped and trusted.

For the first publish of a brand-new package, create `@go-go-golems/rag-evaluation-site` manually with an npm account that has access to the `@go-go-golems` scope, then configure trusted publishing for repository `go-go-golems/rag-evaluation-system`, workflow `publish-npm.yml`, environment `npm-production`.

## Renderer usage

```tsx
import { WidgetRenderer, defaultWidgetRegistry, useWidgetPage } from '@go-go-golems/rag-evaluation-site';
import '@go-go-golems/rag-evaluation-site/styles.css';

export function Page() {
  const { page, loading, error } = useWidgetPage('/api/widget/pages/demo');
  if (loading) return <div>Loading…</div>;
  if (error || !page) return <div>Failed to load page</div>;
  return <WidgetRenderer node={page.root} registry={defaultWidgetRegistry} />;
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
