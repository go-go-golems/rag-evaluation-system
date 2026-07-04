import type { Meta, StoryObj } from "@storybook/react-vite";
import {
	anchoredCommentFixtures,
	contextCourseFixture,
	contextHandoutFixture,
	contextSlides,
	contextWindowSnapshots,
	transcriptFixture,
} from "../context";
import { defaultWidgetRegistry } from "./defaultRegistry";
import { component, text } from "./ir";
import { WidgetRenderer } from "./WidgetRenderer";

const meta = {
	title: "Widget IR/Renderer/Domain Registry Coverage",
	component: WidgetRenderer,
	args: { registry: defaultWidgetRegistry },
} satisfies Meta<typeof WidgetRenderer>;
export default meta;
type Story = StoryObj<typeof meta>;

const snapshot = contextWindowSnapshots[1]!;
const slide = contextSlides[0]!;

export const UiAndDataRegistrySurface: Story = {
	args: {
		node: component("DashboardGrid", { recipe: "twoColumn" }, [
			component(
				"Panel",
				{
					title: "Registered UI widgets",
					actions: component("Button", { variant: "primary" }, [text("Action")]),
				},
				[
					component("Stack", { gap: "sm" }, [
						component("Text", { weight: "bold" }, [
							text("Typography, layout, forms, and tables render through the registry."),
						]),
						component("Inline", { gap: "xs" }, [
							component("StatusText", { status: "ready", icon: true }),
							component("CodeText", {}, [text("ui.dsl")]),
							component("Caption", {}, [text("registered")]),
						]),
						component("DataTable", {
							columns: [
								{ id: "widget", header: "Widget", cell: { kind: "field", field: "widget" } },
								{ id: "module", header: "Module", cell: { kind: "caption", field: "module" } },
							],
							rows: [
								{ id: "panel", widget: "Panel", module: "ui.dsl" },
								{ id: "table", widget: "DataTable", module: "data.dsl" },
							],
							getRowKey: "id",
						}),
					]),
				],
			),
			component("Panel", { title: "Navigation and form controls" }, [
				component("FormRow", {
					label: "Selection",
					control: component("SelectInput", {
						value: "registry",
						options: [{ value: "registry", label: "Registry backed" }],
					}),
					hint: "Rendered by colocated adapters.",
				}),
				component("TabList", {
					activeId: "ui",
					items: [
						{ id: "ui", label: "UI" },
						{ id: "data", label: "Data" },
					],
				}),
			]),
		]),
	},
};

export const ContextWindowRegistrySurface: Story = {
	args: {
		node: component("SplitPane", {
			left: component("Stack", { gap: "md" }, [
				component("ContextBudgetBar", { snapshot, showLegend: true }),
				component("ContextStripDiagram", { snapshot, showLabels: true }),
				component("TranscriptReaderPanel", {
					title: transcriptFixture.title,
					subtitle: transcriptFixture.subtitle,
					messages: transcriptFixture.messages.slice(0, 4),
					annotations: transcriptFixture.annotations,
					selectedAnnotationId: transcriptFixture.selectedAnnotationId,
					showAnnotationChips: true,
				}),
			]),
			right: component("AnchoredCommentRail", {
				comments: anchoredCommentFixtures,
				selectedCommentId: anchoredCommentFixtures[0]?.id,
			}),
			ratio: "rightNarrow",
			gutter: "md",
			divider: true,
		}),
	},
};

export const CourseHandoutRegistrySurface: Story = {
	args: {
		node: component(
			"CourseStudioShell",
			{
				title: "Course registry coverage",
				subtitle: "Course and handout widgets render through course.dsl adapters.",
				activeItemId: "agenda-01",
				sections: [
					{
						id: "course",
						label: "Course",
						items: contextCourseFixture.agenda.slice(0, 3).map((item) => ({
							id: item.id,
							label: item.title,
							icon: component("ContextStudioNavIcon", { id: "lesson" }),
							badge: item.duration,
						})),
					},
				],
			},
			[
				component("SplitPane", {
					left: component("CourseLessonPanel", {
						course: contextCourseFixture,
						activeAgendaItemId: "agenda-01",
					}),
					right: component("Stack", { gap: "md" }, [
						component("CourseSlidePanel", {
							slide,
							snapshot,
							index: 1,
							total: contextSlides.length,
						}),
						component("HandoutDocumentShell", {
							intro: contextHandoutFixture.intro,
							documents: contextHandoutFixture.docs,
							selectedDocumentId: contextHandoutFixture.docs[0]?.id,
							title: "Take-home packet",
						}),
					]),
					ratio: "course",
					gutter: "md",
					divider: true,
				}),
			],
		),
	},
};
