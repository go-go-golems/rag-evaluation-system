import {
	contextCobaltSandStyleSet,
	contextDefaultStyleSet,
	contextSignalOrangeStyleSet,
	contextSlateCoralStyleSet,
	transcriptCobaltSandStyleSet,
	transcriptDefaultStyleSet,
	transcriptSignalOrangeStyleSet,
	transcriptSlateCoralStyleSet,
} from "./fixtures";
import type { ContextStyleSet } from "./types";

export type ContextPaletteName =
	| "Dusty Magenta / Blue"
	| "Signal Orange / Cyan"
	| "Slate / Coral"
	| "Cobalt / Sand";

export const contextPaletteStyleSets: Record<ContextPaletteName, ContextStyleSet> = {
	"Dusty Magenta / Blue": contextDefaultStyleSet,
	"Signal Orange / Cyan": contextSignalOrangeStyleSet,
	"Slate / Coral": contextSlateCoralStyleSet,
	"Cobalt / Sand": contextCobaltSandStyleSet,
};

export const transcriptPaletteStyleSets: Record<ContextPaletteName, ContextStyleSet> = {
	"Dusty Magenta / Blue": transcriptDefaultStyleSet,
	"Signal Orange / Cyan": transcriptSignalOrangeStyleSet,
	"Slate / Coral": transcriptSlateCoralStyleSet,
	"Cobalt / Sand": transcriptCobaltSandStyleSet,
};

export const contextPaletteOptions = Object.keys(contextPaletteStyleSets) as ContextPaletteName[];

export function contextStyleSetForPalette(palette: ContextPaletteName) {
	return contextPaletteStyleSets[palette];
}

export function transcriptStyleSetForPalette(palette: ContextPaletteName) {
	return transcriptPaletteStyleSets[palette];
}
