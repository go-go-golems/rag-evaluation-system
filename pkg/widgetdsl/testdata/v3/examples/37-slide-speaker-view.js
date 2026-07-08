const widget = require("widget.dsl");
const page = widget.page("Slide speaker view", (p) =>
	p.section("Speaker", (s) =>
		s.view(
			widget.course.slideDeck(
				{
					index: 0,
					slides: [
						{ id: "s1", title: "Speaker notes", notes: ["Explain the budget", "Show diagram"] },
					],
					snapshot: { id: "ctx", title: "Context", limit: 2048, parts: [] },
				},
				(d) => d.mode("speaker").visualSide("right"),
			),
		),
	),
);
