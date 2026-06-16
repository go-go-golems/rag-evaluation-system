---
Title: Changelog
Ticket: RAGEVAL-006
Status: active
DocType: changelog
---

# Changelog

- **2026-06-16** Step 1: Added `css.modules.generateScopedName` to vite.config.ts and Storybook main.ts for readable CSS Module class names in dev/Storybook, hashed names in production.
  - `packages/rag-evaluation-site/vite.config.ts`: Added env-conditional generateScopedName
  - `packages/rag-evaluation-site/.storybook/main.ts`: Added viteFinal with readable name pattern

- **2026-06-16** Step 2: Set up Biome v2.5.0 at project root with standard config (tabs, double quotes, semicolons, CSS formatting). Ran format --write across 519 files (517 fixed). Wired biome-format + biome-lint into lefthook pre-commit.
  - `biome.json`: Biome v2 config
  - `package.json`: Root package with biome scripts
  - `lefthook.yml`: Biome pre-commit hooks
