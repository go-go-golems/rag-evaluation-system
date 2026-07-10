const widget = require("widget.dsl");
const palette = widget.context.palette({
	palette: "Dusty Magenta / Blue",
	entries: [
		{ id: "a", label: "Alpha", accent: "a" },
		{ id: "b", label: "Beta", accent: "b", solid: true },
	],
});
const snapshot = {
	id: "palette",
	title: "Palette",
	limit: 1000,
	parts: [
		{ id: "a", label: "Alpha", styleKey: "a", tokens: 300 },
		{ id: "b", label: "Beta", styleKey: "b", tokens: 500 },
	],
};
const page = widget.page("Style palette", (p) =>
	p.section("Palette", (s) => s.view(widget.context.diagram(snapshot, (d) => d.styleSet(palette)))),
);
