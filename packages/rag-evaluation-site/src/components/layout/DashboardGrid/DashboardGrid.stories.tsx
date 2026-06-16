import type { Meta, StoryObj } from "@storybook/react-vite";
import { Panel } from "../Panel";
import { DashboardGrid } from "./DashboardGrid";

const meta = {
	title: "Design System/Layout/DashboardGrid",
	component: DashboardGrid,
} satisfies Meta<typeof DashboardGrid>;
export default meta;
type Story = StoryObj<typeof meta>;

export const SearchWorkbench: Story = {
	render: () => (
		<DashboardGrid recipe="searchWorkbench">
			<Panel title="Controls" density="condensed">
				Query, retriever, filters
			</Panel>
			<Panel title="Results" density="condensed">
				Ranked evidence table
			</Panel>
			<Panel title="Inspector" density="condensed">
				Selected chunk and document context
			</Panel>
		</DashboardGrid>
	),
};

export const CorpusExplorer: Story = {
	render: () => (
		<DashboardGrid recipe="corpusExplorer">
			<Panel title="Sources" density="condensed">
				Corpus source list
			</Panel>
			<Panel title="Documents" density="condensed">
				Document browser
			</Panel>
			<Panel title="Document Inspector" density="condensed">
				Text, chunks, coverage, artifacts
			</Panel>
		</DashboardGrid>
	),
};
