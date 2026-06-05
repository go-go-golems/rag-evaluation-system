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
ExternalSources: []
Summary: ""
LastUpdated: 2026-06-05T12:58:00-04:00
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

**Commit (code):** <built-in function hash> — "Add devctl workflow for xgoja widget site"

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
