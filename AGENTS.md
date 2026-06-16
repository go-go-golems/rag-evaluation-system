# AGENTS.md — LLM Coding Instructions

Guidance for LLM agents working in this repository. Read this before making changes.

## Quick links

- **Design-system guidelines (MUST READ):** [`packages/rag-evaluation-site/GUIDELINES.md`](packages/rag-evaluation-site/GUIDELINES.md)
- **npm release playbook:** [`docs/guides/npm-release-playbook.md`](docs/guides/npm-release-playbook.md) — use this for `@go-go-golems/rag-evaluation-site` releases; publishing is through GitHub Actions Trusted Publishing, not local `npm publish`.
- **Biome config:** [`biome.json`](biome.json)
- **Git hooks:** [`lefthook.yml`](lefthook.yml)

## Formatting rules

This project uses **Biome v2** for all formatting. Always run after editing:

```bash
pnpm biome format --write .
pnpm biome lint --write .      # safe fixes only
```

Or the combined check:

```bash
pnpm biome check --write .
```

### Active style

| Setting | Value |
|---|---|
| Indent | **Tabs** (JS/TS/CSS/JSON) |
| Quotes | **Double** (`"not single"`) |
| Semicolons | **Always** |
| Trailing commas | **All** |
| Line width | **100** |
| CSS formatting | **Enabled** (tabs, 100 line width) |
| Import sorting | **Enabled** (auto-organized) |

### Do not

- Do not introduce Prettier, ESLint formatting rules, or `.editorconfig` — Biome owns formatting.
- Do not fight Biome's output. If it reformats something, that's the canonical style.
- Do not add `// biome-ignore` unless there's a documented reason in a comment.

## Design-system rules (summary)

Full rules in [`packages/rag-evaluation-site/GUIDELINES.md`](packages/rag-evaluation-site/GUIDELINES.md). Key points:

1. **Respect the layer hierarchy:** foundation → atoms → layout → molecules → organisms → widgets. Never import upward across layer boundaries.
2. **Typography:** Use `Text`, `Caption`, `CodeText`, `StatusText` and `--rag-font-role-*` tokens. Never write raw `font-size`/`font-weight`/`letter-spacing` unless adding a new token.
3. **Colors:** Use `--mac-*` and `--rag-*` theme tokens. Never hardcode hex colors when a token exists.
4. **CSS Modules only.** No CSS-in-JS, no `Box`/style-prop components, no utility classes.
5. **Components are API-free.** No RTK Query, no app store, no router, no backend imports in package components.
6. **Storybook required.** New public components need stories before they're part of the design system.
7. **`data-rag-*` attributes** on all public components for visual-diff extraction.
8. **Component folder layout:** `Component.tsx`, `Component.module.css`, `Component.stories.tsx`, `index.ts`.

## Project structure

```
packages/rag-evaluation-site/   # Reusable design-system package (no app deps)
web/                            # App shell — RTK, routing, backend-connected views
cmd/                            # Go backend services
pkg/                            # Go shared libraries
schema/                         # Protobuf / JSON schemas
ttmp/                           # docmgr ticket workspaces (not code)
```

## Release reminders

For `@go-go-golems/rag-evaluation-site` npm releases, follow [`docs/guides/npm-release-playbook.md`](docs/guides/npm-release-playbook.md). Do not publish from a local shell unless the playbook is explicitly changed; bump the package version, validate locally, push to `main`, and trigger `publish-npm.yml` with Trusted Publishing.

## Before committing

- `pnpm biome check .` passes clean
- `pnpm --dir packages/rag-evaluation-site typecheck` passes (for package changes)
- `go test ./...` passes (for Go changes)
- `lefthook` pre-commit hooks run automatically (Biome format + lint on staged TSX/CSS/JSON, Go lint + test on staged Go files)
