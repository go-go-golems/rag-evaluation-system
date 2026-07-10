const widget = require("widget.dsl");

const snapshot = {
	id: "ctx",
	title: "Speaker context",
	limit: 2048,
	parts: [
		{
			id: "system",
			label: "System prompt",
			styleKey: "system",
			tokens: 320,
			contentPreview: "Course facilitation guardrails and output format.",
		},
		{
			id: "retrieval",
			label: "Retrieved notes",
			styleKey: "retrieval",
			tokens: 960,
			contentPreview: "Context budget examples and reclaim-policy notes.",
		},
		{
			id: "speaker",
			label: "Speaker notes",
			styleKey: "assistant",
			tokens: 420,
			contentPreview: "Explain the budget and show the diagram.",
		},
	],
};

const page = widget.page("Slide speaker view", (p) =>
	p.section("Speaker", (s) =>
		s.view(
			widget.course.slideDeck(
				{
					index: 0,
					slides: [
						{
							id: "s1",
							title: "Speaker notes",
							view: "stack",
							notes: ["Explain the budget", "Show diagram"],
						},
					],
					snapshot,
				},
				(d) => d.mode("speaker").visualSide("right"),
			),
		),
	),
);
