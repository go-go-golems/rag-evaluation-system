import type { StorybookConfig } from "@storybook/react-vite";

const config: StorybookConfig = {
	stories: ["../src/**/*.stories.@(ts|tsx)"],
	staticDirs: ["./static"],
	addons: [],
	framework: {
		name: "@storybook/react-vite",
		options: {},
	},
	viteFinal: (config) => {
		config.css = {
			...config.css,
			modules: {
				// Readable class names in Storybook: Button_root, Button_normal
				generateScopedName: "[name]_[local]",
			},
		};
		return config;
	},
};

export default config;
