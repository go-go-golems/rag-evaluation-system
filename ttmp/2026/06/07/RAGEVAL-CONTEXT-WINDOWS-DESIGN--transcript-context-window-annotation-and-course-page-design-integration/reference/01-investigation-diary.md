---
Title: Investigation Diary
Ticket: RAGEVAL-CONTEXT-WINDOWS-DESIGN
Status: active
Topics:
    - design-system
    - frontend-architecture
    - react
    - rag
    - ui-dsl
DocType: reference
Intent: long-term
Owners: []
RelatedFiles:
    - Path: packages/rag-evaluation-site/src/context/fixtures.ts
      Note: Normalized prototype fixtures for package Storybook and future Widget IR
    - Path: packages/rag-evaluation-site/src/context/index.ts
      Note: Context module barrel
    - Path: packages/rag-evaluation-site/src/context/types.ts
      Note: New context-viewer DTO contracts
    - Path: packages/rag-evaluation-site/src/index.ts
      Note: Exports context contracts from package API
    - Path: ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/scripts/01_capture_go_minitrace_help.sh
      Note: Captured go-minitrace help pages for reproducible query setup
    - Path: ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/scripts/02_convert_relevant_pi_sessions.sh
      Note: Converted relevant Pi sessions for transcript search
    - Path: ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/scripts/04_git_frontend_history.sh
      Note: Captured frontend git history after root-path fix
    - Path: ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/scripts/05_catalog_recovered_design_docs.py
      Note: Catalog script with corrected repo-root inference
ExternalSources: []
Summary: ""
LastUpdated: 2026-06-07T00:00:00Z
WhatFor: ""
WhenToUse: ""
---



# Diary

## Goal

Capture the research trail for integrating the new design iteration into a full transcript/context-window/annotation/course page experience, starting from prior Pi session evidence about frontend design-system and web/ playbooks.

## Step 1: Ticket Setup and go-minitrace Bootstrap

I created a dedicated docmgr ticket for the design-integration work and initialized the expected design document, diary, task list, and scripts workspace. The immediate objective is not to design from memory: it is to recover prior design-system/playbook knowledge from local Pi transcripts and make the search reproducible.

I started by following the go-minitrace transcript-analysis workflow: capture the built-in help tree, save reusable scripts under the ticket `scripts/` directory, and convert only relevant Pi session stores for this RAG evaluation workspace.

### Prompt Context

**User prompt (verbatim):** "Create a new ticket to integrate the new design iteration to make a full transacript / context window / annotation page / course page for working on contgext windows and rag and all that.

First, we are going to use `go-minitrace help --all` and its SQL + jsverb scripting capabilities to find the design system playbooks we've surly written in pi sessions in this workspace.

Use go-minitrace and store all scripts you use in the scripts/ folder of the new ticket, keeping a detailed and frequent diary, to search for relevant documents that were written while we were working on web/

You can also use git + sqlite + other tools to do richer searches without having to look at a lot of information yourself."

**Assistant interpretation:** Create a docmgr ticket for a new frontend design-integration effort, then use go-minitrace-driven transcript search to recover previous design-system/web contribution playbooks and record every reusable command/script in the ticket.

**Inferred user intent:** Build the next design iteration on top of prior work and documented conventions rather than re-inventing frontend structure, especially for pages around transcripts, context windows, annotations, courses, and RAG workflows.

### What I did
- Ran `docmgr status --summary-only` to confirm the ticket root and vocabulary.
- Ran `go-minitrace help --all > /tmp/go-minitrace-help-all.txt` as the required first go-minitrace discovery step.
- Created ticket `RAGEVAL-CONTEXT-WINDOWS-DESIGN` titled `Transcript, Context Window, Annotation, and Course Page Design Integration`.
- Added design doc `design-doc/01-design-integration-research-and-implementation-plan.md`.
- Added this diary at `reference/01-investigation-diary.md`.
- Added four initial tasks covering transcript search, convention extraction, requirements mapping, and implementation planning.
- Wrote `scripts/01_capture_go_minitrace_help.sh` to capture `go-minitrace help --all` plus focused query/JS help pages into `sources/go-minitrace/`.
- Wrote `scripts/02_convert_relevant_pi_sessions.sh` to convert two relevant Pi session directories into `sources/minitrace/`.
- Ran both scripts successfully.

### Why
- The user explicitly asked for a new ticket, go-minitrace usage, scripts stored in the ticket, and a detailed/frequent diary.
- Capturing help output locally makes later SQL/JS command choices auditable.
- Converting only the two RAG workspace Pi session stores keeps the archive focused before deeper querying.

### What worked
- `docmgr status --summary-only` reported the repo-local ttmp root and existing vocabulary.
- `go-minitrace help --all` worked and produced 45 lines; the focused help pages produced 1,900 lines total.
- `go-minitrace convert pi` converted 18 total sessions across the outer workspace and nested git checkout session stores.
- The ticket now has a reproducible script trail rather than one-off shell history.

### What didn't work
- No blocker yet.
- Observation: `go-minitrace help --all` is only an index-sized output in this installed version, so I also captured focused help pages: `query-commands`, `query-duckdb`, `structured-query-commands`, `js-api-reference`, and `writing-duckdb-queries`.

### What I learned
- This docmgr root already has topics for `design-system`, `frontend-architecture`, `react`, `rag`, and `ui-dsl`, so the new ticket can stay vocabulary-clean.
- There are two relevant Pi session roots for this workspace:
  - `~/.pi/agent/sessions/--home-manuel-workspaces-2026-05-27-rag-evaluation-system--`
  - `~/.pi/agent/sessions/--home-manuel-workspaces-2026-05-27-rag-evaluation-system-2026-05-27--rag-evaluation-system--`

### What was tricky to build
- The current working directory is the outer workspace, while the git/docmgr root is the nested `2026-05-27--rag-evaluation-system` directory. I kept scripts and docmgr commands anchored under the nested repo so ticket paths resolve correctly.
- Prior Pi sessions may have been recorded from either the outer workspace or nested repo checkout, so the conversion script includes both session directories.

### What warrants a second pair of eyes
- Confirm whether Codex transcript stores should also be included after the first Pi-only pass.
- Confirm whether other workspace directories, such as `ttc-design-system`, should be included if the recovered playbooks reference design-system work outside this repo.

### What should be done in the future
- Write and run saved SQL/JS query commands against the converted archive to rank candidate sessions and extract concrete file/doc references.
- Relate the scripts, help captures, and minitrace outputs to the ticket via docmgr once the first query pass identifies the most important artifacts.

### Code review instructions
- Start with `scripts/01_capture_go_minitrace_help.sh` and `scripts/02_convert_relevant_pi_sessions.sh`.
- Re-run them from the repo root with:
  - `ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/scripts/01_capture_go_minitrace_help.sh`
  - `ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/scripts/02_convert_relevant_pi_sessions.sh`

### Technical details
- Ticket path: `ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration`
- Converted archive list: `sources/minitrace/archive-files.txt`
- Help capture directory: `sources/go-minitrace/`

## Step 2: SQL/JS Transcript Search and Recovered Design Sources

I moved from conversion into actual evidence extraction. The goal of this step was to use go-minitrace's SQL and JS command support to rank prior sessions and recover concrete documents, not to manually skim long transcripts.

The search found the expected prior work: an original RAG React design-system guideline/playbook from June 1, a follow-up audit, a current design reference pair from June 7, and Widget DSL packaging/visual-quality guides from June 4–5. I also used git history to find the strongest commit pointer for the original web design-system playbook.

### Prompt Context

**User prompt (verbatim):** (same as Step 1)

**Assistant interpretation:** Use minitrace SQL/JS and other repository tools to find previously written frontend/design-system guidance and summarize the useful sources for the new ticket.

**Inferred user intent:** Recover prior contribution rules for `web/` and design-system work so the new transcript/context-window/annotation/course design can follow established conventions.

### What I did
- Wrote `scripts/03_candidate_session_search.sql`, a DuckDB query that unnests turns/tool calls and ranks sessions by `design system`, `playbook`, `web/`, `react`, `atoms`, `molecules`, `foundation`, `storybook`, `context window`, `transcript`, `annotation`, and `course` terms.
- Ran the SQL query via:
  - `go-minitrace query duckdb --archive-glob "$TICKET/sources/minitrace/*/active/*/*.minitrace.json" --sql-file "$TICKET/scripts/03_candidate_session_search.sql" --output json > "$TICKET/sources/candidate-session-search.json"`
- Wrote a JS query command at `scripts/query-commands/design/web-playbook-search.js` using `require("minitrace")` and `mt.query(...)`.
- Ran the JS command via `go-minitrace query commands --query-repository ... design web-playbook-search web-playbook-evidence ...` and saved `sources/web-playbook-evidence.json`.
- Wrote `scripts/04_git_frontend_history.sh` to capture git history and tracked frontend files for `web`, `packages/rag-evaluation-site`, `internal/web`, `pkg/widgetdsl`, `pkg/widgetserver`, and `pkg/widgetrunner`.
- Wrote `scripts/05_catalog_recovered_design_docs.py` to build `sources/01-recovered-design-docs.md`, a compact catalog of recovered docs and first headings.
- Updated the main design document with initial recovered guidance and an implementation plan.

### Why
- SQL gives a high-recall pass over session turns/tool outputs without manual transcript reading.
- The JS query command demonstrates the requested jsverb-style path and gives a place to evolve richer ranking logic.
- Git history catches committed design docs even when docmgr indexing misses or vocabulary/frontmatter is imperfect.

### What worked
- The SQL search ranked session `019e8afc-89d8-7f43-95d0-e7f6be234ee4` highest for combined design-system/frontend/playbook evidence.
- The minitrace evidence pointed directly to recent design-reference tickets and the Obsidian article path.
- Git history identified commit `1cb51dd` (`docs: audit RAG design system guidelines`) as the direct commit that added the RAG design-system guideline audit and related docs.
- The catalog script recovered these core documents:
  - `RAG React Design System Guidelines`
  - `RAG Design System Guideline Audit`
  - `Styling and Design Reference`
  - `Widget Hierarchy and Interaction Reference`
  - `Widget DSL Visual Quality Analysis and Implementation Guide`
  - `WidgetRenderer Packaging Architecture and Implementation Guide`
  - `UI DSL and Kanban DSL Design and Implementation Guide`
  - `RAG Widget DSL Design — Component-to-HTML Mapping`
  - `Review and Revised Implementation Guide for the RAG Widget DSL`

### What didn't work
- First JS command run failed with a DuckDB regexp error:
  - `Error: GoError: executing js query: Invalid Input Error: missing ): docmgr doc add|docmgr doc relate|write(|design-doc|reference|.md`
- I fixed it by replacing the problematic regex fragment with a bracket expression: `write[(]` and `[.]md`.
- The first run of `scripts/04_git_frontend_history.sh` failed with:
  - `fatal: not a git repository (or any of the parent directories): .git`
- Root cause: I computed the repo root with too many `..` segments from the ticket directory. I fixed `REPO_ROOT` from `../../../../../../..` to `../../../../..`.
- The first run of `scripts/05_catalog_recovered_design_docs.py` marked repo docs as missing because `REPO = parents[5]` resolved to `ttmp` rather than the repo root. I fixed it to `parents[6]`.

### What I learned
- The actual frontend contribution handbook already exists as `ttmp/2026/06/01/.../design-doc/02-rag-react-design-system-guidelines.md`; it just is not linked from `packages/rag-evaluation-site/README.md` or a top-level frontend `CONTRIBUTING.md`.
- The current design reference is split into styling/design and widget hierarchy/interaction documents under `RAGEVAL-DS-RAGEVAL-QWEN36-HIGH`.
- The older `web/` design-system playbook and the newer `packages/rag-evaluation-site` Widget IR package docs need to be reconciled before implementation starts.

### What was tricky to build
- The minitrace query had to inspect both turn content and tool-call inputs/outputs because the important evidence appears in prompts, commands, generated file paths, and command output.
- Querying the current live session can pollute results because the current ticket creation also contains many target terms. I kept the output but treated it as low-value self-evidence, while relying on older high-score sessions and git commits for actual recovered sources.
- The generated ticket path is deeply nested, so scripts that infer repo root must be careful and should be re-run after edits.

### What warrants a second pair of eyes
- The recovered design guidelines were written for `web/`, while the latest PR work introduced `packages/rag-evaluation-site`; we need a design decision on whether the new pages land in legacy `web/`, the package, or both.
- The catalog currently includes one Obsidian article referenced by transcripts. There may be additional external vault notes with less obvious filenames.

### What should be done in the future
- Read the recovered docs in detail and distill a short `CONTRIBUTING.md`/frontend playbook for future implementers.
- Use the recovered rules to decide exact component placement for transcript/context-window/annotation/course surfaces.
- Optionally add Codex transcript search if Pi transcripts do not cover all design iterations.

### Code review instructions
- Review `scripts/03_candidate_session_search.sql` for ranking terms and false positives.
- Review `scripts/query-commands/design/web-playbook-search.js` as the reusable JS command entrypoint.
- Review `sources/01-recovered-design-docs.md` to confirm the recovered source list is complete enough for planning.
- Validate by re-running scripts 03–05 after scripts 01–02 have created the minitrace archive.

### Technical details
- Candidate session search output: `sources/candidate-session-search.json`.
- JS command output: `sources/web-playbook-evidence.json`.
- Git history output: `sources/git/frontend-history.tsv` and `sources/git/frontend-history-with-files.txt`.
- Recovered source catalog: `sources/01-recovered-design-docs.md`.

## Step 3: Research Logbook for Future Minitrace Searches

I added a dedicated research logbook to capture which resources were useful, stale, incomplete, or need updates for the minitrace/session-log recovery process. This document is intentionally broader than the design-integration plan: it records the meta-process so the next transcript-mining task starts from tested resources and known failure modes.

The logbook covers go-minitrace skills/help pages, ticket-local scripts, generated evidence files, recovered design-system documents, the Obsidian article surfaced by transcripts, and the package README gap. It also records concrete update recommendations such as excluding current-session self-hits, generalizing repo-root discovery, and adding a worked example for “find prior project knowledge in local sessions.”

### Prompt Context

**User prompt (verbatim):** "Create a research logbook document in the ticket that allows us to keep track of which resources    
 are useful / out of date / need updating .                                                          
                                                                                                     
 For all the documents and external resources you read, list:                                        
                                                                                                     
 - what you were researching                                                                         
 - what you were looking for in this document in particular                                          
 - why you chose it
 - how you found the resource itself
 - what you found useful in the document                                                             
 - what you didn't find useful                                                                       
 - what is out of date / what was wrong                                                              
 - what would need updating.                                                                         
                                                                                                     
 Store in ticket and upload to remarkable.     

---

^^ this is regarding th eminitrace process, so that next time we have to look for information in session logs we have a better stsarting point."

**Assistant interpretation:** Add a ticket reference document that evaluates every resource consulted during the minitrace/session-log search process, then upload that logbook to reMarkable.

**Inferred user intent:** Build reusable process memory so future agents can mine session logs faster and know which go-minitrace docs, scripts, and recovered resources are reliable or stale.

### What I did
- Created `reference/02-research-logbook.md` with `docmgr doc add`.
- Wrote detailed entries for go-minitrace skills/help pages, ticket scripts, generated evidence outputs, recovered frontend/design-system docs, the Obsidian article, and the incomplete package README.
- Included the exact evaluation fields requested by the user for each resource.
- Related the logbook to the core ticket scripts and recovered-doc catalog.
- Updated the changelog for the new logbook.

### Why
- The previous design-integration doc explains what we found; this logbook explains how trustworthy each research input was and what should be updated before the process is reused.
- Future transcript searches need known-good starting points and known traps such as current-session self-hits.

### What worked
- `docmgr doc add` created the logbook at `reference/02-research-logbook.md`.
- The logbook now has a quick-reference status table plus detailed per-resource notes.
- `docmgr doc relate` successfully related the document after frontmatter was fixed.

### What didn't work
- The first `docmgr doc relate` failed because I used object entries under `ExternalSources`:
  - `line 25: cannot unmarshal !!map into string`
- I fixed this by changing `ExternalSources` back to `[]` and keeping the Obsidian path in the body.

### What I learned
- In this docmgr frontmatter schema, `ExternalSources` is not the same shape as `RelatedFiles`; map entries with `Path`/`Note` are invalid there.
- The most valuable future improvement is not another manual note: it is making the minitrace scripts exclude the current session and auto-catalog candidate docs.

### What was tricky to build
- The logbook needs to be useful for process reuse, not just a bibliography. I kept each resource entry structured around the user’s requested evaluation questions and added a summary table for quick scanning.
- Some resources are both useful and stale. The older component-to-HTML DSL mapping, for example, remains useful as a widget catalog but is architecturally superseded by Widget IR rendered by React.

### What warrants a second pair of eyes
- The list of resources is based on what was read during this session and the previous recovery pass; there may be additional Codex or Obsidian resources not yet included.
- The “out of date” judgments should be checked against current code before turning them into doc updates.

### What should be done in the future
- Add current-session exclusion to `03_candidate_session_search.sql` and `web-playbook-search.js`.
- Generalize `05_catalog_recovered_design_docs.py` so candidate docs are discovered from query/git output.
- Create a concise minitrace “prior project knowledge recovery” playbook using this logbook as source material.

### Code review instructions
- Start with `reference/02-research-logbook.md`.
- Check that each resource entry answers the user’s requested questions.
- Run `docmgr doctor --ticket RAGEVAL-CONTEXT-WINDOWS-DESIGN --stale-after 30` after edits.

### Technical details
- Logbook path: `ttmp/2026/06/07/RAGEVAL-CONTEXT-WINDOWS-DESIGN--transcript-context-window-annotation-and-course-page-design-integration/reference/02-research-logbook.md`.

## Step 3: Source Management and Context Viewer Analysis

I copied the recovered design-system and Widget DSL documents into this ticket's `sources/` directory and unpacked the new `context-viewer.zip` design iteration into a dedicated source subdirectory. After the user clarified that this pass should be analysis-first rather than an integration plan, I shaped the new guide around understanding the updated work, mapping its concepts, and identifying likely integration surfaces without prescribing a detailed implementation order.

The resulting analysis treats the context-viewer code as a prototype and product/design artifact: it has valuable screens, interaction states, typography adjustments, diagram vocabulary, and annotation patterns, but it should not be copied directly as production React because it uses global scripts, inline styles, CDN React/Babel, and `window` exports.

### Prompt Context

**User prompt (verbatim):** "cp them into the sources/ directory of the new ticket and make an index.md in that directory to explain what they are, how they were found, what they are useful for. That way we capture what we are starting with.

Then, unzip ~/Downloads/context-viewer.zip into a subdirectory in sources/ as well. This is the updated version of our design system, with new widgets and potentially states and such. This is the stuff we want to add to the existing design system, while keeping the rest working. I think there are slight stylistic adjustments in the typography which I personally like. But anyway. Manage the sources, and then Create  a detailed analysis / design / implementation guide that is for a new intern, explaining all the parts of the system needed to understand what it is, with prose paragraphs and bullet point sand pseudocode and diagrams and api references and file references. It should be very clear and technical. Store in the ticket and the nupload to remarkable."

**Assistant interpretation:** Preserve the recovered baseline docs and new prototype zip as local ticket sources, write an intern-friendly technical analysis guide focused on understanding the updated design iteration and integration concepts, and upload the resulting guide bundle to reMarkable.

**Inferred user intent:** Make the starting point durable and reviewable before coding, so the future integration of the context-viewer design system can preserve existing behavior while adopting the new typography, widgets, states, and product surfaces.

### What I did
- Copied the most relevant recovered docs into `sources/02-recovered-documents/`.
- Unpacked `~/Downloads/context-viewer.zip` into `sources/03-context-viewer-design-iteration/`.
- Copied the original zip into `sources/03-context-viewer-design-iteration/context-viewer.zip` for provenance.
- Wrote `sources/index.md` explaining the source corpus, how the docs were found, why each copied document matters, and what each context-viewer file contains.
- Generated file lists and checksums:
  - `sources/02-recovered-documents/00-file-list.txt`
  - `sources/03-context-viewer-file-list.txt`
  - `sources/03-context-viewer-sha256.txt`
- Wrote `scripts/06_analyze_context_viewer_sources.py` to extract a compact source map from the prototype.
- Generated `sources/04-context-viewer-source-map.md` and `.json`.
- Added `design-doc/02-context-viewer-design-iteration-analysis.md`, an analysis-first intern guide with source references, DTO sketches, pseudocode, API references, and integration concept maps.
- Related the analysis doc to the source index, source map, key prototype files, `packages/rag-evaluation-site/src/widgets/ir.ts`, and `pkg/widgetdsl/module.go`.
- Checked tasks 3 and 4.

### Why
- The user asked to capture the starting source material before integration.
- The new prototype has many valuable ideas but also many throwaway implementation choices, so an analysis boundary is necessary before implementation planning.
- Copying sources locally avoids future ambiguity if the original ticket docs or `~/Downloads/context-viewer.zip` move or change.

### What worked
- The source bundle is now self-contained under the ticket.
- The context-viewer zip unpacked cleanly into 17 original files plus the copied zip provenance artifact.
- The analysis identified the key new subsystems: typography adjustments, grouped studio navigation, context-window diagram language, visualizer/budget inspector, upload parser workflow, annotated transcript rail, anchored comments, course/slides/handout surfaces, and design tweak controls.
- The guide explicitly distinguishes implementation concepts from a future detailed implementation plan.

### What didn't work
- The first checksum command failed because paths with spaces were passed through command substitution:
  - `sha256sum: ./Design: No such file or directory`
  - `sha256sum: System.html: No such file or directory`
- I fixed it by switching to a null-delimited pipeline:
  - `find . -type f -print0 | sort -z | xargs -0 sha256sum > ../03-context-viewer-sha256.txt`
- The first source-map script over-collected CSS classes from JSX property access because it applied a generic `.class` regex to all files. I rewrote it to extract functions/constants from JSX and CSS variables/classes only from `.css` files.
- `docmgr doctor` reports warnings because the user explicitly requested `sources/index.md`; docmgr prefers only one `index.md` per ticket and numeric prefixes for docs under subdirectories.

### What I learned
- The prototype's strongest reusable subsystem is the context-window diagram language (`patterns.jsx` + `diagrams.jsx`), not the duplicated primitive components.
- The typography change is systematic: body/compact/metric/code roles are slightly larger and more readable while staying dense.
- The transcript and comments screens provide concrete annotation interaction states that were not present in the old design-system docs.
- The right integration frame is: analyze prototype → promote concepts into layers → only later expose stable concepts through Widget IR recipes.

### What was tricky to build
- The user requested an `index.md` under `sources/`, which is useful for humans but conflicts with docmgr's single-index hygiene rule. I kept the requested filename and will call out the warning in handoff.
- The guide had to avoid becoming a detailed implementation plan while still being useful to an intern. I handled this by writing API sketches, DTOs, component landing-zone maps, and an “integration concept map” rather than ordered coding phases.

### What warrants a second pair of eyes
- Whether the prototype's typography should become the new default or a page-mode/theme variant.
- Whether `MacWindow` and `MenuBar` are appropriate product chrome for the context studio, given earlier RAG dashboard guidance preferred no window chrome.
- Whether new context-window components should land first in `web/`, `packages/rag-evaluation-site`, or both.

### What should be done in the future
- Capture screenshots of every prototype screen.
- Build a component promotion table from prototype function to existing/new component layer.
- Decide the frontend target boundary before any implementation begins.
- Define canonical DTOs for context windows, transcripts, annotations, comments, slides, and handouts.

### Code review instructions
- Start with `sources/index.md` to understand the source corpus.
- Then read `design-doc/02-context-viewer-design-iteration-analysis.md`.
- Use `sources/04-context-viewer-source-map.md` to navigate the prototype quickly.
- Inspect `sources/03-context-viewer-design-iteration/patterns.jsx`, `diagrams.jsx`, `screens2.jsx`, and `screens3.jsx` for the highest-value new concepts.

### Technical details
- Original zip: `~/Downloads/context-viewer.zip`.
- Unpacked copy: `sources/03-context-viewer-design-iteration/`.
- Source index: `sources/index.md`.
- Analysis doc: `design-doc/02-context-viewer-design-iteration-analysis.md`.

## Step 4: Add Web Frontend README Pointer and Glossary

I added a `web/README.md` so future contributors can find the recovered design-system guidance without needing to know the docmgr ticket path first. The README does not duplicate the full analysis; it gives a compact primer, glossary, contribution checklist, and links to the current source documents.

This keeps the detailed analysis in the ticket while making the active frontend directory self-describing.

### Prompt Context

**User prompt (verbatim):** "you have a nice glossary in the current system primer. And just in general the guidance is good. Add it to web/README.md or have README.md point at the proper resources."

**Assistant interpretation:** Add contributor-facing frontend documentation in `web/README.md` that includes the glossary and points to the detailed ticket resources.

**Inferred user intent:** Make the design-system/contribution guidance discoverable from the web frontend directory before anyone starts editing components.

### What I did
- Created `web/README.md`.
- Included commands, architecture primer, glossary, contribution checklist, and links to current recovered references.
- Linked the context-viewer analysis and source index from the README.

### Why
- The guidance was useful but buried in a ticket design doc.
- `web/` had no README, so contributors had no local entrypoint for design-system conventions.

### What worked
- The README now points at the current June 7 design references and context-viewer source bundle.
- The glossary captures both existing frontend layers and new context-window/annotation terminology.

### What didn't work
- N/A

### What I learned
- The best short-term path is a discoverability README rather than moving all long-form design docs into `web/`.

### What was tricky to build
- The README needed to be useful on its own while avoiding drift from the ticket docs. I kept the detailed content in links and only duplicated compact glossary/checklist material.

### What warrants a second pair of eyes
- Confirm whether the `web/README.md` should also link to `packages/rag-evaluation-site/README.md` once package-specific contribution guidance is added.

### What should be done in the future
- Add a similar short README or CONTRIBUTING file under `packages/rag-evaluation-site` if that package becomes the integration target.

### Code review instructions
- Start with `web/README.md`.
- Check that all relative links resolve from the `web/` directory.

## Step 5: Rebase Context-Viewer Plan on the Unified Package Boundary

After closing the design-system unification ticket, I reread the context-viewer diary and analysis docs to re-establish the next implementation boundary. The unification work resolves the biggest open question from the earlier context-viewer analysis: reusable context-window, transcript, annotation, and course presentation vocabulary should start in `packages/rag-evaluation-site`, while `web` should compose package components into app pages and own API/container behavior.

I wrote a new implementation plan that turns the earlier analysis-only concept map into a phase-by-phase coding path. The first recommended coding step is deliberately small: define JSON-compatible package DTOs and fixtures before porting diagram renderers or wiring pages.

### Prompt Context

**User prompt (verbatim):** "ok, close the ticket for the unification, and continue with the context-viewer. What are the next steps from here, potentially reread the diaries if you need to"

**Assistant interpretation:** Confirm the unification ticket is closed, reread the context-viewer diary/plan if needed, and resume by identifying concrete next steps for implementation.

**Inferred user intent:** Move out of design-system cleanup mode and back into context-viewer integration with a clear, current plan that reflects the new package/web ownership boundary.

### What I did
- Confirmed `RAGEVAL-DESIGN-SYSTEM-UNIFY` is already `complete` with 0 open tasks.
- Reread the context-viewer investigation diary and the two existing design docs.
- Created `design-doc/03-context-viewer-integration-plan-after-design-system-unification.md`.
- Updated the plan to reflect that `packages/rag-evaluation-site` is now the canonical reusable layer.
- Completed the open task for writing the phased integration plan.

### Why
- The previous analysis intentionally avoided a detailed implementation plan because the package/web boundary was still under discussion.
- The design-system unification now gives a stable boundary for implementation decisions.

### What worked
- The new plan identifies a package-first, Storybook-first path.
- It provides a component promotion table, DTO sketches, landing zones, validation commands, and phased implementation order.

### What didn't work
- N/A.

### What I learned
- The next safe step is DTOs and fixtures, not SVG renderers or app navigation. Contracts unlock later diagram, transcript, Storybook, web, and Widget IR work without forcing premature UI decisions.

### What was tricky to build
- The plan needed to preserve the earlier analysis-first conclusions while updating the target boundary after unification. The critical distinction is that reusable visuals belong in the package, but upload/parsing workflows, routing, API calls, and persistence stay in `web` until generalized.

### What warrants a second pair of eyes
- Whether readable typography should become global package defaults or a scoped transcript/course mode.
- Whether course/slides/handout surfaces should be package demos, web product pages, or both.

### What should be done in the future
- Start Phase 1 by adding context DTOs and fixtures under `packages/rag-evaluation-site/src/context`.

### Code review instructions
- Start with `design-doc/03-context-viewer-integration-plan-after-design-system-unification.md`.
- Compare its target-boundary table against the recently closed design-system unification ticket.

### Technical details
- First coding target: `packages/rag-evaluation-site/src/context/types.ts`, `fixtures.ts`, and `index.ts`.

## Step 6: Add Package Context DTOs and Fixtures

I started the first coding phase from the post-unification plan by adding JSON-compatible context-viewer contracts and normalized fixture data to the package. This creates a stable contract layer before any SVG diagram, transcript card, annotation rail, or course page component is ported from the prototype.

The fixtures are intentionally derived from the prototype's sample data but normalized for production React: every context part, transcript message, slide, agenda item, and annotation has a stable `id`, and shared metadata is represented with JSON-compatible records so the same contracts can later flow through Widget IR or Goja recipes.

### Prompt Context

**User prompt (verbatim):** (same as Step 5)

**Assistant interpretation:** Continue from planning into the first safe package implementation step.

**Inferred user intent:** Begin integrating context-viewer concepts without breaking the design-system boundary that was just cleaned up.

### What I did
- Added `packages/rag-evaluation-site/src/context/types.ts` with context-window, transcript, annotation, course, slide, and handout DTOs.
- Added `packages/rag-evaluation-site/src/context/fixtures.ts` with normalized fixtures from `data.jsx` and `data2.jsx`.
- Added `packages/rag-evaluation-site/src/context/index.ts`.
- Exported the new context module from `packages/rag-evaluation-site/src/index.ts`.
- Ran `pnpm --dir packages/rag-evaluation-site typecheck`.
- Ran `pnpm --dir packages/rag-evaluation-site build`.

### Why
- DTOs and fixtures are the lowest-risk first increment: they let later Storybook components share a stable vocabulary without prematurely committing to rendering details.
- Keeping the contracts JSON-compatible preserves the future Widget IR/Goja path.

### What worked
- Package typecheck passed.
- Package build passed.

### What didn't work
- First typecheck failed because `src/index.ts` re-exported both Widget IR JSON types and the new context JSON types under the same names:
  - `src/index.ts(6,1): error TS2308: Module './widgets' has already exported a member named 'JsonPrimitive'. Consider explicitly re-exporting to resolve the ambiguity.`
  - `src/index.ts(6,1): error TS2308: Module './widgets' has already exported a member named 'JsonValue'. Consider explicitly re-exporting to resolve the ambiguity.`
- I fixed this by renaming the context-side types to `ContextJsonPrimitive`, `ContextJsonValue`, and `ContextJsonRecord`.

### What I learned
- The package top-level export surface already exposes Widget IR JSON types, so new domain contract modules should avoid generic type names.
- The prototype's kind vocabulary includes existing terms such as `context`, `summary`, `generated`, `evicted`, and `empty`; the production type should preserve those terms rather than forcing an artificial smaller taxonomy.

### What was tricky to build
- The fixtures needed to stay faithful to the prototype while becoming deterministic enough for stories and future visual tests. Stable IDs were added everywhere even when the prototype used array position or inline annotation objects.

### What warrants a second pair of eyes
- Review `ContextPartKind` naming before diagram components rely on it.
- Confirm whether `empty` should remain a first-class context part kind or be derived from `limit - total` in renderers.

### What should be done in the future
- Build package context-window diagram components and stories on top of these DTOs.

### Code review instructions
- Start with `packages/rag-evaluation-site/src/context/types.ts`.
- Then inspect `packages/rag-evaluation-site/src/context/fixtures.ts` for faithful normalization from the prototype.
- Validate with `pnpm --dir packages/rag-evaluation-site typecheck` and `pnpm --dir packages/rag-evaluation-site build`.

### Technical details
- Source prototype files: `sources/03-context-viewer-design-iteration/data.jsx` and `data2.jsx`.
- Public package export: `packages/rag-evaluation-site/src/index.ts` now exports `./context`.

## Step 7: Reconfirm Atomic Split from Original Sources

I reread the first design document and the source index after the user asked whether the context-viewer work was still being split through atoms, molecules, organisms, pages, and containers. The answer is yes: the recovered guidance explicitly says layer ownership is the governing rule, and the context-viewer prototype must be promoted into layers rather than copied as one large screen.

I updated the post-unification integration plan with an explicit atomic layer split so the next implementation step is constrained by that recovered guidance. The DTO/fixture step remains valid because it defines contracts only; the next UI step should start with small package atoms/molecules such as `ContextKindSwatch`, `ContextLegend`, and `ContextBudgetBar`, then compose them into organisms.

### Prompt Context

**User prompt (verbatim):** "what about splitting this out into atoms / molecules / organisms, did you read the first design document again to refresh your memory? and potentially the attached sources / where we got the information form?"

**Assistant interpretation:** Re-check the original design evidence and make sure the implementation plan follows the atomic design/layer ownership rules, rather than jumping straight into coarse components.

**Inferred user intent:** Keep the context-viewer integration aligned with the recovered design-system playbooks and source provenance before further coding.

### What I did
- Reread `design-doc/01-design-integration-research-and-implementation-plan.md`, especially the recovered guidance and candidate component decomposition.
- Reread `sources/index.md` to confirm where the recovered docs and context-viewer zip came from.
- Reread `sources/01-recovered-design-docs.md` to refresh the source list and why each recovered document matters.
- Updated `design-doc/03-context-viewer-integration-plan-after-design-system-unification.md` with an explicit atomic layer split.

### Why
- The context-viewer prototype contains screen-level JSX and duplicated primitives, so the implementation needs a clear layer split before more code is ported.
- The first design document already made this point; the new post-unification plan needed to state it as an implementation constraint.

### What worked
- The original source material cleanly supports a package-first but still atomic/layered implementation plan.
- The source index confirms the provenance: prior docs were recovered via go-minitrace SQL/JS plus git history, and the new design iteration came from `~/Downloads/context-viewer.zip`.

### What didn't work
- N/A.

### What I learned
- The next UI step should begin with small package atoms/molecules, not a `ContextWindowVisualizerPanel` organism first.

### What was tricky to build
- The package-first decision from design-system unification does not replace atomic design; it only decides the landing zone. Inside the package, components still need to be split by layer.

### What warrants a second pair of eyes
- Review whether `ContextKindSwatch` is an atom or molecule. It is small like an atom, but it carries domain meaning from the context-window taxonomy.

### What should be done in the future
- Implement the next UI increment in this order: atom marker/swatch → legend/budget molecules → diagram molecules → visualizer organism.

### Code review instructions
- Review the new “Atomic layer split” section in `design-doc/03-context-viewer-integration-plan-after-design-system-unification.md`.

### Technical details
- Source provenance docs reread: `sources/index.md` and `sources/01-recovered-design-docs.md`.

## Step 8: Add Detailed Handoff Tasks and Context Kind Taxonomy

I replaced the remaining broad context-viewer implementation tasks with a detailed handoff checklist that tracks the work through package atoms, molecules, organisms, Storybook validation, web pages, and eventual Widget IR/Goja recipes. This makes progress auditable at a component-layer level instead of leaving a single vague “build diagrams” task.

I then completed the first implementation task by adding a typed context-kind taxonomy and visual style helpers. These helpers capture the prototype's `patterns.jsx` vocabulary in a React/package-safe form without yet adding rendering components.

### Prompt Context

**User prompt (verbatim):** "ok, add detailed tasks to the ticket so that we can exactly keep track of our progress and potentially hand off the work to somoeone else. then work on them one by one, committing at appropriate intervals, and feel free to ask me for feedback so i can verify the looks in the storybook when you think it's appropriate."

**Assistant interpretation:** Expand the ticket into granular implementation tasks, then start working through them sequentially with commits at sensible boundaries and ask for visual feedback when Storybook output is ready.

**Inferred user intent:** Make the context-viewer implementation handoff-safe and reviewable, especially for visual components.

### What I did
- Removed broad open tasks for diagrams/transcript/web/Widget IR and added detailed layer-by-layer tasks.
- Added `packages/rag-evaluation-site/src/context/kinds.ts`.
- Exported the kind helpers from `packages/rag-evaluation-site/src/context/index.ts`.
- Ran `pnpm --dir packages/rag-evaluation-site typecheck`.

### Why
- The prototype's context kind map is a shared dependency for swatches, legends, budget bars, diagrams, transcript annotations, and future Widget IR recipes.
- Capturing it as pure typed helpers avoids copying prototype globals or JSX before component APIs are designed.

### What worked
- Package typecheck passed.
- The taxonomy preserves the prototype vocabulary while adding production kinds needed by the DTOs: `instruction`, `conversation`, `retrieval`, `tool`, `annotation`, `course`, and `other`.

### What didn't work
- N/A.

### What I learned
- The production taxonomy needs to be slightly broader than the prototype's `KIND` object because the normalized DTOs distinguish transcript/course concepts that the diagram prototype collapsed into broader buckets.

### What was tricky to build
- The style helpers need to support pattern, tone, and outline modes without depending on SVG IDs yet. I represented patterns as semantic names and CSS-variable-like fills so rendering components can decide whether to use SVG patterns, CSS backgrounds, or both.

### What warrants a second pair of eyes
- Review the exact `ContextPartKind` names and labels before they become visible in diagrams and Widget IR.

### What should be done in the future
- Implement `ContextKindSwatch` as the first visual atom using these helpers.

### Code review instructions
- Start with `packages/rag-evaluation-site/src/context/kinds.ts`.
- Validate with `pnpm --dir packages/rag-evaluation-site typecheck`.

### Technical details
- Prototype source: `sources/03-context-viewer-design-iteration/patterns.jsx`.

## Step 9: Add ContextKindSwatch Atom and Stories

I implemented the first visual atom for the context-viewer integration: `ContextKindSwatch`. It renders the context-kind taxonomy as small Classic Mac-style bitmap/pattern markers in pattern, tone, and outline modes, and it is exported from the package atom barrel for reuse by upcoming legends, diagrams, and annotation components.

This is intentionally a tiny visual primitive. It does not know about diagram layout, token budgets, transcripts, or selection state beyond a simple selected outline; those responsibilities belong to molecules and organisms in later tasks.

### Prompt Context

**User prompt (verbatim):** (same as Step 8)

**Assistant interpretation:** Continue working through the detailed ticket tasks one at a time and commit at a sensible component boundary.

**Inferred user intent:** Build the context-viewer components incrementally with handoff-safe task tracking and Storybook validation.

### What I did
- Added `packages/rag-evaluation-site/src/components/atoms/ContextKindSwatch/ContextKindSwatch.tsx`.
- Added CSS Modules for pattern/tone/outline swatch rendering.
- Added Storybook coverage for all context kinds across `pattern`, `tone`, and `outline` modes plus size/selection states.
- Exported `ContextKindSwatch` from `packages/rag-evaluation-site/src/components/atoms/index.ts`.
- Ran `pnpm --dir packages/rag-evaluation-site typecheck`.
- Ran `pnpm --dir packages/rag-evaluation-site build`.
- Ran `pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-package-storybook-context-swatch`.

### Why
- `ContextKindSwatch` is the smallest reusable visual unit needed by context legends, diagram keys, annotation labels, and future transcript/course affordances.
- Storybook now exposes the kind vocabulary for quick visual review before larger components depend on it.

### What worked
- Package typecheck passed.
- Package build passed.
- Package Storybook build passed.

### What didn't work
- First typecheck failed because I used an unsupported `align` prop on the existing `Inline` layout component in the story:
  - `Property 'align' does not exist on type 'IntrinsicAttributes & InlineProps'.`
- I fixed the story by using `style={{ alignItems: 'center' }}` on `Inline`.

### What I learned
- The current `Inline` API supports `gap`, `justify`, and `wrap`, but not cross-axis alignment. Stories should use inline style or a small local wrapper when cross-axis alignment is needed.

### What was tricky to build
- CSS cannot directly consume the prototype's SVG `url(#pattern)` fills in an atom without an enclosing SVG definition. I approximated the prototype bitmap styles using low-level CSS gradients for the atom; the later SVG diagram components can still use SVG patterns where appropriate.

### What warrants a second pair of eyes
- Visual review should check whether the CSS pattern approximations are close enough to the prototype's Classic Mac SVG patterns.
- The `instruction`, `conversation`, `retrieval`, `tool`, `annotation`, `course`, and `other` additions need label review before they spread into diagrams.

### What should be done in the future
- Build `ContextLegend` from `ContextKindSwatch` and use that as the first human-review checkpoint in Storybook.

### Code review instructions
- Start with `packages/rag-evaluation-site/src/components/atoms/ContextKindSwatch/ContextKindSwatch.tsx`.
- Then inspect `ContextKindSwatch.module.css` and `ContextKindSwatch.stories.tsx`.
- Validate with package typecheck/build and the Storybook build command above.

### Technical details
- Storybook build output used for validation: `/tmp/rag-package-storybook-context-swatch`.

## Step 10: Add ContextLegend Molecule and Storybook Checkpoint

I built the first context-viewer molecule, `ContextLegend`, by composing `ContextKindSwatch` with kind labels from the context taxonomy. This keeps the component layering explicit: the swatch remains the atom, while the legend owns the repeated labeled data-display pattern.

I also opened the new Storybook story in the running package Storybook and captured a screenshot for a quick sanity check. The legend renders in pattern, tone, and outline modes, and this is a good point for human visual feedback before the budget bar and diagram molecules start depending on the same visual language.

### Prompt Context

**User prompt (verbatim):** (same as Step 8)

**Assistant interpretation:** Continue through the granular implementation tasks, adding Storybook coverage and asking for feedback when useful.

**Inferred user intent:** Keep visual integration reviewable and avoid building too much on top of unreviewed aesthetics.

### What I did
- Added `packages/rag-evaluation-site/src/components/molecules/ContextLegend/ContextLegend.tsx`.
- Added `ContextLegend.module.css`, `index.ts`, and `ContextLegend.stories.tsx`.
- Exported `ContextLegend` from `packages/rag-evaluation-site/src/components/molecules/index.ts`.
- Ran package typecheck, package build, and package Storybook build.
- Opened `http://127.0.0.1:6007/?path=/story/component-library-molecules-contextlegend--modes` and captured `context-legend-modes.png`.

### Why
- The legend is the first reusable molecule that validates whether the context kind taxonomy and swatches are understandable in context.
- It is the natural visual checkpoint before `ContextBudgetBar`, `ContextStripDiagram`, and other diagram molecules reuse the same swatches and labels.

### What worked
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `pnpm --dir packages/rag-evaluation-site build` passed.
- `pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-package-storybook-context-legend` passed.
- The running Storybook story rendered visible legend panels and swatches.

### What didn't work
- N/A.

### What I learned
- The legend is readable but visually dense; human feedback should decide whether swatch/label spacing should be loosened before the same language appears inside diagrams.

### What was tricky to build
- `ContextLegend` needed to stay a molecule, not absorb diagram concerns. It therefore exposes `kinds`, `mode`, `compact`, and `selectedKind`, but does not know about tokens, budgets, or snapshots.

### What warrants a second pair of eyes
- Please review whether the pattern/tone/outline modes are visually distinct enough and whether the labels are the right domain vocabulary.

### What should be done in the future
- Implement `ContextBudgetBar` next if the legend language looks acceptable.

### Code review instructions
- Start with `packages/rag-evaluation-site/src/components/molecules/ContextLegend/ContextLegend.tsx`.
- Open the Storybook story `Component Library / Molecules / ContextLegend / Modes` for visual review.

### Technical details
- Storybook build output used for validation: `/tmp/rag-package-storybook-context-legend`.

## Step 11: Add Budget Bar and Strip Diagram Molecules

After the legend review looked good, I continued with the next two context-window molecules: `ContextBudgetBar` and `ContextStripDiagram`. Both consume the package context-window DTOs and remain molecule-level components: they visualize one reusable data-display pattern but do not own page state, API loading, or multi-view orchestration.

I validated both in Storybook and captured browser screenshots for sanity checks. The budget bar shows under-budget, near-limit, and over-budget states; the strip diagram shows a dense turn-14 context window with labels and a compact legend.

### Prompt Context

**User prompt (verbatim):** "looks great, continue, make a few more"

**Assistant interpretation:** Continue implementing the next few granular context-viewer tasks now that the legend visuals were approved.

**Inferred user intent:** Keep momentum through the package molecule layer while preserving commit/task checkpoints.

### What I did
- Added `packages/rag-evaluation-site/src/components/molecules/ContextBudgetBar/*`.
- Added `packages/rag-evaluation-site/src/components/molecules/ContextStripDiagram/*`.
- Exported both molecules from `packages/rag-evaluation-site/src/components/molecules/index.ts`.
- Added Storybook states for budget states, selected budget segment, dense strip segments, selected strip segment, and limit comparison.
- Ran package typecheck, package build, and package Storybook build.
- Opened Storybook and captured `context-budget-states.png` and `context-strip-dense.png` for visual sanity checks.

### Why
- Budget and strip views are the next smallest diagram molecules after the legend.
- They exercise token totals, selected part IDs, dense labels, overflow state, and the context-kind visual language before more complex stack/treemap renderers are built.

### What worked
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `pnpm --dir packages/rag-evaluation-site build` passed.
- `pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-package-storybook-context-budget-strip` passed.
- Browser screenshots showed rendered, readable components with no blank/loading/broken UI.

### What didn't work
- First `ContextBudgetBar` story typecheck failed because the Storybook meta had no required `snapshot` arg and TypeScript treated indexed fixture access as possibly undefined.
- I fixed it by destructuring fixtures and adding a non-null `snapshot` arg to the story meta.

### What I learned
- Storybook stories for required-prop components should provide a default `args` object even when individual stories use custom render functions.
- The current strip labels work for moderately dense examples; more extreme examples may need a tooltip-first or hover detail strategy.

### What was tricky to build
- Over-budget bars need to communicate both proportional composition and model limit. I kept the segment widths relative to the limit and added a red limit marker/outline for over-budget state.

### What warrants a second pair of eyes
- Visual review should check whether the over-budget state is obvious enough and whether strip labels are too busy for dense contexts.

### What should be done in the future
- Build stack and treemap molecules next, then compose all diagram molecules into a `ContextDiagramPanel` organism.

### Code review instructions
- Review `ContextBudgetBar.tsx` and `ContextStripDiagram.tsx` first.
- Open Storybook stories:
  - `Component Library / Molecules / ContextBudgetBar / Budget States`
  - `Component Library / Molecules / ContextStripDiagram / Dense Segments`

### Technical details
- Storybook build output used for validation: `/tmp/rag-package-storybook-context-budget-strip`.

## Step 12: Complete Context Diagram Molecules and Panel Organism

I finished the remaining context-window diagram views and composed them into the first package organism. `ContextStackDiagram` renders a layered call view, `ContextTreemap` renders proportional token tiles, and `ContextDiagramPanel` wraps strip/stack/budget/treemap views behind compact view-switching buttons.

This completes the first full atom → molecule → organism slice for context windows: `ContextKindSwatch` feeds `ContextLegend` and diagram molecules, and those molecules now compose into `ContextDiagramPanel`. I also ran Storybook and browser sanity checks; the new panels render rather than blanking or loading indefinitely.

### Prompt Context

**User prompt (verbatim):** "go ahead."

**Assistant interpretation:** Finish and commit the current context diagram WIP after clarifying the atom/molecule/organism split.

**Inferred user intent:** Keep progressing through the ticket tasks while preserving task bookkeeping and commit hygiene.

### What I did
- Added `packages/rag-evaluation-site/src/components/molecules/ContextStackDiagram/*`.
- Added `packages/rag-evaluation-site/src/components/molecules/ContextTreemap/*`.
- Added `packages/rag-evaluation-site/src/components/organisms/ContextDiagramPanel/*`.
- Added/updated component barrels for molecules, organisms, and package components.
- Ran `pnpm --dir packages/rag-evaluation-site typecheck`.
- Ran `pnpm --dir packages/rag-evaluation-site build`.
- Ran `pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-package-storybook-context-panel`.
- Opened the `ContextDiagramPanel` Storybook story and captured `context-diagram-panel-starting-views.png`.

### Why
- Stack and treemap complete the prototype's four context-window diagram views.
- The panel organism validates that the views compose under one controlled feature panel before moving on to transcript/annotation work.

### What worked
- Package typecheck passed.
- Package build passed.
- Package Storybook build passed.
- Browser visual sanity check showed rendered panels, buttons, legends, stack diagrams, treemaps, and budget view content.

### What didn't work
- No blocking failures.
- Minor visual note from screenshot review: the compact view buttons could use a stronger selected state, although the selected button does use the package `Button selected` styling.

### What I learned
- The first context diagram slice is now layered correctly: atom (`ContextKindSwatch`) → molecules (`ContextLegend`, `ContextBudgetBar`, `ContextStripDiagram`, `ContextStackDiagram`, `ContextTreemap`) → organism (`ContextDiagramPanel`).

### What was tricky to build
- The treemap is currently a flex-based proportional approximation rather than the prototype's full squarify algorithm. It is adequate for a first storyable molecule, but exact squarified packing can be revisited if visual review asks for it.

### What warrants a second pair of eyes
- Review whether treemap proportionality is visually good enough or should use the prototype squarify algorithm.
- Review whether the view selector needs a more obvious selected state.

### What should be done in the future
- Run a package Storybook visual sweep for the new context diagram stories.
- Ask for human feedback on the complete context diagram slice before transcript/annotation implementation.

### Code review instructions
- Review `ContextStackDiagram.tsx`, `ContextTreemap.tsx`, and `ContextDiagramPanel.tsx`.
- Open Storybook story `Component Library / Organisms / ContextDiagramPanel / Starting Views`.
- Validate with package typecheck, build, and Storybook build.

### Technical details
- Storybook build output used for validation: `/tmp/rag-package-storybook-context-panel`.

## Step 13: Improve Patterned Diagram Label Readability

The first diagram panel review found that labels over striped and dotted fills were hard to read. I updated the diagram CSS so text over patterned backgrounds gets an opaque white backing, while plain white and grey/solid backgrounds retain their original treatment.

This keeps the Classic Mac pattern language but prevents the pattern from fighting with the text. I validated the package typecheck and build after the CSS-only change.

### Prompt Context

**User prompt (verbatim):** "when overlaid on striped background, use a white background on the text. as it;'s hard to read. /tmp/pi-clipboard-b79b407d-b400-4e3f-a799-d5cf5a9d3b86.png"

**Assistant interpretation:** Improve context diagram label contrast specifically where labels sit on striped/patterned fills.

**Inferred user intent:** Preserve the patterned visual vocabulary while making labels readable in Storybook and future app pages.

### What I did
- Updated `ContextStripDiagram.module.css` so strip labels use an opaque white backing.
- Updated `ContextStackDiagram.module.css` so only patterned kind rows add white text backing.
- Updated `ContextTreemap.module.css` so only patterned kind tiles add white text backing.
- Ran `pnpm --dir packages/rag-evaluation-site typecheck`.
- Ran `pnpm --dir packages/rag-evaluation-site build`.

### Why
- Labels over diagonal/stipple fills had low contrast.
- Applying the backing only to patterned kinds avoids reducing legibility on grey/solid contexts.

### What worked
- Package typecheck passed.
- Package build passed.

### What didn't work
- My first attempt put a white backing on all stack/treemap labels, including grey/solid backgrounds. The user pointed out that grey backgrounds were fine before and that white-on-white/white backing was unnecessary there. I narrowed the CSS selectors to patterned kinds only.

### What I learned
- Label backing should be conditional by visual treatment, not global by component.

### What was tricky to build
- The CSS has separate class names per context kind, so the fix had to enumerate pattern-bearing kinds instead of styling `.label` universally.

### What warrants a second pair of eyes
- Confirm whether strip labels should also become conditional later. The strip labels are small and can sit on many patterns, so for now they keep the stronger white backing.

### What should be done in the future
- Include this readability state in visual sweep review before moving to transcript/annotation work.

### Code review instructions
- Review CSS diffs in `ContextStripDiagram.module.css`, `ContextStackDiagram.module.css`, and `ContextTreemap.module.css`.

### Technical details
- The pattern-backed kind selectors include system, instruction, summary, annotation, retrieval, tool, result, and generated.

## Step 14: Capture Context Diagram Storybook Visual Sweep

I added and ran a ticket-local css-visual-diff sweep for the new context Storybook stories. The sweep filters package Storybook stories whose title contains `Context`, self-compares each story, and stores the resulting screenshots/reports under the context-viewer ticket.

The sweep captured 18 context stories and every self-compare reported zero changed pixels. I removed the generated Storybook static directory from the ticket artifacts and kept the story list, Storybook index, summary, and css-visual-diff report outputs.

### Prompt Context

**User prompt (verbatim):** "good. continue"

**Assistant interpretation:** Continue after the label readability fix by advancing the next ticket task, the package Storybook visual sweep.

**Inferred user intent:** Keep progressing through the planned context diagram review workflow.

### What I did
- Added `scripts/07_context_storybook_visual_sweep.sh`.
- Built package Storybook, served it locally, filtered context stories from `index.json`, and ran `css-visual-diff compare` for each context story.
- Generated `sources/visual-parity/context-story-sweep/summary.tsv` and per-story reports.
- Removed generated `storybook-static/` and transient HTTP server log from the committed artifact set.

### Why
- The context diagram slice now needs visual evidence before moving to transcript/annotation work.
- A ticket-local sweep makes the visual review reproducible for handoff.

### What worked
- 18 context stories were captured.
- `summary.tsv` reports `0` changed pixels and `0` changed percent for all captured stories.

### What didn't work
- The sweep emits verbose css-visual-diff logs to stdout; the script is still usable, but future cleanup could redirect compare logs unless debugging is needed.

### What I learned
- Filtering Storybook `index.json` by title containing `Context` captures the current atom, molecule, and organism context slice cleanly.

### What was tricky to build
- The script should keep review artifacts but avoid committing the full generated Storybook site. I removed `storybook-static/` after the run.

### What warrants a second pair of eyes
- Review the generated screenshots for the context diagram stories, especially treemap proportionality and selected-state clarity.

### What should be done in the future
- Ask for human feedback on the context diagram slice before starting transcript/annotation components.

### Code review instructions
- Start with `scripts/07_context_storybook_visual_sweep.sh`.
- Inspect `sources/visual-parity/context-story-sweep/summary.tsv`.
- Open selected `url1_screenshot.png` files under `sources/visual-parity/context-story-sweep/reports/`.

### Technical details
- Sweep output: `sources/visual-parity/context-story-sweep/`.
- Story count: 18.

## Step 15: Start Transcript/Annotation Atoms and Molecules

After the context diagram slice passed visual review and the user approved continuing, I started the transcript/annotation layer with small atoms first: `TranscriptRoleBadge` for message roles/tool names and `AnnotationBadge` for labeled annotation chips. I then composed those atoms into `TranscriptMessageCard` and `AnnotationNoteCard` molecules.

This keeps the same atomic split as the diagram work: atoms own the small markers, molecules own reusable message/note presentation, and later organisms will compose a transcript reader and annotation rail with controlled selection state.

### Prompt Context

**User prompt (verbatim):** (same as Step 14)

**Assistant interpretation:** Continue past the diagram sweep into the next ticket tasks, starting the transcript/annotation slice.

**Inferred user intent:** Keep implementing the package-first component layers in order, with Storybook validation.

### What I did
- Added `TranscriptRoleBadge` atom and stories.
- Added `AnnotationBadge` atom and stories.
- Added `TranscriptMessageCard` molecule and stories.
- Added `AnnotationNoteCard` molecule and stories.
- Exported the new atoms and molecules from package barrels.
- Ran package typecheck, package build, and package Storybook build.
- Opened `TranscriptMessageCard / Roles` in Storybook and captured `transcript-message-card-roles.png`.

### Why
- Transcript/annotation work needs atom-level role and annotation markers before building larger reader/rail organisms.
- The new molecules consume the normalized transcript fixtures added earlier.

### What worked
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `pnpm --dir packages/rag-evaluation-site build` passed.
- `pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-package-storybook-transcript-atoms` passed.
- Browser sanity check showed readable transcript cards and badges with no blank/broken UI.

### What didn't work
- N/A.

### What I learned
- The transcript role badge can also carry tool names (`read_file`, `run_test`) without needing a separate tool-call atom yet.

### What was tricky to build
- `TranscriptMessageCard` must accept both message-local `annotationIds` and a shared annotation list. This keeps the message DTO simple while letting the future reader organism own annotation selection state.

### What warrants a second pair of eyes
- Review whether the transcript role styling is visually distinct enough for long transcripts.
- Review whether `AnnotationBadge` should stay an atom or become a very small molecule because it composes `ContextKindSwatch`.

### What should be done in the future
- Build `TranscriptReaderPanel` and `AnnotationRailPanel` organisms with controlled selection props.

### Code review instructions
- Start with `TranscriptRoleBadge.tsx` and `AnnotationBadge.tsx`.
- Then review `TranscriptMessageCard.tsx` and `AnnotationNoteCard.tsx`.
- Open Storybook story `Component Library / Molecules / TranscriptMessageCard / Roles`.

### Technical details
- Storybook build output used for validation: `/tmp/rag-package-storybook-transcript-atoms`.

## Step 16: Add Transcript Reader and Annotation Rail Organisms

I composed the transcript/annotation molecules into the first transcript organisms. `TranscriptReaderPanel` renders a panel of transcript message cards with controlled annotation selection, while `AnnotationRailPanel` renders the annotation notes as a selectable side rail.

This completes the transcript reader/rail organism task and gives the next work a stable package-level surface for building a combined transcript annotation page.

### Prompt Context

**User prompt (verbatim):** "continue"

**Assistant interpretation:** Continue implementing the next ticket task after transcript/annotation atoms and molecules.

**Inferred user intent:** Keep progressing through the package component layers toward full transcript annotation pages.

### What I did
- Added `packages/rag-evaluation-site/src/components/organisms/TranscriptReaderPanel/*`.
- Added `packages/rag-evaluation-site/src/components/organisms/AnnotationRailPanel/*`.
- Exported both organisms from the package organism barrel.
- Added Storybook stories for annotated transcript, tool-heavy excerpt, and annotation rail selection.
- Ran package typecheck, package build, and package Storybook build.
- Opened the `TranscriptReaderPanel / Annotated Transcript` story and captured `transcript-reader-panel.png`.

### Why
- The package now needs feature-level transcript panels before web pages or Widget IR recipes can consume transcript/annotation concepts.
- Keeping selection controlled via props lets future pages coordinate reader and rail state.

### What worked
- `pnpm --dir packages/rag-evaluation-site typecheck` passed.
- `pnpm --dir packages/rag-evaluation-site build` passed.
- `pnpm --dir packages/rag-evaluation-site exec storybook build --output-dir /tmp/rag-package-storybook-transcript-panels` passed.
- Browser sanity check showed a rendered readable transcript reader panel.

### What didn't work
- N/A.

### What I learned
- The transcript reader is usable as a package organism without needing web routing or backend state.

### What was tricky to build
- Storybook stories with controlled selection need small local state wrappers; the production components remain controlled and stateless with respect to selected annotation.

### What warrants a second pair of eyes
- Review whether the transcript cards should be denser before long real transcripts are loaded.
- Review whether the annotation rail should use a dedicated CSS module instead of inline button reset styles.

### What should be done in the future
- Build anchored comment card/rail components from the prototype comment states.

### Code review instructions
- Review `TranscriptReaderPanel.tsx` and `AnnotationRailPanel.tsx`.
- Open Storybook stories under `Component Library / Organisms / TranscriptReaderPanel` and `AnnotationRailPanel`.

### Technical details
- Storybook build output used for validation: `/tmp/rag-package-storybook-transcript-panels`.
