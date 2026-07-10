const widget = require("widget.dsl");
const page = widget.page("Dashboard metrics", (p) =>
	p.section("Metrics", (s) =>
		s
			.metric("Documents", "128")
			.metric("Chunks", "42k")
			.metric("Recall", "91%")
			.view(
				widget.ui.callout({ tone: "success", title: "Healthy" }, "The evaluation corpus is ready."),
			),
	),
);
