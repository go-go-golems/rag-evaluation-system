import type { Meta, StoryObj } from "@storybook/react-vite";
import {
	type ContextPaletteName,
	contextDefaultStyleSet,
	contextPaletteOptions,
	contextStyleSetForPalette,
} from "../../../context";
import { Inline } from "../../layout";
import { AnnotationBadge, type AnnotationBadgeProps } from "./AnnotationBadge";

type PaletteControlsArgs = Omit<AnnotationBadgeProps, "visualStyle"> & {
	palette: ContextPaletteName;
	styleKey: string;
};

const meta = {
	title: "Design System/Atoms/AnnotationBadge",
	component: AnnotationBadge,
	args: { visualStyle: contextDefaultStyleSet.styles.context, label: "Task framing" },
} satisfies Meta<typeof AnnotationBadge>;
export default meta;
type Story = StoryObj<typeof meta>;

export const PaletteControls: StoryObj<PaletteControlsArgs> = {
	args: {
		palette: "Dusty Magenta / Blue",
		styleKey: "context",
		label: "Task framing",
		selected: false,
	},
	argTypes: {
		palette: { control: "select", options: contextPaletteOptions },
		styleKey: { control: "select", options: Object.keys(contextDefaultStyleSet.styles) },
		label: { control: "text" },
		selected: { control: "boolean" },
	},
	render: ({ palette, styleKey, ...args }) => {
		const styleSet = contextStyleSetForPalette(palette);
		return (
			<AnnotationBadge
				{...args}
				visualStyle={styleSet.styles[styleKey] ?? styleSet.styles.other!}
			/>
		);
	},
};

export const Styles: Story = {
	args: { visualStyle: contextDefaultStyleSet.styles.context, label: "Task framing" },
	render: () => (
		<Inline>
			<AnnotationBadge visualStyle={contextDefaultStyleSet.styles.context!} label="Task framing" />
			<AnnotationBadge
				visualStyle={contextDefaultStyleSet.styles.result!}
				label="Tool result cost"
			/>
			<AnnotationBadge visualStyle={contextDefaultStyleSet.styles.generated!} label="Scratchpad" />
			<AnnotationBadge
				visualStyle={contextDefaultStyleSet.styles.active!}
				label="Window rebuilt"
				selected
			/>
		</Inline>
	),
};
