const widget = require("widget.dsl");
const page = widget.page("Markdown editor", (p) =>
	p.section("Draft", (s) =>
		s.view(
			widget.cms.markdownEditor(
				`# Draft

Write here...`,
				(e) =>
					e
						.title("Live markdown")
						.placeholder("Start typing")
						.onSubmit(widget.act.server("article.save")),
			),
		),
	),
);
