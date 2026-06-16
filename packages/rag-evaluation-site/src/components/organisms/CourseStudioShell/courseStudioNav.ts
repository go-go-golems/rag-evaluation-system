import { createElement } from "react";
import { ContextStudioNavIcon, type ContextStudioNavIconId } from "../../atoms";
import type { SidebarNavSection } from "../../molecules";

function icon(id: ContextStudioNavIconId) {
	return createElement(ContextStudioNavIcon, { id });
}

export const courseStudioNavSections: SidebarNavSection[] = [
	{
		id: "present",
		label: "Present",
		items: [
			{ id: "course", label: "Course", icon: icon("course") },
			{ id: "slides", label: "Slides", icon: icon("slides") },
		],
	},
	{
		id: "analyze",
		label: "Analyze",
		items: [
			{ id: "visualize", label: "Visualize", icon: icon("visualize") },
			{ id: "upload", label: "Upload", icon: icon("upload") },
		],
	},
	{
		id: "review",
		label: "Review",
		items: [
			{ id: "transcript", label: "Transcript", icon: icon("transcript") },
			{ id: "comments", label: "Comments", icon: icon("comments") },
		],
	},
	{
		id: "take-home",
		label: "Take-home",
		items: [{ id: "handout", label: "Handout", icon: icon("handout") }],
	},
];
