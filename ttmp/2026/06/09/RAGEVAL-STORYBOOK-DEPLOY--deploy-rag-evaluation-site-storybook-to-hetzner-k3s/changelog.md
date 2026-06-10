# Changelog

## 2026-06-09

- Initial workspace created


## 2026-06-09

Created intern-oriented Storybook deployment design guide and investigation diary for K3s static-sites rollout.

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY--deploy-rag-evaluation-site-storybook-to-hetzner-k3s/design-doc/01-intern-guide-deploy-rag-evaluation-site-storybook-to-hetzner-k3s.md — Main deliverable
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY--deploy-rag-evaluation-site-storybook-to-hetzner-k3s/reference/01-investigation-diary.md — Chronological evidence and continuation notes


## 2026-06-09

Uploaded Storybook deployment guide bundle to reMarkable at /ai/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY.

### Related Files

- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY--deploy-rag-evaluation-site-storybook-to-hetzner-k3s/design-doc/01-intern-guide-deploy-rag-evaluation-site-storybook-to-hetzner-k3s.md — Uploaded in reMarkable bundle
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY--deploy-rag-evaluation-site-storybook-to-hetzner-k3s/reference/01-investigation-diary.md — Uploaded in reMarkable bundle


## 2026-06-09

Implemented Storybook static artifact packaging, source CI/GitOps handoff, K3s static-site manifests, and Vault role/policy files; updated guide and diary with commits and validation results.

### Related Files

- /home/manuel/code/wesen/2026-03-27--hetzner-k3s/gitops/kustomize/rag-evaluation-storybook/publish-job.yaml — K3s static-site implementation
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/Dockerfile.storybook-static — Source packaging implementation
- /home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/ttmp/2026/06/09/RAGEVAL-STORYBOOK-DEPLOY--deploy-rag-evaluation-site-storybook-to-hetzner-k3s/reference/01-investigation-diary.md — Playbook-quality implementation diary


## 2026-06-09

Deployed web Storybook at https://rag-evaluation-page.storybook.yolo.scapegoat.dev/ with wildcard Storybook DNS/TLS and recorded rollout diary.

### Related Files

- /tmp/rag-web-storybook-source/Dockerfile.web-storybook-static — Packages the web Storybook static output as /site
- /tmp/storybook-dns-tf/dns/zones/scapegoat-dev/envs/prod/main.tf — Adds wildcard Storybook DNS record
- /tmp/storybook-wildcard-k3s/gitops/kustomize/static-sites-host/storybook-wildcard-ingress.yaml — Adds wildcard Storybook routing to static-sites-host


## 2026-06-09

Fixed web Storybook MIME errors on /%60/ paths by publishing sha-e78da8f with root-absolute Storybook asset URLs.

### Related Files

- /tmp/rag-web-storybook-source/Dockerfile.web-storybook-static — Rewrites generated Storybook asset URLs to host-root absolute paths
- /tmp/storybook-wildcard-k3s/gitops/kustomize/rag-evaluation-page-storybook/publish-job.yaml — Publishes fixed sha-e78da8f image


## 2026-06-09

Fixed remaining backtick-path preview issue by publishing sha-eb4b80c, which redirects /%60?path=... to /?path=... before Storybook manager starts.

### Related Files

- /tmp/rag-web-storybook-source/Dockerfile.web-storybook-static — Injects early Storybook entry path normalization
- /tmp/storybook-wildcard-k3s/gitops/kustomize/rag-evaluation-page-storybook/publish-job.yaml — Publishes fixed sha-eb4b80c image

