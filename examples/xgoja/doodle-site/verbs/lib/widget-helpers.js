function createWidgetHelpers(widget) {
	const act = widget.act;
	const navItems = [
		{ id: "index", label: "All polls", action: act.navigate("/pages/index") },
		{ id: "create", label: "New poll", action: act.navigate("/pages/create") },
	];

	function statusText(status, text) {
		return widget.ui.status(status, text);
	}

	function emptyState(title, description) {
		return widget.ui.emptyState(title, description);
	}

	function formRow(label, control, options = {}) {
		return widget.ui.formRow(label, control, options);
	}

	function textInput(props) {
		return widget.ui.textInput(props);
	}

	function textareaInput(props) {
		return widget.ui.textareaInput(props);
	}

	function selectInput(props) {
		return widget.ui.selectInput(props);
	}

	function collectionTable(name, rows, configureFields, options = {}) {
		const schema = widget.data.fields(name, configureFields).build();
		return widget.data
			.collection(name, rows, (c) => {
				c.schema(schema)
					.empty(options.empty || "No rows")
					.table();
			})
			.toNode();
	}

	function applyPageMeta(p, id, activeNavItemId) {
		const shell = widget.app.shell((builder) =>
			builder
				.brand("Doodle")
				.navigation((navigation) =>
					navigation
						.placement("top")
						.active(activeNavItemId)
						.ariaLabel("Polls")
						.section("polls", "Polls", (items) => {
							navItems.forEach((item) => items.item(item.id, item.label, item.action));
						}),
				)
				.content((content) => content.maxWidth("wide").padding("none")),
		);
		return p.id(id).shell(shell);
	}

	function asPage(pageBuilder) {
		return pageBuilder.toPage();
	}

	return {
		act,
		navItems,
		statusText,
		emptyState,
		formRow,
		textInput,
		textareaInput,
		selectInput,
		collectionTable,
		applyPageMeta,
		asPage,
	};
}

module.exports = { createWidgetHelpers };
