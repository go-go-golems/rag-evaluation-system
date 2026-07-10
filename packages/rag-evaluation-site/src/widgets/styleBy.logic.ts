import type { ContextVisualStyle } from "../context/types";
import type { JsonObject, StyleBySpec } from "./ir";

export function resolveStyleByStyle(
	spec: StyleBySpec | undefined,
	value: unknown,
	row?: JsonObject,
): ContextVisualStyle | undefined {
	if (!spec) return undefined;
	const raw = spec.field != null && row ? row[spec.field] : value;
	const key0 = raw == null ? "" : String(raw);
	const key = spec.map?.[key0] ?? key0;
	return (
		spec.styleSet.styles[key] ??
		(spec.fallbackStyleKey ? spec.styleSet.styles[spec.fallbackStyleKey] : undefined) ??
		spec.styleSet.fallbackStyle
	);
}
