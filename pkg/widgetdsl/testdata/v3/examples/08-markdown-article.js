const widget = require("widget.dsl");
const article = widget.raw.component("MarkdownArticle", {
	source: `# Widget DSL notes

- Builder callbacks stay serializable.
- React owns rendering.
- Actions remain data.`,
});
const page = widget.page("Markdown article", (p) => p.section("Article", (s) => s.view(article)));
