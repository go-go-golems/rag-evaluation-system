import type { Meta, StoryObj } from "@storybook/react-vite";
import {
	type ContextPaletteName,
	contextCobaltSandStyleSet,
	contextPaletteOptions,
	contextStyleSetForPalette,
	contextWindowSnapshots,
} from "../../../context";
import { Panel, Stack } from "../../layout";
import { ContextTreemap, type ContextTreemapProps } from "./ContextTreemap";

const [, deepBug, atLimit] = contextWindowSnapshots;

type StoryArgs = Omit<ContextTreemapProps, "styleSet"> & { palette: ContextPaletteName };

const meta = {
	title: "Component Library/Molecules/ContextTreemap",
	args: { snapshot: deepBug!, palette: "Dusty Magenta / Blue" },
	argTypes: {
		palette: { control: "select", options: contextPaletteOptions },
		snapshot: { control: false },
		selectedPartId: { control: "text" },
	},
} satisfies Meta<StoryArgs>;
export default meta;
type Story = StoryObj<StoryArgs>;

export const ProportionalTokens: Story = {
	render: ({ palette }) => (
		<Panel title={`where the tokens went · ${palette}`}>
			<ContextTreemap snapshot={atLimit!} styleSet={contextStyleSetForPalette(palette)} />
		</Panel>
	),
};

export const SelectedTile: Story = {
	args: { selectedPartId: "t14-file-reads" },
	render: ({ palette, selectedPartId }) => (
		<Panel title={`selected file reads · ${palette}`}>
			<ContextTreemap
				snapshot={deepBug!}
				styleSet={contextStyleSetForPalette(palette)}
				selectedPartId={selectedPartId}
			/>
		</Panel>
	),
};

export const Comparison: Story = {
	render: ({ palette }) => {
		const styleSet = contextStyleSetForPalette(palette);
		return (
			<Stack gap="md">
				<Panel title="turn 14">
					<ContextTreemap snapshot={deepBug!} styleSet={styleSet} />
				</Panel>
				<Panel title="turn 31">
					<ContextTreemap snapshot={atLimit!} styleSet={styleSet} />
				</Panel>
				<Panel title="turn 31 / cobalt sand reference">
					<ContextTreemap snapshot={atLimit!} styleSet={contextCobaltSandStyleSet} />
				</Panel>
			</Stack>
		);
	},
};
