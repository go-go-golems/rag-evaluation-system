import type { Meta, StoryObj } from "@storybook/react-vite";
import {
	type ContextPaletteName,
	contextDefaultStyleSet,
	contextPaletteOptions,
	contextStyleSetForPalette,
} from "../../../context";
import { Panel } from "../../layout";
import { ContextStyleSwatch, type ContextStyleSwatchProps } from "./ContextStyleSwatch";

type PaletteControlsArgs = Omit<ContextStyleSwatchProps, "visualStyle"> & {
	palette: ContextPaletteName;
	styleKey: string;
};

const meta = {
	title: "Component Library/Atoms/ContextStyleSwatch",
	component: ContextStyleSwatch,
	args: { visualStyle: contextDefaultStyleSet.styles.result },
} satisfies Meta<typeof ContextStyleSwatch>;
export default meta;
type Story = StoryObj<typeof meta>;

export const PaletteControls: StoryObj<PaletteControlsArgs> = {
	args: {
		palette: "Dusty Magenta / Blue",
		styleKey: "result",
		size: "lg",
		selected: true,
	},
	argTypes: {
		palette: { control: "select", options: contextPaletteOptions },
		styleKey: { control: "select", options: Object.keys(contextDefaultStyleSet.styles) },
		size: { control: "select", options: ["xs", "sm", "md", "lg"] },
		selected: { control: "boolean" },
	},
	render: ({ palette, styleKey, ...args }) => {
		const styleSet = contextStyleSetForPalette(palette);
		return (
			<ContextStyleSwatch
				{...args}
				visualStyle={styleSet.styles[styleKey] ?? styleSet.styles.other!}
			/>
		);
	},
};

export const PatternGallery: Story = {
	render: () => (
		<Panel title="context style swatches">
			<div style={{ display: "flex", gap: 12, flexWrap: "wrap" }}>
				{Object.entries(contextDefaultStyleSet.styles).map(([key, visualStyle]) => (
					<span key={key} title={key}>
						<ContextStyleSwatch visualStyle={visualStyle} size="lg" />
					</span>
				))}
			</div>
		</Panel>
	),
};
