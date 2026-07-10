const widget = require("widget.dsl");

const articles = [
	{
		id: "a1",
		title: "Context budgets",
		slug: "context-budgets",
		status: "draft",
		tags: ["context"],
		author: "Manuel",
		updatedAt: "2026-07-08T09:00:00Z",
	},
	{
		id: "a2",
		title: "Retrieval eval",
		slug: "retrieval-eval",
		status: "published",
		tags: ["rag"],
		author: "Manuel",
		updatedAt: "2026-07-07T15:30:00Z",
	},
];

const page = widget.page("CMS article queue", (p) =>
	p.section("Queue", (s) =>
		s.view(
			widget.cms.articleQueue(articles, (q) =>
				q
					.selected("a1")
					.onSelect(widget.cms.intent.selectArticle(widget.bind.context("article.id")))
					.onCreate(widget.cms.intent.createArticle()),
			),
		),
	),
);
