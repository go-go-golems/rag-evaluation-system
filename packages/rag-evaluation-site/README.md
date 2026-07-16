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

- foundation/atoms/layout nodes such as `Text`, `CodeText`, `ContextStyleSwatch`, `SectionBlock`, `SplitPane`, `SlideShell`, and `SidebarShell`;
- context diagram nodes such as `ContextBudgetBar`, `ContextStripDiagram`, `ContextTreemap`, and `ContextDiagramPanel`;
- transcript/comment nodes such as `TranscriptWorkspacePanel`, `AnnotationRailPanel`, and `AnchoredCommentRail`;
- course/handout nodes such as `CourseStudioShell`, `CourseSlidePanel`, `MarkdownArticle`, and `HandoutDocumentShell`.

Review package Storybook under `Widget IR/Renderer/...` for IR-authored examples grouped by Foundation and Atoms, Layout Recipes, Context Diagrams, Transcript and Notes, and Course Studio.

Goja authors use only `require("widget.dsl")`. Generic UI/data capabilities and semantic context, course, CMS, CRM, schedule, and time views are organized as typed namespaces on that module.

## Context and transcript style sets

Context-window diagrams and transcript widgets use a hard-cutover `styleKey + ContextStyleSet` contract. Snapshot parts carry a caller-defined `styleKey`; the companion `ContextStyleSet` supplies `styles` and `legend` entries for those keys. Do not send or document legacy context-window `kind`, `ContextPartKind`, `ContextKindSwatch`, `legendKinds`, or `legendMode` values.

```ts
const styleSet = createContextStyleSetFromPalette({
  palette: signalOrangeCyan,
  entries: [
    { id: 'prompt', label: 'Prompt', accent: 'b', pattern: 'checker' },
    { id: 'evidence', label: 'Evidence', accent: 'a', pattern: 'stipple' },
    { id: 'answer', label: 'Answer', accent: 'b', pattern: 'solid', solid: true },
    { id: 'free', label: 'Headroom', accent: 'grid', pattern: 'none', hidden: true },
  ],
});

const snapshot = {
  id: 'rag-window',
  title: 'RAG answer window',
  limit: 32000,
  parts: [
    { id: 'p', label: 'Prompt', styleKey: 'prompt', tokens: 1400 },
    { id: 'e', label: 'Evidence', styleKey: 'evidence', tokens: 9200 },
    { id: 'a', label: 'Draft', styleKey: 'answer', tokens: 1800 },
    { id: 'h', label: 'Headroom', styleKey: 'free', tokens: 19600 },
  ],
};
```

`ContextVisualStyle.labelColor` is part of the foreground contract. Components expose it as `--ctx-label` so labels, token chips, note chips, and selected states remain readable over palette fills. Transcript bodies and side-note bodies intentionally stay neutral by default; palette colors should decorate title bars, borders, glyphs, swatches, and compact chips.

Storybook exposes the four preferred palettes (`Dusty Magenta / Blue`, `Signal Orange / Cyan`, `Slate / Coral`, `Cobalt / Sand`) through a `palette` control on context and transcript stories. For iframe-only screenshots, Storybook args can be encoded in the URL, for example:

```text
/iframe.html?id=widget-ir-renderer-transcript-and-notes--message-card-states&viewMode=story&args=palette:'Signal+Orange+/+Cyan'
```

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

## Page keyboard shortcuts

`RagEvaluationSiteApp` activates optional page-level shortcut bindings returned in the Widget page envelope. A binding carries a stable id, logical key, optional modifiers, human-readable label, and the same `ActionSpec` used by visible controls:

```json
{
  "shortcuts": {
    "bindings": [
      {
        "id": "accept",
        "key": "y",
        "label": "Yes",
        "action": { "kind": "server", "name": "triage.accept" },
        "preventDefault": true,
        "allowRepeat": false
      }
    ]
  }
}
```

The default app:

- installs one page-owned `keydown` listener and dispatches through the existing action handler;
- ignores editable controls, active dialogs, IME composition, disallowed repeats, consumed events, and nested `data-rag-keyboard-scope` regions;
- adds `aria-keyshortcuts` only while shortcuts are enabled;
- renders `KeyboardShortcutHelp` with generated labels and chords;
- persists the user's enable/disable choice under `rag-evaluation-site:page-shortcuts-enabled`.

`DataTable` remains a nested keyboard scope: table navigation and `table.command(...)` take precedence while a row has focus. Page shortcuts are intended for page-level commands such as card triage, navigation, or global workflow actions.

Custom hosts that fetch a `WidgetPageResponse` but do not render `RagEvaluationSiteApp` can integrate the exported `usePageShortcuts` hook, `formatShortcutChord`, and `KeyboardShortcutHelp` component. They must provide an equivalent user-facing disable mechanism for unmodified character shortcuts; rendering `WidgetRenderer` alone does not install page behavior.

Review the `Component Library/Molecules/KeyboardShortcutHelp` Storybook story for the discoverability and preference UI.

## Build outputs

- `dist/`: npm package artifact.
- `app-dist/`: default Vite app build for Go embedding.
