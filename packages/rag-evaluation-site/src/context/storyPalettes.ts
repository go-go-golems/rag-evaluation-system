import {
  contextCobaltSandStyleSet,
  contextDefaultStyleSet,
  contextSignalOrangeStyleSet,
  contextSlateCoralStyleSet,
} from './fixtures';
import type { ContextStyleSet } from './types';

export type ContextPaletteName = 'Dusty Magenta / Blue' | 'Signal Orange / Cyan' | 'Slate / Coral' | 'Cobalt / Sand';

export const contextPaletteStyleSets: Record<ContextPaletteName, ContextStyleSet> = {
  'Dusty Magenta / Blue': contextDefaultStyleSet,
  'Signal Orange / Cyan': contextSignalOrangeStyleSet,
  'Slate / Coral': contextSlateCoralStyleSet,
  'Cobalt / Sand': contextCobaltSandStyleSet,
};

export const contextPaletteOptions = Object.keys(contextPaletteStyleSets) as ContextPaletteName[];

export function contextStyleSetForPalette(palette: ContextPaletteName) {
  return contextPaletteStyleSets[palette];
}
