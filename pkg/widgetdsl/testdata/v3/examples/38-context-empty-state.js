const widget = require("widget.dsl");
const snapshot = { id: "empty", title: "Empty context", limit: 4000, parts: [] };
const page = widget.page("Context empty state", (p) =>
	p.section("Empty", (s) =>
		s.view(
			widget.context.diagram(snapshot, (d) =>
				d.empty((ctx, h) => h.caption("No context parts yet.")),
			),
		),
	),
);
