import type { Meta, StoryObj } from "@storybook/react-vite";
import { AppShell } from "./AppShell";
import { Panel } from "../Panel";

const meta = { title: "Design System/Layout/AppShell", component: AppShell } satisfies Meta<
	typeof AppShell
>;
export default meta;
type Story = StoryObj<typeof meta>;

export const Basic: Story = {
	render: () => (
		<div style={{ height: 360 }}>
			<AppShell
				header={<div style={{ padding: 8 }}>RAG Evaluation</div>}
				sidebar={<div style={{ padding: 8 }}>Navigation</div>}
			>
				<Panel title="Main dashboard">Application content</Panel>
			</AppShell>
		</div>
	),
};
