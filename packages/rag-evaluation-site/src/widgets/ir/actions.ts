import type { JsonObject, JsonValue } from "./core";

export type ActionSpec =
	| NavigateActionSpec
	| DownloadActionSpec
	| ServerActionSpec
	| EventActionSpec
	| CopyActionSpec;

export interface ActionSpecBase {
	/**
	 * Optional confirmation prompt shown before the action dispatches.
	 * String prompts support `${path}` / `$name` interpolation against the action context.
	 * Template prompts are the v2 data form produced by typed builders.
	 * Applies to every action kind; handled centrally in dispatchWidgetAction.
	 */
	confirm?: string | TemplateSpec;
}

export interface TemplateSpec {
	kind?: "template";
	parts: TemplatePartSpec[];
}

export type TemplatePartSpec = TemplateTextPart | TemplatePathPart | TemplateLiteralPart;

export interface TemplateTextPart {
	kind: "text";
	text: string;
}

export interface TemplatePathPart {
	kind: "path";
	path: string;
}

export interface TemplateLiteralPart {
	kind: "literal";
	value: JsonValue;
}

export interface PayloadTemplateSpec {
	kind: "payloadTemplate";
	fields: Record<string, TemplatePartSpec | JsonValue>;
}

export interface NavigateActionSpec extends ActionSpecBase {
	kind: "navigate";
	to: string;
	params?: JsonObject;
}

export interface DownloadActionSpec extends ActionSpecBase {
	kind: "download";
	to: string;
	params?: JsonObject;
}

export interface ServerActionSpec extends ActionSpecBase {
	kind: "server";
	name: string;
	payload?: JsonObject | PayloadTemplateSpec;
}

export interface EventActionSpec extends ActionSpecBase {
	kind: "event";
	event: string;
	detail?: JsonObject;
}

export interface CopyActionSpec extends ActionSpecBase {
	kind: "copy";
	value?: string;
	field?: string;
}
