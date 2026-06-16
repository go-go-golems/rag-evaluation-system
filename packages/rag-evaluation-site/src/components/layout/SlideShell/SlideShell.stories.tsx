import type { Meta, StoryObj } from "@storybook/react-vite";
import {
	contextDefaultStyleSet,
	contextPaletteOptions,
	contextSlides,
	contextStyleSetForPalette,
	contextWindowSnapshots,
	type ContextPaletteName,
} from "../../../context";
import { ContextStackDiagram, FigureBlock, KeyPointList } from "../../molecules";
import { SlideShell, type SlideShellProps } from "./SlideShell";

const slide = contextSlides[1]!;
const snapshot = contextWindowSnapshots.find((item) => item.id === slide.snapshotId)!;

type PaletteControlsArgs = SlideShellProps & { palette: ContextPaletteName };

function figureForPalette(palette: ContextPaletteName) {
	const styleSet = contextStyleSetForPalette(palette);
	return (
		<FigureBlock caption="context window — 200,000 tokens" frame="none">
			<ContextStackDiagram snapshot={snapshot} styleSet={styleSet} />
		</FigureBlock>
	);
}

const figure = (
	<FigureBlock caption="context window — 200,000 tokens" frame="none">
		<ContextStackDiagram snapshot={snapshot} styleSet={contextDefaultStyleSet} />
	</FigureBlock>
);

const points = <KeyPointList items={slide.notes} />;

const meta = {
	title: "Design System/Layout/SlideShell",
	component: SlideShell,
	args: {
		eyebrow: slide.kicker,
		counter: "02 / 06",
		title: slide.title,
		primary: figure,
		secondary: points,
	},
} satisfies Meta<typeof SlideShell>;

export default meta;
type Story = StoryObj<typeof meta>;

export const PaletteControls: StoryObj<PaletteControlsArgs> = {
	args: {
		eyebrow: slide.kicker,
		counter: "02 / 06",
		title: slide.title,
		palette: "Dusty Magenta / Blue",
		primarySide: "right",
		ratio: "primaryWide",
	},
	argTypes: {
		palette: { control: "select", options: contextPaletteOptions },
		primarySide: { control: "select", options: ["left", "right"] },
		ratio: { control: "select", options: ["equal", "primaryWide", "secondaryWide"] },
		primary: { control: false },
		secondary: { control: false },
	},
	render: ({ palette, ...args }) => (
		<SlideShell {...args} primary={figureForPalette(palette)} secondary={points} />
	),
};

export const VisualFirst: Story = {};

export const ContentFirst: Story = {
	args: {
		primary: points,
		secondary: figure,
		primarySide: "left",
		ratio: "secondaryWide",
	},
};

export const VisualOnRight: Story = {
	args: {
		primary: figure,
		secondary: points,
		primarySide: "right",
		ratio: "primaryWide",
	},
};
