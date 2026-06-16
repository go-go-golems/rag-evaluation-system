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
