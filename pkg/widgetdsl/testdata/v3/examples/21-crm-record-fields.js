const widget = require("widget.dsl");
const fields = [
	{
		label: "Basics",
		fields: [
			{ key: "name", type: "text", label: "Name" },
			{ key: "stage", type: "select", label: "Stage" },
			{ key: "value", type: "currency", label: "Value" },
		],
	},
];
const page = widget.page("CRM record fields", (p) =>
	p.section("Record", (s) =>
		s.view(
			widget.raw.component("RecordFieldList", {
				values: { name: "Acme", stage: "Proposal", value: 42000 },
				sections: fields,
				mode: "read",
			}),
		),
	),
);
