# TTC RAG laboratory

The TTC RAG integration has two explicit authorities. Rag-evaluation-system owns immutable corpus/chunk/index/embedding/evaluation artifacts and retrieval semantics. Researchctl owns specifications, runs, attempts, retries, observations, terminal state, verified artifacts, export, and import.

The rag-eval web application exposes only a read-only domain artifact catalog. The former local run/specification API and database-backed JavaScript laboratory were removed after external-import and native-rerun parity review.

## Inspect domain artifacts

Build the SPA and start the server against the TTC catalog:

```bash
pnpm --dir web build
go run ./cmd/rag-eval serve --address 127.0.0.1:8772 \
  --db data/rag-eval.db --log-level info
```

Open `http://127.0.0.1:8772` and select **Evaluation**. The page reports immutable corpus snapshots, chunk sets, embedding sets, and BM25 artifacts. It cannot create or mutate runs.

## Execute through researchctl

Build the observation-only worker:

```bash
go build -o .bin/rag-lab-worker ./cmd/rag-lab-worker
```

Execute a pure `rag-retrieval-spec/v1` program with explicit input references:

```bash
researchctl experiment run-rag experiment.js \
  --project project.yaml \
  --experiment-id EXP-RAG \
  --inputs inputs.json \
  --ttc-database data/rag-eval.db \
  --runner .bin/rag-lab-worker \
  --runner-arg=--db --runner-arg=data/rag-eval.db \
  --runner-has-embedder --timeout 10m

researchctl lab runs list --project project.yaml --output json
researchctl lab runs show RUN_ID --project project.yaml --output json
```

The worker opens the source catalog with WAL-aware `mode=ro` and query-only access. It emits events, query traces, metrics, and artifacts through `researchctl-rag-runner-stdio/v1`; it has no researchctl database handle.

## Frozen domain identities

| Input | Identifier |
| --- | --- |
| selected 200-document corpus snapshot | `sha256:be434a1422487d33e324b5f3833067dcc530efab2df0fea2f7e7bfa9ca86f409` |
| fixed chunk set | `sha256:ef7bdab76583f092d7bc50c9f501fe8c17739d395fcb37d0eaaba5a09c7c9392` |
| BM25 artifact | `sha256:cf6491873ec521135ade41000800751dc8eeaecba52dabbeacda1cf530f7b691` |
| 768D embedding set | `sha256:2665c5249b8352ce6904fc00c934534dd179f3eeef0a6a75429a9034be0e03e0` |
| baseline dataset | `candidate:ttc-baseline-v1` |
| corrected expansion dataset | `candidate:ttc-expansion-v1` |

Candidate datasets are provisional and not human-adjudicated holdouts.

## API reference

| Endpoint | Purpose |
| --- | --- |
| `GET /api/v1/lab/catalog` | Read immutable snapshots, chunk sets, embedding sets, and BM25 artifacts. |

Requests to the retired `/api/v1/lab/specifications`, `/api/v1/lab/runs`, and `/api/v1/lab/comparison` routes return 404. Historical SQLite tables remain append-only migration history but have no supported writer.

## Validation

```bash
go test ./...
pnpm --dir web typecheck
pnpm --dir web build
```

The web build output is generated into `internal/web/dist`; build it before running the embedded server locally or in a release pipeline.
