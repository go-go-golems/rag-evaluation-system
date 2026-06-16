---
Title: Investigation diary
Ticket: RAGEVAL-006
Status: active
Topics:
    - react
    - css
    - packaging
    - storybook
    - vite
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: 2026-05-27--rag-evaluation-system/biome.json
      Note: Biome v2 config with formatter + linter + CSS support
    - Path: 2026-05-27--rag-evaluation-system/lefthook.yml
      Note: Added biome-format and biome-lint pre-commit hooks
    - Path: 2026-05-27--rag-evaluation-system/package.json
      Note: Root package.json with biome scripts and devDependency
    - Path: 2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/.storybook/main.ts
      Note: Added viteFinal with readable CSS module names for Storybook
    - Path: 2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/atoms/Button/Button.module.css
      Note: CSS module with root/normal/compact/primary/selected classes
    - Path: 2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/atoms/Button/Button.tsx
      Note: Component that triggered the CSS module naming investigation
    - Path: 2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/vite.config.ts
      Note: Added css.modules.generateScopedName for readable dev class names
ExternalSources: []
Summary: Chronological diary of RAGEVAL-006 polish work
LastUpdated: 2026-06-16T09:52:00Z
WhatFor: ""
WhenToUse: ""
---




# Diary

## Goal

Polish the React components, CSS modules, and packaging configuration in rag-evaluation-site so that development ergonomics (readable class names, consistent formatting) and production builds (hashed names, clean exports) are both solid.

## Step 1: Readable CSS Module Class Names in Dev and Storybook

The user was working on Button.tsx CSS and noticed that the generated HTML in Storybook showed hashed class names like `_root_1gwp3_1 _normal_1gwp3_41` instead of readable names. CSS Modules in Vite default to scoped hashed names, which is correct for production but terrible for debugging in dev tools and Storybook. The fix is to configure `generateScopedName` in Vite's CSS modules config.

### Prompt Context

**User prompt (verbatim):** "I'm working on the css for @./2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/atoms/Button/Button.tsx and want to a) have proper class names in the generated HTML say in http://localhost:6007/iframe.html?id=design-system-atoms-button--variants&viewMode=story (storybook in packages) as it currently says <button type="button" class="_root_1gwp3_1 _normal_1gwp3_41" data-rag-atom="Button">Default</button> (not sure how to do that if at all), I wonder where the"

**Assistant interpretation:** Make CSS Module class names readable (e.g. `Button_root` instead of `_root_1gwp3_1`) in Storybook and dev builds, while keeping hashed names for production.

**Inferred user intent:** Better developer experience when inspecting DOM in browser dev tools / Storybook; maintain production safety with hashed names.

**Commit (code):** bb2fdc0 — "feat(rag-evaluation-site): readable CSS Module class names in dev/Storybook"
**Commit (docs):** 902f9ad — "docs(ttmp): add RAGEVAL-006 ticket"

### What I did

- Added `css.modules.generateScopedName` to `vite.config.ts` with env-conditional logic: `[name]_[local]` in dev, `[hash:base64:5]` in production
- Added `viteFinal` override in `.storybook/main.ts` with `generateScopedName: '[name]_[local]'` so Storybook always shows readable names
- Both files were previously unchanged (no prior CSS module config existed)

### Why

CSS Modules hash class names by default to guarantee global uniqueness. This is the right behavior for production but makes debugging painful — you can't tell which component a class belongs to by looking at the DOM. The `[name]_[local]` pattern keeps scoping (the filename prefix prevents collisions) while being human-readable.

### What worked

- The Vite `css.modules.generateScopedName` option accepts the same format strings as PostCSS css-loader: `[name]` (filename without `.module.css`), `[local]` (original class name), `[hash:base64:N]`
- Storybook's `viteFinal` hook cleanly overrides the Vite config for Storybook-specific needs

### What didn't work

- `docmgr doc add --ticket RAGEVAL-006` failed with "ambiguous ticket index doc for RAGEVAL-006 (got 2)" — the ticket directory had both `index.md` and `README.md`, which confused the CLI's ticket directory lookup. Worked around by creating docs manually.

### What I learned

- Vite's CSS Modules config is at `css.modules.generateScopedName` (not a plugin option)
- Storybook `viteFinal` receives the full Vite config and can mutate `config.css.modules` directly
- The `[name]_[local]` pattern is sufficient for dev scoping because the filename prefix (e.g. `Button_`) already prevents collisions across components

### What was tricky to build

- The docmgr ambiguity error was unexpected. It seems the ticket directory lookup matches multiple files as "index" docs. Creating docs manually with correct frontmatter works as a workaround but bypasses docmgr validation.

### What warrants a second pair of eyes

- Verify that `[name]_[local]` doesn't cause class name collisions if two CSS modules in different directories share the same filename (e.g. `atoms/Button/Button.module.css` vs `molecules/Button/Button.module.css`). Currently not an issue in this codebase but could be in the future.
- Verify the Storybook restart actually produces `Button_root` / `Button_normal` class names.

### What should be done in the future

- Investigate the docmgr "ambiguous ticket index doc" error and whether it's a known issue
- Consider `[name]_[local]_[hash:base64:5]` as a compromise that's both readable and collision-proof

### Code review instructions

- Start with `packages/rag-evaluation-site/vite.config.ts` — look for the `css.modules.generateScopedName` addition
- Then `packages/rag-evaluation-site/.storybook/main.ts` — look for the `viteFinal` addition
- Verify by restarting Storybook and inspecting the Button component DOM

### Technical details

- Vite CSS Modules format tokens: `[name]` = filename sans `.module.css`, `[local]` = original class, `[hash:base64:N]` = content hash
- `process.env.NODE_ENV` is set by Vite itself: `'production'` during `vite build`, `'development'` during `vite dev`
- Storybook always runs in development mode, so `viteFinal` override is technically redundant but explicit and safer

## Step 2: Set up Biome for consistent TSX/CSS formatting

The user wanted consistent and legible formatting across TSX files and asked whether Biome can format CSS. Biome v2.5.0 supports CSS (including CSS Modules). We set it up at the project root level to cover both `web/` and `packages/rag-evaluation-site/`.

### Prompt Context

**User prompt (verbatim):** "ok, set up biome, standard setup that is widespread."

**Assistant interpretation:** Install Biome at the top level of the rag-evaluation-system monorepo with a standard/recommended configuration, run formatting across all front-end files.

**Inferred user intent:** One consistent formatter for all TSX/TS/CSS/JSON files, wired into git hooks, with standard defaults.

### What I did

- Created root `package.json` with `@biomejs/biome` v2.5.0 as devDependency and format/lint/check scripts
- Created `biome.json` with standard setup:
  - Tab indentation, 100 char line width
  - Double quotes, semicolons, trailing commas for JS/TS
  - CSS formatting enabled (tab indent, 100 line width)
  - `organizeImports` enabled
  - `linter` with `preset: "recommended"`
  - VCS integration (respects .gitignore)
  - `files.includes` with negation patterns (excludes node_modules, dist, ttmp, min files, pnpm-lock)
- Ran `biome migrate --write` to upgrade config to v2 schema (the initial config used v1 `ignore` key which doesn't exist in v2)
- Ran `biome format --write .` across the whole project — 517 files fixed, 1 straggler fixed on second pass
- Updated `lefthook.yml` to add `biome-format` and `biome-lint` pre-commit hooks for `*.{ts,tsx,js,jsx,css,json}` files

### Why

No formatter was configured in the project. The CSS in KeyValueStrip.module.css was minified to single lines, Button.module.css had inconsistent spacing. TSX files used single quotes inconsistently. A single tool (Biome) handles TSX, TS, CSS, and JSON — no need for Prettier + stylelint + eslint for formatting.

### What worked

- `biome migrate --write` automatically converted the v1-style `ignore` array to v2 `includes` with negation patterns (`!**/node_modules`)
- `organizeImports` was auto-migrated from `organizeImports.enabled` to `assist.actions.source.organizeImports: "on"`
- `recommended: true` was auto-migrated to `preset: "recommended"`
- Formatting applied cleanly: tabs, double quotes, semicolons, consistent spacing
- VCS integration means Biome respects `.gitignore` out of the box

### What didn't work

- Initial `biome.json` used v1 config keys (`files.ignore`, `organizeImports.enabled`, `rules.recommended: true`) which caused errors in Biome v2. `biome migrate --write` fixed all of them.
- One file needed a second `format --write` pass (line-width edge case on a function signature)

### What I learned

- Biome v2 has a different config schema from v1 — always run `biome migrate` after writing config by hand
- The v2 `files.includes` uses negation patterns (`!**/node_modules`) instead of v1's `files.ignore`
- 517 of 519 files had formatting differences — the codebase had zero prior formatting discipline

### What was tricky to build

- The v1→v2 config migration was not obvious. The error messages pointed at the `ignore` key but didn't say "use `includes` with negation instead." Running `biome migrate` was the discovery.

### What warrants a second pair of eyes

- Lint findings (2744 errors, 2563 warnings) are not yet addressed — formatting is clean but lint is noisy. The linter `preset: "recommended"` may be too aggressive for an existing codebase. Consider starting with `preset: "weak"` or disabling specific rules.
- The `biome-lint` lefthook hook will currently block commits on lint errors. May need to use `--no-errors-on-unmatched` or downgraded severity initially.

### What should be done in the future

- Address lint findings incrementally (maybe start with `noUnusedLocals`, `noUnusedImports` which are the bulk)
- Consider relaxing lint rules that conflict with existing patterns
- Add `biome check --write` to CI pipeline

### Code review instructions

- Start with `biome.json` at project root — verify config looks reasonable
- Then `lefthook.yml` — verify the two new pre-commit hooks
- Then `package.json` — verify scripts and devDependency
- Spot-check a few formatted files (e.g. `Button.tsx`, `Button.module.css`) for tab indentation and double quotes

### Technical details

- `biome format --write .` reformats in place; `biome format .` is check-only (CI-friendly)
- `biome check` runs both format + lint in one pass
- Lefthook `stage_fixed: true` re-stages files that were fixed by the hook
