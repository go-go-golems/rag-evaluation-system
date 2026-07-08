const widget = require("widget.dsl");

const styleSet = widget.context.styleSet((s) =>
	s
		.style("prompt", { fill: "#f7e8ff", stroke: "#9d4edd" })
		.style("answer", { fill: "#e0f2fe", stroke: "#0284c7" })
		.legend("prompt", "Prompt")
		.legend("answer", "Answer"),
);

const poll = {
	title: "Workshop availability",
	options: [
		{
			id: "wed-9",
			label: "Wed 9",
			startISO: "2026-07-08T09:00:00Z",
			endISO: "2026-07-08T10:00:00Z",
		},
	],
	responses: [{ id: "ana", name: "Ana", availability: { "wed-9": "available" } }],
};

const page = widget.page({ id: "all-modules", title: "All modules gallery" }, (p) =>
	p
		.section("UI", (s) =>
			s.view(
				widget.ui.callout(
					{ tone: "info", title: "Composable UI" },
					"Cards, stacks, captions, badges, and buttons.",
				),
			),
		)
		.section("CMS", (s) =>
			s.view(
				widget.cms.mediaLibrary([{ id: "asset-1", title: "Hero" }], (m) =>
					m.onSelect(widget.cms.intent.selectAsset(widget.bind.context("asset.id"))),
				),
			),
		)
		.section("Course", (s) =>
			s.view(
				widget.course.shell(
					{
						title: "Course",
						sections: [{ id: "intro", items: [{ id: "start", label: "Start" }] }],
					},
					(c) => c.active("start"),
				),
			),
		)
		.section("Context", (s) =>
			s.view(
				widget.context.diagram(
					{
						id: "ctx",
						title: "Context",
						limit: 1000,
						parts: [{ id: "p1", label: "Prompt", styleKey: "prompt", tokens: 100 }],
					},
					(d) => d.styleSet(styleSet),
				),
			),
		)
		.section("Schedule", (s) => s.view(widget.schedule.availabilityPoll(poll, (b) => b.readOnly())))
		.section("Time", (s) =>
			s.view(
				widget.time.week(
					[
						{
							id: "ev1",
							title: "Lab",
							startISO: "2026-07-08T09:00:00Z",
							endISO: "2026-07-08T10:00:00Z",
							styleKey: "answer",
						},
					],
					(w) => w.styleSet(styleSet),
				),
			),
		),
);
