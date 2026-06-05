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
    - Path: examples/xgoja/widget-site/devctl/widget-site.py
      Note: NDJSON devctl plugin for rebuilding and launching the generated widget-site binary
    - Path: packages/rag-evaluation-site/src/app/app.css
      Note: Standalone app shell currently provides only min-height and no page chrome/padding/container
    - Path: packages/rag-evaluation-site/src/theme.css
      Note: Standalone package theme currently lacks mac/design-system tokens consumed by components
    - Path: ttmp/2026/06/05/WIDGETDSL-VISUAL-QUALITY--widget-dsl-visual-quality-and-rich-website-design/scripts/03-capture-visual-evidence.js
      Note: css-visual-diff JS verb for capturing Widget DSL and Storybook visual evidence
    - Path: web/src/styles/tokens.css
      Note: Original RAG app token source that defines --mac and font role variables
ExternalSources: []
Summary: ""
LastUpdated: 2026-06-05T13:40:00-04:00
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
