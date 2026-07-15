# TTC RAG laboratory

The TTC RAG laboratory is a local web application and API for inspecting immutable retrieval artifacts and append-only experiment runs. It treats a corpus snapshot, chunk set, lexical index, embedding set, retrieval configuration, and evaluation dataset identifier as explicit inputs. A run never overwrites its inputs or prior observations.

## Start the laboratory

Build the SPA and start the server against the local experiment database:

```bash
cd /home/manuel/workspaces/2026-07-13/rag-eval-ttc/rag-evaluation-system
pnpm --dir web build
GOWORK=off go run ./cmd/rag-eval serve --address 127.0.0.1:8772 --db data/rag-eval.db --log-level info
```

Open `http://127.0.0.1:8772`, select **Evaluation**, and use the RAG Laboratory view. The page can create a content-addressed specification, create an append-only run, display its lifecycle events, inspect query traces, and fetch a two-run comparison.

## Baseline currently available

The initial local TTC observation uses:

| Input | Identifier |
| --- | --- |
| corpus snapshot | `sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409` |
| fixed chunk set | `sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392` |
| BM25 artifact | `sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691` |
| 768D embedding set | `sha256:2665c5249b8352ce6904fc00c934534dd179f3eeef0a6a75429a9034be0e03e0` |

The ticket-local scripts create and import this observation:

```bash
GOWORK=off go run ./ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/04-run-immutable-retrieval-traces.go
GOWORK=off go run ./ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/05-score-candidate-retrieval-traces.go
GOWORK=off go run ./ttmp/2026/07/13/RAGEVAL-TTC-LAB-001--ttc-rag-laboratory-baseline-and-immutable-experiment-runs/scripts/06-import-candidate-baseline-run.go
```

The candidate cards are source-validated but not human-frozen. Their dataset identifier is therefore `candidate:ttc-baseline-v1`, not `ttc-baseline-eval-v1`. See the ticket’s human adjudication packet before publishing a benchmark.

## API reference

| Endpoint | Purpose |
| --- | --- |
| `GET /api/v1/lab/catalog` | List immutable snapshots, chunk sets, embedding sets, and BM25 artifacts for the form. |
| `GET, POST /api/v1/lab/specifications` | List or create content-addressed specifications. |
| `POST /api/v1/lab/specifications/{id}/runs` | Create a new append-only observation run. |
| `GET /api/v1/lab/runs` | List runs, optionally with `specification_id`. |
| `GET /api/v1/lab/runs/{id}` | Read run events and terminal summary. |
| `GET /api/v1/lab/runs/{id}/traces` | Read immutable per-query traces. |
| `POST /api/v1/lab/runs/{id}/events` | Append a non-terminal lifecycle event. |
| `POST /api/v1/lab/runs/{id}/traces` | Record one immutable trace per query card. |
| `POST /api/v1/lab/runs/{id}/complete` | Add the only permitted terminal summary. |
| `GET /api/v1/lab/comparison?left=...&right=...` | Read two runs and their query traces together. |

The database enforces append-only semantics with SQLite triggers. An event or trace after a terminal summary fails, and `UPDATE` or `DELETE` of specifications, runs, events, summaries, or traces fails.

## Validation

```bash
GOWORK=off go test ./...
pnpm --dir web typecheck
pnpm --dir web build
```

The web build output is generated into `internal/web/dist`; it is intentionally ignored by Git except for the tracked HTML shell. Build it before running the embedded server locally or in a release pipeline.
