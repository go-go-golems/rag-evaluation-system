const widget = require("widget.dsl");
const styleSet = widget.context.palette("Dusty Magenta / Blue", [
	{ id: "system", label: "System", accent: "a" },
	{ id: "retrieval", label: "Retrieval", accent: "b" },
]);
const snapshot = {
	id: "ctx",
	title: "Context budget",
	limit: 16000,
	parts: [
		{ id: "system", label: "System", styleKey: "system", tokens: 1200 },
		{ id: "retrieval", label: "Retrieved docs", styleKey: "retrieval", tokens: 8000 },
	],
};
const page = widget.page("Context budget diagram", (p) =>
	p.section("Diagram", (s) =>
		s.view(
			widget.context.diagram(snapshot, (d) =>
				d.styleSet(styleSet).view("budget").selected("retrieval"),
			),
		),
	),
);
