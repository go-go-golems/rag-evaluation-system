__package__({
  name: "widget-dsl-visual",
  parents: ["widget-dsl"],
  short: "Widget DSL visual evidence capture workflows"
})

async function captureEvidence(values) {
  const cvd = require("css-visual-diff")
  const outDir = values.outDir
  const widgetUrl = values.widgetUrl.replace(/\/$/, "")
  const storybookUrl = values.storybookUrl.replace(/\/$/, "")
  const viewport = { width: values.viewportWidth, height: values.viewportHeight }
  const waitMs = values.waitMs

  const targets = [
    {
      slug: "widget-actions",
      name: "Generated Widget DSL action page",
      url: `${widgetUrl}/pages/actions`,
      root: '[data-rag-page="RagEvaluationSiteApp"]',
      probes: [
        { name: "app-root", selector: '[data-rag-page="RagEvaluationSiteApp"]', props: ["background-color", "color", "font-family", "font-size", "padding", "max-width"] },
        { name: "first-panel", selector: '[data-rag-page="RagEvaluationSiteApp"] section:first-of-type, [data-rag-page="RagEvaluationSiteApp"] article:first-of-type, [data-rag-page="RagEvaluationSiteApp"] > div:first-child', props: ["background-color", "border", "border-radius", "box-shadow", "padding"] },
        { name: "first-table", selector: "table", props: ["border-collapse", "font-size", "background-color", "border-spacing"] },
        { name: "first-button", selector: "button", props: ["height", "padding", "border-radius", "background-color", "color", "font-weight"] },
      ],
      overlay: [
        ["App Root", '[data-rag-page="RagEvaluationSiteApp"]', "#0ea5e9"],
        ["Data Table", "table", "#22c55e"],
        ["Buttons", "button", "#f97316"],
      ],
    },
    {
      slug: "storybook-pipeline-page",
      name: "Original RAG PipelinePage Storybook page",
      url: `${storybookUrl}/iframe.html?id=pages-pipelinepage--populated&viewMode=story`,
      root: "#storybook-root",
      probes: [
        { name: "storybook-root", selector: "#storybook-root", props: ["background-color", "color", "font-family", "font-size", "padding"] },
        { name: "first-panel", selector: "#storybook-root section:first-of-type, #storybook-root article:first-of-type, #storybook-root > div:first-child", props: ["background-color", "border", "border-radius", "box-shadow", "padding"] },
        { name: "first-table", selector: "table", props: ["border-collapse", "font-size", "background-color", "border-spacing"] },
        { name: "first-button", selector: "button", props: ["height", "padding", "border-radius", "background-color", "color", "font-weight"] },
      ],
      overlay: [
        ["Story Root", "#storybook-root", "#0ea5e9"],
        ["Data Table", "table", "#22c55e"],
        ["Buttons", "button", "#f97316"],
      ],
    },
    {
      slug: "storybook-retrieval-results",
      name: "Original RetrievalResultsPanel Storybook organism",
      url: `${storybookUrl}/iframe.html?id=component-library-organisms-retrievalresultspanel--hybrid-results&viewMode=story`,
      root: "#storybook-root",
      probes: [
        { name: "root", selector: "#storybook-root", props: ["background-color", "color", "font-family", "font-size", "padding"] },
        { name: "panel", selector: "#storybook-root section:first-of-type, #storybook-root article:first-of-type, #storybook-root > div:first-child", props: ["background-color", "border", "border-radius", "box-shadow", "padding"] },
        { name: "table", selector: "table", props: ["border-collapse", "font-size", "background-color", "border-spacing"] },
      ],
      overlay: [
        ["Story Root", "#storybook-root", "#0ea5e9"],
        ["Result Table", "table", "#22c55e"],
      ],
    },
    {
      slug: "widget-renderer-storybook",
      name: "WidgetRenderer search workbench Storybook composition",
      url: `${storybookUrl}/iframe.html?id=widget-ir-renderer--search-workbench-composition&viewMode=story`,
      root: "#storybook-root",
      probes: [
        { name: "root", selector: "#storybook-root", props: ["background-color", "color", "font-family", "font-size", "padding"] },
        { name: "first-panel", selector: "#storybook-root section:first-of-type, #storybook-root article:first-of-type, #storybook-root > div:first-child", props: ["background-color", "border", "border-radius", "box-shadow", "padding"] },
        { name: "table", selector: "table", props: ["border-collapse", "font-size", "background-color", "border-spacing"] },
      ],
      overlay: [
        ["Story Root", "#storybook-root", "#0ea5e9"],
        ["Data Table", "table", "#22c55e"],
      ],
    },
  ]

  await cvd.write.json(`${outDir}/targets.json`, targets.map(t => ({ slug: t.slug, name: t.name, url: t.url, root: t.root })))

  const browser = await cvd.browser()
  const rows = []
  try {
    for (const target of targets) {
      const page = await browser.page(target.url, { viewport, waitMs, name: target.slug })
      try {
        await page.css(`
          * { animation: none !important; transition: none !important; }
          html { scroll-behavior: auto !important; }
          body { min-height: 100vh; }
        `)
        const statuses = await page.preflight(target.probes.map(p => ({ name: p.name, selector: p.selector, source: "styles", required: p.name === "app-root" || p.name === "root" || p.name === "storybook-root" })))
        await cvd.write.json(`${outDir}/${target.slug}.preflight.json`, statuses)

        const existingProbes = target.probes.filter((_p, idx) => statuses[idx] && statuses[idx].exists)
        let inspectResult = null
        if (existingProbes.length > 0) {
          inspectResult = await page.inspectAll(existingProbes, { outDir: `${outDir}/artifacts/${target.slug}`, artifacts: "css-json" })
        }

        const overlayBuilder = cvd.overlaySpec()
          .legend(true)
          .screenshot("fullPage")
          .style({
            legend: { position: "bottom-right", background: "rgba(255,255,255,0.92)", color: "#111827" },
            targetDefaults: { borderWidth: 2, contentBackground: "rgba(14,165,233,0.08)", labelColor: "white" },
          })
        for (const [name, selector, color] of target.overlay) {
          const status = await page.locator(selector).status()
          if (status.exists) {
            overlayBuilder.target(cvd.overlayTarget(name).selector(selector).borderColor(color).labelBackground(color))
          }
        }
        const overlay = await page.overlay(overlayBuilder.build()).screenshot(`${outDir}/${target.slug}.overlay.png`)
        rows.push({ slug: target.slug, name: target.name, url: target.url, root: target.root, preflight: statuses, inspectResult, overlay })
      } finally {
        await page.close()
      }
    }
  } finally {
    await browser.close()
  }

  const md = [
    "---",
    "Title: Visual Evidence Capture Summary",
    "DocType: reference",
    "Ticket: WIDGETDSL-VISUAL-QUALITY",
    "Status: active",
    "Intent: short-term",
    "Topics:",
    "  - frontend",
    "  - ui-dsl",
    "  - design-system",
    "---",
    "",
    "# Widget DSL Visual Evidence Capture",
    "",
    `- Widget URL: ${widgetUrl}`,
    `- Storybook URL: ${storybookUrl}`,
    `- Viewport: ${viewport.width}x${viewport.height}`,
    `- Wait: ${waitMs}ms`,
    "",
    "## Captured targets",
    "",
    ...rows.map(row => `- **${row.name}** (${row.slug}) — ${row.url}\n  - overlay: \`${row.slug}.overlay.png\`\n  - preflight: \`${row.slug}.preflight.json\``),
    "",
    "## Notes for analysis",
    "",
    "Use the overlay screenshots for layout-level review and the CSS JSON artifacts for computed style evidence. Missing optional probes are recorded in the preflight files and should guide selector refinement rather than be treated as failures.",
  ].join("\n")
  await cvd.write.markdown(`${outDir}/01-visual-evidence-summary.md`, md)
  return rows.map(row => ({ slug: row.slug, url: row.url, overlay: `${outDir}/${row.slug}.overlay.png` }))
}

__verb__("captureEvidence", {
  parents: ["widget-dsl", "visual"],
  short: "Capture Widget DSL and Storybook visual evidence",
  output: "structured",
  fields: {
    outDir: { argument: true, required: true, help: "Output directory for evidence artifacts" },
    values: { bind: "all" },
    widgetUrl: { type: "string", default: "http://127.0.0.1:18791", help: "Base URL of the generated widget-site" },
    storybookUrl: { type: "string", default: "http://127.0.0.1:6007", help: "Base URL of Storybook or storybook-static" },
    viewportWidth: { type: "int", default: 1440 },
    viewportHeight: { type: "int", default: 1100 },
    waitMs: { type: "int", default: 800 }
  }
})
