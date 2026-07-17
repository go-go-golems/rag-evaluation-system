const widget = require("widget.dsl");

function mainFor(active) {
	if (active === "media") {
		return widget.cms.mediaLibrary(
			[{ id: "asset-1", kind: "image", title: "Course cover", filename: "cover.png" }],
			(m) => m.selection("single").selected(["asset-1"]),
		);
	}
	return widget.cms.articleQueue(
		[
			{
				id: "article-1",
				title: "Typed shell migration",
				slug: "typed-shell-migration",
				status: "draft",
				tags: ["v3"],
				author: "Widget team",
				updatedAt: "2026-07-13T12:00:00Z",
			},
		],
		(q) => q.selected("article-1"),
	);
}

function renderPage(query = {}) {
	const active = query.item === "media" ? "media" : "articles";
	const shell = widget.cms.shell(
		{
			title: "Editorial Studio",
			subtitle: "Typed CMS workspace",
			sections: [
				{
					id: "content",
					label: "Content",
					items: [
						{ id: "articles", label: "Articles" },
						{ id: "media", label: "Media" },
					],
				},
			],
		},
		(s) =>
			s
				.active(active)
				.onNavigate(
					widget.act.navigate("/pages/43-cms-shell-layout", {
						query: { item: widget.bind.context("item.id") },
						omitEmpty: true,
					}),
				)
				.footer(widget.ui.caption("Generated from widget.cms.shell"))
				.main(mainFor(active)),
	);
	return widget.page("CMS shell layout", (p) => p.shell(widget.app.rootOwned()).root(shell));
}

const page = renderPage({});
