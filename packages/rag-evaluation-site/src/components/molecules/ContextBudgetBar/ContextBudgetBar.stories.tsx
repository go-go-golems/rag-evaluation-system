import type { Meta, StoryObj } from "@storybook/react-vite";
import {
	contextCobaltSandStyleSet,
	contextPaletteOptions,
	contextSignalOrangeStyleSet,
	contextStyleSetForPalette,
	contextWindowSnapshots,
	type ContextPaletteName,
} from "../../../context";
import { Panel, Stack } from "../../layout";
import { ContextBudgetBar, type ContextBudgetBarProps } from "./ContextBudgetBar";

const [underBudget, selectedBudget, nearLimit, overBudget] = contextWindowSnapshots;

type StoryArgs = Omit<ContextBudgetBarProps, "styleSet"> & { palette: ContextPaletteName };

const meta = {
	title: "Component Library/Molecules/ContextBudgetBar",
	args: { snapshot: selectedBudget!, palette: "Dusty Magenta / Blue", showLegend: true },
	argTypes: {
		palette: { control: "select", options: contextPaletteOptions },
		snapshot: { control: false },
		selectedPartId: { control: "text" },
		showLegend: { control: "boolean" },
	},
} satisfies Meta<StoryArgs>;

export default meta;
type Story = StoryObj<StoryArgs>;

export const BudgetStates: Story = {
	render: ({ palette, showLegend }) => {
		const styleSet = contextStyleSetForPalette(palette);
		return (
			<Stack gap="md">
				<Panel title="under budget">
					<ContextBudgetBar snapshot={underBudget!} styleSet={styleSet} showLegend={showLegend} />
				</Panel>
				<Panel title="near limit">
					<ContextBudgetBar snapshot={nearLimit!} styleSet={styleSet} showLegend={showLegend} />
				</Panel>
				<Panel title="over budget">
					<ContextBudgetBar snapshot={overBudget!} styleSet={styleSet} showLegend={showLegend} />
				</Panel>
			</Stack>
		);
	},
};

export const SelectedSegment: Story = {
	args: { selectedPartId: "t14-file-reads" },
	render: ({ palette, selectedPartId, showLegend }) => (
		<Panel title={`selected context tenant · ${palette}`}>
			<ContextBudgetBar
				snapshot={selectedBudget!}
				styleSet={contextStyleSetForPalette(palette)}
				selectedPartId={selectedPartId}
				showLegend={showLegend}
			/>
		</Panel>
	),
};

export const SameBudgetDifferentPalettes: Story = {
	render: ({ palette }) => (
		<Stack gap="md">
			<Panel title={`${palette} (controls-selected)`}>
				<ContextBudgetBar
					snapshot={selectedBudget!}
					styleSet={contextStyleSetForPalette(palette)}
				/>
			</Panel>
			<Panel title="Signal Orange / Cyan reference">
				<ContextBudgetBar snapshot={selectedBudget!} styleSet={contextSignalOrangeStyleSet} />
			</Panel>
			<Panel title="Cobalt / Sand reference">
				<ContextBudgetBar snapshot={selectedBudget!} styleSet={contextCobaltSandStyleSet} />
			</Panel>
		</Stack>
	),
};
