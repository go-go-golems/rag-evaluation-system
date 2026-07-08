const widget = require("widget.dsl");

const agenda = [
	{ id: "intro", title: "Introduction", duration: 20, status: "ready" },
	{ id: "lab", title: "Hands-on Lab", duration: 60, status: "draft" },
];

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
	(c) => c.active("lab").onNavigate(widget.course.intent.navigate(widget.bind.context("item.id"))),
);

const metadata = widget.course.metadataForm(
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

const agendaEditor = widget.course
	.agendaEditor(agenda, (c) =>
		c.edit((e) =>
			e
				.submitPost("/api/course/agenda")
				.reorder(widget.course.intent.editAgenda(widget.bind.context("rows"))),
		),
	)
	.toNode();

const materials = widget.course.materialUploads(
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

const media = widget.cms.mediaLibrary(
	[{ id: "asset-hero", title: "Hero diagram", kind: "image" }],
	(m) =>
		m
			.onOpen(widget.cms.intent.openAsset(widget.bind.context("asset.id")))
			.onUpload(widget.cms.intent.uploadAssets()),
);

const page = widget.page("Course CMS admin", (p) =>
	p
		.shell({ kind: "course-admin" })
		.section("Course shell", (s) => s.view(shell))
		.section("Metadata", (s) => s.view(metadata))
		.section("Agenda", (s) => s.view(agendaEditor))
		.section("Materials", (s) => s.view(materials))
		.section("Media", (s) => s.view(media)),
);
