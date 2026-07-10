const widget = require("widget.dsl");

const items = [
	{
		id: "one",
		label: "First",
		badge: "new",
		description: "A freshly imported record ready for review.",
	},
	{
		id: "two",
		label: "Second",
		badge: "locked",
		description: "A disabled record that still renders its metadata.",
	},
];

const page = widget.page("Data list items", (p) =>
	p.section("Items", (s) =>
		s.view(
			widget.ui.stack(
				{ gap: "md" },
				...items.map((item) =>
					widget.ui.card(
						{ title: item.label, density: "condensed" },
						widget.ui.inline(
							{ gap: "sm", align: "center" },
							widget.ui.badge({ label: item.badge, disabled: item.id === "two" }),
							widget.ui.caption(item.description),
						),
					),
				),
			),
		),
	),
);
