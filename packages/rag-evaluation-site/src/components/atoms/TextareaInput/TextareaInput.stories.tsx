import { useState } from "react";
import type { Meta, StoryObj } from "@storybook/react-vite";
import { Caption, StatusText } from "../../foundation";
import { FormRow, Panel, Stack } from "../../layout";
import { TextareaInput } from "./TextareaInput";

const meta = {
	title: "Design System/Atoms/TextareaInput",
	component: TextareaInput,
} satisfies Meta<typeof TextareaInput>;
export default meta;
type Story = StoryObj<typeof meta>;

function ControlledExample() {
	const [value, setValue] = useState(
		"Moins de prompt magique, plus de contexte choisi : apprendre à composer l'entrée du modèle comme une entrée de programme.",
	);
	const remaining = 500 - value.length;

	return (
		<Stack gap="xs">
			<TextareaInput
				aria-describedby="tagline-counter"
				maxLength={500}
				rows={3}
				value={value}
				onChange={(event) => setValue(event.currentTarget.value)}
			/>
			<Caption id="tagline-counter" tone={remaining < 40 ? "warning" : "muted"}>
				{remaining} characters remaining
			</Caption>
		</Stack>
	);
}

export const Examples: Story = {
	render: () => (
		<Stack gap="sm">
			<TextareaInput rows={3} defaultValue="A short multiline value for a course tagline." />
			<TextareaInput rows={5} placeholder="Write a compact course blurb…" />
			<TextareaInput rows={3} defaultValue="Read-only explanation" readOnly />
		</Stack>
	),
};

export const States: Story = {
	render: () => (
		<Panel title="Textarea states" density="condensed">
			<Stack gap="sm">
				<FormRow
					label="Tagline"
					control={
						<TextareaInput rows={3} defaultValue="A longer sentence that benefits from wrapping." />
					}
					hint="Use for sentence-length metadata."
				/>
				<FormRow
					label="Blurb"
					control={<TextareaInput rows={6} placeholder="Describe the page…" />}
					hint="Use for paragraph-length copy."
				/>
				<FormRow
					label="Disabled"
					control={<TextareaInput rows={3} defaultValue="Locked by system" disabled />}
					hint="Native disabled state; not submitted by forms."
				/>
				<FormRow
					label="Invalid"
					control={<TextareaInput rows={3} defaultValue="Too short" aria-invalid="true" />}
					hint={
						<StatusText status="error" icon>
							Description must explain what the learner will do.
						</StatusText>
					}
				/>
			</Stack>
		</Panel>
	),
};

export const ControlledWithCounter: Story = {
	render: () => (
		<Panel title="Course metadata field" density="condensed">
			<FormRow
				orientation="stacked"
				label="Tagline"
				control={<ControlledExample />}
				hint="TextareaInput keeps longer metadata readable without creating page-specific controls."
			/>
		</Panel>
	),
};
