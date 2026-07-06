import type { Meta, StoryObj } from "@storybook/react-vite";
import { useState } from "react";
import { transcriptFixture } from "../../../context";
import { AnnotationRailPanel } from "./AnnotationRailPanel";

const meta = {
	title: "Component Library/Organisms/AnnotationRailPanel",
	component: AnnotationRailPanel,
	args: { annotations: transcriptFixture.annotations },
} satisfies Meta<typeof AnnotationRailPanel>;
export default meta;
type Story = StoryObj<typeof meta>;
export const Rail: Story = {
	render: () => {
		const [selected, setSelected] = useState(transcriptFixture.selectedAnnotationId);
		return (
			<AnnotationRailPanel
				annotations={transcriptFixture.annotations}
				selectedAnnotationId={selected}
				onAnnotationSelect={setSelected}
			/>
		);
	},
};
