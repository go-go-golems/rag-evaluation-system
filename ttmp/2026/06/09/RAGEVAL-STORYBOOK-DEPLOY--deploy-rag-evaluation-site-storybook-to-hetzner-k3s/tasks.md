# Tasks

## Phase 0 — Planning and ticket setup

- [x] Create deployment design guide for package-local Storybook
- [x] Record initial investigation diary and evidence
- [x] Upload initial design bundle to reMarkable
- [x] Expand this ticket task list into phased implementation tasks

## Phase 1 — Source repository static artifact contract

- [x] Add `Dockerfile.storybook-static` that packages `packages/rag-evaluation-site/storybook-static/` as `/site`
- [x] Add `.dockerignore` tuned for the Storybook static artifact build context
- [x] Update `.gitignore` so generated package Storybook output is not accidentally committed
- [x] Validate package-local typecheck and Storybook build
- [x] Validate local Docker artifact image contains `/site/index.html`
- [x] Record Phase 1 diary entry and commit focused source artifact changes

## Phase 2 — Source repository CI and GitOps handoff

- [x] Add `deploy/gitops-targets.json` with `patch_strategy: static-publisher-job`
- [x] Add `.github/workflows/publish-rag-evaluation-storybook.yml` caller workflow for infra-tooling
- [x] Validate workflow YAML and deployment target metadata shape locally
- [x] Record Phase 2 diary entry and commit focused source CI/GitOps handoff changes

## Phase 3 — K3s GitOps static-site manifests

- [x] Add `gitops/kustomize/rag-evaluation-storybook/kustomization.yaml`
- [x] Add static-site ServiceAccount, VaultConnection, VaultAuth, and VaultStaticSecret manifests
- [x] Add publisher Job manifest using a clearly marked initial `sha-0000000` placeholder
- [x] Add Ingress manifest for `rag-evaluation-storybook.yolo.scapegoat.dev`
- [x] Add Argo CD Application manifest in the `static-sites` project
- [x] Validate `kubectl kustomize gitops/kustomize/rag-evaluation-storybook`
- [x] Record Phase 3 diary entry and commit focused K3s static-site manifest changes

## Phase 4 — K3s Vault policy and role manifests

- [ ] Add Kubernetes Vault policy and role for the Storybook image-pull secret
- [ ] Add GitHub Actions Vault policy and role for the source workflow GitOps PR token
- [ ] Validate generated files for placeholder-free secret paths and expected role names
- [ ] Record Phase 4 diary entry and commit focused K3s Vault wiring changes

## Phase 5 — Documentation, validation, and reMarkable update

- [ ] Update design guide with implementation results, exact files created, validation outcomes, and commits
- [ ] Update diary into a playbook-quality chronological packaging narrative
- [ ] Run `docmgr doctor --ticket RAGEVAL-STORYBOOK-DEPLOY --stale-after 30`
- [ ] Upload updated guide bundle to reMarkable
- [ ] Record final changelog entry and mark ticket tasks complete
