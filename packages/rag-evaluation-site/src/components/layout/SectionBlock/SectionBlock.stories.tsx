import type { Meta, StoryObj } from "@storybook/react-vite";
import { Button } from "../../atoms";
import { Text } from "../../foundation";
import { Stack } from "../Stack";
import { SectionBlock } from "./SectionBlock";

const meta = {
	title: "Design System/Layout/SectionBlock",
	component: SectionBlock,
	args: {
		label: "What you'll leave with",
		caption: "A plain section with no black chrome.",
		children: (
			<Text>
				Use this for editorial, course, and marketing surfaces where typography carries hierarchy.
			</Text>
		),
	},
} satisfies Meta<typeof SectionBlock>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Plain: Story = {};

export const WithRule: Story = {
	args: {
		rule: true,
		caption: "The 1px rule under the label is the section boundary; no box needed.",
	},
};

export const Levels: Story = {
	render: () => (
		<Stack gap="lg">
			<SectionBlock label="Level 1 section" level={1} rule density="flush">
				<Text>Primary page section.</Text>
			</SectionBlock>
			<SectionBlock label="Level 2 subsection" level={2} rule density="flush">
				<Text>Nested topic inside a section.</Text>
			</SectionBlock>
			<SectionBlock label="Level 3 group" level={3} density="flush">
				<Text>Small field group or aside.</Text>
			</SectionBlock>
		</Stack>
	),
};

export const WithActions: Story = {
	args: {
		rule: true,
		density: "flush",
		actions: (
			<>
				<Button size="compact">Add</Button>
				<Button size="compact">Reorder</Button>
			</>
		),
		caption: "Section-level tools sit at the right of the label row.",
	},
};

export const Flush: Story = {
	args: {
		density: "flush",
		rule: true,
		caption: "No outer padding; the page or stack provides rhythm.",
	},
};
