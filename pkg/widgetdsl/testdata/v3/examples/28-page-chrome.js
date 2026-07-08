const widget = require("widget.dsl");
const page = widget.page({ id: "page-chrome", title: "Page chrome" }, (p) =>
	p
		.breadcrumb("Home", "/pages/index")
		.breadcrumb("Chrome")
		.density("compact")
		.section("Actions", (s) =>
			s
				.actions((a) => a.button("Refresh", widget.act.event("refresh")))
				.view(widget.ui.caption("Breadcrumbs and section actions.")),
		),
);
