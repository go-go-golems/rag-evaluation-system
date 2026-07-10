const widget = require("widget.dsl");

const rows = [
	{ id: "asset-hero", title: "Hero diagram", kind: "image", status: "ready" },
	{ id: "asset-notes", title: "Workshop notes", kind: "markdown", status: "draft" },
];

const schema = widget.data
	.fields("assets", (f) => f.key("id").primary("title").short("kind").status("status"))
	.build();

const table = widget.data
	.collection("assets", rows, (c) =>
		c.schema(schema).table((t) =>
			t
				.actionColumn(
					"preview",
					"Preview",
					"Open",
					widget.cms.intent.previewArticle(widget.bind.context("row.id")),
				)
				.actionColumn(
					"archive",
					"Archive",
					"Archive",
					widget.cms.intent.archiveArticle(widget.bind.context("row.id"), {
						confirm: "Archive asset?",
					}),
				),
		),
	)
	.toNode();

const page = widget.page("Row actions", (p) =>
	p.section("Media rows", (s) =>
		s.caption("Action columns hide low-level Button/Action wiring.").view(table),
	),
);
