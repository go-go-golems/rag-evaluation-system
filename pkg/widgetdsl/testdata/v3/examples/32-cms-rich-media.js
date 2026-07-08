const widget = require("widget.dsl");
const assets = [
	{ id: "img-1", title: "Diagram", kind: "image", mime: "image/png", filename: "diagram.png" },
	{
		id: "pdf-1",
		title: "Slides",
		kind: "document",
		mime: "application/pdf",
		filename: "slides.pdf",
	},
];
const page = widget.page("CMS rich media", (p) =>
	p.section("Library", (s) =>
		s.view(
			widget.cms.mediaLibrary(assets, (m) =>
				m
					.selection("multi")
					.selected(["img-1"])
					.onOpen(widget.cms.intent.openAsset(widget.bind.context("asset.id"))),
			),
		),
	),
);
