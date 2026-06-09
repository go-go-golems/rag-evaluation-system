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
