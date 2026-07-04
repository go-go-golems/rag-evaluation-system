import { ContextStudioNavIcon } from "../../atoms";
import type { SidebarNavSection } from "../../molecules";

export const cmsNavSections: SidebarNavSection[] = [
	{
		id: "content",
		label: "Content",
		items: [
			{
				id: "articles",
				label: "Articles",
				icon: <ContextStudioNavIcon id="handout" title="Articles" />,
			},
			{
				id: "media",
				label: "Media",
				icon: <ContextStudioNavIcon id="upload" title="Media" />,
			},
		],
	},
	{
		id: "organize",
		label: "Organize",
		items: [
			{
				id: "tags",
				label: "Tags",
				icon: <ContextStudioNavIcon id="comments" title="Tags" />,
			},
			{
				id: "archive",
				label: "Archive",
				icon: <ContextStudioNavIcon id="visualize" title="Archive" />,
			},
		],
	},
];
