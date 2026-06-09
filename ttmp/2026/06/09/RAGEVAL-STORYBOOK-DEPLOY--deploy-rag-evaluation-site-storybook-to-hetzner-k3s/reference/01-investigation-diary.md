---
Title: Investigation diary
Ticket: RAGEVAL-STORYBOOK-DEPLOY
Status: active
Topics:
    - storybook
    - deployment
    - kubernetes
    - static-sites
    - rag-evaluation
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: 2026-05-27--rag-evaluation-system/ttmp/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY--deploy-rag-evaluation-site-storybook-to-hetzner-k3s/design-doc/01-intern-guide-deploy-rag-evaluation-site-storybook-to-hetzner-k3s.md
      Note: Primary design guide produced in this investigation
ExternalSources: []
Summary: Chronological diary for the rag-evaluation-site Storybook deployment design work.
LastUpdated: 2026-06-09T20:59:15-04:00
WhatFor: Use this to understand how the deployment guide was researched, what evidence was inspected, and what remains to implement.
WhenToUse: Before continuing implementation of the Storybook deployment ticket.
---


# Diary

## Goal

This diary records the evidence-gathering and documentation work for deploying the `packages/rag-evaluation-site` Storybook to the Hetzner K3s static-sites platform. It is meant to make the next implementation step easy to resume without re-discovering which Storybook, which deployment pattern, or which repos are involved.

## Step 1: Create ticket and write deployment design guide

I created a new docmgr ticket under the requested `2026-05-27--rag-evaluation-system/ttmp` root and wrote a detailed intern-oriented deployment guide. The guide focuses on the package-local Storybook in `packages/rag-evaluation-site/`, not any parent or app Storybook, and maps it onto the existing K3s static-site release model.

The main outcome is a concrete implementation plan covering source repo workflow files, static artifact Docker packaging, `deploy/gitops-targets.json`, K3s GitOps manifests, Vault roles/policies, bootstrap commands, validation steps, and common failure modes.

### Prompt Context

**User prompt (verbatim):** "Create a new docmgr in 2026-05-27--rag-evaluation-system/ttmp with `docmgr --root 2026-05-27--rag-evaluation-system/ttmp ...` to deploy the storybook to ~/code/wesen/2026-03-27--hetzner-k3s/ (using potentially ~/code/wesen/go-go-golems/infra-tooling/) . Story book is in @2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/ (not the parent one). 

Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Create a new docmgr ticket in the RAG Evaluation System `ttmp` tree, investigate the package-local Storybook and the K3s static-site deployment pattern, write a detailed design/implementation guide for an intern, and upload the result to reMarkable.

**Inferred user intent:** Produce a complete, technically precise handoff document so a junior engineer can safely implement the Storybook deployment through existing infra rather than guessing from scattered repos.

**Commit (code):** N/A — documentation-only work in docmgr ticket workspace.

### What I did

- Created ticket `RAGEVAL-STORYBOOK-DEPLOY` with:
  - `docmgr --root 2026-05-27--rag-evaluation-system/ttmp ticket create-ticket --ticket RAGEVAL-STORYBOOK-DEPLOY --title "Deploy rag-evaluation-site Storybook to Hetzner K3s" --topics storybook,deployment,kubernetes,static-sites,rag-evaluation`
- Added the primary design document:
  - `docmgr --root 2026-05-27--rag-evaluation-system/ttmp doc add --ticket RAGEVAL-STORYBOOK-DEPLOY --doc-type design-doc --title "Intern Guide: Deploy rag-evaluation-site Storybook to Hetzner K3s"`
- Added this diary document:
  - `docmgr --root 2026-05-27--rag-evaluation-system/ttmp doc add --ticket RAGEVAL-STORYBOOK-DEPLOY --doc-type reference --title "Investigation diary"`
- Inspected the package-local Storybook files:
  - `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/package.json`
  - `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/.storybook/main.ts`
  - `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/vite.config.ts`
  - `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/vite.app.config.ts`
  - `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/index.ts`
  - `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx`
  - `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/README.md`
  - `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/GUIDELINES.md`
- Inspected the K3s static-site deployment reference implementation:
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
- Inspected infra-tooling handoff machinery:
  - `/home/manuel/code/wesen/go-go-golems/infra-tooling/.github/workflows/publish-ghcr-image.yml`
  - `/home/manuel/code/wesen/go-go-golems/infra-tooling/actions/open-gitops-pr/action.yml`
  - `/home/manuel/code/wesen/go-go-golems/infra-tooling/docs/platform/source-repo-to-gitops-pr.md`
- Wrote the design guide at:
  - `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY--deploy-rag-evaluation-site-storybook-to-hetzner-k3s/design-doc/01-intern-guide-deploy-rag-evaluation-site-storybook-to-hetzner-k3s.md`

### Why

- The user asked for a guide that a new intern can follow, so the document needed both high-level system explanation and copy/paste-ready implementation details.
- The user explicitly warned that the Storybook is in `packages/rag-evaluation-site/`, so the guide repeatedly distinguishes package Storybook from parent/app build outputs.
- The existing K3s static-sites model already solves static hosting, TLS, GitOps handoff, and immutable release promotion; reusing it is safer than inventing a new Storybook-specific server.

### What worked

- `docmgr --root 2026-05-27--rag-evaluation-system/ttmp status --summary-only` confirmed the requested docs root exists and contains existing RAG Evaluation System tickets.
- The package files clearly identify `build-storybook`, `.storybook/main.ts`, and the package ownership boundary.
- The K3s repo contains a complete static-site reference (`dmeta-examples`) and a detailed static-site playbook.
- The infra-tooling reusable workflow and action already support the needed `static-publisher-job` patch strategy.

### What didn't work

- Looking for `../2026-05-27--rag-evaluation-system` from the current working directory failed:
  - command context: `ls -la ../2026-05-27--rag-evaluation-system`
  - error: `ls: cannot access '../2026-05-27--rag-evaluation-system': No such file or directory`
- I found the requested relative tree under the current repository instead:
  - `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system`
- I also found other copies under `/home/manuel/code/wesen/2026-05-27--rag-evaluation-system` and older workspace paths, but used the relative path requested by the user for docmgr operations.

### What I learned

- The source repo has no root `package.json` in this worktree, but it does have `pnpm-workspace.yaml` listing `web` and `packages/rag-evaluation-site`.
- The package has distinct outputs: `dist/` for npm, `app-dist/` for the default Vite app, and `storybook-static/` for Storybook.
- The K3s `static-sites-host` Caddy config uses the request hostname as part of the content path, so Ingress host and publisher Job host must match exactly.
- The infra-tooling static publisher patch strategy exists specifically because Kubernetes Jobs are immutable and the release token appears in several fields.

### What was tricky to build

- The main tricky part was preventing ambiguity between similarly named front-end targets. The source tree contains package code, a `web` workspace, `dist`, `app-dist`, and Storybook. The guide addresses this by making `pnpm --dir packages/rag-evaluation-site build-storybook` the core invariant and by explicitly saying not to deploy `dist/`, `app-dist/`, or a parent Storybook.
- Another tricky part was the cross-repo credential boundary. The source workflow publishes a GHCR image using `GITHUB_TOKEN`, but opening a PR against the K3s GitOps repo requires a separate GitOps credential obtained through Vault OIDC. The guide separates these credentials and documents both Vault policy/role paths.
- The static publisher Job also has a subtle immutability invariant: all `sha-*` tokens must change together. The guide explains the failure mode and prescribes `patch_strategy: static-publisher-job`.

### What warrants a second pair of eyes

- Confirm the final public hostname. The guide uses `rag-evaluation-storybook.yolo.scapegoat.dev` as a concrete proposal.
- Confirm the actual GitHub repository owner/name for the source workflow before finalizing Vault `bound_claims`. The package metadata says `go-go-golems/rag-evaluation-system`.
- Confirm whether the GHCR package will be public or private. The guide includes image-pull secret wiring for the safer private-package case.
- Confirm the CI package-manager install command. The guide proposes `pnpm install --frozen-lockfile` from the source repo root.

### What should be done in the future

- Implement the source repo files described in the guide.
- Implement the K3s GitOps and Vault files described in the guide.
- Run the local Storybook and Docker validations.
- Bootstrap Vault secrets and apply the first Argo CD Application after the first real `sha-*` image exists.

### Code review instructions

- Start with the design guide:
  - `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY--deploy-rag-evaluation-site-storybook-to-hetzner-k3s/design-doc/01-intern-guide-deploy-rag-evaluation-site-storybook-to-hetzner-k3s.md`
- Then inspect the source evidence files listed in its References section.
- Validate docmgr state with:
  - `docmgr --root 2026-05-27--rag-evaluation-system/ttmp doc list --ticket RAGEVAL-STORYBOOK-DEPLOY`
  - `docmgr --root 2026-05-27--rag-evaluation-system/ttmp doctor --ticket RAGEVAL-STORYBOOK-DEPLOY --stale-after 30`
- When implementation begins, validate source changes with:
  - `pnpm --dir packages/rag-evaluation-site typecheck`
  - `pnpm --dir packages/rag-evaluation-site build-storybook`
  - `docker build -f Dockerfile.storybook-static -t rag-evaluation-storybook:test .`
- Validate K3s changes with:
  - `kubectl kustomize gitops/kustomize/rag-evaluation-storybook`

### Technical details

Key design conclusion:

```text
Use packages/rag-evaluation-site Storybook as a static artifact.
Package storybook-static as /site in GHCR.
Use infra-tooling publish-ghcr-image.yml and open-gitops-pr.
Use static-publisher-job patch strategy.
Serve through K3s static-sites-host.
```

Primary generated document:

```text
/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY--deploy-rag-evaluation-site-storybook-to-hetzner-k3s/design-doc/01-intern-guide-deploy-rag-evaluation-site-storybook-to-hetzner-k3s.md
```

## Step 2: Expand implementation tasks into phased execution plan

I converted the ticket's simple completed checklist into a phased implementation plan. The phases now separate source packaging, source CI/GitOps handoff, K3s static-site manifests, K3s Vault wiring, and final documentation/upload work so the deployment can be implemented and reviewed in focused commits.

This planning step intentionally happened before source or cluster files changed. That makes the ticket useful as a live runbook: each phase has a small set of acceptance criteria, and later diary entries can refer back to the phase being executed.

### Prompt Context

**User prompt (verbatim):** "ok ,add detailed tasks to the ticket's tasks, goruped into different phases, and then execute them in order, committing at appropriate intervals, keeping a frequent diary all along. (we're going to reuse the diary to make it a playbook on how to package storybook for deployment)."

**Assistant interpretation:** Expand the ticket into detailed phased tasks, then implement the deployment in order with focused commits and frequent diary entries that can later become a Storybook deployment playbook.

**Inferred user intent:** Turn the initial design guide into an executable implementation trail, preserving the reasoning and command outcomes so future Storybook deployments can reuse the diary as a procedural playbook.

**Commit (code):** pending — planning-only ticket documentation update.

### What I did
- Rewrote `tasks.md` with phases 0 through 5.
- Kept completed initial design/upload tasks in Phase 0.
- Added source repo packaging tasks for Docker artifact creation and local validation.
- Added source repo CI/GitOps handoff tasks for the infra-tooling workflow and target metadata.
- Added K3s manifest tasks for static-site resources and Argo CD Application wiring.
- Added K3s Vault tasks for image-pull and GitOps PR credential access.
- Added final documentation, doctor, reMarkable upload, and changelog tasks.

### Why
- The implementation touches two Git repositories and a docmgr ticket, so phases reduce the chance of mixing unrelated concerns in one commit.
- The future playbook needs the order of operations, not just the final files.
- The phase structure mirrors the deployment contract boundaries: source artifact, GitOps handoff, cluster manifests, credentials, and delivery docs.

### What worked
- The existing task file was small, so replacing it with a structured phase plan was straightforward.
- The design guide already contained enough detail to derive implementation tasks without another investigation pass.

### What didn't work
- N/A.

### What I learned
- Planning commits are useful here because they preserve the implementation intent before generated validation artifacts or cluster manifests appear.

### What was tricky to build
- The main subtlety was avoiding a task list that says only “deploy Storybook.” The actionable phases need to encode where each contract lives: source repo, infra-tooling API boundary, K3s GitOps repo, and Vault. I split the work along those boundaries so each later commit can be reviewed independently.

### What warrants a second pair of eyes
- Confirm the phase order before any live rollout. In particular, the K3s publisher Job can be committed with a placeholder for render validation, but it should not be applied to the cluster until a real GHCR `sha-*` image exists.

### What should be done in the future
- Execute Phase 1 next: add the static artifact Docker contract, ignore generated Storybook output, and validate package-local Storybook build plus local Docker packaging.

### Code review instructions
- Review `tasks.md` first to confirm the phase boundaries.
- Then review subsequent diary entries and commits phase-by-phase.

### Technical details
- Ticket task file: `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY--deploy-rag-evaluation-site-storybook-to-hetzner-k3s/tasks.md`

## Step 3: Implement the source Storybook static artifact contract

I implemented the first deployable contract in the source repository: after the package-local Storybook build runs, Docker can package the generated `storybook-static` directory as `/site`. This is the minimal artifact shape expected by the K3s static publisher Job.

This step deliberately avoided committing generated Storybook files. The generated output is needed as Docker build input during CI and local validation, but the repository should only track the Dockerfile, Docker context rules, and ignore rules.

### Prompt Context

**User prompt (verbatim):** (see Step 2)

**Assistant interpretation:** Execute Phase 1 from the expanded plan: create the source-side static artifact Docker contract, validate it locally, update the ticket, and commit a focused source change.

**Inferred user intent:** Establish a repeatable, reviewable packaging pattern for Storybook that can become the first concrete section of the future playbook.

**Commit (code):** pending — source artifact contract changes.

### What I did
- Added `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/Dockerfile.storybook-static`.
- Added `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/.dockerignore`.
- Updated `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/.gitignore` to ignore `packages/rag-evaluation-site/storybook-static/`.
- Ran package-local validation:
  - `pnpm --dir packages/rag-evaluation-site typecheck`
  - `pnpm --dir packages/rag-evaluation-site build-storybook`
- Ran Docker artifact validation:
  - `docker build -f Dockerfile.storybook-static -t rag-evaluation-storybook:test .`
  - `docker run --rm rag-evaluation-storybook:test sh -c 'test -f /site/index.html && find /site -maxdepth 2 -type f | sort | head -30'`
- Marked Phase 1 tasks complete in `tasks.md`.

### Why
- The K3s publisher Job copies `/site/.` from an artifact image into the shared static-sites PVC, so the source repo must produce an image with `/site/index.html`.
- Storybook's production build is static, so using an Alpine artifact carrier avoids running a Node/Storybook server in the cluster.
- Ignoring `storybook-static/` prevents a large generated build directory from being accidentally committed.

### What worked
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `pnpm --dir packages/rag-evaluation-site build-storybook` completed successfully and wrote output to `packages/rag-evaluation-site/storybook-static`.
- Docker built the static artifact image successfully.
- Running the image confirmed `/site/index.html` exists and listed Storybook asset files under `/site/assets/`.

### What didn't work
- N/A. The validation commands succeeded on the first run.

### What I learned
- Storybook v10.4.2 emits a large `iframe` bundle and many component/story chunks, but the output is still a normal static directory that fits the `/site` contract.
- Docker's build context was about 7 MB after `.dockerignore`, which is reasonable for this static artifact.

### What was tricky to build
- The important nuance is that `storybook-static/` must be ignored by Git but included in the Docker build context. `.gitignore` and `.dockerignore` have opposite responsibilities here: Git should not track generated Storybook output, while Docker must see it because the artifact Dockerfile copies it into `/site`.
- Another nuance is build ordering. The Dockerfile does not run `pnpm build-storybook`; it assumes CI or the operator already ran it. This keeps the image a simple artifact carrier and lets the reusable workflow's `test_command` own typecheck/build validation before Docker packaging.

### What warrants a second pair of eyes
- Review `.dockerignore` to ensure it excludes large irrelevant directories but still includes `packages/rag-evaluation-site/storybook-static/**`.
- Review whether the Dockerfile should copy a custom Storybook output path in the future. For now it uses the default `storybook-static` path generated by `storybook build`.

### What should be done in the future
- Phase 2 should add `deploy/gitops-targets.json` and the GitHub Actions workflow so CI can publish this artifact image and open a GitOps PR.

### Code review instructions
- Start with `Dockerfile.storybook-static`; verify it copies only `packages/rag-evaluation-site/storybook-static/` to `/site` and checks `/site/index.html`.
- Review `.dockerignore` and `.gitignore` together.
- Re-run:
  - `pnpm --dir packages/rag-evaluation-site typecheck`
  - `pnpm --dir packages/rag-evaluation-site build-storybook`
  - `docker build -f Dockerfile.storybook-static -t rag-evaluation-storybook:test .`
  - `docker run --rm rag-evaluation-storybook:test sh -c 'test -f /site/index.html'`

### Technical details
- Dockerfile artifact contract:
  - input: `packages/rag-evaluation-site/storybook-static/`
  - output inside image: `/site/`
  - invariant: `/site/index.html` must exist
- Successful validation image tag: `rag-evaluation-storybook:test`

## Step 4: Add the source CI workflow and GitOps target metadata

I connected the Storybook artifact contract from Phase 1 to the existing infra-tooling release path. The source repo now has target metadata telling `open-gitops-pr` which K3s manifest to patch, and a GitHub Actions workflow that builds the package-local Storybook, packages it as a GHCR image, and opens a GitOps PR only for pushes to `main`.

This step keeps pull requests safe: PR runs validate the Storybook build and Docker packaging path but do not push images or open GitOps PRs. Only a merged `main` commit can publish `ghcr.io/go-go-golems/rag-evaluation-storybook:sha-*` and hand off to the K3s repo.

### Prompt Context

**User prompt (verbatim):** (see Step 2)

**Assistant interpretation:** Execute Phase 2: add source-side infra-tooling workflow and GitOps target metadata, validate the metadata locally, update the diary, and commit the focused source handoff changes.

**Inferred user intent:** Make Storybook deployment reproducible from CI using the same source-to-GitOps release chain as other static sites.

**Commit (code):** pending — source CI/GitOps handoff changes.

### What I did
- Added `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/deploy/gitops-targets.json`.
- Added `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/.github/workflows/publish-rag-evaluation-storybook.yml`.
- Configured the target metadata to patch `gitops/kustomize/rag-evaluation-storybook/publish-job.yaml` in `wesen/2026-03-27--hetzner-k3s`.
- Set `patch_strategy` to `static-publisher-job` so the Job name, release label, image tag, and shell release token change together.
- Configured the workflow to call `go-go-golems/infra-tooling/.github/workflows/publish-ghcr-image.yml@main`.
- Configured the workflow test command to run:
  - `corepack enable`
  - `pnpm install --frozen-lockfile`
  - `pnpm --dir packages/rag-evaluation-site typecheck`
  - `pnpm --dir packages/rag-evaluation-site build-storybook`
  - `test -f packages/rag-evaluation-site/storybook-static/index.html`
- Validated JSON formatting with `python3 -m json.tool deploy/gitops-targets.json`.
- Validated target metadata with:
  - `python3 /home/manuel/code/wesen/go-go-golems/infra-tooling/scripts/gitops/validate_gitops_targets.py deploy/gitops-targets.json`
- Parsed the workflow with Python/PyYAML and ran sanity checks for the reusable workflow, Storybook build command, and Vault role.

### Why
- The K3s static-sites deployment model depends on a source workflow opening GitOps PRs with immutable image tags.
- `deploy/gitops-targets.json` is the source-owned API contract that tells infra-tooling where the deployment manifest lives.
- Pull request runs should be validation-only to avoid publishing unmerged Storybook artifacts.

### What worked
- `python3 -m json.tool deploy/gitops-targets.json` succeeded.
- `validate_gitops_targets.py` reported: `deploy/gitops-targets.json: OK (1 target(s))`.
- Workflow YAML parsed successfully with PyYAML.

### What didn't work
- N/A. The local metadata and workflow sanity validations passed.

### What I learned
- The shared infra-tooling workflow already has the exact switches needed here: `push_image`, `open_gitops_pr`, Vault-backed GitOps token inputs, and the static publisher patch strategy.
- The deployment target metadata can be very small because the static-site-specific release rewriting lives in infra-tooling.

### What was tricky to build
- The workflow has two different conditions that must stay aligned: `push_image` and `open_gitops_pr`. If `push_image` were true on PRs, CI would publish unmerged artifacts. If `open_gitops_pr` were true outside `main`, GitOps PRs could point at unreviewed commits. Both are now restricted to non-PR pushes, with GitOps PRs further restricted to `refs/heads/main`.
- The workflow relies on the Phase 1 Dockerfile but does not build Storybook inside Docker. The `test_command` must therefore create `storybook-static/` before Docker Buildx runs.

### What warrants a second pair of eyes
- Confirm the source repository's default branch is `main`.
- Confirm the Vault role name `rag-evaluation-system-gitops-pr` and secret path `kv/data/ci/github/rag-evaluation-system/gitops-pr-token` match the K3s Vault files that will be added in Phase 4.
- Confirm whether the GHCR image should be named `rag-evaluation-storybook` or include the repo name, for example `rag-evaluation-system-storybook`.

### What should be done in the future
- Phase 3 should add the K3s static-site manifests at the path referenced by `deploy/gitops-targets.json`.

### Code review instructions
- Review `deploy/gitops-targets.json` and verify `manifest_path` matches the K3s publisher Job file.
- Review `.github/workflows/publish-rag-evaluation-storybook.yml` and verify PRs do not publish images or open GitOps PRs.
- Re-run:
  - `python3 -m json.tool deploy/gitops-targets.json`
  - `python3 /home/manuel/code/wesen/go-go-golems/infra-tooling/scripts/gitops/validate_gitops_targets.py deploy/gitops-targets.json`

### Technical details
- GHCR image: `ghcr.io/go-go-golems/rag-evaluation-storybook`
- GitOps target repo: `wesen/2026-03-27--hetzner-k3s`
- GitOps target manifest: `gitops/kustomize/rag-evaluation-storybook/publish-job.yaml`
- Patch strategy: `static-publisher-job`

## Step 5: Add K3s static-site GitOps manifests

I implemented the cluster-side static-site shape in the K3s GitOps repository. The new `rag-evaluation-storybook` Kustomize package follows the same shared-host pattern as `dmeta-examples`: a ServiceAccount with an image pull secret, Vault Secret Operator resources, a one-shot publisher Job, an Ingress pointing at `static-sites-host`, and an Argo CD Application in the `static-sites` project.

This phase still uses `sha-0000000` as an intentional placeholder in the publisher Job. That lets the manifest render and lets infra-tooling know which tokens to rewrite, but it must be replaced by the first real `sha-*` GHCR image before a live sync is expected to succeed.

### Prompt Context

**User prompt (verbatim):** (see Step 2)

**Assistant interpretation:** Execute Phase 3: add K3s static-site manifests matching the source `deploy/gitops-targets.json` handoff path, validate Kustomize rendering, commit K3s changes, and record the result.

**Inferred user intent:** Put the cluster desired-state scaffold in place so the source workflow has a concrete publisher Job manifest to patch.

**Commit (code):** `7094ddccee1ab5d0c9a5838baf086f0a072780ca` — "gitops: add rag evaluation Storybook static site"

### What I did
- Added K3s Kustomize package:
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/rag-evaluation-storybook/kustomization.yaml`
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/rag-evaluation-storybook/serviceaccount.yaml`
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/rag-evaluation-storybook/vault-connection.yaml`
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/rag-evaluation-storybook/vault-auth.yaml`
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/rag-evaluation-storybook/vault-static-secret-image-pull.yaml`
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/rag-evaluation-storybook/publish-job.yaml`
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/rag-evaluation-storybook/ingress.yaml`
- Added Argo CD Application:
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/applications/rag-evaluation-storybook.yaml`
- Rendered the package with:
  - `kubectl kustomize gitops/kustomize/rag-evaluation-storybook >/tmp/rag-evaluation-storybook-rendered.yaml`
- Grepped the rendered output for the expected app name, placeholder release token, `static-sites-host`, `static-sites-content`, `VaultStaticSecret`, and `Ingress`.
- Committed only the new Storybook GitOps files in the K3s repo. Pre-existing unrelated modified `ttmp` files in that repo were left unstaged.

### Why
- The source workflow added in Phase 2 points at `gitops/kustomize/rag-evaluation-storybook/publish-job.yaml`; that file must exist for `open-gitops-pr` to patch it.
- The static-sites-host model requires a publisher Job to copy `/site` into the shared PVC and an Ingress that routes the hostname to the shared Caddy service.
- The Argo CD Application must be in project `static-sites` to match the cluster's static-site project restrictions.

### What worked
- `kubectl kustomize gitops/kustomize/rag-evaluation-storybook` succeeded.
- The rendered output contained:
  - the ServiceAccount and imagePullSecret name;
  - the publisher Job with `sha-0000000` in the Job name, release label, image tag, and shell release variable;
  - the `static-sites-content` PVC mount;
  - the Ingress backend pointing at `static-sites-host` port `http`;
  - VaultAuth and VaultStaticSecret resources.
- The K3s commit succeeded while leaving unrelated modified docs in the worktree untouched.

### What didn't work
- N/A. Kustomize rendering and the focused K3s commit both succeeded.

### What I learned
- The static-site package is small enough to review as one K3s commit: seven Kustomize resources plus one Argo CD Application.
- Keeping the same host string in the Ingress and publisher Job is essential because `static-sites-host` maps `{host}` directly into `/srv/sites/{host}/current`.

### What was tricky to build
- The placeholder release token is necessary but dangerous. It gives `static-publisher-job` a token to rewrite and makes the manifest structurally valid, but if someone applies the Argo Application before a real image patch lands, the Job will try to pull `ghcr.io/go-go-golems/rag-evaluation-storybook:sha-0000000` and fail. The guide and diary call this out explicitly.
- Another subtlety is that the Argo CD Application is not part of the Kustomize package. It lives under `gitops/applications/` and must be applied once to bootstrap reconciliation.

### What warrants a second pair of eyes
- Confirm `rag-evaluation-storybook.yolo.scapegoat.dev` is the intended hostname.
- Confirm the placeholder token should remain until the first source workflow publishes a real image and opens a GitOps PR.
- Confirm `spec.project: static-sites` is correct for this cluster's current Argo project setup.

### What should be done in the future
- Phase 4 should add the Vault policy and role files that back the VSO image-pull secret and GitHub Actions GitOps PR token access.

### Code review instructions
- In the K3s repo, start at `gitops/kustomize/rag-evaluation-storybook/publish-job.yaml` and verify all `sha-0000000` occurrences are intentional and consistent.
- Review `ingress.yaml` and confirm its host matches the publisher Job `host` variable.
- Review `gitops/applications/rag-evaluation-storybook.yaml` and confirm project `static-sites`.
- Re-run:
  - `cd /home/manuel/code/wesen/2026-03-27--hetzner-k3s`
  - `kubectl kustomize gitops/kustomize/rag-evaluation-storybook >/tmp/rag-evaluation-storybook-rendered.yaml`

### Technical details
- K3s commit: `7094ddccee1ab5d0c9a5838baf086f0a072780ca`
- Hostname: `rag-evaluation-storybook.yolo.scapegoat.dev`
- Placeholder release token: `sha-0000000`
- Publisher image placeholder: `ghcr.io/go-go-golems/rag-evaluation-storybook:sha-0000000`

## Step 6: Add K3s Vault policies and roles for deployment credentials

I added the Vault-side declarative files needed by the Storybook static-site deployment. There are two separate credential paths: Kubernetes reads a GHCR image-pull secret through Vault Secret Operator, and GitHub Actions reads a GitOps PR token through Vault OIDC.

This separation is important for the playbook. The cluster's publisher Job needs only package pull credentials, while the source workflow needs only enough GitHub access to open a PR against the GitOps repo. Neither path should grant broad runtime secret access.

### Prompt Context

**User prompt (verbatim):** (see Step 2)

**Assistant interpretation:** Execute Phase 4: add K3s Vault policy/role manifests for both the cluster-side image-pull secret and the source-side GitOps PR workflow credential, validate the JSON files, and commit focused Vault wiring changes.

**Inferred user intent:** Complete the credential manifests needed before an operator bootstraps Vault and attempts the first CI-to-GitOps rollout.

**Commit (code):** `bbdcab44d59a5f92f6a5a6308f0811146d4089c0` — "vault: add rag evaluation Storybook deployment roles"

### What I did
- Added Kubernetes image-pull policy:
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/vault/policies/kubernetes/rag-evaluation-storybook.hcl`
- Added Kubernetes image-pull role:
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/vault/roles/kubernetes/rag-evaluation-storybook.json`
- Added GitHub Actions GitOps PR policy:
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/vault/policies/github-actions/rag-evaluation-system-gitops-pr.hcl`
- Added GitHub Actions GitOps PR role:
  - `/home/manuel/code/wesen/2026-03-27--hetzner-k3s/vault/roles/github-actions/rag-evaluation-system-gitops-pr.json`
- Validated JSON role syntax with:
  - `python3 -m json.tool vault/roles/kubernetes/rag-evaluation-storybook.json`
  - `python3 -m json.tool vault/roles/github-actions/rag-evaluation-system-gitops-pr.json`
- Grepped the files for expected names and checked that no `TODO` or `PLACEHOLDER` markers remained.
- Committed only the new Vault files in the K3s repo.

### Why
- `vault-static-secret-image-pull.yaml` from Phase 3 references `apps/rag-evaluation-storybook/prod/image-pull`; the Kubernetes Vault policy must allow the service account role to read that path.
- The source workflow from Phase 2 references Vault role `rag-evaluation-system-gitops-pr` and secret path `kv/data/ci/github/rag-evaluation-system/gitops-pr-token`; the GitHub Actions policy/role must make that read possible.
- Keeping these as files in the K3s repo makes bootstrap repeatable via the existing Vault bootstrap scripts.

### What worked
- Both Vault role JSON files parsed successfully.
- The grep check confirmed expected names and paths:
  - `kv/data/apps/rag-evaluation-storybook/prod/image-pull`
  - `kv/data/ci/github/rag-evaluation-system/gitops-pr-token`
  - `go-go-golems/rag-evaluation-system`
  - `gha-rag-evaluation-system-gitops-pr`
- The focused K3s commit succeeded while unrelated modified `ttmp` files remained unstaged.

### What didn't work
- N/A. File creation, syntax validation, and commit succeeded.

### What I learned
- The deployment needs two Vault bootstrap operations later: Kubernetes auth role/policy bootstrap and GitHub Actions OIDC role/policy bootstrap.
- The role names now line up with the source workflow values added in Phase 2.

### What was tricky to build
- The trickiest part is distinguishing Vault path syntax from Vault CLI path syntax. The policy uses `kv/data/...` and `kv/metadata/...` because KV v2 policies operate on API paths. Operator seeding commands later should use `vault kv put kv/apps/...`, not `kv/data/apps/...`.
- The GitHub Actions role is also claim-sensitive. It currently binds to repository `go-go-golems/rag-evaluation-system`, branch `refs/heads/main`, and event `push`. If the actual source repo or default branch differs, Vault authentication will fail even though the files are syntactically valid.

### What warrants a second pair of eyes
- Confirm that `go-go-golems/rag-evaluation-system` is the exact GitHub Actions `repository` claim.
- Confirm the policy name convention `gha-rag-evaluation-system-gitops-pr` matches the bootstrap script's policy naming behavior.
- Confirm the GHCR image-pull username convention before seeding the secret.

### What should be done in the future
- Bootstrap the live Vault policies/roles with the existing scripts after the GitOps changes are merged.
- Seed:
  - `kv/apps/rag-evaluation-storybook/prod/image-pull`
  - `kv/ci/github/rag-evaluation-system/gitops-pr-token`

### Code review instructions
- Review the four Vault files added in commit `bbdcab44d59a5f92f6a5a6308f0811146d4089c0`.
- Re-run JSON syntax checks for both role files.
- Compare the workflow's `vault_role` and `vault_secret_path` against the GitHub Actions Vault role/policy.

### Technical details
- K3s commit: `bbdcab44d59a5f92f6a5a6308f0811146d4089c0`
- Kubernetes role name: `rag-evaluation-storybook`
- GitHub Actions role name: `rag-evaluation-system-gitops-pr`
- GitHub Actions policy name: `gha-rag-evaluation-system-gitops-pr`

## Step 7: Consolidate implementation results for playbook reuse

I updated the design guide with the concrete files, commits, and validation commands produced during the implementation. This turns the original proposed guide into a mixed design-and-playbook artifact: future readers can see both the intended architecture and the exact implementation sequence that worked.

This step stops short of live rollout. The source workflow, K3s manifests, and Vault files are committed, but the publisher Job still uses `sha-0000000` until a real CI-published GHCR image exists and infra-tooling rewrites the Job through a GitOps PR.

### Prompt Context

**User prompt (verbatim):** (see Step 2)

**Assistant interpretation:** Finish Phase 5 documentation work by updating the design guide and diary with implementation results, commits, and validation outcomes, then run docmgr validation and re-upload the updated bundle.

**Inferred user intent:** Preserve enough operational detail that the diary can be refactored into a repeatable playbook for packaging Storybook deployments.

**Commit (code):** pending — final ticket documentation update.

### What I did
- Added an `Implementation results from this session` section to the design guide.
- Listed source repository commits:
  - `2f3e0836ee86034adbbbaf1dc1364f9d966b759f`
  - `8a32d252ad6f779dd47f1b59fee8fa57c5972cc4`
  - `4547a48f0e5f71f86ad8ebb70901a512e0976538`
  - `947b550fe9d1b0fbfb7fd579fcfa97216b07eea9`
  - `4fba6e8fe25a50824297b749db4086fc0acef0b2`
- Listed K3s repository commits:
  - `7094ddccee1ab5d0c9a5838baf086f0a072780ca`
  - `bbdcab44d59a5f92f6a5a6308f0811146d4089c0`
- Added passed validation commands for source, Docker, GitOps target metadata, Kustomize render, and Vault role JSON checks.
- Added remaining rollout steps that require live CI/GitOps/Vault/Argo activity.

### Why
- The diary is intended to become a playbook, so it needs exact command evidence and phase-by-phase lessons rather than just a final summary.
- The design guide should tell reviewers what was actually implemented, not only what was proposed.

### What worked
- The implementation results section could be added without restructuring the entire guide.
- The commits are cleanly separated into source-repo commits and K3s-repo commits.

### What didn't work
- N/A.

### What I learned
- This pattern has a useful reusable sequence: build static output locally, package as `/site`, wire CI to publish immutable image, add static-sites GitOps package, add Vault policies/roles, then only after a real image exists bootstrap Argo and live secrets.

### What was tricky to build
- The documentation has to be clear that “implemented” does not mean “live in cluster.” The GitOps manifests exist, but the placeholder image token intentionally prevents pretending that rollout is complete. The remaining rollout steps are therefore called out as operator work.

### What warrants a second pair of eyes
- Review the implementation results section for any commit hash mistakes before using it as a playbook source.
- Review whether the remaining rollout steps should include actual Vault bootstrap commands in a separate operator-only playbook.

### What should be done in the future
- Run `docmgr doctor`.
- Upload the updated bundle to reMarkable.
- Commit the final ticket documentation update.

### Code review instructions
- Read the implementation results section in the design guide.
- Compare the listed commits with `git log --oneline` in both repositories.
- Confirm the remaining rollout steps accurately describe what is not yet live.

### Technical details
- Updated design doc: `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY--deploy-rag-evaluation-site-storybook-to-hetzner-k3s/design-doc/01-intern-guide-deploy-rag-evaluation-site-storybook-to-hetzner-k3s.md`

## Step 8: Validate ticket and upload the updated playbook bundle

I ran docmgr validation after relating the newly implemented source and K3s files, then uploaded the updated design guide and diary bundle to reMarkable. This is the handoff point for the documentation portion of the work: the ticket now contains the design, implementation evidence, validation outcomes, and remaining live rollout steps.

The upload replaced the previous bundle with the updated implementation-aware version. The generated PDF is suitable for reading as both a design document and an implementation playbook draft.

### Prompt Context

**User prompt (verbatim):** (see Step 2)

**Assistant interpretation:** Complete the final documentation phase by validating docmgr hygiene, uploading the updated bundle, marking the remaining documentation tasks complete, and committing the ticket updates.

**Inferred user intent:** Ensure the reusable playbook material is persisted in docmgr and available on reMarkable after the implementation work.

**Commit (code):** pending — final documentation/upload record.

### What I did
- Related newly implemented files to the primary design doc with `docmgr doc relate`.
- Updated the ticket changelog with the implementation summary.
- Ran:
  - `docmgr --root 2026-05-27--rag-evaluation-system/ttmp doctor --ticket RAGEVAL-STORYBOOK-DEPLOY --stale-after 30`
- Uploaded the updated bundle with:
  - `remarquee upload bundle ... --name "RAGEVAL-STORYBOOK-DEPLOY Guide" --remote-dir "/ai/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY" --toc-depth 2 --force --non-interactive`
- Marked the final Phase 5 tasks complete.

### Why
- The ticket should remain self-contained and navigable from docmgr search/doctor output.
- The reMarkable copy should include the implementation results and not just the initial proposal.
- `--force` was used intentionally because this upload replaces the earlier same-name bundle with the updated version.

### What worked
- `docmgr doctor` passed with all checks green.
- reMarkable upload succeeded with:
  - `OK: uploaded RAGEVAL-STORYBOOK-DEPLOY Guide.pdf -> /ai/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY`

### What didn't work
- N/A.

### What I learned
- The docmgr + reMarkable loop is useful at phase boundaries: the guide now captures both proposed and observed behavior, and the diary has enough detail to become a packaging playbook.

### What was tricky to build
- The main risk was accidentally documenting the deployment as fully live. I kept the distinction explicit: packaging, CI handoff, GitOps manifests, and Vault files are implemented; live rollout still needs a real GHCR image, Vault seeding/bootstrap, GitOps PR merge, Argo Application apply, and HTTP smoke test.

### What warrants a second pair of eyes
- Review the remaining rollout section before using it operationally.
- Confirm whether to split the diary into a standalone “Package Storybook for static-sites-host” playbook after the first live rollout validates the remaining steps.

### What should be done in the future
- Run the source workflow on `main` and confirm the first real image/GitOps PR.
- Replace the placeholder `sha-0000000` through infra-tooling, not by hand if possible.
- Bootstrap and seed Vault secrets, apply the Argo Application once, and smoke test the public URL.

### Code review instructions
- Review the final `tasks.md` to confirm all implementation documentation tasks are marked complete.
- Review the latest reMarkable bundle if reading offline.
- Run `docmgr doctor` again after any follow-up edits.

### Technical details
- reMarkable path: `/ai/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY`
- Uploaded file name: `RAGEVAL-STORYBOOK-DEPLOY Guide.pdf`

## Step 9: Roll out the Storybook static site live

I completed the live rollout after the implementation commits were merged. The happy path was mostly successful, but the source GitHub Actions workflow failed at startup, so I used the same artifact and GitOps contracts manually for the first deployment: build/push the GHCR image, patch the K3s publisher Job from `sha-0000000` to `sha-09d7628`, merge that GitOps PR, apply the Argo CD Application, and smoke test HTTPS.

The site is now live at `https://rag-evaluation-storybook.yolo.scapegoat.dev/`. The Argo CD Application is `Synced Healthy`, the publisher Job completed, Vault Secret Operator created the GHCR pull secret, cert-manager issued the TLS secret, and Caddy serves the Storybook HTML with HTTP 200.

### Prompt Context

**User prompt (verbatim):** "go ahead"

**Assistant interpretation:** Proceed from the planned next steps into the live rollout: merge scaffolding, bootstrap Vault, publish the first static image, patch GitOps to a real release token, apply Argo, and validate the public URL.

**Inferred user intent:** Move beyond documentation/scaffold into an actual deployed Storybook site.

**Commit (code):** pending — rollout diary update.

### What I did
- Created clean worktrees to avoid unrelated local changes:
  - `/tmp/rag-storybook-source-deploy`
  - `/tmp/rag-storybook-k3s-deploy`
- Cherry-picked the source implementation commits onto a clean source branch and opened PR #4:
  - `https://github.com/go-go-golems/rag-evaluation-system/pull/4`
- Cherry-picked the K3s scaffold commits onto a clean K3s branch and opened/merged PR #108:
  - `https://github.com/wesen/2026-03-27--hetzner-k3s/pull/108`
- Bootstrapped live Vault policy/role config:
  - `./scripts/bootstrap-vault-kubernetes-auth.sh`
  - `./scripts/bootstrap-vault-github-actions-oidc.sh`
- Seeded live Vault paths without printing token values:
  - `kv/apps/rag-evaluation-storybook/prod/image-pull`
  - `kv/ci/github/rag-evaluation-system/gitops-pr-token`
- Merged source PR #4 after checks passed.
- Observed that `publish-rag-evaluation-storybook` on `main` failed with `startup_failure` before jobs were created.
- Manually built and pushed the first release image:
  - `ghcr.io/go-go-golems/rag-evaluation-storybook:sha-09d7628`
- Created and merged K3s release PR #109:
  - `https://github.com/wesen/2026-03-27--hetzner-k3s/pull/109`
  - changed all publisher Job release tokens from `sha-0000000` to `sha-09d7628`
- Applied the Argo CD Application once:
  - `kubectl apply -f gitops/applications/rag-evaluation-storybook.yaml`
  - `kubectl -n argocd annotate application rag-evaluation-storybook argocd.argoproj.io/refresh=hard --overwrite`
- Observed the app initially reached `Healthy` but remained `OutOfSync` because `VaultConnection/vault` was shared with `dmeta-examples`.
- Created and merged K3s PR #110 to give Storybook a dedicated VaultConnection:
  - `https://github.com/wesen/2026-03-27--hetzner-k3s/pull/110`
  - `VaultConnection` name: `rag-evaluation-storybook-vault`
  - `VaultAuth.spec.vaultConnectionRef`: `rag-evaluation-storybook-vault`
- Refreshed Argo and confirmed final status:
  - `Synced Healthy 982d63aa2f69d118b7aa7bdd6fc0673beedadc21`
- Smoke tested HTTPS:
  - `curl -fsSI https://rag-evaluation-storybook.yolo.scapegoat.dev/`
  - `curl -fsSL https://rag-evaluation-storybook.yolo.scapegoat.dev/ | head`

### Why
- The static-site scaffold needed a real immutable release image before it could be safely applied.
- The source workflow startup failure blocked the intended automated first publish, but the artifact contract was simple enough to execute manually without changing the target architecture.
- The shared VaultConnection warning would have left Argo in an undesirable state, so I fixed the manifest ownership boundary with a dedicated VaultConnection.

### What worked
- Source PR checks for PR #4 passed before merge.
- K3s PR #108 merged and provided the static-site scaffold and Vault files.
- Vault bootstrap completed and included the new `rag-evaluation-storybook` Kubernetes role plus `rag-evaluation-system-gitops-pr` GitHub Actions role.
- Docker pushed `ghcr.io/go-go-golems/rag-evaluation-storybook:sha-09d7628` successfully.
- K3s release PR #109 patched all `sha-*` tokens consistently.
- Publisher Job `publish-rag-evaluation-storybook-sha-09d7628` completed.
- `VaultAuth` and `VaultStaticSecret` became healthy/ready.
- TLS secret `rag-evaluation-storybook-tls` was issued.
- HTTPS returned `HTTP/2 200` and served Storybook HTML.

### What didn't work
- The source workflow run on `main` failed before creating jobs:
  - workflow: `publish-rag-evaluation-storybook`
  - run: `https://github.com/go-go-golems/rag-evaluation-system/actions/runs/27237183294`
  - conclusion: `startup_failure`
  - jobs: `[]`
- Because no job existed, there was no job log to inspect with `gh run view --log-failed`.
- The first Argo sync reported `OutOfSync Healthy` due to a shared `VaultConnection/vault` resource also owned by `dmeta-examples`. This was fixed in PR #110 by renaming Storybook's VaultConnection to `rag-evaluation-storybook-vault`.

### What I learned
- The static artifact contract is robust enough to recover from CI handoff failure: as long as `/site/index.html` can be built and pushed to GHCR, the K3s static-sites flow works.
- Argo shared-resource ownership matters for reusable VSO support resources. Static-site packages should use unique `VaultConnection` names unless the shared resource is intentionally managed by a separate platform app.
- A completed publisher Job and successful HTTPS smoke test are the strongest rollout proof for this pattern.

### What was tricky to build
- The workflow startup failure happened after merge and before any reusable-workflow job logs existed. I did not block the rollout on it because manual publish plus GitOps PR preserved the same deployment artifact and release-token semantics. The follow-up is to debug/fix the reusable workflow invocation or replace it with an in-repo workflow that calls infra-tooling differently.
- The Argo app was healthy but out of sync because of shared `VaultConnection/vault`, which is a subtle ownership issue rather than a runtime failure. The symptom was `SharedResourceWarning` and `VaultConnection/vault` shown as `OutOfSync`; the fix was a dedicated VaultConnection name and matching VaultAuth reference.

### What warrants a second pair of eyes
- Debug the GitHub Actions `startup_failure` for `.github/workflows/publish-rag-evaluation-storybook.yml`; future automated releases are not proven until that is fixed.
- Decide whether all static-site packages should use dedicated VaultConnection names to avoid shared-resource warnings.
- Confirm GHCR package visibility and long-term token scope for image pulls.

### What should be done in the future
- Fix and re-run the source workflow so future Storybook releases are fully automated.
- Consider extracting the final diary into a standalone static Storybook deployment playbook.
- Add a cleanup policy for old `/srv/sites/rag-evaluation-storybook.yolo.scapegoat.dev/releases/sha-*` directories if releases become frequent.

### Code review instructions
- Review K3s PRs #108, #109, and #110 in order.
- Verify Argo status with:
  - `kubectl -n argocd get application rag-evaluation-storybook`
- Verify runtime resources with:
  - `kubectl -n static-sites get job,pod,vaultauth,vaultstaticsecret,secret,ingress | rg 'rag-evaluation-storybook|NAME'`
- Verify the public site with:
  - `curl -fsSI https://rag-evaluation-storybook.yolo.scapegoat.dev/`

### Technical details
- Live URL: `https://rag-evaluation-storybook.yolo.scapegoat.dev/`
- Release image: `ghcr.io/go-go-golems/rag-evaluation-storybook:sha-09d7628`
- K3s final sync revision: `982d63aa2f69d118b7aa7bdd6fc0673beedadc21`
- Successful HTTP status: `HTTP/2 200`
