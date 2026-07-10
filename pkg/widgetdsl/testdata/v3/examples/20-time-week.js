const widget = require("widget.dsl");
const events = [
	{
		id: "e1",
		title: "Workshop",
		startISO: "2026-07-08T09:00:00Z",
		endISO: "2026-07-08T11:00:00Z",
		styleKey: "busy",
	},
	{
		id: "e2",
		title: "Retro",
		startISO: "2026-07-09T15:00:00Z",
		endISO: "2026-07-09T16:00:00Z",
		styleKey: "focus",
	},
];
const page = widget.page("Time week", (p) =>
	p.section("Week", (s) =>
		s.view(
			widget.time.week(events, (w) =>
				w
					.range(widget.time.range.week("2026-07-08"))
					.hours(8, 18)
					.onSelect(widget.time.intent.selectEvent(widget.bind.context("block.id"))),
			),
		),
	),
);
