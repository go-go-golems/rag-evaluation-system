---
Title: mimimi Ollama tunnel operator playbook
Ticket: RAGEVAL-RAG-DSL-001
Status: active
Topics:
    - rag
    - rag-eval
    - dsl
    - fluent-builder
    - goja
    - xgoja
    - javascript
    - typescript
    - intern-guide
    - playground
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/02-geppetto-ollama-embedding-probe.go
      Note: Direct Geppetto compatibility probe for the tunneled local base URL
    - Path: repo://ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/04-run-immutable-retrieval-traces.go
      Note: Historical vector and RRF trace runner using the tunnel endpoint
    - Path: ws://geppetto/pkg/embeddings/ollama.go
      Note: Geppetto provider endpoint and base URL behavior
ExternalSources: []
Summary: Verified operator procedure for exposing the Mac-hosted Ollama embedding service as a local-only endpoint for reproducible TTC RAG experiments.
LastUpdated: 2026-07-15T14:31:26.995029886-04:00
WhatFor: Recover and verify the SSH loopback tunnel used by Geppetto and the RAG laboratory without exposing Ollama on the LAN.
WhenToUse: Use before a live embedding build or vector/RRF experiment that needs the Mac-hosted nomic-embed-text model.
---


# mimimi Ollama tunnel operator playbook

## Goal

Use the Ollama instance running on Manuel's Mac for fast local embedding while
keeping the model HTTP API private. The laboratory process must talk to a
loopback-only endpoint on the workstation; it must never need to know the Mac
IP address or expose an Ollama listener to the LAN.

This is an operator procedure, not part of an immutable experiment
specification. The selected model, dimensions, corpus artifacts, and retrieval
configuration belong in immutable artifacts and runs. The SSH hostname and
tunnel lifetime are runtime infrastructure.

## Context

### Verified topology

The original setup was recorded as the “mimimi” tunnel. On 2026-07-15 the
reachable SSH alias was verified as `mimimi-2.local`; `mimimi.local` did not
resolve from this workstation. Both aliases use SSH user `manuel`, but only the
former is presently operational. Do not silently substitute a network address:
repair the desired DNS alias separately if `mimimi.local` is required.

The remote Ollama HTTP API was verified live at
`mimimi-2.local:127.0.0.1:11434`. Its installed model list contains
`nomic-embed-text:latest`. The Mac shell does not put an `ollama` executable on
`PATH`, so API health is the authoritative first check; the previous recovery
path used `/Applications/Ollama.app/Contents/Resources/ollama` when a CLI was
needed.

```text
rag-eval / Geppetto / ticket scripts
              |
              | http://127.0.0.1:11435
              v
local SSH client in tmux
              |
              | encrypted SSH forwarding
              v
mimimi-2.local:127.0.0.1:11434 ──► Ollama ──► nomic-embed-text:latest (768D)
```

The two loopback bindings are deliberate:

- The remote Ollama server remains local to the Mac.
- The local forwarding listener remains local to the workstation.
- `11435` prevents a conflict with a locally installed Ollama defaulting to
  `11434`.
- No `-g`, `0.0.0.0:...`, reverse forwarding, or firewall change is needed.

### Established laboratory contract

The completed TTC embedding set was produced through this route with
`nomic-embed-text`, 768 dimensions, and the local base URL
`http://127.0.0.1:11435`. Relevant implementation points are:

- `geppetto/pkg/embeddings/ollama.go` calls
  `POST {base-url}/api/embeddings` for the Go embedding provider. Its normal
  default is `http://localhost:11434`, so the tunnel must be supplied
  explicitly.
- The historical probe
  `RAGEVAL-TTC-LAB-001/scripts/02-geppetto-ollama-embedding-probe.go` accepts
  `--base-url` and verifies the provider plus a real vector.
- The historical trace runner
  `RAGEVAL-TTC-LAB-001/scripts/04-run-immutable-retrieval-traces.go` defaults
  to `http://127.0.0.1:11435` because it was designed for this tunnel.
- The present `raglab.Executor` requires an explicit query embedder for vector
  or RRF execution. The endpoint is runtime capability, not a hidden property
  of an embedding-set artifact.

## Quick Reference

| Item | Value |
| --- | --- |
| Reachable Mac SSH alias | `mimimi-2.local` |
| Intended but currently unresolved alias | `mimimi.local` |
| Remote Ollama listener | `127.0.0.1:11434` |
| Local tunnel listener | `127.0.0.1:11435` |
| Installed embedding model | `nomic-embed-text:latest` |
| Expected vector dimensions | `768` |
| Tunnel tmux session | `rag-ollama-mimimi` |

### 1. Inspect remote health first

This command does not create a tunnel. It verifies that the remote API is
reachable and shows its installed models without reading any provider profile
or secret:

```bash
ssh -o BatchMode=yes -o ConnectTimeout=10 mimimi-2.local \
  'curl -fsS --max-time 5 http://127.0.0.1:11434/api/tags | jq -r "[.models[].name] | join(\",\")"'
```

Expected output includes `nomic-embed-text:latest`. A nonzero result means
either SSH, the Mac, or Ollama is unavailable. Do not start a full corpus build
until this succeeds.

### 2. Start the tunnel in tmux

Run this from the `rag-evaluation-system` worktree. It fails at startup if the
local port is already in use or the remote destination cannot be established.

```bash
tmux new-session -d -s rag-ollama-mimimi \
  'exec ssh -N -o ExitOnForwardFailure=yes -o ServerAliveInterval=30 -o ServerAliveCountMax=3 -L 127.0.0.1:11435:127.0.0.1:11434 mimimi-2.local'
```

If the session name already exists, inspect it before replacing it:

```bash
tmux capture-pane -pt rag-ollama-mimimi:0.0 -S -80
```

### 3. Verify the local endpoint and model

```bash
curl -fsS --max-time 5 http://127.0.0.1:11435/api/tags \
  | jq -r '[.models[].name] | join(",")'
```

The expected model name is `nomic-embed-text:latest`. This checks the actual
path used by the laboratory, not merely the remote server.

For a bounded end-to-end request, use Ollama's current embedding endpoint:

```bash
curl -fsS --max-time 30 http://127.0.0.1:11435/api/embed \
  -H 'Content-Type: application/json' \
  -d '{"model":"nomic-embed-text","input":"tunnel health check"}' \
  | jq '[.embeddings[0] | length]'
```

The expected result is `[768]`. The Geppetto provider itself uses its
compatibility endpoint `/api/embeddings`; validate that exact API with the
ticket probe before a production experiment.

### 4. Use it explicitly

Every command that constructs a Geppetto Ollama provider must receive the
local base URL:

```bash
GOWORK=off go run \
  ./ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/02-geppetto-ollama-embedding-probe.go \
  --base-url http://127.0.0.1:11435 \
  --engine nomic-embed-text \
  --batch-size 1
```

For the next fresh vector/RRF study, create the executor's explicit
`QueryEmbedder` from the same configured Geppetto provider. Do not copy an old
trace and call it fresh, and do not make `127.0.0.1:11435` a committed default
in an immutable experiment specification.

### 5. Stop cleanly

After a short experiment, terminate only the local forwarding session:

```bash
tmux kill-session -t rag-ollama-mimimi
```

This does not stop Ollama on the Mac. Before reclaiming a stuck local port,
inspect the session and port owner; do not indiscriminately terminate remote
processes.

## Recovery and troubleshooting

| Symptom | Evidence command | Likely cause | Safe next action |
| --- | --- | --- | --- |
| `mimimi.local` cannot resolve | `ssh -G mimimi.local` / SSH error | Current DNS or mDNS alias is absent | Use verified `mimimi-2.local`; repair the alias as a separate infrastructure task. |
| SSH cannot connect | `ssh -o BatchMode=yes -o ConnectTimeout=10 mimimi-2.local true` | Mac offline, mDNS unavailable, or SSH service unavailable | Confirm the Mac is reachable; do not alter forwarding options. |
| Remote `/api/tags` fails | remote health command above | Ollama app/service stopped | Start Ollama on the Mac, then repeat the remote health check. Its CLI may be `/Applications/Ollama.app/Contents/Resources/ollama`. |
| Tunnel tmux exits immediately | `tmux capture-pane -pt rag-ollama-mimimi:0.0 -S -80` | Port collision, SSH error, or remote port inaccessible | Read the captured SSH error; use another local port only when all consumers are updated explicitly. |
| Local `/api/tags` refuses connection | `curl -v http://127.0.0.1:11435/api/tags` | Tunnel absent/stale | Recreate the tunnel after remote health succeeds. |
| Model missing | local `/api/tags` output | Ollama registry changed | Pull `nomic-embed-text` on the Mac, then re-verify its 768D embedding response. |
| Build appears to stop | `tmux capture-pane` and process inspection | The caller stopped observing a child, or inference is simply long-running | Treat the tmux pane and the persisted run events as authority; avoid duplicate full builds. |

### Why tmux is mandatory for long work

The earlier TTC corpus build showed that foreground observation can end before
the actual Go process or embedding work ends. A tmux session makes the tunnel
and long running embedding build independently inspectable. Use a separate
session for the build itself and inspect both panes before declaring a run
failed.

```text
tmux:rag-ollama-mimimi  ── SSH tunnel lifecycle
tmux:rag-embedding-run ── laboratory command and its logs
SQLite run events       ── durable experiment lifecycle evidence
```

## Preflight checklist for a fresh vector/RRF run

- [ ] `mimimi-2.local` remote `/api/tags` lists `nomic-embed-text:latest`.
- [ ] The tmux tunnel is alive and local `/api/tags` succeeds on port 11435.
- [ ] A bounded embedding request returns exactly 768 dimensions.
- [ ] The configured Geppetto provider has `Type: ollama`,
  `Engine: nomic-embed-text`, `Dimensions: 768`, and
  `BaseURL: http://127.0.0.1:11435`.
- [ ] The planned run has a new immutable run ID, frozen evaluation manifest,
  explicit `QueryEmbedder`, and a named vector/RRF configuration.
- [ ] Metrics will record quality, latency, model/provider identity, storage,
  and the bounded no-provider-billing cost statement.

## Related

- [Implementation diary](02-implementation-diary.md), especially RAGEVAL-TTC-LAB-001 Step 11, for the original migration and completed 2,024-vector artifact.
- [RAG module API specification](01-rag-laboratory-javascript-module-api-specification.md), for the `rag` builder and its immutable-plan boundary.
- `RAGEVAL-TTC-LAB-001/scripts/02-geppetto-ollama-embedding-probe.go` for the direct Geppetto compatibility test.
- `RAGEVAL-TTC-LAB-001/scripts/04-run-immutable-retrieval-traces.go` for the historical vector/RRF trace route and explicit `--base-url` flag.
