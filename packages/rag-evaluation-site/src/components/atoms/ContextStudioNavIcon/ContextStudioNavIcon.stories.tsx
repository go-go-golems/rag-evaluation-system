import type { Meta, StoryObj } from "@storybook/react-vite";
import { Inline } from "../../layout";
import { ContextStudioNavIcon, type ContextStudioNavIconId } from "./ContextStudioNavIcon";

const ids: ContextStudioNavIconId[] = [
	"course",
	"slides",
	"visualize",
	"upload",
	"transcript",
	"comments",
	"handout",
];

const meta = {
	title: "Design System/Atoms/ContextStudioNavIcon",
	component: ContextStudioNavIcon,
	args: { id: "upload" },
} satisfies Meta<typeof ContextStudioNavIcon>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Upload: Story = {};

export const AllIcons: Story = {
	render: () => (
		<Inline gap="lg">
			{ids.map((id) => (
				<ContextStudioNavIcon key={id} id={id} title={id} />
			))}
		</Inline>
	),
};

export const Inverted: Story = {
	render: () => (
		<div
			style={{
				display: "inline-flex",
				padding: 8,
				background: "var(--mac-bg-dark)",
				color: "var(--mac-text-inv)",
			}}
		>
			<Inline gap="lg">
				{ids.map((id) => (
					<ContextStudioNavIcon key={id} id={id} title={id} />
				))}
			</Inline>
		</div>
	),
};
