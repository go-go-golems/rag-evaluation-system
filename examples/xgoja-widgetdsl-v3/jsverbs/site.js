__package__({ name: "site", short: "widget.dsl v3 example site" });

function start() {
	require("./server");
}

__verb__("start", {
	name: "start",
	short: "Serve the widget.dsl v3 example gallery",
	output: "text",
});
