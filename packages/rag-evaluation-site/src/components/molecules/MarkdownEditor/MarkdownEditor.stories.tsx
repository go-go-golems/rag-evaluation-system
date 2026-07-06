import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { Button } from "../../atoms";
import { SplitPane } from "../../layout";
import { MarkdownArticle } from "../MarkdownArticle";
import { MarkdownEditor } from "./MarkdownEditor";

const SAMPLE = `# Draft article

Introductory prose with **bold** and \`code\`.

- [x] Toolbar wraps the selection
- [ ] Preview stays live
`;

const meta = {
	title: "Component Library/Molecules/MarkdownEditor",
	component: MarkdownEditor,
	args: { value: SAMPLE, onValueChange: () => {} },
} satisfies Meta<typeof MarkdownEditor>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	args: { value: "" },
};

export const WithContent: Story = {};

export const MaxLength: Story = {
	args: { maxLength: 280, value: "Short body with a visible counter." },
};

export const Disabled: Story = {
	args: { disabled: true },
};

export const ToolbarSlot: Story = {
	args: {
		onInsertAsset: () => {},
		toolbarSlot: <Button size="compact">✓ spellcheck</Button>,
	},
};

export const Interactive: Story = {
	render: () => {
		const [value, setValue] = useState(SAMPLE);
		return (
			<SplitPane
				divider
				gutter="md"
				left={<MarkdownEditor value={value} onValueChange={setValue} minRows={12} />}
				right={<MarkdownArticle source={value} />}
			/>
		);
	},
};
