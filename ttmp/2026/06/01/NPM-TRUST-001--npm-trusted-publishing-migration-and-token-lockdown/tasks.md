# Tasks

## Phase 0 — Source capture

- [x] Create ticket workspace for npm trusted publishing migration and token lockdown.
- [x] Add npm/publishing/security topics to docmgr vocabulary.
- [x] Capture npm Trusted Publishers documentation with Defuddle.
- [x] Capture npm `npm trust` CLI documentation with Defuddle.
- [x] Capture npm 2FA and disallow-token package setting documentation with Defuddle.
- [x] Capture npm provenance documentation with Defuddle.
- [x] Capture npm CI/CD token fallback documentation with Defuddle.
- [x] Capture GitHub/npm trusted publishing and token-security changelog context with Defuddle.
- [x] Add source README with document purpose and working conclusion.

## Phase 1 — Audit current packages and workflows

- [x] Identify the public `go-go-os-frontend` npm packages listed by `scripts/packages/package-sets.mjs`.
- [x] Record the pilot package tuple for `@go-go-golems/os-repl`: repo `go-go-golems/go-go-os-frontend`, workflow `publish-npm.yml`, environment `npm-production`.
- [x] Inspect current `go-go-os-frontend` GitHub Actions workflow for `NODE_AUTH_TOKEN`, Vault token use, `--provenance`, and `id-token: write`.
- [ ] Identify packages currently configured with granular tokens or bypass-2FA token publishing.
- [ ] List current npm trusted publisher config for `@go-go-golems/os-repl` from an authenticated npm session.

## Phase 2 — Configure trusted publishing safely

- [ ] Upgrade local/admin npm CLI to a version supporting `npm trust` commands if CLI automation is used.
- [ ] Configure trusted publishers package-by-package, using either npmjs.com or `npm trust github`.
- [ ] Verify OIDC publish or staged publish works without npm tokens.
- [ ] Set package Publishing access to `Require two-factor authentication and disallow tokens` after verification.
- [ ] Revoke obsolete npm automation tokens and remove GitHub secrets.

## Phase 3 — Write playbook

- [ ] Write a repeatable migration playbook for all packages.
- [ ] Include a package inventory table.
- [ ] Include rollback and failure diagnosis notes for common OIDC mismatch errors.
