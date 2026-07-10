const widget = require("widget.dsl");
const upload = widget.raw.component("ContextUploadDropArea", {
	title: "Upload JSON",
	description: "Drop a trace JSON file",
	accept: "application/json,.json",
	onUploadAction: widget.act.server("upload.trace"),
});
const page = widget.page("Upload drop area", (p) => p.section("Upload", (s) => s.view(upload)));
