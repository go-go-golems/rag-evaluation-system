#!/usr/bin/env bash
# Reproduce the non-paid fixture sweep, compact operation custody export, and
# generic researchctl artifact staging/import validation. No provider config or
# source data is copied into this ticket.
set -euo pipefail

repo_root=$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../../../.." && pwd)
researchctl_repo=${RESEARCHCTL_REPO:?set RESEARCHCTL_REPO to the researchctl checkout}
output=${1:-/tmp/rag-ttc-custody-fixture}
spec=/tmp/rag-ttc-custody-spec.json
project=/tmp/rag-ttc-custody-project.yaml
lab_db=/tmp/rag-ttc-custody-lab.sqlite

cd "$repo_root"
GOWORK=off go run "${BASH_SOURCE%/*}/02-build-researchctl-custody-spec.go" >"$spec"
cp "${BASH_SOURCE%/*}/03-researchctl-custody-project.yaml" "$project"
rm -rf "$output"
GOWORK=off go run ./cmd/rag-ttc-v3-sweep \
  --profile fixtures --chunks 16 --concurrency 1,2,4 --maximum-requests 90 \
  --output "$output" --specification "$spec" \
  --researchctl-custody-run-id run_00000000000000000000000000 \
  --researchctl-custody-attempt-id attempt_11111111111111111111111111 \
  --researchctl-custody-external-run-id fixture-custody-e2e \
  --researchctl-custody-recorded-at 2026-07-23T12:00:00Z

cd "$researchctl_repo"
GOWORK=off go build -o /tmp/researchctl-ttc-custody ./cmd/researchctl
rm -f "$lab_db"
/tmp/researchctl-ttc-custody lab init --project "$project" --database "$lab_db" --output json
/tmp/researchctl-ttc-custody experiment import-run "$output/researchctl-run-export.json" \
  --project "$project" --database "$lab_db" --experiment EXP-001 --output json
