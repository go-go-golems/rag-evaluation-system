const widget = require("widget.dsl");

const agenda = [
	{ id: "intro", title: "Introduction", duration: 20, status: "ready" },
	{ id: "lab", title: "Hands-on Lab", duration: 60, status: "draft" },
];

const mediaAssets = [
	{
		id: "asset-hero",
		title: "Hero diagram",
		kind: "image",
		mime: "image/png",
		filename: "hero-diagram.png",
		size: 182400,
	},
];

function metadataWidget() {
	return widget.course.metadataForm(
		{
			title: "Context Window Engineering",
			status: "draft",
			owner: "Manuel",
		},
		(f) =>
			f.title("Course metadata").onSubmit(
				widget.act.server("course.metadata.update", {
					payload: { form: widget.bind.context("form") },
				}),
			),
	);
}

function agendaWidget() {
	return widget.course
		.agendaEditor(agenda, (c) =>
			c.edit((e) =>
				e
					.submitPost("/api/course/agenda")
					.reorder(widget.course.intent.editAgenda(widget.bind.context("rows"))),
			),
		)
		.toNode();
}

function materialsWidget() {
	return widget.course.materialUploads(
		{
			title: "Upload material",
			description: "Slides, handouts, diagrams",
		},
		(u) =>
			u
				.accept(["application/pdf", "image/png"])
				.onUpload(widget.course.intent.uploadMaterial())
				.onDelete(widget.course.intent.deleteMaterial(widget.bind.context("asset.id"))),
	);
}

function mediaWidget() {
	return widget.cms.mediaLibrary(mediaAssets, (m) =>
		m
			.onOpen(widget.cms.intent.openAsset(widget.bind.context("asset.id")))
			.onUpload(widget.cms.intent.uploadAssets()),
	);
}

function mainForItem(item) {
	if (item === "intro") {
		return widget.ui.stack(
			{ gap: "md" },
			widget.ui.card(
				{ title: "Intro module" },
				widget.ui.caption("Edit the course metadata and hero media for the introduction."),
			),
			metadataWidget(),
			mediaWidget(),
		);
	}
	return widget.ui.stack(
		{ gap: "md" },
		widget.ui.card(
			{ title: "Lab module" },
			widget.ui.caption("Edit the lab agenda and upload supporting workshop material."),
		),
		agendaWidget(),
		materialsWidget(),
	);
}

function renderPage(query = {}) {
	const active = query.item === "intro" ? "intro" : "lab";
	const shell = widget.course.shell(
		{
			title: "Context Window Engineering",
			subtitle: "Admin preview",
			sections: [
				{
					id: "modules",
					label: "Modules",
					items: [
						{ id: "intro", label: "Intro" },
						{ id: "lab", label: "Lab" },
					],
				},
			],
		},
		(c) =>
			c
				.active(active)
				.onNavigate(widget.course.intent.navigate(widget.bind.context("item.id")))
				.main(mainForItem(active)),
	);

	return widget.page("Course CMS admin", (p) =>
		p.shell({ kind: "course-admin" }).section("Course shell", (s) => s.view(shell)),
	);
}

const page = renderPage({});
