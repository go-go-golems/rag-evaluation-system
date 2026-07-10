const widget = require("widget.dsl");
const poll = {
	title: "Readonly RSVP",
	options: [
		{
			id: "slot",
			label: "Friday",
			startISO: "2026-07-10T12:00:00Z",
			endISO: "2026-07-10T13:00:00Z",
		},
	],
	responses: [{ id: "team", name: "Team", availability: { slot: "available" } }],
};
const page = widget.page("Readonly schedule", (p) =>
	p.section("Poll", (s) =>
		s.view(
			widget.schedule.availabilityPoll(poll, (b) =>
				b.onToggle(widget.schedule.intent.toggleAvailability("team", "slot")).readOnly(),
			),
		),
	),
);
