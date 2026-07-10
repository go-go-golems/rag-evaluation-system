import type { Meta, StoryObj } from "@storybook/react-vite";
import { TileGrid } from "../../layout";
import { StatTile } from "./StatTile";

const meta = {
	title: "Component Library/Molecules/StatTile",
	component: StatTile,
	args: {
		label: "Open pipeline",
		value: "$2.1M",
		delta: 12,
		progress: 0.62,
	},
	decorators: [
		(Story) => (
			<div style={{ width: 200 }}>
				<Story />
			</div>
		),
	],
} satisfies Meta<typeof StatTile>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

/** Positive trend (green up arrow). */
export const Up: Story = {
	args: { label: "Won this quarter", value: "$340k", delta: 8, tone: "success" },
};

/** Negative trend (red down arrow). */
export const Down: Story = {
	args: { label: "New leads", value: "128", delta: -14, tone: "danger" },
};

/** Flat / no change. */
export const Flat: Story = { args: { label: "Avg. deal size", value: "$18k", delta: 0 } };

/** No delta, no bar — just the number. */
export const ValueOnly: Story = {
	args: { delta: undefined, progress: undefined, label: "Open deals", value: "42" },
};

/** Custom delta label instead of a percentage. */
export const CustomDelta: Story = {
	args: {
		label: "Forecast",
		value: "$1.4M",
		trend: "up",
		deltaLabel: "+$180k vs last Q",
		progress: undefined,
	},
};

/** A metric row laid out with TileGrid — the intended dashboard usage. */
export const MetricRow: Story = {
	render: () => (
		<div style={{ width: 720 }}>
			<TileGrid minTileWidth={160}>
				<StatTile label="Open pipeline" value="$2.1M" delta={12} progress={0.62} />
				<StatTile label="Won (Q3)" value="$340k" delta={8} tone="success" progress={0.4} />
				<StatTile label="New leads" value="128" delta={-14} tone="danger" progress={0.3} />
				<StatTile label="Avg. deal" value="$18k" delta={0} />
			</TileGrid>
		</div>
	),
};
