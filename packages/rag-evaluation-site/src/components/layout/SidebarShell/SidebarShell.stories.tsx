import type { Meta, StoryObj } from "@storybook/react-vite";
import { Text } from "../../foundation";
import { SidebarNav } from "../../molecules";
import { courseStudioNavSections } from "../../organisms/CourseStudioShell/courseStudioNav";
import { SectionBlock } from "../SectionBlock";
import { SidebarShell } from "./SidebarShell";

const meta = {
	title: "Design System/Layout/SidebarShell",
	component: SidebarShell,
	args: {
		sidebar: <SidebarNav sections={courseStudioNavSections} activeItemId="slides" />,
		header: <Text weight="bold">Context Window Engineering</Text>,
		children: (
			<SectionBlock label="Content">
				<Text>The shell owns the fixed sidebar and the scrollable content viewport.</Text>
			</SectionBlock>
		),
	},
} satisfies Meta<typeof SidebarShell>;
export default meta;
type Story = StoryObj<typeof meta>;

export const CourseSidebar: Story = {};
