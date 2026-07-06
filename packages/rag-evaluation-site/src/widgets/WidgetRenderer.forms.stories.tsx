import type { Meta, StoryObj } from "@storybook/react-vite";
import { defaultWidgetRegistry } from "./defaultRegistry";
import { component, text, type WidgetNode } from "./ir";
import { WidgetRenderer } from "./WidgetRenderer";

const meta = {
	title: "Widget IR/Renderer/Forms",
	component: WidgetRenderer,
	args: { registry: defaultWidgetRegistry },
} satisfies Meta<typeof WidgetRenderer>;
export default meta;
type Story = StoryObj<typeof meta>;

function textInput(props: Record<string, unknown>): WidgetNode {
	return component("TextInput", props as never);
}

function textareaInput(props: Record<string, unknown>): WidgetNode {
	return component("TextareaInput", props as never);
}

function selectInput(props: Record<string, unknown>): WidgetNode {
	return component("SelectInput", props as never);
}

function formRow(props: Record<string, unknown>): WidgetNode {
	return component("FormRow", props as never);
}

function displayNameForm(status: "idle" | "saving" | "success" | "error" = "idle"): WidgetNode {
	const invalid = status === "error";
	return component(
		"FormPanel",
		{
			title: "User settings",
			subtitle: "Give this browser a display name for future uploads.",
			formAction: "/settings/profile",
			method: "post",
			submitLabel: "Save settings",
			resetLabel: "Reset",
			status,
			statusMessage:
				status === "success"
					? "Settings saved for this browser."
					: invalid
						? "Fix the highlighted field before saving."
						: undefined,
			footer: "This is browser-local identity, not an authenticated account.",
		} as never,
		[
			formRow({
				label: "Browser identity",
				description: "A read-only system label can appear in the same form anatomy.",
				control: textInput({ name: "sessionLabel", value: "Anonymous visitor", readOnly: true }),
			}),
			formRow({
				label: "Display name",
				required: true,
				hint: "Shown on future uploads from this browser.",
				counter: invalid ? "Needs review" : "6 / 40",
				error: invalid ? "Display name contains unsupported characters." : undefined,
				control: textInput({
					name: "displayName",
					defaultValue: invalid ? "Invalid name" : "Manuel",
					readOnly: false,
					maxLength: 40,
					required: true,
					ariaInvalid: invalid,
				}),
			}),
		],
	);
}

export const DisplayNameSettings: Story = {
	args: { node: displayNameForm("idle") },
};

export const SavingState: Story = {
	args: { node: displayNameForm("saving") },
};

export const SuccessState: Story = {
	args: { node: displayNameForm("success") },
};

export const ErrorState: Story = {
	args: { node: displayNameForm("error") },
};

export const MixedFields: Story = {
	args: {
		node: component(
			"FormPanel",
			{
				title: "Upload defaults",
				subtitle: "Generic Widget IR form with text and select controls.",
				submitLabel: "Save defaults",
				resetLabel: "Reset",
				status: "idle",
			} as never,
			[
				formRow({
					label: "Display name",
					required: true,
					hint: "Editable because readOnly is explicitly false.",
					counter: "6 / 40",
					control: textInput({
						name: "displayName",
						defaultValue: "Manuel",
						readOnly: false,
						maxLength: 40,
					}),
				}),
				formRow({
					label: "Visibility",
					hint: "SelectInput uses an uncontrolled default value for forms.",
					control: selectInput({
						name: "visibility",
						defaultValue: "public",
						options: [
							{ value: "public", label: "Public" },
							{ value: "private", label: "Private" },
						],
					}),
				}),
				formRow({
					orientation: "stacked",
					label: "Course blurb",
					hint: "TextareaInput is for paragraph-length form content.",
					counter: "166 / 1200",
					control: textareaInput({
						name: "blurb",
						defaultValue:
							"Un LLM ne répond pas seulement à votre dernier message : il génère à partir d'une fenêtre remplie de tokens.",
						readOnly: false,
						rows: 4,
						maxLength: 1200,
					}),
				}),
				formRow({
					label: "Status",
					success: "All fields are valid.",
					control: component("StatusText", { status: "succeeded", icon: true } as never, [
						text("Ready to save"),
					]),
				}),
			],
		),
	},
};

export const StackedLongCopy: Story = {
	args: {
		node: component(
			"FormPanel",
			{
				title: "Long copy form",
				subtitle: "This story checks narrow/stacked form rows rendered through Widget IR.",
				submitLabel: "Apply",
				footer:
					"Footer copy remains generic and can explain privacy, persistence, or validation behavior.",
			} as never,
			[
				formRow({
					orientation: "stacked",
					label: "Explanation",
					description:
						"Stacked rows are useful when labels and descriptions need more space than the normal two-column grid.",
					hint: "This helper text intentionally wraps to exercise spacing.",
					counter: "35 / 80",
					control: textInput({
						name: "title",
						defaultValue: "Context workshop browser identity",
						readOnly: false,
					}),
				}),
			],
		),
	},
};
