const widget = require("widget.dsl");
const assets = [
	{ id: "hero", title: "Hero", kind: "image", mime: "image/png", filename: "hero.png" },
];
const events = [
	{
		id: "review",
		title: "Review",
		startISO: "2026-07-08T14:00:00Z",
		endISO: "2026-07-08T15:00:00Z",
		styleKey: "review",
	},
];
const page = widget.page("Complete admin dashboard", (p) =>
	p
		.section("Overview", (s) =>
			s
				.metric("Assets", "1")
				.metric("Events", "1")
				.view(widget.ui.callout({ title: "Ready" }, "All modules are available.")),
		)
		.section("Media", (s) => s.view(widget.cms.mediaLibrary(assets)))
		.section("Calendar", (s) =>
			s.view(widget.time.week(events, (w) => w.range(widget.time.range.week("2026-07-08")))),
		),
);
