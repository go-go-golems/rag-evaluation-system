# Changelog

## 2026-06-07

- Initial workspace created


## 2026-06-07

Created comprehensive design system reference for RAG Evaluation Design System. Two-part reference: Part I (Styling, Tokens, Typography, Spacing) and Part II (Widget IR, Component Catalogue, Actions, Recipes). Uploaded to reMarkable.


## 2026-06-07

Step 2: SSR design system reference renderer working. Two self-contained HTML files produced from actual React components: foundation-and-tokens.html (colours, typography, spacing, StatusText, Caption, Button variants) and widgets-and-composition.html (Panel, Stack, Inline, DashboardGrid, DataTable, MetadataGrid, FormRow, TabList, composition patterns). Uploaded to reMarkable. Commit: 95895aa.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/scripts/01-render-design-reference.cjs — Build script for SSR reference rendering
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/src/reference-pages/FoundationReference.tsx — Part 1 React component
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/src/reference-pages/WidgetReference.tsx — Part 2 React component


## 2026-06-07

Fixed SSR reference CSS extraction so CSS Modules are not tree-shaken; buttons now render as actual design-system buttons. Added css-visual-diff inspection verb and restored clean web typecheck.

### Related Files

- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/ttmp/2026/06/07/DESIGN-REF-001--rag-evaluation-design-system-reference/scripts/02-cssvd-inspect-reference.js — Ticket-local css-visual-diff evidence verb
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/src/reference-pages/WidgetReference.tsx — Typed widget examples and composition reference
- /home/manuel/workspaces/2026-05-27/rag-evaluation-system/2026-05-27--rag-evaluation-system/web/src/reference-pages/css-entry.ts — Keeps CSS Module bindings live so Vite emits matching standalone CSS

