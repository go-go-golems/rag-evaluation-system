import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import {
	type ContextPaletteName,
	contextDefaultStyleSet,
	contextPaletteOptions,
	contextSlides,
	contextStyleSetForPalette,
	contextWindowSnapshots,
} from "../../../context";
import { Caption } from "../../foundation";
import { SlideShell } from "../../layout";
import { ContextStackDiagram, FigureBlock, KeyPointList } from "../../molecules";
import { CourseStudioShell, type CourseStudioShellProps } from "./CourseStudioShell";
import { courseStudioNavSections } from "./courseStudioNav";

const slide = contextSlides[1]!;
const snapshot = contextWindowSnapshots.find((item) => item.id === slide.snapshotId)!;

type PaletteControlsArgs = CourseStudioShellProps & { palette: ContextPaletteName };

function slideContentForPalette(palette: ContextPaletteName) {
	const styleSet = contextStyleSetForPalette(palette);
	return (
		<SlideShell
			eyebrow={slide.kicker}
			counter="02 / 06"
			title={slide.title}
			primary={
				<FigureBlock caption="context window — 200,000 tokens">
					<ContextStackDiagram snapshot={snapshot} styleSet={styleSet} />
				</FigureBlock>
			}
			secondary={<KeyPointList items={slide.notes} />}
			primarySide="right"
		/>
	);
}

const slideContent = (
	<SlideShell
		eyebrow={slide.kicker}
		counter="02 / 06"
		title={slide.title}
		primary={
			<FigureBlock caption="context window — 200,000 tokens">
				<ContextStackDiagram snapshot={snapshot} styleSet={contextDefaultStyleSet} />
			</FigureBlock>
		}
		secondary={<KeyPointList items={slide.notes} />}
		primarySide="right"
	/>
);

const meta = {
	title: "Component Library/Organisms/CourseStudioShell",
	component: CourseStudioShell,
	args: {
		sections: courseStudioNavSections,
		activeItemId: "slides",
		children: slideContent,
		sidebarFooter: <Caption>Workshop draft · v0.4</Caption>,
	},
} satisfies Meta<typeof CourseStudioShell>;
export default meta;
type Story = StoryObj<typeof meta>;

export const PaletteControls: StoryObj<PaletteControlsArgs> = {
	args: {
		sections: courseStudioNavSections,
		activeItemId: "slides",
		palette: "Dusty Magenta / Blue",
		sidebarFooter: <Caption>Workshop draft · v0.4</Caption>,
	},
	argTypes: {
		palette: { control: "select", options: contextPaletteOptions },
		sections: { control: false },
		activeItemId: { control: "text" },
		children: { control: false },
		sidebarFooter: { control: false },
	},
	render: ({ palette, ...args }) => (
		<CourseStudioShell {...args}>{slideContentForPalette(palette)}</CourseStudioShell>
	),
};

export const SlidesActive: Story = {};

export const Interactive: Story = {
	render: () => {
		const [active, setActive] = useState("slides");
		return (
			<CourseStudioShell
				sections={courseStudioNavSections}
				activeItemId={active}
				onNavigate={setActive}
				sidebarFooter={<Caption>Selected: {active}</Caption>}
			>
				{slideContent}
			</CourseStudioShell>
		);
	},
};
