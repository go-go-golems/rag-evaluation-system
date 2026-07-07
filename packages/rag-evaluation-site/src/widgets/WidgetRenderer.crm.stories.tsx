import type { Meta, StoryObj } from "@storybook/react-vite";
import {
	contactFieldDefs,
	sampleCompanies,
	sampleContacts,
	sampleDeals,
	sampleSalesPipeline,
	sampleStageSummaries,
	sampleUsers,
} from "../crm";
import { defaultWidgetRegistry } from "./defaultRegistry";
import { component, type JsonObject, type WidgetNode } from "./ir";
import { buildRefs, fieldSections, pipelineBoard } from "./presets/crm";
import { WidgetRenderer } from "./WidgetRenderer";

const meta = {
	title: "Widget IR/Renderer/CRM",
	component: WidgetRenderer,
	args: {
		registry: defaultWidgetRegistry,
		onAction: (action, context) => {
			// eslint-disable-next-line no-console
			console.log("[widget action]", action, context);
		},
	},
} satisfies Meta<typeof WidgetRenderer>;

export default meta;
type Story = StoryObj<typeof meta>;

const contact = sampleContacts[0]!;
const refs = buildRefs(sampleUsers, sampleCompanies);

function recordFields(mode: "read" | "edit"): WidgetNode {
	return component("RecordFieldList", {
		values: contact.fields as unknown as JsonObject,
		sections: fieldSections(contactFieldDefs),
		mode,
		refs,
		onFieldChangeAction: { kind: "server", name: "field.update" },
	});
}

/** The pipeline board — the signature CRM screen — rendered from IR. */
export const PipelineBoard: Story = {
	args: {
		node: component("Panel", { title: "Pipeline · Sales", density: "condensed" }, [
			pipelineBoard(sampleSalesPipeline, sampleDeals, { summaries: sampleStageSummaries }),
		]) as WidgetNode,
	},
};

/** The field system in read mode — one FieldRenderer per typed value. */
export const RecordFieldsRead: Story = {
	args: {
		node: component("Panel", { title: contact.name, density: "condensed" }, [
			recordFields("read"),
		]) as WidgetNode,
	},
};

/** The same fields flipped to edit mode — every type swaps to its control. */
export const RecordFieldsEdit: Story = {
	args: {
		node: component("Panel", { title: `Edit · ${contact.name}`, density: "condensed" }, [
			recordFields("edit"),
		]) as WidgetNode,
	},
};
