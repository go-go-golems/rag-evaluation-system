const widget = require("widget.dsl");
const poll = {
	title: "Office hours",
	options: [
		{ id: "mon", label: "Mon 9", startISO: "2026-07-06T09:00:00Z", endISO: "2026-07-06T10:00:00Z" },
		{
			id: "tue",
			label: "Tue 10",
			startISO: "2026-07-07T10:00:00Z",
			endISO: "2026-07-07T11:00:00Z",
		},
	],
	responses: [
		{ id: "ana", name: "Ana", availability: { mon: "available", tue: "maybe" } },
		{ id: "noah", name: "Noah", availability: { mon: "unavailable", tue: "available" } },
	],
};
const page = widget.page("Editable poll", (p) =>
	p.section("Poll", (s) =>
		s.view(
			widget.schedule.availabilityPoll(poll, (b) =>
				b
					.editableRow("ana")
					.onToggle(
						widget.schedule.intent.toggleAvailability(
							widget.bind.context("row.id"),
							widget.bind.context("column.id"),
						),
					),
			),
		),
	),
);
