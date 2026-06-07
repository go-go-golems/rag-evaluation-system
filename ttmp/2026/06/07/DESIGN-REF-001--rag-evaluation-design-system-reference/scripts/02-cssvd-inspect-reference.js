__package__({
  name: "design-reference",
  parents: ["design-ref"],
  short: "Inspect generated RAG design-system reference pages"
})

async function inspectReference(url, outDir) {
  const cvd = require("css-visual-diff")
  const browser = await cvd.browser()
  try {
    const page = await browser.page(url, {
      viewport: { width: 1280, height: 1600 },
      waitMs: 250,
      name: "design-reference",
    })

    const probes = [
      {
        name: "button-primary",
        selector: '[data-rag-atom="Button"]',
        props: ["appearance", "background-color", "color", "border", "padding", "font-family", "font-size"],
        attributes: ["class", "data-rag-atom"],
      },
      {
        name: "panel",
        selector: '[data-rag-layout="Panel"]',
        props: ["background-color", "border", "color", "font-family"],
        attributes: ["class", "data-rag-layout"],
      },
      {
        name: "data-table",
        selector: '[data-rag-component="DataTable"]',
        props: ["border-collapse", "background-color", "border", "font-family"],
        attributes: ["class", "data-rag-component"],
      },
    ]

    const preflight = await page.preflight(probes)
    const missing = preflight.filter((p) => !p.exists).map((p) => p.name)
    if (missing.length) {
      throw new cvd.SelectorError(`missing expected probes: ${missing.join(", ")}`)
    }

    const result = await page.inspectAll(probes, {
      outDir,
      artifacts: "css-json",
    })

    return {
      ok: true,
      outputDir: result.outputDir,
      resultCount: result.results.length,
    }
  } finally {
    await browser.close()
  }
}

__verb__("inspectReference", {
  parents: ["design-ref"],
  short: "Inspect button/panel/table CSS on a generated reference page",
  output: "structured",
  fields: {
    url: { argument: true, required: true, help: "Reference page URL" },
    outDir: { argument: true, required: true, help: "Output directory for CSS evidence" },
  },
})
