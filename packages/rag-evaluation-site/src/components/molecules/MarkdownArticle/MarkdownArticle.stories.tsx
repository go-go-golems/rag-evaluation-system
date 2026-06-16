import type { Meta, StoryObj } from "@storybook/react-vite";
import { contextHandoutFixture } from "../../../context";
import { MarkdownArticle } from "./MarkdownArticle";

const meta = {
	title: "Component Library/Molecules/MarkdownArticle",
	component: MarkdownArticle,
	args: { source: contextHandoutFixture.docs[0]!.body },
} satisfies Meta<typeof MarkdownArticle>;

export default meta;
type Story = StoryObj<typeof meta>;

export const FieldGuide: Story = {};

export const Checklist: Story = {
	args: { source: contextHandoutFixture.docs[1]!.body },
};

export const StandaloneImage: Story = {
	args: {
		source: `# Image block

![Context window token budget sketch](/course-assets/context-window-token-budget.svg)

The image should render as a simple borderless figure with an accessible caption derived from alt text.`,
	},
};

export const ImageWithTitle: Story = {
	args: {
		source: `# Titled image

![Budget sketch](/course-assets/context-window-token-budget.svg "A course-authored SVG served from /course-assets")

The title text should be preferred for the visible caption.`,
	},
};

export const MixedMarkdownAndImage: Story = {
	args: {
		source: `# Mixed article

Introductory prose before the visual.

![Context window token budget sketch](/course-assets/context-window-token-budget.svg)

## Checklist

- [x] Prose before the image still renders.
- [x] The image is not swallowed into a paragraph.
- [x] Lists after the image still render as lists.

| Tenant | Action |
| --- | --- |
| Tool output | Summarize after resolution |
| Active task | Keep exact |`,
	},
};

export const ImageBetweenParagraphs: Story = {
	args: {
		source: `A paragraph before the image should stay separate.

![Budget sketch](/course-assets/context-window-token-budget.svg)

A paragraph after the image should not merge into the figure.`,
	},
};
