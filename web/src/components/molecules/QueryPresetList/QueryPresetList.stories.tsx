import type { Meta, StoryObj } from "@storybook/react-vite";
import { QueryPresetList } from "./QueryPresetList";

const meta = {
	title: "Component Library/Molecules/QueryPresetList",
	component: QueryPresetList,
} satisfies Meta<typeof QueryPresetList>;

export default meta;
type Story = StoryObj;

export const TTCQueries: Story = {
	render: () => (
		<div style={{ width: 280 }}>
			<QueryPresetList
				queries={[
					"crape myrtle varieties",
					"how to plant arborvitae",
					"hydrangea pruning",
					"fast growing trees for privacy screen",
				]}
				onSelect={() => undefined}
			/>
		</div>
	),
};
