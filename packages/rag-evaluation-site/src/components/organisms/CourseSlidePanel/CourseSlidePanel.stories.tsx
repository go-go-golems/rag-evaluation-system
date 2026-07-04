import type { Meta, StoryObj } from "@storybook/react-vite";
import {
	type ContextPaletteName,
	contextPaletteOptions,
	contextSlides,
	contextStyleSetForPalette,
	contextWindowSnapshots,
} from "../../../context";
import { Stack } from "../../layout";
import { CourseSlidePanel, type CourseSlidePanelProps } from "./CourseSlidePanel";

const slide = contextSlides[1]!;
const snapshot = contextWindowSnapshots.find((item) => item.id === slide.snapshotId)!;
type PaletteControlsArgs = Omit<CourseSlidePanelProps, "styleSet"> & {
	palette: ContextPaletteName;
};

const meta = {
	title: "Component Library/Organisms/CourseSlidePanel",
	component: CourseSlidePanel,
	args: { slide, snapshot, index: 1, total: contextSlides.length },
} satisfies Meta<typeof CourseSlidePanel>;
export default meta;
type Story = StoryObj<typeof meta>;

export const PaletteControls: StoryObj<PaletteControlsArgs> = {
	args: {
		slide,
		snapshot,
		index: 1,
		total: contextSlides.length,
		visualSide: "left",
		palette: "Dusty Magenta / Blue",
	},
	argTypes: {
		palette: { control: "select", options: contextPaletteOptions },
		slide: { control: false },
		snapshot: { control: false },
		visualSide: { control: "select", options: ["left", "right"] },
	},
	render: ({ palette, ...args }) => (
		<CourseSlidePanel {...args} styleSet={contextStyleSetForPalette(palette)} />
	),
};

export const AnatomySlide: Story = {
	render: () => (
		<CourseSlidePanel slide={slide} snapshot={snapshot} index={1} total={contextSlides.length} />
	),
};
export const VisualOnRight: Story = {
	render: () => (
		<CourseSlidePanel
			slide={slide}
			snapshot={snapshot}
			index={1}
			total={contextSlides.length}
			visualSide="right"
		/>
	),
};
export const SlideDeckSamples: Story = {
	render: () => (
		<Stack gap="md">
			{contextSlides.slice(0, 3).map((s, i) => (
				<CourseSlidePanel
					key={s.id}
					slide={s}
					snapshot={contextWindowSnapshots.find((item) => item.id === s.snapshotId)!}
					index={i}
					total={contextSlides.length}
					visualSide={i % 2 === 0 ? "left" : "right"}
				/>
			))}
		</Stack>
	),
};
