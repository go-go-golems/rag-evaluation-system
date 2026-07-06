import type { Meta, StoryObj } from "@storybook/react-vite";
import { cmsArticleFixtures, cmsAssetFixtures } from "../cms";
import { defaultWidgetRegistry } from "./defaultRegistry";
import { component, text } from "./ir";
import { WidgetRenderer } from "./WidgetRenderer";

const meta = {
	title: "Widget IR/Renderer/CMS",
	component: WidgetRenderer,
	args: { registry: defaultWidgetRegistry },
} satisfies Meta<typeof WidgetRenderer>;
export default meta;
type Story = StoryObj<typeof meta>;

export const AtomsGallery: Story = {
	args: {
		node: component("Stack", { gap: "md" }, [
			component("Inline", { gap: "sm" }, [
				component("ContentStatusBadge", { status: "draft" }),
				component("ContentStatusBadge", { status: "published" }),
				component("ContentStatusBadge", { status: "scheduled" }),
				component("ContentStatusBadge", { status: "archived" }),
			]),
			component("Inline", { gap: "sm" }, [
				component("Tag", { label: "course" }),
				component("Tag", { label: "context-window", selected: true }),
				component("Tag", {
					label: "removable",
					onRemoveAction: { kind: "event", event: "widget-ir:tag-removed" },
				}),
			]),
			component("MeterBar", { value: 0.62, label: "62%", style: { width: 240 } }),
			component("TileGrid", { minTileWidth: 120, style: { maxWidth: 560 } }, [
				component("MediaThumb", {
					src: "/course-assets/context-window-token-budget.svg",
					alt: "Budget sketch",
				}),
				component("MediaThumb", { src: "/media/missing/broken.png", alt: "Broken" }),
				component("MediaThumb", {}),
				component("MediaThumb", {
					src: "/course-assets/context-window-token-budget.svg",
					alt: "Selected",
					selected: true,
				}),
			]),
		]),
	},
};

export const MediaLibraryFromIr: Story = {
	args: {
		node: component("MediaLibraryPanel", {
			assets: cmsAssetFixtures,
			selectedAssetIds: [cmsAssetFixtures[1]?.id ?? ""],
			query: "",
			onQuerySubmitAction: { kind: "navigate", to: "?query=$query" },
			kindFilter: "all",
			onKindFilterChangeAction: { kind: "navigate", to: "?kind=$kind" },
			page: 1,
			pageCount: 3,
			onPageChangeAction: { kind: "navigate", to: "?page=$page" },
			onAssetSelectAction: { kind: "navigate", to: "?asset=$assetId" },
			onFilesSelectedAction: {
				kind: "server",
				name: "admin-upload-course-material",
				payload: { kind: "media" },
			},
			style: { maxWidth: 880 },
		}),
	},
};

export const ArticleListFromIr: Story = {
	args: {
		node: component("ArticleListPanel", {
			articles: cmsArticleFixtures,
			selectedArticleId: cmsArticleFixtures[2]?.id,
			onArticleSelectAction: { kind: "navigate", to: "?article=$articleId" },
			onCreateAction: { kind: "navigate", to: "/pages/admin-article-new" },
			onRowActionAction: { kind: "event", event: "widget-ir:article-row-action" },
			statusFilter: "all",
			onStatusFilterChangeAction: { kind: "navigate", to: "?status=$status" },
			page: 1,
			pageCount: 2,
			onPageChangeAction: { kind: "navigate", to: "?page=$page" },
			style: { maxWidth: 960 },
		}),
	},
};

export const MarkdownEditorWithLivePreview: Story = {
	args: {
		node: component("MarkdownEditor", {
			name: "body",
			defaultValue:
				"# Adapter-local state\n\nThe toolbar and this **live preview** run entirely in the browser — no server round-trip.\n\n- [x] named textarea joins native form posts\n- [ ] save via formPanel formAction",
			minRows: 12,
		}),
	},
};

export const ConfirmedDestructiveAction: Story = {
	args: {
		node: component("Inline", { gap: "sm" }, [
			component(
				"Button",
				{
					action: {
						kind: "event",
						event: "widget-ir:deleted",
						confirm: "Delete “missing-figure.png” permanently?",
					},
				},
				[text("× Delete with confirm")],
			),
			component("Button", { action: { kind: "event", event: "widget-ir:deleted" } }, [
				text("Delete without confirm"),
			]),
		]),
	},
};

export const CmsShellFromIr: Story = {
	args: {
		node: component(
			"CmsShell",
			{
				activeItemId: "media",
				onNavigateAction: { kind: "navigate", to: "/pages/$itemId" },
			},
			[
				component("MediaLibraryPanel", {
					assets: cmsAssetFixtures.slice(0, 6),
					showStatusBadges: true,
				}),
			],
		),
	},
};

export const BreadcrumbsPaginationSearch: Story = {
	args: {
		node: component("Stack", { gap: "md" }, [
			component("Breadcrumbs", {
				items: [
					{ id: "media", label: "Media" },
					{ id: "course", label: "Course assets" },
					{ id: "diagrams", label: "Diagrams" },
				],
				onNavigateAction: { kind: "navigate", to: "?folder=$itemId" },
			}),
			component("SearchField", {
				defaultValue: "context",
				onSubmitAction: { kind: "navigate", to: "?query=$query" },
				style: { width: 240 },
			}),
			component("Pagination", {
				page: 2,
				pageCount: 9,
				pageSize: 24,
				totalItems: 210,
				onPageChangeAction: { kind: "navigate", to: "?page=$page" },
			}),
			component("EmptyState", {
				glyph: "▨",
				title: "No assets match this filter",
				hint: "Try clearing the search.",
				actionSlot: component("Button", { variant: "primary" }, [text("Clear filters")]),
			}),
		]),
	},
};
