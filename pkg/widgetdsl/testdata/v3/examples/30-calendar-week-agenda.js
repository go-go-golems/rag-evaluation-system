const widget = require("widget.dsl");
const events = [
	{
		id: "standup",
		title: "Standup",
		startISO: "2026-07-06T09:00:00Z",
		endISO: "2026-07-06T09:30:00Z",
		styleKey: "meeting",
	},
	{
		id: "demo",
		title: "Demo",
		startISO: "2026-07-10T15:00:00Z",
		endISO: "2026-07-10T16:00:00Z",
		styleKey: "demo",
	},
];
const page = widget.page("Calendar week agenda", (p) =>
	p.section("Calendar", (s) =>
		s.view(
			widget.time.week(events, (w) =>
				w.range(widget.time.range.week("2026-07-06")).now("2026-07-06T10:00:00Z"),
			),
		),
	),
);
