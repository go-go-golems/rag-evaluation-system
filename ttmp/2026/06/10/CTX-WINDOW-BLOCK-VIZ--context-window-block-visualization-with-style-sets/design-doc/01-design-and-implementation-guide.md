---
Title: Design and Implementation Guide
Ticket: CTX-WINDOW-BLOCK-VIZ
Status: active
Topics:
    - frontend
    - visualization
    - context-window
    - storybook
    - widget-ir
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles:
    - Path: ../../../../../../../go-minitrace/pkg/doc/js-api-reference.md
      Note: minitracejs importer/session/view DSL used by uploaded-session flow
    - Path: ../../../../../../../go-minitrace/pkg/doc/minitrace-schema.md
      Note: Authoritative minitrace session JSON schema used for widget mapping
    - Path: ../../../../../../../go-minitrace/pkg/minitracejs/query_view_session.go
      Note: Concrete SQL recipe outputs for TranscriptRows
    - Path: README.md
      Note: Repository architecture and Widget IR runtime overview used by the guide
    - Path: packages/rag-evaluation-site/README.md
      Note: Package docs now describe Session Workspace and upload/visualize/transcript scope
    - Path: packages/rag-evaluation-site/src/components/molecules/ContextStackDiagram/ContextStackDiagram.tsx
      Note: Existing stacked block renderer and readable height scaling
    - Path: packages/rag-evaluation-site/src/components/molecules/ContextStripDiagram/ContextStripDiagram.tsx
      Note: Existing ordered token-proportional strip renderer and readability floor
    - Path: packages/rag-evaluation-site/src/components/molecules/ContextTreemap/ContextTreemap.tsx
      Note: Existing token-area renderer and headroom filtering behavior
    - Path: packages/rag-evaluation-site/src/components/organisms/ContextDiagramPanel/ContextDiagramPanel.tsx
      Note: Existing view switcher
    - Path: packages/rag-evaluation-site/src/components/organisms/ContextUploadDropArea/ContextUploadDropArea.tsx
      Note: Presentational upload entrypoint; parser remains future work
    - Path: packages/rag-evaluation-site/src/context/styles.ts
      Note: Defines palette style-set helpers
    - Path: packages/rag-evaluation-site/src/context/types.ts
      Note: Defines ContextStyleSet
    - Path: packages/rag-evaluation-site/src/widgets/WidgetRenderer.context-diagrams.stories.tsx
      Note: Existing Storybook examples for content blocks and custom legends
    - Path: packages/rag-evaluation-site/src/widgets/WidgetRenderer.course-handout.stories.tsx
      Note: Storybook surface renamed and scoped to Session Workspace
    - Path: packages/rag-evaluation-site/src/widgets/ir.ts
      Note: Widget IR prop contracts for context diagram and transcript widgets
    - Path: pkg/widgetdsl/module.go
      Note: Goja DSL helpers for style sets
    - Path: pkg/xgoja/providers/widgetsite/provider.go
      Note: xgoja provider registration for context_window.dsl
ExternalSources: []
Summary: Plan for visualizing uploaded coding-agent sessions as colored, token-proportional context-window blocks using ContextStyleSet legends and Widget IR.
LastUpdated: 2026-06-10T15:06:51.76997905-04:00
WhatFor: Use before implementing the uploaded-session context-window visualization; explains the existing style-set/Widget IR/DSL system and proposes a phased implementation.
WhenToUse: Use when adding real session parsing, block classification, turn paging, or new context-window widget variants.
---




# Design and Implementation Guide

## Executive Summary

This ticket is a planning-only ticket for using the new `styleKey + ContextStyleSet` functionality to visualize an uploaded coding-agent session as a context window made of colored, token-proportional content blocks. The intended viewer should show block categories such as system, user, agent, tool call, tool result, file read/write, and thinking/scratchpad. The current product scope is intentionally limited to the uploaded-session flow: **Upload → Visualize / Transcript**. Live-session or course-style session views should not be exposed unless a later ticket explicitly reopens that scope.

The existing frontend already has the most important building blocks:

- `ContextWindowPart` carries a required `styleKey`, token count, metadata, and optional preview/content fields.
- `ContextStyleSet` defines visual styles and legend labels independent of baked-in categories.
- `ContextDiagramPanel` can switch between strip, stack, budget, and treemap renderers.
- Widget IR and the Goja DSL can already emit `ContextDiagramPanel`, `ContextStyleSet`, `ContextWindowSnapshot`, and `ContextWindowPart` objects.

The missing piece is not a palette system. The missing piece is a session-to-context-window normalization layer: a parser/classifier that turns the uploaded session into stable block records, estimates or reads token sizes, groups records by turn where possible, and emits one or more `ContextWindowSnapshot` objects plus a matching `ContextStyleSet`.

No implementation should start until the team agrees on the uploaded session schema and on how much turn boundary metadata is available. If the uploaded data has reliable turn IDs and message roles, the first implementation can page by turn. If it only contains flattened text or minitrace spans without context-window membership, the first implementation should provide a global block map and record the schema gap as a design-system/runtime requirement.

## Problem Statement and Scope

The user wants the uploaded session visualization to move from generic context-window examples to a content-aware view. Instead of showing broad categories such as "conversation" or "result" only, the view should expose what actually occupies the model context window:

- system prompt and developer instructions,
- user turns,
- assistant/agent responses,
- tool calls,
- tool results,
- file reads,
- file writes/patches,
- thinking/scratchpad/generated reasoning,
- summaries, evictions, and free/headroom when known.

The requested visualization has two core visual rules:

1. **Color and pattern encode block type.** The legend must be caller-defined and should use the new style-set functionality rather than hardcoded component categories.
2. **Size encodes token size.** A block with 10,000 tokens should occupy much more visual area/width/height than a 100-token block, subject to readability floors so very small blocks remain visible.

The requested interaction should start with one supported mode and one explicitly deferred mode:

1. **Whole-window uploaded-session mode:** show all known context-window blocks at once after an upload.
2. **Deferred turn-aware mode:** point out individual turns through labels/metadata where possible, but do not expose a live/page-through session view unless a later ticket explicitly approves it.

This document intentionally does not prescribe code changes to make immediately. It is an intern-facing implementation guide for a future implementation pass.

## Terms

- **Uploaded session:** The JSON or trace file uploaded by the user. In the current package the upload UI accepts `.json` files, but this guide treats the exact schema as unresolved until inspected.
- **Context window snapshot:** A serializable `ContextWindowSnapshot` object with a token limit and an ordered list of `ContextWindowPart` objects.
- **Block:** A semantic unit in the context window, such as one user message, one tool call, one file read result, or one generated summary.
- **Style key:** A string on each block (`styleKey`) that selects a visual style from the companion `ContextStyleSet`.
- **Turn:** A user/assistant/tool interaction grouping. A turn may contain one user message, one or more assistant messages, tool calls, and tool results depending on the source format.
- **Page-through view:** A UI state where the user navigates turn-by-turn while seeing the context window at that turn.

## Current-State Architecture

### Repository and runtime orientation

The repository is a RAG evaluation system with a Go backend, a React Widget IR frontend, and a Goja DSL for producing Widget IR. The top-level README describes the project as including an "interactive playground" that renders structured Widget IR with semantic panels for context windows, transcripts, annotations, slides, and handouts (`README.md:7-13`). It also explains that the server serves a React UI at `/` and REST endpoints under `/api/v1/*` (`README.md:82-95`).

The same README maps the relevant directories:

```text
pkg/widgetdsl/                         Goja Widget IR DSL
pkg/widgetschema/                      Widget IR schema validation
pkg/xgoja/providers/widgetsite/        xgoja Widget DSL provider + help docs
packages/rag-evaluation-site/          React component library
packages/rag-evaluation-site/src/widgets/ Widget IR types and renderer
```

This matters because the future implementation can happen at two levels:

- frontend-only fixtures/stories, useful for proving the design quickly;
- runtime Widget IR generation, useful for the uploaded session application.

### Context data model

`packages/rag-evaluation-site/src/context/types.ts` defines the data model consumed by the diagram widgets.

Important source facts:

- `ContextVisualStyle` includes `fill`, `line`, `stroke`, `labelColor`, `pattern`, and stroke flags (`types.ts:25-35`).
- `ContextLegendItemSpec` includes `id`, `label`, optional `styleKey`, `order`, and `hidden` (`types.ts:37-44`).
- `ContextStyleSet` includes `styles`, `legend`, fallback style, and size preferences (`types.ts:46-54`).
- `ContextWindowPart` is the atomic block contract: `id`, `label`, `styleKey`, `tokens`, optional content preview/content, token offsets, source ID, and metadata (`types.ts:56-68`).
- `ContextWindowSnapshot` wraps the window title, limit, ordered parts, selection, and metadata (`types.ts:70-78`).
- `TranscriptMessage` currently knows role and optional tokens, but not detailed tool/file subtype (`types.ts:101-110`).

This model is already suitable for content blocks. The future work should not reintroduce legacy `kind` semantics. It should produce `styleKey` values such as `system`, `user`, `agent`, `tool_call`, `tool_result`, `file_read`, `file_write`, `thinking`, `summary`, and `headroom`.

### Style sets and palettes

`packages/rag-evaluation-site/src/context/styles.ts` is the style engine. It provides palette definitions, conversion to CSS variables, style resolution, legend filtering, and palette-derived style-set creation.

Important source facts:

- `contextVisualStyleToCssVars` maps a `ContextVisualStyle` into CSS variables such as `--ctx-fill`, `--ctx-line`, `--ctx-stroke`, and `--ctx-label` (`styles.ts:63-72`).
- `resolveContextVisualStyle` looks up a style by `styleKey` and falls back to a configured fallback style or an overflow/error style (`styles.ts:78-91`).
- `legendItemsForStyleSet` hides legend entries where `hidden` is true and sorts by order (`styles.ts:93-97`).
- `createContextStyleSetFromPalette` builds a `ContextStyleSet` from entries using palette accents, fill percentages, and patterns (`styles.ts:124-150`).
- `transcriptStyleSet` already has role/tone entries for `system`, `developer`, `user`, `assistant`, `tool`, `result`, `context`, `generated`, `active`, and related hidden chrome styles (`styles.ts:152-176`).
- `defaultContextStyleSet` has context-window categories including `system`, `tool`, `result`, `generated`, `active`, `empty`, and `other` (`styles.ts:178-201`).

The future implementation should probably define a new session-specific style set rather than overloading `defaultContextStyleSet`. The block vocabulary is more detailed than the current default and should have explicit labels such as "Tool result" and "File read".

### Existing diagram widgets

The current diagram components already implement most of the requested proportional rendering.

`ContextStripDiagram`:

- receives `snapshot` and `styleSet` (`ContextStripDiagram.tsx:5-10`),
- maps `snapshot.parts` to horizontal segments (`ContextStripDiagram.tsx:27-48`),
- computes `width = Math.max(2, (part.tokens / totalWidth) * 100)` (`ContextStripDiagram.tsx:30`),
- hides labels unless the segment is at least 7% wide (`ContextStripDiagram.tsx:45`).

`ContextStackDiagram`:

- receives `snapshot` and `styleSet` (`ContextStackDiagram.tsx:5-9`),
- maps visible parts into stacked rows/layers (`ContextStackDiagram.tsx:26-48`),
- computes height with a readability floor and ceiling: `Math.max(26, Math.min(74, 22 + (part.tokens / maxTokens) * 52))` (`ContextStackDiagram.tsx:30`).

`ContextTreemap`:

- receives `snapshot` and `styleSet` (`ContextTreemap.tsx:5-9`),
- filters out headroom parts using style-set-aware helpers (`ContextTreemap.tsx:20-24`),
- computes tile basis with a readability floor: `Math.max(7, (part.tokens / total) * 100)` (`ContextTreemap.tsx:29`).

`ContextDiagramPanel`:

- owns the view switcher (`strip`, `stack`, `budget`, `treemap`) (`ContextDiagramPanel.tsx:21`, `ContextDiagramPanel.tsx:76-84`),
- uses `legendItemsForStyleSet` to hide internal legend entries (`ContextDiagramPanel.tsx:74-85`),
- can show selected part details including note, preview, and metadata (`ContextDiagramPanel.tsx:27-50`, `ContextDiagramPanel.tsx:86`).

The existing panel is a strong starting point for whole-window mode. It does not currently provide turn paging, multi-snapshot browsing, grouped turn separators inside a strip, or nested blocks within a turn.

### Widget IR contract

`packages/rag-evaluation-site/src/widgets/ir.ts` defines serializable Widget IR props.

Important source facts:

- Widget types include `ContextLegend`, `ContextBudgetBar`, `ContextStripDiagram`, `ContextStackDiagram`, `ContextTreemap`, and `ContextDiagramPanel` (`ir.ts:24-39`).
- Each context diagram widget prop interface requires both `snapshot` and `styleSet` (`ir.ts:204-238`).
- Transcript widgets accept transcript messages and an optional style set (`ir.ts:249-260` and nearby interfaces).

This means the uploaded-session renderer can emit a pure JSON/Widget IR page without changing the renderer if it uses existing components. For example, a runtime route could return a `ContextDiagramPanel` node with a session-derived `snapshot` and style set.

### Goja DSL contract

`pkg/widgetdsl/module.go` exposes helpers for server-side JavaScript authors. The style helpers are already installed:

- `visualStyle(options)` (`module.go:273-283`),
- `legendItem(id, label, options?)` (`module.go:284-288`),
- `styleSet(options)` (`module.go:289-298`),
- `contextPart(id, label, styleKey, tokens, options?)` (`module.go:299-303`),
- `contextSnapshot(options)` (`module.go:304-319`),
- `paletteStyleSet(options)` (`module.go:320-322`).

The DSL recipe `contextDiagram` requires a `styleSet`, or derives one from `palette + entries`; it then emits a `ContextDiagramPanel` (`module.go:608-629`). The xgoja provider registers the `context_window.dsl` module as the context-window, transcript, annotation, and anchored-comment helper package (`provider.go:36-42`).

For an uploaded session, the runtime author can either:

1. build `snapshot` and `styleSet` in Go, then return Widget IR directly; or
2. run Goja/JS code that uses `context_window.dsl` helpers to build the same JSON.

The implementation should choose one path per host application, not both in the same first pass.

### Upload UI and current examples

`ContextUploadDropArea` is currently a presentational organism around `FileDropZone`; it accepts `.json` by default and emits selected files to a callback (`ContextUploadDropArea.tsx:4-35`). It does not parse the session itself.

Existing Storybook examples already prove that context parts can correspond to transcript/tool blocks. `WidgetRenderer.context-diagrams.stories.tsx` defines `contentSnapshot` with `system`, a user turn, a tool call, a tool result, an assistant answer, and free space (`WidgetRenderer.context-diagrams.stories.tsx:23-37`). The same file shows caller-defined vocabularies and hidden headroom entries (`WidgetRenderer.context-diagrams.stories.tsx:80-187`).

This means the requested visualization is an extension of an existing pattern, not a greenfield feature.

## Proposed User Experience

### Whole-window overview

The default uploaded-session visualization should show a single current context-window snapshot.

Recommended layout:

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ Uploaded session context window                         [strip][stack][map] │
│ 184,230 / 200,000 tokens · current window · 37 turns · 412 blocks            │
├──────────────────────────────────────────────────────────────────────────────┤
│ Legend: ■ System  ■ User  ▧ Agent  ▥ Tool call  ▦ Tool result  ▤ File read   │
│         ▣ File write  ▨ Thinking  □ Summary  ··· Headroom                    │
├──────────────────────────────────────────────────────────────────────────────┤
│ STRIP VIEW                                                                    │
│ ┌──────┬─┬───┬───────┬────────────────────┬──┬──────────────┬─────────────┐ │
│ │system│u│agt│tool:rg│tool result:rg output│rw│file reads     │thinking     │ │
│ └──────┴─┴───┴───────┴────────────────────┴──┴──────────────┴─────────────┘ │
│                 T01        T02                         T03        T04        │
├──────────────────────────────────────────────────────────────────────────────┤
│ Selected block: tool result · rg output · 9,840 tok · turn 03                │
│ packages/.../ContextDiagramPanel.tsx: matching component references...       │
└──────────────────────────────────────────────────────────────────────────────┘
```

The strip view communicates ordering well. It should be the default for session replay because a context window is ordered text. Treemap is useful for "what is biggest?" analysis but does not preserve ordering. Stack is useful for readable labels and details.

### Token-proportional rendering with readability limits

The current components already use readability floors:

- Strip: minimum 2% width per segment, labels only when width >= 7%.
- Stack: minimum 26px and maximum 74px row height.
- Treemap: minimum 7% tile basis.

For the uploaded session, these floors are helpful but can distort extremely many tiny blocks. The implementation should document this explicitly in the caption:

```text
Widths are token-proportional with a 2% minimum so small blocks remain selectable.
```

If more than about 80-120 blocks are in one snapshot, the future implementation should aggregate small adjacent blocks of the same style into summary blocks by default, while preserving raw blocks in metadata or detail view.

### Turn-aware overview

If turn IDs are available, show turn separators. The current `ContextWindowPart.metadata` can carry `turnIndex`, `turnId`, and `turnRole`, but the current diagram components do not render separators or group headers.

A feasible first-pass approximation is to encode turn names in block labels:

```text
T01 system | T01 user | T01 agent | T01 tool:read | T01 result | T02 user | ...
```

A better design-system update would add grouped strip support:

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ Turn timeline                                                                 │
│ T01        T02                         T03                         T04        │
│ ┌────┬────┬──────┐ ┌────┬─────┬───────┐ ┌────┬────┬──────────────┐ ┌───────┐ │
│ │sys │user│agent │ │user│tool │result │ │user│file│file contents │ │answer │ │
│ └────┴────┴──────┘ └────┴─────┴───────┘ └────┴────┴──────────────┘ └───────┘ │
└──────────────────────────────────────────────────────────────────────────────┘
```

This requires either:

- a new `ContextGroupedStripDiagram` component, or
- extending `ContextStripDiagram` with optional `groupByMetadataKey="turnIndex"`.

Do not add this in the first implementation unless the uploaded data has reliable turn metadata and the simple label/metadata approach proves insufficient.

### Deferred: page-through turn view

This section is retained only as a future design note. The current implementation scope should not expose a live/page-through session view. If a later ticket reopens this scope and the session has reliable per-turn context-window snapshots, a page-through version could let the user navigate between turns:

```text
┌──────────────────────────────────────────────────────────────────────────────┐
│ Context window by turn                                ◀ Turn 14 / 37 ▶       │
│ After assistant tool result · 132,640 / 200,000 tokens                       │
├──────────────────────────────────────────────────────────────────────────────┤
│ ┌──────┬────┬──────┬───────────┬───────────────────────────────┬──────────┐ │
│ │system│hist│user14│tool call  │tool result: file read output   │thinking  │ │
│ └──────┴────┴──────┴───────────┴───────────────────────────────┴──────────┘ │
├──────────────────────────────────────────────────────────────────────────────┤
│ Newly added this turn:                                                       │
│ ▥ tool call: read packages/.../fixtures.ts          180 tok                  │
│ ▦ tool result: file content                         7,420 tok                │
│ ▨ thinking / plan                                   1,100 tok                │
└──────────────────────────────────────────────────────────────────────────────┘
```

This can be a useful analysis UX, but it needs more data and is out of current scope. The runtime must know which blocks were present after each turn. If the uploaded file only has final transcript messages, the application can derive "blocks by turn" but not exact context-window membership after compaction/eviction.

## go-minitrace Data Contract for Uploaded Sessions

The uploaded-session flow is not starting from arbitrary JSON once the app uses go-minitrace. The expected pipeline is:

```text
raw uploaded agent transcript/session JSON
  -> go-minitrace importer / converter
  -> .minitrace.json session or temporary minitrace DB
  -> minitracejs DSL views / SQL recipes
  -> Widget IR props for ContextDiagramPanel and TranscriptWorkspacePanel
```

This matters because go-minitrace already provides a normalized schema and query DSL. The widget layer should consume those normalized rows rather than inventing a separate ad-hoc parser for every raw source format.

### Authoritative minitrace session fields

The authoritative schema docs live in `/home/manuel/workspaces/2026-06-07/club-meetup-site/go-minitrace/pkg/doc/minitrace-schema.md`, with Go structs in `/home/manuel/workspaces/2026-06-07/club-meetup-site/go-minitrace/pkg/minitrace/schema.go`. Important fields for widget population are:

- **Session metadata**: `id`, `schema_version`, `profile`, `quality`, `title`, `summary`, `classification`.
- **Provenance**: `provenance.source_format`, `source_path`, `converted_at`, `converter_version`, `original_session_id`.
- **Environment**: `environment.model`, `tools_enabled`, `system_prompt`, `agent_framework`, `platform_type`, `provider_hint`.
- **Operational context**: `operational_context.working_directory`, `git_branch`, `git_ref`, `autonomy_level`, `sandbox`, `framework_config`.
- **Timing**: `started_at`, `ended_at`, `duration_seconds`, `active_duration_seconds`, `hour_of_day`, `day_of_week`.
- **Turns**: ordered `turns[]` with `index`, `timestamp`, `role`, `source`, `model`, `content_type`, `content`, `framework_metadata`, `tool_calls_in_turn`, `thinking`, `streaming`, and `usage`.
- **Per-turn usage**: `input_tokens`, `output_tokens`, `cache_read_tokens`, `cache_creation_tokens`, `reasoning_tokens`, and `tool_tokens`.
- **Tool calls**: ordered `tool_calls[]` with `id`, `emitting_turn_index`, `timestamp`, `tool_name`, normalized `operation_type`, shared `input.file_path`, `input.command`, `input.justification`, raw `input.arguments`, `output.result`, `output.error`, `output.full_bytes`, `output.truncated`, and `output.content_origin`.
- **Annotations**: session/turn/tool annotations with target metadata; these can become transcript notes or side rail cards.
- **Metrics**: summary counts and token totals suitable for headers and compact metadata.

The key correction to the original design is that block classification can use `tool_calls[].operation_type` and the minitrace DSL rows. It does not need to infer file reads/writes only from string matching.

### minitracejs importer flow

The JavaScript API reference in `/home/manuel/workspaces/2026-06-07/club-meetup-site/go-minitrace/pkg/doc/js-api-reference.md` describes the upload/import path:

```js
const importer = mt.importer()
  .Content(uploadText)
  .Name(filename)
  .AutoDetect()
  .Strict()
  .Convert();

const preview = importer.Preview();
const saved = importer.Into(sessionsDir).SessionID(sessionId).Save();
```

Useful preview fields include diagnostics, role/tool counts, sampled turns/tools, `hasSystemPrompt`, thinking signals, image signals, and subagent count. The upload widget can use the preview before saving/rendering to show whether the file is supported and whether sensitive/full-content previews are being used.

### minitracejs view/query outputs

The DSL implementation in `/home/manuel/workspaces/2026-06-07/club-meetup-site/go-minitrace/pkg/minitracejs/query_view_session.go` exposes these recipes/views:

| Builder/view | Output shape | Widget use |
| --- | --- | --- |
| `mt.query().SessionSummary().SessionID(id).Build()` | one `SessionSummary` row | Header metadata and selected-session summary. |
| `mt.query().TurnRows().SessionID(id).Build()` | ordered turn rows | Transcript messages and turn-level context blocks. |
| `mt.query().ToolRows().SessionID(id).Build()` | ordered tool rows | Tool/file context blocks and detail metadata. |
| `mt.query().EventRows()` / `TurnBlockRows()` | renderable event/block rows | Fine-grained transcript frame rendering if needed. |
| `mt.query().TranscriptRows().IncludeTools().SessionID(id).Build()` | flattened transcript rows | Best initial source for `TranscriptWorkspacePanel` and a block strip. |
| `mt.query().TimelineRows().SessionID(id).Build()` | turn-level timeline rows | Turn overview, warning badges, and optional grouping labels. |
| `mt.query().TokenUsageRows().ByTurn()` | token usage grouped by turn | Per-turn budget summary and proportional context blocks. |
| `mt.query().TokenUsageRows().ByRole()` | token usage grouped by role | Summary chart/legend validation. |
| `mt.query().TokenUsageRows().ByTool()` | estimated tool output token rows | Tool/file block sizing when exact token usage is absent. |
| `session.view().Transcript().IncludeTools().Run()` | plain transcript rows | Direct app flow for transcript widgets. |
| `session.view().Timeline().Run()` | timeline rows | Turn navigation metadata without adding a live view. |
| `session.view().TokenUsage().ByTurn().Run()` | token usage rows | Context budget summaries. |
| `session.view().TurnFrames().Run()` | frames with `blocks`, `toolCalls`, `stats` | Future grouped-turn details if needed. |

`TranscriptRows` is especially useful because it already flattens messages, thinking, and optional tools into rows with fields that map almost directly to widgets:

```text
id, session_id, turn_index, ordinal, role, kind, name, title, text,
timestamp, tokens, tool_call_id, severity, collapsed_by_default, metadata
```

The implementation currently labels row kinds such as `message`, `thinking`, `tool_result`, and `tool_error`. Tool rows add stronger operation metadata through `operation_type`, `file_path`, `command`, `success`, `result`, `error`, `truncated`, and `full_bytes`.

### Mapping minitrace rows to transcript widgets

`TranscriptRows` can populate transcript widgets without a bespoke parser:

```ts
function transcriptRowToMessage(row: MinitraceTranscriptRow): TranscriptMessage {
  return {
    id: row.id,
    role: row.role === 'tool' ? 'tool' : row.role,
    name: row.name || row.kind,
    text: row.text ?? '',
    tokens: row.tokens ?? undefined,
    timestamp: row.timestamp ?? undefined,
    metadata: {
      sessionId: row.session_id,
      turnIndex: row.turn_index,
      ordinal: row.ordinal,
      kind: row.kind,
      toolCallId: row.tool_call_id,
      severity: row.severity,
      collapsedByDefault: Boolean(row.collapsed_by_default),
      raw: row.metadata,
    },
  };
}
```

Recommended transcript style mapping:

| minitrace row | Widget role/style |
| --- | --- |
| `role = system`, `kind = message` | `role: 'system'`, `styleKey: 'system'` |
| `role = user`, `kind = message` | `role: 'user'`, `styleKey: 'user'` |
| `role = assistant`, `kind = message` | `role: 'assistant'`, `styleKey: 'assistant'` or `agent` |
| `kind = thinking` | `role: 'assistant'`, `styleKey: 'thinking'`, collapsed by default |
| `kind = tool_result` | `role: 'tool'`, `styleKey: 'result'` or `tool_result` |
| `kind = tool_error` | `role: 'tool'`, `styleKey: 'tool_error'`, danger/error metadata |

### Mapping minitrace rows to context-window blocks

For the context block view, prefer a two-source strategy:

1. Use `TranscriptRows` for ordered visible rows: messages, thinking, and tool result/error text.
2. Join or correlate `ToolRows` when the block needs operation-aware classification: file read/write, shell execute, delegate, or unknown tool operation.

Pseudo-code:

```ts
function transcriptRowToContextPart(row: MinitraceTranscriptRow, toolById: Map<string, MinitraceToolRow>): ContextWindowPart {
  const tool = row.tool_call_id ? toolById.get(row.tool_call_id) : undefined;
  const styleKey = row.kind === 'message'
    ? roleToStyleKey(row.role)
    : row.kind === 'thinking'
      ? 'thinking'
      : tool
        ? toolRowToStyleKey(tool, row.kind)
        : row.kind === 'tool_error'
          ? 'tool_error'
          : 'tool_result';

  return {
    id: row.id,
    label: compactBlockLabel(row, tool),
    styleKey,
    tokens: Math.max(1, Number(row.tokens ?? estimateTokens(row.text))),
    contentPreview: preview(row.text),
    content: row.text,
    sourceId: row.session_id,
    metadata: {
      turnIndex: row.turn_index,
      ordinal: row.ordinal,
      kind: row.kind,
      role: row.role,
      toolCallId: row.tool_call_id,
      toolName: tool?.tool_name,
      operationType: tool?.operation_type,
      filePath: tool?.file_path,
      command: tool?.command,
      success: tool?.success,
      truncated: tool?.truncated,
      fullBytes: tool?.full_bytes,
    },
  };
}

function roleToStyleKey(role: string): string {
  if (role === 'system') return 'system';
  if (role === 'user') return 'user';
  if (role === 'assistant') return 'agent';
  if (role === 'tool') return 'tool_result';
  return 'other';
}

function toolRowToStyleKey(tool: MinitraceToolRow, rowKind: string): string {
  if (rowKind === 'tool_error' || tool.success === false) return 'tool_error';
  switch (tool.operation_type) {
    case 'READ': return 'file_read';
    case 'MODIFY': return 'file_write';
    case 'NEW': return 'file_write';
    case 'EXECUTE': return 'tool_call';
    case 'DELEGATE': return 'agent';
    default: return 'tool_result';
  }
}
```

Recommended context style keys:

| styleKey | Source condition | Label |
| --- | --- | --- |
| `system` | turn/row role `system` | System |
| `user` | turn/row role `user` | User |
| `agent` | assistant message | Agent |
| `thinking` | row kind `thinking` or reasoning text | Thinking |
| `tool_call` | tool `operation_type = EXECUTE` where command itself is the meaningful block | Tool call |
| `tool_result` | generic successful tool output | Tool result |
| `tool_error` | failed tool output or row kind `tool_error` | Tool error |
| `file_read` | `operation_type = READ` | File read |
| `file_write` | `operation_type in (MODIFY, NEW)` | File write |
| `delegate` | `operation_type = DELEGATE` | Delegated agent |
| `summary` | source marks summary/compaction rows | Summary |
| `headroom` | computed model context headroom when known | Headroom |
| `other` | unknown rows/tools | Other |

### What remains unclear in go-minitrace docs/API

The current docs are strong enough for transcript and tool widgets, but the following areas should be expanded if implementation exposes them:

- Whether `TranscriptRows` should include explicit tool-call request rows in addition to tool result rows. Current SQL emits tool result/error rows from `tool_calls`; the command/input can appear in `text` only as fallback when result/error is empty.
- Whether exact context-window membership is available for any adapter. The minitrace schema provides turns, usage, and tools, but not necessarily the exact prompt assembled for a model call after compaction or hidden system injections.
- Whether token usage is adapter-provided or estimated for each source. Turns can have usage fields, while tool rows estimate tokens from `full_bytes`/result length.
- How thinking/reasoning should be labeled for each framework. The schema exposes `turns[].thinking` if captured, but UI language should avoid implying hidden chain-of-thought reconstruction.
- Whether `operation_type` values are complete enough for all file read/write detection, or whether some adapters require fallback inspection of `tool_name` and `input.arguments`.

These should become go-minitrace documentation improvements if they are encountered during implementation.

## Proposed Data Model for Uploaded Session Normalization

The implementation should introduce a small app-side normalized representation before creating Widget IR, but that representation should be derived from go-minitrace/minitracejs rows rather than raw uploaded JSON. This can live in the ClubMeetup/minitrace-viz side first, then move upstream only if multiple applications need the same mapping.

### Block type vocabulary

Recommended initial style keys:

| styleKey | Legend label | Intended source records | Notes |
| --- | --- | --- | --- |
| `system` | System | System prompt, tool policy, developer instruction | Use checker/blue-ish style. |
| `user` | User | User messages | Use solid/accent style for quick scanning. |
| `agent` | Agent | Assistant visible responses | Use diagonal/cool accent. |
| `tool_call` | Tool call | Tool invocation request | Should be distinct from result. |
| `tool_result` | Tool result | Command output, API result | Often large; use high-contrast token labels. |
| `file_read` | File read | Read/file content returned to model | Often the largest tenant. |
| `file_write` | File write | Writes, patches, edits | Usually smaller but important. |
| `thinking` | Thinking | Reasoning/scratchpad/generated plan if present | Only if source explicitly exposes it; do not infer private chain-of-thought. |
| `summary` | Summary | Compacted previous turns | Useful when context compression is represented. |
| `evicted` | Evicted | Dropped/compressed content markers | Optional. |
| `headroom` | Headroom | Free context budget | Hidden legend entry if not analytically useful. |
| `other` | Other | Unknown blocks | Should be rare; emit metadata reason. |

Important privacy/safety note: "thinking" should mean exposed/generated scratchpad or trace metadata that the uploaded session already includes. Do not attempt to reconstruct hidden model reasoning.

### Normalized records

Pseudo-TypeScript:

```ts
type SessionBlockType =
  | 'system'
  | 'user'
  | 'agent'
  | 'tool_call'
  | 'tool_result'
  | 'file_read'
  | 'file_write'
  | 'thinking'
  | 'summary'
  | 'evicted'
  | 'headroom'
  | 'other';

interface NormalizedSessionBlock {
  id: string;
  type: SessionBlockType;
  label: string;
  tokens: number;
  textPreview?: string;
  text?: string;
  turnIndex?: number;
  turnId?: string;
  sourceKind?: string;
  toolName?: string;
  filePath?: string;
  startToken?: number;
  endToken?: number;
  metadata?: Record<string, unknown>;
}

interface NormalizedSessionTurn {
  id: string;
  index: number;
  title: string;
  blockIds: string[];
  addedTokens: number;
  cumulativeTokens?: number;
}

interface NormalizedUploadedSession {
  id: string;
  title: string;
  tokenLimit: number;
  blocks: NormalizedSessionBlock[];
  turns: NormalizedSessionTurn[];
  hasReliableTurnBoundaries: boolean;
  hasPerTurnWindowSnapshots: boolean;
  warnings: string[];
}
```

### Mapping to existing `ContextWindowSnapshot`

Pseudo-code:

```ts
function blockToContextPart(block: NormalizedSessionBlock): ContextWindowPart {
  return {
    id: block.id,
    label: block.label,
    styleKey: block.type,
    tokens: block.tokens,
    contentPreview: block.textPreview,
    content: block.text,
    startToken: block.startToken,
    endToken: block.endToken,
    metadata: {
      turnIndex: block.turnIndex,
      turnId: block.turnId,
      sourceKind: block.sourceKind,
      toolName: block.toolName,
      filePath: block.filePath,
      ...block.metadata,
    },
  };
}

function sessionToSnapshot(session: NormalizedUploadedSession): ContextWindowSnapshot {
  const usedTokens = sum(session.blocks.filter(b => b.type !== 'headroom').map(b => b.tokens));
  const headroom = Math.max(0, session.tokenLimit - usedTokens);

  return {
    id: `${session.id}-current-window`,
    title: `${session.title} context window`,
    subtitle: `${usedTokens.toLocaleString()} / ${session.tokenLimit.toLocaleString()} tokens`,
    limit: session.tokenLimit,
    parts: [
      ...session.blocks.map(blockToContextPart),
      ...(headroom > 0 ? [{ id: 'headroom', label: 'headroom', styleKey: 'headroom', tokens: headroom }] : []),
    ],
    metadata: {
      sourceSessionId: session.id,
      turnCount: session.turns.length,
      hasReliableTurnBoundaries: session.hasReliableTurnBoundaries,
      warnings: session.warnings,
    },
  };
}
```

### Mapping to turn snapshots

Only build these if the source can support them.

```ts
function sessionToTurnSnapshots(session: NormalizedUploadedSession): ContextWindowSnapshot[] {
  if (!session.hasReliableTurnBoundaries) return [];

  return session.turns.map(turn => {
    const blocksThroughTurn = session.blocks.filter(block =>
      block.turnIndex == null || block.turnIndex <= turn.index
    );

    return {
      id: `${session.id}-turn-${turn.index}`,
      title: `Turn ${turn.index}: ${turn.title}`,
      subtitle: `${turn.addedTokens.toLocaleString()} tokens added this turn`,
      limit: session.tokenLimit,
      parts: blocksThroughTurn.map(blockToContextPart),
      selectedPartId: turn.blockIds[0],
      metadata: { turnIndex: turn.index, turnId: turn.id },
    };
  });
}
```

If actual context-window compaction/eviction happens between turns, the above approximation is wrong. In that case, the source format must provide explicit per-turn membership or snapshot boundaries. This is a runtime/schema requirement, not a React styling problem.

## Proposed Style Set

The future implementation should build a session block style set with the palette helper. It should hide headroom by default if the UX focuses on used tokens.

Pseudo-TypeScript:

```ts
const sessionBlockStyleSet = createContextStyleSetFromPalette({
  id: 'uploaded-session-block-types',
  name: 'Uploaded session block types',
  palette: signalOrangeCyan,
  legendSize: 'sm',
  swatchSize: 'sm',
  entries: [
    { id: 'system', label: 'System', accent: 'b', pattern: 'checker', fillPct: 18, linePct: 70 },
    { id: 'user', label: 'User', accent: 'a', pattern: 'solid', solid: true },
    { id: 'agent', label: 'Agent', accent: 'b', pattern: 'diagonalDense', fillPct: 16, linePct: 65 },
    { id: 'tool_call', label: 'Tool call', accent: 'a', pattern: 'stipple', fillPct: 14, linePct: 55 },
    { id: 'tool_result', label: 'Tool result', accent: 'a', pattern: 'cross', fillPct: 16, linePct: 70 },
    { id: 'file_read', label: 'File read', accent: 'grid', pattern: 'diagonal', fillPct: 12, linePct: 45 },
    { id: 'file_write', label: 'File write', accent: 'c', pattern: 'checker', fillPct: 24, linePct: 70 },
    { id: 'thinking', label: 'Thinking', accent: 'c', pattern: 'diagonalDense', fillPct: 18, linePct: 70 },
    { id: 'summary', label: 'Summary', accent: 'b', pattern: 'diagonal', fillPct: 14, linePct: 55 },
    { id: 'evicted', label: 'Evicted', accent: 'shadow', pattern: 'cross', fillPct: 10, linePct: 45 },
    { id: 'headroom', label: 'Headroom', accent: 'grid', pattern: 'none', hidden: true },
    { id: 'other', label: 'Other', accent: 'ink', pattern: 'overflow', fillPct: 10, linePct: 50 },
  ],
});
```

Equivalent Goja DSL sketch:

```js
const contextWindow = require("context_window.dsl")

const styleSet = contextWindow.paletteStyleSet({
  id: "uploaded-session-block-types",
  name: "Uploaded session block types",
  palette: "Signal Orange / Cyan",
  legendSize: "sm",
  swatchSize: "sm",
  entries: [
    { id: "system", label: "System", accent: "b", pattern: "checker", fillPct: 18, linePct: 70 },
    { id: "user", label: "User", accent: "a", pattern: "solid", solid: true },
    { id: "agent", label: "Agent", accent: "b", pattern: "diagonalDense", fillPct: 16, linePct: 65 },
    { id: "tool_call", label: "Tool call", accent: "a", pattern: "stipple", fillPct: 14, linePct: 55 },
    { id: "tool_result", label: "Tool result", accent: "a", pattern: "cross", fillPct: 16, linePct: 70 },
    { id: "file_read", label: "File read", accent: "grid", pattern: "diagonal", fillPct: 12, linePct: 45 },
    { id: "file_write", label: "File write", accent: "c", pattern: "checker", fillPct: 24, linePct: 70 },
    { id: "thinking", label: "Thinking", accent: "c", pattern: "diagonalDense", fillPct: 18, linePct: 70 },
    { id: "headroom", label: "Headroom", accent: "grid", pattern: "none", hidden: true },
  ],
})
```

## Widget Composition Options

### Option A: Existing `ContextDiagramPanel` only

Widget IR sketch:

```js
contextWindow.recipes.contextDiagram({
  snapshot: sessionToSnapshot(uploadedSession),
  styleSet,
  initialView: "strip",
  selectedPartId: "turn-14-tool-result-2",
})
```

ASCII rendering:

```text
┌ Uploaded session context window ─────────────── [strip][stack][budget][map] ┐
│ 184,230 / 200,000 tokens · generated from upload                            │
│                                                                              │
│ ┌───────┬──┬─────┬───┬────────────────────┬──────┬───────────────┬───────┐ │
│ │system │u │agent│call│tool result         │write │file read      │think  │ │
│ └───────┴──┴─────┴───┴────────────────────┴──────┴───────────────┴───────┘ │
│                                                                              │
│ ■ System ■ User ▧ Agent ▥ Tool call ▦ Tool result ▤ File read ▣ File write   │
│                                                                              │
│ Selected: tool result · bash rg output · 9,840 tok · turn 14                 │
│ packages/.../fixtures.ts: ...                                                │
└──────────────────────────────────────────────────────────────────────────────┘
```

Pros:

- fastest implementation,
- uses existing types, components, styles, and Widget IR,
- works with current Storybook and package docs,
- detail panel already shows metadata and preview.

Cons:

- no explicit turn separators,
- no page-through navigation,
- large sessions with hundreds of blocks may become visually crowded.

### Option B: Existing panel plus metadata sidebar

Use a `SplitPane`: left side is `ContextDiagramPanel`, right side is metadata, warnings, and selected-turn details. This pattern already exists in `WidgetRenderer.context-diagrams.stories.tsx:233-253`.

ASCII rendering:

```text
┌─────────────────────────────────────┬──────────────────────────────────────┐
│ Context blocks                       │ Session metadata                     │
│ [strip][stack][budget][map]          │ Turns: 37                            │
│ ┌───┬──┬──────┬──────────────┐       │ Blocks: 412                          │
│ │sys│u │agent │tool result   │       │ Token limit: 200,000                 │
│ └───┴──┴──────┴──────────────┘       │ Warnings: inferred file reads        │
│                                      │                                      │
│ Selected: file read                  │ Turn 14 blocks                       │
│ packages/rag.../styles.ts            │ - tool call: read                    │
│ 7,420 tokens                         │ - file read: styles.ts               │
└─────────────────────────────────────┴──────────────────────────────────────┘
```

Pros:

- still no new diagram component,
- better for explaining inferred boundaries and warnings,
- gives an intern a safe first production target.

Cons:

- still not a true turn-paging UI.

### Deferred Option C: New `ContextTurnPagerPanel`

A new organism receives multiple snapshots and renders navigation around `ContextDiagramPanel`.

Proposed props:

```ts
interface ContextTurnPagerPanelProps extends HTMLAttributes<HTMLDivElement> {
  snapshots: ContextWindowSnapshot[];
  styleSet: ContextStyleSet;
  initialIndex?: number;
  showDeltaList?: boolean;
  showLegend?: boolean;
}
```

Widget IR type addition:

```ts
export interface ContextTurnPagerPanelWidgetProps extends BaseWidgetProps {
  snapshots: ContextWindowSnapshot[];
  styleSet: ContextStyleSet;
  initialIndex?: number;
  showDeltaList?: boolean;
  showLegend?: boolean;
}
```

ASCII rendering:

```text
┌ Context window by turn ──────────────────────────────── ◀ Turn 14 / 37 ▶ ┐
│ After tool result · selected: file read result · 132,640 / 200,000 tok    │
│                                                                           │
│ ┌──────┬──────────┬─────┬────────────────────────┬──────────────┐        │
│ │system│history   │user │tool call               │file read     │        │
│ └──────┴──────────┴─────┴────────────────────────┴──────────────┘        │
│                                                                           │
│ Added this turn                                                           │
│ ▥ tool call: read packages/rag-evaluation-site/src/context/types.ts       │
│ ▤ file read: 170 lines · 5,920 tok                                         │
└───────────────────────────────────────────────────────────────────────────┘
```

Pros:

- directly satisfies "page through individual turns",
- keeps the diagram implementation reused,
- can be introduced as a small organism without rewriting lower-level diagrams.

Cons:

- requires reliable turn snapshots or a clearly documented approximation,
- requires Widget IR type/registry/story/DSL updates,
- adds stateful UI and keyboard/accessibility concerns.

### Option D: New grouped strip component

A grouped strip renders turn boundaries within one snapshot. It should be considered if whole-window ordering matters and paging is too slow.

Potential data extension:

```ts
interface ContextWindowPart {
  // existing fields unchanged
  metadata?: {
    turnIndex?: number;
    turnId?: string;
    groupLabel?: string;
    blockType?: string;
    [key: string]: ContextJsonValue;
  };
}
```

Potential component prop:

```ts
interface ContextGroupedStripDiagramProps {
  snapshot: ContextWindowSnapshot;
  styleSet: ContextStyleSet;
  groupByMetadataKey?: string; // default: "turnIndex"
  groupLabelMetadataKey?: string; // default: "turnLabel"
}
```

ASCII rendering:

```text
T01                  T02                                  T03
┌──────┬────┬──────┐ ┌────┬──────────┬──────────────────┐ ┌────┬───────┐
│system│user│agent │ │user│tool call │tool result       │ │user│agent  │
└──────┴────┴──────┘ └────┴──────────┴──────────────────┘ └────┴───────┘
```

Pros:

- best single-screen expression of turn boundaries,
- useful for uploaded sessions with many compact turns.

Cons:

- more diagram-specific work,
- grouping labels can collide in dense views,
- current `ContextWindowPart` metadata can carry grouping, but the renderer does not yet use it.

## Recommended Implementation Plan

### Phase 0: Confirm uploaded session schema

Goal: determine whether turn boundaries and token counts are reliable.

Tasks:

1. Collect 1-3 real uploaded session JSON examples.
2. Identify where roles, tool calls, tool results, file operations, token counts, and turn boundaries live.
3. Decide whether token counts are source-provided or must be estimated.
4. Decide whether the source represents final context-window membership, per-turn membership, or only a transcript.

Output:

- schema notes,
- examples in `ttmp/.../sources/`,
- a decision on whether Phase 1 can be exact or approximate.

### Phase 1: Build a normalizer with tests

Goal: convert uploaded session data into `NormalizedUploadedSession`.

Pseudo-code:

```ts
function normalizeUploadedSession(raw: unknown): NormalizedUploadedSession {
  const records = detectSchema(raw);
  const blocks = [];
  const turns = new Map();

  for (const record of records) {
    const block = classifyRecord(record);
    block.tokens = readOrEstimateTokens(record);
    block.label = makeCompactLabel(block);
    block.textPreview = makePreview(record.text, 300);
    blocks.push(block);
    addBlockToTurn(turns, block);
  }

  return {
    id: stableSessionId(raw),
    title: titleFromRaw(raw),
    tokenLimit: tokenLimitFromRaw(raw) ?? 200_000,
    blocks: coalesceTinyAdjacentBlocks(blocks),
    turns: Array.from(turns.values()),
    hasReliableTurnBoundaries: detectedTurnIdsAreReliable(records),
    hasPerTurnWindowSnapshots: detectedSnapshotBoundaries(records),
    warnings: collectWarnings(records),
  };
}
```

Classification heuristic sketch:

```ts
function classifyRecord(record): SessionBlockType {
  if (record.role === 'system' || record.kind === 'system_prompt') return 'system';
  if (record.role === 'user') return 'user';
  if (record.role === 'assistant') return record.isThinking ? 'thinking' : 'agent';
  if (record.kind === 'tool_call') return classifyToolCall(record);
  if (record.kind === 'tool_result') return classifyToolResult(record);
  if (record.kind === 'summary') return 'summary';
  if (record.kind === 'evicted') return 'evicted';
  return 'other';
}

function classifyToolCall(record): SessionBlockType {
  if (record.toolName?.match(/write|edit|patch/i)) return 'file_write';
  if (record.toolName?.match(/read|cat|sed|grep|rg/i)) return 'tool_call';
  return 'tool_call';
}

function classifyToolResult(record): SessionBlockType {
  if (record.metadata?.filePath || record.toolName?.match(/read|cat|sed/i)) return 'file_read';
  if (record.toolName?.match(/write|edit|patch/i)) return 'file_write';
  return 'tool_result';
}
```

Token estimation should be explicit and replaceable:

```ts
function readOrEstimateTokens(record): number {
  if (Number.isFinite(record.tokens)) return record.tokens;
  if (Number.isFinite(record.token_count)) return record.token_count;
  return Math.max(1, Math.ceil(String(record.text ?? '').length / 4));
}
```

### Phase 2: Emit existing Widget IR

Goal: produce a usable visualization without new React components.

Widget IR sketch:

```ts
const node = component('SplitPane', {
  ratio: 'rightNarrow',
  divider: true,
  left: component('ContextDiagramPanel', {
    snapshot: sessionToSnapshot(normalized),
    styleSet: sessionBlockStyleSet,
    initialView: 'strip',
    views: ['strip', 'stack', 'budget', 'treemap'],
    showPartDetails: true,
  }),
  right: component('Stack', { gap: 'md' }, [
    component('Panel', { title: 'Upload analysis' }, [metadataGridNode(normalized)]),
    component('Panel', { title: 'Warnings' }, [warningsNode(normalized.warnings)]),
  ]),
});
```

Acceptance criteria:

- all blocks use caller-defined `styleKey`,
- legend labels match the block vocabulary,
- block sizes are proportional to token sizes within current component readability limits,
- selected block detail shows preview and metadata,
- headroom/free blocks do not inflate used-token totals,
- warnings are shown when turns or token counts are inferred.

### Phase 3: Add Storybook coverage

Goal: make the design reviewable before integrating into the upload runtime.

Stories to add:

1. `UploadedSessionBlockMap` — a fixture with system/user/agent/tool/file/thinking blocks.
2. `UploadedSessionDenseBlocks` — many small tool results to validate readability and aggregation.
3. `UploadedSessionUnknownSchemaWarnings` — shows warnings and `other` blocks.
4. Optional `UploadedSessionTurnPagerPrototype` — only if a pager component is implemented.

Story fixtures should include intentionally varied token sizes:

```ts
parts: [
  { id: 'sys', label: 'system', styleKey: 'system', tokens: 7200 },
  { id: 'u12', label: 'T12 user', styleKey: 'user', tokens: 260 },
  { id: 'a12-plan', label: 'T12 plan', styleKey: 'thinking', tokens: 1800 },
  { id: 'tc12-rg', label: 'T12 rg', styleKey: 'tool_call', tokens: 140 },
  { id: 'tr12-rg', label: 'rg output', styleKey: 'tool_result', tokens: 9400 },
  { id: 'fr12-fixtures', label: 'fixtures.ts', styleKey: 'file_read', tokens: 6200 },
  { id: 'fw12-patch', label: 'patch', styleKey: 'file_write', tokens: 900 },
]
```

### Phase 4: Decide whether to add a design-system component

Do not add `ContextTurnPagerPanel` in the current upload/visualize/transcript scope. Reconsider it only if a later ticket explicitly requests page-through analysis and all of the following are true:

- uploaded session data has reliable turn boundaries; and
- stakeholders prefer page-through analysis over a single full-window view; and
- Storybook shows that the existing panel plus metadata sidebar is insufficient.

Add `ContextGroupedStripDiagram` if and only if:

- a single-screen whole-window view must show turn boundaries; and
- labels/separators can remain readable at the target block counts; and
- it can preserve keyboard navigation and selected-part details.

### Phase 5: Runtime integration

The runtime integration depends on where uploads are handled. `ContextUploadDropArea` is currently presentational, so the host app or Widget IR page endpoint should own parsing.

Possible architecture:

```text
Browser upload
  └─ onFilesSelected / POST /api/v1/session-uploads
       └─ backend stores raw upload
       └─ backend normalizes blocks
       └─ backend returns Widget IR page id
            └─ WidgetRenderer fetches/render page
                 └─ ContextDiagramPanel + styleSet + snapshot
```

Or, for a local-only frontend prototype:

```text
Browser upload
  └─ FileReader reads JSON
       └─ normalizeUploadedSession(raw)
       └─ sessionToSnapshot(normalized)
       └─ render ContextDiagramPanel directly
```

Production should prefer backend normalization if the uploaded session can be large, because tokenization and schema detection may become heavy and should be testable outside React.

## Gaps and Design-System Requirements

### Gap: No turn separator rendering inside diagrams

Existing diagrams render parts but do not draw group separators. `ContextWindowPart.metadata` can carry `turnIndex`, but no component reads it for grouping.

Recommendation: start without separators. If needed, add `ContextGroupedStripDiagram` or a grouping option to `ContextStripDiagram`.

### Gap: No multi-snapshot pager

`ContextDiagramPanel` renders one snapshot. It has internal state for selected view and selected part, but not for multiple snapshots.

Recommendation: do not address this gap in the current upload/visualize/transcript scope. If per-turn snapshots become a future approved requirement, add a small wrapper organism (`ContextTurnPagerPanel`) rather than modifying every diagram component.

### Gap: Uploaded session parser is undefined

The frontend has an upload drop area, but no observed parser for real uploaded sessions. This means the exact mapping from upload records to block types is unknown.

Recommendation: Phase 0 must inspect actual upload samples before writing production parsing logic.

### Gap: Token counts may be absent

The source may not provide exact token counts. Estimation by character count is acceptable for a visual approximation, but the UI must label it as estimated.

Recommendation: metadata should include `tokenSource: 'provided' | 'estimated'` and the caption should disclose when estimates are used.

### Gap: Thinking data may be unavailable or sensitive

The requested category includes "thinking". The implementation should only visualize thinking/scratchpad data that is explicitly present in the uploaded session. It should not infer or fabricate hidden model reasoning.

Recommendation: call the legend label "Thinking / scratchpad" only when the source contains explicit scratchpad records; otherwise use "Generated notes" or omit the category.

## Decision Records

### Decision: Use `styleKey + ContextStyleSet` as the block-category contract

- **Context:** The recently updated design system replaced hardcoded context kinds with caller-defined style keys and style sets.
- **Options considered:** Reuse old categories, introduce a new enum, or use caller-defined style keys.
- **Decision:** Use caller-defined `styleKey` values for all uploaded-session block types.
- **Rationale:** The existing types require `styleKey`, and the palette/legend system is designed for caller-defined vocabularies.
- **Consequences:** The normalizer owns semantic classification. React components stay generic.
- **Status:** proposed.

### Decision: Start with `ContextDiagramPanel` before adding new components

- **Context:** The user asked not to start actual work and asked to point out where the design system may need updates.
- **Options considered:** Build a new custom context-window visualizer immediately, or reuse existing diagram widgets first.
- **Decision:** Phase 2 should emit existing `ContextDiagramPanel` Widget IR. New components are Phase 4 decisions.
- **Rationale:** Existing components already support proportional sizing, palettes, legends, detail previews, and view switching.
- **Consequences:** The first implementation may not show turn separators perfectly, but it will be shippable and low-risk.
- **Status:** proposed.

### Decision: Treat turn paging as data-dependent

- **Context:** The request says individual turns might not be possible.
- **Options considered:** Always derive turn pages from transcript order, require exact per-turn context snapshots, or hide turn paging until schema support exists.
- **Decision:** Do not provide page-through UI in the current upload/visualize/transcript scope. If a later ticket reopens it, provide page-through UI only when the source has reliable turn boundaries, and provide exact per-turn context windows only when the source has exact per-turn membership/snapshot data.
- **Rationale:** Derived cumulative views can misrepresent compaction, eviction, or hidden context construction.
- **Consequences:** The normalizer must expose `hasReliableTurnBoundaries` and `hasPerTurnWindowSnapshots` and show warnings when approximating.
- **Status:** proposed.

### Decision: Keep block detail in metadata/content preview

- **Context:** Blocks need to show meaningful contents without rendering huge raw outputs inline.
- **Options considered:** Put full text in labels, put full text in detail only, or use external detail routes.
- **Decision:** Keep labels compact; store preview/content and source metadata on `ContextWindowPart`.
- **Rationale:** `ContextDiagramPanel` already renders selected part previews and metadata.
- **Consequences:** Large content should be truncated or lazy-loaded if production uploads are big.
- **Status:** proposed.

## Testing and Validation Strategy

### Unit tests

Add normalizer tests once implementation begins:

- classifies system/user/assistant/tool records correctly,
- distinguishes tool calls from tool results,
- recognizes file read/write records where schema supports it,
- estimates tokens when missing,
- preserves source IDs and turn IDs,
- emits warnings for inferred token counts and missing turn boundaries,
- coalesces tiny adjacent blocks only when configured.

### Component/Storybook validation

Run:

```bash
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build-storybook
```

Storybook acceptance checklist:

- Legend entries match uploaded-session vocabulary.
- Tiny blocks remain selectable.
- Large tool/file blocks dominate proportional views.
- Selected detail panel shows previews and metadata.
- Warnings are visible for inferred data.
- Page-through controls are keyboard accessible if added.

### Widget IR validation

If a new Widget IR type is added:

- update `packages/rag-evaluation-site/src/widgets/ir.ts`,
- add a widget adapter file,
- register it in `defaultRegistry.ts`,
- update `pkg/widgetschema/schema.go` if schema validation enumerates types,
- update `pkg/widgetdsl/typescript.go` and docs if DSL helpers are added,
- run Go tests for `pkg/widgetdsl`, `pkg/xgoja/providers/widgetsite`, and `pkg/widgetschema`.

## Risks and Open Questions

1. **What is the uploaded session schema?** This is the highest-priority question. Without it, turn grouping and file operation detection are speculative.
2. **Are token counts exact?** If not, the visualization should say "estimated" and avoid overclaiming precision.
3. **Does the uploaded data represent final context window membership?** A transcript is not the same as the actual context window if the runtime summarizes, evicts, or injects hidden state.
4. **Can thinking be shown safely?** Only explicit scratchpad/generated notes should be visualized.
5. **How dense are real sessions?** Hundreds or thousands of blocks may require aggregation, virtualization, or paging.
6. **Should this be frontend-only or backend-normalized?** Backend normalization is better for production, but a frontend Storybook prototype is faster.

## File Reference Map

- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/README.md` — high-level architecture, Widget IR frontend, and semantic recipe overview.
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/context/types.ts` — core `ContextStyleSet`, `ContextWindowPart`, `ContextWindowSnapshot`, and transcript types.
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/context/styles.ts` — palette/style-set creation and style resolution.
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/ContextStripDiagram/ContextStripDiagram.tsx` — ordered strip renderer and current width floor.
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/ContextStackDiagram/ContextStackDiagram.tsx` — readable stacked block renderer and height floor/ceiling.
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/molecules/ContextTreemap/ContextTreemap.tsx` — area-style view and current headroom filtering.
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/organisms/ContextDiagramPanel/ContextDiagramPanel.tsx` — view switcher, legend, and selected part detail panel.
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/widgets/ir.ts` — Widget IR type contracts for context diagram widgets.
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/pkg/widgetdsl/module.go` — Goja DSL helpers and context diagram recipe.
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/pkg/xgoja/providers/widgetsite/provider.go` — xgoja provider registration for `context_window.dsl`.
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/widgets/WidgetRenderer.context-diagrams.stories.tsx` — current Widget IR examples for content blocks and caller-defined legends.
- `/home/manuel/workspaces/2026-06-07/club-meetup-site/2026-05-27--rag-evaluation-system/packages/rag-evaluation-site/src/components/organisms/ContextUploadDropArea/ContextUploadDropArea.tsx` — presentational upload component.

## Implementation Handoff Checklist

Before coding:

- [ ] Add real uploaded session samples to the ticket `sources/` directory or document why samples cannot be stored.
- [ ] Decide frontend-only prototype vs backend normalizer.
- [ ] Decide exact block vocabulary and legend labels.
- [ ] Decide whether `thinking` is a valid source-visible category for the uploaded data.
- [ ] Decide whether turn paging is exact, approximate, or out of scope.

First coding pass:

- [ ] Add a normalizer and tests.
- [ ] Add a session block style set.
- [ ] Emit existing `ContextDiagramPanel` Widget IR.
- [ ] Add Storybook fixtures for realistic session block maps.
- [ ] Validate typecheck and Storybook build.

Only after first pass:

- [ ] Reconsider `ContextTurnPagerPanel` only if a later ticket explicitly reopens page-through/live-session scope.
- [ ] Consider grouped strip rendering.
- [ ] Consider aggregation/virtualization for dense sessions.

## Implementation Outcome — 2026-06-10

The implementation completed the local-first path described in this guide. No upstream React design-system component changes and no Goja DSL/schema changes were required for the first working uploaded-session block visualization.

Implemented in `ClubMedMeetup/minitrace-viz` commit `58ee4d3`:

- The uploaded-session context legend now uses block-specific style keys: `system`, `user`, `agent`, `thinking`, `tool_call`, `tool_result`, `file_read`, `file_write`, `summary`, `tool_error`, `active`, `empty`, and `other`.
- `buildContextWindowModel()` emits role-specific message blocks, thinking/scratchpad blocks, tool-call blocks, tool-result/tool-error blocks, and file read/write blocks.
- File rows and operation metadata drive `file_read` versus `file_write` classification where available.
- Successful file-heavy tool results are split so file blocks carry file-context weight and the residual tool-result block is a bounded status/result block.
- The original upload is preserved as `original-upload.txt` beside the converted minitrace archive so the local app can fall back when normalized minitrace rows are empty.
- The smoke fixture now asserts that uploaded-session context output contains `user`, `agent`, `tool_call`, `tool_result`, and `file_read` style keys.

Validated behavior from the sample upload:

```text
turns 6 tools 5
messages 11
style keys ['agent', 'empty', 'file_read', 'system', 'thinking', 'tool_call', 'tool_result', 'user']
```

Phase 3 and Phase 4 were closed as no-change decisions:

- `@go-go-golems/rag-evaluation-site` already supports the required rendering through `ContextDiagramPanel`, `ContextWindowPart.styleKey`, `ContextStyleSet`, and Widget IR props.
- `context_window.dsl` already supports the necessary style-set, part, snapshot, and panel helpers.
- Future reusable enhancements such as visible grouped turn separators or a multi-snapshot pager should be handled as separate Storybook-first design-system work if users request them.

Final validation included:

```bash
# ClubMed/minitrace-viz
pnpm --dir webapp typecheck
scripts/sync-widget-spa.sh
GOFLAGS=-buildvcs=false make test

# RAG package / Storybook / DSL
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build-storybook
go test ./pkg/widgetdsl ./pkg/xgoja/providers/widgetsite ./pkg/widgetschema -count=1

# Browser and visual smoke
Pi Playwright over course/upload/uploaded visualize/uploaded transcript/slides/handouts and three Storybook stories
css-visual-diff run: sources/visual-smoke/run-phase5-20260610-160514
```

## Post-validation addendum: grouped strip and turn pager components

After the Phase 5 validation pass, the reusable design-system scope was reopened for two Storybook-first context-window components:

1. `ContextGroupedStripDiagram`
2. `ContextTurnPagerPanel`

This supersedes the earlier Phase 4 no-change decision only for these explicitly requested reusable surfaces. The original local-first uploaded-session mapping decision remains unchanged.

### `ContextGroupedStripDiagram`

`ContextGroupedStripDiagram` is a molecule for context-window strips where visible group boundaries matter. It accepts one `ContextWindowSnapshot` and groups `ContextWindowPart` entries by one of three JSON-compatible keys:

```ts
type ContextGroupedStripGroupBy = 'turn' | 'styleKey' | 'sourceId';
```

The final grouping behavior aggregates all parts with the same group key while preserving first-seen group order. This matters because non-contiguous global/system/headroom entries should not produce duplicate `global` groups.

Default use in Widget IR:

```ts
component('ContextGroupedStripDiagram', {
  snapshot,
  styleSet,
  groupBy: 'turn',
  showGroupLabels: true,
  showPartLabels: true,
})
```

Design-system constraints:

- Use `styleKey + ContextStyleSet`; do not reintroduce `kind`/`ContextPartKind`.
- Do not show a blue selected outline in static mode unless `showSelection` is explicitly true.
- Use tokenized CSS module anatomy and foundation typography for hover details.
- Keep hover details inside the widget bounds so Storybook iframes and parent panels cannot clip them.

### `ContextTurnPagerPanel`

`ContextTurnPagerPanel` is an organism for stepping through multiple context-window snapshots. It accepts `snapshots: ContextWindowSnapshot[]`, renders Prev/Next and per-turn buttons, and displays the active snapshot with either the grouped strip or normal strip renderer.

The default behavior is **turn-only paging**, not selection-only paging:

```ts
mode = 'turn-only'
includeGlobalParts = true
```

In `turn-only` mode, the panel derives the active turn from snapshot metadata, filters the visible context window to only that turn's blocks, optionally retains global/system blocks, and recomputes headroom/free space. This avoids a misleading full-session strip where the pager only moves the selected block.

Callers that already provide fully materialized per-page snapshots can opt out:

```ts
component('ContextTurnPagerPanel', {
  snapshots,
  styleSet,
  mode: 'snapshot',
})
```

### Hover detail behavior

The strip and grouped-strip components now render hover details as a docked in-widget panel instead of a floating tooltip. This was chosen because floating popups were clipped by Storybook iframe bounds and by diagram borders.

The docked detail panel:

- appears below the strip/group caption inside the widget root,
- uses foundation `Text` primitives by default,
- can be customized with `renderPartTooltip?: (part: ContextWindowPart) => ReactNode`,
- avoids rounded popup styling,
- renders once per hovered/focused segment.

### Storybook review surfaces

Use these Storybook iframe URLs for focused visual review:

```text
/iframe.html?id=widget-ir-renderer-context-diagrams--context-grouped-strip-by-turn&viewMode=story
/iframe.html?id=widget-ir-renderer-context-diagrams--context-turn-pager-panel-story&viewMode=story
```

The grouped-by-turn story fixture was rebalanced so turn boundaries are visible in cropped widget screenshots. The expected visible groups are:

```text
global
turn 4
turn 5
turn 6
turn headroom
```

### Validation evidence

Validation commands:

```bash
pnpm --dir packages/rag-evaluation-site typecheck
pnpm --dir packages/rag-evaluation-site build-storybook
docmgr doctor --ticket CTX-WINDOW-BLOCK-VIZ --stale-after 30
```

Cropped widget screenshots were captured with Playwright and inspected via `read()` during implementation:

```text
context-grouped-strip-widget.png
context-turn-pager-widget.png
context-turn-pager-hover-widget-2.png
context-grouped-strip-hover-widget.png
```

The temporary screenshot files were removed after inspection; the diary records the observations.
