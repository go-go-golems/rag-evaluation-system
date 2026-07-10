const widget = require("widget.dsl");
const availability = {
	title: "Rooms",
	resources: [
		{ id: "room-a", label: "Room A", availability: { wed9: "available" } },
		{ id: "room-b", label: "Room B", availability: { wed9: "unavailable" } },
	],
	slots: [
		{
			id: "wed9",
			label: "Wed 9",
			startISO: "2026-07-08T09:00:00Z",
			endISO: "2026-07-08T10:00:00Z",
		},
	],
};
const page = widget.page("Booking picker", (p) =>
	p.section("Rooms", (s) =>
		s.view(
			widget.schedule.bookingPicker(availability, (b) =>
				b.onToggle(
					widget.schedule.intent.toggleAvailability(
						widget.bind.context("row.id"),
						widget.bind.context("column.id"),
					),
				),
			),
		),
	),
);
