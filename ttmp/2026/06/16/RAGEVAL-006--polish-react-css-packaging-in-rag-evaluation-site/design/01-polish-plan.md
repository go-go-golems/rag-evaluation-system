---
Title: 'Polish Plan: React, CSS, Packaging'
Ticket: RAGEVAL-006
Status: active
Topics:
    - react
    - css
    - packaging
    - storybook
    - vite
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: 2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/.storybook/main.ts
      Note: Storybook config with viteFinal override
    - Path: 2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/vite.config.ts
      Note: Primary config file for CSS module naming
ExternalSources: []
Summary: Plan for polishing React components, CSS modules, and packaging in rag-evaluation-site
LastUpdated: 2026-06-16T09:52:00Z
WhatFor: ""
WhenToUse: ""
---


# Polish Plan: React, CSS, Packaging

## Overview

This document tracks the polish work for the rag-evaluation-site package: improving CSS module readability, component consistency, and packaging configuration.

## Completed

### Step 1: Readable CSS Module Class Names

- Added `css.modules.generateScopedName` to `vite.config.ts`:
  - Dev: `[name]_[local]` → e.g. `Button_root`, `Button_normal`
  - Prod: `[hash:base64:5]` → short hashes for production
- Added `viteFinal` override in `.storybook/main.ts` with `[name]_[local]` for Storybook readability

## Open Items

- Audit all CSS modules for consistent formatting (expanded vs minified)
- Review Button padding values and alignment
- Evaluate whether `data-rag-atom` attributes should follow a naming convention
- Check that all component exports are clean and consistent
- Verify Storybook stories cover all variants/states
- Production build validation (hashed class names work correctly)
