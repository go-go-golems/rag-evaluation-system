const widget = require("widget.dsl");
const upload = widget.ui.upload({
	title: "Upload JSON",
	description: "Drop a trace JSON file",
	accept: "application/json,.json",
	onFilesSelectedAction: widget.act.server("upload.trace"),
});
const page = widget.page("Upload drop area", (p) => p.section("Upload", (s) => s.view(upload)));
