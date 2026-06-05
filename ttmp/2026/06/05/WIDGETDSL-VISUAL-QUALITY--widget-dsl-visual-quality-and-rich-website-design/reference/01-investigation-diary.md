---
Title: Investigation Diary
Ticket: WIDGETDSL-VISUAL-QUALITY
Status: active
Topics:
    - frontend
    - ui-dsl
    - design-system
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: examples/xgoja/widget-site/.devctl.yaml
      Note: Example-local devctl plugin wiring
    - Path: examples/xgoja/widget-site/.gitignore
      Note: Ignores example-local devctl runtime state
    - Path: examples/xgoja/widget-site/README.devctl.md
      Note: Developer workflow instructions for devctl build/up/logs/down
    - Path: examples/xgoja/widget-site/assets/public/index.html
      Note: Synced rebuilt standalone app into generated xgoja example assets
    - Path: examples/xgoja/widget-site/devctl/widget-site.py
      Note: NDJSON devctl plugin for rebuilding and launching the generated widget-site binary
    - Path: examples/xgoja/widget-site/verbs/sites.js
      Note: Refactored xgoja showcase to use semantic recipes
    - Path: packages/rag-evaluation-site/src/app/App.tsx
      Note: Added default standalone shell rendering
    - Path: packages/rag-evaluation-site/src/app/app.css
      Note: |-
        Standalone app shell currently provides only min-height and no page chrome/padding/container
        Added shell/page chrome styles for root background
    - Path: packages/rag-evaluation-site/src/theme.css
      Note: |-
        Standalone package theme currently lacks mac/design-system tokens consumed by components
        Implemented standalone token bridge for --mac
    - Path: pkg/widgetdsl/module.go
      Note: Added page
    - Path: pkg/widgetdsl/module_test.go
      Note: Covered recipe/action/page helpers and JSON serialization
    - Path: pkg/xgoja/providers/widgetsite/doc/02-widget-dsl-js-api-reference.md
      Note: Documented action/page/recipe API
    - Path: ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/scripts/03-capture-visual-evidence.js
      Note: |-
        css-visual-diff JS verb for capturing Widget DSL and Storybook visual evidence
        Updated visual evidence script to target stable data-rag selectors
    - Path: ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-shell-01/01-visual-evidence-summary.md
      Note: Post-default-shell visual evidence summary
    - Path: ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-token-bridge-02/01-visual-evidence-summary.md
      Note: Post-token-bridge visual evidence summary
    - Path: web/src/styles/tokens.css
      Note: Original RAG app token source that defines --mac and font role variables
ExternalSources: []
Summary: ""
LastUpdated: 2026-06-05T18:35:00-04:00
WhatFor: ""
WhenToUse: ""
---






# Diary

## Goal

Capture the evidence-gathering and design process for improving the visual quality of Widget DSL output so DSL-authored sites can become rich, polished web experiences with little author code.

## Step 1: Create the Visual Quality Ticket and Define the Investigation

This step created a dedicated docmgr workspace for visual-quality research instead of mixing the work into the existing xgoja implementation ticket. The request is not just to tweak one demo page; it asks for a technical analysis, a clear design and implementation guide, visual comparisons, reusable scripts, and a reMarkable-ready deliverable.

The ticket will collect all scripts under its `scripts/` directory, keep screenshots and extracted evidence under ticket-local folders, and use the diary to preserve each command, failure, and conclusion. The intended reader is a new intern who needs enough architecture, API, visual-system, and workflow context to make concrete improvements without rediscovering the whole Widget DSL stack.

### Prompt Context

**User prompt (verbatim):** "ngl, look at the screenshots, this looks a bit like ass compared to the original rag site itself (see storybook, which you can launch and screenshot too, potentially use the tool `css-visual-diff help --all` to write some custom snreenshot / comparison code and jsverbs). Save all scripts in a enew ticket which is about improving the visual look of the widget dsl, and creating solid websites. The point of the DSL is creating rich web experiences that look and work great, without writing a lt of code. 

Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**User prompt (verbatim):** "keep a frequent and detailed diary as you work so we can retrace your steps."

**Assistant interpretation:** Create a new docmgr ticket focused on Widget DSL visual quality, gather visual evidence from the current DSL/demo and original RAG Storybook components, use screenshot/comparison tooling where useful, save scripts in the ticket, write a detailed intern-facing implementation guide, and upload the final bundle to reMarkable.

**Inferred user intent:** Preserve the low-code value of the DSL while raising the resulting site quality so DSL-authored pages look and behave like first-class RAG web experiences instead of crude generated layouts.

**Commit (code):** N/A — documentation/research setup only.

### What I did

- Created docmgr ticket `WIDGETDSL-VISUAL-QUALITY`.
- Created design document `design-doc/01-widget-dsl-visual-quality-analysis-and-implementation-guide.md`.
- Created diary document `reference/01-investigation-diary.md`.
- Added initial tasks for screenshot evidence, scripts, gap analysis, intern guide writing, validation, and reMarkable upload.
- Removed an accidental empty `2026/06/02` ticket stub after noticing docmgr created the real ticket under `2026/06/05`.

### Why

- The work needs isolated artifacts because it includes research scripts, visual evidence, design guidance, and delivery steps.
- A new intern needs a stable entry point: not only screenshots, but also the surrounding architecture and implementation sequence.
- Keeping a detailed diary from the beginning makes later review and continuation possible.

### What worked

- `docmgr ticket create-ticket` created the ticket workspace under `ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design`.
- `docmgr doc add` created the design doc and diary doc.
- The initial task list was added successfully.

### What didn't work

- I initially ran `mkdir -p` against a `2026/06/02/...` path from muscle memory. The real ticket was created under `2026/06/05/...`.
- Fix: removed the empty accidental stub and continued with the actual ticket path.

### What I learned

- For new docmgr tickets, trust the path printed by `docmgr ticket create-ticket` rather than inferring the date path manually.
- This ticket should be treated as a design/research deliverable first, not an implementation patch, because the requested output is an intern-ready guide and visual-quality plan.

### What was tricky to build

- The task spans documentation, visual inspection, browser tooling, and architecture planning. The main risk is collecting screenshots without connecting them to concrete file-level causes.
- To avoid that, the next steps will pair screenshots with file references in the renderer, component library, demo jsverb, and app shell.

### What warrants a second pair of eyes

- Review whether the proposed visual-quality work should result in changes to Widget IR schema, renderer defaults, example authoring conventions, or all three.
- Review whether the intern guide is specific enough to support implementation without hand-holding.

### What should be done in the future

- Gather current widget-site screenshots and original RAG/Storybook screenshots.
- Use `css-visual-diff help --all` and create ticket-local screenshot/comparison scripts.
- Write the final guide and upload it to reMarkable.

### Code review instructions

- Review the ticket setup at `ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design`.
- Confirm all future helper scripts live under that ticket's `scripts/` directory.
- Validate ticket hygiene with `docmgr doctor --ticket WIDGETDSL-VISUAL-QUALITY --stale-after 30`.

### Technical details

- Ticket id: `WIDGETDSL-VISUAL-QUALITY`
- Ticket path: `ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design`
- Primary design doc: `design-doc/01-widget-dsl-visual-quality-analysis-and-implementation-guide.md`
- Diary doc: `reference/01-investigation-diary.md`

## Step 2: Add an Example-Local devctl Rebuild and Launch Workflow

This step added a dedicated `devctl` setup inside `examples/xgoja/widget-site` so the generated Widget DSL site can be rebuilt, launched, supervised, logged, and smoked from one consistent workflow. This is important for the visual-quality ticket because screenshot and comparison work needs a reliable way to start the current widget-site without remembering the xgoja build command, the local `go-go-goja` replace path, or the correct HTTP listen flag.

The setup intentionally lives in the example directory rather than the repository root. The root repository already has a separate `.devctl.yaml` for the normal RAG backend/web stack. The example-local config lets a developer run `devctl up --force --timeout 10m` from `examples/xgoja/widget-site` and get exactly the generated xgoja widget-site service.

### Prompt Context

**User prompt (verbatim):** "commit at appropriate intervals as well,  including your previous work. Also create a devctl setup in the example site directory to easily recompile and launch the site."

**Assistant interpretation:** Commit the already-completed xgoja action-demo changes, then add and validate a devctl workflow in the widget-site example directory for rebuild/launch/smoke operations.

**Inferred user intent:** Keep the work checkpointed in git and make the example site easy to run during the upcoming visual-quality investigation.

**Commit (code):** 9acd175d23e3cc89ebd397b2c6f7e6c5b5c9336f — "Add devctl workflow for xgoja widget site"

### What I did

- Committed the previous xgoja action-demo work as:
  - `694095e7d99bdf3786156162d77e046db30c9420` — `Expand xgoja widget site action demo`.
- Committed the initial visual-quality ticket setup as:
  - `b6c3612b6b19d412762cc66fb83656fe0b43a610` — `Docs: create widget DSL visual quality ticket`.
- Read devctl guidance through:
  - `devctl help --all`
  - `devctl help user-guide`
  - `devctl help scripting-guide`
  - `devctl help plugin-authoring`.
- Added `examples/xgoja/widget-site/.devctl.yaml`.
- Added `examples/xgoja/widget-site/devctl/widget-site.py`.
- Added `examples/xgoja/widget-site/README.devctl.md`.
- Updated `examples/xgoja/widget-site/.gitignore` to ignore `.devctl/` runtime state.
- Added a WIDGETDSL-VISUAL-QUALITY task for the devctl setup and marked it complete.
- Related the devctl files to this diary and updated the ticket changelog.

### Why

- The visual comparison work needs a reproducible way to rebuild and run the current generated site.
- `make smoke` is useful for CI-like validation, but devctl gives a better local development loop: plan, build, up, status, logs, down.
- Keeping the setup in the example directory avoids colliding with the root repo's existing devctl environment.

### What worked

- Plugin discovery worked:

```text
cd examples/xgoja/widget-site
devctl plugins list --timeout 30s
```

- Planning worked and produced a single supervised service:

```text
devctl plan --timeout 30s
```

- Validation worked:

```text
devctl validate --timeout 30s
```

- Build worked after adding the expected per-step `ok` boolean:

```text
devctl build --timeout 10m
```

- Supervised launch worked:

```text
devctl up --force --timeout 10m
devctl status --tail-lines 5
curl -fsS http://127.0.0.1:18791/api/widget/pages/actions | grep -q 'xgoja widget actions demo'
devctl down
```

- Dynamic smoke command worked after stdout cleanup:

```text
devctl smoke --timeout 10m
```

### What didn't work

- First `devctl build --timeout 10m` displayed build steps as `ok: false` even though the underlying commands completed. The plugin returned step objects with `status`, but this devctl renderer expects an `ok` boolean on each step.
- Fix: changed `run_step(...)` to return both fields, for example:

```json
{"name":"build","ok":true,"status":"succeeded","duration_ms":4503}
```

- First `devctl smoke --timeout 10m` failed with protocol contamination:

```text
E_PROTOCOL_STDOUT_CONTAMINATION: cd /home/manuel/workspaces/2026-05-27/rag-evaluation-system/go-go-goja && GOWORK=off go run ./cmd/xgoja doctor ...: invalid character 'c' looking for beginning of value
```

- Cause: `command.run` used `subprocess.run(...)` without redirecting stdout, so `make smoke` wrote ordinary command output to the plugin's stdout. devctl plugin stdout must contain only NDJSON protocol frames.
- Fix: changed command helpers to route child stdout and stderr to `sys.stderr`:

```python
proc = subprocess.run(command_map[name], cwd=str(root), stdout=sys.stderr, stderr=sys.stderr)
```

### What I learned

- devctl's NDJSON boundary is strict in dynamic commands, not only handshake/plan. Any subprocess invoked by `command.run` must not inherit stdout.
- The installed devctl version renders build steps most clearly when each step includes `ok` as well as descriptive status/output fields.
- Example-local devctl configs are useful when a subdirectory has a distinct lifecycle from the root repo.

### What was tricky to build

- Path resolution had to account for three roots:
  - `repo_root`: the example directory passed to devctl.
  - `ragRoot`: `../../..`, the main RAG repo root.
  - `xgojaRoot`: `../../../../go-go-goja`, the sibling checkout used by `make build`.
- The plugin has to use `ctx.repo_root` instead of assuming the caller's process cwd, because devctl can be invoked with `--repo-root`.
- Build and command phases have different output constraints. `build.run` can capture command output and return summaries, while `command.run` must ensure child output goes to stderr to keep stdout protocol-clean.

### What warrants a second pair of eyes

- Review whether the plugin should expose more helper commands, such as `open-actions` or `browser-smoke`, after the visual comparison scripts exist.
- Review whether `devctl up` should always rebuild or whether a faster `--skip-build` workflow should be documented for screenshot iteration.
- Review whether dynamic command output through devctl's structured logs is readable enough for long `make smoke` output.

### What should be done in the future

- Use the devctl-managed site as the source server for screenshot capture scripts in this ticket.
- Add visual-diff helper commands later if the ticket-local scripts become stable enough to expose through devctl.
- Document any chosen ports or overrides once screenshot automation is finalized.

### Code review instructions

- Start with `examples/xgoja/widget-site/README.devctl.md` for intended usage.
- Review `examples/xgoja/widget-site/.devctl.yaml` for plugin wiring.
- Review `examples/xgoja/widget-site/devctl/widget-site.py` for protocol behavior and path resolution.
- Validate with:

```text
cd examples/xgoja/widget-site
devctl plugins list --timeout 30s
devctl plan --timeout 30s
devctl validate --timeout 30s
devctl build --timeout 10m
devctl up --force --timeout 10m
devctl status --tail-lines 5
devctl down
devctl smoke --timeout 10m
```

### Technical details

- Service name: `widget-site`
- Preferred port: `18791`
- Health URL: `http://127.0.0.1:18791/healthz`
- Action demo URL: `http://127.0.0.1:18791/pages/actions`
- Build phase:
  - `make sync-app`
  - `make build`
- Dynamic commands:
  - `devctl sync-app --timeout 10m`
  - `devctl smoke --timeout 10m`
  - `devctl clean`

## Step 3: Capture Visual Evidence and Identify the Token Boundary Failure

This step created the first visual evidence set for the Widget DSL visual-quality investigation. I used the new devctl workflow to rebuild and launch the generated xgoja widget-site, served the existing static Storybook build, and ran a ticket-local `css-visual-diff` JavaScript verb that captures overlay screenshots, selector preflight JSON, and computed CSS JSON for the generated page and several Storybook references.

The most important finding is that the generated page is not merely a weak composition; the standalone package is missing the original design-system token contract that the copied components expect. The component CSS uses `--mac-*`, `--font-body`, `--font-mono`, and `--rag-font-role-*` variables, but `packages/rag-evaluation-site/src/theme.css` only defines `--rag-*` tokens. As a result, CSS declarations such as `background: var(--mac-surface)` and `border: 1px solid var(--mac-border)` become invalid at computed-style time, making buttons, panels, and tables appear transparent or browser-default-like.

### Prompt Context

**User prompt (verbatim):** (same as Step 1 and Step 2)

**Assistant interpretation:** Gather visual evidence with screenshots and css-visual-diff, compare the generated Widget DSL page against Storybook/original RAG references, and keep the evidence and scripts in the ticket.

**Inferred user intent:** Move from subjective “looks bad” feedback to file-backed, screenshot-backed causes and an actionable implementation guide.

**Commit (code):** 8888e0d072f2b45b11c94caf74b3884ab608f2cb — "Docs: capture widget DSL visual evidence"

### What I did

- Read `css-visual-diff help --all` plus relevant help pages:
  - `getting-started`
  - `javascript-api`
  - `javascript-verbs`
  - `pixel-accuracy-scripting-guide`.
- Enumerated Storybook stories and selected visual references:
  - `pages-pipelinepage--populated`
  - `component-library-organisms-retrievalresultspanel--hybrid-results`
  - `widget-ir-renderer--search-workbench-composition`.
- Added ticket-local scripts:
  - `scripts/01-start-widget-and-storybook.sh`
  - `scripts/02-stop-widget-and-storybook.sh`
  - `scripts/03-capture-visual-evidence.js`
  - `scripts/04-run-visual-evidence.sh`.
- Ran the visual evidence workflow twice:
  - `run-01`: failed during a Storybook table selector screenshot.
  - `run-02`: succeeded after switching probe artifacts from full bundles to `css-json`.
- Attempted to use the image QA tool twice with fresh context and all comparison images.
- Inspected generated CSS JSON and source files to identify token mismatch and shell/app layout issues.

### Why

- The visual-quality guide needs evidence beyond subjective screenshot review.
- css-visual-diff artifacts can show both the visible page structure and the browser-computed styles that explain why the page looks wrong.
- The intern guide should explain root causes in terms of architecture boundaries: authoring DSL, renderer, component CSS, theme tokens, app shell, and example composition.

### What worked

- The devctl-managed widget-site startup worked inside the evidence script:

```text
examples/xgoja/widget-site devctl up --force --timeout 10m
```

- Static Storybook served successfully from `web/storybook-static` on port `6007`.
- The successful evidence run wrote artifacts under:

```text
ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-02
```

- Captured targets:
  - generated Widget DSL action page: `widget-actions.overlay.png`
  - original RAG PipelinePage story: `storybook-pipeline-page.overlay.png`
  - original RetrievalResultsPanel story: `storybook-retrieval-results.overlay.png`
  - WidgetRenderer Storybook composition: `widget-renderer-storybook.overlay.png`.
- Computed CSS evidence for the generated widget page showed:

```json
{
  "app-root": {
    "background-color": "rgba(0, 0, 0, 0)",
    "max-width": "none",
    "padding": "0px"
  },
  "first-button": {
    "background-color": "rgba(0, 0, 0, 0)",
    "border-radius": "0px",
    "height": "24px",
    "padding": "2px 12px"
  },
  "first-table": {
    "background-color": "rgba(0, 0, 0, 0)",
    "font-size": "11px"
  }
}
```

- Source inspection found the token mismatch:
  - `packages/rag-evaluation-site/src/theme.css` defines only `--rag-*` tokens.
  - component CSS such as `Button.module.css`, `Panel.module.css`, and `DataTable.module.css` consumes `--mac-*`, `--font-mono`, and role font tokens.
  - `web/src/styles/tokens.css` is the original source defining those variables.

### What didn't work

- The image QA tool failed twice with the same remote TLS error:

```text
pinocchio exited with code 1:
Error: inference failed: Post "https://api.openai.com/v1/responses": remote error: tls: bad record MAC
```

- I continued with direct artifact inspection rather than waiting on the remote service.

- `docmgr doctor` initially failed after evidence capture because css-visual-diff generated Markdown helper files (`index.md`, `computed-css.md`, and `README.md`) without docmgr frontmatter inside the ticket tree. I converted generated helper Markdown files to `.txt`, converted the successful run summary into `01-visual-evidence-summary.md` with frontmatter, and patched `scripts/04-run-visual-evidence.sh` so future runs perform the same sanitization automatically.

- `run-01` failed while trying to capture a selector screenshot for a Storybook `table` probe:

```text
Error: promise rejected: ArtifactError: context canceled
```

- Fix: changed `scripts/03-capture-visual-evidence.js` from `artifacts: "bundle"` to `artifacts: "css-json"`. Overlay screenshots still provide full-page visual evidence, and CSS JSON is enough for token/layout root-cause analysis.
- Generic Storybook selectors such as `table` and `button` sometimes hit hidden Storybook control elements instead of visible story content. The preflight output recorded this as `visible=false` and `bounds={width:0,height:0}` for some Storybook probes.

### What I learned

- The generated page looks bad primarily because copied component CSS depends on design tokens that are absent from the standalone package theme.
- A correct implementation guide must start with the design-token contract before discussing fancier DSL composition improvements.
- `RagEvaluationSiteApp` currently renders only a raw root div around `WidgetRenderer`; it does not provide app-shell chrome, page padding, sidebar/navigation, or a constrained content surface unless the Widget IR author explicitly emits `AppShell`.
- The generated example uses many raw layout widgets (`Stack`, `DashboardGrid`, `Panel`, `Button`, `DataTable`) but no high-level page template. That puts too much visual composition responsibility on the DSL author.

### What was tricky to build

- The tricky part was separating three causes that all contribute to the poor screenshot:
  1. missing CSS variables in the standalone package,
  2. lack of a default app shell/page template,
  3. example composition that uses low-level primitives without visual presets.
- css-visual-diff was useful because the computed CSS exposed missing token symptoms directly: transparent backgrounds, zero panel borders, and root padding of `0px`.
- Storybook iframes include Storybook-managed DOM, so broad selectors can accidentally hit hidden controls. Future scripts should prefer stable `data-rag-*` attributes in components and stories.

### What warrants a second pair of eyes

- Confirm whether the standalone package should import the exact original `web/src/styles/tokens.css` token set or define a cleaner package-local token bridge that maps `--mac-*` to `--rag-*` aliases.
- Review whether `RagEvaluationSiteApp` should always wrap pages in `AppShell`, or whether the server/DSL author should choose a page template explicitly.
- Review whether high-level DSL recipes should be implemented in JS (`rag.dsl`) or React (`WidgetRenderer`/template registry).

### What should be done in the future

- Add `data-rag-component` / `data-rag-layout` attributes to component roots so visual tooling can target visible RAG widgets reliably.
- Add a package-level theme contract test that renders representative components and verifies critical computed values are not transparent/default due to missing variables.
- Build higher-level DSL page recipes for search, corpus explorer, workflow dashboard, and action dashboard so authors can create rich pages without hand-composing every panel.

### Code review instructions

- Start with `packages/rag-evaluation-site/src/theme.css` and compare it with `web/src/styles/tokens.css`.
- Inspect `packages/rag-evaluation-site/src/components/atoms/Button/Button.module.css`, `Panel.module.css`, and `DataTable.module.css` to see the missing token consumers.
- Inspect `packages/rag-evaluation-site/src/app/App.tsx` to see the raw `rag-evaluation-site-root` wrapper around `WidgetRenderer`.
- Inspect `examples/xgoja/widget-site/verbs/sites.js` to see how much low-level layout composition the example author currently writes.
- Re-run evidence with:

```text
ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/scripts/04-run-visual-evidence.sh \
  ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-N
```

### Technical details

- `css-visual-diff` command shape:

```text
css-visual-diff verbs \
  --repository ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/scripts \
  widget-dsl visual capture-evidence \
  ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-02 \
  --widgetUrl http://127.0.0.1:18791 \
  --storybookUrl http://127.0.0.1:6007 \
  --output json
```

- Key source files:
  - `packages/rag-evaluation-site/src/theme.css`
  - `packages/rag-evaluation-site/src/styles.css`
  - `packages/rag-evaluation-site/src/app/app.css`
  - `packages/rag-evaluation-site/src/app/App.tsx`
  - `packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx`
  - `web/src/styles/tokens.css`
  - `examples/xgoja/widget-site/verbs/sites.js`

## Step 4: Write the Intern-Facing Visual Quality Implementation Guide

This step converted the screenshot/CSS evidence and source-code inspection into a structured implementation guide for a new intern. The guide is intentionally technical and explicit: it explains the Widget DSL stack from JavaScript authoring through Goja, HTTP endpoints, React rendering, component CSS, and visual tooling before it proposes concrete fixes.

The guide's central finding is that the generated page looks poor primarily because the standalone package theme does not define the design tokens consumed by the copied component CSS. The recommendations therefore start with a token bridge, then move to app shell defaults, semantic DSL recipes, visual-test selectors, and a refactored xgoja example.

### Prompt Context

**User prompt (verbatim):** (same as Step 1 and Step 2)

**Assistant interpretation:** Write a detailed, technically precise design and implementation guide that teaches a new intern how the system works, what is wrong visually, and how to fix it in phases.

**Inferred user intent:** Make visual-quality work executable by someone new to the project, with enough source references and pseudocode to avoid rediscovery.

**Commit (code):** f3dc7b204b2b404b7ec784f0df1f2ec7542b856f — "Docs: write widget DSL visual quality guide"

### What I did

- Replaced the placeholder design doc with a full guide at:
  - `design-doc/01-widget-dsl-visual-quality-analysis-and-implementation-guide.md`.
- Covered these areas:
  - executive summary,
  - problem statement and scope,
  - glossary,
  - evidence collected,
  - five-layer architecture model,
  - root-cause analysis,
  - proposed architecture,
  - decision records,
  - phased implementation plan,
  - testing strategy,
  - intern checklist,
  - risks and open questions,
  - API reference summary,
  - visual evidence command reference,
  - recommended PR sequence.
- Related the guide to the key source files that shaped it.
- Updated ticket changelog.
- Normalized the task list so validation/reMarkable upload remains open until the actual upload is complete.

### Why

- The project needs more than a one-off CSS fix. The Widget DSL is supposed to produce rich websites with little author code, so the guide needs to explain the system-level path from low-level primitives to polished recipes.
- Intern-facing documentation must explain file ownership and runtime flow before asking the reader to change code.
- The guide makes the implementation order explicit: theme first, shell second, recipes third, visual selectors/tools fourth, example polish fifth.

### What worked

- The guide now anchors the visual root cause to concrete files:
  - `packages/rag-evaluation-site/src/theme.css`
  - `web/src/styles/tokens.css`
  - `packages/rag-evaluation-site/src/components/atoms/Button/Button.module.css`
  - `packages/rag-evaluation-site/src/components/layout/Panel/Panel.module.css`
  - `packages/rag-evaluation-site/src/components/molecules/DataTable/DataTable.module.css`.
- The guide anchors the runtime architecture to concrete files:
  - `pkg/widgetdsl/module.go`
  - `pkg/widgetrunner/runner.go`
  - `pkg/widgetserver/server.go`
  - `packages/rag-evaluation-site/src/app/App.tsx`
  - `packages/rag-evaluation-site/src/widgets/WidgetRenderer.tsx`.
- The guide includes Mermaid diagrams, TypeScript/JavaScript/CSS pseudocode, and command examples.

### What didn't work

- No new tool failure occurred during writing.
- I did notice the task list had a stale checked upload task from earlier ticket setup. I rewrote `tasks.md` so upload remains unchecked until the final upload command succeeds.

### What I learned

- The implementation guide should not frame the problem as “the demo author made an ugly page.” The more accurate framing is that the system has missing default visual infrastructure: tokens, shell, recipes, and visual selectors.
- The current low-level Widget DSL is still valuable. The fix is to layer polished defaults and recipes on top, not to discard the existing primitives.

### What was tricky to build

- The guide had to balance two audiences: a new intern who needs foundational explanation and an experienced reviewer who needs precise file references and actionable phases.
- The strongest recommendation is also the simplest one: restore token compatibility. It would be easy to over-focus on recipes before fixing the missing CSS variables, but computed-style evidence shows that the token bridge is the first implementation slice.

### What warrants a second pair of eyes

- Review the decision to preserve `--mac-*` tokens as a package-local bridge rather than immediately migrating all components to `--rag-*` variables.
- Review the proposed default shell behavior, especially the `meta.shell = 'none'` opt-out and AppShell double-wrap detection.
- Review whether recipes should live under `widget.dsl`, `rag.dsl`, or a separate namespace.

### What should be done in the future

- Implement Phase 1: standalone token bridge.
- Re-run visual evidence after Phase 1 and compare computed CSS.
- Use the guide's PR sequence to split future work into reviewable commits.

### Code review instructions

- Read the guide from the top through the root-cause analysis before reviewing the proposed implementation phases.
- Verify that every major claim has a file reference or captured evidence path.
- Validate docs with:

```text
docmgr doctor --ticket WIDGETDSL-VISUAL-QUALITY --stale-after 30
```

### Technical details

- Primary guide path:
  - `ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/design-doc/01-widget-dsl-visual-quality-analysis-and-implementation-guide.md`
- Evidence path:
  - `ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-02`
- Visual scripts path:
  - `ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/scripts`

## Step 5: Validate and Upload the Visual Quality Guide to reMarkable

This step completed the delivery loop for the visual-quality ticket. After writing and committing the intern-facing guide, I bundled the guide, diary, and visual evidence summary into one PDF and uploaded it to reMarkable.

The upload destination is ticket-specific so the document is easy to find later: `/ai/2026/06/05/WIDGETDSL-VISUAL-QUALITY`. The upload command returned a successful `OK: uploaded ...` response, so no additional cloud listing was required.

### Prompt Context

**User prompt (verbatim):** (same as Step 1 and Step 2)

**Assistant interpretation:** Validate the ticket artifacts and deliver the final guide bundle to reMarkable.

**Inferred user intent:** Make the analysis/design package reviewable on the reMarkable device and close the documentation loop cleanly.

**Commit (code):** 462b3230b2ccd35172ae15f1ae04a40ece4be958 — "Docs: record widget DSL visual guide upload"

### What I did

- Uploaded a reMarkable bundle containing:
  - `design-doc/01-widget-dsl-visual-quality-analysis-and-implementation-guide.md`
  - `reference/01-investigation-diary.md`
  - `sources/visual-evidence/run-02/01-visual-evidence-summary.md`.
- Marked task 8 complete.
- Updated the changelog with the upload destination.

### Why

- The user explicitly requested that the final guide be stored in the ticket and uploaded to reMarkable.
- Bundling guide + diary + evidence summary gives the reviewer both the final recommendation and the trace of how it was produced.

### What worked

- Upload command:

```text
remarquee upload bundle \
  ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/design-doc/01-widget-dsl-visual-quality-analysis-and-implementation-guide.md \
  ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/reference/01-investigation-diary.md \
  ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-02/01-visual-evidence-summary.md \
  --name "Widget DSL Visual Quality Guide" \
  --remote-dir "/ai/2026/06/05/WIDGETDSL-VISUAL-QUALITY" \
  --toc-depth 2 \
  --non-interactive
```

- Upload result:

```text
OK: uploaded Widget DSL Visual Quality Guide.pdf -> /ai/2026/06/05/WIDGETDSL-VISUAL-QUALITY
```

### What didn't work

- N/A.

### What I learned

- The reMarkable upload path works cleanly for a concise bundle of the guide, diary, and evidence summary.
- Keeping generated visual evidence Markdown sanitized earlier prevented docmgr validation problems before delivery.

### What was tricky to build

- The main delivery concern was choosing the right bundle contents. Including all raw visual artifacts would make the PDF noisy, so the bundle uses the primary guide, chronological diary, and compact evidence summary. Raw JSON/PNG artifacts remain in the ticket workspace.

### What warrants a second pair of eyes

- Confirm whether the raw overlay screenshots should also be converted into a separate image-heavy PDF later. The current uploaded bundle references them but does not embed every PNG.

### What should be done in the future

- Start implementation with Phase 1 from the guide: add the standalone token bridge and rerun css-visual-diff evidence.

### Code review instructions

- Validate ticket hygiene with:

```text
docmgr doctor --ticket WIDGETDSL-VISUAL-QUALITY --stale-after 30
```

- Review the uploaded source documents in the ticket if the reMarkable PDF needs regeneration.

### Technical details

- reMarkable destination: `/ai/2026/06/05/WIDGETDSL-VISUAL-QUALITY`
- Uploaded PDF name: `Widget DSL Visual Quality Guide.pdf`

## Step 6: Implement the Standalone Theme Token Bridge

This step implemented Phase 1 from the visual-quality guide: restore the design-token contract that the standalone `@go-go-golems/rag-evaluation-site` package must provide to the copied RAG components. The implementation adds a bridge in `packages/rag-evaluation-site/src/theme.css` that defines the legacy/original `--mac-*`, `--font-*`, and `--rag-font-role-*` variables in terms of the package's canonical `--rag-*` tokens.

The result is that the generated xgoja Widget DSL site no longer renders primary buttons and panels as transparent/default-looking controls. Smoke tests passed, and post-bridge visual evidence shows panel roots computing to a white surface with a black border and primary buttons computing to black background, white text, mono font, and a real border.

### Prompt Context

**User prompt (verbatim):** "Add tasks, implement, commit at appropriate intervals, keep a diary, use smoke tests"

**Assistant interpretation:** Add implementation tasks to the existing visual-quality ticket, implement the next recommended fix, validate with smoke tests and visual evidence, keep the diary current, and commit at focused checkpoints.

**Inferred user intent:** Move from the analysis/design package into concrete implementation while preserving traceability and confidence.

**Commit (code):** bc737b87f65f883cac35c49f3c93da1918992948 — "Add standalone widget site token bridge"

### What I did

- Added WIDGETDSL-VISUAL-QUALITY tasks for:
  - implementing the Phase 1 token bridge,
  - rebuilding/syncing assets and running smoke tests,
  - rerunning css-visual-diff evidence.
- Updated `packages/rag-evaluation-site/src/theme.css` with:
  - canonical standalone tokens,
  - `--mac-*` compatibility variables,
  - `--font-body` / `--font-mono`,
  - `--rag-font-role-*` role aliases.
- Rebuilt the package:

```text
cd packages/rag-evaluation-site
pnpm build
pnpm build:app
```

- Synced app assets into the generated xgoja example:

```text
make -C examples/xgoja/widget-site sync-app
```

- Ran generated-site smoke:

```text
make -C examples/xgoja/widget-site smoke
```

- Ran devctl helper smoke:

```text
cd examples/xgoja/widget-site
devctl smoke --timeout 10m
```

- Ran visual evidence capture:

```text
ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/scripts/04-run-visual-evidence.sh \
  ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-token-bridge-02
```

- Updated `scripts/03-capture-visual-evidence.js` to prefer stable runtime attributes:
  - `[data-rag-layout="Panel"]`
  - `[data-rag-component="DataTable"]`
  - `[data-rag-atom="Button"]`.
- Removed the intermediate `run-token-bridge-01` evidence directory and kept the cleaner `run-token-bridge-02` output.
- Checked tasks 9, 10, and 11.
- Updated changelog and related files.

### Why

- The visual guide identified missing token definitions as the first and most important root cause.
- Component CSS already depended on original RAG tokens; adding the bridge is the smallest safe fix that restores component chrome without rewriting component styles.
- Smoke tests ensure the theme change did not break package builds, generated xgoja builds, API actions, or embedded SPA serving.

### What worked

- `pnpm build` passed for the package build.
- `pnpm build:app` passed for the standalone app build.
- `make -C examples/xgoja/widget-site smoke` passed.
- `devctl smoke --timeout 10m` passed.
- `docmgr doctor --ticket WIDGETDSL-VISUAL-QUALITY --stale-after 30` passed after evidence capture.
- Post-bridge computed CSS assertions passed:

```text
ASSERTIONS: token bridge visual CSS checks passed
```

- Key post-bridge computed CSS for the first Panel:

```json
{
  "background-color": "rgb(255, 255, 255)",
  "border": "1px solid rgb(0, 0, 0)",
  "border-radius": "0px",
  "box-shadow": "none",
  "padding": "0px"
}
```

- Key post-bridge computed CSS for the first primary Button:

```json
{
  "background-color": "rgb(0, 0, 0)",
  "border": "1px solid rgb(0, 0, 0)",
  "color": "rgb(255, 255, 255)",
  "font-family": "SFMono-Regular, Consolas, \"Liberation Mono\", monospace",
  "height": "20.8438px",
  "padding": "2px 12px"
}
```

### What didn't work

- The first comparison against `run-02` suggested the panel had not changed because the old visual script's `first-panel` selector matched the outer stack/root rather than an actual Panel.
- Fix: inspected the runtime DOM with Playwright and found stable attributes already present, including:
  - `data-rag-layout="Panel"`,
  - `data-rag-component="DataTable"`,
  - `data-rag-atom="Button"`.
- I updated the visual script to use those attributes and reran evidence into `run-token-bridge-02`.

### What I learned

- The token bridge immediately fixes component chrome for panels and buttons, which confirms the root-cause analysis.
- The app root still has `padding: 0px` and transparent background. That is expected after Phase 1; Phase 2 in the guide should add the default app shell/page chrome.
- Runtime components already expose useful `data-rag-*` attributes, so visual tooling should use them instead of broad selectors.

### What was tricky to build

- The tricky part was distinguishing a real styling failure from a probe/selector failure. The component did get styled after the bridge, but the original css-visual-diff `first-panel` probe did not target the actual Panel element.
- I used Playwright to inspect live DOM attributes and computed styles, then changed the evidence script so future runs are more reliable.
- Asset sync produced new hashed app assets under `examples/xgoja/widget-site/assets/public/assets`, so old hashed assets appear as deleted and new hashed assets appear as added. This is expected after rebuilding the standalone app.

### What warrants a second pair of eyes

- Review the exact bridge values, especially `--mac-border: var(--rag-color-border-strong)` versus using the softer `--rag-color-border`.
- Review whether the standalone package should preserve the classic black/white RAG look or move toward a softer modern theme after Phase 2.
- Review the updated visual-evidence selectors for Storybook stories; some stories may not contain every `data-rag-*` target, but misses are recorded in preflight JSON.

### What should be done in the future

- Implement Phase 2 from the guide: add app shell/page chrome so the root no longer has `padding: 0px` and raw edge-to-edge layout.
- Consider adding a small computed-style assertion script to CI after the visual selector cleanup settles.
- Optionally regenerate the reMarkable bundle after Phase 2 or after a larger visual-polish pass.

### Code review instructions

- Start with `packages/rag-evaluation-site/src/theme.css` and review the token bridge block.
- Check the rebuilt assets in `examples/xgoja/widget-site/assets/public` only for expected hashed asset replacement.
- Review `scripts/03-capture-visual-evidence.js` to confirm stable `data-rag-*` selectors are used.
- Validate with:

```text
cd packages/rag-evaluation-site && pnpm build && pnpm build:app
make -C examples/xgoja/widget-site sync-app
make -C examples/xgoja/widget-site smoke
cd examples/xgoja/widget-site && devctl smoke --timeout 10m
```

- Validate visual evidence with:

```text
docmgr doctor --ticket WIDGETDSL-VISUAL-QUALITY --stale-after 30
```

### Technical details

- Theme file changed:
  - `packages/rag-evaluation-site/src/theme.css`
- Synced xgoja embedded app files changed:
  - `examples/xgoja/widget-site/assets/public/index.html`
  - `examples/xgoja/widget-site/assets/public/assets/index-*.css`
  - `examples/xgoja/widget-site/assets/public/assets/index-*.js`
- Visual evidence script changed:
  - `ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/scripts/03-capture-visual-evidence.js`
- Post-bridge evidence path:
  - `ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-token-bridge-02`

## Step 7: Add the Default Standalone App Shell

This step implemented Phase 2 from the visual-quality guide: standalone Widget DSL pages now get a default application frame unless they explicitly opt out or already render an `AppShell` as their root. The default frame wraps raw Widget IR pages with the existing React component library's `AppShell` and `AppNav`, adds page padding and a bordered shell, and preserves generated script control through `page.meta`.

The result is that `/pages/actions` no longer looks like a component fragment placed directly against the viewport. Browser smoke confirmed the shell is present, root background and padding are applied, the action workflow still refreshes rows correctly, and the page emits no console warnings or errors.

### Prompt Context

**User prompt (verbatim):** "go ahead"

**Assistant interpretation:** Continue from the previous token-bridge checkpoint and implement the next recommended visual-quality phase.

**Inferred user intent:** Keep moving through the visual-quality implementation plan without waiting for another detailed instruction.

**Commit (code):** de741f7ebc8cc14bb0f71eca8515342b9a332c06 — "Add default shell for widget site app"

### What I did

- Added WIDGETDSL-VISUAL-QUALITY tasks for:
  - Phase 2 default app shell/page chrome,
  - rebuild/sync/smoke after shell changes,
  - post-shell visual evidence capture and comparison.
- Updated `packages/rag-evaluation-site/src/app/App.tsx` to:
  - normalize `page.meta.shell`, `page.meta.navItems`, `page.meta.activeNavItemId`, and `page.meta.maxWidth`,
  - default-wrap pages in `AppShell` unless `meta.shell === "none"` or the root Widget IR node is already an `AppShell`,
  - render a default `AppNav` with Overview/Demo/Actions entries,
  - expose `data-rag-shell="default"` and `data-rag-layout="PageContent"` for smoke and visual tooling,
  - navigate default nav items to `/pages/{id}`.
- Updated `packages/rag-evaluation-site/src/app/app.css` to:
  - set root background/color/font,
  - add 8px shell-mode viewport padding,
  - add a bordered shell frame,
  - add content width modes for `none`, `content`, and `wide`.
- Rebuilt package artifacts:

```text
cd packages/rag-evaluation-site
pnpm build
pnpm build:app
```

- Synced and rebuilt the generated xgoja example:

```text
make -C examples/xgoja/widget-site sync-app
make -C examples/xgoja/widget-site smoke
```

- Ran the devctl smoke workflow:

```text
cd examples/xgoja/widget-site
devctl smoke --timeout 10m
```

- Captured post-shell visual evidence:

```text
ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/scripts/04-run-visual-evidence.sh \
  ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-shell-01
```

- Ran a Pi Playwright browser smoke at `http://127.0.0.1:18798/pages/actions`:
  - verified `data-rag-shell="default"`,
  - verified root padding `8px`,
  - verified root background `rgb(246, 247, 248)`,
  - verified shell border `1px solid rgb(0, 0, 0)`,
  - clicked `Add query`,
  - verified `Rows: 4` and `Follow-up Query 4`,
  - verified console warnings/errors count was zero,
  - captured screenshot `widgetsite-shell-smoke-2026-06-05.png`.

### Why

- Phase 1 restored component-level styling, but the page still lacked a standalone application frame.
- The guide recommended adding app chrome at the app boundary rather than requiring every Widget IR script to author the same shell manually.
- Opt-out and double-wrap detection protect embedded/raw pages and pages that already intentionally render their own `AppShell`.

### What worked

- TypeScript and Vite builds passed.
- Generated xgoja build/smoke passed.
- devctl smoke passed.
- css-visual-diff evidence capture passed.
- Browser action smoke passed with no console warnings/errors.
- Computed CSS improved from the token-bridge run:

```text
run-token-bridge-02 app-root:
  background-color: rgba(0, 0, 0, 0)
  padding: 0px
  font-size: 16px

run-shell-01 app-root:
  background-color: rgb(246, 247, 248)
  padding: 8px
  font-size: 12px
```

- Post-shell assertions passed:

```text
ASSERTIONS: default shell visual CSS checks passed
```

### What didn't work

- The first shell implementation placed a `<main>` element inside `AppShell`, which already renders its children inside a `<main>`.
- Fix: changed the shell content wrapper from `<main>` to `<div data-rag-layout="PageContent">` and reran builds/smokes/evidence.

### What I learned

- The existing `AppShell` component is sufficient for the default shell; the standalone app only needed boundary logic and CSS overrides.
- The post-shell change should be measured on the root element, not only component internals: root background, padding, and shell border are the clearest evidence.
- The action flow remained stable after adding chrome because the `WidgetRenderer` and action callback were preserved unchanged inside the shell.

### What was tricky to build

- The main risk was adding chrome without breaking pages that are already authored as complete app shells. The solution was a narrow `isAppShellNode(page.root)` guard plus `meta.shell = "none"` opt-out.
- Another sharp edge was navigation behavior. The default nav should work without requiring server actions, so items without explicit actions navigate to `/pages/{id}`. Items with explicit actions still call the app action handler.
- The nested-main issue was easy to miss because it compiled and rendered. I corrected it before final evidence so the committed markup is semantically cleaner.

### What warrants a second pair of eyes

- Review whether the default nav entries should remain hard-coded as Overview/Demo/Actions or be omitted unless `page.meta.navItems` is provided.
- Review the default max width: this implementation uses `wide` by default to fit the action dashboard; a more document-like page may prefer `content`.
- Review whether the default shell should expose a prop-level override in addition to `page.meta.shell` for embedding consumers.

### What should be done in the future

- Add focused React tests for shell wrapping, `meta.shell = "none"`, and AppShell-root double-wrap detection.
- Document the `page.meta` shell contract in provider help docs or schema output.
- Continue with Phase 3 semantic `rag.dsl` recipes if more visual polish is desired.

### Code review instructions

- Start in `packages/rag-evaluation-site/src/app/App.tsx`:
  - `renderPage`,
  - `shouldUseDefaultShell`,
  - `normalizeMeta`,
  - `navigateToPage`.
- Then review `packages/rag-evaluation-site/src/app/app.css` for root/shell/content styles.
- Check the rebuilt xgoja assets under `examples/xgoja/widget-site/assets/public` only for expected hashed file replacement.
- Validate with:

```text
cd packages/rag-evaluation-site && pnpm build && pnpm build:app
make -C examples/xgoja/widget-site sync-app
make -C examples/xgoja/widget-site smoke
cd examples/xgoja/widget-site && devctl smoke --timeout 10m
```

- Validate visual evidence with:

```text
docmgr doctor --ticket WIDGETDSL-VISUAL-QUALITY --stale-after 30
```

### Technical details

- App shell implementation:
  - `packages/rag-evaluation-site/src/app/App.tsx`
- App shell styling:
  - `packages/rag-evaluation-site/src/app/app.css`
- Rebuilt embedded app assets:
  - `examples/xgoja/widget-site/assets/public/index.html`
  - `examples/xgoja/widget-site/assets/public/assets/index-*.css`
  - `examples/xgoja/widget-site/assets/public/assets/index-*.js`
- Visual evidence:
  - `ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/sources/visual-evidence/run-shell-01`
- Browser screenshot:
  - `widgetsite-shell-smoke-2026-06-05.png`

## Step 8: Add Semantic Widget DSL Recipes

This step implemented Phase 3 from the visual-quality guide: `widget.dsl`/`rag.dsl` now include semantic helpers for common RAG dashboard patterns. The new helpers keep the low-level Widget IR contract intact, but allow scripts to express page intent through `rag.page`, `rag.action.*`, and `rag.recipes.*` instead of manually assembling every dashboard grid, action toolbar, and master/detail layout.

The xgoja widget-site showcase was refactored to use these recipes while preserving the stateful action dashboard behavior. Go tests, generated xgoja smoke, devctl smoke, generated help checks, and Pi Playwright browser action smoke all passed.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue the visual-quality implementation plan after Phase 2 shell/page chrome, choosing the next practical slice from the guide.

**Inferred user intent:** Keep improving the Widget DSL authoring experience and example quality without stopping at visual shell fixes.

**Commit (code):** 0a4b156fc768023854a714c1a16f572cb633e70e — "Add semantic widget DSL recipes"

### What I did

- Added WIDGETDSL-VISUAL-QUALITY tasks for Phase 3 recipe implementation, showcase refactor, and validation.
- Extended `pkg/widgetdsl/module.go` with:
  - `rag.page(options)`,
  - `rag.action.server`, `rag.action.navigate`, `rag.action.event`, `rag.action.copy`,
  - `rag.recipes.metrics`,
  - `rag.recipes.actionToolbar`,
  - `rag.recipes.masterDetailTable`.
- Added `pkg/widgetdsl/module_test.go` coverage proving recipes/actions/page helpers produce JSON-serializable Widget IR.
- Refactored `examples/xgoja/widget-site/verbs/sites.js` to use:
  - `rag.recipes.metrics` for summary cards,
  - `rag.recipes.actionToolbar` for queue controls,
  - `rag.recipes.masterDetailTable` for the query queue + selected query layout,
  - `rag.page` for the page wrapper and default stacked sections.
- Updated provider-bundled help docs:
  - `pkg/xgoja/providers/widgetsite/doc/01-widget-dsl-getting-started.md`,
  - `pkg/xgoja/providers/widgetsite/doc/02-widget-dsl-js-api-reference.md`.
- Validated generated help content includes the new recipe docs.

### Why

- Phase 1 and Phase 2 fixed visual environment problems, but scripts were still verbose because authors had to manually compose common RAG page patterns.
- Recipes make example code read closer to the dashboard intent while still returning ordinary JSON-compatible Widget IR.
- Keeping recipes in `widget.dsl` means both standalone scripts and generated xgoja binaries get the same authoring conveniences.

### What worked

- `go test ./pkg/widgetdsl ./pkg/xgoja/providers/widgetsite -count=1` passed.
- `make -C examples/xgoja/widget-site smoke` passed after fixing slice normalization.
- `cd examples/xgoja/widget-site && devctl smoke --timeout 10m` passed.
- Generated help checks passed:

```text
examples/xgoja/widget-site/dist/rag-widget-xgoja-site help widget-dsl-js-api-reference | grep -q 'recipes.masterDetailTable'
examples/xgoja/widget-site/dist/rag-widget-xgoja-site help widget-dsl-getting-started | grep -q 'Use recipes for common dashboards'
```

- Pi Playwright browser smoke at `http://127.0.0.1:18800/pages/actions` passed:
  - shell was `default`,
  - `DataTable` existed,
  - initial page had `Rows: 3`, `Queue controls`, `Query queue`, and `Fast Growing Trees`,
  - clicking `Add query` refreshed to `Rows: 4`, `Follow-up Query 4`, and `Added query #4`,
  - console warnings/errors count was zero,
  - screenshot captured as `widgetsite-recipes-smoke-2026-06-05.png`.

### What didn't work

- First `make -C examples/xgoja/widget-site smoke` failed after the recipe refactor.
- Manual reproduction showed the refreshed page had `Rows: 4`, but the recipe-backed `DataTable` had `rows: []`, so smoke greps for `Follow-up Query` and `Bulk retry requested from Widget IR action` failed.
- Root cause: `anySlice` only accepted `[]any`, while database query rows crossed Goja as a typed Go slice.
- Fix: changed `anySlice` to use reflection for slices/arrays and convert them into `[]any`.

### What I learned

- Goja function binding requires functions that accept `goja.FunctionCall` to return `goja.Value` when they need the full call object; the initial `page(call goja.FunctionCall) map[string]any` shape did not receive arguments correctly.
- Recipe helpers need robust normalization because inputs may be plain JS arrays, Go-backed typed slices, or nested exported maps.
- The recipes can simplify example code without changing the renderer contract: the resulting page still contains `Panel`, `DashboardGrid`, `DataTable`, `Button`, and `MetadataGrid` nodes.

### What was tricky to build

- The main tricky part was preserving JSON compatibility while allowing `masterDetailTable.detail` to be a JavaScript callback. The callback is only evaluated while constructing the Widget IR response; the returned value is exported as a plain Widget IR node before crossing the HTTP JSON boundary.
- Another sharp edge was string action shorthand. The recipe layer accepts either a full action spec or a string server action name, but it always expands to a plain `{ kind: "server", name }` object.
- Typed-slice normalization was a runtime-only problem that package tests did not catch until the generated xgoja smoke exercised DB query rows.

### What warrants a second pair of eyes

- Review whether recipes belong permanently in `widget.dsl` or should eventually be separated into a richer `rag.dsl` namespace while keeping `widget.dsl` minimal.
- Review the `masterDetailTable` selected-row lookup, which currently assumes row IDs are under `id` for the callback convenience path.
- Review whether recipe helpers should validate required fields more strictly or remain permissive like the low-level component helpers.

### What should be done in the future

- Add a test that uses a typed Go slice as recipe `rows` input to prevent regressions in `anySlice`.
- Document `rag.page` metadata options for shell behavior in the generated schema/help docs.
- Consider richer recipes for search workbench pages, retrieval result comparison, and evaluation run timelines.

### Code review instructions

- Start in `pkg/widgetdsl/module.go`:
  - `page`,
  - `metricsRecipe`,
  - `actionToolbarRecipe`,
  - `masterDetailTableRecipe`,
  - `anySlice`,
  - `normalizeActionSpec`.
- Then review `pkg/widgetdsl/module_test.go` for the JSON-serialization test.
- Review `examples/xgoja/widget-site/verbs/sites.js` to confirm the recipe refactor did not remove behavior.
- Validate with:

```text
go test ./pkg/widgetdsl ./pkg/xgoja/providers/widgetsite -count=1
make -C examples/xgoja/widget-site smoke
cd examples/xgoja/widget-site && devctl smoke --timeout 10m
```

### Technical details

- Recipe implementation:
  - `pkg/widgetdsl/module.go`
- Recipe tests:
  - `pkg/widgetdsl/module_test.go`
- Recipe-backed showcase:
  - `examples/xgoja/widget-site/verbs/sites.js`
- Provider docs:
  - `pkg/xgoja/providers/widgetsite/doc/01-widget-dsl-getting-started.md`
  - `pkg/xgoja/providers/widgetsite/doc/02-widget-dsl-js-api-reference.md`
- Browser screenshot:
  - `widgetsite-recipes-smoke-2026-06-05.png`
