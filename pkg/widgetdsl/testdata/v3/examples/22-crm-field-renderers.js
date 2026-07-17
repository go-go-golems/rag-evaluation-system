const widget = require("widget.dsl");
const page = widget.page("CRM field renderers", (p) =>
	p.section("Fields", (s) =>
		s.view(
			widget.ui.stack(
				{},
				widget.crm.field("ops@example.com", {
					key: "email",
					type: "email",
					label: "Email",
				}),
				widget.crm.field("Active", {
					key: "status",
					type: "select",
					label: "Status",
				}),
			),
		),
	),
);
