import type { Meta, StoryObj } from "@storybook/react-vite";
import { Inline, Stack } from "../../layout";
import { MediaThumb } from "./MediaThumb";

const SKETCH_SRC = "/course-assets/context-window-token-budget.svg";

const meta = {
	title: "Design System/Atoms/MediaThumb",
	component: MediaThumb,
	args: {
		src: SKETCH_SRC,
		alt: "Context window budget sketch",
		style: { width: 160 },
	},
} satisfies Meta<typeof MediaThumb>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};

export const Broken: Story = {
	args: { src: "/media/missing/does-not-exist.png", alt: "Broken source" },
};

export const Empty: Story = {
	args: { src: undefined },
};

export const Contain: Story = {
	args: { fit: "contain" },
};

export const Wide: Story = {
	args: { aspect: "wide", style: { width: 260 } },
};

export const Selected: Story = {
	args: { selected: true },
};

export const Unframed: Story = {
	args: { frame: "none" },
};

const denseSamples = [
	{ id: "sample-1", src: SKETCH_SRC, selected: false },
	{ id: "sample-2", src: SKETCH_SRC, selected: true },
	{ id: "sample-3", src: SKETCH_SRC, selected: false },
	{ id: "sample-4", src: "/media/missing/broken.png", selected: false },
	{ id: "sample-5", src: SKETCH_SRC, selected: false },
	{ id: "sample-6", src: SKETCH_SRC, selected: false },
];

export const DenseGridSample: Story = {
	render: () => (
		<Stack gap="sm">
			<Inline gap="sm">
				{denseSamples.map((sample) => (
					<MediaThumb
						key={sample.id}
						src={sample.src}
						alt={sample.id}
						selected={sample.selected}
						style={{ width: 96 }}
					/>
				))}
			</Inline>
			<Inline gap="sm">
				<MediaThumb src={undefined} style={{ width: 96 }} />
				<MediaThumb src={SKETCH_SRC} alt="Contain sample" fit="contain" style={{ width: 96 }} />
			</Inline>
		</Stack>
	),
};
