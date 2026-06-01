---
Title: Implementation Diary
Ticket: NPM-TRUST-001
Status: active
Topics:
    - npm
    - publishing
    - security
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/go-go-os-frontend/.github/workflows/publish-npm.yml
      Note: Current npm publish workflow audited for OIDC and Vault token usage
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/go-go-os-frontend/packages/os-repl/package.json
      Note: Pilot npm package metadata for @go-go-golems/os-repl
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/go-go-os-frontend/scripts/packages/package-sets.mjs
      Note: Public package inventory and publish sets
ExternalSources: []
Summary: Detailed diary for migrating npm packages to trusted publishing and token lockdown.
LastUpdated: 2026-06-01T15:25:00-04:00
WhatFor: Use this diary to resume the trusted publishing migration without losing exact commands, decisions, and blockers.
WhenToUse: Read before changing npm package settings, GitHub Actions publish workflows, or revoking npm tokens.
---


# Diary

## Goal

This diary records the attempt to migrate npm packages to trusted publishing and package-level token lockdown, starting with a single `go-go-os-frontend` package as a pilot.

## Step 1: Capture the source pack and create the ticket

I created the `NPM-TRUST-001` ticket to separate npm trusted publishing work from the RAG Evaluation System implementation tickets. The first task was to save the authoritative npm and GitHub source material locally so later package changes can be reviewed against stable references, not memory or search snippets.

The source pack establishes the migration rule we will follow: configure trusted publishing first, verify that OIDC publishing works, then switch the package's Publishing access setting to `Require two-factor authentication and disallow tokens`, and only then revoke obsolete automation tokens.

### Prompt Context

**User prompt (verbatim):** "but actually i see that go-go-os-frontend has trusted publishing set up say for os-repl. .../tmp/pi-clipboard-60fea408-6c61-44ce-9009-fe54c18dabd8.png but it says it's using 2fa or a granular access token with bypass 2fa enabled, but that was the one i now deleted. So OK, how do I properly set up trustede publishing on all my packages at this point, so it's all setup, and disallow tokens."

**Assistant interpretation:** Create a reliable plan for migrating npm packages away from long-lived tokens to npm trusted publishing, including package-level disallow-token settings.

**Inferred user intent:** The user wants all npm packages to publish through OIDC trusted publishing, not granular tokens, and wants package settings hardened so tokens cannot be used accidentally.

**Commit (code):** `1516cf3` — "docs: capture npm trusted publishing sources"

### What I did

- Created ticket `NPM-TRUST-001 — npm Trusted Publishing Migration and Token Lockdown`.
- Added docmgr vocabulary topics: `npm`, `publishing`, and `security`.
- Captured npm and GitHub source documents with Defuddle into `sources/`.
- Added `sources/00-source-pack-readme.md` summarizing the source set and migration conclusion.
- Added task phases for package audit, trusted publishing configuration, token lockdown, and playbook writing.
- Ran `docmgr doctor --ticket NPM-TRUST-001 --stale-after 30` and fixed source frontmatter issues.

### Why

- npm trusted publishing has several moving parts: npm package settings, GitHub Actions OIDC permissions, workflow filename matching, environment matching, provenance, and token revocation.
- The source pack makes the future package-by-package migration auditable.
- Defuddle captures keep the important docs local to the ticket.

### What worked

- Defuddle successfully captured the npm docs for trusted publishers, `npm trust`, package-level 2FA/disallow-token settings, provenance, and CI/CD token fallback.
- Defuddle also captured GitHub/npm changelog context about trusted publishing and token security changes.
- `docmgr doctor` passed after frontmatter was added to source captures.

### What didn't work

- The first `docmgr doctor` run failed because Defuddle output files did not include docmgr frontmatter. The exact class of failure was:

```text
Problem: frontmatter delimiters '---' not found
```

  I fixed this by adding frontmatter to each source file and renaming `sources/README.md` to `sources/00-source-pack-readme.md` to satisfy numeric-prefix hygiene.

### What I learned

- npm's official docs explicitly recommend `Require two-factor authentication and disallow tokens` after trusted publishers are configured.
- The docs state that disallow-token affects traditional token authentication, while trusted publishers continue to work because they use OIDC tokens.
- The `npm trust` command requires npm `11.10.0` or later and account-level 2FA for configuration.

### What was tricky to build

The tricky part was keeping source captures useful to docmgr. Raw Defuddle output is valuable, but docmgr expects Markdown files under the ticket to have valid frontmatter. The fix was to treat each captured source as a reference document with explicit `ExternalSources` metadata.

### What warrants a second pair of eyes

- Whether this ticket should eventually live in a more general release-engineering docs area rather than this repository's `ttmp` tree.
- Whether package migration should prefer npmjs.com UI for the first pilot or `npm trust` CLI from the beginning.

### What should be done in the future

- Audit `go-go-os-frontend` packages and workflows.
- Pilot one package, ideally `@go-go-golems/os-repl`, because the user already observed a trusted publishing setting there.
- Write the final repeatable migration playbook after the pilot.

### Code review instructions

- Review `sources/00-source-pack-readme.md` first.
- Then review the captured npm docs if the migration plan is questioned.
- Validate with:

```bash
docmgr doctor --ticket NPM-TRUST-001 --stale-after 30
```

### Technical details

The key source-backed migration order is:

1. Configure trusted publisher.
2. Ensure GitHub Actions has `permissions: id-token: write`.
3. Verify publish or staged publish works without token auth.
4. Set package Publishing access to `Require two-factor authentication and disallow tokens`.
5. Revoke obsolete automation tokens and remove secrets.

## Step 2: Pilot `@go-go-golems/os-repl` in `go-go-os-frontend`

I started the migration pilot on `@go-go-golems/os-repl`, one of the public workspace packages in `go-go-os-frontend`. The repository already has a `publish-npm` GitHub Actions workflow with OIDC permission and an `npm-production` environment, but it also still reads `NODE_AUTH_TOKEN` from Vault. That means the workflow is close to trusted publishing readiness, but the token path still exists and should not be considered fully locked down.

The pilot reached the npm authentication boundary. I could validate the intended trusted-publisher configuration with `npm trust github --dry-run`, but the actual registry mutation failed because this machine is not logged into npm. That is a useful stopping point: the next action requires an npm-authenticated maintainer session with 2FA, not more repository edits.

### Prompt Context

**User prompt (verbatim):** "ok, let's try it on one one of the go-go-os-frontend packages. Keep a detailed diary as you work."

**Assistant interpretation:** Use one `go-go-os-frontend` npm package as a pilot for trusted publishing/token-lockdown migration, record all commands and blockers, and avoid broad changes until the pilot is understood.

**Inferred user intent:** The user wants a safe, concrete trial on a real package before applying trusted publishing and token-disallow settings across every package.

**Commit (code):** N/A — no `go-go-os-frontend` code changes made in this step.

### What I did

- Located the likely canonical frontend repo at:
  - `/home/manuel/code/wesen/go-go-golems/go-go-os-frontend`
- Checked repository status:

```text
cd /home/manuel/code/wesen/go-go-golems/go-go-os-frontend && git status --short
```

  It returned no changes at that moment.

- Listed workflows:

```text
.github/workflows/launcher-ci.yml
.github/workflows/publish-github-package-canary.yml
.github/workflows/publish-npm.yml
```

- Listed workspace package files and selected the pilot package:
  - `packages/os-repl/package.json`
  - package name: `@go-go-golems/os-repl`
  - version: `0.1.6`
  - repository: `git+https://github.com/go-go-golems/go-go-os-frontend.git`
  - repository directory: `packages/os-repl`
- Inspected `.github/workflows/publish-npm.yml`.
- Confirmed the workflow already has:

```yaml
permissions:
  contents: read
  id-token: write
```

- Confirmed the workflow uses GitHub environment:

```yaml
environment: npm-production
```

- Confirmed the workflow still reads an npm token from Vault:

```yaml
- name: Read NPM token from Vault
  id: vault-npm
  uses: hashicorp/vault-action@v3
```

  and passes it as:

```yaml
NODE_AUTH_TOKEN: ${{ steps.vault-npm.outputs.NODE_AUTH_TOKEN }}
```

- Confirmed package metadata from the registry:

```bash
npm view @go-go-golems/os-repl version dist-tags repository --json
```

  Output:

```json
{
  "version": "0.1.6",
  "dist-tags": {
    "latest": "0.1.6"
  },
  "repository": {
    "type": "git",
    "url": "git+https://github.com/go-go-golems/go-go-os-frontend.git",
    "directory": "packages/os-repl"
  }
}
```

- Checked local npm authentication:

```bash
npm whoami 2>&1 || true
```

  Output:

```text
npm error code E401
npm error 401 Unauthorized - GET https://registry.npmjs.org/-/whoami
npm error A complete log of this run can be found in: /home/manuel/.npm/_logs/2026-06-01T19_21_52_212Z-debug-0.log
```

- Checked that modern npm can be invoked through `npx`:

```bash
npx -y npm@latest --version
```

  Output:

```text
11.16.0
```

- Tried to list current trust configuration:

```bash
npx -y npm@latest trust list @go-go-golems/os-repl --json 2>&1 || true
```

  Output:

```text
npm error code E401
npm error 401 Unauthorized - GET https://registry.npmjs.org/-/package/@go-go-golems%2fos-repl/trust - {"success":false,"error":"You must be logged in to publish packages."}
{
  "error": {
    "code": "E401",
    "summary": "401 Unauthorized - GET https://registry.npmjs.org/-/package/@go-go-golems%2fos-repl/trust - {\"success\":false,\"error\":\"You must be logged in to publish packages.\"}",
    "detail": ""
  }
}
```

- Validated the intended trusted publisher configuration with dry-run:

```bash
npx -y npm@latest trust github @go-go-golems/os-repl \
  --repo go-go-golems/go-go-os-frontend \
  --file publish-npm.yml \
  --env npm-production \
  --allow-publish \
  --dry-run \
  --json 2>&1 || true
```

  Output:

```json
{
  "package": "@go-go-golems/os-repl",
  "file": "publish-npm.yml",
  "repository": "go-go-golems/go-go-os-frontend",
  "environment": "npm-production",
  "permissions": [
    "createPackage"
  ]
}
```

- Tried the actual trust configuration command:

```bash
npx -y npm@latest trust github @go-go-golems/os-repl \
  --repo go-go-golems/go-go-os-frontend \
  --file publish-npm.yml \
  --env npm-production \
  --allow-publish \
  --yes \
  --json 2>&1 || true
```

  Output:

```text
{
  "package": "@go-go-golems/os-repl",
  "file": "publish-npm.yml",
  "repository": "go-go-golems/go-go-os-frontend",
  "environment": "npm-production",
  "permissions": [
    "createPackage"
  ]
}
npm error code E401
npm error 401 Unauthorized - POST https://registry.npmjs.org/-/package/@go-go-golems%2fos-repl/trust - {"success":false,"error":"You must be logged in to publish packages."}
{
  "error": {
    "code": "E401",
    "summary": "401 Unauthorized - POST https://registry.npmjs.org/-/package/@go-go-golems%2fos-repl/trust - {\"success\":false,\"error\":\"You must be logged in to publish packages.\"}",
    "detail": ""
  }
}
```

### Why

- `@go-go-golems/os-repl` is a good pilot because it is already public and the user specifically mentioned seeing trusted publishing state for `os-repl`.
- The workflow already has the two key GitHub-side attributes that npm trusted publishing must match: `publish-npm.yml` and `npm-production`.
- Trying `npm trust github --dry-run` first proves the desired configuration shape before attempting a privileged registry mutation.

### What worked

- The repository and package metadata are consistent with trusted publishing:
  - GitHub repository: `go-go-golems/go-go-os-frontend`
  - workflow filename: `publish-npm.yml`
  - environment: `npm-production`
  - npm package: `@go-go-golems/os-repl`
- The workflow already grants `id-token: write`, which is required for OIDC.
- `npx -y npm@latest` provides npm `11.16.0`, which is new enough for `npm trust` commands.
- The dry-run trusted publisher command produced the expected configuration payload.

### What didn't work

- Local npm is not authenticated:

```text
npm error code E401
npm error 401 Unauthorized - GET https://registry.npmjs.org/-/whoami
```

- Listing and creating trusted publisher configuration both failed with E401 because npm requires an authenticated maintainer session:

```text
npm error 401 Unauthorized - POST https://registry.npmjs.org/-/package/@go-go-golems%2fos-repl/trust - {"success":false,"error":"You must be logged in to publish packages."}
```

- The npm website opened to the login page for package access settings, so there is no browser session available in this agent context:

```text
Page URL: https://www.npmjs.com/login?next=%2Fpackage%2F%40go-go-golems%2Fos-repl%2Faccess
Page Title: npm | Sign In
```

### What I learned

- We can safely compute the exact trusted publisher tuple for the pilot package:
  - Package: `@go-go-golems/os-repl`
  - Provider: GitHub Actions
  - Repository: `go-go-golems/go-go-os-frontend`
  - Workflow filename: `publish-npm.yml`
  - Environment: `npm-production`
  - Permission: publish
- The existing workflow is not tokenless yet because it reads `NODE_AUTH_TOKEN` from Vault.
- Since the workflow can publish multiple package sets, removing the token from the workflow before all packages are trusted could break non-pilot package publishes unless the workflow is adjusted carefully.

### What was tricky to build

The tricky part is that there are two separate states to migrate: npm package settings and GitHub workflow authentication. The GitHub workflow is already OIDC-capable because it has `id-token: write`, but it still carries the old token path through Vault. The npm package may or may not already have a trusted publisher configured; without npm authentication I cannot list the current trust relationship.

The safe pilot sequence is therefore not "delete token wiring first." It is: authenticate to npm, verify or create the trusted publisher for `@go-go-golems/os-repl`, run one OIDC publish path for that package, then set Publishing access to disallow tokens. After that, remove or gate the Vault token path.

### What warrants a second pair of eyes

- The `npm trust github --dry-run` output used `permissions: ["createPackage"]` even though the command specified `--allow-publish`. Before doing the real change, verify whether this is npm's current JSON label for publish permission or whether `--allow-stage-publish` should be used instead for a stage-only hardening flow.
- Whether to keep a single multi-package `publish-npm.yml` workflow or add a temporary one-package OIDC pilot workflow for `@go-go-golems/os-repl`.
- Whether to remove the Vault token step immediately after the pilot package works, or only after all public packages have trusted publishers configured.

### What should be done in the future

- Authenticate npm locally or use npmjs.com in a human browser session.
- Run:

```bash
npx -y npm@latest trust list @go-go-golems/os-repl --json
```

- If no correct trusted publisher exists, run:

```bash
npx -y npm@latest trust github @go-go-golems/os-repl \
  --repo go-go-golems/go-go-os-frontend \
  --file publish-npm.yml \
  --env npm-production \
  --allow-publish \
  --yes
```

- Verify a real OIDC publish path for an unpublished version or a staged publish if stage permissions are chosen.
- Set `@go-go-golems/os-repl` package Publishing access to `Require two-factor authentication and disallow tokens`.
- Decide how to remove the Vault token path without breaking other packages.

### Code review instructions

- Start by reviewing `.github/workflows/publish-npm.yml` in `go-go-os-frontend`.
- Confirm that the trusted publisher tuple exactly matches npm package settings:
  - owner/repo: `go-go-golems/go-go-os-frontend`
  - workflow file: `publish-npm.yml`
  - environment: `npm-production`
- Check `packages/os-repl/package.json` for the package name and repository metadata.
- Do not remove token wiring from the shared workflow until either:
  - all packages in the workflow are trusted, or
  - the workflow has a safe tokenless pilot path for only `@go-go-golems/os-repl`.

### Technical details

The public packages in `go-go-os-frontend` are listed in `scripts/packages/package-sets.mjs`:

```text
packages/os-core
packages/os-repl
packages/os-chat
packages/os-scripting
packages/os-ui-cards
packages/os-confirm
packages/os-shell
packages/os-widgets
packages/os-kanban
```

The pilot package metadata is:

```text
package: @go-go-golems/os-repl
version: 0.1.6
repo: go-go-golems/go-go-os-frontend
workflow: publish-npm.yml
environment: npm-production
```
