import type { CSSProperties } from "react";
import { contextVisualStyleToCssVars } from "../context";
import type { JsonObject, StyleBySpec } from "./ir";
import { resolveStyleByStyle } from "./styleBy.logic";

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
	const style = resolveStyleByStyle(spec, value, row);
	return style ? contextVisualStyleToCssVars(style) : undefined;
}
