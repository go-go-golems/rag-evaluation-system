# Tasks

## Phase 0 — Recovery and scope control

- [x] 0.1 Confirm current git state, mistaken local commits, and uncommitted files before further changes.
- [x] 0.2 Update the design documentation with the actual go-minitrace/minitracejs schema, DSL query outputs, and widget mapping implications.
- [x] 0.3 Restore the over-removed slides/handouts/course Storybook and nav surfaces while preserving the design note that live/realtime service is out of scope.
- [x] 0.4 Remove only the `litellm-live-service.js` live/realtime service integration and any direct references to it.
- [x] 0.5 Validate cleanup with typecheck/builds/tests appropriate to changed files.
- [x] 0.6 Commit cleanup in focused intervals and record each interval in the diary.

## Phase 1 — Minitrace-driven widget mapping design

- [x] 1.1 Expand the design guide with the authoritative minitrace session fields: session metadata, turns, tool calls, per-turn usage, annotations, and metrics.
- [x] 1.2 Document minitracejs DSL outputs relevant to widgets: `TranscriptRows`, `TurnRows`, `ToolRows`, `TimelineRows`, `TokenUsageRows`, and `TurnFrames`.
- [x] 1.3 Define precise mapping tables from minitrace rows to `TranscriptMessage`, `TranscriptAnnotation`, `ContextWindowPart`, and `ContextStyleSet` entries.
- [x] 1.4 Mark any unclear go-minitrace docs/API areas as documentation follow-ups instead of guessing.

## Phase 2 — ClubMeetup/minitrace-viz local widget improvements

- [x] 2.1 Locate the ClubMeetup/minitrace-viz code path that consumes go-minitrace or minitracejs output.
- [x] 2.2 Implement all improvements possible locally: convert minitrace transcript/tool/token rows into uploaded-session context block snapshots and transcript widgets without upstream API changes.
- [x] 2.3 Add local fixtures or examples for uploaded session JSON/minitrace-derived rows.
- [x] 2.4 Validate local minitrace-viz output and document any limits caused by current upstream widget contracts.
- [x] 2.5 Commit local-only widget improvements and diary updates.

## Phase 3 — Upstream rag-evaluation-site design-system changes

- [x] 3.1 Re-read `packages/rag-evaluation-site/GUIDELINES.md` before any upstream component changes.
- [x] 3.2 Identify upstream gaps that cannot be solved locally, such as grouped turn separators, block aggregation, or richer context-detail affordances.
- [x] 3.3 Implement minimal upstream component/type/story changes following the design-system structure and style guidelines.
- [x] 3.4 Validate with `pnpm --dir packages/rag-evaluation-site typecheck` and `pnpm --dir packages/rag-evaluation-site build-storybook`.
- [x] 3.5 Commit upstream widget package changes and diary updates.

## Phase 4 — Goja DSL/schema support

- [x] 4.1 Identify whether the current `context_window.dsl` helpers can express all minitrace-derived widget IR without changes.
- [ ] 4.2 Add only necessary DSL/schema helpers for repeated minitrace mapping patterns; avoid duplicating app-specific logic in the DSL.
- [ ] 4.3 Update TypeScript descriptors, provider docs, schema tests, and Go tests as needed.
- [ ] 4.4 Validate Go packages with focused `go test` commands and frontend typecheck/storybook if contracts changed.
- [ ] 4.5 Commit DSL/schema changes and diary updates.

## Phase 5 — Final validation and delivery

- [ ] 5.1 Run final repo/workspace validation for all changed areas.
- [ ] 5.2 Update ticket design docs, diary, changelog, and related-file links.
- [ ] 5.3 Run `docmgr doctor --ticket CTX-WINDOW-BLOCK-VIZ --stale-after 30`.
- [ ] 5.4 Provide final handoff with commits, validations, remaining risks, and follow-up documentation gaps.
