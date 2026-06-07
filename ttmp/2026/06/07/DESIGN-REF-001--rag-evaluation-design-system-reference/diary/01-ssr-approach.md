# Diary

## Goal

Produce two self-contained HTML reference pages for the RAG Evaluation Design System that render actual React components (not just tables describing them). Page 1: Foundation & Tokens. Page 2: Widgets & Composition. The approach: write a Node script that uses `react-dom/server` to render every widget variant to static HTML, then bundle the CSS inline into a single file.

## Step 1: Explore the build path for SSR-to-HTML

**User prompt (verbatim):** "not but I want to render actual self contained examples of widgets and all that. It's probably best to go through react and capture the output of a render and then export that, no?"

**Assistant interpretation:** The user wants actual rendered React components in the HTML reference, not descriptions. The approach is SSR: render React components to static HTML strings and inline the CSS.

**Inferred user intent:** Produce a visual reference an external designer can open in any browser and see real rendered components.

### What I did
- Explored the component tree, CSS module structure, and Vite/React setup
- Identified that `react-dom/server` + Vite SSR build is the right path
- Need to handle CSS Modules (they generate scoped class names at build time)

### What worked
- Confirmed all components export from `src/components/index.ts`
- Confirmed `react-dom` is installed (v19), `renderToString` is available
- Confirmed design tokens are in `src/styles/tokens.css` and `src/index.css`

### What didn't work
- (nothing yet — just starting)

### What was tricky to build
- CSS Modules produce scoped class names (e.g. `_header_1a2b3`). To get the real CSS into a static HTML file, we need to either:
  - (a) build the Vite project and extract the CSS from the build output
  - (b) use Vite's SSR mode to resolve CSS modules at render time
  Option (a) is more reliable for a static output.

### What warrants a second pair of eyes
- Whether `renderToString` handles CSS module class resolution correctly
- Whether we need to worry about React 19's SSR streaming API

### What should be done in the future
- Automate this as a CI step
- Add visual diff testing against the generated reference

### Code review instructions
- Check `scripts/01-render-design-reference.mjs`
- Run: `cd web && node ../ttmp/.../scripts/01-render-design-reference.mjs`

### Technical details
- React 19 + ReactDOM/server
- CSS Modules via Vite
- Output: two self-contained HTML files
