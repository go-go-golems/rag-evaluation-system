import { useState } from "react";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { FormRow, Panel, Stack } from "../../layout";
import { Caption } from "../../foundation";
import { CheckboxRow } from "./CheckboxRow";

const meta = { title: "Design System/Atoms/CheckboxRow", component: CheckboxRow } satisfies Meta<
	typeof CheckboxRow
>;
export default meta;
type Story = StoryObj<typeof meta>;

function ControlledCheckbox() {
	const [checked, setChecked] = useState(true);

	return (
		<Stack gap="xs">
			<CheckboxRow checked={checked} onChange={(event) => setChecked(event.currentTarget.checked)}>
				Include uploads from this browser
			</CheckboxRow>
			<Caption tone="muted">Current value: {checked ? "checked" : "unchecked"}</Caption>
		</Stack>
	);
}

function Section({ title, children }: { title: string; children: React.ReactNode }) {
	return (
		<Panel title={title} density="condensed">
			<Stack gap="sm">{children}</Stack>
		</Panel>
	);
}

export const Checked: Story = {
	args: { checked: true, onChange: () => {}, children: "TTC Guides" },
};
export const Unchecked: Story = {
	args: { checked: false, onChange: () => {}, children: "TTC Articles" },
};

export const States: Story = {
	args: { checked: false, onChange: () => {} },
	render: () => (
		<Stack gap="md">
			<Section title="Common states">
				<CheckboxRow checked={true} onChange={() => {}}>
					Checked option
				</CheckboxRow>
				<CheckboxRow checked={false} onChange={() => {}}>
					Unchecked option
				</CheckboxRow>
				<CheckboxRow checked={true} onChange={() => {}} disabled>
					Disabled checked option
				</CheckboxRow>
				<CheckboxRow checked={false} onChange={() => {}} disabled>
					Disabled unchecked option
				</CheckboxRow>
				<ControlledCheckbox />
			</Section>
			<Section title="Long label wrapping">
				<CheckboxRow checked={true} onChange={() => {}}>
					Allow this browser-local display name to be shown next to transcript uploads created
					during the workshop. This copy intentionally wraps so spacing and line-height are visible.
				</CheckboxRow>
			</Section>
		</Stack>
	),
};

export const InFormRows: Story = {
	args: { checked: false, onChange: () => {} },
	render: () => (
		<Panel title="Settings checkbox sketch" density="condensed">
			<Stack gap="sm">
				<FormRow
					label="Upload label"
					control={
						<CheckboxRow checked={true} onChange={() => {}}>
							Show my display name on uploads from this browser
						</CheckboxRow>
					}
					hint="Example checkbox in existing form-row anatomy."
				/>
				<FormRow
					label="Visibility"
					control={
						<CheckboxRow checked={false} onChange={() => {}}>
							Hide uploads from the public list by default
						</CheckboxRow>
					}
					hint="Future setting example; no product behavior is wired here."
				/>
			</Stack>
		</Panel>
	),
};
