const widget = require("widget.dsl");
const block = widget.raw.element(
	"div",
	{ className: "custom-note" },
	widget.raw.element("h2", {}, "Raw element"),
	widget.raw.element("p", {}, "Use raw elements sparingly for host-specific markup."),
);
const page = widget.page("Raw HTML elements", (p) => p.section("Raw", (s) => s.view(block)));
