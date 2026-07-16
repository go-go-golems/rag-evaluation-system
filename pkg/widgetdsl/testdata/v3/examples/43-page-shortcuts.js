const widget = require("widget.dsl");

const page = widget.page("Triage", (page) =>
	page
		.shortcuts((keys) =>
			keys
				.bind("accept", "y", widget.act.server("triage.accept"), { label: "Yes" })
				.bind("reject", "n", widget.act.server("triage.reject"), {
					label: "No",
					preventDefault: true,
				})
				.bind("skip", "s", widget.act.server("triage.skip"), { label: "Skip" }),
		)
		.section("Current job", (section) =>
			section
				.caption("Use the visible controls or their keyboard accelerators.")
				.view(
					widget.ui.card(
						{ title: "Go API engineer" },
						widget.ui.inline(
							widget.ui.button("Yes", widget.act.server("triage.accept")),
							widget.ui.button("No", widget.act.server("triage.reject")),
							widget.ui.button("Skip", widget.act.server("triage.skip")),
						),
					),
				),
		),
);

void page;
