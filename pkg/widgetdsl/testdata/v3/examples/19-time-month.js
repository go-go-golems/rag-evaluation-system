const widget = require("widget.dsl");
const events = [
	{
		id: "e1",
		title: "Launch",
		startISO: "2026-07-08T09:00:00Z",
		endISO: "2026-07-08T10:00:00Z",
		styleKey: "event",
	},
	{
		id: "e2",
		title: "Review",
		startISO: "2026-07-16T13:00:00Z",
		endISO: "2026-07-16T14:00:00Z",
		styleKey: "event",
	},
];
const page = widget.page("Time month", (p) =>
	p.section("July", (s) =>
		s.view(
			widget.time.month({ monthISO: "2026-07", events }, (m) =>
				m
					.selected("2026-07-08")
					.onSelect(widget.time.intent.selectDay(widget.bind.context("dayISO"))),
			),
		),
	),
);
