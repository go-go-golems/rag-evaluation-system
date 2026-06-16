import type { Meta, StoryObj } from "@storybook/react-vite";
import { Panel } from "../Panel";
import { Stack } from "./Stack";

const meta = {
	title: "Design System/Layout/Stack",
	component: Stack,
} satisfies Meta<typeof Stack>;

export default meta;
type Story = StoryObj<typeof meta>;

export const PanelStack: Story = {
	render: () => (
		<Stack gap="sm">
			<Panel title="Query" density="condensed">
				Controls
			</Panel>
			<Panel title="Coverage" density="condensed">
				204/204 embedded
			</Panel>
			<Panel title="Presets" density="condensed">
				hydrangea pruning
			</Panel>
		</Stack>
	),
};
