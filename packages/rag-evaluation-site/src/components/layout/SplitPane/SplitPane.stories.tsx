import type { Meta, StoryObj } from "@storybook/react-vite";
import { Text } from "../../foundation";
import { SectionBlock } from "../SectionBlock";
import { SplitPane } from "./SplitPane";
const meta = {
	title: "Design System/Layout/SplitPane",
	component: SplitPane,
	args: {
		left: (
			<SectionBlock label="Left">
				<Text>Outcomes and supporting details.</Text>
			</SectionBlock>
		),
		right: (
			<SectionBlock label="Right">
				<Text>Agenda or related material.</Text>
			</SectionBlock>
		),
	},
} satisfies Meta<typeof SplitPane>;
export default meta;
type Story = StoryObj<typeof meta>;
export const CourseRatio: Story = { args: { ratio: "course", divider: true } };
