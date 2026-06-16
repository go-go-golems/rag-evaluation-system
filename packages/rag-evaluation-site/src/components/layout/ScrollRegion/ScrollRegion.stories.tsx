import type { Meta, StoryObj } from "@storybook/react-vite";
import { ScrollRegion } from "./ScrollRegion";

const meta = { title: "Design System/Layout/ScrollRegion", component: ScrollRegion } satisfies Meta<
	typeof ScrollRegion
>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Vertical: Story = {
	render: () => (
		<ScrollRegion style={{ height: 120, border: "1px solid var(--mac-border)" }}>
			{Array.from({ length: 20 }, (_, i) => (
				<div key={i}>Row {i + 1}</div>
			))}
		</ScrollRegion>
	),
};
