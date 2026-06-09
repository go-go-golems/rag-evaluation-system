---
Title: 'Intern Guide: Deploy rag-evaluation-site Storybook to Hetzner K3s'
Ticket: RAGEVAL-STORYBOOK-DEPLOY
Status: active
Topics:
    - storybook
    - deployment
    - kubernetes
    - static-sites
    - rag-evaluation
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/docs/static-site-packaging-and-gitops-playbook.md
      Note: Reference static-site deployment model
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/gitops/applications/rag-evaluation-storybook.yaml
      Note: Argo CD Application for the Storybook static site
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/dmeta-examples/publish-job.yaml
      Note: Reference publisher Job shape
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/rag-evaluation-storybook/publish-job.yaml
      Note: K3s publisher Job copied /site into static-sites-content PVC
    - Path: ../../../../../../../../../../code/wesen/2026-03-27--hetzner-k3s/vault/roles/github-actions/rag-evaluation-system-gitops-pr.json
      Note: Vault GitHub Actions role referenced by source workflow
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/infra-tooling/.github/workflows/publish-ghcr-image.yml
      Note: Reusable GHCR publish and GitOps PR workflow
    - Path: ../../../../../../../../../../code/wesen/go-go-golems/infra-tooling/actions/open-gitops-pr/action.yml
      Note: GitOps patch action supporting static-publisher-job
    - Path: 2026-05-27--rag-evaluation-system/.github/workflows/publish-rag-evaluation-storybook.yml
      Note: Source CI workflow that builds Storybook
    - Path: 2026-05-27--rag-evaluation-system/Dockerfile.storybook-static
      Note: Static artifact image contract for Storybook /site packaging
    - Path: 2026-05-27--rag-evaluation-system/deploy/gitops-targets.json
      Note: Infra-tooling target metadata for static-publisher-job patching
    - Path: 2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/.storybook/main.ts
      Note: Package-local Storybook story glob and framework
    - Path: 2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/package.json
      Note: Package scripts and Storybook build contract
ExternalSources: []
Summary: Evidence-backed implementation guide for publishing packages/rag-evaluation-site Storybook as a static site through GHCR, infra-tooling, GitOps, Argo CD, and the Hetzner K3s static-sites-host.
LastUpdated: 2026-06-09T20:59:15-04:00
WhatFor: Use this when implementing or reviewing the Storybook deployment for @go-go-golems/rag-evaluation-site.
WhenToUse: Before adding Docker, GitHub Actions, deploy target metadata, K3s manifests, Vault roles, or rollout validation for the package Storybook.
---



# Intern Guide: Deploy `rag-evaluation-site` Storybook to Hetzner K3s

## Executive summary

This guide explains how to deploy the **package-local Storybook** for `@go-go-golems/rag-evaluation-site` to the Hetzner K3s cluster. The important source directory is:

```text
/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/
```

Do **not** deploy a parent Storybook from the repository root or from `web/`. The package Storybook is the canonical reusable design-system review surface. Its Storybook configuration reads stories from `packages/rag-evaluation-site/src/**/*.stories.tsx`, and the package README says this package owns foundation, atom, layout, molecule, organism, and Widget IR stories.

The recommended deployment path is the cluster's existing static-site release model:

```text
rag-evaluation-system source repo
  -> build packages/rag-evaluation-site Storybook
  -> package Storybook static output into a GHCR artifact image containing /site
  -> infra-tooling reusable workflow opens a GitOps PR
  -> K3s GitOps repo updates gitops/kustomize/rag-evaluation-storybook/publish-job.yaml
  -> Argo CD runs a one-shot publisher Job
  -> publisher Job copies /site into static-sites-content PVC
  -> static-sites-host Caddy serves https://rag-evaluation-storybook.yolo.scapegoat.dev/
```

This guide is written for a new intern. It explains the relevant moving parts, why each exists, which files to change, the exact API contracts between repos, implementation sketches, validation commands, and failure modes.

## Problem statement and scope

The RAG Evaluation System has two front-end surfaces:

- product/application UI under `web/`;
- reusable design-system package UI under `packages/rag-evaluation-site/`.

The requested deployment is specifically for the **Storybook in `packages/rag-evaluation-site/`**. That Storybook is useful because it lets reviewers inspect reusable components without running the whole RAG backend, Goja runtime, or widget-page APIs.

The deployment should produce a stable public HTTPS URL backed by the existing Hetzner K3s GitOps infrastructure. The design should not introduce a long-running application server for Storybook. Storybook's production build is static HTML/CSS/JS, so it fits the cluster's shared `static-sites-host` model.

### In scope

- Source repository additions in `2026-05-27--rag-evaluation-system`:
  - `Dockerfile.storybook-static` or equivalent artifact Dockerfile.
  - `.dockerignore` adjustments.
  - `.github/workflows/publish-rag-evaluation-storybook.yml` caller workflow.
  - `deploy/gitops-targets.json` target metadata.
- GitOps repository additions in `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/`:
  - `gitops/kustomize/rag-evaluation-storybook/` static-site package.
  - `gitops/applications/rag-evaluation-storybook.yaml` Argo CD Application.
  - Vault Kubernetes policy/role for the image-pull secret.
  - Vault GitHub Actions policy/role for the source workflow's GitOps PR credential.
- Use of `/home/manuel/code/wesen/go-go-golems/infra-tooling/` reusable workflow and `open-gitops-pr` action.
- Operator commands for bootstrapping Vault secrets, applying the first Argo Application, and validating rollout.

### Out of scope

- Changing the actual React components.
- Reorganizing Storybook stories.
- Publishing the npm package.
- Deploying the `web/` app or default `app-dist/` Vite app.
- Building a custom Storybook Node server.

## Current-state analysis

### 1. The target package is a reusable UI package, not the parent app

The package metadata identifies the deploy target as `@go-go-golems/rag-evaluation-site` and describes it as a React Widget IR renderer and default RAG evaluation site shell (`packages/rag-evaluation-site/package.json:2-7`). It is public, ESM-based, and has peer dependencies on React/React DOM (`package.json:4-5`, `package.json:63-65`).

The package exposes reusable surfaces through `exports`, including the main renderer API, IR types, the default app entry, and CSS files (`package.json:37-51`). The top-level TypeScript entry imports package CSS and re-exports widgets, hooks, components, and context (`src/index.ts:1-6`).

The README makes the ownership boundary explicit:

- `components/foundation`: text, captions, code, status, accessibility helpers.
- `components/atoms`: controls and semantic markers.
- `components/layout`: generic structure primitives.
- `components/molecules`: reusable data/content display patterns.
- `components/organisms`: feature panels with DTO-shaped props/callbacks.
- `widgets`: JSON-compatible Widget IR types and `WidgetRenderer`.

This matters because the deployed Storybook should be the review surface for these package-owned components, not for app containers in `web/`.

### 2. The Storybook configuration is package-local

The package Storybook lives under:

```text
packages/rag-evaluation-site/.storybook/
```

The Storybook config says:

```ts
stories: ['../src/**/*.stories.@(ts|tsx)']
framework: { name: '@storybook/react-vite' }
```

Evidence:

- `packages/rag-evaluation-site/.storybook/main.ts:3-9` declares the config.
- `packages/rag-evaluation-site/package.json:60-61` declares `storybook` and `build-storybook` scripts.
- `packages/rag-evaluation-site/package.json:67-74` depends on `@storybook/react-vite`, `storybook`, Vite, TypeScript, and React plugin packages.

The source tree already contains many package-owned `*.stories.tsx` files under `src/components/...`. That matches the config and makes `pnpm --dir packages/rag-evaluation-site build-storybook` the correct build command.

### 3. The package also has non-Storybook build outputs; do not confuse them

The package has multiple build products:

- `dist/`: npm package artifact.
- `app-dist/`: default Vite app build for Go embedding.
- Storybook static output: default `storybook-static/` unless `storybook build -o <dir>` is used.

Evidence:

- `package.json:55` builds `dist` via `tsc`, `vite build`, CSS copy, and dist preparation.
- `package.json:56` builds the default app with `vite.app.config.ts`.
- `vite.config.ts:6-21` is the library build configuration and emits into `dist`.
- `vite.app.config.ts:6-10` emits the default app into `app-dist`.
- `package.json:61` defines `build-storybook` separately.

For the requested deployment, the artifact must copy the Storybook static output, not `dist/` and not `app-dist/`.

### 4. The Widget IR renderer is part of why Storybook matters

`WidgetRenderer` takes a JSON-compatible `WidgetNode`, a registry, and an optional action handler. It renders text, element, and component nodes and emits an `UnknownWidget` error callout if a component type is missing.

Evidence:

- `WidgetRendererProps` define `node`, `registry`, and `onAction` (`src/widgets/WidgetRenderer.tsx:7-11`).
- `WidgetRenderer` creates a render context and delegates to `renderWidgetNode` (`WidgetRenderer.tsx:13-16`).
- `renderWidgetNode` switches on `text`, `element`, and `component` kinds (`WidgetRenderer.tsx:29-40`).
- Component nodes are looked up in the registry (`WidgetRenderer.tsx:48-53`).
- Unknown nodes produce an explicit error callout (`WidgetRenderer.tsx:67-73`).

A deployed Storybook lets design-system reviewers inspect both direct React components and Widget IR-authored examples without needing to run the backend that would normally produce widget pages.

### 5. The cluster already has a static-site hosting model

The Hetzner K3s repo contains a `static-sites-host` service that serves content from a shared PVC. Its Caddy config uses the request host to choose the directory:

```text
root * /srv/sites/{host}/current
try_files {path} {path}/ /index.html
file_server
```

Evidence:

- `gitops/kustomize/static-sites-host/caddy-config.yaml:16` maps hostnames to `/srv/sites/{host}/current`.
- `caddy-config.yaml:18-20` sets long cache headers for immutable assets and no-cache for other responses.
- `caddy-config.yaml:22-23` serves SPA-style fallback to `/index.html` and files from disk.

This is a good fit for Storybook because Storybook's production output is static and has an `index.html` root plus asset files.

### 6. Existing static sites use one-shot publisher Jobs

The `dmeta-examples` static site is the reference implementation. Its publisher Job:

- has a name containing an immutable `sha-*` release token;
- pulls a GHCR artifact image;
- copies `/site/.` into `/srv/sites/<host>/releases/<release>`;
- atomically repoints `/srv/sites/<host>/current` to the new release;
- mounts the shared `static-sites-content` PVC.

Evidence:

- Job name and release label: `publish-job.yaml:4`, `publish-job.yaml:12-13`.
- GHCR image: `publish-job.yaml:28-30`.
- host/release/base/target variables: `publish-job.yaml:35-40`.
- copy and symlink update: `publish-job.yaml:42-47`.
- shared PVC mount: `publish-job.yaml:48-54`.

The same pattern should be reused for `rag-evaluation-storybook`.

### 7. Argo CD Application objects should use the `static-sites` project

The static-site playbook says static sites are managed in the `static-sites` Argo CD project. The current `dmeta-examples` Application follows that rule:

- `spec.project: static-sites` (`gitops/applications/dmeta-examples.yaml:20`).
- source path points into `gitops/kustomize/dmeta-examples` (`dmeta-examples.yaml:21-24`).
- destination namespace is `static-sites` (`dmeta-examples.yaml:25-27`).
- automated prune/self-heal is enabled (`dmeta-examples.yaml:28-31`).

Use the same shape for the Storybook Application.

### 8. Infra-tooling already provides the image-publish and GitOps PR machinery

The reusable workflow `/home/manuel/code/wesen/go-go-golems/infra-tooling/.github/workflows/publish-ghcr-image.yml` accepts caller inputs for the Dockerfile, build context, test command, image name, target config, Vault role, and Vault secret path.

Evidence:

- workflow inputs include `dockerfile`, `build_context`, `test_command`, `image_name`, and `gitops_target_config` (`publish-ghcr-image.yml:6-37`).
- inputs include `push_image`, `open_gitops_pr`, `tooling_repository`, and `tooling_ref` (`publish-ghcr-image.yml:46-61`).
- Vault inputs include `gitops_pr_token_source`, `vault_addr`, `vault_auth_path`, `vault_role`, and `vault_secret_path` (`publish-ghcr-image.yml:62-89`).
- permissions include `packages: write` and `id-token: write` (`publish-ghcr-image.yml:114-118`).
- the workflow runs the caller-provided test command (`publish-ghcr-image.yml:139-141`).
- it computes immutable `sha-${GITHUB_SHA::7}` image tags (`publish-ghcr-image.yml:143-156`).
- it builds and pushes with Docker Buildx (`publish-ghcr-image.yml:181-191`).
- it validates Vault/GitOps credential configuration (`publish-ghcr-image.yml:215-254`).
- it reads a GitOps PR token from Vault when `gitops_pr_token_source == vault` (`publish-ghcr-image.yml:256-268`).
- it invokes `./.infra-tooling/actions/open-gitops-pr` with config and image (`publish-ghcr-image.yml:311-322`).

The `open-gitops-pr` action declares that targets may set `patch_strategy=static-publisher-job`, which is exactly the strategy needed for one-shot static-site publisher Jobs (`actions/open-gitops-pr/action.yml:3-6`).

## Proposed architecture

### Architecture diagram

```text
┌────────────────────────────────────────────────────────────────────┐
│ Source repo: 2026-05-27--rag-evaluation-system                     │
│                                                                    │
│ packages/rag-evaluation-site/                                      │
│   .storybook/main.ts  -> selects ../src/**/*.stories.tsx           │
│   package.json        -> build-storybook script                    │
│   src/components/...  -> package-owned design-system stories       │
│                                                                    │
│ Dockerfile.storybook-static                                        │
│   COPY packages/rag-evaluation-site/storybook-static/ /site/       │
│                                                                    │
│ .github/workflows/publish-rag-evaluation-storybook.yml             │
│   calls infra-tooling publish-ghcr-image.yml                       │
│                                                                    │
│ deploy/gitops-targets.json                                         │
│   points to K3s publish-job.yaml with static-publisher-job         │
└───────────────────────────────┬────────────────────────────────────┘
                                │ GHCR image: sha-<commit>
                                ▼
┌────────────────────────────────────────────────────────────────────┐
│ infra-tooling reusable workflow/action                             │
│                                                                    │
│ publish-ghcr-image.yml                                             │
│   - runs caller test/build command                                 │
│   - builds artifact image                                          │
│   - pushes ghcr.io/go-go-golems/rag-evaluation-storybook:sha-*     │
│   - obtains GitOps PR token from Vault via GitHub OIDC             │
│   - opens PR against wesen/2026-03-27--hetzner-k3s                 │
│                                                                    │
│ open-gitops-pr action                                              │
│   - patches publish-job.yaml                                       │
│   - rewrites all sha-* release tokens                              │
└───────────────────────────────┬────────────────────────────────────┘
                                │ GitOps PR / merge
                                ▼
┌────────────────────────────────────────────────────────────────────┐
│ GitOps repo: /home/manuel/code/wesen/2026-03-27--hetzner-k3s       │
│                                                                    │
│ gitops/applications/rag-evaluation-storybook.yaml                  │
│ gitops/kustomize/rag-evaluation-storybook/                         │
│   serviceaccount.yaml                                              │
│   vault-connection.yaml                                            │
│   vault-auth.yaml                                                  │
│   vault-static-secret-image-pull.yaml                              │
│   publish-job.yaml                                                 │
│   ingress.yaml                                                     │
│ vault/policies/... and vault/roles/...                             │
└───────────────────────────────┬────────────────────────────────────┘
                                │ Argo CD sync
                                ▼
┌────────────────────────────────────────────────────────────────────┐
│ K3s cluster                                                        │
│                                                                    │
│ namespace: static-sites                                            │
│   Job publish-rag-evaluation-storybook-sha-*                       │
│     image contains /site                                           │
│     copies /site to static-sites-content PVC                       │
│                                                                    │
│   Service static-sites-host                                        │
│   Ingress rag-evaluation-storybook.yolo.scapegoat.dev              │
│   Caddy root /srv/sites/{host}/current                             │
└────────────────────────────────────────────────────────────────────┘
```

### Runtime request flow

```text
browser GET https://rag-evaluation-storybook.yolo.scapegoat.dev/
  -> Traefik Ingress
  -> static-sites-host Service
  -> Caddy container
  -> root /srv/sites/rag-evaluation-storybook.yolo.scapegoat.dev/current
  -> index.html + Storybook assets
```

### Release flow

```text
developer merges source change to main
  -> GitHub Actions builds package Storybook
  -> Docker image contains /site/index.html
  -> image is pushed as ghcr.io/go-go-golems/rag-evaluation-storybook:sha-abc1234
  -> infra-tooling opens GitOps PR changing publish-job.yaml
  -> reviewer merges GitOps PR
  -> Argo CD replaces publisher Job
  -> publisher Job copies /site to releases/sha-abc1234
  -> current symlink moves to releases/sha-abc1234
```

## File-level implementation guide

### Phase 1: Verify the source package build locally

Start in the source repo root:

```bash
cd /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system
corepack enable
pnpm install --frozen-lockfile
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build-storybook
find packages/rag-evaluation-site/storybook-static -maxdepth 2 -type f | sort | head -50
test -f packages/rag-evaluation-site/storybook-static/index.html
```

If `pnpm install --frozen-lockfile` fails because the lockfile is intentionally package-local, use the repo's existing installation convention. The important invariant is that `build-storybook` runs in `packages/rag-evaluation-site`, not in the repository parent.

### Phase 2: Add the static artifact Dockerfile

Create this file in the source repo root:

```text
/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/Dockerfile.storybook-static
```

Suggested content:

```dockerfile
# syntax=docker/dockerfile:1

FROM alpine:3.20

LABEL org.opencontainers.image.title="rag-evaluation-storybook"
LABEL org.opencontainers.image.description="Static Storybook artifact for @go-go-golems/rag-evaluation-site"
LABEL org.opencontainers.image.source="https://github.com/go-go-golems/rag-evaluation-system"

WORKDIR /
COPY packages/rag-evaluation-site/storybook-static/ /site/

RUN find /site -type f | sort > /site-manifest.txt \
  && test -f /site/index.html
```

This image is an artifact carrier. It does not run a server. The cluster's publisher Job mounts the image filesystem, copies `/site` into the shared static-sites PVC, and exits.

Local validation:

```bash
cd /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system
pnpm --dir packages/rag-evaluation-site build-storybook
docker build -f Dockerfile.storybook-static -t rag-evaluation-storybook:test .
docker run --rm rag-evaluation-storybook:test sh -c 'test -f /site/index.html && find /site -maxdepth 2 -type f | sort | head -50'
```

### Phase 3: Add or update `.dockerignore`

The Docker build context should exclude large irrelevant files but **must not exclude** the built Storybook directory. If `.dockerignore` already exists, adjust it carefully.

Suggested content:

```gitignore
.git
.github
.idea
.obsidian
node_modules
web/node_modules
packages/*/node_modules
.dist
build
coverage
**/.DS_Store

# Do not ignore packages/rag-evaluation-site/storybook-static;
# Dockerfile.storybook-static copies it into /site.
```

If the existing `.dockerignore` has a broad `storybook-static` or `packages/**/storybook-static` ignore rule, remove that rule or add a negation:

```gitignore
!packages/rag-evaluation-site/storybook-static/**
```

### Phase 4: Add source-repo deployment target metadata

Create:

```text
/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/deploy/gitops-targets.json
```

Suggested content:

```json
{
  "targets": [
    {
      "name": "rag-evaluation-storybook-prod",
      "gitops_repo": "wesen/2026-03-27--hetzner-k3s",
      "gitops_branch": "main",
      "manifest_path": "gitops/kustomize/rag-evaluation-storybook/publish-job.yaml",
      "container_name": "publish",
      "patch_strategy": "static-publisher-job"
    }
  ]
}
```

The `static-publisher-job` strategy is required because a publisher Job contains the release token in multiple places:

- Job metadata name.
- `static.wesen.dev/release` label.
- container image tag.
- shell `release="sha-*"` variable.

Kubernetes Job pod templates are immutable. Updating only `image:` can leave Argo CD trying to mutate an existing Job. The static publisher strategy rewrites all release tokens together.

### Phase 5: Add the source GitHub Actions caller workflow

Create:

```text
/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/.github/workflows/publish-rag-evaluation-storybook.yml
```

Suggested content:

```yaml
name: publish-rag-evaluation-storybook

on:
  pull_request:
    paths:
      - Dockerfile.storybook-static
      - .dockerignore
      - .github/workflows/publish-rag-evaluation-storybook.yml
      - deploy/gitops-targets.json
      - pnpm-lock.yaml
      - pnpm-workspace.yaml
      - packages/rag-evaluation-site/**
  push:
    branches:
      - main
    paths:
      - Dockerfile.storybook-static
      - .dockerignore
      - .github/workflows/publish-rag-evaluation-storybook.yml
      - deploy/gitops-targets.json
      - pnpm-lock.yaml
      - pnpm-workspace.yaml
      - packages/rag-evaluation-site/**
  workflow_dispatch:

permissions:
  contents: read
  packages: write
  pull-requests: write
  id-token: write

concurrency:
  group: publish-rag-evaluation-storybook-${{ github.ref }}
  cancel-in-progress: true

jobs:
  release:
    uses: go-go-golems/infra-tooling/.github/workflows/publish-ghcr-image.yml@main
    secrets: inherit
    with:
      dockerfile: ./Dockerfile.storybook-static
      build_context: .
      go_version_file: go.mod
      go_cache_dependency_path: go.sum
      test_command: |
        corepack enable
        pnpm install --frozen-lockfile
        pnpm --dir packages/rag-evaluation-site typecheck
        pnpm --dir packages/rag-evaluation-site build-storybook
        test -f packages/rag-evaluation-site/storybook-static/index.html
      image_name: ghcr.io/go-go-golems/rag-evaluation-storybook
      gitops_target_config: deploy/gitops-targets.json
      push_image: ${{ github.event_name != 'pull_request' }}
      open_gitops_pr: ${{ github.event_name != 'pull_request' && github.ref == 'refs/heads/main' }}
      gitops_pr_token_source: vault
      vault_role: rag-evaluation-system-gitops-pr
      vault_secret_path: kv/data/ci/github/rag-evaluation-system/gitops-pr-token
      tooling_repository: go-go-golems/infra-tooling
      tooling_ref: main
```

Notes for the intern:

- The reusable workflow runs `test_command` before Docker build. That is why the Dockerfile can copy `packages/rag-evaluation-site/storybook-static/`.
- `push_image` is disabled on PRs so PR validation can build without publishing.
- `open_gitops_pr` is enabled only on pushes to `main`.
- `id-token: write` is required because the reusable workflow exchanges the GitHub Actions OIDC token for a short-lived Vault token.
- If the source repo does not have `go.mod` in the GitHub checkout, set `go_version_file` to an existing file or update the reusable workflow path. In the current observed repo, `go.mod` exists at the source repo root.

### Phase 6: Add K3s GitOps manifests

Create this directory in the GitOps repo:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/rag-evaluation-storybook/
```

#### `kustomization.yaml`

```yaml
apiVersion: kustomize.config.k8s.io/v1beta1
kind: Kustomization
namespace: static-sites
resources:
  - serviceaccount.yaml
  - vault-connection.yaml
  - vault-auth.yaml
  - vault-static-secret-image-pull.yaml
  - publish-job.yaml
  - ingress.yaml
labels:
  - pairs:
      app.kubernetes.io/part-of: static-sites
      app.kubernetes.io/name: rag-evaluation-storybook
```

#### `serviceaccount.yaml`

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: rag-evaluation-storybook
  annotations:
    argocd.argoproj.io/sync-wave: "-2"
imagePullSecrets:
  - name: rag-evaluation-storybook-ghcr-pull
```

#### `vault-connection.yaml`

```yaml
apiVersion: secrets.hashicorp.com/v1beta1
kind: VaultConnection
metadata:
  name: vault
  annotations:
    argocd.argoproj.io/sync-wave: "-2"
spec:
  address: http://vault.vault.svc.cluster.local:8200
  skipTLSVerify: true
```

#### `vault-auth.yaml`

```yaml
apiVersion: secrets.hashicorp.com/v1beta1
kind: VaultAuth
metadata:
  name: rag-evaluation-storybook
  annotations:
    argocd.argoproj.io/sync-wave: "-2"
spec:
  vaultConnectionRef: vault
  method: kubernetes
  mount: kubernetes
  kubernetes:
    role: rag-evaluation-storybook
    serviceAccount: rag-evaluation-storybook
```

#### `vault-static-secret-image-pull.yaml`

```yaml
apiVersion: secrets.hashicorp.com/v1beta1
kind: VaultStaticSecret
metadata:
  name: rag-evaluation-storybook-ghcr-pull
  annotations:
    argocd.argoproj.io/sync-wave: "-1"
spec:
  vaultAuthRef: rag-evaluation-storybook
  mount: kv
  type: kv-v2
  path: apps/rag-evaluation-storybook/prod/image-pull
  refreshAfter: 30s
  destination:
    name: rag-evaluation-storybook-ghcr-pull
    create: true
    overwrite: true
    type: kubernetes.io/dockerconfigjson
    transformation:
      excludes:
        - ".*"
      templates:
        .dockerconfigjson:
          text: |
            {"auths":{"{{ .Secrets.server }}":{"username":"{{ .Secrets.username }}","password":"{{ .Secrets.password }}","auth":"{{ .Secrets.auth }}"}}}
```

#### `publish-job.yaml`

Use an initial `sha-*` tag that exists in GHCR. During first implementation, this might be a manually built/pushed test image or the first workflow-published image.

```yaml
apiVersion: batch/v1
kind: Job
metadata:
  name: publish-rag-evaluation-storybook-sha-0000000
  annotations:
    argocd.argoproj.io/sync-wave: "1"
    argocd.argoproj.io/compare-options: IgnoreExtraneous
    argocd.argoproj.io/sync-options: Replace=true
  labels:
    app.kubernetes.io/name: rag-evaluation-storybook
    app.kubernetes.io/component: publisher
    static.wesen.dev/host: rag-evaluation-storybook.yolo.scapegoat.dev
    static.wesen.dev/release: sha-0000000
spec:
  backoffLimit: 2
  template:
    metadata:
      labels:
        app.kubernetes.io/name: rag-evaluation-storybook
        app.kubernetes.io/component: publisher
    spec:
      restartPolicy: OnFailure
      enableServiceLinks: false
      serviceAccountName: rag-evaluation-storybook
      imagePullSecrets:
        - name: rag-evaluation-storybook-ghcr-pull
      containers:
        - name: publish
          image: ghcr.io/go-go-golems/rag-evaluation-storybook:sha-0000000
          imagePullPolicy: IfNotPresent
          command:
            - sh
            - -c
            - |
              set -eu
              host="rag-evaluation-storybook.yolo.scapegoat.dev"
              release="sha-0000000"
              base="/srv/sites/${host}"
              target="${base}/releases/${release}"
              tmp="${target}.tmp"

              rm -rf "${tmp}" "${target}"
              mkdir -p "${tmp}"
              cp -a /site/. "${tmp}/"
              mv "${tmp}" "${target}"
              ln -sfn "releases/${release}" "${base}/current"
              find "${target}" -maxdepth 3 -type f | sort | head -100
          volumeMounts:
            - name: static-sites-content
              mountPath: /srv/sites
      volumes:
        - name: static-sites-content
          persistentVolumeClaim:
            claimName: static-sites-content
```

The initial `sha-0000000` is a placeholder in this guide. Do not merge a manifest with a placeholder image unless you intentionally want the first Argo sync to fail with `ImagePullBackOff`. Replace all four tokens with a real tag before applying:

- `metadata.name`
- `static.wesen.dev/release`
- container `image`
- shell `release` variable

#### `ingress.yaml`

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: rag-evaluation-storybook
  annotations:
    argocd.argoproj.io/sync-wave: "2"
    cert-manager.io/cluster-issuer: letsencrypt-prod
  labels:
    app.kubernetes.io/name: rag-evaluation-storybook
    app.kubernetes.io/component: static-site
spec:
  ingressClassName: traefik
  tls:
    - hosts:
        - rag-evaluation-storybook.yolo.scapegoat.dev
      secretName: rag-evaluation-storybook-tls
  rules:
    - host: rag-evaluation-storybook.yolo.scapegoat.dev
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: static-sites-host
                port:
                  name: http
```

### Phase 7: Add the Argo CD Application

Create:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/applications/rag-evaluation-storybook.yaml
```

Suggested content:

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: rag-evaluation-storybook
  namespace: argocd
  annotations:
    argocd.argoproj.io/sync-wave: "1"
    scapegoat.dev/description: "RAG Evaluation Site package Storybook static site"
  labels:
    app.kubernetes.io/name: "rag-evaluation-storybook"
    app.kubernetes.io/part-of: "static-sites"
    app.kubernetes.io/managed-by: "argocd"
    scapegoat.dev/tier: "static-site"
    scapegoat.dev/source-type: "kustomize"
    scapegoat.dev/has-database: "false"
    scapegoat.dev/has-persistent-storage: "false"
    scapegoat.dev/has-ingress: "true"
    scapegoat.dev/database-type: "none"
spec:
  project: static-sites
  source:
    repoURL: https://github.com/wesen/2026-03-27--hetzner-k3s.git
    targetRevision: HEAD
    path: gitops/kustomize/rag-evaluation-storybook
  destination:
    server: https://kubernetes.default.svc
    namespace: static-sites
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
```

### Phase 8: Add Vault policies and roles in the K3s repo

#### Kubernetes image-pull policy

Create:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/vault/policies/kubernetes/rag-evaluation-storybook.hcl
```

```hcl
path "kv/data/apps/rag-evaluation-storybook/prod/image-pull" {
  capabilities = ["read"]
}

path "kv/metadata/apps/rag-evaluation-storybook/prod/image-pull" {
  capabilities = ["read"]
}
```

#### Kubernetes image-pull role

Create:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/vault/roles/kubernetes/rag-evaluation-storybook.json
```

```json
{
  "bound_service_account_names": ["rag-evaluation-storybook"],
  "bound_service_account_namespaces": ["static-sites"],
  "policies": ["rag-evaluation-storybook"],
  "token_ttl": "1h"
}
```

#### GitHub Actions GitOps PR policy

Create:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/vault/policies/github-actions/rag-evaluation-system-gitops-pr.hcl
```

```hcl
# GitHub Actions role for the rag-evaluation-system Storybook release workflow.
# This policy reads only the GitOps PR credential. It does not grant app runtime
# secret access or broad KV traversal.
path "kv/data/ci/github/rag-evaluation-system/gitops-pr-token" {
  capabilities = ["read"]
}

path "auth/token/lookup-self" {
  capabilities = ["read"]
}

path "auth/token/renew-self" {
  capabilities = ["update"]
}

path "auth/token/revoke-self" {
  capabilities = ["update"]
}
```

#### GitHub Actions GitOps PR role

Create:

```text
/home/manuel/code/wesen/2026-03-27--hetzner-k3s/vault/roles/github-actions/rag-evaluation-system-gitops-pr.json
```

```json
{
  "role_type": "jwt",
  "user_claim": "repository",
  "bound_audiences": ["https://vault.yolo.scapegoat.dev"],
  "bound_claims": {
    "repository_owner": "go-go-golems",
    "repository": "go-go-golems/rag-evaluation-system",
    "ref": "refs/heads/main",
    "event_name": "push"
  },
  "policies": ["gha-rag-evaluation-system-gitops-pr"],
  "ttl": "10m",
  "max_ttl": "30m",
  "token_explicit_max_ttl": "30m"
}
```

Check the exact GitHub repository name before merging. The package metadata currently points at `go-go-golems/rag-evaluation-system`, but local worktrees can have nested or alternate origins. The Vault role must match the repository claims emitted by GitHub Actions.

### Phase 9: Bootstrap Vault roles and seed secrets

Run these from the K3s repo:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
export VAULT_TOKEN="$(vault print token 2>/dev/null || cat ~/.vault-token)"

./scripts/bootstrap-vault-kubernetes-auth.sh
./scripts/bootstrap-vault-github-actions-oidc.sh
```

Seed the image-pull credential used by the cluster publisher Job:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
GHCR_TOKEN="$(gh auth token)"
GHCR_AUTH="$(printf 'go-go-golems:%s' "$GHCR_TOKEN" | base64 -w0)"

vault kv put kv/apps/rag-evaluation-storybook/prod/image-pull \
  server=ghcr.io \
  username=go-go-golems \
  password="$GHCR_TOKEN" \
  auth="$GHCR_AUTH"
```

Seed the GitOps PR token used by the source workflow:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
GITOPS_PR_TOKEN="$(gh auth token)"

vault kv put kv/ci/github/rag-evaluation-system/gitops-pr-token \
  token="$GITOPS_PR_TOKEN"
```

Security notes:

- Prefer narrow tokens over a broad operator token.
- The GitOps PR token only needs to push branches and open PRs against `wesen/2026-03-27--hetzner-k3s`.
- The image-pull token only needs package read access if the GHCR package is private.
- If the GHCR package is public, the image-pull secret is still harmless and keeps the pattern consistent, but it may not be strictly required.

## API and contract references

### Source workflow API: `publish-ghcr-image.yml`

Important inputs for this deployment:

| Input | Value | Reason |
|---|---|---|
| `dockerfile` | `./Dockerfile.storybook-static` | Artifact image copies Storybook static output to `/site`. |
| `build_context` | `.` | Dockerfile needs repo-relative `packages/rag-evaluation-site/storybook-static`. |
| `test_command` | `corepack`, `pnpm install`, `typecheck`, `build-storybook`, `test -f index.html` | Proves the exact static output exists before Docker build. |
| `image_name` | `ghcr.io/go-go-golems/rag-evaluation-storybook` | Stable image package for this static site. |
| `gitops_target_config` | `deploy/gitops-targets.json` | Source-owned target metadata. |
| `gitops_pr_token_source` | `vault` | Current preferred credential model. |
| `vault_role` | `rag-evaluation-system-gitops-pr` | Vault JWT role bound to source repo push events. |
| `vault_secret_path` | `kv/data/ci/github/rag-evaluation-system/gitops-pr-token` | KV v2 path containing the GitOps PR token. |

### GitOps target API: `deploy/gitops-targets.json`

Required fields:

```json
{
  "name": "rag-evaluation-storybook-prod",
  "gitops_repo": "wesen/2026-03-27--hetzner-k3s",
  "gitops_branch": "main",
  "manifest_path": "gitops/kustomize/rag-evaluation-storybook/publish-job.yaml",
  "container_name": "publish",
  "patch_strategy": "static-publisher-job"
}
```

Semantics:

- `manifest_path` must point to the publisher Job in the GitOps repo.
- `container_name` must match the `containers[].name` field in the Job (`publish`).
- `patch_strategy: static-publisher-job` rewrites `sha-*` tokens across the Job.

### Static artifact image contract

The publisher Job assumes:

```text
/site/index.html exists
/site/** contains all Storybook static files
```

The image does not need an entrypoint. It only needs a readable filesystem. Alpine is used because it provides shell tools for local inspection and simple `RUN` validation.

### Publisher Job contract

The publisher Job assumes:

```text
image /site/. -> copy to /srv/sites/<host>/releases/<release>/
/srv/sites/<host>/current -> symlink to releases/<release>
```

Pseudo-code:

```sh
set -eu
host="rag-evaluation-storybook.yolo.scapegoat.dev"
release="sha-abc1234"
base="/srv/sites/${host}"
target="${base}/releases/${release}"
tmp="${target}.tmp"

rm -rf "$tmp" "$target"
mkdir -p "$tmp"
cp -a /site/. "$tmp/"
mv "$tmp" "$target"
ln -sfn "releases/${release}" "${base}/current"
```

Why use a temp directory?

- It avoids exposing a partially copied release directory.
- The final `mv` and symlink update are simpler to reason about.
- Rollback is possible by repointing `current` to an older release directory.

### Host routing contract

The static-sites-host Caddy config uses `{host}`. Therefore the Ingress host and publisher Job host variable must be identical:

```text
rag-evaluation-storybook.yolo.scapegoat.dev
```

If the Ingress host is `rag-eval-storybook.yolo.scapegoat.dev` but the publisher writes to `/srv/sites/rag-evaluation-storybook.yolo.scapegoat.dev`, Caddy will look in the wrong directory and return 404.

## Decision records

### Decision: Deploy Storybook as a static site, not a long-running Storybook server

- **Context:** Storybook production output is static. The cluster already has a shared static-sites-host serving files from a PVC.
- **Options considered:** Run `storybook dev` in a Deployment; run a custom Nginx/Caddy container per Storybook; publish static output through the shared static-sites-host.
- **Decision:** Publish Storybook static output through the shared static-sites-host.
- **Rationale:** Reuses existing GitOps, TLS, PVC, Caddy, and publisher Job patterns. Avoids a Node process in production. Matches the K3s static-site playbook.
- **Consequences:** Releases are immutable and cheap. The first rollout needs a publisher Job and image-pull secret. Storybook must be built before Docker packaging.
- **Status:** proposed

### Decision: Build `packages/rag-evaluation-site` Storybook only

- **Context:** The user explicitly requested the Storybook in `packages/rag-evaluation-site/`, not a parent Storybook.
- **Options considered:** Build from repo root; build `web` Storybook; build package-local Storybook with `pnpm --dir packages/rag-evaluation-site build-storybook`.
- **Decision:** Use the package-local build command.
- **Rationale:** `.storybook/main.ts` in the package selects `../src/**/*.stories.@(ts|tsx)`, and the package README defines this package as the reusable UI layer.
- **Consequences:** The deployed site focuses on package components and Widget IR stories. App-level containers remain out of scope.
- **Status:** proposed

### Decision: Use `static-publisher-job` patch strategy

- **Context:** Static publisher Jobs encode the release token in several places. Kubernetes Jobs are immutable after creation.
- **Options considered:** Patch only the image field; manually update all fields; use infra-tooling's `static-publisher-job` strategy.
- **Decision:** Use `patch_strategy: static-publisher-job` in `deploy/gitops-targets.json`.
- **Rationale:** The infra-tooling action explicitly supports this pattern and rewrites all `sha-*` tokens together.
- **Consequences:** The source workflow must publish `sha-*` tags. The publisher Job must keep a standard release-token shape.
- **Status:** proposed

### Decision: Use Vault-backed GitOps PR credentials

- **Context:** The reusable workflow supports Vault, GitHub secret, and GitHub App modes. Current docs prefer Vault over long-lived source-repo secrets.
- **Options considered:** Store `GITOPS_PR_TOKEN` directly in the source repo; use Vault OIDC; use GitHub App credentials from Vault.
- **Decision:** Use Vault OIDC with `gitops_pr_token_source: vault` for the initial implementation.
- **Rationale:** It matches the existing static-site playbook and keeps the source repo from storing long-lived GitOps credentials directly.
- **Consequences:** K3s repo must carry a GitHub Actions Vault policy/role and the operator must seed the KV path.
- **Status:** proposed

### Decision: Keep image-pull secret wiring even if the GHCR image might be public

- **Context:** GHCR packages may be private by default, and static publisher Jobs need to pull the artifact image.
- **Options considered:** Assume public image and omit pull secret; configure VSO-backed image-pull secret consistently.
- **Decision:** Add VSO-backed image-pull secret wiring.
- **Rationale:** It prevents `ImagePullBackOff` if the package remains private and matches the `dmeta-examples` pattern.
- **Consequences:** More files are needed in the GitOps repo, and Vault must be seeded. The pattern is explicit and operationally familiar.
- **Status:** proposed

## End-to-end implementation checklist

### Source repo checklist

- [ ] Confirm the repo origin is `go-go-golems/rag-evaluation-system` or update Vault role claims accordingly.
- [ ] Run `pnpm --dir packages/rag-evaluation-site typecheck`.
- [ ] Run `pnpm --dir packages/rag-evaluation-site build-storybook`.
- [ ] Add `Dockerfile.storybook-static`.
- [ ] Ensure `.dockerignore` does not exclude `packages/rag-evaluation-site/storybook-static`.
- [ ] Add `deploy/gitops-targets.json` with `patch_strategy: static-publisher-job`.
- [ ] Add `.github/workflows/publish-rag-evaluation-storybook.yml`.
- [ ] Open a source PR and verify the workflow builds but does not push image or open GitOps PR on pull requests.
- [ ] Merge to `main` and verify the workflow publishes `ghcr.io/go-go-golems/rag-evaluation-storybook:sha-*`.

### K3s GitOps checklist

- [ ] Add `gitops/kustomize/rag-evaluation-storybook/kustomization.yaml`.
- [ ] Add service account, Vault connection, Vault auth, VSO image-pull secret, publisher Job, and Ingress.
- [ ] Add `gitops/applications/rag-evaluation-storybook.yaml`.
- [ ] Add Vault Kubernetes policy/role.
- [ ] Add Vault GitHub Actions policy/role.
- [ ] Run bootstrap scripts for Kubernetes auth and GitHub Actions OIDC.
- [ ] Seed the GHCR image-pull secret.
- [ ] Seed the GitOps PR token.
- [ ] Run `kubectl kustomize gitops/kustomize/rag-evaluation-storybook`.
- [ ] Merge the first GitOps PR with a real `sha-*` image.
- [ ] Apply the Argo CD Application once.

## Validation plan

### Local source validation

```bash
cd /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build-storybook
test -f packages/rag-evaluation-site/storybook-static/index.html
docker build -f Dockerfile.storybook-static -t rag-evaluation-storybook:test .
docker run --rm rag-evaluation-storybook:test sh -c 'test -f /site/index.html && find /site -maxdepth 2 -type f | sort | head -50'
```

### GitOps render validation

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
kubectl kustomize gitops/kustomize/rag-evaluation-storybook >/tmp/rag-evaluation-storybook-rendered.yaml
rg -n 'rag-evaluation-storybook|sha-|VaultStaticSecret|Ingress|static-sites-host|static-sites-content' /tmp/rag-evaluation-storybook-rendered.yaml
```

Verify these invariants:

- namespace is `static-sites`;
- Application project is `static-sites`;
- publisher Job name, release label, image tag, and shell `release` use the same `sha-*` token;
- publisher container name is `publish`;
- publisher image is `ghcr.io/go-go-golems/rag-evaluation-storybook:sha-*`;
- `serviceAccountName` is `rag-evaluation-storybook`;
- image pull secret is `rag-evaluation-storybook-ghcr-pull`;
- Ingress host is exactly `rag-evaluation-storybook.yolo.scapegoat.dev`;
- Ingress backend is `static-sites-host` port `http`;
- publisher writes to the same hostname used by Ingress.

### First Argo bootstrap

Only required once after the Application manifest is merged:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
export KUBECONFIG=$PWD/.cache/kubeconfig-tailnet.yaml
kubectl apply -f gitops/applications/rag-evaluation-storybook.yaml
kubectl -n argocd annotate application rag-evaluation-storybook argocd.argoproj.io/refresh=hard --overwrite
```

### Cluster rollout validation

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
export KUBECONFIG=$PWD/.cache/kubeconfig-tailnet.yaml

kubectl -n argocd get application rag-evaluation-storybook
kubectl -n static-sites get vaultauth rag-evaluation-storybook
kubectl -n static-sites get vaultstaticsecret rag-evaluation-storybook-ghcr-pull
kubectl -n static-sites get secret rag-evaluation-storybook-ghcr-pull rag-evaluation-storybook-tls
kubectl -n static-sites get job,pod | rg 'rag-evaluation-storybook|NAME'
kubectl -n static-sites logs job/publish-rag-evaluation-storybook-sha-<shortsha>
```

Expected state:

```text
Application: Synced Healthy
VaultAuth: healthy/ready
VaultStaticSecret: synced/healthy/ready
Job: Complete 1/1
Ingress: has cluster ingress address
```

### HTTP smoke test

```bash
curl -fsSI https://rag-evaluation-storybook.yolo.scapegoat.dev/
curl -fsSL https://rag-evaluation-storybook.yolo.scapegoat.dev/ | head
```

Optional asset check:

```bash
curl -fsSL https://rag-evaluation-storybook.yolo.scapegoat.dev/ | rg -o 'src="[^"]+|href="[^"]+' | head
```

If root works but assets fail, inspect whether Storybook emitted absolute asset paths that do not match the hosted root. Hosting at `/` on its own hostname should normally be fine.

## Common failure modes and fixes

| Symptom | Likely cause | Fix |
|---|---|---|
| Workflow builds the wrong Storybook | Command ran from repo root or `web/` | Use `pnpm --dir packages/rag-evaluation-site build-storybook`. |
| Docker build fails with `COPY ... storybook-static: no such file` | `test_command` did not run Storybook build, or `.dockerignore` excluded the output | Build Storybook before Docker build and allow `packages/rag-evaluation-site/storybook-static/**` in context. |
| Workflow publishes image but no GitOps PR appears | `open_gitops_pr` false, run not on `main`, or Vault role/secret path mismatch | Check workflow condition and Vault role bound claims. |
| Vault action fails with permission denied | GitHub repository claims do not match Vault role | Verify `repository_owner`, `repository`, `ref`, and `event_name` in the role JSON. |
| GitOps PR patches only `image:` | Missing `patch_strategy: static-publisher-job` | Add the patch strategy to `deploy/gitops-targets.json`. |
| Argo reports Job immutable field error | Job name/release label/shell release did not change with image tag | Ensure all `sha-*` tokens changed together. |
| Publisher pod has `ImagePullBackOff` | GHCR package private or pull secret wrong | Check VSO secret status and Vault fields under `kv/apps/rag-evaluation-storybook/prod/image-pull`. |
| Site returns 404 | Publisher did not write `/srv/sites/<host>/current`, host mismatch, or Job failed | Check Job logs and compare Ingress host with shell `host` variable. |
| TLS secret missing | cert-manager still solving or Ingress host/DNS wrong | Check `kubectl -n static-sites get certificate,challenge,order`. |
| Storybook loads but assets 404 | Wrong static asset base path | Verify Storybook build output and asset URLs; hosting on root hostname should avoid subpath issues. |

## Intern mental model

Think of the deployment as three contracts, not one script.

### Contract 1: Source build contract

The source repo promises:

```text
Given a commit SHA,
produce a Storybook static directory at packages/rag-evaluation-site/storybook-static
and package it into an image whose /site contains index.html.
```

### Contract 2: GitOps handoff contract

The source workflow promises:

```text
Given image ghcr.io/go-go-golems/rag-evaluation-storybook:sha-abc1234,
open a PR that updates gitops/kustomize/rag-evaluation-storybook/publish-job.yaml
so the cluster desired state points at that immutable image.
```

### Contract 3: Cluster serving contract

The cluster promises:

```text
Given a publisher Job with image /site and host rag-evaluation-storybook.yolo.scapegoat.dev,
copy /site into /srv/sites/<host>/releases/<sha>,
point /srv/sites/<host>/current at that release,
and serve it through static-sites-host and Traefik.
```

When debugging, identify which contract failed before changing code.

## Suggested implementation order for the intern

1. **Read the package docs.** Start with `packages/rag-evaluation-site/README.md` and `GUIDELINES.md` so you understand why this Storybook exists.
2. **Build locally.** Run package typecheck and `build-storybook`. Verify `storybook-static/index.html`.
3. **Package locally.** Add Dockerfile, build the artifact image, and inspect `/site`.
4. **Add source workflow metadata.** Add `deploy/gitops-targets.json` and the caller workflow.
5. **Add K3s manifests.** Copy the `dmeta-examples` shape and rename consistently.
6. **Add Vault policy/role files.** Keep GitHub Actions PR credentials separate from cluster image-pull credentials.
7. **Validate rendered YAML.** Use `kubectl kustomize` before opening any PR.
8. **Run the first source workflow.** Confirm the GHCR `sha-*` image exists.
9. **Merge GitOps PR.** Apply the Argo Application once if this is the first rollout.
10. **Smoke test HTTPS.** Confirm Storybook loads at the chosen hostname.

## Implementation results from this session

The initial implementation now exists across the source repository and K3s GitOps repository. The deployment has not been rolled out to the live cluster yet because the publisher Job still intentionally uses `sha-0000000` until the first CI-published GHCR image exists.

### Source repository commits

- `2f3e0836ee86034adbbbaf1dc1364f9d966b759f` — `docs: plan Storybook deployment ticket`
  - Created the docmgr ticket deliverables and phased task plan.
- `8a32d252ad6f779dd47f1b59fee8fa57c5972cc4` — `build: package rag evaluation Storybook static artifact`
  - Added `Dockerfile.storybook-static`.
  - Added `.dockerignore`.
  - Updated `.gitignore` to ignore generated `packages/rag-evaluation-site/storybook-static/`.
  - Validated package typecheck, Storybook build, Docker build, and `/site/index.html` in the artifact image.
- `4547a48f0e5f71f86ad8ebb70901a512e0976538` — `ci: publish Storybook static artifact`
  - Added `.github/workflows/publish-rag-evaluation-storybook.yml`.
  - Added `deploy/gitops-targets.json` using `patch_strategy: static-publisher-job`.
  - Validated target metadata with infra-tooling's `validate_gitops_targets.py`.
- `947b550fe9d1b0fbfb7fd579fcfa97216b07eea9` — `docs: record Storybook K3s manifest phase`
  - Recorded the K3s static-site manifest phase in the diary/tasks.
- `4fba6e8fe25a50824297b749db4086fc0acef0b2` — `docs: record Storybook Vault wiring phase`
  - Recorded the K3s Vault wiring phase in the diary/tasks.

### K3s GitOps repository commits

- `7094ddccee1ab5d0c9a5838baf086f0a072780ca` — `gitops: add rag evaluation Storybook static site`
  - Added `gitops/kustomize/rag-evaluation-storybook/`.
  - Added `gitops/applications/rag-evaluation-storybook.yaml`.
  - Validated `kubectl kustomize gitops/kustomize/rag-evaluation-storybook`.
- `bbdcab44d59a5f92f6a5a6308f0811146d4089c0` — `vault: add rag evaluation Storybook deployment roles`
  - Added Kubernetes Vault policy/role for the image-pull secret.
  - Added GitHub Actions Vault policy/role for the source workflow GitOps PR token.
  - Validated both role JSON files with `python3 -m json.tool`.

### Validation commands that passed

```bash
cd /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build-storybook
docker build -f Dockerfile.storybook-static -t rag-evaluation-storybook:test .
docker run --rm rag-evaluation-storybook:test sh -c 'test -f /site/index.html && find /site -maxdepth 2 -type f | sort | head -30'
python3 /home/manuel/code/wesen/go-go-golems/infra-tooling/scripts/gitops/validate_gitops_targets.py deploy/gitops-targets.json
```

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
kubectl kustomize gitops/kustomize/rag-evaluation-storybook >/tmp/rag-evaluation-storybook-rendered.yaml
python3 -m json.tool vault/roles/kubernetes/rag-evaluation-storybook.json
python3 -m json.tool vault/roles/github-actions/rag-evaluation-system-gitops-pr.json
```

### Remaining rollout steps

1. Push/merge the source repository changes so the GitHub Actions workflow can publish the first real `ghcr.io/go-go-golems/rag-evaluation-storybook:sha-*` image.
2. Let infra-tooling open a GitOps PR that replaces `sha-0000000` in `publish-job.yaml`.
3. Bootstrap/seed live Vault paths:
   - `kv/apps/rag-evaluation-storybook/prod/image-pull`
   - `kv/ci/github/rag-evaluation-system/gitops-pr-token`
4. Merge the GitOps PR with a real image tag.
5. Apply the Argo CD Application once:

```bash
cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s
export KUBECONFIG=$PWD/.cache/kubeconfig-tailnet.yaml
kubectl apply -f gitops/applications/rag-evaluation-storybook.yaml
kubectl -n argocd annotate application rag-evaluation-storybook argocd.argoproj.io/refresh=hard --overwrite
```

6. Validate `Application` health, `VaultStaticSecret`, publisher Job completion, TLS, and HTTPS smoke tests.


### Live rollout result

The first live rollout completed after the initial implementation. The final public URL is:

```text
https://rag-evaluation-storybook.yolo.scapegoat.dev/
```

Rollout evidence:

- Source PR #4 merged: `https://github.com/go-go-golems/rag-evaluation-system/pull/4`.
- K3s scaffold PR #108 merged: `https://github.com/wesen/2026-03-27--hetzner-k3s/pull/108`.
- First release image pushed manually because the source workflow startup failed: `ghcr.io/go-go-golems/rag-evaluation-storybook:sha-09d7628`.
- K3s release PR #109 merged: `https://github.com/wesen/2026-03-27--hetzner-k3s/pull/109`.
- K3s VaultConnection ownership fix PR #110 merged: `https://github.com/wesen/2026-03-27--hetzner-k3s/pull/110`.
- Argo CD final state: `Synced Healthy 982d63aa2f69d118b7aa7bdd6fc0673beedadc21`.
- Publisher Job final state: `publish-rag-evaluation-storybook-sha-09d7628` completed.
- HTTPS smoke test returned `HTTP/2 200`.

Important follow-up: `.github/workflows/publish-rag-evaluation-storybook.yml` failed with `startup_failure` on the first `main` push before jobs were created. The first deployment was completed manually using the same artifact and GitOps contracts. Future automated releases require debugging that workflow startup failure.

## Open questions

1. **Final hostname.** This guide uses `rag-evaluation-storybook.yolo.scapegoat.dev`. Confirm this is the desired public hostname before live rollout.
2. **GHCR visibility.** If `ghcr.io/go-go-golems/rag-evaluation-storybook` is public, the image-pull secret is optional. If private, the VSO pull secret is required.
3. **GitHub repository identity.** The package metadata says `go-go-golems/rag-evaluation-system`; confirm the actual source repository claims before finalizing the Vault GitHub Actions role.
4. **Workflow environment.** The reusable infra-tooling workflow sets up Go by default. This source repo has `go.mod`, so that should work, but the Storybook build itself depends on Node/Corepack/Pnpm inside `test_command`.
5. **Package manager lockfile mode.** Confirm whether CI should use the repository root `pnpm-lock.yaml` or package-local lockfile if package-local installation conventions changed.

## References

### Source package files

- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/package.json`
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/.storybook/main.ts`
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/vite.config.ts`
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/vite.app.config.ts`
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/index.ts`
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx`
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/README.md`
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/GUIDELINES.md`

### K3s and static-site files

- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/docs/static-site-packaging-and-gitops-playbook.md`
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/static-sites-host/caddy-config.yaml`
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/dmeta-examples/publish-job.yaml`
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/dmeta-examples/kustomization.yaml`
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/dmeta-examples/ingress.yaml`
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/applications/dmeta-examples.yaml`
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/vault/policies/github-actions/dmeta-gitops-pr.hcl`
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/vault/roles/github-actions/dmeta-gitops-pr.json`
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/vault/policies/kubernetes/dmeta-examples.hcl`
- `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/vault/roles/kubernetes/dmeta-examples.json`

### Infra-tooling files

- `/home/manuel/code/wesen/go-go-golems/infra-tooling/.github/workflows/publish-ghcr-image.yml`
- `/home/manuel/code/wesen/go-go-golems/infra-tooling/actions/open-gitops-pr/action.yml`
- `/home/manuel/code/wesen/go-go-golems/infra-tooling/docs/platform/source-repo-to-gitops-pr.md`
