import { type CSSProperties, createContext, type ReactNode, useContext, useMemo } from "react";
import type { ContextPaletteName } from "./storyPalettes";
import {
	cobaltSand,
	defaultContextStyleSet,
	dustyMagentaBlue,
	type PaletteDefinition,
	signalOrangeCyan,
	slateCoral,
	transcriptStyleSet,
} from "./styles";
import type { ContextStyleSet } from "./types";

export const paletteDefinitions: Record<ContextPaletteName, PaletteDefinition> = {
	"Dusty Magenta / Blue": dustyMagentaBlue,
	"Signal Orange / Cyan": signalOrangeCyan,
	"Slate / Coral": slateCoral,
	"Cobalt / Sand": cobaltSand,
};

export const defaultPaletteName: ContextPaletteName = "Dusty Magenta / Blue";

export interface PaletteContextValue {
	paletteName: ContextPaletteName;
	palette: PaletteDefinition;
	contextStyleSet: ContextStyleSet;
	transcriptStyleSet: ContextStyleSet;
}

const defaultPalette = paletteDefinitions[defaultPaletteName];

const PaletteContext = createContext<PaletteContextValue>({
	paletteName: defaultPaletteName,
	palette: defaultPalette,
	contextStyleSet: defaultContextStyleSet(defaultPalette),
	transcriptStyleSet: transcriptStyleSet(defaultPalette),
});

export interface PaletteProviderProps {
	palette?: ContextPaletteName | PaletteDefinition;
	children?: ReactNode;
	className?: string;
	style?: CSSProperties;
}

export function paletteNameForDefinition(
	palette: ContextPaletteName | PaletteDefinition,
): ContextPaletteName {
	if (typeof palette === "string") return palette;
	const found = Object.entries(paletteDefinitions).find(
		([, candidate]) => candidate.name === palette.name,
	);
	return (found?.[0] as ContextPaletteName | undefined) ?? defaultPaletteName;
}

export function paletteDefinition(
	palette: ContextPaletteName | PaletteDefinition = defaultPaletteName,
) {
	return typeof palette === "string" ? paletteDefinitions[palette] : palette;
}

export function paletteCssVars(palette: PaletteDefinition): CSSProperties {
	const colors = palette.colors;
	return {
		"--rag-color-bg": colors.paper,
		"--rag-color-surface": `color-mix(in srgb, ${colors.paper} 90%, #ffffff)`,
		"--rag-color-surface-muted": `color-mix(in srgb, ${colors.paper} 82%, ${colors.grid})`,
		"--rag-color-text": colors.ink,
		"--rag-color-text-muted": colors.shadow,
		"--rag-color-border": colors.grid,
		"--rag-color-border-strong": colors.ink,
		"--rag-color-accent": colors.accent_a,
		"--rag-color-success": colors.accent_b,
		"--rag-color-warning": colors.accent_c,
		"--rag-color-danger": colors.accent_b,
		"--mac-bg-dark": colors.ink,
		"--mac-stripe": colors.ink,
		"--mac-text-inv": colors.paper,
	} as CSSProperties;
}

export function PaletteProvider({
	palette = defaultPaletteName,
	children,
	className,
	style,
}: PaletteProviderProps) {
	const resolvedPalette = paletteDefinition(palette);
	const paletteName = paletteNameForDefinition(palette);
	const value = useMemo<PaletteContextValue>(
		() => ({
			paletteName,
			palette: resolvedPalette,
			contextStyleSet: defaultContextStyleSet(resolvedPalette),
			transcriptStyleSet: transcriptStyleSet(resolvedPalette),
		}),
		[paletteName, resolvedPalette],
	);

	return (
		<PaletteContext.Provider value={value}>
			<div
				className={className}
				data-rag-palette={paletteName}
				style={{ ...paletteCssVars(resolvedPalette), ...style }}
			>
				{children}
			</div>
		</PaletteContext.Provider>
	);
}

export function usePalette() {
	return useContext(PaletteContext);
}

export function useContextStyleSet(explicitStyleSet?: ContextStyleSet) {
	const palette = usePalette();
	return explicitStyleSet ?? palette.contextStyleSet;
}

export function useTranscriptStyleSet(explicitStyleSet?: ContextStyleSet) {
	const palette = usePalette();
	return explicitStyleSet ?? palette.transcriptStyleSet;
}
