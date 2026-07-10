import type { ButtonSize, ButtonVariant } from "../../components/atoms";
import type { CaptionTone, RagStatus } from "../../components/foundation";
import type { ActionSpec } from "./actions";
import type { RenderableValue } from "./core";

export interface DataTableColumnSpec {
	id: string;
	header: RenderableValue;
	cell: CellSpec;
	align?: "start" | "end" | "center";
	maxWidth?: number | string;
}

export type RowKeySpec = string | { field: string } | { template: string };

export type CellSpec =
	| FieldCellSpec
	| NumberCellSpec
	| StatusCellSpec
	| CaptionCellSpec
	| TemplateCellSpec
	| LinkCellSpec
	| LinkButtonCellSpec
	| ActionButtonCellSpec
	| ConstantCellSpec;

export interface FieldCellSpec {
	kind: "field";
	field: string;
	fallback?: string;
}

export interface NumberCellSpec {
	kind: "number";
	field: string;
	format?: "integer" | "fixed";
	digits?: number;
	fallback?: string;
}

export interface StatusCellSpec {
	kind: "status";
	field: string;
	icon?: boolean;
	fallback?: RagStatus | string;
}

export interface CaptionCellSpec {
	kind: "caption";
	field: string;
	tone?: CaptionTone;
	fallback?: string;
}

export interface TemplateCellSpec {
	kind: "template";
	template: string;
}

export interface LinkCellSpec {
	kind: "link";
	hrefField: string;
	labelField: string;
	target?: "_blank" | "_self" | "_parent" | "_top";
	fallbackLabel?: string;
}

export interface LinkButtonCellSpec {
	kind: "linkButton";
	hrefField: string;
	labelField: string;
	variant?: ButtonVariant;
	size?: ButtonSize;
	fallbackLabel?: string;
}

export interface ActionButtonCellSpec {
	kind: "actionButton";
	label: RenderableValue;
	action: ActionSpec;
	variant?: ButtonVariant;
	size?: ButtonSize;
	disabled?: boolean;
}

export interface ConstantCellSpec {
	kind: "constant";
	value: RenderableValue;
}
