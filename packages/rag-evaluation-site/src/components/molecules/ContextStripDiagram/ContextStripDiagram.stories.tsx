import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import {
	type ContextPaletteName,
	type ContextStyleSet,
	contextCobaltSandStyleSet,
	contextDefaultStyleSet,
	contextPaletteOptions,
	contextSignalOrangeStyleSet,
	contextSlateCoralStyleSet,
	contextStyleSetForPalette,
	contextWindowSnapshots,
} from "../../../context";
import { Button } from "../../atoms";
import { Inline, Panel, Stack } from "../../layout";
import { ContextLegend } from "../ContextLegend";
import { ContextStripDiagram, type ContextStripDiagramProps } from "./ContextStripDiagram";

const [, deepBug, atLimit, overBudget] = contextWindowSnapshots;

type StoryArgs = Omit<ContextStripDiagramProps, "styleSet"> & {
	palette: ContextPaletteName;
	showLegend: boolean;
};

const meta = {
	title: "Component Library/Molecules/ContextStripDiagram",
	args: { snapshot: deepBug!, palette: "Dusty Magenta / Blue", showLegend: true },
	argTypes: {
		palette: { control: "select", options: contextPaletteOptions },
		snapshot: { control: false },
		selectedPartId: { control: "text" },
		showLabels: { control: "boolean" },
		showLegend: { control: "boolean" },
	},
} satisfies Meta<StoryArgs>;

export default meta;
type Story = StoryObj<StoryArgs>;

export const DenseSegments: Story = {
	render: ({ palette, showLegend }) => {
		const styleSet = contextStyleSetForPalette(palette);
		return (
			<Panel title={`turn 14 context strip · ${palette}`}>
				<Stack gap="sm">
					<ContextStripDiagram snapshot={deepBug!} styleSet={styleSet} />
					{showLegend && (
						<ContextLegend items={styleSet.legend} styles={styleSet.styles} size="sm" />
					)}
				</Stack>
			</Panel>
		);
	},
};

export const SelectedSegment: Story = {
	args: { selectedPartId: "t14-file-reads" },
	render: ({ palette, selectedPartId }) => (
		<Panel title={`selected file reads · ${palette}`}>
			<ContextStripDiagram
				snapshot={deepBug!}
				styleSet={contextStyleSetForPalette(palette)}
				selectedPartId={selectedPartId}
			/>
		</Panel>
	),
};

export const LimitComparison: Story = {
	render: ({ palette }) => {
		const styleSet = contextStyleSetForPalette(palette);
		return (
			<Stack gap="md">
				<Panel title="at limit">
					<ContextStripDiagram snapshot={atLimit!} styleSet={styleSet} />
				</Panel>
				<Panel title="before reclaim">
					<ContextStripDiagram snapshot={overBudget!} styleSet={styleSet} />
				</Panel>
				<Panel title="at limit / signal orange reference">
					<ContextStripDiagram snapshot={atLimit!} styleSet={contextSignalOrangeStyleSet} />
				</Panel>
			</Stack>
		);
	},
};

const stripStyleSets: ContextStyleSet[] = [
	contextDefaultStyleSet,
	contextSignalOrangeStyleSet,
	contextSlateCoralStyleSet,
	contextCobaltSandStyleSet,
];

export const InteractiveStyleSwitcher: Story = {
	render: ({ palette }) => {
		const [styleSet, setStyleSet] = useState<ContextStyleSet>(contextStyleSetForPalette(palette));
		return (
			<Panel title="switch the same strip between style sets">
				<Stack gap="sm">
					<Inline gap="xs" wrap>
						{stripStyleSets.map((candidate) => (
							<Button
								key={candidate.id ?? candidate.name}
								size="compact"
								selected={candidate.id === styleSet.id}
								onClick={() => setStyleSet(candidate)}
							>
								{candidate.name ?? candidate.id}
							</Button>
						))}
					</Inline>
					<ContextStripDiagram
						snapshot={deepBug!}
						styleSet={styleSet}
						selectedPartId="t14-file-reads"
					/>
					<ContextLegend
						items={styleSet.legend}
						styles={styleSet.styles}
						size="sm"
						selectedId="result"
					/>
				</Stack>
			</Panel>
		);
	},
};
