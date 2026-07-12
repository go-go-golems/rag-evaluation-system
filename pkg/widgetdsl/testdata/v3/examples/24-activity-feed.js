const widget = require("widget.dsl");
const activities = [
	{
		id: "a1",
		kind: "note",
		title: "Note added",
		body: "Follow up next week",
		atISO: "2026-07-08T10:00:00Z",
		actor: { name: "Mira" },
	},
	{
		id: "a2",
		kind: "email",
		title: "Email sent",
		body: "Proposal shared",
		atISO: "2026-07-08T12:00:00Z",
		actor: { name: "Noah" },
	},
];
const page = widget.page("Activity feed", (p) =>
	p.section("Timeline", (s) =>
		s.view(widget.crm.activityFeed(activities, (feed) => feed.groupByDay(true))),
	),
);
