const widget = require("widget.dsl");
const article = widget.ui.markdownArticle(`# Widget DSL notes

- Builder callbacks stay serializable.
- React owns rendering.
- Actions remain data.`);
const page = widget.page("Markdown article", (p) => p.section("Article", (s) => s.view(article)));
