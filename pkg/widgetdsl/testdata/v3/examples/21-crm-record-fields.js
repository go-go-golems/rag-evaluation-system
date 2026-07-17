const widget = require("widget.dsl");
const fields = widget.crm.fields("deal", (f) =>
	f
		.text("name", { label: "Name", section: "Basics" })
		.select("stage", { label: "Stage", section: "Basics" })
		.currency("value", { label: "Value", section: "Basics" }),
);
const page = widget.page("CRM record fields", (p) =>
	p.section("Record", (s) =>
		s.view(
			widget.crm.recordFields({ name: "Acme", stage: "Proposal", value: 42000 }, fields, (r) =>
				r.mode("read"),
			),
		),
	),
);
