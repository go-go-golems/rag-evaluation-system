import type { CSSProperties } from "react";
import { contextVisualStyleToCssVars } from "../context";
import type { JsonObject, StyleBySpec } from "./ir";

/**
 * Resolve a StyleBySpec to CSS variables (the defunctionalized color function):
 * pick the keying value, optionally remap it, look the styleKey up in the
 * styleSet, fall back, and emit `--ctx-*` vars via contextVisualStyleToCssVars.
 * Returns undefined when nothing resolves.
 */
export function resolveStyleByVars(
	spec: StyleBySpec | undefined,
	value: unknown,
	row?: JsonObject,
): CSSProperties | undefined {
	if (!spec) return undefined;
	const raw = spec.field != null && row ? row[spec.field] : value;
	const key0 = raw == null ? "" : String(raw);
	const key = spec.map?.[key0] ?? key0;
	const style =
		spec.styleSet.styles[key] ??
		(spec.fallbackStyleKey ? spec.styleSet.styles[spec.fallbackStyleKey] : undefined) ??
		spec.styleSet.fallbackStyle;
	return style ? contextVisualStyleToCssVars(style) : undefined;
}
