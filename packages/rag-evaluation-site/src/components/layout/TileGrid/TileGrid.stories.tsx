import type { Meta, StoryObj } from "@storybook/react-vite";
import { MediaThumb } from "../../atoms";
import { TileGrid } from "./TileGrid";

const SKETCH_SRC = "/course-assets/context-window-token-budget.svg";

const tileIds = Array.from({ length: 50 }, (_, index) => `tile-${index + 1}`);

const meta = {
	title: "Design System/Layout/TileGrid",
	component: TileGrid,
} satisfies Meta<typeof TileGrid>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {
	render: () => (
		<TileGrid>
			{tileIds.slice(0, 8).map((id) => (
				<MediaThumb key={id} src={SKETCH_SRC} alt={id} />
			))}
		</TileGrid>
	),
};

export const DenseFiftyTiles: Story = {
	render: () => (
		<TileGrid minTileWidth={96} gap="sm">
			{tileIds.map((id) => (
				<MediaThumb key={id} src={SKETCH_SRC} alt={id} />
			))}
		</TileGrid>
	),
};

export const NarrowContainer: Story = {
	render: () => (
		<div style={{ maxWidth: 360, border: "1px solid var(--mac-border-subtle)", padding: 8 }}>
			<TileGrid minTileWidth={120}>
				{tileIds.slice(0, 5).map((id) => (
					<MediaThumb key={id} src={SKETCH_SRC} alt={id} />
				))}
			</TileGrid>
		</div>
	),
};
