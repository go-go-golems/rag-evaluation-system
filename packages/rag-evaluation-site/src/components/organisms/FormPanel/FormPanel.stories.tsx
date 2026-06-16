import { useState } from "react";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { CheckboxRow, SelectInput, TextInput } from "../../atoms";
import { Caption } from "../../foundation";
import { FormRow, Inline, Stack } from "../../layout";
import { FormPanel, type FormPanelStatus } from "./FormPanel";

const meta = {
	title: "Component Library/Organisms/FormPanel",
	component: FormPanel,
	args: {
		title: "Settings",
		submitLabel: "Save",
	},
} satisfies Meta<typeof FormPanel>;
export default meta;
type Story = StoryObj<typeof meta>;

function DisplayNameRows({ invalid = false }: { invalid?: boolean }) {
	return (
		<Stack gap="sm">
			<FormRow
				label="Browser identity"
				control={<TextInput defaultValue="Anonymous visitor" readOnly />}
				description="This is stored in this browser, not an authenticated account."
			/>
			<FormRow
				label="Display name"
				required
				control={
					<TextInput
						defaultValue={invalid ? "Invalid name" : "Manuel"}
						maxLength={40}
						aria-invalid={invalid ? "true" : undefined}
					/>
				}
				hint="Shown on future uploads from this browser."
				counter={invalid ? "Needs review" : "6 / 40"}
				error={invalid ? "Display name contains unsupported characters." : undefined}
			/>
		</Stack>
	);
}

function InteractivePanel() {
	const [name, setName] = useState("Manuel");
	const [status, setStatus] = useState<FormPanelStatus>("idle");
	const remaining = 40 - name.length;

	return (
		<FormPanel
			title="Interactive settings sketch"
			subtitle="Use this story to review status transitions without a backend."
			status={status}
			statusMessage={status === "idle" ? "Ready for review." : undefined}
			submitLabel="Save settings"
			resetLabel="Reset"
			onSubmit={(event) => {
				event.preventDefault();
				setStatus("saving");
				window.setTimeout(() => setStatus("success"), 600);
			}}
			onReset={() => {
				setName("Anonymous visitor");
				setStatus("idle");
			}}
			footer={<Caption tone="muted">Story-only interaction; no network request is sent.</Caption>}
		>
			<FormRow
				label="Display name"
				required
				control={
					<TextInput
						value={name}
						maxLength={40}
						onChange={(event) => setName(event.currentTarget.value)}
					/>
				}
				hint="Shown on future uploads from this browser."
				counter={`${remaining} remaining`}
				error={remaining < 0 ? "Display name is too long." : undefined}
			/>
		</FormPanel>
	);
}

export const Empty: Story = {
	args: { title: "Empty form panel" },
	render: (args) => <FormPanel {...args} subtitle="Panel chrome and footer without fields." />,
};

export const DisplayNameSettings: Story = {
	render: () => (
		<FormPanel
			title="User settings"
			subtitle="Give this browser a display name for future uploads."
			submitLabel="Save settings"
			resetLabel="Reset"
		>
			<DisplayNameRows />
		</FormPanel>
	),
};

export const Saving: Story = {
	render: () => (
		<FormPanel
			title="User settings"
			subtitle="Saving disables submit actions."
			status="saving"
			submitLabel="Save settings"
			resetLabel="Reset"
		>
			<DisplayNameRows />
		</FormPanel>
	),
};

export const Success: Story = {
	render: () => (
		<FormPanel
			title="User settings"
			status="success"
			statusMessage="Settings saved for this browser."
			submitLabel="Save settings"
			resetLabel="Reset"
		>
			<DisplayNameRows />
		</FormPanel>
	),
};

export const Error: Story = {
	render: () => (
		<FormPanel
			title="User settings"
			status="error"
			statusMessage="Fix the highlighted field before saving."
			submitLabel="Save settings"
			resetLabel="Reset"
		>
			<DisplayNameRows invalid />
		</FormPanel>
	),
};

export const MultipleFields: Story = {
	render: () => (
		<FormPanel
			title="Upload defaults"
			subtitle="A generic multi-field form, not MINIVIZ-specific."
			submitLabel="Save defaults"
			resetLabel="Reset"
		>
			<Stack gap="sm">
				<FormRow
					label="Display name"
					control={<TextInput defaultValue="Manuel" />}
					hint="Text input."
					counter="6 / 40"
				/>
				<FormRow
					label="Visibility"
					control={
						<SelectInput defaultValue="public">
							<option value="public">Public</option>
							<option value="private">Private</option>
						</SelectInput>
					}
					hint="Select input."
				/>
				<FormRow
					label="Options"
					control={
						<Stack gap="xs">
							<CheckboxRow checked={true} onChange={() => {}}>
								Show display name on uploads
							</CheckboxRow>
							<CheckboxRow checked={false} onChange={() => {}}>
								Hide uploads from public list by default
							</CheckboxRow>
						</Stack>
					}
					hint="Checkbox rows in the control slot."
				/>
			</Stack>
		</FormPanel>
	),
};

export const LongCopyAndActions: Story = {
	render: () => (
		<FormPanel
			title="Long copy review"
			subtitle="This story checks wrapping behavior for longer explanatory text in a narrow form panel. The panel should preserve the Classic Mac density without hiding the helper copy."
			actions={
				<Inline gap="xs">
					<Caption tone="accent">Preview</Caption>
				</Inline>
			}
			submitLabel="Apply changes"
			resetLabel="Discard"
			footer="Footer copy can explain privacy, persistence, or browser-local behavior without becoming a product-specific component."
		>
			<FormRow
				orientation="stacked"
				label="Explanation"
				control={<TextInput defaultValue="Context workshop browser identity" />}
				description="Stacked rows are useful when labels and descriptions need more room than the normal two-column grid."
				hint="This helper text intentionally wraps to exercise spacing."
				counter="35 / 80"
			/>
		</FormPanel>
	),
};

export const Interactive: Story = { render: () => <InteractivePanel /> };
