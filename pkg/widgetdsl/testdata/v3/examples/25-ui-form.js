const widget = require("widget.dsl");
const form = widget.ui.form(
	{ title: "Contact", formAction: "/api/contact", method: "post" },
	widget.raw.component("FormRow", {
		label: "Name",
		control: widget.raw.component("TextInput", { name: "name", defaultValue: "Ada" }),
	}),
	widget.raw.component("FormRow", {
		label: "Notes",
		control: widget.raw.component("TextareaInput", { name: "notes", defaultValue: "Hello" }),
	}),
);
const page = widget.page("UI form", (p) => p.section("Form", (s) => s.view(form)));
