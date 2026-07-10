const widget = require("widget.dsl");
const page = widget.page("CRM field renderers", (p) =>
	p.section("Fields", (s) =>
		s.view(
			widget.ui.stack(
				{},
				widget.raw.component("FieldRenderer", {
					spec: { key: "email", type: "email", label: "Email" },
					value: "ops@example.com",
				}),
				widget.raw.component("FieldRenderer", {
					spec: { key: "status", type: "select", label: "Status" },
					value: "Active",
				}),
			),
		),
	),
);
