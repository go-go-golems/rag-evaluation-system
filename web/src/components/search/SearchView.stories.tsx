import type { Meta, StoryObj } from "@storybook/react-vite";
import { MockApiProvider } from "../../storybook/MockApiProvider";
import { SearchView } from "./SearchView";

const meta = {
	title: "Views/Search/SearchView",
	component: SearchView,
	decorators: [
		(Story) => (
			<MockApiProvider>
				<Story />
			</MockApiProvider>
		),
	],
	parameters: { layout: "fullscreen" },
} satisfies Meta<typeof SearchView>;

export default meta;
type Story = StoryObj<typeof meta>;

export const Default: Story = {};
