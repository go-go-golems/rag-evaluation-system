const widget = require("widget.dsl");

function mainFor(active) {
	if (active === "slides") {
		return widget.ui.stack(
			{ gap: "md" },
			widget.ui.card(
				{ title: "Slides" },
				widget.ui.caption("Review speaker notes, deck exports, and presentation state."),
			),
			widget.course.slideDeck(
				{
					index: 0,
					slides: [
						{ id: "s1", title: "Course shell slide", notes: ["Navigation changes the shell body"] },
					],
					snapshot: { id: "ctx", title: "Course shell context", limit: 1024, parts: [] },
				},
				(d) => d.mode("speaker"),
			),
		);
	}
	return widget.ui.card(
		{ title: "Overview" },
		widget.ui.caption("Course body inside shell. Use Slides to switch this panel."),
	);
}

function renderPage(query = {}) {
	const active = query.item === "slides" ? "slides" : "overview";
	const shell = widget.course.shell(
		{
			title: "Course Shell",
			sections: [
				{
					id: "nav",
					label: "Navigation",
					items: [
						{ id: "overview", label: "Overview" },
						{ id: "slides", label: "Slides" },
					],
				},
			],
		},
		(s) =>
			s
				.active(active)
				.onNavigate(widget.course.intent.navigate(widget.bind.context("item.id")))
				.main(mainFor(active)),
	);
	return widget.page("Course shell layout", (p) => p.section("Shell", (s) => s.view(shell)));
}

const page = renderPage({});
