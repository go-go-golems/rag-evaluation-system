const widget = require("widget.dsl");
const form = widget.ui.form(
	{ title: "Contact", formAction: "/api/contact", method: "post" },
	widget.ui.formRow("Name", widget.ui.textInput({ name: "name", defaultValue: "Ada" })),
	widget.ui.formRow("Notes", widget.ui.textareaInput({ name: "notes", defaultValue: "Hello" })),
);
const page = widget.page("UI form", (p) => p.section("Form", (s) => s.view(form)));
