---
Title: Native Glazed CLI migration plan
Ticket: RAGEVAL-009
Status: active
Topics:
    - rag
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: repo://Makefile
      Note: Defines the current Glazed lint scope
    - Path: repo://cmd/rag-eval/cmds/corpus/import_ttc.go
      Note: Existing native Glazed pattern
    - Path: repo://cmd/rag-eval/cmds/source/create.go
      Note: Existing native Glazed pattern to preserve
    - Path: repo://cmd/rag-eval/main.go
      Note: Registers the temporary adapter at the application root
ExternalSources: []
Summary: Replace all application-owned direct Cobra flag/handler definitions in rag-eval with typed native Glazed commands and enforce the invariant in lint.
LastUpdated: 2026-07-18T19:50:52-04:00
WhatFor: A safe, command-by-command native Glazed migration.
WhenToUse: When porting or reviewing rag-eval CLI commands.
---





# Native Glazed CLI migration plan

## Executive Summary

`rag-eval` currently contains direct Cobra command implementations across source, corpus, chunk, document, embedding, search, workflow, study, preview, and serve. A temporary `glazedcobra` adapter makes Glazed the outer parser, but it does not meet the desired end state: application code still defines Cobra flags and handlers.

This ticket replaces every leaf with a typed Glazed command, removes the adapter, and extends the mandatory `glazed-lint` scope to `./cmd/rag-eval/...`. Cobra remains the framework transport used by Glazed; it must not remain an application-owned command-definition API.

## Problem Statement and Scope

`cmd/rag-eval/main.go` currently registers ten direct-Cobra command trees through `glazedcobra.WrapTree`. `cmd/rag-eval/cmds/glazedcobra/wrap.go` copies their flag surface dynamically and delegates to their `Run`/`RunE` handlers. This preserves behavior but weakens types, obscures command schemas, and allows direct Cobra declarations to persist.

In scope: all listed rag-eval leaves, root registration, removal of the adapter, tests, and Makefile enforcement. Out of scope: removing Cobra from Glazed itself, changing business behavior, or changing the public flag/argument contract unless a command's existing behavior is demonstrably broken.

## Current-State Evidence

- `cmd/rag-eval/main.go:49-59` wraps direct Cobra command groups instead of registering native Glazed commands.
- `cmd/rag-eval/cmds/glazedcobra/wrap.go` translates pflags into dynamic Glazed fields and invokes legacy handlers.
- `Makefile:17,29,79-84` excludes `cmd/rag-eval/...` from the Glazed lint target.
- Running `GOWORK=off go vet -vettool=/tmp/glazed-lint ./cmd/rag-eval/...` reports direct flag definitions in every current command group.

## Proposed Architecture

Each leaf has:

1. a command struct embedding `*cmds.CommandDescription`;
2. a typed settings struct with `glazed` tags;
3. fields declared through `cmds.WithFlags(fields.New(...))` and positional arguments through `cmds.WithArguments`;
4. `Run` or `RunIntoGlazeProcessor`, which decodes values with `DecodeSectionInto`; and
5. a group root that builds the leaf's Cobra transport with `cli.BuildCobraCommandFromCommand`.

```
user arguments -> Glazed schema/parser -> typed settings -> domain service -> Glazed output
```

No leaf calls `Flags()`, `StringVar`, `RunE`, or accesses Cobra values. Root-level logging/help remains Glazed's supported Cobra integration.

### Decision: Native rewrites rather than the adapter

- **Context:** The adapter preserves old behavior but direct Cobra definitions remain and `glazed-lint` correctly flags them.
- **Options considered:** Retain the adapter; suppress lint; port all leaves natively.
- **Decision:** Port all leaves natively and delete the adapter.
- **Rationale:** This makes schemas typed, inspectable, testable, and mechanically enforced.
- **Consequences:** The migration is sizeable and must preserve each command's established input/output contract.
- **Status:** accepted.

## Implementation Plan

1. Document the baseline and task inventory.
2. Port groups in dependency-light order: source, corpus, document, chunk, embedding, search.
3. Port stateful/complex groups: workflow, study, preview, serve.
4. Remove `cmds/glazedcobra`, register native group roots directly, and run the CLI help regression checks.
5. Add `./cmd/rag-eval/...` to `GLAZED_LINT_DIRS`; remove the application-command allow-list; require passing `make lint`.

## Test Strategy

- Retain and update each group’s existing unit tests.
- Add command-level tests that execute Glazed-built commands with representative flags and positional arguments.
- Compare `--help` for each leaf before/after where it is part of the supported contract.
- Run `GOWORK=off go test ./cmd/rag-eval/... -count=1`, `make lint`, `make test`, and `make logcopter-check`.

## Risks and Open Questions

The largest risk is accidental semantic drift in commands with many flags, especially workflow intake and study execution. Porting must decode typed settings once and pass values to extracted domain helpers; no business logic should be reimplemented merely to satisfy the command framework.

## References

- `cmd/rag-eval/main.go`
- `cmd/rag-eval/cmds/glazedcobra/wrap.go`
- `Makefile`
