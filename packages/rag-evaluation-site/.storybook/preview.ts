import { createElement } from "react";
import type { Preview } from "@storybook/react-vite";
import { PaletteProvider, contextPaletteOptions, type ContextPaletteName } from "../src/context";
import "../src/styles.css";

const preview: Preview = {
	parameters: {
		layout: "padded",
	},
	globalTypes: {
		palette: {
			description: "Design-system palette",
			defaultValue: "Dusty Magenta / Blue",
			toolbar: {
				icon: "paintbrush",
				items: contextPaletteOptions,
				dynamicTitle: true,
			},
		},
	},
	decorators: [
		(Story, context) =>
			createElement(
				PaletteProvider,
				{ palette: context.globals.palette as ContextPaletteName },
				createElement(Story),
			),
	],
};

export default preview;
