/**
 * 01-render-design-reference.cjs
 *
 * Renders two self-contained HTML reference pages using React SSR.
 * Output: foundation-and-tokens.html + widgets-and-composition.html
 *
 * Run:  cd web && node scripts/01-render-design-reference.cjs
 */

const { createServer, build } = require("vite");
const { renderToString } = require("react-dom/server");
const React = require("react");
const fs = require("fs");
const path = require("path");

const WEB_ROOT = path.resolve(__dirname, "..");
const OUTPUT_DIR = path.resolve(
	__dirname,
	"../../ttmp/2026/06/07/DESIGN-REF-001--rag-evaluation-design-system-reference/output",
);
const TOKENS_CSS = path.join(WEB_ROOT, "src/styles/tokens.css");
const INDEX_CSS = path.join(WEB_ROOT, "src/index.css");

async function main() {
	console.log("=== Step 1: Start Vite SSR server ===");
	const server = await createServer({
		root: WEB_ROOT,
		server: { middlewareMode: true },
		appType: "custom",
	});
	console.log("Vite server started");

	console.log("=== Step 2: Render foundation reference ===");
	const fMod = await server.ssrLoadModule("/src/reference-pages/FoundationReference.tsx");
	const FoundationPage = fMod.default || fMod.FoundationReference;
	const foundationHtml = renderToString(React.createElement(FoundationPage));

	console.log("=== Step 3: Render widget reference ===");
	const wMod = await server.ssrLoadModule("/src/reference-pages/WidgetReference.tsx");
	const WidgetPage = wMod.default || wMod.WidgetReference;
	const widgetHtml = renderToString(React.createElement(WidgetPage));

	console.log("=== Step 4: Build client CSS ===");
	const clientBuild = await build({
		root: WEB_ROOT,
		build: {
			write: false,
			rollupOptions: {
				input: [
					path.resolve(WEB_ROOT, "src/reference-pages/css-entry.ts"),
					path.resolve(WEB_ROOT, "src/reference-pages/FoundationReference.tsx"),
					path.resolve(WEB_ROOT, "src/reference-pages/WidgetReference.tsx"),
				],
			},
		},
		logLevel: "warn",
	});
	let css = "";
	for (const chunk of clientBuild.output) {
		if (chunk.type === "asset" && chunk.fileName.endsWith(".css")) {
			css += chunk.source;
		}
	}
	// Append the raw tokens CSS (missing from the build since it's a CSS-only import)
	const tokensCss = fs.readFileSync(TOKENS_CSS, "utf-8");
	css += "\n/* === design tokens === */\n" + tokensCss;
	console.log("CSS extracted:", css.length, "chars");

	console.log("=== Step 5: Write self-contained HTML ===");
	fs.mkdirSync(OUTPUT_DIR, { recursive: true });

	fs.writeFileSync(
		path.resolve(OUTPUT_DIR, "foundation-and-tokens.html"),
		wrapInPage("Design System — Foundation & Tokens", css, foundationHtml),
	);
	fs.writeFileSync(
		path.resolve(OUTPUT_DIR, "widgets-and-composition.html"),
		wrapInPage("Design System — Widgets & Composition", css, widgetHtml),
	);
	console.log("Written 2 files");

	await server.close();
	console.log("=== Done ===");
}

function wrapInPage(title, css, body) {
	// Full page shell matching the real app
	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>${title}</title>
<style>
:root {
  /* Token bridge for standalone — ensures component CSS resolves correctly */
  --mac-bg: #FFFFFF; --mac-bg-dark: #000; --mac-stripe: #000;
  --mac-text: #000; --mac-text-dim: #666; --mac-text-inv: #fff;
  --mac-border: #000; --mac-surface: #FFF; --mac-surface-2: #eee;
  --mac-accent: #0000CC; --mac-accent-2: #CC0000;
  --mac-green: #007722; --mac-amber: #AA7700;
  --mac-mono: 'Chicago','Geneva',Monaco,'Courier',monospace;
  --font-body: Geneva, 'Helvetica Neue', Arial, sans-serif;
  --font-mono: Monaco, 'Courier New', monospace;
  --rag-font-role-body: 400 12px/1.45 Geneva, sans-serif;
  --rag-font-role-compact: 400 11px/1.4 Geneva, sans-serif;
  --rag-font-role-metadata: 400 11px/1.35 Monaco, monospace;
  --rag-font-role-label: 700 11px/1.2 Monaco, monospace;
  --rag-font-role-metric: 700 12px/1.25 Monaco, monospace;
  --rag-font-role-code: 400 11px/1.45 Monaco, monospace;
}
${css}

/* Reset & base */
*, *::before, *::after { box-sizing: border-box; }
body {
  font-family: Geneva, 'Helvetica Neue', Arial, sans-serif;
  font-size: 12px; line-height: 1.45; color: #000;
  background: #FFF; margin: 0; padding:0;
}
.ds-ref-page { max-width: 960px; margin: 0 auto; padding: 2rem 1.5rem; }
.ds-ref-section { margin-top: 2.5rem; }
.ds-ref-section-title { border-bottom:1px solid #000; padding-bottom:4px; margin-bottom:1rem; font:700 12px/1.25 Monaco,monospace; }
.ds-ref-label { font:400 11px/1.35 Monaco,monospace; color:#666; margin-bottom:0.25rem; }
.ds-ref-row { display:flex; gap:1rem; flex-wrap:wrap; margin:.5rem 0; align-items:flex-start; }
.ds-ref-item { flex:0 0 auto; min-width:120px; }
.ds-ref-code { font:400 11px/1.35 Monaco,monospace; background:#f5f5f5; padding:1px 4px; }
</style>
</head>
<body>
<div class="ds-ref-page">
${body}
</div>
</body>
</html>`;
}

main().catch((err) => {
	console.error(err);
	process.exit(1);
});
